// Plik obsługuje warsztat szablonów pism modułu studio: założenie szablonu
// z dokumentu, pola do wypełnienia, zmianę i usunięcie szablonu własnego,
// wniesienie z pliku (`.dotx`, `.ott`), wydanie do pliku oraz wypełnienie pól.
package core

import (
	"context"
	"encoding/json"
	"sort"
	"strconv"
	"strings"
	"time"

	"danacoconsole/server/internal/dane"
	"danacoconsole/shared"
)

// szablonZnacznikPola opisuje, jak pole stoi w treści szablonu. Zapis jest ten
// sam, którym posługuje się `studio.template.apply` (`wypelnijSzablon`) —
// drugiego zapisu znacznika nie zakładam, bo szablony założone starą drogą
// przestałyby się wypełniać.
const (
	szablonOtwarciePola  = "{{"
	szablonZamkniecePola = "}}"
)

// ── Odczyt szablonu ─────────────────────────────────────────────────────────

// szablonPismaSzczegol wczytuje szablon wraz z postacią wzorcową, polami
// i blokadami wzorcowymi, łącząc wiersz szablonu z blokadami z osobnej warstwy.
func (a *adapterStudia) szablonPismaSzczegol(ctx context.Context,
	kod string) (shared.StudioTemplateDetail, error) {

	czysty := strings.TrimSpace(kod)
	if czysty == "" {
		return shared.StudioTemplateDetail{}, bladWskazaniaStudio(
			"czynność na szablonie bez wskazania szablonu")
	}
	skladnica, err := a.wejscieSkladnica()
	if err != nil {
		return shared.StudioTemplateDetail{}, err
	}
	wiersz, err := skladnica.SzablonWarsztatu(ctx, czysty)
	if err != nil {
		if wejscieBrakWiersza(err) {
			return shared.StudioTemplateDetail{}, wejscieBladBraku(
				"szablon nie istnieje: " + czysty)
		}
		return shared.StudioTemplateDetail{}, bladStudio(err)
	}
	return a.szablonZlozKontrakt(ctx, wiersz)
}

// szablonZlozKontrakt składa szablon kontraktu z wiersza warstwy danych,
// łącząc pola i postać wzorcową odczytane z jego zapisu.
func (a *adapterStudia) szablonZlozKontrakt(ctx context.Context,
	wiersz dane.SzablonWarsztatuStudia) (shared.StudioTemplateDetail, error) {

	szablon := shared.StudioTemplateDetail{
		Id:               wiersz.Kod,
		Name:             wiersz.Nazwa,
		Description:      wiersz.Opis,
		Category:         wiersz.Kategoria,
		Format:           shared.StudioDocumentFormat(wiersz.Format),
		Builtin:          wiersz.Fabryczny,
		ThumbnailAssetId: wiersz.MiniaturaZasobKod,
	}
	if wiersz.Zaktualizowano != nil && *wiersz.Zaktualizowano != "" {
		szablon.UpdatedAt = wejscieWskaznikDlugi(chwilaBazy(*wiersz.Zaktualizowano))
	} else {
		szablon.UpdatedAt = wejscieWskaznikDlugi(chwilaBazy(wiersz.Utworzono))
	}

	pola, err := szablonPolaZZapisu(wiersz.PolaJSON, wiersz.Tresc)
	if err != nil {
		return shared.StudioTemplateDetail{}, err
	}
	szablon.Fields = pola

	postac, err := szablonPostacZZapisu(wiersz)
	if err != nil {
		return shared.StudioTemplateDetail{}, err
	}
	szablon.Form = postac

	// Blokady wzorcowe leżą w tabeli obszaru kontroli pracy, drugiej niż tabela szablonu.
	if skladnica, err := a.kontrolaSkladnica(); err == nil {
		wzorce, err := skladnica.BlokadySzablonu(ctx, wiersz.Kod)
		if err != nil {
			return shared.StudioTemplateDetail{}, bladStudio(err)
		}
		blokady := make([]shared.StudioFragmentLock, 0, len(wzorce))
		for _, wzorzec := range wzorce {
			blokady = append(blokady, shared.StudioFragmentLock{
				Id:           wzorzec.Kod,
				DocumentId:   "",
				Name:         wzorzec.Nazwa,
				Reason:       wzorzec.Powod,
				RangeStart:   int(wzorzec.ZakresOd),
				RangeEnd:     int(wzorzec.ZakresDo),
				Scope:        shared.StudioLockScope(wzorzec.Zasieg),
				FromTemplate: wejscieWskaznikLogiczny(true),
			})
		}
		if len(blokady) > 0 {
			szablon.Locks = blokady
		}
	}
	return szablon, nil
}

// szablonPolaZZapisu czyta pola szablonu z ładunku JSON: zapis starszy niesie
// mapę nazwa-wartość, zapis nowy — kształt `StudioTemplateFieldSpec[]`, czytane
// są oba. Ładunek nieczytelny nie schodzi na pusty wykaz pól.
func szablonPolaZZapisu(zapis *string, tresc string) ([]shared.StudioTemplateFieldSpec, error) {
	if zapis == nil || strings.TrimSpace(*zapis) == "" {
		// Pól nie zapisano — wykaz składa się ze znaczników stojących w treści.
		return szablonPolaZTresci(tresc, nil), nil
	}
	surowy := strings.TrimSpace(*zapis)

	var wykaz []shared.StudioTemplateFieldSpec
	if err := json.Unmarshal([]byte(surowy), &wykaz); err == nil {
		return szablonPolaZTresci(tresc, wykaz), nil
	}

	var starszy map[string]string
	if err := json.Unmarshal([]byte(surowy), &starszy); err == nil {
		nazwy := make([]string, 0, len(starszy))
		for nazwa := range starszy {
			nazwy = append(nazwy, nazwa)
		}
		sort.Strings(nazwy)
		przelozone := make([]shared.StudioTemplateFieldSpec, 0, len(nazwy))
		for _, nazwa := range nazwy {
			pole := shared.StudioTemplateFieldSpec{
				Name:  nazwa,
				Label: nazwa,
				Kind:  shared.StudioTemplateFieldKind(shared.StudioTemplateFieldKindText),
			}
			if wartosc := strings.TrimSpace(starszy[nazwa]); wartosc != "" {
				pole.DefaultValue = wejscieWskaznikTekstu(wartosc)
			}
			przelozone = append(przelozone, pole)
		}
		return szablonPolaZTresci(tresc, przelozone), nil
	}

	var nazwy []string
	if err := json.Unmarshal([]byte(surowy), &nazwy); err == nil {
		przelozone := make([]shared.StudioTemplateFieldSpec, 0, len(nazwy))
		for _, nazwa := range nazwy {
			przelozone = append(przelozone, shared.StudioTemplateFieldSpec{
				Name:  nazwa,
				Label: nazwa,
				Kind:  shared.StudioTemplateFieldKind(shared.StudioTemplateFieldKindText),
			})
		}
		return szablonPolaZTresci(tresc, przelozone), nil
	}

	return nil, wejscieBladZaplecza(
		"wykaz pól szablonu jest nieczytelny — ani wykazem pól, ani mapą wartości: " + surowy)
}

