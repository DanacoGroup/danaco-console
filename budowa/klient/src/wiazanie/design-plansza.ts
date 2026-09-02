// Design Board: kompozycja rdzenia rysowana na kanwie prototypu, jej warstwy,
// ramki, siatka, grupowanie i wyniesienie.

import {
  Command,
  DesignLayoutDirection,
  EventType,
  type DesignBoard,
  type DesignBoardLayer,
  type DesignFrame,
} from '../../../shared/contract.ts';
import { zglosUchwyt } from '../polaczenie/rozdzielacz-zdarzen.ts';
import { oglos } from './ogloszenie.ts';
import {
  chip,
  chwila,
  naKlik,
  nieGotowe,
  panel,
  poproszony,
  przycisk,
  tytul,
  wykonany,
  znacznik,
  type Kontekst,
  type Wiersz,
} from './design-wspolne.ts';
import { dolozPlik, pokazWykaz } from './design-wykaz.ts';
import { pokazAdnotacje, pokazWersje } from './design-adnotacje.ts';

const PN = 'http://www.w3.org/2000/svg';

interface StanPlanszy {
  plansza: DesignBoard | null;
  ramki: DesignFrame[];
  zaznaczone: Set<string>;
  wielokrotne: boolean;
  skok: number;
}

const PLANSZE = new Map<string, StanPlanszy>();

export function stanPlanszy(idKarty: string): StanPlanszy {
  const obecny = PLANSZE.get(idKarty);
  if (obecny !== undefined) return obecny;
  const nowy: StanPlanszy = {
    plansza: null, ramki: [], zaznaczone: new Set(), wielokrotne: false, skok: 8,
  };
  PLANSZE.set(idKarty, nowy);
  return nowy;
}

export function zapomnijPlansze(idKarty: string): void {
  PLANSZE.delete(idKarty);
}

export function zwiazPlansze(kontekst: Kontekst): (() => void)[] {
  const okno = panel(kontekst.korzen, 'panel-board');
  if (okno === null) return [];
  const paski = [...okno.querySelectorAll<HTMLElement>('.dg-narzedzia')];
  const pasekGorny = paski[0] ?? null;
  const pasekDolny = paski.at(-1) ?? null;
  const kanwa = okno.querySelector<HTMLElement>('.dg-kanwa');

  const odswiez = (): void => {
    void wczytaj(kontekst, okno, kanwa);
  };
  kontekst.stan.odswiezenia.set('plansza', odswiez);

  naKlik(kontekst, chip(pasekGorny, 'Warstwy'), () => {
    pokazWarstwy(kontekst);
  });
  naKlik(kontekst, chip(pasekGorny, 'Wyrównaj'), () => wyrownaj(kontekst));
  naKlik(kontekst, chip(pasekGorny, 'Siatka'), () => przestawSiatke(kontekst, okno));
  naKlik(kontekst, chip(pasekGorny, 'Szablon układu'), () => pokazRamki(kontekst));
  naKlik(kontekst, chip(pasekGorny, 'Wersje'), () => pokazWersje(kontekst));
  naKlik(kontekst, chip(pasekGorny, 'Adnotacja'), () => pokazAdnotacje(kontekst));
  naKlik(kontekst, przycisk(okno, 'Eksport'), () => wynies(kontekst));

  const chipWielokrotne = chip(pasekDolny, 'Zaznaczenie wielokrotne');
  naKlik(kontekst, chipWielokrotne, () => {
    const stan = stanPlanszy(kontekst.idKarty);
    stan.wielokrotne = !stan.wielokrotne;
    chipWielokrotne?.setAttribute('aria-pressed', String(stan.wielokrotne));
  });
  naKlik(kontekst, chip(pasekDolny, 'Grupuj'), () => zgrupuj(kontekst));

  const chipObecnosci = chip(pasekDolny, 'edycja współdzielona');
  if (chipObecnosci !== null) opiszObecnosc(chipObecnosci, 0);

  const odlaczenia = [
    zglosUchwyt(EventType.DesignBoardChanged, (zdarzenie) => {
      if (zdarzenie.board.windowId !== kontekst.stan.idOkna) return;
      odswiez();
    }),
    zglosUchwyt(EventType.DesignBoardPresence, (zdarzenie) => {
      if (zdarzenie.boardId !== kontekst.stan.idPlanszy) return;
      opiszObecnosc(chipObecnosci, zdarzenie.participants.length);
    }),
  ];

  odswiez();
  return odlaczenia;
}

