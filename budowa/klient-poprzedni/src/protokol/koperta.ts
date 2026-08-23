import { EnvelopeStatus, type Envelope, type MessageType } from '../../../shared/contract';
import { nowyIdentyfikator } from './identyfikator';

/**
 * Składanie i rozpoznawanie koperty kontraktu.
 *
 * Kształt koperty pochodzi w całości z `shared/contract.ts` — plik nie
 * definiuje własnego typu komunikatu i nie powiela ani jednego literału nazwy.
 * Odpowiada wyłącznie za nadanie kopercie wychodzącej identyfikatora oraz
 * czasu nadania i za odczytanie pól odpowiedzi z koperty przychodzącej.
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

/** Czy odpowiedź niesie wynik. Brak statusu traktujemy jak brak wyniku. */
export function czyUdana(koperta: Envelope): boolean {
  return koperta.status === EnvelopeStatus.Ok;
}

/**
 * Treść koperty w kształcie wyznaczonym przez typ komunikatu.
 *
 * Rzutowanie jest świadome: warstwa transportu nie waliduje ładunku, a błąd
 * kształtu dotyczy wyłącznie bieżącego komunikatu.
 */
export function tresc<T>(koperta: Envelope): T | undefined {
  return koperta.payload as T | undefined;
}
