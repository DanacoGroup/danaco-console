// Nazewnictwo: zamiana nazw kontraktu na identyfikatory jezykow docelowych.
// Modul nie zna skladni Go ani TypeScript — wytwarza wylacznie nazwy.

const ROZDZIELACZ = /[^A-Za-z0-9]+/;

/** "session.create" -> "SessionCreate", "tool_use" -> "ToolUse" */
export function naPascal(tekst) {
  return String(tekst)
    .split(ROZDZIELACZ)
    .filter(Boolean)
    .map((czlon) => czlon.charAt(0).toUpperCase() + czlon.slice(1))
    .join('');
}

/** "ExecutionEnv" -> "EXECUTION_ENV" — postac nazwy stalej modulowej TypeScript */
export function naWielkieZPodkresleniem(tekst) {
  return String(tekst)
    .replace(/([a-z0-9])([A-Z])/g, '$1_$2')
    .replace(/[^A-Za-z0-9]+/g, '_')
    .toUpperCase();
}

/** Rozbija zapis typu kontraktu: "Session[]" -> { bazowy: "Session", tablica: true } */
export function rozbijTyp(typ) {
  const tekst = String(typ);
  const tablica = tekst.endsWith('[]');
  return { bazowy: tablica ? tekst.slice(0, -2) : tekst, tablica };
}

/** Nazwa typu tresci zadania komendy: "session.create" -> "SessionCreateRequest" */
export function nazwaZadania(typKomendy) {
  return `${naPascal(typKomendy)}Request`;
}

/** Nazwa typu tresci wyniku komendy: "session.create" -> "SessionCreateResponse" */
export function nazwaWyniku(typKomendy) {
  return `${naPascal(typKomendy)}Response`;
}

/** Nazwa typu tresci zdarzenia: "session.changed" -> "SessionChangedEvent" */
export function nazwaZdarzenia(typZdarzenia) {
  return `${naPascal(typZdarzenia)}Event`;
}

/** Wspolna tresc wszystkich zdarzen *.unknown */
export const NAZWA_TRESCI_NIEZNANEJ = 'UnknownCommandPayload';

/** Nazwa stalej Go dla wartosci wyliczenia: ("ChunkKind", "tool_use") -> "ChunkKindToolUse" */
export function stalaWyliczenia(nazwaTypu, wartosc) {
  return `${nazwaTypu}${naPascal(wartosc)}`;
}
