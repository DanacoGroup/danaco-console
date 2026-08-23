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
 * Rejestr widoków modułów — jedno miejsce, w którym kod modułu z rdzenia
 * spotyka się z widokiem zbudowanym w kliencie.
 *
 * Nawigacja bierze wykaz modułów z rdzenia (`module.list`,
 * `environment.enter`), więc powłoka poznaje kody dopiero w czasie działania
 * i nie może o nich nic zakładać. Bez rejestru każdy moduł musiałby wpinać się
 * sam, a przestrzeń robocza trzymałaby własną kopię katalogu modułów.
 *
 * Moduł nieobecny w rejestrze nie jest błędem: rdzeń może poznać moduł,
 * którego klient jeszcze nie zbudował, a taki moduł dostaje stan pusty z paska
 * uczciwości przestrzeni roboczej zamiast martwego kliknięcia.
 */

/**
 * Widok modułu oddawany powłoce.
 *
 * Element i wczytanie są rozdzielone, bo sesja bywa jeszcze nieznana w chwili
 * wyboru pozycji nawigacji: przestrzeń robocza odkłada wtedy wejście i ponawia
 * je po otwarciu okna przez rdzeń. Gdyby wytwórnia oddawała sam element, każdy
 * moduł musiałby tę samą zwłokę obsłużyć u siebie.
 */
export interface WidokModulu {
  /** Element do postawienia w obszarze roboczym; moduł nie osadza go sam. */
  element: HTMLElement;
  /** Wczytuje zawartość modułu dla wskazanej sesji. */
  wczytaj(idSesji: string): Promise<void>;
  /**
   * Odłącza nasłuch zdarzeń modułu. Szew rozbiórki widoku; powłoka go nie woła,
   * bo nie ma w niej chwili, w której byłby prawdziwy: `przestrzen-modulu.ts`
   * woła `utworzWidok` raz na kod modułu i oddaje ten sam widok przy każdym
   * powrocie, a `powloka/obszar-roboczy.ts` nie zdejmuje elementu z drzewa,
   * tylko przestawia `hidden`. Liczba żywych subskrypcji jest przez to
   * ograniczona z góry liczbą modułów i nie rośnie z liczbą przełączeń pozycji.
   *
   * Pole zostaje opcjonalne i zostaje w umowie: rozbiórkę wykonują sprawdziany
   * modułów. Gdy w powłoce pojawi się pierwsza granica życia widoku, `zamknij`
   * trzeba wołać razem ze zdejmowaniem elementu w `obszar-roboczy.ts` — jedno
   * bez drugiego jest usterką.
   */
  zamknij?(): void;
}

/**
 * Opis modułu — to, co moduł wystawia powłoce ze swojego `indeks.ts`.
 *
 * Kod stoi wewnątrz opisu, nie w mapie rejestru: moduł sam mówi, którym jest
 * modułem, więc dopisanie go tutaj to jeden import zamiast pary kod↔wytwórnia,
 * a rozjazd między nazwą w mapie a rzeczywistością staje się niemożliwy.
 */
export interface OpisModulu {
  /** Kod z kolumny `modul.kod` rdzenia, nie literał klienta. */
  kod: string;
  utworzWidok(kanal: Kanal): WidokModulu;
}

/**
 * Wykaz modułów zbudowanych w kliencie. Każdy wystawia `MODUL` ze swojego
 * `indeks.ts`, więc wpis to jeden import.
 *
 * Wpis MultitaskingAI niesie kod `multitaskingai`, a to kod środowiska
 * (`srodowisko.kod`), nie modułu. Nawigacja tego środowiska podaje sekcje
 * orkiestracji (`powloka/srodowiska.ts`, `SEKCJE_ORKIESTRACJI`), a nie kod
 * `multitaskingai`, więc `opisModulu` tego wpisu z nawigacji nie znajdzie.
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
  // `windowId`; Diagnostics montuje się od razu, bo większość jego komend
  // okna nie wymaga.
  MODUL_DEVELOPER,
  MODUL_DIAGNOSTICS,
  MODUL_ROUNDTABLE,
];

const WEDLUG_KODU: ReadonlyMap<string, OpisModulu> = new Map(
  OPISY.map((opis) => [opis.kod, opis]),
);

/** Opis modułu albo `undefined`, gdy modułu jeszcze nie zbudowano. */
export function opisModulu(kod: string): OpisModulu | undefined {
  return WEDLUG_KODU.get(kod);
}
