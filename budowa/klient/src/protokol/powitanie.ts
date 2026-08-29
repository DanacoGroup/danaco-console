import {
  Command,
  PROTOCOL_VERSION,
  type ConnectionHelloResponse,
} from '../../../shared/contract.ts';
import type { Kanal, Wynik } from './kanal.ts';
import { czyTekst, sprawdzKsztalt } from './ksztalt-odpowiedzi.ts';
import type { TozsamoscKlienta } from './tozsamosc-klienta.ts';
import { wywolaj } from './wywolanie.ts';

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
  return wywolaj(kanal, Command.ConnectionHello, {
    clientId: klient.id,
    clientVersion: klient.wersja,
    protocolVersion: PROTOCOL_VERSION,
    token,
  }).then((wynik) =>
    sprawdzKsztalt(
      wynik,
      Command.ConnectionHello,
      (tresc) => czyTekst(tresc.serverVersion) && czyTekst(tresc.protocolVersion),
    ),
  );
}
