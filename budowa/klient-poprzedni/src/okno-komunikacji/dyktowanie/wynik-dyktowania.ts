import type { SpeechTranscribeResponse } from '../../../../shared/contract';

/**
 * Wynik dyktowania — trzy stany, nie dwa.
 *
 * Kontrakt rozdziela w `SpeechTranscribeResponse` dwa niezależne pola:
 * `processed` (czy nagranie przeszło przez silnik) i `transcript` (co silnik
 * w nim usłyszał). Z tej pary wychodzą trzy stany:
 *
 *   processed=true  + transcript niepusty → `rozpoznano`
 *   processed=true  + transcript pusty    → `bez_mowy`
 *   odmowa komendy albo processed=false   → `nieprzetworzone`
 *
 * Cisza jest prawidłowym wynikiem pomiaru: nagranie, w którym nie padło słowo,
 * przeszło przez silnik tak samo jak nagranie z całym zdaniem. Zlanie tego stanu
 * z odmową pokazywałoby awarię tam, gdzie jej nie ma, i rozmywałoby odmowę
 * prawdziwą (brak Pythona, brak modelu, brak drogi dostarczenia). Dlatego
 * `bez_mowy` ma własny stan i własne zdanie.
 *
 * Pole `powod` wypełnia się tylko przy `nieprzetworzone`; przy dwóch pozostałych
 * stanach pomiar się odbył i nie ma czego uzasadniać.
 */

/** Stan wyniku dyktowania — rozstrzyga, co Operator właśnie zobaczył. */
export type StanWyniku = 'rozpoznano' | 'bez_mowy' | 'nieprzetworzone';

/** Wynik jednego dyktowania w kształcie, którym posługuje się widok. */
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
 * Składa wynik z odpowiedzi kontraktu.
 *
 * Jedyne miejsce, w którym pola `processed` i `transcript` zamieniają się
 * w stan. Gdyby każdy widok czytał `processed` sam, cisza prędzej czy później
 * stałaby się w którymś z nich awarią.
 */
export function zlozWynik(odpowiedz: SpeechTranscribeResponse): WynikDyktowania {
  const tekst = odpowiedz.transcript;
  const wspolne = {
    znakow: odpowiedz.characters,
    trwanieMs: odpowiedz.durationMs,
    model: odpowiedz.model,
    jezyk: odpowiedz.language ?? '',
  };

  // `processed=false` przy udanej odpowiedzi znaczy: rdzeń odpowiedział, ale
  // nagrania nie przetworzył. To wciąż stan trzeci, tylko bez zdania odmowy
  // z warstwy transportu — powód nazywamy sami, żeby pole nie zostało puste.
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
    // Cisza: ani tekstu, ani powodu — powód jest polem stanu `nieprzetworzone`,
    // a tu nic nie zawiodło.
    return { stan: 'bez_mowy', tekst: '', ...wspolne, powod: '' };
  }

  return { stan: 'rozpoznano', tekst, ...wspolne, powod: '' };
}

/**
 * Wynik stanu `nieprzetworzone` złożony z samego zdania odmowy.
 *
 * Używany przez ogniwa, które zatrzymały się przed silnikiem — dostarczenie
 * nagrania albo odmowa komendy. Liczby są zerami, bo pomiaru nie było; wzięcie
 * `trwanieMs` z długości nagrania podstawiałoby wartość, której silnik nie
 * zmierzył.
 */
export function wynikNieprzetworzony(powod: string): WynikDyktowania {
  return { stan: 'nieprzetworzone', tekst: '', znakow: 0, trwanieMs: 0, model: '', jezyk: '', powod };
}

/**
 * Zdanie opisujące wynik — po polsku, jedno na każdy z trzech stanów.
 *
 * Zdanie dla `bez_mowy` mówi, co się stało (nagranie przetworzono) i czego
 * w nagraniu nie było (mowy). Nie zawiera słowa „błąd" ani wezwania do
 * ponownego nagrania — cisza nie jest usterką.
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

/** Długość nagrania po polsku — sekundy z jednym miejscem po przecinku. */
function opisTrwania(trwanieMs: number): string {
  if (trwanieMs <= 0) return 'o nieznanej długości';
  const sekundy = (trwanieMs / 1000).toFixed(1).replace('.', ',');
  return `trwającego ${sekundy} s`;
}

/** Dopisek o modelu i języku; pusty, gdy silnik ich nie podał. */
function opisSilnika(wynik: WynikDyktowania): string {
  const czesci: string[] = [];
  if (wynik.model !== '') czesci.push(`model ${wynik.model}`);
  if (wynik.jezyk !== '') czesci.push(`język ${wynik.jezyk}`);
  return czesci.length === 0 ? '' : ` (${czesci.join(', ')})`;
}

/** Odmiana rzeczownika „znak" przy liczbie — zdanie ma brzmieć po polsku. */
function odmianaZnakow(liczba: number): string {
  if (liczba === 1) return 'znak';
  const dziesiatki = liczba % 100;
  const jednosci = liczba % 10;
  if (jednosci >= 2 && jednosci <= 4 && (dziesiatki < 12 || dziesiatki > 14)) return 'znaki';
  return 'znaków';
}
