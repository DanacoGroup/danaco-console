import {
  Command,
  type AutomationScheduleSetRequest,
  type AutomationScheduleSetResponse,
  type AutomationWorkflowListResponse,
  type AutomationWorkflowSaveRequest,
  type AutomationWorkflowSaveResponse,
  type ScheduleGetResponse,
  type SessionToolAttachRequest,
  type SessionToolAttachResponse,
  type SessionToolDetachRequest,
  type SessionToolDetachResponse,
  type SessionToolListResponse,
  type ToolsCatalogListRequest,
  type ToolsCatalogListResponse,
} from '../../../../shared/contract';
import type { Kanal, Wynik } from '../../protokol/kanal';
import { czyObiekt, czyTablica, sprawdzKsztalt } from '../../protokol/ksztalt-odpowiedzi';
import { wywolaj } from '../../protokol/wywolanie';

/**
 * Narzędzia, dołożenia i rutyny widziane przez okno Command & Tools Hub.
 *
 * Plik odpowiada wyłącznie za warstwę wywołań kontraktu wraz ze sprawdzianem
 * kształtu odpowiedzi; wzorem jest `zrodlo-assistant.ts`. Żadne wywołanie nie
 * rzuca wyjątkiem — odmowa wraca polem `blad` wyniku.
 *
 * Rodziny są trzy i każda odpowiada za inne pytanie okna:
 *
 *   `tools.catalog.list`   co w ogóle da się wywołać po ukośniku — narzędzia,
 *                          umiejętności i komendy akcji w jednym wykazie;
 *   `session.tool.*`       co z tego jest dołożone do tej karty sesji. Dołożenie
 *                          żyje w stanie sesji i NIE rusza definicji eksperta
 *                          (rozstrzygnięcie Właściciela zapisane w opisie
 *                          komendy);
 *   `automation.*`         gdzie kończy makro i rutyna asystenta. Opracowanie
 *                          modułu prowadzi tę drogę wprost: „Przekazanie
 *                          powtarzalnego makra lub rutyny do modułu
 *                          Automations". Rdzeń nie ma innego magazynu sekwencji
 *                          kroków, więc jest to jedyne miejsce, w którym makro
 *                          przeżywa sesję.
 *
 * `schedule.get` stoi po stronie odczytu do pary z `automation.schedule.set` —
 * bez niego okno pokazywałoby harmonogram, który samo wysłało, zamiast tego,
 * który rdzeń trzyma.
 */
export interface ZrodloNarzedzi {
  /** `tools.catalog.list` — wykaz pozycji po ukośniku. */
  katalog(zadanie: ToolsCatalogListRequest): Promise<Wynik<ToolsCatalogListResponse>>;
  /** `session.tool.list` — dołożenia karty sesji. */
  dolozenia(idSesji: string): Promise<Wynik<SessionToolListResponse>>;
  /** `session.tool.attach` — dołożenie pozycji do karty sesji. */
  dolozy(zadanie: SessionToolAttachRequest): Promise<Wynik<SessionToolAttachResponse>>;
  /** `session.tool.detach` — zdjęcie dołożenia z karty sesji. */
  zdejmij(zadanie: SessionToolDetachRequest): Promise<Wynik<SessionToolDetachResponse>>;
  /** `automation.workflow.list` — automatyki dostępne Operatorowi. */
  automatyki(): Promise<Wynik<AutomationWorkflowListResponse>>;
  /** `automation.workflow.save` — zapis makra jako automatyki. */
  zapiszAutomatyke(
    zadanie: AutomationWorkflowSaveRequest,
  ): Promise<Wynik<AutomationWorkflowSaveResponse>>;
  /** `schedule.get` — harmonogramy automatyk. */
  harmonogramy(): Promise<Wynik<ScheduleGetResponse>>;
  /** `automation.schedule.set` — cykliczność rutyny. */
  ustawHarmonogram(
    zadanie: AutomationScheduleSetRequest,
  ): Promise<Wynik<AutomationScheduleSetResponse>>;
}

export function utworzZrodloNarzedzi(kanal: Kanal): ZrodloNarzedzi {
  return {
    async katalog(zadanie) {
      return sprawdzKsztalt(
        await wywolaj(kanal, Command.ToolsCatalogList, zadanie),
        Command.ToolsCatalogList,
        (tresc) => czyTablica(tresc.entries) && czyTablica(tresc.groups),
      );
    },

    async dolozenia(idSesji) {
      return sprawdzKsztalt(
        await wywolaj(kanal, Command.SessionToolList, { sessionId: idSesji }),
        Command.SessionToolList,
        (tresc) => czyTablica(tresc.tools),
      );
    },

    async dolozy(zadanie) {
      return sprawdzKsztalt(
        await wywolaj(kanal, Command.SessionToolAttach, zadanie),
        Command.SessionToolAttach,
        (tresc) => czyObiekt(tresc.tool) && czyTablica(tresc.tools),
      );
    },

    async zdejmij(zadanie) {
      return sprawdzKsztalt(
        await wywolaj(kanal, Command.SessionToolDetach, zadanie),
        Command.SessionToolDetach,
        (tresc) => czyTablica(tresc.tools),
      );
    },

    async automatyki() {
      return sprawdzKsztalt(
        await wywolaj(kanal, Command.AutomationWorkflowList, {}),
        Command.AutomationWorkflowList,
        (tresc) => czyTablica(tresc.workflows),
      );
    },

    async zapiszAutomatyke(zadanie) {
      return sprawdzKsztalt(
        await wywolaj(kanal, Command.AutomationWorkflowSave, zadanie),
        Command.AutomationWorkflowSave,
        (tresc) => czyObiekt(tresc.workflow),
      );
    },

    async harmonogramy() {
      return sprawdzKsztalt(
        await wywolaj(kanal, Command.ScheduleGet, {}),
        Command.ScheduleGet,
        (tresc) => czyTablica(tresc.schedules),
      );
    },

    async ustawHarmonogram(zadanie) {
      return sprawdzKsztalt(
        await wywolaj(kanal, Command.AutomationScheduleSet, zadanie),
        Command.AutomationScheduleSet,
        (tresc) => czyObiekt(tresc.schedule),
      );
    },
  };
}
