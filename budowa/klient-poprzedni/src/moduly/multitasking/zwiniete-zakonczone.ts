import { odmienAgentow, zapiszCzasOdcinka } from './format-zadan';
import { ODZNAKA, type Przeplyw } from './zadania-w-tle';

/**
 * Zwinięty wykaz przepływów domkniętych, zapowiedziany licznikiem ich liczby.
 * Pusty wykaz nie tworzy elementu, a element `details` otwiera się i zamyka bez
 * potwierdzenia, więc brak zakończeń daje krótszy panel, nie kontrolkę martwą.
 */
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

/**
 * Pozycja zwiniętego wykazu: nazwa przepływu, jego stan oraz to samo
 * podsumowanie czasu odcinka i liczby podagentów, które niesie karta przepływu
 * w widoku rozwiniętym.
 */
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
