-- Migracja 316 — ścieżki wektorowe i symbole modułu Design (grupa C
-- opracowania: `design.vector.*`).
--
-- ── Dlaczego ścieżka jest bytem, a nie polem warstwy ─────────────────────────
-- Warstwa kompozycji (`warstwa_kompozycji_design`) opisuje prostokąt: położenie,
-- rozmiar, zasób. Krzywa Béziera nie mieści się w prostokącie, a operacja
-- logiczna na dwóch krzywych potrzebuje obu z osobna. Bez własnego bytu rysunek
-- piórem nie miałby gdzie zostać, a `design.vector.boolean` nie miałoby na czym
-- pracować.
--
-- ── Dlaczego węzły jadą jednym zapisem JSON, a nie wierszem na węzeł ─────────
-- Kontrakt (`DesignVectorPathSetRequest.Nodes`) nadsyła komplet węzłów przy
-- każdej zmianie: zmiana jednego węzła idzie tą samą drogą co narysowanie
-- ścieżki. Wiersz na węzeł dawałby drugi zapis tego samego kształtu, który
-- trzeba by kasować i wstawiać od nowa przy każdym drgnięciu myszy, a żadne
-- zapytanie nie pyta o pojedynczy węzeł. Zapis JSON jest więc tym samym
-- kompletem, którym jedzie kontrakt — jedną prawdą o kształcie.
--
-- ── Dlaczego warstwa wskazuje się kodem, a nie kluczem obcym ────────────────
-- `design.board.update` przepisuje komplet warstw kompozycji od nowa (usuń
-- i wstaw, `dane/design_kompozycje.go`), więc klucz wiersza warstwy zmienia się
-- przy każdym zapisie planszy. Klucz obcy do warstwy zrywałby się wtedy przy
-- zwykłym zapisie kompozycji. Identyfikator zewnętrzny warstwy przeżywa ten
-- zapis, bo klient nadsyła go z powrotem.

CREATE TABLE sciezka_wektorowa_design (
    id                       INTEGER PRIMARY KEY AUTOINCREMENT,
    identyfikator_zewnetrzny TEXT    NOT NULL UNIQUE,
    kompozycja_id            INTEGER NOT NULL REFERENCES kompozycja_design(id) ON DELETE CASCADE,
    warstwa_kod              TEXT,
    nazwa                    TEXT,
    wezly_json               TEXT    NOT NULL,
    zamknieta                INTEGER NOT NULL DEFAULT 0,
    wypelnienie_json         TEXT,
    obrys_json               TEXT,
    kolejnosc                INTEGER NOT NULL DEFAULT 0,
    utworzono                TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now')),
    zaktualizowano           TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now'))
);
CREATE INDEX idx_sciezka_wektorowa_design_kompozycja
    ON sciezka_wektorowa_design(kompozycja_id, kolejnosc, id);
CREATE INDEX idx_sciezka_wektorowa_design_warstwa
    ON sciezka_wektorowa_design(warstwa_kod);

-- ── Symbol i jego członkowie ────────────────────────────────────────────────
-- Symbol jest definicją wielokrotnego użycia złożoną ze ścieżek i warstw.
-- Członkostwo leży w osobnej tabeli, bo `design.vector.symbol.set` podmienia
-- komplet członków przy każdym zapisie, a liczba miejsc, do których zmiana
-- doszła (`propagatedTo`), jest liczbą wierszy naprawdę zapisanych — nie
-- obietnicą.
CREATE TABLE symbol_design (
    id                       INTEGER PRIMARY KEY AUTOINCREMENT,
    identyfikator_zewnetrzny TEXT    NOT NULL UNIQUE,
    kompozycja_id            INTEGER NOT NULL REFERENCES kompozycja_design(id) ON DELETE CASCADE,
    nazwa                    TEXT    NOT NULL,
    utworzono                TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now')),
    zaktualizowano           TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now'))
);
CREATE INDEX idx_symbol_design_kompozycja ON symbol_design(kompozycja_id, id);

CREATE TABLE czlonek_symbolu_design (
    id          INTEGER PRIMARY KEY AUTOINCREMENT,
    symbol_id   INTEGER NOT NULL REFERENCES symbol_design(id) ON DELETE CASCADE,
    rodzaj      TEXT    NOT NULL CHECK(rodzaj IN ('sciezka','warstwa')),
    czlonek_kod TEXT    NOT NULL,
    kolejnosc   INTEGER NOT NULL DEFAULT 0,
    UNIQUE(symbol_id, rodzaj, czlonek_kod)
);
CREATE INDEX idx_czlonek_symbolu_design_symbol ON czlonek_symbolu_design(symbol_id, kolejnosc, id);
