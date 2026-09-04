// Plik sprawdza warunek ukończenia biegu widziany od strony warstwy rozmowy:
// pętla odróżnia pracę skończoną z wynikiem od przerwanej po tej samej
// wartości, co stan wiadomości.
package core

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"testing"
	"time"

	"danacoconsole/server/internal/models"
	"danacoconsole/server/internal/protocol"
	"danacoconsole/server/internal/session"
	"danacoconsole/shared"
)

// zamknieciePrzezKanal jest kanałem modelu, który dowozi turę sprawnie i domyka
// ją zdarzeniem wyniku o zadanym oznaczeniu błędu, jako jedyna droga tego
// pomiaru.
type zamknieciePrzezKanal struct {
	blad bool
}

func (zamknieciePrzezKanal) Kod() string                 { return "zamkniecie" }
func (zamknieciePrzezKanal) Definicja() models.Definicja { return models.Definicja{} }

func (k zamknieciePrzezKanal) Wyslij(ctx context.Context, z models.Zapytanie, u models.Ujscie) error {
	// Kształt fragmentu kończącego turę jest przepisany z kanału głównego.
	dane, err := json.Marshal(map[string]any{
		"cliSessionId": "rozm-zamkniecie",
		"subtype":      "success",
		"isError":      k.blad,
	})
	if err != nil {
		return err
	}
	pusty := ""
	return u.Fragment(ctx, models.Fragment{
		WindowId:  z.Okno(),
		MessageId: z.Wiadomosc,
		Kind:      shared.ChunkKindText,
		Text:      &pusty,
		Data:      dane,
	})
}

// zrodloJednegoKanalu podaje rejestrowi kanałów jeden czynny wiersz, potrzebny
// sprawdzianowi do wskazania kanału tury.
type zrodloJednegoKanalu struct{}

func (zrodloJednegoKanalu) Definicje(context.Context) ([]models.Definicja, error) {
	return []models.Definicja{{
		Id: 1, Kod: "zamkniecie", Nazwa: "zamkniecie", Rodzaj: "zamkniecie", Aktywny: true,
	}}, nil
}

// nadajnikCichy przyjmuje rozgłoszenia i nic z nimi nie robi — sprawdzian mierzy
// stan biegu i stan wiadomości, nie ruch na kanale.
type nadajnikCichy struct{}

func (nadajnikCichy) Rozglos(string, protocol.Koperta) {}

// stanowiskoTuryKoordynatora składa okna jednego biegu, pętlę i warstwę rozmowy
// nad kanałem, który domyka turę zadanym `is_error`.
type stanowiskoTuryKoordynatora struct {
	rozmowa     *adapterRozmowy
	petla       *session.Petla
	dziennik    *dziennikRozmowy
	koordynator session.Okno
	wykonawca   session.Okno
}

func noweStanowiskoTuryKoordynatora(t *testing.T, bladZamkniecia bool) *stanowiskoTuryKoordynatora {
	t.Helper()

	nadzorca := session.NowyNadzorca()
	sesja := nadzorca.ZalozSesje("bieg naprawczy", "", 0)
	koordynator, err := nadzorca.OtworzOkno(sesja.Id, session.Ustawienia{
		RolaOkna: shared.WindowRoleCoordinator, KanalModelu: "zamkniecie",
	})
	if err != nil {
		t.Fatalf("okno koordynatora nie powstało: %v", err)
	}
	wykonawca, err := nadzorca.OtworzOkno(sesja.Id, session.Ustawienia{
		RolaOkna: shared.WindowRoleExecutor, OknoKoordynatora: koordynator.Id,
		KanalModelu: "zamkniecie",
	})
	if err != nil {
		t.Fatalf("okno wykonawcze nie powstało: %v", err)
	}

	// Obieg rusza atrapą uruchomienia: mierzone jest zamknięcie tury koordynatora.
	petla := session.NowaPetla(nadzorca, session.UruchomienieFunkcja(
		func(session.Obieg) error { return nil }), session.UstawieniaPetli{})

	kanaly := models.NowyRejestr(zrodloJednegoKanalu{}, models.Fabryki{
		"zamkniecie": func(models.Definicja) (models.Kanal, error) {
			return zamknieciePrzezKanal{blad: bladZamkniecia}, nil
		},
	})
	if err := kanaly.Odswiez(context.Background()); err != nil {
		t.Fatalf("rejestr kanałów nie powstał: %v", err)
	}
	if _, jest := kanaly.Kanal("zamkniecie"); !jest {
		t.Fatal("kanał atrapy nie wszedł do rejestru — sprawdzian mierzyłby brak kanału")
	}

	dziennik := nowyDziennikRozmowy(context.Background(), nil, log.New(io.Discard, "", 0))
	rozmowa := nowyAdapterRozmowy(context.Background(), nadzorca, kanaly, nadajnikCichy{}, dziennik)
	rozmowa.UstawPetle(petla)

	return &stanowiskoTuryKoordynatora{
		rozmowa: rozmowa, petla: petla, dziennik: dziennik,
		koordynator: koordynator, wykonawca: wykonawca,
	}
}

