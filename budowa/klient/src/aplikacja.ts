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
function zglosGotowosc(): boolean {
  const pole = document.querySelector('[data-ekran-startowy]') as
    (HTMLElement & { ekranStartowy?: { gotowe?: () => void } }) | null;
  const gotowe = pole?.ekranStartowy?.gotowe;
  if (gotowe === undefined) return false;
  gotowe();
  return true;
}

/* Składnik ekranu startowego staje dopiero przy DOMContentLoaded, a pierwszy
   stan transportu potrafi paść wcześniej. Sygnał niedoręczony jest ponawiany —
   zgubiony zostawiałby ekran startowy czekający na `gotowe()` bez końca. */
let gotowoscZgloszona = false;
let ponowienie: ReturnType<typeof setTimeout> | undefined;
function raz(): void {
  if (gotowoscZgloszona) return;
  if (zglosGotowosc()) {
    gotowoscZgloszona = true;
    return;
  }
  if (ponowienie === undefined) {
    ponowienie = globalThis.setTimeout(() => {
      ponowienie = undefined;
      raz();
    }, 100);
  }
}

transport.naStan(raz);
globalThis.setTimeout(raz, 6000);

/* Wiązanie znacznika Właściciela z komendami rdzenia: nasłuchy na jego
   przyciskach i polach. Nie stawia żadnego elementu. */
zwiazWejscie(kanal);

transport.polacz();
