import { WindowRole, type Window } from '../../../../shared/contract';

/** Obsada czterech okien ról modułu MultitaskingAI złożona z pól kontraktu. */

/**
 * Kody okien operacyjnych rejestru rdzenia zapisane bez przedrostka modułu. Kod
 * z przedrostkiem nie trafiłby w żaden wiersz rejestru, więc nawigacja po kodzie
 * okna nie miałaby dokąd skoczyć.
 */
export const KOD_KOORDYNATORA = 'coordinator-chat';
export const KOD_ANALITYKA = 'results-analyzer';

/**
 * Ilu wykonawców niesie wykaz okien operacyjnych: rejestr ma dwa wystąpienia
 * okna wykonawcy, więc obsada przycina zbiór przypiętych okien do tej liczby.
 */
export const LICZBA_WYKONAWCOW = 2;

/**
 * Obsada ról odczytana z okien sesji: koordynator, przypięci do niego wykonawcy,
 * wskazany analityk oraz okna wykonawcze należące do innego koordynatora.
 */
export interface Obsada {
  /** Okno koordynatora; puste, gdy sesja nie ma jeszcze zarządcy. */
  koordynator: Window | null;
  /** Wykonawcy przypięci do tego koordynatora, w kolejności założenia. */
  wykonawcy: readonly Window[];
  /** Okno analityka wskazane na poziomie sesji; puste, gdy nie wskazano. */
  analityk: Window | null;
  /** Okna wykonawcze przypięte do innego koordynatora — obce tej scenie. */
  obce: readonly Window[];
}

/**
 * Pusta obsada — stan przed pierwszym odczytem okien sesji oraz po odmowie
 * rdzenia. Widok dostaje wtedy komplet pól bez wartości, a nie brak obsady.
 */
export const OBSADA_PUSTA: Obsada = {
  koordynator: null,
  wykonawcy: [],
  analityk: null,
  obce: [],
};

/**
 * Składa obsadę z okien sesji uporządkowanych po czasie założenia.
 * Koordynatorem zostaje pierwsze okno roli koordynatora; druga taka rola jest
 * osobną sceną, a nie konkurentem, bo moduł pokazuje jedną parę naraz.
 */
export function zlozObsade(okna: readonly Window[], idAnalityka: string): Obsada {
  const wedlugCzasu = [...okna].sort((a, b) => a.createdAt - b.createdAt);
  const koordynator = wedlugCzasu.find((okno) => okno.windowRole === WindowRole.Coordinator) ?? null;
  const wykonawcy = wedlugCzasu.filter((okno) => okno.windowRole === WindowRole.Executor);
  const moi =
    koordynator === null
      ? []
      : wykonawcy.filter((okno) => okno.coordinatorWindowId === koordynator.id);
  const analityk =
    wedlugCzasu.find(
      (okno) => okno.id === idAnalityka && okno.windowRole === WindowRole.Standalone,
    ) ?? null;

  return {
    koordynator,
    wykonawcy: moi.slice(0, LICZBA_WYKONAWCOW),
    analityk,
    obce: wykonawcy.filter((okno) => !moi.includes(okno)),
  };
}

/**
 * Czy okno o podanym oznaczeniu należy do wykonawców tej obsady. Rozstrzyga
 * zbiór wykonawców przypiętych do koordynatora, a nie sama rola okna.
 */
export function czyWykonawca(obsada: Obsada, idOkna: string): boolean {
  return obsada.wykonawcy.some((okno) => okno.id === idOkna);
}
