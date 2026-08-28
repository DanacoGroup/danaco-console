package models

import (
	"context"
	"fmt"
	"strings"
)

// KanalEcho jest kanałem bez sieci i bez procesu: odsyła treść zapytania
// porcjami, tak jak zrobiłby to model. Służy próbie całej drogi — rejestr,
// zapytanie, prowenancja, strumień, ujście — bez zależności od dostawcy,
// konta i łącza.
type KanalEcho struct {
	def         Definicja
	przedrostek string
	porcja      int
}

// domyslnaPorcjaEcho podaje, ile znaków niesie jeden fragment tekstu, gdy
// wiersz rejestru nie mówi inaczej.
const domyslnaPorcjaEcho = 24

// NowyKanalEcho buduje kanał echo z wiersza rejestru. Parametry wiersza:
// "przedrostek" — tekst doklejany przed odpowiedzią, "porcja" — długość
// jednego fragmentu tekstu.
func NowyKanalEcho(d Definicja) (Kanal, error) {
	porcja := domyslnaPorcjaEcho
	if wskazana := d.Parametr("porcja"); wskazana != "" {
		if _, err := fmt.Sscanf(wskazana, "%d", &porcja); err != nil || porcja <= 0 {
			porcja = domyslnaPorcjaEcho
		}
	}
	return &KanalEcho{def: d, przedrostek: d.Parametr("przedrostek"), porcja: porcja}, nil
}

// Kod zwraca kod kanału z wiersza rejestru, identyfikujący ten kanał wśród
// wszystkich kanałów echo dostawcy.
func (k *KanalEcho) Kod() string {
	return k.def.Kod
}

// Definicja zwraca wiersz rejestru, z którego kanał echo powstał, wraz
// z jego pełnymi parametrami konfiguracji.
func (k *KanalEcho) Definicja() Definicja {
	return k.def
}

// Wyslij nadaje prowenancję wywołania, następnie metadane konta, a dopiero po
// nich tekst odpowiedzi porcjami. Kolejność jest częścią kontraktu strumienia:
// odbiorca poznaje warunki wywołania, zanim zobaczy pierwszy znak odpowiedzi.
func (k *KanalEcho) Wyslij(ctx context.Context, z Zapytanie, u Ujscie) error {
	prowenancja := ProwenancjaZapytania(z, k.def)
	prowenancja.Adres = "echo://" + k.def.Kod
	// Program niesie już klucz adaptera echo, więc argv dokłada wyłącznie
	// model.
	prowenancja.Argv = []string{prowenancja.Model}
	if err := NadajProwenancje(ctx, u, z, prowenancja); err != nil {
		return err
	}
	if wybrane := z.WybraneKonto(k.def); k.def.PoswiadczenieOdwolanie != "" || wybrane != "" {
		konto := MetadaneKonta{Konto: z.WybraneKonto(k.def), Odwolanie: k.def.PoswiadczenieOdwolanie, Powod: "wywołanie kanału echo"}
		if err := NadajKonto(ctx, u, z, konto); err != nil {
			return err
		}
	}
	return k.nadajTresc(ctx, z, u)
}

// nadajTresc dzieli odpowiedź na porcje i nadaje je kolejno, sprawdzając przed
// każdą, czy wywołanie nie zostało odwołane.
func (k *KanalEcho) nadajTresc(ctx context.Context, z Zapytanie, u Ujscie) error {
	tresc := k.przedrostek + z.Tresc
	if strings.TrimSpace(tresc) == "" {
		return nil
	}
	znaki := []rune(tresc)
	for poczatek := 0; poczatek < len(znaki); poczatek += k.porcja {
		if err := ctx.Err(); err != nil {
			return err
		}
		koniec := poczatek + k.porcja
		if koniec > len(znaki) {
			koniec = len(znaki)
		}
		if err := NadajTekst(ctx, u, z, string(znaki[poczatek:koniec])); err != nil {
			return err
		}
	}
	return nil
}
