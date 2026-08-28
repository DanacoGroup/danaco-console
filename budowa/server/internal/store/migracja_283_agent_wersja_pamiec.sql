-- Migracja 283 — poziomy pamięci docierają do migawki wersji także wtedy, gdy
-- zapisano je PO tożsamości.
--
-- Migracja 281 wypełnia `agent_wersja.poziomy_pamieci` podzapytaniem
-- wykonywanym w chwili zakładania migawki. To za wcześnie w dwóch przypadkach,
-- które są normą, a nie wyjątkiem:
--
--   - przy zakładaniu eksperta wyzwalacz `agent_wersja_po_zalozeniu` biegnie
--     AFTER INSERT ON agent, a wiersze `agent_pamiec_poziom` wchodzą krok
--     później — migawka zapisywała pustkę, czyli „pamięć wyłączona", choć
--     wyjściowo obowiązują cztery poziomy;
--   - `agent.update` z samym polem `memoryLevels` nie rusza ani jednej kolumny
--     tabeli `agent`, więc migawki nie zakłada i nie odświeża.
--
-- Skutek byłby dokładnie tym, przed czym broni się `agent.version.get`:
-- podglądem wersji, który KŁAMIE o tym, co w wersji stało.
--
-- Poprawka idzie wyzwalaczami na tabeli poziomów. Każda zmiana składu odświeża
-- poziomy w migawce wersji BIEŻĄCEJ tego eksperta — tej, którą zmiana dotyczy.
-- Wersje wcześniejsze zostają nietknięte: one opisują stan sprzed zmiany.

UPDATE agent_wersja
   SET poziomy_pamieci = COALESCE((SELECT group_concat(p.poziom, ',')
                                     FROM agent_pamiec_poziom p
                                    WHERE p.agent_id = agent_wersja.agent_id), '')
 WHERE numer = (SELECT a.wersja FROM agent a WHERE a.id = agent_wersja.agent_id);

CREATE TRIGGER agent_wersja_pamiec_po_dodaniu
AFTER INSERT ON agent_pamiec_poziom
BEGIN
    UPDATE agent_wersja
       SET poziomy_pamieci = COALESCE((SELECT group_concat(p.poziom, ',')
                                         FROM agent_pamiec_poziom p
                                        WHERE p.agent_id = new.agent_id), '')
     WHERE agent_id = new.agent_id
       AND numer = (SELECT a.wersja FROM agent a WHERE a.id = new.agent_id);
END;

CREATE TRIGGER agent_wersja_pamiec_po_zdjeciu
AFTER DELETE ON agent_pamiec_poziom
BEGIN
    UPDATE agent_wersja
       SET poziomy_pamieci = COALESCE((SELECT group_concat(p.poziom, ',')
                                         FROM agent_pamiec_poziom p
                                        WHERE p.agent_id = old.agent_id), '')
     WHERE agent_id = old.agent_id
       AND numer = (SELECT a.wersja FROM agent a WHERE a.id = old.agent_id);
END;
