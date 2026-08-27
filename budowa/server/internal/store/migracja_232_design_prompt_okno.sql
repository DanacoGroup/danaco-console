-- Migracja 232 dodaje do promptu wydanego kolumny okna i kanału, aby historia
-- promptów mogła zawężać się do jednego okna modułu Design.

ALTER TABLE prompt_design ADD COLUMN okno TEXT NOT NULL DEFAULT '';
ALTER TABLE prompt_design ADD COLUMN kanal TEXT;

-- Indeks wspiera odczyt historii promptów okna w porządku od najświeższego wydania do najstarszego zapisu.
CREATE INDEX idx_prompt_design_okno ON prompt_design(okno, utworzono DESC, id DESC);
