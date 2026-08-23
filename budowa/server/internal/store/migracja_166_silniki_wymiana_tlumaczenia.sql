-- Migracja 166 — nastawy silników i wymiana zewnętrzna modułu Translate:
-- profile silników (`translate.engine.profile.*`), polityka tłumaczenia
-- pivotowego (`translate.pivot.policy.*`), pakiety przekazania
-- (`translate.handoff.*`), pakietowy przebieg operacji (`translate.batch.run`)
-- oraz most do dokumentu (`translate.bridge.*`).
--
-- Profil silnika wskazuje kanały modelu wykazem, więc kanały mają tabelę
-- dziecka: wykaz w jednej kolumnie nie da się złączyć z `kanal_modelu` ani
-- zawęzić zapytaniem „które profile używają tego kanału". Kolejność w wykazie
-- jest znacząca (pierwszy kanał jest kanałem pierwszego wyboru), stąd kolumna
-- `kolejnosc`.
--
-- Polityka pivota jest jedna na zasięg (`ConfigScope` + byt zasięgu), a pary
-- języków są jej dzieckiem. Para mówi: przekład z języka A na B idzie przez
-- język C. `jezyk_domyslny` jest pivotem dla par, których Operator nie wymienił.
--
-- Pakiet przekazania ma stan (`HandoffStatus`), bo jego cykl życia jest realny:
-- złożony → przekazany → zwrot przyjęty. Zwrot (`handoff.receive`) przestawia
-- ten sam wiersz, zamiast zakładać drugi — inaczej nie dałoby się odpowiedzieć,
-- czy materiał wrócił.
--
-- Przebieg pakietowy ma własną kolejkę, nie kolejkę modułu Automations:
-- `batch.run` oddaje `queueId`, po którym okno Translate pyta o postęp, a jego
-- pozycje niosą operacje własne tego modułu (przekład, kontrola jakości,
-- korekta, wydanie). Wpięcie w cudzą kolejkę związałoby moduł z modułem, który
-- Operator ma prawo mieć wyłączony.
--
-- Most do dokumentu jest wiązaniem dwustronnym: okno Translate wie, z którego
-- dokumentu wzięło materiał, a wysyłka wyniku wie, dokąd go odesłać. Bez
-- wiersza `bridge.result.send` musiałby dostać wskazanie dokumentu drugi raz,
-- czego kontrakt nie przewiduje.

-- ── Profil silnika tłumaczenia ──────────────────────────────────────────────
CREATE TABLE profil_silnika_tlumaczenia (
    id                       INTEGER PRIMARY KEY AUTOINCREMENT,
    identyfikator_zewnetrzny TEXT    NOT NULL UNIQUE,
    nazwa                    TEXT    NOT NULL,
    dziedzina                TEXT,
    zasieg                   TEXT    NOT NULL DEFAULT 'global',
    zasieg_id                TEXT,
    adaptacyjny              INTEGER NOT NULL DEFAULT 0 CHECK(adaptacyjny IN (0,1)),
    zasieg_pamieci           TEXT    CHECK(zasieg_pamieci IS NULL OR
                                           zasieg_pamieci IN ('card','project','team')),
    temperatura              REAL,
    zaktualizowano           INTEGER NOT NULL
);
CREATE INDEX idx_profil_silnika_zasieg ON profil_silnika_tlumaczenia(zasieg, zasieg_id);

CREATE TABLE kanal_profilu_silnika (
    id        INTEGER PRIMARY KEY AUTOINCREMENT,
    profil_id INTEGER NOT NULL REFERENCES profil_silnika_tlumaczenia(id) ON DELETE CASCADE,
    kanal_kod TEXT    NOT NULL,
    kolejnosc INTEGER NOT NULL,
    UNIQUE(profil_id, kanal_kod)
);

-- ── Polityka tłumaczenia pivotowego ─────────────────────────────────────────
CREATE TABLE polityka_pivota (
    id              INTEGER PRIMARY KEY AUTOINCREMENT,
    zasieg          TEXT    NOT NULL DEFAULT 'global',
    zasieg_id       TEXT    NOT NULL DEFAULT '',
    jezyk_domyslny  TEXT,
    zaktualizowano  INTEGER NOT NULL,
    UNIQUE(zasieg, zasieg_id)
);

CREATE TABLE para_pivota (
    id           INTEGER PRIMARY KEY AUTOINCREMENT,
    polityka_id  INTEGER NOT NULL REFERENCES polityka_pivota(id) ON DELETE CASCADE,
    jezyk_zrodla TEXT    NOT NULL,
    jezyk_celu   TEXT    NOT NULL,
    jezyk_pivota TEXT    NOT NULL,
    UNIQUE(polityka_id, jezyk_zrodla, jezyk_celu)
);

-- ── Pakiet przekazania wykonawcy ────────────────────────────────────────────
CREATE TABLE pakiet_przekazania (
    id                       INTEGER PRIMARY KEY AUTOINCREMENT,
    identyfikator_zewnetrzny TEXT    NOT NULL UNIQUE,
    okno_id                  INTEGER NOT NULL REFERENCES okno_tlumaczenia(id) ON DELETE CASCADE,
    -- Zawartości pakietu i panele objęte pakietem, po jednym na wiersz kolumny.
    -- Byt bez własnej tożsamości: nikt nie adresuje pojedynczej zawartości.
    zawartosci               TEXT    NOT NULL,
    panele                   TEXT    NOT NULL,
    instrukcje               TEXT,
    sciezka                  TEXT    NOT NULL,
    stan                     TEXT    NOT NULL DEFAULT 'built'
                                     CHECK(stan IN ('built','sent','returned')),
    utworzono                INTEGER NOT NULL,
    zaktualizowano           INTEGER NOT NULL
);
CREATE INDEX idx_pakiet_przekazania_okno ON pakiet_przekazania(okno_id, utworzono DESC);

-- ── Przebieg pakietowy modułu Translate ─────────────────────────────────────
CREATE TABLE zlecenie_pakietu_tlumaczenia (
    id                       INTEGER PRIMARY KEY AUTOINCREMENT,
    identyfikator_zewnetrzny TEXT    NOT NULL UNIQUE,
    okno_id                  INTEGER NOT NULL REFERENCES okno_tlumaczenia(id) ON DELETE CASCADE,
    utworzono                INTEGER NOT NULL
);

CREATE TABLE pozycja_pakietu_tlumaczenia (
    id          INTEGER PRIMARY KEY AUTOINCREMENT,
    zlecenie_id INTEGER NOT NULL REFERENCES zlecenie_pakietu_tlumaczenia(id) ON DELETE CASCADE,
    panel_kod   TEXT    NOT NULL,
    operacja    TEXT    NOT NULL
                        CHECK(operacja IN ('translate','qualityCheck','proofread','export')),
    stan        TEXT    NOT NULL DEFAULT 'pending'
                        CHECK(stan IN ('pending','done','error')),
    szczegol    TEXT,
    zaktualizowano INTEGER NOT NULL
);
CREATE INDEX idx_pozycja_pakietu_zlecenie ON pozycja_pakietu_tlumaczenia(zlecenie_id, id);

-- ── Most do dokumentu ───────────────────────────────────────────────────────
CREATE TABLE most_tlumaczenia (
    id            INTEGER PRIMARY KEY AUTOINCREMENT,
    okno_id       INTEGER NOT NULL UNIQUE REFERENCES okno_tlumaczenia(id) ON DELETE CASCADE,
    dokument_kod  TEXT    NOT NULL,
    zakotwiczenie TEXT,
    zaktualizowano INTEGER NOT NULL
);
