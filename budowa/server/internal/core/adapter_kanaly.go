package core

import (
	"context"
	"encoding/json"
	"strconv"
	"strings"

	"danacoconsole/server/internal/dane"
	"danacoconsole/server/internal/models"
	"danacoconsole/shared"
)

// adapterKanalow wypełnia port Kanaly rejestrem sterowanym danymi.
//
// Zapis idzie do tabeli rejestru, odczyt do rejestru kanałów zbudowanego z tej
// samej tabeli. Po każdej zmianie rejestr jest odświeżany, więc dopisanie
// wiersza natychmiast daje działający kanał — bez zmiany w kodzie i bez restartu.
//
// Identyfikatorem kanału w kontrakcie jest kod wiersza, nie numer wiersza:
// kod przeżywa przeniesienie bazy i jest tym, co widzi okno komunikacji.
type adapterKanalow struct {
	repozytorium dane.RepozytoriumKanalow
	rejestr      *models.Rejestr
	// sejf obsługuje wyłącznie `channel.credential.status`
	// (`adapter_kanaly_sprawdzenie.go`) i widzi z sejfu tylko odczyt. Zależność
	// opcjonalna: bez niej stan poświadczenia mówi „nieustawione" wraz
	// z odwołaniem, pod którym rdzeń szukał.
	sejf sejfPoswiadczen
}

// nowyAdapterKanalow wiąże port z repozytorium i rejestrem kanałów.
func nowyAdapterKanalow(repozytorium dane.RepozytoriumKanalow, rejestr *models.Rejestr) *adapterKanalow {
	return &adapterKanalow{repozytorium: repozytorium, rejestr: rejestr}
}

// Dodaj dopisuje wiersz rejestru kanałów.
func (a *adapterKanalow) Dodaj(ctx context.Context, z shared.ChannelAddRequest) (shared.ChannelAddResponse, error) {
	if brak := brakiWierszaKanalu(z); brak != "" {
		return shared.ChannelAddResponse{}, bladWskazaniaKanalu(brak)
	}
	kanal := dane.Kanal{
		Kod:                    nowyIdentyfikator(przedrostekKanalu),
		Nazwa:                  z.Name,
		Dostawca:               dostawcaKanalu(z),
		IdentyfikatorModelu:    wartoscTekstu(z.Model),
		RodzajKanalu:           z.Kind,
		KontoID:                kontoKanalu(z.Config),
		PoswiadczenieOdwolanie: odwolaniePoswiadczeniaKanalu(z.Config),
		ParametryJSON:          parametryKanalu(z.Config),
		Aktywny:                z.Enabled == nil || *z.Enabled,
	}
	if _, err := a.repozytorium.Dodaj(ctx, kanal); err != nil {
		return shared.ChannelAddResponse{}, err
	}
	a.odswiez(ctx)
	return shared.ChannelAddResponse{Channel: kanalKontraktu(kanal)}, nil
}

// Zmien zmienia wiersz rejestru kanałów wybiórczo.
func (a *adapterKanalow) Zmien(ctx context.Context, z shared.ChannelUpdateRequest) (shared.ChannelUpdateResponse, error) {
	if strings.TrimSpace(z.ChannelId) == "" {
		return shared.ChannelUpdateResponse{},
			bladWskazaniaKanalu("zmiana kanału z pustym polem channelId")
	}
	kanal, err := a.repozytorium.PobierzPoKodzie(ctx, z.ChannelId)
	if err != nil {
		return shared.ChannelUpdateResponse{}, err
	}
	if z.Name != nil {
		kanal.Nazwa = *z.Name
	}
	if z.Model != nil {
		kanal.IdentyfikatorModelu = *z.Model
	}
	if z.Enabled != nil {
		kanal.Aktywny = *z.Enabled
	}
	if len(z.Config) > 0 {
		kanal.ParametryJSON = parametryKanalu(z.Config)
		kanal.PoswiadczenieOdwolanie = odwolaniePoswiadczeniaKanalu(z.Config)
		kanal.KontoID = kontoKanalu(z.Config)
	}
	if err := a.repozytorium.Aktualizuj(ctx, kanal); err != nil {
		return shared.ChannelUpdateResponse{}, err
	}
	a.odswiez(ctx)
	return shared.ChannelUpdateResponse{Channel: kanalKontraktu(kanal)}, nil
}

// Usun wykreśla wiersz rejestru kanałów.
func (a *adapterKanalow) Usun(ctx context.Context, z shared.ChannelRemoveRequest) (shared.ChannelRemoveResponse, error) {
	kanal, err := a.repozytorium.PobierzPoKodzie(ctx, z.ChannelId)
	if err != nil {
		return shared.ChannelRemoveResponse{}, err
	}
	if err := a.repozytorium.Usun(ctx, kanal.ID); err != nil {
		return shared.ChannelRemoveResponse{}, err
	}
	a.odswiez(ctx)
	return shared.ChannelRemoveResponse{ChannelId: z.ChannelId}, nil
}

// Wykaz zwraca rejestr kanałów tak, jak widzi go warstwa modeli — czyli po
// zbudowaniu adapterów. Kanał z wiersza, dla którego nie ma adaptera, nie
// pokazuje się jako gotowy do pracy.
func (a *adapterKanalow) Wykaz(_ context.Context, z shared.ChannelListRequest) (shared.ChannelListResponse, error) {
	tylkoCzynne := z.EnabledOnly != nil && *z.EnabledOnly
	return shared.ChannelListResponse{Channels: a.rejestr.Kontrakt(tylkoCzynne)}, nil
}

