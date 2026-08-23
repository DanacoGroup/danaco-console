import { IdentityLayer, IdentityMode } from '../../../shared/contract';

/**
 * Warstwy nakładki i tryby jej podania — słownik porządku i nazw.
 *
 * Warstwa mówi, jak krytyczna jest treść: konstytucja stoi wyżej niż profil roli,
 * profil wyżej niż ekspertyza. Ta kolejność rządzi zarówno układem wykazu
 * kategorii, jak i porządkiem złożonego promptu. Wartości pochodzą wyłącznie
 * ze stałych kontraktu.
 *
 * Tryb rozstrzyga los ustawień fabrycznych: `ZASTAP` podmienia prompt fabryczny
 * w całości, `DOLACZ` dokłada treść do niego. Wybór decyduje o tym, którym
 * argumentem uruchomienia pojedzie nakładka, więc każdy przełącznik trybu w tej
 * sekcji podaje go wprost.
 */

/** Trzy warstwy w kolejności krytyczności. Konstytucja pierwsza. */
export const WARSTWY_OD_NAJWAZNIEJSZEJ: readonly IdentityLayer[] = [
  IdentityLayer.Constitution,
  IdentityLayer.Profile,
  IdentityLayer.Expertise,
];

/** Nazwa warstwy pokazywana Operatorowi. */
export const NAZWY_WARSTW: Readonly<Record<IdentityLayer, string>> = {
  [IdentityLayer.Constitution]: 'konstytucja',
  [IdentityLayer.Profile]: 'profil roli',
  [IdentityLayer.Expertise]: 'ekspertyza',
};

/** Dwa tryby podania nakładki; ZASTĄP jest wartością domyślną. */
export const TRYBY: readonly IdentityMode[] = [IdentityMode.ZASTAP, IdentityMode.DOLACZ];

/** Nazwa trybu pokazywana Operatorowi. */
export const NAZWY_TRYBOW: Readonly<Record<IdentityMode, string>> = {
  [IdentityMode.ZASTAP]: 'ZASTĄP ustawienia fabryczne',
  [IdentityMode.DOLACZ]: 'DOŁĄCZ do ustawień fabrycznych',
};

/**
 * Ostrzeżenie towarzyszące trybowi, pokazywane przy przełączniku i w podglądzie.
 * Zastąpienie zdejmuje prompt fabryczny w całości.
 */
export const OSTRZEZENIA_TRYBOW: Readonly<Record<IdentityMode, string>> = {
  [IdentityMode.ZASTAP]:
    'ZASTĄPIENIE ZDEJMUJE PROMPT FABRYCZNY W CAŁOŚCI. Model dostanie wyłącznie treść złożoną z kategorii tej sekcji — żadna instrukcja fabryczna producenta nie zostanie podana. Wszystko, co ma obowiązywać, musi stać w tych kategoriach.',
  [IdentityMode.DOLACZ]:
    'DOŁĄCZENIE ZOSTAWIA PROMPT FABRYCZNY. Treść kategorii zostanie dopisana do instrukcji fabrycznej producenta, a nie postawiona w jej miejsce.',
};

/** Pierwszeństwo warstwy; warstwa spoza kontraktu trafia na koniec. */
export function pierwszenstwoWarstwy(warstwa: string): number {
  const miejsce = (WARSTWY_OD_NAJWAZNIEJSZEJ as readonly string[]).indexOf(warstwa);
  return miejsce === -1 ? WARSTWY_OD_NAJWAZNIEJSZEJ.length : miejsce;
}

/** Nazwa warstwy gotowa do wydruku; warstwa spoza kontraktu pokazuje własny kod. */
export function nazwaWarstwy(warstwa: string): string {
  return NAZWY_WARSTW[warstwa as IdentityLayer] ?? warstwa;
}

/** Nazwa trybu gotowa do wydruku. */
export function nazwaTrybu(tryb: string): string {
  return NAZWY_TRYBOW[tryb as IdentityMode] ?? tryb;
}

/** Ostrzeżenie trybu gotowe do wydruku; tryb spoza kontraktu ostrzeżenia nie ma. */
export function ostrzezenieTrybu(tryb: string): string {
  return OSTRZEZENIA_TRYBOW[tryb as IdentityMode] ?? '';
}

/**
 * Tryb odczytany z wartości kontrolki. Wartość spoza kontraktu daje ZASTĄP —
 * wartość domyślną klucza `tozsamosc.tryb_domyslny`.
 */
export function trybZWyboru(wartosc: string): IdentityMode {
  return (TRYBY as readonly string[]).includes(wartosc)
    ? (wartosc as IdentityMode)
    : IdentityMode.ZASTAP;
}
