package session

import (
	"errors"
	"testing"

	"danacoconsole/server/internal/protocol"
	"danacoconsole/shared"
)

// Warunek ukończenia biegu koordynator-wykonawca mierzy to, co pętla naprawdę dostaje.

// stanowiskoPetli jest pętlą wraz z oknami jednego biegu i śladem tego, co
// pętla zrobiła: obiegami oddanymi portowi rozpoczynania i zdarzeniami
// rozesłanymi obserwatorom.
type stanowiskoPetli struct {
	petla       *Petla
	koordynator string
	wykonawca   string
	obiegi      []Obieg
	slady       []ZdarzenieObiegu
	// usterka rozpoczęcia obiegu; pusta znaczy obieg rozpoczęty.
	usterka error
}

// Funkcja noweStanowisko składa nadzorcę, sesję, okno koordynatora, okno wykonawcze i pętlę nad nimi na potrzeby sprawdzianu.
func noweStanowisko(t *testing.T) *stanowiskoPetli {
	t.Helper()

	s := &stanowiskoPetli{}
	nadzorca := NowyNadzorca()
	sesja := nadzorca.ZalozSesje("bieg naprawczy", "", 0)

	koordynator, err := nadzorca.OtworzOkno(sesja.Id, Ustawienia{
		RolaOkna: shared.WindowRoleCoordinator,
	})
	if err != nil {
		t.Fatalf("okno koordynatora nie powstało: %v", err)
	}
	wykonawca, err := nadzorca.OtworzOkno(sesja.Id, Ustawienia{
		RolaOkna:         shared.WindowRoleExecutor,
		OknoKoordynatora: koordynator.Id,
	})
	if err != nil {
		t.Fatalf("okno wykonawcze nie powstało: %v", err)
	}
	s.koordynator, s.wykonawca = koordynator.Id, wykonawca.Id

	s.petla = NowaPetla(nadzorca, UruchomienieFunkcja(func(o Obieg) error {
		s.obiegi = append(s.obiegi, o)
		return s.usterka
	}), UstawieniaPetli{})
	s.petla.Obserwuj(ObserwatorFunkcja(func(z ZdarzenieObiegu) {
		s.slady = append(s.slady, z)
	}))
	return s
}

// Metoda koniecTuryWykonawcy zgłasza koniec tury okna wykonawczego drogą, którą zgłasza go warstwa rozmowy.
func (s *stanowiskoPetli) koniecTuryWykonawcy() bool {
	return s.petla.ZakonczTure(s.wykonawca, PowodWynik)
}

// Metoda koniecTuryKoordynatora zgłasza koniec tury okna koordynatora z podanym powodem zatrzymania albo ukończenia.
func (s *stanowiskoPetli) koniecTuryKoordynatora(powod string) bool {
	return s.petla.ZakonczTure(s.koordynator, powod)
}

// Metoda stan oddaje odpis licznika biegu koordynatora prowadzonego przez to stanowisko sprawdzianu jednostkowego.
func (s *stanowiskoPetli) stan() StanObiegu { return s.petla.Stan(s.koordynator) }

// TestUkonczenieZamykaBiegPoTurzeKoordynatoraZakonczonejWynikiem wykazuje warunek ukończenia biegu w całości, na turze koordynatora zamkniętej wynikiem.
func TestUkonczenieZamykaBiegPoTurzeKoordynatoraZakonczonejWynikiem(t *testing.T) {
	s := noweStanowisko(t)

	if !s.koniecTuryWykonawcy() {
		t.Fatal("koniec tury wykonawcy nie wybudził koordynatora — bieg nie ruszył")
	}
	if len(s.obiegi) != 1 {
		t.Fatalf("port rozpoczynania obiegu dostał %d obiegów, oczekiwano jednego", len(s.obiegi))
	}
	if przed := s.stan(); przed.Zatrzymany {
		t.Fatalf("bieg stanął przed turą koordynatora: powód %q", przed.Powod)
	}

	if !s.koniecTuryKoordynatora(PowodWynik) {
		t.Fatal("koniec tury koordynatora nie poruszył biegu")
	}

	po := s.stan()
	if !po.Zatrzymany {
		t.Fatal("bieg nie stanął po turze koordynatora zamkniętej wynikiem")
	}
	if po.Powod != ZatrzymanieUkonczenie {
		t.Fatalf("bieg stanął powodem %q, a nie ukończeniem z wynikiem", po.Powod)
	}
	if po.Obiegow != 1 {
		t.Errorf("ukończenie zgubiło licznik obiegów: %d", po.Obiegow)
	}

	// Ukończenie idzie do obserwatorów, a nie zostaje w liczniku, czytane tą samą drogą co zatrzymania.
	ostatni := s.slady[len(s.slady)-1]
	if ostatni.Stan.Powod != ZatrzymanieUkonczenie {
		t.Errorf("obserwator nie dostał ukończenia; ostatni ślad niesie powód %q",
			ostatni.Stan.Powod)
	}
	if ostatni.Rozpoczety {
		t.Error("ślad ukończenia melduje rozpoczęty obieg — ukończenie obiegu nie zaczyna")
	}
}

