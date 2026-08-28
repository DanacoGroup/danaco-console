// Plik niesie wspólne zaplecze dwóch narzędzi archiwum — archive.pack i archive.unpack: wołanie 7z przez pakiet zewnetrzne, rozstrzygnięcie katalogu roboczego okna, odczyt wskazanych zasobów i odłożenie wyniku w magazynie.
package core

import (
	"context"
	"errors"
	"os"
	"strings"
	"time"

	"danacoconsole/server/internal/dane"
	"danacoconsole/server/internal/konfig"
	"danacoconsole/server/internal/session"
	"danacoconsole/server/internal/zewnetrzne"
	"danacoconsole/shared"
)

const (
	// granicaNarzedziArchiwum ogranicza jedno uruchomienie 7z. Dziesięć minut starczy na spakowanie katalogu roboczego, a zarazem kończy pracę programu karmionego archiwum uszkodzonym, które potrafi mielić bez końca.
	granicaNarzedziArchiwum = 600 * time.Second

	// oknoZasobowNarzedziArchiwum jest półką na zasoby powstałe z pakowania i rozpakowania. Kolumna zasob_design.okno jest NOT NULL, a żądania archive.* nie niosą okna, bo model woła je z rozmowy, nie z Assets Panelu.
	oknoZasobowNarzedziArchiwum = "narzedzia-archiwum"

	// granicaRozpakowaniaBajty ogranicza rozmiar treści po rozpakowaniu i broni przed bombą dekompresyjną, czyli archiwum kilkukilobajtowym rozwijającym się do gigabajtów. Dwa gibibajty muszą zmieścić się na nośniku dwukrotnie.
	granicaRozpakowaniaBajty int64 = 2 << 30

	// granicaRozpakowaniaPozycji jest drugą granicą tej samej obrony: liczbą pozycji, bo bomba nie musi być wielka w bajtach, a wyczerpuje i-węzły nośnika, których granica rozmiaru by nie dostrzegła.
	granicaRozpakowaniaPozycji = 100000
)

// narzedzie7z opisuje binarium pakujące. Program `7z` dostarcza w tej
// dystrybucji pakiet `7zip`; nazwa `p7zip-full` kierowałaby podpowiedź
// instalacyjną do pakietu, którego nie ma.
func narzedzie7z() zewnetrzne.Narzedzie {
	return zewnetrzne.Narzedzie{
		Nazwa:   "7-Zip",
		Program: "7z",
		Pakiet:  "7zip",
	}
}

// adapterNarzedziArchiwum wypełnia port NarzedziaArchiwum. Zależności są trzy: repozytorium Designu, magazyn bajtów oraz uruchamiacz wraz z dwoma źródłami izolacji. Adapter nie pamięta niczego między wywołaniami.
type adapterNarzedziArchiwum struct {
	repozytorium  dane.RepozytoriumDesignu
	magazyn       *magazynTresciBiblioteki
	katalogDanych string
	uruchamiacz   session.Uruchamiacz
	rozstrzygacz  *konfig.Rozstrzygacz
	katalog       *KatalogRoboczy
}

// nowyAdapterNarzedziArchiwum wiąże port z repozytorium modułu Design
// i magazynem jego zasobów. Uruchamiacz oraz izolacja wchodzą osobno
// (`ZArsenalem`), bo montaż zna je dopiero po złożeniu warstwy kanału.
func nowyAdapterNarzedziArchiwum(repozytorium dane.RepozytoriumDesignu,
	katalogDanych string) *adapterNarzedziArchiwum {

	return &adapterNarzedziArchiwum{
		repozytorium:  repozytorium,
		magazyn:       magazynZasobowDesignu(katalogDanych),
		katalogDanych: strings.TrimSpace(katalogDanych),
	}
}

