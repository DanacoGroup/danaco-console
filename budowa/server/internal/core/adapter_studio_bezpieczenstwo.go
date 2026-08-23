// Odpowiedzialność pliku: bezpieczeństwo dokumentu modułu Studio — szyfrowanie,
// czyszczenie metadanych, redakcja poufności, podpis i jego weryfikacja oraz
// rozpoznanie danych wrażliwych.
//
// ── Biblioteka wkompilowana, nigdy program zewnętrzny ───────────────────────
// Ta rodzina pracuje wyłącznie na `pdfcpu` i na bibliotece standardowej Go
// (`crypto/x509`, `crypto/rsa`, `crypto/sha256`). Nie ma tu ani jednego
// uruchomienia procesu: funkcja zależna od programu, którego instalka nie
// niesie, jest u Operatora odmową, a nie funkcją — a rodzina bezpieczeństwa jest
// tym miejscem, w którym odmowa boli najbardziej, bo Operator dowiaduje się
// o niej dopiero po wysłaniu dokumentu, którego nie oczyścił.
//
// ── Redakcja usuwa, a nie zasłania ─────────────────────────────────────────
// Zamalowanie prostokąta zostawia tekst pod spodem: da się go zaznaczyć,
// skopiować i odczytać wyszukiwarką. To jest redakcja pozorna i tutaj jej nie
// ma. `Zredaguj` wycina z treści strony bloki tekstu leżące w obszarze, dopiero
// potem kładzie prostokąt — i wynik nie zawiera już wyciętych znaków.
package core

import (
	"bytes"
	"context"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"encoding/pem"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/pdfcpu/pdfcpu/pkg/api"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/model"

	"danacoconsole/server/internal/dane"
	"danacoconsole/server/internal/protocol"
	"danacoconsole/shared"
)

// wlasciwoscPodpisu jest nazwą właściwości dokumentu, pod którą leży podpis.
const wlasciwoscPodpisu = "DanacoPodpis"

type adapterBezpieczenstwaStudia struct {
	warsztat *adapterPdfStudia
	studio   dane.RepozytoriumStudia
}

func nowyAdapterBezpieczenstwaStudia(warsztat *adapterPdfStudia) *adapterBezpieczenstwaStudia {
	return &adapterBezpieczenstwaStudia{warsztat: warsztat}
}

// ZeStudiem podaje repozytorium dokumentów: rozpoznanie danych wrażliwych czyta
// treść dokumentu Studia, a nie zasób magazynu.
func (a *adapterBezpieczenstwaStudia) ZeStudiem(r dane.RepozytoriumStudia) *adapterBezpieczenstwaStudia {
	a.studio = r
	return a
}

// odmowaBezpieczenstwa buduje odmowę rodziny.
func odmowaBezpieczenstwa(kod shared.ErrorCode, powod string) error {
	return protocol.JakoError(protocol.NowyBlad(kod, "bezpieczeństwo dokumentu: "+powod))
}

// ── Szyfrowanie ─────────────────────────────────────────────────────────────

// uprawnieniaOdbiorcy odwzorowuje nazwy uprawnień na bity dokumentu.
//
// Nazwy są nazwane po czynnościach, a nie po numerach bitów specyfikacji:
// Operator ustawia „drukowanie", nie „bit 3".
var uprawnieniaOdbiorcy = map[string]model.PermissionFlags{
	"drukowanie":   model.PermissionPrintRev2 + model.PermissionPrintRev3,
	"zmiana":       model.PermissionModify,
	"kopiowanie":   model.PermissionExtract + model.PermissionExtractRev3,
	"komentowanie": model.PermissionModAnnFillForm,
	"wypelnianie":  model.PermissionFillRev3,
	"skladanie":    model.PermissionAssembleRev3,
}

// nazwyUprawnien oddaje przyjmowane nazwy w kolejności ustalonej, żeby odmowa
// wymieniała je zawsze tak samo.
func nazwyUprawnien() []string {
	nazwy := make([]string, 0, len(uprawnieniaOdbiorcy))
	for nazwa := range uprawnieniaOdbiorcy {
		nazwy = append(nazwy, nazwa)
	}
	sort.Strings(nazwy)
	return nazwy
}

// ZaszyfrujPdf nakłada albo zdejmuje szyfrowanie dokumentu.
//
// Puste hasło otwarcia zdejmuje zabezpieczenie — tak mówi kontrakt. Zdjęcie
// wymaga hasła właściciela, bo bez niego dokument zaszyfrowany nie otwiera się
// nikomu, także rdzeniowi.
func (a *adapterBezpieczenstwaStudia) ZaszyfrujPdf(ctx context.Context,
	z shared.StudioSecurityEncryptRequest) (shared.StudioSecurityEncryptResponse, error) {

	bajty, err := a.warsztat.trescMaterialu(ctx, z.AssetId)
	if err != nil {
		return shared.StudioSecurityEncryptResponse{}, err
	}
	hasloOtwarcia := wartoscTekstu(z.UserPassword)
	hasloWlasciciela := wartoscTekstu(z.OwnerPassword)

	nastawy := nastawyPdf()
	nastawy.UserPW = hasloOtwarcia
	nastawy.OwnerPW = hasloWlasciciela
	nastawy.EncryptUsingAES = true

	var wynik bytes.Buffer
	zaszyfrowany := hasloOtwarcia != ""

	if zaszyfrowany {
		uprawnienia := model.PermissionsNone
		for _, nazwa := range z.Permissions {
			bity, znane := uprawnieniaOdbiorcy[strings.ToLower(strings.TrimSpace(nazwa))]
			if !znane {
				return shared.StudioSecurityEncryptResponse{}, odmowaBezpieczenstwa(
					shared.ErrorCodeValidationFailed,
					"uprawnienie „"+nazwa+"” nie jest znane; przyjmowane są: "+
						strings.Join(nazwyUprawnien(), ", "))
			}
			uprawnienia += bity
		}
		nastawy.Permissions = uprawnienia
		if err := api.Encrypt(bytes.NewReader(bajty), &wynik, nastawy); err != nil {
			return shared.StudioSecurityEncryptResponse{}, odmowaBezpieczenstwa(
				shared.ErrorCodeInternalError, "nie można zaszyfrować dokumentu: "+err.Error())
		}
	} else {
		if hasloWlasciciela == "" {
			return shared.StudioSecurityEncryptResponse{}, odmowaBezpieczenstwa(
				shared.ErrorCodeValidationFailed,
				"zdjęcie szyfrowania wymaga hasła właściciela — bez niego dokumentu "+
					"nie otwiera nikt, także rdzeń")
		}
		if err := api.Decrypt(bytes.NewReader(bajty), &wynik, nastawy); err != nil {
			return shared.StudioSecurityEncryptResponse{}, odmowaBezpieczenstwa(
				shared.ErrorCodeValidationFailed,
				"nie można zdjąć szyfrowania — hasło właściciela nie pasuje: "+err.Error())
		}
	}

	nazwa := "dokument zaszyfrowany"
	if !zaszyfrowany {
		nazwa = "dokument bez szyfrowania"
	}
	zasob, err := a.warsztat.odlozPdf(ctx, wynik.Bytes(), nazwa, wartoscTekstu(z.WindowId))
	if err != nil {
		return shared.StudioSecurityEncryptResponse{}, err
	}
	return shared.StudioSecurityEncryptResponse{Asset: zasob, Encrypted: zaszyfrowany}, nil
}

