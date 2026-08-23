// Odpowiedzialność pliku: blokady fragmentów dokumentu — założenie, zdjęcie,
// wykaz — oraz RACHUNEK UZGODNIENIA, którym rdzeń wykonuje zmianę POZA blokadą
// i oddaje bilans tego, co pominął.
//
// ── Blokada obowiązuje w rdzeniu, nie w oknie ────────────────────────────────
// Wymaganie rozstrzygające Właściciela. Model woła komendy rdzenia tak samo jak
// klient, więc blokada pilnowana wyłącznie przez okno jest pozorna — ominąłby ją
// bez wysiłku. Sprawdzenie stoi więc na drodze KAŻDEJ komendy zmieniającej
// dokument (`adapter_modul_studio_blokady_zapora.go`), a nie w poszczególnych
// obsługiwaczach: obsługiwaczy zmieniających dokument jest w Studiu ponad
// setka i pisanych przez czterech wykonawców, a zapora wpięta w rejestr obejmuje
// też te, których jeszcze nie ma.
//
// ── Trzy odpowiedzi, nie dwie ───────────────────────────────────────────────
// Zderzenie zmiany z blokadą nie ma odpowiedzi „wolno / nie wolno", ma trzy:
//
//	CAŁOŚĆ W BLOKADZIE  → odmowa NAZWANA: który fragment i jaka blokada. Cicha
//	                      bezczynność byłaby najgorszą możliwą odpowiedzią, bo
//	                      Operator myślałby, że model wykonał polecenie.
//	CZĘŚĆ W BLOKADZIE   → zmiana wchodzi POZA blokadą, a odpowiedź niesie BILANS:
//	                      co zmienione, co pominięte i przez którą blokadę.
//	                      Odmowa całości byłaby tu nieproporcjonalna (zamiana
//	                      w całym dokumencie trafiająca w jeden cytat), a
//	                      przemilczenie pominięcia — zakazane.
//	POZA BLOKADĄ        → zmiana wchodzi bez słowa.
//
// ── Kogo blokada wiąże ──────────────────────────────────────────────────────
// Blokada jest skierowana przeciw wykonawcom, nie przeciw właścicielowi
// dokumentu: Operator zmienia fragment zablokowany BEZ przeszkód. Blokada
// działająca także na Operatora jest osobnym, jawnym ustawieniem (`scope`
// = `everyone`), a nie zachowaniem domyślnym.
//
// Zdejmuje blokadę WYŁĄCZNIE Operator. Wykonawca, który uzna, że fragment
// wymaga zmiany, zakłada propozycję na marginesie (`studio.markup.add` o
// rodzaju `suggestion`) — i tyle. Dlatego `studio.lock.remove` nie stoi też
// w wykazie narzędzi modelu; sprawdzenie tutaj jest drugą stroną tej samej
// zasady, bo wykaz narzędzi nie jest zaporą.
//
// ── Blokada przechodzi przez wersje i przez szablon ─────────────────────────
// Blokada wisi przy DOKUMENCIE, nie przy wersji, więc przywrócenie wcześniejszej
// wersji jej nie gubi — i to jest cały mechanizm, bez ani jednej linii kodu
// w drodze przywracania. Zakresy nadążają za treścią przez
// `PrzesunBlokadyFragmentow`. Szablon niesie swoje blokady wzorcowe do
// dokumentów z niego zakładanych (`blokadaPrzenieSzablon`).
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

