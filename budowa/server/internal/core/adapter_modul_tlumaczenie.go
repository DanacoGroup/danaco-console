// Typ adaptera modułu Translate, przedrostki identyfikatorów bytów modułu oraz komendy
// translate.source.set, translate.source.segment i translate.source.detect.
package core

import (
	"context"
	"errors"
	"strings"

	"danacoconsole/server/internal/dane"
	"danacoconsole/server/internal/konfig"
	"danacoconsole/server/internal/models"
	"danacoconsole/server/internal/protocol"
	"danacoconsole/server/internal/session"
	"danacoconsole/shared"
)

// Przedrostki całego modułu stoją tu w komplecie, żeby pozostałe pliki modułu
// nie deklarowały ich po raz drugi.
const (
	przedrostekOknaTlumaczenia   = "okt-"
	przedrostekPaneluTlumaczenia = "pan-"
	przedrostekTerminuSlownika   = "trm-"
	przedrostekImportuSlownika   = "imp-"
	przedrostekEksportuSlownika  = "eks-"
	przedrostekPamieciTlumaczen  = "pam-"
	przedrostekSyntezyMowy       = "syn-"
	przedrostekEksportuPanelu    = "epa-"
)

type adapterTlumaczenia struct {
	repozytorium dane.RepozytoriumTlumaczen
	kanaly       *models.Rejestr
	// repozytoriumKanalow rozstrzyga własność kanału z żądania (decyzja 34).
	repozytoriumKanalow dane.RepozytoriumKanalow
	wyjscie             *emiter
	biblioteka          dane.RepozytoriumBiblioteki
	magazynWytworow     *magazynTresciBiblioteki
	uruchamiacz         session.Uruchamiacz
	rozstrzygaczMowy    *konfig.Rozstrzygacz
	katalogIzolacji     *KatalogRoboczy
	katalogDanych       string
}

var (
	errPustyPrzeklad           = errors.New("model oddał pusty przekład — panel nie dostał treści")
	errPusteRozpoznanie        = errors.New("model oddał pustą odpowiedź — język nierozpoznany")
	errPusteTlumaczenieZwrotne = errors.New("model oddał puste tłumaczenie zwrotne — kontrola wierności nie ma wyniku")
)

var errBrakMagazynuWytworow = errors.New(
	"moduł Translate: serwer nie ma wpiętego magazynu wytworów ani repozytorium biblioteki — " +
		"wytwór nie miałby gdzie leżeć; naprawa: podpiąć ZWytworami przy składaniu serwera")

func (a *adapterTlumaczenia) ZWytworami(biblioteka dane.RepozytoriumBiblioteki,
	katalogDanych string) *adapterTlumaczenia {
	a.biblioteka = biblioteka
	if strings.TrimSpace(katalogDanych) != "" {
		a.magazynWytworow = nowyMagazynTresciBiblioteki(katalogDanych)
	}
	return a
}

func nowyAdapterTlumaczenia(repozytorium dane.RepozytoriumTlumaczen) *adapterTlumaczenia {
	return &adapterTlumaczenia{repozytorium: repozytorium}
}

func (a *adapterTlumaczenia) ZKanalami(kanaly *models.Rejestr,
	repozytorium dane.RepozytoriumKanalow) *adapterTlumaczenia {

	a.kanaly, a.repozytoriumKanalow = kanaly, repozytorium
	return a
}

func (a *adapterTlumaczenia) zapytajModel(ctx context.Context, okno, kanal, tresc string) (string, error) {
	if a.kanaly == nil {
		return "", bladBrakuKanalowTlumaczenia()
	}
	if strings.TrimSpace(kanal) == "" {
		return "", bladWskazaniaTlumaczenia("żądanie bez wskazania kanału modelu — serwer nie zgaduje kanału tłumaczenia")
	}
	var zebrane strings.Builder
	ujscie := models.UjscieFunkcji(func(_ context.Context, f models.Fragment) error {
		if f.Kind == shared.ChunkKindText {
			zebrane.WriteString(models.TrescFragmentu(f))
		}
		return nil
	})
	zapytanie := models.Zapytanie{
		Zasiegi:   models.Zasiegi{Okno: okno},
		Wiadomosc: okno,
		Tresc:     tresc,
		Kanal:     kanal,
	}
	if err := a.kanaly.Wyslij(ctx, zapytanie, ujscie); err != nil {
		return "", bladTlumaczenia(err)
	}
	return zebrane.String(), nil
}

