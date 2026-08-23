// Odpowiedzialność pliku: wypełnienie portu Konta rejestrem kont z bazy.
//
// Treść poświadczenia nie ma drogi powrotnej: repozytorium oddaje wyłącznie
// znacznik `MaPoswiadczenie`, a jedyna metoda wynosząca odwołanie
// (`OdwolaniePoswiadczenia`) nie jest tu wywoływana ani razu. Struktura Account
// kontraktu nie ma pola na sekret, więc nie da się go wynieść nawet przez
// pomyłkę.
package core

import (
	"context"
	"strconv"

	"danacoconsole/server/internal/dane"
	"danacoconsole/server/internal/konfiguracja"
	"danacoconsole/shared"
)

// Zgodność adaptera z portem sprawdzana jest przy kompilacji.
var _ Konta = (*adapterKont)(nil)

// adapterKont wypełnia port Konta tabelą `konto`.
type adapterKont struct {
	repozytorium dane.RepozytoriumKont
	sejf         SejfPoswiadczen
}

// nowyAdapterKont wiąże port z rejestrem kont i wpina sejf poświadczeń oparty
// o plik katalogu danych. Bez sejfu poświadczenie z okna kont ginie po cichu,
// a `hasCredential` zawsze mówi „brak".
//
// Katalog obowiązujący wchodzi montażem: `montaz_porty.go` buduje jeden sejf nad
// `Montaz.Konfiguracja.KatalogDanych` — katalogiem, który Operator może
// przestawić przełącznikiem `-dane` albo zmienną `DANACO_KATALOG_DANYCH` — i wpina
// go tu przez `ZSejfem`, tą samą zmienną, którą dostaje magazyn treści biblioteki.
// Sejf domyślny zostaje tu dla wywołania bez montażu, żeby konstruktor nigdy nie
// oddał adaptera bez sejfu.
func nowyAdapterKont(repozytorium dane.RepozytoriumKont) *adapterKont {
	adapter := &adapterKont{repozytorium: repozytorium}
	return adapter.ZSejfem(dane.NowySejfPlikowy(konfiguracja.KatalogDanychDomyslny()))
}

// ZSejfem wpina magazyn sekretów.
func (a *adapterKont) ZSejfem(sejf SejfPoswiadczen) *adapterKont {
	a.sejf = sejf
	return a
}

// Dodaj zakłada konto i — na żądanie — czyni je domyślnym swojego rodzaju.
func (a *adapterKont) Dodaj(ctx context.Context, z shared.AccountAddRequest) (shared.AccountAddResponse, error) {
	if a == nil || a.repozytorium == nil {
		return shared.AccountAddResponse{}, bladBrakuKatalogu("kont")
	}
	konto := dane.Konto{
		Nazwa: z.Name, Rodzaj: z.Kind, Dostawca: z.Provider,
		IdentyfikatorZewnetrzny: z.ExternalId, ModelDomyslny: z.DefaultModel,
		AdresBazowy: z.BaseUrl, KatalogKonfiguracji: z.ConfigDir,
		Stan: dane.StanKontaAktywne, Aktywne: z.Enabled == nil || *z.Enabled,
	}
	odwolanie, err := odwolaniePoswiadczenia(ctx, a.sejf, z.Name, z.Credential)
	if err != nil {
		return shared.AccountAddResponse{}, err
	}
	id, err := a.repozytorium.Dodaj(ctx, konto, odwolanie)
	if err != nil {
		return shared.AccountAddResponse{}, err
	}
	if z.MakeDefault != nil && *z.MakeDefault {
		if _, err := a.repozytorium.UstawDomyslne(ctx, id); err != nil {
			return shared.AccountAddResponse{}, err
		}
	}
	zapisane, err := a.repozytorium.Pobierz(ctx, id)
	if err != nil {
		return shared.AccountAddResponse{}, err
	}
	return shared.AccountAddResponse{Account: kontoKontraktu(zapisane)}, nil
}

// Wykaz zwraca rejestr kont zawężony rodzajem. Poświadczeń wykaz nie niesie.
func (a *adapterKont) Wykaz(ctx context.Context, z shared.AccountListRequest) (shared.AccountListResponse, error) {
	if a == nil || a.repozytorium == nil {
		return shared.AccountListResponse{Accounts: []shared.Account{}}, nil
	}
	wiersze, err := a.repozytorium.Lista(ctx, dane.FiltrKont{
		Rodzaj: z.Kind, TylkoAktywne: z.EnabledOnly != nil && *z.EnabledOnly,
	})
	if err != nil {
		return shared.AccountListResponse{}, err
	}
	konta := make([]shared.Account, 0, len(wiersze))
	for _, wiersz := range wiersze {
		konta = append(konta, kontoKontraktu(wiersz))
	}
	return shared.AccountListResponse{Accounts: konta, Total: len(konta)}, nil
}

