import {
  Command,
  EventType,
  type AccessGrant,
  type AccessGrantAddRequest,
  type AccessGrantAddResponse,
  type AccessGrantChangedEvent,
  type AccessGrantListRequest,
  type AccessGrantRemoveResponse,
  type AccessGrantUpdateRequest,
  type AccessGrantUpdateResponse,
} from '../../../shared/contract';
import type { Odsubskrybuj } from '../polaczenie/magistrala-zdarzen';
import type { Kanal, Wynik } from '../protokol/kanal';
import { czyObiekt, czyTablica, sprawdzKsztalt } from '../protokol/ksztalt-odpowiedzi';
import { wywolaj } from '../protokol/wywolanie';
import { powodOdczytu, type NaNiepowodzenie } from './zrodlo-punktow';

/**
 * Nadania dostępu okna rozmowy.
 *
 * Okno ma zbiór nadań, nie jedno. Kolejność i oznaczenie głównego mają
 * znaczenie, dlatego trzy komendy zapisujące zwracają nie tylko zmienione
 * nadanie, lecz komplet nadań okna po zmianie: przestawienie kolejności albo
 * oznaczenia głównego dotyka pozostałych wierszy i widok musi je zobaczyć
 * w jednej odpowiedzi.
 *
 * Nazwy komend i zdarzeń pochodzą wyłącznie ze stałych kontraktu.
 */
export interface ZrodloNadan {
  /** `access.grant.list` — nadania okna w kolejności. */
  lista(zadanie?: AccessGrantListRequest): Promise<AccessGrant[]>;
  /** `access.grant.add` — nadanie oknu dostępu do punktu. */
  nadaj(zadanie: AccessGrantAddRequest): Promise<Wynik<AccessGrantAddResponse>>;
  /** `access.grant.update` — tryb, korzenie, kolejność albo oznaczenie głównego. */
  zmien(zadanie: AccessGrantUpdateRequest): Promise<Wynik<AccessGrantUpdateResponse>>;
  /** `access.grant.remove` — odebranie nadania. */
  odbierz(nadanieID: string): Promise<Wynik<AccessGrantRemoveResponse>>;
  /** Subskrypcja `access.grant.changed`. */
  naZmiane(sluchacz: (zdarzenie: AccessGrantChangedEvent) => void): Odsubskrybuj;
}

export function utworzZrodloNadan(
  kanal: Kanal,
  naNiepowodzenie: NaNiepowodzenie = () => undefined,
): ZrodloNadan {
  return {
    async lista(zadanie = {}) {
      const wynik = sprawdzKsztalt(
        await wywolaj(kanal, Command.AccessGrantList, zadanie),
        Command.AccessGrantList,
        (tresc) => czyTablica(tresc.grants),
      );
      if (!wynik.udany) {
        console.warn('[dostępy] nadania nie dotarły', wynik.blad?.message ?? '');
        naNiepowodzenie(powodOdczytu(Command.AccessGrantList, wynik.blad?.message));
        return [];
      }
      return uporzadkujNadania(wynik.wynik?.grants ?? []);
    },

    async nadaj(zadanie) {
      return sprawdzKomplet(
        await wywolaj(kanal, Command.AccessGrantAdd, zadanie),
        Command.AccessGrantAdd,
      );
    },

    async zmien(zadanie) {
      return sprawdzKomplet(
        await wywolaj(kanal, Command.AccessGrantUpdate, zadanie),
        Command.AccessGrantUpdate,
      );
    },

    async odbierz(nadanieID) {
      return sprawdzKsztalt(
        await wywolaj(kanal, Command.AccessGrantRemove, { grantId: nadanieID }),
        Command.AccessGrantRemove,
        (tresc) => typeof tresc.removed === 'boolean' && czyTablica(tresc.grants),
      );
    },

    naZmiane(sluchacz) {
      return kanal.naZdarzenie(EventType.AccessGrantChanged, (tresc) => {
        if (czyObiekt(tresc.grant as unknown)) sluchacz(tresc);
      });
    },
  };
}

/** Odpowiedź niosąca komplet nadań okna po zapisie. */
type OdpowiedzKompletu = { grants: AccessGrant[] };

/** Wspólny sprawdzian odpowiedzi zapisu: nadanie i komplet po zmianie. */
function sprawdzKomplet<T extends OdpowiedzKompletu & { grant: AccessGrant }>(
  wynik: Wynik<T>,
  komenda: string,
): Wynik<T> {
  return sprawdzKsztalt(
    wynik,
    komenda,
    (tresc) => czyObiekt(tresc.grant as unknown) && czyTablica(tresc.grants),
  );
}

/**
 * Porządek nadań: kolejność z rdzenia, a przy jej braku czas założenia.
 * Nadanie bez kolejności trafia na koniec zamiast zniknąć.
 */
export function uporzadkujNadania(nadania: readonly AccessGrant[]): AccessGrant[] {
  return [...nadania].sort((pierwsze, drugie) => {
    const roznica = kolejnosc(pierwsze) - kolejnosc(drugie);
    return roznica !== 0 ? roznica : pierwsze.createdAt - drugie.createdAt;
  });
}

function kolejnosc(nadanie: AccessGrant): number {
  return Number.isFinite(nadanie.order) ? nadanie.order : Number.MAX_SAFE_INTEGER;
}
