// Zasięg eksperta jest trzecim zasięgiem: podzbiorem wykazu kontraktu
// wskazanym przez definicję eksperta, zawężonym po to, żeby wykaz nie ważył
// zbyt wiele żetonów.
package narzedzia

const (
	// ZasiegEksperta jest zasięgiem okna z nałożonym ekspertem: podzbiór wykazu
	// wskazany przez definicję eksperta, a nie wykaz okna.
	ZasiegEksperta Zasieg = "ekspert"
	// PrzelacznikEksperta jest nazwą przełącznika niosącego kod eksperta.
	// Osobno od `--zasieg`, bo niesie wartość, a nie nazwę roli.
	PrzelacznikEksperta = "ekspert"
)

// ArgumentyEksperta zwraca argumenty uruchomienia dopisywane do wpisu danaco
// dla okna z nałożonym ekspertem; kod pusty nie daje żadnych argumentów.
func ArgumentyEksperta(kod string) []string {
	if kod == "" {
		return nil
	}
	return []string{
		"--" + PrzelacznikZasiegu, string(ZasiegEksperta),
		"--" + PrzelacznikEksperta, kod,
	}
}
