package core

import (
	"context"

	"danacoconsole/server/internal/dane"
	"danacoconsole/shared"
)

// Silnik wykonania kolejki: zlecenie powstaje, przechodzi stany i kończy się wynikiem. Jeden silnik obsługuje pętlę sesyjną i MultitaskingAI; drugiego przebiegu stanów nie ma nigdzie indziej.
type silnikKolejki struct {
	repozytorium dane.RepozytoriumKolejek
	// wykonawca uruchamia pracę pozycji wchodzącej w stan wykonywana; pusty zostawia silnik przy stanach.
	wykonawca wykonawcaKroku
	// ujscieWyniku odbiera zebraną treść odpowiedzi tury; puste znaczy, że płynie wyłącznie strumieniem.
	ujscieWyniku func(ctx context.Context, pozycjaID int64, tresc string)
}

// ZWykonawca wpina most do realnego wykonania pozycji. Zwraca silnik przez
// wartość, bo `silnikKolejki` trzymany jest w adapterze kolejek jako pole
// wartościowe, a nie wskaźnik — montaż podmienia je w miejscu.
func (s silnikKolejki) ZWykonawca(w wykonawcaKroku) silnikKolejki {
	s.wykonawca = w
	return s
}

// ZUjsciemWyniku wpina odbiorcę zebranej treści tury. Silnik nie wie, kto
// odbiera — dziś jest to wiersz podagenta (`adapter_modul_orkiestracja.go`),
// ale silnik zna wyłącznie pozycję; pozycja bez odbiorcy przechodzi bez śladu
// treści.
func (s silnikKolejki) ZUjsciemWyniku(ujscie func(ctx context.Context, pozycjaID int64, tresc string)) silnikKolejki {
	s.ujscieWyniku = ujscie
	return s
}

// Zasil zakłada zlecenia początkowe kolejki. Wykaz pusty zostawia kolejkę bez
// pozycji — to poprawny stan, nie awaria.
func (s silnikKolejki) Zasil(ctx context.Context, kolejkaID int64, pozycje []dane.Pozycja) error {
	for _, pozycja := range pozycje {
		pozycja.KolejkaID = kolejkaID
		if _, err := s.repozytorium.DodajPozycje(ctx, pozycja); err != nil {
			return err
		}
	}
	return nil
}

// Wykonaj stosuje działanie do pozycji kolejki i zwraca stan, w jakim kolejka
// znajduje się po nim. Stan kolejki jest wyprowadzony z pozycji, nie zadeklarowany:
// dopóki jest co robić, kolejka pracuje; gdy nie ma — jest wyczerpana.
func (s silnikKolejki) Wykonaj(ctx context.Context, kolejkaID int64,
	dzialanie shared.QueueAction, idPozycji *string) (shared.QueueStatus, error) {

	pozycje, err := s.repozytorium.ListaPozycji(ctx, kolejkaID)
	if err != nil {
		return "", err
	}
	if err := s.zastosuj(ctx, dzialanie, pozycje, idPozycji); err != nil {
		return "", err
	}
	poZmianie, err := s.repozytorium.ListaPozycji(ctx, kolejkaID)
	if err != nil {
		return "", err
	}
	return stanKolejkiPoDzialaniu(dzialanie, poZmianie), nil
}

// Pozycje zwraca zlecenia kolejki. Błąd odczytu daje wykaz pusty — odpowiedź
// o kolejce nie ma znikać z powodu jednego zapytania pobocznego.
func (s silnikKolejki) Pozycje(ctx context.Context, kolejkaID int64) []dane.Pozycja {
	pozycje, err := s.repozytorium.ListaPozycji(ctx, kolejkaID)
	if err != nil {
		return nil
	}
	return pozycje
}

// Cykl podaje licznik obiegów pozycji, na której stoi kolejka — pole
// Queue.Cycle kontraktu. Kolejka wyczerpana pokazuje licznik pozycji ostatniej,
// kolejka pusta nie pokazuje żadnego.
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

// zastosuj przeprowadza pozycje przez działanie, jedną po drugiej, aż do wyczerpania całego wykazu kolejki.
func (s silnikKolejki) zastosuj(ctx context.Context, dzialanie shared.QueueAction,
	pozycje []dane.Pozycja, idPozycji *string) error {

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
	}
	return nil
}

// krok posuwa pozycję o jeden stan naprzód; brak pozycji do podjęcia nie jest błędem, bo kolejka pusta albo wyczerpana po prostu nie ma czego posunąć.
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

// biegNaprawczy podnosi licznik obiegów i zawraca pozycję do wykonania; licznik rośnie bez progu, a odmowy nie ma na żadnym obiegu.
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

// wykonaj uruchamia realną pracę pozycji, która właśnie weszła w stan wykonywana, i przesuwa jej stan według wyniku; silnik bez wpiętego wykonawcy zostawia pozycję w wykonywana.
func (s silnikKolejki) wykonaj(ctx context.Context, pozycja dane.Pozycja) error {
	if s.wykonawca == nil {
		return nil
	}
	tresc, err := s.wykonawca.Wykonaj(ctx, pozycja)
	if s.ujscieWyniku != nil && tresc != "" {
		s.ujscieWyniku(ctx, pozycja.ID, tresc)
	}
	if err != nil {
		// Niepowodzenie wykonania jest stanem błędu, nie ukończeniem; da się ponowić biegiem naprawczym.
		return s.domknijPoTurze(ctx, pozycja, stanPozycjiBledna)
	}
	// Praca się zakończyła; wynik czeka na weryfikację, do ukonczona przesuwa go dopiero przyjęcie wyniku.
	return s.domknijPoTurze(ctx, pozycja, stanPozycjiDoWeryfikacji)
}

// domknijPoTurze zapisuje stan pozycji wynikający z tury, chyba że pozycja w międzyczasie weszła w stan końcowy; queue.action stop wydane w trakcie tury ma pierwszeństwo przed spóźnionym werdyktem.
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

// anuluj zamyka pozycje przerwaniem. Wiersz zostaje razem z dziennikiem —
// ślad przerwanego zlecenia jest częścią przejrzystości pętli.
func (s silnikKolejki) anuluj(ctx context.Context, pozycje []dane.Pozycja) error {
	for _, pozycja := range pozycje {
		if err := s.repozytorium.ZmienStanPozycji(ctx, pozycja.ID, stanPozycjiAnulowana, nil); err != nil {
			return err
		}
	}
	return nil
}
