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
 * Czynności Operatora w Sources Manager.
 *
 * Jedna odpowiedzialność: co się dzieje po naciśnięciu. Wydzielone z pliku
 * okna, bo okno składa widok, a to jest jego zachowanie — dwie rzeczy, dwa
 * pliki, oba w rozmiarze do przeczytania naraz.
 *
 * Każde naciśnięcie daje odpowiedź. Akcja bez własnej komendy idzie
 * generycznym `window.action`; brak okna badania i brak treści pola wracają
 * zdaniem, nie ciszą i nie wygaszeniem kontrolki.
 */
export interface KontekstZrodel {
  stan: StanBadania;
  okno: StanOknaBadania;
  odpowiedz: WierszOdpowiedzi;
  formularz: FormularzZrodla;
  przejdz(kodOkna: string): void;
}

/** Rozdziela akcję panelu na drogę własną okna i drogę generyczną. */
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
 * Otwarcie lektury z paska akcji — pierwsze z zaznaczonych źródeł.
 *
 * Przycisk „Czytaj" przy pozycji wskazuje źródło wprost; akcja paska działa na
 * zaznaczeniu, bo pasek nie zna wiersza, nad którym stoi kursor. Czytać można
 * jedno źródło naraz, więc przy wielu zaznaczonych okno mówi, które wzięło,
 * zamiast wybierać po cichu.
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

/** Droga generyczna: `window.action` z kodem akcji i zaznaczeniem w parametrach. */
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
  // Zdanie mówi o wyniku oddanym przez rdzeń, nie o wykonaniu akcji: rdzeń
  // kwituje sukcesem samo przyjęcie zgłoszenia. Odmowę braku wykonawcy
  // pokazuje gałąź wyżej, słowami rdzenia.
  odpowiedz.pokaz(`Rdzeń oddał wynik akcji „${akcja.nazwa}".`, true);
}

/** Katalogowanie źródła — komenda `research.source.add` wraz z metadanymi. */
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
 * Trzy stany wykazu źródeł: pytam, mam treść, nie mam czego pokazać.
 *
 * Pustka wykazu to nie to samo, co brak miejsca na wykaz: przy niewskazanym
 * oknie badania `research.source.add` odmawia, bo `windowId` jest polem
 * obowiązkowym, więc zaproszenie do skatalogowania pierwszego źródła byłoby
 * wtedy mylące. Właściwe zdanie dobiera `pustka-okien.ts`.
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
 * Tekst swobodny okna przekazywany komendom bez własnego formularza.
 *
 * Żądanie składane bez wskazania Operatora wracałoby odmową walidacji, z której
 * nic dla niego nie wynika. Ten jeden krok mówi, skąd okno bierze treść — i gdy
 * jej nie ma, `wywolania-komend.ts` nazywa brak, zamiast wysyłać puste pole.
 */
function tekstDlaKomendy(kontekst: KontekstZrodel): string {
  // Pole tytułu formularza źródła niesie tekst swobodny tej rodziny: adres do
  // pozyskania, ścieżkę bibliografii albo etykiety po przecinku.
  return kontekst.formularz.zlecenie('').tytul.trim();
}
