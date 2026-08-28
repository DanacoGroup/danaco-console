-- Migracja 128 — dobudowa modułu Studio: byty, których dotąd nie było.
--
-- Moduł miał w kontrakcie sześć komend i trzy tabele: dokument, wersję,
-- propozycję zmiany. Opracowanie modułu opisuje osiem okien i dwanaście rodzin
-- funkcji, więc dobudowa wnosi byty, na których te funkcje stoją: komentarze
-- i adnotacje, zmiany śledzone, gałęzie dokumentu, kolejkę cyfryzacji, operacje
-- własne, łańcuchy, profile wydania, szablony i odwołania do wersji.
--
-- ── Dlaczego komentarz i adnotacja są jedną tabelą ───────────────────────────
-- Komentarz redakcyjny wisi przy zakresie znaków treści, adnotacja różnicy przy
-- numerze fragmentu porównania. Poza tym jednym polem oba byty mają te same
-- kolumny: dokument, wątek nadrzędny, autor, treść, stan rozwiązania, czas.
-- Dwie tabele o sześciu wspólnych kolumnach rozjechałyby się przy pierwszej
-- poprawce, a zapytanie „pokaż wszystko, co ktoś napisał przy tym dokumencie"
-- wymagałoby sumy dwóch zapytań.
--
-- ── Dlaczego autor jest kolumną wersji, a nie osobną tabelą ──────────────────
-- Rozdział 3.6 opracowania żąda przy każdej wersji rozróżnienia zmiany
-- Operatora od zmiany modelu. Autor jest cechą wersji, nie bytem samodzielnym:
-- nie ma stanu, nie ma historii i nie istnieje bez wersji. Kolumna dopuszcza
-- NULL, więc wersje założone przed tą migracją pozostają poprawne i czytają się
-- jako autor nieznany — a nie jako Operator, którym mogły nie być.
--
-- ── Dlaczego kroki łańcucha i pola szablonu są tekstem ───────────────────────
-- Krok łańcucha i pole szablonu nie są wyszukiwane, nie mają własnego cyklu
-- życia i nie wiąże się z nimi nic z zewnątrz — istnieją wyłącznie jako
-- zawartość swojego rodzica. Tabela podrzędna dałaby złączenie przy każdym
-- odczycie i nie dałaby w zamian niczego.

-- ── Komentarze redakcyjne i adnotacje przy różnicach ─────────────────────────
CREATE TABLE komentarz_studio (
    id                       INTEGER PRIMARY KEY AUTOINCREMENT,
    identyfikator_zewnetrzny TEXT    NOT NULL UNIQUE,
    dokument_id              INTEGER NOT NULL REFERENCES dokument_studio(id) ON DELETE CASCADE,
    rodzaj                   TEXT    NOT NULL DEFAULT 'komentarz'
                                     CHECK(rodzaj IN ('komentarz','adnotacja')),
    wersja_id                TEXT,
    watek_nadrzedny_id       TEXT,
    autor                    TEXT    NOT NULL DEFAULT 'uzytkownik'
                                     CHECK(autor IN ('uzytkownik','model')),
    -- Zakres znaków niesie komentarz redakcyjny; numer fragmentu — adnotacja
    -- różnicy. Wypełnione jest jedno albo drugie, nigdy oba naraz.
    zakres_od                INTEGER,
    zakres_do                INTEGER,
    fragment_numer           INTEGER,
    wersja_odniesienia_id    TEXT,
    wersja_porownywana_id    TEXT,
    propozycja_id            TEXT,
    tresc                    TEXT    NOT NULL,
    rozwiazany               INTEGER NOT NULL DEFAULT 0 CHECK(rozwiazany IN (0,1)),
    utworzono                TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now'))
);
CREATE INDEX idx_komentarz_studio_dokument ON komentarz_studio(dokument_id, rodzaj, id);

