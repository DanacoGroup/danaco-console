// Odpowiedzialność pliku: szablony moderacji i wydanie transkryptu —
// `roundtable.moderation.template.save`, `roundtable.moderation.template.list`
// i `roundtable.transcript.export`.
//
// ── Cztery formaty, trzy drogi ───────────────────────────────────────────────
// Markdown i JSON składa rdzeń wprost, bez niczego z zewnątrz. PDF powstaje
// biblioteką wkompilowaną w binarium (`pdfcpu`) — dokument i kryptografia są
// w tym produkcie wyjątkiem bezwzględnym od wołania programów serwerowych.
// DOCX powstaje Pandokiem, bo formatu biurowego nie da się złożyć bibliotecznie
// w rdzeniu, a Pandoc jest programem serwerowym zadeklarowanym w sondzie
// zależności (`zaleznosci_zewnetrzne.go`) i używanym już przez Studio,
// Translate i Bibliotekę.
package core

import (
	"bytes"
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/pdfcpu/pdfcpu/pkg/api"

	"danacoconsole/server/internal/dane"
	"danacoconsole/server/internal/konfig"
	"danacoconsole/server/internal/session"
	"danacoconsole/server/internal/zewnetrzne"
	"danacoconsole/shared"
)

const (
	przedrostekSzablonu = "szablon-"

	// granicaZamianyFormatu — Pandoc zamienia transkrypt, a nie przetwarza
	// materiału filmowego; przekroczenie tej granicy znaczy plik uszkodzony albo
	// proces, który utknął.
	granicaZamianyFormatu = 3 * time.Minute

	// znakowWWierszuPdf — ile znaków mieści się w wierszu dokumentu przy foncie
	// o rozmiarze 11 punktów i marginesach domyślnych.
	znakowWWierszuPdf = 96
	// wierszyNaStroniePdf — ile wierszy mieści strona.
	wierszyNaStroniePdf = 46
)

// ZapiszSzablon utrwala format, liczbę tur i kolejność głosu jako szablon.
func (a *adapterDebaty) ZapiszSzablon(ctx context.Context,
	z shared.RoundtableModerationTemplateSaveRequest) (shared.RoundtableModerationTemplateSaveResponse, error) {

	okno := strings.TrimSpace(z.WindowId)
	nazwa := strings.TrimSpace(z.Name)
	if okno == "" {
		return shared.RoundtableModerationTemplateSaveResponse{},
			bladWskazaniaDebaty("komenda bez wskazania okna debaty")
	}
	if nazwa == "" {
		return shared.RoundtableModerationTemplateSaveResponse{},
			bladWskazaniaDebaty("szablon moderacji bez nazwy")
	}

	// Format i granice bierze się z tury ostatniej — to w niej debata biegła
	// tak, jak Operator ją ustawił. Debata bez ani jednej tury nie ma czego
	// zapisać jako szablon.
	tury, err := a.repozytorium.Tury(ctx, okno, 1)
	if err != nil {
		return shared.RoundtableModerationTemplateSaveResponse{}, bladDebaty(err)
	}
	if len(tury) == 0 {
		return shared.RoundtableModerationTemplateSaveResponse{}, odmowaSzablonuBezTury(okno)
	}
	uczestnicy, err := a.repozytorium.Uczestnicy(ctx, okno)
	if err != nil {
		return shared.RoundtableModerationTemplateSaveResponse{}, bladDebaty(err)
	}

	szablon, err := a.repozytorium.ZapiszSzablonDebaty(ctx, dane.SzablonModeracjiDebaty{
		Kod: nowyIdentyfikator(przedrostekSzablonu), Nazwa: nazwa,
		Format: tury[0].Format, GranicaTur: tury[0].GranicaTur,
		GranicaCzasuMs: tury[0].GranicaCzasuMs, KolejnoscGlosu: kodyUczestnikow(uczestnicy),
	})
	if err != nil {
		return shared.RoundtableModerationTemplateSaveResponse{}, bladDebaty(err)
	}
	return shared.RoundtableModerationTemplateSaveResponse{
		Template: szablonKontraktu(szablon),
	}, nil
}

