import type {
  EnvironmentEnterRequest,
  EnvironmentEnterResponse,
  ModuleListRequest,
  ModuleListResponse,
} from '../../../shared/contract';
import type { Wynik } from './kanal';

/** Nawigacja platformy — dwie komendy, których potrzebuje boczna nawigacja powłoki, podane jej jako jeden byt. */
export interface NawigacjaPlatformy {
  /** `environment.enter` — wejście do środowiska wraz z jego nawigacją i kartami sesji. */
  wejdzDoSrodowiska(zadanie: EnvironmentEnterRequest): Promise<Wynik<EnvironmentEnterResponse>>;
  /** `module.list` — moduły platformy albo moduły widoczne w jednym środowisku. */
  wykazModulow(zadanie?: ModuleListRequest): Promise<Wynik<ModuleListResponse>>;
}
