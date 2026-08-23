// Odpowiedzialność pliku: wpięcie trzech komend historii rozmowy okna —
// `history.load`, `history.delete`, `retention.set` — wraz z portem Historia,
// jego wypełnieniem repozytorium i rozgłoszeniem `history.changed`.
//
// Trzy komendy, nie cztery. `history.changed` jest zdarzeniem, nie komendą:
// stoi w dziale zdarzeń kontraktu i niesie `payload`, nie parę żądanie/wynik.
// Rejestr komend rośnie o trzy.
//
// Historia nie jest nowym bytem. Pozycją historii jest wypowiedź okna —
// wiersz `wiadomosc`. Rodzina `history.*` daje jej drugi widok
// (od najnowszej, kursorem czasu, z możliwością skasowania), a nie drugą
// tabelę; wypełnia ją repozytorium `dane/historia.go`. Nowa tabela wchodzi
// wyłącznie pod zasadę przechowywania,
// bo tego bytu w schemacie nie było wcale.
//
// Retencja jest egzekwowana, nie tylko zapisana. Zasada zapisana, której nikt
// nie stosuje, byłaby atrapą: Operator widziałby nastawę i rosnącą
// historię naraz. Egzekucja ma tu dwa miejsca wpięcia, oba na drodze, którą
// rdzeń i tak przechodzi:
//
//	`retention.set`  — zaraz po zapisaniu zasady, na wszystkich oknach zakresu.
//	                   Nastawa działa od chwili ustawienia, a nie od następnego
//	                   restartu.
//	`history.load`   — przed oddaniem wykazu, na czytanym oknie. Wykaz nigdy nie
//	                   pokazuje pozycji, która zasadzie już nie podlega.
//
// Czego te dwa miejsca nie dają. Okno, którego nikt nie czyta i na którym
// nikt nie przestawia zasady, nie zostanie przycięte samo: przemiatania
// okresowego ani przy starcie rdzenia tu nie ma. Wpięcie takiego przemiatania
// siedzi poza tym pakietem (montaż rdzenia, wzorem czyszczenia kosza sesji).
//
// Przy czyszczeniu wielu pozycji zdarzenie nie niesie pozycji. Kontrakt daje
// `entry` jako pole niewymagane właśnie na ten przypadek: po skasowaniu setki
// wypowiedzi nie ma jednej pozycji „po zmianie", a wysyłanie stu zdarzeń
// zamieniłoby odświeżenie wykazu w burzę. Idzie jedno zdarzenie z `windowId`
// i rodzajem `deleted` — okno przeładowuje wykaz.
package core

import (
	"context"
	"strings"

	"danacoconsole/server/internal/dane"
	"danacoconsole/server/internal/protocol"
	"danacoconsole/shared"
)

// Historia jest portem historii rozmowy okna i jej retencji.
type Historia interface {
	Wczytaj(ctx context.Context, z shared.HistoryLoadRequest) (shared.HistoryLoadResponse, error)
	Usun(ctx context.Context, z shared.HistoryDeleteRequest) (shared.HistoryDeleteResponse, error)
	UstawZasade(ctx context.Context, z shared.RetentionSetRequest) (shared.RetentionSetResponse, error)
	// OknaZasady oddaje okna, na których zasada właśnie zadziałała — rdzeń
	// rozgłasza dla nich `history.changed`. Bez tego Operator ustawiałby
	// retencję i patrzył na wykaz nieodświeżony do czasu ponownego wejścia.
	OknaZasady(ctx context.Context, z shared.RetentionSetRequest) []string
}

