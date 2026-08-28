import { ChunkKind } from '../../../../shared/contract';
import type { StanDebaty } from './stan-debaty';
import type { FragmentDebaty } from './zrodlo-strumienia-debaty';

/** Wypowiedzi rosnące na żywo — gromadzenie fragmentów strumienia i przypisanie ich uczestnikom debaty: głos jednego uczestnika w postaci, w jakiej dotąd przyszedł. */
export interface GlosNaZywo {
  /** Tekst narosły z fragmentów rodzaju `text`. */
  tekst: string;
  /** Tok rozumowania z fragmentów rodzaju `thinking`; osobno, bo to nie wypowiedź. */
  rozumowanie: string;
  /** Czy strumień tego głosu został domknięty (`done` koperty). */
  domkniety: boolean;
  /** Treść fragmentu rodzaju `error`; pusta, gdy błędu nie było. */
  przyczyna: string;
}

/** Głos, którego nie dało się przypisać uczestnikowi — wraz z surowym kluczem, identyfikatorem nieznanym ani składowi, ani wypowiedziom. */
export interface GlosNieprzypisany {
  identyfikator: string;
  glos: GlosNaZywo;
}

export interface StrumienWypowiedzi {
  /** Głos uczestnika; `null`, gdy dla tej tożsamości nic nie płynęło. */
  glos(idUczestnika: string): GlosNaZywo | null;
  /** Głosy o identyfikatorze nieznanym ani składowi, ani wykazowi wypowiedzi. */
  nieprzypisane(): GlosNieprzypisany[];
  /** Przyjmuje fragment ze źródła. */
  przyjmij(fragment: FragmentDebaty): void;
  /** Porzuca całe gromadzenie — wołane przy zmianie tury i przy zmianie okna. */
  wyczysc(): void;
}

export function utworzStrumienWypowiedzi(stan: StanDebaty): StrumienWypowiedzi {
  const glosy = new Map<string, GlosNaZywo>();
  const oczekujace = new Map<string, GlosNaZywo>();

  /** Przenosi oczekujące głosy, które stan zdążył już rozpoznać. */
  function przeniesRozpoznane(): void {
    for (const [identyfikator, glos] of [...oczekujace]) {
      const uczestnik = rozpoznaj(stan, identyfikator);
      if (uczestnik === '') continue;
      oczekujace.delete(identyfikator);
      glosy.set(uczestnik, scal(glosy.get(uczestnik), glos));
    }
  }

  return {
    glos: (idUczestnika) => glosy.get(idUczestnika) ?? null,

    nieprzypisane: () =>
      [...oczekujace].map(([identyfikator, glos]) => ({ identyfikator, glos })),

    przyjmij(fragment) {
      przeniesRozpoznane();
      const uczestnik = rozpoznaj(stan, fragment.identyfikator);
      const wykaz = uczestnik === '' ? oczekujace : glosy;
      const klucz = uczestnik === '' ? fragment.identyfikator : uczestnik;
      wykaz.set(klucz, dopisz(wykaz.get(klucz), fragment));
    },

    wyczysc() {
      glosy.clear();
      oczekujace.clear();
    },
  };
}

/**
 * Uczestnik stojący za surowym identyfikatorem fragmentu; pusty napis, gdy stan
 * go nie zna. Kolejność sprawdzeń jest kolejnością pewności: tożsamość
 * uczestnika rozstrzyga wprost, wypowiedź — przez swojego mówcę.
 */
function rozpoznaj(stan: StanDebaty, identyfikator: string): string {
  if (identyfikator === '') return '';
  if (stan.uczestnik(identyfikator) !== null) return identyfikator;
  const wypowiedz = stan.wypowiedzi().find((pozycja) => pozycja.id === identyfikator);
  return wypowiedz === undefined ? '' : wypowiedz.participantId;
}

/** Pusty głos — jedno miejsce, żeby cztery pola nie rozjechały się przy dopisaniu fragmentu strumienia. */
function pustyGlos(): GlosNaZywo {
  return { tekst: '', rozumowanie: '', domkniety: false, przyczyna: '' };
}

/**
 * Dokłada fragment do głosu, rozdzielając rodzaje: tekst, rozumowanie i błąd nie mieszają się
 * w jednym polu.
 */
function dopisz(poprzedni: GlosNaZywo | undefined, fragment: FragmentDebaty): GlosNaZywo {
  const glos = poprzedni ?? pustyGlos();
  const nastepny: GlosNaZywo = {
    tekst: fragment.rodzaj === ChunkKind.Text ? glos.tekst + fragment.tekst : glos.tekst,
    rozumowanie:
      fragment.rodzaj === ChunkKind.Thinking ? glos.rozumowanie + fragment.tekst : glos.rozumowanie,
    domkniety: glos.domkniety || fragment.ostatni,
    przyczyna: fragment.rodzaj === ChunkKind.Error ? fragment.tekst : glos.przyczyna,
  };
  return nastepny;
}

/** Scala głos oczekujący z głosem już przypisanym — oczekujący dokleja się na koniec treści i rozumowania. */
function scal(przypisany: GlosNaZywo | undefined, oczekujacy: GlosNaZywo): GlosNaZywo {
  const glos = przypisany ?? pustyGlos();
  return {
    tekst: glos.tekst + oczekujacy.tekst,
    rozumowanie: glos.rozumowanie + oczekujacy.rozumowanie,
    domkniety: glos.domkniety || oczekujacy.domkniety,
    przyczyna: oczekujacy.przyczyna === '' ? glos.przyczyna : oczekujacy.przyczyna,
  };
}
