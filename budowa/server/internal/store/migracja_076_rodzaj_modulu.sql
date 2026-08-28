-- Migracja 076 — kolumna `rodzaj` w tabeli `modul` oraz treść ikon piętnastu modułów.
--
-- Moduły nie są jednym gatunkiem: część to środowiska robocze (okna czatu
-- z narzędziami i oknami pomocniczymi), część kompozytory wytwarzające pozycje
-- używane w innych modułach, jeden jest repozytorium plików, a czwartym
-- rodzajem jest sekcja konfiguracyjna. Bez kolumny `rodzaj` tabela `modul`
-- nie odróżnia ich niczym, więc podziału nie da się wyprowadzić z rdzenia
-- i strona główna musi trzymać własny wykaz kafli strefy 2.
--
-- Kolumna zostaje opcjonalna, tym samym prawem co `srodowisko.motto`: moduł
-- dołożony później, którego rodzaju jeszcze nie ustalono, ma prawo stać
-- z NULL-em. NULL znaczy „nie ustalono" i jest odróżnialny od każdej z czterech
-- wartości; NOT NULL wymuszałoby rodzaj wpisany byle jak przy zakładaniu
-- wiersza. Rodzaj `sekcja_konfiguracyjna` nie dostaje tu ani jednego modułu.
--
-- Bez warunku CHECK: SQLite nie umie dołożyć go przez ALTER TABLE, a
-- przepisanie tabeli pociągnęłoby za sobą klucze obce wskazujące `modul(id)`.
-- Zbioru czterech wartości —
--     srodowisko_robocze · kompozytor · repozytorium_plikow · sekcja_konfiguracyjna
-- — pilnuje treść tej migracji i warstwa danych, nie schemat.
--
-- Ikony przenoszą się z mapy `IKONY_MODULOW` klienta do bazy: kolumna
-- `modul.ikona` stała pusta, więc boczna nawigacja musiała trzymać kopię faktu
-- należącego do bazy. Każda z nazw istnieje w zestawie ikon klienta
-- (`client/src/ikony/zrodla/`).
--
-- Automations traci widoczność w bocznej nawigacji, bo moduł działa ze strefy 2
-- strony głównej i nie otwiera własnego okna modułowego. Droga wejścia zostaje:
-- kafel strefy 2 niesie kod modułu zawsze, a klient otwiera moduł spoza wykazu
-- nawigacji — `powloka/nawigacja-modulow.ts` zgłasza wskazanie bez odpowiednika
-- (`naBrakPozycji`), a `aplikacja/widok-srodowiska.ts` dobiera moduł z katalogu
-- `module.list`, zwracającego komplet modułów platformy niezależnie od macierzy.
-- Wiersze zostają w macierzy z `widoczny = 0`, a nie znikają: macierz ma być
-- kompletem par, a odsłonięcie modułu — zmianą bitu, nie wstawieniem wiersza.
--
-- `ALTER TABLE ... ADD COLUMN` w SQLite przepisuje sam nagłówek schematu, nie
-- tabelę; kolumna dopuszczająca NULL i bez DEFAULT nie dotyka ani jednego
-- wiersza danych. UPDATE-y idą w tej samej transakcji co ALTER (`migracje.go`
-- stosuje krok jednym `Exec` w jednej transakcji), więc schemat nie zostanie
-- zastosowany bez treści ani treść bez schematu. Dopasowanie idzie po `kod` —
-- jedynej kolumnie z warunkiem UNIQUE — więc każdy UPDATE trafia w dokładnie
-- jeden wiersz albo w żaden. Brak wiersza o danym kodzie nie jest błędem:
-- wierszy ten krok nie zakłada.

ALTER TABLE modul ADD COLUMN rodzaj TEXT;

