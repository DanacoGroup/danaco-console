import './diagnostics.css';

import { utworzPasPomocniczych } from '../../okna-pomocnicze/indeks';
import type { Kanal } from '../../protokol/kanal';
import { widokZMontazu, type OpisModulu } from '../rejestracja';
import { utworzOknoDiagnosticsCenter } from './okno-diagnostics-center';
import { utworzOknoErrorsPanel } from './okno-errors-panel';
import { utworzOknoLogsViewera } from './okno-logs-viewer';
import { utworzZrodloProwenancji } from './prowenancja-zrodlo';
import { utworzOknoObservabilityTools } from './okno-observability-tools';
import { utworzOknoRecommendationsPanel } from './okno-recommendations-panel';
import { utworzStanDiagnostyki, type StanDiagnostyki } from './stan-diagnostyki';
import { utworzZrodloDiagnostics } from './zrodlo-diagnostics';
import { utworzZrodloObserwowalnosci } from './zrodlo-obserwowalnosci';
import { utworzZrodloZuzycia } from './zuzycie-zrodlo';

/**
 * Moduł Diagnostics — złożenie pięciu okien wokół jednego stanu systemu,
 * montowane bez okna, bo większość komend obszaru go nie wymaga.
 */
export interface ZamontowaneDiagnostics {
  /** Element osadzony w dokumencie. */
  element: HTMLElement;
  /** Stan diagnostyki wspólny czterem oknom. */
  stan: StanDiagnostyki;
  /** Odczytuje wszystkie okna z rdzenia. */
  odswiez(): void;
  /** Zamyka nasłuch zdarzeń modułu wraz z nasłuchem okien. */
  zamknij(): void;
}

export function zamontujDiagnostics(
  gospodarz: HTMLElement,
  kanal: Kanal,
  idOkna = '',
): ZamontowaneDiagnostics {
  const zrodlo = utworzZrodloDiagnostics(kanal);
  const obserwowalnosc = utworzZrodloObserwowalnosci(kanal);
  const stan = utworzStanDiagnostyki(zrodlo);

  const obszar = document.createElement('div');
  obszar.className = 'dg-modul';
  obszar.dataset['modul'] = 'diagnostics';

  const centrum = utworzOknoDiagnosticsCenter(zrodlo, stan, idOkna);
  const dziennik = utworzOknoLogsViewera(zrodlo, stan);
  const bledy = utworzOknoErrorsPanel(zrodlo, stan);
  const rekomendacje = utworzOknoRecommendationsPanel(zrodlo, stan);
  const narzedzia = utworzOknoObservabilityTools(
    zrodlo,
    obserwowalnosc,
    stan,
    idOkna,
    utworzZrodloZuzycia(kanal),
    utworzZrodloProwenancji(kanal),
  );

  const gorny = document.createElement('div');
  gorny.className = 'dg-modul__gora';
  gorny.append(
    oznacz(centrum.element, 'diagnostics-center'),
    oznacz(dziennik.element, 'logs-viewer'),
  );

  const dolny = document.createElement('div');
  dolny.className = 'dg-modul__dol';
  dolny.append(
    oznacz(bledy.element, 'errors-panel'),
    oznacz(rekomendacje.element, 'recommendations-panel'),
  );

  const ekspercki = document.createElement('div');
  ekspercki.className = 'dg-modul__ekspercki';
  ekspercki.append(oznacz(narzedzia.element, 'observability-tools'));

  const pomocnicze = utworzPasPomocniczych({
    kanal,
    modul: KOD_MODULU,
    nazwaModulu: NAZWA_MODULU,
    okno: idOkna,
    przedrostek: 'dg',
  });

  obszar.append(gorny, dolny, ekspercki, pomocnicze.element);
  gospodarz.replaceChildren(obszar);

  function odswiez(): void {
    centrum.odswiez();
    dziennik.odswiez();
    bledy.odswiez();
    rekomendacje.odswiez();
    narzedzia.odswiez();
    pomocnicze.odswiez();
  }

  odswiez();

  return {
    element: obszar,
    stan,
    odswiez,
    zamknij() {
      // Pas zamyka się pierwszy: trzyma subskrypcję, która żyje niezależnie od stanu modułu.
      pomocnicze.zamknij();
      narzedzia.zamknij();
      rekomendacje.zamknij();
      bledy.zamknij();
      dziennik.zamknij();
      centrum.zamknij();
      stan.zamknij();
    },
  };
}

/** Znakuje okno kodem rejestru okien operacyjnych, po którym skacze nawigacja między oknami modułu Diagnostics. */
function oznacz(element: HTMLElement, kod: string): HTMLElement {
  element.dataset['okno'] = kod;
  element.tabIndex = -1;
  return element;
}

/**
 * Samoopisujący się moduł dla rejestru powłoki: kod siedzi w module, nie
 * w mapie po stronie powłoki systemu.
 */
const KOD_MODULU = 'diagnostics';

/** Nazwa modułu w etykietach dostępności pasa okien pomocniczych modułu Diagnostics w oknie sesji rdzenia. */
const NAZWA_MODULU = 'Diagnostics';

export const MODUL: OpisModulu = {
  kod: KOD_MODULU,
  utworzWidok: (kanal) => widokZMontazu(zamontujDiagnostics, kanal),
};
