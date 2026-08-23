import { Command, type Window } from '../../../shared/contract';
import type { OpisOkna } from '../okno-komunikacji/opis-okna';
import {
  ID_GNIAZD,
  zamontujUkladOkien,
  type IdGniazda,
  type UkladOkien,
} from '../okna-rownolegle/indeks';
import { utworzRejestrKanalow } from '../sterowanie/indeks';
import { pokazKomunikat } from './komunikaty';
import { utworzKorzen } from './korzen-dokumentu';
import { zamowOknoRdzenia } from './otwarcie-okna';
import type { PolaczenieZRdzeniem } from './polaczenie-z-rdzeniem';
import { zwiazGniazdoZOknem, type WiazanieGniazda } from './wiazanie-gniazda';

/** Zależności sceny sesji. */
export interface ZaleznosciSceny {
  /** Droga do rdzenia: transport, kanał, uzgodnienie. */
  rdzen: PolaczenieZRdzeniem;
  /** Opis okna wspólny dla całej sceny. */
  opis: OpisOkna;
  /** Miejsce akcji paska górnego powłoki. */
  akcjePaska: HTMLElement;
}

/** Scena sesji: okna równoległe, ich rozmowy i ich sterowanie. */
export interface ScenaSesji {
  /** Rama sceny montowana w obszarze roboczym powłoki. */
  element: HTMLElement;
  /** Układ okien równoległych — liczba, role, więź koordynator–wykonawca. */
  uklad: UkladOkien;
  /** Wprowadza na scenę kolejne okno komunikacji; zwraca jego gniazdo. */
  dodajOkno(): IdGniazda | null;
  /** Liczba okien obecnych na scenie. */
  liczbaOkien(): number;
  /**
   * Zdejmuje ze sceny okno ostatnie i zamyka je w rdzeniu komendą
   * `window.close`. Samo zmniejszenie układu zostawiłoby okno w `window.list`
   * razem z jego procesem.
   */
  zamknijOstatnie(): void;
  /** Sesja, do której należą okna sceny; pusta, dopóki scena nie ma okna. */
  sesjaNaScenie(): string;
  /**
   * Wprowadza na scenę okno, które POWSTAŁO W RDZENIU poza tą sceną.
   *
   * Droga dla okien otwartych przez asystenta innym połączeniem. Scena nie
   * zamawia wtedy niczego — okno już istnieje — tylko odsłania dla niego
   * gniazdo i wiąże je z rozmową i sterowaniem. Okno cudzej sesji, okno już
   * związane i scena pełna kończą wywołanie bez skutku i bez odmowy: to nie
   * jest komenda Operatora, tylko doniesienie o stanie.
   */
  przyjmijOknoZRdzenia(okno: Window): void;
}

/**
 * Scena sesji — miejsce, w którym dzieje się praca.
 *
 * Jedna odpowiedzialność: związanie trzech warstw w jedną scenę — układu okien
 * równoległych (`okna-rownolegle/`), rozmowy każdego okna (`rozmowa/`) oraz
 * kompletu sterowania każdego okna (`widok-sterowania/`, `sterowanie/`).
 * Scena nie buduje ani okna, ani kontrolki, ani wpisu rozmowy.
 *
 * Okno komunikacji to nie karta sesji. Karta w pasie powłoki jest sesją
 * rdzenia i rządzi się komendami `session.*`; okno sceny jest bytem podrzędnym
 * wobec sesji i rządzi się komendami `window.*`. Liczbę okien ustawia wyłącznie
 * przełącznik „Okna komunikacji: 1 2 3" — pas kart nie dokłada okien i ich nie
 * zdejmuje.
 *
 * Okno powstaje, gdy wchodzi na scenę. Pierwsze okno otwiera uzgodnienie
 * z rdzeniem. Drugie i trzecie zamawiane są komendą `window.create` dokładnie
 * w chwili, gdy Operator wprowadza je na scenę przełącznikiem liczby okien.
 * Wejście gniazda na scenę rozpoznaje obserwator atrybutu `hidden`: o
 * widoczności rozstrzyga układ okien, a scena obserwuje tylko jego skutek.
 *
 * Wszystkie okna sceny należą do jednej sesji. Po powiązaniu połączenia z inną
 * sesją (`session.bind` z Centrum dowodzenia) scena wciąż niesie okna sesji
 * poprzedniej. Kolejne okno powstałoby wtedy w sesji innej niż jego sąsiedzi,
 * więc scena go nie zamawia i mówi o tym wprost.
 *
 * Uzgodnienia scena nie rozpoczyna. Robi to przepływ komunikatów podpięty przez
 * `zamontujUkladOkien` w chwili, gdy transport zgłosi stan „połączony".
 * Wywołanie stąd dałoby drugie powitanie i podwojenie całej historii.
 */
