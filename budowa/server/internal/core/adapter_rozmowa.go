package core

import (
	"context"
	"encoding/json"
	"runtime/debug"
	"strings"
	"sync"
	"time"

	"danacoconsole/server/internal/models"
	"danacoconsole/server/internal/protocol"
	"danacoconsole/server/internal/session"
	"danacoconsole/shared"
)

// adapterRozmowy wypełnia port Rozmowa: przyjmuje wiadomość okna, kieruje ją do
// kanału modelu z rejestru i odsyła odpowiedź strumieniem.
//
// Odpowiedź komendy potwierdza przyjęcie wiadomości od razu, a treść modelu idzie
// osobno — tura trwa dłużej niż wykonanie komendy, a klient nie może czekać na
// potwierdzenie do końca odpowiedzi.
type adapterRozmowy struct {
	nadzorca *session.Nadzorca
	kanaly   *models.Rejestr
	nadajnik Nadajnik
	dziennik *dziennikRozmowy
	zycie    context.Context
	// petla prowadzi bieg koordynator–wykonawca. Pusta pętla nie blokuje
	// rozmowy — okno samodzielne pracuje bez niej.
	petla *session.Petla
	// tozsamosc podaje nakładkę obowiązującą okna.
	tozsamosc Tozsamosc
	// agenci podają tożsamość eksperta wskazanego przez okno: warstwy promptu,
	// model bazowy i nastawy procesu. Port pusty znaczy okno na modelu surowym.
	agenci ZrodloTozsamosciAgenta
	// dolozeniaSesji podaje doraźne dołożenia — drugie źródło zestawu narzędzi
	// tury (adapter_rozmowa_zestaw.go). Pierwszym jest ekspert okna, wskazany
	// portem `agenci` wyżej. Port pusty nie odbiera modelowi narzędzi: tura
	// jedzie samą podstawą eksperta i melduje powód do dziennika rdzenia.
	dolozeniaSesji DolozeniaNarzedziSesji
	// zgloszoneZestawy pilnuje, by powód zestawu szedł do dziennika raz na
	// powód, a nie raz na turę.
	zgloszoneZestawy sync.Map
	// katalog ustala katalog roboczy sesji, mosty — konfigurację MCP ze zbioru
	// nadań okna.
	katalog *KatalogRoboczy
	mosty   *mostyOkna
	// ciaglosc utrwala identyfikator rozmowy programu CLI przy oknie, dzięki
	// czemu kolejna tura wznawia rozmowę zamiast zaczynać od zera.
	// Port pusty nie blokuje rozmowy — znika tylko ciągłość.
	ciaglosc CiagloscRozmowy
	// wykonanie podaje nakład rozumowania i model zapasowy z konfiguracji.
	// Port pusty zostawia decyzję kanałowi.
	wykonanie ParametryWykonania
	// konfiguracja odczytuje obowiązującą konfigurację sesji i tłumaczy jej
	// obszary na wejście procesu (adapter_rozmowa_konfiguracja.go). Port pusty
	// zostawia turę na wartościach okna i wiersza rejestru.
	konfiguracja CzytelnikKonfiguracjiSesji
	// zdarzenia utrwala zamknięcia tur. Port pusty nie blokuje rozmowy — znika
	// wyłącznie dowód rozstrzygnięcia.
	zdarzenia *zdarzeniaWykonawcze
	// bloki utrwala nietekstowe fragmenty strumienia w trakcie tury:
	// rozumowanie, narzędzia, prowenancję. Port pusty nie blokuje rozmowy —
	// historia wraca wtedy samym tekstem.
	bloki *rejestratorBlokow
	// zalaczniki jest magazynem bajtów załącznika wiadomości
	// (adapter_rozmowa_zalaczniki.go). Magazyn pusty nie blokuje tury, ale
	// załącznik niosący bajty zostaje wtedy niedoręczony i tura mówi to
	// modelowi wprost, zamiast milczeć.
	zalaczniki *magazynTresciBiblioteki

	mu       sync.Mutex
	biegnace map[string]context.CancelFunc
}

