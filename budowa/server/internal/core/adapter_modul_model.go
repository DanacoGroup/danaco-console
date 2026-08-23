// Odpowiedzialność pliku: rodzina `model.*` — jedna komenda kontraktu,
// `model.channel.set`, czyli wybór kanału modelu obsługującego okno komunikacji
// albo kartę sesji. Kanały zakłada i zmienia rodzina `channel.*`; ta komenda
// niczego nie zakłada, wyłącznie wskazuje jeden z wierszy rejestru.
//
// To rozszerzenie portu okien, nie drugi adapter. Metody poniżej stoją na
// `adapterOkien` — tym samym bycie, którym jedzie `window.update`. Kanał okna ma
// w rdzeniu jednego właściciela: gdyby rodzina `model.*` dostała własny adapter
// nad własnym dojściem do rejestru okien, ten sam kanał miałby dwa miejsca
// zmiany i dwie prawdy o tym, który wygrywa.
//
// Kanał okna mieszka w dwóch miejscach naraz i tylko jedno z nich ma wpływ na
// wywołanie modelu:
//
//	okno_komunikacji.kanal_modelu_id  — kolumna wiersza okna oraz jej
//	                                    odpowiednik pamięciowy
//	                                    `session.Okno.KanalModelu`. To ona jedzie
//	                                    do wywołania tury i to ją odbudowuje
//	                                    `odtworzenie_stanu.go` po restarcie.
//	ustawienie.klucz = 'kanal_modelu' — pozycja katalogu ustawień, ustawialna na
//	                                    wszystkich ośmiu poziomach zasięgu,
//	                                    w tym `karta_sesji` i `okno`.
//	                                    Rozstrzygacz konfiguracji jej nie czyta —
//	                                    żaden kod nie sięga po ten klucz,
//	                                    `konfig/definicje_wykonania.go` tylko go
//	                                    deklaruje.
//
// Komenda zapisuje więc miejsce pierwsze — to, które działa. Zapis do pozycji
// konfiguracji byłby ciszą udającą skutek: Operator przestawiłby kanał, panel
// pokazałby zmianę, a tura poszłaby starym kanałem.
package core

import (
	"context"
	"errors"

	"danacoconsole/server/internal/dane"
	"danacoconsole/server/internal/protocol"
	"danacoconsole/server/internal/session"
	"danacoconsole/shared"
)

// UstawKanalModelu obsługuje `model.channel.set`.
//
// Zwraca odpowiedź kontraktu oraz wykaz okien, których kanał został zmieniony —
// rozgłoszeniem `window.changed` zajmuje się wpięcie (handlers_model.go), tak
// samo jak przy `window.update`. Żądanie karty sesji dotyka wielu okien naraz,
// więc pojedyncze okno odpowiedzi by tego nie pokazało.
func (a *adapterOkien) UstawKanalModelu(ctx context.Context,
	z shared.ModelChannelSetRequest) (shared.ModelChannelSetResponse, []shared.Window, error) {

	kanal, err := a.kanalRejestru(ctx, z.ModelChannelId)
	if err != nil {
		return shared.ModelChannelSetResponse{}, nil, err
	}
	idOkna := wartoscTekstu(z.WindowId)
	idSesji := wartoscTekstu(z.SessionId)
	switch {
	case idOkna != "":
		return a.kanalOkna(ctx, kanal, idOkna, z.AgentId)
	case idSesji != "":
		return a.kanalKartySesji(ctx, kanal, idSesji, z.AgentId)
	default:
		return shared.ModelChannelSetResponse{}, nil, bladKanaluModelu(
			"żądanie bez wskazania okna i bez wskazania karty sesji; " +
				"nie ma czemu nadać kanału")
	}
}

// kanalRejestru odnajduje wiersz rejestru kanałów po kodzie z żądania.
//
// Kanał wybrany z rejestru — nie zakładany. Kodu, którego w rejestrze nie ma,
// nie da się nadać oknu: kolumna `okno_komunikacji.kanal_modelu_id` jest kluczem
// obcym z ON DELETE RESTRICT, więc zapis i tak by nie przeszedł, a okno w pamięci
// wskazywałoby kanał widmo. Odmowa niesie kod `not_found` i nazwę bytu.
//
// Wykaz idzie po wszystkie wiersze, także nieczynne: kanał wyłączony istnieje
// i Operator ma prawo go wskazać, a odpowiedź mówi o tym wprost polem `enabled`.
// Odmowa wyboru kanału wyłączonego byłaby zasadą, której kontrakt nie stawia.
func (a *adapterOkien) kanalRejestru(ctx context.Context, kod string) (dane.Kanal, error) {
	if kod == "" {
		return dane.Kanal{}, bladKanaluModelu("żądanie bez wskazania kanału modelu")
	}
	if a.kanaly == nil {
		return dane.Kanal{}, protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeInternalError,
			"rodzina model.*: rejestr kanałów niewpięty, kanału nie da się rozpoznać"))
	}
	kanaly, err := a.kanaly.Lista(ctx, false)
	if err != nil {
		return dane.Kanal{}, err
	}
	for _, kanal := range kanaly {
		if kanal.Kod == kod {
			return kanal, nil
		}
	}
	return dane.Kanal{}, protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeNotFound,
		"rodzina model.*: kanał modelu "+kod+" nie istnieje w rejestrze"))
}

