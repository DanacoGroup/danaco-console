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

/** Miejsca powłoki, w które montuje się widok sterowania. */
export interface MiejscaWidoku {
  /** Kolumna sterowania obok sceny okna (`KorzenAplikacji.panel`). */
  panel: HTMLElement;
  /** Miejsce akcji paska górnego (`KorzenAplikacji.akcje`). */
  akcjePaska: HTMLElement;
}

/** Zależności widoku: droga do rdzenia, okno oraz wspólny wykaz kanałów. */
export interface ZaleznosciWidoku extends MiejscaWidoku {
  kanal: Kanal;
  okno: Window;
  rejestrKanalow: RejestrKanalow;
}

/** Widok sterowania jednego okna komunikacji. */
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

/**
 * Oprawa czyniąca komplet sterowania widocznym.
 *
 * Kontrolki nie powstają tutaj — pochodzą w całości z katalogu `sterowanie/`
 * (`utworzPanelSterowania`: środowisko wykonania, host, moduł, model, model
 * zapasowy, nakład rozumowania, tryb uprawnień, rola okna, katalogi robocze).
 * Ten plik dokłada im miejsce, nagłówek, podsumowanie wartości i dwa uchwyty
 * rozwijania.
 *
 * Kolumna montuje się obok sceny okna, a nie w prawym rogu paska górnego —
 * dziewięć kontrolek nie ma jak się tam zmieścić.
 *
 * Widok powstaje osobno dla każdego okna i domyka się na jego identyfikatorze,
 * więc dwa okna obok siebie dostają dwie niezależne kolumny bez ani jednej
 * wspólnej zmiennej. Szuflada startuje rozwinięta, żeby komplet ustawień był
 * widoczny bez szukania, co nacisnąć.
 */
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
