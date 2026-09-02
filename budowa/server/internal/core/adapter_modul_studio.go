// Moduł Studio: otwieranie i zapis dokumentu Studio Editor; jedyną drogą zapisu treści
// jest `document.save`. Operacje kontekstowe i porównania stoją w osobnym pliku adaptera.
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

const (
	przedrostekDokumentuStudio = "studio-dok-"
	przedrostekWersjiStudio    = "studio-wer-"
)

type adapterStudia struct {
	repozytorium dane.RepozytoriumStudia
	kanaly       *models.Rejestr
	okna         *session.Rejestr
	uruchamiacz  session.Uruchamiacz
	rozstrzygacz *konfig.Rozstrzygacz
	katalog      *KatalogRoboczy
	zasoby       dane.RepozytoriumDesignu
	magazyn      *magazynTresciBiblioteki
	biblioteka   dane.RepozytoriumBiblioteki
}

func nowyAdapterStudia(repozytorium dane.RepozytoriumStudia) *adapterStudia {
	return &adapterStudia{repozytorium: repozytorium}
}

func (a *adapterStudia) ZKanalami(kanaly *models.Rejestr, okna *session.Rejestr) *adapterStudia {
	a.kanaly, a.okna = kanaly, okna
	return a
}

func (a *adapterStudia) ZNarzedziami(uruchamiacz session.Uruchamiacz,
	rozstrzygacz *konfig.Rozstrzygacz, katalog *KatalogRoboczy) *adapterStudia {

	a.uruchamiacz, a.rozstrzygacz, a.katalog = uruchamiacz, rozstrzygacz, katalog
	return a
}

func (a *adapterStudia) ZZasobami(zasoby dane.RepozytoriumDesignu, katalogDanych string) *adapterStudia {
	a.zasoby = zasoby
	if magazyn := magazynZasobowDesignu(katalogDanych); magazyn != nil {
		a.magazyn = magazyn
	}
	return a
}

func (a *adapterStudia) ZBiblioteka(biblioteka dane.RepozytoriumBiblioteki) *adapterStudia {
	a.biblioteka = biblioteka
	return a
}

// Bez wskazania wraca ostatnio zmieniany dokument okna: wznowienie sesji nie zna kodu, a zakładanie nowego mnożyło puste dokumenty.
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

// Znaczniki czasu warstwa danych oddaje surowym tekstem bazy; przekład na milisekundy robi `chwilaBazy`.
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

func bladStudio(err error) error {
	if err == nil {
		return nil
	}
	if errors.Is(err, dane.ErrKolizjaWiersza) {
		return protocol.JakoError(protocol.BladZeZrodla(shared.ErrorCodeConflict, err))
	}
	return protocol.JakoError(protocol.BladZeZrodla(shared.ErrorCodeInternalError, err))
}

func bladWskazaniaStudio(powod string) error {
	return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeValidationFailed,
		"moduł Studio: "+powod))
}

func bladNieznanegoDokumentu(kod string, err error) error {
	if errors.Is(err, dane.ErrBrakWiersza) {
		return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeNotFound,
			"moduł Studio: dokument nie istnieje: "+kod))
	}
	return bladStudio(err)
}
