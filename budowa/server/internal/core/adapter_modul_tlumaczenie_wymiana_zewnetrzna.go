// Wymiana modułu Translate ze światem: budowa i przyjęcie pakietu przekazania,
// most z dokumentem zewnętrznym i wydanie wytworu do biblioteki. Pakiet
// przekazania jest archiwum ZIP złożonym biblioteką archive/zip Go, bez
// programu pakującego.
package core

import (
	"archive/zip"
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"danacoconsole/server/internal/dane"
	"danacoconsole/shared"
)

const (
	// przedrostekPakietuPrzekazania znakuje identyfikator pakietu przekazania
	// nadawany przy zapisie translate.handoff.build.
	przedrostekPakietuPrzekazania = "prz-"
	// przedrostekWytworuTlumaczenia znakuje identyfikator pliku biblioteki
	// założonego przez `artifact.publish`.
	przedrostekWytworuTlumaczenia = "wyt-"
)

// ZlozPakietPrzekazania obsługuje translate.handoff.build: składa archiwum
// przekazania z paneli wskazanego okna i zapisuje je na dysku.
func (a *adapterTlumaczenia) ZlozPakietPrzekazania(ctx context.Context,
	z shared.TranslateHandoffBuildRequest) (shared.TranslateHandoffBuildResponse, error) {

	okno, err := a.repozytorium.Okno(ctx, z.WindowId)
	if err != nil {
		return shared.TranslateHandoffBuildResponse{}, bladNieznanegoOkna(z.WindowId, err)
	}
	if len(z.Contents) == 0 {
		return shared.TranslateHandoffBuildResponse{}, bladWskazaniaTlumaczenia(
			"pakiet bez wskazania zawartości byłby pustym archiwum")
	}
	panele, err := a.repozytorium.Panele(ctx, okno.ID)
	if err != nil {
		return shared.TranslateHandoffBuildResponse{}, bladTlumaczenia(err)
	}
	wskazane := map[string]struct{}{}
	for _, kod := range z.PanelIds {
		wskazane[kod] = struct{}{}
	}
	objete := []dane.PanelTlumaczenia{}
	for _, panel := range panele {
		if len(wskazane) > 0 {
			if _, jest := wskazane[panel.Kod]; !jest {
				continue
			}
		}
		objete = append(objete, panel)
	}
	if len(objete) == 0 {
		return shared.TranslateHandoffBuildResponse{}, bladWskazaniaTlumaczenia(
			"okno nie ma paneli objętych wskazaniem — pakiet nie miałby czego nieść")
	}

	sciezka := strings.TrimSpace(napisZeWskaznika(z.Path))
	if sciezka == "" {
		sciezka = filepath.Join(os.TempDir(), "przekazanie-"+okno.Kod+".zip")
	}

	segmenty, err := a.segmentyOkna(ctx, okno)
	if err != nil {
		return shared.TranslateHandoffBuildResponse{}, err
	}
	bajty, err := a.zawartoscPakietu(ctx, okno, objete, segmenty, z.Contents, z.Instructions)
	if err != nil {
		return shared.TranslateHandoffBuildResponse{}, err
	}
	if err := os.MkdirAll(filepath.Dir(sciezka), 0o700); err != nil {
		return shared.TranslateHandoffBuildResponse{}, bladPlikuTlumaczenia(sciezka, err)
	}
	if err := zapiszPlikWyniku(sciezka, bajty); err != nil {
		return shared.TranslateHandoffBuildResponse{}, err
	}

	kodyPaneli := make([]string, 0, len(objete))
	for _, panel := range objete {
		kodyPaneli = append(kodyPaneli, panel.Kod)
	}
	zawartosci := make([]string, 0, len(z.Contents))
	for _, zawartosc := range z.Contents {
		zawartosci = append(zawartosci, string(zawartosc))
	}
	pakiet, err := a.repozytorium.ZapiszPakietPrzekazania(ctx, dane.PakietPrzekazania{
		Kod:        nowyIdentyfikator(przedrostekPakietuPrzekazania),
		OknoID:     okno.ID,
		Zawartosci: zawartosci,
		Panele:     kodyPaneli,
		Instrukcje: z.Instructions,
		Sciezka:    sciezka,
		Stan:       string(shared.HandoffStatusBuilt),
	})
	if err != nil {
		return shared.TranslateHandoffBuildResponse{}, bladTlumaczenia(err)
	}
	return shared.TranslateHandoffBuildResponse{
		Package: zlozPakietPrzekazania(pakiet),
		Path:    sciezka,
	}, nil
}