// nowyAdapterRozmowy wiąże port z pakietem sesji i rejestrem kanałów. Kontekst
// życia jest kontekstem rdzenia, nie połączenia: rozłączenie klienta nie
// przerywa rozpoczętej tury.
func nowyAdapterRozmowy(zycie context.Context, nadzorca *session.Nadzorca,
	kanaly *models.Rejestr, nadajnik Nadajnik, dziennik *dziennikRozmowy) *adapterRozmowy {
	return &adapterRozmowy{
		nadzorca: nadzorca, kanaly: kanaly, nadajnik: nadajnik, dziennik: dziennik,
		zycie: zycie, biegnace: map[string]context.CancelFunc{},
	}
}

// UstawPetle wpina pętlę koordynator–wykonawca po jej złożeniu. Wiązanie jest
// dwuetapowe, bo adapter i pętla znają się nawzajem: pętla potrzebuje portu
// rozpoczynania obiegu, którym jest ten adapter.
func (a *adapterRozmowy) UstawPetle(p *session.Petla) { a.petla = p }

// ZCiagloscia wpina utrwalanie identyfikatora rozmowy programu CLI.
// Adapter bez tego portu prowadzi rozmowę bez ciągłości — każda tura zaczyna od
// nowa. Zwraca siebie, żeby dało się złożyć w jednym wierszu montażu.
func (a *adapterRozmowy) ZCiagloscia(c CiagloscRozmowy) *adapterRozmowy {
	a.ciaglosc = c
	return a
}

// ZZdarzeniami wpina odbiornik zdarzeń wykonawczych: zamknięcie tury
// jedzie do tabeli zamknięć razem ze stanem, który z niego wyprowadzono.
func (a *adapterRozmowy) ZZdarzeniami(z *zdarzeniaWykonawcze) *adapterRozmowy {
	a.zdarzenia = z
	return a
}

// ZBlokami wpina rejestrator bloków wiadomości: fragmenty nietekstowe
// strumienia jadą do bazy w trakcie tury, zamiast ginąć z procesem.
func (a *adapterRozmowy) ZBlokami(b *rejestratorBlokow) *adapterRozmowy {
	a.bloki = b
	return a
}

// ZZalacznikami wpina magazyn załączników rozmowy nad katalogiem danych
// wskazanym konfiguracją procesu — tą samą wartością, którą montaż podaje
// sejfowi poświadczeń i magazynowi biblioteki (jedna prawda o katalogu).
// Adapter bez tego portu prowadzi rozmowę, ale załącznik niosący bajty jedzie
// do modelu jako niedoręczony, nie jako cisza.
func (a *adapterRozmowy) ZZalacznikami(katalogDanych string) *adapterRozmowy {
	a.zalaczniki = magazynZalacznikow(katalogDanych)
	return a
}

