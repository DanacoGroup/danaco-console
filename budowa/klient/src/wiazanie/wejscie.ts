// Wiązanie drogi wejścia z komendami rdzenia. Znacznik i przełączanie widoków
// należą do biblioteki `design/zasoby/okna/wejscie/`.

import {
  AuthMethodKind,
  Command,
  type AuthSession,
  type ErrorInfo,
} from '../../../shared/contract.ts';
import type { Kanal } from '../protokol/kanal.ts';
import { zadajPowitanie } from '../protokol/powitanie.ts';
import { tozsamoscKlienta } from '../protokol/tozsamosc-klienta.ts';
import { wywolaj } from '../protokol/wywolanie.ts';
import { zwiazPowloke } from './powloka.ts';
import { zwiazStudio } from './studio.ts';

interface EkranStartowy {
  gotowe(): void;
}

interface Katalog {
  tekst(sciezka: string): unknown;
}

// Skrypty biblioteki są funkcjami domkniętymi, nie modułami — importu z nich nie ma.
interface Styk {
  DanacoKanal?: Kanal;
  DanacoNarzedzia?: { zwiaz(katalog: unknown): Katalog };
  DanacoWejscie?: { tresci?: unknown };
  dnPrzelaczWidok?: (widok: string, grupa?: string | null) => void;
  dnPrzejdz?: (widok: string, grupa?: string | null) => void;
  dnToast?: (tytul: string, tresc: string, rodzaj?: string) => void;
}

function styk(): Styk {
  return globalThis as unknown as Styk;
}

const ETAPY: Readonly<Record<string, string>> = {
  uruchomienie: 'laczenie',
  dostep: 'uwierzytelnienie',
  przygotowanie: 'przygotowanie',
};

/** Widoki, których czynność główna woła rdzeń. */
type Widok =
  | 'logowanie'
  | 'logowanie-blad'
  | 'rejestracja'
  | 'kod'
  | 'odzyskiwanie-adres'
  | 'odzyskiwanie-kod'
  | 'odzyskiwanie-haslo';

// Kod z widoku `odzyskiwanie-kod` zużywa dopiero `auth.reset`, razem z hasłem.
let drogaOdzyskania = '';

let tokenSesji = '';
let ostatniLogin = '';

export function zwiazWejscie(podany?: Kanal): void {
  oznaczObszary();
  zalozPrzejscie();
  document.addEventListener('click', naKlikniecie);
  gotowosc(() => {
    void powitaj(kanal(podany));
  });
}

function kanal(podany?: Kanal): Kanal | undefined {
  return podany ?? styk().DanacoKanal;
}

function gotowosc(bieg: () => void): void {
  if (document.readyState === 'loading') {
    document.addEventListener('DOMContentLoaded', bieg, { once: true });
    return;
  }
  bieg();
}

// Dokument klienta stawia obszary bez znaczników widoku; biblioteka przełącza po nich.
function oznaczObszary(): void {
  for (const miejsce of document.querySelectorAll<HTMLElement>('[data-wejscie-okno]')) {
    const etap = ETAPY[miejsce.dataset.wejscieOkno ?? ''];
    const obszar = miejsce.closest<HTMLElement>('.we-scena');
    if (etap === undefined || obszar === null) continue;
    obszar.dataset.widok = etap;
    obszar.dataset.grupaWidoku = 'etap';
    obszar.dataset.widokAktywny = etap === 'laczenie' ? 'tak' : 'nie';
    if (etap === 'laczenie') obszar.dataset.uruchomienieScena = '';
  }
}

// `dnPrzejdz` wystawia strona podglądu biblioteki, której produkt nie wczytuje.
function zalozPrzejscie(): void {
  const s = styk();
  if (s.dnPrzejdz !== undefined) return;
  s.dnPrzejdz = (widok, grupa) => {
    s.dnPrzelaczWidok?.(widok, grupa ?? null);
  };
}

// Sprawdzenie pól z `przeplyw-wejscia.js` idzie w fazie przechwytywania i rozstrzyga pierwsze.
function naKlikniecie(zdarzenie: MouseEvent): void {
  const cel = zdarzenie.target;
  if (!(cel instanceof Element)) return;
  const czynnosc = cel.closest<HTMLElement>('[data-idz]');
  if (czynnosc === null) return;
  if (czynnosc.tagName === 'A') zdarzenie.preventDefault();

  const dokad = czynnosc.dataset.idz ?? '';
  const grupa = czynnosc.dataset.idzGrupa ?? null;
  const pas = czynnosc.closest<HTMLElement>('.we-pas');
  const widok = (pas?.dataset.widok ?? '') as Widok;
  const most = kanal();

  if (pas === null || most === undefined || !wolaRdzen(widok)) {
    idz(dokad, grupa);
    return;
  }
  zdarzenie.preventDefault();
  void wykonaj(most, widok, dokad, grupa);
}

