// Odpowiedzialność pliku: przemiatanie retencji historii przy starcie rdzenia —
// zastosowanie zasady przechowywania (tabela `zasada_przechowywania`) do
// wszystkich okien, także tych, których nikt nie czyta.
//
// Pozostałe drogi egzekwowania retencji działają przy okazji innej pracy: po
// `retention.set` na oknach zakresu oraz przed oddaniem wykazu w `history.load`
// (`core/handlers_historia.go`). Okno, którego nikt nie otwiera i na którym
// nikt nie przestawia zasady, żadnej z nich nie napotyka — przemiatanie
// startowe jest jedynym miejscem, w którym zasada globalna sięga i tam.
//
// Start jest jedynym momentem, w którym rdzeń i tak przechodzi po stanie
// trwałym, więc nie potrzeba budzika ani wątku; tą samą drogą sprząta kosz
// sesji (`trwalosc_kosza.go`). Osobny zegar tylko dla retencji byłby drugim
// mechanizmem sprzątania obok już istniejącego.
//
// Każde przemiatanie, które coś zabrało, zostawia zdanie w dzienniku rdzenia —
// wzorem czyszczenia kosza. Ubytek bez śladu jest nie do odróżnienia od utraty
// danych.
package core

import (
	"context"
	"log"

	"danacoconsole/server/internal/dane"
)

// przemiecRetencjeHistorii stosuje zasadę przechowywania do wszystkich okien
// rozmowy. Wywoływane raz, przy starcie rdzenia. Niepowodzenie nie przerywa
// montażu — przemiatanie poczeka do następnego startu.
//
// Przemiatanie nie zna progów i ich nie zgaduje: pyta o nie `Egzekwuj`, a ten
// dla okna bez zasady oraz dla zasady bez progów nie usuwa niczego i zwraca
// zero (`dane/historia_retencja.go`). Brak zasady znaczy „trzymaj wszystko",
// nie „skasuj wszystko".
func przemiecRetencjeHistorii(kontekst context.Context, repozytoria *dane.Zestaw, dziennik *log.Logger) {
	if repozytoria == nil || repozytoria.Historia == nil {
		return
	}
	// Zakres global wylicza wszystkie okna z identyfikatorem kontraktowym —
	// zasadę obowiązującą każde z nich rozstrzyga potem `Egzekwuj` osobno,
	// kolejnością window → session → global. Przemiatanie tej
	// kolejności nie zna i nie powiela.
	okna, err := repozytoria.Historia.OknaZakresu(kontekst, "global", "")
	if err != nil {
		if dziennik != nil {
			dziennik.Printf("przemiatanie retencji historii: %v", err)
		}
		return
	}
	przyciete, oknaPrzyciete := 0, 0
	for _, oknoKod := range okna {
		usuniete, err := repozytoria.Historia.Egzekwuj(kontekst, oknoKod)
		if err != nil {
			// Okno oporne nie zatrzymuje przemiatania pozostałych, ale nie
			// przechodzi też bez śladu.
			if dziennik != nil {
				dziennik.Printf("przemiatanie retencji historii, okno %s: %v", oknoKod, err)
			}
			continue
		}
		if usuniete > 0 {
			przyciete += usuniete
			oknaPrzyciete++
		}
	}
	if przyciete > 0 && dziennik != nil {
		dziennik.Printf("retencja historii: przycięto trwale %d pozycji w %d oknach (przemiatanie startowe, %d okien sprawdzonych)",
			przyciete, oknaPrzyciete, len(okna))
	}
}
