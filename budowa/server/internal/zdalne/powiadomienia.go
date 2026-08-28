// Silnik kolejki powiadomień prowadzi cztery wejścia od faktu, o którym operator ma wiedzieć, do doręczenia.
package zdalne

import (
	"context"
	"fmt"
	"sync"
	"time"

	"danacoconsole/server/internal/dane"
)

// FormatZnacznika jest formatem znacznika czasu bazy — trzy cyfry części ułamkowej, zgodnie z zapisem bazy.
const FormatZnacznika = "2006-01-02T15:04:05.000Z"

// odstepyPonowien to rosnące odstępy między podejściami, liczone od numeru próby już odbytej danego powiadomienia.
var odstepyPonowien = []time.Duration{
	30 * time.Second,
	time.Minute,
	2 * time.Minute,
	5 * time.Minute,
	10 * time.Minute,
	15 * time.Minute,
}

// nadawanie trzyma nadajnik podany przez kompozycję. Stoi osobno od `zasilenie`,
// bo uchwyt bazy i droga doręczenia to dwa niezależne wpięcia.
var nadawanie struct {
	sync.RWMutex
	nadajnik Nadajnik
}

// ZasilNadajnik podaje pakietowi drogę doręczenia powiadomień. Wywołanie należy
// do kompozycji, obok `zdalne.Zasil`.
func ZasilNadajnik(n Nadajnik) {
	nadawanie.Lock()
	defer nadawanie.Unlock()
	nadawanie.nadajnik = n
}

// Funkcja OdetnijNadajnik zdejmuje nadajnik zasilony wcześniej — dla porządku sprawdzianów jednostkowych.
func OdetnijNadajnik() {
	ZasilNadajnik(nil)
}

func nadajnik() Nadajnik {
	nadawanie.RLock()
	defer nadawanie.RUnlock()
	return nadawanie.nadajnik
}

// powiadomienia składa repozytorium z uchwytu bazy, który pakiet już ma.
// Brak zasilenia jest odmową nazwaną, wspólną dla wszystkich dróg pakietu.
func powiadomienia() (dane.RepozytoriumPowiadomien, error) {
	db := baza()
	if db == nil {
		return nil, odmowaBrakuZasilenia()
	}
	return dane.NowePowiadomienia(db), nil
}

// Zgloszenie jest tym, co wołający wnosi do kolejki powiadomień przy jego pierwotnym pełnym zakładaniu.
type Zgloszenie struct {
	// Tytul jest zdaniem, które Operator zobaczy pierwsze.
	Tytul string
	// Tresc dopowiada; może być pusta.
	Tresc string
	// Priorytet: `zwykly` albo `pilny`. Puste bierze `zwykly`.
	Priorytet string
	// Powod mówi, po co dzwonimy; powiadomienie bez powodu jest budzikiem bez oceny operatora.
	Powod string
	// BytRodzaj i BytID wskazują, czego dotyczy powiadomienie; kotwica jest miękka, bez klucza obcego.
	BytRodzaj string
	BytID     string
	// Waznosc jest czasem, po którym powiadomienie przestaje mieć sens; pole jest obowiązkowe.
	Waznosc time.Duration
}

// Zglos wnosi powiadomienie do kolejki i oddaje jego klucz. Sama droga nie
// doręcza — wysyłką zajmuje się takt (Wyslij). Gdyby zgłoszenie doręczało od
// razu, decyzja podjęta w tej samej chwili nie zdążyłaby powiadomienia odwołać.
func Zglos(ctx context.Context, z Zgloszenie) (int64, error) {
	if z.Tytul == "" {
		return 0, fmt.Errorf("zdalne: powiadomienie bez tytułu nie weszło do kolejki — " +
			"Operator zobaczyłby wołanie, którego nie da się przeczytać")
	}
	if z.Waznosc <= 0 {
		return 0, fmt.Errorf("zdalne: powiadomienie %q nie weszło do kolejki, bo nie ma "+
			"terminu ważności — wołający musi odpowiedzieć, do kiedy ta sprawa ma sens; "+
			"powiadomienie bez terminu wisiałoby w kolejce na zawsze", z.Tytul)
	}
	repozytorium, err := powiadomienia()
	if err != nil {
		return 0, err
	}

	teraz := time.Now().UTC()
	p := dane.Powiadomienie{
		Tytul:         z.Tytul,
		Tresc:         z.Tresc,
		Priorytet:     z.Priorytet,
		Powod:         z.Powod,
		NastepnaProba: teraz.Format(FormatZnacznika),
		Wygasa:        teraz.Add(z.Waznosc).Format(FormatZnacznika),
	}
	if z.BytRodzaj != "" && z.BytID != "" {
		rodzaj, klucz := z.BytRodzaj, z.BytID
		p.BytRodzaj = &rodzaj
		p.BytID = &klucz
	}
	return repozytorium.Wstaw(ctx, p)
}

