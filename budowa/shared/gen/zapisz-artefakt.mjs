// Zapis wytworzonego artefaktu na dysk wraz z krotkim sprawozdaniem.

import { writeFileSync } from 'node:fs';
import { join } from 'node:path';

/** Zapisuje tresc do katalogu shared/ i zwraca opis wyniku. */
export function zapiszArtefakt(katalog, nazwa, tresc) {
  const sciezka = join(katalog, nazwa);
  const zakonczona = tresc.endsWith('\n') ? tresc : `${tresc}\n`;
  writeFileSync(sciezka, zakonczona, 'utf8');
  return { sciezka, linie: zakonczona.split('\n').length - 1, bajty: Buffer.byteLength(zakonczona) };
}

/** Jedna linia sprawozdania dla wypisu na wyjscie standardowe. */
export function opisArtefaktu({ sciezka, linie, bajty }) {
  return `  ${sciezka} — ${linie} linii, ${bajty} B`;
}
