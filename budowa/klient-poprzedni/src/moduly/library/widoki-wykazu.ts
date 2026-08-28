import type { LibraryFile } from '../../../../shared/contract';
import type { StanTresci } from './dostepnosc-tresci';
import { rodzinaTresci, utworzKafelPliku } from './kafel-pliku';
import type { TrafienieZnaczenia, WidokWykazu } from './magazyn-biblioteki';
import { utworzWierszPliku } from './wiersz-pliku';

/**
 * Pięć form prezentacji tego samego wykazu plików przelicza wyłącznie
 * odpowiedź, którą rdzeń już oddał, i żadna nie wysyła komendy ani nie
 * unieważnia zaznaczenia.
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

/** Nazwy widoków w przełączniku, w kolejności pokazywanej użytkownikowi, od siatki miniatur po mapę zasobów. */
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

/** Widok listy szczegółowej — wiersz z metryką, ścieżką i etykietami, czytelny przy pełnej szerokości okna. */
function listaSzczegolowa(pliki: readonly LibraryFile[], opis: OpisWidoku): HTMLElement {
  const element = document.createElement('ul');
  element.className = 'ml-pliki';
  element.append(...pliki.map((plik) => wiersz(plik, opis)));
  return element;
}

/** Widok domyślny: siatka kafli, po jednym na każdy plik, każdy z miniaturą i skróconą nazwą własną zasobu. */
function siatkaMiniatur(pliki: readonly LibraryFile[], opis: OpisWidoku): HTMLElement {
  const element = document.createElement('ul');
  element.className = 'ml-siatka';
  element.append(...pliki.map((plik) => kafel(plik, opis)));
  return element;
}

/**
 * Widok galerii pokazuje wyłącznie materiał graficzny: zawężenie idzie po
 * rodzaju treści, a plik bez tego pola do galerii nie wchodzi, co jest
 * wypowiedziane wprost.
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
 * Widok osi czasu grupuje pliki po dacie dodania do repozytorium, bo opisuje
 * napływ, a nie ostatnią zmianę, z dobą jako najmniejszą jednostką bez
 * ustawienia strefy.
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
 * Widok mapy nie ma dziś źródła współrzędnych: kontrakt nie niesie pola
 * współrzędnych ani komendy metadanych technicznych, więc widok zostaje
 * w przełączniku i mówi to wprost.
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

/** Jeden wiersz wykazu wraz z jego czynnościami zaznaczenia, wskazania i odczytanym stanem treści pliku. */
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

/** Jeden kafel siatki albo galerii wraz z jego czynnościami zaznaczenia i wskazania pliku jako czynnego. */
function kafel(plik: LibraryFile, opis: OpisWidoku): HTMLElement {
  return utworzKafelPliku(plik, {
    zaznaczony: opis.zaznaczone.includes(plik.id),
    czynny: opis.czynny === plik.id,
    znaczenie: opis.znaczenie(plik.id),
    naZaznaczenie: () => opis.naZaznaczenie(plik.id),
    naWskazanie: () => opis.naWskazanie(plik.id),
  });
}
