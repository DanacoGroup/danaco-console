import {
  AssistantActionControl,
  type AssistantAction,
} from '../../../../shared/contract';
import { NAZWY_DROG, NAZWY_STANOW, PLAKIETKI_STANOW } from './etykiety-assistant';

/**
 * Jeden wiersz tabeli zleceń Actions Monitora.
 *
 * Plik odpowiada wyłącznie za zamianę zlecenia w wiersz wraz z jego panelem
 * akcji. Wiersz niczego nie wywołuje — zamiar oddaje oknu, które prowadzi jedno
 * wywołanie `assistant.action.status`.
 *
 * Wszystkie sterowania zostają klikalne niezależnie od stanu zlecenia.
 * Wstrzymania zlecenia już wykonanego widok nie blokuje — odpowiada na nie
 * rdzeń, a odpowiedź trafia do wiersza odpowiedzi okna. Wygaszona kontrolka
 * kazałaby zgadywać, czy przycisk nie działa, czy tylko nie odpowiada.
 */

/** Zamiar wydany z wiersza: sterowanie albo zmiana priorytetu. */
export type NaSterowanie = (
  idZlecenia: string,
  sterowanie: AssistantActionControl,
  priorytet?: number,
) => void;

/** Podgląd szczegółów i wyniku zlecenia — obsługiwany po stronie okna. */
export type NaPodglad = (zlecenie: AssistantAction, wynik: boolean) => void;

/** Sterowania panelu akcji wiersza w kolejności pokazywanej na ekranie. */
const STEROWANIA: readonly (readonly [string, AssistantActionControl])[] = [
  ['Wstrzymaj', AssistantActionControl.Pause],
  ['Wznów', AssistantActionControl.Resume],
  ['Anuluj', AssistantActionControl.Cancel],
  ['Ponów', AssistantActionControl.Retry],
];

export function wierszZlecenia(
  zlecenie: AssistantAction,
  naSterowanie: NaSterowanie,
  naPodglad: NaPodglad,
): HTMLTableRowElement {
  const wiersz = document.createElement('tr');
  wiersz.dataset['zlecenie'] = zlecenie.id;

  wiersz.append(
    komorka(zlecenie.title ?? zlecenie.id),
    komorkaStanu(zlecenie),
    komorka(opisEtapu(zlecenie)),
    komorkaPriorytetu(zlecenie, naSterowanie),
    komorkaAkcji(zlecenie, naSterowanie, naPodglad),
  );
  return wiersz;
}

function komorka(tresc: string): HTMLTableCellElement {
  const element = document.createElement('td');
  element.textContent = tresc;
  return element;
}

/** Stan zlecenia jako plakietka — barwa nigdy nie niesie znaczenia sama. */
function komorkaStanu(zlecenie: AssistantAction): HTMLTableCellElement {
  const plakietka = document.createElement('span');
  plakietka.className = PLAKIETKI_STANOW[zlecenie.status];
  plakietka.dataset['stan'] = zlecenie.status;
  plakietka.textContent = NAZWY_STANOW[zlecenie.status];

  const element = document.createElement('td');
  element.append(plakietka, ` · ${NAZWY_DROG[zlecenie.origin]}`);
  return element;
}

/** Etap bieżący zlecenia wieloetapowego; bez etapów — zdanie o ich braku. */
function opisEtapu(zlecenie: AssistantAction): string {
  if (zlecenie.totalSteps === undefined || zlecenie.totalSteps === 0) {
    return 'rdzeń nie podał etapów';
  }
  return `${String(zlecenie.currentStep ?? 0)} z ${String(zlecenie.totalSteps)}`;
}

/**
 * Priorytet zlecenia — pole liczbowe zapisywane komendą sterowania.
 *
 * Pole wysyła sterowanie `none`, choć rdzeń zapisuje priorytet wyłącznie przy
 * sterowaniu innym niż `none`. Zmiana kolejności obsługi nie jest zmianą stanu
 * zlecenia, a kontrakt nie zna wartości `AssistantActionControl` znaczącej
 * „zapisz sam priorytet" — są tylko `none`, `pause`, `resume`, `cancel`
 * i `retry`, a `assistant.action.status` jest jedyną komendą obszaru niosącą
 * pole `priority`.
 *
 * Podszycie się pod `pause` albo `retry` po to, żeby przemycić priorytet,
 * przestawiłoby stan zlecenia, o który nikt nie prosił. Pole zostaje więc przy
 * `none`, a rozbieżność nazywa wiersz odpowiedzi okna (`skutek-sterowania.ts`).
 */
function komorkaPriorytetu(
  zlecenie: AssistantAction,
  naSterowanie: NaSterowanie,
): HTMLTableCellElement {
  const pole = document.createElement('input');
  pole.type = 'number';
  pole.className = 'dn-pole-kontrolka ma-zlecenia__priorytet';
  pole.value = String(zlecenie.priority ?? 0);
  pole.setAttribute('aria-label', `Priorytet zlecenia ${zlecenie.title ?? zlecenie.id}`);
  pole.addEventListener('change', () => {
    naSterowanie(zlecenie.id, AssistantActionControl.None, Number(pole.value));
  });

  const element = document.createElement('td');
  element.append(pole);
  return element;
}

/** Panel akcji wiersza: cztery sterowania rdzenia oraz dwa podglądy lokalne. */
function komorkaAkcji(
  zlecenie: AssistantAction,
  naSterowanie: NaSterowanie,
  naPodglad: NaPodglad,
): HTMLTableCellElement {
  const panel = document.createElement('div');
  panel.className = 'ma-zlecenia__akcje';

  for (const [etykieta, sterowanie] of STEROWANIA) {
    panel.append(guzik(etykieta, () => naSterowanie(zlecenie.id, sterowanie)));
  }
  panel.append(
    guzik('Szczegóły', () => naPodglad(zlecenie, false)),
    guzik('Pokaż wynik', () => naPodglad(zlecenie, true)),
  );

  const element = document.createElement('td');
  element.append(panel);
  return element;
}

function guzik(etykieta: string, naKlik: () => void): HTMLButtonElement {
  const element = document.createElement('button');
  element.type = 'button';
  element.className = 'dn-btn dn-btn--sm dn-btn--zarys';
  element.textContent = etykieta;
  element.addEventListener('click', naKlik);
  return element;
}
