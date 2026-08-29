package konfiguracja

import (
	"fmt"
	"strconv"
)

// Nazwy zmiennych środowiska ustalających konfigurację rdzenia.
// Nazwy są polskie tak samo jak pola konfiguracji, którym odpowiadają,
// i zgodne ze wzorcem `budowa/.env.example`.
const (
	zmiennaRola           = "DANACO_ROLA"
	zmiennaPort           = "DANACO_PORT"
	zmiennaKatalogDanych  = "DANACO_KATALOG_DANYCH"
	zmiennaKatalogKlienta = "DANACO_KATALOG_KLIENTA"
	zmiennaKatalogProfili = "DANACO_KATALOG_PROFILI"
	// Nazwy zmiennych brzegu transportu: adres nasłuchu, wystawienie na
	// wszystkie interfejsy, certyfikat i klucz TLS, pochodzenia i wymóg
	// logowania.
	zmiennaAdres               = "DANACO_ADRES"
	zmiennaWszystkieInterfejsy = "DANACO_WSZYSTKIE_INTERFEJSY"
	zmiennaCertyfikatTLS       = "DANACO_TLS_CERTYFIKAT"
	zmiennaKluczTLS            = "DANACO_TLS_KLUCZ"
	zmiennaPochodzenia         = "DANACO_POCHODZENIA"
	zmiennaWymogLogowania      = "DANACO_WYMOG_LOGOWANIA"

	// Konto nadawcze platformy, nie skrzynka operatora. Idą nim dwa listy
	// systemowe: potwierdzenie adresu przy rejestracji i droga odzyskania
	// konta. Jest osobne od skrzynki operatora, aby utrata dostępu do niej
	// nie odcinała drogi odzyskania.
	zmiennaNadawcaHost       = "DANACO_NADAWCA_HOST"
	zmiennaNadawcaPort       = "DANACO_NADAWCA_PORT"
	zmiennaNadawcaUzytkownik = "DANACO_NADAWCA_UZYTKOWNIK"
	zmiennaNadawcaSekret     = "DANACO_NADAWCA_SEKRET"
	zmiennaNadawcaAdres      = "DANACO_NADAWCA_ADRES"
	zmiennaNadawcaNazwa      = "DANACO_NADAWCA_NAZWA"
	zmiennaNadawcaStartTLS   = "DANACO_NADAWCA_STARTTLS"

	// Publiczny adres Konsoli, pod ktory kieruja odsylacze z listow.
	// Osobny od DANACO_ADRES, ktory niesie gniazdo nasluchu rdzenia.
	zmiennaAdresKonsoli = "DANACO_ADRES_KONSOLI"
)

// ZmienneSrodowiska zwraca nazwy zmiennych czytanych przez rdzeń
// w kolejności prezentacji. Jest to jedyny wykaz tych nazw: zasila pomoc
// wiersza poleceń i sprawdzenie zgodności ze wzorcem .env.example.
func ZmienneSrodowiska() []string {
	return []string{
		zmiennaRola, zmiennaPort, zmiennaKatalogDanych,
		zmiennaKatalogKlienta, zmiennaKatalogProfili,
		zmiennaAdres, zmiennaWszystkieInterfejsy,
		zmiennaCertyfikatTLS, zmiennaKluczTLS, zmiennaPochodzenia,
		zmiennaWymogLogowania,
		zmiennaAdresKonsoli,
		zmiennaNadawcaHost, zmiennaNadawcaPort, zmiennaNadawcaUzytkownik,
		zmiennaNadawcaSekret, zmiennaNadawcaAdres, zmiennaNadawcaNazwa,
		zmiennaNadawcaStartTLS,
	}
}

