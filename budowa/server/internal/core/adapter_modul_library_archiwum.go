// Pakiet obsługuje rodzinę komend `library.*`: utrwalenie archiwalne zasobu
// w postaci PDF/A, BagIt albo PREMIS/METS, politykę retencji zasobów wraz
// z jej rejestrem oraz wywóz paczki archiwum z repozytorium biblioteki.
package core

import (
	"archive/tar"
	"archive/zip"
	"bytes"
	"compress/gzip"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/pdfcpu/pdfcpu/pkg/api"

	"danacoconsole/server/internal/dane"
	"danacoconsole/server/internal/konfig"
	"danacoconsole/server/internal/session"
	"danacoconsole/server/internal/zewnetrzne"
	"danacoconsole/shared"
)

// profilPdfaDomyslny jest odmianą profilu PDF/A przyjmowaną do utrwalenia
// archiwalnego, gdy żądanie nie podaje własnej odmiany profilu wprost.
const profilPdfaDomyslny = "PDF/A-2b"

// granicaPakowaniaArsenalem jest granicą czasu pakowania formatem 7z, które
// idzie procesem zewnętrznym arsenału serwerowego, nie kodem wkompilowanym.
const granicaPakowaniaArsenalem = 5 * time.Minute

// narzedziePakowaniaBiblioteki opisuje binarium arsenału serwerowego, którego
// wywołanie moduł Library zleca wyłącznie przy pakowaniu formatem 7z.
var narzedziePakowaniaBiblioteki = zewnetrzne.Narzedzie{
	Nazwa: "7-Zip", Program: "7z", Pakiet: "p7zip-full",
}

// UtrwalArchiwalnie obsługuje `library.preservation.run`: utrwala wskazane
// zasoby jedną z trzech postaci archiwalnych i oddaje wynik każdego z nich.
func (a *adapterBiblioteki) UtrwalArchiwalnie(ctx context.Context,
	z shared.LibraryPreservationRunRequest) (shared.LibraryPreservationRunResponse, error) {

	zasoby, err := a.zasobyZbioru(ctx, z.FileIds, z.CollectionId)
	if err != nil {
		return shared.LibraryPreservationRunResponse{}, err
	}
	if len(zasoby) == 0 {
		return shared.LibraryPreservationRunResponse{}, bladWskazaniaBiblioteki(
			"utrwalenie bez ani jednego zasobu — wskaż zasoby albo kolekcję")
	}
	profil := strings.TrimSpace(wartoscTekstu(z.Profile))
	nowaWersja := z.ReplaceOriginal != nil && *z.ReplaceOriginal

	wyniki := make([]shared.LibraryPreservationResult, 0, len(zasoby))
	poprawne, niepoprawne := 0, 0
	for _, zasob := range zasoby {
		wynik, err := a.utrwalZasob(ctx, zasob, z.Kind, profil, nowaWersja)
		if err != nil {
			return shared.LibraryPreservationRunResponse{}, err
		}
		if wynik.Valid {
			poprawne++
		} else {
			niepoprawne++
		}
		wyniki = append(wyniki, wynik)
	}
	return shared.LibraryPreservationRunResponse{
		Results: wyniki, ValidCount: poprawne, InvalidCount: niepoprawne,
	}, nil
}

// utrwalZasob wykonuje jedno utrwalenie wskazaną postacią i odkłada po nim
// ślad w bazie danych, wraz z zapisem albo nową wersją wytworu.
func (a *adapterBiblioteki) utrwalZasob(ctx context.Context, zasob dane.PlikBiblioteki,
	rodzaj shared.LibraryPreservationKind, profil string,
	nowaWersja bool) (shared.LibraryPreservationResult, error) {

	bajty, err := a.bajtyZasobuBiblioteki(zasob)
	if err != nil {
		return shared.LibraryPreservationResult{}, err
	}

	var wytworzone []byte
	var nazwa, raport string
	var poprawne bool

	switch rodzaj {
	case shared.LibraryPreservationKindPdfa:
		if profil == "" {
			profil = profilPdfaDomyslny
		}
		wytworzone, poprawne, raport = normalizujDokumentArchiwalnie(bajty, profil)
		nazwa = nazwaZUtrwaleniaBiblioteki(zasob.Nazwa, "pdf")

	case shared.LibraryPreservationKindBagit:
		wytworzone, poprawne, raport, err = pakietBagIt(zasob, bajty)
		if err != nil {
			return shared.LibraryPreservationResult{}, err
		}
		nazwa = nazwaZUtrwaleniaBiblioteki(zasob.Nazwa, "zip")

	case shared.LibraryPreservationKindPremis:
		wytworzone, poprawne, raport, err = opisPremisMets(zasob, bajty)
		if err != nil {
			return shared.LibraryPreservationResult{}, err
		}
		nazwa = nazwaZUtrwaleniaBiblioteki(zasob.Nazwa, "xml")

	default:
		return shared.LibraryPreservationResult{}, bladWskazaniaBiblioteki(
			"postać utrwalenia " + string(rodzaj) + " nie ma odwzorowania w serwerze")
	}

	wynik := shared.LibraryPreservationResult{
		FileId: zasob.Kod, Kind: rodzaj, Valid: poprawne, Report: raport,
		ProducedAt: time.Now().UnixMilli(),
	}

	if nowaWersja && rodzaj == shared.LibraryPreservationKindPdfa {
		// Utrwalenie „w miejsce oryginału" jest dołożeniem wersji, nie nadpisaniem.
		if err := a.dolozWersjeUtrwalenia(ctx, zasob, wytworzone, profil); err != nil {
			return shared.LibraryPreservationResult{}, err
		}
		wynik.ProducedFileId = wskazanieBiblioteki(zasob.Kod)
	} else {
		kod, err := a.odlozZasobBiblioteki(ctx, wytworzone, nazwa, zasob.ProjektID)
		if err != nil {
			return shared.LibraryPreservationResult{}, err
		}
		wynik.ProducedFileId = &kod
	}

	_, _ = a.repozytorium.ZapiszUtrwalenie(ctx, dane.ZadanieUtrwaleniaBiblioteki{
		Kod: nowyIdentyfikator(przedrostekUtrwaleniaBiblioteki), PlikKod: zasob.Kod,
		Rodzaj: rodzajUtrwaleniaBazy(rodzaj), PlikWynikowyKod: wynik.ProducedFileId,
		Profil: wskazanieBiblioteki(profil), Poprawne: poprawne, Raport: raport,
	})
	a.odnotuj(ctx, shared.LibraryAuditActionPreserve, wskazanieBiblioteki(zasob.Kod),
		"utrwalenie w postaci "+string(rodzaj))
	return wynik, nil
}

