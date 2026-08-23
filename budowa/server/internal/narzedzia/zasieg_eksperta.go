// Odpowiedzialność pliku: trzeci zasięg — narzędzia tego eksperta.
//
// Zasięg eksperta to podzbiór wykazu kontraktu wskazany przez definicję eksperta,
// a nie wykaz okna (dwa pozostałe zasięgi wynikają z roli okna — `zasieg_roli.go`
// — i tylko dokładają). Zawężenie jest potrzebne, bo pełny wykaz kontraktu
// złożony w kształt `tools/list` waży rzędu dziesiątek tysięcy żetonów zjadanych,
// zanim padnie pierwsze słowo zadania; przy nadmiarze narzędzia docierają do
// modelu bez opisów — model widzi same nazwy i po nie nie sięga. Cała ta droga
// istnieje po to, żeby cichej degradacji nie było, więc sama nie ma prawa
// degradować po cichu (`ekspert_wykaz.go`).
//
// Kod eksperta przychodzi przełącznikiem `--ekspert`, czytanym raz przy
// uruchomieniu, z wpisu `danaco` ułożonego przez rdzeń. Model nie ma czym o
// zasięg poprosić ani go poszerzyć: nie ma narzędzia zmieniającego zasięg, a wpis
// powstaje zanim proces modelu wystartuje. Rdzeń bierze kod z pola
// `session.Okno.Ustawienia.Agent` — tego samego, które nakłada eksperta na okno.
//
// `core/sprawca.go` rozdziela dwie ręce po zasięgu: klawiatura Operatora to
// asystent, wszystko inne to model roboczy. Zasięg eksperta spada tam na model
// roboczy: okno eksperta jest ręką modelu wykonującego zlecenie, nie ręką
// Operatora.
package narzedzia

const (
	// ZasiegEksperta jest zasięgiem okna z nałożonym ekspertem: podzbiór wykazu
	// wskazany przez definicję eksperta, a nie wykaz okna.
	ZasiegEksperta Zasieg = "ekspert"
	// PrzelacznikEksperta jest nazwą przełącznika niosącego kod eksperta.
	// Osobno od `--zasieg`, bo niesie wartość, a nie nazwę roli.
	PrzelacznikEksperta = "ekspert"
)

// ArgumentyEksperta zwraca argumenty uruchomienia dopisywane do wpisu `danaco`
// dla okna z nałożonym ekspertem: nazwa zasięgu i kod eksperta razem.
//
// Kod pusty nie daje żadnych argumentów — ani nazwy zasięgu. Zasięg eksperta
// bez eksperta nie jest zasięgiem węższym, tylko zasięgiem bez treści; wpis
// z samą nazwą roli kazałby serwerowi meldować brak przy każdym `tools/list`
// zamiast po prostu pracować w zasięgu okna.
func ArgumentyEksperta(kod string) []string {
	if kod == "" {
		return nil
	}
	return []string{
		"--" + PrzelacznikZasiegu, string(ZasiegEksperta),
		"--" + PrzelacznikEksperta, kod,
	}
}
