// Odpowiedzialność pliku: raport badania (Report Builder), jego eksport i
// przestrzeń badania — metody `ZbudujRaport`, `WyeksportujRaport`,
// `UstawPrzestrzen` na `*adapterBadan`. Źródła i ustalenia leżą w
// `adapter_modul_badania.go` — jeden adapter, dwa pliki wedle
// odpowiedzialności.
//
// Rdzeń nie prowadzi badań i nie zmyśla treści. Sekcje redagowane wprost
// (`z.Sections`) przechodzą bez zmiany — to redakcja Operatora, bez modelu.
// `z.FindingIds` bez sekcji własnych uruchamia redakcję modelem: most
// `zapytajModel` streszcza zebrane ustalenia w jedną sekcję nadrzędną
// („Streszczenie ustaleń"), po której idą sekcje szczegółowe — po jednej na
// ustalenie, z treścią ustalenia skopiowaną dosłownie dla śladu (tytuł to kod
// ustalenia, bo kontrakt nie daje ustaleniu własnego tytułu). Streszczenie
// bierze domyślny czynny kanał rejestru (`domyslnyKanalBadania`), a brak
// czynnego kanału to odmowa wprost.
//
// `WyeksportujRaport` zapisuje ślad eksportu w bycie trwałym
// `eksport_raportu_badania`, bo `LibraryFileId`/`Path`/`SizeBytes` muszą
// przeżyć wywołanie. Pliku na dysk rdzeń nie zapisuje — nie ma magazynu blobów
// (ten sam brak, co w module Library). Odpowiedź zostawia `LibraryFileId`,
// `Path` i `SizeBytes` puste, zamiast udawać wytworzony plik.
//
// `UstawPrzestrzen` operuje na przestrzeni jednowierszowej — kontrakt
// `research.workspace.set` nie niesie identyfikatora okna ani sesji
// (sprawdzone w `dane/badania_raport.go`), więc jedna przestrzeń na
// instalację jest jedynym uczciwym odczytaniem.
package core

import (
	"context"
	"errors"
	"strconv"
	"strings"

	"danacoconsole/server/internal/dane"
	"danacoconsole/server/internal/protocol"
	"danacoconsole/shared"
)

// tytulStreszczeniaRaportu nazywa sekcję nadrzędną złożoną przez model ze
// zredagowanych ustaleń — odróżnia ją od sekcji szczegółowych, których tytułem
// jest kod ustalenia.
const tytulStreszczeniaRaportu = "Streszczenie ustaleń"

// ZbudujRaport zakłada raport nowy albo rozbudowuje zastany o wskazanym
// identyfikatorze i zwraca stan po złożeniu — patrz nagłówek pliku o źródle
// sekcji.
func (a *adapterBadan) ZbudujRaport(ctx context.Context,
	z shared.ResearchReportBuildRequest) (shared.ResearchReportBuildResponse, error) {

	if z.WindowId == "" {
		return shared.ResearchReportBuildResponse{}, bladWskazaniaBadan("report.build bez okna")
	}
	kod := wartoscTekstu(z.ReportId)
	if kod == "" {
		kod = nowyIdentyfikator(przedrostekRaportuBadania)
	}
	tytul, err := a.tytulRaportu(ctx, kod, z.Title)
	if err != nil {
		return shared.ResearchReportBuildResponse{}, err
	}

	sekcje, err := a.sekcjeDoZapisu(ctx, z)
	if err != nil {
		return shared.ResearchReportBuildResponse{}, err
	}

	zapisany, err := a.repozytorium.ZapiszRaport(ctx,
		dane.RaportBadania{Kod: kod, Okno: z.WindowId, Tytul: tytul}, sekcje)
	if err != nil {
		return shared.ResearchReportBuildResponse{}, bladRaportuBadania(err)
	}

	raport, err := a.przelozRaport(ctx, zapisany)
	if err != nil {
		return shared.ResearchReportBuildResponse{}, err
	}
	return shared.ResearchReportBuildResponse{Report: raport}, nil
}

// tytulRaportu bierze tytuł z żądania, a przy jego braku zachowuje tytuł
// zastany (raport rozbudowywany) — zapis w `dane` nadpisuje `tytul` zawsze
// wartością przekazaną, więc utrzymanie starego tytułu jest sprawą adaptera.
func (a *adapterBadan) tytulRaportu(ctx context.Context, kod string, tytul *string) (string, error) {
	if tytul != nil {
		return *tytul, nil
	}
	zastany, err := a.repozytorium.Raport(ctx, kod)
	if errors.Is(err, dane.ErrBrakWiersza) {
		return "", nil
	}
	if err != nil {
		return "", bladRaportuBadania(err)
	}
	return zastany.Tytul, nil
}

