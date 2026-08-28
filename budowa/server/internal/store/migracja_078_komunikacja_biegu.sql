-- Migracja 078 wprowadza komunikację między modelami: zakres tego, co następny
-- widzi z pracy poprzednika, oraz ślad narzędzi, niezależny od zakresu i od
-- profilu izolacji.

-- Tabela przekazanie_biegu przechowuje plan rozmowy: wiązania między miejscami
-- obsady wraz z zakresem i warunkiem przekazania.
CREATE TABLE przekazanie_biegu (
    id                 INTEGER PRIMARY KEY AUTOINCREMENT,
    bieg_id            INTEGER NOT NULL REFERENCES bieg_orkiestracji(id) ON DELETE CASCADE,
    od_miejsca         INTEGER NOT NULL CHECK(od_miejsca BETWEEN 1 AND 4),
    do_miejsca         INTEGER NOT NULL CHECK(do_miejsca BETWEEN 1 AND 4),
    zakres             TEXT    NOT NULL DEFAULT 'artefakt'
                               CHECK(zakres IN ('artefakt','streszczenie','wypowiedz')),
    ze_sladem_narzedzi INTEGER NOT NULL DEFAULT 0 CHECK(ze_sladem_narzedzi IN (0,1)),
    -- Warunek przekazania; puste znaczy zawsze. Ten sam kształt co warunek kroku automatyki.
    warunek            TEXT,
    kolejnosc          INTEGER NOT NULL DEFAULT 0,
    utworzono          TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now')),
    -- Stanowisko nie przekazuje samo sobie — to byłaby pętla bez wyjścia.
    CHECK (od_miejsca <> do_miejsca),
    UNIQUE (bieg_id, od_miejsca, do_miejsca)
);

CREATE INDEX idx_przekazanie_biegu_od
    ON przekazanie_biegu(bieg_id, od_miejsca, kolejnosc);

-- Tabela wiadomosc_biegu przechowuje ślad tego, co jedno stanowisko rzeczywiście
-- przekazało drugiemu, wraz ze źródłową pozycją kolejki.
CREATE TABLE wiadomosc_biegu (
    id                 INTEGER PRIMARY KEY AUTOINCREMENT,
    bieg_id            INTEGER NOT NULL REFERENCES bieg_orkiestracji(id) ON DELETE CASCADE,
    od_miejsca         INTEGER NOT NULL CHECK(od_miejsca BETWEEN 1 AND 4),
    do_miejsca         INTEGER NOT NULL CHECK(do_miejsca BETWEEN 1 AND 4),
    zakres             TEXT    NOT NULL
                               CHECK(zakres IN ('artefakt','streszczenie','wypowiedz')),
    tresc              TEXT    NOT NULL,
    slad_narzedzi      TEXT,
    -- Pozycja kolejki, która wytworzyła treść; puste dla przekazania spoza kolejki.
    pozycja_kolejki_id INTEGER REFERENCES pozycja_kolejki(id) ON DELETE SET NULL,
    -- Podagent, który treść wytworzył; puste, gdy wytworzyło ją samo stanowisko.
    podagent_id        INTEGER REFERENCES podagent(id) ON DELETE SET NULL,
    utworzono          TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now')),
    CHECK (od_miejsca <> do_miejsca)
);

-- Indeks porządkuje odczyt rozmowy biegu w kolejności powstawania wiadomości,
-- w takiej, w jakiej ogląda ją operator okna.
CREATE INDEX idx_wiadomosc_biegu_bieg
    ON wiadomosc_biegu(bieg_id, id);

-- Indeks wspiera pytanie, co dostało dane stanowisko: wejście, które otrzymuje
-- ono przy rozpoczęciu swojej tury.
CREATE INDEX idx_wiadomosc_biegu_do
    ON wiadomosc_biegu(bieg_id, do_miejsca, id);
