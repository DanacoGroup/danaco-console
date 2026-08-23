// Emiter TypeScript: sklada contract.ts z deklaracji wytworzonych przez
// typy-typescript.mjs. Zaden literal nazwy kontraktu nie powstaje tutaj —
// wszystkie pochodza z modelu.

import { naPascal, naWielkieZPodkresleniem } from './nazwy.mjs';
import { linieNaglowka } from './naglowek.mjs';
import { opisDoBazy, opisZBazy } from './odwzorowanie-bazy.mjs';
import { narzedziaTS } from './narzedzia-typescript.mjs';
import {
  interfejsTS,
  obiektStalychTS,
  slownikTS,
  tablicaTS,
} from './typy-typescript.mjs';

/** Naglowek pliku jako komentarz blokowy. */
function naglowek(kontrakt) {
  const tresc = linieNaglowka(kontrakt, kontrakt.artefakty.typescript);
  return ['/*', ...tresc.map((l) => (l ? ` * ${l}` : ' *')), ' */', ''];
}

/** Wyliczenia kontraktu jako obiekty stalych z suma literalow. */
function wyliczenia(model) {
  return model.wyliczenia.flatMap((w) =>
    obiektStalychTS(
      w.nazwa,
      w.opis,
      w.wartosci.map((v) => ({ klucz: naPascal(v.wartosc), opis: v.opis, wartosc: v.wartosc })),
    ),
  );
}

/**
 * Odwzorowanie wartosci wyliczenia na wartosc kolumny modelu danych i z powrotem.
 * Klient siega po te slowniki zamiast wpisywac przeklad u siebie.
 */
function odwzorowania(model) {
  return model.odwzorowaniaBazy.flatMap((o) => {
    const przyrostek = naWielkieZPodkresleniem(o.nazwa);
    return [
      ...slownikTS(`WARTOSCI_BAZY_${przyrostek}`, opisDoBazy(o), o.nazwa, 'string',
        o.wartosci.map((w) => [w.wartosc, `'${w.baza}'`]), !o.pelne),
      ...slownikTS(`WARTOSCI_KONTRAKTU_${przyrostek}`, opisZBazy(o), 'string', o.nazwa,
        o.wartosci.map((w) => [w.baza, `${o.nazwa}.${naPascal(w.wartosc)}`])),
    ];
  });
}

/** Nazwy komend, nazwy zdarzen i wspolny typ komunikatu. */
function nazwyKomunikatow(model) {
  const komendy = model.komendy.map((k) => ({ klucz: k.staly, opis: k.opis, wartosc: k.typ }));
  const zdarzenia = model.zdarzenia.map((z) => ({ klucz: z.staly, opis: z.opis, wartosc: z.typ }));
  return [
    ...obiektStalychTS('Command', 'Nazwy komend kontraktu', komendy),
    ...obiektStalychTS('EventType', 'Nazwy zdarzen kontraktu', zdarzenia),
    '/** Typ komunikatu koperty: komenda albo zdarzenie. */',
    'export type MessageType = Command | EventType;',
    '',
  ];
}

/** Kody bledow wraz z informacja o sensownosci ponowienia. */
function kodyBledow(model) {
  const wpisy = model.kodyBledow.map((b) => ({ klucz: naPascal(b.kod), opis: b.opis, wartosc: b.kod }));
  return [
    ...obiektStalychTS('ErrorCode', 'Kody bledow kontraktu', wpisy),
    ...slownikTS(
      'KODY_PONAWIALNE',
      'Czy ponowienie zadania po danym kodzie ma sens',
      'ErrorCode',
      'boolean',
      model.kodyBledow.map((b) => [b.kod, String(b.retryable)]),
    ),
  ];
}

/** Koperta, struktury dziedzinowe i tresci komunikatow. */
function typyDanych(model) {
  const koperta = { nazwa: 'Envelope', opis: model.koperta.opis, pola: model.koperta.pola };
  return [
    ...interfejsTS(koperta, '<T = unknown>', { payload: 'T' }),
    ...model.struktury.flatMap((s) => interfejsTS(s)),
    ...model.tresci.flatMap((t) => interfejsTS(t)),
  ];
}

