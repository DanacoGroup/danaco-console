// Odpowiedzialność pliku: podstawa odcinka wejścia, wyjścia i szablonów modułu
// Studio — kontrakt warstwy danych (kolumny warsztatu szablonów z migracji 367
// i tabela pochodzenia fragmentów z 368), JEDNA droga do bajtów materiału
// wskazanego przez Operatora, rozpoznanie formatu pliku oraz utrwalenie POSTACI
// dokumentu wraz z tym, co warstwa danych trzyma wierszami.
//
// ── Postać zapisuje obszar postaci, nie ten odcinek ─────────────────────────
// Droga zapisu postaci stoi w `adapter_modul_studio_postac.go`
// (`postacWczytaj`, `postacZapisz`, `postacZakoncz`) i ten odcinek ją WOŁA.
// Dokłada do niej jedno: wiersze, których `postacZapisz` świadomie z drzewa
// wycina — arkusz stylów, sekcje, obiekty i pola. Wycina je, bo dla dokumentu
// już istniejącego one w bazie stoją i drzewo nie ma być ich drugą prawdą.
// Dokument WNOSZONY z pliku albo KOPIOWANY jest przypadkiem odwrotnym: wierszy
// jeszcze nie ma, a postać przyszła z zewnątrz. Gdyby ten odcinek zapisał samo
// drzewo, Operator dostałby dokument, który po ponownym wczytaniu traci arkusz
// stylów wniesiony z `.docx` — czyli dokładnie tę cichą stratę, której zlecenie
// zakazuje.
//
// ── Dlaczego własny, węższy kontrakt danych ─────────────────────────────────
// `dane.RepozytoriumStudia` deklaruje `dane/studio.go` w całości i dopisanie tam
// metod byłoby wejściem w plik cudzego odcinka. Rdzeń bierze więc dokładnie te
// metody, których używa, rzutowaniem DWUWARTOŚCIOWYM — brak nazywa się wprost,
// zamiast wywracać montaż rdzenia. Ten sam wzór trzyma odcinek kontroli pracy.
package core

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"time"
	"unicode/utf8"

	"danacoconsole/server/internal/dane"
	"danacoconsole/server/internal/protocol"
	"danacoconsole/shared"
)

// WarsztatWejsciaStudia jest kontraktem warstwy danych tego odcinka: warsztat
// szablonów pism (kolumny migracji 367) i pochodzenie fragmentów (tabela z 368).
type WarsztatWejsciaStudia interface {
	dane.WarsztatSzablonowStudia
	dane.PochodzenieWejsciaStudia
}

// wejscieSkladnica zdejmuje kontrakt odcinka z repozytorium Studia.
func (a *adapterStudia) wejscieSkladnica() (WarsztatWejsciaStudia, error) {
	if a == nil || a.repozytorium == nil {
		return nil, wejscieBladZaplecza("repozytorium Studia nie zostało podane przy montażu rdzenia")
	}
	skladnica, jest := a.repozytorium.(WarsztatWejsciaStudia)
	if !jest {
		return nil, wejscieBladZaplecza("repozytorium Studia nie niesie warsztatu szablonów " +
			"z migracji 367 ani pochodzenia fragmentów z migracji 368")
	}
	return skladnica, nil
}

// ── Odmowy odcinka ──────────────────────────────────────────────────────────

// wejscieBladZaplecza nazywa niepełny montaż rdzenia — usterkę wdrożenia
// serwera, nie brak funkcji produktu.
func wejscieBladZaplecza(powod string) error {
	return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeInternalError,
		"moduł Studio, wejście i wydanie: "+powod))
}

// wejscieBladBraku nazywa byt, którego nie ma.
func wejscieBladBraku(powod string) error {
	return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeNotFound,
		"moduł Studio: "+powod))
}

// ── Jedna droga do bajtów ───────────────────────────────────────────────────

// granicaPlikuWejsciowego chroni rdzeń przed plikiem, którego rozbiór wypełniłby
// pamięć procesu. Pismo Operatora bywa obszerne; sto megabajtów mieści dokument
// z setkami obrazów.
const granicaPlikuWejsciowego = 128 << 20

