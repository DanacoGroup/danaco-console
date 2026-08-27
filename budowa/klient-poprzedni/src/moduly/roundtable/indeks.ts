import './roundtable.css';
import './debata.css';
import './analiza-debaty.css';

import type { Kanal } from '../../protokol/kanal';
import { utworzRejestrKanalow, type RejestrKanalow } from '../../sterowanie/rejestr-kanalow';
import { widokZOknaSesji, type OpisModulu } from '../rejestracja';
import { utworzOknoArgumentMap } from './okno-argument-map';
import { utworzOknoConsensusPanel } from './okno-consensus-panel';
import { utworzOknoDebatePanel } from './okno-debate-panel';
import { utworzOknoModelPanels } from './okno-model-panels';
import { utworzOknoModeratorPanel } from './okno-moderator-panel';
import { utworzOknoVotingEvaluation } from './okno-voting-evaluation';
import { utworzPasekUczciwosci } from './pasek-uczciwosci';
import { utworzPasRozszerzen } from './rozszerzenia-boczne';
import { utworzStanDebaty, type StanDebaty } from './stan-debaty';
import { utworzStrumienWypowiedzi } from './strumien-wypowiedzi';
import { utworzZrodloArsenaluRoundtable } from './zrodlo-arsenalu';
import { utworzZrodloRoundtable } from './zrodlo-roundtable';
import { utworzZrodloStrumieniaDebaty } from './zrodlo-strumienia-debaty';

/**
 * Moduł Roundtable składa sześć okien wokół jednej debaty wielu modeli, po wspólnym rejestrze
 * kanałów, którym jedzie okno rozmowy.
 */
export interface ZamontowanyRoundtable {
  /** Element osadzony w dokumencie. */
  element: HTMLElement;
  /** Stan debaty wspólny sześciu oknom. */
  stan: StanDebaty;
  /** Odczytuje wszystkie okna z rdzenia. */
  odswiez(): void;
  /** Zamyka nasłuch zdarzeń modułu wraz z nasłuchem okien. */
  zamknij(): void;
}

