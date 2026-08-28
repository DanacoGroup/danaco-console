import type { Kanal } from '../../protokol/kanal';

/**
 * Dostarczenie nagrania na maszynę silnika mowy, zamieniające bajty nagrania na ścieżkę
 * pliku oczekiwaną przez komendę transkrypcji.
 */
export interface Dostarczenie {
  /** Czy klient ma czym dostarczyć nagranie; bez komendy kontraktu `false`. */
  dostepne(): boolean;
  /** Oddaje `audioRef` dla `speech.transcribe`; odmawia `OdmowaDostarczenia`. */
  dostarcz(nagranie: Blob, rodzajTresci: string): Promise<string>;
}

/**
 * Odmowa dostarczenia nagrania: osobny typ błędu, żeby wywołujący odróżnił świadomą
 * odmowę ogniwa od wyjątku środowiska, z gotowym do pokazania zdaniem.
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

/**
 * Nazwa komendy kontraktu przyjmującej nagranie w formacie base64; pozostaje pusta,
 * dopóki taka komenda nie zostanie dodana do kontraktu.
 */
const KOMENDA_NAGRANIA: string | null = null;

export function utworzDostarczenie(kanal: Kanal): Dostarczenie {
  // Kanał jest przyjmowany bez użycia dziś, bo nową komendę poniesie właśnie on, gdy powstanie.
  void kanal;

  return {
    dostepne: () => KOMENDA_NAGRANIA !== null,

    async dostarcz(nagranie: Blob, rodzajTresci: string): Promise<string> {
      void nagranie;
      void rodzajTresci;
      // Miejsce na wywołanie komendy przyjmującej nagranie; do czasu jej dodania ogniwo odmawia.
      throw new OdmowaDostarczenia(ZDANIE_BRAKU_DROGI);
    },
  };
}

/**
 * Zamienia bajty nagrania na tekst base64 porcjami, żeby uniknąć przekroczenia limitu
 * argumentów wywołania na dłuższych nagraniach.
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
