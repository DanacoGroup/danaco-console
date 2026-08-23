import type { RoundtableStatement, RoundtableTurn } from '../../../../shared/contract';

/**
 * Metryki odpowiedzi uczestnika — wyłącznie to, co da się zmierzyć z kontraktu.
 *
 * Opracowanie modułu wymienia przy Model Panels cztery metryki: długość, czas
 * odpowiedzi, koszt tokenów i szacowaną pewność. Zmierzyć da się dwie pierwsze:
 * długość liczy się z treści, którą okno ma przed sobą, a czas z dwóch znaczników
 * kontraktu — `RoundtableTurn.startedAt` i `RoundtableStatement.createdAt`.
 * Pewność ma dziś w kontrakcie pole (`RoundtableStatement.confidence`), lecz
 * rdzeń go jeszcze nie wypełnia; licznika tokenów nie niesie żadne pole obszaru
 * `roundtable.*`. Okno mówi o obu wprost, zamiast wypełniać rubrykę kreską.
 *
 * Czas mierzy odstęp między otwarciem tury a zapisaniem wypowiedzi, a nie czas
 * pracy modelu: chwili rozpoczęcia wypowiedzi kontrakt nie niesie. Nazwa metryki
 * mówi dokładnie to, co metryka liczy.
 *
 * Zero i wartość niebędąca liczbą znaczą „rdzeń chwili nie podał", nie „zdarzyło
 * się w chwili zero" — obie prowadzą do braku pomiaru, nie do liczby ujemnej
 * albo do daty z początku epoki.
 */

/** Pomiar jednej odpowiedzi; wartość `null` znaczy „nie ma z czego zmierzyć". */
export interface MetrykiOdpowiedzi {
  znakow: number;
  wyrazow: number;
  /** Milisekundy od otwarcia tury do zapisania wypowiedzi; `null`, gdy nieznane. */
  odOtwarciaTury: number | null;
}

/** Czy znacznik czasu kontraktu niesie chwilę, czy brak wiedzy. */
function chwilaZnana(znacznik: number | undefined): znacznik is number {
  return znacznik !== undefined && Number.isFinite(znacznik) && znacznik > 0;
}

export function zmierzOdpowiedz(
  tekst: string,
  wypowiedz: RoundtableStatement | null,
  tura: RoundtableTurn | null,
): MetrykiOdpowiedzi {
  const przyciety = tekst.trim();
  const wyrazow = przyciety === '' ? 0 : przyciety.split(/\s+/u).length;
  const start = tura === null ? undefined : tura.startedAt;
  const zapis = wypowiedz === null ? undefined : wypowiedz.createdAt;
  const odOtwarciaTury =
    chwilaZnana(start) && chwilaZnana(zapis) && zapis >= start ? zapis - start : null;
  return { znakow: tekst.length, wyrazow, odOtwarciaTury };
}

/** Odstęp czasu słowem — sekundy z jednym miejscem, minuty powyżej minuty. */
export function opisOdstepu(milisekundy: number): string {
  const sekundy = milisekundy / 1000;
  if (sekundy < 60) return `${sekundy.toFixed(1)} s`;
  const minuty = Math.floor(sekundy / 60);
  const reszta = Math.round(sekundy - minuty * 60);
  return `${minuty} min ${reszta} s`;
}

/**
 * Zdania metryk gotowe do wypisania w panelu uczestnika.
 *
 * Ostatnie zdanie mówi o dwóch metrykach, których okno nie pokazuje, i o tym,
 * że każdej brakuje z innego powodu. Stoi zawsze, bo rubryka pokazująca dwie
 * metryki z czterech wygląda na komplet, dopóki nie powie, że kompletem nie jest.
 */
export function zdaniaMetryk(metryki: MetrykiOdpowiedzi): string[] {
  const zdania = [`Długość: ${metryki.znakow} znaków, ${metryki.wyrazow} wyrazów`];
  zdania.push(
    metryki.odOtwarciaTury === null
      ? 'Czas od otwarcia tury: niezmierzony — rdzeń nie podał chwili otwarcia tury albo chwili zapisania wypowiedzi'
      : `Czas od otwarcia tury do zapisania wypowiedzi: ${opisOdstepu(metryki.odOtwarciaTury)}`,
  );
  zdania.push(
    'Pewność uczestnika niesie pole confidence wypowiedzi, którego rdzeń jeszcze nie wypełnia; kosztu tokenów nie niesie żadne pole obszaru roundtable.',
  );
  return zdania;
}
