import {
  ProgressStatus,
  QueueAction,
  QueueStatus,
  SessionStatus,
  WindowRole,
} from '../../../shared/contract';
import { RodzajRelacji, RodzajUtworzenia } from './model-danych';

/**
 * Polskie nazwy stanów i działań pulpitu.
 *
 * Jedna odpowiedzialność: przekład wartości kontraktu na słowo widoczne dla
 * operatora. Widok nie zna literałów kontraktu — pyta o etykietę tutaj.
 *
 * Rekordy są kompletne (`Record<T, string>`), więc dopisanie wartości do
 * kontraktu przerywa kompilację tego pliku. To zamierzone: nowy stan ma dostać
 * polską nazwę, a nie wypaść z widoku po cichu.
 */

/**
 * Etykieta miary, której kontrakt dziś nie niesie. Widok wypisuje
 * ją zamiast liczby — pulpit nigdy nie pokazuje wartości wymyślonej.
 */
export const ETYKIETA_BRAKU_ZRODLA = 'brak źródła danych w kontrakcie';

/** Stan sesji. */
export const ETYKIETA_STANU_SESJI: Readonly<Record<SessionStatus, string>> = {
  [SessionStatus.Active]: 'czynna',
  [SessionStatus.Paused]: 'wstrzymana',
  [SessionStatus.Finished]: 'zakończona',
  [SessionStatus.Archived]: 'zarchiwizowana',
};

/** Rola okna w pętli koordynator–wykonawca. */
export const ETYKIETA_ROLI_OKNA: Readonly<Record<WindowRole, string>> = {
  [WindowRole.Coordinator]: 'koordynator',
  [WindowRole.Executor]: 'wykonawca',
  [WindowRole.Standalone]: 'samodzielne',
};

/** Stan kolejki. */
export const ETYKIETA_STANU_KOLEJKI: Readonly<Record<QueueStatus, string>> = {
  [QueueStatus.Idle]: 'gotowa',
  [QueueStatus.Running]: 'w biegu',
  [QueueStatus.Paused]: 'wstrzymana',
  [QueueStatus.Stopped]: 'zatrzymana',
  [QueueStatus.Done]: 'wyczerpana',
};

/** Stan procesu w telemetrii postępu. */
export const ETYKIETA_STANU_PROCESU: Readonly<Record<ProgressStatus, string>> = {
  [ProgressStatus.Pending]: 'oczekuje',
  [ProgressStatus.Running]: 'w biegu',
  [ProgressStatus.Paused]: 'wstrzymany',
  [ProgressStatus.Stopped]: 'zatrzymany',
  [ProgressStatus.Done]: 'zakończony',
  [ProgressStatus.Failed]: 'zakończony błędem',
};

/** Działanie transportu kolejki. */
export const ETYKIETA_DZIALANIA_KOLEJKI: Readonly<Record<QueueAction, string>> = {
  [QueueAction.Start]: 'Uruchom',
  [QueueAction.Pause]: 'Wstrzymaj',
  [QueueAction.Resume]: 'Wznów',
  [QueueAction.Stop]: 'Zatrzymaj',
  [QueueAction.Retry]: 'Ponów',
  [QueueAction.Clear]: 'Opróżnij',
  [QueueAction.Enqueue]: 'Wstaw do kolejki',
  [QueueAction.Dequeue]: 'Zdejmij z kolejki',
  [QueueAction.Delay]: 'Odłóż w czasie',
  [QueueAction.Split]: 'Rozdziel',
  [QueueAction.Merge]: 'Scal',
  [QueueAction.Route]: 'Skieruj',
  [QueueAction.Branch]: 'Rozgałęź',
  [QueueAction.Condition]: 'Warunek',
};

/** Rodzaj powiązania sesji w pasie „Relacje". */
export const ETYKIETA_RODZAJU_RELACJI: Readonly<Record<RodzajRelacji, string>> = {
  [RodzajRelacji.Wymiana]: 'wymiana dwustronna',
  [RodzajRelacji.Przekazanie]: 'przekazanie jednostronne',
};

/** Znak wiodący relacji: dwustronny lub jednostronny. */
export const STRZALKA_RELACJI: Readonly<Record<RodzajRelacji, string>> = {
  [RodzajRelacji.Wymiana]: '⇄',
  [RodzajRelacji.Przekazanie]: '→',
};

/** Podpis kafla w rzędzie „Utwórz". */
export const ETYKIETA_UTWORZENIA: Readonly<Record<RodzajUtworzenia, string>> = {
  [RodzajUtworzenia.Sesja]: 'Nowa sesja',
  [RodzajUtworzenia.Projekt]: 'Załóż projekt',
  [RodzajUtworzenia.Agent]: 'Skonfiguruj agenta',
  [RodzajUtworzenia.Automatyka]: 'Zbuduj automatykę',
  [RodzajUtworzenia.Kolejka]: 'Ustaw kolejkę',
  [RodzajUtworzenia.Zespol]: 'Złóż zespół',
};
