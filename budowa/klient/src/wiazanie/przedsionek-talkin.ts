/**
 * Przedsionek TalkIn. Wspólne wiązanie zostawia kafle bez miary, listwę bez
 * stanu środowiska, a sesję bez czynności poza wznowieniem; ten plik dopełnia
 * te braki odpowiedziami rdzenia. Komendy: `module.list`, `session.list`,
 * `session.stop`, `session.archive`, `session.delete`.
 */

import {
  Command,
  type Module,
  type Session,
  type SessionPresence,
} from '../../../shared/contract.ts';
import { miaraSesji } from '../model/miary.ts';
import type { Kanal } from '../protokol/kanal.ts';
import { wywolaj } from '../protokol/wywolanie.ts';
import { oglos } from './ogloszenie.ts';

const NAGLOWEK = 'Przedsionek TalkIn';
const KOD_SRODOWISKA = 'talkin';
const NAPIS_USUN = 'Usuń sesję';
const NAPIS_POTWIERDZ = 'Potwierdź usunięcie';

let sesjeSrodowiska: Session[] = [];
let odpisySesji: SessionPresence[] = [];
let wskazanaSesja = '';
/* Usunięcie sesji jest jedyną drogą utraty jej zapisu, więc pierwsze
   naciśnięcie wyłącznie uzbraja przycisk, a wykonuje dopiero drugie. */
let usuniecieUzbrojone = false;
let przedsionekOtwarty = false;

/** Dopełnia przedsionek TalkIn danymi rdzenia i wiąże czynności na sesji. */
export function zwiazPrzedsionekTalkin(kanal: Kanal): void {
  sledzOtwarcie(kanal);
  document.addEventListener('pointerdown', (zdarzenie) => {
    const cel = zdarzenie.target;
    if (!(cel instanceof Element) || przedsionekTalkin() === null) return;
    const wiersz = cel.closest<HTMLElement>('.pd-sesja[data-id-sesji]');
    if (wiersz !== null) wskaz(wiersz.dataset.idSesji ?? '');
  }, true);
  document.addEventListener('click', (zdarzenie) => {
    const cel = zdarzenie.target;
    if (!(cel instanceof Element) || przedsionekTalkin() === null) return;
    const przycisk = cel.closest<HTMLElement>('.pd-szyna-czynnosci .dn-btn');
    if (przycisk === null) return;
    zdarzenie.stopPropagation();
    zdarzenie.preventDefault();
    void wykonajCzynnosc(kanal, przycisk);
  }, true);
}

/* Widok przedsionka powstaje w płótnie dopiero przy wejściu w środowisko i ginie
   przy powrocie do Centrum, więc wiązanie idzie za przebudową płótna. */
function sledzOtwarcie(kanal: Kanal): void {
  const obserwator = new MutationObserver(() => {
    const widok = przedsionekTalkin();
    if (widok === null) {
      przedsionekOtwarty = false;
      return;
    }
    if (przedsionekOtwarty) return;
    przedsionekOtwarty = true;
    wskaz('');
    void odswiez(kanal);
  });
  obserwator.observe(document.body, { childList: true, subtree: true });
}

function przedsionekTalkin(): HTMLElement | null {
  const widok = document.querySelector<HTMLElement>('.cd-tresc--przedsionek');
  if (widok === null || widok.hidden) return null;
  const znak = widok.querySelector<HTMLElement>('[data-nowy-projekt-srodowisko]');
  return znak?.dataset.nowyProjektSrodowisko === KOD_SRODOWISKA ? widok : null;
}

async function odswiez(kanal: Kanal): Promise<void> {
  const [moduly, wykazSesji] = await Promise.all([pobierzModuly(kanal), pobierzSesje(kanal)]);
  if (moduly === null || !wykazSesji) return;
  const widok = przedsionekTalkin();
  if (widok === null) return;
  opiszKafle(widok, moduly);
  opiszStanSrodowiska(widok);
  ustawCzynnosci(widok);
}

/* Wykaz modułów bierze się z rdzenia, a nie z kafli prototypu: prototyp niesie
   dziewięć kafli TalkIn niezależnie od tego, co środowisko naprawdę udostępnia. */
async function pobierzModuly(kanal: Kanal): Promise<Module[] | null> {
  const wynik = await wywolaj(kanal, Command.ModuleList, { environmentId: KOD_SRODOWISKA });
  if (!wynik.udany || wynik.wynik === undefined) {
    zglosOdmowe('Rdzeń nie wydał wykazu modułów środowiska.', wynik.blad?.message);
    return null;
  }
  return wynik.wynik.modules;
}

