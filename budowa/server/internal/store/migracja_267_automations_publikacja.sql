-- Migracja 267 — publikacja, udostępnienie i budżety czasu automatyki
-- (`automation.workflow.publish`, `.share`, `automation.execution.budget.set`).
--
-- Cztery kolumny na tabeli `automatyka`, nie cztery tabele: każda jest polem
-- pojedynczym o krotności jeden do jednego z automatyką i nie ma własnego
-- cyklu życia. Osobna tabela na jedną liczbę byłaby złączeniem bez powodu.
--
-- `wersja_opublikowana` puste znaczy „automatyka nigdy nie została
-- opublikowana”, a nie „opublikowano wersję zerową”: wykonywana produkcyjnie
-- jest wtedy wersja bieżąca, bo Operator rozdziału wersji nie wprowadził.
--
-- Budżet zero znaczy „bez granicy” — tak samo mówi kontrakt.
ALTER TABLE automatyka ADD COLUMN wersja_opublikowana INTEGER;
ALTER TABLE automatyka ADD COLUMN udostepniona INTEGER NOT NULL DEFAULT 0
    CHECK(udostepniona IN (0,1));
ALTER TABLE automatyka ADD COLUMN budzet_przebiegu_sekundy INTEGER NOT NULL DEFAULT 0;
ALTER TABLE automatyka ADD COLUMN budzet_kroku_sekundy INTEGER NOT NULL DEFAULT 0;
ALTER TABLE automatyka ADD COLUMN regula_budzetu TEXT;
