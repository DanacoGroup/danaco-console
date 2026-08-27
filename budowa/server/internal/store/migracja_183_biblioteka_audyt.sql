-- Migracja 183 zakłada tabelę dziennika audytu repozytorium modułu Library, przyrostowego i niezmiennego, z luźnym wskazaniem zasobu kodem zewnętrznym.

CREATE TABLE wpis_audytu_biblioteki (
    id                       INTEGER PRIMARY KEY AUTOINCREMENT,
    identyfikator_zewnetrzny TEXT NOT NULL UNIQUE,
    plik_kod                 TEXT,
    czynnosc                 TEXT NOT NULL
                             CHECK(czynnosc IN ('dostep','zmiana','archiwizacja',
                                                'przywrocenie','eksport','usuniecie',
                                                'utrwalenie')),
    -- Sprawca: Operator, moduł albo model; wolny tekst, bo bywa bytem spoza katalogu kont.
    sprawca                  TEXT NOT NULL,
    opis                     TEXT,
    chwila                   TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now'))
);
-- Dziennik czyta się od najnowszego wpisu, zwykle zawężony do jednego zasobu albo do jednej czynności repozytorium.
CREATE INDEX idx_wpis_audytu_biblioteki_chwila ON wpis_audytu_biblioteki(chwila DESC, id DESC);
CREATE INDEX idx_wpis_audytu_biblioteki_plik ON wpis_audytu_biblioteki(plik_kod, chwila DESC);
CREATE INDEX idx_wpis_audytu_biblioteki_czynnosc ON wpis_audytu_biblioteki(czynnosc, chwila DESC);
