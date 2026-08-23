/**
 * Drobne części układu nakładki — wykaz pól i nagłówek sekcji.
 *
 * Cztery sekcje okna (stan, podpowiedzi, obecność, kontekst) wypisują pola tym
 * samym wykazem `dl`; wspólne miejsce trzyma pętlę `dt`/`dd` w jednym egzemplarzu.
 *
 * Plik nie zna barw ani odstępów — nadaje wyłącznie klasy z przedrostkiem `ao-`,
 * które pokrywa `aod.css` żetonami `--dn-*`.
 */

/** Zdanie wypisywane, gdy pole jest puste. Pustka bywa poprawna. */
export const BRAK = '(brak)';

/** Buduje pusty wykaz pól. */
export function utworzWykazPol(): HTMLDListElement {
  const lista = document.createElement('dl');
  lista.className = 'ao-pola';
  return lista;
}

/** Dokłada do wykazu jedno pole: etykietę i wartość. */
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

/** Buduje podtytuł sekcji okna. */
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