function opiszObecnosc(wezel: HTMLElement | null, ilu: number): void {
  if (wezel === null) return;
  const kropka = wezel.querySelector('.dn-kropka');
  wezel.replaceChildren();
  if (kropka !== null) wezel.append(kropka, ' ');
  wezel.append(ilu === 0
    ? 'edycja współdzielona · brak innych kursorów'
    : `edycja współdzielona · ${String(ilu)} kursorów`);
}

async function wczytaj(
  kontekst: Kontekst,
  okno: HTMLElement,
  kanwa: HTMLElement | null,
): Promise<void> {
  if (kontekst.stan.idOkna === '') return;
  const odpowiedz = await poproszony(kontekst.kanal, Command.DesignBoardList, {
    windowId: kontekst.stan.idOkna,
  });
  const stan = stanPlanszy(kontekst.idKarty);
  const plansza = odpowiedz?.boards.find(
    (kandydat) => kontekst.stan.idPlanszy === '' || kandydat.id === kontekst.stan.idPlanszy,
  ) ?? null;
  stan.plansza = plansza;
  if (plansza === null) {
    kontekst.stan.idPlanszy = '';
    tytul(okno, 'Design Board');
    znacznik(okno, '');
    nieGotowe(kanwa, 'Okno nie ma jeszcze kompozycji w rdzeniu.');
    return;
  }
  kontekst.stan.idPlanszy = plansza.id;
  tytul(okno, `Design Board — ${plansza.name ?? plansza.id}`);
  znacznik(okno, `zmieniona ${chwila(plansza.updatedAt)}`);

  const ramki = await poproszony(kontekst.kanal, Command.DesignFrameList, { boardId: plansza.id });
  stan.ramki = ramki?.frames ?? [];
  if (kontekst.stan.idRamki === '' && stan.ramki.length > 0) {
    kontekst.stan.idRamki = stan.ramki[0].id;
  }
  narysuj(kontekst, kanwa, plansza, stan);
}

function narysuj(
  kontekst: Kontekst,
  kanwa: HTMLElement | null,
  plansza: DesignBoard,
  stan: StanPlanszy,
): void {
  if (kanwa === null) return;
  const warstwy = plansza.layers ?? [];
  if (warstwy.length === 0 && stan.ramki.length === 0) {
    nieGotowe(kanwa, 'Kompozycja nie ma jeszcze żadnej warstwy.');
    return;
  }
  const dokument = kanwa.ownerDocument;
  const szerokosc = Math.max(640, ...stan.ramki.map((ramka) => (ramka.x ?? 0) + ramka.width));
  const wysokosc = Math.max(300, ...warstwy.map(
    (warstwa) => (warstwa.y ?? 0) + (warstwa.height ?? 0),
  ));
  const rysunek = dokument.createElementNS(PN, 'svg');
  rysunek.setAttribute('viewBox', `0 0 ${String(szerokosc)} ${String(wysokosc)}`);
  rysunek.setAttribute('role', 'img');
  rysunek.setAttribute('aria-label', `Kompozycja ${plansza.name ?? plansza.id}`);

  for (const ramka of stan.ramki) {
    const prostokat = dokument.createElementNS(PN, 'rect');
    prostokat.setAttribute('x', String(ramka.x ?? 0));
    prostokat.setAttribute('y', String(ramka.y ?? 0));
    prostokat.setAttribute('width', String(ramka.width));
    prostokat.setAttribute('height', String(ramka.height));
    prostokat.setAttribute('rx', '10');
    prostokat.setAttribute('fill', 'var(--dn-powierzchnia)');
    prostokat.setAttribute('stroke', 'var(--dn-obrys)');
    prostokat.setAttribute('stroke-width', '1.5');
    const podpis = dokument.createElementNS(PN, 'text');
    podpis.setAttribute('class', 'dg-etyk');
    podpis.setAttribute('x', String((ramka.x ?? 0) + 8));
    podpis.setAttribute('y', String((ramka.y ?? 0) - 6));
    podpis.textContent = ramka.name;
    rysunek.append(prostokat, podpis);
  }

  const uporzadkowane = [...warstwy].sort((a, b) => (a.order ?? 0) - (b.order ?? 0));
  for (const warstwa of uporzadkowane) {
    rysunek.append(ksztaltWarstwy(kontekst, dokument, warstwa, stan));
  }
  kanwa.replaceChildren(rysunek);
}

