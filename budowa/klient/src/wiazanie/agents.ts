// Wiązanie karty modułu Agents z rdzeniem: biblioteka ekspertów, edytor
// definicji oraz panele modelu, umiejętności, rozszerzeń, zakresu i wersji.

import {
  ChangeKind,
  Command,
  EventType,
  type Agent,
} from '../../../shared/contract.ts';
import type { Odsubskrybuj } from '../polaczenie/magistrala-zdarzen.ts';
import { zglosUchwyt } from '../polaczenie/rozdzielacz-zdarzen.ts';
import type { Kanal } from '../protokol/kanal.ts';
import { wywolaj } from '../protokol/wywolanie.ts';
import { oglos } from './ogloszenie.ts';
import {
  uzgodnijPrzelacznikiPaneli,
  zdejmijSterowanieWspolne,
  zdejmijTrescWspolna,
} from './okno-modulu.ts';
import { zwiazBiblioteke } from './agents-biblioteka.ts';
import { zwiazKonektory } from './agents-konektory.ts';
import { zwiazModel } from './agents-model.ts';
import { zwiazTozsamosc } from './agents-tozsamosc.ts';
import { zwiazUmiejetnosci } from './agents-umiejetnosci.ts';
import { zwiazUprawnienia } from './agents-uprawnienia.ts';
import { zwiazWersje } from './agents-wersje.ts';
import { zwiazZespoly } from './agents-zespoly.ts';
import { pustka, tekst, type StanowiskoEkspertow } from './agents-wspolne.ts';

interface WiazanieKarty {
  korzen: Element;
  odlaczenia: Odsubskrybuj[];
}

// Każda karta ma własny ekspert bieżący; zamknięcie zwalnia wyłącznie jej wiązanie.
const WIAZANIA = new Map<string, WiazanieKarty>();

export function zwiazAgents(
  kanal: Kanal,
  nazwaSrodowiska: string,
  idOknaStojacego: string,
  wskazanieKorzenia: Element | string,
): boolean {
  const korzenMoze = korzenKarty(wskazanieKorzenia);
  if (korzenMoze === null) return false;
  const korzen = korzenMoze;
  const idKarty = korzen.getAttribute('data-karta') ?? '';
  if (idKarty === '') return false;
  if (WIAZANIA.get(idKarty)?.korzen === korzen) return false;
  zwolnijAgents(idKarty);

  const odlaczenia: Odsubskrybuj[] = [];
  const sterowanie = new AbortController();
  const przy = { signal: sterowanie.signal };
  odlaczenia.push(() => {
    sterowanie.abort();
  });
  WIAZANIA.set(idKarty, { korzen, odlaczenia });

  if (korzen instanceof HTMLElement) {
    zdejmijTrescWspolna(korzen);
    zdejmijSterowanieWspolne(korzen);
    uzgodnijPrzelacznikiPaneli(korzen, []);
  }
  zdejmijTrescPrzykladowa(korzen, nazwaSrodowiska);

  let ekspert: Agent | null = null;

  const stan: StanowiskoEkspertow = {
    kanal,
    korzen,
    idOkna: idOknaStojacego,
    przy,
    ekspert: () => ekspert,
    postaw: (agent: Agent) => {
      ekspert = agent;
      pokazEksperta();
    },
    odswiez: () => {
      biblioteka?.odswiez();
    },
  };

  const tozsamosc = zwiazTozsamosc(stan);
  const model = zwiazModel(stan);
  const umiejetnosci = zwiazUmiejetnosci(stan);
  const konektory = zwiazKonektory(stan);
  const uprawnienia = zwiazUprawnienia(stan);
  const wersje = zwiazWersje(stan, otworzEksperta);
  const panele = zwiazZespoly(stan);
  const biblioteka = zwiazBiblioteke(stan, otworzEksperta);

  function pokazEksperta(): void {
    tozsamosc?.pokaz(ekspert);
    model?.pokaz(ekspert);
    umiejetnosci?.pokaz(ekspert);
    konektory?.pokaz(ekspert);
    uprawnienia?.pokaz(ekspert);
    wersje?.pokaz(ekspert);
    panele?.pokaz(ekspert);
    opiszPodsumowanie(korzen, ekspert);
  }

  function otworzEksperta(agent: Agent): void {
    ekspert = agent;
    przelaczWidok(korzen, 'widok-edytor');
    pokazEksperta();
  }

  async function zalozEksperta(): Promise<void> {
    const wynik = await wywolaj(kanal, Command.AgentCreate, { name: 'Nowy ekspert' });
    if (!wynik.udany || wynik.wynik === undefined) {
      oglos('Nowy ekspert', wynik.blad?.message ?? 'Rdzeń odrzucił założenie.', 'blad');
      return;
    }
    otworzEksperta(wynik.wynik.agent);
    biblioteka?.odswiez();
  }

  korzen.addEventListener('click', (zdarzenie) => {
    const cel = zdarzenie.target;
    if (!(cel instanceof Element)) return;
    if (cel.closest('[data-nowy-ekspert]') !== null) {
      void zalozEksperta();
      return;
    }
    if (cel.closest('[data-do-biblioteki]') !== null) {
      zdarzenie.preventDefault();
      przelaczWidok(korzen, 'widok-biblioteka');
      biblioteka?.odswiez();
      return;
    }
    const przelacznikWidoku = cel.closest<HTMLElement>('[data-widok]');
    if (przelacznikWidoku?.dataset.widok !== undefined) {
      przelaczWidok(korzen, przelacznikWidoku.dataset.widok);
      return;
    }
    const skok = cel.closest<HTMLElement>('[data-jump], [data-jump-panel]');
    const idPanelu = skok?.dataset.jump ?? skok?.dataset.jumpPanel;
    if (idPanelu !== undefined) korzen.querySelector(`#${idPanelu}`)?.removeAttribute('hidden');
    const przelacznikTabu = cel.closest<HTMLElement>('[data-tab]');
    if (przelacznikTabu?.dataset.tab !== undefined) {
      przelaczZakladke(korzen, przelacznikTabu.dataset.tab);
    }
  }, przy);

  odlaczenia.push(zglosUchwyt(EventType.AgentChanged, (tresc) => {
    if (tresc.change === ChangeKind.Deleted && tresc.agent.id === ekspert?.id) {
      ekspert = null;
      pokazEksperta();
    } else if (tresc.agent.id === ekspert?.id) {
      ekspert = tresc.agent;
      pokazEksperta();
    }
    biblioteka?.odswiez();
  }));

  odlaczenia.push(zglosUchwyt(EventType.TeamChanged, () => {
    panele?.pokaz(ekspert);
  }));

  odlaczenia.push(zglosUchwyt(EventType.SubagentChanged, () => {
    panele?.odswiezPodagentow();
  }));

  pokazEksperta();
  biblioteka?.odswiez();
  return true;
}

