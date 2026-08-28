-- Migracja 278 zakłada tabele zakresów uprawnień i zużycia limitu wywołań
-- pozycji katalogu narzędzi dla profilu asystenta.

CREATE TABLE narzedzie_zakres_profilu (
    id             INTEGER PRIMARY KEY AUTOINCREMENT,
    profil_id      INTEGER NOT NULL REFERENCES profil_asystenta(id) ON DELETE CASCADE,
    -- Pełna nazwa pozycji katalogu ze źródłem; nazwa skrócona nie jest jednoznaczna.
    nazwa_pelna    TEXT    NOT NULL,
    dostepne       INTEGER NOT NULL DEFAULT 1 CHECK (dostepne IN (0, 1)),
    -- Określa, czy wywołanie wymaga jawnego potwierdzenia operatora.
    potwierdzenie  INTEGER NOT NULL DEFAULT 0 CHECK (potwierdzenie IN (0, 1)),
    -- Górna granica liczby wywołań w oknie czasu; zero oznacza brak granicy.
    limit_wywolan  INTEGER NOT NULL DEFAULT 0 CHECK (limit_wywolan >= 0),
    okno_sekund    INTEGER NOT NULL DEFAULT 3600 CHECK (okno_sekund > 0),
    -- Dopuszczalne wartości argumentów, po jednej w wierszu; pusta oznacza brak zawężenia.
    dopuszczone    TEXT    NOT NULL DEFAULT '',
    uzasadnienie   TEXT    NOT NULL DEFAULT '',
    zaktualizowano TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now')),
    UNIQUE (profil_id, nazwa_pelna)
);

CREATE TABLE narzedzie_wywolanie (
    id          INTEGER PRIMARY KEY AUTOINCREMENT,
    profil_id   INTEGER NOT NULL REFERENCES profil_asystenta(id) ON DELETE CASCADE,
    nazwa_pelna TEXT    NOT NULL,
    -- Karta sesji, z której wywołanie przyszło; pusta wartość oznacza pracę własną rdzenia.
    sesja_kod   TEXT,
    wykonano    TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now'))
);

CREATE INDEX idx_narzedzie_wywolanie_okno
    ON narzedzie_wywolanie(profil_id, nazwa_pelna, wykonano DESC);
