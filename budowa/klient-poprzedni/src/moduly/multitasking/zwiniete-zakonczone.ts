import { odmienAgentow, zapiszCzasOdcinka } from './format-zadan';
import { ODZNAKA, type Przeplyw } from './zadania-w-tle';

/**
 * Zwinięty licznik przepływów domkniętych — „Zakończone 48".
 *
 * Przepływów zakończonych bywa kilkadziesiąt i wypisane w całości spychałyby
 * poza krawędź panelu pracę trwającą. Zwinięcie nie jest ukryciem: licznik
 * podaje ich liczbę, a rozwinięcie oddaje każdy po nazwie i stanie.
 *
 * Element `details` otwiera się i zamyka bez potwierdzenia; brak przepływów
 * zakończonych daje krótszy panel, a nie kontrolkę nieczynną.
 */

/** Zwinięty wykaz przepływów domkniętych; pusty wykaz nie tworzy elementu. */
export function zwinieteZakonczone(przeplywy: readonly Przeplyw[]): HTMLElement | null {
  if (przeplywy.length === 0) return null;

  const element = document.createElement('details');
  element.className = 'dm-zwiniete';

  const naglowek = document.createElement('summary');
  naglowek.className = 'dm-zwiniete__naglowek';
  naglowek.textContent = `Zakończone ${przeplywy.length}`;
  element.append(naglowek);

  const lista = document.createElement('ul');
  lista.className = 'dm-zwiniete__lista';
  for (const przeplyw of przeplywy) lista.append(pozycja(przeplyw));
  element.append(lista);
  return element;
}

/** Pozycja zwiniętego wykazu: nazwa, stan i to samo podsumowanie co na karcie. */
function pozycja(przeplyw: Przeplyw): HTMLElement {
  const nazwa = document.createElement('strong');
  nazwa.className = 'dm-zwiniete__nazwa';
  nazwa.textContent = przeplyw.nazwa === '' ? 'zadanie bez treści w kontrakcie' : przeplyw.nazwa;
  nazwa.title = przeplyw.nazwa;

  const opis = document.createElement('span');
  opis.className = 'dm-zwiniete__opis';
  opis.textContent = `${ODZNAKA[przeplyw.stan]} · ${zapiszCzasOdcinka(przeplyw.czas)} · ${odmienAgentow(przeplyw.podagenci.length)}`;

  const element = document.createElement('li');
  element.className = 'dm-zwiniete__pozycja';
  element.dataset['stan'] = przeplyw.stan;
  element.append(nazwa, opis);
  return element;
}
