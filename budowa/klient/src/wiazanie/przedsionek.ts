// Wspólne wiązanie przedsionka: wejście w środowisko, szyna sesji, powrót do
// pracy i jedna droga z kafla w moduł. Cztery przedsionki mają tu wspólną
// warstwę, a pliki `przedsionek-talkin`, `-workspace`, `-codestudio`
// i `-multitaskingai` dokładają wyłącznie to, czym środowiska się różnią.
// Wykaz modułów i sesji bierze się wyłącznie z odpowiedzi rdzenia.
import {
  Command,
  SessionStatus,
  type Module,
  type Session,
  type SessionPresence,
} from '../../../shared/contract.ts';
import { miaraSesji } from '../model/miary.ts';
import type { Kanal } from '../protokol/kanal.ts';
import { tozsamoscKlienta } from '../protokol/tozsamosc-klienta.ts';
import { wywolaj } from '../protokol/wywolanie.ts';
import { oglos } from './ogloszenie.ts';
import { zdanieOdmowy } from './wejscie-odmowa.ts';

const NAGLOWEK = 'Przedsionek';

const ZAKRESY: readonly (SessionStatus | '')[] = [SessionStatus.Active, '', SessionStatus.Finished];

const STANY: Readonly<Record<SessionStatus, string>> = {
  [SessionStatus.Active]: 'czynna',
  [SessionStatus.Paused]: 'wstrzymana',
  [SessionStatus.Finished]: 'zakończona',
  [SessionStatus.Archived]: 'archiwalna',
};

/* Kafle projektu, repozytorium i przebiegu prowadzą we własne byty środowiska,
   więc rozpoznawanie modułu ich nie dotyka. */
const ZNAKI_OBCE: readonly string[] = ['data-projekt', 'data-repozytorium', 'data-przebieg'];

let zakres: SessionStatus | '' = SessionStatus.Active;
let sesje: Session[] = [];
let odpisy: SessionPresence[] = [];
let katalogModulow = new Map<string, Module>();
let kodOstatni = '';
/* Wzorzec wiersza pochodzi ze znacznika i ginie przy pierwszym pustym wykazie,
   więc szyna zapamiętuje go, zanim po raz pierwszy wymieni swoją zawartość. */
let wzorzecWiersza: HTMLElement | null = null;

/* Widok przedsionka powstaje w płótnie przy każdym wejściu w środowisko, więc
   wiązanie idzie za jego przebudową obserwatorem, nie odczytem znacznika. */
export function zwiazPrzedsionek(kanal: Kanal): void {
  sledzWidok(kanal);
  document.addEventListener('click', (zdarzenie) => {
    const cel = zdarzenie.target;
    if (!(cel instanceof Element)) return;
    const filtr = cel.closest<HTMLElement>('.pd-szyna-filtr .pd-filtr');
    if (filtr !== null) {
      zdarzenie.stopPropagation();
      przestawZakres(filtr);
      return;
    }
    const kafel = cel.closest<HTMLElement>('.cd-tresc--przedsionek .pd-siatka .pd-kafel');
    if (kafel !== null) {
      ujednolicKafel(kafel);
      return;
    }
    const wiersz = cel.closest<HTMLElement>('.pd-sesja[data-id-sesji]');
    if (wiersz !== null) wznowWiersz(kanal, wiersz, zdarzenie);
  }, true);
}

function sledzWidok(kanal: Kanal): void {
  const obserwator = new MutationObserver(() => {
    const kod = kodSrodowiska();
    if (kod === '' || kod === kodOstatni) return;
    kodOstatni = kod;
    void wejdz(kanal, kod);
  });
  obserwator.observe(document.body, { childList: true, subtree: true });
}

function kodSrodowiska(): string {
  const widok = document.querySelector<HTMLElement>('.cd-tresc--przedsionek');
  if (widok === null || widok.hidden) return '';
  const listwa = widok.querySelector<HTMLElement>('.pd-listwa [data-nowy-projekt-srodowisko]');
  return listwa?.dataset.nowyProjektSrodowisko ?? '';
}

