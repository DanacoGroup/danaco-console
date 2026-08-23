// Emisja deklaracji narzedzi modelu w Go. Kazda deklaracja zajmuje
// jedna linie: literal zlozony zapisany w jednej linii nie podlega wyrownaniu
// kolumn, wiec wynik jest zgodny z gofmt bez uruchamiania gofmt.

import { mapaGo } from './blok-mapy-go.mjs';

/** Napis Go z zachowaniem znakow specjalnych. */
function napis(tekst) {
  return JSON.stringify(String(tekst));
}

/** Wycinek napisow Go. */
function wycinek(wartosci) {
  return `[]string{${wartosci.map((w) => napis(w)).join(', ')}}`;
}

/** Jedno pole schematu wejscia jako literal ToolParameter. */
function parametrGo(p) {
  const czesci = [
    `Name: ${napis(p.nazwa)}`,
    `Type: ${napis(p.typ)}`,
    `Items: ${napis(p.element)}`,
    `Required: ${p.wymagane}`,
    `Description: ${napis(p.opis)}`,
  ];
  if (p.wartosci.length) czesci.push(`Enum: ${wycinek(p.wartosci)}`);
  return `{${czesci.join(', ')}}`;
}

/** Jedna deklaracja narzedzia jako literal ToolDeclaration. */
function narzedzieGo(n) {
  const parametry = n.parametry.map(parametrGo).join(', ');
  return [
    `{Name: ${napis(n.nazwa)}`,
    `Command: Command${n.staly}`,
    `Description: ${napis(n.opis)}`,
    `Parameters: []ToolParameter{${parametry}}}`,
  ].join(', ');
}

/** Funkcja wykazu narzedzi oraz odwzorowanie nazwy narzedzia na nazwe komendy. */
export function narzedziaGo(model) {
  return [
    '// NarzedziaModelu zwraca komplet deklaracji narzedzi modelu.',
    '// Kanal modelu dostaje sterowanie platforma jako narzedzia, a nie jako',
    '// osobny parser intencji; kazda deklaracja odpowiada komendzie kontraktu.',
    'func NarzedziaModelu() []ToolDeclaration {',
    '\treturn []ToolDeclaration{',
    ...model.narzedzia.map((n) => `\t\t${narzedzieGo(n)},`),
    '\t}',
    '}',
    '',
    ...mapaGo(
      'KomendyNarzedzi',
      'nazwa narzedzia modelu odwzorowana na komende kontraktu',
      'string',
      'MessageType',
      model.narzedzia.map((n) => [napis(n.nazwa), `Command${n.staly}`]),
    ),
  ];
}
