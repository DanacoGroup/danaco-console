// Odpowiedzialność pliku: wejście profilu asystenta do tury zlecenia. Typ
// `adapterAsystenta` deklaruje `adapter_modul_asystent.go`, turę prowadzi
// `adapter_modul_asystent_wykonawca.go`; ten plik dokłada metody na tym samym
// typie.
//
// Profil niesie warstwę promptu, która mówi modelowi, że jest klawiaturą
// Operatora, a nie autorem. To ona rozstrzyga, czy model sięgnie po narzędzia
// platformy (`most_narzedzi.go`), czy odpisze tekstem we własnym oknie. Bez niej
// asystent zachowuje się jak zwykły czat, dlatego warstwa ma byt trwały, a nie
// żyje w polu żądania, które ginie razem z odpowiedzią na komendę.
//
// Warstwa profilu ląduje w `models.Zapytanie.Nakladka` — dokładnie tam, gdzie
// rozmowa kładzie warstwy osi (`adapter_rozmowa_tozsamosc.go`) i warstwy
// eksperta (`tozsamosc_agenta_nakladka.go`). Dalej jedzie tą samą trasą:
// `Nakladka` → `PromptSystemowy()` → nakładka silnika → przełącznik promptu
// kanału, a `models.ProwenancjaZapytania` pokazuje ją w polu `systemPrompt`
// fragmentu `provenance`. Druga droga do promptu znaczyłaby dwie prawdy o tym,
// co model naprawdę dostał.
//
// Profil dopisuje, nigdy nie zastępuje — tak samo jak ekspert. Prompt wbudowany
// programu kanału zostaje w mocy, a profil dokłada do niego paczkę. Stąd tryb
// nakładki ustawiony w jedną stronę (`injection.TrybDopisz`); profil nie ma
// żadnej drogi, którą mógłby prompt platformy zdjąć.
//
// Syntezy mowy tu nie ma. Kolumna `glos_syntezy` niesie nastawę głosu odczytu,
// ale rdzeń syntezy nie wykonuje — nastawa jest odkładana i czytana, a nie
// udawana wywołaniem, którego nikt nie spełnia.
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

// profilZlecenia rozstrzyga, który profil obowiązuje w tej turze.
//
// Kolejność pytań jest rozstrzygnięciem, nie wygodą:
//
//  1. profil wskazany przez żądanie (`profileId`, odłożony w kolumnie
//     `zlecenie_asystenta.profil_kod`) — Operator powiedział wprost, którym
//     profilem pracuje;
//  2. profil domyślny (`profil_asystenta.domyslny = 1`) — bo okno asystenta ma
//     ruszać bez wybierania profilu przy każdym poleceniu;
//  3. brak — i to nie jest usterka.
//
// Brak profilu nie wstrzymuje pracy. Instalacja świeża nie ma ani jednego
// wiersza w `profil_asystenta`, bo komendy zakładającej profil kontrakt nie ma.
// Gdyby brak warstwy odmawiał tury, asystent nie ruszyłby na takiej instalacji
// w ogóle. Zlecenie idzie więc bez warstwy, a fakt zostaje nazwany wpisem
// dziennika (`opiszBrakWarstwy`), zamiast zniknąć.
//
// Wskazanie na profil, którego nie ma, kończy się tak samo jak brak wskazania:
// bez warstwy i z nazwanym powodem. Zerwanie zlecenia byłoby tu karą za nastawę
// okna, na którą Operator w tej komendzie nie ma wpływu.
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