// wejscieWskazanieZrodla zbiera cztery drogi, którymi kontrakt podaje materiał:
// ścieżkę pliku, zasób magazynu rdzenia, plik Biblioteki i bajty wprost.
type wejscieWskazanieZrodla struct {
	Sciezka        *string
	ZasobKod       *string
	PlikBiblioteki *string
	BajtyBase64    *string
	// Czynnosc nazywa czynność w odmowie — Operator ma wiedzieć, czego brakło
	// przy której komendzie.
	Czynnosc string
}

// wejscieBajtyZrodla oddaje bajty materiału wraz z nazwą pliku, którą materiał
// przyniósł. Nazwa służy rozpoznaniu formatu i nazwaniu dokumentu — plik
// `umowa najmu.docx` ma zostać dokumentem „umowa najmu", nie dokumentem bez nazwy.
//
// Kolejność dróg jest rozstrzygnięciem, nie przypadkiem: bajty wprost są
// najpewniejsze (Operator je właśnie wysłał), potem zasób rdzenia, potem plik
// Biblioteki, a ścieżka NA KOŃCU, bo jest ścieżką na maszynie serwera.
func (a *adapterStudia) wejscieBajtyZrodla(ctx context.Context,
	wskazanie wejscieWskazanieZrodla) ([]byte, string, error) {

	czynnosc := wskazanie.Czynnosc
	if czynnosc == "" {
		czynnosc = "wniesienie pliku"
	}

	if zapis := strings.TrimSpace(wartoscTekstu(wskazanie.BajtyBase64)); zapis != "" {
		bajty, err := base64.StdEncoding.DecodeString(zapis)
		if err != nil {
			return nil, "", bladWskazaniaStudio(czynnosc +
				": bajty pliku nie są poprawnym zapisem base64: " + err.Error())
		}
		if len(bajty) == 0 {
			return nil, "", bladWskazaniaStudio(czynnosc + ": plik przyszedł pusty")
		}
		if len(bajty) > granicaPlikuWejsciowego {
			return nil, "", bladWskazaniaStudio(czynnosc + ": plik jest większy niż granica " +
				"wniesienia (128 MB)")
		}
		return bajty, strings.TrimSpace(wartoscTekstu(wskazanie.Sciezka)), nil
	}

	if kod := strings.TrimSpace(wartoscTekstu(wskazanie.ZasobKod)); kod != "" {
		bajty, err := a.bajtyZasobuStudia(ctx, kod)
		if err != nil {
			return nil, "", err
		}
		return bajty, kod, nil
	}

	if kod := strings.TrimSpace(wartoscTekstu(wskazanie.PlikBiblioteki)); kod != "" {
		if a.biblioteka == nil {
			return nil, "", wejscieBladZaplecza(czynnosc + " z Biblioteki: rdzeń złożony bez " +
				"repozytorium Library; naprawa: podpiąć je przy składaniu rdzenia")
		}
		plik, err := a.biblioteka.Plik(ctx, kod)
		if err != nil {
			return nil, "", wejscieBladBraku("plik Biblioteki nie istnieje: " + kod)
		}
		if plik.TrescOdwolanie == nil || strings.TrimSpace(*plik.TrescOdwolanie) == "" {
			return nil, "", wejscieBladBraku("plik Biblioteki " + kod +
				" nie ma odłożonej treści — nie ma czego wnieść")
		}
		bajty, err := os.ReadFile(*plik.TrescOdwolanie)
		if err != nil {
			return nil, "", wejscieBladZaplecza("treści pliku Biblioteki " + kod +
				" nie da się odczytać: " + err.Error())
		}
		return bajty, plik.Nazwa, nil
	}

	if sciezka := strings.TrimSpace(wartoscTekstu(wskazanie.Sciezka)); sciezka != "" {
		bajty, err := os.ReadFile(sciezka)
		if err != nil {
			return nil, "", wejscieBladBraku(czynnosc + ": pliku " + sciezka +
				" nie da się odczytać na maszynie serwera: " + err.Error())
		}
		if len(bajty) > granicaPlikuWejsciowego {
			return nil, "", bladWskazaniaStudio(czynnosc + ": plik " + sciezka +
				" jest większy niż granica wniesienia (128 MB)")
		}
		return bajty, filepath.Base(sciezka), nil
	}

	return nil, "", bladWskazaniaStudio(czynnosc + " bez wskazania materiału — żądanie nie " +
		"niosło ani bajtów, ani zasobu, ani pliku Biblioteki, ani ścieżki")
}

