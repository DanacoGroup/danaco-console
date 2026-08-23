// Odpowiedzialność pliku: jedna prawda o tym, co dzieje się z wynikiem
// czynności arsenału — cztery rodziny narzędzi modelu (`image.*`, `media.*`,
// `document.*`, `archive.*`) odkładają bajty i zakładają wiersz zasobu tą samą
// drogą, nazywają rodzaj powstałej rzeczy tą samą tablicą i rozstrzygają okno
// wyniku tą samą regułą. Cztery kopie jednej reguły byłyby czterema okazjami do
// rozjazdu; tutaj reguła jest jedna:
//
//  1. Okno wyniku bierze się z nieobowiązkowego pola `windowId` żądania, a gdy
//     go nie ma — z okna zasobu źródłowego. Patrz `oknoWynikuArsenalu`.
//
//  2. Rodzaj zasobu rozstrzyga format wyniku, nie nazwa komendy. Warunek CHECK
//     kolumny rodzaju poszerza `migracja_113_rodzaje_zasobow_arsenalu.sql`, więc
//     film nie musi jechać jako `image`, a PDF jako `vector`. Patrz
//     `rodzajZasobuArsenalu`.
//
//  3. Odłożenie bajtów i założenie wiersza idą jedną sekwencją: zapisz do
//     magazynu, złóż `dane.ZasobDesignu`, zapisz wiersz, przełóż na kontrakt.
//     Patrz `odlozWynikArsenalu`.
//
//  4. Magazynem bajtów jest magazyn zasobów Designu — ten sam, którym jedzie
//     `design.asset.upload`. Wynikiem każdej z tych czynności jest `DesignAsset`
//     i Assets Panel ma prawo go otworzyć, a dwa składy bajtów byłyby dwiema
//     prawdami o tym, gdzie rdzeń trzyma treść poza bazą.
//
// Odmowy zostają przy rodzinach. Treść odmowy nazywa brak właściwy rodzinie
// („brakuje pola assetId albo sourcePath — nie ma czego zmierzyć" mówi co innego
// niż „…nie ma czego spakować"), a wspólna formułka zamieniłaby cztery zdania
// mówiące na jedno milczące. Dlatego wspólne funkcje przyjmują opakowywacz
// odmowy rodziny i nie znają żadnego kodu błędu z własnej głowy.
package core

import (
	"context"
	"strings"
	"time"

	"danacoconsole/server/internal/dane"
	"danacoconsole/shared"
)

// odmowaRodziny jest opakowywaczem odmowy wnoszonym przez rodzinę wołającą.
// Wspólny kod arsenału nie zna kodów błędu ani nagłówków komunikatów — zna
// wyłącznie powód i oddaje jego ubranie rodzinie, która wie, czyim brakiem to
// jest i jak on brzmi po polsku.
type odmowaRodziny func(powod string) error

// oknoWynikuArsenalu rozstrzyga, do którego okna należy wynik czynności —
// jedna reguła dla wszystkich czterech rodzin.
//
// Pusty wynik znaczy „bez wiersza", a nie „bez okna". Kolumna
// `zasob_design.okno` jest NOT NULL (`migracja_048_design.sql`) i `ZapiszZasob`
// odmawia zasobowi bez okna. Gdy nie wiadomo, do którego okna wynik należy, nie
// zakładamy wiersza wcale — i to jest zamierzone: bajty leżą w magazynie pod
// sumą kontrolną, model dostaje odwołanie, którym plik da się otworzyć,
// a Assets Panel nie dostaje kafelka o zmyślonej przynależności. Dokładnie to
// obiecuje opis pola `windowId` w kontrakcie: „Brak nie wstrzymuje czynności:
// bajty i tak trafiają do magazynu pod sumą kontrolną, ale zasób nie pojawi się
// w wykazie okna".
//
// Kolejność pytań jest zamierzona. Najpierw `windowId` z żądania — model zna
// własny zasięg i podaje okno, w którym Operator na wynik czeka. Potem okno
// zasobu źródłowego: gdy czynność przerabia zasób leżący już w jakimś oknie,
// wynik zostaje tam, gdzie Operator widzi materiał, z którego powstał. To nie
// jest zgadywanie — to okno odczytane z wiersza źródła, a nie wymyślone.
//
// Nazwa okna, którego nie ma w rejestrze okien, byłaby przynależnością
// zmyśloną: Operator nie ma jak takiego wykazu otworzyć, a wiersz i tak nie
// trafiłby do żadnego panelu. Półka nazwana wygląda porządniej niż brak wiersza,
// ale mówi nieprawdę, a brak wiersza mówi prawdę.
func oknoWynikuArsenalu(zadane *string, oknoZrodla string) string {
	if !bezWartosci(zadane) {
		return strings.TrimSpace(*zadane)
	}
	return strings.TrimSpace(oknoZrodla)
}

