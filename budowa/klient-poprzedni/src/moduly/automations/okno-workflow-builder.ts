import { Command, type AutomationWorkflow } from '../../../../shared/contract';
import { utworzEdytorKrokow, type EdytorKrokow } from './edytor-krokow';
import { KODY_OKIEN } from './kody-okien';
import { utworzRameOkna } from '../../komponenty/rama-okna';
import {
  pobierzPlik,
  pole,
  przelacznik,
  przyciskAkcji as przycisk,
  pozycjaWykazu,
  wiersz,
  wykaz,
} from '../../modele/kontrolki-formularza';
import type { PokrycieKomend } from '../pokrycie-komend';
import type { StanAutomatyki } from './stan-automatyki';
import { utworzStanTresci, type StanTresci } from './stany-okna';
import {
  zastrzezeniaDefinicji,
  zdanieWalidacji,
  type ZastrzezenieDefinicji,
} from './walidacja-definicji';
import { wczytajPlik } from './wymiana-definicji';
import type { ZrodloAutomations } from './zrodlo-automations';

/**
 * Workflow Builder — okno kreatora modułu Automations. Kroki, warunek
 * i kolejność prowadzi edytor (`edytor-krokow.ts`); to okno odpowiada za
 * tożsamość automatyki, zapis do rdzenia i panel akcji.
 *
 * Cofnij i ponów działają miejscowo, na pracy sprzed zapisu. Po zapisie śladem
 * jest numer wersji definicji, a powrót do wersji wcześniejszej należy do
 * pozycji „Historia wersji" — ta nazywa swoją komendę i mówi, czy rdzeń ma dla
 * niej uchwyt.
 *
 * Zdanie potwierdzenia powstaje z odpowiedzi rdzenia — z identyfikatora, numeru
 * wersji i liczby kroków. Przy duplikowaniu okno porównuje identyfikator sprzed
 * czynności z tym, który wrócił, bo dopiero różnica dowodzi, że powstała nowa
 * definicja.
 */
export interface OknoWorkflowBuilder {
  element: HTMLElement;
  /** Odczytuje wykaz automatyk z rdzenia i wypełnia edytor. */
  odswiez(): void;
}

/** Skok do innego okna modułu — nawigacja układu. */
export interface SkokiOkien {
  pokaz(kodOkna: string): void;
}

