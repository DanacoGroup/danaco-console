import { naPunkty, zPunktow, type JednostkaMiary } from './nastawy-strony';

/**
 * Podziałka linijek — rachunek kresek, przyciągania i przeliczeń.
 *
 * ── Dlaczego rachunek stoi osobno od linijki ─────────────────────────────────
 * Chwyt marginesu i zakładanie tabulatora liczą się w milimetrach, a rysują
 * w punktach ekranu, przy skali widoku, która bywa inna niż 100 %. Pomyłka o krok
 * podziałki przesuwa margines pisma urzędowego — więc rachunek ma dać się
 * sprawdzić bez przeglądarki, a linijka ma go wyłącznie rysować.
 *
 * ── Jednostka jest wyborem Operatora ────────────────────────────────────────
 * Milimetry albo cale, wedle zasady ogólnej zlecenia. Podziałka calowa nie jest
 * podziałką milimetrową z inną etykietą: kreski stoją co 1/8 cala, a przyciąganie
 * chodzi po 1/16 cala, bo tak działa linijka w pakiecie biurowym i tak Operator
 * spodziewa się, że tabulator wskoczy.
 */

/** Rodzaj kreski podziałki — trzy długości, wzorem linijki. */
export type RodzajKreski = 'duza' | 'srednia' | 'mala';

/** Jedna kreska podziałki. */
export interface KreskaPodzialki {
  /** Odległość od początku pola w milimetrach. */
  milimetry: number;
  rodzaj: RodzajKreski;
  /** Napis przy kresce; puste znaczy „kreska bez napisu". */
  napis: string;
}

/** Najmniejszy odstęp kresek, przy którym jeszcze da się je rozróżnić. */
const ODSTEP_ROZROZNIALNY = 3;

/**
 * Składa kreski podziałki dla pola o podanej długości.
 *
 * Kreska, która przy tej skali stanęłaby bliżej niż {@link ODSTEP_ROZROZNIALNY}
 * punktów od poprzedniej, jest pomijana — podziałka zlana w szarą wstęgę nie
 * mówi nic, a udaje, że mówi. Napisy niosą pełne centymetry albo cale, bo to one
 * są miarą, w której Operator myśli o marginesie.
 *
 * @param dlugoscMm Długość podziałki w milimetrach.
 * @param jednostka Jednostka wybrana przez Operatora.
 * @param skala Skala widoku jako krotność (1 znaczy 100 %).
 */
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
    // Kreski małe schodzą pierwsze, potem średnie: przy skali 30 % zostaje sama
    // siatka pełnych centymetrów, która jeszcze jest czytelna.
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

/** Krok przyciągania w milimetrach — pół milimetra albo 1/16 cala. */
export function krokPrzyciaganiaLinijki(jednostka: JednostkaMiary): number {
  return jednostka === 'cal' ? 25.4 / 16 : 0.5;
}

/**
 * Przyciąga długość do kroku podziałki i przycina do zakresu.
 *
 * Chwyt puszczony między kreskami ma wskoczyć na kreskę, bo margines 19,73 mm
 * nie jest nastawą, którą ktoś świadomie wybiera. Zakres jest przycinany, a nie
 * odrzucany: chwyt wleczony za krawędź kartki zostaje na krawędzi.
 */
