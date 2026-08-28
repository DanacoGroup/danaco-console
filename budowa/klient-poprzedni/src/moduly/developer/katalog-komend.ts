import { Command } from '../../../../shared/contract';
import { opisOdmowyBledu } from '../../komponenty/odmowa';
import type { ZrodloDeveloper } from './zrodlo-developer';

/** Katalog komend obszaru `developer.` wzięty z rdzenia, a nie wpisany w kod na stałe. */

/** Przedrostek obszaru, po którym poznajemy komendy tego modułu w wykazie zwróconym przez rdzeń aplikacji. */
const OBSZAR = 'developer.';

/**
 * Komendy, na których stoją okna modułu — te i tylko te są wołane. Wykaz jest
 * pisany ręcznie, nie wywiedziony z kontraktu.
 */
const KOMENDY_OKIEN: readonly string[] = [
  // Code Editor, Project Tree, Git Panel, Build Output.
  Command.DeveloperFileOpen,
  Command.DeveloperFileSave,
  Command.DeveloperTreeGet,
  Command.DeveloperGitAction,
  Command.DeveloperBuildRun,
  Command.DeveloperGitStatus,
  Command.DeveloperGitDiff,
  Command.DeveloperGitLog,
  Command.DeveloperGitBranchList,
  Command.DeveloperGitConflictGet,
  Command.DeveloperGitConflictResolve,
  Command.DeveloperFileCreate,
  Command.DeveloperFileRename,
  Command.DeveloperFileDelete,
  Command.DeveloperFileMove,
  Command.DeveloperGrepSearch,
  Command.DeveloperGrepReplace,
  // Warsztat kodu — pasek operacji Code Editora.
  Command.DeveloperSymbolNavigate,
  Command.DeveloperFormatRun,
  Command.DeveloperLintGet,
  Command.DeveloperRefactorApply,
  Command.DeveloperFileVersionList,
  Command.DeveloperFileVersionRestore,
  Command.DeveloperContextualOp,
  Command.DeveloperToolchainCheck,
  // Kolumna monitora — historia przebiegów i pomiary.
  Command.DeveloperBuildList,
  Command.DeveloperBuildLogGet,
  Command.DeveloperTestResultGet,
  Command.DeveloperCoverageGet,
  // Run & Debug.
  Command.DeveloperDebugSessionStart,
  Command.DeveloperDebugSessionControl,
  Command.DeveloperBreakpointSet,
  Command.DeveloperDebugScopeGet,
  Command.DeveloperDebugEvaluate,
  // Dev Tools — API Client.
  Command.DeveloperApiRequest,
  Command.DeveloperApiCollectionSave,
  Command.DeveloperApiCollectionList,
  Command.DeveloperApiOpenapiImport,
  // Dev Tools — Data Console.
  Command.DeveloperDataConnectionSet,
  Command.DeveloperDataConnectionList,
  Command.DeveloperDataSchemaGet,
  Command.DeveloperDataQueryRun,
  Command.DeveloperDataMigrationRun,
  // Dev Tools — Containers.
  Command.DeveloperContainerList,
  Command.DeveloperContainerAction,
  Command.DeveloperImageBuild,
  Command.DeveloperComposeUp,
  // Dev Tools — Dependencies & Security.
  Command.DeveloperDependencyList,
  Command.DeveloperScanRun,
  Command.DeveloperScanResultList,
];

/** Zdanie wypowiadane, dopóki rdzeń nie odpowiedział na pytanie o katalog komend tego całego obszaru modułu. */
export const KATALOG_W_ODCZYCIE = 'Katalog komend rdzenia — odczyt w toku…';

export interface KatalogKomend {
  /** Zdanie o odmowie albo o milczeniu rdzenia; puste, gdy wykaz przyszedł. */
  odmowa: string;
  /** Komendy obszaru `developer.` zarejestrowane w rdzeniu. */
  obszar: readonly string[];
  /** Komendy obszaru, których żadne okno modułu nie woła. */
  niewolane: readonly string[];
  /** Komendy okien, których rdzeń nie zarejestrował — port martwy. */
  brakujace: readonly string[];
}

/**
 * Pyta rdzeń o jego rejestr komend i zestawia go z komendami okien modułu.
 *
 * Nie rzuca wyjątkiem i nie odrzuca obietnicy: odmowa wraca zdaniem o odmowie,
 * bo okno ma mówić także wtedy, gdy nie udało się ustalić niczego.
 */