func bladBrakuKanalowTlumaczenia() error {
	return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeChannelUnavailable,
		"moduł Translate: rejestr kanałów modelu nie jest wpięty — tłumaczenie nie ma czym wołać modelu"))
}

func (a *adapterTlumaczenia) UstawZrodlo(ctx context.Context,
	z shared.TranslateSourceSetRequest) (shared.TranslateSourceSetResponse, error) {

	if z.WindowId == "" {
		return shared.TranslateSourceSetResponse{}, bladWskazaniaTlumaczenia("żądanie bez okna tłumaczenia")
	}
	if z.Text == "" {
		return shared.TranslateSourceSetResponse{}, bladWskazaniaTlumaczenia("żądanie bez tekstu źródłowego")
	}

	tekst := z.Text
	okno := dane.OknoTlumaczenia{
		Kod:           z.WindowId,
		TekstZrodlowy: &tekst,
		JezykZrodlowy: z.SourceLanguage,
	}

	// Segmentacja liczy się od nowa, gdy operator o to poprosił albo gdy
	// okno jeszcze nie istnieje.
	if z.Resegment == nil || *z.Resegment {
		liczba := int64(len(podzielNaZdania(z.Text)))
		okno.LiczbaSegmentow = &liczba
	}

	zapisane, err := a.repozytorium.ZapiszOkno(ctx, okno)
	if err != nil {
		return shared.TranslateSourceSetResponse{}, bladTlumaczenia(err)
	}

	panele, err := a.repozytorium.Panele(ctx, zapisane.ID)
	if err != nil {
		return shared.TranslateSourceSetResponse{}, bladTlumaczenia(err)
	}

	return shared.TranslateSourceSetResponse{
		SourceLanguage: jezykZrodlowyZadania(zapisane.JezykZrodlowy),
		SegmentCount:   liczbaSegmentowOdpowiedzi(zapisane.LiczbaSegmentow),
		Panels:         zlozPaneleTlumaczenia(panele),
	}, nil
}

// Kontrakt każe oddać string, nie wskaźnik; brak rozpoznania zostaje pustym napisem.
func jezykZrodlowyZadania(jezyk *string) string {
	if jezyk != nil {
		return *jezyk
	}
	return ""
}

func liczbaSegmentowOdpowiedzi(liczba *int64) *int {
	if liczba == nil {
		return nil
	}
	n := int(*liczba)
	return &n
}

func (a *adapterTlumaczenia) PodzielNaSegmenty(ctx context.Context,
	z shared.TranslateSourceSegmentRequest) (shared.TranslateSourceSegmentResponse, error) {

	tekst := ""
	if z.Text != nil {
		tekst = *z.Text
	}
	if tekst == "" {
		return shared.TranslateSourceSegmentResponse{}, bladWskazaniaTlumaczenia(
			"żądanie bez tekstu do podziału — translate.source.segment nie ma dostępu do okna bez tekstu wskazanego wprost")
	}

	return shared.TranslateSourceSegmentResponse{Segments: podzielNaZdania(tekst)}, nil
}

