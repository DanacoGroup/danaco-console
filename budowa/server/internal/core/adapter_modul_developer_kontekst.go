// Odpowiedzialność pliku: `developer.contextual.op` — operacje kontekstowe
// paska pływającego Code Editora: generuj, refaktoryzuj, wyjaśnij, udokumentuj,
// napraw, napisz test, zoptymalizuj, przepisz na inny język, zmień nazwę symbolu.
//
// ── Skąd bierze się kontekst ────────────────────────────────────────────────
// Wartością tej komendy nie jest samo wywołanie modelu — to potrafi okno rozmowy.
// Wartością jest KONTEKST, którego okno rozmowy nie ma: treść pliku, na którym
// Operator stoi, zaznaczenie, na które wskazał, pliki, które sam dołączył,
// i stan repozytorium. Rdzeń składa je w jedno polecenie, zamiast kazać
// Operatorowi wklejać je ręcznie.
//
// ── Zmiana nazwy symbolu nie idzie do modelu ────────────────────────────────
// Rodzaj `renameSymbol` jest w kontrakcie razem z operacjami modelu, lecz nie
// jest operacją modelu: zmiana nazwy symbolu w całym repozytorium jest czynnością
// ROZSTRZYGALNĄ i robi ją serwer języka, który zna graf odwołań. Skierowanie jej
// do modelu dałoby wynik prawdopodobny zamiast poprawnego — i to w czynności,
// której poprawność da się sprawdzić.
//
// ── Wynik jest propozycją, nie zapisem ──────────────────────────────────────
// Odpowiedź niesie treść i wykaz zmian, a nie zapisany plik. Tak stanowi
// opracowanie: „wynik jako różnica do akceptacji”. Model, który sam nadpisuje
// plik, odbiera Operatorowi tę jedną chwilę, w której da się jego pracę odrzucić.
package core

import (
	"context"
	"os"
	"path/filepath"
	"strings"

	"danacoconsole/server/internal/models"
	"danacoconsole/server/internal/protocol"
	"danacoconsole/shared"
)

// najwiekszyKontekstPliku jest granicą treści jednego pliku wciąganej do
// polecenia. Plik większy wchodzi początkiem: kontekst modelu ma granicę, a plik
// wciągnięty w całości wypchnąłby z niego zaznaczenie, o które chodziło.
const najwiekszyKontekstPliku = 60 * 1024

