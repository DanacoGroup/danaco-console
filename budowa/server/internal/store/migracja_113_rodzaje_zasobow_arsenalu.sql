-- Migracja 113 — rodzaj zasobu przestaje kłamać o tym, co powstało.
--
-- Kolumna `zasob_design.rodzaj` miała warunek CHECK dopuszczający trzy wartości:
-- 'image', 'vector', 'composition'. Rdzeń wytwarza jednak także film, dźwięk,
-- dokument i archiwum (rodziny narzędzi `image.*`, `media.*`, `document.*`,
-- `archive.*`), a kontrakt niesie cztery brakujące wartości (`document`,
-- `audio`, `video`, `archive`), które bez tej migracji rozbiłyby się o warunek
-- schematu.
--
-- Dlaczego przebudowa, a nie ALTER. SQLite nie zna zmiany warunku CHECK
-- w miejscu — warunek jest częścią tekstu CREATE TABLE. Jedyną drogą jest
-- przebudowa. Kolumny, typy, wartości domyślne i pozostałe warunki są
-- przepisane co do znaku; ta migracja poszerza dokładnie jedną
-- listę wartości i nie zmienia niczego innego.
--
-- Dlaczego etykiety też są przebudowywane. `etykieta_zasobu_design.zasob_id`
-- wskazuje `zasob_design(id)` z ON DELETE CASCADE, a rdzeń trzyma
-- `PRAGMA foreign_keys` włączoną i stosuje
-- migracje w transakcji, więc pragmy tu wyłączyć się nie da. Przy włączonych
-- kluczach DROP starej tabeli zasobów wywołałby kaskadę i skasował wszystkie
-- etykiety. Dlatego kolejność jest odwrotna do naturalnej: najpierw kopiujemy
-- etykiety do tabeli tymczasowej, potem porzucamy oryginał (nie ma już wtedy
-- dzieci), potem tabelę zasobów, a na końcu zmieniamy nazwy. SQLite przy
-- zmianie nazwy tabeli przepisuje odwołania w tabelach potomnych, więc klucz
-- obcy kopii etykiet sam trafia na nową nazwę.
--
-- Wiersze zastane zostają nietknięte. Kusi, żeby przy okazji naprawić stare
-- kłamstwa (film zapisany jako 'image' na 'video'). Nie robimy tego: z kolumny
-- `format` da się to wywnioskować tylko czasem, a wiersz poprawiony
-- zgadywaniem jest w panelu nieodróżnialny od prawdziwego. Stare
-- wiersze mówią to, co mówiły; prawdę mówią wiersze zakładane od tej chwili.
-- Wartości `id` przepisujemy co do wartości, więc powiązania etykiet po
-- podmianie wskazują te same zasoby, co przed nią.

CREATE TABLE zasob_design_nowy (
    id                       INTEGER PRIMARY KEY AUTOINCREMENT,
    identyfikator_zewnetrzny TEXT    NOT NULL UNIQUE,
    okno                     TEXT    NOT NULL,
    nazwa                    TEXT,
    rodzaj                   TEXT    NOT NULL
                                     CHECK(rodzaj IN ('image','vector','composition',
                                                      'document','audio','video','archive')),
    format                   TEXT,
    uri                      TEXT,
    prompt_id                INTEGER REFERENCES prompt_design(id) ON DELETE SET NULL,
    wariant_zasobu_id        TEXT,
    ulubiony                 INTEGER NOT NULL DEFAULT 0 CHECK(ulubiony IN (0,1)),
    szerokosc                INTEGER,
    wysokosc                 INTEGER,
    utworzono                TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now'))
);

INSERT INTO zasob_design_nowy
    (id, identyfikator_zewnetrzny, okno, nazwa, rodzaj, format, uri, prompt_id,
     wariant_zasobu_id, ulubiony, szerokosc, wysokosc, utworzono)
SELECT id, identyfikator_zewnetrzny, okno, nazwa, rodzaj, format, uri, prompt_id,
       wariant_zasobu_id, ulubiony, szerokosc, wysokosc, utworzono
FROM zasob_design;

CREATE TABLE etykieta_zasobu_design_nowa (
    zasob_id   INTEGER NOT NULL REFERENCES zasob_design_nowy(id) ON DELETE CASCADE,
    etykieta   TEXT    NOT NULL,
    utworzono  TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now')),
    PRIMARY KEY (zasob_id, etykieta)
);

INSERT INTO etykieta_zasobu_design_nowa (zasob_id, etykieta, utworzono)
SELECT zasob_id, etykieta, utworzono FROM etykieta_zasobu_design;

DROP TABLE etykieta_zasobu_design;
DROP TABLE zasob_design;

ALTER TABLE zasob_design_nowy RENAME TO zasob_design;
ALTER TABLE etykieta_zasobu_design_nowa RENAME TO etykieta_zasobu_design;

-- Indeksy giną razem ze starymi tabelami — odtwarzamy je co do znaku:
-- `design.asset.list` filtruje po oknie i po etykiecie,
-- sortując od najnowszych.
CREATE INDEX idx_zasob_design_okno ON zasob_design(okno, utworzono DESC, id DESC);
CREATE INDEX idx_zasob_design_prompt ON zasob_design(prompt_id);
CREATE INDEX idx_etykieta_zasobu_design_etykieta ON etykieta_zasobu_design(etykieta, zasob_id);
