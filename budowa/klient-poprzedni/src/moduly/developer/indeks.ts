import './developer.css';

import { utworzPasPomocniczych } from '../../okna-pomocnicze/indeks';
import type { Kanal } from '../../protokol/kanal';
import { widokZOknaSesji, type OpisModulu } from '../rejestracja';
import { utworzOknoBudowania } from './okno-build-output';
import { utworzOknoCodeEditora } from './okno-code-editor';
import { utworzOknoDevTools } from './okno-dev-tools';
import { utworzOknoGitPanelu } from './okno-git-panel';
import { utworzOknoProjectTree } from './okno-project-tree';
import { utworzOknoWarsztatuKodu } from './okno-warsztatu-kodu';
import { utworzStanDevelopera, type StanDevelopera } from './stan-developer';
import { utworzZrodloDeveloper } from './zrodlo-developer';
import { utworzZrodloWarsztatu } from './zrodlo-warsztatu';

/**
 * Moduł Developer — złożenie pięciu okien operacyjnych wokół jednego
 * repozytorium, montowane w oknie identyfikowanym kodem, nie w sesji.
 */
export interface ZamontowanyDeveloper {
  /** Element osadzony w dokumencie. */
  element: HTMLElement;
  /** Stan repozytorium wspólny czterem oknom. */
  stan: StanDevelopera;
  /** Odczytuje wszystkie okna z rdzenia. */
  odswiez(): void;
  /** Zamyka nasłuch zdarzeń modułu wraz z nasłuchem okien. */
  zamknij(): void;
}

/** Nazwa modułu w etykietach dostępności pasa okien pomocniczych modułu Developer w oknie sesji rdzenia. */
const NAZWA_MODULU = 'Developer';

export function zamontujDeveloper(
  gospodarz: HTMLElement,
  kanal: Kanal,
  okno: string,
): ZamontowanyDeveloper {
  const zrodlo = utworzZrodloDeveloper(kanal);
  // Warsztat jest drugim źródłem, bo ma innych odbiorców niż Project Tree i Code Editor.
  const warsztat = utworzZrodloWarsztatu(kanal);
  const stan = utworzStanDevelopera(zrodlo, { okno });

  const obszar = document.createElement('div');
  obszar.className = 'mdev-modul';
  obszar.dataset['modul'] = 'developer';

  const edytor = utworzOknoCodeEditora(zrodlo, warsztat, stan);
  const drzewo = utworzOknoProjectTree(zrodlo, stan);
  const repozytorium = utworzOknoGitPanelu(zrodlo, stan);
  const budowanie = utworzOknoBudowania(zrodlo, warsztat, stan);
  const narzedzia = utworzOknoDevTools(warsztat, stan);
  const warsztatKodu = utworzOknoWarsztatuKodu(warsztat, stan);

  const gorny = document.createElement('div');
  gorny.className = 'mdev-modul__gora';
  gorny.append(
    oznacz(edytor.element, 'code-editor'),
    oznacz(drzewo.element, 'project-tree'),
  );

  const srodkowy = document.createElement('div');
  srodkowy.className = 'mdev-modul__srodek';
  srodkowy.append(
    oznacz(repozytorium.element, 'git-panel'),
    oznacz(budowanie.element, 'build-output'),
  );

  // Dev Tools dostaje własny pas, a nie miejsce obok Git Panelu, bo potrzebuje pełnej szerokości.
  const dolny = document.createElement('div');
  dolny.className = 'mdev-modul__dol';
  dolny.append(
    oznacz(warsztatKodu.element, 'warsztat-kodu'),
    oznacz(narzedzia.element, 'dev-tools'),
  );

  const pomocnicze = utworzPasPomocniczych({
    kanal,
    modul: KOD_MODULU,
    nazwaModulu: NAZWA_MODULU,
    okno,
    przedrostek: 'mdev',
  });

  obszar.append(gorny, srodkowy, dolny, pomocnicze.element);
  gospodarz.replaceChildren(obszar);

  function odswiez(): void {
    edytor.odswiez();
    drzewo.odswiez();
    repozytorium.odswiez();
    budowanie.odswiez();
    narzedzia.odswiez();
    warsztatKodu.odswiez();
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
      warsztatKodu.zamknij();
      narzedzia.zamknij();
      budowanie.zamknij();
      repozytorium.zamknij();
      drzewo.zamknij();
      edytor.zamknij();
      stan.zamknij();
    },
  };
}

/** Znakuje okno kodem rejestru okien operacyjnych, po którym skacze nawigacja między oknami modułu Developer. */
function oznacz(element: HTMLElement, kod: string): HTMLElement {
  element.dataset['okno'] = kod;
  element.tabIndex = -1;
  return element;
}

/**
 * Samoopisujący się moduł dla rejestru powłoki.
 *
 * Kod siedzi W MODULE, nie w mapie po stronie powłoki: dodanie modułu to jeden
 * wpis, a nie dwa, więc nie da się dodać modułu i zapomnieć o wytwórni.
 */
const KOD_MODULU = 'developer';

export const MODUL: OpisModulu = {
  kod: KOD_MODULU,
  utworzWidok: (kanal) =>
    widokZOknaSesji(
      kanal,
      (gospodarz, kanalOkna, okno) => zamontujDeveloper(gospodarz, kanalOkna, okno),
      KOD_MODULU,
    ),
};
