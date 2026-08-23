import { RoundtableTurnStatus, type RoundtableTurn } from '../../../../shared/contract';

/**
 * Zegar tury debaty — czas mierzony, nie odliczany.
 *
 * Opracowanie modułu opisuje w Moderator Panelu timer tury z paskiem postępu
 * i limitem czasu. Limitu nie ma z czego wziąć: `RoundtableTurn` niesie chwilę
 * rozpoczęcia i chwilę zamknięcia, a pola limitu czasu nie ma ani tura, ani
 * żądanie tury, ani czynność moderatora. Pasek postępu wymaga dwóch końców,
 * więc zamiast niego stoi licznik: ile czasu upłynęło od otwarcia tury.
 *
 * Licznik nie jest ozdobą — moderator zamyka turę ręcznie i to jest jedyna
 * wielkość, po której poznaje, jak długo tura trwa. Tura zamknięta ma czas
 * trwania ostateczny (od `startedAt` do `closedAt`), więc licznik zatrzymuje
 * się na nim zamiast rosnąć w nieskończoność.
 *
 * `startedAt` równe zeru znaczy „rdzeń chwili nie podał”, nie „tura zaczęła się
 * w chwili zero epoki” — wtedy zegar mówi o braku pomiaru, a nie o pięćdziesięciu
 * latach trwania tury.
 */

/** Odstęp odświeżania licznika — sekunda, bo licznik pokazuje sekundy. */
export const ODSTEP_ZEGARA_MS = 1000;

/** Czy znacznik czasu kontraktu niesie chwilę, czy brak wiedzy. */
function chwilaZnana(znacznik: number | undefined): znacznik is number {
  return znacznik !== undefined && Number.isFinite(znacznik) && znacznik > 0;
}

/** Czas trwania w postaci minuty:sekundy, z zerem wiodącym w sekundach. */
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