// ZalozBlokade obsługuje `studio.lock.add`.
//
// Blokadę zakłada Operator; wykonawca, który po nią sięgnie, dostaje odmowę
// nazywającą powód. Wykonawca zakładający blokadę mógłby unieruchomić fragment
// przed Operatorem — a blokada jest narzędziem Operatora przeciw wykonawcom,
// nie odwrotnie.
func (a *adapterStudia) ZalozBlokade(ctx context.Context,
	z shared.StudioLockAddRequest) (shared.StudioLockAddResponse, error) {

	dokument, err := a.dokumentDoCzynnosci(ctx, z.DocumentId)
	if err != nil {
		return shared.StudioLockAddResponse{}, err
	}
	skladnica, err := a.kontrolaSkladnica()
	if err != nil {
		return shared.StudioLockAddResponse{}, err
	}
	if strings.TrimSpace(z.Name) == "" {
		return shared.StudioLockAddResponse{},
			bladWskazaniaStudio("blokada bez nazwy — Operator nie odróżniłby jej od innych w wykazie")
	}
	if z.RangeEnd < z.RangeStart {
		return shared.StudioLockAddResponse{},
			bladWskazaniaStudio("blokada o zakresie odwróconym: koniec " +
				strconv.Itoa(z.RangeEnd) + " leży przed początkiem " + strconv.Itoa(z.RangeStart))
	}
	wykonawca := kontrolaRozpoznajWykonawce(ctx, kontrolaPodpisZadania{})
	if wykonawca.czyWykonawca() {
		return shared.StudioLockAddResponse{}, protocol.JakoError(protocol.NowyBlad(
			shared.ErrorCodePermissionDenied,
			"moduł Studio: blokadę fragmentu zakłada wyłącznie Operator — "+
				wykonawca.nazwaWykonawcy()+" nie unieruchamia fragmentów dokumentu. "+
				"Fragment wymagający uwagi wskazuje się propozycją na marginesie."))
	}

	zasieg := string(shared.StudioLockScopeModel)
	if z.Scope != nil && *z.Scope != "" {
		zasieg = string(*z.Scope)
	}
	zapisana, err := skladnica.ZapiszBlokadeFragmentu(ctx, dokument.ID, dane.BlokadaFragmentuStudia{
		Kod:           nowyIdentyfikator(przedrostekBlokadyStudia),
		Nazwa:         strings.TrimSpace(z.Name),
		Powod:         z.Reason,
		ZakresOd:      int64(z.RangeStart),
		ZakresDo:      int64(z.RangeEnd),
		Zasieg:        zasieg,
		ZalozylRodzaj: string(shared.StudioAuthorUzytkownik),
	})
	if err != nil {
		return shared.StudioLockAddResponse{}, bladStudio(err)
	}
	wszystkie, err := a.blokadyDokumentu(ctx, dokument.ID)
	if err != nil {
		return shared.StudioLockAddResponse{}, err
	}
	return shared.StudioLockAddResponse{
		Lock: blokadaZlozKontrakt(zapisana), Locks: wszystkie,
	}, nil
}

// ZdejmijBlokade obsługuje `studio.lock.remove`.
func (a *adapterStudia) ZdejmijBlokade(ctx context.Context,
	z shared.StudioLockRemoveRequest) (shared.StudioLockRemoveResponse, error) {

	dokument, err := a.dokumentDoCzynnosci(ctx, z.DocumentId)
	if err != nil {
		return shared.StudioLockRemoveResponse{}, err
	}
	skladnica, err := a.kontrolaSkladnica()
	if err != nil {
		return shared.StudioLockRemoveResponse{}, err
	}
	if z.LockId == "" {
		return shared.StudioLockRemoveResponse{},
			bladWskazaniaStudio("zdjęcie blokady bez wskazania blokady")
	}

	wykonawca := kontrolaRozpoznajWykonawce(ctx, kontrolaPodpisZadania{
		Author: z.Author, AgentId: z.AgentId, AgentName: z.AgentName, SubagentId: z.SubagentId,
	})
	if wykonawca.czyWykonawca() {
		return shared.StudioLockRemoveResponse{}, protocol.JakoError(protocol.NowyBlad(
			shared.ErrorCodePermissionDenied,
			"moduł Studio: blokadę zdejmuje WYŁĄCZNIE Operator — "+
				wykonawca.nazwaWykonawcy()+" jej nie zdejmie i nie prosi o zdjęcie obejściem. "+
				"Fragment wymagający zmiany wskazuje się propozycją na marginesie "+
				"(studio.markup.add o rodzaju suggestion)."))
	}

	blokada, err := skladnica.BlokadaFragmentu(ctx, z.LockId)
	if err != nil {
		if errors.Is(err, dane.ErrBrakWiersza) {
			return shared.StudioLockRemoveResponse{},
				bladBrakuStudio("blokada nie istnieje: " + z.LockId)
		}
		return shared.StudioLockRemoveResponse{}, bladStudio(err)
	}
	// Blokada cudzego dokumentu nie zdejmuje się przez wskazanie tego dokumentu:
	// odpowiedź „zdjęto" na blokadę, która stała przy innym dokumencie, byłaby
	// nieprawdą o obu.
	if blokada.DokumentKod != dokument.Kod {
		return shared.StudioLockRemoveResponse{},
			bladWskazaniaStudio("blokada " + z.LockId + " stoi przy dokumencie " +
				blokada.DokumentKod + ", nie przy " + dokument.Kod)
	}

	zdjeta, err := skladnica.UsunBlokadeFragmentu(ctx, z.LockId)
	if err != nil {
		return shared.StudioLockRemoveResponse{}, bladStudio(err)
	}
	wszystkie, err := a.blokadyDokumentu(ctx, dokument.ID)
	if err != nil {
		return shared.StudioLockRemoveResponse{}, err
	}
	return shared.StudioLockRemoveResponse{Removed: zdjeta, Locks: wszystkie}, nil
}