// ── Metadane ────────────────────────────────────────────────────────────────

// WyczyscMetadane zdejmuje z dokumentu opis, historię i dane ukryte.
//
// Czyszczenie idzie po samym drzewie dokumentu, a nie po jego przepisaniu:
// usuwane są pola opisu, właściwości własne i strumień metadanych XMP, w którym
// leży ta sama treść w postaci równoległej. Zdjęcie samego opisu z pozostawieniem
// XMP byłoby czyszczeniem pozornym — narzędzia czytają XMP przed opisem.
func (a *adapterBezpieczenstwaStudia) WyczyscMetadane(ctx context.Context,
	z shared.StudioSecurityMetadataStripRequest) (shared.StudioSecurityMetadataStripResponse, error) {

	bajty, err := a.warsztat.trescMaterialu(ctx, z.AssetId)
	if err != nil {
		return shared.StudioSecurityMetadataStripResponse{}, err
	}
	zostawiane := map[string]bool{}
	for _, pole := range z.KeepFields {
		zostawiane[strings.ToLower(strings.TrimSpace(pole))] = true
	}

	drzewo, err := api.ReadValidateAndOptimize(bytes.NewReader(bajty), nastawyPdf())
	if err != nil {
		return shared.StudioSecurityMetadataStripResponse{}, odmowaBezpieczenstwa(
			shared.ErrorCodeValidationFailed, "nie można otworzyć dokumentu: "+err.Error())
	}
	usuniete := make([]string, 0)

	if drzewo.Info != nil {
		opis, err := drzewo.DereferenceDict(*drzewo.Info)
		if err == nil && opis != nil {
			for nazwa := range opis {
				if zostawiane[strings.ToLower(nazwa)] {
					continue
				}
				delete(opis, nazwa)
				usuniete = append(usuniete, nazwa)
			}
		}
	}

	// Katalog niesie strumień XMP oraz ślad historii narzędzi, którymi dokument
	// szedł. Jedno i drugie schodzi razem z opisem, bo razem opisuje pochodzenie.
	if katalog, err := drzewo.Catalog(); err == nil && katalog != nil {
		for _, nazwa := range []string{"Metadata", "PieceInfo"} {
			if zostawiane[strings.ToLower(nazwa)] {
				continue
			}
			if _, jest := katalog.Find(nazwa); jest {
				delete(katalog, nazwa)
				usuniete = append(usuniete, nazwa)
			}
		}
	}

	sort.Strings(usuniete)
	var wynik bytes.Buffer
	if err := api.WriteContext(drzewo, &wynik); err != nil {
		return shared.StudioSecurityMetadataStripResponse{}, odmowaBezpieczenstwa(
			shared.ErrorCodeInternalError, "nie można zapisać dokumentu: "+err.Error())
	}
	zasob, err := a.warsztat.odlozPdf(ctx, wynik.Bytes(), "dokument bez metadanych",
		wartoscTekstu(z.WindowId))
	if err != nil {
		return shared.StudioSecurityMetadataStripResponse{}, err
	}
	return shared.StudioSecurityMetadataStripResponse{Asset: zasob, Removed: usuniete}, nil
}

// ── Redakcja ────────────────────────────────────────────────────────────────

// obszarRedakcji jest prostokątem wskazanym przez Operatora, sprowadzonym do
// współrzędnych dokumentu.
//
// Początek układu leży w LEWYM DOLNYM rogu strony, tak jak w samym dokumencie.
// Wskazanie liczone od góry dawałoby zamazanie przesunięte o wysokość strony,
// a Operator zobaczyłby, że zamazało nie to miejsce, dopiero po otwarciu wyniku.
type obszarRedakcji struct {
	x, y, szerokosc, wysokosc float64
}

func (o obszarRedakcji) zawiera(x, y float64) bool {
	return x >= o.x && x <= o.x+o.szerokosc && y >= o.y && y <= o.y+o.wysokosc
}

func (o obszarRedakcji) przecina(x1, y1, x2, y2 float64) bool {
	return x1 <= o.x+o.szerokosc && x2 >= o.x && y1 <= o.y+o.wysokosc && y2 >= o.y
}

