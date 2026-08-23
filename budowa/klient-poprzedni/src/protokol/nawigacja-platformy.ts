import type {
  EnvironmentEnterRequest,
  EnvironmentEnterResponse,
  ModuleListRequest,
  ModuleListResponse,
} from '../../../shared/contract';
import type { Wynik } from './kanal';

/**
 * Nawigacja platformy — dwie komendy, których potrzebuje boczna nawigacja
 * powłoki, podane jej jako jeden byt.
 *
 * Plik niesie wyłącznie kształt zależności, nie drogę do rdzenia. Wywołania
 * stoją w plikach jednokomendowych (`wejscie-do-srodowiska.ts`,
 * `wykaz-modulow.ts`), a obiekt tego kształtu składa korzeń montażu klienta
 * (`aplikacja/polaczenie-z-rdzeniem.ts`) — jedyne miejsce, które ma kanał.
 * Powłoka dostaje go gotowego i nie zna nazwy ani jednej komendy.
 *
 * Metody są dwie, bo tyle woła powłoka. `home.enter` woła wprost widok strony
 * głównej, `workspace.enter` — przestrzeń modułu, a `environment.list` czyta
 * strona główna kanałem.
 */
export interface NawigacjaPlatformy {
  /** `environment.enter` — wejście do środowiska wraz z jego nawigacją i kartami sesji. */
  wejdzDoSrodowiska(zadanie: EnvironmentEnterRequest): Promise<Wynik<EnvironmentEnterResponse>>;
  /** `module.list` — moduły platformy albo moduły widoczne w jednym środowisku. */
  wykazModulow(zadanie?: ModuleListRequest): Promise<Wynik<ModuleListResponse>>;
}
