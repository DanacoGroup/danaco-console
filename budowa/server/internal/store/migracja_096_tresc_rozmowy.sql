-- Migracja 096 wprowadza tabelę blok_wiadomosci przechowującą nietekstowe bloki wiadomości oraz indeks pełnotekstowy FTS5 wiadomosc_szukanie nad treścią wiadomości.

CREATE TABLE blok_wiadomosci (
    id            INTEGER PRIMARY KEY,
    chwila        INTEGER NOT NULL,                -- epoka w milisekundach
    okno_kod      TEXT    NOT NULL,
    wiadomosc_kod TEXT    NOT NULL,
    kolejnosc     INTEGER NOT NULL,                -- porządek w obrębie wiadomości
    rodzaj        TEXT    NOT NULL
                  CHECK (rodzaj IN ('thinking','tool_use','tool_result',
                                    'image','audio','error',
                                    'provenance','account')),
    tresc         TEXT    NOT NULL DEFAULT '',
    ladunek       TEXT
);

-- Indeksy wspierają odczyt całym oknem rozmowy przy odtwarzaniu historii oraz odczyt ograniczony do pojedynczej wiadomości, zachowując porządek bloków.
CREATE INDEX idx_blok_wiadomosci_okno      ON blok_wiadomosci (okno_kod, wiadomosc_kod, kolejnosc);
CREATE INDEX idx_blok_wiadomosci_wiadomosc ON blok_wiadomosci (wiadomosc_kod, kolejnosc);

-- Indeks pełnotekstowy FTS5 wiadomosc_szukanie obejmuje kolumnę tresc tabeli wiadomosc, z tabelą zewnętrzną utrzymywaną przez wyzwalacze zgodności poniżej.
CREATE VIRTUAL TABLE wiadomosc_szukanie USING fts5(
    tresc,
    content='wiadomosc',
    content_rowid='id',
    tokenize="unicode61 remove_diacritics 2"
);

-- Zasiew wypełnia indeks pełnotekstowy treścią wiadomości zapisanych przed tą migracją, aby pozostały odnajdywalne przez wyszukiwanie.
INSERT INTO wiadomosc_szukanie(rowid, tresc)
    SELECT id, COALESCE(tresc, '') FROM wiadomosc;

CREATE TRIGGER wiadomosc_szukanie_po_wstawieniu AFTER INSERT ON wiadomosc BEGIN
    INSERT INTO wiadomosc_szukanie(rowid, tresc)
        VALUES (new.id, COALESCE(new.tresc, ''));
END;

CREATE TRIGGER wiadomosc_szukanie_po_usunieciu AFTER DELETE ON wiadomosc BEGIN
    INSERT INTO wiadomosc_szukanie(wiadomosc_szukanie, rowid, tresc)
        VALUES ('delete', old.id, COALESCE(old.tresc, ''));
END;

CREATE TRIGGER wiadomosc_szukanie_po_zmianie AFTER UPDATE OF tresc ON wiadomosc BEGIN
    INSERT INTO wiadomosc_szukanie(wiadomosc_szukanie, rowid, tresc)
        VALUES ('delete', old.id, COALESCE(old.tresc, ''));
    INSERT INTO wiadomosc_szukanie(rowid, tresc)
        VALUES (new.id, COALESCE(new.tresc, ''));
END;
