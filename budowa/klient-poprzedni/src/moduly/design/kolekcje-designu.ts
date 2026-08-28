import { Command, type DesignCollection } from '../../../../shared/contract';
import { opisOdmowyBledu } from '../../komponenty/odmowa';
import {
  poleTekstowe,
  poleWyboru,
  przycisk,
  ustawPozycje,
  utworzWierszOdpowiedzi,
} from '../../modele/kontrolki-formularza';
import { nazwaZasobu } from './karta-zasobu';
import type { StanDesignu } from './stan-designu';

/**
 * Kolekcje zasobów w Assets Panel: byt osobny od etykiety, z własną nazwą, opisem i kolejnością,
 * przypisywany dokładką albo odjęciem, nigdy zastąpieniem stanu.
 */
export interface KolekcjeDesignu {
  element: HTMLElement;
  /** Przepisuje podpisy pod zasób wskazany. */
  odswiez(): void;
  /** Zleca odczyt kolekcji okna. */
  wczytaj(): Promise<void>;
}

export function utworzKolekcjeDesignu(stan: StanDesignu): KolekcjeDesignu {
  let zbior: readonly DesignCollection[] = [];

  const nazwa = poleTekstowe({
    etykieta: 'Nazwa nowej kolekcji',
    podpowiedz: 'np. Kampania Q3',
    opis: 'Pole name żądania design.collection.create.',
  });
  const opis = poleTekstowe({
    etykieta: 'Opis kolekcji',
    podpowiedz: 'do czego kolekcja służy',
    opis: 'Pole description; puste zostawia kolekcję bez opisu.',
  });
  const wybor = poleWyboru(
    {
      etykieta: 'Kolekcja',
      opis: 'Pole collectionId żądania design.collection.assign — cel dołożenia i zdjęcia.',
    },
    [],
  );

  const odczytaj = przycisk('Odczytaj kolekcje', 'dn-btn dn-btn--sm dn-btn--zarys');
  const zaloz = przycisk('Załóż kolekcję', 'dn-btn dn-btn--sm dn-btn--atrament');
  const dolóż = przycisk('Dołóż wskazany zasób', 'dn-btn dn-btn--sm dn-btn--zarys');
  const zdejmij = przycisk('Zdejmij wskazany zasób', 'dn-btn dn-btn--sm dn-btn--zarys');
  const odpowiedz = utworzWierszOdpowiedzi();

  const pasek = document.createElement('div');
  pasek.className = 'md-czynnosci__pasek';
  pasek.append(odczytaj, zaloz, dolóż, zdejmij);

  const wykaz = document.createElement('ul');
  wykaz.className = 'md-kolekcje__wykaz';

  const element = document.createElement('div');
  element.className = 'md-kolekcje';
  element.append(nazwa.element, opis.element, wybor.element, pasek, wykaz, odpowiedz.element);

  odczytaj.addEventListener('click', () => void wczytaj());
  zaloz.addEventListener('click', () => void zalozKolekcje());
  dolóż.addEventListener('click', () => void zmienPrzypisanie(false));
  zdejmij.addEventListener('click', () => void zmienPrzypisanie(true));

  function pokazWykaz(): void {
    ustawPozycje(
      wybor.kontrolka,
      zbior.map((kolekcja) => ({
        wartosc: kolekcja.id,
        etykieta: `${kolekcja.name} (${kolekcja.assetCount})`,
      })),
    );
    wykaz.replaceChildren(
      ...zbior.map((kolekcja) => {
        const wiersz = document.createElement('li');
        wiersz.className = 'md-kolekcje__pozycja';
        wiersz.dataset['kolekcja'] = kolekcja.id;
        wiersz.textContent =
          `${kolekcja.name} — zasobów: ${kolekcja.assetCount}` +
          (kolekcja.description === undefined ? '' : ` · ${kolekcja.description}`);
        return wiersz;
      }),
    );
  }

  async function wczytaj(): Promise<void> {
    if (stan.idOkna() === '') {
      odpowiedz.pokaz(
        `Komenda ${Command.DesignCollectionList} wymaga okna modułu. ${stan.opisOkna()}`,
        false,
      );
      return;
    }
    odpowiedz.pokaz('Odczyt kolekcji okna…', true);
    const wynik = await stan.czuwanie.prowadz(
      'odczyt kolekcji',
      // Zawężenia do zasobu tu nie ma: panel pokazuje komplet kolekcji okna, do których można dołożyć.
      stan.zrodlo.kolekcje(stan.idOkna(), ''),
      {
        wToku: (zdanie) => odpowiedz.pokaz(zdanie, true),
        cisza: (zdanie) => odpowiedz.pokaz(zdanie, false),
        powrot: (zdanie) => odpowiedz.pokaz(zdanie, false),
        spozniona: (zdanie) => odpowiedz.pokaz(zdanie, false),
      },
    );
    if (wynik === null) return;
    if (!wynik.udany || wynik.wynik === undefined) {
      odpowiedz.pokaz(opisOdmowyBledu('Odczyt kolekcji', wynik.blad), false);
      return;
    }
    zbior = wynik.wynik.collections;
    pokazWykaz();
    odpowiedz.pokaz(
      zbior.length === 0
        ? 'Rdzeń nie zna ani jednej kolekcji tego okna.'
        : `Kolekcji w oknie: ${wynik.wynik.total}.`,
      true,
    );
  }

  async function zalozKolekcje(): Promise<void> {
    if (stan.idOkna() === '') {
      odpowiedz.pokaz(
        `Komenda ${Command.DesignCollectionCreate} wymaga okna modułu. ${stan.opisOkna()}`,
        false,
      );
      return;
    }
    if (nazwa.kontrolka.value.trim() === '') {
      odpowiedz.pokaz('Kolekcja bez nazwy nie da się odróżnić od żadnej innej.', false);
      return;
    }
    odpowiedz.pokaz('Zakładanie kolekcji…', true);
    const wynik = await stan.czuwanie.prowadz(
      'założenie kolekcji',
      stan.zrodlo.zalozKolekcje(stan.idOkna(), nazwa.kontrolka.value, opis.kontrolka.value),
      {
        wToku: (zdanie) => odpowiedz.pokaz(zdanie, true),
        cisza: (zdanie) => odpowiedz.pokaz(zdanie, false),
        powrot: (zdanie) => odpowiedz.pokaz(zdanie, false),
        spozniona: (zdanie) => odpowiedz.pokaz(zdanie, false),
      },
    );
    if (wynik === null) return;
    if (!wynik.udany || wynik.wynik === undefined) {
      odpowiedz.pokaz(opisOdmowyBledu('Założenie kolekcji', wynik.blad), false);
      return;
    }
    const kolekcja = wynik.wynik.collection;
    zbior = [kolekcja, ...zbior];
    pokazWykaz();
    wybor.kontrolka.value = kolekcja.id;
    nazwa.kontrolka.value = '';
    opis.kontrolka.value = '';
    odpowiedz.pokaz(`Rdzeń założył kolekcję „${kolekcja.name}" (${kolekcja.id}).`, true);
  }

  async function zmienPrzypisanie(zdejmowanie: boolean): Promise<void> {
    const zasob = stan.wybrany();
    if (zasob === null) {
      odpowiedz.pokaz('Wskaż zasób w wykazie — przypisanie dotyczy wskazanych zasobów.', false);
      return;
    }
    if (wybor.kontrolka.value === '') {
      odpowiedz.pokaz('Wskaż kolekcję — najpierw odczytaj kolekcje albo załóż pierwszą.', false);
      return;
    }
    const czynnosc = zdejmowanie ? 'zdjęcie z kolekcji' : 'dołożenie do kolekcji';
    odpowiedz.pokaz(`Zmiana przypisania zasobu „${nazwaZasobu(zasob)}"…`, true);
    const wynik = await stan.czuwanie.prowadz(
      czynnosc,
      stan.zrodlo.przypiszDoKolekcji({
        idKolekcji: wybor.kontrolka.value,
        idZasobow: [zasob.id],
        zdejmij: zdejmowanie,
      }),
      {
        wToku: (zdanie) => odpowiedz.pokaz(zdanie, true),
        cisza: (zdanie) => odpowiedz.pokaz(zdanie, false),
        powrot: (zdanie) => odpowiedz.pokaz(zdanie, false),
        spozniona: (zdanie) => odpowiedz.pokaz(zdanie, false),
      },
    );
    if (wynik === null) return;
    if (!wynik.udany || wynik.wynik === undefined) {
      odpowiedz.pokaz(opisOdmowyBledu(czynnosc, wynik.blad), false);
      return;
    }
    const kolekcja = wynik.wynik.collection;
    zbior = zbior.map((pozycja) => (pozycja.id === kolekcja.id ? kolekcja : pozycja));
    pokazWykaz();
    // Zero zmian jest odpowiedzią udaną: zasób, który w kolekcji już był, nie zmienił niczego.
    odpowiedz.pokaz(
      wynik.wynik.changed === 0
        ? `Nic się nie zmieniło — zasób „${nazwaZasobu(zasob)}" był już w tym stanie. ` +
          `Kolekcja „${kolekcja.name}" ma ${kolekcja.assetCount} zasobów.`
        : `Kolekcja „${kolekcja.name}" ma teraz ${kolekcja.assetCount} zasobów ` +
          `(zmienionych przypisań: ${wynik.wynik.changed}).`,
      true,
    );
  }

  return {
    element,
    wczytaj,

    odswiez() {
      const zasob = stan.wybrany();
      const podpis = zasob === null ? 'wskazany zasób' : `„${nazwaZasobu(zasob)}"`;
      dolóż.textContent = `Dołóż ${podpis}`;
      zdejmij.textContent = `Zdejmij ${podpis}`;
    },
  };
}
