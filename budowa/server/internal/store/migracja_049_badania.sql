-- Migracja 049 zakłada trwałość modułu badań: źródła, ustalenia, raport
-- składany z ustaleń wraz z eksportem oraz jednowierszową przestrzeń badania.

-- ── Źródło badania — Sources Manager ───────────────────────────────────────
-- Wartości `rodzaj` i `wiarygodnosc` są wartościami kontraktu
-- (ResearchSourceKind, ResearchCredibility) wprost, bez tłumaczenia — wzór
-- `rodzaj` w migracji 039.
CREATE TABLE zrodlo_badania (
    id                       INTEGER PRIMARY KEY AUTOINCREMENT,
    identyfikator_zewnetrzny TEXT    NOT NULL UNIQUE,
    okno                     TEXT    NOT NULL,
    tytul                    TEXT    NOT NULL,
    rodzaj                   TEXT    NOT NULL DEFAULT 'web'
                                     CHECK(rodzaj IN ('web','document','note','dataset')),
    adres                    TEXT,
    pochodzenie              TEXT,
    wiarygodnosc             TEXT    NOT NULL DEFAULT 'unverified'
                                     CHECK(wiarygodnosc IN ('high','medium','low','unverified')),
    plik_biblioteki_id       TEXT,
    pozyskano_o              TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now'))
);
CREATE INDEX idx_zrodlo_badania_okno ON zrodlo_badania(okno, pozyskano_o DESC);

-- Ustalenie badania niesie treść przez parę kolumn tresc i tresc_odwolanie,
-- bo bywa krótką notatką albo długim akapitem przeniesionym z dokumentu.
CREATE TABLE ustalenie_badania (
    id                       INTEGER PRIMARY KEY AUTOINCREMENT,
    identyfikator_zewnetrzny TEXT    NOT NULL UNIQUE,
    okno                     TEXT    NOT NULL,
    tresc                    TEXT,
    tresc_odwolanie          TEXT,
    stan                     TEXT    NOT NULL DEFAULT 'open'
                                     CHECK(stan IN ('open','resolved')),
    utworzono                TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now')),
    zaktualizowano           TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now'))
);
CREATE INDEX idx_ustalenie_badania_okno ON ustalenie_badania(okno, zaktualizowano DESC);

-- Powiązanie ustalenia ze źródłami jest tabelą złącznikową relacji wiele do
-- wielu: jedno źródło zasila wiele ustaleń, jedno ustalenie może mieć wiele źródeł.
CREATE TABLE zrodlo_ustalenia_badania (
    ustalenie_id  INTEGER NOT NULL REFERENCES ustalenie_badania(id) ON DELETE CASCADE,
    zrodlo_id     INTEGER NOT NULL REFERENCES zrodlo_badania(id) ON DELETE CASCADE,
    utworzono     TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now')),
    PRIMARY KEY (ustalenie_id, zrodlo_id)
);
CREATE INDEX idx_zrodlo_ustalenia_badania_zrodlo ON zrodlo_ustalenia_badania(zrodlo_id, ustalenie_id);

-- Raport badania jest bytem własnym budowanym z ustaleń; sekcje raportu i ich
-- powiązania z ustaleniami mają osobne tabele niżej.
CREATE TABLE raport_badania (
    id                       INTEGER PRIMARY KEY AUTOINCREMENT,
    identyfikator_zewnetrzny TEXT    NOT NULL UNIQUE,
    okno                     TEXT    NOT NULL,
    tytul                    TEXT    NOT NULL,
    utworzono                TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now')),
    zaktualizowano           TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now'))
);
CREATE INDEX idx_raport_badania_okno ON raport_badania(okno, zaktualizowano DESC);

-- Sekcja raportu niesie treść tą samą parą kolumn co ustalenie, bo tak samo
-- bywa krótką notatką albo długim akapitem z dokumentu.
CREATE TABLE sekcja_raportu_badania (
    id                       INTEGER PRIMARY KEY AUTOINCREMENT,
    identyfikator_zewnetrzny TEXT    NOT NULL UNIQUE,
    raport_id                INTEGER NOT NULL REFERENCES raport_badania(id) ON DELETE CASCADE,
    tytul                    TEXT    NOT NULL,
    tresc                    TEXT,
    tresc_odwolanie          TEXT,
    kolejnosc                INTEGER NOT NULL DEFAULT 0,
    utworzono                TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now'))
);
CREATE INDEX idx_sekcja_raportu_badania_raport ON sekcja_raportu_badania(raport_id, kolejnosc, id);

-- Powiązanie sekcji raportu z ustaleniami jest tabelą złącznikową: jedno
-- ustalenie zasila wiele sekcji, jedna sekcja opiera się na wielu ustaleniach.
CREATE TABLE ustalenie_sekcji_raportu_badania (
    sekcja_id     INTEGER NOT NULL REFERENCES sekcja_raportu_badania(id) ON DELETE CASCADE,
    ustalenie_id  INTEGER NOT NULL REFERENCES ustalenie_badania(id) ON DELETE CASCADE,
    utworzono     TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now')),
    PRIMARY KEY (sekcja_id, ustalenie_id)
);
CREATE INDEX idx_ustalenie_sekcji_raportu_badania_ustalenie ON ustalenie_sekcji_raportu_badania(ustalenie_id, sekcja_id);

-- ── Eksport raportu — ślad wyjścia (Report Builder) ─────────────────────────
-- Wartości `format` są wartościami kontraktu (ExportFormat) wprost.
CREATE TABLE eksport_raportu_badania (
    id                       INTEGER PRIMARY KEY AUTOINCREMENT,
    identyfikator_zewnetrzny TEXT    NOT NULL UNIQUE,
    raport_id                INTEGER NOT NULL REFERENCES raport_badania(id) ON DELETE CASCADE,
    format                   TEXT    NOT NULL
                                     CHECK(format IN ('pdf','docx','markdown','html','txt')),
    sciezka_docelowa         TEXT,
    plik_biblioteki_id       TEXT,
    sciezka_wyniku           TEXT,
    rozmiar_bajtow           INTEGER,
    utworzono                TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now'))
);
CREATE INDEX idx_eksport_raportu_badania_raport ON eksport_raportu_badania(raport_id, utworzono DESC);

-- Przestrzeń badania jest bytem własnym bez okna, jednowierszowym: klucz
-- główny jest wymuszony na jeden, a zapis nadpisuje ten jeden wiersz zamiast zakładać nowy.
CREATE TABLE przestrzen_badania (
    id             INTEGER NOT NULL PRIMARY KEY CHECK(id = 1),
    zakres         TEXT    NOT NULL DEFAULT '',
    zaktualizowano TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now'))
);

-- Etapy badania — lista uporządkowana, osobna tabela bo `Stages []string` nie
-- ma tożsamości własnej (jak etykieta pliku w migracji 045), a kolejność się
-- liczy (kontrakt zwraca `Stages` w kolejności zapisu).
CREATE TABLE etap_przestrzeni_badania (
    kolejnosc  INTEGER NOT NULL,
    etap       TEXT    NOT NULL,
    PRIMARY KEY (kolejnosc)
);