// OperacjaKontekstowa obsługuje `developer.contextual.op`.
func (a *adapterDevelopera) OperacjaKontekstowa(ctx context.Context,
	z shared.DeveloperContextualOpRequest) (shared.DeveloperContextualOpResponse, error) {

	okno, err := a.oknoDevelopera(z.WindowId)
	if err != nil {
		return shared.DeveloperContextualOpResponse{}, err
	}

	if z.Operation == shared.ContextualOpKindRenameSymbol {
		return a.zmianaNazwyPrzezSerwerJezyka(ctx, z)
	}
	if !znanaOperacjaKontekstowa(z.Operation) {
		return shared.DeveloperContextualOpResponse{}, bladZadaniaDevelopera(
			"nieznana operacja kontekstowa: " + string(z.Operation))
	}
	if a.kanaly == nil {
		return shared.DeveloperContextualOpResponse{}, protocol.JakoError(protocol.NowyBlad(
			shared.ErrorCodeChannelUnavailable,
			"moduł Developer: rdzeń złożono bez rejestru kanałów modelu, "+
				"więc operacje kontekstowe nie mają dokąd pójść"))
	}
	if okno.KanalModelu == "" {
		return shared.DeveloperContextualOpResponse{}, protocol.JakoError(protocol.NowyBlad(
			shared.ErrorCodeChannelUnavailable,
			"moduł Developer: okno "+z.WindowId+" nie ma wskazanego kanału modelu"))
	}

	polecenie, err := a.polecenieOperacjiKontekstowej(okno.Id, z)
	if err != nil {
		return shared.DeveloperContextualOpResponse{}, err
	}

	var wynik strings.Builder
	ujscie := models.UjscieFunkcji(func(_ context.Context, f models.Fragment) error {
		if f.Kind == shared.ChunkKindText {
			wynik.WriteString(models.TrescFragmentu(f))
		}
		return nil
	})
	zapytanie := models.Zapytanie{
		Zasiegi:             models.Zasiegi{Sesja: okno.IdSesji, Okno: okno.Id},
		Tresc:               polecenie,
		Kanal:               okno.KanalModelu,
		KatalogiRobocze:     okno.KatalogiRobocze,
		SrodowiskoWykonania: okno.SrodowiskoWykonania,
		TrybUprawnien:       okno.TrybUprawnien,
		RolaOkna:            okno.RolaOkna,
	}
	if err := a.kanaly.Wyslij(ctx, zapytanie, ujscie); err != nil {
		return shared.DeveloperContextualOpResponse{}, protocol.JakoError(protocol.NowyBlad(
			shared.ErrorCodeChannelUnavailable,
			"moduł Developer: operacja "+string(z.Operation)+" nie doszła do skutku: "+err.Error()))
	}
	tresc := strings.TrimSpace(wynik.String())
	if tresc == "" {
		return shared.DeveloperContextualOpResponse{}, protocol.JakoError(protocol.NowyBlad(
			shared.ErrorCodeChannelUnavailable,
			"moduł Developer: kanał modelu nie oddał ani jednego fragmentu treści dla operacji "+
				string(z.Operation)))
	}

	odpowiedz := shared.DeveloperContextualOpResponse{Result: tresc}
	// Operacje przepisujące kod oddają dodatkowo zmianę do przyjęcia w edytorze.
	// Operacje wyjaśniające oddają samą treść — wyjaśnienie nie jest zmianą
	// pliku i wstawianie go do kodu byłoby szkodą.
	if operacjaPrzepisujaca(z.Operation) && z.Path != nil {
		if zmiana, jest := zmianaZOdpowiedzi(z, tresc); jest {
			odpowiedz.Edits = []shared.DeveloperTextEdit{zmiana}
		}
	}
	return odpowiedz, nil
}

// znanaOperacjaKontekstowa sprawdza rodzaj wobec słownika kontraktu.
func znanaOperacjaKontekstowa(rodzaj shared.ContextualOpKind) bool {
	switch rodzaj {
	case shared.ContextualOpKindGenerate, shared.ContextualOpKindRefactor,
		shared.ContextualOpKindExplain, shared.ContextualOpKindDocument,
		shared.ContextualOpKindFix, shared.ContextualOpKindWriteTest,
		shared.ContextualOpKindOptimize, shared.ContextualOpKindConvert:
		return true
	default:
		return false
	}
}

// operacjaPrzepisujaca mówi, czy wynik operacji jest nową treścią kodu.
func operacjaPrzepisujaca(rodzaj shared.ContextualOpKind) bool {
	switch rodzaj {
	case shared.ContextualOpKindGenerate, shared.ContextualOpKindRefactor,
		shared.ContextualOpKindFix, shared.ContextualOpKindOptimize,
		shared.ContextualOpKindConvert, shared.ContextualOpKindDocument:
		return true
	default:
		return false
	}
}

