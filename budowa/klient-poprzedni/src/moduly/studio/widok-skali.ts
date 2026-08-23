import { naPunkty } from './nastawy-strony';
import { SKALA_DOLNA, SKALA_GORNA, type MagazynNastawWidoku } from './widok-nastawy-operatora';

/**
 * Skala widoku — procenty, nastawy gotowe i pamięć skali przy dokumencie.
 *
 * ── Dlaczego rachunek, a nie sam suwak ──────────────────────────────────────
 * „Do szerokości strony" nie jest wartością, jest RACHUNKIEM: zależy od
 * szerokości pola widoku, od nośnika, od orientacji i od liczby kartek w rzędzie.
 * Wpisana na sztywno rozjechałaby się przy pierwszej zmianie nośnika. Stoi więc
 * jako funkcja czysta, którą da się sprawdzić bez przeglądarki — wołający podaje
 * zmierzone pole widoku, a plik oddaje procenty.
 *
 * ── Skala jest pamiętana przy dokumencie ────────────────────────────────────
 * Wymóg zlecenia. Pamięć jest przy dokumencie, nie przy oknie: Operator wraca do
 * pisma po dwóch dniach i ma je zobaczyć w skali, w której nad nim pracował,
 * a drugie pismo w swojej.
 */

/** Nastawa gotowa skali. */
export interface NastawaSkali {
  kod: string;
  nazwa: string;
  opis: string;
}

/** Nastawy gotowe skali, wzorem pakietu biurowego. */
export const NASTAWY_SKALI: readonly NastawaSkali[] = [
  {
    kod: 'sto',
    nazwa: '100 %',
    opis: 'Kartka w rozmiarze, w jakim wyjdzie z drukarki — milimetr nośnika na milimetr ekranu.',
  },
  {
    kod: 'szerokosc-strony',
    nazwa: 'Do szerokości strony',
    opis: 'Cała szerokość kartki wraz z marginesami mieści się w polu widoku.',
  },
  {
    kod: 'cala-strona',
    nazwa: 'Cała strona',
    opis: 'Cała kartka widoczna naraz — do oceny łamania i rozkładu, nie do pisania.',
  },
  {
    kod: 'szerokosc-tekstu',
    nazwa: 'Do szerokości tekstu',
    opis:
      'Pole pisania wypełnia widok, marginesy zostają za krawędzią — najwięcej liter przy ' +
      'zachowanej postaci kartki.',
  },
];

/** Skale wpisane w wykaz — obok pola liczbowego i suwaka. */
export const SKALE_GOTOWE: readonly number[] = [50, 75, 100, 125, 150, 200];

/** Wymiary potrzebne do policzenia skali z nastawy gotowej. */
export interface WymiaryPolaWidoku {
  /** Szerokość pola widoku w punktach ekranu, bez pasków przewijania. */
  szerokoscWidokuPx: number;
  /** Wysokość pola widoku w punktach ekranu. */
  wysokoscWidokuPx: number;
  /** Szerokość kartki w milimetrach, z uwzględnieniem orientacji. */
  szerokoscKartkiMm: number;
  /** Wysokość kartki w milimetrach, z uwzględnieniem orientacji. */
  wysokoscKartkiMm: number;
  /** Szerokość pola pisania w milimetrach — kartka bez marginesów. */
  szerokoscTekstuMm: number;
  /** Liczba kartek w rzędzie; rozkładówka podaje dwie. */
  kartekWRzedzie: number;
  /** Odstęp między kartkami w rzędzie, w punktach ekranu. */
  odstepPx: number;
}

/**
 * Skala w procentach dla nastawy gotowej.
 *
 * Nastawa nieznana oddaje `null`, a nie 100 %: cicha podmiana nastawy na inną
 * byłaby kłamstwem o tym, co Operator wybrał. Pole widoku o zerowej szerokości
 * (okno jeszcze nieosadzone) też oddaje `null` — skala nieskończona nie jest
 * skalą.
 */
