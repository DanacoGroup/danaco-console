-- Migracja 046 — trwałość modułu Studio: dokument otwarty w Studio Editor,
-- historia wersji dokumentu (repozytorium sesji) i propozycja zmiany
-- wypracowana przez operację kontekstową Tools Panel.
--
-- Repozytorium i historia wersji są jednym bytem, nie dwoma. Kontrakt niesie
-- `studio.repository.list` / `studio.repository.restore` obok pola
-- `StudioDocument.versionId` i odpowiedzi `studio.document.save`, która oddaje
-- `StudioVersion` przy `createVersion=true` — ale wszystkie trzy mówią o tej
-- samej liście wersji dokumentu. Dwie tabele niosłyby dwie prawdy o jednej
-- historii, więc prawda jest jedna — `wersja_dokumentu_studio` — a `document.save`
-- z `createVersion=true` dopisuje do niej wiersz zamiast zakładać własną.
--
-- Propozycja zmiany nie jest wersją. `studio.contextual.op` oddaje `proposalId`,
-- a `studio.diff.compare` przyjmuje go zamiennie z `targetVersionId` — więc
-- propozycja i wersja są w kontrakcie dwoma różnymi rzeczami porównywalnymi z tym
-- samym punktem odniesienia. Propozycja to wynik operacji Tools Panel
-- (streszczenie, przepisanie, tłumaczenie), zanim Operator zdecyduje, czy trafi
-- do repozytorium — commitowanie jej robi `document.save`, nie `contextual.op`.
-- Dlatego ma własną tabelę `propozycja_zmiany_studio`, osobną od wersji.
--
-- Treść dzieli się na kolumnę `tresc` obok `tresc_odwolanie`: treść krótka
-- zostaje w bazie, treść obszerna (PDF/DOCX odczytany do tekstu, długi dokument)
-- trafia do pliku, a baza trzyma tylko odwołanie. Dokument Studio bywa plikiem
-- PDF/DOCX wczytanym z repozytorium Library, więc para `tresc`/`tresc_odwolanie`
-- jest właściwym wyborem.
--
-- Okno i plik repozytorium są kolumnami tekstowymi bez więzu obcego: `windowId`
-- jest identyfikatorem kontraktu okna modułu, nie wiersza `okno_komunikacji`,
-- a `libraryFileId` wskazuje plik modułu Library, który żyje w osobnym,
-- równolegle budowanym module — więz obcy do niego wiązałby kolejność migracji,
-- której ta migracja nie kontroluje.

-- ── Dokument otwarty w Studio Editor ─────────────────────────────────────────
-- Wartości kolumny `format` są wartościami kontraktu (StudioDocumentFormat),
-- nie ich tłumaczeniem.
CREATE TABLE dokument_studio (
    id                       INTEGER PRIMARY KEY AUTOINCREMENT,
    identyfikator_zewnetrzny TEXT    NOT NULL UNIQUE,
    okno                     TEXT    NOT NULL,
    tytul                    TEXT,
    format                   TEXT    NOT NULL DEFAULT 'txt'
                                     CHECK(format IN ('pdf','docx','txt','markdown')),
    tresc                    TEXT,
    tresc_odwolanie          TEXT,
    plik_repozytorium_id     TEXT,
    -- Wskazuje ostatnią wersję z `wersja_dokumentu_studio` identyfikatorem
    -- zewnętrznym, nie kluczem wewnętrznym: dokument istnieje i bez wersji
    -- (przed pierwszym zapisem z `createVersion=true`), więc więz obcy
    -- musiałby dopuszczać NULL i tak, a kolumna tekstowa unika zależności od
    -- kolejności wstawienia wiersza wersji względem dokumentu.
    wersja_biezaca_id        TEXT,
    utworzono                TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now')),
    zaktualizowano           TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now'))
);
CREATE INDEX idx_dokument_studio_okno ON dokument_studio(okno, id);

-- ── Wersja dokumentu — repozytorium sesji (Repository Panel) ────────────────
CREATE TABLE wersja_dokumentu_studio (
    id                       INTEGER PRIMARY KEY AUTOINCREMENT,
    identyfikator_zewnetrzny TEXT    NOT NULL UNIQUE,
    dokument_id              INTEGER NOT NULL REFERENCES dokument_studio(id) ON DELETE CASCADE,
    etykieta                 TEXT,
    podsumowanie             TEXT,
    skrot_tresci             TEXT,
    tresc                    TEXT,
    tresc_odwolanie          TEXT,
    utworzono                TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now'))
);
CREATE INDEX idx_wersja_dokumentu_studio_dokument ON wersja_dokumentu_studio(dokument_id, utworzono DESC, id DESC);

-- ── Propozycja zmiany — wynik operacji kontekstowej (Tools Panel) ───────────
-- `akcja_id` trzyma pozycję rejestru akcji, która wygenerowała
-- propozycję — bez niej ślad operacji kontekstowej byłby nie do odróżnienia
-- od zwykłej wersji przy przeglądzie repozytorium.
CREATE TABLE propozycja_zmiany_studio (
    id                       INTEGER PRIMARY KEY AUTOINCREMENT,
    identyfikator_zewnetrzny TEXT    NOT NULL UNIQUE,
    dokument_id              INTEGER NOT NULL REFERENCES dokument_studio(id) ON DELETE CASCADE,
    akcja_id                 TEXT    NOT NULL,
    wiadomosc_id             TEXT,
    tresc_wyniku             TEXT,
    tresc                    TEXT,
    tresc_odwolanie          TEXT,
    utworzono                TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now'))
);
CREATE INDEX idx_propozycja_zmiany_studio_dokument ON propozycja_zmiany_studio(dokument_id, utworzono DESC);
