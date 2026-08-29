// Plik definiuje nastawy poziomu aplikacja, opisujące sam program, a nie
// prowadzoną w nim treść; klucz i wartość mieszkają w tym samym rejestrze
// definicji i tabeli ustawienie, obsługiwane rodziną komend config.
package konfig

// KluczWymogLogowania jest kluczem nastawy „Wymóg logowania”. Nazwa klucza
// należy do tabeli `ustawienie`, nie do kontraktu — kontrakt zna ją jako
// dowolny `key` komendy `config.set`.
const KluczWymogLogowania = "gateway.requireLogin"

// Wartości logiczne kluczy poziomu aplikacji. Rodzaj `logiczna` odpowiada
// kolumnie ustawienie.rodzaj_wartosci, więc zapis i odczyt mówią jednym słowem.
const (
	WymogLogowaniaWlaczony = "true"
	wymogLogowaniaZAdresu  = ""
)

// Klucze konta nadawczego platformy, z którego rdzeń wysyła listy systemowe
// potwierdzenia adresu i odzyskania konta; nie jest to skrzynka operatora,
// prowadzona osobno przez moduł poczty.
const (
	KluczNadawcaHost     = "mailer.host"
	KluczNadawcaPort     = "mailer.port"
	KluczNadawcaAdres    = "mailer.address"
	KluczNadawcaNazwa    = "mailer.displayName"
	KluczNadawcaUzytkow  = "mailer.username"
	KluczNadawcaSekret   = "mailer.secret"
	KluczNadawcaStartTLS = "mailer.startTLS"
)

// Klucze sekcji „Powiadomienia" okna Ustawień — przełącznik główny, czynność
// klasy zdarzenia i kanały dodatkowe klasy. Katalog wnosi je migracją 377;
// kody klas są kodami z modelu danych (`powiadomienie.klasa`), a nie drugim
// zestawem nazw.
const KluczPowiadomieniaWlaczone = "powiadomienia.wlaczone"

// KluczKlasyPowiadomien składa klucz czynności jednej klasy zdarzenia, łącząc
// prefiks sekcji powiadomień z kodem klasy z modelu danych.
func KluczKlasyPowiadomien(klasa string) string {
	return "powiadomienia.klasa." + klasa
}

// KluczKanalowPowiadomien składa klucz kanałów dodatkowych jednej klasy
// zdarzenia, rozszerzając klucz czynności o segment kanałów.
func KluczKanalowPowiadomien(klasa string) string {
	return KluczKlasyPowiadomien(klasa) + ".kanaly"
}

// definicjeAplikacji zwraca nastawy poziomu aplikacja. Wartość domyślna
// wymogu logowania jest pusta, nie fałsz, ponieważ nastawa niesie trzy
// stany zależne od adresu nasłuchu.
func definicjeAplikacji() []Definicja {
	return []Definicja{
		{
			Klucz:    KluczWymogLogowania,
			Domyslna: wymogLogowaniaZAdresu,
			Rodzaj:   RodzajLogiczna,
			Objasnienie: "Czy gniazdo musi przedstawić token sesji bramki, zanim wykona " +
				"cokolwiek poza connection.hello, auth.login i auth.register. " +
				"Bez wskazania rozstrzyga adres nasłuchu: pętla zwrotna bez wymogu, " +
				"nasłuch szerszy z wymogiem. To nie jest bramka uprawnień — po przejściu " +
				"logowania nie ma już ani jednego pytania.",
			DozwoloneZasiegi: []Poziom{PoziomAplikacji},
		},
		{
			Klucz:  KluczNadawcaHost,
			Rodzaj: RodzajTekst,
			Objasnienie: "Adres serwera poczty wychodzącej, którym platforma nadaje listy " +
				"systemowe. Bez niego rejestracja i odzyskanie konta odmawiają, bo droga " +
				"potwierdzenia nie ma czym dojść do Operatora.",
			DozwoloneZasiegi: []Poziom{PoziomAplikacji},
		},
		{
			Klucz:  KluczNadawcaPort,
			Rodzaj: RodzajLiczba,
			Objasnienie: "Port serwera poczty wychodzącej. Pusto znaczy 587 — port zgłoszenia " +
				"z szyfrowaniem STARTTLS.",
			DozwoloneZasiegi: []Poziom{PoziomAplikacji},
		},
		{
			Klucz:            KluczNadawcaAdres,
			Rodzaj:           RodzajTekst,
			Objasnienie:      "Adres w kopercie zwrotnej listu. Serwer poczty odrzuca list bez niego.",
			DozwoloneZasiegi: []Poziom{PoziomAplikacji},
		},
		{
			Klucz:            KluczNadawcaNazwa,
			Rodzaj:           RodzajTekst,
			Objasnienie:      "Nazwa widoczna przy adresie nadawcy w skrzynce odbiorcy.",
			DozwoloneZasiegi: []Poziom{PoziomAplikacji},
		},
		{
			Klucz:  KluczNadawcaUzytkow,
			Rodzaj: RodzajTekst,
			Objasnienie: "Nazwa logowania do serwera poczty wychodzącej. Pusto znaczy serwer " +
				"bez logowania — tak stoi przekaźnik na tej samej maszynie.",
			DozwoloneZasiegi: []Poziom{PoziomAplikacji},
		},
		{
			Klucz:            KluczNadawcaSekret,
			Rodzaj:           RodzajTekst,
			Objasnienie:      "Hasło konta nadawczego platformy. Wartość jest tajemnicą — okno Konfiguracji pokazuje ją zakrytą.",
			DozwoloneZasiegi: []Poziom{PoziomAplikacji},
		},
		{
			Klucz:    KluczNadawcaStartTLS,
			Domyslna: "true",
			Rodzaj:   RodzajLogiczna,
			Objasnienie: "Czy rozmowa z serwerem poczty ma zostać zaszyfrowana. Zdjęcie tego ma " +
				"sens wyłącznie dla przekaźnika na tej samej maszynie i jest zapisem jawnym, " +
				"nie przeoczeniem.",
			DozwoloneZasiegi: []Poziom{PoziomAplikacji},
		},
	}
}
