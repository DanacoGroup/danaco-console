// Plik obsługuje sterowanie pojedynczym krokiem zlecenia: wstrzymanie kroku,
// zapisanie decyzji Operatora oraz wznowienie procesu z tą decyzją zastosowaną,
// zamiast wznowienia od początku pracy.
package core

import (
	"context"
	"strconv"
	"strings"

	"danacoconsole/server/internal/dane"
	"danacoconsole/server/internal/protocol"
	"danacoconsole/shared"
)

// Stany sterowania krokiem widziane przez Operatora: trzy zapisuje tabela
// wstrzymanie_kroku, a czeka i biegnie wyprowadza się z kolumny pozycja_kolejki.stan.
const (
	stanKrokuCzeka        = "czeka"
	stanKrokuBiegnie      = "biegnie"
	stanKrokuWstrzymany   = dane.StanKrokuWstrzymany
	stanKrokuZatwierdzony = dane.StanKrokuZatwierdzony
	stanKrokuOdrzucony    = dane.StanKrokuOdrzucony
	stanKrokuZamkniety    = "zamkniety"
)

// Decyzje Operatora przyjmowane z ładunku. Nazwy angielskie, bo to wartości
// koperty kontraktu; ich odpowiedniki po stronie zapisu są w `dane`.
const (
	decyzjaZatwierdz = "approve"
	decyzjaOdrzuc    = "reject"
)

// Stan pracy kroku w kopercie niesie wartości wyliczenia QueueStepWorkStatus;
// koperta kontraktu jest angielska, a kolumna bazy polska, więc rdzeń trzyma
// osobny słownik przekładu między nimi.
const (
	krokPracaOczekuje     = "pending"
	krokPracaPrzydzielona = "assigned"
	krokPracaWykonywana   = "running"
	krokPracaWeryfikacja  = "review"
	krokPracaUkonczona    = "done"
	krokPracaBledna       = "failed"
	krokPracaAnulowana    = "cancelled"
)

// Stan sterowania krokiem w kopercie niesie wartości wyliczenia QueueStepControlStatus,
// widziane przez Operatora niezależnie od stanu pracy nad krokiem.
const (
	krokSterowanieCzeka        = "waiting"
	krokSterowanieBiegnie      = "running"
	krokSterowanieWstrzymany   = "held"
	krokSterowanieZatwierdzony = "approved"
	krokSterowanieOdrzucony    = "rejected"
	krokSterowanieZamkniety    = "closed"
)

// stanPracyKrokuNaKontrakt wiąże siedem wartości kolumny `pozycja_kolejki.stan`
// dopuszczonych jej więzem CHECK z wartościami koperty.
var stanPracyKrokuNaKontrakt = map[string]string{
	stanPozycjiOczekuje:      krokPracaOczekuje,
	stanPozycjiPrzydzielona:  krokPracaPrzydzielona,
	stanPozycjiWykonywana:    krokPracaWykonywana,
	stanPozycjiDoWeryfikacji: krokPracaWeryfikacja,
	stanPozycjiUkonczona:     krokPracaUkonczona,
	stanPozycjiBledna:        krokPracaBledna,
	stanPozycjiAnulowana:     krokPracaAnulowana,
}

// stanSterowaniaKrokuNaKontrakt wiąże stany sterowania z kopertą. Trzy z nich
// mają kolumnę (`wstrzymanie_kroku.stan`), pozostałe trzy są wyprowadzane.
var stanSterowaniaKrokuNaKontrakt = map[string]string{
	stanKrokuCzeka:        krokSterowanieCzeka,
	stanKrokuBiegnie:      krokSterowanieBiegnie,
	stanKrokuWstrzymany:   krokSterowanieWstrzymany,
	stanKrokuZatwierdzony: krokSterowanieZatwierdzony,
	stanKrokuOdrzucony:    krokSterowanieOdrzucony,
	stanKrokuZamkniety:    krokSterowanieZamkniety,
}

// stanPracyKontraktu i stanSterowaniaKontraktu odmawiają przy wartości spoza
// słownika, zamiast puścić ją na drut: cichy przepust wystawiłby polskie słowo
// pod angielskim wyliczeniem, czyli kontrakt mówiłby jedno, a drut niósł drugie.
func stanPracyKontraktu(stan string) (string, error) {
	wartosc, znana := stanPracyKrokuNaKontrakt[stan]
	if !znana {
		return "", bladSterowaniaKrokiem(shared.ErrorCodeInternalError,
			"wartość kolumny pozycja_kolejki.stan nie ma odpowiednika w kontrakcie: "+stan)
	}
	return wartosc, nil
}

