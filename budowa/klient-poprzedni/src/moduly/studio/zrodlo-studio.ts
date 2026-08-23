import {
  Command,
  EventType,
  type StudioContextualOpRequest,
  type StudioContextualOpResponse,
  type StudioDiffCompareRequest,
  type StudioDiffCompareResponse,
  type StudioDocumentChangedEvent,
  type StudioDocumentOpenRequest,
  type StudioDocumentOpenResponse,
  type StudioDocumentSaveRequest,
  type StudioDocumentSaveResponse,
  type StudioRepositoryListRequest,
  type StudioRepositoryListResponse,
  type StudioRepositoryRestoreRequest,
  type StudioRepositoryRestoreResponse,
} from '../../../../shared/contract';
import type { Odsubskrybuj } from '../../polaczenie/magistrala-zdarzen';
import type { Kanal, Wynik } from '../../protokol/kanal';
import { czyObiekt, czyTablica, sprawdzKsztalt } from '../../protokol/ksztalt-odpowiedzi';
import { wywolajUczciwie } from './odmowa-rdzenia';

/**
 * Sześć komend obszaru `studio.*` widzianych przez okna modułu.
 *
 * Źródło nie ma własnego stanu i nie buduje ani jednego elementu — jest
 * warstwą wywołań i sprawdzianu kształtu odpowiedzi. Dokument czynny mieszka
 * w `stan-studio.ts`, żeby pięć okien patrzyło na jeden dokument, a nie na
 * pięć jego kopii.
 *
 * Wszystkie sześć komend ma uchwyt w rdzeniu, który wpina port Studia
 * (`core/kompozycja.go` → `zarejestrujStudio`). `document.open` zakłada albo
 * wczytuje dokument, `document.save` zapisuje treść i zakłada wersję,
 * `repository.list` oddaje wykaz wersji, `repository.restore` przestawia treść
 * dokumentu na wskazaną wersję i zostawia wersje nowsze, `diff.compare` liczy
 * fragmenty różnicy i trafienia wzorca, `contextual.op` wychodzi do kanału
 * modelu okna. Zapis i przywrócenie rozgłaszają `studio.document.changed`.
 *
 * Osłona odmowy zostaje mimo to: `wywolajUczciwie` broni przed kopertą
 * `studio.unknown`, która nie niesie pola `status` i nigdy by się nie
 * skorelowała, przez co okno stałoby w ładowaniu bez końca. Kosztuje jedną
 * subskrypcję.
 */
export interface ZrodloStudio {
  otworz(zadanie: StudioDocumentOpenRequest): Promise<Wynik<StudioDocumentOpenResponse>>;
  zapisz(zadanie: StudioDocumentSaveRequest): Promise<Wynik<StudioDocumentSaveResponse>>;
  operacja(zadanie: StudioContextualOpRequest): Promise<Wynik<StudioContextualOpResponse>>;
  porownaj(zadanie: StudioDiffCompareRequest): Promise<Wynik<StudioDiffCompareResponse>>;
  wersje(zadanie: StudioRepositoryListRequest): Promise<Wynik<StudioRepositoryListResponse>>;
  przywroc(
    zadanie: StudioRepositoryRestoreRequest,
  ): Promise<Wynik<StudioRepositoryRestoreResponse>>;
  /** Subskrypcja zdarzenia `studio.document.changed` — zasila podgląd i repozytorium. */
  naZmianeDokumentu(sluchacz: (tresc: StudioDocumentChangedEvent) => void): Odsubskrybuj;
}

export function utworzZrodloStudio(kanal: Kanal): ZrodloStudio {
  return {
    async otworz(zadanie) {
      return sprawdzKsztalt(
        await wywolajUczciwie(kanal, Command.StudioDocumentOpen, zadanie),
        Command.StudioDocumentOpen,
        (tresc) => czyObiekt(tresc.document),
      );
    },

    async zapisz(zadanie) {
      return sprawdzKsztalt(
        await wywolajUczciwie(kanal, Command.StudioDocumentSave, zadanie),
        Command.StudioDocumentSave,
        (tresc) => czyObiekt(tresc.document),
      );
    },

    async operacja(zadanie) {
      return sprawdzKsztalt(
        await wywolajUczciwie(kanal, Command.StudioContextualOp, zadanie),
        Command.StudioContextualOp,
        (tresc) => typeof tresc.messageId === 'string',
      );
    },

    async porownaj(zadanie) {
      // Kontrakt dopuszcza odpowiedź bez fragmentów i bez trafień — jedna
      // komenda obsługuje i porównanie wersji, i wyszukiwanie wzorca, więc
      // brak jednej z tablic jest wynikiem poprawnym, nie usterką kształtu.
      return sprawdzKsztalt(
        await wywolajUczciwie(kanal, Command.StudioDiffCompare, zadanie),
        Command.StudioDiffCompare,
        (tresc) =>
          (tresc.hunks === undefined || czyTablica(tresc.hunks)) &&
          (tresc.matches === undefined || czyTablica(tresc.matches)),
      );
    },

    async wersje(zadanie) {
      return sprawdzKsztalt(
        await wywolajUczciwie(kanal, Command.StudioRepositoryList, zadanie),
        Command.StudioRepositoryList,
        (tresc) => czyTablica(tresc.versions),
      );
    },

    async przywroc(zadanie) {
      return sprawdzKsztalt(
        await wywolajUczciwie(kanal, Command.StudioRepositoryRestore, zadanie),
        Command.StudioRepositoryRestore,
        (tresc) => czyObiekt(tresc.document),
      );
    },

    naZmianeDokumentu(sluchacz) {
      return kanal.naZdarzenie(EventType.StudioDocumentChanged, (tresc) => sluchacz(tresc));
    },
  };
}
