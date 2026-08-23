// Odpowiedzialność pliku: ODWRACALNY DZIENNIK CZYNNOŚCI dokumentu — wykaz
// (`studio.journal.list`), cofnięcie pojedynczej czynności także NIE PO KOLEI
// (`studio.journal.revert`) i ponowienie cofniętej (`studio.journal.redo`).
//
// ── Dlaczego wpis niesie całe drzewo postaci, a nie samą treść ───────────────
// Wymaganie Właściciela: „pomyłkowa zmiana kroju w całym dokumencie musi się
// cofać tak samo jak skasowany akapit". Taka zmiana nie rusza ANI JEDNEJ litery,
// więc treść sprzed czynności jej nie odtworzy. Dlatego wpis dziennika niesie
// `StanPrzed` i `StanPo` jako pełne drzewo postaci (`StudioDocumentForm`),
// a cofnięcie przywraca drzewo, nie napis.
//
// ── Dlaczego cofnięcie NIE jest przywróceniem wersji ────────────────────────
// Przywrócenie wersji cofa wszystko, co po niej weszło — także pracę Operatora.
// Cofnięcie czynności ze ŚRODKA dziennika ma zdjąć wyłącznie to, co zrobiła ta
// jedna czynność. Dlatego rdzeń nie wgrywa `StanPrzed` w miejsce dokumentu, a
// liczy RÓŻNICĘ między `StanPo` i `StanPrzed` i nakłada ją odwrotnie na stan
// BIEŻĄCY, byt po bycie (blok, styl, sekcja, obiekt, pole, nastawy strony).
// Byt, którego czynność nie tknęła, zostaje nietknięty.
//
// ── Dlaczego czynność z zależnością ODMAWIA, a nie cofa się po cichu ─────────
// Tabela zależności czynności (migracja 364) istnieje wyłącznie po to. Cofnięcie
// czynności, na której stoi późniejsza, zostawiłoby dokument w stanie
// niespójnym — na przykład zdjęcie tabeli, do której później wstawiono wiersze.
// Odmowa NAZYWA zależność: Operator ma wiedzieć, co cofnąć wcześniej, a nie
// zgadywać, dlaczego nic się nie stało.
//
// ── Dwa słowniki, których nie wolno pomieszać ───────────────────────────────
// Rodzaj czynności dziennika to `shared.StudioActionKind` (jedenaście wartości).
// Rodzaj zmiany śledzonej to `shared.StudioChangeKind` (trzy). Stan wpisu
// dziennika to WYŁĄCZNIE `active` albo `reverted` — tabela z migracji 364 nie
// przyjmie niczego innego.
package core

import (
	"context"
	"encoding/json"
	"errors"
	"sort"
	"strconv"
	"strings"

	"danacoconsole/server/internal/dane"
	"danacoconsole/server/internal/protocol"
	"danacoconsole/shared"
)

// DziennikCzynnosci obsługuje `studio.journal.list`.
func (a *adapterStudia) DziennikCzynnosci(ctx context.Context,
	z shared.StudioJournalListRequest) (shared.StudioJournalListResponse, error) {

	dokument, err := a.dokumentDoCzynnosci(ctx, z.DocumentId)
	if err != nil {
		return shared.StudioJournalListResponse{}, err
	}
	skladnica, err := a.kontrolaSkladnica()
	if err != nil {
		return shared.StudioJournalListResponse{}, err
	}
	wiersze, err := skladnica.CzynnosciDokumentu(ctx, dokument.ID)
	if err != nil {
		return shared.StudioJournalListResponse{}, bladStudio(err)
	}

	wybrane := make([]shared.StudioDocumentAction, 0, len(wiersze))
	for _, wiersz := range wiersze {
		if z.Author != nil && *z.Author != "" && wiersz.AutorRodzaj != string(*z.Author) {
			continue
		}
		if z.Kind != nil && *z.Kind != "" && wiersz.Rodzaj != string(*z.Kind) {
			continue
		}
		if z.State != nil && *z.State != "" && wiersz.Stan != string(*z.State) {
			continue
		}
		if kontrolaTekstNiepusty(z.AgentId) && wartoscTekstu(wiersz.AutorAgentKod) != *z.AgentId {
			continue
		}
		if kontrolaTekstNiepusty(z.SubagentId) &&
			wartoscTekstu(wiersz.AutorPodagentKod) != *z.SubagentId {

			continue
		}
		wybrane = append(wybrane, dziennikZlozCzynnosc(wiersz))
	}
	// Licznik mówi, ile wpisów dziennik niesie PO zawężeniu, a nie ile ich
	// wyszło po obcięciu granicą: Operator pytający o dziesięć ostatnich ma
	// wiedzieć, że stoi ich sto.
	wszystkich := len(wybrane)
	if z.Limit != nil && *z.Limit > 0 && *z.Limit < len(wybrane) {
		wybrane = wybrane[:*z.Limit]
	}
	return shared.StudioJournalListResponse{Actions: wybrane, Total: wszystkich}, nil
}

