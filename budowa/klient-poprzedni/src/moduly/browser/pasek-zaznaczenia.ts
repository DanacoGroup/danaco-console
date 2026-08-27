import { opisOdmowy } from '../../komponenty/odmowa';
import { pole, przycisk } from '../../modele/kontrolki-formularza';
import { skutekPytania } from './skutek-zapisu';
import type { StanPrzegladania } from './stan-przegladania';

/**
 * Pasek pływający zaznaczenia modułu Browser: pięć czynności wykonywanych na
 * zaznaczonym fragmencie strony. Pasek nie znika i nie gaśnie, a przy pustym
 * zaznaczeniu przyciski zostają naciskalne i nazywają, czego brakuje.
 */
export interface PasekZaznaczenia {
  element: HTMLElement;
  /** Nanosi bieżące zaznaczenie na pasek. */
  odswiez(): void;
}

/**
 * Ujścia, których pasek potrzebuje z zewnątrz: przeniesienie fragmentu do
 * formularza notatki, przekazanie go do tłumaczenia oraz zdanie odpowiedzi
 * pokazywane po każdym naciśnięciu przycisku paska.
 */
export interface UjsciaZaznaczenia {
  /** Przenosi fragment do formularza notatki (Notes Panel). */
  naNotatke(fragment: string): void;
  /** Przekazuje zaznaczony fragment do modułu Translate — czynność strony. */
  naTlumaczenie(): Promise<void>;
  /** Odpowiedź pokazywana Operatorowi po każdym naciśnięciu. */
  powiedz(tresc: string, powodzenie: boolean): void;
}

export function utworzPasekZaznaczenia(
  stan: StanPrzegladania,
  ujscia: UjsciaZaznaczenia,
): PasekZaznaczenia {
  const podpis = document.createElement('span');
  podpis.className = 'mb-zaznaczenie__podpis';

  const wyjasnij = przycisk('Wyjaśnij', 'dn-btn dn-btn--sm dn-btn--zarys');
  const wyodrebnij = przycisk('Wyodrębnij', 'dn-btn dn-btn--sm dn-btn--zarys');
  const notatka = przycisk('Notatka', 'dn-btn dn-btn--sm dn-btn--zarys');
  const tlumacz = przycisk('Tłumacz', 'dn-btn dn-btn--sm dn-btn--zarys');
  const pytanie = pole('Pytanie o zaznaczony fragment', 'o co zapytać?');
  const zapytaj = przycisk('Zapytaj', 'dn-btn dn-btn--sm dn-btn--zarys');

  const element = document.createElement('div');
  element.className = 'mb-zaznaczenie';
  element.dataset['zaznaczenie'] = 'nie';
  element.setAttribute('aria-label', 'Pasek zaznaczonego fragmentu strony');
  element.append(podpis, wyjasnij, wyodrebnij, notatka, tlumacz, pytanie, zapytaj);

  /** Fragment albo pusty napis; jedno miejsce sprawdzenia dla trzech akcji. */
  function fragment(): string {
    return stan.zaznaczenie();
  }

  function brakZaznaczenia(): boolean {
    if (fragment() !== '') return false;
    ujscia.powiedz(
      'Zaznacz fragment w podglądzie strony — pasek pracuje na zaznaczeniu, nie na całej stronie.',
      false,
    );
    return true;
  }

  // Pyta o zaznaczony fragment; wskazane polecenie jest zdaniem otwierajacym.
  async function zapytajOFragment(polecenie: string): Promise<void> {
    if (brakZaznaczenia()) return;
    const idOkna = stan.idOkna();
    if (idOkna === '') {
      ujscia.powiedz(stan.powod(), false);
      return;
    }
    ujscia.powiedz('Pytanie o zaznaczony fragment idzie do okna rozmowy…', true);
    const wynik = await stan.zapisy.zapytaj(idOkna, `${polecenie}\n\n${fragment()}`);
    if (!wynik.udany || wynik.wynik === undefined) {
      ujscia.powiedz(opisOdmowy('Pytanie o fragment', wynik.blad?.code, wynik.blad?.message), false);
      return;
    }
    const skutek = skutekPytania(wynik.wynik.message, idOkna);
    ujscia.powiedz(skutek.zdanie, skutek.udany);
  }

  wyjasnij.addEventListener('click', () => {
    void zapytajOFragment('Wyjaśnij zaznaczony fragment oglądanej strony:');
  });

  zapytaj.addEventListener('click', () => {
    const tresc = pytanie.value.trim();
    if (tresc === '') {
      ujscia.powiedz('Wpisz pytanie — „Zapytaj" wysyła je razem z zaznaczeniem.', false);
      return;
    }
    void zapytajOFragment(tresc);
  });

  tlumacz.addEventListener('click', () => void ujscia.naTlumaczenie());

  wyodrebnij.addEventListener('click', () => {
    if (brakZaznaczenia()) return;
    // Zdanie bierze sie ze skutku dopisania, a nie z samego nacisniecia.
    const dopisany = stan.zebrane.dopiszWyodrebniony(fragment());
    ujscia.powiedz(
      dopisany
        ? 'Fragment dopisany do wyodrębnionej treści okna.'
        : 'Ten fragment jest już w wyodrębnionej treści okna — wykaz się nie zmienił.',
      dopisany,
    );
  });

  notatka.addEventListener('click', () => {
    if (brakZaznaczenia()) return;
    ujscia.naNotatke(fragment());
    ujscia.powiedz('Fragment przeniesiony do formularza notatki w Notes Panel.', true);
  });

  return {
    element,

    odswiez() {
      const tresc = fragment();
      element.dataset['zaznaczenie'] = tresc === '' ? 'nie' : 'tak';
      podpis.textContent =
        tresc === ''
          ? 'Brak zaznaczenia — przyciski nazwą, czego brakuje.'
          : `Zaznaczono ${tresc.length} znaków.`;
    },
  };
}
