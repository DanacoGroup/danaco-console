// Odpowiedzialność pliku: budowniczowie rejestrów i źródeł, z których korzysta
// montaż rdzenia, wszystkie sterowane danymi: nowy wiersz, nie nowa gałąź
// w kodzie. Osobno od montaz.go, bo montaż mówi, co z czym się wiąże.
package core

import (
	"context"
	"log"
	"strconv"
	"strings"
	"time"

	"danacoconsole/server/internal/dane"
	"danacoconsole/server/internal/injection"
	"danacoconsole/server/internal/konfig"
	"danacoconsole/server/internal/models"
	"danacoconsole/shared"
)

func rejestrKanalow(kontekst context.Context, m Montaz, repozytoria *dane.Zestaw,
	przejmowanie *przejmowanieProcesow, zdarzenia *zdarzeniaWykonawcze) *models.Rejestr {
	rejestr := models.NowyRejestr(models.NoweZrodloBazy(m.Baza.DB), models.FabrykiWbudowane())

	fabrykaProcesu := fabrykaKanaluGlownego(pulaKont(kontekst, m, repozytoria), przejmowanie,
		zdarzenia, nazwaKontaZKatalogu(kontekst, repozytoria))
	rejestr.UstawFabryke(models.AdapterCLI, fabrykaProcesu)
	rejestr.UstawFabryke(rodzajKanaluLokalny, fabrykaProcesu)
	if err := rejestr.Odswiez(kontekst); err != nil && m.Dziennik != nil {
		m.Dziennik.Printf("rejestr kanałów: %v", err)
	}
	return rejestr
}

// rodzajKanaluLokalny jest wartością kolumny kanal_modelu.rodzaj_kanalu (KnownChannelKinds).
const rodzajKanaluLokalny = "lokalny"

func nazwaKontaZKatalogu(kontekst context.Context, repozytoria *dane.Zestaw) func(string) string {
	if repozytoria == nil || repozytoria.Konta == nil {
		return nil
	}
	repo := repozytoria.Konta
	return func(wskazanie string) string {
		id, err := strconv.ParseInt(wskazanie, 10, 64)
		if err != nil {
			return wskazanie
		}
		konto, err := repo.Pobierz(kontekst, id)
		if err != nil {
			return wskazanie
		}
		return konto.Nazwa
	}
}

func pulaKont(kontekst context.Context, m Montaz, repozytoria *dane.Zestaw) *injection.PulaKont {
	if repozytoria != nil && repozytoria.Konta != nil {
		konta, err := repozytoria.Konta.KontaRotacji(kontekst, shared.AccountKindCli)
		if err != nil && m.Dziennik != nil {
			m.Dziennik.Printf("katalog kont: %v", err)
		}

		pula := injection.NowaPula(naPuleKont(konta)...)
		wyposazPuleKont(kontekst, pula, m, repozytoria.Konta)
		return pula
	}
	if m.KatalogProfili == "" {
		return injection.NowaPula()
	}
	konta, err := injection.KontaZKatalogu(m.KatalogProfili)
	if err != nil && m.Dziennik != nil {
		m.Dziennik.Printf("profile kanału głównego: %v", err)
	}
	return injection.NowaPula(konta...)
}

func wyposazPuleKont(kontekst context.Context, pula *injection.PulaKont, m Montaz,
	repo dane.RepozytoriumKont) {

	pula.UstawUtrwalanie(func(kod string, doChwili time.Time) {
		konto, err := repo.PobierzPoNazwie(kontekst, kod)
		if err != nil {
			if m.Dziennik != nil {
				m.Dziennik.Printf("utrwalenie wyczerpania konta %q: %v", kod, err)
			}
			return
		}
		iso := doChwili.UTC().Format(time.RFC3339)
		if err := repo.OznaczStan(kontekst, konto.ID, dane.StanKontaWyczerpane, &iso); err != nil && m.Dziennik != nil {
			m.Dziennik.Printf("utrwalenie wyczerpania konta %q: %v", kod, err)
		}
	})

	pula.UstawZrodlo(func() ([]injection.Konto, bool) {
		konta, err := repo.KontaRotacji(kontekst, shared.AccountKindCli)
		if err != nil {
			if m.Dziennik != nil {
				m.Dziennik.Printf("odświeżenie puli kont: %v", err)
			}
			return nil, false
		}
		return naPuleKont(konta), true
	})
}

