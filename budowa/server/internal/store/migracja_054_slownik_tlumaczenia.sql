-- Migracja 054 — trwałość fasety Translate „słownik": terminy glosariusza,
-- ślady importu i eksportu słownika oraz pamięć tłumaczeń (Glossary Panel,
-- Translation Memory).
--
-- Wystąpienia terminu (`translate.glossary.occurrences`) nie mają tabeli.
-- Wystąpienie to fragment tekstu w otoczeniu terminu, odtwarzalny z
-- `okno_tlumaczenia.tekst_zrodlowy`/`tekst_zrodlowy_odwolanie` i
-- `panel_tlumaczenia.tresc`/`tresc_odwolanie` przeszukanych względem
-- `termin_slownika.zrodlo`; osobne wiersze trzeba by unieważniać przy każdej
-- korekcie tłumaczenia panelu. Warstwa `dane` liczy je w locie.
--
-- Import i eksport słownika mają dwie wąskie tabele śladu zamiast jednej
-- wspólnej — to różne kierunki i różne kolumny wyniku. Liczba przetworzonych
-- terminów jest jedynym faktem o wyniku, którego nie da się odtworzyć bez
-- powtórzenia operacji. Żadna z tabel nie ma kolumny na treść ani na
-- identyfikator pliku repozytorium: rdzeń nie ma magazynu blobów, a
-- `glossary.export` pliku nie wytwarza.
--
-- Pamięć tłumaczeń jest tabelą gromadzoną z zatwierdzonych par segmentów, nie
-- odczytem po panelach. `translate.memory.suggest` dopasowuje przybliżenie
-- segmentu źródłowego względem wszystkich języków i paneli naraz, a podpowiedź
-- ma się opierać na treści zatwierdzonej, nie na bieżącej treści panelu
-- w trakcie edycji. Warstwa `dane` dopisuje wpis, gdy panel przechodzi w stan
-- `ready`.

-- ── Termin glosariusza — Glossary Panel ─────────────────────────────────────
-- Para (zrodlo, jezyk) nie jest kluczem, bo `translate.glossary.set` pozwala
-- edytować istniejący termin przez `TermId` — tożsamość terminu jest więc
-- własna (`identyfikator_zewnetrzny`), a nie wyprowadzona z treści.
CREATE TABLE termin_slownika (
    id                       INTEGER PRIMARY KEY AUTOINCREMENT,
    identyfikator_zewnetrzny TEXT    NOT NULL UNIQUE,
    zrodlo                   TEXT    NOT NULL,
    jezyk                    TEXT    NOT NULL,
    cel                      TEXT,
    nie_tlumaczyc            INTEGER NOT NULL DEFAULT 0 CHECK(nie_tlumaczyc IN (0,1)),
    uwaga                    TEXT,
    zaktualizowano           INTEGER NOT NULL
);
-- glossary.apply i glossary.occurrences szukają terminów danego języka po
-- treści źródłowej; indeks obsługuje obie ścieżki wyszukiwania.
CREATE INDEX idx_termin_slownika_jezyk_zrodlo ON termin_slownika(jezyk, zrodlo);

-- ── Ślad importu słownika — Glossary Panel (Import TBX/CSV) ────────────────
CREATE TABLE slad_importu_slownika (
    id                       INTEGER PRIMARY KEY AUTOINCREMENT,
    identyfikator_zewnetrzny TEXT    NOT NULL UNIQUE,
    sciezka                  TEXT    NOT NULL,
    liczba_zaimportowanych   INTEGER NOT NULL,
    utworzono                INTEGER NOT NULL
);
CREATE INDEX idx_slad_importu_slownika_czas ON slad_importu_slownika(utworzono DESC);

-- ── Ślad eksportu słownika — Glossary Panel (Export TBX/CSV) ────────────────
-- Bez kolumny na treść ani identyfikator pliku repozytorium — rdzeń nie ma
-- magazynu blobów i `glossary.export` nie wytwarza pliku. `sciezka` niesie
-- wyłącznie ścieżkę docelową zgłoszoną w żądaniu, nie dowód powstania pliku.
CREATE TABLE slad_eksportu_slownika (
    id                       INTEGER PRIMARY KEY AUTOINCREMENT,
    identyfikator_zewnetrzny TEXT    NOT NULL UNIQUE,
    sciezka                  TEXT    NOT NULL,
    liczba_wyeksportowanych  INTEGER NOT NULL,
    utworzono                INTEGER NOT NULL
);
CREATE INDEX idx_slad_eksportu_slownika_czas ON slad_eksportu_slownika(utworzono DESC);

-- ── Pamięć tłumaczeń — Translation Memory (`translate.memory.suggest`) ─────
-- Wpis gromadzony z zatwierdzonej pary segmentów, nie odczyt po panelu.
-- Segmenty to zwykle pojedyncze zdania (kontrakt: `Segment string`,
-- `Suggestions []string`), więc treść zostaje kolumną TEXT wprost, bez pary
-- `_odwolanie` — inaczej niż długa treść panelu czy okna, która bywa całym
-- dokumentem.
CREATE TABLE pamiec_tlumaczen (
    id                       INTEGER PRIMARY KEY AUTOINCREMENT,
    identyfikator_zewnetrzny TEXT    NOT NULL UNIQUE,
    panel_id                 INTEGER NOT NULL REFERENCES panel_tlumaczenia(id) ON DELETE CASCADE,
    jezyk                    TEXT    NOT NULL,
    segment_zrodlowy         TEXT    NOT NULL,
    segment_docelowy         TEXT    NOT NULL,
    utworzono                INTEGER NOT NULL
);
-- memory.suggest dopasowuje przybliżenie segmentu źródłowego w obrębie języka
-- docelowego panelu — indeks po (jezyk, segment_zrodlowy) obsługuje wstępne
-- zawężenie kandydatów przed dopasowaniem przybliżonym w warstwie `dane`.
CREATE INDEX idx_pamiec_tlumaczen_jezyk_segment ON pamiec_tlumaczen(jezyk, segment_zrodlowy);
CREATE INDEX idx_pamiec_tlumaczen_panel ON pamiec_tlumaczen(panel_id);