// CofnijCzynnosc obsługuje `studio.journal.revert`.
//
// Cofa POJEDYNCZO i nie po kolei. Kolejność wykonania jest odwrotna do
// kolejności dziennika: gdy Operator wskazał kilka czynności naraz, zdejmuje się
// najpierw najświeższą, bo starsza mogła być jej podstawą.
func (a *adapterStudia) CofnijCzynnosc(ctx context.Context,
	z shared.StudioJournalRevertRequest) (shared.StudioJournalRevertResponse, error) {

	stan, skladnica, wykaz, err := a.dziennikStanIWykaz(ctx, z.DocumentId)
	if err != nil {
		return shared.StudioJournalRevertResponse{}, err
	}
	if len(z.ActionIds) == 0 {
		return shared.StudioJournalRevertResponse{},
			bladWskazaniaStudio("cofnięcie bez wskazania czynności — dziennik cofa się " +
				"pojedynczo, a nie „ostatnim ruchem”")
	}

	objete, wybrane, err := dziennikWybierz(wykaz, z.ActionIds, stan.dokument.Kod,
		string(shared.StudioActionStateActive))
	if err != nil {
		return shared.StudioJournalRevertResponse{}, err
	}
	sort.Slice(wybrane, func(i, j int) bool { return wybrane[i].Kolejnosc > wybrane[j].Kolejnosc })

	stanWpisu := map[string]string{}
	for _, wpis := range wykaz {
		stanWpisu[wpis.Kod] = wpis.Stan
	}

	bilans := shared.StudioActionBalance{Skipped: []shared.StudioSkippedItem{}}
	cofniete := []string{}
	for _, czynnosc := range wybrane {
		if przeszkoda := dziennikStojaceNaNiej(czynnosc, wykaz, stanWpisu, objete); przeszkoda != "" {
			bilans.Skipped = append(bilans.Skipped,
				dziennikPominiecieZaleznosci(czynnosc, przeszkoda))
			continue
		}
		pominiete, err := a.dziennikNalozOdwrotnie(ctx, stan, czynnosc.StanPo, czynnosc.StanPrzed)
		if err != nil {
			return shared.StudioJournalRevertResponse{}, err
		}
		bilans.Skipped = append(bilans.Skipped, pominiete...)
		przestawiona, err := skladnica.PrzestawStanCzynnosci(ctx, czynnosc.Kod,
			string(shared.StudioActionStateActive), string(shared.StudioActionStateReverted))
		if err != nil {
			return shared.StudioJournalRevertResponse{}, bladStudio(err)
		}
		if !przestawiona {
			// Wiersz zmienił stan między odczytem wykazu a zapisem — inne okno
			// cofnęło tę czynność pierwsze. To nie jest usterka, ale nie wolno
			// tego policzyć jako skutku tego wywołania.
			bilans.Skipped = append(bilans.Skipped, shared.StudioSkippedItem{
				Reason: "czynność była już cofnięta",
				Detail: kontrolaWskaznikTekstu("czynność „" + czynnosc.Opis +
					"” cofnął w tym czasie ktoś inny — drugie cofnięcie nie zrobiło nic"),
			})
			continue
		}
		stanWpisu[czynnosc.Kod] = string(shared.StudioActionStateReverted)
		cofniete = append(cofniete, czynnosc.Kod)
		bilans.Applied++
	}

	if bilans.Applied == 0 {
		return shared.StudioJournalRevertResponse{}, dziennikBladZaleznosci(bilans.Skipped)
	}
	if err := a.postacZapisz(ctx, stan); err != nil {
		return shared.StudioJournalRevertResponse{}, err
	}
	dokument := stan.dokument
	if z.CreateVersion == nil || *z.CreateVersion {
		dokument, err = a.zalozWersjeDokumentu(ctx, dokument,
			wartoscTekstu(dokument.Tresc), shared.StudioAuthorUzytkownik, nil)
		if err != nil {
			return shared.StudioJournalRevertResponse{}, err
		}
	}
	bilans.SkippedCount = len(bilans.Skipped)
	bilans.Note = kontrolaWskaznikTekstu("Cofnięto " + strconv.Itoa(bilans.Applied) +
		" czynność(-i) dziennika wraz z postacią; wersje nowsze zostają, więc samo " +
		"cofnięcie jest odwracalne poleceniem studio.journal.redo.")
	return shared.StudioJournalRevertResponse{
		Reverted: cofniete,
		Document: a.zlozDokument(dokument),
		Form:     stan.forma,
		Balance:  bilans,
	}, nil
}

