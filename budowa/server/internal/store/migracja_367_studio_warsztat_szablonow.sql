-- Migracja 367 — warsztat szablonów pism: postać wzorcowa i pola do wypełnienia.
--
-- ── Czego rdzeń NIE miał ─────────────────────────────────────────────────────
-- `studio.template.list` oddawał wykaz FABRYCZNY (pismo, umowa, raport, notatka,
-- oferta), a `studio.template.apply` zakładał z niego dokument. Nie było czym
-- szablonu założyć ani zmienić: tabela `szablon_studio` niosła nazwę, opis,
-- format, treść i pola, ale nie niosła POSTACI — więc szablon nie mógł przenieść
-- arkusza stylów, nastaw strony, nagłówka, stopki, logo ani tabel. A to jest
-- wprost wymagane: szablon niesie wszystkie te rzeczy naraz.
--
-- ── Dlaczego DOBUDOWA tabeli, a nie druga tabela szablonów ──────────────────
-- Szablon fabryczny i szablon Operatora są tym samym bytem widzianym z dwóch
-- stron — jeden jest do usunięcia, drugi nie, i to rozstrzyga kolumna
-- `fabryczny`, która już stoi. Druga tabela znaczyłaby dwa wykazy szablonów,
-- dwie ścieżki w `template.apply` i pytanie, który wygrywa przy tej samej
-- nazwie. To ten sam błąd, co dwa wykazy nośników druku.
--
-- ── Dlaczego pola zostają w `pola_json`, a nie idą do własnej tabeli ────────
-- `pola_json` jest już czytane przez warstwę danych (`SzablonStudia.PolaJSON`)
-- i przez `template.apply`. Wykaz pól jest krótki, czyta się go CAŁY przy
-- wypełnianiu i nikt nie pyta o jedno pole osobno. Własna tabela dawałaby drugie
-- miejsce prawdy o polach: `template.apply` czytałby jedno, warsztat drugie,
-- a rozjazd wyszedłby przy pierwszym szablonie założonym starą drogą. Kształt
-- zapisu rośnie (rodzaj pola, opis, wykaz wartości do wyboru, miejsce w treści),
-- ale miejsce zostaje jedno.
--
-- Blokady wzorcowe szablonu leżą w `blokada_szablonu_studio` (migracja 363) —
-- tam, bo są zakresami znaków, o które PYTA się po zakresie, a nie wykazem
-- czytanym w całości.
--
-- ── Format szablonu ──────────────────────────────────────────────────────────
-- Warunku CHECK na kolumnie `format` nie ruszamy: SQLite nie umie zmienić
-- warunku bez przepisania tabeli, a przepisanie tabeli szablonów zabrałoby
-- szablony Operatora założone wcześniej. Szablon wniesiony z pliku `.dotx` albo
-- `.ott` zapisuje się jako `docx` albo jako `markdown` wedle tego, czym jego
-- treść jest po odczycie — a to, że przyszedł z pliku szablonu, mówi
-- `zrodlo_pliku`.

ALTER TABLE szablon_studio ADD COLUMN kategoria TEXT;

-- Postać wzorcowa: arkusz stylów, nastawy strony, nagłówek, stopka, logo,
-- tabele. Jednym zapisem, tym samym kształtem, którym jedzie postać dokumentu
-- (`postac_dokumentu_studio.postac_json`) — szablon zakłada dokument, więc jego
-- postać musi dać się podstawić bez przekładu.
ALTER TABLE szablon_studio ADD COLUMN postac_json TEXT;

-- Miniatura podglądu do galerii szablonów. Odwołanie do magazynu zasobów, nie
-- bajty: obrazek w wierszu szablonu obciążyłby każdy odczyt wykazu.
ALTER TABLE szablon_studio ADD COLUMN miniatura_zasob_kod TEXT;

-- Skąd szablon przyszedł: puste znaczy założony w Studiu, wypełnione — wniesiony
-- z pliku Operatora. Bez tego oddanie szablonu do pliku nie wie, jakim formatem
-- go oddać, żeby Operator dostał to, co wniósł.
ALTER TABLE szablon_studio ADD COLUMN zrodlo_pliku TEXT;

-- Dokument, z którego szablon powstał. Służy jednej rzeczy: Operator, który
-- poprawia wzór pisma, ma wiedzieć, gdzie stoi pismo wzorcowe.
ALTER TABLE szablon_studio ADD COLUMN dokument_zrodlowy_kod TEXT;

ALTER TABLE szablon_studio ADD COLUMN zaktualizowano TEXT;

CREATE INDEX idx_szablon_studio_kategoria ON szablon_studio(kategoria, nazwa);

-- ── Wersje dokumentu: szereg Operatora i szereg autozapisu ──────────────────
-- Autozapis nie ma zaśmiecać historii wersji. Wersje nazwane i kluczowe zakłada
-- Operator; zapisy samoczynne mają być OSOBNYM szeregiem, odróżnialnym
-- w wykazie, z własną zasadą wygasania. Rozróżnienie idzie kolumną, nie drugą
-- tabelą wersji: wersja autozapisu jest wersją — przywraca się ją tą samą
-- drogą (`studio.repository.restore`) i porównuje tym samym `diff.compare`.
--
-- Domyślnie `operator`, bo wszystkie wersje zapisane przed tą dobudową założył
-- Operator — i tak trzeba je pokazać.
ALTER TABLE wersja_dokumentu_studio ADD COLUMN szereg TEXT NOT NULL DEFAULT 'operator'
    CHECK(szereg IN ('operator','autosave'));

-- Postać dokumentu w chwili założenia wersji. Bez niej porównanie POSTACI dwóch
-- wersji nie miałoby czego porównać, a przywrócenie wersji oddawałoby treść
-- sprzed i postać bieżącą — czyli dokument, którego nigdy nie było.
ALTER TABLE wersja_dokumentu_studio ADD COLUMN postac_json TEXT;

-- Tożsamość wykonawcy, który wersję założył. Kolumna `autor` (rodzaj) już stoi
-- od migracji 128; tożsamością jest kod agenta z modułu Agents.
ALTER TABLE wersja_dokumentu_studio ADD COLUMN autor_agent_kod TEXT;
ALTER TABLE wersja_dokumentu_studio ADD COLUMN autor_agent_nazwa TEXT;

CREATE INDEX idx_wersja_dokumentu_studio_szereg
    ON wersja_dokumentu_studio(dokument_id, szereg, utworzono DESC, id DESC);
