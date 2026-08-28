import { zdanieKwitu, type Kwit, type MagazynKwitow } from './kwit-decyzji';

/**
 * Pasek kwitu — wykaz decyzji wysłanych z tego urządzenia, umieszczony u dołu
 * ekranu w zasięgu kciuka. Wiersz niesie zdanie rdzenia: godzinę, treść decyzji
 * i potwierdzony stan. Kwit nieudany trafia na ten sam pasek wraz z treścią
 * odmowy.
 */
export interface PasekKwitu {
  element: HTMLElement;
  /** Dopisuje kwit do pamięci i wypisuje go na wierzchu. */
  dopisz(kwit: Kwit): void;
  /** Wypisuje ostatnie kwity z pamięci — wołane przy otwarciu ekranu. */
  odswiez(): void;
}

/**
 * Ile kwitów widać naraz na pasku; wykaz wypisuje tylko tyle najnowszych
 * wierszy, a pozostałe zostają w pamięci trwałej urządzenia.
 */
const ILE_WIDOCZNYCH = 3;

export function utworzPasekKwitu(magazyn: MagazynKwitow): PasekKwitu {
  const wykaz = document.createElement('ul');
  wykaz.className = 'mb-kwity__wykaz';

  const naglowek = document.createElement('p');
  naglowek.className = 'mb-kwity__naglowek';
  naglowek.textContent = 'Co stąd wysłano';

  const element = document.createElement('section');
  element.className = 'mb-kwity';
  element.setAttribute('role', 'status');
  element.setAttribute('aria-live', 'polite');
  element.append(naglowek, wykaz);

  function wypisz(kwity: readonly Kwit[]): void {
    if (kwity.length === 0) {
      const pusto = document.createElement('li');
      pusto.className = 'mb-kwity__pusto';
      pusto.textContent = 'Z tego urządzenia nie wysłano jeszcze żadnej decyzji.';
      wykaz.replaceChildren(pusto);
      return;
    }

    wykaz.replaceChildren(
      ...kwity.slice(0, ILE_WIDOCZNYCH).map((kwit) => {
        const wiersz = document.createElement('li');
        wiersz.className = 'mb-kwity__wiersz';
        wiersz.dataset['udany'] = String(kwit.udany);
        wiersz.textContent = zdanieKwitu(kwit);
        return wiersz;
      }),
    );
  }

  return {
    element,
    dopisz: (kwit) => wypisz(magazyn.dopisz(kwit)),
    odswiez: () => wypisz(magazyn.odczytaj()),
  };
}