// zawartoscPakietu składa archiwum ZIP z wskazanych części: jednostek XLIFF,
// pamięci TMX, terminologii TBX i pliku instrukcji.
func (a *adapterTlumaczenia) zawartoscPakietu(ctx context.Context, okno dane.OknoTlumaczenia,
	panele []dane.PanelTlumaczenia, segmenty []string,
	zawartosci []shared.HandoffContent, instrukcje *string) ([]byte, error) {

	var bufor bytes.Buffer
	archiwum := zip.NewWriter(&bufor)
	dolóz := func(nazwa string, tresc []byte) error {
		strumien, err := archiwum.Create(nazwa)
		if err != nil {
			return bladTlumaczenia(err)
		}
		_, err = strumien.Write(tresc)
		return err
	}

	jezykZrodlowy := jezykZrodlowyZadania(okno.JezykZrodlowy)
	for _, zawartosc := range zawartosci {
		switch zawartosc {
		case shared.HandoffContentXliff:
			for _, panel := range panele {
				jednostki := jednostkiPanelu(segmenty, panel)
				if err := dolóz("xliff/"+panel.Jezyk+".xlf",
					zlozXliff(jezykZrodlowy, panel.Jezyk, jednostki)); err != nil {
					return nil, err
				}
			}
		case shared.HandoffContentMemory:
			wpisy, _, err := a.repozytorium.WpisyPamieci(ctx, dane.FiltrPamieciTlumaczen{})
			if err != nil {
				return nil, bladTlumaczenia(err)
			}
			if err := dolóz("pamiec.tmx", tmxZWpisow(wpisy)); err != nil {
				return nil, err
			}
		case shared.HandoffContentTermbase:
			terminy, err := a.repozytorium.Terminy(ctx)
			if err != nil {
				return nil, bladTlumaczenia(err)
			}
			if err := dolóz("terminologia.tbx", tbxZTerminow(terminy)); err != nil {
				return nil, err
			}
		case shared.HandoffContentInstructions:
			tresc := napisZeWskaznika(instrukcje)
			if strings.TrimSpace(tresc) == "" {
				return nil, bladWskazaniaTlumaczenia(
					"pakiet miał nieść instrukcje, a żądanie ich nie podało")
			}
			if err := dolóz("instrukcje.txt", []byte(tresc)); err != nil {
				return nil, err
			}
		}
	}
	if err := archiwum.Close(); err != nil {
		return nil, bladTlumaczenia(err)
	}
	return bufor.Bytes(), nil
}

// jednostkiPanelu paruje segmenty źródła z akapitami przekładu panelu,
// tworząc jednostki gotowe do zapisu w formacie XLIFF.
func jednostkiPanelu(segmenty []string, panel dane.PanelTlumaczenia) []jednostkaXliff {
	przeklad := []string{}
	if panel.Tresc != nil {
		przeklad = rozdzielAkapity(*panel.Tresc)
	}
	jednostki := make([]jednostkaXliff, 0, len(segmenty))
	for numer, segment := range segmenty {
		jednostka := jednostkaXliff{Zrodlo: segment}
		if numer < len(przeklad) {
			jednostka.Cel = strings.TrimSpace(przeklad[numer])
		}
		jednostki = append(jednostki, jednostka)
	}
	return jednostki
}