// ── Rozpoznanie formatu ─────────────────────────────────────────────────────

// wejscieRozpoznajFormat rozstrzyga format pliku. Wskazanie Operatora ma
// pierwszeństwo; bez niego rozstrzyga ZAWARTOŚĆ, a rozszerzenie nazwy dopiero
// na końcu — plik przemianowany jest zwykłą rzeczą w archiwum Operatora,
// a wniesienie go wedle rozszerzenia dałoby rozbiór na ślepo.
func wejscieRozpoznajFormat(nazwa string, bajty []byte,
	wskazanie *shared.StudioImportFormat) (shared.StudioImportFormat, error) {

	if wskazanie != nil && strings.TrimSpace(string(*wskazanie)) != "" {
		for _, znany := range shared.WartosciStudioImportFormat() {
			if *wskazanie == znany {
				return *wskazanie, nil
			}
		}
		return "", bladWskazaniaStudio("format pliku " + string(*wskazanie) +
			" nie jest formatem, który rdzeń wnosi; rdzeń wnosi docx, dotx, odt, ott, " +
			"txt, md, rtf, html i pdf")
	}

	rozszerzenie := strings.ToLower(strings.TrimPrefix(filepath.Ext(nazwa), "."))

	switch {
	case len(bajty) >= 5 && string(bajty[:5]) == "%PDF-":
		return shared.StudioImportFormatPdf, nil
	case len(bajty) >= 4 && bajty[0] == 'P' && bajty[1] == 'K' &&
		(bajty[2] == 3 || bajty[2] == 5 || bajty[2] == 7):
		// Archiwum ZIP: rozstrzyga jego zawartość, bo `.docx`, `.dotx`, `.odt`
		// i `.ott` mają ten sam nagłówek.
		return wejscieFormatArchiwum(bajty, rozszerzenie)
	case len(bajty) >= 5 && strings.HasPrefix(string(bajty[:5]), "{\\rtf"):
		return shared.StudioImportFormatRtf, nil
	}

	switch rozszerzenie {
	case "docx":
		return shared.StudioImportFormatDocx, nil
	case "dotx":
		return shared.StudioImportFormatDotx, nil
	case "odt":
		return shared.StudioImportFormatOdt, nil
	case "ott":
		return shared.StudioImportFormatOtt, nil
	case "rtf":
		return shared.StudioImportFormatRtf, nil
	case "pdf":
		return shared.StudioImportFormatPdf, nil
	case "md", "markdown":
		return shared.StudioImportFormatMd, nil
	case "html", "htm", "xhtml":
		return shared.StudioImportFormatHtml, nil
	case "txt", "text", "":
		// Rozstrzygnięcie po treści: znacznik HTML na początku pliku bez
		// rozszerzenia jest HTML-em, a nie tekstem ze znacznikami w środku.
		if wejscieWygladaNaHtml(bajty) {
			return shared.StudioImportFormatHtml, nil
		}
		return shared.StudioImportFormatTxt, nil
	}
	if wejscieWygladaNaHtml(bajty) {
		return shared.StudioImportFormatHtml, nil
	}
	if utf8.Valid(bajty) {
		return shared.StudioImportFormatTxt, nil
	}
	return "", bladWskazaniaStudio("formatu pliku " + nazwa + " nie da się rozpoznać ani po " +
		"zawartości, ani po nazwie; naprawa: wskazać format polem żądania")
}

