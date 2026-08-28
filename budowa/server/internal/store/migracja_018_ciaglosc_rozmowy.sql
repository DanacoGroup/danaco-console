-- Migracja dokłada kolumnę identyfikatora rozmowy narzędzia wiersza poleceń do okna
-- komunikacji, utrwalającą ciągłość rozmowy między turami.

ALTER TABLE okno_komunikacji ADD COLUMN id_rozmowy_cli TEXT NOT NULL DEFAULT '';

-- Odczyt idzie zawsze po oknie, więc indeks zakłada się na kolumnie wiodącej
-- tylko tam, gdzie identyfikator jest niepusty — okna bez rozmowy nie zaśmiecają
-- indeksu.
CREATE INDEX IF NOT EXISTS idx_okno_rozmowa_cli
    ON okno_komunikacji (id_rozmowy_cli)
    WHERE id_rozmowy_cli <> '';