function ksztaltWarstwy(
  kontekst: Kontekst,
  dokument: Document,
  warstwa: DesignBoardLayer,
  stan: StanPlanszy,
): SVGGElement {
  const grupa = dokument.createElementNS(PN, 'g');
  grupa.dataset.warstwa = warstwa.id;
  const prostokat = dokument.createElementNS(PN, 'rect');
  prostokat.setAttribute('x', String(warstwa.x ?? 0));
  prostokat.setAttribute('y', String(warstwa.y ?? 0));
  prostokat.setAttribute('width', String(warstwa.width ?? 0));
  prostokat.setAttribute('height', String(warstwa.height ?? 0));
  prostokat.setAttribute('rx', '6');
  prostokat.setAttribute('fill', 'var(--dn-powierzchnia-2)');
  prostokat.setAttribute(
    'stroke',
    stan.zaznaczone.has(warstwa.id) ? 'var(--dn-sygnal-obrys)' : 'var(--dn-obrys)',
  );
  prostokat.setAttribute('stroke-width', stan.zaznaczone.has(warstwa.id) ? '2' : '1.5');
  grupa.append(prostokat);
  if (warstwa.note !== undefined && warstwa.note !== '') {
    const opis = dokument.createElementNS(PN, 'text');
    opis.setAttribute('class', 'dg-etyk');
    opis.setAttribute('x', String((warstwa.x ?? 0) + 8));
    opis.setAttribute('y', String((warstwa.y ?? 0) + (warstwa.height ?? 0) - 8));
    opis.textContent = warstwa.note;
    grupa.append(opis);
  }
  grupa.addEventListener('click', () => {
    if (!stan.wielokrotne) stan.zaznaczone.clear();
    if (stan.zaznaczone.has(warstwa.id)) stan.zaznaczone.delete(warstwa.id);
    else stan.zaznaczone.add(warstwa.id);
    kontekst.stan.odswiezenia.get('plansza')?.();
  }, kontekst.przy);
  if (warstwa.locked !== true) zwiazPrzesuwanie(kontekst, grupa, warstwa, stan);
  return grupa;
}

// Rdzeń nie zna przeciągania cząstkowego: położenie idzie dopiero po puszczeniu.
function zwiazPrzesuwanie(
  kontekst: Kontekst,
  grupa: SVGGElement,
  warstwa: DesignBoardLayer,
  stan: StanPlanszy,
): void {
  let start: { x: number; y: number } | null = null;
  grupa.addEventListener('pointerdown', (zdarzenie) => {
    start = { x: zdarzenie.clientX, y: zdarzenie.clientY };
  }, kontekst.przy);
  grupa.addEventListener('pointerup', (zdarzenie) => {
    if (start === null) return;
    const przesuniecieX = Math.round((zdarzenie.clientX - start.x) / stan.skok) * stan.skok;
    const przesuniecieY = Math.round((zdarzenie.clientY - start.y) / stan.skok) * stan.skok;
    start = null;
    if (przesuniecieX === 0 && przesuniecieY === 0) return;
    void przesun(kontekst, warstwa.id, przesuniecieX, przesuniecieY);
  }, kontekst.przy);
}

async function przesun(
  kontekst: Kontekst,
  idWarstwy: string,
  dx: number,
  dy: number,
): Promise<void> {
  const stan = stanPlanszy(kontekst.idKarty);
  const plansza = stan.plansza;
  if (plansza === null) return;
  const warstwy = (plansza.layers ?? []).map((warstwa) => (warstwa.id === idWarstwy
    ? { ...warstwa, x: (warstwa.x ?? 0) + dx, y: (warstwa.y ?? 0) + dy }
    : warstwa));
  await wykonany(kontekst.kanal, Command.DesignBoardUpdate, {
    windowId: kontekst.stan.idOkna,
    boardId: plansza.id,
    layers: warstwy,
  }, 'Przesunięcie warstwy');
}

export async function dolozZasobDoPlanszy(kontekst: Kontekst, idZasobu: string): Promise<void> {
  if (idZasobu === '') {
    oglos('Design', 'Najpierw wskaż zasób w Assets Panel.', 'ostrzezenie');
    return;
  }
  const stan = stanPlanszy(kontekst.idKarty);
  const warstwy = [...(stan.plansza?.layers ?? [])];
  warstwy.push({
    id: `warstwa-${String(warstwy.length + 1)}`,
    assetId: idZasobu,
    x: 40,
    y: 40,
    width: 240,
    height: 150,
    order: warstwy.length,
  });
  const zadanie = {
    windowId: kontekst.stan.idOkna,
    ...(stan.plansza === null ? {} : { boardId: stan.plansza.id }),
    layers: warstwy,
  };
  await wykonany(kontekst.kanal, Command.DesignBoardUpdate, zadanie, 'Dołożenie warstwy');
  kontekst.stan.odswiezenia.get('plansza')?.();
}