// tmxZWpisow składa plik TMX z par pamięci — ta sama postać, którą wypisuje
// `memory.export`, bo to ten sam standard i ten sam odbiorca.
func tmxZWpisow(wpisy []dane.WpisPamieciTlumaczenPelny) []byte {
	var b strings.Builder
	b.WriteString(`<?xml version="1.0" encoding="UTF-8"?>` + "\n")
	b.WriteString(`<tmx version="1.4"><header creationtool="Danaco Console" segtype="sentence" ` +
		`srclang="*all*" o-tmf="DanacoTM" datatype="plaintext"/><body>` + "\n")
	for _, wpis := range wpisy {
		b.WriteString("<tu><tuv xml:lang=\"x-source\"><seg>" + zabezpieczHtml(wpis.SegmentZrodlowy) +
			"</seg></tuv><tuv xml:lang=\"" + wpis.Jezyk + "\"><seg>" +
			zabezpieczHtml(wpis.SegmentDocelowy) + "</seg></tuv></tu>\n")
	}
	b.WriteString("</body></tmx>\n")
	return []byte(b.String())
}

// tbxZTerminow składa bazę terminologiczną w standardzie TBX z wykazu
// terminów, parując źródło z przekładem w każdym wpisie.
func tbxZTerminow(terminy []dane.TerminSlownika) []byte {
	var b strings.Builder
	b.WriteString(`<?xml version="1.0" encoding="UTF-8"?>` + "\n<martif type=\"TBX\"><text><body>\n")
	for numer, termin := range terminy {
		b.WriteString("<termEntry id=\"t" + strings.TrimSpace(termin.Kod) + "\">\n")
		b.WriteString("  <langSet xml:lang=\"x-source\"><tig><term>" +
			zabezpieczHtml(termin.Zrodlo) + "</term></tig></langSet>\n")
		if termin.Cel != nil && strings.TrimSpace(*termin.Cel) != "" {
			b.WriteString("  <langSet xml:lang=\"" + termin.Jezyk + "\"><tig><term>" +
				zabezpieczHtml(*termin.Cel) + "</term></tig></langSet>\n")
		}
		b.WriteString("</termEntry>\n")
		_ = numer
	}
	b.WriteString("</body></text></martif>\n")
	return []byte(b.String())
}

