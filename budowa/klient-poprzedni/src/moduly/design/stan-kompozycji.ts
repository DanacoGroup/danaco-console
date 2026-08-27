import type { DesignAsset, DesignBoardLayer } from '../../../../shared/contract';
import {
  dolozElement,
  dolozRamke,
  dolozZasob,
  naniesWarstwy,
  objeteWarstwy,
  przestawZaznaczenie,
  pustyZapis,
  rozmiescWarstwy,
  ulozWedlugSzablonu,
  usunWarstwe,
  wczytajUklad,
  wyrownajWarstwy,
  zmienWarstwe,
  type OsWyrownania,
  type SzablonUkladu,
  type WidokKanwy,
  type ZapisKompozycji,
} from './zapis-kompozycji';

// Kompozycja Design Board widziana przez okno: zapis układu i rozgłaszanie jego zmian.

/**
 * Położenie i rozmiar warstwy, czyli cztery liczby opisujące jeden prostokąt
 * w układzie współrzędnych kanwy. Wartości są podawane w jednostkach kompozycji,
 * a nie w pikselach ekranu, więc nie zależą od powiększenia widoku.
 */
export interface Prostokat {
  x: number;
  y: number;
  width: number;
  height: number;
}

export interface StanKompozycji {
  warstwy(): readonly DesignBoardLayer[];
  zaznaczone(): readonly string[];
  widok(): WidokKanwy;
  nazwa(): string;
  /** Identyfikator kompozycji nadany przez rdzeń; pusty zakłada nową. */
  idKompozycji(): string;
  ustawNazwe(nazwa: string): void;
  ustawIdKompozycji(id: string): void;
  /** Wstawia kompozycję odczytaną z rdzenia; zastępuje dotychczasowy układ kanwy. */
  wczytaj(identyfikator: string, nazwa: string, warstwy: readonly DesignBoardLayer[]): void;
  dolozZasob(zasob: DesignAsset): void;
  dolozElement(nazwa: string): void;
  /** Dokłada ramkę obszaru roboczego o szerokości punktu łamania produktu. */
  dolozRamke(nazwa: string, szerokosc: number): void;
  usun(idWarstwy: string): void;
  /** Przełącza zaznaczenie warstwy; `dolacz` zachowuje dotychczasowe. */
  zaznacz(idWarstwy: string, dolacz: boolean): void;
  zaznaczWszystkie(): void;
  odznacz(): void;
  przestawBlokade(idWarstwy: string): void;
  ustawAdnotacje(idWarstwy: string, tresc: string): void;
  przesun(idWarstwy: string, x: number, y: number): void;
  /** Ustawia położenie i rozmiar warstwy naraz — droga inspektora właściwości. */
  ustawGeometrie(idWarstwy: string, prostokat: Prostokat): void;
  wyrownaj(os: OsWyrownania): void;
  /** Rozmieszcza warstwy zaznaczone w równych odstępach. */
  rozmiesc(pionowo: boolean): void;
  ulozSzablon(szablon: SzablonUkladu): void;
  ustawWidok(zmiana: Partial<WidokKanwy>): void;
  obserwuj(sluchacz: () => void): () => void;
}

export function utworzStanKompozycji(): StanKompozycji {
  const zapis: ZapisKompozycji = pustyZapis();
  const sluchacze = new Set<() => void>();

  /** Wykonuje czynność na zapisie i rozgłasza jedną zmianę. */
  function zmien(czynnosc: (zapis: ZapisKompozycji) => void): void {
    czynnosc(zapis);
    for (const sluchacz of [...sluchacze]) sluchacz();
  }

  return {
    warstwy: () => zapis.warstwy,
    zaznaczone: () => zapis.wybor,
    widok: () => zapis.widok,
    nazwa: () => zapis.nazwa,
    idKompozycji: () => zapis.identyfikator,
    ustawNazwe: (nazwa) => zmien((z) => void (z.nazwa = nazwa)),
    ustawIdKompozycji: (id) => zmien((z) => void (z.identyfikator = id)),
    wczytaj: (identyfikator, nazwa, warstwy) =>
      zmien((z) => wczytajUklad(z, identyfikator, nazwa, warstwy)),
    dolozZasob: (zasob) => zmien((z) => dolozZasob(z, zasob)),
    dolozElement: (nazwa) => zmien((z) => dolozElement(z, nazwa)),
    dolozRamke: (nazwa, szerokosc) => zmien((z) => dolozRamke(z, nazwa, szerokosc)),
    usun: (idWarstwy) => zmien((z) => usunWarstwe(z, idWarstwy)),
    zaznacz: (idWarstwy, dolacz) => zmien((z) => przestawZaznaczenie(z, idWarstwy, dolacz)),
    zaznaczWszystkie: () => zmien((z) => void (z.wybor = z.warstwy.map((w) => w.id))),
    odznacz: () => zmien((z) => void (z.wybor = [])),
    ustawAdnotacje: (idWarstwy, tresc) => zmien((z) => zmienWarstwe(z, idWarstwy, { note: tresc })),
    przesun: (idWarstwy, x, y) => zmien((z) => zmienWarstwe(z, idWarstwy, { x, y })),
    ustawGeometrie: (idWarstwy, prostokat) =>
      zmien((z) => zmienWarstwe(z, idWarstwy, { ...prostokat })),
    wyrownaj: (os) => zmien((z) => naniesWarstwy(z, wyrownajWarstwy(objeteWarstwy(z), os))),
    rozmiesc: (pionowo) => zmien((z) => naniesWarstwy(z, rozmiescWarstwy(objeteWarstwy(z), pionowo))),
    ulozSzablon: (szablon) => zmien((z) => void (z.warstwy = ulozWedlugSzablonu(z.warstwy, szablon))),
    ustawWidok: (zmiana) => zmien((z) => void (z.widok = { ...z.widok, ...zmiana })),

    przestawBlokade(idWarstwy) {
      const warstwa = zapis.warstwy.find((w) => w.id === idWarstwy);
      if (warstwa === undefined) return;
      zmien((z) => zmienWarstwe(z, idWarstwy, { locked: warstwa.locked !== true }));
    },

    obserwuj(sluchacz) {
      sluchacze.add(sluchacz);
      return () => sluchacze.delete(sluchacz);
    },
  };
}
