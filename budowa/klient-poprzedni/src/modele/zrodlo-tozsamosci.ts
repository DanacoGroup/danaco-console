import {
  Command,
  EventType,
  type IdentityCategory,
  type IdentityCategoryListRequest,
  type IdentityChangedEvent,
  type IdentityDocument,
  type IdentityDocumentGetRequest,
  type IdentityDocumentRemoveRequest,
  type IdentityDocumentSetRequest,
  type IdentityDocumentSetResponse,
  type IdentityEffectiveGetRequest,
  type IdentityEffectiveGetResponse,
} from '../../../shared/contract';
import type { Odsubskrybuj } from '../polaczenie/magistrala-zdarzen';
import type { Kanal, Wynik } from '../protokol/kanal';
import { czyObiekt, czyTablica, czyTekst, sprawdzKsztalt } from '../protokol/ksztalt-odpowiedzi';
import { wywolaj } from '../protokol/wywolanie';
import { pierwszenstwoWarstwy } from './warstwy-tozsamosci';

/**
 * Źródło tożsamości modelu udostępnia klientowi pięć komend obszaru
 * `identity.*` oraz zdarzenie `identity.changed`. Katalog kategorii przychodzi
 * z rdzenia jako dane, a nakładkę obowiązującą składa rdzeń, nie widok.
 */
export interface ZrodloTozsamosci {
  /** `identity.category.list` — katalog kategorii zasad. */
  kategorie(zadanie?: IdentityCategoryListRequest): Promise<IdentityCategory[]>;
  /** `identity.document.get` — treści zapisane dla wskazanej osi. */
  dokumenty(zadanie: IdentityDocumentGetRequest): Promise<IdentityDocument[]>;
  /** `identity.document.set` — zapis treści kategorii wraz z trybem podania. */
  zapisz(
    zadanie: IdentityDocumentSetRequest,
  ): Promise<Wynik<IdentityDocumentSetResponse>>;
  /** `identity.document.remove` — zdjęcie zapisu; treść wraca z osi szerszej. */
  usun(zadanie: IdentityDocumentRemoveRequest): Promise<Wynik<unknown>>;
  /** `identity.effective.get` — nakładka obowiązująca wraz ze złożonym promptem. */
  nakladka(
    zadanie: IdentityEffectiveGetRequest,
  ): Promise<Wynik<IdentityEffectiveGetResponse>>;
  /** Subskrypcja zdarzenia `identity.changed`. */
  naZmiane(sluchacz: (tresc: IdentityChangedEvent) => void): Odsubskrybuj;
}

export function utworzZrodloTozsamosci(kanal: Kanal): ZrodloTozsamosci {
  return {
    async kategorie(zadanie = {}) {
      const wynik = sprawdzKsztalt(
        await wywolaj(kanal, Command.IdentityCategoryList, zadanie),
        Command.IdentityCategoryList,
        (tresc) => czyTablica(tresc.categories),
      );
      if (!wynik.udany) {
        console.warn('[modele] katalog kategorii zasad nie dotarł', wynik.blad?.message ?? '');
        return [];
      }
      return uporzadkujKategorie(wynik.wynik?.categories ?? []);
    },

    async dokumenty(zadanie) {
      const wynik = sprawdzKsztalt(
        await wywolaj(kanal, Command.IdentityDocumentGet, zadanie),
        Command.IdentityDocumentGet,
        (tresc) => czyTablica(tresc.documents),
      );
      if (!wynik.udany) {
        console.warn('[modele] treści kategorii nie dotarły', wynik.blad?.message ?? '');
        return [];
      }
      return [...(wynik.wynik?.documents ?? [])];
    },

    zapisz: (zadanie) =>
      wywolaj(kanal, Command.IdentityDocumentSet, zadanie).then((wynik) =>
        sprawdzKsztalt(wynik, Command.IdentityDocumentSet, (tresc) =>
          czyObiekt(tresc.document),
        ),
      ),

    usun: (zadanie) => wywolaj(kanal, Command.IdentityDocumentRemove, zadanie),

    nakladka: (zadanie) =>
      wywolaj(kanal, Command.IdentityEffectiveGet, zadanie).then((wynik) =>
        sprawdzKsztalt(
          wynik,
          Command.IdentityEffectiveGet,
          (tresc) => czyTablica(tresc.layers) && czyTekst(tresc.prompt),
        ),
      ),

    naZmiane: (sluchacz) => kanal.naZdarzenie(EventType.IdentityChanged, sluchacz),
  };
}

/**
 * Porządek kategorii: najpierw warstwa nadana przez katalog, potem kolejność
 * wewnątrz warstwy. Kategoria bez kolejności trafia na koniec swojej warstwy,
 * zamiast zniknąć.
 */
export function uporzadkujKategorie(
  kategorie: readonly IdentityCategory[],
): IdentityCategory[] {
  return [...kategorie].sort((pierwsza, druga) => {
    const warstwy =
      pierwszenstwoWarstwy(pierwsza.layer) - pierwszenstwoWarstwy(druga.layer);
    return warstwy !== 0 ? warstwy : kolejnosc(pierwsza) - kolejnosc(druga);
  });
}

function kolejnosc(kategoria: IdentityCategory): number {
  return Number.isFinite(kategoria.order) ? kategoria.order : Number.MAX_SAFE_INTEGER;
}