// sekcjeDoZapisu wybiera źródło sekcji: redakcja wprost ma pierwszeństwo,
// a bez niej sekcje powstają po jednej na ustalenie z `FindingIds`.
func (a *adapterBadan) sekcjeDoZapisu(ctx context.Context,
	z shared.ResearchReportBuildRequest) ([]dane.SekcjaRaportu, error) {

	if z.Sections != nil {
		sekcje := make([]dane.SekcjaRaportu, 0, len(z.Sections))
		for numer, sekcja := range z.Sections {
			kod := sekcja.Id
			if kod == "" {
				kod = nowyIdentyfikator(przedrostekSekcjiRaportu)
			}
			kolejnosc := numer + 1
			if sekcja.Order != nil {
				kolejnosc = *sekcja.Order
			}
			sekcje = append(sekcje, dane.SekcjaRaportu{
				Kod: kod, Tytul: sekcja.Title, Tresc: sekcja.Content,
				Kolejnosc: kolejnosc, UstalenieKody: sekcja.FindingIds,
			})
		}
		return sekcje, nil
	}
	if len(z.FindingIds) == 0 {
		return []dane.SekcjaRaportu{}, nil
	}

	ustalenia := make([]dane.UstalenieBadania, 0, len(z.FindingIds))
	for _, kodUstalenia := range z.FindingIds {
		ustalenie, err := a.repozytorium.Ustalenie(ctx, kodUstalenia)
		if errors.Is(err, dane.ErrBrakWiersza) {
			return nil, bladNieznanegoUstaleniaRaportu(kodUstalenia)
		}
		if err != nil {
			return nil, bladRaportuBadania(err)
		}
		ustalenia = append(ustalenia, ustalenie)
	}

	// Redakcja modelem przed zapisem sekcji szczegółowych: brak czynnego kanału
	// odmawia całej budowie, zamiast zapisać raport bez streszczenia.
	streszczenie, err := a.sekcjaStreszczenia(ctx, z.WindowId, ustalenia)
	if err != nil {
		return nil, err
	}

	sekcje := make([]dane.SekcjaRaportu, 0, len(ustalenia)+1)
	sekcje = append(sekcje, streszczenie)
	for numer, ustalenie := range ustalenia {
		sekcje = append(sekcje, dane.SekcjaRaportu{
			Kod: nowyIdentyfikator(przedrostekSekcjiRaportu), Tytul: ustalenie.Kod,
			Tresc: ustalenie.Tresc, TrescOdwolanie: ustalenie.TrescOdwolanie,
			Kolejnosc: numer + 2, UstalenieKody: []string{ustalenie.Kod},
		})
	}
	return sekcje, nil
}

// sekcjaStreszczenia redaguje modelem sekcję nadrzędną raportu ze zebranych
// ustaleń. Operacja wymaga modelu: bierze domyślny czynny kanał
// rejestru, składa z treści ustaleń polecenie redakcji i zwraca sekcję z
// odpowiedzią modelu. Brak kanału albo błąd kanału to odmowa wprost,
// nie streszczenie zmyślone bez modelu.
func (a *adapterBadan) sekcjaStreszczenia(ctx context.Context, okno string,
	ustalenia []dane.UstalenieBadania) (dane.SekcjaRaportu, error) {

	kanal, err := a.domyslnyKanalBadania()
	if err != nil {
		return dane.SekcjaRaportu{}, err
	}
	tresc, err := a.zapytajModel(ctx, okno, kanal, poleceniStreszczenia(ustalenia))
	if err != nil {
		return dane.SekcjaRaportu{}, err
	}
	kody := make([]string, 0, len(ustalenia))
	for _, u := range ustalenia {
		kody = append(kody, u.Kod)
	}
	streszczenie := tresc
	return dane.SekcjaRaportu{
		Kod:           nowyIdentyfikator(przedrostekSekcjiRaportu),
		Tytul:         tytulStreszczeniaRaportu,
		Tresc:         &streszczenie,
		Kolejnosc:     1,
		UstalenieKody: kody,
	}, nil
}