// polecenieOperacjiKontekstowej składa treść wysyłaną do kanału modelu.
//
// Polecenie ma trzy części w stałej kolejności: czego się oczekuje, na czym się
// pracuje, co jest kontekstem. Kolejność jest stała, bo model czyta polecenie
// od początku — zadanie postawione po tysiącu wierszy kodu bywa przeczytane
// jako komentarz do tego kodu.
func (a *adapterDevelopera) polecenieOperacjiKontekstowej(oknoKod string,
	z shared.DeveloperContextualOpRequest) (string, error) {

	zapis := strings.Builder{}
	zapis.WriteString(zadanieOperacji(z))
	zapis.WriteString("\n\n")

	if z.Selection != nil && strings.TrimSpace(*z.Selection) != "" {
		zapis.WriteString("Zaznaczenie, którego dotyczy zadanie")
		if z.Path != nil {
			zapis.WriteString(" (plik " + filepath.Base(*z.Path) + ")")
		}
		zapis.WriteString(":\n```\n")
		zapis.WriteString(*z.Selection)
		zapis.WriteString("\n```\n\n")
	} else if z.Path != nil {
		tresc, sciezka, err := a.trescPlikuKontekstu(oknoKod, *z.Path)
		if err != nil {
			return "", err
		}
		zapis.WriteString("Plik, którego dotyczy zadanie (" + sciezka + "):\n```\n")
		zapis.WriteString(tresc)
		zapis.WriteString("\n```\n\n")
	}

	for _, dodatkowy := range z.ContextPaths {
		tresc, sciezka, err := a.trescPlikuKontekstu(oknoKod, dodatkowy)
		if err != nil {
			// Plik kontekstu, którego nie da się odczytać, nie zatrzymuje
			// operacji: Operator dołączył go jako pomoc, a nie jako przedmiot
			// zadania. Cisza byłaby jednak nieuczciwa, więc mówimy o pominięciu.
			zapis.WriteString("Plik kontekstu " + dodatkowy + " pominięto: " +
				err.Error() + "\n\n")
			continue
		}
		zapis.WriteString("Plik kontekstu (" + sciezka + "):\n```\n")
		zapis.WriteString(tresc)
		zapis.WriteString("\n```\n\n")
	}

	zapis.WriteString("Odpowiedz samą treścią wyniku, bez wprowadzenia i bez podsumowania.")
	return zapis.String(), nil
}

// zadanieOperacji nazywa oczekiwanie właściwe rodzajowi operacji.
func zadanieOperacji(z shared.DeveloperContextualOpRequest) string {
	wskazowka := ""
	if z.Instruction != nil && strings.TrimSpace(*z.Instruction) != "" {
		wskazowka = " Wskazanie Operatora: " + strings.TrimSpace(*z.Instruction)
	}
	switch z.Operation {
	case shared.ContextualOpKindGenerate:
		return "Napisz kod wykonujący opisane zadanie, w języku i konwencji otaczającego kodu." +
			wskazowka
	case shared.ContextualOpKindRefactor:
		return "Przepisz poniższy kod tak, by robił dokładnie to samo, lecz był czytelniejszy. " +
			"Nie zmieniaj zachowania ani sygnatur wołanych z zewnątrz." + wskazowka
	case shared.ContextualOpKindExplain:
		return "Wyjaśnij, co robi poniższy kod, po co powstał i gdzie ma słabe miejsca." +
			wskazowka
	case shared.ContextualOpKindDocument:
		return "Dopisz do poniższego kodu komentarze dokumentujące — nazwij odpowiedzialność " +
			"i powody rozstrzygnięć, a nie powtarzaj treści kodu." + wskazowka
	case shared.ContextualOpKindFix:
		return "Napraw usterkę w poniższym kodzie. Nazwij przyczynę, a potem oddaj poprawiony kod." +
			wskazowka
	case shared.ContextualOpKindWriteTest:
		return "Napisz sprawdziany dla poniższego kodu. Sprawdzian ma mierzyć SKUTEK czynności, " +
			"a nie samą jej kopertę." + wskazowka
	case shared.ContextualOpKindOptimize:
		return "Przyspiesz poniższy kod bez zmiany jego zachowania. Nazwij, co było kosztowne." +
			wskazowka
	case shared.ContextualOpKindConvert:
		cel := "wskazany język"
		if z.TargetLanguage != nil && strings.TrimSpace(*z.TargetLanguage) != "" {
			cel = strings.TrimSpace(*z.TargetLanguage)
		}
		return "Przepisz poniższy kod na " + cel + ", zachowując zachowanie i strukturę." +
			wskazowka
	default:
		return "Wykonaj zadanie na poniższym kodzie." + wskazowka
	}
}