// odswiez przebudowuje rejestr kanałów po zmianie wiersza. Niepowodzenie
// odświeżenia nie unieważnia zapisu — wiersz jest w bazie, a rejestr odbuduje
// się przy kolejnej zmianie albo starcie.
func (a *adapterKanalow) odswiez(ctx context.Context) {
	if a.rejestr == nil {
		return
	}
	_ = a.rejestr.Odswiez(ctx)
}

// przedrostekKanalu znakuje kod wiersza rejestru nadany przez rdzeń.
const przedrostekKanalu = "kanal-"

// dostawcaKanalu wybiera dostawcę: parametr „provider" wiersza, a w jego braku
// rodzaj kanału. Kontrakt nie ma osobnego pola dostawcy, a kolumna rejestru go
// wymaga — rdzeń nie zmyśla wartości, tylko powtarza rodzaj.
func dostawcaKanalu(z shared.ChannelAddRequest) string {
	var parametry map[string]any
	if len(z.Config) > 0 && json.Unmarshal(z.Config, &parametry) == nil {
		if dostawca, jest := parametry["provider"].(string); jest && dostawca != "" {
			return dostawca
		}
	}
	return z.Kind
}

// parametryKanalu zwraca parametry wiersza. Parametry niosą wyłącznie odwołania
// do danych dostępowych, nigdy ich treść.
func parametryKanalu(config json.RawMessage) string {
	if len(config) == 0 {
		return "{}"
	}
	return string(config)
}

// kluczOdwolaniaKanalu to nazwa parametru konfiguracji kanału niosącego odwołanie
// do danych dostępowych — nazwę zmiennej środowiskowej, pod którą Operator trzyma
// klucz kanału API. Sama nazwa nie jest sekretem, więc jedzie w parametrach;
// wartość mieszka poza bazą, a kanał API czyta ją przy wysyłce.
const kluczOdwolaniaKanalu = "credentialRef"

// odwolaniePoswiadczeniaKanalu wyjmuje z konfiguracji kanału odwołanie do jego
// poświadczenia i podaje je do kolumny poswiadczenie_odwolanie. Bez tej drogi
// kolumna zostaje pusta i kanał API nie ma skąd wziąć nazwy zmiennej z kluczem.
// Brak parametru albo pusta wartość znaczy kanał bez uwierzytelnienia i daje
// brak odwołania.
//
// Kontrakt nie ma osobnego pola na to odwołanie, więc jedzie ono parametrem
// konfiguracji.
func odwolaniePoswiadczeniaKanalu(config json.RawMessage) *string {
	if len(config) == 0 {
		return nil
	}
	var parametry map[string]any
	if json.Unmarshal(config, &parametry) != nil {
		return nil
	}
	odwolanie, jest := parametry[kluczOdwolaniaKanalu].(string)
	if !jest {
		return nil
	}
	if strings.TrimSpace(odwolanie) == "" {
		return nil
	}
	return &odwolanie
}

// kluczKontaKanalu to nazwa parametru konfiguracji kanału niosącego powiązanie
// z kontem — identyfikator wiersza rejestru kont (kanal_modelu.konto_id).
// Kontrakt ChannelAdd/Update nie ma pola accountId, więc powiązanie jedzie
// parametrem konfiguracji, tą samą drogą co credentialRef.
const kluczKontaKanalu = "accountId"

// kontoKanalu wyjmuje z konfiguracji kanału powiązanie z kontem i podaje je do
// kolumny konto_id. Wartość jest identyfikatorem wiersza rejestru kont; przyjmuje
// postać liczby albo napisu liczbowego (JSON koduje liczby jako float64). Brak
// parametru, wartość pusta albo nieliczbowa znaczy kanał bez powiązania i daje
// brak konta — dana pomocnicza nie może wywrócić zapisu.
func kontoKanalu(config json.RawMessage) *int64 {
	if len(config) == 0 {
		return nil
	}
	var parametry map[string]any
	if json.Unmarshal(config, &parametry) != nil {
		return nil
	}
	wartosc, jest := parametry[kluczKontaKanalu]
	if !jest || wartosc == nil {
		return nil
	}
	switch typowa := wartosc.(type) {
	case float64:
		id := int64(typowa)
		return &id
	case string:
		id, err := strconv.ParseInt(strings.TrimSpace(typowa), 10, 64)
		if err != nil || id == 0 {
			return nil
		}
		return &id
	default:
		return nil
	}
}

// brakiWierszaKanalu nazywa pola wiersza kanału, które przyszły puste. Wiersz
// bez rodzaju kanału nie przechodził dotąd więzu schematu i wracał jako usterka
// wewnętrzna z treścią zapytania SQL — Operator dostawał nazwę kolumny bazy
// zamiast nazwy pola, którego nie wypełnił. Brak samego pola w treści żądania
// odsiewa brama kontraktu (`brama_kontraktu.go`); tutaj rozstrzyga się wartość
// pusta.
func brakiWierszaKanalu(z shared.ChannelAddRequest) string {
	var puste []string
	if strings.TrimSpace(z.Name) == "" {
		puste = append(puste, "name")
	}
	if strings.TrimSpace(z.Kind) == "" {
		puste = append(puste, "kind")
	}
	if len(puste) == 0 {
		return ""
	}
	return "wiersz rejestru kanałów z pustymi polami: " + strings.Join(puste, ", ")
}

// kanalKontraktu przekłada wiersz repozytorium na kanał kontraktu.
func kanalKontraktu(k dane.Kanal) shared.Channel {
	kanal := shared.Channel{
		Id: k.Kod, Name: k.Nazwa, Kind: k.RodzajKanalu, Enabled: k.Aktywny,
		Config: json.RawMessage(k.ParametryJSON),
	}
	if k.IdentyfikatorModelu != "" {
		model := k.IdentyfikatorModelu
		kanal.Model = &model
	}
	return kanal
}
