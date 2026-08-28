import {
  AssistantActionControl,
  type AssistantAction,
} from '../../../../shared/contract';
import { NAZWY_DROG, NAZWY_STANOW, PLAKIETKI_STANOW } from './etykiety-assistant';

// Wiersz tabeli zleceń Actions Monitora wraz z panelem akcji; sam niczego nie wywołuje.

/**
 * Zamiar wydany z wiersza: sterowanie zleceniem albo zmiana jego priorytetu.
 * Wiersz oddaje zamiar oknu, a okno prowadzi jedno wywołanie komendy
 * `assistant.action.status`.
 */
export type NaSterowanie = (
  idZlecenia: string,
  sterowanie: AssistantActionControl,
  priorytet?: number,
) => void;

/**
 * Podgląd szczegółów albo wyniku zlecenia, obsługiwany po stronie okna. Wartość
 * logiczna rozstrzyga, który z dwóch podglądów ma się otworzyć; oba są czynnością
 * lokalną i na drut do rdzenia nie idą.
 */
export type NaPodglad = (zlecenie: AssistantAction, wynik: boolean) => void;

/**
 * Sterowania panelu akcji wiersza w kolejności pokazywanej na ekranie: napis
 * przycisku wraz z wartością wyliczenia `AssistantActionControl`, która pojedzie
 * do rdzenia po naciśnięciu.
 */
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

/**
 * Składa komórkę stanu: plakietkę z nazwą stanu zlecenia oraz nazwę drogi, którą
 * zlecenie przyszło. Stan stoi w plakietce napisem, więc barwa nigdy nie niesie
 * znaczenia sama.
 */
function komorkaStanu(zlecenie: AssistantAction): HTMLTableCellElement {
  const plakietka = document.createElement('span');
  plakietka.className = PLAKIETKI_STANOW[zlecenie.status];
  plakietka.dataset['stan'] = zlecenie.status;
  plakietka.textContent = NAZWY_STANOW[zlecenie.status];

  const element = document.createElement('td');
  element.append(plakietka, ` · ${NAZWY_DROG[zlecenie.origin]}`);
  return element;
}

/**
 * Składa opis etapu zlecenia wieloetapowego z pól `currentStep` oraz `totalSteps`.
 * Zlecenie bez podanej liczby etapów dostaje zdanie o tym, że rdzeń etapów nie
 * podał, zamiast pustej komórki.
 */
function opisEtapu(zlecenie: AssistantAction): string {
  if (zlecenie.totalSteps === undefined || zlecenie.totalSteps === 0) {
    return 'rdzeń nie podał etapów';
  }
  return `${String(zlecenie.currentStep ?? 0)} z ${String(zlecenie.totalSteps)}`;
}

/**
 * Składa komórkę priorytetu: pole liczbowe, którego zmiana wydaje zamiar
 * sterowania `none` wraz z nową wartością pola `priority`. Kontrakt nie zna
 * sterowania znaczącego wyłącznie zapis priorytetu.
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

/**
 * Składa komórkę akcji: cztery sterowania rdzenia z wykazu `STEROWANIA` oraz dwa
 * podglądy prowadzone lokalnie przez okno. Wszystkie przyciski zostają klikalne
 * niezależnie od stanu zlecenia.
 */
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
