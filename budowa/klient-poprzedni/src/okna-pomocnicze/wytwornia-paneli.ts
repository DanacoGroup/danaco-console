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
 * Mapowanie kod pozycji → wytwórnia panelu; jedno miejsce w całym kliencie.
 *
 * Odbiorcy paneli są dwaj — pas okien pomocniczych modułu i kolumna paneli
 * sceny okien równoległych. Gdyby każdy z nich rozstrzygał kod pozycji własnym
 * `if`, byłyby to dwa miejsca do rozjechania się; tutaj kod pozycji spotyka się
 * z kodem wykonawczym raz.
 *
 * Wytwórnie wpisuje się strukturalnie, bez przejściówek. Panel wystarczy, że
 * niesie `element`, `odswiez()` i `zamknij()` — nadmiarowe pola (jak
 * `ustawOkno` terminala) nie przeszkadzają, a pole wymagane przechodzi
 * w miejsce opcjonalnego. TypeScript wiąże strukturalnie, więc przejściówka
 * byłaby warstwą bez treści: gdy okno rozjedzie się z umową, kompilator
 * zatrzyma się na wpisie w mapie poniżej.
 *
 * Plik nie czyta rejestru i nie zna stanów pozycji. „Ma wytwórnię" i „rejestr
 * nazywa ją zbudowaną" to dwie różne prawdy; zestawia je
 * `panele-otwieralne.ts`, a rozjazd między nimi zgłasza pas.
 */
export type WytworniaPanelu = (opcje: OpcjePanelu) => PanelPomocniczy;

/**
 * Wykaz wytwórni. Napisany wprost, pozycja po pozycji — bez pętli po rejestrze
 * i bez składania nazw funkcji z kodu pozycji. Dołożenie panelu to jeden wpis
 * tutaj i jeden w rejestrze, i nic poza tym.
 */
const WYTWORNIE: ReadonlyMap<string, WytworniaPanelu> = new Map<string, WytworniaPanelu>([
  ['podglad-bash', (opcje) => utworzOknoPodgladuBash(opcje)],
  // Historia rozmowy stoi na trzech komendach, nie na jednej: panel woła
  // `history.load`, `history.delete` i `retention.set`.
  ['historia-rozmowy', (opcje) => utworzOknoHistoriiRozmowy(opcje)],
  ['terminal', (opcje) => utworzOknoTerminalaPomocnicze(opcje)],
  ['przebieg-debaty', (opcje) => utworzPanelDebaty(opcje)],
  [KOD_PANELU_ZASOBOW, (opcje) => utworzPanelZasobow(opcje)],
]);

/** Wytwórnia panelu o danym kodzie; `null` = panelu nie ma czym zbudować. */
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
