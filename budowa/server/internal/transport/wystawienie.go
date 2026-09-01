package transport

import (
	"net"
	"strings"
)

// Rdzeń miał cztery decyzje bezpieczeństwa nierozsądne poza pętlą zwrotną; trzy są już zdjęte.

// ostrzezenieWystawienia to szkielet linii wypisywanej, gdy nasłuch wychodzi poza pętlę zwrotną maszyny.
const ostrzezenieWystawienia = "transport: nasłuch %s wychodzi poza pętlę zwrotną — " +
	"od tej chwili połączenie wykonuje wyłącznie connection.hello, auth.login i auth.register, " +
	"dopóki nie przedstawi tokenu sesji bramki w powitaniu (odmowa: not_authenticated); " +
	"%s; " +
	"na Operatorze zostaje wtedy jedno: siła sekretu bramki, bo po jego zdobyciu " +
	"wołający wykonuje komendy tak samo jak Operator — łącznie z torem zdalnym do jego maszyn."

// Funkcja czlonWarstwy nazywa stan warstwy transportu, szyfrowanej albo otwartej, w treści zawiadomienia.
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

// Funkcja petlaZwrotna rozstrzyga, czy adres nasłuchu jest pętlą zwrotną osiągalną wyłącznie z tej samej maszyny.
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

// Funkcja ostrzezJezeliWystawiony wypisuje do dziennika rdzenia stan bramki dla tego nasłuchu; każdy z czterech stanów ma własny wiersz.
func ostrzezJezeliWystawiony(u Ustawienia) {
	if u.Dziennik == nil {
		return
	}
	// Nastawa zdejmująca wymóg nie ma skutku, a Operator ma to przeczytać, nie
	// wywnioskować z zachowania rdzenia.
	if u.WymogLogowania != nil && !*u.WymogLogowania {
		u.Dziennik.Printf(zawiadomienieNastawyOdrzuconej, ZmiennaZniesieniaBramki)
	}
	naPetli := petlaZwrotna(u.Adres)
	wymagana := u.dopuszczenie().wymagana
	switch {
	case !naPetli && !wymagana:
		// Zniesienie wymogu przy nasłuchu szerszym jest głośne; Operator zdejmuje dźwignię, nie po cichu.
		u.Dziennik.Printf(ostrzezenieZniesienia, opisWystawienia(u.Adres), ZmiennaZniesieniaBramki)
	case naPetli && wymagana:
		// Pętla zwrotna z wymogiem włączonym ręcznie też ma ślad w dzienniku dla Operatora.
		u.Dziennik.Printf(zawiadomienieWymoguNaPetli, opisWystawienia(u.Adres))
	case naPetli && !wymagana:
		// Przypadek najczęstszy i najcichszy: bramka zdjęta z samego adresu nasłuchu.
		u.Dziennik.Printf(zawiadomieniePetliBezWymogu, opisWystawienia(u.Adres), ZmiennaZniesieniaBramki)
	default:
		u.Dziennik.Printf(ostrzezenieWystawienia, opisWystawienia(u.Adres), czlonWarstwy(u.zTLS()))
	}
}

// ostrzezenieZniesienia to jedyna linia tego pliku, która jest ostrzeżeniem
// w pełnym znaczeniu: mówi o rdzeniu stojącym w sieci otworem.
const ostrzezenieZniesienia = "transport: OSTRZEŻENIE — nasłuch %s wychodzi poza pętlę zwrotną, " +
	"a wymóg logowania został ZNIESIONY zmienną %s wskazaną przy starcie procesu; " +
	"znaczy to, że serwer wykona komendy każdego, kto dosięgnie tego adresu w sieci, " +
	"tak samo jak komendy Operatora — łącznie z torem zdalnym do jego maszyn (tabela host_zdalny " +
	"wydaje zgodę HOSTOWI, nie wołającemu)."

// zawiadomienieWymoguNaPetli odnotowuje dźwignię wymogu podniesioną tam, gdzie sama z siebie by nie stanęła.
const zawiadomienieWymoguNaPetli = "transport: nasłuch %s stoi na pętli zwrotnej, " +
	"a wymóg logowania został WŁĄCZONY wskazaniem Operatora — połączenie wykonuje wyłącznie " +
	"connection.hello, auth.login i auth.register, dopóki nie przedstawi tokenu sesji bramki."

// zawiadomieniePetliBezWymogu nazywa stan, w którym rdzeń wykonuje komendy
// każdego procesu maszyny. Bez tego wiersza przypadek najgroźniejszy z lokalnych
// jest jedynym, który nie zostawia w dzienniku żadnego śladu.
const zawiadomieniePetliBezWymogu = "transport: nasłuch %s stoi na pętli zwrotnej BEZ wymogu logowania — " +
	"każdy proces tej maszyny wykonuje komendy tak samo jak Operator, łącznie z torem zdalnym do jego maszyn; " +
	"wymóg podnosi nastawa gateway.requireLogin=true, a zdejmuje wyłącznie zmienna %s przy starcie procesu."

// zawiadomienieNastawyOdrzuconej odnotowuje nastawę, która nie ma skutku:
// wymóg logowania wolno nastawą wyłącznie podnieść.
const zawiadomienieNastawyOdrzuconej = "transport: nastawa gateway.requireLogin=false pominięta — " +
	"nastawę zmienia komenda config.set wykonywana z gniazda, więc wymóg logowania wolno nią wyłącznie podnieść; " +
	"zdjęcie bramki bierze się wyłącznie ze zmiennej %s wskazanej przy starcie procesu."
