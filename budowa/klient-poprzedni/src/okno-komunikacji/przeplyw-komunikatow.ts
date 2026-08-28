import {
  Command,
  EventType,
  MessageRole,
  type ErrorInfo,
} from '../../../shared/contract';
import type { Transport } from '../polaczenie/gniazdo';
import type { Kanal } from '../protokol/kanal';
import type { PostepUzgodnienia, Uzgodnienie } from '../protokol/uzgodnienie';
import { utworzDyktowanie } from './dyktowanie/dyktowanie';
import { nazwaEtapu } from './etykiety-okna';
import type { OknoKomunikacji } from './okno';
import type { OpisOkna } from './opis-okna';
import { nazwaNadawcy, wpisOperatora, wpisSystemowy } from './wpis';
import { personaRoli, pokazFragment } from './wpisy-strumienia';

/** Elementy, które przepływ komunikatów spina w jedną całość: transport, kontrakt, warstwę dyktowania oraz widok okna. */
export interface CzesciPrzeplywu {
  okno: OknoKomunikacji;
  kanal: Kanal;
  transport: Transport;
  uzgodnienie: Uzgodnienie;
  opis: OpisOkna;
}

/**
 * Przepływ komunikatów okna łączy transport i kontrakt po jednej stronie z widokiem po drugiej; treść wpisana przed otwarciem okna czeka i wysyła się po uzgodnieniu, bez blokad interfejsu.
 */
export function polaczPrzeplyw(czesci: CzesciPrzeplywu): void {
  const { okno, kanal, transport, uzgodnienie, opis } = czesci;

  // Dyktowanie powstaje w przepływie, bo tu jest kanał — przepływ składa warstwę i oddaje ją oknu.
  okno.podlaczDyktowanie(utworzDyktowanie(kanal));
  const oczekujace: string[] = [];
  const zestrumieniowane = new Set<string>();
  const persona = opis.kanalModelu;

  /** Identyfikator okna nadany przez rdzeń; pusty do chwili uzgodnienia. */
  function idOkna(): string {
    return uzgodnienie.okno()?.id ?? '';
  }

  /** Znacznik tury jest przybliżeniem: rozstrzyga tylko, czy poprzedzić wysłanie zatrzymaniem tury. */
  let turaWBiegu = false;

  function nadajTresc(tresc: string): void {
    kanal.wyslij(
      Command.MessageSend,
      { windowId: idOkna(), content: tresc, stream: true },
      (wynik) => {
        if (wynik.udany) {
          turaWBiegu = true;
          return;
        }
        turaWBiegu = false;
        okno.dopisz(wpisSystemowy(`Wiadomość niewysłana — ${opisBledu(wynik.blad)}`));
      },
    );
  }

  /** Przerwanie jest jawne: wysyłka do zajętego okna odmawia konfliktem zamiast skasować odpowiedź. */
  function wyslijTresc(tresc: string): void {
    if (!turaWBiegu) {
      nadajTresc(tresc);
      return;
    }
    kanal.wyslij(Command.MessageStop, { windowId: idOkna() }, () => {
      turaWBiegu = false;
      nadajTresc(tresc);
    });
  }

  // Ponowne połączenie nie zakłada drugiego okna, jeżeli rdzeń otworzył już wcześniej to samo okno.
  transport.naStan((stan) => {
    okno.pokazStan(stan, transport.oczekujace());
    if (stan === 'polaczony' && uzgodnienie.okno() === null) uzgodnienie.rozpocznij();
  });

  uzgodnienie.naPostep((postep) => okno.dopisz(wpisSystemowy(opisPostepu(postep))));

  uzgodnienie.naOtwarcieOkna(() => {
    for (const tresc of oczekujace.splice(0)) wyslijTresc(tresc);
  });

  okno.naWpisanie((tresc) => {
    okno.dopisz(wpisOperatora(tresc, nazwaNadawcy(MessageRole.User)));
    if (idOkna().length === 0) {
      oczekujace.push(tresc);
      return;
    }
    wyslijTresc(tresc);
  });

  kanal.naZdarzenie(EventType.StreamChunk, (fragment, koperta) => {
    if (fragment.windowId !== idOkna()) return;
    zestrumieniowane.add(fragment.messageId);
    // Zdarzenie domykające gasi znacznik tury raz na turę, także po błędzie i po przerwaniu strumienia.
    if (koperta.done === true) turaWBiegu = false;
    pokazFragment(okno, persona, fragment);
  });

  kanal.naZdarzenie(EventType.MessageChanged, ({ message }) => {
    if (message.windowId !== idOkna()) return;
    if (message.role === MessageRole.User || zestrumieniowane.has(message.id)) return;
    okno.dopisz({
      nadawca: message.role,
      persona: personaRoli(message.role, persona),
      tresc: message.content,
      znacznikCzasu: message.createdAt,
    });
  });

  // Zmiana okna niesie też jego moduł: okno rekonfiguruje pasek narzędzi, panel akcji i kontekst.
  kanal.naZdarzenie(EventType.WindowChanged, ({ change, window }) => {
    if (window.id !== idOkna()) return;
    okno.ustawModul(window.moduleId);
    okno.dopisz(wpisSystemowy(`Okno komunikacji — ${change}, stan: ${window.status}`));
  });
}

/** Komunikat o etapie uzgodnienia z rdzeniem, pokazywany operatorowi w trakcie zakładania okna komunikacji. */
function opisPostepu(postep: PostepUzgodnienia): string {
  const etap = nazwaEtapu(postep.etap);
  return postep.udany ? etap : `${etap} — ${opisBledu(postep.blad)}`;
}

/** Opis błędu kontraktu przedstawiany operatorowi, gdy uzgodnienie okna z rdzeniem całkiem się nie powiedzie. */
function opisBledu(blad: ErrorInfo | undefined): string {
  if (blad === undefined) return 'rdzeń nie podał przyczyny';
  return `${blad.message} (${blad.code})`;
}
