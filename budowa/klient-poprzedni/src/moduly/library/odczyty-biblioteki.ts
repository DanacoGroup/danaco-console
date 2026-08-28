import { Command, type KnowledgeHit, type LibraryFile } from '../../../../shared/contract';
import { opisOdmowy } from '../../komponenty/odmowa';
import type { MagazynBiblioteki, TrybWyszukiwania } from './magazyn-biblioteki';
import type { ZrodloBiblioteki } from './zrodlo-biblioteki';
import type { ZrodloZnaczenia } from './zrodlo-znaczenia';

// Drogi odczytu wykazu osobno od pamięci i widoku: trzy komendy, każda z osobnym wejściem okna.

/** Górna granica wyników wyszukiwania — okno pokazuje ograniczoną liczbę trafień, nie cały zbiór pasujących plików. */
const GRANICA_TRAFIEN = 50;

/**
 * Górna granica wykazu czytanego po to, żeby przypisać trafieniom znaczenia
 * wiersze plików. `knowledge.search` oddaje identyfikator źródła, a nie plik,
 * więc metryka wiersza musi przyjść drugim odczytem.
 */
const GRANICA_WYKAZU_ZNACZEN = 500;

export async function odczytajWykaz(
  magazyn: MagazynBiblioteki,
  zrodlo: ZrodloBiblioteki,
  fraza: string,
  etykieta: string,
): Promise<void> {
  zapowiedz(magazyn);
  const wynik = await zrodlo.wykaz({
    ...(fraza.trim() === '' ? {} : { query: fraza.trim() }),
    ...(etykieta.trim() === '' ? {} : { tags: [etykieta.trim()] }),
  });
  przyjmij(
    magazyn,
    wynik.wynik?.files,
    opisOdmowy('Odczyt wykazu plików', wynik.blad?.code, wynik.blad?.message),
  );
}

/**
 * Wyszukiwanie w trybie wybranym przez użytkownika zwraca zdanie o drodze, którą
 * wynik powstał, bo trzy tryby dają wyniki nieporównywalne między sobą.
 */
export async function szukajWykaz(
  magazyn: MagazynBiblioteki,
  zrodlo: ZrodloBiblioteki,
  wskaznik: ZrodloZnaczenia,
  fraza: string,
  tryb: TrybWyszukiwania,
): Promise<string> {
  if (tryb === 'pelnotekstowy') return szukajSlowami(magazyn, zrodlo, fraza);
  if (tryb === 'semantyczny') return szukajZnaczeniem(magazyn, zrodlo, wskaznik, fraza);
  return szukajHybrydowo(magazyn, zrodlo, wskaznik, fraza);
}

/** Dopasowanie słów w indeksie pełnotekstowym rdzenia, ograniczone górną granicą liczby trafień wyszukiwania. */
async function szukajSlowami(
  magazyn: MagazynBiblioteki,
  zrodlo: ZrodloBiblioteki,
  fraza: string,
): Promise<string> {
  zapowiedz(magazyn);
  const wynik = await zrodlo.szukaj(fraza.trim(), GRANICA_TRAFIEN);
  const powod = opisOdmowy('Wyszukiwanie w bibliotece', wynik.blad?.code, wynik.blad?.message);
  przyjmij(magazyn, wynik.wynik?.files, powod);
  if (wynik.wynik === undefined) return powod;
  return `Wyszukiwanie po słowach (library.file.search): trafień ${wynik.wynik.files.length}.`;
}

/**
 * Dopasowanie znaczenia w zakresie biblioteki: trafienie niesie identyfikator
 * źródła, a nie plik, więc wiersze biorą się z osobnego odczytu wykazu, a
 * trafienie bez pliku jest wypowiedziane, nie pominięte.
 */
