// Odpowiedzialność pliku: część drukarska modułu Design — profile wydania
// (`design.print.profile.set`, `design.print.profile.list`), kontrola
// przeddrukowa (`design.print.preflight`), wydanie do druku
// (`design.print.export`) i podział materiału wielkoformatowego
// (`design.largeformat.tile`). Wykaz nośników, profile ICC serwera i rachunek
// milimetrów leżą w `adapter_modul_design_druk_wspolne.go`; wykresy i schematy
// w `adapter_modul_design_wykresy.go`.
//
// ── Jednostka kompozycji przy druku to MILIMETR ─────────────────────────────
// Kompozycja Design Board nie ma jednostki w kontrakcie — warstwa niesie liczby.
// Część drukarska czyta je jako MILIMETRY i jest to rozstrzygnięcie tego pliku,
// obowiązujące spójnie: kontrolę przeddrukową, wydanie i podział na kafle.
// Dzięki temu kompozycja 210×297 jest arkuszem A4, a nie prostokątem, którego
// rozmiaru nikt nie umie nazwać. Wyrys ekranowy (`design.board.export`) liczy te
// same liczby jako piksele i to jest w porządku: tam nie ma nośnika, więc nie ma
// czego mierzyć w milimetrach.
//
// ── `design.print.export` ODMAWIA przy wadzie o wadze błędu ─────────────────
// Plik nie do druku wydany jako gotowy do druku jest gorszy niż odmowa: idzie do
// drukarni, wraca po dniu i kosztuje nakład. Pominięcie kontroli jest jawnym
// wyborem Operatora (`skipPreflight`) i WRACA w odpowiedzi polem
// `preflightSkipped`, więc nikt nie powie potem, że nie wiedział.
//
// ── Zastrzeżenia w kolejności WAGI ──────────────────────────────────────────
// Kontrakt tak opisuje pole `issues` i tak Operator pracuje: najpierw naprawia
// to, co blokuje druk, potem to, co grozi jakością. Kolejność jest stabilna,
// żeby dwie kontrole tego samego materiału nie różniły się porządkiem.
//
// ── Profil ICC: mówimy, czego serwer NIE MA ─────────────────────────────────
// Rdzeń nie rozkłada profili ICC i nie przelicza barw przez nie. Wydanie w CMYK
// idzie więc bez osadzonego profilu i kontrola przeddrukowa mówi to jako
// zastrzeżenie — nie milczy. Powód stoi w nagłówku
// `adapter_modul_design_druk_wspolne.go`.
package core

import (
	"bytes"
	"context"
	"fmt"
	"image"
	"io"
	"math"
	"sort"
	"strings"

	"github.com/pdfcpu/pdfcpu/pkg/api"
	"golang.org/x/image/draw"

	"danacoconsole/server/internal/dane"
	"danacoconsole/shared"
)

const (
	// przedrostekProfiluDrukuDesign znakuje identyfikatory zewnętrzne profili.
	przedrostekProfiluDrukuDesign = "profil-druku-"

	// udzialRozdzielczosciBledoweDesignu wyznacza granicę między ostrzeżeniem
	// a błędem: materiał poniżej POŁOWY wymaganej rozdzielczości nie da się
	// wydrukować z sensem, a materiał poniżej wymaganej — da się, tylko gorzej.
	udzialRozdzielczosciBledowejDesignu = 0.5

	// granicaKafliDesignu chroni odpowiedź przed podziałem na dziesiątki tysięcy
	// kafli: baner sześciometrowy na kafle dziesięciocentymetrowe to sześćset
	// kafli w jednym rzędzie, a każdy kafel jest osobnym zasobem w magazynie.
	granicaKafliDesignu = 400
)

// UstawProfilDruku zapisuje profil wydania — obsługuje
// `design.print.profile.set`.
func (a *adapterDesignu) UstawProfilDruku(ctx context.Context,
	z shared.DesignPrintProfileSetRequest) (shared.DesignPrintProfileSetResponse, error) {

	if strings.TrimSpace(z.WindowId) == "" {
		return shared.DesignPrintProfileSetResponse{}, bladWskazaniaDesignu(
			"komenda design.print.profile.set bez wskazania okna")
	}
	if err := sprawdzProfilDrukuDesignu("design.print.profile.set", z.Profile); err != nil {
		return shared.DesignPrintProfileSetResponse{}, err
	}

	kod := nowyIdentyfikator(przedrostekProfiluDrukuDesign)
	if z.ProfileId != nil && strings.TrimSpace(*z.ProfileId) != "" {
		kod = strings.TrimSpace(*z.ProfileId)
		zastany, err := a.repozytorium.ProfilDrukuDesignuPoKodzie(ctx, kod)
		if err != nil {
			return shared.DesignPrintProfileSetResponse{}, bladNieznanegoProfiluDrukuDesignu(kod, err)
		}
		if zastany.Okno != z.WindowId {
			return shared.DesignPrintProfileSetResponse{}, bladWskazaniaDesignu(fmt.Sprintf(
				"profil %s należy do okna %s, a komenda design.print.profile.set przyszła z okna %s",
				kod, zastany.Okno, z.WindowId))
		}
	}

	wiersz := dane.ProfilDrukuDesignu{
		Kod:                kod,
		Okno:               z.WindowId,
		Nazwa:              z.Profile.Name,
		PrzestrzenBarw:     string(z.Profile.ColorSpace),
		SpadMm:             z.Profile.BleedMm,
		ZnacznikiCiecia:    z.Profile.CropMarks != nil && *z.Profile.CropMarks,
		ZnacznikiPasowania: z.Profile.RegistrationMarks != nil && *z.Profile.RegistrationMarks,
		PasekBarw:          z.Profile.ColorBar != nil && *z.Profile.ColorBar,
		ProfilICC:          z.Profile.IccProfile,
		NadrukCzerni:       z.Profile.OverprintBlack != nil && *z.Profile.OverprintBlack,
		Nosnik:             z.Profile.PaperSize,
	}
	if z.Profile.Standard != nil {
		norma := string(*z.Profile.Standard)
		wiersz.Norma = &norma
	}
	if z.Profile.Dpi != nil {
		rozdzielczosc := int64(*z.Profile.Dpi)
		wiersz.Rozdzielczosc = &rozdzielczosc
	}

	zapisany, err := a.repozytorium.ZapiszProfilDrukuDesignu(ctx, wiersz)
	if err != nil {
		return shared.DesignPrintProfileSetResponse{}, bladDesignu(err)
	}
	return shared.DesignPrintProfileSetResponse{Profile: profilDrukuKontraktu(zapisany)}, nil
}

// sprawdzProfilDrukuDesignu odrzuca nastawy spoza kontraktu i bezsensowne PRZED
// zapisem.
func sprawdzProfilDrukuDesignu(komenda string, profil shared.DesignPrintProfile) error {
	if err := sprawdzWyliczenieDesignu(komenda, "profile.colorSpace", profil.ColorSpace,
		shared.WartosciDesignPrintColorSpace()); err != nil {
		return err
	}
	if profil.Standard != nil {
		if err := sprawdzWyliczenieDesignu(komenda, "profile.standard", *profil.Standard,
			shared.WartosciDesignPrintStandard()); err != nil {
			return err
		}
	}
	if profil.BleedMm != nil && *profil.BleedMm < 0 {
		return bladWskazaniaDesignu(fmt.Sprintf(
			"komenda %s ze spadem %v mm: spad ujemny nie istnieje", komenda, *profil.BleedMm))
	}
	if profil.Dpi != nil && (*profil.Dpi < 1 || *profil.Dpi > 4800) {
		return bladWskazaniaDesignu(fmt.Sprintf(
			"komenda %s z rozdzielczością %d dpi: rdzeń przyjmuje od 1 do 4800 — powyżej materiał "+
				"nie zmieści się w pamięci, a żadna maszyna drukarska tego nie odda",
			komenda, *profil.Dpi))
	}
	if profil.PaperSize != nil && strings.TrimSpace(*profil.PaperSize) != "" {
		if _, znany := nosnikPoNazwie(*profil.PaperSize); !znany {
			wykaz := nosnikiDruku()
			nazwy := make([]string, 0, len(wykaz))
			for _, nosnik := range wykaz {
				nazwy = append(nazwy, nosnik.Name)
			}
			return bladWskazaniaDesignu(fmt.Sprintf(
				"komenda %s z nośnikiem %q, którego rdzeń nie zna; nośniki znane: %s",
				komenda, *profil.PaperSize, strings.Join(nazwy, ", ")))
		}
	}
	return nil
}

