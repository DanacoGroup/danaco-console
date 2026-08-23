import type { Transport } from '../polaczenie/gniazdo';
import type { OdczytLacznosci, PortPonawiania } from './lacznosc-okna';

/**
 * Doprowadzenie stanu łącza do sceny okien równoległych.
 *
 * Jedna odpowiedzialność: zamiana subskrypcji transportu na strumień odczytów,
 * które układ rozsyła do nagłówków gniazd. Nic tu nie jest źródłem prawdy —
 * prawdą jest transport, a ten plik wyłącznie go odpytuje.
 *
 * Licznik kolejki musi być prawdziwy, a nie zamrożony. Transport ogłasza
 * `polaczony` przed opróżnieniem kolejki (`polaczenie/gniazdo.ts` →
 * `obsluzOtwarcie`: najpierw `zapiszStan`, potem `oproznijKolejke`), więc
 * odczyt zrobiony w chwili zmiany stanu zamarza na wartości sprzed wysłania.
 * Dlatego po połączeniu odczyt dobija się cyklicznie, aż kolejka spadnie do
 * zera. Tę samą rachubę prowadzi pasek górny
 * (`aplikacja/wskaznik-lacznosci.ts`) — nie druga prawda, tylko drugi pytający
 * tego samego transportu. Wspólnego miejsca dla niej nie ma: pasek zwraca
 * element, nie strumień, a `polaczenie/` leży poza tym pakietem.
 */

/** Odstęp dobijania licznika kolejki po połączeniu, w milisekundach. */
const KROK_DOBIJANIA_MS = 150;

/** Ustawienia wiązania; każde ma wartość domyślną. */
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

  // Transport ogłasza subskrybentowi stan bieżący, więc pierwszy odczyt idzie
  // bez czekania na zmianę.
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
 * Dojście do przebiegu ponowienia, jeśli transport je wystawia.
 *
 * Interfejs `Transport` ma sześć pozycji (`polacz`, `wyslij`, `naRamke`,
 * `naStan`, `stan`, `oczekujace`), a `Gniazdo` trzyma `numerProby`
 * i `zaplanowane` prywatnie — port ponawiania jest więc pytaniem, nie
 * założeniem. Wywołanie `transport.numerProby()` wpisane na sztywno wywróciłoby
 * scenę wszędzie tam, gdzie transport tych metod nie ma; sprawdzenie jest jedno
 * i jawne, a zwrócone `null` scena obsługuje jako brak przebiegu.
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