async function szukajZnaczeniem(
  magazyn: MagazynBiblioteki,
  zrodlo: ZrodloBiblioteki,
  wskaznik: ZrodloZnaczenia,
  fraza: string,
): Promise<string> {
  zapowiedz(magazyn);
  const trafienia = await wskaznik.szukajZnaczeniem(fraza.trim(), GRANICA_TRAFIEN);
  if (!trafienia.udany || trafienia.wynik === undefined) {
    const powod = opisOdmowy(
      'Wyszukiwanie po znaczeniu',
      trafienia.blad?.code,
      trafienia.blad?.message,
    );
    przyjmij(magazyn, undefined, powod);
    return powod;
  }
  const wykaz = await zrodlo.wykaz({ limit: GRANICA_WYKAZU_ZNACZEN });
  if (!wykaz.udany || wykaz.wynik === undefined) {
    const powod = opisOdmowy(
      'Odczyt wykazu plików dla trafień znaczenia',
      wykaz.blad?.code,
      wykaz.blad?.message,
    );
    przyjmij(magazyn, undefined, powod);
    return powod;
  }
  const { pliki, bezPliku } = zlozTrafienia(trafienia.wynik.results, wykaz.wynik.files);
  zapiszZnaczenia(magazyn, trafienia.wynik.results);
  przyjmij(magazyn, pliki, '');
  return zdanieOZnaczeniu(trafienia.wynik.results.length, pliki.length, bezPliku);
}

/**
 * Tryb hybrydowy składa obie odpowiedzi po stronie okna: dopasowanie słów przed
 * dopasowaniem znaczenia, bo trafienia obu dróg nie są sumowane w jedną wspólną
 * trafność.
 */
async function szukajHybrydowo(
  magazyn: MagazynBiblioteki,
  zrodlo: ZrodloBiblioteki,
  wskaznik: ZrodloZnaczenia,
  fraza: string,
): Promise<string> {
  zapowiedz(magazyn);
  const [slowa, trafienia] = await Promise.all([
    zrodlo.szukaj(fraza.trim(), GRANICA_TRAFIEN),
    wskaznik.szukajZnaczeniem(fraza.trim(), GRANICA_TRAFIEN),
  ]);
  if (!slowa.udany || slowa.wynik === undefined) {
    const powod = opisOdmowy('Wyszukiwanie w bibliotece', slowa.blad?.code, slowa.blad?.message);
    przyjmij(magazyn, undefined, powod);
    return powod;
  }
  if (!trafienia.udany || trafienia.wynik === undefined) {
    // Połowa hybrydy padła, druga przyszła: zdanie mówi, że to już nie jest tryb hybrydowy.
    zapiszZnaczenia(magazyn, []);
    przyjmij(magazyn, slowa.wynik.files, '');
    return (
      `Wyszukiwanie po słowach: trafień ${slowa.wynik.files.length}. Znaczenie nie weszło ` +
      `do wyniku — ${opisOdmowy(Command.KnowledgeSearch, trafienia.blad?.code, trafienia.blad?.message)}`
    );
  }
  const wykaz = await zrodlo.wykaz({ limit: GRANICA_WYKAZU_ZNACZEN });
  const zeZnaczenia =
    wykaz.udany && wykaz.wynik !== undefined
      ? zlozTrafienia(trafienia.wynik.results, wykaz.wynik.files)
      : { pliki: [], bezPliku: trafienia.wynik.results.length };
  const znane = new Set(slowa.wynik.files.map((plik) => plik.id));
  const dolozone = zeZnaczenia.pliki.filter((plik) => !znane.has(plik.id));
  zapiszZnaczenia(magazyn, trafienia.wynik.results);
  przyjmij(magazyn, [...slowa.wynik.files, ...dolozone], '');
  return (
    `Wyszukiwanie hybrydowe: po słowach ${slowa.wynik.files.length}, po znaczeniu ` +
    `dołożono ${dolozone.length} spoza tamtego zbioru` +
    (zeZnaczenia.bezPliku === 0
      ? '.'
      : `; trafień znaczenia bez pliku w wykazie: ${zeZnaczenia.bezPliku}.`)
  );
}

