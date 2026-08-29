/** Rozstrzygnięcie i wykonanie chwili przekazania sterowania ramie po drodze wejścia. */

import type { Environment } from '../../../shared/contract.ts';
import type { Kanal } from '../protokol/kanal.ts';
import type { StanPrzebiegu } from '../wejscie/przebieg.ts';
import { zamontujRame } from './montaz.ts';

/** Stan przebiegu, w którym środowisko jest już znane — warunek konieczny wykonania przekazania. */
type StanZeSrodowiskiem = StanPrzebiegu & { srodowisko: Environment };

/** Czy przebieg doszedł do miejsca, w którym rama aplikacji przejmuje ekran po drodze wejścia. */
export function gotowaDoPrzekazania(stan: StanPrzebiegu): stan is StanZeSrodowiskiem {
  return stan.etap === 'przygotowanie' && stan.srodowisko !== undefined;
}

/** Węzły dokumentu i działanie, których przekazanie potrzebuje od wywołującego. */
export interface ZaleznosciPrzekazania {
  dokument: Document;
  zdejmijOknoWejscia: () => void;
  /** Kanał, którym rama woła komendy rdzenia dla okna modułu; brak — rama bez okna modułu. */
  kanal?: Kanal;
}

/**
 * Montuje ramę, a dopiero po powodzeniu zdejmuje scenę wejścia — nieudany
 * montaż zostaje odmową nazwaną w dzienniku, scena wejścia zostaje na ekranie.
 */
export function wykonajPrzekazanie(stan: StanZeSrodowiskiem, zaleznosci: ZaleznosciPrzekazania): boolean {
  const scenaWejscia = zaleznosci.dokument.querySelector('[data-wejscie]');
  const miejsceRamy = zaleznosci.dokument.querySelector('[data-rama-aplikacji]');
  if (scenaWejscia === null || miejsceRamy === null) {
    console.error('[rama] przekazanie odrzucone: brak węzła montażu w dokumencie');
    return false;
  }

  try {
    zamontujRame({
      miejsce: miejsceRamy as HTMLElement,
      srodowisko: stan.srodowisko,
      moduly: stan.moduly,
      sesje: stan.sesje,
      kanal: zaleznosci.kanal,
      /* Konto, którym Operator wszedł — pasek stanu ma nazwać, czyim kontem
         pracuje. Rdzeń tego nie oddaje, więc idzie z drogi wejścia. */
      kontoOperatora: stan.kontoOperatora,
    });
  } catch (blad) {
    console.error('[rama] przekazanie odrzucone: montaż rzucił wyjątkiem', blad);
    return false;
  }

  zaleznosci.zdejmijOknoWejscia();
  (scenaWejscia as HTMLElement).hidden = true;
  (miejsceRamy as HTMLElement).hidden = false;
  return true;
}

/** Przekazanie jednorazowe: znacznik ustawia się dopiero po powodzeniu, więc odmowa nie zamyka drogi na kolejną zmianę stanu. */
export function utworzPrzekazanieJednorazowe(
  zaleznosci: ZaleznosciPrzekazania,
): (stan: StanPrzebiegu) => void {
  let przekazano = false;
  return function naZmianePrzebiegu(stan: StanPrzebiegu): void {
    if (przekazano || !gotowaDoPrzekazania(stan)) return;
    przekazano = wykonajPrzekazanie(stan, zaleznosci);
  };
}
