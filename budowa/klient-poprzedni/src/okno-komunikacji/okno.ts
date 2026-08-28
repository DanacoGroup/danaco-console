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
 * Okno komunikacji łączy nagłówek, powierzchnię modułową, historię wątku i pole wpisywania, nie znając kontraktu ani transportu — wystawia treść operatora i przyjmuje treść przychodzącą, a przełączenie modułu zachowuje rozmowę bez dotykania historii.
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
  /** Warstwa dyktowania stoi na oknie, bo potrzebuje kanału do sprawdzenia silnika i wysłania nagrania. */
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

    // Poprzednia warstwa dyktowania rozłącza się przed wejściem nowej, by nie zostawić zapalonej lampki.
    podlaczDyktowanie(warstwa) {
      dyktowanie?.rozlacz();
      dyktowanie = warstwa;
    },
  };
}