// ZArsenalem wpina drogę do binarium: uruchamiacz procesów i te same dwa źródła izolacji, którymi pracują Terminal, Developer, silniki mowy oraz arsenały obrazu i mediów.
func (a *adapterNarzedziArchiwum) ZArsenalem(uruchamiacz session.Uruchamiacz,
	rozstrzygacz *konfig.Rozstrzygacz, katalog *KatalogRoboczy) *adapterNarzedziArchiwum {

	a.uruchamiacz = uruchamiacz
	a.rozstrzygacz = rozstrzygacz
	a.katalog = katalog
	return a
}

// zasiegNarzedzi składa trójkę okno–zasady–obszar dla zasięgu platformy: żądanie archive.* niesie samo archiwum, a nie okno rozmowy, więc adresem jest najszerszy z ośmiu poziomów zasięgu.
func (a *adapterNarzedziArchiwum) zasiegNarzedzi() (session.Okno, session.Zasady, session.Obszar) {
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
	return okno, zasady, obszar
}

// katalogRoboczyOkna oddaje katalog, względem którego rozstrzygane są wszystkie ścieżki tej rodziny. Katalog roboczy okna jest jedynym miejscem, w którym model ma prawo pisać.
func (a *adapterNarzedziArchiwum) katalogRoboczyOkna() (string, error) {
	_, _, obszar := a.zasiegNarzedzi()
	sciezka := strings.TrimSpace(obszar.KatalogRoboczy)
	if sciezka == "" {
		return "", bladZapleczaArchiwum("rdzeń nie ma ustalonego katalogu roboczego okna — " +
			"narzędzia archiwum nie mają gdzie pisać ani względem czego rozstrzygać ścieżek")
	}
	return sciezka, nil
}

// wolaj7z przeprowadza jedno uruchomienie `7z` i oddaje to, co program wypisał
// na wyjście (spis archiwum przy `l`, sprawozdanie przy `a` i `x`).
func (a *adapterNarzedziArchiwum) wolaj7z(ctx context.Context,
	argumenty []string, katalog string) (string, error) {

	if a.uruchamiacz == nil {
		return "", bladArchiwumNiedostepnego("rdzeń nie ma uruchamiacza procesów — " +
			"narzędzia archiwum nie mają czym wystartować; " +
			"naprawa: podpiąć warstwę kanału (injection) przy składaniu rdzenia")
	}
	okno, zasady, obszar := a.zasiegNarzedzi()
	if strings.TrimSpace(katalog) == "" {
		katalog = strings.TrimSpace(obszar.KatalogRoboczy)
	}

	wynik, err := zewnetrzne.Wolaj(ctx, a.uruchamiacz, okno, zasady, obszar,
		narzedzie7z(), argumenty, katalog, granicaNarzedziArchiwum)
	if err != nil {
		return string(wynik.Wyjscie), bladArsenaluArchiwum(err)
	}
	return string(wynik.Wyjscie), nil
}

// zasobArchiwum odczytuje wiersz zasobu i sprawdza, że jego bajty leżą tam,
// gdzie wskazuje `uri`. Wiersz bez treści na dysku jest odmową, nie pustym
// archiwum.
func (a *adapterNarzedziArchiwum) zasobArchiwum(ctx context.Context,
	kod string) (dane.ZasobDesignu, string, error) {

	if a.repozytorium == nil {
		return dane.ZasobDesignu{}, "", bladZapleczaArchiwum(
			"rdzeń nie ma magazynu zasobów — nie ma gdzie szukać wskazanego zasobu")
	}
	kod = strings.TrimSpace(kod)
	if kod == "" {
		return dane.ZasobDesignu{}, "", bladWskazaniaArchiwum("wskazanie zasobu jest puste")
	}
	zasob, err := a.repozytorium.Zasob(ctx, kod)
	if err != nil {
		if errors.Is(err, dane.ErrBrakWiersza) {
			return dane.ZasobDesignu{}, "", bladWskazaniaArchiwum("nie ma zasobu o identyfikatorze " + kod)
		}
		return dane.ZasobDesignu{}, "", bladZapleczaArchiwum("nie można odczytać zasobu " + kod + ": " + err.Error())
	}
	if zasob.URI == nil || strings.TrimSpace(*zasob.URI) == "" {
		return dane.ZasobDesignu{}, "", bladWskazaniaArchiwum("zasób " + kod +
			" nie ma treści (puste uri) — nie ma czego pakować")
	}
	if _, err := os.Stat(*zasob.URI); err != nil {
		return dane.ZasobDesignu{}, "", bladZapleczaArchiwum("treść zasobu " + kod +
			" nie leży pod zapisanym odwołaniem: " + err.Error())
	}
	return zasob, *zasob.URI, nil
}

