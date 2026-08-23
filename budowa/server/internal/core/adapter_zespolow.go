// Odpowiedzialność pliku: wypełnienie portu Zespoly trwałością zespołów
// ekspertów z bazy — zapis, wczytanie, wykaz i kopia.
//
// Zespół jest składem, a nie drugą biblioteką ekspertów: adapter przenosi
// wyłącznie kody, którymi posługuje się kontrakt (`Agent.id` = `agent.kod`),
// i nie sprawdza tożsamości eksperta. Dostępność eksperta ustala warstwa
// `dane` jednym złączeniem.
//
// Brak eksperta nie przerywa wczytania: ekspert usunięty albo zarchiwizowany
// wypada ze składu oddawanego klientowi, a zespół wraca kompletny w pozostałej
// części. Fakt pominięcia trafia do dziennika rdzenia, ponieważ struktura
// `Team` kontraktu nie ma pola na taką wiadomość, a milczenie ukryłoby zmianę
// składu.
//
// Kopia nie zmienia źródła: `team.duplicate` czyta zespół źródłowy i zakłada
// obok nowy wiersz z tym samym składem; źródło nie przechodzi ani przez
// `Zapisz`, ani przez żadne inne polecenie zmieniające.
package core

import (
	"context"
	"log"
	"strings"

	"danacoconsole/server/internal/dane"
	"danacoconsole/server/internal/protocol"
	"danacoconsole/shared"
)

// Zgodność adaptera z portem sprawdzana jest przy kompilacji.
var _ Zespoly = (*adapterZespolow)(nil)

// przedrostekZespolu znakuje kod zespołu nadany przez rdzeń.
const przedrostekZespolu = "zsp-"

// przyrostekKopiiZespolu dopisuje się do nazwy źródła, gdy `team.duplicate`
// przychodzi bez nazwy — tak stanowi kontrakt.
const przyrostekKopiiZespolu = " (kopia)"

// adapterZespolow wypełnia port Zespoly tabelą `zespol`.
type adapterZespolow struct {
	repozytorium dane.RepozytoriumZespolow
	// dziennik przyjmuje wiadomość o składzie okrojonym z ekspertów, których
	// już nie ma. Nil znaczy „montaż nie wpiął dziennika" — zespoły pracują
	// wtedy bez zmiany, milknie wyłącznie ta jedna wiadomość.
	dziennik *log.Logger
}

// nowyAdapterZespolow wiąże port z repozytorium zespołów.
func nowyAdapterZespolow(repozytorium dane.RepozytoriumZespolow) *adapterZespolow {
	return &adapterZespolow{repozytorium: repozytorium}
}

// ZDziennikiem wpina dziennik rdzenia, którym adapter mówi o pominiętych
// ekspertach.
func (a *adapterZespolow) ZDziennikiem(dziennik *log.Logger) *adapterZespolow {
	a.dziennik = dziennik
	return a
}

// Zapisz zakłada zespół albo zmienia istniejący. Brak identyfikatora znaczy
// „nowy" — tak brzmi kontrakt `team.save` i tak zachowuje się okno składu.
func (a *adapterZespolow) Zapisz(ctx context.Context,
	z shared.TeamSaveRequest) (shared.TeamSaveResponse, error) {

	if a == nil || a.repozytorium == nil {
		return shared.TeamSaveResponse{}, bladBrakuZespolow(shared.CommandTeamSave)
	}
	nazwa := strings.TrimSpace(z.Name)
	if nazwa == "" {
		return shared.TeamSaveResponse{}, bladZadaniaZespolu("zespół wymaga nazwy")
	}
	// Skład jest polem wymaganym (kontrakt `team.save`, pole `agentIds`). Brak
	// pola nie jest wyborem składu pustego, tylko brakiem w żądaniu, więc wraca
	// odmową — inaczej sama zmiana nazwy kasowałaby cały skład. Pole wypełnione
	// wykazem pustym (`"agentIds": []`) jest czynnością poprawną i skład opróżnia.
	if z.AgentIds == nil {
		return shared.TeamSaveResponse{}, bladZadaniaZespolu(
			"żądanie zapisu zespołu bez pola `agentIds`; Operator poda skład zespołu — " +
				"wykaz pusty opróżnia skład, brak pola nie znaczy nic")
	}
	wpisywany := dane.Zespol{
		Nazwa: nazwa,
		Opis:  wartoscTekstu(z.Description),
		Sklad: z.AgentIds,
	}
	if z.TeamId == nil || strings.TrimSpace(*z.TeamId) == "" {
		wpisywany.Kod = nowyIdentyfikator(przedrostekZespolu)
		zapisany, err := a.repozytorium.Dodaj(ctx, wpisywany)
		if err != nil {
			return shared.TeamSaveResponse{}, err
		}
		return shared.TeamSaveResponse{Team: a.zespolKontraktu(zapisany)}, nil
	}
	wpisywany.Kod = strings.TrimSpace(*z.TeamId)
	zastany, err := a.repozytorium.PoKodzie(ctx, wpisywany.Kod)
	if err != nil {
		return shared.TeamSaveResponse{}, bladWskazania(err, "zespół", wpisywany.Kod)
	}
	a.dolozStanZastany(&wpisywany, zastany, z.Description)
	zapisany, err := a.repozytorium.Zapisz(ctx, wpisywany)
	if err != nil {
		return shared.TeamSaveResponse{}, bladWskazania(err, "zespół", wpisywany.Kod)
	}
	return shared.TeamSaveResponse{Team: a.zespolKontraktu(zapisany)}, nil
}

