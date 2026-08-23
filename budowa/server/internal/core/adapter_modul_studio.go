// Odpowiedzialność pliku: moduł Studio — otwieranie i zapis dokumentu Studio
// Editor. Operacje kontekstowe Tools Panel i porównania Diff/Grep Panel stoją
// w osobnym pliku adaptera.
//
// Jedna droga zapisu: `document.save` jest jedynym miejscem, które zapisuje
// treść dokumentu — `document.open` tylko czyta albo zakłada wiersz pusty.
// Dwie ścieżki zapisu tej samej treści byłyby dwiema prawdami o tym samym
// bycie.
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
	// kanaly i okna dają Studiu silnik modelu — ten sam rejestr, którym jedzie
	// okno rozmowy i moduł Roundtable. Bez nich operacja kontekstowa odmawia
	// wprost zamiast udawać pracę.
	kanaly *models.Rejestr
	okna   *session.Rejestr
	// uruchamiacz, rozstrzygacz i katalog dają Studiu dostęp do arsenału
	// programów zewnętrznych — rozpoznania pisma, zamiany formatów, pakowania.
	// Ten sam komplet trzech zależności, którym jedzie adapter narzędzi
	// dokumentu; rodzina cyfryzacji bez niego nie ma czym wystartować procesu.
	uruchamiacz  session.Uruchamiacz
	rozstrzygacz *konfig.Rozstrzygacz
	katalog      *KatalogRoboczy
	// zasoby i magazyn dają Studiu jedyną drogę do bajtów: magazyn trzyma
	// treść pod sumą kontrolną, repozytorium Designu — wiersz, po którym
	// Assets Panel i Preview Window widzą wynik. Bez tej pary wydanie archiwum,
	// wyrys strony i osadzenie grafiki nie mają gdzie odłożyć tego, co zrobiły,
	// a odczyt zasobu nie ma skąd wziąć materiału.
	zasoby  dane.RepozytoriumDesignu
	magazyn *magazynTresciBiblioteki
	// biblioteka służy jednej czynności: porównaniu dokumentu roboczego
	// z materiałem wejściowym (`studio.diff.source`). Studio nie zarządza
	// biblioteką i niczego w niej nie zapisuje — czyta odwołanie do treści
	// pliku i tyle.
	biblioteka dane.RepozytoriumBiblioteki
}

// nowyAdapterStudia wiąże port z repozytorium modułu.
func nowyAdapterStudia(repozytorium dane.RepozytoriumStudia) *adapterStudia {
	return &adapterStudia{repozytorium: repozytorium}
}

// ZKanalami podaje Studiu rejestr kanałów modelu i rejestr okien.
//
// Bez tych rejestrów operacja kontekstowa — sedno modułu Studio — odmawia
// kodem `channel_unavailable`, bo nie ma czym wywołać modelu. Podaje się tu
// ten sam rejestr, którym jedzie okno rozmowy i moduł Roundtable.
//
// Rejestr okien jest potrzebny, bo kanał modelu należy do okna: Studio pracuje
// w imieniu okna komunikacji i ma sięgnąć po ten kanał, który Operator ustawił
// temu oknu, a nie po żaden własny.
func (a *adapterStudia) ZKanalami(kanaly *models.Rejestr, okna *session.Rejestr) *adapterStudia {
	a.kanaly, a.okna = kanaly, okna
	return a
}

// ZNarzedziami podaje Studiu komplet, bez którego nie da się uruchomić programu
// zewnętrznego: uruchamiacz procesów, rozstrzygacz zasięgu izolacji i ustalacz
// katalogu roboczego.
//
// Trzy zależności, nie jedna, bo każda odpowiada za co innego i żadnej nie da
// się wyprowadzić z pozostałych: uruchamiacz startuje proces, rozstrzygacz mówi,
// jakie zasady obowiązują okno, a katalog wskazuje obszar, w którym proces wolno
// puścić. Ten sam komplet bierze adapter narzędzi dokumentu (`ZIzolacja`) —
// nazwa jest tu inna, bo Studio bierze wszystkie trzy naraz, a tamten adapter
// dostaje uruchamiacz konstruktorem.
//
// Zależność jest opcjonalna na tych samych zasadach co rejestr kanałów: bez
// niej czynności sięgające po arsenał odmawiają zdaniem nazywającym brak,
// a reszta modułu pracuje dalej.
func (a *adapterStudia) ZNarzedziami(uruchamiacz session.Uruchamiacz,
	rozstrzygacz *konfig.Rozstrzygacz, katalog *KatalogRoboczy) *adapterStudia {

	a.uruchamiacz, a.rozstrzygacz, a.katalog = uruchamiacz, rozstrzygacz, katalog
	return a
}

