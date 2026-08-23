import { ResearchCredibility, ResearchSourceKind } from '../../../../shared/contract';
import { poleTekstowe } from '../../modele/kontrolki-formularza';
import { zDymkiem } from './dymek-badania';
import { utworzWyborNastawy, wierszNastawy } from './wybor-nastawy';
import type { ZlecenieZrodla } from './zlecenia-badania';

/**
 * Formularz katalogowania źródła — sześć pól żądania `research.source.add`.
 *
 * Jedna odpowiedzialność: zebranie metadanych źródła. Wydzielony z okna, bo
 * „katalogowanie" i „ocena wiarygodności" to dwie z trzech funkcji operatora
 * Sources Manager i mają swoje pola w kontrakcie (`kind`, `origin`, `url`,
 * `libraryFileId`, `credibility`).
 *
 * Każde pole niesie dymek [?], bo każde jest elementem konfiguracji zlecenia.
 *
 * Rodzaj źródła i ocena wiarygodności idą przez `wybor-nastawy` — obsadę
 * bibliotecznego `menu-drzewo`, gdzie uchwyt niesie wartość bieżącą („strona
 * internetowa", „wysoka"), a kliknięcie ją zmienia. Natywny `<select>`
 * pokazywałby ją dopiero po rozwinięciu. Pola tekstowe zostają polami
 * tekstowymi: nie ma w nich czego rozwijać.
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
 * Wartość steru zawężona do wyliczenia kontraktu.
 *
 * Rzutowanie jest tutaj i tylko tutaj: ster oddaje napis, a jego pozycje
 * pochodzą wyłącznie z wyliczeń `shared/contract`, więc innej wartości wydać
 * nie może. Wariant pusty (wykaz bez pozycji) wraca wartością zastępczą, nie
 * napisem pustym.
 */
function wartosc<T extends string>(wybrana: string, zastepcza: T): T {
  return wybrana === '' ? zastepcza : (wybrana as T);
}
