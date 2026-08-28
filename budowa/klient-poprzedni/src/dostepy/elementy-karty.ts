/**
 * Powtarzalne elementy kart sekcji dostępów: przycisk czynności, plakietka
 * znaku i rozpychacz nagłówka. Wygląd pochodzi wyłącznie z klas biblioteki
 * `komponenty/` oraz z klas modyfikujących podanych przez wywołującego.
 */

/**
 * Przycisk czynności karty. Element ma ustawiony `type`, więc naciśnięcie nie
 * wysyła formularza, a komplet klas wyglądu podaje wywołujący.
 */
export function przyciskKarty(tresc: string, klasa: string): HTMLButtonElement {
  const element = document.createElement('button');
  element.type = 'button';
  element.className = klasa;
  element.append(document.createTextNode(tresc));
  return element;
}

/**
 * Plakietka znaku karty. Klasa modyfikująca jest osobnym parametrem, bo mówi
 * o miejscu w układzie, a klasa biblioteki o wyglądzie — dwie różne decyzje.
 */
export function plakietkaZnaku(tresc: string, klasa: string, klasaMiejsca: string): HTMLElement {
  const element = document.createElement('span');
  element.className = `${klasa} ${klasaMiejsca}`;
  element.textContent = tresc;
  return element;
}

/**
 * Rozpychacz nagłówka karty: pusty element rozciągliwy, który odsuwa czynności
 * na prawą krawędź wiersza. Szerokość wynika z klasy podanej przez wywołującego.
 */
export function rozpychaczNaglowka(klasa: string): HTMLElement {
  const element = document.createElement('span');
  element.className = klasa;
  return element;
}