// Funkcja ZarejestrujUrzadzenie zapisuje zgodę urządzenia na wołanie kanałem połączenia i oddaje klucz rejestracji.
func ZarejestrujUrzadzenie(ctx context.Context, urzadzenieID int64, kluczKanalu, etykieta string) (int64, error) {
	if urzadzenieID <= 0 {
		return 0, fmt.Errorf("zdalne: rejestracja powiadomień bez wskazania urządzenia — " +
			"adresatem może być dziś wyłącznie maszyna z wykazu `urzadzenie`, bo tabeli " +
			"użytkowników w tym schemacie nie ma")
	}
	if kluczKanalu == "" {
		return 0, fmt.Errorf("zdalne: rejestracja powiadomień urządzenia %d nie ma adresu — "+
			"dla kanału %q adresem jest identyfikator klienta z powitania gniazda",
			urzadzenieID, dane.KanalPolaczenie)
	}
	repozytorium, err := powiadomienia()
	if err != nil {
		return 0, err
	}
	rejestracja := dane.RejestracjaPowiadomien{
		UrzadzenieID: urzadzenieID,
		Kanal:        dane.KanalPolaczenie,
		KluczKanalu:  kluczKanalu,
	}
	if etykieta != "" {
		nazwa := etykieta
		rejestracja.Etykieta = &nazwa
	}
	return repozytorium.Zarejestruj(ctx, rejestracja)
}

// Funkcja WyrejestrujUrzadzenie cofa zgodę urządzenia, zostawiając ślad po tym, że było ono kiedyś wołane.
func WyrejestrujUrzadzenie(ctx context.Context, kluczKanalu string) error {
	if kluczKanalu == "" {
		return fmt.Errorf("zdalne: wyrejestrowanie powiadomień bez wskazania adresu")
	}
	repozytorium, err := powiadomienia()
	if err != nil {
		return err
	}
	return repozytorium.Wyrejestruj(ctx, dane.KanalPolaczenie, kluczKanalu,
		time.Now().UTC().Format(FormatZnacznika))
}

// Funkcja Odwolaj gasi wszystkie oczekujące powiadomienia o wskazanym bycie i oddaje ich łączną liczbę.
func Odwolaj(ctx context.Context, bytRodzaj, bytID, powod string) (int, error) {
	if bytRodzaj == "" || bytID == "" {
		return 0, fmt.Errorf("zdalne: odwołanie powiadomień bez wskazania bytu nie ma " +
			"adresata — rodzaj i klucz bytu chodzą parą")
	}
	repozytorium, err := powiadomienia()
	if err != nil {
		return 0, err
	}
	return repozytorium.Odwolaj(ctx, bytRodzaj, bytID, powod,
		time.Now().UTC().Format(FormatZnacznika))
}

// Wygas zamyka powiadomienia, którym minął termin ważności, i oddaje ich liczbę.
// Droga jest osobna od Wyslij, bo termin ważności ma mijać także wtedy, gdy
// nadajnika nie ma.
func Wygas(ctx context.Context) (int, error) {
	repozytorium, err := powiadomienia()
	if err != nil {
		return 0, err
	}
	return repozytorium.Wygas(ctx, time.Now().UTC().Format(FormatZnacznika))
}

// PodsumowanieWysylki zbiera wynik jednego przebiegu kolejki powiadomień, wykonanego przez takt rdzenia.
type PodsumowanieWysylki struct {
	// Rozpatrzonych to liczba powiadomień, które takt wziął pod zamek.
	Rozpatrzonych int
	// Doreczonych to liczba powiadomień zamkniętych doręczeniem.
	Doreczonych int
	// Pominietych to liczba powiadomień, które w chwili zamka nie czekały już na wysyłkę.
	Pominietych int
	// Przelozonych to liczba powiadomień bez odbiorcy, odłożonych na później.
	Przelozonych int
	// Wygaslych to liczba powiadomień, którym termin minął w trakcie taktu.
	Wygaslych int
}

