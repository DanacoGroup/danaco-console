-- Migracja 052 — trwałość pętli koordynator–wykonawca dla komend
-- `window.handoff` i `window.action`. Bez tej migracji więź koordynatora z oknem
-- wykonawczym oraz treść zlecenia, które ją uzasadnia, żyją wyłącznie w pamięci
-- przeglądarki i giną z jej zamknięciem — Mission Control po restarcie nie ma
-- czego pokazać.
--
-- Więź koordynator–wykonawca już ma dom. Kolumna
-- `okno_komunikacji.okno_koordynatora_id` jest dokładnym odpowiednikiem
-- `Window.coordinatorWindowId` z kontraktu. Ta migracja nie zakłada drugiej
-- tabeli na tę samą więź — `window.handoff` wyłącznie ją wypełnia (warstwa
-- `dane`, nie ten plik), gdy okno wykonawcy przyjmuje zlecenie. Czego tu
-- naprawdę brakuje: samego zlecenia — polecenia i kompletu kontekstu — oraz
-- jego powiązania z pozycją kolejki, która je niesie.
--
-- Komplet kontekstu jako zapis strukturalny, nie sześć tabel. `ContextBundle`
-- niesie sześć pól: prompt (skalar) oraz pięć list identyfikatorów obcych bytów
-- (dokumenty, agenci, historia, źródła wiedzy) plus json.RawMessage parametrów
-- wykonania. Żadne z tych pól nie jest spisem z własną logiką zapytań — to
-- zamrożony migawkowy odczyt „co Operator miał pod ręką w chwili przekazania”,
-- przenoszony jedną komendą i nieedytowalny później. Tabela na rodzaj odwołania
-- (pięć tabel złączeniowych) dałaby pięć pustych powtórzeń tego samego wzorca
-- bez korzyści zapytania — nic tu nie filtruje „wszystkie zlecenia z dokumentem
-- X”. Dlatego komplet kontekstu trafia do jednej kolumny `komplet_kontekstu`
-- (JSON kompletnego `ContextBundle`, zapisywany i odczytywany bez rozbioru
-- w SQL), tak samo jak `krok_automatyki` trzyma `parametry` i `warunek` jako
-- TEXT.
--
-- Wykonanie akcji zostawia trwały ślad. Dziennik `log_akcji_kolejki` już przyjął
-- zasadę „przejrzystość zamiast bramy” — dziennik zamiast blokowania akcji przez
-- brak zapisu. Akcje katalogu (`core/akcje.go`) potrafią zmieniać stan okna
-- i sesji, a ten sam problem zniknięcia po zamknięciu przeglądarki dotyczy ich
-- efektu tak samo jak więzi koordynatora. Dziennik `log_akcji_okna` odnotowuje
-- więc każde wykonanie — parametry, wynik, znacznik czasu — nie tylko te uznane
-- za istotne; wybór, co pokazać w interfejsie, należy do warstwy prezentacji.

-- ── Zlecenie przekazania — treść komendy window.handoff ────────────────────────
-- Jeden wiersz na jedno wywołanie window.handoff. `okno_zrodlowe_id` jest oknem
-- koordynatora w chwili wysyłki (może nim przestać być później — zlecenie to
-- migawka, nie odczyt na żywo). `pozycja_kolejki_id` wiąże zlecenie z jedynym
-- silnikiem kolejek — ta migracja go nie dubluje: pozycję kolejki zakłada
-- warstwa `dane` w tej samej transakcji, tu tylko przechowujemy odwołanie do
-- niej, które odpowiedź window.handoff zwraca jako `queueItemId`.
CREATE TABLE zlecenie_przekazania (
    id                     INTEGER PRIMARY KEY AUTOINCREMENT,
    sesja_id               INTEGER NOT NULL REFERENCES sesja(id) ON DELETE CASCADE,
    okno_zrodlowe_id       INTEGER NOT NULL REFERENCES okno_komunikacji(id) ON DELETE CASCADE,
    okno_docelowe_id       INTEGER NOT NULL REFERENCES okno_komunikacji(id) ON DELETE CASCADE,
    pozycja_kolejki_id     INTEGER NOT NULL REFERENCES pozycja_kolejki(id) ON DELETE CASCADE,
    polecenie              TEXT    NOT NULL,
    -- ContextBundle w całości, zapis JSON — patrz nagłówek pliku. Puste, gdy
    -- window.handoff nie niesie pola `bundle` (kontrakt: `omitempty`).
    komplet_kontekstu      TEXT,
    utworzono              TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now'))
);
CREATE INDEX idx_zlecenie_przekazania_docelowe ON zlecenie_przekazania(okno_docelowe_id, utworzono);
CREATE INDEX idx_zlecenie_przekazania_zrodlowe ON zlecenie_przekazania(okno_zrodlowe_id, utworzono);
CREATE INDEX idx_zlecenie_przekazania_pozycja ON zlecenie_przekazania(pozycja_kolejki_id);

-- ── Dziennik wykonania akcji — treść komendy window.action ──────────────────────
-- Katalog akcji mówi, jakie akcje istnieją; ta tabela mówi, kiedy i z jakim
-- skutkiem konkretne okno je wykonało. `akcja_id` odpowiada `Action.id`
-- katalogu (TEXT w kontrakcie — kod akcji, nie klucz obcy liczbowy), więc
-- kolumna jest TEXT, nie odwołaniem do innej tabeli tej migracji: katalog akcji
-- żyje poza schematem (`core/akcje.go`).
CREATE TABLE log_akcji_okna (
    id                     INTEGER PRIMARY KEY AUTOINCREMENT,
    okno_komunikacji_id    INTEGER NOT NULL REFERENCES okno_komunikacji(id) ON DELETE CASCADE,
    akcja_id               TEXT    NOT NULL,
    parametry              TEXT,
    wynik                  TEXT,
    utworzono              TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now'))
);
CREATE INDEX idx_log_akcji_okna_okno ON log_akcji_okna(okno_komunikacji_id, utworzono);
CREATE INDEX idx_log_akcji_okna_akcja ON log_akcji_okna(akcja_id);
