-- Migracja dokłada ekspertowi nośnik nastaw procesu nakładanych obok promptu: zaczepy
-- i reguły narzędzi w postaci napisu JSON.

ALTER TABLE agent ADD COLUMN ustawienia_json TEXT NOT NULL DEFAULT '';
