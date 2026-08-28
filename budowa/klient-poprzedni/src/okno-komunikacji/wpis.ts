import { MessageRole } from '../../../shared/contract';

/**
 * Pojedynczy wpis historii okna komunikacji niesie rolę nadawcy wprost z kontraktu oraz personę — tożsamość mówiącego wewnątrz tej roli, bo w jednym oknie wypowiada się wiele person po stronie modelu.
 */
export interface Wpis {
  nadawca: MessageRole;
  persona: string;
  tresc: string;
  znacznikCzasu: number;
}

/** Tożsamość, którą podpisuje się sama powłoka aplikacji, gdy zwraca się bezpośrednio do operatora okna. */
export const PERSONA_POWLOKI = 'Danaco Console';

/** Nazwa nadawcy w postaci prezentacyjnej, pokazywana operatorowi zamiast surowej wartości roli kontraktu. */
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

// Przekład roli kontraktu na klasę biblioteki wystawia historia.ts wraz z wymaganym medalionem.

/** Wpis operatora z bieżącym znacznikiem czasu, tworzony w chwili wysłania treści z pola wypowiedzi okna. */
export function wpisOperatora(tresc: string, persona: string): Wpis {
  return { nadawca: MessageRole.User, persona, tresc, znacznikCzasu: Date.now() };
}

/** Wpis modelu — treść odebrana z rdzenia albo złożona ze strumienia kolejnych fragmentów jego odpowiedzi. */
export function wpisModelu(tresc: string, persona: string): Wpis {
  return { nadawca: MessageRole.Assistant, persona, tresc, znacznikCzasu: Date.now() };
}

/** Wpis systemowy obejmujący zdarzenia transportu oraz komunikaty pochodzące bezpośrednio od rdzenia platformy. */
export function wpisSystemowy(tresc: string): Wpis {
  return {
    nadawca: MessageRole.System,
    persona: PERSONA_POWLOKI,
    tresc,
    znacznikCzasu: Date.now(),
  };
}

/** Godzina wpisu w formacie prezentacyjnym, czytelnym dla operatora przeglądającego historię tej rozmowy. */
export function godzinaWpisu(wpis: Wpis): string {
  return new Date(wpis.znacznikCzasu).toLocaleTimeString('pl-PL', {
    hour: '2-digit',
    minute: '2-digit',
    second: '2-digit',
  });
}
