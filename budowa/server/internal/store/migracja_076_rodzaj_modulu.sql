-- Migracja dodaje do tabeli modul opcjonalną kolumnę rodzaj o czterech wartościach
-- oraz przenosi ikony piętnastu modułów z mapy klienta do bazy.

ALTER TABLE modul ADD COLUMN rodzaj TEXT;

-- Kompozytor nie jest miejscem pracy, lecz wytwórnią: jego wynik staje się pozycją
-- do wyboru w innym module, jak agent w Studio albo automatyka w WorkSpace.
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

-- Ikony piętnastu modułów przenoszą się tu z mapy IKONY_MODULOW klienta do kolumny
-- modul.ikona w bazie.
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
