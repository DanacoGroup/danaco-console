/**
 * Okna robocze aplikacji wraz z kartami, które w nich stoją. Okno robocze jest
 * kontenerem jednego działania Operatora: otwiera się na Centrum dowodzenia,
 * niesie najwyżej jedną sesję rdzenia i trzyma własne karty. Karta jest jedną
 * pracą, nie modułem — niesie kod modułu i okno komunikacji rdzenia osobno,
 * więc w oknie stoi obok siebie wiele kart tego samego modułu. Kontrakt nie ma
 * pola na wykaz okien roboczych: odtwarza się przy starcie z `home.enter`
 * i `window.list`.
 */
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

/** Karta okna roboczego: jedna praca Operatora wraz z jej oknem w rdzeniu. */
export interface KartaRobocza {
  /** Oznaczenie karty nadane przy otwarciu; nie wychodzi poza klienta. */
  id: string;
  /** Kod modułu, którego wnętrze karta pokazuje. */
  kodModulu: string;
  /** Okno komunikacji w rdzeniu; pustka znaczy okno jeszcze niezałożone. */
  idOknaKomunikacji: string;
  /** Nazwa karty widoczna w paśmie. */
  nazwa: string;
}

/** Okno robocze widziane przez wiązanie: nazwa dla przełącznika, sesja i karta bieżąca. */
export interface OknoRobocze {
  /** Oznaczenie okna nadane przy otwarciu; nie wychodzi poza klienta. */
  id: string;
  /** Nazwa pokazywana w przełączniku; okno bez pracy nazywa się Centrum dowodzenia. */
  nazwa: string;
  /** Sesja rdzenia niesiona przez okno; pustka znaczy sesję jeszcze niezałożoną — powstaje pierwszą wiadomością. */
  idSesji: string;
  /** Identyfikator karty bieżącej; pustka znaczy kartę główną — Centrum. */
  kartaBiezaca: string;
  /** Identyfikatory kart otwartych w oknie, w kolejności otwarcia; karta główna stoi poza wykazem. */
  karty: string[];
}

const OKNA: OknoRobocze[] = [];
/* Karty wszystkich okien roboczych. Okno wymienia je identyfikatorem, bo karta
   przeżywa zamknięcie pasma i odszukuje się ją po identyfikatorze, nie po
   miejscu w wykazie. */
const KARTY = new Map<string, KartaRobocza>();
let biezace = '';
let licznik = 0;
let licznikKart = 0;
/* Kanał wpięty przy wiązaniu Centrum. Zamknięcie karty i okna sięga rdzenia,
   a woła je czynność okna, która kanału nie niesie. */
let kanalRdzenia: Kanal | null = null;

/** Nazwa okna roboczego bez podjętej pracy. */
const NAZWA_POCZATKOWA = 'Centrum dowodzenia';

/** Wpina kanał rdzenia, którym okna robocze zamykają okna komunikacji swoich kart. */
export function wskazKanalRdzenia(kanal: Kanal): void {
  kanalRdzenia = kanal;
}

/** Wykaz okien roboczych w kolejności otwarcia. */
export function oknaRobocze(): OknoRobocze[] {
  return OKNA;
}

/** Okno robocze, w którym stoi wskazana karta; brak znaczy kartę zdjętą. */
export function oknoKarty(idKarty: string): OknoRobocze | undefined {
  return OKNA.find((okno) => okno.karty.includes(idKarty));
}

/** Okno robocze niosące wskazaną sesję; brak znaczy sesję bez okna w tym kliencie. */
export function oknoSesji(idSesji: string): OknoRobocze | undefined {
  if (idSesji === '') return undefined;
  return OKNA.find((okno) => okno.idSesji === idSesji);
}

/**
 * Zapisuje na oknie roboczym sesję, którą rdzeń dla niego założył albo
 * otworzył. Okno niesie najwyżej jedną sesję, więc wpis stoi raz — kolejny
 * na oknie z sesją wraca fałszem i nic nie zmienia.
 */
export function przypiszSesjeOkna(idOkna: string, idSesji: string, nazwa = ''): boolean {
  const okno = OKNA.find((kandydat) => kandydat.id === idOkna);
  if (okno === undefined || idSesji === '' || (okno.idSesji !== '' && okno.idSesji !== idSesji)) {
    return false;
  }
  okno.idSesji = idSesji;
  if (nazwa !== '' && okno.kartaBiezaca === '') okno.nazwa = nazwa;
  return true;
}

