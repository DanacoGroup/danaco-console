/**
 * Przekład wartości kontraktu na polskie nazwy stanów i działań pulpitu. Plik
 * trzyma zupełne rekordy etykiet stanu sesji, roli okna, stanu kolejki, stanu
 * procesu, działania kolejki oraz rodzaju relacji i jej utworzenia.
 */
import {
  ProgressStatus,
  QueueAction,
  QueueStatus,
  SessionStatus,
  WindowRole,
} from '../../../shared/contract';
import { RodzajRelacji, RodzajUtworzenia } from './model-danych';

/**
 * Etykieta miary, której kontrakt dziś nie niesie. Widok wypisuje
 * ją zamiast liczby — pulpit nigdy nie pokazuje wartości wymyślonej.
 */
export const ETYKIETA_BRAKU_ZRODLA = 'brak źródła danych w kontrakcie';

/**
 * Polska nazwa stanu sesji dla każdej wartości kontraktu: czynna, wstrzymana,
 * zakończona oraz zarchiwizowana. Rekord jest zupełny, więc nowa wartość
 * w kontrakcie przerywa kompilację zamiast wypaść z widoku.
 */
export const ETYKIETA_STANU_SESJI: Readonly<Record<SessionStatus, string>> = {
  [SessionStatus.Active]: 'czynna',
  [SessionStatus.Paused]: 'wstrzymana',
  [SessionStatus.Finished]: 'zakończona',
  [SessionStatus.Archived]: 'zarchiwizowana',
};

/**
 * Polska nazwa roli okna w pętli koordynator–wykonawca. Kontrakt niesie trzy
 * role: koordynatora, wykonawcę oraz okno samodzielne, czyli stojące poza tą
 * pętlą.
 */
export const ETYKIETA_ROLI_OKNA: Readonly<Record<WindowRole, string>> = {
  [WindowRole.Coordinator]: 'koordynator',
  [WindowRole.Executor]: 'wykonawca',
  [WindowRole.Standalone]: 'samodzielne',
};

/**
 * Polska nazwa stanu kolejki zleceń dla pięciu wartości kontraktu: gotowa,
 * w biegu, wstrzymana, zatrzymana oraz wyczerpana. Widok nie zna literałów
 * kontraktu i pyta o nazwę tutaj.
 */
export const ETYKIETA_STANU_KOLEJKI: Readonly<Record<QueueStatus, string>> = {
  [QueueStatus.Idle]: 'gotowa',
  [QueueStatus.Running]: 'w biegu',
  [QueueStatus.Paused]: 'wstrzymana',
  [QueueStatus.Stopped]: 'zatrzymana',
  [QueueStatus.Done]: 'wyczerpana',
};

/**
 * Polska nazwa stanu procesu w telemetrii postępu. Wykaz odróżnia zakończenie
 * pomyślne od zakończenia błędem, ponieważ pulpit pokazuje je jako dwa różne
 * stany procesu.
 */
export const ETYKIETA_STANU_PROCESU: Readonly<Record<ProgressStatus, string>> = {
  [ProgressStatus.Pending]: 'oczekuje',
  [ProgressStatus.Running]: 'w biegu',
  [ProgressStatus.Paused]: 'wstrzymany',
  [ProgressStatus.Stopped]: 'zatrzymany',
  [ProgressStatus.Done]: 'zakończony',
  [ProgressStatus.Failed]: 'zakończony błędem',
};

/**
 * Podpis przycisku działania transportu kolejki. Podpisy zapisane są wielką
 * literą początkową, ponieważ trafiają wprost na kontrolkę, a nie w zdanie
 * opisujące stan.
 */
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

/**
 * Nazwa rodzaju powiązania sesji w pasie relacji. Wymiana jest dwustronna,
 * przekazanie jednostronne, a kierunek niesie osobno znak wiodący stojący obok
 * nazwy.
 */
export const ETYKIETA_RODZAJU_RELACJI: Readonly<Record<RodzajRelacji, string>> = {
  [RodzajRelacji.Wymiana]: 'wymiana dwustronna',
  [RodzajRelacji.Przekazanie]: 'przekazanie jednostronne',
};

/**
 * Znak wiodący relacji: strzałka dwustronna dla wymiany i jednostronna dla
 * przekazania. Znak towarzyszy nazwie rodzaju i sam jej nie zastępuje.
 */
export const STRZALKA_RELACJI: Readonly<Record<RodzajRelacji, string>> = {
  [RodzajRelacji.Wymiana]: '⇄',
  [RodzajRelacji.Przekazanie]: '→',
};

/**
 * Podpis kafla w rzędzie tworzenia. Każdy rodzaj utworzenia niesie własne zdanie
 * w trybie rozkazującym, ponieważ kafel jest wezwaniem do czynności, a nie nazwą
 * stanu.
 */
export const ETYKIETA_UTWORZENIA: Readonly<Record<RodzajUtworzenia, string>> = {
  [RodzajUtworzenia.Sesja]: 'Nowa sesja',
  [RodzajUtworzenia.Projekt]: 'Załóż projekt',
  [RodzajUtworzenia.Agent]: 'Skonfiguruj agenta',
  [RodzajUtworzenia.Automatyka]: 'Zbuduj automatykę',
  [RodzajUtworzenia.Kolejka]: 'Ustaw kolejkę',
  [RodzajUtworzenia.Zespol]: 'Złóż zespół',
};