// zarejestrujHistorie wpina trzy komendy rodziny historii.
func zarejestrujHistorie(r *Rejestr, h Historia, e *emiter) {
	if r == nil || h == nil {
		return
	}

	r.Zarejestruj(shared.CommandHistoryLoad, obsluz(h.Wczytaj))

	r.Zarejestruj(shared.CommandHistoryDelete,
		obsluz(func(ctx context.Context, z shared.HistoryDeleteRequest) (shared.HistoryDeleteResponse, error) {
			w, err := h.Usun(ctx, z)
			// Zdarzenie idzie wyłącznie po czynności, która coś zmieniła:
			// usunięcie zera pozycji nie jest zmianą historii.
			if err == nil && w.Deleted > 0 {
				e.historia(shared.ChangeKindDeleted, w.WindowId, nil)
			}
			return w, err
		}))

	r.Zarejestruj(shared.CommandRetentionSet,
		obsluz(func(ctx context.Context, z shared.RetentionSetRequest) (shared.RetentionSetResponse, error) {
			w, err := h.UstawZasade(ctx, z)
			if err != nil {
				return w, err
			}
			for _, oknoKod := range h.OknaZasady(ctx, z) {
				e.historia(shared.ChangeKindDeleted, oknoKod, nil)
			}
			return w, nil
		}))
}

// historia rozgłasza zmianę historii rozmowy okna. Sesji komunikatu rdzeń tu
// nie wyznacza: zmiana dotyczy okna, a jedno rozgłoszenie potrafi
// objąć okna wielu sesji naraz — zasada zakresu globalnego tnie je wszystkie.
func (e *emiter) historia(zmiana shared.ChangeKind, oknoKod string, pozycja *shared.HistoryEntry) {
	e.wyslij(shared.EventHistoryChanged, "",
		shared.HistoryChangedEvent{Change: zmiana, WindowId: oknoKod, Entry: pozycja})
}

// adapterHistorii wypełnia port Historia repozytorium historii i retencji.
// Stoi w tym samym pliku co port — jest jego jedynym wypełnieniem i niesie
// wyłącznie przekład kontraktu, wzorem `handlers_agent_wersje.go`.
type adapterHistorii struct {
	historia dane.RepozytoriumHistorii
}

var _ Historia = (*adapterHistorii)(nil)

// NowyPortHistorii wiąże port z repozytorium historii rozmowy okna.
func NowyPortHistorii(historia dane.RepozytoriumHistorii) *adapterHistorii {
	return &adapterHistorii{historia: historia}
}

// dlugoscSkrotuHistorii jest progiem skrótu treści w wykazie. Skrót składa się
// tu, a nie w bazie: kontrakt nazywa go `preview` i to kontrakt rozstrzyga, ile
// go jest — kolumna z przyciętą treścią byłaby drugą prawdą o wypowiedzi.
const dlugoscSkrotuHistorii = 200

// Wczytaj oddaje wykaz historii okna od pozycji najnowszej, stronicowany
// kursorem czasu. Przed odczytem stosuje zasadę przechowywania — wykaz nie
// pokazuje pozycji, która zasadzie już nie podlega.
//
// Egzekucja nie przerywa odczytu. Gdyby czyszczenie zawiodło, Operator ma
// dostać historię, jaka jest, zamiast odmowy odczytu — niepowodzenie
// zabiera tę jedną czynność, nie całą komendę.
func (a *adapterHistorii) Wczytaj(ctx context.Context,
	z shared.HistoryLoadRequest) (shared.HistoryLoadResponse, error) {

	if a == nil || a.historia == nil {
		return shared.HistoryLoadResponse{WindowId: z.WindowId, Entries: []shared.HistoryEntry{}}, nil
	}
	if z.WindowId == "" {
		return shared.HistoryLoadResponse{}, bladZadaniaHistorii("wskazanie okna jest puste")
	}
	_, _ = a.historia.Egzekwuj(ctx, z.WindowId)

	kursor := int64(0)
	if z.Before != nil {
		kursor = *z.Before
	}
	limit := 0
	if z.Limit != nil && *z.Limit > 0 {
		limit = *z.Limit
	}
	wiersze, err := a.historia.Pozycje(ctx, z.WindowId, kursor, limit)
	if err != nil {
		return shared.HistoryLoadResponse{}, err
	}
	// `total` jest rozmiarem całej historii okna, nie strony i nie wycinka za
	// kursorem: panel podstawia tę liczbę pod ostrzeżenie „ile zniknie" przy
	// czyszczeniu, a `history.delete` bez wskazania pozycji kasuje okno w
	// całości. Liczba zawężona kursorem malała przy każdym dociągnięciu strony
	// i obiecywała utratę mniejszą niż rzeczywista.
	razem, err := a.historia.Policz(ctx, z.WindowId)
	if err != nil {
		return shared.HistoryLoadResponse{}, err
	}
	pozycje := make([]shared.HistoryEntry, 0, len(wiersze))
	for _, wiersz := range wiersze {
		pozycje = append(pozycje, pozycjaHistoriiKontraktu(wiersz))
	}
	return shared.HistoryLoadResponse{WindowId: z.WindowId, Entries: pozycje, Total: razem}, nil
}

