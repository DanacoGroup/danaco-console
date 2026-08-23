// Odpowiedzialność pliku: zadania pętli wykonawczej modułu Studio — rozkład
// zlecenia dokumentowego Operatora na zadania, magazyn rozkładów i złożenie
// obu bytów w kształt kontraktu (`StudioTaskPlan`, `StudioDocumentTask`).
//
// Pętla, która te zadania prowadzi, stoi w `adapter_modul_studio_petla.go`.
// Podział jest taki: tutaj mieszka to, CZYM pętla pracuje, tam to, JAK pracuje.
//
// ── Rozkład jest robotą, nie ozdobą ─────────────────────────────────────────
// Operator mówi „napisz pismo w tej sprawie na podstawie tych materiałów,
// sprawdź terminologię i przygotuj wersję do druku" i ma dostać cztery zadania,
// nie jedno. Rozkład idzie dwiema drogami po kolei: najpierw modelem (kontrakt
// mówi wprost „brak znaczy rozkład ułożony przez model"), a gdy modelu nie ma
// albo jego odpowiedź nie składa się w zadania — rozkładem własnym, po czynności
// nazwanych w zleceniu. Druga droga nie jest atrapą: rozpoznaje czynności
// z wykazu `StudioTaskKind` po słowach, którymi Operator o nich mówi, i zachowuje
// ich kolejność ze zdania. Zlecenie, w którym nie da się rozpoznać ani jednej
// czynności, daje jedno zadanie rodzaju `custom` niosące całe zlecenie — bo
// „nie rozpoznałem" nie może znaczyć „rozkład pusty".
//
// ── Magazyn jest w rdzeniu, nie w bazie — i to jest zapisane ─────────────────
// Rozkłady leżą w magazynie pamięciowym tego procesu. Warstwa danych nie ma
// dziś tabeli rozkładu ani zadania, a zakładanie jej należy do odcinka
// fundamentu, nie tutaj. Skutek jest nazwany, nie przemilczany: rozkład nie
// przeżywa ponownego uruchomienia rdzenia. Tabela `plan_studio` i
// `zadanie_studio` są wypisane w sprawozdaniu jako brak do domknięcia.
package core

import (
	"sort"
	"strings"
	"sync"
	"time"

	"danacoconsole/shared"
)

// Przedrostki identyfikatorów bytów pętli wykonawczej.
const (
	przedrostekRozkladuStudia = "studio-plan-"
	przedrostekZadaniaStudia  = "studio-zad-"
)

// petlaZadanieStudia jest zadaniem rozkładu w postaci magazynu. Kształt
// kontraktu składa `petlaZlozZadanie`; tutaj leżą te same pola plus wskazanie
// okna, którym zadanie jedzie do modelu.
type petlaZadanieStudia struct {
	Kod           string
	KodRozkladu   string
	KodDokumentu  string
	Kolejnosc     int
	Rodzaj        shared.StudioTaskKind
	Nazwa         string
	Polecenie     string
	ZakresOd      *int
	ZakresDo      *int
	Stan          shared.StudioTaskState
	Wykonawca     *shared.StudioActor
	StoiNa        []string
	Wynik         *string
	PowodPorazki  *string
	KodyCzynnosci []string
	Podjeto       *int64
	Domknieto     *int64
}

// petlaRozkladStudia jest rozkładem zlecenia w postaci magazynu.
type petlaRozkladStudia struct {
	Kod              string
	KodDokumentu     string
	Zlecenie         string
	Stan             shared.StudioPlanState
	Zadania          []*petlaZadanieStudia
	Ulozyl           *shared.StudioActor
	Obiegi           int
	ObiegiBezPostepu int
	PowodZatrzymania *string
	Utworzono        int64
	Zaktualizowano   int64
	// zatrzymanie stoi na `true` od chwili, w której Operator zawołał
	// `studio.plan.stop`. Pętla sprawdza je PRZED każdym zadaniem i po każdym
	// obiegu, więc zatrzymanie naprawdę zatrzymuje, a nie tylko przestawia
	// napis w oknie.
	zatrzymanie bool
}

