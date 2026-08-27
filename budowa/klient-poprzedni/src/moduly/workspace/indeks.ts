import './workspace.css';

import type { Kanal } from '../../protokol/kanal';
import type { OpisModulu } from '../rejestracja';
import { utworzWykazBrakow } from './braki-kontraktu';
import { czynnosciPlanowania } from './czynnosci-planowania';
import { czynnosciWiedzy } from './czynnosci-wiedzy';
import { utworzOknoHubuPlanowania } from './hub-planowania';
import { utworzOknoWikiProjektu } from './wiki-projektu';
import { utworzOknoWspolpracy } from './wspolpraca-projektu';
import { utworzOknoAgentow } from './zarzadca-agentow';
import { utworzOknoBiblioteki } from './biblioteka-projektu';
import { utworzOknoInstrukcji } from './panel-instrukcji';
import { utworzOknoPamieci } from './pamiec-kontekstu';
import { odnajdzOknoRozmowy } from './okno-rozmowy';
import { utworzOknoPulpitu } from './pulpit-projektu';
import { utworzStanProjektu, type StanProjektu } from './stan-projektu';
import { utworzZrodloWorkspace } from './zrodlo-workspace';

/**
 * Moduł Workspace jest złożeniem pięciu okien operacyjnych wokół jednego projektu, ułożonych
 * według roli: pulpit wiodący, pozostałe pomocnicze wokół niego.
 */
export interface ZamontowanyWorkspace {
  /** Element osadzony w dokumencie. */
  element: HTMLElement;
  /** Stan projektu wspólny pięciu oknom. */
  stan: StanProjektu;
  /** Odczytuje wszystkie okna z rdzenia. */
  odswiez(): void;
  /** Zamyka nasłuch zdarzeń modułu. */
  zamknij(): void;
}

/** Zależności złożenia niesie projekt otwierany od razu oraz okno rozmowy jako nośnik przeniesienia kontekstu sesji. */
export interface OpcjeWorkspace {
  /** Projekt otwierany od razu; pusty zostawia wskazanie Operatorowi. */
  projekt?: string;
  /** Okno rozmowy modułu — nośnik przeniesienia kontekstu. */
  oknoRozmowy?: string;
}

export function zamontujWorkspace(
  gospodarz: HTMLElement,
  kanal: Kanal,
  opcje: OpcjeWorkspace = {},
): ZamontowanyWorkspace {
  const zrodlo = utworzZrodloWorkspace(kanal);
  const stan = utworzStanProjektu(zrodlo, opcje);
  // Jeden wykaz braków na pięć okien: rdzeń pytany o komendy raz, kontrolki biorą z odpowiedzi zdanie.
  const braki = utworzWykazBrakow(kanal);

  const obszar = document.createElement('div');
  obszar.className = 'dw-modul';
  obszar.dataset['modul'] = 'workspace';

  // Trzy okna planowania i wiedzy jadą własnymi czynnościami, nie źródłem pięciu okien pierwotnych.
  const planowanie = czynnosciPlanowania(kanal);
  const wiedza = czynnosciWiedzy(kanal);

  const pulpit = utworzOknoPulpitu(
    zrodlo,
    wiedza,
    stan,
    {
      pokaz(kodOkna) {
        const okno = obszar.querySelector<HTMLElement>(`[data-okno="${kodOkna}"]`);
        okno?.scrollIntoView({ behavior: 'smooth', block: 'start' });
        okno?.focus({ preventScroll: true });
      },
    },
    braki,
  );
  const instrukcje = utworzOknoInstrukcji(zrodlo, stan, braki);
  const pamiec = utworzOknoPamieci(zrodlo, stan);
  const biblioteka = utworzOknoBiblioteki(zrodlo, stan, braki);
  const agenci = utworzOknoAgentow(zrodlo, stan, braki);
  const hub = utworzOknoHubuPlanowania(planowanie, wiedza, stan);
  const wiki = utworzOknoWikiProjektu(wiedza, stan);
  const wspolpraca = utworzOknoWspolpracy(wiedza, stan);

  const gorny = document.createElement('div');
  gorny.className = 'dw-modul__gora';
  gorny.append(oznacz(pulpit.element, 'project-dashboard'), oznacz(instrukcje.element, 'instructions-panel'));

  const dolny = document.createElement('div');
  dolny.className = 'dw-modul__dol';
  dolny.append(
    oznacz(pamiec.element, 'context-memory'),
    oznacz(biblioteka.element, 'project-library'),
    oznacz(agenci.element, 'agent-manager'),
  );

  // Pas trzeci: hub planowania, wiki i czynności przekrojowe, stoją pod pasem zarządców materiału.
  const planowanieRzad = document.createElement('div');
  planowanieRzad.className = 'dw-modul__dol';
  planowanieRzad.append(
    oznacz(hub.element, 'planning-hub'),
    oznacz(wiki.element, 'project-wiki'),
    oznacz(wspolpraca.element, 'project-collaboration'),
  );

  obszar.append(gorny, dolny, planowanieRzad);
  gospodarz.replaceChildren(obszar);

  function odswiez(): void {
    pulpit.odswiez();
    instrukcje.odswiez();
    pamiec.odswiez();
    biblioteka.odswiez();
    agenci.odswiez();
    hub.odswiez();
    wiki.odswiez();
    wspolpraca.odswiez();
  }

  odswiez();
  // Pytanie o komendy idzie równolegle z odczytem okien: kontrolki mówią, że pytanie jest w drodze.
  void braki.odczytaj();

  return {
    element: obszar,
    stan,
    odswiez,
    zamknij: () => stan.zamknij(),
  };
}

/**
 * Czyni okno celem skoku nawigacji.
 *
 * Kod okna nadaje rama, której każde z pięciu okien podaje swój `kod`
 * (`komponenty/rama-okna.ts`); ta funkcja tylko go potwierdza, gdy rama go nie
 * ustawiła, i ustawia `tabIndex` wymagany przez skok ogniska.
 */
function oznacz(element: HTMLElement, kod: string): HTMLElement {
  if (element.dataset['okno'] !== kod) element.dataset['okno'] = kod;
  element.tabIndex = -1;
  return element;
}

/**
 * Samoopisujący się moduł dla rejestru powłoki.
 *
 * Kod siedzi w module, nie w mapie po stronie powłoki: dodanie modułu to jeden
 * wpis, a nie dwa, więc nie da się dodać modułu i zapomnieć o wytwórni.
 */
export const MODUL: OpisModulu = {
  kod: 'workspace',
  utworzWidok: (kanal) => {
    const gospodarz = document.createElement('div');
    const zamontowane = zamontujWorkspace(gospodarz, kanal);
    return {
      element: gospodarz,
      // Karta sesji wchodzi wyłącznie tędy, moduł nie zgaduje jej z ogniska ani nie pyta o nią drugą drogą.
      wczytaj: async (idSesji: string) => {
        zamontowane.stan.ustawSesje(idSesji);
        if (zamontowane.stan.oknoRozmowy() !== '') return;
        zamontowane.stan.ustawOknoRozmowy(await odnajdzOknoRozmowy(kanal, idSesji));
      },
      zamknij: () => zamontowane.zamknij(),
    };
  },
};
