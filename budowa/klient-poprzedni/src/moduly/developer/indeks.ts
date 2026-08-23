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
 * repozytorium.
 *
 * Moduł pracuje w oknie, nie w sesji: wszystkie pięć komend obszaru
 * (`developer.file.open`, `file.save`, `tree.get`, `git.action`, `build.run`)
 * wymaga `windowId`. Umowa `WidokModulu.wczytaj` niesie natomiast
 * identyfikator sesji, więc złożenie jedzie przez `widokZOknaSesji`: przejście
 * pyta rdzeń o okna sesji i odracza montaż do chwili, gdy okno tego modułu
 * jest znane. Kod modułu jest przejściu podawany, bo bez niego przejście bierze
 * pierwsze okno w wykazie — także cudze (`moduly/rejestracja.ts`).
 *
 * Układ wynika z ról. W pasie górnym wiodące okno Code Editor wraz
 * z pomocniczym Project Tree: drzewo wskazuje plik (`stan.wskazPlik`), edytor
 * go odczytuje. W pasie środkowym zarządca repozytorium (Git Panel) i monitor
 * (Build Output i Run & Debug). W pasie dolnym Dev Tools — kolumna czterech
 * integracji deweloperskich. Okno rozmowy modułu nie należy do tego złożenia:
 * jest bytem sesji i składa je warstwa rozmowy, tak samo jak Execution Loop
 * Window.
 *
 * Opracowanie wymienia siedem okien modułu. Pięć składa to złożenie, dwa
 * pozostałe są wspólne platformie i stoją poza katalogiem modułu. Monitor jest
 * JEDNYM oknem o dwóch częściach (Build Output i Run & Debug) przełączanych
 * zakładkami w nagłówku kolumny, nie dwoma oknami — dlatego ma jeden kod
 * katalogu rdzenia.
 *
 * Ostatni pas niesie okna pomocnicze — zbudowane (podgląd w tle bash na
 * `terminal.output.stream`) oraz spis pozycji jeszcze nieistniejących wraz
 * z powodem każdej (`okna-pomocnicze/rejestr-pomocniczych.ts`). Terminal ma
 * w tym spisie miejsce i nie jest tu budowany drugi raz; nie powiela go też
 * zakładka Containers okna Dev Tools.
 *
 * Zdarzenie `developer.build.changed` ma dwie subskrypcje, bo każda bierze co
 * innego. Stan modułu unieważnia po nim drzewo i edytor — budowanie generuje
 * pliki. Build Output bierze przyrost logu (`logLine`), którego stan nie
 * przenosi; kontrakt nie ma komendy rozwijającej `logRef`, więc log narasta
 * wyłącznie ze zdarzenia.
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

/** Nazwa modułu w etykietach dostępności pasa okien pomocniczych. */
const NAZWA_MODULU = 'Developer';

export function zamontujDeveloper(
  gospodarz: HTMLElement,
  kanal: Kanal,
  okno: string,
): ZamontowanyDeveloper {
  const zrodlo = utworzZrodloDeveloper(kanal);
  // Warsztat jest drugim źródłem, bo ma innych odbiorców: zakładki Dev Tools,
  // panel Run & Debug, historię przebiegów i pasek operacji Code Editora.
  // Jedna umowa na wszystko kazałaby Project Tree przyjmować zależność od
  // czterdziestu metod, z których używa czterech.
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

  // Dev Tools dostaje własny pas, a nie miejsce obok Git Panelu: kolumna
  // czterech zakładek integracji jest w opracowaniu rozszerzeniem bocznym
  // obszaru roboczego, nie sąsiadem zarządcy repozytorium, i potrzebuje pełnej
  // szerokości na wykaz zależności zewnętrznych.
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
      // Pas zamyka się PIERWSZY: trzyma subskrypcję `stream.chunk`, która żyje
      // niezależnie od stanu modułu i po zejściu ze sceny nikt by jej nie zdjął.
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

/** Znakuje okno kodem rejestru okien operacyjnych — po nim skacze nawigacja. */
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
