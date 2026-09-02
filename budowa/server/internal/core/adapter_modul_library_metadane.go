// Moduł Library obsługuje opis zasobu i schemat metadanych, w tym scalanie oraz podmianę
// pól: `library.metadata.get`, `library.metadata.set`, `library.schema.get`, `library.schema.set`.
package core

import (
	"context"
	"encoding/json"
	"strings"

	"danacoconsole/server/internal/dane"
	"danacoconsole/shared"
)

// Opis obsługuje `library.metadata.get` i dokłada metadane techniczne osadzone w pliku
// wyłącznie na wyraźne żądanie, bo ich odczyt wymaga otwarcia bajtów zasobu.
func (a *adapterBiblioteki) Opis(ctx context.Context,
	z shared.LibraryMetadataGetRequest) (shared.LibraryMetadataGetResponse, error) {

	plik, err := a.plik(ctx, z.FileId)
	if err != nil {
		return shared.LibraryMetadataGetResponse{}, err
	}
	wiersz, err := a.repozytorium.Opis(ctx, plik.ID)
	if err != nil {
		return shared.LibraryMetadataGetResponse{}, bladBiblioteki(err)
	}
	opis := opisKontraktuBiblioteki(plik.Kod, wiersz)
	if z.IncludeTechnical != nil && *z.IncludeTechnical {
		opis.Technical = a.metadaneTechniczne(plik)
	}
	// Odczyt opisu jest dostępem do zasobu i zostawia ślad w dzienniku audytu.
	a.odnotuj(ctx, shared.LibraryAuditActionAccess, &plik.Kod, "odczyt opisu zasobu")
	return shared.LibraryMetadataGetResponse{Metadata: opis}, nil
}

// ZapiszOpis obsługuje `library.metadata.set` i utrwala opis zasobu w tabeli towarzyszącej, scalając albo podmieniając pola.
func (a *adapterBiblioteki) ZapiszOpis(ctx context.Context,
	z shared.LibraryMetadataSetRequest) (shared.LibraryMetadataSetResponse, error) {

	plik, err := a.plik(ctx, z.FileId)
	if err != nil {
		return shared.LibraryMetadataSetResponse{}, err
	}
	zastany, err := a.repozytorium.Opis(ctx, plik.ID)
	if err != nil {
		return shared.LibraryMetadataSetResponse{}, bladBiblioteki(err)
	}

	podmien := z.Replace != nil && *z.Replace
	poZmianie := dane.OpisZasobuBiblioteki{}
	if !podmien {
		poZmianie = zastany
	}
	nadpiszOpisBiblioteki(&poZmianie, z.Metadata, podmien)

	pola, err := polaNiestandardoweBiblioteki(zastany.PolaNiestandardowe, z.Metadata.Custom, podmien)
	if err != nil {
		return shared.LibraryMetadataSetResponse{}, err
	}
	poZmianie.PolaNiestandardowe = pola

	if err := a.repozytorium.ZapiszOpis(ctx, plik.ID, poZmianie); err != nil {
		return shared.LibraryMetadataSetResponse{}, bladBiblioteki(err)
	}
	zapisany, err := a.repozytorium.Opis(ctx, plik.ID)
	if err != nil {
		return shared.LibraryMetadataSetResponse{}, bladBiblioteki(err)
	}
	kontrakt, err := a.zloz(ctx, plik)
	if err != nil {
		return shared.LibraryMetadataSetResponse{}, err
	}
	a.odnotuj(ctx, shared.LibraryAuditActionChange, &plik.Kod, "zapis opisu zasobu")
	a.zglosNasluchom(ctx, shared.LibraryWebhookEventFileChanged, plik.Kod)
	return shared.LibraryMetadataSetResponse{
		Metadata: opisKontraktuBiblioteki(plik.Kod, zapisany), File: kontrakt,
	}, nil
}

// SchematMetadanych obsługuje `library.schema.get` i zwraca definicje pól dopasowane do rodzaju oraz kolekcji zasobu.
func (a *adapterBiblioteki) SchematMetadanych(ctx context.Context,
	z shared.LibrarySchemaGetRequest) (shared.LibrarySchemaGetResponse, error) {

	pola, err := a.repozytorium.PolaSchematu(ctx, z.MimeType, z.CollectionId)
	if err != nil {
		return shared.LibrarySchemaGetResponse{}, bladBiblioteki(err)
	}
	definicje := make([]shared.LibraryFieldDefinition, 0, len(pola))
	for _, pole := range pola {
		definicje = append(definicje, definicjaPolaBiblioteki(pole))
	}
	return shared.LibrarySchemaGetResponse{Fields: definicje}, nil
}

