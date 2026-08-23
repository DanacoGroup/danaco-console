import { elementIkony, type NazwaIkony } from '../ikony/ikony';
import './czynnosci.css';

/**
 * Jeden wiersz sekcji czynności sesji w menu `⋮`.
 *
 * Stoi osobno od `menu-paneli.ts`, bo wiersz paneli jest przełącznikiem: niesie
 * `aria-checked`, ptaszek i stałe miejsce na ten ptaszek, bo panel bywa otwarty
 * albo zamknięty. Czynność sesji nie ma stanu „włączona" — naciśnięcie ją
 * wykonuje, więc ptaszek, który nigdy się nie zapali, kłamałby o tym, czym
 * pozycja jest.
 *
 * Układ jest ten sam co w sekcji paneli: ikona · nazwa i przeznaczenie · `<kbd>`
 * skrótu wyrównany do prawej. Dwie sekcje jednego menu nie mają dwóch rytmów.
 *
 * Skrót jest prawdziwy, ale miejscowy. Sekcja paneli skrótów nie pokazuje, bo
 * nikt nie rejestruje klawiszy i napis byłby atrapą; tutaj napis zostaje, bo
 * klawisz działa — nasłuch stoi na sekcji i łapie literę, dopóki menu jest
 * rozwinięte (`czynnosci-sesji-menu.ts`). Skrótu globalnego nie rejestrujemy:
 * `R`, `A`, `D` wciśnięte w polu wypowiedzi mają pisać litery, a nie kasować
 * sesję.
 *
 * Wyróżnienie ostrzegawcze jest barwą wiersza, nie blokadą: „Usuń" idzie
 * czerwienią, bo kasuje zapis bez odwrotu, ale wiersz jest tak samo klikalny
 * jak każdy inny.
 */

/** Czynność gotowa do postawienia w menu. */
export interface PozycjaCzynnosciMenu {
  /** Klucz techniczny — po nim idzie wybór skrótu i porządek. */
  klucz: string;
  /** Nazwa widoczna w menu. */
  nazwa: string;
  /** Zdanie pod nazwą: co czynność naprawdę robi. */
  przeznaczenie: string;
  ikona: NazwaIkony;
  /**
   * Skrót — pojedyncza litera, na przykład `R`.
   *
   * Pominięty znaczy, że pozycja klawisza nie ma, a nie że go jeszcze nie
   * podpięto: litery niosą cztery pozycje (`R`, `F`, `A`, `D`), a `Otwórz w ›`
   * stoi bez klawisza. Wiersz bez skrótu nie dostaje pustego `<kbd>` — ramka
   * bez litery obiecywałaby klawisz, którego nikt nie nasłuchuje.
   */
  skrot?: string;
  /** Czy wiersz idzie kolorem ostrzegawczym (utrata danych bez odwrotu). */
  grozna?: boolean;
  wykonaj(): void;
}

/**
 * Buduje wiersz czynności.
 *
 * Naciśnięcie nie zwija menu — tak samo jak przy przełączaniu paneli nad
 * kreską. Menu `⋮` zwija naciśnięcie poza nim albo `Escape`
 * (`menu-rozwijane.ts`), a modal pytania o nazwę i modal potwierdzenia
 * usunięcia to natywne `<dialog>` z własną nakładką — stają nad menu, więc
 * nie ma czego chować. Jedno zachowanie w obu sekcjach jednego menu.
 */
export function zbudujWierszCzynnosci(pozycja: PozycjaCzynnosciMenu): HTMLElement {
  const wiersz = document.createElement('button');
  wiersz.type = 'button';
  wiersz.className = 'dn-czynnosci-sesji__pozycja';
  if (pozycja.grozna === true) wiersz.classList.add('dn-czynnosci-sesji__pozycja--grozna');
  wiersz.setAttribute('role', 'menuitem');
  wiersz.dataset.czynnosc = pozycja.klucz;
  wiersz.title = pozycja.przeznaczenie;

  const ikona = document.createElement('span');
  ikona.className = 'dn-czynnosci-sesji__ikona';
  ikona.append(elementIkony(pozycja.ikona, { rozmiar: 16 }));

  const tresc = document.createElement('span');
  tresc.className = 'dn-czynnosci-sesji__tresc';

  const nazwa = document.createElement('span');
  nazwa.className = 'dn-czynnosci-sesji__nazwa';
  nazwa.textContent = pozycja.nazwa;

  const przeznaczenie = document.createElement('span');
  przeznaczenie.className = 'dn-czynnosci-sesji__przeznaczenie';
  przeznaczenie.textContent = pozycja.przeznaczenie;

  tresc.append(nazwa, przeznaczenie);

  wiersz.append(ikona, tresc);

  if (pozycja.skrot !== undefined) {
    const skrot = document.createElement('kbd');
    skrot.className = 'dn-czynnosci-sesji__skrot';
    skrot.textContent = pozycja.skrot;
    wiersz.append(skrot);
  }

  wiersz.addEventListener('click', () => pozycja.wykonaj());

  return wiersz;
}