// katalogRoboczyTymczasowy zakłada katalog na pracę pośrednią — rusztowanie pakowania, kwarantannę rozpakowania. Podstawa wskazana trzyma pracę na tym samym nośniku, co jej cel.
func katalogRoboczyTymczasowy(podstawa, przedrostek string) (string, error) {
	podstawa = strings.TrimSpace(podstawa)
	if podstawa != "" {
		if err := os.MkdirAll(podstawa, 0o755); err != nil {
			return "", bladZapleczaArchiwum("nie można założyć katalogu " + podstawa + ": " + err.Error())
		}
	}
	katalog, err := os.MkdirTemp(podstawa, przedrostek)
	if err != nil {
		return "", bladZapleczaArchiwum("nie można założyć katalogu roboczego: " + err.Error())
	}
	return katalog, nil
}

// odlozArchiwumJakoZasob utrwala gotowy plik archiwum w magazynie i zakłada jego wiersz, tą samą drogą, którą idzie design.asset.upload. Kolejność jest zamierzona: najpierw bajty, potem wiersz.
func (a *adapterNarzedziArchiwum) odlozArchiwumJakoZasob(ctx context.Context,
	sciezka, nazwa string, format formatArchiwum) (shared.DesignAsset, int64, error) {

	if a.magazyn == nil {
		return shared.DesignAsset{}, 0, bladZapleczaArchiwum(
			"magazyn treści nie jest wpięty — nie ma gdzie odłożyć bajtów archiwum")
	}
	if a.repozytorium == nil {
		return shared.DesignAsset{}, 0, bladZapleczaArchiwum(
			"rdzeń nie ma rejestru zasobów — nie ma gdzie założyć wiersza archiwum")
	}

	odwolanie, rozmiar, _, err := a.magazyn.ZapiszZePliku(sciezka)
	if err != nil {
		return shared.DesignAsset{}, 0, bladZapleczaArchiwum(
			"nie można utrwalić treści archiwum: " + err.Error())
	}
	if rozmiar == 0 {
		// Program zakończył się powodzeniem i zostawił plik pusty; puste
		// archiwum nie jest wynikiem pakowania.
		return shared.DesignAsset{}, 0, bladPrzetwarzaniaArchiwum(
			"7-Zip zakończył się powodzeniem, ale archiwum wyszło puste")
	}

	zapisany, err := a.repozytorium.ZapiszZasob(ctx, dane.ZasobDesignu{
		Kod:    nowyIdentyfikator(przedrostekZasobuDesign),
		Okno:   oknoZasobowNarzedziArchiwum,
		Nazwa:  wskaznikTekstu(nazwa),
		Rodzaj: shared.DesignAssetKindImage,
		Format: wskaznikTekstu(string(format)),
		URI:    &odwolanie,
	})
	if err != nil {
		return shared.DesignAsset{}, 0, bladZapleczaArchiwum(
			"nie można założyć wiersza zasobu archiwum: " + err.Error())
	}
	return zasobKontraktu(zapisany, nil), rozmiar, nil
}

// sprzatnij usuwa katalog pracy pośredniej. Niepowodzenie sprzątania nie jest
// odmową komendy: praca się udała, a zostawiony katalog tymczasowy jest
// śmieciem, nie utratą wyniku.
func sprzatnij(katalog string) {
	if strings.TrimSpace(katalog) != "" {
		_ = os.RemoveAll(katalog)
	}
}
