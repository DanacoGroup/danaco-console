import { Command, ResearchCredibility } from '../../../../shared/contract';
import { opisOdmowyBledu } from '../../komponenty/odmowa';
import type { WierszOdpowiedzi } from '../../modele/kontrolki-formularza';
import type { AkcjaBadania } from './akcje-okien';
import { KODY_OKIEN } from './kody-okien';
import { pokazPustke, PUSTKA_ODKRYWANIA } from './pustka-okien';
import type { StanBadania } from './stan-badania';
import { czyKomendaBadania, wykonajKomendeBadania } from './wywolania-komend';
import type { StanOknaBadania } from './stan-okna-badania';
import { zlecenieZWyniku, type WynikOdkrycia } from './wynik-odkrycia';

/**
 * Czynności Operatora w Discovery Panel: wyszukanie i przeniesienie pozycji do
 * katalogu źródeł.
 *
 * Tryb rozstrzyga o drodze. Dwa tryby mają własną komendę kontraktu i idą nią
 * naprawdę; dwa pozostałe komendy nie mają i idą generycznym `window.action`,
 * którego odmowę okno wypisuje słowami rdzenia. Rozdział stoi tutaj, nie
 * w widoku — widok zbiera zapytanie, czynność wie, czym się ono kończy.
 */

/** Tryb zapytania Discovery Panel — cztery z opracowania modułu, rozdz. 3.4. */
export type TrybOdkrywania = 'znaczenie' | 'tresc' | 'web' | 'naukowy';

/** Górna granica liczby pozycji jednego zapytania. */
const LIMIT_WYNIKOW = 20;

export interface KontekstOdkrywania {
  stan: StanBadania;
  okno: StanOknaBadania;
  odpowiedz: WierszOdpowiedzi;
  zapytanie: HTMLInputElement;
  /** Tryb wybrany sterem nastawy panelu. */
  tryb(): TrybOdkrywania;
  /** Pozycje ostatniego wyniku; puste znaczy „jeszcze nie pytano". */
  wyniki(): readonly WynikOdkrycia[];
  /** Wymienia pozycje wyniku i przerysowuje wykaz. */
  ustawWyniki(pozycje: readonly WynikOdkrycia[]): void;
  przejdz(kodOkna: string): void;
}

