import {
  Command,
  EventType,
  type Account,
  type AccountAddRequest,
  type AccountAddResponse,
  type AccountChangedEvent,
  type AccountDefaultSetRequest,
  type AccountDefaultSetResponse,
  type AccountListRequest,
  type AccountRemoveRequest,
  type AccountRemoveResponse,
  type AccountUpdateRequest,
  type AccountUpdateResponse,
} from '../../../shared/contract';
import type { Odsubskrybuj } from '../polaczenie/magistrala-zdarzen';
import type { Kanal, Wynik } from '../protokol/kanal';
import { czyLiczba, czyObiekt, czyTablica, sprawdzKsztalt } from '../protokol/ksztalt-odpowiedzi';
import { wywolaj } from '../protokol/wywolanie';

/**
 * Rejestr kont modeli i kont programów code CLI widziany przez klienta —
 * pięć komend obszaru `account.*` oraz zdarzenie `account.changed`.
 *
 * Poświadczenie idzie w jedną stronę: kontrakt przyjmuje `credential` w treści
 * żądania i nie zwraca go nigdy, a odpowiedź niesie wyłącznie znacznik
 * `hasCredential`. Moduł nie ma ścieżki odczytu poświadczenia, bo nie istnieje
 * komenda, którą mógłby o nie spytać.
 *
 * Odczyt jest fail-open: wykaz, który nie dotarł, wraca pusty i zostawia wpis
 * w dzienniku, a widok podaje, że rdzeń nie odpowiedział, i pozostaje czynny.
 * Zapisy oddają pełny `Wynik`, bo ich niepowodzenie musi stanąć przy formularzu,
 * a nie zniknąć w konsoli.
 */
export interface ZrodloKont {
  /** `account.list` — wykaz kont w kolejności katalogu. */
  lista(zadanie?: AccountListRequest): Promise<Account[]>;
  /** `account.add` — założenie konta wraz z poświadczeniem, jeżeli podane. */
  dodaj(zadanie: AccountAddRequest): Promise<Wynik<AccountAddResponse>>;
  /** `account.update` — zmiana konta; pole pominięte zostaje bez zmiany. */
  zmien(zadanie: AccountUpdateRequest): Promise<Wynik<AccountUpdateResponse>>;
  /** `account.remove` — usunięcie konta wraz z odłączeniem kanałów. */
  usun(zadanie: AccountRemoveRequest): Promise<Wynik<AccountRemoveResponse>>;
  /** `account.default.set` — wskazanie konta domyślnego swojego rodzaju. */
  ustawDomyslne(
    zadanie: AccountDefaultSetRequest,
  ): Promise<Wynik<AccountDefaultSetResponse>>;
  /** Subskrypcja zdarzenia `account.changed`. */
  naZmiane(sluchacz: (tresc: AccountChangedEvent) => void): Odsubskrybuj;
}

/**
 * Powiadomienie o niepowodzeniu odczytu rejestru.
 *
 * Odczyt pozostaje fail-open — wykaz wraca pusty, a sekcja jest czynna — lecz
 * o niepowodzeniu dowiaduje się także widok: rejestr pusty po odmowie rdzenia
 * i rejestr pusty na świeżej instalacji to dwa różne stany.
 */
export type NaNiepowodzenie = (powod: string) => void;

export function utworzZrodloKont(
  kanal: Kanal,
  naNiepowodzenie: NaNiepowodzenie = () => undefined,
): ZrodloKont {
  return {
    async lista(zadanie = {}) {
      const wynik = sprawdzKsztalt(
        await wywolaj(kanal, Command.AccountList, zadanie),
        Command.AccountList,
        (tresc) => czyTablica(tresc.accounts) && czyLiczba(tresc.total),
      );
      if (!wynik.udany) {
        console.warn('[modele] wykaz kont nie dotarł', wynik.blad?.message ?? '');
        naNiepowodzenie(powodOdczytu(Command.AccountList, wynik.blad?.message));
        return [];
      }
      return uporzadkujKonta(wynik.wynik?.accounts ?? []);
    },

    dodaj: (zadanie) =>
      wywolaj(kanal, Command.AccountAdd, zadanie).then((wynik) =>
        sprawdzKsztalt(wynik, Command.AccountAdd, (tresc) => czyObiekt(tresc.account)),
      ),

    zmien: (zadanie) =>
      wywolaj(kanal, Command.AccountUpdate, zadanie).then((wynik) =>
        sprawdzKsztalt(wynik, Command.AccountUpdate, (tresc) => czyObiekt(tresc.account)),
      ),

    usun: (zadanie) => wywolaj(kanal, Command.AccountRemove, zadanie),

    ustawDomyslne: (zadanie) =>
      wywolaj(kanal, Command.AccountDefaultSet, zadanie).then((wynik) =>
        sprawdzKsztalt(wynik, Command.AccountDefaultSet, (tresc) =>
          czyObiekt(tresc.account),
        ),
      ),

    naZmiane: (sluchacz) => kanal.naZdarzenie(EventType.AccountChanged, sluchacz),
  };
}

/**
 * Porządek wykazu bierzemy z kolejności nadanej przez rdzeń, a przy jej braku
 * z nazwy. Konto bez kolejności nie znika — trafia na koniec.
 */
export function uporzadkujKonta(konta: readonly Account[]): Account[] {
  return [...konta].sort((pierwsze, drugie) => {
    const roznica = kolejnoscKonta(pierwsze) - kolejnoscKonta(drugie);
    return roznica !== 0 ? roznica : pierwsze.name.localeCompare(drugie.name, 'pl');
  });
}

function kolejnoscKonta(konto: Account): number {
  return Number.isFinite(konto.order) ? konto.order : Number.MAX_SAFE_INTEGER;
}

/**
 * Zdanie o niepowodzeniu odczytu pokazywane w widoku. Gdy rdzeń nie podał powodu,
 * zdanie wskazuje przynajmniej komendę, która nie odpowiedziała.
 */
export function powodOdczytu(komenda: string, powod: string | undefined): string {
  const tresc = powod !== undefined && powod !== '' ? powod : 'rdzeń nie podał powodu';
  return `Komenda ${komenda} nie odpowiedziała: ${tresc}`;
}