// UstawPoleSchematu obsługuje `library.schema.set` i liczy zasoby z wartością pola przed
// zdjęciem jego definicji, tak by Operator widział zasięg zmiany.
func (a *adapterBiblioteki) UstawPoleSchematu(ctx context.Context,
	z shared.LibrarySchemaSetRequest) (shared.LibrarySchemaSetResponse, error) {

	kod := strings.TrimSpace(z.Field.Code)
	if kod == "" {
		return shared.LibrarySchemaSetResponse{}, bladWskazaniaBiblioteki("definicja pola bez kodu")
	}
	dotkniete, err := a.repozytorium.ZasobyZPolem(ctx, kod)
	if err != nil {
		return shared.LibrarySchemaSetResponse{}, bladBiblioteki(err)
	}

	if z.Remove != nil && *z.Remove {
		if _, err := a.repozytorium.UsunPoleSchematu(ctx, kod); err != nil {
			return shared.LibrarySchemaSetResponse{}, bladBiblioteki(err)
		}
		a.odnotuj(ctx, shared.LibraryAuditActionChange, nil, "zdjęcie pola schematu "+kod)
		return shared.LibrarySchemaSetResponse{Field: z.Field, AffectedFiles: dotkniete}, nil
	}

	if strings.TrimSpace(z.Field.Label) == "" {
		return shared.LibrarySchemaSetResponse{}, bladWskazaniaBiblioteki(
			"pole schematu " + kod + " bez nazwy widocznej dla Operatora")
	}
	opcje, err := opcjePolaBiblioteki(z.Field)
	if err != nil {
		return shared.LibrarySchemaSetResponse{}, err
	}
	zapisane, err := a.repozytorium.ZapiszPoleSchematu(ctx, dane.PoleSchematuBiblioteki{
		Kod: kod, Etykieta: z.Field.Label, Rodzaj: rodzajPolaBazy(z.Field.Kind),
		Wymagane: z.Field.Required, MimeType: z.Field.AppliesToMimeType,
		KolekcjaKod: z.Field.AppliesToCollectionId, Opcje: opcje,
	})
	if err != nil {
		return shared.LibrarySchemaSetResponse{}, bladBiblioteki(err)
	}
	a.odnotuj(ctx, shared.LibraryAuditActionChange, nil, "zapis pola schematu "+kod)
	return shared.LibrarySchemaSetResponse{
		Field: definicjaPolaBiblioteki(zapisane), AffectedFiles: dotkniete,
	}, nil
}

// nadpiszOpisBiblioteki przenosi pola żądania na wiersz opisu, po wskaźnikach: puste pole
// kasuje wartość, pole pominięte przy scalaniu ją zostawia.
func nadpiszOpisBiblioteki(cel *dane.OpisZasobuBiblioteki, zrodlo shared.LibraryMetadata, podmien bool) {
	przypisz := func(docelowe **string, przyslane *string) {
		if przyslane == nil {
			if podmien {
				*docelowe = nil
			}
			return
		}
		if strings.TrimSpace(*przyslane) == "" {
			*docelowe = nil
			return
		}
		wartosc := *przyslane
		*docelowe = &wartosc
	}
	przypisz(&cel.Tytul, zrodlo.Title)
	przypisz(&cel.Tworca, zrodlo.Creator)
	przypisz(&cel.Temat, zrodlo.Subject)
	przypisz(&cel.Opis, zrodlo.Description)
	przypisz(&cel.Wydawca, zrodlo.Publisher)
	przypisz(&cel.Wspoltworca, zrodlo.Contributor)
	przypisz(&cel.DataZasobu, zrodlo.Date)
	przypisz(&cel.Rodzaj, zrodlo.Type)
	przypisz(&cel.Format, zrodlo.Format)
	przypisz(&cel.Identyfikator, zrodlo.Identifier)
	przypisz(&cel.Zrodlo, zrodlo.Source)
	przypisz(&cel.Jezyk, zrodlo.Language)
	przypisz(&cel.Powiazanie, zrodlo.Relation)
	przypisz(&cel.Zakres, zrodlo.Coverage)
	przypisz(&cel.Prawa, zrodlo.Rights)
}

