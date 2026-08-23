-- Migracja 318 — komponenty makiety wraz z wariantami i ich instancje
-- (`design.component.save`, `design.component.list`,
-- `design.component.instance.add`).
--
-- ── Dlaczego komponent wisi na oknie, a nie na kompozycji ───────────────────
-- Kontrakt (`DesignComponentSaveRequest.WindowId`, `DesignComponent.WindowId`)
-- wiąże komponent z oknem modułu. Tak jest, bo biblioteka UI jest wspólna dla
-- wszystkich plansz stanowiska: przycisk zdefiniowany przy makiecie ekranu
-- logowania ma dać się wstawić na planszę ekranu ustawień. Komponent przypięty
-- do kompozycji byłby biblioteką na jedną planszę, czyli nie byłby biblioteką.
--
-- ── Dlaczego instancja ma wiersz ────────────────────────────────────────────
-- `instanceCount` i `propagatedTo` są liczbami mierzonymi, nie deklarowanymi.
-- Bez wiersza na instancję rdzeń mógłby oddać co najwyżej zero albo zmyślenie,
-- a zdanie „zmiana doszła do siedmiu instancji" ma znaczyć siedem wierszy.
-- Instancja wskazuje warstwę identyfikatorem zewnętrznym z tego samego powodu,
-- co przynależność do ramki (migracja 317): `design.board.update` przepisuje
-- warstwy od nowa i klucz wiersza warstwy nie przeżywa zapisu planszy.
--
-- Warianty jadą jednym zapisem JSON, bo kontrakt nadsyła ich komplet przy
-- każdym zapisie komponentu (`DesignComponentSaveRequest.Variants`) i żadne
-- zapytanie nie pyta o pojedynczy wariant w oderwaniu od komponentu.

CREATE TABLE komponent_design (
    id                       INTEGER PRIMARY KEY AUTOINCREMENT,
    identyfikator_zewnetrzny TEXT    NOT NULL UNIQUE,
    okno                     TEXT    NOT NULL,
    nazwa                    TEXT    NOT NULL,
    warianty_json            TEXT,
    zestaw_zetonow_kod       TEXT,
    utworzono                TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now')),
    zaktualizowano           TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now'))
);
CREATE INDEX idx_komponent_design_okno ON komponent_design(okno, id);

CREATE TABLE instancja_komponentu_design (
    id            INTEGER PRIMARY KEY AUTOINCREMENT,
    komponent_id  INTEGER NOT NULL REFERENCES komponent_design(id) ON DELETE CASCADE,
    kompozycja_id INTEGER NOT NULL REFERENCES kompozycja_design(id) ON DELETE CASCADE,
    warstwa_kod   TEXT    NOT NULL UNIQUE,
    wariant       TEXT,
    utworzono     TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now'))
);
CREATE INDEX idx_instancja_komponentu_design_komponent
    ON instancja_komponentu_design(komponent_id, id);