// ProfileDruku zwraca profile okna wraz z nośnikami i profilami ICC obecnymi na
// serwerze — obsługuje `design.print.profile.list`.
func (a *adapterDesignu) ProfileDruku(ctx context.Context,
	z shared.DesignPrintProfileListRequest) (shared.DesignPrintProfileListResponse, error) {

	if strings.TrimSpace(z.WindowId) == "" {
		return shared.DesignPrintProfileListResponse{}, bladWskazaniaDesignu(
			"komenda design.print.profile.list bez wskazania okna")
	}
	wiersze, err := a.repozytorium.ProfileDrukuDesignu(ctx, z.WindowId)
	if err != nil {
		return shared.DesignPrintProfileListResponse{}, bladDesignu(err)
	}
	profile := make([]shared.DesignPrintProfile, 0, len(wiersze))
	for _, wiersz := range wiersze {
		profile = append(profile, profilDrukuKontraktu(wiersz))
	}
	// `iccProfiles` jest POMIAREM tej maszyny: pusty wykaz znaczy, że wydanie
	// w CMYK pójdzie bez osadzonego profilu, i to jest odpowiedź, nie brak
	// odpowiedzi.
	return shared.DesignPrintProfileListResponse{
		Profiles: profile, PaperSizes: nosnikiDruku(), IccProfiles: profileICCSerwera(),
	}, nil
}

// KontrolaPrzeddrukowa mierzy materiał wobec profilu — obsługuje
// `design.print.preflight`.
func (a *adapterDesignu) KontrolaPrzeddrukowa(ctx context.Context,
	z shared.DesignPrintPreflightRequest) (shared.DesignPrintPreflightResponse, error) {

	profil, err := a.profilWydaniaDesignu(ctx, "design.print.preflight", z.ProfileId, z.Profile)
	if err != nil {
		return shared.DesignPrintPreflightResponse{}, err
	}
	zastrzezenia, err := a.zastrzezeniaPrzeddrukoweDesignu(ctx, "design.print.preflight",
		z.BoardId, z.AssetId, profil)
	if err != nil {
		return shared.DesignPrintPreflightResponse{}, err
	}
	bledow, ostrzezen := zlicZastrzezeniaDesignu(zastrzezenia)
	return shared.DesignPrintPreflightResponse{
		Issues: zastrzezenia, Errors: bledow, Warnings: ostrzezen, Ready: bledow == 0,
	}, nil
}

// profilWydaniaDesignu rozstrzyga profil, wobec którego mierzy kontrola i w którym
// wychodzi wydanie.
//
// Wskazanie wprost (`profile`) bije profil z bazy: Operator, który podał nastawy
// w żądaniu, chce tych nastaw. Brak jednego i drugiego bierze nastawy domyślne —
// CMYK w rozdzielczości drukarskiej, bez spadu, bo spad dołożony po cichu
// zmieniałby wymiar strony (powód przy `spadProfilu`).
func (a *adapterDesignu) profilWydaniaDesignu(ctx context.Context, komenda string,
	kod *string, wskazany *shared.DesignPrintProfile) (shared.DesignPrintProfile, error) {

	if wskazany != nil {
		if err := sprawdzProfilDrukuDesignu(komenda, *wskazany); err != nil {
			return shared.DesignPrintProfile{}, err
		}
		return *wskazany, nil
	}
	if kod != nil && strings.TrimSpace(*kod) != "" {
		wiersz, err := a.repozytorium.ProfilDrukuDesignuPoKodzie(ctx, strings.TrimSpace(*kod))
		if err != nil {
			return shared.DesignPrintProfile{}, bladNieznanegoProfiluDrukuDesignu(*kod, err)
		}
		return profilDrukuKontraktu(wiersz), nil
	}
	return shared.DesignPrintProfile{ColorSpace: shared.DesignPrintColorSpaceCmyk}, nil
}

// zastrzezeniaPrzeddrukoweDesignu składa wykaz zastrzeżeń w kolejności wagi.
func (a *adapterDesignu) zastrzezeniaPrzeddrukoweDesignu(ctx context.Context, komenda string,
	kompozycja, zasob *string,
	profil shared.DesignPrintProfile) ([]shared.DesignPreflightIssue, error) {

	maKompozycje := kompozycja != nil && strings.TrimSpace(*kompozycja) != ""
	maZasob := zasob != nil && strings.TrimSpace(*zasob) != ""
	if !maKompozycje && !maZasob {
		return nil, bladWskazaniaDesignu(fmt.Sprintf(
			"komenda %s bez wskazania kompozycji ani zasobu: kontrola przeddrukowa mierzy "+
				"MATERIAŁ, a nie sam profil", komenda))
	}

	zastrzezenia := []shared.DesignPreflightIssue{}
	rozdzielczosc := rozdzielczoscProfilu(profil)

	if maZasob {
		wiersz, err := a.repozytorium.Zasob(ctx, strings.TrimSpace(*zasob))
		if err != nil {
			return nil, bladNieznanegoZasobuDesignu(*zasob, err)
		}
		zastrzezenia = append(zastrzezenia,
			zastrzezeniaZasobuPrzeddrukoweDesignu(wiersz, profil, rozdzielczosc)...)
	}

	if maKompozycje {
		wiersz, err := a.repozytorium.Kompozycja(ctx, strings.TrimSpace(*kompozycja))
		if err != nil {
			return nil, bladNieznanejKompozycjiDesignu(*kompozycja, err)
		}
		warstwy, err := a.repozytorium.Warstwy(ctx, wiersz.ID)
		if err != nil {
			return nil, bladDesignu(err)
		}
		if len(warstwy) == 0 {
			zastrzezenia = append(zastrzezenia, shared.DesignPreflightIssue{
				Severity: shared.DesignPreflightSeverityBlad,
				Code:     "kompozycja-bez-warstw",
				Message: fmt.Sprintf("kompozycja %s nie ma ani jednej warstwy — wydanie byłoby "+
					"pustym arkuszem", wiersz.Kod),
			})
		}
		for _, warstwa := range warstwy {
			zastrzezenia = append(zastrzezenia, a.zastrzezeniaWarstwyPrzeddrukoweDesignu(ctx,
				warstwa, profil, rozdzielczosc)...)
		}
		zastrzezenia = append(zastrzezenia,
			zastrzezeniaSpaduDesignu(warstwy, profil)...)
	}

	zastrzezenia = append(zastrzezenia, zastrzezeniaProfiluDesignu(profil)...)
	uporzadkujZastrzezeniaDesignu(zastrzezenia)
	return zastrzezenia, nil
}