// Szablony oddaje zapisane szablony moderacji.
func (a *adapterDebaty) Szablony(ctx context.Context,
	z shared.RoundtableModerationTemplateListRequest) (shared.RoundtableModerationTemplateListResponse, error) {

	szablony, err := a.repozytorium.SzablonyDebaty(ctx, wartoscTekstu(z.Query))
	if err != nil {
		return shared.RoundtableModerationTemplateListResponse{}, bladDebaty(err)
	}
	wykaz := make([]shared.RoundtableModerationTemplate, 0, len(szablony))
	for _, szablon := range szablony {
		wykaz = append(wykaz, szablonKontraktu(szablon))
	}
	return shared.RoundtableModerationTemplateListResponse{Templates: wykaz}, nil
}

// WydajTranskrypt wydaje pełny zapis debaty jako artefakt.
func (a *adapterDebaty) WydajTranskrypt(ctx context.Context,
	z shared.RoundtableTranscriptExportRequest) (shared.RoundtableTranscriptExportResponse, error) {

	okno := strings.TrimSpace(z.WindowId)
	if okno == "" {
		return shared.RoundtableTranscriptExportResponse{},
			bladWskazaniaDebaty("komenda bez wskazania okna debaty")
	}
	format := strings.TrimSpace(string(z.Format))
	turaKod := strings.TrimSpace(wartoscTekstu(z.TurnId))
	zUstaleniami := z.IncludeAnalysis != nil && *z.IncludeAnalysis

	var bajty []byte
	var err error
	switch format {
	case shared.RoundtableTranscriptFormatMarkdown:
		var tekst string
		if tekst, err = a.transkryptMarkdown(ctx, okno, turaKod, zUstaleniami); err == nil {
			bajty = []byte(tekst)
		}
	case shared.RoundtableTranscriptFormatJson:
		bajty, err = a.transkryptJson(ctx, okno, turaKod, zUstaleniami)
	case shared.RoundtableTranscriptFormatPdf:
		var tekst string
		if tekst, err = a.transkryptMarkdown(ctx, okno, turaKod, zUstaleniami); err == nil {
			bajty, err = dokumentPdfZTekstu(tekst)
		}
	case shared.RoundtableTranscriptFormatDocx:
		var tekst string
		if tekst, err = a.transkryptMarkdown(ctx, okno, turaKod, zUstaleniami); err == nil {
			bajty, err = a.zamienNaDocx(ctx, tekst)
		}
	default:
		return shared.RoundtableTranscriptExportResponse{},
			bladWskazaniaDebaty("format transkryptu " + format + " nie jest formatem znanym kontraktowi")
	}
	if err != nil {
		return shared.RoundtableTranscriptExportResponse{}, err
	}

	artefakt, err := a.wydajArtefaktDebaty(ctx, okno, rodzajArtefaktuTranskryptu, format, bajty, 0)
	if err != nil {
		return shared.RoundtableTranscriptExportResponse{}, err
	}
	return shared.RoundtableTranscriptExportResponse{
		ArtifactId: artefakt.Kod, Uri: odwolanieArtefaktu(artefakt),
	}, nil
}

// transkryptMarkdown składa zapis debaty: tury, pytania, wypowiedzi z podpisami
// i — na żądanie — ustalenia analizy wraz z grafem argumentów.
func (a *adapterDebaty) transkryptMarkdown(ctx context.Context, okno, turaKod string,
	zUstaleniami bool) (string, error) {

	uczestnicy, err := a.repozytorium.Uczestnicy(ctx, okno)
	if err != nil {
		return "", bladDebaty(err)
	}
	podpisy := podpisyUczestnikow(uczestnicy)

	tury, err := a.turyOknaOdPierwszej(ctx, okno)
	if err != nil {
		return "", err
	}
	czesci := make([]string, 0, len(tury)+3)
	czesci = append(czesci, "# Zapis debaty okna "+okno)
	for _, tura := range tury {
		if turaKod != "" && tura.Kod != turaKod {
			continue
		}
		wypowiedzi, err := a.repozytorium.Wypowiedzi(ctx, tura.Kod)
		if err != nil {
			return "", bladDebaty(err)
		}
		czesci = append(czesci, zapisTury(tura, wypowiedzi, podpisy))
	}
	if len(czesci) == 1 {
		return "", odmowaTranskryptuBezZapisu(okno)
	}

	if zUstaleniami {
		ustalenia, err := a.repozytorium.UstaleniaDebaty(ctx, okno, "", turaKod)
		if err != nil {
			return "", bladDebaty(err)
		}
		if len(ustalenia) > 0 {
			wiersze := make([]string, 0, len(ustalenia)+1)
			wiersze = append(wiersze, "## Ustalenia analizy")
			for _, ustalenie := range ustalenia {
				wiersze = append(wiersze, "- ["+ustalenie.Rodzaj+"] "+ustalenie.Tresc)
			}
			czesci = append(czesci, strings.Join(wiersze, "\n"))
		}
		graf, err := a.zlozGraf(ctx, okno, turaKod, false)
		if err != nil {
			return "", err
		}
		if len(graf.Nodes) > 0 {
			wiersze := make([]string, 0, len(graf.Nodes)+1)
			wiersze = append(wiersze, "## Graf argumentów")
			for _, wezel := range graf.Nodes {
				wiersze = append(wiersze, "- "+string(wezel.SpeechAct)+": "+wezel.Text)
			}
			czesci = append(czesci, strings.Join(wiersze, "\n"))
		}
	}
	return strings.Join(czesci, "\n\n"), nil
}

