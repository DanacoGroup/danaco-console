// Skutek bezpieczeństwa dokumentu: czy czynność naprawdę zmienia dokument; żaden sprawdzian tu nie kończy się na odpowiedzi komendy, każdy schodzi do bajtów wyniku.
package core

import (
	"bytes"
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"math/big"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/pdfcpu/pdfcpu/pkg/api"

	"danacoconsole/shared"
)

// trescStronWyniku skleja treść wszystkich stron dokumentu leżącego pod
// odwołaniem zasobu — to jest miara redakcji, bo tekst wycięty znika właśnie
// stąd, a nie z odpowiedzi komendy.
func trescStronWyniku(t *testing.T, sciezka string) string {
	t.Helper()

	bajty, err := os.ReadFile(sciezka)
	if err != nil {
		t.Fatalf("nie można odczytać wyniku: %v", err)
	}
	drzewo, err := api.ReadValidateAndOptimize(bytes.NewReader(bajty), nastawyPdf())
	if err != nil {
		t.Fatalf("nie można otworzyć wyniku: %v", err)
	}
	var zebrane strings.Builder
	for numer := 1; numer <= drzewo.PageCount; numer++ {
		strona, _, _, err := drzewo.PageDict(numer, false)
		if err != nil || strona == nil {
			t.Fatalf("nie można sięgnąć po stronę %d wyniku", numer)
		}
		tresc, err := drzewo.PageContent(strona, numer)
		if err != nil {
			continue
		}
		zebrane.Write(tresc)
	}
	return zebrane.String()
}

// certyfikatProbny wytwarza parę klucz–certyfikat w postaci PEM.
//
// Certyfikat powstaje na czas sprawdzianu, a nie leży w drzewie: klucz wniesiony
// do repozytorium byłby kluczem, który wyciekł w chwili wniesienia.
func certyfikatProbny(t *testing.T, nazwa string) []byte {
	t.Helper()

	klucz, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("nie można wytworzyć klucza: %v", err)
	}
	wzorzec := x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject:      pkix.Name{CommonName: nazwa},
		NotBefore:    time.Now().Add(-time.Hour),
		NotAfter:     time.Now().Add(time.Hour),
	}
	surowy, err := x509.CreateCertificate(rand.Reader, &wzorzec, &wzorzec,
		&klucz.PublicKey, klucz)
	if err != nil {
		t.Fatalf("nie można wytworzyć certyfikatu: %v", err)
	}
	var zapis bytes.Buffer
	if err := pem.Encode(&zapis, &pem.Block{Type: "CERTIFICATE", Bytes: surowy}); err != nil {
		t.Fatalf("nie można zapisać certyfikatu: %v", err)
	}
	if err := pem.Encode(&zapis, &pem.Block{
		Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(klucz),
	}); err != nil {
		t.Fatalf("nie można zapisać klucza: %v", err)
	}
	return zapis.Bytes()
}

// wniesMaterial wnosi dowolną treść do magazynu zasobów okna, gotową do dalszego przetworzenia przez komendę.
func wniesMaterial(t *testing.T, zmontowany *Zmontowany, zycie context.Context,
	bajty []byte, nazwa, format string, rodzaj shared.DesignAssetKind) string {
	t.Helper()

	var wynik shared.DesignAssetUploadResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandDesignAssetUpload,
		shared.DesignAssetUploadRequest{
			WindowId:      "okno-sprawdzianu",
			Name:          wskaznik(nazwa),
			Kind:          rodzaj,
			Format:        wskaznik(format),
			ContentBase64: wskaznik(wBase64(bajty)),
		}, &wynik)
	return wynik.Asset.Id
}

// TestSzyfrowanieZamykaDokumentNaHaslo wykazuje skutek: wynik nie otwiera się
// bez hasła, a odpowiedź nie kłamie o stanie zabezpieczenia.
func TestSzyfrowanieZamykaDokumentNaHaslo(t *testing.T) {
	zmontowany, zycie, katalog := zmontujDoPomiaruSkutku(t)

	kod := wniesPdf(t, zmontowany, zycie, pdfProbny(t, 2), "jawny")

	var wynik shared.StudioSecurityEncryptResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandStudioSecurityEncrypt,
		shared.StudioSecurityEncryptRequest{
			AssetId:       kod,
			UserPassword:  wskaznik("otwarcie"),
			OwnerPassword: wskaznik("wlasciciel"),
			Permissions:   []string{"drukowanie"},
			WindowId:      wskaznik("okno-sprawdzianu"),
		}, &wynik)

	if !wynik.Encrypted {
		t.Fatal("odpowiedź melduje dokument niezaszyfrowany po nałożeniu haseł")
	}
	sciezka := sciezkaZasobuSprawdzianu(t, zmontowany, zycie, katalog, wynik.Asset.Id)
	bajty, err := os.ReadFile(sciezka)
	if err != nil {
		t.Fatalf("nie można odczytać wyniku: %v", err)
	}
	// Miara jest niezależna od odpowiedzi: dokument zaszyfrowany nie daje się otworzyć bez hasła.
	if _, err := api.PageCount(bytes.NewReader(bajty), nastawyPdf()); err == nil {
		t.Fatal("wynik otwiera się bez hasła — szyfrowania w bajtach nie ma")
	}
}

