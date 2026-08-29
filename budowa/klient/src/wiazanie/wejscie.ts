/**
 * Wiązanie drogi wejścia: znacznik Właściciela po jednej stronie, komendy
 * rdzenia po drugiej.
 *
 * Ani jeden węzeł nie powstaje tutaj. Okna składa `okna/wejscie/montaz.js`,
 * odsłony przełącza `prototyp.js`, sprawdzenie pól przed wysłaniem prowadzi
 * `okna/przeplyw-wejscia.js`. Ten plik dokłada wyłącznie to, czego prototyp
 * z natury nie ma: rozmowę z rdzeniem i skutek jej wyniku w znaczniku.
 */

import { zwiazPowloke } from './powloka.ts';
import { zwiazStudio } from './studio.ts';
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

/** Ekran startowy trzymany atrybutem `data-czekaj`; `gotowe()` pozwala mu domknąć bieg. */
interface EkranStartowy {
  gotowe(): void;
}

/** Katalog treści okna związany narzędziami biblioteki; napisy pochodzą stamtąd, nie stąd. */
interface Katalog {
  tekst(sciezka: string): unknown;
}

/**
 * Styk z biblioteką Właściciela. Biblioteka to funkcje domknięte wystawiające
 * się na obiekcie globalnym — importu z nich nie ma. Rozpoznanie jest lokalne,
 * bez rozszerzania typu globalnego: inne wiązania robią to samo u siebie.
 */
interface Seam {
  DanacoKanal?: Kanal;
  DanacoNarzedzia?: { zwiaz(katalog: unknown): Katalog };
  DanacoWejscie?: { tresci?: unknown };
  dnPrzelaczWidok?: (widok: string, grupa?: string | null) => void;
  dnPrzejdz?: (widok: string, grupa?: string | null) => void;
  dnToast?: (tytul: string, tresc: string, rodzaj?: string) => void;
}

function seam(): Seam {
  return globalThis as unknown as Seam;
}

/* Sceny drogi wejścia w kolejności etapów. Nazwy są te, którymi posługuje się
   biblioteka: `data-po-polaczeniu` wskazuje „uwierzytelnienie”, a pasy działań
   prowadzą przez `data-idz` do „przygotowania”. */
const ETAPY: Readonly<Record<string, string>> = {
  uruchomienie: 'laczenie',
  dostep: 'uwierzytelnienie',
  przygotowanie: 'przygotowanie',
};

/** Odsłony okna dostępu, w których czynność główna woła rdzeń, a nie tylko przewija widok. */
type Odslona =
  | 'logowanie'
  | 'logowanie-blad'
  | 'rejestracja'
  | 'kod'
  | 'odzyskiwanie-adres'
  | 'odzyskiwanie-kod'
  | 'odzyskiwanie-haslo';

/* Droga potwierdzenia z listu. Odzyskiwanie zbiera ją w odsłonie kodu, a zużywa
   dopiero w odsłonie hasła — między nimi nie ma komendy, która by ją przeniosła. */
let drogaOdzyskania = '';

/** Token sesji wydany przez rdzeń; niesie go powitanie kolejnego połączenia. */
let tokenSesji = '';

/* Login, którym Operator wszedł — pasek stanu nie ma go skąd wziąć od rdzenia. */
let ostatniLogin = '';

export function zwiazWejscie(podany?: Kanal): void {
  oznaczSceny();
  zalozPrzejscie();
  document.addEventListener('click', naKlikniecie);
  gotowosc(() => {
    void powitaj(kanal(podany));
  });
}

/** Kanał rozmowy z rdzeniem: podany argumentem albo ten wystawiony przez `aplikacja.ts`. */
function kanal(podany?: Kanal): Kanal | undefined {
  return podany ?? seam().DanacoKanal;
}

/** Odkłada bieg do chwili, gdy montaż biblioteki zbudował okna i zgłosił gotowość dokumentu. */
function gotowosc(bieg: () => void): void {
  if (document.readyState === 'loading') {
    document.addEventListener('DOMContentLoaded', bieg, { once: true });
    return;
  }
  bieg();
}

/* ── Znacznik sceny ──────────────────────────────────────────────────────── */