// wejscieFormatArchiwum rozstrzyga, którym archiwum biurowym jest plik ZIP.
func wejscieFormatArchiwum(bajty []byte,
	rozszerzenie string) (shared.StudioImportFormat, error) {

	skladniki, err := wejscieOtworzArchiwum(bajty)
	if err != nil {
		return "", err
	}
	_, maOoxml := skladniki[ooxmlSkladnikDokumentu]
	_, maOdf := skladniki[odfSkladnikTresci]

	switch {
	case maOoxml && rozszerzenie == "dotx":
		return shared.StudioImportFormatDotx, nil
	case maOoxml:
		return shared.StudioImportFormatDocx, nil
	case maOdf && rozszerzenie == "ott":
		return shared.StudioImportFormatOtt, nil
	case maOdf:
		// Rodzaj dokumentu ODF stoi w `mimetype` — szablon niesie
		// `…text-template`. Rozszerzenie nazwy tu nie rozstrzyga, bo plik
		// szablonu bywa przemianowany.
		if rodzaj := strings.TrimSpace(string(skladniki[odfSkladnikRodzaju])); rodzaj != "" {
			if strings.Contains(rodzaj, "text-template") {
				return shared.StudioImportFormatOtt, nil
			}
		}
		return shared.StudioImportFormatOdt, nil
	}
	return "", bladWskazaniaStudio("plik jest archiwum ZIP, ale nie niesie ani " +
		ooxmlSkladnikDokumentu + " (Word), ani " + odfSkladnikTresci +
		" (OpenDocument) — rdzeń nie ma czego z niego wnieść")
}

// wejscieWygladaNaHtml sprawdza, czy treść zaczyna się znacznikiem dokumentu.
func wejscieWygladaNaHtml(bajty []byte) bool {
	poczatek := bajty
	if len(poczatek) > 512 {
		poczatek = poczatek[:512]
	}
	nizej := strings.ToLower(strings.TrimSpace(string(poczatek)))
	return strings.HasPrefix(nizej, "<!doctype html") || strings.HasPrefix(nizej, "<html") ||
		strings.HasPrefix(nizej, "<?xml") && strings.Contains(nizej, "<html")
}

// wejscieCzyFormatSzablonu mówi, czy plik jest plikiem SZABLONU, a nie dokumentu.
func wejscieCzyFormatSzablonu(format shared.StudioImportFormat) bool {
	return format == shared.StudioImportFormatDotx || format == shared.StudioImportFormatOtt
}

// ── Utrwalenie postaci wraz z wierszami ─────────────────────────────────────

// wejscieUtrwalPostac zapisuje postać dokumentu WRAZ z tym, co warstwa danych
// trzyma wierszami: arkuszem stylów, sekcjami, obiektami i polami. Drzewo
// i treść idą drogą obszaru postaci (`postacZapisz`), a nie drugą własną.
//
// Styl fabryczny dokumentu zapisuje się jako fabryczny wtedy i tylko wtedy, gdy
// tak przyszedł: arkusz wniesiony z `.docx` Operatora jest jego arkuszem, nie
// arkuszem platformy, i ma dać się zmienić bez odmowy „styl fabryczny".
func (a *adapterStudia) wejscieUtrwalPostac(ctx context.Context, stan *stanPostaci) error {
	skladnica, err := a.postacSkladnica()
	if err != nil {
		return err
	}

	for _, styl := range stan.forma.Styles {
		wiersz, err := postacStylDoWiersza(stan.dokument.ID, styl)
		if err != nil {
			return err
		}
		if _, err := skladnica.ZapiszStylNazwany(ctx, wiersz); err != nil {
			return bladStudio(err)
		}
	}

	for kolejnosc, sekcja := range stan.forma.Sections {
		wiersz, err := wejscieSekcjaDoWiersza(stan.dokument.ID, kolejnosc, sekcja)
		if err != nil {
			return err
		}
		if _, err := skladnica.ZapiszSekcje(ctx, wiersz); err != nil {
			return bladStudio(err)
		}
	}

	for _, obiekt := range stan.forma.Objects {
		wiersz, err := wejscieObiektDoWiersza(stan.dokument.ID, obiekt)
		if err != nil {
			return err
		}
		if _, err := skladnica.ZapiszObiektDokumentu(ctx, wiersz); err != nil {
			return bladStudio(err)
		}
	}

	for _, pole := range stan.forma.Fields {
		if _, err := skladnica.ZapiszPoleDokumentu(ctx,
			wejsciePoleDoWiersza(stan.dokument.ID, pole)); err != nil {
			return bladStudio(err)
		}
	}

	return a.postacZapisz(ctx, stan)
}

