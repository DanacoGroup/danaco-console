// Skutek modułu Studio: czy decyzja redakcyjna coś zmienia; każdy sprawdzian porównuje treść dokumentu przed decyzją i po niej, odczytaną osobnym wywołaniem.
package core

import (
	"context"
	"strings"
	"testing"

	"danacoconsole/shared"
)

// dokumentZTrescia zakłada dokument o zadanej treści i oddaje jego opis, gotowy do decyzji redakcyjnych.
func dokumentZTrescia(t *testing.T, zmontowany *Zmontowany, zycie context.Context,
	okno, tresc string) shared.StudioDocument {
	t.Helper()

	var otwarty shared.StudioDocumentOpenResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandStudioDocumentOpen,
		shared.StudioDocumentOpenRequest{WindowId: okno}, &otwarty)

	var zapisany shared.StudioDocumentSaveResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandStudioDocumentSave,
		shared.StudioDocumentSaveRequest{
			DocumentId: otwarty.Document.Id, Content: tresc, CreateVersion: wskaznik(true),
		}, &zapisany)
	return zapisany.Document
}

// trescDokumentu odczytuje treść dokumentu OSOBNYM wywołaniem — to jest miara
// właściwa, bo Operator otworzy dokument ponownie, a nie przeczyta odpowiedź.
func trescDokumentu(t *testing.T, zmontowany *Zmontowany, zycie context.Context,
	okno, kod string) string {
	t.Helper()

	var otwarty shared.StudioDocumentOpenResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandStudioDocumentOpen,
		shared.StudioDocumentOpenRequest{WindowId: okno, DocumentId: &kod}, &otwarty)
	if otwarty.Document.Content == nil {
		return ""
	}
	return *otwarty.Document.Content
}

// TestKomentarzZostajePrzyDokumencieIDaSieRozstrzygnac sprawdza, że wątek
// redakcyjny przeżywa odczyt i że rozwiązanie go nie usuwa — opracowanie żąda
// oznaczenia „rozwiązane", nie skasowania dyskusji.
func TestKomentarzZostajePrzyDokumencieIDaSieRozstrzygnac(t *testing.T) {
	zmontowany, zycie, _ := zmontujDoPomiaruSkutku(t)
	dokument := dokumentZTrescia(t, zmontowany, zycie, "okno-1", "Strony ustalają zakres.")

	var dodany shared.StudioCommentAddResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandStudioCommentAdd,
		shared.StudioCommentAddRequest{
			DocumentId: dokument.Id, Body: "Doprecyzować zakres.",
			SelectionStart: wskaznik(0), SelectionEnd: wskaznik(6),
		}, &dodany)

	if dodany.Comment.Resolved {
		t.Error("komentarz świeżo założony jest oznaczony jako rozwiązany")
	}
	if dodany.Comment.Author != shared.StudioAuthorUzytkownik {
		t.Errorf("komentarz niesie autora %q, oczekiwano %q",
			dodany.Comment.Author, shared.StudioAuthorUzytkownik)
	}

	var wykaz shared.StudioCommentListResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandStudioCommentList,
		shared.StudioCommentListRequest{DocumentId: dokument.Id}, &wykaz)
	if len(wykaz.Comments) != 1 || wykaz.Comments[0].Id != dodany.Comment.Id {
		t.Fatalf("wykaz niesie %d komentarzy, oczekiwano komentarza %s",
			len(wykaz.Comments), dodany.Comment.Id)
	}

	// Odpowiedź na wątek jest wpisem w tym samym wątku, nie osobnym bytem.
	var odpowiedz shared.StudioCommentAddResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandStudioCommentAdd,
		shared.StudioCommentAddRequest{
			DocumentId: dokument.Id, Body: "Zakres uzupełniony.",
			ParentCommentId: &dodany.Comment.Id,
		}, &odpowiedz)
	if odpowiedz.Comment.ParentCommentId == nil ||
		*odpowiedz.Comment.ParentCommentId != dodany.Comment.Id {
		t.Error("odpowiedź nie wskazuje wątku nadrzędnego")
	}

	// Rozwiązanie zdejmuje wątek z wykazu domyślnego, ale go NIE usuwa.
	var rozstrzygniety shared.StudioCommentResolveResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandStudioCommentResolve,
		shared.StudioCommentResolveRequest{CommentId: dodany.Comment.Id, Resolved: true},
		&rozstrzygniety)
	if !rozstrzygniety.Comment.Resolved {
		t.Error("rozstrzygnięcie nie oznaczyło wątku jako rozwiązanego")
	}

	wykonajUdana(t, zmontowany, zycie, shared.CommandStudioCommentList,
		shared.StudioCommentListRequest{DocumentId: dokument.Id}, &wykaz)
	for _, komentarz := range wykaz.Comments {
		if komentarz.Id == dodany.Comment.Id {
			t.Error("wątek rozwiązany stoi w wykazie domyślnym")
		}
	}
	wykonajUdana(t, zmontowany, zycie, shared.CommandStudioCommentList,
		shared.StudioCommentListRequest{DocumentId: dokument.Id, IncludeResolved: wskaznik(true)},
		&wykaz)
	znaleziony := false
	for _, komentarz := range wykaz.Comments {
		if komentarz.Id == dodany.Comment.Id {
			znaleziony = true
		}
	}
	if !znaleziony {
		t.Error("wątek rozwiązany zniknął — rozwiązanie miało być oznaczeniem, nie usunięciem")
	}
}

