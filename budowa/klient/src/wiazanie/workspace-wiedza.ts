// Panel „W tle" projektu WorkSpace: tablica wizualna, graf wiedzy
// i wydanie kalendarza. Komendy `workspace.canvas.get`,
// `workspace.knowledge.graph.get` i `workspace.calendar.export`.
import { Command } from '../../../shared/contract.ts';
import type { Kanal } from '../protokol/kanal.ts';
import { wywolaj } from '../protokol/wywolanie.ts';
import { oglos } from './ogloszenie.ts';
import { cialoPanelu, przyjmij, type WezlyWorkspace } from './workspace-wezly.ts';

const NAGLOWEK = 'Wiedza projektu';

export interface StanWiedzy {
  odswiez: () => Promise<void>;
}

/*
zwiazWiedze dokłada do panelu biblioteki trzy zakresy, których okno nie miało.

Tablica i graf są odczytem, wydanie kalendarza czynnością: zapis iCal idzie do
schowka wraz z nazwą pliku podaną przez rdzeń, bo okno nie ma miejsca, w którym
mogłoby zostawić plik.
*/
export function zwiazWiedze(
  wezly: WezlyWorkspace,
  kanal: Kanal,
  idProjektu: () => string,
  przy: AddEventListenerOptions,
): StanWiedzy {
  const cialo = cialoPanelu(wezly.wTle);
  const odswiez = async (): Promise<void> => {
    const projekt = idProjektu();
    if (cialo === null || projekt === '') return;
    cialo.replaceChildren();
    postawPasWydania(cialo);
    await Promise.all([
      opiszTablice(kanal, cialo, projekt),
      opiszGraf(kanal, cialo, projekt),
    ]);
  };

  wezly.wTle?.addEventListener('click', (zdarzenie) => {
    const cel = zdarzenie.target;
    if (!(cel instanceof Element)) return;
    if (cel.closest('[data-kalendarz-wydaj]') !== null) {
      zdarzenie.stopPropagation();
      void wydajKalendarz(kanal, idProjektu());
      return;
    }
    if (cel.closest('[data-kalendarz-wciagnij]') !== null) {
      zdarzenie.stopPropagation();
      void wciagnijKalendarz(kanal, wezly, idProjektu(), odswiez);
      return;
    }
    if (cel.closest('[data-tablica-zaloz]') !== null) {
      zdarzenie.stopPropagation();
      void zalozTablice(kanal, wezly, idProjektu(), odswiez);
      return;
    }
  }, przy);

  return { odswiez };
}

function postawPasWydania(cialo: HTMLElement): void {
  const pole = cialo.ownerDocument.createElement('input');
  pole.type = 'text';
  pole.className = 'dn-pole dn-pole--sm';
  pole.placeholder = 'Zapis iCal albo nazwa tablicy';
  pole.setAttribute('aria-label', 'Zapis iCal albo nazwa tablicy');
  pole.dataset.wiedzaWpis = '';
  cialo.append(
    pole,
    przyciskPasa(cialo, 'kalendarzWydaj', 'Wydaj kalendarz'),
    przyciskPasa(cialo, 'kalendarzWciagnij', 'Wciągnij kalendarz'),
    przyciskPasa(cialo, 'tablicaZaloz', 'Załóż tablicę'),
  );
}

function przyciskPasa(cialo: HTMLElement, znacznik: string, etykieta: string): HTMLButtonElement {
  const przycisk = cialo.ownerDocument.createElement('button');
  przycisk.type = 'button';
  przycisk.className = 'dn-btn dn-btn--duch dn-btn--sm';
  przycisk.dataset[znacznik] = '';
  przycisk.textContent = etykieta;
  return przycisk;
}

// Tablica bez pozycji nie jest odmową: projekt może jej po prostu nie mieć.
async function opiszTablice(kanal: Kanal, cialo: HTMLElement, idProjektu: string): Promise<void> {
  const wynik = await wywolaj(kanal, Command.WorkspaceCanvasGet, { projectId: idProjektu });
  const odpowiedz = przyjmij('Tablica projektu', wynik);
  if (odpowiedz === null) return;
  const tablice = odpowiedz.canvasIds?.length ?? 0;
  /* Scena tablicy jest nieprzejrzysta dla rdzenia — jej kształt należy do
     widoku, który ją rysuje — więc wiersz podaje nazwę, nie liczbę elementów. */
  dopiszWiersz(cialo, 'Tablica wizualna',
    odpowiedz.canvas === undefined
      ? 'projekt nie ma jeszcze tablicy'
      : `${odpowiedz.canvas.name} · tablic: ${tablice}`);
}

