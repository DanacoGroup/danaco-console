// Odpowiedzialność pliku: wspólny warsztat części drukarskiej i marketingowej
// modułu Design — wykaz nośników, odczyt profili ICC, przekład profilu
// wydania, rachunek milimetrów oraz zakładanie zasobu magazynu.
package core

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"danacoconsole/server/internal/dane"
	"danacoconsole/shared"
)

const (
	// milimetryNaCal jest przelicznikiem rachunku drukarskiego. Rozdzielczość
	// podaje się w punktach na CAL, a materiał mierzy się w milimetrach —
	// bez tej stałej nie da się powiedzieć, ile pikseli ma baner.
	milimetryNaCal = 25.4

	// punktyNaCal jest jednostką strony PDF (punkt typograficzny). Wymiary
	// strony w `pdfcpu` podaje się w punktach, a Operator mierzy w milimetrach.
	punktyNaCal = 72.0

	// rozdzielczoscDrukuDomyslna jest rozdzielczością, wobec której mierzy
	// kontrola przeddrukowa, gdy profil jej nie podaje. 300 dpi to próg, poniżej
	// którego drukarnia arkuszowa oddaje obraz z widocznym rastrem.
	rozdzielczoscDrukuDomyslna = 300

	// granicaPikseliWydaniaDesignu chroni rdzeń przed żądaniem materiału nie do
	// zmieszczenia w pamięci: baner sześciometrowy w 300 dpi to ponad 70 tysięcy
	// pikseli boku. Granica jest liczona na CAŁE płótno, bo tam leży koszt.
	granicaPikseliWydaniaDesignu = 120_000_000

	// katalogiProfiliICC wylicza miejsca, w których profile ICC leżą na
	// systemach, na jakich stoi serwer produktu. Odczyt jest wyłącznie odczytem
	// katalogu — rdzeń profili nie rozkłada.
	katalogProfiliICCSystemowy = "/usr/share/color/icc"
	katalogProfiliICCLokalny   = "/usr/local/share/color/icc"
)

// nosnikiDruku oddaje nośniki w postaci kontraktu Designu, czytając wykaz
// wspólny rdzenia. Wykazu własnego ten moduł nie ma.
func nosnikiDruku() []shared.DesignPaperSize {
	wykaz := make([]shared.DesignPaperSize, 0, len(wykazNosnikowDruku))
	for _, nosnik := range wykazNosnikowDruku {
		wykaz = append(wykaz, shared.DesignPaperSize{
			Name:     nosnik.Nazwa,
			WidthMm:  nosnik.SzerokoscMm,
			HeightMm: nosnik.WysokoscMm,
			Family:   wskazTekstDesignu(nosnik.Rodzina),
		})
	}
	return wykaz
}

// wskazTekstDesignu oddaje wskaźnik na napis — pola opcjonalne kontraktu są
// wskaźnikami, a literału adresu wziąć nie można.
func wskazTekstDesignu(wartosc string) *string {
	return &wartosc
}

// NosnikiDruku oddaje nośniki spełniające zawężenie rodziną — obsługuje
// `design.print.paper.list`, licząc wymiary z wykazu wspólnego rdzenia.
func (a *adapterDesignu) NosnikiDruku(_ context.Context,
	z shared.DesignPrintPaperListRequest) (shared.DesignPrintPaperListResponse, error) {

	rodzina := ""
	if z.Family != nil {
		rodzina = strings.ToLower(strings.TrimSpace(*z.Family))
	}
	wykaz := nosnikiDruku()
	wybrane := make([]shared.DesignPaperSize, 0, len(wykaz))
	for _, nosnik := range wykaz {
		if rodzina != "" && (nosnik.Family == nil || strings.ToLower(*nosnik.Family) != rodzina) {
			continue
		}
		wybrane = append(wybrane, nosnik)
	}
	if rodzina != "" && len(wybrane) == 0 {
		return shared.DesignPrintPaperListResponse{}, bladWskazaniaDesignu(fmt.Sprintf(
			"komenda design.print.paper.list z rodziną %q, której rdzeń nie zna; rodziny znane: %s",
			*z.Family, strings.Join(rodzinyNosnikow(), ", ")))
	}
	return shared.DesignPrintPaperListResponse{Sizes: wybrane, Total: len(wybrane)}, nil
}

