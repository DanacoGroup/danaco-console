import { znajdzWzorce } from './wzorce-placeholderow';

/** Pseudolokalizacja przekształca tekst źródłowy do testu interfejsu przed właściwym przekładem. */

/** Wydłużenie domyślne pola dopełnienia jest punktem wyjścia, a nie granicą — wywołanie może podać wartość inną. */
export const WYDLUZENIE_DOMYSLNE = 30;

/** Ramka odcinająca początek i koniec treści wyniku sprawia, że sklejenie kawałków tekstu widać od razu. */
const RAMKA_POCZATEK = '⟦';
const RAMKA_KONIEC = '⟧';

/** Znak dopełnienia treści wyniku nie jest literą alfabetu, więc wizualnie nie udaje rzeczywistej treści źródła. */
const ZNAK_DOPELNIENIA = '·';

/**
 * Zamiana liter na warianty z diakrytykami.
 *
 * Wykaz obejmuje litery alfabetu łacińskiego, dla których istnieje wariant
 * o zachowanej sylwetce; litera bez takiego wariantu zostaje sobą, bo celem
 * jest test kroju i szerokości, a nie nieczytelność.
 */
const DIAKRYTYKI = new Map<string, string>([
  ['a', 'á'], ['b', 'ḃ'], ['c', 'ç'], ['d', 'ð'], ['e', 'é'], ['f', 'ƒ'], ['g', 'ĝ'],
  ['h', 'ĥ'], ['i', 'í'], ['j', 'ĵ'], ['k', 'ķ'], ['l', 'ļ'], ['m', 'ṁ'], ['n', 'ñ'],
  ['o', 'ó'], ['p', 'ṗ'], ['r', 'ŕ'], ['s', 'š'], ['t', 'ţ'],
  ['u', 'ú'], ['v', 'ṽ'], ['w', 'ŵ'], ['x', 'ẋ'], ['y', 'ý'], ['z', 'ž'],
  ['A', 'Á'], ['B', 'Ḃ'], ['C', 'Ç'], ['D', 'Ð'], ['E', 'É'], ['F', 'Ḟ'], ['G', 'Ĝ'],
  ['H', 'Ĥ'], ['I', 'Í'], ['J', 'Ĵ'], ['K', 'Ķ'], ['L', 'Ļ'], ['M', 'Ṁ'], ['N', 'Ñ'],
  ['O', 'Ó'], ['P', 'Ṗ'], ['R', 'Ŕ'], ['S', 'Š'], ['T', 'Ţ'], ['U', 'Ú'], ['V', 'Ṽ'],
  ['W', 'Ŵ'], ['X', 'Ẋ'], ['Y', 'Ý'], ['Z', 'Ž'],
]);

/** Wynik przekształcenia niesie tekst po pseudolokalizacji wraz z liczbami, które pozwalają ocenić skalę zmiany. */
export interface WynikPseudolokalizacji {
  /** Tekst po przekształceniu. */
  readonly tekst: string;
  /** Długość tekstu wejściowego w znakach. */
  readonly dlugoscZrodla: number;
  /** Długość wyniku w znakach. */
  readonly dlugoscWyniku: number;
  /** Ile symboli zastępczych i znaczników przeszło nietkniętych. */
  readonly chronione: number;
}

/**
 * Przekształca tekst źródłowy.
 *
 * `wydluzenie` jest wyrażone w procentach długości źródła; wartość ujemna albo
 * nieliczbowa schodzi do zera, bo skracania pseudolokalizacja nie testuje.
 */
export function pseudolokalizuj(tekst: string, wydluzenie: number): WynikPseudolokalizacji {
  const chronione = znajdzWzorce(tekst);
  const procent = Number.isFinite(wydluzenie) && wydluzenie > 0 ? wydluzenie : 0;

  let przepisany = '';
  let ogon = 0;
  for (const wpis of chronione) {
    przepisany += zdiakrytyzuj(tekst.slice(ogon, wpis.poczatek)) + wpis.zapis;
    ogon = wpis.koniec;
  }
  przepisany += zdiakrytyzuj(tekst.slice(ogon));

  const dopelnienie = Math.round((tekst.length * procent) / 100);
  const tresc =
    dopelnienie > 0
      ? `${przepisany} ${ZNAK_DOPELNIENIA.repeat(dopelnienie)}`
      : przepisany;
  const wynik = tekst === '' ? '' : `${RAMKA_POCZATEK}${tresc}${RAMKA_KONIEC}`;

  return {
    tekst: wynik,
    dlugoscZrodla: tekst.length,
    dlugoscWyniku: wynik.length,
    chronione: chronione.length,
  };
}

function zdiakrytyzuj(odcinek: string): string {
  let wynik = '';
  for (const znak of odcinek) wynik += DIAKRYTYKI.get(znak) ?? znak;
  return wynik;
}

/** Zdanie o wyniku pseudolokalizacji podaje liczby przekształcenia, a nie zapewnienia o jego jakości albo skutku. */
export function zdanieOPseudolokalizacji(wynik: WynikPseudolokalizacji): string {
  if (wynik.dlugoscZrodla === 0) {
    return (
      'Pseudolokalizacja: rdzeń nie ma jeszcze tekstu źródłowego tego okna, więc nie ma czego ' +
      'przekształcić. Przekształcenie bierze materiał zapisany, nie treść wpisaną do pola — ' +
      'zapisz źródło w Source Panel.'
    );
  }
  const wzrost = Math.round(
    ((wynik.dlugoscWyniku - wynik.dlugoscZrodla) / wynik.dlugoscZrodla) * 100,
  );
  return (
    `Pseudolokalizacja: ${String(wynik.dlugoscZrodla)} znaków źródła → ` +
    `${String(wynik.dlugoscWyniku)} znaków wyniku (wzrost o ${String(wzrost)} procent). ` +
    `Symboli zastępczych i znaczników pozostawionych bez zmiany: ${String(wynik.chronione)}. ` +
    'Wynik powstaje w oknie — kontrakt nie ma komendy pseudolokalizacji, więc w rdzeniu nic ' +
    'się przy tym nie zapisało.'
  );
}
