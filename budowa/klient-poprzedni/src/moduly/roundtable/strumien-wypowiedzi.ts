import { ChunkKind } from '../../../../shared/contract';
import type { StanDebaty } from './stan-debaty';
import type { FragmentDebaty } from './zrodlo-strumienia-debaty';

/**
 * Wypowiedzi rosnące na żywo — gromadzenie fragmentów `stream.chunk`
 * i przypisanie ich uczestnikom debaty.
 *
 * `zrodlo-strumienia-debaty.ts` oddaje fragment z surowym `messageId`, bo
 * subskrypcja nie zna składu debaty. Przypisanie fragmentu do uczestnika wymaga
 * wiedzy o składzie i o wypowiedziach tury — czyli `StanDebaty` — więc stoi
 * tutaj, a nie w źródle.
 *
 * Rozpoznanie pyta stan, zamiast patrzeć na przedrostek identyfikatora. Rdzeń
 * nadaje identyfikatory z przedrostkami (`uczest-`, `wypow-`, `tura-` —
 * `adapter_modul_roundtable.go`), ale przedrostków tych nie ma w kontrakcie:
 * ani w `RoundtableParticipant`, ani w `RoundtableStatement`, ani
 * w `StreamChunkEvent`. Klient oparty na nich orzekałby o rdzeniu rzecz, której
 * kontrakt nie obiecuje, i zamilkłby przy ich zmianie bez błędu kompilacji.
 * Dlatego pytamy stan: identyfikator znany składowi jest uczestnikiem,
 * identyfikator znany wykazowi wypowiedzi tury oddaje swojego mówcę,
 * a identyfikator nieznany żadnemu z nich zostaje nieprzypisany i widoczny.
 *
 * W debacie `messageId` niesie kod wypowiedzi, nie kod uczestnika — tak mówi
 * opis `StreamChunkEvent.messageId` w kontrakcie, a mówcę klient bierze
 * z `RoundtableStatement.participantId`. Rozpoznanie po składzie zostaje mimo to
 * jako droga obronna: kosztuje jedno przeszukanie wykazu tury i odpowiada
 * poprawnie także wtedy, gdy fragment przyjdzie pod kodem uczestnika.
 *
 * Fragmenty nieprzypisane czekają, zamiast być porzucane. Rozpoznanie po kodzie
 * wypowiedzi wymaga, żeby stan tę wypowiedź już znał, a wykaz wypowiedzi tury
 * napełnia dopiero zdarzenie `roundtable.debate.changed`; dwa strumienie idą
 * osobnymi biegami, więc pierwszy fragment potrafi wyprzedzić zdarzenie, które
 * nazywa jego wypowiedź. Okno otwarte w trakcie debaty nie zna też jeszcze
 * całego składu: odczyt uczestników jest w kontrakcie (`roundtable.model.list`),
 * ale żadne okno modułu go dziś nie wywołuje. Taki fragment czeka pod własnym
 * identyfikatorem i przechodzi do uczestnika w chwili, gdy stan go pozna.
 *
 * Moduł nie rysuje niczego i nie zna klas CSS. Nie utrwala też wypowiedzi —
 * utrwala je rdzeń, a wykaz zamknięty niesie `StanDebaty`.
 */

/** Głos jednego uczestnika w postaci, w jakiej dotąd przyszedł. */
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

/** Głos, którego nie dało się przypisać uczestnikowi — wraz z surowym kluczem. */
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

/** Pusty głos — jedno miejsce, żeby cztery pola nie rozjechały się przy dopisaniu. */
function pustyGlos(): GlosNaZywo {
  return { tekst: '', rozumowanie: '', domkniety: false, przyczyna: '' };
}

/**
 * Dokłada fragment do głosu.
 *
 * Rozdzielenie rodzajów jest celowe. Tok rozumowania (`thinking`) nie jest
 * wypowiedzią uczestnika: rdzeń liczy do treści wypowiedzi wyłącznie fragmenty
 * rodzaju `text` (`adapter_modul_roundtable_glos.go`), więc doklejenie
 * rozumowania do tekstu dałoby na żywo zdanie inne niż to, które za chwilę
 * utrwali rdzeń. Rodzaje pozostałe — narzędzia, obraz, dźwięk, prowenancja,
 * konto — nie niosą słów wypowiedzi i nie są tu gromadzone; panel debaty ich
 * nie pokazuje.
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

/** Scala głos oczekujący z głosem już przypisanym — oczekujący dokleja się na koniec. */
function scal(przypisany: GlosNaZywo | undefined, oczekujacy: GlosNaZywo): GlosNaZywo {
  const glos = przypisany ?? pustyGlos();
  return {
    tekst: glos.tekst + oczekujacy.tekst,
    rozumowanie: glos.rozumowanie + oczekujacy.rozumowanie,
    domkniety: glos.domkniety || oczekujacy.domkniety,
    przyczyna: oczekujacy.przyczyna === '' ? glos.przyczyna : oczekujacy.przyczyna,
  };
}
