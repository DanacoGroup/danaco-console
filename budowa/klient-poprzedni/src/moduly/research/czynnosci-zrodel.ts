import { opisOdmowyBledu } from '../../komponenty/odmowa';
import type { WierszOdpowiedzi } from '../../modele/kontrolki-formularza';
import type { AkcjaBadania } from './akcje-okien';
import type { FormularzZrodla } from './formularz-zrodla';
import { KODY_OKIEN } from './kody-okien';
import { pokazPustke, PUSTKA_ZRODEL } from './pustka-okien';
import type { StanBadania } from './stan-badania';
import { czyKomendaBadania, wykonajKomendeBadania } from './wywolania-komend';
import type { StanOknaBadania } from './stan-okna-badania';

/**
 * Plik odpowiada na pytanie, co dzieje się po naciśnięciu kontrolki w oknie
 * Sources Manager, podczas gdy samo składanie widoku należy do pliku okna.
 */
export interface KontekstZrodel {
  stan: StanBadania;
  okno: StanOknaBadania;
  odpowiedz: WierszOdpowiedzi;
  formularz: FormularzZrodla;
  przejdz(kodOkna: string): void;
}

/**
 * Rozdziela akcję panelu źródeł na drogę własną okna, obsługiwaną komendą
 * dedykowaną, oraz drogę generyczną przekazywaną rdzeniowi jako window.action.
 */
export async function wykonajAkcjeZrodel(
  kontekst: KontekstZrodel,
  akcja: AkcjaBadania,
): Promise<void> {
  if (akcja.kod === 'research.source.new') {
    kontekst.formularz.ognisko();
    kontekst.odpowiedz.pokaz(
      'Formularz katalogowania jest wyżej — wypełnij tytuł i naciśnij „Skataloguj źródło".',
      true,
    );
    return;
  }
  if (akcja.kod === 'research.source.link') {
    kontekst.przejdz(KODY_OKIEN.ustalenia);
    kontekst.odpowiedz.pokaz(
      'Zaznaczone źródła są do wzięcia w Findings Panel — powiązanie jedzie w polu sourceIds.',
      true,
    );
    return;
  }
  if (akcja.kod === 'research.source.read') {
    otworzLekture(kontekst);
    return;
  }
  if (akcja.droga === 'komenda' && czyKomendaBadania(akcja.kod)) {
    const wynik = await wykonajKomendeBadania(
      { stan: kontekst.stan, tekst: tekstDlaKomendy(kontekst) },
      akcja,
    );
    kontekst.odpowiedz.pokaz(wynik.opis, wynik.udany);
    return;
  }
  await przezPanelAkcji(kontekst, akcja);
}

/**
 * Otwiera lekturę pierwszego z zaznaczonych źródeł, ponieważ pasek akcji działa
 * na zaznaczeniu i nie zna wiersza, nad którym stoi kursor.
 */
function otworzLekture(kontekst: KontekstZrodel): void {
  const { stan, odpowiedz } = kontekst;
  const wybrane = stan.wybraneZrodla.wybrane();
  const pierwsze = wybrane[0] ?? '';
  if (pierwsze === '') {
    odpowiedz.pokaz(
      'Zaznacz źródło albo naciśnij „Czytaj" przy pozycji wykazu — czytnik potrzebuje wskazania materiału.',
      false,
    );
    return;
  }
  stan.lektura.wskaz(pierwsze);
  kontekst.przejdz(KODY_OKIEN.lektura);
  const tytul = stan.zrodla().find((zrodlo) => zrodlo.id === pierwsze)?.title ?? pierwsze;
  odpowiedz.pokaz(
    wybrane.length === 1
      ? `Reading View wczytuje źródło „${tytul}".`
      : `Zaznaczonych źródeł jest ${String(wybrane.length)}; czytać można jedno naraz — Reading View wczytuje pierwsze, „${tytul}".`,
    true,
  );
}

