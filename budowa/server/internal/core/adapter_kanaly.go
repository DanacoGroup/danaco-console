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

// adapterKanalow wypełnia port Kanaly rejestrem sterowanym danymi: dopisanie
// wiersza rejestru natychmiast daje działający kanał, bez zmiany w kodzie.
type adapterKanalow struct {
	repozytorium dane.RepozytoriumKanalow
	rejestr      *models.Rejestr
	// sejf obsługuje wyłącznie status poświadczenia, opcjonalnie.
	sejf sejfPoswiadczen
}

// nowyAdapterKanalow wiąże port z repozytorium i rejestrem kanałów, gotowy do
// obsługi całej rodziny komend `channel.*`.
func nowyAdapterKanalow(repozytorium dane.RepozytoriumKanalow, rejestr *models.Rejestr) *adapterKanalow {
	return &adapterKanalow{repozytorium: repozytorium, rejestr: rejestr}
}

// Dodaj dopisuje wiersz rejestru kanałów i odświeża rejestr sterowany danymi,
// żeby nowy kanał zaczął działać bez restartu rdzenia.
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
		KontoDostawcyID:        kontoKanalu(z.Config),
		PoswiadczenieOdwolanie: odwolaniePoswiadczeniaKanalu(z.Config),
		ParametryJSON:          parametryKanalu(z.Config),
		Aktywny:                z.Enabled == nil || *z.Enabled,
	}
	if _, err := a.repozytorium.Dodaj(ctx, kanal); err != nil {
		return shared.ChannelAddResponse{}, bladZapisuKanalu(err, z.Kind)
	}
	a.odswiez(ctx)
	return shared.ChannelAddResponse{Channel: kanalKontraktu(kanal)}, nil
}

// Zmien zmienia wiersz rejestru kanałów wybiórczo, dotykając wyłącznie pól
// obecnych w żądaniu, i odświeża rejestr sterowany danymi.
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
		kanal.KontoDostawcyID = kontoKanalu(z.Config)
	}
	if err := a.repozytorium.Aktualizuj(ctx, kanal); err != nil {
		return shared.ChannelUpdateResponse{}, err
	}
	a.odswiez(ctx)
	return shared.ChannelUpdateResponse{Channel: kanalKontraktu(kanal)}, nil
}

// Usun wykreśla wiersz rejestru kanałów i odświeża rejestr sterowany danymi,
// żeby kanał usunięty przestał być widoczny bez restartu.
func (a *adapterKanalow) Usun(ctx context.Context, z shared.ChannelRemoveRequest) (shared.ChannelRemoveResponse, error) {
	kanal, err := a.repozytorium.PobierzPoKodzie(ctx, z.ChannelId)
	if err != nil {
		return shared.ChannelRemoveResponse{}, err
	}
	if err := a.repozytorium.Usun(ctx, kanal.ID); err != nil {
		return shared.ChannelRemoveResponse{}, bladUsunieciaKanalu(err, z.ChannelId)
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

// przedrostekKanalu znakuje kod wiersza rejestru nadany przez rdzeń, użyty
// przy tworzeniu nowego identyfikatora kanału.
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

// kluczOdwolaniaKanalu to nazwa parametru konfiguracji kanału niosącego
// odwołanie do danych dostępowych, nie samą wartość sekretu.
const kluczOdwolaniaKanalu = "credentialRef"

// odwolaniePoswiadczeniaKanalu wyjmuje z konfiguracji kanału odwołanie do jego
// poświadczenia. Brak parametru albo pusta wartość znaczy kanał bez
// uwierzytelnienia.
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
// z kontem, tą samą drogą co credentialRef, bo kontrakt nie ma pola accountId.
const kluczKontaKanalu = "accountId"

// kontoKanalu wyjmuje z konfiguracji kanału powiązanie z kontem, w postaci
// liczby albo napisu liczbowego. Wartość nieliczbowa znaczy brak powiązania.
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

// brakiWierszaKanalu nazywa pola wiersza kanału, które przyszły puste, zamiast
// zostawiać rozstrzygnięcie usterce wewnętrznej z treścią zapytania SQL.
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

// kanalKontraktu przekłada wiersz repozytorium na kanał kontraktu, wypełniając
// pole modelu tylko wtedy, gdy wiersz niesie identyfikator modelu.
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
