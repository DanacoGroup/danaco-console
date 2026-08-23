import type { Command, RequestOf, ResponseOf } from '../../../shared/contract';
import type { Kanal, Wynik } from '../protokol/kanal';

/**
 * Wysłanie komendy `isolation.*` jako obietnica — jeden byt dla całego okna.
 *
 * `Kanal.wyslij` przyjmuje wywołanie zwrotne, a obszary tego okna czytają po
 * kilka komend po kolei (najpierw wykaz, potem odczyt, potem zapis) — splot
 * wywołań zwrotnych w takim ciągu jest nieczytelny. Wszystkie obszary okna
 * wołają więc stąd, zamiast powtarzać u siebie ten sam opakowujący zapis.
 *
 * Typowanie zostaje pełne: parametr `K extends Command` wiąże żądanie
 * (`RequestOf<K>`) z odpowiedzią (`ResponseOf<K>`) wprost z kontraktu, więc
 * pomyłka w kształcie żądania przerywa kompilację klienta zamiast wracać
 * odmową rdzenia w czasie działania.
 *
 * Odmowa nie jest wyjątkiem: obietnica spełnia się zawsze, także gdy rdzeń
 * odmówił. Odmowa siedzi w `wynik.udany === false` i w `wynik.blad`, żeby
 * wywołujący rozstrzygnął ją zdaniem dla Operatora (`opisOdmowyBledu`),
 * a nie `catch`-em, który łatwo pominąć.
 */
export function posijKomende<K extends Command>(
  kanal: Kanal,
  komenda: K,
  zadanie: RequestOf<K>,
): Promise<Wynik<ResponseOf<K>>> {
  return new Promise((rozstrzygnij) => {
    kanal.wyslij(komenda, zadanie, (wynik) => rozstrzygnij(wynik));
  });
}

/** Identyfikator sesji do żądań albo `undefined`, gdy rdzeń jeszcze sesji nie założył. */
export function idSesji(kanal: Kanal): string | undefined {
  const id = kanal.sesja().id();
  return id === '' ? undefined : id;
}
