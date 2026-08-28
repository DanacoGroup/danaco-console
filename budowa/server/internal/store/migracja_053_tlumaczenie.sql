-- Migracja 053 — trwałość rdzenia modułu Translate: okno tłumaczenia (tekst
-- źródłowy), panel tłumaczenia (byt centralny modułu) i tłumaczenie zwrotne.
--
-- Segment nie jest bytem trwałym. `translate.source.segment` oddaje w kontrakcie
-- zwykłą listę napisów (`Segments []string`) — bez identyfikatora, bez kolejności
-- trwałej między wywołaniami, bez pola, po którym dałoby się go zaktualizować
-- pojedynczo. `source.set` niesie tylko `SegmentCount`, czyli liczbę, nie
-- odwołanie do wierszy. Tabela segmentów byłaby więc bytem, którego nikt nie
-- odczytuje po identyfikatorze i nikt nie aktualizuje punktowo. Segmentacja jest
-- funkcją tekstu źródłowego w chwili wywołania — `okno_tlumaczenia.tekst_zrodlowy`
-- już go trzyma, więc `source.segment` dzieli go w locie i zwraca wynik bez
-- zapisu. `liczba_segmentow` w oknie jest jedyną trwałą pamięcią wyniku
-- ostatniego podziału, dokładnie tak, jak wymaga `SegmentCount` kontraktu.
--
-- Tłumaczenie zwrotne mieszka jako kolumna panelu, nie osobna tabela. Kontrakt
-- (`TranslateBacktranslationRunResponse`) oddaje jeden tekst na panel — nie ma
-- historii wielu przebiegów, nie ma własnego stanu ani znacznika czasu osobnego
-- od panelu. Panel jest już bytem „jeden wiersz na (okno, język docelowy)" —
-- tłumaczenie zwrotne to kolejne pole tego samego wiersza, tak jak `ton` jest
-- polem panelu mimo że ustawia je inna komenda (`translate.panel.tone.set`).
-- Kolumna jest nadpisywana przy każdym `backtranslation.run` — rdzeń nie ma
-- silnika przekładu (patrz niżej), więc trzyma wyłącznie to, co realnie dostał.
--
-- Rdzeń nie tłumaczy i nie rozpoznaje języka. Nie ma tu silnika przekładu ani
-- kanału modelu. `translation.set` zapisuje wyłącznie korektę Operatora — treść,
-- którą realnie przysłano. `source.detect` nie ma na czym oprzeć rozpoznania:
-- kolumna `okno_tlumaczenia.jezyk_zrodlowy` zostaje NULL, dopóki Operator albo
-- `source.set` jej nie poda wprost; wynik `source.detect` nie jest tu utrwalany
-- osobno, bo nie ma źródła, z którego miałby wynikać uczciwie.
-- `backtranslation.run` bez tekstu panelu nie ma co tłumaczyć z powrotem —
-- `tresc_zwrotna` zostaje wtedy pusta, a nie zmyślona.
--
-- Czas jest liczbą (ms epoki), tak jak w reszcie modułu Translate.

-- ── Okno tłumaczenia (tekst źródłowy) ─────────────────────────────────────────
-- Treść źródłowa bywa obszerna → plik na dysku, baza trzyma odwołanie.
CREATE TABLE okno_tlumaczenia (
    id                       INTEGER PRIMARY KEY AUTOINCREMENT,
    identyfikator_zewnetrzny TEXT    NOT NULL UNIQUE,
    tekst_zrodlowy           TEXT,
    tekst_zrodlowy_odwolanie TEXT,
    jezyk_zrodlowy           TEXT,
    liczba_segmentow         INTEGER,
    utworzono                INTEGER NOT NULL,
    zaktualizowano           INTEGER NOT NULL
);

-- ── Panel tłumaczenia — byt centralny modułu Translate ────────────────────────
-- Wartości kolumny `stan` są wartościami kontraktu (TranslationStatus) wprost,
-- bez tłumaczenia. `ton` ustawia komenda `translate.panel.tone.set` — kolumna
-- mieszka tu, bo to pole tego panelu, nie osobny byt. `tresc_zwrotna` też jest
-- polem panelu; patrz uzasadnienie na górze pliku.
CREATE TABLE panel_tlumaczenia (
    id                       INTEGER PRIMARY KEY AUTOINCREMENT,
    identyfikator_zewnetrzny TEXT    NOT NULL UNIQUE,
    okno_id                  INTEGER NOT NULL REFERENCES okno_tlumaczenia(id) ON DELETE CASCADE,
    jezyk                    TEXT    NOT NULL,
    tresc                    TEXT,
    tresc_odwolanie          TEXT,
    stan                     TEXT    NOT NULL DEFAULT 'pending'
                                     CHECK(stan IN ('pending','translating','ready','error')),
    ton                      TEXT,
    tresc_zwrotna            TEXT,
    zaktualizowano           INTEGER NOT NULL
);
CREATE INDEX idx_panel_tlumaczenia_okno ON panel_tlumaczenia(okno_id, jezyk);
CREATE INDEX idx_panel_tlumaczenia_stan ON panel_tlumaczenia(stan, zaktualizowano DESC);
