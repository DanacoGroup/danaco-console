import { opisOdmowyBledu } from '../../komponenty/odmowa';
import type { WierszOdpowiedzi } from '../../modele/kontrolki-formularza';
import type { AkcjaBadania } from './akcje-okien';
import { KODY_OKIEN } from './kody-okien';
import type { KreatorRaportu } from './kreator-raportu';
import { pokazPustke, PUSTKA_RAPORTU } from './pustka-okien';
import { opisZlozenia } from './skutek-zlozenia';
import type { StanBadania } from './stan-badania';
import { czyKomendaBadania, wykonajKomendeBadania } from './wywolania-komend';
import type { StanOknaBadania } from './stan-okna-badania';

/**
 * Czynności Operatora w Report Builderze.
 *
 * Jedna odpowiedzialność: kompozycja raportu i rozdział akcji panelu.
 * Materiał wejściowy pochodzi z Findings Panel — pole `findingIds` bierze
 * zaznaczenie ustaleń wspólne obu oknom przez `stan-badania`.
 */
export interface KontekstRaportu {
  stan: StanBadania;
  okno: StanOknaBadania;
  odpowiedz: WierszOdpowiedzi;
  kreator: KreatorRaportu;
  przejdz(kodOkna: string): void;
}

/** Rozdziela akcję panelu na drogę własną okna i drogę generyczną. */
export async function wykonajAkcjeRaportu(
  kontekst: KontekstRaportu,
  akcja: AkcjaBadania,
): Promise<void> {
  if (akcja.kod === 'research.report.section.add') {
    kontekst.kreator.otworz(kontekst.stan.zakres());
    kontekst.kreator.sekcje.dodaj();
    return;
  }
  if (akcja.kod === 'research.report.toExport') {
    kontekst.przejdz(KODY_OKIEN.eksport);
    kontekst.odpowiedz.pokaz(
      kontekst.stan.raport() === null
        ? 'Export Panel jest czynny zawsze; bez złożonego raportu powie o brakującym warunku, nie zablokuje przycisku.'
        : 'Export Panel weźmie raport złożony w tym oknie.',
      true,
    );
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

/** Droga generyczna: `window.action` z identyfikatorem raportu w parametrach. */
async function przezPanelAkcji(kontekst: KontekstRaportu, akcja: AkcjaBadania): Promise<void> {
  const { stan, odpowiedz } = kontekst;
  if (stan.idOkna() === '') {
    odpowiedz.pokaz(
      `Akcja „${akcja.nazwa}" wymaga okna badania, którego rdzeń jeszcze nie wskazał.`,
      false,
    );
    return;
  }
  odpowiedz.pokaz(`Wysłano window.action ${akcja.kod}…`, true);
  const wynik = await stan.okna.akcja(stan.idOkna(), akcja.kod, {
    reportId: stan.raport()?.id ?? '',
  });
  if (!wynik.udany) {
    odpowiedz.pokaz(opisOdmowyBledu(`Akcja „${akcja.nazwa}"`, wynik.blad, wynik.nieznanyTyp), false);
    return;
  }
  // Zdanie mówi o wyniku oddanym przez rdzeń, nie o wykonaniu akcji: rdzeń
  // kwituje sukcesem samo przyjęcie zgłoszenia. Odmowę braku wykonawcy
  // pokazuje gałąź wyżej, słowami rdzenia.
  odpowiedz.pokaz(`Rdzeń oddał wynik akcji „${akcja.nazwa}".`, true);
}

/**
 * Złożenie raportu; przycisk kreatora zostaje czynny przez cały czas.
 *
 * Rdzeń ma dwie drogi budowy, a rozstrzyga o nich redakcja:
 * `research.report.build` z polem `sections` składa raport z samych podanych
 * sekcji i odkłada `findingIds` na bok; to samo żądanie bez pola `sections`
 * woła kanał modelu po sekcję nadrzędną, a po niej idzie sekcja na każde
 * zaznaczone ustalenie. Okno nie obchodzi tej reguły po swojej stronie —
 * nazywa ją Operatorowi, żeby zaznaczenie ustaleń nie znikało bez słowa.
 *
 * Droga modelu kończy się sukcesem także wtedy, gdy model nic nie powiedział:
 * bez czynnego logowania do programu `claude` budowa wraca ze `status: "ok"`,
 * raport zostaje zapisany, a treścią sekcji nadrzędnej jest komunikat procesu
 * kanału. Rdzeń zbiera z kanału same fragmenty `text`, więc odpowiedź modelu
 * i komunikat jego procesu docierają tą samą drogą i jako ta sama treść.
 * Czynność nie orzeka o tym po napisie — skutek nazywa `skutek-zlozenia.ts`,
 * z odpowiedzi.
 */
export async function zlozRaport(kontekst: KontekstRaportu): Promise<void> {
  const { stan, kreator, odpowiedz } = kontekst;
  if (stan.idOkna() === '') {
    kreator.stan.blad(
      'Rdzeń nie wskazał okna badania; komenda research.report.build wymaga pola windowId.',
    );
    return;
  }
  const ustalenia = stan.wybraneUstalenia.wybrane();
  // Wiersze puste nie liczą się do niczego: kreator dokłada jeden przy każdym
  // otwarciu, więc licząc wiersze zamiast sekcji okno przepuszczałoby budowę
  // pustego dokumentu i kwitowało ją sukcesem.
  const sekcje = kreator.sekcje.zebrane();
  if (ustalenia.length === 0 && sekcje.length === 0) {
    kreator.stan.blad(
      'Raport byłby pusty: nie ma ani zaznaczonych ustaleń, ani wypełnionej sekcji — pusty wiersz redakcji sekcją nie jest. Zaznacz ustalenia w Findings Panel albo wpisz tytuł lub treść sekcji. Kreator zostaje otwarty.',
    );
    return;
  }
  // Komunikat czekania nazywa drogę, którą pójdzie rdzeń, a nie skutek, którego
  // jeszcze nie ma: bez czynnego logowania rdzeń streszczenia modelu nie
  // dostarcza.
  kreator.stan.ladowanie(
    sekcje.length === 0
      ? 'Składanie raportu — rdzeń woła kanał modelu po streszczenie zaznaczonych ustaleń…'
      : 'Składanie raportu z sekcji redakcji…',
  );
  const wynik = await stan.zrodlo.zlozRaport({
    idOkna: stan.idOkna(),
    idRaportu: stan.raport()?.id ?? '',
    tytul: kreator.tytul(),
    idUstalen: ustalenia,
    sekcje,
  });
  if (!wynik.udany || wynik.wynik === undefined) {
    kreator.stan.blad(opisOdmowyBledu('Złożenie raportu', wynik.blad, wynik.nieznanyTyp));
    return;
  }
  stan.wchlonRaport(wynik.wynik.report);
  // Redakcja przejmuje sekcje potwierdzone przez rdzeń — razem z sekcjami,
  // które rdzeń napisał sam, i z wiązaniem sekcja↔ustalenia. Bez tego kolejne
  // złożenie wysyłałoby zastane wiersze kreatora i wycierało to, co rdzeń
  // dopiero co zbudował.
  kreator.sekcje.wczytaj(wynik.wynik.report.sections ?? []);
  kreator.stan.gotowe();
  kreator.zamknij();
  // Zdanie o skutku powstaje z raportu, który wrócił, porównanego
  // z zamówieniem — patrz `skutek-zlozenia.ts`. Liczby policzone przed
  // wysłaniem nie wystarczą: nie mówią nic o tym, co rdzeń dopisał sam.
  const skutek = opisZlozenia(wynik.wynik.report, sekcje, ustalenia);
  odpowiedz.pokaz(skutek.zdanie, skutek.udany);
}

/**
 * Trzy stany podglądu raportu: pytam, mam dokument, nie mam czego pokazać.
 *
 * Zdanie pustki dobiera `pustka-okien.ts`: przy niewskazanym oknie badania
 * `research.report.build` odmawia i żadna z dwóch dróg budowy nie ruszy, więc
 * tłumaczenie tam reguły dwóch dróg myliłoby Operatora.
 */
export function ustawStanRaportu(kontekst: KontekstRaportu): void {
  const { stan, okno } = kontekst;
  if (stan.faza() === 'odczyt') {
    okno.ladowanie('Odczyt okna badania z rdzenia…');
    return;
  }
  if (stan.raport() !== null) {
    okno.gotowe();
    return;
  }
  if (stan.faza() === 'blad') {
    okno.blad(stan.powod());
    return;
  }
  pokazPustke(okno, stan, PUSTKA_RAPORTU);
}

/**
 * Tekst swobodny okna przekazywany komendom bez własnego formularza.
 *
 * Żądanie składane bez wskazania Operatora wracałoby odmową walidacji, z której
 * nic dla niego nie wynika. Ten jeden krok mówi, skąd okno bierze treść — i gdy
 * jej nie ma, `wywolania-komend.ts` nazywa brak, zamiast wysyłać puste pole.
 */
function tekstDlaKomendy(_kontekst: KontekstRaportu): string {
  // Kreator raportu nie ma pola swobodnego; komendy tej rodziny pracują na
  // raporcie bieżącym, którego wskazanie niesie stan badania.
  return '';
}
