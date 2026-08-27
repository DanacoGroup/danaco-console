-- Migracja 251 dodaje do karty terminala kolumny adresu, portu i wpisu
-- książki hostów celu zdalnego, aby karta przetrwała restart rdzenia.

ALTER TABLE terminal_karta ADD COLUMN cel_zdalny  TEXT NOT NULL DEFAULT '';
ALTER TABLE terminal_karta ADD COLUMN port_zdalny INTEGER;
ALTER TABLE terminal_karta ADD COLUMN host_kod    TEXT;
