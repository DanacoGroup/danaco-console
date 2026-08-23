import { elementIkony } from '../ikony/ikony';
import './czynnosci.css';

/**
 * `Widok transkryptu ›` — jedyne podmenu sekcji czynności sesji. Podmenu
 * kosztuje kliknięcie i drugi poziom wędrówki ognisk, więc opłaca się tylko
 * tam, gdzie celów jest wiele; pozostałe czynności są płaskie.
 *
 * Tryby przychodzą z zewnątrz: buduje je potok `client/src/rozmowa/`, a ten
 * plik zna je wyłącznie przez port opisany napisami (`kod: string`). Własny
 * spis trybów byłby drugą prawdą o tym samym. Bez portu sekcja czynności jest
 * krótsza o tę pozycję — nie ma wiersza wygaszonego ani strzałki do pustki.
 *
 * Tryb bieżący nosi ptaszek. Wybór jest jednokrotny, więc rola wiersza brzmi
 * `menuitemradio`, a nie `menuitemcheckbox`.
 */

/** Jeden tryb pokazywania zapisu — kształt zgodny z `rozmowa/WIDOKI_ZAPISU`. */
export interface TrybZapisu {
  kod: string;
  nazwa: string;
  przeznaczenie: string;
}

/** Dojście do trybu okna; wypełnia je powłoka, gdy potok rozmowy odda API. */
export interface PortWidokuZapisu {
  tryby: readonly TrybZapisu[];
  /** Kod trybu, w którym okno pokazuje wątek w tej chwili. */
  biezacy(): string;
  /** Przestawia okno na wskazany tryb. */
  ustaw(kod: string): void;
}

export interface PodmenuWidokuZapisu {
  /** Wiersz `Widok transkryptu ›` wraz z listą trybów pod nim. */
  element: HTMLElement;
  /** Przerysowuje ptaszki po zmianie trybu z zewnątrz. */
  odswiez(): void;
  /** Zwija podmenu — wołane przy zwijaniu całego menu `⋮`. */
  zwin(): void;
}

export function utworzPodmenuWidokuZapisu(port: PortWidokuZapisu): PodmenuWidokuZapisu {
  let rozwiniete = false;

  const element = document.createElement('div');
  element.className = 'dn-czynnosci-sesji__podmenu';

  const uchwyt = document.createElement('button');
  uchwyt.type = 'button';
  uchwyt.className = 'dn-czynnosci-sesji__pozycja';
  uchwyt.setAttribute('role', 'menuitem');
  uchwyt.setAttribute('aria-haspopup', 'menu');
  uchwyt.setAttribute('aria-expanded', 'false');
  uchwyt.dataset.czynnosc = 'widok-zapisu';
  uchwyt.title = 'Ile z pracy modelu widać w zapisie tego okna.';

  const ikona = document.createElement('span');
  ikona.className = 'dn-czynnosci-sesji__ikona';
  ikona.append(elementIkony('warstwy', { rozmiar: 16 }));

  const tresc = document.createElement('span');
  tresc.className = 'dn-czynnosci-sesji__tresc';
  const nazwa = document.createElement('span');
  nazwa.className = 'dn-czynnosci-sesji__nazwa';
  nazwa.textContent = 'Widok transkryptu';
  const przeznaczenie = document.createElement('span');
  przeznaczenie.className = 'dn-czynnosci-sesji__przeznaczenie';
  przeznaczenie.textContent = 'Ile z pracy modelu widać w zapisie tego okna.';
  tresc.append(nazwa, przeznaczenie);

  const strzalka = document.createElement('span');
  strzalka.className = 'dn-czynnosci-sesji__strzalka';
  strzalka.append(elementIkony('grot-prawo', { rozmiar: 14 }));

  uchwyt.append(ikona, tresc, strzalka);

  const lista = document.createElement('div');
  lista.className = 'dn-czynnosci-sesji__tryby';
  lista.setAttribute('role', 'menu');
  lista.setAttribute('aria-label', 'Widok transkryptu');
  lista.hidden = true;

  const wiersze = port.tryby.map((tryb) => zbudujTryb(tryb, () => wybierz(tryb.kod)));
  lista.append(...wiersze);
  element.append(uchwyt, lista);

  function wybierz(kod: string): void {
    port.ustaw(kod);
    odswiez();
  }

  /**
   * Zwinięcie chowa listę i każdy jej wiersz. Wędrówka strzałkami
   * w `menu-rozwijane.ts` pomija pozycje z `hidden`, ale pyta o atrybut samej
   * pozycji, nie jej przodka — bez tego ognisko wpadałoby w tryby schowane pod
   * zwiniętą strzałką.
   */
  function ustaw(nowe: boolean): void {
    rozwiniete = nowe;
    lista.hidden = !nowe;
    for (const wiersz of wiersze) wiersz.hidden = !nowe;
    uchwyt.setAttribute('aria-expanded', String(nowe));
  }

  function odswiez(): void {
    const biezacy = port.biezacy();
    for (const wiersz of wiersze) {
      wiersz.setAttribute('aria-checked', String(wiersz.dataset.tryb === biezacy));
    }
  }

  uchwyt.addEventListener('click', () => ustaw(!rozwiniete));
  uchwyt.addEventListener('keydown', (zdarzenie) => {
    if (zdarzenie.key !== 'ArrowRight') return;
    zdarzenie.preventDefault();
    ustaw(true);
    wiersze[0]?.focus();
  });

  // Strzałka w lewo wraca z listy trybów na wiersz nadrzędny — droga powrotna
  // musi istnieć, bo Escape zwija całe menu, a nie samo podmenu.
  lista.addEventListener('keydown', (zdarzenie) => {
    if (zdarzenie.key !== 'ArrowLeft') return;
    zdarzenie.preventDefault();
    zdarzenie.stopPropagation();
    ustaw(false);
    uchwyt.focus();
  });

  ustaw(false);
  odswiez();

  return {
    element,
    odswiez,
    zwin: () => ustaw(false),
  };
}

/** Jeden tryb: ptaszek, nazwa i przeznaczenie. Bez skrótu klawiszowego. */
function zbudujTryb(tryb: TrybZapisu, naWybor: () => void): HTMLElement {
  const wiersz = document.createElement('button');
  wiersz.type = 'button';
  wiersz.className = 'dn-czynnosci-sesji__tryb';
  wiersz.setAttribute('role', 'menuitemradio');
  wiersz.setAttribute('aria-checked', 'false');
  wiersz.dataset.tryb = tryb.kod;
  wiersz.title = tryb.przeznaczenie;

  const znacznik = document.createElement('span');
  znacznik.className = 'dn-czynnosci-sesji__ptaszek';
  znacznik.append(elementIkony('ptaszek', { rozmiar: 14 }));

  const tresc = document.createElement('span');
  tresc.className = 'dn-czynnosci-sesji__tresc';
  const nazwa = document.createElement('span');
  nazwa.className = 'dn-czynnosci-sesji__nazwa';
  nazwa.textContent = tryb.nazwa;
  const przeznaczenie = document.createElement('span');
  przeznaczenie.className = 'dn-czynnosci-sesji__przeznaczenie';
  przeznaczenie.textContent = tryb.przeznaczenie;
  tresc.append(nazwa, przeznaczenie);

  wiersz.append(znacznik, tresc);
  wiersz.addEventListener('click', naWybor);
  return wiersz;
}
