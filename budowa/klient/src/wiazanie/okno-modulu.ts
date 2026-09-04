// Wiązanie karty okna roboczego z oknem komunikacji rdzenia wraz z wpisami
// wspólnymi wiązaniom modułów. Znacznik niesie biblioteka Właściciela.
import {
  Command,
  ExecutionEnv,
  PermissionMode,
  WindowRole,
  WindowStatus,
} from '../../../shared/contract.ts';
import type { Kanal } from '../protokol/kanal.ts';
import { wywolaj } from '../protokol/wywolanie.ts';
import { oglos } from './ogloszenie.ts';
import { kartyOkna, oknaRobocze, przypiszOknoKomunikacji } from './okna-robocze.ts';
import { zapewnijSesje } from './sesja-biezaca.ts';

export function zwiazOkno(
  kanal: Kanal,
  _kodModulu: string,
  nazwaSrodowiska: string,
  korzen: ParentNode,
): void {
  void opiszNaglowek(kanal, nazwaSrodowiska, korzen);
}

export function zwiazOknoStojace(
  kanal: Kanal,
  idKarty: string,
  nazwaSrodowiska: string,
  idOkna: string,
  korzen: ParentNode,
): void {
  if (!przypiszOknoKomunikacji(idKarty, idOkna)) {
    oglos('Okno modułu', 'Karta zeszła z pasma, zanim okno rdzenia do niej doszło — '
      + 'zamknięcie karty nie zamknie tego okna.', 'ostrzezenie');
  }
  void opiszNaglowek(kanal, nazwaSrodowiska, korzen);
}

/* Nagłówek prototypu niesie środowisko, model i wysiłek wpisane wprost. Pole
   o znanej wartości dostaje ją z rdzenia, pole bez pokrycia schodzi — inaczej
   okno w sesji TalkIn twierdzi, że stoi w WorkSpace. */
export async function opiszNaglowek(
  kanal: Kanal,
  nazwaSrodowiska: string,
  korzen: ParentNode,
): Promise<void> {
  const naglowek = korzen.querySelector('.sta-kom-naglowek');
  if (naglowek === null) return;
  const wynik = await wywolaj(kanal, Command.ChannelList, { enabledOnly: true });
  const kanalModelu = wynik.udany ? wynik.wynik?.channels[0] : undefined;
  const wartosci: Record<string, string> = {
    'Środowisko': nazwaSrodowiska,
    Model: kanalModelu?.model ?? kanalModelu?.name ?? '',
  };
  for (const pole of [...naglowek.querySelectorAll('.sta-kom-pole')]) {
    const podpis = Object.keys(wartosci).find(
      (kandydat) => pole.textContent?.startsWith(kandydat) === true,
    );
    wpiszPole([pole], podpis ?? '', podpis === undefined ? '' : wartosci[podpis]);
  }
}

export function wpiszPole(pola: Element[], podpis: string, wartosc: string): void {
  const pole = podpis === ''
    ? pola[0]
    : pola.find((kandydat) => kandydat.textContent?.startsWith(podpis) === true);
  if (pole === undefined) return;
  if (wartosc === '') {
    pole.remove();
    return;
  }
  // Pole środowiska w prototypie nie niesie klasy `.dane`, tylko wyróżnienie.
  const dane = pole.querySelector('.dane') ?? pole.querySelector('b');
  if (dane !== null) dane.textContent = wartosc;
}

export function wpiszTekst(wezel: Element, tekst: string): boolean {
  for (const dziecko of wezel.childNodes) {
    if (dziecko.nodeType !== Node.TEXT_NODE) continue;
    if ((dziecko.nodeValue ?? '').trim() === '') continue;
    dziecko.nodeValue = tekst;
    return true;
  }
  return false;
}

export function opiszWezel(wezel: Element | null, wartosc: string): void {
  if (wezel === null) return;
  const teksty = [...wezel.childNodes].filter((dziecko) => dziecko.nodeType === Node.TEXT_NODE);
  const cel = teksty.find((dziecko) => (dziecko.nodeValue ?? '').trim() !== '')
    ?? teksty[teksty.length - 1];
  if (cel === undefined) return;
  cel.nodeValue = wartosc;
}

export function sklonuj(wezel: Element | null): HTMLElement | null {
  return wezel instanceof HTMLElement ? (wezel.cloneNode(true) as HTMLElement) : null;
}

export function godzina(znacznik: number): string {
  return new Date(znacznik).toLocaleTimeString('pl-PL', { hour: '2-digit', minute: '2-digit' });
}

export function data(znacznik: number): string {
  return new Date(znacznik).toLocaleDateString('pl-PL');
}

