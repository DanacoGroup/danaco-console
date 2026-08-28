import {
  Command,
  type ModuleListResponse,
  type RequestOf,
  type ResponseOf,
  type SessionListResponse,
  type WindowListResponse,
} from '../../../shared/contract';
import type { Kanal, Wynik } from '../protokol/kanal';
import { utworzZrodloAod } from './zrodlo-komend';

/** Byt, którego dotyczy wyciszenie kontekstowe: identyfikator wraz z pełną nazwą czytelną dla Operatora. */
export interface BytKontekstu {
  id: string;
  /** Pełna nazwa widziana przez Operatora; nigdy kod wymyślony przez nakładkę. */
  nazwa: string;
}

/** Kontekst nakładki odczytany z rdzenia: bieżący moduł, karta sesji oraz moduły okien komunikacji rdzenia. */
export interface KontekstWyciszenia {
  /** Ponawia odczyt kontekstu; odmowa nie jest awarią, zostaje kontekst poprzedni. */
  odswiez(): Promise<void>;
  /** Bieżący moduł albo `null`, gdy rdzeń go nie wskazał. */
  biezacyModul(): BytKontekstu | null;
  /** Bieżąca karta sesji albo `null`. */
  biezacaKarta(): BytKontekstu | null;
  /** Moduł okna sugestii; `undefined`, gdy nakładka go nie rozpoznaje. */
  modulOkna(idOkna?: string): string | undefined;
  /** Pełna nazwa modułu; identyfikator, gdy nazwy nie znamy. */
  nazwaModulu(idModulu: string): string;
  // Zdanie mówiące, dlaczego bieżący byt nie jest znany — menu nie wyszarza pozycji.
  powodBrakuModulu(): string;
  powodBrakuKarty(): string;
}

/** Opakowuje `kanal.wyslij` w Promise, tak samo jak pozostałe źródła katalogu kontekstu wyciszenia nakładki. */
function poslijKomende<K extends Command>(
  kanal: Kanal,
  komenda: K,
  zadanie: RequestOf<K>,
): Promise<Wynik<ResponseOf<K>>> {
  return new Promise((rozstrzygnij) => {
    kanal.wyslij(komenda, zadanie, (wynik) => rozstrzygnij(wynik));
  });
}

export function utworzKontekstWyciszenia(kanal: Kanal): KontekstWyciszenia {
  const zrodlo = utworzZrodloAod(kanal);

  /** Moduł każdego okna komunikacji — `Window.moduleId`. */
  const modulOknaWykaz = new Map<string, string>();
  /** Sesja każdego okna komunikacji — droga zapasowa, gdy nakładka nie zna sesji. */
  const sesjaOkna = new Map<string, string>();
  /** Pełna nazwa modułu — `Module.name`. */
  const nazwyModulow = new Map<string, string>();
  /** Nazwa karty sesji — `Session.title`. */
  const nazwySesji = new Map<string, string>();

  let idOknaNakladki: string | undefined;
  let idSesjiNakladki: string | undefined;
  /** Powód, dla którego ostatni odczyt niczego nie przyniósł; pusty, gdy przyniósł. */
  let powodOdmowy = '';

  function okna(): Promise<Wynik<WindowListResponse>> {
    return poslijKomende(kanal, Command.WindowList, {});
  }

  function moduly(): Promise<Wynik<ModuleListResponse>> {
    return poslijKomende(kanal, Command.ModuleList, {});
  }

  function sesje(): Promise<Wynik<SessionListResponse>> {
    return poslijKomende(kanal, Command.SessionList, {});
  }

  async function odswiez(): Promise<void> {
    const [wynikStanu, wynikOkien, wynikModulow, wynikSesji] = await Promise.all([
      zrodlo.stan(),
      okna(),
      moduly(),
      sesje(),
    ]);

    if (wynikStanu.udany && wynikStanu.wynik !== undefined) {
      idOknaNakladki = wynikStanu.wynik.status.activeWindowId;
      idSesjiNakladki = wynikStanu.wynik.status.activeSessionId;
    }

    if (wynikOkien.udany && wynikOkien.wynik !== undefined) {
      modulOknaWykaz.clear();
      sesjaOkna.clear();
      for (const okno of wynikOkien.wynik.windows) {
        if (okno.moduleId !== '') modulOknaWykaz.set(okno.id, okno.moduleId);
        if (okno.sessionId !== '') sesjaOkna.set(okno.id, okno.sessionId);
      }
    }

    if (wynikModulow.udany && wynikModulow.wynik !== undefined) {
      nazwyModulow.clear();
      for (const modul of wynikModulow.wynik.modules) nazwyModulow.set(modul.id, modul.name);
    }

    if (wynikSesji.udany && wynikSesji.wynik !== undefined) {
      nazwySesji.clear();
      for (const sesja of wynikSesji.wynik.sessions) {
        if (sesja.title !== undefined && sesja.title !== '') nazwySesji.set(sesja.id, sesja.title);
      }
    }

    powodOdmowy = wynikStanu.udany
      ? ''
      : `Rdzeń odmówił odczytu stanu nakładki (${wynikStanu.blad?.code ?? 'bez kodu'}).`;
  }

  function nazwaModulu(idModulu: string): string {
    return nazwyModulow.get(idModulu) ?? idModulu;
  }

  function biezacyModul(): BytKontekstu | null {
    if (idOknaNakladki === undefined) return null;
    const idModulu = modulOknaWykaz.get(idOknaNakladki);
    if (idModulu === undefined) return null;
    return { id: idModulu, nazwa: nazwaModulu(idModulu) };
  }

  function biezacaKarta(): BytKontekstu | null {
    const id =
      idSesjiNakladki ??
      (idOknaNakladki === undefined ? undefined : sesjaOkna.get(idOknaNakladki)) ??
      kanal.sesja().id() ??
      undefined;
    if (id === undefined || id === '') return null;
    return { id, nazwa: nazwySesji.get(id) ?? id };
  }

  return {
    odswiez,
    biezacyModul,
    biezacaKarta,

    modulOkna: (idOkna) => (idOkna === undefined ? undefined : modulOknaWykaz.get(idOkna)),

    nazwaModulu,

    powodBrakuModulu: () =>
      powodOdmowy === ''
        ? `Nakładka nie wie, w którym module Operator pracuje: \`${Command.AodStatusGet}\` nie wskazał okna ` +
          `albo \`${Command.WindowList}\` nie zna jego modułu. Kontrakt nie niesie modułu przy stanie ` +
          'nakładki — brak jest po stronie kontraktu. Wyciszenie karty sesji działa niezależnie.'
        : `${powodOdmowy} Otwórz powierzchnię interakcji i odczytaj stan ponownie; wyciszenie ` +
          'czasowe i wyciszenie klasy zdarzeń działają bez tego odczytu.',

    powodBrakuKarty: () =>
      powodOdmowy === ''
        ? 'Nakładka nie ma dziś karty sesji: rdzeń nie wskazał sesji w nakładce, a połączenie nie ' +
          'niesie sesji. Otwórz kartę sesji, a pozycja nazwie ją wprost.'
        : `${powodOdmowy} Otwórz powierzchnię interakcji i odczytaj stan ponownie.`,
  };
}