// PrzyjmijPakietPrzekazania obsługuje `translate.handoff.receive`. Czyta zwrot
// wykonawcy — archiwum pakietu albo pojedynczy plik XLIFF — i wnosi przekłady
// do paneli okna.
func (a *adapterTlumaczenia) PrzyjmijPakietPrzekazania(ctx context.Context,
	z shared.TranslateHandoffReceiveRequest) (shared.TranslateHandoffReceiveResponse, error) {

	okno, err := a.repozytorium.Okno(ctx, z.WindowId)
	if err != nil {
		return shared.TranslateHandoffReceiveResponse{}, bladNieznanegoOkna(z.WindowId, err)
	}
	sciezka := strings.TrimSpace(z.Path)
	if sciezka == "" {
		return shared.TranslateHandoffReceiveResponse{}, bladWskazaniaTlumaczenia(
			"żądanie translate.handoff.receive bez ścieżki zwrotu")
	}
	bajty, err := os.ReadFile(filepath.Clean(sciezka))
	if err != nil {
		return shared.TranslateHandoffReceiveResponse{}, bladPlikuTlumaczenia(sciezka, err)
	}

	pliki := [][]byte{}
	if archiwum, err := zip.NewReader(bytes.NewReader(bajty), int64(len(bajty))); err == nil {
		for _, wpis := range archiwum.File {
			if !strings.HasSuffix(strings.ToLower(wpis.Name), ".xlf") &&
				!strings.HasSuffix(strings.ToLower(wpis.Name), ".xliff") {
				continue
			}
			strumien, err := wpis.Open()
			if err != nil {
				return shared.TranslateHandoffReceiveResponse{}, bladTlumaczenia(err)
			}
			tresc, err := io.ReadAll(strumien)
			strumien.Close()
			if err != nil {
				return shared.TranslateHandoffReceiveResponse{}, bladTlumaczenia(err)
			}
			pliki = append(pliki, tresc)
		}
	} else {
		pliki = append(pliki, bajty)
	}
	if len(pliki) == 0 {
		return shared.TranslateHandoffReceiveResponse{}, bladWskazaniaTlumaczenia(
			"zwrot nie niesie ani jednego pliku XLIFF — nie ma czego przyjąć")
	}

	wniesione := 0
	for _, plik := range pliki {
		odczytany := rozbierzXliff(plik)
		if odczytany.JezykDocelowy == "" || len(odczytany.Jednostki) == 0 {
			continue
		}
		tresci := make([]string, 0, len(odczytany.Jednostki))
		for _, jednostka := range odczytany.Jednostki {
			if strings.TrimSpace(jednostka.Cel) == "" {
				tresci = append(tresci, jednostka.Zrodlo)
				continue
			}
			tresci = append(tresci, jednostka.Cel)
		}
		if _, err := a.panelJezyka(ctx, okno, odczytany.JezykDocelowy, true,
			strings.Join(tresci, "\n\n")); err != nil {
			return shared.TranslateHandoffReceiveResponse{}, err
		}
		wniesione += len(odczytany.Jednostki)
	}

	// Pakiet, którego zwrot dotyczy, przechodzi w stan `returned`, żeby
	// wiadomo było, czy wrócił.
	if kod := strings.TrimSpace(napisZeWskaznika(z.PackageId)); kod != "" {
		pakiet, err := a.repozytorium.PakietPrzekazania(ctx, kod)
		if err == nil {
			pakiet.Stan = string(shared.HandoffStatusReturned)
			if _, err := a.repozytorium.ZapiszPakietPrzekazania(ctx, pakiet); err != nil {
				return shared.TranslateHandoffReceiveResponse{}, bladTlumaczenia(err)
			}
		}
	}

	if z.RunQualityCheck != nil && *z.RunQualityCheck {
		panele, err := a.repozytorium.Panele(ctx, okno.ID)
		if err != nil {
			return shared.TranslateHandoffReceiveResponse{}, bladTlumaczenia(err)
		}
		tekstZrodlowy := ""
		if okno.TekstZrodlowy != nil {
			tekstZrodlowy = *okno.TekstZrodlowy
		}
		for _, panel := range panele {
			_ = a.repozytorium.ZapiszNiezgodnosci(ctx, panel.ID, zbadajPanel(panel, tekstZrodlowy))
		}
	}

	panele, err := a.repozytorium.Panele(ctx, okno.ID)
	if err != nil {
		return shared.TranslateHandoffReceiveResponse{}, bladTlumaczenia(err)
	}
	return shared.TranslateHandoffReceiveResponse{
		Panels:        zlozPaneleTlumaczenia(panele),
		ImportedCount: wniesione,
	}, nil
}

// zlozPakietPrzekazania przekłada wiersz pakietu przekazania z bazy na byt
// kontraktu zwracany w odpowiedzi komendy.
func zlozPakietPrzekazania(pakiet dane.PakietPrzekazania) shared.HandoffPackage {
	zawartosci := make([]shared.HandoffContent, 0, len(pakiet.Zawartosci))
	for _, zawartosc := range pakiet.Zawartosci {
		zawartosci = append(zawartosci, shared.HandoffContent(zawartosc))
	}
	panele := pakiet.Panele
	if panele == nil {
		panele = []string{}
	}
	return shared.HandoffPackage{
		Id:        pakiet.Kod,
		WindowId:  pakiet.OknoKod,
		Contents:  zawartosci,
		PanelIds:  panele,
		Status:    shared.HandoffStatus(pakiet.Stan),
		CreatedAt: pakiet.Utworzono,
	}
}

