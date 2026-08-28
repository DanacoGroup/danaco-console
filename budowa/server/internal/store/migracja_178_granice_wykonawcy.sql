-- Migracja 178 — granice działania Wykonawcy na stronie
-- (`browser.executor.limits.set`, `browser.executor.limits.get`).
--
-- Granice są konfiguracją zasięgu, nie bytem okna: kontrakt niesie je z polem
-- `scope` (`ConfigScope`) i `scopeId`, a żądanie wskazuje okno ALBO kartę
-- sesji. Dlatego kluczem jest para (zasięg, wskazanie zasięgu), a nie samo
-- okno — dwa wiersze o tym samym zasięgu byłyby dwiema odpowiedziami na jedno
-- pytanie „ile kroków wolno Wykonawcy tutaj".
--
-- Domeny dozwolone i zablokowane idą kolumnami JSON, bo są wykazem wewnątrz
-- jednego ustawienia, a nie bytem wyszukiwanym osobno.
--
-- Wiersza domyślnego migracja nie zakłada. Brak wiersza znaczy „granice
-- domyślne rdzenia" i tak też odpowiada `browser.executor.limits.get` — wartość
-- domyślna należy do kodu, który ją stosuje, a nie do schematu; wiersz
-- zaszczepiony w migracji rozjechałby się z nią przy pierwszej zmianie.
CREATE TABLE granica_wykonawcy_przegladania (
    id                       INTEGER PRIMARY KEY AUTOINCREMENT,
    zasieg                   TEXT    NOT NULL,
    zasieg_id                TEXT    NOT NULL DEFAULT '',
    max_krokow               INTEGER NOT NULL,
    max_czas_sekund          INTEGER NOT NULL,
    domeny_dozwolone_json    TEXT,
    domeny_zablokowane_json  TEXT,
    potwierdzaj_wyslanie     INTEGER NOT NULL DEFAULT 1 CHECK(potwierdzaj_wyslanie IN (0,1)),
    zaktualizowano           TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now')),
    UNIQUE(zasieg, zasieg_id)
);
