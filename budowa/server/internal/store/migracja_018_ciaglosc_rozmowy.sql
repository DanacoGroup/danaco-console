-- Migracja 018 — ciągłość rozmowy okna.
--
-- Program `claude` podaje własny identyfikator rozmowy w strumieniu zdarzeń;
-- bez utrwalenia istniałby przez jedną turę i przepadał, a przełącznik
-- `--resume` nie miałby wartości. Identyfikator rozmowy należy do okna
-- komunikacji, nie do sesji ani do wiadomości: dwa okna jednej sesji prowadzą
-- dwie niezależne rozmowy z modelem i muszą mieć osobne wznowienia.
--
-- Kolumna jest pusta do pierwszej tury. Okno nowo założone nie ma jeszcze
-- rozmowy po stronie programu `claude`; pusta wartość znaczy „zacznij nową
-- rozmowę" i jest stanem poprawnym, nie brakiem danych.

ALTER TABLE okno_komunikacji ADD COLUMN id_rozmowy_cli TEXT NOT NULL DEFAULT '';

-- Odczyt idzie zawsze po oknie, więc indeks zakłada się na kolumnie wiodącej
-- tylko tam, gdzie identyfikator jest niepusty — okna bez rozmowy nie zaśmiecają
-- indeksu.
CREATE INDEX IF NOT EXISTS idx_okno_rozmowa_cli
    ON okno_komunikacji (id_rozmowy_cli)
    WHERE id_rozmowy_cli <> '';