// turaKoordynatora przeprowadza jedną turę okna koordynatora tą samą drogą, co
// wiadomość Operatora, i oddaje wpis odpowiedzi po jej zamknięciu.
func (s *stanowiskoTuryKoordynatora) turaKoordynatora(t *testing.T) shared.Message {
	t.Helper()

	pytanie := shared.Message{
		Id: "msg-pyt", WindowId: s.koordynator.Id, SessionId: s.koordynator.IdSesji,
		Role: shared.MessageRoleUser, Status: shared.MessageStatusComplete,
		CreatedAt: time.Now().UnixMilli(),
	}
	odpowiedz := shared.Message{
		Id: "msg-odp", WindowId: s.koordynator.Id, SessionId: s.koordynator.IdSesji,
		Role: shared.MessageRoleAssistant, Status: shared.MessageStatusStreaming,
		CreatedAt: time.Now().UnixMilli(),
	}
	s.dziennik.Dopisz(pytanie)
	s.dziennik.Dopisz(odpowiedz)
	s.rozmowa.prowadzTure(context.Background(), s.koordynator, pytanie, odpowiedz, "zad-zamkniecie", nil)

	wiadomosci, _ := s.dziennik.Wykaz(s.koordynator.Id, nil, nil)
	for _, w := range wiadomosci {
		if w.Id == odpowiedz.Id {
			return w
		}
	}
	t.Fatal("odpowiedź tury nie wróciła z dziennika rozmowy — pomiar nie doszedł do skutku")
	return shared.Message{}
}

// TestTuraZamknietaBledemNieOglaszaUkonczenia mierzy zgodność dwóch
// rozstrzygnięć wyprowadzonych z tego samego zamknięcia tury: stanu wiadomości
// i powodu zatrzymania biegu.
func TestTuraZamknietaBledemNieOglaszaUkonczenia(t *testing.T) {
	s := noweStanowiskoTuryKoordynatora(t, true)

	// Bieg musi mieć za sobą obieg, inaczej sprawdzian przejdzie niezależnie od naprawy.
	s.petla.ZakonczTure(s.wykonawca.Id, session.PowodWynik)
	if s.petla.Stan(s.koordynator.Id).Obiegow == 0 {
		t.Fatal("licznik obiegów stoi na zerze — bieg nie ruszył, nie ma czego mierzyć")
	}

	odpowiedz := s.turaKoordynatora(t)
	if odpowiedz.Status != shared.MessageStatusError {
		t.Fatalf("stan wiadomości: chciano %q, jest %q", shared.MessageStatusError, odpowiedz.Status)
	}

	stan := s.petla.Stan(s.koordynator.Id)
	if stan.Zatrzymany && stan.Powod == session.ZatrzymanieUkonczenie {
		t.Fatalf("bieg ogłoszony jako ukończony z wynikiem, choć wiadomość ma stan %q",
			odpowiedz.Status)
	}
}

// TestTuraZamknietaWynikiemKonczyBieg pilnuje drugiej strony rozróżnienia:
// zamknięcie bez błędu nadal kończy bieg ukończeniem z wynikiem. Bez tego
// sprawdzianu naprawa mogłaby odebrać pętli ukończenie w ogóle i nikt by tego
// nie zauważył.
func TestTuraZamknietaWynikiemKonczyBieg(t *testing.T) {
	s := noweStanowiskoTuryKoordynatora(t, false)

	s.petla.ZakonczTure(s.wykonawca.Id, session.PowodWynik)
	if s.petla.Stan(s.koordynator.Id).Obiegow == 0 {
		t.Fatal("licznik obiegów stoi na zerze — bieg nie ruszył, nie ma czego mierzyć")
	}

	odpowiedz := s.turaKoordynatora(t)
	if odpowiedz.Status != shared.MessageStatusComplete {
		t.Fatalf("stan wiadomości: chciano %q, jest %q", shared.MessageStatusComplete, odpowiedz.Status)
	}

	stan := s.petla.Stan(s.koordynator.Id)
	if !stan.Zatrzymany || stan.Powod != session.ZatrzymanieUkonczenie {
		t.Fatalf("bieg nie stanął ukończeniem z wynikiem: zatrzymany=%v powód=%q",
			stan.Zatrzymany, stan.Powod)
	}
}
