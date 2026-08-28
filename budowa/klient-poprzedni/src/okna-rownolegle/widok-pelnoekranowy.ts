import './panele.css';

import { utworzNaglowekPanelu } from './naglowek-panelu';

/**
 * Widok pełnoekranowy: trzeci rodzaj obszaru sceny okien równoległych, zajmujący całą
 * scenę bez sąsiada, przenoszący treść panelu tymczasowo i oddający ją klawiszem Escape.
 */
export interface WidokPelnoekranowy {
  /** Element osadzany na scenie; chowa się sam, gdy nic nie pokazuje. */
  element: HTMLElement;
  /** Przejmuje treść panelu na całą scenę. */
  pokaz(tytul: string, tresc: HTMLElement): void;
  /** Oddaje treść z powrotem wołającemu i chowa się. */
  ukryj(): void;
  widoczny(): boolean;
}

/**
 * Miejsce, z którego treść panelu została zabrana na czas widoku pełnoekranowego,
 * zapamiętane po to, żeby po wyjściu wróciła dokładnie tam, skąd pochodzi.
 */
interface MiejscePochodzenia {
  rodzic: ParentNode;
  /** Węzeł, przed którym treść stała. `null` = stała na końcu. */
  nastepnik: ChildNode | null;
}

export function utworzWidokPelnoekranowy(naWyjscie: () => void): WidokPelnoekranowy {
  let pochodzenie: MiejscePochodzenia | null = null;
  let trzymana: HTMLElement | null = null;

  const gniazdo = document.createElement('div');
  gniazdo.className = 'dn-okna__pelny-ekran-tresc';

  const element = document.createElement('section');
  element.className = 'dn-okna__pelny-ekran';
  element.hidden = true;
  // Ognisko z kodu — inaczej nasłuch `Escape` na elemencie nigdy by nie usłyszał.
  element.tabIndex = -1;

  let naglowek = utworzNaglowekPanelu('', []);
  element.append(naglowek, gniazdo);

  function ukryj(): void {
    if (trzymana !== null && pochodzenie !== null) {
      // Odłożenie przed zapamiętanym następnikiem; brak następnika znaczy koniec listy.
      pochodzenie.rodzic.insertBefore(trzymana, pochodzenie.nastepnik);
    }
    trzymana = null;
    pochodzenie = null;
    element.hidden = true;
  }

  /** Wyjście Operatora: treść wraca, potem gospodarz dowiaduje się o zmianie. */
  function wyjdz(): void {
    ukryj();
    naWyjscie();
  }

  element.addEventListener('keydown', (zdarzenie) => {
    if (zdarzenie.key !== 'Escape' || element.hidden) return;
    wyjdz();
  });

  return {
    element,
    pokaz(tytul, tresc) {
      // Pokazanie drugiej treści bez wyjścia oddaje pierwszą na swoje miejsce, żeby nie zniknęła bez śladu.
      if (trzymana !== null) ukryj();

      // Treść bez rodzica nie ma dokąd wrócić — widok tego nie zmyśla i nie podstawia własnego gniazda.
      const rodzic = tresc.parentNode;
      pochodzenie = rodzic === null ? null : { rodzic, nastepnik: tresc.nextSibling };
      trzymana = tresc;

      const nowy = utworzNaglowekPanelu(tytul, [
        {
          // Ikona pod czynność: czynnością jest zamknięcie tego obszaru sceny.
          ikona: 'zamknij',
          etykieta: `Zamknij widok pełnoekranowy panelu ${tytul} — klawisz Escape`,
          dzialanie: wyjdz,
        },
      ]);
      element.replaceChild(nowy, naglowek);
      naglowek = nowy;

      gniazdo.replaceChildren(tresc);
      element.setAttribute('aria-label', `Panel ${tytul} na całej scenie`);
      element.hidden = false;
      element.focus();
    },
    ukryj,
    widoczny: () => !element.hidden,
  };
}