// normalizujDokumentArchiwalnie porządkuje strukturę dokumentu PDF i waliduje
// wynik strukturalnie, bez orzekania zgodności z profilem PDF/A.
func normalizujDokumentArchiwalnie(bajty []byte, profil string) ([]byte, bool, string) {
	if !bytes.HasPrefix(bajty, []byte("%PDF-")) {
		return nil, false, "materiał nie jest dokumentem PDF — normalizacja archiwalna " +
			"dotyczy dokumentów, a nie dowolnej treści"
	}
	var wynik bytes.Buffer
	if err := api.Optimize(bytes.NewReader(bajty), &wynik, nastawyPdf()); err != nil {
		return nil, false, "normalizacja nie powiodła się: " + err.Error()
	}
	raport := "sprawdzono: struktura dokumentu po normalizacji (pdfcpu), " +
		"liczba stron zachowana, docelowa odmiana profilu: " + profil + ". " +
		"NIE sprawdzono: zgodności z profilem " + profil +
		" — serwer nie niesie walidatora profilu i zgodności nie orzeka"
	if err := api.Validate(bytes.NewReader(wynik.Bytes()), nastawyPdf()); err != nil {
		return wynik.Bytes(), false, raport + ". Walidacja struktury: niepowodzenie — " + err.Error()
	}
	stronPrzed, err := api.PageCount(bytes.NewReader(bajty), nastawyPdf())
	if err == nil {
		stronPo, err := api.PageCount(bytes.NewReader(wynik.Bytes()), nastawyPdf())
		if err == nil && stronPrzed != stronPo {
			return wynik.Bytes(), false, raport + ". Liczba stron zmieniła się z " +
				strconv.Itoa(stronPrzed) + " na " + strconv.Itoa(stronPo)
		}
	}
	return wynik.Bytes(), true, raport + ". Walidacja struktury: pomyślna"
}

// pakietBagIt składa pakiet archiwalny w postaci BagIt i sprawdza jego
// manifest przeliczeniem sumy kontrolnej bajtów już spakowanych.
func pakietBagIt(zasob dane.PlikBiblioteki, bajty []byte) ([]byte, bool, string, error) {
	skrot := sha256.Sum256(bajty)
	suma := hex.EncodeToString(skrot[:])
	sciezkaWPakiecie := "data/" + zapisBezpiecznyNazwyBiblioteki(zasob.Nazwa, true)

	pozycje := []pozycjaArchiwumBiblioteki{
		{nazwa: "bagit.txt", tresc: []byte("BagIt-Version: 1.0\nTag-File-Character-Encoding: UTF-8\n")},
		{nazwa: "bag-info.txt", tresc: []byte(
			"Source-Organization: Danaco Holding Group\n" +
				"Bagging-Date: " + time.Now().UTC().Format("2006-01-02") + "\n" +
				"External-Identifier: " + zasob.Kod + "\n" +
				"Payload-Oxum: " + strconv.Itoa(len(bajty)) + ".1\n")},
		{nazwa: "manifest-sha256.txt", tresc: []byte(suma + "  " + sciezkaWPakiecie + "\n")},
		{nazwa: sciezkaWPakiecie, tresc: bajty},
	}
	pakiet, err := archiwumZipBiblioteki(pozycje)
	if err != nil {
		return nil, false, "", bladBiblioteki(err)
	}
	// Walidacja idzie po bajtach już spakowanych, nie po tych źródłowych.
	zgodny, powod := sprawdzManifestBagIt(pakiet, sciezkaWPakiecie, suma)
	raport := "sprawdzono: odczyt pakietu, obecność bagit.txt, bag-info.txt " +
		"i manifest-sha256.txt oraz zgodność sumy sha256 pozycji data/ z manifestem. " + powod
	return pakiet, zgodny, raport, nil
}

