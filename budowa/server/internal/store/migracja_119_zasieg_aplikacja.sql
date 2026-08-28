-- Migracja 119 — dziewiąty poziom zasięgu: `aplikacja`.
--
-- Osiem poziomów z migracji 001 opisuje treść prowadzoną w platformie: globalny,
-- środowisko, moduł, para modułów, projekt, karta sesji, rola, okno. Nie ma
-- wśród nich poziomu opisującego sam program — a nastawy takie istnieją i są
-- używane: wymóg logowania, adres i postać nasłuchu. Bez tego poziomu wymóg
-- logowania da się ustawić wyłącznie przy starcie rdzenia (przełącznik wiersza
-- poleceń `--wymog-logowania` albo zmienna `DANACO_WYMOG_LOGOWANIA`), bo nie ma
-- poziomu, na którym wolno by go zapisać, ani drogi, którą zapis dotarłby do
-- warstwy nasłuchu bez restartu.
--
-- Dlaczego nie nowa tabela nastaw. Osobne miejsce na „ustawienia serwera”
-- byłoby drugim źródłem prawdy o nastawach: `config.get`/`config.set` widziałyby
-- jedne wartości, nastawy nasłuchu drugie, a Operator nie miałby jak zobaczyć
-- jednych obok drugich. Dokładamy więc wiersz do słownika, który już istnieje,
-- i cała rodzina `config.*` obsługuje nowy poziom bez jednej nowej komendy.
--
-- Dlaczego pierwszeństwo 0, czyli poziom najszerszy. Kolumna `pierwszenstwo`
-- rośnie ku poziomom węższym (okno ma 8 i wygrywa
-- z każdym). Aplikacja jest szersza niż globalny: „globalny” znaczy „wszędzie
-- w treści platformy”, a „aplikacja” — „w samym programie, niezależnie od
-- jakiejkolwiek treści”. Zero stawia ją poniżej wszystkich, więc zapis na
-- którymkolwiek z ośmiu dotychczasowych poziomów nadal wygrywa i rozstrzyganie
-- sprzed tej migracji nie zmienia się ani o jeden klucz. Bytu ten poziom nie ma
-- (`klucz_zasiegu` pusty), tak samo jak globalny — programu nie ma czym zawęzić.
--
-- Dlaczego poprawka schematu, a nie przebudowa tabeli. `poziom_zasiegu.kod`
-- niesie z migracji 001 warunek CHECK wyliczający osiem
-- kodów, a SQLite nie umie zdjąć warunku poleceniem ALTER. Zwykła przebudowa
-- (nowa tabela, przepisanie wierszy, podmiana nazw) jest tu wykluczona: przy
-- `PRAGMA foreign_keys = ON` — a rdzeń
-- otwiera bazę właśnie tak (`store/baza.go`) — polecenie ALTER TABLE RENAME
-- przepisuje klauzule REFERENCES w schematach WSZYSTKICH tabel wskazujących
-- przenoszoną (`ustawienie`, `akcja`, `profil_izolacji_zasieg`,
-- `przestrzen_robocza`), więc po podmianie wskazywały tabelę odstawioną.
-- `PRAGMA legacy_alter_table` tego nie zmienia — flaga rządzi wyzwalaczami
-- i widokami, nie kluczami obcymi — a `PRAGMA foreign_keys` nie daje się
-- przestawić wewnątrz transakcji, w której migracja biegnie.
--
-- Zostaje droga, którą dokumentacja SQLite podaje dla zmiany samego warunku
-- CHECK: poprawka tekstu schematu przy `writable_schema`. Zmienia się WYŁĄCZNIE
-- napis opisujący tę jedną tabelę — nie rusza się ani jednego wiersza, ani
-- jednego identyfikatora, ani jednego klucza obcego. Numer wersji schematu
-- rośnie o jeden, żeby połączenia trzymające schemat w pamięci przeczytały go
-- na nowo; bez tego stary warunek obowiązywałby do końca biegu procesu.

PRAGMA writable_schema = ON;

UPDATE sqlite_master
   SET sql = replace(sql,
                     '''globalny'',''srodowisko''',
                     '''aplikacja'',''globalny'',''srodowisko''')
 WHERE type = 'table' AND name = 'poziom_zasiegu';

PRAGMA schema_version = 1000119;

PRAGMA writable_schema = OFF;

INSERT INTO poziom_zasiegu (kod, nazwa, pierwszenstwo) VALUES
    ('aplikacja', 'Aplikacja', 0);
