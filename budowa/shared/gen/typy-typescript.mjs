// Zapis pojedynczych deklaracji TypeScript: typy pol, interfejsy, obiekty
// stalych. Modul nie zna ukladu calego pliku — sklada go emiter.

import { rozbijTyp } from './nazwy.mjs';

/** Odwzorowanie typow kontraktu na typy TypeScript. */
const TYPY_PODSTAWOWE = {
  string: 'string',
  int: 'number',
  int64: 'number',
  float: 'number',
  bool: 'boolean',
  json: 'unknown',
  MessageType: 'MessageType',
};

/** Zapis typu kontraktu; "Session[]" -> "Session[]", "int64" -> "number". */
export function typTS(typ) {
  const { bazowy, tablica } = rozbijTyp(typ);
  const nazwa = TYPY_PODSTAWOWE[bazowy] ?? bazowy;
  return tablica ? `${nazwa}[]` : nazwa;
}

/** Dwie linie pola: komentarz dokumentacyjny i deklaracja. */
function poleTS(pole, nadpisania) {
  const znak = pole.wymagane ? '' : '?';
  const typ = nadpisania[pole.nazwa] ?? typTS(pole.typ);
  return [`  /** ${pole.opis} */`, `  ${pole.nazwa}${znak}: ${typ};`];
}

/**
 * Interfejs kontraktu. `parametry` niesie ewentualna liste parametrow typu,
 * `nadpisania` zamienia typ wybranego pola (uzywane dla payload koperty).
 */
export function interfejsTS({ nazwa, opis, pola }, parametry = '', nadpisania = {}) {
  const linie = [`/** ${opis} */`, `export interface ${nazwa}${parametry} {`];
  for (const pole of pola) linie.push(...poleTS(pole, nadpisania));
  linie.push('}', '');
  return linie;
}

/**
 * Obiekt stalych wraz z typem sumy literalow. Jedna postac dla wyliczen,
 * komend, zdarzen i kodow bledow — nazwa kontraktu zyje w jednym miejscu.
 */
export function obiektStalychTS(nazwa, opis, wpisy) {
  const linie = [`/** ${opis} */`, `export const ${nazwa} = {`];
  for (const w of wpisy) linie.push(`  /** ${w.opis} */`, `  ${w.klucz}: '${w.wartosc}',`);
  linie.push('} as const;');
  linie.push(`export type ${nazwa} = (typeof ${nazwa})[keyof typeof ${nazwa}];`, '');
  return linie;
}

/** Tablica literalow tylko do odczytu, np. lista modulow znanych kontraktowi. */
export function tablicaTS(nazwa, opis, wartosci) {
  const pozycje = wartosci.map((w) => `'${w}'`).join(', ');
  return [`/** ${opis} */`, `export const ${nazwa} = [${pozycje}] as const;`, ''];
}

/**
 * Odwzorowanie klucz -> wartosc o kluczach bedacych literalami kontraktu.
 *
 * `czesciowy` oznacza slownik, ktory z zalozenia nie pokrywa wszystkich kluczy
 * — tak jest przy wyliczeniu z wartosciami przelotowymi, ktore nie maja
 * odpowiednika w kolumnie modelu danych. TypeScript ma wtedy wiedziec, ze
 * odczyt moze nie trafic, zamiast obiecywac wartosc dla kazdego klucza.
 */
export function slownikTS(nazwa, opis, typKlucza, typWartosci, pary, czesciowy = false) {
  const zapis = czesciowy
    ? `Readonly<Partial<Record<${typKlucza}, ${typWartosci}>>>`
    : `Readonly<Record<${typKlucza}, ${typWartosci}>>`;
  const linie = [
    `/** ${opis} */`,
    `export const ${nazwa}: ${zapis} = {`,
  ];
  for (const [klucz, wartosc] of pary) linie.push(`  '${klucz}': ${wartosc},`);
  linie.push('};', '');
  return linie;
}
