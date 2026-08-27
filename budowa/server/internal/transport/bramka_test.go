package transport

import (
	"context"
	"strings"
	"testing"

	"danacoconsole/server/internal/protocol"
	"danacoconsole/shared"
)

// Straż bramki jest jedyną granicą bezpieczeństwa tego pakietu i jedynym miejscem odmowy produktu.

// ujscieSprawdzianu jest najuboższym ujściem, jakie straż widzi: identyfikatorem
// i niczym więcej. Straż pyta wyłącznie o identyfikator, więc bogatsze ujście
// mierzyłoby coś innego niż straż.
type ujscieSprawdzianu struct {
	id string
}

func (u ujscieSprawdzianu) Id() string                    { return u.id }
func (u ujscieSprawdzianu) Konto() string                 { return KontoDomyslne }
func (u ujscieSprawdzianu) PrzypiszKonto(string)          {}
func (u ujscieSprawdzianu) Tozsamosc() Tozsamosc          { return Tozsamosc{} }
func (u ujscieSprawdzianu) Wyslij(protocol.Koperta) error { return nil }

// rdzenBezStanuBramki nie zna rozszerzenia StanBramki — tak wygląda rdzeń, który
// nie umie odpowiedzieć „czy to gniazdo przeszło bramkę".
type rdzenBezStanuBramki struct{}

func (rdzenBezStanuBramki) Obsluz(context.Context, protocol.Request, Ujscie) protocol.Koperta {
	return protocol.Koperta{}
}

// rdzenZeStanemBramki jest rdzeniem, który odpowiada na pytanie o więź gniazda z ważną sesją tej bramki.
type rdzenZeStanemBramki struct {
	zwiazane map[string]bool
}

func (rdzenZeStanemBramki) Obsluz(context.Context, protocol.Request, Ujscie) protocol.Koperta {
	return protocol.Koperta{}
}

func (r rdzenZeStanemBramki) PolaczenieZwiazane(id string) bool { return r.zwiazane[id] }

// TestWymogLogowaniaRozstrzygaAdresAlboWskazanie sprawdza regułę wprost:
// wskazanie Operatora wygrywa w obie strony, a jego brak oddaje głos adresowi
// nasłuchu. Domyślna ma być bezpieczna, więc nasłuch szerszy bez wskazania musi
// dać wymóg.
func TestWymogLogowaniaRozstrzygaAdresAlboWskazanie(t *testing.T) {
	prawda, falsz := true, false

	przypadki := []struct {
		nazwa    string
		adres    string
		nastawa  *bool
		wymagany bool
	}{
		{"pętla zwrotna bez wskazania", "127.0.0.1", nil, false},
		{"pętla zwrotna nazwana bez wskazania", "localhost", nil, false},
		{"pętla zwrotna szóstej wersji bez wskazania", "::1", nil, false},
		{"nasłuch szerszy bez wskazania", "0.0.0.0", nil, true},
		{"adres maszyny bez wskazania", "51.75.62.180", nil, true},
		{"pętla zwrotna ze wskazaniem wymogu", "127.0.0.1", &prawda, true},
		{"nasłuch szerszy ze zniesieniem wymogu", "0.0.0.0", &falsz, false},
	}

	for _, przypadek := range przypadki {
		t.Run(przypadek.nazwa, func(t *testing.T) {
			if wymagany := wymogLogowania(przypadek.adres, przypadek.nastawa); wymagany != przypadek.wymagany {
				t.Errorf("wymóg logowania rozstrzygnięty jako %t, oczekiwane %t",
					wymagany, przypadek.wymagany)
			}
		})
	}
}

