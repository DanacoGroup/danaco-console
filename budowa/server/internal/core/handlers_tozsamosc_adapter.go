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

type adapterTozsamosci struct {
	repozytorium dane.RepozytoriumTozsamosci
	skladacz     *SkladaczTozsamosci
	rozstrzygacz *konfig.Rozstrzygacz
	osie         *wskazanieOsiOkna
}

func nowyAdapterTozsamosci(repozytorium dane.RepozytoriumTozsamosci,
	rozstrzygacz *konfig.Rozstrzygacz) *adapterTozsamosci {

	return &adapterTozsamosci{
		repozytorium: repozytorium,
		skladacz:     NowySkladaczTozsamosci(NoweZrodloTozsamosci(repozytorium)),
		rozstrzygacz: rozstrzygacz,
	}
}

func (a *adapterTozsamosci) ZOknami(okna dane.RepozytoriumOkien,
	kanaly dane.RepozytoriumKanalow) *adapterTozsamosci {

	a.osie = nowaWskazanieOsiOkna(okna, kanaly)
	return a
}

func (a *adapterTozsamosci) Kategorie(ctx context.Context,
	z shared.IdentityCategoryListRequest) (shared.IdentityCategoryListResponse, error) {

	if a == nil || a.repozytorium == nil {
		return shared.IdentityCategoryListResponse{}, bladBrakuKatalogu("tożsamości modelu")
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

func (a *adapterTozsamosci) Dokumenty(ctx context.Context,
	z shared.IdentityDocumentGetRequest) (shared.IdentityDocumentGetResponse, error) {

	if a == nil || a.repozytorium == nil {
		return shared.IdentityDocumentGetResponse{}, bladBrakuKatalogu("tożsamości modelu")
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

func (a *adapterTozsamosci) Usun(ctx context.Context,
	z shared.IdentityDocumentRemoveRequest) (shared.IdentityDocumentRemoveResponse, error) {

	if a == nil || a.repozytorium == nil {
		return shared.IdentityDocumentRemoveResponse{}, bladBrakuKatalogu("tożsamości modelu")
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

func (a *adapterTozsamosci) Obowiazujaca(ctx context.Context,
	z shared.IdentityEffectiveGetRequest) (shared.IdentityEffectiveGetResponse, error) {

	if a == nil || a.skladacz == nil {
		return shared.IdentityEffectiveGetResponse{Layers: []shared.IdentityLayerContent{}}, nil
	}
	zapytanie, err := a.zapytanieOsi(ctx, z)
	if err != nil {
		return shared.IdentityEffectiveGetResponse{}, err
	}
	zapytanie.TrybDomyslny = a.trybDomyslny(ctx, zapytanie)
	return a.skladacz.Zloz(ctx, zapytanie)
}

func (a *adapterTozsamosci) trybZapisu(wskazany *shared.IdentityMode,
	kategorii shared.IdentityMode) shared.IdentityMode {

	if wskazany != nil && *wskazany != "" {
		return trybLubDomyslny(*wskazany, kategorii)
	}
	return trybLubDomyslny(kategorii, shared.IdentityModeZASTAP)
}

func (a *adapterTozsamosci) trybDomyslny(ctx context.Context, z ZapytanieTozsamosci) shared.IdentityMode {
	if a.rozstrzygacz == nil {
		return shared.IdentityModeZASTAP
	}
	kontekst := konfig.Kontekst{Model: z.Model, Konto: z.Konto, KontoOperatora: dane.KontoOperatora(ctx)}
	wynik := a.rozstrzygacz.Rozstrzygnij(kontekst, KluczTozsamoscTrybDomyslny)
	return trybLubDomyslny(shared.IdentityMode(wynik.Wartosc), shared.IdentityModeZASTAP)
}

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
