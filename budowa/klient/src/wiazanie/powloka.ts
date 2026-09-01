/**
 * Wiązanie powłoki Właściciela z rdzeniem. Znacznik stawia
 * `design/zasoby/powloka.js`; tutaj wyłącznie wypełniane są istniejące węzły
 * wartościami z kontraktu i chowane pozycje, dla których rdzeń nie ma pokrycia.
 * Żaden element nie powstaje po tej stronie.
 */

import { Command, SessionStatus } from '../../../shared/contract.ts';
import type { Environment, Module } from '../../../shared/contract.ts';
import type { Kanal } from '../protokol/kanal.ts';
import { wywolaj } from '../protokol/wywolanie.ts';

/** Dane znane warstwie wejścia w chwili zalogowania, nieosiągalne z listy komend powłoki. */
export interface DanePowloki {
  /** Login Operatora podany przy wejściu; pusty zostawia pozycję paska bez zmiany. */
  login?: string;
  /** Nazwa okna bieżącego wystawiana w belce tytułowej. */
  nazwaOkna?: string;
  /** Kod środowiska bieżącego; puste zostawia w pasku znak pustej wartości prototypu. */
  kodSrodowiska?: string;
  /** Kod modułu bieżącego wewnątrz środowiska. */
  kodModulu?: string;
}

/** Rozliczenie wiązania: ile pozycji rdzeń pokrył, ile znacznik ma ponad to, czego zabrakło. */
export interface RozliczeniePowloki {
  /** Czy powłoka była zamontowana i wiązanie doszło do skutku. */
  zwiazana: boolean;
  /** Liczba środowisk oddanych przez rdzeń. */
  srodowiskaRdzenia: number;
  /** Liczba pozycji środowisk stojących w znaczniku. */
  srodowiskaZnacznika: number;
  /** Pozycje modułów w znaczniku ukryte, bo rdzeń nie miał dla nich pokrycia. */
  ukrytePozycje: number;
  /** Środowiska i moduły rdzenia bez pozycji w znaczniku; nie są dorabiane. */
  bezPozycji: string[];
  /** Napotkane niepowodzenia komend, opisem błędu z kontraktu. */
  bledy: string[];
}

/** Kanał wystawiony przez `aplikacja.ts`; skrypty biblioteki nie są modułami, więc importu z nich nie ma. */
function kanalGlobalny(): Kanal | undefined {
  return (globalThis as { DanacoKanal?: Kanal }).DanacoKanal;
}

/**
 * Wypełnia powłokę wartościami rdzenia. Woła to warstwa wejścia po zalogowaniu,
 * bo dopiero wtedy znany jest login i okno, do którego Operator wchodzi.
 */
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

/** Kolejność wyświetlania jest własnością rdzenia, nie kolejnością odpowiedzi. */
function uporzadkuj(wykaz: Environment[]): Environment[] {
  return [...wykaz].sort((a, b) => a.order - b.order);
}

/** Nazwa pod pozycją szyny: pogrubiona nazwa i opis; oba węzły stoją już w znaczniku. */
function opiszPozycje(pozycja: Element, nazwa: string, opis: string | undefined): void {
  pozycja.setAttribute('aria-label', nazwa);
  const etykieta = pozycja.querySelector('.dn-szyna-etyk');
  if (etykieta === null) return;
  const tytul = etykieta.querySelector('b');
  if (tytul !== null) tytul.textContent = nazwa;
  const podpis = etykieta.querySelector('span');
  if (podpis !== null) podpis.textContent = opis ?? '';
}

/** Pozycja bez pokrycia w rdzeniu znika z widoku; dorabianie brakującej byłoby atrapą. */
function ukryj(element: Element, rozliczenie: RozliczeniePowloki): void {
  element.setAttribute('hidden', '');
  rozliczenie.ukrytePozycje += 1;
}

/**
 * Wpisuje środowiska i ich moduły w gotowe pozycje szyny. Grupa modułów wiąże
 * się z pozycją przez aria-controls — atrybut zostaje nietknięty, bo po nim
 * rozwija grupy `zasoby/rama.js`.
 *
 * Pozycję z rejestrem wiąże kod w `data-srodowisko`, nie miejsce w znaczniku:
 * kolejność pozycji nie jest kolejnością rejestru, więc wiązanie po numerze
 * dawało pozycji podpis cudzego środowiska.
 */