// TestCzyszczenieMetadanychUsuwaOpisZBajtow wykazuje, że opis znika z dokumentu,
// a nie tylko z odpowiedzi.
func TestCzyszczenieMetadanychUsuwaOpisZBajtow(t *testing.T) {
	zmontowany, zycie, katalog := zmontujDoPomiaruSkutku(t)

	// Opis wchodzi do materiału tą samą biblioteką, żeby sprawdzian mierzył usunięcie realnej treści.
	var zOpisem bytes.Buffer
	if err := api.AddProperties(bytes.NewReader(pdfProbny(t, 1)), &zOpisem,
		map[string]string{"Sprawa": "poufna", "Autor": "Operator"}, nastawyPdf()); err != nil {
		t.Fatalf("nie można założyć opisu materiałowi: %v", err)
	}
	kod := wniesPdf(t, zmontowany, zycie, zOpisem.Bytes(), "z-opisem")

	var wynik shared.StudioSecurityMetadataStripResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandStudioSecurityMetadataStrip,
		shared.StudioSecurityMetadataStripRequest{
			AssetId:  kod,
			WindowId: wskaznik("okno-sprawdzianu"),
		}, &wynik)

	sciezka := sciezkaZasobuSprawdzianu(t, zmontowany, zycie, katalog, wynik.Asset.Id)
	bajty, err := os.ReadFile(sciezka)
	if err != nil {
		t.Fatalf("nie można odczytać wyniku: %v", err)
	}
	wlasciwosci, err := api.Properties(bytes.NewReader(bajty), nastawyPdf())
	if err != nil {
		t.Fatalf("nie można odczytać właściwości wyniku: %v", err)
	}
	if _, jest := wlasciwosci["Sprawa"]; jest {
		t.Fatal("właściwość Sprawa została w wyniku mimo czyszczenia metadanych")
	}
}

// TestRedakcjaWycinaTekstZTresciStrony wykazuje, że redakcja USUWA, a nie
// zasłania: po czynności tekst nie leży już w treści strony.
func TestRedakcjaWycinaTekstZTresciStrony(t *testing.T) {
	zmontowany, zycie, katalog := zmontujDoPomiaruSkutku(t)

	material := pdfProbny(t, 1)
	kod := wniesPdf(t, zmontowany, zycie, material, "do-redakcji")

	var wynik shared.StudioSecurityRedactResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandStudioSecurityRedact,
		shared.StudioSecurityRedactRequest{
			AssetId:  kod,
			WindowId: wskaznik("okno-sprawdzianu"),
			// Obszar obejmuje całą stronę, bo materiał ma jeden napis, a miejsce jego wpisania jest nieznane.
			Regions: []shared.StudioRedactionRegion{{
				Page:   wskaznik(1),
				X:      wskaznik(0),
				Y:      wskaznik(0),
				Width:  wskaznik(2000),
				Height: wskaznik(2000),
			}},
		}, &wynik)

	if wynik.Redacted != 1 {
		t.Fatalf("odpowiedź melduje %d zamazanych obszarów przy jednym wskazanym",
			wynik.Redacted)
	}
	sciezka := sciezkaZasobuSprawdzianu(t, zmontowany, zycie, katalog, wynik.Asset.Id)
	tresc := trescStronWyniku(t, sciezka)
	if strings.Contains(tresc, "Strona") {
		t.Fatal("napis został w treści strony — to jest redakcja pozorna: " +
			"prostokąt zasłania, a tekst nadal da się skopiować")
	}
	if !strings.Contains(tresc, " re f") {
		t.Fatal("w treści strony nie ma prostokąta zamazania")
	}
}

