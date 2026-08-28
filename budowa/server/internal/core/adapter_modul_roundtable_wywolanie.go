// Odpowiedzialność pliku: wywołania kanału modelu, które NIE są głosem
// uczestnika — analiza zapisu, werdykt sędziego i powtórzenie wypowiedzi.
// Wynik nie jest niczyim głosem, więc nie wchodzi do transkryptu jako
// wypowiedź i nie idzie strumieniem.
package core

import (
	"context"
	"strings"

	"danacoconsole/server/internal/dane"
	"danacoconsole/server/internal/models"
	"danacoconsole/shared"
)

// wywolajModelDebaty prowadzi jedno wywołanie kanału i oddaje całą odpowiedź.
// Ujście zbiera tekst zamiast go rozgłaszać: odbiorcą jest rdzeń, nie okno.
func (a *adapterDebaty) wywolajModelDebaty(ctx context.Context,
	okno, kanal, promptSystemowy, tresc string) (string, error) {

	if a.kanaly == nil {
		return "", bladBrakuKanalow()
	}
	kanal = strings.TrimSpace(kanal)
	if kanal == "" {
		return "", bladWskazaniaDebaty("wywołanie modelu bez wskazania kanału")
	}
	if _, jest := a.kanaly.Kanal(kanal); !jest {
		return "", bladNieznanegoKanalu(kanal)
	}

	var odpowiedz strings.Builder
	ujscie := models.UjscieFunkcji(func(_ context.Context, f models.Fragment) error {
		if f.Kind == shared.ChunkKindText {
			odpowiedz.WriteString(models.TrescFragmentu(f))
		}
		return nil
	})
	zapytanie := models.Zapytanie{
		Zasiegi:   models.Zasiegi{Okno: okno},
		Wiadomosc: nowyIdentyfikator(przedrostekWypowiedzi),
		Tresc:     tresc,
		Kanal:     kanal,
		Nakladka:  models.Nakladka{ProfilRoli: promptSystemowy},
	}
	if err := a.kanaly.Wyslij(ctx, zapytanie, ujscie); err != nil {
		return "", odmowaKanaluDebaty(kanal, err)
	}
	return odpowiedz.String(), nil
}

// kanalAnalizy rozstrzyga, który kanał wykonuje pracę nad zapisem debaty.
// Kolejność jest rozmyślna: wskazanie z żądania, potem kanał pierwszego
// uczestnika składu.
func (a *adapterDebaty) kanalAnalizy(ctx context.Context, okno string, wskazany *string) (string, error) {
	if kanal := strings.TrimSpace(wartoscTekstu(wskazany)); kanal != "" {
		return kanal, nil
	}
	uczestnicy, err := a.repozytorium.Uczestnicy(ctx, okno)
	if err != nil {
		return "", bladDebaty(err)
	}
	if len(uczestnicy) == 0 {
		return "", bladWskazaniaDebaty(
			"debata okna " + okno + " nie ma uczestnika, którego kanałem można wykonać analizę — " +
				"dodaj model w Model Panels albo wskaż kanał w żądaniu (channelId)")
	}
	return uczestnicy[0].KanalModelu, nil
}

// trescPonownegoGlosu wywołuje kanał uczestnika po raz drugi nad tym samym
// pytaniem tury i oddaje samą treść.
func (a *adapterDebaty) trescPonownegoGlosu(ctx context.Context, tura dane.TuraDebaty,
	uczestnik dane.UczestnikDebaty, kodWypowiedzi string) string {

	zapytanie := zapytanieUczestnika(tura, uczestnik, kodWypowiedzi, tura.Pytanie, "")
	var tresc strings.Builder
	ujscie := models.UjscieFunkcji(func(_ context.Context, f models.Fragment) error {
		if f.Kind == shared.ChunkKindText {
			tresc.WriteString(models.TrescFragmentu(f))
		}
		return nil
	})
	if err := a.kanaly.Wyslij(ctx, zapytanie, ujscie); err != nil {
		return ""
	}
	return tresc.String()
}

// zapisDebatyDoAnalizy składa materiał, nad którym pracuje model: tura po
// turze, mówca po mówcy. Wskazanie tury zawęża materiał do niej jednej.
func (a *adapterDebaty) zapisDebatyDoAnalizy(ctx context.Context,
	okno, turaKod string) (string, []dane.WypowiedzDebaty, error) {

	uczestnicy, err := a.repozytorium.Uczestnicy(ctx, okno)
	if err != nil {
		return "", nil, bladDebaty(err)
	}
	podpisy := podpisyUczestnikow(uczestnicy)

	wypowiedzi, err := a.wypowiedziZakresu(ctx, okno, turaKod)
	if err != nil {
		return "", nil, err
	}
	istotne := make([]dane.WypowiedzDebaty, 0, len(wypowiedzi))
	var zapis strings.Builder
	for _, wypowiedz := range wypowiedzi {
		if strings.TrimSpace(wypowiedz.Tresc) == "" {
			continue // uczestnik nie odpowiedział — pustki nie podaje się do analizy
		}
		istotne = append(istotne, wypowiedz)
		zapis.WriteString("[")
		zapis.WriteString(wypowiedz.Kod)
		zapis.WriteString("] ")
		zapis.WriteString(podpis(podpisy, wypowiedz.Uczestnik))
		zapis.WriteString(": ")
		zapis.WriteString(strings.TrimSpace(wypowiedz.Tresc))
		zapis.WriteString("\n\n")
	}
	return strings.TrimSpace(zapis.String()), istotne, nil
}

// wierszeOdpowiedzi rozbija odpowiedź modelu na pozycje wykazu. Model
// poproszony o wykaz oddaje go zwykle wierszami, czasem z myślnikiem albo
// numerem na początku.
func wierszeOdpowiedzi(odpowiedz string) []string {
	pozycje := make([]string, 0, 8)
	for _, wiersz := range strings.Split(odpowiedz, "\n") {
		przyciety := strings.TrimSpace(wiersz)
		if przyciety == "" {
			continue
		}
		przyciety = strings.TrimLeft(przyciety, "-*• \t")
		przyciety = zdejmijNumerPozycji(przyciety)
		if przyciety = strings.TrimSpace(przyciety); przyciety != "" {
			pozycje = append(pozycje, przyciety)
		}
	}
	if len(pozycje) == 0 {
		if calosc := strings.TrimSpace(odpowiedz); calosc != "" {
			return []string{calosc}
		}
		return nil
	}
	return pozycje
}

// zdejmijNumerPozycji usuwa wiodący numer wykazu („1.", „2)"). Ciąg cyfr, po
// którym nie ma kropki ani nawiasu, zostaje: bywa treścią.
func zdejmijNumerPozycji(wiersz string) string {
	koniec := 0
	for koniec < len(wiersz) && wiersz[koniec] >= '0' && wiersz[koniec] <= '9' {
		koniec++
	}
	if koniec == 0 || koniec >= len(wiersz) {
		return wiersz
	}
	if wiersz[koniec] == '.' || wiersz[koniec] == ')' {
		return wiersz[koniec+1:]
	}
	return wiersz
}
