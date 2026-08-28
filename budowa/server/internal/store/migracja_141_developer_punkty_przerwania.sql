-- Migracja 141 — punkty przerwania okna Run & Debug modułu Developer.
--
-- Punkt przerwania przeżywa sesję debugowania i musi ją przeżyć. Operator
-- stawia go w marginesie Code Editora zanim cokolwiek uruchomi, a potem
-- uruchamia debugowanie po raz drugi i trzeci — punkt postawiony w pamięci
-- sesji zniknąłby razem z nią i trzeba by go stawiać od nowa przy każdym biegu.
-- Dlatego punkt należy do okna i pliku, nie do sesji debugowania.
--
-- `zweryfikowany` mówi, czy adapter DAP potwierdził, że pod tym wierszem da się
-- zatrzymać. Punkt niezweryfikowany nie jest usterką: plik bywa jeszcze
-- nieskompilowany, a Operator ma prawo postawić punkt zanim program powstanie.
--
-- Warunek, warunek trafień i treść wpisu stoją osobnymi kolumnami, bo są
-- osobnymi rodzajami punktu (BreakpointKind kontraktu) i pytanie o nie zadaje
-- się osobno — sklejenie ich w jedno pole kazałoby rdzeniowi zgadywać, które
-- z trzech znaczeń niesie zapisany tekst.
CREATE TABLE developer_punkt_przerwania (
    id            INTEGER PRIMARY KEY AUTOINCREMENT,
    kod           TEXT    NOT NULL UNIQUE,
    okno_kod      TEXT    NOT NULL,
    sciezka       TEXT    NOT NULL,
    wiersz        INTEGER NOT NULL,
    rodzaj        TEXT    NOT NULL DEFAULT 'line'
                          CHECK(rodzaj IN ('line','conditional','logpoint','exception')),
    warunek       TEXT,
    warunek_trafien TEXT,
    wpis          TEXT,
    zweryfikowany INTEGER NOT NULL DEFAULT 0,
    utworzono     TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now')),
    UNIQUE(okno_kod, sciezka, wiersz)
);
CREATE INDEX idx_developer_punkt_okno ON developer_punkt_przerwania(okno_kod, sciezka, wiersz);
