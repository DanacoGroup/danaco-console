import { elementIkony } from '../ikony/ikony';
import { nazwaStanuPary } from './etykiety-ukladu';
import { wygladStanu, type StanPary } from './stan-pary';

/** Plakietka stanu pętli — kropka, ikona i etykieta. */
export interface PlakietkaStanu {
  element: HTMLElement;
  /** Przestawia plakietkę na inny stan pary. */
  ustaw(stan: StanPary): void;
}

/**
 * Plakietka stanu pary koordynator–wykonawca.
 *
 * Każdy stan ma odrębną ikonę, odrębną etykietę i odrębną kropkę — odczyt nie
 * zależy od rozróżnienia barw.
 *
 * Praca trwająca teraz dokłada `.dn-spinner` obok etykiety. Wskaźnik nie
 * zastępuje ani ikony, ani napisu.
 */
export function utworzPlakietkeStanu(stan: StanPary): PlakietkaStanu {
  const element = document.createElement('span');
  element.setAttribute('role', 'status');

  const kropka = document.createElement('span');
  const nosnik = document.createElement('span');
  nosnik.className = 'dn-okna__stan-ikona';
  const napis = document.createElement('span');

  const wskaznik = document.createElement('span');
  wskaznik.className = 'dn-spinner dn-okna__stan-wskaznik';
  wskaznik.setAttribute('aria-hidden', 'true');
  wskaznik.hidden = true;

  element.append(kropka, nosnik, napis, wskaznik);

  function ustaw(nowy: StanPary): void {
    const wyglad = wygladStanu(nowy);
    element.className = `dn-plakietka dn-okna__stan ${wyglad.wariantPlakietki}`;
    element.dataset.stanPary = nowy;
    kropka.className = `dn-kropka ${wyglad.wariantKropki}`;
    nosnik.replaceChildren(elementIkony(wyglad.ikona, { rozmiar: 16 }));
    napis.textContent = nazwaStanuPary(nowy);
    wskaznik.hidden = !wyglad.wTrakcie;
  }

  ustaw(stan);
  return { element, ustaw };
}