// Funkcja Wyslij wykonuje jeden przebieg kolejki na wskazaną chwilę, doręczając powiadomienia przez nadajnik.
func Wyslij(ctx context.Context, teraz time.Time) (PodsumowanieWysylki, error) {
	nad := nadajnik()
	if nad == nil {
		return PodsumowanieWysylki{}, fmt.Errorf("zdalne: przebieg kolejki powiadomień nie " +
			"ruszył, bo kompozycja nie podała nadajnika (zdalne.ZasilNadajnik) — " +
			"powiadomienia zostają w kolejce nietknięte, z pełnym budżetem prób, " +
			"i nic nie udaje wysłanego")
	}
	repozytorium, err := powiadomienia()
	if err != nil {
		return PodsumowanieWysylki{}, err
	}

	znacznik := teraz.UTC().Format(FormatZnacznika)
	nalezne, err := repozytorium.Nalezne(ctx, znacznik)
	if err != nil {
		return PodsumowanieWysylki{}, err
	}
	if len(nalezne) == 0 {
		return PodsumowanieWysylki{}, nil
	}

	cele, err := celeWolania(ctx, repozytorium)
	if err != nil {
		return PodsumowanieWysylki{}, err
	}

	podsumowanie := PodsumowanieWysylki{}
	for _, p := range nalezne {
		wynik, err := repozytorium.Takt(ctx, p.ID, dane.PlanTaktu{
			Teraz:         znacznik,
			NastepnaProba: func(prob int) string { return nastepnaProba(teraz, prob) },
			Wyslij: func(ctx context.Context, wiersz dane.Powiadomienie) ([]int64, error) {
				return rozeslij(ctx, nad, cele, wiersz), nil
			},
		})
		if err != nil {
			// Awaria taktu przerywa przebieg; podsumowanie oddaje to, co
			// zdążyło się rozstrzygnąć przed błędem.
			return podsumowanie, err
		}
		podsumowanie.Rozpatrzonych++
		switch {
		case wynik.Pominieto:
			podsumowanie.Pominietych++
		case wynik.Stan == dane.StanPowiadomieniaDostarczone:
			podsumowanie.Doreczonych++
		case wynik.Stan == dane.StanPowiadomieniaWygasle:
			podsumowanie.Wygaslych++
		default:
			podsumowanie.Przelozonych++
		}
	}
	return podsumowanie, nil
}

// Funkcja celeWolania składa wykaz adresów z czynnych rejestracji, odczytywany raz na cały przebieg kolejki.
func celeWolania(ctx context.Context, repozytorium dane.RepozytoriumPowiadomien) ([]Cel, error) {
	rejestracje, err := repozytorium.AktywneRejestracje(ctx)
	if err != nil {
		return nil, err
	}
	cele := make([]Cel, 0, len(rejestracje))
	for _, r := range rejestracje {
		cel := Cel{Rejestracja: r.ID, Kanal: r.Kanal, KluczKanalu: r.KluczKanalu}
		if r.Etykieta != nil {
			cel.Etykieta = *r.Etykieta
		}
		cele = append(cele, cel)
	}
	return cele, nil
}

// rozeslij woła każdy czynny adres i oddaje klucze rejestracji, które kopertę
// przyjęły. Urządzenie nieobecne odpowiada fałszem — to nie jest błąd przebiegu,
// tylko powód ponowienia.
func rozeslij(ctx context.Context, nad Nadajnik, cele []Cel, p dane.Powiadomienie) []int64 {
	przyjeli := make([]int64, 0, len(cele))
	for _, cel := range cele {
		if nad.Zawolaj(ctx, cel, p) {
			przyjeli = append(przyjeli, cel.Rejestracja)
		}
	}
	return przyjeli
}

// Funkcja nastepnaProba wylicza termin kolejnego podejścia z liczby prób już odbytych przez to powiadomienie.
func nastepnaProba(teraz time.Time, probIle int) string {
	odstep := odstepyPonowien[len(odstepyPonowien)-1]
	if probIle >= 1 && probIle <= len(odstepyPonowien) {
		odstep = odstepyPonowien[probIle-1]
	}
	return teraz.UTC().Add(odstep).Format(FormatZnacznika)
}
