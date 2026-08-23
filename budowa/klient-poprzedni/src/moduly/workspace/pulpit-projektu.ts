import {
  Command,
  WorkspaceProjectStatus,
  type WorkspaceDashboard,
} from '../../../../shared/contract';
import type { WykazBrakow } from './braki-kontraktu';
import type { CzynnosciWiedzy } from './czynnosci-wiedzy';
import { utworzRameOkna, type WagaZnacznika } from '../../komponenty/rama-okna';
import {
  pobierzPlik,
  pole,
  pozycjaWykazu,
  przyciskAkcji as przycisk,
  wiersz,
  wykaz,
} from '../../modele/kontrolki-formularza';
import type { StanProjektu } from './stan-projektu';
import { utworzStanTresci } from './stany-okna';
import type { ZrodloWorkspace } from './zrodlo-workspace';

/**
 * Project Dashboard — okno wiodące modułu Workspace: przegląd stanu projektu
 * i punkt wejścia po jego wybraniu. Stąd trzy części okna: wskazanie projektu,
 * zestawienie liczb i skoki do czterech pozostałych okien modułu.
 *
 * Wskazanie projektu jest polem, nie wykazem — okno nie zna komendy zwracającej
 * listę projektów. Wpisany identyfikator idzie do rdzenia komendą
 * `workspace.dashboard.get`, a co rdzeń z nim zrobi, mówi jego odpowiedź.
 *
 * Powody bezczynności przycisków bez pokrycia bierze `braki-kontraktu.ts`
 * z wykazu komend kontraktu, a nie z napisu na stałe.
 */
export interface OknoPulpitu {
  element: HTMLElement;
  /** Odczytuje zestawienie z rdzenia. */
  odswiez(): void;
}

/** Skoki do pozostałych okien modułu. */
export interface SkokiOkien {
  pokaz(kodOkna: string): void;
}