func naPuleKont(konta []dane.Konto) []injection.Konto {
	pula := make([]injection.Konto, 0, len(konta))
	for _, konto := range konta {
		if konto.KatalogKonfiguracji == nil || *konto.KatalogKonfiguracji == "" {
			continue
		}
		wpis := injection.Konto{
			Kod:                 konto.Nazwa,
			KatalogKonfiguracji: *konto.KatalogKonfiguracji,
		}
		if konto.Stan == dane.StanKontaWyczerpane && konto.WyczerpaneDo != nil {
			if chwila, err := time.Parse(time.RFC3339, strings.TrimSpace(*konto.WyczerpaneDo)); err == nil {
				wpis.WyczerpaneDo = chwila
			}
		}
		pula = append(pula, wpis)
	}
	return pula
}

// Adres złożony poziom×oś od migracji 012; błąd odczytu nie zatrzymuje rozstrzygania,
// konto idzie z kontekstu rozstrzygania (ustawienie.konto_id od migracji 484).
func zrodloUstawienOsiZBazy(kontekst context.Context,
	repozytorium dane.RepozytoriumKonfiguracjiOsi) konfig.Zrodlo {

	return konfig.NoweZrodloZOdczytuOsi(func(konto int64, adres konfig.Adres) ([]konfig.Wpis, error) {
		ustawienia, err := repozytorium.ListaOsi(dane.ZKontemOperatora(kontekst, konto),
			adres.Poziom, adres.KluczZasiegu, adres.Os, adres.KluczOsi)
		if err != nil {
			return nil, err
		}
		wpisy := make([]konfig.Wpis, 0, len(ustawienia))
		for _, u := range ustawienia {
			wpisy = append(wpisy, konfig.Wpis{
				Poziom: u.Poziom, KluczZasiegu: u.KluczZasiegu,
				Os: u.Os, KluczOsi: u.KluczOsi, Klucz: u.Klucz,
				Wartosc: wartoscTekstu(u.Wartosc), Rodzaj: konfig.Rodzaj(u.RodzajWartosci),
			})
		}
		return wpisy, nil
	})
}

type katalogUstawienZBazy struct {
	kontekst     context.Context
	repozytorium dane.RepozytoriumKatalogUstawien
}

func (k katalogUstawienZBazy) Definicje() ([]shared.SettingDefinition, error) {
	if k.repozytorium == nil {
		return nil, nil
	}
	return k.repozytorium.Definicje(k.kontekst, true)
}

func rejestrUstawien(kontekst context.Context, repozytoria *dane.Zestaw,
	dziennik *log.Logger) *konfig.Rejestr {

	rejestr, zKatalogu := konfig.RejestrZKatalogu(katalogUstawienZBazy{kontekst, repozytoria.KatalogUstawien})
	if zKatalogu {
		return rejestr
	}
	if dziennik != nil {
		dziennik.Printf("katalog ustawień: pusty albo niedostępny, obowiązuje rejestr wbudowany")
	}
	rejestr.Dodaj(DefinicjeKataloguRoboczego()...)
	return rejestr
}

func rejestrAkcji(kontekst context.Context, repozytoria *dane.Zestaw, dziennik *log.Logger) *RejestrAkcji {
	rejestr := NowyRejestrAkcji(ZrodloAkcjiZRepozytorium(repozytoria.Akcje))
	if err := rejestr.Odswiez(kontekst); err != nil && dziennik != nil {
		dziennik.Printf("katalog akcji: %v", err)
	}
	return rejestr
}
