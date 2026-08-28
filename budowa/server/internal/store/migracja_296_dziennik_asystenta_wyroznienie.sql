-- Migracja 296 — wyróżnienie wpisu dziennika asystenta (`assistant.activity.flag`).
--
-- Bez tych dwóch kolumn Activity Feed pokazywałby wyróżnienie, którego rdzeń nie
-- pamięta: `AssistantActivityEntry` niesie pole `important`, a tabela
-- `wpis_dziennika_asystenta` (migracja 050) nie ma go gdzie odłożyć. Okno
-- oznaczyłoby wpis, a po odświeżeniu wykazu oznaczenie znikałoby bez śladu.
--
-- Kolumny są dwie, bo kontrakt wnosi dwie rzeczy: sam znacznik oraz powód
-- wyróżnienia zapisywany razem z nim. Powód bez znacznika byłby notatką do
-- wpisu, którego nikt nie wyróżnił; znacznik bez powodu jest w porządku, więc
-- notatka zostaje pusta.
--
-- Wartość domyślna 0 znaczy „niewyróżniony": wpisy zastane, założone przed tą
-- migracją, nie stają się nagle ważne.

ALTER TABLE wpis_dziennika_asystenta ADD COLUMN wazny INTEGER NOT NULL DEFAULT 0;
ALTER TABLE wpis_dziennika_asystenta ADD COLUMN notatka_wyroznienia TEXT;

CREATE INDEX idx_wpis_dziennika_asystenta_wazny
    ON wpis_dziennika_asystenta(okno_kod, wazny, utworzono DESC);