// TestTuraWykonawcyWBieguWstrzymujeUkonczenie wykazuje, że tura wykonawcy prowadzona w biegu wstrzymuje ukończenie mimo tury koordynatora zamkniętej wynikiem.
func TestTuraWykonawcyWBieguWstrzymujeUkonczenie(t *testing.T) {
	s := noweStanowisko(t)
	s.koniecTuryWykonawcy()

	// Fragment wykonawcy jest zgłoszeniem tury: pętla dowiaduje się z niego, że okno wykonawcze pracuje.
	s.petla.ObserwujFragment(protocol.ChunkTekstu(s.wykonawca, "wiadomosc-1", "liczę"))

	if s.koniecTuryKoordynatora(PowodWynik) {
		t.Fatal("bieg został ukończony, choć okno wykonawcze prowadziło turę")
	}
	if stan := s.stan(); stan.Zatrzymany {
		t.Fatalf("bieg stanął mimo tury wykonawcy: powód %q", stan.Powod)
	}

	// Koniec tury wykonawcy zdejmuje jego turę i wybudza koordynatora, dopiero potem tura kończy bieg.
	s.koniecTuryWykonawcy()
	if !s.koniecTuryKoordynatora(PowodWynik) {
		t.Fatal("bieg nie został ukończony po zamknięciu tury wykonawcy")
	}
	if stan := s.stan(); stan.Powod != ZatrzymanieUkonczenie {
		t.Fatalf("bieg stanął powodem %q, a nie ukończeniem", stan.Powod)
	}
}

// TestTuraKoordynatoraZamknietaInaczejNizWynikiemNieKonczyBiegu wykazuje, że
// ukończenie bierze się z WYNIKU tury, a nie z samego jej końca. Tura przerwana
// i tura zakończona błędem kończą się tak samo — i żadna nie jest ukończeniem.
func TestTuraKoordynatoraZamknietaInaczejNizWynikiemNieKonczyBiegu(t *testing.T) {
	przypadki := map[string]string{
		"tura przerwana":    "zatrzymanie tury",
		"tura z błędem":     "błąd tury",
		"powód nienazwany":  "",
		"powód spoza pętli": "koniec procesu okna",
	}
	for nazwa, powod := range przypadki {
		t.Run(nazwa, func(t *testing.T) {
			s := noweStanowisko(t)
			s.koniecTuryWykonawcy()

			if s.koniecTuryKoordynatora(powod) {
				t.Fatalf("tura zamknięta powodem %q ukończyła bieg", powod)
			}
			if stan := s.stan(); stan.Zatrzymany {
				t.Fatalf("bieg stanął po turze zamkniętej powodem %q: %q", powod, stan.Powod)
			}
		})
	}
}

// TestTuraKoordynatoraPrzedPierwszymObiegiemNieKonczyBiegu wykazuje, że
// ukończenie zamyka bieg, a nie okno: koordynator, który jeszcze ani razu nie
// obiegł pętli, nie zapala kontrolki biegu zatrzymanego.
func TestTuraKoordynatoraPrzedPierwszymObiegiemNieKonczyBiegu(t *testing.T) {
	s := noweStanowisko(t)

	if s.koniecTuryKoordynatora(PowodWynik) {
		t.Fatal("tura koordynatora przed pierwszym obiegiem ukończyła bieg, którego nie ma")
	}
	if stan := s.stan(); stan.Zatrzymany {
		t.Fatalf("bieg bez ani jednego obiegu stoi zatrzymany powodem %q", stan.Powod)
	}
}

