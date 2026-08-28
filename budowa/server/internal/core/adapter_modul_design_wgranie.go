// Odpowiedzialność pliku: wniesienie zasobu do Assets Panel
// (`design.asset.upload`) — magazyn bajtów zasobu, rozstrzygnięcie treści
// żądania i rozpoznanie formatu oraz wymiarów z nagłówka pliku.
package core

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"image"
	// Rejestracja dekoderów: bez importów `image.DecodeConfig` nie rozpozna formatu.
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"os"
	"path/filepath"
	"strings"

	"danacoconsole/server/internal/dane"
	"danacoconsole/server/internal/protocol"
	"danacoconsole/shared"
)

const (
	// przedrostekZasobuDesign znakuje identyfikatory zewnętrzne zasobów —
	// obok `plansza-` i `warstwa-` z `adapter_modul_design_kompozycje.go`.
	przedrostekZasobuDesign = "zasob-"

	// przedrostekPromptuDesign znakuje identyfikatory zewnętrzne promptów
	// wydanych — tych, które wchodzą do `design.prompt.history.list` i które
	// zasób wskazuje polem `promptId`.
	przedrostekPromptuDesign = "prompt-"

	// Magazyn zasobów leży w katalogu danych rdzenia, obok magazynu biblioteki
	// i sejfu poświadczeń — przeżywa restart tak samo jak wiersz w bazie.
	podkatalogDesignu        = "design"
	podkatalogZasobowDesignu = "zasoby"
)

// magazynZasobowDesignu składa magazyn bajtów zasobów nad katalogiem danych
// rdzenia. Katalog nie musi istnieć — powstaje przy pierwszym zapisie.
func magazynZasobowDesignu(katalogDanych string) *magazynTresciBiblioteki {
	if strings.TrimSpace(katalogDanych) == "" {
		return nil
	}
	return &magazynTresciBiblioteki{
		katalog: filepath.Join(katalogDanych, podkatalogDesignu, podkatalogZasobowDesignu),
	}
}

// korzenZasobowDesignu jest korzeniem magazynu zasobów liczonym od katalogu
// danych — ta sama para podkatalogów, którą składa `magazynZasobowDesignu`,
// więc postać odwołania i miejsce zapisu nie mają jak się rozjechać.
const korzenZasobowDesignu = podkatalogDesignu + "/" + podkatalogZasobowDesignu

// odwolanieZasobuDesignu oddaje `uri` zasobu w postaci, którą wolno wypuścić
// z rdzenia: ścieżkę względną magazynu zasobów, nie ścieżkę na dysku Operatora.
func odwolanieZasobuDesignu(sciezka string) string {
	return odwolanieMagazynu(sciezka, korzenZasobowDesignu)
}

// WniesZasob utrwala bajty zasobu, zakłada jego wiersz i oddaje zasób odczytany
// z bazy wraz z etykietami — obsługuje `design.asset.upload`. Kolejność jest
// zamierzona: najpierw bajty, potem wiersz.
func (a *adapterDesignu) WniesZasob(ctx context.Context,
	z shared.DesignAssetUploadRequest) (shared.DesignAssetUploadResponse, error) {

	if z.WindowId == "" {
		return shared.DesignAssetUploadResponse{}, bladWskazaniaDesignu(
			"komenda design.asset.upload bez wskazania okna")
	}
	if strings.TrimSpace(string(z.Kind)) == "" {
		return shared.DesignAssetUploadResponse{}, bladWskazaniaDesignu(
			"komenda design.asset.upload bez rodzaju zasobu")
	}
	// Sprawdzenie idzie przed wniesieniem treści, bo odmowa po zapisie
	// zostawiałaby blob bez wiersza.
	if err := sprawdzRodzajZasobu("design.asset.upload", z.Kind); err != nil {
		return shared.DesignAssetUploadResponse{}, err
	}

	odwolanie, err := a.trescZasobu(z.ContentBase64, z.SourcePath)
	if err != nil {
		return shared.DesignAssetUploadResponse{}, err
	}

	// Format i wymiary rozpoznaje się z tego, co naprawdę leży w magazynie,
	// nie z żądania.
	format, szerokosc, wysokosc := rozpoznajObrazZasobu(odwolanie)
	// Zmierzone bije zadeklarowane, gdy oba są dostępne.
	format = pierwszyTekst(format, z.Format)
	szerokosc = pierwszaLiczba(szerokosc, z.Width)
	wysokosc = pierwszaLiczba(wysokosc, z.Height)

	zapisany, err := a.repozytorium.ZapiszZasob(ctx, dane.ZasobDesignu{
		Kod:       nowyIdentyfikator(przedrostekZasobuDesign),
		Okno:      z.WindowId,
		Nazwa:     z.Name,
		Rodzaj:    string(z.Kind),
		Format:    format,
		URI:       &odwolanie,
		Szerokosc: szerokosc,
		Wysokosc:  wysokosc,
	})
	if err != nil {
		return shared.DesignAssetUploadResponse{}, bladDesignu(err)
	}

	if len(z.Tags) > 0 {
		if err := a.repozytorium.UstawEtykietyZasobu(ctx, zapisany.ID, z.Tags); err != nil {
			return shared.DesignAssetUploadResponse{}, bladDesignu(err)
		}
	}
	etykiety, err := a.repozytorium.EtykietyZasobu(ctx, zapisany.ID)
	if err != nil {
		return shared.DesignAssetUploadResponse{}, bladDesignu(err)
	}
	return shared.DesignAssetUploadResponse{Asset: zasobKontraktu(zapisany, etykiety)}, nil
}