// szablonPolaZTresci dokłada do wykazu pola, które stoją w treści szablonu,
// a w wykazie ich nie ma — inaczej byłyby niemożliwe do wypełnienia.
func szablonPolaZTresci(tresc string,
	wykaz []shared.StudioTemplateFieldSpec) []shared.StudioTemplateFieldSpec {

	znane := map[string]bool{}
	for _, pole := range wykaz {
		znane[pole.Name] = true
	}
	if wykaz == nil {
		wykaz = []shared.StudioTemplateFieldSpec{}
	}
	for _, nazwa := range szablonNazwyZeTresci(tresc) {
		if znane[nazwa] {
			continue
		}
		znane[nazwa] = true
		wykaz = append(wykaz, shared.StudioTemplateFieldSpec{
			Name:  nazwa,
			Label: nazwa,
			Kind:  shared.StudioTemplateFieldKind(shared.StudioTemplateFieldKindText),
		})
	}
	// Miejsce pola w treści liczone znakami — okno stawia po nim kursor przy wypełnianiu.
	for i := range wykaz {
		if wykaz[i].AnchorOffset != nil {
			continue
		}
		znacznik := szablonOtwarciePola + wykaz[i].Name + szablonZamkniecePola
		if wskazanie := strings.Index(tresc, znacznik); wskazanie >= 0 {
			wykaz[i].AnchorOffset = wejscieWskaznikCalkowity(
				len([]rune(tresc[:wskazanie])))
		}
	}
	return wykaz
}

// szablonNazwyZeTresci wymienia nazwy pól stojących w treści szablonu,
// odczytane ze znaczników postaci wzorcowej.
func szablonNazwyZeTresci(tresc string) []string {
	nazwy := []string{}
	reszta := tresc
	for {
		poczatek := strings.Index(reszta, szablonOtwarciePola)
		if poczatek < 0 {
			return nazwy
		}
		reszta = reszta[poczatek+len(szablonOtwarciePola):]
		koniec := strings.Index(reszta, szablonZamkniecePola)
		if koniec < 0 {
			return nazwy
		}
		nazwa := strings.TrimSpace(reszta[:koniec])
		reszta = reszta[koniec+len(szablonZamkniecePola):]
		if nazwa == "" || strings.ContainsAny(nazwa, "{}\n") {
			continue
		}
		nazwy = append(nazwy, nazwa)
	}
}

// szablonPostacZZapisu czyta postać wzorcową szablonu z zapisu; brakujący zapis
// buduje postać wprost z treści szablonu.
func szablonPostacZZapisu(wiersz dane.SzablonWarsztatuStudia) (*shared.StudioDocumentForm, error) {
	if wiersz.PostacJSON == nil || strings.TrimSpace(*wiersz.PostacJSON) == "" {
		// Szablon bez zapisanej postaci, na przykład fabryczny — postać składa się z treści szablonu.
		if strings.TrimSpace(wiersz.Tresc) == "" {
			return nil, nil
		}
		postac := wejsciePostacZTekstu("", wiersz.Tresc)
		return &postac, nil
	}
	var postac shared.StudioDocumentForm
	if err := json.Unmarshal([]byte(*wiersz.PostacJSON), &postac); err != nil {
		return nil, wejscieBladZaplecza("postać wzorcowa szablonu " + wiersz.Kod +
			" jest nieczytelna: " + err.Error())
	}
	return &postac, nil
}

// ── Zapis szablonu z dokumentu ──────────────────────────────────────────────

