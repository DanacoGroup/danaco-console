import { poleTekstowe, przycisk, utworzWierszOdpowiedzi } from '../../modele/kontrolki-formularza';
import { AKCJE_ODKRYWANIA } from './akcje-okien';
import {
  przeniesDoZrodel,
  ustawStanOdkrywania,
  wykonajAkcjeOdkrywania,
  wyszukaj,
  type KontekstOdkrywania,
  type TrybOdkrywania,
} from './czynnosci-odkrywania';
import { zDymkiem } from './dymek-badania';
import { KODY_OKIEN } from './kody-okien';
import { utworzRameBadania } from './rama-badania';
import type { StanBadania } from './stan-badania';
import { utworzStanOknaBadania } from './stan-okna-badania';
import { utworzWierszWyniku } from './wiersz-wyniku';
import { utworzWyborNastawy, wierszNastawy } from './wybor-nastawy';
import type { WynikOdkrycia } from './wynik-odkrycia';

/**
 * Discovery Panel wyszukuje i odkrywa źródła badania, kierując zapytanie do komendy wyszukiwania semantycznego, pełnotekstowego, webowego albo naukowego według wybranego trybu.
 */
export interface OknoDiscoveryPanel {
  element: HTMLElement;
  odswiez(): void;
}

/** Stała wylicza cztery tryby zapytania Discovery Panel wraz z komendą kontraktu, którą dziś obsługuje każdy z nich. */
const TRYBY: readonly { wartosc: TrybOdkrywania; etykieta: string; opis: string }[] = [
  {
    wartosc: 'znaczenie',
    etykieta: 'własne semantyczne',
    opis: 'Szuka po znaczeniu w bibliotece wiedzy, historii rozmów i przestrzeni roboczej okna badania. Idzie komendą knowledge.search.',
  },
  {
    wartosc: 'tresc',
    etykieta: 'pełnotekstowe',
    opis: 'Szuka po treści w zasobach repozytorium Library. Idzie komendą library.file.search.',
  },
  {
    wartosc: 'web',
    etykieta: 'webowe',
    opis: 'Wyszukiwanie w sieci. Komenda jest już w kontrakcie — o trybie rozstrzyga jej pole mode — ale rdzeń nie ma dla niej uchwytu, więc zapytanie wraca odmową.',
  },
  {
    wartosc: 'naukowy',
    etykieta: 'naukowe',
    opis: 'Wyszukiwanie w bazach publikacji. Ta sama komenda co tryb webowy, inna wartość pola mode; rdzeń nie ma dla niej jeszcze uchwytu.',
  },
];

export function utworzOknoDiscoveryPanel(
  stan: StanBadania,
  przejdz: (kodOkna: string) => void,
): OknoDiscoveryPanel {
  const tryb = utworzWyborNastawy('Tryb wyszukiwania', TRYBY);
  const polePytania = poleTekstowe({
    etykieta: 'Zapytanie',
    podpowiedz: 'czego szukamy w materiale',
  });

  let wyniki: readonly WynikOdkrycia[] = [];

  const wykaz = document.createElement('ul');
  wykaz.className = 'mr-wykaz';

  const kontekst: KontekstOdkrywania = {
    stan,
    okno: utworzStanOknaBadania(),
    odpowiedz: utworzWierszOdpowiedzi(),
    zapytanie: polePytania.kontrolka,
    tryb: () => wartoscTrybu(tryb.wartosc()),
    wyniki: () => wyniki,
    ustawWyniki: (pozycje) => {
      wyniki = pozycje;
      odswiez();
    },
    przejdz,
  };

  const szukaj = przycisk('Szukaj', 'dn-btn dn-btn--sm dn-btn--atrament');
  szukaj.addEventListener('click', () => void wyszukaj(kontekst));
  // Enter w polu zapytania uruchamia wyszukiwanie tak samo jak przycisk Szukaj.
  polePytania.kontrolka.addEventListener('keydown', (zdarzenie) => {
    if (zdarzenie.key === 'Enter') void wyszukaj(kontekst);
  });

  kontekst.okno.tresc.append(
    zDymkiem(
      wierszNastawy('Tryb wyszukiwania', tryb),
      'Tryb rozstrzyga, którą komendą pójdzie zapytanie. Dwa tryby mają komendę kontraktu, dwa nie mają i kończą się odmową rdzenia.',
    ),
    zDymkiem(
      polePytania.element,
      'Pole query komendy knowledge.search albo library.file.search — zależnie od trybu. Enter uruchamia wyszukiwanie.',
    ),
    szukaj,
    kontekst.odpowiedz.element,
    wykaz,
  );

  const rama = utworzRameBadania(
    KODY_OKIEN.odkrywanie,
    'Discovery Panel',
    'pomocnicze',
    AKCJE_ODKRYWANIA,
    (akcja) => void wykonajAkcjeOdkrywania(kontekst, akcja),
  );
  rama.cialo.append(kontekst.okno.element);

  function odswiez(): void {
    wykaz.replaceChildren(
      ...wyniki.map((pozycja) =>
        utworzWierszWyniku(pozycja, (wybrana) => void przeniesDoZrodel(kontekst, wybrana)),
      ),
    );
    ustawStanOdkrywania(kontekst, wyniki.length);
  }

  odswiez();
  return { element: rama.element, odswiez };
}

/**
 * Wartość steru zawężona do trybu panelu.
 *
 * Ster oddaje napis, a jego pozycje pochodzą wyłącznie z wykazu wyżej, więc
 * innej wartości wydać nie może. Wariant pusty wraca trybem początkowym, nie
 * napisem pustym — ten sam zabieg co w `formularz-zrodla.ts`.
 */
function wartoscTrybu(wybrana: string): TrybOdkrywania {
  const pozycja = TRYBY.find((tryb) => tryb.wartosc === wybrana);
  return pozycja?.wartosc ?? 'znaczenie';
}