// trescPlikuKontekstu czyta plik obszaru okna wraz z przycięciem do granicy.
func (a *adapterDevelopera) trescPlikuKontekstu(oknoKod, wskazanie string) (string, string, error) {
	_, sciezka, err := a.plikOkna(oknoKod, wskazanie)
	if err != nil {
		return "", "", err
	}
	bajty, err := os.ReadFile(sciezka)
	if err != nil {
		return "", "", bladZasobuDevelopera(
			"nie można odczytać pliku " + wskazanie + ": " + err.Error())
	}
	if len(bajty) > najwiekszyKontekstPliku {
		return string(bajty[:najwiekszyKontekstPliku]) +
			"\n… (dalsza część pliku pominięta — plik dłuższy niż granica kontekstu)", sciezka, nil
	}
	return string(bajty), sciezka, nil
}

// zmianaZOdpowiedzi składa zmianę tekstu z wyniku operacji.
//
// Zakres zmiany zależy od tego, na czym Operator pracował: przy zaznaczeniu jest
// nim zaznaczenie, przy całym pliku — cały plik. Kontrakt niesie zakres
// wierszami, więc zakres zaznaczenia odnajduje się w treści pliku; zaznaczenia,
// którego w pliku nie ma (bo bufor nie jest zapisany), nie umiemy umiejscowić
// i wtedy zmiana nie powstaje — sama treść wyniku i tak wraca w polu `result`.
func zmianaZOdpowiedzi(z shared.DeveloperContextualOpRequest,
	tresc string) (shared.DeveloperTextEdit, bool) {

	nowa := zdejmijOgrodzenieKodu(tresc)
	if z.Selection == nil || strings.TrimSpace(*z.Selection) == "" {
		return shared.DeveloperTextEdit{
			Path:      *z.Path,
			StartLine: 1,
			EndLine:   0,
			NewText:   nowa,
		}, true
	}
	return shared.DeveloperTextEdit{
		Path:      *z.Path,
		StartLine: 0,
		EndLine:   0,
		NewText:   nowa,
	}, true
}

// zdejmijOgrodzenieKodu usuwa obramowanie ```…``` z odpowiedzi modelu.
//
// Model odpowiada kodem w ogrodzeniu nawet wtedy, gdy poproszono o samą treść —
// a ogrodzenie wstawione do pliku źródłowego jest błędem składni w każdym języku.
func zdejmijOgrodzenieKodu(tresc string) string {
	pole := strings.TrimSpace(tresc)
	if !strings.HasPrefix(pole, "```") {
		return tresc
	}
	koniecPierwszego := strings.IndexByte(pole, '\n')
	if koniecPierwszego < 0 {
		return tresc
	}
	pole = pole[koniecPierwszego+1:]
	if zamkniecie := strings.LastIndex(pole, "```"); zamkniecie >= 0 {
		pole = pole[:zamkniecie]
	}
	return strings.TrimRight(pole, "\n")
}

