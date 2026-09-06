// Wiązanie drogi wejścia z komendami rdzenia. Znacznik i przełączanie widoków
// należą do biblioteki `design/zasoby/okna/wejscie/`.
import {
  AuthChangeReason,
  AuthMethodKind,
  Command,
  EventType,
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
import { napis, urzadzenieTrwale } from './wejscie-katalog.ts';
import { odmowaLogowania, odmowaRejestracji, pokazOdmowe } from './wejscie-odmowa.ts';
import { przygotujSrodowisko, zwiazPasPrzygotowania } from './wejscie-przygotowanie.ts';
import { powlokaStoi, stanRdzenia, wskazRdzen, wskazanieRdzenia } from './wejscie-powloka.ts';
import { zapamietajMetody, zwiazMetody } from './logowanie-metody.ts';
import { zwiazPrzedsionek } from './przedsionek.ts';
import { zwiazPrzedsionekTalkin } from './przedsionek-talkin.ts';
import { zwiazPrzedsionekWorkspace } from './przedsionek-workspace.ts';
import { zwiazPrzedsionekCodeStudio } from './przedsionek-codestudio.ts';
import { zwiazPrzedsionekMultitaskingAI } from './przedsionek-multitaskingai.ts';

interface EkranStartowy {
  gotowe(): void;
}

// Skrypty biblioteki są funkcjami domkniętymi, nie modułami — importu z nich nie ma.
interface Styk {
  DanacoKanal?: Kanal;
  dnPrzelaczWidok?: (widok: string, grupa?: string | null) => void;
  dnPrzejdz?: (widok: string, grupa?: string | null) => void;
}

function styk(): Styk {
  return globalThis as unknown as Styk;
}

const ETAPY: Readonly<Record<string, string>> = {
  uruchomienie: 'laczenie',
  dostep: 'uwierzytelnienie',
  przygotowanie: 'przygotowanie',
};

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
let powitanoRaz = false;
let powlokaOdsloniona = false;

export function zwiazWejscie(podany?: Kanal): void {
  oznaczObszary();
  zalozPrzejscie();
  document.addEventListener('click', naKlikniecie);
  document.addEventListener('click', naPonowienie);
  document.addEventListener('keydown', naKlawisz);
  zwiazStanPolaczenia(kanal(podany));
  zwiazZmianeUwierzytelnienia();
  const most = kanal(podany);
  poMontazu(() => {
    if (most === undefined) return;
    zwiazMetody(most);
    zwiazPasPrzygotowania(most);
    zwiazPrzedsionek(most);
    // Każde środowisko wnosi własne miary i czynności ponad wspólną szyną sesji;
    // wiązania stoją po wspólnym, bo przepuszczają do niego kliknięcia.
    zwiazPrzedsionekTalkin(most);
    zwiazPrzedsionekWorkspace(most);
    zwiazPrzedsionekCodeStudio(most);
    zwiazPrzedsionekMultitaskingAI(most);
  });
  gotowosc(() => {
    void powitaj(most);
  });
}

// Zerwanie po powitaniu pokazuje wariant błędu; powrót łączności odgrywa łączenie.
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

function odegrajLaczenie(wariant: 'w-laczenie' | 'w-token' = 'w-laczenie'): void {
  ustawWariant(wariant);
  document
    .querySelector('[data-ekran-startowy]')
    ?.dispatchEvent(new CustomEvent('ekran-startowy-koniec', { bubbles: true }));
}

// Rdzeń adresuje `auth.changed` do konta, nie do jednego urządzenia sesji.
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

function poMontazu(bieg: () => void): void {
  if (document.querySelector('[data-wejscie-okno="dostep"] .we-okno') !== null) {
    bieg();
    return;
  }
  document.addEventListener('wejscie-gotowe', bieg, { once: true });
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
    // Okno przygotowania ma jedną odsłonę, więc jest czynna zawsze.
    if (etap === 'przygotowanie') miejsce.dataset.widokAktywny = 'tak';
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

// Enter w polu wchodzi tak samo jak naciśnięcie czynności głównej pasa.
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
  domknijEkranStartowy();
  if (!wynik.udany) {
    ustawWariant('w-blad');
    return;
  }
  // Panel z celem przejścia staje w gnieździe dopiero po powrocie powitania.
  odegrajLaczenie(wynik.wynik?.authenticated === true ? 'w-token' : 'w-laczenie');
}

function ustawWariant(widok: string): void {
  styk().dnPrzelaczWidok?.(widok, 'wariant');
  if (widok === 'w-blad') void opiszWskazanieRdzenia();
}

// Wariant błędu nazywa serwer, którego powłoka nie doprosiła; poza powłoką wskazania nie ma.
async function opiszWskazanieRdzenia(): Promise<void> {
  if (!powlokaStoi()) return;
  const lid = document.querySelector<HTMLElement>('.we-panel[data-widok="w-blad"] .we-lid');
  if (lid === null) return;
  lid.dataset.lidPierwotny ??= lid.textContent ?? '';
  const wskazanie = await wskazanieRdzenia();
  if (wskazanie === null) return;
  const stan = await stanRdzenia();
  const zdanie = wskazanie.adres === null
    ? 'Powłoka nie ma wskazania rdzenia.'
    : `Rdzeń wskazany: ${wskazanie.adres}. ${stan?.opis ?? ''}`.trim();
  lid.textContent = `${lid.dataset.lidPierwotny} ${zdanie}`;
}

/* Wskazanie wpisane w pole adresu idzie do powłoki, która sprawdza łączność
   i zapisuje je trwale; bez wpisu ponowienie łączy pod adres obowiązujący. */
function naPonowienie(zdarzenie: MouseEvent): void {
  const cel = zdarzenie.target;
  if (!(cel instanceof Element) || cel.closest('#btn-ponow') === null) return;
  const most = kanal();
  if (most === undefined) return;
  const pole = document.querySelector<HTMLInputElement>(
    '.we-panel[data-widok="w-blad"] input[data-adres-rdzenia]',
  );
  const adres = pole?.value.trim() ?? '';
  if (adres !== '' && powlokaStoi()) {
    void przyjmijWskazanie(adres);
    return;
  }
  most.wznowPolaczenie();
}

async function przyjmijWskazanie(adres: string): Promise<void> {
  try {
    await wskazRdzen(adres);
    globalThis.location.reload();
  } catch (blad) {
    const odmowa = blad as { zdanie?: string } | undefined;
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
    deviceId: urzadzenieTrwale(),
    keepSignedIn: trwalaSesja(panel),
  });
  if (!wynik.udany) {
    odmowaLogowania(wynik.blad);
    return;
  }
  zapamietajMetody(wynik.wynik?.methods);
  await wejdz(most);
}