// Wyslij przyjmuje wiadomość użytkownika i otwiera turę okna.
func (a *adapterRozmowy) Wyslij(ctx context.Context, z shared.MessageSendRequest) (shared.MessageSendResponse, error) {
	okno, err := a.nadzorca.Rejestr().Okno(z.WindowId)
	if err != nil {
		return shared.MessageSendResponse{}, bladSesji(err)
	}
	// Zajęcie okna idzie przed dziennikiem. Okno prowadzące turę odmawia
	// (adapter_rozmowa_przerwanie.go), a odmowa nie może zostawić po sobie
	// dwóch wierszy w dzienniku rozmowy — pytania, które nie pojechało, i
	// odpowiedzi, która nigdy nie ruszy. Odwołanie tury powstaje więc tutaj,
	// a nie po zapisie: wcześniej nie ma czego zajmować.
	kontekst, anuluj := context.WithCancel(a.zycie)
	if !a.zajmijBieg(okno.Id, anuluj) {
		anuluj()
		return shared.MessageSendResponse{}, odmowaTuryWBiegu(okno.Id)
	}

	// Atrybucja tury: obie wiadomości powstają w tym samym oknie i w tej samej
	// roli pętli, więc niosą tę samą metrykę — inaczej pytanie i odpowiedź jednego
	// obiegu rozjechałyby się w kolumnach `persona` i `okno_zrodlowe_id`.
	atrybucja := atrybucjaOkna(okno)
	pytanie := shared.Message{
		Id: identyfikatorWiadomosci(), WindowId: okno.Id, SessionId: okno.IdSesji,
		Role: shared.MessageRoleUser, Content: z.Content, Status: shared.MessageStatusComplete,
		Attachments: z.Attachments, CreatedAt: time.Now().UnixMilli(),
		Metadata: atrybucja,
	}
	odpowiedz := shared.Message{
		Id: identyfikatorWiadomosci(), WindowId: okno.Id, SessionId: okno.IdSesji,
		Role: shared.MessageRoleAssistant, Status: shared.MessageStatusStreaming,
		CreatedAt: time.Now().UnixMilli(),
		Metadata:  atrybucja,
	}
	a.dziennik.Dopisz(pytanie)
	a.dziennik.Dopisz(odpowiedz)

	// Tożsamość strumienia zdejmuje się tutaj, a nie w turze: kontekst tury
	// wisi na życiu adaptera (`a.zycie`), więc wpisu żądania już nie niesie.
	// Zastępczo — dla tury powołanej poza drogą komendy — idzie identyfikator
	// pytania.
	idZadania := protocol.TozsamoscStrumienia(ctx, pytanie.Id)

	go a.prowadzTure(kontekst, okno, pytanie, odpowiedz, idZadania)

	return shared.MessageSendResponse{Message: pytanie}, nil
}

// Wykaz zwraca wiadomości okna.
func (a *adapterRozmowy) Wykaz(_ context.Context, z shared.MessageListRequest) (shared.MessageListResponse, error) {
	wiadomosci, sastarsze := a.dziennik.Wykaz(z.WindowId, z.Before, z.Limit)
	return shared.MessageListResponse{Messages: wiadomosci, HasMore: sastarsze}, nil
}