// zastosujSrodowisko nakłada na konfigurację wartości ze zmiennych środowiska.
// Zmienna nieustawiona lub pusta pozostawia wartość dotychczasową.
func zastosujSrodowisko(kon *Konfiguracja, odczyt func(string) string) error {
	if tekst := odczyt(zmiennaRola); tekst != "" {
		rola, err := RolaZTekstu(tekst)
		if err != nil {
			return fmt.Errorf("%s: %w", zmiennaRola, err)
		}
		kon.Rola = rola
	}
	if tekst := odczyt(zmiennaPort); tekst != "" {
		port, err := strconv.Atoi(tekst)
		if err != nil {
			return fmt.Errorf("%s: wartość %q nie jest liczbą", zmiennaPort, tekst)
		}
		kon.Port = port
	}
	if katalog := odczyt(zmiennaKatalogDanych); katalog != "" {
		kon.KatalogDanych = katalog
	}
	if katalog := odczyt(zmiennaKatalogKlienta); katalog != "" {
		kon.KatalogKlienta = katalog
	}
	if katalog := odczyt(zmiennaKatalogProfili); katalog != "" {
		kon.KatalogProfili = katalog
	}
	if adres := odczyt(zmiennaAdres); adres != "" {
		kon.Adres = adres
	}
	if adres := odczyt(zmiennaAdresKonsoli); adres != "" {
		kon.AdresKonsoli = adres
	}
	if host := odczyt(zmiennaNadawcaHost); host != "" {
		kon.NadawcaHost = host
	}
	if tekst := odczyt(zmiennaNadawcaPort); tekst != "" {
		port, err := strconv.Atoi(tekst)
		if err != nil {
			return fmt.Errorf("%s: wartość %q nie jest liczbą", zmiennaNadawcaPort, tekst)
		}
		kon.NadawcaPort = port
	}
	if uzytkownik := odczyt(zmiennaNadawcaUzytkownik); uzytkownik != "" {
		kon.NadawcaUzytkownik = uzytkownik
	}
	if sekret := odczyt(zmiennaNadawcaSekret); sekret != "" {
		kon.NadawcaSekret = sekret
	}
	if adres := odczyt(zmiennaNadawcaAdres); adres != "" {
		kon.NadawcaAdres = adres
	}
	if nazwa := odczyt(zmiennaNadawcaNazwa); nazwa != "" {
		kon.NadawcaNazwa = nazwa
	}
	if tekst := odczyt(zmiennaNadawcaStartTLS); tekst != "" {
		// Wartość nieczytelna zatrzymuje start, zamiast po cichu znaczyć nie.
		startTLS, err := strconv.ParseBool(tekst)
		if err != nil {
			return fmt.Errorf("%s: wartość %q nie jest wartością logiczną (true|false)",
				zmiennaNadawcaStartTLS, tekst)
		}
		kon.NadawcaStartTLS = &startTLS
	}
	if tekst := odczyt(zmiennaWszystkieInterfejsy); tekst != "" {
		// Wartość nieczytelna zatrzymuje start, zamiast po cichu znaczyć nie,
		// jak przy wystawieniu na sieć.
		wszystkie, err := strconv.ParseBool(tekst)
		if err != nil {
			return fmt.Errorf("%s: wartość %q nie jest wartością logiczną (true|false)",
				zmiennaWszystkieInterfejsy, tekst)
		}
		kon.WszystkieInterfejsy = wszystkie
	}
	if plik := odczyt(zmiennaCertyfikatTLS); plik != "" {
		kon.CertyfikatTLS = plik
	}
	if plik := odczyt(zmiennaKluczTLS); plik != "" {
		kon.KluczTLS = plik
	}
	if tekst := odczyt(zmiennaWymogLogowania); tekst != "" {
		// Wartość nieczytelna zatrzymuje start z tego samego powodu, co
		// wystawienie na wszystkie interfejsy.
		wymog, err := strconv.ParseBool(tekst)
		if err != nil {
			return fmt.Errorf("%s: wartość %q nie jest wartością logiczną (true|false)",
				zmiennaWymogLogowania, tekst)
		}
		kon.WymogLogowania = &wymog
	}
	if wykaz := odczyt(zmiennaPochodzenia); wykaz != "" {
		kon.PochodzeniaDozwolone = wykazPoPrzecinku(wykaz)
	}
	return nil
}
