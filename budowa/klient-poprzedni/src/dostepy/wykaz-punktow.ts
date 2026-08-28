import type { AccessPoint, AccessPointKind } from '../../../shared/contract';
import { utworzDodanieKatalogu } from './dodanie-katalogu';
import { utworzKartePunktu, type KartaPunktu } from './karta-punktu';
import { nazwaGrupy, RODZAJE } from './nazwy-dostepow';
import type { StanDostepow } from './stan-dostepow';

/**
 * Wykaz punktów dostępu rozdzielony rodzajem: maszyny osobno, katalogi lokalne
 * osobno. Karty przeżywają przeliczenie, bo zmiana stanu nanosi wartości na
 * karty już zbudowane, a przebudowa idzie po zmianie składu wykazu.
 */
export interface WykazPunktow {
  /** Element osadzany w sekcji dostępów. */
  element: HTMLElement;
  /** Nanosi stan na karty; przebudowuje wykaz, gdy zmienił się jego skład. */
  odswiez(): void;
  /** Odłącza nasłuch powłoki używany przez dodawanie katalogu. */
  rozlacz(): void;
}

export function utworzWykazPunktow(stan: StanDostepow): WykazPunktow {
  const dodanie = utworzDodanieKatalogu(stan);

  const karty = new Map<string, KartaPunktu>();
  let skladWykazu = '';

  const cialo = document.createElement('div');
  cialo.className = 'dd-wykaz__cialo';

  const pusty = stanPusty();

  const element = document.createElement('section');
  element.className = 'dd-wykaz';
  element.append(naglowek(), pusty, cialo, dodanie.element);

  /** Przebudowa wykazu: grupy rodzajów, w każdej karty punktów tego rodzaju. */
  function przebuduj(punkty: readonly AccessPoint[]): void {
    karty.clear();
    cialo.replaceChildren();

    for (const rodzaj of rodzajeWykazu(punkty)) {
      const nalezace = punkty.filter((punkt) => punkt.kind === rodzaj);
      if (nalezace.length === 0) continue;

      const grupa = document.createElement('div');
      grupa.className = 'dd-grupa';
      grupa.append(naglowekGrupy(rodzaj, nalezace.length));

      for (const punkt of nalezace) {
        const karta = utworzKartePunktu(punkt, stan);
        karty.set(punkt.id, karta);
        grupa.append(karta.element);
      }

      cialo.append(grupa);
    }
  }

  return {
    element,

    odswiez() {
      const punkty = stan.punkty();
      // Sklad obejmuje korzenie; stan punktu i czas sprawdzenia do niego nie wchodza.
      const sklad = punkty
        .map((punkt) => `${punkt.kind}:${punkt.id}:${punkt.roots.join(',')}`)
        .join('|');
      if (sklad !== skladWykazu) {
        skladWykazu = sklad;
        przebuduj(punkty);
      }

      // Stan pusty nalezy sie wylacznie odczytowi zakonczonemu powodzeniem.
      pusty.hidden = punkty.length > 0 || stan.faza() !== 'gotowe';

      const nadane = new Set(stan.nadania().map((nadanie) => nadanie.accessPointId));
      for (const punkt of punkty) {
        karty.get(punkt.id)?.odswiez(punkt, nadane.has(punkt.id));
      }
    },

    rozlacz: dodanie.rozlacz,
  };
}

/**
 * Rodzaje w kolejności stałej, a po nich rodzaje spoza wyliczenia, gdyby
 * rdzeń przysłał wartość nowszą od klienta. Nierozpoznany rodzaj nie znika
 * z wykazu — dostaje własną grupę pod własną nazwą.
 */
function rodzajeWykazu(punkty: readonly AccessPoint[]): AccessPointKind[] {
  const nieznane = punkty
    .map((punkt) => punkt.kind)
    .filter((rodzaj) => !RODZAJE.includes(rodzaj));
  return [...RODZAJE, ...new Set(nieznane)];
}

/**
 * Stan pusty wykazu (`.dn-pusty-stan` z biblioteki komponentów).
 *
 * Mówi, że stan jest oczekiwany i naturalny — świeża instalacja nie ma ani
 * jednego punktu — i wskazuje pierwszą czynność, zamiast zostawiać pusty
 * prostokąt.
 */
function stanPusty(): HTMLElement {
  const element = document.createElement('div');
  element.className = 'dn-pusty-stan dd-wykaz__pusty';
  element.hidden = true;

  const tytul = document.createElement('p');
  tytul.className = 'dn-pusty-stan-tytul';
  tytul.textContent = 'Brak punktów dostępu';

  const opis = document.createElement('p');
  opis.className = 'dn-pusty-stan-opis';
  opis.textContent =
    'Rdzeń odpowiedział, ale nie zna ani jednego punktu. Dodaj katalog lokalny poniżej.';

  element.append(tytul, opis);
  return element;
}

function naglowek(): HTMLElement {
  const element = document.createElement('header');
  element.className = 'dd-wykaz__naglowek';

  const tytul = document.createElement('h3');
  tytul.className = 'dd-wykaz__tytul';
  tytul.textContent = 'Punkty dostępu';

  const opis = document.createElement('p');
  opis.className = 'dd-wykaz__opis';
  opis.textContent =
    'Maszyny i katalogi, do których model może sięgnąć. Punkt jest bytem platformy; dostęp do niego nadaje się osobno, oknu rozmowy.';

  element.append(tytul, opis);
  return element;
}

function naglowekGrupy(rodzaj: AccessPointKind, liczba: number): HTMLElement {
  const element = document.createElement('h4');
  element.className = 'dd-grupa__tytul';
  element.textContent = `${nazwaGrupy(rodzaj)} (${liczba})`;
  return element;
}