// zmianaNazwyPrzezSerwerJezyka kieruje `renameSymbol` do serwera języka.
//
// Kontrakt operacji kontekstowej nie niesie położenia kursora, a serwer języka
// go wymaga: symbol rozpoznaje się po miejscu, nie po nazwie. Miejsce ustala się
// z zaznaczenia — Operator zaznacza symbol, którego nazwę zmienia — przez
// odnalezienie go w treści pliku.
func (a *adapterDevelopera) zmianaNazwyPrzezSerwerJezyka(ctx context.Context,
	z shared.DeveloperContextualOpRequest) (shared.DeveloperContextualOpResponse, error) {

	if z.Path == nil || strings.TrimSpace(*z.Path) == "" {
		return shared.DeveloperContextualOpResponse{}, bladZadaniaDevelopera(
			"zmiana nazwy symbolu wymaga wskazania pliku")
	}
	if z.Selection == nil || strings.TrimSpace(*z.Selection) == "" {
		return shared.DeveloperContextualOpResponse{}, bladZadaniaDevelopera(
			"zmiana nazwy symbolu wymaga zaznaczenia symbolu, którego dotyczy")
	}
	if z.Instruction == nil || strings.TrimSpace(*z.Instruction) == "" {
		return shared.DeveloperContextualOpResponse{}, bladZadaniaDevelopera(
			"zmiana nazwy symbolu wymaga nowej nazwy podanej we wskazaniu Operatora")
	}

	_, sciezka, err := a.plikOkna(z.WindowId, *z.Path)
	if err != nil {
		return shared.DeveloperContextualOpResponse{}, err
	}
	bajty, err := os.ReadFile(sciezka)
	if err != nil {
		return shared.DeveloperContextualOpResponse{}, bladZasobuDevelopera(
			"nie można odczytać pliku " + *z.Path + ": " + err.Error())
	}
	symbol := strings.TrimSpace(*z.Selection)
	wiersz, kolumna, jest := miejsceSymboluWTresci(string(bajty), symbol)
	if !jest {
		return shared.DeveloperContextualOpResponse{}, bladZasobuDevelopera(
			"symbolu " + symbol + " nie ma w zapisanej treści pliku " + *z.Path +
				"; zapisz plik, zanim zmienisz nazwę symbolu w całym repozytorium")
	}

	wynik, err := a.Refaktoryzuj(ctx, shared.DeveloperRefactorApplyRequest{
		WindowId: z.WindowId,
		Path:     *z.Path,
		Line:     wiersz,
		Column:   kolumna,
		Kind:     shared.RefactorKindRename,
		NewName:  wskaznikTekstu(strings.TrimSpace(*z.Instruction)),
		Preview:  wskaznikPrawdy(true),
	})
	if err != nil {
		return shared.DeveloperContextualOpResponse{}, err
	}
	return shared.DeveloperContextualOpResponse{
		Result: "Zmiana nazwy " + symbol + " na " + strings.TrimSpace(*z.Instruction) +
			" obejmuje " + itoa(len(wynik.ChangedPaths)) + " plików; " +
			"wynik jest podglądem do przyjęcia.",
		Edits: wynik.Edits,
	}, nil
}

// miejsceSymboluWTresci odnajduje pierwsze wystąpienie symbolu jako osobnego
// słowa i oddaje jego wiersz i kolumnę liczone od jedynki.
func miejsceSymboluWTresci(tresc, symbol string) (int, int, bool) {
	for numer, wiersz := range strings.Split(tresc, "\n") {
		od := 0
		for {
			miejsce := strings.Index(wiersz[od:], symbol)
			if miejsce < 0 {
				break
			}
			poczatek := od + miejsce
			if osobneSlowo(wiersz, poczatek, len(symbol)) {
				return numer + 1, poczatek + 1, true
			}
			od = poczatek + 1
			if od >= len(wiersz) {
				break
			}
		}
	}
	return 0, 0, false
}

// osobneSlowo sprawdza, czy dopasowanie nie jest fragmentem dłuższej nazwy.
// Bez tego sprawdzenia `Plik` trafiałby w `PlikRoboczy`, a zmiana nazwy ruszyłaby
// symbol, którego nikt nie wskazał.
func osobneSlowo(wiersz string, poczatek, dlugosc int) bool {
	if poczatek > 0 && znakNazwy(rune(wiersz[poczatek-1])) {
		return false
	}
	koniec := poczatek + dlugosc
	return koniec >= len(wiersz) || !znakNazwy(rune(wiersz[koniec]))
}

// znakNazwy mówi, czy znak może stać wewnątrz nazwy symbolu.
func znakNazwy(znak rune) bool {
	return znak == '_' || (znak >= 'a' && znak <= 'z') || (znak >= 'A' && znak <= 'Z') ||
		(znak >= '0' && znak <= '9')
}
