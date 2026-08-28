import {
  ChangeKind,
  Command,
  EventType,
  ProgressStatus,
  SessionStatus,
  type Session,
  type SessionPresence,
} from '../../../shared/contract';
import { utworzMagistrale, utworzObserwatorOgniska, type Odsubskrybuj } from '../polaczenie/indeks';
import type { Kanal } from '../protokol/kanal';
import { czyTablica, sprawdzKsztalt } from '../protokol/ksztalt-odpowiedzi';
import type { TozsamoscKlienta } from '../protokol/tozsamosc-klienta';
import { wywolaj } from '../protokol/wywolanie';

// Źródło pasa kart sesji — jedyna prawda powłoki o tym, jakie sesje rdzenia w tej chwili trwają.

/** Jedna zakładka pasa: sesja rdzenia w postaci gotowej dla karty widocznej w tym samym pasie sesji klienta. */
export interface WpisKarty {
  /** Identyfikator sesji nadany przez rdzeń — on jest identyfikatorem karty. */
  id: string;
  tytul: string;
  stan: ProgressStatus;
}

/** Stan źródła: przed pierwszą odpowiedzią rdzenia, zaraz po niej, albo po jego pełnej odmowie wykonania. */
export type StanZrodlaKart = 'oczekiwanie' | 'gotowe' | 'blad';

/** Migawka pasa — komplet danych do jednego przerysowania pasa kart sesji tego samego klienta tej powłoki. */
export interface MigawkaKart {
  stan: StanZrodlaKart;
  wpisy: readonly WpisKarty[];
  /** Sesja ogniskowana przez tego klienta; pusta, gdy rdzeń jej nie wskazał. */
  ogniskowana: string;
  /** Treść odmowy rdzenia — wyłącznie przy stanie `blad`. */
  blad?: string;
}

export interface ZrodloKartSesji {
  /** Bieżąca migawka — pierwsze rysowanie przed nadejściem zmiany. */
  migawka(): MigawkaKart;
  /** Subskrypcja kolejnych migawek. */
  naZmiane(sluchacz: (migawka: MigawkaKart) => void): Odsubskrybuj;
}

/** Nazwa zastępcza sesji bez tytułu; nie udaje danych, których wcale nie ma, tylko wprost nazywa ich brak. */
const BEZ_NAZWY = 'Sesja bez nazwy';

export function utworzZrodloKartSesji(kanal: Kanal, klient: TozsamoscKlienta): ZrodloKartSesji {
  const zmiany = utworzMagistrale<MigawkaKart>();
  const sesje = new Map<string, Session>();
  const obecnosci = new Map<string, SessionPresence>();
  const ognisko = utworzObserwatorOgniska(kanal);
  let stan: StanZrodlaKart = 'oczekiwanie';
  let blad: string | undefined;
  let numerZapytania = 0;

  function zbudujMigawke(): MigawkaKart {
    const wpisy = [...sesje.values()]
      .filter(czyZakladka)
      .sort((pierwsza, druga) => pierwsza.createdAt - druga.createdAt)
      .map((sesja) => wpisKarty(sesja, obecnosci.get(sesja.id)));
    return { stan, wpisy, ogniskowana: ogniskowana(), blad };
  }

  /** Karta czynna: ognisko tego klienta, a przed pierwszym zdarzeniem — sesja połączenia. */
  function ogniskowana(): string {
    const zOgniska = ognisko.ogniskowanaSesja();
    return zOgniska.length > 0 ? zOgniska : kanal.sesja().id();
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
      sesje.clear();
      obecnosci.clear();
      for (const odpis of sprawdzony.wynik?.presence ?? []) obecnosci.set(odpis.sessionId, odpis);
      for (const sesja of sprawdzony.wynik?.sessions ?? []) sesje.set(sesja.id, sesja);
      oglos();
    });
  }

  kanal.naZdarzenie(EventType.SessionChanged, (zdarzenie) => {
    if (zdarzenie.change === ChangeKind.Deleted) {
      sesje.delete(zdarzenie.session.id);
      obecnosci.delete(zdarzenie.session.id);
    } else {
      sesje.set(zdarzenie.session.id, zdarzenie.session);
      if (zdarzenie.presence !== undefined) obecnosci.set(zdarzenie.session.id, zdarzenie.presence);
    }
    oglos();
  });

  // Ognisko innego klienta tego samego konta nie przestawia kart tutaj.
  ognisko.naZmiane((zmiana) => {
    if (zmiana.clientId !== klient.id) return;
    oglos();
  });

  // Pierwsze odpytanie od razu: ramka przy braku połączenia trafia do kolejki wychodzącej transportu.
  odpytaj();

  return {
    migawka: zbudujMigawke,
    naZmiane: (sluchacz) => zmiany.subskrybuj(sluchacz),
  };
}

/** Zakładką jest sesja otwarta w tym rdzeniu: czynna albo wstrzymana, nigdy zakończona ani zarchiwizowana. */
function czyZakladka(sesja: Session): boolean {
  return sesja.status === SessionStatus.Active || sesja.status === SessionStatus.Paused;
}

/** Przekład sesji rdzenia na zakładkę pasa, gotową do narysowania jako karta w tym samym pasie kart sesji. */
function wpisKarty(sesja: Session, obecnosc?: SessionPresence): WpisKarty {
  const tytul = sesja.title !== undefined && sesja.title.length > 0 ? sesja.title : BEZ_NAZWY;
  return { id: sesja.id, tytul, stan: stanKarty(sesja, obecnosc) };
}

/** Stan karty liczony wyłącznie z odpowiedzi rdzenia, nigdy z lokalnego domysłu tego samego widoku klienta. */
function stanKarty(sesja: Session, obecnosc?: SessionPresence): ProgressStatus {
  if ((obecnosc?.streamingWindowCount ?? 0) > 0) return ProgressStatus.Running;
  if (sesja.status === SessionStatus.Paused) return ProgressStatus.Paused;
  return ProgressStatus.Pending;
}
