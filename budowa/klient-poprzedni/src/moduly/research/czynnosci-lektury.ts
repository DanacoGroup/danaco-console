import {
  ResearchFindingStatus,
  type LibraryPreview,
  type ResearchSource,
} from '../../../../shared/contract';
import { opisOdmowyBledu } from '../../komponenty/odmowa';
import type { WierszOdpowiedzi } from '../../modele/kontrolki-formularza';
import type { AkcjaBadania } from './akcje-okien';
import { KODY_OKIEN } from './kody-okien';
import { pokazPustke, PUSTKA_LEKTURY } from './pustka-okien';
import type { StanBadania } from './stan-badania';
import { czyKomendaBadania, wykonajKomendeBadania } from './wywolania-komend';
import type { StanOknaBadania } from './stan-okna-badania';

/**
 * Czynności Operatora w Reading View: lektura źródła i zamiana zaznaczonego
 * fragmentu w ustalenie.
 *
 * Dwie czynności mają dziś pokrycie w kontrakcie i obie idą nim naprawdę:
 * wczytanie strony (`library.file.preview`) i zapis wypisu jako ustalenia
 * (`research.finding.add` z cytatem w polu `content` i źródłem w `sourceIds`).
 * Podświetlenia trwałe, notatki na marginesie, OCR, ekstrakcja tabel,
 * streszczenie źródła i stan lektury komendy nie mają — idą panelem akcji
 * i wracają odmową rdzenia.
 */

/** Górna granica długości podglądu jednej strony — pole `maxChars` kontraktu. */
const LIMIT_ZNAKOW = 20000;

export interface KontekstLektury {
  stan: StanBadania;
  okno: StanOknaBadania;
  odpowiedz: WierszOdpowiedzi;
  /** Obszar treści źródła — z niego pochodzi zaznaczony fragment. */
  tresc: HTMLElement;
  /** Numer strony w lekturze; liczony od jednego, tak jak `page` kontraktu. */
  strona(): number;
  /** Przestawia numer strony i wczytuje ją ponownie. */
  ustawStrone(numer: number): void;
  /** Wchłania podgląd oddany przez rdzeń albo jego brak. */
  wchlonPodglad(podglad: LibraryPreview | null): void;
  przejdz(kodOkna: string): void;
}

