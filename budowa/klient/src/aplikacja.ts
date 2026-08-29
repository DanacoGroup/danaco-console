// Punkt wejścia klienta. Znacznik okien pochodzi z biblioteki `design/zasoby/`.

import { adresGniazdaRdzenia, adresRdzeniaLokalnego } from './polaczenie/adres-rdzenia.ts';
import { utworzTransport } from './polaczenie/gniazdo.ts';
import { utworzKanal } from './protokol/kanal.ts';
import { utworzSesje } from './protokol/sesja.ts';
import { zwiazWejscie } from './wiazanie/wejscie.ts';

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

// Ekran startowy niesie `data-czekaj` i stoi do wywołania `gotowe()`.
function zglosGotowosc(): void {
  const pole = document.querySelector('[data-ekran-startowy]') as
    (HTMLElement & { ekranStartowy?: { gotowe?: () => void } }) | null;
  pole?.ekranStartowy?.gotowe?.();
}

let gotowoscZgloszona = false;
function raz(): void {
  if (gotowoscZgloszona) return;
  gotowoscZgloszona = true;
  zglosGotowosc();
}

transport.naStan(raz);
globalThis.setTimeout(raz, 6000);

/* Wiązanie znacznika Właściciela z komendami rdzenia: nasłuchy na jego
   przyciskach i polach. Nie stawia żadnego elementu. */
zwiazWejscie(kanal);

transport.polacz();
