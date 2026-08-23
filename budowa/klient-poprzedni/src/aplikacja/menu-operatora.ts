import { utworzMenuRozwijane, type MenuRozwijane } from '../okna-rownolegle/menu-rozwijane';
import { POZYCJE_USTAWIEN, type KodUstawienia, type PozycjaUstawienia } from '../strona-glowna/pozycje-ustawien';

/**
 * Menu Operatora w pasku aplikacji — tożsamość i czynności jej dotyczące.
 *
 * Nie ma tu nazwy, adresu, awatara ani przełączania kont, bo platforma kont
 * nie prowadzi: `auth.register` nie zakłada konta i nie przyjmuje adresu
 * e-mail. Nie ma też pozycji „Wyloguj" — rdzeń nie odcina komend po wygaśnięciu
 * sesji, więc taka pozycja mogłaby jedynie skasować zapis sesji i przeładować
 * aplikację.
 *
 * Zostają dwie pozycje, które dotyczą Operatora i mają dokąd prowadzić:
 * „Ustawienia" (hasło, metody wejścia) oraz „Punkty izolacji" (zakres pracy).
 * Obie pochodzą z tego samego wykazu, co listwa strony głównej, i idą tą samą
 * drogą skutku.
 */
export interface MenuOperatora {
  element: HTMLElement;
  /** Zdejmuje nasłuchy dokumentu. Obowiązkowe przy zejściu gospodarza. */
  zamknij(): void;
}

export interface OpcjeMenuOperatora {
  naPozycje(pozycja: PozycjaUstawienia): void;
}

/** Pozycje wykazu, które dotyczą Operatora, a nie platformy ani modeli. */
const KODY_OPERATORA: readonly KodUstawienia[] = ['ustawienia', 'punkty-izolacji'];

export function utworzMenuOperatora(opcje: OpcjeMenuOperatora): MenuOperatora {
  const menu: MenuRozwijane = utworzMenuRozwijane({
    ikona: 'uzytkownik',
    etykieta: 'Operator',
  });

  const pozycje = POZYCJE_USTAWIEN.filter((pozycja) => KODY_OPERATORA.includes(pozycja.kod));

  menu.ustawTresc([nagl(), ...pozycje.map((pozycja) => wpis(pozycja, opcje.naPozycje, menu)), granica()]);

  return {
    element: menu.element,
    zamknij: () => menu.zamknij(),
  };
}

/** Nagłówek menu — nazywa rolę zalogowanego, jedyną, jaką platforma zna. */
function nagl(): HTMLElement {
  const element = document.createElement('p');
  element.className = 'dn-menu__naglowek';
  element.textContent = 'Operator';
  return element;
}

/**
 * Zdanie zamykające menu. Stoi na ekranie, a nie tylko w komentarzu, żeby
 * szukający profilu osobowego dostał odpowiedź: czego nie ma, dlaczego
 * i gdzie znajdzie to, co istnieje.
 */
function granica(): HTMLElement {
  const element = document.createElement('p');
  element.className = 'dn-menu__granica';
  element.textContent =
    'Platforma nie prowadzi kont ani profili osobowych — bramka jest jedna, ' +
    'a Operator bezimienny. Hasło i metody wejścia zmienisz w Ustawieniach.';
  return element;
}

function wpis(
  pozycja: PozycjaUstawienia,
  naPozycje: (pozycja: PozycjaUstawienia) => void,
  menu: MenuRozwijane,
): HTMLElement {
  const element = document.createElement('button');
  element.type = 'button';
  element.className = 'dn-menu__pozycja';
  element.setAttribute('role', 'menuitem');
  element.dataset['ustawienie'] = pozycja.kod;
  element.title = pozycja.wyjasnienie;
  element.textContent = pozycja.nazwa;
  element.addEventListener('click', () => {
    menu.ustaw(false);
    naPozycje(pozycja);
  });
  return element;
}
