// Okna robocze wraz z kartami, które w nich stoją: kontener jednego działania
// Operatora. Kontrakt wykazu okien roboczych nie ma — odtwarza się ze stanu.
import {
  ChangeKind,
  Command,
  EventType,
  WindowStatus,
  type PanelSection,
  type Session,
  type SessionPresence,
  type Window,
} from '../../../shared/contract.ts';
import type { Kanal } from '../protokol/kanal.ts';
import { wywolaj } from '../protokol/wywolanie.ts';
import { oglos } from './ogloszenie.ts';
import { zglosUchwyt } from './zdarzenia.ts';

export interface KartaRobocza {
  id: string;
  kodModulu: string;
  idOknaKomunikacji: string;
  nazwa: string;
}

export interface OknoRobocze {
  id: string;
  nazwa: string;
  idSesji: string;
  kartaBiezaca: string;
  karty: string[];
}

const OKNA: OknoRobocze[] = [];
// Karta przeżywa zamknięcie pasma, więc okno wymienia ją identyfikatorem.
const KARTY = new Map<string, KartaRobocza>();
let biezace = '';
let licznik = 0;
let licznikKart = 0;
// Zamknięcie karty sięga rdzenia, a czynność okna kanału nie niesie.
let kanalRdzenia: Kanal | null = null;

const NAZWA_POCZATKOWA = 'Centrum dowodzenia';

export function wskazKanalRdzenia(kanal: Kanal): void {
  kanalRdzenia = kanal;
}

export function oknaRobocze(): OknoRobocze[] {
  return OKNA;
}

export function oknoKarty(idKarty: string): OknoRobocze | undefined {
  return OKNA.find((okno) => okno.karty.includes(idKarty));
}

export function oknoSesji(idSesji: string): OknoRobocze | undefined {
  if (idSesji === '') return undefined;
  return OKNA.find((okno) => okno.idSesji === idSesji);
}

export function przypiszSesjeOkna(idOkna: string, idSesji: string, nazwa = ''): boolean {
  const okno = OKNA.find((kandydat) => kandydat.id === idOkna);
  if (okno === undefined || idSesji === '' || (okno.idSesji !== '' && okno.idSesji !== idSesji)) {
    return false;
  }
  okno.idSesji = idSesji;
  if (nazwa !== '' && okno.kartaBiezaca === '') okno.nazwa = nazwa;
  return true;
}

export function zdejmijSesjeOkien(idSesji: string): string[] {
  const dotkniete: string[] = [];
  for (const okno of OKNA) {
    if (okno.idSesji !== idSesji || idSesji === '') continue;
    okno.idSesji = '';
    dotkniete.push(okno.id);
  }
  return dotkniete;
}

export function oknoBiezace(): OknoRobocze {
  if (OKNA.length === 0) otworzOkno();
  return OKNA.find((okno) => okno.id === biezace) ?? OKNA[0];
}

export function otworzOkno(): OknoRobocze {
  licznik += 1;
  const okno: OknoRobocze = {
    id: 'okr-' + String(licznik), nazwa: NAZWA_POCZATKOWA, idSesji: '', kartaBiezaca: '', karty: [],
  };
  OKNA.push(okno);
  biezace = okno.id;
  return okno;
}

export function przelaczOkno(id: string): OknoRobocze | undefined {
  const okno = OKNA.find((kandydat) => kandydat.id === id);
  if (okno === undefined) return undefined;
  biezace = okno.id;
  return okno;
}

export function zamknijOkno(): OknoRobocze {
  const numer = OKNA.findIndex((okno) => okno.id === biezace);
  if (OKNA.length === 1 || numer < 0) {
    const jedyne = OKNA[0] ?? otworzOkno();
    for (const idKarty of jedyne.karty) zwolnijKarte(idKarty);
    jedyne.nazwa = NAZWA_POCZATKOWA;
    jedyne.idSesji = '';
    jedyne.kartaBiezaca = '';
    jedyne.karty = [];
    biezace = jedyne.id;
    return jedyne;
  }
  for (const idKarty of OKNA[numer].karty) zwolnijKarte(idKarty);
  OKNA.splice(numer, 1);
  const poprzednie = OKNA[Math.max(0, numer - 1)];
  biezace = poprzednie.id;
  return poprzednie;
}

