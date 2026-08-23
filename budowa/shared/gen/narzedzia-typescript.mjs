// Emisja deklaracji narzedzi modelu w TypeScript. Klient wystawia ten
// sam wykaz co rdzen — nazwa narzedzia nie powstaje po zadnej ze stron osobno.

import { slownikTS } from './typy-typescript.mjs';

/** Napis TypeScript w apostrofach, z zabezpieczeniem znakow specjalnych. */
function napis(tekst) {
  return `'${String(tekst).replace(/\\/g, '\\\\').replace(/'/g, "\\'")}'`;
}

/** Tablica napisow TypeScript. */
function tablica(wartosci) {
  return `[${wartosci.map((w) => napis(w)).join(', ')}]`;
}

/** Jedno pole schematu wejscia jako literal ToolParameter. */
function parametrTS(p) {
  const czesci = [
    `name: ${napis(p.nazwa)}`,
    `type: ${napis(p.typ)}`,
    `items: ${napis(p.element)}`,
    `required: ${p.wymagane}`,
    `description: ${napis(p.opis)}`,
  ];
  if (p.wartosci.length) czesci.push(`enum: ${tablica(p.wartosci)}`);
  return `{ ${czesci.join(', ')} }`;
}

/** Jedna deklaracja narzedzia jako literal ToolDeclaration. */
function narzedzieTS(n) {
  const parametry = n.parametry.map(parametrTS).join(', ');
  return [
    `{ name: ${napis(n.nazwa)}`,
    `command: Command.${n.staly}`,
    `description: ${napis(n.opis)}`,
    `parameters: [${parametry}] }`,
  ].join(', ');
}

/** Wykaz narzedzi oraz odwzorowanie nazwy narzedzia na nazwe komendy. */
export function narzedziaTS(model) {
  return [
    '/**',
    ' * Komplet deklaracji narzedzi modelu. Kanal modelu dostaje sterowanie',
    ' * platforma jako narzedzia; kazda deklaracja odpowiada komendzie kontraktu.',
    ' */',
    'export const NARZEDZIA_MODELU: readonly ToolDeclaration[] = [',
    ...model.narzedzia.map((n) => `  ${narzedzieTS(n)},`),
    '];',
    '',
    ...slownikTS(
      'KOMENDY_NARZEDZI',
      'Nazwa narzedzia modelu odwzorowana na komende kontraktu',
      'string',
      'MessageType',
      model.narzedzia.map((n) => [n.nazwa, `Command.${n.staly}`]),
    ),
  ];
}
