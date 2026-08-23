import { WindowRole } from '../../../shared/contract';
import { ID_GNIAZD, type IdGniazda } from './identyfikatory';

/**
 * Role przyjmowane przy zmianie liczby okien.
 *
 * Brak ustawienia oznacza wartość domyślną, nie blokadę.
 * Domyślne nie są bramą: `nadajRole` przykrywa je w każdej chwili, a raz
 * nadana ręcznie rola nie zostaje odebrana przy kolejnej zmianie liczby.
 *
 * Dobór odpowiada obsadzie pętli multitaskingu:
 *
 *   jedno okno   → samodzielne, bo pętla nie ma z kim biec;
 *   dwa okna     → koordynator i wykonawca;
 *   trzy okna    → koordynator i dwaj wykonawcy;
 *   cztery okna  → koordynator, dwaj wykonawcy i analityk.
 *
 * Czwarte gniazdo bierze rolę samodzielną, nie wykonawczą. Obsada multitaskingu
 * szuka analityka właśnie wśród okien samodzielnych, a wykonawców liczy do
 * dwóch (`LICZBA_WYKONAWCOW`) — trzeci wykonawca byłby oknem, którego pętla
 * nie obsłuży, i zniknąłby z obsady jako nadmiarowy.
 */
export function rolaDomyslna(id: IdGniazda, liczba: number): WindowRole {
  if (liczba <= 1) return WindowRole.Standalone;
  if (id === ID_GNIAZD[0]) return WindowRole.Coordinator;
  return id === ID_GNIAZD[3] ? WindowRole.Standalone : WindowRole.Executor;
}