export function utworzSceneSesji(zaleznosci: ZaleznosciSceny): ScenaSesji {
  const { rdzen, opis, akcjePaska } = zaleznosci;

  const korzen = utworzKorzen(akcjePaska);
  const uklad = zamontujUkladOkien(korzen, rdzen, opis);

  const rejestrKanalow = utworzRejestrKanalow(rdzen.kanal);
  rejestrKanalow.odswiez();

  const wiazania = new Map<IdGniazda, WiazanieGniazda>();
  /** Gniazda z zamówieniem w drodze — zapora przed drugim `window.create`. */
  const zamawiane = new Set<IdGniazda>();
  /** Gniazda, które weszły na scenę przed zakończeniem uzgodnienia. */
  const oczekujace = new Set<IdGniazda>();

  /** Czy rdzeń założył już sesję, w której mogą powstawać kolejne okna. */
  function sesjaGotowa(): boolean {
    return rdzen.kanal.sesja().id().length > 0;
  }

  /** Sesja okien już związanych; pusta, dopóki scena nie ma ani jednego okna. */
  function sesjaNaScenie(): string {
    return [...wiazania.values()][0]?.idSesji ?? '';
  }

  /** Doprowadza do gniazda rozmowę i komplet sterowania okna. */
  function zwiaz(id: IdGniazda, okno: Window): void {
    const gniazdo = uklad.gniazdo(id);
    if (gniazdo === null || wiazania.has(id)) return;

    const wiazanie = zwiazGniazdoZOknem({
      kanal: rdzen.kanal,
      rejestrKanalow,
      gniazdo,
      okno,
      persona: opis.kanalModelu,
      panel: korzen.panel,
    });
    wiazanie.ustawWidocznosc(gniazdo.widoczne());
    wiazania.set(id, wiazanie);
  }

  /** Zamawia okno rdzenia dla gniazda, które właśnie weszło na scenę. */
  function zamow(id: IdGniazda): void {
    const gniazdo = uklad.gniazdo(id);
    if (gniazdo === null || wiazania.has(id) || zamawiane.has(id)) return;

    // Zamówienie przed założeniem sesji nie jest odrzucane — czeka na
    // uzgodnienie i rusza zaraz po nim (żadnych blokad w interfejsie).
    if (!sesjaGotowa()) {
      oczekujace.add(id);
      return;
    }

    if (rozjazdSesji(sesjaNaScenie(), rdzen.kanal.sesja().id())) {
      ostrzezORozjezdzieSesji();
      return;
    }

    zamawiane.add(id);
    const tytul = `${opis.tytul} ${ID_GNIAZD.indexOf(id) + 1}`;
    zamowOknoRdzenia(rdzen.kanal, opis, gniazdo.rola(), tytul, ({ okno }) => {
      zamawiane.delete(id);
      if (okno !== null) zwiaz(id, okno);
    });
  }

  // Gniazdo pierwsze bierze okno uzgodnione z rdzeniem; gniazda, które zdążyły
  // wejść na scenę wcześniej, dostają swoje okno zaraz po założeniu sesji.
  rdzen.uzgodnienie.naOtwarcieOkna((okno) => {
    zwiaz(ID_GNIAZD[0], okno);
    for (const id of [...oczekujace]) {
      oczekujace.delete(id);
      zamow(id);
    }
  });

  for (const id of ID_GNIAZD.slice(1)) {
    pilnujGniazda(uklad, id, (widoczne) => {
      wiazania.get(id)?.ustawWidocznosc(widoczne);
      if (widoczne) zamow(id);
    });
  }

  /** Gniazdo ostatnie spośród związanych — to ono schodzi ze sceny. */
  function ostatnieZwiazane(): IdGniazda | null {
    for (let miejsce = ID_GNIAZD.length - 1; miejsce >= 0; miejsce -= 1) {
      const id = ID_GNIAZD[miejsce];
      if (id !== undefined && wiazania.has(id)) return id;
    }
    return null;
  }

  return {
    element: korzen.element,

    zamknijOstatnie() {
      const id = ostatnieZwiazane();
      if (id === null) {
        uklad.ustawLiczbe(Math.max(1, uklad.liczba() - 1));
        return;
      }
      const wiazanie = wiazania.get(id);
      wiazania.delete(id);
      uklad.ustawLiczbe(Math.max(1, uklad.liczba() - 1));
      if (wiazanie === undefined) return;
      // Kolejność: najpierw rdzeń, potem widok. Odwrotna zostawiałaby Operatora
      // z pustym gniazdem i oknem, które nadal pracuje po drugiej stronie.
      rdzen.kanal.wyslij(Command.WindowClose, { windowId: wiazanie.idOkna }, () => {});
      wiazanie.rozlacz();
    },
    uklad,

    dodajOkno() {
      const nastepne = ID_GNIAZD.at(uklad.liczba());
      if (nastepne === undefined) return null;
      uklad.ustawLiczbe(uklad.liczba() + 1);
      return nastepne;
    },

    liczbaOkien: () => uklad.liczba(),
    sesjaNaScenie,

    /**
     * Wprowadza na scenę okno założone w rdzeniu poza nią — tak, jak zgłasza je
     * `window.changed` ze zmianą `created`. Bez tej drogi okno asystenta
     * pracowałoby w rdzeniu, mając proces i kanał modelu, a ekran pokazywałby
     * dalej stan sprzed jego powstania.
     *
     * Scena nie rozstrzyga, kto okno otworzył, i o nic nie pyta: okno sesji
     * Operatora ma być na jego ekranie niezależnie od sprawcy.
     */
    przyjmijOknoZRdzenia(okno) {
      if (okno.id === '' || okno.sessionId === '') return;
      // Okno własne wchodzi tą samą drogą, którą je zamówiono; wpuszczenie go
      // tu drugi raz dołożyłoby scenie gniazdo bez pokrycia.
      for (const wiazanie of wiazania.values()) if (wiazanie.idOkna === okno.id) return;
      // Okna sceny należą do jednej sesji — okno cudzej sesji nie ma tu
      // miejsca, a Operator o nic nie prosił, więc nie ma też odmowy.
      const naScenie = sesjaNaScenie();
      if (naScenie !== '' && naScenie !== okno.sessionId) return;

      const wolne = ID_GNIAZD.find((id) => !wiazania.has(id) && !zamawiane.has(id));
      if (wolne === undefined) return;

      // Gniazdo musi być NA SCENIE, zanim dostanie rozmowę: `zwiaz` bierze
      // gniazdo z układu, a gniazdo poza sceną jest ukryte atrybutem `hidden`.
      // Podniesienie liczby okien odsłania je tą samą drogą, którą odsłania je
      // przełącznik Operatora — scena nie ma drugiego sposobu na pokazanie okna.
      const miejsce = ID_GNIAZD.indexOf(wolne) + 1;
      if (uklad.liczba() < miejsce) uklad.ustawLiczbe(miejsce);
      zwiaz(wolne, okno);
    },
  };
}

