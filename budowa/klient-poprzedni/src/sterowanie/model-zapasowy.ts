import { KluczUstawieniaOkna } from './klucze-ustawien';
import { utworzListeWyboru, type OpcjaWyboru } from './lista-wyboru';
import { naniesStanRejestru, opcje } from './model-glowny';
import type { RejestrKanalow } from './rejestr-kanalow';
import type { StanSterowania } from './stan-sterowania';
import type { ZmianaUstawienia } from './zmiana-ustawienia';

const NAZWA = 'Model zapasowy';

/** Pozycja pusta: okno bez modelu zapasowego. Brak ustawienia to nie brak działania, tylko wartość domyślna. */
const BEZ_ZAPASOWEGO: OpcjaWyboru = { wartosc: '', nazwa: 'Bez modelu zapasowego' };

/** Sterowanie kanałem zapasowym okna, zapisywane ustawieniem poziomu okna, z wykazem tego samego rejestru kanałów co model główny. */
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