export function utworzOknoPulpitu(
  zrodlo: ZrodloWorkspace,
  wiedza: CzynnosciWiedzy,
  stan: StanProjektu,
  skoki: SkokiOkien,
  braki: WykazBrakow,
): OknoPulpitu {
  const rama = utworzRameOkna({
    tytul: 'Project Dashboard',
    kod: 'project-dashboard',
    rola: 'wiodące',
    przeznaczenie: 'Przegląd stanu projektu i punkt wejścia do pozostałych okien modułu.',
    przedrostek: 'dw',
  });
  const tresc = utworzStanTresci();

  const wskazanie = pole('Identyfikator projektu', 'np. sprawa-2026-114');
  wskazanie.value = stan.projekt();
  const wejdz = przycisk('Wejdź do projektu', 'dn-btn dn-btn--atrament');
  const odswiez = przycisk('Odśwież zestawienie');
  const eksport = przycisk('Eksportuj podsumowanie');
  // Archiwizacja jest zmianą stanu projektu, nie osobną czynnością: kontrakt ma
  // na nią `workspace.project.status.set`, a stan `archived` znaczy „projekt
  // zamknięty, dostępny do podglądu". Dlatego przycisk woła komendę wprost,
  // zamiast stać jako kontrolka bez pokrycia.
  const zarchiwizuj = przycisk('Zarchiwizuj projekt');
  zarchiwizuj.addEventListener('click', () => {
    const projekt = stan.projekt();
    if (projekt === '') {
      tresc.potwierdzenie('Archiwizacja bez wskazania projektu nie ma przedmiotu.', false);
      return;
    }
    void wiedza
      .ustawStanProjektu({ projectId: projekt, status: WorkspaceProjectStatus.Archived })
      .then((wynik) => {
        if (!wynik.udany || wynik.wynik === undefined) {
          tresc.blad('Nie udało się zarchiwizować projektu.', wynik.blad);
          return;
        }
        tresc.potwierdzenie(`Projekt stoi w stanie ${wynik.wynik.status}.`, true);
        odczytaj();
      });
  });

  rama.akcje.append(
    wejdz,
    odswiez,
    eksport,
    braki.przyciskBraku('Duplikuj projekt', 'Duplikowanie projektu', Command.ComponentCreate),
    zarchiwizuj,
    braki.przyciskBraku('Usuń', 'Usunięcie projektu', Command.ComponentDelete),
  );
  const objasnienie =
    'Identyfikator idzie do rdzenia komendą workspace.dashboard.get; co rdzeń z nim ' +
    'zrobił, mówi zestawienie poniżej albo jego odmowa.';
  rama.cialo.append(wiersz('Projekt', wskazanie, { klasa: 'dw-wiersz', objasnienie }), tresc.element);
  const skokiPaska = pasekSkokow(skoki);
  rama.cialo.append(skokiPaska.element);

  let ostatnie: WorkspaceDashboard | null = null;

  function pokaz(zestawienie: WorkspaceDashboard): void {
    ostatnie = zestawienie;
    // Plakietka niesie stan wprost z odpowiedzi rdzenia, bez tłumaczenia go na
    // własne słowo: ten sam napis stoi w wierszu „Stan prac” poniżej.
    rama.ustawZnacznik(zestawienie.project.status, wagaStanu(zestawienie.project.status));
    skokiPaska.ustawLiczniki(zestawienie);
    const lista = wykaz('Zestawienie projektu', 'dw-wykaz');
    const wiersze: ReadonlyArray<[string, string]> = [
      ['Projekt', `${zestawienie.project.name} (${zestawienie.project.id})`],
      ['Stan prac', zestawienie.project.status],
      ['Karty sesji w projekcie', opisWykazu(zestawienie.openSessionIds)],
      ['Eksperci przypisani', opisWykazu(zestawienie.assignedAgentIds)],
      ['Pliki biblioteki projektu', opisLiczby(zestawienie.libraryFileCount)],
      ['Wpisy pamięci projektu', opisLiczby(zestawienie.memoryEntryCount)],
      ['Ostatnia czynność', opisChwili(zestawienie.lastActivityAt)],
    ];
    for (const [nazwa, wartosc] of wiersze) {
      lista.append(pozycjaWykazu(nazwa, wartosc, 'dw').element);
    }
    tresc.tresc().append(lista);
  }

  function odczytaj(): void {
    const projekt = stan.projekt();
    if (projekt === '') {
      tresc.pusto('Wskaż projekt, aby zobaczyć jego stan. Bez projektu okna modułu nie mają zakresu.');
      return;
    }
    tresc.ladowanie('Odczyt zestawienia projektu…');
    void zrodlo.pulpit(projekt).then((wynik) => {
      if (!wynik.udany || wynik.wynik === undefined) {
        // Plakietka i liczniki gasną razem z zestawieniem: pokazane nad odmową
        // mówiłyby o odczycie, który się nie powiódł.
        ostatnie = null;
        rama.ustawZnacznik('');
        skokiPaska.ustawLiczniki(null);
        tresc.blad('Rdzeń nie oddał zestawienia projektu.', wynik.blad);
        return;
      }
      pokaz(wynik.wynik);
      stan.ustawPulpit(wynik.wynik);
    });
  }

  wejdz.addEventListener('click', () => {
    stan.ustawProjekt(wskazanie.value);
    odczytaj();
  });
  odswiez.addEventListener('click', odczytaj);
  eksport.addEventListener('click', () => {
    if (ostatnie === null) {
      tresc.potwierdzenie('Nie ma czego wyeksportować — zestawienie nie zostało jeszcze odczytane.', false);
      return;
    }
    pobierzPlik(`${ostatnie.project.id}-podsumowanie.md`, podsumowanie(ostatnie), 'text/markdown');
    tresc.potwierdzenie('Podsumowanie projektu pobrane jako plik Markdown.', true);
  });

  stan.naZmiane(() => {
    if (wskazanie.value !== stan.projekt()) wskazanie.value = stan.projekt();
    // Zestawienie skasowane — zmianą projektu albo zdarzeniem rdzenia — unieważnia
    // plakietkę stanu, liczniki kafli i materiał eksportu: wszystkie trzy mówiły
    // o projekcie, który nie jest już projektem bieżącym okna.
    if (stan.pulpit() !== null) return;
    ostatnie = null;
    rama.ustawZnacznik('');
    skokiPaska.ustawLiczniki(null);
  });

  return { element: rama.element, odswiez: odczytaj };
}

/**
 * Waga plakietki stanu projektu — trzy stany kontraktu wedle rozdz. 6.2
 * opracowania: czynny sukcesem, wstrzymany ostrzeżeniem, zarchiwizowany
 * neutralnie. Stan spoza kontraktu dostaje wagę ostrzegawczą, bo okno nie wie,
 * co rdzeń przez niego rozumie — przypisanie mu wagi „sukces” albo „neutralna”
 * byłoby zgadywaniem.
 */
const WAGA_STANU: Readonly<Record<string, WagaZnacznika>> = {
  [WorkspaceProjectStatus.Active]: 'sukces',
  [WorkspaceProjectStatus.Paused]: 'ostrzezenie',
  [WorkspaceProjectStatus.Archived]: 'neutralna',
};

function wagaStanu(status: string): WagaZnacznika {
  return WAGA_STANU[status] ?? 'ostrzezenie';
}

/** Kafel nawigacyjny: okno docelowe i — gdy zestawienie go niesie — licznik pozycji. */
interface OpisSkoku {
  kod: string;
  nazwa: string;
  /** Liczony byt okna docelowego; pominięty znaczy „kontrakt licznika nie niesie”. */
  licznik?: {
    nazwa: string;
    ile(zestawienie: WorkspaceDashboard): number | undefined;
  };
}

