import {
  AutomationExecutionStatus,
  type AutomationExecution,
} from '../../../../shared/contract';

/**
 * Przegląd przebiegów okna Execution Monitor — zawężenie wykazu i miary
 * niezawodności.
 *
 * Obie czynności idą po stronie klienta i idą tam z tego samego powodu:
 * `automation.execution.subscribe` przyjmuje wyłącznie automatykę, pojedynczy
 * przebieg, okno i górną granicę, a oddaje przebiegi w całości. Nie ma pola
 * zapytania, zakresu dat ani zestawienia miar — a wykaz jest już w oknie, więc
 * pytanie rdzenia o to samo drugi raz nic by nie wniosło.
 *
 * Miary liczą się z pól, które kontrakt naprawdę niesie: stanu przebiegu oraz
 * znaczników czasu rozpoczęcia i zakończenia. Miary kosztu — tokeny i wywołania
 * modelu — składa widok, bo ich pola są nieobowiązkowe i suma wymaga podania
 * obok siebie liczby przebiegów, które je wypełniły (`widok-przebiegow.ts`).
 *
 * Plik jest czysty: nie dotyka dokumentu i nie woła rdzenia.
 */

/** Zawężenie wykazu przebiegów — wartości pól paska narzędzi okna. */
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

/** Zawężenie puste — wykaz w całości. */
export const BEZ_ZAWEZENIA: ZawezeniePrzebiegow = {
  stan: '',
  odDnia: '',
  doDnia: '',
  szukane: '',
};

/** Miary wykazu przebiegów; pole puste znaczy „nie ma z czego policzyć”. */
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

/** Zawęża wykaz przebiegów; zawężenie puste oddaje wykaz w całości. */
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

/** Liczy miary wykazu; wykaz pusty daje same zera i pola puste. */
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

/** Czas trwania w mowie Operatora: minuty i sekundy albo same sekundy. */
export function czasTrwania(milisekundy: number): string {
  const sekundy = Math.round(milisekundy / 1000);
  if (sekundy < 60) return `${sekundy} s`;
  return `${Math.floor(sekundy / 60)} min ${String(sekundy % 60).padStart(2, '0')} s`;
}

/** Początek dnia w czasie miejscowym; zapis nieczytelny daje `null`. */
function poczatekDnia(dzien: string): number | null {
  if (dzien.trim() === '') return null;
  const czas = new Date(`${dzien.trim()}T00:00:00`).getTime();
  return Number.isNaN(czas) ? null : czas;
}

/** Koniec dnia w czasie miejscowym; zapis nieczytelny daje `null`. */
function koniecDnia(dzien: string): number | null {
  if (dzien.trim() === '') return null;
  const czas = new Date(`${dzien.trim()}T23:59:59.999`).getTime();
  return Number.isNaN(czas) ? null : czas;
}