// PonowCzynnosc obsługuje `studio.journal.redo`.
//
// Ponowienie idzie w kolejności dziennika, nie odwrotnie: czynność starsza
// wchodzi pierwsza, bo młodsza mogła na niej stać.
func (a *adapterStudia) PonowCzynnosc(ctx context.Context,
	z shared.StudioJournalRedoRequest) (shared.StudioJournalRedoResponse, error) {

	stan, skladnica, wykaz, err := a.dziennikStanIWykaz(ctx, z.DocumentId)
	if err != nil {
		return shared.StudioJournalRedoResponse{}, err
	}
	if len(z.ActionIds) == 0 {
		return shared.StudioJournalRedoResponse{},
			bladWskazaniaStudio("ponowienie bez wskazania czynności")
	}
	objete, wybrane, err := dziennikWybierz(wykaz, z.ActionIds, stan.dokument.Kod,
		string(shared.StudioActionStateReverted))
	if err != nil {
		return shared.StudioJournalRedoResponse{}, err
	}
	sort.Slice(wybrane, func(i, j int) bool { return wybrane[i].Kolejnosc < wybrane[j].Kolejnosc })

	stanWpisu := map[string]string{}
	for _, wpis := range wykaz {
		stanWpisu[wpis.Kod] = wpis.Stan
	}

	bilans := shared.StudioActionBalance{Skipped: []shared.StudioSkippedItem{}}
	ponowione := []string{}
	for _, czynnosc := range wybrane {
		if przeszkoda := dziennikPodstawaCofnieta(czynnosc, stanWpisu, objete); przeszkoda != "" {
			bilans.Skipped = append(bilans.Skipped, shared.StudioSkippedItem{
				Reason: "podstawa czynności jest cofnięta",
				Detail: kontrolaWskaznikTekstu("czynności „" + czynnosc.Opis +
					"” nie da się ponowić, dopóki cofnięta stoi czynność " + przeszkoda +
					", na której ona stoi"),
			})
			continue
		}
		// Ponowienie jest cofnięciem cofnięcia: nakłada się różnicę w kierunku
		// od stanu sprzed czynności do stanu po niej.
		pominiete, err := a.dziennikNalozOdwrotnie(ctx, stan, czynnosc.StanPrzed, czynnosc.StanPo)
		if err != nil {
			return shared.StudioJournalRedoResponse{}, err
		}
		bilans.Skipped = append(bilans.Skipped, pominiete...)
		przestawiona, err := skladnica.PrzestawStanCzynnosci(ctx, czynnosc.Kod,
			string(shared.StudioActionStateReverted), string(shared.StudioActionStateActive))
		if err != nil {
			return shared.StudioJournalRedoResponse{}, bladStudio(err)
		}
		if !przestawiona {
			bilans.Skipped = append(bilans.Skipped, shared.StudioSkippedItem{
				Reason: "czynność stała już w dokumencie",
				Detail: kontrolaWskaznikTekstu("czynność „" + czynnosc.Opis +
					"” ponowił w tym czasie ktoś inny"),
			})
			continue
		}
		stanWpisu[czynnosc.Kod] = string(shared.StudioActionStateActive)
		ponowione = append(ponowione, czynnosc.Kod)
		bilans.Applied++
	}
	if bilans.Applied == 0 {
		return shared.StudioJournalRedoResponse{}, dziennikBladZaleznosci(bilans.Skipped)
	}
	if err := a.postacZapisz(ctx, stan); err != nil {
		return shared.StudioJournalRedoResponse{}, err
	}
	bilans.SkippedCount = len(bilans.Skipped)
	bilans.Note = kontrolaWskaznikTekstu("Ponowiono " + strconv.Itoa(bilans.Applied) +
		" czynność(-i) dziennika wraz z postacią.")
	return shared.StudioJournalRedoResponse{
		Redone:   ponowione,
		Document: a.zlozDokument(stan.dokument),
		Form:     stan.forma,
		Balance:  bilans,
	}, nil
}

// ── Wspólne dla cofania i ponawiania ────────────────────────────────────────

// dziennikStanIWykaz wczytuje postać dokumentu, składnicę odcinka i dziennik —
// trzy rzeczy, których obie drogi potrzebują naraz.
func (a *adapterStudia) dziennikStanIWykaz(ctx context.Context,
	kodDokumentu string) (*stanPostaci, KontrolaPracyStudia,
	[]dane.CzynnoscDokumentuStudia, error) {

	stan, err := a.postacWczytaj(ctx, kodDokumentu)
	if err != nil {
		return nil, nil, nil, err
	}
	skladnica, err := a.kontrolaSkladnica()
	if err != nil {
		return nil, nil, nil, err
	}
	wykaz, err := skladnica.CzynnosciDokumentu(ctx, stan.dokument.ID)
	if err != nil {
		return nil, nil, nil, bladStudio(err)
	}
	return stan, skladnica, wykaz, nil
}

