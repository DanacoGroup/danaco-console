package konfig

// Wartości ustawień izolacji. Odpowiadają dosłownie kolumnom
// regula_izolacji_kontekstu.wymiar oraz regula_izolacji_technicznej.zakres.
const (
	izolacjaOdrebna       = "odrebna"
	IzolacjaWspoldzielona = "wspoldzielona"
	IzolacjaWlaczona      = "wlaczony"
	izolacjaWylaczona     = "wylaczony"
)

// Klucze jedenastu punktów izolacji: trzy wymiary kontekstu i osiem zakresów
// technicznych. Każdy rozstrzygany osobno — polityka zasięgu jest zbiorem
// niezależnych ustawień, nie wyborem jednej opcji z listy.
const (
	KluczIzolacjaHistoria = "izolacja_historia"
	KluczIzolacjaPamiec   = "izolacja_pamiec"
	KluczIzolacjaKontekst = "izolacja_kontekst"

	KluczIzolacjaKatalogRoboczy    = "izolacja_katalog_roboczy_sesji"
	KluczIzolacjaSrodowiskoProcesu = "izolacja_srodowisko_procesu"
	KluczIzolacjaKatalogDanych     = "izolacja_katalog_danych_modelu"
	KluczIzolacjaDostepSieciowy    = "izolacja_dostep_sieciowy"
	KluczIzolacjaPliki             = "izolacja_odczyt_zapis_plikow"
	KluczIzolacjaKontoIToken       = "izolacja_konto_i_token"
	KluczIzolacjaModelProcesu      = "izolacja_model_procesu"
	KluczIzolacjaSerwerWykonania   = "izolacja_serwer_wykonania"
)

// definicjeIzolacji zwraca jedenaście punktów izolacji z wartościami stanu
// wyjściowego platformy: kontekst odrębny, żaden zakres techniczny niewłączony.
// Wszelka dalsza izolacja jest decyzją Operatora, nigdy ustawieniem narzuconym.
func definicjeIzolacji() []Definicja {
	return append(definicjeIzolacjiKontekstu(), definicjeIzolacjiTechnicznej()...)
}

func definicjeIzolacjiKontekstu() []Definicja {
	return []Definicja{
		{
			Klucz: KluczIzolacjaHistoria, Domyslna: izolacjaOdrebna, Rodzaj: RodzajTekst,
			Objasnienie: "Zapis wymiany wiadomości. Odrębna: nowe okno zaczyna z pustą historią. " +
				"Współdzielona: ten sam zapis widoczny w kilku oknach lub modułach.",
		},
		{
			Klucz: KluczIzolacjaPamiec, Domyslna: izolacjaOdrebna, Rodzaj: RodzajTekst,
			Objasnienie: "Pamięć długoterminowa zasięgu. Odrębna: pamięć jednego zasięgu " +
				"niewidoczna w innym. Współdzielona: jeden zasób zasila kilka zasięgów.",
		},
		{
			Klucz: KluczIzolacjaKontekst, Domyslna: izolacjaOdrebna, Rodzaj: RodzajTekst,
			Objasnienie: "Bieżący stan roboczy: aktywne pliki, projekt, załączniki, zmienne. " +
				"Współdzielony przenosi się między oknami bez przeładowania.",
		},
	}
}

func definicjeIzolacjiTechnicznej() []Definicja {
	zakresy := []struct{ klucz, opis string }{
		{KluczIzolacjaKatalogRoboczy, "Fizyczny katalog plików procesu. Włączony: własny katalog, " +
			"niewidoczny dla innych sesji zasięgu."},
		{KluczIzolacjaSrodowiskoProcesu, "Zmienne środowiskowe i kontekst uruchomieniowy. " +
			"Włączony: własny zestaw zmiennych zamiast wspólnego środowiska serwera."},
		{KluczIzolacjaKatalogDanych, "Dane pomocnicze kanału modelu: konfiguracja, dane tymczasowe, " +
			"ustawienia dostawcy. Włączony: własny katalog danych modelu."},
		{KluczIzolacjaDostepSieciowy, "Połączenia wychodzące: API, strony, serwery MCP, rozszerzenia. " +
			"Włączony: odrębny, ograniczony dostęp sieciowy."},
		{KluczIzolacjaPliki, "Uprawnienia do plików poza własnym katalogiem roboczym. " +
			"Włączony: dostęp wyłącznie do ścieżek jawnie dozwolonych."},
		{KluczIzolacjaKontoIToken, "Token dostępu, klucz API, dane logowania kanału modelu. " +
			"Włączony: własne dane dostępowe zamiast platformowych (odwołania, nie sekrety)."},
		{KluczIzolacjaModelProcesu, "Instancja procesu wykonawczego modelu. Włączony: własna, " +
			"niezależna instancja zamiast wspólnej puli."},
		{KluczIzolacjaSerwerWykonania, "Serwer wykonania procesu — istotne przy kanale zdalnym. " +
			"Włączony: serwer dedykowany zamiast współdzielonego."},
	}
	definicje := make([]Definicja, 0, len(zakresy))
	for _, zakres := range zakresy {
		definicje = append(definicje, Definicja{
			Klucz:       zakres.klucz,
			Domyslna:    izolacjaWylaczona,
			Rodzaj:      RodzajTekst,
			Objasnienie: zakres.opis + " Stan wyjściowy platformy: wyłączony.",
		})
	}
	return definicje
}
