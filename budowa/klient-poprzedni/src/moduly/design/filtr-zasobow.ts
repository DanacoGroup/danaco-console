import { DesignAssetKind } from '../../../../shared/contract';
import { poleTekstowe, poleWyboru, przycisk } from '../../modele/kontrolki-formularza';
import { dopnijDymek } from './dymek-objasnienia';
import { przytnijKod, rozbijEtykiety } from './przyciecie-pol';
import type { ZapytanieZasobow } from './zrodlo-designu';

/**
 * Filtr wykazu zasobów — warunki, którymi Assets Panel zawęża odczyt.
 *
 * Jedna odpowiedzialność: zebranie warunków i oddanie ich oknu.
 *
 * Cztery pola idą do rdzenia, piąte zostaje tutaj. Rodzaj, etykiety, ulubione
 * i granica są polami żądania `design.asset.list` — zawężają odczyt. Fraza
 * wyszukiwania nie jest polem tego żądania, więc zawęża wyłącznie to, co już
 * wróciło; pole mówi o tym wprost, zamiast pozorować wyszukiwanie po stronie
 * rdzenia.
 */
export interface FiltrZasobow {
  element: HTMLElement;
  /** Warunki odczytu złożone z pól idących do rdzenia. */
  warunki(): Partial<ZapytanieZasobow>;
  /** Fraza zawężająca miejscowo; pusta znaczy „bez zawężenia". */
  fraza(): string;
  /** Widok wykazu: siatka albo lista. */
  widok(): 'siatka' | 'lista';
}

/** Wywołania zwrotne filtra — okno decyduje, co zrobić ze zmianą. */
export interface ObslugaFiltra {
  naOdczyt(): void;
  naZawezenie(): void;
  naWidok(): void;
}

export function utworzFiltrZasobow(obsluga: ObslugaFiltra): FiltrZasobow {
  const rodzaj = poleWyboru({ etykieta: 'Rodzaj zasobu' }, [
    { wartosc: '', etykieta: 'wszystkie rodzaje' },
    { wartosc: DesignAssetKind.Image, etykieta: 'grafika rastrowa' },
    { wartosc: DesignAssetKind.Vector, etykieta: 'grafika wektorowa' },
    { wartosc: DesignAssetKind.Composition, etykieta: 'kompozycja tablicy' },
  ]);
  dopnijDymek(rodzaj.element, 'Pole kind żądania design.asset.list — zawęża odczyt po stronie rdzenia.');

  const etykiety = poleTekstowe({
    etykieta: 'Etykiety i kolekcje',
    podpowiedz: 'etykieta, druga etykieta',
    opis:
      'Rozdzielone przecinkiem. Zawężają odczyt; nadaje się je zasobowi wskazanemu ' +
      'kontrolką pod wykazem. Kolekcje zasobów nadal nie mają komendy.',
  });
  dopnijDymek(
    etykiety.element,
    'Pole tags żądania design.asset.list. Zawężenie po WSZYSTKICH wskazanych etykietach naraz. ' +
      'Rdzeń porównuje etykietę co do znaku i sam jej nie przycina — pole przycina brzegi ' +
      'tak samo jak kontrolka nadająca etykiety, żeby jedno trafiało w drugie.',
  );

  const granica = poleTekstowe({
    etykieta: 'Górna granica liczby zasobów',
    podpowiedz: '60',
  });
  granica.kontrolka.value = '60';
  granica.kontrolka.inputMode = 'numeric';
  dopnijDymek(granica.element, 'Pole limit żądania. Zero i pusta wartość znaczą „bez granicy".');

  const fraza = poleTekstowe({
    etykieta: 'Szukaj w wykazie wczytanym',
    podpowiedz: 'fragment nazwy albo etykiety',
    opis: 'Zawężenie miejscowe — żądanie design.asset.list nie ma pola frazy.',
  });

  const ulubione = document.createElement('input');
  ulubione.type = 'checkbox';
  ulubione.className = 'dn-check';
  ulubione.id = 'md-filtr-ulubione';

  const etykietaUlubionych = document.createElement('label');
  etykietaUlubionych.className = 'dn-pole-etykieta';
  etykietaUlubionych.htmlFor = ulubione.id;
  etykietaUlubionych.textContent = 'Tylko ulubione';

  const rzadUlubionych = document.createElement('div');
  rzadUlubionych.className = 'md-filtr__przelacz';
  rzadUlubionych.append(ulubione, etykietaUlubionych);

  const widok = poleWyboru({ etykieta: 'Widok wykazu' }, [
    { wartosc: 'siatka', etykieta: 'siatka' },
    { wartosc: 'lista', etykieta: 'lista' },
  ]);

  const odczytaj = przycisk('Odczytaj zasoby', 'dn-btn dn-btn--sm dn-btn--atrament');

  const element = document.createElement('section');
  element.className = 'md-filtr';
  element.setAttribute('aria-label', 'Filtr wykazu zasobów');
  element.append(
    rodzaj.element,
    etykiety.element,
    granica.element,
    fraza.element,
    rzadUlubionych,
    widok.element,
    odczytaj,
  );

  odczytaj.addEventListener('click', () => obsluga.naOdczyt());
  for (const kontrolka of [rodzaj.kontrolka, ulubione]) {
    kontrolka.addEventListener('change', () => obsluga.naOdczyt());
  }
  fraza.kontrolka.addEventListener('input', () => obsluga.naZawezenie());
  widok.kontrolka.addEventListener('change', () => obsluga.naWidok());

  return {
    element,

    warunki() {
      const liczba = Number.parseInt(granica.kontrolka.value, 10);
      return {
        rodzaj: rodzaj.kontrolka.value as DesignAssetKind | '',
        // Ten sam rozbiór co przy nadawaniu etykiet — inaczej filtr nie trafia
        // we własny zapis.
        etykiety: rozbijEtykiety(etykiety.kontrolka.value),
        tylkoUlubione: ulubione.checked,
        granica: Number.isFinite(liczba) && liczba > 0 ? liczba : 0,
      };
    },

    fraza: () => przytnijKod(fraza.kontrolka.value).toLowerCase(),
    widok: () => (widok.kontrolka.value === 'lista' ? 'lista' : 'siatka'),
  };
}