// zastrzezeniaZasobuPrzeddrukoweDesignu mierzy jeden zasób wobec profilu.
//
// Rachunek jest prawdziwym pomiarem: rozdzielczość skuteczna to piksele zasobu
// rozłożone na wymiar nośnika. Bez nośnika w profilu nie ma czego mierzyć
// i rdzeń mówi to wprost, zamiast wystawić ocenę bez podstawy.
func zastrzezeniaZasobuPrzeddrukoweDesignu(zasob dane.ZasobDesignu,
	profil shared.DesignPrintProfile, rozdzielczosc int) []shared.DesignPreflightIssue {

	zastrzezenia := []shared.DesignPreflightIssue{}
	kod := zasob.Kod
	if zasob.URI == nil || strings.TrimSpace(*zasob.URI) == "" {
		return append(zastrzezenia, shared.DesignPreflightIssue{
			Severity: shared.DesignPreflightSeverityBlad,
			Code:     "zasob-bez-tresci",
			Message: fmt.Sprintf("zasób %s nie ma treści w magazynie — nie ma czego wydrukować",
				kod),
			AssetId: &kod,
		})
	}
	if zasob.Szerokosc == nil || zasob.Wysokosc == nil {
		return append(zastrzezenia, shared.DesignPreflightIssue{
			Severity: shared.DesignPreflightSeverityOstrzezenie,
			Code:     "zasob-bez-wymiarow",
			Message: fmt.Sprintf("rdzeń nie zmierzył wymiarów zasobu %s (format spoza png, jpeg, "+
				"gif) — rozdzielczości skutecznej nie da się policzyć", kod),
			AssetId: &kod,
		})
	}
	if profil.PaperSize == nil || strings.TrimSpace(*profil.PaperSize) == "" {
		return append(zastrzezenia, shared.DesignPreflightIssue{
			Severity: shared.DesignPreflightSeverityInformacja,
			Code:     "profil-bez-nosnika",
			Message: "profil nie wskazuje nośnika, więc rozdzielczości skutecznej zasobu nie ma " +
				"do czego odnieść — wskaż paperSize, żeby kontrola zmierzyła materiał",
			AssetId: &kod,
		})
	}
	nosnik, znany := nosnikPoNazwie(*profil.PaperSize)
	if !znany {
		return zastrzezenia
	}
	skuteczna := float64(*zasob.Szerokosc) / nosnik.WidthMm * milimetryNaCal
	if pionowa := float64(*zasob.Wysokosc) / nosnik.HeightMm * milimetryNaCal; pionowa < skuteczna {
		skuteczna = pionowa
	}
	if skuteczna >= float64(rozdzielczosc) {
		return zastrzezenia
	}
	waga := shared.DesignPreflightSeverity(shared.DesignPreflightSeverityOstrzezenie)
	if skuteczna < float64(rozdzielczosc)*udzialRozdzielczosciBledowejDesignu {
		waga = shared.DesignPreflightSeverity(shared.DesignPreflightSeverityBlad)
	}
	zmierzona := fmt.Sprintf("%.0f dpi", skuteczna)
	oczekiwana := fmt.Sprintf("%d dpi", rozdzielczosc)
	return append(zastrzezenia, shared.DesignPreflightIssue{
		Severity: waga,
		Code:     "rozdzielczosc-ponizej-progu",
		Message: fmt.Sprintf("zasób %s rozciągnięty na nośnik %s daje %.0f dpi, a profil wymaga "+
			"%d dpi — wydruk pokaże raster", kod, nosnik.Name, skuteczna, rozdzielczosc),
		AssetId: &kod, Measured: &zmierzona, Expected: &oczekiwana,
	})
}

// zastrzezeniaWarstwyPrzeddrukoweDesignu mierzy jedną warstwę kompozycji.
//
// Rozdzielczość skuteczna warstwy liczy się z pikseli jej zasobu i z jej wymiaru
// na kompozycji, czytanego w MILIMETRACH (nagłówek pliku). To jest ta liczba,
// którą drukarnia mierzy jako pierwszą.
func (a *adapterDesignu) zastrzezeniaWarstwyPrzeddrukoweDesignu(ctx context.Context,
	warstwa dane.WarstwaKompozycji, profil shared.DesignPrintProfile,
	rozdzielczosc int) []shared.DesignPreflightIssue {

	kodWarstwy := warstwa.Kod
	if warstwa.ZasobID == nil || strings.TrimSpace(*warstwa.ZasobID) == "" {
		// Warstwa bez zasobu jest w kompozycji normalna (ramka, prowadnica), więc
		// to informacja, nie ostrzeżenie.
		return []shared.DesignPreflightIssue{{
			Severity: shared.DesignPreflightSeverityInformacja,
			Code:     "warstwa-bez-zasobu",
			Message:  fmt.Sprintf("warstwa %s nie wskazuje zasobu — w wydaniu nie będzie jej widać", kodWarstwy),
			LayerId:  &kodWarstwy,
		}}
	}
	kodZasobu := strings.TrimSpace(*warstwa.ZasobID)
	zasob, err := a.repozytorium.Zasob(ctx, kodZasobu)
	if err != nil {
		if czyBrakZasobuDesignu(err) {
			return []shared.DesignPreflightIssue{{
				Severity: shared.DesignPreflightSeverityBlad,
				Code:     "warstwa-wskazuje-zasob-usuniety",
				Message: fmt.Sprintf("warstwa %s wskazuje zasób %s, którego nie ma w tym rdzeniu "+
					"— w wydaniu zostanie po niej dziura", kodWarstwy, kodZasobu),
				LayerId: &kodWarstwy, AssetId: &kodZasobu,
			}}
		}
		return nil
	}
	if zasob.URI == nil || strings.TrimSpace(*zasob.URI) == "" {
		return []shared.DesignPreflightIssue{{
			Severity: shared.DesignPreflightSeverityBlad,
			Code:     "zasob-bez-tresci",
			Message: fmt.Sprintf("warstwa %s wskazuje zasób %s bez bajtów w magazynie",
				kodWarstwy, kodZasobu),
			LayerId: &kodWarstwy, AssetId: &kodZasobu,
		}}
	}
	if zasob.Szerokosc == nil || warstwa.Szerokosc == nil || *warstwa.Szerokosc <= 0 {
		return nil
	}
	skuteczna := float64(*zasob.Szerokosc) / *warstwa.Szerokosc * milimetryNaCal
	if skuteczna >= float64(rozdzielczosc) {
		return nil
	}
	waga := shared.DesignPreflightSeverity(shared.DesignPreflightSeverityOstrzezenie)
	if skuteczna < float64(rozdzielczosc)*udzialRozdzielczosciBledowejDesignu {
		waga = shared.DesignPreflightSeverity(shared.DesignPreflightSeverityBlad)
	}
	zmierzona := fmt.Sprintf("%.0f dpi", skuteczna)
	oczekiwana := fmt.Sprintf("%d dpi", rozdzielczosc)
	return []shared.DesignPreflightIssue{{
		Severity: waga,
		Code:     "rozdzielczosc-ponizej-progu",
		Message: fmt.Sprintf("warstwa %s ma %.0f dpi (zasób %d px na %.1f mm), a profil wymaga "+
			"%d dpi", kodWarstwy, skuteczna, *zasob.Szerokosc, *warstwa.Szerokosc, rozdzielczosc),
		LayerId: &kodWarstwy, AssetId: &kodZasobu,
		Measured: &zmierzona, Expected: &oczekiwana,
	}}
}