// sprawdzManifestBagIt odczytuje pakiet i porównuje sumę kontrolną treści
// pozycji `data/` z sumą zapisaną w manifeście pakietu.
func sprawdzManifestBagIt(pakiet []byte, sciezka, suma string) (bool, string) {
	czytnik, err := zip.NewReader(bytes.NewReader(pakiet), int64(len(pakiet)))
	if err != nil {
		return false, "Wynik: pakietu nie da się odczytać — " + err.Error()
	}
	for _, pozycja := range czytnik.File {
		if pozycja.Name != sciezka {
			continue
		}
		strumien, err := pozycja.Open()
		if err != nil {
			return false, "Wynik: pozycji data/ nie da się otworzyć — " + err.Error()
		}
		defer strumien.Close()
		skrot := sha256.New()
		if _, err := io.Copy(skrot, strumien); err != nil {
			return false, "Wynik: pozycji data/ nie da się odczytać"
		}
		if hex.EncodeToString(skrot.Sum(nil)) != suma {
			return false, "Wynik: suma pozycji data/ nie zgadza się z manifestem"
		}
		return true, "Wynik: zgodny"
	}
	return false, "Wynik: pakiet nie niesie pozycji " + sciezka
}

// opisPremisMets składa dokument METS z sekcją PREMIS opisującą obiekt, jego
// sumę kontrolną i zdarzenie utrwalenia — wynikiem jest opis, nie kopia zasobu.
func opisPremisMets(zasob dane.PlikBiblioteki, bajty []byte) ([]byte, bool, string, error) {
	skrot := sha256.Sum256(bajty)
	suma := hex.EncodeToString(skrot[:])
	rodzaj := ""
	if zasob.MimeType != nil {
		rodzaj = *zasob.MimeType
	}

	var zapis strings.Builder
	zapis.WriteString(xml.Header)
	zapis.WriteString(`<mets xmlns="http://www.loc.gov/METS/"` + "\n")
	zapis.WriteString(`      xmlns:premis="http://www.loc.gov/premis/v3"` + "\n")
	zapis.WriteString(`      OBJID="` + ucieczkaXmlTezaurusa(zasob.Kod) + `">` + "\n")
	zapis.WriteString("  <metsHdr CREATEDATE=\"" + time.Now().UTC().Format(time.RFC3339) +
		"\"/>\n")
	zapis.WriteString("  <amdSec>\n    <techMD>\n      <mdWrap MDTYPE=\"PREMIS:OBJECT\">\n")
	zapis.WriteString("        <xmlData>\n")
	zapis.WriteString("          <premis:object>\n")
	zapis.WriteString("            <premis:objectIdentifier>\n")
	zapis.WriteString("              <premis:objectIdentifierType>Danaco Console</premis:objectIdentifierType>\n")
	zapis.WriteString("              <premis:objectIdentifierValue>" +
		ucieczkaXmlTezaurusa(zasob.Kod) + "</premis:objectIdentifierValue>\n")
	zapis.WriteString("            </premis:objectIdentifier>\n")
	zapis.WriteString("            <premis:objectCharacteristics>\n")
	zapis.WriteString("              <premis:fixity>\n")
	zapis.WriteString("                <premis:messageDigestAlgorithm>SHA-256</premis:messageDigestAlgorithm>\n")
	zapis.WriteString("                <premis:messageDigest>" + suma + "</premis:messageDigest>\n")
	zapis.WriteString("              </premis:fixity>\n")
	zapis.WriteString("              <premis:size>" + strconv.Itoa(len(bajty)) + "</premis:size>\n")
	zapis.WriteString("              <premis:format><premis:formatDesignation>" +
		"<premis:formatName>" + ucieczkaXmlTezaurusa(rodzaj) + "</premis:formatName>" +
		"</premis:formatDesignation></premis:format>\n")
	zapis.WriteString("            </premis:objectCharacteristics>\n")
	zapis.WriteString("            <premis:originalName>" + ucieczkaXmlTezaurusa(zasob.Nazwa) +
		"</premis:originalName>\n")
	zapis.WriteString("          </premis:object>\n")
	zapis.WriteString("          <premis:event>\n")
	zapis.WriteString("            <premis:eventType>fixity check</premis:eventType>\n")
	zapis.WriteString("            <premis:eventDateTime>" + time.Now().UTC().Format(time.RFC3339) +
		"</premis:eventDateTime>\n")
	zapis.WriteString("            <premis:eventOutcomeInformation><premis:eventOutcome>" +
		"suma przeliczona z bajtów</premis:eventOutcome></premis:eventOutcomeInformation>\n")
	zapis.WriteString("          </premis:event>\n")
	zapis.WriteString("        </xmlData>\n      </mdWrap>\n    </techMD>\n  </amdSec>\n")
	zapis.WriteString("</mets>\n")

	tresc := []byte(zapis.String())
	// Walidacja opisu sprawdza, że dokument daje się odczytać jako XML.
	poprawny := true
	powod := "Wynik: zgodny"
	if err := xml.Unmarshal(tresc, new(any)); err != nil {
		poprawny, powod = false, "Wynik: opis nie jest poprawnym XML — "+err.Error()
	}
	raport := "sprawdzono: poprawność składniową dokumentu METS, obecność sekcji " +
		"PREMIS object wraz z sumą sha256 i rozmiarem oraz zdarzenia utrwalenia. " + powod
	return tresc, poprawny, raport, nil
}