// kanalOkna nadaje kanał jednemu oknu komunikacji wraz z ekspertem żądania.
//
// Ekspert jedzie razem z kanałem, bo kontrakt tak go podaje polem
// `ModelChannelSetRequest.agentId`. Pominięcie go dawałoby Operatorowi
// wybierającemu eksperta z menu „Modele" potwierdzenie bez skutku: kanał by się
// zmieniał, ekspert nie.
//
// Wskaźnik pusty zostawia wybór bez zmiany, wskaźnik na pusty napis zdejmuje
// eksperta — znaczenie jest jedno i pochodzi z `session.Zmiana`, żeby ta sama
// wartość nie znaczyła tu czegoś innego niż w `window.update`.
func (a *adapterOkien) kanalOkna(ctx context.Context, kanal dane.Kanal,
	idOkna string, agent *string) (shared.ModelChannelSetResponse, []shared.Window, error) {

	okno, err := a.nadzorca.Rejestr().ZmienOkno(idOkna,
		session.Zmiana{KanalModelu: &kanal.Kod, Agent: agent})
	if err != nil {
		return shared.ModelChannelSetResponse{}, nil, bladSesji(err)
	}
	if err := a.utrwalKanalOkna(ctx, idOkna, kanal.ID); err != nil {
		return shared.ModelChannelSetResponse{}, nil, err
	}
	if err := a.utrwalAgentaOkna(ctx, idOkna, agent); err != nil {
		return shared.ModelChannelSetResponse{}, nil, err
	}
	zmienione := []shared.Window{oknoKontraktu(okno)}
	return shared.ModelChannelSetResponse{
		Channel: kanalOdpowiedzi(kanal), WindowId: &idOkna,
	}, zmienione, nil
}

// kanalKartySesji nadaje kanał karcie sesji, czyli wszystkim jej otwartym oknom.
//
// Kontrakt dopuszcza wskazanie karty sesji zamiast okna, a schemat nie ma dla
// karty ani kolumny kanału, ani żadnego innego miejsca, w którym „kanał karty"
// mógłby zamieszkać z wpływem na pracę. Jedynym bytem, który kanał naprawdę
// obsługuje, jest okno. Nadanie kanału karcie znaczy więc: nadaj go każdemu
// oknu, które ta karta prowadzi.
//
// Dwa inne znaczenia są tu świadomie nieprzyjęte:
//   - zapis pozycji `kanal_modelu` na poziomie zasięgu `karta_sesji` — pozycja
//     istnieje w katalogu ustawień, lecz nikt jej nie czyta, więc Operator
//     dostałby potwierdzenie bez skutku;
//   - „kanał domyślny dla okien zakładanych później" — okno zakładane bez
//     wskazania bierze kanał z rejestru (`adapterOkien.kanalDomyslny`), a nie
//     z karty; druga reguła domyślności zrobiłaby z jednego pytania dwa źródła
//     prawdy.
//
// Okna zamknięte zostają nietknięte: zamknięte okno nie prowadzi pracy, a jego
// kanał jest zapisem tego, czym pracowało.
func (a *adapterOkien) kanalKartySesji(ctx context.Context, kanal dane.Kanal,
	idSesji string, agent *string) (shared.ModelChannelSetResponse, []shared.Window, error) {

	okna, err := a.nadzorca.Rejestr().OknaSesji(idSesji)
	if err != nil {
		return shared.ModelChannelSetResponse{}, nil, bladSesji(err)
	}
	zmienione := make([]shared.Window, 0, len(okna))
	for _, okno := range okna {
		if !okno.CzyOtwarte() {
			continue
		}
		zmienione, err = a.dopiszZmianeKanalu(ctx, zmienione, okno.Id, kanal, agent)
		if err != nil {
			return shared.ModelChannelSetResponse{}, nil, err
		}
	}
	// Karta bez ani jednego otwartego okna nie ma czemu nadać kanału. Odpowiedź
	// „kanał ustawiony" byłaby tu ciszą udającą skutek.
	if len(zmienione) == 0 {
		return shared.ModelChannelSetResponse{}, nil, protocol.JakoError(protocol.NowyBlad(
			shared.ErrorCodeConflict,
			"rodzina model.*: karta sesji "+idSesji+" nie prowadzi ani jednego otwartego okna; "+
				"kanał modelu obsługuje okno, nie kartę"))
	}
	return shared.ModelChannelSetResponse{
		Channel: kanalOdpowiedzi(kanal), SessionId: &idSesji,
	}, zmienione, nil
}

