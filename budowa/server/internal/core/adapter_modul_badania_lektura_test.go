package core

import (
	"context"
	"encoding/base64"
	"path/filepath"
	"testing"

	"danacoconsole/server/internal/dane"
	"danacoconsole/server/internal/store"
	"danacoconsole/shared"
)

// fakeDokumentyLektury podstawia port Dokumenty pod adapter Badań: zamiast
// wołać programy zewnętrzne, zapamiętuje ostatnie zadanie, żeby dało się
// zmierzyć, czy pola `preprocess` i `languages` doszły do portu, zamiast
// zostać po drodze upuszczone.
type fakeDokumentyLektury struct {
	ostatnie shared.DocumentTextExtractRequest
}

func (f *fakeDokumentyLektury) Przeksztalc(ctx context.Context,
	z shared.DocumentConvertRequest) (shared.DocumentConvertResponse, error) {
	return shared.DocumentConvertResponse{}, nil
}

func (f *fakeDokumentyLektury) WyciagnijTekst(ctx context.Context,
	z shared.DocumentTextExtractRequest) (shared.DocumentTextExtractResponse, error) {
	f.ostatnie = z
	return shared.DocumentTextExtractResponse{Text: "tekst rozpoznany", UsedOcr: true}, nil
}

// zlozAdapterBadan składa adapter modułu Research nad świeżą bazą SQLite,
// bez pełnego montażu rdzenia — sprawdzianowi tego pliku wystarczy jeden port.
func zlozAdapterBadan(t *testing.T) (*adapterBadan, dane.RepozytoriumBadan, context.Context) {
	t.Helper()

	katalog := t.TempDir()
	baza, err := store.Otworz(filepath.Join(katalog, "dane.sqlite"))
	if err != nil {
		t.Fatalf("nie można otworzyć bazy sprawdzianu: %v", err)
	}
	t.Cleanup(func() { _ = baza.Zamknij() })
	if err := baza.Migruj(); err != nil {
		t.Fatalf("migracje nie doszły do skutku: %v", err)
	}
	zycie := context.Background()
	repozytoria, err := dane.Otworz(zycie, baza)
	if err != nil {
		t.Fatalf("nie można złożyć repozytoriów: %v", err)
	}

	adapter := nowyAdapterBadan(repozytoria.Badania).ZKatalogiemDanych(katalog)
	return adapter, repozytoria.Badania, zycie
}

// TestRozpoznaniePismaZrodlaPrzekazujePreprocessIJezyki wykazuje, że
// `research.source.ocr` przekazuje pola `preprocess` i `languages` żądania
// nietknięte do `document.text.extract` — bez tego uchwytu jedno albo drugie
// pole ginie po drodze, a Operator dostaje rozpoznanie inną drogą niż zamówił.
// Miarę skutku obróbki wstępnej na materiale produkcyjnym niesie
// TestWyciagnijTekstProstujeSkosPrzedRozpoznaniem.
func TestRozpoznaniePismaZrodlaPrzekazujePreprocessIJezyki(t *testing.T) {
	adapter, _, zycie := zlozAdapterBadan(t)
	fake := &fakeDokumentyLektury{}
	adapter.ZDokumentami(fake)

	dodane, err := adapter.DodajZrodlo(zycie, shared.ResearchSourceAddRequest{
		WindowId: "okno-obrobki-wstepnej", Title: "skan protokołu odbioru",
	})
	if err != nil {
		t.Fatalf("dodanie źródła nie powiodło się: %v", err)
	}
	bajty := base64.StdEncoding.EncodeToString([]byte("bajty skanu"))
	if _, err := adapter.DodajZalacznikZrodla(zycie, shared.ResearchSourceAttachmentAddRequest{
		SourceId: dodane.Source.Id, Kind: shared.ResearchAttachmentKindFulltext,
		ContentBase64: &bajty,
	}); err != nil {
		t.Fatalf("dodanie załącznika nie powiodło się: %v", err)
	}

	if _, err := adapter.RozpoznajPismoZrodla(zycie, shared.ResearchSourceOcrRequest{
		SourceId: dodane.Source.Id,
	}); err != nil {
		t.Fatalf("source.ocr bez obróbki odmówiło: %v", err)
	}
	if fake.ostatnie.Preprocess != nil && *fake.ostatnie.Preprocess {
		t.Fatalf("preprocess doszedł do portu jako prawda, mimo że żądanie go nie zamówiło")
	}
	if fake.ostatnie.Language != nil {
		t.Fatalf("language doszedł do portu, mimo że żądanie nie podało languages: %q",
			*fake.ostatnie.Language)
	}

	prawda := true
	if _, err := adapter.RozpoznajPismoZrodla(zycie, shared.ResearchSourceOcrRequest{
		SourceId: dodane.Source.Id, Preprocess: &prawda, Languages: []string{"pol", "eng"},
	}); err != nil {
		t.Fatalf("source.ocr z obróbką odmówiło: %v", err)
	}
	if fake.ostatnie.Preprocess == nil || !*fake.ostatnie.Preprocess {
		t.Fatalf("pole preprocess nie doszło do portu dokumentów jako prawda: %+v",
			fake.ostatnie.Preprocess)
	}
	if fake.ostatnie.Language == nil || *fake.ostatnie.Language != "pol+eng" {
		t.Fatalf("pole languages nie doszło do portu dokumentów jako language: %+v",
			fake.ostatnie.Language)
	}
}