/** Powiazanie nazwy komunikatu z typem jego tresci. */
function powiazania(model) {
  const linie = ['/** Tresc zadania i wyniku kazdej komendy. */', 'export interface CommandPayloads {'];
  for (const k of model.komendy) {
    linie.push(`  /** ${k.opis} */`);
    linie.push(`  '${k.typ}': { request: ${k.typZadania}; response: ${k.typWyniku} };`);
  }
  linie.push('}', '');
  linie.push('/** Tresc kazdego zdarzenia. */', 'export interface EventPayloads {');
  for (const z of model.zdarzenia) {
    linie.push(`  /** ${z.opis} */`, `  '${z.typ}': ${z.typTresci};`);
  }
  linie.push('}', '');
  linie.push(
    '/** Tresc zadania wskazanej komendy. */',
    'export type RequestOf<K extends Command> = CommandPayloads[K][\'request\'];',
    '/** Tresc wyniku wskazanej komendy. */',
    'export type ResponseOf<K extends Command> = CommandPayloads[K][\'response\'];',
    '/** Tresc wskazanego zdarzenia. */',
    'export type EventPayloadOf<K extends EventType> = EventPayloads[K];',
    '',
  );
  return linie;
}

/** Listy nazw i rozpoznanie typu — bez powielania literalow po stronie klienta. */
function rozpoznanie(model) {
  const komendy = model.komendy.map((k) => `Command.${k.staly}`).join(', ');
  const zdarzenia = model.zdarzenia.map((z) => `EventType.${z.staly}`).join(', ');
  return [
    '/** Komplet nazw komend kontraktu. */',
    `export const KOMENDY: readonly Command[] = [${komendy}];`,
    '',
    '/** Komplet nazw zdarzen kontraktu. */',
    `export const ZDARZENIA: readonly EventType[] = [${zdarzenia}];`,
    '',
    '/** Czy nazwa jest komenda kontraktu. */',
    'export function czyKomenda(typ: string): typ is Command {',
    '  return (KOMENDY as readonly string[]).includes(typ);',
    '}',
    '',
    '/** Czy nazwa jest zdarzeniem kontraktu. */',
    'export function czyZdarzenie(typ: string): typ is EventType {',
    '  return (ZDARZENIA as readonly string[]).includes(typ);',
    '}',
    '',
    ...slownikTS(
      'ZDARZENIA_NIEZNANEJ',
      'Zdarzenie *.unknown wlasciwe dla obszaru nazwy',
      'string',
      'EventType',
      model.zdarzenia.filter((z) => z.nieznane).map((z) => [z.obszar, `EventType.${z.staly}`]),
    ),
    '/**',
    ' * Zdarzenie zwracane na nieznana komende (fail-open):',
    ' * polaczenie nie jest zrywane, obszar spoza kontraktu dostaje zdarzenie zapasowe.',
    ' */',
    'export function zdarzenieNieznanej(typ: string): EventType {',
    '  const obszar = typ.split(SEPARATOR_OBSZARU)[0];',
    `  return ZDARZENIA_NIEZNANEJ[obszar] ?? EventType.${model.zdarzenieZapasowe};`,
    '}',
    '',
  ];
}

/** Listy wartosci znanych w chwili wydania kontraktu — informacyjne, nie bramy. */
function listyZnane(model) {
  return model.listyZnane.flatMap((l) => tablicaTS(l.nazwa, l.opis, l.wartosci));
}

export function emitujTypeScript(model) {
  return [
    ...naglowek(model.kontrakt),
    '/** Wersja protokolu kontraktu. */',
    `export const PROTOCOL_VERSION = '${model.kontrakt.protokol}';`,
    '',
    '/** Separator obszaru w notacji nazwy kontraktu. */',
    `export const SEPARATOR_OBSZARU = '${model.kontrakt.separator}';`,
    '',
    ...wyliczenia(model),
    ...odwzorowania(model),
    ...nazwyKomunikatow(model),
    ...kodyBledow(model),
    ...typyDanych(model),
    ...powiazania(model),
    ...rozpoznanie(model),
    ...narzedziaTS(model),
    ...listyZnane(model),
  ].join('\n');
}