export function utworzOknoWorkflowBuilder(
  zrodlo: ZrodloAutomations,
  stan: StanAutomatyki,
  skoki: SkokiOkien,
  pokrycie: PokrycieKomend,
): OknoWorkflowBuilder {
  const rama = utworzRameOkna({
    tytul: 'Workflow Builder',
    rola: 'kreator',
    przeznaczenie:
      'Budowa definicji automatyki: kroki, warunki i kolejność wykonania. Zapis tworzy komponent własny, który wykonuje silnik kolejek.',
    przedrostek: 'da',
  });
  const tresc = utworzStanTresci();
  const edytor = utworzEdytorKrokow(() => {
    stos.odnotuj();
    sprawdzDefinicje();
  });
  const stos = utworzStosZmianKrokow(edytor, tresc);
  const powierzchnia = zlozPowierzchnieBudowniczego(rama, edytor.element, tresc.element, pokrycie);
  const { wskazanie, nazwa, opis, czynna } = powierzchnia;

  /**
   * Walidacja definicji po stronie okna.
   *
   * Idzie przy każdej zmianie kroków, nie dopiero przy zapisie: zastrzeżenie
   * postawione w chwili wpisywania jest poprawką, a to samo zastrzeżenie
   * postawione po zapisie jest już tylko wiadomością o definicji wadliwej,
   * która zdążyła trafić do magazynu automatyk.
   *
   * Sygnalizacja jest podwójna, bo służy dwóm czynnościom: znacznik przy kroku
   * pokazuje, który krok poprawić, a wykaz w treści okna mówi, co dokładnie.
   * Oddaje zastrzeżenia, żeby zapis mógł dopowiedzieć o nich w potwierdzeniu.
   */
  function sprawdzDefinicje(): ZastrzezenieDefinicji[] {
    const zastrzezenia = zastrzezeniaDefinicji(edytor.kroki());
    edytor.oznaczZastrzezenia(zastrzezenia);
    powierzchnia.ocenaDefinicji.textContent = zdanieWalidacji(zastrzezenia);
    powierzchnia.ocenaDefinicji.dataset['ocena'] =
      zastrzezenia.length === 0 ? 'bez-zastrzezen' : 'z-zastrzezeniami';
    return zastrzezenia;
  }

  function pokaz(definicja: AutomationWorkflow): void {
    wskazanie.value = definicja.id;
    nazwa.value = definicja.name;
    opis.value = definicja.description ?? '';
    czynna.checked = definicja.enabled;
    edytor.wczytaj(definicja.steps ?? []);
    stos.przyjmijStanBiezacy();
    const zastrzezenia = sprawdzDefinicje();
    const miejsce = tresc.tresc();
    miejsce.append(podsumowanieAutomatyki(definicja));
    const wykazZastrzezen = listaZastrzezen(zastrzezenia);
    if (wykazZastrzezen !== null) miejsce.append(wykazZastrzezen);
  }

  function zadanieZapisu(zId: boolean) {
    const zadanie: Parameters<ZrodloAutomations['zapiszAutomatyke']>[0] = {
      name: nazwa.value.trim() === '' ? 'Automatyka bez nazwy' : nazwa.value.trim(),
      steps: edytor.kroki(),
      enabled: czynna.checked,
    };
    if (opis.value.trim() !== '') zadanie.description = opis.value.trim();
    if (zId && wskazanie.value.trim() !== '') zadanie.workflowId = wskazanie.value.trim();
    return zadanie;
  }

  /**
   * Zapis definicji. `zId` mówi, czy w żądaniu idzie identyfikator zastany —
   * bez niego rdzeń zakłada definicję nową, i to jest duplikowanie.
   */
  function zapiszDefinicje(zId: boolean): void {
    const przedZapisem = wskazanie.value.trim();
    // Zastrzeżenia liczymy przed wysłaniem, bo po zapisie okno pokazuje już
    // definicję oddaną przez rdzeń. Zapis idzie mimo nich — walidacja nie jest
    // bramą, tak samo jak walidacja układu po stronie rdzenia.
    const zastrzezenia = sprawdzDefinicje();
    tresc.ladowanie('Zapis definicji automatyki…');
    void zrodlo.zapiszAutomatyke(zadanieZapisu(zId)).then((wynik) => {
      if (!wynik.udany || wynik.wynik === undefined) {
        tresc.blad('Rdzeń nie zapisał definicji automatyki.', wynik.blad);
        return;
      }
      pokaz(wynik.wynik);
      stan.ustawDefinicje(wynik.wynik);
      const dopowiedzenie =
        zastrzezenia.length === 0
          ? ''
          : ` Definicja miała przy zapisie ${zastrzezenia.length} zastrzeżeń okna — zapis ich nie wstrzymał.`;
      tresc.potwierdzenie(`${zdanieZapisu(wynik.wynik, przedZapisem, zId)}${dopowiedzenie}`, true);
    });
  }

  function odczytaj(): void {
    tresc.ladowanie('Odczyt wykazu automatyk…');
    void zrodlo.automatyki({}).then((wynik) => {
      if (!wynik.udany || wynik.wynik === undefined) {
        tresc.blad('Rdzeń nie oddał wykazu automatyk.', wynik.blad);
        return;
      }
      const wybrana =
        wynik.wynik.find((automatyka) => automatyka.id === stan.automatyka()) ?? wynik.wynik[0];
      if (wybrana === undefined) {
        tresc.pusto('Nie ma jeszcze ani jednej automatyki. Dodaj krok i zapisz — powstanie nowa.');
        return;
      }
      pokaz(wybrana);
      stan.ustawDefinicje(wybrana);
    });
  }

  function eksportuj(): void {
    const definicja = { ...zadanieZapisu(true), id: wskazanie.value.trim() };
    const nazwaPliku = `${definicja.id === '' ? 'automatyka' : definicja.id}.json`;
    pobierzPlik(nazwaPliku, JSON.stringify(definicja, null, 2));
    tresc.potwierdzenie(
      `Pobrano plik ${nazwaPliku} — treść z pól okna, kroków: ${(definicja.steps ?? []).length}. ` +
        'Rdzeń nie brał w tym udziału.',
      true,
    );
  }

  /**
   * Walidacja na żądanie: to samo sprawdzenie, które idzie przy każdej zmianie,
   * lecz wypisane wykazem w treści okna. Zastrzeżenia dotyczące kroków rozstrzyga
   * okno; cykle i ścieżkę krytyczną całego układu rozstrzyga rdzeń w Orchestratorze,
   * więc zdanie mówi wprost, gdzie szukać drugiej połowy oceny.
   */
  function walidujDefinicje(): void {
    const zastrzezenia = sprawdzDefinicje();
    const miejsce = tresc.tresc();
    const wykazZastrzezen = listaZastrzezen(zastrzezenia);
    if (wykazZastrzezen !== null) miejsce.append(wykazZastrzezen);
    tresc.potwierdzenie(
      `${zdanieWalidacji(zastrzezenia)} Ocenę układu zależności — cykle i ścieżkę krytyczną — ` +
        'wydaje rdzeń w oknie Orchestratora.',
      zastrzezenia.length === 0,
    );
  }

  podepnijAkcjeBudowniczego(powierzchnia, {
    dodaj: () => edytor.dodajKrok(),
    zapisz: () => zapiszDefinicje(true),
    duplikuj: () => zapiszDefinicje(false),
    cofnij: () => stos.cofnij(),
    ponow: () => stos.ponow(),
    waliduj: walidujDefinicje,
    uklad: () => skoki.pokaz(KODY_OKIEN.orchestrator),
    wykazAutomatyk: odczytaj,
    eksport: eksportuj,
    importuj: () => wczytajPlik(edytor, tresc),
  });

  stan.naZmiane(() => {
    const definicja = stan.definicja();
    if (definicja !== null && definicja.id !== wskazanie.value) pokaz(definicja);
  });

  return { element: rama.element, odswiez: odczytaj };
}

