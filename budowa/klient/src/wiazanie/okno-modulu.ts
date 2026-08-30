/**
 * Wiązanie okna modułu z rdzeniem. Znacznik niesie biblioteka Właściciela;
 * ten plik zakłada oknu sesję w rdzeniu i oddaje je wiązaniu szczegółowemu
 * modułu, a treść przykładową zdejmuje, zanim Operator ją zobaczy.
 */

import {
  Command,
  ExecutionEnv,
  PermissionMode,
  WindowRole,
} from '../../../shared/contract.ts';
import type { Kanal } from '../protokol/kanal.ts';
import { wywolaj } from '../protokol/wywolanie.ts';
import { zwiazBrowser } from './okno-browser.ts';
import { zwiazDeveloper } from './okno-developer.ts';
import { zwiazLibrary } from './okno-library.ts';
import { zwiazWorkspace } from './okno-workspace.ts';

/** Wiązania szczegółowe modułów, po kodzie rejestru rdzenia; moduł bez wpisu dostaje samo okno rdzenia. */
type WiazanieModulu = (kanal: Kanal, idOkna: string) => void;
const WIAZANIA = new Map<string, WiazanieModulu>([
  ['browser', zwiazBrowser],
  ['developer', zwiazDeveloper],
  ['library', zwiazLibrary],
  ['workspace', zwiazWorkspace],
]);

/**
 * Czy moduł ma wiązanie wypełniające jego wnętrze odpowiedzią rdzenia. Wnętrze
 * bez wiązania pokazałoby treść przykładową prototypu jako pracę Operatora,
 * więc takiego okna wydanie nie stawia wcale.
 */
export function maWiazanie(kod: string): boolean {
  return WIAZANIA.has(kod);
}

/** Rejestruje wiązanie szczegółowe modułu. Woła to moduł wiązania przy wczytaniu, więc kolejność plików nie ma znaczenia. */
export function zglosWiazanieModulu(kod: string, wiazanie: WiazanieModulu): void {
  WIAZANIA.set(kod, wiazanie);
}

/**
 * Zakłada oknu modułu sesję i okno komunikacji w rdzeniu, po czym oddaje je
 * wiązaniu szczegółowemu. Nazwa środowiska wchodzi w nagłówek okna.
 */
export function zwiazOkno(kanal: Kanal, kodModulu: string, nazwaSrodowiska: string): void {
  void opiszNaglowek(kanal, nazwaSrodowiska);
  void otworzOkno(kanal, kodModulu).then((idOkna) => {
    if (idOkna === '') return;
    WIAZANIA.get(kodModulu)?.(kanal, idOkna);
  });
}

/** Zakłada sesję i okno komunikacji modułu; pusty wynik znaczy, że rdzeń odmówił i wiązania szczegółowego nie ma po co wołać. */
async function otworzOkno(kanal: Kanal, kodModulu: string): Promise<string> {
  const sesja = await wywolaj(kanal, Command.SessionCreate, {});
  if (!sesja.udany || sesja.wynik === undefined) return '';

  const moduly = await wywolaj(kanal, Command.ModuleList, {});
  if (!moduly.udany || moduly.wynik === undefined) return '';
  const modul = moduly.wynik.modules.find((pozycja) => pozycja.code === kodModulu);
  if (modul === undefined) return '';

  const kanaly = await wywolaj(kanal, Command.ChannelList, { enabledOnly: true });
  const kanalModelu = kanaly.udany ? (kanaly.wynik?.channels[0]?.id ?? '') : '';
  if (kanalModelu === '') return '';

  const okno = await wywolaj(kanal, Command.WindowCreate, {
    sessionId: sesja.wynik.session.id,
    moduleId: modul.id,
    modelChannelId: kanalModelu,
    workingDirs: [],
    executionEnv: ExecutionEnv.Local,
    permissionMode: PermissionMode.Manual,
    windowRole: WindowRole.Standalone,
  });
  return okno.udany && okno.wynik !== undefined ? okno.wynik.window.id : '';
}

/**
 * Opisuje nagłówek okna komunikacji środowiskiem wejścia i modelem kanału.
 * Pozostałe podpisy prototypu — wysiłek, wykonawca, pamięć, rola, format,
 * tura — schodzą: ich wartości niesie dopiero praca podjęta w oknie, a okno
 * dopiero co stanęło, więc podpis pokazywałby wartość wymyśloną.
 */
