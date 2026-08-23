// Literal mapy w Go. gofmt wyrownuje wartosci w kolumnie, wiec generator musi
// dopelnic klucze sam — inaczej `gofmt -l` zglosi plik.

/** Szerokosc kolumny klucza: najdluzszy klucz wraz z dwukropkiem. */
function szerokoscKolumny(pary) {
  return pary.reduce((maks, [klucz]) => Math.max(maks, klucz.length + 1), 0);
}

/** Zmienna pakietowa z literalem mapy; `pary` to lista [klucz, wartosc]. */
export function mapaGo(nazwa, opis, typKlucza, typWartosci, pary, eksport = true) {
  const szerokosc = szerokoscKolumny(pary);
  const komentarz = eksport ? `// ${nazwa} — ${opis}` : `// ${opis}`;
  const linie = [komentarz, `var ${nazwa} = map[${typKlucza}]${typWartosci}{`];
  for (const [klucz, wartosc] of pary) {
    linie.push(`\t${`${klucz}:`.padEnd(szerokosc)} ${wartosc},`);
  }
  linie.push('}', '');
  return linie;
}
