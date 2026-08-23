import {
  AutomationDependencyKind,
  Command,
  type AutomationDependency,
  type AutomationOrchestratorDefineResponse,
} from '../../../../shared/contract';
import { utworzRameOkna } from '../../komponenty/rama-okna';
import {
  pobierzPlik,
  pole,
  przelacznikWidoku,
  przestaw,
  przyciskAkcji as przycisk,
  wiersz,
  wybor,
} from '../../modele/kontrolki-formularza';
import type { PokrycieKomend } from '../pokrycie-komend';
import { podpisKrawedzi, zapisDot, zapisMermaid } from './eksport-grafu';
import { utworzGrafKrokow, type OpisGrafu } from './graf-krokow';
import { wykazPrzejsc, zdanieOMaszynieStanow } from './maszyna-stanow';
import { utworzPanelDopelnienUkladu } from './panel-dopelnien-ukladu';
import { utworzPanelZaleznosci } from './panel-zaleznosci';
import type { StanAutomatyki } from './stan-automatyki';
import { utworzStanTresci } from './stany-okna';
import {
  listaZaleznosci,
  ocenaUkladu,
  opisSciezkiKrytycznej,
  opisZaleznosci,
  pustyUklad,
  zastrzezeniaUkladu,
  zawiera,
} from './widok-ukladu';
import type { ZrodloAutomations } from './zrodlo-automations';

/**
 * Orchestrator — okno kreatora modułu Automations: ustalenie zależności między
 * krokami, zmiana ich kolejności, walidacja układu, podświetlenie ścieżki
 * krytycznej i eksport mapy zależności.
 *
 * Walidacja nie jest bramą — rdzeń zapisuje układ także wtedy, gdy ma cykl,
 * i oddaje zastrzeżenia, które okno pokazuje wprost. Zastrzeżenie do układu nie
 * jest niepowodzeniem czynności: werdykt o czynności należy do odpowiedzi na
 * wywołanie, werdykt o układzie stoi osobno — w ocenie układu i w wykazie
 * zastrzeżeń.
 *
 * Pole `criticalPathStepIds` liczy rdzeń, który zna cały układ; drugi rachunek
 * w oknie byłby drugą prawdą o tym samym grafie.
 *
 * Układ pokazują trzy widoki tej samej treści: kanwa grafu (`graf-krokow.ts`),
 * wykaz zależności i model stanów przebiegu (`maszyna-stanow.ts`). Kanwa mówi,
 * jak długi jest łańcuch i gdzie tory się rozchodzą, wykaz — co z czym jest
 * związane i którą zależność usunąć, model — co przebiegowi wolno dalej.
 */
export interface OknoOrchestratora {
  element: HTMLElement;
  odswiez(): void;
  /** Zamyka nasłuch zmiany układu zależności. */
  zamknij(): void;
}

/** Rodzaje zależności — grupowanie sekwencyjne, równoległe i warunkowe. */
const RODZAJE_ZALEZNOSCI: ReadonlyArray<[string, string]> = [
  [AutomationDependencyKind.Sequential, 'sekwencyjna (krok po kroku)'],
  [AutomationDependencyKind.Parallel, 'równoległa (tor obok toru)'],
  [AutomationDependencyKind.Conditional, 'warunkowa (przejście pod warunkiem)'],
];

