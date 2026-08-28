-- Tworzy tabelę konto_wlasciciela, przechowującą jedno konto właściciela platformy identyfikowane loginem i adresem e-mail, bez hasła w tym schemacie.
CREATE TABLE konto_wlasciciela (
	id           INTEGER PRIMARY KEY CHECK (id = 1),
	login        TEXT    NOT NULL,
	email        TEXT    NOT NULL,
	potwierdzone INTEGER NOT NULL DEFAULT 0,
	utworzono    INTEGER NOT NULL
);

-- Tworzy tabelę potwierdzenie_tozsamosci, przechowującą skrót jednorazowego materiału potwierdzającego, ważnego przez czas ograniczony, dla weryfikacji adresu i odzyskania konta.
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