/**
 * Cztery kafle nawigacyjne pulpitu wraz z licznikami pozycji. Każdy licznik
 * bierze się z pola, które zestawienie pulpitu niesie; liczby wymyślonej nie ma
 * tu żadnej, a pole pominięte przez rdzeń mówi o sobie wprost.
 */
const SKOKI: readonly OpisSkoku[] = [
  {
    kod: 'instructions-panel',
    nazwa: 'Instructions Panel',
    licznik: { nazwa: 'zestawy instrukcji', ile: (zestawienie) => zestawienie.instructionSetCount },
  },
  {
    kod: 'context-memory',
    nazwa: 'Context Memory',
    licznik: { nazwa: 'wpisy pamięci', ile: (zestawienie) => zestawienie.memoryEntryCount },
  },
  {
    kod: 'project-library',
    nazwa: 'Project Library',
    licznik: { nazwa: 'pliki', ile: (zestawienie) => zestawienie.libraryFileCount },
  },
  {
    kod: 'agent-manager',
    nazwa: 'Agent Manager',
    licznik: { nazwa: 'eksperci przypisani', ile: (zestawienie) => zestawienie.assignedAgentIds?.length },
  },
];

/** Pasek kafli nawigacyjnych do pozostałych okien modułu. */
interface PasekSkokow {
  element: HTMLElement;
  /** Przepisuje liczniki; zestawienie nieodczytane mówi to wprost, nie zeruje. */
  ustawLiczniki(zestawienie: WorkspaceDashboard | null): void;
}

function pasekSkokow(skoki: SkokiOkien): PasekSkokow {
  const pasek = document.createElement('nav');
  pasek.className = 'dw-skoki';
  pasek.setAttribute('aria-label', 'Nawigacja do okien modułu Workspace');
  const kafle = SKOKI.map((opis) => {
    const skok = przycisk(etykietaSkoku(opis, null), 'dn-btn dn-btn--zarys');
    skok.addEventListener('click', () => skoki.pokaz(opis.kod));
    pasek.append(skok);
    return { opis, skok };
  });
  return {
    element: pasek,
    ustawLiczniki(zestawienie) {
      for (const kafel of kafle) kafel.skok.textContent = etykietaSkoku(kafel.opis, zestawienie);
    },
  };
}

/**
 * Napis kafla: nazwa okna docelowego i licznik pozycji. Zestawienie
 * nieodczytane nie daje zera — o liczbie pozycji nie wiadomo wtedy nic.
 */
function etykietaSkoku(opis: OpisSkoku, zestawienie: WorkspaceDashboard | null): string {
  if (opis.licznik === undefined) return `→ ${opis.nazwa}`;
  const wartosc =
    zestawienie === null ? 'zestawienie nieodczytane' : opisLiczby(opis.licznik.ile(zestawienie));
  return `→ ${opis.nazwa} (${opis.licznik.nazwa}: ${wartosc})`;
}

/** Wykaz identyfikatorów w jednym zdaniu; wykaz pusty mówi to wprost. */
function opisWykazu(pozycje?: string[]): string {
  if (pozycje === undefined || pozycje.length === 0) return 'brak';
  return `${pozycje.length} — ${pozycje.join(', ')}`;
}

/**
 * Liczba z zestawienia. Pole pominięte przez rdzeń nie jest zerem — `?? 0`
 * stawiałoby przed oczami liczbę, której rdzeń nie podał, nie do odróżnienia
 * od zmierzonego zera.
 */
function opisLiczby(ile?: number): string {
  return ile === undefined ? 'rdzeń nie podał' : String(ile);
}

/** Chwila w postaci lokalnej; brak znacznika mówi „nie odnotowano”. */
function opisChwili(chwila?: number): string {
  if (chwila === undefined || chwila === 0) return 'nie odnotowano';
  return new Date(chwila).toLocaleString('pl-PL');
}

/** Podsumowanie projektu w Markdown — treść pliku eksportu. */
function podsumowanie(zestawienie: WorkspaceDashboard): string {
  return [
    `# Projekt ${zestawienie.project.name}`,
    '',
    `- Identyfikator: ${zestawienie.project.id}`,
    `- Stan: ${zestawienie.project.status}`,
    `- Karty sesji: ${opisWykazu(zestawienie.openSessionIds)}`,
    `- Eksperci: ${opisWykazu(zestawienie.assignedAgentIds)}`,
    `- Pliki biblioteki: ${opisLiczby(zestawienie.libraryFileCount)}`,
    `- Wpisy pamięci: ${opisLiczby(zestawienie.memoryEntryCount)}`,
    `- Ostatnia czynność: ${opisChwili(zestawienie.lastActivityAt)}`,
    '',
  ].join('\n');
}
