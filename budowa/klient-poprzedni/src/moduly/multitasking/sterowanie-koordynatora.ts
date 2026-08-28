import { QueueAction, type ErrorInfo } from '../../../../shared/contract';
import { przyciskAkcji as przycisk } from '../../modele/kontrolki-formularza';
import { rozdziel, type Zlecenie } from './tryby-wspolpracy';
import type { StanMultitaskingu } from './stan-multitaskingu';
import type { ZrodloBiegu } from './zrodlo-biegu';

// Siedem przycisków sterowania koordynatora: Start, Stop, Pauza, Wznów, Przekaż, Powtórz, Waliduj.

/** Interfejs zestawia element sterowania koordynatora wraz z funkcją odświeżania, która przestawia dostępność przycisków wedle obsady i kolejki. */
export interface SterowanieKoordynatora {
  element: HTMLElement;
  /** Przestawia dostępność przycisków wedle obsady i kolejki. */
  odswiez(): void;
}

/** Interfejs zestawia zależności sterowania koordynatora: źródło biegu, stan multitaskingu oraz funkcje odczytu treści zlecenia, walidacji i potwierdzenia. */
export interface OpcjeSterowania {
  zrodlo: ZrodloBiegu;
  stan: StanMultitaskingu;
  /** Treść zlecenia wpisana przez Operatora w kreatorze promptu. */
  trescZlecenia(): string;
  /** Odczyt stanu okien na żądanie — nośnik walidacji przebiegu. */
  zwaliduj(): void;
  /** Krótkie potwierdzenie czynności widoczne w oknie. */
  potwierdz(zdanie: string, udane: boolean): void;
}

export function utworzSterowanie(opcje: OpcjeSterowania): SterowanieKoordynatora {
  const { zrodlo, stan } = opcje;

  const { element, naKolejce, przekaz, waliduj } = zlozPasSterowaniaKoordynatora({
    kolejka: (dzialanie) => {
      void steruj(dzialanie);
    },
    stop: () => {
      void zatrzymajBiegKoordynatora(zrodlo, stan, opcje.potwierdz);
    },
    przekazanie: () => {
      void przekazZlecenie();
    },
    walidacja: () => opcje.zwaliduj(),
  });

  // Sterowanie kolejką jednym wywołaniem: zdanie zestawia stan sprzed wywołania ze stanem oddanym.
  async function steruj(dzialanie: QueueAction): Promise<void> {
    const kolejka = stan.kolejkaBiezaca();
    if (kolejka === null) {
      opcje.potwierdz('Nie ma etapu do sterowania — załóż etap w planie.', false);
      return;
    }
    const stanPrzed = kolejka.status;
    const wynik = await zrodlo.sterujKolejka({ queueId: kolejka.id, action: dzialanie });
    if (!wynik.udany || wynik.wynik === undefined) {
      opcje.potwierdz(`Rdzeń odmówił działania ${dzialanie}. ${powod(wynik.blad)}`, false);
      return;
    }
    const po = wynik.wynik;
    stan.zapiszKolejke(po);
    opcje.potwierdz(
      po.status === stanPrzed
        ? `Etap „${po.name ?? po.id}" pozostał w stanie ${po.status} — rdzeń przyjął żądanie, ale stan kolejki się nie zmienił.`
        : `Etap „${po.name ?? po.id}": rdzeń przestawił kolejkę ze stanu ${stanPrzed} na ${po.status}.`,
      true,
    );
  }

  // Doręczenie treści i utrwalenie więzi to dwa osobne niepowodzenia; wracają osobnymi wykazami.
  async function przekazZlecenie(): Promise<void> {
    const tresc = opcje.trescZlecenia().trim();
    if (tresc === '') {
      opcje.potwierdz('Zlecenie puste — koordynator nie przekazuje pustej tury.', false);
      return;
    }
    const zlecenia = rozdziel(stan.tryb(), tresc, {
      wykonawcy: stan.obsada().wykonawcy.map((okno) => okno.id),
      poprzedni: stan.poprzedniWykonawca(),
      wynik: (okno) => stan.strumien().wynik(okno),
    });
    if (zlecenia.length === 0) {
      opcje.potwierdz('Obsada nie ma wykonawcy — załóż okno wykonawcy.', false);
      return;
    }
    const koordynator = stan.obsada().koordynator;
    const przebieg = await wyslijZleceniaWykonawcom(zrodlo, zlecenia, koordynator?.id ?? '');
    const ostatnieDoreczone = przebieg.doreczone[przebieg.doreczone.length - 1];
    if (ostatnieDoreczone !== undefined) stan.zapamietajPrzekazanie(ostatnieDoreczone);
    opcje.potwierdz(zdanieOPrzekazaniu(przebieg, zlecenia.length), przebieg.niedoreczone.length === 0);
  }

  function odswiez(): void {
    // Przyciski zostają czynne bez przesłanki: każda czynność ma własną odmowę z powodem.
    const brakEtapu = stan.kolejkaBiezaca() === null;
    for (const kontrolka of naKolejce) {
      oznaczCzynnosc(kontrolka, brakEtapu ? 'plan nie ma jeszcze etapu' : '');
    }
    oznaczCzynnosc(
      przekaz,
      stan.obsada().wykonawcy.length === 0 ? 'obsada nie ma wykonawcy' : '',
    );
    oznaczCzynnosc(waliduj, stan.obsada().koordynator === null ? 'obsada nie ma koordynatora' : '');
  }

  odswiez();
  return { element, odswiez };
}

