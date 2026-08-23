package core

import (
	"encoding/json"
	"strings"
	"testing"

	"danacoconsole/shared"
)

// Skutek modułu Studio: czy wpis katalogowy ZOSTAJE i czy szablon coś zakłada.
//
// Szkoda, którą ten plik ma wykluczyć: zapis, który wraca `ok` i nie zostawia
// wiersza. Rodzina katalogowa jest na nią podatna, bo jej odpowiedzi niosą byt
// złożony z żądania — wystarczyłoby oddać to, co przyszło, żeby wszystko
// wyglądało poprawnie. Dlatego każdy sprawdzian pyta o wpis OSOBNYM odczytem,
// a szablon mierzy treścią dokumentu, który z niego powstał.

// TestOperacjaWlasnaZostajeWKatalogu sprawdza pełny obieg: zapis, odczyt,
// nadpisanie, usunięcie — i to, że usunięcie bytu nieistniejącego odmawia.
func TestOperacjaWlasnaZostajeWKatalogu(t *testing.T) {
	zmontowany, zycie, _ := zmontujDoPomiaruSkutku(t)

	var zapisana shared.StudioOperationSaveResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandStudioOperationSave,
		shared.StudioOperationSaveRequest{
			Name: "Skrót zarządczy", Category: "streszczenie",
			Prompt: "Streść dokument w trzech zdaniach dla zarządu.",
		}, &zapisana)
	if zapisana.Operation.Id == "" {
		t.Fatal("zapis operacji nie nadał identyfikatora")
	}
	if zapisana.Operation.Builtin {
		t.Error("operacja zapisana przez Operatora jest oznaczona jako fabryczna")
	}

	var wykaz shared.StudioOperationListResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandStudioOperationList,
		shared.StudioOperationListRequest{}, &wykaz)
	if len(wykaz.Operations) != 1 || wykaz.Operations[0].Id != zapisana.Operation.Id {
		t.Fatalf("katalog niesie %d operacji, oczekiwano zapisanej %s",
			len(wykaz.Operations), zapisana.Operation.Id)
	}
	if wykaz.Operations[0].Prompt != "Streść dokument w trzech zdaniach dla zarządu." {
		t.Errorf("operacja w katalogu niesie prompt %q", wykaz.Operations[0].Prompt)
	}

	// Nadpisanie po identyfikatorze zmienia wpis, a nie zakłada drugiego.
	wykonajUdana(t, zmontowany, zycie, shared.CommandStudioOperationSave,
		shared.StudioOperationSaveRequest{
			OperationId: &zapisana.Operation.Id, Name: "Skrót zarządczy",
			Category: "streszczenie", Prompt: "Streść dokument w pięciu zdaniach.",
		}, &shared.StudioOperationSaveResponse{})
	wykonajUdana(t, zmontowany, zycie, shared.CommandStudioOperationList,
		shared.StudioOperationListRequest{}, &wykaz)
	if len(wykaz.Operations) != 1 {
		t.Fatalf("nadpisanie założyło drugi wpis: katalog niesie %d operacji", len(wykaz.Operations))
	}
	if !strings.Contains(wykaz.Operations[0].Prompt, "pięciu") {
		t.Errorf("nadpisanie nie zmieniło promptu: %q", wykaz.Operations[0].Prompt)
	}

	// Zawężenie po kategorii ma odsiewać.
	wykonajUdana(t, zmontowany, zycie, shared.CommandStudioOperationList,
		shared.StudioOperationListRequest{Category: wskaznik("korekta")}, &wykaz)
	if len(wykaz.Operations) != 0 {
		t.Errorf("zawężenie po obcej kategorii oddało %d operacji", len(wykaz.Operations))
	}

	var usunieta shared.StudioOperationDeleteResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandStudioOperationDelete,
		shared.StudioOperationDeleteRequest{OperationId: zapisana.Operation.Id}, &usunieta)
	if !usunieta.Deleted {
		t.Error("usunięcie oddało „nie usunięto” mimo istniejącego wpisu")
	}
	wykonajUdana(t, zmontowany, zycie, shared.CommandStudioOperationList,
		shared.StudioOperationListRequest{}, &wykaz)
	if len(wykaz.Operations) != 0 {
		t.Errorf("po usunięciu katalog niesie %d operacji", len(wykaz.Operations))
	}

	// Usunięcie bytu, którego nie ma, odmawia — a nie melduje „nie usunięto”.
	odpowiedz := wykonajKomende(t, zmontowany, zycie, shared.CommandStudioOperationDelete,
		shared.StudioOperationDeleteRequest{OperationId: zapisana.Operation.Id})
	if odpowiedz.Error == nil || odpowiedz.Error.Code != shared.ErrorCodeNotFound {
		t.Errorf("powtórne usunięcie nie odmówiło kodem not_found: %+v", odpowiedz.Error)
	}
}