async function pobierzSesje(kanal: Kanal): Promise<boolean> {
  const wynik = await wywolaj(kanal, Command.SessionList, { includePresence: true });
  if (!wynik.udany || wynik.wynik === undefined) {
    zglosOdmowe('Rdzeń nie wydał wykazu sesji.', wynik.blad?.message);
    return false;
  }
  sesjeSrodowiska = wynik.wynik.sessions.filter(
    (sesja) => sesja.environmentCode === KOD_SRODOWISKA,
  );
  odpisySesji = wynik.wynik.presence ?? [];
  return true;
}

/* Miara kafla stała w prototypie jako liczba wymyślona („1 240 pozycji”),
   więc każdy kafel dostaje liczbę sesji tego modułu albo zdanie o jej braku. */
function opiszKafle(widok: HTMLElement, moduly: Module[]): void {
  const katalog = new Map(moduly.map((modul) => [modul.code, modul]));
  const liczby = policzSesjeModulow();
  for (const kafel of widok.querySelectorAll<HTMLElement>('.pd-siatka .pd-kafel')) {
    const kod = kafel.dataset.modul ?? '';
    const modul = katalog.get(kod);
    const miara = kafel.querySelector('.pd-kafel-meta') ?? dopiszMiareKafla(kafel);
    if (modul === undefined) {
      kafel.setAttribute('aria-disabled', 'true');
      miara.textContent = 'poza środowiskiem';
      continue;
    }
    const bezOkna = modul.operationalWindowCodes.length === 0;
    kafel.setAttribute('aria-disabled', String(bezOkna));
    miara.textContent = bezOkna ? 'bez okna operacyjnego' : opisMiaryModulu(liczby.get(kod) ?? 0);
  }
}

/* Wspólne wypełnianie kafli zdejmuje miarę prototypu, bo nie zna liczby sesji
   modułu; przedsionek TalkIn zna ją z odpisów i stawia węzeł z powrotem. */
function dopiszMiareKafla(kafel: HTMLElement): Element {
  const miara = kafel.ownerDocument.createElement('span');
  miara.className = 'pd-kafel-meta';
  kafel.appendChild(miara);
  return miara;
}

function opisMiaryModulu(ile: number): string {
  return ile === 0 ? 'bez sesji' : miaraSesji(ile);
}

/* Moduł sesji niesie odpis okna ogniskowanego, więc liczba sesji modułu jest
   liczbą odpisów wskazujących ten moduł wśród sesji tego środowiska. */
function policzSesjeModulow(): Map<string, number> {
  const znane = new Set(sesjeSrodowiska.map((sesja) => sesja.id));
  const liczby = new Map<string, number>();
  for (const odpis of odpisySesji) {
    if (!znane.has(odpis.sessionId) || odpis.moduleCode === undefined) continue;
    liczby.set(odpis.moduleCode, (liczby.get(odpis.moduleCode) ?? 0) + 1);
  }
  return liczby;
}

/* Stan środowiska nie stał nigdzie w przedsionku: Operator widział liczbę sesji,
   ale nie to, ile z nich pracuje i w ilu oknach trwa odpowiedź. */
function opiszStanSrodowiska(widok: HTMLElement): void {
  const napis = widok.querySelector('.pd-listwa-stan');
  if (napis === null) return;
  const znane = new Set(sesjeSrodowiska.map((sesja) => sesja.id));
  const wlasne = odpisySesji.filter((odpis) => znane.has(odpis.sessionId));
  const pracujace = wlasne.filter((odpis) => odpis.live).length;
  const strumienie = wlasne.reduce((suma, odpis) => suma + odpis.streamingWindowCount, 0);
  napis.textContent = `${miaraSesji(sesjeSrodowiska.length)} · pracuje ${pracujace}`
    + ` · odpowiedź trwa w ${strumienie} oknach`;
}

function wskaz(idSesji: string): void {
  wskazanaSesja = idSesji;
  usuniecieUzbrojone = false;
  const widok = przedsionekTalkin();
  if (widok === null) return;
  for (const wiersz of widok.querySelectorAll<HTMLElement>('.pd-sesja[data-id-sesji]')) {
    wiersz.dataset.sesjaWskazana = wiersz.dataset.idSesji === idSesji ? 'tak' : 'nie';
  }
  ustawCzynnosci(widok);
}

/* Czynność bez wskazanej sesji nie ma przedmiotu, a zatrzymanie tury ma sens
   wyłącznie dla sesji trwającej na rdzeniu — stąd dwa różne warunki. */