/** Interfejs zestawia pas sterowania koordynatora wraz z przyciskami, które wytwórnia gasi albo znakuje wedle stanu obsady i kolejki. */
interface PasSterowaniaKoordynatora {
  element: HTMLElement;
  /** Cztery przyciski jednego silnika kolejek — gasną, gdy nie ma etapu. */
  naKolejce: readonly HTMLButtonElement[];
  przekaz: HTMLButtonElement;
  waliduj: HTMLButtonElement;
}

/** Interfejs nazywa cztery czynności pasa sterowania: pas ich nie wykonuje, tylko o nie prosi wywołanie zwrotne, które je zna. */
interface CzynnosciSterowania {
  kolejka(dzialanie: QueueAction): void;
  stop(): void;
  przekazanie(): void;
  walidacja(): void;
}

/** Funkcja składa siedem kontrolek pasa sterowania. Pas jest czystą konstrukcją: nie zna źródła biegu ani stanu wspólnego, każde kliknięcie oddaje wywołaniu zwrotnemu. */
function zlozPasSterowaniaKoordynatora(
  czynnosci: CzynnosciSterowania,
): PasSterowaniaKoordynatora {
  const element = document.createElement('div');
  element.className = 'dm-sterowanie';
  element.setAttribute('aria-label', 'Sterowanie koordynatora');

  const opisy = [
    ['Start', 'dn-btn dn-btn--atrament', QueueAction.Start],
    ['Pauza', 'dn-btn', QueueAction.Pause],
    ['Wznów', 'dn-btn', QueueAction.Resume],
    ['Powtórz', 'dn-btn', QueueAction.Retry],
  ] as const;

  const naKolejce = opisy.map(([etykieta, klasa, dzialanie]) => {
    const kontrolka = przycisk(etykieta, klasa);
    kontrolka.dataset['dzialanie'] = dzialanie;
    kontrolka.addEventListener('click', () => {
      czynnosci.kolejka(dzialanie);
    });
    element.append(kontrolka);
    return kontrolka;
  });

  // Zatrzymanie jest czynne bez etapu i przed pierwszym obiegiem, gasi turę bez powstałej kolejki.
  const stop = przycisk('Stop', 'dn-btn dn-btn--niebezpieczny');
  stop.dataset['dzialanie'] = QueueAction.Stop;
  stop.addEventListener('click', () => {
    czynnosci.stop();
  });
  element.append(stop);

  const przekaz = przycisk('Przekaż', 'dn-btn dn-btn--sygnal');
  przekaz.dataset['dzialanie'] = 'przekazanie';
  przekaz.addEventListener('click', () => {
    czynnosci.przekazanie();
  });

  const waliduj = przycisk('Waliduj');
  waliduj.dataset['dzialanie'] = 'walidacja';
  waliduj.addEventListener('click', () => {
    czynnosci.walidacja();
  });

  element.append(przekaz, waliduj);

  return { element, naKolejce, przekaz, waliduj };
}

/** Funkcja zatrzymuje bieg koordynatora: naraz kolejkę etapu i tury wykonawców, najpierw turę, która zużywa czas modelu, a dopiero potem kolejkę mogącą ją wznowić. */
async function zatrzymajBiegKoordynatora(
  zrodlo: ZrodloBiegu,
  stan: StanMultitaskingu,
  potwierdz: (zdanie: string, udane: boolean) => void,
): Promise<void> {
  const zgaszone: string[] = [];
  const odmowy: string[] = [];
  for (const okno of stan.obsada().wykonawcy) {
    const wynik = await zrodlo.zatrzymaj({ windowId: okno.id });
    // Odmowa gaszenia idzie do osobnego wykazu, by zdanie o niczym niebiegnącym nie zakryło tury.
    if (!wynik.udany) odmowy.push(`okno ${okno.id}: ${powod(wynik.blad)}`);
    else if (wynik.wynik?.stopped === true) zgaszone.push(`tura okna ${okno.id}`);
  }
  const kolejka = stan.kolejkaBiezaca();
  if (kolejka !== null) {
    const wynik = await zrodlo.sterujKolejka({ queueId: kolejka.id, action: QueueAction.Stop });
    if (wynik.udany && wynik.wynik !== undefined) {
      stan.zapiszKolejke(wynik.wynik);
      zgaszone.push(`etap „${kolejka.name ?? kolejka.id}"`);
    } else if (!wynik.udany) {
      potwierdz(`Rdzeń odmówił zatrzymania etapu. ${powod(wynik.blad)}`, false);
      return;
    }
  }
  const zdanie =
    zgaszone.length === 0
      ? 'Rdzeń nie zgłosił zatrzymania niczego — nic nie biegło.'
      : `Rdzeń zatrzymał: ${zgaszone.join(', ')}.`;
  potwierdz(
    odmowy.length === 0 ? zdanie : `${zdanie} Odmowy zatrzymania (${odmowy.length}): ${odmowy.join(' ')}`,
    odmowy.length === 0,
  );
}