// trescZasobu utrwala treść żądania w magazynie i oddaje odwołanie do bajtów —
// ścieżkę bloba, która staje się `uri` zasobu. Ścieżka ma pierwszeństwo przed
// treścią, gdy przyszły obie.
func (a *adapterDesignu) trescZasobu(trescBase64, sciezkaZrodlowa *string) (string, error) {
	if a.magazyn == nil {
		return "", bladZapisuZasobuDesignu(
			"magazyn treści nie jest wpięty — nie ma gdzie odłożyć bajtów zasobu")
	}

	if !bezWartosci(sciezkaZrodlowa) {
		// Wciąga bajty, nie dowiązuje pliku.
		odwolanie, _, _, err := a.magazyn.ZapiszZePliku(*sciezkaZrodlowa)
		if err != nil {
			return "", bladZapisuZasobuDesignu("nie można wciągnąć treści spod ścieżki: " + err.Error())
		}
		return odwolanie, nil
	}

	if bezWartosci(trescBase64) {
		return "", bladWskazaniaDesignu(
			"komenda design.asset.upload bez treści zasobu: brakuje pola contentBase64 albo sourcePath — " +
				"wniesienie nie ma skąd wziąć bajtów; obraz wytworzony przez rdzeń zamawia się " +
				"komendą design.asset.generate, wskazując kanał obrazowy")
	}
	// Rozbiór base64 idzie tu, wyłącznie dla odmowy mówiącej o module Design.
	bajty, err := base64.StdEncoding.DecodeString(*trescBase64)
	if err != nil {
		return "", bladWskazaniaDesignu("treść zasobu nie jest poprawnym base64: " + err.Error())
	}
	if len(bajty) == 0 {
		return "", bladWskazaniaDesignu("treść zasobu jest pusta — zasób bez bajtów nie ma czego pokazać")
	}
	suma := sha256.Sum256(bajty)
	odwolanie, err := a.magazyn.Zapisz(bajty, hex.EncodeToString(suma[:]))
	if err != nil {
		return "", bladZapisuZasobuDesignu("nie można utrwalić treści zasobu: " + err.Error())
	}
	return odwolanie, nil
}

// rozpoznajObrazZasobu czyta nagłówek utrwalonego pliku i oddaje format wraz
// z wymiarami. Formatu nierozpoznanego nie zgaduje: trójka pustych wartości
// znaczy „plik nie jest obrazem, który rdzeń potrafi zmierzyć".
func rozpoznajObrazZasobu(sciezka string) (*string, *int, *int) {
	plik, err := os.Open(sciezka)
	if err != nil {
		return nil, nil, nil
	}
	defer plik.Close()

	opis, format, err := image.DecodeConfig(plik)
	if err != nil {
		return nil, nil, nil
	}
	szerokosc, wysokosc := opis.Width, opis.Height
	return &format, &szerokosc, &wysokosc
}

// pierwszyTekst oddaje wartość zmierzoną, a gdy jej nie ma — zadeklarowaną
// przez Operatora. Wartość pusta jest tu brakiem, nie tekstem: „format: " bez
// formatu byłby cechą, której nikt nie wskazał.
func pierwszyTekst(zmierzona, zadeklarowana *string) *string {
	if zmierzona != nil && strings.TrimSpace(*zmierzona) != "" {
		return zmierzona
	}
	if zadeklarowana != nil && strings.TrimSpace(*zadeklarowana) != "" {
		return zadeklarowana
	}
	return nil
}

// pierwszaLiczba działa jak `pierwszyTekst` dla wymiarów. Liczba niedodatnia
// jest odrzucana po obu stronach — obraz o szerokości zero nie istnieje, a
// zapisany w wierszu wyglądałby na zmierzony.
func pierwszaLiczba(zmierzona, zadeklarowana *int) *int {
	if zmierzona != nil && *zmierzona > 0 {
		return zmierzona
	}
	if zadeklarowana != nil && *zadeklarowana > 0 {
		return zadeklarowana
	}
	return nil
}

// bladZapisuZasobuDesignu nazywa niepowodzenie utrwalenia bajtów: nośnik
// pełny, brak praw do katalogu danych, ścieżka źródłowa nieczytelna, magazyn
// niewpięty. To awaria zapisu po stronie rdzenia, nie wina Operatora.
func bladZapisuZasobuDesignu(powod string) error {
	return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeInternalError, "moduł Design: "+powod))
}
