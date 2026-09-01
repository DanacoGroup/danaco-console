// Odpowiedzialność pliku: moduł Studio — otwieranie i zapis dokumentu Studio
// Editor. Operacje kontekstowe Tools Panel i porównania Diff/Grep Panel stoją
// w osobnym pliku adaptera. Jedyną drogą zapisu treści dokumentu jest
// `document.save`.
package core

import (
	"context"
	"errors"

	"danacoconsole/server/internal/dane"
	"danacoconsole/server/internal/konfig"
	"danacoconsole/server/internal/models"
	"danacoconsole/server/internal/protocol"
	"danacoconsole/server/internal/session"
	"danacoconsole/shared"
)

// Przedrostki identyfikatorów bytów modułu. Byt nadany przez rdzeń wychodzi
// kontraktem pod własnym identyfikatorem, nie pod kluczem wiersza (wzór:
// `adapter_modul_automations.go`).
const (
	przedrostekDokumentuStudio = "studio-dok-"
	przedrostekWersjiStudio    = "studio-wer-"
)

// adapterStudia wypełnia część portu Studio dotyczącą dokumentu i wersji.
// Drugą część (operacje kontekstowe, porównania) dokładają metody z sąsiednich
// plików na tym samym typie; typ i konstruktor deklaruje wyłącznie ten plik.
type adapterStudia struct {
	repozytorium dane.RepozytoriumStudia
	// kanaly i okna dają Studiu silnik modelu, bez którego operacja odmawia.
	kanaly *models.Rejestr
	okna   *session.Rejestr
	// uruchamiacz, rozstrzygacz i katalog dają dostęp do programów zewnętrznych.
	uruchamiacz  session.Uruchamiacz
	rozstrzygacz *konfig.Rozstrzygacz
	katalog      *KatalogRoboczy
	// zasoby i magazyn dają Studiu jedyną drogę do bajtów wyników.
	zasoby  dane.RepozytoriumDesignu
	magazyn *magazynTresciBiblioteki
	// biblioteka służy porównaniu dokumentu z materiałem `studio.diff.source`.
	biblioteka dane.RepozytoriumBiblioteki
}

// nowyAdapterStudia zakłada adapter portu Studio i wiąże go z repozytorium
// modułu, jedyną zależnością wymaganą konstruktorem.
func nowyAdapterStudia(repozytorium dane.RepozytoriumStudia) *adapterStudia {
	return &adapterStudia{repozytorium: repozytorium}
}

// ZKanalami podaje Studiu rejestr kanałów modelu i rejestr okien, ten sam,
// którym jedzie okno rozmowy i moduł Roundtable.
func (a *adapterStudia) ZKanalami(kanaly *models.Rejestr, okna *session.Rejestr) *adapterStudia {
	a.kanaly, a.okna = kanaly, okna
	return a
}

// ZNarzedziami podaje Studiu komplet potrzebny programowi zewnętrznemu:
// uruchamiacz procesów, rozstrzygacz zasięgu izolacji i katalog roboczy.
func (a *adapterStudia) ZNarzedziami(uruchamiacz session.Uruchamiacz,
	rozstrzygacz *konfig.Rozstrzygacz, katalog *KatalogRoboczy) *adapterStudia {

	a.uruchamiacz, a.rozstrzygacz, a.katalog = uruchamiacz, rozstrzygacz, katalog
	return a
}

// ZZasobami podaje Studiu magazyn bajtów wyników i repozytorium zasobów,
// nad tym samym katalogiem danych, nad którym stoi magazyn zasobów Designu.
func (a *adapterStudia) ZZasobami(zasoby dane.RepozytoriumDesignu, katalogDanych string) *adapterStudia {
	a.zasoby = zasoby
	if magazyn := magazynZasobowDesignu(katalogDanych); magazyn != nil {
		a.magazyn = magazyn
	}
	return a
}

// ZBiblioteka podaje Studiu repozytorium Library — wyłącznie do odczytu
// materiału wejściowego w `studio.diff.source`.
func (a *adapterStudia) ZBiblioteka(biblioteka dane.RepozytoriumBiblioteki) *adapterStudia {
	a.biblioteka = biblioteka
	return a
}

