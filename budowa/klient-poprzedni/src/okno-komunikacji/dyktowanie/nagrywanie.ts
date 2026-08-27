/**
 * Nagrywanie to surowy silnik obsługi mikrofonu bez widoku: otwiera strumień, zbiera bajty, zamyka go i oddaje wynik wywołującemu, nie tworząc interfejsu ani nie znając transkrypcji.
 */

/** Stan nagrywania — pole `konczy` mierzy czas od polecenia zatrzymania do chwili, w której bajty nagrania są gotowe do odczytu. */
export type StanNagrywania = 'bezczynne' | 'nagrywa' | 'konczy';

/** Gotowe nagranie oddane wywołującemu, niosące bajty zapisu wraz z rodzajem treści i zmierzonym czasem trwania nagrywania. */
export interface Nagranie {
  bajty: Blob;
  /** Rodzaj MIME faktycznie użyty przez `MediaRecorder`. */
  rodzajTresci: string;
  /** Czas mierzony zegarem od startu do stopu, w milisekundach. */
  trwanieMs: number;
}

export interface Nagrywanie {
  stan(): StanNagrywania;
  /** Otwiera strumień i zaczyna zapis; rzuca błąd ze zdaniem po polsku. */
  rozpocznij(idUrzadzenia: string): Promise<void>;
  /** Domyka zapis; `null`, gdy nie było czego domykać. */
  zakoncz(): Promise<Nagranie | null>;
  /** Porzuca nagranie — bajty idą do kosza, wynik nie powstaje. */
  przerwij(): void;
  naStan(sluchacz: (stan: StanNagrywania) => void): () => void;
}

/**
 * Kolejność preferencji rodzaju treści odzwierciedla wsparcie silnika transkrypcji, sprawdzane każdorazowo metodą `isTypeSupported` przed użyciem.
 */
const PREFEROWANE_RODZAJE: readonly string[] = [
  'audio/webm;codecs=opus',
  'audio/webm',
  'audio/ogg;codecs=opus',
  'audio/ogg',
  'audio/mp4',
  'audio/wav',
];

/** Zdania odmowy nagrywania — brak zgody użytkownika i brak sprzętu audio dostają osobne, rozróżnialne komunikaty błędu. */
const POWOD_ODMOWA_ZGODY =
  'Nie ma zgody na mikrofon. Przeglądarka odmówiła dostępu do dźwięku dla tej ' +
  'konsoli. Otwórz ustawienia witryny w przeglądarce i zezwól na mikrofon, a potem ' +
  'spróbuj nagrać jeszcze raz.';

const POWOD_BRAK_URZADZENIA =
  'Nie znaleziono mikrofonu. Wskazane wejście dźwięku zniknęło albo nigdy nie było ' +
  'podłączone. Wepnij mikrofon i wybierz go z wykazu urządzeń.';

const POWOD_BRAK_SILNIKA =
  'Ta przeglądarka nie umie nagrywać dźwięku. Brakuje jej zapisu strumienia, którego ' +
  'używa dyktowanie. Otwórz konsolę pod adresem „https” w nowszym oknie przeglądarki.';

const POWOD_NIEZNANY =
  'Nagrywanie nie ruszyło. Przeglądarka przerwała otwieranie mikrofonu bez podania ' +
  'rozpoznanego powodu. Spróbuj ponownie, a jeśli błąd wraca — wybierz inne wejście ' +
  'dźwięku z wykazu.';

/**
 * Najlepszy rodzaj treści obsługiwany przez to okno, ustalony metodą `isTypeSupported`, a nie zakładany na podstawie samej preferencji `audio/webm`.
 */
export function wybierzRodzajTresci(): string {
  const silnik = (globalThis as { MediaRecorder?: typeof MediaRecorder }).MediaRecorder;
  if (silnik === undefined || typeof silnik.isTypeSupported !== 'function') return '';
  return PREFEROWANE_RODZAJE.find((rodzaj) => silnik.isTypeSupported(rodzaj)) ?? '';
}

/**
 * Zdanie błędu dla `getUserMedia` rozróżnia odmowę zgody użytkownika od braku sprzętu audio, ponieważ każda z nich wymaga innej naprawy.
 */
export function opiszOdmoweNagrywania(blad: unknown): string {
  const nazwa = (blad as { name?: string } | null)?.name ?? '';
  if (nazwa === 'NotAllowedError' || nazwa === 'SecurityError') return POWOD_ODMOWA_ZGODY;
  if (nazwa === 'NotFoundError' || nazwa === 'OverconstrainedError') return POWOD_BRAK_URZADZENIA;
  return POWOD_NIEZNANY;
}

