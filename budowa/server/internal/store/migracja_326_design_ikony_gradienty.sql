-- Migracja 326 — obszar koloru, ikon i typografii modułu Design: ikona własna
-- rysowana na siatce (`design.icon.set`, `design.icon.generate`) oraz gradient
-- wypełnienia kompozycji (`design.color.gradient.set`).
--
-- Czego ta migracja NIE zakłada, i dlaczego:
--
--   · katalog ikon otwartoźródłowych nie ma tu tabeli. Katalog jest wniesiony
--     do rdzenia jako zasób wkompilowany w binarium (`ikony_katalogu/`), a nie
--     ściągany ani zasiewany do bazy: to samo źródło i ta sama treść na każdym
--     stanowisku, bez kroku zasiewu, który mógłby się nie wykonać. Baza trzyma
--     wyłącznie ikony WŁASNE, czyli te, których rdzeń nie ma skąd odtworzyć.
--   · paleta, wynik kontrastu, przeliczenie barwy i podgląd kroju nie mają
--     tabel, bo nie są bytami — są rachunkiem z danych żądania. Zapisywanie
--     wyniku rachunku dawałoby drugą prawdę o barwie obok samej barwy.
--
-- ── Ikona własna ──────────────────────────────────────────────────────────────
-- Ikona należy do okna modułu, tak jak zasób i kompozycja: `design.icon.set`
-- niesie `windowId` jako pole wymagane. Treść SVG leży w kolumnie, nie
-- w magazynie bloków: ikona jest tekstem rzędu kilkuset bajtów, a sprite składa
-- się z wielu ikon naraz — czytanie ich z osobnych plików byłoby tyloma
-- otwarciami pliku, ile symboli w pakiecie.
--
-- `zestaw` odróżnia ikony własne okna od ikon przyniesionych z katalogu rdzenia
-- i jest tym samym polem, co `DesignIcon.Set` kontraktu.
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

-- Wykaz ikon okna idzie od ostatnio zmienianej — porządek ten sam, co w wykazie
-- zasobów i kompozycji, żeby ikona dopiero poprawiona stała na górze.
CREATE INDEX idx_ikona_design_okno ON ikona_design(okno, zaktualizowano DESC, id DESC);

-- Nazwa ikony jest jedyna w obrębie okna. `design.icon.set` bez `iconId`
-- zakłada ikonę nową, a druga ikona o tej samej nazwie w tym samym oknie
-- czyniłaby wskazanie po nazwie niejednoznacznym — w sprite'cie SVG nazwa staje
-- się `id` symbolu, a dwa symbole o jednym `id` w jednym pliku to plik zepsuty.
CREATE UNIQUE INDEX idx_ikona_design_nazwa ON ikona_design(okno, nazwa);

-- Etykieta ikony — para (ikona, etykieta) bez surogatu, wzorem
-- `etykieta_zasobu_design` z migracji 048.
CREATE TABLE etykieta_ikony_design (
    ikona_id  INTEGER NOT NULL REFERENCES ikona_design(id) ON DELETE CASCADE,
    etykieta  TEXT    NOT NULL,
    PRIMARY KEY (ikona_id, etykieta)
);

-- ── Gradient wypełnienia ──────────────────────────────────────────────────────
-- Gradient wisi na kompozycji, a wskazanie ścieżki albo warstwy zawęża go do
-- jednego bytu na niej (`design.color.gradient.set`: `boardId` wymagane,
-- `pathId` i `layerId` opcjonalne). Oba wskazania są TEKSTEM, nie więzem obcym:
-- warstwa jest przepisywana w całości przy każdym `design.board.update`
-- (migracja 048), więc więz kasowałby gradient przy każdym zapisie planszy,
-- a ścieżka należy do obszaru wektorowego, który ma własny cykl życia.
--
-- Klucz jest trójką (kompozycja, ścieżka, warstwa). Wskazanie pominięte zapisuje
-- się pustym napisem, nie NULL-em: SQLite uznaje dwa NULL-e za różne w indeksie
-- UNIQUE, więc gradient całej kompozycji zakładany dwa razy powstałby dwa razy
-- zamiast nadpisać się raz.
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

-- Stopień gradientu ma własną tabelę, nie zapis strukturalny w kolumnie:
-- `DesignGradientStop` niesie położenie, barwę i krycie, a kolejność stopni
-- rozstrzyga o wyniku interpolacji — porządkowanie po polu wewnątrz JSON-a
-- odbierałoby bazie możliwość oddania stopni już ułożonych.
CREATE TABLE stopien_gradientu_design (
    gradient_id  INTEGER NOT NULL REFERENCES gradient_design(id) ON DELETE CASCADE,
    kolejnosc    INTEGER NOT NULL,
    polozenie    REAL    NOT NULL,
    barwa        TEXT    NOT NULL,
    krycie       REAL,
    PRIMARY KEY (gradient_id, kolejnosc)
);