/** Rozdziela akcję panelu na drogę własną okna i drogę generyczną. */
export async function wykonajAkcjeOdkrywania(
  kontekst: KontekstOdkrywania,
  akcja: AkcjaBadania,
): Promise<void> {
  if (akcja.kod === 'research.discovery.run') {
    await wyszukaj(kontekst);
    return;
  }
  if (akcja.kod === 'research.discovery.toSources') {
    kontekst.przejdz(KODY_OKIEN.zrodla);
    kontekst.odpowiedz.pokaz(
      'Pozycje przeniesione do katalogu stoją w Sources Manager — panel odkrywania trzyma sam wynik zapytania, nie materiał badania.',
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

/** Droga generyczna: `window.action` z zapytaniem i trybem w parametrach. */
async function przezPanelAkcji(
  kontekst: KontekstOdkrywania,
  akcja: AkcjaBadania,
): Promise<void> {
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
    query: kontekst.zapytanie.value.trim(),
    mode: kontekst.tryb(),
  });
  if (!wynik.udany) {
    odpowiedz.pokaz(opisOdmowyBledu(`Akcja „${akcja.nazwa}"`, wynik.blad, wynik.nieznanyTyp), false);
    return;
  }
  odpowiedz.pokaz(`Rdzeń oddał wynik akcji „${akcja.nazwa}".`, true);
}

/**
 * Wyszukanie w trybie bieżącym.
 *
 * Tryb semantyczny i pełnotekstowy idą komendami, których rdzeń słucha, więc
 * oddają wynik. Tryb webowy i naukowy mają w kontrakcie własną komendę —
 * `research.discovery.search`, gdzie o trybie rozstrzyga pole `mode` — ale rdzeń
 * nie ma dla niej jeszcze uchwytu, więc zapytanie idzie zgłoszeniem
 * `window.action` pod jej nazwą i wraca odmową. Okno nie podstawia pod nie
 * wyszukiwania semantycznego: dwa różne pytania oddające ten sam wynik byłyby
 * zmyśleniem zdolności, której moduł nie ma.
 */
export async function wyszukaj(kontekst: KontekstOdkrywania): Promise<void> {
  const { stan, okno, odpowiedz } = kontekst;
  const zapytanie = kontekst.zapytanie.value.trim();
  if (zapytanie === '') {
    kontekst.zapytanie.focus();
    odpowiedz.pokaz('Wpisz zapytanie — rdzeń odmówi wyszukania bez treści pola query.', false);
    return;
  }

  const tryb = kontekst.tryb();
  if (tryb === 'web' || tryb === 'naukowy') {
    await przezPanelAkcji(kontekst, {
      kod: Command.ResearchDiscoverySearch,
      nazwa: tryb === 'web' ? 'Wyszukiwanie webowe' : 'Wyszukiwanie naukowe',
      droga: 'komenda',
      objasnienie: '',
    });
    return;
  }

  okno.ladowanie(
    tryb === 'znaczenie'
      ? 'Wyszukiwanie po znaczeniu w wiedzy Operatora…'
      : 'Wyszukiwanie pełnotekstowe w zasobach repozytorium…',
  );
  const wynik =
    tryb === 'znaczenie'
      ? await stan.odkrywanie.poZnaczeniu(zapytanie, stan.idOkna(), LIMIT_WYNIKOW)
      : await stan.odkrywanie.poTresci(zapytanie, LIMIT_WYNIKOW);

  if (!wynik.udany || wynik.wynik === undefined) {
    const opis = opisOdmowyBledu('Wyszukiwanie źródeł', wynik.blad, wynik.nieznanyTyp);
    okno.blad(opis);
    odpowiedz.pokaz(opis, false);
    return;
  }
  kontekst.ustawWyniki(wynik.wynik);
  odpowiedz.pokaz(
    wynik.wynik.length === 0
      ? `Rdzeń nie znalazł ani jednej pozycji dla zapytania „${zapytanie}".`
      : `Rdzeń oddał ${String(wynik.wynik.length)} pozycji dla zapytania „${zapytanie}".`,
    true,
  );
}

/**
 * Przeniesienie pozycji wyniku do katalogu źródeł badania.
 *
 * To jest ta sama komenda, którą wysyła formularz Sources Manager — panel
 * odkrywania nie ma własnej drogi wstawiania i nie potrzebuje jej mieć.
 * Wiarygodność wchodzi jako „niezweryfikowane", bo pozycja wyszukiwania nie
 * niesie oceny, a wpisanie tu wartości innej niż ta byłoby oceną zmyśloną.
 */
export async function przeniesDoZrodel(
  kontekst: KontekstOdkrywania,
  pozycja: WynikOdkrycia,
): Promise<void> {
  const { stan, odpowiedz } = kontekst;
  if (stan.idOkna() === '') {
    odpowiedz.pokaz(
      'Rdzeń nie wskazał okna badania; komenda research.source.add wymaga pola windowId.',
      false,
    );
    return;
  }
  odpowiedz.pokaz(`Przenoszenie pozycji „${pozycja.tytul}" do katalogu źródeł…`, true);
  const wynik = await stan.zrodlo.dodajZrodlo(
    zlecenieZWyniku(pozycja, stan.idOkna(), ResearchCredibility.Unverified),
  );
  if (!wynik.udany || wynik.wynik === undefined) {
    odpowiedz.pokaz(
      opisOdmowyBledu('Przeniesienie pozycji do źródeł', wynik.blad, wynik.nieznanyTyp),
      false,
    );
    return;
  }
  stan.wchlonZrodlo(wynik.wynik.source);
  odpowiedz.pokaz(
    `Rdzeń skatalogował źródło „${wynik.wynik.source.title}" — pozycja jest w Sources Manager.`,
    true,
  );
}

/**
 * Trzy stany panelu: pytam, mam wynik, nie mam czego pokazać.
 *
 * Pustka wyniku to nie pustka badania: zapytanie bez trafień jest odpowiedzią
 * rdzenia, a nie brakiem materiału, więc zdanie mówi o zapytaniu. Zdanie
 * o samym panelu — sprzed pierwszego zapytania — stoi w `pustka-okien.ts`.
 */
export function ustawStanOdkrywania(kontekst: KontekstOdkrywania, liczba: number): void {
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
  pokazPustke(okno, stan, PUSTKA_ODKRYWANIA);
}

/**
 * Tekst swobodny okna przekazywany komendom bez własnego formularza.
 *
 * Żądanie składane bez wskazania Operatora wracałoby odmową walidacji, z której
 * nic dla niego nie wynika. Ten jeden krok mówi, skąd okno bierze treść — i gdy
 * jej nie ma, `wywolania-komend.ts` nazywa brak, zamiast wysyłać puste pole.
 */
function tekstDlaKomendy(kontekst: KontekstOdkrywania): string {
  return kontekst.zapytanie.value.trim();
}