/**
 * Sceny etapów dostają znaczniki widoku, którymi biblioteka już się posługuje.
 * Dokument klienta stawia je obok siebie bez nazw, bo w prototypie nazwy niosły
 * sekcje podglądu, których w produkcie nie ma. Bez nich `dnPrzelaczWidok` nie
 * miałby czego przełączać, a `przeplyw-wejscia.js` nie znalazłby swojej sceny.
 */
function oznaczSceny(): void {
  const miejsca = document.querySelectorAll<HTMLElement>('[data-wejscie-okno]');
  for (const miejsce of miejsca) {
    const nazwa = miejsce.dataset.wejscieOkno ?? '';
    const etap = ETAPY[nazwa];
    const scena = miejsce.closest<HTMLElement>('.we-scena');
    if (etap === undefined || scena === null) continue;
    scena.dataset.widok = etap;
    scena.dataset.grupaWidoku = 'etap';
    scena.dataset.widokAktywny = etap === 'laczenie' ? 'tak' : 'nie';
    if (etap === 'laczenie') scena.dataset.uruchomienieScena = '';
  }
}

/**
 * Szew `dnPrzejdz` — mechanika okna woła go po ostatnim etapie łączenia.
 * W prototypie stawiało go rusztowanie podglądu, którego produkt nie wczytuje.
 */
function zalozPrzejscie(): void {
  const s = seam();
  if (s.dnPrzejdz !== undefined) return;
  s.dnPrzejdz = (widok, grupa) => {
    s.dnPrzelaczWidok?.(widok, grupa ?? null);
  };
}

/* ── Czynności okna ──────────────────────────────────────────────────────── */

/**
 * Czynności okna niosą `data-idz`. Wiązanie wchodzi w fazie bąbelkowania, więc
 * sprawdzenie pól z `przeplyw-wejscia.js` — przechwytujące — rozstrzyga
 * pierwsze: gdy pola są puste, klik do nas nie dochodzi.
 */
function naKlikniecie(zdarzenie: MouseEvent): void {
  const cel = zdarzenie.target;
  if (!(cel instanceof Element)) return;
  const czynnosc = cel.closest<HTMLElement>('[data-idz]');
  if (czynnosc === null) return;
  if (czynnosc.tagName === 'A') zdarzenie.preventDefault();

  const dokad = czynnosc.dataset.idz ?? '';
  const grupa = czynnosc.dataset.idzGrupa ?? null;
  const pas = czynnosc.closest<HTMLElement>('.we-pas');
  const odslona = (pas?.dataset.widok ?? '') as Odslona;
  const most = kanal();

  if (pas === null || most === undefined || !wolaRdzen(odslona)) {
    idz(dokad, grupa);
    return;
  }
  zdarzenie.preventDefault();
  void wykonaj(most, odslona, dokad, grupa);
}

/** Czy odsłona ma za sobą komendę rdzenia; pozostałe czynności są samą nawigacją. */
function wolaRdzen(odslona: string): odslona is Odslona {
  return (
    odslona === 'logowanie' ||
    odslona === 'logowanie-blad' ||
    odslona === 'rejestracja' ||
    odslona === 'kod' ||
    odslona === 'odzyskiwanie-adres' ||
    odslona === 'odzyskiwanie-kod' ||
    odslona === 'odzyskiwanie-haslo'
  );
}

function idz(dokad: string, grupa: string | null): void {
  const s = seam();
  if (s.dnPrzejdz !== undefined) s.dnPrzejdz(dokad, grupa);
  else s.dnPrzelaczWidok?.(dokad, grupa);
}

/* ── Rozmowa z rdzeniem ──────────────────────────────────────────────────── */

/**
 * Powitanie połączenia. Do jego rozstrzygnięcia ekran startowy stoi wstrzymany
 * atrybutem `data-czekaj` — dopiero odpowiedź rdzenia pozwala mu domknąć bieg,
 * więc Operator nie ogląda okna, zanim wiadomo, czy jest z czym rozmawiać.
 */
