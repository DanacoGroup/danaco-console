// Odpowiedzialność pliku: definicje nastaw poziomu `aplikacja` — tych, które
// opisują sam program, a nie treść w nim prowadzoną.
//
// Nastawy wykonania (definicje_wykonania.go) i izolacji (definicje_izolacji.go)
// rozstrzygają, jak platforma prowadzi rozmowę. Te rozstrzygają, jak stoi sam
// rdzeń: czy nasłuch wymaga logowania. Poziom zasięgu, na którym mieszkają, jest
// najszerszy i nie ma bytu — programu nie ma czym zawęzić.
//
// Klucz siedzi w tym samym rejestrze definicji, wartość w tej samej tabeli
// `ustawienie`, a odczyt idzie tym samym rozstrzygaczem; rodzina
// `config.get` / `config.set` / `config.reset` obsługuje go bez nowej komendy.
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

// Klucze konta nadawczego platformy — skrzynki, z której rdzeń pisze dwa listy
// systemowe: potwierdzenie adresu przy rejestracji i drogę odzyskania konta.
//
// Nastawy wchodzą także zmiennymi środowiska przy starcie
// (`konfiguracja/srodowisko.go`); zapis w tabeli `ustawienie` je przesłania,
// bo pomyłki w adresie serwera poczty nie da się naprawić bez zatrzymania
// rdzenia, jeżeli jedyną drogą jest środowisko. To NIE jest skrzynka Operatora
// — skrzynki Operatora prowadzi moduł Poczty z własnym sejfem poświadczeń.
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

// KluczKlasyPowiadomien składa klucz czynności jednej klasy zdarzenia.
func KluczKlasyPowiadomien(klasa string) string {
	return "powiadomienia.klasa." + klasa
}

// KluczKanalowPowiadomien składa klucz kanałów dodatkowych jednej klasy.
func KluczKanalowPowiadomien(klasa string) string {
	return KluczKlasyPowiadomien(klasa) + ".kanaly"
}

// definicjeAplikacji zwraca nastawy poziomu `aplikacja`.
//
// Wartość domyślna wymogu logowania jest pusta, a nie „false”, bo nastawa ma
// trzy stany (zob. transport/ustawienia.go): brak wskazania — rozstrzyga adres
// nasłuchu, wskazanie „tak” i wskazanie „nie”. Domyślna `false` skasowałaby stan
// pierwszy i zniosłaby wymóg na nasłuchu wystawionym poza pętlę zwrotną.
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
				"bramki nie ma już ani jednego pytania.",
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
