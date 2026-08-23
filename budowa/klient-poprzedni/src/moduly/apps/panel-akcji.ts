import type { Action } from '../../../../shared/contract';
import { pokazKomunikat } from '../../aplikacja/komunikaty';
import { zrodloAkcjiKanalu, type ZrodloAkcjiModulu } from '../../okno-komunikacji/katalog-akcji';
import type { Kanal } from '../../protokol/kanal';
import { KOD_MODULU } from './zrodlo-okna-modulu';

/**
 * Panel akcji modułu zbudowany z rejestru rdzenia, nie z listy w kliencie.
 *
 * Pozycje przychodzą komendą `action.list` o zasięgu `module` i kluczu równym
 * kodowi modułu. Nowa akcja modułu to nowy wiersz rejestru, nie zmiana kodu
 * klienta — dlatego panel nie ma ani jednej pozycji zapisanej tutaj.
 *
 * Wiersz rejestru niesie nazwę komendy, którą akcja wywołuje, ale generycznej
 * drogi wywołania (`window.action`) nie woła dziś żaden widok. Naciśnięcie
 * nazywa więc komendę wiersza i mówi, że droga generyczna czeka na konsumenta,
 * zamiast milczeć albo udawać wykonanie.
 */
export interface PanelAkcji {
  element: HTMLElement;
  /** Odczytuje katalog akcji modułu z rejestru rdzenia. */
  wczytaj(): void;
}

export function utworzPanelAkcji(kanal: Kanal): PanelAkcji {
  const zrodlo: ZrodloAkcjiModulu = zrodloAkcjiKanalu(kanal);

  const tytul = document.createElement('p');
  tytul.className = 'mp-akcje__tytul';
  tytul.textContent = 'Panel akcji modułu (rejestr rdzenia)';

  const lista = document.createElement('ul');
  lista.className = 'mp-akcje';

  const powod = document.createElement('p');
  powod.className = 'mp-akcje__powod';

  const element = document.createElement('div');
  element.className = 'mp-akcje__powloka';
  element.append(tytul, lista, powod);

  function nanies(akcje: readonly Action[], zdanie: string): void {
    lista.replaceChildren(...akcje.map(wiersz));
    powod.textContent = zdanie;
    powod.hidden = zdanie === '';
  }

  nanies([], 'Katalog akcji jeszcze nieodczytany.');

  return {
    element,
    wczytaj() {
      nanies([], 'Czytam katalog akcji modułu…');
      zrodlo(KOD_MODULU, (akcje, zdanie) => nanies(akcje, zdanie));
    },
  };
}

/** Jedna pozycja katalogu: przycisk nazywający komendę wiersza rejestru. */
function wiersz(akcja: Action): HTMLElement {
  const przycisk = document.createElement('button');
  przycisk.type = 'button';
  przycisk.className = 'dn-btn dn-btn--zarys dn-btn--sm';
  przycisk.textContent = akcja.name;
  przycisk.dataset['akcja'] = akcja.id;
  przycisk.addEventListener('click', () => {
    pokazKomunikat({
      tytul: akcja.name,
      tresc:
        `Wiersz rejestru wskazuje komendę ${akcja.command}. Generycznej drogi ` +
        'wywołania akcji (window.action) nie woła dziś żaden widok — moduł nie ' +
        'udaje wykonania.',
      waga: 'ostrz',
    });
  });

  const element = document.createElement('li');
  element.className = 'mp-akcje__wiersz';
  element.append(przycisk);
  if (akcja.description !== undefined) {
    const opis = document.createElement('span');
    opis.className = 'mp-akcje__opis';
    opis.textContent = akcja.description;
    element.append(opis);
  }
  return element;
}
