-- Migracja dokłada kolumnę wiążącą okno komunikacji z kodem eksperta, bez klucza obcego,
-- jako nastawę procesu, nie przepustkę.

ALTER TABLE okno_komunikacji ADD COLUMN agent_kod TEXT;

-- Wyszukiwanie okien jednego eksperta, potrzebne przy kasowaniu eksperta, żeby podać
-- liczbę okien, których to dotknie.
CREATE INDEX IF NOT EXISTS idx_okno_agent
    ON okno_komunikacji(agent_kod)
 WHERE agent_kod IS NOT NULL;