/**
 * Wykaz zastrzeżeń walidacji; `null` znaczy „definicja bez zastrzeżeń", nie
 * „nie sprawdzono". Miejsce kroku idzie w tytule pozycji, bo po nim Operator
 * odnajduje wiersz w edytorze.
 */
function listaZastrzezen(zastrzezenia: readonly ZastrzezenieDefinicji[]): HTMLElement | null {
  if (zastrzezenia.length === 0) return null;
  const lista = wykaz('Zastrzeżenia okna do definicji', 'da-wykaz');
  for (const zastrzezenie of zastrzezenia) {
    const nazwaKroku = zastrzezenie.idKroku === '' ? 'krok bez identyfikatora' : zastrzezenie.idKroku;
    const pozycja = pozycjaWykazu(
      `Krok ${zastrzezenie.miejsce} — ${nazwaKroku}`,
      zastrzezenie.zdanie,
      'da',
    );
    pozycja.element.dataset['zastrzezenie'] = zastrzezenie.waga;
    lista.append(pozycja.element);
  }
  return lista;
}

/** Kontrolki okna Workflow Buildera: tożsamość automatyki i pasek akcji. */
interface PowierzchniaBudowniczego {
  wskazanie: HTMLInputElement;
  nazwa: HTMLInputElement;
  opis: HTMLInputElement;
  czynna: HTMLInputElement;
  dodaj: HTMLButtonElement;
  zapisz: HTMLButtonElement;
  waliduj: HTMLButtonElement;
  uklad: HTMLButtonElement;
  duplikuj: HTMLButtonElement;
  cofnij: HTMLButtonElement;
  ponow: HTMLButtonElement;
  eksport: HTMLButtonElement;
  importuj: HTMLButtonElement;
  wykazAutomatyk: HTMLButtonElement;
  /** Zdanie oceny definicji pod edytorem — odpowiedź walidacji na żywo. */
  ocenaDefinicji: HTMLElement;
}

/**
 * Składa kontrolki tożsamości, pasek akcji i ciało okna. Fragment nie domyka się
 * ani na stosie zmian, ani na źródle — dostaje gotowy element edytora i miejsce
 * stanu treści. Kolejność dokładania przycisków jest znacząca, bo po niej idą
 * sprawdziany widoku.
 */