// polaNiestandardoweBiblioteki składa mapę pól niestandardowych po zmianie, scalając po
// kluczach: klucz z wartością pustą znika, klucz pominięty zostaje.
func polaNiestandardoweBiblioteki(zastane *string, przyslane json.RawMessage,
	podmien bool) (*string, error) {

	if len(przyslane) == 0 || string(przyslane) == "null" {
		if podmien {
			return nil, nil
		}
		return zastane, nil
	}
	nowe := map[string]string{}
	if err := json.Unmarshal(przyslane, &nowe); err != nil {
		return nil, bladWskazaniaBiblioteki(
			"pola niestandardowe podaje się jako zbiór par kod-wartość: " + err.Error())
	}
	wynik := map[string]string{}
	if !podmien && zastane != nil && *zastane != "" {
		if err := json.Unmarshal([]byte(*zastane), &wynik); err != nil {
			// Zapis zastany nieczytelny ustępuje wartości przysłanej.
			wynik = map[string]string{}
		}
	}
	for kod, wartosc := range nowe {
		if strings.TrimSpace(wartosc) == "" {
			delete(wynik, kod)
			continue
		}
		wynik[kod] = wartosc
	}
	if len(wynik) == 0 {
		return nil, nil
	}
	tresc, err := json.Marshal(wynik)
	if err != nil {
		return nil, bladBiblioteki(err)
	}
	zapis := string(tresc)
	return &zapis, nil
}

// opisKontraktuBiblioteki przenosi wiersz opisu z bazy na strukturę odpowiedzi kontraktu, dołączając pola niestandardowe.
func opisKontraktuBiblioteki(kodPliku string, wiersz dane.OpisZasobuBiblioteki) shared.LibraryMetadata {
	opis := shared.LibraryMetadata{
		FileId: kodPliku, Title: wiersz.Tytul, Creator: wiersz.Tworca, Subject: wiersz.Temat,
		Description: wiersz.Opis, Publisher: wiersz.Wydawca, Contributor: wiersz.Wspoltworca,
		Date: wiersz.DataZasobu, Type: wiersz.Rodzaj, Format: wiersz.Format,
		Identifier: wiersz.Identyfikator, Source: wiersz.Zrodlo, Language: wiersz.Jezyk,
		Relation: wiersz.Powiazanie, Coverage: wiersz.Zakres, Rights: wiersz.Prawa,
	}
	if wiersz.PolaNiestandardowe != nil && *wiersz.PolaNiestandardowe != "" {
		opis.Custom = json.RawMessage(*wiersz.PolaNiestandardowe)
	}
	return opis
}

// definicjaPolaBiblioteki przenosi wiersz definicji pola z bazy na strukturę kontraktu, dołączając słownik dopuszczalnych wartości.
func definicjaPolaBiblioteki(pole dane.PoleSchematuBiblioteki) shared.LibraryFieldDefinition {
	definicja := shared.LibraryFieldDefinition{
		Code: pole.Kod, Label: pole.Etykieta, Kind: rodzajPolaKontraktu(pole.Rodzaj),
		Required: pole.Wymagane, AppliesToMimeType: pole.MimeType,
		AppliesToCollectionId: pole.KolekcjaKod, CreatedAt: chwilaBazy(pole.Utworzono),
	}
	if pole.Opcje != nil && *pole.Opcje != "" {
		var opcje []string
		if err := json.Unmarshal([]byte(*pole.Opcje), &opcje); err == nil {
			definicja.Options = opcje
		}
	}
	return definicja
}

// opcjePolaBiblioteki sprawdza słownik wartości pola i składa go do zapisu; pole rodzaju
// „lista” bez słownika jest odmową, bo nie da się go wypełnić.
func opcjePolaBiblioteki(pole shared.LibraryFieldDefinition) (*string, error) {
	if pole.Kind != shared.LibraryFieldKindList {
		if len(pole.Options) == 0 {
			return nil, nil
		}
		return nil, bladWskazaniaBiblioteki(
			"słownik wartości ma sens wyłącznie przy polu rodzaju „lista”")
	}
	if len(pole.Options) == 0 {
		return nil, bladWskazaniaBiblioteki(
			"pole rodzaju „lista” bez słownika wartości nie da się wypełnić")
	}
	tresc, err := json.Marshal(pole.Options)
	if err != nil {
		return nil, bladBiblioteki(err)
	}
	zapis := string(tresc)
	return &zapis, nil
}

// rodzajPolaBazy i rodzajPolaKontraktu przekładają wyliczenie kontraktu na wartość kolumny
// bazy i z powrotem, w dwie strony tego samego odwzorowania.
func rodzajPolaBazy(rodzaj shared.LibraryFieldKind) string {
	switch rodzaj {
	case shared.LibraryFieldKindNumber:
		return "liczba"
	case shared.LibraryFieldKindDate:
		return "data"
	case shared.LibraryFieldKindBoolean:
		return "logiczna"
	case shared.LibraryFieldKindList:
		return "lista"
	default:
		return "tekst"
	}
}

func rodzajPolaKontraktu(rodzaj string) shared.LibraryFieldKind {
	switch rodzaj {
	case "liczba":
		return shared.LibraryFieldKindNumber
	case "data":
		return shared.LibraryFieldKindDate
	case "logiczna":
		return shared.LibraryFieldKindBoolean
	case "lista":
		return shared.LibraryFieldKindList
	default:
		return shared.LibraryFieldKindText
	}
}
