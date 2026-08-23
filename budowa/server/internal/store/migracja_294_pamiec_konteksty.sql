-- Migracja 294 — nazwane konteksty pamięci i zasady retencji (rodzina `memory.*`).
--
-- Kontekst jest ZESTAWEM WSKAZAŃ, nie właścicielem treści: wskazuje poziomy
-- pamięci, wpisy pamięci wprost oraz warstwę promptu systemowego. Dlatego
-- usunięcie kontekstu nie kasuje ani jednego wpisu pamięci — kasuje wskazanie.
-- Poziomy i wpisy idą zapisem strukturalnym w kolumnach `poziomy_json`
-- i `wpisy_json`: są wykazami identyfikatorów bez własnych atrybutów, a wpis
-- pamięci ma własną tabelę i własne komendy (`memory.set`, `memory.delete`),
-- z którymi kontekst nie ma prawa się rozjechać przez klucz obcy kasujący
-- kaskadowo.
--
-- Kontekst czynny karty sesji jest wierszem osobnej tabeli, a nie kolumną
-- w kontekście: ten sam kontekst bywa czynny w kilku kartach naraz, a kolumna
-- „aktywny" w kontekście pozwoliłaby na jedną kartę.
--
-- Zasada retencji obowiązuje ZAPISY KOLEJNE i nie rusza wstecz wpisów zastanych
-- — tak mówi kontrakt `memory.retention.set`. Schemat tego nie wymusza (to
-- rozstrzygnięcie adaptera), ale kolumna `zaktualizowano` daje granicę, od
-- której zasada obowiązuje, więc wygaszanie ma po czym odróżnić jedno od
-- drugiego.
--
-- Trójka (zasięg, byt zasięgu, profil) jest kluczem zasady: dwie zasady o tym
-- samym zasięgu byłyby dwiema odpowiedziami na jedno pytanie. Byt zasięgu
-- i profil zapisują się pustym napisem, nie NULL — z tego samego powodu co przy
-- skrótach tekstowych.

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

-- Karta sesji ma najwyżej jeden kontekst czynny, więc kluczem jest sama karta.
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
