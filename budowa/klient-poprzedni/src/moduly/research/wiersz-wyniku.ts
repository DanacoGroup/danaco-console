import { przycisk } from '../../modele/kontrolki-formularza';
import type { WynikOdkrycia } from './wynik-odkrycia';

/**
 * Jedna pozycja wykazu wyników Discovery Panel: wynik przełożony na wiersz z metadanymi,
 * fragmentem i akcją dodania do źródeł.
 */
export function utworzWierszWyniku(
  pozycja: WynikOdkrycia,
  naPrzeniesienie: (pozycja: WynikOdkrycia) => void,
): HTMLElement {
  const element = document.createElement('li');
  element.className = 'mr-wykaz__wiersz';
  element.dataset['wynik'] = pozycja.klucz;

  const tytul = document.createElement('span');
  tytul.className = 'mr-wykaz__tytul';
  tytul.textContent = pozycja.tytul;

  const rodzaj = document.createElement('span');
  rodzaj.className = 'dn-plakietka dn-plakietka--informacja';
  rodzaj.textContent = `wejdzie jako: ${pozycja.rodzaj}`;

  const dodaj = przycisk('→ Dodaj do źródeł', 'dn-btn dn-btn--sm dn-btn--duch');
  dodaj.addEventListener('click', () => naPrzeniesienie(pozycja));

  const metadane = document.createElement('span');
  metadane.className = 'mr-wykaz__meta';
  metadane.textContent = pozycja.metadane;

  element.append(tytul, rodzaj, dodaj, metadane);
  if (pozycja.fragment !== '') element.append(fragment(pozycja.fragment));
  return element;
}

/** Fragment treści oddany przez rdzeń jako cytat wyniku wyszukiwania, wypisany bez skracania i bez zmian jego treści. */
function fragment(tresc: string): HTMLElement {
  const element = document.createElement('p');
  element.className = 'mr-wykaz__fragment';
  element.textContent = tresc;
  return element;
}
