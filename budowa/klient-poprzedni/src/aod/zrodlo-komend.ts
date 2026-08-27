/**
 * Źródło komend nakładki, będące jedynym miejscem katalogu nakładki, które zna
 * nazwy komend kontraktu i kształt ich żądań. Kontrakt niesie siedem komend
 * rodziny nakładki i to źródło wywołuje wszystkie siedem.
 */
import {
  Command,
  type AodChatSendRequest,
  type AodChatSendResponse,
  type AodContextGetRequest,
  type AodContextGetResponse,
  type AodObserveAttachResponse,
  type AodObserveDetachResponse,
  type AodStatusGetResponse,
  type AodSuggestionRequest,
  type AodSuggestionResponse,
  type AodVoiceCommandRequest,
  type AodVoiceCommandResponse,
  type RequestOf,
  type ResponseOf,
} from '../../../shared/contract';
import type { Kanal, Wynik } from '../protokol/kanal';

/**
 * Opakowuje wysłanie komendy kanałem w obietnicę, ponieważ kanał sam daje
 * wyłącznie wersję z wywołaniem zwrotnym, a sekcje okna czekają na wynik
 * składnią asynchroniczną.
 */
function poslijKomende<K extends Command>(
  kanal: Kanal,
  komenda: K,
  zadanie: RequestOf<K>,
): Promise<Wynik<ResponseOf<K>>> {
  return new Promise((rozstrzygnij) => {
    kanal.wyslij(komenda, zadanie, (wynik) => rozstrzygnij(wynik));
  });
}

/**
 * Czynności nakładki widziane przez sekcje okna. Sekcje wołają czynności,
 * a nie stałe komend, dzięki czemu żadna z nich nie powtarza opakowania kanału
 * ani nie zna nazw komend kontraktu.
 */
export interface ZrodloAod {
  /** Identyfikator sesji bieżącej albo `undefined`, gdy sesji jeszcze nie ma. */
  idSesji(): string | undefined;
  /** `aod.status.get` — stan nakładki wraz z wykazem procesów przypiętych. */
  stan(): Promise<Wynik<AodStatusGetResponse>>;
  /** `aod.context.get` — komplet kontekstu okna ogniskowanego. */
  kontekst(zadanie: AodContextGetRequest): Promise<Wynik<AodContextGetResponse>>;
  /** `aod.suggestion` — podpowiedzi następnego kroku, od bytu najwęższego do platformy. */
  podpowiedzi(zadanie: AodSuggestionRequest): Promise<Wynik<AodSuggestionResponse>>;
  /** `aod.observe.attach` — przypina proces do obserwacji; zwraca wykaz po zmianie. */
  przypnij(processId: string): Promise<Wynik<AodObserveAttachResponse>>;
  /** `aod.observe.detach` — odpina proces od obserwacji; zwraca wykaz po zmianie. */
  odepnij(processId: string): Promise<Wynik<AodObserveDetachResponse>>;
  /** `aod.chat.send` — wiadomość z nakładki do okna rozmowy. */
  wyslijRozmowe(zadanie: AodChatSendRequest): Promise<Wynik<AodChatSendResponse>>;
  /** `aod.voice.command` — polecenie kierowane do asystenta. */
  wydajPolecenieGlosowe(
    zadanie: AodVoiceCommandRequest,
  ): Promise<Wynik<AodVoiceCommandResponse>>;
}

/**
 * Buduje źródło czynności nad kanałem rdzenia. Żadne żądanie nie niesie
 * identyfikatora urządzenia, ponieważ żadna komenda kontraktu nie mówi
 * klientowi, którym numerem katalogu maszyn jest ta maszyna.
 */
export function utworzZrodloAod(kanal: Kanal): ZrodloAod {
  const idSesji = (): string | undefined => kanal.sesja().id() || undefined;

  return {
    idSesji,

    stan: () => poslijKomende(kanal, Command.AodStatusGet, {}),

    kontekst: (zadanie) => poslijKomende(kanal, Command.AodContextGet, zadanie),

    podpowiedzi: (zadanie) => poslijKomende(kanal, Command.AodSuggestion, zadanie),

    przypnij: (processId) => poslijKomende(kanal, Command.AodObserveAttach, { processId }),

    odepnij: (processId) => poslijKomende(kanal, Command.AodObserveDetach, { processId }),

    wyslijRozmowe: (zadanie) => poslijKomende(kanal, Command.AodChatSend, zadanie),

    wydajPolecenieGlosowe: (zadanie) => poslijKomende(kanal, Command.AodVoiceCommand, zadanie),
  };
}
