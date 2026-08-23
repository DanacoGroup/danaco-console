import {
  Command,
  EventType,
  type Envelope,
  type Message,
  type MessageChangedEvent,
  type MessageListRequest,
  type MessageSendRequest,
  type MessageStopRequest,
  type MessageStopResponse,
  type MonitorStatusRequest,
  type MonitorStatusResponse,
  type MonitorSubscribeRequest,
  type MonitorSubscribeResponse,
  type Queue,
  type QueueActionRequest,
  type QueueChangedEvent,
  type QueueCreateRequest,
  type QueueListRequest,
  type StreamChunkEvent,
  type WindowHandoffRequest,
  type WindowHandoffResponse,
  type WindowStateChangedEvent,
} from '../../../../shared/contract';
import type { Odsubskrybuj } from '../../polaczenie/magistrala-zdarzen';
import type { Kanal, Wynik } from '../../protokol/kanal';
import { czyObiekt, czyTablica, sprawdzKsztalt } from '../../protokol/ksztalt-odpowiedzi';
import { przenies } from '../../protokol/wynik-czastkowy';
import { wywolaj } from '../../protokol/wywolanie';

/**
 * Bieg pracy widziany przez klienta — ta połowa źródła, która porusza rolami.
 *
 * Koordynator widzi pełny strumień wykonawcy. Nośnikiem jest `stream.chunk` —
 * jeden strumień całej platformy, w którym `windowId` mówi, czyja to tura.
 * Okno koordynatora nie zakłada drugiego kanału podglądu: subskrybuje ten sam
 * strumień i odsiewa okna swoich wykonawców.
 *
 * Kolejka jest jedna: `queue.create` i `queue.action` to ten sam silnik,
 * którym jedzie pętla sesyjna i moduł Automations — moduł ról nie buduje
 * własnego wykonawcy zleceń.
 */
export interface ZrodloBiegu {
  /** `message.send` — wydanie polecenia oknu roli. */
  wyslij(zadanie: MessageSendRequest): Promise<Wynik<Message>>;
  /**
   * `window.handoff` — utrwalenie zlecenia koordynator→wykonawca.
   *
   * Doręczenie i zapis to dwie różne czynności i obie są potrzebne.
   * `message.send` niesie treść do okna wykonawcy, więc wykonawca rusza do
   * pracy; sama ta droga zostawia jednak więź wyłącznie w pamięci przeglądarki
   * i gubi ją z jej zamknięciem. `window.handoff` zapisuje w bazie rdzenia,
   * kto komu co zlecił, i zakłada pozycję kolejki — dzięki temu para przeżywa
   * restart, a Mission Control ma co pokazać.
   */
  przekaz(zadanie: Omit<WindowHandoffRequest, 'sessionId'>): Promise<Wynik<WindowHandoffResponse>>;
  /** `message.stop` — przerwanie tury okna; przycisk czynny zawsze. */
  zatrzymaj(zadanie: MessageStopRequest): Promise<Wynik<MessageStopResponse>>;
  /** `message.list` — historia okna; stąd bierze się wynik wykonawcy. */
  wiadomosci(zadanie: MessageListRequest): Promise<Wynik<Message[]>>;
  /** `queue.create` — etap planu jako kolejka jednego silnika. */
  zalozKolejke(zadanie: QueueCreateRequest): Promise<Wynik<Queue>>;
  /** `queue.action` — sterowanie kolejką: start, pauza, wznowienie, stop, powtórzenie. */
  sterujKolejka(zadanie: QueueActionRequest): Promise<Wynik<Queue>>;
  /**
   * `queue.list` — kolejki jednego silnika, zawężone sesją albo oknem.
   *
   * Tędy idzie zatrzymanie podagenta. Kontrakt nie ma `subagent.stop`,
   * a `queue.action` żąda `queueId`, którego `Subagent` nie niesie. Ogniwem
   * jest nazwa kolejki: rdzeń zakłada kolejkę podagenta pod jego
   * identyfikatorem (`core/adapter_modul_orkiestracja.go`), więc wykaz kolejek
   * sesji pozwala odnaleźć `queueId` po `name === subagent.id` — bez komendy
   * wymyślonej poza kontraktem.
   */
  wykazKolejek(zadanie: QueueListRequest): Promise<Wynik<Queue[]>>;
  /**
   * `monitor.status` — zbiorczy stan bieżący procesów telemetrii, bez zakładania
   * obserwacji. Rdzeń oddaje `{ statuses }`; wykaz jest pusty, gdy sito nic nie
   * łapie.
   */
  stanMonitora(zadanie: MonitorStatusRequest): Promise<Wynik<MonitorStatusResponse>>;
  /**
   * `monitor.subscribe` — zapisuje wskazane okno na telemetrię postępu procesów
   * i oddaje `{ statuses, subscribed }`. Żądanie z `windowId` zakłada
   * obserwację (`subscribed` prawdziwe), bez niego jest zwykłym odczytem.
   */
  subskrybujMonitor(zadanie: MonitorSubscribeRequest): Promise<Wynik<MonitorSubscribeResponse>>;
  /** Subskrypcja `stream.chunk` — strumień wykonawcy widziany przez koordynatora. */
  naFragment(sluchacz: (tresc: StreamChunkEvent, koperta: Envelope) => void): Odsubskrybuj;
  /** Subskrypcja `message.changed` — domknięta odpowiedź okna roli. */
  naWiadomosc(sluchacz: (tresc: MessageChangedEvent) => void): Odsubskrybuj;
  /** Subskrypcja `window.state.changed` — zmiana stanu okna roli. */
  naStanOkna(sluchacz: (tresc: WindowStateChangedEvent) => void): Odsubskrybuj;
  /** Subskrypcja `queue.changed` — stan kolejki etapu na żywo. */
  naKolejke(sluchacz: (tresc: QueueChangedEvent) => void): Odsubskrybuj;
}

