import { WindowRole } from '../../../shared/contract';
import { ID_GNIAZD, type IdGniazda } from './identyfikatory';

/**
 * Moduł przypisuje domyślną rolę oknu na podstawie liczby otwartych okien i jego identyfikatora, gdy wywołanie nie ustawia roli jawnie.
 */
export function rolaDomyslna(id: IdGniazda, liczba: number): WindowRole {
  if (liczba <= 1) return WindowRole.Standalone;
  if (id === ID_GNIAZD[0]) return WindowRole.Coordinator;
  return id === ID_GNIAZD[3] ? WindowRole.Standalone : WindowRole.Executor;
}
