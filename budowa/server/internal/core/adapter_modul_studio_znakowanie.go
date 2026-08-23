// Odpowiedzialność pliku: ZNAKOWANIE FRAGMENTÓW — wyróżnienie barwą, znacznik
// własny Operatora i propozycja zmiany na marginesie — wraz z wykazem do
// przejścia, rozstrzygnięciem propozycji i warsztatem rodzajów znaczników.
//
// ── Trzy różne byty, których nie wolno zlać ─────────────────────────────────
// Właściciel każe odróżnić je w oknie wyraźnie, więc rdzeń odróżnia je w danych:
//
//	WYRÓŻNIENIE (`highlight`) — cecha POSTACI dokumentu. Barwa tła fragmentu
//	    idzie tą samą drogą co reszta formatowania, czyli w runy drzewa postaci;
//	    wiersz znakowania jest wykazem do przejścia, nie drugim miejscem, w którym
//	    trzymana jest barwa.
//	ZNACZNIK (`mark`) — nazwa własna Operatora („do sprawdzenia", „wymaga
//	    źródła", „gotowe"). Treści dokumentu NIE rusza.
//	PROPOZYCJA (`suggestion`) — brzmienie proponowane, stojące NA MARGINESIE.
//	    Nie jest zmianą śledzoną, bo tamta jest już w treści, i nie jest
//	    komentarzem, bo tamten nie niesie brzmienia. Wchodzi do treści dopiero
//	    decyzją Operatora — i wtedy odkłada się jako zmiana śledzona jej autora.
//
// ── Dlaczego propozycja jest jedyną drogą wykonawcy do fragmentu pod blokadą ──
// Blokadę zdejmuje wyłącznie Operator. Wykonawca, który uzna, że fragment wymaga
// zmiany, zakłada propozycję — i dlatego założenie PROPOZYCJI na fragmencie
// zablokowanym jest dozwolone, a wniesienie jej do treści już nie: decyzja
// należy do Operatora, a jego blokada o zasięgu `model` nie wiąże.
package core

import (
	"context"
	"errors"
	"sort"
	"strconv"
	"strings"

	"danacoconsole/server/internal/dane"
	"danacoconsole/server/internal/protocol"
	"danacoconsole/shared"
)

