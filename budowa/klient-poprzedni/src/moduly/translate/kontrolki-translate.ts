/**
 * Dwie kontrolki, których biblioteka nie ma: dymek objaśnienia [?] przy
 * elemencie konfiguracji i podpowiedź przy polu otwartym.
 *
 * Dymek jest budowany tutaj, a nie brany z biblioteki, bo nosi klasy arkusza
 * modułu (`mt-…`); wersja biblioteczna wymaga klas układu, których ten arkusz
 * nie wciąga, więc dałaby dymek bez pozycjonowania.
 *
 * Dymek pokazuje się na `:hover` albo `:focus-within`, bez klikania i bez
 * zamykania; znak jest przyciskiem, więc naciśnięcie prowadzi ognisko i daje
 * odpowiedź; treść siedzi w `aria-label`, więc dymek nie potrzebuje
 * identyfikatora i nie zderza się między oknami.
 */

export function utworzDymek(objasnienie: string): HTMLElement {
  const dymek = document.createElement('span');
  dymek.className = 'dn-tooltip mt-dymek';

  const znak = document.createElement('button');
  znak.type = 'button';
  znak.className = 'mt-dymek__znak';
  znak.textContent = '?';
  znak.setAttribute('aria-label', objasnienie);

  const tresc = document.createElement('span');
  tresc.className = 'dn-tooltip-tresc';
  tresc.setAttribute('aria-hidden', 'true');
  tresc.textContent = objasnienie;

  dymek.append(znak, tresc);
  return dymek;
}

/**
 * Dopina dymek do etykiety pola formularza.
 *
 * Dymek stoi przy etykiecie, nie pod polem: pod polem mieszka tekst pomocy
 * i tekst błędu walidacji, a trzeci wiersz w tym samym miejscu zacierałby
 * różnicę między objaśnieniem a usterką.
 */
export function dopnijDymek(pole: HTMLElement, objasnienie: string): HTMLElement {
  const etykieta = pole.querySelector('.dn-pole-etykieta');
  if (etykieta === null) pole.append(utworzDymek(objasnienie));
  else etykieta.append(utworzDymek(objasnienie));
  return pole;
}

/**
 * Podpowiedź do pola otwartego.
 *
 * Kontrakt nie ma komendy zwracającej wykaz języków ani tonów, więc lista jest
 * podpowiedzią, a pole zostaje edytowalne: wartość spoza wykazu wolno wpisać
 * wprost. `datalist` robi dokładnie to i niczego nie zamyka.
 */
export function podepnijPodpowiedz(
  kontrolka: HTMLInputElement,
  identyfikator: string,
  wartosci: readonly string[],
): HTMLDataListElement {
  const wykaz = document.createElement('datalist');
  wykaz.id = identyfikator;
  wykaz.replaceChildren(
    ...wartosci.map((wartosc) => {
      const pozycja = document.createElement('option');
      pozycja.value = wartosc;
      return pozycja;
    }),
  );
  kontrolka.setAttribute('list', identyfikator);
  return wykaz;
}

/** Nagłówek okna operacyjnego wraz z jego rolą z wykazu okien. */
export function naglowekOkna(nazwa: string, rola: string): HTMLElement {
  const tytul = document.createElement('h3');
  tytul.className = 'mt-okno__tytul';
  tytul.textContent = nazwa;

  const plakietka = document.createElement('span');
  plakietka.className = 'dn-plakietka dn-plakietka--rola mt-okno__rola';
  plakietka.textContent = rola;

  const naglowek = document.createElement('header');
  naglowek.className = 'mt-okno__naglowek';
  naglowek.append(tytul, plakietka);
  return naglowek;
}
