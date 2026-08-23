import {
  AutomationExecutionStatus,
  Command,
  QueueAction,
  type AutomationExecution,
} from '../../../../shared/contract';
import { utworzRameOkna } from '../../komponenty/rama-okna';
import {
  pobierzPlik,
  pole,
  poleLiczbowe,
  przelacznikWidoku,
  przestaw,
  przyciskAkcji as przycisk,
  pozycjaWykazu,
  wiersz,
  wybor,
  wykaz,
} from '../../modele/kontrolki-formularza';
import type { PokrycieKomend } from '../pokrycie-komend';
import {
  utworzDziennikPrzebiegow,
  zdanieODzienniku,
} from './dziennik-przebiegow';
import { nazwaStanu } from './maszyna-stanow';
import {
  BEZ_ZAWEZENIA,
  zawezonePrzebiegi,
  type ZawezeniePrzebiegow,
} from './przeglad-przebiegow';
import type { StanAutomatyki } from './stan-automatyki';
import { utworzStanTresci } from './stany-okna';
import {
  osCzasuPrzebiegow,
  panelMetryk,
  szufladaBledow,
  wykazDziennika,
  zapisDziennika,
  zestawieniePrzebiegow,
} from './widok-przebiegow';
import type { ZrodloAutomations } from './zrodlo-automations';

/**
 * Execution Monitor — okno monitora modułu Automations: stan każdego przebiegu,
 * interwencja przy błędzie i ponowne uruchomienie. Wiersz przebiegu niesie etap
 * bieżący, liczbę etapów, stopień ukończenia, stan i licznik obiegów.
 *
 * Przycisk zatrzymania nie ma warunku — nie jest wyszarzany ani przy przebiegu
 * zakończonym, ani przy braku odczytu. Odmowa rdzenia jest widoczna w stanie
 * błędu okna.
 *
 * „Na żywo" znaczy ze zdarzenia, nie z odpytywania: okno słucha
 * `automation.execution.status` i przerysowuje wiersz przebiegu, gdy zdarzenie
 * przyjdzie — także z pracy innego okna albo innego urządzenia tego konta.
 */
export interface OknoExecutionMonitora {
  element: HTMLElement;
  odswiez(): void;
  /** Zamyka nasłuch zdarzenia stanu przebiegu. */
  zamknij(): void;
}

