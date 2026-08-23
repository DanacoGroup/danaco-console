import type { ConfigEntry } from '../../../shared/contract';
import { adresWpisu, opisAdresu } from './adres-ustawienia';
import { napis } from './kontrolka';
import { czyZrodlo, type Rozstrzygniecie } from './rozstrzygniecie';
import { nazwaZasiegu } from './zasiegi';

/**
 * Podgląd polityki efektywnej z dziedziczeniem.
 *
 * Wykaz wszystkich zapisów jednego klucza, od najwęższego. Wiersz, z którego
 * pochodzi wartość obowiązująca w punkcie widzenia okna, jest oznaczony
 * sygnałem; pozostałe zapisy pozostają widoczne, bo Operator ma wiedzieć nie
 * tylko, co obowiązuje, ale i co czeka na innym poziomie.
 *
 * Ostatni wiersz należy do wartości domyślnej katalogu. Brak zapisu nie jest
 * dziurą — jest wartością domyślną i tak został nazwany.
 *
 * Każdy zapis wolno usunąć. Usunięcie nie jest niszczeniem ustawienia, tylko
 * zdjęciem nadpisania: wartość wraca do poziomu szerszego.
 */
export interface LancuchZasiegow {
  /** Wykaz osadzany w panelu zasięgu pola. */
  element: HTMLElement;
  /** Przerysowuje wykaz na podstawie świeżego rozstrzygnięcia. */
  odswiez(rozstrzygniecie: Rozstrzygniecie): void;
}

export function utworzLancuchZasiegow(
  naPrzywrocenie: (wpis: ConfigEntry) => void,
): LancuchZasiegow {
  const element = document.createElement('ul');
  element.className = 'dk-lancuch';

  return {
    element,

    odswiez(rozstrzygniecie) {
      const wiersze = rozstrzygniecie.lancuch.map((wpis) =>
        wierszWpisu(wpis, czyZrodlo(rozstrzygniecie, wpis), naPrzywrocenie),
      );
      element.replaceChildren(...wiersze, wierszDomyslny(rozstrzygniecie));
    },
  };
}

/** Wiersz jednego zapisu konfiguracji. */
function wierszWpisu(
  wpis: ConfigEntry,
  obowiazuje: boolean,
  naPrzywrocenie: (wpis: ConfigEntry) => void,
): HTMLElement {
  const wiersz = document.createElement('li');
  wiersz.className = 'dk-lancuch__wiersz';
  if (obowiazuje) wiersz.dataset.obowiazuje = 'tak';

  const adres = document.createElement('span');
  adres.className = 'dk-lancuch__adres';
  adres.textContent = opisAdresu(adresWpisu(wpis));

  const wartosc = document.createElement('code');
  wartosc.className = 'dk-lancuch__wartosc';
  wartosc.textContent = skroc(napis(wpis.value));

  const znak = document.createElement('span');
  znak.className = obowiazuje
    ? 'dn-plakietka dn-plakietka--sygnal dk-lancuch__znak'
    : 'dn-plakietka dk-lancuch__znak';
  znak.textContent = obowiazuje ? 'obowiązuje' : 'zapis';

  const usun = document.createElement('button');
  usun.type = 'button';
  usun.className = 'dn-btn dn-btn--duch dn-btn--sm dk-lancuch__usun';
  usun.textContent = 'Zdejmij zapis';
  usun.title = `Usuwa zapis z poziomu „${nazwaZasiegu(wpis.scope)}"; wartość wróci z poziomu szerszego`;
  usun.addEventListener('click', () => naPrzywrocenie(wpis));

  wiersz.append(adres, wartosc, znak, usun);
  return wiersz;
}

/** Wiersz wartości domyślnej katalogu — zawsze ostatni, nigdy usuwalny. */
function wierszDomyslny(rozstrzygniecie: Rozstrzygniecie): HTMLElement {
  const wiersz = document.createElement('li');
  wiersz.className = 'dk-lancuch__wiersz dk-lancuch__wiersz--domyslna';
  if (rozstrzygniecie.domyslna) wiersz.dataset.obowiazuje = 'tak';

  const adres = document.createElement('span');
  adres.className = 'dk-lancuch__adres';
  adres.textContent = 'wartość domyślna katalogu';

  const wartosc = document.createElement('code');
  wartosc.className = 'dk-lancuch__wartosc';
  wartosc.textContent = rozstrzygniecie.domyslna
    ? skroc(napis(rozstrzygniecie.wartosc))
    : '—';

  const znak = document.createElement('span');
  znak.className = rozstrzygniecie.domyslna
    ? 'dn-plakietka dn-plakietka--sygnal dk-lancuch__znak'
    : 'dn-plakietka dk-lancuch__znak';
  znak.textContent = rozstrzygniecie.domyslna ? 'obowiązuje' : 'zaplecze';

  wiersz.append(adres, wartosc, znak);
  return wiersz;
}

/**
 * Skraca wartość do jednego wiersza, żeby długi zapis nie rozpychał wykazu.
 * Wartość pusta dostaje napis `(pusto)`, dłuższa niż 72 znaki — wielokropek.
 */
function skroc(tekst: string): string {
  const jednowierszowy = tekst.replace(/\s+/gu, ' ').trim();
  if (jednowierszowy === '') return '(pusto)';
  return jednowierszowy.length > 72 ? `${jednowierszowy.slice(0, 71)}…` : jednowierszowy;
}
