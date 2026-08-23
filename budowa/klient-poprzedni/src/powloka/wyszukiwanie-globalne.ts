import type { NazwaIkony } from '../ikony/ikony';
import { POZYCJE_USTAWIEN } from '../strona-glowna/pozycje-ustawien';
import type { KartaSesji } from './karty-sesji';
import type { Srodowisko } from './srodowiska';
import { wykonajZaczep, type OdbiorcaZdania, type ZaczepyPaska } from './zaczepy-paska';

/**
 * Wyszukiwanie globalne — materiał, nie mechanizm i nie wygląd.
 *
 * Plik zbiera to, co da się w tej chwili przeszukać, i oddaje jako płaskie
 * wpisy. Nie rysuje ani jednego wiersza, nie zna menu i nie zna rdzenia;
 * rysowanie należy do mechanizmu z biblioteki, obsadzonego w
 * `wykaz-wynikow.ts`.
 *
 * Wpisy powstają wyłącznie z tego, co powłoka już ma u siebie: z wykazu
 * środowiska pobranego przez `module.list` do bocznej nawigacji, z pasa kart
 * sesji i z klientowego `POZYCJE_USTAWIEN` — ani jednej nowej komendy, ani
 * jednego nowego odczytu. Drugi odczyt `module.list` obok tego, który zrobiła
 * nawigacja, byłby drugą prawdą o tym samym wykazie, a wykazy rozjeżdżają się,
 * gdy rdzeń zmieni jeden z nich.
 *
 * Czego wyszukiwanie nie obejmuje — nazwane wprost, żeby nikt nie odczytał
 * granicy jako usterki:
 *   — nie przeszukuje treści: ani wiadomości, ani transkryptów, ani plików, ani
 *     dokumentów. Kontrakt (`shared/contract.json`) nie ma rodziny `search.*`;
 *     są wyłącznie dwie komendy modułowe o innym zakresie —
 *     `library.file.search` i `knowledge.search`;
 *   — tylko środowisko otwarte. Moduły spoza macierzy widoczności bieżącego
 *     środowiska i moduły innych środowisk nie wchodzą, bo powłoka ich nie ma;
 *     `module.list` bez `environmentId` oddaje komplet platformy, ale ten
 *     odczyt robi się poza powłoką;
 *   — okna operacyjne wychodzą kodem, nie nazwą polską:
 *     `Module.operationalWindowCodes` niesie same kody, a rdzeń nazw okien nie
 *     oddaje;
 *   — sesje to karty otwarte w tym oknie, a nie `session.list` rdzenia; sesje
 *     zamknięte i sesje innych klientów nie wchodzą;
 *   — bez składni zapytań (`typ:`, `w:`), bez operatorów logicznych i bez
 *     tolerancji literówek. Fraza jest podciągiem — tak działa mechanizm
 *     biblioteki, który ten wykaz przycina;
 *   — bez paginacji. Przy kilkunastu modułach, kilkudziesięciu oknach i kilku
 *     kartach wykaz jest kompletem; przy tysiącach pozycji potrzebna będzie
 *     komenda rdzenia, a nie ten plik.
 */

/** Rodzaj wpisu — on rozstrzyga o grupie, w której wpis stanie w wykazie. */
export type RodzajWpisu = 'srodowisko' | 'modul' | 'okno' | 'sesja' | 'ustawienie';

/** Jedna rzecz, którą Operator może znaleźć i wybrać. */
export interface WpisWyszukiwania {
  /** Klucz niepowtarzalny w całym wykazie; niesie rodzaj, żeby się nie zderzał. */
  klucz: string;
  rodzaj: RodzajWpisu;
  /**
   * Nazwa widoczna w wykazie.
   *
   * Bez drugiej nazwy („nazwaKrotka" mechanizmu). Mechanizm biblioteki umie
   * nieść dwie nazwy — pełną ze źródłem i skróconą bez niego — i rysuje obie
   * w jednym wierszu. Ma to sens w wykazie płaskim na setki pozycji, gdzie bez
   * źródła dwie pozycje o tej samej nazwie są nie do rozróżnienia. Tutaj źródło
   * niesie nagłówek grupy („Moduły", „Sesje otwarte"), więc druga nazwa dałaby
   * wiersz w rodzaju „Studio · studio.editor (studio.editor)".
   */
  nazwa: string;
  /** Zdanie mówiące, co wybór tej pozycji zrobi. Pusty napis = bez opisu. */
  opis: string;
  ikona: NazwaIkony;
  /** Skutek wyboru. Wołany przez wykaz, nigdy przez ten plik. */
  wykonaj(): void;
}