export function zdejmijSterowanieWspolne(cialo: HTMLElement): void {
  // Rejestr uprawnień nie zna domyślności, którą opisuje podpis piątego trybu.
  cialo.querySelector('#pop-tryb-upr .sta-popover-wiersz small')?.remove();
}

// Panel bez pokrycia w rdzeniu zostaje: układ okna jest kompozycją Właściciela.
export function uzgodnijPrzelacznikiPaneli(cialo: HTMLElement, _bezPokrycia: string[]): void {
  for (const przelacznik of cialo.querySelectorAll<HTMLElement>('[data-panel-toggle]')) {
    const panel = cialo.querySelector(`#${przelacznik.dataset.panelToggle ?? ''}`);
    /* Przycisk bez roli nie przyjmuje `aria-checked`; stan wciśnięcia niesie
       `aria-pressed`, jedyny atrybut przełącznika dopuszczony na przycisku. */
    przelacznik.setAttribute('aria-pressed', String(panel?.hasAttribute('hidden') === false));
    przelacznik.removeAttribute('aria-checked');
  }
}

export function zdejmijTrescWspolna(cialo: HTMLElement): void {
  const historia = cialo.querySelector('.sta-kom-historia');
  historia?.replaceChildren();
  /* Wiadomość wysyła wyłącznie Studio, a Studio tej ścieżki nie wywołuje.
     Puste pole nadawania obok wybranego modelu obiecuje odpowiedź, której
     nikt nie da, więc miejsce rozmowy mówi wprost, jak jest. */
  if (historia !== null) {
    const zdanie = cialo.ownerDocument.createElement('p');
    zdanie.className = 'dn-tekst-ciagly dn-tekst-ciagly--drobny';
    zdanie.dataset.pustka = 'tak';
    zdanie.textContent = 'Rozmowa tego modułu nie jest jeszcze związana z rdzeniem.';
    historia.appendChild(zdanie);
  }
  /* Monitor pokazuje, czym model zajmuje się w tej chwili. Bez wiązanej rozmowy
     nie ma czego pokazywać, a jego tytuł prototypu opisuje cudzą pracę. */
  cialo.querySelector('.sta-kom-monitor')?.remove();
  /* Żetony kontekstu i znacznik pracy prototyp wpisuje wprost — nazwy plików,
     adresy stron, „pracuje”. Bez zdjęcia stoją obok prawdziwych zer wykazu
     jako dwa stany wykluczające się, oba wzięte znikąd. */
  cialo.querySelector('.sta-kom-kontekst')?.replaceChildren();
  /* Belka panelu rozmowy niesie nazwę pracy wymyśloną na pokaz — raz żetonem,
     raz dopiskiem w cudzysłowie przy tytule panelu. */
  for (const zeton of cialo.querySelectorAll('.sta-kom > .sta-okno-belka .sta-chip')) {
    zeton.remove();
  }
  const tytulRozmowy = cialo.querySelector('.sta-kom > .sta-okno-belka .sta-okno-tytul b');
  if (tytulRozmowy !== null) {
    tytulRozmowy.textContent = (tytulRozmowy.textContent ?? '')
      .replace(/\s*[—–-]?\s*„[^”]*”/u, '').trim();
  }
  cialo.querySelector('.sta-kontekst-akcji')?.replaceChildren();
  cialo.querySelector('.sta-kom-stan')?.remove();
  /* Znaczniki paneli niosą w prototypie liczby wzięte znikąd. Wiązanie panelu,
     które ma czym je wypełnić, robi to po tym zdjęciu. */
  for (const znacznik of cialo.querySelectorAll('.sta-okno-znacznik')) {
    znacznik.textContent = '';
  }
  for (const znacznik of cialo.querySelectorAll('#menu-filtr .sta-menu-poz[aria-checked]')) {
    znacznik.removeAttribute('aria-checked');
  }
}

const MIEJSCA = new WeakMap<Element, { rodzic: Element; przed: Node | null }>();

export function zdejmijAlboPostaw(wezel: Element | null, obecny: boolean): void {
  if (wezel === null) return;
  if (!obecny) {
    if (wezel.parentElement === null) return;
    MIEJSCA.set(wezel, { rodzic: wezel.parentElement, przed: wezel.nextSibling });
    wezel.remove();
    return;
  }
  if (wezel.parentElement !== null) return;
  const miejsce = MIEJSCA.get(wezel);
  miejsce?.rodzic.insertBefore(wezel, miejsce.przed);
}

export function wpiszAlboZdejmij(wezel: Element | null, wartosc: string): void {
  if (wezel === null) return;
  if (wartosc !== '') wezel.textContent = wartosc;
  zdejmijAlboPostaw(wezel, wartosc !== '');
}

