-- Migracja 235 — zestawy żetonów systemu projektowego (Tokens & System Panel).
--
-- Motyw produktu jest własnością powłoki i rdzeń go NIE nadpisuje. Zestaw
-- żetonów jest bytem obok motywu: Operator zakłada go, wczytuje z zapisu
-- zewnętrznego (`design.tokenset.import`), wydaje do kodu
-- (`design.tokenset.export`) i porównuje z motywem obowiązującym. Bez własnej
-- tabeli panel czytałby żetony z motywu i nie miałby ich gdzie odłożyć — czyli
-- każda praca nad systemem projektowym kończyłaby się z zamknięciem karty.
--
-- Żeton jest wierszem, nie polem zapisu strukturalnego: wydanie do CSS, SCSS,
-- Swift i Kotlin idzie żeton po żetonie, a odsyłacz (`odsylacz_do`) rozstrzyga
-- się po nazwie roli w obrębie zestawu. Para (zestaw, nazwa roli) jest kluczem
-- bez surogatu — dwie wartości tej samej roli w jednym zestawie byłyby
-- systemem, który sam sobie przeczy.
--
-- Motyw jest kolumną zestawu, nie żetonu: oba motywy są równoprawne i niosą
-- WŁASNE wartości tych samych ról (kontrakt, `DesignTokensetSaveRequest.Theme`),
-- więc motyw jasny i ciemny to dwa zestawy, a nie jeden z podwójnymi wierszami.

CREATE TABLE zestaw_zetonow_design (
    id                       INTEGER PRIMARY KEY AUTOINCREMENT,
    identyfikator_zewnetrzny TEXT    NOT NULL UNIQUE,
    okno                     TEXT    NOT NULL,
    nazwa                    TEXT    NOT NULL,
    motyw                    TEXT,
    zaktualizowano           TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now'))
);
-- design.tokenset.list czyta zestawy okna, od ostatnio zmienianego.
CREATE INDEX idx_zestaw_zetonow_design_okno
    ON zestaw_zetonow_design(okno, zaktualizowano DESC, id DESC);

-- Rodzaj żetonu jest wartością kontraktu (DesignTokenKind), nie jej
-- tłumaczeniem. Warunku CHECK tu nie ma: kontrakt bywa rozszerzany o kolejny
-- rodzaj, a wartość spoza wykazu sprawdza adapter i odmawia nią wołającemu —
-- odbicie od schematu wracałoby jako awaria rdzenia oznaczona jako ponawialna
-- (ten sam rozstrzyg co przy `sprawdzRodzajZasobu`).
CREATE TABLE zeton_design (
    zestaw_id   INTEGER NOT NULL REFERENCES zestaw_zetonow_design(id) ON DELETE CASCADE,
    nazwa       TEXT    NOT NULL,
    rodzaj      TEXT    NOT NULL,
    wartosc     TEXT    NOT NULL,
    opis        TEXT,
    odsylacz_do TEXT,
    kolejnosc   INTEGER NOT NULL DEFAULT 0,
    PRIMARY KEY (zestaw_id, nazwa)
);
CREATE INDEX idx_zeton_design_zestaw ON zeton_design(zestaw_id, kolejnosc);