export function utworzOknoOrchestratora(
  zrodlo: ZrodloAutomations,
  stan: StanAutomatyki,
  pokrycie: PokrycieKomend,
): OknoOrchestratora {
  const rama = utworzRameOkna({
    tytul: 'Orchestrator',
    rola: 'kreator',
    przeznaczenie:
      'Układ zależności między krokami automatyki wraz z jego walidacją i ścieżką krytyczną.',
    przedrostek: 'da',
  });
  const tresc = utworzStanTresci();

  // Rodzina `orchestration.*` obok `automation.orchestrator.define`: odczyt
  // układu bez zapisu, zmiana pojedynczej krawędzi i sprawdzenie bez dotykania
  // grafu. Panel stoi w narzędziach kontekstowych, bo zawęża i bada ten sam
  // układ, którym gospodaruje panel akcji okna.
  const zaleznosci = utworzPanelZaleznosci(zrodlo, () => stan.automatyka());
  // Cztery dopełnienia układu, których krawędź nie wyraża: bramka dołączenia,
  // grupa kroków, krok wycofujący i spięcie kolejek z silnikiem środowiska
  // MultitaskingAI. Panel stoi obok panelu zależności, bo gospodaruje tym samym
  // układem — innym jego wymiarem.
  const dopelnienia = utworzPanelDopelnienUkladu(zrodlo, () => stan.automatyka());
  rama.narzedzia.append(zaleznosci.element, dopelnienia.element);

  // Kanwa powstaje raz i przeżywa przerysowania treści: `tresc.tresc()` czyści
  // miejsce, lecz dokładany element jest wciąż ten sam, więc powiększenie
  // ustawione przez Operatora nie wraca do wartości wyjściowej po każdym
  // odczycie układu.
  const graf = utworzGrafKrokow();

  const powierzchnia = zlozPowierzchnie(rama, tresc.element, pokrycie);
  const { krokZ, krokDo, rodzaj, warunek, dodaj, zapisz, waliduj } = powierzchnia;

  let uklad: AutomationDependency[] = [];
  /**
   * Ostatnia odpowiedź rdzenia o układzie. Trzymamy ją w całości, bo widoki
   * przełączają się bez pytania rdzenia, a ocena, zastrzeżenia i ścieżka
   * krytyczna należą do niego — okno ma je powtórzyć, nie wyliczyć ponownie.
   */
  let ostatniOdczyt: AutomationOrchestratorDefineResponse = { dependencies: [], valid: true };

  /**
   * Układ w postaci, w której czyta go kanwa i eksporty.
   *
   * Węzły biorą nazwy z definicji automatyki, gdy ta jest w module odczytana;
   * zależność potrafi jednak wskazać krok, którego w odczytanej definicji nie
   * ma — wtedy węzeł zostaje przy samym identyfikatorze, zamiast zniknąć
   * z rysunku i zerwać krawędź.
   */
  function opisUkladu(): OpisGrafu {
    const nazwy = new Map<string, string>();
    for (const krok of stan.definicja()?.steps ?? []) nazwy.set(krok.id, krok.name ?? '');
    const kody = new Set<string>([...nazwy.keys()]);
    for (const zaleznosc of uklad) {
      kody.add(zaleznosc.fromStepId);
      kody.add(zaleznosc.toStepId);
    }
    return {
      wezly: [...kody].map((kod) => ({ kod, nazwa: nazwy.get(kod) ?? '' })),
      krawedzie: uklad.map((zaleznosc) => ({
        odKroku: zaleznosc.fromStepId,
        doKroku: zaleznosc.toStepId,
        podpis: podpisKrawedzi(zaleznosc.kind, zaleznosc.condition),
      })),
      sciezkaKrytyczna: new Set(ostatniOdczyt.criticalPathStepIds ?? []),
    };
  }

  function pokaz(wynik: AutomationOrchestratorDefineResponse): void {
    uklad = [...wynik.dependencies];
    ostatniOdczyt = wynik;
    const miejsce = tresc.tresc();
    const sciezka = new Set(wynik.criticalPathStepIds ?? []);

    miejsce.append(ocenaUkladu(wynik));
    const zastrzezenia = zastrzezeniaUkladu(wynik.issues ?? []);
    if (zastrzezenia !== null) miejsce.append(zastrzezenia);
    if (uklad.length === 0) miejsce.append(pustyUklad());
    if (powierzchnia.kanwa.dataset['wlaczony'] === 'true') {
      graf.pokaz(opisUkladu());
      miejsce.append(graf.element);
    }
    miejsce.append(listaZaleznosci(uklad, sciezka, usunZaleznosc));
    miejsce.append(opisSciezkiKrytycznej(sciezka));
    if (powierzchnia.stany.dataset['wlaczony'] === 'true') {
      miejsce.append(zdanieOStanach(), wykazPrzejsc(''));
    }
  }

  /**
   * Wysyła układ i opisuje to, co oddał rdzeń.
   *
   * `czekana` to zależność, której obecności w odpowiedzi należy oczekiwać po
   * dodaniu — bez niej zdanie sukcesu mówiłoby o zamiarze okna, a nie o skutku
   * po stronie rdzenia.
   */
  function wyslij(
    zaleznosci: AutomationDependency[] | undefined,
    czynnosc: string,
    czekana?: AutomationDependency,
  ): void {
    const automatyka = stan.automatyka();
    if (automatyka === '') {
      tresc.pusto('Wskaż automatykę w Workflow Builderze — układ dotyczy jej kroków.');
      return;
    }
    tresc.ladowanie('Praca nad układem zależności…');
    const zadanie: Parameters<ZrodloAutomations['ustawZaleznosci']>[0] = { workflowId: automatyka };
    if (zaleznosci !== undefined) zadanie.dependencies = zaleznosci;
    void zrodlo.ustawZaleznosci(zadanie).then((wynik) => {
      if (!wynik.udany || wynik.wynik === undefined) {
        tresc.blad('Rdzeń nie oddał układu zależności.', wynik.blad);
        return;
      }
      pokaz(wynik.wynik);
      const oddane = wynik.wynik.dependencies;
      if (czekana !== undefined && !zawiera(oddane, czekana)) {
        tresc.potwierdzenie(
          `Rdzeń oddał układ ${oddane.length} zależności, lecz ${opisZaleznosci(czekana)} ` +
            'nie ma wśród nich — czynność się nie odbyła.',
          false,
        );
        return;
      }
      tresc.potwierdzenie(`${czynnosc} Rdzeń oddał układ ${oddane.length} zależności.`, true);
    });
  }

  /**
   * Usunięcie pojedynczej zależności komendą `orchestration.dependency.remove`.
   *
   * Komenda oddaje pole `removed` oraz `dependencies` po usunięciu. Pełny obraz
   * układu — ocenę i ścieżkę krytyczną — odczytuje się po niej ponownym
   * `automation.orchestrator.define` bez zmiany. Gdy `removed` jest fałszem,
   * czynność się nie odbyła i okno tego nie zataja.
   */
  function usunZaleznosc(zaleznosc: AutomationDependency): void {
    const automatyka = stan.automatyka();
    if (automatyka === '') {
      tresc.pusto('Wskaż automatykę w Workflow Builderze — usunięcie dotyczy jej układu.');
      return;
    }
    tresc.ladowanie('Usuwanie zależności…');
    void zrodlo
      .usunZaleznosc({
        workflowId: automatyka,
        fromStepId: zaleznosc.fromStepId,
        toStepId: zaleznosc.toStepId,
      })
      .then((wynik) => {
        if (!wynik.udany || wynik.wynik === undefined) {
          tresc.blad('Rdzeń nie usunął zależności.', wynik.blad);
          return;
        }
        if (!wynik.wynik.removed) {
          tresc.potwierdzenie(
            `Rdzeń oddał układ ${wynik.wynik.dependencies.length} zależności, lecz ` +
              `${opisZaleznosci(zaleznosc)} nie usunął — czynność się nie odbyła.`,
            false,
          );
          uklad = [...wynik.wynik.dependencies];
          return;
        }
        // Pełny obraz (ocena, zastrzeżenia, ścieżka krytyczna) czytamy ponownym
        // odczytem układu — komenda usunięcia oddaje same zależności.
        wyslij(undefined, `Zależność ${opisZaleznosci(zaleznosc)} usunięta.`);
      });
  }

  dodaj.addEventListener('click', () => {
    if (krokZ.value.trim() === '' || krokDo.value.trim() === '') {
      tresc.potwierdzenie('Zależność wymaga wskazania obu kroków.', false);
      return;
    }
    const zaleznosc: AutomationDependency = {
      fromStepId: krokZ.value.trim(),
      toStepId: krokDo.value.trim(),
      kind: rodzaj.value as AutomationDependency['kind'],
    };
    if (warunek.value.trim() !== '') zaleznosc.condition = warunek.value.trim();
    uklad.push(zaleznosc);
    krokZ.value = '';
    krokDo.value = '';
    warunek.value = '';
    wyslij(uklad, `Zależność ${opisZaleznosci(zaleznosc)} dodana.`, zaleznosc);
  });

  zapisz.addEventListener('click', () => wyslij(uklad, 'Układ zależności zapisany.'));
  waliduj.addEventListener('click', () => wyslij(undefined, 'Układ sprawdzony bez zmiany.'));

  /**
   * Eksport mapy zależności w jednej z czterech postaci.
   *
   * Treść składa się w oknie z układu, który rdzeń już oddał, więc eksport nie
   * pyta rdzenia o nic. Nazwa pliku bierze identyfikator automatyki; automatyka
   * niewskazana daje nazwę rodzajową, bo plik ma się zapisać mimo wszystko.
   */
  function eksportuj(rozszerzenie: string, rodzajTresci: string, zloz: () => string): void {
    if (uklad.length === 0) {
      tresc.potwierdzenie('Nie ma czego wyeksportować — układ jest pusty.', false);
      return;
    }
    const podstawa = stan.automatyka() === '' ? 'uklad' : stan.automatyka();
    const nazwaPliku = `${podstawa}-zaleznosci.${rozszerzenie}`;
    pobierzPlik(nazwaPliku, zloz(), rodzajTresci);
    tresc.potwierdzenie(
      `Pobrano plik ${nazwaPliku} — ${uklad.length} zależności z układu odczytanego w oknie. ` +
        'Rdzeń nie brał w tym udziału.',
      true,
    );
  }

  powierzchnia.eksportJson.addEventListener('click', () =>
    eksportuj('json', 'application/json', () => JSON.stringify(uklad, null, 2)));
  powierzchnia.eksportDot.addEventListener('click', () =>
    eksportuj('dot', 'text/vnd.graphviz', () => zapisDot(opisUkladu())));
  powierzchnia.eksportMermaid.addEventListener('click', () =>
    eksportuj('mmd', 'text/plain', () => zapisMermaid(opisUkladu())));
  powierzchnia.eksportRysunku.addEventListener('click', () =>
    eksportuj('svg', 'image/svg+xml', () => {
      // Kanwa rysuje się dopiero przy włączonym przełączniku, więc przed
      // zapisem odświeżamy ją z układu bieżącego — inaczej plik niósłby
      // rysunek sprzed ostatniej zmiany albo pustkę.
      graf.pokaz(opisUkladu());
      return graf.zapisWektorowy();
    }));

  powierzchnia.kanwa.addEventListener('click', () => {
    przestaw(powierzchnia.kanwa);
    przerysuj();
  });
  powierzchnia.stany.addEventListener('click', () => {
    przestaw(powierzchnia.stany);
    przerysuj();
  });
  powierzchnia.powieksz.addEventListener('click', () => graf.powieksz(1));
  powierzchnia.pomniejsz.addEventListener('click', () => graf.powieksz(-1));

  /**
   * Przerysowanie widoków bez pytania rdzenia — po przestawieniu przełącznika.
   * Ocena układu pochodzi z ostatniej odpowiedzi, więc przełącznik nie zmienia
   * werdyktu, wyłącznie to, które widoki stoją na ekranie.
   */
  function przerysuj(): void {
    pokaz(ostatniOdczyt);
  }

  // Układ zmienia się także z pracy innego okna albo innego urządzenia tego
  // konta; zdarzenie mówi, że obraz w oknie jest już nieaktualny.
  const odsubskrybuj = zrodlo.naZmianeUkladu(() => {
    if (stan.automatyka() === '') return;
    wyslij(undefined, 'Układ odczytany ponownie po zmianie zgłoszonej przez rdzeń.');
  });

  return {
    element: rama.element,
    odswiez() {
      if (stan.automatyka() === '') {
        tresc.pusto('Wskaż automatykę w Workflow Builderze, aby ułożyć zależności jej kroków.');
        return;
      }
      wyslij(undefined, 'Układ odczytany z rdzenia.');
    },
    zamknij: odsubskrybuj,
  };
}