// wejscieSekcjaDoWiersza składa wiersz sekcji z sekcji kontraktu.
func wejscieSekcjaDoWiersza(dokumentID int64, kolejnosc int,
	sekcja shared.StudioSection) (dane.SekcjaDokumentuStudia, error) {

	wiersz := dane.SekcjaDokumentuStudia{
		Kod:        sekcja.Id,
		DokumentID: dokumentID,
		Kolejnosc:  int64(kolejnosc),
		Tytul:      sekcja.Title,
		ZakresOd:   int64(sekcja.RangeStart),
		ZakresDo:   int64(sekcja.RangeEnd),
	}
	if wiersz.Kod == "" {
		wiersz.Kod = nowyIdentyfikator(przedrostekSekcjiStudia)
	}
	if sekcja.Start != nil && strings.TrimSpace(string(*sekcja.Start)) != "" {
		wiersz.Rozpoczecie = string(*sekcja.Start)
	} else {
		wiersz.Rozpoczecie = string(shared.StudioSectionStartContinuous)
	}
	for _, zapis := range []struct {
		wartosc any
		cel     **string
		nazwa   string
	}{
		{sekcja.PageSetup, &wiersz.NastawyStronyJSON, "nastawy strony sekcji"},
		{sekcja.Numbering, &wiersz.NumeracjaJSON, "numeracja stron sekcji"},
		{sekcja.Watermark, &wiersz.ZnakWodnyJSON, "znak wodny sekcji"},
	} {
		zapisany, err := wejscieZapisNieobowiazkowy(zapis.wartosc, zapis.nazwa)
		if err != nil {
			return dane.SekcjaDokumentuStudia{}, err
		}
		*zapis.cel = zapisany
	}
	if len(sekcja.HeadersFooters) > 0 {
		zapis, err := json.Marshal(sekcja.HeadersFooters)
		if err != nil {
			return dane.SekcjaDokumentuStudia{}, wejscieBladZaplecza(
				"nagłówka i stopki sekcji nie da się zapisać: " + err.Error())
		}
		tekst := string(zapis)
		wiersz.NaglowkiJSON = &tekst
	}
	return wiersz, nil
}

// wejscieObiektDoWiersza składa wiersz obiektu osadzonego z obiektu kontraktu.
//
// Postać obiektu (obramowanie, opływanie, przycięcie, obrót) idzie ładunkiem
// JSON, bo kolumn na nią nie ma — i tak samo ją czyta `postacZlozObiekty`.
func wejscieObiektDoWiersza(dokumentID int64,
	obiekt shared.StudioDocumentObject) (dane.ObiektDokumentuStudia, error) {

	wiersz := dane.ObiektDokumentuStudia{
		Kod:               obiekt.Id,
		DokumentID:        dokumentID,
		Rodzaj:            string(obiekt.Kind),
		ZasobKod:          obiekt.AssetId,
		DesignWezelKod:    obiekt.DesignNodeId,
		BibliotekaPlikKod: obiekt.LibraryFileId,
		AdresZrodla:       obiekt.SourceUrl,
		TekstZastepczy:    obiekt.AltText,
		Zakotwiczenie:     string(shared.StudioAnchorKindCharacter),
	}
	if wiersz.Kod == "" {
		wiersz.Kod = nowyIdentyfikator(przedrostekObiektuStudia)
	}
	if obiekt.Source != nil && strings.TrimSpace(string(*obiekt.Source)) != "" {
		zrodlo := string(*obiekt.Source)
		wiersz.Zrodlo = &zrodlo
	}
	if obiekt.Anchor != nil && strings.TrimSpace(string(*obiekt.Anchor)) != "" {
		wiersz.Zakotwiczenie = string(*obiekt.Anchor)
	}
	if obiekt.AnchorOffset != nil {
		wiersz.ZakotwiczeniePozycja = int64(*obiekt.AnchorOffset)
	}
	if obiekt.ZOrder != nil {
		wiersz.Warstwa = int64(*obiekt.ZOrder)
	}
	zapis, err := wejscieZapisNieobowiazkowy(obiekt, "postać obiektu osadzonego")
	if err != nil {
		return dane.ObiektDokumentuStudia{}, err
	}
	wiersz.PostacJSON = zapis
	return wiersz, nil
}