async function zarejestruj(
  most: Kanal,
  panel: HTMLElement,
  dokad: string,
  grupa: string | null,
): Promise<void> {
  // Droga przez kod aktywacji pomija ekran logowania, a pasek stanu nazywa Operatora.
  ostatniLogin = wartosc('rej-login');
  const wynik = await wywolaj(most, Command.AuthRegister, {
    login: ostatniLogin,
    email: wartosc('rej-email'),
    password: wartosc('rej-haslo'),
    deviceId: urzadzenieTrwale(),
  });
  if (!wynik.udany) {
    odmowaRejestracji(panel, wynik.blad);
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
    deviceId: urzadzenieTrwale(),
    keepSignedIn: trwalaSesja(panel),
  });
  if (!wynik.udany) {
    pokazOdmowe(panel, napis('usterki.naglowekKod'), wynik.blad);
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
    pokazOdmowe(panel, napis('usterki.naglowekKod'), wynik.blad);
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
    pokazOdmowe(panel, napis('usterki.naglowekHaslo'), wynik.blad);
    return;
  }
  // `auth.reset` unieważnia tokeny wydane wcześniej.
  drogaOdzyskania = '';
  zapomnijTokenSesji();
  idz(dokad, grupa);
}

/* Wejście po wydaniu tokenu prowadzi przez etap przygotowania: dopiero jego
   domknięcie odsłania powłokę, w której stoi wybór środowiska. */
async function wejdz(most: Kanal): Promise<void> {
  idz('przygotowanie', 'etap');
  await przygotujSrodowisko(most, () => {
    void otworzPowloke(most);
  });
}

async function otworzPowloke(most: Kanal): Promise<void> {
  if (powlokaOdsloniona) return;
  powlokaOdsloniona = true;
  await wywolaj(most, Command.EnvironmentList, { includeModules: true });
  document.querySelector('[data-rama-aplikacji]')?.removeAttribute('hidden');
  document.querySelector('[data-wejscie]')?.setAttribute('hidden', '');
  await zwiazPowloke({ login: ostatniLogin, nazwaOkna: 'Danaco Console' }, most);
  zwiazCentrum(most);
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
