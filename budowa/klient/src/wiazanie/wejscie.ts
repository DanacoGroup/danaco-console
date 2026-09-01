// Wiązanie drogi wejścia z komendami rdzenia. Znacznik i przełączanie widoków
// należą do biblioteki `design/zasoby/okna/wejscie/`.

import { invoke, isTauri } from '@tauri-apps/api/core';
import {
  AuthChangeReason,
  AuthMethodKind,
  Command,
  EventType,
  type ErrorInfo,
} from '../../../shared/contract.ts';
import type { Kanal } from '../protokol/kanal.ts';
import { zadajPowitanie } from '../protokol/powitanie.ts';
import { tokenSesji, urzadzenieSesji, zapomnijTokenSesji } from '../protokol/token-sesji.ts';
import { tozsamoscKlienta } from '../protokol/tozsamosc-klienta.ts';
import { wywolaj } from '../protokol/wywolanie.ts';
import { oglos } from './ogloszenie.ts';
import { zwiazPowloke } from './powloka.ts';
import { zwiazCentrum } from './centrum.ts';
import { zglosUchwyt } from './zdarzenia.ts';

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

let ostatniLogin = '';
/* Pierwsze powitanie należy do `powitaj`; dopiero po nim zmiany stanu
   transportu mają w oknie co odzwierciedlać. */
let powitanoRaz = false;

/** Polecenia powłoki z `desktop/src-tauri/src/polecenia.rs`. */
const POLECENIE_WSKAZANIA_RDZENIA = 'wskazanie_rdzenia';
const POLECENIE_WSKAZ_RDZEN = 'wskaz_rdzen';

/** Stan wskazania rdzenia oddawany przez powłokę; nazwy pól jak w `wskazanie.rs`, bo serde ich nie przemianowuje. */
interface WskazanieRdzenia {
  schemat: string;
  host: string | null;
  port: number;
  adres: string | null;
  warstwa: string;
}

/** Odmowa wskazania rdzenia; nazwy pól jak w `wskazanie.rs`. */
interface OdmowaWskazania {
  powod: string;
  zdanie: string;
}

export function zwiazWejscie(podany?: Kanal): void {
  oznaczObszary();
  zalozPrzejscie();
  document.addEventListener('click', naKlikniecie);
  document.addEventListener('click', naPonowienie);
  document.addEventListener('keydown', naKlawisz);
  zwiazStanPolaczenia(kanal(podany));
  zwiazZmianeUwierzytelnienia();
  gotowosc(() => {
    void powitaj(kanal(podany));
  });
}

/**
 * Zerwanie po pierwszym powitaniu pokazuje w oknie wejścia wariant błędu,
 * a powrót łączności — po ponowionym powitaniu — odgrywa łączenie od nowa.
 * Okno już schowane za powłoką nie ma czego pokazywać.
 */
function zwiazStanPolaczenia(most: Kanal | undefined): void {
  if (most === undefined) return;
  most.naStan((stan) => {
    if (!powitanoRaz || !oknoWejsciaStoi()) return;
    if (stan === 'polaczony') {
      void powitajPonownie(most);
      return;
    }
    if (stan === 'rozlaczony' || stan === 'ponawianie') ustawWariant('w-blad');
  });
}

async function powitajPonownie(most: Kanal): Promise<void> {
  const wynik = await zadajPowitanie(most, tozsamoscKlienta(), tokenSesji() || undefined);
  if (!oknoWejsciaStoi()) return;
  if (!wynik.udany) {
    ustawWariant('w-blad');
    return;
  }
  if (etapLaczeniaAktywny()) odegrajLaczenie();
}

function oknoWejsciaStoi(): boolean {
  const okno = document.querySelector('[data-wejscie]');
  return okno !== null && !okno.hasAttribute('hidden');
}

function etapLaczeniaAktywny(): boolean {
  return document.querySelector('.we-scena[data-widok="laczenie"][data-widok-aktywny="tak"]') !== null;
}

/* Przebieg łączenia rusza w bibliotece na `ekran-startowy-koniec`; po powrocie
   łączności to samo zdarzenie odgrywa etapy i przechodzi do etapu następnego. */
function odegrajLaczenie(): void {
  ustawWariant('w-laczenie');
  document
    .querySelector('[data-ekran-startowy]')
    ?.dispatchEvent(new CustomEvent('ekran-startowy-koniec', { bubbles: true }));
}

/**
 * Token schodzi przy odzyskaniu konta i unieważnieniu sesji. Rdzeń adresuje
 * `auth.changed` do konta, więc zdarzenie dotyczy konta bieżącego; unieważnienie
 * nazywające urządzenie zdejmuje token tylko wtedy, gdy to urządzenie tej sesji.
 */