// kanalOdpowiedzi składa kanał kontraktu na odpowiedź tej komendy.
//
// Przekład wiersza rejestru jest jeden (`kanalKontraktu`, adapter_kanaly.go) —
// drugiego tu nie ma. Dochodzi wyłącznie czas utworzenia, którego tamten
// przekład nie nanosi: `channel.list` podaje go wypełnionego (rejestr kanałów
// warstwy modeli), więc zero w odpowiedzi tej komendy przeczyłoby tej samej
// wartości na tym samym ekranie. Znacznik nieczytelny zostawia zero — data
// pomocnicza nie wywraca odpowiedzi.
func kanalOdpowiedzi(k dane.Kanal) shared.Channel {
	kanal := kanalKontraktu(k)
	if chwila, jest := chwilaZapisu(k.Utworzono); jest {
		kanal.CreatedAt = chwila
	}
	return kanal
}

// dopiszZmianeKanalu nanosi kanał i eksperta na jedno okno karty i dopisuje je
// do wykazu zmienionych.
//
// Ekspert obejmuje te same okna co kanał — czyli wszystkie otwarte okna karty.
// Nadanie kanału karcie znaczy „nadaj go każdemu oknu, które ta karta prowadzi";
// wyłączenie eksperta z tej samej reguły zrobiłoby z jednego żądania dwa
// zasięgi i Operator nie miałby jak zgadnąć, który obowiązuje.
func (a *adapterOkien) dopiszZmianeKanalu(ctx context.Context, zmienione []shared.Window,
	idOkna string, kanal dane.Kanal, agent *string) ([]shared.Window, error) {

	okno, err := a.nadzorca.Rejestr().ZmienOkno(idOkna,
		session.Zmiana{KanalModelu: &kanal.Kod, Agent: agent})
	if err != nil {
		return nil, bladSesji(err)
	}
	if err := a.utrwalKanalOkna(ctx, idOkna, kanal.ID); err != nil {
		return nil, err
	}
	if err := a.utrwalAgentaOkna(ctx, idOkna, agent); err != nil {
		return nil, err
	}
	return append(zmienione, oknoKontraktu(okno)), nil
}

// utrwalKanalOkna zapisuje wybór w wierszu okna, żeby przeżył restart rdzenia.
//
// Parametry modelu są pamiętane per okno. Rejestr pamięciowy odtwarza się po
// restarcie z wierszy (`odtworzenie_stanu.go`), więc wybór zapisany wyłącznie
// w pamięci znikałby przy najbliższym uruchomieniu — a Operator nie miałby jak
// się o tym dowiedzieć.
//
// Brak wiersza nie jest błędem. Wiersz okna powstaje leniwie, przy pierwszej
// wiadomości (`dane/rozmowa_lancuch.go`), i bierze kanał wprost z opisu okna
// w rejestrze rdzenia — czyli z wartości ustawionej właśnie teraz. Zapis nie ma
// więc czego dogonić.
//
// Błąd zapisu jest błędem komendy. Wiersz istnieje, a nie przyjął zmiany —
// wybór nie przeżyje restartu i milczenie o tym byłoby obietnicą bez pokrycia.
func (a *adapterOkien) utrwalKanalOkna(ctx context.Context, idOkna string, kanalID int64) error {
	if a.trwalosc == nil || a.trwalosc.okna == nil {
		return nil
	}
	wiersz, err := a.trwalosc.okna.PoIdentyfikatorze(ctx, idOkna)
	if errors.Is(err, dane.ErrBrakWiersza) {
		return nil
	}
	if err != nil {
		return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeInternalError,
			"rodzina model.*: nie da się odczytać wiersza okna "+idOkna+": "+err.Error()))
	}
	if wiersz.KanalModeluID == kanalID {
		return nil
	}
	wiersz.KanalModeluID = kanalID
	if err := a.trwalosc.okna.Aktualizuj(ctx, wiersz); err != nil {
		return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeInternalError,
			"rodzina model.*: nie da się zapisać kanału okna "+idOkna+": "+err.Error()))
	}
	return nil
}

// bladKanaluModelu składa odmowę żądania niezgodnego z kontraktem.
func bladKanaluModelu(powod string) error {
	return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeValidationFailed,
		"rodzina model.*: "+powod))
}
