/**
 * Nagłówek okna operacyjnego — górny pas okna: nazwa, rola, kontrolki.
 *
 * Jeden kształt pasa dla wszystkich modułów: Operator dostaje tę samą odpowiedź
 * na pytanie „na które okno patrzę", niezależnie od modułu.
 *
 * Granica wobec ramy okna (`komponenty/rama-okna.ts`): rama to cały kontener
 * okna wraz z ciałem i stanami, nagłówek to wyłącznie jego górny pas i nie wie
 * nic o treści pod sobą. Nagłówek nie osadza się sam — dostaje go rama albo
 * plik okna.
 *
 * Rola stoi przy nazwie, bo nazwy okien powtarzają się między modułami —
 * Process Monitor (Terminal) i Execution Monitor (Automations) to odrębne
 * okna. Podpis roli jest tym, co je na ekranie rozróżnia.
 *
 * Wygląd pochodzi w całości z biblioteki: pas stoi na `dn-karta-naglowek`,
 * nazwa na `dn-karta-tytul`, rola na `dn-plakietka--rola` — bez klas własnych
 * i bez barw w kodzie. Moduł, który potrzebuje odstępstwa, podaje własną klasę
 * polem `klasa`, zamiast powielać komponent.
 */

/** Opis nagłówka. Pola opcjonalne pominięte = elementu nie ma w drzewie DOM. */
export interface OpisNaglowkaOkna {
  /** Nazwa okna operacyjnego — pierwszy element pasa. */
  tytul: string;
  /**
   * Rola okna widoczna plakietką obok nazwy. Napis dowolny, bo moduły niosą
   * tu i samo określenie roli („wiodące"), i zdanie doprecyzowujące
   * („monitor · zlecenia wieloetapowe asystenta na żywo").
   */
  rola?: string;
  /**
   * Elementy sterowania osadzane w pasie, w podanej kolejności, za plakietką.
   * Kontrolki przychodzą gotowe — nagłówek ich nie buduje i nie zna ich stanu.
   * Wyrównanie do prawej krawędzi robi wołający, opakowując je
   * w `dn-pasek-prawa`.
   */
  kontrolki?: readonly HTMLElement[];
  /**
   * Klasa modułu dopisywana za klasą biblioteczną (np. `mb-okno__naglowek`).
   * Miejsce na odstępstwo modułu, nie na drugi układ pasa; nazwę wewnątrz
   * pasa moduł styluje selektorem potomka.
   */
  klasa?: string;
}

/** Pas nagłówka gotowy do osadzenia w ramie okna albo w pliku okna. */
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
