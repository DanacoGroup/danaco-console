-- Migracja 186 zakłada tabelę sugestii porządkujących modułu Library, niosącą uzasadnienie jako pole obowiązkowe oraz stan decyzji Operatora.

CREATE TABLE sugestia_biblioteki (
    id                       INTEGER PRIMARY KEY AUTOINCREMENT,
    identyfikator_zewnetrzny TEXT    NOT NULL UNIQUE,
    plik_kod                 TEXT    NOT NULL,
    rodzaj                   TEXT    NOT NULL
                             CHECK(rodzaj IN ('etykieta','kolekcja','duplikat',
                                              'osierocony','wrazliwy')),
    wartosc                  TEXT,
    uzasadnienie             TEXT    NOT NULL,
    -- Pewność w setnych, jeśli model ją podał; brak znaczy „model nie mówi".
    pewnosc                  INTEGER,
    stan                     TEXT    NOT NULL DEFAULT 'oczekujaca'
                             CHECK(stan IN ('oczekujaca','przyjeta','odrzucona')),
    rozstrzygnieto           TEXT,
    utworzono                TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now'))
);
CREATE INDEX idx_sugestia_biblioteki_stan ON sugestia_biblioteki(stan, utworzono DESC);
CREATE INDEX idx_sugestia_biblioteki_plik ON sugestia_biblioteki(plik_kod, stan);
