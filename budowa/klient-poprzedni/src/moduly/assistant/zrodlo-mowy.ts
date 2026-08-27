import {
  Command,
  EventType,
  type SpeechAudioFetchResponse,
  type SpeechAudioUploadRequest,
  type SpeechAudioUploadResponse,
  type SpeechListenPartialEvent,
  type SpeechListenStartRequest,
  type SpeechListenStartResponse,
  type SpeechListenStopResponse,
  type SpeechWakeDetectedEvent,
  type SpeechWakeGetResponse,
  type SpeechWakeSetRequest,
  type SpeechWakeSetResponse,
} from '../../../../shared/contract';
import type { Odsubskrybuj } from '../../polaczenie/magistrala-zdarzen';
import type { Kanal, Wynik } from '../../protokol/kanal';
import { czyLiczba, czyLogiczna, czyObiekt, czyTekst, sprawdzKsztalt } from '../../protokol/ksztalt-odpowiedzi';
import { wywolaj } from '../../protokol/wywolanie';

/** Sześć komend mowy dobudowanych obok transkrypcji oraz dwa zdarzenia nasłuchu ciągłego. */

/** Nagranie przyjmowane przez rdzeń niesie bajty wraz z typem treści, na przykład audio/webm albo audio/wav. */
export interface NagranieDoWyslania {
  /** Bajty nagrania w postaci base64, bez przedrostka schematu danych. */
  base64: string;
  /** Typ treści, na przykład `audio/webm` albo `audio/wav`. */
  typTresci: string;
  /** Okno, z którego nagranie wyszło; puste zostawia je bez przypisania. */
  idOkna?: string;
  /** Karta sesji, której nagranie dotyczy. */
  idSesji?: string;
  /** Czy nagranie ma przeżyć transkrypcję. Puste bierze nastawę rdzenia. */
  zachowaj?: boolean;
}

/** Nastawa wybudzania zapisywana wybiórczo, niosąca frazę, tryb, próg i odszumianie; pola pominięte zostają bez zmian. */
export interface NastawaWybudzania {
  fraza?: string;
  tryb?: SpeechWakeSetRequest['mode'];
  progDetekcji?: number;
  odszumianie?: boolean;
}

export interface ZrodloMowy {
  /** `speech.audio.upload` — bajty nagrania w zamian za odnośnik. */
  przeslijNagranie(nagranie: NagranieDoWyslania): Promise<Wynik<SpeechAudioUploadResponse>>;
  /** `speech.audio.fetch` — odsłuch nagrania albo odpowiedzi syntezowanej. */
  pobierzNagranie(odnosnik: string): Promise<Wynik<SpeechAudioFetchResponse>>;
  /** `speech.wake.get` — nastawa wybudzania wraz z uczciwym stanem wykonalności. */
  nastawaWybudzania(): Promise<Wynik<SpeechWakeGetResponse>>;
  /** `speech.wake.set` — zapis frazy, trybu, progu i odszumiania. */
  zapiszWybudzanie(nastawa: NastawaWybudzania): Promise<Wynik<SpeechWakeSetResponse>>;
  /** `speech.listen.start` — nasłuch ciągły rdzenia na wskazanym oknie. */
  uruchomNasluch(idOkna: string, idSesji: string): Promise<Wynik<SpeechListenStartResponse>>;
  /** `speech.listen.stop` — zatrzymanie nasłuchu. */
  zatrzymajNasluch(idNasluchu: string, idOkna: string): Promise<Wynik<SpeechListenStopResponse>>;
  /** Subskrypcja `speech.listen.partial` — rozpoznanego odcinka nasłuchu. */
  naOdcinek(sluchacz: (tresc: SpeechListenPartialEvent) => void): Odsubskrybuj;
  /** Subskrypcja `speech.wake.detected` — wykrycia frazy wybudzającej. */
  naWybudzenie(sluchacz: (tresc: SpeechWakeDetectedEvent) => void): Odsubskrybuj;
}

export function utworzZrodloMowy(kanal: Kanal): ZrodloMowy {
  return {
    async przeslijNagranie(nagranie) {
      const zadanie: SpeechAudioUploadRequest = {
        audio: nagranie.base64,
        contentType: nagranie.typTresci,
      };
      if (nagranie.idOkna !== undefined && nagranie.idOkna !== '') {
        zadanie.windowId = nagranie.idOkna;
      }
      if (nagranie.idSesji !== undefined && nagranie.idSesji !== '') {
        zadanie.sessionId = nagranie.idSesji;
      }
      if (nagranie.zachowaj !== undefined) zadanie.retain = nagranie.zachowaj;
      return sprawdzKsztalt(
        await wywolaj(kanal, Command.SpeechAudioUpload, zadanie),
        Command.SpeechAudioUpload,
        (tresc) => czyTekst(tresc.audioRef) && czyLiczba(tresc.sizeBytes),
      );
    },

    async pobierzNagranie(odnosnik) {
      return sprawdzKsztalt(
        await wywolaj(kanal, Command.SpeechAudioFetch, { audioRef: odnosnik }),
        Command.SpeechAudioFetch,
        (tresc) => czyTekst(tresc.audio) && czyTekst(tresc.contentType),
      );
    },

    async nastawaWybudzania() {
      return sprawdzKsztalt(
        await wywolaj(kanal, Command.SpeechWakeGet, {}),
        Command.SpeechWakeGet,
        (tresc) => czyObiekt(tresc.config) && czyLogiczna(tresc.available),
      );
    },

    async zapiszWybudzanie(nastawa) {
      const zadanie: SpeechWakeSetRequest = {};
      if (nastawa.fraza !== undefined) zadanie.phrase = nastawa.fraza;
      if (nastawa.tryb !== undefined) zadanie.mode = nastawa.tryb;
      if (nastawa.progDetekcji !== undefined) zadanie.vadThreshold = nastawa.progDetekcji;
      if (nastawa.odszumianie !== undefined) zadanie.noiseSuppression = nastawa.odszumianie;
      return sprawdzKsztalt(
        await wywolaj(kanal, Command.SpeechWakeSet, zadanie),
        Command.SpeechWakeSet,
        (tresc) => czyObiekt(tresc.config),
      );
    },

    async uruchomNasluch(idOkna, idSesji) {
      const zadanie: SpeechListenStartRequest = { windowId: idOkna };
      if (idSesji !== '') zadanie.sessionId = idSesji;
      return sprawdzKsztalt(
        await wywolaj(kanal, Command.SpeechListenStart, zadanie),
        Command.SpeechListenStart,
        (tresc) => czyLogiczna(tresc.listening) && czyTekst(tresc.listenerId),
      );
    },

    async zatrzymajNasluch(idNasluchu, idOkna) {
      const zadanie: Record<string, string> = {};
      if (idNasluchu !== '') zadanie['listenerId'] = idNasluchu;
      if (idOkna !== '') zadanie['windowId'] = idOkna;
      return sprawdzKsztalt(
        await wywolaj(kanal, Command.SpeechListenStop, zadanie),
        Command.SpeechListenStop,
        (tresc) => czyLogiczna(tresc.stopped),
      );
    },

    naOdcinek(sluchacz) {
      return kanal.naZdarzenie(EventType.SpeechListenPartial, (tresc) => sluchacz(tresc));
    },

    naWybudzenie(sluchacz) {
      return kanal.naZdarzenie(EventType.SpeechWakeDetected, (tresc) => sluchacz(tresc));
    },
  };
}
