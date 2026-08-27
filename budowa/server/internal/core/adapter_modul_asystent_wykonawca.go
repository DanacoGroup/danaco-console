// Odpowiedzialność pliku: wykonawca zleceń asystenta — podejmuje zlecenie ze
// stanu `queued` i prowadzi je do `done`/`failed`. Tura idzie rejestrem
// kanałów wspólnym z rozmową; wykonawca prowadzi wyłącznie drogę modelu, nie
// zgaduje komend z tekstu.
package core

import (
	"context"
	"strings"
	"time"

	"danacoconsole/server/internal/dane"
	"danacoconsole/server/internal/models"
	"danacoconsole/server/internal/session"
	"danacoconsole/shared"
)

// ZKanalami wpina rejestr kanałów modelu. Bez niego moduł przyjmuje polecenia
// i prowadzi ich stan ręcznie, lecz zlecenia z kolejki nie podejmuje — tura nie
// ma dokąd pójść, więc `queued` zostaje `queued`.
func (a *adapterAsystenta) ZKanalami(kanaly *models.Rejestr) *adapterAsystenta {
	a.kanaly = kanaly
	return a
}

// ZSesjami wpina nadzorcę sesji. Wykonawca bierze z niego okno zlecenia — kanał
// modelu, katalogi robocze i środowisko wykonania. Bez niego nie ma z czego
// złożyć tury, więc zlecenie zostaje w kolejce.
func (a *adapterAsystenta) ZSesjami(nadzorca *session.Nadzorca) *adapterAsystenta {
	a.nadzorca = nadzorca
	return a
}

// ZWyjsciem wpina nadajnik zdarzeń i kontekst życia rdzenia. Nadajnik niesie
// strumień tury i zdarzenie `assistant.action.changed`; kontekst życia trzyma
// turę przy życiu mimo rozłączenia klienta. Bez nadajnika okno nie dostaje
// strumienia na żywo.
func (a *adapterAsystenta) ZWyjsciem(nadajnik Nadajnik, zycie context.Context) *adapterAsystenta {
	a.nadajnik = nadajnik
	a.zycie = zycie
	return a
}

// podejmij uruchamia wykonawcę zlecenia świeżo złożonego albo ponowionego,
// stojącego w stanie `queued`. Brak rejestru kanałów, nadzorcy sesji albo
// kontekstu życia znaczy brak wykonawcy — zlecenie zostaje `queued`. Bieg
// idzie gorutyną.
func (a *adapterAsystenta) podejmij(kodZlecenia string) {
	a.podejmijZeStanu(kodZlecenia, string(shared.AssistantActionStatusQueued))
}

// wznow podejmuje zlecenie wznowione sterowaniem `resume`. Sterowanie
// przestawia zlecenie na `running`, a wykonawca podejmuje stan `queued`,
// więc bez tego wejścia zlecenie stałoby bez końca. Idzie tą samą drogą co
// podjęcie, ze stanem wejścia `running`.
func (a *adapterAsystenta) wznow(kodZlecenia string) {
	a.podejmijZeStanu(kodZlecenia, string(shared.AssistantActionStatusRunning))
}

// przerwijBieg przerywa turę zlecenia, jeśli jakaś biegnie. Zapisany anulator
// należy do kontekstu tury (`wykonajZlecenie`); jego brak znaczy „nic nie
// biegnie" i nie jest błędem — Operator mógł sterować zleceniem stojącym.
func (a *adapterAsystenta) przerwijBieg(kodZlecenia string) {
	a.zamki.Lock()
	if a.odwolane == nil {
		a.odwolane = map[string]bool{}
	}
	// Znacznik zostaje, gdy nic jeszcze nie biegnie: zapis stanu przez
	// sterowanie mógłby go nadpisać.
	a.odwolane[kodZlecenia] = true
	przerwij := a.biegi[kodZlecenia]
	a.zamki.Unlock()
	if przerwij != nil {
		przerwij()
	}
}

// odwolany mówi, czy Operator zdążył odwołać to zlecenie; zapomnijOdwolanie
// zdejmuje znacznik, gdy zlecenie dostaje świeży mandat (podjęcie, ponowna
// próba, wznowienie) — dawne odwołanie nie unieważnia nowej decyzji.
func (a *adapterAsystenta) odwolany(kodZlecenia string) bool {
	a.zamki.Lock()
	defer a.zamki.Unlock()
	return a.odwolane[kodZlecenia]
}

