package mail

import (
	"io"
	"mime"
	"mime/multipart"
	netmail "net/mail"
	"os"
	"strings"
	"testing"
	"time"
)

// katalogSzablonow wskazuje katalog 06-poczta-transakcyjna względem pakietu.
const katalogSzablonow = "../.."

// zmienneNiezaufane to nazwy pochodzące spoza rdzenia.
var zmienneNiezaufane = map[string]bool{"original_subject": true, "run_name": true}

func wczytaj(t *testing.T) *Set {
	t.Helper()
	set, err := Load(os.DirFS(katalogSzablonow))
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	return set
}

// wartosci wypełnia wszystkie zmienne szablonu wartościami zastępczymi.
func wartosci(t *testing.T, set *Set, kind Kind) Content {
	t.Helper()
	tpl, err := set.Template(kind)
	if err != nil {
		t.Fatalf("Template: %v", err)
	}
	c := Content{Values: map[string]string{}, Untrusted: map[string]string{}}
	for _, name := range tpl.Vars {
		switch {
		case name == "logo_src" || name == "logo_src_dark":
			// podstawiane przez Build
		case zmienneNiezaufane[name]:
			c.Untrusted[name] = "wartość niezaufana"
		default:
			c.Values[name] = "wartość-" + name
		}
	}
	return c
}

func koperta() Envelope {
	return Envelope{
		From: netmail.Address{Name: "Danaco Console", Address: "noreply@danaco-group.pl"},
		To:   netmail.Address{Address: "operator@przyklad.pl"},
		Date: time.Date(2026, 8, 29, 17, 22, 0, 0, time.UTC),
	}
}

func znaki() Logos {
	return Logos{
		Light: []byte("PNG-jasny"), LightName: "znak-poczty-jasny@2x.png",
		Dark: []byte("PNG-ciemny"), DarkName: "znak-poczty-ciemny@2x.png",
	}
}

func TestLoadKompletSzablonow(t *testing.T) {
	set := wczytaj(t)
	for _, kind := range Kinds() {
		tpl, err := set.Template(kind)
		if err != nil {
			t.Fatalf("%s: %v", kind, err)
		}
		if tpl.HTML == "" || tpl.Text == "" {
			t.Errorf("%s: pusta część szablonu", kind)
		}
		if len(tpl.Vars) == 0 {
			t.Errorf("%s: szablon bez zmiennych", kind)
		}
		if strings.Contains(tpl.HTML, "<!--") {
			t.Errorf("%s: komentarz wewnętrzny pozostał w treści HTML", kind)
		}
		if !strings.Contains(tpl.HTML, "{{logo_src}}") {
			t.Errorf("%s: brak znaku w nagłówku", kind)
		}
	}
}

func TestLoadNieznanyRodzaj(t *testing.T) {
	if _, err := wczytaj(t).Template("list-99-nieistniejacy"); err == nil {
		t.Fatal("oczekiwano błędu dla nieznanego rodzaju")
	}
}

