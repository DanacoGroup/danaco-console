import type { Kanal } from '../protokol/kanal';
import { MODUL as MODUL_AGENTS } from '../moduly/agents/indeks';
import { MODUL as MODUL_WORKSPACE } from '../moduly/workspace/indeks';
import { MODUL as MODUL_AUTOMATIONS } from '../moduly/automations/indeks';
import { MODUL as MODUL_TERMINAL } from '../moduly/terminal/indeks';
import { MODUL as MODUL_MULTITASKING } from '../moduly/multitasking/indeks';
import { MODUL as MODUL_STUDIO } from '../moduly/studio/indeks';
import { MODUL as MODUL_LIBRARY } from '../moduly/library/indeks';
import { MODUL as MODUL_TRANSLATE } from '../moduly/translate/indeks';
import { MODUL as MODUL_BROWSER } from '../moduly/browser/indeks';
import { MODUL as MODUL_RESEARCH } from '../moduly/research/indeks';
import { MODUL as MODUL_DESIGN } from '../moduly/design/indeks';
import { MODUL as MODUL_ASSISTANT } from '../moduly/assistant/indeks';
import { MODUL as MODUL_APPS } from '../moduly/apps/indeks';
import { MODUL as MODUL_DEVELOPER } from '../moduly/developer/indeks';
import { MODUL as MODUL_DIAGNOSTICS } from '../moduly/diagnostics/indeks';
import { MODUL as MODUL_ROUNDTABLE } from '../moduly/roundtable/indeks';

/**
 * Rejestr widoków modułów wiąże kod modułu podany przez rdzeń z widokiem
 * zbudowanym w kliencie. Widok oddawany powłoce rozdziela element od wczytania,
 * ponieważ sesja bywa nieznana w chwili wyboru pozycji nawigacji.
 */
export interface WidokModulu {
  /** Element do postawienia w obszarze roboczym; moduł nie osadza go sam. */
  element: HTMLElement;
  /** Wczytuje zawartość modułu dla wskazanej sesji. */
  wczytaj(idSesji: string): Promise<void>;
  /** Odłącza nasłuch zdarzeń modułu; rozbiórkę widoku wykonują sprawdziany modułów. */
  zamknij?(): void;
}

/**
 * Opis modułu niesie to, co moduł wystawia powłoce ze swojego pliku `indeks.ts`.
 * Kod modułu stoi wewnątrz opisu, a nie w mapie rejestru, więc dopisanie modułu
 * sprowadza się do jednego importu i nie rozjeżdża się z rzeczywistością.
 */
export interface OpisModulu {
  /** Kod z kolumny `modul.kod` rdzenia, nie literał klienta. */
  kod: string;
  utworzWidok(kanal: Kanal): WidokModulu;
}

/**
 * Wykaz modułów zbudowanych w kliencie. Każdy moduł wystawia stałą `MODUL` ze
 * swojego pliku `indeks.ts`, więc wpis sprowadza się do jednego importu, a mapa
 * budowana z wykazu daje odczyt opisu po kodzie modułu.
 */
const OPISY: readonly OpisModulu[] = [
  // Moduły treściowe.
  MODUL_STUDIO,
  MODUL_LIBRARY,
  MODUL_TRANSLATE,
  MODUL_BROWSER,
  MODUL_RESEARCH,
  MODUL_DESIGN,
  MODUL_ASSISTANT,
  MODUL_APPS,

  // Moduły procesowe.
  MODUL_AGENTS,
  MODUL_WORKSPACE,
  MODUL_AUTOMATIONS,
  MODUL_TERMINAL,
  MODUL_MULTITASKING,

  // Developer i Roundtable montują się z okna sesji, bo ich komendy wymagają
  // identyfikatora okna.
  MODUL_DEVELOPER,
  MODUL_DIAGNOSTICS,
  MODUL_ROUNDTABLE,
];

const WEDLUG_KODU: ReadonlyMap<string, OpisModulu> = new Map(
  OPISY.map((opis) => [opis.kod, opis]),
);

/**
 * Zwraca opis modułu o podanym kodzie albo wartość `undefined`, gdy klient tego
 * modułu jeszcze nie zbudował. Brak opisu nie jest błędem, ponieważ rdzeń zna
 * moduły niezależnie od stanu budowy klienta.
 */
export function opisModulu(kod: string): OpisModulu | undefined {
  return WEDLUG_KODU.get(kod);
}