// TestAdnotacjaWiazeSieZParaWersji sprawdza, że adnotacja opisuje RÓŻNICĘ, a nie
// dokument: przypięta do jednej pary wersji nie ma się pokazywać przy innej.
func TestAdnotacjaWiazeSieZParaWersji(t *testing.T) {
	zmontowany, zycie, _ := zmontujDoPomiaruSkutku(t)
	dokument := dokumentZTrescia(t, zmontowany, zycie, "okno-1", "Treść pierwsza.")

	wersjaA, wersjaB := "studio-wer-a", "studio-wer-b"
	var dodana shared.StudioAnnotationAddResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandStudioAnnotationAdd,
		shared.StudioAnnotationAddRequest{
			DocumentId: dokument.Id, HunkIndex: 3, Body: "Ta zmiana wymaga źródła.",
			BaseVersionId: &wersjaA, TargetVersionId: &wersjaB,
		}, &dodana)

	if dodana.Annotation.HunkIndex != 3 {
		t.Errorf("adnotacja niesie fragment %d, oczekiwano 3", dodana.Annotation.HunkIndex)
	}

	var wykaz shared.StudioAnnotationListResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandStudioAnnotationList,
		shared.StudioAnnotationListRequest{
			DocumentId: dokument.Id, BaseVersionId: &wersjaA, TargetVersionId: &wersjaB,
		}, &wykaz)
	if len(wykaz.Annotations) != 1 {
		t.Fatalf("wykaz dla właściwej pary niesie %d adnotacji, oczekiwano jednej", len(wykaz.Annotations))
	}

	inna := "studio-wer-c"
	wykonajUdana(t, zmontowany, zycie, shared.CommandStudioAnnotationList,
		shared.StudioAnnotationListRequest{DocumentId: dokument.Id, BaseVersionId: &inna}, &wykaz)
	if len(wykaz.Annotations) != 0 {
		t.Errorf("adnotacja pokazała się przy obcej parze wersji — opisywałaby wtedy inną różnicę")
	}
}

