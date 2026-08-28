// Import wprost do pliku panelu, nie do `moduly/<modul>/indeks` — tamten
// wciągnąłby gospodarzowi całe złożenie modułu wraz z jego oknami operacyjnymi
// i rejestracją w powłoce. Wytwórnia potrzebuje jednej funkcji, nie modułu.
import { KOD_PANELU_ZASOBOW, utworzPanelZasobow } from '../moduly/design/panel-zasobow';
import { utworzPanelDebaty } from '../moduly/roundtable/panel-debaty';
import { utworzOknoTerminalaPomocnicze } from '../moduly/terminal/okno-pomocnicze';
import { utworzOknoHistoriiRozmowy } from './okno-historii-rozmowy';
import { utworzOknoPodgladuBash } from './okno-podglad-bash';
import type { OpcjePanelu, PanelPomocniczy } from './panel-pomocniczy';

/**
 * Mapowanie kodu pozycji na wytwórnię panelu istnieje w jednym miejscu klienta, bo dwaj odbiorcy paneli — pas okien pomocniczych i kolumna paneli sceny — dzielą ten sam kod wykonawczy zamiast rozstrzygać kod pozycji osobno.
 */
export type WytworniaPanelu = (opcje: OpcjePanelu) => PanelPomocniczy;

/**
 * Wykaz wytwórni. Napisany wprost, pozycja po pozycji — bez pętli po rejestrze
 * i bez składania nazw funkcji z kodu pozycji. Dołożenie panelu to jeden wpis
 * tutaj i jeden w rejestrze, i nic poza tym.
 */
const WYTWORNIE: ReadonlyMap<string, WytworniaPanelu> = new Map<string, WytworniaPanelu>([
  ['podglad-bash', (opcje) => utworzOknoPodgladuBash(opcje)],
  // Historia rozmowy stoi na trzech komendach, nie jednej — odczyt, usunięcie i zapis zasady retencji.
  ['historia-rozmowy', (opcje) => utworzOknoHistoriiRozmowy(opcje)],
  ['terminal', (opcje) => utworzOknoTerminalaPomocnicze(opcje)],
  ['przebieg-debaty', (opcje) => utworzPanelDebaty(opcje)],
  [KOD_PANELU_ZASOBOW, (opcje) => utworzPanelZasobow(opcje)],
]);

/** Wytwórnia panelu o danym kodzie; wartość pusta znaczy, że panelu nie ma czym zbudować w tym kliencie. */
export function wytworniaPanelu(kod: string): WytworniaPanelu | null {
  return WYTWORNIE.get(kod) ?? null;
}

/**
 * Kody pozycji, które realnie się otworzą.
 *
 * Wykaz bierze się z mapy, a nie z drugiej listy obok niej: druga lista
 * mogłaby przeżyć usunięcie wytwórni i obiecywać panel, którego nie ma.
 */
export const KODY_Z_WYTWORNIA: readonly string[] = [...WYTWORNIE.keys()];
