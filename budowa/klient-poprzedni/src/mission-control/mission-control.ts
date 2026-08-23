import './mission-control.css';

import { utworzKolumneKolejek } from './kolumna-kolejek';
import { utworzKolumneProcesow } from './kolumna-procesow';
import { utworzKolumneZespolu } from './kolumna-zespolu';
import type { DanePulpitu } from './model-danych';
import { utworzNaglowekPulpitu } from './naglowek-pulpitu';
import { utworzNadajnik, type Odpiecie } from './nadajnik';
import {
  opiszWejscie,
  opiszZamiarDecyzji,
  opiszZamiarKolejki,
  opiszZamiarUtworzenia,
} from './opis-zamiaru';
import { utworzPasDecyzji } from './pas-decyzji';
import { utworzPasRelacji } from './pas-relacji';
import { utworzPasekDzialan } from './pasek-dzialan';
import { utworzSekcjeAktywnosci } from './sekcja-aktywnosci';
import { utworzSekcjeMatrycy } from './sekcja-matrycy';
import { utworzSekcjeOperacji } from './sekcja-operacji';
import { utworzSekcjeUtworz } from './sekcja-utworz';
import type {
  WejscieDoSesji,
  ZamiarDecyzji,
  ZamiarKolejki,
  ZamiarUtworzenia,
} from './zdarzenia-pulpitu';

/**
 * Pulpit operacyjny Mission Control — złożenie sekcji.
 *
 * Jedna odpowiedzialność: kompozycja. Plik nie rysuje treści — buduje nadajniki,
 * powołuje sekcje w ustalonej kolejności i spina odświeżanie:
 *   • Aktywność AI z pasem decyzji — eskalacja koordynatora,
 *   • Operacje AI — kanały modelu i koszt,
 *   • Matryca sesji z pasem relacji — procesy biegnące równolegle,
 *   • Utwórz — wejście do pracy jeszcze nierozpoczętej,
 *   • trzy kolumny — procesy, kolejki, zespół.
 */
export interface MissionControl {
  /** Element montowany w powłoce. */
  element: HTMLElement;
  /** Podpina odbiorcę wejścia do sesji; zwraca odpięcie. */
  naWejscieDoSesji(odbiorca: (wejscie: WejscieDoSesji) => void): Odpiecie;
  /** Podpina odbiorcę działań na kolejkach ról. */
  naZamiarKolejki(odbiorca: (zamiar: ZamiarKolejki) => void): Odpiecie;
  /** Podpina odbiorcę żądań utworzenia bytu. */
  naZamiarUtworzenia(odbiorca: (zamiar: ZamiarUtworzenia) => void): Odpiecie;
  /** Podpina odbiorcę wezwania do rozstrzygnięcia wstrzymanych przepływów. */
  naZamiarDecyzji(odbiorca: (zamiar: ZamiarDecyzji) => void): Odpiecie;
  /** Podmienia komplet danych bez przebudowy widoku. */
  odswiez(dane: DanePulpitu): void;
  /** Odpina wszystkich odbiorców. */
  rozlacz(): void;
}

/**
 * Buduje pulpit operacyjny.
 *
 * Komplet danych jest wymagany i pochodzi ze źródła pulpitu (`zrodlo-pulpitu`),
 * które buduje go wyłącznie z odczytów i zdarzeń rdzenia. Komplet sprzed
 * pierwszego odczytu daje `pustyKomplet()` — stan oczekiwania bez wartości
 * zastępczych.
 */
export function utworzMissionControl(dane: DanePulpitu): MissionControl {
  const wejscia = utworzNadajnik<WejscieDoSesji>('wejście do sesji');
  const kolejki = utworzNadajnik<ZamiarKolejki>('zamiar kolejki');
  const utworzenia = utworzNadajnik<ZamiarUtworzenia>('zamiar utworzenia');
  const decyzje = utworzNadajnik<ZamiarDecyzji>('zamiar decyzji');

  const czolo = utworzNaglowekPulpitu(dane.zrodloDanych);
  const dzialania = utworzPasekDzialan();

  const aktywnosc = utworzSekcjeAktywnosci(dane.aktywnosc);
  const pasDecyzji = utworzPasDecyzji(dane.decyzje, decyzje.nadaj);
  aktywnosc.podKaflami.append(pasDecyzji.element);

  const operacje = utworzSekcjeOperacji(dane.kanaly);

  const matryca = utworzSekcjeMatrycy(dane.matryca, dane.pozaSrodowiskami, wejscia.nadaj);
  const relacje = utworzPasRelacji(dane.relacje);
  matryca.podKolumnami.append(relacje.element);

  const utworz = utworzSekcjeUtworz(utworzenia.nadaj);

  const procesy = utworzKolumneProcesow(dane.procesy);
  const kolejkiWidok = utworzKolumneKolejek(dane.kolejki, kolejki.nadaj);
  const zespol = utworzKolumneZespolu(dane.zespol);

  const trojkolumna = document.createElement('div');
  trojkolumna.className = 'mc-trojkolumna';
  trojkolumna.append(procesy.element, kolejkiWidok.element, zespol.element);

  const element = document.createElement('div');
  element.className = 'mc-pulpit';
  element.setAttribute('aria-labelledby', 'mc-tytul-pulpitu');
  element.append(
    czolo.element,
    aktywnosc.element,
    operacje.element,
    matryca.element,
    utworz.element,
    trojkolumna,
    dzialania.element,
  );

  // Pasek działań słucha własnych nadajników, żeby każde działanie pulpitu
  // miało widoczny skutek niezależnie od odbiorców zewnętrznych.
  const wlasne: Odpiecie[] = [
    wejscia.sluchaj((wejscie) => dzialania.zapisz(opiszWejscie(wejscie))),
    kolejki.sluchaj((zamiar) => dzialania.zapisz(opiszZamiarKolejki(zamiar))),
    utworzenia.sluchaj((zamiar) => dzialania.zapisz(opiszZamiarUtworzenia(zamiar))),
    decyzje.sluchaj((zamiar) => dzialania.zapisz(opiszZamiarDecyzji(zamiar))),
  ];

  return {
    element,
    naWejscieDoSesji: wejscia.sluchaj,
    naZamiarKolejki: kolejki.sluchaj,
    naZamiarUtworzenia: utworzenia.sluchaj,
    naZamiarDecyzji: decyzje.sluchaj,
    odswiez(nowe) {
      czolo.oznacz(nowe.zrodloDanych);
      aktywnosc.odswiez(nowe.aktywnosc);
      pasDecyzji.odswiez(nowe.decyzje);
      operacje.odswiez(nowe.kanaly);
      matryca.odswiez(nowe.matryca, nowe.pozaSrodowiskami);
      relacje.odswiez(nowe.relacje);
      procesy.odswiez(nowe.procesy);
      kolejkiWidok.odswiez(nowe.kolejki);
      zespol.odswiez(nowe.zespol);
    },
    rozlacz() {
      for (const odepnij of wlasne) {
        odepnij();
      }
      wejscia.rozlacz();
      kolejki.rozlacz();
      utworzenia.rozlacz();
      decyzje.rozlacz();
    },
  };
}
