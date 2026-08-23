import type { KodSekcjiUstawien } from './sekcje';
import { REJESTR_SEKCJI_USTAWIEN } from './sekcje';

/**
 * Stan ramy Okna Ustawień: wyłącznie sekcja czynna.
 *
 * Wykaz sekcji jest stały (rejestr w `sekcje.ts`, nie katalog rdzenia), więc
 * rama nie potrzebuje odczytu ani fazy ładowania na swoim poziomie — każda
 * sekcja niesie własny stan odczytu wewnątrz siebie. Jedyna rzecz, którą okno
 * pamięta między otwarciami, to sekcja, na której Operator skończył: tak samo
 * jak okno konfiguracji, które nie wraca do pierwszej kategorii przy każdym
 * otwarciu.
 */
export interface StanUstawien {
  /** Kod sekcji czynnej. */
  czynna(): KodSekcjiUstawien;
  /** Ustawia sekcję czynną; kod spoza rejestru jest ignorowany. */
  ustawCzynna(kod: KodSekcjiUstawien): void;
  /** Subskrypcja zmiany sekcji czynnej. */
  naZmiane(sluchacz: (kod: KodSekcjiUstawien) => void): void;
}

export function utworzStanUstawien(): StanUstawien {
  const znaneKody = new Set(REJESTR_SEKCJI_USTAWIEN.map((sekcja) => sekcja.kod));
  const domyslna = REJESTR_SEKCJI_USTAWIEN[0]?.kod ?? 'konta';

  let czynna: KodSekcjiUstawien = domyslna;
  const sluchacze: Array<(kod: KodSekcjiUstawien) => void> = [];

  return {
    czynna: () => czynna,

    ustawCzynna(kod) {
      if (!znaneKody.has(kod) || kod === czynna) return;
      czynna = kod;
      for (const sluchacz of [...sluchacze]) sluchacz(czynna);
    },

    naZmiane: (sluchacz) => void sluchacze.push(sluchacz),
  };
}