func (a *adapterAsystenta) zapomnijOdwolanie(kodZlecenia string) {
	a.zamki.Lock()
	defer a.zamki.Unlock()
	delete(a.odwolane, kodZlecenia)
}

// zapiszBieg i zdejmijBieg prowadzą wykaz tur w locie. Wykaz jest jedyną drogą,
// którą sterowanie Operatora (`cancel`, `pause`) sięga do pracy już trwającej —
// bez niego odwołanie zmieniałoby sam wiersz, a tura pracowałaby dalej.
func (a *adapterAsystenta) zapiszBieg(kodZlecenia string, przerwij context.CancelFunc) {
	a.zamki.Lock()
	defer a.zamki.Unlock()
	if a.biegi == nil {
		a.biegi = map[string]context.CancelFunc{}
	}
	a.biegi[kodZlecenia] = przerwij
}

func (a *adapterAsystenta) zdejmijBieg(kodZlecenia string) {
	a.zamki.Lock()
	defer a.zamki.Unlock()
	delete(a.biegi, kodZlecenia)
}

// podejmijZeStanu jest jedynym wejściem wykonawcy. `stanWejscia` mówi, w jakim
// stanie zlecenie ma stać, żeby wolno je było podjąć — to chroni przed
// podwójnym biegiem tego samego zlecenia dwiema drogami naraz.
func (a *adapterAsystenta) podejmijZeStanu(kodZlecenia, stanWejscia string) {
	if a.kanaly == nil || a.nadzorca == nil || a.zycie == nil || kodZlecenia == "" {
		return
	}
	// Świeży mandat kasuje dawne odwołanie: `retry` i `resume` są decyzją
	// późniejszą niż `cancel`.
	a.zapomnijOdwolanie(kodZlecenia)
	// Sprawcą biegu jest rdzeń: zlecenie leży w kolejce, a tutaj zmienia je
	// własny wykonawca, nie gniazdo.
	go a.wykonajZlecenie(zSprawcaRdzenia(a.zycie), kodZlecenia, stanWejscia)
}

// wykonajZlecenie podejmuje jedno zlecenie z kolejki i prowadzi je do końca.
// Kolejność jest rozmyślna: stan, treść polecenia, okno, dopiero potem
// `running`. Brak treści albo kanału zamyka zlecenie stanem `failed` zamiast
// wiersza w `queued` bez powodu.
func (a *adapterAsystenta) wykonajZlecenie(ctx context.Context, kodZlecenia, stanWejscia string) {
	zlecenie, err := a.repozytorium.Zlecenie(ctx, kodZlecenia)
	if err != nil || zlecenie.Stan != stanWejscia {
		return
	}

	// Okno idzie przed treścią: z okna bierze się identyfikator sesji
	// rozgłaszający każdą zmianę.
	okno, err := a.nadzorca.Rejestr().Okno(zlecenie.OknoKod)
	if err != nil {
		a.zerwijZlecenie(ctx, zlecenie, kodZlecenia, "",
			"okno zlecenia "+zlecenie.OknoKod+" jest nieczytelne dla rdzenia: "+err.Error()+
				"; naprawa: złożyć polecenie w istniejącym oknie asystenta")
		return
	}
	if strings.TrimSpace(okno.KanalModelu) == "" {
		// Okno bez kanału modelu nie ma czym poprowadzić tury — brak do
		// naprawienia przez Operatora.
		a.zerwijZlecenie(ctx, zlecenie, kodZlecenia, okno.IdSesji,
			"okno "+zlecenie.OknoKod+" nie ma wpiętego kanału modelu, więc tura nie ma dokąd pojechać"+
				"; naprawa: wskazać oknu kanał modelu (channel.add + window.update) i ponowić zlecenie")
		return
	}

	polecenie := a.trescPolecenia(ctx, kodZlecenia)
	if polecenie == "" {
		// Zlecenie ze śladem nagrania, którego nikt nie przepisał na tekst, nie
		// niesie polecenia do wykonania.
		a.zerwijZlecenie(ctx, zlecenie, kodZlecenia, okno.IdSesji,
			"zlecenie nie niesie treści polecenia — dziennik ma sam ślad nagrania, a rdzeń nie ma"+
				" czym go przepisać; naprawa: przysłać polecenie tekstem (pole transcript)")
		return
	}

	// Ostatnie spojrzenie przed przejęciem: Operator mógł zdążyć odwołać
	// zlecenie przed zapisem `running`.
	if a.odwolany(kodZlecenia) {
		return
	}
	if biezace, err := a.repozytorium.UstawStanZlecenia(ctx, kodZlecenia,
		string(shared.AssistantActionStatusRunning)); err == nil {
		a.rozglosZlecenie(ctx, shared.ChangeKindUpdated, okno.IdSesji, zlozZlecenie(biezace))
	}

	// Tura dostaje kontekst, by sterowanie miało co przerwać; rozłączenie
	// klienta nie przerywa pracy.
	kontekstTury, przerwij := context.WithCancel(ctx)
	a.zapiszBieg(kodZlecenia, przerwij)
	tresc, bladTury := a.turaModelu(kontekstTury, okno, zlecenie, polecenie)
	a.zdejmijBieg(kodZlecenia)
	przerwij()

	stan := string(shared.AssistantActionStatusDone)
	wynik := tresc
	if bladTury != nil {
		stan = string(shared.AssistantActionStatusFailed)
		wynik = "tura modelu nie powiodła się: " + bladTury.Error()
	}
	a.domknijZlecenie(ctx, zlecenie, kodZlecenia, okno.IdSesji, stan, wynik, tresc)
}

