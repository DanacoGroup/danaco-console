import { utworzNaglowekSterowania } from './naglowek-sterowania';
import type { StanSterowania } from './stan-sterowania';
import type { ZmianaOkna } from './zmiana-okna';

const NAZWA = 'Katalogi robocze';

/** Sterowanie listą katalogów roboczych okna, zapisywaną komendą aktualizacji okna z pełną listą po każdej zmianie. */
export function utworzSterowanieKatalogow(
  stan: StanSterowania,
  zmiana: ZmianaOkna,
): HTMLElement {
  const element = document.createElement('div');
  element.className = 'dc-ster-pole dc-ster-pole--lista';

  // Identyfikator niesie identyfikator okna, by dwa komplety nie dzieliły jednego elementu.
  const identyfikator = `dc-ster-katalogi-${stan.idOkna()}`;
  const naglowek = utworzNaglowekSterowania(NAZWA, identyfikator);

  const wykaz = document.createElement('ul');
  wykaz.className = 'dc-ster-katalogi';

  /** Wysyła listę po zmianie; widok wraca ze stanu potwierdzonego. */
  function wyslij(katalogi: string[]): void {
    zmiana.zastosuj(NAZWA, { workingDirs: katalogi });
  }

  function odrysuj(katalogi: string[]): void {
    wykaz.replaceChildren();
    if (katalogi.length === 0) wykaz.append(pustaLista());
    katalogi.forEach((katalog, numer) => {
      wykaz.append(
        pozycja(katalog, () => wyslij(katalogi.filter((_, i) => i !== numer))),
      );
    });
  }

  const dodawanie = utworzDodawanie(identyfikator, (katalog) => {
    wyslij([...stan.migawka().okno.workingDirs, katalog]);
  });

  element.append(naglowek, wykaz, dodawanie);

  stan.naZmiane((migawka) => odrysuj(migawka.okno.workingDirs));
  odrysuj(stan.migawka().okno.workingDirs);

  return element;
}

/** Pozycja listy katalogów: ścieżka katalogu wraz z przyciskiem usuwającym wyłącznie ten jeden wiersz listy. */
function pozycja(katalog: string, przyUsunieciu: () => void): HTMLLIElement {
  const element = document.createElement('li');
  element.className = 'dc-ster-katalogi__pozycja';

  const sciezka = document.createElement('span');
  sciezka.className = 'dc-ster-katalogi__sciezka';
  sciezka.textContent = katalog;
  sciezka.title = katalog;

  const usun = document.createElement('button');
  usun.type = 'button';
  usun.className = 'dc-ster-przycisk dc-ster-przycisk--usun';
  usun.textContent = 'Usuń';
  usun.setAttribute('aria-label', `Usuń katalog ${katalog}`);
  usun.addEventListener('click', przyUsunieciu);

  element.append(sciezka, usun);
  return element;
}

/** Informacja o pustej liście katalogów, pokazywana w wykazie zamiast wiersza pozycji; nie jest komunikatem błędu. */
function pustaLista(): HTMLLIElement {
  const element = document.createElement('li');
  element.className = 'dc-ster-katalogi__pusta';
  element.textContent = 'Brak katalogów — okno pracuje bez wskazanego katalogu';
  return element;
}

/** Wiersz dodawania katalogu: pole ścieżki i przycisk dodania, przycisk pozostaje czynny niezależnie od treści pola. */
function utworzDodawanie(
  identyfikator: string,
  przyDodaniu: (katalog: string) => void,
): HTMLElement {
  const element = document.createElement('div');
  element.className = 'dc-ster-katalogi__dodawanie';

  const pole = document.createElement('input');
  pole.id = identyfikator;
  pole.type = 'text';
  pole.className = 'dc-ster-pole__wpis';
  pole.placeholder = 'ścieżka katalogu roboczego';

  const dodaj = document.createElement('button');
  dodaj.type = 'button';
  dodaj.className = 'dc-ster-przycisk';
  dodaj.textContent = 'Dodaj';

  function zglos(): void {
    const katalog = pole.value.trim();
    if (katalog.length === 0) return;
    pole.value = '';
    przyDodaniu(katalog);
  }

  dodaj.addEventListener('click', zglos);
  pole.addEventListener('keydown', (zdarzenie) => {
    if (zdarzenie.key === 'Enter') zglos();
  });

  element.append(pole, dodaj);
  return element;
}
