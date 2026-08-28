import type { NazwaIkony } from '../ikony/ikony';
import { POZYCJE_USTAWIEN } from '../strona-glowna/pozycje-ustawien';
import type { KartaSesji } from './karty-sesji';
import type { Srodowisko } from './srodowiska';
import { wykonajZaczep, type OdbiorcaZdania, type ZaczepyPaska } from './zaczepy-paska';

// Wyszukiwanie globalne — materiał, nie mechanizm i nie wygląd; oddaje płaskie wpisy do przeszukania.

/** Rodzaj wpisu — on rozstrzyga o grupie, w której dany wpis stanie w całym wykazie tego wyszukiwania globalnego. */
export type RodzajWpisu = 'srodowisko' | 'modul' | 'okno' | 'sesja' | 'ustawienie';

/** Jedna rzecz, którą Operator może znaleźć w wykazie wyszukiwania globalnego i od razu ją tutaj wybrać. */
export interface WpisWyszukiwania {
  /** Klucz niepowtarzalny w całym wykazie; niesie rodzaj, żeby się nie zderzał. */
  klucz: string;
  rodzaj: RodzajWpisu;
  // Nazwa widoczna w wykazie, bez drugiej nazwy skróconej — źródło niesie już nagłówek grupy.
  nazwa: string;
  /** Zdanie mówiące, co wybór tej pozycji zrobi. Pusty napis = bez opisu. */
  opis: string;
  ikona: NazwaIkony;
  /** Skutek wyboru. Wołany przez wykaz, nigdy przez ten plik. */
  wykonaj(): void;
}

/** Nagłówek grupy wpisów jednego rodzaju, wypisywany nad wierszami tego samego rodzaju wpisów tego wykazu. */
export const NAZWY_GRUP: Readonly<Record<RodzajWpisu, string>> = {
  srodowisko: 'Środowisko',
  modul: 'Moduły',
  okno: 'Okna operacyjne',
  sesja: 'Sesje otwarte',
  ustawienie: 'Ustawienia i okna platformy',
};

/** Porządek grup w wykazie wyszukiwania przy pustej frazie, od modułów aż po samo to otwarte środowisko. */
export const PORZADEK_GRUP: readonly RodzajWpisu[] = [
  'modul',
  'okno',
  'sesja',
  'ustawienie',
  'srodowisko',
];

/** Zakres przeszukania wypowiedziany zdaniem — dla stanu braku trafień w tym samym wykazie wyszukiwania. */
export const ZAKRES_PRZESZUKANIA =
  'otwarte środowisko, jego moduły, kody okien operacyjnych tych modułów, '
  + 'karty sesji tego okna oraz ustawienia i okna platformy';

/** Co powłoka podaje wyszukiwaniu globalnemu; wszystko czytane na żądanie, nic tutaj nigdy nie kopiowane. */
export interface ZaleznosciWyszukiwania {
  /** Wykaz środowiska w postaci, w jakiej ma go boczna nawigacja. */
  srodowisko(): Srodowisko;
  /** Karty sesji otwarte w tym oknie. */
  karty(): readonly KartaSesji[];
  /** Przestawia boczną nawigację na moduł o tym kluczu. */
  wybierzModul(kluczPozycji: string): void;
  /** Czyni kartę sesji czynną. */
  wybierzKarte(id: string): void;
  // Czynności paska; czytane w chwili wyboru pozycji, brak zaczepu znaczy, że pasek nie dostał drogi.
  zaczepy?: ZaczepyPaska;
  // Zgłasza zdanie, którego wykaz nie ma jak wykonać; bez odbiorcy zdanie przepada świadomie.
  naKomunikat?: OdbiorcaZdania;
}

export interface ZrodloWyszukiwania {
  /** Komplet wpisów znanych w tej chwili, w porządku grup. */
  wpisy(): readonly WpisWyszukiwania[];
  /** Stan wykazu modułów — pochodzi wprost z bocznej nawigacji. */
  stan(): Srodowisko['stan'];
  /** Treść odmowy rdzenia; pusty napis znaczy „nic nie odmówiło". */
  odmowa(): string;
}

