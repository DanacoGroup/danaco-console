import { poleWyboru, przycisk } from '../../modele/kontrolki-formularza';
import type { StanBiblioteki, TrybWyszukiwania, WidokWykazu } from './stan-biblioteki';
import { WIDOKI } from './widoki-wykazu';

/**
 * Pasek narzędzi kontekstowych Library Explorera: przełącznik widoku, tryb
 * wyszukiwania i zdjęcie zawężenia wykazu.
 *
 * Trzy kontrolki stoją razem, bo wszystkie odpowiadają na jedno pytanie — „co
 * i w jakiej postaci widzę" — i żadna z nich nie zmienia zawartości
 * repozytorium. Przełącznik widoku i zawężenie są w całości kliencke; tryb
 * wyszukiwania rozstrzyga, którą komendą pójdzie następne szukanie
 * (`library.file.search`, `knowledge.search` albo obie).
 *
 * Widok mapy zostaje w przełączniku mimo braku źródła współrzędnych: pozycja
 * usunięta wyglądałaby na widok, którego dokumentacja nigdy nie przewidywała,
 * a wybrana mówi wprost, czego kontrakt nie niesie (`widoki-wykazu.ts`).
 */
export interface NarzedziaExplorera {
  element: HTMLElement;
  /** Przerysowuje stan przełączników i pasek zawężenia. */
  odswiez(): void;
}

/** Pozycje trybu wyszukiwania wraz z komendą, którą każdy z nich idzie. */
const TRYBY: ReadonlyArray<{ kod: TrybWyszukiwania; etykieta: string }> = [
  { kod: 'pelnotekstowy', etykieta: 'Pełnotekstowy (library.file.search)' },
  { kod: 'semantyczny', etykieta: 'Semantyczny (knowledge.search)' },
  { kod: 'hybrydowy', etykieta: 'Hybrydowy (obie komendy, złożenie w oknie)' },
];

export function utworzNarzedziaExplorera(stan: StanBiblioteki): NarzedziaExplorera {
  const przelacznik = document.createElement('div');
  przelacznik.className = 'ml-widoki';
  przelacznik.setAttribute('role', 'group');
  przelacznik.setAttribute('aria-label', 'Widok wykazu plików');

  const przyciski = new Map<WidokWykazu, HTMLButtonElement>();
  for (const pozycja of WIDOKI) {
    const kontrolka = przycisk(pozycja.nazwa, 'dn-btn dn-btn--sm dn-btn--duch');
    kontrolka.dataset['widok'] = pozycja.kod;
    kontrolka.addEventListener('click', () => stan.ustawWidok(pozycja.kod));
    przyciski.set(pozycja.kod, kontrolka);
    przelacznik.append(kontrolka);
  }

  const tryb = poleWyboru(
    {
      etykieta: 'Tryb wyszukiwania',
      opis:
        'Słowa dopasowuje indeks pełnotekstowy repozytorium, znaczenie — wskaźnik osadzeń. ' +
        'Wskaźnik nie buduje się przy wgraniu pliku; przelicza go zakładka Higiena.',
    },
    TRYBY.map((pozycja) => ({ wartosc: pozycja.kod, etykieta: pozycja.etykieta })),
  );
  tryb.kontrolka.dataset['ster'] = 'tryb-wyszukiwania';
  tryb.kontrolka.addEventListener('change', () => {
    stan.ustawTryb(odczytajTryb(tryb.kontrolka.value));
  });

  const zdanieZawezenia = document.createElement('span');
  zdanieZawezenia.className = 'ml-zawezenie__zdanie';

  const zdejmij = przycisk('Zdejmij zawężenie', 'dn-btn dn-btn--sm dn-btn--zarys');
  zdejmij.dataset['czynnosc'] = 'zawezenie-zdejmij';
  zdejmij.addEventListener('click', () => stan.ustawZawezenie(null));

  const zawezenie = document.createElement('div');
  zawezenie.className = 'ml-zawezenie';
  zawezenie.hidden = true;
  zawezenie.append(zdanieZawezenia, zdejmij);

  const element = document.createElement('div');
  element.className = 'ml-narzedzia';
  element.append(przelacznik, tryb.element, zawezenie);

  return {
    element,

    odswiez() {
      for (const [kod, kontrolka] of przyciski) {
        const czynny = stan.widok() === kod;
        kontrolka.setAttribute('aria-pressed', String(czynny));
        kontrolka.className = czynny
          ? 'dn-btn dn-btn--sm dn-btn--wybrany'
          : 'dn-btn dn-btn--sm dn-btn--duch';
      }
      tryb.kontrolka.value = stan.tryb();
      const wskazane = stan.zawezenie();
      zawezenie.hidden = wskazane === null;
      if (wskazane !== null) {
        // Liczba widocznych bierze się z wykazu, nie z długości zbioru
        // wskazanego: raport potrafi wskazać plik, którego świeży odczyt już
        // nie zawiera, a katalog struktury zawęża wynik dodatkowo.
        zdanieZawezenia.textContent =
          `Wykaz zawężony: ${wskazane.opis} — widocznych pozycji ${stan.widoczne().length} ` +
          `z ${wskazane.kody.length} wskazanych.`;
      }
    },
  };
}

/** Przekład wartości kontrolki na tryb; wartość spoza wykazu bierze pełnotekstowy. */
function odczytajTryb(wartosc: string): TrybWyszukiwania {
  const pozycja = TRYBY.find((wpis) => wpis.kod === wartosc);
  return pozycja === undefined ? 'pelnotekstowy' : pozycja.kod;
}
