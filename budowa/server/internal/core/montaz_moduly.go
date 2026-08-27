package core

import (
	"context"

	"danacoconsole/server/internal/dane"
	"danacoconsole/server/internal/injection"
	"danacoconsole/server/internal/konfig"
	"danacoconsole/server/internal/models"
	"danacoconsole/server/internal/mowa"
	"danacoconsole/server/internal/podagenci"
	"danacoconsole/server/internal/session"
)

// adapteryModulow niesie adaptery modułów, które muszą powstać przed rdzeniem:
// dwa uruchamiają procesy, jeden wykonuje zlecenia, jeden przyjmuje odmowy.
type adapteryModulow struct {
	terminal    *adapterTerminala
	developer   *adapterDevelopera
	kolejki     *adapterKolejek
	automatyki  *adapterAutomatyk
	diagnostyka *adapterDiagnostyki
	// mowa uruchamia pomocnika transkrypcji z tego samego uruchamiacza procesów co pozostałe moduły.
	mowa *adapterMowy
	// doradcy prowadzi konsultację silniejszym modelem; stoi na rejestrze kanałów i dzienniku bazy.
	doradcy *adapterDoradcow
}

// zlozAdapteryModulow składa adaptery modułowe wymagające bytów montażu, w kolejności zależności
// między nimi.
func zlozAdapteryModulow(kontekst context.Context, m Montaz, repozytoria *dane.Zestaw,
	nadzorca *session.Nadzorca, rozstrzygacz *konfig.Rozstrzygacz,
	katalogRoboczy *KatalogRoboczy, telemetria *telemetriaPostepu,
	szynaZdarzen Nadajnik, kanaly *models.Rejestr) adapteryModulow {

	// Terminal i Developer stoją na tym samym uruchamiaczu procesów i rozstrzygaczu zasad izolacji.
	terminal := nowyAdapterTerminala(nadzorca.Rejestr(), injection.UruchamiaczOkien()).
		ZTrwaloscia(repozytoria.Terminal).
		ZIzolacja(rozstrzygacz, katalogRoboczy).
		ZKatalogiemDanych(m.Konfiguracja.KatalogDanych).
		ZWyjsciem(szynaZdarzen)
	przygotujTerminal(kontekst, terminal, m.Dziennik)

	developer := nowyAdapterDevelopera(nadzorca.Rejestr(), injection.UruchamiaczOkien()).
		ZTrwaloscia(repozytoria.Developer).
		ZIzolacja(rozstrzygacz, katalogRoboczy).
		ZKanalami(kanaly)
	przygotujDevelopera(kontekst, developer, m.Dziennik)

	// Jeden adapter kolejek służy domenie kolejek, MultitaskingAI i Automations — silnik jest jeden.
	kolejki := nowyAdapterKolejek(repozytoria.Kolejki).
		ZSesjami(repozytoria.Sesje).
		ZTelemetria(telemetria).
		ZWykonawcaModelu(kanaly, szynaZdarzen)

	// Automatyki stoją na tym samym adapterze kolejek, dzieląc port modułu i budzik harmonogramu.
	automatyki := nowyAdapterAutomatyk(repozytoria.Automatyki).ZKolejkami(kolejki).
		ZUkladem(repozytoria.UkladOrkiestracji).ZOknami(repozytoria.Okna).
		ZSejfem(dane.NowySejfPlikowy(m.Konfiguracja.KatalogDanych))

	// Diagnostics pisze dziennik rdzenia do bazy i przyjmuje odmowy wykonania komend przed rdzeniem.
	diagnostyka := nowyAdapterDiagnostyki(kontekst, repozytoria.Diagnostyka)
	diagnostyka.PodepnijDziennik(m.Dziennik)

	// Silnik mowy stoi na uruchamiaczu i źródłach izolacji tych samych co Terminal i Developer.
	mowaModulu := nowyAdapterMowy(injection.UruchamiaczOkien()).
		ZIzolacja(rozstrzygacz, katalogRoboczy)
	if m.Baza != nil {
		mowaModulu = mowaModulu.ZDziennikiem(mowa.NowyDziennik(m.Baza.DB))
	}
	// Magazyn nagrań, katalog danych i konfiguracja domykają cztery komendy obok transkrypcji.
	if repozytoria != nil {
		mowaModulu = mowaModulu.ZMagazynemNagran(repozytoria.NagraniaMowy,
			m.Konfiguracja.KatalogDanych, repozytoria.Konfiguracja)
	}
	mowaModulu = mowaModulu.ZWyjsciem(nowyEmiter(szynaZdarzen))

	// Doradca prowadzi konsultację u modelu równego albo słabszego, na rejestrze kanałów okna rozmowy.
	doradcy := nowyAdapterDoradcow(kanaly).ZOknami(nadzorca.Rejestr()).ZNadajnikiem(szynaZdarzen)
	if m.Baza != nil {
		doradcy = doradcy.ZDziennikiem(podagenci.NowyDziennik(m.Baza.DB))
	}

	return adapteryModulow{
		terminal: terminal, developer: developer,
		kolejki: kolejki, automatyki: automatyki, diagnostyka: diagnostyka,
		mowa: mowaModulu, doradcy: doradcy,
	}
}
