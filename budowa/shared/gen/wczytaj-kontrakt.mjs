// Wczytanie zrodla prawdy: shared/contract.json.

import { readFileSync } from 'node:fs';
import { resolve } from 'node:path';

/** Katalog shared/ wyliczony wzgledem katalogu gen/. */
export function katalogShared(katalogGen) {
  return resolve(katalogGen, '..');
}

/** Zwraca zawartosc contract.json jako obiekt. Blad parsowania przerywa generacje. */
export function wczytajKontrakt(sciezka) {
  const tekst = readFileSync(sciezka, 'utf8');
  try {
    return JSON.parse(tekst);
  } catch (blad) {
    throw new Error(`contract.json nie jest poprawnym JSON: ${blad.message}`);
  }
}
