-- Migracja 161 zakłada tabele zestawów reguł segmentacji, reguł zestawu oraz trwałych segmentów okna tłumaczenia, zapisywanych leniwie od pierwszego scalenia albo podziału.

-- Zakłada tabelę zestaw_regul_segmentacji niosącą nastawę wielokrotnego użytku z treścią pliku SRX podaną przez Operatora.
CREATE TABLE zestaw_regul_segmentacji (
    id                       INTEGER PRIMARY KEY AUTOINCREMENT,
    identyfikator_zewnetrzny TEXT    NOT NULL UNIQUE,
    nazwa                    TEXT    NOT NULL,
    jezyk                    TEXT,
    srx                      TEXT,
    zaktualizowano           INTEGER NOT NULL
);
CREATE INDEX idx_zestaw_regul_segmentacji_jezyk ON zestaw_regul_segmentacji(jezyk, nazwa);

-- ── Reguła zestawu (dziecko zestawu) ────────────────────────────────────────
-- `lamie` mówi, czy zderzenie wzorców jest granicą zdania (SRX: break), czy
-- wyjątkiem od granicy (SRX: no-break) — na przykład kropka po skrócie.
CREATE TABLE regula_segmentacji (
    id         INTEGER PRIMARY KEY AUTOINCREMENT,
    zestaw_id  INTEGER NOT NULL REFERENCES zestaw_regul_segmentacji(id) ON DELETE CASCADE,
    kolejnosc  INTEGER NOT NULL,
    przed      TEXT    NOT NULL,
    po         TEXT    NOT NULL,
    lamie      INTEGER NOT NULL DEFAULT 1 CHECK(lamie IN (0,1))
);
CREATE INDEX idx_regula_segmentacji_zestaw ON regula_segmentacji(zestaw_id, kolejnosc);

-- ── Segment okna tłumaczenia ────────────────────────────────────────────────
-- Numer (`kolejnosc`) jest tożsamością segmentu w obrębie okna — kontrakt
-- adresuje segmenty numerem (`segmentIndexes`, `segmentIndex`), nie kodem.
CREATE TABLE segment_okna_tlumaczenia (
    id        INTEGER PRIMARY KEY AUTOINCREMENT,
    okno_id   INTEGER NOT NULL REFERENCES okno_tlumaczenia(id) ON DELETE CASCADE,
    kolejnosc INTEGER NOT NULL,
    tresc     TEXT    NOT NULL,
    UNIQUE(okno_id, kolejnosc)
);