// Usun kasuje wskazane pozycje historii okna, a bez wskazania — całą historię
// okna. Oddaje liczbę usuniętych pozycji: zero znaczy „nie było czego usuwać"
// i nie jest odmową.
func (a *adapterHistorii) Usun(ctx context.Context,
	z shared.HistoryDeleteRequest) (shared.HistoryDeleteResponse, error) {

	if a == nil || a.historia == nil {
		return shared.HistoryDeleteResponse{}, bladBrakuKatalogu("historii rozmowy")
	}
	if z.WindowId == "" {
		return shared.HistoryDeleteResponse{}, bladZadaniaHistorii("wskazanie okna jest puste")
	}
	usuniete, err := a.historia.Usun(ctx, z.WindowId, z.EntryIds)
	if err != nil {
		return shared.HistoryDeleteResponse{}, err
	}
	return shared.HistoryDeleteResponse{WindowId: z.WindowId, Deleted: usuniete}, nil
}

// UstawZasade zapisuje zasadę przechowywania zakresu i oddaje ją w kształcie
// obowiązującym po zmianie. Egzekucję na oknach zakresu wykonuje OknaZasady,
// wołane zaraz po zapisie przez wpięcie komendy.
func (a *adapterHistorii) UstawZasade(ctx context.Context,
	z shared.RetentionSetRequest) (shared.RetentionSetResponse, error) {

	if a == nil || a.historia == nil {
		return shared.RetentionSetResponse{}, bladBrakuKatalogu("zasad przechowywania")
	}
	if err := sprawdzZakresZasady(z); err != nil {
		return shared.RetentionSetResponse{}, err
	}
	kod := ""
	if z.ScopeId != nil {
		kod = *z.ScopeId
	}
	// Zasada dla bytu, którego nie ma, jest nastawą bez skutku. Zakres okna
	// nie pyta bazy przy egzekucji (`OknaZakresu` oddaje samo wskazanie), więc
	// zasada zapisana na nieistniejące okno albo sesję leżałaby w tabeli i nie
	// przycięła nigdy niczego, a Operator miałby na nią potwierdzenie „ok".
	// Odmowa nazywa brak — czego w bazie nie ma — a nie zakaz.
	jest, err := a.historia.IstniejeByt(ctx, z.Scope, kod)
	if err != nil {
		return shared.RetentionSetResponse{}, err
	}
	if !jest {
		nazwa := "okna"
		if z.Scope == "session" {
			nazwa = "sesji"
		}
		return shared.RetentionSetResponse{}, bladZadaniaHistorii(
			"nie ma " + nazwa + " o wskazaniu " + kod + " — zasada przechowywania zapisana na byt " +
				"nieistniejący nigdy niczego nie przytnie; naprawa: sprawdzić wskazanie " +
				"wykazem (window.list albo session.list)")
	}
	zapisana, err := a.historia.ZapiszZasade(ctx, dane.ZasadaPrzechowywania{
		Zakres:          z.Scope,
		ZakresKod:       kod,
		DniTrzymania:    z.KeepDays,
		PozycjeTrzymane: z.KeepEntries,
	})
	if err != nil {
		return shared.RetentionSetResponse{}, err
	}
	return shared.RetentionSetResponse{Policy: zasadaKontraktu(zapisana)}, nil
}

