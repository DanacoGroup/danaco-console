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

/**
 * Wywołanie komendy obszaru `translate.*` — jedyna droga modułu do rdzenia.
 *
 * Rdzeń odpowiada na komendę, której nie obsługuje, kopertą `translate.unknown`
 * (`server/internal/protocol/zadanie.go`). Koperta niesie identyfikator
 * żądania, ale nie niesie pola `status`, więc korelacja klienta nie rozpoznaje
 * jej jako odpowiedzi (`protokol/koperta.ts` — `czyOdpowiedz`). Zwykłe
 * `wywolaj` czekałoby na nią bez końca, a okno stałoby w stanie ładowania na
 * zawsze.
 *
 * Dlatego wywołanie nasłuchuje `translate.unknown` równolegle z odpowiedzią
 * i rozstrzyga się na tym, co przyjdzie pierwsze. Odmowa wraca nazwana: pole
 * `odmowa.zadanyTyp` niesie `requestedType` z rdzenia, więc okno mówi wprost,
 * której komendy rdzeń nie zna, zamiast pokazać pustą listę jako wynik.
 *
 * To nie jest druga droga do rdzenia: wysyłka idzie tym samym `kanal.wyslij`,
 * a nazwa komendy pochodzi wyłącznie ze stałych kontraktu. Obietnica nigdy nie
 * jest odrzucana.
 */

/** Odmowa rdzenia nazwana żądanym typem — treść zdarzenia `translate.unknown`. */
export interface OdmowaRdzenia {
  /** Typ, którego rdzeń nie rozpoznał (`requestedType`). */
  zadanyTyp: string;
  /** Powód podany przez rdzeń; pusty, gdy go nie podał. */
  powod: string;
}

/**
 * Wynik wywołania modułu — `Wynik` protokołu poszerzony o nazwaną odmowę.
 *
 * Pole `blad` jest wypełnione także przy odmowie, więc widok, który zna tylko
 * `Wynik`, zachowuje się poprawnie; `odmowa` pozwala oknu odróżnić „rdzeń nie
 * zna tej komendy" od „rdzeń ją zna i zawiódł".
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

    // Subskrypcja przed wysyłką: odpowiedź rdzenia bywa natychmiastowa,
    // a odmowa przepuszczona byłaby zawieszeniem okna bez wyjścia.
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
 * Czy odmowa dotyczy tego właśnie żądania.
 *
 * Rdzeń podaje `requestId` i po nim rozpoznajemy odmowę w pierwszej
 * kolejności. Gdy pola zabrakło — starszy rdzeń albo pośrednik — zostaje
 * dopasowanie po żądanym typie: mylna zbieżność jest mniej szkodliwa niż okno
 * czekające bez końca.
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
 * Odmowa jako wynik nieudany.
 *
 * Kod `not_found` jest wyborem świadomym: kontrakt nie ma kodu „komenda
 * nieobsługiwana", a spośród ośmiu kodów to ten mówi „wskazanego bytu nie ma".
 * Rzeczą, której nie ma, jest tutaj uchwyt komendy. Zdanie błędu nazywa typ
 * wprost, więc Operator nie musi znać kodu, żeby zrozumieć odmowę.
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
