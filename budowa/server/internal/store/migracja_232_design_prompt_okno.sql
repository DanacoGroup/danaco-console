-- Migracja 232 — prompt wydany dostaje okno i kanał.
--
-- `design.prompt.history.list` zwraca prompty wydane W TYM OKNIE, a
-- `prompt_design` (migracja 048) okna nie zna: prompt powstawał jako parametr
-- wywołania, nie jako zapis zdarzenia w oknie. Bez tej kolumny historia albo
-- mieszałaby prompty wszystkich okien, albo musiałaby wnioskować okno z zasobów,
-- które z promptu powstały — a prompt bez ani jednego zasobu (kanał odmówił)
-- wypadłby wtedy z historii, choć jest dokładnie tym, czego Operator w historii
-- szuka.
--
-- Kanał zapisujemy obok, bo kontrakt (`DesignPromptRecord.ChannelId`) o niego
-- pyta, a rdzeń wie go wyłącznie w chwili wydania.
--
-- Wartość domyślna pusta jest zamierzona: wiersze zastane powstały przed tą
-- kolumną i ich okna nie ma skąd wziąć. Pusty tekst nie jest oknem żadnego
-- modułu, więc takie prompty nie wejdą do żadnej historii — i to jest prawda
-- o nich, nie ukrycie.

ALTER TABLE prompt_design ADD COLUMN okno TEXT NOT NULL DEFAULT '';
ALTER TABLE prompt_design ADD COLUMN kanal TEXT;

-- design.prompt.history.list czyta prompty okna, od najświeższego.
CREATE INDEX idx_prompt_design_okno ON prompt_design(okno, utworzono DESC, id DESC);