export async function odczytajKatalogKomend(zrodlo: ZrodloDeveloper): Promise<KatalogKomend> {
  const pusty = { obszar: [], niewolane: [], brakujace: [] };
  const wynik = await zrodlo.komendyRdzenia();
  if (!wynik.udany || wynik.wynik === undefined) {
    return { odmowa: opisOdmowyBledu('Odczyt katalogu komend rdzenia', wynik.blad), ...pusty };
  }
  // Milczenie nie jest orzeczeniem o braku: rdzeń bez wykazu nie powiedział, że komend nie ma.
  if (wynik.wynik.commands === undefined) {
    return {
      odmowa:
        'Rdzeń przyjął powitanie, ale nie podał wykazu komend — czym rozporządza w obszarze ' +
        'developer, nie da się dziś od niego ustalić',
      ...pusty,
    };
  }
  const obszar = wynik.wynik.commands.filter((komenda) => komenda.startsWith(OBSZAR));
  const wolane = new Set(KOMENDY_OKIEN);
  return {
    odmowa: '',
    obszar,
    niewolane: obszar.filter((komenda) => !wolane.has(komenda)),
    brakujace: KOMENDY_OKIEN.filter((komenda) => !obszar.includes(komenda)),
  };
}

/**
 * Zdanie stanu pustego Git Panelu — skąd okno bierze (albo nie bierze) wiedzę
 * o repozytorium. Nie orzeka o niczym, czego nie powiedział rdzeń.
 */
export function zdanieWiedzyGitPanelu(katalog: KatalogKomend): string {
  const ogon =
    'Wybierz czynność z panelu akcji, a jej wynik będzie pierwszą wiedzą tego okna o repozytorium.';
  if (katalog.odmowa !== '') return `${katalog.odmowa}. ${ogon}`;
  if (katalog.brakujace.length > 0) {
    return (
      `Rdzeń nie zarejestrował komend, na których stoi to okno: ${katalog.brakujace.join(', ')}. ` +
      'Czynności panelu pojadą do rdzenia, który ich nie zna — odpowie odmową nieznanej komendy.'
    );
  }
  if (katalog.niewolane.length > 0) {
    return (
      `Rdzeń rejestruje w obszarze developer komendy, których to okno NIE woła: ` +
      `${katalog.niewolane.join(', ')}. Stan repozytorium może być już do odczytania, ale ten ` +
      `panel jeszcze po niego nie sięga. ${ogon}`
    );
  }
  return (
    `Rdzeń rejestruje w obszarze developer wyłącznie: ${katalog.obszar.join(', ')}. Ani jedna ` +
    `z nich nie oddaje stanu repozytorium bez wykonania czynności — nie ma tu odpowiednika ` +
    `git status, wykazu gałęzi ani historii zatwierdzeń. ${ogon}`
  );
}

/**
 * Zdanie stanu pustego Build Output — dlaczego okno nie zna przebiegów sprzed
 * swojego otwarcia. Także ono wynika z rejestru rdzenia, nie z napisu.
 */
export function zdanieHistoriiBudowania(katalog: KatalogKomend): string {
  const czolo = 'Nie uruchomiono jeszcze budowania w tym oknie w bieżącej sesji.';
  if (katalog.odmowa !== '') return `${czolo} ${katalog.odmowa}.`;
  if (katalog.brakujace.length > 0) {
    return (
      `${czolo} Rdzeń nie zarejestrował komend, na których stoi to okno: ` +
      `${katalog.brakujace.join(', ')} — uruchomienie budowania nie ma dziś dokąd pojechać.`
    );
  }
  if (katalog.niewolane.length > 0) {
    return (
      `${czolo} Rdzeń rejestruje w obszarze developer komendy, których to okno NIE woła: ` +
      `${katalog.niewolane.join(', ')} — historia przebiegów może być już do odczytania.`
    );
  }
  return (
    `${czolo} Rdzeń rejestruje w obszarze developer wyłącznie: ${katalog.obszar.join(', ')} — ` +
    'wykazu przebiegów wśród nich nie ma, więc historia sprzed otwarcia okna jest niewidoczna.'
  );
}