/** Rozdziela akcję panelu na drogę własną okna i drogę generyczną. */
export async function wykonajAkcjeLektury(
  kontekst: KontekstLektury,
  akcja: AkcjaBadania,
): Promise<void> {
  if (akcja.kod === 'research.reading.toFinding') {
    await zapiszWypis(kontekst);
    return;
  }
  if (akcja.kod === 'research.reading.reload') {
    await wczytajZrodlo(kontekst);
    return;
  }
  if (akcja.kod === 'research.reading.toFindings') {
    kontekst.przejdz(KODY_OKIEN.ustalenia);
    kontekst.odpowiedz.pokaz(
      'Wypisy zapisane z tego okna stoją w Findings Panel razem z ustaleniami wpisanymi ręcznie — to ta sama komenda i ten sam wykaz.',
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

/** Droga generyczna: `window.action` ze wskazanym źródłem i stroną w parametrach. */
async function przezPanelAkcji(kontekst: KontekstLektury, akcja: AkcjaBadania): Promise<void> {
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
    sourceIds: [stan.lektura.wskazane()],
    page: kontekst.strona(),
  });
  if (!wynik.udany) {
    odpowiedz.pokaz(opisOdmowyBledu(`Akcja „${akcja.nazwa}"`, wynik.blad, wynik.nieznanyTyp), false);
    return;
  }
  odpowiedz.pokaz(`Rdzeń oddał wynik akcji „${akcja.nazwa}".`, true);
}

/**
 * Wczytanie treści wskazanego źródła.
 *
 * Droga do treści prowadzi przez dokument repozytorium, więc źródło bez pola
 * `libraryFileId` nie ma czego pokazać — i to okno mówi wprost, zamiast
 * pokazywać pusty czytnik i zostawiać Operatora z domysłem, że materiał się nie
 * wczytał.
 */
export async function wczytajZrodlo(kontekst: KontekstLektury): Promise<void> {
  const { stan, okno, odpowiedz } = kontekst;
  const zrodlo = wskazaneZrodlo(stan);
  if (zrodlo === null) {
    kontekst.wchlonPodglad(null);
    odpowiedz.pokaz(
      'Nie wskazano, co czytać — naciśnij „Czytaj" przy pozycji w Sources Manager.',
      false,
    );
    return;
  }
  const idPliku = zrodlo.libraryFileId ?? '';
  if (idPliku === '') {
    kontekst.wchlonPodglad(null);
    okno.puste(
      'Źródło bez dokumentu repozytorium',
      `Źródło „${zrodlo.title}" nie wskazuje dokumentu repozytorium (pole libraryFileId jest puste). ` +
        'Własna komenda odczytu treści źródła jest już w kontrakcie, ale rdzeń nie ma dla niej ' +
        'uchwytu, więc przejezdny zostaje wyłącznie podgląd zasobu repozytorium — a materiału ' +
        'spoza repozytorium czytnik nie ma dziś skąd wziąć. Kontrolki zostają czynne.',
    );
    return;
  }

  okno.ladowanie(`Wczytywanie strony ${String(kontekst.strona())} źródła „${zrodlo.title}"…`);
  const wynik = await stan.lekturaZrodla.wczytajStrone(idPliku, kontekst.strona(), LIMIT_ZNAKOW);
  if (!wynik.udany || wynik.wynik === undefined) {
    const opis = opisOdmowyBledu('Wczytanie treści źródła', wynik.blad, wynik.nieznanyTyp);
    kontekst.wchlonPodglad(null);
    okno.blad(opis);
    odpowiedz.pokaz(opis, false);
    return;
  }
  kontekst.wchlonPodglad(wynik.wynik);
  odpowiedz.pokaz(opisPodgladu(wynik.wynik, zrodlo), true);
}

/**
 * Zamiana zaznaczonego fragmentu w ustalenie.
 *
 * To jest ta sama komenda, którą wysyła formularz Findings Panel — czytnik nie
 * ma własnej drogi zapisu i nie potrzebuje jej mieć. Wypis wchodzi z powiązaniem
 * do czytanego źródła, więc kotwica ustalenia sięga materiału, a nie samego
 * napisu.
 *
 * Kotwica fragmentu jest w kontrakcie połowicznie i to jest tu rzecz istotna:
 * żądanie zapisu ustalenia ma pole `anchor`, ale `ResearchFinding`, które rdzeń
 * oddaje, nie niesie ani numeru strony, ani zakresu znaków. Kotwicę da się więc
 * wysłać, a nie da się jej odczytać z powrotem — dopóki byt ustalenia jej nie
 * niesie, okno nie ma jak pokazać, dokąd cytat sięga. Numer strony wchodzi
 * zatem w TREŚĆ cytatu, żeby nie przepadł, i okno mówi o tym wprost.
 */
export async function zapiszWypis(kontekst: KontekstLektury): Promise<void> {
  const { stan, odpowiedz } = kontekst;
  const fragment = zaznaczonyFragment(kontekst.tresc);
  if (fragment === '') {
    odpowiedz.pokaz(
      'Zaznacz fragment w treści źródła — wypis bez cytatu byłby ustaleniem bez zakotwiczenia.',
      false,
    );
    return;
  }
  if (stan.idOkna() === '') {
    odpowiedz.pokaz(
      'Rdzeń nie wskazał okna badania; komenda research.finding.add wymaga pola windowId.',
      false,
    );
    return;
  }
  const zrodlo = wskazaneZrodlo(stan);
  const idZrodel = zrodlo === null ? [] : [zrodlo.id];

  odpowiedz.pokaz('Zapis wypisu jako ustalenia…', true);
  const wynik = await stan.zrodlo.zapiszUstalenie({
    idOkna: stan.idOkna(),
    tresc: trescWypisu(fragment, kontekst.strona(), zrodlo),
    idUstalenia: '',
    idZrodel,
    stan: ResearchFindingStatus.Open,
  });
  if (!wynik.udany || wynik.wynik === undefined) {
    odpowiedz.pokaz(opisOdmowyBledu('Zapis wypisu', wynik.blad, wynik.nieznanyTyp), false);
    return;
  }
  stan.wchlonUstalenie(wynik.wynik.finding);
  const zwiazane = (wynik.wynik.finding.sourceIds ?? []).length;
  odpowiedz.pokaz(
    `Rdzeń zapisał wypis jako ustalenie ${wynik.wynik.finding.id}` +
      (idZrodel.length === 0
        ? ' bez powiązania ze źródłem — czytane źródło nie było wskazane.'
        : zwiazane === 0
          ? '. Rdzeń NIE związał go z czytanym źródłem — powiązania nie ma w ustaleniu, które oddał.'
          : ' wraz z powiązaniem do czytanego źródła.'),
    idZrodel.length === 0 || zwiazane > 0,
  );
}

/** Źródło wskazane do lektury; `null` znaczy „nie wskazano albo już go nie ma". */
export function wskazaneZrodlo(stan: StanBadania): ResearchSource | null {
  const wskazane = stan.lektura.wskazane();
  if (wskazane === '') return null;
  return stan.zrodla().find((zrodlo) => zrodlo.id === wskazane) ?? null;
}

/**
 * Trzy stany czytnika: pytam, mam treść, nie mam czego czytać.
 *
 * Rozstrzygnięcie „źródło bez dokumentu repozytorium" nie należy tutaj — to
 * odpowiedź na konkretną lekturę, a nie stan okna sprzed niej; stawia je
 * `wczytajZrodlo`.
 */
export function ustawStanLektury(kontekst: KontekstLektury, maTresc: boolean): void {
  const { stan, okno } = kontekst;
  if (stan.faza() === 'odczyt') {
    okno.ladowanie('Odczyt okna badania z rdzenia…');
    return;
  }
  if (maTresc) {
    okno.gotowe();
    return;
  }
  if (stan.faza() === 'blad') {
    okno.blad(stan.powod());
    return;
  }
  pokazPustke(okno, stan, PUSTKA_LEKTURY);
}

/** Zaznaczenie Operatora ograniczone do obszaru treści źródła. */
function zaznaczonyFragment(tresc: HTMLElement): string {
  const zaznaczenie = document.getSelection();
  if (zaznaczenie === null || zaznaczenie.isCollapsed) return '';
  // Zaznaczenie spoza obszaru treści nie jest cytatem ze źródła: fragment
  // etykiety okna albo komunikatu zapisany jako wypis byłby cytatem zmyślonym.
  const kotwica = zaznaczenie.anchorNode;
  if (kotwica === null || !tresc.contains(kotwica)) return '';
  return zaznaczenie.toString().trim();
}

/**
 * Treść wypisu: cytat wraz z miejscem, z którego pochodzi.
 *
 * Numer strony idzie w treści, bo `ResearchFinding` nie ma pola na kotwicę
 * pozycji. Przemilczenie strony zostawiłoby cytat, którego nie da się odnaleźć
 * w materiale — a to jest dokładnie ta rzecz, którą ustalenie ma zapewniać.
 */
function trescWypisu(fragment: string, strona: number, zrodlo: ResearchSource | null): string {
  const skad =
    zrodlo === null
      ? `strona ${String(strona)}`
      : `${zrodlo.title}, strona ${String(strona)}`;
  return `„${fragment}" [${skad}]`;
}

/** Zdanie o wczytanym podglądzie — złożone z pól, które rdzeń naprawdę oddał. */
function opisPodgladu(podglad: LibraryPreview, zrodlo: ResearchSource): string {
  const czesci = [
    `Rdzeń wczytał podgląd źródła „${zrodlo.title}" (rodzaj: ${podglad.kind}).`,
    podglad.pageCount === undefined
      ? ''
      : `Strona ${String(podglad.page ?? 1)} z ${String(podglad.pageCount)}.`,
    podglad.truncated === true
      ? `Podgląd jest SKRÓCONY do ${String(LIMIT_ZNAKOW)} znaków — cytat spoza tego zakresu nie jest tu widoczny.`
      : '',
    (podglad.text ?? '') === ''
      ? 'Treści tekstowej rdzeń nie oddał; podgląd graficzny czytnik pokazuje odnośnikiem, a nie tekstem do zaznaczenia.'
      : '',
  ];
  return czesci.filter((czesc) => czesc !== '').join(' ');
}

/**
 * Tekst swobodny okna przekazywany komendom bez własnego formularza.
 *
 * Żądanie składane bez wskazania Operatora wracałoby odmową walidacji, z której
 * nic dla niego nie wynika. Ten jeden krok mówi, skąd okno bierze treść — i gdy
 * jej nie ma, `wywolania-komend.ts` nazywa brak, zamiast wysyłać puste pole.
 */
function tekstDlaKomendy(_kontekst: KontekstLektury): string {
  // Fragment zaznaczony w treści materiału jest tym, czego dotyczy adnotacja
  // i pytanie o korpus — bierzemy go z zaznaczenia okna, nie z osobnego pola.
  return (window.getSelection()?.toString() ?? '').trim();
}
