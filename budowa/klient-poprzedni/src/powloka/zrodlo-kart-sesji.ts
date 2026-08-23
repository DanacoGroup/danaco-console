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

/**
 * Źródło pasa kart sesji — jedyna prawda powłoki o tym, jakie sesje trwają.
 *
 * Jedna odpowiedzialność: odpytanie `session.list` i nasłuch zdarzeń
 * `session.changed` oraz `session.focus.changed`, złożone w jedną migawkę pasa.
 *
 * Karta sesji jest sesją rdzenia. Pas nie nadaje kartom identyfikatorów
 * miejscowych i nie zakłada ich sam: karta powstaje, bo rdzeń ma sesję, i znika,
 * bo rdzeń ją zamknął. Sesja i okno komunikacji pozostają dwoma bytami —
 * liczba okien na scenie nie ma tu wpływu.
 *
 * Zakładkami są sesje otwarte: czynne i wstrzymane. Sesja zakończona albo
 * zarchiwizowana przestaje być zakładką, tak samo jak po `session.close`.
 *
 * Ognisko jest właściwością klienta. Kartę czynną wskazuje `session.focus`
 * tego klienta, nie konta — zdarzenia z innego `clientId` pas pomija. Zanim
 * padnie pierwsze zdarzenie, czynna jest sesja powiązana z tym połączeniem.
 *
 * Odmowa albo odpowiedź o złym kształcie daje stan `blad` z treścią odmowy,
 * a nie pusty pas udający brak sesji.
 */

/** Jedna zakładka pasa: sesja rdzenia w postaci gotowej dla karty. */
export interface WpisKarty {
  /** Identyfikator sesji nadany przez rdzeń — on jest identyfikatorem karty. */
  id: string;
  tytul: string;
  stan: ProgressStatus;
}

/** Stan źródła: przed pierwszą odpowiedzią rdzenia, po niej, albo po odmowie. */
export type StanZrodlaKart = 'oczekiwanie' | 'gotowe' | 'blad';

/** Migawka pasa — komplet danych do jednego przerysowania. */
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

/** Nazwa zastępcza sesji bez tytułu; nie udaje danych, nazywa ich brak. */
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

  // Pierwsze odpytanie od razu: ramka przy braku połączenia trafia do kolejki
  // wychodzącej transportu i wychodzi z chwilą nawiązania łączności.
  odpytaj();

  return {
    migawka: zbudujMigawke,
    naZmiane: (sluchacz) => zmiany.subskrybuj(sluchacz),
  };
}

/** Zakładką jest sesja otwarta: czynna albo wstrzymana. */
function czyZakladka(sesja: Session): boolean {
  return sesja.status === SessionStatus.Active || sesja.status === SessionStatus.Paused;
}

/** Przekład sesji rdzenia na zakładkę pasa. */
function wpisKarty(sesja: Session, obecnosc?: SessionPresence): WpisKarty {
  const tytul = sesja.title !== undefined && sesja.title.length > 0 ? sesja.title : BEZ_NAZWY;
  return { id: sesja.id, tytul, stan: stanKarty(sesja, obecnosc) };
}

/**
 * Stan karty liczony wyłącznie z odpowiedzi rdzenia.
 *
 * Kropka „praca w tle" zapala się od strumienia zgłoszonego w żywym stanie
 * sesji, a nie od samego faktu, że sesja jest czynna — inaczej wskaźnik
 * mówiłby o pracy, której nie ma.
 */
function stanKarty(sesja: Session, obecnosc?: SessionPresence): ProgressStatus {
  if ((obecnosc?.streamingWindowCount ?? 0) > 0) return ProgressStatus.Running;
  if (sesja.status === SessionStatus.Paused) return ProgressStatus.Paused;
  return ProgressStatus.Pending;
}
