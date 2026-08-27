import { utworzNaglowekOkna } from '../../komponenty/naglowek-okna';
import { utworzZakladkiSekcji } from '../../modele/zakladki-sekcji';
import { utworzEdytorMakra } from './edytor-makra';
import { utworzPanelNarzedzi, type PanelNarzedzi } from './panel-narzedzi';
import { utworzPanelRutyn, type PanelRutyn } from './panel-rutyn';
import { utworzPanelSchowka, type PanelSchowka } from './panel-schowka';
import { utworzSiatkeAkcji, type SiatkaAkcji } from './siatka-akcji';
import type { StanAssistant } from './stan-assistant';
import type { ZrodloNarzedzi } from './zrodlo-narzedzi';

/**
 * Kod okna operacyjnego modułu. Rdzeń nie ma go dziś w katalogu okien, więc okno
 * przedstawia się tym napisem samo, a wpis kontraktu zajmie jego miejsce bez
 * zmiany składu okna.
 */
export const KOD_OKNA = 'command-tools-hub';

/**
 * Command & Tools Hub — okno zarządcy szybkich akcji, narzędzi i rutyn modułu
 * Assistant. Cztery pierwsze zakładki stoją nad rzeczywistymi komendami
 * kontraktu: akcje, makra, narzędzia wraz z protokołem MCP oraz rutyny.
Piąta zakładka — „Skróty i schowek" — stoi nad trzema rodzinami, które
 * kontrakt niesie w całości: trwałą historią schowka, słownikiem skrótów
 * rozwijanych we wszystkich polach platformy oraz skrótem globalnym
 * wywoływacza. Skład okna jest jedyną odpowiedzialnością tego pliku.
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
  // Siatka akcji jest drugim widokiem katalogu `action.list`, a nie drugim źródłem.
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
    // Odczyty idą równolegle; żaden nie warunkuje pozostałych i każdy melduje osobno.
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

