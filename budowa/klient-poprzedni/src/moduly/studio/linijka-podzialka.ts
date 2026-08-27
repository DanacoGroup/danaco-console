import { naPunkty, zPunktow, type JednostkaMiary } from './nastawy-strony';

/** Podziałka linijek liczy kreski, przyciąganie i przeliczenia w milimetrach; rodzaj kreski ma trzy długości: dużą, średnią i małą. */
export type RodzajKreski = 'duza' | 'srednia' | 'mala';

/** Jedna kreska podziałki linijki wraz z jej odległością od początku pola, rodzajem i ewentualnym napisem. */
export interface KreskaPodzialki {
  /** Odległość od początku pola w milimetrach. */
  milimetry: number;
  rodzaj: RodzajKreski;
  /** Napis przy kresce; puste znaczy „kreska bez napisu". */
  napis: string;
}

/** Najmniejszy odstęp między kreskami podziałki, przy którym jeszcze da się je od siebie rozróżnić na ekranie. */
const ODSTEP_ROZROZNIALNY = 3;

/** Składa kreski podziałki dla pola o podanej długości, pomijając kreski, które przy bieżącej skali stanęłyby zbyt blisko siebie. */
export function zlozPodzialkeLinijki(
  dlugoscMm: number,
  jednostka: JednostkaMiary,
  skala: number,
): KreskaPodzialki[] {
  if (!Number.isFinite(dlugoscMm) || dlugoscMm <= 0) return [];
  const krotnosc = skala > 0 && Number.isFinite(skala) ? skala : 1;
  const krokMm = jednostka === 'cal' ? 25.4 / 8 : 1;
  const naSrednia = jednostka === 'cal' ? 4 : 5;
  const naDuza = jednostka === 'cal' ? 8 : 10;

  const kreski: KreskaPodzialki[] = [];
  const odstepPunktow = naPunkty(krokMm) * krotnosc;
  const liczba = Math.floor(dlugoscMm / krokMm);
  for (let numer = 0; numer <= liczba; numer += 1) {
    const duza = numer % naDuza === 0;
    const srednia = !duza && numer % naSrednia === 0;
    // Kreski małe schodzą pierwsze, potem średnie, zostawiając przy niskiej skali siatkę centymetrów.
    if (!duza) {
      if (odstepPunktow < ODSTEP_ROZROZNIALNY && !srednia) continue;
      if (odstepPunktow * naSrednia < ODSTEP_ROZROZNIALNY) continue;
    }
    const numerDuzej = numer / naDuza;
    kreski.push({
      milimetry: numer * krokMm,
      rodzaj: duza ? 'duza' : srednia ? 'srednia' : 'mala',
      napis: duza && numerDuzej > 0 ? String(numerDuzej) : '',
    });
  }
  return kreski;
}

/** Krok przyciągania na linijce w milimetrach: pół milimetra dla jednostki metrycznej albo szesnastą cala. */
export function krokPrzyciaganiaLinijki(jednostka: JednostkaMiary): number {
  return jednostka === 'cal' ? 25.4 / 16 : 0.5;
}

/** Przyciąga wskazaną długość do kroku podziałki linijki i przycina wynik do dozwolonego zakresu wartości. */
export function przyciagnijNaLinijce(
  milimetry: number,
  jednostka: JednostkaMiary,
  dolnaMm: number,
  gornaMm: number,
): number {
  const krok = krokPrzyciaganiaLinijki(jednostka);
  const przyciete = Math.min(Math.max(milimetry, dolnaMm), gornaMm);
  const przyciagniete = Math.round(przyciete / krok) * krok;
  // Zaokrąglenie na dwie cyfry po przecinku usuwa ogon przypadkowych cyfr niesiony przez krok calowy.
  return Math.round(Math.min(Math.max(przyciagniete, dolnaMm), gornaMm) * 100) / 100;
}

/**
 * Milimetry z odległości w punktach ekranu, z uwzględnieniem skali widoku.
 *
 * Skala niedodatnia albo nieskończona bierze 1, zamiast dzielić przez zero:
 * chwyt linijki nie jest miejscem, w którym wolno oddać nieskończoność.
 */
export function milimetryZPunktowLinijki(punkty: number, skala: number): number {
  const krotnosc = skala > 0 && Number.isFinite(skala) ? skala : 1;
  return zPunktow(punkty / krotnosc);
}

/** Punkty ekranu odpowiadające podanej odległości w milimetrach, z uwzględnieniem bieżącej skali widoku dokumentu. */
export function punktyNaLinijce(milimetry: number, skala: number): number {
  const krotnosc = skala > 0 && Number.isFinite(skala) ? skala : 1;
  return naPunkty(milimetry) * krotnosc;
}

/** Rodzaj tabulatora akapitu, cztery wartości wzorem pakietu biurowego: lewy, prawy, środkowy i dziesiętny. */
export type RodzajTabulatora = 'lewy' | 'prawy' | 'srodkowy' | 'dziesietny';

