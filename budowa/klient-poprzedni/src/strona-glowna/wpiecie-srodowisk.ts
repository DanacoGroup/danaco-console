import { Command } from '../../../shared/contract';
import type { Kanal } from '../protokol/kanal';
import { pozycjaZeSrodowiska } from './pozycje-srodowisk';
import type { StronaGlowna } from './okno-strona-glowna';

/**
 * Wpięcie wykazu środowisk w stronę główną — źródłem prawdy jest `environment.list`
 * z rdzenia, nie stała klienta.
 *
 * Wykaz zastany w kliencie zostaje jako wartość początkowa: odpowiedź rdzenia
 * przychodzi po pierwszym rysowaniu i wtedy przerysowuje strefę.
 *
 * Odmowa i wykaz pusty nie gaszą ekranu — na miejscu zostaje wykaz zastany,
 * po którym da się wejść do pracy.
 */
export function wepnijSrodowiska(strona: StronaGlowna, kanal: Kanal): void {
  kanal.wyslij(Command.EnvironmentList, {}, (wynik) => {
    const srodowiska = wynik.wynik?.environments;
    if (!wynik.udany || srodowiska === undefined || srodowiska.length === 0) return;
    strona.srodowiska.ustawWykaz(srodowiska.map(pozycjaZeSrodowiska));
  });
}
