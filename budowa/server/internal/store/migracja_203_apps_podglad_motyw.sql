-- Migracja 203 — moduł Apps, warsztaty frontendu i backendu: podgląd na żywo
-- oraz motyw produktu.
--
-- Podgląd ma jeden wiersz na okno, nie na warstwę. `apps.preview.stop` niesie
-- samo `windowId` — nie ma czym wskazać, którą z dwóch warstw zatrzymać — więc
-- w oknie stoi najwyżej jeden serwer podglądu naraz, a `apps.preview.start`
-- z inną warstwą przestawia ten sam wiersz. Kolumna `warstwa` mówi, co ten
-- serwer dziś pokazuje.
--
-- Adres podglądu jest adresem serwera, który naprawdę stoi. Rdzeń podnosi go
-- na pętli zwrotnej i oddaje pod nim treść plików warsztatu — kolumna trzyma
-- adres wydany przez system operacyjny przy nasłuchu, a nie adres wymyślony.
-- Po restarcie rdzenia serwer nie stoi, więc wiersz zostaje w stanie `stopped`:
-- odtwarzanie nasłuchów z bazy przy starcie obiecywałoby podgląd, którego nikt
-- nie zamawiał.
--
-- Motyw jest surowym JSON-em kontraktu (`apps.theme.set` przyjmuje pole `json`,
-- `apps.theme.get` oddaje je z powrotem). Rdzeń go nie rozkłada na zmienne
-- stylistyczne: kształt motywu należy do systemu wizualnego produktu, nie do
-- platformy, a rozłożenie go tutaj byłoby drugą prawdą o cudzym kształcie.

-- ── Serwer podglądu warstwy produktu ─────────────────────────────────────────
-- Wartości kolumny `stan` są wartościami kontraktu (AppPreviewStatus).
CREATE TABLE podglad_apps (
    okno       TEXT PRIMARY KEY,
    warstwa    TEXT NOT NULL CHECK(warstwa IN ('frontend','backend')),
    adres      TEXT NOT NULL DEFAULT '',
    stan       TEXT NOT NULL DEFAULT 'stopped'
                    CHECK(stan IN ('starting','running','failed','stopped')),
    -- Milisekundy epoki — `AppsPreviewStartResponse.startedAt` niesie je wprost.
    rozpoczeto INTEGER NOT NULL DEFAULT 0,
    zatrzymano INTEGER
);

-- ── Motyw produktu — edytor motywu i stylów Frontend Workspace ───────────────
CREATE TABLE motyw_apps (
    okno           TEXT PRIMARY KEY,
    tresc          TEXT NOT NULL DEFAULT '{}',
    zaktualizowano TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now'))
);