// Metoda TestStrazPrzepuszczaWylacznieWedlugReguly przechodzi wszystkie układy trzech wejść straży bramki.
func TestStrazPrzepuszczaWylacznieWedlugReguly(t *testing.T) {
	zeStanem := rdzenZeStanemBramki{zwiazane: map[string]bool{"pol-zwiazane": true}}

	przypadki := []struct {
		nazwa       string
		wymagana    bool
		rdzen       Rdzen
		komenda     shared.MessageType
		ujscie      Ujscie
		przepuszcza bool
	}{
		{"straż wyłączona przepuszcza komendę dowolną", false, zeStanem,
			shared.CommandTerminalCommandExec, ujscieSprawdzianu{"pol-obce"}, true},
		{"powitanie przechodzi zawsze", true, zeStanem,
			shared.CommandConnectionHello, ujscieSprawdzianu{"pol-obce"}, true},
		{"logowanie przechodzi zawsze", true, zeStanem,
			shared.CommandAuthLogin, ujscieSprawdzianu{"pol-obce"}, true},
		{"rejestracja przechodzi zawsze", true, zeStanem,
			shared.CommandAuthRegister, ujscieSprawdzianu{"pol-obce"}, true},
		{"potwierdzenie adresu przechodzi zawsze", true, zeStanem,
			shared.CommandAuthVerify, ujscieSprawdzianu{"pol-obce"}, true},
		{"prośba o odzyskanie przechodzi zawsze", true, zeStanem,
			shared.CommandAuthRecover, ujscieSprawdzianu{"pol-obce"}, true},
		{"ustawienie hasła po odzyskaniu przechodzi zawsze", true, zeStanem,
			shared.CommandAuthReset, ujscieSprawdzianu{"pol-obce"}, true},
		{"przedłużenie sesji przechodzi zawsze", true, zeStanem,
			shared.CommandAuthTokenRefresh, ujscieSprawdzianu{"pol-obce"}, true},
		{"gniazdo związane przechodzi", true, zeStanem,
			shared.CommandTerminalCommandExec, ujscieSprawdzianu{"pol-zwiazane"}, true},
		{"gniazdo nieprzedstawione nie przechodzi", true, zeStanem,
			shared.CommandTerminalCommandExec, ujscieSprawdzianu{"pol-obce"}, false},
		{"rdzeń bez stanu bramki nie przepuszcza", true, rdzenBezStanuBramki{},
			shared.CommandTerminalCommandExec, ujscieSprawdzianu{"pol-zwiazane"}, false},
		{"brak ujścia nie przepuszcza", true, zeStanem,
			shared.CommandTerminalCommandExec, nil, false},
	}

	for _, przypadek := range przypadki {
		t.Run(przypadek.nazwa, func(t *testing.T) {
			straz := straznikBramki{wymagana: przypadek.wymagana}
			if przepuszcza := straz.przepusc(przypadek.rdzen, przypadek.komenda, przypadek.ujscie); przepuszcza != przypadek.przepuszcza {
				t.Errorf("straż rozstrzygnęła %t, oczekiwane %t", przepuszcza, przypadek.przepuszcza)
			}
		})
	}
}

// TestKomendyWejsciaToWylacznieBramka pilnuje wykazu komend wejścia, który sam siebie nazywa wyczerpującym.
func TestKomendyWejsciaToWylacznieBramka(t *testing.T) {
	oczekiwane := map[shared.MessageType]bool{
		shared.CommandConnectionHello:  true,
		shared.CommandAuthLogin:        true,
		shared.CommandAuthRegister:     true,
		shared.CommandAuthVerify:       true,
		shared.CommandAuthRecover:      true,
		shared.CommandAuthReset:        true,
		shared.CommandAuthTokenRefresh: true,
	}
	if len(komendyWejscia) != len(oczekiwane) {
		t.Errorf("wykaz komend wejścia ma %d pozycji, oczekiwane %d",
			len(komendyWejscia), len(oczekiwane))
	}
	for komenda := range komendyWejscia {
		if !oczekiwane[komenda] {
			t.Errorf("komenda %q wykonuje się przed przejściem przez bramkę", komenda)
		}
	}
}

// TestOdmowaBramkiNiesieKodOtwierajacyOknoLogowania sprawdza kod, po którym
// klient rozpoznaje, że ma pokazać okno logowania zamiast błędu komendy.
func TestOdmowaBramkiNiesieKodOtwierajacyOknoLogowania(t *testing.T) {
	zadanie := protocol.Request{
		Komenda: shared.CommandTerminalCommandExec,
		Id:      "jeden",
	}
	odmowa := odmowaBezBramki(zadanie)

	if odmowa.Id != "jeden" {
		t.Errorf("odmowa zgubiła identyfikator żądania: %q", odmowa.Id)
	}
	if odmowa.Status == nil || *odmowa.Status != shared.EnvelopeStatusError {
		t.Error("odmowa bez stanu błędu")
	}
	if odmowa.Error == nil || odmowa.Error.Code != shared.ErrorCodeNotAuthenticated {
		t.Errorf("odmowa niesie %+v, oczekiwany kod %q", odmowa.Error, shared.ErrorCodeNotAuthenticated)
	}
	if odmowa.Error != nil && !strings.Contains(odmowa.Error.Message, "zaloguj") {
		t.Errorf("odmowa nie podaje drogi naprawy: %q", odmowa.Error.Message)
	}
}

