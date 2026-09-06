package core

import (
	"context"

	"danacoconsole/server/internal/dane"
	"danacoconsole/server/internal/protocol"
	"danacoconsole/shared"
)

// Jeden silnik na pętlę sesyjną i MultitaskingAI; drugiego przebiegu stanów nie ma.
type silnikKolejki struct {
	repozytorium dane.RepozytoriumKolejek
	wykonawca    wykonawcaKroku
	ujscieWyniku func(ctx context.Context, pozycjaID int64, tresc string)
}

func (s silnikKolejki) ZWykonawca(w wykonawcaKroku) silnikKolejki {
	s.wykonawca = w
	return s
}

func (s silnikKolejki) ZUjsciemWyniku(ujscie func(ctx context.Context, pozycjaID int64, tresc string)) silnikKolejki {
	s.ujscieWyniku = ujscie
	return s
}

// Wykaz pusty zostawia kolejkę bez pozycji — to poprawny stan, nie awaria.
func (s silnikKolejki) Zasil(ctx context.Context, kolejkaID int64, pozycje []dane.Pozycja) error {
	for _, pozycja := range pozycje {
		pozycja.KolejkaID = kolejkaID
		if _, err := s.repozytorium.DodajPozycje(ctx, pozycja); err != nil {
			return err
		}
	}
	return nil
}

// Stan kolejki jest wyprowadzony z pozycji, nie zadeklarowany.
func (s silnikKolejki) Wykonaj(ctx context.Context, kolejkaID int64,
	dzialanie shared.QueueAction, idPozycji *string) (shared.QueueStatus, error) {

	pozycje, err := s.repozytorium.ListaPozycji(ctx, kolejkaID)
	if err != nil {
		return "", err
	}
	if err := s.zastosuj(ctx, kolejkaID, dzialanie, pozycje, idPozycji); err != nil {
		return "", err
	}
	poZmianie, err := s.repozytorium.ListaPozycji(ctx, kolejkaID)
	if err != nil {
		return "", err
	}
	return stanKolejkiPoDzialaniu(dzialanie, poZmianie), nil
}

// Błąd odczytu daje wykaz pusty, żeby kolejka nie znikała z odpowiedzi.
func (s silnikKolejki) Pozycje(ctx context.Context, kolejkaID int64) []dane.Pozycja {
	pozycje, err := s.repozytorium.ListaPozycji(ctx, kolejkaID)
	if err != nil {
		return nil
	}
	return pozycje
}

// Kolejka wyczerpana pokazuje licznik pozycji ostatniej, pusta — żadnego.
func (s silnikKolejki) Cykl(pozycje []dane.Pozycja) *int {
	if len(pozycje) == 0 {
		return nil
	}
	biezaca := pierwszaCzynna(pozycje)
	if biezaca == nil {
		biezaca = &pozycje[len(pozycje)-1]
	}
	obieg := biezaca.LicznikObiegow
	return &obieg
}

func (s silnikKolejki) zastosuj(ctx context.Context, kolejkaID int64,
	dzialanie shared.QueueAction, pozycje []dane.Pozycja, idPozycji *string) error {

	cel, err := celDzialania(pozycje, idPozycji)
	if err != nil {
		return err
	}
	switch dzialanie {
	case shared.QueueActionStart, shared.QueueActionResume:
		return s.krok(ctx, cel)
	case shared.QueueActionRetry:
		return s.biegNaprawczy(ctx, celNaprawy(pozycje, cel))
	case shared.QueueActionStop, shared.QueueActionClear:
		return s.anuluj(ctx, przerywane(pozycje, cel, czyWskazano(idPozycji)))
	case shared.QueueActionEnqueue:
		return wstawZlecenie(ctx, s.repozytorium, kolejkaID, pozycje, cel, czyWskazano(idPozycji))
	case shared.QueueActionDequeue:
		return zdejmijZlecenie(ctx, s.repozytorium, cel, czyWskazano(idPozycji))
	case shared.QueueActionPause:
		// Bieg podjęty dobiega do końca; zmienia się stan kolejki, nie pozycji.
		return nil
	}
	// Bez tej gałęzi `switch` oddawał sukces, choć nic się nie ruszyło.
	return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeValidationFailed,
		"Działanie kolejki „"+string(dzialanie)+"” nie jest prowadzone przez silnik kolejek."))
}

// Brak pozycji do podjęcia nie jest błędem — kolejka nie ma czego posunąć.
func (s silnikKolejki) krok(ctx context.Context, pozycja *dane.Pozycja) error {
	if pozycja == nil {
		return nil
	}
	nastepny, jest := krokNaprzod[pozycja.Stan]
	if !jest {
		return nil
	}
	if err := s.repozytorium.ZmienStanPozycji(ctx, pozycja.ID, nastepny, werdyktKroku(nastepny)); err != nil {
		return err
	}
	if nastepny == stanPozycjiWykonywana {
		pozycja.Stan = stanPozycjiWykonywana
		return s.wykonaj(ctx, *pozycja)
	}
	return nil
}

