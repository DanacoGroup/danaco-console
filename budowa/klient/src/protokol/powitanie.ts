import {
  Command,
  ErrorCode,
  PROTOCOL_VERSION,
  type ConnectionHelloResponse,
  type ErrorInfo,
} from '../../../shared/contract.ts';
import type { Kanal, Wynik } from './kanal.ts';
import { czyTekst, sprawdzKsztalt } from './ksztalt-odpowiedzi.ts';
import { zapomnijTokenSesji } from './token-sesji.ts';
import type { TozsamoscKlienta } from './tozsamosc-klienta.ts';
import { wywolaj } from './wywolanie.ts';

/*
Powitanie bieżącego połączenia. Rdzeń wiąże sesję bramki z gniazdem właśnie
w powitaniu, więc jest ono jedno na gniazdo: drugie uzgadniałoby tę samą sesję
po raz wtóry. Wołający, który przyjdzie po odpowiedź później, dostaje tę samą
obietnicę zamiast wysyłać komunikat powtórnie.
*/
let biezace: Promise<Wynik<ConnectionHelloResponse>> | null = null;

/**
 * `connection.hello` — powitanie połączenia i uzgodnienie wersji protokołu;
 * jedyna komenda, którą warstwa protokołu wysyła z własnej woli. Poprzedza
 * każdą inną rozmowę z rdzeniem, ustalając wersję protokołu i obsługiwane
 * komendy.
 */
export function zadajPowitanie(
  kanal: Kanal,
  klient: TozsamoscKlienta,
  token?: string,
): Promise<Wynik<ConnectionHelloResponse>> {
  biezace ??= zloz(kanal, klient, token);
  return biezace;
}

/** Znosi zapamiętane powitanie; nowe gniazdo uzgadnia je od nowa. */
export function zapomnijPowitanie(): void {
  biezace = null;
}

/** Składa powitanie i uzgadnia jego odpowiedź: sprawdza kształt, rozstrzyga wersję protokołu i zapamiętuje wykaz komend rdzenia. */
function zloz(
  kanal: Kanal,
  klient: TozsamoscKlienta,
  token?: string,
): Promise<Wynik<ConnectionHelloResponse>> {
  return wywolaj(kanal, Command.ConnectionHello, {
    clientId: klient.id,
    clientVersion: klient.wersja,
    protocolVersion: PROTOCOL_VERSION,
    token,
  })
    .then((wynik) =>
      sprawdzKsztalt(
        wynik,
        Command.ConnectionHello,
        (tresc) => czyTekst(tresc.serverVersion) && czyTekst(tresc.protocolVersion),
      ),
    )
    .then((wynik) => uzgodnij(kanal, wynik, token ?? ''));
}

/**
 * Rozstrzyga powitanie po stronie klienta.
 *
 * Wersja protokołu rozbieżna zatrzymuje wejście: kształt komunikatu zmienia się
 * razem z wersją, więc rozjazd wyszedłby dopiero na pierwszej komendzie, której
 * pola się rozeszły. Wykaz komend rdzenia zostaje w kanale jako sprawdzenie
 * przed wywołaniem.
 */
function uzgodnij(
  kanal: Kanal,
  wynik: Wynik<ConnectionHelloResponse>,
  token: string,
): Wynik<ConnectionHelloResponse> {
  const tresc = wynik.wynik;
  if (!wynik.udany || tresc === undefined) return wynik;
  if (tresc.protocolVersion !== PROTOCOL_VERSION) {
    return { udany: false, blad: bladWersji(tresc.protocolVersion) };
  }
  /* Token odrzucony przy powitaniu jest martwy — rdzeń mówi o tym wprost polem
     `authenticated`. Milczenie tego pola znaczy rdzeń bez wpiętej bramki, więc
     tokenu nie zdejmuje. Zostawiony wracałby przy każdym kolejnym gnieździe. */
  if (token !== '' && tresc.authenticated === false) zapomnijTokenSesji();
  kanal.zapamietajKomendy(tresc.commands ?? []);
  return wynik;
}

/** Odmowa powitania, gdy rdzeń mówi inną wersją protokołu niż ta, którą zna klient. */
function bladWersji(wersjaRdzenia: string): ErrorInfo {
  return {
    code: ErrorCode.Conflict,
    message: `Rdzeń mówi protokołem ${wersjaRdzenia}, klient protokołem ${PROTOCOL_VERSION}`,
    retryable: false,
  };
}
