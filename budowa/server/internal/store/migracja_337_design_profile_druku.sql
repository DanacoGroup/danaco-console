-- Migracja 337 — profile wydania do druku (spady, znaczniki, przestrzeń barw,
-- rozdzielczość, profil ICC, nośnik).
--
-- Profil jest bytem okna, nie parametrem jednego wywołania: ta sama drukarnia
-- dostaje od Operatora te same nastawy przez cały rok, a przepisywanie ich
-- przy każdym wydaniu jest drogą do pliku wydanego bez spadu. Kontrola
-- przeddrukowa (`design.print.preflight`) i wydanie (`design.print.export`)
-- czytają stąd JEDNE nastawy, więc nie ma jak się rozjechać to, wobec czego
-- mierzymy, z tym, co wydajemy.
--
-- Nastawy podane wprost w żądaniu (`profile`) nie mają tu wiersza — są
-- jednorazowe z zamysłu i nie mają prawa zostać po sobie w oknie.
--
-- Rozdzielczość i spad dopuszczają NULL: brak wskazania znaczy „bierz
-- domyślne rdzenia", a zapisane zero znaczyłoby „drukuj bez spadu", co jest
-- innym rozstrzygnięciem i innym wynikiem w drukarni.

CREATE TABLE profil_druku_design (
    id                       INTEGER PRIMARY KEY AUTOINCREMENT,
    identyfikator_zewnetrzny TEXT    NOT NULL UNIQUE,
    okno                     TEXT    NOT NULL,
    nazwa                    TEXT,
    przestrzen_barw          TEXT    NOT NULL
                                     CHECK(przestrzen_barw IN ('cmyk','rgb','grayscale','spot')),
    norma                    TEXT    CHECK(norma IS NULL OR
                                           norma IN ('pdfX1a','pdfX3','pdfX4','pdfA2b','brak')),
    spad_mm                  REAL,
    znaczniki_ciecia         INTEGER NOT NULL DEFAULT 0 CHECK(znaczniki_ciecia IN (0,1)),
    znaczniki_pasowania      INTEGER NOT NULL DEFAULT 0 CHECK(znaczniki_pasowania IN (0,1)),
    pasek_barw               INTEGER NOT NULL DEFAULT 0 CHECK(pasek_barw IN (0,1)),
    rozdzielczosc            INTEGER,
    profil_icc               TEXT,
    nadruk_czerni            INTEGER NOT NULL DEFAULT 0 CHECK(nadruk_czerni IN (0,1)),
    nosnik                   TEXT,
    zaktualizowano           TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now'))
);

-- design.print.profile.list czyta profile okna, od ostatnio zmienianego.
CREATE INDEX idx_profil_druku_design_okno
    ON profil_druku_design(okno, zaktualizowano DESC, id DESC);