// Zredaguj trwale usuwa treść z obszarów wskazanych przez Operatora.
//
// Materiał zostaje nietknięty, bo czynność jest nieodwracalna dla wyniku: gdyby
// szła w miejscu, pomyłka we współrzędnych kasowałaby treść bezpowrotnie.
//
// Ziarno wycięcia to blok tekstu (BT…ET), a nie pojedynczy znak. Wycięcie idzie
// więc czasem szerzej, niż wskazał Operator, i nigdy węziej — a przy redakcji
// tylko ten kierunek błędu jest dopuszczalny.
func (a *adapterBezpieczenstwaStudia) Zredaguj(ctx context.Context,
	z shared.StudioSecurityRedactRequest) (shared.StudioSecurityRedactResponse, error) {

	bajty, err := a.warsztat.trescMaterialu(ctx, z.AssetId)
	if err != nil {
		return shared.StudioSecurityRedactResponse{}, err
	}
	if len(z.Regions) == 0 {
		return shared.StudioSecurityRedactResponse{}, odmowaBezpieczenstwa(
			shared.ErrorCodeValidationFailed,
			"żądanie nie wskazuje ani jednego obszaru — nie ma czego zredagować")
	}

	obszary := map[int][]obszarRedakcji{}
	for _, obszar := range z.Regions {
		if obszar.Page == nil {
			return shared.StudioSecurityRedactResponse{}, odmowaBezpieczenstwa(
				shared.ErrorCodeValidationFailed,
				"obszar bez wskazania strony opisuje fragment tekstu, a materiałem jest "+
					"dokument PDF; wskaż stronę oraz prostokąt")
		}
		if obszar.X == nil || obszar.Y == nil || obszar.Width == nil || obszar.Height == nil {
			return shared.StudioSecurityRedactResponse{}, odmowaBezpieczenstwa(
				shared.ErrorCodeValidationFailed,
				"obszar na stronie "+strconv.Itoa(*obszar.Page)+" nie niesie kompletu "+
					"współrzędnych: potrzebne są x, y, szerokość i wysokość")
		}
		if *obszar.Width <= 0 || *obszar.Height <= 0 {
			return shared.StudioSecurityRedactResponse{}, odmowaBezpieczenstwa(
				shared.ErrorCodeValidationFailed,
				"obszar na stronie "+strconv.Itoa(*obszar.Page)+" ma zerowy albo ujemny "+
					"wymiar — nie zamaże niczego")
		}
		obszary[*obszar.Page] = append(obszary[*obszar.Page], obszarRedakcji{
			x: float64(*obszar.X), y: float64(*obszar.Y),
			szerokosc: float64(*obszar.Width), wysokosc: float64(*obszar.Height),
		})
	}

	drzewo, err := api.ReadValidateAndOptimize(bytes.NewReader(bajty), nastawyPdf())
	if err != nil {
		return shared.StudioSecurityRedactResponse{}, odmowaBezpieczenstwa(
			shared.ErrorCodeValidationFailed, "nie można otworzyć dokumentu: "+err.Error())
	}
	strony := make([]int, 0, len(obszary))
	for numer := range obszary {
		strony = append(strony, numer)
	}
	sort.Ints(strony)

	for _, numer := range strony {
		if numer < 1 || numer > drzewo.PageCount {
			return shared.StudioSecurityRedactResponse{}, odmowaBezpieczenstwa(
				shared.ErrorCodeValidationFailed,
				"dokument ma "+strconv.Itoa(drzewo.PageCount)+" stron, a obszar wskazuje "+
					"stronę "+strconv.Itoa(numer))
		}
		strona, _, _, err := drzewo.PageDict(numer, false)
		if err != nil || strona == nil {
			return shared.StudioSecurityRedactResponse{}, odmowaBezpieczenstwa(
				shared.ErrorCodeInternalError,
				"nie można sięgnąć po stronę "+strconv.Itoa(numer))
		}
		tresc, err := drzewo.PageContent(strona, numer)
		if err != nil {
			// Strona bez treści nie ma czego stracić; prostokąt i tak na nią idzie,
			// żeby wynik wyglądał tak samo niezależnie od tego, czym stronę wypełniono.
			tresc = nil
		}

		oczyszczona := wytnijBlokiTekstu(tresc, obszary[numer])
		oczyszczona = append(oczyszczona, prostokatyZamazania(obszary[numer])...)

		strumien, err := drzewo.NewStreamDictForBuf(oczyszczona)
		if err != nil {
			return shared.StudioSecurityRedactResponse{}, odmowaBezpieczenstwa(
				shared.ErrorCodeInternalError,
				"nie można złożyć treści strony "+strconv.Itoa(numer)+": "+err.Error())
		}
		if err := strumien.Encode(); err != nil {
			return shared.StudioSecurityRedactResponse{}, odmowaBezpieczenstwa(
				shared.ErrorCodeInternalError,
				"nie można zakodować treści strony "+strconv.Itoa(numer)+": "+err.Error())
		}
		odwolanie, err := drzewo.IndRefForNewObject(*strumien)
		if err != nil {
			return shared.StudioSecurityRedactResponse{}, odmowaBezpieczenstwa(
				shared.ErrorCodeInternalError,
				"nie można odłożyć treści strony "+strconv.Itoa(numer)+": "+err.Error())
		}
		strona.Update("Contents", *odwolanie)
	}

	var wynik bytes.Buffer
	if err := api.WriteContext(drzewo, &wynik); err != nil {
		return shared.StudioSecurityRedactResponse{}, odmowaBezpieczenstwa(
			shared.ErrorCodeInternalError, "nie można zapisać dokumentu: "+err.Error())
	}
	zasob, err := a.warsztat.odlozPdf(ctx, wynik.Bytes(), "dokument po redakcji",
		wartoscTekstu(z.WindowId))
	if err != nil {
		return shared.StudioSecurityRedactResponse{}, err
	}
	return shared.StudioSecurityRedactResponse{Asset: zasob, Redacted: len(z.Regions)}, nil
}

// prostokatyZamazania składa polecenia rysunkowe kładące czarne pola.
//
// Prostokąt idzie PO wycięciu, więc zamazuje to, co pod obszarem zostało —
// obrazy i grafikę. Tekst już wtedy z treści nie istnieje.
func prostokatyZamazania(obszary []obszarRedakcji) []byte {
	var polecenia bytes.Buffer
	for _, obszar := range obszary {
		polecenia.WriteString("\nq 0 g ")
		polecenia.WriteString(strconv.FormatFloat(obszar.x, 'f', 2, 64) + " ")
		polecenia.WriteString(strconv.FormatFloat(obszar.y, 'f', 2, 64) + " ")
		polecenia.WriteString(strconv.FormatFloat(obszar.szerokosc, 'f', 2, 64) + " ")
		polecenia.WriteString(strconv.FormatFloat(obszar.wysokosc, 'f', 2, 64) + " ")
		polecenia.WriteString("re f Q\n")
	}
	return polecenia.Bytes()
}

