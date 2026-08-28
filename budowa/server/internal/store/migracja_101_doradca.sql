-- Migracja 101 — DORADCA: trwałość konsultacji silniejszym modelem.
--
-- CZYM JEST DORADCA. Doradca to kanał modelu, którego agent woła po RADĘ, a nie
-- po wykonanie pracy. Pojęcie stoi na trzech własnościach i ta migracja utrwala
-- wszystkie trzy, bo bez śladu żadna z nich nie przetrwa zamknięcia okna:
--
--   (a) RADA JEST JAWNA — kolumny `pytanie`, `doradca_kanal`, `doradca_model`
--       i `rada` stoją OBOK SIEBIE w jednym wierszu. Nie da się odczytać rady
--       bez odczytania, kto ją dał i o co był pytany. Gdyby treść rady lądowała
--       w tabeli wiadomości agenta, po tygodniu nikt by nie odróżnił rady cudzej
--       od odpowiedzi własnej — i to jest powód, dla którego ta tabela istnieje
--       osobno, a nie jako kolejny rodzaj wiersza `wiadomosc`.
--
--   (b) KONSULTACJA LĄDUJE W PROWENANCJI — kolumna `prowenancja` niesie TEN SAM
--       napis JSON, który poszedł strumieniem jako fragment `provenance`.
--       Nie jest to drugi opis wywołania złożony przy zapisie:
--       złożenie go tu po raz drugi rozjechałoby ślad w bazie ze śladem
--       w oknie. Rdzeń kopiuje fragment, a nie odtwarza go.
--
--   (c) KONSULTACJA, NIE DELEGACJA — w tej tabeli NIE MA kolumny stanu pracy,
--       kolumny wyniku ani wiązania z pozycją kolejki. Rada niczego nie zleca
--       i nie ma własnego biegu; odpowiedzialność za wynik zostaje przy agencie
--       pytającym, a jego praca ma własny ślad gdzie indziej. Kolumna „czy radę
--       wykonano" byłaby pytaniem o delegację i tu jej nie będzie.
--
-- KTO MOŻE BYĆ DORADCĄ — TEJ MIGRACJI NIE OBCHODZI. Rozstrzygają to DANE
-- rejestru kanałów: parametry `doradca` i `sila` w kolumnie
-- `kanal_modelu.parametry`. Nowy doradca to nowy wiersz istniejącej tabeli, nie
-- nowa tabela i nie nowa kolumna — dlatego ta migracja nie dokłada ani jednego
-- pola do `kanal_modelu`.
--
-- ODMOWY ZAPISUJEMY RAZEM Z RADAMI. Kolumna `stan` niesie dwie wartości, bo
-- pierwsze pytanie Operatora brzmi „dlaczego nie zapytano doradcy", a odpowiedź
-- na nie jest możliwa tylko wtedy, gdy nieodbyta konsultacja zostawia wiersz.
-- CHECK tabelowy wiąże stan z treścią w obie strony: rada MUSI nieść radę,
-- odmowa MUSI nieść powód. Wiersz „rada bez rady" ma być niemożliwy w bazie,
-- a nie tylko niezalecany w kodzie.
--
-- CZAS JEST LICZBĄ, A ZEGAR JEST JEDEN. `utworzono` to milisekundy
-- epoki podane przez wołającego — tak samo jak `transkrypcja.utworzono`
-- (`migracja_075_mowa.sql`) i `rozszerzenie.zaktualizowano`
-- (`migracja_070_katalog_rozszerzen.sql`). Baza własnego „teraz" nie wstawia;
-- drugi zegar rozjeżdżałby chwilę rady z chwilą zapisu.
--
-- `okno_id` JEST NAPISEM, NIE KLUCZEM OBCYM — z tego samego powodu, co przy
-- dzienniku transkrypcji: konsultację prowadzi też praca spoza okna rozmowy,
-- a odmowa zapisu śladu z powodu nieznanego okna kasowałaby dowód zdarzenia,
-- które się wydarzyło.

CREATE TABLE konsultacja_doradcy (
    id                       INTEGER PRIMARY KEY AUTOINCREMENT,
    -- Tożsamość wpisu widoczna poza bazą; klucz sztuczny zostaje sprawą bazy.
    identyfikator_zewnetrzny TEXT    NOT NULL UNIQUE,
    -- Okno, w którym padła konsultacja. NULL znaczy „spoza okna".
    okno_id                  TEXT,
    -- Kanał agenta pytającego — po nim rozstrzygał się próg siły doradcy.
    pytajacy_kanal           TEXT    NOT NULL,
    -- Kanał i model doradcy zapisane TAK, JAK OBOWIĄZYWAŁY W CHWILI PYTANIA.
    -- Operator jutro przestawi model kanału, a ślad ma mówić, kto radził wtedy.
    doradca_kanal            TEXT    NOT NULL,
    doradca_model            TEXT    NOT NULL,
    -- Pełna treść wysłana doradcy, wraz z ramą konsultacji. Skrócenie jej do
    -- samej sprawy zabrałoby dowód, że doradcy powiedziano wprost, iż ma radzić,
    -- a nie wykonywać (własność c).
    pytanie                  TEXT    NOT NULL,
    -- Treść rady. NULL przy odmowie — nie ma czego zapisać.
    rada                     TEXT,
    -- Skrót z pytania i rady; diagnostyka, nigdy bramka.
    skrot                    TEXT,
    -- Fragment `provenance` strumienia, przepisany co do znaku (własność b).
    prowenancja              TEXT,
    stan                     TEXT    NOT NULL CHECK(stan IN ('rada', 'odmowa')),
    -- Nazwany powód odmowy; przy radzie nie ma czego opisywać.
    powod                    TEXT,
    -- Milisekundy epoki podane przez wołającego.
    utworzono                INTEGER NOT NULL,
    CHECK((stan = 'rada'   AND rada  IS NOT NULL AND rada  <> '' AND powod IS NULL)
       OR (stan = 'odmowa' AND powod IS NOT NULL AND powod <> '' AND rada  IS NULL))
);

-- Wykaz konsultacji okna od najnowszej — dokładnie to pytanie zadaje okno
-- pokazujące, co agent tego okna konsultował. Kolumny w kolejności (okno,
-- utworzono, id) czynią indeks pokrywającym porządek zapytania, więc „najnowsze
-- najpierw" nie kosztuje sortowania wyniku.
CREATE INDEX idx_konsultacja_doradcy_wykaz
    ON konsultacja_doradcy (okno_id, utworzono, id);

-- Wykaz konsultacji jednego doradcy — odpowiedź na pytanie „ile i o co pytano
-- tego doradcę", które rozstrzyga o koszcie silniejszego modelu.
CREATE INDEX idx_konsultacja_doradcy_doradca
    ON konsultacja_doradcy (doradca_kanal, utworzono);