/** Zdejmuje sesję usuniętą w rdzeniu z okien, które ją niosły; zwraca identyfikatory tych okien. */
export function zdejmijSesjeOkien(idSesji: string): string[] {
  const dotkniete: string[] = [];
  for (const okno of OKNA) {
    if (okno.idSesji !== idSesji || idSesji === '') continue;
    okno.idSesji = '';
    dotkniete.push(okno.id);
  }
  return dotkniete;
}

/** Okno robocze bieżące; przy pierwszym pytaniu otwiera pierwsze okno. */
export function oknoBiezace(): OknoRobocze {
  if (OKNA.length === 0) otworzOkno();
  return OKNA.find((okno) => okno.id === biezace) ?? OKNA[0];
}

/** Otwiera okno robocze na Centrum dowodzenia i czyni je bieżącym. */
export function otworzOkno(): OknoRobocze {
  licznik += 1;
  const okno: OknoRobocze = {
    id: 'okr-' + String(licznik), nazwa: NAZWA_POCZATKOWA, idSesji: '', kartaBiezaca: '', karty: [],
  };
  OKNA.push(okno);
  biezace = okno.id;
  return okno;
}

/** Czyni wskazane okno bieżącym; okno nieznane zostawia stan bez zmiany. */
export function przelaczOkno(id: string): OknoRobocze | undefined {
  const okno = OKNA.find((kandydat) => kandydat.id === id);
  if (okno === undefined) return undefined;
  biezace = okno.id;
  return okno;
}

/**
 * Zamyka okno bieżące wraz z oknami komunikacji jego kart i wraca do
 * poprzedniego. Ostatnie okno nie znika — aplikacja bez okna roboczego nie
 * miałaby czego pokazać; wraca na Centrum.
 */
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

/**
 * Zapisuje na oknie bieżącym kartę, na której Operator stanął, i zwraca jej
 * identyfikator. Kartę rozpoznaje para modułu i okna komunikacji: wejście
 * w moduł bez wskazanego okna wraca na kartę stojącą, a okno wskazane wprost
 * otwiera kartę własną — dwie rozmowy tego samego modułu stoją obok siebie.
 * Kod pusty wraca na kartę główną i zwraca pustkę.
 */
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

/**
 * Otwiera w oknie bieżącym nową kartę modułu i czyni ją bieżącą — zawsze nową,
 * także gdy karta tego modułu już stoi: znak „+" pasma zakłada drugą pracę
 * obok pierwszej. Okno komunikacji podaje się przy wznowieniu sesji; pustka
 * znaczy okno zakładane pierwszą wiadomością. Zwraca identyfikator karty.
 */
export function otworzKarte(kodModulu: string, nazwa: string, idOknaKomunikacji = ''): string {
  const okno = oknoBiezace();
  const karta = zalozKarte(okno, kodModulu, idOknaKomunikacji, nazwa);
  okno.kartaBiezaca = karta.id;
  okno.nazwa = nazwa;
  return karta.id;
}

/**
 * Wiąże wskazaną kartę z oknem komunikacji, które stanęło w rdzeniu. Fałsz
 * znaczy kartę już zdjętą albo wskazanie puste: po tym przypisaniu idzie
 * zamknięcie okna w rdzeniu, a wiązania raz zawiązanego nie zdejmuje wskazanie
 * puste — zdejmuje je dopiero zamknięcie karty.
 */
export function przypiszOknoKomunikacji(idKarty: string, idOknaKomunikacji: string): boolean {
  const karta = KARTY.get(idKarty);
  if (karta === undefined || idOknaKomunikacji === '') return false;
  karta.idOknaKomunikacji = idOknaKomunikacji;
  return true;
}

/** Opis karty po identyfikatorze; brak znaczy kartę już zdjętą. */
export function kartaOkna(idKarty: string): KartaRobocza | undefined {
  return KARTY.get(idKarty);
}

/**
 * Czyni wskazaną kartę bieżącą w oknie bieżącym i zwraca jej opis. Brak znaczy
 * kartę zdjętą albo stojącą w innym oknie roboczym — takiej nie wolno postawić
 * w tym oknie, bo pasmo pokazuje karty jednego okna.
 */
