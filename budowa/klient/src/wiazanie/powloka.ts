/* Znacznik powłoki niesie `design/zasoby/powloka.js`; szyna powstaje tu
   z pozycji wzorcowej środowiska i modułu, powielanych wykazem rejestru. */

import { Command, SessionStatus } from '../../../shared/contract.ts';
import type { Environment, Module } from '../../../shared/contract.ts';
import type { Kanal } from '../protokol/kanal.ts';
import { wywolaj } from '../protokol/wywolanie.ts';

export interface DanePowloki {
  login?: string;
  nazwaOkna?: string;
  kodSrodowiska?: string;
  kodModulu?: string;
}

export interface RozliczeniePowloki {
  zwiazana: boolean;
  srodowiskaRdzenia: number;
  srodowiskaZnacznika: number;
  ukrytePozycje: number;
  bezPozycji: string[];
  bledy: string[];
}

/** Kanał wystawiony przez `aplikacja.ts`; skrypty biblioteki nie są modułami, więc importu z nich nie ma. */
function kanalGlobalny(): Kanal | undefined {
  return (globalThis as { DanacoKanal?: Kanal }).DanacoKanal;
}

export async function zwiazPowloke(
  dane: DanePowloki = {},
  kanalWskazany?: Kanal,
): Promise<RozliczeniePowloki> {
  const rozliczenie: RozliczeniePowloki = {
    zwiazana: false,
    srodowiskaRdzenia: 0,
    srodowiskaZnacznika: 0,
    ukrytePozycje: 0,
    bezPozycji: [],
    bledy: [],
  };

  const powloka = document.getElementById('dn-powloka-montaz');
  const kanal = kanalWskazany ?? kanalGlobalny();
  if (powloka === null || kanal === undefined) return rozliczenie;

  const [srodowiska, moduly, sesje] = await Promise.all([
    wywolaj(kanal, Command.EnvironmentList, { includeModules: true }),
    wywolaj(kanal, Command.ModuleList, {}),
    wywolaj(kanal, Command.SessionList, { status: SessionStatus.Active }),
  ]);

  for (const wynik of [srodowiska, moduly, sesje]) {
    if (!wynik.udany && wynik.blad !== undefined) rozliczenie.bledy.push(wynik.blad.message);
  }

  const wykazSrodowisk = uporzadkuj(srodowiska.udany ? (srodowiska.wynik?.environments ?? []) : []);
  const wykazModulow = moduly.udany ? (moduly.wynik?.modules ?? []) : [];
  const katalogModulow = new Map(wykazModulow.map((modul) => [modul.code, modul]));

  rozliczenie.srodowiskaRdzenia = wykazSrodowisk.length;
  wypelnijSzyne(powloka, wykazSrodowisk, katalogModulow, rozliczenie);

  const stan = powloka.querySelector('.dn-stan');
  if (stan !== null) {
    usunMiaryBezZrodla(stan);
    wpiszKonto(stan, dane.login);
    wpiszPolozenie(stan, wykazSrodowisk, katalogModulow, dane);
    wpiszSesje(stan, sesje.udany ? (sesje.wynik?.total ?? 0) : 0);
  }

  if (dane.nazwaOkna !== undefined && dane.nazwaOkna.length > 0) {
    const belka = powloka.querySelector('.dn-belka-tytul');
    if (belka !== null) belka.textContent = dane.nazwaOkna;
  }

  rozliczenie.zwiazana = true;
  return rozliczenie;
}

function uporzadkuj(wykaz: Environment[]): Environment[] {
  return [...wykaz].sort((a, b) => a.order - b.order);
}

interface WzorySzyny {
  lista: HTMLElement;
  srodowisko: HTMLElement;
  grupa: HTMLElement;
  modul: HTMLElement;
  znaki: Map<string, Element>;
}

/* Znak jest jedyną cechą pozycji, której rejestr nie niesie: `Module.icon` podaje
   kod, a zestawu znaków klient nie ma. Znak bierze się więc z pozycji znacznika
   o tym samym kodzie. */
function zdejmijWzory(powloka: HTMLElement): WzorySzyny | null {
  const srodowisko = powloka.querySelector<HTMLElement>('.dn-szyna-poz--srodowisko');
  const lista = srodowisko?.parentElement ?? null;
  const grupa = document.getElementById(srodowisko?.getAttribute('aria-controls') ?? '');
  const modul = grupa?.querySelector<HTMLElement>('.dn-szyna-poz--modul') ?? null;
  if (srodowisko === null || lista === null || grupa === null || modul === null) return null;
  const znaki = new Map<string, Element>();
  for (const pozycja of powloka.querySelectorAll<HTMLElement>('.dn-szyna-poz--srodowisko, .dn-szyna-poz--modul')) {
    const kod = (pozycja.dataset.srodowisko ?? pozycja.dataset.modul ?? '').toLowerCase();
    const znak = pozycja.querySelector(':scope > svg');
    if (kod !== '' && znak !== null && !znaki.has(kod)) znaki.set(kod, znak.cloneNode(true) as Element);
  }
  return {
    lista,
    srodowisko: srodowisko.cloneNode(true) as HTMLElement,
    grupa: grupa.cloneNode(true) as HTMLElement,
    modul: modul.cloneNode(true) as HTMLElement,
    znaki,
  };
}