// ZapiszSzablonPisma obsługuje `studio.template.save` — zakłada szablon pisma
// z bieżącego dokumentu wraz z arkuszem stylów, nastawami strony, nagłówkiem,
// stopką, logo, tabelami i blokadami wzorcowymi.
func (a *adapterStudia) ZapiszSzablonPisma(ctx context.Context,
	z shared.StudioTemplateSaveRequest) (shared.StudioTemplateSaveResponse, error) {

	if strings.TrimSpace(z.Name) == "" {
		return shared.StudioTemplateSaveResponse{}, bladWskazaniaStudio(
			"zapis szablonu bez nazwy — szablon bez nazwy jest niewybieralny w galerii")
	}
	if _, err := wejscieAutorCzynnosci(z.Author); err != nil {
		return shared.StudioTemplateSaveResponse{}, err
	}
	skladnica, err := a.wejscieSkladnica()
	if err != nil {
		return shared.StudioTemplateSaveResponse{}, err
	}

	kod := strings.TrimSpace(wartoscTekstu(z.TemplateId))
	if kod != "" {
		zastany, err := skladnica.SzablonWarsztatu(ctx, kod)
		if err != nil {
			if wejscieBrakWiersza(err) {
				return shared.StudioTemplateSaveResponse{}, wejscieBladBraku(
					"szablon nie istnieje: " + kod)
			}
			return shared.StudioTemplateSaveResponse{}, bladStudio(err)
		}
		if zastany.Fabryczny {
			return shared.StudioTemplateSaveResponse{}, bladWskazaniaStudio(
				"szablon „" + zastany.Nazwa + "” jest fabryczny i zmianie nie podlega. " +
					"Naprawa: zapisać go pod nową nazwą jako szablon własny — wtedy " +
					"fabryczny zostaje wzorcem, a Operator pracuje na swoim")
		}
	} else {
		kod = nowyIdentyfikator(przedrostekSzablonuStudia)
	}

	wiersz := dane.SzablonWarsztatuStudia{
		Kod:       kod,
		Nazwa:     strings.TrimSpace(z.Name),
		Opis:      z.Description,
		Kategoria: z.Category,
		Format:    string(shared.StudioDocumentFormatMarkdown),
		// Szablon zakładany bez dokumentu jest szablonem samych pól i nastaw — treść dostaje pustą.
		Tresc:             "",
		MiniaturaZasobKod: z.ThumbnailAssetId,
	}

	var stan *stanPostaci
	if kodDokumentu := strings.TrimSpace(wartoscTekstu(z.DocumentId)); kodDokumentu != "" {
		wczytany, err := a.postacWczytaj(ctx, kodDokumentu)
		if err != nil {
			return shared.StudioTemplateSaveResponse{}, err
		}
		stan = wczytany
		wiersz.Tresc = wartoscTekstu(stan.dokument.Tresc)
		wiersz.Format = string(stan.dokument.Format)
		wiersz.DokumentZrodlowyKod = wejscieWskaznikTekstu(stan.dokument.Kod)

		postac := stan.forma
		// Postać wzorcowa idzie bez identyfikatorów dokumentu źródłowego — szablon dostanie własne.
		postac.DocumentId = ""
		postac.Locks = nil
		postac.Revision = nil
		postac.UpdatedAt = nil
		zapisPostaci, err := json.Marshal(postac)
		if err != nil {
			return shared.StudioTemplateSaveResponse{}, wejscieBladZaplecza(
				"postaci wzorcowej szablonu nie da się zapisać: " + err.Error())
		}
		zapis := string(zapisPostaci)
		wiersz.PostacJSON = &zapis
	}

	pola, err := szablonPolaZadania(z.Fields, wiersz.Tresc)
	if err != nil {
		return shared.StudioTemplateSaveResponse{}, err
	}
	zapisPol, err := json.Marshal(pola)
	if err != nil {
		return shared.StudioTemplateSaveResponse{}, wejscieBladZaplecza(
			"wykazu pól szablonu nie da się zapisać: " + err.Error())
	}
	zapisanePola := string(zapisPol)
	wiersz.PolaJSON = &zapisanePola

	zapisany, err := skladnica.ZapiszSzablonWarsztatu(ctx, wiersz)
	if err != nil {
		return shared.StudioTemplateSaveResponse{}, bladStudio(err)
	}

	// Blokady dokumentu wzorcowego stają się blokadami wzorcowymi szablonu domyślnie.
	if stan != nil && (z.IncludeLocks == nil || *z.IncludeLocks) {
		if err := a.szablonPrzeniesBlokady(ctx, stan, zapisany.Kod); err != nil {
			return shared.StudioTemplateSaveResponse{}, err
		}
	}

	szablon, err := a.szablonZlozKontrakt(ctx, zapisany)
	if err != nil {
		return shared.StudioTemplateSaveResponse{}, err
	}
	return shared.StudioTemplateSaveResponse{Template: szablon}, nil
}

// szablonPolaZadania czyta pola podane w żądaniu. Kształt jest surowym JSON-em,
// bo kontrakt opisuje go typem `json` — szablon nie ma z góry znanego zbioru pól
// i mieć nie może.
func szablonPolaZadania(surowe json.RawMessage,
	tresc string) ([]shared.StudioTemplateFieldSpec, error) {

	if len(surowe) == 0 {
		return szablonPolaZTresci(tresc, nil), nil
	}
	var wykaz []shared.StudioTemplateFieldSpec
	if err := json.Unmarshal(surowe, &wykaz); err != nil {
		return nil, bladWskazaniaStudio(
			"wykaz pól szablonu w żądaniu jest nieczytelny: " + err.Error())
	}
	for i := range wykaz {
		if strings.TrimSpace(wykaz[i].Name) == "" {
			return nil, bladWskazaniaStudio(
				"pole szablonu bez nazwy — nazwa jest tym, co podstawia się w treści")
		}
		if strings.TrimSpace(wykaz[i].Label) == "" {
			wykaz[i].Label = wykaz[i].Name
		}
		if strings.TrimSpace(string(wykaz[i].Kind)) == "" {
			wykaz[i].Kind = shared.StudioTemplateFieldKindText
		}
		if err := szablonSprawdzRodzajPola(wykaz[i]); err != nil {
			return nil, err
		}
	}
	return szablonPolaZTresci(tresc, wykaz), nil
}

// szablonSprawdzRodzajPola sprawdza rodzaj pola wykazem kontraktu i pilnuje, żeby
// pole wyboru miało z czego wybierać.
func szablonSprawdzRodzajPola(pole shared.StudioTemplateFieldSpec) error {
	znany := false
	for _, rodzaj := range shared.WartosciStudioTemplateFieldKind() {
		if pole.Kind == rodzaj {
			znany = true
			break
		}
	}
	if !znany {
		return bladWskazaniaStudio("pole szablonu „" + pole.Name + "” ma rodzaj " +
			string(pole.Kind) + ", którego kontrakt nie zna; rodzaje to tekst, data, " +
			"liczba i wybór z wykazu")
	}
	if pole.Kind == shared.StudioTemplateFieldKindChoice && len(pole.Choices) == 0 {
		return bladWskazaniaStudio("pole szablonu „" + pole.Name +
			"” jest wyborem z wykazu, a wykazu wartości nie ma — nie byłoby z czego wybrać")
	}
	return nil
}

