import type { LibraryFile } from '../../../../shared/contract';
import type { StanTresci } from './dostepnosc-tresci';
import { rodzinaTresci, utworzKafelPliku } from './kafel-pliku';
import type { TrafienieZnaczenia, WidokWykazu } from './magazyn-biblioteki';
import { utworzWierszPliku } from './wiersz-pliku';

/**
 * Pięć form prezentacji tego samego wykazu plików.
 *
 * Widok przelicza wyłącznie prezentację odpowiedzi, którą rdzeń już oddał —
 * żaden z nich nie wysyła komendy i żaden nie zawęża zbioru inaczej niż tym,
 * co niesie sam plik. Dlatego wybór widoku nie unieważnia zaznaczenia ani
 * wskazania pliku czynnego.
 *
 * Widok, który nie ma czego pokazać mimo niepustego wykazu, oddaje zdanie
 * w polu `brak` zamiast pustego prostokąta. Rozróżnienie jest tu istotne:
 * „galeria bez obrazów" i „repozytorium bez plików" to dwa różne stany, a widok
 * mapy nie ma dziś źródła współrzędnych w ogóle.
 */
export interface OpisWidoku {
  zaznaczone: readonly string[];
  czynny: string | null;
  tresc(idPliku: string): StanTresci;
  znaczenie(idPliku: string): TrafienieZnaczenia | null;
  naZaznaczenie(idPliku: string): void;
  naWskazanie(idPliku: string): void;
}

export interface WynikWidoku {
  /** Element osadzany w ciele wykazu. */
  element: HTMLElement;
  /** Zdanie, gdy widok nie ma czego pokazać mimo niepustego wykazu. */
  brak: string;
}

/** Nazwy widoków w przełączniku — kolejność jak w dokumentacji modułu. */
export const WIDOKI: ReadonlyArray<{ kod: WidokWykazu; nazwa: string }> = [
  { kod: 'siatka', nazwa: 'Siatka' },
  { kod: 'lista', nazwa: 'Lista' },
  { kod: 'galeria', nazwa: 'Galeria' },
  { kod: 'os-czasu', nazwa: 'Oś czasu' },
  { kod: 'mapa', nazwa: 'Mapa' },
];

export function zbudujWidok(
  widok: WidokWykazu,
  pliki: readonly LibraryFile[],
  opis: OpisWidoku,
): WynikWidoku {
  if (widok === 'lista') return { element: listaSzczegolowa(pliki, opis), brak: '' };
  if (widok === 'siatka') return { element: siatkaMiniatur(pliki, opis), brak: '' };
  if (widok === 'galeria') return galeria(pliki, opis);
  if (widok === 'os-czasu') return osCzasu(pliki, opis);
  return mapa(pliki);
}

/** Widok listy szczegółowej — wiersz z metryką, ścieżką i etykietami. */
function listaSzczegolowa(pliki: readonly LibraryFile[], opis: OpisWidoku): HTMLElement {
  const element = document.createElement('ul');
  element.className = 'ml-pliki';
  element.append(...pliki.map((plik) => wiersz(plik, opis)));
  return element;
}

/** Widok domyślny: siatka kafli, po jednym na plik. */
function siatkaMiniatur(pliki: readonly LibraryFile[], opis: OpisWidoku): HTMLElement {
  const element = document.createElement('ul');
  element.className = 'ml-siatka';
  element.append(...pliki.map((plik) => kafel(plik, opis)));
  return element;
}

/**
 * Widok galerii — wyłącznie materiał graficzny.
 *
 * Zawężenie idzie po `mimeType`, jedynym polu kontraktu mówiącym o rodzaju
 * treści. Plik bez tego pola do galerii nie wchodzi i jest to wypowiedziane:
 * wciągnięcie go „na wszelki wypadek" stawiałoby w galerii dokumenty.
 */
function galeria(pliki: readonly LibraryFile[], opis: OpisWidoku): WynikWidoku {
  const obrazy = pliki.filter((plik) => rodzinaTresci(plik.mimeType) === 'image');
  const element = document.createElement('ul');
  element.className = 'ml-siatka ml-siatka--galeria';
  element.append(...obrazy.map((plik) => kafel(plik, opis)));
  if (obrazy.length > 0) return { element, brak: '' };
  const bezRodzaju = pliki.filter((plik) => (plik.mimeType ?? '') === '').length;
  return {
    element,
    brak:
      `Wykaz liczy plików: ${pliki.length}, a żaden nie ma rodzaju treści z rodziny obrazów. ` +
      (bezRodzaju === 0
        ? 'Galeria pokazuje wyłącznie materiał graficzny — przełącz widok na siatkę albo listę.'
        : `Plików bez pola mimeType: ${bezRodzaju} — o ich rodzaju rdzeń nic nie powiedział, ` +
          'więc galeria ich nie zgaduje.'),
  };
}

