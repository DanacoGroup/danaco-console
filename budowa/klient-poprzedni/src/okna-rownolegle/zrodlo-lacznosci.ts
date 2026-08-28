import type { Transport } from '../polaczenie/gniazdo';
import type { OdczytLacznosci, PortPonawiania } from './lacznosc-okna';

/**
 * Odstęp w milisekundach, w którym licznik kolejki dobija po połączeniu do wartości
 * prawdziwej, bo transport ogłasza połączenie przed opróżnieniem kolejki.
 */
const KROK_DOBIJANIA_MS = 150;

/**
 * Ustawienia wiązania łącza ze sceną okien równoległych; każde pole ma wartość
 * domyślną używaną, gdy wołający jej nie poda.
 */
export interface OpcjeZrodlaLacznosci {
  krokDobijaniaMs?: number;
}

/**
 * Wiąże transport z odbiorcą odczytów. Zwraca odłączenie, które trzeba wywołać
 * przy zdejmowaniu sceny, bo subskrypcja transportu przeżywa usunięcie węzła
 * z dokumentu.
 */
export function zwiazLacznoscUkladu(
  transport: Transport,
  przyOdczycie: (odczyt: OdczytLacznosci) => void,
  opcje: OpcjeZrodlaLacznosci = {},
): () => void {
  const krok = opcje.krokDobijaniaMs ?? KROK_DOBIJANIA_MS;
  let uchwyt: ReturnType<typeof setInterval> | null = null;
  let czynne = true;

  function zatrzymaj(): void {
    if (uchwyt !== null) {
      clearInterval(uchwyt);
      uchwyt = null;
    }
  }

  function oglos(): void {
    if (!czynne) return;
    przyOdczycie({ stan: transport.stan(), oczekujace: transport.oczekujace() });
  }

  // Transport ogłasza subskrybentowi stan bieżący; pierwszy odczyt idzie bez czekania na zmianę.
  const odsubskrybuj = transport.naStan((stan) => {
    oglos();
    if (stan === 'polaczony' && transport.oczekujace() > 0) {
      if (uchwyt === null) {
        uchwyt = setInterval(() => {
          oglos();
          if (transport.oczekujace() === 0) zatrzymaj();
        }, krok);
      }
      return;
    }
    zatrzymaj();
  });

  return () => {
    czynne = false;
    zatrzymaj();
    odsubskrybuj();
  };
}

/**
 * Dojście do przebiegu ponowienia, jeśli transport je wystawia; port ponawiania jest
 * pytaniem o metody, nie założeniem o ich obecności.
 */
export function portPonawiania(transport: Transport): PortPonawiania | null {
  const kandydat = transport as Partial<PortPonawiania>;
  if (
    typeof kandydat.numerProby !== 'function'
    || typeof kandydat.zaPonowieniem !== 'function'
    || typeof kandydat.ponowTeraz !== 'function'
  ) {
    return null;
  }
  return {
    numerProby: () => Number(kandydat.numerProby?.() ?? 0),
    zaPonowieniem: () => Number(kandydat.zaPonowieniem?.() ?? 0),
    ponowTeraz: () => kandydat.ponowTeraz?.(),
  };
}
