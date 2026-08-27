import type { Action } from '../../../shared/contract';
import type { ProfilModulu } from './profil-modulu';

/** Panel akcji modułu pokazuje pozycje pobrane z rejestru rdzenia, uporządkowane w kolejności, w jakiej rejestr je zwraca. */
export interface PanelAkcji {
  /** Element montowany w oknie komunikacji. */
  element: HTMLElement;
  /** Buduje panel na nowo dla wskazanego modułu i jego katalogu akcji. */
  ustaw(profil: ProfilModulu, akcje: readonly Action[], powod: string): void;
}

/**
 * Panel akcji budowany dynamicznie z rejestru akcji na moduł nie zna ani jednej akcji z nazwy: kliknięcie wstawia do pola wypowiedzi polecenie wykonania wraz z nazwą komendy kontraktu, a treść żądania buduje model, nie panel.
 */
export function utworzPanelAkcji(naAkcje: (akcja: Action) => void): PanelAkcji {
  const element = document.createElement('section');
  element.className = 'dc-panel-akcji';
  element.setAttribute('aria-label', 'Panel akcji modułu');

  const naglowek = document.createElement('h2');
  naglowek.className = 'dc-panel-akcji__naglowek';

  const lista = document.createElement('div');
  lista.className = 'dc-panel-akcji__lista';

  const komunikat = document.createElement('p');
  komunikat.className = 'dc-panel-akcji__komunikat';

  element.append(naglowek, lista, komunikat);

  return {
    element,
    ustaw(profil, akcje, powod) {
      element.dataset['modul'] = profil.kod;
      naglowek.textContent = `Panel akcji — ${profil.nazwa}`;
      lista.replaceChildren(...akcje.map((a) => przycisk(a, naAkcje)));
      komunikat.textContent = powod;
      komunikat.hidden = powod.length === 0;
    },
  };
}

/** Pozycja panelu akcji niosąca etykietę katalogu, opis działania oraz komendę wstawianą do pola podpowiedzi. */
function przycisk(akcja: Action, naAkcje: (akcja: Action) => void): HTMLButtonElement {
  const element = document.createElement('button');
  element.type = 'button';
  element.className = 'dn-btn dn-btn--duch dn-btn--sm dc-panel-akcji__akcja';
  element.dataset['akcja'] = akcja.id;
  element.dataset['komenda'] = akcja.command;
  const opis = akcja.description ?? akcja.command;
  element.title = `${opis} · komenda ${akcja.command}`;
  element.setAttribute('aria-label', `${akcja.name} — ${opis}`);
  element.textContent = akcja.name;
  element.addEventListener('click', () => naAkcje(akcja));
  return element;
}
