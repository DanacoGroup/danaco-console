-- Migracja 055 — trwałość fasety „mowa i jakość” modułu Translate: niezgodności
-- kontroli jakości panelu, ślad syntezy mowy i ślad eksportu panelu. Czas jest
-- liczbą (ms epoki), spójnie z resztą modułu Translate.
--
-- Ton panelu nie ma własnej tabeli. `translate.panel.tone.set` zmienia jedno
-- pole jednego wiersza `panel_tlumaczenia` — kolumnę `ton`. Ton nie ma własnej
-- tożsamości ani cyklu życia odrębnego od panelu: nie da się go odczytać,
-- wylistować ani skasować niezależnie od panelu, którego dotyczy.
--
-- Niezgodności kontroli jakości są wymieniane w całości przy każdej kontroli,
-- nie dopisywane. `TranslateQualityCheckResponse.Issues` to płaski wykaz
-- zastrzeżeń wobec aktualnego stanu panelu; kontrakt nie niesie pola
-- przyrostowego (brak `Since`, brak identyfikatora poprzedniej kontroli), więc
-- nie ma do czego dopisywać. `quality.check` liczy niezgodności na nowo
-- z treści panelu w chwili wywołania (liczby, daty, waluty, symbole zastępcze,
-- długość, segmenty pominięte) i zastępuje poprzedni wykaz tej samej kontroli.
-- Każda kontrola jest więc migawką bieżącą, nie strumieniem zdarzeń.
--
-- Synteza mowy zostawia ślad trwały, bo `TranslateSpeechSynthesizeResponse`
-- oddaje `Path` — ścieżkę, po którą kontrakt każe wrócić przy kolejnym
-- zapytaniu o ten sam panel. Rdzeń mowy nie syntezuje i nie ma magazynu
-- blobów, więc `nagranie_odnosnik` jest odwołaniem do pliku dostarczonego
-- z zewnątrz, a przy jego braku zostaje NULL; kolumny na sam dźwięk nie ma.
--
-- Eksport panelu zostawia ślad z tego samego powodu: `TranslatePanelExportResponse`
-- oddaje `Path`. `panel.export` zapisuje ślad żądania i odwołanie do pliku,
-- jeśli je dostał, ale pliku sam nie wytwarza. `format` jest kolumną własną
-- (wartości ExportFormat kontraktu wprost), bo bez niej ślad nie odpowiada na
-- pytanie, w jakim kształcie panel eksportowano ostatnio.

-- ── Niezgodność kontroli jakości panelu (dziecko panelu) ─────────────────────
-- Klucz obcy twardy (`ON DELETE CASCADE`): niezgodność bez panelu, którego
-- dotyczy, nie ma sensu — usunięcie panelu zabiera ze sobą jego zastrzeżenia.
-- Wartości kolumny `rodzaj` są wartościami kontraktu (TranslationIssueKind)
-- wprost, bez tłumaczenia.
CREATE TABLE panel_tlumaczenia_niezgodnosc (
    id        INTEGER PRIMARY KEY AUTOINCREMENT,
    panel_id  INTEGER NOT NULL REFERENCES panel_tlumaczenia(id) ON DELETE CASCADE,
    rodzaj    TEXT    NOT NULL
                       CHECK(rodzaj IN ('number','date','currency','placeholder','length','omission')),
    segment   TEXT,
    szczegol  TEXT,
    utworzono INTEGER NOT NULL
);
CREATE INDEX idx_panel_tlumaczenia_niezgodnosc_panel ON panel_tlumaczenia_niezgodnosc(panel_id, utworzono DESC);

-- ── Ślad syntezy mowy panelu (odsłuch) ────────────────────────────────────────
-- `nagranie_odnosnik` jest odwołaniem do pliku dostarczonego z zewnątrz
-- — rdzeń nie syntezuje mowy, więc kolumna zostaje NULL, dopóki
-- nagranie realnie nie istnieje. Klucz obcy twardy, tak jak wyżej.
CREATE TABLE panel_tlumaczenia_synteza_mowy (
    id                INTEGER PRIMARY KEY AUTOINCREMENT,
    panel_id          INTEGER NOT NULL REFERENCES panel_tlumaczenia(id) ON DELETE CASCADE,
    nagranie_odnosnik TEXT,
    utworzono         INTEGER NOT NULL
);
CREATE INDEX idx_panel_tlumaczenia_synteza_mowy_panel ON panel_tlumaczenia_synteza_mowy(panel_id, utworzono DESC);

-- ── Ślad eksportu panelu ──────────────────────────────────────────────────────
-- `plik_odnosnik` jest odwołaniem do pliku wyniku, jeśli realnie powstał
-- (rdzeń nie ma magazynu blobów, więc zwykle zostaje NULL) — nigdy treścią
-- zmyśloną. `format` niesie wartości ExportFormat kontraktu wprost.
CREATE TABLE panel_tlumaczenia_eksport (
    id             INTEGER PRIMARY KEY AUTOINCREMENT,
    panel_id       INTEGER NOT NULL REFERENCES panel_tlumaczenia(id) ON DELETE CASCADE,
    format         TEXT    NOT NULL CHECK(format IN ('pdf','docx','markdown','html','txt')),
    plik_odnosnik  TEXT,
    utworzono      INTEGER NOT NULL
);
CREATE INDEX idx_panel_tlumaczenia_eksport_panel ON panel_tlumaczenia_eksport(panel_id, utworzono DESC);