async function powitaj(most: Kanal | undefined): Promise<void> {
  if (most === undefined) {
    odslonWariant('w-blad');
    domknijEkranStartowy();
    return;
  }
  const wynik = await zadajPowitanie(most, tozsamoscKlienta(), tokenSesji || undefined);
  odslonWariant(wynik.udany ? 'w-laczenie' : 'w-blad');
  domknijEkranStartowy();
  /*
  Przejście na etap uwierzytelnienia po dobiegnięciu animacji uruchomienia.

  Animacja gra na scenie łączenia i sama zgłasza koniec zdarzeniem
  `ekran-startowy-koniec`. Przejście przed nim zostawiało ją rysującą się na
  ekranie logowania — Operator widział znak kreślony po formularzu zamiast
  własnej sceny uruchomienia.

  Zabezpieczenie czasowe przepuszcza dalej także wtedy, gdy zdarzenie nie
  przyjdzie: okno bez wyjścia jest gorsze niż animacja ucięta.
  */
  if (!wynik.udany) return;
  let przeszedl = false;
  const dalej = (): void => {
    if (przeszedl) return;
    przeszedl = true;
    seam().dnPrzelaczWidok?.('uwierzytelnienie', 'etap');
  };
  document.addEventListener('ekran-startowy-koniec', dalej, { once: true });
  globalThis.setTimeout(dalej, 4000);
}

/** Odsłona okna uruchomienia zgodna z wynikiem powitania; przebieg etapów prowadzi biblioteka. */
function odslonWariant(widok: string): void {
  seam().dnPrzelaczWidok?.(widok, 'wariant');
}

function domknijEkranStartowy(): void {
  const pole = document.querySelector('[data-ekran-startowy]');
  const ekran = (pole as unknown as { ekranStartowy?: EkranStartowy } | null)?.ekranStartowy;
  ekran?.gotowe();
}

/** Wykonanie czynności głównej odsłony: komenda rdzenia, a po jej powodzeniu przejście dalej. */
async function wykonaj(
  most: Kanal,
  odslona: Odslona,
  dokad: string,
  grupa: string | null,
): Promise<void> {
  const panel = document.querySelector<HTMLElement>(`.we-panel[data-widok="${odslona}"]`);
  if (panel === null) return;

  switch (odslona) {
    case 'logowanie':
    case 'logowanie-blad':
      await zaloguj(most, panel, odslona);
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
      /* Kod odzyskiwania jest drogą potwierdzenia dla `auth.reset`, a nie osobną
         komendą — rdzeń sprawdza go dopiero razem z nowym hasłem. */
      drogaOdzyskania = zbierzKod(panel);
      idz(dokad, grupa);
      return;
    case 'odzyskiwanie-haslo':
      await ustawHaslo(most, panel, dokad, grupa);
      return;
  }
}

