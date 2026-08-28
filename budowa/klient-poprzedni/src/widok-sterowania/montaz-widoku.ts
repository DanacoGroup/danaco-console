import './widok-sterowania.css';

import type { Window } from '../../../shared/contract';
import type { Kanal } from '../protokol/kanal';
import {
  utworzPanelSterowania,
  type PanelSterowania,
  type RejestrKanalow,
} from '../sterowanie/indeks';
import { utworzNaglowekWidoku } from './naglowek-widoku';
import { utworzObserwatorUstawien } from './obserwator-ustawien';
import { utworzPodsumowanie } from './podsumowanie-ustawien';
import { utworzSzuflade } from './szuflada';
import { utworzUchwytPaska } from './uchwyt-paska';

/** Definiuje miejsca powłoki aplikacji, w które montuje się widok sterowania: kolumnę panelu obok sceny okna oraz miejsce akcji na pasku górnym. */
export interface MiejscaWidoku {
  /** Kolumna sterowania obok sceny okna (`KorzenAplikacji.panel`). */
  panel: HTMLElement;
  /** Miejsce akcji paska górnego (`KorzenAplikacji.akcje`). */
  akcjePaska: HTMLElement;
}

/** Zależności widoku sterowania: droga do rdzenia, okno, którego widok dotyczy, wspólny wykaz kanałów oraz miejsca powłoki, w które widok się montuje. */
export interface ZaleznosciWidoku extends MiejscaWidoku {
  kanal: Kanal;
  okno: Window;
  rejestrKanalow: RejestrKanalow;
}

/** Widok sterowania jednego okna komunikacji, udostępniający jego element, identyfikator okna oraz czynności rozwijania, przyjęcia okna i odłączenia. */
export interface WidokSterowania {
  /** Element kolumny sterowania. */
  element: HTMLElement;
  /** Okno, którego dotyczy widok. */
  idOkna(): string;
  /** Rozwija komplet kontrolek. */
  otworz(): void;
  /** Przyjmuje okno potwierdzone przez rdzeń poza widokiem. */
  przyjmijOkno(okno: Window): void;
  /** Odłącza subskrypcje i zdejmuje widok z powłoki. */
  rozlacz(): void;
}

/** Widok sterowania buduje oprawę dla kontrolek katalogu sterowania: montuje kolumnę obok sceny okna, nagłówek, podsumowanie wartości oraz dwa uchwyty rozwijania, osobno dla każdego okna. */
export function zamontujWidokSterowania(
  zaleznosci: ZaleznosciWidoku,
): WidokSterowania {
  const { kanal, okno, rejestrKanalow, panel, akcjePaska } = zaleznosci;

  const komplet: PanelSterowania = utworzPanelSterowania({
    kanal,
    okno,
    rejestrKanalow,
  });
  const obserwator = utworzObserwatorUstawien(kanal, okno);

  const tresc = document.createElement('div');
  tresc.className = 'dc-widok-ster__tresc';
  tresc.tabIndex = -1;
  tresc.append(komplet.element);

  const podsumowanie = utworzPodsumowanie(obserwator, rejestrKanalow, () =>
    rozwinIWskaz(),
  );

  const szuflada = utworzSzuflade({
    podsumowanie: podsumowanie.element,
    tresc,
    idOkna: komplet.idOkna(),
    otwartaNaStarcie: true,
  });

  const naglowek = utworzNaglowekWidoku(obserwator, komplet.idOkna());
  const uchwytPaska = utworzUchwytPaska(szuflada);

  const element = document.createElement('div');
  element.className = 'dc-widok-ster';
  element.dataset.okno = komplet.idOkna();
  element.append(naglowek.element, szuflada.element);

  panel.append(element);
  akcjePaska.append(uchwytPaska.element);

  /** Rozwija szufladę i prowadzi ognisko do kompletu kontrolek. */
  function rozwinIWskaz(): void {
    szuflada.ustaw(true);
    tresc.scrollIntoView({ block: 'nearest' });
    tresc.focus();
  }

  return {
    element,
    idOkna: komplet.idOkna,
    otworz: () => szuflada.ustaw(true),
    przyjmijOkno: komplet.przyjmijOkno,
    rozlacz() {
      uchwytPaska.rozlacz();
      naglowek.rozlacz();
      podsumowanie.rozlacz();
      obserwator.rozlacz();
      komplet.rozlacz();
      element.remove();
      uchwytPaska.element.remove();
    },
  };
}