// DodajZnakowanie obsługuje `studio.markup.add`.
func (a *adapterStudia) DodajZnakowanie(ctx context.Context,
	z shared.StudioMarkupAddRequest) (shared.StudioMarkupAddResponse, error) {

	stan, err := a.postacWczytaj(ctx, z.DocumentId)
	if err != nil {
		return shared.StudioMarkupAddResponse{}, err
	}
	skladnica, err := a.kontrolaSkladnica()
	if err != nil {
		return shared.StudioMarkupAddResponse{}, err
	}
	if err := znakowanieSprawdzRodzaj(z.Kind); err != nil {
		return shared.StudioMarkupAddResponse{}, err
	}
	if z.RangeEnd < z.RangeStart {
		return shared.StudioMarkupAddResponse{}, bladWskazaniaStudio(
			"znakowanie o zakresie odwróconym: koniec " + strconv.Itoa(z.RangeEnd) +
				" leży przed początkiem " + strconv.Itoa(z.RangeStart))
	}
	wykonawca := kontrolaRozpoznajWykonawce(ctx, kontrolaPodpisZadania{
		Author: z.Author, AgentId: z.AgentId, AgentName: z.AgentName, SubagentId: z.SubagentId,
	})

	// Rodzaj znacznika własnego musi istnieć w wykazie rodzajów. Znacznik o
	// nazwie spoza wykazu byłby znacznikiem, którego Operator nie odfiltruje
	// i którego barwy nikt nie zna.
	if z.Kind == shared.StudioMarkupKindMark {
		if !kontrolaTekstNiepusty(z.MarkType) {
			return shared.StudioMarkupAddResponse{}, bladWskazaniaStudio(
				"znacznik własny bez wskazania rodzaju — rodzaje zakłada " +
					"studio.markup.type.save, a wykaz oddaje studio.markup.type.list")
		}
		if _, err := skladnica.RodzajZnacznika(ctx, *z.MarkType); err != nil {
			return shared.StudioMarkupAddResponse{}, dziennikBrakWiersza(err,
				"wykaz rodzajów znaczników nie zna rodzaju „"+*z.MarkType+"”")
		}
	}
	if z.Kind == shared.StudioMarkupKindSuggestion && !kontrolaTekstNiepusty(z.SuggestedText) {
		return shared.StudioMarkupAddResponse{}, bladWskazaniaStudio(
			"propozycja bez proponowanego brzmienia — propozycja bez brzmienia jest " +
				"komentarzem i zakłada się ją przez studio.comment.add")
	}

	dlugosc := postacDlugosc(&stan.forma)
	od, do := z.RangeStart, z.RangeEnd
	if od < 0 {
		od = 0
	}
	if do > dlugosc {
		do = dlugosc
	}
	zastane := znakowanieTekstZakresu(&stan.forma, od, do)

	wiersz := dane.ZnakowanieStudia{
		Kod:                  nowyIdentyfikator(przedrostekZnakowaniaStudia),
		Rodzaj:               string(z.Kind),
		AutorRodzaj:          string(wykonawca.Rodzaj),
		AutorAgentKod:        wykonawca.AgentKod,
		AutorAgentNazwa:      wykonawca.AgentNazwa,
		AutorAgentWersja:     wykonawca.AgentWersja,
		AutorPodagentKod:     wykonawca.PodagentKod,
		ZakresOd:             int64(od),
		ZakresDo:             int64(do),
		Barwa:                z.Color,
		ZnacznikNazwa:        z.MarkType,
		Tresc:                z.Body,
		BrzmienieProponowane: z.SuggestedText,
		Stan:                 string(shared.StudioMarkupStateOpen),
	}
	if zastane != "" {
		wiersz.BrzmienieZastane = &zastane
	}

	odpowiedz := shared.StudioMarkupAddResponse{}

	// Wyróżnienie barwą jest cechą postaci — nakłada się na runy drzewa, tą samą
	// drogą co reszta formatowania. Blokada wiążąca tę rękę zatrzymuje wyłącznie
	// nałożenie barwy; sam wiersz znakowania stoi POZA treścią i blokady nie
	// narusza, więc odmowa całości byłaby nieproporcjonalna.
	if z.Kind == shared.StudioMarkupKindHighlight {
		if !kontrolaTekstNiepusty(z.Color) {
			return shared.StudioMarkupAddResponse{}, bladWskazaniaStudio(
				"wyróżnienie bez barwy — wyróżnienie bez barwy nie jest widoczne")
		}
		odcinki, pominiete := postacOdcinkiDozwolone(&stan.forma, od, do, wykonawca.Rodzaj)
		nalozone := 0
		for _, odcinek := range odcinki {
			postacRozetnij(&stan.forma, odcinek[0], odcinek[1])
			for _, wskazanie := range postacFragmentyZakresu(&stan.forma, odcinek[0], odcinek[1]) {
				run := &stan.forma.Blocks[wskazanie[0]].Runs[wskazanie[1]]
				run.Format = postacScalZnak(run.Format,
					shared.StudioCharacterFormat{HighlightColor: z.Color})
				nalozone++
			}
		}
		if nalozone == 0 {
			return shared.StudioMarkupAddResponse{}, znakowanieBladBlokady(pominiete, od, do)
		}
		if err := a.postacZapisz(ctx, stan); err != nil {
			return shared.StudioMarkupAddResponse{}, err
		}
		forma := stan.forma
		odpowiedz.Form = &forma
	}

	zapisane, err := skladnica.ZapiszZnakowanie(ctx, stan.dokument.ID, wiersz)
	if err != nil {
		return shared.StudioMarkupAddResponse{}, bladStudio(err)
	}

	drzewo, err := dziennikZapisDrzewa(stan.forma)
	if err != nil {
		return shared.StudioMarkupAddResponse{}, err
	}
	czynnosc, err := a.dziennikOdlozCzynnosc(ctx, stan.dokument, wykonawca,
		shared.StudioActionKindMarkupChange,
		znakowanieOpisCzynnosci(z.Kind, wykonawca, od, do), &od, &do, drzewo, drzewo, nil)
	if err != nil {
		return shared.StudioMarkupAddResponse{}, err
	}
	odpowiedz.Markup = znakowanieZlozKontrakt(zapisane)
	odpowiedz.ActionId = czynnosc
	return odpowiedz, nil
}