// wytnijBlokiTekstu usuwa z treści strony te bloki tekstu, których położenie
// wypada w którymkolwiek obszarze.
//
// Czytanie idzie po składni treści, a nie po prostym szukaniu napisów: napis
// „BT" bywa treścią rysowanego tekstu, a wycięcie w niewłaściwym miejscu
// rozsypałoby stronę.
func wytnijBlokiTekstu(tresc []byte, obszary []obszarRedakcji) []byte {
	if len(tresc) == 0 {
		return nil
	}
	stan := nowyCzytnikTresci(tresc)
	wyciete := stan.blokiWObszarach(obszary)
	if len(wyciete) == 0 {
		return tresc
	}

	var wynik bytes.Buffer
	poprzedni := 0
	for _, blok := range wyciete {
		wynik.Write(tresc[poprzedni:blok.poczatek])
		poprzedni = blok.koniec
	}
	wynik.Write(tresc[poprzedni:])
	return wynik.Bytes()
}

// blokTekstu jest zakresem bajtów jednego bloku BT…ET wraz z jego położeniem.
type blokTekstu struct {
	poczatek, koniec int
	x, y             float64
}

// czytnikTresci przechodzi treść strony operator po operatorze.
type czytnikTresci struct {
	tresc []byte
}

func nowyCzytnikTresci(tresc []byte) *czytnikTresci {
	return &czytnikTresci{tresc: tresc}
}

// przeksztalcenie jest stanem układu współrzędnych: przesunięcie i skala.
//
// Obrót i pochylenie nie są tu brane pod uwagę: strona obrócona daje położenie
// przybliżone, a przybliżenie przy redakcji wychodzi na wycięcie szersze, nie
// węższe.
type przeksztalcenie struct {
	przesuniecieX, przesuniecieY float64
	skalaX, skalaY               float64
}

// blokiWObszarach oddaje bloki tekstu leżące w którymkolwiek z obszarów,
// w kolejności występowania.
func (c *czytnikTresci) blokiWObszarach(obszary []obszarRedakcji) []blokTekstu {
	bloki := make([]blokTekstu, 0)
	stos := []przeksztalcenie{{skalaX: 1, skalaY: 1}}
	biezace := func() przeksztalcenie { return stos[len(stos)-1] }

	var liczby []float64
	wTekscie := false
	poczatekBloku := 0
	var tekstX, tekstY float64

	i := 0
	for i < len(c.tresc) {
		poczatekTokenu := i
		token, dalej := c.token(i)
		if dalej == i {
			break
		}
		i = dalej
		if token == "" {
			continue
		}

		if liczba, err := strconv.ParseFloat(token, 64); err == nil {
			liczby = append(liczby, liczba)
			continue
		}

		switch token {
		case "q":
			stos = append(stos, biezace())
		case "Q":
			if len(stos) > 1 {
				stos = stos[:len(stos)-1]
			}
		case "cm":
			if len(liczby) >= 6 {
				aktualne := biezace()
				aktualne.przesuniecieX += aktualne.skalaX * liczby[len(liczby)-2]
				aktualne.przesuniecieY += aktualne.skalaY * liczby[len(liczby)-1]
				aktualne.skalaX *= liczby[len(liczby)-6]
				aktualne.skalaY *= liczby[len(liczby)-3]
				stos[len(stos)-1] = aktualne
			}
		case "BT":
			wTekscie = true
			poczatekBloku = poczatekTokenu
			tekstX, tekstY = 0, 0
		case "Tm":
			if len(liczby) >= 6 {
				tekstX = liczby[len(liczby)-2]
				tekstY = liczby[len(liczby)-1]
			}
		case "Td", "TD":
			if len(liczby) >= 2 {
				tekstX += liczby[len(liczby)-2]
				tekstY += liczby[len(liczby)-1]
			}
		case "ET":
			if wTekscie {
				aktualne := biezace()
				x := aktualne.przesuniecieX + aktualne.skalaX*tekstX
				y := aktualne.przesuniecieY + aktualne.skalaY*tekstY
				for _, obszar := range obszary {
					if obszar.zawiera(x, y) {
						bloki = append(bloki, blokTekstu{
							poczatek: poczatekBloku, koniec: i, x: x, y: y,
						})
						break
					}
				}
			}
			wTekscie = false
		}
		liczby = liczby[:0]
	}
	return bloki
}

// token oddaje kolejny token treści oraz położenie za nim.
//
// Napisy, napisy szesnastkowe i komentarze przechodzą w całości: ich zawartość
// nie jest składnią i nie wolno jej czytać jako operatorów.
func (c *czytnikTresci) token(i int) (string, int) {
	for i < len(c.tresc) && bialyZnak(c.tresc[i]) {
		i++
	}
	if i >= len(c.tresc) {
		return "", i
	}
	switch c.tresc[i] {
	case '%':
		for i < len(c.tresc) && c.tresc[i] != '\n' && c.tresc[i] != '\r' {
			i++
		}
		return "", i
	case '(':
		zaglebienie := 0
		for i < len(c.tresc) {
			switch c.tresc[i] {
			case '\\':
				i++
			case '(':
				zaglebienie++
			case ')':
				zaglebienie--
				if zaglebienie == 0 {
					return "", i + 1
				}
			}
			i++
		}
		return "", i
	case '<':
		if i+1 < len(c.tresc) && c.tresc[i+1] == '<' {
			return "<<", i + 2
		}
		for i < len(c.tresc) && c.tresc[i] != '>' {
			i++
		}
		return "", i + 1
	case '>':
		if i+1 < len(c.tresc) && c.tresc[i+1] == '>' {
			return ">>", i + 2
		}
		return ">", i + 1
	case '[', ']', '{', '}':
		return string(c.tresc[i]), i + 1
	case '/':
		poczatek := i
		i++
		for i < len(c.tresc) && !bialyZnak(c.tresc[i]) && !ogranicznik(c.tresc[i]) {
			i++
		}
		return string(c.tresc[poczatek:i]), i
	}

	poczatek := i
	for i < len(c.tresc) && !bialyZnak(c.tresc[i]) && !ogranicznik(c.tresc[i]) {
		i++
	}
	return string(c.tresc[poczatek:i]), i
}

