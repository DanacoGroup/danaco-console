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

/*
Zgłoszenie gotowości ekranowi startowemu.

Ekran biblioteki (`zasoby/ekran-startowy.js`) niesie atrybut `data-czekaj`
i trzyma zasłonę, dopóki nie usłyszy `gotowe()`. Zgłasza to warstwa połączenia,
bo tylko ona wie, kiedy rdzeń odpowiedział — bez tego zasłona nie schodzi nigdy
i Operator widzi puste tło.

Zabezpieczenie czasowe zdejmuje ją także wtedy, gdy rdzeń nie odpowiada:
zasłona bez końca ukrywałaby przed Operatorem samą informację o usterce.
*/
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
  pokazScene('dostep');
}

/*
Wybór sceny drogi wejścia.

Prototyp `przeplyw-wejscia.html` jest stroną przeglądową: stawia wszystkie trzy
sceny jedna pod drugą, żeby dało się je ocenić naraz. W aplikacji naraz stoi
jedna — pozostałe schodzą z widoku. Bez tego pierwsza scena zasłania resztę
i Operator widzi puste tło.
*/
function pokazScene(okno: 'uruchomienie' | 'dostep' | 'przygotowanie'): void {
  for (const scena of document.querySelectorAll<HTMLElement>('.we-scena')) {
    scena.hidden = scena.querySelector(`[data-wejscie-okno="${okno}"]`) === null;
  }
}

pokazScene('uruchomienie');

transport.naStan(raz);
globalThis.setTimeout(raz, 6000);

/* Wiązanie znacznika Właściciela z komendami rdzenia: nasłuchy na jego
   przyciskach i polach. Nie stawia żadnego elementu. */
zwiazWejscie(kanal);

transport.polacz();