async function opiszNaglowek(kanal: Kanal, nazwaSrodowiska: string): Promise<void> {
  const naglowek = document.querySelector('.sta-kom-naglowek');
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

/** Wpisuje wartość w pole nagłówka rozpoznane po jego podpisie; wartość pusta zdejmuje całe pole, bo pole bez wartości niczego nie mówi. */
export function wpiszPole(pola: Element[], podpis: string, wartosc: string): void {
  const pole = podpis === ''
    ? pola[0]
    : pola.find((kandydat) => kandydat.textContent?.startsWith(podpis) === true);
  if (pole === undefined) return;
  if (wartosc === '') {
    pole.remove();
    return;
  }
  /* Wartość stoi w prototypie raz w `.dane`, raz w samym wyróżnieniu — pole
     środowiska nie niesie klasy danych, a wpisać się musi tak samo. */
  const dane = pole.querySelector('.dane') ?? pole.querySelector('b');
  if (dane !== null) dane.textContent = wartosc;
}

/** Wpisuje wartość w pierwszy niepusty węzeł tekstowy, zostawiając ikonę i przyciski znacznika; fałsz znaczy węzeł bez miejsca na tekst. */
export function wpiszTekst(wezel: Element, tekst: string): boolean {
  for (const dziecko of wezel.childNodes) {
    if (dziecko.nodeType !== Node.TEXT_NODE) continue;
    if ((dziecko.nodeValue ?? '').trim() === '') continue;
    dziecko.nodeValue = tekst;
    return true;
  }
  return false;
}

/**
 * Wpisuje wartość w węzeł tekstowy znacznika, zostawiając ikonę i przyciski.
 * Wartość pusta zostawia węzeł pusty, nie zdejmuje go: pole, które kontrakt
 * zna, czeka na wartość następną — a wpis pusty zachowuje węzeł tekstowy,
 * więc wartość następna ma dokąd wejść.
 */
export function opiszWezel(wezel: Element | null, wartosc: string): void {
  if (wezel === null) return;
  const teksty = [...wezel.childNodes].filter((dziecko) => dziecko.nodeType === Node.TEXT_NODE);
  const cel = teksty.find((dziecko) => (dziecko.nodeValue ?? '').trim() !== '')
    ?? teksty[teksty.length - 1];
  if (cel === undefined) return;
  cel.nodeValue = wartosc;
}

/** Klon węzła wzorcowego, odporny na jego brak w znaczniku. */
export function sklonuj(wezel: Element | null): HTMLElement | null {
  return wezel instanceof HTMLElement ? (wezel.cloneNode(true) as HTMLElement) : null;
}

/** Godzina znacznika czasu rdzenia w zapisie, którego używa prototyp — godziny i minuty. */
export function godzina(znacznik: number): string {
  return new Date(znacznik).toLocaleTimeString('pl-PL', { hour: '2-digit', minute: '2-digit' });
}

/** Data znacznika czasu rdzenia w zapisie dziennym prototypu. */
export function data(znacznik: number): string {
  return new Date(znacznik).toLocaleDateString('pl-PL');
}

/** Zdejmuje sterowanie wspólne wszystkim oknom modułów: wybór modelu i nakładu, urządzenia dźwięku oraz konektory i wtyczki menu dodawania. */
export function zdejmijSterowanieWspolne(cialo: HTMLElement): void {
  cialo.querySelector('#pop-mik')?.closest('.sta-nrz')?.remove();
  cialo.querySelector('#pop-model')?.closest('.sta-nrz')?.remove();
  cialo.querySelector('#pop-wysilek')?.closest('.sta-nrz')?.remove();
  // Podpis piątego trybu opisuje domyślność, której rejestr uprawnień nie zna.
  cialo.querySelector('#pop-tryb-upr .sta-popover-wiersz small')?.remove();
  const menu = cialo.querySelector('#pop-plus');
  if (menu === null) return;
  // Menu dodawania zostaje przy tytule i trzech pozycjach plików; konektory
  // i wtyczki należą do rodzin spoza obszaru modułu.
  for (const [numer, pozycja] of [...menu.children].entries()) {
    if (numer > 3) pozycja.remove();
  }
}

/** Ustawia znaczniki przełączników paneli tak, by zgadzały się z panelami, które w znaczniku zostały. */
export function uzgodnijPrzelacznikiPaneli(cialo: HTMLElement, bezPokrycia: string[]): void {
  for (const nazwa of bezPokrycia) {
    cialo.querySelector(`#${nazwa}`)?.remove();
    cialo.querySelector(`[data-panel-toggle="${nazwa}"]`)?.remove();
  }
  for (const przelacznik of cialo.querySelectorAll<HTMLElement>('[data-panel-toggle]')) {
    const panel = cialo.querySelector(`#${przelacznik.dataset.panelToggle ?? ''}`);
    przelacznik.setAttribute('aria-checked', String(panel?.hasAttribute('hidden') === false));
  }
}

/** Zdejmuje z okna komunikacji treść przykładową wspólną wszystkim oknom modułów: znacznik pracy, historię rozmowy i znaczniki menu zawężania. */
export function zdejmijTrescWspolna(cialo: HTMLElement): void {
  cialo.querySelector('.sta-kom-naglowek .sta-kom-stan')?.remove();
  cialo.querySelector('.sta-kom-historia')?.replaceChildren();
  cialo.querySelector('.sta-kom-monitor')?.remove();
  for (const znacznik of cialo.querySelectorAll('#menu-filtr .sta-menu-poz[aria-checked]')) {
    znacznik.removeAttribute('aria-checked');
  }
}