// dolozStanZastany uzupełnia zapis zespołu o to, czego żądanie nieść nie mogło.
//
// Opis: `description` jest w kontrakcie polem nieobowiązkowym, a brak pola znaczy
// „nie zmieniaj opisu" — tak zapowiada panel składu zespołu i tak zachowuje się
// klient, który pustego napisu nie wysyła. Opis zmienia się wtedy i tylko wtedy,
// gdy pole przyszło, także wypełnione pustką.
//
// Skład: odczyt zespołu okrawa skład z ekspertów zarchiwizowanych i usuniętych,
// więc żądanie złożone z tego, co widać w oknie, ich nie niesie. Adapter dokleja
// ich z powrotem — tak samo jak kopia zespołu (`Skopiuj`) — dzięki czemu wracają
// do składu po przywróceniu do biblioteki. Ekspert widoczny i wypisany ze składu
// wypada normalnie.
func (a *adapterZespolow) dolozStanZastany(wpisywany *dane.Zespol, zastany dane.Zespol,
	opisZadania *string) {

	if opisZadania == nil {
		wpisywany.Opis = zastany.Opis
	}
	if len(zastany.Pominieci) == 0 {
		return
	}
	wpisywany.Sklad = append(append([]string{}, wpisywany.Sklad...), zastany.Pominieci...)
	if a.dziennik != nil {
		a.dziennik.Printf("zespoły: zapis zespołu %q (%s) zachowuje %d ekspertów spoza "+
			"biblioteki: %s — żądanie ich nie niosło, bo odczyt ich nie oddaje; wrócą "+
			"do składu sami po agent.restore",
			zastany.Nazwa, zastany.Kod, len(zastany.Pominieci), strings.Join(zastany.Pominieci, ", "))
	}
}

// Wczytaj oddaje jeden zespół wraz ze składem dostępnym.
func (a *adapterZespolow) Wczytaj(ctx context.Context,
	z shared.TeamLoadRequest) (shared.TeamLoadResponse, error) {

	if a == nil || a.repozytorium == nil {
		return shared.TeamLoadResponse{}, bladBrakuZespolow(shared.CommandTeamLoad)
	}
	wiersz, err := a.repozytorium.PoKodzie(ctx, z.TeamId)
	if err != nil {
		return shared.TeamLoadResponse{}, bladWskazania(err, "zespół", z.TeamId)
	}
	return shared.TeamLoadResponse{Team: a.zespolKontraktu(wiersz)}, nil
}

// Wykaz zwraca zespoły Operatora. Brak zespołów nie jest odmową — okno pokazuje
// wtedy stan pusty i zachętę do zapisania pierwszego składu.
func (a *adapterZespolow) Wykaz(ctx context.Context,
	z shared.TeamListRequest) (shared.TeamListResponse, error) {

	if a == nil || a.repozytorium == nil {
		return shared.TeamListResponse{Teams: []shared.Team{}}, nil
	}
	filtr := dane.FiltrZespolow{Fraza: wartoscTekstu(z.Query)}
	if z.Limit != nil {
		filtr.Granica = *z.Limit
	}
	if z.Offset != nil {
		filtr.Przesuniecie = *z.Offset
	}
	wiersze, razem, err := a.repozytorium.Lista(ctx, filtr)
	if err != nil {
		return shared.TeamListResponse{}, err
	}
	zespoly := make([]shared.Team, 0, len(wiersze))
	for _, wiersz := range wiersze {
		zespoly = append(zespoly, a.zespolKontraktu(wiersz))
	}
	return shared.TeamListResponse{Teams: zespoly, Total: razem}, nil
}

