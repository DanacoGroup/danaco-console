// Odpowiedzialność pliku: port Rozmowa — przyjęcie wiadomości okna, tura kanału modelu gorutyną i odesłanie odpowiedzi strumieniem.
package core

import (
	"context"
	"encoding/json"
	"runtime/debug"
	"strings"
	"sync"
	"time"

	"danacoconsole/server/internal/dane"
	"danacoconsole/server/internal/models"
	"danacoconsole/server/internal/protocol"
	"danacoconsole/server/internal/session"
	"danacoconsole/shared"
)

type adapterRozmowy struct {
	nadzorca *session.Nadzorca
	kanaly   *models.Rejestr
	// repozytoriumKanalow rozstrzyga własność kanału z żądania (decyzja 34).
	repozytoriumKanalow dane.RepozytoriumKanalow
	nadajnik            Nadajnik
	dziennik            *dziennikRozmowy
	zycie               context.Context
	petla               *session.Petla
	tozsamosc           Tozsamosc
	agenci              ZrodloTozsamosciAgenta
	dolozeniaSesji      DolozeniaNarzedziSesji
	zgloszoneZestawy    sync.Map
	katalog             *KatalogRoboczy
	mosty               *mostyOkna
	ciaglosc            CiagloscRozmowy
	wykonanie           ParametryWykonania
	konfiguracja        CzytelnikKonfiguracjiSesji
	zdarzenia           *zdarzeniaWykonawcze
	bloki               *rejestratorBlokow
	zalaczniki          *magazynTresciBiblioteki
	tor                 *torStrumieni

	mu       sync.Mutex
	biegnace map[string]*biegTury
}

// Kontekst życia jest kontekstem rdzenia, nie połączenia: rozłączenie klienta nie przerywa rozpoczętej tury.
func nowyAdapterRozmowy(zycie context.Context, nadzorca *session.Nadzorca,
	kanaly *models.Rejestr, nadajnik Nadajnik, dziennik *dziennikRozmowy) *adapterRozmowy {
	return &adapterRozmowy{
		nadzorca: nadzorca, kanaly: kanaly, nadajnik: nadajnik, dziennik: dziennik,
		zycie: zycie, biegnace: map[string]*biegTury{},
	}
}

func (a *adapterRozmowy) UstawPetle(p *session.Petla) { a.petla = p }

func (a *adapterRozmowy) ZCiagloscia(c CiagloscRozmowy) *adapterRozmowy {
	a.ciaglosc = c
	return a
}

func (a *adapterRozmowy) ZZdarzeniami(z *zdarzeniaWykonawcze) *adapterRozmowy {
	a.zdarzenia = z
	return a
}

func (a *adapterRozmowy) ZBlokami(b *rejestratorBlokow) *adapterRozmowy {
	a.bloki = b
	return a
}

func (a *adapterRozmowy) ZZalacznikami(katalogDanych string) *adapterRozmowy {
	a.zalaczniki = magazynZalacznikow(katalogDanych)
	return a
}

func (a *adapterRozmowy) ZTorem(t *torStrumieni) *adapterRozmowy {
	a.tor = t
	return a
}

func (a *adapterRozmowy) Wyslij(ctx context.Context, z shared.MessageSendRequest) (shared.MessageSendResponse, error) {
	okno, err := a.nadzorca.Rejestr().Okno(z.WindowId)
	if err != nil {
		return shared.MessageSendResponse{}, bladSesji(err)
	}
	// Zajęcie okna idzie przed dziennikiem; tura biegnie na koncie zamawiającego (decyzja 34).
	kontekst, anuluj := context.WithCancel(zKontemZadania(a.zycie, ctx))
	bieg := a.zajmijBieg(okno.Id, anuluj)
	if bieg == nil {
		anuluj()
		return shared.MessageSendResponse{}, odmowaTuryWBiegu(okno.Id)
	}

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

	idZadania := protocol.TozsamoscStrumienia(ctx, pytanie.Id)
	a.tor.zwiaz(idZadania, ujscieZKontekstu(ctx))

	go a.prowadzTure(kontekst, okno, pytanie, odpowiedz, idZadania, bieg)

	return shared.MessageSendResponse{Message: pytanie}, nil
}

func (a *adapterRozmowy) Wykaz(_ context.Context, z shared.MessageListRequest) (shared.MessageListResponse, error) {
	wiadomosci, sastarsze := a.dziennik.Wykaz(z.WindowId, z.Before, z.Limit)
	return shared.MessageListResponse{Messages: wiadomosci, HasMore: sastarsze}, nil
}