// dziennikWybierz przesiewa dziennik do czynności wskazanych żądaniem.
//
// Wskazanie czynności, której nie ma, jest ODMOWĄ, nie ciszą: odpowiedź
// pomyślna na cofnięcie czegoś, czego dziennik nie zna, kazałaby Operatorowi
// czytać to jako wykonane. Tak samo wskazanie czynności w stanie innym niż
// wymagany — cofnięcie już cofniętej nie jest cofnięciem.
func dziennikWybierz(wykaz []dane.CzynnoscDokumentuStudia, kody []string,
	kodDokumentu, stanWymagany string) (map[string]bool,
	[]dane.CzynnoscDokumentuStudia, error) {

	poKodzie := map[string]dane.CzynnoscDokumentuStudia{}
	for _, wpis := range wykaz {
		poKodzie[wpis.Kod] = wpis
	}
	objete := map[string]bool{}
	wybrane := make([]dane.CzynnoscDokumentuStudia, 0, len(kody))
	for _, kod := range kody {
		wpis, jest := poKodzie[strings.TrimSpace(kod)]
		if !jest {
			return nil, nil, bladBrakuStudio("dziennik dokumentu " + kodDokumentu +
				" nie zna czynności " + kod)
		}
		if wpis.Stan != stanWymagany {
			return nil, nil, bladWskazaniaStudio("czynność " + kod + " jest w stanie „" +
				wpis.Stan + "”, a ta czynność wymaga stanu „" + stanWymagany + "”")
		}
		if objete[wpis.Kod] {
			continue
		}
		objete[wpis.Kod] = true
		wybrane = append(wybrane, wpis)
	}
	return objete, wybrane, nil
}

// dziennikStojaceNaNiej oddaje kod pierwszej czynności późniejszej, która stoi
// na tej i nie jest cofana razem z nią. Pusty napis znaczy „nic nie stoi".
//
// ── Dwa źródła zależności, nie jedno ────────────────────────────────────────
// Zależność ZAPISANA (tabela `zaleznosc_czynnosci_studio`) mówi o powiązaniach,
// których z zakresu nie widać: wstawienie tabeli jest podstawą scalenia komórki,
// choć zakresy mogą się nie stykać. Tę wiedzę ma czynność, która wpis odkładała,
// i ona ją zapisuje.
//
// Zależność WYWIEDZIONA z zakresu jest siatką pod tym: czynność późniejsza,
// która ruszyła TEN SAM fragment, stoi na tej wcześniejszej z samej natury
// rzeczy — cofnięcie wcześniejszej wgrałoby w to miejsce stan sprzed niej,
// a więc zabrałoby ze sobą pracę późniejszą i zostawiło dokument w stanie,
// którego nigdy nie było.
//
// Dlaczego oba, a nie tylko zapisane: wykaz zapisanych zależności wypełniają
// czynności czterech odcinków i jest dziś NIEPEŁNY — a niepełny wykaz zależności
// znaczy ciche cofnięcie, które psuje dokument. Wywiedzenie z zakresu myli się
// w drugą stronę: odmawia czasem cofnięcia, które byłoby bezpieczne, i mówi
// wtedy wprost, która czynność stoi na drodze. Odmowa nazwana jest tu tańsza
// niż niespójny dokument.
func dziennikStojaceNaNiej(czynnosc dane.CzynnoscDokumentuStudia,
	wykaz []dane.CzynnoscDokumentuStudia, stanWpisu map[string]string,
	objete map[string]bool) string {

	for _, kod := range czynnosc.StojaceNaNiej {
		if objete[kod] {
			continue
		}
		if stanWpisu[kod] == string(shared.StudioActionStateActive) {
			return kod
		}
	}
	return dziennikStojaceZZakresu(czynnosc, wykaz, stanWpisu, objete)
}

// dziennikStojaceZZakresu wywodzi zależność ze stykających się zakresów.
//
// Czynność bez zakresu (nastawa strony, odświeżenie aparatu) nie wywodzi
// zależności ani w jedną, ani w drugą stronę: nie ma czego porównać, a
// przyjęcie „brak zakresu znaczy cały dokument" zablokowałoby cofanie
// czegokolwiek po pierwszej takiej czynności.
func dziennikStojaceZZakresu(czynnosc dane.CzynnoscDokumentuStudia,
	wykaz []dane.CzynnoscDokumentuStudia, stanWpisu map[string]string,
	objete map[string]bool) string {

	od, do, jest := dziennikZakresCzynnosci(czynnosc)
	if !jest {
		return ""
	}
	for _, pozniejsza := range wykaz {
		if pozniejsza.Kolejnosc <= czynnosc.Kolejnosc || objete[pozniejsza.Kod] {
			continue
		}
		if stanWpisu[pozniejsza.Kod] != string(shared.StudioActionStateActive) {
			continue
		}
		odPo, doPo, jestPo := dziennikZakresCzynnosci(pozniejsza)
		if !jestPo {
			continue
		}
		if kontrolaZakresyStykaja(od, do, odPo, doPo) {
			return pozniejsza.Kod
		}
	}
	return ""
}

// dziennikZakresCzynnosci oddaje zakres wpisu dziennika w znakach.
func dziennikZakresCzynnosci(czynnosc dane.CzynnoscDokumentuStudia) (int, int, bool) {
	if czynnosc.ZakresOd == nil || czynnosc.ZakresDo == nil {
		return 0, 0, false
	}
	return int(*czynnosc.ZakresOd), int(*czynnosc.ZakresDo), true
}

// dziennikPodstawaCofnieta oddaje kod pierwszej podstawy, która jest cofnięta
// i nie jest ponawiana razem z tą czynnością.
func dziennikPodstawaCofnieta(czynnosc dane.CzynnoscDokumentuStudia,
	stanWpisu map[string]string, objete map[string]bool) string {

	for _, kod := range czynnosc.PodstawyKody {
		if objete[kod] {
			continue
		}
		if stanWpisu[kod] == string(shared.StudioActionStateReverted) {
			return kod
		}
	}
	return ""
}

