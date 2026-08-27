-- Migracja 271 wprowadza trwały zapis stanu kroku przebiegu wraz z ładunkami wejścia i wyjścia,
-- redagowanymi przy odczycie, a nie przy zapisie, aby regułę redakcji dało się zmieniać.
CREATE TABLE krok_przebiegu_automatyki (
    id              INTEGER PRIMARY KEY AUTOINCREMENT,
    przebieg_id     INTEGER NOT NULL REFERENCES przebieg_automatyki(id) ON DELETE CASCADE,
    krok_kod        TEXT    NOT NULL,
    stan            TEXT    NOT NULL DEFAULT 'pending',
    proba           INTEGER NOT NULL DEFAULT 0,
    komunikat_bledu TEXT,
    slad_stosu      TEXT,
    ladunek_wejscia TEXT,
    ladunek_wyjscia TEXT,
    kolejnosc       INTEGER NOT NULL DEFAULT 0,
    rozpoczeto      TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now')),
    zakonczono      TEXT,
    UNIQUE(przebieg_id, krok_kod)
);
CREATE INDEX idx_krok_przebiegu_automatyki_porzadek
    ON krok_przebiegu_automatyki(przebieg_id, kolejnosc, id);
