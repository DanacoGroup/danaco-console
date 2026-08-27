import './research.css';

import type { Kanal } from '../../protokol/kanal';
import { utworzOknoDiscoveryPanel } from './okno-discovery-panel';
import { utworzOknoExportPanel } from './okno-export-panel';
import { utworzOknoFindingsPanel } from './okno-findings-panel';
import { utworzOknoReadingView } from './okno-reading-view';
import { utworzOknoReportBuilder } from './okno-report-builder';
import { utworzOknoResearchWorkspace } from './okno-research-workspace';
import { utworzOknoSourcesManager } from './okno-sources-manager';
import { utworzStanBadania, type StanBadania } from './stan-badania';

/**
 * Moduł Research łączy siedem okien operacyjnych w jednym układzie zgodnym z porządkiem pracy badawczej, dzieląc jeden wspólny stan badania.
 */
export interface ModulResearch {
  /** Element osadzany w obszarze roboczym powłoki. */
  element: HTMLElement;
  /** Zleca odczyt okna badania z rdzenia. */
  wczytaj(idSesji: string): Promise<void>;
  /** Odłącza subskrypcje zdarzeń rdzenia i zdejmuje kreator z dokumentu. */
  rozlacz(): void;
}

export function utworzModulResearch(kanal: Kanal): ModulResearch {
  const stan: StanBadania = utworzStanBadania(kanal);

  const workspace = utworzOknoResearchWorkspace(stan, (kod) => przenieOgnisko(kod));
  const odkrywanie = utworzOknoDiscoveryPanel(stan, (kod) => przenieOgnisko(kod));
  const zrodla = utworzOknoSourcesManager(stan, (kod) => przenieOgnisko(kod));
  const lektura = utworzOknoReadingView(stan, (kod) => przenieOgnisko(kod));
  const ustalenia = utworzOknoFindingsPanel(stan, (kod) => przenieOgnisko(kod));
  const raport = utworzOknoReportBuilder(stan, (kod) => przenieOgnisko(kod));
  const eksport = utworzOknoExportPanel(stan);

  const element = document.createElement('div');
  element.className = 'mr-modul';
  element.dataset['modul'] = 'research';
  element.setAttribute('aria-label', 'Moduł Research — okna operacyjne');
  element.append(
    workspace.element,
    pas([odkrywanie.element, zrodla.element]),
    pas([lektura.element, ustalenia.element]),
    pas([raport.element, eksport.element]),
  );

  /** Przenosi ognisko do okna wskazanego kodem katalogu rdzenia. */
  function przenieOgnisko(kodOkna: string): void {
    const cel = element.querySelector<HTMLElement>(`[data-okno="${kodOkna}"]`);
    if (cel === null) return;
    for (const inne of element.querySelectorAll<HTMLElement>('[data-okno]')) {
      delete inne.dataset['ognisko'];
    }
    cel.dataset['ognisko'] = 'tak';
    cel.focus();
    cel.scrollIntoView({ block: 'nearest' });
  }

  function odswiezWszystkie(): void {
    workspace.odswiez();
    odkrywanie.odswiez();
    zrodla.odswiez();
    lektura.odswiez();
    ustalenia.odswiez();
    raport.odswiez();
    eksport.odswiez();
  }

  const odsubskrybuj = stan.obserwuj(odswiezWszystkie);

  return {
    element,

    async wczytaj(idSesji) {
      await stan.odswiez(idSesji);
    },

    rozlacz() {
      odsubskrybuj();
      raport.rozlacz();
      stan.rozlacz();
    },
  };
}

/** Funkcja tworzy pas układu z dwóch okien ułożonych obok siebie, zwijany do jednej kolumny na wąskim ekranie. */
function pas(okna: readonly HTMLElement[]): HTMLElement {
  const element = document.createElement('div');
  element.className = 'mr-modul__pas';
  element.append(...okna);
  return element;
}
