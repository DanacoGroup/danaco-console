-- Migracja 281 utrwala widoczność i poziomy pamięci w wersji eksperta, żeby
-- migawka AgentVersionSnapshot spełniała pola wymagane przez kontrakt.

ALTER TABLE agent_wersja ADD COLUMN widocznosc TEXT NOT NULL DEFAULT 'global';
ALTER TABLE agent_wersja ADD COLUMN poziomy_pamieci TEXT NOT NULL DEFAULT '';

-- Migawki zastane dostają stan bieżący swojego eksperta. To jedyna odpowiedź,
-- jaką da się podać zgodnie z prawdą: historii tych dwóch pól nikt dotąd nie
-- zapisywał i nie da się jej zmyślić.
UPDATE agent_wersja
   SET widocznosc = COALESCE((SELECT a.widocznosc FROM agent a WHERE a.id = agent_wersja.agent_id), 'global'),
       poziomy_pamieci = COALESCE((SELECT group_concat(p.poziom, ',')
                                     FROM agent_pamiec_poziom p
                                    WHERE p.agent_id = agent_wersja.agent_id), '');

-- Oba wyzwalacze idą przez zdjęcie i założenie od nowa, ponieważ SQLite nie
-- zmienia treści wyzwalacza w miejscu. Warunek WHEN i lista pól zostają te
-- same — dochodzą wyłącznie dwie kolumny.
DROP TRIGGER agent_wersja_po_zalozeniu;
DROP TRIGGER agent_wersja_po_zmianie;

CREATE TRIGGER agent_wersja_po_zalozeniu
AFTER INSERT ON agent
BEGIN
    INSERT INTO agent_wersja
        (agent_id, numer, nazwa, opis, instrukcje_systemowe, kanal_kod, model, transport,
         parametry_json, imie_wlasne, favikon, ustawienia_json, tryb_nakladki, aktywny,
         widocznosc, poziomy_pamieci, powod)
    VALUES
        (new.id, new.wersja, new.nazwa, new.opis, new.instrukcje_systemowe, new.kanal_kod,
         new.model, new.transport, new.parametry_json, new.imie_wlasne, new.favikon,
         new.ustawienia_json, new.tryb_nakladki, new.aktywny,
         new.widocznosc,
         COALESCE((SELECT group_concat(p.poziom, ',') FROM agent_pamiec_poziom p
                    WHERE p.agent_id = new.id), ''),
         'zalozenie eksperta');
END;

CREATE TRIGGER agent_wersja_po_zmianie
AFTER UPDATE ON agent
WHEN new.nazwa            IS NOT old.nazwa
  OR new.opis             IS NOT old.opis
  OR new.instrukcje_systemowe IS NOT old.instrukcje_systemowe
  OR new.kanal_kod        IS NOT old.kanal_kod
  OR new.model            IS NOT old.model
  OR new.transport        IS NOT old.transport
  OR new.parametry_json   IS NOT old.parametry_json
  OR new.imie_wlasne      IS NOT old.imie_wlasne
  OR new.favikon          IS NOT old.favikon
  OR new.ustawienia_json  IS NOT old.ustawienia_json
  OR new.tryb_nakladki    IS NOT old.tryb_nakladki
  OR new.aktywny          IS NOT old.aktywny
  OR new.widocznosc       IS NOT old.widocznosc
  OR new.wersja           IS NOT old.wersja
BEGIN
    INSERT INTO agent_wersja
        (agent_id, numer, nazwa, opis, instrukcje_systemowe, kanal_kod, model, transport,
         parametry_json, imie_wlasne, favikon, ustawienia_json, tryb_nakladki, aktywny,
         widocznosc, poziomy_pamieci, powod)
    VALUES
        (new.id, new.wersja, new.nazwa, new.opis, new.instrukcje_systemowe, new.kanal_kod,
         new.model, new.transport, new.parametry_json, new.imie_wlasne, new.favikon,
         new.ustawienia_json, new.tryb_nakladki, new.aktywny,
         new.widocznosc,
         COALESCE((SELECT group_concat(p.poziom, ',') FROM agent_pamiec_poziom p
                    WHERE p.agent_id = new.id), ''),
         '')
    ON CONFLICT (agent_id, numer) DO UPDATE SET
        nazwa                = excluded.nazwa,
        opis                 = excluded.opis,
        instrukcje_systemowe = excluded.instrukcje_systemowe,
        kanal_kod            = excluded.kanal_kod,
        model                = excluded.model,
        transport            = excluded.transport,
        parametry_json       = excluded.parametry_json,
        imie_wlasne          = excluded.imie_wlasne,
        favikon              = excluded.favikon,
        ustawienia_json      = excluded.ustawienia_json,
        tryb_nakladki        = excluded.tryb_nakladki,
        aktywny              = excluded.aktywny,
        widocznosc           = excluded.widocznosc,
        poziomy_pamieci      = excluded.poziomy_pamieci,
        zapisano             = strftime('%Y-%m-%dT%H:%M:%fZ', 'now');
END;
