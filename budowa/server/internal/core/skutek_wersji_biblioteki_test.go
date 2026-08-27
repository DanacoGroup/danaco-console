// Sprawdza, czy wskaźnik wersji biblioteki wskazuje właściwy plik: treść czytana
// przez podgląd ma odpowiadać sumie kontrolnej bajtów zapisanych do repozytorium,
// nie samemu identyfikatorowi wpisu w bazie danych.
package core

import (
	"context"
	"testing"

	"danacoconsole/shared"
)

// trescBiezaca odczytuje treść widzianą dziś przez czytelnika pliku. Podgląd
// tekstowy czyta bajty spod odwołania pliku, więc jego wynik jest odpowiedzią na
// pytanie „co naprawdę leży pod tym plikiem teraz".
func trescBiezaca(t *testing.T, zmontowany *Zmontowany, zycie context.Context, kodPliku string) string {
	t.Helper()

	var podglad shared.LibraryFilePreviewResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandLibraryFilePreview,
		shared.LibraryFilePreviewRequest{FileId: kodPliku}, &podglad)

	if podglad.Preview.Kind != shared.LibraryPreviewKindText {
		t.Fatalf("podgląd pliku %s ma rodzaj %q — sprawdzian zakłada treść tekstową",
			kodPliku, podglad.Preview.Kind)
	}
	if podglad.Preview.Text == nil {
		t.Fatalf("podgląd pliku %s nie niesie treści", kodPliku)
	}
	if podglad.Preview.Truncated != nil && *podglad.Preview.Truncated {
		t.Fatalf("podgląd pliku %s został skrócony — sprawdzian porównywałby ogryzek", kodPliku)
	}
	return *podglad.Preview.Text
}

// wgrajPlikTekstowy wnosi do repozytorium plik o zadanej treści tekstowej i oddaje
// jego opis kontraktowy, gotowy do dalszego porównania w sprawdzianach tego pliku.
func wgrajPlikTekstowy(t *testing.T, zmontowany *Zmontowany, zycie context.Context,
	nazwa, tresc string) shared.LibraryFile {
	t.Helper()

	var wynik shared.LibraryFileUploadResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandLibraryFileUpload,
		shared.LibraryFileUploadRequest{
			Name:          nazwa,
			MimeType:      wskaznik("text/plain"),
			ContentBase64: wskaznik(wBase64([]byte(tresc))),
		}, &wynik)
	return wynik.File
}

// TestWgraniePlikuUtrwalaTrescPodWlasnaSumaKontrolna sprawdza wejście do rodziny:
// po wgraniu plik niesie sumę kontrolną wgranych bajtów, a czytelnik widzi
// dokładnie tę treść.
func TestWgraniePlikuUtrwalaTrescPodWlasnaSumaKontrolna(t *testing.T) {
	zmontowany, zycie, _ := zmontujDoPomiaruSkutku(t)

	tresc := "pierwsza treść repozytorium wiedzy"
	plik := wgrajPlikTekstowy(t, zmontowany, zycie, "notatka.txt", tresc)

	if plik.Checksum == nil || *plik.Checksum != sumaSha256([]byte(tresc)) {
		t.Errorf("plik niesie sumę kontrolną %v, treść ma %s", plik.Checksum, sumaSha256([]byte(tresc)))
	}
	if plik.SizeBytes == nil || *plik.SizeBytes != int64(len(tresc)) {
		t.Errorf("plik niesie rozmiar %v bajtów, treść ma %d", plik.SizeBytes, len(tresc))
	}
	if plik.VersionId == nil || *plik.VersionId == "" {
		t.Error("wgranie nie zostawiło wersji bieżącej — Versioning Panel zaczyna od pustki")
	}
	if widziana := trescBiezaca(t, zmontowany, zycie, plik.Id); widziana != tresc {
		t.Errorf("czytelnik widzi %q, wgrano %q", widziana, tresc)
	}
}

