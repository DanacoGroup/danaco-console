#!/usr/bin/env node
// Odprawa terytorium — skorowidz wskazanego obszaru drzewa, wypisywany na wyjście
// standardowe: wykaz plików z liczbą linii i tym, co każdy eksportuje, komendy
// kontraktu, których obszar dotyka, oraz rodziny komend w nim nietknięte.
//
// Skorowidz powstaje z odczytu plików, nie z dokumentu opisującego drzewo, więc
// nie rozjeżdża się ze stanem faktycznym.
//
// Użycie:
//   node budowa/pomocniki/odprawa.mjs client/src/moduly/wiedza
//   node budowa/pomocniki/odprawa.mjs server/internal/poczta server/internal/wiedza

import { readFileSync, readdirSync, statSync, existsSync } from 'node:fs';
import { join, relative, extname } from 'node:path';
import { fileURLToPath } from 'node:url';

const KORZEN = join(fileURLToPath(new URL('.', import.meta.url)), '..');
const KOD = new Set(['.ts', '.go', '.rs', '.mjs', '.js']);

// Pliki pochodne i zależności zewnętrzne — nie należą do żadnego terytorium.
const POMIJANE = /node_modules|\/dist\/|\/target\/|\.przejazd|\/gen\/|\.min\./;

function pliki(katalog) {
  const wynik = [];
  const chodz = (k) => {
    let wpisy;
    try { wpisy = readdirSync(k, { withFileTypes: true }); } catch { return; }
    for (const w of wpisy) {
      const p = join(k, w.name);
      if (POMIJANE.test(p)) continue;
      if (w.isDirectory()) chodz(p);
      else wynik.push(p);
    }
  };
  chodz(katalog);
  return wynik.sort();
}

// Co plik oddaje na zewnątrz. Wykaz eksportów pozwala rozpoznać, co w obszarze już
// istnieje, bez otwierania każdego pliku z osobna.
function eksporty(sciezka) {
  let tresc;
  try { tresc = readFileSync(sciezka, 'utf8'); } catch { return []; }
  const rozsz = extname(sciezka);
  const nazwy = [];
  if (rozsz === '.ts' || rozsz === '.mjs' || rozsz === '.js') {
    for (const m of tresc.matchAll(/^export\s+(?:async\s+)?(?:function|const|class|interface|type|enum)\s+([A-Za-z0-9_]+)/gm)) nazwy.push(m[1]);
  } else if (rozsz === '.go') {
    for (const m of tresc.matchAll(/^func\s+(?:\([^)]*\)\s*)?([A-Z][A-Za-z0-9_]*)/gm)) nazwy.push(m[1] + '()');
    for (const m of tresc.matchAll(/^type\s+([A-Z][A-Za-z0-9_]*)/gm)) nazwy.push(m[1]);
  } else if (rozsz === '.rs') {
    for (const m of tresc.matchAll(/^pub\s+(?:async\s+)?(?:fn|struct|enum|trait)\s+([a-zA-Z0-9_]+)/gm)) nazwy.push(m[1]);
  }
  return [...new Set(nazwy)];
}

function linie(sciezka) {
  try { return readFileSync(sciezka, 'utf8').split('\n').length; } catch { return 0; }
}

// Komendy kontraktu, których terytorium dotyka — z obu stron: te, które woła, i te,
// które obsługuje. Trafienie rozpoznaje się po nazwie komendy występującej w treści pliku.
function komendyTerytorium(sciezki) {
  const kontrakt = JSON.parse(readFileSync(join(KORZEN, 'shared/contract.json'), 'utf8'));
  const wszystkie = kontrakt.komendy.map((k) => k.typ);
  const trafione = new Map();
  for (const s of sciezki) {
    let tresc;
    try { tresc = readFileSync(s, 'utf8'); } catch { continue; }
    for (const nazwa of wszystkie) {
      if (tresc.includes(nazwa)) {
        if (!trafione.has(nazwa)) trafione.set(nazwa, []);
        trafione.get(nazwa).push(relative(KORZEN, s));
      }
    }
  }
  return { trafione, wszystkie };
}

