import {
  ErrorCode,
  EventType,
  type Command,
  type RequestOf,
  type ResponseOf,
  type UnknownCommandPayload,
} from '../../../../shared/contract';
import type { Odsubskrybuj } from '../../polaczenie/magistrala-zdarzen';
import type { Kanal, Wynik } from '../../protokol/kanal';
import { sprawdzKsztalt } from '../../protokol/ksztalt-odpowiedzi';

/** Wywołanie komendy obszaru translate jest jedyną drogą modułu do rdzenia, z odmową nazwaną typem. */

/** Odmowa rdzenia nazwana żądanym typem stanowi treść zdarzenia nieznanej komendy, gdy rdzeń nie rozpoznał żądania. */
export interface OdmowaRdzenia {
  /** Typ, którego rdzeń nie rozpoznał (`requestedType`). */
  zadanyTyp: string;
  /** Powód podany przez rdzeń; pusty, gdy go nie podał. */
  powod: string;
}

/**
 * Wynik wywołania modułu poszerza wynik protokołu o nazwaną odmowę, pozwalając oknu odróżnić
 * nieznaną komendę od zwykłego zawodu rdzenia.
 */
export interface WynikTranslate<T> extends Wynik<T> {
  odmowa?: OdmowaRdzenia;
}

export function zadaj<K extends Command>(
  kanal: Kanal,
  komenda: K,
  zadanie: RequestOf<K>,
  sprawdzian: (tresc: ResponseOf<K>) => boolean,
): Promise<WynikTranslate<ResponseOf<K>>> {
  return new Promise((rozstrzygnij) => {
    let rozstrzygniete = false;
    let idZadania = '';
    let odsubskrybuj: Odsubskrybuj = () => undefined;

    function zakoncz(wynik: WynikTranslate<ResponseOf<K>>): void {
      if (rozstrzygniete) return;
      rozstrzygniete = true;
      odsubskrybuj();
      rozstrzygnij(wynik);
    }

    // Subskrypcja przed wysyłką: odpowiedź bywa natychmiastowa, odmowa przepuszczona zawiesiłaby okno.
    odsubskrybuj = kanal.naZdarzenie(EventType.TranslateUnknown, (tresc) => {
      if (!dotyczyZadania(tresc, idZadania, komenda)) return;
      zakoncz(zbudujOdmowe(komenda, tresc));
    });

    idZadania = kanal.wyslij(komenda, zadanie, (wynik) =>
      zakoncz(sprawdzKsztalt(wynik, komenda, sprawdzian)),
    );
  });
}

/**
 * Dopasowanie odmowy do żądania idzie po identyfikatorze żądania w pierwszej kolejności, a po
 * żądanym typie, gdy pola brak.
 */
function dotyczyZadania(
  tresc: UnknownCommandPayload,
  idZadania: string,
  komenda: string,
): boolean {
  const identyfikator = (tresc.requestId ?? '').trim();
  if (identyfikator !== '') return identyfikator === idZadania;
  return (tresc.requestedType ?? '').trim() === komenda;
}

/**
 * Odmowa jako wynik nieudany bierze kod braku bytu, bo kontrakt nie ma osobnego kodu na komendę
 * nieobsługiwaną przez rdzeń.
 */
function zbudujOdmowe<T>(komenda: string, tresc: UnknownCommandPayload): WynikTranslate<T> {
  const zadanyTyp = ((tresc.requestedType ?? '').trim() || komenda);
  const powod = (tresc.reason ?? '').trim();
  const zdanie = `Rdzeń nie obsługuje komendy ${zadanyTyp} (odmowa translate.unknown)`;
  return {
    udany: false,
    odmowa: { zadanyTyp, powod },
    blad: {
      code: ErrorCode.NotFound,
      message: powod === '' ? zdanie : `${zdanie}: ${powod}`,
      retryable: false,
    },
  };
}
