import { AppComponentKind, type AppComponent } from '../../../../shared/contract';
import { poleTekstowe, poleWielowierszowe } from '../../modele/kontrolki-formularza';
import { opiszPole } from './dymek-objasnienia';
import { utworzWyborZMenu, wierszWyboru, type PozycjaWyboruMenu } from './wybor-z-menu';

/**
 * Formularz jednego komponentu architektury — definiowanie komponentów
 * rozwiązania i ustalanie zależności w oknie Architecture Designer.
 *
 * Formularz zbiera pola komponentu i wydaje je w kształcie `AppComponent`
 * z kontraktu. Nie zna kanału i niczego nie wysyła — wysyłką zajmuje się okno.
 *
 * Identyfikator składa się z nazwy, o czym mówi objaśnienie pola. Kontrakt
 * wymaga `id` w każdym komponencie żądania, a rdzeń nadaje własne dopiero
 * w odpowiedzi; identyfikator roboczy jest więc kluczem zależności wewnątrz
 * jednego zapisu, nie obietnicą trwałości.
 */
export interface FormularzKomponentu {
  element: HTMLElement;
  /** Komponent z pól formularza albo `null`, gdy brak nazwy. */
  zbierz(): AppComponent | null;
  /** Odświeża podpowiedzi zależności o komponenty już zestawione. */
  ustawZaleznosci(komponenty: readonly AppComponent[]): void;
  wyczysc(): void;
}

const RODZAJE: readonly PozycjaWyboruMenu[] = [
  { wartosc: AppComponentKind.Frontend, etykieta: 'Warstwa interfejsu' },
  { wartosc: AppComponentKind.Backend, etykieta: 'Warstwa serwerowa' },
  { wartosc: AppComponentKind.Service, etykieta: 'Usługa wydzielona' },
  { wartosc: AppComponentKind.Database, etykieta: 'Baza danych' },
  { wartosc: AppComponentKind.Queue, etykieta: 'Kolejka' },
  { wartosc: AppComponentKind.External, etykieta: 'Składnik zewnętrzny' },
];

export function utworzFormularzKomponentu(): FormularzKomponentu {
  const nazwa = poleTekstowe({
    etykieta: 'Nazwa komponentu',
    podpowiedz: 'np. Panel operatora',
  });
  // Rozwijanie z biblioteki (`komponenty/menu-drzewo.ts` przez obsadę
  // `wybor-z-menu.ts`), nie natywny `<select>`.
  const rodzaj = utworzWyborZMenu('Rodzaj komponentu', RODZAJE);
  const stos = poleTekstowe({
    etykieta: 'Stos technologiczny',
    podpowiedz: 'np. TypeScript · Vite',
  });
  const opis = poleTekstowe({ etykieta: 'Opis komponentu' });
  const zaleznosci = poleTekstowe({
    etykieta: 'Zależy od',
    podpowiedz: 'identyfikatory po przecinku',
  });
  const podpowiedzi = document.createElement('datalist');
  podpowiedzi.id = 'mp-komponenty-podpowiedzi';
  zaleznosci.kontrolka.setAttribute('list', podpowiedzi.id);

  const kontrakt = poleWielowierszowe({ etykieta: 'Kontrakt API komponentu' }, 3);

  const element = document.createElement('div');
  element.className = 'mp-formularz';
  element.append(
    opiszPole(
      nazwa.element,
      'Nazwa buduje też identyfikator roboczy komponentu. Rdzeń nadaje własny ' +
        'identyfikator dopiero w odpowiedzi na apps.architecture.define.',
    ),
    opiszPole(
      wierszWyboru('Rodzaj komponentu', rodzaj),
      'Rodzaj z wyliczenia kontraktu AppComponentKind — sześć wartości.',
    ),
    opiszPole(stos.element, 'Selektor stosu z panelu akcji; pole otwarte, kontrakt nie zamyka listy.'),
    opis.element,
    opiszPole(
      zaleznosci.element,
      'Linie zależności diagramu. Wpisz identyfikatory komponentów już zestawionych; ' +
        'podpowiedź składa się z nich samych.',
    ),
    podpowiedzi,
    opiszPole(kontrakt.element, 'Panel kontraktu API — kontrakt niesie go jako tekst komponentu.'),
  );

  return {
    element,

    zbierz() {
      const miano = nazwa.kontrolka.value.trim();
      if (miano === '') return null;
      const komponent: AppComponent = {
        id: identyfikatorRoboczy(miano),
        name: miano,
        kind: rodzaj.wartosc() as AppComponent['kind'],
      };
      if (stos.kontrolka.value.trim() !== '') komponent.stack = stos.kontrolka.value.trim();
      if (opis.kontrolka.value.trim() !== '') komponent.description = opis.kontrolka.value.trim();
      const wskazane = rozdziel(zaleznosci.kontrolka.value);
      if (wskazane.length > 0) komponent.dependsOn = wskazane;
      if (kontrakt.kontrolka.value.trim() !== '') komponent.apiContract = kontrakt.kontrolka.value;
      return komponent;
    },

    ustawZaleznosci(komponenty) {
      podpowiedzi.replaceChildren(
        ...komponenty.map((komponent) => {
          const pozycja = document.createElement('option');
          pozycja.value = komponent.id;
          pozycja.textContent = komponent.name;
          return pozycja;
        }),
      );
    },

    wyczysc() {
      nazwa.kontrolka.value = '';
      stos.kontrolka.value = '';
      opis.kontrolka.value = '';
      zaleznosci.kontrolka.value = '';
      kontrakt.kontrolka.value = '';
    },
  };
}

/** Identyfikator roboczy komponentu wyprowadzony z nazwy. */
function identyfikatorRoboczy(nazwa: string): string {
  const rdzen = nazwa
    .toLowerCase()
    .replace(/[^a-z0-9]+/g, '-')
    .replace(/^-+|-+$/g, '');
  return rdzen === '' ? 'komponent' : rdzen;
}

/** Rozdziela wpisane identyfikatory zależności, pomijając puste człony. */
function rozdziel(tekst: string): string[] {
  return tekst
    .split(',')
    .map((czlon) => czlon.trim())
    .filter((czlon) => czlon !== '');
}