// zerwijZlecenie zamyka zlecenie stanem `failed` z nazwanym powodem. Idzie tą
// samą drogą co domknięcie po turze — różni je tylko to, że tura nie ruszyła,
// więc nie ma treści do zapisania w dzienniku.
func (a *adapterAsystenta) zerwijZlecenie(ctx context.Context, zlecenie dane.ZlecenieAsystenta,
	kodZlecenia, idSesji, powod string) {

	po, err := a.repozytorium.ZakonczZlecenie(ctx, kodZlecenia,
		string(shared.AssistantActionStatusFailed), "zlecenia nie da się wykonać: "+powod)
	if err != nil {
		return
	}
	a.rozglosZlecenie(ctx, shared.ChangeKindUpdated, idSesji, zlozZlecenie(po))
}

// trescPolecenia wyjmuje tekst polecenia z dziennika zlecenia — wpis rodzaju
// `command` niesie transkrypcję zapisaną przez PrzyjmijPolecenie. Bierze
// wpis najstarszy z niepustą treścią: to pierwotne polecenie, nie wynik
// dopisany później.
func (a *adapterAsystenta) trescPolecenia(ctx context.Context, kodZlecenia string) string {
	wpisy, err := a.repozytorium.WpisyZlecenia(ctx, kodZlecenia)
	if err != nil {
		return ""
	}
	// WpisyZlecenia oddaje wpisy od najnowszego — polecenie jest najstarsze,
	// więc idziemy od końca listy.
	for i := len(wpisy) - 1; i >= 0; i-- {
		if wpisy[i].Rodzaj == string(shared.AssistantActivityKindCommand) {
			if tresc := strings.TrimSpace(wpisy[i].Tresc); tresc != "" {
				return tresc
			}
		}
	}
	return ""
}

// turaModelu prowadzi jedną turę kanału modelu dla zlecenia i zwraca jej
// treść. Strumień idzie tą samą drogą co w rozmowie: `stream.chunk`
// z identyfikatorem zlecenia jako identyfikatorem wiadomości. Niepowodzenie
// kanału wraca w błędzie.
func (a *adapterAsystenta) turaModelu(ctx context.Context, okno session.Okno,
	zlecenie dane.ZlecenieAsystenta, polecenie string) (string, error) {

	strumien := nowyNadawcaStrumienia(a.nadajnik, zlecenie.Kod, okno.IdSesji)
	var tresc strings.Builder
	ujscie := models.UjscieFunkcji(func(c context.Context, f models.Fragment) error {
		if f.Kind == shared.ChunkKindText {
			tresc.WriteString(models.TrescFragmentu(f))
		}
		return strumien.Fragment(c, f)
	})

	// Sterowanie platformą dojeżdża do modelu tędy: bez tej linii model nie
	// widzi narzędzi kontraktu.
	zapytanie := zapytanieZlecenia(okno, zlecenie, polecenie)

	// Warstwa promptu profilu wchodzi tylko tutaj: po złożeniu zapytania,
	// przed narzędziami i wysyłką.
	profil, jestProfil := a.profilZlecenia(ctx, zlecenie)
	if jestProfil {
		uzupelnijProfil(profil, &zapytanie)
	}
	// Brak warstwy zostaje nazwany: z dziennika widać, dlaczego asystent
	// zachował się jak czat.
	switch {
	case !jestProfil:
		a.opiszBrakWarstwy(ctx, zlecenie, "żadnego profilu nie wskazano i nie ma profilu domyślnego")
	case warstwaProfilu(profil) == "":
		a.opiszBrakWarstwy(ctx, zlecenie, "profil "+profil.Nazwa+" nie ma zapisanej warstwy promptu")
	}

	a.uzupelnijNarzedzia(ctx, okno, &zapytanie)

	blad := a.kanaly.Wyslij(ctx, zapytanie, ujscie)
	strumien.Zakoncz(okno.Id, zlecenie.Kod, tresc.String(), blad)
	return tresc.String(), blad
}