// OtworzDokument wczytuje dokument istniejący (po `documentId`), a bez wskazania
// wraca do ostatnio zmienianego dokumentu okna — wznowienie sesji nie ma skąd
// znać kodu, a zakładanie nowego przy każdym wejściu mnożyło puste dokumenty.
// Nowy dokument powstaje, gdy okno nie ma żadnego albo wskazano plik biblioteki.
func (a *adapterStudia) OtworzDokument(ctx context.Context,
	z shared.StudioDocumentOpenRequest) (shared.StudioDocumentOpenResponse, error) {

	if z.WindowId == "" {
		return shared.StudioDocumentOpenResponse{}, bladWskazaniaStudio("otwarcie dokumentu bez okna")
	}

	if kod := wartoscTekstu(z.DocumentId); kod != "" {
		wiersz, err := a.repozytorium.Dokument(ctx, kod)
		if err != nil {
			return shared.StudioDocumentOpenResponse{}, bladNieznanegoDokumentu(kod, err)
		}
		return shared.StudioDocumentOpenResponse{Document: a.zlozDokument(wiersz)}, nil
	}

	if wartoscTekstu(z.LibraryFileId) == "" {
		if ostatni, jest := a.ostatniDokumentOkna(ctx, z.WindowId); jest {
			return shared.StudioDocumentOpenResponse{Document: a.zlozDokument(ostatni)}, nil
		}
	}

	nowy := dane.DokumentStudia{
		Kod:                nowyIdentyfikator(przedrostekDokumentuStudio),
		Okno:               z.WindowId,
		Format:             shared.StudioDocumentFormatTxt,
		PlikRepozytoriumID: z.LibraryFileId,
	}
	zapisany, err := a.repozytorium.ZapiszDokument(ctx, nowy)
	if err != nil {
		return shared.StudioDocumentOpenResponse{}, bladStudio(err)
	}
	return shared.StudioDocumentOpenResponse{Document: a.zlozDokument(zapisany)}, nil
}

// ostatniDokumentOkna oddaje dokument okna o najpóźniejszej zmianie; brak
// dokumentów albo błąd odczytu oddaje false — wtedy powstaje nowy.
func (a *adapterStudia) ostatniDokumentOkna(ctx context.Context, okno string) (dane.DokumentStudia, bool) {
	dokumenty, err := a.repozytorium.Dokumenty(ctx, okno)
	if err != nil || len(dokumenty) == 0 {
		return dane.DokumentStudia{}, false
	}
	ostatni := dokumenty[0]
	for _, dokument := range dokumenty[1:] {
		if dokument.Zaktualizowano > ostatni.Zaktualizowano {
			ostatni = dokument
		}
	}
	return ostatni, true
}

// ZapiszDokument zapisuje treść dokumentu i, gdy Operator o to poprosi,
// zakłada wersję w repozytorium sesji, przestawiając wskaźnik wersji bieżącej.
func (a *adapterStudia) ZapiszDokument(ctx context.Context,
	z shared.StudioDocumentSaveRequest) (shared.StudioDocumentSaveResponse, error) {

	if z.DocumentId == "" {
		return shared.StudioDocumentSaveResponse{}, bladWskazaniaStudio("zapis dokumentu bez wskazania dokumentu")
	}
	istniejacy, err := a.repozytorium.Dokument(ctx, z.DocumentId)
	if err != nil {
		return shared.StudioDocumentSaveResponse{}, bladNieznanegoDokumentu(z.DocumentId, err)
	}

	// Postać z pola `form` jest utrwalana przed zapisem treści.
	if len(z.Form) > 0 {
		zPostacia, err := a.wejsciePrzyjmijPostacZapisu(ctx, istniejacy, z.Form)
		if err != nil {
			return shared.StudioDocumentSaveResponse{}, err
		}
		istniejacy = zPostacia
	}

	tresc := z.Content
	istniejacy.Tresc = &tresc
	if z.Title != nil {
		istniejacy.Tytul = z.Title
	}

	zapisany, err := a.repozytorium.ZapiszDokument(ctx, istniejacy)
	if err != nil {
		return shared.StudioDocumentSaveResponse{}, bladStudio(err)
	}

	var wersjaOdpowiedzi *shared.StudioVersion
	if z.CreateVersion != nil && *z.CreateVersion {
		nowaWersja, err := a.repozytorium.ZapiszWersje(ctx, zapisany.ID, dane.WersjaDokumentu{
			Kod:   nowyIdentyfikator(przedrostekWersjiStudio),
			Tresc: &tresc,
		})
		if err != nil {
			return shared.StudioDocumentSaveResponse{}, bladStudio(err)
		}
		zapisany.WersjaBiezacaKod = &nowaWersja.Kod
		zapisany, err = a.repozytorium.ZapiszDokument(ctx, zapisany)
		if err != nil {
			return shared.StudioDocumentSaveResponse{}, bladStudio(err)
		}
		zlozona := a.zlozWersje(nowaWersja)
		wersjaOdpowiedzi = &zlozona
	}

	return shared.StudioDocumentSaveResponse{
		Document: a.zlozDokument(zapisany),
		Version:  wersjaOdpowiedzi,
	}, nil
}

