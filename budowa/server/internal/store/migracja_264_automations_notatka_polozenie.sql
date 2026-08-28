-- Migracja 264 wprowadza notatkę przy kroku i położenie węzła na kanwie: wiersz jest kluczowany
-- kodem kroku, a nie kluczem wiersza tabeli kroków, ponieważ zapis definicji podmienia komplet
-- wierszy kroków przy każdym zapisie przepływu.
CREATE TABLE adnotacja_kroku_automatyki (
    automatyka_id INTEGER NOT NULL REFERENCES automatyka(id) ON DELETE CASCADE,
    krok_kod      TEXT    NOT NULL,
    notatka       TEXT,
    wspolrzedna_x INTEGER NOT NULL DEFAULT 0,
    wspolrzedna_y INTEGER NOT NULL DEFAULT 0,
    PRIMARY KEY (automatyka_id, krok_kod)
);

-- Odwołania do skarbca idą kolumną kroku, a nie adnotacją, ponieważ odwołanie do poświadczenia
-- jest częścią definicji kroku i znika razem z krokiem, gdy krok przestaje wywoływać usługę.
ALTER TABLE krok_automatyki ADD COLUMN odwolania_sekretow TEXT;
