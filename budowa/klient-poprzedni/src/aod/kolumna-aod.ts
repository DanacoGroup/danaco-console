import './aod.css';

import { EventType } from '../../../shared/contract';
import { utworzRameOkna } from '../komponenty/rama-okna';
import { utworzStanTresci } from '../komponenty/stan-tresci';
import { otworzOknoKonfiguracji } from '../konfiguracja/indeks';
import type { Odsubskrybuj } from '../polaczenie/magistrala-zdarzen';
import type { Kanal } from '../protokol/kanal';
import type { Queue } from '../../../shared/contract';
import type { OpisCzterechSterow, TozsamoscDlaOgniska } from './cztery-stery';
import type { KolejkaDecyzji } from './kolejka-decyzji';
import { opisOdmowyAod } from './odmowy-aod';
import { obserwacjaZeStanuOkna, obserwacjaZPostepu } from './rozpoznanie-decyzji';
import { utworzSekcjeDecyzji } from './sekcja-decyzji';
import { utworzSekcjeKontekstu } from './sekcja-kontekstu';
import { utworzSekcjeObecnosci } from './sekcja-obecnosci';
import { utworzSekcjePodpowiedzi } from './sekcja-podpowiedzi';
import { utworzSekcjeGlosu } from './sekcja-glosu';
import { utworzSekcjeRozmowy } from './sekcja-rozmowy';
import { utworzSekcjeStanu } from './sekcja-stanu';
import { utworzSterowanieObecnoscia } from './sterowanie-obecnoscia';
import type { StanObecnosci } from './tryb-obecnosci';
import type { KontekstWyciszenia } from './wyciszenie-kontekst';
import { utworzZrodloDecyzji } from './zrodlo-decyzji';
import { utworzZrodloAod } from './zrodlo-komend';

/**
 * Powierzchnia interakcji Always On Display — KOLUMNA BOCZNA, nie okno.
 *
 * Opracowanie (`docs/funkcje-globalne/always-on-display.md`, rozdz. 2.1, 2.3,
 * 2.6, 8.3) mówi wprost: funkcja nie ma własnego okna. Całą jej powierzchnią
 * jest awatar oraz powierzchnia otwierana jako ROZSZERZENIE BOCZNE po prawej
 * stronie obszaru roboczego — kolumna na pełną wysokość, o regulowanej
 * szerokości. Otwarcie zwęża kolumny obszaru roboczego, nie przesłania ich
 * i nie zamyka żadnego z okien komunikacji operacyjnej.
 *
 * Stąd trzy różnice wobec modala: nie ma przyciemnienia tła, nie ma pułapki
 * ogniska i nie ma zabranej pracy pod spodem. Kolumna leży na warstwie AOD
 * (`--dn-z-aod`, 1200 wg `design/KANON.md`), a nie na warstwie modala.
 *
 * Kolumna woła `aod.status.get`, `aod.context.get`, `aod.suggestion`,
 * `aod.observe.attach`, `aod.observe.detach`, `aod.chat.send`
 * i `aod.voice.command`; obsługiwacze stoją w `core/handlers_aod.go`.
 *
 * Ten plik wyłącznie składa powierzchnię: rama z biblioteki, stan treści
 * z biblioteki, sekcje z osobnych plików, sterowanie obecnością
 * z `sterowanie-obecnoscia.ts` i dwa źródła komend (`zrodlo-komend` — rodzina
 * `aod.*`, `zrodlo-decyzji` — komendy cudzych rodzin potrzebne kolejce decyzji).
 *
 * Kolumna subskrybuje dwa zdarzenia i zdejmuje subskrypcje w `rozlacz()`:
 *   • `progress.changed` — rozgłaszane do każdego połączenia, niesie proces,
 *     etap, procent, stan i okno;
 *   • `window.state.changed` — jedyny żywy nośnik `LoopState`, bo pola `loop`
 *     w `progress.changed` telemetria nie wypełnia.
 * To, co przeleciało przed otwarciem kolumny, nadrabia `monitor.status`.
 *
 * Subskrypcje żyją między zbudowaniem warstwy a `rozlacz()`, nie między
 * otwarciem a zamknięciem kolumny: proces zmienia stan niezależnie od tego, czy
 * Operator patrzy. Dzięki temu plakietka awatara jest prawdziwa także wtedy,
 * gdy kolumna stoi zamknięta.
 */