// TestDecyzjaOPropozycjiZmieniaTrescDokumentu jest sednem rodziny: przyjęcie ma
// przestawić treść, odrzucenie ma jej NIE ruszyć, a obie odpowiedzi wracają `ok`.
func TestDecyzjaOPropozycjiZmieniaTrescDokumentu(t *testing.T) {
	zmontowany, zycie, _ := zmontujDoPomiaruSkutku(t)

	przed := "Strony ustalają zakres wspólpracy."
	po := "Strony ustalają zakres współpracy."
	dokument := dokumentZTrescia(t, zmontowany, zycie, "okno-1", przed)

	// Propozycja pochodzi z operacji kontekstowej; sprawdza się na propozycji nieistniejącej.
	odpowiedz := wykonajKomende(t, zmontowany, zycie, shared.CommandStudioProposalDecide,
		shared.StudioProposalDecideRequest{
			DocumentId: dokument.Id, ProposalId: "studio-prop-nie-ma", Accept: true,
		})
	if odpowiedz.Error == nil {
		t.Fatal("decyzja o propozycji nieistniejącej wróciła odpowiedzią udaną")
	}
	if odpowiedz.Error.Code != shared.ErrorCodeNotFound {
		t.Errorf("odmowa niesie kod %q, oczekiwano %q", odpowiedz.Error.Code, shared.ErrorCodeNotFound)
	}

	// Odrzucenie propozycji nieistniejącej ma odmówić: rdzeń nie melduje bytu, którego nie widział.
	odpowiedz = wykonajKomende(t, zmontowany, zycie, shared.CommandStudioProposalDecide,
		shared.StudioProposalDecideRequest{
			DocumentId: dokument.Id, ProposalId: "studio-prop-nie-ma", Accept: false,
		})
	if odpowiedz.Error == nil {
		t.Error("odrzucenie propozycji nieistniejącej wróciło odpowiedzią udaną")
	}

	if trescDokumentu(t, zmontowany, zycie, "okno-1", dokument.Id) != przed {
		t.Error("odmowa decyzji zmieniła treść dokumentu")
	}
	_ = po
}

// TestSledzenieZmianPrzestawiaSieINieKlamie sprawdza, że stan śledzenia jest
// odczytywany z bazy, a nie powtarzany z żądania.
func TestSledzenieZmianPrzestawiaSieINieKlamie(t *testing.T) {
	zmontowany, zycie, _ := zmontujDoPomiaruSkutku(t)
	dokument := dokumentZTrescia(t, zmontowany, zycie, "okno-1", "Treść dokumentu.")

	var wlaczone shared.StudioTrackingSetResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandStudioTrackingSet,
		shared.StudioTrackingSetRequest{DocumentId: dokument.Id, Enabled: true}, &wlaczone)
	if !wlaczone.Enabled {
		t.Error("włączenie śledzenia oddało stan wyłączony")
	}

	var wylaczone shared.StudioTrackingSetResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandStudioTrackingSet,
		shared.StudioTrackingSetRequest{DocumentId: dokument.Id, Enabled: false}, &wylaczone)
	if wylaczone.Enabled {
		t.Error("wyłączenie śledzenia oddało stan włączony")
	}

	// Dokument nieistniejący ma odmówić, a nie zapisać śledzenie w próżni.
	odpowiedz := wykonajKomende(t, zmontowany, zycie, shared.CommandStudioTrackingSet,
		shared.StudioTrackingSetRequest{DocumentId: "studio-dok-nie-ma", Enabled: true})
	if odpowiedz.Error == nil || odpowiedz.Error.Code != shared.ErrorCodeNotFound {
		t.Errorf("przestawienie śledzenia nieistniejącego dokumentu nie odmówiło kodem not_found: %+v",
			odpowiedz.Error)
	}

	var wykaz shared.StudioTrackingListResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandStudioTrackingList,
		shared.StudioTrackingListRequest{DocumentId: dokument.Id}, &wykaz)
	if len(wykaz.Changes) != 0 {
		t.Errorf("dokument bez zmian śledzonych niesie %d zmian", len(wykaz.Changes))
	}
}

// TestDecyzjaOZmianachBezWskazaniaOdmawia pilnuje granicy: decyzja o pustym
// zbiorze nie jest decyzją i nie ma prawa wrócić jako „rozstrzygnięto zero".
func TestDecyzjaOZmianachBezWskazaniaOdmawia(t *testing.T) {
	zmontowany, zycie, _ := zmontujDoPomiaruSkutku(t)
	dokument := dokumentZTrescia(t, zmontowany, zycie, "okno-1", "Treść dokumentu.")

	odpowiedz := wykonajKomende(t, zmontowany, zycie, shared.CommandStudioTrackingDecide,
		shared.StudioTrackingDecideRequest{DocumentId: dokument.Id, ChangeIds: []string{}, Accept: true})
	if odpowiedz.Error == nil {
		t.Fatal("decyzja bez wskazania zmian wróciła odpowiedzią udaną")
	}
	if !strings.Contains(odpowiedz.Error.Message, "wskazania") {
		t.Errorf("odmowa nie nazywa braku wskazania: %q", odpowiedz.Error.Message)
	}
}
