import { przycisk } from '../../modele/kontrolki-formularza';

/**
 * Przycisk czynności modułu Browser — jeden kształt dla paneli i wierszy
 * wykazu.
 *
 * Panele pomocnicze i wiersze wykazu potrzebują tego samego: zbudować przycisk
 * i podpiąć nasłuch. Kopia tej czynności stoi tu jedna, żeby cztery pliki
 * modułu nie rozjechały się przy pierwszej poprawce.
 *
 * Nazwa funkcji niesie słowo „przycisk", bo po nim rozpoznają kontrolki
 * narzędzia zestawiające etykiety czynności okien; opakowanie nazwane inaczej
 * chowa etykietę przed takim zestawieniem.
 */

/**
 * Trzy odmiany przycisku, których moduł używa — nazwa zamiast napisu klasy
 * powtarzanego przy każdym wywołaniu.
 *
 * Odmiana mówi o randze czynności, nie o barwie: `glowny` stoi tam, gdzie panel
 * ma jedną czynność wiodącą (zapis, dodanie), `zarys` przy pozostałych
 * czynnościach panelu, `duch` przy pozycjach wykazu, gdzie przycisków jest wiele
 * w jednym wierszu.
 */
export const KLASA_PRZYCISKU = {
  glowny: 'dn-btn dn-btn--sm dn-btn--atrament',
  zarys: 'dn-btn dn-btn--sm dn-btn--zarys',
  duch: 'dn-btn dn-btn--sm dn-btn--duch',
} as const;

/** Przycisk czynności wraz z podpiętym naciśnięciem. */
export function przyciskCzynnosci(
  nazwa: string,
  klasa: string,
  naNacisniecie: () => void,
): HTMLButtonElement {
  const element = przycisk(nazwa, klasa);
  element.addEventListener('click', naNacisniecie);
  return element;
}