export function policzSkaleWidoku(
  kod: string,
  wymiary: WymiaryPolaWidoku,
): number | null {
  const wRzedzie = Math.max(1, Math.round(wymiary.kartekWRzedzie));
  const odstepy = wymiary.odstepPx * (wRzedzie - 1);
  const doDyspozycji = wymiary.szerokoscWidokuPx - odstepy;
  if (kod === 'sto') return 100;
  if (doDyspozycji <= 0) return null;

  if (kod === 'szerokosc-strony') {
    const potrzebne = naPunkty(wymiary.szerokoscKartkiMm) * wRzedzie;
    if (potrzebne <= 0) return null;
    return przytnijSkale((doDyspozycji / potrzebne) * 100);
  }
  if (kod === 'szerokosc-tekstu') {
    const potrzebne = naPunkty(wymiary.szerokoscTekstuMm) * wRzedzie;
    if (potrzebne <= 0) return null;
    return przytnijSkale((doDyspozycji / potrzebne) * 100);
  }
  if (kod === 'cala-strona') {
    const potrzebneSzerz = naPunkty(wymiary.szerokoscKartkiMm) * wRzedzie;
    const potrzebneWzwyz = naPunkty(wymiary.wysokoscKartkiMm);
    if (potrzebneSzerz <= 0 || potrzebneWzwyz <= 0 || wymiary.wysokoscWidokuPx <= 0) return null;
    // Cała strona znaczy: mieści się i wszerz, i wzwyż — więc rozstrzyga wymiar
    // ciaśniejszy. Wzięcie samej szerokości ucinałoby dół kartki i nazywałoby to
    // „całą stroną".
    const wszerz = doDyspozycji / potrzebneSzerz;
    const wzwyz = wymiary.wysokoscWidokuPx / potrzebneWzwyz;
    return przytnijSkale(Math.min(wszerz, wzwyz) * 100);
  }
  return null;
}

/** Przycina skalę do granic i zaokrągla do pełnego procentu. */
export function przytnijSkale(procent: number): number {
  if (!Number.isFinite(procent)) return 100;
  return Math.min(Math.max(Math.round(procent), SKALA_DOLNA), SKALA_GORNA);
}

/** Klucz zapisu skali dokumentu. */
const PRZEDROSTEK_SKALI = 'dn.studio.skala.';

function magazynDomyslny(): MagazynNastawWidoku | null {
  try {
    return globalThis.localStorage ?? null;
  } catch {
    return null;
  }
}

/**
 * Skala zapamiętana przy dokumencie; brak zapisu oddaje `null`.
 *
 * `null`, a nie 100 %: wołający ma wtedy wziąć skalę ostatnią Operatora, a nie
 * wracać do stu procent przy każdym nowym dokumencie.
 */
export function skalaDokumentu(
  idDokumentu: string,
  magazyn: MagazynNastawWidoku | null = magazynDomyslny(),
): number | null {
  if (idDokumentu === '' || magazyn === null) return null;
  try {
    const zapis = magazyn.getItem(PRZEDROSTEK_SKALI + idDokumentu);
    if (zapis === null || zapis === '') return null;
    const wartosc = Number(zapis);
    return Number.isFinite(wartosc) ? przytnijSkale(wartosc) : null;
  } catch {
    return null;
  }
}

/** Zapamiętuje skalę przy dokumencie. */
export function zapamietajSkaleDokumentu(
  idDokumentu: string,
  procent: number,
  magazyn: MagazynNastawWidoku | null = magazynDomyslny(),
): void {
  if (idDokumentu === '' || magazyn === null) return;
  try {
    magazyn.setItem(PRZEDROSTEK_SKALI + idDokumentu, String(przytnijSkale(procent)));
  } catch {
    // Nastawa widoku nie jest powodem, żeby cokolwiek przerywać.
  }
}
