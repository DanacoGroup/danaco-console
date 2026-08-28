import { ExecutionEnv, WindowRole } from '../../../shared/contract';
import type { EtapUzgodnienia } from '../protokol/uzgodnienie';

/**
 * Nazwy prezentacyjne wartości kontraktu.
 *
 * Kontrakt niesie nazwy angielskie; warstwa widoku pokazuje je operatorowi
 * po polsku. Tłumaczenie żyje w jednym pliku, żeby nie rozproszyć się po
 * elementach interfejsu.
 */
export function nazwaRoli(rola: WindowRole): string {
  switch (rola) {
    case WindowRole.Executor:
      return 'Wykonawca';
    case WindowRole.Coordinator:
      return 'Koordynator';
    case WindowRole.Standalone:
      return 'Samodzielne';
  }
}

/** Nazwa środowiska wykonania modelu — lokalne urządzenie operatora, rdzeń albo host zdalny — pokazywana w interfejsie po polsku. */
export function nazwaSrodowiska(srodowisko: ExecutionEnv): string {
  switch (srodowisko) {
    case ExecutionEnv.Local:
      return 'Urządzenie operatora';
    case ExecutionEnv.Core:
      return 'Rdzeń';
    case ExecutionEnv.Remote:
      return 'Host zdalny';
  }
}

/** Nazwa etapu uzgodnienia z rdzeniem, pokazywana operatorowi po polsku jako czytelny odpowiednik wartości kontraktu. */
export function nazwaEtapu(etap: EtapUzgodnienia): string {
  switch (etap) {
    case 'powitanie':
      return 'Powitanie połączenia';
    case 'sesja':
      return 'Założenie sesji';
    case 'okno':
      return 'Otwarcie okna komunikacji';
    case 'gotowe':
      return 'Okno gotowe do pracy';
  }
}