export interface KolumnaAod {
  /** Kolumna do osadzenia w warstwie AOD. */
  element: HTMLElement;
  otworz(): void;
  zamknij(): void;
  przelacz(): void;
  czyOtwarta(): boolean;
  /** Szerokość zajmowana przez kolumnę; 0, gdy zamknięta. */
  szerokosc(): number;
  /** Kolejka decyzji — warstwa dosypuje do niej sygnały i czyta z niej plakietkę. */
  kolejka: KolejkaDecyzji;
  /** Opis sterów wspólny z dymkiem kontekstowym. */
  stery: OpisCzterechSterow;
  /** Kolejka rdzenia dopasowana do decyzji. */
  dopasujKolejke(idOkna?: string, idSesji?: string): Queue | undefined;
  /** Przerysowuje wykaz decyzji bez pytania rdzenia. */
  przerysujDecyzje(): void;
  /** Ponawia odczyt kolejki decyzji z rdzenia. */
  odswiezDecyzje(): Promise<void>;
  /** Prowadzi ognisko na pole toru głosowego — skrót `Ctrl/Cmd + Shift + V`. */
  ogniskujGlos(): void;
  /** Prowadzi ognisko na listę oczekujących — skrót `Ctrl/Cmd + Shift + Q`. */
  ogniskujSugestie(): void;
  /** Zgłasza obserwatora zmiany szerokości i stanu otwarcia. */
  obserwujUklad(sluchacz: () => void): () => void;
  rozlacz(): void;
}

/** Nastawy powierzchni wykraczające poza kanał. */
export interface OpisKolumnyAod {
  /**
   * Tożsamość klienta Z POWITANIA POŁĄCZENIA.
   *
   * Pole jest opcjonalne, bo wołacz podaje sam kanał. Bez tożsamości czynny
   * zostaje cały ster przejęcia poza jedną rzeczą — przestawieniem ogniska —
   * i powierzchnia mówi wprost, czego nie robi (`cztery-stery.ts`).
   */
  klient?: TozsamoscDlaOgniska;
}

/** Szerokość początkowa kolumny w pikselach — punkt wyjścia regulacji. */
const SZEROKOSC_POCZATKOWA = 420;

/** Najwęższa dopuszczalna szerokość kolumny. */
const SZEROKOSC_MIN = 280;

/** Najszersza dopuszczalna szerokość kolumny — połowa okna, nigdy więcej. */
function szerokoscMax(): number {
  return Math.max(SZEROKOSC_MIN, Math.round(globalThis.innerWidth / 2));
}