async function wejdz(kanal: Kanal, kod: string): Promise<void> {
  const wynik = await wywolaj(kanal, Command.EnvironmentEnter, {
    environmentId: kod,
    clientId: tozsamoscKlienta().id,
  });
  if (!wynik.udany || wynik.wynik === undefined) {
    oglos(NAGLOWEK, zdanieOdmowy(wynik.blad), 'blad');
    /* Bez odpowiedzi rdzenia w szynie stoją wiersze i miara ze znacznika.
       Zostawione mówiłyby o cudzych sesjach jako o sesjach Operatora. */
    sesje = [];
    odpisy = [];
    katalogModulow = new Map();
    wypelnijSzyne('');
    return;
  }
  sesje = wynik.wynik.sessions;
  katalogModulow = spisModulow(wynik.wynik.modules);
  odpisy = await pobierzOdpisy(kanal);
  wypelnijSzyne(wynik.wynik.focusedSessionId ?? '');
}

function spisModulow(moduly: Module[]): Map<string, Module> {
  return new Map(moduly.map((modul) => [modul.code.toLowerCase(), modul]));
}

async function pobierzOdpisy(kanal: Kanal): Promise<SessionPresence[]> {
  const wynik = await wywolaj(kanal, Command.SessionList, { includePresence: true });
  if (!wynik.udany || wynik.wynik === undefined) {
    /* Bez odpisów wiersze stoją dalej, ale milczą o pracy sesji; odmowa
       nazwana Operatorowi odróżnia sesję bezczynną od nieznanej. */
    oglos(NAGLOWEK, `Rdzeń nie wydał stanu sesji. ${zdanieOdmowy(wynik.blad)}`, 'ostrzezenie');
    return [];
  }
  return wynik.wynik.presence ?? [];
}

/* Centrum otwiera moduł po znaczniku `data-modul`, a prototypy niosą w nim raz
   kod modułu, raz polską nazwę sekcji. Przedsionek wpisuje tam kod z wykazu
   rdzenia, zanim Centrum znacznik odczyta, więc cztery przedsionki wchodzą
   w moduł jedną drogą. */
/* Kafel bez pokrycia w wykazie zostaje nietknięty: odpowiada za niego wiązanie
   własne środowiska. */
function ujednolicKafel(kafel: HTMLElement): void {
  if (ZNAKI_OBCE.some((znak) => kafel.hasAttribute(znak))) return;
  const modul = rozpoznajModul(kafel);
  if (modul !== undefined) kafel.dataset.modul = modul.code;
}

/* Kod ze znacznika ma pierwszeństwo; kafel bez kodu albo z kodem spoza wykazu
   rozpoznaje się nazwą, bo to ona stoi Operatorowi na kaflu. */
function rozpoznajModul(kafel: HTMLElement): Module | undefined {
  const znak = (kafel.dataset.modul ?? '').trim().toLowerCase();
  const poKodzie = katalogModulow.get(znak);
  if (poKodzie !== undefined) return poKodzie;
  const napis = kafel.querySelector('.pd-kafel-nazwa')?.textContent ?? '';
  const nazwa = napis.trim().toLowerCase() || znak;
  if (nazwa === '') return undefined;
  return [...katalogModulow.values()].find(
    (modul) => modul.name.trim().toLowerCase() === nazwa,
  );
}

function przestawZakres(filtr: HTMLElement): void {
  const przyciski = [...document.querySelectorAll<HTMLElement>('.pd-szyna-filtr .pd-filtr')];
  const numer = przyciski.indexOf(filtr);
  if (numer < 0) return;
  zakres = ZAKRESY[numer] ?? '';
  for (const [pozycja, przycisk] of przyciski.entries()) {
    przycisk.setAttribute('aria-pressed', pozycja === numer ? 'true' : 'false');
  }
  wypelnijSzyne('');
}

function wypelnijSzyne(ogniskowana: string): void {
  const lista = document.querySelector<HTMLElement>('.cd-tresc--przedsionek .pd-szyna-lista');
  if (lista === null) return;
  const wzor = lista.querySelector<HTMLElement>('.pd-sesja');
  if (wzor !== null) wzorzecWiersza = wzor.cloneNode(true) as HTMLElement;
  const wzorzec = wzorzecWiersza;
  if (wzorzec === null) return;
  const widoczne = sesje.filter((sesja) => zakres === '' || sesja.status === zakres);
  lista.replaceChildren(...widoczne.map(
    (sesja) => zbudujWiersz(wzorzec, sesja, sesja.id === ogniskowana),
  ));
  if (widoczne.length === 0) lista.append(zdaniePustej(lista, zakres !== ''));
  opiszMiare(widoczne.length);
}