// zastrzezeniaSpaduDesignu sprawdza, czy materiał dochodzi do spadu.
//
// Rachunek jest prawdziwy: warstwa, która kończy się dokładnie na krawędzi
// nośnika, po obcięciu zostawi biały pasek, bo maszyna cięcia ma tolerancję.
// Spad wymaga, żeby treść WYCHODZIŁA poza krawędź.
func zastrzezeniaSpaduDesignu(warstwy []dane.WarstwaKompozycji,
	profil shared.DesignPrintProfile) []shared.DesignPreflightIssue {

	spad := spadProfilu(profil)
	if spad <= 0 || profil.PaperSize == nil || strings.TrimSpace(*profil.PaperSize) == "" {
		return nil
	}
	nosnik, znany := nosnikPoNazwie(*profil.PaperSize)
	if !znany || len(warstwy) == 0 {
		return nil
	}
	// Prostokąt obejmujący treść: warstwy bez położenia liczone od zera.
	lewa, gora := math.MaxFloat64, math.MaxFloat64
	prawa, dol := -math.MaxFloat64, -math.MaxFloat64
	for _, warstwa := range warstwy {
		x, y := 0.0, 0.0
		if warstwa.X != nil {
			x = *warstwa.X
		}
		if warstwa.Y != nil {
			y = *warstwa.Y
		}
		szerokosc, wysokosc := 0.0, 0.0
		if warstwa.Szerokosc != nil {
			szerokosc = *warstwa.Szerokosc
		}
		if warstwa.Wysokosc != nil {
			wysokosc = *warstwa.Wysokosc
		}
		lewa = mniejszaDesignu(lewa, x)
		gora = mniejszaDesignu(gora, y)
		prawa = wiekszaDesignu(prawa, x+szerokosc)
		dol = wiekszaDesignu(dol, y+wysokosc)
	}
	brakujace := []string{}
	if lewa > -spad {
		brakujace = append(brakujace, "lewa")
	}
	if gora > -spad {
		brakujace = append(brakujace, "górna")
	}
	if prawa < nosnik.WidthMm+spad {
		brakujace = append(brakujace, "prawa")
	}
	if dol < nosnik.HeightMm+spad {
		brakujace = append(brakujace, "dolna")
	}
	if len(brakujace) == 0 {
		return nil
	}
	zmierzona := fmt.Sprintf("treść od %.1f;%.1f do %.1f;%.1f mm", lewa, gora, prawa, dol)
	oczekiwana := fmt.Sprintf("od %.1f;%.1f do %.1f;%.1f mm",
		-spad, -spad, nosnik.WidthMm+spad, nosnik.HeightMm+spad)
	return []shared.DesignPreflightIssue{{
		Severity: shared.DesignPreflightSeverityOstrzezenie,
		Code:     "tresc-nie-dochodzi-do-spadu",
		Message: fmt.Sprintf("profil żąda spadu %.1f mm, a treść nie wychodzi poza krawędź na "+
			"stronach: %s — po obcięciu zostanie tam biały pasek",
			spad, strings.Join(brakujace, ", ")),
		Measured: &zmierzona, Expected: &oczekiwana,
	}}
}

// zastrzezeniaProfiluDesignu wylicza zastrzeżenia wynikające z samego profilu
// i z tego, czym rdzeń NAPRAWDĘ dysponuje.
func zastrzezeniaProfiluDesignu(profil shared.DesignPrintProfile) []shared.DesignPreflightIssue {
	zastrzezenia := []shared.DesignPreflightIssue{}

	if profil.ColorSpace == shared.DesignPrintColorSpaceCmyk ||
		profil.ColorSpace == shared.DesignPrintColorSpaceSpot {

		obecne := profileICCSerwera()
		if profil.IccProfile != nil && strings.TrimSpace(*profil.IccProfile) != "" {
			znany := false
			for _, nazwa := range obecne {
				if strings.EqualFold(nazwa, strings.TrimSpace(*profil.IccProfile)) {
					znany = true
					break
				}
			}
			if !znany {
				zmierzona := fmt.Sprintf("profili na serwerze: %d", len(obecne))
				zastrzezenia = append(zastrzezenia, shared.DesignPreflightIssue{
					Severity: shared.DesignPreflightSeverityOstrzezenie,
					Code:     "profil-icc-nieobecny",
					Message: fmt.Sprintf("profil wskazuje ICC %q, którego na tym serwerze nie ma "+
						"— wydanie pójdzie bez osadzonego profilu", *profil.IccProfile),
					Measured: &zmierzona,
				})
			}
		}
		// Rdzeń nie rozdziela barw przez profil ICC. To jest granica produktu
		// i mówi się ją WPROST, jako informacja przy każdym wydaniu w CMYK —
		// drukarnia ma wiedzieć, że dostała przeliczenie naiwne.
		zastrzezenia = append(zastrzezenia, shared.DesignPreflightIssue{
			Severity: shared.DesignPreflightSeverityInformacja,
			Code:     "cmyk-bez-rozdzialu-icc",
			Message: "rdzeń przelicza RGB na CMYK WPROST, bez profilu ICC — wartości są punktem " +
				"wyjścia dla drukarni, a nie barwą rozdzieloną pod maszynę",
		})
	}
	if profil.Standard != nil && *profil.Standard != shared.DesignPrintStandardBrak {
		zastrzezenia = append(zastrzezenia, shared.DesignPreflightIssue{
			Severity: shared.DesignPreflightSeverityInformacja,
			Code:     "norma-bez-weryfikacji",
			Message: fmt.Sprintf("profil żąda normy %s; rdzeń składa dokument PDF biblioteką "+
				"wkompilowaną i NIE weryfikuje zgodności z normą — weryfikacja wymaga narzędzia "+
				"spoza instalki", string(*profil.Standard)),
		})
	}
	if profil.OverprintBlack != nil && *profil.OverprintBlack {
		zastrzezenia = append(zastrzezenia, shared.DesignPreflightIssue{
			Severity: shared.DesignPreflightSeverityInformacja,
			Code:     "nadruk-czerni-poza-rdzeniem",
			Message: "profil żąda nadruku czerni; ustawienie nadruku należy do rozdziału barw " +
				"w drukarni, a rdzeń przekazuje je jako nastawę profilu, nie jako właściwość pliku",
		})
	}
	return zastrzezenia
}

// uporzadkujZastrzezeniaDesignu porządkuje zastrzeżenia po wadze, a w obrębie
// wagi po kodzie — kolejność stabilna, żeby dwie kontrole tego samego materiału
// nie różniły się porządkiem.
func uporzadkujZastrzezeniaDesignu(zastrzezenia []shared.DesignPreflightIssue) {
	waga := func(w shared.DesignPreflightSeverity) int {
		switch w {
		case shared.DesignPreflightSeverityBlad:
			return 0
		case shared.DesignPreflightSeverityOstrzezenie:
			return 1
		}
		return 2
	}
	sort.SliceStable(zastrzezenia, func(i, j int) bool {
		if waga(zastrzezenia[i].Severity) != waga(zastrzezenia[j].Severity) {
			return waga(zastrzezenia[i].Severity) < waga(zastrzezenia[j].Severity)
		}
		return zastrzezenia[i].Code < zastrzezenia[j].Code
	})
}

// zlicZastrzezeniaDesignu liczy zastrzeżenia o wadze błędu i ostrzeżenia.
func zlicZastrzezeniaDesignu(zastrzezenia []shared.DesignPreflightIssue) (int, int) {
	bledow, ostrzezen := 0, 0
	for _, zastrzezenie := range zastrzezenia {
		switch zastrzezenie.Severity {
		case shared.DesignPreflightSeverityBlad:
			bledow++
		case shared.DesignPreflightSeverityOstrzezenie:
			ostrzezen++
		}
	}
	return bledow, ostrzezen
}

