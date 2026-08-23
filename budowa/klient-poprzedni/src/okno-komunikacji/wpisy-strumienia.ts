import { ChunkKind, MessageRole, type StreamChunkEvent } from '../../../shared/contract';
import type { OknoKomunikacji } from './okno';
import { nazwaNadawcy, PERSONA_POWLOKI, wpisSystemowy } from './wpis';

/**
 * Wyświetlenie ruchu przychodzącego rdzenia w oknie komunikacji.
 *
 * Miejsce wspólne dla dwóch odbiorców: `przeplyw-komunikatow.ts` prowadzi jedno
 * okno sesji uzgodnione z rdzeniem przy starcie powłoki, a panel modułu prowadzi
 * okna zakładane osobno. Obaj biorą obie funkcje stąd, więc rozpoznanie persony
 * i pokazanie fragmentu mają jedną postać.
 */

/**
 * Tożsamość mówiącego przypisana roli wiadomości.
 *
 * Kanał modelu podpisuje wyłącznie rolę `assistant`. Wiadomość systemowa mówi
 * głosem powłoki, a wynik narzędzia — głosem narzędzia; podpisanie ich kanałem
 * modelu przypisywałoby modelowi zdania, których nie wypowiedział, i kazałoby
 * plakietce kłamać przy wpisie klasy neutralnej.
 *
 * Rola `tool` powinna nieść w plakietce nazwę narzędzia, ale kontrakt tej nazwy
 * przy wiadomości nie przenosi — `Message` nie ma takiego pola. Do czasu, aż je
 * dostanie, plakietka niesie nazwę roli, a nie zmyśloną nazwę narzędzia.
 */
export function personaRoli(rola: MessageRole, kanalModelu: string): string {
  if (rola === MessageRole.Assistant) return kanalModelu;
  if (rola === MessageRole.System) return PERSONA_POWLOKI;
  return nazwaNadawcy(rola);
}

/**
 * Fragment strumienia w historii okna.
 *
 * Tekst dokłada się do bieżącej wypowiedzi persony; fragment innego rodzaju
 * (tok rozumowania, wywołanie narzędzia, błąd kanału) zostaje pokazany jako
 * wpis systemowy, żeby żadna treść strumienia nie znikła bez śladu.
 */
export function pokazFragment(
  okno: OknoKomunikacji,
  persona: string,
  fragment: StreamChunkEvent,
): void {
  const tekst = fragment.text ?? '';
  if (fragment.kind === ChunkKind.Text) {
    okno.dopiszFragment(persona, tekst);
    return;
  }
  okno.dopisz(wpisSystemowy(tekst.length > 0 ? `${fragment.kind}: ${tekst}` : fragment.kind));
}
