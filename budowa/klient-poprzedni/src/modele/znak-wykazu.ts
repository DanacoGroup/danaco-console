/**
 * Plakietka znaku w wierszu wykazu sekcji modeli — wspólna dla wykazu kont
 * i wykazu kategorii tożsamości.
 *
 * Klasa biblioteki mówi o wyglądzie znaku, klasa miejsca — o jego położeniu
 * w wierszu; to dwie różne decyzje i stoją w dwóch osobnych parametrach.
 */
export function znakWykazu(tresc: string, klasa: string, klasaMiejsca: string): HTMLElement {
  const element = document.createElement('span');
  element.className = `${klasa} ${klasaMiejsca}`;
  element.textContent = tresc;
  return element;
}
