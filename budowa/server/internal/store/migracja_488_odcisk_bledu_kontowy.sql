-- Odcisk błędu (fingerprint treści) jest wspólny między kontami: dwa konta,
-- które trafią na ten sam błąd, mają go zliczyć osobno, a nie dzielić jeden
-- wiersz z licznikiem wystąpień. Jednoznaczność odcisku obowiązuje w koncie,
-- wzorem migracji 485. `kod` zostaje jednoznaczny globalnie — nadaje go rdzeń
-- z licznika i losu, więc dwa konta nie zajmą tej samej wartości.
--
-- Odcisk stoi jako UNIQUE wbudowany w kolumnę, nie osobny indeks, więc tabelę
-- trzeba przebudować. Kontrola przejazdu pilnuje liczebności: niezgodność
-- wywraca się o CHECK i cofa transakcję migracji.

CREATE TABLE IF NOT EXISTS kontrola_przejazdu (
    wersja  INTEGER NOT NULL,
    tabela  TEXT    NOT NULL,
    roznica INTEGER NOT NULL CHECK (roznica = 0)
);

CREATE TABLE diagnostyka_blad_nowa (
    id          INTEGER PRIMARY KEY AUTOINCREMENT,
    kod         TEXT    NOT NULL UNIQUE,
    odcisk      TEXT    NOT NULL,
    tresc       TEXT    NOT NULL,
    zrodlo      TEXT,
    kod_bledu   TEXT    NOT NULL DEFAULT 'internal_error'
                        CHECK(kod_bledu IN ('validation_failed','not_found','not_authenticated',
                                            'permission_denied','conflict','channel_unavailable',
                                            'rate_limited','internal_error')),
    stan        TEXT    NOT NULL DEFAULT 'new'
                        CHECK(stan IN ('new','analyzing','resolved','ignored')),
    priorytet   TEXT    NOT NULL DEFAULT 'medium'
                        CHECK(priorytet IN ('low','medium','high','critical')),
    wystapienia INTEGER NOT NULL DEFAULT 1,
    notatka     TEXT,
    kontekst    TEXT,
    pierwsze    INTEGER NOT NULL,
    ostatnie    INTEGER NOT NULL,
    konto_id    INTEGER
);

INSERT INTO diagnostyka_blad_nowa
    (id, kod, odcisk, tresc, zrodlo, kod_bledu, stan, priorytet, wystapienia,
     notatka, kontekst, pierwsze, ostatnie, konto_id)
    SELECT id, kod, odcisk, tresc, zrodlo, kod_bledu, stan, priorytet, wystapienia,
           notatka, kontekst, pierwsze, ostatnie, konto_id
      FROM diagnostyka_blad;

INSERT INTO kontrola_przejazdu (wersja, tabela, roznica) VALUES (
    488, 'diagnostyka_blad',
    (SELECT COUNT(*) FROM diagnostyka_blad_nowa) - (SELECT COUNT(*) FROM diagnostyka_blad)
);

DROP TABLE diagnostyka_blad;
ALTER TABLE diagnostyka_blad_nowa RENAME TO diagnostyka_blad;

CREATE UNIQUE INDEX idx_diagnostyka_blad_odcisk
    ON diagnostyka_blad (odcisk, COALESCE(konto_id, 0));
CREATE INDEX idx_diagnostyka_blad_ostatnie  ON diagnostyka_blad(ostatnie DESC);
CREATE INDEX idx_diagnostyka_blad_stan       ON diagnostyka_blad(stan, ostatnie DESC);
CREATE INDEX idx_diagnostyka_blad_priorytet  ON diagnostyka_blad(priorytet, ostatnie DESC);
CREATE INDEX idx_diagnostyka_blad_zrodlo     ON diagnostyka_blad(zrodlo, ostatnie DESC);