export function wskazKarte(idKarty: string): KartaRobocza | undefined {
  const okno = oknoBiezace();
  const karta = KARTY.get(idKarty);
  if (karta === undefined || !okno.karty.includes(idKarty)) return undefined;
  okno.kartaBiezaca = karta.id;
  okno.nazwa = karta.nazwa;
  return karta;
}

/**
 * Nanosi zmianę okna komunikacji zgłoszoną przez rdzeń na karty okien roboczych.
 * Okno zamknięte poza tym klientem zdejmuje swoją kartę bez wołania `window.close`
 * — okno już nie stoi; zmiana nazwy przechodzi na kartę. Zwraca identyfikatory
 * kart zdjętych.
 */
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

/**
 * Zgłasza uchwyt `window.changed`: zmiana okna komunikacji w rdzeniu wchodzi na
 * karty okien roboczych, a wołający dostaje identyfikatory kart zdjętych, żeby
 * zdjąć je z pasma i z płótna.
 */
export function zwiazZdarzeniaOkien(naZmiane: (zdjete: string[]) => void): void {
  zglosUchwyt(EventType.WindowChanged, (tresc) => {
    naZmiane(zastosujZmianeOkna(tresc.change, tresc.window));
  });
}

/** Karty wskazanego okna roboczego w kolejności otwarcia. */
export function kartyOkna(okno: OknoRobocze = oknoBiezace()): KartaRobocza[] {
  const karty: KartaRobocza[] = [];
  for (const idKarty of okno.karty) {
    const karta = KARTY.get(idKarty);
    if (karta !== undefined) karty.push(karta);
  }
  return karty;
}

/**
 * Zdejmuje kartę z okna bieżącego i zamyka jej okno komunikacji w rdzeniu;
 * zwraca identyfikatory kart, które zostały.
 */
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

/** Zostawia w oknie wyłącznie kartę wskazaną; zwraca identyfikatory kart zdjętych. */
export function zostawKarte(idKarty: string): string[] {
  const okno = oknoBiezace();
  const zdjete = okno.karty.filter((id) => id !== idKarty);
  okno.karty = okno.karty.filter((id) => id === idKarty);
  for (const id of zdjete) zwolnijKarte(id);
  return zdjete;
}

/** Zdejmuje karty stojące po wskazanej; zwraca identyfikatory kart zdjętych. */
export function zdejmijKartyPoPrawej(idKarty: string): string[] {
  const okno = oknoBiezace();
  const numer = okno.karty.indexOf(idKarty);
  if (numer < 0) return [];
  const zdjete = okno.karty.slice(numer + 1);
  okno.karty = okno.karty.slice(0, numer + 1);
  for (const id of zdjete) zwolnijKarte(id);
  return zdjete;
}

/**
 * Układ sekcji panelu karty zapisany w rdzeniu; `null` znaczy odmowę rdzenia
 * albo kartę bez okna komunikacji, dla której układu nie ma gdzie szukać.
 */
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

/**
 * Zapisuje w rdzeniu układ sekcji panelu karty. Fałsz znaczy odmowę rdzenia
 * albo kartę bez okna komunikacji — układ nie ma wtedy do czego przylgnąć.
 */
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

/**
 * Odtwarza wykaz okien roboczych ze stanu rdzenia: sesja czynna konta daje
 * okno robocze niosące tę sesję, a każde jej otwarte okno komunikacji — kartę.
 * Karta ogniskowana ostatnio staje się bieżącą, sesja ogniskowana — oknem
 * bieżącym. Fałsz znaczy odmowę rdzenia: wykaz zostaje pusty i wołający ma to
 * nazwać.
 */
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
  /* Okna otwarte, zanim rdzeń odpowiedział, są puste — praca stoi w tych
     odtworzonych, więc puste schodzą z przełącznika. */
  OKNA.splice(0, przed);
  biezace = odtworzone !== '' ? odtworzone : OKNA[0].id;
  return true;
}

/** Zakłada w oknie kartę o własnym identyfikatorze i wpisuje ją do rejestru kart. */
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

/**
 * Zdejmuje kartę z rejestru i zamyka jej okno komunikacji w rdzeniu. Karta
 * zdjęta bez tego zostawiłaby okno otwarte na zawsze: rdzeń liczy okna sesji
 * i wybiera po nich okno, do którego Operator wraca.
 */
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
