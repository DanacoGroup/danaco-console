-- Migracja 194 — głosowanie nad stanowiskami: warianty, głosy i wynik agregacji.
--
-- Wynik agregacji NIE jest kolumną. Liczy się go z głosów przy każdym odczycie
-- (`roundtable.vote.get`), bo metoda agregacji jest własnością głosowania,
-- a głos może dojść po pierwszym odczycie. Kolumna z wynikiem byłaby drugą
-- prawdą, rozjeżdżającą się z pierwszą przy każdym kolejnym głosie.
--
-- Głos jest jeden na wyborcę. Powtórne oddanie zastępuje poprzedni (ON CONFLICT
-- w zapisie), bo zmiana zdania w trakcie otwartego głosowania jest czynnością
-- dozwoloną, a dwa głosy tej samej osoby liczone dwukrotnie nie są.
--
-- Kształt głosu zależy od metody, więc kolumny są trzy i wszystkie mogą być
-- puste: aprobata wypełnia `aprobaty`, metody rankingowe `ranking`, skala
-- punktowa i metoda kwadratowa `punkty_json`.

CREATE TABLE debata_glosowanie (
    id                       INTEGER PRIMARY KEY AUTOINCREMENT,
    identyfikator_zewnetrzny TEXT    NOT NULL UNIQUE,
    okno                     TEXT    NOT NULL,
    tura                     TEXT    NOT NULL DEFAULT '',
    metoda                   TEXT    NOT NULL
                                     CHECK(metoda IN ('approval','irv','schulze','score','quadratic')),
    stan                     TEXT    NOT NULL DEFAULT 'open'
                                     CHECK(stan IN ('open','closed','tied')),
    -- Próg zgody ujemny znaczy „bez progu".
    prog                     REAL    NOT NULL DEFAULT -1.0,
    -- Uprawnieni do głosu, rozdzieleni znakiem nowego wiersza; pusto znaczy
    -- „cały skład".
    uprawnieni               TEXT    NOT NULL DEFAULT '',
    rozpoczeto               TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now')),
    zamknieto                TEXT
);
CREATE INDEX idx_debata_glosowanie_okno ON debata_glosowanie(okno, id DESC);

CREATE TABLE debata_wariant (
    id                       INTEGER PRIMARY KEY AUTOINCREMENT,
    identyfikator_zewnetrzny TEXT    NOT NULL UNIQUE,
    glosowanie               TEXT    NOT NULL,
    etykieta                 TEXT    NOT NULL,
    wypowiedz                TEXT    NOT NULL DEFAULT '',
    kolejnosc                INTEGER NOT NULL DEFAULT 0
);
CREATE INDEX idx_debata_wariant_glosowanie ON debata_wariant(glosowanie, kolejnosc, id);

CREATE TABLE debata_glos (
    id                       INTEGER PRIMARY KEY AUTOINCREMENT,
    identyfikator_zewnetrzny TEXT    NOT NULL UNIQUE,
    glosowanie               TEXT    NOT NULL,
    wyborca                  TEXT    NOT NULL,
    aprobaty                 TEXT    NOT NULL DEFAULT '',
    ranking                  TEXT    NOT NULL DEFAULT '',
    punkty_json              TEXT    NOT NULL DEFAULT '',
    oddano                   TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now')),
    UNIQUE (glosowanie, wyborca)
);
CREATE INDEX idx_debata_glos_glosowanie ON debata_glos(glosowanie, id);