// szablonPrzeniesBlokady przenosi blokady dokumentu wzorcowego do szablonu,
// jako blokady wzorcowe tego szablonu.
func (a *adapterStudia) szablonPrzeniesBlokady(ctx context.Context, stan *stanPostaci,
	kodSzablonu string) error {

	skladnica, err := a.kontrolaSkladnica()
	if err != nil {
		return err
	}
	blokady, err := skladnica.BlokadyFragmentow(ctx, stan.dokument.ID)
	if err != nil {
		return bladStudio(err)
	}
	for _, blokada := range blokady {
		if err := skladnica.ZapiszBlokadeSzablonu(ctx, dane.BlokadaSzablonuStudia{
			Kod:        nowyIdentyfikator(przedrostekBlokadyStudia),
			SzablonKod: kodSzablonu,
			Nazwa:      blokada.Nazwa,
			Powod:      blokada.Powod,
			ZakresOd:   blokada.ZakresOd,
			ZakresDo:   blokada.ZakresDo,
			Zasieg:     blokada.Zasieg,
		}); err != nil {
			return bladStudio(err)
		}
	}
	return nil
}

// ── Usunięcie szablonu ──────────────────────────────────────────────────────

// UsunSzablonPisma obsługuje `studio.template.delete` — usuwa szablon własny.
//
// Szablonu fabrycznego nie usuwa i mówi to WPROST, nazywając powód. Wzór:
// `studio.operation.delete`, który robi to samo dla operacji fabrycznych.
func (a *adapterStudia) UsunSzablonPisma(ctx context.Context,
	z shared.StudioTemplateDeleteRequest) (shared.StudioTemplateDeleteResponse, error) {

	kod := strings.TrimSpace(z.TemplateId)
	if kod == "" {
		return shared.StudioTemplateDeleteResponse{}, bladWskazaniaStudio(
			"usunięcie szablonu bez wskazania szablonu")
	}
	skladnica, err := a.wejscieSkladnica()
	if err != nil {
		return shared.StudioTemplateDeleteResponse{}, err
	}
	zastany, err := skladnica.SzablonWarsztatu(ctx, kod)
	if err != nil {
		if wejscieBrakWiersza(err) {
			return shared.StudioTemplateDeleteResponse{}, wejscieBladBraku(
				"szablon nie istnieje: " + kod)
		}
		return shared.StudioTemplateDeleteResponse{}, bladStudio(err)
	}
	if zastany.Fabryczny {
		return shared.StudioTemplateDeleteResponse{}, bladWskazaniaStudio(
			"szablon „" + zastany.Nazwa + "” jest fabryczny i nie da się go usunąć: " +
				"wykaz fabryczny jest wzorcem platformy, którego nie da się odtworzyć bez " +
				"ponownego wdrożenia. Naprawa: usunąć szablon własny albo — jeśli " +
				"fabryczny przeszkadza — założyć własny i posługiwać się nim")
	}
	usuniety, err := skladnica.UsunSzablonWlasny(ctx, kod)
	if err != nil {
		return shared.StudioTemplateDeleteResponse{}, bladStudio(err)
	}
	if !usuniety {
		return shared.StudioTemplateDeleteResponse{}, wejscieBladZaplecza(
			"szablon „" + zastany.Nazwa + "” istnieje, a zapytanie usuwające nie zdjęło " +
				"ani jednego wiersza — to jest usterka warstwy danych, nie odmowa")
	}
	return shared.StudioTemplateDeleteResponse{Deleted: true}, nil
}

// ── Pola do wypełnienia ─────────────────────────────────────────────────────

