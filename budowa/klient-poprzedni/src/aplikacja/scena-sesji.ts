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

/** Zależności sceny sesji: droga do rdzenia, opis okna oraz miejsce akcji na pasku górnym samej powłoki. */
export interface ZaleznosciSceny {
  /** Droga do rdzenia: transport, kanał, uzgodnienie. */
  rdzen: PolaczenieZRdzeniem;
  /** Opis okna wspólny dla całej sceny. */
  opis: OpisOkna;
  /** Miejsce akcji paska górnego powłoki. */
  akcjePaska: HTMLElement;
}

/** Scena sesji: okna równoległe, ich rozmowy i ich sterowanie, montowane razem w obszarze roboczym powłoki. */
export interface ScenaSesji {
  /** Rama sceny montowana w obszarze roboczym powłoki. */
  element: HTMLElement;
  /** Układ okien równoległych — liczba, role, więź koordynator–wykonawca. */
  uklad: UkladOkien;
  /** Wprowadza na scenę kolejne okno komunikacji; zwraca jego gniazdo. */
  dodajOkno(): IdGniazda | null;
  /** Liczba okien obecnych na scenie. */
  liczbaOkien(): number;
  // Zdejmuje ze sceny okno ostatnie i zamyka je w rdzeniu komendą `window.close`.
  zamknijOstatnie(): void;
  /** Sesja, do której należą okna sceny; pusta, dopóki scena nie ma okna. */
  sesjaNaScenie(): string;
  // Wprowadza na scenę okno z rdzenia — droga dla okien otwartych innym połączeniem.
  przyjmijOknoZRdzenia(okno: Window): void;
}

/** Scena sesji — miejsce, w którym dzieje się praca: układ okien równoległych, rozmowy i sterowanie każdego okna. */
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

    // Zamówienie przed założeniem sesji nie jest odrzucane — czeka na uzgodnienie i rusza zaraz po nim.
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

  // Gniazdo pierwsze bierze okno z uzgodnienia; wcześniejsze gniazda dostają okno zaraz po sesji.
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
      // Kolejność: najpierw rdzeń, potem widok — odwrotna zostawiałaby okno pracujące po drugiej stronie.
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

    // Wprowadza na scenę okno z rdzenia, zgłoszone przez `window.changed` ze zmianą `created`.
    przyjmijOknoZRdzenia(okno) {
      if (okno.id === '' || okno.sessionId === '') return;
      // Okno własne wchodzi tą drogą, którą zamówiono; drugi raz dołożyłoby gniazdo bez pokrycia.
      for (const wiazanie of wiazania.values()) if (wiazanie.idOkna === okno.id) return;
      // Okna sceny należą do jednej sesji — okno cudzej sesji nie ma tu miejsca ani odmowy.
      const naScenie = sesjaNaScenie();
      if (naScenie !== '' && naScenie !== okno.sessionId) return;

      const wolne = ID_GNIAZD.find((id) => !wiazania.has(id) && !zamawiane.has(id));
      if (wolne === undefined) return;

      // Gniazdo musi być na scenie przed rozmową — podniesienie liczby okien odsłania je jak przełącznik.
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

/** Uczciwa odpowiedź na naciśnięcie, którego scena nie umie dziś wykonać z powodu rozjazdu sesji okien. */
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
