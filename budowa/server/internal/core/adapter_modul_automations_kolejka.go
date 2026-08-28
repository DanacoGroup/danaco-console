// Odpowiedzialność pliku: okno Queue Manager — działanie silnika kolejek na
// kolejce automatyki wraz z zasileniem kolejki krokami definicji, w porządku
// topologicznym układu zależności, nie w kolejności ich wpisywania.
package core

import (
	"context"
	"strconv"

	"danacoconsole/server/internal/dane"
	"danacoconsole/shared"
)

// DzialanieKolejki wykonuje działanie silnika kolejek na kolejce automatyki
// i oddaje kolejkę po działaniu.
func (a *adapterAutomatyk) DzialanieKolejki(ctx context.Context,
	z shared.AutomationQueueActionRequest) (shared.AutomationQueueActionResponse, error) {

	if a.kolejki == nil {
		return shared.AutomationQueueActionResponse{}, errBrakSilnikaKolejek
	}
	id, err := strconv.ParseInt(z.QueueId, 10, 64)
	if err != nil {
		return shared.AutomationQueueActionResponse{}, bladWskazaniaAutomatyki(
			"kolejka o nieznanym identyfikatorze: " + z.QueueId)
	}
	automatyka, err := a.automatykaZadania(ctx, z.WorkflowId)
	if err != nil {
		return shared.AutomationQueueActionResponse{}, err
	}
	if automatyka != nil {
		if err := a.zasilKolejke(ctx, id, *automatyka); err != nil {
			return shared.AutomationQueueActionResponse{}, bladAutomatyki(err)
		}
		a.przypniOknaObserwatorow(id, automatyka.Kod)
	}
	if err := a.ulozPozycje(ctx, z); err != nil {
		return shared.AutomationQueueActionResponse{}, bladAutomatyki(err)
	}
	wynik, err := a.kolejki.Wykonaj(ctx, shared.QueueActionRequest{
		QueueId: z.QueueId, Action: z.Action, ItemId: z.ItemId,
	})
	if err != nil {
		return shared.AutomationQueueActionResponse{}, err
	}
	if err := a.odnotujPrzebieg(ctx, id, automatyka, wynik.Queue); err != nil {
		return shared.AutomationQueueActionResponse{}, bladAutomatyki(err)
	}
	return shared.AutomationQueueActionResponse{Queue: wynik.Queue}, nil
}

// rodzajKolejkiHarmonogramu znakuje kolejkę założoną przez budzik harmonogramu
// — odróżnia przebieg odpalony z zegara od kolejki sesyjnej i przekazania.
const rodzajKolejkiHarmonogramu = "harmonogram"

// UruchomAutomatyke zakłada kolejkę automatyki, zasila ją krokami w porządku
// układu zależności i rusza jej wykonanie — tą samą drogą, którą automatykę
// odpala Operator z Queue Managera, a budzik harmonogramu o wyliczonej godzinie.
func (a *adapterAutomatyk) UruchomAutomatyke(ctx context.Context,
	automatyka dane.Automatyka) (shared.Queue, error) {

	if a.kolejki == nil {
		return shared.Queue{}, errBrakSilnikaKolejek
	}
	id, err := a.kolejki.repozytorium.UtworzKolejke(ctx, dane.Kolejka{
		Nazwa: automatyka.Kod, Rodzaj: rodzajKolejkiHarmonogramu, Stan: shared.QueueStatusIdle,
	})
	if err != nil {
		return shared.Queue{}, bladAutomatyki(err)
	}
	if err := a.zasilKolejke(ctx, id, automatyka); err != nil {
		return shared.Queue{}, bladAutomatyki(err)
	}
	a.przypniOknaObserwatorow(id, automatyka.Kod)
	wynik, err := a.kolejki.Wykonaj(ctx, shared.QueueActionRequest{
		QueueId: strconv.FormatInt(id, 10), Action: shared.QueueActionStart,
	})
	if err != nil {
		return shared.Queue{}, err
	}
	if err := a.odnotujPrzebieg(ctx, id, &automatyka, wynik.Queue); err != nil {
		return shared.Queue{}, bladAutomatyki(err)
	}
	return wynik.Queue, nil
}

