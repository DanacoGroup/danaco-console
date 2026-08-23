import { PermissionMode } from '../../../shared/contract';

/**
 * Nazwy prezentacyjne wartości używanych wyłącznie przez komplet sterowania.
 *
 * Nazwy roli okna i środowiska wykonania są już nazwane w `etykiety-okna.ts`
 * okna komunikacji — komplet czyta je stamtąd, zamiast zakładać drugi słownik
 * tych samych wartości.
 *
 * Ta sama wartość jest pokazywana w trzech miejscach o trzech różnych
 * szerokościach, więc każda ma trzy postacie w jednym wykazie:
 *   — postać pełna (`nazwa…`) — wiersz listy wyboru w kolumnie sterowania
 *     i podsumowanie szuflady; miejsca jest tam na całe zdanie,
 *   — postać krótka (`nazwaKrotka…`) — uchwyt steru w pasku zlecenia, gdzie
 *     etykieta jest bieżącą wartością, a pasek ma zostać jednym rzędem:
 *     „Ręczny", nie „Ręczny — pytanie przed każdą zmianą",
 *   — opis (`opis…`) — zdanie przy pozycji menu; mówi o skutku wyboru,
 *     nie powtarza nazwy.
 * Postać pełna nie jest osobnym napisem, tylko złożeniem krótkiej z dopiskiem
 * — inaczej pasek i szuflada rozjechałyby się przy pierwszej zmianie słownika.
 */

/** Trzy postacie nazwy jednej wartości nastawy. */
interface PostacieNazwy {
  /** Nazwa krótka — uchwyt steru w pasku zlecenia. */
  krotka: string;
  /** Dopisek postaci pełnej; pusty znaczy „postać pełna równa się krótkiej". */
  dopisek: string;
  /** Zdanie o skutku wyboru — pozycja menu. */
  opis: string;
}

/** Postać pełna: nazwa krótka złożona z dopiskiem. */
function pelna(postacie: PostacieNazwy): string {
  return postacie.dopisek === ''
    ? postacie.krotka
    : `${postacie.krotka} — ${postacie.dopisek}`;
}

const TRYBY_UPRAWNIEN: Record<PermissionMode, PostacieNazwy> = {
  [PermissionMode.Manual]: {
    krotka: 'Ręczny',
    dopisek: 'pytanie przed każdą zmianą',
    opis: 'Model zatrzymuje się przed każdym krokiem zmieniającym cokolwiek i czeka na zgodę.',
  },
  [PermissionMode.AcceptEdits]: {
    krotka: 'Akceptuj zmiany plików',
    dopisek: '',
    opis: 'Zmiany w plikach idą bez pytania; pozostałe kroki nadal wymagają zgody.',
  },
  [PermissionMode.Plan]: {
    krotka: 'Plan',
    dopisek: 'bez zmian w systemie',
    opis: 'Model układa plan i niczego nie wykonuje — ani pliku, ani polecenia powłoki.',
  },
  [PermissionMode.Auto]: {
    krotka: 'Automatyczny',
    dopisek: 'decyduje model',
    opis: 'Model sam rozstrzyga, o co zapytać, a co wykonać bez pytania.',
  },
  [PermissionMode.DontAsk]: {
    krotka: 'Bez pytania',
    dopisek: 'z zachowaniem ograniczeń',
    opis: 'Model nie pyta o zgodę, ale zostaje w granicach ustawionych dla okna.',
  },
  [PermissionMode.BypassPermissions]: {
    krotka: 'Pominięcie kontroli uprawnień',
    dopisek: '',
    opis: 'Kontrola uprawnień nie jest sprawdzana — model wykonuje każdy krok, o który poprosi.',
  },
};

export function nazwaTrybuUprawnien(tryb: PermissionMode): string {
  return pelna(TRYBY_UPRAWNIEN[tryb]);
}

