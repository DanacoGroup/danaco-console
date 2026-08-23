package core

import (
	"context"

	"danacoconsole/server/internal/dane"
	"danacoconsole/shared"
)

// Silnik wykonania kolejki: zlecenie powstaje, przechodzi stany i kończy się
// wynikiem. Jeden silnik obsługuje pętlę sesyjną i MultitaskingAI;
// drugiego przebiegu stanów nie ma nigdzie indziej.
//
// Protokół działań. Kontrakt daje sześć działań (QueueAction) i nie ma osobnego
// działania „zamknij pozycję werdyktem". Silnik czyta więc działania dosłownie
// tak, jak są nazwane, i posuwa pozycje po tabeli przejść `krokNaprzod`:
//
//	start · resume  krok naprzód: oczekuje → wykonywana → do_weryfikacji → ukonczona
//	retry           bieg naprawczy: licznik obiegów +1, werdykt do_poprawy,
//	                powrót do wykonywana; bez limitu i bez warunku
//	stop            przerwanie: pozycja wskazana albo wszystkie czynne → anulowana
//	pause           wstrzymanie kolejki; pozycje zostają, gdzie były
//	clear           opróżnienie: pozycje czynne → anulowana, dziennik zostaje
//	                (przejrzystość zamiast kasowania śladu)
//
// Przyjęcie wyniku kroku wyraża się wyborem działania: krok naprzód znaczy
// przyjęcie, retry znaczy odesłanie do poprawy. Silnik nie wystawia werdyktu,
// którego nie wywołało działanie Operatora, i nie zmyśla postępu.
//
// Działanie wskazujące pozycję (`itemId`) dotyczy tej pozycji; bez wskazania
// dotyczy pozycji, na której kolejka stoi.
//
// Most do realnego wykonania. Sam przebieg stanów nie wykonuje pracy — pozycja
// wchodząca w stan `wykonywana` musi zostać naprawdę wykonana, a jej wynik, nie
// klik Operatora, przesuwa ją dalej: powodzenie do `do_weryfikacji`,
// niepowodzenie do `bledna`, czyli do realnego stanu błędu zamiast cichego
// ukończenia. Robi to wpięty `wykonawca`. Silnik bez wykonawcy zostaje czystą
// maszyną stanów: pozycję posuwa wtedy działanie Operatora, pętli sesyjnej albo
// MultitaskingAI, tak jak przed wpięciem mostu.
type silnikKolejki struct {
	repozytorium dane.RepozytoriumKolejek
	// wykonawca uruchamia realną pracę pozycji wchodzącej w stan wykonywana.
	// Pusty zostawia silnik przy samym przebiegu stanów.
	wykonawca wykonawcaKroku
	// ujscieWyniku odbiera zebraną treść odpowiedzi tury pozycji.
	// Puste znaczy: treść płynie wyłącznie strumieniem. Ujście dostaje
	// treść także przy pozycji przerwanej w locie — praca, która się odbyła,
	// zostaje widoczna, a nie wymazana.
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

// zastosuj przeprowadza pozycje przez działanie.
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

// krok posuwa pozycję o jeden stan naprzód. Brak pozycji do podjęcia nie jest
// błędem — kolejka pusta albo wyczerpana po prostu nie ma czego posunąć.
//
// Wejście w stan `wykonywana` nie kończy kroku: pozycja jest wtedy naprawdę
// wykonywana, a jej wynik przesuwa ją dalej.
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

// biegNaprawczy podnosi licznik obiegów i zawraca pozycję do wykonania.
// Licznik rośnie bez progu, a odmowy nie ma na żadnym obiegu.
// Zawrócona pozycja wchodzi w `wykonywana`, więc jest wykonywana od nowa —
// powtórzenie kroku ma powtórzyć pracę, nie samo przełożyć etykietę stanu.
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

// wykonaj uruchamia realną pracę pozycji, która właśnie weszła w stan
// wykonywana, i przesuwa jej stan według wyniku. Silnik bez wpiętego wykonawcy
// nie ma czym wykonać kroku — zostawia pozycję w `wykonywana`, a dalej posuwa ją
// działanie Operatora, pętli albo MultitaskingAI.
func (s silnikKolejki) wykonaj(ctx context.Context, pozycja dane.Pozycja) error {
	if s.wykonawca == nil {
		return nil
	}
	tresc, err := s.wykonawca.Wykonaj(ctx, pozycja)
	if s.ujscieWyniku != nil && tresc != "" {
		s.ujscieWyniku(ctx, pozycja.ID, tresc)
	}
	if err != nil {
		// Niepowodzenie wykonania jest realnym stanem błędu, nie cichym
		// ukończeniem. Pozycję da się ponowić biegiem naprawczym.
		return s.domknijPoTurze(ctx, pozycja, stanPozycjiBledna)
	}
	// Praca się zakończyła; wynik czeka na weryfikację. Do `ukonczona` pozycję
	// przesuwa dopiero przyjęcie wyniku — krok naprzód po weryfikacji.
	return s.domknijPoTurze(ctx, pozycja, stanPozycjiDoWeryfikacji)
}

// domknijPoTurze zapisuje stan pozycji wynikający z tury — chyba że pozycja
// w międzyczasie weszła w stan końcowy.
//
// Tura pozycji biegnie synchronicznie wewnątrz `Wykonaj`, a `queue.action stop`
// wydane w trakcie tury zapisuje pozycji stan `anulowana`. Bez tego sprawdzenia
// zapis poniżej nadpisałby go chwilę później stanem `do_weryfikacji`, jak gdyby
// przerwania nie było — łamiąc własny protokół silnika („stop · przerwanie:
// pozycja … → anulowana") i zamykając jedyną kontraktową drogę zatrzymania
// pracy podagenta, bo komendy `subagent.stop` kontrakt nie ma. Przerwanie ma
// pierwszeństwo przed spóźnionym werdyktem tury; sama treść odpowiedzi trafiła
// już do ujścia wyniku, więc nic z pracy nie ginie po cichu.
//
// Odczyt idzie listą pozycji kolejki, bo repozytorium nie ma odczytu jednej
// pozycji — a dorabianie go dla tego jednego miejsca byłoby drugim zapytaniem
// o to samo. Pozycja nieodnaleziona przechodzi na zapis wprost:
// lepiej zapisać stan wynikający z tury, niż zgubić go z powodu błędu odczytu.
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