// automatykaZadania dobiera automatykę wskazaną w żądaniu. Żądanie bez
// wskazania jest zwykłym działaniem na kolejce — nie każda kolejka wykonuje
// automatykę, bo silnik jest wspólny z pętlą sesyjną.
func (a *adapterAutomatyk) automatykaZadania(ctx context.Context, kod *string) (*dane.Automatyka, error) {
	if wartoscTekstu(kod) == "" {
		return nil, nil
	}
	wiersz, err := a.wiersz(ctx, *kod)
	if err != nil {
		return nil, err
	}
	return &wiersz, nil
}

// zasilKolejke zamienia kroki automatyki na zlecenia kolejki. Kolejka, która
// ma już pozycje, zostaje nietknięta: powtórne uruchomienie tej samej kolejki
// jest biegiem naprawczym po zastanych zleceniach, nie ich podwojeniem.
func (a *adapterAutomatyk) zasilKolejke(ctx context.Context, kolejkaID int64,
	automatyka dane.Automatyka) error {

	zastane, err := a.kolejki.repozytorium.ListaPozycji(ctx, kolejkaID)
	if err != nil {
		return err
	}
	if len(zastane) > 0 {
		return nil
	}
	kroki, err := a.repozytorium.Kroki(ctx, automatyka.ID)
	if err != nil {
		return err
	}
	zaleznosci, err := a.repozytorium.Zaleznosci(ctx, automatyka.ID)
	if err != nil {
		return err
	}
	for _, krok := range ulozoneKroki(kroki, zaleznosci) {
		if _, err := a.kolejki.repozytorium.DodajPozycje(ctx, pozycjaZKroku(kolejkaID, krok)); err != nil {
			return err
		}
	}
	return nil
}

// ulozoneKroki podaje kroki w kolejności wykonania: najpierw porządek
// topologiczny układu zależności, a kroki uwikłane w cykl na końcu.
func ulozoneKroki(kroki []dane.KrokAutomatyki, zaleznosci []dane.ZaleznoscKroku) []dane.KrokAutomatyki {
	kolejnosc, pozostale := porzadekTopologiczny(kroki, zaleznosci)
	wedlugKodu := make(map[string]dane.KrokAutomatyki, len(kroki))
	for _, krok := range kroki {
		wedlugKodu[krok.Kod] = krok
	}
	ulozone := make([]dane.KrokAutomatyki, 0, len(kroki))
	for _, kod := range append(kolejnosc, pozostale...) {
		if krok, jest := wedlugKodu[kod]; jest {
			ulozone = append(ulozone, krok)
		}
	}
	return ulozone
}

// pozycjaZKroku składa zlecenie kolejki z kroku automatyki. Treścią zlecenia
// jest komenda kroku wraz z jej ładunkiem — to ona ma zostać wykonana.
func pozycjaZKroku(kolejkaID int64, krok dane.KrokAutomatyki) dane.Pozycja {
	tytul := krok.Kod
	if krok.Nazwa != nil && *krok.Nazwa != "" {
		tytul = *krok.Nazwa
	}
	pozycja := dane.Pozycja{KolejkaID: kolejkaID, Tytul: tytul, Stan: stanPozycjiOczekuje}
	if tresc := trescZlecenia(krok); tresc != "" {
		pozycja.TrescZlecenia = &tresc
	}
	return pozycja
}

// trescZlecenia składa polecenie kroku: komendę, ładunek i warunek wykonania.
// Krok bez komendy niesie sam warunek albo nic — i tak też zostaje zapisany,
// zamiast dostać treść wymyśloną przez rdzeń.
func trescZlecenia(krok dane.KrokAutomatyki) string {
	czesci := []string{}
	if krok.Komenda != nil && *krok.Komenda != "" {
		czesci = append(czesci, *krok.Komenda)
	}
	if krok.Parametry != nil && *krok.Parametry != "" {
		czesci = append(czesci, *krok.Parametry)
	}
	if krok.Warunek != nil && *krok.Warunek != "" {
		czesci = append(czesci, "warunek: "+*krok.Warunek)
	}
	tresc := ""
	for numer, czesc := range czesci {
		if numer > 0 {
			tresc += "\n"
		}
		tresc += czesc
	}
	return tresc
}
