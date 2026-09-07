// Czynności kontroli w oknie Developer: repozytorium, budowanie wraz z
// dziennikiem, wyniki prób, pokrycie, przegląd bezpieczeństwa i praca krokowa.
import {
  Command,
  ConflictResolutionKind,
  DebugStepKind,
  GitActionKind,
  ScanKind,
} from '../../../shared/contract.ts';
import type {
  DebugFrame,
  DeveloperCoverage,
  DeveloperTestResult,
  GitCommit,
  GitStatusEntry,
} from '../../../shared/contract.ts';
import type { Kanal } from '../protokol/kanal.ts';
import { wywolaj } from '../protokol/wywolanie.ts';
import { oglos } from './ogloszenie.ts';
import {
  czesciWpisu,
  potwierdzone,
  wpisPasa,
  wypelnijPanel,
  zalozPas,
  zwiazPas,
} from './developer-pas.ts';

const NAGLOWEK = 'Kontrola projektu';
const ZNACZNIK = 'devkontrola';

/* Praca krokowa żyje w oknie, bo kontrakt nie ma wykazu jej przebiegów:
   oznaczenie bierze się z rozpoczęcia i służy krokom, odczytowi i wyliczeniu. */
const PRZEBIEGI = new Map<string, string>();

const PRZYCISKI: ReadonlyArray<readonly [string, string]> = [
  ['stan', 'Stan repozytorium'],
  ['roznica', 'Różnica zmian'],
  ['dzieje', 'Dzieje gałęzi'],
  ['zapisz', 'Zapisz zmiany'],
  ['spor', 'Odczytaj spór'],
  ['spor-rozstrzygnij', 'Rozstrzygnij spór'],
  ['budowanie', 'Uruchom budowanie'],
  ['budowanie-dziennik', 'Dziennik budowania'],
  ['proby', 'Wyniki prób'],
  ['pokrycie', 'Pokrycie prób'],
  ['przeglad', 'Przegląd bezpieczeństwa'],
  ['narzedzia', 'Sprawdź narzędzia'],
  ['krokowa', 'Zacznij pracę krokową'],
  ['krok', 'Krok dalej'],
  ['zatrzymanie', 'Postaw zatrzymanie'],
  ['wylicz', 'Wylicz wyrażenie'],
  ['zasieg', 'Odczytaj zasięg'],
];

export function zwiazKontroleDevelopera(
  kanal: Kanal,
  korzen: Element,
  idOkna: () => string,
  odswiez: () => void,
  przy: AddEventListenerOptions,
): void {
  zalozPas(korzen, 'panel-build', ZNACZNIK, 'Zlecenie, ścieżka albo wyrażenie', PRZYCISKI);
  zwiazPas(korzen, ZNACZNIK, (czynnosc) => {
    void wykonaj(kanal, korzen, czynnosc, idOkna(), odswiez);
  }, przy);
}

function wpis(korzen: Element): string {
  return wpisPasa(korzen, ZNACZNIK);
}

function czesci(korzen: Element): string[] {
  return czesciWpisu(korzen, ZNACZNIK);
}

/* Dziennik, wyniki prób i pokrycie idą z budowania stojącego w wykazie
   najwyżej: okno nie prowadzi wskazania budowania, a wykaz podaje kolejność. */
async function ostatnieBudowanie(kanal: Kanal, idOkna: string): Promise<string> {
  const wykaz = await wywolaj(kanal, Command.DeveloperBuildList, { windowId: idOkna });
  return wykaz.wynik?.builds[0]?.id ?? '';
}

