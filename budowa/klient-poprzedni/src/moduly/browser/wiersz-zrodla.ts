import type { BrowserSource } from '../../../../shared/contract';
import { KLASA_PRZYCISKU, przyciskCzynnosci } from './przyciski-browser';

/**
 * Jedna pozycja wykazu Sources Panel wraz z jej panelem akcji. Wykaz,
 * formularz dodania i przekazanie do Research mieszkają w oknie.
 *
 * Każda akcja wiersza jest naciskalna — także „Usuń", której kontrakt nie
 * niesie: przycisk odpowiada nazwaniem brakującej komendy zamiast znikać albo
 * gasnąć. Wiersz nie usuwa pozycji z widoku na własną rękę, bo zniknięcie bez
 * zapisu w rdzeniu byłoby udawaniem wykonania.
 */
export interface AkcjeZrodla {
  otworz(zrodlo: BrowserSource): void;
  migawka(zrodlo: BrowserSource): void;
  oznaczKluczowe(zrodlo: BrowserSource): void;
  usun(zrodlo: BrowserSource): void;
  przelaczWybor(zrodlo: BrowserSource): void;
}

export function utworzWierszZrodla(
  zrodlo: BrowserSource,
  wybrane: boolean,
  akcje: AkcjeZrodla,
): HTMLElement {
  const wybor = document.createElement('input');
  wybor.type = 'checkbox';
  wybor.className = 'dn-check mb-zrodla__wybor';
  wybor.checked = wybrane;
  wybor.setAttribute('aria-label', `Zaznacz źródło ${zrodlo.url} do przekazania`);
  wybor.addEventListener('change', () => akcje.przelaczWybor(zrodlo));

  const opis = document.createElement('span');
  opis.className = 'mb-zrodla__opis';
  opis.textContent = (zrodlo.title ?? '').trim() === '' ? zrodlo.url : `${zrodlo.title} — ${zrodlo.url}`;

  const plakietka = document.createElement('span');
  plakietka.className = zrodlo.key === true ? 'dn-plakietka dn-plakietka--sygnal' : 'dn-plakietka';
  plakietka.textContent = zrodlo.key === true ? 'kluczowe' : 'zwykłe';

  const element = document.createElement('li');
  element.className = 'mb-zrodla__wiersz';
  element.dataset['zrodlo'] = zrodlo.id;
  element.append(
    wybor,
    opis,
    plakietka,
    przyciskCzynnosci('Otwórz', KLASA_PRZYCISKU.duch, () => akcje.otworz(zrodlo)),
    przyciskCzynnosci('Podgląd migawki', KLASA_PRZYCISKU.duch, () => akcje.migawka(zrodlo)),
    przyciskCzynnosci('Oznacz jako kluczowe', KLASA_PRZYCISKU.duch, () =>
      akcje.oznaczKluczowe(zrodlo),
    ),
    przyciskCzynnosci('Usuń', KLASA_PRZYCISKU.duch, () => akcje.usun(zrodlo)),
  );
  return element;
}
