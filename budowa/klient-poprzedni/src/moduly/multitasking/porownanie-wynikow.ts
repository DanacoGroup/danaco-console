import type { Message, Window } from '../../../../shared/contract';

/**
 * Wynik jednego wykonawcy przygotowany do zestawienia. Numer wykonawcy jest
 * jego miejscem na scenie, a nie identyfikatorem okna; okno oraz ostatnia
 * domknięta odpowiedź pochodzą wprost z kontraktu.
 */
export interface WynikWykonawcy {
  /** Numer wykonawcy na scenie: 1 albo 2. */
  numer: number;
  /** Okno wykonawcy. */
  okno: Window;
  /** Ostatnia domknięta odpowiedź; pusta, gdy wykonawca nie zwrócił jeszcze nic. */
  wiadomosc: Message | null;
}

/**
 * Rozbieżność wyników wykonawców: prawda wtedy, gdy co najmniej dwa niepuste
 * wyniki różnią się treścią. Pojedynczy wynik ani sama pustka rozbieżnością
 * nie są — nie ma wtedy konfliktu do rozstrzygnięcia.
 */
export function czyRozbiezne(wyniki: readonly WynikWykonawcy[]): boolean {
  const tresci = wyniki
    .map((pozycja) => pozycja.wiadomosc?.content ?? '')
    .filter((tresc) => tresc !== '');
  if (tresci.length < 2) return false;
  return tresci.some((tresc) => tresc !== tresci[0]);
}

/**
 * Kolumna wyniku jednego wykonawcy: tytuł z numerem wykonawcy, plakietka stanu
 * wiadomości oraz treść odpowiedzi. Brak wyniku daje zdanie mówiące o tym
 * wprost, zamiast pustej kolumny.
 */
export function kolumnaWyniku(wynik: WynikWykonawcy): HTMLElement {
  const element = document.createElement('article');
  element.className = 'dm-kolumna';
  element.dataset['wykonawca'] = String(wynik.numer);
  element.dataset['okno'] = wynik.okno.id;

  const naglowek = document.createElement('h3');
  naglowek.className = 'dm-kolumna__tytul';
  naglowek.textContent = `Executor ${wynik.numer}`;

  const stan = document.createElement('span');
  stan.className = 'dn-plakietka';
  stan.textContent = wynik.wiadomosc?.status ?? 'bez wyniku';

  const tresc = document.createElement('pre');
  tresc.className = 'dm-kolumna__tresc';
  tresc.textContent =
    wynik.wiadomosc === null || wynik.wiadomosc.content === ''
      ? 'Wykonawca nie zwrócił jeszcze wyniku.'
      : wynik.wiadomosc.content;

  element.append(naglowek, stan, tresc);
  return element;
}

/**
 * Zestawienie wyników: kolumna na każdego wykonawcę, a przy pustej obsadzie
 * zdanie o braku wyników do oceny. Atrybut `data-rozbiezne` niesie wynik
 * porównania treści.
 */
export function zestawienie(wyniki: readonly WynikWykonawcy[]): HTMLElement {
  const element = document.createElement('div');
  element.className = 'dm-zestawienie';
  element.dataset['rozbiezne'] = String(czyRozbiezne(wyniki));

  if (wyniki.length === 0) {
    const pusto = document.createElement('p');
    pusto.className = 'dm-stan';
    pusto.dataset['stan'] = 'pusto';
    pusto.textContent = 'Obsada bez wykonawców — nie ma wyników do oceny.';
    element.append(pusto);
    return element;
  }

  for (const wynik of wyniki) element.append(kolumnaWyniku(wynik));
  return element;
}