// transkryptJson wydaje zapis debaty w postaci nadającej się do dalszego
// przetworzenia: migawka stanu wraz z ustaleniami, gdy zażądano.
func (a *adapterDebaty) transkryptJson(ctx context.Context, okno, turaKod string,
	zUstaleniami bool) ([]byte, error) {

	stan, err := a.StanDebaty(ctx, shared.RoundtableDebateGetRequest{
		WindowId: okno, TurnId: wskaznikTekstu(turaKod),
	})
	if err != nil {
		return nil, err
	}
	if len(stan.Snapshot.Turns) == 0 {
		return nil, odmowaTranskryptuBezZapisu(okno)
	}

	wydanie := struct {
		Snapshot shared.RoundtableDebateSnapshot    `json:"snapshot"`
		Findings []shared.RoundtableAnalysisFinding `json:"findings,omitempty"`
		Graph    *shared.RoundtableArgumentGraph    `json:"graph,omitempty"`
	}{Snapshot: stan.Snapshot}

	if zUstaleniami {
		ustalenia, err := a.repozytorium.UstaleniaDebaty(ctx, okno, "", turaKod)
		if err != nil {
			return nil, bladDebaty(err)
		}
		wydanie.Findings = ustaleniaKontraktu(ustalenia)
		graf, err := a.zlozGraf(ctx, okno, turaKod, false)
		if err != nil {
			return nil, err
		}
		wydanie.Graph = &graf
	}
	bajty, err := json.MarshalIndent(wydanie, "", "  ")
	if err != nil {
		return nil, bladDebaty(err)
	}
	return bajty, nil
}

// dokumentPdfZTekstu składa dokument z tekstu, biblioteką wkompilowaną w rdzeń.
//
// Rdzeń nie woła tu żadnego programu i nie ma prawa go wołać: dokument
// i kryptografia są w tym produkcie wyjątkiem bezwzględnym — robi je biblioteka
// Go albo nikt.
func dokumentPdfZTekstu(tekst string) ([]byte, error) {
	wiersze := zawinWiersze(tekst, znakowWWierszuPdf)
	if len(wiersze) == 0 {
		return nil, odmowaPustegoArtefaktu(rodzajArtefaktuTranskryptu)
	}

	var opis strings.Builder
	opis.WriteString(`{"pages":{`)
	numerStrony := 0
	for poczatek := 0; poczatek < len(wiersze); poczatek += wierszyNaStroniePdf {
		koniec := poczatek + wierszyNaStroniePdf
		if koniec > len(wiersze) {
			koniec = len(wiersze)
		}
		numerStrony++
		if numerStrony > 1 {
			opis.WriteString(",")
		}
		opis.WriteString(`"`)
		opis.WriteString(itoa(numerStrony))
		opis.WriteString(`":{"content":{"text":[{"value":`)
		wartosc, err := json.Marshal(strings.Join(wiersze[poczatek:koniec], "\n"))
		if err != nil {
			return nil, bladDebaty(err)
		}
		opis.Write(wartosc)
		opis.WriteString(`,"font":{"name":"Helvetica","size":11},"position":[0.08,0.94],` +
			`"anchor":"topleft"}]}}`)
	}
	opis.WriteString(`}}`)

	var dokument bytes.Buffer
	if err := api.Create(nil, strings.NewReader(opis.String()), &dokument, nastawyPdf()); err != nil {
		return nil, odmowaSkladaniaDokumentu(err)
	}
	return dokument.Bytes(), nil
}

