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
 * a których sam nie prowadzi: okna komunikacji sesji, katalog akcji zasięgu
 * modułu, komponenty profilu asystenta oraz stan i wywołanie silnika mowy.
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
        // Komponenty wyłączone wchodzą do wykazu; wyłączenie nazywa wiersz selektora.
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
        // Pusta transkrypcja przy `processed` prawdziwym jest wynikiem prawidłowym.
        (tresc) => czyLogiczna(tresc.processed) && czyTekst(tresc.transcript),
      );
    },
  };
}