// UstawPoleSzablonu obsługuje `studio.template.field.set` — wskazuje pole
// szablonu do wypełnienia: nazwę, opis, wartość domyślną, rodzaj i to, czy jest
// wymagane. Tą samą drogą pole się USUWA (`remove`).
func (a *adapterStudia) UstawPoleSzablonu(ctx context.Context,
	z shared.StudioTemplateFieldSetRequest) (shared.StudioTemplateFieldSetResponse, error) {

	if strings.TrimSpace(z.Name) == "" {
		return shared.StudioTemplateFieldSetResponse{}, bladWskazaniaStudio(
			"czynność na polu szablonu bez nazwy pola")
	}
	skladnica, err := a.wejscieSkladnica()
	if err != nil {
		return shared.StudioTemplateFieldSetResponse{}, err
	}
	kod := strings.TrimSpace(z.TemplateId)
	wiersz, err := skladnica.SzablonWarsztatu(ctx, kod)
	if err != nil {
		if wejscieBrakWiersza(err) {
			return shared.StudioTemplateFieldSetResponse{}, wejscieBladBraku(
				"szablon nie istnieje: " + kod)
		}
		return shared.StudioTemplateFieldSetResponse{}, bladStudio(err)
	}
	if wiersz.Fabryczny {
		return shared.StudioTemplateFieldSetResponse{}, bladWskazaniaStudio(
			"szablon „" + wiersz.Nazwa + "” jest fabryczny — jego pól nie da się zmienić. " +
				"Naprawa: zapisać go jako szablon własny (studio.template.save) i pracować " +
				"na własnym")
	}

	pola, err := szablonPolaZZapisu(wiersz.PolaJSON, wiersz.Tresc)
	if err != nil {
		return shared.StudioTemplateFieldSetResponse{}, err
	}
	nazwa := strings.TrimSpace(z.Name)

	if z.Remove != nil && *z.Remove {
		zostawione := make([]shared.StudioTemplateFieldSpec, 0, len(pola))
		zdjete := false
		for _, pole := range pola {
			if pole.Name == nazwa {
				zdjete = true
				continue
			}
			zostawione = append(zostawione, pole)
		}
		if !zdjete {
			return shared.StudioTemplateFieldSetResponse{}, wejscieBladBraku(
				"szablon „" + wiersz.Nazwa + "” nie ma pola o nazwie „" + nazwa + "”")
		}
		pola = zostawione
	} else {
		wskazanie := -1
		for i := range pola {
			if pola[i].Name == nazwa {
				wskazanie = i
				break
			}
		}
		if wskazanie < 0 {
			pola = append(pola, shared.StudioTemplateFieldSpec{
				Name:  nazwa,
				Label: nazwa,
				Kind:  shared.StudioTemplateFieldKind(shared.StudioTemplateFieldKindText),
			})
			wskazanie = len(pola) - 1
		}
		if z.Label != nil {
			pola[wskazanie].Label = strings.TrimSpace(*z.Label)
		}
		if strings.TrimSpace(pola[wskazanie].Label) == "" {
			pola[wskazanie].Label = nazwa
		}
		if z.Description != nil {
			pola[wskazanie].Description = z.Description
		}
		if z.Kind != nil {
			pola[wskazanie].Kind = *z.Kind
		}
		if z.Required != nil {
			pola[wskazanie].Required = *z.Required
		}
		if z.DefaultValue != nil {
			pola[wskazanie].DefaultValue = z.DefaultValue
		}
		if z.Choices != nil {
			pola[wskazanie].Choices = z.Choices
		}
		if z.AnchorOffset != nil {
			pola[wskazanie].AnchorOffset = z.AnchorOffset
		}
		if err := szablonSprawdzRodzajPola(pola[wskazanie]); err != nil {
			return shared.StudioTemplateFieldSetResponse{}, err
		}
	}

	zapis, err := json.Marshal(pola)
	if err != nil {
		return shared.StudioTemplateFieldSetResponse{}, wejscieBladZaplecza(
			"wykazu pól szablonu nie da się zapisać: " + err.Error())
	}
	zapisane := string(zapis)
	wiersz.PolaJSON = &zapisane
	zapisany, err := skladnica.ZapiszSzablonWarsztatu(ctx, wiersz)
	if err != nil {
		return shared.StudioTemplateFieldSetResponse{}, bladStudio(err)
	}
	szablon, err := a.szablonZlozKontrakt(ctx, zapisany)
	if err != nil {
		return shared.StudioTemplateFieldSetResponse{}, err
	}
	return shared.StudioTemplateFieldSetResponse{
		Template: szablon,
		Fields:   szablon.Fields,
	}, nil
}

// PolaSzablonu obsługuje `studio.template.field.list` — oddaje pola szablonu
// wraz z rodzajem, wartością domyślną i wymagalnością.
func (a *adapterStudia) PolaSzablonu(ctx context.Context,
	z shared.StudioTemplateFieldListRequest) (shared.StudioTemplateFieldListResponse, error) {

	szablon, err := a.szablonPismaSzczegol(ctx, z.TemplateId)
	if err != nil {
		return shared.StudioTemplateFieldListResponse{}, err
	}
	pola := szablon.Fields
	if pola == nil {
		pola = []shared.StudioTemplateFieldSpec{}
	}
	return shared.StudioTemplateFieldListResponse{Fields: pola}, nil
}

// ── Wniesienie szablonu z pliku Operatora ───────────────────────────────────

// WniesSzablonZPliku obsługuje `studio.template.import` — wnosi szablon z pliku
// (`.dotx`, `.ott`) wraz z postacią i polami do wypełnienia. Plik dokumentu
// (`.docx`, `.odt`) też wchodzi, a bilans mówi wprost, że szablon powstał
// z dokumentu.
func (a *adapterStudia) WniesSzablonZPliku(ctx context.Context,
	z shared.StudioTemplateImportRequest) (shared.StudioTemplateImportResponse, error) {

	bajty, nazwaPliku, err := a.wejscieBajtyZrodla(ctx, wejscieWskazanieZrodla{
		Sciezka:        z.Path,
		PlikBiblioteki: z.LibraryFileId,
		BajtyBase64:    z.BytesBase64,
		Czynnosc:       "wniesienie szablonu z pliku",
	})
	if err != nil {
		return shared.StudioTemplateImportResponse{}, err
	}
	format, err := wejscieRozpoznajFormat(nazwaPliku, bajty, z.Format)
	if err != nil {
		return shared.StudioTemplateImportResponse{}, err
	}
	skladnica, err := a.wejscieSkladnica()
	if err != nil {
		return shared.StudioTemplateImportResponse{}, err
	}

	postac, tresc, bilans, err := a.wejscieRozbierzPlik("", bajty, format, nil)
	if err != nil {
		return shared.StudioTemplateImportResponse{}, err
	}
	if !wejscieCzyFormatSzablonu(format) {
		bilans.Skipped = append(bilans.Skipped, shared.StudioSkippedItem{
			Reason: "plik nie był plikiem szablonu, a dokumentem",
			Detail: wejscieWskaznikTekstu("szablon powstał z jego postaci i treści; " +
				"format pliku: " + string(format)),
		})
	}

	nazwa := strings.TrimSpace(wartoscTekstu(z.Name))
	if nazwa == "" {
		if zPliku := wejscieNazwaZPliku(nazwaPliku); zPliku != nil {
			nazwa = *zPliku
		}
	}
	if nazwa == "" {
		return shared.StudioTemplateImportResponse{}, bladWskazaniaStudio(
			"wniesienie szablonu bez nazwy — ani żądanie, ani nazwa pliku jej nie niosły")
	}

	postac.DocumentId = ""
	postac.Locks = nil
	postac.Revision = nil
	postac.UpdatedAt = nil
	zapisPostaci, err := json.Marshal(postac)
	if err != nil {
		return shared.StudioTemplateImportResponse{}, wejscieBladZaplecza(
			"postaci wzorcowej szablonu nie da się zapisać: " + err.Error())
	}
	zapis := string(zapisPostaci)

	// Pola do wypełnienia wychodzą ze znaczników stojących w treści szablonu.
	pola := szablonPolaZTresci(tresc, nil)
	zapisPol, err := json.Marshal(pola)
	if err != nil {
		return shared.StudioTemplateImportResponse{}, wejscieBladZaplecza(
			"wykazu pól szablonu nie da się zapisać: " + err.Error())
	}
	zapisanePola := string(zapisPol)
	zrodlo := string(format)

	zapisany, err := skladnica.ZapiszSzablonWarsztatu(ctx, dane.SzablonWarsztatuStudia{
		Kod:         nowyIdentyfikator(przedrostekSzablonuStudia),
		Nazwa:       nazwa,
		Kategoria:   z.Category,
		Format:      string(wejscieFormatDokumentu(format)),
		Tresc:       tresc,
		PolaJSON:    &zapisanePola,
		PostacJSON:  &zapis,
		ZrodloPliku: &zrodlo,
	})
	if err != nil {
		return shared.StudioTemplateImportResponse{}, bladStudio(err)
	}

	bilans.Note = wejscieWskaznikTekstu("szablon „" + nazwa + "” wniesiony z pliku " +
		string(format) + ": przejęto arkusz stylów (" +
		strconv.Itoa(len(postac.Styles)) + " stylów), nastawy strony, nagłówek i stopkę " +
		"oraz " + strconv.Itoa(len(pola)) + " pól do wypełnienia")

	szablon, err := a.szablonZlozKontrakt(ctx, zapisany)
	if err != nil {
		return shared.StudioTemplateImportResponse{}, err
	}
	return shared.StudioTemplateImportResponse{Template: szablon, Balance: bilans}, nil
}