// petlaMagazynStudia trzyma rozkłady procesu. Zamek jest jeden na magazyn, bo
// pętla czyta i zapisuje ten sam rozkład w jednym obiegu, a dwa zamki na tym
// samym bycie znaczyłyby dwie prawdy o jego stanie.
type petlaMagazynStudia struct {
	zamek     sync.Mutex
	rozklady  map[string]*petlaRozkladStudia
	zadania   map[string]*petlaZadanieStudia
	kolejnosc []string
}

// petlaNowyMagazyn zakłada pusty magazyn rozkładów.
func petlaNowyMagazyn() *petlaMagazynStudia {
	return &petlaMagazynStudia{
		rozklady: make(map[string]*petlaRozkladStudia),
		zadania:  make(map[string]*petlaZadanieStudia),
	}
}

// petlaDolozRozklad wnosi rozkład do magazynu wraz z jego zadaniami.
func (m *petlaMagazynStudia) petlaDolozRozklad(rozklad *petlaRozkladStudia) {
	if m == nil || rozklad == nil {
		return
	}
	m.zamek.Lock()
	defer m.zamek.Unlock()
	m.rozklady[rozklad.Kod] = rozklad
	m.kolejnosc = append(m.kolejnosc, rozklad.Kod)
	for _, zadanie := range rozklad.Zadania {
		m.zadania[zadanie.Kod] = zadanie
	}
}

// petlaRozkladPoKodzie oddaje rozkład wskazany kodem.
func (m *petlaMagazynStudia) petlaRozkladPoKodzie(kod string) (*petlaRozkladStudia, bool) {
	if m == nil {
		return nil, false
	}
	m.zamek.Lock()
	defer m.zamek.Unlock()
	rozklad, jest := m.rozklady[kod]
	return rozklad, jest
}

// petlaZadaniePoKodzie oddaje zadanie wskazane kodem wraz z jego rozkładem.
func (m *petlaMagazynStudia) petlaZadaniePoKodzie(kod string) (*petlaZadanieStudia, *petlaRozkladStudia, bool) {
	if m == nil {
		return nil, nil, false
	}
	m.zamek.Lock()
	zadanie, jest := m.zadania[kod]
	m.zamek.Unlock()
	if !jest {
		return nil, nil, false
	}
	rozklad, maRozklad := m.petlaRozkladPoKodzie(zadanie.KodRozkladu)
	if !maRozklad {
		return nil, nil, false
	}
	return zadanie, rozklad, true
}

// petlaRozkladyDokumentu oddaje rozkłady dokumentu, od najświeższego.
func (m *petlaMagazynStudia) petlaRozkladyDokumentu(kodDokumentu string) []*petlaRozkladStudia {
	if m == nil {
		return nil
	}
	m.zamek.Lock()
	defer m.zamek.Unlock()
	znalezione := make([]*petlaRozkladStudia, 0, len(m.kolejnosc))
	for _, kod := range m.kolejnosc {
		rozklad, jest := m.rozklady[kod]
		if !jest {
			continue
		}
		if kodDokumentu != "" && rozklad.KodDokumentu != kodDokumentu {
			continue
		}
		znalezione = append(znalezione, rozklad)
	}
	sort.SliceStable(znalezione, func(i, j int) bool {
		return znalezione[i].Utworzono > znalezione[j].Utworzono
	})
	return znalezione
}

// ── Rozkład zlecenia na zadania ─────────────────────────────────────────────