// Znakowania obsługuje `studio.markup.list`.
//
// Wykaz jest SPISEM DO PRZEJŚCIA: idzie w kolejności wystąpienia w treści, żeby
// Operator skakał po nim od góry dokumentu do dołu i odhaczał pozycje. Bez tego
// znakowanie w długim dokumencie ginie.
func (a *adapterStudia) Znakowania(ctx context.Context,
	z shared.StudioMarkupListRequest) (shared.StudioMarkupListResponse, error) {

	dokument, err := a.dokumentDoCzynnosci(ctx, z.DocumentId)
	if err != nil {
		return shared.StudioMarkupListResponse{}, err
	}
	skladnica, err := a.kontrolaSkladnica()
	if err != nil {
		return shared.StudioMarkupListResponse{}, err
	}
	wiersze, err := skladnica.Znakowania(ctx, dokument.ID)
	if err != nil {
		return shared.StudioMarkupListResponse{}, bladStudio(err)
	}

	odpowiedz := shared.StudioMarkupListResponse{Markups: []shared.StudioMarkup{}}
	for _, wiersz := range wiersze {
		if z.Kind != nil && *z.Kind != "" && wiersz.Rodzaj != string(*z.Kind) {
			continue
		}
		if z.Author != nil && *z.Author != "" && wiersz.AutorRodzaj != string(*z.Author) {
			continue
		}
		if z.State != nil && *z.State != "" && wiersz.Stan != string(*z.State) {
			continue
		}
		if kontrolaTekstNiepusty(z.MarkType) &&
			wartoscTekstu(wiersz.ZnacznikNazwa) != *z.MarkType {

			continue
		}
		if kontrolaTekstNiepusty(z.AgentId) && wartoscTekstu(wiersz.AutorAgentKod) != *z.AgentId {
			continue
		}
		if kontrolaTekstNiepusty(z.SubagentId) &&
			wartoscTekstu(wiersz.AutorPodagentKod) != *z.SubagentId {

			continue
		}
		odpowiedz.Markups = append(odpowiedz.Markups, znakowanieZlozKontrakt(wiersz))
		if wiersz.Stan == string(shared.StudioMarkupStateOpen) {
			odpowiedz.OpenCount++
		}
	}
	odpowiedz.Total = len(odpowiedz.Markups)

	// Komentarze redakcyjne wchodzą do tego samego spisu — Operator ma widzieć
	// WSZYSTKO, czym dokument jest znaczony, w jednym wykazie. Kontrakt mówi
	// „brak znaczy tak", bo bez komentarzy spis nie byłby spisem znakowań.
	if z.IncludeComments == nil || *z.IncludeComments {
		komentarze, err := a.repozytorium.Komentarze(ctx, dokument.ID, rodzajKomentarza)
		if err != nil {
			return shared.StudioMarkupListResponse{}, bladStudio(err)
		}
		odpowiedz.Comments = []shared.StudioComment{}
		for _, wiersz := range komentarze {
			if z.Author != nil && *z.Author != "" && wiersz.Autor != string(*z.Author) {
				continue
			}
			odpowiedz.Comments = append(odpowiedz.Comments, złóżKomentarz(wiersz))
			odpowiedz.Total++
			if !wiersz.Rozwiazany {
				odpowiedz.OpenCount++
			}
		}
	}
	if z.IncludeTrackedChanges != nil && *z.IncludeTrackedChanges {
		zmiany, err := a.repozytorium.ZmianySledzone(ctx, dokument.ID)
		if err != nil {
			return shared.StudioMarkupListResponse{}, bladStudio(err)
		}
		odpowiedz.TrackedChanges = []shared.StudioTrackedChange{}
		for _, wiersz := range zmiany {
			if z.Author != nil && *z.Author != "" && wiersz.Autor != string(*z.Author) {
				continue
			}
			odpowiedz.TrackedChanges = append(odpowiedz.TrackedChanges, złóżZmianeSledzona(wiersz))
			odpowiedz.Total++
			if wiersz.Decyzja == "oczekuje" {
				odpowiedz.OpenCount++
			}
		}
	}
	return odpowiedz, nil
}

