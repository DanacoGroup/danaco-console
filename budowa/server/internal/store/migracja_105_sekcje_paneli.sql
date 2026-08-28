-- Migracja 105 — układ sekcji panelu okna: kolejność, zwinięcie, zdjęcie
-- z widoku. Trwałość rodziny komend `panel.sections.get` / `panel.sections.set`.
--
-- Czego dotyczy. Panel okna operacyjnego pokazuje kilka sekcji jedna pod drugą.
-- Operator przestawia je kolejnością, zwija te, których w danej pracy nie
-- czyta, a niepotrzebne zdejmuje z widoku. Ta tabela jest nośnikiem tego
-- układu — bez niej ginąłby z zamknięciem okna.
--
-- Kluczem jest para (okno, panel), nie sam panel. Ten sam panel stoi w wielu
-- oknach naraz i w każdym z nich Operator układa go pod inną pracę. Układ
-- wspólny dla wszystkich okien byłby jedną prawdą narzuconą na wiele niezależnych
-- stanowisk; układ per okno jest tym, o co komenda pyta (`windowId` i `panelId`
-- w żądaniu obu komend kontraktu).
--
-- `okno_id` jest napisem, nie kluczem obcym — z tego samego powodu, co przy
-- dzienniku transkrypcji (075) i konsultacji doradcy (101). Wiersz okna
-- komunikacji powstaje LENIWIE, przy pierwszej wiadomości, a układ panelu
-- Operator przestawia natychmiast po otwarciu okna. Klucz obcy odmawiałby zapisu
-- układu dokładnie wtedy, gdy Operator go ustawia — i kasowałby ustawienie, które
-- się wydarzyło. Ten sam identyfikator zewnętrzny okna, którym posługuje się
-- kontrakt (`Window.id`).
--
-- Kolejność liczy się od 1 i jest ciągła. Wymusza to CHECK na dolnej granicy;
-- ciągłość ustala zapis całościowy w warstwie danych, bo baza nie zna liczby
-- sekcji panelu przed zapisem. Zero i liczby ujemne są w bazie niemożliwe, a nie
-- tylko niezalecane w kodzie — pierwsza sekcja ma numer pierwszy.
--
-- Dwa stany ukrycia, bo są to dwie różne rzeczy. `zwinieta` znaczy „sekcja jest
-- na widoku, ale zawinięta do nagłówka"; `zdjeta` znaczy „sekcji na widoku nie
-- ma wcale". Jedna kolumna trójstanowa zlewałaby oba w jedno i odbierała
-- Operatorowi możliwość zdjęcia sekcji rozwiniętej — po przywróceniu wróciłaby
-- zwinięta. Kontrakt też trzyma je osobno (`PanelSection.collapsed`,
-- `PanelSection.hidden`).
--
-- Czego tu nie ma. Nie ma kolumny z NAZWĄ sekcji ani z jej treścią: katalog
-- sekcji panelu należy do warstwy widoku, a ta tabela zapisuje wyłącznie UKŁAD
-- sekcji, które widok już zna. Wpisanie tu nazwy założyłoby drugi katalog
-- sekcji obok tego, którym posługuje się klient.

CREATE TABLE sekcja_panelu (
    id             INTEGER PRIMARY KEY AUTOINCREMENT,
    -- Identyfikator zewnętrzny okna — ten sam, którym posługuje się kontrakt.
    okno_id        TEXT    NOT NULL,
    -- Panel w obrębie okna; identyfikator nadaje widok.
    panel_id       TEXT    NOT NULL,
    -- Sekcja w obrębie panelu; identyfikator nadaje widok.
    sekcja_id      TEXT    NOT NULL,
    -- Miejsce na widoku, liczone od 1.
    kolejnosc      INTEGER NOT NULL CHECK(kolejnosc >= 1),
    -- Sekcja zawinięta do nagłówka, ale obecna na widoku.
    zwinieta       INTEGER NOT NULL DEFAULT 0 CHECK(zwinieta IN (0, 1)),
    -- Sekcja zdjęta z widoku w całości.
    zdjeta         INTEGER NOT NULL DEFAULT 0 CHECK(zdjeta IN (0, 1)),
    -- Chwila ostatniej zmiany układu, w postaci ISO 8601 UTC — tak samo jak
    -- `okno_komunikacji.zaktualizowano`, żeby oba wiersze mówiły jednym zegarem.
    zaktualizowano TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now')),
    -- Sekcja występuje w panelu okna dokładnie raz. Bez tego dwa zapisy tej samej
    -- sekcji dałyby dwa miejsca na widoku dla jednego bytu.
    UNIQUE(okno_id, panel_id, sekcja_id)
);

-- Odczyt układu jednego panelu w kolejności widoku — dokładnie to jedno pytanie
-- zadaje `panel.sections.get`. Kolumny w kolejności (okno, panel, kolejnosc)
-- czynią indeks pokrywającym porządek zapytania, więc sortowanie nie kosztuje.
CREATE INDEX idx_sekcja_panelu_uklad
    ON sekcja_panelu (okno_id, panel_id, kolejnosc);
