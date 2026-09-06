/**
 * Przedsionek środowiska MultitaskingAI. Miary kafli bierze z rejestrów rdzenia
 * — `team.list`, `agent.list`, `channel.list` i `queue.list` — a wskazanie kafla
 * rozstrzyga wykazem `module.list`: sekcja bez modułu melduje brak okna zamiast
 * milczeć.
 */

import { Command, type Queue } from '../../../shared/contract.ts';
import type { Kanal, Wynik } from '../protokol/kanal.ts';
import { wywolaj } from '../protokol/wywolanie.ts';
import { oglos } from './ogloszenie.ts';
import { zdanieOdmowy } from './wejscie-odmowa.ts';

const NAGLOWEK = 'MultitaskingAI';
const KOD_SRODOWISKA = 'multitaskingai';

interface MiaryZespolu {
  zespoly: number | null;
  eksperci: number | null;
  kanaly: number | null;
  kolejki: Queue[] | null;
}

let siatkaWypelniona: HTMLElement | null = null;
let kodyModulow: Set<string> | null = null;

/* Widok przedsionka powstaje w płótnie przy każdym wejściu w środowisko, więc
   wiązanie idzie za jego przebudową obserwatorem, nie odczytem znacznika. */
export function zwiazPrzedsionekMultitaskingAI(kanal: Kanal): void {
  sledzWidok(kanal);
  document.addEventListener('click', (zdarzenie) => {
    const cel = zdarzenie.target;
    if (!(cel instanceof Element) || siatkaSrodowiska() === null) return;
    const kafel = cel.closest<HTMLElement>('.pd-kafel[data-modul]');
    if (kafel !== null) opiszWejscieDoSekcji(kafel);
  }, true);
}

function sledzWidok(kanal: Kanal): void {
  const obserwator = new MutationObserver(() => {
    const siatka = siatkaSrodowiska();
    if (siatka === null) {
      siatkaWypelniona = null;
      return;
    }
    if (siatka === siatkaWypelniona) return;
    siatkaWypelniona = siatka;
    void wypelnij(kanal, siatka);
  });
  obserwator.observe(document.body, { childList: true, subtree: true });
}

/* Cztery przedsionki dzielą jeden węzeł płótna, a rozpoznaje je dopiero kod
   środowiska niesiony przez listwę działań; bez tej kontroli wiązanie sięgałoby
   po kafle środowiska sąsiedniego. */
function siatkaSrodowiska(): HTMLElement | null {
  const widok = document.querySelector<HTMLElement>('.cd-tresc--przedsionek');
  if (widok === null || widok.hidden) return null;
  const listwa = widok.querySelector<HTMLElement>('.pd-listwa [data-nowy-projekt-srodowisko]');
  if (listwa?.dataset.nowyProjektSrodowisko !== KOD_SRODOWISKA) return null;
  return widok.querySelector<HTMLElement>('.pd-siatka');
}

async function wypelnij(kanal: Kanal, siatka: HTMLElement): Promise<void> {
  opiszKafle(siatka, await pobierzMiary(kanal));
  kodyModulow = await pobierzKodyModulow(kanal);
}

async function pobierzMiary(kanal: Kanal): Promise<MiaryZespolu> {
  const [zespoly, eksperci, kanaly, kolejki] = await Promise.all([
    wywolaj(kanal, Command.TeamList, {}),
    wywolaj(kanal, Command.AgentList, { enabledOnly: true }),
    wywolaj(kanal, Command.ChannelList, { enabledOnly: true }),
    wywolaj(kanal, Command.QueueList, {}),
  ]);
  const odmowy: string[] = [];
  const miary: MiaryZespolu = {
    zespoly: przyjmij(zespoly, 'zespołów', odmowy)?.teams.length ?? null,
    eksperci: przyjmij(eksperci, 'ekspertów', odmowy)?.agents.length ?? null,
    kanaly: przyjmij(kanaly, 'kanałów modelu', odmowy)?.channels.length ?? null,
    kolejki: przyjmij(kolejki, 'kolejek', odmowy)?.queues ?? null,
  };
  if (odmowy.length > 0) oglos(NAGLOWEK, zdanieOdmow(odmowy), 'ostrzezenie');
  return miary;
}

