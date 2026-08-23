import { elementIkony } from '../ikony/ikony';

/**
 * Węzły jednej karty sesji — sama budowa i etykiety czynności.
 *
 * Jedna odpowiedzialność: kształt karty. Stan, wybór i cykl życia zostają
 * w `karta-sesji.ts`; te funkcje są czyste i nie domykają się na stanie karty,
 * więc wytwórnia karty zostaje krótka.
 */

/** Węzły jednej karty; złożone osobno, żeby wytwórnia została krótka. */
export interface CzesciKarty {
  element: HTMLElement;
  kropka: HTMLElement;
  napis: HTMLElement;
  plakietka: HTMLElement;
  usun: HTMLButtonElement;
  zamknij: HTMLButtonElement;
}

/** Przycisk ikonowy karty — obie czynności karty mają ten sam kształt. */
function przyciskKarty(ikona: 'zamknij' | 'kosz', klasa: string): HTMLButtonElement {
  const przycisk = document.createElement('button');
  przycisk.type = 'button';
  przycisk.className = `dn-btn-ikona ${klasa}`;
  przycisk.append(elementIkony(ikona, { rozmiar: 16 }));
  return przycisk;
}

/**
 * Budowa węzłów karty.
 *
 * Karta jest zakładką z czynnościami w środku, więc nośnikiem roli jest element
 * bez własnej semantyki — przycisk w przycisku byłby zapisem wadliwym.
 *
 * Usunięcie trwałe stoi obok zamknięcia, nie zamiast niego: kosz mówi
 * o utracie zapisu, krzyżyk o zamknięciu karty. Jedna kontrolka dla obu
 * czynności kazałaby Operatorowi zgadywać, którą właśnie wykonał.
 */
export function zlozCzesciKarty(id: string): CzesciKarty {
  const element = document.createElement('div');
  element.className = 'dn-zakladka dn-sesja';
  element.setAttribute('role', 'tab');
  element.setAttribute('aria-selected', 'false');
  element.tabIndex = -1;
  element.dataset.karta = id;

  const kropka = document.createElement('span');
  kropka.setAttribute('role', 'img');

  const napis = document.createElement('span');
  napis.className = 'dn-sesja__tytul';

  const plakietka = document.createElement('span');
  plakietka.className = 'dn-plakietka dn-plakietka--rola dn-sesja__stan';

  const usun = przyciskKarty('kosz', 'dn-sesja__usun');
  const zamknij = przyciskKarty('zamknij', 'dn-sesja__zamknij');

  element.append(kropka, napis, plakietka, usun, zamknij);
  return { element, kropka, napis, plakietka, usun, zamknij };
}

/**
 * Etykiety obu czynności karty. Etykieta nazywa skutek, nie gest: „usuń” bez
 * „trwale” brzmi jak zamknięcie, a chodzi o utratę zapisu.
 */
export function ubierzCzynnosciKarty(czesci: CzesciKarty, tytul: string): void {
  czesci.zamknij.setAttribute('aria-label', `Zamknij kartę sesji: ${tytul}`);
  czesci.zamknij.title = `Zamknij kartę sesji: ${tytul}`;
  const opisUsuniecia = `Usuń trwale sesję: ${tytul} — zapis ginie bezpowrotnie`;
  czesci.usun.setAttribute('aria-label', opisUsuniecia);
  czesci.usun.title = opisUsuniecia;
}