// WydajDoDruku wydaje materiał do druku — obsługuje `design.print.export`.
func (a *adapterDesignu) WydajDoDruku(ctx context.Context,
	z shared.DesignPrintExportRequest) (shared.DesignPrintExportResponse, error) {

	format := normalizujFormatWydaniaDesignu(z.Format)
	switch format {
	case "pdf", "tiff", "eps":
	default:
		return shared.DesignPrintExportResponse{}, bladWskazaniaDesignu(fmt.Sprintf(
			"komenda design.print.export z formatem %q: wydanie do druku wychodzi jako pdf, "+
				"tiff albo eps", z.Format))
	}
	profil, err := a.profilWydaniaDesignu(ctx, "design.print.export", z.ProfileId, z.Profile)
	if err != nil {
		return shared.DesignPrintExportResponse{}, err
	}

	// Publikacja wielostronicowa: szablon materiału o wielu stronach wychodzi
	// JEDNYM plikiem, tą samą drogą i przez tę samą kontrolę przeddrukową —
	// z odmową przy wadzie o wadze błędu obowiązującą także tutaj.
	if z.TemplateId != nil && strings.TrimSpace(*z.TemplateId) != "" {
		return a.zlozWydaniePublikacjiDesignu(ctx, z, format, profil)
	}

	pominieta := z.SkipPreflight != nil && *z.SkipPreflight
	zastrzezenia := []shared.DesignPreflightIssue{}
	if !pominieta {
		zastrzezenia, err = a.zastrzezeniaPrzeddrukoweDesignu(ctx, "design.print.export",
			z.BoardId, z.AssetId, profil)
		if err != nil {
			return shared.DesignPrintExportResponse{}, err
		}
		bledow, _ := zlicZastrzezeniaDesignu(zastrzezenia)
		if bledow > 0 {
			// ODMOWA, nie wydanie z ostrzeżeniem: plik nie do druku wydany jako
			// gotowy do druku idzie do drukarni i kosztuje nakład (nagłówek pliku).
			return shared.DesignPrintExportResponse{}, bladWskazaniaDesignu(fmt.Sprintf(
				"kontrola przeddrukowa znalazła %d wad o wadze błędu — rdzeń nie wyda pliku "+
					"nie do druku jako gotowego do druku; pierwsza wada: %s; naprawa: usunąć wady "+
					"albo świadomie pominąć kontrolę polem skipPreflight (pominięcie wraca "+
					"w odpowiedzi)", bledow, pierwszyBladZastrzezenDesignu(zastrzezenia)))
		}
	}

	bajty, typTresci, nazwa, okno, stron, err := a.zlozWydanieDrukarskieDesignu(ctx, z, format, profil)
	if err != nil {
		return shared.DesignPrintExportResponse{}, err
	}
	zasob, err := a.zalozZasobZBajtowDesignu(ctx, oknoWytworu(z.WindowId, okno), nazwa,
		rodzajWydaniaDrukarskiegoDesignu(format), format, bajty)
	if err != nil {
		return shared.DesignPrintExportResponse{}, err
	}
	return shared.DesignPrintExportResponse{
		Asset: zasobWytworzonyKontraktu(zasob), FileName: nazwa, MediaType: typTresci,
		SizeBytes: len(bajty), PreflightSkipped: pominieta, Issues: zastrzezenia,
		PageCount: &stron,
	}, nil
}

// zlozWydaniePublikacjiDesignu wydaje publikację wielostronicową jednym plikiem.
//
// Kontrola przeddrukowa obowiązuje TAK SAMO: każda strona jest arkuszem i każda
// przechodzi te same pomiary, a wada o wadze błędu na jednej stronie odmawia
// całego wydania. Publikacja przepuszczona z jedną stroną nie do druku wraca
// z drukarni tak samo jak pojedynczy arkusz — tylko drożej.
func (a *adapterDesignu) zlozWydaniePublikacjiDesignu(ctx context.Context,
	z shared.DesignPrintExportRequest, format string,
	profil shared.DesignPrintProfile) (shared.DesignPrintExportResponse, error) {

	szablon, err := a.repozytorium.SzablonMaterialuDesignuPoKodzie(ctx,
		strings.TrimSpace(*z.TemplateId))
	if err != nil {
		return shared.DesignPrintExportResponse{}, bladNieznanegoSzablonuMaterialu(*z.TemplateId, err)
	}
	strony, err := a.repozytorium.StronySzablonuMaterialuDesignu(ctx, szablon.ID)
	if err != nil {
		return shared.DesignPrintExportResponse{}, bladDesignu(err)
	}
	if len(strony) == 0 {
		return shared.DesignPrintExportResponse{}, bladWskazaniaDesignu(fmt.Sprintf(
			"szablon %s nie ma ani jednej strony — publikacja wielostronicowa powstaje z szablonu "+
				"o wielu stronach (design.template.save z polem pages); szablon bez stron wydaje "+
				"się przez wskazanie kompozycji, nie szablonu", szablon.Kod))
	}
	oprawa := shared.DesignPrintBinding(shared.DesignPrintBindingBrak)
	if z.Binding != nil {
		if err := sprawdzWyliczenieDesignu("design.print.export", "binding", *z.Binding,
			shared.WartosciDesignPrintBinding()); err != nil {
			return shared.DesignPrintExportResponse{}, err
		}
		oprawa = *z.Binding
	}
	uszeregowane, err := uszeregujStronyPublikacjiDesignu(strony, z.PageOrder, oprawa)
	if err != nil {
		return shared.DesignPrintExportResponse{}, err
	}

	pominieta := z.SkipPreflight != nil && *z.SkipPreflight
	zastrzezenia := []shared.DesignPreflightIssue{}
	rozdzielczosc := rozdzielczoscProfilu(profil)
	skala := float64(rozdzielczosc) / milimetryNaCal

	obrazy := make([]image.Image, 0, len(uszeregowane))
	for numer, strona := range uszeregowane {
		warstwy, err := a.repozytorium.WarstwyStronySzablonuMaterialuDesignu(ctx, strona.ID)
		if err != nil {
			return shared.DesignPrintExportResponse{}, bladDesignu(err)
		}
		if !pominieta {
			numerStrony := numer + 1
			zastrzezeniaStrony := zastrzezeniaSpaduDesignu(warstwy, profil)
			for _, warstwa := range warstwy {
				zastrzezeniaStrony = append(zastrzezeniaStrony,
					a.zastrzezeniaWarstwyPrzeddrukoweDesignu(ctx, warstwa, profil, rozdzielczosc)...)
			}
			if len(warstwy) == 0 {
				zastrzezeniaStrony = append(zastrzezeniaStrony, shared.DesignPreflightIssue{
					Severity: shared.DesignPreflightSeverityBlad,
					Code:     "strona-bez-warstw",
					Message: fmt.Sprintf("strona %d publikacji nie ma ani jednej warstwy — "+
						"w pliku zostałaby po niej pusta kartka", strona.Numer),
				})
			}
			// Numer strony wchodzi w zastrzeżenie: bez niego Operator dostawałby
			// wykaz wad publikacji dwudziestostronicowej bez wskazania, której
			// strony dotyczą.
			for numerZastrzezenia := range zastrzezeniaStrony {
				zastrzezeniaStrony[numerZastrzezenia].Page = &numerStrony
			}
			zastrzezenia = append(zastrzezenia, zastrzezeniaStrony...)
		}

		kafle, pominietych, err := a.kafleWyrysuDesignu(ctx, warstwy)
		if err != nil {
			return shared.DesignPrintExportResponse{}, err
		}
		// Strona bez bajtów wychodzi jako CZYSTA kartka o wymiarach szablonu, a nie
		// jako odmowa: publikacja ma strony celowo puste (wakat, strona redakcyjna),
		// a kontrola przeddrukowa powiedziała już wyżej, że warstw tam nie ma.
		obszar := shared.DesignBoardRegion{Width: szablon.Szerokosc, Height: szablon.Wysokosc}
		if len(kafle) > 0 {
			obszar = obszarWyrysuDesignu(kafle, &obszar)
		}
		_ = pominietych
		if pikseli := obszar.Width * skala * obszar.Height * skala; pikseli >
			granicaPikseliWydaniaDesignu {

			return shared.DesignPrintExportResponse{}, bladWskazaniaDesignu(fmt.Sprintf(
				"strona %d o wymiarach %.0f×%.0f mm w %d dpi to %.0f pikseli płótna, a granica "+
					"wydania to %d", strona.Numer, obszar.Width, obszar.Height, rozdzielczosc,
				pikseli, granicaPikseliWydaniaDesignu))
		}
		obrazy = append(obrazy, zlozWyrysRastrowyDesignu(kafle, obszar, skala))
	}

	uporzadkujZastrzezeniaDesignu(zastrzezenia)
	bledow, _ := zlicZastrzezeniaDesignu(zastrzezenia)
	if !pominieta && bledow > 0 {
		return shared.DesignPrintExportResponse{}, bladWskazaniaDesignu(fmt.Sprintf(
			"kontrola przeddrukowa publikacji %s znalazła %d wad o wadze błędu — rdzeń nie wyda "+
				"pliku nie do druku jako gotowego do druku; pierwsza wada: %s",
			szablon.Kod, bledow, pierwszyBladZastrzezenDesignu(zastrzezenia)))
	}

	bajty, typTresci, err := zakodujWydanieDrukarskieDesignu(obrazy, format)
	if err != nil {
		return shared.DesignPrintExportResponse{}, bladWydaniaDesignu(err.Error())
	}
	nazwa := oczyscNazwePlikuDesignu(szablon.Nazwa) + "." + rozszerzenieWydaniaDesignu(format)
	zasob, err := a.zalozZasobZBajtowDesignu(ctx, oknoWytworu(z.WindowId, szablon.Okno), nazwa,
		rodzajWydaniaDrukarskiegoDesignu(format), format, bajty)
	if err != nil {
		return shared.DesignPrintExportResponse{}, err
	}
	stron := len(obrazy)
	return shared.DesignPrintExportResponse{
		Asset: zasobWytworzonyKontraktu(zasob), FileName: nazwa, MediaType: typTresci,
		SizeBytes: len(bajty), PreflightSkipped: pominieta, Issues: zastrzezenia,
		PageCount: &stron,
	}, nil
}

