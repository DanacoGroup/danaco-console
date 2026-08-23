-- Migracja 206 — moduł Apps, Publisher Panel: pakiety rozszerzenia zbudowane
-- z produktu.
--
-- Pakiet jest archiwum na dysku i wierszem obok niego. `apps.package.build`
-- składa archiwum z artefaktu budowania i manifestu, kładzie je w magazynie
-- treści rdzenia i zapisuje tu odwołanie, rozmiar i format. Kolejne komendy
-- rodziny pracują na tym samym wierszu: `apps.package.manifest.save` wymienia
-- manifest, `apps.package.validate` czyta go do raportu zastrzeżeń,
-- `apps.package.sign` dopisuje podpis, `apps.package.publish` — kod pozycji
-- katalogu, która z pakietu powstała.
--
-- Manifest i podpis leżą jako surowy JSON kontraktu (`AppPackageManifest`,
-- `ExtensionSignature`). Rozłożenie manifestu na kolumny znaczyłoby drugą
-- definicję kształtu, którego jedynym źródłem jest kontrakt, a narzędzia
-- i uprawnienia pakietu wychodzą zawsze w komplecie razem z pakietem — nie ma
-- po czym filtrować.
--
-- `rozszerzenie_kod` wskazuje pozycję katalogu kodem, nie więzem obcym: pozycja
-- żyje w tabeli `rozszerzenie` (migracja 070) własnym cyklem życia i jej
-- odinstalowanie nie ma prawa skasować pakietu, z którego powstała.

CREATE TABLE pakiet_apps (
    id                       INTEGER PRIMARY KEY AUTOINCREMENT,
    identyfikator_zewnetrzny TEXT    NOT NULL UNIQUE,
    okno                     TEXT    NOT NULL,
    -- Surowy JSON `AppPackageManifest`; puste, dopóki manifestu nie zapisano.
    manifest                 TEXT,
    -- Artefakt, z którego pakiet zbudowano (`AppArtifact.id`).
    artefakt_odwolanie       TEXT,
    -- Wartości kontraktu (AppPackageFormat) wprost.
    format                   TEXT    NOT NULL DEFAULT 'zip' CHECK(format IN ('zip','targz')),
    -- Odwołanie względne magazynu treści rdzenia do archiwum pakietu.
    sciezka                  TEXT,
    rozmiar                  INTEGER,
    -- Surowy JSON `ExtensionSignature`; puste, dopóki pakietu nie podpisano.
    podpis                   TEXT,
    -- Kod pozycji katalogu powstałej z pakietu po publikacji.
    rozszerzenie_kod         TEXT,
    utworzono                TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now')),
    zaktualizowano           TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now'))
);
CREATE INDEX idx_pakiet_apps_okno ON pakiet_apps(okno, id DESC);