func bialyZnak(znak byte) bool {
	return znak == ' ' || znak == '\n' || znak == '\r' || znak == '\t' ||
		znak == '\f' || znak == 0
}

func ogranicznik(znak byte) bool {
	return znak == '(' || znak == ')' || znak == '<' || znak == '>' ||
		znak == '[' || znak == ']' || znak == '{' || znak == '}' ||
		znak == '/' || znak == '%'
}

// ── Podpis ──────────────────────────────────────────────────────────────────

// podpisDanaco jest treścią podpisu odkładaną we właściwościach dokumentu.
type podpisDanaco struct {
	Podpisujacy string `json:"podpisujacy"`
	Czas        int64  `json:"czas"`
	Powod       string `json:"powod,omitempty"`
	Miejsce     string `json:"miejsce,omitempty"`
	Skrot       string `json:"skrot"`
	Podpis      string `json:"podpis"`
	Certyfikat  string `json:"certyfikat"`
}

// skrotTresciDokumentu liczy skrót po treści stron, nie po bajtach pliku.
//
// Bajty pliku zmieniają się przy każdym zapisie — także przy dołożeniu samego
// podpisu — więc podpis liczony po nich nie dałby się zweryfikować ani razu.
// Skrót po treści stron opisuje to, co Operator widzi, i przeżywa dołożenie
// właściwości.
func skrotTresciDokumentu(drzewo *model.Context) (string, error) {
	suma := sha256.New()
	suma.Write([]byte(strconv.Itoa(drzewo.PageCount) + "\n"))
	for numer := 1; numer <= drzewo.PageCount; numer++ {
		strona, _, _, err := drzewo.PageDict(numer, false)
		if err != nil || strona == nil {
			return "", odmowaBezpieczenstwa(shared.ErrorCodeInternalError,
				"nie można sięgnąć po stronę "+strconv.Itoa(numer))
		}
		tresc, err := drzewo.PageContent(strona, numer)
		if err != nil {
			tresc = nil
		}
		suma.Write(tresc)
	}
	return hex.EncodeToString(suma.Sum(nil)), nil
}

// kluczICertyfikat rozkłada zasób certyfikatu na klucz prywatny i certyfikat.
//
// Zasób niesie obie części w postaci PEM: bez klucza nie ma czym podpisać,
// a bez certyfikatu nie ma czym sprawdzić.
func kluczICertyfikat(bajty []byte) (*rsa.PrivateKey, *x509.Certificate, error) {
	var klucz *rsa.PrivateKey
	var certyfikat *x509.Certificate

	reszta := bajty
	for {
		blok, dalej := pem.Decode(reszta)
		if blok == nil {
			break
		}
		reszta = dalej
		switch blok.Type {
		case "RSA PRIVATE KEY":
			odczytany, err := x509.ParsePKCS1PrivateKey(blok.Bytes)
			if err != nil {
				return nil, nil, odmowaBezpieczenstwa(shared.ErrorCodeValidationFailed,
					"klucz prywatny certyfikatu jest nieczytelny: "+err.Error())
			}
			klucz = odczytany
		case "PRIVATE KEY":
			odczytany, err := x509.ParsePKCS8PrivateKey(blok.Bytes)
			if err != nil {
				return nil, nil, odmowaBezpieczenstwa(shared.ErrorCodeValidationFailed,
					"klucz prywatny certyfikatu jest nieczytelny: "+err.Error())
			}
			rsaKlucz, jest := odczytany.(*rsa.PrivateKey)
			if !jest {
				return nil, nil, odmowaBezpieczenstwa(shared.ErrorCodeValidationFailed,
					"podpis dokumentu idzie kluczem RSA; podany klucz jest innego rodzaju")
			}
			klucz = rsaKlucz
		case "CERTIFICATE":
			odczytany, err := x509.ParseCertificate(blok.Bytes)
			if err != nil {
				return nil, nil, odmowaBezpieczenstwa(shared.ErrorCodeValidationFailed,
					"certyfikat jest nieczytelny: "+err.Error())
			}
			certyfikat = odczytany
		}
	}
	if klucz == nil || certyfikat == nil {
		return nil, nil, odmowaBezpieczenstwa(shared.ErrorCodeValidationFailed,
			"zasób certyfikatu niesie obie części w postaci PEM: certyfikat oraz "+
				"klucz prywatny; brakuje jednej z nich")
	}
	return klucz, certyfikat, nil
}

