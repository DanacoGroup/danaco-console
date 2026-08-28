import { zdarzenieNieznanej, type Envelope } from '../../../shared/contract';

/** Ramka tekstowa WebSocket: zapis i odczyt koperty w formacie JSON, fail-open przy ramce nieczytelnej. */
export function zapiszRamke(koperta: Envelope): string {
  return JSON.stringify(koperta);
}

export function odczytajRamke(ramka: string): Envelope {
  try {
    const odczytane: unknown = JSON.parse(ramka);
    return czyKoperta(odczytane) ? odczytane : kopertaNierozpoznana(odczytane);
  } catch {
    return kopertaNierozpoznana(ramka);
  }
}

/** Sprawdza, czy odczytana wartość niesie pola obowiązkowe koperty kontraktu, przede wszystkim jej typ. */
function czyKoperta(wartosc: unknown): wartosc is Envelope {
  if (typeof wartosc !== 'object' || wartosc === null) return false;
  const kandydat = wartosc as Record<string, unknown>;
  return typeof kandydat['type'] === 'string' && kandydat['type'].length > 0;
}

/**
 * Koperta zastępcza dla ramki nierozpoznanej.
 *
 * Zdarzenie zapasowe wskazuje kontrakt — klient nie wybiera go samodzielnie.
 */
function kopertaNierozpoznana(tresc: unknown): Envelope {
  return {
    type: zdarzenieNieznanej(''),
    id: '',
    payload: tresc,
    timestamp: Date.now(),
  };
}