function zlozPowierzchnieBudowniczego(
  rama: { akcje: HTMLElement; cialo: HTMLElement },
  edytor: HTMLElement,
  stanTresci: HTMLElement,
  pokrycie: PokrycieKomend,
): PowierzchniaBudowniczego {
  const wskazanie = pole('Identyfikator automatyki', 'puste zakłada nową automatykę');
  const nazwa = pole('Nazwa automatyki', 'np. Kompletowanie akt sprawy');
  const opis = pole('Opis automatyki', 'do czego służy ta automatyka');
  const czynna = przelacznik('Automatyka czynna');
  czynna.checked = true;

  const dodaj = przycisk('+ Dodaj krok', 'dn-btn dn-btn--atrament');
  const zapisz = przycisk('Zapisz automatykę', 'dn-btn dn-btn--atrament');
  const waliduj = przycisk('Waliduj definicję');
  const uklad = przycisk('Układ zależności');
  const duplikuj = przycisk('Duplikuj');
  const cofnij = przycisk('Cofnij');
  const ponow = przycisk('Ponów');
  const eksport = przycisk('Eksportuj definicję (JSON)');
  const importuj = przycisk('Importuj definicję (JSON)');
  const wykazAutomatyk = przycisk('Odśwież wykaz automatyk');

  rama.akcje.append(
    dodaj,
    zapisz,
    waliduj,
    uklad,
    duplikuj,
    cofnij,
    ponow,
    eksport,
    importuj,
    wykazAutomatyk,
    pokrycie.przycisk(
      'Test i symulacja',
      Command.AutomationWorkflowSimulate,
      'przebieg próbny bez wpięcia produkcyjnego, z podglądem wyniku każdego kroku',
    ),
    pokrycie.przycisk(
      'Historia wersji',
      Command.AutomationWorkflowVersionList,
      'wersje wcześniejsze definicji wraz z ich porównaniem i przywróceniem',
    ),
    pokrycie.przycisk(
      'Szablon z definicji',
      Command.AutomationTemplateSave,
      'zapis przepływu jako szablonu z parametrami do wielokrotnego użycia',
    ),
    pokrycie.przycisk(
      'Tagowanie automatyki',
      Command.AutomationWorkflowTagSet,
      'kategoryzacja automatyk i wyszukiwanie po tagu',
    ),
    pokrycie.przycisk(
      'Notatka przy kroku',
      Command.AutomationStepNoteSet,
      'dokumentacja opisowa dołączana do pojedynczego kroku',
    ),
    pokrycie.przycisk(
      'Zmienne przepływu i mapowanie danych',
      Command.AutomationWorkflowVariablesSet,
      'zmienne procesu przekazywane między krokami wraz z mapowaniem wyjść na wejścia',
    ),
  );

  const tozsamosc = document.createElement('div');
  tozsamosc.className = 'da-tozsamosc';
  tozsamosc.append(
    wiersz('Automatyka', wskazanie, {
      klasa: 'da-wiersz',
      // `automation.workflow.save` z kodem, którego rdzeń nie zna, nie jest
      // odmową — zakłada definicję pod tym kodem w wersji 1.
      objasnienie:
        'Puste pole zakłada nową automatykę. Kod zastany zmienia jego definicję; ' +
        'kod nieznany rdzeniowi również zakłada nową — pod tym właśnie kodem.',
    }),
    wiersz('Nazwa', nazwa, { klasa: 'da-wiersz' }),
    wiersz('Opis', opis, { klasa: 'da-wiersz' }),
    wiersz('Czynna', czynna, {
      klasa: 'da-wiersz',
      objasnienie: 'Automatyka wyłączona zostaje w wykazie, lecz nie pracuje.',
    }),
  );
  // Ocena definicji stoi pod edytorem, a nie w pasie akcji: mówi o treści
  // kroków, więc ma być tam, gdzie te kroki się wpisuje.
  const ocenaDefinicji = document.createElement('p');
  ocenaDefinicji.className = 'dn-pole-opis da-ocena-definicji';
  ocenaDefinicji.setAttribute('aria-live', 'polite');

  rama.cialo.append(tozsamosc, edytor, ocenaDefinicji, stanTresci);

  return {
    wskazanie, nazwa, opis, czynna,
    dodaj, zapisz, waliduj, uklad, duplikuj, cofnij, ponow, eksport, importuj, wykazAutomatyk,
    ocenaDefinicji,
  };
}

/**
 * Podpina pasek akcji do czynności okna. Fragment bierze wywołania zwrotne, więc
 * nie zna ani źródła, ani stosu zmian — wie tylko, który przycisk co uruchamia.
 */
