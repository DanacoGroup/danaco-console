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

/**
 * Trzy komendy i jedno zdarzenie obszaru `assistant.*` widziane przez okna
 * modułu.
 *
 * Plik odpowiada wyłącznie za warstwę wywołań kontraktu wraz ze sprawdzianem
 * kształtu odpowiedzi. Źródło nie ma własnego stanu i nie buduje ani jednego
 * elementu — stan zleceń mieszka w `stan-assistant.ts`, żeby trzy okna modułu
 * patrzyły na jeden zbiór, a nie na trzy kopie.
 *
 * Żadne wywołanie nie rzuca wyjątkiem: niepowodzenie wraca polem `blad` wyniku,
 * a okno pokazuje je w swoim stanie błędu. Tą samą drogą wraca odmowa
 * merytoryczna rdzenia i koperta `assistant.unknown` drogi bez uchwytu, więc
 * okno nazywa je wprost zamiast udawać wykonanie.
 */

/** Polecenie wydane asystentowi z Voice Console. */
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

/** Sterowanie zleceniem wydane z Actions Monitora. */
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
  /**
   * `assistant.activity.flag` — wyróżnienie wpisu dziennika wraz z powodem.
   *
   * Wyróżnienie ma gdzie zamieszkać po stronie rdzenia, więc okno go nie udaje:
   * znacznik przeżywa odświeżenie wykazu, a zdjęcie wyróżnienia kasuje też
   * powód — powód bez znacznika byłby notatką do wpisu, którego nikt nie
   * wyróżnił.
   */
  oznaczWpis(
    idWpisu: string,
    wazny: boolean,
    powod: string,
  ): Promise<Wynik<AssistantActivityFlagResponse>>;
  /**
   * Subskrypcja `assistant.action.changed` — jedynego zdarzenia obszaru.
   *
   * Słuchacz dostaje także kopertę, bo treść zdarzenia nie niesie sesji, a
   * koperta ją niesie (`shared/contract.ts`, pole `sessionId`; rdzeń wypełnia
   * je sesją okna zlecenia — `adapter_modul_asystent_wykonawca.go`,
   * `okno.IdSesji`). Bez niej zlecenie założone w innym oknie tej samej sesji
   * byłoby nie do odróżnienia od zlecenia cudzej sesji i trzeba by odrzucać oba.
   */
  naZmianeZlecenia(
    sluchacz: (tresc: AssistantActionChangedEvent, koperta: Envelope) => void,
  ): Odsubskrybuj;
}

/** Górna granica dziennika — okno pokazuje historię, nie cały zapis rdzenia. */
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
