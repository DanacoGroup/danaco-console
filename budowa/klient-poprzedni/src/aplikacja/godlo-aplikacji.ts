import { elementGodla } from '../ikony/ikony';

/**
 * Godło marki i tożsamość projektu pokazywane na pasku aplikacji. Znak stawia
 * warstwa marki wołana przez `ikony/ikony.ts`, a nazwa produktu pada tutaj raz
 * i tylko stąd trafia na pasek.
 */
const NAZWA_PRODUKTU = 'Danaco Console';

export function utworzGodlo(projekt: string): HTMLElement {
  const element = document.createElement('div');
  element.className = 'dn-pasek-godlo dn-powloka__godlo';

  const godlo = elementGodla({
    rozmiar: 24,
    podloze: 'ciemne',
    etykieta: 'Danaco Holding Group',
  });

  const nazwa = document.createElement('span');
  nazwa.className = 'dn-powloka__nazwa';
  nazwa.textContent = NAZWA_PRODUKTU;

  element.append(godlo, nazwa, ...nazwaProjektu(projekt));
  return element;
}

/**
 * Nazwa projektu oddana jako wykaz elementów: pozostaje pusty, gdy nazwa nie
 * została podana albo powtarza nazwę produktu, ponieważ powtórzenie zajmuje
 * miejsce na pasku i niczego nie wnosi.
 */
function nazwaProjektu(projekt: string): HTMLElement[] {
  const tresc = projekt.trim();
  if (tresc === '' || tresc.toLocaleLowerCase('pl') === NAZWA_PRODUKTU.toLocaleLowerCase('pl')) {
    return [];
  }
  const element = document.createElement('span');
  element.className = 'dn-powloka__projekt';
  element.textContent = tresc;
  return [element];
}