async function zaloguj(most: Kanal, panel: HTMLElement, odslona: Odslona): Promise<void> {
  const przedrostek = odslona === 'logowanie' ? 'log' : 'blad';
  const wynik = await wywolaj(most, Command.AuthLogin, {
    method: AuthMethodKind.Password,
    login: (ostatniLogin = wartosc(`${przedrostek}-login`)),
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
  /* Rejestracja sesji nie zakłada — konto czeka na potwierdzenie adresu, więc
     przejście do odsłony kodu jest jedyną drogą dalej. */
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
  /* Odzyskanie unieważnia tokeny wydane wcześniej, więc droga wraca do
     logowania — tak samo, jak zapowiada pas działań tej odsłony. */
  drogaOdzyskania = '';
  tokenSesji = '';
  idz(dokad, grupa);
}

/* ── Wejście do środowiska ───────────────────────────────────────────────── */

/**
 * Po wydaniu tokenu klient pyta rdzeń o to, z czego powłoka się składa, i
 * dopiero mając odpowiedzi odsłania ramę. Odsłonięcie przed nimi pokazywałoby
 * puste okno robocze.
 */
async function wejdz(most: Kanal, sesja: AuthSession | undefined): Promise<void> {
  if (sesja !== undefined) tokenSesji = sesja.token;
  await Promise.all([
    wywolaj(most, Command.EnvironmentList, {}),
    wywolaj(most, Command.ModuleList, {}),
    wywolaj(most, Command.SessionList, {}),
  ]);
  odslonPowloke();
  /* Powłoka i okno robocze przychodzą z prototypu wraz z jego treścią
     przykładową — nazwiskiem z przykładu, zmyślonymi miarami maszyny, cudzym
     dokumentem. Bez tych dwóch wywołań Operator ogląda przykład podany jako
     jego własna praca. */
  /* Nazwa okna: biblioteka powłoki niesie w belce nazwę okna logowania, bo jest
     wspólna dla wszystkich okien. Po wejściu stoi tu okno robocze. */
  await zwiazPowloke({ login: ostatniLogin, nazwaOkna: 'Danaco Console' }, most);
  zwiazStudio(most);
}

/** Zamiana sceny drogi wejścia na ramę aplikacji; oba węzły stoją w dokumencie od startu. */
function odslonPowloke(): void {
  document.querySelector('[data-rama-aplikacji]')?.removeAttribute('hidden');
  document.querySelector('[data-wejscie]')?.setAttribute('hidden', '');
}

/* ── Odmowy rdzenia ──────────────────────────────────────────────────────── */

/**
 * Odmowa logowania ma w oknie własną odsłonę z gotowym banerem, więc wystarczy
 * ją pokazać. Adres niepotwierdzony jest wyjątkiem: to nie jest zła para
 * login–hasło, tylko konto zatrzymane przed potwierdzeniem — droga prowadzi do
 * odsłony kodu, nie do komunikatu o nierozpoznanych danych.
 */
function odmowaLogowania(blad: ErrorInfo | undefined): void {
  if (powod(blad) === 'adres-niepotwierdzony') {
    seam().dnPrzelaczWidok?.('kod', 'stan');
    return;
  }
  seam().dnPrzelaczWidok?.('logowanie-blad', 'stan');
}

/**
 * Odmowa pozostałych czynności trafia w baner odsłony. Napis bierze się
 * z katalogu okna po nazwanym powodzie — treść od serwera zostaje w dzienniku,
 * bo niesie szczegóły, których Operatorowi pokazywać nie wolno.
 */
function odmowa(panel: HTMLElement, naglowek: string, blad: ErrorInfo | undefined): void {
  const klucz = powod(blad) === 'kolizja-danych' ? 'usterki.loginZajety' : '';
  const glowa = klucz === '' ? napis(naglowek) : napis(`${klucz}.glowa`);
  const tresc = klucz === '' ? napis('usterki.wiele') : napis(`${klucz}.tresc`);
  console.warn('[wejście] odmowa rdzenia', blad?.code, powod(blad));
  if (!wpiszWBaner(panel, glowa, tresc)) seam().dnToast?.(glowa, tresc, 'blad');
}

/** Nazwany powód odmowy podany przez rdzeń w danych diagnostycznych; puste, gdy go nie ma. */
function powod(blad: ErrorInfo | undefined): string {
  const dane = blad?.details as { powod?: unknown } | undefined;
  return typeof dane?.powod === 'string' ? dane.powod : '';
}

/**
 * Wypełnienie banera, który odsłona już ma: głowa i zdanie pod nią. Baner
 * zmienia odmianę na błąd klasą, którą biblioteka rozumie. Gdy odsłona banera
 * nie ma, wiązanie go nie dostawia — komunikat idzie wtedy powiadomieniem.
 */
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

/* ── Odczyt pól ──────────────────────────────────────────────────────────── */

function wartosc(id: string): string {
  const pole = document.getElementById(id);
  return pole instanceof HTMLInputElement ? pole.value.trim() : '';
}

/** Zgoda „nie wyloguj mnie” z pola wyboru odsłony; brak pola znaczy brak zgody. */
function trwalaSesja(panel: HTMLElement): boolean {
  const pole = panel.querySelector<HTMLInputElement>('.au-opcja input[type="checkbox"]');
  return pole?.checked ?? false;
}

/** Kod potwierdzający złożony z sześciu pól odsłony w kolejności ich stania. */
function zbierzKod(panel: HTMLElement): string {
  const pola = panel.querySelectorAll<HTMLInputElement>('.au-kod-pole');
  let kod = '';
  for (const pole of pola) kod += pole.value.trim();
  return kod;
}

/* ── Napisy ──────────────────────────────────────────────────────────────── */

/** Napis z katalogu treści okna; ani jedno zdanie dla Operatora nie stoi w tym pliku. */
function napis(sciezka: string): string {
  const s = seam();
  const katalog = s.DanacoNarzedzia?.zwiaz(s.DanacoWejscie?.tresci ?? {});
  const wartoscNapisu = katalog?.tekst(sciezka);
  return typeof wartoscNapisu === 'string' ? wartoscNapisu : '';
}

zwiazWejscie();
