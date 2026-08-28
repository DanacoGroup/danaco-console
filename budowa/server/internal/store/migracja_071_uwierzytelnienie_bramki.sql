-- Migracja zakłada bramkę uwierzytelnienia Operatora: katalog metod wejścia oraz sesje
-- bramki, jako byt osobny od zastanych tabel.

CREATE TABLE metoda_uwierzytelnienia (
    id                       INTEGER PRIMARY KEY AUTOINCREMENT,
    identyfikator_zewnetrzny TEXT    NOT NULL UNIQUE,
    rodzaj                   TEXT    NOT NULL
                                     CHECK(rodzaj IN ('password','pin','hello')),
    etykieta                 TEXT,
    urzadzenie_kod           TEXT,
    nazwa_urzadzenia         TEXT,
    kotwica                  INTEGER NOT NULL DEFAULT 0 CHECK(kotwica IN (0,1)),
    sekret_odwolanie         TEXT    NOT NULL DEFAULT '',
    utworzono                INTEGER NOT NULL,
    ostatnio_uzyto           INTEGER,
    -- Kotwicą jest tylko hasło: otwiera bramkę z każdej maszyny, a PIN i klucz — tylko ze swojej.
    CHECK(kotwica = 0 OR rodzaj = 'password'),
    CHECK(rodzaj <> 'password' OR urzadzenie_kod IS NULL)
);

CREATE UNIQUE INDEX idx_metoda_uwierzytelnienia_kotwica
    ON metoda_uwierzytelnienia(kotwica) WHERE kotwica = 1;

CREATE UNIQUE INDEX idx_metoda_uwierzytelnienia_urzadzenie
    ON metoda_uwierzytelnienia(urzadzenie_kod, rodzaj) WHERE urzadzenie_kod IS NOT NULL;

CREATE INDEX idx_metoda_uwierzytelnienia_wykaz
    ON metoda_uwierzytelnienia(rodzaj, utworzono, id);

CREATE TABLE sesja_bramki (
    id             INTEGER PRIMARY KEY AUTOINCREMENT,
    token_skrot    TEXT    NOT NULL UNIQUE,
    metoda_rodzaj  TEXT    CHECK(metoda_rodzaj IS NULL
                                 OR metoda_rodzaj IN ('password','pin','hello')),
    urzadzenie_kod TEXT,
    wygasa         INTEGER NOT NULL,
    utworzono      INTEGER NOT NULL,
    uniewazniono   INTEGER
);

CREATE INDEX idx_sesja_bramki_czynne ON sesja_bramki(uniewazniono, wygasa);