// PrzyjmijZrodloMostu obsługuje `translate.bridge.source.receive`. Wiąże okno
// z dokumentem, z którego materiał przyszedł, i wpisuje treść jako źródło.
func (a *adapterTlumaczenia) PrzyjmijZrodloMostu(ctx context.Context,
	z shared.TranslateBridgeSourceReceiveRequest) (shared.TranslateBridgeSourceReceiveResponse, error) {

	if strings.TrimSpace(z.Text) == "" {
		return shared.TranslateBridgeSourceReceiveResponse{}, bladWskazaniaTlumaczenia(
			"most bez treści nie wnosi materiału")
	}
	if strings.TrimSpace(z.DocumentId) == "" {
		return shared.TranslateBridgeSourceReceiveResponse{}, bladWskazaniaTlumaczenia(
			"most bez wskazania dokumentu nie wiedziałby, dokąd odesłać wynik")
	}
	segmenty := podzielNaZdania(z.Text)
	liczba := int64(len(segmenty))
	tekst := z.Text
	okno, err := a.repozytorium.ZapiszOkno(ctx, dane.OknoTlumaczenia{
		Kod:             z.WindowId,
		TekstZrodlowy:   &tekst,
		LiczbaSegmentow: &liczba,
	})
	if err != nil {
		return shared.TranslateBridgeSourceReceiveResponse{}, bladTlumaczenia(err)
	}
	if err := a.repozytorium.UstawSegmentyOkna(ctx, okno.ID, segmenty); err != nil {
		return shared.TranslateBridgeSourceReceiveResponse{}, bladTlumaczenia(err)
	}
	if err := a.repozytorium.ZapiszMost(ctx, dane.MostTlumaczenia{
		OknoID:        okno.ID,
		DokumentKod:   z.DocumentId,
		Zakotwiczenie: z.Anchor,
	}); err != nil {
		return shared.TranslateBridgeSourceReceiveResponse{}, bladTlumaczenia(err)
	}
	return shared.TranslateBridgeSourceReceiveResponse{
		WindowId:     okno.Kod,
		SegmentCount: len(segmenty),
	}, nil
}

// OdesljWynikMostu obsługuje translate.bridge.result.send. Składa treść
// wskazanych paneli w postaci żądanej przez tryb i odkłada ją jako wytwór
// mostu obok dokumentu źródłowego, bo dokument należy do innego modułu.
func (a *adapterTlumaczenia) OdesljWynikMostu(ctx context.Context,
	z shared.TranslateBridgeResultSendRequest) (shared.TranslateBridgeResultSendResponse, error) {

	okno, err := a.repozytorium.Okno(ctx, z.WindowId)
	if err != nil {
		return shared.TranslateBridgeResultSendResponse{}, bladNieznanegoOkna(z.WindowId, err)
	}
	most, err := a.repozytorium.Most(ctx, okno.ID)
	if err != nil {
		return shared.TranslateBridgeResultSendResponse{}, bladWskazaniaTlumaczenia(
			"okno " + okno.Kod + " nie jest związane z żadnym dokumentem — " +
				"most zakłada `translate.bridge.source.receive`")
	}
	if len(z.PanelIds) == 0 {
		return shared.TranslateBridgeResultSendResponse{}, bladWskazaniaTlumaczenia(
			"odesłanie bez wskazania paneli nie niosłoby przekładu")
	}
	segmenty, err := a.segmentyOkna(ctx, okno)
	if err != nil {
		return shared.TranslateBridgeResultSendResponse{}, err
	}

	czesci := []string{}
	wyslane := 0
	for _, kod := range z.PanelIds {
		panel, err := a.repozytorium.Panel(ctx, kod)
		if err != nil {
			return shared.TranslateBridgeResultSendResponse{}, bladNieznanegoPanelu(kod, err)
		}
		tresc, err := trescPanelu(panel)
		if err != nil {
			return shared.TranslateBridgeResultSendResponse{}, err
		}
		switch z.Mode {
		case shared.BridgeResultModeBilingual:
			przeklad := rozdzielAkapity(tresc)
			for numer, segment := range segmenty {
				czesci = append(czesci, segment)
				if numer < len(przeklad) {
					czesci = append(czesci, przeklad[numer])
				}
			}
		case shared.BridgeResultModeTargetOnly:
			czesci = append(czesci, tresc)
		case shared.BridgeResultModeAppendix:
			czesci = append(czesci, "— przekład na język "+panel.Jezyk+" —", tresc)
		default:
			return shared.TranslateBridgeResultSendResponse{}, bladWskazaniaTlumaczenia(
				"nieznany tryb odesłania wyniku: " + string(z.Mode))
		}
		wyslane++
	}

	sciezka := filepath.Join(os.TempDir(), "most-"+most.DokumentKod+"-"+okno.Kod+".md")
	if err := zapiszPlikWyniku(sciezka, []byte(strings.Join(czesci, "\n\n")+"\n")); err != nil {
		return shared.TranslateBridgeResultSendResponse{}, err
	}
	return shared.TranslateBridgeResultSendResponse{
		DocumentId: most.DokumentKod,
		SentCount:  wyslane,
	}, nil
}

