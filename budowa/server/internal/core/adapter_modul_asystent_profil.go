// Odpowiedzialność pliku: wejście profilu asystenta do tury zlecenia, jako warstwa promptu
// dopisywana do nakładki zapytania, obok metod z adapter_modul_asystent_wykonawca.go.
package core

import (
	"context"
	"strings"
	"time"

	"danacoconsole/server/internal/dane"
	"danacoconsole/server/internal/injection"
	"danacoconsole/server/internal/models"
	"danacoconsole/shared"
)

// profilZlecenia rozstrzyga, który profil obowiązuje w tej turze: profil wskazany żądaniem,
// w braku wskazania profil domyślny, a w braku obu — praca bez warstwy profilu.
func (a *adapterAsystenta) profilZlecenia(ctx context.Context,
	zlecenie dane.ZlecenieAsystenta) (dane.ProfilAsystenta, bool) {

	if a == nil || a.repozytorium == nil {
		return dane.ProfilAsystenta{}, false
	}
	if zlecenie.ProfilKod != nil && strings.TrimSpace(*zlecenie.ProfilKod) != "" {
		profil, err := a.repozytorium.Profil(ctx, *zlecenie.ProfilKod)
		if err != nil {
			return dane.ProfilAsystenta{}, false
		}
		return profil, true
	}
	profil, err := a.repozytorium.ProfilDomyslny(ctx)
	if err != nil {
		return dane.ProfilAsystenta{}, false
	}
	return profil, true
}

// uzupelnijProfil nakłada profil na zapytanie kanału złożone z okna: treść promptu dokłada
// do warstwy profilu nakładki, a nastawy zasięgu tylko wtedy, gdy profil je niesie.
func uzupelnijProfil(profil dane.ProfilAsystenta, zapytanie *models.Zapytanie) {
	if zapytanie == nil {
		return
	}
	if warstwa := warstwaProfilu(profil); warstwa != "" {
		zapytanie.Nakladka.ProfilRoli = zlozWarstwe(zapytanie.Nakladka.ProfilRoli, warstwa)
		// Prompt wbudowany zostaje w mocy, tryb rusza się wyłącznie, gdy profil ma co dopisać.
		zapytanie.Nakladka.Tryb = injection.TrybDopisz
	}
	if profil.KanalModelu != nil && strings.TrimSpace(*profil.KanalModelu) != "" {
		zapytanie.Kanal = strings.TrimSpace(*profil.KanalModelu)
	}
	// Kolumny profilu trzymają wartości kontraktu wprost, więc przekładu tu nie ma.
	if srodowisko := strings.TrimSpace(profil.SrodowiskoWykonania); srodowisko != "" {
		zapytanie.SrodowiskoWykonania = shared.ExecutionEnv(srodowisko)
	}
	if tryb := strings.TrimSpace(profil.TrybUprawnien); tryb != "" {
		zapytanie.TrybUprawnien = shared.PermissionMode(tryb)
	}
}

// wskazanyProfil sprawdza profil wskazany przez żądanie i oddaje go w kształcie gotowym do
// zapisu; sprawdzenie wcześniejsze pozwala nazwać powód odmowy zamiast usterki więzów.
func (a *adapterAsystenta) wskazanyProfil(ctx context.Context, kod *string) (*string, error) {
	if kod == nil {
		return nil, nil
	}
	wskazany := strings.TrimSpace(*kod)
	if wskazany == "" {
		return nil, nil
	}
	if a == nil || a.repozytorium == nil {
		return nil, nil
	}
	if _, err := a.repozytorium.Profil(ctx, wskazany); err != nil {
		return nil, bladWskazaniaAsystenta("wskazany profil " + wskazany + " nie istnieje" +
			"; naprawa: wydać polecenie bez profilu (pracuje wtedy profil domyślny albo żaden)" +
			" albo założyć profil o tym kodzie")
	}
	return &wskazany, nil
}

// opiszBrakWarstwy dopisuje do dziennika zlecenia zdanie o tym, że tura poszła bez warstwy
// promptu profilu; nieudany zapis wpisu nie ma prawa zerwać pracy zlecenia.
func (a *adapterAsystenta) opiszBrakWarstwy(ctx context.Context, zlecenie dane.ZlecenieAsystenta, powod string) {
	if a == nil || a.repozytorium == nil {
		return
	}
	kod := zlecenie.Kod
	_, _ = a.repozytorium.ZapiszWpis(ctx, dane.WpisDziennikaAsystenta{
		Kod:         nowyIdentyfikator(przedrostekWpisuAsystenta),
		OknoKod:     zlecenie.OknoKod,
		ZlecenieKod: &kod,
		Rodzaj:      string(shared.AssistantActivityKindNote),
		Tresc: "Asystent pracuje BEZ warstwy promptu profilu: " + powod +
			". Zlecenie idzie dalej w zakresie, w jakim serwer może je wykonać.",
		Utworzono: time.Now().UnixMilli(),
	})
}

// warstwaProfilu wyjmuje treść warstwy promptu profilu; wskaźnik pusty i napis z samych
// odstępów są tu tym samym — profilem bez treści do dołożenia.
func warstwaProfilu(profil dane.ProfilAsystenta) string {
	if profil.WarstwaPromptu == nil {
		return ""
	}
	return strings.TrimSpace(*profil.WarstwaPromptu)
}
