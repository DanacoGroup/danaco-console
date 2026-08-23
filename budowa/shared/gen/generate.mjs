#!/usr/bin/env node
// Punkt wejscia generatora kontraktu. Wylacznie kompozycja: wczytanie
// zrodla prawdy, zbudowanie modelu, emisja bindingow i zapis artefaktow.
//
// Uruchomienie:  node shared/gen/generate.mjs

import { dirname, join } from 'node:path';
import { fileURLToPath } from 'node:url';

import { katalogShared, wczytajKontrakt } from './wczytaj-kontrakt.mjs';
import { zbudujModel } from './model-kontraktu.mjs';
import { emitujTypeScript } from './emiter-typescript.mjs';
import { emitujGo } from './emiter-go.mjs';
import { opisArtefaktu, zapiszArtefakt } from './zapisz-artefakt.mjs';

const katalogGen = dirname(fileURLToPath(import.meta.url));
const katalog = katalogShared(katalogGen);
const zrodlo = wczytajKontrakt(join(katalog, 'contract.json'));
const model = zbudujModel(zrodlo);

const artefakty = [
  zapiszArtefakt(katalog, model.kontrakt.artefakty.typescript, emitujTypeScript(model)),
  zapiszArtefakt(katalog, model.kontrakt.artefakty.go, emitujGo(model)),
];

console.log(`Kontrakt ${model.kontrakt.produkt} — protokol ${model.kontrakt.protokol}`);
console.log(`  komendy: ${model.komendy.length} · zdarzenia: ${model.zdarzenia.length} · kody bledow: ${model.kodyBledow.length}`);
console.log(`  narzedzia modelu: ${model.narzedzia.length}`);
for (const artefakt of artefakty) console.log(opisArtefaktu(artefakt));