-- ── Kompozytory ───────────────────────────────────────────────────────────────
-- Kompozytor nie jest miejscem, w którym się pracuje — jest wytwórnią: jego
-- wynik staje się pozycją do wyboru w innym module. Agent zbudowany w Agents
-- pojawia się na liście modeli w Studio; pętla złożona w Automations staje się
-- gotową automatyką do uruchomienia jednym kliknięciem w WorkSpace.
UPDATE modul SET rodzaj = 'kompozytor' WHERE kod = 'agents';
UPDATE modul SET rodzaj = 'kompozytor' WHERE kod = 'automations';

-- ── Repozytorium plików ───────────────────────────────────────────────────────
-- Menedżer zasobów bez AI — jedyny moduł tego rodzaju.
UPDATE modul SET rodzaj = 'repozytorium_plikow' WHERE kod = 'library';

-- ── Środowiska robocze — czat, narzędzia i okna pomocnicze ────────────────────
-- Każdy z nich jest oknem pracy z czatem i narzędziami: nie komponuje elementów
-- do użycia w innych modułach i nie jest menedżerem plików bez AI.
UPDATE modul SET rodzaj = 'srodowisko_robocze' WHERE kod = 'studio';
UPDATE modul SET rodzaj = 'srodowisko_robocze' WHERE kod = 'developer';
UPDATE modul SET rodzaj = 'srodowisko_robocze' WHERE kod = 'browser';
UPDATE modul SET rodzaj = 'srodowisko_robocze' WHERE kod = 'workspace';
UPDATE modul SET rodzaj = 'srodowisko_robocze' WHERE kod = 'research';
UPDATE modul SET rodzaj = 'srodowisko_robocze' WHERE kod = 'translate';
UPDATE modul SET rodzaj = 'srodowisko_robocze' WHERE kod = 'roundtable';
UPDATE modul SET rodzaj = 'srodowisko_robocze' WHERE kod = 'design';
UPDATE modul SET rodzaj = 'srodowisko_robocze' WHERE kod = 'assistant';
UPDATE modul SET rodzaj = 'srodowisko_robocze' WHERE kod = 'terminal';
UPDATE modul SET rodzaj = 'srodowisko_robocze' WHERE kod = 'diagnostics';
UPDATE modul SET rodzaj = 'srodowisko_robocze' WHERE kod = 'apps';

-- ── Automations znika z bocznej nawigacji, nie z platformy ───────────────────
-- Pary modułu dostają `widoczny = 0`. Nie ma tu DELETE: para zostaje
-- w macierzy, bo macierz ma być kompletem, a odsłonięcie modułu — zmianą bitu.
UPDATE srodowisko_modul
   SET widoczny = 0
 WHERE modul_id = (SELECT id FROM modul WHERE kod = 'automations');

-- ── Ikony piętnastu modułów (przeniesione z IKONY_MODULOW klienta) ────────────
UPDATE modul SET ikona = 'dokument'      WHERE kod = 'studio';
UPDATE modul SET ikona = 'folder'        WHERE kod = 'workspace';
UPDATE modul SET ikona = 'automatyzacja' WHERE kod = 'automations';
UPDATE modul SET ikona = 'karta-okna'    WHERE kod = 'browser';
UPDATE modul SET ikona = 'badanie'       WHERE kod = 'research';
UPDATE modul SET ikona = 'biblioteka'    WHERE kod = 'library';
UPDATE modul SET ikona = 'tlumacz'       WHERE kod = 'translate';
UPDATE modul SET ikona = 'debata'        WHERE kod = 'roundtable';
UPDATE modul SET ikona = 'paleta'        WHERE kod = 'design';
UPDATE modul SET ikona = 'mikrofon'      WHERE kod = 'assistant';
UPDATE modul SET ikona = 'terminal'      WHERE kod = 'terminal';
UPDATE modul SET ikona = 'kod'           WHERE kod = 'developer';
UPDATE modul SET ikona = 'diagnostyka'   WHERE kod = 'diagnostics';
UPDATE modul SET ikona = 'aplikacje'     WHERE kod = 'apps';
UPDATE modul SET ikona = 'agent'         WHERE kod = 'agents';