function wolaRdzen(widok: string): widok is Widok {
  return (
    widok === 'logowanie' ||
    widok === 'logowanie-blad' ||
    widok === 'rejestracja' ||
    widok === 'kod' ||
    widok === 'odzyskiwanie-adres' ||
    widok === 'odzyskiwanie-kod' ||
    widok === 'odzyskiwanie-haslo'
  );
}

function idz(dokad: string, grupa: string | null): void {
  const s = styk();
  if (s.dnPrzejdz !== undefined) s.dnPrzejdz(dokad, grupa);
  else s.dnPrzelaczWidok?.(dokad, grupa);
}

async function powitaj(most: Kanal | undefined): Promise<void> {
  if (most === undefined) {
    ustawWariant('w-blad');
    domknijEkranStartowy();
    return;
  }
  const wynik = await zadajPowitanie(most, tozsamoscKlienta(), tokenSesji || undefined);
  ustawWariant(wynik.udany ? 'w-laczenie' : 'w-blad');
  domknijEkranStartowy();
  /* Dalej okno przechodzi samo: po animacji biblioteka odgrywa etapy łączenia
     i przechodzi do etapu z `data-po-polaczeniu`. Przełączenie stąd wyprzedzałoby
     ten przebieg i pomijało okno uruchomienia. */
}

function ustawWariant(widok: string): void {
  styk().dnPrzelaczWidok?.(widok, 'wariant');
}

function domknijEkranStartowy(): void {
  const pole = document.querySelector('[data-ekran-startowy]');
  const ekran = (pole as unknown as { ekranStartowy?: EkranStartowy } | null)?.ekranStartowy;
  ekran?.gotowe();
}

async function wykonaj(
  most: Kanal,
  widok: Widok,
  dokad: string,
  grupa: string | null,
): Promise<void> {
  const panel = document.querySelector<HTMLElement>(`.we-panel[data-widok="${widok}"]`);
  if (panel === null) return;

  switch (widok) {
    case 'logowanie':
    case 'logowanie-blad':
      await zaloguj(most, panel, widok);
      return;
    case 'rejestracja':
      await zarejestruj(most, panel, dokad, grupa);
      return;
    case 'kod':
      await potwierdz(most, panel);
      return;
    case 'odzyskiwanie-adres':
      await odzyskaj(most, panel, dokad, grupa);
      return;
    case 'odzyskiwanie-kod':
      drogaOdzyskania = zbierzKod(panel);
      idz(dokad, grupa);
      return;
    case 'odzyskiwanie-haslo':
      await ustawHaslo(most, panel, dokad, grupa);
      return;
  }
}

async function zaloguj(most: Kanal, panel: HTMLElement, widok: Widok): Promise<void> {
  const przedrostek = widok === 'logowanie' ? 'log' : 'blad';
  ostatniLogin = wartosc(`${przedrostek}-login`);
  const wynik = await wywolaj(most, Command.AuthLogin, {
    method: AuthMethodKind.Password,
    login: ostatniLogin,
    secret: wartosc(`${przedrostek}-haslo`),
    keepSignedIn: trwalaSesja(panel),
  });
  if (!wynik.udany) {
    odmowaLogowania(wynik.blad);
    return;
  }
  await wejdz(most, wynik.wynik?.session);
}

async function zarejestruj(
  most: Kanal,
  panel: HTMLElement,
  dokad: string,
  grupa: string | null,
): Promise<void> {
  const wynik = await wywolaj(most, Command.AuthRegister, {
    login: wartosc('rej-login'),
    email: wartosc('rej-email'),
    password: wartosc('rej-haslo'),
  });
  if (!wynik.udany) {
    odmowa(panel, 'usterki.naglowekKonto', wynik.blad);
    return;
  }
  idz(dokad, grupa);
}

async function potwierdz(most: Kanal, panel: HTMLElement): Promise<void> {
  const wynik = await wywolaj(most, Command.AuthVerify, {
    token: zbierzKod(panel),
    keepSignedIn: trwalaSesja(panel),
  });
  if (!wynik.udany) {
    odmowa(panel, 'usterki.naglowekKod', wynik.blad);
    return;
  }
  await wejdz(most, wynik.wynik?.session);
}

