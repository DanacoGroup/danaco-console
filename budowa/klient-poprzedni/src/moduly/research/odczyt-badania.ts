import { opisOdmowyBledu } from '../../komponenty/odmowa';
import type { PamiecBadania } from './pamiec-badania';
import type { ZrodloOknaBadania } from './zrodlo-okna-badania';

/**
 * Funkcja odczytuje okno badania z rdzenia dwoma krokami, wskazując okno sesji modułu Research, a następnie pobierając jego stan.
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
