// Odpowiedzialność pliku: konto nadawcze platformy wpisane w rdzeń na stałe,
// żeby świeża instalacja wysyłała kody potwierdzenia bez żadnej nastawy.
package konfiguracja

// Konto nadawcze Danaco. Skrzynka `noreply@danaco-group.pl` stoi na własnym
// serwerze pocztowym platformy (Stalwart, usługa `danaco-hotmail` na maszynie
// wdrożenia) i jest jedyną drogą, którą rdzeń pisze listy systemowe:
// potwierdzenie adresu przy rejestracji i odzyskanie konta.
//
// Wartości stoją tu, a nie w pliku środowiska, bo instalacja bez nich wchodzi
// w drogę „bez poczty": zakłada konto, zostawia adres niepotwierdzony i odcina
// odzyskanie listem. Rozstrzygnięcie Właściciela z 29 sierpnia 2026: konto
// nadawcze ma jechać z pakietem, nie być czynnością wdrożeniową.
//
// Zmienne `DANACO_NADAWCA_*` i nastawy `mailer.*` te wartości nadpisują —
// instalacja u klienta z własną skrzynką nie musi ruszać kodu.
const (
	// NadawcaHostDomyslny to serwer wysyłkowy platformy.
	NadawcaHostDomyslny = "mail.danaco-group.pl"
	// NadawcaPortDomyslny to port zgłoszenia z szyfrowaniem STARTTLS.
	NadawcaPortDomyslny = 587
	// NadawcaAdresDomyslny to adres, z którego przychodzą listy systemowe.
	NadawcaAdresDomyslny = "noreply@danaco-group.pl"
	// NadawcaUzytkownikDomyslny to nazwa logowania skrzynki nadawczej.
	NadawcaUzytkownikDomyslny = "noreply@danaco-group.pl"
	// NadawcaNazwaDomyslna staje w polu nadawcy listu obok adresu.
	NadawcaNazwaDomyslna = "Danaco Console"
)

/*
NadawcaSekretWbudowany niesie hasło skrzynki nadawczej wpisane przy składaniu
pakietu. Puste znaczy pakiet złożony bez hasła — rdzeń bierze je wtedy ze
zmiennej `DANACO_NADAWCA_SEKRET`, a bez niej wchodzi w drogę bez poczty.

Hasła nie ma w źródle i nie może w nim być: wartość wchodzi przy składaniu
przełącznikiem konsolidatora (`-ldflags -X`), tą samą drogą, którą instalator
dostaje poświadczenia kanału wydań. Repozytorium go nie niesie, więc kopia
źródeł nie daje dostępu do skrzynki platformy.
*/
var NadawcaSekretWbudowany = ""
