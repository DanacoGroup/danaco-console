import { Command } from '../../../shared/contract';
import type { Kanal } from '../protokol/kanal';
import { pozycjaZeSrodowiska } from './pozycje-srodowisk';
import type { StronaGlowna } from './okno-strona-glowna';

/** Wpięcie wykazu środowisk do strony głównej pobiera listę z rdzenia po pierwszym rysowaniu i zastępuje nią wykaz początkowy, zachowując go przy odmowie. */
export function wepnijSrodowiska(strona: StronaGlowna, kanal: Kanal): void {
  kanal.wyslij(Command.EnvironmentList, {}, (wynik) => {
    const srodowiska = wynik.wynik?.environments;
    if (!wynik.udany || srodowiska === undefined || srodowiska.length === 0) return;
    strona.srodowiska.ustawWykaz(srodowiska.map(pozycjaZeSrodowiska));
  });
}