// TestLancuchZostajeZKrokamiWKolejnosci sprawdza, że kroki przeżywają zapis
// i odczyt w tej samej kolejności — sekwencja przestawiona to inna sekwencja.
func TestLancuchZostajeZKrokamiWKolejnosci(t *testing.T) {
	zmontowany, zycie, _ := zmontujDoPomiaruSkutku(t)

	kroki := []shared.StudioChainStep{
		{ActionId: "studio.korekta.ortografia"},
		{ActionId: "studio.streszczenie.akapitowe"},
		{ActionId: "studio.styl.rejestr"},
	}
	var zapisany shared.StudioChainSaveResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandStudioChainSave,
		shared.StudioChainSaveRequest{Name: "Redakcja finalna", Steps: kroki}, &zapisany)

	var wykaz shared.StudioChainListResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandStudioChainList,
		shared.StudioChainListRequest{}, &wykaz)
	if len(wykaz.Chains) != 1 {
		t.Fatalf("katalog niesie %d łańcuchów, oczekiwano jednego", len(wykaz.Chains))
	}
	odczytane := wykaz.Chains[0].Steps
	if len(odczytane) != len(kroki) {
		t.Fatalf("łańcuch odczytany ma %d kroków, zapisano %d", len(odczytane), len(kroki))
	}
	for i, krok := range kroki {
		if odczytane[i].ActionId != krok.ActionId {
			t.Errorf("krok %d to %q, zapisano %q", i, odczytane[i].ActionId, krok.ActionId)
		}
	}

	// Łańcuch bez kroków nie jest łańcuchem.
	odpowiedz := wykonajKomende(t, zmontowany, zycie, shared.CommandStudioChainSave,
		shared.StudioChainSaveRequest{Name: "Pusty", Steps: []shared.StudioChainStep{}})
	if odpowiedz.Error == nil {
		t.Error("zapis łańcucha bez kroków wrócił odpowiedzią udaną")
	}

	// Uruchomienie liczy kroki z DEFINICJI, więc musi się zgadzać z zapisem.
	dokument := dokumentZTrescia(t, zmontowany, zycie, "okno-1", "Treść dokumentu.")
	var przebieg shared.StudioChainRunResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandStudioChainRun,
		shared.StudioChainRunRequest{
			WindowId: "okno-1", DocumentId: dokument.Id, ChainId: zapisany.Chain.Id,
			Scope: shared.StudioOperationScopeDocument,
		}, &przebieg)
	if przebieg.Steps != len(kroki) {
		t.Errorf("uruchomienie zapowiada %d kroków, łańcuch ma %d", przebieg.Steps, len(kroki))
	}
	if przebieg.RunId == "" {
		t.Error("uruchomienie nie nadało identyfikatora przebiegu")
	}

	// Zakres „zaznaczenie” bez granic jest żądaniem sprzecznym.
	odpowiedz = wykonajKomende(t, zmontowany, zycie, shared.CommandStudioChainRun,
		shared.StudioChainRunRequest{
			WindowId: "okno-1", DocumentId: dokument.Id, ChainId: zapisany.Chain.Id,
			Scope: shared.StudioOperationScopeSelection,
		})
	if odpowiedz.Error == nil {
		t.Error("uruchomienie z zakresem zaznaczenia bez granic wróciło odpowiedzią udaną")
	}
}

