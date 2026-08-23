// Narzedzia modelu: deklaracja narzedzia powstaje z komendy kontraktu.
// Nazwa, opis i schemat wejscia sa wyprowadzone, nie przepisane — dzieki temu
// zmiana ksztaltu zadania komendy zmienia schemat narzedzia w tej samej chwili.
// Modul nie zna skladni jezykow docelowych; wytwarza wylacznie model.

import { naPascal, rozbijTyp } from './nazwy.mjs';

/** Odwzorowanie typow kontraktu na typy schematu wejscia (JSON Schema). */
const TYPY_SCHEMATU = {
  string: 'string',
  int: 'integer',
  int64: 'integer',
  float: 'number',
  bool: 'boolean',
  json: 'object',
  MessageType: 'string',
};

/** Typ schematu dla nazwy typu kontraktu; wyliczenie idzie napisem, struktura obiektem. */
function typSchematu(bazowy, wyliczenia, struktury) {
  if (TYPY_SCHEMATU[bazowy]) return TYPY_SCHEMATU[bazowy];
  if (wyliczenia.has(bazowy)) return 'string';
  if (struktury.has(bazowy)) return 'object';
  throw new Error(`Narzedzia modelu: typ "${bazowy}" nie ma odpowiednika w schemacie wejscia`);
}

/**
 * Pojedyncze pole schematu wejscia. Pole tablicowe niesie typ elementu w `element`,
 * pole o typie wyliczenia — komplet dopuszczalnych wartosci w `wartosci`.
 */
function parametr(pole, wyliczenia, struktury) {
  const { bazowy, tablica } = rozbijTyp(pole.typ);
  const typElementu = typSchematu(bazowy, wyliczenia, struktury);
  return {
    nazwa: pole.nazwa,
    typ: tablica ? 'array' : typElementu,
    element: tablica ? typElementu : '',
    wymagane: Boolean(pole.wymagane),
    opis: pole.opis,
    wartosci: wyliczenia.get(bazowy) ?? [],
  };
}

/** Nazwa narzedzia z nazwy komendy: "session.create" -> "danaco_session_create". */
export function nazwaNarzedzia(typKomendy, przedrostek, separator) {
  return [przedrostek, ...String(typKomendy).split('.')].join(separator);
}

/** Komplet deklaracji narzedzi modelu; kazda pozycja wskazuje istniejaca komende. */
export function narzedziaModelu(zrodlo) {
  const { przedrostek, separatorNazwy, pozycje } = zrodlo.narzedzia;
  const komendy = new Map(zrodlo.komendy.map((k) => [k.typ, k]));
  const wyliczenia = new Map(zrodlo.wyliczenia.map((w) => [w.nazwa, w.wartosci.map((v) => v.wartosc)]));
  const struktury = new Set(zrodlo.struktury.map((s) => s.nazwa));
  return pozycje.map((p) => {
    const komenda = komendy.get(p.komenda);
    if (!komenda) {
      throw new Error(`Narzedzia modelu: komenda "${p.komenda}" nie wystepuje w kontrakcie`);
    }
    return {
      nazwa: nazwaNarzedzia(komenda.typ, przedrostek, separatorNazwy),
      komenda: komenda.typ,
      staly: naPascal(komenda.typ),
      opis: `${komenda.opis}. ${p.zastosowanie}`,
      parametry: komenda.zadanie.map((pole) => parametr(pole, wyliczenia, struktury)),
    };
  });
}