// uzupelnijProfil nakłada profil na zapytanie kanału złożone z okna.
//
// Warstwa promptu idzie do warstwy `profil` nakładki, a nie do konstytucji ani
// do ekspertyzy. „Jesteś klawiaturą Operatora, nie autorem" opisuje rolę,
// w której model ma wystąpić; konstytucja jest warstwą platformy (nie profilu
// jednego okna), a ekspertyza — wiedzą zadaniową. Włożenie tego zdania
// w konstytucję postawiłoby nastawę jednego okna wyżej niż zasady całej
// platformy.
//
// `zlozWarstwe` dokleja treść po tym, co w warstwie już stoi. Dziś zlecenie
// asystenta wchodzi tu z nakładką pustą (tura zlecenia nie liczy osi tożsamości),
// więc profil jest w niej pierwszy i jedyny — ale dokładanie, a nie podstawianie,
// znaczy, że dołożenie osi w przyszłości nie skasuje warstwy profilu po cichu.
//
// Nastawy poza promptem: kanał, zasięg urządzenia i zasięg pracy. Wszystkie
// trzy są nastawą zasięgu, nie bramką — profil mówi, dokąd sięga praca, a nie
// czego zabrania. Nakładane są tylko wtedy, gdy profil naprawdę coś w danej
// sprawie mówi: pole puste znaczy „profil nie ma zdania" i wtedy zostaje
// nastawa okna. Bez tego strażnika profil bez wskazanego kanału zabrałby turze
// jedyny kanał, jaki miała.
func uzupelnijProfil(profil dane.ProfilAsystenta, zapytanie *models.Zapytanie) {
	if zapytanie == nil {
		return
	}
	if warstwa := warstwaProfilu(profil); warstwa != "" {
		zapytanie.Nakladka.ProfilRoli = zlozWarstwe(zapytanie.Nakladka.ProfilRoli, warstwa)
		// Prompt wbudowany zostaje w mocy. Tryb rusza się wyłącznie wtedy, gdy
		// profil ma co dopisać — profil bez warstwy nie ma prawa zmieniać sposobu
		// podania promptu, bo nie wnosi ani zdania (spójnie z `wnosiTresc`).
		zapytanie.Nakladka.Tryb = injection.TrybDopisz
	}
	if profil.KanalModelu != nil && strings.TrimSpace(*profil.KanalModelu) != "" {
		zapytanie.Kanal = strings.TrimSpace(*profil.KanalModelu)
	}
	// Kolumny profilu trzymają wartości kontraktu wprost (pilnuje tego warunek
	// CHECK tabeli), więc przekładu tu nie ma — jest rzutowanie na typ
	// kontraktu. Gdyby kolumna niosła napis polski, jak
	// `okno_komunikacji.tryb_uprawnien`, musiałby tu stać słownik
	// `shared.WartosciKontraktuPermissionMode`.
	if srodowisko := strings.TrimSpace(profil.SrodowiskoWykonania); srodowisko != "" {
		zapytanie.SrodowiskoWykonania = shared.ExecutionEnv(srodowisko)
	}
	if tryb := strings.TrimSpace(profil.TrybUprawnien); tryb != "" {
		zapytanie.TrybUprawnien = shared.PermissionMode(tryb)
	}
}

// wskazanyProfil sprawdza profil wskazany przez żądanie i oddaje go w kształcie
// gotowym do odłożenia w kolumnie `zlecenie_asystenta.profil_kod`.
//
// Sprawdzenie jest konieczne: kolumna ma klucz obcy do `profil_asystenta`,
// a pragma `foreign_keys` jest na połączeniu włączona (`store/baza.go`) — kod
// nieznany wywróciłby cały zapis polecenia usterką więzów, czyli
// `internal_error` bez powodu czytelnego dla Operatora. Sprawdzamy wcześniej,
// żeby powód nazwać.
//
// To odmowa, a nie ciche pominięcie: brak profilu w ogóle pracy nie wstrzymuje,
// ale wskazanie profilu, którego nie ma, jest czym innym — Operator powiedział
// wprost, którą warstwą promptu ma pracować model. Ciche pominięcie puściłoby
// turę z warstwą profilu domyślnego albo bez żadnej, meldując wykonanie
// polecenia, którego nikt nie wydał.
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

// opiszBrakWarstwy dopisuje do dziennika zlecenia zdanie o tym, że tura poszła
// bez warstwy promptu profilu.
//
// Kontrakt nie ma na to ani jednego pola: `AssistantVoiceCommandResponse`
// niesie transkrypcję i zlecenie, a `AssistantAction.result` jest miejscem na
// odpowiedź modelu — dopisanie tam zdania rdzenia zmieszałoby dwa głosy
// w jednym polu i zafałszowało wynik tury.
// Dziennik (`assistant.activity.list`) jest tą samą rozmową, w której stoi
// polecenie i wynik, więc Operator czyta powód dokładnie tam, gdzie patrzy —
// i zostaje mu ślad po zleceniu, a nie sam komunikat, który znika.
//
// Wpis nie jest warunkiem tury: nieudany zapis nie ma prawa zerwać pracy, którą
// zlecenie ma wykonać. Stąd wynik zapisu porzucony rozmyślnie.
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
			". Zlecenie idzie dalej w zakresie, w jakim rdzeń może je wykonać.",
		Utworzono: time.Now().UnixMilli(),
	})
}

// warstwaProfilu wyjmuje treść warstwy promptu profilu. Wskaźnik pusty i napis
// z samych odstępów są tu tym samym: profilem, który nie ma nic do powiedzenia.
func warstwaProfilu(profil dane.ProfilAsystenta) string {
	if profil.WarstwaPromptu == nil {
		return ""
	}
	return strings.TrimSpace(*profil.WarstwaPromptu)
}
