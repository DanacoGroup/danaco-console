import {
  AssistantActionControl,
  AssistantActionStatus,
  type AssistantAction,
} from '../../../../shared/contract';
import { NAZWY_STANOW } from './etykiety-assistant';

/**
 * Zdanie o skutku sterowania zleceniem, oparte na odpowiedzi rdzenia, nie na
 * zamówieniu wysłanym z okna.
 *
 * Plik odpowiada wyłącznie za przekład tego, co wróciło
 * z `assistant.action.status`, na zdanie dla czytającego. Stoi poza oknem
 * monitora, bo okno składa tabelę, a to jest ocena odpowiedzi rdzenia.
 *
 * Reguła: porównaj zamówienie ze zleceniem, które wróciło. Odpowiedź rdzenia
 * niesie zlecenie po zmianie, a przekład sterowania na stan stoi w rdzeniu
 * w jednej funkcji (`adapter_modul_asystent_czynnosci.go`, `stanDlaSterowania`),
 * więc sprawdzenie jest darmowe i obowiązuje tak samo dla priorytetu, jak dla
 * czterech przycisków panelu akcji. Dwa powody, dla których zamówienie
 * i odpowiedź potrafią się rozejść:
 *
 *   · `assistant.action.status` zapisuje priorytet wyłącznie przy `control`
 *     innym niż `none`, a sama zmiana kolejności wysyła `none` — żądanie idzie
 *     wtedy torem odczytu i wraca ze zleceniem bez zmian;
 *   · wykonawca zlecenia domyka je własną gorutyną
 *     (`adapter_modul_asystent_wykonawca.go`, `domknijZlecenie`), a rdzeń oddaje
 *     wiersz odczytany po zapisie, więc zapis wykonawcy potrafi wejść między
 *     zapis sterowania a jego odczyt.
 *
 * Zdanie odmowy powstaje wyłącznie z rozbieżności między zamówieniem
 * a odpowiedzią rdzenia — okno nie orzeka o braku, którego rdzeń nie pokazał,
 * tak samo jak nie potwierdza skutku, którego rdzeń nie oddał.
 */
export interface SkutekSterowania {
  zdanie: string;
  udany: boolean;
}

/**
 * Stan, który rdzeń nadaje zleceniu przy danym sterowaniu — jeden do jednego
 * z `stanDlaSterowania` rdzenia. `none` nie zamawia żadnego stanu: to sam
 * odczyt, więc nie ma czego porównywać ze stanem, który wrócił.
 *
 * Słownik jest kluczowany stałymi kontraktu, więc dopisanie sterowania przerwie
 * kompilację tutaj, zamiast po cichu wpaść w gałąź „nic nie zamówiono".
 */
const STAN_ZAMOWIONY: Readonly<Record<AssistantActionControl, AssistantActionStatus | null>> = {
  [AssistantActionControl.None]: null,
  [AssistantActionControl.Pause]: AssistantActionStatus.Paused,
  [AssistantActionControl.Resume]: AssistantActionStatus.Running,
  [AssistantActionControl.Cancel]: AssistantActionStatus.Cancelled,
  [AssistantActionControl.Retry]: AssistantActionStatus.Queued,
  // Zatwierdzenie bramy potwierdzeń puszcza zlecenie w bieg: czekało wyłącznie
  // na rękę Operatora, a nie na zasób. Odmowa zdejmuje je z kolejki — to jest
  // anulowanie z powodem, nie osobny stan.
  [AssistantActionControl.Confirm]: AssistantActionStatus.Running,
  [AssistantActionControl.Reject]: AssistantActionStatus.Cancelled,
};

export function opisSkutku(
  zlecenia: readonly AssistantAction[],
  idZlecenia: string,
  sterowanie: AssistantActionControl,
  priorytet?: number,
): SkutekSterowania {
  const po = zlecenia.find((wpis) => wpis.id === idZlecenia);
  if (po === undefined) {
    return {
      zdanie: 'Rdzeń przyjął żądanie, ale nie oddał tego zlecenia — skutku nie ma czym potwierdzić.',
      udany: false,
    };
  }

  // Każde zamówione pole daje jedno zdanie; brak zamówienia nie daje żadnego,
  // bo o polu, którego żądanie nie ruszyło, okno nie ma nic do powiedzenia.
  const czesci: SkutekSterowania[] = [];
  const zamowionyStan = STAN_ZAMOWIONY[sterowanie];
  if (zamowionyStan !== null) czesci.push(skutekStanu(po, zamowionyStan));
  if (priorytet !== undefined) czesci.push(skutekPriorytetu(po, priorytet));

  if (czesci.length === 0) {
    return {
      zdanie:
        `Rdzeń oddał zlecenie w stanie: ${NAZWY_STANOW[po.status]}. ` +
        'Sterowania ani priorytetu nie zamówiono, więc nie ma czego potwierdzać.',
      udany: true,
    };
  }
  return {
    zdanie: czesci.map((czesc) => czesc.zdanie).join(' '),
    udany: czesci.every((czesc) => czesc.udany),
  };
}

/** Stan sprawdzony w zleceniu, które wróciło — a nie w sterowaniu, które wysłano. */
function skutekStanu(po: AssistantAction, zamowiony: AssistantActionStatus): SkutekSterowania {
  if (po.status !== zamowiony) {
    return {
      zdanie:
        `Rdzeń NIE przestawił zlecenia na stan: ${NAZWY_STANOW[zamowiony]} — ` +
        `oddał je w stanie: ${NAZWY_STANOW[po.status]}. ` +
        'Zlecenie zmienił w międzyczasie ktoś inny niż to sterowanie ' +
        '(wykonawca rdzenia domyka zlecenie własnym zapisem).',
      udany: false,
    };
  }
  return { zdanie: `Rdzeń przestawił zlecenie na stan: ${NAZWY_STANOW[zamowiony]}.`, udany: true };
}

/** Priorytet sprawdzony w zleceniu, które wróciło — a nie w tym, co wysłano. */
function skutekPriorytetu(po: AssistantAction, priorytet: number): SkutekSterowania {
  const zapisany = po.priority ?? 0;
  if (zapisany !== priorytet) {
    return {
      zdanie:
        `Rdzeń NIE zapisał priorytetu ${String(priorytet)} — zlecenie ma nadal ${String(zapisany)}. ` +
        'Komenda assistant.action.status zapisuje pole priority tylko razem ze sterowaniem ' +
        'innym niż „none", a sama zmiana kolejności stanu nie rusza. Pole wróciło do wartości rdzenia.',
      udany: false,
    };
  }
  return { zdanie: `Rdzeń zapisał priorytet ${String(zapisany)}.`, udany: true };
}
