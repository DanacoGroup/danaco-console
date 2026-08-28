-- Migracja 190 rozszerza skład debaty do kształtu kontraktu o wagę, rolę i oznaczenie kluczowego uczestnika oraz zakłada zespoły i bibliotekę ról.

ALTER TABLE debata_uczestnik ADD COLUMN kluczowy INTEGER NOT NULL DEFAULT 0
    CHECK(kluczowy IN (0,1));
ALTER TABLE debata_uczestnik ADD COLUMN waga REAL NOT NULL DEFAULT 1.0;
-- Rola pusta znaczy „bez wskazania"; wartości niepuste są wartościami kontraktu
-- (RoundtableParticipantRole).
ALTER TABLE debata_uczestnik ADD COLUMN rola TEXT NOT NULL DEFAULT ''
    CHECK(rola IN ('','debater','jury','observer','moderator'));
ALTER TABLE debata_uczestnik ADD COLUMN agent TEXT;
ALTER TABLE debata_uczestnik ADD COLUMN awatar TEXT;
ALTER TABLE debata_uczestnik ADD COLUMN opis_roli TEXT;
-- Liczba wariantów odpowiedzi zbieranych od jednej persony (zespół spójności
-- własnej). Zero i jeden znaczą jedną odpowiedź.
ALTER TABLE debata_uczestnik ADD COLUMN liczba_probek INTEGER NOT NULL DEFAULT 0;

-- Zakłada tabelę debata_zespol niosącą zapisany skład debaty do ponownego użycia wraz z uczestnikami zespołu.
CREATE TABLE debata_zespol (
    id                       INTEGER PRIMARY KEY AUTOINCREMENT,
    identyfikator_zewnetrzny TEXT    NOT NULL UNIQUE,
    nazwa                    TEXT    NOT NULL,
    -- Format zapisany razem ze składem; pusty znaczy „zespół bez formatu".
    format                   TEXT    NOT NULL DEFAULT '',
    utworzono                TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now'))
);
CREATE INDEX idx_debata_zespol_nazwa ON debata_zespol(nazwa);

CREATE TABLE debata_zespol_uczestnik (
    id             INTEGER PRIMARY KEY AUTOINCREMENT,
    zespol_id      INTEGER NOT NULL REFERENCES debata_zespol(id) ON DELETE CASCADE,
    kanal_modelu   TEXT    NOT NULL,
    nazwa_tozsamosci TEXT,
    prompt_systemowy TEXT,
    rola           TEXT    NOT NULL DEFAULT '',
    waga           REAL    NOT NULL DEFAULT 1.0,
    awatar         TEXT,
    opis_roli      TEXT,
    kolejnosc      INTEGER NOT NULL DEFAULT 0
);
CREATE INDEX idx_debata_zespol_uczestnik ON debata_zespol_uczestnik(zespol_id, kolejnosc, id);

-- Rola jest daną: Operator ją czyta, kopiuje i zmienia. Zestaw wnoszony migracją odpowiada siedmiu rolom fabrycznym debaty.
CREATE TABLE debata_rola (
    id                       INTEGER PRIMARY KEY AUTOINCREMENT,
    identyfikator_zewnetrzny TEXT    NOT NULL UNIQUE,
    nazwa                    TEXT    NOT NULL,
    prompt_systemowy         TEXT    NOT NULL,
    opis                     TEXT,
    fabryczna                INTEGER NOT NULL DEFAULT 0 CHECK(fabryczna IN (0,1)),
    utworzono                TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now'))
);

INSERT INTO debata_rola (identyfikator_zewnetrzny, nazwa, prompt_systemowy, opis, fabryczna) VALUES
    ('debata-rola-adwokat-diabla', 'Adwokat diabła',
     'Podważaj każdą tezę, której nikt nie podważył. Szukaj najmocniejszego zarzutu wobec stanowiska większości i podaj go wprost.',
     'Kontrgłos przeciw myśleniu grupowemu', 1),
    ('debata-rola-glos-ostroznosci', 'Głos ostrożności',
     'Wskazuj ryzyka, koszty i skutki uboczne omawianych rozwiązań. Podawaj warunki, przy których proponowane rozwiązanie zawodzi.',
     'Ryzyka i skutki uboczne', 1),
    ('debata-rola-ekspert-techniczny', 'Ekspert techniczny',
     'Oceniaj wykonalność techniczną. Podawaj wymagania, ograniczenia i nakład potrzebny do wykonania omawianego rozwiązania.',
     'Wykonalność i nakład', 1),
    ('debata-rola-sceptyk', 'Sceptyk',
     'Żądaj dowodów. Przy każdym twierdzeniu pytaj o źródło i rozróżniaj to, co zmierzono, od tego, co założono.',
     'Żądanie dowodów', 1),
    ('debata-rola-optymista', 'Optymista',
     'Szukaj wariantu, w którym rozwiązanie działa. Wskazuj korzyści i drogi ich osiągnięcia, nie pomijając warunków.',
     'Korzyści i drogi ich osiągnięcia', 1),
    ('debata-rola-arbiter', 'Arbiter',
     'Rozstrzygaj spory między stanowiskami. Nazywaj punkt sporny i wskazuj, które argumenty go rozstrzygają.',
     'Rozstrzyganie sporów', 1),
    ('debata-rola-moderator', 'Moderator',
     'Prowadź dyskusję. Pilnuj tematu, wskazuj powtórzenia i kieruj rozmowę na punkty nierozstrzygnięte.',
     'Prowadzenie dyskusji', 1);
