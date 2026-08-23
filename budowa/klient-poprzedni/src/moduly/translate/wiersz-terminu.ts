import type { GlossaryTerm } from '../../../../shared/contract';
import { przycisk } from '../../modele/kontrolki-formularza';

/**
 * Jeden termin glosariusza w wykazie okna zarządcy.
 *
 * Jeden wiersz na plik — tak samo jak wiersz biblioteki ekspertów w module
 * Agents. Wiersz niczego nie woła: oddaje dwa naciśnięcia oknu, które wie,
 * co z nimi zrobić. Dzięki temu ten sam wiersz nadaje się i do wykazu, i do
 * zestawienia wystąpień.
 */
export interface CzynnosciWiersza {
  /** „Edytuj" — wciąga termin z powrotem do formularza definicji. */
  naEdycje(termin: GlossaryTerm): void;
  /** „Pokaż wystąpienia" — pyta rdzeń o miejsca użycia terminu. */
  naWystapienia(termin: GlossaryTerm): void;
}

export function utworzWierszTerminu(
  termin: GlossaryTerm,
  czynnosci: CzynnosciWiersza,
): HTMLElement {
  const nazwa = document.createElement('span');
  nazwa.className = 'mt-termin__zrodlo';
  nazwa.textContent = termin.source;

  const odpowiednik = document.createElement('span');
  odpowiednik.className = 'mt-termin__odpowiednik';
  odpowiednik.textContent = termin.doNotTranslate === true
    ? `${termin.language}: nie tłumacz`
    : `${termin.language}: ${termin.target ?? '—'}`;

  const uwaga = document.createElement('span');
  uwaga.className = 'mt-termin__uwaga';
  uwaga.textContent = termin.note ?? '';

  const edytuj = przycisk('Edytuj', 'dn-btn dn-btn--sm dn-btn--duch');
  edytuj.addEventListener('click', () => czynnosci.naEdycje(termin));

  const wystapienia = przycisk('Pokaż wystąpienia', 'dn-btn dn-btn--sm dn-btn--duch');
  wystapienia.addEventListener('click', () => czynnosci.naWystapienia(termin));

  const element = document.createElement('li');
  element.className = 'mt-termin';
  element.dataset['termin'] = termin.id;
  if (termin.doNotTranslate === true) element.dataset['nietlumaczony'] = 'tak';
  element.append(nazwa, odpowiednik, uwaga, edytuj, wystapienia);
  return element;
}