async function wykonaj(
  kanal: Kanal,
  korzen: Element,
  czynnosc: string,
  idOkna: string,
  odswiez: () => void,
): Promise<void> {
  if (idOkna === '') {
    oglos(NAGLOWEK, 'Rdzeń nie dał okna projektu dla tej karty.', 'ostrzezenie');
    return;
  }
  if (czynnosc === 'stan') return stanRepozytorium(kanal, korzen, idOkna);
  if (czynnosc === 'roznica') return roznicaZmian(kanal, korzen, idOkna);
  if (czynnosc === 'dzieje') return dziejeGalezi(kanal, korzen, idOkna);
  if (czynnosc === 'zapisz') return zapiszZmiany(kanal, korzen, idOkna, odswiez);
  if (czynnosc === 'spor') return odczytajSpor(kanal, korzen, idOkna);
  if (czynnosc === 'spor-rozstrzygnij') return rozstrzygnijSpor(kanal, korzen, idOkna);
  if (czynnosc === 'budowanie') return uruchomBudowanie(kanal, korzen, idOkna, odswiez);
  if (czynnosc === 'budowanie-dziennik') return dziennikBudowania(kanal, korzen, idOkna);
  if (czynnosc === 'proby') return wynikiProb(kanal, korzen, idOkna);
  if (czynnosc === 'pokrycie') return pokrycieProb(kanal, korzen, idOkna);
  if (czynnosc === 'przeglad') return przegladBezpieczenstwa(kanal, idOkna);
  if (czynnosc === 'narzedzia') return sprawdzNarzedzia(kanal, korzen);
  if (czynnosc === 'krokowa') return zacznijPraceKrokowa(kanal, korzen, idOkna);
  if (czynnosc === 'krok') return krokDalej(kanal, idOkna);
  if (czynnosc === 'zatrzymanie') return postawZatrzymanie(kanal, korzen, idOkna);
  if (czynnosc === 'wylicz') return wyliczWyrazenie(kanal, korzen, idOkna);
  if (czynnosc === 'zasieg') return odczytajZasieg(kanal, korzen, idOkna);
  /* Czynność spoza obsłużonych odmawia zamiast milczeć: przycisk bez gałęzi
     wyglądałby jak działający. */
  oglos(NAGLOWEK, `Czynność „${czynnosc}" nie jest prowadzona przez to okno.`, 'ostrzezenie');
}

async function stanRepozytorium(kanal: Kanal, korzen: Element, idOkna: string): Promise<void> {
  const wynik = await wywolaj(kanal, Command.DeveloperGitStatus, { windowId: idOkna });
  if (!wynik.udany || wynik.wynik === undefined) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił odczytu stanu.', 'ostrzezenie');
    return;
  }
  const stan = wynik.wynik.status;
  if (!stan.isRepository) {
    oglos(NAGLOWEK, 'Katalog roboczy tego okna nie jest repozytorium.', 'ostrzezenie');
    return;
  }
  wypelnijPanel(korzen, 'panel-git', stan.entries.map((pozycja: GitStatusEntry) =>
    `${pozycja.index}/${pozycja.worktree} ${pozycja.path}`));
  oglos(NAGLOWEK, `Gałąź ${stan.branch ?? 'bez nazwy'}: ${stan.entries.length} zmian, `
    + `${stan.hasConflicts ? 'są spory' : 'bez sporów'}.`);
}

async function roznicaZmian(kanal: Kanal, korzen: Element, idOkna: string): Promise<void> {
  const wynik = await wywolaj(kanal, Command.DeveloperGitDiff, { windowId: idOkna });
  if (!wynik.udany || wynik.wynik === undefined) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił odczytu różnicy.', 'ostrzezenie');
    return;
  }
  const fragmenty = wynik.wynik.hunks;
  wypelnijPanel(korzen, 'panel-git', fragmenty.map((fragment) =>
    `${fragment.path} · ${fragment.header ?? ''}`));
  if (fragmenty.length === 0) oglos(NAGLOWEK, 'Katalog roboczy nie ma niezapisanych zmian.');
}

async function dziejeGalezi(kanal: Kanal, korzen: Element, idOkna: string): Promise<void> {
  const wynik = await wywolaj(kanal, Command.DeveloperGitLog, { windowId: idOkna, limit: 50 });
  if (!wynik.udany || wynik.wynik === undefined) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił odczytu dziejów.', 'ostrzezenie');
    return;
  }
  const zapisy = wynik.wynik.commits;
  wypelnijPanel(korzen, 'panel-git', zapisy.map((zapis: GitCommit) =>
    `${zapis.shortId} ${zapis.author} · ${zapis.message.split('\n')[0] ?? ''}`));
  if (zapisy.length === 0) oglos(NAGLOWEK, 'Ta gałąź nie ma jeszcze zapisów.');
}

/* Zapis zmian bierze wszystkie pozycje stanu i zamyka je jednym opisem z pola:
   okno nie prowadzi wyboru pozycji, więc dzielenie ich byłoby zgadywaniem. */