export function zamontujRoundtable(
  gospodarz: HTMLElement,
  kanal: Kanal,
  okno: string,
  rejestrPodany?: RejestrKanalow,
): ZamontowanyRoundtable {
  // Rejestr przyjmowany z zewnątrz, gdy wywołujący go ma — inaczej zakładany tutaj.
  const rejestr = rejestrPodany ?? utworzRejestrKanalow(kanal);
  const zrodlo = utworzZrodloRoundtable(kanal);
  // Arsenał obszaru: czterdzieści dwie komendy poza czwórką prowadzącą debatę, jedno źródło na złożenie.
  const arsenal = utworzZrodloArsenaluRoundtable(kanal);
  const stan = utworzStanDebaty(zrodlo, rejestr, { okno });
  const strumien = utworzStrumienWypowiedzi(stan);
  // Okno czytane w chwili nadejścia fragmentu, nie w chwili subskrypcji filtra strumienia.
  const zrodloStrumienia = utworzZrodloStrumieniaDebaty(kanal, () => stan.okno());

  // Zmiana tury czyści gromadzenie przed oknami, żeby nie dokleić zdań tury poprzedniej do mówców nowej.
  let turaWidziana = stan.tura();
  const odsubskrybujTure = stan.naZmiane(() => {
    if (stan.tura() === turaWidziana) return;
    turaWidziana = stan.tura();
    strumien.wyczysc();
  });

  const obszar = document.createElement('div');
  obszar.className = 'dr-modul';
  obszar.dataset['modul'] = 'roundtable';

  const sklad = utworzOknoModelPanels(zrodlo, arsenal, stan, strumien, okno);
  const moderator = utworzOknoModeratorPanel(zrodlo, arsenal, stan, okno);
  const stanowisko = utworzOknoConsensusPanel(zrodlo, arsenal, stan, okno);
  const analiza = utworzOknoArgumentMap(arsenal, stan, strumien);
  const ocena = utworzOknoVotingEvaluation(arsenal, stan, strumien);

  // Pas rozszerzeń powstaje przed Debate Panelem, który stawia jego przyciski w swoim pasku akcji.
  const rozszerzenia = utworzPasRozszerzen([
    {
      kod: 'argument-map-analysis',
      nazwa: 'Argument Map & Analysis',
      element: oznacz(analiza.element, 'argument-map-analysis'),
    },
    {
      kod: 'voting-evaluation-center',
      nazwa: 'Voting & Evaluation Center',
      element: oznacz(ocena.element, 'voting-evaluation-center'),
    },
    {
      kod: 'moderator-panel',
      nazwa: 'Moderator Panel',
      element: oznacz(moderator.element, 'moderator-panel'),
    },
    {
      kod: 'consensus-panel',
      nazwa: 'Consensus Panel',
      element: oznacz(stanowisko.element, 'consensus-panel'),
    },
  ]);

  const przebieg = utworzOknoDebatePanel(zrodlo, arsenal, stan, strumien, [
    rozszerzenia.przyciskOtwarcia('argument-map-analysis'),
    rozszerzenia.przyciskOtwarcia('voting-evaluation-center'),
    rozszerzenia.przyciskOtwarcia('moderator-panel'),
    rozszerzenia.przyciskOtwarcia('consensus-panel'),
  ]);

  const gorny = document.createElement('div');
  gorny.className = 'dr-modul__gora';
  gorny.append(
    oznacz(sklad.element, 'model-panels'),
    oznacz(przebieg.element, 'debate-panel'),
  );

  const uczciwosc = utworzPasekUczciwosci(kanal);

  obszar.append(uczciwosc.element, gorny, rozszerzenia.element);
  gospodarz.replaceChildren(obszar);

  function odswiez(): void {
    sklad.odswiez();
    przebieg.odswiez();
    moderator.odswiez();
    stanowisko.odswiez();
    analiza.odswiez();
    ocena.odswiez();
  }

  // Katalog kanałów zamawia złożenie, nie okno, przed odczytem okien, żeby uniknąć podwójnego zapytania.
  rejestr.odswiez();
  odswiez();
  // Katalog okien rdzenia czyta się raz na montaż — pasek uczciwości liczy z niego rozjazd okien.
  void uczciwosc.katalog.odczytaj();

  // Fragment budzi wyłącznie okna rysujące treść wypowiedzi, nie Moderator Panel ani Consensus Panel.
  const odsubskrybujStrumien = zrodloStrumienia.naFragmentWypowiedzi((fragment) => {
    strumien.przyjmij(fragment);
    sklad.odswiezGlosy();
    przebieg.odswiezGlosy();
    analiza.odswiezGlosy();
    ocena.odswiezGlosy();
  });

  return {
    element: obszar,
    stan,
    odswiez,
    zamknij() {
      odsubskrybujStrumien();
      odsubskrybujTure();
      uczciwosc.katalog.zamknij();
      ocena.zamknij();
      analiza.zamknij();
      stanowisko.zamknij();
      moderator.zamknij();
      przebieg.zamknij();
      sklad.zamknij();
      stan.zamknij();
    },
  };
}

/** Znakuje okno kodem rejestru okien operacyjnych, po którym skacze nawigacja, znacząc też jego indeks skupienia. */
function oznacz(element: HTMLElement, kod: string): HTMLElement {
  element.dataset['okno'] = kod;
  element.tabIndex = -1;
  return element;
}

/**
 * Kod modułu dla rejestru powłoki.
 *
 * Kod siedzi w module, a nie w mapie po stronie powłoki: dodanie modułu to jeden
 * wpis, a nie dwa, więc nie da się dodać modułu i zapomnieć o wytwórni.
 */
const KOD_MODULU = 'roundtable';

/**
 * Panel debaty dla stosu paneli pomocniczych — wyłącznie reeksport dla czytającego, bo panelu
 * tu się nie stawia.
 */
export { utworzPanelDebaty, KOD_PANELU_DEBATY } from './panel-debaty';

export const MODUL: OpisModulu = {
  kod: KOD_MODULU,
  utworzWidok: (kanal) =>
    widokZOknaSesji(
      kanal,
      (gospodarz, kanalOkna, okno) => zamontujRoundtable(gospodarz, kanalOkna, okno),
      KOD_MODULU,
    ),
};
