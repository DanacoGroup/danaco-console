/**
 * Katalog tras aplikacji.
 *
 * Jedna odpowiedzialność: ustalenie, ile widoków najwyższego rzędu ma
 * aplikacja i jak się nazywają. Bez elementów, bez logiki przełączania —
 * przełączaniem zajmuje się `router.ts`.
 *
 * Trzy trasy:
 *
 *   strona-glowna → Centrum dowodzenia; wejście do produktu
 *   srodowisko    → powłoka środowiska z kartami sesji i modułami
 *   pulpit        → Mission Control jako widok obok strony głównej
 *
 * Nazwa trasy jest zarazem wartością zapisywaną w adresie dokumentu, więc
 * odświeżenie strony wraca do tego samego widoku, a nie na początek.
 */

/** Nazwy tras — jedno źródło prawdy dla routera i wszystkich przycisków. */
export const Trasa = {
  StronaGlowna: 'strona-glowna',
  Srodowisko: 'srodowisko',
  Pulpit: 'pulpit',
} as const;

/** Nazwa trasy. */
export type Trasa = (typeof Trasa)[keyof typeof Trasa];

/** Trasa otwierana przy uruchomieniu — Centrum dowodzenia, nie okno pracy. */
export const TRASA_POCZATKOWA: Trasa = Trasa.StronaGlowna;

/** Wszystkie trasy w kolejności przepływu. */
export const TRASY: readonly Trasa[] = [Trasa.StronaGlowna, Trasa.Srodowisko, Trasa.Pulpit];

/** Nazwa widoku pokazywana Operatorowi. */
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