// TestPrzywrocenieWersjiWracaDoJejTresciAWskaznikNaNiaWskazuje mierzy trzy rzeczy
// naraz: czytelnik po przywróceniu widzi treść przywróconej wersji, suma pliku
// opisuje tę treść, a wskaźnik wersji bieżącej wskazuje wersję, o którą proszono.
func TestPrzywrocenieWersjiWracaDoJejTresciAWskaznikNaNiaWskazuje(t *testing.T) {
	zmontowany, zycie, _ := zmontujDoPomiaruSkutku(t)

	pierwsza := "treść pierwsza — ta, do której się wraca"
	druga := "treść druga — dołożona po pierwszej"

	plik := wgrajPlikTekstowy(t, zmontowany, zycie, "dokument.txt", pierwsza)
	wersjaPierwsza := *plik.VersionId

	var dolozona shared.LibraryVersionAddResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandLibraryVersionAdd,
		shared.LibraryVersionAddRequest{
			FileId:        plik.Id,
			ContentBase64: wskaznik(wBase64([]byte(druga))),
			Label:         wskaznik("po korekcie"),
		}, &dolozona)

	// Dołożenie wersji zmienia treść bieżącą — inaczej sprawdzian niżej przechodziłby zawsze.
	if widziana := trescBiezaca(t, zmontowany, zycie, plik.Id); widziana != druga {
		t.Fatalf("po dołożeniu wersji czytelnik widzi %q, dołożono %q", widziana, druga)
	}
	if dolozona.Version.Checksum == nil || *dolozona.Version.Checksum != sumaSha256([]byte(druga)) {
		t.Errorf("dołożona wersja niesie sumę %v, jej treść ma %s",
			dolozona.Version.Checksum, sumaSha256([]byte(druga)))
	}

	var historia shared.LibraryVersionListResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandLibraryVersionList,
		shared.LibraryVersionListRequest{FileId: plik.Id}, &historia)
	if len(historia.Versions) != 2 {
		t.Fatalf("historia pliku ma %d wersji, zapisano 2", len(historia.Versions))
	}

	var przywrocona shared.LibraryVersionRestoreResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandLibraryVersionRestore,
		shared.LibraryVersionRestoreRequest{FileId: plik.Id, VersionId: wersjaPierwsza}, &przywrocona)

	if widziana := trescBiezaca(t, zmontowany, zycie, plik.Id); widziana != pierwsza {
		t.Errorf("po przywróceniu czytelnik widzi %q, przywracano %q", widziana, pierwsza)
	}
	if przywrocona.File.Checksum == nil || *przywrocona.File.Checksum != sumaSha256([]byte(pierwsza)) {
		t.Errorf("plik po przywróceniu niesie sumę %v, przywrócona treść ma %s",
			przywrocona.File.Checksum, sumaSha256([]byte(pierwsza)))
	}
	if przywrocona.File.VersionId == nil || *przywrocona.File.VersionId != wersjaPierwsza {
		t.Errorf("wskaźnik wersji bieżącej pokazuje %v, przywracano %s",
			przywrocona.File.VersionId, wersjaPierwsza)
	}
}

// TestPrzywroceniePilnujeGranicPlikuIZostawiaTrescNietknieta sprawdza, że wersja
// szukana jest wyłącznie w obrębie wskazanego pliku, więc żądanie z kodem wersji
// cudzego pliku ma odmówić, a treść obu plików ma pozostać nietknięta.
func TestPrzywroceniePilnujeGranicPlikuIZostawiaTrescNietknieta(t *testing.T) {
	zmontowany, zycie, _ := zmontujDoPomiaruSkutku(t)

	swoj := "treść pliku, którego dotyczy żądanie"
	cudzy := "treść pliku obcego — nie ma prawa tu trafić"

	plikSwoj := wgrajPlikTekstowy(t, zmontowany, zycie, "swoj.txt", swoj)
	plikCudzy := wgrajPlikTekstowy(t, zmontowany, zycie, "cudzy.txt", cudzy)

	odmowa := wykonajOdmowna(t, zmontowany, zycie, shared.CommandLibraryVersionRestore,
		shared.LibraryVersionRestoreRequest{
			FileId:    plikSwoj.Id,
			VersionId: *plikCudzy.VersionId,
		})
	if odmowa.Code != shared.ErrorCodeNotFound {
		t.Errorf("odmowa niesie kod %q, wersja spoza historii pliku to %q",
			odmowa.Code, shared.ErrorCodeNotFound)
	}

	if widziana := trescBiezaca(t, zmontowany, zycie, plikSwoj.Id); widziana != swoj {
		t.Errorf("po odmowie plik widzi %q zamiast swojej treści %q", widziana, swoj)
	}
	if widziana := trescBiezaca(t, zmontowany, zycie, plikCudzy.Id); widziana != cudzy {
		t.Errorf("plik obcy widzi %q zamiast %q", widziana, cudzy)
	}
}

