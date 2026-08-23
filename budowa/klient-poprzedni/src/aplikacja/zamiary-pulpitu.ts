import { Command } from '../../../shared/contract';
import type {
  IdSrodowiska,
  WejscieDoSesji,
  ZamiarDecyzji,
  ZamiarKolejki,
  ZamiarUtworzenia,
} from '../mission-control/indeks';
import type { Kanal } from '../protokol/kanal';
import { pokazKomunikat } from './komunikaty';

/** Zależności obsługi zamiarów pulpitu. */
export interface ObslugaZamiarow {
  /** Kanał komunikatów kontraktu. */
  kanal: Kanal;
  /** Wejście do środowiska wskazanego przez pulpit. */
  naSrodowisko(kod: IdSrodowiska): void;
}

/**
 * Skutki zamiarów zgłaszanych przez Mission Control.
 *
 * Jedna odpowiedzialność: przełożenie zamiaru na komendę kontraktu albo na
 * uczciwą odpowiedź, gdy komendy dla niego nie ma. Pulpit sam niczego nie
 * wysyła — nadaje zamiar wyrażony nazwami kontraktu, a kopertę składa ta
 * warstwa.
 *
 * Dwa zamiary nie mają dziś drogi przez kontrakt i żaden z nich nie jest
 * udawany:
 *
 *   priorytet     — brak odpowiednika w `QueueAction` i wśród komend;
 *   przekazanie   — `context.transfer` wymaga kompletu `ContextBundle`,
 *                   którego zamiar kolejki nie niesie i którego nie wolno
 *                   zmyślić po stronie klienta.
 *
 * Operator dostaje w obu przypadkach komunikat mówiący wprost, czego brakuje.
 */
export function obsluzWejscieDoSesji(obsluga: ObslugaZamiarow, wejscie: WejscieDoSesji): void {
  obsluga.kanal.wyslij(Command.SessionOpen, { sessionId: wejscie.sessionId }, (wynik) => {
    if (wynik.udany) return;
    pokazKomunikat({
      tytul: `Sesja „${wejscie.tytul}" nieotwarta`,
      tresc: opisBledu(wynik.blad?.message, wynik.blad?.code),
      waga: 'blad',
    });
  });
  obsluga.naSrodowisko(wejscie.srodowisko);
}

export function obsluzZamiarKolejki(obsluga: ObslugaZamiarow, zamiar: ZamiarKolejki): void {
  if (zamiar.rodzaj === 'kolejka') {
    obsluga.kanal.wyslij(
      Command.QueueAction,
      { queueId: zamiar.idKolejki, action: zamiar.dzialanie },
      (wynik) => {
        if (wynik.udany) return;
        pokazKomunikat({
          tytul: `Kolejka ${zamiar.rola} — działanie ${zamiar.dzialanie} nieudane`,
          tresc: opisBledu(wynik.blad?.message, wynik.blad?.code),
          waga: 'blad',
        });
      },
    );
    return;
  }

  const tresc =
    zamiar.rodzaj === 'przekazanie'
      ? `Komenda ${Command.ContextTransfer} wymaga kompletu kontekstu (ContextBundle), którego pas kolejek nie niesie. Przekazanie idzie z okna koordynatora, nie z pulpitu.`
      : 'Podniesienie priorytetu nie ma dziś odpowiednika w kontrakcie — ani w QueueAction, ani wśród komend.';

  pokazKomunikat({
    tytul: `Kolejka ${zamiar.rola}`,
    tresc,
    waga: 'ostrz',
  });
}

export function obsluzZamiarUtworzenia(zamiar: ZamiarUtworzenia): void {
  pokazKomunikat({
    tytul: zamiar.etykieta,
    tresc: `Utworzenie bytu „${zamiar.rodzaj}" prowadzi przez kreator, którego jeszcze nie ma. Sesję zakładasz przyciskiem ＋ w pasie kart środowiska.`,
    waga: 'info',
  });
}

export function obsluzZamiarDecyzji(zamiar: ZamiarDecyzji): void {
  pokazKomunikat({
    tytul: 'Wstrzymane przepływy czekają na decyzję',
    tresc: `Czeka ${zamiar.przeplywy}; najdłużej: ${zamiar.najstarszy}. Rozstrzygnięcie należy do Operatora i zapada w oknie koordynatora.`,
    waga: 'ostrz',
  });
}

/** Przyczyna niepowodzenia komendy podana przez rdzeń. */
function opisBledu(wiadomosc: string | undefined, kod: string | undefined): string {
  if (wiadomosc === undefined) return 'Rdzeń nie podał przyczyny.';
  return kod === undefined ? wiadomosc : `${wiadomosc} (${kod})`;
}
