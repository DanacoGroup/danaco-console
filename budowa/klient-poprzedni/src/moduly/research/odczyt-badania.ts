import { opisOdmowyBledu } from '../../komponenty/odmowa';
import type { PamiecBadania } from './pamiec-badania';
import type { ZrodloOknaBadania } from './zrodlo-okna-badania';

/**
 * Odczyt okna badania z rdzenia — jedyna droga modułu Research do treści.
 *
 * Jedna odpowiedzialność: dwa kroki odczytu i przełożenie ich niepowodzeń na
 * fazę pamięci. Wydzielone z pamięci, bo pamięć nie zna kontraktu, a odczyt zna
 * wyłącznie kontrakt — rozdzielenie trzyma obie rzeczy w rozmiarze, w którym
 * dają się przeczytać naraz.
 *
 * Kroki są dwa, bo mówią o dwóch różnych rzeczach: `window.list` mówi, które
 * okno sesji należy do modułu Research, a `window.state.get` — co w nim jest.
 * Klient nie zgaduje
 * identyfikatora okna — bez wskazania rdzenia zostaje uczciwy stan pusty
 * nazywający brak, nie zaszyta wartość.
 */
export function utworzOdczytBadania(
  okna: ZrodloOknaBadania,
  pamiec: PamiecBadania,
): (idSesji: string) => Promise<void> {
  return async (idSesji) => {
    pamiec.ustawFaze('odczyt', '');

    const wykaz = await okna.oknaBadania(idSesji);
    if (!wykaz.udany || wykaz.wynik === undefined) {
      pamiec.ustawFaze('blad', opisOdmowyBledu('Odczyt okien badania', wykaz.blad));
      return;
    }

    const wskazane = wykaz.wynik[0]?.id ?? '';
    pamiec.ustawOkno(wskazane);
    if (wskazane === '') {
      pamiec.ustawFaze(
        'gotowe',
        'Rdzeń nie wskazał ani jednego okna modułu Research w tej sesji.',
      );
      return;
    }

    const stan = await okna.stanOkna(wskazane);
    if (!stan.udany || stan.wynik === undefined) {
      pamiec.ustawFaze('blad', opisOdmowyBledu('Odczyt stanu okna badania', stan.blad, stan.nieznanyTyp));
      return;
    }
    if (pamiec.zakres() === '') pamiec.wchlonZakres(stan.wynik.window.title ?? '', pamiec.etapy());
    pamiec.ustawFaze('gotowe', '');
  };
}
