package mail

import (
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"mime"
	"mime/multipart"
	"mime/quotedprintable"
	netmail "net/mail"
	"net/textproto"
	"strings"
	"time"
)

// MaxMessageSize to próg, powyżej którego Gmail przycina treść i podstawia
// odnośnik „pokaż całą wiadomość". List transakcyjny nie może być przycięty.
const MaxMessageSize = 102 * 1024

// Content-ID obu wariantów znaku. Wartości muszą zgadzać się z tym,
// co szablony podstawiają pod {{logo_src}} i {{logo_src_dark}}.
const (
	logoContentID     = "danaco-lockup"
	logoDarkContentID = "danaco-lockup-dark"
)

// Logos to oba warianty znaku dołączane do każdej wiadomości. Który z nich
// zobaczy Operator, rozstrzyga jego klient pocztowy.
type Logos struct {
	Light     []byte
	LightName string
	Dark      []byte
	DarkName  string
}

// Envelope to dane koperty wiadomości.
type Envelope struct {
	From netmail.Address
	To   netmail.Address
	// Date pusta oznacza czas złożenia wiadomości.
	Date time.Time
	// MessageID bez nawiasów ostrych; pusty zostanie wytworzony.
	MessageID string
}

// Content to wartości podstawiane w szablonie.
type Content struct {
	// Values pochodzą z rdzenia i są podstawiane bez zmian.
	Values map[string]string
	// Untrusted pochodzą spoza rdzenia (temat wiadomości nadesłanej, nazwa
	// automatyki nadana przez Operatora) i przechodzą sanityzację.
	Untrusted map[string]string
}

// Message to złożona wiadomość gotowa do podania serwerowi poczty.
type Message struct {
	Kind      Kind
	To        string
	MessageID string
	Data      []byte
}

// LogFields zwraca pola, które wolno zapisać w dzienniku. Treść, temat
// i kody uwierzytelniające nie występują — celowo.
func (m *Message) LogFields() map[string]string {
	return map[string]string{
		"kind":       string(m.Kind),
		"to":         m.To,
		"message_id": m.MessageID,
		"size_bytes": fmt.Sprintf("%d", len(m.Data)),
	}
}

// Build składa wiadomość: multipart/related z zagnieżdżonym
// multipart/alternative (tekst przed HTML) oraz dwoma obrazami cid.
func (s *Set) Build(kind Kind, env Envelope, content Content, logos Logos) (*Message, error) {
	tpl, err := s.Template(kind)
	if err != nil {
		return nil, err
	}
	if err := validateLogos(logos); err != nil {
		return nil, err
	}

	// Trzy komplety wartości: tekst bez ucieczki, HTML z ucieczką,
	// temat bez ucieczki (nagłówek kodujemy osobno wg RFC 2047).
	textValues, htmlValues, subjectValues, err := prepareValues(content)
	if err != nil {
		return nil, err
	}
	htmlValues["logo_src"] = "cid:" + logoContentID
	htmlValues["logo_src_dark"] = "cid:" + logoDarkContentID

	subject, err := substitute(tpl.Subject, subjectValues)
	if err != nil {
		return nil, fmt.Errorf("temat listu %s: %w", kind, err)
	}
	textBody, err := substitute(tpl.Text, textValues)
	if err != nil {
		return nil, fmt.Errorf("część tekstowa listu %s: %w", kind, err)
	}
	htmlBody, err := substitute(tpl.HTML, htmlValues)
	if err != nil {
		return nil, fmt.Errorf("część HTML listu %s: %w", kind, err)
	}

	date := env.Date
	if date.IsZero() {
		date = time.Now()
	}
	messageID := env.MessageID
	if messageID == "" {
		messageID, err = newMessageID(env.From.Address)
		if err != nil {
			return nil, err
		}
	}

	body := &bytes.Buffer{}
	related := multipart.NewWriter(body)

	if err := writeAlternative(related, textBody, htmlBody); err != nil {
		return nil, err
	}
	if err := writeImage(related, logoContentID, logos.LightName, logos.Light); err != nil {
		return nil, err
	}
	if err := writeImage(related, logoDarkContentID, logos.DarkName, logos.Dark); err != nil {
		return nil, err
	}
	if err := related.Close(); err != nil {
		return nil, err
	}

	head := headers(env, subject, messageID, date, kind, related.Boundary())
	data := append(head, body.Bytes()...)

	if len(data) > MaxMessageSize {
		return nil, fmt.Errorf("wiadomość %s ma %d B — próg przycięcia to %d B",
			kind, len(data), MaxMessageSize)
	}
	return &Message{Kind: kind, To: env.To.Address, MessageID: messageID, Data: data}, nil
}