// RozstrzygnijZnakowanie obsługuje `studio.markup.decide`.
//
// Przyjęcie propozycji WNOSI proponowane brzmienie do treści i odkłada zmianę
// śledzoną AUTORA PROPOZYCJI — nie Operatora, który ją przyjął. Inaczej
// przełącznik „pokaż wszystko, co zrobił model" przestałby pokazywać zmianę,
// którą model naprawdę wniósł, tylko dlatego, że Operator ją zaakceptował.
//
// Decyzja idzie od KOŃCA dokumentu: zakresy propozycji liczone są w treści
// sprzed decyzji, więc przyjęcie od początku przesuwałoby zakresy propozycji
// jeszcze nierozpatrzonych.
func (a *adapterStudia) RozstrzygnijZnakowanie(ctx context.Context,
	z shared.StudioMarkupDecideRequest) (shared.StudioMarkupDecideResponse, error) {

	stan, err := a.postacWczytaj(ctx, z.DocumentId)
	if err != nil {
		return shared.StudioMarkupDecideResponse{}, err
	}
	skladnica, err := a.kontrolaSkladnica()
	if err != nil {
		return shared.StudioMarkupDecideResponse{}, err
	}
	if len(z.MarkupIds) == 0 {
		return shared.StudioMarkupDecideResponse{},
			bladWskazaniaStudio("decyzja bez wskazania znakowania")
	}
	if z.EditedText != nil && len(z.MarkupIds) > 1 {
		return shared.StudioMarkupDecideResponse{}, bladWskazaniaStudio(
			"poprawione brzmienie podano dla " + strconv.Itoa(len(z.MarkupIds)) +
				" propozycji naraz — jedno brzmienie nie może zastąpić kilku różnych " +
				"fragmentów; rozstrzygnij je po kolei")
	}

	wybrane := make([]dane.ZnakowanieStudia, 0, len(z.MarkupIds))
	for _, kod := range z.MarkupIds {
		wiersz, err := skladnica.Znakowanie(ctx, strings.TrimSpace(kod))
		if err != nil {
			return shared.StudioMarkupDecideResponse{}, dziennikBrakWiersza(err,
				"znakowanie nie istnieje: "+kod)
		}
		if wiersz.DokumentKod != stan.dokument.Kod {
			return shared.StudioMarkupDecideResponse{}, bladWskazaniaStudio("znakowanie " +
				kod + " stoi przy dokumencie " + wiersz.DokumentKod + ", nie przy " +
				stan.dokument.Kod)
		}
		if wiersz.Stan != string(shared.StudioMarkupStateOpen) {
			return shared.StudioMarkupDecideResponse{}, bladWskazaniaStudio("znakowanie " +
				kod + " jest już rozstrzygnięte („" + wiersz.Stan + "”) — druga decyzja " +
				"nie zmienia pierwszej")
		}
		wybrane = append(wybrane, wiersz)
	}
	sort.Slice(wybrane, func(i, j int) bool { return wybrane[i].ZakresOd > wybrane[j].ZakresOd })

	stanDocelowy := shared.StudioMarkupStateRejected
	if z.Accept {
		stanDocelowy = shared.StudioMarkupStateAccepted
	}

	bilans := shared.StudioActionBalance{Skipped: []shared.StudioSkippedItem{}}
	rozstrzygniete := 0
	for _, wiersz := range wybrane {
		if z.Accept && wiersz.Rodzaj == string(shared.StudioMarkupKindSuggestion) {
			brzmienie := wartoscTekstu(wiersz.BrzmienieProponowane)
			if z.EditedText != nil {
				brzmienie = *z.EditedText
			}
			od, do := int(wiersz.ZakresOd), int(wiersz.ZakresDo)
			przed := znakowanieTekstZakresu(&stan.forma, od, do)
			autor := shared.StudioAuthor(wiersz.AutorRodzaj)
			postacRozetnij(&stan.forma, od, do)
			postacZamienTresc(&stan.forma, od, do, brzmienie, nil, &autor)
			po := brzmienie
			if _, err := a.postacOdlozZmiane(ctx, stan, autor,
				shared.StudioChangeKindWstawienie, od, od+len([]rune(brzmienie)),
				&przed, &po); err != nil {

				return shared.StudioMarkupDecideResponse{}, err
			}
			bilans.Applied++
		}
		przestawione, err := skladnica.PrzestawStanZnakowania(ctx, wiersz.Kod,
			string(stanDocelowy))
		if err != nil {
			return shared.StudioMarkupDecideResponse{}, bladStudio(err)
		}
		if !przestawione {
			bilans.Skipped = append(bilans.Skipped, shared.StudioSkippedItem{
				Reason: "znakowanie rozstrzygnął w tym czasie ktoś inny",
				Detail: kontrolaWskaznikTekstu("znakowanie " + wiersz.Kod +
					" zmieniło stan między odczytem a zapisem"),
			})
			continue
		}
		rozstrzygniete++
	}
	if rozstrzygniete == 0 {
		return shared.StudioMarkupDecideResponse{}, protocol.JakoError(protocol.NowyBlad(
			shared.ErrorCodeConflict, "moduł Studio, znakowanie: nie rozstrzygnięto ani "+
				"jednej pozycji — wszystkie zmieniły stan w tym czasie"))
	}

	if err := a.postacZapisz(ctx, stan); err != nil {
		return shared.StudioMarkupDecideResponse{}, err
	}
	dokument := stan.dokument
	if bilans.Applied > 0 {
		drzewo, err := dziennikZapisDrzewa(stan.forma)
		if err != nil {
			return shared.StudioMarkupDecideResponse{}, err
		}
		if _, err := a.dziennikOdlozCzynnosc(ctx, dokument,
			kontrolaWykonawca{Rodzaj: shared.StudioAuthorUzytkownik},
			shared.StudioActionKindProposalAccept,
			"przyjęcie "+strconv.Itoa(bilans.Applied)+" propozycji z marginesu",
			nil, nil, drzewo, drzewo, nil); err != nil {

			return shared.StudioMarkupDecideResponse{}, err
		}
	}
	if z.CreateVersion == nil || *z.CreateVersion {
		if bilans.Applied > 0 {
			dokument, err = a.zalozWersjeDokumentu(ctx, dokument,
				wartoscTekstu(dokument.Tresc), shared.StudioAuthorUzytkownik, nil)
			if err != nil {
				return shared.StudioMarkupDecideResponse{}, err
			}
		}
	}
	bilans.SkippedCount = len(bilans.Skipped)
	bilans.Note = kontrolaWskaznikTekstu("Rozstrzygnięto " + strconv.Itoa(rozstrzygniete) +
		" pozycji znakowania; do treści weszło " + strconv.Itoa(bilans.Applied) +
		" proponowanych brzmień.")
	return shared.StudioMarkupDecideResponse{
		Decided:  rozstrzygniete,
		Document: a.zlozDokument(dokument),
		Form:     stan.forma,
		Balance:  bilans,
	}, nil
}