// UstawRetencje obsługuje `library.retention.set`: zapisuje albo usuwa
// politykę retencji i oddaje liczbę zasobów, których polityka dotyczy.
func (a *adapterBiblioteki) UstawRetencje(ctx context.Context,
	z shared.LibraryRetentionSetRequest) (shared.LibraryRetentionSetResponse, error) {

	if z.Remove != nil && *z.Remove {
		kod := strings.TrimSpace(z.Policy.Id)
		if kod == "" {
			return shared.LibraryRetentionSetResponse{}, bladWskazaniaBiblioteki(
				"usunięcie polityki bez wskazania jej identyfikatora")
		}
		if _, err := a.repozytorium.UsunPolitykeRetencji(ctx, kod); err != nil {
			return shared.LibraryRetentionSetResponse{}, bladBiblioteki(err)
		}
		a.odnotuj(ctx, shared.LibraryAuditActionChange, nil, "usunięcie polityki retencji "+kod)
		return shared.LibraryRetentionSetResponse{Policy: z.Policy}, nil
	}
	if z.Policy.KeepDays <= 0 {
		return shared.LibraryRetentionSetResponse{}, bladWskazaniaBiblioteki(
			"okres przechowywania liczy się w dniach i musi być dodatni")
	}
	kod := strings.TrimSpace(z.Policy.Id)
	if kod == "" {
		kod = nowyIdentyfikator(przedrostekPolitykiBiblioteki)
	}
	zapisana, err := a.repozytorium.ZapiszPolitykeRetencji(ctx, dane.PolitykaRetencjiBiblioteki{
		Kod: kod, Zasieg: zasiegRetencjiBazy(z.Policy.Scope), ZasiegID: z.Policy.ScopeId,
		DniPrzechowywania: z.Policy.KeepDays, Czynnosc: czynnoscRetencjiBazy(z.Policy.Action),
	})
	if err != nil {
		return shared.LibraryRetentionSetResponse{}, bladBiblioteki(err)
	}
	// Liczba objętych zasobów jest liczbą rzeczywistą, przeliczoną, nie zapowiedzią.
	objete, err := a.zasobyObjetePolityka(ctx, zapisana)
	if err != nil {
		return shared.LibraryRetentionSetResponse{}, err
	}
	a.odnotuj(ctx, shared.LibraryAuditActionChange, nil, "zapis polityki retencji "+kod)
	return shared.LibraryRetentionSetResponse{
		Policy: politykaKontraktuBiblioteki(zapisana), AffectedFiles: len(objete),
	}, nil
}

// WykazRetencji obsługuje `library.retention.list`.
//
// Raport terminów liczy się w chwili pytania: termin zasobu wynika z jego
// ostatniej zmiany i okresu polityki, a obie wartości zmieniają się w czasie.
func (a *adapterBiblioteki) WykazRetencji(ctx context.Context,
	z shared.LibraryRetentionListRequest) (shared.LibraryRetentionListResponse, error) {

	var zasieg *string
	if z.Scope != nil {
		wartosc := zasiegRetencjiBazy(*z.Scope)
		zasieg = &wartosc
	}
	polityki, err := a.repozytorium.PolitykiRetencji(ctx, zasieg)
	if err != nil {
		return shared.LibraryRetentionListResponse{}, bladBiblioteki(err)
	}
	odpowiedz := shared.LibraryRetentionListResponse{
		Policies: make([]shared.LibraryRetentionPolicy, 0, len(polityki)),
	}
	for _, polityka := range polityki {
		odpowiedz.Policies = append(odpowiedz.Policies, politykaKontraktuBiblioteki(polityka))
	}
	if z.DueWithinDays == nil {
		return odpowiedz, nil
	}

	granica := time.Now().AddDate(0, 0, *z.DueWithinDays)
	terminy := []shared.LibraryRetentionDue{}
	for _, polityka := range polityki {
		zasoby, err := a.zasobyObjetePolityka(ctx, polityka)
		if err != nil {
			return shared.LibraryRetentionListResponse{}, err
		}
		for _, zasob := range zasoby {
			zmiana := chwilaBazy(zasob.Zaktualizowano)
			termin := time.UnixMilli(zmiana).AddDate(0, 0, polityka.DniPrzechowywania)
			if termin.After(granica) {
				continue
			}
			terminy = append(terminy, shared.LibraryRetentionDue{
				FileId: zasob.Kod, PolicyId: polityka.Kod, DueAt: termin.UnixMilli(),
				Action: czynnoscRetencjiKontraktu(polityka.Czynnosc),
			})
		}
	}
	sort.SliceStable(terminy, func(pierwszy, drugi int) bool {
		return terminy[pierwszy].DueAt < terminy[drugi].DueAt
	})
	if z.Limit != nil && *z.Limit > 0 && len(terminy) > *z.Limit {
		terminy = terminy[:*z.Limit]
	}
	odpowiedz.Due = terminy
	return odpowiedz, nil
}