// Rodziny komend kontraktu pokrewne terytorium — grupowane po przedrostku nazwy.
// Wskazują komendy z tych samych rodzin, których obszar nie dotyka.
function rodziny(trafione, wszystkie) {
  const przedrostki = new Set([...trafione.keys()].map((n) => n.split('.')[0]));
  const wynik = new Map();
  for (const przedrostek of przedrostki) {
    const rodzina = wszystkie.filter((n) => n.startsWith(przedrostek + '.'));
    const brakujace = rodzina.filter((n) => !trafione.has(n));
    if (brakujace.length) wynik.set(przedrostek, brakujace);
  }
  return wynik;
}

const cele = process.argv.slice(2);
if (!cele.length) {
  console.error('Podaj terytorium, np.: node budowa/pomocniki/odprawa.mjs client/src/moduly/wiedza');
  process.exit(2);
}

const wszystkiePliki = [];
for (const cel of cele) {
  const pelna = join(KORZEN, cel);
  if (!existsSync(pelna)) { console.error(`NIE MA TAKIEGO TERYTORIUM: ${cel}`); process.exit(2); }
  wszystkiePliki.push(...(statSync(pelna).isDirectory() ? pliki(pelna) : [pelna]));
}

const kodowe = wszystkiePliki.filter((p) => KOD.has(extname(p)));
const pozostale = wszystkiePliki.filter((p) => !KOD.has(extname(p)));
const sumaLinii = kodowe.reduce((s, p) => s + linie(p), 0);

console.log(`# Odprawa terytorium: ${cele.join(' + ')}\n`);
console.log(`**${kodowe.length} plików kodu · ${sumaLinii} linii** (plus ${pozostale.length} plików nie-kodowych: style, dane, opisy).`);
console.log(`\n> To jest CAŁE twoje terytorium. Nie przeszukuj drzewa — nie ma tam nic więcej, co należy do ciebie.`);
console.log(`> Potrzebujesz pliku spoza tego wykazu? Odprawa jest zła — ZAMELDUJ TO, zamiast szukać.\n`);

console.log(`## Pliki i to, co oddają na zewnątrz\n`);
for (const p of kodowe) {
  const e = eksporty(p);
  const nazwa = relative(KORZEN, p);
  console.log(`- \`${nazwa}\` — ${linie(p)} lin.${e.length ? ' → ' + e.slice(0, 12).join(', ') + (e.length > 12 ? ` (+${e.length - 12})` : '') : ' → nic nie eksportuje'}`);
}
if (pozostale.length) {
  console.log(`\n### Nie-kodowe\n`);
  for (const p of pozostale) console.log(`- \`${relative(KORZEN, p)}\` — ${linie(p)} lin.`);
}

const { trafione, wszystkie } = komendyTerytorium(kodowe);
console.log(`\n## Komendy kontraktu, których to terytorium dotyka (${trafione.size} z ${wszystkie.length})\n`);
if (!trafione.size) {
  console.log(`ŻADNEJ. To terytorium nie styka się dziś z kontraktem — jeśli ma wykonywać pracę rdzenia, to jest brak, nie cecha.`);
} else {
  for (const [nazwa, gdzie] of [...trafione].sort()) console.log(`- \`${nazwa}\` — w ${gdzie.length === 1 ? gdzie[0] : gdzie.length + ' plikach'}`);
}

const brakujaceRodziny = rodziny(trafione, wszystkie);
if (brakujaceRodziny.size) {
  console.log(`\n### Z tych samych rodzin, a NIE tknięte tutaj — sprawdź, czy to brak\n`);
  for (const [przedrostek, lista] of brakujaceRodziny) {
    console.log(`- **${przedrostek}.\\*** → ${lista.join(', ')}`);
  }
}
console.log(`\n---\nOdprawa wygenerowana z drzewa, nie z dokumentu. Dokument bywa starszy o dwa dni.`);
