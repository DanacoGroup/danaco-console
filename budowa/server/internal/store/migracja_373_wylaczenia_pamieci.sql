-- Migracja 373 — wyłączenia pamięci w zasięgu.
--
-- Powód jest jeden: rozstrzygnięcie Właściciela z 17.08.2026 — „niech będzie
-- opcja wyłączenia: wyłączenia całkiem, wyłączenia tylko dla niektórych modułów
-- itp., ale to wszystko ma się sterować z pozycji Operatora w konfiguracji".
-- Do tej migracji schemat nie znał ŻADNEJ tabeli wyłączeń, więc żądanie
-- wyłączenia ustalenia wspólnego w jednym projekcie albo module kończyło się
-- odmową `conflict`, nazwaną wprost w `core/adapter_modul_workspace_pamiec_
-- komendy.go`. Odmowa była prawdziwa i dlatego trzeba było tabeli, nie zmiany
-- brzmienia odmowy.
--
-- WYŁĄCZENIE NIE JEST USUNIĘCIEM i to jest cała różnica wobec `memory.delete`.
-- Wiersz tutaj nie rusza `wpis_pamieci_projektu`: treść stoi nietknięta,
-- a zniesienie wyłączenia (usunięcie wiersza) wraca wpis do kontekstu
-- w całości. Gdyby wyłączenie kasowało treść, Operator wyłączający pamięć „na
-- czas jednego zadania" traciłby ustalenia bezpowrotnie.
--
-- WYŁĄCZENIE NIE JEST TEŻ ODPIĘCIEM (`memory.detach`). Odpięcie zwęża zasięg
-- SAMEGO WPISU do jego projektu — zmienia wiersz pamięci. Wyłączenie zostawia
-- zasięg wpisu nietknięty i wstrzymuje go w zasięgu, w którym Operator go nie
-- chce; ten sam wpis obowiązuje dalej wszędzie indziej.
--
-- Dwa byty wyłączane, nie jeden: albo pojedynczy wpis (`wpis_id`), albo cały
-- poziom pamięci (`poziom_pamieci`). „Wyłącz pamięć projektu w module Developer"
-- nie wskazuje żadnego wpisu i musi się zapisać bez niego; „wyłącz to jedno
-- ustalenie w tej karcie sesji" wskazuje wpis i nie dotyczy poziomu. Wiersz
-- z obydwoma polami pustymi nie wyłącza niczego i schemat go nie przyjmuje
-- (CHECK poniżej).
--
-- Zasięg idzie tabelą `poziom_zasiegu` — tą samą, którą idzie zasięg wpisu
-- pamięci i każde ustawienie platformy. Sześć zasięgów wymienionych przez
-- Właściciela to sześć jej wierszy: `globalny` znaczy „całkiem", dalej
-- `srodowisko`, `projekt`, `modul`, `para_modulow`, `karta_sesji`. Drugiego
-- porządku poziomów w tym pliku nie ma i nie ma po co być.
--
-- Poziom pamięci zapisuje się wartością kontraktu (`global`, `environment`,
-- `project`, `session`), nie polskim kodem: poziomy pamięci eksperta jadą tak
-- samo (`agent.memoryLevels`), a drugie odwzorowanie tego samego wyliczenia
-- byłoby drugą prawdą o czterech wartościach.
--
-- Wymaga restartu: nie.

CREATE TABLE IF NOT EXISTS wylaczenie_pamieci (
    id                     INTEGER PRIMARY KEY AUTOINCREMENT,
    identyfikator_zewnetrzny TEXT  NOT NULL UNIQUE,
    wpis_id                INTEGER          REFERENCES wpis_pamieci_projektu(id) ON DELETE CASCADE,
    poziom_pamieci         TEXT    NOT NULL DEFAULT ''
                                   CHECK(poziom_pamieci IN ('', 'global', 'environment',
                                                            'project', 'session')),
    poziom_zasiegu_id      INTEGER NOT NULL REFERENCES poziom_zasiegu(id) ON DELETE CASCADE,
    klucz_zasiegu          TEXT    NOT NULL DEFAULT '',
    utworzono              TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now')),
    zaktualizowano         TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now')),
    -- Wyłączenie bez wskazanego bytu nie wyłącza niczego; wyłączenie wskazujące
    -- oba naraz nie mówi, czy dotyczy wpisu, czy poziomu. Jedno albo drugie.
    CHECK((wpis_id IS NOT NULL AND poziom_pamieci = '')
          OR (wpis_id IS NULL AND poziom_pamieci <> ''))
);

-- Jedno wyłączenie na byt i zasięg: powtórzone żądanie nie zakłada drugiego
-- wiersza, bo „wyłączone dwa razy" nie jest stanem, który da się znieść jednym
-- ruchem.
CREATE UNIQUE INDEX IF NOT EXISTS idx_wylaczenie_pamieci_byt
    ON wylaczenie_pamieci(IFNULL(wpis_id, 0), poziom_pamieci,
                          poziom_zasiegu_id, klucz_zasiegu);

CREATE INDEX IF NOT EXISTS idx_wylaczenie_pamieci_zasieg
    ON wylaczenie_pamieci(poziom_zasiegu_id, klucz_zasiegu);
