import type { Kanal } from '../../protokol/kanal';

/**
 * Dostarczenie nagrania do silnika mowy — jedyna luka tego ogniwa.
 *
 * `SpeechTranscribeRequest.audioRef` to ścieżka do pliku na maszynie silnika,
 * a `MediaRecorder` w przeglądarce oddaje `Blob`, czyli bajty w pamięci karty.
 * Kontrakt nie ma komendy przyjmującej dźwięk: `speech.transcribe` bierze
 * gotową ścieżkę, a `translate.speech.synthesize` idzie w stronę odwrotną.
 * Powłoka Tauri wystawia trzy polecenia (`wybierz_katalog_roboczy`,
 * `stan_rdzenia`, `adres_rdzenia`) i nie ma wtyczki systemu plików, więc pliku
 * nie ma czym zapisać.
 *
 * Ogniwo odmawia zamiast obchodzić ten brak: nagranie nie jest nigdzie
 * zapisywane, ścieżka nie jest zmyślana, a dźwięk nie opuszcza maszyny.
 *
 * Drogę otworzy komenda kontraktu przyjmująca nagranie w base64. Zmienia się
 * wtedy wyłącznie `utworzDostarczenie`: stała `KOMENDA_NAGRANIA` dostaje nazwę
 * komendy, `dostepne()` odwraca się samo, a `dostarcz()` woła komendę kanałem
 * i oddaje `audioRef`. Pomocnik `naBase64` jest gotowy, bo kodowanie bajtów to
 * jedyna część tej drogi możliwa do zbudowania bez kontraktu. Poza tym plikiem
 * nie zmienia się nic: `transkrypcja-nagrania.ts` i `dostepnosc-dyktowania.ts`
 * pytają wyłącznie przez ten interfejs.
 */

/** Dostarczenie nagrania na maszynę silnika — zamiana bajtów w `audioRef`. */
export interface Dostarczenie {
  /** Czy klient ma czym dostarczyć nagranie; bez komendy kontraktu `false`. */
  dostepne(): boolean;
  /** Oddaje `audioRef` dla `speech.transcribe`; odmawia `OdmowaDostarczenia`. */
  dostarcz(nagranie: Blob, rodzajTresci: string): Promise<string>;
}

/**
 * Odmowa dostarczenia — osobny typ, żeby powód nie zgubił się po drodze.
 *
 * Zwykły `Error` przeszedłby przez `catch` tak samo, ale wywołujący nie umiałby
 * odróżnić odmowy ogniwa od wyjątku środowiska (np. `TypeError` z pomyłki
 * w kodzie). Pole `zdanie` jest gotowe do pokazania bez obróbki.
 */
export class OdmowaDostarczenia extends Error {
  readonly zdanie: string;

  constructor(zdanie: string) {
    super(zdanie);
    this.name = 'OdmowaDostarczenia';
    this.zdanie = zdanie;
  }
}

/**
 * Zdanie odmowy: co się nie stało, dlaczego i co to zmieni. Stała jest
 * eksportowana, bo sięga po nią sprawdzian.
 */
export const ZDANIE_BRAKU_DROGI =
  'Nagranie nie zostało przekazane silnikowi mowy. ' +
  'Kontrakt przyjmuje w speech.transcribe pole audioRef, czyli ŚCIEŻKĘ PLIKU na maszynie silnika, ' +
  'a nagranie z mikrofonu istnieje wyłącznie jako bajty w pamięci karty — kontrakt nie ma dziś komendy ' +
  'przyjmującej dźwięk, a powłoka Tauri wystawia tylko trzy polecenia (wybierz_katalog_roboczy, ' +
  'stan_rdzenia, adres_rdzenia) i nie ma wtyczki systemu plików, więc pliku nie ma czym zapisać. ' +
  'Zmieni to dopiero dopisanie do kontraktu komendy przyjmującej nagranie w base64 i uchwytu na nią ' +
  'w rdzeniu — do tego czasu dyktowanie pozostaje wyłączone, a tekst wpisuje się z klawiatury.';

/** Nazwa komendy kontraktu przyjmującej nagranie; `null`, dopóki jej nie ma. */
const KOMENDA_NAGRANIA: string | null = null;

export function utworzDostarczenie(kanal: Kanal): Dostarczenie {
  // Kanał jest przyjmowany, choć nie ma dziś czym go użyć: nową komendę poniesie
  // właśnie on, a zmiana sygnatury pociągnęłaby za sobą każde miejsce wywołania.
  void kanal;

  return {
    dostepne: () => KOMENDA_NAGRANIA !== null,

    async dostarcz(nagranie: Blob, rodzajTresci: string): Promise<string> {
      void nagranie;
      void rodzajTresci;
      // Miejsce na wywołanie komendy przyjmującej nagranie. Bez niej ogniwo
      // odmawia: nie zapisuje nagrania, nie zmyśla ścieżki i nie wysyła dźwięku
      // poza maszynę.
      throw new OdmowaDostarczenia(ZDANIE_BRAKU_DROGI);
    },
  };
}

/**
 * Bajty nagrania jako base64 — gotowe pod nową komendę, na razie bez odbiorcy.
 *
 * Przemiana `Blob` → `ArrayBuffer` → base64 dzieje się w całości w karcie i nic
 * nie opuszcza maszyny. Kodowanie idzie porcjami, bo
 * `String.fromCharCode(...bajty)` na nagraniu kilkusekundowym przekracza limit
 * argumentów wywołania i wywraca się na pozornie losowej długości.
 */
export async function naBase64(nagranie: Blob): Promise<string> {
  const bajty = new Uint8Array(await nagranie.arrayBuffer());
  const PORCJA = 0x8000;
  let znaki = '';
  for (let poczatek = 0; poczatek < bajty.length; poczatek += PORCJA) {
    znaki += String.fromCharCode(...bajty.subarray(poczatek, poczatek + PORCJA));
  }
  return btoa(znaki);
}