export function utworzOknoExecutionMonitora(
  zrodlo: ZrodloAutomations,
  stan: StanAutomatyki,
  idOkna: string,
  pokrycie: PokrycieKomend,
): OknoExecutionMonitora {
  const rama = utworzRameOkna({
    tytul: 'Execution Monitor',
    rola: 'monitor',
    przeznaczenie:
      'Stan przebiegów automatyki na żywo: etapy, licznik obiegów naprawczych i interwencja przy błędzie.',
    przedrostek: 'da',
  });
  const tresc = utworzStanTresci();
  const przebiegi = new Map<string, AutomationExecution>();
  const dziennik = utworzDziennikPrzebiegow();
  const powierzchnia = zlozPowierzchnieMonitora(rama, tresc.element, pokrycie);

  /** Zawężenie wykazu wzięte z pól paska narzędzi. */
  function zawezenie(): ZawezeniePrzebiegow {
    return {
      ...BEZ_ZAWEZENIA,
      stan: powierzchnia.stanPrzebiegu.value,
      odDnia: powierzchnia.odDnia.value,
      doDnia: powierzchnia.doDnia.value,
      szukane: powierzchnia.szukanie.value,
    };
  }

  /**
   * Rysuje komplet widoków w jednym miejscu treści.
   *
   * Miejsce czyści się przy każdym wskazaniu (`stan-tresci.ts`), więc widoki
   * dokładają się jednym przebiegiem: wykaz przebiegów, oś czasu, metryki,
   * szuflada błędów i dziennik. Przełączniki paska decydują, które z nich stoją
   * na ekranie — wszystkie naraz zajęłyby kolumnę okna pomocniczego bez reszty.
   */
  function rysuj(): void {
    const wszystkie = [...przebiegi.values()];
    const widoczne = zawezonePrzebiegi(wszystkie, zawezenie());
    if (wszystkie.length === 0) {
      tresc.pusto('Ta automatyka nie ma jeszcze ani jednego przebiegu. Uruchom ją w Queue Managerze.');
      return;
    }
    if (widoczne.length === 0) {
      tresc.pusto(
        `Żaden z ${wszystkie.length} przebiegów nie mieści się w zawężeniu. ` +
          'Zawężanie idzie po stronie okna — komenda przebiegów przyjmuje wyłącznie automatykę, ' +
          'pojedynczy przebieg i górną granicę. Wyczyść pola zawężenia, aby zobaczyć wykaz w całości.',
      );
      return;
    }
    const uporzadkowane = [...widoczne].sort((pierwszy, drugi) => drugi.startedAt - pierwszy.startedAt);
    const miejsce = tresc.tresc();
    miejsce.append(listaPrzebiegow(uporzadkowane));
    if (powierzchnia.osCzasu.dataset['wlaczony'] === 'true') {
      miejsce.append(osCzasuPrzebiegow(uporzadkowane));
    }
    if (powierzchnia.metryki.dataset['wlaczony'] === 'true') {
      miejsce.append(panelMetryk(uporzadkowane));
    }
    const bledy = szufladaBledow(uporzadkowane);
    if (bledy !== null) miejsce.append(bledy);
    if (powierzchnia.dziennikWidoczny.dataset['wlaczony'] === 'true') {
      const wpisy = dziennik.wpisy(powierzchnia.stanPrzebiegu.value);
      const wykazWpisow = wykazDziennika(wpisy);
      if (wykazWpisow !== null) miejsce.append(wykazWpisow);
      miejsce.append(akapitOpisowy(zdanieODzienniku(wpisy.length, dziennik.odrzucone())));
    }
  }

  /**
   * Odczyt przebiegów. Oddaje obietnicę, żeby interwencja mogła odświeżyć wykaz
   * i dopiero potem wypowiedzieć swoje zdanie — inaczej odczyt nadpisałby
   * potwierdzenie interwencji zdaniem o obserwacji.
   */
  function odczytaj(): Promise<void> {
    tresc.ladowanie('Odczyt przebiegów automatyki…');
    const zadanie = zadaniePrzebiegow(idOkna, stan.automatyka(), powierzchnia.granica.value);
    return zrodlo.przebiegi(zadanie).then((wynik) => {
      if (!wynik.udany || wynik.wynik === undefined) {
        tresc.blad('Rdzeń nie oddał przebiegów automatyki.', wynik.blad);
        return;
      }
      przebiegi.clear();
      for (const przebieg of wynik.wynik.executions) przebiegi.set(przebieg.id, przebieg);
      rysuj();
      tresc.potwierdzenie(
        wynik.wynik.subscribed
          ? 'Obserwacja przebiegów założona — stan aktualizuje się na żywo.'
          : 'Wykaz odczytany; obserwacji nie założono (rdzeń nie przyjął wskazania okna).',
        wynik.wynik.subscribed,
      );
    });
  }

  /**
   * Interwencja w przebiegu.
   *
   * Zdanie potwierdzenia niesie licznik oddany przez rdzeń: odpowiedź ma pole
   * `cycle`, więc podajemy jego wartość zamiast orzekać o zmianie, której okno
   * nie zmierzyło.
   */
  function interwencja(dzialanie: QueueAction, zdanie: string): void {
    const kolejka = stan.kolejka();
    if (kolejka === '') {
      tresc.blad('Nie wiadomo, którą kolejkę przerwać — załóż albo wskaż kolejkę w Queue Managerze.');
      return;
    }
    const zadanie: Parameters<ZrodloAutomations['dzialanieKolejki']>[0] = {
      queueId: kolejka,
      action: dzialanie,
    };
    if (stan.automatyka() !== '') zadanie.workflowId = stan.automatyka();
    tresc.ladowanie('Interwencja w przebiegu…');
    void zrodlo.dzialanieKolejki(zadanie).then((wynik) => {
      if (!wynik.udany || wynik.wynik === undefined) {
        tresc.blad('Rdzeń odmówił interwencji w przebiegu.', wynik.blad);
        return;
      }
      const oddana = wynik.wynik;
      void odczytaj().then(() => {
        tresc.potwierdzenie(
          `${zdanie} Kolejka ${oddana.id} po interwencji: ${oddana.status}, ` +
            `licznik obiegów naprawczych ${oddana.cycle ?? 0}.`,
          true,
        );
      });
    });
  }

  /**
   * Przerwanie przebiegu wraz z cofnięciem.
   *
   * Opracowanie żąda, by przerwanie było natychmiastowe, a „Cofnij" stało przez
   * krótki czas po akcji. Cofnięciem jest wznowienie tej samej kolejki — innej
   * drogi kontrakt nie ma, bo `QueueAction` nie zna działania odwracającego
   * zatrzymanie. Okno mówi to wprost i nie obiecuje przywrócenia stanu sprzed
   * zatrzymania: wznawia kolejkę, a co z niej zostało, rozstrzyga rdzeń.
   */
  function przerwijZCofnieciem(): void {
    interwencja(QueueAction.Stop, 'Rdzeń przyjął zatrzymanie kolejki przebiegu.');
    powierzchnia.cofnij.hidden = false;
    if (odliczanieCofniecia !== undefined) clearTimeout(odliczanieCofniecia);
    odliczanieCofniecia = setTimeout(() => {
      powierzchnia.cofnij.hidden = true;
    }, CZAS_NA_COFNIECIE_MS);
  }

  function cofnijPrzerwanie(): void {
    powierzchnia.cofnij.hidden = true;
    if (odliczanieCofniecia !== undefined) clearTimeout(odliczanieCofniecia);
    interwencja(
      QueueAction.Resume,
      'Rdzeń przyjął wznowienie kolejki po zatrzymaniu — cofnięcie idzie wznowieniem, ' +
        'bo działania odwracającego zatrzymanie kontrakt nie ma.',
    );
  }

  /** Uchwyt odliczania widoczności cofnięcia; zdejmowany przy zamknięciu okna. */
  let odliczanieCofniecia: ReturnType<typeof setTimeout> | undefined;

  /**
   * Eksport zestawienia przebiegów.
   *
   * Dwie postacie, bo służą dwóm czynnościom: opis w Markdown czyta człowiek,
   * a zestawienie rozdzielane średnikiem wchodzi do arkusza. Postaci
   * przenośnego dokumentu okno nie składa — wymagałaby biblioteki składu, której
   * warstwa kliencka platformy nie ma.
   */
  function raportuj(rozszerzenie: string, rodzajTresci: string, zloz: () => string): void {
    if (przebiegi.size === 0) {
      tresc.potwierdzenie('Nie ma czego wyeksportować — wykaz przebiegów jest pusty.', false);
      return;
    }
    const podstawa = stan.automatyka() === '' ? 'przebiegi' : stan.automatyka();
    const nazwaPliku = `${podstawa}-raport.${rozszerzenie}`;
    pobierzPlik(nazwaPliku, zloz(), rodzajTresci);
    tresc.potwierdzenie(`Pobrano plik ${nazwaPliku}. Rdzeń nie brał w tym udziału.`, true);
  }

  /** Przebiegi po zawężeniu — eksporty biorą to, co Operator widzi na ekranie. */
  function widoczne(): AutomationExecution[] {
    return zawezonePrzebiegi([...przebiegi.values()], zawezenie());
  }

  podepnijAkcjeMonitora(powierzchnia, {
    odczytaj,
    interwencja,
    przerwij: przerwijZCofnieciem,
    cofnij: cofnijPrzerwanie,
    raportMarkdown: () => raportuj('md', 'text/markdown', () => raportPrzebiegow(widoczne())),
    raportZestawienia: () => raportuj('csv', 'text/csv', () => zestawieniePrzebiegow(widoczne())),
    dziennikPliku: () => {
      const wpisy = dziennik.wpisy(powierzchnia.stanPrzebiegu.value);
      if (wpisy.length === 0) {
        tresc.potwierdzenie(
          'Dziennik jest pusty — nie ma czego zapisać. Wiersze logu przychodzą wyłącznie ' +
            'zdarzeniem stanu przebiegu.',
          false,
        );
        return;
      }
      const podstawa = stan.automatyka() === '' ? 'przebiegi' : stan.automatyka();
      pobierzPlik(`${podstawa}-dziennik.txt`, zapisDziennika(wpisy), 'text/plain');
      tresc.potwierdzenie(`Pobrano dziennik: ${wpisy.length} wierszy zebranych przez to okno.`, true);
    },
    przerysuj: rysuj,
  });

  // Aktualizacja na żywo: zdarzenie zastępuje wiersz przebiegu, nie cały wykaz —
  // dzięki temu odczyt Operatora nie znika przy każdej zmianie stanu.
  const odsubskrybuj = zrodlo.naStanPrzebiegu((zdarzenie) => {
    if (stan.automatyka() !== '' && zdarzenie.execution.workflowId !== stan.automatyka()) return;
    przebiegi.set(zdarzenie.execution.id, zdarzenie.execution);
    // Wiersz logu przychodzi wyłącznie tą drogą i drugi raz nie przyjdzie,
    // więc dziennik przyjmuje go przed przerysowaniem widoku.
    dziennik.przyjmij(zdarzenie);
    rysuj();
    if (zdarzenie.stepLabel !== undefined) {
      tresc.potwierdzenie(`Etap: ${zdarzenie.stepLabel}.`, true);
    }
  });

  return {
    element: rama.element,
    odswiez: () => void odczytaj(),
    zamknij() {
      if (odliczanieCofniecia !== undefined) clearTimeout(odliczanieCofniecia);
      odsubskrybuj();
    },
  };
}

