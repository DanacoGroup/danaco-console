import { ExtensionKind, type Extension } from '../../../../shared/contract';
import { opisOdmowyBledu } from '../../komponenty/odmowa';
import { utworzRameOkna, type RamaOkna } from '../../komponenty/rama-okna';
import {
  poleTekstowe,
  poleWyboru,
  przycisk,
  utworzWierszOdpowiedzi,
} from '../../modele/kontrolki-formularza';
import { utworzStanOkna, type StanOkna } from './stan-okna';
import { nazwaRodzaju, utworzWierszRozszerzenia } from './wiersz-rozszerzenia';
import type { ZrodloRozszerzen } from './zrodlo-rozszerzen';

/**
 * Katalog rozszerzeń to okno wejściowe do rodziny komend extension — wykaz, instalacja,
 * konfiguracja, przełączenie i odinstalowanie pozycji platformy.
 */
export interface OknoKatalogRozszerzen {
  element: HTMLElement;
  /** Odczyt katalogu z rdzenia; wywoływany przy wczytaniu modułu i po zmianach. */
  wczytaj(): Promise<void>;
}

/**
 * Pozycje selektora rodzaju rozszerzenia wraz z etykietami; pusty kod oznacza brak
 * zawężenia wykazu do jednego rodzaju.
 */
const RODZAJE: readonly { wartosc: string; etykieta: string }[] = [
  { wartosc: '', etykieta: 'wszystkie rodzaje' },
  ...Object.values(ExtensionKind).map((rodzaj) => ({
    wartosc: rodzaj,
    etykieta: nazwaRodzaju(rodzaj),
  })),
];

