import { opisOdmowy } from '../../komponenty/odmowa';
import {
  poleTekstowe,
  przycisk,
  ustawPozycje,
  poleWyboru,
  utworzWierszOdpowiedzi,
} from '../../modele/kontrolki-formularza';
import type { StanBiblioteki } from './stan-biblioteki';
import { przypiszZaznaczenie } from './zapisy-zbiorcze';

/**
 * Formularz zakłada kolekcję komendą `library.collection.create` i przypisuje
 * do niej zasoby zaznaczone w wykazie modułu. Kolekcja jest swobodna, a nie
 * regułowa: kontrakt przyjmuje nazwę oraz opis, bez reguły składającej.
 */
export interface FormularzKolekcji {
  element: HTMLElement;
  /** Przerysowuje wykaz kolekcji w liście wyboru. */
  odswiez(): void;
}

export function utworzFormularzKolekcji(stan: StanBiblioteki): FormularzKolekcji {
  const odpowiedz = utworzWierszOdpowiedzi();

  const nazwa = poleTekstowe({
    etykieta: 'Nazwa nowej kolekcji',
    podpowiedz: 'np. Materiały sprawy',
    opis: 'Kolekcje i etykiety definiuje Operator swobodnie — rdzeń nie narzuca słownika.',
  });
  const opis = poleTekstowe({ etykieta: 'Opis kolekcji', podpowiedz: 'do czego służy' });
  const wybor = poleWyboru({ etykieta: 'Kolekcja docelowa przypisania' }, []);

  const zaloz = przycisk('Nowa kolekcja', 'dn-btn dn-btn--sm dn-btn--atrament');
  zaloz.dataset['czynnosc'] = 'kolekcja-nowa';
  zaloz.addEventListener('click', () => void zalozKolekcje());

  const przypisz = przycisk('Przypisz zaznaczone', 'dn-btn dn-btn--sm dn-btn--zarys');
  przypisz.dataset['czynnosc'] = 'kolekcja-przypisz';
  przypisz.addEventListener('click', () => void przypiszZasoby());

  async function zalozKolekcje(): Promise<void> {
    const wpisane = nazwa.kontrolka.value;
    const miano = wpisane.trim();
    if (miano === '') {
      // Zdanie mówi o przycięciu po stronie okna, a nie o odmowie rdzenia.
      oznaczBlad(
        nazwa.element,
        'Nazwa jest pusta po przycięciu przez okno — nic nie zostało wysłane do rdzenia.',
      );
      return;
    }
    zdejmijBlad(nazwa.element);
    odpowiedz.pokaz(`Zakładanie kolekcji „${miano}"…`, true);
    const wynik = await stan.zrodlo.zalozKolekcje(miano, opis.kontrolka.value.trim());
    if (!wynik.udany || wynik.wynik === undefined) {
      odpowiedz.pokaz(
        opisOdmowy('Utworzenie kolekcji', wynik.blad?.code, wynik.blad?.message),
        false,
      );
      return;
    }
    stan.dopiszKolekcje(wynik.wynik.collectionId);
    // Nazwa idzie z odpowiedzi, a rozbieżność z wpisaną jest wypowiedziana.
    const oddana = wynik.wynik.name;
    const dopisek =
      oddana === wpisane ? '' : ` Rdzeń zapisał nazwę „${oddana}", wpisano „${wpisane}".`;
    odpowiedz.pokaz(
      `Kolekcja „${oddana}" założona (${wynik.wynik.collectionId}).${dopisek}`,
      true,
    );
  }

  async function przypiszZasoby(): Promise<void> {
    const zaznaczone = stan.zaznaczone();
    const pliki = stan.pliki().filter((plik) => zaznaczone.includes(plik.id));
    odpowiedz.pokaz('Przypisywanie zaznaczenia do kolekcji…', true);
    const wynik = await przypiszZaznaczenie(stan, wybor.kontrolka.value, pliki);
    odpowiedz.pokaz(wynik.tresc, wynik.powodzenie);
  }

  const pasek = document.createElement('div');
  pasek.className = 'ml-etykiety__pasek';
  pasek.append(zaloz, przypisz);

  const element = document.createElement('div');
  element.className = 'ml-kolekcje';
  element.append(nazwa.element, opis.element, wybor.element, pasek, odpowiedz.element);

  return {
    element,
    odswiez() {
      ustawPozycje(
        wybor.kontrolka,
        stan.kolekcje().map((kod) => ({ wartosc: kod, etykieta: kod })),
      );
    },
  };
}

/**
 * Ostrzeżenie walidacji przy polu.
 *
 * Pole pozostaje w pełni edytowalne, a formularz wypełniony: błąd walidacji nie
 * blokuje pola, bo Operator ma je poprawić, a nie zaczynać od nowa.
 */
function oznaczBlad(pole: HTMLElement, tresc: string): void {
  pole.dataset['blad'] = 'tak';
  let opis = pole.querySelector<HTMLElement>('.dn-pole-blad');
  if (opis === null) {
    opis = document.createElement('p');
    opis.className = 'dn-pole-blad';
    pole.append(opis);
  }
  opis.textContent = tresc;
}

/**
 * Zdejmuje ostrzeżenie walidacji z pola po zatwierdzeniu, które przeszło:
 * usuwa znacznik błędu oraz opis pod polem. Ostrzeżenie pozostawione przy polu
 * poprawionym mówiłoby o stanie, którego formularz już nie ma.
 */
function zdejmijBlad(pole: HTMLElement): void {
  delete pole.dataset['blad'];
  pole.querySelector('.dn-pole-blad')?.remove();
}
