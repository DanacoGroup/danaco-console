-- Migracja 176 zakłada tabelę pobrań modułu przeglądarki, niosącą stan cyklu życia pobrania, postęp w bajtach oraz komunikat błędu przy niepowodzeniu.

CREATE TABLE pobranie_przegladania (
    id                       INTEGER PRIMARY KEY AUTOINCREMENT,
    identyfikator_zewnetrzny TEXT    NOT NULL UNIQUE,
    okno                     TEXT    NOT NULL,
    url                      TEXT    NOT NULL,
    nazwa_pliku              TEXT,
    sciezka_docelowa         TEXT,
    typ_mime                 TEXT,
    stan                     TEXT    NOT NULL DEFAULT 'queued',
    odebrano_bajtow          INTEGER,
    razem_bajtow             INTEGER,
    komunikat_bledu          TEXT,
    rozpoczeto               TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now')),
    zakonczono               TEXT
);
CREATE INDEX idx_pobranie_przegladania_okno ON pobranie_przegladania(okno, rozpoczeto DESC, id);