// uszeregujStronyPublikacjiDesignu ustawia strony w kolejności wydania.
//
// Kolejność wskazana wprost (`pageOrder`) bije wszystko: Operator, który podał
// numery, chce tej kolejności. Brak wskazania bierze kolejność numerów strony.
// Oprawa zeszytowa wymaga liczby stron podzielnej przez cztery — arkusz zgięty
// na pół daje cztery strony, więc publikacja o dwudziestu dwóch stronach nie da
// się w ten sposób zszyć i jest to ODMOWA, a nie ciche dołożenie dwóch wakatów.
func uszeregujStronyPublikacjiDesignu(strony []dane.StronaSzablonuMaterialuDesignu,
	kolejnosc []int, oprawa shared.DesignPrintBinding) (
	[]dane.StronaSzablonuMaterialuDesignu, error) {

	if oprawa == shared.DesignPrintBindingZeszytowa && len(strony)%4 != 0 {
		return nil, bladWskazaniaDesignu(fmt.Sprintf(
			"oprawa zeszytowa wymaga liczby stron podzielnej przez cztery (arkusz zgięty na pół "+
				"daje cztery strony), a publikacja ma %d — rdzeń nie dokłada wakatów za Operatora, "+
				"bo strona pusta w środku książki jest rozstrzygnięciem, nie zaokrągleniem",
			len(strony)))
	}
	if len(kolejnosc) == 0 {
		return strony, nil
	}
	poNumerze := make(map[int]dane.StronaSzablonuMaterialuDesignu, len(strony))
	for _, strona := range strony {
		poNumerze[strona.Numer] = strona
	}
	uszeregowane := make([]dane.StronaSzablonuMaterialuDesignu, 0, len(kolejnosc))
	uzyte := map[int]bool{}
	for _, numer := range kolejnosc {
		strona, jest := poNumerze[numer]
		if !jest {
			return nil, bladWskazaniaDesignu(fmt.Sprintf(
				"kolejność stron wskazuje numer %d, którego publikacja nie ma", numer))
		}
		if uzyte[numer] {
			return nil, bladWskazaniaDesignu(fmt.Sprintf(
				"kolejność stron wskazuje numer %d dwa razy — ta sama strona dwukrotnie w wydaniu "+
					"jest rozstrzygnięciem, którego rdzeń nie zgadnie; wskaż ją raz albo zduplikuj "+
					"stronę w szablonie", numer))
		}
		uzyte[numer] = true
		uszeregowane = append(uszeregowane, strona)
	}
	// Strona pominięta w kolejności NIE wchodzi do wydania i to jest wybór
	// Operatora — ale liczba stron w odpowiedzi (`pageCount`) mówi wtedy prawdę
	// o pliku, więc pominięcie nie przechodzi w ciszy.
	return uszeregowane, nil
}

// rodzajWydaniaDrukarskiegoDesignu nazywa rodzaj zasobu wydania. PDF i EPS są
// dokumentami, TIFF obrazem — rodzaj rozstrzyga o tym, czym Assets Panel go
// pokaże.
func rodzajWydaniaDrukarskiegoDesignu(format string) shared.DesignAssetKind {
	if format == "tiff" {
		return shared.DesignAssetKindImage
	}
	return shared.DesignAssetKindDocument
}

// pierwszyBladZastrzezenDesignu oddaje treść pierwszej wady o wadze błędu —
// wchodzi do odmowy, żeby Operator nie musiał puszczać drugiej komendy, aby się
// dowiedzieć, co jest nie tak.
func pierwszyBladZastrzezenDesignu(zastrzezenia []shared.DesignPreflightIssue) string {
	for _, zastrzezenie := range zastrzezenia {
		if zastrzezenie.Severity == shared.DesignPreflightSeverityBlad {
			return zastrzezenie.Message
		}
	}
	return "brak"
}

// zlozWydanieDrukarskieDesignu składa bajty wydania wraz z typem treści, nazwą
// pliku i oknem materiału źródłowego.
func (a *adapterDesignu) zlozWydanieDrukarskieDesignu(ctx context.Context,
	z shared.DesignPrintExportRequest, format string,
	profil shared.DesignPrintProfile) ([]byte, string, string, string, int, error) {

	if z.AssetId != nil && strings.TrimSpace(*z.AssetId) != "" {
		obraz, wiersz, err := a.obrazZasobuPoKodzieDesignu(ctx, "design.print.export", *z.AssetId)
		if err != nil {
			return nil, "", "", "", 0, err
		}
		bajty, typTresci, err := zakodujWydanieDrukarskieDesignu([]image.Image{obraz}, format)
		if err != nil {
			return nil, "", "", "", 0, bladWydaniaDesignu(err.Error())
		}
		nazwa := oczyscNazwePlikuDesignu(nazwaZasobuDesignu(wiersz)) + "." +
			rozszerzenieWydaniaDesignu(format)
		return bajty, typTresci, nazwa, wiersz.Okno, 1, nil
	}

	if z.BoardId == nil || strings.TrimSpace(*z.BoardId) == "" {
		return nil, "", "", "", 0, bladWskazaniaDesignu(
			"komenda design.print.export bez wskazania kompozycji, zasobu ani szablonu publikacji")
	}
	kompozycja, err := a.repozytorium.Kompozycja(ctx, strings.TrimSpace(*z.BoardId))
	if err != nil {
		return nil, "", "", "", 0, bladNieznanejKompozycjiDesignu(*z.BoardId, err)
	}
	strony, err := a.stronyWydaniaDrukarskiegoDesignu(ctx, kompozycja, z.FrameId, profil)
	if err != nil {
		return nil, "", "", "", 0, err
	}
	bajty, typTresci, err := zakodujWydanieDrukarskieDesignu(strony, format)
	if err != nil {
		return nil, "", "", "", 0, bladWydaniaDesignu(err.Error())
	}
	nazwa := oczyscNazwePlikuDesignu(nazwaKompozycjiDoPlikuDesignu(kompozycja)) + "." +
		rozszerzenieWydaniaDesignu(format)
	return bajty, typTresci, nazwa, kompozycja.Okno, len(strony), nil
}