// ── Oddanie szablonu do pliku ───────────────────────────────────────────────

// OddajSzablonDoPliku obsługuje `studio.template.export` — oddaje szablon do
// pliku wraz z arkuszem stylów, nastawami strony i polami. Format domyślny to
// `docx`; formaty tekstowe też są dopuszczone, z wykazem cech pominiętych.
func (a *adapterStudia) OddajSzablonDoPliku(ctx context.Context,
	z shared.StudioTemplateExportRequest) (shared.StudioTemplateExportResponse, error) {

	szablon, err := a.szablonPismaSzczegol(ctx, z.TemplateId)
	if err != nil {
		return shared.StudioTemplateExportResponse{}, err
	}
	format := shared.StudioExportFormat(shared.StudioExportFormatDocx)
	if z.Format != nil && strings.TrimSpace(string(*z.Format)) != "" {
		sprawdzony, err := wydanieFormatDocelowy(*z.Format)
		if err != nil {
			return shared.StudioTemplateExportResponse{}, err
		}
		format = sprawdzony
	}

	postac := shared.StudioDocumentForm{}
	if szablon.Form != nil {
		postac = *szablon.Form
	}
	tresc := wejscieTrescZPostaci(&postac)
	if strings.TrimSpace(tresc) == "" && len(postac.Blocks) == 0 {
		return shared.StudioTemplateExportResponse{}, wejscieBladBraku(
			"szablon „" + szablon.Name + "” nie niesie ani postaci wzorcowej, ani treści — " +
				"nie ma czego oddać do pliku")
	}

	var bajty []byte
	var pominiete []shared.StudioSkippedItem
	switch format {
	case shared.StudioExportFormatDocx:
		bajty, pominiete, err = wejscieZlozOoxml(&postac, tresc, szablon.Name, true)
	case shared.StudioExportFormatOdt:
		bajty, pominiete, err = wejscieZlozOdf(&postac, tresc, szablon.Name, true)
	case shared.StudioExportFormatTxt:
		bajty, pominiete = wydanieTekstem(&postac, tresc)
	case shared.StudioExportFormatMd:
		bajty, pominiete = wydanieMarkdownem(&postac, tresc)
	case shared.StudioExportFormatHtml:
		bajty, pominiete = wydanieHtmlem(&postac, tresc, szablon.Name)
	case shared.StudioExportFormatRtf:
		bajty, pominiete = wydanieRtfem(&postac, tresc)
	default:
		return shared.StudioTemplateExportResponse{}, bladWskazaniaStudio(
			"szablonu nie oddaje się do formatu " + string(format) +
				"; plik szablonu to dotx (docx) albo ott (odt)")
	}
	if err != nil {
		return shared.StudioTemplateExportResponse{}, err
	}
	if len(bajty) == 0 {
		return shared.StudioTemplateExportResponse{}, wejscieBladZaplecza(
			"oddanie szablonu do formatu " + string(format) +
				" nie złożyło ani jednego bajtu")
	}

	if len(szablon.Fields) > 0 {
		pominiete = append(pominiete, shared.StudioSkippedItem{
			Reason: "pola do wypełnienia wyszły znacznikami w treści, bez wykazu pól pliku",
			Detail: wejscieWskaznikTekstu(strconv.Itoa(len(szablon.Fields)) +
				" pól w zapisie " + szablonOtwarciePola + "nazwa" + szablonZamkniecePola +
				"; wniesienie tego pliku z powrotem odtworzy wykaz z tych znaczników"),
		})
	}
	if len(szablon.Locks) > 0 {
		pominiete = append(pominiete, shared.StudioSkippedItem{
			Reason: "blokady wzorcowe nie weszły do pliku — żaden format biurowy nie niesie " +
				"blokady skierowanej przeciw modelowi",
			Detail: wejscieWskaznikTekstu(strconv.Itoa(len(szablon.Locks)) +
				" blokad; zostają przy szablonie w serwerze i przechodzą do dokumentów " +
				"z niego zakładanych"),
		})
	}

	rozszerzenie := string(format)
	if format == shared.StudioExportFormatDocx {
		rozszerzenie = "dotx"
	}
	if format == shared.StudioExportFormatOdt {
		rozszerzenie = "ott"
	}
	wynik := shared.StudioExportResult{
		DocumentId:      szablon.Id,
		Format:          format,
		Bytes:           wejscieWskaznikDlugi(int64(len(bajty))),
		DroppedFeatures: pominiete,
		Note: wejscieWskaznikTekstu("szablon „" + szablon.Name + "” oddany plikiem " +
			rozszerzenie + "; " + wydanieZdanieBilansu(format, pominiete)),
	}
	// Okno bierze się z dokumentu wzorcowego — zasób bez okna dałby kod bez dostępu do pliku.
	zasob, err := a.odlozTrescStudia(ctx, bajty, szablon.Name+"."+rozszerzenie,
		rozszerzenie, a.szablonOknoWzorca(ctx, szablon.Id))
	if err != nil {
		return shared.StudioTemplateExportResponse{}, err
	}
	wynik.AssetId = wejscieWskaznikTekstu(zasob.Id)

	if sciezka := strings.TrimSpace(wartoscTekstu(z.Path)); sciezka != "" {
		if err := wydanieZapiszPlik(sciezka, bajty); err != nil {
			return shared.StudioTemplateExportResponse{}, err
		}
		wynik.Path = wejscieWskaznikTekstu(sciezka)
	}
	return shared.StudioTemplateExportResponse{Result: wynik}, nil
}

