// Odpowiedzialność pliku: wypełnienie portu Tozsamosc katalogiem kategorii
// i treścią zapisaną per oś. Składanie nakładki obowiązującej należy do
// SkladaczTozsamosci — adapter go woła, nie powtarza.
package core

import (
	"context"
	"strconv"

	"danacoconsole/server/internal/dane"
	"danacoconsole/server/internal/konfig"
	"danacoconsole/shared"
)

// Zgodność adaptera z portem sprawdzana jest przy kompilacji, bez próby wykonania kodu rdzenia platformy.
var _ Tozsamosc = (*adapterTozsamosci)(nil)

// adapterTozsamosci wypełnia port Tozsamosc tabelami kategorii tożsamości i dokumentu tożsamości modelu.
type adapterTozsamosci struct {
	repozytorium dane.RepozytoriumTozsamosci
	skladacz     *SkladaczTozsamosci
	rozstrzygacz *konfig.Rozstrzygacz
	osie         *wskazanieOsiOkna
}

// nowyAdapterTozsamosci wiąże port z repozytorium, składaczem nakładki
// i rozstrzygaczem klucza `tozsamosc.tryb_domyslny`.
func nowyAdapterTozsamosci(repozytorium dane.RepozytoriumTozsamosci,
	rozstrzygacz *konfig.Rozstrzygacz) *adapterTozsamosci {

	return &adapterTozsamosci{
		repozytorium: repozytorium,
		skladacz:     NowySkladaczTozsamosci(NoweZrodloTozsamosci(repozytorium)),
		rozstrzygacz: rozstrzygacz,
	}
}

// ZOknami wpina odczyt osi okna: model i konto, dla których liczona jest
// nakładka okna rozmowy. Bez tego wpięcia nakładkę liczy się dla osi podanej
// wprost w żądaniu.
func (a *adapterTozsamosci) ZOknami(okna dane.RepozytoriumOkien,
	kanaly dane.RepozytoriumKanalow) *adapterTozsamosci {

	a.osie = nowaWskazanieOsiOkna(okna, kanaly)
	return a
}

// Kategorie zwraca katalog kategorii zasad tożsamości, zawężony warstwą nakładki modelu w tej rozmowie.
func (a *adapterTozsamosci) Kategorie(ctx context.Context,
	z shared.IdentityCategoryListRequest) (shared.IdentityCategoryListResponse, error) {

	if a == nil || a.repozytorium == nil {
		return shared.IdentityCategoryListResponse{Categories: []shared.IdentityCategory{}}, nil
	}
	wiersze, err := a.repozytorium.Kategorie(ctx, tylkoAktywneKatalogu(z.IncludeDisabled))
	if err != nil {
		return shared.IdentityCategoryListResponse{}, err
	}
	kategorie := make([]shared.IdentityCategory, 0, len(wiersze))
	for _, wiersz := range wiersze {
		kategoria := kategoriaKontraktu(wiersz)
		if z.Layer != nil && *z.Layer != "" && kategoria.Layer != *z.Layer {
			continue
		}
		kategorie = append(kategorie, kategoria)
	}
	return shared.IdentityCategoryListResponse{Categories: kategorie}, nil
}

// Dokumenty zwraca zapisy treści tożsamości zawężone kategorią zasad i wskazaną osią modelu platformy.
func (a *adapterTozsamosci) Dokumenty(ctx context.Context,
	z shared.IdentityDocumentGetRequest) (shared.IdentityDocumentGetResponse, error) {

	if a == nil || a.repozytorium == nil {
		return shared.IdentityDocumentGetResponse{Documents: []shared.IdentityDocument{}}, nil
	}
	filtr := dane.FiltrTozsamosci{
		KodKategorii: wartoscTekstu(z.CategoryId),
		OsByt:        wartoscTekstu(z.AxisId),
	}
	if z.Axis != nil && *z.Axis != "" {
		filtr.Os = konfig.OsLubPlatforma(*z.Axis)
	}
	wiersze, err := a.repozytorium.Dokumenty(ctx, filtr)
	if err != nil {
		return shared.IdentityDocumentGetResponse{}, err
	}
	dokumenty := make([]shared.IdentityDocument, 0, len(wiersze))
	for _, wiersz := range wiersze {
		dokumenty = append(dokumenty, dokumentKontraktu(wiersz))
	}
	return shared.IdentityDocumentGetResponse{Documents: dokumenty}, nil
}