async function opiszGraf(kanal: Kanal, cialo: HTMLElement, idProjektu: string): Promise<void> {
  const wynik = await wywolaj(kanal, Command.WorkspaceKnowledgeGraphGet, { projectId: idProjektu });
  const odpowiedz = przyjmij('Graf wiedzy', wynik);
  if (odpowiedz === null) return;
  const graf = odpowiedz.graph;
  dopiszWiersz(cialo, 'Graf wiedzy',
    `${graf.nodes?.length ?? 0} bytów · ${graf.edges?.length ?? 0} powiązań`);
}

function dopiszWiersz(cialo: HTMLElement, tytul: string, podpis: string): void {
  const wiersz = cialo.ownerDocument.createElement('div');
  wiersz.className = 'dn-wykaz-modulu-poz';
  const nazwa = cialo.ownerDocument.createElement('span');
  nazwa.textContent = tytul;
  const meta = cialo.ownerDocument.createElement('span');
  meta.className = 'dn-meta';
  meta.textContent = podpis;
  wiersz.append(nazwa, meta);
  cialo.appendChild(wiersz);
}

function wpisPasa(wezly: WezlyWorkspace): string {
  const pole = wezly.wTle?.querySelector<HTMLInputElement>('[data-wiedza-wpis]');
  return (pole?.value ?? '').trim();
}

/* Kalendarz wchodzi zapisem iCal wklejonym w polu: okno nie ma wskazania pliku,
   a rdzeń przyjmuje treść, nie ścieżkę. Wydarzenia zakładają zadania — taki
   jest domyślny kształt komendy. */
async function wciagnijKalendarz(
  kanal: Kanal,
  wezly: WezlyWorkspace,
  idProjektu: string,
  odswiez: () => Promise<void>,
): Promise<void> {
  const zapis = wpisPasa(wezly);
  if (idProjektu === '' || zapis === '') {
    oglos(NAGLOWEK, 'Wciągnięcie kalendarza potrzebuje zapisu iCal wklejonego w polu.',
      'ostrzezenie');
    return;
  }
  const wynik = await wywolaj(kanal, Command.WorkspaceCalendarImport, {
    projectId: idProjektu,
    content: zapis,
  });
  const odpowiedz = przyjmij(NAGLOWEK, wynik);
  if (odpowiedz === null) return;
  oglos(NAGLOWEK, 'Kalendarz wciągnięty do projektu.');
  await odswiez();
}

/* Scena nowej tablicy jest pusta: kształt sceny należy do widoku, który ją
   rysuje, więc okno zakłada tablicę nazwaną i nic w niej nie zmyśla. */
async function zalozTablice(
  kanal: Kanal,
  wezly: WezlyWorkspace,
  idProjektu: string,
  odswiez: () => Promise<void>,
): Promise<void> {
  const nazwa = wpisPasa(wezly);
  if (idProjektu === '' || nazwa === '') {
    oglos(NAGLOWEK, 'Nowa tablica potrzebuje nazwy wpisanej w polu.', 'ostrzezenie');
    return;
  }
  const wynik = await wywolaj(kanal, Command.WorkspaceCanvasSave, {
    projectId: idProjektu,
    name: nazwa,
    scene: {},
  });
  if (przyjmij(NAGLOWEK, wynik) === null) return;
  oglos(NAGLOWEK, `Tablica „${nazwa}" założona.`);
  await odswiez();
}

async function wydajKalendarz(kanal: Kanal, idProjektu: string): Promise<void> {
  if (idProjektu === '') return;
  const wynik = await wywolaj(kanal, Command.WorkspaceCalendarExport, { projectId: idProjektu });
  if (!wynik.udany || wynik.wynik === undefined) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił wydania kalendarza.', 'ostrzezenie');
    return;
  }
  const { content, fileName, exported } = wynik.wynik;
  try {
    await navigator.clipboard.writeText(content);
    oglos(NAGLOWEK, `Kalendarz w schowku: ${exported} pozycji, nazwa pliku „${fileName}".`);
  } catch {
    oglos(NAGLOWEK, 'Przeglądarka nie dała dostępu do schowka, więc zapis nie został przeniesiony.',
      'ostrzezenie');
  }
}