/** Interfejs rozdziela przebieg przekazania zlecenia na trzy wykazy okien: doręczone, niedoręczone oraz doręczone bez utrwalenia więzi. */
interface PrzebiegPrzekazania {
  /** Okna, do których treść dotarła — wedle odpowiedzi `message.send`. */
  doreczone: string[];
  /** Okna, do których treść nie dotarła, wraz z powodem odmowy. */
  niedoreczone: string[];
  /** Okna doręczone, których więzi rdzeń nie zapisał, wraz z powodem. */
  nieutrwalone: string[];
}

/** Funkcja przekazuje zlecenie wykonawcom: doręcza treść oknu wykonawcy i utrwala w rdzeniu więź zlecenia, dwiema osobnymi czynnościami liczonymi osobno. */
async function wyslijZleceniaWykonawcom(
  zrodlo: ZrodloBiegu,
  zlecenia: readonly Zlecenie[],
  idKoordynatora: string,
): Promise<PrzebiegPrzekazania> {
  const przebieg: PrzebiegPrzekazania = { doreczone: [], niedoreczone: [], nieutrwalone: [] };
  for (const zlecenie of zlecenia) {
    if (idKoordynatora !== '') {
      const zapis = await zrodlo.przekaz({
        fromWindowId: idKoordynatora,
        toWindowId: zlecenie.okno,
        instruction: zlecenie.tresc,
      });
      if (!zapis.udany) przebieg.nieutrwalone.push(`${zlecenie.okno}: ${powod(zapis.blad)}`);
    }
    const wynik = await zrodlo.wyslij({ windowId: zlecenie.okno, content: zlecenie.tresc });
    // Dowodem doręczenia jest wiadomość, która wróciła, nie sama zgoda bez treści odpowiedzi.
    if (wynik.udany && wynik.wynik !== undefined) przebieg.doreczone.push(zlecenie.okno);
    else przebieg.niedoreczone.push(`${zlecenie.okno}: ${powod(wynik.blad)}`);
  }
  return przebieg;
}

/**
 * Zdanie o przekazaniu — liczy okna i rozdziela doręczenie od utrwalenia.
 *
 * Czysta zamiana przebiegu na zdanie: nie zna ani obsady, ani rdzenia.
 */
function zdanieOPrzekazaniu(przebieg: PrzebiegPrzekazania, zleconych: number): string {
  const czesci: string[] = [];
  czesci.push(
    przebieg.doreczone.length === 0
      ? `Zlecenie nie dotarło do żadnego z ${zleconych} okien wykonawcy.`
      : `Rdzeń przyjął zlecenie dla ${przebieg.doreczone.length} z ${zleconych} okien wykonawcy: ${przebieg.doreczone.join(', ')}.`,
  );
  if (przebieg.niedoreczone.length > 0) {
    czesci.push(`Bez doręczenia (${przebieg.niedoreczone.length}): ${przebieg.niedoreczone.join(' ')}`);
  }
  if (przebieg.nieutrwalone.length > 0) {
    czesci.push(
      `Doręczone, ale NIEUTRWALONE w rdzeniu (${przebieg.nieutrwalone.length}) — po restarcie więzi nie będzie: ${przebieg.nieutrwalone.join(' ')}`,
    );
  }
  return czesci.join(' ');
}

/** Funkcja znakuje czynność, której przesłanka jeszcze nie zaszła, bez wygaszania przycisku: znacznik pokazuje powód przed naciśnięciem, zamiast blokować kliknięcie. */
function oznaczCzynnosc(kontrolka: HTMLButtonElement, brakujacaPrzeslanka: string): void {
  if (brakujacaPrzeslanka === '') {
    kontrolka.removeAttribute('title');
    kontrolka.removeAttribute('aria-description');
    delete kontrolka.dataset['przeslanka'];
    return;
  }
  const zdanie = `Czynność nie ruszy, dopóki ${brakujacaPrzeslanka} — naciśnięcie powie to samo wprost.`;
  kontrolka.title = zdanie;
  kontrolka.setAttribute('aria-description', zdanie);
  kontrolka.dataset['przeslanka'] = brakujacaPrzeslanka;
}

/** Funkcja składa treść odmowy z komunikatu i kodu kontraktu, po którym odróżnia się rodzaj odmowy zgłoszonej przez rdzeń. */
function powod(blad?: ErrorInfo): string {
  if (blad === undefined) return 'Rdzeń nie podał przyczyny.';
  return `Powód: ${blad.message} (kod ${blad.code}).`;
}
