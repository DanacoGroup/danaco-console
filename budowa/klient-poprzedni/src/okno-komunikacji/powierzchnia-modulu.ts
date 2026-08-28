import './powierzchnia-modulu.css';

import type { Action } from '../../../shared/contract';
import type { ZrodloAkcjiModulu } from './katalog-akcji';
import { utworzPanelAkcji } from './panel-akcji';
import { utworzPanelKontekstu, type KontekstOkna } from './panel-kontekstu';
import { utworzPasekNarzedzi } from './pasek-narzedzi-promptu';
import type { ProfilModulu } from './profil-modulu';
import { profilModulu } from './rejestr-profilow';
import { utworzWskaznikModulu } from './wskaznik-modulu';

/** Zależności powierzchni modułowej okna komunikacji, wstrzykiwane z zewnątrz jako kanał, katalog akcji i kontekst. */
export interface OpcjePowierzchni {
  /** Wstawia gotowe polecenie do pola wypowiedzi operatora. */
  naPolecenie(tekst: string): void;
  /** Źródło katalogu akcji modułu; brak = panel akcji nie powstaje. */
  katalogAkcji?: ZrodloAkcjiModulu;
  /** Kontekst okna pokazywany w panelu kontekstu. */
  kontekst?: KontekstOkna;
  /** Moduł, w którym okno zaczyna pracę; pusty znaczy „jeszcze nieustalony". */
  modul?: string;
}

/**
 * Powierzchnia modułowa okna komunikacji — wszystko, co okno przestawia przy zmianie modułu, z zachowaniem historii wątku.
 */
export interface PowierzchniaModulu {
  /** Element montowany w oknie, nad historią wątku. */
  element: HTMLElement;
  /** Przestawia okno na wskazany moduł; historii wątku nie dotyka. */
  ustawModul(kod: string): void;
  /** Kod modułu, w którym okno pracuje w tej chwili. */
  modul(): string;
  /** Odświeża panel kontekstu bez zmiany modułu, na przykład po dopisaniu wpisu. */
  odswiezKontekst(): void;
}

/**
 * Rekonfiguracja okna komunikacji przy zmianie modułu przestawia pasek narzędzi promptu, panel akcji i panel kontekstu, zachowując historię wątku; katalog akcji spóźniony o kolejną zmianę modułu jest odrzucany.
 */
export function utworzPowierzchnieModulu(opcje: OpcjePowierzchni): PowierzchniaModulu {
  const wskaznik = utworzWskaznikModulu();
  const pasek = utworzPasekNarzedzi((n) => opcje.naPolecenie(n.polecenie));
  const panelKontekstu = utworzPanelKontekstu();
  const panelAkcji =
    opcje.katalogAkcji === undefined
      ? null
      : utworzPanelAkcji((a) => opcje.naPolecenie(polecenieAkcji(a)));

  const element = document.createElement('section');
  element.className = 'dc-powierzchnia';
  element.setAttribute('aria-label', 'Powierzchnia modułowa okna komunikacji');

  const naglowek = document.createElement('div');
  naglowek.className = 'dc-powierzchnia__naglowek';
  naglowek.append(wskaznik.element, pasek.element);

  const panele = document.createElement('div');
  panele.className = 'dc-powierzchnia__panele';
  if (panelAkcji !== null) panele.append(panelAkcji.element);
  panele.append(panelKontekstu.element);

  element.append(naglowek, panele);

  const kontekst = opcje.kontekst ?? {};
  let biezacy = profilModulu(opcje.modul ?? '');

  function pokazAkcje(profil: ProfilModulu): void {
    if (panelAkcji === null || opcje.katalogAkcji === undefined) return;
    panelAkcji.ustaw(profil, [], 'Katalog akcji modułu — zapytanie w toku.');
    opcje.katalogAkcji(profil.kod, (akcje, powod) => {
      if (biezacy.kod !== profil.kod) return;
      panelAkcji.ustaw(profil, akcje, powod);
    });
  }

  function przestaw(profil: ProfilModulu): void {
    biezacy = profil;
    element.dataset['modul'] = profil.kod;
    wskaznik.ustaw(profil);
    pasek.ustaw(profil);
    panelKontekstu.ustaw(profil, kontekst);
    pokazAkcje(profil);
  }

  przestaw(biezacy);

  return {
    element,
    ustawModul(kod) {
      if (kod === biezacy.kod) return;
      przestaw(profilModulu(kod));
    },
    modul: () => biezacy.kod,
    odswiezKontekst: () => panelKontekstu.ustaw(biezacy, kontekst),
  };
}

/** Polecenie wykonania akcji katalogu wywoływane przez model po kliknięciu pozycji panelu akcji operatora. */
function polecenieAkcji(akcja: Action): string {
  return `Wykonaj akcję „${akcja.name}" (komenda kontraktu ${akcja.command}).`;
}