func stanSterowaniaKontraktu(stan string) (string, error) {
	wartosc, znana := stanSterowaniaKrokuNaKontrakt[stan]
	if !znana {
		return "", bladSterowaniaKrokiem(shared.ErrorCodeInternalError,
			"stan sterowania krokiem nie ma odpowiednika w kontrakcie: "+stan)
	}
	return wartosc, nil
}

// werdyktOdrzucone to wartość kolumny `pozycja_kolejki.werdykt_weryfikacji`
// zapisywana wyłącznie przy odrzuceniu kroku przez Operatora. Werdykt wystawiony
// bez czynności Operatora byłby oceną wymyśloną przez silnik.
const werdyktOdrzucone = "odrzucone"

// Adnotacje doręczenia. Wychodzą w odpowiedzi zamiast milczenia, bo „decyzja
// zapisana, ale nikt jej nie wykonał" to fakt, który Operator ma zobaczyć,
// a nie domyślić się z braku efektu.
const (
	doreczenieNiepotrzebne = "krok nie wymagał wykonania — decyzja zastosowana wprost do jego stanu"
	doreczenieBezWykonawcy = "brak wpiętego wykonawcy kroku — decyzja zapisana i czeka na doręczenie"
)

// wstrzymania sięga po rozszerzenie repozytorium kolejek o sterowanie krokiem.
// Repozytorium bez tego rozszerzenia zostawia całą rodzinę nieczynną — port
// odmawia wtedy wprost, zamiast udawać wstrzymanie.
func (a *adapterKolejek) wstrzymania() (dane.RepozytoriumWstrzymanKroku, bool) {
	rozszerzone, ok := a.repozytorium.(dane.RepozytoriumWstrzymanKroku)
	return rozszerzone, ok
}

// WstrzymajKrok zatrzymuje jeden krok i zostawia go czekającego na decyzję
// Operatora; kolejka przechodzi w paused tylko wtedy, gdy wstrzymany krok jest
// tym, na którym kolejka aktualnie stoi.
func (a *adapterKolejek) WstrzymajKrok(ctx context.Context,
	z zadanieWstrzymaniaKroku) (odpowiedzKroku, error) {

	repozytorium, ok := a.wstrzymania()
	if !ok {
		return odpowiedzKroku{}, bladSterowaniaKrokiem(shared.ErrorCodeInternalError,
			"repozytorium kolejek nie niesie sterowania krokiem")
	}
	kolejkaID, pozycje, pozycja, err := a.odnajdzKrok(ctx, z.QueueId, z.StepId)
	if err != nil {
		return odpowiedzKroku{}, err
	}
	// Krok zamknięty nie ma czego wstrzymać — nie ma jak opuścić stanu pracy nad nim.
	if czyStanKoncowyPozycji(pozycja.Stan) {
		return odpowiedzKroku{}, bladSterowaniaKrokiem(shared.ErrorCodeConflict,
			"krok "+z.StepId+" jest zamknięty ("+pozycja.Stan+") — nie ma czego wstrzymać")
	}
	if _, jest, err := repozytorium.CzynneWstrzymanie(ctx, pozycja.ID); err != nil {
		return odpowiedzKroku{}, err
	} else if jest {
		return odpowiedzKroku{}, bladSterowaniaKrokiem(shared.ErrorCodeConflict,
			"krok "+z.StepId+" jest już wstrzymany i czeka na decyzję")
	}

	wstrzymanie, err := repozytorium.WstrzymajKrok(ctx, pozycja.ID, pozycja.Stan,
		tekstNiepusty(z.Reason))
	if err != nil {
		return odpowiedzKroku{}, err
	}
	if biezaca := pierwszaCzynna(pozycje); biezaca != nil && biezaca.ID == pozycja.ID {
		if err := a.repozytorium.ZmienStanKolejki(ctx, kolejkaID,
			stanKolejkiPoDzialaniu(shared.QueueActionPause, pozycje), shared.QueueActionPause); err != nil {
			return odpowiedzKroku{}, err
		}
	}
	kolejka, err := a.odpowiedz(ctx, kolejkaID, etapDzialanieNaKolejce)
	if err != nil {
		return odpowiedzKroku{}, err
	}
	krok, err := krokZleceniaKontraktu(kolejkaID, *pozycja, wstrzymanie, true)
	if err != nil {
		return odpowiedzKroku{}, err
	}
	return odpowiedzKroku{Step: krok, Queue: kolejka}, nil
}

