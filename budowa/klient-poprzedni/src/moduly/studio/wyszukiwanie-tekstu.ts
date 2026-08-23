/**
 * Wyszukiwanie i zamiana w treści dokumentu — rachunek na napisie, bez DOM.
 *
 * Czynność jest bliźniacza wobec narzędzi znacznikowych: pracuje na buforze
 * edytora, więc rdzenia nie potrzebuje. Komenda `studio.diff.compare` też umie
 * szukać wzorca, ale szuka go w treści WERSJI zapisanej w repozytorium — a to
 * inna rzecz niż wzorzec w tekście, który Operator właśnie pisze i którego
 * jeszcze nie zapisał. Oba wyszukiwania istnieją obok siebie z zamysłem
 * i każde mówi, po czym szuka.
 *
 * Wzorzec niepoprawny nie jest wyciszany. `RegExp` rzuca na złej składni,
 * a przechwycenie tego bez słowa zamieniłoby literówkę Operatora w ciszę
 * wyglądającą jak „brak trafień". Wynik niesie więc powód wprost.
 */

/** Jedno trafienie wzorca w treści. */
export interface TrafienieTekstu {
  /** Położenie początku trafienia w treści. */
  poczatek: number;
  /** Położenie końca trafienia w treści. */
  koniec: number;
  /** Dopasowany fragment. */
  tekst: string;
}

/** Nastawy wyszukiwania odpowiadające przełącznikom panelu. */
export interface NastawyWyszukiwania {
  wzorzec: string;
  /** Czy wzorzec jest wyrażeniem regularnym; `false` znaczy frazę dosłowną. */
  regularne: boolean;
  /** Czy wielkość liter ma znaczenie. */
  wielkoscLiter: boolean;
}

/** Wynik wyszukiwania albo powód, dla którego wyszukiwania nie da się wykonać. */
export interface WynikWyszukiwania {
  trafienia: readonly TrafienieTekstu[];
  /** Powód niepowodzenia; pusty, gdy wyszukiwanie się wykonało. */
  powod: string;
}

/** Wynik zamiany wraz z liczbą podmienionych trafień. */
export interface WynikZamiany {
  tresc: string;
  liczba: number;
  powod: string;
}

/**
 * Buduje wyrażenie z nastaw albo oddaje powód, dla którego się nie da.
 *
 * Fraza dosłowna jest osłaniana znak po znaku, więc kropka we frazie znaczy
 * kropkę, a nie „dowolny znak". Bez osłony wyszukanie „wersja 2.0" trafiałoby
 * także w „wersja 250".
 */
function zbudujWyrazenie(nastawy: NastawyWyszukiwania): RegExp | string {
  const zrodlo = nastawy.regularne ? nastawy.wzorzec : oslon(nastawy.wzorzec);
  const znaczniki = nastawy.wielkoscLiter ? 'gu' : 'giu';
  try {
    return new RegExp(zrodlo, znaczniki);
  } catch (blad) {
    const powod = blad instanceof Error ? blad.message : String(blad);
    return `Wzorzec nie jest poprawnym wyrażeniem regularnym: ${powod}`;
  }
}

/** Osłania znaki o znaczeniu składniowym, żeby fraza znaczyła samą siebie. */
function oslon(fraza: string): string {
  return fraza.replace(/[.*+?^${}()|[\]\\]/gu, '\\$&');
}

/** Znajduje wszystkie trafienia wzorca w treści. */
export function znajdzTrafienia(
  tresc: string,
  nastawy: NastawyWyszukiwania,
): WynikWyszukiwania {
  if (nastawy.wzorzec === '') return { trafienia: [], powod: '' };
  const wyrazenie = zbudujWyrazenie(nastawy);
  if (typeof wyrazenie === 'string') return { trafienia: [], powod: wyrazenie };

  const trafienia: TrafienieTekstu[] = [];
  for (const dopasowanie of tresc.matchAll(wyrazenie)) {
    const poczatek = dopasowanie.index;
    if (poczatek === undefined) continue;
    // Wzorzec dopasowujący pustkę („a*") przesunąłby pętlę o zero znaków
    // i zapełnił wykaz trafieniami bez treści. Pomijamy je, zamiast zawieszać
    // przeglądarkę albo oddawać tysiąc pustych wierszy.
    if (dopasowanie[0] === '') continue;
    trafienia.push({
      poczatek,
      koniec: poczatek + dopasowanie[0].length,
      tekst: dopasowanie[0],
    });
  }
  return { trafienia, powod: '' };
}

/**
 * Zamienia wszystkie trafienia i oddaje treść nową wraz z ich liczbą.
 *
 * Treść wejściowa zostaje nietknięta — wynik jest nowym napisem, więc podgląd
 * przed zatwierdzeniem ma co pokazać, a Operator ma do czego wrócić.
 */
export function zamienWszystkie(
  tresc: string,
  nastawy: NastawyWyszukiwania,
  zamiennik: string,
): WynikZamiany {
  const znalezione = znajdzTrafienia(tresc, nastawy);
  if (znalezione.powod !== '') return { tresc, liczba: 0, powod: znalezione.powod };
  if (znalezione.trafienia.length === 0) return { tresc, liczba: 0, powod: '' };

  // Składamy od tyłu, bo podmiana od przodu przesuwałaby położenia trafień
  // jeszcze nieprzetworzonych o różnicę długości zamiennika.
  let wynik = tresc;
  for (const trafienie of [...znalezione.trafienia].reverse()) {
    wynik = `${wynik.slice(0, trafienie.poczatek)}${zamiennik}${wynik.slice(trafienie.koniec)}`;
  }
  return { tresc: wynik, liczba: znalezione.trafienia.length, powod: '' };
}

/** Fragment treści wokół trafienia — materiał podglądu przed zatwierdzeniem. */
export function otoczenieTrafienia(
  tresc: string,
  trafienie: TrafienieTekstu,
  promien = 40,
): string {
  const od = Math.max(0, trafienie.poczatek - promien);
  const az = Math.min(tresc.length, trafienie.koniec + promien);
  const przed = od > 0 ? '…' : '';
  const po = az < tresc.length ? '…' : '';
  return `${przed}${tresc.slice(od, az).replace(/\n/gu, '⏎')}${po}`;
}
