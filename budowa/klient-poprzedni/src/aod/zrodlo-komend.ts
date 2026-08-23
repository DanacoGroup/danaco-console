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
 * Źródło komend nakładki — jedyne miejsce w katalogu `aod/`, które zna nazwy
 * komend kontraktu i kształt ich żądań. Sekcje okna dostają stąd czynności,
 * nie `Command.*`, więc żadna sekcja nie powtarza opakowania kanału.
 *
 * Kontrakt niesie siedem komend `aod.*` i to źródło wywołuje wszystkie siedem.
 * `aod.voice.command` idzie bez mikrofonu: `audioRef` i `transcript` są w
 * kontrakcie polami opcjonalnymi, więc polecenie wydane samą treścią jest
 * wywołaniem pełnoprawnym. Nagrania nakładka nie wytwarza — granicę nazywa
 * `sekcja-glosu.ts`.
 */

/** Opakowuje `kanal.wyslij` w Promise — kanał sam daje wyłącznie wersję z wywołaniem zwrotnym. */
function poslijKomende<K extends Command>(
  kanal: Kanal,
  komenda: K,
  zadanie: RequestOf<K>,
): Promise<Wynik<ResponseOf<K>>> {
  return new Promise((rozstrzygnij) => {
    kanal.wyslij(komenda, zadanie, (wynik) => rozstrzygnij(wynik));
  });
}

/** Czynności nakładki widziane przez sekcje okna. */
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
 * Buduje źródło czynności nad kanałem rdzenia.
 *
 * Żadne żądanie nie niesie `deviceId`. Rdzeń rozumie przez identyfikator
 * urządzenia numer wiersza katalogu maszyn — `adapterNakladkiAod.urzadzenieNakladki`
 * czyta go przez `strconv.ParseInt`, tym samym prawem co rodzina `accessPoint.*`.
 * Identyfikator sesji („ses_…") dostaje odmowę `validation_failed`, a
 * `kolumna-aod.ts` przerywa wczytywanie na pierwszym błędzie stanu, więc jedna zła
 * wartość zabiera całą treść nakładki.
 *
 * Pole puste rdzeń przyjmuje. Żadna komenda kontraktu nie mówi klientowi,
 * którym numerem katalogu maszyn jest ta maszyna, więc podstawienie zgadniętej
 * jedynki byłoby wartością zmyśloną.
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