// ZdejmijZnakowanie obsługuje `studio.markup.remove`.
//
// Zdjęcie wyróżnienia zdejmuje też BARWĘ z runów — inaczej wiersz znikałby
// z wykazu, a fragment zostawałby podświetlony i nikt nie wiedziałby, czym.
func (a *adapterStudia) ZdejmijZnakowanie(ctx context.Context,
	z shared.StudioMarkupRemoveRequest) (shared.StudioMarkupRemoveResponse, error) {

	stan, err := a.postacWczytaj(ctx, z.DocumentId)
	if err != nil {
		return shared.StudioMarkupRemoveResponse{}, err
	}
	skladnica, err := a.kontrolaSkladnica()
	if err != nil {
		return shared.StudioMarkupRemoveResponse{}, err
	}
	if strings.TrimSpace(z.MarkupId) == "" {
		return shared.StudioMarkupRemoveResponse{},
			bladWskazaniaStudio("zdjęcie znakowania bez wskazania znakowania")
	}
	wiersz, err := skladnica.Znakowanie(ctx, strings.TrimSpace(z.MarkupId))
	if err != nil {
		return shared.StudioMarkupRemoveResponse{}, dziennikBrakWiersza(err,
			"znakowanie nie istnieje: "+z.MarkupId)
	}
	if wiersz.DokumentKod != stan.dokument.Kod {
		return shared.StudioMarkupRemoveResponse{}, bladWskazaniaStudio("znakowanie " +
			z.MarkupId + " stoi przy dokumencie " + wiersz.DokumentKod + ", nie przy " +
			stan.dokument.Kod)
	}
	wykonawca := kontrolaRozpoznajWykonawce(ctx, kontrolaPodpisZadania{
		Author: z.Author, AgentId: z.AgentId, AgentName: z.AgentName, SubagentId: z.SubagentId,
	})

	odpowiedz := shared.StudioMarkupRemoveResponse{}
	if wiersz.Rodzaj == string(shared.StudioMarkupKindHighlight) {
		od, do := int(wiersz.ZakresOd), int(wiersz.ZakresDo)
		odcinki, _ := postacOdcinkiDozwolone(&stan.forma, od, do, wykonawca.Rodzaj)
		for _, odcinek := range odcinki {
			postacRozetnij(&stan.forma, odcinek[0], odcinek[1])
			for _, wskazanie := range postacFragmentyZakresu(&stan.forma, odcinek[0], odcinek[1]) {
				run := &stan.forma.Blocks[wskazanie[0]].Runs[wskazanie[1]]
				if run.Format != nil {
					run.Format.HighlightColor = nil
				}
			}
		}
		if err := a.postacZapisz(ctx, stan); err != nil {
			return shared.StudioMarkupRemoveResponse{}, err
		}
		forma := stan.forma
		odpowiedz.Form = &forma
	}
	zdjete, err := skladnica.UsunZnakowanie(ctx, wiersz.Kod)
	if err != nil {
		return shared.StudioMarkupRemoveResponse{}, bladStudio(err)
	}
	odpowiedz.Removed = zdjete
	return odpowiedz, nil
}