/**
 * Czy sesja połączenia rozjechała się z sesją okien sceny. Rozjazd nie jest
 * błędem — bierze się z powrotu do sesji w tle (`session.bind`) — ale kolejne
 * okno powstałoby w sesji innej niż jego sąsiedzi, więc zamówienie nie idzie.
 */
function rozjazdSesji(naScenie: string, wPolaczeniu: string): boolean {
  return naScenie.length > 0 && naScenie !== wPolaczeniu;
}

/** Uczciwa odpowiedź na naciśnięcie, którego scena nie umie dziś wykonać. */
function ostrzezORozjezdzieSesji(): void {
  pokazKomunikat({
    tytul: 'Okno stanęłoby w innej sesji',
    tresc:
      'Okna na scenie należą do sesji sprzed powrotu, a to połączenie jest powiązane z inną. Wejdź w środowisko od nowa, zanim dołożysz kolejne okno.',
    waga: 'ostrz',
  });
}

/**
 * Obserwuje wejście gniazda na scenę i zejście z niej. Nośnikiem zmiany jest
 * atrybut `hidden` gniazda, bo to układ okien rozstrzyga o widoczności — także
 * wtedy, gdy nikt nie naciska przełącznika liczby okien.
 */
function pilnujGniazda(
  uklad: UkladOkien,
  id: IdGniazda,
  przyZmianie: (widoczne: boolean) => void,
): void {
  const gniazdo = uklad.gniazdo(id);
  if (gniazdo === null) return;
  const reaguj = (): void => przyZmianie(gniazdo.widoczne());
  new MutationObserver(reaguj).observe(gniazdo.element, { attributeFilter: ['hidden'] });
}
