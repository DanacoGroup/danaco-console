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
	// mowa uruchamia pomocnika transkrypcji, więc powstaje tu, razem z dwoma
	// pozostałymi modułami startującymi procesy — z tego samego uruchamiacza.
	mowa *adapterMowy
	// doradcy prowadzi konsultację silniejszym modelem. Powstaje tu, bo stoi na
	// rejestrze kanałów i dzienniku bazy — tak samo jak mowa wyżej.
	doradcy *adapterDoradcow
}

// zlozAdapteryModulow składa adaptery modułowe wymagające bytów montażu.
//
// Zebrane są tu moduły stojące na wspólnych bytach rdzenia — tym samym
// uruchamiaczu procesów, tym samym rozstrzygaczu zasięgu, tym samym ustalaczu
// katalogu roboczego. Drugiego uruchamiacza ani drugiego silnika kolejek nie
// ma nigdzie.
func zlozAdapteryModulow(kontekst context.Context, m Montaz, repozytoria *dane.Zestaw,
	nadzorca *session.Nadzorca, rozstrzygacz *konfig.Rozstrzygacz,
	katalogRoboczy *KatalogRoboczy, telemetria *telemetriaPostepu,
	szynaZdarzen Nadajnik, kanaly *models.Rejestr) adapteryModulow {

	// Terminal i Developer stoją na tym samym uruchamiaczu procesów co okna
	// rozmowy i na tym samym rozstrzygaczu, z którego biorą zasady izolacji
	// egzekwowane przy każdym poleceniu.
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

	// Jeden adapter kolejek dla domeny kolejek i dla modułu Automations —
	// drugiego wykonawcy zleceń nie ma nigdzie. Most do wykonania wpięty jest
	// tu, w silnik: pozycja wchodząca w stan `wykonywana` jedzie turą kanału
	// modelu tym samym rejestrem i nadajnikiem, co tura okna. Ten sam most służy
	// domenie kolejek, MultitaskingAI i Automations, bo silnik jest jeden.
	kolejki := nowyAdapterKolejek(repozytoria.Kolejki).
		ZSesjami(repozytoria.Sesje).
		ZTelemetria(telemetria).
		ZWykonawcaModelu(kanaly, szynaZdarzen)

	// Automatyki stoją na tym samym adapterze kolejek. Instancja powstaje tu,
	// bo dzieli ją port modułu (montaz_porty.go) i budzik harmonogramu: jedna
	// mapa obserwatorów przebiegów dla Execution Monitora.
	//
	// Sejf poświadczeń idzie tym samym katalogiem danych, nad którym stoi sejf
	// kont i punktów dostępu (montaz_porty.go). Wartość poświadczenia kroku
	// leży POZA bazą — bez tego wpięcia `automation.secret.set` odmawia wprost,
	// zamiast zapisywać referencję wskazującą na nic.
	automatyki := nowyAdapterAutomatyk(repozytoria.Automatyki).ZKolejkami(kolejki).
		ZUkladem(repozytoria.UkladOrkiestracji).ZOknami(repozytoria.Okna).
		ZSejfem(dane.NowySejfPlikowy(m.Konfiguracja.KatalogDanych))

	// Diagnostics pisze dziennik rdzenia do bazy i przyjmuje odmowy wykonania
	// komend, więc powstaje przed rdzeniem i zostaje mu podpięty. Bez wywołania
	// `PodepnijDziennik` okno dziennika diagnostyki nie dostałoby ani jednej
	// linii dziennika rdzenia; metoda znosi dziennik pusty sama.
	diagnostyka := nowyAdapterDiagnostyki(kontekst, repozytoria.Diagnostyka)
	diagnostyka.PodepnijDziennik(m.Dziennik)

	// Silnik mowy stoi na tym samym uruchamiaczu i tych samych dwóch źródłach
	// izolacji co Terminal i Developer: pomocnik Pythona jest procesem drzewa
	// jak każdy inny i przechodzi przez tę samą bramę.
	//
	// Dziennik transkrypcji stoi nad bazą montażu. `mowa.NowyDziennik` znosi
	// bazę pustą sam (oddaje wtedy nil), a silnik bez dziennika rozpoznaje mowę
	// tak samo — traci wyłącznie ślad. Baza idzie z `Montaz.Baza`, bo
	// `dane.Zestaw` uchwytu `*sql.DB` nie wystawia; tabelę zakłada
	// `migracja_075_mowa.sql` wraz z czterema kluczami katalogu ustawień.
	mowaModulu := nowyAdapterMowy(injection.UruchamiaczOkien()).
		ZIzolacja(rozstrzygacz, katalogRoboczy)
	if m.Baza != nil {
		mowaModulu = mowaModulu.ZDziennikiem(mowa.NowyDziennik(m.Baza.DB))
	}
	// Magazyn nagrań, katalog danych i konfiguracja domykają cztery komendy
	// dobudowane obok transkrypcji: przyjęcie i oddanie bajtów nagrania
	// (`speech.audio.*`) oraz nastawę wybudzania (`speech.wake.*`). Nadajnik
	// domyka nasłuch ciągły: to nim idą `speech.listen.partial`
	// i `speech.wake.detected`.
	if repozytoria != nil {
		mowaModulu = mowaModulu.ZMagazynemNagran(repozytoria.NagraniaMowy,
			m.Konfiguracja.KatalogDanych, repozytoria.Konfiguracja)
	}
	mowaModulu = mowaModulu.ZWyjsciem(nowyEmiter(szynaZdarzen))

	// Doradca — konsultacja u modelu równego albo słabszego (podagenci/doradca.go).
	// Stoi na tym samym rejestrze kanałów, którym jedzie okno rozmowy, Roundtable
	// i Research — drugiego silnika modelu w rdzeniu nie ma. Dziennik
	// konsultacji bierze `*sql.DB` z montażu, tak samo jak dziennik transkrypcji
	// wyżej, bo `dane.Zestaw` uchwytu bazy nie wystawia; tabelę
	// `konsultacja_doradcy` zakłada `migracja_101_doradca.sql`. Baza pusta znosi
	// się sama: rada zostaje wtedy jawna w strumieniu i traci wyłącznie ślad
	// w bazie.
	//
	// Rejestr okien jest tu częścią sufitu, nie wygodą. Z okna bierze się kanał
	// pytającego, więc bez tego wpięcia komenda `advisor.consult` musiałaby
	// wierzyć modelowi, kim jest — a wtedy sufit siły dałby się obejść jednym
	// polem żądania. Szyna zdarzeń jest z kolei częścią jawności: nią jadą do
	// okna fragmenty rady (`stream.chunk`), bez niej Operator widziałby wyłącznie
	// skrót w zdarzeniu `advisor.consulted`.
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