// zapytanieZlecenia składa wywołanie kanału z parametrów okna asystenta. Kanał
// nie sięga po konfigurację sam — parametry przychodzą z okna, tak jak
// w `zapytanieKanalu` rozmowy.
func zapytanieZlecenia(okno session.Okno, zlecenie dane.ZlecenieAsystenta, polecenie string) models.Zapytanie {
	return models.Zapytanie{
		Zasiegi:             models.Zasiegi{Sesja: okno.IdSesji, Okno: okno.Id},
		Wiadomosc:           zlecenie.Kod,
		Tresc:               polecenie,
		Kanal:               okno.KanalModelu,
		KatalogiRobocze:     okno.KatalogiRobocze,
		SrodowiskoWykonania: okno.SrodowiskoWykonania,
		TrybUprawnien:       okno.TrybUprawnien,
		RolaOkna:            okno.RolaOkna,
	}
}

// domknijZlecenie zapisuje odpowiedź modelu w dzienniku, przestawia zlecenie
// na stan końcowy z wynikiem i rozgłasza zmianę. Wpis dziennika i domknięcie
// stanu idą w tej kolejności: treść przybyłą przed zerwaniem zapisuje się
// zawsze, także przy `failed`.
func (a *adapterAsystenta) domknijZlecenie(ctx context.Context, zlecenie dane.ZlecenieAsystenta,
	kodZlecenia, idSesji, stan, wynik, tresc string) {

	if strings.TrimSpace(tresc) != "" {
		_, _ = a.repozytorium.ZapiszWpis(ctx, dane.WpisDziennikaAsystenta{
			Kod:         nowyIdentyfikator(przedrostekWpisuAsystenta),
			OknoKod:     zlecenie.OknoKod,
			ZlecenieKod: &kodZlecenia,
			Rodzaj:      string(shared.AssistantActivityKindResult),
			Tresc:       tresc,
			Utworzono:   time.Now().UnixMilli(),
		})
	}

	// Decyzja Operatora wygrywa z wynikiem: sprawdzenie stanu chroni przed
	// zameldowaniem odwołanej pracy.
	if biezace, err := a.repozytorium.Zlecenie(ctx, kodZlecenia); err == nil &&
		biezace.Stan != string(shared.AssistantActionStatusRunning) {

		a.rozglosZlecenie(ctx, shared.ChangeKindUpdated, idSesji, zlozZlecenie(biezace))
		return
	}

	po, err := a.repozytorium.ZakonczZlecenie(ctx, kodZlecenia, stan, wynik)
	if err != nil {
		return
	}
	a.rozglosZlecenie(ctx, shared.ChangeKindUpdated, idSesji, zlozZlecenie(po))
}

// rozglosZlecenie rozgłasza zmianę stanu zlecenia zdarzeniem
// `assistant.action.changed`. Nadajnik niepodłączony nie jest błędem: wykonawca
// domyka stan także wtedy, gdy nikt nie słucha zdarzeń.
func (a *adapterAsystenta) rozglosZlecenie(ctx context.Context, zmiana shared.ChangeKind,
	idSesji string, z shared.AssistantAction) {

	if a == nil || a.nadajnik == nil {
		return
	}
	nowyEmiter(a.nadajnik).zlecenieAsystenta(ctx, zmiana, idSesji, z)
}
