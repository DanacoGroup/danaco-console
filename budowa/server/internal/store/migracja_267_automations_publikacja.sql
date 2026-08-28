-- Migracja 267 wprowadza publikację, udostępnienie i budżety czasu automatyki jako cztery
-- kolumny na tabeli automatyki, ponieważ każda jest polem pojedynczym bez własnego cyklu życia.
ALTER TABLE automatyka ADD COLUMN wersja_opublikowana INTEGER;
ALTER TABLE automatyka ADD COLUMN udostepniona INTEGER NOT NULL DEFAULT 0
    CHECK(udostepniona IN (0,1));
ALTER TABLE automatyka ADD COLUMN budzet_przebiegu_sekundy INTEGER NOT NULL DEFAULT 0;
ALTER TABLE automatyka ADD COLUMN budzet_kroku_sekundy INTEGER NOT NULL DEFAULT 0;
ALTER TABLE automatyka ADD COLUMN regula_budzetu TEXT;