// prepareValues rozdziela wartości na trzy komplety i pilnuje, żeby ta sama
// nazwa nie występowała jednocześnie jako zaufana i niezaufana.
func prepareValues(content Content) (text, html, subject map[string]string, err error) {
	text = make(map[string]string, len(content.Values)+len(content.Untrusted))
	html = make(map[string]string, len(content.Values)+len(content.Untrusted)+2)
	subject = make(map[string]string, len(content.Values)+len(content.Untrusted))

	for k, v := range content.Values {
		text[k], html[k], subject[k] = v, v, v
	}
	for k, v := range content.Untrusted {
		if _, dup := content.Values[k]; dup {
			return nil, nil, nil, fmt.Errorf("zmienna %s podana jako zaufana i niezaufana", k)
		}
		clean := cleanUntrusted(v)
		text[k], subject[k] = clean, clean
		html[k] = escapeForHTML(clean)
	}
	return text, html, subject, nil
}

func validateLogos(l Logos) error {
	switch {
	case len(l.Light) == 0 || len(l.Dark) == 0:
		return fmt.Errorf("brak jednego z wariantów znaku")
	case l.LightName == "" || l.DarkName == "":
		return fmt.Errorf("brak nazwy pliku znaku")
	}
	return nil
}

// writeAlternative dopisuje część multipart/alternative. Kolejność jest
// znacząca: klient wybiera ostatnią część, którą potrafi wyświetlić,
// więc tekst idzie przed HTML.
func writeAlternative(related *multipart.Writer, textBody, htmlBody string) error {
	inner := &bytes.Buffer{}
	alternative := multipart.NewWriter(inner)

	if err := writeQuotedPrintable(alternative, "text/plain", textBody); err != nil {
		return err
	}
	if err := writeQuotedPrintable(alternative, "text/html", htmlBody); err != nil {
		return err
	}
	if err := alternative.Close(); err != nil {
		return err
	}

	head := textproto.MIMEHeader{
		"Content-Type": {"multipart/alternative; boundary=" + alternative.Boundary()},
	}
	part, err := related.CreatePart(head)
	if err != nil {
		return err
	}
	_, err = part.Write(inner.Bytes())
	return err
}

// writeQuotedPrintable koduje treść. Quoted-printable jest konieczne:
// polskie znaki diakrytyczne to bajty spoza ASCII, a wiersze HTML
// przekraczają limit 998 oktetów na wiersz.
func writeQuotedPrintable(w *multipart.Writer, contentType, body string) error {
	head := textproto.MIMEHeader{
		"Content-Type":              {contentType + "; charset=utf-8"},
		"Content-Transfer-Encoding": {"quoted-printable"},
	}
	part, err := w.CreatePart(head)
	if err != nil {
		return err
	}
	enc := quotedprintable.NewWriter(part)
	if _, err := enc.Write([]byte(toCRLF(body))); err != nil {
		return err
	}
	return enc.Close()
}

func writeImage(w *multipart.Writer, contentID, filename string, data []byte) error {
	head := textproto.MIMEHeader{
		"Content-Type":              {"image/png"},
		"Content-Transfer-Encoding": {"base64"},
		"Content-ID":                {"<" + contentID + ">"},
		"Content-Disposition":       {fmt.Sprintf("inline; filename=%q", filename)},
	}
	part, err := w.CreatePart(head)
	if err != nil {
		return err
	}
	return writeBase64Lines(part, data)
}

// headers składa nagłówki wiadomości. Auto-Submitted i
// X-Auto-Response-Suppress powstrzymują automatyczne odpowiedzi po drugiej
// stronie — bez nich autoresponder odbiorcy odpisuje na kod logowania.
func headers(env Envelope, subject, messageID string, date time.Time, kind Kind, boundary string) []byte {
	autoSubmitted := "auto-generated"
	if autoReplied[kind] {
		autoSubmitted = "auto-replied"
	}
	var b strings.Builder
	writeHeader(&b, "From", env.From.String())
	writeHeader(&b, "To", env.To.String())
	writeHeader(&b, "Subject", mime.QEncoding.Encode("utf-8", subject))
	writeHeader(&b, "Date", date.Format(time.RFC1123Z))
	writeHeader(&b, "Message-ID", "<"+messageID+">")
	writeHeader(&b, "MIME-Version", "1.0")
	writeHeader(&b, "Auto-Submitted", autoSubmitted)
	writeHeader(&b, "X-Auto-Response-Suppress", "All")
	writeHeader(&b, "Content-Type",
		`multipart/related; type="multipart/alternative"; boundary=`+boundary)
	b.WriteString("\r\n")
	return []byte(b.String())
}

func writeHeader(b *strings.Builder, name, value string) {
	b.WriteString(name)
	b.WriteString(": ")
	b.WriteString(value)
	b.WriteString("\r\n")
}

// newMessageID wytwarza identyfikator z 16 bajtów losowych i domeny nadawcy.
func newMessageID(from string) (string, error) {
	buf := make([]byte, 16)
	if _, err := rand.Read(buf); err != nil {
		return "", fmt.Errorf("identyfikator wiadomości: %w", err)
	}
	domain := from
	if at := strings.LastIndex(from, "@"); at >= 0 {
		domain = from[at+1:]
	}
	return hex.EncodeToString(buf) + "@" + domain, nil
}

// toCRLF ujednolica złamania wiersza. Szablony są zapisane z LF, a poczta
// wymaga CRLF.
func toCRLF(s string) string {
	return strings.ReplaceAll(strings.ReplaceAll(s, "\r\n", "\n"), "\n", "\r\n")
}
