-- Migracja 283 dodaje wyzwalacze aktualizujące poziomy pamięci w migawce
-- bieżącej wersji eksperta przy każdej zmianie składu poziomów.

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
