/**
 * Wyszukiwanie i zamiana w treści dokumentu działa na buforze edytora, bez rdzenia; wzorzec
 * niepoprawny nie jest wyciszany, wynik niesie powód wprost.
 */

/** Jedno trafienie wzorca w treści dokumentu wraz z jego dokładnym położeniem i dopasowanym fragmentem tekstu. */
export interface TrafienieTekstu {
  /** Położenie początku trafienia w treści. */
  poczatek: number;
  /** Położenie końca trafienia w treści. */
  koniec: number;
  /** Dopasowany fragment. */
  tekst: string;
}

/** Nastawy wyszukiwania odpowiadające przełącznikom panelu: wzorzec, tryb wyrażenia regularnego i wielkość liter. */
export interface NastawyWyszukiwania {
  wzorzec: string;
  /** Czy wzorzec jest wyrażeniem regularnym; `false` znaczy frazę dosłowną. */
  regularne: boolean;
  /** Czy wielkość liter ma znaczenie. */
  wielkoscLiter: boolean;
}

/** Wynik wyszukiwania niosący listę wszystkich trafień albo powód, dla którego wyszukiwania nie dało się wykonać. */
export interface WynikWyszukiwania {
  trafienia: readonly TrafienieTekstu[];
  /** Powód niepowodzenia; pusty, gdy wyszukiwanie się wykonało. */
  powod: string;
}

/** Wynik zamiany niosący treść dokumentu po podmianie wraz z liczbą podmienionych trafień wzorca wyszukiwania. */
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

/** Osłania znaki o znaczeniu składniowym wyrażenia regularnego, żeby fraza dosłowna znaczyła samą siebie. */
function oslon(fraza: string): string {
  return fraza.replace(/[.*+?^${}()|[\]\\]/gu, '\\$&');
}

/** Znajduje wszystkie trafienia wzorca w treści dokumentu zgodnie z podanymi nastawami wyszukiwania panelu. */
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
    // Wzorzec dopasowujący pustkę przesunąłby pętlę o zero znaków — pomijamy takie trafienia.
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

  // Składamy od tyłu, bo podmiana od przodu przesunęłaby położenia trafień jeszcze nieprzetworzonych.
  let wynik = tresc;
  for (const trafienie of [...znalezione.trafienia].reverse()) {
    wynik = `${wynik.slice(0, trafienie.poczatek)}${zamiennik}${wynik.slice(trafienie.koniec)}`;
  }
  return { tresc: wynik, liczba: znalezione.trafienia.length, powod: '' };
}

/** Fragment treści dokumentu wokół trafienia — materiał podglądu wyświetlany przed zatwierdzeniem zamiany. */
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