function wypelnijSzyne(
  powloka: HTMLElement,
  wykazSrodowisk: Environment[],
  katalogModulow: Map<string, Module>,
  rozliczenie: RozliczeniePowloki,
): void {
  const pozycje = [...powloka.querySelectorAll<HTMLElement>('.dn-szyna-poz--srodowisko')];
  rozliczenie.srodowiskaZnacznika = pozycje.length;
  const poKodzie = spisPoKodzie(pozycje, 'srodowisko');
  const zajete = new Set<HTMLElement>();

  for (const srodowisko of wykazSrodowisk) {
    const pozycja = poKodzie.get(srodowisko.code.toLowerCase());
    if (pozycja === undefined) {
      rozliczenie.bezPozycji.push(`środowisko ${srodowisko.code}`);
      continue;
    }
    zajete.add(pozycja);
    const grupa = document.getElementById(pozycja.getAttribute('aria-controls') ?? '');
    pozycja.removeAttribute('hidden');
    pozycja.setAttribute('data-srodowisko', srodowisko.code);
    opiszPozycje(pozycja, srodowisko.name, srodowisko.description);
    przestaw(pozycja, grupa);
    if (grupa !== null) {
      /* Zasłona grupy zostaje nietknięta: rozwija ją i zwija `zasoby/rama.js`
         przy naciśnięciu pozycji środowiska. */
      grupa.setAttribute('aria-label', `Moduły środowiska ${srodowisko.name}`);
      wypelnijModuly(grupa, srodowisko, katalogModulow, rozliczenie);
    }
  }

  for (const pozycja of pozycje) {
    if (zajete.has(pozycja)) continue;
    ukryj(pozycja, rozliczenie);
    const grupa = document.getElementById(pozycja.getAttribute('aria-controls') ?? '');
    if (grupa !== null) grupa.setAttribute('hidden', '');
  }
}

/** Moduły widoczne w środowisku podaje samo środowisko; katalog daje im nazwy. */
function wypelnijModuly(
  grupa: HTMLElement,
  srodowisko: Environment,
  katalogModulow: Map<string, Module>,
  rozliczenie: RozliczeniePowloki,
): void {
  const pozycje = [...grupa.querySelectorAll<HTMLElement>('.dn-szyna-poz--modul')];
  const poKodzie = spisPoKodzie(pozycje, 'modul');
  const zajete = new Set<HTMLElement>();

  for (const kod of srodowisko.moduleCodes ?? []) {
    const modul = katalogModulow.get(kod);
    const pozycja = poKodzie.get(kod.toLowerCase());
    if (modul === undefined || pozycja === undefined) {
      rozliczenie.bezPozycji.push(`${srodowisko.code}/${kod}`);
      continue;
    }
    zajete.add(pozycja);
    pozycja.removeAttribute('hidden');
    pozycja.setAttribute('data-modul', modul.code);
    opiszPozycje(pozycja, modul.name, modul.description);
    grupa.appendChild(pozycja);
  }

  for (const pozycja of pozycje) {
    if (!zajete.has(pozycja)) ukryj(pozycja, rozliczenie);
  }
}

/** Spis pozycji po kodzie z ich cechy; kod sprowadza się do małych liter, bo znacznik pisze nazwy wielką, a rejestr trzyma kody małymi. */
function spisPoKodzie(
  pozycje: HTMLElement[],
  cecha: 'srodowisko' | 'modul',
): Map<string, HTMLElement> {
  const spis = new Map<string, HTMLElement>();
  for (const pozycja of pozycje) {
    const kod = (pozycja.dataset[cecha] ?? '').toLowerCase();
    if (kod !== '' && !spis.has(kod)) spis.set(kod, pozycja);
  }
  return spis;
}

/* Kolejność pozycji szyny jest własnością rejestru: pozycja wraz ze swoją grupą
   modułów idzie na koniec listy, więc po przejściu całego rejestru lista stoi
   w jego porządku. Grupa musi iść za swoją pozycją, bo rama rozwija ją
   sąsiedztwem w znaczniku. */
function przestaw(pozycja: HTMLElement, grupa: HTMLElement | null): void {
  const lista = pozycja.parentElement;
  if (lista === null) return;
  lista.appendChild(pozycja);
  if (grupa !== null && grupa.parentElement === lista) lista.appendChild(grupa);
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

/** Pierwsza pozycja paska niesie rolę w pogrubieniu i login w tekście za nią. */
function wpiszKonto(stan: Element, login: string | undefined): void {
  if (login === undefined || login.length === 0) return;
  const pozycja = stan.querySelector('.dn-stan-poz');
  const rola = pozycja === null ? null : pozycja.querySelector('b');
  const tekst = rola === null ? null : rola.nextSibling;
  if (tekst !== null) tekst.nodeValue = ` · ${login}`;
}

/** Druga pozycja paska mówi, gdzie Operator stoi: środowisko w pogrubieniu, moduł za nim. */
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

/** Liczba sesji czynnych; rdzeń podaje sumę, nie długość odcinka wykazu. */
function wpiszSesje(stan: Element, liczba: number): void {
  const licznik = stan.querySelector('[data-stan-sesje]');
  if (licznik !== null) licznik.textContent = String(liczba);
}