/** Jak długo po przerwaniu stoi przycisk cofnięcia. */
const CZAS_NA_COFNIECIE_MS = 15_000;

/** Stany przebiegu w wykazie zawężenia; pusta wartość znaczy „wszystkie”. */
const STANY_ZAWEZENIA: ReadonlyArray<[string, string]> = [
  ['', 'wszystkie stany'],
  ...Object.values(AutomationExecutionStatus).map(
    (stan): [string, string] => [stan, nazwaStanu(stan)],
  ),
];

/** Akapit opisowy pod treścią okna — zdanie o dzienniku i o pochodzeniu danych. */
function akapitOpisowy(zdanie: string): HTMLElement {
  const akapit = document.createElement('p');
  akapit.className = 'dn-pole-opis';
  akapit.textContent = zdanie;
  return akapit;
}

/** Zdanie o przebiegu: etapy, stopień ukończenia, licznik obiegów i powód błędu. */
function opisPrzebiegu(przebieg: AutomationExecution): string {
  const etapy =
    przebieg.totalSteps === undefined || przebieg.totalSteps === 0
      ? 'etapy nieznane'
      : `etap ${przebieg.currentStep ?? 0} z ${przebieg.totalSteps} (${stopien(przebieg)}%)`;
  const obiegi = `obieg ${przebieg.attempt ?? 0}`;
  const powod = przebieg.errorMessage === undefined ? '' : ` — ${przebieg.errorMessage}`;
  return `${przebieg.status}; ${etapy}; ${obiegi}${powod}`;
}