export function przyciagnijNaLinijce(
  milimetry: number,
  jednostka: JednostkaMiary,
  dolnaMm: number,
  gornaMm: number,
): number {
  const krok = krokPrzyciaganiaLinijki(jednostka);
  const przyciete = Math.min(Math.max(milimetry, dolnaMm), gornaMm);
  const przyciagniete = Math.round(przyciete / krok) * krok;
  // Zaokrąglenie na dwie cyfry po przecinku: krok calowy jest liczbą niewymierną
  // w milimetrach i bez tego nastawa niosłaby ogon przypadkowych cyfr.
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

/** Punkty ekranu z milimetrów, z uwzględnieniem skali widoku. */
export function punktyNaLinijce(milimetry: number, skala: number): number {
  const krotnosc = skala > 0 && Number.isFinite(skala) ? skala : 1;
  return naPunkty(milimetry) * krotnosc;
}

/** Rodzaj tabulatora — cztery, wzorem pakietu biurowego. */
export type RodzajTabulatora = 'lewy' | 'prawy' | 'srodkowy' | 'dziesietny';

/** Znak wiodący tabulatora — czym wypełnić drogę do niego. */
export type ZnakWiodacy = 'brak' | 'kropka' | 'kreska' | 'podkreslenie';

/** Tabulator założony na linijce poziomej. */
export interface TabulatorAkapitu {
  /** Odległość od lewej krawędzi pola pisania w milimetrach. */
  milimetry: number;
  rodzaj: RodzajTabulatora;
  znakWiodacy: ZnakWiodacy;
}

/** Kolejność, w której naciśnięcie znaku tabulatora przestawia jego rodzaj. */
const KOLEJNOSC_TABULATOROW: readonly RodzajTabulatora[] = [
  'lewy',
  'srodkowy',
  'prawy',
  'dziesietny',
];

/** Rodzaj następny w obiegu — jedno naciśnięcie znaku tabulatora. */
export function nastepnyRodzajTabulatora(rodzaj: RodzajTabulatora): RodzajTabulatora {
  const numer = KOLEJNOSC_TABULATOROW.indexOf(rodzaj);
  return KOLEJNOSC_TABULATOROW[(numer + 1) % KOLEJNOSC_TABULATOROW.length] ?? 'lewy';
}

/** Znak rysowany na linijce dla rodzaju tabulatora. */
export function znakTabulatora(rodzaj: RodzajTabulatora): string {
  if (rodzaj === 'prawy') return '⌐';
  if (rodzaj === 'srodkowy') return '⊥';
  if (rodzaj === 'dziesietny') return '⊦';
  return 'L';
}

/** Nazwa rodzaju tabulatora dla objaśnienia i dla odczytu ekranowego. */
export function nazwaRodzajuTabulatora(rodzaj: RodzajTabulatora): string {
  if (rodzaj === 'prawy') return 'tabulator prawy';
  if (rodzaj === 'srodkowy') return 'tabulator środkowy';
  if (rodzaj === 'dziesietny') return 'tabulator dziesiętny';
  return 'tabulator lewy';
}

/** Nazwa znaku wiodącego. */
export function nazwaZnakuWiodacego(znak: ZnakWiodacy): string {
  if (znak === 'kropka') return 'znak wiodący: kropki';
  if (znak === 'kreska') return 'znak wiodący: kreski';
  if (znak === 'podkreslenie') return 'znak wiodący: podkreślenie';
  return 'bez znaku wiodącego';
}

/**
 * Wcięcia akapitu chwytane na linijce poziomej.
 *
 * Trzy osobne znaczniki, jak w pakiecie biurowym: wcięcie pierwszego wiersza
 * (górny trójkąt), wcięcie lewe wiersza dalszego (dolny trójkąt) i wcięcie prawe.
 * Wcięcie pierwszego wiersza jest liczone WZGLĘDEM wcięcia lewego, nie od
 * krawędzi pola — dzięki temu przeciągnięcie wcięcia lewego zabiera pierwszy
 * wiersz ze sobą, a to jest zachowanie, którego Operator się spodziewa.
 */
export interface WciecieAkapitu {
  /** Wcięcie lewe wierszy dalszych, od lewej krawędzi pola, w milimetrach. */
  leweMm: number;
  /** Wcięcie prawe, od prawej krawędzi pola, w milimetrach. */
  praweMm: number;
  /** Wcięcie pierwszego wiersza względem wcięcia lewego; ujemne znaczy wysunięcie. */
  pierwszyWierszMm: number;
}

/** Wcięcia zerowe — akapit wyrównany do marginesów. */
export function zerowaWciecieAkapitu(): WciecieAkapitu {
  return { leweMm: 0, praweMm: 0, pierwszyWierszMm: 0 };
}

/**
 * Zdanie o wcięciach akapitu — do objaśnienia chwytu i do paska stanu.
 */
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
