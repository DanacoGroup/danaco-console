package konfig

import "danacoconsole/shared"

// Klucze ustawień wykonania. Odpowiadają kolumnie ustawienie.klucz; ustawialne
// na każdym z ośmiu poziomów, najczęściej na poziomie okna komunikacji.
const (
	kluczTrybUprawnien       = "tryb_uprawnien"
	kluczRolaOkna            = "rola_okna"
	kluczSrodowiskoWykonania = "srodowisko_wykonania"
	kluczHostWykonania       = "host_wykonania"
	kluczKanalModelu         = "kanal_modelu"
	KluczKanalZapasowy       = "kanal_modelu_zapasowy"
	KluczNakladRozumowania   = "naklad_rozumowania"
	// KluczPulapKosztu jest górną granicą kosztu jednego wywołania modelu,
	// wyrażoną w dolarach. Zero — i tak samo wartość pusta — znaczy brak pułapu,
	// nie zakaz wywołania.
	KluczPulapKosztu     = "pulap_kosztu_usd"
	kluczKatalogiRobocze = "katalogi_robocze"
)

// definicjeWykonania zwraca ustawienia sterujące wykonaniem okna komunikacji.
// Wartości domyślne pochodzą ze stałych kontraktu. Wartość pusta znaczy brak
// wskazania na tym poziomie i pozostawia decyzję rejestrowi kanałów albo
// samemu kanałowi.
func definicjeWykonania() []Definicja {
	return []Definicja{
		{
			Klucz: kluczTrybUprawnien, Domyslna: shared.PermissionModeManual, Rodzaj: RodzajTekst,
			Objasnienie: "Tryb uprawnień okna komunikacji — zakres zgody wydanej modelowi " +
				"przed zmianą w systemie. Wartości kontraktu PermissionMode.",
		},
		{
			Klucz: kluczRolaOkna, Domyslna: shared.WindowRoleStandalone, Rodzaj: RodzajTekst,
			Objasnienie: "Rola okna w pętli koordynator–wykonawca. Okno samodzielne " +
				"pozostaje poza pętlą.",
		},
		{
			Klucz: kluczSrodowiskoWykonania, Domyslna: shared.ExecutionEnvLocal, Rodzaj: RodzajTekst,
			Objasnienie: "Zasięg wykonania modelu: urządzenie użytkownika, " +
				"host serwera albo host zdalny. Niezależny od umiejscowienia serwera.",
		},
		{
			Klucz: kluczHostWykonania, Domyslna: "", Rodzaj: RodzajTekst,
			Objasnienie: "Nazwa hosta wykonania przy zasięgu zdalnym — na przykład danaco-system. " +
				"Nazwa hosta jest ustawieniem poziomu zasięgu, nie wartością wyliczenia.",
		},
		{
			Klucz: kluczKanalModelu, Domyslna: "", Rodzaj: RodzajTekst,
			Objasnienie: "Kanał modelu z rejestru kanałów. Brak wskazania znaczy: " +
				"rozstrzyga rejestr kanałów, nie odmowa uruchomienia.",
		},
		{
			Klucz: KluczKanalZapasowy, Domyslna: "", Rodzaj: RodzajTekst,
			Objasnienie: "Kanał modelu używany po niepowodzeniu kanału głównego. Brak wskazania " +
				"znaczy brak zapasu, nie blokadę wywołania.",
		},
		{
			Klucz: KluczNakladRozumowania, Domyslna: "", Rodzaj: RodzajTekst,
			Objasnienie: "Nakład rozumowania modelu — od szybciej do mądrzej. Brak wskazania " +
				"zostawia rozstrzygnięcie kanałowi modelu.",
		},
		{
			Klucz: KluczPulapKosztu, Domyslna: "0", Rodzaj: RodzajLiczba,
			Objasnienie: "Górna granica kosztu jednego wywołania modelu w dolarach, " +
				"przekazywana programowi przełącznikiem --max-budget-usd. Zero znaczy " +
				"brak pułapu — wywołanie idzie bez ograniczenia, nie zostaje wstrzymane.",
		},
		{
			Klucz: kluczKatalogiRobocze, Domyslna: "[]", Rodzaj: RodzajJSON,
			Objasnienie: "Lista katalogów roboczych okna komunikacji. Lista pusta znaczy: " +
				"katalog wskaże Operator przy otwarciu okna.",
		},
	}
}
