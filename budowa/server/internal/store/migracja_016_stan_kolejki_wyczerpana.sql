-- Migracja 016 — stan kolejki `wyczerpana` w więzie CHECK.
--
-- POWÓD. Kontrakt zna PIĘĆ stanów kolejki: QueueStatus.idle · running · paused ·
-- stopped · done, a mapa przekładu `WartosciBazyQueueStatus` wiąże je kolejno
-- z kolumnami 'bezczynna' · 'pracuje' · 'wstrzymana' · 'zatrzymana' ·
-- 'wyczerpana'. Więz CHECK z `migracja_003_kolejki.sql` dopuszcza CZTERY — bez
-- 'wyczerpana', więc zapis stanu, który kontrakt zna, kończy się odmową
-- schematu. To usterka schematu, nie kontraktu: kolejka, w której nie została
-- ani jedna pozycja czynna, MUSI mieć jak się nazwać, inaczej silnik wykonania
-- nie ma stanu końcowego i po ostatniej pozycji kłamie, że nadal pracuje.
--
-- ZASIĘG PRZEBUDOWY. SQLite nie zna polecenia zdejmującego więz CHECK, więc
-- tabela `kolejka` idzie przez przebudowę (nowa tabela → przepisanie wierszy →
-- podmiana nazwy) tak samo jak `ustawienie` w `migracja_012_katalog_ustawien.sql`.
-- Różnica jest jedna
-- i przesądza o kształcie tego kroku: `kolejka` JEST wskazywana kluczem obcym —
-- przez `pozycja_kolejki` i przez `log_akcji_kolejki`, w obu razach z ON DELETE
-- CASCADE. Przy włączonym więzie kluczy obcych (`store/baza.go`, pragma
-- `foreign_keys(1)` w DSN każdego połączenia) polecenie DROP TABLE wykonuje
-- niejawne DELETE wszystkich wierszy, a to uruchamia kaskadę: zniknęłyby
-- wszystkie pozycje kolejek i cały dziennik akcji. Pragmy nie da się wyłączyć
-- na czas kroku — migracja biegnie w transakcji (`store/migracje.go`), a
-- `PRAGMA foreign_keys` jest w transakcji bez skutku. Dlatego przebudowa
-- obejmuje wszystkie trzy tabele obszaru, w porządku bezpiecznym dla danych:
--
--   1. kopie tabel z poprawionym więzem, wiersze przepisane co do kolumny;
--   2. skasowanie tabel starych od dziecka do rodzica — gdy pada `kolejka`,
--      nie wskazuje jej już żaden wiersz, więc kaskada nie ma czego zabrać;
--   3. podmiana nazw; ALTER TABLE ... RENAME przepisuje odwołania kluczy obcych
--      w tabelach pozostałych, więc dzieci wracają do rodzica pod nazwą własną;
--   4. odtworzenie indeksów pod nazwami z `migracja_003_kolejki.sql` (indeksy
--      tabeli skasowanej giną razem z nią, więc nazwy są wolne).
--
-- CZEGO TU NIE MA. Słownik stanów POZYCJI kolejki zostaje bez zmiany: siedem
-- wartości z `migracja_003_kolejki.sql` wystarcza silnikowi wykonania na pełny cykl życia
-- zlecenia (oczekuje → wykonywana → do_weryfikacji → ukonczona, z biegiem
-- naprawczym i przerwaniem). Kolumny nowej nie dokłada tu ani jedna linia —
-- migracja naprawia więz, a nie rozszerza model.

-- ── 1. Kopie tabel ───────────────────────────────────────────────────────────

CREATE TABLE kolejka_nowa (
    id                     INTEGER PRIMARY KEY AUTOINCREMENT,
    nazwa                  TEXT    NOT NULL,
    rodzaj                 TEXT    NOT NULL DEFAULT 'sesyjna'
                                   CHECK(rodzaj IN ('sesyjna','multitasking')),
    sesja_id               INTEGER REFERENCES sesja(id) ON DELETE CASCADE,
    okno_koordynatora_id   INTEGER REFERENCES okno_komunikacji(id) ON DELETE SET NULL,
    -- Piąta wartość odpowiada QueueStatus.done kontraktu: kolejka wyczerpana to
    -- kolejka, w której nie została ani jedna pozycja czynna.
    stan                   TEXT    NOT NULL DEFAULT 'bezczynna'
                                   CHECK(stan IN ('bezczynna','pracuje','wstrzymana',
                                                  'zatrzymana','wyczerpana')),
    utworzono              TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now')),
    zaktualizowano         TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now'))
);