/** Zdanie nad wykazem przejść maszyny stanów przebiegu. */
function zdanieOStanach(): HTMLElement {
  const akapit = document.createElement('p');
  akapit.className = 'dn-pole-opis';
  akapit.textContent = zdanieOMaszynieStanow('');
  return akapit;
}

/** Kontrolki okna Orchestratora. */
interface PowierzchniaOrchestratora {
  krokZ: HTMLInputElement;
  krokDo: HTMLInputElement;
  rodzaj: HTMLSelectElement;
  warunek: HTMLInputElement;
  dodaj: HTMLButtonElement;
  zapisz: HTMLButtonElement;
  waliduj: HTMLButtonElement;
  eksportJson: HTMLButtonElement;
  eksportDot: HTMLButtonElement;
  eksportMermaid: HTMLButtonElement;
  eksportRysunku: HTMLButtonElement;
  kanwa: HTMLButtonElement;
  stany: HTMLButtonElement;
  powieksz: HTMLButtonElement;
  pomniejsz: HTMLButtonElement;
}

/**
 * Składa kontrolki, pasek akcji i ciało okna.
 *
 * Czysta konstrukcja: nie domyka się na stanie okna. Przełączniki widoku stoją
 * w narzędziach kontekstowych, bo zmieniają obraz układu, a nie sam układ.
 * Kolejność dokładania jest znacząca — po niej idą sprawdziany widoku.
 */
