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
 * Straż odmów — jedyna droga modułu Library do rdzenia.
 *
 * Rdzeń odpowiada na komendę bez uchwytu kopertą zdarzenia `<obszar>.unknown`
 * z identyfikatorem żądania, ale bez pola `status`
 * (`server/internal/protocol/zadanie.go`). Korelacja klienta rozstrzyga
 * wyłącznie koperty ze statusem (`protokol/koperta.ts`), więc obietnica
 * zwykłego `wywolaj()` po takiej odmowie nigdy się nie rozstrzyga, a okno
 * zostaje w stanie ładowania.
 *
 * Straż wiąże odmowę z żądaniem po `requestId` i zamienia ją w zwykły `Wynik`
 * z polem `blad`. Dzięki temu okno obsługuje odmowę tą samą drogą co każde inne
 * niepowodzenie i mówi wprost, której komendy rdzeń nie zna.
 *
 * Sześć obszarów własnych, bo tyle woła moduł: własny `library`, okno
 * komunikacji (`window.action`, `window.state.get`), przenoszenie kontekstu
 * (`context`), katalog akcji (`action`), katalog modułów (`module.list` —
 * obsadza ster modułu docelowego) oraz komplet kontekstu okna
 * (`aod.context.get` — jedyna komenda kontraktu czytająca to, co przyniosło
 * przekazanie).
 *
 * Siódma subskrypcja idzie na obszar zapasowy `connection`, bo rodzina
 * `knowledge.*` — wyszukiwanie po znaczeniu i wskaźnik znaczenia biblioteki —
 * nie ma własnego zdarzenia odmowy. Wykaz `zdarzeniaNieznanej`
 * (`shared/contract.go`) nie zna klucza `knowledge`, więc rdzeń bez wpiętego
 * portu Wiedzy odpowie `connection.unknown`. Bez tej subskrypcji obietnica
 * takiego wywołania nigdy by się nie rozstrzygnęła, a okno zostałoby
 * w ładowaniu. Wiązanie idzie po `requestId`, więc cudza odmowa obszaru
 * zapasowego niczego tu nie rozstrzyga.
 */
export interface StrazOdmow {
  /** Wysyła komendę kontraktu; odmowa rdzenia wraca jako `Wynik` z błędem. */
  wywolaj<K extends Command>(komenda: K, zadanie: RequestOf<K>): Promise<Wynik<ResponseOf<K>>>;
  /** Odpina subskrypcje zdarzeń odmowy. */
  rozlacz(): void;
}

/** Odbiorca wyniku bez wiedzy o kształcie treści — mapa oczekujących jest jedna. */
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

  // Subskrypcje wypisane po jednej, a nie złożone pętlą: kształt treści
  // zdarzenia bierze się z jego nazwy, więc pętla po wykazie zgubiłaby typ
  // ładunku i kazałaby go rzutować na ślepo.
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
 * Kod `not_found`: brakuje nie bytu wskazanego w żądaniu, lecz uchwytu
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
