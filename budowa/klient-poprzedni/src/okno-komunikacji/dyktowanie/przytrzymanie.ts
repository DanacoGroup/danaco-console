/**
 * Gest „przytrzymaj, aby nagrać”.
 *
 * Jedna odpowiedzialność: zamienić trzymanie przycisku — palcem, myszą albo
 * klawiszem — na trzy wywołania: start, koniec, porzucenie. Plik nie nagrywa,
 * nie tworzy elementu i nie zna mikrofonu; dostaje cudzy element i oddaje
 * odpięcie.
 */

export interface OpcjePrzytrzymania {
  /** Trzymanie zaczęte — czas ruszyć nagrywanie. */
  naStart(): void;
  /** Trzymanie puszczone — wypowiedź domknięta i idzie dalej. */
  naKoniec(): void;
  /** Trzymanie zerwane rozmyślnie — nagranie do kosza, nie do transkrypcji. */
  naPorzucenie(): void;
}

export interface Przytrzymanie {
  /** Podpina gest do elementu; zwrócona funkcja odpina wszystkie nasłuchy. */
  podepnij(cel: HTMLElement): () => void;
}

/** Klawisze uznawane za naciśnięcie przycisku — umowa platformy, nie wybór. */
const KLAWISZE_TRZYMANIA: readonly string[] = [' ', 'Spacebar', 'Enter'];

export function utworzPrzytrzymanie(opcje: OpcjePrzytrzymania): Przytrzymanie {
  return {
    podepnij(cel) {
      // Jeden znacznik trzymania na podpięcie. Źródło (palec albo klawisz)
      // zapamiętujemy, żeby puszczenie klawisza nie kończyło gestu zaczętego
      // myszą — Operator potrafi trzymać przycisk i pisać drugą ręką.
      let zrodlo: 'wskaznik' | 'klawiatura' | null = null;
      let idWskaznika: number | null = null;

      function zacznij(nowe: 'wskaznik' | 'klawiatura'): void {
        if (zrodlo !== null) return;
        zrodlo = nowe;
        opcje.naStart();
      }

      function domknij(porzuc: boolean): void {
        if (zrodlo === null) return;
        zrodlo = null;
        idWskaznika = null;
        if (porzuc) opcje.naPorzucenie();
        else opcje.naKoniec();
      }

      function naNacisniecie(zdarzenie: PointerEvent): void {
        // Tylko przycisk główny. Prawy przycisk otwiera menu kontekstowe i nie
        // zostawia po sobie `pointerup` na elemencie — gest zawisłby na zawsze.
        if (zdarzenie.button !== 0) return;
        idWskaznika = zdarzenie.pointerId;
        // Przechwycenie wskaźnika: bez niego zjechanie palcem poza krawędź
        // przycisku zabiera `pointerup` innemu elementowi, a nagrywanie zostaje
        // włączone mimo puszczonego palca.
        try {
          cel.setPointerCapture(zdarzenie.pointerId);
        } catch {
          // Przechwycenie bywa odmówione (wskaźnik już zwolniony). Gest działa
          // dalej — traci tylko odporność na zjechanie poza przycisk.
        }
        zacznij('wskaznik');
      }

      function naPuszczenie(zdarzenie: PointerEvent): void {
        if (zrodlo !== 'wskaznik' || zdarzenie.pointerId !== idWskaznika) return;
        const trzymany = idWskaznika;
        domknij(false);
        if (trzymany !== null && cel.hasPointerCapture?.(trzymany)) {
          cel.releasePointerCapture(trzymany);
        }
      }

      function naPrzerwanieWskaznika(zdarzenie: PointerEvent): void {
        if (zrodlo !== 'wskaznik' || zdarzenie.pointerId !== idWskaznika) return;
        // `pointercancel` to systemowe wyrwanie gestu (przewijanie, telefon).
        // Kończymy, a nie porzucamy: Operator zdążył coś powiedzieć i słowa
        // wyrzucamy wyłącznie wtedy, gdy sam o to prosi (Escape).
        domknij(false);
      }

      /**
       * Escape stoi na dokumencie, nie na elemencie.
       *
       * Podczas trzymania myszą ognisko bywa gdzie indziej — przycisk wcale nie
       * musi je mieć. Nasłuch przy elemencie przepuściłby wtedy Escape i nie
       * byłoby czym cofnąć nagrania w połowie zdania.
       */
      function naEscape(zdarzenie: Event): void {
        if ((zdarzenie as KeyboardEvent).key !== 'Escape' || zrodlo === null) return;
        zdarzenie.preventDefault();
        domknij(true);
      }

      function naKlawiszWdol(zdarzenie: KeyboardEvent): void {
        if (!KLAWISZE_TRZYMANIA.includes(zdarzenie.key)) return;
        // Auto-powtarzanie: przytrzymany klawisz sypie zdarzeniami kilkanaście
        // razy na sekundę, a bez tego odcięcia każde z nich próbowałoby zacząć
        // nagranie od nowa.
        if (zdarzenie.repeat) return;
        // Spacja przewija stronę, a Enter po puszczeniu wysyła `click` — oba
        // zachowania kolidują z trzymaniem.
        zdarzenie.preventDefault();
        zacznij('klawiatura');
      }

      function naKlawiszWgore(zdarzenie: KeyboardEvent): void {
        if (zrodlo !== 'klawiatura' || !KLAWISZE_TRZYMANIA.includes(zdarzenie.key)) return;
        zdarzenie.preventDefault();
        domknij(false);
      }

      function naUtrateOgniska(): void {
        // Ognisko uciekło w trakcie trzymania klawiszem (przełączenie okna) —
        // `keyup` już nie przyjdzie. Domykamy jak puszczenie, bo cisza w tym
        // miejscu zostawiłaby zapalony mikrofon bez żadnego wyjścia.
        if (zrodlo === 'klawiatura') domknij(false);
      }

      cel.addEventListener('pointerdown', naNacisniecie);
      cel.addEventListener('pointerup', naPuszczenie);
      cel.addEventListener('pointercancel', naPrzerwanieWskaznika);
      cel.addEventListener('keydown', naKlawiszWdol);
      cel.addEventListener('keyup', naKlawiszWgore);
      cel.addEventListener('blur', naUtrateOgniska);
      cel.ownerDocument.addEventListener('keydown', naEscape);

      return () => {
        cel.removeEventListener('pointerdown', naNacisniecie);
        cel.removeEventListener('pointerup', naPuszczenie);
        cel.removeEventListener('pointercancel', naPrzerwanieWskaznika);
        cel.removeEventListener('keydown', naKlawiszWdol);
        cel.removeEventListener('keyup', naKlawiszWgore);
        cel.removeEventListener('blur', naUtrateOgniska);
        cel.ownerDocument.removeEventListener('keydown', naEscape);
        // Odpięcie w trakcie trzymania nie może zostawić nagrywania włączonego
        // — nikt już nie doniesie o puszczeniu przycisku.
        if (zrodlo !== null) domknij(true);
      };
    },
  };
}
