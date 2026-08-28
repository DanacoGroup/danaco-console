import {
  Command,
  EventType,
  type AssistantActionChangedEvent,
  type AssistantActionControl,
  type AssistantActionStatusRequest,
  type AssistantActionStatusResponse,
  type AssistantActivityListRequest,
  type AssistantActivityFlagRequest,
  type AssistantActivityFlagResponse,
  type AssistantActivityListResponse,
  type AssistantVoiceCommandRequest,
  type AssistantVoiceCommandResponse,
  type Envelope,
} from '../../../../shared/contract';
import type { Odsubskrybuj } from '../../polaczenie/magistrala-zdarzen';
import type { Kanal, Wynik } from '../../protokol/kanal';
import { czyObiekt, czyTablica, czyTekst, sprawdzKsztalt } from '../../protokol/ksztalt-odpowiedzi';
import { wywolaj } from '../../protokol/wywolanie';

/** Trzy komendy i jedno zdarzenie obszaru assistant widziane przez okna modułu. */

/** Polecenie wydane asystentowi z Voice Console, niosące okno modułu, transkrypcję, profil i nastawę syntezy mowy. */
export interface PolecenieAsystenta {
  /** Okno modułu; kontrakt wymaga go w każdym poleceniu. */
  idOkna: string;
  /** Transkrypcja po edycji w pasku promptu — treść polecenia. */
  transkrypcja: string;
  /** Profil asystenta; pusty znaczy profil domyślny rdzenia. */
  profil: string;
  /** Czy odpowiedź ma zostać odczytana syntezą mowy. */
  czytaj: boolean;
}

/** Sterowanie zleceniem wydane z Actions Monitora, niosące identyfikator zlecenia i nowy priorytet albo jego brak. */
export interface SterowanieZleceniem {
  idZlecenia: string;
  sterowanie: AssistantActionControl;
  /** Nowy priorytet; pominięty zostawia priorytet bez zmiany. */
  priorytet?: number;
}

export interface ZrodloAssistant {
  /** `assistant.voice.command` — polecenie przez potok STT → model → TTS. */
  polecenie(zlecenie: PolecenieAsystenta): Promise<Wynik<AssistantVoiceCommandResponse>>;
  /** `assistant.action.status` — sam odczyt stanu zleceń. */
  zlecenia(idOkna: string): Promise<Wynik<AssistantActionStatusResponse>>;
  /** `assistant.action.status` — sterowanie zleceniem i priorytetem. */
  steruj(zlecenie: SterowanieZleceniem): Promise<Wynik<AssistantActionStatusResponse>>;
  /** `assistant.activity.list` — dziennik działań, opcjonalnie jednego zlecenia. */
  dziennik(idOkna: string, idZlecenia: string): Promise<Wynik<AssistantActivityListResponse>>;
  /** Wyróżnienie wpisu dziennika wraz z powodem; znacznik przeżywa odświeżenie, zdjęcie kasuje też powód. */
  oznaczWpis(
    idWpisu: string,
    wazny: boolean,
    powod: string,
  ): Promise<Wynik<AssistantActivityFlagResponse>>;
  /** Subskrypcja jedynego zdarzenia obszaru; słuchacz dostaje kopertę, bo treść nie niesie sesji. */
  naZmianeZlecenia(
    sluchacz: (tresc: AssistantActionChangedEvent, koperta: Envelope) => void,
  ): Odsubskrybuj;
}

/** Górna granica dziennika — okno pokazuje historię, nie cały zapis rdzenia, ograniczony do dwustu wpisów. */
const GRANICA_DZIENNIKA = 200;

export function utworzZrodloAssistant(kanal: Kanal): ZrodloAssistant {
  return {
    async polecenie(zlecenie) {
      const zadanie: AssistantVoiceCommandRequest = {
        windowId: zlecenie.idOkna,
        transcript: zlecenie.transkrypcja,
        speak: zlecenie.czytaj,
      };
      if (zlecenie.profil.trim() !== '') zadanie.profileId = zlecenie.profil.trim();
      return sprawdzKsztalt(
        await wywolaj(kanal, Command.AssistantVoiceCommand, zadanie),
        Command.AssistantVoiceCommand,
        (tresc) => czyTekst(tresc.transcript) && czyObiekt(tresc.action),
      );
    },

    async zlecenia(idOkna) {
      const zadanie: AssistantActionStatusRequest = {};
      if (idOkna !== '') zadanie.windowId = idOkna;
      return sprawdzKsztalt(
        await wywolaj(kanal, Command.AssistantActionStatus, zadanie),
        Command.AssistantActionStatus,
        (tresc) => czyTablica(tresc.actions),
      );
    },

    async steruj(zlecenie) {
      const zadanie: AssistantActionStatusRequest = {
        actionId: zlecenie.idZlecenia,
        control: zlecenie.sterowanie,
      };
      if (zlecenie.priorytet !== undefined) zadanie.priority = zlecenie.priorytet;
      return sprawdzKsztalt(
        await wywolaj(kanal, Command.AssistantActionStatus, zadanie),
        Command.AssistantActionStatus,
        (tresc) => czyTablica(tresc.actions),
      );
    },

    async dziennik(idOkna, idZlecenia) {
      const zadanie: AssistantActivityListRequest = { limit: GRANICA_DZIENNIKA };
      if (idOkna !== '') zadanie.windowId = idOkna;
      if (idZlecenia !== '') zadanie.actionId = idZlecenia;
      return sprawdzKsztalt(
        await wywolaj(kanal, Command.AssistantActivityList, zadanie),
        Command.AssistantActivityList,
        (tresc) => czyTablica(tresc.entries),
      );
    },

    async oznaczWpis(idWpisu, wazny, powod) {
      const zadanie: AssistantActivityFlagRequest = { entryId: idWpisu, important: wazny };
      if (powod.trim() !== '') zadanie.note = powod.trim();
      return sprawdzKsztalt(
        await wywolaj(kanal, Command.AssistantActivityFlag, zadanie),
        Command.AssistantActivityFlag,
        (tresc) => czyObiekt(tresc.entry),
      );
    },

    naZmianeZlecenia(sluchacz) {
      return kanal.naZdarzenie(EventType.AssistantActionChanged, (tresc, koperta) =>
        sluchacz(tresc, koperta),
      );
    },
  };
}
