import { KluczUstawieniaOkna } from './klucze-ustawien';
import { utworzListeWyboru, type OpcjaWyboru } from './lista-wyboru';
import { naniesStanRejestru, opcje } from './model-glowny';
import type { RejestrKanalow } from './rejestr-kanalow';
import type { StanSterowania } from './stan-sterowania';
import type { ZmianaUstawienia } from './zmiana-ustawienia';

const NAZWA = 'Model zapasowy';

/** Pozycja pusta: okno bez modelu zapasowego. Brak ustawienia to nie brak działania. */
const BEZ_ZAPASOWEGO: OpcjaWyboru = { wartosc: '', nazwa: 'Bez modelu zapasowego' };

/**
 * Sterowanie kanałem zapasowym okna.
 *
 * Treść `window.update` nie ma pola na kanał zapasowy, więc wartość idzie
 * ustawieniem poziomu okna (`config.set`, zasięg `window`) — poziom najwęższy,
 * wygrywający z każdym szerszym. Wykaz pochodzi z tego samego rejestru
 * kanałów, co model główny.
 */
export function utworzSterowanieModeluZapasowego(
  stan: StanSterowania,
  ustawienia: ZmianaUstawienia,
  rejestr: RejestrKanalow,
): HTMLElement {
  const lista = utworzListeWyboru(NAZWA, (wartosc) => {
    ustawienia.zapisz(NAZWA, KluczUstawieniaOkna.KanalZapasowy, wartosc);
  });

  function odrysuj(): void {
    lista.pokaz(
      [BEZ_ZAPASOWEGO, ...opcje(rejestr)],
      stan.migawka().ustawienia.kanalZapasowy,
    );
    naniesStanRejestru(lista, rejestr);
  }

  stan.naZmiane(odrysuj);
  rejestr.naZmiane(odrysuj);
  odrysuj();

  return lista.element;
}