export function utworzOknoKatalogRozszerzen(zrodlo: ZrodloRozszerzen): OknoKatalogRozszerzen {
  const okno: StanOkna = utworzStanOkna();

  const rama: RamaOkna = utworzRameOkna({
    tytul: 'Katalog rozszerzeń',
    rola: 'zarządca',
    modul: 'agents',
    przedrostek: 'da',
    przeznaczenie:
      'Wykaz rozszerzeń platformy — serwery MCP, wtyczki, integracje API i umiejętności. ' +
      'Instalacja, włączenie i odinstalowanie pozycji.',
  });

  const filtr = poleWyboru(
    {
      etykieta: 'Rodzaj rozszerzenia',
      opis: 'Zawężenie idzie do rdzenia jako pole kind komendy extension.list.',
    },
    RODZAJE.map((r) => ({ wartosc: r.wartosc, etykieta: r.etykieta })),
  );

  const kod = poleTekstowe({
    etykieta: 'Kod pozycji do zainstalowania',
    podpowiedz: 'kod rozszerzenia',
    opis: 'Kod jest stały między wydaniami — to po nim rdzeń rozpoznaje pozycję katalogu.',
  });

  const zrodloPaczki = poleTekstowe({
    etykieta: 'Źródło paczki albo adres serwera',
    podpowiedz: 'pole opcjonalne',
    opis: 'Pominięte, gdy pozycja katalogu niesie własne źródło.',
  });

  const instaluj = przycisk('Zainstaluj pozycję', 'dn-btn dn-btn--sm dn-btn--atrament');
  const odpowiedz = utworzWierszOdpowiedzi();

  const granica = document.createElement('p');
  granica.className = 'dn-pole-opis da-granica';
  granica.textContent =
    'Znaczenie instalacji — pobranie paczki, zarejestrowanie adresu czy zapis punktu ' +
    'dostępu — nie zostało rozstrzygnięte przez Właściciela. Okno przekazuje wywołanie ' +
    'i pokazuje odpowiedź rdzenia bez dopowiadania, co się wydarzyło.';

  const lista = document.createElement('ul');
  lista.className = 'da-rozszerzenia';
  okno.tresc.append(lista);

  rama.narzedzia.append(filtr.element);
  rama.akcje.append(kod.element, zrodloPaczki.element, instaluj);
  rama.cialo.append(okno.element, odpowiedz.element, granica);

  const element = rama.element;

  async function wczytaj(): Promise<void> {
    const rodzaj = filtr.kontrolka.value;
    okno.ladowanie('Odczyt katalogu rozszerzeń…');
    const wynik = await zrodlo.katalog(
      rodzaj === '' ? {} : { rodzaj: rodzaj as ExtensionKind },
    );
    if (!wynik.udany || wynik.wynik === undefined) {
      lista.replaceChildren();
      rama.ustawZnacznik('odmowa', 'blad');
      okno.blad(opisOdmowyBledu('Odczyt katalogu rozszerzeń', wynik.blad));
      return;
    }
    pokaz(wynik.wynik.extensions);
  }

  function pokaz(pozycje: readonly Extension[]): void {
    lista.replaceChildren(...pozycje.map((p) => utworzWierszRozszerzenia(p, obsluz)));
    rama.ustawZnacznik(`${pozycje.length} pozycji`, pozycje.length === 0 ? 'ostrzezenie' : 'neutralna');
    if (pozycje.length === 0) {
      okno.puste(
        filtr.kontrolka.value === ''
          ? 'Rdzeń oddał katalog pusty — platforma nie zna dziś żadnego rozszerzenia.'
          : 'Rdzeń nie zna rozszerzeń tego rodzaju.',
      );
      return;
    }
    okno.gotowe();
  }

  function obsluz(czynnosc: 'instaluj' | 'przelacz' | 'odinstaluj', pozycja: Extension): void {
    if (czynnosc === 'przelacz') {
      void wykonaj(
        `${pozycja.enabled ? 'Wyłączenie' : 'Włączenie'} rozszerzenia`,
        zrodlo.przelacz(pozycja.id, !pozycja.enabled),
      );
      return;
    }
    if (czynnosc === 'odinstaluj') {
      void wykonaj('Odinstalowanie rozszerzenia', zrodlo.odinstaluj(pozycja.id));
      return;
    }
    void wykonaj(
      'Instalacja rozszerzenia',
      zrodlo.zainstaluj(pozycja.code, pozycja.kind),
    );
  }

  /** Odpowiedź rdzenia trafia do wiersza, po czym następuje ponowny odczyt całego wykazu. */
  async function wykonaj(czynnosc: string, bieg: Promise<{ udany: boolean; blad?: unknown }>): Promise<void> {
    odpowiedz.pokaz(`${czynnosc}…`, true);
    const wynik = (await bieg) as { udany: boolean; blad?: { code?: string; message?: string } };
    if (!wynik.udany) {
      odpowiedz.pokaz(opisOdmowyBledu(czynnosc, wynik.blad ?? null), false);
      return;
    }
    odpowiedz.pokaz(`${czynnosc} — rdzeń przyjął wywołanie.`, true);
    await wczytaj();
  }

  async function zainstalujZPola(): Promise<void> {
    const wpisany = kod.kontrolka.value.trim();
    if (wpisany === '') {
      odpowiedz.pokaz('Wskaż kod pozycji — rdzeń odmówi instalacji bez niego.', false);
      return;
    }
    const rodzaj = filtr.kontrolka.value;
    if (rodzaj === '') {
      odpowiedz.pokaz(
        'Wybierz rodzaj rozszerzenia — kontrakt wymaga pola kind przy instalacji.',
        false,
      );
      return;
    }
    await wykonaj(
      'Instalacja rozszerzenia',
      zrodlo.zainstaluj(wpisany, rodzaj as ExtensionKind, zrodloPaczki.kontrolka.value.trim()),
    );
    kod.kontrolka.value = '';
  }

  instaluj.addEventListener('click', () => void zainstalujZPola());
  filtr.kontrolka.addEventListener('change', () => void wczytaj());

  return { element, wczytaj };
}
