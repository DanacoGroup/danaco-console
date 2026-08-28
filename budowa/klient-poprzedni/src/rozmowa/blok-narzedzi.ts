import { blokKodu, utworzBlokZwijany, type BlokZwijany } from './blok-zwijany';
import { NAPISY } from './etykiety-rozmowy';
import type { WywolanieNarzedzia } from './wpis-rozmowy';

/** Blok wywołań narzędzi jednego wpisu rozmowy, pokazywany zawsze wraz z pełnymi wynikami tych wywołań. */
export interface BlokNarzedzi {
  /** Element montowany we wpisie. */
  element: HTMLElement;
  /** Wypełnia blok wywołaniami albo go ukrywa, gdy ich nie było. */
  aktualizuj(narzedzia: WywolanieNarzedzia[]): void;
  /** Rozwija blok — tryb `pelny` kładzie argumenty i wyniki na wierzchu. */
  ustawRozwiniecie(otwarty: boolean): void;
}

/** Wywołania narzędzi tury — wynik narzędzia stoi zawsze w tej samej pozycji, w której padło wywołanie. */
export function utworzBlokNarzedzi(): BlokNarzedzi {
  const blok: BlokZwijany = utworzBlokZwijany({ tytul: NAPISY.narzedzia, ikona: 'kod' });

  function aktualizuj(narzedzia: WywolanieNarzedzia[]): void {
    if (narzedzia.length === 0) {
      blok.pokaz(false);
      return;
    }
    blok.pokaz(true);
    blok.ustawPodtytul(podtytul(narzedzia));
    blok.tresc.replaceChildren(...narzedzia.map(pozycja));
  }

  return { element: blok.element, aktualizuj, ustawRozwiniecie: blok.ustawRozwiniecie };
}

/** Podtytuł bloku: liczba wszystkich wywołań w tej turze oraz liczba wyników zakończonych pełnym błędem. */
function podtytul(narzedzia: WywolanieNarzedzia[]): string {
  const bledne = narzedzia.filter((pozycja) => pozycja.bledne).length;
  const opis = `wywołań: ${narzedzia.length}`;
  return bledne > 0 ? `${opis} · błędnych: ${bledne}` : opis;
}

/** Jedno wywołanie narzędzia wraz z jego wynikiem, pokazywane jako osobna pozycja tego samego bloku wywołań. */
function pozycja(narzedzie: WywolanieNarzedzia): HTMLElement {
  const element = document.createElement('article');
  element.className = 'dc-narzedzie';

  const nazwa = document.createElement('span');
  nazwa.className = 'dn-plakietka' + (narzedzie.bledne ? ' dn-plakietka--blad' : '');
  nazwa.textContent = narzedzie.bledne ? `${narzedzie.nazwa} — błąd` : narzedzie.nazwa;
  element.append(nazwa);

  if (narzedzie.wejscie.length > 0) element.append(blokKodu(narzedzie.wejscie));
  if (narzedzie.wynik.length > 0) element.append(blokKodu(narzedzie.wynik));

  return element;
}
