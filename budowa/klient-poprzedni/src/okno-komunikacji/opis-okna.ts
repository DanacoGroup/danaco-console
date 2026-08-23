import {
  ExecutionEnv,
  PermissionMode,
  WindowRole,
} from '../../../shared/contract';

/**
 * Opis okna komunikacji.
 *
 * Okno jest bytem pośrednim między sesją a wiadomością i niesie własny moduł,
 * kanał modelu, katalogi robocze, środowisko wykonania, tryb uprawnień oraz
 * rolę. Zasięg wykonania jest parametrem okna, nie właściwością wdrożenia.
 * Wartości pól pochodzą z kontraktu — okno nie prowadzi własnego katalogu nazw.
 */
export interface OpisOkna {
  /** Projekt, w którego kontekście pracuje okno. */
  projekt: string;
  /** Tytuł okna widoczny dla operatora. */
  tytul: string;
  /** Moduł przypisany oknu (`moduleId` kontraktu). */
  modul: string;
  /** Kanał modelu obsługujący okno (`modelChannelId` kontraktu). */
  kanalModelu: string;
  /** Katalogi robocze udostępnione oknu. */
  katalogiRobocze: string[];
  /** Środowisko wykonania modelu: urządzenie, rdzeń albo host zdalny. */
  srodowiskoWykonania: ExecutionEnv;
  /** Tryb uprawnień okna. */
  trybUprawnien: PermissionMode;
  /** Rola okna w pętli koordynator–wykonawca. */
  rola: WindowRole;
}

/**
 * Opis początkowy okna.
 *
 * Brak ustawienia oznacza wartość domyślną, nie blokadę uruchomienia. Domyślnie:
 * wykonanie na urządzeniu operatora, tryb uprawnień ręczny, okno samodzielne.
 *
 * Moduł i kanał modelu zostają puste. `KnownModuleIds[0]` to kod środowiska, nie
 * modułu, a `KnownChannelKinds[0]` to rodzaj kanału, nie kod wiersza rejestru —
 * podstawienie któregokolwiek wskazywałoby moduł spoza katalogu i kanał, którego
 * rejestr nie ma, więc pierwsza wypowiedź wracałaby odmową. Widok nie zna
 * rejestru rdzenia w chwili zakładania okna; pierwszy kanał czynny podstawia rdzeń.
 */
export function opisPoczatkowy(): OpisOkna {
  return {
    projekt: ustawienie(import.meta.env.VITE_PROJEKT, 'Danaco Console'),
    tytul: ustawienie(import.meta.env.VITE_TYTUL_OKNA, 'Okno komunikacji'),
    modul: ustawienie(import.meta.env.VITE_MODUL, ''),
    kanalModelu: ustawienie(import.meta.env.VITE_KANAL_MODELU, ''),
    katalogiRobocze: katalogiRobocze(import.meta.env.VITE_KATALOGI_ROBOCZE),
    srodowiskoWykonania: ExecutionEnv.Local,
    trybUprawnien: PermissionMode.Manual,
    rola: WindowRole.Standalone,
  };
}

/** Wartość z konfiguracji budowania albo wartość domyślna. */
function ustawienie(wartosc: unknown, domyslna: string): string {
  return typeof wartosc === 'string' && wartosc.length > 0 ? wartosc : domyslna;
}

/** Katalogi robocze zapisane w konfiguracji jako lista rozdzielona średnikiem. */
function katalogiRobocze(wartosc: unknown): string[] {
  if (typeof wartosc !== 'string') return [];
  return wartosc
    .split(';')
    .map((katalog) => katalog.trim())
    .filter((katalog) => katalog.length > 0);
}
