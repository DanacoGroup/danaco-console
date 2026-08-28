import { RoundtableTurnStatus, type RoundtableTurn } from '../../../../shared/contract';

/** Zegar tury debaty — czas mierzony, nie odliczany: odstęp odświeżania licznika to sekunda, bo licznik pokazuje sekundy. */
export const ODSTEP_ZEGARA_MS = 1000;

/** Czy znacznik czasu kontraktu niesie chwilę, czy brak wiedzy — zero i wartość ujemna znaczą brak, nie chwilę zero. */
function chwilaZnana(znacznik: number | undefined): znacznik is number {
  return znacznik !== undefined && Number.isFinite(znacznik) && znacznik > 0;
}

/** Czas trwania w postaci minuty:sekundy, z zerem wiodącym w sekundach, licząc od zera przy wartości ujemnej. */
export function czasTrwania(milisekundy: number): string {
  const sekundy = Math.max(0, Math.floor(milisekundy / 1000));
  const minuty = Math.floor(sekundy / 60);
  const reszta = sekundy - minuty * 60;
  return `${minuty}:${String(reszta).padStart(2, '0')}`;
}

/**
 * Zdanie zegara dla tury bieżącej.
 *
 * @param teraz chwila odczytu, wnoszona parametrem — funkcja zostaje czysta,
 *   a okno decyduje, jak często ją woła.
 */
export function zdanieZegara(tura: RoundtableTurn | null, teraz: number): string {
  if (tura === null) {
    return 'Debata nie ma otwartej tury — zegar nie ma czego mierzyć.';
  }
  if (!chwilaZnana(tura.startedAt)) {
    return `Tura #${tura.index}: rdzeń nie podał chwili otwarcia, więc czasu trwania nie ma z czego policzyć.`;
  }
  if (tura.status === RoundtableTurnStatus.Closed) {
    if (!chwilaZnana(tura.closedAt)) {
      return `Tura #${tura.index} jest zamknięta, ale rdzeń nie podał chwili zamknięcia — czasu trwania nie ma z czego policzyć.`;
    }
    return `Tura #${tura.index} zamknięta. Trwała ${czasTrwania(tura.closedAt - tura.startedAt)}.`;
  }
  return (
    `Tura #${tura.index} otwarta. Od jej otwarcia upłynęło ${czasTrwania(teraz - tura.startedAt)}. ` +
    'Zegar mierzy czas, nie odlicza go: limitu czasu tury nie niesie żadne pole kontraktu.'
  );
}