async function zapiszZmiany(
  kanal: Kanal,
  korzen: Element,
  idOkna: string,
  odswiez: () => void,
): Promise<void> {
  const opis = wpis(korzen);
  if (opis === '') {
    oglos(NAGLOWEK, 'Zapis zmian potrzebuje opisu wpisanego w polu.', 'ostrzezenie');
    return;
  }
  const stan = await wywolaj(kanal, Command.DeveloperGitStatus, { windowId: idOkna });
  const sciezki = (stan.wynik?.status.entries ?? []).map((pozycja) => pozycja.path);
  if (sciezki.length === 0) {
    oglos(NAGLOWEK, 'Nie ma zmian do zapisania.', 'ostrzezenie');
    return;
  }
  const przygotowanie = await wywolaj(kanal, Command.DeveloperGitAction, {
    windowId: idOkna,
    action: GitActionKind.Stage,
    paths: sciezki,
  });
  if (!przygotowanie.udany) {
    oglos(NAGLOWEK, przygotowanie.blad?.message ?? 'Rdzeń odmówił przygotowania zmian.', 'ostrzezenie');
    return;
  }
  const wynik = await wywolaj(kanal, Command.DeveloperGitAction, {
    windowId: idOkna,
    action: GitActionKind.Commit,
    message: opis,
  });
  if (!wynik.udany) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił zapisu zmian.', 'ostrzezenie');
    return;
  }
  oglos(NAGLOWEK, `Zapisano ${sciezki.length} zmienionych ścieżek.`);
  odswiez();
}

async function odczytajSpor(kanal: Kanal, korzen: Element, idOkna: string): Promise<void> {
  const sciezka = wpis(korzen);
  if (sciezka === '') {
    oglos(NAGLOWEK, 'Odczyt sporu potrzebuje ścieżki wpisanej w polu.', 'ostrzezenie');
    return;
  }
  const wynik = await wywolaj(kanal, Command.DeveloperGitConflictGet, {
    windowId: idOkna,
    path: sciezka,
  });
  if (!wynik.udany || wynik.wynik === undefined) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił odczytu sporu.', 'ostrzezenie');
    return;
  }
  const spor = wynik.wynik.conflict;
  wypelnijPanel(korzen, 'panel-edytor', [
    '— strona bieżąca —',
    ...spor.current.split('\n'),
    '— strona przychodząca —',
    ...spor.incoming.split('\n'),
  ]);
  oglos(NAGLOWEK, `Spór w pliku ${spor.path} stoi w edytorze.`);
}

/* Rozstrzygnięcie nadpisuje plik, więc pierwsze naciśnięcie uzbraja przycisk;
   okno bierze stronę bieżącą, bo tylko ją Operator ma przed sobą w edytorze. */
async function rozstrzygnijSpor(kanal: Kanal, korzen: Element, idOkna: string): Promise<void> {
  const sciezka = wpis(korzen);
  if (sciezka === '') {
    oglos(NAGLOWEK, 'Rozstrzygnięcie potrzebuje ścieżki wpisanej w polu.', 'ostrzezenie');
    return;
  }
  if (!potwierdzone(korzen, ZNACZNIK, 'spor-rozstrzygnij')) return;
  const wynik = await wywolaj(kanal, Command.DeveloperGitConflictResolve, {
    windowId: idOkna,
    path: sciezka,
    resolution: ConflictResolutionKind.TakeCurrent,
  });
  if (!wynik.udany || wynik.wynik === undefined) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił rozstrzygnięcia sporu.', 'ostrzezenie');
    return;
  }
  oglos(NAGLOWEK, wynik.wynik.resolved
    ? `Spór w ${sciezka} rozstrzygnięty; zostaje ${wynik.wynik.remainingPaths.length} spornych plików.`
    : 'Rdzeń nie uznał sporu za rozstrzygnięty.',
  wynik.wynik.resolved ? 'informacja' : 'ostrzezenie');
}

async function uruchomBudowanie(
  kanal: Kanal,
  korzen: Element,
  idOkna: string,
  odswiez: () => void,
): Promise<void> {
  const zlecenie = wpis(korzen);
  if (zlecenie === '') {
    oglos(NAGLOWEK, 'Budowanie potrzebuje nazwy zlecenia wpisanej w polu.', 'ostrzezenie');
    return;
  }
  const wynik = await wywolaj(kanal, Command.DeveloperBuildRun, {
    windowId: idOkna,
    task: zlecenie,
  });
  if (!wynik.udany || wynik.wynik === undefined) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił uruchomienia budowania.', 'ostrzezenie');
    return;
  }
  oglos(NAGLOWEK, `Budowanie „${zlecenie}" w stanie ${wynik.wynik.build.status}.`);
  odswiez();
}

async function dziennikBudowania(kanal: Kanal, korzen: Element, idOkna: string): Promise<void> {
  const cel = await ostatnieBudowanie(kanal, idOkna);
  if (cel === '') {
    oglos(NAGLOWEK, 'Żadne budowanie nie ruszyło w tym oknie.', 'ostrzezenie');
    return;
  }
  const wynik = await wywolaj(kanal, Command.DeveloperBuildLogGet, { buildId: cel, tail: 200 });
  if (!wynik.udany || wynik.wynik === undefined) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił odczytu dziennika.', 'ostrzezenie');
    return;
  }
  wypelnijPanel(korzen, 'panel-build', wynik.wynik.lines);
  if (wynik.wynik.truncated) {
    oglos(NAGLOWEK, 'Dziennik przycięty granicą bufora — wierszy jest więcej.', 'ostrzezenie');
  }
}

