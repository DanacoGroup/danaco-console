import { elementIkony, type NazwaIkony } from '../ikony/ikony';
import { ETYKIETA_UTWORZENIA } from './etykiety-pulpitu';
import { RodzajUtworzenia } from './model-danych';
import { utworzSekcje } from './naglowek-sekcji';
import type { ZamiarUtworzenia } from './zdarzenia-pulpitu';

/**
 * Sekcja tworzenia na pulpicie, czyli rząd przycisków szybkiego utworzenia
 * sesji, projektu, agenta, automatyki, kolejki oraz zespołu. Każdy kafel jest
 * czynny bez warunku, a podpis nazywa byt, który powstanie po naciśnięciu.
 */
export interface SekcjaUtworz {
  element: HTMLElement;
}

/**
 * Ikona przypisana każdemu rodzajowi tworzonego bytu. Wykaz obejmuje wszystkie
 * rodzaje utworzenia, więc kafel zawsze ma ikonę, a dodanie nowego rodzaju
 * wymusza wskazanie ikony dla niego.
 */
const IKONA_RODZAJU: Readonly<Record<RodzajUtworzenia, NazwaIkony>> = {
  [RodzajUtworzenia.Sesja]: 'plus',
  [RodzajUtworzenia.Projekt]: 'folder',
  [RodzajUtworzenia.Agent]: 'uzytkownik',
  [RodzajUtworzenia.Automatyka]: 'kalendarz',
  [RodzajUtworzenia.Kolejka]: 'menu',
  [RodzajUtworzenia.Zespol]: 'tarcza',
};

/**
 * Kolejność kafli w rzędzie, od sesji do zespołu. Rząd jest stały i niezależny
 * od stanu pulpitu, dzięki czemu położenie kafla nie zmienia się między
 * kolejnymi odsłonami sekcji.
 */
const RZAD: readonly RodzajUtworzenia[] = [
  RodzajUtworzenia.Sesja,
  RodzajUtworzenia.Projekt,
  RodzajUtworzenia.Agent,
  RodzajUtworzenia.Automatyka,
  RodzajUtworzenia.Kolejka,
  RodzajUtworzenia.Zespol,
];

/**
 * Buduje sekcję szybkiego tworzenia: nagłówek sekcji wraz z dopiskiem oraz
 * rząd kafli w ustalonej kolejności. Naciśnięcie kafla oddaje zamiar
 * utworzenia procedurze podanej przez pulpit.
 */
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

/**
 * Buduje jeden kafel tworzenia jako przycisk z ikoną nad podpisem. Rodzaj bytu
 * trafia do zbioru danych przycisku, a naciśnięcie nadaje zamiar utworzenia
 * wraz z etykietą tego rodzaju.
 */
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

/**
 * Buduje podpis kafla jako element pisany krojem bazowym półgrubym. Podpis
 * niesie etykietę rodzaju, czyli nazwę bytu, który powstanie po naciśnięciu
 * kafla stojącego nad nim.
 */
function napis(tekst: string): HTMLElement {
  const element = document.createElement('span');
  element.className = 'mc-utworz__napis';
  element.textContent = tekst;
  return element;
}
