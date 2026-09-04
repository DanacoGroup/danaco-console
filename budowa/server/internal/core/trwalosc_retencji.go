// Plik przemiata retencję historii przy starcie rdzenia, stosując zasadę przechowywania do wszystkich okien, także tych, których nikt nie czyta. Start jest jedynym momentem, w którym rdzeń i tak przechodzi po stanie trwałym.
package core

import (
	"context"
	"log"

	"danacoconsole/server/internal/dane"
)

// przemiecRetencjeHistorii stosuje zasadę przechowywania do wszystkich okien rozmowy, raz przy starcie rdzenia. Niepowodzenie nie przerywa montażu. Przemiatanie nie zna progów, pyta o nie Egzekwuj, a brak zasady znaczy trzymaj wszystko.
func przemiecRetencjeHistorii(kontekst context.Context, repozytoria *dane.Zestaw, dziennik *log.Logger) {
	if repozytoria == nil || repozytoria.Historia == nil {
		return
	}
	// Przemiatanie idzie po kontach po kolei (decyzja 34): wykaz okien jest
	// zawężony kontem, więc jeden przebieg objąłby wyłącznie konto najstarsze.
	konteksty, err := kontekstyKont(kontekst, repozytoria.KontoWlasciciela)
	if err != nil {
		if dziennik != nil {
			dziennik.Printf("przemiatanie retencji historii, odczyt kont: %v", err)
		}
		return
	}
	przyciete, oknaPrzyciete, sprawdzonych := 0, 0, 0
	for _, kontekstKonta := range konteksty {
		p, o, s := przemiecRetencjeKonta(kontekstKonta, repozytoria, dziennik)
		przyciete += p
		oknaPrzyciete += o
		sprawdzonych += s
	}
	if przyciete > 0 && dziennik != nil {
		dziennik.Printf("retencja historii: przycięto trwale %d pozycji w %d oknach (przemiatanie startowe, %d okien sprawdzonych)",
			przyciete, oknaPrzyciete, sprawdzonych)
	}
}

// przemiecRetencjeKonta przemiata okna jednego konta i oddaje liczby do sumy.
func przemiecRetencjeKonta(kontekst context.Context, repozytoria *dane.Zestaw,
	dziennik *log.Logger) (przyciete, oknaPrzyciete, sprawdzonych int) {

	// Zakres global wylicza wszystkie okna kontraktowe; kolejność zasad rozstrzyga Egzekwuj osobno.
	okna, err := repozytoria.Historia.OknaZakresu(kontekst, "global", "")
	if err != nil {
		if dziennik != nil {
			dziennik.Printf("przemiatanie retencji historii: %v", err)
		}
		return 0, 0, 0
	}
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
	return przyciete, oknaPrzyciete, len(okna)
}
