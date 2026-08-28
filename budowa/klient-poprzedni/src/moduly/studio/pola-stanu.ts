import {
  StudioOperationScope,
  type StudioDocument,
  type StudioVersion,
} from '../../../../shared/contract';

/** Zaznaczony fragment dokumentu w edytorze, będący przedmiotem operacji kontekstowej wywoływanej na tym fragmencie. */
export interface Zaznaczenie {
  poczatek: number;
  koniec: number;
}

/** Wynik operacji sztucznej inteligencji czekający na decyzję operatora: przyjęcie treści albo jej odrzucenie. */
export interface PropozycjaZmiany {
  /** Propozycja po stronie rdzenia; zasila `studio.diff.compare`. */
  idPropozycji: string;
  /** Treść wyniku operacji, o ile rdzeń ją oddał. */
  tresc: string;
  /** Operacja, która propozycję wywołała — do nazwania jej w panelu różnicy. */
  idAkcji: string;
}

/** Pola stanu modułu Studio wraz z przejściami, które je zmieniają; treść robocza i zaakceptowana to dwa różne pola. */
export interface PolaStanu {
  okno: string;
  dokument: StudioDocument | null;
  robocza: string;
  zakres: Zaznaczenie | null;
  /** Wybór ręczny zakresu operacji; null oznacza podążanie za zaznaczeniem, nastawa liczona tylko tutaj. */
  zakresReczny: StudioOperationScope | null;
  historia: readonly StudioVersion[];
  propozycja: PropozycjaZmiany | null;
  zaakceptowana: string;
  /** Powód, dla którego ostatnia operacja kontekstowa nie dała propozycji; null, gdy żadna nie odmówiła. */
  odmowaOperacji: string | null;
  /** Para wersji wskazana do porównania w Diff/Grep Panelu; null, gdy żadnej pary jeszcze nie wskazano. */
  paraPorownania: ParaPorownania | null;
}

/** Dwie strony porównania wskazane w Session Repository: wersja odniesienia i wersja porównywana z nią. */
export interface ParaPorownania {
  /** Wersja odniesienia — starsza strona porównania. */
  odniesienie: string;
  /** Wersja porównywana; pusta znaczy „wersja bieżąca dokumentu". */
  porownywana: string;
}

export function pustePolaStanu(): PolaStanu {
  return {
    okno: '',
    dokument: null,
    robocza: '',
    zakres: null,
    zakresReczny: null,
    historia: [],
    propozycja: null,
    zaakceptowana: '',
    odmowaOperacji: null,
    paraPorownania: null,
  };
}

/** Przyjmuje zaznaczenie edytora i rozstrzyga los wyboru ręcznego zakresu operacji kontekstowej modelu. */
export function ustawZaznaczenieStanu(pola: PolaStanu, zakres: Zaznaczenie | null): void {
  const inne =
    zakres !== null &&
    (pola.zakres === null ||
      pola.zakres.poczatek !== zakres.poczatek ||
      pola.zakres.koniec !== zakres.koniec);
  pola.zakres = zakres;
  if (inne) pola.zakresReczny = null;
}

/** Zwraca zakres pokazywany przez ster operacji kontekstowej: wybór ręczny, a bez niego ślad zaznaczenia. */
export function zakresZadany(pola: PolaStanu): StudioOperationScope {
  if (pola.zakresReczny !== null) return pola.zakresReczny;
  return pola.zakres === null ? StudioOperationScope.Document : StudioOperationScope.Selection;
}

/** Zwraca zakres, który faktycznie pojedzie do rdzenia, różniący się od żądanego tylko przy pustym zaznaczeniu. */
export function zakresSkuteczny(pola: PolaStanu): StudioOperationScope {
  return zakresZadany(pola) === StudioOperationScope.Selection && pola.zakres !== null
    ? StudioOperationScope.Selection
    : StudioOperationScope.Document;
}

/** Wciąga dokument oddany przez rdzeń do stanu modułu jako treść zaakceptowaną i roboczą naraz, od razu. */
export function wchlonDokument(pola: PolaStanu, dokument: StudioDocument): void {
  pola.dokument = dokument;
  pola.robocza = dokument.content ?? '';
  pola.zaakceptowana = pola.robocza;
  pola.zakres = null;
  pola.zakresReczny = null;
}

/** Zdejmuje ze stanu dokument usunięty w rdzeniu; ślad po bycie, którego już nie ma, myliłby operatora. */
export function zapomnijDokument(pola: PolaStanu): void {
  pola.dokument = null;
  pola.robocza = '';
  pola.zaakceptowana = '';
  pola.historia = [];
  pola.propozycja = null;
  pola.zakres = null;
  pola.zakresReczny = null;
  // Odmowa i para porównania dotyczyły dokumentu, którego już nie ma, więc obie schodzą razem z nim.
  pola.odmowaOperacji = null;
  pola.paraPorownania = null;
}

/** Migawka pól dotyczących jednego dokumentu, odkładana przy przełączeniu zakładki okna pracy z dokumentem. */
export interface MigawkaDokumentu {
  dokument: StudioDocument | null;
  robocza: string;
  zaakceptowana: string;
  zakres: Zaznaczenie | null;
  zakresReczny: StudioOperationScope | null;
  propozycja: PropozycjaZmiany | null;
  odmowaOperacji: string | null;
  paraPorownania: ParaPorownania | null;
}

/** Zdejmuje ze stanu modułu migawkę pól dokumentu obecnie czynnego, gotową do odłożenia w jego zakładce. */
export function migawkaPolStanu(pola: PolaStanu): MigawkaDokumentu {
  return {
    dokument: pola.dokument,
    robocza: pola.robocza,
    zaakceptowana: pola.zaakceptowana,
    zakres: pola.zakres,
    zakresReczny: pola.zakresReczny,
    propozycja: pola.propozycja,
    odmowaOperacji: pola.odmowaOperacji,
    paraPorownania: pola.paraPorownania,
  };
}

/** Migawka dokumentu jeszcze niewczytanego, przyjmowana przez nową zakładkę bez własnego dokumentu Studia. */
export function pustaMigawka(): MigawkaDokumentu {
  return {
    dokument: null,
    robocza: '',
    zaakceptowana: '',
    zakres: null,
    zakresReczny: null,
    propozycja: null,
    odmowaOperacji: null,
    paraPorownania: null,
  };
}

/** Wstawia migawkę odłożoną wcześniej na miejsce pól dokumentu, który staje się teraz dokumentem czynnym. */
export function przywrocPolaStanu(pola: PolaStanu, migawka: MigawkaDokumentu): void {
  pola.dokument = migawka.dokument;
  pola.robocza = migawka.robocza;
  pola.zaakceptowana = migawka.zaakceptowana;
  pola.zakres = migawka.zakres;
  pola.zakresReczny = migawka.zakresReczny;
  pola.propozycja = migawka.propozycja;
  pola.odmowaOperacji = migawka.odmowaOperacji;
  pola.paraPorownania = migawka.paraPorownania;
  // Historia dotyczy dokumentu, więc przy przełączeniu zakładki jest nieaktualna do odczytu z rdzenia.
  pola.historia = [];
}

/** Przyjmuje propozycję zmiany: jej treść staje się jednocześnie treścią roboczą i zaakceptowaną dokumentu. */
export function przyjmijPropozycje(pola: PolaStanu): void {
  if (pola.propozycja === null) return;
  if (pola.propozycja.tresc !== '') {
    pola.robocza = pola.propozycja.tresc;
    pola.zaakceptowana = pola.propozycja.tresc;
  }
  pola.propozycja = null;
}