function zwiazZmianeUwierzytelnienia(): void {
  zglosUchwyt(EventType.AuthChanged, (tresc) => {
    const powodZdjecia =
      tresc.reason === AuthChangeReason.PasswordReset ||
      tresc.reason === AuthChangeReason.SessionRevoked;
    if (!powodZdjecia || tokenSesji() === '') return;
    if (tresc.deviceId !== undefined && tresc.deviceId !== urzadzenieSesji()) return;
    zapomnijTokenSesji();
    oglos('Sesja', 'Sesja została unieważniona — po następnym połączeniu trzeba zalogować się ponownie.', 'ostrzezenie');
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

/* Enter w polu wchodzi tak samo jak naciśnięcie czynności głównej. Pas działań
   domyka okno u dołu, więc w oknie niższym niż panel czynność bywa poza
   widokiem — wpisane dane zostawałyby wtedy bez drogi zatwierdzenia. */
function naKlawisz(zdarzenie: KeyboardEvent): void {
  if (zdarzenie.key !== 'Enter' || zdarzenie.shiftKey) return;
  const cel = zdarzenie.target;
  if (!(cel instanceof HTMLInputElement)) return;
  const czynnosc = document.querySelector<HTMLElement>(
    '.we-pas[data-widok-aktywny="tak"] button[data-idz]',
  );
  if (czynnosc === null) return;
  zdarzenie.preventDefault();
  czynnosc.click();
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
  const wynik = await zadajPowitanie(most, tozsamoscKlienta(), tokenSesji() || undefined);
  powitanoRaz = true;
  ustawWariant(wynik.udany ? 'w-laczenie' : 'w-blad');
  domknijEkranStartowy();
  /* Dalej okno przechodzi samo: po animacji biblioteka odgrywa etapy łączenia
     i przechodzi do etapu z `data-po-polaczeniu`. Przełączenie stąd wyprzedzałoby
     ten przebieg i pomijało okno uruchomienia. */
}

function ustawWariant(widok: string): void {
  styk().dnPrzelaczWidok?.(widok, 'wariant');
  if (widok === 'w-blad') void opiszWskazanieRdzenia();
}

/* Wariant błędu nazywa serwer, którego powłoka nie doprosiła: bez tego
   Operator nie wie, czy zawiódł adres, czy sieć. Poza powłoką wskazania nie ma. */
async function opiszWskazanieRdzenia(): Promise<void> {
  if (!isTauri()) return;
  const lid = document.querySelector<HTMLElement>('.we-panel[data-widok="w-blad"] .we-lid');
  if (lid === null) return;
  lid.dataset.lidPierwotny ??= lid.textContent ?? '';
  try {
    const wskazanie = await invoke<WskazanieRdzenia>(POLECENIE_WSKAZANIA_RDZENIA);
    const zdanie = wskazanie.adres === null
      ? 'Powłoka nie ma wskazania rdzenia.'
      : `Rdzeń wskazany: ${wskazanie.adres}.`;
    lid.textContent = `${lid.dataset.lidPierwotny} ${zdanie}`;
  } catch (blad) {
    console.warn('[wejście] powłoka nie oddała wskazania rdzenia', blad);
  }
}

/*
Ponowienie z wariantu błędu. Wskazanie wpisane w pole adresu idzie do powłoki,
która sprawdza łączność i zapisuje je trwale; wskazanie przyjęte wchodzi w stronę
dopiero przy wczytaniu, bo powłoka podaje je skryptem wstępnym okna. Bez pola
albo bez wpisu ponowienie łączy pod adres obowiązujący od razu, bez czekania
na zaplanowane opóźnienie.
*/
function naPonowienie(zdarzenie: MouseEvent): void {
  const cel = zdarzenie.target;
  if (!(cel instanceof Element) || cel.closest('#btn-ponow') === null) return;
  const most = kanal();
  if (most === undefined) return;
  const pole = document.querySelector<HTMLInputElement>(
    '.we-panel[data-widok="w-blad"] input[data-adres-rdzenia]',
  );
  const adres = pole?.value.trim() ?? '';
  if (adres !== '' && isTauri()) {
    void wskazRdzen(adres);
    return;
  }
  most.wznowPolaczenie();
}

async function wskazRdzen(adres: string): Promise<void> {
  try {
    await invoke<WskazanieRdzenia>(POLECENIE_WSKAZ_RDZEN, { adres });
    globalThis.location.reload();
  } catch (blad) {
    const odmowa = blad as Partial<OdmowaWskazania> | undefined;
    oglos('Wskazanie rdzenia', odmowa?.zdanie ?? 'Powłoka odrzuciła wskazanie rdzenia.', 'blad');
  }
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
  await wejdz(most);
}

async function zarejestruj(
  most: Kanal,
  panel: HTMLElement,
  dokad: string,
  grupa: string | null,
): Promise<void> {
  // Login zapamiętany przy rejestracji: droga przez kod aktywacji wchodzi do
  // powłoki z pominięciem ekranu logowania, a pasek stanu nazywa Operatora.
  ostatniLogin = wartosc('rej-login');
  const wynik = await wywolaj(most, Command.AuthRegister, {
    login: ostatniLogin,
    email: wartosc('rej-email'),
    password: wartosc('rej-haslo'),
  });
  if (!wynik.udany) {
    odmowa(panel, 'usterki.naglowekKonto', wynik.blad);
    return;
  }
  /* `pendingVerification` fałszywe oznacza rejestrację bez listu z kodem —
     widok kodu byłby ślepym zaułkiem, wejście idzie od razu hasłem. */
  if (wynik.wynik?.pendingVerification !== true) {
    idz('logowanie', grupa);
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
  await wejdz(most);
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
  zapomnijTokenSesji();
  idz(dokad, grupa);
}

// Token sesji przejmuje z odpowiedzi rdzenia `pilnujTokenu` w warstwie protokołu.
async function wejdz(most: Kanal): Promise<void> {
  await Promise.all([
    wywolaj(most, Command.EnvironmentList, {}),
    wywolaj(most, Command.ModuleList, {}),
    wywolaj(most, Command.SessionList, {}),
  ]);
  odslonPowloke();
  await zwiazPowloke({ login: ostatniLogin, nazwaOkna: 'Danaco Console' }, most);
  zwiazCentrum(most);
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