function ustawCzynnosci(widok: HTMLElement): void {
  const odpis = odpisySesji.find((pozycja) => pozycja.sessionId === wskazanaSesja);
  const wybrana = wskazanaSesja !== '';
  ustawPrzycisk(widok, 'sesjaZatrzymaj', wybrana && odpis?.live === true);
  ustawPrzycisk(widok, 'sesjaArchiwizuj', wybrana);
  const usun = ustawPrzycisk(widok, 'sesjaUsun', wybrana);
  if (usun !== null) usun.textContent = usuniecieUzbrojone ? NAPIS_POTWIERDZ : NAPIS_USUN;
}

function ustawPrzycisk(
  widok: HTMLElement,
  czynnosc: string,
  czynny: boolean,
): HTMLButtonElement | null {
  const przycisk = widok.querySelector<HTMLButtonElement>(
    `.pd-szyna-czynnosci [data-czynnosc-sesji="${czynnosc}"]`,
  );
  if (przycisk === null) return null;
  przycisk.disabled = !czynny;
  return przycisk;
}

async function wykonajCzynnosc(kanal: Kanal, przycisk: HTMLElement): Promise<void> {
  const czynnosc = przycisk.dataset.czynnoscSesji ?? '';
  if (wskazanaSesja === '') return;
  if (czynnosc === 'sesjaZatrzymaj') await zatrzymaj(kanal);
  if (czynnosc === 'sesjaArchiwizuj') await archiwizuj(kanal);
  if (czynnosc === 'sesjaUsun') await usun(kanal);
}

async function zatrzymaj(kanal: Kanal): Promise<void> {
  const wynik = await wywolaj(kanal, Command.SessionStop, { sessionId: wskazanaSesja });
  if (!wynik.udany || wynik.wynik === undefined) {
    zglosOdmowe('Rdzeń nie zatrzymał tur sesji.', wynik.blad?.message);
    return;
  }
  oglos(NAGLOWEK, `Zatrzymano tury w ${wynik.wynik.stoppedWindowIds.length} oknach.`);
  await odswiez(kanal);
}

async function archiwizuj(kanal: Kanal): Promise<void> {
  const idSesji = wskazanaSesja;
  const wynik = await wywolaj(kanal, Command.SessionArchive, { sessionIds: [idSesji] });
  if (!wynik.udany || wynik.wynik === undefined) {
    zglosOdmowe('Rdzeń nie przeniósł sesji do archiwum.', wynik.blad?.message);
    return;
  }
  if (wynik.wynik.archivedIds.length === 0) {
    zglosOdmowe('Rdzeń nie przeniósł sesji do archiwum.', undefined);
    return;
  }
  zdejmijWiersz(idSesji);
  oglos(NAGLOWEK, 'Sesja stoi w archiwum; jej zapis został w całości.');
  await odswiez(kanal);
}

async function usun(kanal: Kanal): Promise<void> {
  if (!usuniecieUzbrojone) {
    usuniecieUzbrojone = true;
    const widok = przedsionekTalkin();
    if (widok !== null) ustawCzynnosci(widok);
    oglos(NAGLOWEK, 'Usunięcie sesji jest nieodwracalne. Naciśnij ponownie, aby je wykonać.',
      'ostrzezenie');
    return;
  }
  const idSesji = wskazanaSesja;
  const wynik = await wywolaj(kanal, Command.SessionDelete, {
    sessionIds: [idSesji],
    confirm: true,
  });
  usuniecieUzbrojone = false;
  if (!wynik.udany || wynik.wynik === undefined) {
    zglosOdmowe('Rdzeń nie usunął sesji.', wynik.blad?.message);
    return;
  }
  if (wynik.wynik.deletedCount === 0) {
    zglosOdmowe('Rdzeń nie znalazł tej sesji w historii.', undefined);
    return;
  }
  zdejmijWiersz(idSesji);
  oglos(NAGLOWEK, 'Sesja usunięta wraz z całym zapisem.');
  await odswiez(kanal);
}

/* Szyna sesji należy do wspólnego wiązania i nie odświeża się sama po czynności,
   więc wiersz sesji, której już nie ma, schodzi z ekranu tutaj. */
function zdejmijWiersz(idSesji: string): void {
  const widok = przedsionekTalkin();
  widok?.querySelector(`.pd-sesja[data-id-sesji="${idSesji}"]`)?.remove();
  wskaz('');
}

function zglosOdmowe(zdanie: string, opisBledu: string | undefined): void {
  oglos(NAGLOWEK, opisBledu === undefined ? zdanie : `${zdanie} ${opisBledu}`, 'ostrzezenie');
}