// Zapisz utrwala treść kategorii dla wskazanej osi wraz z trybem podania.
// Odcisk treści liczony jest tym samym skrótem, co odcisk promptu — po nim
// poznaje się, że nakładka nie zmieniła się mimo zapisu.
func (a *adapterTozsamosci) Zapisz(ctx context.Context,
	z shared.IdentityDocumentSetRequest) (shared.IdentityDocumentSetResponse, error) {

	if a == nil || a.repozytorium == nil {
		return shared.IdentityDocumentSetResponse{}, bladBrakuKatalogu("tożsamości modelu")
	}
	kategoria, err := a.repozytorium.KategoriaPoKodzie(ctx, z.CategoryId)
	if err != nil {
		return shared.IdentityDocumentSetResponse{}, bladWskazania(err, "kategoria zasad", z.CategoryId)
	}
	os := konfig.OsPlatformy
	if z.Axis != nil && *z.Axis != "" {
		os = konfig.OsLubPlatforma(*z.Axis)
	}
	dokument := dane.DokumentTozsamosci{
		KodKategorii: kategoria.Kod, Os: os, OsByt: wartoscTekstu(z.AxisId),
		Tryb:  a.trybZapisu(z.Mode, kategoria.TrybDomyslny),
		Tresc: z.Content, OdciskTresci: odciskPromptu(z.Content),
		Aktywny: z.Enabled == nil || *z.Enabled,
	}
	zapisany, err := a.repozytorium.ZapiszDokument(ctx, dokument)
	if err != nil {
		return shared.IdentityDocumentSetResponse{}, err
	}
	return shared.IdentityDocumentSetResponse{Document: dokumentKontraktu(zapisany)}, nil
}

// Usun kasuje zapis treści. Brak zapisu znaczy treść z osi szerszej,
// więc usunięcie zapisu, którego nie ma, nie jest błędem.
func (a *adapterTozsamosci) Usun(ctx context.Context,
	z shared.IdentityDocumentRemoveRequest) (shared.IdentityDocumentRemoveResponse, error) {

	if a == nil || a.repozytorium == nil {
		return shared.IdentityDocumentRemoveResponse{Removed: false}, nil
	}
	id, err := strconv.ParseInt(z.DocumentId, 10, 64)
	if err != nil {
		return shared.IdentityDocumentRemoveResponse{Removed: false}, nil
	}
	usuniete, err := a.repozytorium.UsunDokument(ctx, id)
	if err != nil {
		return shared.IdentityDocumentRemoveResponse{}, err
	}
	return shared.IdentityDocumentRemoveResponse{Removed: usuniete}, nil
}

// Obowiazujaca zwraca nakładkę obowiązującą: tryb, warstwy wg krytyczności,
// złożony prompt i wykaz kategorii obowiązkowych bez treści.
func (a *adapterTozsamosci) Obowiazujaca(ctx context.Context,
	z shared.IdentityEffectiveGetRequest) (shared.IdentityEffectiveGetResponse, error) {

	if a == nil || a.skladacz == nil {
		return shared.IdentityEffectiveGetResponse{Layers: []shared.IdentityLayerContent{}}, nil
	}
	zapytanie, err := a.zapytanieOsi(ctx, z)
	if err != nil {
		return shared.IdentityEffectiveGetResponse{}, err
	}
	zapytanie.TrybDomyslny = a.trybDomyslny(zapytanie)
	return a.skladacz.Zloz(ctx, zapytanie)
}

// trybZapisu rozstrzyga tryb podania zapisywanej treści: wskazany w żądaniu,
// a w jego braku — proponowany przez kategorię katalogu.
func (a *adapterTozsamosci) trybZapisu(wskazany *shared.IdentityMode,
	kategorii shared.IdentityMode) shared.IdentityMode {

	if wskazany != nil && *wskazany != "" {
		return trybLubDomyslny(*wskazany, kategorii)
	}
	return trybLubDomyslny(kategorii, shared.IdentityModeZASTAP)
}

// trybDomyslny odczytuje klucz `tozsamosc.tryb_domyslny` w kontekście osi.
// Brak rozstrzygacza i brak zapisu dają tryb ZASTAP.
func (a *adapterTozsamosci) trybDomyslny(z ZapytanieTozsamosci) shared.IdentityMode {
	if a.rozstrzygacz == nil {
		return shared.IdentityModeZASTAP
	}
	kontekst := konfig.Kontekst{Model: z.Model, Konto: z.Konto}
	wynik := a.rozstrzygacz.Rozstrzygnij(kontekst, KluczTozsamoscTrybDomyslny)
	return trybLubDomyslny(shared.IdentityMode(wynik.Wartosc), shared.IdentityModeZASTAP)
}

// zapytanieOsi ustala, dla czego liczona jest nakładka. Wskazanie okna ma
// pierwszeństwo, bo okno zna i model, i konto naraz; bez okna obowiązuje oś
// podana wprost w żądaniu.
func (a *adapterTozsamosci) zapytanieOsi(ctx context.Context,
	z shared.IdentityEffectiveGetRequest) (ZapytanieTozsamosci, error) {

	if z.WindowId != nil && *z.WindowId != "" && a.osie != nil {
		return a.osie.dlaOkna(ctx, *z.WindowId)
	}
	zapytanie := ZapytanieTozsamosci{}
	if z.Axis == nil {
		return zapytanie, nil
	}
	switch konfig.OsLubPlatforma(*z.Axis) {
	case konfig.OsModelu:
		zapytanie.Model = wartoscTekstu(z.AxisId)
	case konfig.OsKonta:
		zapytanie.Konto = wartoscTekstu(z.AxisId)
	}
	return zapytanie, nil
}