// wejsciePoleDoWiersza składa wiersz pola dokumentu z pola kontraktu.
func wejsciePoleDoWiersza(dokumentID int64,
	pole shared.StudioDocumentField) dane.PoleDokumentuStudia {

	wiersz := dane.PoleDokumentuStudia{
		Kod:              pole.Id,
		DokumentID:       dokumentID,
		Rodzaj:           string(pole.Kind),
		Format:           pole.Format,
		Wyrazenie:        pole.Expression,
		NazwaWlasciwosci: pole.PropertyName,
		Wartosc:          pole.Value,
	}
	if wiersz.Kod == "" {
		wiersz.Kod = nowyIdentyfikator(przedrostekPolaPostaci)
	}
	if pole.AnchorOffset != nil {
		wiersz.Kotwica = int64(*pole.AnchorOffset)
	}
	if pole.Stale != nil {
		wiersz.Nieswieze = *pole.Stale
	}
	return wiersz
}

// wejscieZapisNieobowiazkowy składa ładunek JSON pola nieobowiązkowego. Wartość
// pusta nie zapisuje się jako `null` w kolumnie tekstowej — kolumna zostaje
// pusta, bo `null` w niej znaczyłby „zapisano brak", a nie „nic nie zapisano".
func wejscieZapisNieobowiazkowy(wartosc any, nazwa string) (*string, error) {
	if wartosc == nil {
		return nil, nil
	}
	zapis, err := json.Marshal(wartosc)
	if err != nil {
		return nil, wejscieBladZaplecza(nazwa + " nie da się zapisać: " + err.Error())
	}
	tekst := string(zapis)
	if tekst == "null" || tekst == "{}" {
		return nil, nil
	}
	return &tekst, nil
}

// ── Pochodzenie wniesionego fragmentu ───────────────────────────────────────

// wejscieOdlozPochodzenie utrwala zapis, skąd fragment przyszedł, i oddaje go
// w kształcie kontraktu.
//
// Zapis pochodzenia jest obowiązkowy przy każdym wniesieniu z zewnątrz: bez
// niego za tydzień nikt nie odtworzy, na czym pismo się opiera, a powołanie
// bibliograficzne dopisywane z ręki po tygodniu jest zgadywaniem.
func (a *adapterStudia) wejscieOdlozPochodzenie(ctx context.Context,
	dokument dane.DokumentStudia, rodzaj shared.StudioProvenanceKind, od, do int,
	adres, plikBiblioteki, tytul *string, autor shared.StudioAuthor,
	agentKod, agentNazwa *string) (*shared.StudioProvenance, error) {

	skladnica, err := a.wejscieSkladnica()
	if err != nil {
		return nil, err
	}
	siegnieto := wejscieZnacznikChwili()
	wiersz, err := skladnica.ZapiszPochodzenieFragmentu(ctx, dane.PochodzenieFragmentuStudia{
		Kod:               nowyIdentyfikator(przedrostekPochodzeniaStudia),
		DokumentID:        dokument.ID,
		Rodzaj:            string(rodzaj),
		ZakresOd:          int64(od),
		ZakresDo:          int64(do),
		AdresZrodla:       adres,
		BibliotekaPlikKod: plikBiblioteki,
		TytulZrodla:       tytul,
		Siegnieto:         &siegnieto,
		AutorRodzaj:       string(autor),
		AutorAgentKod:     agentKod,
		AutorAgentNazwa:   agentNazwa,
	})
	if err != nil {
		return nil, bladStudio(err)
	}
	pochodzenie := shared.StudioProvenance{
		Id:            wiersz.Kod,
		DocumentId:    dokument.Kod,
		Kind:          shared.StudioProvenanceKind(wiersz.Rodzaj),
		RangeStart:    int(wiersz.ZakresOd),
		RangeEnd:      int(wiersz.ZakresDo),
		SourceUrl:     wiersz.AdresZrodla,
		LibraryFileId: wiersz.BibliotekaPlikKod,
		SourceVersion: wiersz.WersjaZrodla,
		SourceTitle:   wiersz.TytulZrodla,
		RetrievedAt:   wejscieWskaznikDlugi(chwilaBazy(wiersz.Utworzono)),
	}
	rodzajAutora := shared.StudioAuthor(wiersz.AutorRodzaj)
	pochodzenie.Author = &rodzajAutora
	return &pochodzenie, nil
}

