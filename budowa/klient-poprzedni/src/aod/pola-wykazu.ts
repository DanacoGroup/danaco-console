/**
 * Drobne części układu nakładki — wykaz pól i nagłówek sekcji. Cztery sekcje
 * okna wypisują pola tym samym wykazem `dl`, a plik nadaje wyłącznie klasy
 * z przedrostkiem `ao-`, które pokrywa `aod.css` żetonami `--dn-*`.
 */

/**
 * Zdanie wypisywane w miejsce wartości, gdy pole jest puste; pustka bywa stanem
 * poprawnym, więc wykaz nazywa ją wprost zamiast zostawiać puste miejsce.
 */
export const BRAK = '(brak)';

/**
 * Buduje pusty wykaz pól: element `dl` z klasą `ao-pola`, do którego kolejne
 * pola dokłada `dodajPole`. Sekcja okna zaczyna się od tego wykazu.
 */
export function utworzWykazPol(): HTMLDListElement {
  const lista = document.createElement('dl');
  lista.className = 'ao-pola';
  return lista;
}

/**
 * Dokłada do wykazu jedno pole: element `dt` z etykietą i element `dd`
 * z wartością. Wartość pusta ustępuje miejsca zdaniu o braku, żeby wiersz nie
 * został bez treści.
 */
export function dodajPole(lista: HTMLDListElement, etykieta: string, wartosc: string): void {
  const dt = document.createElement('dt');
  dt.textContent = etykieta;
  const dd = document.createElement('dd');
  dd.textContent = wartosc === '' ? BRAK : wartosc;
  lista.append(dt, dd);
}

/**
 * Dokłada pole niosące wykaz identyfikatorów.
 *
 * Wykaz pusty ORAZ wykaz nieobecny znaczą dla Operatora to samo — rdzeń nic
 * nie przysłał — i oba wypisujemy jako `(brak)`, nigdy jako `0` czy `[]`.
 */
export function dodajPoleWykazu(
  lista: HTMLDListElement,
  etykieta: string,
  wartosci?: readonly string[],
): void {
  const pozycje = wartosci ?? [];
  dodajPole(lista, etykieta, pozycje.length === 0 ? BRAK : pozycje.join(', '));
}

/**
 * Buduje podtytuł sekcji okna: nagłówek czwartego rzędu z klasą `ao-podtytul`,
 * którą pokrywa `aod.css`. Podtytuł dzieli nakładkę na cztery sekcje.
 */
export function utworzPodtytul(tekst: string): HTMLElement {
  const naglowek = document.createElement('h4');
  naglowek.className = 'ao-podtytul';
  naglowek.textContent = tekst;
  return naglowek;
}

/**
 * Buduje akapit zdania stanu wewnątrz sekcji.
 *
 * Wykaz klas jest zamknięty: nazwa podana napisem swobodnym przeszłaby
 * sprawdzenie typów także wtedy, gdy arkusz nie ma dla niej reguły. Cztery
 * zdania, cztery klasy, wszystkie z regułą w `aod.css`.
 */
export function utworzAkapit(
  klasa: 'ao-odmowa' | 'ao-pusto' | 'ao-granica' | 'ao-pokwitowanie',
  tekst: string,
): HTMLParagraphElement {
  const akapit = document.createElement('p');
  akapit.className = klasa;
  akapit.textContent = tekst;
  return akapit;
}
