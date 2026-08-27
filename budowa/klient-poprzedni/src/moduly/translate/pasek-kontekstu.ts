import { przycisk } from '../../modele/kontrolki-formularza';
import { BRAK_OKNA } from './etykiety-translate';
import type { StanTranslate } from './stan-translate';

/**
 * Pas kontekstu modułu pokazujący okno wykonania odczytane z rdzenia, wraz
 * z przyciskiem ponowienia odczytu po niepowodzeniu.
 */
export interface PasekKontekstu {
  element: HTMLElement;
  odswiez(): void;
}

export function utworzPasekKontekstu(
  stan: StanTranslate,
  naPonowienie: () => void,
): PasekKontekstu {
  const zdanie = document.createElement('p');
  zdanie.className = 'mt-kontekst__zdanie';

  const ponow = przycisk('Spróbuj ponownie', 'dn-btn dn-btn--sm dn-btn--zarys');
  ponow.classList.add('mt-kontekst__ponow');
  ponow.addEventListener('click', naPonowienie);

  const dane = document.createElement('dl');
  dane.className = 'mt-kontekst__dane';

  const element = document.createElement('header');
  element.className = 'mt-kontekst';
  element.setAttribute('aria-label', 'Kontekst okna modułu Translate');
  element.append(zdanie, ponow, dane);

  function odswiez(): void {
    const faza = stan.fazaKontekstu();
    element.dataset['faza'] = faza;
    // Widocznością steruje klasa, bo przycisk biblioteki narzuca własny styl wyświetlania.
    ponow.classList.toggle('mt-kontekst__ponow--widoczny', faza === 'blad');
    if (faza === 'spoczynek') {
      pokazZdanie('Kontekst okna nieodczytany — moduł pyta o niego przy wejściu.');
      return;
    }
    if (faza === 'odczyt') {
      pokazZdanie('Odczyt okna modułu z rdzenia…');
      return;
    }
    if (faza === 'blad') {
      pokazZdanie(`Odczyt okna modułu nie powiódł się: ${stan.powodKontekstu()}`);
      return;
    }
    const okno = stan.oknoModulu();
    if (okno === null) {
      pokazZdanie(BRAK_OKNA);
      return;
    }
    // Pas mówi wyłącznie, w imieniu którego okna moduł działa.
    zdanie.textContent = `Okno ${okno.title ?? okno.id}`;
    dane.replaceChildren();
  }

  function pokazZdanie(tresc: string): void {
    zdanie.textContent = tresc;
    dane.replaceChildren();
  }

  odswiez();
  return { element, odswiez };
}