// szablonOknoWzorca oddaje okno dokumentu, z którego szablon powstał, albo
// pustkę, gdy szablon dokumentu źródłowego nie ma (na przykład fabryczny).
func (a *adapterStudia) szablonOknoWzorca(ctx context.Context, kodSzablonu string) string {
	skladnica, err := a.wejscieSkladnica()
	if err != nil {
		return ""
	}
	wiersz, err := skladnica.SzablonWarsztatu(ctx, kodSzablonu)
	if err != nil || wiersz.DokumentZrodlowyKod == nil {
		return ""
	}
	dokument, err := a.repozytorium.Dokument(ctx, *wiersz.DokumentZrodlowyKod)
	if err != nil {
		return ""
	}
	return dokument.Okno
}

// ── Wypełnienie pól ─────────────────────────────────────────────────────────

// WypelnijSzablon obsługuje `studio.template.fill` — wypełnia pola szablonu
// wartościami i oddaje dokument gotowy. Pole wymagane bez wartości nie znika
// z treści: znacznik zostaje widoczny, a odpowiedź wymienia pola brakujące
// w `missingRequired`.
func (a *adapterStudia) WypelnijSzablon(ctx context.Context,
	z shared.StudioTemplateFillRequest) (shared.StudioTemplateFillResponse, error) {

	autor, err := wejscieAutorCzynnosci(z.Author)
	if err != nil {
		return shared.StudioTemplateFillResponse{}, err
	}
	szablon, err := a.szablonPismaSzczegol(ctx, z.TemplateId)
	if err != nil {
		return shared.StudioTemplateFillResponse{}, err
	}

	wartosci := map[string]string{}
	if len(z.Values) > 0 {
		var surowe map[string]any
		if err := json.Unmarshal(z.Values, &surowe); err != nil {
			return shared.StudioTemplateFillResponse{}, bladWskazaniaStudio(
				"wartości pól szablonu są nieczytelne: " + err.Error())
		}
		for nazwa, wartosc := range surowe {
			wartosci[nazwa] = szablonWartoscNapisem(wartosc)
		}
	}

	// Wartość domyślna wchodzi tam, gdzie Operator wartości nie podał — po to
	// jest domyślna.
	brakujace := []string{}
	for _, pole := range szablon.Fields {
		if _, jest := wartosci[pole.Name]; jest {
			continue
		}
		if pole.DefaultValue != nil && strings.TrimSpace(*pole.DefaultValue) != "" {
			wartosci[pole.Name] = *pole.DefaultValue
			continue
		}
		if pole.Kind == shared.StudioTemplateFieldKindDate {
			// Data bez wartości i bez domyślnej jest datą dzisiejszą — tak zachowuje się każde pismo urzędowe.
			wartosci[pole.Name] = time.Now().Format("2006-01-02")
			continue
		}
		if pole.Required {
			brakujace = append(brakujace, pole.Name)
		}
	}

	bilans := shared.StudioActionBalance{Skipped: []shared.StudioSkippedItem{}}
	if len(brakujace) > 0 {
		bilans.Skipped = append(bilans.Skipped, shared.StudioSkippedItem{
			Reason: "pola wymagane zostały niewypełnione — ich znaczniki są w dokumencie " +
				"WIDOCZNE, a nie usunięte",
			Detail: wejscieWskaznikTekstu(strings.Join(brakujace, ", ")),
		})
	}

	postac := shared.StudioDocumentForm{}
	if szablon.Form != nil {
		postac = *szablon.Form
	}

	dokument, err := a.wejscieDokumentAlboNowy(ctx, z.DocumentId,
		strings.TrimSpace(wartoscTekstu(z.WindowId)), szablonTytulDokumentu(z, szablon),
		szablon.Format)
	if err != nil {
		return shared.StudioTemplateFillResponse{}, err
	}

	stan := &stanPostaci{dokument: dokument}
	if strings.TrimSpace(wartoscTekstu(z.DocumentId)) != "" {
		// Wypełnienie dokumentu istniejącego zachowuje jego postać — podstawiają się tylko wartości pól.
		wczytany, err := a.postacWczytaj(ctx, dokument.Kod)
		if err != nil {
			return shared.StudioTemplateFillResponse{}, err
		}
		stan = wczytany
		bilans.Applied = szablonPodstawWPostaci(&stan.forma, wartosci, brakujace)
	} else {
		postac.DocumentId = dokument.Kod
		wejsciePrzepiszIdentyfikatory(&postac)
		stan.forma = postac
		bilans.Applied = szablonPodstawWPostaci(&stan.forma, wartosci, brakujace)
		if z.Title != nil && strings.TrimSpace(*z.Title) != "" {
			nazwa := strings.TrimSpace(*z.Title)
			stan.dokument.Tytul = &nazwa
		}
	}
	if bilans.Applied == 0 && len(wartosci) > 0 {
		bilans.Skipped = append(bilans.Skipped, shared.StudioSkippedItem{
			Reason: "żadna wartość nie miała gdzie wejść — treść nie niosła znaczników pól",
			Detail: wejscieWskaznikTekstu("podano " + strconv.Itoa(len(wartosci)) +
				" wartości; znacznik pola ma postać " + szablonOtwarciePola + "nazwa" +
				szablonZamkniecePola),
		})
	}
	// Wypełnienie szablonu jest czynnością, więc bilans nie może wyjść zerowy — dokument powstał.
	if bilans.Applied == 0 {
		bilans.Applied = 1
	}
	bilans.SkippedCount = len(bilans.Skipped)
	bilans.Note = wejscieWskaznikTekstu("dokument z szablonu „" + szablon.Name +
		"” gotowy; podstawiono " + strconv.Itoa(len(wartosci)) + " wartości, pól wymaganych " +
		"bez wartości: " + strconv.Itoa(len(brakujace)))

	stan.opisCzynnosci = "wypełnienie szablonu „" + szablon.Name + "”"
	if err := a.wejscieUtrwalPostac(ctx, stan); err != nil {
		return shared.StudioTemplateFillResponse{}, err
	}
	if strings.TrimSpace(wartoscTekstu(z.DocumentId)) == "" {
		if _, err := a.blokadaPrzenieSzablon(ctx, stan.dokument.ID, szablon.Id); err != nil {
			return shared.StudioTemplateFillResponse{}, err
		}
	}
	if err := a.wejscieOdlozCzynnosc(ctx, stan, autor); err != nil {
		return shared.StudioTemplateFillResponse{}, err
	}

	odpowiedz := shared.StudioTemplateFillResponse{
		Document: a.zlozDokument(stan.dokument),
		Form:     stan.forma,
		Balance:  bilans,
	}
	if len(brakujace) > 0 {
		odpowiedz.MissingRequired = brakujace
	}
	return odpowiedz, nil
}

