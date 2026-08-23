// Naglowek pliku generowanego: jedna tresc ostrzezenia dla obu jezykow.

/** Linie naglowka bez znakow komentarza — emiter dokleja wlasne. */
export function linieNaglowka(kontrakt, nazwaArtefaktu) {
  return [
    `${kontrakt.produkt} — ${nazwaArtefaktu}`,
    '',
    'PLIK GENEROWANY. Nie edytuj recznie — zmiany przepadna przy nastepnej generacji.',
    '',
    'Zrodlo prawdy: shared/contract.json',
    'Generacja:     node shared/gen/generate.mjs',
    '',
    `Protokol: ${kontrakt.protokol} · Transport: ${kontrakt.transport}`,
    `Notacja:  ${kontrakt.notacja}`,
  ];
}
