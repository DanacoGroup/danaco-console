import { wykonajZamiarUstawien } from './akcje-ustawien';
import {
  pustyKomplet,
  utworzMissionControl,
  utworzZrodloPulpitu,
  type IdSrodowiska,
} from '../mission-control/indeks';
import type { Kanal } from '../protokol/kanal';
import { utworzPasekAplikacji } from './pasek-aplikacji';
import type { WidokTrasy } from './router';
import { Trasa } from './trasy';
import {
  obsluzWejscieDoSesji,
  obsluzZamiarDecyzji,
  obsluzZamiarKolejki,
  obsluzZamiarUtworzenia,
} from './zamiary-pulpitu';

/** Zależności widoku pulpitu. */
export interface ZaleznosciPulpitu {
  /** Nazwa projektu na pasku aplikacji. */
  projekt: string;
  /** Kanał komunikatów kontraktu — zamiary pulpitu idą do rdzenia. */
  kanal: Kanal;
  /** Przejście na inną trasę. */
  naTrase(trasa: Trasa): void;
  /** Wejście do środowiska wskazanego z matrycy sesji. */
  naSrodowisko(kod: IdSrodowiska): void;
}

/**
 * Widok trasy Mission Control.
 *
 * Jedna odpowiedzialność: związanie pulpitu operacyjnego z kanałem rdzenia
 * i z routerem. Pulpit stoi obok strony głównej, a nie zamiast niej: Centrum
 * dowodzenia jest wejściem przy rozpoczynaniu pracy, pulpit — przy powrocie
 * do niej. Obydwa są dostępne z przełącznika widoków.
 *
 * Dane pulpitu pochodzą z rdzenia: źródło pulpitu odpytuje `session.list`,
 * `window.list` i `channel.list` oraz subskrybuje zdarzenia `session.changed`,
 * `window.changed`, `queue.changed` i `progress.changed`; do pierwszego odczytu
 * pulpit pokazuje stany puste. Zamiary idą do rdzenia tym samym kanałem —
 * `session.open` i `queue.action` przechodzą kanałem kontraktu.
 */
export function utworzWidokPulpitu(zaleznosci: ZaleznosciPulpitu): WidokTrasy {
  const { kanal, naSrodowisko } = zaleznosci;

  const pulpit = utworzMissionControl(pustyKomplet());
  const zrodlo = utworzZrodloPulpitu(kanal, (dane) => pulpit.odswiez(dane));
  zrodlo.uruchom();

  const pasek = utworzPasekAplikacji({
    naUstawienie: (pozycja) => wykonajZamiarUstawien(pozycja, zaleznosci.kanal),
    projekt: zaleznosci.projekt,
    naTrase: zaleznosci.naTrase,
  });

  const element = document.createElement('div');
  element.className = 'dn-widok dn-widok--pulpit';
  element.append(pasek.element, pulpit.element);

  const obsluga = { kanal, naSrodowisko };

  pulpit.naWejscieDoSesji((wejscie) => obsluzWejscieDoSesji(obsluga, wejscie));
  pulpit.naZamiarKolejki((zamiar) => obsluzZamiarKolejki(obsluga, zamiar));
  pulpit.naZamiarUtworzenia(obsluzZamiarUtworzenia);
  pulpit.naZamiarDecyzji(obsluzZamiarDecyzji);

  return {
    element,
    przyWejsciu: () => pasek.trasy.ustawBiezaca(Trasa.Pulpit),
  };
}