// Licznik obiegów rośnie bez progu, a odmowy nie ma na żadnym obiegu.
func (s silnikKolejki) biegNaprawczy(ctx context.Context, pozycja *dane.Pozycja) error {
	if pozycja == nil {
		return nil
	}
	if _, err := s.repozytorium.ZwiekszObieg(ctx, pozycja.ID); err != nil {
		return err
	}
	werdykt := werdyktDoPoprawy
	if err := s.repozytorium.ZmienStanPozycji(ctx, pozycja.ID, stanPozycjiWykonywana, &werdykt); err != nil {
		return err
	}
	pozycja.Stan = stanPozycjiWykonywana
	return s.wykonaj(ctx, *pozycja)
}

// Wskazane zlecenie wraca do kolejki nowym wierszem: zamknięte zostaje z wynikiem.
func wstawZlecenie(ctx context.Context, repozytorium dane.RepozytoriumKolejek, kolejkaID int64,
	pozycje []dane.Pozycja, cel *dane.Pozycja, wskazano bool) error {

	if !wskazano {
		// Zlecenia automatyki dokłada Queue Manager przed wejściem w silnik.
		if pierwszaCzynna(pozycje) == nil {
			return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeValidationFailed,
				"Wstawienie zlecenia wymaga wskazania zlecenia (itemId) albo automatyki (workflowId) wnoszącej kroki."))
		}
		return nil
	}
	if cel == nil {
		return nil
	}
	wstawiane := dane.Pozycja{
		KolejkaID: kolejkaID, OknoWykonawcyID: cel.OknoWykonawcyID, Tytul: cel.Tytul,
		TrescZlecenia: cel.TrescZlecenia, TrescOdwolanie: cel.TrescOdwolanie,
		Stan: stanPozycjiOczekuje,
	}
	_, err := repozytorium.DodajPozycje(ctx, wstawiane)
	return err
}

// Wiersz zostaje anulowany, nie skasowany: kolejka jest zapisem przebiegu.
func zdejmijZlecenie(ctx context.Context, repozytorium dane.RepozytoriumKolejek, cel *dane.Pozycja, wskazano bool) error {
	if !wskazano {
		return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeValidationFailed,
			"Zdjęcie zlecenia z kolejki wymaga wskazania zlecenia (itemId)."))
	}
	if cel == nil {
		return nil
	}
	return repozytorium.ZmienStanPozycji(ctx, cel.ID, stanPozycjiAnulowana, nil)
}

// Silnik bez wpiętego wykonawcy zostawia pozycję w stanie wykonywana.
func (s silnikKolejki) wykonaj(ctx context.Context, pozycja dane.Pozycja) error {
	if s.wykonawca == nil {
		return nil
	}
	tresc, err := s.wykonawca.Wykonaj(ctx, pozycja)
	if s.ujscieWyniku != nil && tresc != "" {
		s.ujscieWyniku(ctx, pozycja.ID, tresc)
	}
	if err != nil {
		// Niepowodzenie to stan błędu, nie ukończenie; ponawia je bieg naprawczy.
		return s.domknijPoTurze(ctx, pozycja, stanPozycjiBledna)
	}
	return s.domknijPoTurze(ctx, pozycja, stanPozycjiDoWeryfikacji)
}

// Zatrzymanie wydane w trakcie tury ma pierwszeństwo przed spóźnionym werdyktem.
func (s silnikKolejki) domknijPoTurze(ctx context.Context, pozycja dane.Pozycja, stan string) error {
	pozycje, err := s.repozytorium.ListaPozycji(ctx, pozycja.KolejkaID)
	if err == nil {
		for _, biezaca := range pozycje {
			if biezaca.ID != pozycja.ID {
				continue
			}
			if _, koncowy := stanyKoncowePozycji[biezaca.Stan]; koncowy {
				return nil
			}
			break
		}
	}
	return s.repozytorium.ZmienStanPozycji(ctx, pozycja.ID, stan, nil)
}

// Wiersz zostaje — ślad przerwanego zlecenia należy do przejrzystości pętli.
func (s silnikKolejki) anuluj(ctx context.Context, pozycje []dane.Pozycja) error {
	for _, pozycja := range pozycje {
		if err := s.repozytorium.ZmienStanPozycji(ctx, pozycja.ID, stanPozycjiAnulowana, nil); err != nil {
			return err
		}
	}
	return nil
}
