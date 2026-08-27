-- Migracja 326 dodaje ikony własne modułu Design rysowane na siatce, ich
-- etykiety oraz gradient wypełnienia ze stopniami.
CREATE TABLE ikona_design (
    id                       INTEGER PRIMARY KEY AUTOINCREMENT,
    identyfikator_zewnetrzny TEXT    NOT NULL UNIQUE,
    okno                     TEXT    NOT NULL,
    nazwa                    TEXT    NOT NULL,
    zestaw                   TEXT,
    svg                      TEXT    NOT NULL,
    siatka                   INTEGER,
    grubosc_obrysu           REAL,
    utworzono                TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now')),
    zaktualizowano           TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now'))
);

-- Wykaz ikon okna idzie od ostatnio zmienianej, tym samym porządkiem co
-- wykazy zasobów oraz kompozycji.
CREATE INDEX idx_ikona_design_okno ON ikona_design(okno, zaktualizowano DESC, id DESC);

-- Nazwa ikony jest jedyna w obrębie okna, ponieważ w pliku sprite nazwa
-- staje się identyfikatorem symbolu.
CREATE UNIQUE INDEX idx_ikona_design_nazwa ON ikona_design(okno, nazwa);

-- Etykieta ikony trzyma się parą ikony i etykiety bez osobnego
-- identyfikatora, wzorem etykiety zasobu.
CREATE TABLE etykieta_ikony_design (
    ikona_id  INTEGER NOT NULL REFERENCES ikona_design(id) ON DELETE CASCADE,
    etykieta  TEXT    NOT NULL,
    PRIMARY KEY (ikona_id, etykieta)
);

-- Gradient wypełnienia wisi na kompozycji, a wskazanie ścieżki lub warstwy
-- zawęża go do jednego bytu tekstem, nie więzem obcym.
CREATE TABLE gradient_design (
    id                       INTEGER PRIMARY KEY AUTOINCREMENT,
    identyfikator_zewnetrzny TEXT    NOT NULL UNIQUE,
    kompozycja_id            INTEGER NOT NULL REFERENCES kompozycja_design(id) ON DELETE CASCADE,
    sciezka_id               TEXT    NOT NULL DEFAULT '',
    warstwa_id               TEXT    NOT NULL DEFAULT '',
    rodzaj                   TEXT    NOT NULL
                                     CHECK(rodzaj IN ('linear','radial','conic')),
    kat                      REAL,
    zaktualizowano           TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now'))
);

CREATE UNIQUE INDEX idx_gradient_design_cel
    ON gradient_design(kompozycja_id, sciezka_id, warstwa_id);

-- Stopień gradientu ma własną tabelę, ponieważ kolejność stopni rozstrzyga
-- o wyniku interpolacji barwy.
CREATE TABLE stopien_gradientu_design (
    gradient_id  INTEGER NOT NULL REFERENCES gradient_design(id) ON DELETE CASCADE,
    kolejnosc    INTEGER NOT NULL,
    polozenie    REAL    NOT NULL,
    barwa        TEXT    NOT NULL,
    krycie       REAL,
    PRIMARY KEY (gradient_id, kolejnosc)
);
