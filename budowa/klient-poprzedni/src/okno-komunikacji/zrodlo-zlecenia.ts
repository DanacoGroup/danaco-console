import type { Window } from '../../../shared/contract';
import type { Odsubskrybuj } from '../polaczenie/magistrala-zdarzen';
import type { Kanal } from '../protokol/kanal';
import type { UstawieniaOkna } from '../sterowanie/klucze-ustawien';
import type { KomunikatZmiany } from '../sterowanie/komunikat-zmiany';
import { utworzStanSterowania } from '../sterowanie/stan-sterowania';
import { utworzZmianeOkna, type ZmianaPolOkna } from '../sterowanie/zmiana-okna';
import { utworzZmianeUstawienia } from '../sterowanie/zmiana-ustawienia';

// Port paska zlecenia do kompletu sterowania okna. Jedno źródło na cały pasek, nie jedno na ster.

/** Migawka, z której stery paska zlecenia czerpią całą swoją prawdę o stanie okna i jego ustawień poziomu. */
export interface MigawkaZlecenia {
  /** Okno komunikacji potwierdzone przez rdzeń. */
  okno: Window;
  /** Ustawienia z zasięgu okna. */
  ustawienia: UstawieniaOkna;
}

/** Odczyt i zapis nastaw okna widziane przez stery paska, obejmujące migawkę, subskrypcję i dwie drogi zapisu. */
export interface ZrodloZlecenia {
  /** Bieżąca migawka stanu okna. */
  migawka(): MigawkaZlecenia;
  /** Subskrypcja zmian stanu (`window.changed`, `config.changed`). */
  naZmiane(sluchacz: () => void): void;
  /** Zmiana pól okna komendą window.update; odrzucenie niesie zdanie odmowy wprost z kontraktu. */
  zastosuj(nazwa: string, zmiana: ZmianaPolOkna): Promise<void>;
  /** Zapis ustawienia poziomu okna komendą `config.set`. */
  zapisz(nazwa: string, klucz: string, wartosc: unknown): Promise<void>;
  /** Odłącza subskrypcje kanału i stanu. */
  rozlacz(): void;
}

/**
 * Kolejka wywołań czekających na los zmiany jest przekładem między wywołaniem zwrotnym rdzenia a obietnicą, jakiej żądają stery.
 */
type Kolejka = ((komunikat: KomunikatZmiany) => void)[];

/** Zamienia meldunek wywołania zwrotnego rdzenia na rozstrzygnięcie obietnicy oczekiwanej przez ster paska. */
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
  // Ustawienia okna wczytane od razu, inaczej ster nakładu pokazywałby brak wskazania przy stopniu.
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