/** Stopień ukończenia w procentach; przebieg zamknięty pokazuje sto. */
function stopien(przebieg: AutomationExecution): number {
  if (przebieg.status === AutomationExecutionStatus.Succeeded) return 100;
  const etapow = przebieg.totalSteps ?? 0;
  if (etapow === 0) return 0;
  return Math.round((((przebieg.currentStep ?? 1) - 1) / etapow) * 100);
}

/** Raport przebiegów w Markdown — treść pliku eksportu. */
function raportPrzebiegow(przebiegi: readonly AutomationExecution[]): string {
  const wiersze = ['# Raport przebiegów automatyki', ''];
  for (const przebieg of przebiegi) {
    wiersze.push(
      `## ${przebieg.id}`,
      `- Automatyka: ${przebieg.workflowId}`,
      `- Stan: ${przebieg.status}`,
      `- ${opisPrzebiegu(przebieg)}`,
      `- Rozpoczęto: ${new Date(przebieg.startedAt).toLocaleString('pl-PL')}`,
      przebieg.finishedAt === undefined
        ? '- Zakończono: przebieg trwa'
        : `- Zakończono: ${new Date(przebieg.finishedAt).toLocaleString('pl-PL')}`,
      '',
    );
  }
  return wiersze.join('\n');
}