-- ── Zmiany zarejestrowane przez śledzenie zmian ──────────────────────────────
CREATE TABLE zmiana_sledzona_studio (
    id                       INTEGER PRIMARY KEY AUTOINCREMENT,
    identyfikator_zewnetrzny TEXT    NOT NULL UNIQUE,
    dokument_id              INTEGER NOT NULL REFERENCES dokument_studio(id) ON DELETE CASCADE,
    rodzaj                   TEXT    NOT NULL
                                     CHECK(rodzaj IN ('wstawienie','usuniecie','formatowanie')),
    autor                    TEXT    NOT NULL DEFAULT 'uzytkownik'
                                     CHECK(autor IN ('uzytkownik','model')),
    zakres_od                INTEGER NOT NULL,
    zakres_do                INTEGER NOT NULL,
    tresc_przed              TEXT,
    tresc_po                 TEXT,
    decyzja                  TEXT    NOT NULL DEFAULT 'oczekuje'
                                     CHECK(decyzja IN ('oczekuje','przyjeta','odrzucona')),
    utworzono                TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now'))
);
CREATE INDEX idx_zmiana_sledzona_studio_dokument ON zmiana_sledzona_studio(dokument_id, decyzja, zakres_od);

-- ── Gałęzie dokumentu ────────────────────────────────────────────────────────
CREATE TABLE galaz_studio (
    id                       INTEGER PRIMARY KEY AUTOINCREMENT,
    identyfikator_zewnetrzny TEXT    NOT NULL UNIQUE,
    dokument_id              INTEGER NOT NULL REFERENCES dokument_studio(id) ON DELETE CASCADE,
    nazwa                    TEXT    NOT NULL,
    wersja_startowa_id       TEXT    NOT NULL,
    wersja_biezaca_id        TEXT,
    scalona                  INTEGER NOT NULL DEFAULT 0 CHECK(scalona IN (0,1)),
    utworzono                TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now'))
);
CREATE INDEX idx_galaz_studio_dokument ON galaz_studio(dokument_id, id);

-- ── Kolejka wczytywania i cyfryzacji ─────────────────────────────────────────
-- Warstwa słów rozpoznanych i bloki układu stoją tekstem w formacie JSON: są
-- odczytem z materiału, nie bytem samodzielnym, i giną razem z pozycją.
CREATE TABLE pozycja_wczytywania_studio (
    id                       INTEGER PRIMARY KEY AUTOINCREMENT,
    identyfikator_zewnetrzny TEXT    NOT NULL UNIQUE,
    okno                     TEXT    NOT NULL,
    sciezka_zrodlowa         TEXT,
    zasob_id                 TEXT,
    stan                     TEXT    NOT NULL DEFAULT 'oczekuje'
                                     CHECK(stan IN ('oczekuje','przetwarzanie','gotowa','ponowienie','odmowa')),
    tekst                    TEXT,
    stron                    INTEGER,
    uzyto_rozpoznania        INTEGER NOT NULL DEFAULT 0 CHECK(uzyto_rozpoznania IN (0,1)),
    pewnosc                  REAL,
    powod_odmowy             TEXT,
    slowa_json               TEXT,
    uklad_json               TEXT,
    nastawy_json             TEXT,
    utworzono                TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now'))
);
CREATE INDEX idx_pozycja_wczytywania_studio_okno ON pozycja_wczytywania_studio(okno, id);

-- ── Operacje własne Tools Panel ──────────────────────────────────────────────
CREATE TABLE operacja_studio (
    id                       INTEGER PRIMARY KEY AUTOINCREMENT,
    identyfikator_zewnetrzny TEXT    NOT NULL UNIQUE,
    nazwa                    TEXT    NOT NULL,
    kategoria                TEXT    NOT NULL,
    prompt                   TEXT    NOT NULL,
    zasieg                   TEXT    NOT NULL DEFAULT 'global',
    zasieg_id                TEXT,
    utworzono                TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now'))
);
CREATE INDEX idx_operacja_studio_zasieg ON operacja_studio(zasieg, zasieg_id, kategoria, nazwa);

-- ── Łańcuchy operacji ────────────────────────────────────────────────────────
CREATE TABLE lancuch_studio (
    id                       INTEGER PRIMARY KEY AUTOINCREMENT,
    identyfikator_zewnetrzny TEXT    NOT NULL UNIQUE,
    nazwa                    TEXT    NOT NULL,
    kroki_json               TEXT    NOT NULL,
    zasieg                   TEXT    NOT NULL DEFAULT 'project',
    zasieg_id                TEXT,
    utworzono                TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now'))
);
CREATE INDEX idx_lancuch_studio_zasieg ON lancuch_studio(zasieg, zasieg_id, nazwa);

