import { elementIkony, type NazwaIkony } from '../ikony/ikony';

/**
 * Stan pusty sekcji pulpitu na komponencie `dn-pusty-stan`.
 *
 * Sekcja bez danych dostaje tytuł i zdanie mówiące, czego brakuje: odczytu,
 * który jeszcze nie nadszedł, albo źródła, którego kontrakt nie ma.
 */
export function utworzStanPusty(
  tytul: string,
  opis: string,
  ikona: NazwaIkony = 'info',
): HTMLElement {
  const element = document.createElement('div');
  element.className = 'dn-pusty-stan mc-pusto';

  const naglowek = document.createElement('span');
  naglowek.className = 'dn-pusty-stan-tytul';
  naglowek.textContent = tytul;

  const zdanie = document.createElement('span');
  zdanie.className = 'dn-pusty-stan-opis';
  zdanie.textContent = opis;

  element.append(elementIkony(ikona, { rozmiar: 28 }), naglowek, zdanie);
  return element;
}
