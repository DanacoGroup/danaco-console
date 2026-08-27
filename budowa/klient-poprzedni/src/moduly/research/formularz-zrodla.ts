import { ResearchCredibility, ResearchSourceKind } from '../../../../shared/contract';
import { poleTekstowe } from '../../modele/kontrolki-formularza';
import { zDymkiem } from './dymek-badania';
import { utworzWyborNastawy, wierszNastawy } from './wybor-nastawy';
import type { ZlecenieZrodla } from './zlecenia-badania';

/**
 * Formularz katalogowania źródła zbiera pola opisujące tytuł, rodzaj, adres, pochodzenie, ocenę wiarygodności oraz dokument repozytorium źródła.
 */
export interface FormularzZrodla {
  element: HTMLElement;
  /** Zlecenie złożone z pól; okno dokłada identyfikator okna badania. */
  zlecenie(idOkna: string): ZlecenieZrodla;
  /** Czy tytuł — jedyne pole obowiązkowe kontraktu — jest wypełniony. */
  czyKompletny(): boolean;
  /** Przenosi ognisko do pierwszego pola. */
  ognisko(): void;
  /** Czyści pola po udanym zapisie. */
  wyczysc(): void;
}

const RODZAJE = [
  { wartosc: ResearchSourceKind.Web, etykieta: 'strona internetowa' },
  { wartosc: ResearchSourceKind.Document, etykieta: 'dokument repozytorium' },
  { wartosc: ResearchSourceKind.Note, etykieta: 'notatka' },
  { wartosc: ResearchSourceKind.Dataset, etykieta: 'zbiór danych' },
];

const OCENY = [
  { wartosc: ResearchCredibility.Unverified, etykieta: 'niezweryfikowane' },
  { wartosc: ResearchCredibility.High, etykieta: 'wysoka' },
  { wartosc: ResearchCredibility.Medium, etykieta: 'średnia' },
  { wartosc: ResearchCredibility.Low, etykieta: 'niska' },
];

export function utworzFormularzZrodla(): FormularzZrodla {
  const tytul = poleTekstowe({ etykieta: 'Tytuł źródła', podpowiedz: 'nazwa materiału' });
  const rodzaj = utworzWyborNastawy('Rodzaj źródła', RODZAJE);
  const adres = poleTekstowe({ etykieta: 'Adres źródła', podpowiedz: 'https://…' });
  const pochodzenie = poleTekstowe({ etykieta: 'Pochodzenie', podpowiedz: 'skąd materiał pochodzi' });
  const wiarygodnosc = utworzWyborNastawy('Ocena wiarygodności', OCENY);
  const plik = poleTekstowe({ etykieta: 'Dokument repozytorium', podpowiedz: 'identyfikator pliku Library' });

  const element = document.createElement('div');
  element.className = 'mr-formularz';
  element.append(
    zDymkiem(tytul.element, 'Pole title — jedyne obowiązkowe obok rodzaju; bez niego rdzeń odmówi zapisu.'),
    zDymkiem(
      wierszNastawy('Rodzaj źródła', rodzaj),
      'Pole kind: web, document, note albo dataset. Rozstrzyga, jak rdzeń kataloguje materiał.',
    ),
    zDymkiem(adres.element, 'Pole url — puste nie idzie do rdzenia.'),
    zDymkiem(pochodzenie.element, 'Pole origin: skąd materiał pochodzi. Metadana wykazu okien obok typu i daty pozyskania.'),
    zDymkiem(
      wierszNastawy('Ocena wiarygodności', wiarygodnosc),
      'Pole credibility: ocena wiarygodności źródła. Domyślnie „niezweryfikowane".',
    ),
    zDymkiem(plik.element, 'Pole libraryFileId — wiąże źródło z dokumentem przyjętym z modułu Library.'),
  );

  return {
    element,

    zlecenie: (idOkna) => ({
      idOkna,
      tytul: tytul.kontrolka.value,
      rodzaj: wartosc(rodzaj.wartosc(), ResearchSourceKind.Web),
      adres: adres.kontrolka.value,
      pochodzenie: pochodzenie.kontrolka.value,
      wiarygodnosc: wartosc(wiarygodnosc.wartosc(), ResearchCredibility.Unverified),
      idPlikuRepozytorium: plik.kontrolka.value,
    }),

    czyKompletny: () => tytul.kontrolka.value.trim() !== '',
    ognisko: () => tytul.kontrolka.focus(),

    wyczysc() {
      for (const pole of [tytul, adres, pochodzenie, plik]) pole.kontrolka.value = '';
    },
  };
}

/**
 * Funkcja zawęża wartość sterowania listą wyboru do wyliczenia kontraktu, zwracając wartość zastępczą, gdy lista nie ma zaznaczonej pozycji.
 */
function wartosc<T extends string>(wybrana: string, zastepcza: T): T {
  return wybrana === '' ? zastepcza : (wybrana as T);
}
