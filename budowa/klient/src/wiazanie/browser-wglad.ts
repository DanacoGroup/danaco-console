// Wgląd techniczny w przeglądaną stronę: odbitka, zrzut ekranu, drzewo
// dokumentu, konsola, ruch sieciowy, kontrola dostępności i emulacja
// urządzenia. Wyniki stają pod płótnem strony, bo okno nie ma gdzie zostawić
// pliku: treść tekstowa idzie do schowka, obraz staje na miejscu.

import {
  BrowserDevicePreset,
  BrowserScreenshotMode,
  Command,
} from '../../../shared/contract.ts';
import type { Kanal } from '../protokol/kanal.ts';
import { wywolaj } from '../protokol/wywolanie.ts';
import {
  doSchowka,
  dolozPrzycisk,
  odmowa,
  powiedz,
  wyborPasa,
  wypiszWynik,
  zalozPas,
} from './browser-wspolne.ts';

const GRANICA_WYKAZU = 100;
const GLEBOKOSC_DRZEWA = 4;

const URZADZENIA: readonly (readonly [string, string])[] = [
  [BrowserDevicePreset.Desktop, 'Biurko'],
  [BrowserDevicePreset.Tablet, 'Tablet'],
  [BrowserDevicePreset.Mobile, 'Telefon'],
  ['reset', 'Bez emulacji'],
];

export function zwiazWgladTechniczny(
  kanal: Kanal,
  korzen: Element,
  idOkna: string,
  przy: AddEventListenerOptions,
): void {
  const wyniki = postawPas(korzen);

  korzen.addEventListener('click', (zdarzenie) => {
    const cel = zdarzenie.target;
    if (!(cel instanceof Element)) return;
    const przycisk = cel.closest<HTMLElement>('[data-wglad]');
    if (przycisk === null) return;
    void wykonaj(kanal, korzen, idOkna, wyniki, przycisk.dataset.wglad ?? '');
  }, przy);
}

function postawPas(korzen: Element): HTMLElement | null {
  const panel = korzen.querySelector('#panel-browser');
  const stopka = panel?.querySelector('.brw-stopka') ?? null;
  if (panel === null || stopka === null) return null;
  const pas = zalozPas('Wgląd techniczny');
  dolozPrzycisk(pas, 'Odbitka strony', 'wglad', 'odbitka');
  dolozPrzycisk(pas, 'Zrzut ekranu', 'wglad', 'zrzut');
  dolozPrzycisk(pas, 'Drzewo dokumentu', 'wglad', 'drzewo');
  dolozPrzycisk(pas, 'Konsola', 'wglad', 'konsola');
  dolozPrzycisk(pas, 'Ruch sieciowy', 'wglad', 'siec');
  dolozPrzycisk(pas, 'Dostępność', 'wglad', 'dostepnosc');
  const wybor = wyborPasa(pas, 'Emulowane urządzenie', URZADZENIA);
  if (wybor !== null) wybor.dataset.wyborUrzadzenia = '';
  dolozPrzycisk(pas, 'Emuluj urządzenie', 'wglad', 'urzadzenie');
  const wyniki = document.createElement('div');
  wyniki.className = 'dn-wykaz-modulu';
  wyniki.dataset.wgladWynik = '';
  panel.insertBefore(pas, stopka);
  panel.insertBefore(wyniki, stopka);
  return wyniki;
}

async function wykonaj(
  kanal: Kanal,
  korzen: Element,
  idOkna: string,
  wyniki: HTMLElement | null,
  czynnosc: string,
): Promise<void> {
  if (czynnosc === 'odbitka') return odbitkaStrony(kanal, idOkna);
  if (czynnosc === 'zrzut') return zrzutEkranu(kanal, idOkna, wyniki);
  if (czynnosc === 'drzewo') return drzewoDokumentu(kanal, idOkna, wyniki);
  if (czynnosc === 'konsola') return odczytKonsoli(kanal, idOkna, wyniki);
  if (czynnosc === 'siec') return ruchSieciowy(kanal, idOkna, wyniki);
  if (czynnosc === 'dostepnosc') return kontrolaDostepnosci(kanal, idOkna, wyniki);
  if (czynnosc === 'urzadzenie') return emulujUrzadzenie(kanal, korzen, idOkna);
}

async function odbitkaStrony(kanal: Kanal, idOkna: string): Promise<void> {
  const wynik = await wywolaj(kanal, Command.BrowserSnapshotGet, {
    windowId: idOkna,
    includeHtml: false,
  });
  if (!wynik.udany) {
    odmowa(wynik.blad, 'Rdzeń odmówił odbitki strony.');
    return;
  }
  const migawka = wynik.wynik?.snapshot;
  await doSchowka(migawka?.text ?? '', migawka?.title ?? migawka?.url ?? 'odbitka');
}

/* Wykonanie zrzutu oddaje samo odwołanie do treści; obraz przychodzi dopiero
   odczytem migawki z żądaniem treści, więc czynność idzie dwoma wywołaniami. */