/** Przyciski paska akcji Execution Monitora. */
interface AkcjeMonitora {
  odswiezPrzycisk: HTMLButtonElement;
  uruchomPonownie: HTMLButtonElement;
  odNieudanego: HTMLButtonElement;
  przerwij: HTMLButtonElement;
  cofnij: HTMLButtonElement;
  raportMarkdown: HTMLButtonElement;
  raportZestawienia: HTMLButtonElement;
  dziennikPliku: HTMLButtonElement;
}

/**
 * Składa pasek akcji okna i osadza go w ramie.
 *
 * Pozycje bez pokrycia nazywają komendę, której im brakuje, i biorą swoje
 * rozstrzygnięcie z wykazu komend rdzenia — zdanie wpisane tu na sztywno
 * przestałoby być prawdą w dniu, w którym rdzeń tę komendę dostanie.
 */
function zlozAkcjeMonitora(gospodarz: HTMLElement, pokrycie: PokrycieKomend): AkcjeMonitora {
  const odswiezPrzycisk = przycisk('Odśwież przebiegi', 'dn-btn dn-btn--atrament');
  const uruchomPonownie = przycisk('Uruchom ponownie (bieg naprawczy)');
  const odNieudanego = przycisk('Wznów od nieudanego kroku');
  // Przerwanie biegu przerywa pracę, więc bierze wariant `--niebezpieczny`
  // z biblioteki kontrolek.
  const przerwij = przycisk('Przerwij przebieg', 'dn-btn dn-btn--niebezpieczny');
  // Cofnięcie stoi obok przerwania i pokazuje się dopiero po nim: przycisk
  // widoczny stale obiecywałby czynność, która bez przerwania nie ma przedmiotu.
  const cofnij = przycisk('Cofnij przerwanie');
  cofnij.hidden = true;
  const raportMarkdown = przycisk('Eksportuj raport (Markdown)');
  const raportZestawienia = przycisk('Eksportuj zestawienie (CSV)');
  const dziennikPliku = przycisk('Zapisz dziennik');

  gospodarz.append(
    odswiezPrzycisk,
    uruchomPonownie,
    odNieudanego,
    przerwij,
    cofnij,
    raportMarkdown,
    raportZestawienia,
    dziennikPliku,
    pokrycie.przycisk(
      'Log przebiegu z rdzenia',
      Command.AutomationExecutionLog,
      'pełny zapis zdarzeń pojedynczego uruchomienia, także sprzed otwarcia tego okna',
    ),
    pokrycie.przycisk(
      'Stan krokowy przebiegu',
      Command.AutomationExecutionSteps,
      'drążenie przebiegu do stanu każdego kroku osobno',
    ),
    pokrycie.przycisk(
      'Punkty wznowienia',
      Command.AutomationExecutionCheckpointList,
      'wznowienie przebiegu od miejsca przerwania bez powtarzania kroków ukończonych',
    ),
    pokrycie.przycisk(
      'Ładunek kroku',
      Command.AutomationExecutionPayloadGet,
      'wgląd w dane wejściowe i wyjściowe kroku oraz odtworzenie przebiegu z tym ładunkiem',
    ),
    pokrycie.przycisk(
      'Reguły alarmowania',
      Command.AutomationAlertRuleSet,
      'warunki i kanały powiadomień: błąd, przekroczenie czasu, brak uruchomienia',
    ),
    pokrycie.przycisk(
      'Budżety czasu przebiegu',
      Command.AutomationExecutionBudgetSet,
      'maksymalny czas przebiegu i kroku wraz z alarmem przy jego przekroczeniu',
    ),
  );
  return {
    odswiezPrzycisk, uruchomPonownie, odNieudanego, przerwij, cofnij,
    raportMarkdown, raportZestawienia, dziennikPliku,
  };
}

