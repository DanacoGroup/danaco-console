import { EnvelopeStatus, type Envelope, type MessageType } from '../../../shared/contract.ts';
import { nowyIdentyfikator } from './identyfikator.ts';

/**
 * Składanie i rozpoznawanie koperty kontraktu; kształt koperty pochodzi
 * w całości ze współdzielonego kontraktu. Odpowiada za nadanie kopercie
 * wychodzącej identyfikatora i czasu nadania oraz za odczyt pól odpowiedzi.
 */
export function zbudujKoperte<T>(typ: MessageType, idSesji: string, tresc: T): Envelope<T> {
  return {
    type: typ,
    id: nowyIdentyfikator('zadanie'),
    ...(idSesji.length > 0 ? { sessionId: idSesji } : {}),
    payload: tresc,
    timestamp: Date.now(),
  };
}

/**
 * Czy koperta jest odpowiedzią na komendę.
 *
 * Kontrakt wypełnia `status` wyłącznie w odpowiedzi, więc obecność tego pola
 * odróżnia odpowiedź od zdarzenia i od fragmentu strumienia.
 */
export function czyOdpowiedz(koperta: Envelope): boolean {
  return koperta.status !== undefined;
}

/** Czy odpowiedź niesie wynik — sprawdzenie działa na podstawie statusu koperty; brak statusu traktowany jest jak brak wyniku. */
export function czyUdana(koperta: Envelope): boolean {
  return koperta.status === EnvelopeStatus.Ok;
}

/**
 * Treść koperty w kształcie wyznaczonym przez typ komunikatu.
 *
 * Rzutowanie jest świadome: warstwa połączenia nie waliduje ładunku, a błąd
 * kształtu dotyczy wyłącznie bieżącego komunikatu.
 */
export function tresc<T>(koperta: Envelope): T | undefined {
  return koperta.payload as T | undefined;
}
