package transport

import (
	"net"
	"strings"
)

// Rdzeń miał cztery zapisane decyzje, z których każda z osobna była rozsądna,
// dopóki nasłuch nie wychodził poza pętlę zwrotną: nasłuch pusty = wszystkie
// interfejsy, uwierzytelnianie celowo nie działa, OriginPatterns "*" bez
// sprawdzania pochodzenia, brak TLS. Trzy z nich są już zdjęte: adres pusty
// znaczy pętlę zwrotną (`ustawienia.go`), pochodzenie jest sprawdzane wykazem
// wzorców (`nawiazanie.go`), a warstwa TLS włącza się parą plików wskazaną
// przełącznikiem albo zmienną środowiska (`konfiguracja/argumenty.go`,
// `konfiguracja/srodowisko.go`).
//
// Czwartą — brak uwierzytelniania — wygasza straż bramki (`bramka.go`): poza
// pętlą zwrotną gniazdo wykonuje wyłącznie `connection.hello`, `auth.login`
// i `auth.register`, dopóki nie przedstawi tokenu sesji. Na pętli zwrotnej straż
// nie istnieje.
//
// Ostrzeżenie zostaje, ale mówi o stanie bieżącym: składa się z tego, co
// naprawdę zastane w nastawach, nie z braków, których już nie ma. Ten plik nie
// zatrzymuje startu — rdzeń ma wstać i działać, tylko ma powiedzieć głośno,
// na czym staje.

// ostrzezenieWystawienia to szkielet linii wypisywanej, gdy nasłuch wychodzi
// poza pętlę zwrotną. Mówi stan faktyczny: co od tej chwili obowiązuje i co
// nadal zostaje na Operatorze. Nie jest to już ostrzeżenie o braku bramki —
// bramka jest — lecz zawiadomienie o zmianie zachowania rdzenia (nic
// istotnego nie dzieje się po cichu).
const ostrzezenieWystawienia = "transport: nasłuch %s wychodzi poza pętlę zwrotną — " +
	"od tej chwili połączenie wykonuje wyłącznie connection.hello, auth.login i auth.register, " +
	"dopóki nie przedstawi tokenu sesji bramki w powitaniu (odmowa: not_authenticated); " +
	"%s; " +
	"na Operatorze zostaje wtedy jedno: siła sekretu bramki, bo po jego zdobyciu " +
	"wołający wykonuje komendy tak samo jak Operator — łącznie z torem zdalnym do jego maszyn."

// czlonWarstwy nazywa stan warstwy transportu w treści zawiadomienia.
//
// Bez TLS zdanie jest mocniejsze. Token sesji jedzie tym łączem przy każdym
// powitaniu, a otwartym tekstem jedzie w postaci czytelnej dla każdego po
// drodze. Warstwa szyfrowana nie rozstrzyga, kto się łączy (od tego jest
// bramka), ale bez niej bramka broni wejścia, którego klucz leci obok, na
// wierzchu.
func czlonWarstwy(zTLS bool) string {
	if zTLS {
		return "połączenie idzie warstwą TLS, więc token sesji nie jedzie na wierzchu"
	}
	return "nie ma TLS — połączenie idzie otwartym tekstem, więc token sesji bramki " +
		"jest czytelny dla każdego pośrednika po drodze; poza siecią zaufaną wskaż parę plików TLS"
}

// opisWystawienia nazywa adres nasłuchu tak, by czytający dziennik widział
// różnicę między adresem wskazanym a brakiem wskazania. Pusty adres to nie
// jest „bez adresu" — to wszystkie interfejsy maszyny, i tak ma być napisane.
func opisWystawienia(adres string) string {
	if strings.TrimSpace(adres) == "" {
		return "na wszystkich interfejsach (adres nasłuchu pusty)"
	}
	return adres
}

