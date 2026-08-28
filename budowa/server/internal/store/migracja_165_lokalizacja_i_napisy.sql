-- Migracja 165 — lokalizacja oprogramowania (`translate.resource.*`) i napisy
-- (`translate.subtitle.*`, `translate.dubbing.script.build`).
--
-- Zasób lokalizacyjny to plik kluczy wniesiony do okna: JSON, YAML, properties,
-- Android XML, iOS strings i stringsdict, RESX, gettext PO. Klucz jest bytem
-- adresowanym po nazwie w obrębie zasobu (`resource.key.context.set` przyjmuje
-- `key`, nie identyfikator), więc para (zasób, klucz) jest kluczem naturalnym.
--
-- `formy_mnogie` trzyma JSON, i to jest wyjątek świadomy: liczba form mnogich
-- zależy od języka (angielski ma dwie, polski trzy, arabski sześć), a nazwy
-- form są nazwami CLDR (`one`, `few`, `many`, `other`). Tabela dziecka miałaby
-- tyle wierszy, ile form, i ani jednego zapytania, które by po nich zawężało —
-- formy czyta się zawsze kompletem razem z kluczem.
--
-- `zrzut_zasob_id` wskazuje zasób modułu Design (`zasob_design`) ze zrzutem
-- ekranu, na którym klucz widać. Odwołanie jest miękkie (kod zewnętrzny, bez
-- klucza obcego), bo zrzut należy do innego modułu i jego usunięcie nie ma
-- prawa skasować kontekstu klucza.
--
-- Kwestia napisów mieszka przy panelu, nie przy oknie: napisy są tłumaczone,
-- więc każdy język ma własne taktowanie i własny podział linii. Import napisów
-- wnosi kwestie źródłowe do okna, dlatego `panel_id` bywa pusty — wtedy kwestia
-- jest kwestią materiału źródłowego, nie przekładu.

-- ── Zasób lokalizacyjny ─────────────────────────────────────────────────────
CREATE TABLE zasob_lokalizacji (
    id                       INTEGER PRIMARY KEY AUTOINCREMENT,
    identyfikator_zewnetrzny TEXT    NOT NULL UNIQUE,
    okno_id                  INTEGER NOT NULL REFERENCES okno_tlumaczenia(id) ON DELETE CASCADE,
    sciezka                  TEXT    NOT NULL,
    format                   TEXT    NOT NULL
                                     CHECK(format IN ('json','yaml','properties','androidXml',
                                                      'iosStrings','iosStringsdict','resx','gettextPo')),
    jezyk_zrodlowy           TEXT,
    utworzono                INTEGER NOT NULL,
    zaktualizowano           INTEGER NOT NULL
);
CREATE INDEX idx_zasob_lokalizacji_okno ON zasob_lokalizacji(okno_id, zaktualizowano DESC);

CREATE TABLE klucz_lokalizacji (
    id             INTEGER PRIMARY KEY AUTOINCREMENT,
    zasob_id       INTEGER NOT NULL REFERENCES zasob_lokalizacji(id) ON DELETE CASCADE,
    klucz          TEXT    NOT NULL,
    tresc          TEXT    NOT NULL,
    -- Znaczniki podstawienia wypisane z treści, rozdzielone znakiem nowej
    -- linii. Kontrakt oddaje je wykazem; kolumna trzyma je w postaci, która nie
    -- wymaga drugiej tabeli na byt bez własnej tożsamości.
    znaczniki      TEXT,
    kontekst       TEXT,
    zrzut_zasob_id TEXT,
    formy_mnogie   TEXT,
    kolejnosc      INTEGER NOT NULL DEFAULT 0,
    UNIQUE(zasob_id, klucz)
);

-- ── Kwestia napisów ─────────────────────────────────────────────────────────
CREATE TABLE kwestia_napisow (
    id         INTEGER PRIMARY KEY AUTOINCREMENT,
    okno_id    INTEGER NOT NULL REFERENCES okno_tlumaczenia(id) ON DELETE CASCADE,
    panel_id   INTEGER          REFERENCES panel_tlumaczenia(id) ON DELETE CASCADE,
    kolejnosc  INTEGER NOT NULL,
    poczatek_ms INTEGER NOT NULL,
    koniec_ms   INTEGER NOT NULL,
    tresc      TEXT    NOT NULL,
    mowca      TEXT
);
CREATE INDEX idx_kwestia_napisow_okno ON kwestia_napisow(okno_id, kolejnosc);
CREATE INDEX idx_kwestia_napisow_panel ON kwestia_napisow(panel_id, kolejnosc);
