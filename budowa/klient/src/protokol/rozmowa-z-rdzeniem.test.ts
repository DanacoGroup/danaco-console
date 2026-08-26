import {
  PROTOCOL_VERSION,
  type ConnectionHelloResponse,
} from '../../../shared/contract.ts';
import { adresRdzeniaLokalnego } from '../polaczenie/adres-rdzenia.ts';
import { utworzTransport } from '../polaczenie/gniazdo.ts';
import { bieg, sprawdz } from '../sprawdzian.ts';
import { utworzKanal } from './kanal.ts';
import { zadajPowitanie } from './powitanie.ts';
import { utworzSesje } from './sesja.ts';
import { tozsamoscKlienta } from './tozsamosc-klienta.ts';

/**
 * Rozmowa z rdzeniem uruchomionym naprawdę.
 *
 * Sprawdzian wymaga rdzenia nasłuchującego pod adresem lokalnym i bez niego
 * nie ma czego zmierzyć. Milczenie rdzenia kończy się tu niepowodzeniem
 * nazywającym przeszkodę, nie pominięciem: sprawdzian, który sam siebie
 * odpuszcza przy braku rdzenia, wygląda potem tak samo jak sprawdzian zdany.
 */

/** Górna granica oczekiwania na odpowiedź rdzenia. */
const GRANICA_MS = 10_000;

/** Rozstrzyga obietnicę albo przerywa ją z nazwaniem przeszkody. */
function wGranicyCzasu<T>(obietnica: Promise<T>, czynnosc: string): Promise<T> {
  return Promise.race([
    obietnica,
    new Promise<T>((_, odrzuc) =>
      setTimeout(
        () =>
          odrzuc(
            new Error(
              `${czynnosc}: rdzeń nie odpowiedział w ciągu ${GRANICA_MS} ms pod adresem ` +
                `${adresRdzeniaLokalnego()} — uruchom rdzeń przed tym sprawdzianem`,
            ),
          ),
        GRANICA_MS,
      ),
    ),
  ]);
}

await bieg('rozmowa z rdzeniem', {
  async 'klient łączy się z rdzeniem, wita go i odczytuje wersję protokołu'() {
    const transport = utworzTransport(adresRdzeniaLokalnego());
    try {
      await wGranicyCzasu(
        new Promise<void>((rozstrzygnij) => {
          transport.naStan((stan) => {
            if (stan === 'polaczony') rozstrzygnij();
          });
          transport.polacz();
        }),
        'nawiązanie połączenia',
      );

      const kanal = utworzKanal(transport, utworzSesje());
      const wynik = await wGranicyCzasu(
        zadajPowitanie(kanal, tozsamoscKlienta()),
        'powitanie połączenia',
      );

      sprawdz(wynik.udany, `rdzeń odmówił powitania: ${JSON.stringify(wynik.blad)}`);
      const odpowiedz = wynik.wynik;
      sprawdz(odpowiedz !== undefined, 'rdzeń oddał powitanie bez treści');
      console.log(`           odpowiedź rdzenia: ${JSON.stringify(zwiezle(odpowiedz))}`);
      sprawdz(
        odpowiedz?.protocolVersion === PROTOCOL_VERSION,
        `wersja protokołu rdzenia ${odpowiedz?.protocolVersion} wobec klienta ${PROTOCOL_VERSION}`,
      );
      sprawdz(
        typeof odpowiedz?.serverVersion === 'string' && odpowiedz.serverVersion.length > 0,
        'rdzeń nie podał własnej wersji',
      );
    } finally {
      transport.rozlacz();
    }
  },
});

/** Odpowiedź powitania bez wykazu komend — ten liczy setki pozycji. */
function zwiezle(odpowiedz: ConnectionHelloResponse | undefined): Record<string, unknown> {
  if (odpowiedz === undefined) return {};
  const { commands, ...reszta } = odpowiedz;
  return { ...reszta, commands: commands === undefined ? commands : `${commands.length} nazw` };
}