/** Znak wiodący tabulatora akapitu, czyli czym wypełnia się drogę od tekstu do pozycji tabulatora na wierszu. */
export type ZnakWiodacy = 'brak' | 'kropka' | 'kreska' | 'podkreslenie';

/** Tabulator założony na linijce poziomej wraz z jego odległością od lewej krawędzi, rodzajem i znakiem wiodącym. */
export interface TabulatorAkapitu {
  /** Odległość od lewej krawędzi pola pisania w milimetrach. */
  milimetry: number;
  rodzaj: RodzajTabulatora;
  znakWiodacy: ZnakWiodacy;
}

/** Kolejność rodzajów tabulatora, w której kolejne naciśnięcie znaku tabulatora przestawia jego rodzaj na linijce. */
const KOLEJNOSC_TABULATOROW: readonly RodzajTabulatora[] = [
  'lewy',
  'srodkowy',
  'prawy',
  'dziesietny',
];

/** Zwraca rodzaj tabulatora następny w obiegu po bieżącym, wywoływany jednym naciśnięciem znaku tabulatora. */
export function nastepnyRodzajTabulatora(rodzaj: RodzajTabulatora): RodzajTabulatora {
  const numer = KOLEJNOSC_TABULATOROW.indexOf(rodzaj);
  return KOLEJNOSC_TABULATOROW[(numer + 1) % KOLEJNOSC_TABULATOROW.length] ?? 'lewy';
}

/** Znak rysowany na linijce poziomej odpowiadający rodzajowi tabulatora ustawionego w danym miejscu wiersza. */
export function znakTabulatora(rodzaj: RodzajTabulatora): string {
  if (rodzaj === 'prawy') return '⌐';
  if (rodzaj === 'srodkowy') return '⊥';
  if (rodzaj === 'dziesietny') return '⊦';
  return 'L';
}

/** Nazwa rodzaju tabulatora przeznaczona do objaśnienia w interfejsie oraz do odczytu przez czytnik ekranowy. */
export function nazwaRodzajuTabulatora(rodzaj: RodzajTabulatora): string {
  if (rodzaj === 'prawy') return 'tabulator prawy';
  if (rodzaj === 'srodkowy') return 'tabulator środkowy';
  if (rodzaj === 'dziesietny') return 'tabulator dziesiętny';
  return 'tabulator lewy';
}

/** Nazwa znaku wiodącego tabulatora, czyli znaku wypełniającego drogę do jego pozycji na linijce poziomej. */
export function nazwaZnakuWiodacego(znak: ZnakWiodacy): string {
  if (znak === 'kropka') return 'znak wiodący: kropki';
  if (znak === 'kreska') return 'znak wiodący: kreski';
  if (znak === 'podkreslenie') return 'znak wiodący: podkreślenie';
  return 'bez znaku wiodącego';
}

/** Wcięcia akapitu chwytane na linijce poziomej trzema osobnymi znacznikami, jak w pakiecie biurowym: pierwszy wiersz, wiersze dalsze i wcięcie prawe. */
export interface WciecieAkapitu {
  /** Wcięcie lewe wierszy dalszych, od lewej krawędzi pola, w milimetrach. */
  leweMm: number;
  /** Wcięcie prawe, od prawej krawędzi pola, w milimetrach. */
  praweMm: number;
  /** Wcięcie pierwszego wiersza względem wcięcia lewego; ujemne znaczy wysunięcie. */
  pierwszyWierszMm: number;
}

/** Wcięcia zerowe akapitu, czyli akapit wyrównany do obu marginesów bez żadnego wcięcia ani wysunięcia. */
export function zerowaWciecieAkapitu(): WciecieAkapitu {
  return { leweMm: 0, praweMm: 0, pierwszyWierszMm: 0 };
}

/** Zdanie opisujące wcięcia akapitu do objaśnienia chwytu na linijce oraz do wyświetlenia w pasku stanu. */
export function opiszWciecieAkapitu(
  wciecie: WciecieAkapitu,
  jednostka: JednostkaMiary,
): string {
  const opis = (milimetry: number): string =>
    jednostka === 'cal'
      ? `${(milimetry / 25.4).toFixed(2)}″`
      : `${Math.round(milimetry * 10) / 10} mm`;
  const pierwszy =
    wciecie.pierwszyWierszMm === 0
      ? 'pierwszy wiersz bez wcięcia'
      : wciecie.pierwszyWierszMm > 0
        ? `pierwszy wiersz +${opis(wciecie.pierwszyWierszMm)}`
        : `pierwszy wiersz wysunięty o ${opis(-wciecie.pierwszyWierszMm)}`;
  return `wcięcie lewe ${opis(wciecie.leweMm)} · prawe ${opis(wciecie.praweMm)} · ${pierwszy}`;
}
