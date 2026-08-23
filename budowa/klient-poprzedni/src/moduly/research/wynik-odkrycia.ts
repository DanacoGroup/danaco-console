import {
  KnowledgeScope,
  ResearchSourceKind,
  type KnowledgeHit,
  type LibraryFile,
  type ResearchCredibility,
} from '../../../../shared/contract';
import type { ZlecenieZrodla } from './zlecenia-badania';

/**
 * Pozycja wyniku Discovery Panel — jeden kształt dla dwóch różnych odpowiedzi
 * kontraktu.
 *
 * `knowledge.search` oddaje fragmenty treści (`KnowledgeHit`), a
 * `library.file.search` — zasoby repozytorium (`LibraryFile`). Panel pokazuje
 * jedną listę wyników i wstawia z niej źródła jedną komendą, więc obie
 * odpowiedzi sprowadza do jednego bytu tutaj, a nie w widoku: widok nie ma
 * rozstrzygać, z której komendy pochodzi wiersz, który rysuje.
 *
 * Przeniesienie do badania to `research.source.add`. Kształt zlecenia powstaje
 * z pozycji, więc reguła „fragment wiedzy wchodzi jako notatka, zasób
 * repozytorium jako dokument" stoi w jednym miejscu.
 */
export interface WynikOdkrycia {
  /** Klucz pozycji w obrębie jednego wyniku wyszukiwania. */
  klucz: string;
  /** Nazwa pozycji widoczna w wykazie. */
  tytul: string;
  /** Metadane wiersza: skąd pozycja pochodzi i czym jest. */
  metadane: string;
  /** Fragment treści, gdy odpowiedź go niesie; pusty znaczy brak fragmentu. */
  fragment: string;
  /** Rodzaj źródła, jakim pozycja stanie się po przeniesieniu do badania. */
  rodzaj: ResearchSourceKind;
  /** Pochodzenie źródła — pole `origin` komendy katalogowania. */
  pochodzenie: string;
  /** Dokument repozytorium, gdy pozycja go wskazuje; pusty znaczy brak. */
  idPlikuRepozytorium: string;
}

/** Fragmenty wiedzy Operatora jako pozycje wyniku. */
export function zHitowWiedzy(trafienia: readonly KnowledgeHit[]): WynikOdkrycia[] {
  return trafienia.map((trafienie, kolejnosc) => ({
    klucz: trafienie.sourceId ?? `${trafienie.scope}-${String(kolejnosc)}`,
    tytul: pierwszeZdanie(trafienie.text),
    metadane: metadaneTrafienia(trafienie),
    fragment: trafienie.text,
    // Fragment znaleziony po znaczeniu jest wypisem z korpusu, a nie plikiem:
    // wchodzi do katalogu jako notatka. Gdy trafienie wskazuje zasób
    // repozytorium, wiązanie niesie osobne pole, nie rodzaj źródła.
    rodzaj: ResearchSourceKind.Note,
    pochodzenie: trafienie.source,
    idPlikuRepozytorium: idPlikuTrafienia(trafienie),
  }));
}

/** Zasoby repozytorium Library jako pozycje wyniku. */
export function zPlikowRepozytorium(pliki: readonly LibraryFile[]): WynikOdkrycia[] {
  return pliki.map((plik) => ({
    klucz: plik.id,
    tytul: plik.name,
    metadane: metadanePliku(plik),
    fragment: '',
    rodzaj: ResearchSourceKind.Document,
    pochodzenie: plik.path ?? plik.name,
    idPlikuRepozytorium: plik.id,
  }));
}

/**
 * Zlecenie katalogowania złożone z pozycji wyniku.
 *
 * Adresu pozycje wyszukiwania nie niosą — ani `KnowledgeHit`, ani `LibraryFile`
 * nie mają pola z adresem sieciowym — więc pole `url` zostaje puste zamiast być
 * dopowiedziane ze ścieżki repozytorium, która adresem nie jest.
 */
export function zlecenieZWyniku(
  pozycja: WynikOdkrycia,
  idOkna: string,
  wiarygodnosc: ResearchCredibility,
): ZlecenieZrodla {
  return {
    idOkna,
    tytul: pozycja.tytul,
    rodzaj: pozycja.rodzaj,
    adres: '',
    pochodzenie: pozycja.pochodzenie,
    wiarygodnosc,
    idPlikuRepozytorium: pozycja.idPlikuRepozytorium,
  };
}

/** Metadane trafienia wiedzy: zakres, źródło i trafność oddana przez rdzeń. */
function metadaneTrafienia(trafienie: KnowledgeHit): string {
  const czesci = [
    `zakres: ${trafienie.scope}`,
    `pochodzenie: ${trafienie.source}`,
    trafienie.score === undefined ? '' : `trafność: ${String(trafienie.score)}/100`,
  ];
  return czesci.filter((czesc) => czesc !== '').join(' · ');
}

/** Metadane zasobu repozytorium: ścieżka, rodzaj treści, rozmiar i etykiety. */
function metadanePliku(plik: LibraryFile): string {
  const czesci = [
    plik.path === undefined ? '' : plik.path,
    plik.mimeType === undefined ? '' : plik.mimeType,
    plik.sizeBytes === undefined ? '' : `${String(plik.sizeBytes)} B`,
    (plik.tags ?? []).length === 0 ? '' : `etykiety: ${(plik.tags ?? []).join(', ')}`,
  ];
  return czesci.filter((czesc) => czesc !== '').join(' · ');
}

/**
 * Dokument repozytorium wskazany przez trafienie.
 *
 * `KnowledgeHit.sourceId` jest identyfikatorem tego, z czego fragment pochodzi,
 * ale pochodzić może z pliku biblioteki, z wiadomości albo z pliku przestrzeni
 * roboczej — zakres to rozstrzyga. Wiązanie z dokumentem repozytorium zakładamy
 * wyłącznie dla zakresu biblioteki; poza nim identyfikator wskazuje byt, którego
 * `research.source.add` w polu `libraryFileId` nie przyjmie.
 */
function idPlikuTrafienia(trafienie: KnowledgeHit): string {
  const zBiblioteki = trafienie.scope === KnowledgeScope.Library;
  return zBiblioteki ? (trafienie.sourceId ?? '') : '';
}

/** Pierwsze zdanie fragmentu jako nazwa pozycji; fragment bez kropki idzie w całości. */
function pierwszeZdanie(tresc: string): string {
  const zwiniete = tresc.replace(/\s+/g, ' ').trim();
  if (zwiniete === '') return 'fragment bez treści';
  const koniec = zwiniete.indexOf('. ');
  const zdanie = koniec === -1 ? zwiniete : zwiniete.slice(0, koniec + 1);
  return zdanie.length <= 120 ? zdanie : `${zdanie.slice(0, 119)}…`;
}
