/**
 * Wejście katalogu wiedzy, udostępniające wyszukiwanie po znaczeniu i złożone
 * do sceny przez moduł poczty, bez własnego wpisu w rejestrze modułów.
 */
import './wiedza.css';

export {
  utworzOknoWyszukiwaniaZnaczenia,
  type OknoWyszukiwaniaZnaczenia,
} from './okno-wyszukiwania-znaczenia';
export { utworzZrodloWiedzy, type OwocSzukania, type OwocWskaznika, type ZrodloWiedzy } from './zrodlo-wiedzy';
