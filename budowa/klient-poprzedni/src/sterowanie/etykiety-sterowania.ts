import { PermissionMode } from '../../../shared/contract';

// Nazwy prezentacyjne wartości sterowania: postać pełna, krótka i opis budowane z jednego słownika.

/** Trzy postacie nazwy jednej wartości nastawy: krótka na uchwyt steru, pełna ze złożonym dopiskiem i opis skutku wyboru na pozycję menu. */
interface PostacieNazwy {
  /** Nazwa krótka — uchwyt steru w pasku zlecenia. */
  krotka: string;
  /** Dopisek postaci pełnej; pusty znaczy „postać pełna równa się krótkiej". */
  dopisek: string;
  /** Zdanie o skutku wyboru — pozycja menu. */
  opis: string;
}

/** Postać pełna nazwy: krótka złożona z dopiskiem myślnikiem, albo sama krótka, gdy dopisek pozostaje pusty. */
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

/** Nazwa trybu uprawnień na uchwyt steru w pasku zlecenia, bez dopisku, bo pasek ma zostać jednym rzędem. */
export function nazwaKrotkaTrybuUprawnien(tryb: PermissionMode): string {
  return TRYBY_UPRAWNIEN[tryb].krotka;
}

/** Zdanie o skutku wyboru trybu uprawnień, pokazywane jako pozycja menu zamiast powtórzenia nazwy trybu. */
export function opisTrybuUprawnien(tryb: PermissionMode): string {
  return TRYBY_UPRAWNIEN[tryb].opis;
}

/** Stopnie nakładu rozumowania, od odpowiedzi szybkiej po najstaranniejszą, jako napisy wyliczenia katalogu rdzenia. */
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

/** Nazwa stopnia nakładu: brak wartości schodzi na napis pusty jako pełnoprawny stopień, a stopień spoza wykazu pokazuje się dosłownie. */
export function nazwaNakladu(stopien: string | undefined | null): string {
  const kod = typeof stopien === 'string' ? stopien : '';
  const postacie = NAKLADY[kod];
  if (postacie !== undefined) return pelna(postacie);
  return kod === '' ? pelna(NAKLADY['']) : kod;
}

/** Nazwa stopnia nakładu na uchwyt steru w pasku zlecenia, bez dopisku, tak samo jak dla trybu uprawnień. */
export function nazwaKrotkaNakladu(stopien: string): string {
  return NAKLADY[stopien]?.krotka ?? nazwaNakladu(stopien);
}

/** Zdanie o skutku wyboru stopnia nakładu, pokazywane jako pozycja menu, tak samo jak dla trybu uprawnień. */
export function opisNakladu(stopien: string): string {
  return NAKLADY[stopien]?.opis ?? '';
}

/** Nazwa modułu okna, zapasowa do chwili, aż rdzeń przyśle nazwę własną; do tej pory oddaje sam identyfikator bez zmiany. */
export function nazwaModulu(identyfikator: string): string {
  return identyfikator;
}