// PodpiszPdf podpisuje dokument certyfikatem wskazanym przez Operatora.
//
// Podpis leży we właściwości dokumentu i jest podpisem tego produktu: cudze
// czytniki go nie pokażą. Nazwane wprost, bo podpis, o którym Operator myśli,
// że jest podpisem kwalifikowanym, byłby gorszy niż brak podpisu.
func (a *adapterBezpieczenstwaStudia) PodpiszPdf(ctx context.Context,
	z shared.StudioSecuritySignRequest) (shared.StudioSecuritySignResponse, error) {

	bajty, err := a.warsztat.trescMaterialu(ctx, z.AssetId)
	if err != nil {
		return shared.StudioSecuritySignResponse{}, err
	}
	materialCertyfikatu, err := a.warsztat.bajtyZasobu(ctx, z.CertificateId)
	if err != nil {
		return shared.StudioSecuritySignResponse{}, err
	}
	klucz, certyfikat, err := kluczICertyfikat(materialCertyfikatu)
	if err != nil {
		return shared.StudioSecuritySignResponse{}, err
	}

	drzewo, err := api.ReadValidateAndOptimize(bytes.NewReader(bajty), nastawyPdf())
	if err != nil {
		return shared.StudioSecuritySignResponse{}, odmowaBezpieczenstwa(
			shared.ErrorCodeValidationFailed, "nie można otworzyć dokumentu: "+err.Error())
	}
	skrot, err := skrotTresciDokumentu(drzewo)
	if err != nil {
		return shared.StudioSecuritySignResponse{}, err
	}
	surowy, _ := hex.DecodeString(skrot)
	podpis, err := rsa.SignPKCS1v15(rand.Reader, klucz, crypto.SHA256, surowy)
	if err != nil {
		return shared.StudioSecuritySignResponse{}, odmowaBezpieczenstwa(
			shared.ErrorCodeInternalError, "nie można złożyć podpisu: "+err.Error())
	}

	czas := time.Now().UnixMilli()
	tresc := podpisDanaco{
		Podpisujacy: certyfikat.Subject.CommonName,
		Czas:        czas,
		Powod:       wartoscTekstu(z.Reason),
		Miejsce:     wartoscTekstu(z.Location),
		Skrot:       skrot,
		Podpis:      base64.StdEncoding.EncodeToString(podpis),
		Certyfikat:  base64.StdEncoding.EncodeToString(certyfikat.Raw),
	}
	zapis, err := json.Marshal(tresc)
	if err != nil {
		return shared.StudioSecuritySignResponse{}, odmowaBezpieczenstwa(
			shared.ErrorCodeInternalError, "nie można zapisać podpisu: "+err.Error())
	}

	var wynik bytes.Buffer
	if err := api.WriteContext(drzewo, &wynik); err != nil {
		return shared.StudioSecuritySignResponse{}, odmowaBezpieczenstwa(
			shared.ErrorCodeInternalError, "nie można zapisać dokumentu: "+err.Error())
	}
	istniejace, _ := api.Properties(bytes.NewReader(wynik.Bytes()), nastawyPdf())
	wlasciwosci := map[string]string{}
	for nazwa, wartosc := range istniejace {
		wlasciwosci[nazwa] = wartosc
	}
	wlasciwosci[wlasciwoscPodpisu] = dopiszPodpis(wlasciwosci[wlasciwoscPodpisu], zapis)

	var podpisany bytes.Buffer
	if err := api.AddProperties(bytes.NewReader(wynik.Bytes()), &podpisany, wlasciwosci,
		nastawyPdf()); err != nil {
		return shared.StudioSecuritySignResponse{}, odmowaBezpieczenstwa(
			shared.ErrorCodeInternalError, "nie można osadzić podpisu: "+err.Error())
	}

	zasob, err := a.warsztat.odlozPdf(ctx, podpisany.Bytes(), "dokument podpisany",
		wartoscTekstu(z.WindowId))
	if err != nil {
		return shared.StudioSecuritySignResponse{}, err
	}
	return shared.StudioSecuritySignResponse{Asset: zasob, SignedAt: czas}, nil
}

// dopiszPodpis dokłada podpis do listy podpisów już leżących w dokumencie.
//
// Dokument bywa podpisywany wielokrotnie i drugi podpis nie zastępuje
// pierwszego: kontrakt weryfikacji mówi o podpisach w liczbie mnogiej.
func dopiszPodpis(zastane string, nowy []byte) string {
	wpis := base64.StdEncoding.EncodeToString(nowy)
	if strings.TrimSpace(zastane) == "" {
		return wpis
	}
	return zastane + " " + wpis
}

// SprawdzPodpisy weryfikuje podpisy dokumentu.
//
// Weryfikacja przelicza skrót treści z dokumentu, który ma w ręku, i porównuje
// go ze skrótem podpisanym. Zmiana choćby jednej strony rozjeżdża oba i podpis
// przestaje być poprawny — właśnie po to on jest.
func (a *adapterBezpieczenstwaStudia) SprawdzPodpisy(ctx context.Context,
	z shared.StudioSecuritySignVerifyRequest) (shared.StudioSecuritySignVerifyResponse, error) {

	bajty, err := a.warsztat.trescMaterialu(ctx, z.AssetId)
	if err != nil {
		return shared.StudioSecuritySignVerifyResponse{}, err
	}
	wlasciwosci, err := api.Properties(bytes.NewReader(bajty), nastawyPdf())
	if err != nil {
		return shared.StudioSecuritySignVerifyResponse{}, odmowaBezpieczenstwa(
			shared.ErrorCodeValidationFailed, "nie można odczytać dokumentu: "+err.Error())
	}
	surowe := strings.Fields(wlasciwosci[wlasciwoscPodpisu])
	if len(surowe) == 0 {
		return shared.StudioSecuritySignVerifyResponse{
			Signatures: []shared.StudioSignature{}, AllValid: false,
		}, nil
	}

	drzewo, err := api.ReadValidateAndOptimize(bytes.NewReader(bajty), nastawyPdf())
	if err != nil {
		return shared.StudioSecuritySignVerifyResponse{}, odmowaBezpieczenstwa(
			shared.ErrorCodeValidationFailed, "nie można otworzyć dokumentu: "+err.Error())
	}
	skrotBiezacy, err := skrotTresciDokumentu(drzewo)
	if err != nil {
		return shared.StudioSecuritySignVerifyResponse{}, err
	}

	podpisy := make([]shared.StudioSignature, 0, len(surowe))
	wszystkiePoprawne := true
	for _, wpis := range surowe {
		podpis, powod := odczytajPodpis(wpis, skrotBiezacy)
		if powod != "" {
			podpis.Valid = false
			przyczyna := powod
			podpis.FailureReason = &przyczyna
		}
		if !podpis.Valid {
			wszystkiePoprawne = false
		}
		podpisy = append(podpisy, podpis)
	}
	return shared.StudioSecuritySignVerifyResponse{
		Signatures: podpisy, AllValid: wszystkiePoprawne,
	}, nil
}