// zlozDokument składa dokument kontraktu z wiersza warstwy danych. Znaczniki
// czasu warstwa danych oddaje jako surowy tekst bazy — przekład na
// milisekundy epoki robi `chwilaBazy`, tak jak w module Automations.
func (a *adapterStudia) zlozDokument(wiersz dane.DokumentStudia) shared.StudioDocument {
	return shared.StudioDocument{
		Id:            wiersz.Kod,
		WindowId:      wiersz.Okno,
		Title:         wiersz.Tytul,
		Format:        wiersz.Format,
		Content:       wiersz.Tresc,
		LibraryFileId: wiersz.PlikRepozytoriumID,
		VersionId:     wiersz.WersjaBiezacaKod,
		CreatedAt:     chwilaBazy(wiersz.Utworzono),
		UpdatedAt:     chwilaBazy(wiersz.Zaktualizowano),
	}
}

// zlozWersje składa wersję kontraktu z wiersza warstwy danych — używane też
// przez `studio_wersje.go` (ta sama metoda, jeden sposób przekładu).
func (a *adapterStudia) zlozWersje(wiersz dane.WersjaDokumentu) shared.StudioVersion {
	wersja := shared.StudioVersion{
		Id:          wiersz.Kod,
		DocumentId:  wiersz.DokumentKod,
		Label:       wiersz.Etykieta,
		Summary:     wiersz.Podsumowanie,
		ContentHash: wiersz.SkrotTresci,
		CreatedAt:   chwilaBazy(wiersz.Utworzono),
		Milestone:   wskaznikLogiczny(wiersz.KamienMilowy),
		BranchId:    wiersz.GalazKod,
		ProposalId:  wiersz.PropozycjaKod,
	}
	// Autor wychodzi kontraktem wyłącznie wtedy, gdy wiersz go niesie.
	if wiersz.Autor != nil && *wiersz.Autor != "" {
		autor := shared.StudioAuthor(*wiersz.Autor)
		wersja.Author = &autor
	}
	return wersja
}

// bladStudio znakuje usterkę wewnętrzną modułu kodem kontraktu, wzorem
// funkcji `bladAutomatyki` z modułu Automations.
func bladStudio(err error) error {
	if err == nil {
		return nil
	}
	return protocol.JakoError(protocol.BladZeZrodla(shared.ErrorCodeInternalError, err))
}

// bladWskazaniaStudio nazywa brak danych wymaganych w żądaniu — jest to
// błąd żądania Operatora, a nie usterka rdzenia platformy.
func bladWskazaniaStudio(powod string) error {
	return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeValidationFailed,
		"moduł Studio: "+powod))
}

// bladNieznanegoDokumentu odróżnia stan „dokumentu nie ma” od stanu
// „odczyt dokumentu się nie powiódł”, niosącego przyczynę usterki.
func bladNieznanegoDokumentu(kod string, err error) error {
	if errors.Is(err, dane.ErrBrakWiersza) {
		return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeNotFound,
			"moduł Studio: dokument nie istnieje: "+kod))
	}
	return bladStudio(err)
}