// Blokady obsługuje `studio.lock.list`.
func (a *adapterStudia) Blokady(ctx context.Context,
	z shared.StudioLockListRequest) (shared.StudioLockListResponse, error) {

	dokument, err := a.dokumentDoCzynnosci(ctx, z.DocumentId)
	if err != nil {
		return shared.StudioLockListResponse{}, err
	}
	wszystkie, err := a.blokadyDokumentu(ctx, dokument.ID)
	if err != nil {
		return shared.StudioLockListResponse{}, err
	}
	if z.RangeStart == nil && z.RangeEnd == nil {
		return shared.StudioLockListResponse{Locks: wszystkie}, nil
	}
	// Zawężenie do fragmentu oddaje blokady STYKAJĄCE SIĘ z nim, nie tylko
	// zawarte w nim w całości: Operator pytający „co stoi na tym zaznaczeniu"
	// pyta o wszystko, co go dotyka, także o blokadę wychodzącą poza zaznaczenie.
	od := 0
	if z.RangeStart != nil {
		od = *z.RangeStart
	}
	do := od
	if z.RangeEnd != nil {
		do = *z.RangeEnd
	}
	wybrane := []shared.StudioFragmentLock{}
	for _, blokada := range wszystkie {
		if kontrolaZakresyStykaja(od, do, blokada.RangeStart, blokada.RangeEnd) {
			wybrane = append(wybrane, blokada)
		}
	}
	return shared.StudioLockListResponse{Locks: wybrane}, nil
}

// blokadyDokumentu składa wykaz blokad dokumentu w kształcie kontraktu.
func (a *adapterStudia) blokadyDokumentu(ctx context.Context,
	dokumentID int64) ([]shared.StudioFragmentLock, error) {

	skladnica, err := a.kontrolaSkladnica()
	if err != nil {
		return nil, err
	}
	wiersze, err := skladnica.BlokadyFragmentow(ctx, dokumentID)
	if err != nil {
		return nil, bladStudio(err)
	}
	blokady := make([]shared.StudioFragmentLock, 0, len(wiersze))
	for _, wiersz := range wiersze {
		blokady = append(blokady, blokadaZlozKontrakt(wiersz))
	}
	return blokady, nil
}

// blokadaZlozKontrakt składa blokadę kontraktu z wiersza warstwy danych.
func blokadaZlozKontrakt(wiersz dane.BlokadaFragmentuStudia) shared.StudioFragmentLock {
	blokada := shared.StudioFragmentLock{
		Id:         wiersz.Kod,
		DocumentId: wiersz.DokumentKod,
		Name:       wiersz.Nazwa,
		Reason:     wiersz.Powod,
		RangeStart: int(wiersz.ZakresOd),
		RangeEnd:   int(wiersz.ZakresDo),
		Scope:      shared.StudioLockScope(wiersz.Zasieg),
	}
	if wiersz.ZalozylRodzaj != "" {
		autor := shared.StudioAuthor(wiersz.ZalozylRodzaj)
		blokada.CreatedBy = &autor
	}
	if wiersz.ZSzablonu {
		blokada.FromTemplate = wskaznikLogiczny(true)
	}
	if chwila := chwilaBazy(wiersz.Utworzono); chwila != 0 {
		blokada.CreatedAt = &chwila
	}
	return blokada
}

