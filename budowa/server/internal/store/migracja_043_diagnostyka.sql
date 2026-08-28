-- Migracja 043 — trwałość modułu Diagnostics: dziennik (Logs Viewer), błędy
-- zgrupowane po odcisku (Errors Panel), analiza stanu systemu (Diagnostics
-- Center) i rekomendacje z niej wyprowadzone (Recommendations Panel).
--
-- Dziennik ma wiersze, a nie bufor w pamięci: Logs Viewer pokazuje także to, co
-- działo się przed restartem rdzenia, a bufor pamięci ginie razem z procesem.
--
-- Czas jest liczbą, a nie napisem ISO. Kontrakt niesie czas wpisu
-- (`LogEntry.timestamp`), zakres zapytania (`fromTime`, `toTime`) i czasy błędu
-- w milisekundach epoki. Kolumna INTEGER przenosi dokładnie wartość kontraktu,
-- bez przekładu w obie strony przy każdym porównaniu zakresu.

-- ── Wpis dziennika (Logs Viewer) ─────────────────────────────────────────────
-- `odcisk` jest kolumną, a nie wyrażeniem liczonym przy odczycie: deduplikacja
-- z licznikiem grupuje po nim setki tysięcy wierszy i musi mieć indeks, którego
-- nie da się założyć na wyrażeniu sklejającym trzy kolumny w locie.
CREATE TABLE diagnostyka_wpis (
    id         INTEGER PRIMARY KEY AUTOINCREMENT,
    kod        TEXT    NOT NULL UNIQUE,
    chwila     INTEGER NOT NULL,
    poziom     TEXT    NOT NULL DEFAULT 'info'
                       CHECK(poziom IN ('debug','info','warn','error')),
    zrodlo     TEXT,
    tresc      TEXT    NOT NULL,
    sesja_kod  TEXT,
    okno_kod   TEXT,
    proces_kod TEXT,
    odcisk     TEXT    NOT NULL
);
CREATE INDEX idx_diagnostyka_wpis_chwila ON diagnostyka_wpis(chwila DESC);
CREATE INDEX idx_diagnostyka_wpis_poziom ON diagnostyka_wpis(poziom, chwila DESC);
CREATE INDEX idx_diagnostyka_wpis_zrodlo ON diagnostyka_wpis(zrodlo, chwila DESC);
CREATE INDEX idx_diagnostyka_wpis_odcisk ON diagnostyka_wpis(odcisk, chwila DESC);

-- ── Błąd zgrupowany po odcisku (Errors Panel) ────────────────────────────────
-- Odcisk jest unikalny, wystąpienia są licznikiem. Ta sama odmowa powtórzona
-- wielokrotnie jest jednym błędem o wielu wystąpieniach, a nie wieloma wierszami
-- — inaczej panel pokazywałby wyłącznie ostatnią minutę pracy i gubił błąd
-- rzadki, lecz istotny.
--
-- Kod błędu kontraktu ma własną kolumnę, mimo że `DiagnosticError` go nie ma.
-- Kolumna jest źródłem, z którego kod wychodzi do klienta polem `context` —
-- dzięki temu da się po nim także zawęzić zapytanie, czego pole JSON by nie dało.
CREATE TABLE diagnostyka_blad (
    id          INTEGER PRIMARY KEY AUTOINCREMENT,
    kod         TEXT    NOT NULL UNIQUE,
    odcisk      TEXT    NOT NULL UNIQUE,
    tresc       TEXT    NOT NULL,
    zrodlo      TEXT,
    kod_bledu   TEXT    NOT NULL DEFAULT 'internal_error'
                        CHECK(kod_bledu IN ('validation_failed','not_found','not_authenticated',
                                            'permission_denied','conflict','channel_unavailable',
                                            'rate_limited','internal_error')),
    stan        TEXT    NOT NULL DEFAULT 'new'
                        CHECK(stan IN ('new','analyzing','resolved','ignored')),
    priorytet   TEXT    NOT NULL DEFAULT 'medium'
                        CHECK(priorytet IN ('low','medium','high','critical')),
    wystapienia INTEGER NOT NULL DEFAULT 1,
    notatka     TEXT,
    kontekst    TEXT,
    pierwsze    INTEGER NOT NULL,
    ostatnie    INTEGER NOT NULL
);
CREATE INDEX idx_diagnostyka_blad_ostatnie ON diagnostyka_blad(ostatnie DESC);
CREATE INDEX idx_diagnostyka_blad_stan ON diagnostyka_blad(stan, ostatnie DESC);
CREATE INDEX idx_diagnostyka_blad_priorytet ON diagnostyka_blad(priorytet, ostatnie DESC);

-- ── Analiza stanu systemu (Diagnostics Center) ───────────────────────────────
-- Analiza jest migawką, nie widokiem liczonym na bieżąco: wiersz przechowuje
-- wynik z chwili uruchomienia, bo bez niego „porównanie z migawką wcześniejszą"
-- z panelu akcji nie miałoby z czym porównywać.
CREATE TABLE diagnostyka_analiza (
    id            INTEGER PRIMARY KEY AUTOINCREMENT,
    kod           TEXT    NOT NULL UNIQUE,
    okno_kod      TEXT,
    zakres_od     INTEGER,
    zakres_do     INTEGER,
    podsumowanie  TEXT,
    bledy         TEXT,
    porownana_kod TEXT,
    utworzono     INTEGER NOT NULL
);
CREATE INDEX idx_diagnostyka_analiza_utworzono ON diagnostyka_analiza(utworzono DESC);

-- ── Rekomendacja wyprowadzona z analizy (Recommendations Panel) ──────────────
-- Klucz obcy do analizy jest twardy: rekomendacja bez analizy nie ma faktu,
-- z którego wynika, a taka rekomendacja jest zgadywaniem.
CREATE TABLE diagnostyka_rekomendacja (
    id          INTEGER PRIMARY KEY AUTOINCREMENT,
    kod         TEXT    NOT NULL UNIQUE,
    analiza_kod TEXT    NOT NULL REFERENCES diagnostyka_analiza(kod) ON DELETE CASCADE,
    blad_kod    TEXT,
    tytul       TEXT    NOT NULL,
    szczegol    TEXT,
    priorytet   TEXT    NOT NULL DEFAULT 'medium'
                        CHECK(priorytet IN ('low','medium','high','critical')),
    stan        TEXT    NOT NULL DEFAULT 'proposed'
                        CHECK(stan IN ('proposed','applied','rejected')),
    sciezka     TEXT,
    poprawka    TEXT,
    utworzono   INTEGER NOT NULL
);
CREATE INDEX idx_diagnostyka_rekomendacja_analiza ON diagnostyka_rekomendacja(analiza_kod, utworzono DESC);
CREATE INDEX idx_diagnostyka_rekomendacja_stan ON diagnostyka_rekomendacja(stan, utworzono DESC);