// OknaZasady stosuje świeżo zapisaną zasadę na oknach jej zakresu i oddaje te,
// na których coś odpadło. Okno, na którym zasada nic nie zmieniła, nie trafia
// do wykazu — zdarzenie po czynności bez skutku byłoby fałszywym
// powiadomieniem.
func (a *adapterHistorii) OknaZasady(ctx context.Context, z shared.RetentionSetRequest) []string {
	if a == nil || a.historia == nil {
		return nil
	}
	kod := ""
	if z.ScopeId != nil {
		kod = *z.ScopeId
	}
	okna, err := a.historia.OknaZakresu(ctx, z.Scope, kod)
	if err != nil {
		return nil
	}
	zmienione := []string{}
	for _, oknoKod := range okna {
		usuniete, err := a.historia.Egzekwuj(ctx, oknoKod)
		if err == nil && usuniete > 0 {
			zmienione = append(zmienione, oknoKod)
		}
	}
	return zmienione
}

// pozycjaHistoriiKontraktu przekłada wiersz historii na pozycję kontraktu.
func pozycjaHistoriiKontraktu(p dane.PozycjaHistorii) shared.HistoryEntry {
	return shared.HistoryEntry{
		Id:        p.Identyfikator,
		WindowId:  p.OknoKod,
		SessionId: wskaznikTekstu(p.SesjaKod),
		Role:      string(p.Rola),
		Preview:   wskaznikTekstu(skrotTresci(p.Tresc)),
		CreatedAt: p.Chwila,
	}
}

// zasadaKontraktu przekłada wiersz zasady na strukturę RetentionPolicy.
func zasadaKontraktu(z dane.ZasadaPrzechowywania) shared.RetentionPolicy {
	return shared.RetentionPolicy{
		Scope:       z.Zakres,
		ScopeId:     wskaznikTekstu(z.ZakresKod),
		KeepDays:    z.DniTrzymania,
		KeepEntries: z.PozycjeTrzymane,
	}
}

// skrotTresci przycina treść do długości skrótu wykazu. Cięcie idzie po
// znakach, nie po bajtach — wypowiedź polska ucięta w środku znaku
// wielobajtowego wracałaby z rombem zamiast litery.
func skrotTresci(tresc string) string {
	tresc = strings.TrimSpace(tresc)
	znaki := []rune(tresc)
	if len(znaki) <= dlugoscSkrotuHistorii {
		return tresc
	}
	return string(znaki[:dlugoscSkrotuHistorii]) + "…"
}

// zakresyZasady wylicza trzy zakresy zasady przechowywania. Wartości są
// kontraktowe (`RetentionPolicy.scope`), a kontrakt nie daje im ani typu
// wyliczeniowego, ani polskiego odwzorowania — dlatego sprawdzenie stoi tutaj,
// przy jedynym czytelniku, a nie jako kopia słownika.
var zakresyZasady = map[string]bool{"window": true, "session": true, "global": true}

// sprawdzZakresZasady odmawia zapisu zasady, której zakresu schemat nie uniesie.
// Odmowa jest merytoryczna i dotyczy jednego wywołania — bez niej
// warunek CHECK schematu oddawałby Operatorowi błąd bazy zamiast powodu.
func sprawdzZakresZasady(z shared.RetentionSetRequest) error {
	if !zakresyZasady[z.Scope] {
		return bladZadaniaHistorii("zakres " + z.Scope + " nie należy do kontraktu")
	}
	pustyKod := z.ScopeId == nil || *z.ScopeId == ""
	if z.Scope == "global" && !pustyKod {
		return bladZadaniaHistorii("zakres global nie wskazuje bytu")
	}
	if z.Scope != "global" && pustyKod {
		return bladZadaniaHistorii("zakres " + z.Scope + " wymaga wskazania bytu")
	}
	if (z.KeepDays != nil && *z.KeepDays <= 0) || (z.KeepEntries != nil && *z.KeepEntries <= 0) {
		return bladZadaniaHistorii("próg przechowywania musi być większy od zera; " +
			"brak progu znaczy brak ograniczenia")
	}
	return nil
}

// bladZadaniaHistorii odmawia żądaniu, którego rodzina historii nie uniesie.
func bladZadaniaHistorii(powod string) error {
	return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeValidationFailed,
		"historia rozmowy: "+powod))
}
