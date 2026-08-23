/**
 * Powtarzalne elementy kart sekcji dostępów: przycisk czynności, plakietka
 * znaku i rozpychacz nagłówka.
 *
 * Jedno miejsce budowy tych elementów utrzymuje spójny wygląd karty punktu,
 * wiersza nadania i obszaru dodawania katalogu.
 *
 * Wygląd pochodzi wyłącznie z klas biblioteki `komponenty/` i z klas
 * modyfikujących podanych przez wywołującego — plik nie zna barw, odstępów
 * ani reguł widoku.
 */

/** Przycisk czynności karty; klasy wyglądu podaje wywołujący. */
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

/** Rozpychacz nagłówka: odsuwa czynności na prawą krawędź wiersza. */
export function rozpychaczNaglowka(klasa: string): HTMLElement {
  const element = document.createElement('span');
  element.className = klasa;
  return element;
}