// ── Rachunek uzgodnienia ────────────────────────────────────────────────────

// blokadaWiazaca mówi, czy blokada wiąże tę rękę.
//
// Zasięg `everyone` wiąże każdego, także Operatora — to jest to OSOBNE, JAWNE
// ustawienie. Zasięg `model` (postać domyślna) wiąże wyłącznie wykonawców.
func blokadaWiazaca(blokada dane.BlokadaFragmentuStudia, wykonawca kontrolaWykonawca) bool {
	if blokada.Zasieg == string(shared.StudioLockScopeEveryone) {
		return true
	}
	return wykonawca.czyWykonawca()
}

// blokadaWiazaceDlaRak przesiewa blokady dokumentu do tych, które wiążą tę rękę.
func blokadaWiazaceDlaRak(blokady []dane.BlokadaFragmentuStudia,
	wykonawca kontrolaWykonawca) []dane.BlokadaFragmentuStudia {

	wiazace := make([]dane.BlokadaFragmentuStudia, 0, len(blokady))
	for _, blokada := range blokady {
		if blokadaWiazaca(blokada, wykonawca) {
			wiazace = append(wiazace, blokada)
		}
	}
	return wiazace
}

// blokadaNaZakresie oddaje pierwszą blokadę wiążącą, która styka się z zakresem.
// Pierwszą, nie wszystkie: odmowa ma nazwać JEDNĄ blokadę, tę najbliższą
// początkowi fragmentu, bo wykaz trzech blokad w treści błędu nie pomaga
// bardziej niż jedna, a wykaz pełny Operator ma w `studio.lock.list`.
func blokadaNaZakresie(blokady []dane.BlokadaFragmentuStudia, od, do int,
	wykonawca kontrolaWykonawca) (dane.BlokadaFragmentuStudia, bool) {

	for _, blokada := range blokady {
		if !blokadaWiazaca(blokada, wykonawca) {
			continue
		}
		if kontrolaZakresyStykaja(od, do, int(blokada.ZakresOd), int(blokada.ZakresDo)) {
			return blokada, true
		}
	}
	return dane.BlokadaFragmentuStudia{}, false
}

// blokadaUzgodnienie jest wynikiem rachunku uzgodnienia treści z blokadami.
type blokadaUzgodnienie struct {
	// Tresc jest treścią po uzgodnieniu: zmiana wniesiona poza blokadami,
	// fragmenty zablokowane przywrócone do brzmienia zastanego.
	Tresc string
	// Zmienione i Pominiete liczą WIERSZE, bo wierszami idzie uzgodnienie.
	Zmienione int
	Pominiete []shared.StudioSkippedItem
	// CalkiemWBlokadzie znaczy, że po uzgodnieniu nie zostało nic do zmiany —
	// czyli zmiana leżała w blokadzie w CAŁOŚCI. Wtedy odpowiedzią jest odmowa
	// nazwana, nie cicha zgoda na brak skutku.
	CalkiemWBlokadzie bool
	// PierwszaBlokada nazywa blokadę do treści odmowy.
	PierwszaBlokada dane.BlokadaFragmentuStudia
}