-- ── Profile wydania ──────────────────────────────────────────────────────────
CREATE TABLE profil_wydania_studio (
    id                       INTEGER PRIMARY KEY AUTOINCREMENT,
    identyfikator_zewnetrzny TEXT    NOT NULL UNIQUE,
    nazwa                    TEXT    NOT NULL,
    format                   TEXT    NOT NULL,
    ustawienia_json          TEXT,
    zasieg                   TEXT    NOT NULL DEFAULT 'global',
    zasieg_id                TEXT,
    utworzono                TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now'))
);
CREATE INDEX idx_profil_wydania_studio_zasieg ON profil_wydania_studio(zasieg, zasieg_id, nazwa);

-- ── Szablony dokumentów ──────────────────────────────────────────────────────
CREATE TABLE szablon_studio (
    id                       INTEGER PRIMARY KEY AUTOINCREMENT,
    identyfikator_zewnetrzny TEXT    NOT NULL UNIQUE,
    nazwa                    TEXT    NOT NULL,
    opis                     TEXT,
    format                   TEXT    NOT NULL DEFAULT 'markdown'
                                     CHECK(format IN ('pdf','docx','txt','markdown')),
    tresc                    TEXT    NOT NULL,
    pola_json                TEXT,
    fabryczny                INTEGER NOT NULL DEFAULT 0 CHECK(fabryczny IN (0,1)),
    utworzono                TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now'))
);
CREATE INDEX idx_szablon_studio_nazwa ON szablon_studio(nazwa);

-- ── Odwołania do wersji ──────────────────────────────────────────────────────
CREATE TABLE odwolanie_wersji_studio (
    id                       INTEGER PRIMARY KEY AUTOINCREMENT,
    identyfikator_zewnetrzny TEXT    NOT NULL UNIQUE,
    wersja_id                TEXT    NOT NULL,
    zasieg                   TEXT    NOT NULL DEFAULT 'session',
    wygasa                   TEXT,
    utworzono                TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now'))
);
CREATE INDEX idx_odwolanie_wersji_studio_wersja ON odwolanie_wersji_studio(wersja_id);

-- ── Wersja dokumentu: autor i przynależność ──────────────────────────────────
-- Cztery kolumny dopuszczające NULL, więc historia zastana pozostaje poprawna
-- bez przepisywania danych.
ALTER TABLE wersja_dokumentu_studio ADD COLUMN autor TEXT
    CHECK(autor IS NULL OR autor IN ('uzytkownik','model'));
ALTER TABLE wersja_dokumentu_studio ADD COLUMN kamien_milowy INTEGER NOT NULL DEFAULT 0
    CHECK(kamien_milowy IN (0,1));
ALTER TABLE wersja_dokumentu_studio ADD COLUMN galaz_id TEXT;
ALTER TABLE wersja_dokumentu_studio ADD COLUMN propozycja_id TEXT;

-- ── Dokument: stan śledzenia zmian ───────────────────────────────────────────
ALTER TABLE dokument_studio ADD COLUMN sledzenie_zmian INTEGER NOT NULL DEFAULT 0
    CHECK(sledzenie_zmian IN (0,1));

-- ── Ingest/OCR Panel wchodzi do katalogu okien ───────────────────────────────
-- Ósme okno modułu z opracowania. Katalog rdzenia znał dla Studia pięć okien,
-- więc okno zbudowane w kliencie nie miało wiersza, do którego mogłoby się
-- odwołać. Kategoria „narzedzia" kończyła się na pozycji 15, więc 16 dokleja
-- się na końcu i nie przestawia żadnego wiersza zastanego.
INSERT INTO okno_operacyjne (kod, nazwa, rola, kategoria, kolejnosc)
VALUES ('ingest-ocr-panel', 'Ingest/OCR Panel', 'pomocnicze', 'narzedzia', 16)
ON CONFLICT(kod) DO NOTHING;

INSERT INTO okno_operacyjne_modul (okno_operacyjne_id, modul_id, kolejnosc)
SELECT o.id, m.id, 6
  FROM okno_operacyjne o
  JOIN modul m ON m.kod = 'studio'
 WHERE o.kod = 'ingest-ocr-panel'
ON CONFLICT(okno_operacyjne_id, modul_id) DO NOTHING;
