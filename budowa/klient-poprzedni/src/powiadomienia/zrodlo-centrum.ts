import {
  Command,
  EventType,
  type Notification,
  type NotificationClass,
  type NotificationListResponse,
} from '../../../shared/contract';
import type { Odsubskrybuj } from '../polaczenie/magistrala-zdarzen';
import type { Kanal } from '../protokol/kanal';
import { wywolaj } from '../protokol/wywolanie';

// Źródło centrum powiadomień to jedyna droga klienta do rodziny notification, czytana raz i nadążana.

export interface StanCentrum {
  zdarzenia: Notification[];
  /** Liczba zdarzeń nowych w całym rejestrze, nie w widoku. */
  nowe: number;
  /** Zdanie odmowy rdzenia; puste znaczy „odczyt się udał". */
  odmowa: string;
}

export interface FiltrCentrum {
  klasy?: NotificationClass[];
  /** Tylko zdarzenia wymagające decyzji — filtr wagi z karty komponentu. */
  tylkoDecyzje?: boolean;
}

export interface ZrodloCentrum {
  odczytaj(filtr?: FiltrCentrum): Promise<StanCentrum>;
  /** Oznacza wskazane zdarzenia jako odczytane; brak wykazu bierze wszystkie nowe. */
  odczytajZdarzenia(identyfikatory?: string[]): Promise<string>;
  zamknij(identyfikator: string): Promise<string>;
  odloz(identyfikator: string, doKiedy: number): Promise<string>;
  /** Nasłuch obu zdarzeń rodziny; oddaje licznik po każdej zmianie. */
  naZmiane(sluchacz: (nowe: number) => void): Odsubskrybuj;
}

export function utworzZrodloCentrum(kanal: Kanal): ZrodloCentrum {
  async function odczytaj(filtr: FiltrCentrum = {}): Promise<StanCentrum> {
    const odpowiedz = await wywolaj(kanal, Command.NotificationList, {
      ...(filtr.klasy === undefined || filtr.klasy.length === 0 ? {} : { classes: filtr.klasy }),
      ...(filtr.tylkoDecyzje === true ? { weights: ['wymagajaca_decyzji'] } : {}),
    });
    if (!odpowiedz.udany || odpowiedz.wynik === undefined) {
      return { zdarzenia: [], nowe: 0, odmowa: zdanieOdmowy(odpowiedz.blad?.message) };
    }
    const wynik: NotificationListResponse = odpowiedz.wynik;
    return { zdarzenia: wynik.notifications, nowe: wynik.unread, odmowa: '' };
  }

  /** Wspólna postać odpowiedzi czynności: puste zdanie znaczy powodzenie. */
  async function wykonaj(komenda: string, tresc: Record<string, unknown>): Promise<string> {
    const odpowiedz = await wywolaj(kanal, komenda as never, tresc as never);
    return odpowiedz.udany ? '' : zdanieOdmowy(odpowiedz.blad?.message);
  }

  return {
    odczytaj,

    odczytajZdarzenia: (identyfikatory) =>
      wykonaj(
        Command.NotificationAcknowledge,
        identyfikatory === undefined ? {} : { ids: identyfikatory },
      ),

    zamknij: (identyfikator) => wykonaj(Command.NotificationResolve, { id: identyfikator }),

    odloz: (identyfikator, doKiedy) =>
      wykonaj(Command.NotificationSnooze, { id: identyfikator, until: doKiedy }),

    naZmiane(sluchacz) {
      const zdjecia = [
        kanal.naZdarzenie(EventType.NotificationRaised, (tresc) => sluchacz(tresc.unread)),
        kanal.naZdarzenie(EventType.NotificationChanged, (tresc) => sluchacz(tresc.unread)),
      ];
      return () => zdjecia.forEach((zdejmij) => zdejmij());
    },
  };
}

/** Zdanie odmowy rdzenia albo nazwanie milczenia rdzenia po nieudanej czynności; nigdy nie zwraca pustki. */
function zdanieOdmowy(wiadomosc: string | undefined): string {
  return wiadomosc === undefined || wiadomosc === ''
    ? 'Rdzeń odmówił bez podania powodu.'
    : wiadomosc;
}
