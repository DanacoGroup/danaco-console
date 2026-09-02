// Wiązanie przedsionka środowiska: wejście w środowisko, szyna sesji i powrót
// do pracy. Wykaz modułów i sesji bierze się wyłącznie z odpowiedzi rdzenia.
import {
  Command,
  SessionStatus,
  type Module,
  type Session,
  type SessionPresence,
} from '../../../shared/contract.ts';
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

let zakres: SessionStatus | '' = SessionStatus.Active;
let sesje: Session[] = [];
let odpisy: SessionPresence[] = [];
let nazwyModulow = new Map<string, string>();
let kodOstatni = '';

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
    return;
  }
  sesje = wynik.wynik.sessions;
  nazwyModulow = spisModulow(wynik.wynik.modules);
  odpisy = await pobierzOdpisy(kanal);
  wypelnijSzyne(wynik.wynik.focusedSessionId ?? '');
}

function spisModulow(moduly: Module[]): Map<string, string> {
  return new Map(moduly.map((modul) => [modul.code, modul.name]));
}

async function pobierzOdpisy(kanal: Kanal): Promise<SessionPresence[]> {
  const wynik = await wywolaj(kanal, Command.SessionList, { includePresence: true });
  return wynik.wynik?.presence ?? [];
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
  const wzor = lista?.querySelector<HTMLElement>('.pd-sesja');
  if (lista === null || wzor === null || wzor === undefined) return;
  const wzorzec = wzor.cloneNode(true) as HTMLElement;
  const widoczne = sesje.filter((sesja) => zakres === '' || sesja.status === zakres);
  lista.replaceChildren(...widoczne.map(
    (sesja) => zbudujWiersz(wzorzec, sesja, sesja.id === ogniskowana),
  ));
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
  const modul = odpis?.moduleCode === undefined
    ? ''
    : nazwyModulow.get(odpis.moduleCode) ?? odpis.moduleCode;
  const stan = odpis?.live === true ? 'pracuje' : STANY[sesja.status];
  return modul === '' ? stan : `${modul} · ${stan}`;
}

function opiszMiare(ile: number): void {
  const miara = document.querySelector(
    '.cd-tresc--przedsionek .dn-szyna-modulu-glowa .dn-plakietka',
  );
  if (miara !== null) miara.textContent = `sesji: ${ile}`;
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
