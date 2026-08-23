import { przycisk } from '../../modele/kontrolki-formularza';
import { BRAK_OKNA } from './etykiety-translate';
import type { StanTranslate } from './stan-translate';

/**
 * Pas kontekstu modułu — okno wykonania odczytane komendą `window.state.get`.
 *
 * Trzy pustki są trzema różnymi zdaniami: „jeszcze nie pytałem", „pytam"
 * i „rdzeń nie zna ani jednego okna tej sesji". Pierwsze każe czekać na wejście
 * do modułu, drugie na odpowiedź, trzecie mówi, że komendy Source Panel nie
 * mają czym zaadresować żądania.
 *
 * „Spróbuj ponownie" stoi tylko tutaj, bo kontekst okna jest jedynym odczytem
 * modułu; resztę wyzwala zapis, którego przycisk zostaje klikalny także po
 * odmowie. Ponowienie wyzwala odczyt i nie zdejmuje komunikatu — ten znika
 * dopiero, gdy odczyt się powiedzie.
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
    // Klasa steruje widocznością zamiast atrybutu `hidden`: przycisk
    // biblioteki niesie własny `display`, który `[hidden]` przegrywa.
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
    // Pas mówi, w imieniu którego okna moduł działa, i nic ponadto. Tryb
    // uprawnień, zasięg wykonania, rola okna i katalogi robocze stoją w panelu
    // „Sterowanie okna" na tym samym ekranie.
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