func podzielNaZdania(tekst string) []string {
	var zdania []string
	poczatek := 0
	for i, r := range tekst {
		if r != '.' && r != '!' && r != '?' {
			continue
		}
		koniecZnaku := i + len(string(r))
		czyGraniczne := koniecZnaku >= len(tekst) || tekst[koniecZnaku] == ' ' ||
			tekst[koniecZnaku] == '\n' || tekst[koniecZnaku] == '\t'
		if !czyGraniczne {
			continue
		}
		kandydat := strings.TrimSpace(tekst[poczatek:koniecZnaku])
		if kandydat != "" {
			zdania = append(zdania, kandydat)
		}
		poczatek = koniecZnaku
	}
	ogon := strings.TrimSpace(tekst[poczatek:])
	if ogon != "" {
		zdania = append(zdania, ogon)
	}
	return zdania
}

func (a *adapterTlumaczenia) RozpoznajJezyk(ctx context.Context,
	z shared.TranslateSourceDetectRequest) (shared.TranslateSourceDetectResponse, error) {

	tekst := ""
	if z.Text != nil {
		tekst = strings.TrimSpace(*z.Text)
	}
	if tekst == "" {
		return shared.TranslateSourceDetectResponse{}, bladWskazaniaTlumaczenia(
			"żądanie translate.source.detect bez tekstu — kontrakt nie niesie okna, więc serwer nie ma czego rozpoznać")
	}

	jezyk, err := a.rozpoznajJezykModelem(ctx, tekst)
	if err != nil {
		return shared.TranslateSourceDetectResponse{}, err
	}
	return shared.TranslateSourceDetectResponse{Language: jezyk}, nil
}

func zlozPaneleTlumaczenia(panele []dane.PanelTlumaczenia) []shared.TranslationPanel {
	wynik := make([]shared.TranslationPanel, 0, len(panele))
	for _, p := range panele {
		wynik = append(wynik, zlozPanelTlumaczenia(p))
	}
	return wynik
}

func zlozPanelTlumaczenia(p dane.PanelTlumaczenia) shared.TranslationPanel {
	return shared.TranslationPanel{
		Id:            p.Kod,
		WindowId:      p.OknoKod,
		Language:      p.Jezyk,
		Text:          p.Tresc,
		Status:        shared.TranslationStatus(p.Stan),
		Tone:          p.Ton,
		UpdatedAt:     p.Zaktualizowano,
		ApprovalStage: etapZatwierdzeniaPanelu(p.EtapZatwierdzenia),
		ApprovedBy:    p.Zatwierdzil,
		ApprovedAt:    p.Zatwierdzono,
	}
}

func etapZatwierdzeniaPanelu(etap *string) *shared.ApprovalStage {
	if etap == nil || strings.TrimSpace(*etap) == "" {
		return nil
	}
	wartosc := shared.ApprovalStage(*etap)
	return &wartosc
}

func bladTlumaczenia(err error) error {
	if err == nil {
		return nil
	}
	if errors.Is(err, dane.ErrKolizjaWiersza) {
		return protocol.JakoError(protocol.BladZeZrodla(shared.ErrorCodeConflict, err))
	}
	return protocol.JakoError(protocol.BladZeZrodla(shared.ErrorCodeInternalError, err))
}

func bladWskazaniaTlumaczenia(powod string) error {
	return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeValidationFailed, "moduł Translate: "+powod))
}

func bladNieznanegoOkna(kod string, err error) error {
	if errors.Is(err, dane.ErrBrakWiersza) {
		return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeNotFound,
			"moduł Translate: nie ma okna tłumaczenia o identyfikatorze "+kod))
	}
	return bladTlumaczenia(err)
}

func bladNieznanegoPanelu(kod string, err error) error {
	if errors.Is(err, dane.ErrBrakWiersza) {
		return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeNotFound,
			"moduł Translate: nie ma panelu tłumaczenia o identyfikatorze "+kod))
	}
	return bladTlumaczenia(err)
}

func bladNieznanegoUstaleniaKorekty(kod string, err error) error {
	if errors.Is(err, dane.ErrBrakWiersza) {
		return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeNotFound,
			"moduł Translate: nie ma ustalenia korekty o identyfikatorze "+kod))
	}
	return bladTlumaczenia(err)
}