// dziennikPominiecieZaleznosci składa pozycję bilansu NAZYWAJĄCĄ zależność.
func dziennikPominiecieZaleznosci(czynnosc dane.CzynnoscDokumentuStudia,
	przeszkoda string) shared.StudioSkippedItem {

	pozycja := shared.StudioSkippedItem{
		Reason: "czynność jest podstawą późniejszej",
		Detail: kontrolaWskaznikTekstu("czynności „" + czynnosc.Opis +
			"” nie da się cofnąć, bo stoi na niej czynność " + przeszkoda +
			" — cofnij najpierw ją, a potem tę"),
	}
	if czynnosc.ZakresOd != nil {
		od := int(*czynnosc.ZakresOd)
		pozycja.RangeStart = &od
	}
	if czynnosc.ZakresDo != nil {
		do := int(*czynnosc.ZakresDo)
		pozycja.RangeEnd = &do
	}
	return pozycja
}

// dziennikBladZaleznosci zamienia bilans samych pominięć w odmowę nazwaną.
//
// Odpowiedź pomyślna bez ani jednej cofniętej czynności byłaby najgorszym
// możliwym wynikiem: Operator zobaczyłby „gotowe" i uznał, że dokument wrócił
// do stanu sprzed pomyłki.
func dziennikBladZaleznosci(pominiete []shared.StudioSkippedItem) error {
	powody := make([]string, 0, len(pominiete))
	for _, pozycja := range pominiete {
		if pozycja.Detail != nil {
			powody = append(powody, *pozycja.Detail)
			continue
		}
		powody = append(powody, pozycja.Reason)
	}
	tresc := "moduł Studio, dziennik czynności: nie cofnięto ani jednej czynności"
	if len(powody) > 0 {
		tresc += " — " + strings.Join(powody, "; ")
	}
	return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeConflict, tresc))
}

// ── Nakładanie różnicy drzew postaci ────────────────────────────────────────

// dziennikNalozOdwrotnie nakłada na stan BIEŻĄCY różnicę między drzewem `zStanu`
// i drzewem `naStan`, byt po bycie.
//
// Kierunek podaje wołający: cofnięcie idzie od `StanPo` do `StanPrzed`,
// ponowienie od `StanPrzed` do `StanPo`. Rachunek jest ten sam, więc stoi
// w jednym miejscu — dwie kopie rozjechałyby się przy pierwszej poprawce.
//
// Nakłada się RÓŻNICĘ, a nie całe drzewo: wgranie `StanPrzed` w miejsce
// dokumentu zabrałoby ze sobą wszystko, co po tej czynności weszło, czyli
// byłoby przywróceniem wersji, a nie cofnięciem pojedynczej czynności.
func (a *adapterStudia) dziennikNalozOdwrotnie(ctx context.Context, stan *stanPostaci,
	zStanu, naStan *string) ([]shared.StudioSkippedItem, error) {

	zrodlo, err := dziennikDrzewo(zStanu)
	if err != nil {
		return nil, err
	}
	cel, err := dziennikDrzewo(naStan)
	if err != nil {
		return nil, err
	}
	if zrodlo == nil || cel == nil {
		return []shared.StudioSkippedItem{{
			Reason: "wpis dziennika nie niesie stanu postaci",
			Detail: kontrolaWskaznikTekstu("czynność odłożono bez drzewa postaci sprzed " +
				"albo po zmianie, więc nie ma z czego jej odtworzyć — cofnięcie takiej " +
				"czynności wymagałoby zgadywania"),
		}}, nil
	}

	pominiete := []shared.StudioSkippedItem{}

	// Bloki — treść i postać znaku oraz akapitu. Tożsamością bloku jest jego
	// identyfikator; blok, którego czynność nie tknęła, zostaje nietknięty.
	dziennikNalozBloki(&stan.forma, zrodlo, cel)

	// Nastawy strony są jednym bytem, więc nakłada się je w całości — ale tylko
	// wtedy, gdy czynność je naprawdę zmieniła.
	if dziennikRoznePola(zrodlo.PageSetup, cel.PageSetup) {
		stan.forma.PageSetup = cel.PageSetup
	}

	// Styl nazwany, sekcja, obiekt i pole mają WŁASNE WIERSZE i drzewo ich nie
	// zapisuje — dlatego wracają osobnym zapisem.
	odrzucone, err := a.dziennikPrzywrocByty(ctx, stan, zrodlo, cel)
	if err != nil {
		return nil, err
	}
	pominiete = append(pominiete, odrzucone...)
	return pominiete, nil
}

// dziennikDrzewo czyta drzewo postaci z ładunku wpisu dziennika.
func dziennikDrzewo(zapis *string) (*shared.StudioDocumentForm, error) {
	if zapis == nil || strings.TrimSpace(*zapis) == "" {
		return nil, nil
	}
	var drzewo shared.StudioDocumentForm
	if err := json.Unmarshal([]byte(*zapis), &drzewo); err != nil {
		// Ładunek nieczytelny NIE zamienia się w brak: to zamieniłoby uszkodzenie
		// wpisu w ciszę, a Operator uznałby, że czynność nie miała czego cofać.
		return nil, kontrolaBladZaplecza("wpis dziennika niesie nieczytelne drzewo " +
			"postaci: " + err.Error())
	}
	return &drzewo, nil
}

