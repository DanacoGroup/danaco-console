// Plik niesie jedną prawdę o wyniku czynności arsenału: cztery rodziny
// narzędzi modelu odkładają bajty i zakładają wiersz zasobu tą samą drogą,
// nazywają rodzaj rzeczy tą samą tablicą i rozstrzygają okno wyniku tą samą
// regułą.
package core

import (
	"context"
	"strings"
	"time"

	"danacoconsole/server/internal/dane"
	"danacoconsole/shared"
)

// odmowaRodziny jest opakowywaczem odmowy wnoszonym przez rodzinę wołającą:
// wspólny kod arsenału nie zna kodów błędu ani nagłówków komunikatów, zna
// wyłącznie powód.
type odmowaRodziny func(powod string) error

// oknoWynikuArsenalu rozstrzyga, do którego okna należy wynik czynności,
// jedną regułą dla wszystkich czterech rodzin: najpierw pole żądania, potem
// okno zasobu źródłowego, a bez obu wiersz świadomie nie powstaje.
func oknoWynikuArsenalu(zadane *string, oknoZrodla string) string {
	if !bezWartosci(zadane) {
		return strings.TrimSpace(*zadane)
	}
	return strings.TrimSpace(oknoZrodla)
}

// rodzajeZasobuPoFormacie przekłada format wyniku na rodzaj zasobu kontraktu:
// rozstrzyga format, a nie nazwa komendy wołanej przez rodzinę narzędzi.
var rodzajeZasobuPoFormacie = map[string]shared.DesignAssetKind{
	// Rastry — to, co ma piksele i wymiary.
	"png": shared.DesignAssetKindImage, "jpg": shared.DesignAssetKindImage,
	"jpeg": shared.DesignAssetKindImage, "webp": shared.DesignAssetKindImage,
	"gif": shared.DesignAssetKindImage, "bmp": shared.DesignAssetKindImage,
	"tif": shared.DesignAssetKindImage, "tiff": shared.DesignAssetKindImage,
	"avif": shared.DesignAssetKindImage, "heic": shared.DesignAssetKindImage,
	"ppm": shared.DesignAssetKindImage,

	// Grafika opisana krzywymi, osobno od rastra, bo vector znaczy dokładnie to.
	"svg": shared.DesignAssetKindVector, "eps": shared.DesignAssetKindVector,
	"ai": shared.DesignAssetKindVector,

	// Film — kontener bywa wspólny dla obrazu i dźwięku, rozstrzyga format,
	// o który poproszono narzędzie.
	"mp4": shared.DesignAssetKindVideo, "mkv": shared.DesignAssetKindVideo,
	"webm": shared.DesignAssetKindVideo, "mov": shared.DesignAssetKindVideo,
	"avi": shared.DesignAssetKindVideo, "m4v": shared.DesignAssetKindVideo,
	"mpg": shared.DesignAssetKindVideo, "mpeg": shared.DesignAssetKindVideo,

	// Dźwięk.
	"mp3": shared.DesignAssetKindAudio, "wav": shared.DesignAssetKindAudio,
	"flac": shared.DesignAssetKindAudio, "ogg": shared.DesignAssetKindAudio,
	"oga": shared.DesignAssetKindAudio, "opus": shared.DesignAssetKindAudio,
	"aac": shared.DesignAssetKindAudio, "m4a": shared.DesignAssetKindAudio,
	"wma": shared.DesignAssetKindAudio,

	// Dokumenty — to, co się czyta: krzywe PDF, znaczniki HTML czy zwykły tekst.
	"pdf": shared.DesignAssetKindDocument, "docx": shared.DesignAssetKindDocument,
	"doc": shared.DesignAssetKindDocument, "odt": shared.DesignAssetKindDocument,
	"rtf": shared.DesignAssetKindDocument, "md": shared.DesignAssetKindDocument,
	"markdown": shared.DesignAssetKindDocument, "html": shared.DesignAssetKindDocument,
	"htm": shared.DesignAssetKindDocument, "txt": shared.DesignAssetKindDocument,
	"tex": shared.DesignAssetKindDocument, "epub": shared.DesignAssetKindDocument,
	"csv": shared.DesignAssetKindDocument, "xlsx": shared.DesignAssetKindDocument,
	"pptx": shared.DesignAssetKindDocument,

	// Archiwa — pojemniki; tar.gz i tar.bz2 przychodzą tu już rozłożone na
	// ostatni człon.
	"zip": shared.DesignAssetKindArchive, "7z": shared.DesignAssetKindArchive,
	"tar": shared.DesignAssetKindArchive, "gz": shared.DesignAssetKindArchive,
	"tgz": shared.DesignAssetKindArchive, "bz2": shared.DesignAssetKindArchive,
	"xz": shared.DesignAssetKindArchive, "rar": shared.DesignAssetKindArchive,
	"zst": shared.DesignAssetKindArchive,
}