INSERT INTO kolejka_nowa
    (id, nazwa, rodzaj, sesja_id, okno_koordynatora_id, stan, utworzono, zaktualizowano)
SELECT id, nazwa, rodzaj, sesja_id, okno_koordynatora_id, stan, utworzono, zaktualizowano
  FROM kolejka;

CREATE TABLE pozycja_kolejki_nowa (
    id                     INTEGER PRIMARY KEY AUTOINCREMENT,
    kolejka_id             INTEGER NOT NULL REFERENCES kolejka_nowa(id) ON DELETE CASCADE,
    okno_wykonawcy_id      INTEGER REFERENCES okno_komunikacji(id) ON DELETE SET NULL,
    kolejnosc              INTEGER NOT NULL DEFAULT 0,
    tytul                  TEXT    NOT NULL,
    tresc_zlecenia         TEXT,
    tresc_odwolanie        TEXT,
    stan                   TEXT    NOT NULL DEFAULT 'oczekuje'
                                   CHECK(stan IN ('oczekuje','przydzielona','wykonywana',
                                                  'do_weryfikacji','ukonczona','bledna','anulowana')),
    werdykt_weryfikacji    TEXT    CHECK(werdykt_weryfikacji IS NULL OR
                                         werdykt_weryfikacji IN ('przyjete','do_poprawy','odrzucone')),
    licznik_obiegow        INTEGER NOT NULL DEFAULT 0,
    utworzono              TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now')),
    zaktualizowano         TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now'))
);

INSERT INTO pozycja_kolejki_nowa
    (id, kolejka_id, okno_wykonawcy_id, kolejnosc, tytul, tresc_zlecenia, tresc_odwolanie,
     stan, werdykt_weryfikacji, licznik_obiegow, utworzono, zaktualizowano)
SELECT id, kolejka_id, okno_wykonawcy_id, kolejnosc, tytul, tresc_zlecenia, tresc_odwolanie,
       stan, werdykt_weryfikacji, licznik_obiegow, utworzono, zaktualizowano
  FROM pozycja_kolejki;

CREATE TABLE log_akcji_kolejki_nowa (
    id                     INTEGER PRIMARY KEY AUTOINCREMENT,
    kolejka_id             INTEGER NOT NULL REFERENCES kolejka_nowa(id) ON DELETE CASCADE,
    pozycja_kolejki_id     INTEGER REFERENCES pozycja_kolejki_nowa(id) ON DELETE CASCADE,
    akcja                  TEXT    NOT NULL,
    stan_przed             TEXT,
    stan_po                TEXT,
    numer_obiegu           INTEGER NOT NULL DEFAULT 0,
    szczegoly              TEXT,
    utworzono              TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now'))
);

INSERT INTO log_akcji_kolejki_nowa
    (id, kolejka_id, pozycja_kolejki_id, akcja, stan_przed, stan_po, numer_obiegu,
     szczegoly, utworzono)
SELECT id, kolejka_id, pozycja_kolejki_id, akcja, stan_przed, stan_po, numer_obiegu,
       szczegoly, utworzono
  FROM log_akcji_kolejki;

-- ── 2. Skasowanie tabel starych od dziecka do rodzica ────────────────────────

DROP TABLE log_akcji_kolejki;
DROP TABLE pozycja_kolejki;
DROP TABLE kolejka;

-- ── 3. Podmiana nazw ─────────────────────────────────────────────────────────

ALTER TABLE kolejka_nowa           RENAME TO kolejka;
ALTER TABLE pozycja_kolejki_nowa   RENAME TO pozycja_kolejki;
ALTER TABLE log_akcji_kolejki_nowa RENAME TO log_akcji_kolejki;

-- ── 4. Odtworzenie indeksów ──────────────────────────────────────────────────

CREATE INDEX idx_kolejka_sesja ON kolejka(sesja_id);
CREATE INDEX idx_kolejka_koordynator ON kolejka(okno_koordynatora_id);
CREATE INDEX idx_pozycja_kolejki_kolejka ON pozycja_kolejki(kolejka_id, kolejnosc, id);
CREATE INDEX idx_pozycja_kolejki_wykonawca ON pozycja_kolejki(okno_wykonawcy_id);
CREATE INDEX idx_log_akcji_kolejki_kolejka ON log_akcji_kolejki(kolejka_id, id);
CREATE INDEX idx_log_akcji_kolejki_pozycja ON log_akcji_kolejki(pozycja_kolejki_id);