// blokadaUzgodnijTresc wykonuje zmianę POZA blokadami i liczy bilans pominięć.
//
// ── Dlaczego wierszami ──────────────────────────────────────────────────────
// Uzgodnienie musi rozdzielić „ten fragment zmiany wolno" od „tego nie wolno",
// a do tego trzeba miejsc, w których zmianę da się PRZERWAĆ. Wiersz jest takim
// miejscem i jest nim w tym module od początku: `policzFragmentyRoznicy` liczy
// różnicę wierszami, a `studio.diff.hunk.apply` wierszami przenosi fragmenty.
// Znak byłby miejscem gęstszym, ale uzgodnienie znakami wymagałoby dopasowania
// klasy LCS, którego rdzeń nie ma i którego ten odcinek nie zakłada drugi raz.
//
// ── Dlaczego dwa tryby ──────────────────────────────────────────────────────
// Gdy strony mają TYLE SAMO wierszy w środku (a tak wygląda zamiana w całym
// dokumencie — sedno wymagania o bilansie), wiersze dają się sparować jeden do
// jednego i uzgodnienie jest dokładne: wiersz trafiony blokadą zostaje zastany,
// pozostałe wchodzą.
//
// Gdy liczba wierszy się różni, sparować ich nie sposób bez zgadywania, KTÓRY
// wiersz odpowiada któremu. Wtedy środek jest jedną całością: styka się
// z blokadą — nie wchodzi w ogóle i wraca pominięciem; nie styka — wchodzi
// w całości. Zgadywanie parowania byłoby tu gorsze niż odmowa, bo wstawiłoby
// treść w środek zablokowanego cytatu i nazwało to bilansem.
func blokadaUzgodnijTresc(zastana, proponowana string,
	blokady []dane.BlokadaFragmentuStudia, wykonawca kontrolaWykonawca) blokadaUzgodnienie {

	if zastana == proponowana {
		return blokadaUzgodnienie{Tresc: proponowana}
	}
	wiazace := blokadaWiazaceDlaRak(blokady, wykonawca)
	if len(wiazace) == 0 {
		return blokadaUzgodnienie{Tresc: proponowana, Zmienione: 1}
	}

	przed := strings.Split(zastana, "\n")
	po := strings.Split(proponowana, "\n")

	przedrostek := 0
	for przedrostek < len(przed) && przedrostek < len(po) && przed[przedrostek] == po[przedrostek] {
		przedrostek++
	}
	sufiks := 0
	for sufiks < len(przed)-przedrostek && sufiks < len(po)-przedrostek &&
		przed[len(przed)-1-sufiks] == po[len(po)-1-sufiks] {
		sufiks++
	}
	srodekPrzed := przed[przedrostek : len(przed)-sufiks]
	srodekPo := po[przedrostek : len(po)-sufiks]

	// Początki wierszy liczone w ZNAKACH, bo w znakach liczy zakresy kontrakt
	// i w znakach stoją zakresy blokad.
	poczatki := blokadaPoczatkiWierszy(przed)

	uzgodnienie := blokadaUzgodnienie{}
	wynik := append([]string{}, przed[:przedrostek]...)

	if len(srodekPrzed) == len(srodekPo) {
		for i := range srodekPrzed {
			numer := przedrostek + i
			if srodekPrzed[i] == srodekPo[i] {
				wynik = append(wynik, srodekPo[i])
				continue
			}
			od := poczatki[numer]
			do := od + len([]rune(srodekPrzed[i]))
			blokada, trafiona := blokadaNaZakresie(wiazace, od, do, wykonawca)
			if !trafiona {
				wynik = append(wynik, srodekPo[i])
				uzgodnienie.Zmienione++
				continue
			}
			wynik = append(wynik, srodekPrzed[i])
			uzgodnienie.Pominiete = append(uzgodnienie.Pominiete,
				blokadaPominiecie(blokada, od, do))
			if uzgodnienie.PierwszaBlokada.Kod == "" {
				uzgodnienie.PierwszaBlokada = blokada
			}
		}
	} else {
		od := poczatki[przedrostek]
		do := od
		if len(srodekPrzed) > 0 {
			do = poczatki[przedrostek+len(srodekPrzed)-1] +
				len([]rune(srodekPrzed[len(srodekPrzed)-1]))
		}
		blokada, trafiona := blokadaNaZakresie(wiazace, od, do, wykonawca)
		if trafiona {
			wynik = append(wynik, srodekPrzed...)
			uzgodnienie.Pominiete = append(uzgodnienie.Pominiete,
				blokadaPominiecie(blokada, od, do))
			uzgodnienie.PierwszaBlokada = blokada
		} else {
			wynik = append(wynik, srodekPo...)
			uzgodnienie.Zmienione++
		}
	}

	wynik = append(wynik, przed[len(przed)-sufiks:]...)
	uzgodnienie.Tresc = strings.Join(wynik, "\n")
	uzgodnienie.CalkiemWBlokadzie = uzgodnienie.Zmienione == 0 && len(uzgodnienie.Pominiete) > 0
	return uzgodnienie
}