// rodzajZasobuArsenalu nazywa rodzaj zasobu po formacie tego, co naprawdę
// powstało. Format nieznany spada na `document`, nie na `image`, bo prawie
// zawsze jest to plik do otwarcia programem, nie obraz z pikselami.
func rodzajZasobuArsenalu(format string) shared.DesignAssetKind {
	format = strings.ToLower(strings.TrimSpace(format))
	format = strings.TrimPrefix(format, ".")
	if format == "" {
		return shared.DesignAssetKindDocument
	}
	// tar.gz, tar.bz2, tar.xz — liczy się ostatni człon nazwy rzeczy.
	if i := strings.LastIndex(format, "."); i >= 0 {
		format = format[i+1:]
	}
	if rodzaj, znany := rodzajeZasobuPoFormacie[format]; znany {
		return rodzaj
	}
	return shared.DesignAssetKindDocument
}

// wynikArsenalu opisuje to, co powstało, w słowach wspólnych czterem rodzinom:
// gdzie leżą bajty, jak się rzecz nazywa, w jakim jest formacie i do którego
// okna należy.
type wynikArsenalu struct {
	// odwolanie wskazuje bajty już utrwalone w magazynie zasobów Designu.
	odwolanie string
	// okno pustym tekstem znaczy „bez wiersza".
	okno string
	// nazwa jest podpisem kafelka w Assets Panelu.
	nazwa string
	// format nazywa rzecz i jednocześnie rozstrzyga o rodzaju zasobu.
	format string
	// Wymiary niesie wyłącznie to, co je ma i co zmierzono. Puste znaczy
	// „nie wiem" i tak ma zostać.
	szerokosc, wysokosc *int
}

// odlozWynikArsenalu zakłada wiersz zasobu dla wyniku czynności, albo
// świadomie go nie zakłada, gdy nie wiadomo, do którego okna wynik należy.
// Bajty są utrwalone, zanim ta funkcja ruszy.
func odlozWynikArsenalu(ctx context.Context, repozytorium dane.RepozytoriumDesignu,
	w wynikArsenalu, odmowa odmowaRodziny) (shared.DesignAsset, error) {

	rodzaj := rodzajZasobuArsenalu(w.format)
	nazwa := strings.TrimSpace(w.nazwa)
	format := strings.TrimSpace(w.format)

	if w.okno == "" || repozytorium == nil {
		// Brak okna jest rozstrzygnięciem wołającego; brak repozytorium jest
		// niepełnym montażem.
		return shared.DesignAsset{
			Id:        nowyIdentyfikator(przedrostekZasobuDesign),
			WindowId:  "",
			Name:      wskaznikTekstu(nazwa),
			Kind:      rodzaj,
			Format:    wskaznikTekstu(format),
			Uri:       &w.odwolanie,
			Width:     w.szerokosc,
			Height:    w.wysokosc,
			CreatedAt: time.Now().UTC().UnixMilli(),
		}, nil
	}

	zapisany, err := repozytorium.ZapiszZasob(ctx, dane.ZasobDesignu{
		Kod:       nowyIdentyfikator(przedrostekZasobuDesign),
		Okno:      w.okno,
		Nazwa:     wskaznikTekstu(nazwa),
		Rodzaj:    string(rodzaj),
		Format:    wskaznikTekstu(format),
		URI:       &w.odwolanie,
		Szerokosc: w.szerokosc,
		Wysokosc:  w.wysokosc,
	})
	if err != nil {
		return shared.DesignAsset{}, odmowa("nie można założyć wiersza zasobu wynikowego: " + err.Error())
	}
	return zasobKontraktu(zapisany, nil), nil
}