// zawinWiersze łamie tekst na wiersze mieszczące się w szerokości strony.
// Łamanie idzie po słowach — łamanie w środku słowa dawałoby zapis, którego nie
// da się przeczytać ani przeszukać.
func zawinWiersze(tekst string, szerokosc int) []string {
	wiersze := make([]string, 0, 64)
	for _, akapit := range strings.Split(tekst, "\n") {
		if strings.TrimSpace(akapit) == "" {
			wiersze = append(wiersze, "")
			continue
		}
		biezacy := ""
		for _, slowo := range strings.Fields(akapit) {
			if biezacy == "" {
				biezacy = slowo
				continue
			}
			if len([]rune(biezacy))+1+len([]rune(slowo)) > szerokosc {
				wiersze = append(wiersze, biezacy)
				biezacy = slowo
				continue
			}
			biezacy += " " + slowo
		}
		if biezacy != "" {
			wiersze = append(wiersze, biezacy)
		}
	}
	return wiersze
}

// zamienNaDocx składa dokument biurowy Pandokiem.
//
// Materiał wchodzi plikiem, nie strumieniem: Pandoc rozpoznaje format wyjściowy
// po rozszerzeniu pliku docelowego, a zapis do strumienia wymagałby wskazania
// go osobno i tak samo tworzyłby plik pośredni.
func (a *adapterDebaty) zamienNaDocx(ctx context.Context, tekst string) ([]byte, error) {
	if a.uruchamiacz == nil {
		return nil, odmowaBrakuUruchamiacza()
	}
	katalog, err := os.MkdirTemp("", "roundtable-docx-")
	if err != nil {
		return nil, bladDebaty(err)
	}
	defer os.RemoveAll(katalog)

	zrodlo := filepath.Join(katalog, "transkrypt.md")
	docelowy := filepath.Join(katalog, "transkrypt.docx")
	if err := os.WriteFile(zrodlo, []byte(tekst), 0o600); err != nil {
		return nil, bladDebaty(err)
	}

	okno, zasady, obszar := a.zasiegNarzedziDebaty()
	if _, err := zewnetrzne.Wolaj(ctx, a.uruchamiacz, okno, zasady, obszar,
		narzedziePandoc, []string{zrodlo, "-o", docelowy}, katalog,
		granicaZamianyFormatu); err != nil {
		return nil, odmowaZamianyFormatu(err)
	}
	bajty, err := os.ReadFile(docelowy)
	if err != nil {
		return nil, bladDebaty(err)
	}
	return bajty, nil
}

// zasiegNarzedziDebaty składa trójkę okno–zasady–obszar dla wywołań arsenału.
//
// Żądania modułu niosą okno debaty, a nie okno sesji terminalowej, więc zasady
// izolacji bierze się z zasięgu platformy — tak samo jak w rodzinie narzędzi
// mediów. Gdy Operator włączy punkt izolacji globalnie, brama zadziała tu tak
// samo jak dla Terminala.
func (a *adapterDebaty) zasiegNarzedziDebaty() (session.Okno, session.Zasady, session.Obszar) {
	okno := session.Okno{Ustawienia: session.Ustawienia{
		SrodowiskoWykonania: shared.ExecutionEnvCore,
	}}
	zasady := session.Zasady{}
	if a.rozstrzygacz != nil {
		zasady = ZasadyIzolacji(a.rozstrzygacz, konfig.Kontekst{})
	}
	return okno, zasady, session.Obszar{}
}

// szablonKontraktu przekłada szablon moderacji na byt kontraktu.
func szablonKontraktu(s dane.SzablonModeracjiDebaty) shared.RoundtableModerationTemplate {
	szablon := shared.RoundtableModerationTemplate{
		Id: s.Kod, Name: s.Nazwa, Format: shared.RoundtableFormat(s.Format),
		SpeakingOrder: s.KolejnoscGlosu, CreatedAt: chwilaBazy(s.Utworzono),
	}
	if s.GranicaTur > 0 {
		granica := s.GranicaTur
		szablon.TurnLimit = &granica
	}
	if s.GranicaCzasuMs > 0 {
		czas := s.GranicaCzasuMs
		szablon.TimeLimitMs = &czas
	}
	return szablon
}