// dziennikNalozBloki nakłada różnicę bloków: blok zmieniony wraca do brzmienia
// docelowego, blok dodany czynnością znika, blok zdjęty czynnością wraca na
// swoje miejsce.
func dziennikNalozBloki(biezaca, zrodlo, cel *shared.StudioDocumentForm) {
	wZrodle := dziennikBlokiPoKodzie(zrodlo.Blocks)
	wCelu := dziennikBlokiPoKodzie(cel.Blocks)

	nowe := make([]shared.StudioDocumentBlock, 0, len(biezaca.Blocks)+len(cel.Blocks))
	obecne := map[string]bool{}
	for _, blok := range biezaca.Blocks {
		zrodlowy, byłWZrodle := wZrodle[blok.Id]
		docelowy, jestWCelu := wCelu[blok.Id]
		switch {
		case byłWZrodle && !jestWCelu:
			// Czynność ten blok dodała — nakładając różnicę w drugą stronę,
			// blok schodzi.
			continue
		case jestWCelu && (!byłWZrodle || !dziennikRownyBlok(zrodlowy, docelowy)):
			nowe = append(nowe, docelowy)
		default:
			nowe = append(nowe, blok)
		}
		obecne[blok.Id] = true
	}
	// Blok, który czynność zdjęła, wraca na miejsce wyznaczone kolejnością
	// drzewa docelowego.
	for numer, blok := range cel.Blocks {
		if obecne[blok.Id] {
			continue
		}
		if numer >= len(nowe) {
			nowe = append(nowe, blok)
			continue
		}
		nowe = append(nowe[:numer], append([]shared.StudioDocumentBlock{blok},
			nowe[numer:]...)...)
	}
	biezaca.Blocks = nowe
}

// dziennikBlokiPoKodzie układa bloki drzewa pod ich identyfikatorami.
func dziennikBlokiPoKodzie(
	bloki []shared.StudioDocumentBlock) map[string]shared.StudioDocumentBlock {

	poKodzie := make(map[string]shared.StudioDocumentBlock, len(bloki))
	for _, blok := range bloki {
		poKodzie[blok.Id] = blok
	}
	return poKodzie
}

// dziennikRownyBlok mówi, czy dwa bloki są tym samym stanem. Porównanie idzie
// zapisem JSON, bo blok niesie zagnieżdżone postaci znaku i akapitu, a
// porównanie pole po polu rozjechałoby się przy pierwszym dołożonym polu
// kontraktu.
func dziennikRownyBlok(pierwszy, drugi shared.StudioDocumentBlock) bool {
	return dziennikRoznePola(pierwszy, drugi) == false
}

// dziennikRoznePola mówi, czy dwie wartości kontraktu różnią się treścią.
func dziennikRoznePola(pierwsza, druga any) bool {
	zapisPierwszej, blad := json.Marshal(pierwsza)
	if blad != nil {
		return true
	}
	zapisDrugiej, blad := json.Marshal(druga)
	if blad != nil {
		return true
	}
	return string(zapisPierwszej) != string(zapisDrugiej)
}

