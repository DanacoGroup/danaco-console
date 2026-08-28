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

/**
 * Zależności widoku pulpitu: nazwa projektu na pasku aplikacji, kanał kontraktu
 * niosący zamiary do rdzenia oraz dwa przejścia — na inną trasę i do środowiska
 * wskazanego z matrycy sesji.
 */
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
 * Widok trasy Mission Control wiąże pulpit operacyjny z kanałem rdzenia
 * i z routerem: dane pulpitu pochodzą wyłącznie z odczytu rdzenia, a zamiary
 * Operatora wracają do rdzenia tym samym kanałem kontraktu.
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
