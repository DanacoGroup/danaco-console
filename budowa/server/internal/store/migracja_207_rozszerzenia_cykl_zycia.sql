-- Migracja 207 — rodzina `extension.*`, cykl życia pozycji katalogu: kolekcje
-- kuratorskie, dziennik cyklu życia, wersje pozycji wraz z przypięciem oraz
-- paczki przesłane instalacją Personal.
--
-- Migracja 070 dała katalogowi jeden wiersz na pozycję i nic poza nim. Wszystko,
-- co rodzina `extension.*` robi z pozycją w czasie — kolekcjonuje ją, odnotowuje
-- zmiany, przypina wersję, cofa do wcześniejszej, przyjmuje przesłaną paczkę —
-- nie miało dotąd gdzie usiąść. Cztery tabele niżej są tymi miejscami.
--
-- KOLEKCJA JEST NAZWANYM ZESTAWEM, NIE ETYKIETĄ POZYCJI. `extension.collection.save`
-- nadsyła `extensionIds` w komplecie przy każdym zapisie, a
-- `extension.collection.apply` włącza albo wyłącza cały zestaw jednym
-- wywołaniem — więc związek ma tabelę złącznikową wymienianą „usuń, wstaw od
-- nowa", a nie kolumnę listy w wierszu kolekcji.
--
-- DZIENNIK CYKLU ŻYCIA JEST DZIENNIKIEM, NIE STANEM. `ExtensionHistoryEntry`
-- niesie czynność, wersję przed i po oraz czas — wiersz na zdarzenie, nigdy
-- nadpisywany. Wartości kolumny `czynnosc` są wartościami kontraktu
-- (ExtensionLifecycleAction).
--
-- WERSJA POZYCJI MA WIERSZ, BO INACZEJ COFNIĘCIE NIE MA DOKĄD WRÓCIĆ.
-- `extension.version.rollback` przyjmuje `targetVersion` i ma przywrócić stan
-- tamtej wersji; pozycja z jedną kolumną `wersja` pamięta wyłącznie tę bieżącą.
-- Wiersz wersji trzyma numer, dziennik zmian i odwołanie do paczki, z której
-- wersja powstała — to wystarcza, żeby cofnięcie było przywróceniem, a nie
-- przepisaniem napisu.
--
-- PRZYPIĘCIE JEST KOLUMNĄ POZYCJI, NIE WIERSZEM WERSJI. `extension.version.pin`
-- przypina JEDNĄ wersję pozycji, a wersja przypięta w dwóch wierszach naraz
-- byłaby sprzecznością, której nikt by nie wykrył. Kolumna `wersja_przypieta`
-- w tabeli `rozszerzenie` niesie tę jedną wartość; pusta znaczy „bez przypięcia".
--
-- PACZKA PRZESŁANA LEŻY NA DYSKU. `extension.package.upload` przyjmuje bajty
-- i oddaje `uploadRef`, którym woła się potem instalację. Kolumna `sciezka`
-- trzyma odwołanie względne magazynu treści rdzenia; wiersz bez pliku byłby
-- meldunkiem o przesyłce, której nie ma.

-- ── Kolekcja kuratorska ──────────────────────────────────────────────────────
CREATE TABLE kolekcja_rozszerzen (
    id                       INTEGER PRIMARY KEY AUTOINCREMENT,
    identyfikator_zewnetrzny TEXT    NOT NULL UNIQUE,
    nazwa                    TEXT    NOT NULL,
    opis                     TEXT,
    -- Oznaczenie barwne kolekcji; wartość podana przez Operatora, nie żeton
    -- systemu wizualnego — rdzeń jej nie interpretuje.
    oznaczenie_barwne        TEXT,
    zaktualizowano           INTEGER NOT NULL DEFAULT 0
);

CREATE TABLE pozycja_kolekcji_rozszerzen (
    kolekcja_id      INTEGER NOT NULL REFERENCES kolekcja_rozszerzen(id) ON DELETE CASCADE,
    -- `Extension.id` pozycji; wartość danych, nie więz obcy: kolekcja ma prawo
    -- wskazywać pozycję odinstalowaną, bo odinstalowanie nie kasuje wiersza.
    rozszerzenie_kod TEXT    NOT NULL,
    kolejnosc        INTEGER NOT NULL DEFAULT 0,
    PRIMARY KEY (kolekcja_id, rozszerzenie_kod)
);

-- ── Dziennik cyklu życia pozycji ─────────────────────────────────────────────
CREATE TABLE historia_rozszerzenia (
    id                       INTEGER PRIMARY KEY AUTOINCREMENT,
    identyfikator_zewnetrzny TEXT    NOT NULL UNIQUE,
    rozszerzenie_kod         TEXT    NOT NULL,
    czynnosc                 TEXT    NOT NULL
                                     CHECK(czynnosc IN ('installed','updated','enabled','disabled',
                                                        'rolledBack','uninstalled','configured')),
    wersja_przed             TEXT,
    wersja_po                TEXT,
    szczegol                 TEXT,
    -- Milisekundy epoki — `ExtensionHistoryEntry.occurredAt` niesie je wprost.
    zaszlo                   INTEGER NOT NULL
);
CREATE INDEX idx_historia_rozszerzenia ON historia_rozszerzenia(rozszerzenie_kod, zaszlo DESC, id DESC);

-- ── Wersja pozycji katalogu ──────────────────────────────────────────────────
CREATE TABLE wersja_rozszerzenia (
    id               INTEGER PRIMARY KEY AUTOINCREMENT,
    rozszerzenie_kod TEXT    NOT NULL,
    wersja           TEXT    NOT NULL,
    dziennik_zmian   TEXT,
    -- Odwołanie do paczki, z której wersja powstała; puste dla wersji
    -- zarejestrowanej bez przesyłki.
    paczka_odwolanie TEXT,
    utworzono        INTEGER NOT NULL DEFAULT 0,
    UNIQUE (rozszerzenie_kod, wersja)
);

-- Przypięcie wersji — patrz rozstrzygnięcie na czole pliku.
ALTER TABLE rozszerzenie ADD COLUMN wersja_przypieta TEXT;

-- ── Paczka przesłana instalacją Personal ─────────────────────────────────────
CREATE TABLE paczka_rozszerzenia (
    id                       INTEGER PRIMARY KEY AUTOINCREMENT,
    -- `uploadRef` oddawany Operatorowi i przyjmowany z powrotem przy instalacji.
    identyfikator_zewnetrzny TEXT    NOT NULL UNIQUE,
    nazwa_pliku              TEXT    NOT NULL,
    -- Odwołanie względne magazynu treści rdzenia — patrz czoło pliku.
    sciezka                  TEXT    NOT NULL,
    rozmiar                  INTEGER NOT NULL DEFAULT 0,
    suma_kontrolna           TEXT    NOT NULL,
    utworzono                INTEGER NOT NULL DEFAULT 0
);
