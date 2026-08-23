-- Migracja 049 — trwałość modułu Research: źródła badania (Sources Manager),
-- ustalenia (Findings), raport składany z ustaleń (Report Builder) wraz z
-- eksportem, oraz przestrzeń badania (Research Workspace).
--
-- `ResearchSource` nie jest `BrowserSource`. Migracja 047 (Browser) kataloguje
-- odcisk odwiedzonej strony — cykl życia zaczyna się od nawigacji przeglądarki
-- i kończy z kartą. `ResearchSource` tu zaczyna się od decyzji Operatora
-- „to źródło wchodzi do badania" i niesie własną, ręczną ocenę wiarygodności
-- (`ResearchCredibility`), której `BrowserSource` nie ma i mieć nie może — to
-- inna prawda o innym momencie. Stąd `zrodlo_badania` jest tabelą
-- osobną, bez więzu do tabeli źródeł modułu Browser.
--
-- `LibraryFileId` jest tekstem bez więzu obcego, po wzorze migracji 046
-- (`dokument_studio.plik_repozytorium_id`). Plik biblioteki żyje w module
-- Library (migracja 045), budowanym równolegle — więz obcy
-- do niego wiązałby kolejność migracji, której migracja Research nie
-- kontroluje, i wymagałby, żeby wiersz `plik_biblioteki` istniał już w chwili
-- katalogowania źródła. Kontrakt dopuszcza źródło bez pliku repozytorium
-- (pole opcjonalne — źródło może być samym adresem `Url` albo notatką), więc
-- kolumna tekstowa dopuszczająca NULL bez sztywnego więzu jest właściwym
-- wyborem, nie ustępstwem.
--
-- Treść ustalenia i sekcji raportu idzie przez parę `tresc`/`tresc_odwolanie`,
-- po wzorze migracji 046: `ResearchFinding.Content` i
-- `ResearchReportSection.Content` bywają krótką notatką albo długim akapitem
-- przeniesionym z dokumentu — para kolumn obsługuje oba przypadki bez osobnej
-- ścieżki dla „długiej" treści. Migracja 045 (Library) używa samego
-- `tresc_odwolanie`, bo tam treść zawsze jest plikiem (dokument repozytorium);
-- tu treść bywa krótkim zdaniem wpisanym wprost przez Operatora, więc krótka
-- ścieżka (kolumna `tresc`) musi zostać dostępna tak jak w Studio.
--
-- Eksport raportu jest bytem trwałym, nie czynnością bez śladu. Kontrakt
-- `research.report.export` oddaje `LibraryFileId`, `Path` i `SizeBytes` —
-- trzy fakty o wyniku, które muszą przeżyć samo wywołanie komendy, inaczej
-- Operator traci możliwość odpowiedzieć na pytanie „czy i dokąd ten raport
-- już wyeksportowałem" bez ponownego eksportu. Panel Report Builder pokazuje
-- historię eksportów obok wersji raportu, więc ślad ma własną tabelę
-- (`eksport_raportu_badania`), osobną od `raport_badania` — dwa eksporty tego
-- samego raportu do różnych formatów to dwa wiersze, nie nadpisanie jednego.
--
-- Przestrzeń badania jest bytem własnym bez okna, jak `automatyka` w migracji
-- 039. `ResearchWorkspaceSetRequest` nie niesie `windowId` ani identyfikatora
-- rozbudowywanego bytu — zakres i etapy badania są jedną, bieżącą definicją
-- całego modułu, nie stanem karty. Stąd `przestrzen_badania` jest tabelą
-- jednowierszową (klucz `id` wymuszony na 1), a `research.workspace.set`
-- nadpisuje ten jeden wiersz zamiast zakładać nowy — tak jak zapis definicji
-- automatyki nadpisuje jej pola, nie mnoży wierszy.
--
-- Związek ustalenie↔źródło i sekcja↔ustalenie mają własne tabele złącznikowe,
-- po wzorze `zaleznosc_kroku_automatyki` (migracja 039). `SourceIds` przy
-- ustaleniu i `FindingIds` przy sekcji raportu to relacje wiele-do-wielu —
-- jedno źródło zasila wiele ustaleń, jedno ustalenie zasila wiele sekcji.

-- ── Źródło badania — Sources Manager ───────────────────────────────────────
-- Wartości `rodzaj` i `wiarygodnosc` są wartościami kontraktu
-- (ResearchSourceKind, ResearchCredibility) wprost, bez tłumaczenia — wzór
-- `rodzaj` w migracji 039.
CREATE TABLE zrodlo_badania (
    id                       INTEGER PRIMARY KEY AUTOINCREMENT,
    identyfikator_zewnetrzny TEXT    NOT NULL UNIQUE,
    okno                     TEXT    NOT NULL,
    tytul                    TEXT    NOT NULL,
    rodzaj                   TEXT    NOT NULL DEFAULT 'web'
                                     CHECK(rodzaj IN ('web','document','note','dataset')),
    adres                    TEXT,
    pochodzenie              TEXT,
    wiarygodnosc             TEXT    NOT NULL DEFAULT 'unverified'
                                     CHECK(wiarygodnosc IN ('high','medium','low','unverified')),
    plik_biblioteki_id       TEXT,
    pozyskano_o              TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now'))
);
CREATE INDEX idx_zrodlo_badania_okno ON zrodlo_badania(okno, pozyskano_o DESC);

