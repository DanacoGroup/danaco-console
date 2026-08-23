import { przycisk } from '../../modele/kontrolki-formularza';
import type { WynikOdkrycia } from './wynik-odkrycia';

/**
 * Jedna pozycja wykazu wyników Discovery Panel.
 *
 * Jedna odpowiedzialność: przełożenie pozycji wyniku na wiersz z metadanymi,
 * fragmentem i akcją podstawową „→ Dodaj do źródeł". Przycisk jest czynny
 * zawsze, także dla pozycji już przeniesionej: powtórzenie kończy się
 * odpowiedzią rdzenia, nie odebraniem klikalności.
 *
 * Wiersz nie niesie pola wyboru. Zaznaczenie wielokrotne w module dotyczy
 * źródeł i ustaleń — bytów badania — a pozycja wyniku badaniem jeszcze nie
 * jest; staje się nim dopiero po skatalogowaniu.
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

/** Fragment treści oddany przez rdzeń — cytat, więc bez skracania i bez zmian. */
function fragment(tresc: string): HTMLElement {
  const element = document.createElement('p');
  element.className = 'mr-wykaz__fragment';
  element.textContent = tresc;
  return element;
}
