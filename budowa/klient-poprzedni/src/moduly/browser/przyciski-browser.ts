/**
 * Budowa przycisków czynności modułu Browser: jeden kształt przycisku dla
 * paneli pomocniczych i dla wierszy wykazu, wraz z nazwami trzech odmian
 * niosących rangę czynności.
 */
import { przycisk } from '../../modele/kontrolki-formularza';

/**
 * Trzy odmiany przycisku używane w module. Odmiana mówi o randze czynności,
 * nie o barwie: `glowny` stoi przy czynności wiodącej panelu, `zarys` przy
 * pozostałych czynnościach panelu, `duch` przy pozycjach wykazu.
 */
export const KLASA_PRZYCISKU = {
  glowny: 'dn-btn dn-btn--sm dn-btn--atrament',
  zarys: 'dn-btn dn-btn--sm dn-btn--zarys',
  duch: 'dn-btn dn-btn--sm dn-btn--duch',
} as const;

/**
 * Przycisk czynności wraz z podpiętym naciśnięciem. Nazwa funkcji niesie wyraz
 * „przycisk", po którym rozpoznają go narzędzia zestawiające etykiety czynności
 * okien.
 */
export function przyciskCzynnosci(
  nazwa: string,
  klasa: string,
  naNacisniecie: () => void,
): HTMLButtonElement {
  const element = przycisk(nazwa, klasa);
  element.addEventListener('click', naNacisniecie);
  return element;
}