export function zapiszKarte(kodModulu: string, nazwa: string, idOknaKomunikacji = ''): string {
  const okno = oknoBiezace();
  if (kodModulu === '') {
    okno.kartaBiezaca = '';
    okno.nazwa = NAZWA_POCZATKOWA;
    return '';
  }
  const karty = kartyOkna(okno).filter((kandydat) => kandydat.kodModulu === kodModulu);
  const karta = (idOknaKomunikacji === ''
    ? karty.find((kandydat) => kandydat.id === okno.kartaBiezaca) ?? karty[0]
    : karty.find((kandydat) => kandydat.idOknaKomunikacji === idOknaKomunikacji))
    ?? zalozKarte(okno, kodModulu, idOknaKomunikacji, nazwa);
  karta.nazwa = nazwa;
  okno.kartaBiezaca = karta.id;
  okno.nazwa = nazwa;
  return karta.id;
}

export function otworzKarte(kodModulu: string, nazwa: string, idOknaKomunikacji = ''): string {
  const okno = oknoBiezace();
  const karta = zalozKarte(okno, kodModulu, idOknaKomunikacji, nazwa);
  okno.kartaBiezaca = karta.id;
  okno.nazwa = nazwa;
  return karta.id;
}

export function przypiszOknoKomunikacji(idKarty: string, idOknaKomunikacji: string): boolean {
  const karta = KARTY.get(idKarty);
  if (karta === undefined || idOknaKomunikacji === '') return false;
  karta.idOknaKomunikacji = idOknaKomunikacji;
  return true;
}

export function kartaOkna(idKarty: string): KartaRobocza | undefined {
  return KARTY.get(idKarty);
}

export function wskazKarte(idKarty: string): KartaRobocza | undefined {
  const okno = oknoBiezace();
  const karta = KARTY.get(idKarty);
  if (karta === undefined || !okno.karty.includes(idKarty)) return undefined;
  okno.kartaBiezaca = karta.id;
  okno.nazwa = karta.nazwa;
  return karta;
}

export function zastosujZmianeOkna(zmiana: ChangeKind, okno: Window): string[] {
  const zdjete: string[] = [];
  for (const karta of KARTY.values()) {
    if (karta.idOknaKomunikacji !== okno.id) continue;
    if (zmiana === ChangeKind.Deleted || okno.status === WindowStatus.Closed) {
      zdjete.push(karta.id);
      continue;
    }
    if (okno.title !== undefined && okno.title !== '') karta.nazwa = okno.title;
  }
  for (const idKarty of zdjete) {
    KARTY.delete(idKarty);
    for (const robocze of OKNA) {
      if (!robocze.karty.includes(idKarty)) continue;
      robocze.karty = robocze.karty.filter((id) => id !== idKarty);
      if (robocze.kartaBiezaca === idKarty) {
        robocze.kartaBiezaca = '';
        robocze.nazwa = NAZWA_POCZATKOWA;
      }
    }
  }
  return zdjete;
}

export function zwiazZdarzeniaOkien(naZmiane: (zdjete: string[]) => void): void {
  zglosUchwyt(EventType.WindowChanged, (tresc) => {
    naZmiane(zastosujZmianeOkna(tresc.change, tresc.window));
  });
}

export function kartyOkna(okno: OknoRobocze = oknoBiezace()): KartaRobocza[] {
  const karty: KartaRobocza[] = [];
  for (const idKarty of okno.karty) {
    const karta = KARTY.get(idKarty);
    if (karta !== undefined) karty.push(karta);
  }
  return karty;
}

export function zdejmijKarteOkna(idKarty: string): string[] {
  const okno = oknoBiezace();
  if (!okno.karty.includes(idKarty)) return okno.karty;
  okno.karty = okno.karty.filter((id) => id !== idKarty);
  zwolnijKarte(idKarty);
  if (okno.kartaBiezaca === idKarty) {
    okno.kartaBiezaca = '';
    okno.nazwa = NAZWA_POCZATKOWA;
  }
  return okno.karty;
}

export function zostawKarte(idKarty: string): string[] {
  const okno = oknoBiezace();
  const zdjete = okno.karty.filter((id) => id !== idKarty);
  okno.karty = okno.karty.filter((id) => id === idKarty);
  for (const id of zdjete) zwolnijKarte(id);
  return zdjete;
}

export function zdejmijKartyPoPrawej(idKarty: string): string[] {
  const okno = oknoBiezace();
  const numer = okno.karty.indexOf(idKarty);
  if (numer < 0) return [];
  const zdjete = okno.karty.slice(numer + 1);
  okno.karty = okno.karty.slice(0, numer + 1);
  for (const id of zdjete) zwolnijKarte(id);
  return zdjete;
}

export async function odczytajUkladPaneli(
  kanal: Kanal,
  idKarty: string,
  idPanelu: string,
): Promise<PanelSection[] | null> {
  const okno = KARTY.get(idKarty)?.idOknaKomunikacji ?? '';
  if (okno === '') return null;
  const wynik = await wywolaj(kanal, Command.PanelSectionsGet, { windowId: okno, panelId: idPanelu });
  if (!wynik.udany || wynik.wynik === undefined) return null;
  return wynik.wynik.sections;
}