// dziennikPrzywrocByty przywraca byty o własnych wierszach: styl nazwany,
// sekcję, obiekt osadzony, pole dokumentu i aparat dokumentu.
//
// Przekłady „byt kontraktu → wiersz" bierze się z obszarów, które te wiersze
// zakładają (`postacStylDoWiersza`, `wejscieSekcjaDoWiersza`,
// `wejscieObiektDoWiersza`, `wejsciePoleDoWiersza`, `aparatDoWiersza`). Drugi
// przekład tej samej tabeli byłby drugą prawdą o jednym wierszu — i właśnie
// dlatego cofnięcie woła cudze przekłady, a nie pisze własnych.
func (a *adapterStudia) dziennikPrzywrocByty(ctx context.Context, stan *stanPostaci,
	zrodlo, cel *shared.StudioDocumentForm) ([]shared.StudioSkippedItem, error) {

	skladnica, err := a.postacSkladnica()
	if err != nil {
		return nil, err
	}
	pominiete := []shared.StudioSkippedItem{}

	// Styl nazwany — tożsamością jest nazwa.
	wCelu := map[string]shared.StudioNamedStyle{}
	for _, styl := range cel.Styles {
		wCelu[styl.Name] = styl
	}
	for _, styl := range zrodlo.Styles {
		docelowy, jest := wCelu[styl.Name]
		if !jest {
			if _, err := skladnica.UsunStylNazwany(ctx, stan.dokument.ID, styl.Name); err != nil {
				return nil, bladStudio(err)
			}
			continue
		}
		if !dziennikRoznePola(styl, docelowy) {
			continue
		}
		wiersz, err := postacStylDoWiersza(stan.dokument.ID, docelowy)
		if err != nil {
			return nil, err
		}
		if _, err := skladnica.ZapiszStylNazwany(ctx, wiersz); err != nil {
			return nil, bladStudio(err)
		}
	}
	for nazwa, styl := range wCelu {
		if dziennikMaStyl(zrodlo.Styles, nazwa) {
			continue
		}
		wiersz, err := postacStylDoWiersza(stan.dokument.ID, styl)
		if err != nil {
			return nil, err
		}
		if _, err := skladnica.ZapiszStylNazwany(ctx, wiersz); err != nil {
			return nil, bladStudio(err)
		}
	}

	// Sekcje — tożsamością jest identyfikator sekcji.
	sekcjeCelu := map[string]shared.StudioSection{}
	for _, sekcja := range cel.Sections {
		sekcjeCelu[sekcja.Id] = sekcja
	}
	for _, sekcja := range zrodlo.Sections {
		if _, jest := sekcjeCelu[sekcja.Id]; jest {
			continue
		}
		if _, err := skladnica.UsunSekcje(ctx, sekcja.Id); err != nil {
			return nil, bladStudio(err)
		}
	}
	for numer, sekcja := range cel.Sections {
		wiersz, err := wejscieSekcjaDoWiersza(stan.dokument.ID, numer, sekcja)
		if err != nil {
			return nil, err
		}
		if _, err := skladnica.ZapiszSekcje(ctx, wiersz); err != nil {
			return nil, bladStudio(err)
		}
	}

	// Obiekty osadzone.
	obiektyCelu := map[string]shared.StudioDocumentObject{}
	for _, obiekt := range cel.Objects {
		obiektyCelu[obiekt.Id] = obiekt
	}
	for _, obiekt := range zrodlo.Objects {
		if _, jest := obiektyCelu[obiekt.Id]; jest {
			continue
		}
		if _, err := skladnica.UsunObiektDokumentu(ctx, obiekt.Id); err != nil {
			return nil, bladStudio(err)
		}
	}
	for _, obiekt := range cel.Objects {
		wiersz, err := wejscieObiektDoWiersza(stan.dokument.ID, obiekt)
		if err != nil {
			return nil, err
		}
		if _, err := skladnica.ZapiszObiektDokumentu(ctx, wiersz); err != nil {
			return nil, bladStudio(err)
		}
	}

	// Pola dokumentu.
	polaCelu := map[string]shared.StudioDocumentField{}
	for _, pole := range cel.Fields {
		polaCelu[pole.Id] = pole
	}
	for _, pole := range zrodlo.Fields {
		if _, jest := polaCelu[pole.Id]; jest {
			continue
		}
		if _, err := skladnica.UsunPoleDokumentu(ctx, pole.Id); err != nil {
			return nil, bladStudio(err)
		}
	}
	for _, pole := range cel.Fields {
		if _, err := skladnica.ZapiszPoleDokumentu(ctx,
			wejsciePoleDoWiersza(stan.dokument.ID, pole)); err != nil {

			return nil, bladStudio(err)
		}
	}

	// Aparat dokumentu — tożsamością jest identyfikator elementu.
	//
	// Do 17.08.2026 stała tu ODMOWA nazwana: cofnięcie czynności odtwarzało treść
	// i postać, a spis treści, przypis albo bibliografia zostawały w brzmieniu
	// bieżącym, bo dziennik nie miał przekładu „element kontraktu → wiersz".
	// Przekład istniał — `aparatDoWiersza` obszaru aparatu — i był w tym pliku
	// nieużyty. Aparat wraca więc tak samo jak sekcja i obiekt, a bilans
	// nazywający brak zniknął razem z brakiem.
	//
	// Kotwica i znacznik nieświeżości idą Z ELEMENTU, nie z rachunku: cofnięcie
	// ma przywrócić stan sprzed czynności, a nie przeliczyć aparat od nowa.
	// Przeliczenie wstawiłoby numery stron liczone na treści BIEŻĄCEJ, czyli
	// numery, których w stanie sprzed nie było.
	aparatCelu := map[string]shared.StudioApparatusItem{}
	for _, element := range cel.Apparatus {
		aparatCelu[element.Id] = element
	}
	for _, element := range zrodlo.Apparatus {
		if _, jest := aparatCelu[element.Id]; jest {
			continue
		}
		if _, err := skladnica.UsunElementAparatu(ctx, element.Id); err != nil {
			return nil, bladStudio(err)
		}
	}
	for _, element := range cel.Apparatus {
		wiersz, err := aparatDoWiersza(stan.dokument.ID, element.Id, element,
			aparatWartoscLiczby(element.AnchorStart), aparatWartoscLiczby(element.AnchorEnd),
			element.Stale != nil && *element.Stale)
		if err != nil {
			return nil, err
		}
		if _, err := skladnica.ZapiszElementAparatu(ctx, wiersz); err != nil {
			return nil, bladStudio(err)
		}
	}

	// Wykaz bytów w drzewie musi zgadzać się z tym, co stoi w wierszach —
	// inaczej następny odczyt postaci pokazałby stan sprzed przywrócenia.
	stan.forma.Styles = cel.Styles
	stan.forma.Sections = cel.Sections
	stan.forma.Objects = cel.Objects
	stan.forma.Fields = cel.Fields
	stan.forma.Apparatus = cel.Apparatus
	return pominiete, nil
}