function podepnijAkcjeBudowniczego(
  powierzchnia: PowierzchniaBudowniczego,
  obsluga: Record<
    'dodaj' | 'zapisz' | 'duplikuj' | 'cofnij' | 'ponow' | 'waliduj' | 'uklad' | 'wykazAutomatyk' | 'eksport' | 'importuj',
    () => void
  >,
): void {
  powierzchnia.dodaj.addEventListener('click', obsluga.dodaj);
  powierzchnia.zapisz.addEventListener('click', obsluga.zapisz);
  powierzchnia.duplikuj.addEventListener('click', obsluga.duplikuj);
  powierzchnia.cofnij.addEventListener('click', obsluga.cofnij);
  powierzchnia.ponow.addEventListener('click', obsluga.ponow);
  powierzchnia.waliduj.addEventListener('click', obsluga.waliduj);
  powierzchnia.uklad.addEventListener('click', obsluga.uklad);
  powierzchnia.wykazAutomatyk.addEventListener('click', obsluga.wykazAutomatyk);
  powierzchnia.eksport.addEventListener('click', obsluga.eksport);
  powierzchnia.importuj.addEventListener('click', obsluga.importuj);
}

/**
 * Zdanie potwierdzenia zapisu — składane z odpowiedzi rdzenia, nie z żądania.
 *
 * Przy duplikowaniu okno wysyła żądanie bez identyfikatora i o nowej definicji
 * mówi dopiero wtedy, gdy identyfikator w odpowiedzi różni się od tego, który
 * stał w oknie przed czynnością.
 */
function zdanieZapisu(
  definicja: AutomationWorkflow,
  przedZapisem: string,
  zId: boolean,
): string {
  const opis =
    `${definicja.id} — wersja ${definicja.version ?? 1}, kroków: ${(definicja.steps ?? []).length}`;
  if (zId) return `Rdzeń zapisał automatykę ${opis}.`;
  if (przedZapisem !== '' && definicja.id === przedZapisem) {
    return (
      `Rdzeń oddał tę samą automatykę ${opis} — nowa definicja NIE powstała, ` +
      'duplikowanie się nie odbyło.'
    );
  }
  return `Rdzeń założył nową automatykę ${opis}.`;
}

/** Podsumowanie definicji — czysta konstrukcja z bytu `AutomationWorkflow`. */
function podsumowanieAutomatyki(definicja: AutomationWorkflow): HTMLElement {
  const podsumowanie = document.createElement('p');
  podsumowanie.className = 'dn-pole-opis';
  podsumowanie.textContent =
    `Automatyka ${definicja.id} — wersja ${definicja.version ?? 1}, kroków: ${(definicja.steps ?? []).length}.`;
  return podsumowanie;
}

/** Stos zmian kroków — cofnij i ponów miejscowe, sprzed zapisu do rdzenia. */
interface StosZmianKrokow {
  /** Odnotowuje zmianę kroków, jeżeli różni się od stanu ostatnio zapamiętanego. */
  odnotuj(): void;
  /** Przyjmuje stan bieżący za ostatni — po wczytaniu definicji z rdzenia. */
  przyjmijStanBiezacy(): void;
  cofnij(): void;
  ponow(): void;
}

/**
 * Stos zmian kroków okna. Domyka się wyłącznie na edytorze kroków i stanie
 * treści — nie zna ani źródła, ani stanu modułu, ani pól tożsamości. Migawka
 * jest zapisem JSON kroków, bo porównanie stanów ma być całościowe.
 */
function utworzStosZmianKrokow(edytor: EdytorKrokow, tresc: StanTresci): StosZmianKrokow {
  const historia: string[] = [];
  const cofniete: string[] = [];

  function migawka(): string {
    return JSON.stringify(edytor.kroki());
  }

  let ostatnia = migawka();

  function przywroc(zrodloStosu: string[], celStosu: string[]): void {
    const poprzednia = zrodloStosu.pop();
    if (poprzednia === undefined) {
      tresc.potwierdzenie('Nie ma czego przywrócić — stos zmian jest pusty.', false);
      return;
    }
    celStosu.push(ostatnia);
    edytor.wczytaj(JSON.parse(poprzednia));
    ostatnia = migawka();
    tresc.potwierdzenie('Stan kroków przywrócony ze stosu zmian okna.', true);
  }

  return {
    odnotuj() {
      const biezaca = migawka();
      if (biezaca === ostatnia) return;
      historia.push(ostatnia);
      cofniete.length = 0;
      ostatnia = biezaca;
    },

    przyjmijStanBiezacy() {
      ostatnia = migawka();
    },

    cofnij: () => przywroc(historia, cofniete),
    ponow: () => przywroc(cofniete, historia),
  };
}
