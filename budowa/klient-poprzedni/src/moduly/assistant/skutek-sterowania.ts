import {
  AssistantActionControl,
  AssistantActionStatus,
  type AssistantAction,
} from '../../../../shared/contract';
import { NAZWY_STANOW } from './etykiety-assistant';

/**
 * Zdanie o skutku sterowania zleceniem, oparte na zleceniu zwróconym przez
 * komendę `assistant.action.status`, a nie na zamówieniu wysłanym z okna.
 * Skutek powstaje z porównania zamówionego stanu i priorytetu z tym, co wróciło.
 */
export interface SkutekSterowania {
  zdanie: string;
  udany: boolean;
}

/**
 * Stan, który rdzeń nadaje zleceniu przy danym sterowaniu, jeden do jednego
 * z przekładem rdzenia. Sterowanie `none` nie zamawia stanu, bo jest samym
 * odczytem. Słownik kluczują stałe kontraktu `AssistantActionControl`.
 */
const STAN_ZAMOWIONY: Readonly<Record<AssistantActionControl, AssistantActionStatus | null>> = {
  [AssistantActionControl.None]: null,
  [AssistantActionControl.Pause]: AssistantActionStatus.Paused,
  [AssistantActionControl.Resume]: AssistantActionStatus.Running,
  [AssistantActionControl.Cancel]: AssistantActionStatus.Cancelled,
  [AssistantActionControl.Retry]: AssistantActionStatus.Queued,
  // Zatwierdzenie bramy puszcza zlecenie w bieg, odmowa zdejmuje je z kolejki.
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

  // Kazde zamowione pole daje jedno zdanie; pole niezamowione nie daje zadnego.
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

/**
 * Sprawdza stan w zleceniu, które wróciło z rdzenia, i zestawia go ze stanem
 * zamówionym; rozbieżność daje zdanie odmowy z nazwami obu stanów.
 */
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

/**
 * Sprawdza pole `priority` w zleceniu, które wróciło z rdzenia, i zestawia je
 * z priorytetem zamówionym; rozbieżność daje zdanie odmowy z obiema wartościami.
 */
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
