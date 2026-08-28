import {
  ChangeKind,
  Command,
  EventType,
  SessionStatus,
  type Session,
  type SessionPresence,
} from '../../../shared/contract';
import { utworzMagistrale, type Odsubskrybuj } from '../polaczenie/indeks';
import type { Kanal } from '../protokol/kanal';
import { czyTablica, sprawdzKsztalt } from '../protokol/ksztalt-odpowiedzi';
import { wywolaj } from '../protokol/wywolanie';

/** Źródło danych sekcji sesji w tle jest jedyną prawdą strony głównej o sesjach: odpytuje wykaz sesji z żywym stanem i nasłuchuje zdarzeń rdzenia, składając je w jedną migawkę. */
export interface WpisSesji {
  sesja: Session;
  obecnosc?: SessionPresence;
}

/** Stan źródła: przed pierwszą odpowiedzią rdzenia, po niej z gotowym wykazem, albo po odmowie rdzenia. */
export type StanZrodlaSesji = 'oczekiwanie' | 'gotowe' | 'blad';

/** Migawka sekcji niesie komplet danych potrzebny do jednego pełnego przerysowania widoku sesji w tle strony. */
export interface MigawkaSesji {
  stan: StanZrodlaSesji;
  /** Sesje w tle: czynne i wstrzymane, bez sesji bieżącego połączenia. */
  wpisy: readonly WpisSesji[];
  /** Treść odmowy rdzenia — wyłącznie przy stanie `blad`. */
  blad?: string;
}

export interface ZrodloSesji {
  /** Bieżąca migawka — pierwsze rysowanie przed nadejściem zmiany. */
  migawka(): MigawkaSesji;
  /** Subskrypcja kolejnych migawek. */
  naZmiane(sluchacz: (migawka: MigawkaSesji) => void): Odsubskrybuj;
}

/** Zwłoka scalania zdarzeń okien w jedno ponowne odpytanie rdzenia, liczona w milisekundach realnego czasu. */
const ZWLOKA_ODSWIEZENIA = 300;

export function utworzZrodloSesji(kanal: Kanal): ZrodloSesji {
  const zmiany = utworzMagistrale<MigawkaSesji>();
  const wpisy = new Map<string, WpisSesji>();
  let stan: StanZrodlaSesji = 'oczekiwanie';
  let blad: string | undefined;
  let numerZapytania = 0;
  let zaplanowane: ReturnType<typeof setTimeout> | null = null;

  function zbudujMigawke(): MigawkaSesji {
    const wlasna = kanal.sesja().id();
    const widoczne = [...wpisy.values()]
      .filter((wpis) => wpis.sesja.id !== wlasna && czyWTle(wpis))
      .sort((pierwszy, drugi) => ostatniaCzynnosc(drugi) - ostatniaCzynnosc(pierwszy));
    return { stan, wpisy: widoczne, blad };
  }

  function oglos(): void {
    zmiany.oglos(zbudujMigawke());
  }

  /** Pyta rdzeń o pełny wykaz; odpowiedź przedawniona nie nadpisuje świeższej. */
  function odpytaj(): void {
    const numer = ++numerZapytania;
    void wywolaj(kanal, Command.SessionList, { includePresence: true }).then((wynik) => {
      if (numer !== numerZapytania) return;
      const sprawdzony = sprawdzKsztalt(wynik, Command.SessionList, (tresc) =>
        czyTablica(tresc.sessions),
      );
      if (!sprawdzony.udany) {
        stan = 'blad';
        blad = sprawdzony.blad?.message ?? 'Rdzeń odmówił wykazu sesji.';
        oglos();
        return;
      }
      stan = 'gotowe';
      blad = undefined;
      wpisy.clear();
      const obecnosci = new Map(
        (sprawdzony.wynik?.presence ?? []).map((odpis) => [odpis.sessionId, odpis]),
      );
      for (const sesja of sprawdzony.wynik?.sessions ?? []) {
        wpisy.set(sesja.id, { sesja, obecnosc: obecnosci.get(sesja.id) });
      }
      oglos();
    });
  }

  /** Zdarzenia okien scala w jedno odpytanie — liczby okien niesie żywy stan. */
  function zaplanujOdswiezenie(): void {
    if (zaplanowane !== null) return;
    zaplanowane = setTimeout(() => {
      zaplanowane = null;
      odpytaj();
    }, ZWLOKA_ODSWIEZENIA);
  }

  kanal.naZdarzenie(EventType.SessionChanged, (zdarzenie) => {
    if (zdarzenie.change === ChangeKind.Deleted) {
      wpisy.delete(zdarzenie.session.id);
    } else {
      const dotychczas = wpisy.get(zdarzenie.session.id);
      wpisy.set(zdarzenie.session.id, {
        sesja: zdarzenie.session,
        obecnosc: zdarzenie.presence ?? dotychczas?.obecnosc,
      });
    }
    oglos();
  });

  kanal.naZdarzenie(EventType.WindowChanged, () => zaplanujOdswiezenie());

  // Pierwsze odpytanie od razu: ramka czeka w kolejce transportu do chwili nawiązania łączności.
  odpytaj();

  return {
    migawka: zbudujMigawke,
    naZmiane: (sluchacz) => zmiany.subskrybuj(sluchacz),
  };
}

/** Sesją w tle jest sesja czynna albo wstrzymana; żywy odpis rdzenia przesądza o tym samodzielnie zawsze. */
function czyWTle(wpis: WpisSesji): boolean {
  if (wpis.obecnosc?.live === true) return true;
  return wpis.sesja.status === SessionStatus.Active || wpis.sesja.status === SessionStatus.Paused;
}

/** Porządek wykazu sesji w tle: od sesji o najświeższej czynności do sesji najdawniej używanej w historii. */
function ostatniaCzynnosc(wpis: WpisSesji): number {
  return wpis.obecnosc?.lastActivityAt ?? wpis.sesja.updatedAt;
}