// wejscieZnacznikChwili oddaje czas w zapisie, którym baza trzyma znaczniki —
// tym samym, który czyta `chwilaBazy`. Zapis niezgodny z tym formatem odczytałby
// się jako zero, czyli jako „czas nieznany".
func wejscieZnacznikChwili() string {
	return time.Now().UTC().Format(formatZnacznikaBazy)
}

// ── Dokument, na którym czynność pracuje ────────────────────────────────────

// wejscieDokumentAlboNowy wczytuje dokument wskazany albo zakłada nowy w oknie.
// Wniesienie pliku do edytora nie wymaga dokumentu wcześniejszego — plik ma
// stanąć w edytorze wprost, a nie po dwóch komendach.
func (a *adapterStudia) wejscieDokumentAlboNowy(ctx context.Context, kodDokumentu *string,
	okno string, tytul *string, format shared.StudioDocumentFormat) (dane.DokumentStudia, error) {

	if kod := strings.TrimSpace(wartoscTekstu(kodDokumentu)); kod != "" {
		dokument, err := a.repozytorium.Dokument(ctx, kod)
		if err != nil {
			return dane.DokumentStudia{}, bladNieznanegoDokumentu(kod, err)
		}
		return dokument, nil
	}
	if strings.TrimSpace(okno) == "" {
		return dane.DokumentStudia{}, bladWskazaniaStudio(
			"założenie dokumentu bez okna, w którym ma stanąć")
	}
	pusta := ""
	nowy, err := a.repozytorium.ZapiszDokument(ctx, dane.DokumentStudia{
		Kod:    nowyIdentyfikator(przedrostekDokumentuStudio),
		Okno:   strings.TrimSpace(okno),
		Tytul:  tytul,
		Format: format,
		Tresc:  &pusta,
	})
	if err != nil {
		return dane.DokumentStudia{}, bladStudio(err)
	}
	return nowy, nil
}

// wejscieFormatDokumentu przekłada format pliku wniesionego na format dokumentu
// platformy. Wykaz formatów dokumentu jest krótszy niż wykaz formatów pliku —
// `.docx` i `.dotx` są dokumentem `docx`, markdown i HTML są dokumentem
// `markdown`, a reszta tekstem.
func wejscieFormatDokumentu(format shared.StudioImportFormat) shared.StudioDocumentFormat {
	switch format {
	case shared.StudioImportFormatDocx, shared.StudioImportFormatDotx,
		shared.StudioImportFormatOdt, shared.StudioImportFormatOtt:
		return shared.StudioDocumentFormatDocx
	case shared.StudioImportFormatMd, shared.StudioImportFormatHtml:
		return shared.StudioDocumentFormatMarkdown
	case shared.StudioImportFormatPdf:
		return shared.StudioDocumentFormatTxt
	default:
		return shared.StudioDocumentFormatTxt
	}
}

// wejscieNazwaZPliku wyjmuje nazwę dokumentu z nazwy pliku — bez rozszerzenia,
// bo Operator nazwał pismo „umowa najmu", a nie „umowa najmu.docx".
func wejscieNazwaZPliku(nazwa string) *string {
	czysta := strings.TrimSpace(filepath.Base(strings.TrimSpace(nazwa)))
	if czysta == "" || czysta == "." || czysta == string(filepath.Separator) {
		return nil
	}
	if rozszerzenie := filepath.Ext(czysta); rozszerzenie != "" {
		czysta = strings.TrimSuffix(czysta, rozszerzenie)
	}
	if czysta == "" {
		return nil
	}
	return &czysta
}

// wejscieBrakWiersza mówi, czy błąd warstwy danych jest brakiem wiersza.
func wejscieBrakWiersza(err error) bool {
	return errors.Is(err, dane.ErrBrakWiersza)
}
