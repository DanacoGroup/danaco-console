-- Migracja 073 wprowadza tożsamość własną eksperta, warstwy jego promptu
-- oraz wtyczki, dodając kolumny i tabele bez zmiany istniejących.

-- ── (a) Tożsamość własna eksperta ─────────────────────────────────────────────
-- Imię własne jest tym, co widzi Operator, i pozostaje niezależne od nazwy
-- technicznej, po której idzie sortowanie i wyszukiwanie wykazu.
ALTER TABLE agent ADD COLUMN imie_wlasne TEXT NOT NULL DEFAULT '';
ALTER TABLE agent ADD COLUMN favikon     TEXT NOT NULL DEFAULT '';

-- ── (b) Warstwy promptu eksperta ──────────────────────────────────────────────
-- Warstwy promptu mają osobną tabelę, nie stałe kolumny, bo katalog warstw
-- silnika nakładki jest otwarty na warstwę spoza trzech znanych dziś nazw.
CREATE TABLE agent_warstwa (
    agent_id       INTEGER NOT NULL REFERENCES agent(id) ON DELETE CASCADE,
    warstwa        TEXT    NOT NULL CHECK (warstwa IN ('constitution', 'profile', 'expertise')),
    tresc          TEXT    NOT NULL DEFAULT '',
    tryb           TEXT    NOT NULL DEFAULT 'DOLACZ' CHECK (tryb IN ('ZASTAP', 'DOLACZ')),
    -- Warstwa wyłączona zostaje zapisana i widoczna w oknie, nie wchodzi
    -- do złożonego promptu.
    aktywna        INTEGER NOT NULL DEFAULT 1 CHECK (aktywna IN (0, 1)),
    zaktualizowano TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ', 'now')),
    PRIMARY KEY (agent_id, warstwa)
);

-- ── (c) Wtyczki eksperta ──────────────────────────────────────────────────────
-- Wtyczki mają osobną tabelę od konektorów, bo są odrębnym bytem: konektor
-- wskazuje most usługi, a wtyczka niesie nazwę, źródło i wersję rozszerzenia
-- powłoki.
CREATE TABLE agent_wtyczka (
    id        INTEGER PRIMARY KEY AUTOINCREMENT,
    kod       TEXT    NOT NULL UNIQUE,
    agent_id  INTEGER NOT NULL REFERENCES agent(id) ON DELETE CASCADE,
    nazwa     TEXT    NOT NULL,
    -- Skąd wtyczka pochodzi: katalog na dysku, adres repozytorium albo
    -- nazwa rejestru.
    zrodlo    TEXT,
    wersja    TEXT,
    aktywna   INTEGER NOT NULL DEFAULT 1 CHECK (aktywna IN (0, 1)),
    utworzono TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ', 'now'))
);

-- Wykaz wtyczek jednego eksperta jest porządkowany po nazwie, tak samo jak
-- wykaz konektorów tego samego agenta.
CREATE INDEX idx_agent_wtyczka_agent ON agent_wtyczka(agent_id, nazwa);