// zasobyObjetePolityka wskazuje zasoby, których polityka dotyczy.
//
// Zasięg projektowy zawęża do projektu, zasięg kolekcji do kolekcji, każdy inny
// obejmuje repozytorium czynne — polityka założona globalnie ma dotyczyć
// wszystkiego i tak też jest liczona.
func (a *adapterBiblioteki) zasobyObjetePolityka(ctx context.Context,
	polityka dane.PolitykaRetencjiBiblioteki) ([]dane.PlikBiblioteki, error) {

	stan := dane.StanZasobuCzynny
	filtr := dane.FiltrPlikow{Stan: &stan, Limit: granicaZbioruHigieny}
	if polityka.ZasiegID != nil && *polityka.ZasiegID != "" {
		switch polityka.Zasieg {
		case "projekt":
			filtr.ProjektID = polityka.ZasiegID
		case "modul", "para_modulow", "srodowisko", "aplikacja", "okno", "rola", "karta_sesji":
			// Te poziomy nie zawężają zbioru zasobów: polityka obowiązuje całe repozytorium.
		default:
			filtr.KolekcjaKod = polityka.ZasiegID
		}
	}
	wiersze, _, err := a.repozytorium.Pliki(ctx, filtr)
	if err != nil {
		return nil, bladBiblioteki(err)
	}
	return wiersze, nil
}

// WywiezPaczke obsługuje `library.package.export`: składa z zasobów paczkę
// archiwum wraz z manifestem sum kontrolnych i indeksem opisującym zawartość.
func (a *adapterBiblioteki) WywiezPaczke(ctx context.Context,
	z shared.LibraryPackageExportRequest) (shared.LibraryPackageExportResponse, error) {

	kody := z.FileIds
	kolekcja := z.CollectionId
	if z.Kind == shared.LibraryPackageKindSnapshot {
		// Migawka jest zapisem stanu całego repozytorium — zawężenia żądania nie obowiązują.
		kody, kolekcja = nil, nil
	}
	zasoby, err := a.zasobyZbioru(ctx, kody, kolekcja)
	if err != nil {
		return shared.LibraryPackageExportResponse{}, err
	}
	if len(zasoby) == 0 {
		return shared.LibraryPackageExportResponse{}, bladWskazaniaBiblioteki(
			"paczka bez ani jednego zasobu — wywóz pustki nie jest wywozem")
	}

	pozycje := []pozycjaArchiwumBiblioteki{}
	manifest := strings.Builder{}
	indeks := make([]map[string]any, 0, len(zasoby))
	zWersjami := z.IncludeVersions != nil && *z.IncludeVersions

	for _, zasob := range zasoby {
		bajty, err := a.bajtyZasobuBiblioteki(zasob)
		if err != nil {
			// Zasób bez treści nie wywraca paczki: wchodzi do indeksu jako pozycja bez bajtów.
			indeks = append(indeks, opisZasobuWPaczce(zasob, "", false))
			continue
		}
		sciezka := "data/" + zasob.Kod + "-" + zapisBezpiecznyNazwyBiblioteki(zasob.Nazwa, true)
		skrot := sha256.Sum256(bajty)
		suma := hex.EncodeToString(skrot[:])
		pozycje = append(pozycje, pozycjaArchiwumBiblioteki{nazwa: sciezka, tresc: bajty})
		manifest.WriteString(suma + "  " + sciezka + "\n")
		indeks = append(indeks, opisZasobuWPaczce(zasob, sciezka, true))

		if !zWersjami {
			continue
		}
		wersje, err := a.repozytorium.Wersje(ctx, zasob.ID)
		if err != nil {
			return shared.LibraryPackageExportResponse{}, bladBiblioteki(err)
		}
		for _, wersja := range wersje {
			if wersja.TrescOdwolanie == nil || *wersja.TrescOdwolanie == "" {
				continue
			}
			trescWersji, err := os.ReadFile(*wersja.TrescOdwolanie)
			if err != nil {
				continue
			}
			sciezkaWersji := "wersje/" + zasob.Kod + "/" + wersja.Kod
			skrotWersji := sha256.Sum256(trescWersji)
			pozycje = append(pozycje, pozycjaArchiwumBiblioteki{
				nazwa: sciezkaWersji, tresc: trescWersji,
			})
			manifest.WriteString(hex.EncodeToString(skrotWersji[:]) + "  " + sciezkaWersji + "\n")
		}
	}

	opisIndeksu, err := json.MarshalIndent(map[string]any{
		"kind":      string(z.Kind),
		"createdAt": time.Now().UTC().Format(time.RFC3339),
		"entries":   indeks,
	}, "", "  ")
	if err != nil {
		return shared.LibraryPackageExportResponse{}, bladBiblioteki(err)
	}
	trescManifestu := manifest.String()
	skrotManifestu := sha256.Sum256([]byte(trescManifestu))
	sumaManifestu := hex.EncodeToString(skrotManifestu[:])

	pozycje = append(pozycje,
		pozycjaArchiwumBiblioteki{nazwa: "indeks.json", tresc: opisIndeksu},
		pozycjaArchiwumBiblioteki{nazwa: "manifest-sha256.txt", tresc: []byte(trescManifestu)},
	)

	format := strings.ToLower(strings.TrimSpace(wartoscTekstu(z.Format)))
	if format == "" {
		format = "zip"
	}
	paczka, err := a.spakuj(ctx, pozycje, format)
	if err != nil {
		return shared.LibraryPackageExportResponse{}, err
	}
	nazwa := "paczka-" + string(z.Kind) + "-" +
		time.Now().UTC().Format("20060102-150405") + "." + format
	kod, err := a.odlozZasobBiblioteki(ctx, paczka, nazwa, nil)
	if err != nil {
		return shared.LibraryPackageExportResponse{}, err
	}
	a.odnotuj(ctx, shared.LibraryAuditActionExport, wskazanieBiblioteki(kod),
		"wywóz paczki: pozycji "+strconv.Itoa(len(pozycje)))
	return shared.LibraryPackageExportResponse{
		FileId: kod, Entries: len(pozycje), SizeBytes: int64(len(paczka)),
		ManifestChecksum: sumaManifestu,
	}, nil
}