func TestBuildStrukturaMIME(t *testing.T) {
	set := wczytaj(t)
	msg, err := set.Build(KindLoginCode, koperta(), wartosci(t, set, KindLoginCode), znaki())
	if err != nil {
		t.Fatalf("Build: %v", err)
	}

	parsed, err := netmail.ReadMessage(strings.NewReader(string(msg.Data)))
	if err != nil {
		t.Fatalf("wiadomość nie parsuje się: %v", err)
	}
	if got := parsed.Header.Get("Auto-Submitted"); got != "auto-generated" {
		t.Errorf("Auto-Submitted = %q, oczekiwano auto-generated", got)
	}
	if got := parsed.Header.Get("X-Auto-Response-Suppress"); got != "All" {
		t.Errorf("X-Auto-Response-Suppress = %q", got)
	}

	mediaType, params, err := mime.ParseMediaType(parsed.Header.Get("Content-Type"))
	if err != nil {
		t.Fatalf("Content-Type: %v", err)
	}
	if mediaType != "multipart/related" {
		t.Fatalf("typ nadrzędny = %s, oczekiwano multipart/related", mediaType)
	}

	var typyRelated []string
	var contentIDs []string
	var alternative []string
	related := multipart.NewReader(parsed.Body, params["boundary"])
	for {
		part, err := related.NextPart()
		if err == io.EOF {
			break
		}
		if err != nil {
			t.Fatalf("część related: %v", err)
		}
		mt, p, err := mime.ParseMediaType(part.Header.Get("Content-Type"))
		if err != nil {
			t.Fatalf("Content-Type części: %v", err)
		}
		typyRelated = append(typyRelated, mt)
		if id := part.Header.Get("Content-ID"); id != "" {
			contentIDs = append(contentIDs, id)
		}
		if mt == "multipart/alternative" {
			inner := multipart.NewReader(part, p["boundary"])
			for {
				sub, err := inner.NextPart()
				if err == io.EOF {
					break
				}
				if err != nil {
					t.Fatalf("część alternative: %v", err)
				}
				smt, _, err := mime.ParseMediaType(sub.Header.Get("Content-Type"))
				if err != nil {
					t.Fatalf("Content-Type podczęści: %v", err)
				}
				alternative = append(alternative, smt)
			}
		}
	}

	// multipart.Reader odkodowuje quoted-printable i ukrywa nagłówek
	// kodowania, więc obecność kodowania sprawdzamy na surowych bajtach.
	if n := strings.Count(string(msg.Data), "Content-Transfer-Encoding: quoted-printable"); n != 2 {
		t.Errorf("części quoted-printable = %d, oczekiwano 2", n)
	}
	if n := strings.Count(string(msg.Data), "Content-Transfer-Encoding: base64"); n != 2 {
		t.Errorf("części base64 = %d, oczekiwano 2", n)
	}

	if len(alternative) != 2 || alternative[0] != "text/plain" || alternative[1] != "text/html" {
		t.Errorf("kolejność części alternative = %v, oczekiwano [text/plain text/html]", alternative)
	}
	wantRelated := []string{"multipart/alternative", "image/png", "image/png"}
	if len(typyRelated) != 3 {
		t.Fatalf("części related = %v, oczekiwano %v", typyRelated, wantRelated)
	}
	for i, want := range wantRelated {
		if typyRelated[i] != want {
			t.Errorf("część %d = %s, oczekiwano %s", i, typyRelated[i], want)
		}
	}
	wantIDs := []string{"<danaco-lockup>", "<danaco-lockup-dark>"}
	for i, want := range wantIDs {
		if i >= len(contentIDs) || contentIDs[i] != want {
			t.Errorf("Content-ID %d = %v, oczekiwano %s", i, contentIDs, want)
		}
	}
}

func TestBuildWszystkieRodzaje(t *testing.T) {
	set := wczytaj(t)
	for _, kind := range Kinds() {
		msg, err := set.Build(kind, koperta(), wartosci(t, set, kind), znaki())
		if err != nil {
			t.Fatalf("%s: %v", kind, err)
		}
		if strings.Contains(string(msg.Data), "{{") {
			t.Errorf("%s: w wiadomości pozostała zmienna", kind)
		}
		if len(msg.Data) > MaxMessageSize {
			t.Errorf("%s: %d B przekracza próg", kind, len(msg.Data))
		}
		if msg.MessageID == "" {
			t.Errorf("%s: pusty identyfikator wiadomości", kind)
		}
	}
}

func TestBuildBrakWartosci(t *testing.T) {
	set := wczytaj(t)
	c := wartosci(t, set, KindLoginCode)
	delete(c.Values, "code")
	_, err := set.Build(KindLoginCode, koperta(), c, znaki())
	if err == nil {
		t.Fatal("oczekiwano błędu przy braku wartości")
	}
	if !strings.Contains(err.Error(), "code") {
		t.Errorf("komunikat nie wskazuje brakującej zmiennej: %v", err)
	}
}

func TestBuildBrakZnaku(t *testing.T) {
	set := wczytaj(t)
	l := znaki()
	l.Dark = nil
	if _, err := set.Build(KindLoginCode, koperta(), wartosci(t, set, KindLoginCode), l); err == nil {
		t.Fatal("oczekiwano błędu przy braku wariantu znaku")
	}
}

func TestBuildZmiennaZaufanaINiezaufana(t *testing.T) {
	set := wczytaj(t)
	c := wartosci(t, set, KindNoMailbox)
	c.Values["original_subject"] = "z rdzenia"
	if _, err := set.Build(KindNoMailbox, koperta(), c, znaki()); err == nil {
		t.Fatal("oczekiwano błędu przy zmiennej podanej dwa razy")
	}
}