// rodzinyNosnikow oddaje rodziny rozmiarów w kolejności alfabetycznej — wykaz
// wchodzi do odmowy, żeby Operator dostał drogę wyjścia, a nie samo „nie".
func rodzinyNosnikow() []string {
	zbior := map[string]struct{}{}
	for _, nosnik := range nosnikiDruku() {
		if nosnik.Family != nil {
			zbior[*nosnik.Family] = struct{}{}
		}
	}
	nazwy := make([]string, 0, len(zbior))
	for nazwa := range zbior {
		nazwy = append(nazwy, nazwa)
	}
	sort.Strings(nazwy)
	return nazwy
}

// nosnikPoNazwie odnajduje nośnik po nazwie bez względu na wielkość liter,
// przeszukując wykaz wspólny rdzenia.
func nosnikPoNazwie(nazwa string) (shared.DesignPaperSize, bool) {
	szukana := strings.ToLower(strings.TrimSpace(nazwa))
	for _, nosnik := range nosnikiDruku() {
		if strings.ToLower(nosnik.Name) == szukana {
			return nosnik, true
		}
	}
	return shared.DesignPaperSize{}, false
}

// profileICCSerwera oddaje nazwy profili ICC leżących na tej maszynie. Wykaz
// jest pomiarem, nie zapowiedzią: pusty znaczy, że serwer nie ma ani jednego
// profilu.
func profileICCSerwera() []string {
	nazwy := []string{}
	for _, katalog := range []string{katalogProfiliICCSystemowy, katalogProfiliICCLokalny} {
		wpisy, err := os.ReadDir(katalog)
		if err != nil {
			continue
		}
		for _, wpis := range wpisy {
			if wpis.IsDir() {
				continue
			}
			rozszerzenie := strings.ToLower(filepath.Ext(wpis.Name()))
			if rozszerzenie != ".icc" && rozszerzenie != ".icm" {
				continue
			}
			nazwy = append(nazwy, wpis.Name())
		}
	}
	sort.Strings(nazwy)
	return nazwy
}

// profilDrukuKontraktu składa `DesignPrintProfile` kontraktu z wiersza
// zapisanego w bazie danych rdzenia.
func profilDrukuKontraktu(p dane.ProfilDrukuDesignu) shared.DesignPrintProfile {
	profil := shared.DesignPrintProfile{
		Id:                &p.Kod,
		WindowId:          &p.Okno,
		Name:              p.Nazwa,
		ColorSpace:        shared.DesignPrintColorSpace(p.PrzestrzenBarw),
		BleedMm:           p.SpadMm,
		CropMarks:         wskazLogiczneDesignu(p.ZnacznikiCiecia),
		RegistrationMarks: wskazLogiczneDesignu(p.ZnacznikiPasowania),
		ColorBar:          wskazLogiczneDesignu(p.PasekBarw),
		IccProfile:        p.ProfilICC,
		OverprintBlack:    wskazLogiczneDesignu(p.NadrukCzerni),
		PaperSize:         p.Nosnik,
	}
	if p.Norma != nil {
		norma := shared.DesignPrintStandard(*p.Norma)
		profil.Standard = &norma
	}
	if p.Rozdzielczosc != nil {
		wartosc := int(*p.Rozdzielczosc)
		profil.Dpi = &wartosc
	}
	return profil
}

// wskazLogiczneDesignu oddaje wskaźnik na wartość logiczną — pole opcjonalne
// kontraktu niesie wtedy rozstrzygnięcie wprost, a nie brak wartości.
func wskazLogiczneDesignu(wartosc bool) *bool {
	return &wartosc
}

// rozdzielczoscProfilu oddaje rozdzielczość, wobec której mierzy kontrola
// i w której wychodzi wydanie. Brak wskazania bierze próg drukarski.
func rozdzielczoscProfilu(profil shared.DesignPrintProfile) int {
	if profil.Dpi != nil && *profil.Dpi > 0 {
		return *profil.Dpi
	}
	return rozdzielczoscDrukuDomyslna
}