// poleceniStreszczenia składa polecenie redakcji z treści zebranych ustaleń.
// Wiąże model regułą „tylko z podanych ustaleń", żeby streszczenie zostało
// redakcją zebranych bytów, nie dopisaniem faktów spoza nich.
func poleceniStreszczenia(ustalenia []dane.UstalenieBadania) string {
	var b strings.Builder
	b.WriteString("Zredaguj zwięzłą sekcję raportu badawczego streszczającą poniższe ustalenia ")
	b.WriteString("w jeden spójny tekst. Opieraj się WYŁĄCZNIE na podanych ustaleniach — ")
	b.WriteString("nie dodawaj faktów spoza nich.\n\nUstalenia:\n")
	for i, u := range ustalenia {
		b.WriteString(strconv.Itoa(i + 1))
		b.WriteString(". ")
		b.WriteString(trescUstaleniaDoPromptu(u))
		b.WriteString("\n")
	}
	return b.String()
}

// trescUstaleniaDoPromptu wybiera treść ustalenia do polecenia: treść wprost,
// a przy jej braku treść odwołania. Puste ustalenie idzie pustym wierszem —
// most nie dopisuje niczego, czego w ustaleniu nie ma.
func trescUstaleniaDoPromptu(u dane.UstalenieBadania) string {
	if u.Tresc != nil && strings.TrimSpace(*u.Tresc) != "" {
		return *u.Tresc
	}
	if u.TrescOdwolanie != nil {
		return *u.TrescOdwolanie
	}
	return ""
}

// przelozRaport dobudowuje sekcje do wiersza raportu i przekłada całość na
// byt kontraktu.
func (a *adapterBadan) przelozRaport(ctx context.Context, raport dane.RaportBadania) (shared.ResearchReport, error) {
	sekcje, err := a.repozytorium.Sekcje(ctx, raport.ID)
	if err != nil {
		return shared.ResearchReport{}, bladRaportuBadania(err)
	}
	przelozone := make([]shared.ResearchReportSection, 0, len(sekcje))
	for _, sekcja := range sekcje {
		przelozone = append(przelozone, shared.ResearchReportSection{
			Id: sekcja.Kod, Title: sekcja.Tytul, Content: sekcja.Tresc,
			FindingIds: sekcja.UstalenieKody, Order: wskaznikLiczby(sekcja.Kolejnosc),
		})
	}
	return shared.ResearchReport{
		Id: raport.Kod, WindowId: raport.Okno, Title: raport.Tytul,
		Sections: przelozone, UpdatedAt: chwilaBazy(raport.Zaktualizowano),
	}, nil
}

// UstawPrzestrzen nadpisuje jedyną przestrzeń badania instalacji i wymienia
// jej etapy w całości.
func (a *adapterBadan) UstawPrzestrzen(ctx context.Context,
	z shared.ResearchWorkspaceSetRequest) (shared.ResearchWorkspaceSetResponse, error) {

	if z.Scope == "" {
		return shared.ResearchWorkspaceSetResponse{}, bladWskazaniaBadan("workspace.set bez zakresu")
	}
	zakres, etapy, err := a.repozytorium.UstawPrzestrzen(ctx, z.Scope, z.Stages)
	if err != nil {
		return shared.ResearchWorkspaceSetResponse{}, bladRaportuBadania(err)
	}
	return shared.ResearchWorkspaceSetResponse{Scope: zakres, Stages: etapy}, nil
}

// bladRaportuBadania znakuje usterkę kodem kontraktu (wzór: bladBiblioteki).
func bladRaportuBadania(err error) error {
	if err == nil {
		return nil
	}
	return protocol.JakoError(protocol.BladZeZrodla(shared.ErrorCodeInternalError, err))
}

// bladNieznanegoRaportuBadania odróżnia „raportu nie ma" od „zapis się nie powiódł".
func bladNieznanegoRaportuBadania(kod string) error {
	return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeNotFound,
		"moduł Research: raport nie istnieje: "+kod))
}

// bladNieznanegoUstaleniaRaportu odróżnia „ustalenia nie ma" od usterki wewnętrznej,
// gdy `report.build` składa sekcje z samych `FindingIds` (patrz nagłówek pliku).
func bladNieznanegoUstaleniaRaportu(kod string) error {
	return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeNotFound,
		"moduł Research: ustalenie nie istnieje: "+kod))
}