// rodzajeZasobuPoFormacie przekłada format wyniku na rodzaj zasobu kontraktu.
//
// Rozstrzyga format, a nie nazwa komendy, i to jest cała istota tej tablicy.
// `media.transcode` z czynnością `frame` daje obraz, a z `extractAudio` —
// dźwięk, mimo że komenda jest ta sama; `document.convert` do `png` daje obraz,
// a nie dokument. Rodzaj wyprowadzony z nazwy rodziny byłby etykietą wołacza,
// nie opisem tego, co naprawdę leży w magazynie.
var rodzajeZasobuPoFormacie = map[string]shared.DesignAssetKind{
	// Rastry — to, co ma piksele i wymiary.
	"png": shared.DesignAssetKindImage, "jpg": shared.DesignAssetKindImage,
	"jpeg": shared.DesignAssetKindImage, "webp": shared.DesignAssetKindImage,
	"gif": shared.DesignAssetKindImage, "bmp": shared.DesignAssetKindImage,
	"tif": shared.DesignAssetKindImage, "tiff": shared.DesignAssetKindImage,
	"avif": shared.DesignAssetKindImage, "heic": shared.DesignAssetKindImage,
	"ppm": shared.DesignAssetKindImage,

	// Grafika opisana krzywymi — osobno od rastra, bo `vector` istniał
	// w kontrakcie od początku i znaczy dokładnie to.
	"svg": shared.DesignAssetKindVector, "eps": shared.DesignAssetKindVector,
	"ai": shared.DesignAssetKindVector,

	// Film. Kontener bywa wspólny dla obrazu i dźwięku (`mkv`, `webm`) —
	// rozstrzyga tu format, o który poproszono narzędzie, a rodziny pytające
	// o strumienie (`media.inspect`) mają swoją własną, dokładniejszą drogę.
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

	// Dokumenty — to, co się czyta, niezależnie od tego, czy w środku są
	// krzywe (PDF), znaczniki (HTML) czy zwykły tekst.
	"pdf": shared.DesignAssetKindDocument, "docx": shared.DesignAssetKindDocument,
	"doc": shared.DesignAssetKindDocument, "odt": shared.DesignAssetKindDocument,
	"rtf": shared.DesignAssetKindDocument, "md": shared.DesignAssetKindDocument,
	"markdown": shared.DesignAssetKindDocument, "html": shared.DesignAssetKindDocument,
	"htm": shared.DesignAssetKindDocument, "txt": shared.DesignAssetKindDocument,
	"tex": shared.DesignAssetKindDocument, "epub": shared.DesignAssetKindDocument,
	"csv": shared.DesignAssetKindDocument, "xlsx": shared.DesignAssetKindDocument,
	"pptx": shared.DesignAssetKindDocument,

	// Archiwa — pojemniki. `tar.gz` i `tar.bz2` przychodzą tu już rozłożone
	// na ostatni człon przez `rodzajZasobuArsenalu`.
	"zip": shared.DesignAssetKindArchive, "7z": shared.DesignAssetKindArchive,
	"tar": shared.DesignAssetKindArchive, "gz": shared.DesignAssetKindArchive,
	"tgz": shared.DesignAssetKindArchive, "bz2": shared.DesignAssetKindArchive,
	"xz": shared.DesignAssetKindArchive, "rar": shared.DesignAssetKindArchive,
	"zst": shared.DesignAssetKindArchive,
}

