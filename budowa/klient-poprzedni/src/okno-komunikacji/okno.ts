import { utworzMagistrale, type Odsubskrybuj } from '../polaczenie/magistrala-zdarzen';
import type { StanPolaczenia } from '../polaczenie/stan-polaczenia';
import type { Dyktowanie } from './dyktowanie/dyktowanie';
import { utworzHistorie } from './historia';
import { utworzNaglowek } from './naglowek';
import type { OpisOkna } from './opis-okna';
import { utworzPoleWpisywania } from './pole-wpisywania';
import { utworzPowierzchnieModulu } from './powierzchnia-modulu';
import type { Wpis } from './wpis';

/**
 * Okno komunikacji — złożenie nagłówka, powierzchni modułowej, historii i pola
 * wpisywania.
 *
 * Okno nie zna kontraktu ani transportu: wystawia treść wpisaną przez
 * operatora i przyjmuje treść przychodzącą. Powiązaniem z rdzeniem zajmuje się
 * przepływ komunikatów.
 *
 * Moduł jest właściwością okna, nie wdrożenia: `ustawModul`
 * przestawia wskaźnik modułu, pasek narzędzi promptu i panel kontekstu,
 * a historii wątku nie dotyka — okno zachowuje rozmowę przy zmianie modułu.
 *
 * Panelu akcji w tym złożeniu nie ma: pozycje panelu pochodzą z rejestru rdzenia
 * (`action.list`), a złożenie nie ma kanału — panel bez katalogu byłby atrapą.
 * Okno z kanałem składa warstwa rozmowy (`rozmowa/montaz-rozmowy`) i tam panel
 * akcji jest pełny.
 */
export interface OknoKomunikacji {
  /** Element montowany w dokumencie. */
  element: HTMLElement;
  /** Dokłada wpis do historii. */
  dopisz(wpis: Wpis): void;
  /** Dokłada fragment strumienia do wypowiedzi tej samej persony. */
  dopiszFragment(persona: string, fragment: string): void;
  /** Odświeża wskaźnik stanu połączenia. */
  pokazStan(stan: StanPolaczenia, oczekujace: number): void;
  /** Subskrypcja treści wpisanych przez operatora. */
  naWpisanie(sluchacz: (tresc: string) => void): Odsubskrybuj;
  /** Ustawia ognisko na polu wpisywania. */
  ustawOgnisko(): void;
  /** Przestawia okno na wskazany moduł; historia wątku zostaje. */
  ustawModul(kod: string): void;
  /** Moduł, w którym okno pracuje w tej chwili. */
  modul(): string;
  /**
   * Warstwa dyktowania okna; pusta do chwili uzgodnienia kanału.
   *
   * Warstwa stoi na oknie, a nie w pasku, bo dyktowanie potrzebuje kanału, żeby
   * zapytać rdzeń o dostępność silnika i wysłać nagranie — a pasek polecenia
   * kanału nie widzi. Okno zna oba końce, więc to ono przechowuje warstwę
   * i podaje ją temu, kto rysuje mikrofon.
   */
  dyktowanie(): Dyktowanie | null;
  /** Podłącza warstwę dyktowania; robi to przepływ, gdy kanał jest już znany. */
  podlaczDyktowanie(warstwa: Dyktowanie | null): void;
}

export function utworzOkno(opis: OpisOkna): OknoKomunikacji {
  const wpisane = utworzMagistrale<string>();
  let dyktowanie: Dyktowanie | null = null;
  const naglowek = utworzNaglowek(opis);
  const historia = utworzHistorie();
  const pole = utworzPoleWpisywania((tresc) => wpisane.oglos(tresc));

  const powierzchnia = utworzPowierzchnieModulu({
    naPolecenie: (tekst) => pole.wstaw(tekst),
    modul: opis.modul,
    kontekst: {
      kanalModelu: opis.kanalModelu,
      srodowiskoWykonania: opis.srodowiskoWykonania,
      katalogiRobocze: opis.katalogiRobocze,
      wpisyWatku: () => historia.liczbaWpisow(),
    },
  });

  const element = document.createElement('main');
  element.className = 'dc-okno';
  element.append(naglowek.element, powierzchnia.element, historia.element, pole.element);

  return {
    element,
    dopisz(wpis) {
      historia.dopisz(wpis);
      powierzchnia.odswiezKontekst();
    },
    dopiszFragment: historia.dopiszDoOstatniego,
    pokazStan: naglowek.pokazStan,
    naWpisanie: (sluchacz) => wpisane.subskrybuj(sluchacz),
    ustawOgnisko: pole.ustawOgnisko,
    ustawModul: powierzchnia.ustawModul,
    modul: powierzchnia.modul,

    dyktowanie: () => dyktowanie,

    // Poprzednia warstwa jest rozłączana, zanim wejdzie nowa. Bez tego zmiana
    // kanału zostawiałaby żywy nasłuch sprzętu i — gdyby akurat nagrywał —
    // zapaloną lampkę mikrofonu bez właściciela.
    podlaczDyktowanie(warstwa) {
      dyktowanie?.rozlacz();
      dyktowanie = warstwa;
    },
  };
}
