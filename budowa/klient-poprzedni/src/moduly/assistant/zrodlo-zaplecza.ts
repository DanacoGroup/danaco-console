import {
  Command,
  ComponentKind,
  ConfigScope,
  WindowStatus,
  type Action,
  type Component,
  type SpeechAvailabilityGetResponse,
  type SpeechTranscribeRequest,
  type SpeechTranscribeResponse,
  type Window,
} from '../../../../shared/contract';
import type { Kanal, Wynik } from '../../protokol/kanal';
import { czyLogiczna, czyTablica, czyTekst, sprawdzKsztalt } from '../../protokol/ksztalt-odpowiedzi';
import { wywolaj } from '../../protokol/wywolanie';

/**
 * Zaplecze modułu Assistant — rejestry i czynności, z których moduł korzysta,
 * a których sam nie prowadzi.
 *
 *   `window.list`                — okna komunikacji sesji. `assistant.voice.command`
 *                                  wymaga pola `windowId`, a moduł okien nie
 *                                  zakłada: bierze to, które rdzeń przypisał
 *                                  modułowi.
 *   `action.list`                — katalog akcji zasięgu modułu. Siatka szybkich
 *                                  akcji Voice Console jest widokiem tego
 *                                  katalogu, nie listą zaszytą w kliencie: nowa
 *                                  akcja to nowy wiersz rdzenia, nie zmiana kodu.
 *   `component.list`             — komponenty własne rodzaju `assistant`, czyli
 *                                  Profile asystenta ze strefy 2 strony głównej.
 *                                  To jedyny odczyt kontraktu, którym da się
 *                                  wypełnić selektor profilu Voice Console.
 *   `speech.availability.get`    — czy silnik mowy stoi na tej maszynie. Brak
 *                                  silnika jest odpowiedzią, nie awarią, więc
 *                                  okno mówi o nim wprost.
 *   `speech.transcribe`          — rozpoznanie nagrania wskazanego ścieżką na
 *                                  maszynie silnika.
 */
export interface ZrodloZaplecza {
  okna(idSesji: string): Promise<Wynik<{ windows: Window[] }>>;
  akcje(kodModulu: string): Promise<Wynik<{ actions: Action[] }>>;
  /** Profile asystenta zapisane jako komponenty własne strefy 2. */
  profile(): Promise<Wynik<{ components: Component[] }>>;
  /** Stan silnika rozpoznawania mowy wraz z powodem niedostępności. */
  silnikMowy(): Promise<Wynik<SpeechAvailabilityGetResponse>>;
  /** Rozpoznanie nagrania leżącego na maszynie silnika. */
  transkrypcja(zadanie: SpeechTranscribeRequest): Promise<Wynik<SpeechTranscribeResponse>>;
}

export function utworzZrodloZaplecza(kanal: Kanal): ZrodloZaplecza {
  return {
    async okna(idSesji) {
      return sprawdzKsztalt(
        await wywolaj(kanal, Command.WindowList, {
          sessionId: idSesji,
          status: WindowStatus.Open,
        }),
        Command.WindowList,
        (tresc) => czyTablica(tresc.windows),
      );
    },

    async akcje(kodModulu) {
      return sprawdzKsztalt(
        await wywolaj(kanal, Command.ActionList, {
          scope: ConfigScope.Module,
          scopeId: kodModulu,
          enabledOnly: true,
        }),
        Command.ActionList,
        (tresc) => czyTablica(tresc.actions),
      );
    },

    async profile() {
      return sprawdzKsztalt(
        // Komponenty wyłączone wchodzą do wykazu: profil wyłączony w strefie 2
        // wciąż istnieje w rdzeniu i wskazanie go jest decyzją Operatora, a nie
        // pomyłką okna. Wyłączenie nazywa wiersz selektora.
        await wywolaj(kanal, Command.ComponentList, {
          kind: ComponentKind.Assistant,
          includeDisabled: true,
        }),
        Command.ComponentList,
        (tresc) => czyTablica(tresc.components),
      );
    },

    async silnikMowy() {
      return sprawdzKsztalt(
        await wywolaj(kanal, Command.SpeechAvailabilityGet, {}),
        Command.SpeechAvailabilityGet,
        (tresc) => czyLogiczna(tresc.available),
      );
    },

    async transkrypcja(zadanie) {
      return sprawdzKsztalt(
        await wywolaj(kanal, Command.SpeechTranscribe, zadanie),
        Command.SpeechTranscribe,
        // Pusta transkrypcja jest wynikiem prawidłowym, gdy `processed` jest
        // prawdą (cisza i szum), więc sprawdzian pyta o obecność pól, nie
        // o niepustą treść.
        (tresc) => czyLogiczna(tresc.processed) && czyTekst(tresc.transcript),
      );
    },
  };
}