/** Nazwa trybu uprawnień na uchwyt steru — bez dopisku. */
export function nazwaKrotkaTrybuUprawnien(tryb: PermissionMode): string {
  return TRYBY_UPRAWNIEN[tryb].krotka;
}

/** Zdanie o skutku wyboru trybu uprawnień — pozycja menu. */
export function opisTrybuUprawnien(tryb: PermissionMode): string {
  return TRYBY_UPRAWNIEN[tryb].opis;
}

/**
 * Stopnie nakładu rozumowania: od odpowiedzi szybkiej po najstaranniejszą.
 *
 * Etykiety przepisane z kolumny `opcja_ustawienia.etykieta` katalogu rdzenia dla
 * klucza `naklad_rozumowania`. Stopień jest napisem wyliczenia, nie liczbą,
 * a napis pusty jest pełnoprawnym stopniem katalogu.
 */
const NAKLADY: Record<string, PostacieNazwy> = {
  '': {
    krotka: 'Bez wskazania',
    dopisek: 'rozstrzyga kanał modelu',
    opis: 'Okno nie narzuca nakładu — obowiązuje ustawienie kanału modelu.',
  },
  low: {
    krotka: 'Niski',
    dopisek: 'szybciej',
    opis: 'Najkrótszy namysł: odpowiedź przychodzi najszybciej i kosztuje najmniej.',
  },
  medium: {
    krotka: 'Średni',
    dopisek: 'równowaga',
    opis: 'Namysł pośredni — domyślny wybór do pracy bieżącej.',
  },
  high: {
    krotka: 'Wysoki',
    dopisek: 'mądrzej',
    opis: 'Dłuższy namysł przed odpowiedzią; zadania wielokrokowe wychodzą staranniej.',
  },
  xhigh: {
    krotka: 'Bardzo wysoki',
    dopisek: 'mądrzej, dłużej',
    opis: 'Namysł znacznie dłuższy — do zadań, w których błąd kosztuje więcej niż czas.',
  },
  max: {
    krotka: 'Najwyższy',
    dopisek: 'pełny namysł',
    opis: 'Pełny namysł bez oszczędzania czasu — najwolniej i najdrożej.',
  },
};

/**
 * Nazwa stopnia nakładu.
 *
 * Brak wartości schodzi na napis pusty, bo to pełnoprawny stopień katalogu
 * o nazwie „Bez wskazania" — nie na etykietę wymyśloną w tym pliku.
 *
 * Stopień spoza wykazu, ale niepusty, zostaje pokazany dosłownie: rdzeń może
 * znać stopień, którego ten słownik jeszcze nie zna, a jego kod jest prawdą
 * o oknie.
 */
export function nazwaNakladu(stopien: string | undefined | null): string {
  const kod = typeof stopien === 'string' ? stopien : '';
  const postacie = NAKLADY[kod];
  if (postacie !== undefined) return pelna(postacie);
  return kod === '' ? pelna(NAKLADY['']) : kod;
}

/** Nazwa stopnia nakładu na uchwyt steru — bez dopisku. */
export function nazwaKrotkaNakladu(stopien: string): string {
  return NAKLADY[stopien]?.krotka ?? nazwaNakladu(stopien);
}

/** Zdanie o skutku wyboru stopnia nakładu — pozycja menu. */
export function opisNakladu(stopien: string): string {
  return NAKLADY[stopien]?.opis ?? '';
}

/**
 * Nazwa modułu okna — zapasowa, gdy nazwa rdzenia jeszcze nie jest znana.
 *
 * Rdzeń przysyła moduły komendą `module.list` z gotową nazwą (`Module.name`);
 * ta funkcja jest wyłącznie zapasem na kod, dopóki katalog rdzenia się nie
 * naładuje. Oddaje identyfikator bez zmiany, bo moduł spoza wykazu ma zachować
 * swój kod zamiast zniknąć z widoku pod nazwą zmyśloną po stronie klienta.
 */
export function nazwaModulu(identyfikator: string): string {
  return identyfikator;
}
