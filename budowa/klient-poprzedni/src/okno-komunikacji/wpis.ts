import { MessageRole } from '../../../shared/contract';

/**
 * Pojedynczy wpis historii okna komunikacji.
 *
 * `nadawca` jest rolą wiadomości wprost z kontraktu (`MessageRole`), dzięki
 * czemu wpis odebrany zdarzeniem `message.changed` nie wymaga przekładu.
 * `persona` jest tożsamością mówiącego wewnątrz tej roli — nazwą modelu,
 * agenta albo roli przypisanej oknu. Rozróżnienie jest konieczne, ponieważ
 * w jednym oknie wypowiada się wiele person po stronie modelu.
 */
export interface Wpis {
  nadawca: MessageRole;
  persona: string;
  tresc: string;
  znacznikCzasu: number;
}

/** Tożsamość, którą podpisuje się sama powłoka, gdy mówi do operatora. */
export const PERSONA_POWLOKI = 'Danaco Console';

/** Nazwa nadawcy w postaci prezentacyjnej. */
export function nazwaNadawcy(nadawca: MessageRole): string {
  switch (nadawca) {
    case MessageRole.User:
      return 'Operator';
    case MessageRole.Assistant:
      return 'Model';
    case MessageRole.System:
      return 'System';
    case MessageRole.Tool:
      return 'Narzędzie';
  }
}

/* Przekład roli kontraktu na klasę biblioteki (`user → .dn-wpis--czlowiek`,
   `assistant → .dn-wpis--inteligencja`, `system` i `tool` → `.dn-wpis--system`)
   wystawia `historia.ts`, wraz z medalionem, którego wymaga dwukolumnowa siatka
   `.dn-wpis`. */

/** Wpis operatora z bieżącym znacznikiem czasu. */
export function wpisOperatora(tresc: string, persona: string): Wpis {
  return { nadawca: MessageRole.User, persona, tresc, znacznikCzasu: Date.now() };
}

/** Wpis modelu — treść odebrana z rdzenia albo złożona ze strumienia. */
export function wpisModelu(tresc: string, persona: string): Wpis {
  return { nadawca: MessageRole.Assistant, persona, tresc, znacznikCzasu: Date.now() };
}

/** Wpis systemowy — zdarzenia transportu i komunikaty rdzenia. */
export function wpisSystemowy(tresc: string): Wpis {
  return {
    nadawca: MessageRole.System,
    persona: PERSONA_POWLOKI,
    tresc,
    znacznikCzasu: Date.now(),
  };
}

/** Godzina wpisu w formacie prezentacyjnym. */
export function godzinaWpisu(wpis: Wpis): string {
  return new Date(wpis.znacznikCzasu).toLocaleTimeString('pl-PL', {
    hour: '2-digit',
    minute: '2-digit',
    second: '2-digit',
  });
}