// ZZasobami podaje Studiu magazyn bajtów wyników i repozytorium zasobów.
//
// Katalog danych jest ten sam, nad którym stoi magazyn zasobów Designu —
// wynik wydania Studia i wynik warsztatu PDF mają leżeć w jednym miejscu, bo
// jedno i drugie jest zasobem tej samej platformy, a dwa magazyny znaczyłyby
// dwa katalogi, z których jeden prędzej czy później zostałby przy kopii.
//
// Zależność jest opcjonalna na tych samych zasadach co rejestr kanałów: bez
// niej czynności wydania, wyrysu i osadzenia odmawiają zdaniem nazywającym
// brak, a reszta modułu pracuje dalej.
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

// OtworzDokument wczytuje dokument istniejący (po `documentId`) albo zakłada
// nowy. Wskazanie pliku repozytorium zapamiętuje odwołanie, ale treści z
// Library nie doczytuje — to zrobiłoby z adaptera Studio klienta modułu
// Library, którego konstruktor nie zna (jedna zależność: RepozytoriumStudia).
// Ścieżka urządzenia (`path`) jest lokalna dla klienta; rdzeń nie ma dostępu
// do systemu plików Operatora, więc dokument otwiera się bez treści, a
// pierwszy `document.save` ją dostarcza.
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

// ZapiszDokument zapisuje treść dokumentu i, gdy Operator o to poprosi,
// zakłada wersję w repozytorium sesji.
//
// Wskaźnik wersji bieżącej przestawia ten adapter, nie warstwa danych:
// `dane.ZapiszWersje` zostawia `wersja_biezaca_id` nietknięty (patrz komentarz
// w `studio_wersje.go`) i oddaje wywołującemu decyzję, kiedy nowa wersja staje
// się bieżącą. Tutaj `document.save` z `createVersion=true` zakłada wersję po
// to, żeby Operator dalej edytował od niej, więc od razu staje się bieżącą
// (drugie wywołanie `ZapiszDokument` z nowym `WersjaBiezacaKod`). Inaczej niż
// w `PrzywrocWersje` (`studio_wersje.go`), gdzie bieżącą staje się wersja
// wskazana przez Operatora, a nie najświeższa zapisana.
func (a *adapterStudia) ZapiszDokument(ctx context.Context,
	z shared.StudioDocumentSaveRequest) (shared.StudioDocumentSaveResponse, error) {

	if z.DocumentId == "" {
		return shared.StudioDocumentSaveResponse{}, bladWskazaniaStudio("zapis dokumentu bez wskazania dokumentu")
	}
	istniejacy, err := a.repozytorium.Dokument(ctx, z.DocumentId)
	if err != nil {
		return shared.StudioDocumentSaveResponse{}, bladNieznanegoDokumentu(z.DocumentId, err)
	}

	// Postać dokumentu przyjeżdża polem `form` i jest utrwalana PRZED zapisem
	// treści — utrwalenie postaci składa treść z bloków i wpisuje ją do
	// wiersza, więc treść z żądania musi wejść po nim, a nie przed. Brak pola
	// znaczy „bez zmiany postaci", nie „postać na zero": zwykły zapis treści
	// nie ma prawa zetrzeć arkusza stylów ani tabel.
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
	// Autor wychodzi kontraktem WYŁĄCZNIE wtedy, gdy wiersz go niesie. Wersje
	// założone przed dobudową autora nie mają, a podstawienie tu Operatora
	// zamieniłoby brak wiedzy w twierdzenie — i to twierdzenie fałszywe dla
	// każdej wersji, którą naprawdę zapisał model.
	if wiersz.Autor != nil && *wiersz.Autor != "" {
		autor := shared.StudioAuthor(*wiersz.Autor)
		wersja.Author = &autor
	}
	return wersja
}

// bladStudio znakuje usterkę wewnętrzną kodem kontraktu (wzór: `bladAutomatyki`).
func bladStudio(err error) error {
	if err == nil {
		return nil
	}
	return protocol.JakoError(protocol.BladZeZrodla(shared.ErrorCodeInternalError, err))
}

// bladWskazaniaStudio nazywa brak danych w żądaniu — błąd Operatora, nie rdzenia.
func bladWskazaniaStudio(powod string) error {
	return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeValidationFailed,
		"moduł Studio: "+powod))
}

// bladNieznanegoDokumentu odróżnia „dokumentu nie ma” od „odczyt się nie powiódł”.
func bladNieznanegoDokumentu(kod string, err error) error {
	if errors.Is(err, dane.ErrBrakWiersza) {
		return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeNotFound,
			"moduł Studio: dokument nie istnieje: "+kod))
	}
	return bladStudio(err)
}