async function odzyskaj(
  most: Kanal,
  panel: HTMLElement,
  dokad: string,
  grupa: string | null,
): Promise<void> {
  const wynik = await wywolaj(most, Command.AuthRecover, { email: wartosc('odz-email') });
  if (!wynik.udany) {
    odmowa(panel, 'usterki.naglowekKod', wynik.blad);
    return;
  }
  idz(dokad, grupa);
}

async function ustawHaslo(
  most: Kanal,
  panel: HTMLElement,
  dokad: string,
  grupa: string | null,
): Promise<void> {
  const wynik = await wywolaj(most, Command.AuthReset, {
    token: drogaOdzyskania,
    newPassword: wartosc('odz-haslo'),
  });
  if (!wynik.udany) {
    odmowa(panel, 'usterki.naglowekHaslo', wynik.blad);
    return;
  }
  // `auth.reset` unieważnia tokeny wydane wcześniej.
  drogaOdzyskania = '';
  tokenSesji = '';
  idz(dokad, grupa);
}

async function wejdz(most: Kanal, sesja: AuthSession | undefined): Promise<void> {
  if (sesja !== undefined) tokenSesji = sesja.token;
  await Promise.all([
    wywolaj(most, Command.EnvironmentList, {}),
    wywolaj(most, Command.ModuleList, {}),
    wywolaj(most, Command.SessionList, {}),
  ]);
  odslonPowloke();
  await zwiazPowloke({ login: ostatniLogin, nazwaOkna: 'Danaco Console' }, most);
  zwiazStudio(most);
}

function odslonPowloke(): void {
  document.querySelector('[data-rama-aplikacji]')?.removeAttribute('hidden');
  document.querySelector('[data-wejscie]')?.setAttribute('hidden', '');
}

// Adres niepotwierdzony to konto zatrzymane przed aktywacją, nie zła para login–hasło.
function odmowaLogowania(blad: ErrorInfo | undefined): void {
  if (powod(blad) === 'adres-niepotwierdzony') {
    styk().dnPrzelaczWidok?.('kod', 'stan');
    return;
  }
  styk().dnPrzelaczWidok?.('logowanie-blad', 'stan');
}

// Widok bez banera dostaje powiadomienie; wiązanie banera nie dostawia.
function odmowa(panel: HTMLElement, naglowek: string, blad: ErrorInfo | undefined): void {
  const klucz = powod(blad) === 'kolizja-danych' ? 'usterki.loginZajety' : '';
  const glowa = klucz === '' ? napis(naglowek) : napis(`${klucz}.glowa`);
  const tresc = klucz === '' ? napis('usterki.wiele') : napis(`${klucz}.tresc`);
  if (!wpiszWBaner(panel, glowa, tresc)) styk().dnToast?.(glowa, tresc, 'blad');
}

function powod(blad: ErrorInfo | undefined): string {
  const dane = blad?.details as { powod?: unknown } | undefined;
  return typeof dane?.powod === 'string' ? dane.powod : '';
}

function wpiszWBaner(panel: HTMLElement, glowa: string, tresc: string): boolean {
  const baner = panel.querySelector<HTMLElement>('.we-komunikaty .dn-alert');
  const pole = baner?.querySelector<HTMLElement>('.dn-alert-tresc');
  const czolo = pole?.querySelector('b');
  if (baner == null || pole == null || czolo == null) return false;
  baner.classList.remove('dn-alert--info', 'dn-alert--ostrzezenie', 'dn-alert--sukces');
  baner.classList.add('dn-alert--blad');
  baner.setAttribute('role', 'alert');
  baner.setAttribute('data-usterka-formularza', '');
  czolo.textContent = glowa;
  const ostatni = pole.lastChild;
  if (ostatni !== null && ostatni.nodeType === Node.TEXT_NODE) ostatni.textContent = tresc;
  else pole.appendChild(document.createTextNode(tresc));
  return true;
}

function wartosc(id: string): string {
  const pole = document.getElementById(id);
  return pole instanceof HTMLInputElement ? pole.value.trim() : '';
}

function trwalaSesja(panel: HTMLElement): boolean {
  const pole = panel.querySelector<HTMLInputElement>('.au-opcja input[type="checkbox"]');
  return pole?.checked ?? false;
}

function zbierzKod(panel: HTMLElement): string {
  const pola = panel.querySelectorAll<HTMLInputElement>('.au-kod-pole');
  let kod = '';
  for (const pole of pola) kod += pole.value.trim();
  return kod;
}

function napis(sciezka: string): string {
  const s = styk();
  const katalog = s.DanacoNarzedzia?.zwiaz(s.DanacoWejscie?.tresci ?? {});
  const wartoscNapisu = katalog?.tekst(sciezka);
  return typeof wartoscNapisu === 'string' ? wartoscNapisu : '';
}
