// Pakiet konfiguracja ustala ustawienia rdzenia Danaco Console
// na podstawie wartości domyślnych, zmiennych środowiska i argumentów wywołania.
package konfiguracja

import "fmt"

// PortDomyslny to port nasłuchu rdzenia przyjmowany bez wskazania Operatora.
const PortDomyslny = 17870

// Konfiguracja to komplet ustawień rdzenia ustalony w chwili startu procesu.
type Konfiguracja struct {
	Rola          Rola   // która część rdzenia pracuje w tym procesie
	Port          int    // port nasłuchu rdzenia
	KatalogDanych string // katalog danych rdzenia
	// KatalogKlienta wskazuje pakiet interfejsu serwowany obok gniazda.
	// Brak pakietu nie wstrzymuje startu — gniazdo działa bez plików.
	KatalogKlienta string
	// KatalogProfili wskazuje katalog profili kanału głównego. Poświadczenia
	// zostają w profilach na dysku; rdzeń zna wyłącznie odwołanie.
	// Wartość pusta jest dopuszczalna: pula kont startuje wtedy pusta.
	KatalogProfili string

	// Ustawienia brzegu transportu — adres nasłuchu, TLS i wykaz pochodzeń —
	// stoją tu razem z resztą ustawień startu, bo tylko stąd sięgają po nie
	// warstwy wartości domyślnych, zmiennych środowiska i argumentów wywołania.

	// Adres wskazuje interfejs nasłuchu. Puste = pętla zwrotna.
	Adres string

	// Konto nadawcze platformy — nim idą dwa listy systemowe: potwierdzenie
	// adresu przy rejestracji i droga odzyskania konta. Nie jest to skrzynka
	// Operatora: gdyby platforma pisała jego kontem, utrata dostępu do skrzynki
	// odcinałaby drogę odzyskania dokładnie wtedy, gdy jest potrzebna.
	//
	// Brak tych wartości nie wstrzymuje startu rdzenia — wstrzymuje wyłącznie
	// rejestrację, i to odmową nazywającą brak wprost. Rdzeń bez konta
	// nadawczego pracuje dla Operatora już zalogowanego.
	NadawcaHost       string
	NadawcaPort       int
	NadawcaUzytkownik string
	NadawcaSekret     string
	NadawcaAdres      string
	NadawcaNazwa      string
	// NadawcaStartTLS ma trzy stany, stąd wskaźnik: nil = wartość domyślna
	// (szyfrowanie włączone), false = jawne zejście do rozmowy otwartym tekstem
	// dla przekaźnika na tej samej maszynie.
	NadawcaStartTLS *bool
	// WszystkieInterfejsy wystawia nasłuch na wszystkich interfejsach maszyny.
	// Osobne pole, a nie pusty adres: „nie wskazałem" i „chcę wszędzie" to dwa
	// różne zdania i mają wyglądać różnie w miejscu wywołania.
	WszystkieInterfejsy bool
	// CertyfikatTLS i KluczTLS wskazują parę plików warstwy TLS. Wskazanie obu
	// przełącza nasłuch na wss; wskazanie jednego zatrzymuje start, bo cicha
	// praca otwartym tekstem po wskazaniu certyfikatu byłaby zejściem poniżej
	// wskazania. Rozstrzyga to warstwa transportu, tu wartość tylko przechodzi.
	CertyfikatTLS string
	KluczTLS      string
	// WymogLogowania jest dźwignią Operatora nad strażą bramki warstwy
	// transportu. Trzy stany, stąd wskaźnik: nil = rozstrzyga adres nasłuchu
	// (poza pętlą zwrotną wymóg obowiązuje sam z siebie), true = wymóg także na
	// pętli zwrotnej, false = wymóg zniesiony, a dziennik mówi o tym wprost.
	// Rozstrzyga warstwa transportu, tu wartość tylko przechodzi. Nastawa
	// obowiązuje od startu procesu: warstwa nasłuchu nie przyjmuje zmiany
	// wymogu na żywo.
	WymogLogowania *bool
	// PochodzeniaDozwolone dopisuje wzorce nagłówka Origin przyjmowane przy
	// nawiązaniu gniazda. Pochodzeń własnych produktu nie zastępuje.
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
	}
}

// Opis zwraca jednowierszowy zapis konfiguracji przeznaczony do dziennika.
func (k Konfiguracja) Opis() string {
	return fmt.Sprintf("rola=%s port=%d dane=%s", k.Rola, k.Port, k.KatalogDanych)
}