// WydajWytwor obsługuje translate.artifact.publish: odkłada wytwór do
// magazynu treści rdzenia i zakłada wiersz pliku biblioteki — jedyną drogę,
// którą wytwór staje się widoczny dla reszty platformy.
func (a *adapterTlumaczenia) WydajWytwor(ctx context.Context,
	z shared.TranslateArtifactPublishRequest) (shared.TranslateArtifactPublishResponse, error) {

	if a.biblioteka == nil || a.magazynWytworow == nil {
		return shared.TranslateArtifactPublishResponse{}, bladTlumaczenia(
			errBrakMagazynuWytworow)
	}
	okno, err := a.repozytorium.Okno(ctx, z.WindowId)
	if err != nil {
		return shared.TranslateArtifactPublishResponse{}, bladNieznanegoOkna(z.WindowId, err)
	}
	panele, err := a.repozytorium.Panele(ctx, okno.ID)
	if err != nil {
		return shared.TranslateArtifactPublishResponse{}, bladTlumaczenia(err)
	}
	wskazane := map[string]struct{}{}
	for _, kod := range z.PanelIds {
		wskazane[kod] = struct{}{}
	}
	objete := []dane.PanelTlumaczenia{}
	for _, panel := range panele {
		if len(wskazane) > 0 {
			if _, jest := wskazane[panel.Kod]; !jest {
				continue
			}
		}
		objete = append(objete, panel)
	}

	segmenty, err := a.segmentyOkna(ctx, okno)
	if err != nil {
		return shared.TranslateArtifactPublishResponse{}, err
	}

	bajty, nazwa, err := a.trescWytworu(ctx, okno, objete, segmenty, z.Kind)
	if err != nil {
		return shared.TranslateArtifactPublishResponse{}, err
	}

	suma := sha256.Sum256(bajty)
	odcisk := hex.EncodeToString(suma[:])
	sciezka, err := a.magazynWytworow.Zapisz(bajty, odcisk)
	if err != nil {
		return shared.TranslateArtifactPublishResponse{}, bladTlumaczenia(err)
	}

	rozmiar := int64(len(bajty))
	modul := "translate"
	plik, err := a.biblioteka.ZapiszPlik(ctx, dane.PlikBiblioteki{
		Kod:             nowyIdentyfikator(przedrostekWytworuTlumaczenia),
		Nazwa:           nazwa,
		RozmiarBajtow:   &rozmiar,
		ModulZrodlowyID: &modul,
		SumaKontrolna:   &odcisk,
		TrescOdwolanie:  &sciezka,
	})
	if err != nil {
		return shared.TranslateArtifactPublishResponse{}, bladTlumaczenia(err)
	}
	if kolekcja := strings.TrimSpace(napisZeWskaznika(z.CollectionId)); kolekcja != "" {
		if _, err := a.biblioteka.PrzypiszDoKolekcji(ctx, kolekcja, []string{plik.Kod}); err != nil {
			return shared.TranslateArtifactPublishResponse{}, bladTlumaczenia(err)
		}
	}
	if len(z.Tags) > 0 {
		if _, err := a.biblioteka.UstawEtykiety(ctx, plik.Kod, z.Tags); err != nil {
			return shared.TranslateArtifactPublishResponse{}, bladTlumaczenia(err)
		}
	}
	return shared.TranslateArtifactPublishResponse{FileId: plik.Kod, Path: sciezka}, nil
}