// prowadzTure wykonuje jedną turę kanału modelu i odsyła jej strumień.
func (a *adapterRozmowy) prowadzTure(kontekst context.Context, okno session.Okno,
	pytanie, odpowiedz shared.Message, idZadania string) {

	// Osłona gorutyny tury. Rejestrowana PIERWSZA, więc przy panice biegnie
	// OSTATNIA — dopiero po tym, jak strumień domknie się awaryjnie, a bieg
	// zostanie zapomniany. Klient dostaje więc znacznik końca i wpis rozmowy
	// wychodzi ze stanu `strumien`, a dopiero potem panika przestaje lecieć.
	//
	// Bez tej osłony KAŻDA usterka adaptera w torze tury gasi cały rdzeń, a nie
	// jedną turę: panika w gorutynie nie ma kto przechwycić i proces ginie
	// z sesją, kolejką i połączeniem Operatora naraz. Tak wywracał go pusty
	// wskaźnik portu doraźnych dołożeń — jedna wada jednego adaptera zabierała
	// całą pracę.
	//
	// Osłona NIE jest zgodą na usterki i niczego nie ucisza: powód idzie do
	// dziennika ze śladem stosu, żeby wada została zgłoszona jako wada, a nie
	// zniknęła w ciszy. Miarą jest to, że Operator traci jedną odpowiedź zamiast
	// całej sesji.
	defer func() {
		powod := recover()
		if powod == nil {
			return
		}
		if a != nil && a.mosty != nil && a.mosty.dziennik != nil {
			a.mosty.dziennik.Printf("tura okna %q przerwana usterką rdzenia: %v\n%s",
				okno.Id, powod, debug.Stack())
		}
	}()

	defer a.zapomnijBieg(okno.Id)

	strumien := nowyNadawcaStrumienia(a.nadajnik, idZadania, okno.IdSesji)
	// Wywołanie niżej domyka strumień także wtedy, gdy tura wyleci stąd
	// panicznie: bez niego klient czekałby na znacznik końca, który nigdy by
	// nie przyszedł, a wpis rozmowy stałby w stanie `strumien` na zawsze.
	defer strumien.DomknijAwaryjnie(okno.Id, odpowiedz.Id)
	var tresc strings.Builder
	var idRozmowy string
	var zamkniecie *zamkniecieTury
	ujscie := models.UjscieFunkcji(func(ctx context.Context, f models.Fragment) error {
		if f.Kind == shared.ChunkKindText {
			tresc.WriteString(models.TrescFragmentu(f))
		}
		// Zamknięcie tury przychodzi w ostatnim fragmencie. Zdejmujemy je
		// tutaj, bo po zamknięciu strumienia nie ma już do niego drogi:
		// identyfikator rozmowy domyka ciągłość, a całe zdarzenie rozstrzyga
		// stan odpowiedzi i jedzie do tabeli zamknięć.
		if z, jest := zamkniecieZFragmentu(f.Data); jest {
			zamkniecie = &z
			idRozmowy = z.IdRozmowyCLI
		}
		// Bloki nietekstowe (rozumowanie, narzędzia, prowenancja) jadą do bazy
		// w trakcie tury, nie po niej — po zamknięciu strumienia nie ma już do
		// nich drogi, a tekst i tak domknie dziennik.
		if a.bloki != nil {
			a.bloki.Zanotuj(f)
		}
		if a.petla != nil {
			a.petla.ObserwujFragment(f) // koordynator widzi strumień wykonawcy
		}
		return strumien.Fragment(ctx, f)
	})

	// Załączniki rozstrzygamy przed złożeniem zapytania: odwołanie niosące bajty
	// staje się ścieżką na nośniku, a ścieżka wchodzi do treści pytania — kanał
	// CLI nie ma pola na załącznik (adapter_rozmowa_zalaczniki.go).
	zapytanie := zapytanieKanalu(okno, pytanie, odpowiedz, rozwiazZalaczniki(a.zalaczniki, pytanie.Attachments))
	// Pamięć rozmowy: wcześniejsze wypowiedzi tego okna, bez bieżącej (ta jedzie
	// w Tresc). Kanał bezstanowy (api) składa z niej kontekst wywołania, dzięki
	// czemu odpowiada na całą rozmowę, a nie na turę od zera. Kanał
	// z własną pamięcią (cli przez Wznowienie) może ją pominąć.
	zapytanie.Historia = a.historiaRozmowy(okno.Id, pytanie.Id)
	// Wznowienie ustawiamy przed turą: puste znaczy „rozmowa nowa", niepuste
	// trafia do `--resume` w argumentach procesu (injection/argumenty.go).
	if a.ciaglosc != nil {
		zapytanie.Wznowienie = a.ciaglosc.Przypomnij(kontekst, okno.Id)
	}
	zapytanie.Nakladka = a.nakladkaOkna(kontekst, okno)
	a.uzupelnijSrodowisko(kontekst, okno, &zapytanie)
	// Konfiguracja sesji dokłada rozstrzygnięcia obowiązujące: model i konto
	// wywołania oraz powierzchnię procesu (ustawienia, środowisko, MCP). Idzie po
	// uzupelnijSrodowisko, bo wskazanie per sesja/okno wygrywa z wartością okna
	// tam, gdzie obszar jest wypełniony.
	a.uzupelnijKonfiguracje(kontekst, okno, &zapytanie)
	// Ekspert nakłada się na końcu — model, ustawienia i mosty MCP wskazane
	// przez niego mają ostatnie słowo, tak samo jak jego warstwy promptu.
	a.uzupelnijAgenta(kontekst, okno, &zapytanie)

	err := a.kanaly.Wyslij(kontekst, zapytanie, ujscie)
	// Kanał odmówił, zanim cokolwiek doszło do odbiorcy — tura może pojechać
	// kanałem zapasowym wiersza rejestru, jawnie (adapter_rozmowa_zapas.go).
	if zapasMozliwy(kontekst, err, tresc.Len() > 0, zamkniecie != nil) {
		err, _ = a.pojedzZapasem(kontekst, zapytanie, ujscie, err)
	}
	// Utrwalenie idzie po turze i niezależnie od jej powodzenia: rozmowa mogła
	// powstać po stronie programu `claude` nawet wtedy, gdy tura skończyła się
	// błędem, a jej porzucenie kosztowałoby całą historię kontekstu.
	if a.ciaglosc != nil && idRozmowy != "" {
		// Nieutrwalona ciągłość nie przerywa tury — następna zacznie rozmowę
		// od nowa, co jest gorsze, ale nie jest awarią.
		_ = a.ciaglosc.Zapamietaj(kontekst, okno.Id, idRozmowy)
	}
	strumien.Zakoncz(okno.Id, odpowiedz.Id, tresc.String(), err)
	if a.petla != nil {
		a.petla.ZakonczTure(okno.Id, powodTury(kontekst, err, zamkniecie)) // koniec tury wybudza koordynatora
	}

	odpowiedz.Content = tresc.String()
	// Stan odpowiedzi rozstrzyga zdarzenie zamknięcia tury, nie słowo modelu:
	// `is_error` zdarzenia `result` wygrywa z treścią.
	odpowiedz.Status = stanOdpowiedziZeZdarzen(kontekst, err, zamkniecie)
	if a.zdarzenia != nil {
		a.zdarzenia.ZanotujZamkniecie(okno.Id, odpowiedz.Id, zamkniecie, odpowiedz.Status)
	}
	a.dziennik.Zmien(odpowiedz)
	a.rozglosWiadomosc(kontekst, odpowiedz)
}

