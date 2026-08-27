import { ChangeKind, type LibraryFile, type Module } from '../../../../shared/contract';
import { zbadajTresc, type StanTresci } from './dostepnosc-tresci';
import {
  odczytajTresc,
  usunPlik,
  utworzMagazyn,
  wchlonPlik,
  zapiszTresc,
  zbierz,
  type FazaWykazu,
  type MagazynBiblioteki,
  type TrafienieZnaczenia,
  type TrybWyszukiwania,
  type WidokWykazu,
  type ZawezenieWykazu,
} from './magazyn-biblioteki';
import { odczytajWykaz, szukajWykaz } from './odczyty-biblioteki';
import type { ZrodloBiblioteki } from './zrodlo-biblioteki';
import type { ZrodloZnaczenia } from './zrodlo-znaczenia';

/**
 * Stan biblioteki utrzymuje jeden zbiór plików i jedno zaznaczenie na cały
 * moduł: pięć okien pracują na tym samym pliku czynnym, a forma prezentacji
 * i tryb wyszukiwania mieszkają tu, nie w oknach.
 */
export type {
  FazaWykazu,
  TrafienieZnaczenia,
  TrybWyszukiwania,
  WidokWykazu,
  ZawezenieWykazu,
} from './magazyn-biblioteki';

export interface StanBiblioteki {
  /** Źródło komend obszaru `library.*` — okna wołają je wprost. */
  zrodlo: ZrodloBiblioteki;
  /** Wskaźnik znaczenia biblioteki — rodzina `knowledge.*` w zakresie `library`. */
  wskaznik: ZrodloZnaczenia;
  pliki(): readonly LibraryFile[];
  faza(): FazaWykazu;
  /** Powód ostatniej odmowy; pusty, gdy odczyt się udał. */
  powod(): string;
  /** Katalog zawężający wykaz; pusty znaczy cały zbiór. */
  katalog(): string;
  ustawKatalog(sciezka: string): void;
  /** Forma prezentacji wykazu; nie dotyka rdzenia. */
  widok(): WidokWykazu;
  ustawWidok(widok: WidokWykazu): void;
  /** Tryb wyszukiwania: słowa, znaczenie albo obie drogi naraz. */
  tryb(): TrybWyszukiwania;
  ustawTryb(tryb: TrybWyszukiwania): void;
  /** Zawężenie wykazu do wskazanego zbioru; `null` znaczy cały wykaz. */
  zawezenie(): ZawezenieWykazu | null;
  ustawZawezenie(zawezenie: ZawezenieWykazu | null): void;
  /** Trafienie wskaźnika znaczenia przy pliku; `null`, gdy go nie ma. */
  znaczenie(idPliku: string): TrafienieZnaczenia | null;
  /** Pliki po zawężeniu katalogiem i zbiorem wskazanym — to widzi Explorer. */
  widoczne(): readonly LibraryFile[];
  zaznaczone(): readonly string[];
  przelaczZaznaczenie(idPliku: string): void;
  /** Plik czynny — źródło podglądu, wersji i etykiet. */
  czynny(): LibraryFile | null;
  wskaz(idPliku: string | null): void;
  /** Etykiety i kolekcje złożone ze zbioru plików; kontrakt nie ma ich katalogu. */
  etykiety(): readonly string[];
  kolekcje(): readonly string[];
  /** Dopisuje kolekcję założoną w tej sesji — bez pliku nie wyszłaby ze zbioru. */
  dopiszKolekcje(idKolekcji: string): void;
  /** Żywe okno komunikacji sesji — źródło przenoszenia kontekstu. */
  idOkna(): string;
  ustawOkno(idOkna: string): void;
  /** Katalog modułów z `module.list`; pusty, dopóki rdzeń go nie oddał. */
  moduly(): readonly Module[];
  /** Powód pustego katalogu modułów; pusty napis znaczy „odczyt się udał". */
  powodModulow(): string;
  ustawModuly(moduly: readonly Module[], powod: string): void;
  // Nastawa modułu docelowego jest jedna na cały moduł: ten sam ster czytają Explorer i File Preview.
  modulDocelowy(): string;
  ustawModulDocelowy(kod: string): void;
  // Co rdzeń odpowiedział o treści pliku, osobno od metryki; plik nieodpytany zwraca nieznana, nie brak.
  tresc(idPliku: string): StanTresci;
  /** Odkłada werdykt poznany przy okazji innego wywołania, na przykład podglądu okna. */
  zapiszTresc(idPliku: string, stan: StanTresci): void;
  /** Pyta rdzeń o treść pliku i odkłada jego odpowiedź; zwraca ją wołającemu. */
  zbadajTresc(idPliku: string): Promise<StanTresci>;
  odczytaj(fraza: string, etykieta: string): Promise<void>;
  /** Wyszukuje w trybie bieżącym i oddaje zdanie o drodze, którą wynik powstał. */
  szukaj(fraza: string): Promise<string>;
  wchlon(plik: LibraryFile): void;
  obserwuj(sluchacz: () => void): () => void;
  rozlacz(): void;
}

