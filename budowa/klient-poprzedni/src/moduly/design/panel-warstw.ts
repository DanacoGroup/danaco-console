import type { DesignBoardLayer } from '../../../../shared/contract';
import type { StanKompozycji } from './stan-kompozycji';

/**
 * Panel warstw Design Board: wykaz warstw uporządkowany malejąco po polu
 * porządku, z zaznaczeniem, adnotacją, plakietką pochodzenia, przestawieniem
 * blokady i zdjęciem warstwy z kanwy. Pusty wykaz zastępuje zdanie o braku
 * zasobów.
 */
export interface PanelWarstw {
  element: HTMLElement;
  odswiez(): void;
}

export function utworzPanelWarstw(stan: StanKompozycji): PanelWarstw {
  const tytul = document.createElement('h4');
  tytul.className = 'md-panel__tytul';
  tytul.textContent = 'Warstwy kompozycji';

  const wykaz = document.createElement('ul');
  wykaz.className = 'md-warstwy';

  const pustka = document.createElement('p');
  pustka.className = 'dn-pusty-stan-opis md-panel__opis';
  pustka.textContent = 'Kompozycja jest bez zestawionych jeszcze zasobów.';

  const element = document.createElement('section');
  element.className = 'md-panel md-panel-warstw';
  element.append(tytul, pustka, wykaz);

  return {
    element,

    odswiez() {
      const warstwy = stan.warstwy();
      pustka.hidden = warstwy.length > 0;
      wykaz.replaceChildren(
        ...[...warstwy]
          .sort((a, b) => (b.order ?? 0) - (a.order ?? 0))
          .map((warstwa) => wierszWarstwy(warstwa, stan)),
      );
    },
  };
}

function wierszWarstwy(warstwa: DesignBoardLayer, stan: StanKompozycji): HTMLElement {
  const zaznaczenie = document.createElement('input');
  zaznaczenie.type = 'checkbox';
  zaznaczenie.className = 'dn-check';
  zaznaczenie.checked = stan.zaznaczone().includes(warstwa.id);
  zaznaczenie.setAttribute('aria-label', `Zaznacz warstwę ${warstwa.id}`);
  zaznaczenie.addEventListener('change', () => stan.zaznacz(warstwa.id, true));

  const adnotacja = document.createElement('input');
  adnotacja.type = 'text';
  adnotacja.className = 'dn-pole-kontrolka md-warstwy__adnotacja';
  adnotacja.value = warstwa.note ?? '';
  adnotacja.placeholder = 'adnotacja warstwy';
  adnotacja.setAttribute('aria-label', `Adnotacja warstwy ${warstwa.id}`);
  adnotacja.addEventListener('change', () => stan.ustawAdnotacje(warstwa.id, adnotacja.value));

  const zrodlo = document.createElement('span');
  zrodlo.className = 'dn-plakietka';
  zrodlo.textContent = warstwa.assetId === undefined ? 'element pomocniczy' : 'zasób';

  const blokada = document.createElement('button');
  blokada.type = 'button';
  blokada.className = 'dn-btn dn-btn--duch dn-btn--sm';
  blokada.textContent = warstwa.locked === true ? 'Odblokuj' : 'Zablokuj';
  blokada.addEventListener('click', () => stan.przestawBlokade(warstwa.id));

  const usun = document.createElement('button');
  usun.type = 'button';
  usun.className = 'dn-btn dn-btn--duch dn-btn--sm';
  usun.textContent = 'Zdejmij z kanwy';
  usun.addEventListener('click', () => stan.usun(warstwa.id));

  const element = document.createElement('li');
  element.className = 'md-warstwy__wiersz';
  element.dataset['warstwa'] = warstwa.id;
  element.dataset['zablokowana'] = String(warstwa.locked === true);
  element.append(zaznaczenie, adnotacja, zrodlo, blokada, usun);
  return element;
}
