import type { Transport } from '../polaczenie/gniazdo';
import { utworzPrzelacznikTras, type PrzelacznikTras } from './przelacznik-tras';
import type { Trasa } from './trasy';
import { utworzWskaznikLacznosci, type WskaznikLacznosci } from './wskaznik-lacznosci';

/**
 * Grupa dokładana do akcji paska górnego powłoki środowiska. Niesie własny
 * element montażowy, przełącznik tras oraz wskaźnik łączności z rdzeniem.
 */
export interface AkcjeTras {
  /** Grupa montowana przy prawej krawędzi paska, przed akcjami powłoki. */
  element: HTMLElement;
  /** Przełącznik widoków — oznaczenie trasy bieżącej. */
  trasy: PrzelacznikTras;
  /** Wskaźnik łączności z rdzeniem. */
  lacznosc: WskaznikLacznosci;
}

/**
 * Składa w jedną grupę dwie rzeczy, których pasek powłoki środowiska sam nie
 * niesie: wskaźnik łączności z rdzeniem oraz przełącznik widoków najwyższego
 * rzędu. Grupa uzupełnia pasek istniejący, zamiast zakładać drugi.
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