// stronyWydaniaDrukarskiegoDesignu wyrysowuje strony wydania w rozdzielczości
// profilu.
//
// Kompozycja daje jedną stronę. Wielostronicowość wchodzi przez szablon materiału
// (`adapter_modul_design_szablony_materialu.go`) i tam ma swoją drogę — tutaj
// kompozycja jest arkuszem.
func (a *adapterDesignu) stronyWydaniaDrukarskiegoDesignu(ctx context.Context,
	kompozycja dane.KompozycjaDesignu, ramka *string,
	profil shared.DesignPrintProfile) ([]image.Image, error) {

	warstwy, err := a.repozytorium.Warstwy(ctx, kompozycja.ID)
	if err != nil {
		return nil, bladDesignu(err)
	}
	kafle, pominietych, err := a.kafleWyrysuDesignu(ctx, warstwy)
	if err != nil {
		return nil, err
	}
	if len(kafle) == 0 {
		return nil, bladWskazaniaDesignu(fmt.Sprintf(
			"kompozycja %s nie ma ani jednej warstwy z bajtami (warstw pominiętych: %d) — "+
				"wydanie do druku byłoby pustym arkuszem", kompozycja.Kod, pominietych))
	}

	obszar := obszarWyrysuDesignu(kafle, nil)
	if ramka != nil && strings.TrimSpace(*ramka) != "" {
		wiersz, err := a.repozytorium.RamkaDesignuPoKodzie(ctx, strings.TrimSpace(*ramka))
		if err != nil {
			return nil, bladNieznanejRamkiDesignu(*ramka, err)
		}
		if wiersz.KompozycjaID != kompozycja.ID {
			return nil, bladWskazaniaDesignu(fmt.Sprintf(
				"ramka %s leży na innej kompozycji niż wydawana", wiersz.Kod))
		}
		x, y := 0.0, 0.0
		if wiersz.X != nil {
			x = *wiersz.X
		}
		if wiersz.Y != nil {
			y = *wiersz.Y
		}
		obszar = shared.DesignBoardRegion{
			X: x, Y: y, Width: wiersz.Szerokosc, Height: wiersz.Wysokosc,
		}
	}
	// Spad rozszerza kadr: materiał ze spadem ma wychodzić poza krawędź, a
	// wydanie musi ten nadmiar unieść.
	spad := spadProfilu(profil)
	if spad > 0 {
		obszar = shared.DesignBoardRegion{
			X: obszar.X - spad, Y: obszar.Y - spad,
			Width: obszar.Width + 2*spad, Height: obszar.Height + 2*spad,
		}
	}

	// Jednostki kompozycji to milimetry (nagłówek pliku), więc skala z milimetrów
	// na piksele jest rozdzielczością profilu.
	rozdzielczosc := rozdzielczoscProfilu(profil)
	skala := float64(rozdzielczosc) / milimetryNaCal
	if pikseli := obszar.Width * skala * obszar.Height * skala; pikseli > granicaPikseliWydaniaDesignu {
		return nil, bladWskazaniaDesignu(fmt.Sprintf(
			"materiał %.0f×%.0f mm w %d dpi to %.0f pikseli płótna, a granica wydania to %d — "+
				"zmniejsz rozdzielczość profilu albo podziel materiał komendą design.largeformat.tile",
			obszar.Width, obszar.Height, rozdzielczosc, pikseli, granicaPikseliWydaniaDesignu))
	}
	return []image.Image{zlozWyrysRastrowyDesignu(kafle, obszar, skala)}, nil
}

// zakodujWydanieDrukarskieDesignu składa bajty wydania z gotowych stron.
func zakodujWydanieDrukarskieDesignu(strony []image.Image,
	format string) ([]byte, string, error) {

	if len(strony) == 0 {
		return nil, "", fmt.Errorf("wydanie bez ani jednej strony nie ma czego nieść")
	}
	switch format {
	case "pdf":
		return zakodujDokumentWielostronicowyDesignu(strony)
	case "tiff":
		// TIFF niesie jedną stronę: format wielostronicowy TIFF istnieje, ale
		// koder biblioteki go nie zapisuje, a wydanie pierwszej strony pod nazwą
		// całości byłoby cichą utratą reszty.
		if len(strony) > 1 {
			return nil, "", fmt.Errorf(
				"wydanie ma %d stron, a koder tiff biblioteki wkompilowanej zapisuje jedną — "+
					"wydaj publikację jako pdf albo wskaż jedną stronę", len(strony))
		}
		bajty, err := zakodujTiffDesignu(strony[0])
		if err != nil {
			return nil, "", err
		}
		return bajty, "image/tiff", nil
	case "eps":
		if len(strony) > 1 {
			return nil, "", fmt.Errorf(
				"wydanie ma %d stron, a EPS jest formatem jednostronicowym — wydaj publikację "+
					"jako pdf", len(strony))
		}
		bajty, err := zakodujEpsDesignu(strony[0])
		if err != nil {
			return nil, "", err
		}
		return bajty, "application/postscript", nil
	}
	return nil, "", fmt.Errorf("formatu %q rdzeń nie wydaje do druku", format)
}

