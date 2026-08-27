// Plik niesie silnik wybudzeń automatyki: zawieszanie biegu na sygnał ze świata, doręczanie sygnału i rozstrzyganie, co się dzieje, gdy reakcja nie przyjdzie.
package core

import (
	"context"
	"errors"
	"log"
	"time"

	"danacoconsole/server/internal/dane"
)

// interwalWybudzen wyznacza takt sprawdzania terminów oczekiwań. Minuta zgadza
// się z taktem budzika harmonogramu — oba zegary pytają bazę o to samo: co się
// należy teraz. Częstsze pytanie nie trafiłoby w żaden nowy termin.
const interwalWybudzen = time.Minute

// Powody wybudzenia — słownik zamknięty, zgodny z więzem CHECK kolumny powod w tabeli oczekiwań pętli.
const (
	powodSygnal   = "sygnal"
	powodTermin   = "termin"
	powodOperator = "operator"
)

// Rozstrzygnięcia po upływie terminu — słownik zamknięty, zgodny z więzem CHECK
// kolumny `po_terminie`.
const (
	poTerminieWznow    = "wznow"
	poTerminiePonow    = "ponow"
	poTerminiePrzerwij = "przerwij"
)

// repozytoriumWybudzen to wszystko, czego silnik potrzebuje od warstwy danych: byty pętli plus odczyt i zapis przebiegu; interfejs łączy dwa zakresy w jeden widok jednego odbiorcy.
type repozytoriumWybudzen interface {
	dane.RepozytoriumPetli
	Przebieg(ctx context.Context, kod string) (dane.Przebieg, error)
	ZapiszPrzebieg(ctx context.Context, przebieg dane.Przebieg) (dane.Przebieg, error)
}

// silnikWybudzen prowadzi oczekiwania biegów: zakłada je, doręcza sygnały i pilnuje terminów zegara pętli.
type silnikWybudzen struct {
	repozytorium repozytoriumWybudzen
	dziennik     *log.Logger
}

// nowySilnikWybudzen wiąże silnik z repozytorium adaptera Automations.
// Repozytorium, które nie niesie bytów pętli, daje silnik pusty zamiast awarii:
// rdzeń wstaje, automatyka pracuje, a wybudzeń po prostu nie ma.
func nowySilnikWybudzen(automatyki *adapterAutomatyk, dziennik *log.Logger) *silnikWybudzen {
	if automatyki == nil || automatyki.repozytorium == nil {
		return nil
	}
	petla, spelnia := automatyki.repozytorium.(repozytoriumWybudzen)
	if !spelnia {
		return nil
	}
	return &silnikWybudzen{repozytorium: petla, dziennik: dziennik}
}

// Uruchom wpina silnik w cykl życia rdzenia; pierwsze sprawdzenie idzie od razu po starcie, żeby terminy minione w czasie postoju rdzenia zostały rozstrzygnięte bez czekania na pełny takt.
func (s *silnikWybudzen) Uruchom(zycie context.Context) {
	if s == nil {
		return
	}
	go s.petla(zycie)
}

func (s *silnikWybudzen) petla(zycie context.Context) {
	zegar := time.NewTicker(interwalWybudzen)
	defer zegar.Stop()
	s.sprawdzTerminy(zycie, time.Now().UTC())
	for {
		select {
		case <-zycie.Done():
			return
		case teraz := <-zegar.C:
			s.sprawdzTerminy(zycie, teraz.UTC())
		}
	}
}

// Zawies zatrzymuje bieg na sygnał ze świata i przestawia go w stan oczekuje; kolejność ma znaczenie, najpierw wiersz oczekiwania, potem stan biegu, żeby awaria między nimi nie zostawiła biegu bez oczekiwania.
func (s *silnikWybudzen) Zawies(ctx context.Context, przebiegKod, krok, sygnal string,
	termin *string, poTerminie string) (dane.OczekiwanieBiegu, error) {

	if s == nil {
		return dane.OczekiwanieBiegu{}, errors.New("silnik wybudzeń niewpięty")
	}
	przebieg, err := s.repozytorium.Przebieg(ctx, przebiegKod)
	if err != nil {
		return dane.OczekiwanieBiegu{}, err
	}
	oczekiwanie, err := s.repozytorium.ZalozOczekiwanie(ctx, dane.OczekiwanieBiegu{
		PrzebiegID:     przebieg.ID,
		KrokZewnetrzny: krok,
		Sygnal:         sygnal,
		Termin:         termin,
		PoTerminie:     poTerminie,
	})
	if err != nil {
		return dane.OczekiwanieBiegu{}, err
	}
	przebieg.Stan = "oczekuje"
	if _, err := s.repozytorium.ZapiszPrzebieg(ctx, przebieg); err != nil {
		return dane.OczekiwanieBiegu{}, err
	}
	return oczekiwanie, nil
}

// Wybudz doręcza sygnał wszystkim biegom, które na niego czekają, i zwraca
// liczbę wybudzonych. Sygnał bez adresata nie jest błędem: świat zewnętrzny nie
// wie, które biegi czekają, więc zero wybudzonych to poprawna odpowiedź.
func (s *silnikWybudzen) Wybudz(ctx context.Context, sygnal string, tresc *string) (int, error) {
	if s == nil {
		return 0, errors.New("silnik wybudzeń niewpięty")
	}
	czekajace, err := s.repozytorium.OczekiwaniaNaSygnal(ctx, sygnal)
	if err != nil {
		return 0, err
	}
	wybudzonych := 0
	for _, oczekiwanie := range czekajace {
		if err := s.wznow(ctx, oczekiwanie, powodSygnal, tresc); err != nil {
			s.zapisz("silnik wybudzeń: nie można wybudzić biegu %s: %v", oczekiwanie.PrzebiegKod, err)
			continue
		}
		wybudzonych++
	}
	return wybudzonych, nil
}