// TestSzablonFabrycznyZakladaDokumentZWypelnionymiPolami mierzy szablon tam,
// gdzie ma skutek: w treści dokumentu, który z niego powstał.
func TestSzablonFabrycznyZakladaDokumentZWypelnionymiPolami(t *testing.T) {
	zmontowany, zycie, _ := zmontujDoPomiaruSkutku(t)

	var szablony shared.StudioTemplateListResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandStudioTemplateList,
		shared.StudioTemplateListRequest{}, &szablony)
	if len(szablony.Templates) < 5 {
		t.Fatalf("katalog niesie %d szablonów, opracowanie wymienia pięć układów fabrycznych",
			len(szablony.Templates))
	}

	var raport *shared.StudioTemplate
	for i, szablon := range szablony.Templates {
		if szablon.Name == "Raport" {
			raport = &szablony.Templates[i]
		}
		if !szablon.Builtin {
			t.Errorf("szablon %q nie jest oznaczony jako fabryczny", szablon.Name)
		}
		if len(szablon.Fields) == 0 {
			t.Errorf("szablon %q nie niesie wykazu pól — okno musiałoby zgadywać je z treści",
				szablon.Name)
		}
	}
	if raport == nil {
		t.Fatal("katalog nie niesie szablonu „Raport”")
	}

	var powstaly shared.StudioTemplateApplyResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandStudioTemplateApply,
		shared.StudioTemplateApplyRequest{
			WindowId: "okno-1", TemplateId: raport.Id,
			Title: wskaznik("Raport z przeglądu"),
			Values: json.RawMessage(`{"tytul":"Przegląd stanu wdrożenia",` +
				`"autor":"Operator","streszczenie":"Wdrożenie przebiega zgodnie z planem."}`),
		}, &powstaly)

	if powstaly.Document.Content == nil || *powstaly.Document.Content == "" {
		t.Fatal("zastosowanie szablonu założyło dokument bez treści")
	}
	tresc := *powstaly.Document.Content
	if !strings.Contains(tresc, "Przegląd stanu wdrożenia") {
		t.Errorf("pole „tytul” nie zostało podstawione: %q", tresc[:min(120, len(tresc))])
	}
	if strings.Contains(tresc, "{{tytul}}") {
		t.Error("znacznik pola wypełnionego został w treści")
	}
	// Pole niewypełnione ZOSTAJE widoczne — dokument z pustym miejscem po polu
	// wyglądałby na kompletny, a nie jest.
	if !strings.Contains(tresc, "{{wnioski}}") {
		t.Error("znacznik pola niewypełnionego zniknął z treści zamiast zostać widoczny")
	}
	if powstaly.Document.VersionId == nil || *powstaly.Document.VersionId == "" {
		t.Error("dokument z szablonu nie ma wersji pierwszej")
	}

	// Dokument odczytany osobno niesie tę samą treść.
	if trescDokumentu(t, zmontowany, zycie, "okno-1", powstaly.Document.Id) != tresc {
		t.Error("dokument z szablonu po ponownym otwarciu niesie inną treść")
	}
}

// TestProfilWydaniaZachowujeUstawieniaStrony sprawdza, że ustawienia strony
// przeżywają zapis i odczyt — profil bez nich jest samą nazwą formatu.
func TestProfilWydaniaZachowujeUstawieniaStrony(t *testing.T) {
	zmontowany, zycie, _ := zmontujDoPomiaruSkutku(t)

	var orientacja shared.StudioPageOrientation = shared.StudioPageOrientationPionowa
	var zapisany shared.StudioExportProfileSaveResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandStudioExportProfileSave,
		shared.StudioExportProfileSaveRequest{
			Name: "PDF do druku", Format: "pdf",
			PageSetup: &shared.StudioPageSetup{
				PageSize: wskaznik("A4"), Orientation: &orientacja,
				MarginTop: wskaznik(20), MarginBottom: wskaznik(20),
				PageNumbers: wskaznik(true),
			},
		}, &zapisany)

	var wykaz shared.StudioExportProfileListResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandStudioExportProfileList,
		shared.StudioExportProfileListRequest{}, &wykaz)
	if len(wykaz.Profiles) != 1 {
		t.Fatalf("katalog niesie %d profili, oczekiwano jednego", len(wykaz.Profiles))
	}
	profil := wykaz.Profiles[0]
	if profil.PageSetup == nil {
		t.Fatal("profil odczytany nie niesie ustawień strony — zapisano je i zniknęły")
	}
	if profil.PageSetup.PageSize == nil || *profil.PageSetup.PageSize != "A4" {
		t.Errorf("profil niesie rozmiar strony %v, zapisano A4", profil.PageSetup.PageSize)
	}
	if profil.PageSetup.MarginTop == nil || *profil.PageSetup.MarginTop != 20 {
		t.Errorf("profil niesie margines górny %v, zapisano 20", profil.PageSetup.MarginTop)
	}
	if profil.PageSetup.PageNumbers == nil || !*profil.PageSetup.PageNumbers {
		t.Error("profil zgubił numerowanie stron")
	}
}

