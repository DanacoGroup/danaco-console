import type { Command, RequestOf, ResponseOf } from '../../../shared/contract';
import type { Kanal, Wynik } from '../protokol/kanal';

/** Wysłanie komendy izolacji do rdzenia jako obietnica — jeden wspólny byt dla całego okna Punktów Izolacji. */
export function posijKomende<K extends Command>(
  kanal: Kanal,
  komenda: K,
  zadanie: RequestOf<K>,
): Promise<Wynik<ResponseOf<K>>> {
  return new Promise((rozstrzygnij) => {
    kanal.wyslij(komenda, zadanie, (wynik) => rozstrzygnij(wynik));
  });
}

/** Identyfikator sesji do żądań albo wartość pusta, gdy rdzeń jeszcze w ogóle tej sesji nie założył wcale. */
export function idSesji(kanal: Kanal): string | undefined {
  const id = kanal.sesja().id();
  return id === '' ? undefined : id;
}