// TestPodpisPrzezywaOdczytAZmianaTresciGoUniewaznia wykazuje obie strony podpisu: dokument nietknięty przechodzi weryfikację, a zmieniony ją oblewa.
func TestPodpisPrzezywaOdczytAZmianaTresciGoUniewaznia(t *testing.T) {
	zmontowany, zycie, _ := zmontujDoPomiaruSkutku(t)

	kod := wniesPdf(t, zmontowany, zycie, pdfProbny(t, 2), "do-podpisu")
	kodCertyfikatu := wniesMaterial(t, zmontowany, zycie,
		certyfikatProbny(t, "Operator Sprawdzianu"), "certyfikat", "pem",
		shared.DesignAssetKindDocument)

	var podpisany shared.StudioSecuritySignResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandStudioSecuritySign,
		shared.StudioSecuritySignRequest{
			AssetId:       kod,
			CertificateId: kodCertyfikatu,
			Reason:        wskaznik("zgodność z oryginałem"),
			WindowId:      wskaznik("okno-sprawdzianu"),
		}, &podpisany)

	var sprawdzony shared.StudioSecuritySignVerifyResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandStudioSecuritySignVerify,
		shared.StudioSecuritySignVerifyRequest{AssetId: podpisany.Asset.Id}, &sprawdzony)

	if len(sprawdzony.Signatures) != 1 || !sprawdzony.AllValid {
		t.Fatalf("dokument tuż po podpisaniu nie przechodzi weryfikacji: %+v", sprawdzony)
	}
	if sprawdzony.Signatures[0].Signer != "Operator Sprawdzianu" {
		t.Fatalf("podpis niesie podpisującego %q, a certyfikat mówił co innego",
			sprawdzony.Signatures[0].Signer)
	}

	// Zmiana treści: usunięcie strony podpisanego dokumentu. Podpis ma po niej przestać być poprawny.
	var pozmianie shared.StudioPdfPagesReorderResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandStudioPdfPagesReorder,
		shared.StudioPdfPagesReorderRequest{
			AssetId:  podpisany.Asset.Id,
			WindowId: wskaznik("okno-sprawdzianu"),
			Operations: []shared.StudioPdfPageOperation{{
				Kind:  shared.StudioPdfPageOperationKindWyodrebnienie,
				Pages: "1",
			}},
		}, &pozmianie)

	var poZmianie shared.StudioSecuritySignVerifyResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandStudioSecuritySignVerify,
		shared.StudioSecuritySignVerifyRequest{AssetId: pozmianie.Asset.Id}, &poZmianie)

	if poZmianie.AllValid {
		t.Fatal("podpis został poprawny po usunięciu strony — nie pilnuje treści")
	}
}

// TestRozpoznanieDanychWrazliwychOdrzucaLiczbyBezCyfryKontrolnej pilnuje, żeby rozpoznanie nie meldowało każdego ciągu cyfr jako numeru.
func TestRozpoznanieDanychWrazliwychOdrzucaLiczbyBezCyfryKontrolnej(t *testing.T) {
	zmontowany, zycie, _ := zmontujDoPomiaruSkutku(t)

	var otwarty shared.StudioDocumentOpenResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandStudioDocumentOpen,
		shared.StudioDocumentOpenRequest{WindowId: "okno-sprawdzianu"}, &otwarty)

	var zapisany shared.StudioDocumentSaveResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandStudioDocumentSave,
		shared.StudioDocumentSaveRequest{
			DocumentId: otwarty.Document.Id,
			Title:      wskaznik("pismo"),
			Content:    "Kontakt: biuro@danaco.pl, sprawa 12345678901.",
		}, &zapisany)

	var wynik shared.StudioSecuritySensitiveDetectResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandStudioSecuritySensitiveDetect,
		shared.StudioSecuritySensitiveDetectRequest{DocumentId: zapisany.Document.Id}, &wynik)

	poczta := 0
	for _, znalezisko := range wynik.Findings {
		if znalezisko.Category == "numer-ewidencyjny" {
			t.Fatalf("ciąg %q uznano za numer ewidencyjny mimo błędnej cyfry kontrolnej",
				znalezisko.Text)
		}
		if znalezisko.Category == "poczta" {
			poczta++
			if znalezisko.Text != "biuro@danaco.pl" {
				t.Fatalf("znalezisko poczty niesie %q", znalezisko.Text)
			}
		}
	}
	if poczta != 1 {
		t.Fatalf("adres poczty rozpoznano %d razy, a w treści jest jeden", poczta)
	}
}
