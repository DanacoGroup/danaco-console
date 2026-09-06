-- Znacznik listu o zakończeniu przebiegu.
--
-- Stan przebiegu jest wyprowadzany z kolejki przy każdym naniesieniu postępu,
-- a nie zapisywany zdarzeniem. Bez znacznika list o zakończeniu wychodziłby
-- przy każdym odczycie okna przebiegów — tyle razy, ile razy ktoś na nie
-- spojrzy. Kolumna zamyka powiadomienie po pierwszym nadaniu.
ALTER TABLE przebieg_automatyki ADD COLUMN powiadomiono INTEGER NOT NULL DEFAULT 0;

-- Wybór biegów czekających na list idzie po parze (stan, powiadomiono).
CREATE INDEX idx_przebieg_powiadomienie ON przebieg_automatyki (stan, powiadomiono);
