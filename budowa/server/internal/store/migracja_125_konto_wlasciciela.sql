-- Konto właściciela i drogi potwierdzenia tożsamości.
--
-- Konto jest JEDNO i schemat tego pilnuje warunkiem `id = 1`, a nie umową
-- w kodzie: platforma prowadzi jednego właściciela na dowolnej liczbie urządzeń,
-- więc drugi wiersz nie jest stanem, który wolno osiągnąć.
--
-- Login i adres e-mail są tożsamością; hasła tu nie ma. Skrót hasła leży poza
-- bazą, w sejfie poświadczeń, a wiersz metody uwierzytelnienia niesie wyłącznie
-- odwołanie do niego — tak samo jak przed dołożeniem konta. Kolumny na sekret
-- w tym schemacie nie ma i nie może być.
--
-- `potwierdzone` rozdziela dwa stany konta wymagane przy rejestracji: konto
-- powstaje niepotwierdzone i pozostaje takie do chwili potwierdzenia adresu.
-- Dopiero potwierdzenie wydaje urządzeniu token dostępu.
CREATE TABLE konto_wlasciciela (
	id           INTEGER PRIMARY KEY CHECK (id = 1),
	login        TEXT    NOT NULL,
	email        TEXT    NOT NULL,
	potwierdzone INTEGER NOT NULL DEFAULT 0,
	utworzono    INTEGER NOT NULL
);

-- Droga potwierdzenia tożsamości — jedna tabela na dwa cele, bo mechanizm jest
-- ten sam: jednorazowy materiał wysłany listem, ważny przez czas ograniczony.
--
-- W bazie leży SKRÓT drogi, nigdy sama droga. Wyciek kopii bazy nie daje wtedy
-- możliwości potwierdzenia cudzej tożsamości, bo ze skrótu nie da się odtworzyć
-- materiału, który poszedł listem.
--
-- `uzyte` zamyka drogę po pierwszym użyciu. Bez tej kolumny ta sama droga
-- otwierałaby konto wielokrotnie, a list zostaje w skrzynce na zawsze.
--
-- Cel przechowywany jest tekstem, bo są dwa i nie rosną: weryfikacja adresu przy
-- rejestracji oraz odzyskanie konta.
CREATE TABLE potwierdzenie_tozsamosci (
	skrot     TEXT    PRIMARY KEY,
	cel       TEXT    NOT NULL CHECK (cel IN ('weryfikacja', 'odzyskanie')),
	wygasa    INTEGER NOT NULL,
	uzyte     INTEGER NOT NULL DEFAULT 0,
	utworzono INTEGER NOT NULL
);

-- Wyszukiwanie idzie po skrócie (klucz główny), a sprzątanie po czasie
-- wygaśnięcia — stąd indeks na kolumnie, po której chodzi sprzątanie.
CREATE INDEX idx_potwierdzenie_wygasa ON potwierdzenie_tozsamosci (wygasa);
