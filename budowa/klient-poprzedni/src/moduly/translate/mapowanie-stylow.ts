import { poleWyboru, przycisk } from '../../modele/kontrolki-formularza';
import { utworzRozwiniecie } from './warstwy-translate';
import { STYLE_DOCELOWE, przelozStyl, type StylDocelowy } from './wzorce-placeholderow';

/**
 * Mapowanie stylów symboli zastępczych przekłada zapis zmiennej między platformami na tekście
 * źródłowym, licząc okno na przeniesienie ręczne, bo rdzeń nie ma komendy tej zmiany.
 */
export interface MapowanieStylow {
  element: HTMLElement;
}

export function utworzMapowanieStylow(tekstZrodlowy: () => string): MapowanieStylow {
  const rozwiniecie = utworzRozwiniecie({
    warstwa: 4,
    nazwa: 'Mapowanie stylów symboli zastępczych',
    wyjasnienie:
      'Przekład zapisu zmiennych między platformami. Liczy go okno — kontrakt nie ma komendy ' +
      'zamiany stylu, więc wynik służy przeniesieniu ręcznemu.',
    znacznik: '☰',
  });

  const styl = poleWyboru({ etykieta: 'Styl docelowy' }, [
    { wartosc: STYLE_DOCELOWE.printf, etykieta: 'wzorzec printf — %s' },
    { wartosc: STYLE_DOCELOWE.indeks, etykieta: 'indeks w nawiasach klamrowych — {0}' },
    { wartosc: STYLE_DOCELOWE.nazwa, etykieta: 'podwójne nawiasy klamrowe — {{nazwa}}' },
  ]);

  const zdanie = document.createElement('p');
  zdanie.className = 'mt-mapowanie__zdanie';

  const wynik = document.createElement('pre');
  wynik.className = 'mt-mapowanie__wynik';

  const zamiany = document.createElement('ul');
  zamiany.className = 'mt-mapowanie__zamiany';

  const przeloz = przycisk('Przełóż styl w tekście źródłowym', 'dn-btn dn-btn--sm dn-btn--zarys');
  przeloz.addEventListener('click', () => {
    const przeklad = przelozStyl(tekstZrodlowy(), styl.kontrolka.value as StylDocelowy);
    wynik.textContent = przeklad.tekst;
    zamiany.replaceChildren(...przeklad.zamiany.map(wierszZamiany));
    zdanie.textContent = zdanieOPrzekladzie(przeklad.zamiany.length, przeklad.nazwyZNumeru);
  });

  const pasek = document.createElement('div');
  pasek.className = 'mt-pasek';
  pasek.append(przeloz);

  rozwiniecie.tresc.append(styl.element, pasek, zdanie, zamiany, wynik);
  return { element: rozwiniecie.element };
}

function zdanieOPrzekladzie(ile: number, nazwyZNumeru: boolean): string {
  if (ile === 0) {
    return (
      'Materiał zapisany w rdzeniu nie ma symboli zastępczych w zapisach znanych oknu, więc nie ' +
      'było czego przełożyć. Przekład bierze tekst źródłowy zapisany, nie treść wpisaną do pola. ' +
      'Znaczniki formatu zostają nietknięte zawsze — to struktura treści, nie zmienne.'
    );
  }
  const podstawa =
    `Przełożono ${String(ile)} wystąpień. Wynik powstał w oknie i nigdzie się nie zapisał — ` +
    'przenieś go do materiału ręcznie.';
  if (!nazwyZNumeru) return podstawa;
  return (
    `${podstawa} Część zapisów źródłowych nie niosła nazwy zmiennej, więc nazwą został numer ` +
    'kolejny wystąpienia; nazwy zmiennej okno nie zgaduje.'
  );
}

function wierszZamiany(zamiana: string): HTMLElement {
  const element = document.createElement('li');
  element.className = 'mt-mapowanie__zamiana';
  element.textContent = zamiana;
  return element;
}