// szablonTytulDokumentu rozstrzyga tytuł dokumentu zakładanego z szablonu: tytuł
// z żądania, a bez niego — nazwa szablonu.
func szablonTytulDokumentu(z shared.StudioTemplateFillRequest,
	szablon shared.StudioTemplateDetail) *string {

	if z.Title != nil && strings.TrimSpace(*z.Title) != "" {
		nazwa := strings.TrimSpace(*z.Title)
		return &nazwa
	}
	nazwa := szablon.Name
	return &nazwa
}

// szablonWartoscNapisem przekłada wartość pola na napis.
//
// Liczba i wartość logiczna wchodzą zapisem, a nie znikają: pole „kwota" podane
// liczbą jest normalnym żądaniem, a pominięcie go dałoby pismo bez kwoty.
func szablonWartoscNapisem(wartosc any) string {
	switch dokladna := wartosc.(type) {
	case string:
		return dokladna
	case float64:
		if dokladna == float64(int64(dokladna)) {
			return strconv.FormatInt(int64(dokladna), 10)
		}
		return strconv.FormatFloat(dokladna, 'f', -1, 64)
	case bool:
		if dokladna {
			return "tak"
		}
		return "nie"
	case nil:
		return ""
	default:
		zapis, err := json.Marshal(dokladna)
		if err != nil {
			return ""
		}
		return string(zapis)
	}
}

// szablonPodstawWPostaci podstawia wartości pól w postaci dokumentu i oddaje
// liczbę podstawień. Podstawienie idzie po fragmentach postaci, nie po napisie
// treści — zamiana w napisie zgubiłaby kroje, wcięcia i granice akapitów.
func szablonPodstawWPostaci(forma *shared.StudioDocumentForm, wartosci map[string]string,
	brakujace []string) int {

	pominac := map[string]bool{}
	for _, nazwa := range brakujace {
		pominac[nazwa] = true
	}
	zamien := func(tekst string) (string, int) {
		weszlo := 0
		for nazwa, wartosc := range wartosci {
			if pominac[nazwa] {
				continue
			}
			znacznik := szablonOtwarciePola + nazwa + szablonZamkniecePola
			ile := strings.Count(tekst, znacznik)
			if ile == 0 {
				continue
			}
			tekst = strings.ReplaceAll(tekst, znacznik, wartosc)
			weszlo += ile
		}
		return tekst, weszlo
	}

	weszlo := 0
	for i := range forma.Blocks {
		for j := range forma.Blocks[i].Runs {
			nowy, ile := zamien(forma.Blocks[i].Runs[j].Text)
			forma.Blocks[i].Runs[j].Text = nowy
			weszlo += ile
		}
	}
	for i := range forma.Tables {
		for j := range forma.Tables[i].Cells {
			if forma.Tables[i].Cells[j].Text == nil {
				continue
			}
			nowy, ile := zamien(*forma.Tables[i].Cells[j].Text)
			forma.Tables[i].Cells[j].Text = &nowy
			weszlo += ile
		}
	}
	// Nagłówek i stopka też niosą pola — sygnatura pisma stoi zwykle właśnie
	// tam, a nie w treści.
	for i := range forma.Sections {
		for j := range forma.Sections[i].HeadersFooters {
			wpis := &forma.Sections[i].HeadersFooters[j]
			if wpis.HeaderText != nil {
				nowy, ile := zamien(*wpis.HeaderText)
				wpis.HeaderText = &nowy
				weszlo += ile
			}
			if wpis.FooterText != nil {
				nowy, ile := zamien(*wpis.FooterText)
				wpis.FooterText = &nowy
				weszlo += ile
			}
		}
	}
	if forma.PageSetup != nil {
		if forma.PageSetup.Header != nil {
			nowy, ile := zamien(*forma.PageSetup.Header)
			forma.PageSetup.Header = &nowy
			weszlo += ile
		}
		if forma.PageSetup.Footer != nil {
			nowy, ile := zamien(*forma.PageSetup.Footer)
			forma.PageSetup.Footer = &nowy
			weszlo += ile
		}
	}
	postacPrzeliczZakresy(forma)
	return weszlo
}