// TestAdresPustyZnaczyPetleZwrotna pilnuje wartości domyślnej, która jest tu
// rozstrzygnięciem bezpieczeństwa, nie wygodą. Pusty adres znaczy dla net.Listen
// wszystkie interfejsy, więc brak normalizacji wystawiałby rdzeń przez
// przeoczenie.
func TestAdresPustyZnaczyPetleZwrotna(t *testing.T) {
	ustawienia := Ustawienia{}.zNormalizowane()
	if ustawienia.Adres != adresDomyslny {
		t.Errorf("brak wskazania adresu dał nasłuch na %q, oczekiwana pętla zwrotna", ustawienia.Adres)
	}
	if wymogLogowania(ustawienia.Adres, ustawienia.WymogLogowania) {
		t.Error("nasłuch domyślny wymaga logowania — Operator na własnej maszynie dostał pytanie")
	}

	szeroki := Ustawienia{WszystkieInterfejsy: true}.zNormalizowane()
	if szeroki.Adres != "" {
		t.Errorf("wskazanie wszystkich interfejsów dało adres %q", szeroki.Adres)
	}
	if !wymogLogowania(szeroki.Adres, szeroki.WymogLogowania) {
		t.Error("nasłuch na wszystkich interfejsach nie wymaga logowania")
	}
}

// TestNiekompletnaParaTlsZatrzymujeStart sprawdza jedyne miejsce w tym pakiecie,
// które celowo nie wstaje. Praca otwartym tekstem po wskazaniu certyfikatu
// byłaby cichym zejściem poniżej wskazania Operatora.
func TestNiekompletnaParaTlsZatrzymujeStart(t *testing.T) {
	if err := (Ustawienia{CertyfikatTLS: "cert.pem"}).sprawdzTLS(); err == nil {
		t.Error("certyfikat bez klucza został przyjęty")
	}
	if err := (Ustawienia{KluczTLS: "klucz.pem"}).sprawdzTLS(); err == nil {
		t.Error("klucz bez certyfikatu został przyjęty")
	}
	if err := (Ustawienia{}).sprawdzTLS(); err != nil {
		t.Errorf("para pusta została odrzucona: %v", err)
	}
	if err := (Ustawienia{CertyfikatTLS: "cert.pem", KluczTLS: "klucz.pem"}).sprawdzTLS(); err != nil {
		t.Errorf("para kompletna została odrzucona: %v", err)
	}
}

// TestPochodzeniaWskazaneDopisujaSieDoWlasnych pilnuje reguły z opisu:
// wystawienie pod domenę nie jest powodem, żeby produkt przestał wpuszczać
// własną powłokę. Zastąpienie wykazu odcięłoby Operatora od jego własnego okna.
func TestPochodzeniaWskazaneDopisujaSieDoWlasnych(t *testing.T) {
	serwer := Nowy(Ustawienia{PochodzeniaDozwolone: []string{"https://konsola.example", "  ", ""}})
	wzorce := serwer.pochodzeniaDozwolone()

	zbior := make(map[string]bool, len(wzorce))
	for _, wzorzec := range wzorce {
		zbior[wzorzec] = true
	}
	for _, wlasne := range pochodzeniaWlasne {
		if !zbior[wlasne] {
			t.Errorf("pochodzenie własne %q zniknęło po wskazaniu wykazu", wlasne)
		}
	}
	if !zbior["https://konsola.example"] {
		t.Error("pochodzenie wskazane nie weszło do wykazu")
	}
	if zbior[""] || zbior["  "] {
		t.Error("puste pochodzenie weszło do wykazu")
	}
}

// TestWzorcePochodzenBezGolegoNawiasu pilnuje szczegółu, który potrafi odciąć
// pochodzenia sprawdzane po nim: dopasowanie idzie przez path.Match, gdzie
// nawias kwadratowy otwiera klasę znaków, a wzorzec wadliwy przerywa przegląd
// całego wykazu.
func TestWzorcePochodzenBezGolegoNawiasu(t *testing.T) {
	for _, wzorzec := range pochodzeniaWlasne {
		if !strings.Contains(wzorzec, "[") && !strings.Contains(wzorzec, "]") {
			continue
		}
		if strings.Contains(wzorzec, `\[`) && strings.Contains(wzorzec, `\]`) {
			continue
		}
		t.Errorf("wzorzec %q niesie goły nawias kwadratowy — path.Match odrzuci cały wykaz", wzorzec)
	}
}
