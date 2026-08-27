import {
  KnowledgeScope,
  ResearchSourceKind,
  type KnowledgeHit,
  type LibraryFile,
  type ResearchCredibility,
} from '../../../../shared/contract';
import type { ZlecenieZrodla } from './zlecenia-badania';

/**
 * Pozycja wyniku Discovery Panel jest jednym kształtem dla dwóch różnych odpowiedzi kontraktu:
 * fragmentów wiedzy oraz zasobów repozytorium, sprowadzonych tu do jednego bytu.
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

/** Fragmenty wiedzy Operatora przełożone na pozycje wyniku wyszukiwania, gotowe do wyświetlenia w wykazie Discovery Panel. */
export function zHitowWiedzy(trafienia: readonly KnowledgeHit[]): WynikOdkrycia[] {
  return trafienia.map((trafienie, kolejnosc) => ({
    klucz: trafienie.sourceId ?? `${trafienie.scope}-${String(kolejnosc)}`,
    tytul: pierwszeZdanie(trafienie.text),
    metadane: metadaneTrafienia(trafienie),
    fragment: trafienie.text,
    // Fragment znaleziony po znaczeniu wchodzi jako notatka; wiązanie z repozytorium niesie osobne pole.
    rodzaj: ResearchSourceKind.Note,
    pochodzenie: trafienie.source,
    idPlikuRepozytorium: idPlikuTrafienia(trafienie),
  }));
}

/** Zasoby repozytorium Library przełożone na pozycje wyniku wyszukiwania, gotowe do wyświetlenia w wykazie Discovery Panel. */
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
 * Zlecenie katalogowania złożone z pozycji wyniku; adres sieciowy zostaje pusty, bo pozycje
 * wyszukiwania nie niosą takiego pola.
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

/** Metadane trafienia wiedzy: zakres wyszukiwania, źródło pochodzenia oraz trafność oddana przez rdzeń, złożone w jeden opis. */
function metadaneTrafienia(trafienie: KnowledgeHit): string {
  const czesci = [
    `zakres: ${trafienie.scope}`,
    `pochodzenie: ${trafienie.source}`,
    trafienie.score === undefined ? '' : `trafność: ${String(trafienie.score)}/100`,
  ];
  return czesci.filter((czesc) => czesc !== '').join(' · ');
}

/** Metadane zasobu repozytorium: ścieżka, rodzaj treści, rozmiar oraz etykiety, złożone w jeden wiersz opisowy. */
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
 * Dokument repozytorium wskazany przez trafienie, ustalany wyłącznie dla zakresu biblioteki, bo
 * poza nim identyfikator wskazuje inny byt.
 */
function idPlikuTrafienia(trafienie: KnowledgeHit): string {
  const zBiblioteki = trafienie.scope === KnowledgeScope.Library;
  return zBiblioteki ? (trafienie.sourceId ?? '') : '';
}

/** Pierwsze zdanie fragmentu jako nazwa pozycji wyniku; fragment bez kropki idzie w całości, skrócony po stu dwudziestu znakach. */
function pierwszeZdanie(tresc: string): string {
  const zwiniete = tresc.replace(/\s+/g, ' ').trim();
  if (zwiniete === '') return 'fragment bez treści';
  const koniec = zwiniete.indexOf('. ');
  const zdanie = koniec === -1 ? zwiniete : zwiniete.slice(0, koniec + 1);
  return zdanie.length <= 120 ? zdanie : `${zdanie.slice(0, 119)}…`;
}