// blokadaPoczatkiWierszy liczy początek każdego wiersza w ZNAKACH. Znak nowego
// wiersza liczy się do długości wiersza poprzedniego, bo tak liczy je zakres
// zaznaczenia w oknie.
func blokadaPoczatkiWierszy(wiersze []string) []int {
	poczatki := make([]int, len(wiersze)+1)
	biezacy := 0
	for i, wiersz := range wiersze {
		poczatki[i] = biezacy
		biezacy += len([]rune(wiersz)) + 1
	}
	poczatki[len(wiersze)] = biezacy
	return poczatki
}

// blokadaPominiecie składa pozycję bilansu: co pominięte i przez którą blokadę.
func blokadaPominiecie(blokada dane.BlokadaFragmentuStudia, od, do int) shared.StudioSkippedItem {
	szczegol := "fragment od znaku " + strconv.Itoa(od) + " do " + strconv.Itoa(do) +
		" został pominięty — stoi na nim blokada „" + blokada.Nazwa + "”"
	if blokada.Powod != nil && *blokada.Powod != "" {
		szczegol += " (" + *blokada.Powod + ")"
	}
	nazwa := blokada.Nazwa
	kod := blokada.Kod
	return shared.StudioSkippedItem{
		Reason:     "blokada fragmentu",
		Detail:     &szczegol,
		LockId:     &kod,
		LockName:   &nazwa,
		RangeStart: &od,
		RangeEnd:   &do,
	}
}

// blokadaBilans składa bilans czynności z uzgodnienia.
func blokadaBilans(uzgodnienie blokadaUzgodnienie) shared.StudioActionBalance {
	bilans := shared.StudioActionBalance{
		Applied:      uzgodnienie.Zmienione,
		SkippedCount: len(uzgodnienie.Pominiete),
		Skipped:      uzgodnienie.Pominiete,
	}
	if len(uzgodnienie.Pominiete) > 0 {
		zdanie := "Zmiana weszła poza blokadami; " + strconv.Itoa(len(uzgodnienie.Pominiete)) +
			" fragment(ów) pominięto, bo stoją na nich blokady Operatora."
		bilans.Note = &zdanie
	}
	return bilans
}

// ── Blokady wzorcowe szablonu ───────────────────────────────────────────────

// blokadaPrzenieSzablon kopiuje blokady wzorcowe szablonu do dokumentu z niego
// zakładanego — fragmenty wzorcowe pisma mają zostać wzorcowe.
//
// Metoda stoi tutaj i jest wołana przez odcinek szablonów; wołanie jej z drogi
// zakładania dokumentu należy do tamtego odcinka, bo tam stoi ta droga.
func (a *adapterStudia) blokadaPrzenieSzablon(ctx context.Context, dokumentID int64,
	szablonKod string) (int, error) {

	if szablonKod == "" {
		return 0, nil
	}
	skladnica, err := a.kontrolaSkladnica()
	if err != nil {
		return 0, err
	}
	wzorce, err := skladnica.BlokadySzablonu(ctx, szablonKod)
	if err != nil {
		return 0, bladStudio(err)
	}
	szablon := szablonKod
	przeniesione := 0
	for _, wzorzec := range wzorce {
		_, err := skladnica.ZapiszBlokadeFragmentu(ctx, dokumentID, dane.BlokadaFragmentuStudia{
			Kod:           nowyIdentyfikator(przedrostekBlokadyStudia),
			Nazwa:         wzorzec.Nazwa,
			Powod:         wzorzec.Powod,
			ZakresOd:      wzorzec.ZakresOd,
			ZakresDo:      wzorzec.ZakresDo,
			Zasieg:        wzorzec.Zasieg,
			ZalozylRodzaj: string(shared.StudioAuthorUzytkownik),
			ZSzablonu:     true,
			SzablonKod:    &szablon,
		})
		if err != nil {
			return przeniesione, bladStudio(err)
		}
		przeniesione++
	}
	return przeniesione, nil
}