// RodzajeZnacznika obsługuje `studio.markup.type.list`.
func (a *adapterStudia) RodzajeZnacznika(ctx context.Context,
	z shared.StudioMarkupTypeListRequest) (shared.StudioMarkupTypeListResponse, error) {

	skladnica, err := a.kontrolaSkladnica()
	if err != nil {
		return shared.StudioMarkupTypeListResponse{}, err
	}
	// Dokument zerowy znaczy „nie zawężaj licznika do jednego dokumentu"; wtedy
	// licznik jest liczbą użyć rodzaju we wszystkich dokumentach.
	var dokumentID int64
	if kontrolaTekstNiepusty(z.DocumentId) {
		dokument, err := a.dokumentDoCzynnosci(ctx, *z.DocumentId)
		if err != nil {
			return shared.StudioMarkupTypeListResponse{}, err
		}
		dokumentID = dokument.ID
	}
	wiersze, err := skladnica.RodzajeZnacznika(ctx, dokumentID)
	if err != nil {
		return shared.StudioMarkupTypeListResponse{}, bladStudio(err)
	}
	rodzaje := make([]shared.StudioMarkupType, 0, len(wiersze))
	for _, wiersz := range wiersze {
		rodzaje = append(rodzaje, znakowanieZlozRodzaj(wiersz))
	}
	return shared.StudioMarkupTypeListResponse{MarkupTypes: rodzaje}, nil
}

// ZapiszRodzajZnacznika obsługuje `studio.markup.type.save`.
//
// Rodzaju FABRYCZNEGO nie przepisuje: nazwa fabryczna jest tym samym pojęciem
// w każdym dokumencie i przestawienie jej barwy albo etykiety u jednego
// Operatora zmieniłoby ją wszystkim.
func (a *adapterStudia) ZapiszRodzajZnacznika(ctx context.Context,
	z shared.StudioMarkupTypeSaveRequest) (shared.StudioMarkupTypeSaveResponse, error) {

	skladnica, err := a.kontrolaSkladnica()
	if err != nil {
		return shared.StudioMarkupTypeSaveResponse{}, err
	}
	nazwa := strings.TrimSpace(z.Name)
	if nazwa == "" {
		return shared.StudioMarkupTypeSaveResponse{},
			bladWskazaniaStudio("rodzaj znacznika bez nazwy")
	}
	zastany, err := skladnica.RodzajZnacznika(ctx, nazwa)
	switch {
	case err == nil && zastany.Fabryczny:
		return shared.StudioMarkupTypeSaveResponse{}, protocol.JakoError(protocol.NowyBlad(
			shared.ErrorCodePermissionDenied, "moduł Studio: rodzaj znacznika „"+nazwa+
				"” jest fabryczny — nazwa fabryczna znaczy to samo w każdym dokumencie, "+
				"więc jej barwy ani etykiety nie przestawia się miejscowo. Załóż rodzaj "+
				"własny o innej nazwie."))
	case err != nil && !errors.Is(err, dane.ErrBrakWiersza):
		return shared.StudioMarkupTypeSaveResponse{}, bladStudio(err)
	}

	zapisany, err := skladnica.ZapiszRodzajZnacznika(ctx, dane.RodzajZnacznikaStudia{
		Nazwa:         nazwa,
		NazwaWidoczna: wartoscTekstu(z.Label),
		Barwa:         z.Color,
	})
	if err != nil {
		return shared.StudioMarkupTypeSaveResponse{}, bladStudio(err)
	}
	return shared.StudioMarkupTypeSaveResponse{MarkupType: znakowanieZlozRodzaj(zapisany)}, nil
}

