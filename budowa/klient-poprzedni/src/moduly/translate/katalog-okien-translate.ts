/**
 * Kod modułu i kody jego okien operacyjnych — jedno miejsce na oznaczenia,
 * którymi moduł zgłasza się rdzeniowi.
 *
 * Kody nie są wymyślone: `translate` oraz `source-panel`, `translation-panels`
 * i `glossary-manager` pochodzą z rejestru okien operacyjnych rdzenia. Trzy
 * pozostałe okna opracowania — pamięć tłumaczeń, studio formatów i centrum
 * kontroli jakości — w rejestrze rdzenia nie występują, więc ich kody powstają
 * tutaj wedle tej samej reguły: nazwa własna okna zapisana małymi literami
 * z łącznikami.
 *
 * Rozjazd nie jest przemilczany. Moduł zgłasza katalogowi okien komplet
 * sześciu kodów, a byt wspólny (`moduly/katalog-okien.ts`) wypowiada obie
 * strony różnicy: okna rejestru, których moduł nie buduje, oraz okna budowane
 * spoza rejestru. Dopisanie trzech okien do rejestru rdzenia zdejmie tę drugą
 * połowę bez zmiany ani jednej linii tutaj.
 */

/** Kod modułu w rejestrze rdzenia. */
export const KOD_MODULU = 'translate';

/** Kody okien operacyjnych budowanych przez moduł. */
export const KODY_OKIEN = {
  zrodlo: 'source-panel',
  panele: 'translation-panels',
  glosariusz: 'glossary-manager',
  pamiec: 'translation-memory-panel',
  formaty: 'format-studio',
  jakosc: 'qa-review-center',
  warsztat: 'translation-workshop',
} as const;