// zapytanieKanalu składa zapytanie kanału z parametrów okna. Kanał nie sięga po
// konfigurację sam — parametry przychodzą z okna, najwęższego poziomu zasięgu.
//
// Załączniki jadą w treści, bo kanał nie ma na nie pola: `models.Zapytanie`
// przekłada się na `injection.Zapytanie`, które niesie jedno wejście rozmowy
// (`Tekst`, stdin JSON-lines). Wplecione są odwołaniami nazwanymi (ścieżkami),
// nie treścią: model sięga po plik narzędziem.
func zapytanieKanalu(okno session.Okno, pytanie, odpowiedz shared.Message, zalaczniki []zalacznikTury) models.Zapytanie {
	return models.Zapytanie{
		Zasiegi:             models.Zasiegi{Sesja: okno.IdSesji, Okno: okno.Id},
		Wiadomosc:           odpowiedz.Id,
		Tresc:               wplecZalaczniki(pytanie.Content, zalaczniki),
		Kanal:               okno.KanalModelu,
		KatalogiRobocze:     okno.KatalogiRobocze,
		SrodowiskoWykonania: okno.SrodowiskoWykonania,
		TrybUprawnien:       okno.TrybUprawnien,
		RolaOkna:            okno.RolaOkna,
	}
}

// metadaneNadania są tą częścią obszaru `Metadata` wiadomości, którą warstwa
// trwałości przenosi do kolumn tabeli `wiadomosc` (`persona`, `okno_zrodlowe_id`).
// Kształt jest powtórzony po nazwach kluczy, a nie zaimportowany: odpowiadający
// mu typ warstwy danych jest jej prywatny, a rdzeń do środka warstwy nie sięga.
// Klucze muszą zgadzać się co do znaku z `metadaneWiadomosci` w `dane/rozmowa.go`,
// bo to one zamykają obieg zapis → kolumna → odczyt.
//
// Zużycia tokenów tu nie ma: podsumowanie tury kanału głównego
// (`injection.ZakonczenieTury`) niesie koszt, liczbę tur i czas trwania, ale
// surowe `usage` zdarzenia `result` zostaje w pakiecie `injection`. Wpisanie zer
// wyglądałoby na zmierzone „zero tokenów", więc pola zostają nieobecne.
type metadaneNadania struct {
	Persona        string `json:"persona,omitempty"`
	OknoZrodloweId string `json:"sourceWindowId,omitempty"`
}