async function zrzutEkranu(
  kanal: Kanal,
  idOkna: string,
  wyniki: HTMLElement | null,
): Promise<void> {
  const zrobiony = await wywolaj(kanal, Command.BrowserScreenshotCapture, {
    windowId: idOkna,
    mode: BrowserScreenshotMode.Viewport,
  });
  if (!zrobiony.udany) {
    odmowa(zrobiony.blad, 'Rdzeń odmówił wykonania zrzutu ekranu.');
    return;
  }
  const odwolanie = zrobiony.wynik?.screenshot.ref ?? '';
  const odczyt = await wywolaj(kanal, Command.BrowserSnapshotScreenshotGet, {
    screenshotRef: odwolanie,
    includeContent: true,
  });
  if (!odczyt.udany) {
    odmowa(odczyt.blad, 'Zrzut powstał, ale rdzeń odmówił oddania jego treści.');
    return;
  }
  const zrzut = odczyt.wynik?.screenshot;
  if (zrzut === undefined || zrzut.contentBase64 === undefined) {
    wypiszWynik(wyniki, 'Zrzut ekranu', [
      `Zrzut ${odwolanie} stoi w magazynie rdzenia, ale treści obrazu rdzeń nie oddał.`,
    ]);
    return;
  }
  pokazObraz(wyniki, zrzut.format, zrzut.contentBase64, `${zrzut.width}×${zrzut.height}`);
}

function pokazObraz(
  wyniki: HTMLElement | null,
  format: string,
  zapis: string,
  wymiary: string,
): void {
  if (wyniki === null) return;
  wypiszWynik(wyniki, `Zrzut ekranu — ${wymiary}`, []);
  const obraz = document.createElement('img');
  obraz.src = `data:image/${format};base64,${zapis}`;
  obraz.alt = 'Zrzut ekranu przeglądanej strony';
  obraz.style.maxWidth = '100%';
  wyniki.append(obraz);
}

async function drzewoDokumentu(
  kanal: Kanal,
  idOkna: string,
  wyniki: HTMLElement | null,
): Promise<void> {
  const wynik = await wywolaj(kanal, Command.BrowserDomInspect, {
    windowId: idOkna,
    depth: GLEBOKOSC_DRZEWA,
  });
  if (!wynik.udany) {
    odmowa(wynik.blad, 'Rdzeń odmówił wglądu w drzewo dokumentu.');
    return;
  }
  const wezly = wynik.wynik?.nodes ?? [];
  wypiszWynik(wyniki, `Drzewo dokumentu — węzłów: ${wezly.length}`, wezly.map(
    (wezel) => `${'· '.repeat(wezel.depth)}${wezel.tagName} ${wezel.selector ?? ''}`.trimEnd(),
  ));
}

async function odczytKonsoli(
  kanal: Kanal,
  idOkna: string,
  wyniki: HTMLElement | null,
): Promise<void> {
  const wynik = await wywolaj(kanal, Command.BrowserConsoleRead, {
    windowId: idOkna,
    limit: GRANICA_WYKAZU,
  });
  if (!wynik.udany) {
    odmowa(wynik.blad, 'Rdzeń odmówił odczytu konsoli.');
    return;
  }
  const wpisy = wynik.wynik?.entries ?? [];
  wypiszWynik(wyniki, `Konsola — wpisów: ${wpisy.length}`, wpisy.map(
    (wpis) => `${wpis.level.toUpperCase()} · ${wpis.text}`,
  ));
}

async function ruchSieciowy(
  kanal: Kanal,
  idOkna: string,
  wyniki: HTMLElement | null,
): Promise<void> {
  const wynik = await wywolaj(kanal, Command.BrowserNetworkHar, {
    windowId: idOkna,
    limit: GRANICA_WYKAZU,
  });
  if (!wynik.udany) {
    odmowa(wynik.blad, 'Rdzeń odmówił zapisu ruchu sieciowego.');
    return;
  }
  const wpisy = wynik.wynik?.entries ?? [];
  wypiszWynik(wyniki, `Ruch sieciowy — żądań: ${wpisy.length}`, wpisy.map(
    (wpis) => `${wpis.method} ${wpis.statusCode ?? '—'} · ${wpis.url}`,
  ));
}

async function kontrolaDostepnosci(
  kanal: Kanal,
  idOkna: string,
  wyniki: HTMLElement | null,
): Promise<void> {
  const wynik = await wywolaj(kanal, Command.BrowserAccessibilityAudit, {
    windowId: idOkna,
    includeWarnings: true,
  });
  if (!wynik.udany) {
    odmowa(wynik.blad, 'Rdzeń odmówił kontroli dostępności.');
    return;
  }
  const odpowiedz = wynik.wynik;
  const naglowek = `Dostępność ${odpowiedz?.standard ?? ''} — naruszeń: `
    + `${odpowiedz?.errorCount ?? 0}, ostrzeżeń: ${odpowiedz?.warningCount ?? 0}`;
  wypiszWynik(wyniki, naglowek, (odpowiedz?.issues ?? []).map(
    (zgloszenie) => `${zgloszenie.level} · ${zgloszenie.code} · ${zgloszenie.message}`,
  ));
}

async function emulujUrzadzenie(kanal: Kanal, korzen: Element, idOkna: string): Promise<void> {
  const wybor = korzen.querySelector<HTMLSelectElement>('[data-wybor-urzadzenia]');
  const wskazanie = wybor?.value ?? BrowserDevicePreset.Desktop;
  const zadanie = wskazanie === 'reset'
    ? { windowId: idOkna, reset: true }
    : { windowId: idOkna, preset: wskazanie as BrowserDevicePreset };
  const wynik = await wywolaj(kanal, Command.BrowserDeviceEmulate, zadanie);
  if (!wynik.udany) {
    odmowa(wynik.blad, 'Rdzeń odmówił emulacji urządzenia.');
    return;
  }
  const miary = wynik.wynik?.metrics;
  powiedz(`Strona stoi w wymiarach ${miary?.width ?? '?'}×${miary?.height ?? '?'} `
    + `(${miary?.preset ?? wskazanie}).`);
}
