import type { StudioDiffHunk, StudioTextMatch } from '../../../../shared/contract';

/**
 * Wyrys wyniku `studio.diff.compare` — fragmenty różnicy i trafienia wzorca.
 *
 * Kontrakt łączy porównanie wersji i wyszukiwanie wzorca w jednej komendzie:
 * odpowiedź niesie `hunks`, `matches` albo oba. Wyrys jest więc jeden i nie
 * zakłada, która tablica przyszła — brak obu jest poprawnym wynikiem, który
 * panel nazywa pustką merytoryczną.
 *
 * Rodzaj fragmentu jest daną, nie barwą w kodzie: `DiffHunkKind` trafia do
 * `data-rodzaj`, a barwę dobiera arkusz modułu z żetonów motywu.
 */
export function wyrysFragmentow(fragmenty: readonly StudioDiffHunk[]): HTMLElement[] {
  return fragmenty.map((fragment) => {
    const naglowek = document.createElement('p');
    naglowek.className = 'ms-roznica__naglowek';
    naglowek.textContent = `Fragment ${fragment.index} · ${fragment.kind}${zakresWierszy(fragment)}`;

    const element = document.createElement('li');
    element.className = 'ms-roznica__wiersz';
    element.dataset['rodzaj'] = fragment.kind;
    element.append(naglowek);

    if (fragment.before !== undefined) element.append(blok('przed', fragment.before));
    if (fragment.after !== undefined) element.append(blok('po', fragment.after));
    return element;
  });
}

export function wyrysTrafien(trafienia: readonly StudioTextMatch[]): HTMLElement[] {
  return trafienia.map((trafienie) => {
    const polozenie = document.createElement('span');
    polozenie.className = 'dn-plakietka ms-roznica__polozenie';
    polozenie.textContent =
      trafienie.column === undefined
        ? `wiersz ${trafienie.line}`
        : `wiersz ${trafienie.line}, kolumna ${trafienie.column}`;

    const tekst = document.createElement('code');
    tekst.className = 'ms-roznica__tekst';
    tekst.textContent = trafienie.text;

    const element = document.createElement('li');
    element.className = 'ms-roznica__wiersz';
    element.dataset['rodzaj'] = 'trafienie';
    element.append(polozenie, tekst);
    return element;
  });
}

/** Zakres wierszy fragmentu; pusty, gdy rdzeń go nie podał. */
function zakresWierszy(fragment: StudioDiffHunk): string {
  if (fragment.startLine === undefined) return '';
  const koniec = fragment.endLine ?? fragment.startLine;
  return ` · wiersze ${fragment.startLine}–${koniec}`;
}

/** Blok treści fragmentu wraz z etykietą strony porównania. */
function blok(strona: string, tresc: string): HTMLElement {
  const etykieta = document.createElement('span');
  etykieta.className = 'ms-roznica__strona';
  etykieta.textContent = strona;

  const tekst = document.createElement('pre');
  tekst.className = 'ms-roznica__tekst';
  tekst.textContent = tresc;

  const element = document.createElement('div');
  element.className = 'ms-roznica__blok';
  element.append(etykieta, tekst);
  return element;
}