export function utworzNagrywanie(): Nagrywanie {
  const sluchacze = new Set<(stan: StanNagrywania) => void>();
  let stanBiezacy: StanNagrywania = 'bezczynne';
  let zapis: MediaRecorder | null = null;
  let strumien: MediaStream | null = null;
  let kawalki: Blob[] = [];
  let poczatek = 0;
  let porzucone = false;

  function ustawStan(nowy: StanNagrywania): void {
    if (stanBiezacy === nowy) return;
    stanBiezacy = nowy;
    for (const sluchacz of [...sluchacze]) {
      try {
        sluchacz(nowy);
      } catch (blad) {
        console.error('[dyktowanie] błąd słuchacza stanu', blad);
      }
    }
  }

  /** Zamknięcie strumienia biegnie każdą drogą wyjścia, by nie zostawić zapalonej lampki mikrofonu. */
  function zwolnijSciezki(): void {
    strumien?.getTracks().forEach((sciezka) => sciezka.stop());
    strumien = null;
    zapis = null;
  }

  async function rozpocznij(idUrzadzenia: string): Promise<void> {
    if (stanBiezacy !== 'bezczynne') return;

    const nawigator = (globalThis as { navigator?: Navigator }).navigator;
    const silnik = (globalThis as { MediaRecorder?: typeof MediaRecorder }).MediaRecorder;
    if (nawigator?.mediaDevices?.getUserMedia === undefined || silnik === undefined) {
      throw new Error(POWOD_BRAK_SILNIKA);
    }

    // Puste `idUrzadzenia` oznacza urządzenie domyślne, nie pusty ciąg w zapytaniu `deviceId`.
    const wiezy: MediaStreamConstraints = {
      audio: idUrzadzenia === '' ? true : { deviceId: { exact: idUrzadzenia } },
    };

    let otwarty: MediaStream;
    try {
      otwarty = await nawigator.mediaDevices.getUserMedia(wiezy);
    } catch (blad) {
      throw new Error(opiszOdmoweNagrywania(blad));
    }

    const rodzaj = wybierzRodzajTresci();
    strumien = otwarty;
    kawalki = [];
    porzucone = false;
    try {
      zapis = new silnik(otwarty, rodzaj === '' ? undefined : { mimeType: rodzaj });
    } catch (blad) {
      // Zwolnienie żyjącego strumienia zapobiega zapalonej lampce mikrofonu po błędzie Operatora.
      zwolnijSciezki();
      console.error('[dyktowanie] zapis nie ruszył', blad);
      throw new Error(POWOD_BRAK_SILNIKA);
    }

    zapis.addEventListener('dataavailable', (zdarzenie) => {
      if (zdarzenie.data.size > 0) kawalki.push(zdarzenie.data);
    });
    poczatek = performance.now();
    zapis.start();
    ustawStan('nagrywa');
  }

  function zakoncz(): Promise<Nagranie | null> {
    if (stanBiezacy !== 'nagrywa' || zapis === null) return Promise.resolve(null);
    const biezacy = zapis;
    // Czas trwania mierzy zegar, nie rozmiar bajtów — przepływność zależy od kodeka i ciszy w nagraniu.
    const trwanieMs = Math.max(0, Math.round(performance.now() - poczatek));
    ustawStan('konczy');

    return new Promise<Nagranie | null>((rozstrzygnij) => {
      biezacy.addEventListener('stop', () => {
        const rodzajTresci = biezacy.mimeType !== '' ? biezacy.mimeType : wybierzRodzajTresci();
        const zebrane = kawalki;
        kawalki = [];
        zwolnijSciezki();
        ustawStan('bezczynne');

        // Zero bajtów oznacza brak nagrania, nie ciszę — pusty `Blob` nie idzie do transkrypcji.
        if (porzucone || zebrane.length === 0) {
          rozstrzygnij(null);
          return;
        }
        rozstrzygnij({ bajty: new Blob(zebrane, { type: rodzajTresci }), rodzajTresci, trwanieMs });
      });
      biezacy.stop();
    });
  }

  function przerwij(): void {
    if (stanBiezacy === 'bezczynne') return;
    porzucone = true;
    kawalki = [];
    try {
      zapis?.stop();
    } catch (blad) {
      console.error('[dyktowanie] przerwanie zapisu', blad);
    }
    zwolnijSciezki();
    ustawStan('bezczynne');
  }

  return {
    stan: () => stanBiezacy,
    rozpocznij,
    zakoncz,
    przerwij,
    naStan(sluchacz) {
      sluchacze.add(sluchacz);
      return () => {
        sluchacze.delete(sluchacz);
      };
    },
  };
}