// petlaZnacznikCzynnosci wiąże słowa, którymi Operator nazywa czynność, z jej
// rodzajem z wykazu kontraktu. Wykaz jest po polsku, bo zlecenie przychodzi po
// polsku; dopisanie wariantu jest dopisaniem wiersza, nie zmianą rozkładu.
//
// Kolejność w tablicy nie ma znaczenia — kolejność zadań bierze się z miejsca,
// w którym słowo stoi w zleceniu, bo Operator wymienia czynności w tej
// kolejności, w jakiej mają się wykonać.
var petlaZnacznikiCzynnosci = []struct {
	slowa  []string
	rodzaj shared.StudioTaskKind
	nazwa  string
}{
	{
		slowa:  []string{"na podstawie", "zbierz", "materiał", "materiały", "materialy", "źródł", "zrodl", "poszukaj", "sprawdź w źródłach"},
		rodzaj: shared.StudioTaskKindResearch,
		nazwa:  "Zebranie źródeł i materiału",
	},
	{
		slowa:  []string{"napisz", "przygotuj pismo", "sporządź", "sporzadz", "zredaguj", "ułóż treść", "uloz tresc", "brzmienie"},
		rodzaj: shared.StudioTaskKindDraft,
		nazwa:  "Napisanie brzmienia",
	},
	{
		slowa:  []string{"sformatuj", "formatowanie", "postać", "postac", "styl", "układ", "uklad", "tabel", "marginesy", "do druku", "wersję do druku", "wersje do druku"},
		rodzaj: shared.StudioTaskKindFormat,
		nazwa:  "Doprowadzenie postaci dokumentu",
	},
	{
		slowa:  []string{"spis treści", "spis tresci", "przypis", "bibliograf", "powołani", "powolani", "indeks", "odwołani", "odwolani", "aparat"},
		rodzaj: shared.StudioTaskKindApparatus,
		nazwa:  "Aparat dokumentu",
	},
	{
		slowa:  []string{"przejrzyj", "przejrzenie", "oceń", "ocen", "zaznacz", "znakuj", "recenzj", "opini"},
		rodzaj: shared.StudioTaskKindReview,
		nazwa:  "Przejrzenie i znakowanie",
	},
	{
		slowa:  []string{"terminolog", "korekt", "popraw język", "popraw jezyk", "ortograf", "interpunkc", "sprawdź pisownię", "sprawdz pisownie", "sprawdź terminologię", "sprawdz terminologie"},
		rodzaj: shared.StudioTaskKindProofread,
		nazwa:  "Poprawa językowa i terminologiczna",
	},
	{
		slowa:  []string{"wydaj", "wydanie", "eksport", "zapisz jako pdf", "pdf", "docx", "odt", "wyślij do druku", "wyslij do druku"},
		rodzaj: shared.StudioTaskKindExport,
		nazwa:  "Wydanie do formatu",
	},
}

// petlaRozlozZlecenie rozkłada zlecenie Operatora na zadania rozpoznaniem
// czynności nazwanych w jego treści.
//
// Rozpoznanie idzie po POŁOŻENIU słowa w zleceniu, nie po kolejności wykazu:
// „sprawdź terminologię i napisz pismo" ma dać korektę po napisaniu tylko wtedy,
// gdy Operator tak powiedział. Rodzaj rozpoznany dwa razy wchodzi raz — dwa
// zadania tej samej czynności z jednego zdania byłyby zdublowaną robotą.
//
// Zależności między zadaniami są łańcuchem: zadanie stoi na poprzednim. Rozkład
// zlecenia dokumentowego jest z natury szeregowy — nie da się poprawić języka
// pisma, którego jeszcze nie ma. Rozkład podany wprost przez Operatora albo
// przez model swoje zależności niesie własne i ich nie nadpisujemy.
func petlaRozlozZlecenie(zlecenie string) []petlaZadanieStudia {
	tresc := strings.ToLower(zlecenie)
	type trafienie struct {
		miejsce int
		rodzaj  shared.StudioTaskKind
		nazwa   string
	}
	trafienia := []trafienie{}
	widziane := map[shared.StudioTaskKind]bool{}

	for _, znacznik := range petlaZnacznikiCzynnosci {
		najblizsze := -1
		for _, slowo := range znacznik.slowa {
			if miejsce := strings.Index(tresc, slowo); miejsce >= 0 {
				if najblizsze < 0 || miejsce < najblizsze {
					najblizsze = miejsce
				}
			}
		}
		if najblizsze < 0 || widziane[znacznik.rodzaj] {
			continue
		}
		widziane[znacznik.rodzaj] = true
		trafienia = append(trafienia, trafienie{
			miejsce: najblizsze, rodzaj: znacznik.rodzaj, nazwa: znacznik.nazwa,
		})
	}
	sort.SliceStable(trafienia, func(i, j int) bool {
		return trafienia[i].miejsce < trafienia[j].miejsce
	})

	if len(trafienia) == 0 {
		// Zlecenie bez rozpoznanej czynności NIE daje rozkładu pustego: daje
		// jedno zadanie niosące zlecenie w całości. Pusty rozkład byłby
		// odpowiedzią udaną bez treści, a takiej w tym produkcie nie ma.
		return []petlaZadanieStudia{{
			Kolejnosc: 1,
			Rodzaj:    shared.StudioTaskKindCustom,
			Nazwa:     petlaSkrocZlecenie(zlecenie),
			Polecenie: zlecenie,
			Stan:      shared.StudioTaskStatePending,
		}}
	}

	zadania := make([]petlaZadanieStudia, 0, len(trafienia))
	for numer, t := range trafienia {
		zadanie := petlaZadanieStudia{
			Kolejnosc: numer + 1,
			Rodzaj:    t.rodzaj,
			Nazwa:     t.nazwa,
			Polecenie: zlecenie,
			Stan:      shared.StudioTaskStatePending,
		}
		zadania = append(zadania, zadanie)
	}
	return zadania
}

