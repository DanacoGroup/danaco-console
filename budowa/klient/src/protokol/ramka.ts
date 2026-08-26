import { zdarzenieNieznanej, type Envelope } from '../../../shared/contract.ts';

/**
 * Ramka tekstowa WebSocket: zapis i odczyt koperty w formacie JSON.
 *
 * Odczyt jest fail-open: ramka nieczytelna albo o kształcie
 * niezgodnym z kopertą nie zrywa połączenia i nie blokuje sesji — wraca jako
 * zdarzenie `*.unknown` z zachowaniem treści surowej w ładunku.
 */
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

/** Sprawdza, czy odczytana wartość ma pola obowiązkowe koperty kontraktu. */
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
