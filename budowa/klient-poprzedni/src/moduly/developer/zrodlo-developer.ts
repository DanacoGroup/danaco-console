import {
  Command,
  EventType,
  PROTOCOL_VERSION,
  type ConnectionHelloResponse,
  type DeveloperBuild,
  type DeveloperBuildChangedEvent,
  type DeveloperBuildRunRequest,
  type DeveloperFile,
  type DeveloperFileOpenRequest,
  type DeveloperFileSaveRequest,
  type DeveloperGitActionRequest,
  type DeveloperTreeGetRequest,
  type DeveloperTreeGetResponse,
  type GitActionResult,
} from '../../../../shared/contract';
import type { Odsubskrybuj } from '../../polaczenie/magistrala-zdarzen';
import type { Kanal, Wynik } from '../../protokol/kanal';
import { czyObiekt, czyTablica, czyTekst, sprawdzKsztalt } from '../../protokol/ksztalt-odpowiedzi';
import { tozsamoscKlienta } from '../../protokol/tozsamosc-klienta';
import { przenies } from '../../protokol/wynik-czastkowy';
import { wywolaj } from '../../protokol/wywolanie';

/**
 * Nazwa odczytu przedstawiana rdzeniowi w powitaniu katalogu komend, służąca
 * wyłącznie czytelności dziennika rdzenia.
 */
const NAZWA_ODCZYTU_KATALOGU = 'developer-katalog-komend';

/**
 * Moduł Developer widziany przez klienta — pięć komend obszaru `developer.*`
 * oraz zdarzenie przyrostu budowania, każda oddająca wynik.
 */
export interface ZrodloDeveloper {
  /** `developer.file.open` — wczytanie pliku repozytorium do Code Editor. */
  otworzPlik(zadanie: DeveloperFileOpenRequest): Promise<Wynik<DeveloperFile>>;
  /** `developer.file.save` — utrwalenie treści pliku repozytorium. */
  zapiszPlik(zadanie: DeveloperFileSaveRequest): Promise<Wynik<DeveloperFile>>;
  /** `developer.tree.get` — drzewo projektu wraz z katalogiem korzenia. */
  drzewo(zadanie: DeveloperTreeGetRequest): Promise<Wynik<DeveloperTreeGetResponse>>;
  /** `developer.git.action` — czynność panelu repozytorium. */
  czynnoscRepozytorium(zadanie: DeveloperGitActionRequest): Promise<Wynik<GitActionResult>>;
  /** `developer.build.run` — uruchomienie albo przerwanie budowania. */
  budowanie(zadanie: DeveloperBuildRunRequest): Promise<Wynik<DeveloperBuild>>;
  /** Subskrypcja `developer.build.changed` — Build Output na żywo. */
  naZmianeBudowania(sluchacz: (tresc: DeveloperBuildChangedEvent) => void): Odsubskrybuj;
  /** `connection.hello` — wykaz komend zarejestrowanych w rdzeniu, z rejestru, nie z kontraktu. */
  komendyRdzenia(): Promise<Wynik<ConnectionHelloResponse>>;
}

export function utworzZrodloDeveloper(kanal: Kanal): ZrodloDeveloper {
  return {
    async otworzPlik(zadanie) {
      const wynik = sprawdzKsztalt(
        await wywolaj(kanal, Command.DeveloperFileOpen, zadanie),
        Command.DeveloperFileOpen,
        (tresc) => czyObiekt(tresc.file) && czyTekst(tresc.file.path),
      );
      return przenies(wynik, (tresc) => tresc.file);
    },

    async zapiszPlik(zadanie) {
      const wynik = sprawdzKsztalt(
        await wywolaj(kanal, Command.DeveloperFileSave, zadanie),
        Command.DeveloperFileSave,
        (tresc) => czyObiekt(tresc.file) && czyTekst(tresc.file.path),
      );
      return przenies(wynik, (tresc) => tresc.file);
    },

    async drzewo(zadanie) {
      return sprawdzKsztalt(
        await wywolaj(kanal, Command.DeveloperTreeGet, zadanie),
        Command.DeveloperTreeGet,
        (tresc) => czyTekst(tresc.root) && czyTablica(tresc.nodes),
      );
    },

    async czynnoscRepozytorium(zadanie) {
      const wynik = sprawdzKsztalt(
        await wywolaj(kanal, Command.DeveloperGitAction, zadanie),
        Command.DeveloperGitAction,
        (tresc) => czyObiekt(tresc.result) && czyTekst(tresc.result.action),
      );
      return przenies(wynik, (tresc) => tresc.result);
    },

    async budowanie(zadanie) {
      const wynik = sprawdzKsztalt(
        await wywolaj(kanal, Command.DeveloperBuildRun, zadanie),
        Command.DeveloperBuildRun,
        (tresc) => czyObiekt(tresc.build) && czyTekst(tresc.build.id),
      );
      return przenies(wynik, (tresc) => tresc.build);
    },

    naZmianeBudowania(sluchacz) {
      return kanal.naZdarzenie(EventType.DeveloperBuildChanged, (tresc) => sluchacz(tresc));
    },

    async komendyRdzenia() {
      return sprawdzKsztalt(
        await wywolaj(kanal, Command.ConnectionHello, {
          clientId: NAZWA_ODCZYTU_KATALOGU,
          clientVersion: tozsamoscKlienta().wersja,
          protocolVersion: PROTOCOL_VERSION,
        }),
        Command.ConnectionHello,
        // `commands` jest w kontrakcie opcjonalne: jego brak to milczenie rdzenia, nie brak komend.
        (tresc) => tresc.commands === undefined || czyTablica(tresc.commands),
      );
    },
  };
}
