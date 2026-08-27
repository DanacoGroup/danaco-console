import {
  ChangeKind,
  Command,
  EventType,
  type QueueChangedEvent,
  type SessionChangedEvent,
  type WindowChangedEvent,
} from '../../../shared/contract';
import type { Odsubskrybuj } from '../polaczenie/indeks';
import type { Kanal } from '../protokol/kanal';
import type { DanePulpitu } from './model-danych';
import { pustyStanZrodla } from './stan-zrodla';
import { zlozDanePulpitu } from './zlozenie-danych';

/**
 * Źródło danych pulpitu — jedyny dostawca kompletu `DanePulpitu`. Zbiera stan
 * z rdzenia odczytami kontraktu i po każdym zdarzeniu zmiany podaje odbiorcy
 * świeży komplet; przełożenie stanu na komplet należy do `zlozenie-danych.ts`.
 */
export interface ZrodloPulpitu {
  /** Wysyła odczyty i podpina subskrypcje; wywołanie jednokrotne. */
  uruchom(): void;
  /** Odpina subskrypcje zdarzeń. */
  zatrzymaj(): void;
}

/**
 * Buduje źródło danych pulpitu na kanale kontraktu, wiążąc odczyty stanu
 * i subskrypcje zdarzeń zmiany w jednego dostawcę kompletu dla pulpitu.
 */
export function utworzZrodloPulpitu(
  kanal: Kanal,
  odbiorca: (dane: DanePulpitu) => void,
): ZrodloPulpitu {
  const stan = pustyStanZrodla();
  const subskrypcje: Odsubskrybuj[] = [];

  const emituj = (): void => {
    odbiorca(zlozDanePulpitu(stan));
  };

  const przyjmijSesje = (zdarzenie: SessionChangedEvent): void => {
    stan.odczytano = true;
    if (zdarzenie.change === ChangeKind.Deleted) {
      stan.sesje.delete(zdarzenie.session.id);
      stan.obecnosc.delete(zdarzenie.session.id);
    } else {
      stan.sesje.set(zdarzenie.session.id, zdarzenie.session);
      if (zdarzenie.presence !== undefined) {
        stan.obecnosc.set(zdarzenie.session.id, zdarzenie.presence);
      }
    }
    emituj();
  };

  const przyjmijOkno = (zdarzenie: WindowChangedEvent): void => {
    stan.odczytano = true;
    if (zdarzenie.change === ChangeKind.Deleted) {
      stan.okna.delete(zdarzenie.window.id);
    } else {
      stan.okna.set(zdarzenie.window.id, zdarzenie.window);
    }
    emituj();
  };

  const przyjmijKolejke = (zdarzenie: QueueChangedEvent): void => {
    stan.odczytano = true;
    if (zdarzenie.change === ChangeKind.Deleted) {
      stan.kolejki.delete(zdarzenie.queue.id);
    } else {
      stan.kolejki.set(zdarzenie.queue.id, zdarzenie.queue);
    }
    emituj();
  };

  return {
    uruchom(): void {
      subskrypcje.push(
        kanal.naZdarzenie(EventType.SessionChanged, przyjmijSesje),
        kanal.naZdarzenie(EventType.WindowChanged, przyjmijOkno),
        kanal.naZdarzenie(EventType.QueueChanged, przyjmijKolejke),
        kanal.naZdarzenie(EventType.ProgressChanged, (zdarzenie) => {
          stan.odczytano = true;
          stan.procesy.set(zdarzenie.processId, zdarzenie);
          emituj();
        }),
      );

      kanal.wyslij(Command.SessionList, { includePresence: true }, (wynik) => {
        if (!wynik.udany || wynik.wynik === undefined) return;
        stan.odczytano = true;
        stan.sesje = new Map(wynik.wynik.sessions.map((sesja) => [sesja.id, sesja]));
        stan.obecnosc = new Map(
          (wynik.wynik.presence ?? []).map((obecnosc) => [obecnosc.sessionId, obecnosc]),
        );
        emituj();
      });

      kanal.wyslij(Command.WindowList, {}, (wynik) => {
        if (!wynik.udany || wynik.wynik === undefined) return;
        stan.odczytano = true;
        stan.okna = new Map(wynik.wynik.windows.map((okno) => [okno.id, okno]));
        emituj();
      });

      kanal.wyslij(Command.ChannelList, {}, (wynik) => {
        if (!wynik.udany || wynik.wynik === undefined) return;
        stan.odczytano = true;
        stan.kanaly = wynik.wynik.channels;
        emituj();
      });

      kanal.wyslij(Command.EnvironmentList, {}, (wynik) => {
        if (!wynik.udany || wynik.wynik === undefined) return;
        stan.odczytano = true;
        stan.srodowiska = new Map(
          wynik.wynik.environments.map((srodowisko) => [srodowisko.code, srodowisko]),
        );
        emituj();
      });
    },

    zatrzymaj(): void {
      for (const odepnij of subskrypcje) {
        odepnij();
      }
      subskrypcje.length = 0;
    },
  };
}
