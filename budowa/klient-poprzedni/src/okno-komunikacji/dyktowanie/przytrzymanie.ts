/**
 * Gest przytrzymania łączy trzymanie przycisku — palcem, myszą albo klawiszem — z trzema wywołaniami: rozpoczęciem, zakończeniem i porzuceniem nagrania.
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

/** Klawisze uznawane za naciśnięcie przycisku — umowa platformy interfejsu, a nie dowolny wybór implementacji tego gestu. */
const KLAWISZE_TRZYMANIA: readonly string[] = [' ', 'Spacebar', 'Enter'];

export function utworzPrzytrzymanie(opcje: OpcjePrzytrzymania): Przytrzymanie {
  return {
    podepnij(cel) {
      // Znacznik trzymania zapamiętuje źródło, by puszczenie klawisza nie kończyło gestu zaczętego myszą.
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
        // Reaguje tylko przycisk główny — prawy przycisk otwiera menu i nie zostawia zdarzenia `pointerup`.
        if (zdarzenie.button !== 0) return;
        idWskaznika = zdarzenie.pointerId;
        // Przechwycenie wskaźnika zapobiega przejęciu `pointerup` przez inny element przy zjechaniu palcem.
        try {
          cel.setPointerCapture(zdarzenie.pointerId);
        } catch {
          // Odmowa przechwycenia nie przerywa gestu — traci tylko odporność na zjechanie poza przycisk.
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
        // Zdarzenie `pointercancel` kończy nagranie, nie porzuca go — słowa odrzuca wyłącznie klawisz Escape.
        domknij(false);
      }

      /** Escape nasłuchiwany jest na dokumencie, bo ognisko podczas trzymania myszą bywa gdzie indziej. */
      function naEscape(zdarzenie: Event): void {
        if ((zdarzenie as KeyboardEvent).key !== 'Escape' || zrodlo === null) return;
        zdarzenie.preventDefault();
        domknij(true);
      }

      function naKlawiszWdol(zdarzenie: KeyboardEvent): void {
        if (!KLAWISZE_TRZYMANIA.includes(zdarzenie.key)) return;
        // Auto-powtarzanie klawisza jest odcinane, by seria zdarzeń nie zaczynała nagrania wielokrotnie.
        if (zdarzenie.repeat) return;
        // Spacja i Enter są przechwytywane, bo ich domyślne działanie koliduje z trzymaniem przycisku.
        zdarzenie.preventDefault();
        zacznij('klawiatura');
      }

      function naKlawiszWgore(zdarzenie: KeyboardEvent): void {
        if (zrodlo !== 'klawiatura' || !KLAWISZE_TRZYMANIA.includes(zdarzenie.key)) return;
        zdarzenie.preventDefault();
        domknij(false);
      }

      function naUtrateOgniska(): void {
        // Ucieczka ogniska przy trzymaniu klawiszem domyka gest jak puszczenie, nie zostawia mikrofonu.
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
        // Odpięcie w trakcie trzymania kończy nagrywanie, bo puszczenie przycisku nie zostanie już zgłoszone.
        if (zrodlo !== null) domknij(true);
      };
    },
  };
}