/**
 * Widok osi czasu — chronologia dodania do repozytorium.
 *
 * Grupowanie idzie po dacie z `createdAt`, bo oś czasu opisuje napływ, a nie
 * ostatnią zmianę. Doba jest jednostką najmniejszą, którą da się nazwać bez
 * ustawienia strefy: znacznik kontraktu jest liczbą milisekund epoki, a
 * przeglądarka zna strefę Operatora.
 */
function osCzasu(pliki: readonly LibraryFile[], opis: OpisWidoku): WynikWidoku {
  const element = document.createElement('div');
  element.className = 'ml-os-czasu';

  const doby = new Map<string, LibraryFile[]>();
  for (const plik of [...pliki].sort((pierwszy, drugi) => drugi.createdAt - pierwszy.createdAt)) {
    const doba = new Date(plik.createdAt).toLocaleDateString('pl');
    const zbior = doby.get(doba);
    if (zbior === undefined) doby.set(doba, [plik]);
    else zbior.push(plik);
  }

  for (const [doba, zbior] of doby) {
    const naglowek = document.createElement('h4');
    naglowek.className = 'ml-os-czasu__doba';
    naglowek.textContent = `${doba} — plików ${zbior.length}`;

    const lista = document.createElement('ul');
    lista.className = 'ml-pliki';
    lista.append(...zbior.map((plik) => wiersz(plik, opis)));

    const grupa = document.createElement('section');
    grupa.className = 'ml-os-czasu__grupa';
    grupa.append(naglowek, lista);
    element.append(grupa);
  }
  return { element, brak: '' };
}

/**
 * Widok mapy — dziś bez źródła współrzędnych.
 *
 * Dokumentacja modułu opiera mapę na geolokalizacji z metadanych EXIF. Kontrakt
 * nie niesie ani pola współrzędnych przy pliku (`LibraryFile`), ani komendy
 * oddającej metadane techniczne zasobu biblioteki, więc okno nie ma czego
 * nanieść. Widok zostaje w przełączniku i mówi to wprost — pozycja usunięta
 * z przełącznika wyglądałaby na widok, którego nigdy nie było.
 */
function mapa(pliki: readonly LibraryFile[]): WynikWidoku {
  const element = document.createElement('div');
  element.className = 'ml-mapa';
  return {
    element,
    brak:
      `Widok mapy nie ma dziś czego nanieść. Współrzędne bierze się z metadanych EXIF, ` +
      'a kontrakt nie niesie ich ani przy pliku repozytorium, ani żadną komendą odczytu ' +
      `metadanych technicznych zasobu biblioteki. Wykaz liczy plików: ${pliki.length} — ` +
      'żeby zobaczyć je na osi czasu albo w siatce, przełącz widok.',
  };
}

/** Jeden wiersz wykazu wraz z jego czynnościami. */
function wiersz(plik: LibraryFile, opis: OpisWidoku): HTMLElement {
  return utworzWierszPliku(plik, {
    zaznaczony: opis.zaznaczone.includes(plik.id),
    czynny: opis.czynny === plik.id,
    tresc: opis.tresc(plik.id),
    znaczenie: opis.znaczenie(plik.id),
    naZaznaczenie: () => opis.naZaznaczenie(plik.id),
    naWskazanie: () => opis.naWskazanie(plik.id),
  });
}

/** Jeden kafel siatki albo galerii wraz z jego czynnościami. */
function kafel(plik: LibraryFile, opis: OpisWidoku): HTMLElement {
  return utworzKafelPliku(plik, {
    zaznaczony: opis.zaznaczone.includes(plik.id),
    czynny: opis.czynny === plik.id,
    znaczenie: opis.znaczenie(plik.id),
    naZaznaczenie: () => opis.naZaznaczenie(plik.id),
    naWskazanie: () => opis.naWskazanie(plik.id),
  });
}
