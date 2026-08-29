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
// kanału modelu z rejestru i odsyła odpowiedź strumieniem, podczas gdy komenda
// potwierdza przyjęcie wiadomości od razu, bez czekania na koniec odpowiedzi.
type adapterRozmowy struct {
	nadzorca *session.Nadzorca
	kanaly   *models.Rejestr
	nadajnik Nadajnik
	dziennik *dziennikRozmowy
	zycie    context.Context
	// petla prowadzi bieg koordynator–wykonawca; okno samodzielne pracuje bez niej.
	petla *session.Petla
	// tozsamosc podaje nakładkę obowiązującą okna.
	tozsamosc Tozsamosc
	// agenci podają tożsamość eksperta okna; port pusty znaczy model surowy.
	agenci ZrodloTozsamosciAgenta
	// dolozeniaSesji podaje doraźne dołożenia narzędzi tury, obok eksperta okna.
	dolozeniaSesji DolozeniaNarzedziSesji
	// zgloszoneZestawy pilnuje, by powód zestawu szedł do dziennika raz na
	// powód, a nie raz na turę.
	zgloszoneZestawy sync.Map
	// katalog ustala katalog roboczy sesji, mosty — konfigurację MCP ze zbioru
	// nadań okna.
	katalog *KatalogRoboczy
	mosty   *mostyOkna
	// ciaglosc utrwala identyfikator rozmowy CLI; port pusty zaczyna każdą turę od zera.
	ciaglosc CiagloscRozmowy
	// wykonanie podaje nakład rozumowania i model zapasowy z konfiguracji.
	wykonanie ParametryWykonania
	// konfiguracja tłumaczy obowiązującą konfigurację sesji na wejście procesu.
	konfiguracja CzytelnikKonfiguracjiSesji
	// zdarzenia utrwala zamknięcia tur; port pusty zostawia rozmowę bez dowodu rozstrzygnięcia.
	zdarzenia *zdarzeniaWykonawcze
	// bloki utrwala fragmenty nietekstowe tury; port pusty zwraca historię samym tekstem.
	bloki *rejestratorBlokow
	// zalaczniki jest magazynem bajtów załącznika; pusty zostawia go niedoręczonym, nazwanym wprost.
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

// ZZalacznikami wpina magazyn załączników nad katalogiem danych procesu —
// tym samym, który montaż podaje sejfowi poświadczeń i bibliotece. Adapter
// bez tego portu wysyła załącznik do modelu jako niedoręczony, nie jako ciszę.
func (a *adapterRozmowy) ZZalacznikami(katalogDanych string) *adapterRozmowy {
	a.zalaczniki = magazynZalacznikow(katalogDanych)
	return a
}

// Wyslij przyjmuje wiadomość użytkownika, dopisuje ją do dziennika rozmowy
// i otwiera turę okna gorutyną osobną od odpowiedzi tej komendy.
func (a *adapterRozmowy) Wyslij(ctx context.Context, z shared.MessageSendRequest) (shared.MessageSendResponse, error) {
	okno, err := a.nadzorca.Rejestr().Okno(z.WindowId)
	if err != nil {
		return shared.MessageSendResponse{}, bladSesji(err)
	}
	// Zajęcie okna idzie przed dziennikiem, żeby odmowa nie zostawiła pytania bez odpowiedzi.
	kontekst, anuluj := context.WithCancel(a.zycie)
	if !a.zajmijBieg(okno.Id, anuluj) {
		anuluj()
		return shared.MessageSendResponse{}, odmowaTuryWBiegu(okno.Id)
	}

	// Obie wiadomości tury niosą tę samą atrybucję okna i roli pętli.
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

	// Tożsamość strumienia zdejmuje się tu — kontekst tury nie niesie już wpisu żądania.
	idZadania := protocol.TozsamoscStrumienia(ctx, pytanie.Id)

	go a.prowadzTure(kontekst, okno, pytanie, odpowiedz, idZadania)

	return shared.MessageSendResponse{Message: pytanie}, nil
}

// Wykaz zwraca wiadomości okna z dziennika rozmowy wraz ze znacznikiem,
// czy poza zwróconą stroną stoją wiadomości starsze.
func (a *adapterRozmowy) Wykaz(_ context.Context, z shared.MessageListRequest) (shared.MessageListResponse, error) {
	wiadomosci, sastarsze := a.dziennik.Wykaz(z.WindowId, z.Before, z.Limit)
	return shared.MessageListResponse{Messages: wiadomosci, HasMore: sastarsze}, nil
}

// prowadzTure wykonuje jedną turę kanału modelu i odsyła jej strumień; biegnie
// własną gorutyną, więc każda usterka w jej torze musi zostać przechwycona tutaj.
func (a *adapterRozmowy) prowadzTure(kontekst context.Context, okno session.Okno,
	pytanie, odpowiedz shared.Message, idZadania string) {

	// Osłona gorutyny tury: bez niej panika w torze tury gasi cały rdzeń.
	defer func() {
		powod := recover()
		if powod == nil {
			return
		}
		if a != nil && a.mosty != nil && a.mosty.dziennik != nil {
			a.mosty.dziennik.Printf("tura okna %q przerwana usterką serwera: %v\n%s",
				okno.Id, powod, debug.Stack())
		}
	}()

	defer a.zapomnijBieg(okno.Id)

	strumien := nowyNadawcaStrumienia(a.nadajnik, idZadania, okno.IdSesji)
	// Domyka strumień także wtedy, gdy tura wyleci stąd panicznie.
	defer strumien.DomknijAwaryjnie(okno.Id, odpowiedz.Id)
	var tresc strings.Builder
	var idRozmowy string
	var zamkniecie *zamkniecieTury
	ujscie := models.UjscieFunkcji(func(ctx context.Context, f models.Fragment) error {
		if f.Kind == shared.ChunkKindText {
			tresc.WriteString(models.TrescFragmentu(f))
		}
		// Zamknięcie tury zdejmuje się tu — po zamknięciu strumienia nie ma już do niego drogi.
		if z, jest := zamkniecieZFragmentu(f.Data); jest {
			zamkniecie = &z
			idRozmowy = z.IdRozmowyCLI
		}
		// Bloki nietekstowe jadą do bazy w trakcie tury, nie po niej.
		if a.bloki != nil {
			a.bloki.Zanotuj(f)
		}
		if a.petla != nil {
			a.petla.ObserwujFragment(f) // koordynator widzi strumień wykonawcy
		}
		return strumien.Fragment(ctx, f)
	})

	// Załączniki rozstrzyga się przed zapytaniem: bajty stają się ścieżką w treści pytania.
	zapytanie := zapytanieKanalu(okno, pytanie, odpowiedz, rozwiazZalaczniki(a.zalaczniki, pytanie.Attachments))
	// Pamięć rozmowy niesie wcześniejsze wypowiedzi okna dla kanału bezstanowego.
	zapytanie.Historia = a.historiaRozmowy(okno.Id, pytanie.Id)
	// Wznowienie puste znaczy „rozmowa nowa", niepuste trafia do `--resume`.
	if a.ciaglosc != nil {
		zapytanie.Wznowienie = a.ciaglosc.Przypomnij(kontekst, okno.Id)
	}
	zapytanie.Nakladka = a.nakladkaOkna(kontekst, okno)
	a.uzupelnijSrodowisko(kontekst, okno, &zapytanie)
	// Konfiguracja sesji idzie po środowisku okna — wskazanie sesji wygrywa.
	a.uzupelnijKonfiguracje(kontekst, okno, &zapytanie)
	// Ekspert nakłada się na końcu — jego ustawienia mają ostatnie słowo.
	a.uzupelnijAgenta(kontekst, okno, &zapytanie)

	err := a.kanaly.Wyslij(kontekst, zapytanie, ujscie)
	// Kanał odmówił bez dostawy — tura może pojechać kanałem zapasowym.
	if zapasMozliwy(kontekst, err, tresc.Len() > 0, zamkniecie != nil) {
		err, _ = a.pojedzZapasem(kontekst, zapytanie, ujscie, err)
	}
	// Utrwalenie idzie po turze niezależnie od jej powodzenia.
	if a.ciaglosc != nil && idRozmowy != "" {
		// Nieutrwalona ciągłość nie przerywa tury — następna zaczyna od zera.
		_ = a.ciaglosc.Zapamietaj(kontekst, okno.Id, idRozmowy)
	}
	strumien.Zakoncz(okno.Id, odpowiedz.Id, tresc.String(), err)
	if a.petla != nil {
		a.petla.ZakonczTure(okno.Id, powodTury(kontekst, err, zamkniecie)) // koniec tury wybudza koordynatora
	}

	odpowiedz.Content = tresc.String()
	// Stan odpowiedzi rozstrzyga zdarzenie zamknięcia tury, nie słowo modelu.
	odpowiedz.Status = stanOdpowiedziZeZdarzen(kontekst, err, zamkniecie)
	if a.zdarzenia != nil {
		a.zdarzenia.ZanotujZamkniecie(okno.Id, odpowiedz.Id, zamkniecie, odpowiedz.Status)
	}
	a.dziennik.Zmien(odpowiedz)
	a.rozglosWiadomosc(kontekst, odpowiedz)
}

// zapytanieKanalu składa zapytanie kanału z parametrów okna, najwęższego
// poziomu zasięgu; kanał nie sięga po konfigurację sam. Załączniki jadą
// w treści jako odwołania nazwane, ścieżkami, nie samą treścią bajtów.
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
// trwałości przenosi do kolumn tabeli `wiadomosc`; kształt jest powtórzony po
// nazwach kluczy warstwy danych, nie zaimportowany, bo ten typ jest jej prywatny.
type metadaneNadania struct {
	Persona        string `json:"persona,omitempty"`
	OknoZrodloweId string `json:"sourceWindowId,omitempty"`
}

// atrybucjaOkna składa metadane wiadomości nadawanej w oknie, czytane z samego
// okna; okno bez atrybucji daje obszar pusty, nie obiekt z pustymi polami.
func atrybucjaOkna(okno session.Okno) json.RawMessage {
	meta := metadaneNadania{
		Persona: personaOkna(okno.RolaOkna),
		// Okno źródłowe wypełnia się wyłącznie dla wykonawcy, zleceniodawcą tury.
		OknoZrodloweId: okno.OknoKoordynatora,
	}
	if meta.Persona == "" && meta.OknoZrodloweId == "" {
		return nil
	}
	surowe, err := json.Marshal(meta)
	if err != nil {
		// Niezłożona atrybucja nie może zabrać wiadomości — jedzie bez metadanych.
		return nil
	}
	return surowe
}

// personaOkna przekłada rolę okna na personę wiadomości; słownik kolumny jest
// węższy niż słownik kontraktu, więc okno poza pętlą zostaje bez persony.
func personaOkna(rola shared.WindowRole) string {
	switch rola {
	case shared.WindowRoleCoordinator, shared.WindowRoleExecutor:
		return string(rola)
	default:
		return ""
	}
}

// historiaRozmowy składa pamięć wcześniejszych tur okna dla kanału bezstanowego,
// pomijając wypowiedzi bez roli user/assistant oraz puste. Pusta historia nie
// jest błędem — pierwsza tura okna po prostu nie ma pamięci.
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

// rozglosWiadomosc zawiadamia urządzenia konta o zamkniętej odpowiedzi;
// sprawcą jest rdzeń, nie ten, kto turę otworzył, bo nikt jej nie wpisał ręką.
func (a *adapterRozmowy) rozglosWiadomosc(ctx context.Context, w shared.Message) {
	nowyEmiter(a.nadajnik).wiadomosc(zSprawcaRdzenia(ctx), shared.ChangeKindUpdated, w)
}