function przyjmij<T>(wynik: Wynik<T>, rejestr: string, odmowy: string[]): T | null {
  if (wynik.udany && wynik.wynik !== undefined) return wynik.wynik;
  odmowy.push(`${rejestr} — ${zdanieOdmowy(wynik.blad)}`);
  return null;
}

function zdanieOdmow(odmowy: string[]): string {
  return `Rdzeń nie wydał wykazu ${odmowy.join('; ')}`;
}

/* Kafel bez miary z rdzenia traci swoją miarę prototypową: liczba wpisana
   w prototypie mówiłaby o zespołach i zleceniach, których w bazie nie ma. */
function opiszKafle(siatka: HTMLElement, miary: MiaryZespolu): void {
  for (const kafel of siatka.querySelectorAll<HTMLElement>('.pd-kafel')) {
    const meta = kafel.querySelector('.pd-kafel-meta');
    if (meta === null) continue;
    const zdanie = miaraKafla(kafel.dataset.miaraKafla ?? '', miary);
    if (zdanie === '') meta.remove();
    else meta.textContent = zdanie;
  }
}

function miaraKafla(rodzaj: string, miary: MiaryZespolu): string {
  if (rodzaj === 'zespoly') return liczebnik(miary.zespoly, 'zespół', 'zespoły', 'zespołów');
  if (rodzaj === 'eksperci') return liczebnik(miary.eksperci, 'ekspert', 'ekspertów', 'ekspertów');
  if (rodzaj === 'kanaly') {
    return liczebnik(miary.kanaly, 'kanał modelu', 'kanały modelu', 'kanałów modelu');
  }
  if (rodzaj === 'kolejki') return miaraKolejek(miary.kolejki);
  return '';
}

/* Sama liczba kolejek nie mówi o obciążeniu zespołu, więc miara niesie obok
   niej sumę zleceń czekających na podjęcie. */
function miaraKolejek(kolejki: Queue[] | null): string {
  if (kolejki === null) return '';
  const wykaz = liczebnik(kolejki.length, 'kolejka', 'kolejki', 'kolejek');
  if (kolejki.length === 0) return wykaz;
  const oczekuje = kolejki.reduce((suma, kolejka) => suma + (kolejka.pendingCount ?? 0), 0);
  return `${wykaz} · ${oczekuje} w oczekiwaniu`;
}

/* Polszczyzna odmienia rzeczownik po liczbie w trzech postaciach. */
function liczebnik(ile: number | null, jedna: string, kilka: string, wiele: string): string {
  if (ile === null) return '';
  if (ile === 1) return `1 ${jedna}`;
  const reszta = ile % 10;
  const setka = ile % 100;
  const mnoga = reszta >= 2 && reszta <= 4 && (setka < 12 || setka > 14);
  return `${ile} ${mnoga ? kilka : wiele}`;
}

async function pobierzKodyModulow(kanal: Kanal): Promise<Set<string> | null> {
  const wynik = await wywolaj(kanal, Command.ModuleList, {});
  if (!wynik.udany || wynik.wynik === undefined) return null;
  return new Set(wynik.wynik.modules.map((modul) => modul.code.toLowerCase()));
}

/* MultitaskingAI prowadzi się sekcjami sterowania zespołem, a rejestr modułów
   zna spośród nich tylko część; kafel bez modułu zostawał dotąd bez odpowiedzi,
   więc Operator nie wiedział, czy wskazanie w ogóle doszło. */
function opiszWejscieDoSekcji(kafel: HTMLElement): void {
  if (kodyModulow === null) return;
  const kod = (kafel.dataset.modul ?? '').toLowerCase();
  if (kod === '' || kodyModulow.has(kod)) return;
  const nazwa = kafel.querySelector('.pd-kafel-nazwa')?.textContent?.trim() ?? '';
  const sekcja = nazwa === '' ? 'Ta sekcja' : `Sekcja „${nazwa}”`;
  oglos(NAGLOWEK, `${sekcja} nie ma własnego okna w tym wydaniu.`, 'ostrzezenie');
}