// atrybucjaOkna składa metadane wiadomości nadawanej w oknie. Obie wartości
// czyta z samego okna — rdzeń nie dolicza tu niczego, czego okno nie wie.
//
// Okno bez atrybucji (samodzielne, poza pętlą) daje obszar pusty, a nie obiekt
// z pustymi polami: kontrakt zna `metadata` jako pole opcjonalne, a wiadomość
// zwykłej rozmowy nie ma czego atrybuować.
func atrybucjaOkna(okno session.Okno) json.RawMessage {
	meta := metadaneNadania{
		Persona: personaOkna(okno.RolaOkna),
		// Okno źródłowe wypełnia się wyłącznie dla wykonawcy: to okno
		// koordynatora zleciło mu turę, więc ono jest źródłem wypowiedzi —
		// atrybucję w pętli niesie kolumna `okno_zrodlowe_id`. Dla koordynatora
		// i okna samodzielnego pole zostaje puste.
		OknoZrodloweId: okno.OknoKoordynatora,
	}
	if meta.Persona == "" && meta.OknoZrodloweId == "" {
		return nil
	}
	surowe, err := json.Marshal(meta)
	if err != nil {
		// Niezłożona atrybucja nie może zabrać wiadomości — wiadomość jedzie
		// bez metadanych, kolumny zostają puste.
		return nil
	}
	return surowe
}

// personaOkna przekłada rolę okna na personę wiadomości.
//
// Słownik kolumny jest węższy niż słownik kontraktu: `persona` przyjmuje
// wyłącznie role pętli i tory pochodne, okna samodzielnego w nim nie ma,
// a wpisanie „standalone" odbiłoby się warunkiem CHECK i wywróciłoby zapis całej
// wiadomości. Dlatego okno poza pętlą zostaje bez persony — kolumna pusta znaczy
// zwykłą wypowiedź Operatora albo modelu.
func personaOkna(rola shared.WindowRole) string {
	switch rola {
	case shared.WindowRoleCoordinator, shared.WindowRoleExecutor:
		return string(rola)
	default:
		return ""
	}
}

// historiaRozmowy składa pamięć wcześniejszych tur okna dla kanału bezstanowego.
// Bierze wiadomości sprzed bieżącego pytania (dziennik ma je już zapisane —
// Wyslij dopisał pytanie i odpowiedź przed startem tury), pomija wypowiedzi bez
// roli user/assistant oraz puste, i oddaje je w kolejności nadania. Pusta
// historia nie jest błędem — pierwsza tura okna po prostu nie ma pamięci.
func (a *adapterRozmowy) historiaRozmowy(idOkna, idPytania string) []models.WiadomoscHistorii {
	if a.dziennik == nil {
		return nil
	}
	wczesniejsze, _ := a.dziennik.Wykaz(idOkna, &idPytania, nil)
	historia := make([]models.WiadomoscHistorii, 0, len(wczesniejsze))
	for _, w := range wczesniejsze {
		rola := rolaHistorii(w.Role)
		if rola == "" || strings.TrimSpace(w.Content) == "" {
			continue
		}
		historia = append(historia, models.WiadomoscHistorii{Rola: rola, Tresc: w.Content})
	}
	return historia
}

// rolaHistorii przekłada rolę wiadomości kontraktu na rolę pamięci wywołania.
// Wypowiedzi spoza rozmowy (systemowe, narzędziowe) nie wchodzą do pamięci —
// zwrócony pusty napis odsiewa je w historiaRozmowy.
func rolaHistorii(rola shared.MessageRole) string {
	switch rola {
	case shared.MessageRoleUser:
		return models.RolaUzytkownika
	case shared.MessageRoleAssistant:
		return models.RolaAsystenta
	default:
		return ""
	}
}

// rozglosWiadomosc zawiadamia urządzenia konta o zamkniętej odpowiedzi.
//
// Sprawcą jest rdzeń, nie ten, kto turę otworzył: rozgłaszana wiadomość to
// odpowiedź zamknięta przez rdzeń po turze modelu, więc nikt jej nie wpisał
// ręką. Autor promptu jest opisany zdarzeniem założenia wiadomości
// w `handlers_message.go`.
func (a *adapterRozmowy) rozglosWiadomosc(ctx context.Context, w shared.Message) {
	nowyEmiter(a.nadajnik).wiadomosc(zSprawcaRdzenia(ctx), shared.ChangeKindUpdated, w)
}
