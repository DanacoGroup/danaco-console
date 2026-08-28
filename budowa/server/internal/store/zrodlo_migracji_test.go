package store

import (
	"strings"
	"testing"
)

// TestNazwaMigracjiPozaWzorcemJestOdmowa sprawdza rozbiór nazwy pliku: krok o nazwie spoza wzorca nie ma jak dostać numeru wersji.
func TestNazwaMigracjiPozaWzorcemJestOdmowa(t *testing.T) {
	przypadki := []string{
		"schemat.sql",
		"migracja_.sql",
		"migracja_bezcyfr_nazwa.sql",
		"migracja_012.sql",
	}
	for _, nazwa := range przypadki {
		t.Run(nazwa, func(t *testing.T) {
			if _, err := zbudujKrok(nazwa); err == nil {
				t.Errorf("nazwa %q spoza wzorca została przyjęta", nazwa)
			}
		})
	}
}

// TestDwaKrokiOTymSamymNumerzeSaOdmowa pilnuje jednoznaczności numeracji.
// Rejestr ma na kolumnie wersji warunek UNIQUE, więc drugi krok wykonałby swój
// schemat, ale nie zostałby odnotowany — i wykonałby się ponownie przy każdym
// kolejnym starcie.
func TestDwaKrokiOTymSamymNumerzeSaOdmowa(t *testing.T) {
	kroki := []migracja{
		{Wersja: 1, Nazwa: "fundament"},
		{Wersja: 2, Nazwa: "okna"},
		{Wersja: 2, Nazwa: "okna_powtorzone"},
	}
	err := sprawdzUnikalnoscWersji(kroki)
	if err == nil {
		t.Fatal("dwa kroki o tym samym numerze zostały przyjęte")
	}
	if !strings.Contains(err.Error(), "okna_powtorzone") {
		t.Errorf("odmowa nie nazywa kroku powtórzonego: %v", err)
	}
}

// TestNumeracjaMozeMiecLuki utrwala regułę odwrotną: ciągłość numeracji nie jest wymagana, a luki powstają przy pracy równoległej.
func TestNumeracjaMozeMiecLuki(t *testing.T) {
	kroki := []migracja{
		{Wersja: 1, Nazwa: "fundament"},
		{Wersja: 5, Nazwa: "okna"},
		{Wersja: 93, Nazwa: "dziennik"},
	}
	if err := sprawdzUnikalnoscWersji(kroki); err != nil {
		t.Errorf("numeracja z lukami została odrzucona: %v", err)
	}
}

// TestKrokZasobuNiesieSumeTresci sprawdza, że suma kontrolna liczy się z treści
// kroku, a nie z jego nazwy — inaczej sprawdzenie niezmienności nie zauważyłoby
// podmienionej treści pod tą samą nazwą.
func TestKrokZasobuNiesieSumeTresci(t *testing.T) {
	kroki, err := wczytajMigracje()
	if err != nil {
		t.Fatalf("nie można odczytać wykazu migracji: %v", err)
	}

	sumy := make(map[string]string, len(kroki))
	for _, krok := range kroki {
		if len(krok.SumaKontrolna) != 64 {
			t.Errorf("krok %03d (%s) ma sumę o długości %d, oczekiwane 64 znaki zapisu szesnastkowego",
				krok.Wersja, krok.Nazwa, len(krok.SumaKontrolna))
		}
		if wczesniejszy, powtorzona := sumy[krok.SumaKontrolna]; powtorzona {
			t.Errorf("kroki %q i %q mają tę samą sumę kontrolną — dwa razy ta sama treść",
				wczesniejszy, krok.Nazwa)
		}
		sumy[krok.SumaKontrolna] = krok.Nazwa
	}
}
