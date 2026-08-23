import './wiedza.css';

export {
  utworzOknoWyszukiwaniaZnaczenia,
  type OknoWyszukiwaniaZnaczenia,
} from './okno-wyszukiwania-znaczenia';
export { utworzZrodloWiedzy, type OwocSzukania, type OwocWskaznika, type ZrodloWiedzy } from './zrodlo-wiedzy';

/**
 * Wejście katalogu WIEDZA — wyszukiwanie po znaczeniu (`knowledge.*`).
 *
 * Katalog stoi osobno od `poczta/`, bo poczta i wiedza to dwie różne dziedziny
 * rdzenia: `mail.*` sięga po skrzynkę Operatora stojącą poza urządzeniem,
 * `knowledge.*` po jego własne treści leżące w rdzeniu. Wspólny katalog zlepiłby
 * dwa słowniki pojęć w jeden.
 *
 * Własnego wpisu w rejestrze modułów katalog mimo to nie ma. Rejestr wiąże widok
 * z kodem modułu rdzenia (`modul.kod`), a tabela `modul`
 * (`migracja_007_zaczyn_slownikow.sql`) nie niesie ani modułu poczty, ani modułu
 * wiedzy — drugi wpis rejestru czekałby na kod, którego nawigacja nigdy nie poda.
 * Okno wchodzi do sceny przez złożenie modułu Poczta (`moduly/poczta/indeks.ts`),
 * które ma w rejestrze jeden wpis. Dlatego ten plik nie wystawia `MODUL`.
 */