// opisZasobuWPaczce składa pozycję indeksu paczki: tożsamość zasobu, jego
// ścieżkę w archiwum, gdy zapisano treść, oraz typ i sumę kontrolną.
func opisZasobuWPaczce(zasob dane.PlikBiblioteki, sciezka string, zTrescia bool) map[string]any {
	wpis := map[string]any{
		"id": zasob.Kod, "name": zasob.Nazwa, "hasContent": zTrescia,
	}
	if sciezka != "" {
		wpis["path"] = sciezka
	}
	if zasob.MimeType != nil {
		wpis["mimeType"] = *zasob.MimeType
	}
	if zasob.SumaKontrolna != nil {
		wpis["checksum"] = *zasob.SumaKontrolna
	}
	return wpis
}

// pozycjaArchiwumBiblioteki to jeden wpis składanego archiwum: nazwa jego
// ścieżki wewnątrz archiwum wraz z treścią bajtową tego wpisu.
type pozycjaArchiwumBiblioteki struct {
	nazwa string
	tresc []byte
}

// spakuj składa archiwum we wskazanym formacie: zip, tar.gz albo 7z przez
// arsenał serwerowy, gdy wdrożenie ma program 7z wpięty.
func (a *adapterBiblioteki) spakuj(ctx context.Context, pozycje []pozycjaArchiwumBiblioteki,
	format string) ([]byte, error) {

	switch format {
	case "zip":
		bajty, err := archiwumZipBiblioteki(pozycje)
		if err != nil {
			return nil, bladBiblioteki(err)
		}
		return bajty, nil
	case "tar.gz", "targz", "tgz":
		bajty, err := archiwumTarGzBiblioteki(pozycje)
		if err != nil {
			return nil, bladBiblioteki(err)
		}
		return bajty, nil
	case "7z":
		return a.archiwum7z(ctx, pozycje)
	default:
		return nil, bladWskazaniaBiblioteki(
			"format archiwum " + format + " nie jest znany; paczka idzie w zip, tar.gz albo 7z")
	}
}

// archiwumZipBiblioteki składa archiwum ZIP w pamięci z podanej listy wpisów,
// bez zapisu na dysk pośredniego.
func archiwumZipBiblioteki(pozycje []pozycjaArchiwumBiblioteki) ([]byte, error) {
	var bufor bytes.Buffer
	zapis := zip.NewWriter(&bufor)
	for _, pozycja := range pozycje {
		wpis, err := zapis.Create(pozycja.nazwa)
		if err != nil {
			return nil, err
		}
		if _, err := wpis.Write(pozycja.tresc); err != nil {
			return nil, err
		}
	}
	if err := zapis.Close(); err != nil {
		return nil, err
	}
	return bufor.Bytes(), nil
}

// archiwumTarGzBiblioteki składa archiwum tar spakowane gzipem z podanej
// listy wpisów, bez zapisu na dysk pośredniego.
func archiwumTarGzBiblioteki(pozycje []pozycjaArchiwumBiblioteki) ([]byte, error) {
	var bufor bytes.Buffer
	spakowanie := gzip.NewWriter(&bufor)
	zapis := tar.NewWriter(spakowanie)
	for _, pozycja := range pozycje {
		naglowek := &tar.Header{
			Name: pozycja.nazwa, Mode: 0o600, Size: int64(len(pozycja.tresc)),
			ModTime: time.Now(),
		}
		if err := zapis.WriteHeader(naglowek); err != nil {
			return nil, err
		}
		if _, err := zapis.Write(pozycja.tresc); err != nil {
			return nil, err
		}
	}
	if err := zapis.Close(); err != nil {
		return nil, err
	}
	if err := spakowanie.Close(); err != nil {
		return nil, err
	}
	return bufor.Bytes(), nil
}

