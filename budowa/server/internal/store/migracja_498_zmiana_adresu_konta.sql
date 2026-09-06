-- Zmiana adresu e-mail uwierzytelniającego: adres oczekujący i dwie drogi.
--
-- Czynność ma dwa sekrety naraz i jeden adres, którego konto jeszcze nie
-- używa, więc nie mieści się w `potwierdzenie_tozsamosci`: tamta tabela zna
-- jedną drogę na wiersz i żadnego adresu. Kod potwierdzenia idzie na adres
-- nowy, droga wycofania na dotychczasowy, a konto stoi przy dotychczasowym,
-- dopóki kod nie wróci.
--
-- Terminy są dwa i różne z zamysłu. Kod żyje tyle, co każda droga
-- potwierdzenia — godzinę. Wycofanie musi przeżyć weekend, bo list na adres
-- dotychczasowy jest jedyną obroną Operatora, któremu przejęto konto, a ten
-- czyta pocztę wtedy, kiedy czyta.
CREATE TABLE zmiana_adresu_konta (
	konto_id         INTEGER NOT NULL,
	adres_nowy       TEXT    NOT NULL,
	adres_poprzedni  TEXT    NOT NULL,
	skrot_kodu       TEXT    NOT NULL,
	skrot_wycofania  TEXT    NOT NULL,
	wygasa_kod       INTEGER NOT NULL,
	wygasa_wycofanie INTEGER NOT NULL,
	proby            INTEGER NOT NULL DEFAULT 0,
	zamkniete        INTEGER NOT NULL DEFAULT 0,
	utworzono        INTEGER NOT NULL,
	zrodlo_ip        TEXT,
	urzadzenie       TEXT,
	PRIMARY KEY (konto_id, utworzono),
	FOREIGN KEY (konto_id) REFERENCES konto_wlasciciela (id) ON DELETE CASCADE
);

-- Wycofanie przychodzi bez sesji, samą drogą z listu: wiersz odnajduje się
-- wyłącznie po jej skrócie.
CREATE UNIQUE INDEX idx_zmiana_adresu_wycofanie ON zmiana_adresu_konta (skrot_wycofania);

-- Potwierdzenie przychodzi z sesji Operatora, ale kod trzeba znaleźć wśród
-- zmian czynnych tego konta — stąd para (konto, zamknięte).
CREATE INDEX idx_zmiana_adresu_czynne ON zmiana_adresu_konta (konto_id, zamkniete);