// ZamknijRecznie wybudza bieg na żądanie operatora, bez czekania na sygnał ani
// na termin. Jest to trzecia — obok sygnału i terminu — droga wyjścia biegu ze
// stanu oczekiwania.
func (s *silnikWybudzen) ZamknijRecznie(ctx context.Context, przebiegID int64) error {
	if s == nil {
		return errors.New("silnik wybudzeń niewpięty")
	}
	oczekiwanie, err := s.repozytorium.OczekiwanieCzynne(ctx, przebiegID)
	if errors.Is(err, dane.ErrBrakWiersza) {
		return nil // bieg na nic nie czeka — stan, o który chodziło
	}
	if err != nil {
		return err
	}
	return s.wznow(ctx, oczekiwanie, powodOperator, nil)
}

// sprawdzTerminy rozstrzyga oczekiwania, którym minął już termin zegara pętli automatyki tego rdzenia.
func (s *silnikWybudzen) sprawdzTerminy(ctx context.Context, teraz time.Time) {
	spoznione, err := s.repozytorium.OczekiwaniaPoTerminie(ctx, teraz.Format(formatZnacznikaBazy))
	if err != nil {
		s.zapisz("silnik wybudzeń: nie można odczytać oczekiwań po terminie: %v", err)
		return
	}
	for _, oczekiwanie := range spoznione {
		if err := s.rozstrzygnijTermin(ctx, oczekiwanie); err != nil {
			s.zapisz("silnik wybudzeń: nie można rozstrzygnąć terminu biegu %s: %v",
				oczekiwanie.PrzebiegKod, err)
		}
	}
}

// rozstrzygnijTermin stosuje politykę zapisaną przy zakładaniu oczekiwania.
// Awaria jednego biegu nie zatrzymuje zegara — pętla idzie dalej.
func (s *silnikWybudzen) rozstrzygnijTermin(ctx context.Context, o dane.OczekiwanieBiegu) error {
	switch o.PoTerminie {
	case poTerminiePrzerwij:
		return s.przerwij(ctx, o)
	case poTerminiePonow:
		return s.wznowZProba(ctx, o)
	case poTerminieWznow:
		return s.wznow(ctx, o, powodTermin, nil)
	default:
		// Wartość spoza słownika nie zawiesza biegu na zawsze: silnik wraca do zachowania domyślnego i śladu.
		s.zapisz("silnik wybudzeń: nieznane rozstrzygnięcie %q biegu %s — wznawiam",
			o.PoTerminie, o.PrzebiegKod)
		return s.wznow(ctx, o, powodTermin, nil)
	}
}

// wznow zamyka oczekiwanie i przywraca bieg do pracy — tak, jak gdyby oczekiwany sygnał właśnie przyszedł.
func (s *silnikWybudzen) wznow(ctx context.Context, o dane.OczekiwanieBiegu,
	powod string, tresc *string) error {

	return s.zamknijIZapisz(ctx, o, powod, tresc, func(przebieg *dane.Przebieg) {
		przebieg.Stan = "running"
		przebieg.KomunikatBledu = nil
	})
}

// wznowZProba wznawia bieg, licząc próbę. `proba` nie ma górnej granicy:
// licznik rośnie i jest widoczny, ale sam nie zatrzymuje pracy.
func (s *silnikWybudzen) wznowZProba(ctx context.Context, o dane.OczekiwanieBiegu) error {
	return s.zamknijIZapisz(ctx, o, powodTermin, nil, func(przebieg *dane.Przebieg) {
		przebieg.Stan = "running"
		przebieg.Proba++
	})
}

// przerwij kończy bieg, zostawiając powód. Komunikat nazywa sygnał i krok —
// bez tego Operator zobaczyłby bieg przerwany bez przyczyny.
func (s *silnikWybudzen) przerwij(ctx context.Context, o dane.OczekiwanieBiegu) error {
	powod := "bieg czekał na sygnał " + o.Sygnal + " w kroku " + o.KrokZewnetrzny +
		"; termin minął, a rozstrzygnięcie oczekiwania brzmi „przerwij”"
	return s.zamknijIZapisz(ctx, o, powodTermin, nil, func(przebieg *dane.Przebieg) {
		przebieg.Stan = "stopped"
		przebieg.KomunikatBledu = &powod
	})
}

// zamknijIZapisz domyka oczekiwanie i nanosi zmianę na przebieg; zamknięcie idzie pierwsze i jest idempotentne, dzięki czemu sygnał doręczony dwa razy wybudza bieg raz.
func (s *silnikWybudzen) zamknijIZapisz(ctx context.Context, o dane.OczekiwanieBiegu,
	powod string, tresc *string, zmiana func(*dane.Przebieg)) error {

	if err := s.repozytorium.ZamknijOczekiwanie(ctx, o.ID, powod, tresc); err != nil {
		return err
	}
	przebieg, err := s.repozytorium.Przebieg(ctx, o.PrzebiegKod)
	if err != nil {
		return err
	}
	zmiana(&przebieg)
	_, err = s.repozytorium.ZapiszPrzebieg(ctx, przebieg)
	return err
}

// zapisz nanosi wiersz do dziennika rdzenia, znosząc dziennik pusty milczącym pominięciem samego zapisu.
func (s *silnikWybudzen) zapisz(wzorzec string, argumenty ...any) {
	if s == nil || s.dziennik == nil {
		return
	}
	s.dziennik.Printf(wzorzec, argumenty...)
}
