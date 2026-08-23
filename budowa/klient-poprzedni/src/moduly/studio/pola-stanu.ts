import {
  StudioOperationScope,
  type StudioDocument,
  type StudioVersion,
} from '../../../../shared/contract';

/** Zaznaczony fragment dokumentu — przedmiot operacji kontekstowej. */
export interface Zaznaczenie {
  poczatek: number;
  koniec: number;
}

/** Wynik operacji AI czekający na przyjęcie albo odrzucenie. */
export interface PropozycjaZmiany {
  /** Propozycja po stronie rdzenia; zasila `studio.diff.compare`. */
  idPropozycji: string;
  /** Treść wyniku operacji, o ile rdzeń ją oddał. */
  tresc: string;
  /** Operacja, która propozycję wywołała — do nazwania jej w panelu różnicy. */
  idAkcji: string;
}

/**
 * Pola stanu modułu Studio wraz z przejściami, które je zmieniają.
 *
 * Wydzielone od wytwórni stanu, bo to dwie odpowiedzialności: tutaj kształt
 * stanu i jego przejścia, tam subskrypcja rdzenia i powiadamianie widoków.
 * Podział pozwala sprawdzić przejścia bez stawiania kanału.
 *
 * Treść robocza i treść zaakceptowana to dwa różne pola. Edytor pracuje na
 * pierwszej, podgląd czyta drugą — podgląd odświeża się po zmianach
 * zaakceptowanych, a nie po każdym naciśnięciu klawisza.
 */
export interface PolaStanu {
  okno: string;
  dokument: StudioDocument | null;
  robocza: string;
  zakres: Zaznaczenie | null;
  /**
   * Wybór ręczny zakresu operacji; `null` znaczy „idź za zaznaczeniem".
   *
   * Nastawa jest jedna i mieszka tutaj, a `zakresZadany` (co pokazuje ster)
   * i `zakresSkuteczny` (co poleci do rdzenia) są z niej wyliczane, nie
   * przechowywane osobno. Trzymanie jej w `value` listy wyboru Tools Panelu
   * rozjeżdżałoby wyświetlaną wartość z wysyłaną: odświeżenie panelu nadpisuje
   * listę wedle zaznaczenia w edytorze.
   */
  zakresReczny: StudioOperationScope | null;
  historia: readonly StudioVersion[];
  propozycja: PropozycjaZmiany | null;
  zaakceptowana: string;
  /**
   * Powód, dla którego ostatnia operacja kontekstowa nie dała propozycji;
   * `null`, gdy żadna operacja nie odmówiła od czasu ostatniego wyniku.
   *
   * Odmowa i pustka to dwie różne rzeczy, a samo `propozycja === null` niesie
   * obie naraz: „nikt jeszcze nie uruchomił operacji" i „operacja wróciła
   * odmową". Odmowę widzi pas stanu Tools Panelu, ale kanwa tekstowa stoi obok
   * i musi wiedzieć to samo — stąd osobne pole.
   */
  odmowaOperacji: string | null;
  /**
   * Para wersji wskazana do porównania w Diff/Grep Panelu; `null`, gdy żadnej
   * nie wskazano.
   *
   * Pole stoi w stanie modułu, a nie w polach Diff/Grep Panelu, bo wskazuje je
   * INNE okno: opracowanie każe przyciskowi „Porównaj" przy wersji otworzyć
   * panel z wypełnioną parą. Sięganie z repozytorium do kontrolek cudzego okna
   * byłoby drugą drogą do tej samej nastawy.
   */
  paraPorownania: ParaPorownania | null;
}

/** Dwie strony porównania wskazane w Session Repository. */
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

/**
 * Przyjmuje zaznaczenie edytora i rozstrzyga los wyboru ręcznego.
 *
 * Zaznaczenie nowe kasuje wybór ręczny, zaznaczenie to samo — nie. Wybór ręczny
 * ma pierwszeństwo do następnego zaznaczenia, a „następne" znaczy inne: edytor
 * woła tę drogę przy każdym `select`, `keyup` i `mouseup`, więc kasowanie na
 * każdym wywołaniu zdejmowałoby wybór ręczny przy pierwszym kliknięciu w treść.
 */
export function ustawZaznaczenieStanu(pola: PolaStanu, zakres: Zaznaczenie | null): void {
  const inne =
    zakres !== null &&
    (pola.zakres === null ||
      pola.zakres.poczatek !== zakres.poczatek ||
      pola.zakres.koniec !== zakres.koniec);
  pola.zakres = zakres;
  if (inne) pola.zakresReczny = null;
}

/**
 * Zakres pokazywany przez ster: wybór ręczny, a bez niego ślad zaznaczenia.
 *
 * To jest wartość listy wyboru w Tools Panelu — i wyłącznie ona. Panel jej nie
 * przechowuje u siebie, bo dwa miejsca przechowujące jedną nastawę rozjeżdżają
 * się przy pierwszym odświeżeniu wywołanym z zewnątrz.
 */