// archiwum7z składa archiwum formatem 7z jedyną drogą, którą rdzeń zna:
// wywołaniem programu `7z` arsenału serwerowego, bo biblioteka standardowa
// Go tego formatu nie zapisuje.
func (a *adapterBiblioteki) archiwum7z(ctx context.Context,
	pozycje []pozycjaArchiwumBiblioteki) ([]byte, error) {

	if a.uruchamiacz == nil {
		return nil, bladWskazaniaBiblioteki(
			"format 7z żąda arsenału serwerowego, którego serwer nie ma wpiętego; " +
				"paczka w postaci zip i tar.gz powstaje bez niego")
	}
	katalog, err := os.MkdirTemp("", "danaco-paczka-")
	if err != nil {
		return nil, bladBiblioteki(err)
	}
	defer os.RemoveAll(katalog)

	for _, pozycja := range pozycje {
		sciezka := filepath.Join(katalog, "zawartosc", filepath.FromSlash(pozycja.nazwa))
		if err := os.MkdirAll(filepath.Dir(sciezka), 0o700); err != nil {
			return nil, bladBiblioteki(err)
		}
		if err := os.WriteFile(sciezka, pozycja.tresc, 0o600); err != nil {
			return nil, bladBiblioteki(err)
		}
	}
	wynikowa := filepath.Join(katalog, "paczka.7z")

	okno := session.Okno{Ustawienia: session.Ustawienia{
		SrodowiskoWykonania: shared.ExecutionEnvCore,
	}}
	zasady := session.Zasady{}
	if a.rozstrzygacz != nil {
		zasady = ZasadyIzolacji(a.rozstrzygacz, konfig.Kontekst{})
	}
	obszar := session.Obszar{}
	if a.katalog != nil {
		obszar = ObszarOkna(a.katalog.Ustal(konfig.Kontekst{}, ""), "")
	}
	_, err = zewnetrzne.Wolaj(ctx, a.uruchamiacz, okno, zasady, obszar,
		narzedziePakowaniaBiblioteki,
		[]string{"a", "-t7z", wynikowa, filepath.Join(katalog, "zawartosc", "*")},
		"", granicaPakowaniaArsenalem)
	if err != nil {
		return nil, bladWskazaniaBiblioteki(
			"pakowanie formatem 7z nie powiodło się: " + err.Error() +
				"; paczka w postaci zip i tar.gz powstaje bez programu zewnętrznego")
	}
	bajty, err := os.ReadFile(wynikowa)
	if err != nil {
		return nil, bladBiblioteki(err)
	}
	return bajty, nil
}

// bajtyZasobuBiblioteki oddaje treść zasobu spod odwołania magazynu treści,
// zapisanego przy jego wersji bieżącej.
func (a *adapterBiblioteki) bajtyZasobuBiblioteki(zasob dane.PlikBiblioteki) ([]byte, error) {
	if zasob.TrescOdwolanie == nil || *zasob.TrescOdwolanie == "" {
		return nil, bladBrakuTresciBiblioteki(zasob.Kod)
	}
	bajty, err := os.ReadFile(*zasob.TrescOdwolanie)
	if err != nil {
		return nil, bladBrakuTresciBiblioteki(zasob.Kod)
	}
	return bajty, nil
}

// odlozZasobBiblioteki utrwala wytworzoną treść w magazynie i zakłada jej zasób
// repozytorium wraz z pierwszą wersją.
func (a *adapterBiblioteki) odlozZasobBiblioteki(ctx context.Context, bajty []byte,
	nazwa string, projektID *string) (string, error) {

	if a.magazyn == nil {
		return "", bladBiblioteki(fmt.Errorf("moduł Library nie ma magazynu treści"))
	}
	skrot := sha256.Sum256(bajty)
	suma := hex.EncodeToString(skrot[:])
	odwolanie, err := a.magazyn.Zapisz(bajty, suma)
	if err != nil {
		return "", bladBiblioteki(err)
	}
	rozmiar := int64(len(bajty))
	modul := "library"
	plik, err := a.repozytorium.ZapiszPlik(ctx, dane.PlikBiblioteki{
		Kod: nowyIdentyfikator(przedrostekPlikuBiblioteki), Nazwa: nazwa,
		RozmiarBajtow: &rozmiar, SumaKontrolna: &suma, TrescOdwolanie: &odwolanie,
		ProjektID: projektID, ModulZrodlowyID: &modul, Stan: dane.StanZasobuCzynny,
	})
	if err != nil {
		return "", bladBiblioteki(err)
	}
	wersja, err := a.repozytorium.ZapiszWersje(ctx, plik.ID, dane.WersjaPlikuBiblioteki{
		Kod: nowyIdentyfikator(przedrostekWersjiBiblioteki), PlikID: plik.ID,
		RozmiarBajtow: &rozmiar, SumaKontrolna: &suma, TrescOdwolanie: &odwolanie,
	})
	if err != nil {
		return "", bladBiblioteki(err)
	}
	if _, err := a.repozytorium.ZapiszPlik(ctx, plikZBiezacaWersja(plik, wersja.ID)); err != nil {
		return "", bladBiblioteki(err)
	}
	a.zaindeksujTresc(ctx, plik)
	a.zglosNasluchom(shared.LibraryWebhookEventFileAdded, plik.Kod)
	return plik.Kod, nil
}