// petlaSkrocZlecenie oddaje zlecenie skrócone do nazwy zadania. Nazwa jest
// jednym zdaniem, tak jak każe kontrakt (`StudioDocumentTask.title`).
func petlaSkrocZlecenie(zlecenie string) string {
	tresc := strings.TrimSpace(zlecenie)
	if tresc == "" {
		return "Zlecenie bez treści"
	}
	runy := []rune(tresc)
	const granica = 80
	if len(runy) <= granica {
		return tresc
	}
	return strings.TrimSpace(string(runy[:granica])) + "…"
}

// petlaZlozZadania nadaje zadaniom kody, wiąże je z rozkładem i układa łańcuch
// zależności tam, gdzie zadanie własnej zależności nie niesie.
func petlaZlozZadania(kodRozkladu, kodDokumentu string,
	surowe []petlaZadanieStudia) []*petlaZadanieStudia {

	zlozone := make([]*petlaZadanieStudia, 0, len(surowe))
	poprzedni := ""
	for numer := range surowe {
		zadanie := surowe[numer]
		zadanie.Kod = nowyIdentyfikator(przedrostekZadaniaStudia)
		zadanie.KodRozkladu = kodRozkladu
		zadanie.KodDokumentu = kodDokumentu
		if zadanie.Kolejnosc == 0 {
			zadanie.Kolejnosc = numer + 1
		}
		if zadanie.Stan == "" {
			zadanie.Stan = shared.StudioTaskStatePending
		}
		if zadanie.Rodzaj == "" {
			zadanie.Rodzaj = shared.StudioTaskKindCustom
		}
		if strings.TrimSpace(zadanie.Nazwa) == "" {
			zadanie.Nazwa = string(zadanie.Rodzaj)
		}
		if len(zadanie.StoiNa) == 0 && poprzedni != "" {
			zadanie.StoiNa = []string{poprzedni}
		}
		poprzedni = zadanie.Kod
		zlozone = append(zlozone, &zadanie)
	}
	return zlozone
}

// ── Złożenie w kształt kontraktu ────────────────────────────────────────────