// dziennikMaStyl mówi, czy arkusz niesie styl o tej nazwie.
func dziennikMaStyl(style []shared.StudioNamedStyle, nazwa string) bool {
	for _, styl := range style {
		if styl.Name == nazwa {
			return true
		}
	}
	return false
}

// ── Przekład wpisu dziennika ────────────────────────────────────────────────

// dziennikZlozCzynnosc składa wpis dziennika kontraktu z wiersza warstwy danych.
func dziennikZlozCzynnosc(wiersz dane.CzynnoscDokumentuStudia) shared.StudioDocumentAction {
	czynnosc := shared.StudioDocumentAction{
		Id:                 wiersz.Kod,
		DocumentId:         wiersz.DokumentKod,
		Kind:               shared.StudioActionKind(wiersz.Rodzaj),
		Author:             shared.StudioAuthor(wiersz.AutorRodzaj),
		Description:        wiersz.Opis,
		RangeStart:         kontrolaWskaznikLiczby(wiersz.ZakresOd),
		RangeEnd:           kontrolaWskaznikLiczby(wiersz.ZakresDo),
		DependsOn:          wiersz.PodstawyKody,
		BlockedBy:          wiersz.StojaceNaNiej,
		State:              shared.StudioActionState(wiersz.Stan),
		ChangeId:           wiersz.ZmianaSledzonaK,
		Sequence:           wiersz.Kolejnosc,
		CreatedAt:          chwilaBazy(wiersz.Utworzono),
		AuthorAgentId:      wiersz.AutorAgentKod,
		AuthorAgentName:    wiersz.AutorAgentNazwa,
		AuthorAgentVersion: wiersz.AutorAgentWersja,
		AuthorSubagentId:   wiersz.AutorPodagentKod,
	}
	return czynnosc
}

// dziennikOdlozCzynnosc odkłada wpis dziennika dla czynności TEGO odcinka —
// znakowania, schowka, przeniesienia fragmentu różnicy, przywrócenia kopii.
//
// Czynności postaci odkłada `postacOdlozCzynnosc`; ta droga jest jej
// odpowiednikiem dla czynności, które nie idą przez drzewo postaci. Przedrostek
// kodu jest ten sam, bo dziennik jest JEDEN — dwa przedrostki rozdzieliłyby go
// na dwa nieporównywalne szeregi.
func (a *adapterStudia) dziennikOdlozCzynnosc(ctx context.Context, dokument dane.DokumentStudia,
	wykonawca kontrolaWykonawca, rodzaj shared.StudioActionKind, opis string,
	od, do *int, drzewoPrzed, drzewoPo *string, zmianaKod *string) (*string, error) {

	skladnica, err := a.kontrolaSkladnica()
	if err != nil {
		return nil, err
	}
	wpis := dane.CzynnoscDokumentuStudia{
		Kod:              nowyIdentyfikator(przedrostekCzynnosciStudia),
		DokumentKod:      dokument.Kod,
		Rodzaj:           string(rodzaj),
		AutorRodzaj:      string(wykonawca.Rodzaj),
		AutorAgentKod:    wykonawca.AgentKod,
		AutorAgentNazwa:  wykonawca.AgentNazwa,
		AutorAgentWersja: wykonawca.AgentWersja,
		AutorPodagentKod: wykonawca.PodagentKod,
		Opis:             opis,
		ZakresOd:         kontrolaLiczbaZWskaznika(od),
		ZakresDo:         kontrolaLiczbaZWskaznika(do),
		StanPrzed:        drzewoPrzed,
		StanPo:           drzewoPo,
		ZmianaSledzonaK:  zmianaKod,
		// Stan wpisu jest słownikiem tabeli dziennika: `active` albo `reverted`.
		Stan: string(shared.StudioActionStateActive),
	}
	zapisana, err := skladnica.ZapiszCzynnoscDokumentu(ctx, dokument.ID, wpis)
	if err != nil {
		return nil, bladStudio(err)
	}
	return &zapisana.Kod, nil
}

// dziennikZapisDrzewa składa ładunek drzewa postaci do wpisu dziennika.
func dziennikZapisDrzewa(forma shared.StudioDocumentForm) (*string, error) {
	zapis, err := json.Marshal(forma)
	if err != nil {
		return nil, kontrolaBladZaplecza("drzewa postaci nie da się odłożyć w dzienniku: " +
			err.Error())
	}
	tekst := string(zapis)
	return &tekst, nil
}

// dziennikBrakWiersza odróżnia „bytu nie ma" od „odczyt się nie powiódł" dla
// bytów tego odcinka.
func dziennikBrakWiersza(err error, powod string) error {
	if errors.Is(err, dane.ErrBrakWiersza) {
		return bladBrakuStudio(powod)
	}
	return bladStudio(err)
}
