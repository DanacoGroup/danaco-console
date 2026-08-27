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

/**
 * Zależności obsługi zamiarów pulpitu: kanał, którym idą komendy kontraktu, oraz
 * wejście przenoszące widok na środowisko wskazane przez pulpit. Obie warstwa
 * powłoki podaje przy montażu, bo pulpit sam żadnej z nich nie tworzy.
 */
export interface ObslugaZamiarow {
  /** Kanał komunikatów kontraktu. */
  kanal: Kanal;
  /** Wejście do środowiska wskazanego przez pulpit. */
  naSrodowisko(kod: IdSrodowiska): void;
}

/**
 * Skutki zamiarów zgłaszanych przez pulpit: przełożenie zamiaru na komendę
 * kontraktu albo na komunikat nazywający brak drogi. Pulpit nadaje zamiar
 * nazwami kontraktu, a kopertę komendy składa dopiero ta warstwa.
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

/**
 * Przyczyna niepowodzenia komendy podana przez rdzeń: wiadomość rdzenia
 * uzupełniona kodem błędu, gdy rdzeń kod dołączył. Odpowiedź bez wiadomości
 * daje zdanie mówiące wprost, że przyczyny nie podano.
 */
function opisBledu(wiadomosc: string | undefined, kod: string | undefined): string {
  if (wiadomosc === undefined) return 'Rdzeń nie podał przyczyny.';
  return kod === undefined ? wiadomosc : `${wiadomosc} (${kod})`;
}