// trescWytworu składa bajty wytworu wskazanego rodzaju — pliku docelowego,
// dwujęzycznego, pamięci albo terminologii — wraz z nazwą pliku.
func (a *adapterTlumaczenia) trescWytworu(ctx context.Context, okno dane.OknoTlumaczenia,
	panele []dane.PanelTlumaczenia, segmenty []string,
	rodzaj shared.TranslationArtifactKind) ([]byte, string, error) {

	switch rodzaj {
	case shared.TranslationArtifactKindTargetFile:
		if len(panele) == 0 {
			return nil, "", bladWskazaniaTlumaczenia("okno nie ma paneli — wytwór nie miałby treści")
		}
		czesci := []string{}
		for _, panel := range panele {
			tresc, err := trescPanelu(panel)
			if err != nil {
				return nil, "", err
			}
			czesci = append(czesci, tresc)
		}
		return []byte(strings.Join(czesci, "\n\n")), okno.Kod + "-przeklad.md", nil

	case shared.TranslationArtifactKindBilingualFile:
		if len(panele) == 0 {
			return nil, "", bladWskazaniaTlumaczenia("okno nie ma paneli — wytwór nie miałby treści")
		}
		var b strings.Builder
		for _, panel := range panele {
			b.WriteString("## " + panel.Jezyk + "\n\n")
			przeklad := []string{}
			if panel.Tresc != nil {
				przeklad = rozdzielAkapity(*panel.Tresc)
			}
			for numer, segment := range segmenty {
				b.WriteString("- " + segment + "\n")
				if numer < len(przeklad) {
					b.WriteString("  → " + strings.TrimSpace(przeklad[numer]) + "\n")
				}
			}
			b.WriteString("\n")
		}
		return []byte(b.String()), okno.Kod + "-dwujezyczny.md", nil

	case shared.TranslationArtifactKindMemory:
		wpisy, _, err := a.repozytorium.WpisyPamieci(ctx, dane.FiltrPamieciTlumaczen{})
		if err != nil {
			return nil, "", bladTlumaczenia(err)
		}
		return tmxZWpisow(wpisy), okno.Kod + "-pamiec.tmx", nil

	case shared.TranslationArtifactKindTermbase:
		terminy, err := a.repozytorium.Terminy(ctx)
		if err != nil {
			return nil, "", bladTlumaczenia(err)
		}
		return tbxZTerminow(terminy), okno.Kod + "-terminologia.tbx", nil
	}
	return nil, "", bladWskazaniaTlumaczenia("nieznany rodzaj wytworu: " + string(rodzaj))
}

