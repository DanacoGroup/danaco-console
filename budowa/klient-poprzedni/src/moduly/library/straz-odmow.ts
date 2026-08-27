import {
  ErrorCode,
  EventType,
  type Command,
  type RequestOf,
  type ResponseOf,
  type UnknownCommandPayload,
} from '../../../../shared/contract';
import type { Odsubskrybuj } from '../../polaczenie/magistrala-zdarzen';
import type { Kanal, Wynik } from '../../protokol/kanal';

/**
 * Straż odmów — jedyna droga modułu Library do rdzenia. Odmowa komendy wraca
 * kopertą `<obszar>.unknown` bez pola `status`, więc korelacja jej nie
 * rozstrzyga; straż wiąże ją po `requestId` i zwraca jako `Wynik` z błędem.
 */
export interface StrazOdmow {
  /** Wysyła komendę kontraktu; odmowa rdzenia wraca jako `Wynik` z błędem. */
  wywolaj<K extends Command>(komenda: K, zadanie: RequestOf<K>): Promise<Wynik<ResponseOf<K>>>;
  /** Odpina subskrypcje zdarzeń odmowy. */
  rozlacz(): void;
}

/**
 * Odbiorca wyniku bez wiedzy o kształcie treści, ponieważ mapa oczekujących na
 * odmowę jest jedna dla wszystkich obserwowanych obszarów kontraktu.
 */
type RozstrzygnijNieznane = (wynik: Wynik<never>) => void;

export function utworzStrazOdmow(kanal: Kanal): StrazOdmow {
  const oczekujace = new Map<string, RozstrzygnijNieznane>();

  function przyjmijOdmowe(tresc: UnknownCommandPayload): void {
    const idZadania = tresc.requestId ?? '';
    const rozstrzygnij = oczekujace.get(idZadania);
    if (rozstrzygnij === undefined) return;
    oczekujace.delete(idZadania);
    rozstrzygnij({ udany: false, blad: bladOdmowy(tresc) });
  }

  // Subskrypcje po jednej: pętla po wykazie zgubiłaby typ ładunku
  // i wymusiłaby rzutowanie.
  const odsubskrybowania: Odsubskrybuj[] = [
    kanal.naZdarzenie(EventType.LibraryUnknown, (tresc) => przyjmijOdmowe(tresc)),
    kanal.naZdarzenie(EventType.WindowUnknown, (tresc) => przyjmijOdmowe(tresc)),
    kanal.naZdarzenie(EventType.ContextUnknown, (tresc) => przyjmijOdmowe(tresc)),
    kanal.naZdarzenie(EventType.ActionUnknown, (tresc) => przyjmijOdmowe(tresc)),
    kanal.naZdarzenie(EventType.ModuleUnknown, (tresc) => przyjmijOdmowe(tresc)),
    kanal.naZdarzenie(EventType.AodUnknown, (tresc) => przyjmijOdmowe(tresc)),
    kanal.naZdarzenie(EventType.ConnectionUnknown, (tresc) => przyjmijOdmowe(tresc)),
  ];

  return {
    wywolaj(komenda, zadanie) {
      return new Promise((rozstrzygnij) => {
        let idZadania = '';
        idZadania = kanal.wyslij(komenda, zadanie, (wynik) => {
          oczekujace.delete(idZadania);
          rozstrzygnij(wynik);
        });
        oczekujace.set(idZadania, rozstrzygnij as RozstrzygnijNieznane);
      });
    },

    rozlacz() {
      for (const odsubskrybuj of odsubskrybowania.splice(0)) odsubskrybuj();
      oczekujace.clear();
    },
  };
}

/**
 * Błąd opisujący odmowę rdzenia.
 *
 * Kod `not_found`: brakuje nie bytu z żądania, lecz uchwytu
 * komendy — stan nieponawialny, więc `retryable` jest fałszem. Nazwa żądanego
 * typu zostaje w treści, bo bez niej nie widać, której komendy rdzeń nie zna.
 */
function bladOdmowy(tresc: UnknownCommandPayload) {
  const powod = (tresc.reason ?? '').trim();
  const ogon = powod === '' ? '' : ` (${powod})`;
  return {
    code: ErrorCode.NotFound,
    message: `Rdzeń nie ma dziś uchwytu komendy ${tresc.requestedType}${ogon}`,
    retryable: false,
  };
}