/** Kontrolki okna: pasek akcji, pola zawężenia i przełączniki widoków. */
interface PowierzchniaMonitora extends AkcjeMonitora {
  granica: HTMLInputElement;
  stanPrzebiegu: HTMLSelectElement;
  odDnia: HTMLInputElement;
  doDnia: HTMLInputElement;
  szukanie: HTMLInputElement;
  osCzasu: HTMLButtonElement;
  metryki: HTMLButtonElement;
  dziennikWidoczny: HTMLButtonElement;
}

/** Składa kontrolki okna, pasek akcji i ciało ramy. */
function zlozPowierzchnieMonitora(
  rama: { akcje: HTMLElement; narzedzia: HTMLElement; cialo: HTMLElement },
  stanTresci: HTMLElement,
  pokrycie: PokrycieKomend,
): PowierzchniaMonitora {
  const granica = poleLiczbowe('Górna granica liczby przebiegów', 'domyślnie wszystkie');
  const stanPrzebiegu = wybor('Stan przebiegu', STANY_ZAWEZENIA);
  const odDnia = pole('Od dnia', 'RRRR-MM-DD');
  odDnia.type = 'date';
  const doDnia = pole('Do dnia', 'RRRR-MM-DD');
  doDnia.type = 'date';
  const szukanie = pole('Szukaj w przebiegach', 'fragment identyfikatora albo powodu błędu');
  szukanie.type = 'search';

  const osCzasu = przelacznikWidoku('Oś czasu', true);
  const metryki = przelacznikWidoku('Metryki', true);
  const dziennikWidoczny = przelacznikWidoku('Dziennik', false);

  const akcje = zlozAkcjeMonitora(rama.akcje, pokrycie);
  rama.narzedzia.append(osCzasu, metryki, dziennikWidoczny);
  rama.cialo.append(
    wiersz('Granica wykazu', granica, {
      klasa: 'da-wiersz',
      objasnienie: 'Puste pole zwraca wszystkie przebiegi znane rdzeniowi.',
    }),
    wiersz('Stan', stanPrzebiegu, {
      klasa: 'da-wiersz',
      objasnienie: 'Zawęża wykaz i dziennik po stronie okna — komenda przebiegów nie ma pola stanu.',
    }),
    wiersz('Od dnia', odDnia, { klasa: 'da-wiersz' }),
    wiersz('Do dnia', doDnia, { klasa: 'da-wiersz' }),
    wiersz('Szukanie', szukanie, {
      klasa: 'da-wiersz',
      objasnienie: 'Idzie po identyfikatorze przebiegu, automatyce i powodzie niepowodzenia.',
    }),
    stanTresci,
  );
  return { granica, stanPrzebiegu, odDnia, doDnia, szukanie, osCzasu, metryki, dziennikWidoczny, ...akcje };
}

