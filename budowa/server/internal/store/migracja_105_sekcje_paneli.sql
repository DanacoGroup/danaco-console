-- Tworzy tabelę sekcja_panelu, przechowującą dla pary okna i panelu kolejność, zwinięcie i zdjęcie z widoku każdej sekcji, osobno dla każdego okna.

CREATE TABLE sekcja_panelu (
    id             INTEGER PRIMARY KEY AUTOINCREMENT,
    -- Identyfikator zewnętrzny okna — ten sam, którym posługuje się kontrakt.
    okno_id        TEXT    NOT NULL,
    -- Panel w obrębie okna; identyfikator nadaje widok.
    panel_id       TEXT    NOT NULL,
    -- Sekcja w obrębie panelu; identyfikator nadaje widok.
    sekcja_id      TEXT    NOT NULL,
    -- Miejsce na widoku, liczone od 1.
    kolejnosc      INTEGER NOT NULL CHECK(kolejnosc >= 1),
    -- Sekcja zawinięta do nagłówka, ale obecna na widoku.
    zwinieta       INTEGER NOT NULL DEFAULT 0 CHECK(zwinieta IN (0, 1)),
    -- Sekcja zdjęta z widoku w całości.
    zdjeta         INTEGER NOT NULL DEFAULT 0 CHECK(zdjeta IN (0, 1)),
    -- Chwila ostatniej zmiany układu w postaci ISO 8601 UTC, tym samym zegarem co okno_komunikacji.
    zaktualizowano TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now')),
    -- Sekcja występuje w panelu okna dokładnie raz, więc nie powstają dwa miejsca dla jednego bytu.
    UNIQUE(okno_id, panel_id, sekcja_id)
);

-- Odczyt układu jednego panelu w kolejności widoku — dokładnie to jedno pytanie
-- zadaje `panel.sections.get`. Kolumny w kolejności (okno, panel, kolejnosc)
-- czynią indeks pokrywającym porządek zapytania, więc sortowanie nie kosztuje.
CREATE INDEX idx_sekcja_panelu_uklad
    ON sekcja_panelu (okno_id, panel_id, kolejnosc);
