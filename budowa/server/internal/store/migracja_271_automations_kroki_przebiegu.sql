-- Migracja 271 — stan kroku przebiegu wraz z ładunkami
-- (`automation.execution.steps`, `automation.execution.payload.get`).
--
-- Stan przebiegu jest wyprowadzany z pozycji kolejki i tak zostaje. Ta tabela
-- jest czym innym: zapisem TRWAŁYM stanu kroku wraz z jego danymi wejściowymi
-- i wyjściowymi. Kolejka po opróżnieniu (`clear`) traci pozycje, a drążenie do
-- poziomu kroku i podgląd ładunku mają dalej mieć co pokazać — inaczej historia
-- przebiegu kończyłaby się w chwili sprzątnięcia kolejki.
--
-- Ładunki są redagowane przy odczycie, nie przy zapisie: reguła redakcji
-- sekretów bywa zmieniana, a zapis zredagowany nie da się odredagować.
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
