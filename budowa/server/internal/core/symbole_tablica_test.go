package core

import (
	"strings"
	"testing"

	"danacoconsole/shared"
)

// Sprawdziany tablicy znaków mierzą, czy znak wskazany trzema drogami jest tym samym znakiem.

// TestSymbolTablicaNiesieZnakiWymienione mierzy wprost wykaz znaków wymaganych w siedmiu grupach tablicy.
func TestSymbolTablicaNiesieZnakiWymienione(t *testing.T) {
	tablica := symbolTablica()
	poZnaku := make(map[string]shared.StudioSymbol, len(tablica))
	for _, znak := range tablica {
		poZnaku[znak.Character] = znak
	}
	wymagane := []struct{ znak, grupa string }{
		{"§", symbolGrupaPrawnicze},
		{"¶", symbolGrupaPrawnicze},
		{"©", symbolGrupaPrawnicze},
		{"®", symbolGrupaPrawnicze},
		{"™", symbolGrupaPrawnicze},
		{"–", symbolGrupaInterpunkcja}, // półpauza
		{"—", symbolGrupaInterpunkcja}, // pauza
		{"„", symbolGrupaInterpunkcja},
		{"”", symbolGrupaInterpunkcja},
		{" ", symbolGrupaInterpunkcja}, // twarda spacja
		{"‑", symbolGrupaInterpunkcja}, // twardy dywiz
		{"­", symbolGrupaInterpunkcja}, // znak podziału wyrazu
		{"€", symbolGrupaWaluty},
		{"α", symbolGrupaGreckie},
		{"Ω", symbolGrupaGreckie},
		{"→", symbolGrupaStrzalki},
		{"±", symbolGrupaMatematyczne},
		{"´", symbolGrupaDiakrytyczne},
	}
	for _, pozycja := range wymagane {
		znak, jest := poZnaku[pozycja.znak]
		if !jest {
			t.Errorf("tablica znaków nie niesie znaku %q", pozycja.znak)
			continue
		}
		if znak.Category == nil || *znak.Category != pozycja.grupa {
			t.Errorf("znak %q stoi w grupie %v, a należy do grupy %q",
				pozycja.znak, znak.Category, pozycja.grupa)
		}
		if strings.TrimSpace(znak.Name) == "" {
			t.Errorf("znak %q nie ma nazwy, więc nie da się go wskazać nazwą", pozycja.znak)
		}
		if !strings.HasPrefix(znak.Code, "U+") {
			t.Errorf("znak %q ma kod %q, a kod zapisuje się jako U+XXXX",
				pozycja.znak, znak.Code)
		}
	}
	// Wszystkie siedem grup, które zlecenie wymienia, musi mieć swoje znaki.
	grupy := symbolGrupy()
	if len(grupy) != 7 {
		t.Errorf("grup znaków jest %d, a zlecenie wymienia siedem: %v", len(grupy), grupy)
	}
	// Żaden znak nie powtarza się dwa razy, bo powtórka dawałaby dwie pozycje robiące to samo.
	widziane := map[string]int{}
	for _, znak := range tablica {
		widziane[znak.Character]++
	}
	for znak, ile := range widziane {
		if ile > 1 {
			t.Errorf("znak %q stoi w tablicy %d razy", znak, ile)
		}
	}
}

// TestSymbolKodIZnakSaOdwracalne mierzy rachunek punktu kodowego w obie strony: ze znaku na kod i z kodu na znak.
func TestSymbolKodIZnakSaOdwracalne(t *testing.T) {
	if kod := symbolKodZnaku("§"); kod != "U+00A7" {
		t.Errorf("paragraf ma kod %q, a ma mieć U+00A7", kod)
	}
	// Zapisy, w jakich Operator wkleja kod. Wszystkie mają dać ten sam znak.
	for _, zapis := range []string{"U+00A7", "u+00a7", "0x00A7", "00a7", "A7"} {
		znak, jest := symbolZnakZKodu(zapis)
		if !jest || znak != "§" {
			t.Errorf("zapis %q dał %q (odczytany: %v), a ma dać paragraf", zapis, znak, jest)
		}
	}
	// Znak złożony z kilku punktów kodowych oddaje je rozdzielone spacją, bo widoczny jest jeden znak.
	zlozony := "á"
	if kod := symbolKodZnaku(zlozony); kod != "U+0061 U+0301" {
		t.Errorf("znak złożony ma kod %q", kod)
	}
	if znak, jest := symbolZnakZKodu("U+0061 U+0301"); !jest || znak != zlozony {
		t.Errorf("kod złożony dał %q (odczytany: %v)", znak, jest)
	}
	for _, zapis := range []string{"", "   ", "U+", "paragraf", "0xZZZZ", "U+110000"} {
		if _, jest := symbolZnakZKodu(zapis); jest {
			t.Errorf("zapis %q przeszedł jako punkt kodowy", zapis)
		}
	}
}

