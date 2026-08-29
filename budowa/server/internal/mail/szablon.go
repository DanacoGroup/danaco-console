// Pakiet mail składa wiadomości poczty transakcyjnej Danaco Console
// z szablonów katalogu 06-poczta-transakcyjna.
//
// Pakiet nie wysyła wiadomości i nie prowadzi dziennika — zwraca gotowe
// bajty. Wysyłkę i zapis zdarzenia prowadzi warstwa wywołująca; kody
// uwierzytelniające nie mogą trafić do dziennika (patrz Message.LogFields).
package mail

import (
	"fmt"
	"io/fs"
	"path"
	"regexp"
	"sort"
	"strings"
)

// Kind to rodzaj listu; wartość jest rdzeniem nazw obu plików szablonu.
type Kind string

const (
	KindAccountActivation  Kind = "list-01-aktywacja-konta"
	KindLoginCode          Kind = "list-02-uwierzytelnienie-logowania"
	KindNoMailbox          Kind = "list-03-autoresponder-brak-skrzynki"
	KindPasswordReset      Kind = "list-04-reset-hasla"
	KindAddressChangeCode  Kind = "list-05-zmiana-adresu-potwierdzenie"
	KindAddressChangeAlert Kind = "list-06-zmiana-adresu-powiadomienie"
	KindRunFinished        Kind = "list-07-przebieg-zakonczony"
)

// subjects to wzorce tematów. Kodu w temacie nie ma: temat zapisuje każdy
// element trasy — filtr, brama, kopia kolejki — a to rozszerza powierzchnię
// wycieku poza serwer, którym władamy (runbook-wdrozenia.md, rozdz. 5.3).
// Rozstrzygnięcie Właściciela z 29.08.2026; kod stoi wyłącznie w treści listu.
var subjects = map[Kind]string{
	KindAccountActivation:  "Danaco Console — kod aktywacji konta",
	KindLoginCode:          "Danaco Console — kod logowania",
	KindNoMailbox:          "Adres noreply@danaco-group.pl nie przyjmuje korespondencji",
	KindPasswordReset:      "Danaco Console — kod resetu hasła",
	KindAddressChangeCode:  "Danaco Console — kod potwierdzenia nowego adresu",
	KindAddressChangeAlert: "Danaco Console — zamówiono zmianę adresu konta",
	KindRunFinished:        "Danaco Console — przebieg {{run_name}}: {{run_status}}",
}

// autoReplied rozstrzyga o wartości nagłówka Auto-Submitted (RFC 3834).
var autoReplied = map[Kind]bool{KindNoMailbox: true}

// Kinds zwraca wszystkie rodzaje listów w kolejności numerów.
func Kinds() []Kind {
	out := make([]Kind, 0, len(subjects))
	for k := range subjects {
		out = append(out, k)
	}
	sort.Slice(out, func(i, j int) bool { return out[i] < out[j] })
	return out
}

// Template to para plików jednego listu wraz z wzorcem tematu.
type Template struct {
	Kind    Kind
	Subject string
	HTML    string
	Text    string
	Vars    []string
}

// Set to komplet szablonów wczytany z katalogu.
type Set struct {
	templates map[Kind]*Template
}

var (
	rePlaceholder = regexp.MustCompile(`\{\{(\w+)\}\}`)
	// Komentarz zwykły; komentarz warunkowy Outlooka (<!--[if …) zachowujemy.
	reComment = regexp.MustCompile(`(?s)<!--(?:[^\[].*?|)-->`)
)

// Load wczytuje szablony z systemu plików. Oczekuje układu katalogu
// 06-poczta-transakcyjna: html/<kind>.html oraz text/<kind>.txt.
//
// W rdzeniu wywołanie wygląda tak:
//
//	//go:embed html/*.html text/*.txt
//	var templateFiles embed.FS
//	set, err := mail.Load(templateFiles)
func Load(fsys fs.FS) (*Set, error) {
	set := &Set{templates: make(map[Kind]*Template, len(subjects))}
	for kind, subject := range subjects {
		htmlBytes, err := fs.ReadFile(fsys, path.Join("html", string(kind)+".html"))
		if err != nil {
			return nil, fmt.Errorf("szablon HTML %s: %w", kind, err)
		}
		textBytes, err := fs.ReadFile(fsys, path.Join("text", string(kind)+".txt"))
		if err != nil {
			return nil, fmt.Errorf("szablon tekstowy %s: %w", kind, err)
		}
		// Komentarze szablonu są dokumentacją wewnętrzną i nie mogą trafić
		// do skrzynki Operatora. Usuwamy je przy wczytaniu, nie przy składaniu.
		htmlClean := strings.TrimSpace(reComment.ReplaceAllString(string(htmlBytes), ""))
		tpl := &Template{
			Kind:    kind,
			Subject: subject,
			HTML:    htmlClean,
			Text:    string(textBytes),
		}
		tpl.Vars = placeholders(subject + htmlClean + tpl.Text)
		set.templates[kind] = tpl
	}
	return set, nil
}

// Template zwraca szablon danego rodzaju.
func (s *Set) Template(kind Kind) (*Template, error) {
	tpl, ok := s.templates[kind]
	if !ok {
		return nil, fmt.Errorf("nieznany rodzaj listu: %s", kind)
	}
	return tpl, nil
}

// placeholders zwraca posortowaną listę nazw zmiennych bez powtórzeń.
func placeholders(src string) []string {
	seen := make(map[string]struct{})
	for _, m := range rePlaceholder.FindAllStringSubmatch(src, -1) {
		seen[m[1]] = struct{}{}
	}
	out := make([]string, 0, len(seen))
	for name := range seen {
		out = append(out, name)
	}
	sort.Strings(out)
	return out
}

// substitute podstawia wartości i zwraca błąd przy zmiennej bez wartości —
// inaczej list wyszedłby z ciągiem {{code}} zamiast kodu.
func substitute(src string, values map[string]string) (string, error) {
	var missing []string
	out := rePlaceholder.ReplaceAllStringFunc(src, func(m string) string {
		name := rePlaceholder.FindStringSubmatch(m)[1]
		v, ok := values[name]
		if !ok {
			missing = append(missing, name)
			return m
		}
		return v
	})
	if len(missing) > 0 {
		sort.Strings(missing)
		return "", fmt.Errorf("zmienne bez wartości: %s", strings.Join(missing, ", "))
	}
	return out, nil
}