func (a *adapterRozmowy) prowadzTure(kontekst context.Context, okno session.Okno,
	pytanie, odpowiedz shared.Message, idZadania string, bieg *biegTury) {

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

	defer a.zapomnijBieg(okno.Id, bieg)

	strumien := nowyNadawcaStrumienia(a.nadajnik, kontoAdresata(kontekst), idZadania, okno.IdSesji)
	defer strumien.DomknijAwaryjnie(okno.Id, odpowiedz.Id)
	var tresc strings.Builder
	var idRozmowy string
	var zamkniecie *zamkniecieTury
	ujscie := models.UjscieFunkcji(func(ctx context.Context, f models.Fragment) error {
		if f.Kind == shared.ChunkKindText {
			tresc.WriteString(models.TrescFragmentu(f))
		}
		if z, jest := zamkniecieZFragmentu(f.Data); jest {
			zamkniecie = &z
			idRozmowy = z.IdRozmowyCLI
		}
		if a.bloki != nil {
			a.bloki.Zanotuj(kontekst, f)
		}
		if a.petla != nil {
			a.petla.ObserwujFragment(f)
		}
		return strumien.Fragment(ctx, f)
	})

	zapytanie := zapytanieKanalu(okno, pytanie, odpowiedz, rozwiazZalaczniki(a.zalaczniki, pytanie.Attachments))
	zapytanie.Historia = a.historiaRozmowy(okno.Id, pytanie.Id)
	// Wznowienie puste znaczy „rozmowa nowa", niepuste trafia do `--resume`.
	if a.ciaglosc != nil {
		zapytanie.Wznowienie = a.ciaglosc.Przypomnij(kontekst, okno.Id)
	}
	zapytanie.Nakladka = a.nakladkaOkna(kontekst, okno)
	a.uzupelnijSrodowisko(kontekst, okno, &zapytanie)
	// Kolejność nakładek: środowisko okna, konfiguracja sesji, ekspert — późniejsza wygrywa.
	a.uzupelnijKonfiguracje(kontekst, okno, &zapytanie)
	a.uzupelnijAgenta(kontekst, okno, &zapytanie)

	err := a.kanaly.Wyslij(kontekst, zapytanie, ujscie)
	if zapasMozliwy(kontekst, err, tresc.Len() > 0, zamkniecie != nil) {
		err, _ = a.pojedzZapasem(kontekst, zapytanie, ujscie, err)
	}
	if a.ciaglosc != nil && idRozmowy != "" {
		// Nieutrwalona ciągłość nie przerywa tury; następna zaczyna od zera.
		_ = a.ciaglosc.Zapamietaj(kontekst, okno.Id, idRozmowy)
	}
	strumien.Zakoncz(okno.Id, odpowiedz.Id, tresc.String(), err)
	if a.petla != nil {
		a.petla.ZakonczTure(okno.Id, powodTury(kontekst, err, zamkniecie))
	}

	odpowiedz.Content = tresc.String()
	// Powód odmowy wchodzi w pusty wpis: Operator czyta wpis rozmowy, nie dziennik serwera.
	if err != nil && odpowiedz.Content == "" {
		odpowiedz.Content = err.Error()
	}
	// Stan odpowiedzi rozstrzyga zdarzenie zamknięcia tury, nie słowo modelu.
	odpowiedz.Status = stanOdpowiedziZeZdarzen(kontekst, err, zamkniecie)
	if a.zdarzenia != nil {
		a.zdarzenia.ZanotujZamkniecie(okno.Id, odpowiedz.Id, zamkniecie, odpowiedz.Status)
	}
	a.dziennik.Zmien(odpowiedz)
	a.rozglosWiadomosc(kontekst, odpowiedz)
}

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

// metadaneNadania powtarza klucze warstwy danych zamiast je importować, bo tamten typ jest prywatny.
type metadaneNadania struct {
	Persona        string `json:"persona,omitempty"`
	OknoZrodloweId string `json:"sourceWindowId,omitempty"`
}

func atrybucjaOkna(okno session.Okno) json.RawMessage {
	meta := metadaneNadania{
		Persona:        personaOkna(okno.RolaOkna),
		OknoZrodloweId: okno.OknoKoordynatora,
	}
	if meta.Persona == "" && meta.OknoZrodloweId == "" {
		return nil
	}
	surowe, err := json.Marshal(meta)
	if err != nil {
		// Niezłożona atrybucja nie zabiera wiadomości; jedzie bez metadanych.
		return nil
	}
	return surowe
}

func personaOkna(rola shared.WindowRole) string {
	switch rola {
	case shared.WindowRoleCoordinator, shared.WindowRoleExecutor:
		return string(rola)
	default:
		return ""
	}
}

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

func (a *adapterRozmowy) rozglosWiadomosc(ctx context.Context, w shared.Message) {
	nowyEmiter(a.nadajnik).wiadomosc(zSprawcaRdzenia(ctx), shared.ChangeKindUpdated, w)
}
