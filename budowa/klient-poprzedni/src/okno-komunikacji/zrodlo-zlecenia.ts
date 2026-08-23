import type { Window } from '../../../shared/contract';
import type { Odsubskrybuj } from '../polaczenie/magistrala-zdarzen';
import type { Kanal } from '../protokol/kanal';
import type { UstawieniaOkna } from '../sterowanie/klucze-ustawien';
import type { KomunikatZmiany } from '../sterowanie/komunikat-zmiany';
import { utworzStanSterowania } from '../sterowanie/stan-sterowania';
import { utworzZmianeOkna, type ZmianaPolOkna } from '../sterowanie/zmiana-okna';
import { utworzZmianeUstawienia } from '../sterowanie/zmiana-ustawienia';

/**
 * Port paska zlecenia do kompletu sterowania okna.
 *
 * Jedno źródło na cały pasek, nie jedno na ster. Cztery stery paska — model,
 * nakład rozumowania, tryb zatwierdzania i katalog roboczy — czytają jedną
 * migawkę i piszą jedną drogą. Gdyby każdy zakładał własny `StanSterowania`,
 * pasek miałby cztery kopie tej samej prawdy i cztery komplety subskrypcji na
 * okno, a rozjazd między nimi byłby kwestią czasu, nie możliwości.
 *
 * Drogi zapisu są dwie, bo kontrakt ma dwie. Model, tryb uprawnień i katalogi
 * robocze mają pola w treści `window.update`. Nakład rozumowania takiego pola
 * nie ma i idzie ustawieniem poziomu okna (`config.set`, klucz katalogu rdzenia
 * `naklad_rozumowania`). Port oddaje obie drogi osobno, zamiast zlewać je
 * w jedną „zapisz cokolwiek" — zlanie ukryłoby przed wołającym fakt, że to dwie
 * różne komendy o dwóch różnych potwierdzeniach.
 *
 * Te same nastawy stoją w kolumnie sterowania (`sterowanie/`), ale reguła
 * protokołu nie jest tu pisana po raz drugi: odczyt idzie przez
 * `StanSterowania`, zapis przez `ZmianaOkna` i `ZmianaUstawienia` — te same
 * byty, te same komendy, ten sam `windowId`, te same nazwy zmian
 * w komunikatach. Oba widoki przyjmują wyłącznie stan potwierdzony przez
 * rdzeń: odpowiedź na komendę oraz zdarzenia `window.changed` i
 * `config.changed`, które rdzeń rozgłasza do wszystkich połączeń konta
 * (`core/zdarzenia.go` → `Nadajnik.Rozglos`). Dlatego zmiana dokonana w pasku
 * dochodzi do kolumny sterowania tą samą drogą, którą dochodzi do drugiego
 * urządzenia konta — i odwrotnie.
 *
 * To ten sam odczyt równoległy, co w `zrodlo-srodowiska.ts`
 * i `widok-sterowania/obserwator-ustawien.ts`: własny egzemplarz stanu przestaje
 * być potrzebny, gdy `PanelSterowania` wystawi `migawka()` i `naZmiane()`.
 */

/** Migawka, z której stery paska czerpią całą swoją prawdę. */
export interface MigawkaZlecenia {
  /** Okno komunikacji potwierdzone przez rdzeń. */
  okno: Window;
  /** Ustawienia z zasięgu okna. */
  ustawienia: UstawieniaOkna;
}

/** Odczyt i zapis nastaw okna widziane przez stery paska. */
export interface ZrodloZlecenia {
  /** Bieżąca migawka stanu okna. */
  migawka(): MigawkaZlecenia;
  /** Subskrypcja zmian stanu (`window.changed`, `config.changed`). */
  naZmiane(sluchacz: () => void): void;
  /**
   * Zmiana pól okna komendą `window.update`.
   *
   * Odrzucenie niesie zdanie odmowy złożone przez `komunikat-zmiany.ts` — kod
   * i treść wprost z kontraktu, bez parafrazy.
   */
  zastosuj(nazwa: string, zmiana: ZmianaPolOkna): Promise<void>;
  /** Zapis ustawienia poziomu okna komendą `config.set`. */
  zapisz(nazwa: string, klucz: string, wartosc: unknown): Promise<void>;
  /** Odłącza subskrypcje kanału i stanu. */
  rozlacz(): void;
}

/**
 * Kolejka wywołań czekających na los swojej zmiany.
 *
 * `ZmianaOkna` i `ZmianaUstawienia` meldują wynik wywołaniem zwrotnym, a stery
 * żądają obietnicy — kolejka jest całym przekładem między jednym a drugim.
 * Odbiorca dostaje dokładnie jeden komunikat na jedną wysyłkę, więc kolejka nie
 * rośnie: każdy meldunek zdejmuje z niej najstarsze oczekiwanie.
 */
type Kolejka = ((komunikat: KomunikatZmiany) => void)[];

/** Zamienia meldunek wywołania zwrotnego na rozstrzygnięcie obietnicy. */
function rozstrzygnij(kolejka: Kolejka, komunikat: KomunikatZmiany): void {
  kolejka.shift()?.(komunikat);
}

export function utworzZrodloZlecenia(kanal: Kanal, okno: Window): ZrodloZlecenia {
  const stan = utworzStanSterowania(okno);

  const oczekujaceOkna: Kolejka = [];
  const oczekujaceUstawien: Kolejka = [];

  const zmiana = utworzZmianeOkna(kanal, stan, (komunikat) => {
    rozstrzygnij(oczekujaceOkna, komunikat);
  });
  const ustawienia = utworzZmianeUstawienia(kanal, stan, (komunikat) => {
    rozstrzygnij(oczekujaceUstawien, komunikat);
  });
  // Ustawienia poziomu okna wczytane od razu — bez tego ster nakładu pokazywałby
  // „Bez wskazania" przy stopniu zapisanym w bazie.
  ustawienia.wczytaj();

  const odsubskrybuj: Odsubskrybuj[] = [];

  /** Obietnica spełniona potwierdzeniem, odrzucona zdaniem odmowy rdzenia. */
  function obietnica(kolejka: Kolejka, wyslij: () => void): Promise<void> {
    return new Promise<void>((spelnij, odrzuc) => {
      kolejka.push((komunikat) => {
        if (komunikat.udany) spelnij();
        else odrzuc(new Error(komunikat.tresc));
      });
      wyslij();
    });
  }

  return {
    migawka: () => stan.migawka(),

    naZmiane(sluchacz) {
      odsubskrybuj.push(stan.naZmiane(() => sluchacz()));
    },

    zastosuj(nazwa, tresc) {
      return obietnica(oczekujaceOkna, () => zmiana.zastosuj(nazwa, tresc));
    },

    zapisz(nazwa, klucz, wartosc) {
      return obietnica(oczekujaceUstawien, () => ustawienia.zapisz(nazwa, klucz, wartosc));
    },

    rozlacz() {
      for (const zdejmij of odsubskrybuj.splice(0)) zdejmij();
      zmiana.rozlacz();
      ustawienia.rozlacz();
    },
  };
}