export function utworzKolumneAod(
  kanal: Kanal,
  stanObecnosci: StanObecnosci,
  opisKolumny: OpisKolumnyAod = {},
  kontekstWyciszenia?: KontekstWyciszenia,
): KolumnaAod {
  const zrodlo = utworzZrodloAod(kanal);
  const zrodloDecyzji = utworzZrodloDecyzji(kanal);
  const teraz = (): number => Date.now();

  const rama = utworzRameOkna({
    kod: 'aod',
    tytul: 'Always On Display',
    rola: 'monitor',
    modul: 'AOD',
    przeznaczenie:
      'Powierzchnia interakcji funkcji globalnej otwarta jako rozszerzenie boczne — rozmowa ' +
      'z agentem towarzyszącym, lista oczekujących sugestii, tor głosowy i komplet kontekstu.',
    przedrostek: 'ao',
  });

  const stany = utworzStanTresci('ao');
  rama.cialo.append(stany.element);

  const zamelduj = (zdanie: string, udane: boolean): void => stany.potwierdzenie(zdanie, udane);

  const sekcjaDecyzji = utworzSekcjeDecyzji({
    zrodlo: zrodloDecyzji,
    zamelduj,
    otworzOknoKonfiguracji: () => void otworzOknoKonfiguracji(kanal),
    ...(opisKolumny.klient === undefined ? {} : { klient: opisKolumny.klient }),
  });

  const sekcjaStanu = utworzSekcjeStanu();
  const sekcjaPodpowiedzi = utworzSekcjePodpowiedzi(zrodlo);
  const sekcjaObecnosci = utworzSekcjeObecnosci({ zrodlo, zamelduj });
  const sekcjaRozmowy = utworzSekcjeRozmowy({ zrodlo, zamelduj });
  const sekcjaGlosu = utworzSekcjeGlosu({ zrodlo, zamelduj });
  const sekcjaKontekstu = utworzSekcjeKontekstu();

  /**
   * Trzon treści budowany raz i wstawiany po każdym odczycie.
   *
   * `stany.tresc()` czyści miejsce treści, ale sekcje są tymi samymi
   * elementami — wpisana i niewysłana wiadomość ani identyfikator procesu
   * nie giną przy ponownym odczycie.
   */
  const trzon = document.createElement('div');
  trzon.className = 'ao-trzon';
  trzon.append(
    // Decyzje na początku: sprawa wymagająca reakcji ma być widoczna od razu,
    // a nie po przewinięciu sekcji opisujących bieg normalny.
    sekcjaDecyzji.element,
    sekcjaStanu.element,
    sekcjaObecnosci.element,
    sekcjaPodpowiedzi.element,
    sekcjaRozmowy.element,
    sekcjaGlosu.element,
    sekcjaKontekstu.element,
  );

  const odswiez = document.createElement('button');
  odswiez.type = 'button';
  odswiez.className = 'dn-btn dn-btn--zarys';
  odswiez.textContent = 'Odczytaj ponownie';
  odswiez.addEventListener('click', () => void wczytaj());
  rama.narzedzia.append(odswiez);

  // --- nagłówek powierzchni: tryb obecności, tor głosowy, menu kebab ------
  // Rozdz. 2.3 opracowania: `[ tryb ▼ ]  [ 🎙 ]  [ ⋮ ]` w nagłówku.

  // Menu kebab nagłówka jest tym samym komponentem, który stoi przy awatarze;
  // kontekst wyciszenia kontekstowego dojeżdża tu z warstwy, bo tam mieszka
  // jedna kopia odczytu (`wyciszenie-kontekst.ts`).
  const sterowanie = utworzSterowanieObecnoscia(stanObecnosci, teraz, {
    ...(kontekstWyciszenia === undefined ? {} : { kontekst: kontekstWyciszenia }),
    zamelduj: (zdanie: string) => zamelduj(zdanie, true),
    naOtwarcieMenu: () => {
      if (kontekstWyciszenia !== undefined) void kontekstWyciszenia.odswiez();
    },
  });

  const glos = document.createElement('button');
  glos.type = 'button';
  glos.className = 'dn-btn-ikona ao-glos-skrot';
  glos.textContent = '🎙';
  glos.setAttribute(
    'aria-label',
    'Tor głosowy — prowadzi do pola polecenia głosowego (aod.voice.command)',
  );
  glos.title =
    'Tor głosowy funkcji. Nagrania nakładka nie wytwarza — kontrakt nie ma komendy przesyłu ' +
    'nagrania z przeglądarki. Polecenie jedzie treścią, komendą aod.voice.command.';
  glos.addEventListener('click', () => ogniskujGlos());

  rama.pasek.append(sterowanie.przelacznik, glos, sterowanie.kebab);

  // --- obudowa kolumny ----------------------------------------------------

  const element = document.createElement('aside');
  element.className = 'ao-kolumna';
  element.hidden = true;
  element.setAttribute('aria-label', 'Always On Display — powierzchnia interakcji');

  /** Uchwyt regulacji szerokości; kolumna jest rozszerzeniem o regulowanej szerokości. */
  const uchwyt = document.createElement('div');
  uchwyt.className = 'ao-kolumna__uchwyt';
  uchwyt.setAttribute('role', 'separator');
  uchwyt.setAttribute('aria-orientation', 'vertical');
  uchwyt.setAttribute('aria-label', 'Szerokość kolumny Always On Display');
  uchwyt.tabIndex = 0;

  const cialo = document.createElement('div');
  cialo.className = 'ao-kolumna__cialo';
  cialo.append(rama.element);

  const zamknij = document.createElement('button');
  zamknij.type = 'button';
  zamknij.className = 'dn-btn dn-btn--zarys ao-kolumna__zamknij';
  zamknij.textContent = 'Zwiń powierzchnię';
  zamknij.title = 'Zamknięcie przywraca stan spoczynku bez utraty historii rozmowy.';
  zamknij.addEventListener('click', () => zamknijKolumne());

  const stopka = document.createElement('footer');
  stopka.className = 'ao-kolumna__stopka';
  stopka.append(zamknij);

  element.append(uchwyt, cialo, stopka);

  let szerokosc = SZEROKOSC_POCZATKOWA;
  element.style.setProperty('--ao-szerokosc', `${szerokosc}px`);

  const sluchaczeUkladu = new Set<() => void>();
  function rozglosUklad(): void {
    for (const sluchacz of sluchaczeUkladu) sluchacz();
  }

  // --- regulacja szerokości ----------------------------------------------

  function ustawSzerokosc(nowa: number): void {
    const ograniczona = Math.min(szerokoscMax(), Math.max(SZEROKOSC_MIN, Math.round(nowa)));
    if (ograniczona === szerokosc) return;
    szerokosc = ograniczona;
    element.style.setProperty('--ao-szerokosc', `${szerokosc}px`);
    rozglosUklad();
  }

  function naRuch(zdarzenie: PointerEvent): void {
    // Kolumna wyrasta od prawej krawędzi, więc szerokość to odległość wskaźnika
    // od tej krawędzi.
    ustawSzerokosc(globalThis.innerWidth - zdarzenie.clientX);
  }

  function naKoniec(): void {
    document.removeEventListener('pointermove', naRuch);
    document.removeEventListener('pointerup', naKoniec);
  }

  uchwyt.addEventListener('pointerdown', (zdarzenie) => {
    zdarzenie.preventDefault();
    document.addEventListener('pointermove', naRuch);
    document.addEventListener('pointerup', naKoniec);
  });

  // Regulacja z klawiatury — kolumna ma być szeroka także bez wskaźnika.
  uchwyt.addEventListener('keydown', (zdarzenie) => {
    if (zdarzenie.key === 'ArrowLeft') ustawSzerokosc(szerokosc + 24);
    else if (zdarzenie.key === 'ArrowRight') ustawSzerokosc(szerokosc - 24);
    else return;
    zdarzenie.preventDefault();
  });

  function naZmianeOkna(): void {
    ustawSzerokosc(szerokosc);
    rozglosUklad();
  }

  globalThis.addEventListener('resize', naZmianeOkna);

  // --- odczyt -------------------------------------------------------------

  /**
   * Odczyt powierzchni.
   *
   * Odmowa jednego odczytu nie wywraca powierzchni: każda sekcja niesie swoją
   * odmowę u siebie, a kolejka decyzji — stojąca na komendach innych rodzin —
   * nie znika razem z odczytem rodziny `aod.*`.
   */
  async function wczytaj(): Promise<void> {
    stany.ladowanie('Pytam rdzeń o stan nakładki (aod.status.get) i o decyzje czekające…');

    const idSesji = zrodlo.idSesji();
    sekcjaRozmowy.ustawSesje(idSesji);
    sekcjaGlosu.ustawSesje(idSesji);

    // Trzon wchodzi przed odczytami: sekcje mają gdzie nanieść i treść,
    // i odmowę, a powierzchnia jest widoczna od razu.
    stany.tresc().replaceChildren(trzon);

    // Odczyt nadrabiający kolejki decyzji — niezależny od rodziny `aod.*`.
    await sekcjaDecyzji.odswiez();

    const wynikStanu = await zrodlo.stan();
    if (!wynikStanu.udany || wynikStanu.wynik === undefined) {
      const zdanie = opisOdmowyAod('Odczyt stanu nakładki', 'stan', wynikStanu.blad);
      sekcjaStanu.odmowa(zdanie);
      sekcjaObecnosci.odmowa(zdanie);
      stany.potwierdzenie(zdanie, false);
      return;
    }

    const status = wynikStanu.wynik.status;
    sekcjaStanu.pokaz(status);
    sekcjaObecnosci.pokaz(status);

    const wynikKontekstu = await zrodlo.kontekst({
      windowId: status.activeWindowId,
      sessionId: idSesji,
      historyLimit: 20,
    });
    sekcjaKontekstu.pokaz(wynikKontekstu);

    stany.potwierdzenie('Stan nakładki odczytany.', true);

    await sekcjaPodpowiedzi.odswiez({
      windowId: status.activeWindowId,
      sessionId: idSesji,
    });
  }

  // --- subskrypcje zdarzeń rdzenia ---------------------------------------

  const subskrypcje: Odsubskrybuj[] = [
    kanal.naZdarzenie(EventType.ProgressChanged, (tresc) => {
      if (sekcjaDecyzji.kolejka.nanies(obserwacjaZPostepu(tresc, Date.now()), Date.now())) {
        sekcjaDecyzji.przerysuj();
        rozglosUklad();
      }
    }),
    kanal.naZdarzenie(EventType.WindowStateChanged, (tresc) => {
      const obserwacja = obserwacjaZeStanuOkna(tresc, Date.now());
      if (obserwacja === null) return;
      if (sekcjaDecyzji.kolejka.nanies(obserwacja, Date.now())) {
        sekcjaDecyzji.przerysuj();
        rozglosUklad();
      }
    }),
  ];

  function ogniskujGlos(): void {
    otworzKolumne();
    const pole = sekcjaGlosu.element.querySelector('textarea');
    if (pole instanceof HTMLTextAreaElement) pole.focus();
  }

  function ogniskujSugestie(): void {
    otworzKolumne();
    sekcjaDecyzji.element.scrollIntoView({ block: 'start' });
  }

  function otworzKolumne(): void {
    if (!element.hidden) return;
    element.hidden = false;
    sterowanie.odswiez(teraz());
    rozglosUklad();
    void wczytaj();
  }

  function zamknijKolumne(): void {
    if (element.hidden) return;
    element.hidden = true;
    rozglosUklad();
  }

  return {
    element,

    otworz: otworzKolumne,
    zamknij: zamknijKolumne,
    przelacz: () => (element.hidden ? otworzKolumne() : zamknijKolumne()),
    czyOtwarta: () => !element.hidden,
    szerokosc: () => (element.hidden ? 0 : szerokosc),

    kolejka: sekcjaDecyzji.kolejka,
    stery: sekcjaDecyzji.stery,
    dopasujKolejke: sekcjaDecyzji.dopasujKolejke,
    przerysujDecyzje: sekcjaDecyzji.przerysuj,
    odswiezDecyzje: sekcjaDecyzji.odswiez,

    ogniskujGlos,
    ogniskujSugestie,

    obserwujUklad(sluchacz) {
      sluchaczeUkladu.add(sluchacz);
      return () => sluchaczeUkladu.delete(sluchacz);
    },

    rozlacz() {
      for (const odsubskrybuj of subskrypcje) odsubskrybuj();
      subskrypcje.length = 0;
      globalThis.removeEventListener('resize', naZmianeOkna);
      naKoniec();
      sterowanie.rozlacz();
      element.remove();
    },
  };
}