async function wynikiProb(kanal: Kanal, korzen: Element, idOkna: string): Promise<void> {
  const cel = await ostatnieBudowanie(kanal, idOkna);
  if (cel === '') {
    oglos(NAGLOWEK, 'Żadne budowanie nie ruszyło w tym oknie.', 'ostrzezenie');
    return;
  }
  const wynik = await wywolaj(kanal, Command.DeveloperTestResultGet, { buildId: cel });
  if (!wynik.udany || wynik.wynik === undefined) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił odczytu wyników prób.', 'ostrzezenie');
    return;
  }
  wypelnijPanel(korzen, 'panel-zadania', wynik.wynik.results.map((proba: DeveloperTestResult) =>
    `${proba.status} ${proba.name}`));
  oglos(NAGLOWEK, `Próby: ${wynik.wynik.passed} zdanych, ${wynik.wynik.failed} niezdanych, `
    + `${wynik.wynik.skipped} pominiętych.`,
  wynik.wynik.failed === 0 ? 'informacja' : 'ostrzezenie');
}

async function pokrycieProb(kanal: Kanal, korzen: Element, idOkna: string): Promise<void> {
  const cel = await ostatnieBudowanie(kanal, idOkna);
  if (cel === '') {
    oglos(NAGLOWEK, 'Żadne budowanie nie ruszyło w tym oknie.', 'ostrzezenie');
    return;
  }
  const wynik = await wywolaj(kanal, Command.DeveloperCoverageGet, { buildId: cel });
  if (!wynik.udany || wynik.wynik === undefined) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił odczytu pokrycia.', 'ostrzezenie');
    return;
  }
  wypelnijPanel(korzen, 'panel-zadania', wynik.wynik.files.map((plik: DeveloperCoverage) =>
    `${plik.percent}% ${plik.path}`));
  oglos(NAGLOWEK, `Pokrycie prób: ${wynik.wynik.percent}%.`);
}

async function przegladBezpieczenstwa(kanal: Kanal, idOkna: string): Promise<void> {
  const wynik = await wywolaj(kanal, Command.DeveloperScanRun, {
    windowId: idOkna,
    kinds: [ScanKind.Dependencies, ScanKind.Secrets, ScanKind.Code, ScanKind.Licenses],
  });
  if (!wynik.udany || wynik.wynik === undefined) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił przeglądu.', 'ostrzezenie');
    return;
  }
  const przeglad = wynik.wynik.scan;
  oglos(NAGLOWEK, `Przegląd w stanie ${przeglad.status}; `
    + `ustaleń ${przeglad.findingCount ?? 0}.`);
}

/* Wykaz narzędzi bierze się z pola rozdzielony spacjami; bez wpisu okno pyta
   o zestaw, którym samo się posługuje. */
async function sprawdzNarzedzia(kanal: Kanal, korzen: Element): Promise<void> {
  const wpisane = wpis(korzen).split(/\s+/).filter((nazwa) => nazwa !== '');
  const programy = wpisane.length === 0 ? ['git', 'go', 'node', 'python3'] : wpisane;
  const wynik = await wywolaj(kanal, Command.DeveloperToolchainCheck, { programs: programy });
  if (!wynik.udany || wynik.wynik === undefined) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił sprawdzenia narzędzi.', 'ostrzezenie');
    return;
  }
  wypelnijPanel(korzen, 'panel-terminal', wynik.wynik.programs.map((program) =>
    `${program.present ? '✓' : '✗'} ${program.program} ${program.version ?? ''}`.trim()));
  const brakujace = wynik.wynik.programs.filter((program) => !program.present).length;
  if (brakujace > 0) {
    oglos(NAGLOWEK, `${brakujace} narzędzi nie stoi na maszynie rdzenia.`, 'ostrzezenie');
  }
}

async function zacznijPraceKrokowa(kanal: Kanal, korzen: Element, idOkna: string): Promise<void> {
  const program = wpis(korzen);
  const wynik = await wywolaj(kanal, Command.DeveloperDebugSessionStart, {
    windowId: idOkna,
    ...(program === '' ? {} : { program }),
    stopOnEntry: true,
  });
  if (!wynik.udany || wynik.wynik === undefined) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił rozpoczęcia pracy krokowej.', 'ostrzezenie');
    return;
  }
  PRZEBIEGI.set(idOkna, wynik.wynik.session.id);
  oglos(NAGLOWEK, `Praca krokowa stoi na ${wynik.wynik.session.adapter}, `
    + `stan ${wynik.wynik.session.status}.`);
}