// UsunRodzajZnacznika obsługuje `studio.markup.type.delete`.
//
// Rodzaju fabrycznego NIE usuwa — odpowiada odmową nazywającą powód, wzorem
// `studio.operation.delete`. Rodzaju będącego w użyciu też nie: usunięcie
// zostawiłoby znakowania wskazujące na rodzaj, którego nie ma, a Operator
// przestałby je odfiltrować.
func (a *adapterStudia) UsunRodzajZnacznika(ctx context.Context,
	z shared.StudioMarkupTypeDeleteRequest) (shared.StudioMarkupTypeDeleteResponse, error) {

	skladnica, err := a.kontrolaSkladnica()
	if err != nil {
		return shared.StudioMarkupTypeDeleteResponse{}, err
	}
	nazwa := strings.TrimSpace(z.Name)
	if nazwa == "" {
		return shared.StudioMarkupTypeDeleteResponse{},
			bladWskazaniaStudio("usunięcie rodzaju znacznika bez nazwy")
	}
	wiersze, err := skladnica.RodzajeZnacznika(ctx, 0)
	if err != nil {
		return shared.StudioMarkupTypeDeleteResponse{}, bladStudio(err)
	}
	for _, wiersz := range wiersze {
		if wiersz.Nazwa != nazwa {
			continue
		}
		if wiersz.Fabryczny {
			return shared.StudioMarkupTypeDeleteResponse{}, protocol.JakoError(protocol.NowyBlad(
				shared.ErrorCodePermissionDenied, "moduł Studio: rodzaj znacznika „"+nazwa+
					"” jest fabryczny i nie da się go usunąć — fabryczne stoją w produkcie "+
					"na stałe, tak samo jak operacje fabryczne"))
		}
		if wiersz.IleUzyc > 0 {
			return shared.StudioMarkupTypeDeleteResponse{}, protocol.JakoError(protocol.NowyBlad(
				shared.ErrorCodeConflict, "moduł Studio: rodzaju znacznika „"+nazwa+
					"” używa "+strconv.FormatInt(wiersz.IleUzyc, 10)+" znakowań — usunięcie "+
					"zostawiłoby je wskazujące na rodzaj, którego nie ma. Zdejmij je "+
					"najpierw przez studio.markup.remove."))
		}
	}
	usuniety, err := skladnica.UsunRodzajZnacznika(ctx, nazwa)
	if err != nil {
		return shared.StudioMarkupTypeDeleteResponse{}, bladStudio(err)
	}
	if !usuniety {
		return shared.StudioMarkupTypeDeleteResponse{},
			bladBrakuStudio("wykaz rodzajów znaczników nie zna rodzaju „" + nazwa + "”")
	}
	return shared.StudioMarkupTypeDeleteResponse{Deleted: true}, nil
}

// ── Wspólne ─────────────────────────────────────────────────────────────────

// znakowanieSprawdzRodzaj odbija rodzaj spoza wyliczenia kontraktu — tabela
// z migracji 365 i tak go nie przyjmie, a odmowa nazwana tutaj mówi, co wolno.
func znakowanieSprawdzRodzaj(rodzaj shared.StudioMarkupKind) error {
	for _, dozwolony := range shared.WartosciStudioMarkupKind() {
		if rodzaj == dozwolony {
			return nil
		}
	}
	nazwy := make([]string, 0, 3)
	for _, dozwolony := range shared.WartosciStudioMarkupKind() {
		nazwy = append(nazwy, string(dozwolony))
	}
	return bladWskazaniaStudio("znakowanie o nieznanym rodzaju „" + string(rodzaj) +
		"” — wolno: " + strings.Join(nazwy, ", "))
}

