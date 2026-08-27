-- Migracja 367 dodaje do warsztatu szablonów pism postać wzorcową dokumentu
-- oraz osobny szereg wersji autozapisu.

ALTER TABLE szablon_studio ADD COLUMN kategoria TEXT;

-- Kolumna niesie postać wzorcową szablonu tym samym zapisem, którym jedzie
-- postać dokumentu, aby dała się podstawić bez przekładu.
ALTER TABLE szablon_studio ADD COLUMN postac_json TEXT;

-- Kolumna niesie odwołanie do zasobu miniatury podglądu w galerii szablonów,
-- a nie same bajty obrazka.
ALTER TABLE szablon_studio ADD COLUMN miniatura_zasob_kod TEXT;

-- Kolumna wskazuje, czy szablon powstał w warsztacie, czy został wniesiony
-- z pliku, co rozstrzyga format oddania pliku.
ALTER TABLE szablon_studio ADD COLUMN zrodlo_pliku TEXT;

-- Kolumna wskazuje dokument, z którego szablon powstał, aby było wiadomo,
-- gdzie stoi to pismo wzorcowe.
ALTER TABLE szablon_studio ADD COLUMN dokument_zrodlowy_kod TEXT;

ALTER TABLE szablon_studio ADD COLUMN zaktualizowano TEXT;

CREATE INDEX idx_szablon_studio_kategoria ON szablon_studio(kategoria, nazwa);

-- Kolumna szeregu odróżnia wersje nazwane od zapisów autozapisu, aby
-- autozapis nie zaśmiecał historii wersji dokumentu.
ALTER TABLE wersja_dokumentu_studio ADD COLUMN szereg TEXT NOT NULL DEFAULT 'operator'
    CHECK(szereg IN ('operator','autosave'));

-- Kolumna niesie postać dokumentu w chwili założenia wersji, potrzebną do
-- porównania i przywrócenia postaci.
ALTER TABLE wersja_dokumentu_studio ADD COLUMN postac_json TEXT;

-- Kolumny niosą tożsamość wykonawcy, który założył wersję dokumentu, kodem
-- i nazwą agenta modułu Agents.
ALTER TABLE wersja_dokumentu_studio ADD COLUMN autor_agent_kod TEXT;
ALTER TABLE wersja_dokumentu_studio ADD COLUMN autor_agent_nazwa TEXT;

CREATE INDEX idx_wersja_dokumentu_studio_szereg
    ON wersja_dokumentu_studio(dokument_id, szereg, utworzono DESC, id DESC);
