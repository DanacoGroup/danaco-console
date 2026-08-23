import type { Transport } from '../polaczenie/gniazdo';
import { utworzPrzelacznikTras, type PrzelacznikTras } from './przelacznik-tras';
import type { Trasa } from './trasy';
import { utworzWskaznikLacznosci, type WskaznikLacznosci } from './wskaznik-lacznosci';

/** Grupa dokładana do akcji paska górnego powłoki środowiska. */
export interface AkcjeTras {
  /** Grupa montowana przy prawej krawędzi paska, przed akcjami powłoki. */
  element: HTMLElement;
  /** Przełącznik widoków — oznaczenie trasy bieżącej. */
  trasy: PrzelacznikTras;
  /** Wskaźnik łączności z rdzeniem. */
  lacznosc: WskaznikLacznosci;
}

/**
 * Dwie rzeczy, których pasek powłoki środowiska sam nie niesie: stan łączności
 * z rdzeniem i przełącznik widoków najwyższego rzędu.
 *
 * Jedna odpowiedzialność: złożenie ich w jedną grupę. Powłoka środowiska ma
 * własny pasek górny wraz z przełącznikiem motywu, powiadomieniami i awatarem;
 * nie dokłada się więc drugiego paska, tylko uzupełnia grupę akcji istniejącego.
 */
export function utworzAkcjeTras(
  transport: Transport,
  naTrase: (trasa: Trasa) => void,
): AkcjeTras {
  const element = document.createElement('div');
  element.className = 'dn-akcje-tras';

  const lacznosc = utworzWskaznikLacznosci(transport);
  const trasy = utworzPrzelacznikTras(naTrase);

  element.append(lacznosc.element, trasy.element);

  return { element, trasy, lacznosc };
}
