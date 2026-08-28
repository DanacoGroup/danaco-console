import { ExtensionKind, ExtensionOrigin, type Extension } from '../../../../shared/contract';

/**
 * Zawęża katalog rozszerzeń po stronie klienta fasetami rodzaju, pochodzenia
 * i stanu oraz tekstem szukanym w nazwie, kodzie i opisie pozycji. Filtr jest
 * czynnością czystą — bez elementu i bez kanału — więc daje się sprawdzić bez
 * budowania dokumentu.
 */
export interface ZawezenieWidoku {
  /** Rodzaj pozycji; pusty łańcuch znaczy „wszystkie rodzaje". */
  rodzaj: string;
  /** Pochodzenie pozycji; pusty łańcuch znaczy „oba źródła". */
  pochodzenie: string;
  /** Stan pozycji: `''` wszystkie, `dostepne`, `zainstalowane`, `wlaczone`. */
  stan: string;
  /** Tekst szukany w nazwie, kodzie i opisie; pusty nie zawęża. */
  tekst: string;
}

/**
 * Zawężenie puste, czyli punkt wyjścia okna, w którym widać cały katalog. Każda
 * faseta stoi tu na wartości pustej, a wartość pusta w każdej z nich znaczy „bez
 * zawężenia tą fasetą".
 */
export const BEZ_ZAWEZENIA: ZawezenieWidoku = {
  rodzaj: '',
  pochodzenie: '',
  stan: '',
  tekst: '',
};

/**
 * Wartości fasety stanu wraz z etykietami: wartość jedzie do zawężenia, etykieta
 * na ekran. Stan pozycji składa się z pól zainstalowania i włączenia, więc fasety
 * nie da się wziąć wprost z wyliczenia kontraktu.
 */
export const STANY_WIDOKU: ReadonlyArray<{ wartosc: string; etykieta: string }> = [
  { wartosc: '', etykieta: 'każdy stan' },
  { wartosc: 'dostepne', etykieta: 'dostępne, niezainstalowane' },
  { wartosc: 'zainstalowane', etykieta: 'zainstalowane' },
  { wartosc: 'wlaczone', etykieta: 'włączone' },
];

/**
 * Wartości fasety rodzaju wzięte wprost z wyliczenia `ExtensionKind` kontraktu,
 * uzupełnione o wartość pustą oznaczającą każdy rodzaj. Etykiety są zdaniem
 * Operatora i z nazwami kontraktu się nie pokrywają.
 */
export const RODZAJE_WIDOKU: ReadonlyArray<{ wartosc: string; etykieta: string }> = [
  { wartosc: '', etykieta: 'każdy rodzaj' },
  { wartosc: ExtensionKind.Mcp, etykieta: 'serwer MCP' },
  { wartosc: ExtensionKind.Plugin, etykieta: 'wtyczka' },
  { wartosc: ExtensionKind.Api, etykieta: 'integracja API' },
  { wartosc: ExtensionKind.Skill, etykieta: 'umiejętność' },
];

/**
 * Wartości fasety pochodzenia wzięte wprost z wyliczenia `ExtensionOrigin`
 * kontraktu, uzupełnione o wartość pustą oznaczającą oba źródła katalogu
 * rozszerzeń.
 */
export const POCHODZENIA_WIDOKU: ReadonlyArray<{ wartosc: string; etykieta: string }> = [
  { wartosc: '', etykieta: 'oba źródła' },
  { wartosc: ExtensionOrigin.Danaco, etykieta: 'Danaco Plugin' },
  { wartosc: ExtensionOrigin.Personal, etykieta: 'Personal' },
];

/**
 * Rozstrzyga, czy pozycja przechodzi fasetę stanu. Stan składa się z pól
 * `installed` oraz `enabled` pozycji katalogu, a wartość pusta fasety przepuszcza
 * każdą pozycję.
 */
function zgodnyStan(pozycja: Extension, stan: string): boolean {
  if (stan === '') return true;
  if (stan === 'dostepne') return !pozycja.installed;
  if (stan === 'zainstalowane') return pozycja.installed;
  return pozycja.enabled;
}

/**
 * Rozstrzyga, czy tekst szukany występuje w nazwie, kodzie albo opisie pozycji.
 * Porównanie idzie bez rozróżnienia wielkości liter, ale bez normalizacji
 * znaków diakrytycznych.
 */
function zgodnyTekst(pozycja: Extension, tekst: string): boolean {
  if (tekst === '') return true;
  const szukane = tekst.toLowerCase();
  return [pozycja.name, pozycja.code, pozycja.description ?? ''].some((pole) =>
    pole.toLowerCase().includes(szukane),
  );
}

/**
 * Oddaje pozycje katalogu spełniające wszystkie cztery fasety zawężenia,
 * w kolejności oddanej przez rdzeń. Kolejność zostaje nietknięta, bo porządek
 * wykazu należy do rdzenia, a nie do okna.
 */
export function zawez(
  pozycje: readonly Extension[],
  zawezenie: ZawezenieWidoku,
): readonly Extension[] {
  return pozycje.filter(
    (pozycja) =>
      (zawezenie.rodzaj === '' || pozycja.kind === zawezenie.rodzaj) &&
      (zawezenie.pochodzenie === '' || pozycja.origin === zawezenie.pochodzenie) &&
      zgodnyStan(pozycja, zawezenie.stan) &&
      zgodnyTekst(pozycja, zawezenie.tekst),
  );
}

/**
 * Składa zdanie o wyniku zawężenia: ile pozycji pokazano z ilu, a przy czynnym
 * szukaniu również jaki jest jego zakres. Katalog niezawężony dostaje samą liczbę
 * pozycji, a katalog pusty — zdanie puste.
 */
export function zdanieZawezenia(
  pokazane: number,
  wszystkie: number,
  zawezenie: ZawezenieWidoku,
): string {
  if (wszystkie === 0) return '';
  const czynne = zawezenie !== BEZ_ZAWEZENIA && czyZawezone(zawezenie);
  if (!czynne) return `Katalog liczy ${wszystkie} pozycji.`;
  const oTekscie =
    zawezenie.tekst === ''
      ? ''
      : ' Szukanie obejmuje nazwę, kod i opis pozycji — nie obejmuje udostępnianych ' +
        'narzędzi ani znaczników, bo pozycja katalogu ich nie niesie.';
  return `Pokazano ${pokazane} z ${wszystkie} pozycji katalogu.${oTekscie}`;
}

/**
 * Rozstrzyga, czy zawężenie cokolwiek zawęża, czyli czy którakolwiek z czterech
 * faset stoi na wartości innej niż pusta. Okno bierze stąd rozstrzygnięcie
 * o brzmieniu zdania podsumowującego wykaz.
 */
export function czyZawezone(zawezenie: ZawezenieWidoku): boolean {
  return (
    zawezenie.rodzaj !== '' ||
    zawezenie.pochodzenie !== '' ||
    zawezenie.stan !== '' ||
    zawezenie.tekst !== ''
  );
}
