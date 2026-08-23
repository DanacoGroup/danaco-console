// Odpowiedzialność pliku: głos jednego uczestnika debaty — wywołanie jego
// kanału modelu, strumień odpowiedzi i zapis wypowiedzi w turze.
//
// Strumień jedzie wspólną drogą. Kontrakt nie ma zdarzenia niosącego fragment
// wypowiedzi uczestnika, a `roundtable.debate.changed` niesie wypowiedź
// w całości. Fragmenty idą więc zdarzeniem `stream.chunk`, opisanym
// w kontrakcie jako jedna droga dla wszystkich kanałów: `windowId` wskazuje
// okno debaty, `messageId` — identyfikator wypowiedzi. Obie drogi (fragmenty
// treści i fragment domykający) niosą ten sam kod wypowiedzi; uzasadnienie
// wyboru stoi przy `zapytanieUczestnika`.
//
// Kolejność jest rozmyślna: wypowiedź pusta powstaje i rozgłasza się przed
// wywołaniem kanału. Bez tego klient dostawałby fragmenty opatrzone
// identyfikatorem, którego jeszcze nie zna, i nie miałby ich do czego przypiąć.
package core

import (
	"context"
	"strings"

	"danacoconsole/server/internal/dane"
	"danacoconsole/server/internal/models"
	"danacoconsole/shared"
)

// wypowiedz prowadzi jeden głos w turze i zwraca jego treść.
//
// Niepowodzenie kanału nie wpisuje się w treść wypowiedzi: przypisanie
// uczestnikowi słów, których nie powiedział, byłoby wytworzeniem zapisu.
// Przyczyna jedzie fragmentem błędu w strumieniu, a wypowiedź zostaje pusta —
// Model Panel pokazuje wtedy stan błędu tego jednego panelu.
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
	a.rozglos(shared.ChangeKindCreated, kontrakt, &wpis)

	strumien := nowyNadawcaStrumienia(a.nadajnik, wypowiedz.Kod, "")
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
		// Utrwalenie idzie po strumieniu i niezależnie od jego powodzenia: tekst,
		// który zdążył przyjść przed zerwaniem, jest zapisem tury tak samo jak
		// odpowiedź pełna.
		_ = a.repozytorium.UzupelnijWypowiedz(kontekst, wypowiedz.Kod, wypowiedz.Tresc)
	}
	po := wypowiedzKontraktu(wypowiedz)
	a.rozglos(shared.ChangeKindUpdated, kontrakt, &po)
	return wypowiedz.Tresc
}

// zapytanieUczestnika składa wywołanie kanału dla jednego uczestnika.
//
// `Wiadomosc` niesie kod wypowiedzi, nie kod uczestnika. Pole
// `models.Zapytanie.Wiadomosc` jest wprost polem `messageId` kontraktu
// (`zapytanie.go`, znacznik `json:"messageId"`) i zasila wszystkie cztery
// wytwórnie fragmentów (`models/fragment.go`: tekst, prowenancja, konto, błąd).
// Kod uczestnika w tym polu rozjechałby strumień: fragmenty treści szłyby
// z kodem uczestnika, a fragment domykający i błąd z kodem wypowiedzi, bo te
// buduje `nadawcaStrumienia.Zakoncz` z własnego identyfikatora.
//
// Wypowiedź wygrywa z uczestnikiem z trzech powodów, każdy sam wystarczający:
//
//  1. Uczestnik nie jest jednoznaczny w czasie — zabiera głos w każdej turze
//     debaty, więc jego kod wskazuje dowolną z wielu wypowiedzi. Strumień
//     opisuje jedno wywołanie kanału, więc klucz ma być jednorazowy.
//  2. `messageId` znaczy wiadomość. Wypowiedź jest wiadomością tury; uczestnik
//     jest jej autorem, a autor w polu identyfikatora wiadomości to inny byt.
//  3. Fragmentu domykającego nie da się przypisać uczestnikowi: błąd kanału
//     dotyczy tej jednej próby, nie osoby, a próba jest wypowiedzią.
//
// Drugiego pola nie dokładamy. Klient potrzebuje mówcy, ale ma go już bez
// pytania: rdzeń rozgłasza wypowiedź zdarzeniem `roundtable.debate.changed`
// jako `created` przed wywołaniem kanału, a `RoundtableStatement` niesie
// `participantId`. Dopisanie go do `stream.chunk` byłoby drugą drogą do wiedzy,
// którą klient już posiada.
//
// Tożsamość idzie warstwą nakładki, nie kodem kanału.
// `models.Nakladka.ProfilRoli` niesie prompt systemowy uczestnika wprost, bez
// ani jednego słowa dopisanego przez rdzeń. Uczestnik dodany bez promptu
// systemowego dostaje nakładkę pustą — dwaj tacy uczestnicy na jednym kanale
// odpowiedzą podobnie i jest to stan poprawny, bo Operator nie dał im różnych
// instrukcji.
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

// trescPytania buduje treść skierowaną do uczestnika: zagadnienie tury,
// wypowiedzi poprzedników i pytanie.
//
// Wiersz z nazwą tożsamości jest adresowaniem głosu, tak jak moderator udziela
// głosu przy stole, a nie instrukcją wymyśloną za Operatora: nazwę podał on sam
// przy dodaniu uczestnika, a wiersz widać w transkrypcie tury.
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
