// Odpowiedzialność pliku: list jako dokument — rozbiór przysłanego MIME na
// treść i załączniki oraz złożenie MIME wychodzącego. Jedno i drugie stoi tu
// razem, bo to jest ta sama wiedza czytana w dwie strony.
package poczta

import (
	"bytes"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/emersion/go-imap/v2"
	"github.com/emersion/go-imap/v2/imapclient"
	// Rejestracja dekoderów stron kodowych, skutkiem ubocznym importu pustego.
	_ "github.com/emersion/go-message/charset"
	"github.com/emersion/go-message/mail"
)

// zlozNaglowek przekłada odpowiedź FETCH na nagłówek pakietu, samą kopertę
// bez treści ani załączników.
func zlozNaglowek(folder string, w *imapclient.FetchMessageBuffer) Naglowek {
	n := Naglowek{
		Identyfikator:  zlozIdentyfikator(folder, w.UID),
		Folder:         folder,
		Chwila:         chwilaListu(w.Envelope, w.InternalDate),
		Nieprzeczytana: czyNieprzeczytana(w.Flags),
	}
	if w.Envelope != nil {
		n.Od = pierwszyAdres(w.Envelope.From)
		n.Do = adresy(w.Envelope.To)
		n.Kopia = adresy(w.Envelope.Cc)
		n.Temat = w.Envelope.Subject
		// Wątek wiąże identyfikator listu, na który ten odpowiada, albo własny.
		n.Watek = pierwszyNiepusty(w.Envelope.InReplyTo...)
		if n.Watek == "" {
			n.Watek = w.Envelope.MessageID
		}
	}
	n.Zapowiedz = zapowiedz(w.FindBodySection(sekcjaZapowiedzi()))
	if w.BodyStructure != nil {
		n.NazwyZalacznikow = nazwyZBudowy(w.BodyStructure)
	}
	return n
}

// zapowiedz oczyszcza początek treści do jednej linii. Surowe bajty części
// niosą łamania wierszy i bywają urwane w połowie znaku wielobajtowego —
// jedno i drugie wyglądałoby w oknie na uszkodzoną treść.
func zapowiedz(surowe []byte) string {
	if len(surowe) == 0 {
		return ""
	}
	// Odkodowanie transportowe zostaje nietknięte świadomie.
	tekst := strings.TrimSpace(strings.ToValidUTF8(string(surowe), ""))
	tekst = strings.Join(strings.Fields(tekst), " ")
	return tekst
}

// nazwyZBudowy wyciąga nazwy załączników z BODYSTRUCTURE — bez ściągania
// choćby jednego bajtu treści. To jest cała różnica między wykazem a odczytem.
func nazwyZBudowy(budowa imap.BodyStructure) []string {
	var nazwy []string
	budowa.Walk(func(_ []int, czesc imap.BodyStructure) bool {
		pojedyncza, jest := czesc.(*imap.BodyStructureSinglePart)
		if !jest {
			return true
		}
		if nazwa := pojedyncza.Filename(); nazwa != "" {
			nazwy = append(nazwy, nazwa)
		}
		return true
	})
	return nazwy
}

// rozbierzList rozkłada surowy dokument MIME na treść i załączniki. Dokument
// bez MIME też jest listem: `mail.CreateReader` odda go jako jedną część
// wpisaną, więc obsługuje go ta sama pętla bez przypadku szczególnego.
func rozbierzList(surowy []byte) (string, []Zalacznik, error) {
	czytnik, err := mail.CreateReader(bytes.NewReader(surowy))
	if err != nil {
		return "", nil, err
	}
	defer czytnik.Close()

	var (
		zwykla     strings.Builder
		zapasowa   strings.Builder
		zalaczniki []Zalacznik
	)
	for {
		czesc, err := czytnik.NextPart()
		if err == io.EOF {
			break
		}
		if err != nil {
			return "", nil, err
		}
		switch naglowek := czesc.Header.(type) {
		case *mail.InlineHeader:
			bajty, err := io.ReadAll(czesc.Body)
			if err != nil {
				return "", nil, err
			}
			typ, _, _ := naglowek.ContentType()
			if typ == "text/plain" {
				zwykla.Write(bajty)
			} else {
				zapasowa.Write(bajty)
			}
		case *mail.AttachmentHeader:
			bajty, err := io.ReadAll(czesc.Body)
			if err != nil {
				return "", nil, err
			}
			nazwa, _ := naglowek.Filename()
			if strings.TrimSpace(nazwa) == "" {
				// Załącznik bez nazwy dostaje etykietę zastępczą.
				nazwa = fmt.Sprintf("zalacznik-%d", len(zalaczniki)+1)
			}
			typ, _, _ := naglowek.ContentType()
			zalaczniki = append(zalaczniki, Zalacznik{Nazwa: nazwa, TypTresci: typ, Bajty: bajty})
		}
	}
	if tresc := strings.TrimSpace(zwykla.String()); tresc != "" {
		return zwykla.String(), zalaczniki, nil
	}
	return zapasowa.String(), zalaczniki, nil
}