function pokazWarstwy(kontekst: Kontekst): void {
  const stan = stanPlanszy(kontekst.idKarty);
  const warstwy = stan.plansza?.layers ?? [];
  const wiersze: Wiersz[] = warstwy.map((warstwa) => ({
    tekst: warstwa.note ?? warstwa.id,
    meta: `${String(warstwa.width ?? 0)}×${String(warstwa.height ?? 0)}`,
    plakietka: warstwa.locked === true ? 'zablokowana' : '',
    kropka: stan.zaznaczone.has(warstwa.id) ? 'sygnal' : 'neutralna',
    naKlik: () => {
      if (!stan.wielokrotne) stan.zaznaczone.clear();
      stan.zaznaczone.add(warstwa.id);
      kontekst.stan.odswiezenia.get('plansza')?.();
    },
  }));
  pokazWykaz(kontekst, 'Warstwy kompozycji', wiersze, 'Kompozycja nie ma warstw.');
}

async function pokazRamki(kontekst: Kontekst): Promise<void> {
  if (kontekst.stan.idPlanszy === '') return;
  const odpowiedz = await poproszony(kontekst.kanal, Command.DesignFrameList, {
    boardId: kontekst.stan.idPlanszy,
  });
  const ramki = odpowiedz?.frames ?? [];
  const wiersze: Wiersz[] = ramki.map((ramka) => ({
    tekst: ramka.name,
    meta: `${String(ramka.width)}×${String(ramka.height)}`,
    plakietka: ramka.devicePreset ?? '',
    kropka: ramka.id === kontekst.stan.idRamki ? 'sygnal' : 'neutralna',
    naKlik: () => {
      kontekst.stan.idRamki = ramka.id;
      void zastosujRozmiar(kontekst, ramka);
    },
  }));
  pokazWykaz(kontekst, 'Ramki kompozycji', wiersze, 'Kompozycja nie ma ramek.');
}

async function zastosujRozmiar(kontekst: Kontekst, ramka: DesignFrame): Promise<void> {
  await wykonany(kontekst.kanal, Command.DesignFrameResizeApply, {
    frameId: ramka.id,
    width: ramka.width,
    height: ramka.height,
  }, 'Przeliczenie więzi ramki');
  await wykonany(kontekst.kanal, Command.DesignConstraintSet, {
    frameId: ramka.id,
    constraints: [],
  }, 'Więzi ramki');
}

async function wyrownaj(kontekst: Kontekst): Promise<void> {
  const stan = stanPlanszy(kontekst.idKarty);
  if (kontekst.stan.idRamki === '') {
    oglos('Design', 'Wyrównanie działa w ramce — wskaż ją w „Szablon układu".', 'ostrzezenie');
    return;
  }
  await wykonany(kontekst.kanal, Command.DesignLayoutAuto, {
    frameId: kontekst.stan.idRamki,
    layout: { direction: DesignLayoutDirection.Vertical },
    layerIds: [...stan.zaznaczone],
  }, 'Układ automatyczny');
  kontekst.stan.odswiezenia.get('plansza')?.();
}

async function przestawSiatke(kontekst: Kontekst, okno: HTMLElement): Promise<void> {
  const stan = stanPlanszy(kontekst.idKarty);
  stan.skok = stan.skok === 8 ? 16 : 8;
  if (kontekst.stan.idPlanszy === '') return;
  await wykonany(kontekst.kanal, Command.DesignGridSet, {
    grid: { columns: 12, gutter: stan.skok * 2, margin: stan.skok * 3, baseline: stan.skok },
    boardId: kontekst.stan.idPlanszy,
  }, 'Siatka kompozycji');
  znacznik(okno, `siatka: 12 kolumn, linia bazowa ${String(stan.skok)} px`);
}

async function zgrupuj(kontekst: Kontekst): Promise<void> {
  const stan = stanPlanszy(kontekst.idKarty);
  if (kontekst.stan.idPlanszy === '' || stan.zaznaczone.size === 0) {
    oglos('Design', 'Zaznacz warstwy, które mają wejść w symbol.', 'ostrzezenie');
    return;
  }
  await wykonany(kontekst.kanal, Command.DesignVectorSymbolSet, {
    boardId: kontekst.stan.idPlanszy,
    name: `symbol-${String(stan.zaznaczone.size)}`,
    layerIds: [...stan.zaznaczone],
  }, 'Grupowanie w symbol');
}

async function wynies(kontekst: Kontekst): Promise<void> {
  if (kontekst.stan.idPlanszy === '') return;
  const odpowiedz = await wykonany(kontekst.kanal, Command.DesignBoardExport, {
    boardId: kontekst.stan.idPlanszy,
    format: 'pdf',
  }, 'Wyniesienie kompozycji');
  if (odpowiedz === null) return;
  dolozPlik(kontekst, {
    fileName: odpowiedz.fileName,
    mediaType: odpowiedz.mediaType,
    contentBase64: odpowiedz.contentBase64,
  });
}
