-- Migracja 319 — graf prototypu modułu Design: połączenia między ramkami
-- (`design.prototype.link.set`, `design.prototype.get`,
-- `design.prototype.link.remove`).
--
-- ── Dlaczego połączenie ma wiersz, a nie pole ramki ─────────────────────────
-- Z jednej ramki wychodzi tyle przejść, ile jest na niej elementów klikalnych,
-- a do jednej ramki wchodzi ich dowolnie wiele. Pole przy ramce wyraziłoby
-- jedno wyjście i milczałoby o reszcie.
--
-- ── Dlaczego ramki wskazuje się kodem ───────────────────────────────────────
-- Kontrakt nazywa ramki identyfikatorami zewnętrznymi (`fromFrameId`,
-- `toFrameId`) i tymi samymi wartościami wraca `design.prototype.get`. Klucz
-- obcy wymagałby przekładu w obie strony przy każdym odczycie grafu, a graf
-- czyta się w całości — przekład byłby robotą bez odbiorcy. Spójność pilnuje
-- adapter: obie ramki muszą leżeć w tej samej kompozycji, co połączenie, i to
-- jest sprawdzane przed zapisem.
--
-- ── Bilans zamiast ciszy ────────────────────────────────────────────────────
-- `design.prototype.get` oddaje ponadto `unreachableFrameIds` — ramki, do
-- których nie wchodzi żadne połączenie. Liczy je odczyt, nie zapis: ramka
-- osierocona dziś bywa jutro ramką początkową, więc utrwalanie tej cechy
-- w kolumnie znaczyłoby przechowywanie wniosku zamiast faktów.

CREATE TABLE polaczenie_prototypu_design (
    id                       INTEGER PRIMARY KEY AUTOINCREMENT,
    identyfikator_zewnetrzny TEXT    NOT NULL UNIQUE,
    kompozycja_id            INTEGER NOT NULL REFERENCES kompozycja_design(id) ON DELETE CASCADE,
    ramka_od_kod             TEXT    NOT NULL,
    ramka_do_kod             TEXT    NOT NULL,
    wyzwalacz                TEXT    NOT NULL,
    przejscie                TEXT    NOT NULL,
    czas_ms                  INTEGER,
    warstwa_kod              TEXT,
    utworzono                TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now')),
    zaktualizowano           TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now'))
);
CREATE INDEX idx_polaczenie_prototypu_design_kompozycja
    ON polaczenie_prototypu_design(kompozycja_id, id);
CREATE INDEX idx_polaczenie_prototypu_design_od
    ON polaczenie_prototypu_design(ramka_od_kod);
