/**
 * Przegląd przebiegów okna Execution Monitor: zawężenie wykazu i miary
 * niezawodności liczone po stronie klienta z przebiegów już obecnych w oknie.
 * Plik nie dotyka dokumentu i nie woła rdzenia.
 */

import {
  AutomationExecutionStatus,
  type AutomationExecution,
} from '../../../../shared/contract';

/**
 * Zawężenie wykazu przebiegów złożone z wartości pól paska narzędzi okna:
 * stanu przebiegu, granic zakresu dni oraz napisu szukanego w wykazie.
 */
export interface ZawezeniePrzebiegow {
  /** Stan przebiegu; pusty napis znaczy „wszystkie stany”. */
  stan: string;
  /** Dzień początkowy w zapisie ISO (RRRR-MM-DD); pusty znaczy „bez granicy”. */
  odDnia: string;
  /** Dzień końcowy w zapisie ISO (RRRR-MM-DD); pusty znaczy „bez granicy”. */
  doDnia: string;
  /** Napis szukany w identyfikatorze przebiegu i w powodzie niepowodzenia. */
  szukane: string;
}

/**
 * Zawężenie puste, w którym wszystkie pola są pustymi napisami; jako stan
 * wyjściowy paska narzędzi oddaje wykaz przebiegów w całości.
 */
export const BEZ_ZAWEZENIA: ZawezeniePrzebiegow = {
  stan: '',
  odDnia: '',
  doDnia: '',
  szukane: '',
};

/**
 * Miary wykazu przebiegów: liczby przebiegów w poszczególnych stanach oraz
 * wskaźnik powodzenia, średni czas trwania i opóźnienie kwantyla; pole puste
 * znaczy, że nie ma z czego liczyć.
 */
export interface MiaryPrzebiegow {
  liczba: number;
  udane: number;
  nieudane: number;
  wToku: number;
  /** Udział przebiegów udanych wśród zamkniętych, w procentach. */
  wskaznikPowodzenia: number | null;
  /** Średni czas trwania przebiegu zamkniętego, w milisekundach. */
  sredniCzasMs: number | null;
  /** Czas trwania, poniżej którego mieści się 95 setnych przebiegów zamkniętych. */
  opoznienieP95Ms: number | null;
}

/**
 * Zawęża wykaz przebiegów wartościami zawężenia: stanem, granicami dni oraz
 * napisem szukanym w identyfikatorze i w powodzie niepowodzenia; zawężenie
 * puste oddaje wykaz w całości.
 */
export function zawezonePrzebiegi(
  przebiegi: readonly AutomationExecution[],
  zawezenie: ZawezeniePrzebiegow,
): AutomationExecution[] {
  const igla = zawezenie.szukane.trim().toLocaleLowerCase('pl-PL');
  const od = poczatekDnia(zawezenie.odDnia);
  const doo = koniecDnia(zawezenie.doDnia);
  return przebiegi.filter((przebieg) => {
    if (zawezenie.stan !== '' && przebieg.status !== zawezenie.stan) return false;
    if (od !== null && przebieg.startedAt < od) return false;
    if (doo !== null && przebieg.startedAt > doo) return false;
    if (igla === '') return true;
    const stog = `${przebieg.id} ${przebieg.workflowId} ${przebieg.errorMessage ?? ''}`;
    return stog.toLocaleLowerCase('pl-PL').includes(igla);
  });
}

/**
 * Liczy miary wykazu przebiegów: liczby w stanach, wskaźnik powodzenia wśród
 * zamkniętych, średni czas trwania oraz opóźnienie kwantyla; wykaz pusty daje
 * same zera i pola puste.
 */
export function miaryPrzebiegow(przebiegi: readonly AutomationExecution[]): MiaryPrzebiegow {
  const udane = przebiegi.filter(
    (przebieg) => przebieg.status === AutomationExecutionStatus.Succeeded,
  ).length;
  const nieudane = przebiegi.filter(
    (przebieg) => przebieg.status === AutomationExecutionStatus.Failed,
  ).length;
  const wToku = przebiegi.filter(
    (przebieg) => przebieg.status === AutomationExecutionStatus.Running,
  ).length;
  const zamkniete = udane + nieudane;
  const czasy = przebiegi
    .filter((przebieg) => przebieg.finishedAt !== undefined)
    .map((przebieg) => (przebieg.finishedAt ?? 0) - przebieg.startedAt)
    .filter((czas) => czas >= 0)
    .sort((pierwszy, drugi) => pierwszy - drugi);

  return {
    liczba: przebiegi.length,
    udane,
    nieudane,
    wToku,
    wskaznikPowodzenia: zamkniete === 0 ? null : Math.round((udane / zamkniete) * 100),
    sredniCzasMs:
      czasy.length === 0
        ? null
        : Math.round(czasy.reduce((suma, czas) => suma + czas, 0) / czasy.length),
    opoznienieP95Ms: czasy.length === 0 ? null : (czasy[kwantyl(czasy.length, 95)] ?? null),
  };
}

/**
 * Miejsce kwantyla w wykazie uporządkowanym — metoda najbliższej pozycji.
 * Wykaz o jednej pozycji ma ją jedyną także dla kwantyla 95.
 */
function kwantyl(dlugosc: number, procent: number): number {
  const miejsce = Math.ceil((procent / 100) * dlugosc) - 1;
  return Math.min(dlugosc - 1, Math.max(0, miejsce));
}

/**
 * Zapisuje czas trwania przebiegu w postaci czytelnej w oknie: dla czasów od
 * minuty wzwyż minuty wraz z sekundami, a poniżej minuty same sekundy.
 */
export function czasTrwania(milisekundy: number): string {
  const sekundy = Math.round(milisekundy / 1000);
  if (sekundy < 60) return `${sekundy} s`;
  return `${Math.floor(sekundy / 60)} min ${String(sekundy % 60).padStart(2, '0')} s`;
}

/**
 * Początek podanego dnia w czasie miejscowym, w milisekundach epoki; napis
 * pusty albo zapis nieczytelny dla konstruktora daty daje `null`.
 */
function poczatekDnia(dzien: string): number | null {
  if (dzien.trim() === '') return null;
  const czas = new Date(`${dzien.trim()}T00:00:00`).getTime();
  return Number.isNaN(czas) ? null : czas;
}

/**
 * Koniec podanego dnia w czasie miejscowym, w milisekundach epoki i z
 * dokładnością do milisekundy; napis pusty albo zapis nieczytelny daje `null`.
 */
function koniecDnia(dzien: string): number | null {
  if (dzien.trim() === '') return null;
  const czas = new Date(`${dzien.trim()}T23:59:59.999`).getTime();
  return Number.isNaN(czas) ? null : czas;
}
