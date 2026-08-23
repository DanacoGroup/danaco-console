import { elementIkony, type NazwaIkony } from '../ikony/ikony';
import { ETYKIETA_UTWORZENIA } from './etykiety-pulpitu';
import { RodzajUtworzenia } from './model-danych';
import { utworzSekcje } from './naglowek-sekcji';
import type { ZamiarUtworzenia } from './zdarzenia-pulpitu';

/**
 * Sekcja „Utwórz" pulpitu — rząd przycisków szybkiego tworzenia.
 *
 * Stoi na pulpicie, bo bez niej jedyną drogą do pracy jest kliknięcie sesji już
 * biegnącej w matrycy; ten rząd otwiera wejście tam, gdzie nic jeszcze nie
 * biegnie. Każdy kafel jest czynny, żaden nie czeka na spełnienie warunku,
 * a podpis mówi, co powstanie po naciśnięciu.
 *
 * Etykiety idą krojem bazowym półgrubym, nie szeryfowym: różnica krojów niesie
 * znaczenie — szeryfowy to wejście do środowiska, bazowy to zbudowanie komponentu.
 */
export interface SekcjaUtworz {
  element: HTMLElement;
}

/** Ikona przypisana rodzajowi tworzonego bytu. */
const IKONA_RODZAJU: Readonly<Record<RodzajUtworzenia, NazwaIkony>> = {
  [RodzajUtworzenia.Sesja]: 'plus',
  [RodzajUtworzenia.Projekt]: 'folder',
  [RodzajUtworzenia.Agent]: 'uzytkownik',
  [RodzajUtworzenia.Automatyka]: 'kalendarz',
  [RodzajUtworzenia.Kolejka]: 'menu',
  [RodzajUtworzenia.Zespol]: 'tarcza',
};

/** Kolejność kafli w rzędzie. */
const RZAD: readonly RodzajUtworzenia[] = [
  RodzajUtworzenia.Sesja,
  RodzajUtworzenia.Projekt,
  RodzajUtworzenia.Agent,
  RodzajUtworzenia.Automatyka,
  RodzajUtworzenia.Kolejka,
  RodzajUtworzenia.Zespol,
];

/** Buduje rząd przycisków szybkiego tworzenia. */
export function utworzSekcjeUtworz(nadaj: (zamiar: ZamiarUtworzenia) => void): SekcjaUtworz {
  const { element, cialo } = utworzSekcje(
    'mc-sekcja--utworz',
    {
      tytul: 'Utwórz',
      dopisek: 'Wejście do pracy, w której nic jeszcze nie biegnie.',
      ikona: 'plus',
    },
    'mc-tytul-utworz',
  );

  const rzad = document.createElement('div');
  rzad.className = 'mc-utworz';
  rzad.append(...RZAD.map((rodzaj) => kafel(rodzaj, nadaj)));
  cialo.append(rzad);

  return { element };
}

/** Jeden kafel tworzenia: ikona nad podpisem, cały kafel jest przyciskiem. */
function kafel(
  rodzaj: RodzajUtworzenia,
  nadaj: (zamiar: ZamiarUtworzenia) => void,
): HTMLButtonElement {
  const etykieta = ETYKIETA_UTWORZENIA[rodzaj];

  const element = document.createElement('button');
  element.type = 'button';
  element.className = 'dn-btn dn-btn--zarys mc-utworz__kafel';
  element.dataset.rodzaj = rodzaj;
  element.append(
    elementIkony(IKONA_RODZAJU[rodzaj], { rozmiar: 18 }),
    napis(etykieta),
  );
  element.addEventListener('click', () => nadaj({ rodzaj, etykieta }));
  return element;
}

/** Podpis kafla — krój bazowy półgruby. */
function napis(tekst: string): HTMLElement {
  const element = document.createElement('span');
  element.className = 'mc-utworz__napis';
  element.textContent = tekst;
  return element;
}
