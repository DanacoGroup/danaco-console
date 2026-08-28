/**
 * Warstwy nakładki tożsamości i tryby jej podania: słownik porządku, nazw
 * i ostrzeżeń. Warstwa mówi, jak krytyczna jest treść, a tryb rozstrzyga los
 * promptu fabrycznego. Wszystkie wartości pochodzą wyłącznie ze stałych
 * kontraktu.
 */
import { IdentityLayer, IdentityMode } from '../../../shared/contract';

/**
 * Trzy warstwy w kolejności krytyczności; konstytucja pierwsza, ekspertyza
 * ostatnia. Kolejność wykazu rządzi układem kategorii w oknie i porządkiem
 * złożonego promptu.
 */
export const WARSTWY_OD_NAJWAZNIEJSZEJ: readonly IdentityLayer[] = [
  IdentityLayer.Constitution,
  IdentityLayer.Profile,
  IdentityLayer.Expertise,
];

/**
 * Nazwa warstwy pokazywana w oknie. Odwzorowanie pokrywa wszystkie trzy warstwy
 * kontraktu, więc odczyt po kluczu warstwy kontraktowej nie ma przypadku pustego.
 */
export const NAZWY_WARSTW: Readonly<Record<IdentityLayer, string>> = {
  [IdentityLayer.Constitution]: 'konstytucja',
  [IdentityLayer.Profile]: 'profil roli',
  [IdentityLayer.Expertise]: 'ekspertyza',
};

/**
 * Dwa tryby podania nakładki w kolejności wyboru; zastąpienie jest wartością
 * domyślną, do której wraca odczyt wartości nierozpoznanej.
 */
export const TRYBY: readonly IdentityMode[] = [IdentityMode.ZASTAP, IdentityMode.DOLACZ];

/**
 * Nazwa trybu pokazywana w oknie. Brzmienie niesie skutek wyboru wprost,
 * ponieważ przełącznik trybu rozstrzyga o losie promptu fabrycznego.
 */
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

/**
 * Pierwszeństwo warstwy; warstwa spoza kontraktu trafia na koniec, bo brak
 * w wykazie daje liczbę równą jego długości. Wynik służy porządkowaniu kategorii.
 */
export function pierwszenstwoWarstwy(warstwa: string): number {
  const miejsce = (WARSTWY_OD_NAJWAZNIEJSZEJ as readonly string[]).indexOf(warstwa);
  return miejsce === -1 ? WARSTWY_OD_NAJWAZNIEJSZEJ.length : miejsce;
}

/**
 * Nazwa warstwy gotowa do wydruku; warstwa spoza kontraktu pokazuje własny kod,
 * więc wykaz nie gubi wiersza o nieznanej warstwie.
 */
export function nazwaWarstwy(warstwa: string): string {
  return NAZWY_WARSTW[warstwa as IdentityLayer] ?? warstwa;
}

/**
 * Nazwa trybu gotowa do wydruku. Tryb spoza kontraktu pokazuje własny kod,
 * ponieważ odczyt odwzorowania zastępuje brak wartością wejściową.
 */
export function nazwaTrybu(tryb: string): string {
  return NAZWY_TRYBOW[tryb as IdentityMode] ?? tryb;
}

/**
 * Ostrzeżenie trybu gotowe do wydruku; tryb spoza kontraktu ostrzeżenia nie ma,
 * więc miejsce ostrzeżenia zostaje puste zamiast pokazywać kod trybu.
 */
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
