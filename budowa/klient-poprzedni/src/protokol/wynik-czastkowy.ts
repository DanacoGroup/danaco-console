import type { Wynik } from './kanal';

/** Przeniesienie udanego wyniku komendy na jego wycinek, z zachowaniem odmowy rdzenia całkiem nietkniętej. */
export function przenies<Z, W>(wynik: Wynik<Z>, wybierz: (tresc: Z) => W): Wynik<W> {
  if (!wynik.udany || wynik.wynik === undefined) {
    return { udany: false, ...(wynik.blad === undefined ? {} : { blad: wynik.blad }) };
  }
  return { udany: true, wynik: wybierz(wynik.wynik) };
}
