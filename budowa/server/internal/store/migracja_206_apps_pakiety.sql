-- Migracja 206 zakłada tabelę pakietów rozszerzenia modułu Apps, niosącą odwołanie do archiwum, manifest i podpis jako surowy JSON kontraktu.

CREATE TABLE pakiet_apps (
    id                       INTEGER PRIMARY KEY AUTOINCREMENT,
    identyfikator_zewnetrzny TEXT    NOT NULL UNIQUE,
    okno                     TEXT    NOT NULL,
    -- Surowy JSON `AppPackageManifest`; puste, dopóki manifestu nie zapisano.
    manifest                 TEXT,
    -- Artefakt, z którego pakiet zbudowano (`AppArtifact.id`).
    artefakt_odwolanie       TEXT,
    -- Wartości kontraktu (AppPackageFormat) wprost.
    format                   TEXT    NOT NULL DEFAULT 'zip' CHECK(format IN ('zip','targz')),
    -- Odwołanie względne magazynu treści rdzenia do archiwum pakietu.
    sciezka                  TEXT,
    rozmiar                  INTEGER,
    -- Surowy JSON `ExtensionSignature`; puste, dopóki pakietu nie podpisano.
    podpis                   TEXT,
    -- Kod pozycji katalogu powstałej z pakietu po publikacji.
    rozszerzenie_kod         TEXT,
    utworzono                TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now')),
    zaktualizowano           TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now'))
);
CREATE INDEX idx_pakiet_apps_okno ON pakiet_apps(okno, id DESC);
