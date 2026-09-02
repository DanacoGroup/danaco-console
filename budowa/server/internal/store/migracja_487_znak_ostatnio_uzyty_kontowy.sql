-- Znak ostatnio użyty jest pracą jednego konta, a `kod` (codepoint znaku) jest
-- wspólny między kontami: unikat na samym kodzie kazałby dwóm kontom dzielić
-- jeden licznik użyć. Jednoznaczność obowiązuje w koncie, wzorem migracji 485.
--
-- UNIQUE stoi tu jako ograniczenie wbudowane kolumny, a nie osobny indeks, więc
-- nie da się go zdjąć samym DROP INDEX — tabelę trzeba przebudować. Kontrola
-- przejazdu pilnuje, że przebudowa nie zgubiła ani nie dorobiła wiersza:
-- niezgodność liczebności wywraca się o CHECK, transakcja migracji się cofa,
-- a wpis w rejestrze migracji nie powstaje.

CREATE TABLE IF NOT EXISTS kontrola_przejazdu (
    wersja  INTEGER NOT NULL,
    tabela  TEXT    NOT NULL,
    roznica INTEGER NOT NULL CHECK (roznica = 0)
);

CREATE TABLE znak_ostatnio_uzyty_studio_nowa (
    id        INTEGER PRIMARY KEY AUTOINCREMENT,
    kod       TEXT    NOT NULL,
    znak      TEXT    NOT NULL,
    ile_uzyc  INTEGER NOT NULL DEFAULT 1,
    uzyto     TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now')),
    konto_id  INTEGER
);

INSERT INTO znak_ostatnio_uzyty_studio_nowa (id, kod, znak, ile_uzyc, uzyto, konto_id)
    SELECT id, kod, znak, ile_uzyc, uzyto, konto_id FROM znak_ostatnio_uzyty_studio;

INSERT INTO kontrola_przejazdu (wersja, tabela, roznica) VALUES (
    487, 'znak_ostatnio_uzyty_studio',
    (SELECT COUNT(*) FROM znak_ostatnio_uzyty_studio_nowa)
        - (SELECT COUNT(*) FROM znak_ostatnio_uzyty_studio)
);

DROP TABLE znak_ostatnio_uzyty_studio;
ALTER TABLE znak_ostatnio_uzyty_studio_nowa RENAME TO znak_ostatnio_uzyty_studio;

CREATE UNIQUE INDEX idx_znak_ostatnio_uzyty_studio_kod
    ON znak_ostatnio_uzyty_studio (kod, COALESCE(konto_id, 0));
CREATE INDEX idx_znak_ostatnio_uzyty_studio_kolejnosc
    ON znak_ostatnio_uzyty_studio (uzyto DESC, ile_uzyc DESC);