// ZdecydujOKroku przyjmuje rozstrzygnięcie Operatora i wznawia proces z tą
// decyzją zastosowaną; decyzja zapisuje się przed zastosowaniem, więc
// przerwanie rdzenia między jednym a drugim zostawia ją widoczną i czekającą,
// nie zgubioną.
func (a *adapterKolejek) ZdecydujOKroku(ctx context.Context,
	z zadanieDecyzjiKroku) (odpowiedzDecyzjiKroku, error) {

	repozytorium, ok := a.wstrzymania()
	if !ok {
		return odpowiedzDecyzjiKroku{}, bladSterowaniaKrokiem(shared.ErrorCodeInternalError,
			"repozytorium kolejek nie niesie sterowania krokiem")
	}
	stanDecyzji, err := stanDecyzjiZLadunku(z.Decision)
	if err != nil {
		return odpowiedzDecyzjiKroku{}, err
	}
	kolejkaID, _, pozycja, err := a.odnajdzKrok(ctx, z.QueueId, z.StepId)
	if err != nil {
		return odpowiedzDecyzjiKroku{}, err
	}
	// Krok zamknięty w międzyczasie nie wraca do pracy przez decyzję nad nim.
	if czyStanKoncowyPozycji(pozycja.Stan) {
		return odpowiedzDecyzjiKroku{}, bladSterowaniaKrokiem(shared.ErrorCodeConflict,
			"krok "+z.StepId+" został zamknięty ("+pozycja.Stan+") — decyzji nie ma do czego zastosować")
	}
	wstrzymanie, jest, err := repozytorium.CzynneWstrzymanie(ctx, pozycja.ID)
	if err != nil {
		return odpowiedzDecyzjiKroku{}, err
	}
	if !jest {
		return odpowiedzDecyzjiKroku{}, bladSterowaniaKrokiem(shared.ErrorCodeConflict,
			"krok "+z.StepId+" nie jest wstrzymany — nie ma czego wznawiać")
	}
	// Decyzja zapada raz — rozstrzygnięty, lecz niezastosowany epizod wraca bez nowego werdyktu.
	if !wstrzymanie.CzyZdecydowane() {
		wstrzymanie, err = repozytorium.ZapiszDecyzje(ctx, wstrzymanie.ID, stanDecyzji,
			tekstNiepusty(z.Note))
		if err != nil {
			return odpowiedzDecyzjiKroku{}, err
		}
	} else if wstrzymanie.Stan != stanDecyzji {
		return odpowiedzDecyzjiKroku{}, bladSterowaniaKrokiem(shared.ErrorCodeConflict,
			"krok "+z.StepId+" ma już decyzję „"+wstrzymanie.Stan+"”, niezastosowaną — powtórz ją, nie zmieniaj")
	}

	doreczono, adnotacja, err := a.zastosujDecyzje(ctx, *pozycja, wstrzymanie)
	if err != nil {
		return odpowiedzDecyzjiKroku{}, err
	}
	wstrzymanie, err = repozytorium.OznaczZastosowanie(ctx, wstrzymanie.ID, doreczono)
	if err != nil {
		return odpowiedzDecyzjiKroku{}, err
	}

	poZmianie := a.silnik.Pozycje(ctx, kolejkaID)
	stanKolejki := stanKolejkiPoDzialaniu(shared.QueueActionResume, poZmianie)
	if err := a.repozytorium.ZmienStanKolejki(ctx, kolejkaID, stanKolejki, shared.QueueActionResume); err != nil {
		return odpowiedzDecyzjiKroku{}, err
	}
	kolejka, err := a.odpowiedz(ctx, kolejkaID, etapDzialanieNaKolejce)
	if err != nil {
		return odpowiedzDecyzjiKroku{}, err
	}
	krok, err := krokZleceniaKontraktu(kolejkaID, *a.pozycjaPoID(poZmianie, pozycja.ID, *pozycja), wstrzymanie, true)
	if err != nil {
		return odpowiedzDecyzjiKroku{}, err
	}
	return odpowiedzDecyzjiKroku{
		Step: krok, Queue: kolejka, Delivered: doreczono, DeliveryNote: adnotacja,
	}, nil
}

