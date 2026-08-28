import type { SpeechTranscribeResponse } from '../../../../shared/contract';

/** Wynik dyktowania różnicuje trzy stany kontraktu: rozpoznano, bez_mowy oraz nieprzetworzone nagranie. */

/** Stan wyniku dyktowania rozstrzyga, co użytkownik właśnie zobaczył na ekranie po zakończeniu nagrywania i próbie rozpoznania mowy. */
export type StanWyniku = 'rozpoznano' | 'bez_mowy' | 'nieprzetworzone';

/** Wynik jednego dyktowania w kształcie, którym posługuje się widok, niosący stan, tekst, długość nagrania oraz dane modelu. */
export interface WynikDyktowania {
  /** Który z trzech stanów zaszedł. */
  stan: StanWyniku;
  /** Rozpoznany tekst; pusty przy `bez_mowy` i przy `nieprzetworzone`. */
  tekst: string;
  /** Liczba znaków transkrypcji podana przez silnik. */
  znakow: number;
  /** Długość nagrania w milisekundach. */
  trwanieMs: number;
  /** Model, którym rozpoznano; pusty, gdy do rozpoznania nie doszło. */
  model: string;
  /** Język wykryty albo zadany; pusty, gdy silnik go nie podał. */
  jezyk: string;
  /** Zdanie odmowy — wypełnione tylko przy `nieprzetworzone`. */
  powod: string;
}

/**
 * Składanie wyniku z odpowiedzi kontraktu jest jedynym miejscem zamiany pól `processed` i `transcript` w jeden z trzech stanów wyniku.
 */
export function zlozWynik(odpowiedz: SpeechTranscribeResponse): WynikDyktowania {
  const tekst = odpowiedz.transcript;
  const wspolne = {
    znakow: odpowiedz.characters,
    trwanieMs: odpowiedz.durationMs,
    model: odpowiedz.model,
    jezyk: odpowiedz.language ?? '',
  };

  // `processed=false` przy udanej odpowiedzi to wciąż stan nieprzetworzony z dopisanym powodem.
  if (!odpowiedz.processed) {
    return {
      stan: 'nieprzetworzone',
      tekst: '',
      ...wspolne,
      powod:
        'Nagranie nie zostało przetworzone: rdzeń odpowiedział na speech.transcribe, ale zwrócił processed=false, ' +
        'czyli silnik mowy nie wykonał pracy nad tym plikiem — Operator odczyta powód w dzienniku rdzenia albo ' +
        'sprawdzi dostępność silnika w ustawieniach mowy.',
    };
  }

  if (tekst === '') {
    // Cisza: ani tekstu, ani powodu — pole powodu należy do stanu nieprzetworzonego, tu nic nie zawiodło.
    return { stan: 'bez_mowy', tekst: '', ...wspolne, powod: '' };
  }

  return { stan: 'rozpoznano', tekst, ...wspolne, powod: '' };
}

/**
 * Wynik stanu nieprzetworzonego składa się z samego zdania odmowy dla ogniw zatrzymanych przed silnikiem, z licznikami wyzerowanymi wobec braku pomiaru.
 */
export function wynikNieprzetworzony(powod: string): WynikDyktowania {
  return { stan: 'nieprzetworzone', tekst: '', znakow: 0, trwanieMs: 0, model: '', jezyk: '', powod };
}

/**
 * Zdanie opisujące wynik po polsku obejmuje jeden komunikat na każdy z trzech stanów, a dla braku mowy nazywa fakt bez słowa błąd.
 */
export function zdanieWyniku(wynik: WynikDyktowania): string {
  switch (wynik.stan) {
    case 'rozpoznano':
      return `Rozpoznano ${wynik.znakow} ${odmianaZnakow(wynik.znakow)} z nagrania ${opisTrwania(wynik.trwanieMs)}${opisSilnika(wynik)}.`;
    case 'bez_mowy':
      return `Nagranie ${opisTrwania(wynik.trwanieMs)} zostało przetworzone i mowy w nim nie było${opisSilnika(wynik)}.`;
    case 'nieprzetworzone':
      return wynik.powod === '' ? 'Nagranie nie zostało przetworzone; powodu nie podano.' : wynik.powod;
  }
}

/** Długość nagrania podana po polsku jako liczba sekund z jednym miejscem po przecinku, czytelna dla użytkownika interfejsu. */
function opisTrwania(trwanieMs: number): string {
  if (trwanieMs <= 0) return 'o nieznanej długości';
  const sekundy = (trwanieMs / 1000).toFixed(1).replace('.', ',');
  return `trwającego ${sekundy} s`;
}

/** Dopisek o modelu rozpoznawania i wykrytym języku pozostaje pusty, gdy silnik transkrypcji ich nie podał w odpowiedzi. */
function opisSilnika(wynik: WynikDyktowania): string {
  const czesci: string[] = [];
  if (wynik.model !== '') czesci.push(`model ${wynik.model}`);
  if (wynik.jezyk !== '') czesci.push(`język ${wynik.jezyk}`);
  return czesci.length === 0 ? '' : ` (${czesci.join(', ')})`;
}

/** Odmiana rzeczownika „znak” dobierana jest do liczby, aby zdanie o długości transkrypcji brzmiało poprawnie po polsku. */
function odmianaZnakow(liczba: number): string {
  if (liczba === 1) return 'znak';
  const dziesiatki = liczba % 100;
  const jednosci = liczba % 10;
  if (jednosci >= 2 && jednosci <= 4 && (dziesiatki < 12 || dziesiatki > 14)) return 'znaki';
  return 'znaków';
}
