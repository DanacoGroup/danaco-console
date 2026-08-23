import { utworzNaglowekOkna } from '../../komponenty/naglowek-okna';
import { utworzZakladkiSekcji } from '../../modele/zakladki-sekcji';
import { utworzEdytorMakra } from './edytor-makra';
import { utworzPanelNarzedzi, type PanelNarzedzi } from './panel-narzedzi';
import { utworzPanelRutyn, type PanelRutyn } from './panel-rutyn';
import { utworzPanelSchowka, type PanelSchowka } from './panel-schowka';
import { utworzSiatkeAkcji, type SiatkaAkcji } from './siatka-akcji';
import type { StanAssistant } from './stan-assistant';
import type { ZrodloNarzedzi } from './zrodlo-narzedzi';

/** Kod okna operacyjnego modułu; rdzeń nie ma go dziś w katalogu okien. */
export const KOD_OKNA = 'command-tools-hub';

/**
 * Command & Tools Hub — okno zarządcy szybkich akcji, narzędzi i rutyn modułu
 * Assistant.
 *
 * Okno zamyka obszary szybkich akcji, wywołań narzędzi oraz proaktywności
 * modułu. Zakładek jest cztery i każda stoi nad rzeczywistą komendą kontraktu:
 *
 *   Akcje            `action.list` — katalog akcji zasięgu modułu;
 *   Makra            `automation.workflow.save` — jedyne trwałe miejsce dla
 *                    sekwencji kroków, wskazane wprost przez opracowanie modułu
 *                    („→ Wyślij do Automations");
 *   Narzędzia i MCP  `tools.catalog.list` wraz z `session.tool.*`;
 *   Rutyny           `automation.workflow.list`, `schedule.get`
 *                    i `automation.schedule.set`.
 *
 * Zakładki „Umiejętności" osobno nie ma z rozstrzygnięcia, nie z przeoczenia:
 * kontrakt trzyma narzędzia i umiejętności w jednym katalogu i rozróżnia je
 * polem `kind` oraz przedrostkiem źródła. Druga zakładka nad tą samą komendą
 * udawałaby drugie źródło; rozróżnienie robi filtr rodzaju w zakładce narzędzi.
 *
Piąta zakładka — „Skróty i schowek" — stoi nad trzema rodzinami, które
 * kontrakt niesie w całości: `clipboard.*` (trwała historia schowka),
 * `snippet.*` (słownik skrótów rozwijanych we wszystkich polach platformy)
 * oraz `launcher.hotkey.*` (skrót globalny wywoływacza). Schowka maszyny
 * Operatora rdzeń nie czyta i zakładka tego nie udaje — podział ról jest
 * widoczny na ekranie (`panel-schowka.ts`).
 *
 * Plik odpowiada wyłącznie za skład okna; wywołania mieszkają
 * w `zrodlo-narzedzi.ts`, a każda zakładka ma własny plik obszaru.
 */
export interface OknoCommandToolsHub {
  element: HTMLElement;
  /** Odczyt katalogu akcji, katalogu narzędzi i wykazu rutyn. */
  wczytaj(): Promise<void>;
}

export function utworzOknoCommandToolsHub(
  stan: StanAssistant,
  zrodlo: ZrodloNarzedzi,
  naAkcje: (tresc: string) => void,
): OknoCommandToolsHub {
  // Siatka akcji jest tym samym widokiem katalogu, który stoi w Voice Console.
  // Drugi widok, a nie drugie źródło: obie sięgają `action.list` zasięgu modułu,
  // a naciśnięcie kafla wypełnia to samo pole polecenia.
  const akcje: SiatkaAkcji = utworzSiatkeAkcji(stan, naAkcje);
  const makra = utworzEdytorMakra(zrodlo);
  const narzedzia: PanelNarzedzi = utworzPanelNarzedzi(stan, zrodlo);
  const rutyny: PanelRutyn = utworzPanelRutyn(zrodlo);
  const schowek: PanelSchowka = utworzPanelSchowka(stan);

  const zakladki = utworzZakladkiSekcji([
    { kod: 'akcje', nazwa: 'Akcje', element: akcje.element },
    { kod: 'makra', nazwa: 'Makra', element: makra.element },
    { kod: 'narzedzia', nazwa: 'Narzędzia i MCP', element: narzedzia.element },
    { kod: 'rutyny', nazwa: 'Rutyny', element: rutyny.element },
    { kod: 'schowek-i-skroty', nazwa: 'Skróty i schowek', element: schowek.element },
  ]);

  const element = document.createElement('section');
  element.className = 'ma-okno ma-okno--zarzadca';
  element.dataset['okno'] = KOD_OKNA;
  element.append(
    utworzNaglowekOkna({
      tytul: 'Command & Tools Hub',
      rola: 'zarządca · akcje, makra, narzędzia i rutyny asystenta',
    }),
    zakladki.element,
  );

  return {
    element,
    // Trzy odczyty idą równolegle: katalog akcji, katalog narzędzi i wykaz
    // rutyn dotyczą trzech różnych rodzin komend i żaden nie warunkuje
    // pozostałych. Każdy nazywa swoje niepowodzenie w swoim obszarze.
    wczytaj: async () => {
      await Promise.all([
        akcje.wczytaj(),
        narzedzia.wczytaj(),
        rutyny.wczytaj(),
        schowek.wczytaj(),
      ]);
    },
  };
}