// WykazKrokow oddaje kroki zlecenia ze stanem sterowania. Bez tego wykazu
// żadne okno nie ma czego pokazać, a przycisk „zatwierdź krok" prowadziłby
// donikąd.
func (a *adapterKolejek) WykazKrokow(ctx context.Context,
	z zadanieWykazuKrokow) (odpowiedzWykazuKrokow, error) {

	kolejkaID, err := identyfikatorKolejki(z.QueueId)
	if err != nil {
		return odpowiedzWykazuKrokow{}, err
	}
	if _, err := a.repozytorium.PobierzKolejke(ctx, kolejkaID); err != nil {
		return odpowiedzWykazuKrokow{}, bladSterowaniaKrokiem(shared.ErrorCodeNotFound,
			"kolejka nie istnieje: "+z.QueueId)
	}
	pozycje, err := a.repozytorium.ListaPozycji(ctx, kolejkaID)
	if err != nil {
		return odpowiedzWykazuKrokow{}, err
	}
	// Brak rozszerzenia repozytorium znaczy nieczynne sterowanie krokiem, nie błąd odczytu.
	wstrzymania := map[int64]dane.WstrzymanieKroku{}
	if repozytorium, ok := a.wstrzymania(); ok {
		wykaz, err := repozytorium.CzynneWstrzymaniaKolejki(ctx, kolejkaID)
		if err != nil {
			return odpowiedzWykazuKrokow{}, err
		}
		wstrzymania = wykaz
	}
	kroki := make([]krokZlecenia, 0, len(pozycje))
	for _, pozycja := range pozycje {
		wstrzymanie, jest := wstrzymania[pozycja.ID]
		krok, err := krokZleceniaKontraktu(kolejkaID, pozycja, wstrzymanie, jest)
		if err != nil {
			return odpowiedzWykazuKrokow{}, err
		}
		kroki = append(kroki, krok)
	}
	return odpowiedzWykazuKrokow{Steps: kroki}, nil
}

// zastosujDecyzje przeprowadza krok tam, dokąd prowadzi go decyzja Operatora,
// i mówi, czy decyzja dojechała do wykonawcy.
func (a *adapterKolejek) zastosujDecyzje(ctx context.Context, pozycja dane.Pozycja,
	wstrzymanie dane.WstrzymanieKroku) (bool, string, error) {

	if wstrzymanie.Stan == stanKrokuOdrzucony {
		werdykt := werdyktOdrzucone
		if err := a.repozytorium.ZmienStanPozycji(ctx, pozycja.ID, stanPozycjiAnulowana, &werdykt); err != nil {
			return false, "", err
		}
		return false, doreczenieNiepotrzebne, nil
	}

	// Zatwierdzenie prowadzi krok dalej tabelą przejść, od stanu sprzed wstrzymania.
	docelowy, jest := krokNaprzod[wstrzymanie.StanPozycjiPrzed]
	if !jest {
		return false, "", bladSterowaniaKrokiem(shared.ErrorCodeConflict,
			"krok wstrzymany w stanie "+wstrzymanie.StanPozycjiPrzed+" nie ma dokąd pójść naprzód")
	}
	// Praca przerwana w locie wykonuje się ponownie, zamiast przeskoczyć do weryfikacji.
	if wstrzymanie.StanPozycjiPrzed == stanPozycjiWykonywana {
		docelowy = stanPozycjiWykonywana
	}
	if err := a.repozytorium.ZmienStanPozycji(ctx, pozycja.ID, docelowy, werdyktKroku(docelowy)); err != nil {
		return false, "", err
	}
	if docelowy != stanPozycjiWykonywana {
		// Krok zatwierdzony po weryfikacji jest zamknięty przyjęciem wyniku, bez wykonawcy.
		return false, doreczenieNiepotrzebne, nil
	}
	if a.silnik.wykonawca == nil {
		return false, doreczenieBezWykonawcy, nil
	}
	// Treść zlecenia niesie rozstrzygnięcie Operatora do wykonawcy tą samą drogą co zawsze.
	zPolecenia := pozycja
	zPolecenia.Stan = stanPozycjiWykonywana
	zPolecenia.TrescZlecenia = trescZDecyzja(pozycja.TrescZlecenia, wstrzymanie)
	return true, "", a.silnik.wykonaj(ctx, zPolecenia)
}

