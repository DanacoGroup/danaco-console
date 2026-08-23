package konfiguracja

import (
	"flag"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

// nazwaProgramu to nazwa binarki rdzenia używana w pomocy wiersza poleceń.
const nazwaProgramu = "danaco-console"

// zastosujArgumenty nakłada na konfigurację wartości z wiersza poleceń.
// Wartości domyślne przełączników pochodzą z warstw wcześniejszych,
// dzięki czemu argument wygrywa ze zmienną środowiska, a zmienna z wartością domyślną.
func zastosujArgumenty(kon *Konfiguracja, argumenty []string) error {
	zestaw := flag.NewFlagSet(nazwaProgramu, flag.ContinueOnError)
	zestaw.SetOutput(os.Stdout)
	zestaw.Usage = func() { wypiszPomoc(zestaw.Output(), zestaw) }

	rola := zestaw.String("role", string(kon.Rola), "rola procesu: "+strings.Join(NazwyRol(), "|"))
	port := zestaw.Int("port", kon.Port, "port nasłuchu rdzenia")
	dane := zestaw.String("dane", kon.KatalogDanych, "katalog danych rdzenia")
	klient := zestaw.String("klient", kon.KatalogKlienta, "katalog pakietu interfejsu")
	profile := zestaw.String("profile", kon.KatalogProfili, "katalog profili kanału głównego")
	adres := zestaw.String("adres", kon.Adres, "interfejs nasłuchu; puste = pętla zwrotna 127.0.0.1")
	wszystkie := zestaw.Bool("wszystkie-interfejsy", kon.WszystkieInterfejsy,
		"nasłuch na wszystkich interfejsach maszyny (wystawienie poza pętlę zwrotną)")
	cert := zestaw.String("tls-certyfikat", kon.CertyfikatTLS, "plik certyfikatu TLS; razem z -tls-klucz włącza wss")
	klucz := zestaw.String("tls-klucz", kon.KluczTLS, "plik klucza TLS; razem z -tls-certyfikat włącza wss")
	// Przełącznik jest tekstem, nie flagą logiczną, bo niesie trzy stany:
	// niewskazany (puste), wymuszony i zniesiony. flag.Bool umiałby dwa i
	// zamieniłby brak wskazania we wskazanie „nie", zdejmując wymóg
	// wystawionemu rdzeniowi przez samo pominięcie przełącznika.
	wymog := zestaw.String("wymog-logowania", tekstWymogu(kon.WymogLogowania),
		"czy połączenie musi przedstawić token bramki: true|false; "+
			"puste = rozstrzyga adres nasłuchu (poza pętlą zwrotną wymóg obowiązuje)")
	pochodzenia := zestaw.String("pochodzenia", strings.Join(kon.PochodzeniaDozwolone, ","),
		"wykaz wzorców Origin przyjmowanych przy nawiązaniu, po przecinku; dopisuje się do pochodzeń własnych")

	if err := zestaw.Parse(argumenty); err != nil {
		return err
	}
	kon.Adres = *adres
	kon.WszystkieInterfejsy = *wszystkie
	kon.CertyfikatTLS = *cert
	kon.KluczTLS = *klucz
	kon.PochodzeniaDozwolone = wykazPoPrzecinku(*pochodzenia)
	wskazany, err := wymogZTekstu(*wymog)
	if err != nil {
		return err
	}
	kon.WymogLogowania = wskazany
	wybrana, err := RolaZTekstu(*rola)
	if err != nil {
		return fmt.Errorf("--role: %w", err)
	}
	kon.Rola = wybrana
	kon.Port = *port
	kon.KatalogDanych = *dane
	kon.KatalogKlienta = *klient
	kon.KatalogProfili = *profile
	return nil
}

// tekstWymogu zapisuje trójstan jako wartość domyślną przełącznika: pusty napis
// znaczy „nie wskazano", więc powtórzone przejście warstw nie zamienia braku
// wskazania we wskazanie.
func tekstWymogu(wymog *bool) string {
	if wymog == nil {
		return ""
	}
	return strconv.FormatBool(*wymog)
}

// wymogZTekstu odczytuje trójstan z przełącznika. Wartość nieczytelna zatrzymuje
// start, zamiast po cichu znaczyć jedno z dwóch.
func wymogZTekstu(tekst string) (*bool, error) {
	tekst = strings.TrimSpace(tekst)
	if tekst == "" {
		return nil, nil
	}
	wymog, err := strconv.ParseBool(tekst)
	if err != nil {
		return nil, fmt.Errorf("--wymog-logowania: wartość %q nie jest wartością logiczną (true|false)", tekst)
	}
	return &wymog, nil
}

// wykazPoPrzecinku rozbija jeden przełącznik na wykaz wzorców. Przełącznik
// powtarzalny wymagałby własnego typu flagi; przecinek wystarcza, bo wzorzec
// Origin przecinka nie zawiera. Człony puste wypadają — wzorzec pusty pasowałby
// do niczego albo, w bibliotece gniazda, do czegokolwiek.
func wykazPoPrzecinku(tekst string) []string {
	wykaz := make([]string, 0)
	for _, czlon := range strings.Split(tekst, ",") {
		if czlon = strings.TrimSpace(czlon); czlon != "" {
			wykaz = append(wykaz, czlon)
		}
	}
	if len(wykaz) == 0 {
		return nil
	}
	return wykaz
}

// wypiszPomoc wypisuje nagłówek, katalog przełączników i nazwy zmiennych środowiska.
func wypiszPomoc(wyjscie io.Writer, zestaw *flag.FlagSet) {
	fmt.Fprintf(wyjscie, "%s — rdzeń Danaco Console\n\n", nazwaProgramu)
	fmt.Fprintf(wyjscie, "Użycie:\n  %s [przełączniki]\n\nPrzełączniki:\n", nazwaProgramu)
	zestaw.PrintDefaults()
	fmt.Fprintf(wyjscie, "\nZmienne środowiska: %s\n", strings.Join(ZmienneSrodowiska(), ", "))
}