// dolozWersjeUtrwalenia dokłada wynik utrwalenia jako nową wersję zasobu,
// zamiast zakładać zasób osobny — poprzednia treść zostaje w historii.
func (a *adapterBiblioteki) dolozWersjeUtrwalenia(ctx context.Context, zasob dane.PlikBiblioteki,
	bajty []byte, profil string) error {

	skrot := sha256.Sum256(bajty)
	suma := hex.EncodeToString(skrot[:])
	odwolanie, err := a.magazyn.Zapisz(bajty, suma)
	if err != nil {
		return bladBiblioteki(err)
	}
	rozmiar := int64(len(bajty))
	autor := "model"
	etykieta := "utrwalenie " + profil
	_, _, err = a.repozytorium.DolozWersje(ctx, zasob.ID, dane.WersjaPlikuBiblioteki{
		Kod: nowyIdentyfikator(przedrostekWersjiBiblioteki), PlikID: zasob.ID,
		Etykieta: &etykieta, Autor: &autor, RozmiarBajtow: &rozmiar,
		SumaKontrolna: &suma, TrescOdwolanie: &odwolanie,
	})
	if err != nil {
		return bladBiblioteki(err)
	}
	a.zglosNasluchom(shared.LibraryWebhookEventVersionAdded, zasob.Kod)
	return nil
}

// nazwaZUtrwaleniaBiblioteki składa nazwę wytworu utrwalenia z trzonu nazwy
// zasobu źródłowego i rozszerzenia właściwego postaci utrwalenia.
func nazwaZUtrwaleniaBiblioteki(nazwa, rozszerzenie string) string {
	trzon, _ := trzonIRozszerzenieBiblioteki(nazwa)
	return trzon + "-utrwalony." + rozszerzenie
}

// politykaKontraktuBiblioteki przenosi wiersz polityki retencji z bazy danych
// na kształt polityki zwracany kontraktem komunikacji.
func politykaKontraktuBiblioteki(wiersz dane.PolitykaRetencjiBiblioteki) shared.LibraryRetentionPolicy {
	return shared.LibraryRetentionPolicy{
		Id: wiersz.Kod, Scope: zasiegRetencjiKontraktu(wiersz.Zasieg), ScopeId: wiersz.ZasiegID,
		KeepDays: wiersz.DniPrzechowywania, Action: czynnoscRetencjiKontraktu(wiersz.Czynnosc),
		CreatedAt: chwilaBazy(wiersz.Utworzono),
	}
}

// rodzajUtrwaleniaBazy przekłada postać utrwalenia na wartość kolumny
// (odwzorowanie: `zadanie_utrwalenia_biblioteki.rodzaj`).
func rodzajUtrwaleniaBazy(rodzaj shared.LibraryPreservationKind) string {
	switch rodzaj {
	case shared.LibraryPreservationKindBagit:
		return "bagit"
	case shared.LibraryPreservationKindPremis:
		return "premis"
	default:
		return "pdfa"
	}
}

// czynnoscRetencjiBazy i czynnoscRetencjiKontraktu przekładają czynność polityki
// (odwzorowanie: `polityka_retencji_biblioteki.czynnosc`).
func czynnoscRetencjiBazy(czynnosc shared.LibraryRetentionAction) string {
	switch czynnosc {
	case shared.LibraryRetentionActionArchive:
		return "archiwizacja"
	case shared.LibraryRetentionActionDelete:
		return "usuniecie"
	default:
		return "przeglad"
	}
}

func czynnoscRetencjiKontraktu(czynnosc string) shared.LibraryRetentionAction {
	switch czynnosc {
	case "archiwizacja":
		return shared.LibraryRetentionActionArchive
	case "usuniecie":
		return shared.LibraryRetentionActionDelete
	default:
		return shared.LibraryRetentionActionReview
	}
}

// zasiegRetencjiBazy i zasiegRetencjiKontraktu przekładają poziom zasięgu
// kontraktu na kod katalogu `poziom_zasiegu` i z powrotem.
func zasiegRetencjiBazy(zasieg shared.ConfigScope) string {
	switch zasieg {
	case shared.ConfigScopeApplication:
		return "aplikacja"
	case shared.ConfigScopeEnvironment:
		return "srodowisko"
	case shared.ConfigScopeModule:
		return "modul"
	case shared.ConfigScopeModulePair:
		return "para_modulow"
	case shared.ConfigScopeProject:
		return "projekt"
	case shared.ConfigScopeSession:
		return "karta_sesji"
	case shared.ConfigScopeRole:
		return "rola"
	case shared.ConfigScopeWindow:
		return "okno"
	default:
		return "globalny"
	}
}

func zasiegRetencjiKontraktu(zasieg string) shared.ConfigScope {
	switch zasieg {
	case "aplikacja":
		return shared.ConfigScopeApplication
	case "srodowisko":
		return shared.ConfigScopeEnvironment
	case "modul":
		return shared.ConfigScopeModule
	case "para_modulow":
		return shared.ConfigScopeModulePair
	case "projekt":
		return shared.ConfigScopeProject
	case "karta_sesji":
		return shared.ConfigScopeSession
	case "rola":
		return shared.ConfigScopeRole
	case "okno":
		return shared.ConfigScopeWindow
	default:
		return shared.ConfigScopeGlobal
	}
}