// TestSymbolRozstrzygnijTrzyDrogi mierzy trzy drogi wskazania znaku: znakiem,
// kodem i nazwą. Polecenie modelu „wstaw tu paragraf" idzie drogą trzecią.
func TestSymbolRozstrzygnijTrzyDrogi(t *testing.T) {
	poZnaku, err := symbolRozstrzygnij(shared.StudioSymbolInsertRequest{
		Character: postacWskaznikTekstu("§"),
	})
	if err != nil {
		t.Fatalf("wskazanie znakiem odmówiło: %v", err)
	}
	if poZnaku.Name != "paragraf" {
		t.Errorf("znak wskazany znakiem dostał nazwę %q", poZnaku.Name)
	}

	poKodzie, err := symbolRozstrzygnij(shared.StudioSymbolInsertRequest{
		Code: postacWskaznikTekstu("U+00A7"),
	})
	if err != nil {
		t.Fatalf("wskazanie kodem odmówiło: %v", err)
	}
	poNazwie, err := symbolRozstrzygnij(shared.StudioSymbolInsertRequest{
		Name: postacWskaznikTekstu("paragraf"),
	})
	if err != nil {
		t.Fatalf("wskazanie nazwą odmówiło: %v", err)
	}
	if poZnaku.Character != poKodzie.Character || poKodzie.Character != poNazwie.Character {
		t.Errorf("trzy drogi dały trzy różne znaki: %q, %q, %q",
			poZnaku.Character, poKodzie.Character, poNazwie.Character)
	}

	// Znak spoza tablicy jest znakiem prawdziwym — Unicode ma ich więcej niż mieści wykaz okna.
	spozaTablicy, err := symbolRozstrzygnij(shared.StudioSymbolInsertRequest{
		Code: postacWskaznikTekstu("U+2603"),
	})
	if err != nil {
		t.Fatalf("znak spoza tablicy odmówił: %v", err)
	}
	if spozaTablicy.Character != "☃" {
		t.Errorf("znak spoza tablicy wyszedł jako %q", spozaTablicy.Character)
	}

	// Nazwa dwuznaczna nazywa dwuznaczność, zamiast wstawiać pierwszy napotkany.
	_, err = symbolRozstrzygnij(shared.StudioSymbolInsertRequest{
		Name: postacWskaznikTekstu("strzałka"),
	})
	if err == nil {
		t.Error("nazwa wskazująca wiele znaków przeszła bez odmowy")
	} else if !strings.Contains(err.Error(), "więcej niż jeden") {
		t.Errorf("odmowa nie nazywa dwuznaczności: %v", err)
	}
	// Nazwa pełna trafia w jeden znak, choć jest zawarta w nazwach innych.
	strzalka, err := symbolRozstrzygnij(shared.StudioSymbolInsertRequest{
		Name: postacWskaznikTekstu("strzałka w prawo"),
	})
	if err != nil || strzalka.Character != "→" {
		t.Errorf("nazwa pełna dała %q, błąd %v", strzalka.Character, err)
	}

	// Żądanie bez ani jednego wskazania jest odmową nazwaną.
	if _, err := symbolRozstrzygnij(shared.StudioSymbolInsertRequest{}); err == nil {
		t.Error("wstawienie znaku bez wskazania znaku przeszło bez odmowy")
	}
	// Pole character niosące napis, a nie znak, jest odmową ze wskazaniem drogi.
	_, err = symbolRozstrzygnij(shared.StudioSymbolInsertRequest{
		Character: postacWskaznikTekstu("całe zdanie do wstawienia"),
	})
	if err == nil || !strings.Contains(err.Error(), "studio.text.edit") {
		t.Errorf("napis w polu character: %v", err)
	}
	// Nazwa, której tablica nie niesie, wskazuje drogę wyjścia przez kod.
	_, err = symbolRozstrzygnij(shared.StudioSymbolInsertRequest{
		Name: postacWskaznikTekstu("znak, którego nie ma"),
	})
	if err == nil || !strings.Contains(err.Error(), "code") {
		t.Errorf("nazwa nieznana nie wskazuje drogi wyjścia: %v", err)
	}
}

// TestSymbolPasujeSzukaPoNazwieIKodzie mierzy wyszukiwanie wymagane wprost: po nazwie i po kodzie jednocześnie.
func TestSymbolPasujeSzukaPoNazwieIKodzie(t *testing.T) {
	paragraf := shared.StudioSymbol{
		Code: "U+00A7", Character: "§", Name: "paragraf",
	}
	for _, szukane := range []string{"paragraf", "para", "u+00a7", "00a7", "§"} {
		if !symbolPasuje(paragraf, strings.ToLower(szukane)) {
			t.Errorf("szukane %q nie trafiło w paragraf", szukane)
		}
	}
	for _, szukane := range []string{"euro", "u+20ac"} {
		if symbolPasuje(paragraf, szukane) {
			t.Errorf("szukane %q trafiło w paragraf", szukane)
		}
	}
}