// TestBuildTematNiezaufanyBezWstrzykniecia sprawdza, że temat nadesłany
// z zewnątrz nie dopisuje własnych nagłówków ani znaczników.
func TestBuildTematNiezaufanyBezWstrzykniecia(t *testing.T) {
	set := wczytaj(t)
	c := wartosci(t, set, KindNoMailbox)
	c.Untrusted["original_subject"] = "Zwykły\r\nBcc: napastnik@przyklad.pl\r\n<script>alert(1)</script>"

	msg, err := set.Build(KindNoMailbox, koperta(), c, znaki())
	if err != nil {
		t.Fatalf("Build: %v", err)
	}
	parsed, err := netmail.ReadMessage(strings.NewReader(string(msg.Data)))
	if err != nil {
		t.Fatalf("wiadomość nie parsuje się: %v", err)
	}
	if parsed.Header.Get("Bcc") != "" {
		t.Error("wstrzyknięto nagłówek Bcc")
	}

	html := czescHTML(t, parsed)
	if strings.Contains(html, "<script>") {
		t.Error("znacznik script trafił do części HTML bez ucieczki")
	}
	if !strings.Contains(html, "&lt;script&gt;") {
		t.Error("brak ucieczki znaczników w części HTML")
	}
}

// TestBuildCzescTekstowaBezEncji pilnuje rozdziału reguł sanityzacji:
// ucieczka HTML nie może wyciec do części tekstowej.
func TestBuildCzescTekstowaBezEncji(t *testing.T) {
	set := wczytaj(t)
	c := wartosci(t, set, KindNoMailbox)
	c.Untrusted["original_subject"] = `Umowa "najmu" & aneks`

	msg, err := set.Build(KindNoMailbox, koperta(), c, znaki())
	if err != nil {
		t.Fatalf("Build: %v", err)
	}
	parsed, err := netmail.ReadMessage(strings.NewReader(string(msg.Data)))
	if err != nil {
		t.Fatalf("wiadomość nie parsuje się: %v", err)
	}
	text := czescTekstowa(t, parsed)
	if strings.Contains(text, "&amp;") || strings.Contains(text, "&quot;") {
		t.Errorf("encje HTML w części tekstowej: %q", text)
	}
	if !strings.Contains(text, `Umowa "najmu" & aneks`) {
		t.Error("temat nie trafił do części tekstowej w postaci zwykłej")
	}
}

func TestLogFieldsBezTresci(t *testing.T) {
	set := wczytaj(t)
	c := wartosci(t, set, KindLoginCode)
	c.Values["code"] = "418402"
	msg, err := set.Build(KindLoginCode, koperta(), c, znaki())
	if err != nil {
		t.Fatalf("Build: %v", err)
	}
	for klucz, wartosc := range msg.LogFields() {
		if strings.Contains(wartosc, "418402") {
			t.Errorf("kod trafił do pola dziennika %s", klucz)
		}
	}
}

// czescTekstowa i czescHTML wyciągają odkodowaną treść wskazanej części.
func czescTekstowa(t *testing.T, m *netmail.Message) string {
	t.Helper()
	return czesc(t, m, "text/plain")
}

func czescHTML(t *testing.T, m *netmail.Message) string {
	t.Helper()
	return czesc(t, m, "text/html")
}

func czesc(t *testing.T, m *netmail.Message, szukany string) string {
	t.Helper()
	_, params, err := mime.ParseMediaType(m.Header.Get("Content-Type"))
	if err != nil {
		t.Fatalf("Content-Type: %v", err)
	}
	related := multipart.NewReader(m.Body, params["boundary"])
	for {
		part, err := related.NextPart()
		if err == io.EOF {
			break
		}
		if err != nil {
			t.Fatalf("część related: %v", err)
		}
		mt, p, err := mime.ParseMediaType(part.Header.Get("Content-Type"))
		if err != nil || mt != "multipart/alternative" {
			continue
		}
		inner := multipart.NewReader(part, p["boundary"])
		for {
			sub, err := inner.NextPart()
			if err == io.EOF {
				break
			}
			if err != nil {
				t.Fatalf("część alternative: %v", err)
			}
			smt, _, _ := mime.ParseMediaType(sub.Header.Get("Content-Type"))
			if smt != szukany {
				continue
			}
			// Part odkodowuje quoted-printable samoczynnie.
			body, err := io.ReadAll(sub)
			if err != nil {
				t.Fatalf("odczyt części %s: %v", szukany, err)
			}
			return string(body)
		}
	}
	t.Fatalf("nie znaleziono części %s", szukany)
	return ""
}