/* Lista pustoszeje przed wypełnieniem: pozycja znacznika, której rejestr nie
   potwierdza, byłaby drogą do środowiska spoza platformy. */
function wypelnijSzyne(
  powloka: HTMLElement,
  wykazSrodowisk: Environment[],
  katalogModulow: Map<string, Module>,
  rozliczenie: RozliczeniePowloki,
): void {
  const wzory = zdejmijWzory(powloka);
  if (wzory === null) return;
  wzory.lista.replaceChildren();
  for (const srodowisko of wykazSrodowisk) {
    if (srodowisko.code === '') {
      rozliczenie.ukrytePozycje += 1;
      continue;
    }
    const oznaczenie = 'moduly-' + srodowisko.code.toLowerCase();
    const pozycja = wzory.srodowisko.cloneNode(true) as HTMLElement;
    pozycja.removeAttribute('hidden');
    pozycja.dataset.srodowisko = srodowisko.code;
    pozycja.setAttribute('aria-controls', oznaczenie);
    pozycja.setAttribute('aria-expanded', 'false');
    wstawZnak(pozycja, wzory.znaki.get(srodowisko.code.toLowerCase()));
    opiszPozycje(pozycja, srodowisko.name, srodowisko.description);
    const grupa = wzory.grupa.cloneNode(true) as HTMLElement;
    grupa.id = oznaczenie;
    grupa.setAttribute('hidden', '');
    grupa.setAttribute('aria-label', `Moduły środowiska ${srodowisko.name}`);
    grupa.replaceChildren();
    wypelnijModuly(grupa, wzory, srodowisko, katalogModulow, rozliczenie);
    wzory.lista.appendChild(pozycja);
    wzory.lista.appendChild(grupa);
    rozliczenie.srodowiskaZnacznika += 1;
  }
}

function wypelnijModuly(
  grupa: HTMLElement,
  wzory: WzorySzyny,
  srodowisko: Environment,
  katalogModulow: Map<string, Module>,
  rozliczenie: RozliczeniePowloki,
): void {
  for (const kod of srodowisko.moduleCodes ?? []) {
    const modul = katalogModulow.get(kod);
    if (modul === undefined) {
      rozliczenie.bezPozycji.push(`${srodowisko.code}/${kod}`);
      continue;
    }
    const pozycja = wzory.modul.cloneNode(true) as HTMLElement;
    pozycja.removeAttribute('hidden');
    pozycja.dataset.modul = modul.code;
    wstawZnak(pozycja, wzory.znaki.get(modul.code.toLowerCase()));
    opiszPozycje(pozycja, modul.name, modul.description);
    grupa.appendChild(pozycja);
  }
}

function wstawZnak(pozycja: HTMLElement, znak: Element | undefined): void {
  const stojacy = pozycja.querySelector(':scope > svg');
  if (znak === undefined || stojacy === null) return;
  stojacy.replaceWith(znak.cloneNode(true));
}

function opiszPozycje(pozycja: Element, nazwa: string, opis: string | undefined): void {
  pozycja.setAttribute('aria-label', nazwa);
  const etykieta = pozycja.querySelector('.dn-szyna-etyk');
  if (etykieta === null) return;
  const tytul = etykieta.querySelector('b');
  if (tytul !== null) tytul.textContent = nazwa;
  const podpis = etykieta.querySelector('span');
  if (podpis !== null) podpis.textContent = opis ?? '';
}

/** Miara bez źródła w kontrakcie jest atrapą; znika wraz z rozdzielnikiem przed nią. */
function usunMiaryBezZrodla(stan: Element): void {
  for (const znacznik of ['[data-stan-cpu]', '[data-stan-ram]']) {
    const miara = stan.querySelector(znacznik);
    const pozycja = miara === null ? null : miara.closest('.dn-stan-poz');
    if (pozycja === null) continue;
    const rozdzielnik = pozycja.previousElementSibling;
    if (rozdzielnik !== null && rozdzielnik.classList.contains('dn-stan-sep')) rozdzielnik.remove();
    pozycja.remove();
  }
}

function wpiszKonto(stan: Element, login: string | undefined): void {
  if (login === undefined || login.length === 0) return;
  const pozycja = stan.querySelector('.dn-stan-poz');
  const rola = pozycja === null ? null : pozycja.querySelector('b');
  const tekst = rola === null ? null : rola.nextSibling;
  if (tekst !== null) tekst.nodeValue = ` · ${login}`;
}

function wpiszPolozenie(
  stan: Element,
  wykazSrodowisk: Environment[],
  katalogModulow: Map<string, Module>,
  dane: DanePowloki,
): void {
  const srodowisko = wykazSrodowisk.find((pozycja) => pozycja.code === dane.kodSrodowiska);
  if (srodowisko === undefined) return;
  const pozycja = stan.querySelector('.dn-stan-poz--drugorzedna');
  const nazwa = pozycja === null ? null : pozycja.querySelector('b');
  if (nazwa === null) return;
  nazwa.textContent = srodowisko.name;
  const modul = dane.kodModulu === undefined ? undefined : katalogModulow.get(dane.kodModulu);
  const tekst = nazwa.nextSibling;
  if (tekst !== null && modul !== undefined) tekst.nodeValue = ` · ${modul.name}`;
}

function wpiszSesje(stan: Element, liczba: number): void {
  const licznik = stan.querySelector('[data-stan-sesje]');
  if (licznik !== null) licznik.textContent = String(liczba);
}