// trescZDecyzja dokleja rozstrzygnięcie Operatora do treści zlecenia kroku, wyłącznie
// do treści tej jednej tury; wiersz pozycja_kolejki zostaje przy tym nietknięty.
func trescZDecyzja(tresc *string, wstrzymanie dane.WstrzymanieKroku) *string {
	czesci := []string{strings.TrimSpace(wartoscTekstu(tresc))}
	czesci = append(czesci, "Decyzja Operatora: krok zatwierdzony do dalszego wykonania.")
	if powod := strings.TrimSpace(wartoscTekstu(wstrzymanie.Powod)); powod != "" {
		czesci = append(czesci, "Powód wstrzymania: "+powod)
	}
	if uzasadnienie := strings.TrimSpace(wartoscTekstu(wstrzymanie.Uzasadnienie)); uzasadnienie != "" {
		czesci = append(czesci, "Uzasadnienie decyzji: "+uzasadnienie)
	}
	pelna := strings.TrimSpace(strings.Join(czesci, "\n\n"))
	return &pelna
}

// odnajdzKrok rozwiązuje kolejkę i krok ze wskazań żądania. Krok nieznany
// kolejce jest odmową, nie cichym przejściem na inny — wykonanie czegoś innego
// niż zlecono byłoby gorsze od odmowy.
func (a *adapterKolejek) odnajdzKrok(ctx context.Context, idKolejki, idKroku string) (
	int64, []dane.Pozycja, *dane.Pozycja, error) {

	kolejkaID, err := identyfikatorKolejki(idKolejki)
	if err != nil {
		return 0, nil, nil, err
	}
	if _, err := a.repozytorium.PobierzKolejke(ctx, kolejkaID); err != nil {
		return 0, nil, nil, bladSterowaniaKrokiem(shared.ErrorCodeNotFound,
			"kolejka nie istnieje: "+idKolejki)
	}
	pozycje, err := a.repozytorium.ListaPozycji(ctx, kolejkaID)
	if err != nil {
		return 0, nil, nil, err
	}
	wskazanie := strings.TrimSpace(idKroku)
	if wskazanie == "" {
		return 0, nil, nil, bladSterowaniaKrokiem(shared.ErrorCodeValidationFailed,
			"sterowanie krokiem wymaga wskazania kroku — całą kolejką steruje queue.action")
	}
	krokID, err := strconv.ParseInt(wskazanie, 10, 64)
	if err == nil {
		for i := range pozycje {
			if pozycje[i].ID == krokID {
				return kolejkaID, pozycje, &pozycje[i], nil
			}
		}
	}
	return 0, nil, nil, bladSterowaniaKrokiem(shared.ErrorCodeNotFound,
		"krok nie istnieje w kolejce "+idKolejki+": "+idKroku)
}

// pozycjaPoID odnajduje krok w świeżo odczytanym wykazie; brak trafienia
// oddaje wersję sprzed zmiany, żeby odpowiedź nie zniknęła z powodu wyścigu.
func (a *adapterKolejek) pozycjaPoID(pozycje []dane.Pozycja, id int64,
	zapasowa dane.Pozycja) *dane.Pozycja {

	for i := range pozycje {
		if pozycje[i].ID == id {
			return &pozycje[i]
		}
	}
	return &zapasowa
}

