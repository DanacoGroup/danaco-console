-- Migracja 161 — segmentacja modułu Translate: zestawy reguł podziału
-- (`translate.segmentation.rules.*`) i trwałe segmenty okna
-- (`translate.segment.merge`, `translate.segment.split`).
--
-- Migracja 053 stwierdzała, że segment nie jest bytem trwałym, bo jedyna
-- komenda, która go dotykała (`source.segment`), dzieliła tekst w locie.
-- Scalanie i podział segmentów zmieniają to wprost: po scaleniu dwóch zdań
-- podział wynikający z tekstu źródłowego przestaje być prawdą o segmentach
-- okna, a kolejne wywołanie musi zobaczyć wynik poprzedniego. Bez tabeli
-- `segment_okna_tlumaczenia` obie komendy oddawałyby wykaz, który znika razem
-- z odpowiedzią — czyli meldunek zamiast skutku.
--
-- Wiersze zakładane są leniwie: dopóki nikt nie scalał ani nie dzielił,
-- okno nie ma ani jednego wiersza i podział liczy się z tekstu źródłowego.
-- Pierwsze scalenie albo pierwszy podział utrwala cały bieżący wykaz, a potem
-- zmienia w nim jedną rzecz — inaczej numer segmentu w żądaniu wskazywałby na
-- inny segment niż ten, który Operator widział.
--
-- Zestaw reguł segmentacji jest odrębny od okna: to nastawa wielokrotnego
-- użytku (kontrakt: `SegmentationRuleset` z własnym identyfikatorem i językiem).
-- `srx` niesie treść pliku SRX podaną przez Operatora — standard branżowy,
-- którego rdzeń nie rozkłada na własne reguły; przechowywany w całości, żeby
-- eksport oddał to samo, co przyszło.

-- ── Zestaw reguł segmentacji ────────────────────────────────────────────────
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
