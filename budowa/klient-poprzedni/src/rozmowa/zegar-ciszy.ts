import type { StanRozmowy } from './stan-rozmowy';

/**
 * Zegar ciszy — jedyna warstwa okna rozmowy, która mierzy czas.
 *
 * Domknięcie tury przychodzi wyłącznie zdarzeniem z rdzenia, a zerwany
 * strumień nie ma w kontrakcie własnego zdarzenia. Rdzeń ubity w połowie tury
 * zostawiłby więc wpis na stanie „wysyłanie" bez końca, dlatego warstwa
 * rozmowy rozpoznaje zerwanie po jedynym śladzie, jaki zostaje: po tym, że nic
 * nie przychodzi.
 *
 * Zegar niczego nie przerywa i nie zamyka tury — tura wraca do życia, gdy
 * nadejdzie fragment. Zegar wyłącznie nazywa to, co widać; różnicę między
 * długim rozumowaniem a martwym połączeniem rozstrzyga Operator, mając przed
 * sobą liczbę sekund.
 */
export interface ZegarCiszy {
  /** Uruchamia odliczanie; wywołanie powtórne nie zakłada drugiego interwału. */
  uruchom: () => void;
  /** Zatrzymuje odliczanie; wywołanie na zatrzymanym nic nie robi. */
  zatrzymaj: () => void;
}

/** Zależności zegara — stan tury i droga jego zmiany. */
export interface OtoczenieZegara {
  /** Stan bieżący rozmowy; zegar go czyta, nie trzyma własnej kopii. */
  stan: () => StanRozmowy;
  /** Chwila ostatniego znaku życia tury, w milisekundach epoki. */
  znakZycia: () => number;
  /** Przestawia stan — jedyną drogą, którą zegar zmienia cokolwiek. */
  ustawStan: (zmiana: Partial<StanRozmowy>) => void;
}

/** Takt odliczania. Sekunda: krócej byłoby migotaniem, dłużej — nieuwagą. */
const taktCiszy = 1000;

/** Zakłada zegar ciszy nad podanym otoczeniem. */
export function utworzZegarCiszy(o: OtoczenieZegara): ZegarCiszy {
  let uchwyt: ReturnType<typeof setInterval> | null = null;

  function zatrzymaj(): void {
    if (uchwyt === null) return;
    clearInterval(uchwyt);
    uchwyt = null;
  }

  function uruchom(): void {
    if (uchwyt !== null) return;
    uchwyt = setInterval(() => {
      const stan = o.stan();
      if (!stan.wysyla) {
        zatrzymaj();
        return;
      }
      const sekundy = Math.floor((Date.now() - o.znakZycia()) / taktCiszy);
      if (sekundy !== stan.ciszaSekundy) o.ustawStan({ ciszaSekundy: sekundy });
    }, taktCiszy);
  }

  return { uruchom, zatrzymaj };
}
