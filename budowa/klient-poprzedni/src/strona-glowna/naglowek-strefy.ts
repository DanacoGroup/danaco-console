/**
 * Nagłówek strefy — etykieta wersalikowa i jedno zdanie wyjaśnienia.
 *
 * Nagłówek pozostaje lekki: hierarchię niesie proporcja kart, a nie napis nad
 * nimi. Stąd wzorzec `.dn-etykieta-wersalikowa` — 11 px, tracking 0,12em,
 * barwa --dn-tekst-3.
 */
export function utworzNaglowekStrefy(etykieta: string, wyjasnienie: string): HTMLElement {
  const element = document.createElement('header');
  element.className = 'dn-strona__naglowek-strefy';

  const napis = document.createElement('p');
  napis.className = 'dn-etykieta-wersalikowa dn-strona__etykieta-strefy';
  napis.textContent = etykieta;

  const opis = document.createElement('p');
  opis.className = 'dn-strona__wyjasnienie-strefy';
  opis.textContent = wyjasnienie;

  element.append(napis, opis);
  return element;
}