/**
 * Przekazuje kod akcji oraz zaznaczone źródła do rdzenia komendą window.action,
 * bez własnej obsługi po stronie okna Sources Manager.
 */
async function przezPanelAkcji(kontekst: KontekstZrodel, akcja: AkcjaBadania): Promise<void> {
  const { stan, odpowiedz } = kontekst;
  if (stan.idOkna() === '') {
    odpowiedz.pokaz(
      `Akcja „${akcja.nazwa}" wymaga okna badania, którego rdzeń jeszcze nie wskazał (window.list nie oddał okna modułu Research).`,
      false,
    );
    return;
  }
  const wybrane = stan.wybraneZrodla.wybrane();
  odpowiedz.pokaz(`Wysłano window.action ${akcja.kod} dla ${wybrane.length} zaznaczonych źródeł…`, true);
  const wynik = await stan.okna.akcja(stan.idOkna(), akcja.kod, { sourceIds: wybrane });
  if (!wynik.udany) {
    odpowiedz.pokaz(opisOdmowyBledu(`Akcja „${akcja.nazwa}"`, wynik.blad, wynik.nieznanyTyp), false);
    return;
  }
  // Zdanie mówi o wyniku oddanym przez rdzeń, nie o wykonaniu samej akcji.
  odpowiedz.pokaz(`Rdzeń oddał wynik akcji „${akcja.nazwa}".`, true);
}

/**
 * Kataloguje źródło komendą research.source.add wraz z metadanymi pobranymi
 * z formularza, po sprawdzeniu, że pole tytułu zostało wypełnione.
 */
export async function skatalogujZrodlo(kontekst: KontekstZrodel): Promise<void> {
  const { stan, okno, odpowiedz, formularz } = kontekst;
  if (!formularz.czyKompletny()) {
    formularz.ognisko();
    odpowiedz.pokaz('Wpisz tytuł źródła — pole title jest w kontrakcie obowiązkowe.', false);
    return;
  }
  if (stan.idOkna() === '') {
    odpowiedz.pokaz(
      'Rdzeń nie wskazał okna badania; komenda research.source.add wymaga pola windowId.',
      false,
    );
    return;
  }
  okno.ladowanie('Katalogowanie źródła…');
  const wynik = await stan.zrodlo.dodajZrodlo(formularz.zlecenie(stan.idOkna()));
  if (!wynik.udany || wynik.wynik === undefined) {
    const opis = opisOdmowyBledu('Skatalogowanie źródła', wynik.blad, wynik.nieznanyTyp);
    okno.blad(opis);
    odpowiedz.pokaz(opis, false);
    return;
  }
  stan.wchlonZrodlo(wynik.wynik.source);
  formularz.wyczysc();
  odpowiedz.pokaz(`Rdzeń skatalogował źródło „${wynik.wynik.source.title}".`, true);
}

/**
 * Ustawia jeden z trzech stanów wykazu źródeł: odczyt trwa, wykaz niesie treść,
 * albo nie ma czego pokazać; pustkę odróżnia od braku wskazanego okna badania.
 */
export function ustawStanZrodel(kontekst: KontekstZrodel, liczba: number): void {
  const { stan, okno } = kontekst;
  if (stan.faza() === 'odczyt') {
    okno.ladowanie('Odczyt okna badania z rdzenia…');
    return;
  }
  if (liczba > 0) {
    okno.gotowe();
    return;
  }
  if (stan.faza() === 'blad') {
    okno.blad(stan.powod());
    return;
  }
  pokazPustke(okno, stan, PUSTKA_ZRODEL);
}

/**
 * Podaje tekst swobodny okna przekazywany komendom, które nie mają własnego
 * formularza; brak treści nazywa plik wywolania-komend.ts.
 */
function tekstDlaKomendy(kontekst: KontekstZrodel): string {
  // Pole tytułu niesie adres, ścieżkę bibliografii albo etykiety po przecinku.
  return kontekst.formularz.zlecenie('').tytul.trim();
}
