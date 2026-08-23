import { WindowRole, type Window } from '../../../../shared/contract';

/**
 * Obsada czterech okien ról MultitaskingAI zbudowana z pól kontraktu.
 *
 * Rola pochodzi z `WindowRole`, nie z domysłu. Koordynatorem jest okno
 * `coordinator`, wykonawcą okno `executor`, którego `coordinatorWindowId`
 * wskazuje tego koordynatora — to jest cała definicja pętli. Znakowanie ról po
 * tytule albo po kolejności rozjechałoby rdzeń z widokiem przy pierwszym
 * przepięciu wykonawcy pod innego koordynatora.
 *
 * Analityk stoi poza pętlą. Results Analyzer nie jest wykonawcą: nie kończy
 * tury, która miałaby wybudzić koordynatora, więc jego rolą kontraktową jest
 * `standalone`. Odróżnia go wskazanie zapisane na poziomie sesji — kontrakt nie
 * ma czwartej wartości `WindowRole`, a wymyślanie jej po stronie klienta
 * rozjechałoby wykaz ról z bazą.
 *
 * Kody okien są bezmodułowe: rejestr okien operacyjnych niesie `coordinator-chat`,
 * `executor-chat-1`, `executor-chat-2` i `results-analyzer`, bez przedrostka
 * modułu. Przynależność do modułu niesie osobna macierz, w której czterech okien
 * ról nie ma — liczą się tam jako okna pozamodułowe. Kod z przedrostkiem
 * wypisany w `data-okno` nie trafiłby w żaden wiersz rejestru i nawigacja po
 * kodzie okna nie miałaby w co skoczyć.
 */

/** Kody okien operacyjnych rejestru rdzenia — bez przedrostka modułu. */
export const KOD_KOORDYNATORA = 'coordinator-chat';
export const KOD_ANALITYKA = 'results-analyzer';

/** Ilu wykonawców niesie wykaz okien: dwa wystąpienia Executor Chat. */
export const LICZBA_WYKONAWCOW = 2;

/** Obsada ról odczytana z okien sesji. */
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

/** Pusta obsada — stan przed odczytem i po odmowie rdzenia. */
export const OBSADA_PUSTA: Obsada = {
  koordynator: null,
  wykonawcy: [],
  analityk: null,
  obce: [],
};

/**
 * Składa obsadę z okien sesji.
 *
 * Koordynatorem zostaje pierwsze okno roli `coordinator` w kolejności założenia;
 * druga taka rola w sesji jest osobną sceną, a nie konkurentem — moduł pokazuje
 * jedną parę naraz, tak jak pas relacji układu okien równoległych.
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

/** Czy okno należy do wykonawców tej obsady. */
export function czyWykonawca(obsada: Obsada, idOkna: string): boolean {
  return obsada.wykonawcy.some((okno) => okno.id === idOkna);
}
