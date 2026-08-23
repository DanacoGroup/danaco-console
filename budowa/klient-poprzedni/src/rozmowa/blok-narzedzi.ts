import { blokKodu, utworzBlokZwijany, type BlokZwijany } from './blok-zwijany';
import { NAPISY } from './etykiety-rozmowy';
import type { WywolanieNarzedzia } from './wpis-rozmowy';

/** Blok wywołań narzędzi jednego wpisu. */
export interface BlokNarzedzi {
  /** Element montowany we wpisie. */
  element: HTMLElement;
  /** Wypełnia blok wywołaniami albo go ukrywa, gdy ich nie było. */
  aktualizuj(narzedzia: WywolanieNarzedzia[]): void;
  /** Rozwija blok — tryb `pelny` kładzie argumenty i wyniki na wierzchu. */
  ustawRozwiniecie(otwarty: boolean): void;
}

/**
 * Wywołania narzędzi tury — rodzaje `tool_use` i `tool_result` strumienia.
 *
 * Wynik narzędzia jest w kontrakcie osobną rolą wiadomości, ale w historii
 * należy do tury, w której padł. Dlatego wywołanie i jego wynik stoją w jednej
 * pozycji, bez przeskakiwania między wpisami.
 */
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

/** Podtytuł: liczba wywołań i liczba wyników błędnych. */
function podtytul(narzedzia: WywolanieNarzedzia[]): string {
  const bledne = narzedzia.filter((pozycja) => pozycja.bledne).length;
  const opis = `wywołań: ${narzedzia.length}`;
  return bledne > 0 ? `${opis} · błędnych: ${bledne}` : opis;
}

/** Jedno wywołanie wraz z wynikiem. */
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
