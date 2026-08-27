import type { StanRozmowy } from './stan-rozmowy';

/**
 * Zegar ciszy — jedyna warstwa okna rozmowy, która mierzy czas, bo rdzeń nie ma zdarzenia
 * zerwanego strumienia.
 */
export interface ZegarCiszy {
  /** Uruchamia odliczanie; wywołanie powtórne nie zakłada drugiego interwału. */
  uruchom: () => void;
  /** Zatrzymuje odliczanie; wywołanie na zatrzymanym nic nie robi. */
  zatrzymaj: () => void;
}

/**
 * Zależności zegara ciszy: bieżący stan tury okna rozmowy oraz kierunek jego zmiany w
 * czasie trwania.
 */
export interface OtoczenieZegara {
  /** Stan bieżący rozmowy; zegar go czyta, nie trzyma własnej kopii. */
  stan: () => StanRozmowy;
  /** Chwila ostatniego znaku życia tury, w milisekundach epoki. */
  znakZycia: () => number;
  /** Przestawia stan — jedyną drogą, którą zegar zmienia cokolwiek. */
  ustawStan: (zmiana: Partial<StanRozmowy>) => void;
}

/**
 * Takt odliczania zegara ciszy wynosi jedną sekundę: krócej byłoby migotaniem, dłużej
 * nieuwagą operatora.
 */
const taktCiszy = 1000;

/**
 * Zakłada zegar ciszy nad podanym otoczeniem stanu bieżącej rozmowy okna komunikacji z
 * rdzeniem serwera.
 */
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