// spadProfilu oddaje spad w milimetrach. Brak wskazania znaczy „bez spadu" —
// rdzeń nie dokłada trzech milimetrów za Operatora, bo materiał bez spadu bywa
// zamierzony (wydruk biurowy), a spad dołożony po cichu zmienia wymiar strony.
func spadProfilu(profil shared.DesignPrintProfile) float64 {
	if profil.BleedMm != nil && *profil.BleedMm > 0 {
		return *profil.BleedMm
	}
	return 0
}

// pikseleZMilimetrow przelicza wymiar materiału na piksele przy zadanej
// rozdzielczości — rachunek, na którym stoi cała część drukarska.
func pikseleZMilimetrow(milimetry float64, dpi int) int {
	liczba := int(milimetry/milimetryNaCal*float64(dpi) + 0.5)
	if liczba < 1 {
		return 1
	}
	return liczba
}

// punktyZMilimetrow przelicza milimetry na punkty strony PDF, przez cale,
// jednostkę pośrednią obu miar.
func punktyZMilimetrow(milimetry float64) float64 {
	return milimetry / milimetryNaCal * punktyNaCal
}

// zalozZasobZBajtowDesignu odkłada bajty wytworzone przez rdzeń w magazynie
// i zakłada wiersz zasobu — droga wspólna dla kafli, wykresów, schematów,
// makiet i wydań do druku. Kolejność jest ta sama, co przy wniesieniu:
// najpierw bajty, potem wiersz.
func (a *adapterDesignu) zalozZasobZBajtowDesignu(ctx context.Context, okno, nazwa string,
	rodzaj shared.DesignAssetKind, format string, bajty []byte) (dane.ZasobDesignu, error) {

	if a.magazyn == nil {
		return dane.ZasobDesignu{}, bladZapisuZasobuDesignu(
			"magazyn treści nie jest wpięty — nie ma gdzie odłożyć bajtów wytworzonego materiału")
	}
	if len(bajty) == 0 {
		return dane.ZasobDesignu{}, bladWydaniaDesignu(
			"rdzeń nie wytworzył ani jednego bajtu — zasób bez treści nie ma czego pokazać")
	}
	suma := sha256.Sum256(bajty)
	odwolanie, err := a.magazyn.Zapisz(bajty, hex.EncodeToString(suma[:]))
	if err != nil {
		return dane.ZasobDesignu{}, bladZapisuZasobuDesignu(
			"nie można utrwalić treści wytworzonego materiału: " + err.Error())
	}

	// Format i wymiary pochodzą z bajtów utrwalonych, nie z zapowiedzi.
	zmierzonyFormat, szerokosc, wysokosc := rozpoznajObrazZasobu(odwolanie)
	wybranyFormat := pierwszyTekst(zmierzonyFormat, &format)

	zapisany, err := a.repozytorium.ZapiszZasob(ctx, dane.ZasobDesignu{
		Kod:       nowyIdentyfikator(przedrostekZasobuDesign),
		Okno:      okno,
		Nazwa:     &nazwa,
		Rodzaj:    string(rodzaj),
		Format:    wybranyFormat,
		URI:       &odwolanie,
		Szerokosc: szerokosc,
		Wysokosc:  wysokosc,
	})
	if err != nil {
		return dane.ZasobDesignu{}, bladDesignu(err)
	}
	return zapisany, nil
}

// zasobWytworzonyKontraktu składa odpowiedź kontraktu z wiersza świeżo
// założonego zasobu. Etykiet nie ma — zasób właśnie powstał i żadnych nie
// nadano; pusty wykaz jest tu prawdą, a nie brakiem odczytu.
func zasobWytworzonyKontraktu(zasob dane.ZasobDesignu) shared.DesignAsset {
	return zasobKontraktu(zasob, []string{})
}

// oknoWytworu rozstrzyga okno, w którym staje wytworzony materiał. Wskazanie
// żądania ma pierwszeństwo; brak bierze okno bytu źródłowego, bo materiał ma
// stanąć tam, skąd wyszedł, a nie w oknie bez nazwy.
func oknoWytworu(wskazane *string, zrodlowe string) string {
	if wskazane != nil && strings.TrimSpace(*wskazane) != "" {
		return strings.TrimSpace(*wskazane)
	}
	return zrodlowe
}