/** Odwzorowanie trafień wskaźnika na wiersze wykazu plików, zachowane w kolejności malejącej trafności. */
function zlozTrafienia(
  trafienia: readonly KnowledgeHit[],
  wykaz: readonly LibraryFile[],
): { pliki: LibraryFile[]; bezPliku: number } {
  const poKodzie = new Map(wykaz.map((plik) => [plik.id, plik]));
  const pliki: LibraryFile[] = [];
  const wziete = new Set<string>();
  let bezPliku = 0;
  for (const trafienie of trafienia) {
    const kod = trafienie.sourceId ?? '';
    const plik = poKodzie.get(kod);
    if (plik === undefined) {
      bezPliku += 1;
      continue;
    }
    // Wskaźnik dzieli dokument na fragmenty, jeden plik trafia wielokrotnie; liczy się trafienie pierwsze.
    if (wziete.has(plik.id)) continue;
    wziete.add(plik.id);
    pliki.push(plik);
  }
  return { pliki, bezPliku };
}

/** Odkłada trafność i fragment przy pliku w magazynie stanu; poprzednie trafienia znikają przy każdym zapisie. */
function zapiszZnaczenia(magazyn: MagazynBiblioteki, trafienia: readonly KnowledgeHit[]): void {
  magazyn.znaczenia.clear();
  for (const trafienie of trafienia) {
    const kod = trafienie.sourceId ?? '';
    if (kod === '' || magazyn.znaczenia.has(kod)) continue;
    magazyn.znaczenia.set(kod, {
      trafnosc: trafienie.score ?? null,
      fragment: trafienie.text,
    });
  }
}

/** Zdanie o wyniku wyszukiwania po znaczeniu — trzy różne stany wyniku dostają trzy osobne zdania opisowe. */
function zdanieOZnaczeniu(trafien: number, plikow: number, bezPliku: number): string {
  if (trafien === 0) {
    return (
      'Wyszukiwanie po znaczeniu (knowledge.search): rdzeń nie oddał ani jednego fragmentu. ' +
      'Wskaźnik znaczenia nie buduje się przy wgraniu pliku — przelicz go w zakładce Higiena ' +
      'panelu Metadata & Archive Panel, jeżeli repozytorium ma treść.'
    );
  }
  if (plikow === 0) {
    return (
      `Wyszukiwanie po znaczeniu: fragmentów ${trafien}, ale żaden nie wskazuje pliku ` +
      'obecnego w wykazie. Wskaźnik opisuje stan sprzed ostatniej zmiany repozytorium.'
    );
  }
  return (
    `Wyszukiwanie po znaczeniu: fragmentów ${trafien}, plików ${plikow}` +
    (bezPliku === 0 ? '.' : `, trafień bez pliku w wykazie ${bezPliku}.`)
  );
}

/** Zapowiedź trwającego odczytu — okno ma pokazać stan pytania rdzenia, a nie stan pustego wykazu plików. */
function zapowiedz(magazyn: MagazynBiblioteki): void {
  magazyn.faza = 'odczyt';
  magazyn.powod = '';
  magazyn.oglos();
}

/** Wspólne zakończenie odczytów wykazu: jeden powód niepowodzenia i jedno wspólne ogłoszenie zmiany magazynu. */
function przyjmij(
  magazyn: MagazynBiblioteki,
  pliki: LibraryFile[] | undefined,
  powod: string,
): void {
  if (pliki === undefined) {
    magazyn.faza = 'blad';
    magazyn.powod = powod;
    magazyn.oglos();
    return;
  }
  magazyn.zbior = pliki;
  magazyn.zaznaczenie = magazyn.zaznaczenie.filter((id) => pliki.some((plik) => plik.id === id));
  if (magazyn.wskazany !== null && !pliki.some((plik) => plik.id === magazyn.wskazany)) {
    magazyn.wskazany = null;
  }
  // Nowy odczyt orzeka o całym wykazie, więc wcześniejsze zawężenie przestaje o nim mówić.
  magazyn.zawezenie = null;
  magazyn.faza = 'gotowe';
  magazyn.powod = '';
  magazyn.oglos();
}