// odczytajPodpis sprawdza jeden wpis podpisu i oddaje jego opis wraz z powodem
// niepowodzenia, gdy weryfikacja nie wyszła.
func odczytajPodpis(wpis, skrotBiezacy string) (shared.StudioSignature, string) {
	zapis, err := base64.StdEncoding.DecodeString(wpis)
	if err != nil {
		return shared.StudioSignature{}, "wpis podpisu jest nieczytelny"
	}
	var tresc podpisDanaco
	if err := json.Unmarshal(zapis, &tresc); err != nil {
		return shared.StudioSignature{}, "treść podpisu jest nieczytelna"
	}
	opis := shared.StudioSignature{
		Signer: tresc.Podpisujacy, SignedAt: tresc.Czas,
	}
	if tresc.Powod != "" {
		powod := tresc.Powod
		opis.Reason = &powod
	}
	if tresc.Miejsce != "" {
		miejsce := tresc.Miejsce
		opis.Location = &miejsce
	}

	surowyCertyfikat, err := base64.StdEncoding.DecodeString(tresc.Certyfikat)
	if err != nil {
		return opis, "certyfikat podpisu jest nieczytelny"
	}
	certyfikat, err := x509.ParseCertificate(surowyCertyfikat)
	if err != nil {
		return opis, "certyfikat podpisu jest nieczytelny: " + err.Error()
	}
	klucz, jest := certyfikat.PublicKey.(*rsa.PublicKey)
	if !jest {
		return opis, "certyfikat nie niesie klucza RSA"
	}
	surowyPodpis, err := base64.StdEncoding.DecodeString(tresc.Podpis)
	if err != nil {
		return opis, "podpis jest nieczytelny"
	}
	skrot, err := hex.DecodeString(tresc.Skrot)
	if err != nil {
		return opis, "skrót podpisanej treści jest nieczytelny"
	}
	if err := rsa.VerifyPKCS1v15(klucz, crypto.SHA256, skrot, surowyPodpis); err != nil {
		return opis, "podpis nie zgadza się z certyfikatem"
	}
	if tresc.Skrot != skrotBiezacy {
		return opis, "treść dokumentu zmieniła się po złożeniu podpisu"
	}
	opis.Valid = true
	return opis, ""
}

// ── Dane wrażliwe ───────────────────────────────────────────────────────────

// wzorzecWrazliwy jest jedną kategorią danych wraz z wyrażeniem, które ją
// rozpoznaje, i pewnością tego rozpoznania.
type wzorzecWrazliwy struct {
	kategoria string
	wyrazenie *regexp.Regexp
	pewnosc   float64
	sprawdz   func(string) bool
}

// wzorceWrazliwe wymienia rozpoznawane kategorie.
//
// Pewność jest różna, bo różna jest siła rozpoznania: adres poczty rozpoznaje
// się z kształtu i nie bywa czymś innym, a jedenaście cyfr obok siebie bywa
// numerem ewidencyjnym albo numerem zamówienia — i dlatego numer ewidencyjny
// przechodzi jeszcze sprawdzenie cyfry kontrolnej.
var wzorceWrazliwe = []wzorzecWrazliwy{
	{
		kategoria: "poczta",
		wyrazenie: regexp.MustCompile(`[\p{L}0-9._%+\-]+@[\p{L}0-9.\-]+\.[\p{L}]{2,}`),
		pewnosc:   0.95,
	},
	{
		kategoria: "karta-platnicza",
		wyrazenie: regexp.MustCompile(`\b(?:\d[ \-]?){12,18}\d\b`),
		pewnosc:   0.9,
		sprawdz:   przechodziLuhna,
	},
	{
		kategoria: "numer-ewidencyjny",
		wyrazenie: regexp.MustCompile(`\b\d{11}\b`),
		pewnosc:   0.9,
		sprawdz:   przechodziPesel,
	},
	{
		kategoria: "numer-podatkowy",
		wyrazenie: regexp.MustCompile(`\b\d{10}\b`),
		pewnosc:   0.85,
		sprawdz:   przechodziNip,
	},
	{
		kategoria: "rachunek-bankowy",
		wyrazenie: regexp.MustCompile(`\b[A-Z]{2}\d{2}(?:[ ]?\d{4}){4,7}\b`),
		pewnosc:   0.85,
	},
	{
		kategoria: "telefon",
		wyrazenie: regexp.MustCompile(`(?:\+\d{1,3}[ \-]?)?(?:\d{3}[ \-]?){2}\d{3}\b`),
		pewnosc:   0.6,
	},
}

// przechodziLuhna sprawdza cyfrę kontrolną numeru karty.
func przechodziLuhna(tekst string) bool {
	cyfry := make([]int, 0, len(tekst))
	for _, znak := range tekst {
		if znak >= '0' && znak <= '9' {
			cyfry = append(cyfry, int(znak-'0'))
		}
	}
	if len(cyfry) < 13 {
		return false
	}
	suma := 0
	podwajaj := false
	for i := len(cyfry) - 1; i >= 0; i-- {
		cyfra := cyfry[i]
		if podwajaj {
			cyfra *= 2
			if cyfra > 9 {
				cyfra -= 9
			}
		}
		suma += cyfra
		podwajaj = !podwajaj
	}
	return suma%10 == 0
}

// przechodziPesel sprawdza cyfrę kontrolną numeru ewidencyjnego.
func przechodziPesel(tekst string) bool {
	if len(tekst) != 11 {
		return false
	}
	wagi := []int{1, 3, 7, 9, 1, 3, 7, 9, 1, 3}
	suma := 0
	for i, waga := range wagi {
		suma += waga * int(tekst[i]-'0')
	}
	kontrolna := (10 - suma%10) % 10
	return kontrolna == int(tekst[10]-'0')
}

// przechodziNip sprawdza cyfrę kontrolną numeru podatkowego.
func przechodziNip(tekst string) bool {
	if len(tekst) != 10 {
		return false
	}
	wagi := []int{6, 5, 7, 2, 3, 4, 5, 6, 7}
	suma := 0
	for i, waga := range wagi {
		suma += waga * int(tekst[i]-'0')
	}
	kontrolna := suma % 11
	return kontrolna != 10 && kontrolna == int(tekst[9]-'0')
}

