-- Migracja 094 — ślad przełączeń kanału.
--
-- Zapas: gdy kanał odmawia, tura próbuje kanałem następnym, ale nie po cichu.
-- Wskazanie zapasu jest wartością danych wiersza rejestru kanałów: parametr
-- `kanal_zapasowy` w `kanal_modelu.parametry_json`, tak samo jak `program`
-- i `adapter` — Operator ustawia go komendą `channel.update` polem config, bez
-- zmiany kontraktu i bez nowej kolumny (parametry wiersza są miejscem na nastawy
-- kanału; druga kolumna obok parametru byłaby drugą prawdą).
--
-- Ta tabela jest śladem jawności. Każde wykonane przełączenie zostawia wiersz:
-- z którego kanału, na który, w której turze i dlaczego. Obok wiersza Operator
-- dostaje fragment metadanych konta w strumieniu tury (pole `powod` i `kolejne`)
-- oraz prowenancję drugiego wywołania — widzi, czym faktycznie pojechało,
-- z odpowiedzi rdzenia, nie z konfiguracji.
--
-- Zapis następuje po przełączeniu i niczego nie steruje. Identyfikatory
-- kontraktowe zamiast kluczy obcych: ślad ma przeżyć byt, którego dotyczył.

CREATE TABLE przelaczenie_kanalu (
    id            INTEGER PRIMARY KEY,
    chwila        INTEGER NOT NULL,          -- epoka w milisekundach
    okno_kod      TEXT    NOT NULL,
    wiadomosc_kod TEXT    NOT NULL,
    z_kanalu      TEXT    NOT NULL,          -- kod kanału, który odmówił
    na_kanal      TEXT    NOT NULL,          -- kod kanału, którym tura pojechała dalej
    powod         TEXT    NOT NULL DEFAULT ''
);

CREATE INDEX idx_przelaczenie_okno ON przelaczenie_kanalu (okno_kod, chwila);
