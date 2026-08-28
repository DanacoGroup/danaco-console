-- Migracja 369 dodaje tożsamość wykonawcy przy zmianie śledzonej i przy
-- komentarzu, kodem i nazwą agenta.

ALTER TABLE zmiana_sledzona_studio ADD COLUMN autor_agent_kod TEXT;
ALTER TABLE zmiana_sledzona_studio ADD COLUMN autor_agent_nazwa TEXT;
ALTER TABLE zmiana_sledzona_studio ADD COLUMN autor_agent_wersja TEXT;
ALTER TABLE zmiana_sledzona_studio ADD COLUMN autor_podagent_kod TEXT;

-- Kolumny niosą wycinek postaci sprzed i po zmianie, aby zmiana postaci była
-- widoczna w podświetleniu tak samo jak zmiana treści.
ALTER TABLE zmiana_sledzona_studio ADD COLUMN postac_przed_json TEXT;
ALTER TABLE zmiana_sledzona_studio ADD COLUMN postac_po_json TEXT;

-- Kolumna wiąże zmianę śledzoną z czynnością dziennika, która ją odłożyła,
-- łącząc dwie drogi cofania zmiany.
ALTER TABLE zmiana_sledzona_studio ADD COLUMN czynnosc_kod TEXT;

CREATE INDEX idx_zmiana_sledzona_studio_agent
    ON zmiana_sledzona_studio(dokument_id, autor_agent_kod, zakres_od);
-- Indeks podświetlenia zmian wykonawców pyta o wszystkie zmiany autora model
-- w całym dokumencie, licząc je do licznika.
CREATE INDEX idx_zmiana_sledzona_studio_autor
    ON zmiana_sledzona_studio(dokument_id, autor, decyzja, zakres_od);

ALTER TABLE komentarz_studio ADD COLUMN autor_agent_kod TEXT;
ALTER TABLE komentarz_studio ADD COLUMN autor_agent_nazwa TEXT;
ALTER TABLE komentarz_studio ADD COLUMN autor_agent_wersja TEXT;
ALTER TABLE komentarz_studio ADD COLUMN autor_podagent_kod TEXT;

CREATE INDEX idx_komentarz_studio_agent
    ON komentarz_studio(dokument_id, autor_agent_kod, rodzaj);