// WykryjWrazliwe oznacza w treści dokumentu dane wymagające redakcji.
//
// Czynność CZYTA i niczego nie zmienia — tak mówi kontrakt. Położenie jest
// liczone w znakach treści, więc wynik wchodzi wprost do żądania redakcji jako
// zakres, bez przeliczania po drodze.
func (a *adapterBezpieczenstwaStudia) WykryjWrazliwe(ctx context.Context,
	z shared.StudioSecuritySensitiveDetectRequest) (shared.StudioSecuritySensitiveDetectResponse, error) {

	if a.studio == nil {
		return shared.StudioSecuritySensitiveDetectResponse{}, odmowaBezpieczenstwa(
			shared.ErrorCodeInternalError,
			"rdzeń nie ma wpiętego repozytorium Studia — naprawa: podpiąć je przy "+
				"składaniu rdzenia")
	}

	var tresc string
	if z.VersionId != nil && strings.TrimSpace(*z.VersionId) != "" {
		wersja, err := a.studio.Wersja(ctx, strings.TrimSpace(*z.VersionId))
		if err != nil {
			return shared.StudioSecuritySensitiveDetectResponse{}, odmowaBezpieczenstwa(
				shared.ErrorCodeNotFound, "wersji "+*z.VersionId+" nie ma w rdzeniu")
		}
		if wersja.Tresc != nil {
			tresc = *wersja.Tresc
		}
	} else {
		dokument, err := a.studio.Dokument(ctx, strings.TrimSpace(z.DocumentId))
		if err != nil {
			return shared.StudioSecuritySensitiveDetectResponse{}, odmowaBezpieczenstwa(
				shared.ErrorCodeNotFound, "dokumentu "+z.DocumentId+" nie ma w rdzeniu")
		}
		if dokument.Tresc != nil {
			tresc = *dokument.Tresc
		}
	}

	brane := map[string]bool{}
	for _, kategoria := range z.Categories {
		brane[strings.ToLower(strings.TrimSpace(kategoria))] = true
	}
	if len(brane) > 0 {
		for _, wzorzec := range wzorceWrazliwe {
			if brane[wzorzec.kategoria] {
				delete(brane, wzorzec.kategoria)
			}
		}
		if len(brane) > 0 {
			nieznane := make([]string, 0, len(brane))
			for kategoria := range brane {
				nieznane = append(nieznane, kategoria)
			}
			sort.Strings(nieznane)
			znane := make([]string, 0, len(wzorceWrazliwe))
			for _, wzorzec := range wzorceWrazliwe {
				znane = append(znane, wzorzec.kategoria)
			}
			return shared.StudioSecuritySensitiveDetectResponse{}, odmowaBezpieczenstwa(
				shared.ErrorCodeValidationFailed,
				"kategorie nieznane: "+strings.Join(nieznane, ", ")+
					"; rozpoznawane są: "+strings.Join(znane, ", "))
		}
	}
	wybrane := map[string]bool{}
	for _, kategoria := range z.Categories {
		wybrane[strings.ToLower(strings.TrimSpace(kategoria))] = true
	}

	znaleziska := make([]shared.StudioSensitiveFinding, 0)
	for _, wzorzec := range wzorceWrazliwe {
		if len(wybrane) > 0 && !wybrane[wzorzec.kategoria] {
			continue
		}
		for _, zakres := range wzorzec.wyrazenie.FindAllStringIndex(tresc, -1) {
			fragment := tresc[zakres[0]:zakres[1]]
			if wzorzec.sprawdz != nil && !wzorzec.sprawdz(fragment) {
				continue
			}
			pewnosc := wzorzec.pewnosc
			znaleziska = append(znaleziska, shared.StudioSensitiveFinding{
				Category:   wzorzec.kategoria,
				Text:       fragment,
				RangeStart: zakres[0],
				RangeEnd:   zakres[1],
				Confidence: &pewnosc,
			})
		}
	}
	sort.SliceStable(znaleziska, func(i, j int) bool {
		return znaleziska[i].RangeStart < znaleziska[j].RangeStart
	})
	return shared.StudioSecuritySensitiveDetectResponse{Findings: znaleziska}, nil
}

// ── Rejestracja ─────────────────────────────────────────────────────────────

// BezpieczenstwoDokumentu wypełnia rodzinę `studio.security.*`.
type BezpieczenstwoDokumentu interface {
	ZaszyfrujPdf(ctx context.Context, z shared.StudioSecurityEncryptRequest) (shared.StudioSecurityEncryptResponse, error)
	WyczyscMetadane(ctx context.Context, z shared.StudioSecurityMetadataStripRequest) (shared.StudioSecurityMetadataStripResponse, error)
	Zredaguj(ctx context.Context, z shared.StudioSecurityRedactRequest) (shared.StudioSecurityRedactResponse, error)
	PodpiszPdf(ctx context.Context, z shared.StudioSecuritySignRequest) (shared.StudioSecuritySignResponse, error)
	SprawdzPodpisy(ctx context.Context, z shared.StudioSecuritySignVerifyRequest) (shared.StudioSecuritySignVerifyResponse, error)
	WykryjWrazliwe(ctx context.Context, z shared.StudioSecuritySensitiveDetectRequest) (shared.StudioSecuritySensitiveDetectResponse, error)
}

// zarejestrujBezpieczenstwoDokumentu wpina czynności rodziny bezpieczeństwa.
func zarejestrujBezpieczenstwoDokumentu(r *Rejestr, b BezpieczenstwoDokumentu) {
	if r == nil || b == nil {
		return
	}
	r.Zarejestruj(shared.CommandStudioSecurityEncrypt, obsluz(b.ZaszyfrujPdf))
	r.Zarejestruj(shared.CommandStudioSecurityMetadataStrip, obsluz(b.WyczyscMetadane))
	r.Zarejestruj(shared.CommandStudioSecurityRedact, obsluz(b.Zredaguj))
	r.Zarejestruj(shared.CommandStudioSecuritySign, obsluz(b.PodpiszPdf))
	r.Zarejestruj(shared.CommandStudioSecuritySignVerify, obsluz(b.SprawdzPodpisy))
	r.Zarejestruj(shared.CommandStudioSecuritySensitiveDetect, obsluz(b.WykryjWrazliwe))
}