function zbudujWiersz(wzorzec: HTMLElement, sesja: Session, ogniskowana: boolean): HTMLElement {
  const wiersz = wzorzec.cloneNode(true) as HTMLElement;
  wiersz.dataset.idSesji = sesja.id;
  wiersz.dataset.stanSesji = sesja.status;
  const odpis = odpisy.find((pozycja) => pozycja.sessionId === sesja.id);
  if (odpis?.live !== true) wiersz.querySelector('.pt-tetno')?.remove();
  if (ogniskowana) wiersz.setAttribute('aria-current', 'true');
  else wiersz.removeAttribute('aria-current');
  const tytul = wiersz.querySelector('.pd-sesja-tytul');
  if (tytul !== null) tytul.textContent = sesja.title ?? sesja.id;
  const meta = wiersz.querySelector('.pd-sesja-meta');
  if (meta !== null) meta.textContent = opisSesji(sesja, odpis);
  return wiersz;
}

function opisSesji(sesja: Session, odpis: SessionPresence | undefined): string {
  const kod = odpis?.moduleCode;
  const modul = kod === undefined
    ? ''
    : katalogModulow.get(kod.toLowerCase())?.name ?? kod;
  const stan = odpis?.live === true ? 'pracuje' : STANY[sesja.status];
  return modul === '' ? stan : `${modul} · ${stan}`;
}

function zdaniePustej(lista: HTMLElement, zawezona: boolean): HTMLElement {
  const zdanie = lista.ownerDocument.createElement('p');
  zdanie.className = 'dn-pusty';
  zdanie.textContent = zawezona
    ? 'Żadna sesja nie ma tego stanu.'
    : 'Nie masz jeszcze sesji w tym środowisku.';
  return zdanie;
}

function opiszMiare(ile: number): void {
  const miara = document.querySelector(
    '.cd-tresc--przedsionek .dn-szyna-modulu-glowa .dn-plakietka',
  );
  if (miara !== null) miara.textContent = `sesji: ${ile}`;
  /* Listwa dolna niesie w prototypie miarę projektów i sesji naraz; kontrakt
     nie wiąże projektu ze środowiskiem, więc zostaje sama liczba sesji.
     Opisuje ją przedsionek, bo stoi także w środowisku bez modułów. */
  const listwa = document.querySelector('.cd-tresc--przedsionek .pd-listwa-meta');
  if (listwa !== null) listwa.textContent = miaraSesji(ile);
}

/* Sesja czynna idzie w otwarcie prowadzone przez Centrum; wstrzymana wymaga
   wpierw wznowienia, więc naciśnięcie wraca dopiero po odpowiedzi rdzenia. */
function wznowWiersz(kanal: Kanal, wiersz: HTMLElement, zdarzenie: Event): void {
  const stan = wiersz.dataset.stanSesji ?? '';
  if (stan === SessionStatus.Active || wiersz.dataset.wznowiona === 'tak') return;
  zdarzenie.stopPropagation();
  zdarzenie.preventDefault();
  void wznow(kanal, wiersz);
}

async function wznow(kanal: Kanal, wiersz: HTMLElement): Promise<void> {
  const idSesji = wiersz.dataset.idSesji ?? '';
  const wynik = await wywolaj(kanal, Command.SessionResume, { sessionId: idSesji });
  if (!wynik.udany || wynik.wynik === undefined) {
    oglos(NAGLOWEK, zdanieOdmowy(wynik.blad), 'blad');
    return;
  }
  const powiazanie = await wywolaj(kanal, Command.SessionBind, {
    sessionId: idSesji,
    clientId: tozsamoscKlienta().id,
  });
  if (!powiazanie.udany) {
    oglos(NAGLOWEK, zdanieOdmowy(powiazanie.blad), 'blad');
    return;
  }
  wiersz.dataset.wznowiona = 'tak';
  wiersz.dataset.stanSesji = wynik.wynik.session.status;
  wiersz.click();
}
