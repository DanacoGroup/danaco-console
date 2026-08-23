-- Wiązanie okna komunikacji z agentem.
--
-- Kolumna niesie `agent.kod`, nie `agent.id`, i nie ma więzu obcego — tak samo
-- jak `agent.kanal_kod`. Agent bywa kasowany niezależnie od okien, w których
-- pracował: `ON DELETE CASCADE` skasowałby wtedy okna wraz z całą rozmową,
-- a `ON DELETE SET NULL` po cichu przestawiłby okno na model surowy w środku
-- pracy. Kod, którego nikt już nie rozpoznaje, jest faktem czytelnym: okno
-- pamięta, kogo wskazano, a warstwa składająca nakładkę rusza z samą osią
-- modelu, gdy agenta nie znajdzie.
--
-- NULL znaczy model surowy — kontrakt opisuje pole jako „puste znaczy model
-- surowy", a pusty napis byłby drugim sposobem powiedzenia tego samego. Kolumna
-- jest więc bez NOT NULL i bez DEFAULT.
--
-- Wskazanie agenta niczego nie zamyka i niczego nie sprawdza: nie ma warunku
-- CHECK na istnienie agenta ani reguły „okno bez agenta nie może wysłać tury".
-- Agent jest nastawą procesu, nie przepustką.

ALTER TABLE okno_komunikacji ADD COLUMN agent_kod TEXT;

-- Wyszukiwanie okien jednego agenta: potrzebne przy kasowaniu agenta, żeby
-- podać liczbę okien, których to dotknie, zanim kasowanie ruszy.
-- Indeks częściowy — okna bez agenta są większością i nie mają czego wnosić.
CREATE INDEX IF NOT EXISTS idx_okno_agent
    ON okno_komunikacji(agent_kod)
 WHERE agent_kod IS NOT NULL;
