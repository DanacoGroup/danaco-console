import './panele.css';

import { utworzNaglowekPanelu } from './naglowek-panelu';

/**
 * Widok pełnoekranowy — trzeci rodzaj obszaru sceny okien równoległych.
 *
 * Scena dzieli szerokość między trzy byty (`rodzaje-obszaru.ts`): kolumnę
 * rozmowy, kolumnę paneli i widok pełnoekranowy. Ten trzeci jako jedyny nie ma
 * sąsiada — bierze całą scenę, bo panel w kolumnie trzyma minimum szerokości,
 * a treść, która się w nim nie mieści (różnica plikowa, drzewo, renderowana
 * strona), potrzebuje sceny, nie kolumny.
 *
 * Treść jest przenoszona, nie kopiowana. Drugi egzemplarz panelu znaczyłby
 * drugą subskrypcję rdzenia (panel podglądu bash trzyma `stream.chunk`),
 * a przy zamknięciu jednego z nich — subskrypcję osieroconą. Dlatego `pokaz`
 * zapamiętuje miejsce, z którego treść zabrał (rodzic i następnik), a `ukryj`
 * odkłada ją dokładnie tam. Bez zapamiętanego następnika panel wracałby na
 * koniec stosu i kolejność paneli zmieniałaby się po każdym wyjściu z pełnego
 * ekranu.
 *
 * Klawisz `Escape` wychodzi — wzorzec zwijania z
 * `widok-sterowania/szuflada.ts`: nasłuch `keydown` stoi na elemencie widoku,
 * nie na dokumencie. Wynikają z tego dwie rzeczy, obie tu potrzebne: sprzątanie
 * jest zbędne, bo nasłuch ginie razem z elementem, a klawisz działa tylko
 * wtedy, gdy ognisko jest wewnątrz widoku — nie zabiera `Escape` niczemu innemu
 * na scenie. Dlatego `pokaz` przenosi ognisko na widok, tak jak szuflada oddaje
 * je uchwytowi przy zwinięciu.
 *
 * Wyjście nie jest zamknięciem: oddaje treść z powrotem do stosu, a panel żyje
 * dalej ze swoją subskrypcją. Zamknięcie panelu należy do jego obudowy w stosie
 * i tam zostaje.
 *
 * Widok nie buduje treści, nie zna paneli po kodzie, nie zna rejestru ani
 * wytwórni. Nie rozstrzyga też, co robi scena, gdy się pokazuje — chowanie
 * kolumn należy do układu, który ten widok osadził.
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

/** Miejsce, z którego treść została zabrana — żeby wróciła dokładnie tam. */
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
      // Odłożenie przed zapamiętanym następnikiem; `insertBefore` z `null`
      // znaczy „na koniec" — czyli dokładnie stan, w którym treść była ostatnia.
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
      // Pokazanie drugiej treści bez wyjścia oddaje pierwszą na swoje miejsce;
      // inaczej poprzedni panel zostałby bez wnętrza i bez śladu, gdzie stał.
      if (trzymana !== null) ukryj();

      // Treść bez rodzica nie ma dokąd wrócić — widok tego nie zmyśla i nie
      // podstawia własnego gniazda, bo odłożyłby ją do siebie samego.
      const rodzic = tresc.parentNode;
      pochodzenie = rodzic === null ? null : { rodzic, nastepnik: tresc.nextSibling };
      trzymana = tresc;

      const nowy = utworzNaglowekPanelu(tytul, [
        {
          // Ikona pod czynność, nie pod wygląd narzędzia: czynnością jest
          // zamknięcie tego obszaru sceny.
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