// WykazKrokow obsługuje translate.step.list: oddaje kroki pracy modułu wraz
// z komendą i polami, które ona przyjmuje. Wykaz stoi w kodzie, nie w bazie,
// bo opisuje zdolności rdzenia, nie dane Operatora.
func (a *adapterTlumaczenia) WykazKrokow(_ context.Context,
	z shared.TranslateStepListRequest) (shared.TranslateStepListResponse, error) {

	kroki := []shared.TranslationStep{
		{Command: string(shared.CommandTranslateSourceSet), Name: "Ustal tekst źródłowy",
			Kind:        shared.TranslationStepKindTranslation,
			Description: "Zapisuje tekst źródłowy okna i liczy jego segmenty",
			Parameters:  []string{"windowId", "text", "sourceLanguage", "resegment"}},
		{Command: string(shared.CommandTranslateTargetAdd), Name: "Dodaj język docelowy",
			Kind:        shared.TranslationStepKindTranslation,
			Description: "Zakłada panel języka i tłumaczy do niego tekst źródłowy kanałem modelu",
			Parameters:  []string{"windowId", "language", "tone", "channelId"}},
		{Command: string(shared.CommandTranslateMemoryPretranslate), Name: "Tłumaczenie wstępne z pamięci",
			Kind:        shared.TranslationStepKindMemory,
			Description: "Wypełnia panele parami pamięci powyżej progu dopasowania",
			Parameters:  []string{"windowId", "panelId", "threshold", "onlyEmpty"}},
		{Command: string(shared.CommandTranslateMemorySet), Name: "Dopisz parę do pamięci",
			Kind:        shared.TranslationStepKindMemory,
			Description: "Wnosi parę segmentów do pamięci tłumaczeń",
			Parameters:  []string{"language", "sourceSegment", "targetSegment", "project"}},
		{Command: string(shared.CommandTranslateQualityCheck), Name: "Kontrola jakości",
			Kind:        shared.TranslationStepKindQuality,
			Description: "Sprawdza liczby, daty, waluty, znaczniki, długość i pominięcia",
			Parameters:  []string{"panelId"}},
		{Command: string(shared.CommandTranslateProofreadRun), Name: "Korekta językowa",
			Kind:        shared.TranslationStepKindQuality,
			Description: "Zakłada ustalenia korekty wraz z propozycjami poprawek",
			Parameters:  []string{"panelId", "checks"}},
		{Command: string(shared.CommandTranslateConsistencyCheck), Name: "Kontrola spójności",
			Kind:        shared.TranslationStepKindQuality,
			Description: "Szuka zdań i terminów przełożonych różnie w obrębie jednego języka",
			Parameters:  []string{"windowId", "panelId"}},
		{Command: string(shared.CommandTranslateApprovalSet), Name: "Etap zatwierdzenia",
			Kind:        shared.TranslationStepKindQuality,
			Description: "Przestawia panel na kolejny etap obiegu i zapisuje krok obiegu",
			Parameters:  []string{"panelId", "stage", "note"}},
		{Command: string(shared.CommandTranslateDocumentRender), Name: "Złóż dokument wyniku",
			Kind:        shared.TranslationStepKindExport,
			Description: "Składa plik dokumentu z treści panelu",
			Parameters:  []string{"documentId", "panelId", "format", "path"}},
		{Command: string(shared.CommandTranslateSubtitleExport), Name: "Wydaj napisy",
			Kind:        shared.TranslationStepKindExport,
			Description: "Zapisuje napisy w formacie SRT, WebVTT, TTML albo EBU STL",
			Parameters:  []string{"panelId", "path", "format"}},
		{Command: string(shared.CommandTranslateResourceExport), Name: "Wydaj zasób lokalizacyjny",
			Kind:        shared.TranslationStepKindExport,
			Description: "Zapisuje plik kluczy w języku panelu",
			Parameters:  []string{"resourceId", "panelId", "path", "format"}},
		{Command: string(shared.CommandTranslateArtifactPublish), Name: "Wydaj wytwór do biblioteki",
			Kind:        shared.TranslationStepKindExport,
			Description: "Odkłada wytwór w magazynie rdzenia i zakłada wiersz pliku biblioteki",
			Parameters:  []string{"windowId", "kind", "panelIds", "collectionId", "tags"}},
		{Command: string(shared.CommandTranslateMemoryExport), Name: "Wydaj pamięć tłumaczeń",
			Kind:        shared.TranslationStepKindMemory,
			Description: "Wypisuje pary pamięci do pliku TMX albo CSV",
			Parameters:  []string{"path", "language", "project"}},
	}
	if z.Kind == nil {
		return shared.TranslateStepListResponse{Steps: kroki}, nil
	}
	zawezone := []shared.TranslationStep{}
	for _, krok := range kroki {
		if krok.Kind == *z.Kind {
			zawezone = append(zawezone, krok)
		}
	}
	sort.SliceStable(zawezone, func(i, j int) bool { return zawezone[i].Name < zawezone[j].Name })
	return shared.TranslateStepListResponse{Steps: zawezone}, nil
}
