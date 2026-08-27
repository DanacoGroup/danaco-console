-- Migracja rozszerza tabelę konto o pola wymagane przez strukturę konta kontraktu:
-- dostawcę, model domyślny, adres punktu końcowego, oznaczenie domyślności i stan rotacji.

ALTER TABLE konto ADD COLUMN dostawca       TEXT    NOT NULL DEFAULT '';
ALTER TABLE konto ADD COLUMN model_domyslny TEXT;
ALTER TABLE konto ADD COLUMN adres_bazowy   TEXT;
ALTER TABLE konto ADD COLUMN domyslne       INTEGER NOT NULL DEFAULT 0
                                            CHECK(domyslne IN (0,1));
ALTER TABLE konto ADD COLUMN stan           TEXT    NOT NULL DEFAULT 'aktywne'
                                            CHECK(stan IN ('aktywne','wyczerpane','zawieszone'));
ALTER TABLE konto ADD COLUMN wyczerpane_do  TEXT;

-- Konto domyślne jest dokładnie jedno na rodzaj — pilnuje tego indeks częściowy
-- w bazie, nie warunek w kodzie aplikacji.
CREATE UNIQUE INDEX idx_konto_domyslne_rodzaj ON konto(rodzaj) WHERE domyslne = 1;

-- Odczyt puli rotacji idzie po rodzaju konta: wskazuje konta czynne w kolejności rotacji dla tego samego rodzaju konta.
CREATE INDEX idx_konto_rotacja ON konto(rodzaj, aktywne, kolejnosc);