export function utworzStanBiblioteki(
  zrodlo: ZrodloBiblioteki,
  wskaznik: ZrodloZnaczenia,
): StanBiblioteki {
  const magazyn: MagazynBiblioteki = utworzMagazyn();

  const odsubskrybuj = zrodlo.naZmianePliku((tresc) => {
    if (tresc.change === ChangeKind.Deleted) {
      usunPlik(magazyn, tresc.file.id);
      return;
    }
    wchlonPlik(magazyn, tresc.file);
  });

  return {
    zrodlo,
    wskaznik,
    pliki: () => magazyn.zbior,
    faza: () => magazyn.faza,
    powod: () => magazyn.powod,
    katalog: () => magazyn.sciezka,
    widok: () => magazyn.widok,
    tryb: () => magazyn.tryb,
    zawezenie: () => magazyn.zawezenie,
    znaczenie: (idPliku) => magazyn.znaczenia.get(idPliku) ?? null,
    zaznaczone: () => magazyn.zaznaczenie,
    idOkna: () => magazyn.okno,
    etykiety: () => zbierz(magazyn.zbior, (plik) => plik.tags ?? []),

    kolekcje: () =>
      [
        ...new Set([
          ...zbierz(magazyn.zbior, (plik) => plik.collectionIds ?? []),
          ...magazyn.zalozone,
        ]),
      ].sort((pierwsza, druga) => pierwsza.localeCompare(druga, 'pl')),

    // Dwa zawężenia składają się koniunkcyjnie; zbiór wskazany idzie pierwszy, bo bywa znacznie węższy.
    widoczne() {
      const wskazane =
        magazyn.zawezenie === null
          ? magazyn.zbior
          : magazyn.zbior.filter((plik) => magazyn.zawezenie?.kody.includes(plik.id) === true);
      return magazyn.sciezka === ''
        ? wskazane
        : wskazane.filter((plik) => (plik.path ?? '').startsWith(magazyn.sciezka));
    },

    czynny: () => magazyn.zbior.find((plik) => plik.id === magazyn.wskazany) ?? null,

    ustawKatalog(sciezka) {
      magazyn.sciezka = sciezka;
      magazyn.oglos();
    },

    ustawWidok(widok) {
      if (magazyn.widok === widok) return;
      magazyn.widok = widok;
      magazyn.oglos();
    },

    ustawTryb(tryb) {
      if (magazyn.tryb === tryb) return;
      magazyn.tryb = tryb;
      magazyn.oglos();
    },

    ustawZawezenie(zawezenie) {
      magazyn.zawezenie = zawezenie;
      magazyn.oglos();
    },

    przelaczZaznaczenie(idPliku) {
      magazyn.zaznaczenie = magazyn.zaznaczenie.includes(idPliku)
        ? magazyn.zaznaczenie.filter((wpis) => wpis !== idPliku)
        : [...magazyn.zaznaczenie, idPliku];
      magazyn.oglos();
    },

    wskaz(idPliku) {
      if (magazyn.wskazany === idPliku) return;
      magazyn.wskazany = idPliku;
      magazyn.oglos();
    },

    dopiszKolekcje(idKolekcji) {
      if (idKolekcji === '' || magazyn.zalozone.includes(idKolekcji)) return;
      magazyn.zalozone = [...magazyn.zalozone, idKolekcji];
      magazyn.oglos();
    },

    ustawOkno(idOkna) {
      magazyn.okno = idOkna;
    },

    moduly: () => magazyn.moduly,
    powodModulow: () => magazyn.powodModulow,

    ustawModuly(moduly, powod) {
      magazyn.moduly = [...moduly];
      magazyn.powodModulow = powod;
      magazyn.oglos();
    },

    modulDocelowy: () => magazyn.modulDocelowy,

    ustawModulDocelowy(kod) {
      if (magazyn.modulDocelowy === kod) return;
      magazyn.modulDocelowy = kod;
      magazyn.oglos();
    },

    tresc: (idPliku) => odczytajTresc(magazyn, idPliku),
    zapiszTresc: (idPliku, stanTresci) => zapiszTresc(magazyn, idPliku, stanTresci),

    async zbadajTresc(idPliku) {
      const stanTresci = await zbadajTresc(zrodlo, idPliku);
      zapiszTresc(magazyn, idPliku, stanTresci);
      return stanTresci;
    },

    odczytaj: (fraza, etykieta) => odczytajWykaz(magazyn, zrodlo, fraza, etykieta),
    szukaj: (fraza) => szukajWykaz(magazyn, zrodlo, wskaznik, fraza, magazyn.tryb),
    wchlon: (plik) => wchlonPlik(magazyn, plik),
    obserwuj: (sluchacz) => magazyn.obserwuj(sluchacz),

    rozlacz() {
      odsubskrybuj();
      magazyn.zamknij();
    },
  };
}
