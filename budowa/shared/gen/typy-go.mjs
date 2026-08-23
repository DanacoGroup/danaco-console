// Zapis pojedynczych deklaracji Go: typy pol, struktury, bloki stalych, listy.
// Wciecie tabulatorem i brak wyrownania kolumn sa celowe — pole poprzedzone
// komentarzem tworzy wlasna sekcje wyrownania, wiec wynik jest zgodny z gofmt.

import { naPascal, rozbijTyp } from './nazwy.mjs';

/** Odwzorowanie typow kontraktu na typy Go. */
const TYPY_PODSTAWOWE = {
  string: 'string',
  int: 'int',
  int64: 'int64',
  float: 'float64',
  bool: 'bool',
  json: 'json.RawMessage',
  MessageType: 'MessageType',
};

/** Zapis typu kontraktu; "Session[]" -> "[]Session", "json" -> "json.RawMessage". */
export function typGo(typ) {
  const { bazowy, tablica } = rozbijTyp(typ);
  const nazwa = TYPY_PODSTAWOWE[bazowy] ?? bazowy;
  return tablica ? `[]${nazwa}` : nazwa;
}

/**
 * Pole opcjonalne o typie prostym idzie wskaznikiem, aby brak wartosci dal sie
 * odroznic od wartosci zerowej. Wycinek i surowy JSON niosa brak przez nil.
 */
function przezWskaznik(pole) {
  if (pole.wymagane) return false;
  const { bazowy, tablica } = rozbijTyp(pole.typ);
  return !tablica && bazowy !== 'json';
}

/** Dwie linie pola struktury: komentarz i deklaracja ze znacznikiem JSON. */
function poleGo(pole) {
  const typ = przezWskaznik(pole) ? `*${typGo(pole.typ)}` : typGo(pole.typ);
  const znacznik = pole.wymagane ? pole.nazwa : `${pole.nazwa},omitempty`;
  return [`\t// ${pole.opis}`, `\t${naPascal(pole.nazwa)} ${typ} \`json:"${znacznik}"\``];
}

/** Struktura kontraktu wraz z komentarzem dokumentacyjnym Go. */
export function strukturaGo({ nazwa, opis, pola }) {
  const linie = [`// ${nazwa} — ${opis}`, `type ${nazwa} struct {`];
  for (const pole of pola) linie.push(...poleGo(pole));
  linie.push('}', '');
  return linie;
}

/** Nazwany typ napisowy wyliczenia. */
export function typWyliczeniaGo(nazwa, opis) {
  return [`// ${nazwa} — ${opis}`, `type ${nazwa} string`, ''];
}

/**
 * Blok stalych napisowych. Stale bez jawnego typu pozostaja nietypowane,
 * dzieki czemu wchodza wprost tam, gdzie warstwa wyzej ma wlasny typ napisowy.
 */
export function blokStalychGo(opis, wpisy, typ = '') {
  const przedrostek = typ ? ` ${typ}` : '';
  const linie = [`// ${opis}`, 'const ('];
  for (const w of wpisy) linie.push(`\t// ${w.opis}`, `\t${w.nazwa}${przedrostek} = "${w.wartosc}"`);
  linie.push(')', '');
  return linie;
}

/** Funkcja zwracajaca komplet nazw jako wycinek. */
export function funkcjaListyGo(nazwa, opis, typElementu, elementy) {
  return [
    `// ${nazwa} ${opis}`,
    `func ${nazwa}() []${typElementu} {`,
    `\treturn []${typElementu}{`,
    ...elementy.map((e) => `\t\t${e},`),
    '\t}',
    '}',
    '',
  ];
}
