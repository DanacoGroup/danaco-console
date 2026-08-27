import {
  BuildStatus,
  ChangeKind,
  GitActionKind,
  TreeNodeKind,
  type DeveloperBuild,
  type DeveloperBuildChangedEvent,
  type DeveloperBuildRunRequest,
  type DeveloperFile,
  type DeveloperFileOpenRequest,
  type DeveloperFileSaveRequest,
  type DeveloperGitActionRequest,
  type DeveloperTreeGetRequest,
  type DeveloperTreeGetResponse,
  type DeveloperTreeNode,
  type GitActionResult,
} from '../../../../shared/contract';
import type { ZrodloDeveloper } from './zrodlo-developer';

/**
 * Źródło próbne modułu Developer zastępuje rdzeń zapisem wysłanych żądań oraz
 * gotowymi odpowiedziami podstawianymi przez sprawdzian. Umowa `ZrodloDeveloper`
 * obowiązuje je tak samo jak źródło produktowe.
 */
export interface ZrodloProbne extends ZrodloDeveloper {
  /** Żądania, które okna naprawdę wysłały — w kolejności wysłania. */
  zadania: {
    otwarcia: DeveloperFileOpenRequest[];
    zapisy: DeveloperFileSaveRequest[];
    drzewa: DeveloperTreeGetRequest[];
    czynnosci: DeveloperGitActionRequest[];
    budowania: DeveloperBuildRunRequest[];
  };
  /** Odpowiedzi podstawiane oknom; sprawdzian podmienia je przed czynnością. */
  odpowiedzi: {
    plik: DeveloperFile | null;
    drzewo: DeveloperTreeGetResponse | null;
    czynnosc: GitActionResult | null;
    przebieg: DeveloperBuild | null;
    komendy: string[] | undefined;
  };
  /** Wysyła oknom zdarzenie przyrostu budowania, tak jak zrobiłby to rdzeń. */
  wyslijZdarzenieBudowania(zdarzenie: DeveloperBuildChangedEvent): void;
}

/**
 * Składa plik gotowy do podstawienia oknu. Pole `updatedAt` jest wymagane przez
 * kontrakt, więc stoi na stałej chwili, a oznaczenie wersji wchodzi wyłącznie
 * wtedy, gdy sprawdzian je poda.
 */
export function plikProbny(path: string, content: string, versionId?: string): DeveloperFile {
  const plik: DeveloperFile = { path, content, updatedAt: 1_755_000_000_000 };
  if (versionId !== undefined) plik.versionId = versionId;
  return plik;
}

/**
 * Składa węzeł drzewa plików o wskazanej ścieżce, nazwie i rodzaju. Pominięta
 * ścieżka rodzica znaczy węzeł korzenia, więc pole wchodzi do węzła dopiero po
 * podaniu wartości.
 */
export function wezelProbny(
  path: string,
  name: string,
  kind: TreeNodeKind = TreeNodeKind.File,
  parentPath?: string,
): DeveloperTreeNode {
  const wezel: DeveloperTreeNode = { path, name, kind };
  if (parentPath !== undefined) wezel.parentPath = parentPath;
  return wezel;
}

/**
 * Składa przebieg budowania w stanie trwającym, o ile sprawdzian nie poda stanu
 * innego. Okno i chwila rozpoczęcia stoją na stałych wartościach, ponieważ
 * sprawdziany badają stan przebiegu, a nie jego pochodzenie.
 */
export function przebiegProbny(
  id: string,
  task: string,
  status: BuildStatus = BuildStatus.Running,
): DeveloperBuild {
  return { id, windowId: 'okno-1', task, status, startedAt: 1_755_000_000_000 };
}

/**
 * Zdarzenie przyrostu budowania z wypełnionym `change`.
 *
 * Pole jest wymagane przez kontrakt, choć okna go nie czytają; wypełnia je ta
 * funkcja, aby sprawdziany nie powielały tego u siebie.
 */
export function zdarzenieBudowania(
  build: DeveloperBuild,
  logLine?: string,
): DeveloperBuildChangedEvent {
  const zdarzenie: DeveloperBuildChangedEvent = { change: ChangeKind.Updated, build };
  if (logLine !== undefined) zdarzenie.logLine = logLine;
  return zdarzenie;
}

export function utworzZrodloProbne(): ZrodloProbne {
  const sluchacze = new Set<(tresc: DeveloperBuildChangedEvent) => void>();
  const zrodlo: ZrodloProbne = {
    zadania: { otwarcia: [], zapisy: [], drzewa: [], czynnosci: [], budowania: [] },
    odpowiedzi: {
      plik: null,
      drzewo: null,
      czynnosc: null,
      przebieg: null,
      komendy: undefined,
    },

    async otworzPlik(zadanie) {
      zrodlo.zadania.otwarcia.push(zadanie);
      const plik = zrodlo.odpowiedzi.plik ?? plikProbny(zadanie.path, `treść ${zadanie.path}`);
      return { udany: true, wynik: plik };
    },

    async zapiszPlik(zadanie) {
      zrodlo.zadania.zapisy.push(zadanie);
      const plik = zrodlo.odpowiedzi.plik ?? plikProbny(zadanie.path, zadanie.content);
      return { udany: true, wynik: plik };
    },

    async drzewo(zadanie) {
      zrodlo.zadania.drzewa.push(zadanie);
      return { udany: true, wynik: zrodlo.odpowiedzi.drzewo ?? { root: '', nodes: [] } };
    },

    async czynnoscRepozytorium(zadanie) {
      zrodlo.zadania.czynnosci.push(zadanie);
      const wynik = zrodlo.odpowiedzi.czynnosc ?? {
        action: zadanie.action ?? GitActionKind.Stage,
        succeeded: true,
      };
      return { udany: true, wynik };
    },

    async budowanie(zadanie) {
      zrodlo.zadania.budowania.push(zadanie);
      const przebieg = zrodlo.odpowiedzi.przebieg ?? przebiegProbny('bud-1', zadanie.task);
      return { udany: true, wynik: przebieg };
    },

    naZmianeBudowania(sluchacz) {
      sluchacze.add(sluchacz);
      return () => sluchacze.delete(sluchacz);
    },

    async komendyRdzenia() {
      const odpowiedz =
        zrodlo.odpowiedzi.komendy === undefined
          ? { protocolVersion: '1', serverVersion: 'probny' }
          : { protocolVersion: '1', serverVersion: 'probny', commands: zrodlo.odpowiedzi.komendy };
      return { udany: true, wynik: odpowiedz };
    },

    wyslijZdarzenieBudowania(zdarzenie) {
      for (const sluchacz of [...sluchacze]) sluchacz(zdarzenie);
    },
  };
  return zrodlo;
}