// rodzajZasobuArsenalu nazywa rodzaj zasobu po formacie tego, co naprawdę
// powstało.
//
// Format nieznany spada na `document`, a nie na `image`, i to jest wybór, nie
// przeoczenie. Format spoza tablicy to prawie zawsze plik do otwarcia jakimś
// programem, a nie coś o pikselach i wymiarach; wpisanie `image` obiecywałoby
// Assets Panelowi podgląd, którego nie ma czym narysować. Poszerzenie tablicy
// jest tańsze niż odmowa, a odmowa byłaby tu nie na miejscu: bajty powstały
// i model ma prawo dostać do nich odwołanie.
func rodzajZasobuArsenalu(format string) shared.DesignAssetKind {
	format = strings.ToLower(strings.TrimSpace(format))
	format = strings.TrimPrefix(format, ".")
	if format == "" {
		return shared.DesignAssetKindDocument
	}
	// `tar.gz`, `tar.bz2`, `tar.xz` — liczy się ostatni człon, bo to on nazywa
	// rzecz leżącą w magazynie.
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
	// Utrwalenie idzie przed założeniem wiersza — patrz `odlozWynikArsenalu`.
	odwolanie string
	// okno pustym tekstem znaczy „bez wiersza" — patrz `oknoWynikuArsenalu`.
	okno string
	// nazwa jest podpisem kafelka w Assets Panelu.
	nazwa string
	// format nazywa rzecz i jednocześnie rozstrzyga o rodzaju zasobu.
	format string
	// Wymiary niesie wyłącznie to, co je ma i co zmierzono. Puste znaczy
	// „nie wiem" i tak ma zostać.
	szerokosc, wysokosc *int
}

// odlozWynikArsenalu zakłada wiersz zasobu dla wyniku czynności — albo świadomie
// go nie zakłada, gdy nie wiadomo, do którego okna wynik należy.
//
// Bajty są już utrwalone, zanim ta funkcja ruszy, i tak ma być: kolejność
// „najpierw treść, potem wiersz" jest wspólna całemu rdzeniowi (`design.asset.upload`),
// bo wiersz wskazujący odwołanie, za którym nic nie leży, jest w panelu
// kafelkiem bez zawartości.
//
// Zasób bez wiersza jest co do bajtów zasobem prawdziwym. Oddajemy wtedy
// `DesignAsset` złożony z ręki: `uri` wskazuje blob pod sumą sha256, a pusty
// `windowId` jest uczciwym oświadczeniem „ten wynik do żadnego okna nie
// należy". Model dostaje odwołanie, którym plik da się odczytać i oddać
// Operatorowi; Operator nie dostaje w panelu nic — bo o nic nie prosił.
func odlozWynikArsenalu(ctx context.Context, repozytorium dane.RepozytoriumDesignu,
	w wynikArsenalu, odmowa odmowaRodziny) (shared.DesignAsset, error) {

	rodzaj := rodzajZasobuArsenalu(w.format)
	nazwa := strings.TrimSpace(w.nazwa)
	format := strings.TrimSpace(w.format)

	if w.okno == "" || repozytorium == nil {
		// Brak okna jest rozstrzygnięciem wołającego; brak repozytorium jest
		// niepełnym montażem rdzenia. Skutek dla wyniku jest ten sam — zasób
		// bez wiersza — więc i droga jest ta sama, zamiast dwóch prawie
		// identycznych gałęzi.
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