function zlozPowierzchnie(
  rama: { akcje: HTMLElement; narzedzia: HTMLElement; cialo: HTMLElement },
  stanTresci: HTMLElement,
  pokrycie: PokrycieKomend,
): PowierzchniaOrchestratora {
  const krokZ = pole('Krok poprzedzający', 'identyfikator kroku');
  const krokDo = pole('Krok następujący', 'identyfikator kroku');
  const rodzaj = wybor('Rodzaj zależności', RODZAJE_ZALEZNOSCI);
  const warunek = pole('Warunek przejścia', 'wymagany przy zależności warunkowej');

  const dodaj = przycisk('+ Dodaj zależność', 'dn-btn dn-btn--atrament');
  const zapisz = przycisk('Zapisz układ', 'dn-btn dn-btn--atrament');
  const waliduj = przycisk('Waliduj graf');
  const eksportJson = przycisk('Eksportuj mapę (JSON)');
  const eksportDot = przycisk('Eksportuj mapę (DOT)');
  const eksportMermaid = przycisk('Eksportuj mapę (Mermaid)');
  const eksportRysunku = przycisk('Eksportuj rysunek (SVG)');

  // Usunięcie pojedynczej zależności ma drogę do rdzenia z wykazu: każdy wiersz
  // niesie przycisk „Usuń” wołający komendę `orchestration.dependency.remove`.
  rama.akcje.append(
    dodaj,
    zapisz,
    waliduj,
    eksportJson,
    eksportDot,
    eksportMermaid,
    eksportRysunku,
    pokrycie.przycisk(
      'Bramki dołączenia',
      Command.OrchestrationGateSet,
      'scalenie torów równoległych regułą: wszystkie, dowolny albo licznik',
    ),
    pokrycie.przycisk(
      'Grupy równoległe i sekwencyjne',
      Command.OrchestrationGroupSet,
      'oznaczenie zbioru kroków jako wykonywanych równolegle albo w ścisłej kolejności',
    ),
    pokrycie.przycisk(
      'Transakcje kompensujące',
      Command.OrchestrationCompensationSet,
      'kroki wycofujące skutki przy błędzie w połowie przebiegu',
    ),
    pokrycie.przycisk(
      'Spięcie z kolejkami MultitaskingAI',
      Command.OrchestrationMultitaskingLink,
      'respektowanie zależności układu przez warstwę orkiestracji środowiska MultitaskingAI',
    ),
  );

  const kanwa = przelacznikWidoku('Kanwa grafu', true);
  const stany = przelacznikWidoku('Stany przebiegu', false);
  const powieksz = przycisk('Powiększ graf', 'dn-btn dn-btn--sm dn-btn--zarys');
  const pomniejsz = przycisk('Pomniejsz graf', 'dn-btn dn-btn--sm dn-btn--zarys');
  rama.narzedzia.append(kanwa, stany, powieksz, pomniejsz);

  rama.cialo.append(
    wiersz('Krok poprzedzający', krokZ, { klasa: 'da-wiersz' }),
    wiersz('Krok następujący', krokDo, { klasa: 'da-wiersz' }),
    wiersz('Rodzaj', rodzaj, {
      klasa: 'da-wiersz',
      objasnienie: 'Zależność warunkowa bez warunku zostanie zgłoszona jako zastrzeżenie.',
    }),
    wiersz('Warunek', warunek, { klasa: 'da-wiersz' }),
    stanTresci,
  );

  return {
    krokZ, krokDo, rodzaj, warunek,
    dodaj, zapisz, waliduj,
    eksportJson, eksportDot, eksportMermaid, eksportRysunku,
    kanwa, stany, powieksz, pomniejsz,
  };
}
