-- Migracja 261 — etykiety automatyki (`automation.workflow.tag.set`).
--
-- Etykieta jest wierszem, nie polem z przecinkami: wyszukiwanie po etykiecie
-- ma iść indeksem, a nie dopasowaniem podciągu, które trafiałoby „raport”
-- wewnątrz „raportowanie”.
CREATE TABLE etykieta_automatyki (
    automatyka_id INTEGER NOT NULL REFERENCES automatyka(id) ON DELETE CASCADE,
    etykieta      TEXT    NOT NULL,
    PRIMARY KEY (automatyka_id, etykieta)
);
-- Wyszukiwanie automatyk po etykiecie idzie od etykiety do automatyki.
CREATE INDEX idx_etykieta_automatyki_nazwa ON etykieta_automatyki(etykieta);