// zlozWychodzacy buduje dokument MIME listu do nadania, tą samą wiedzą, którą
// `rozbierzList` czyta go z powrotem.
func zlozWychodzacy(w Wychodzacy) ([]byte, error) {
	if strings.TrimSpace(w.Od) == "" {
		return nil, fmt.Errorf("list bez nadawcy — skrzynka nie ma zapisanego adresu")
	}
	naglowek := mail.Header{}
	naglowek.SetAddressList("From", []*mail.Address{{Name: w.NazwaWyswietlana, Address: w.Od}})
	if adresaci := adresyMail(w.Do); len(adresaci) > 0 {
		naglowek.SetAddressList("To", adresaci)
	}
	if kopia := adresyMail(w.Kopia); len(kopia) > 0 {
		naglowek.SetAddressList("Cc", kopia)
	}
	naglowek.SetSubject(w.Temat)
	naglowek.SetDate(time.Now())
	if err := naglowek.GenerateMessageID(); err != nil {
		return nil, fmt.Errorf("nie można nadać listowi identyfikatora: %w", err)
	}
	if odpowiedz := bezNawiasowKatowych(w.WOdpowiedziNa); odpowiedz != "" {
		// Dwa nagłówki, nie jeden: `In-Reply-To` wiąże list bezpośredni,
		// `References` cały wątek.
		naglowek.SetMsgIDList("In-Reply-To", []string{odpowiedz})
		naglowek.SetMsgIDList("References", []string{odpowiedz})
	}

	var bufor bytes.Buffer
	pisarz, err := mail.CreateWriter(&bufor, naglowek)
	if err != nil {
		return nil, fmt.Errorf("nie można złożyć listu: %w", err)
	}
	wpisana, err := pisarz.CreateSingleInline(inlineTekst())
	if err != nil {
		return nil, fmt.Errorf("nie można złożyć treści listu: %w", err)
	}
	if _, err := io.WriteString(wpisana, w.Tresc); err != nil {
		return nil, fmt.Errorf("nie można zapisać treści listu: %w", err)
	}
	if err := wpisana.Close(); err != nil {
		return nil, fmt.Errorf("nie można domknąć treści listu: %w", err)
	}

	for _, z := range w.Zalaczniki {
		naglowekZalacznika := mail.AttachmentHeader{}
		naglowekZalacznika.SetContentType(typTresci(z.TypTresci), nil)
		naglowekZalacznika.SetFilename(z.Nazwa)
		strumien, err := pisarz.CreateAttachment(naglowekZalacznika)
		if err != nil {
			return nil, fmt.Errorf("nie można dołączyć %q: %w", z.Nazwa, err)
		}
		if _, err := strumien.Write(z.Bajty); err != nil {
			return nil, fmt.Errorf("nie można zapisać załącznika %q: %w", z.Nazwa, err)
		}
		if err := strumien.Close(); err != nil {
			return nil, fmt.Errorf("nie można domknąć załącznika %q: %w", z.Nazwa, err)
		}
	}
	if err := pisarz.Close(); err != nil {
		return nil, fmt.Errorf("nie można domknąć listu: %w", err)
	}
	return bufor.Bytes(), nil
}

// inlineTekst opisuje część treściową listu — zwykły tekst w UTF-8, jedyny
// typ, jaki `zlozWychodzacy` składa.
func inlineTekst() mail.InlineHeader {
	h := mail.InlineHeader{}
	h.SetContentType("text/plain", map[string]string{"charset": "utf-8"})
	return h
}

// typTresci bierze typ wskazany, a przy jego braku — najogólniejszy możliwy,
// `application/octet-stream`, zamiast zgadywać.
func typTresci(wskazany string) string {
	if s := strings.TrimSpace(wskazany); s != "" {
		return s
	}
	return "application/octet-stream"
}

// adresyMail przekłada adresy tekstowe na kształt biblioteki, pomijając puste,
// do pól `To` i `Cc` nagłówka.
func adresyMail(lista []string) []*mail.Address {
	wynik := make([]*mail.Address, 0, len(lista))
	for _, a := range lista {
		if a = strings.TrimSpace(a); a != "" {
			wynik = append(wynik, &mail.Address{Address: a})
		}
	}
	return wynik
}

// bezNawiasowKatowych zdejmuje z identyfikatora listu nawiasy `<>`, jeśli je
// niesie: `SetMsgIDList` zakłada identyfikator goły i sam dokłada nawiasy,
// więc wartość już w nawiasach dałaby zapis podwójny, którego RFC 5322 nie
// zna.
func bezNawiasowKatowych(identyfikator string) string {
	identyfikator = strings.TrimSpace(identyfikator)
	identyfikator = strings.TrimPrefix(identyfikator, "<")
	return strings.TrimSuffix(identyfikator, ">")
}

// pierwszyNiepusty oddaje pierwszą niepustą wartość z wykazu tekstów — pustka
// znaczy brak wartości w ogóle, nie wpis pusty.
func pierwszyNiepusty(wartosci ...string) string {
	for _, w := range wartosci {
		if s := strings.TrimSpace(w); s != "" {
			return s
		}
	}
	return ""
}