// petlaZwrotna rozstrzyga, czy adres nasłuchu jest pętlą zwrotną, czyli czy
// rdzeń jest osiągalny wyłącznie z tej samej maszyny.
//
// Pusty adres to nie pętla zwrotna. net.Listen z pustym hostem wiąże wszystkie
// interfejsy — i IPv4, i IPv6 — więc pusty jest najszerszym z możliwych
// wystawień, nie najwęższym. To samo dotyczy 0.0.0.0 i ::, które ParseIP wyda
// jako adresy nienależące do pętli zwrotnej.
//
// IPv4 i IPv6 idą jedną drogą: net.IP.IsLoopback obejmuje całą sieć 127.0.0.0/8,
// adres ::1 oraz postać ::ffff:127.0.0.1. Zone id (fe80::1%eth0) odcinamy przed
// rozbiorem, bo ParseIP go nie przyjmuje, a bez odcięcia adres łączowy zostałby
// wzięty za nazwę.
//
// Nazwa, która nie jest adresem, jest traktowana jako wystawienie, nie jako
// pętla zwrotna — jedynym wyjątkiem jest nazwa `localhost` i nazwy w jej
// domenie, których rozwiązanie na pętlę zwrotną gwarantuje RFC 6761. Rdzeń nie
// pyta o to systemu nazw: odpytanie DNS przy starcie wstrzymywałoby start,
// a ostrzeżenie ma być tanie i pewne. Wynik z tego jest asymetryczny celowo —
// przy wątpliwości ostrzeżenie pada, bo ostrzeżenie zbędne kosztuje linię
// dziennika, a ostrzeżenie pominięte kosztuje wystawiony rdzeń.
func petlaZwrotna(adres string) bool {
	host := strings.TrimSpace(adres)
	if host == "" {
		return false
	}
	// Postać [::1] z zapisu adres:port — nawiasy zdejmujemy przed rozbiorem.
	host = strings.TrimSuffix(strings.TrimPrefix(host, "["), "]")
	if procent := strings.Index(host, "%"); procent >= 0 {
		host = host[:procent]
	}
	if ip := net.ParseIP(host); ip != nil {
		return ip.IsLoopback()
	}
	nazwa := strings.ToLower(strings.TrimSuffix(host, "."))
	return nazwa == "localhost" || strings.HasSuffix(nazwa, ".localhost")
}

// ostrzezJezeliWystawiony wypisuje ostrzeżenie do dziennika rdzenia, jeżeli
// nasłuch nie stoi na pętli zwrotnej. Nic nie zwraca i nic nie zatrzymuje —
// jedynym skutkiem jest linia w dzienniku (brak zabezpieczenia nie
// wstrzymuje startu; nic istotnego nie dzieje się po cichu).
func ostrzezJezeliWystawiony(u Ustawienia) {
	if u.Dziennik == nil {
		return
	}
	// Zniesienie wymogu przy nasłuchu szerszym jest głośne. Operator ma prawo
	// zdjąć dźwignię — to jego maszyna i jego rozstrzygnięcie — ale nie ma prawa
	// zrobić tego po cichu: nieuwierzytelnione połączenie sięga po hosty z tabeli
	// `host_zdalny`, czyli po VPS i komputer w biurze. Linia idzie zamiast
	// zawiadomienia zwykłego, bo mówi o tym samym nasłuchu rzecz ważniejszą.
	if !petlaZwrotna(u.Adres) && !wymogLogowania(u.Adres, u.WymogLogowania) {
		u.Dziennik.Printf(ostrzezenieZniesienia, opisWystawienia(u.Adres))
		return
	}
	if petlaZwrotna(u.Adres) {
		// Pętla zwrotna z wymogiem WŁĄCZONYM ręcznie też ma ślad — Operator ma
		// wiedzieć, dlaczego jego własna maszyna prosi o logowanie.
		if wymogLogowania(u.Adres, u.WymogLogowania) {
			u.Dziennik.Printf(zawiadomienieWymoguNaPetli, opisWystawienia(u.Adres))
		}
		return
	}
	u.Dziennik.Printf(ostrzezenieWystawienia, opisWystawienia(u.Adres), czlonWarstwy(u.zTLS()))
}

// ostrzezenieZniesienia to jedyna linia tego pliku, która jest ostrzeżeniem
// w pełnym znaczeniu: mówi o rdzeniu stojącym w sieci otworem.
const ostrzezenieZniesienia = "transport: OSTRZEŻENIE — nasłuch %s wychodzi poza pętlę zwrotną, " +
	"a wymóg logowania został ZNIESIONY wskazaniem Operatora; " +
	"znaczy to, że rdzeń wykona komendy każdego, kto dosięgnie tego adresu w sieci, " +
	"tak samo jak komendy Operatora — łącznie z torem zdalnym do jego maszyn (tabela host_zdalny " +
	"wydaje zgodę HOSTOWI, nie wołającemu)."

// zawiadomienieWymoguNaPetli odnotowuje dźwignię podniesioną tam, gdzie sama
// z siebie by nie stanęła.
const zawiadomienieWymoguNaPetli = "transport: nasłuch %s stoi na pętli zwrotnej, " +
	"a wymóg logowania został WŁĄCZONY wskazaniem Operatora — połączenie wykonuje wyłącznie " +
	"connection.hello, auth.login i auth.register, dopóki nie przedstawi tokenu sesji bramki."