// znakowanieTekstZakresu wycina brzmienie zastane fragmentu. Liczy w ZNAKACH,
// tak jak nazywa zakresy kontrakt.
func znakowanieTekstZakresu(forma *shared.StudioDocumentForm, od, do int) string {
	znaki := []rune(postacTekstFormy(forma))
	od, do, poprawny := kontrolaZakresWTresci(znaki, od, do)
	if !poprawny {
		return ""
	}
	return string(znaki[od:do])
}

// znakowanieOpisCzynnosci składa zdanie dziennika — Operator ma wiedzieć, co
// cofa, a nie odczytywać to z rodzaju wpisu.
func znakowanieOpisCzynnosci(rodzaj shared.StudioMarkupKind, wykonawca kontrolaWykonawca,
	od, do int) string {

	czynnosc := "znakowanie fragmentu"
	switch rodzaj {
	case shared.StudioMarkupKindHighlight:
		czynnosc = "wyróżnienie fragmentu barwą"
	case shared.StudioMarkupKindMark:
		czynnosc = "znacznik na fragmencie"
	case shared.StudioMarkupKindSuggestion:
		czynnosc = "propozycja brzmienia na marginesie"
	}
	return czynnosc + " od znaku " + strconv.Itoa(od) + " do " + strconv.Itoa(do) +
		"; założył: " + wykonawca.nazwaWykonawcy()
}

// znakowanieBladBlokady jest odmową nazywającą blokadę, która zatrzymała
// nałożenie wyróżnienia na cały wskazany fragment.
func znakowanieBladBlokady(pominiete []shared.StudioSkippedItem, od, do int) error {
	return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodePermissionDenied,
		"moduł Studio: wyróżnienia nie da się nałożyć na fragment od znaku "+
			strconv.Itoa(od)+" do "+strconv.Itoa(do)+" — cały leży pod blokadą "+
			postacNazwaBlokad(pominiete)+
			". Fragment wymagający uwagi wskazuje się propozycją na marginesie "+
			"(studio.markup.add o rodzaju suggestion)."))
}

// znakowanieZlozKontrakt składa znakowanie kontraktu z wiersza warstwy danych.
func znakowanieZlozKontrakt(wiersz dane.ZnakowanieStudia) shared.StudioMarkup {
	return shared.StudioMarkup{
		Id:                 wiersz.Kod,
		DocumentId:         wiersz.DokumentKod,
		Kind:               shared.StudioMarkupKind(wiersz.Rodzaj),
		Author:             shared.StudioAuthor(wiersz.AutorRodzaj),
		RangeStart:         int(wiersz.ZakresOd),
		RangeEnd:           int(wiersz.ZakresDo),
		Color:              wiersz.Barwa,
		MarkType:           wiersz.ZnacznikNazwa,
		Body:               wiersz.Tresc,
		SuggestedText:      wiersz.BrzmienieProponowane,
		OriginalText:       wiersz.BrzmienieZastane,
		State:              shared.StudioMarkupState(wiersz.Stan),
		CreatedAt:          chwilaBazy(wiersz.Utworzono),
		AuthorAgentId:      wiersz.AutorAgentKod,
		AuthorAgentName:    wiersz.AutorAgentNazwa,
		AuthorAgentVersion: wiersz.AutorAgentWersja,
		AuthorSubagentId:   wiersz.AutorPodagentKod,
	}
}

// znakowanieZlozRodzaj składa rodzaj znacznika kontraktu z wiersza.
func znakowanieZlozRodzaj(wiersz dane.RodzajZnacznikaStudia) shared.StudioMarkupType {
	ileUzyc := int(wiersz.IleUzyc)
	return shared.StudioMarkupType{
		Name:       wiersz.Nazwa,
		Label:      wiersz.NazwaWidoczna,
		Color:      wiersz.Barwa,
		Builtin:    wskaznikLogiczny(wiersz.Fabryczny),
		UsageCount: &ileUzyc,
	}
}
