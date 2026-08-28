-- Migracja 336 — szablony materiału (marketing: format społecznościowy, baner,
-- slajd, materiał do druku, nagłówek wiadomości).
--
-- Szablon materiału to NIE szablon promptu (`szablon_promptu_design`,
-- migracja 231). Tamten jest gotowym poleceniem dla silnika obrazu, ten jest
-- gotowym UKŁADEM: rozmiarem materiału i kompletem warstw, z którego
-- `design.template.apply` zakłada kompozycję. Wspólna tabela wymagałaby kolumn
-- pustych po obu stronach — prompt nie ma szerokości, układ nie ma tematu.
--
-- Warstwy leżą osobno, tak samo jak warstwy kompozycji
-- (`warstwa_kompozycji_design`, migracja 048): szablon bez warstw jest stanem
-- poprawnym (sam rozmiar materiału), a zapis podmienia komplet warstw naraz.
--
-- `zasob_id` jest tekstem, nie więzem obcym — dokładnie z tego powodu, co
-- w warstwie kompozycji: szablon przeżywa usunięcie zasobu z Assets Panel,
-- a wyrys pomija warstwę bez bajtów zamiast odmawiać całości.

CREATE TABLE szablon_materialu_design (
    id                       INTEGER PRIMARY KEY AUTOINCREMENT,
    identyfikator_zewnetrzny TEXT    NOT NULL UNIQUE,
    okno                     TEXT    NOT NULL,
    nazwa                    TEXT    NOT NULL,
    rodzaj                   TEXT    NOT NULL
                                     CHECK(rodzaj IN ('social','banner','presentation',
                                                      'print','email')),
    szerokosc                REAL    NOT NULL,
    wysokosc                 REAL    NOT NULL,
    opis                     TEXT,
    zaktualizowano           TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now'))
);

CREATE TABLE warstwa_szablonu_materialu_design (
    id                       INTEGER PRIMARY KEY AUTOINCREMENT,
    identyfikator_zewnetrzny TEXT    NOT NULL,
    szablon_id               INTEGER NOT NULL
                                     REFERENCES szablon_materialu_design(id) ON DELETE CASCADE,
    zasob_id                 TEXT,
    x                        REAL,
    y                        REAL,
    szerokosc                REAL,
    wysokosc                 REAL,
    kolejnosc                INTEGER NOT NULL DEFAULT 0,
    zablokowana              INTEGER NOT NULL DEFAULT 0 CHECK(zablokowana IN (0,1)),
    adnotacja                TEXT
);

-- design.template.list czyta szablony okna, od ostatnio zmienianego.
CREATE INDEX idx_szablon_materialu_design_okno
    ON szablon_materialu_design(okno, zaktualizowano DESC, id DESC);

-- Warstwy czyta się zawsze kompletem jednego szablonu, w kolejności wyrysu.
CREATE INDEX idx_warstwa_szablonu_materialu_design_szablon
    ON warstwa_szablonu_materialu_design(szablon_id, kolejnosc, id);