// krokZleceniaKontraktu składa krok zlecenia widziany z zewnątrz; stan sterowania
// jest wyprowadzany, nie zapisywany drugi raz, więc pozycja zamknięta zawsze
// przedstawia się jako zamkniety.
func krokZleceniaKontraktu(kolejkaID int64, pozycja dane.Pozycja,
	wstrzymanie dane.WstrzymanieKroku, jestWstrzymanie bool) (krokZlecenia, error) {

	stanPracy, err := stanPracyKontraktu(pozycja.Stan)
	if err != nil {
		return krokZlecenia{}, err
	}
	// Stan sterowania liczy się z wstrzymania czynnego — zastosowany epizod jest już historią.
	stanSterowania, err := stanSterowaniaKontraktu(stanSterowaniaKrokiem(pozycja, wstrzymanie,
		jestWstrzymanie && !wstrzymanie.CzyZastosowane()))
	if err != nil {
		return krokZlecenia{}, err
	}

	krok := krokZlecenia{
		Id:            strconv.FormatInt(pozycja.ID, 10),
		QueueId:       strconv.FormatInt(kolejkaID, 10),
		Title:         pozycja.Tytul,
		Order:         pozycja.Kolejnosc,
		WorkStatus:    stanPracy,
		ControlStatus: stanSterowania,
		Cycle:         pozycja.LicznikObiegow,
		UpdatedAt:     chwilaBazy(pozycja.Zaktualizowano),
	}
	if !jestWstrzymanie {
		return krok, nil
	}
	// resumeFrom niesie wartość tego samego wyliczenia co workStatus, więc idzie tym samym przekładem.
	stanPowrotu, err := stanPracyKontraktu(wstrzymanie.StanPozycjiPrzed)
	if err != nil {
		return krokZlecenia{}, err
	}
	krok.HoldReason = wstrzymanie.Powod
	krok.DecisionNote = wstrzymanie.Uzasadnienie
	krok.ResumeFrom = &stanPowrotu
	krok.HeldAt = &wstrzymanie.WstrzymanoO
	krok.DecidedAt = wstrzymanie.ZdecydowanoO
	krok.AppliedAt = wstrzymanie.ZastosowanoO
	krok.DeliveredAt = wstrzymanie.DoreczonoO
	return krok, nil
}

// stanSterowaniaKrokiem wyprowadza stan kroku widziany przez Operatora; sprawdza
// najpierw stan końcowy pozycji, dopiero potem wstrzymanie.
func stanSterowaniaKrokiem(pozycja dane.Pozycja, wstrzymanie dane.WstrzymanieKroku,
	jestWstrzymanie bool) string {

	if czyStanKoncowyPozycji(pozycja.Stan) {
		return stanKrokuZamkniety
	}
	if jestWstrzymanie {
		return wstrzymanie.Stan
	}
	if pozycja.Stan == stanPozycjiWykonywana {
		return stanKrokuBiegnie
	}
	return stanKrokuCzeka
}

// stanDecyzjiZLadunku przekłada rozstrzygnięcie z koperty na stan zapisu.
// Wartość spoza dwóch znanych jest odmową — zgadywanie intencji Operatora przy
// decyzji o wstrzymanym kroku byłoby najgorszym możliwym miejscem na domysł.
func stanDecyzjiZLadunku(decyzja string) (string, error) {
	switch strings.ToLower(strings.TrimSpace(decyzja)) {
	case decyzjaZatwierdz:
		return stanKrokuZatwierdzony, nil
	case decyzjaOdrzuc:
		return stanKrokuOdrzucony, nil
	}
	return "", bladSterowaniaKrokiem(shared.ErrorCodeValidationFailed,
		"decyzja o kroku musi brzmieć „"+decyzjaZatwierdz+"” albo „"+decyzjaOdrzuc+"”, jest: "+decyzja)
}

// identyfikatorKolejki czyta numer kolejki ze wskazania kontraktu; wskazanie
// nieliczbowe jest odmową, bo kontrakt niesie identyfikator kolejki jako tekst.
func identyfikatorKolejki(id string) (int64, error) {
	kolejkaID, err := strconv.ParseInt(strings.TrimSpace(id), 10, 64)
	if err != nil {
		return 0, bladSterowaniaKrokiem(shared.ErrorCodeNotFound, "kolejka nie istnieje: "+id)
	}
	return kolejkaID, nil
}

// tekstNiepusty oddaje wskaźnik na tekst po obcięciu białych znaków; pusty
// zostaje pusty, żeby baza nie zapamiętała spacji jako powodu wstrzymania.
func tekstNiepusty(wartosc string) *string {
	przyciety := strings.TrimSpace(wartosc)
	if przyciety == "" {
		return nil
	}
	return &przyciety
}

// bladSterowaniaKrokiem nazywa odmowę rodziny sterowania krokiem, ze wspólnym
// kodem i opisem dla wstrzymania, decyzji i wykazu kroków tego samego zlecenia.
func bladSterowaniaKrokiem(kod protocol.KodBledu, opis string) error {
	return protocol.JakoError(protocol.NowyBlad(kod, opis))
}
