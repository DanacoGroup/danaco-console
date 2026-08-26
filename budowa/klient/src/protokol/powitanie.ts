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
 * `connection.hello` — powitanie połączenia i uzgodnienie wersji protokołu.
 *
 * Jedyna komenda, którą warstwa protokołu wysyła z własnej woli. Powitanie
 * poprzedza każdą inną rozmowę z rdzeniem: dopiero z odpowiedzi klient
 * dowiaduje się, jaką wersję protokołu zna rdzeń i które komendy ta wersja
 * rdzenia obsługuje. Rozstrzygnięcie, co zrobić z rozjazdem wersji, należy do
 * warstwy wyższej — protokół oddaje odpowiedź w kształcie kontraktu.
 *
 * Wersja protokołu w żądaniu pochodzi ze stałej `PROTOCOL_VERSION` artefaktu
 * kontraktu, nie z literału: klient przedstawia się tą wersją, z którą został
 * zbudowany.
 *
 * Token wiąże połączenie z sesją bramki. Jego brak nie jest błędem — rdzeń
 * odpowiada wtedy `authenticated: false`. Token dostarcza wołający, bo magazyn
 * sesji bramki nie należy do warstwy protokołu.
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