export function wpiszTekstAlboZdejmij(wezel: Element | null, wartosc: string): void {
  if (wezel === null) return;
  if (wartosc !== '' && !wpiszTekst(wezel, wartosc)) {
    wezel.remove();
    return;
  }
  zdejmijAlboPostaw(wezel, wartosc !== '');
}

/* Moduł bez okna nie ma czym wołać komend rdzenia: wykaz kart, źródeł czy
   ustaleń odmawia bez wskazania okna. Droga jest jedna dla wszystkich modułów. */
export async function zapewnijOknoModulu(
  kanal: Kanal,
  idKarty: string,
  kodModulu: string,
  nazwaSesji: string,
  stojace: string,
): Promise<string> {
  if (stojace !== '') return stojace;
  const idSesji = await zapewnijSesje(kanal, idKarty, nazwaSesji);
  if (idSesji === '') return '';
  const moduly = await wywolaj(kanal, Command.ModuleList, {});
  const idModulu = moduly.wynik?.modules.find((m) => m.code === kodModulu)?.id ?? '';
  if (idModulu === '') return '';
  const wolne = await wskazOknoWolneModulu(kanal, idSesji, idModulu);
  if (wolne !== '') {
    przypiszOknoKomunikacji(idKarty, wolne);
    return wolne;
  }
  const kanaly = await wywolaj(kanal, Command.ChannelList, { enabledOnly: true });
  const idKanalu = kanaly.wynik?.channels[0]?.id ?? '';
  if (idKanalu === '') return '';
  const okno = await wywolaj(kanal, Command.WindowCreate, {
    sessionId: idSesji,
    moduleId: idModulu,
    modelChannelId: idKanalu,
    workingDirs: [],
    executionEnv: ExecutionEnv.Local,
    permissionMode: PermissionMode.Manual,
    windowRole: WindowRole.Standalone,
  });
  const powstale = okno.wynik?.window.id ?? '';
  if (powstale !== '') przypiszOknoKomunikacji(idKarty, powstale);
  return powstale;
}

async function wskazOknoWolneModulu(kanal: Kanal, idSesji: string, idModulu: string): Promise<string> {
  const odpowiedz = await wywolaj(kanal, Command.WindowList, {
    sessionId: idSesji,
    status: WindowStatus.Open,
  });
  if (!odpowiedz.udany || odpowiedz.wynik === undefined) return '';
  const zajete = new Set<string>();
  for (const okno of oknaRobocze()) {
    for (const karta of kartyOkna(okno)) {
      if (karta.idOknaKomunikacji !== '') zajete.add(karta.idOknaKomunikacji);
    }
  }
  return odpowiedz.wynik.windows.find(
    (okno) => okno.moduleId === idModulu && !zajete.has(okno.id),
  )?.id ?? '';
}

/* Panel modułu wypełnia się zawsze tak samo: odmowa rdzenia idzie zdaniem,
   wykaz pusty własnym zdaniem, a pozycja niesie tytuł i podpis. Pomocnik
   trzyma ten porządek w jednym miejscu dla wszystkich modułów. */
export function wykazPanelu<T>(
  cialo: Element | null,
  udany: boolean,
  blad: string | undefined,
  pozycje: T[] | undefined,
  pusty: string,
  mapuj: (pozycja: T) => readonly [string, string],
): void {
  if (cialo === null) return;
  if (!udany || pozycje === undefined) {
    niegotowyPanel(cialo, blad ?? 'Wykaz nie doszedł.');
    return;
  }
  if (pozycje.length === 0) {
    niegotowyPanel(cialo, pusty);
    return;
  }
  cialo.replaceChildren(...pozycje.map((p) => {
    const [tytul, podpis] = mapuj(p);
    const wiersz = cialo.ownerDocument.createElement('div');
    wiersz.className = 'dn-wykaz-modulu-poz';
    const nazwa = cialo.ownerDocument.createElement('span');
    nazwa.textContent = tytul;
    wiersz.append(nazwa);
    if (podpis !== '') {
      const meta = cialo.ownerDocument.createElement('span');
      meta.className = 'dn-meta';
      meta.textContent = podpis;
      wiersz.append(meta);
    }
    return wiersz;
  }));
}

export function niegotowyPanel(cialo: Element | null, zdanie: string): void {
  if (cialo === null) return;
  const napis = cialo.ownerDocument.createElement('div');
  napis.className = 'dn-meta';
  napis.textContent = zdanie;
  cialo.replaceChildren(napis);
}

export function cialoPanelu(korzen: ParentNode, identyfikator: string): Element | null {
  return korzen.querySelector(`#${identyfikator} .sta-okno-tresc`);
}