export function utworzZrodloBiegu(kanal: Kanal): ZrodloBiegu {
  return {
    async przekaz(zadanie) {
      // Sesję zna kanał, nie sterowanie modułu — podawanie jej z widoku byłoby
      // przepisywaniem wartości, którą warstwa protokołu i tak trzyma.
      return sprawdzKsztalt(
        await wywolaj(kanal, Command.WindowHandoff, { ...zadanie, sessionId: kanal.sesja().id() }),
        Command.WindowHandoff,
        (tresc) => czyObiekt(tresc.window),
      );
    },

    async wyslij(zadanie) {
      const wynik = sprawdzKsztalt(
        await wywolaj(kanal, Command.MessageSend, zadanie),
        Command.MessageSend,
        (tresc) => czyObiekt(tresc.message),
      );
      return przenies(wynik, (tresc) => tresc.message);
    },

    async zatrzymaj(zadanie) {
      return sprawdzKsztalt(
        await wywolaj(kanal, Command.MessageStop, zadanie),
        Command.MessageStop,
        (tresc) => typeof tresc.stopped === 'boolean',
      );
    },

    async wiadomosci(zadanie) {
      const wynik = sprawdzKsztalt(
        await wywolaj(kanal, Command.MessageList, zadanie),
        Command.MessageList,
        (tresc) => czyTablica(tresc.messages),
      );
      return przenies(wynik, (tresc) => tresc.messages);
    },

    async zalozKolejke(zadanie) {
      const wynik = sprawdzKsztalt(
        await wywolaj(kanal, Command.QueueCreate, zadanie),
        Command.QueueCreate,
        (tresc) => czyObiekt(tresc.queue),
      );
      return przenies(wynik, (tresc) => tresc.queue);
    },

    async sterujKolejka(zadanie) {
      const wynik = sprawdzKsztalt(
        await wywolaj(kanal, Command.QueueAction, zadanie),
        Command.QueueAction,
        (tresc) => czyObiekt(tresc.queue),
      );
      return przenies(wynik, (tresc) => tresc.queue);
    },

    async wykazKolejek(zadanie) {
      const wynik = sprawdzKsztalt(
        await wywolaj(kanal, Command.QueueList, zadanie),
        Command.QueueList,
        (tresc) => czyTablica(tresc.queues),
      );
      return przenies(wynik, (tresc) => tresc.queues);
    },

    async stanMonitora(zadanie) {
      return sprawdzKsztalt(
        await wywolaj(kanal, Command.MonitorStatus, zadanie),
        Command.MonitorStatus,
        (tresc) => czyTablica(tresc.statuses),
      );
    },

    async subskrybujMonitor(zadanie) {
      return sprawdzKsztalt(
        await wywolaj(kanal, Command.MonitorSubscribe, zadanie),
        Command.MonitorSubscribe,
        (tresc) => czyTablica(tresc.statuses) && typeof tresc.subscribed === 'boolean',
      );
    },

    naFragment: (sluchacz) => kanal.naZdarzenie(EventType.StreamChunk, sluchacz),
    naWiadomosc: (sluchacz) => kanal.naZdarzenie(EventType.MessageChanged, sluchacz),
    naStanOkna: (sluchacz) => kanal.naZdarzenie(EventType.WindowStateChanged, sluchacz),
    naKolejke: (sluchacz) => kanal.naZdarzenie(EventType.QueueChanged, sluchacz),
  };
}