// Skopiuj zakłada nowy zespół o składzie źródła. Źródło zostaje nietknięte:
// adapter je czyta i zakłada zespół obok, bez zapisu po źródle.
func (a *adapterZespolow) Skopiuj(ctx context.Context,
	z shared.TeamDuplicateRequest) (shared.TeamDuplicateResponse, error) {

	if a == nil || a.repozytorium == nil {
		return shared.TeamDuplicateResponse{}, bladBrakuZespolow(shared.CommandTeamDuplicate)
	}
	zrodlo, err := a.repozytorium.PoKodzie(ctx, z.TeamId)
	if err != nil {
		return shared.TeamDuplicateResponse{}, bladWskazania(err, "zespół", z.TeamId)
	}
	// Skład kopii bierze także ekspertów pominiętych, żeby kopia była pełna:
	// ekspert przywrócony do biblioteki wraca wtedy do obu zespołów, nie tylko
	// do źródłowego.
	sklad := append(append([]string{}, zrodlo.Sklad...), zrodlo.Pominieci...)
	kopia, err := a.repozytorium.Dodaj(ctx, dane.Zespol{
		Kod:   nowyIdentyfikator(przedrostekZespolu),
		Nazwa: nazwaKopiiZespolu(z.Name, zrodlo.Nazwa),
		Opis:  zrodlo.Opis,
		Sklad: sklad,
	})
	if err != nil {
		return shared.TeamDuplicateResponse{}, err
	}
	return shared.TeamDuplicateResponse{Team: a.zespolKontraktu(kopia)}, nil
}

// nazwaKopiiZespolu rozstrzyga nazwę kopii: podana wygrywa, pusta bierze nazwę
// źródła z przyrostkiem.
func nazwaKopiiZespolu(podana *string, zrodlowa string) string {
	if podana != nil && strings.TrimSpace(*podana) != "" {
		return strings.TrimSpace(*podana)
	}
	return zrodlowa + przyrostekKopiiZespolu
}

// zespolKontraktu przekłada wiersz zespołu na strukturę Team i — jeżeli skład
// kogoś stracił — mówi o tym wprost w dzienniku rdzenia.
func (a *adapterZespolow) zespolKontraktu(z dane.Zespol) shared.Team {
	a.powiedzOPominietych(z)
	sklad := z.Sklad
	if sklad == nil {
		sklad = []string{}
	}
	return shared.Team{
		Id:          z.Kod,
		Name:        z.Nazwa,
		Description: wskaznikTekstu(z.Opis),
		AgentIds:    sklad,
		CreatedAt:   z.Utworzono,
		UpdatedAt:   z.Zaktualizowano,
	}
}

// powiedzOPominietych zapisuje w dzienniku, których ekspertów skład pomija
// i dlaczego. Bez tej wiadomości zespół wracałby krótszy niż zapisany, bez
// wskazania, co stało się z resztą składu.
func (a *adapterZespolow) powiedzOPominietych(z dane.Zespol) {
	if a == nil || a.dziennik == nil || len(z.Pominieci) == 0 {
		return
	}
	a.dziennik.Printf("zespoły: skład zespołu %q (%s) pomija %d ekspertów — nie ma ich "+
		"w bibliotece albo są zarchiwizowani: %s; wrócą do składu sami, gdy wrócą "+
		"do biblioteki (agent.restore)",
		z.Nazwa, z.Kod, len(z.Pominieci), strings.Join(z.Pominieci, ", "))
}

// bladZadaniaZespolu odmawia wykonania komendy o niepoprawnej treści.
// Odmowa jest trójczęściowa: nazywa obszar, powód i czynność Operatora.
func bladZadaniaZespolu(powod string) error {
	return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeValidationFailed,
		"zespoły ekspertów: "+powod+"; popraw treść żądania i wyślij komendę ponownie"))
}

// bladBrakuZespolow odmawia wykonania, gdy montaż nie wpiął repozytorium:
// nazywa komendę, powód i miejsce, w którym brak się usuwa. Odpowiedź pomyślna
// udawałaby tu zapis, który nigdzie nie trafia.
func bladBrakuZespolow(komenda shared.MessageType) error {
	return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeInternalError,
		"rdzeń: komenda "+string(komenda)+" odmawia wykonania, ponieważ repozytorium "+
			"zespołów ekspertów nie jest wpięte do portu; wpięcie zakłada się w pliku "+
			"server/internal/core/montaz_porty.go wierszem "+
			"Zespoly: nowyAdapterZespolow(s.repozytoria.Zespoly).ZDziennikiem(s.montaz.Dziennik)"))
}
