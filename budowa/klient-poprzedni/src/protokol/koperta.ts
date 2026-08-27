import { EnvelopeStatus, type Envelope, type MessageType } from '../../../shared/contract';
import { nowyIdentyfikator } from './identyfikator';

/** Składanie i rozpoznawanie koperty kontraktu, bez definiowania żadnego własnego typu komunikatu tutaj. */
export function zbudujKoperte<T>(typ: MessageType, idSesji: string, tresc: T): Envelope<T> {
  return {
    type: typ,
    id: nowyIdentyfikator('zadanie'),
    ...(idSesji.length > 0 ? { sessionId: idSesji } : {}),
    payload: tresc,
    timestamp: Date.now(),
  };
}

/** Czy koperta jest odpowiedzią na komendę; kontrakt wypełnia pole statusu wyłącznie w odpowiedzi rdzenia. */
export function czyOdpowiedz(koperta: Envelope): boolean {
  return koperta.status !== undefined;
}

/** Czy odpowiedź niesie wynik powodzenia; brak pola statusu traktujemy tutaj jak brak wyniku tej komendy. */
export function czyUdana(koperta: Envelope): boolean {
  return koperta.status === EnvelopeStatus.Ok;
}

/** Treść koperty w kształcie wyznaczonym przez typ komunikatu; rzutowanie jest świadome, transport nie waliduje ładunku. */
export function tresc<T>(koperta: Envelope): T | undefined {
  return koperta.payload as T | undefined;
}
