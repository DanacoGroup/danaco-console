// Pakiet konfiguracja ustala ustawienia rdzenia Danaco Console
// na podstawie wartości domyślnych, zmiennych środowiska i argumentów wywołania.
package konfiguracja

import "fmt"

// PortDomyslny to port nasłuchu rdzenia przyjmowany bez wskazania Operatora
// w argumentach uruchomienia.
const PortDomyslny = 17870

// Konfiguracja to komplet ustawień rdzenia ustalony w chwili startu procesu
// i niezmienny do jego końca.
type Konfiguracja struct {
	Rola          Rola   // która część rdzenia pracuje w tym procesie
	Port          int    // port nasłuchu rdzenia
	KatalogDanych string // katalog danych rdzenia
	// KatalogKlienta wskazuje pakiet interfejsu serwowany obok gniazda; brak
	// nie wstrzymuje startu.
	KatalogKlienta string
	// KatalogProfili wskazuje katalog profili kanału głównego; wartość pusta
	// jest dopuszczalna.
	KatalogProfili string

	// Ustawienia brzegu transportu: adres nasłuchu, TLS i wykaz pochodzeń.

	// Adres wskazuje interfejs nasłuchu. Puste = pętla zwrotna.
	Adres string

	// Konto nadawcze platformy, nie skrzynka operatora; nim idą listy
	// potwierdzenia i odzyskania konta.
	NadawcaHost       string
	NadawcaPort       int
	NadawcaUzytkownik string
	NadawcaSekret     string
	NadawcaAdres      string
	AdresKonsoli      string
	NadawcaNazwa      string
	// NadawcaStartTLS ma trzy stany: nil oznacza szyfrowanie włączone, false
	// zejście do tekstu otwartego.
	NadawcaStartTLS *bool
	// WszystkieInterfejsy wystawia nasłuch na wszystkich interfejsach, nie
	// tylko na pętli zwrotnej.
	WszystkieInterfejsy bool
	// CertyfikatTLS i KluczTLS wskazują parę plików warstwy TLS
	// przełączających nasłuch na wss.
	CertyfikatTLS string
	KluczTLS      string
	// WymogLogowania jest dźwignią operatora nad strażą bramki warstwy
	// transportu; ma trzy stany.
	WymogLogowania *bool
	// PochodzeniaDozwolone dopisuje wzorce nagłówka Origin przyjmowane przy
	// nawiązaniu gniazda.
	PochodzeniaDozwolone []string
}

// Domyslna zwraca konfigurację obowiązującą przy braku jakiegokolwiek ustawienia
// (brak ustawienia = wartość domyślna; brak nie blokuje startu).
func Domyslna() Konfiguracja {
	return Konfiguracja{
		Rola:           RolaWszystko,
		Port:           PortDomyslny,
		KatalogDanych:  KatalogDanychDomyslny(),
		KatalogKlienta: KatalogKlientaDomyslny(),
		KatalogProfili: "",
		// Konto nadawcze platformy jedzie z rdzeniem: świeża instalacja ma
		// wysyłać kody potwierdzenia bez żadnej nastawy wdrożeniowej.
		// Zmienne `DANACO_NADAWCA_*` te wartości nadpisują.
		NadawcaHost:       NadawcaHostDomyslny,
		NadawcaPort:       NadawcaPortDomyslny,
		NadawcaUzytkownik: NadawcaUzytkownikDomyslny,
		NadawcaSekret:     NadawcaSekretWbudowany,
		NadawcaAdres:      NadawcaAdresDomyslny,
		AdresKonsoli:      AdresKonsoliDomyslny,
		NadawcaNazwa:      NadawcaNazwaDomyslna,
	}
}

// Opis zwraca jednowierszowy zapis konfiguracji przeznaczony do dziennika
// startu i diagnostyki rdzenia.
func (k Konfiguracja) Opis() string {
	return fmt.Sprintf("rola=%s port=%d dane=%s", k.Rola, k.Port, k.KatalogDanych)
}
