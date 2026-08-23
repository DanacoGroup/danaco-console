-- Migracja 347 — nastawy warsztatu fotografii modułu Design
-- (`design.photo.preset.save`, `design.photo.preset.list`,
-- `design.photo.batch.apply`).
--
-- ── Nastawa jest ZESTAWEM CZYNNOŚCI, nie jedną czynnością ───────────────────
-- Operator zapisuje pod nazwą całą drogę: „prostuj, rozjaśnij o pół działki,
-- wyostrzyj". Wsad (`design.photo.batch.apply`) puszcza tę samą drogę na wielu
-- zasobach. Czynności leżą więc jednym zapisem JSON, w kolejności wykonania —
-- kolejność jest treścią nastawy, a nie jej ozdobą.
--
-- ── Nastawa należy do OKNA, nie do zasobu ───────────────────────────────────
-- Nastawa jest sposobem pracy Operatora, a nie cechą zdjęcia: zapisuje się raz
-- i stosuje do wszystkiego, co przez to okno przechodzi. Stąd unikat na parze
-- (okno, nazwa) — dwie nastawy o tej samej nazwie w jednym oknie byłyby dwiema
-- odpowiedziami na jedno wskazanie.

CREATE TABLE nastawa_fotografii_design (
    id                       INTEGER PRIMARY KEY AUTOINCREMENT,
    identyfikator_zewnetrzny TEXT    NOT NULL UNIQUE,
    okno                     TEXT    NOT NULL,
    nazwa                    TEXT    NOT NULL,
    czynnosci_json           TEXT    NOT NULL,
    utworzono                TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now')),
    zaktualizowano           TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now')),
    UNIQUE(okno, nazwa)
);

CREATE INDEX idx_nastawa_fotografii_design_okno
    ON nastawa_fotografii_design(okno, nazwa);