// TestTrescWgranaSciezkaNieIdzieZaNadpisanymPlikiemZrodlowym sprawdza, że ścieżka
// źródłowa jest jedynie źródłem bajtów w chwili wgrania, nie miejscem ich
// składowania: nadpisanie pliku na dysku nie zmienia treści już utrwalonej.
func TestTrescWgranaSciezkaNieIdzieZaNadpisanymPlikiemZrodlowym(t *testing.T) {
	zmontowany, zycie, _ := zmontujDoPomiaruSkutku(t)

	pierwotna := "treść, którą Operator wskazał ścieżką"
	sciezka := t.TempDir() + "/zrodlo.txt"
	zapiszPlikSprawdzianu(t, sciezka, []byte(pierwotna))

	var wgrany shared.LibraryFileUploadResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandLibraryFileUpload,
		shared.LibraryFileUploadRequest{
			Name:       "ze-sciezki.txt",
			MimeType:   wskaznik("text/plain"),
			SourcePath: &sciezka,
		}, &wgrany)
	wersjaPierwotna := *wgrany.File.VersionId

	// Nadpisanie pliku źródłowego po wgraniu nie może zmienić treści już utrwalonej.
	zapiszPlikSprawdzianu(t, sciezka, []byte("treść podmieniona za plecami rdzenia"))

	if widziana := trescBiezaca(t, zmontowany, zycie, wgrany.File.Id); widziana != pierwotna {
		t.Fatalf("treść poszła za nadpisanym plikiem źródłowym: widać %q, wgrano %q", widziana, pierwotna)
	}

	// Ta sama próba dla wersji dołożonej ścieżką — druga droga, ten sam wymóg.
	druga := "treść drugiej wersji, też wskazana ścieżką"
	sciezkaDrugiej := t.TempDir() + "/zrodlo-drugiej.txt"
	zapiszPlikSprawdzianu(t, sciezkaDrugiej, []byte(druga))

	wykonajUdana(t, zmontowany, zycie, shared.CommandLibraryVersionAdd,
		shared.LibraryVersionAddRequest{FileId: wgrany.File.Id, SourcePath: &sciezkaDrugiej}, nil)
	zapiszPlikSprawdzianu(t, sciezkaDrugiej, []byte("i to też podmienione"))

	if widziana := trescBiezaca(t, zmontowany, zycie, wgrany.File.Id); widziana != druga {
		t.Errorf("wersja dołożona ścieżką poszła za nadpisaniem: widać %q, dołożono %q", widziana, druga)
	}

	wykonajUdana(t, zmontowany, zycie, shared.CommandLibraryVersionRestore,
		shared.LibraryVersionRestoreRequest{FileId: wgrany.File.Id, VersionId: wersjaPierwotna}, nil)
	if widziana := trescBiezaca(t, zmontowany, zycie, wgrany.File.Id); widziana != pierwotna {
		t.Errorf("przywrócenie wersji wgranej ścieżką dało %q, utrwalono %q", widziana, pierwotna)
	}
}

// TestWersjaZnacznikNieWymazujeTresciBiezacej sprawdza, że wersja założona bez
// nowej treści zachowuje treść bieżącą pliku, ponieważ znacznik ma opisywać stan
// istniejący, a nie zastępować go pustką.
func TestWersjaZnacznikNieWymazujeTresciBiezacej(t *testing.T) {
	zmontowany, zycie, _ := zmontujDoPomiaruSkutku(t)

	tresc := "treść, na której stawiany jest kamień milowy"
	plik := wgrajPlikTekstowy(t, zmontowany, zycie, "kamien.txt", tresc)

	var znacznik shared.LibraryVersionAddResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandLibraryVersionAdd,
		shared.LibraryVersionAddRequest{FileId: plik.Id, Label: wskaznik("wydanie przyjęte")}, &znacznik)

	if znacznik.Version.Checksum == nil || *znacznik.Version.Checksum != sumaSha256([]byte(tresc)) {
		t.Errorf("wersja-znacznik niesie sumę %v, treść bieżąca ma %s",
			znacznik.Version.Checksum, sumaSha256([]byte(tresc)))
	}
	if widziana := trescBiezaca(t, zmontowany, zycie, plik.Id); widziana != tresc {
		t.Errorf("po założeniu znacznika czytelnik widzi %q zamiast %q", widziana, tresc)
	}

	// Przywrócenie znacznika zostawia treść — to ta sama treść.
	wykonajUdana(t, zmontowany, zycie, shared.CommandLibraryVersionRestore,
		shared.LibraryVersionRestoreRequest{FileId: plik.Id, VersionId: znacznik.Version.Id}, nil)
	if widziana := trescBiezaca(t, zmontowany, zycie, plik.Id); widziana != tresc {
		t.Errorf("przywrócenie znacznika wymazało treść: widać %q zamiast %q", widziana, tresc)
	}
}

// TestDolozenieWersjiDoNieznanegoPlikuNiczegoNieZaklada domyka rodzinę od strony
// wskazania: dołożenie wersji do pliku, którego nie ma, ma odmówić i nie
// zostawić po sobie ani pliku, ani wersji. Ciche powodzenie założyłoby historię
// bytu, którego nie ma.
func TestDolozenieWersjiDoNieznanegoPlikuNiczegoNieZaklada(t *testing.T) {
	zmontowany, zycie, _ := zmontujDoPomiaruSkutku(t)

	odmowa := wykonajOdmowna(t, zmontowany, zycie, shared.CommandLibraryVersionAdd,
		shared.LibraryVersionAddRequest{
			FileId:        "plik-ktorego-nie-ma",
			ContentBase64: wskaznik(wBase64([]byte("treść bez pliku"))),
		})
	if odmowa.Code != shared.ErrorCodeNotFound {
		t.Errorf("odmowa niesie kod %q, plik nieznany to %q", odmowa.Code, shared.ErrorCodeNotFound)
	}

	var wykaz shared.LibraryFileListResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandLibraryFileList,
		shared.LibraryFileListRequest{}, &wykaz)
	if len(wykaz.Files) != 0 {
		t.Errorf("po odmowie w repozytorium leży %d plików", len(wykaz.Files))
	}
}
