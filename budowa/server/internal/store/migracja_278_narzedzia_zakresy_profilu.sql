-- Migracja 278 — zakresy uprawnień i limity wywołań pozycji katalogu narzędzi
-- dla profilu asystenta (`tools.scope.list`, `tools.scope.set`).
--
-- Zakres jest nastawą zasięgu, nie bramką wbudowaną: stanem wyjściowym jest
-- pełny dostęp bez limitu, a wiersz powstaje dopiero wtedy, gdy Operator coś
-- przestawił. Brak wiersza znaczy więc „pozycja dostępna, bez granicy" i to
-- jest jedyne znaczenie braku.
--
-- Zapis, który nikt nie czyta przy wykonaniu, byłby gorszy niż jego brak:
-- Operator widziałby ograniczenie, którego nikt nie egzekwuje. Dlatego obok
-- zakresu stoi tabela zużycia — bez niej pole `callLimit` byłoby etykietą,
-- a `callsUsed` z kontraktu nie miałoby z czego powstać.
--
-- Zużycie liczymy wierszem na wywołanie, a nie licznikiem w kolumnie zakresu.
-- Limit obowiązuje w OKNIE CZASU (`okno_sekund`), więc licznik narastający
-- musiałby być zerowany przez coś, co wie, kiedy okno się przesunęło —
-- czyli przez zadanie w tle, którego przy jednoosobowym produkcie nie ma po co
-- zakładać. Wiersze ze znacznikiem czasu odpowiadają na pytanie „ile wywołań
-- w ostatnich N sekundach" jednym zapytaniem i nie wymagają niczego w tle.

CREATE TABLE narzedzie_zakres_profilu (
    id             INTEGER PRIMARY KEY AUTOINCREMENT,
    profil_id      INTEGER NOT NULL REFERENCES profil_asystenta(id) ON DELETE CASCADE,
    -- Nazwa PEŁNA pozycji katalogu, ze źródłem: `danaco:session.rename`,
    -- `anthropic-skills:skill-creator`. Nazwa skrócona nie jest jednoznaczna —
    -- ta sama może przyjść z dwóch źródeł.
    nazwa_pelna    TEXT    NOT NULL,
    dostepne       INTEGER NOT NULL DEFAULT 1 CHECK (dostepne IN (0, 1)),
    -- Czy wywołanie wymaga jawnego potwierdzenia Operatora.
    potwierdzenie  INTEGER NOT NULL DEFAULT 0 CHECK (potwierdzenie IN (0, 1)),
    -- Górna granica liczby wywołań w oknie czasu; zero znaczy bez granicy.
    limit_wywolan  INTEGER NOT NULL DEFAULT 0 CHECK (limit_wywolan >= 0),
    okno_sekund    INTEGER NOT NULL DEFAULT 3600 CHECK (okno_sekund > 0),
    -- Dopuszczalne wartości argumentów, po jednej w wierszu. Pusta znaczy bez
    -- zawężenia.
    dopuszczone    TEXT    NOT NULL DEFAULT '',
    uzasadnienie   TEXT    NOT NULL DEFAULT '',
    zaktualizowano TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now')),
    UNIQUE (profil_id, nazwa_pelna)
);

CREATE TABLE narzedzie_wywolanie (
    id          INTEGER PRIMARY KEY AUTOINCREMENT,
    profil_id   INTEGER NOT NULL REFERENCES profil_asystenta(id) ON DELETE CASCADE,
    nazwa_pelna TEXT    NOT NULL,
    -- Karta sesji, z której wywołanie przyszło. NULL znaczy wywołanie spoza
    -- karty — na przykład pracę własną rdzenia.
    sesja_kod   TEXT,
    wykonano    TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now'))
);

CREATE INDEX idx_narzedzie_wywolanie_okno
    ON narzedzie_wywolanie(profil_id, nazwa_pelna, wykonano DESC);
