-- Migracja 203 zakłada tabele serwera podglądu na żywo warstwy produktu oraz motywu produktu warsztatów frontendu i backendu modułu Apps.

-- ── Serwer podglądu warstwy produktu ─────────────────────────────────────────
-- Wartości kolumny `stan` są wartościami kontraktu (AppPreviewStatus).
CREATE TABLE podglad_apps (
    okno       TEXT PRIMARY KEY,
    warstwa    TEXT NOT NULL CHECK(warstwa IN ('frontend','backend')),
    adres      TEXT NOT NULL DEFAULT '',
    stan       TEXT NOT NULL DEFAULT 'stopped'
                    CHECK(stan IN ('starting','running','failed','stopped')),
    -- Milisekundy epoki — `AppsPreviewStartResponse.startedAt` niesie je wprost.
    rozpoczeto INTEGER NOT NULL DEFAULT 0,
    zatrzymano INTEGER
);

-- Zakłada tabelę motyw_apps niosącą surowy JSON motywu edytowanego w warsztacie stylów Frontend Workspace.
CREATE TABLE motyw_apps (
    okno           TEXT PRIMARY KEY,
    tresc          TEXT NOT NULL DEFAULT '{}',
    zaktualizowano TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now'))
);
