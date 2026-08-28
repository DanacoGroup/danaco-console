/**
 * Katalog tras aplikacji ustala, ile widoków najwyższego rzędu ma aplikacja
 * i jak się nazywają. Nie buduje elementów i nie przełącza widoków —
 * przełączaniem zajmuje się `router.ts`.
 */

/**
 * Nazwy tras stanowią jedyne źródło prawdy dla routera oraz dla każdego
 * przycisku otwierającego widok. Nazwa trasy jest zarazem wartością zapisywaną
 * w adresie dokumentu, więc odświeżenie strony wraca do tego samego widoku.
 */
export const Trasa = {
  StronaGlowna: 'strona-glowna',
  Srodowisko: 'srodowisko',
  Pulpit: 'pulpit',
} as const;

/**
 * Nazwa trasy jako typ zawężony do wartości katalogu. Wartość spoza katalogu nie
 * przejdzie kontroli typów, więc router i przyciski operują wyłącznie na trasach
 * istniejących.
 */
export type Trasa = (typeof Trasa)[keyof typeof Trasa];

/**
 * Trasa otwierana przy uruchomieniu aplikacji. Jest nią Centrum dowodzenia,
 * a nie okno pracy: wejście do produktu prowadzi przez widok zbiorczy, z którego
 * Operator wybiera środowisko.
 */
export const TRASA_POCZATKOWA: Trasa = Trasa.StronaGlowna;

/**
 * Wszystkie trasy ułożone w kolejności przepływu pracy. Kolejność wiąże rozpoznanie
 * nazwy trasy oraz porządek widoków najwyższego rzędu w interfejsie.
 */
export const TRASY: readonly Trasa[] = [Trasa.StronaGlowna, Trasa.Srodowisko, Trasa.Pulpit];

/**
 * Nazwy widoków pokazywane w interfejsie. Nazwa techniczna trasy pozostaje
 * wartością adresu, a nazwa z tego odwzorowania trafia do nagłówków i przycisków.
 */
export const NAZWY_TRAS: Readonly<Record<Trasa, string>> = {
  [Trasa.StronaGlowna]: 'Centrum dowodzenia',
  [Trasa.Srodowisko]: 'Środowisko',
  [Trasa.Pulpit]: 'Mission Control',
};

/**
 * Rozstrzyga, czy dowolny tekst jest nazwą trasy.
 *
 * Adres z nieznaną nazwą nie zatrzymuje uruchomienia — router sprowadza go do
 * trasy początkowej.
 */
export function czyTrasa(tekst: string): tekst is Trasa {
  return (TRASY as readonly string[]).includes(tekst);
}
