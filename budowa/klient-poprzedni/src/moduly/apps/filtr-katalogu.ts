import { ExtensionKind, ExtensionOrigin, type Extension } from '../../../../shared/contract';

/**
 * Zawężenie katalogu rozszerzeń po stronie klienta — fasety i szukanie
 * App Catalogu.
 *
 * Filtrowanie idzie tutaj, a nie w żądaniu, bo `extension.list` zawęża
 * wyłącznie rodzajem i stanem zainstalowania. Pozostałych faset opracowania —
 * źródła, stanu włączenia, tekstu — kontrakt w żądaniu nie ma, a cztery odczyty
 * po jednym na fasetę dałyby cztery migawki z czterech różnych chwil.
 *
 * Zakres szukania jest ograniczony i okno musi to powiedzieć: przeszukujemy
 * nazwę, kod i opis, bo tyle niesie pozycja katalogu. Szukanie po udostępnianych
 * narzędziach i znacznikach stoi w wykazie braków, zamiast udawać, że pusty
 * wynik znaczy „nie ma takiego rozszerzenia".
 *
 * Filtr jest czynnością czystą — bez elementu i bez kanału — więc daje się
 * sprawdzić bez budowania dokumentu.
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

/** Zawężenie puste — punkt wyjścia okna, w którym widać cały katalog. */
export const BEZ_ZAWEZENIA: ZawezenieWidoku = {
  rodzaj: '',
  pochodzenie: '',
  stan: '',
  tekst: '',
};

/** Wartości fasety stanu; klucz jedzie do zawężenia, napis na ekran. */
export const STANY_WIDOKU: ReadonlyArray<{ wartosc: string; etykieta: string }> = [
  { wartosc: '', etykieta: 'każdy stan' },
  { wartosc: 'dostepne', etykieta: 'dostępne, niezainstalowane' },
  { wartosc: 'zainstalowane', etykieta: 'zainstalowane' },
  { wartosc: 'wlaczone', etykieta: 'włączone' },
];

/** Wartości fasety rodzaju, wprost z wyliczenia kontraktu. */
export const RODZAJE_WIDOKU: ReadonlyArray<{ wartosc: string; etykieta: string }> = [
  { wartosc: '', etykieta: 'każdy rodzaj' },
  { wartosc: ExtensionKind.Mcp, etykieta: 'serwer MCP' },
  { wartosc: ExtensionKind.Plugin, etykieta: 'wtyczka' },
  { wartosc: ExtensionKind.Api, etykieta: 'integracja API' },
  { wartosc: ExtensionKind.Skill, etykieta: 'umiejętność' },
];

/** Wartości fasety pochodzenia, wprost z wyliczenia kontraktu. */
export const POCHODZENIA_WIDOKU: ReadonlyArray<{ wartosc: string; etykieta: string }> = [
  { wartosc: '', etykieta: 'oba źródła' },
  { wartosc: ExtensionOrigin.Danaco, etykieta: 'Danaco Plugin' },
  { wartosc: ExtensionOrigin.Personal, etykieta: 'Personal' },
];

/** Czy pozycja przechodzi fasetę stanu. */
function zgodnyStan(pozycja: Extension, stan: string): boolean {
  if (stan === '') return true;
  if (stan === 'dostepne') return !pozycja.installed;
  if (stan === 'zainstalowane') return pozycja.installed;
  return pozycja.enabled;
}

/**
 * Czy tekst szukany występuje w polach, które pozycja niesie.
 *
 * Porównanie idzie bez rozróżnienia wielkości liter, ale bez normalizacji
 * diakrytyków: „Umiejętność" i „Umiejetnosc" zostają dwoma różnymi napisami,
 * bo zrównanie ich w kliencie kazałoby oknu twierdzić coś o dopasowaniu, czego
 * rdzeń przy własnym szukaniu nie potwierdzi.
 */
function zgodnyTekst(pozycja: Extension, tekst: string): boolean {
  if (tekst === '') return true;
  const szukane = tekst.toLowerCase();
  return [pozycja.name, pozycja.code, pozycja.description ?? ''].some((pole) =>
    pole.toLowerCase().includes(szukane),
  );
}

/** Pozycje katalogu spełniające zawężenie, w kolejności oddanej przez rdzeń. */
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
 * Zdanie o wyniku zawężenia — ile pokazano z ilu i czym zawężono.
 *
 * Bez tego zdania pusty wykaz przy czynnym filtrze czyta się jak pusty rejestr.
 * Zdanie wymienia zakres szukania, bo tekst nieznaleziony w opisie bywa nazwą
 * narzędzia, której pozycja katalogu nie niesie.
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

/** Czy zawężenie cokolwiek zawęża. */
export function czyZawezone(zawezenie: ZawezenieWidoku): boolean {
  return (
    zawezenie.rodzaj !== '' ||
    zawezenie.pochodzenie !== '' ||
    zawezenie.stan !== '' ||
    zawezenie.tekst !== ''
  );
}