/**
 * Podpina pasek akcji do czynności okna.
 *
 * Fragment bierze wywołania zwrotne, więc nie musi znać ani źródła, ani stanu
 * modułu: dobór działania kolejki i zdanie potwierdzenia zostają przy przycisku,
 * a praca zostaje w oknie.
 */
function podepnijAkcjeMonitora(
  powierzchnia: PowierzchniaMonitora,
  obsluga: {
    odczytaj: () => void;
    interwencja: (dzialanie: QueueAction, zdanie: string) => void;
    przerwij: () => void;
    cofnij: () => void;
    raportMarkdown: () => void;
    raportZestawienia: () => void;
    dziennikPliku: () => void;
    przerysuj: () => void;
  },
): void {
  powierzchnia.odswiezPrzycisk.addEventListener('click', obsluga.odczytaj);
  powierzchnia.uruchomPonownie.addEventListener('click', () =>
    obsluga.interwencja(QueueAction.Retry, 'Rdzeń przyjął bieg naprawczy.'));
  powierzchnia.odNieudanego.addEventListener('click', () =>
    obsluga.interwencja(QueueAction.Resume, 'Rdzeń przyjął wznowienie kolejki przebiegu.'));
  powierzchnia.przerwij.addEventListener('click', obsluga.przerwij);
  powierzchnia.cofnij.addEventListener('click', obsluga.cofnij);
  powierzchnia.raportMarkdown.addEventListener('click', obsluga.raportMarkdown);
  powierzchnia.raportZestawienia.addEventListener('click', obsluga.raportZestawienia);
  powierzchnia.dziennikPliku.addEventListener('click', obsluga.dziennikPliku);

  // Zawężenie i przełączniki widoku pracują na wykazie już odczytanym, więc
  // przerysowują treść bez pytania rdzenia.
  powierzchnia.stanPrzebiegu.addEventListener('change', obsluga.przerysuj);
  for (const kontrolka of [powierzchnia.odDnia, powierzchnia.doDnia, powierzchnia.szukanie]) {
    kontrolka.addEventListener('input', obsluga.przerysuj);
  }
  for (const przelacznik of [
    powierzchnia.osCzasu,
    powierzchnia.metryki,
    powierzchnia.dziennikWidoczny,
  ]) {
    przelacznik.addEventListener('click', () => {
      przestaw(przelacznik);
      obsluga.przerysuj();
    });
  }
}

/**
 * Zadanie odczytu przebiegów — czysta konstrukcja z wartości pól. Granica pusta
 * albo niepoprawna znaczy „wszystkie", więc nie trafia do żądania wcale.
 */
function zadaniePrzebiegow(
  idOkna: string,
  automatyka: string,
  granica: string,
): Parameters<ZrodloAutomations['przebiegi']>[0] {
  const zadanie: Parameters<ZrodloAutomations['przebiegi']>[0] = { windowId: idOkna };
  if (automatyka !== '') zadanie.workflowId = automatyka;
  const limit = Number.parseInt(granica, 10);
  if (Number.isInteger(limit) && limit > 0) zadanie.limit = limit;
  return zadanie;
}

/**
 * Buduje wykaz przebiegów z danych, bez dostępu do stanu okna. Stan przebiegu
 * ląduje w `data-stan-przebiegu`, żeby sprawdziany widoku i arkusz stylów mogły
 * się o niego zaczepić bez czytania treści.
 */
function listaPrzebiegow(przebiegi: readonly AutomationExecution[]): HTMLElement {
  const lista = wykaz('Przebiegi automatyki', 'da-wykaz');
  for (const przebieg of przebiegi) {
    const pozycja = pozycjaWykazu(przebieg.id, opisPrzebiegu(przebieg), 'da');
    pozycja.element.dataset['stanPrzebiegu'] = przebieg.status;
    lista.append(pozycja.element);
  }
  return lista;
}