/** Nagłówek grupy wpisów jednego rodzaju. */
export const NAZWY_GRUP: Readonly<Record<RodzajWpisu, string>> = {
  srodowisko: 'Środowisko',
  modul: 'Moduły',
  okno: 'Okna operacyjne',
  sesja: 'Sesje otwarte',
  ustawienie: 'Ustawienia i okna platformy',
};

/** Porządek grup w wykazie przy pustej frazie. */
export const PORZADEK_GRUP: readonly RodzajWpisu[] = [
  'modul',
  'okno',
  'sesja',
  'ustawienie',
  'srodowisko',
];

/**
 * Zakres przeszukania wypowiedziany zdaniem — dla stanu „brak trafień".
 *
 * Stoi tutaj, a nie w wykazie, bo to ten plik wie, co naprawdę weszło do
 * materiału. Zdanie w wykazie rozjechałoby się z zakresem przy pierwszej
 * zmianie źródeł.
 */
export const ZAKRES_PRZESZUKANIA =
  'otwarte środowisko, jego moduły, kody okien operacyjnych tych modułów, '
  + 'karty sesji tego okna oraz ustawienia i okna platformy';

/** Co powłoka podaje wyszukiwaniu; wszystko czytane na żądanie, nic kopiowane. */
export interface ZaleznosciWyszukiwania {
  /** Wykaz środowiska w postaci, w jakiej ma go boczna nawigacja. */
  srodowisko(): Srodowisko;
  /** Karty sesji otwarte w tym oknie. */
  karty(): readonly KartaSesji[];
  /** Przestawia boczną nawigację na moduł o tym kluczu. */
  wybierzModul(kluczPozycji: string): void;
  /** Czyni kartę sesji czynną. */
  wybierzKarte(id: string): void;
  /**
   * Czynności paska; czytane w chwili wyboru pozycji (`zaczepy-paska.ts`).
   *
   * Brak zaczepu znaczy „ten pasek nie dostał drogi", a nie „okna nie ma".
   * Rozróżnienie jest istotne: powierzchnie Always On Display i Mobile istnieją
   * (`aod/kolumna-aod.ts`, `mobile/okno-mobile.ts`) i są osiągalne z listwy
   * strony głównej, więc zdanie „nie jest zbudowane" kłamałoby o bycie, który
   * jest.
   */
  zaczepy?: ZaczepyPaska;
  /**
   * Zgłasza zdanie, którego wykaz nie ma jak wykonać — np. wskazanie okna
   * operacyjnego, do którego powłoka nie ma drogi. Bez odbiorcy zdanie
   * przepada świadomie i jest to jedyne miejsce, w którym tak wolno.
   */
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
      // `globus`, a nie `srodowisko-<kod>`: zestaw ma godła czterech konkretnych
      // środowisk, a wpis powstaje dla dowolnego kodu, jaki poda rdzeń.
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

  /**
   * Okna operacyjne modułów — po kodzie, bo nazwy rdzeń nie oddaje.
   *
   * Okno rozmowy odpada z wykazu: montuje je scena sesji przy każdym module
   * i nie jest oknem operacyjnym. Rdzeń niesie je dziś w dwóch postaciach
   * (`chat-window` oraz `<kod>.chat-window`) — odsiewamy obie, tak samo jak
   * robi to `moduly/katalog-okien.ts`. Warunek jest przepisany, a nie
   * zaimportowany, bo `moduly/` jest katalogiem innej grupy roboczej.
   */
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

  /**
   * Ustawienia i okna platformy — z `POZYCJE_USTAWIEN`, nie z drugiego wykazu.
   *
   * Ten sam wykaz zasila listwę strony głównej, menu aplikacji i menu
   * Operatora. Druga kopia rozjechałaby się przy pierwszej zmianie nazwy.
   */
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
