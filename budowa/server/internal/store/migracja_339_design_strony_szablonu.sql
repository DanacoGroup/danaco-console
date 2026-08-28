-- Migracja 339 — strony szablonu materiału modułu Design: publikacje
-- wielostronicowe (książka, broszura, katalog) na istniejącym druku.
--
-- ── Dlaczego NIE ma rodziny `design.publication.*` ──────────────────────────
-- Właściciel wybrał szablon materiału o wielu stronach, a nie nową rodzinę
-- komend i nie montaż arkusza. Publikacja jest więc szablonem, który ma strony,
-- i wychodzi tą samą drogą, co każdy inny materiał — przez `design.print.export`,
-- przez tę samą kontrolę przeddrukową i z tą samą odmową przy wadzie o wadze
-- błędu. Druga rodzina komend znaczyłaby drugą kontrolę przeddrukową i drugie
-- wydanie, a te dwa musiałyby się potem zgadzać.
--
-- ── Strona jest bytem, warstwa należy do strony ─────────────────────────────
-- Szablon jednostronicowy ma warstwy w `warstwa_szablonu_materialu_design`
-- (migracja 336) i tak zostaje — wsteczna zgodność nie jest tu ustępstwem, tylko
-- prawdą: baner nie ma stron. Publikacja dokłada wiersze STRON, a warstwa strony
-- wskazuje stronę kluczem obcym. Szablon bez ani jednej strony jest szablonem
-- jednostronicowym i jego warstwy leżą tam, gdzie leżały.
--
-- ── Numer strony jest jawny, nie wynika z klucza ────────────────────────────
-- Kolejność stron publikacji zmienia się (przestawienie rozkładówki), a klucz
-- wiersza nie. Numer jest więc kolumną z unikatem na parze (szablon, numer):
-- dwie strony o tym samym numerze dałyby przy wydaniu kolejność zależną od
-- porządku odczytu, czyli żadną.

CREATE TABLE strona_szablonu_materialu_design (
    id                       INTEGER PRIMARY KEY AUTOINCREMENT,
    identyfikator_zewnetrzny TEXT    NOT NULL UNIQUE,
    szablon_id               INTEGER NOT NULL
                                     REFERENCES szablon_materialu_design(id) ON DELETE CASCADE,
    numer                    INTEGER NOT NULL,
    nazwa                    TEXT,
    UNIQUE(szablon_id, numer)
);

CREATE INDEX idx_strona_szablonu_materialu_design_szablon
    ON strona_szablonu_materialu_design(szablon_id, numer, id);

CREATE TABLE warstwa_strony_szablonu_design (
    id                       INTEGER PRIMARY KEY AUTOINCREMENT,
    identyfikator_zewnetrzny TEXT    NOT NULL,
    strona_id                INTEGER NOT NULL
                                     REFERENCES strona_szablonu_materialu_design(id)
                                     ON DELETE CASCADE,
    zasob_id                 TEXT,
    x                        REAL,
    y                        REAL,
    szerokosc                REAL,
    wysokosc                 REAL,
    kolejnosc                INTEGER NOT NULL DEFAULT 0,
    zablokowana              INTEGER NOT NULL DEFAULT 0 CHECK(zablokowana IN (0,1)),
    adnotacja                TEXT
);

CREATE INDEX idx_warstwa_strony_szablonu_design_strona
    ON warstwa_strony_szablonu_design(strona_id, kolejnosc, id);