// TestUkonczenieNieJestBramkaAkceptacji wykazuje, że bieg ukończony podejmuje się sam, gdy praca dostaje ciąg dalszy, bez potwierdzenia i bez wznowienia.
func TestUkonczenieNieJestBramkaAkceptacji(t *testing.T) {
	s := noweStanowisko(t)
	s.koniecTuryWykonawcy()
	s.koniecTuryKoordynatora(PowodWynik)

	if stan := s.stan(); stan.Powod != ZatrzymanieUkonczenie {
		t.Fatalf("bieg nie został ukończony: %q", stan.Powod)
	}
	obiegowPoUkonczeniu := len(s.obiegi)

	// Praca ma ciąg dalszy: kolejna tura wykonawcy kończy się i wybudza koordynatora bez wznowienia.
	if !s.koniecTuryWykonawcy() {
		t.Fatal("koniec tury wykonawcy po ukończeniu nie wybudził koordynatora")
	}
	if len(s.obiegi) != obiegowPoUkonczeniu+1 {
		t.Fatalf("bieg ukończony nie podjął pracy: obiegów %d, przed zgłoszeniem %d — "+
			"ukończenie zachowuje się jak bramka akceptacji",
			len(s.obiegi), obiegowPoUkonczeniu)
	}
	po := s.stan()
	if po.Zatrzymany {
		t.Errorf("bieg po podjęciu nadal stoi zatrzymany powodem %q", po.Powod)
	}
	if po.Obiegow != 2 {
		t.Errorf("podjęcie biegu zgubiło historię obiegów: %d zamiast 2", po.Obiegow)
	}
}

// TestTrzyZatrzymaniaCzekajaNaOperatoraAUkonczenieNie wykazuje różnicę zachowań: trzy powody zatrzymania czekają na Operatora, a ukończenie podejmuje się samo.
func TestTrzyZatrzymaniaCzekajaNaOperatoraAUkonczenieNie(t *testing.T) {
	t.Run("zatrzymanie przez Operatora", func(t *testing.T) {
		s := noweStanowisko(t)
		s.koniecTuryWykonawcy()

		if stan := s.petla.Zatrzymaj(s.koordynator); stan.Powod != ZatrzymanieRecznie {
			t.Fatalf("bieg stanął powodem %q", stan.Powod)
		}
		sprawdzZatrzymanieTrzyma(t, s, ZatrzymanieRecznie)
	})

	t.Run("brak postępu wykonawcy", func(t *testing.T) {
		s := noweStanowisko(t)
		// Wykonawca powtarza się co do znaku: strumień jest pusty przy obiegu, odcisk pozostaje ten sam.
		for obieg := 0; obieg <= ProgBrakuPostepuDomyslny; obieg++ {
			s.koniecTuryWykonawcy()
		}
		if stan := s.stan(); stan.Powod != ZatrzymanieBrakPostepu {
			t.Fatalf("bieg stanął powodem %q, a nie brakiem postępu", stan.Powod)
		}
		sprawdzZatrzymanieTrzyma(t, s, ZatrzymanieBrakPostepu)
	})

	t.Run("usterka rozpoczęcia obiegu", func(t *testing.T) {
		s := noweStanowisko(t)
		s.usterka = errors.New("kanał odmówił tury")

		s.koniecTuryWykonawcy()
		if stan := s.stan(); stan.Powod != ZatrzymanieUsterka {
			t.Fatalf("bieg stanął powodem %q, a nie usterką", stan.Powod)
		}
		s.usterka = nil
		sprawdzZatrzymanieTrzyma(t, s, ZatrzymanieUsterka)
	})

	t.Run("ukończenie z wynikiem", func(t *testing.T) {
		s := noweStanowisko(t)
		s.koniecTuryWykonawcy()
		s.koniecTuryKoordynatora(PowodWynik)
		if stan := s.stan(); stan.Powod != ZatrzymanieUkonczenie {
			t.Fatalf("bieg stanął powodem %q, a nie ukończeniem", stan.Powod)
		}

		przed := len(s.obiegi)
		s.koniecTuryWykonawcy()
		if len(s.obiegi) == przed {
			t.Fatal("bieg ukończony nie podjął pracy bez wznowienia — " +
				"ukończenie nie odróżnia się od zatrzymania")
		}
	})
}

// sprawdzZatrzymanieTrzyma wykazuje, że bieg zatrzymany podanym powodem NIE
// rusza sam, a podejmuje go dopiero wznowienie — decyzja Operatora.
func sprawdzZatrzymanieTrzyma(t *testing.T, s *stanowiskoPetli, powod PowodZatrzymania) {
	t.Helper()

	przed := len(s.obiegi)
	s.koniecTuryWykonawcy()
	if len(s.obiegi) != przed {
		t.Fatalf("bieg zatrzymany powodem %q ruszył sam — zatrzymanie nie trzyma", powod)
	}
	if stan := s.stan(); !stan.Zatrzymany || stan.Powod != powod {
		t.Fatalf("bieg po zgłoszeniu stoi w stanie zatrzymany=%v powód=%q",
			stan.Zatrzymany, stan.Powod)
	}

	s.petla.Wznow(s.koordynator)
	s.koniecTuryWykonawcy()
	if len(s.obiegi) == przed {
		t.Fatalf("bieg zatrzymany powodem %q nie ruszył po wznowieniu", powod)
	}
}