// TestEtykietaWersjiNadajeSieIZdejmuje sprawdza obie strony jednej komendy:
// nadanie i zdjęcie. Etykieta pusta ma zdejmować, a nie zapisywać pustki.
func TestEtykietaWersjiNadajeSieIZdejmuje(t *testing.T) {
	zmontowany, zycie, _ := zmontujDoPomiaruSkutku(t)
	dokument := dokumentZTrescia(t, zmontowany, zycie, "okno-1", "Treść pierwsza.")
	if dokument.VersionId == nil {
		t.Fatal("dokument nie ma wersji, więc nie ma czego etykietować")
	}

	var nadana shared.StudioVersionLabelSetResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandStudioVersionLabelSet,
		shared.StudioVersionLabelSetRequest{
			VersionId: *dokument.VersionId,
			Label:     wskaznik("do akceptacji klienta"), Milestone: wskaznik(true),
		}, &nadana)
	if nadana.Version.Label == nil || *nadana.Version.Label != "do akceptacji klienta" {
		t.Errorf("wersja niesie etykietę %v, nadano „do akceptacji klienta”", nadana.Version.Label)
	}
	if nadana.Version.Milestone == nil || !*nadana.Version.Milestone {
		t.Error("wersja nie została oznaczona jako kluczowa")
	}

	// Etykieta ma być widoczna w historii, a nie tylko w odpowiedzi.
	var historia shared.StudioRepositoryListResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandStudioRepositoryList,
		shared.StudioRepositoryListRequest{DocumentId: dokument.Id}, &historia)
	if len(historia.Versions) == 0 || historia.Versions[0].Label == nil {
		t.Fatal("etykieta nie pojawiła się w historii wersji")
	}

	var zdjeta shared.StudioVersionLabelSetResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandStudioVersionLabelSet,
		shared.StudioVersionLabelSetRequest{VersionId: *dokument.VersionId, Label: wskaznik("")},
		&zdjeta)
	if zdjeta.Version.Label != nil {
		t.Errorf("etykieta pusta nie zdjęła etykiety: %v", zdjeta.Version.Label)
	}

	odpowiedz := wykonajKomende(t, zmontowany, zycie, shared.CommandStudioVersionLabelSet,
		shared.StudioVersionLabelSetRequest{VersionId: "studio-wer-nie-ma"})
	if odpowiedz.Error == nil || odpowiedz.Error.Code != shared.ErrorCodeNotFound {
		t.Errorf("etykietowanie wersji nieistniejącej nie odmówiło kodem not_found: %+v", odpowiedz.Error)
	}
}

// TestFormatDokumentuPrzestawiaSieBezZamianyTresci pilnuje granicy komendy:
// przestawia cechę, a zamiany treści nie udaje.
func TestFormatDokumentuPrzestawiaSieBezZamianyTresci(t *testing.T) {
	zmontowany, zycie, _ := zmontujDoPomiaruSkutku(t)
	tresc := "# Nagłówek\n\nTreść dokumentu."
	dokument := dokumentZTrescia(t, zmontowany, zycie, "okno-1", tresc)

	var przestawiony shared.StudioDocumentFormatSetResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandStudioDocumentFormatSet,
		shared.StudioDocumentFormatSetRequest{
			DocumentId: dokument.Id, Format: shared.StudioDocumentFormatMarkdown,
		}, &przestawiony)
	if przestawiony.Document.Format != shared.StudioDocumentFormatMarkdown {
		t.Errorf("dokument niesie format %q, przestawiono na markdown", przestawiony.Document.Format)
	}
	if trescDokumentu(t, zmontowany, zycie, "okno-1", dokument.Id) != tresc {
		t.Error("przestawienie formatu zmieniło treść dokumentu")
	}

	// Zamiana treści nie należy do tej komendy i ma to powiedzieć wprost,
	// zamiast po cichu jej nie wykonać.
	odpowiedz := wykonajKomende(t, zmontowany, zycie, shared.CommandStudioDocumentFormatSet,
		shared.StudioDocumentFormatSetRequest{
			DocumentId: dokument.Id, Format: shared.StudioDocumentFormatPdf,
			ConvertContent: wskaznik(true),
		})
	if odpowiedz.Error == nil {
		t.Fatal("żądanie zamiany treści wróciło odpowiedzią udaną, choć zamiany nie wykonano")
	}
	if !strings.Contains(odpowiedz.Error.Message, "obszaru dokumentów") {
		t.Errorf("odmowa nie wskazuje właściwej drogi: %q", odpowiedz.Error.Message)
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
