// Plik obsługuje głos jednego uczestnika debaty: wywołanie jego kanału modelu, strumień odpowiedzi i zapis wypowiedzi w turze. Fragmenty idą zdarzeniem `stream.chunk`, wspólną drogą dla wszystkich kanałów.
package core

import (
	"context"
	"strings"

	"danacoconsole/server/internal/dane"
	"danacoconsole/server/internal/models"
	"danacoconsole/shared"
)

// wypowiedz prowadzi jeden głos w turze i zwraca jego treść. Niepowodzenie kanału nie wpisuje się w treść wypowiedzi; przyczyna jedzie fragmentem błędu w strumieniu, a wypowiedź zostaje pusta.
func (a *adapterDebaty) wypowiedz(kontekst context.Context, tura dane.TuraDebaty,
	uczestnik dane.UczestnikDebaty, pytanie, tlo string) string {

	wypowiedz, err := a.repozytorium.ZapiszWypowiedz(kontekst, dane.WypowiedzDebaty{
		Kod: nowyIdentyfikator(przedrostekWypowiedzi), TuraKod: tura.Kod,
		Uczestnik: uczestnik.Kod,
	})
	if err != nil {
		return ""
	}
	kontrakt := turaKontraktu(tura)
	wpis := wypowiedzKontraktu(wypowiedz)
	a.rozglos(kontekst, shared.ChangeKindCreated, kontrakt, &wpis)

	strumien := nowyNadawcaStrumienia(a.nadajnik, kontoAdresata(kontekst), wypowiedz.Kod, "")
	var tresc strings.Builder
	ujscie := models.UjscieFunkcji(func(ctx context.Context, f models.Fragment) error {
		if f.Kind == shared.ChunkKindText {
			tresc.WriteString(models.TrescFragmentu(f))
		}
		return strumien.Fragment(ctx, f)
	})

	zapytanie := zapytanieUczestnika(tura, uczestnik, wypowiedz.Kod, pytanie, tlo)
	blad := a.kanaly.Wyslij(kontekst, zapytanie, ujscie)
	strumien.Zakoncz(tura.Okno, wypowiedz.Kod, tresc.String(), blad)

	wypowiedz.Tresc = tresc.String()
	if wypowiedz.Tresc != "" {
		// Utrwalenie idzie po strumieniu niezależnie od powodzenia: tekst przed zerwaniem to zapis tury.
		_ = a.repozytorium.UzupelnijWypowiedz(kontekst, wypowiedz.Kod, wypowiedz.Tresc)
	}
	po := wypowiedzKontraktu(wypowiedz)
	a.rozglos(kontekst, shared.ChangeKindUpdated, kontrakt, &po)
	return wypowiedz.Tresc
}

// zapytanieUczestnika składa wywołanie kanału dla jednego uczestnika. `Wiadomosc` niesie kod wypowiedzi, nie kod uczestnika, bo wypowiedź jest kluczem jednorazowym strumienia, a uczestnik może zabierać głos wielokrotnie.
func zapytanieUczestnika(tura dane.TuraDebaty, uczestnik dane.UczestnikDebaty,
	kodWypowiedzi, pytanie, tlo string) models.Zapytanie {

	return models.Zapytanie{
		Zasiegi:   models.Zasiegi{Okno: tura.Okno},
		Wiadomosc: kodWypowiedzi,
		Tresc:     trescPytania(tura, uczestnik, pytanie, tlo),
		Kanal:     uczestnik.KanalModelu,
		Nakladka:  models.Nakladka{ProfilRoli: wartoscTekstu(uczestnik.PromptSystemowy)},
	}
}

// trescPytania buduje treść skierowaną do uczestnika: zagadnienie tury, wypowiedzi poprzedników i pytanie. Wiersz z nazwą tożsamości jest adresowaniem głosu, tak jak moderator udziela głosu przy stole.
func trescPytania(tura dane.TuraDebaty, uczestnik dane.UczestnikDebaty,
	pytanie, tlo string) string {

	czesci := make([]string, 0, 4)
	if nazwa := strings.TrimSpace(wartoscTekstu(uczestnik.NazwaTozsamosci)); nazwa != "" {
		czesci = append(czesci, "Głos w debacie zabiera: "+nazwa+".")
	}
	if tura.Zagadnienie != nil && strings.TrimSpace(*tura.Zagadnienie) != "" {
		czesci = append(czesci, "Zagadnienie tury: "+strings.TrimSpace(*tura.Zagadnienie)+".")
	}
	if strings.TrimSpace(tlo) != "" {
		czesci = append(czesci, "Dotychczasowe wypowiedzi w tej turze:\n"+strings.TrimSpace(tlo))
	}
	czesci = append(czesci, pytanie)
	return strings.Join(czesci, "\n\n")
}

// nazwaUczestnika zwraca podpis uczestnika w transkrypcie: nazwę tożsamości,
// a w jej braku kod kanału z dopiskiem identyfikatora — dwaj uczestnicy na tym
// samym kanale muszą być rozróżnialni także wtedy, gdy nikt ich nie nazwał.
func nazwaUczestnika(uczestnik dane.UczestnikDebaty) string {
	if nazwa := strings.TrimSpace(wartoscTekstu(uczestnik.NazwaTozsamosci)); nazwa != "" {
		return nazwa
	}
	return uczestnik.KanalModelu + " (" + uczestnik.Kod + ")"
}