-- ── Ustalenie badania — Findings ────────────────────────────────────────────
CREATE TABLE ustalenie_badania (
    id                       INTEGER PRIMARY KEY AUTOINCREMENT,
    identyfikator_zewnetrzny TEXT    NOT NULL UNIQUE,
    okno                     TEXT    NOT NULL,
    tresc                    TEXT,
    tresc_odwolanie          TEXT,
    stan                     TEXT    NOT NULL DEFAULT 'open'
                                     CHECK(stan IN ('open','resolved')),
    utworzono                TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now')),
    zaktualizowano           TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now'))
);
CREATE INDEX idx_ustalenie_badania_okno ON ustalenie_badania(okno, zaktualizowano DESC);

-- ── Powiązanie ustalenia ze źródłami ────────────────────────────────────────
CREATE TABLE zrodlo_ustalenia_badania (
    ustalenie_id  INTEGER NOT NULL REFERENCES ustalenie_badania(id) ON DELETE CASCADE,
    zrodlo_id     INTEGER NOT NULL REFERENCES zrodlo_badania(id) ON DELETE CASCADE,
    utworzono     TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now')),
    PRIMARY KEY (ustalenie_id, zrodlo_id)
);
CREATE INDEX idx_zrodlo_ustalenia_badania_zrodlo ON zrodlo_ustalenia_badania(zrodlo_id, ustalenie_id);

-- ── Raport badania — Report Builder ─────────────────────────────────────────
CREATE TABLE raport_badania (
    id                       INTEGER PRIMARY KEY AUTOINCREMENT,
    identyfikator_zewnetrzny TEXT    NOT NULL UNIQUE,
    okno                     TEXT    NOT NULL,
    tytul                    TEXT    NOT NULL,
    utworzono                TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now')),
    zaktualizowano           TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now'))
);
CREATE INDEX idx_raport_badania_okno ON raport_badania(okno, zaktualizowano DESC);

-- ── Sekcja raportu ───────────────────────────────────────────────────────────
CREATE TABLE sekcja_raportu_badania (
    id                       INTEGER PRIMARY KEY AUTOINCREMENT,
    identyfikator_zewnetrzny TEXT    NOT NULL UNIQUE,
    raport_id                INTEGER NOT NULL REFERENCES raport_badania(id) ON DELETE CASCADE,
    tytul                    TEXT    NOT NULL,
    tresc                    TEXT,
    tresc_odwolanie          TEXT,
    kolejnosc                INTEGER NOT NULL DEFAULT 0,
    utworzono                TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now'))
);
CREATE INDEX idx_sekcja_raportu_badania_raport ON sekcja_raportu_badania(raport_id, kolejnosc, id);

-- ── Powiązanie sekcji raportu z ustaleniami ─────────────────────────────────
CREATE TABLE ustalenie_sekcji_raportu_badania (
    sekcja_id     INTEGER NOT NULL REFERENCES sekcja_raportu_badania(id) ON DELETE CASCADE,
    ustalenie_id  INTEGER NOT NULL REFERENCES ustalenie_badania(id) ON DELETE CASCADE,
    utworzono     TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now')),
    PRIMARY KEY (sekcja_id, ustalenie_id)
);
CREATE INDEX idx_ustalenie_sekcji_raportu_badania_ustalenie ON ustalenie_sekcji_raportu_badania(ustalenie_id, sekcja_id);

-- ── Eksport raportu — ślad wyjścia (Report Builder) ─────────────────────────
-- Wartości `format` są wartościami kontraktu (ExportFormat) wprost.
CREATE TABLE eksport_raportu_badania (
    id                       INTEGER PRIMARY KEY AUTOINCREMENT,
    identyfikator_zewnetrzny TEXT    NOT NULL UNIQUE,
    raport_id                INTEGER NOT NULL REFERENCES raport_badania(id) ON DELETE CASCADE,
    format                   TEXT    NOT NULL
                                     CHECK(format IN ('pdf','docx','markdown','html','txt')),
    sciezka_docelowa         TEXT,
    plik_biblioteki_id       TEXT,
    sciezka_wyniku           TEXT,
    rozmiar_bajtow           INTEGER,
    utworzono                TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now'))
);
CREATE INDEX idx_eksport_raportu_badania_raport ON eksport_raportu_badania(raport_id, utworzono DESC);

-- ── Przestrzeń badania — Research Workspace (byt jednowierszowy) ───────────
CREATE TABLE przestrzen_badania (
    id             INTEGER NOT NULL PRIMARY KEY CHECK(id = 1),
    zakres         TEXT    NOT NULL DEFAULT '',
    zaktualizowano TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now'))
);

-- Etapy badania — lista uporządkowana, osobna tabela bo `Stages []string` nie
-- ma tożsamości własnej (jak etykieta pliku w migracji 045), a kolejność się
-- liczy (kontrakt zwraca `Stages` w kolejności zapisu).
CREATE TABLE etap_przestrzeni_badania (
    kolejnosc  INTEGER NOT NULL,
    etap       TEXT    NOT NULL,
    PRIMARY KEY (kolejnosc)
);