export async function zapiszUkladPaneli(
  kanal: Kanal,
  idKarty: string,
  idPanelu: string,
  sekcje: PanelSection[],
): Promise<boolean> {
  const okno = KARTY.get(idKarty)?.idOknaKomunikacji ?? '';
  if (okno === '') return false;
  const wynik = await wywolaj(kanal, Command.PanelSectionsSet, {
    windowId: okno, panelId: idPanelu, sections: sekcje,
  });
  return wynik.udany;
}

export async function odtworzOknaRobocze(
  kanal: Kanal,
  sesje: Session[],
  obecnosc: SessionPresence[],
  ogniskowana: string,
): Promise<boolean> {
  // Wykaz niosący pracę nie jest nadpisywany: odtworzenie należy do startu.
  if (KARTY.size > 0) return true;
  const okna = await wywolaj(kanal, Command.WindowList, { status: WindowStatus.Open });
  const moduly = await wywolaj(kanal, Command.ModuleList, {});
  if (!okna.udany || okna.wynik === undefined) return false;
  if (!moduly.udany || moduly.wynik === undefined) return false;
  const kodyModulow = new Map(moduly.wynik.modules.map((modul) => [modul.id, modul.code]));
  const nazwyModulow = new Map(moduly.wynik.modules.map((modul) => [modul.id, modul.name]));
  const oknaSesji = new Map<string, Window[]>();
  for (const okno of okna.wynik.windows) {
    const wykaz = oknaSesji.get(okno.sessionId) ?? [];
    wykaz.push(okno);
    oknaSesji.set(okno.sessionId, wykaz);
  }
  const przed = OKNA.length;
  const przedBiezace = biezace;
  let odtworzone = '';
  for (const sesja of sesje) {
    const wykaz = oknaSesji.get(sesja.id) ?? [];
    if (wykaz.length === 0) continue;
    const okno = otworzOkno();
    okno.idSesji = sesja.id;
    okno.nazwa = sesja.title ?? NAZWA_POCZATKOWA;
    const ogniskowaneOkno = obecnosc.find((odpis) => odpis.sessionId === sesja.id)?.focusedWindowId;
    for (const okienko of wykaz) {
      const kod = kodyModulow.get(okienko.moduleId);
      // Okno modułu spoza rejestru nie ma wnętrza, które karta mogłaby postawić.
      if (kod === undefined) continue;
      const karta = zalozKarte(okno, kod, okienko.id,
        okienko.title ?? nazwyModulow.get(okienko.moduleId) ?? kod);
      if (okienko.id === ogniskowaneOkno) {
        okno.kartaBiezaca = karta.id;
        okno.nazwa = karta.nazwa;
      }
    }
    // Sesja o samych oknach spoza rejestru nie ma czego postawić w karcie.
    if (okno.karty.length === 0) OKNA.pop();
    else if (sesja.id === ogniskowana) odtworzone = okno.id;
  }
  if (OKNA.length === przed) {
    biezace = przedBiezace;
    return true;
  }
  // Okna otwarte, zanim rdzeń odpowiedział, są puste i schodzą z przełącznika.
  OKNA.splice(0, przed);
  biezace = odtworzone !== '' ? odtworzone : OKNA[0].id;
  return true;
}

function zalozKarte(
  okno: OknoRobocze,
  kodModulu: string,
  idOknaKomunikacji: string,
  nazwa: string,
): KartaRobocza {
  licznikKart += 1;
  const karta: KartaRobocza = {
    id: 'krt-' + String(licznikKart), kodModulu, idOknaKomunikacji, nazwa,
  };
  KARTY.set(karta.id, karta);
  okno.karty.push(karta.id);
  return karta;
}

function zwolnijKarte(idKarty: string): void {
  const karta = KARTY.get(idKarty);
  KARTY.delete(idKarty);
  if (karta === undefined || karta.idOknaKomunikacji === '') return;
  if (kanalRdzenia === null) {
    oglos('Karta okna', 'Karta zeszła z pasma, a okno rozmowy zostało w rdzeniu: '
      + 'kanał do rdzenia nie stoi.', 'ostrzezenie');
    return;
  }
  void wywolaj(kanalRdzenia, Command.WindowClose, { windowId: karta.idOknaKomunikacji })
    .then((wynik) => {
      if (wynik.udany) return;
      oglos('Karta okna', wynik.blad?.message ?? 'Rdzeń odmówił zamknięcia okna rozmowy.', 'blad');
    });
}