export function zakresZadany(pola: PolaStanu): StudioOperationScope {
  if (pola.zakresReczny !== null) return pola.zakresReczny;
  return pola.zakres === null ? StudioOperationScope.Document : StudioOperationScope.Selection;
}

/**
 * Zakres, który pojedzie do rdzenia.
 *
 * Różni się od żądanego dokładnie w jednym przypadku: wybrano „zaznaczenie",
 * a w edytorze nic nie jest zaznaczone. Wybór nie jest wtedy blokowany
 * — schodzi na cały dokument, a okna mówią o tym wprost, zamiast
 * milcząco wysyłać `selection` bez granic i prosić się o odmowę walidacji.
 */
export function zakresSkuteczny(pola: PolaStanu): StudioOperationScope {
  return zakresZadany(pola) === StudioOperationScope.Selection && pola.zakres !== null
    ? StudioOperationScope.Selection
    : StudioOperationScope.Document;
}

/**
 * Wciąga dokument oddany przez rdzeń.
 *
 * Treść z rdzenia staje się i roboczą, i zaakceptowaną: dokument właśnie
 * wczytany albo zapisany jest stanem zaakceptowanym, więc podgląd ma go pokazać
 * od razu. Zaznaczenie znika, bo odnosiło się do treści poprzedniej — a razem
 * z nim wybór ręczny zakresu, bo dotyczył tamtego zaznaczenia i tamtej treści.
 */
export function wchlonDokument(pola: PolaStanu, dokument: StudioDocument): void {
  pola.dokument = dokument;
  pola.robocza = dokument.content ?? '';
  pola.zaakceptowana = pola.robocza;
  pola.zakres = null;
  pola.zakresReczny = null;
}

/** Zdejmuje dokument usunięty w rdzeniu; ślad po bycie, którego nie ma, myli. */
export function zapomnijDokument(pola: PolaStanu): void {
  pola.dokument = null;
  pola.robocza = '';
  pola.zaakceptowana = '';
  pola.historia = [];
  pola.propozycja = null;
  pola.zakres = null;
  pola.zakresReczny = null;
  // Odmowa dotyczyła operacji na dokumencie, którego już nie ma — trzymanie jej
  // dalej opisywałoby byt nieistniejący. Para porównania schodzi z tego samego
  // powodu: wskazywała wersje dokumentu, który zniknął.
  pola.odmowaOperacji = null;
  pola.paraPorownania = null;
}

/**
 * Migawka pól dotyczących JEDNEGO dokumentu.
 *
 * Okno pracy prowadzi dwa dokumenty naraz w zakładkach, a stan modułu ma jeden
 * dokument czynny — bo Tools Panel, Session Repository i Ingest/OCR Panel patrzą
 * właśnie na czynny. Zakładka nieczynna trzyma więc swoje pola w migawce
 * i wraca z nimi przy przełączeniu. Okno modułu jest wciąż jedno i jeden
 * dokument jest w nim czynny — zmienia się to, KTÓRY.
 *
 * Migawka nie niesie okna osadzenia ani historii wersji: okno jest wspólne
 * całemu modułowi, a historia dotyczy dokumentu i doczyta się przy przełączeniu.
 */
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

/** Zdejmuje migawkę pól dokumentu czynnego. */
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

/** Migawka dokumentu jeszcze niewczytanego — zakładka nowa. */
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

/** Wstawia migawkę na miejsce pól dokumentu czynnego. */
export function przywrocPolaStanu(pola: PolaStanu, migawka: MigawkaDokumentu): void {
  pola.dokument = migawka.dokument;
  pola.robocza = migawka.robocza;
  pola.zaakceptowana = migawka.zaakceptowana;
  pola.zakres = migawka.zakres;
  pola.zakresReczny = migawka.zakresReczny;
  pola.propozycja = migawka.propozycja;
  pola.odmowaOperacji = migawka.odmowaOperacji;
  pola.paraPorownania = migawka.paraPorownania;
  // Historia dotyczy dokumentu, więc przy przełączeniu zakładki jest nieaktualna
  // do czasu odczytu. Zostawienie jej pokazywałoby wersje cudzego dokumentu.
  pola.historia = [];
}

/**
 * Przyjmuje propozycję: jej treść staje się roboczą i zaakceptowaną.
 *
 * Propozycja bez treści (rdzeń oddał samo `proposalId`) też zostaje przyjęta —
 * decyzja Operatora jest wiążąca niezależnie od tego, czy rdzeń dołączył wynik
 * wprost, czy odesłał do porównania w Diff/Grep Panelu.
 */
export function przyjmijPropozycje(pola: PolaStanu): void {
  if (pola.propozycja === null) return;
  if (pola.propozycja.tresc !== '') {
    pola.robocza = pola.propozycja.tresc;
    pola.zaakceptowana = pola.propozycja.tresc;
  }
  pola.propozycja = null;
}
