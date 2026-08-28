import { ChunkKind, MessageRole, type StreamChunkEvent } from '../../../shared/contract';
import type { OknoKomunikacji } from './okno';
import { nazwaNadawcy, PERSONA_POWLOKI, wpisSystemowy } from './wpis';

// Wyświetlenie ruchu rdzenia w oknie komunikacji, wspólne dla dwóch odbiorców fragmentów strumienia.

/**
 * Tożsamość mówiącego przypisana roli wiadomości: kanał modelu podpisuje wyłącznie rolę assistant, a wiadomość systemowa i wynik narzędzia mówią własnym głosem, nie głosem modelu, który ich nie wypowiedział.
 */
export function personaRoli(rola: MessageRole, kanalModelu: string): string {
  if (rola === MessageRole.Assistant) return kanalModelu;
  if (rola === MessageRole.System) return PERSONA_POWLOKI;
  return nazwaNadawcy(rola);
}

/**
 * Fragment strumienia w historii okna: tekst dokłada się do bieżącej wypowiedzi persony, a fragment innego rodzaju zostaje pokazany jako wpis systemowy, żeby żadna treść strumienia nie znikła bez śladu.
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