export function utworzZrodloWyszukiwania(
  zaleznosci: ZaleznosciWyszukiwania,
): ZrodloWyszukiwania {
  /** Środowisko otwarte — jedna pozycja, bo powłoka zna tylko to jedno. */
  function wpisSrodowiska(srodowisko: Srodowisko): WpisWyszukiwania {
    return {
      klucz: `srodowisko:${srodowisko.klucz}`,
      rodzaj: 'srodowisko',
      nazwa: srodowisko.nazwa,
      opis:
        srodowisko.motto === ''
          ? 'Środowisko otwarte w tym oknie.'
          : `${srodowisko.motto} — środowisko otwarte w tym oknie.`,
      // Nazwa ikony globus, a nie nazwa środowiska: zestaw ma godła czterech konkretnych środowisk.
      ikona: 'globus',
      wykonaj: () => {
        const pierwsza = srodowisko.pozycje[0];
        if (pierwsza !== undefined) zaleznosci.wybierzModul(pierwsza.klucz);
      },
    };
  }

  function wpisyModulow(srodowisko: Srodowisko): WpisWyszukiwania[] {
    return srodowisko.pozycje.map((pozycja) => ({
      klucz: `modul:${pozycja.klucz}`,
      rodzaj: 'modul' as const,
      nazwa: pozycja.nazwa,
      opis:
        pozycja.opis === ''
          ? 'Moduł środowiska — wybór przestawia boczną nawigację i obszar roboczy.'
          : pozycja.opis,
      ikona: pozycja.ikona,
      wykonaj: () => zaleznosci.wybierzModul(pozycja.klucz),
    }));
  }

  // Okna operacyjne modułów — po kodzie, bo nazwy rdzeń nie oddaje; okno rozmowy odpada z wykazu.
  function wpisyOkien(srodowisko: Srodowisko): WpisWyszukiwania[] {
    const wynik: WpisWyszukiwania[] = [];
    for (const pozycja of srodowisko.pozycje) {
      for (const kod of pozycja.okna) {
        if (kod === 'chat-window' || kod.endsWith('.chat-window')) continue;
        wynik.push({
          klucz: `okno:${pozycja.klucz}:${kod}`,
          rodzaj: 'okno',
          nazwa: `${pozycja.nazwa} · ${kod}`,
          opis:
            `Okno operacyjne modułu ${pozycja.nazwa}, kod katalogu rdzenia. `
            + 'Wybór otwiera moduł — pojedynczego okna powłoka nie przywołuje, '
            + 'robi to sam moduł.',
          ikona: 'karta-okna',
          wykonaj: () => {
            zaleznosci.wybierzModul(pozycja.klucz);
            zaleznosci.naKomunikat?.(
              `Okno ${kod}`,
              `Otwarto moduł ${pozycja.nazwa}. Samego okna „${kod}" powłoka nie `
              + 'przywołuje: kody okien operacyjnych niesie katalog rdzenia, a wejście '
              + 'do pojedynczego okna należy do widoku modułu.',
            );
          },
        });
      }
    }
    return wynik;
  }

  function wpisySesji(): WpisWyszukiwania[] {
    return zaleznosci.karty().map((karta) => ({
      klucz: `sesja:${karta.id}`,
      rodzaj: 'sesja' as const,
      nazwa: karta.tytul(),
      opis: 'Karta sesji otwarta w tym oknie — wybór czyni ją czynną.',
      ikona: 'rozmowa',
      wykonaj: () => zaleznosci.wybierzKarte(karta.id),
    }));
  }

  // Ustawienia i okna platformy z jednego wspólnego wykazu, nie z osobnej, drugiej jego kopii.
  function wpisyUstawien(): WpisWyszukiwania[] {
    return POZYCJE_USTAWIEN.map((pozycja) => ({
      klucz: `ustawienie:${pozycja.kod}`,
      rodzaj: 'ustawienie' as const,
      nazwa: pozycja.nazwa,
      opis: pozycja.wyjasnienie,
      ikona: pozycja.ikona,
      // Zaczep czytany dopiero w chwili wyboru (`zaczepy-paska.ts`).
      wykonaj: () => wykonajZaczep(zaleznosci.zaczepy ?? {}, pozycja, zaleznosci.naKomunikat),
    }));
  }

  return {
    wpisy() {
      const srodowisko = zaleznosci.srodowisko();
      const wedlugRodzaju: Readonly<Record<RodzajWpisu, WpisWyszukiwania[]>> = {
        modul: wpisyModulow(srodowisko),
        okno: wpisyOkien(srodowisko),
        sesja: wpisySesji(),
        ustawienie: wpisyUstawien(),
        srodowisko: [wpisSrodowiska(srodowisko)],
      };
      return PORZADEK_GRUP.flatMap((rodzaj) => wedlugRodzaju[rodzaj]);
    },

    stan: () => zaleznosci.srodowisko().stan,
    odmowa: () => zaleznosci.srodowisko().blad ?? '',
  };
}
