import type { RoundtableStatement, RoundtableTurn } from '../../../../shared/contract';

/** Pomiar jednej odpowiedzi, wyłącznie to, co da się zmierzyć z kontraktu; wartość null znaczy, że nie ma z czego zmierzyć. */
export interface MetrykiOdpowiedzi {
  znakow: number;
  wyrazow: number;
  /** Milisekundy od otwarcia tury do zapisania wypowiedzi; `null`, gdy nieznane. */
  odOtwarciaTury: number | null;
}

/** Czy znacznik czasu kontraktu niesie chwilę, czy brak wiedzy — zero i wartość niebędąca liczbą znaczą brak, nie zdarzenie w chwili zero. */
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

/** Odstęp czasu słowem — sekundy z jednym miejscem po przecinku, minuty i sekundy powyżej pełnej minuty. */
export function opisOdstepu(milisekundy: number): string {
  const sekundy = milisekundy / 1000;
  if (sekundy < 60) return `${sekundy.toFixed(1)} s`;
  const minuty = Math.floor(sekundy / 60);
  const reszta = Math.round(sekundy - minuty * 60);
  return `${minuty} min ${reszta} s`;
}

/** Zdania metryk gotowe do wypisania w panelu uczestnika, wraz ze zdaniem o dwóch metrykach, których okno nie pokazuje. */
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
