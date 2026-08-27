import {
  ExecutionEnv,
  PermissionMode,
  WindowRole,
} from '../../../shared/contract';

/**
 * Opis okna komunikacji jest bytem pośrednim między sesją a wiadomością, niosącym moduł, kanał modelu, katalogi robocze, środowisko wykonania, tryb uprawnień oraz rolę, przy czym wartości pól pochodzą z kontraktu.
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
 * Opis początkowy okna zakłada wykonanie na urządzeniu operatora, tryb uprawnień ręczny i okno samodzielne, z pustym modułem i kanałem modelu do chwili podstawienia ich przez pierwszy kanał czynny w rejestrze rdzenia.
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

/** Wartość z konfiguracji budowania używana przez okno, a w razie jej braku zastosowana zostaje wartość domyślna. */
function ustawienie(wartosc: unknown, domyslna: string): string {
  return typeof wartosc === 'string' && wartosc.length > 0 ? wartosc : domyslna;
}

/** Katalogi robocze przekazane oknu, zapisane w konfiguracji budowania jako lista ścieżek rozdzielona średnikiem. */
function katalogiRobocze(wartosc: unknown): string[] {
  if (typeof wartosc !== 'string') return [];
  return wartosc
    .split(';')
    .map((katalog) => katalog.trim())
    .filter((katalog) => katalog.length > 0);
}
