/**
 * Jedyny punkt wejścia klienta: składa gniazdo, kanał, sesję i tożsamość, po
 * czym nawiązuje rozmowę z rdzeniem.
 *
 * Widoku ten plik nie buduje i budować nie ma. Znacznik okien pochodzi
 * z biblioteki prototypu — powłokę stawia `zasoby/powloka.js`, wnętrze okna
 * wstrzykuje `vite.config.js` z `design/05-okna/`, a zachowania niosą skrypty
 * biblioteki wczytane przed tym modułem. Warstwa własna klienta ogranicza się
 * do tego, czego prototyp z natury nie niesie: protokołu i połączenia.
 */

import { adresGniazdaRdzenia, adresRdzeniaLokalnego } from './polaczenie/adres-rdzenia.ts';
import { utworzTransport } from './polaczenie/gniazdo.ts';
import { utworzKanal } from './protokol/kanal.ts';
import { utworzSesje } from './protokol/sesja.ts';

/**
 * Adres gniazda rdzenia pochodzi z dokumentu wczytanego po HTTP; w pozostałych
 * przypadkach zostaje pętla zwrotna z portem domyślnym.
 */
function adresRdzenia(): string {
  return adresGniazdaRdzenia(globalThis.location.origin) ?? adresRdzeniaLokalnego();
}

const transport = utworzTransport(adresRdzenia());
const kanal = utworzKanal(transport, utworzSesje());

/*
Kanał wystawiony na obiekcie globalnym: skrypty biblioteki są funkcjami
domkniętymi wczytywanymi przed tym modułem, nie modułami — importu z nich nie ma
i mieć nie może. Wystawienie jest jedynym stykiem między nimi a rdzeniem.
*/
declare global {
  // eslint-disable-next-line no-var
  var DanacoKanal: typeof kanal | undefined;
}
globalThis.DanacoKanal = kanal;

transport.polacz();