// Zmien zapisuje zmienione pola konta. Oznaczenia domyślnego nie rusza — ma
// własną komendę, bo dotyczy całego rodzaju kont, nie jednego wiersza.
func (a *adapterKont) Zmien(ctx context.Context, z shared.AccountUpdateRequest) (shared.AccountUpdateResponse, error) {
	if a == nil || a.repozytorium == nil {
		return shared.AccountUpdateResponse{}, bladBrakuKatalogu("kont")
	}
	konto, err := a.konto(ctx, z.AccountId)
	if err != nil {
		return shared.AccountUpdateResponse{}, bladWskazania(err, "konto", z.AccountId)
	}
	zastosujZmianeKonta(&konto, z)
	if err := a.repozytorium.Aktualizuj(ctx, konto); err != nil {
		return shared.AccountUpdateResponse{}, err
	}
	if err := a.zapiszPoswiadczenie(ctx, konto, z.Credential); err != nil {
		return shared.AccountUpdateResponse{}, err
	}
	zapisane, err := a.repozytorium.Pobierz(ctx, konto.ID)
	if err != nil {
		return shared.AccountUpdateResponse{}, err
	}
	return shared.AccountUpdateResponse{Account: kontoKontraktu(zapisane)}, nil
}

// Usun kasuje konto i oddaje kanały, które utraciły powiązanie. Konto nieznane
// nie jest błędem — wynik mówi wtedy, że nic nie usunięto.
func (a *adapterKont) Usun(ctx context.Context, z shared.AccountRemoveRequest) (shared.AccountRemoveResponse, error) {
	if a == nil || a.repozytorium == nil {
		return shared.AccountRemoveResponse{Removed: false}, nil
	}
	konto, err := a.konto(ctx, z.AccountId)
	if brakWiersza(err) {
		return shared.AccountRemoveResponse{Removed: false}, nil
	}
	if err != nil {
		return shared.AccountRemoveResponse{}, err
	}
	odlaczone, err := a.repozytorium.Usun(ctx, konto.ID)
	if err != nil {
		if brakWiersza(err) {
			return shared.AccountRemoveResponse{Removed: false}, nil
		}
		return shared.AccountRemoveResponse{}, err
	}
	usunPoswiadczenie(ctx, a.sejf, konto.Nazwa)
	return shared.AccountRemoveResponse{Removed: true, DetachedChannelIds: numeryTekstem(odlaczone)}, nil
}

// UstawDomyslne wskazuje konto domyślne swojego rodzaju; domyślne jest dokładnie
// jedno na rodzaj i pilnuje tego indeks bazy, nie warunek w rdzeniu.
func (a *adapterKont) UstawDomyslne(ctx context.Context,
	z shared.AccountDefaultSetRequest) (shared.AccountDefaultSetResponse, error) {

	if a == nil || a.repozytorium == nil {
		return shared.AccountDefaultSetResponse{}, bladBrakuKatalogu("kont")
	}
	konto, err := a.konto(ctx, z.AccountId)
	if err != nil {
		return shared.AccountDefaultSetResponse{}, bladWskazania(err, "konto", z.AccountId)
	}
	poprzednie, err := a.repozytorium.UstawDomyslne(ctx, konto.ID)
	if err != nil {
		return shared.AccountDefaultSetResponse{}, err
	}
	zapisane, err := a.repozytorium.Pobierz(ctx, konto.ID)
	if err != nil {
		return shared.AccountDefaultSetResponse{}, err
	}
	wynik := shared.AccountDefaultSetResponse{Account: kontoKontraktu(zapisane)}
	if poprzednie != nil {
		wynik.PreviousDefaultId = tekstOpcjonalny(strconv.FormatInt(*poprzednie, 10))
	}
	return wynik, nil
}

// konto odczytuje wiersz wskazany identyfikatorem kontraktu. Identyfikator
// niebędący numerem wiersza znaczy dokładnie tyle, co wiersz nieistniejący —
// rozstrzygnięcie, czy to odmowa, czy pustka, należy do wywołującego.
func (a *adapterKont) konto(ctx context.Context, identyfikator string) (dane.Konto, error) {
	id, err := strconv.ParseInt(identyfikator, 10, 64)
	if err != nil {
		return dane.Konto{}, dane.ErrBrakWiersza
	}
	return a.repozytorium.Pobierz(ctx, id)
}

// zapiszPoswiadczenie utrwala odwołanie do nowego sekretu, gdy żądanie je niosło.
func (a *adapterKont) zapiszPoswiadczenie(ctx context.Context, konto dane.Konto, poswiadczenie *string) error {
	odwolanie, err := odwolaniePoswiadczenia(ctx, a.sejf, konto.Nazwa, poswiadczenie)
	if err != nil || odwolanie == nil {
		return err
	}
	return a.repozytorium.UstawPoswiadczenie(ctx, konto.ID, odwolanie)
}
