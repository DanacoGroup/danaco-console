-- Migracja 293 — słownik skrótów tekstowych (rodzina `snippet.*`).
--
-- Skrót rozwija się we WSZYSTKICH polach tekstowych platformy, a nie w jednym
-- oknie — dlatego słownik należy do rdzenia. Gdyby mieszkał w kliencie, ta sama
-- fraza rozwijałaby się inaczej na dwóch maszynach tego samego Operatora.
--
-- Para (profil, skrót) jest kluczem: jeden profil asystenta nie może mieć dwóch
-- rozwinięć tego samego skrótu, bo wtedy rozwinięcie rozstrzygałby przypadek
-- kolejności odczytu. Profil pusty (skrót wspólny) zapisuje się jako pusty
-- napis, nie NULL — w SQLite NULL nie jest równy NULL, więc warunek UNIQUE
-- przepuściłby dowolną liczbę duplikatów skrótu bez profilu.
--
-- Pola szablonu (`variables` kontraktu) idą zapisem strukturalnym w jednej
-- kolumnie: są wykazem nazw bez własnych atrybutów, a wypełnia je okno przy
-- rozwinięciu.

CREATE TABLE skrot_tekstowy (
    id                       INTEGER PRIMARY KEY AUTOINCREMENT,
    identyfikator_zewnetrzny TEXT    NOT NULL UNIQUE,
    profil_kod               TEXT    NOT NULL DEFAULT '',
    skrot                    TEXT    NOT NULL,
    tresc                    TEXT    NOT NULL,
    opis                     TEXT,
    pola_json                TEXT    NOT NULL DEFAULT '[]',
    czynny                   INTEGER NOT NULL DEFAULT 1,
    utworzono                INTEGER NOT NULL,
    zaktualizowano           INTEGER NOT NULL,
    UNIQUE (profil_kod, skrot)
);
CREATE INDEX idx_skrot_tekstowy_wykaz ON skrot_tekstowy(profil_kod, skrot);