function przebieg(idOkna: string): string {
  const zapamietany = PRZEBIEGI.get(idOkna) ?? '';
  if (zapamietany === '') {
    oglos(NAGLOWEK, 'Okno nie prowadzi pracy krokowej — zacznij ją najpierw.', 'ostrzezenie');
  }
  return zapamietany;
}

async function krokDalej(kanal: Kanal, idOkna: string): Promise<void> {
  const cel = przebieg(idOkna);
  if (cel === '') return;
  const wynik = await wywolaj(kanal, Command.DeveloperDebugSessionControl, {
    sessionId: cel,
    step: DebugStepKind.StepOver,
  });
  if (!wynik.udany || wynik.wynik === undefined) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił kroku.', 'ostrzezenie');
    return;
  }
  oglos(NAGLOWEK, `Praca krokowa w stanie ${wynik.wynik.session.status}`
    + `${wynik.wynik.session.stoppedReason === undefined ? '' : `: ${wynik.wynik.session.stoppedReason}`}.`);
}

async function postawZatrzymanie(kanal: Kanal, korzen: Element, idOkna: string): Promise<void> {
  const [sciezka = '', wiersz = ''] = czesci(korzen);
  const numer = Number.parseInt(wiersz, 10);
  if (sciezka === '' || !Number.isFinite(numer)) {
    oglos(NAGLOWEK,
      'Zatrzymanie potrzebuje ścieżki i wiersza, na przykład „src/plik.ts | 42".',
      'ostrzezenie');
    return;
  }
  const wynik = await wywolaj(kanal, Command.DeveloperBreakpointSet, {
    windowId: idOkna,
    path: sciezka,
    line: numer,
  });
  if (!wynik.udany || wynik.wynik === undefined) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił postawienia zatrzymania.', 'ostrzezenie');
    return;
  }
  oglos(NAGLOWEK, `Zatrzymań w oknie: ${wynik.wynik.breakpoints.length}.`);
}

/* Wyliczenie idzie w ramce stojącej najwyżej: okno nie prowadzi wyboru ramki,
   a zasięg podaje ich kolejność od miejsca zatrzymania. */
async function wyliczWyrazenie(kanal: Kanal, korzen: Element, idOkna: string): Promise<void> {
  const wyrazenie = wpis(korzen);
  if (wyrazenie === '') {
    oglos(NAGLOWEK, 'Wyliczenie potrzebuje wyrażenia wpisanego w polu.', 'ostrzezenie');
    return;
  }
  const cel = przebieg(idOkna);
  if (cel === '') return;
  const zasieg = await wywolaj(kanal, Command.DeveloperDebugScopeGet, { sessionId: cel });
  const ramka = zasieg.wynik?.frames[0]?.id ?? '';
  if (ramka === '') {
    oglos(NAGLOWEK, 'Praca krokowa nie stoi w żadnej ramce.', 'ostrzezenie');
    return;
  }
  const wynik = await wywolaj(kanal, Command.DeveloperDebugEvaluate, {
    sessionId: cel,
    frameId: ramka,
    expression: wyrazenie,
  });
  if (!wynik.udany || wynik.wynik === undefined) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił wyliczenia.', 'ostrzezenie');
    return;
  }
  oglos(NAGLOWEK, `${wyrazenie} = ${wynik.wynik.value}`);
}

async function odczytajZasieg(kanal: Kanal, korzen: Element, idOkna: string): Promise<void> {
  const cel = przebieg(idOkna);
  if (cel === '') return;
  const wynik = await wywolaj(kanal, Command.DeveloperDebugScopeGet, { sessionId: cel });
  if (!wynik.udany || wynik.wynik === undefined) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił odczytu zasięgu.', 'ostrzezenie');
    return;
  }
  wypelnijPanel(korzen, 'panel-kolejka', [
    ...wynik.wynik.frames.map((ramka: DebugFrame) =>
      `ramka ${ramka.name} · ${ramka.path ?? ''}:${ramka.line ?? 0}`),
    ...wynik.wynik.variables.map((zmienna) => `${zmienna.name} = ${zmienna.value}`),
  ]);
  if (wynik.wynik.frames.length === 0) oglos(NAGLOWEK, 'Praca krokowa nie stoi w żadnej ramce.');
}
