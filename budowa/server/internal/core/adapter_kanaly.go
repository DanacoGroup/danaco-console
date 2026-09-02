package core

import (
	"context"
	"encoding/json"
	"strconv"
	"strings"

	"danacoconsole/server/internal/dane"
	"danacoconsole/server/internal/models"
	"danacoconsole/server/internal/protocol"
	"danacoconsole/shared"
)

type adapterKanalow struct {
	repozytorium dane.RepozytoriumKanalow
	rejestr      *models.Rejestr
	sejf         sejfPoswiadczen
}

func nowyAdapterKanalow(repozytorium dane.RepozytoriumKanalow, rejestr *models.Rejestr) *adapterKanalow {
	return &adapterKanalow{repozytorium: repozytorium, rejestr: rejestr}
}

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

// Wykaz oddaje kanały konta: rejestr jest jeden na proces (decyzja 34), więc
// jego kontrakt przecina się z kodami zawężonego repozytorium.
func (a *adapterKanalow) Wykaz(ctx context.Context, z shared.ChannelListRequest) (shared.ChannelListResponse, error) {
	tylkoCzynne := z.EnabledOnly != nil && *z.EnabledOnly
	konta, err := kluczeKanalowKonta(ctx, a.repozytorium)
	if err != nil {
		return shared.ChannelListResponse{}, err
	}
	rejestr := a.rejestr.Kontrakt(tylkoCzynne)
	kanaly := make([]shared.Channel, 0, len(rejestr))
	for _, kanal := range rejestr {
		if _, jest := konta[kanal.Id]; jest {
			kanaly = append(kanaly, kanal)
		}
	}
	return shared.ChannelListResponse{Channels: kanaly}, nil
}

func kluczeKanalowKonta(ctx context.Context, repozytorium dane.RepozytoriumKanalow) (map[string]struct{}, error) {
	if repozytorium == nil {
		return nil, protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeInternalError,
			"kanały: serwer nie ma wpiętego rejestru kanałów"))
	}
	wiersze, err := repozytorium.Lista(ctx, false)
	if err != nil {
		return nil, err
	}
	klucze := make(map[string]struct{}, 2*len(wiersze))
	for _, wiersz := range wiersze {
		klucze[wiersz.Kod] = struct{}{}
		klucze[strconv.FormatInt(wiersz.ID, 10)] = struct{}{}
	}
	return klucze, nil
}

// Własność kodu kanału rozstrzyga zawężone repozytorium, nie rejestr procesu
// (decyzja 34); brak wiersza znaczy kanał nieznany temu kontu.
func kanalKonta(ctx context.Context, repozytorium dane.RepozytoriumKanalow,
	rejestr *models.Rejestr, klucz string) (models.Kanal, bool) {

	if rejestr == nil || repozytorium == nil {
		return nil, false
	}
	definicja, jest := rejestr.Definicja(klucz)
	if !jest {
		return nil, false
	}
	if _, err := repozytorium.PobierzPoKodzie(ctx, definicja.Kod); err != nil {
		return nil, false
	}
	return rejestr.Kanal(klucz)
}

func (a *adapterKanalow) odswiez(ctx context.Context) {
	if a.rejestr == nil {
		return
	}
	_ = a.rejestr.Odswiez(ctx)
}

const przedrostekKanalu = "kanal-"

// Kontrakt nie ma pola dostawcy, a kolumna rejestru go wymaga — rdzeń powtarza rodzaj.
func dostawcaKanalu(z shared.ChannelAddRequest) string {
	var parametry map[string]any
	if len(z.Config) > 0 && json.Unmarshal(z.Config, &parametry) == nil {
		if dostawca, jest := parametry["provider"].(string); jest && dostawca != "" {
			return dostawca
		}
	}
	return z.Kind
}

func parametryKanalu(config json.RawMessage) string {
	if len(config) == 0 {
		return "{}"
	}
	return string(config)
}

// Parametr niesie odwołanie do danych dostępowych, nie wartość sekretu.
const kluczOdwolaniaKanalu = "credentialRef"

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

// Kontrakt nie ma pola accountId; powiązanie z kontem idzie parametrem jak credentialRef.
const kluczKontaKanalu = "accountId"

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
