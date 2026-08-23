/**
 * Plakietka znaku w wierszu wykazu sekcji modeli — wspólna dla wykazu kont
 * i wykazu kategorii tożsamości.
 *
 * Klasa biblioteki mówi o wyglądzie znaku (`dn-plakietka` i jej odmiany), klasa
 * miejsca — o jego położeniu w wierszu; to dwie różne decyzje i stoją w dwóch
 * parametrach. Plik nie wymienia żadnej barwy.
 */
export function znakWykazu(tresc: string, klasa: string, klasaMiejsca: string): HTMLElement {
  const element = document.createElement('span');
  element.className = `${klasa} ${klasaMiejsca}`;
  element.textContent = tresc;
  return element;
}
