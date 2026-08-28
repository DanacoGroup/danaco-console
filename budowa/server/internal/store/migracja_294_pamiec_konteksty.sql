-- Migracja 294 dodaje nazwane konteksty pamięci, ich kontekst czynny na
-- karcie sesji oraz zasady retencji wpisów pamięci.

CREATE TABLE kontekst_pamieci (
    id                       INTEGER PRIMARY KEY AUTOINCREMENT,
    identyfikator_zewnetrzny TEXT    NOT NULL UNIQUE,
    nazwa                    TEXT    NOT NULL,
    opis                     TEXT,
    profil_kod               TEXT    NOT NULL DEFAULT '',
    poziomy_json             TEXT    NOT NULL DEFAULT '[]',
    wpisy_json               TEXT    NOT NULL DEFAULT '[]',
    prompt_systemowy         TEXT,
    czynny                   INTEGER NOT NULL DEFAULT 1,
    utworzono                INTEGER NOT NULL,
    zaktualizowano           INTEGER NOT NULL
);
CREATE INDEX idx_kontekst_pamieci_profil ON kontekst_pamieci(profil_kod, nazwa);

-- Karta sesji ma najwyżej jeden kontekst pamięci czynny naraz, dlatego
-- kluczem tej tabeli jest sama karta sesji.
CREATE TABLE kontekst_pamieci_czynny (
    sesja_kod    TEXT    NOT NULL PRIMARY KEY,
    kontekst_kod TEXT    NOT NULL
                         REFERENCES kontekst_pamieci(identyfikator_zewnetrzny) ON DELETE CASCADE,
    uaktywniono  INTEGER NOT NULL
);
CREATE INDEX idx_kontekst_pamieci_czynny_kontekst ON kontekst_pamieci_czynny(kontekst_kod);

CREATE TABLE zasada_retencji_pamieci (
    id                       INTEGER PRIMARY KEY AUTOINCREMENT,
    identyfikator_zewnetrzny TEXT    NOT NULL UNIQUE,
    zasieg                   TEXT    NOT NULL,
    zasieg_kod               TEXT    NOT NULL DEFAULT '',
    profil_kod               TEXT    NOT NULL DEFAULT '',
    dni_wygasania            INTEGER NOT NULL DEFAULT 0,
    wrazliwe_domyslnie       INTEGER NOT NULL DEFAULT 0,
    wzorce_json              TEXT    NOT NULL DEFAULT '[]',
    czynna                   INTEGER NOT NULL DEFAULT 1,
    zaktualizowano           INTEGER NOT NULL,
    UNIQUE (zasieg, zasieg_kod, profil_kod)
);