// PodzielMaterialWielkoformatowy dzieli materiał na kafle — obsługuje
// `design.largeformat.tile`.
//
// Kafle są ZASOBAMI w magazynie, nie zapowiedzią: każdy niesie własne bajty,
// bo drukarnia wielkoformatowa dostaje pliki, nie wykaz prostokątów. Zakładka na
// sklejenie wchodzi w wymiar kafla, więc sąsiednie kafle mają wspólny pas obrazu.
func (a *adapterDesignu) PodzielMaterialWielkoformatowy(ctx context.Context,
	z shared.DesignLargeformatTileRequest) (shared.DesignLargeformatTileResponse, error) {

	if z.TileWidthMm <= 0 || z.TileHeightMm <= 0 {
		return shared.DesignLargeformatTileResponse{}, bladWskazaniaDesignu(fmt.Sprintf(
			"komenda design.largeformat.tile z kaflem %v×%v mm: kafel o niedodatnim boku nie istnieje",
			z.TileWidthMm, z.TileHeightMm))
	}
	zakladka := 0.0
	if z.OverlapMm != nil {
		if *z.OverlapMm < 0 {
			return shared.DesignLargeformatTileResponse{}, bladWskazaniaDesignu(fmt.Sprintf(
				"komenda design.largeformat.tile z zakładką %v mm: zakładka ujemna zostawiłaby "+
					"między kaflami dziurę", *z.OverlapMm))
		}
		zakladka = *z.OverlapMm
	}
	if zakladka*2 >= z.TileWidthMm || zakladka*2 >= z.TileHeightMm {
		return shared.DesignLargeformatTileResponse{}, bladWskazaniaDesignu(fmt.Sprintf(
			"zakładka %v mm nie mieści się w kaflu %v×%v mm — kafel byłby w całości zakładką",
			zakladka, z.TileWidthMm, z.TileHeightMm))
	}

	obraz, zrodlo, err := a.obrazZasobuPoKodzieDesignu(ctx, "design.largeformat.tile", z.AssetId)
	if err != nil {
		return shared.DesignLargeformatTileResponse{}, err
	}
	granice := obraz.Bounds()

	// Rozmiar docelowy całości: wskazany żądaniem albo wzięty z pikseli materiału
	// przy rozdzielczości drukarskiej. Bez jednego i drugiego nie da się
	// powiedzieć, ile kafli wyjdzie.
	docelowaSzerokosc := float64(granice.Dx()) / float64(rozdzielczoscDrukuDomyslna) * milimetryNaCal
	docelowaWysokosc := float64(granice.Dy()) / float64(rozdzielczoscDrukuDomyslna) * milimetryNaCal
	if z.TargetWidthMm != nil && *z.TargetWidthMm > 0 {
		docelowaSzerokosc = *z.TargetWidthMm
	}
	if z.TargetHeightMm != nil && *z.TargetHeightMm > 0 {
		docelowaWysokosc = *z.TargetHeightMm
	}

	// Krok jest kaflem pomniejszonym o zakładkę: to on rozstrzyga o liczbie kafli.
	krokX := z.TileWidthMm - zakladka
	krokY := z.TileHeightMm - zakladka
	kolumn := int(math.Ceil(docelowaSzerokosc / krokX))
	wierszy := int(math.Ceil(docelowaWysokosc / krokY))
	if kolumn < 1 {
		kolumn = 1
	}
	if wierszy < 1 {
		wierszy = 1
	}
	if kolumn*wierszy > granicaKafliDesignu {
		return shared.DesignLargeformatTileResponse{}, bladWskazaniaDesignu(fmt.Sprintf(
			"podział %.0f×%.0f mm na kafle %v×%v mm daje %d kafli, a granica to %d — każdy kafel "+
				"jest osobnym zasobem w magazynie; zwiększ kafel albo zmniejsz materiał",
			docelowaSzerokosc, docelowaWysokosc, z.TileWidthMm, z.TileHeightMm,
			kolumn*wierszy, granicaKafliDesignu))
	}

	// Rozdzielczość skuteczna: piksele materiału rozłożone na rozmiar docelowy.
	// To POMIAR, a nie życzenie — baner z obrazka 800 px ma w sześciu metrach
	// trzy dpi i Operator ma to wiedzieć.
	skuteczna := int(float64(granice.Dx()) / docelowaSzerokosc * milimetryNaCal)
	// Piksele na milimetr materiału źródłowego — tym przelicza się granice kafla
	// na wycinek obrazu.
	pikseliNaMmX := float64(granice.Dx()) / docelowaSzerokosc
	pikseliNaMmY := float64(granice.Dy()) / docelowaWysokosc

	okno := oknoWytworu(z.WindowId, zrodlo.Okno)
	znaczniki := z.Markers != nil && *z.Markers
	kafle := make([]shared.DesignTile, 0, kolumn*wierszy)
	for wiersz := 0; wiersz < wierszy; wiersz++ {
		for kolumna := 0; kolumna < kolumn; kolumna++ {
			lewaMm := float64(kolumna) * krokX
			goraMm := float64(wiersz) * krokY
			szerokoscMm := mniejszaDesignu(z.TileWidthMm, docelowaSzerokosc-lewaMm)
			wysokoscMm := mniejszaDesignu(z.TileHeightMm, docelowaWysokosc-goraMm)
			if szerokoscMm <= 0 || wysokoscMm <= 0 {
				continue
			}
			wycinek := image.Rect(
				granice.Min.X+int(lewaMm*pikseliNaMmX),
				granice.Min.Y+int(goraMm*pikseliNaMmY),
				granice.Min.X+int((lewaMm+szerokoscMm)*pikseliNaMmX),
				granice.Min.Y+int((goraMm+wysokoscMm)*pikseliNaMmY),
			).Intersect(granice)
			if wycinek.Empty() {
				continue
			}
			plotno := image.NewRGBA(image.Rect(0, 0, wycinek.Dx(), wycinek.Dy()))
			draw.Draw(plotno, plotno.Bounds(), obraz, wycinek.Min, draw.Src)
			if znaczniki {
				narysujZnacznikiSklejeniaDesignu(plotno, zakladka*pikseliNaMmX,
					zakladka*pikseliNaMmY, kolumna > 0, wiersz > 0,
					kolumna < kolumn-1, wiersz < wierszy-1)
			}
			bajty, _, err := zakodujObrazDesignu(plotno, "png", nil)
			if err != nil {
				return shared.DesignLargeformatTileResponse{}, bladWydaniaDesignu(err.Error())
			}
			nazwa := fmt.Sprintf("%s — kafel %d-%d", nazwaZasobuDesignu(zrodlo),
				wiersz+1, kolumna+1)
			zasob, err := a.zalozZasobZBajtowDesignu(ctx, okno, nazwa,
				shared.DesignAssetKindImage, "png", bajty)
			if err != nil {
				return shared.DesignLargeformatTileResponse{}, err
			}
			kod := zasob.Kod
			kafle = append(kafle, shared.DesignTile{
				Row: wiersz + 1, Column: kolumna + 1,
				WidthMm: szerokoscMm, HeightMm: wysokoscMm, AssetId: &kod,
			})
		}
	}
	if len(kafle) == 0 {
		return shared.DesignLargeformatTileResponse{}, bladWydaniaDesignu(
			"podział nie dał ani jednego kafla")
	}
	return shared.DesignLargeformatTileResponse{
		Tiles: kafle, Rows: wierszy, Columns: kolumn, EffectiveDpi: &skuteczna,
	}, nil
}

// narysujZnacznikiSklejeniaDesignu nakłada na kafel linie pokazujące, którędy
// biegnie zakładka.
//
// Znaczniki są rysowane na obrazie kafla, nie dokładane osobnym plikiem: kafel
// idzie do drukarki jako jeden plik i to na nim monter musi widzieć, gdzie
// nakłada sąsiada.
func narysujZnacznikiSklejeniaDesignu(plotno *image.RGBA, zakladkaX, zakladkaY float64,
	lewa, gora, prawa, dol bool) {

	granice := plotno.Bounds()
	kreska := image.NewUniform(barwaZnacznikaSklejeniaDesignu())
	if lewa && zakladkaX >= 1 {
		draw.Draw(plotno, image.Rect(int(zakladkaX), granice.Min.Y, int(zakladkaX)+1, granice.Max.Y),
			kreska, image.Point{}, draw.Src)
	}
	if prawa && zakladkaX >= 1 {
		x := granice.Max.X - int(zakladkaX)
		draw.Draw(plotno, image.Rect(x-1, granice.Min.Y, x, granice.Max.Y),
			kreska, image.Point{}, draw.Src)
	}
	if gora && zakladkaY >= 1 {
		draw.Draw(plotno, image.Rect(granice.Min.X, int(zakladkaY), granice.Max.X, int(zakladkaY)+1),
			kreska, image.Point{}, draw.Src)
	}
	if dol && zakladkaY >= 1 {
		y := granice.Max.Y - int(zakladkaY)
		draw.Draw(plotno, image.Rect(granice.Min.X, y-1, granice.Max.X, y),
			kreska, image.Point{}, draw.Src)
	}
}

// zakodujDokumentWielostronicowyDesignu składa dokument PDF z wielu stron
// biblioteką `pdfcpu` — tą samą, którą pracuje warsztat dokumentu modułu Studio.
//
// Strony jadą jako strumienie PNG w pamięci, w kolejności wykazu. Plik pośredni
// byłby trzecim miejscem, w którym ta sama treść żyje.
func zakodujDokumentWielostronicowyDesignu(strony []image.Image) ([]byte, string, error) {
	zrodla := make([]bytes.Reader, 0, len(strony))
	for _, strona := range strony {
		bajty, _, err := zakodujObrazDesignu(strona, "png", nil)
		if err != nil {
			return nil, "", err
		}
		zrodla = append(zrodla, *bytes.NewReader(bajty))
	}
	czytelniki := make([]io.Reader, 0, len(zrodla))
	for numer := range zrodla {
		czytelniki = append(czytelniki, &zrodla[numer])
	}
	var dokument bytes.Buffer
	if err := api.ImportImages(nil, &dokument, czytelniki, nil, nastawyPdf()); err != nil {
		return nil, "", fmt.Errorf("nie można złożyć dokumentu wydania: %w", err)
	}
	return dokument.Bytes(), "application/pdf", nil
}

// bladNieznanegoProfiluDrukuDesignu nazywa profil, którego rdzeń nie zna.
func bladNieznanegoProfiluDrukuDesignu(kod string, err error) error {
	if czyBrakZasobuDesignu(err) {
		return bladNieznanegoBytuDesignu("profilu druku " + kod + " nie ma w tym rdzeniu")
	}
	return bladDesignu(err)
}