export function zwolnijAgents(idKarty: string): void {
  const wiazanie = WIAZANIA.get(idKarty);
  if (wiazanie === undefined) return;
  for (const odlacz of wiazanie.odlaczenia) odlacz();
  WIAZANIA.delete(idKarty);
}

function korzenKarty(wskazanie: Element | string): Element | null {
  if (typeof wskazanie !== 'string') return wskazanie;
  return document.querySelector(`.cd-tresc--modul[data-karta="${wskazanie}"]`);
}

function przelaczWidok(korzen: Element, idWidoku: string): void {
  for (const widok of korzen.querySelectorAll('.ab-widok')) {
    if (widok.id === idWidoku) widok.removeAttribute('hidden');
    else widok.setAttribute('hidden', '');
  }
  for (const przelacznik of korzen.querySelectorAll<HTMLElement>('[data-widok]')) {
    przelacznik.setAttribute('aria-selected', String(przelacznik.dataset.widok === idWidoku));
  }
}

function przelaczZakladke(korzen: Element, idZakladki: string): void {
  for (const zakladka of korzen.querySelectorAll('.ab-tab')) {
    if (zakladka.id === idZakladki) zakladka.removeAttribute('hidden');
    else zakladka.setAttribute('hidden', '');
  }
  for (const przelacznik of korzen.querySelectorAll<HTMLElement>('.ab-zakl [data-tab]')) {
    przelacznik.setAttribute('aria-selected', String(przelacznik.dataset.tab === idZakladki));
  }
  for (const krok of korzen.querySelectorAll<HTMLElement>('.ab-cykl-krok[data-tab]')) {
    if (krok.dataset.tab === idZakladki) krok.setAttribute('aria-current', 'step');
    else krok.removeAttribute('aria-current');
  }
}

// Podsumowanie definicji zbiera to, co rdzeń oddał o ekspercie; wiersz bez
// pokrycia zostaje pusty, bo jego wartość z prototypu byłaby wymyślona.
function opiszPodsumowanie(korzen: Element, agent: Agent | null): void {
  const podsumowanie = korzen.querySelector('.ab-sum:not(.ab-historia)');
  if (podsumowanie === null) return;
  const wartosci: Record<string, string> = agent === null ? {} : {
    Model: agent.model ?? '',
    'Kanał': agent.transport ?? '',
    'Umiejętności': String(agent.skillIds?.length ?? 0),
    Rozszerzenia: String((agent.connectorIds?.length ?? 0) + (agent.pluginIds?.length ?? 0)),
    Uprawnienia: String(agent.permissions?.filter((wpis) => !wpis.granted).length ?? 0),
    'Pamięć': agent.memoryLevels.join(', '),
  };
  for (const wiersz of podsumowanie.querySelectorAll('.ab-sum-w')) {
    const podpis = (wiersz.querySelector('span')?.textContent ?? '').trim();
    tekst(wiersz.querySelector('b'), wartosci[podpis] ?? '');
  }
}

function zdejmijTrescPrzykladowa(korzen: Element, nazwaSrodowiska: string): void {
  for (const baner of korzen.querySelectorAll('.ab-tab .pk-baner')) baner.textContent = '';
  const historia = korzen.querySelector('.sta-kom-historia');
  historia?.replaceChildren();
  if (historia !== null) {
    pustka(historia, `${nazwaSrodowiska} — rozmowa modułu nie jest jeszcze związana z rdzeniem.`);
  }
  korzen.querySelector('.ab-sum:not(.ab-historia) .ab-sum-tytul')?.removeAttribute('hidden');
  for (const nazwa of korzen.querySelectorAll('.pt-pozycja-tytul')) {
    if ((nazwa.textContent ?? '').includes('Agent Redaktor')) nazwa.textContent = '';
  }
}
