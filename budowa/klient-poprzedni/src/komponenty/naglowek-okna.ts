/**
 * Nagłówek okna operacyjnego, czyli górny pas z nazwą okna, plakietką roli
 * i kontrolkami. Jeden kształt pasa obowiązuje we wszystkich modułach.
 * Osadzenie pasa należy do ramy okna albo do pliku okna, a cały wygląd
 * pochodzi z klas biblioteki.
 */

/**
 * Opis nagłówka okna. Pominięcie pola opcjonalnego oznacza, że odpowiadający
 * mu element nie powstaje w drzewie dokumentu, więc kształt pasa wynika
 * wyłącznie z pól podanych przez moduł wołający.
 */
export interface OpisNaglowkaOkna {
  /** Nazwa okna operacyjnego — pierwszy element pasa. */
  tytul: string;
  /** Rola okna widoczna plakietką obok nazwy; napis jest dowolny. */
  rola?: string;
  /** Gotowe elementy sterowania osadzane w pasie w podanej kolejności, za plakietką roli. */
  kontrolki?: readonly HTMLElement[];
  /** Klasa modułu dopisywana za klasą biblioteczną, przeznaczona na odstępstwo wyglądu. */
  klasa?: string;
}

/**
 * Buduje pas nagłówka jako element header z klasą biblioteczną, nazwą okna
 * w nagłówku trzeciego stopnia, opcjonalną plakietką roli oraz kontrolkami
 * dopisanymi na końcu pasa. Zwrócony element jest gotowy do osadzenia.
 */
export function utworzNaglowekOkna(opis: OpisNaglowkaOkna): HTMLElement {
  const element = document.createElement('header');
  element.className = 'dn-karta-naglowek';
  if (opis.klasa !== undefined && opis.klasa !== '') element.classList.add(opis.klasa);

  const nazwa = document.createElement('h3');
  nazwa.className = 'dn-karta-tytul';
  nazwa.textContent = opis.tytul;
  element.append(nazwa);

  if (opis.rola !== undefined && opis.rola !== '') {
    const plakietka = document.createElement('span');
    plakietka.className = 'dn-plakietka dn-plakietka--rola';
    plakietka.textContent = opis.rola;
    element.append(plakietka);
  }

  element.append(...(opis.kontrolki ?? []));
  return element;
}