// petlaZlozZadanie oddaje zadanie w kształcie kontraktu.
func petlaZlozZadanie(zadanie *petlaZadanieStudia) shared.StudioDocumentTask {
	if zadanie == nil {
		return shared.StudioDocumentTask{}
	}
	kształt := shared.StudioDocumentTask{
		Id:         zadanie.Kod,
		PlanId:     zadanie.KodRozkladu,
		DocumentId: zadanie.KodDokumentu,
		Order:      zadanie.Kolejnosc,
		Kind:       zadanie.Rodzaj,
		Title:      zadanie.Nazwa,
		State:      zadanie.Stan,
		RangeStart: zadanie.ZakresOd,
		RangeEnd:   zadanie.ZakresDo,
		AssignedTo: zadanie.Wykonawca,
		Result:     zadanie.Wynik,
		StartedAt:  zadanie.Podjeto,
		FinishedAt: zadanie.Domknieto,
	}
	if zadanie.Polecenie != "" {
		polecenie := zadanie.Polecenie
		kształt.Instruction = &polecenie
	}
	if zadanie.PowodPorazki != nil {
		kształt.FailureReason = zadanie.PowodPorazki
	}
	if len(zadanie.StoiNa) > 0 {
		kształt.DependsOn = append([]string(nil), zadanie.StoiNa...)
	}
	if len(zadanie.KodyCzynnosci) > 0 {
		kształt.ActionIds = append([]string(nil), zadanie.KodyCzynnosci...)
	}
	return kształt
}

// petlaZlozRozklad oddaje rozkład w kształcie kontraktu, z zadaniami
// zawężonymi do wskazanego stanu, gdy Operator o zawężenie poprosił.
func petlaZlozRozklad(rozklad *petlaRozkladStudia,
	stan *shared.StudioTaskState) shared.StudioTaskPlan {

	if rozklad == nil {
		return shared.StudioTaskPlan{}
	}
	obiegi, bezPostepu := rozklad.Obiegi, rozklad.ObiegiBezPostepu
	utworzono, zaktualizowano := rozklad.Utworzono, rozklad.Zaktualizowano
	kształt := shared.StudioTaskPlan{
		Id:                        rozklad.Kod,
		DocumentId:                rozklad.KodDokumentu,
		Order:                     rozklad.Zlecenie,
		State:                     rozklad.Stan,
		CreatedBy:                 rozklad.Ulozyl,
		Iterations:                &obiegi,
		IterationsWithoutProgress: &bezPostepu,
		StopReason:                rozklad.PowodZatrzymania,
		CreatedAt:                 &utworzono,
		UpdatedAt:                 &zaktualizowano,
	}
	zadania := make([]shared.StudioDocumentTask, 0, len(rozklad.Zadania))
	for _, zadanie := range rozklad.Zadania {
		if stan != nil && zadanie.Stan != *stan {
			continue
		}
		zadania = append(zadania, petlaZlozZadanie(zadanie))
	}
	kształt.Tasks = zadania
	return kształt
}

// petlaTeraz oddaje chwilę bieżącą w milisekundach epoki — jednostce, którą
// kontrakt nazywa w polach czasu rozkładu i zadania.
func petlaTeraz() int64 {
	return time.Now().UnixMilli()
}

// petlaWykonawcaZlecenia składa tożsamość wykonawcy z pól żądania. Rodzaj
// autora jest `model`, bo rozkład układa wykonawca, nie Operator — a rozkład
// bez zapisanego autora nie dałby się odróżnić w podświetleniu zmian modelu.
func petlaWykonawcaZlecenia(kodEksperta, nazwaEksperta, okno *string) *shared.StudioActor {
	wykonawca := shared.StudioActor{Kind: shared.StudioAuthorModel}
	pusty := true
	if kodEksperta != nil && strings.TrimSpace(*kodEksperta) != "" {
		wykonawca.AgentId = kodEksperta
		pusty = false
	}
	if nazwaEksperta != nil && strings.TrimSpace(*nazwaEksperta) != "" {
		wykonawca.AgentName = nazwaEksperta
		pusty = false
	}
	if okno != nil && strings.TrimSpace(*okno) != "" {
		wykonawca.WindowId = okno
		pusty = false
	}
	if pusty {
		// Wykonawca bez ani jednego wskazania nadal jest wykonawcą — rodzaj
		// autora niesie prawdę o tym, że rozkładu nie ułożył Operator.
		return &wykonawca
	}
	return &wykonawca
}
