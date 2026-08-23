/**
 * Nagrywanie dźwięku z mikrofonu — surowy silnik, bez widoku.
 *
 * Jedna odpowiedzialność: otworzyć strumień, zebrać bajty, oddać je i zamknąć
 * strumień. Plik nie tworzy DOM-u, nie rysuje przycisku, nie wysyła nagrania
 * nigdzie dalej i nie zna transkrypcji — wynik oddaje wywołującemu.
 */

/** Stan nagrywania; `konczy` to czas między poleceniem stopu a gotowymi bajtami. */
export type StanNagrywania = 'bezczynne' | 'nagrywa' | 'konczy';

/** Gotowe nagranie oddane wywołującemu. */
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
 * Kolejność preferencji rodzaju treści.
 *
 * Silnik transkrypcji przyjmuje `.wav .ogg .m4a .webm`. `audio/webm` stoi
 * pierwszy, bo w oknie osadzonym opartym o Chromium bywa jedynym wspieranym
 * zapisem — ale pierwszeństwo to nie pewność, więc każdą pozycję i tak
 * przepuszczamy przez `isTypeSupported`. Kodek podajemy tam,
 * gdzie przeglądarki go wymagają do rozstrzygnięcia; wpis bez kodeka stoi
 * zaraz za nim jako zapasowy.
 */
const PREFEROWANE_RODZAJE: readonly string[] = [
  'audio/webm;codecs=opus',
  'audio/webm',
  'audio/ogg;codecs=opus',
  'audio/ogg',
  'audio/mp4',
  'audio/wav',
];

/** Zdania odmowy — dwa różne powody dostają dwa różne zdania. */
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
 * Najlepszy rodzaj treści obsługiwany przez to okno.
 *
 * Zwraca pusty napis, gdy żadna pozycja wykazu nie przechodzi — wtedy
 * `MediaRecorder` dostaje wybór własny, a `rodzajTresci` bierzemy z niego po
 * fakcie. `audio/webm` nie wraca stąd na wiarę: napis rodzaju, którego zapis
 * nie użył, byłby atrapą podpisu pod bajtami.
 */
export function wybierzRodzajTresci(): string {
  const silnik = (globalThis as { MediaRecorder?: typeof MediaRecorder }).MediaRecorder;
  if (silnik === undefined || typeof silnik.isTypeSupported !== 'function') return '';
  return PREFEROWANE_RODZAJE.find((rodzaj) => silnik.isTypeSupported(rodzaj)) ?? '';
}

/**
 * Zdanie dla błędu `getUserMedia`.
 *
 * Odmowa i brak sprzętu to dwie różne rzeczy. Pierwsza jest decyzją, którą
 * Operator cofa w ustawieniach przeglądarki; druga jest stanem biurka, który
 * naprawia kabel. Wspólne zdanie „nie udało się nagrać” kazałoby szukać po
 * omacku w obu przypadkach.
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

  /**
   * Zamknięcie strumienia.
   *
   * Dopóki ścieżka żyje, system pokazuje, że konsola słucha. Ścieżka zostawiona
   * po nagraniu zapala lampkę mikrofonu, choć nikt nie nagrywa, więc `stop()`
   * idzie po każdej drodze wyjścia — udanej, przerwanej i błędnej.
   */
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

    // Puste `idUrzadzenia` znaczy „wejście domyślne systemu” — nie jest błędem
    // i nie zawęża zapytania, bo `deviceId: ''` odrzuciłoby każde urządzenie.
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
      // Strumień już żyje — bez tego zwolnienia lampka zostałaby zapalona po
      // błędzie, którego Operator nawet nie spowodował.
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
    // Trwanie mierzy zegar, nie rozmiar bajtów: przepływność zależy od kodeka
    // i ciszy, więc dzielenie rozmiaru przez cokolwiek dawałoby liczbę zmyśloną.
    const trwanieMs = Math.max(0, Math.round(performance.now() - poczatek));
    ustawStan('konczy');

    return new Promise<Nagranie | null>((rozstrzygnij) => {
      biezacy.addEventListener('stop', () => {
        const rodzajTresci = biezacy.mimeType !== '' ? biezacy.mimeType : wybierzRodzajTresci();
        const zebrane = kawalki;
        kawalki = [];
        zwolnijSciezki();
        ustawStan('bezczynne');

        // Zero bajtów to nie jest nagranie ciszy — to nagranie, którego nie
        // było. Pusty `Blob` udawałby wypowiedź i poszedłby do transkrypcji.
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
