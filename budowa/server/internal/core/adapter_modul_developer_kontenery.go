// Pakiet obsługuje zakładkę Containers w Dev Tools: wykaz, czynności, log,
// budowanie obrazu i kompozycję kontenerów biblioteką klienta OCI
// wkompilowaną w rdzeń, bez uruchamiania programu `docker`.
package core

import (
	"archive/tar"
	"context"
	"encoding/json"
	"errors"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/docker/docker/api/types/build"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"go.yaml.in/yaml/v3"

	"danacoconsole/shared"
)

const (
	// czasSilnikaKontenerow jest granicą jednej rozmowy z silnikiem kontenerów,
	// wspólną wszystkim czynnościom poza budowaniem obrazu.
	czasSilnikaKontenerow = 60 * time.Second
	// czasBudowaniaObrazu jest granicą budowania obrazu. Budowanie bywa długie,
	// lecz nie nieskończone — obraz budujący się godzinę trzyma połączenie okna.
	czasBudowaniaObrazu = 30 * time.Minute
	// najwiecejWierszyLoguKontenera jest domyślną głębokością logu kontenera,
	// gdy żądanie nie poda własnej głębokości.
	najwiecejWierszyLoguKontenera = 500
)

// silnikKontenerow otwiera rozmowę z silnikiem kontenerów.
//
// Wskazanie gniazda bierze się ze środowiska serwera (`DOCKER_HOST`), a przy
// jego braku — z miejsca domyślnego. Dzięki temu Podman podpina się nastawą
// serwera, bez zmiany w kodzie.
func silnikKontenerow() (*client.Client, error) {
	silnik, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, err
	}
	return silnik, nil
}

// brakSilnikaKontenerow rozpoznaje odmowę „silnika nie ma”.
//
// Odróżnienie jest istotne, bo naprawa jest inna: brak silnika usuwa się
// instalacją po stronie serwera, a odmowa silnika — poprawką żądania.
func brakSilnikaKontenerow(err error) bool {
	if err == nil {
		return false
	}
	if client.IsErrConnectionFailed(err) {
		return true
	}
	tresc := strings.ToLower(err.Error())
	return strings.Contains(tresc, "no such file or directory") ||
		strings.Contains(tresc, "connection refused") ||
		strings.Contains(tresc, "cannot connect to the docker daemon") ||
		strings.Contains(tresc, "permission denied")
}

// bladSilnikaKontenerow składa zdanie odmowy dla braku silnika kontenerów,
// z nazwą czynności, której silnik zabrakło.
func bladSilnikaKontenerow(czynnosc string) error {
	return bladZasobuDevelopera("na serwerze nie odpowiada żaden silnik kontenerów " +
		"(Docker ani Podman), więc " + czynnosc + " nie ma czym się wykonać; " +
		"naprawa po stronie serwera: " + narzedzieSilnikaKontenerow.Pakiet)
}

// WykazKontenerow obsługuje `developer.container.list`: wykaz kontenerów
// widzianych przez silnik wraz ze stanem dostępności silnika samego.
func (a *adapterDevelopera) WykazKontenerow(ctx context.Context,
	z shared.DeveloperContainerListRequest) (shared.DeveloperContainerListResponse, error) {

	if _, err := a.oknoDevelopera(z.WindowId); err != nil {
		return shared.DeveloperContainerListResponse{}, err
	}
	pusta := shared.DeveloperContainerListResponse{
		Containers:      []shared.ContainerInfo{},
		EngineAvailable: false,
	}

	silnik, err := silnikKontenerow()
	if err != nil {
		return pusta, nil
	}
	defer silnik.Close()

	kontekst, przerwij := context.WithTimeout(ctx, czasSilnikaKontenerow)
	defer przerwij()

	wykaz, err := silnik.ContainerList(kontekst, container.ListOptions{
		All: z.All != nil && *z.All,
	})
	if err != nil {
		if brakSilnikaKontenerow(err) {
			return pusta, nil
		}
		return shared.DeveloperContainerListResponse{}, bladWykonaniaDevelopera(
			"silnik kontenerów odmówił wykazu: " + err.Error())
	}

	kontenery := make([]shared.ContainerInfo, 0, len(wykaz))
	for _, pozycja := range wykaz {
		kontenery = append(kontenery, kontenerKontraktu(pozycja))
	}
	sort.SliceStable(kontenery, func(i, j int) bool { return kontenery[i].Name < kontenery[j].Name })

	odpowiedz := shared.DeveloperContainerListResponse{
		Containers:      kontenery,
		EngineAvailable: true,
	}
	if z.IncludeImages != nil && *z.IncludeImages {
		obrazy, err := silnik.ImageList(kontekst, image.ListOptions{})
		if err == nil {
			wykazObrazow := make([]shared.ImageInfo, 0, len(obrazy))
			for _, obraz := range obrazy {
				rozmiar, utworzono := obraz.Size, obraz.Created*1000
				wykazObrazow = append(wykazObrazow, shared.ImageInfo{
					Id:        obraz.ID,
					Tags:      obraz.RepoTags,
					SizeBytes: &rozmiar,
					CreatedAt: &utworzono,
				})
			}
			odpowiedz.Images = wykazObrazow
		}
	}
	return odpowiedz, nil
}

// kontenerKontraktu przekłada opis kontenera zwrócony silnikiem na kształt
// kontenera zwracany kontraktem komunikacji.
func kontenerKontraktu(pozycja container.Summary) shared.ContainerInfo {
	opis := shared.ContainerInfo{
		Id:     pozycja.ID,
		Name:   nazwaKontenera(pozycja.Names, pozycja.ID),
		Status: stanKontenera(pozycja.State),
	}
	if pozycja.Image != "" {
		opis.Image = wskaznikTekstu(pozycja.Image)
	}
	if pozycja.Created > 0 {
		utworzono := pozycja.Created * 1000
		opis.CreatedAt = &utworzono
	}
	porty := make([]string, 0, len(pozycja.Ports))
	for _, port := range pozycja.Ports {
		zapis := itoa(int(port.PrivatePort)) + "/" + port.Type
		if port.PublicPort > 0 {
			zapis = itoa(int(port.PublicPort)) + "→" + zapis
		}
		porty = append(porty, zapis)
	}
	if len(porty) > 0 {
		sort.Strings(porty)
		opis.Ports = porty
	}
	return opis
}

// nazwaKontenera bierze pierwszą nazwę silnika bez wiodącego ukośnika.
// Kontener bez nazwy zostaje przy skróconym identyfikatorze — pusta nazwa
// w oknie byłaby wierszem, którego nie da się wskazać.
func nazwaKontenera(nazwy []string, identyfikator string) string {
	for _, nazwa := range nazwy {
		if tresc := strings.TrimPrefix(strings.TrimSpace(nazwa), "/"); tresc != "" {
			return tresc
		}
	}
	if len(identyfikator) > 12 {
		return identyfikator[:12]
	}
	return identyfikator
}

// stanKontenera przekłada stan kontenera zwrócony silnikiem na stan kontenera
// wymieniony w kontrakcie komunikacji.
func stanKontenera(stan string) shared.ContainerStatus {
	switch strings.ToLower(strings.TrimSpace(stan)) {
	case "running":
		return shared.ContainerStatusRunning
	case "paused":
		return shared.ContainerStatusPaused
	case "dead":
		return shared.ContainerStatusDead
	case "exited", "removing", "restarting":
		// Kontener wznawiany zgłasza się jako zakończony do chwili, w której wystartuje.
		return shared.ContainerStatusExited
	default:
		return shared.ContainerStatusCreated
	}
}

// CzynnoscKontenera obsługuje `developer.container.action`: start, stop,
// restart, wznowienie albo usunięcie kontenera wskazanego identyfikatorem.
func (a *adapterDevelopera) CzynnoscKontenera(ctx context.Context,
	z shared.DeveloperContainerActionRequest) (shared.DeveloperContainerActionResponse, error) {

	identyfikator := strings.TrimSpace(z.ContainerId)
	if identyfikator == "" {
		return shared.DeveloperContainerActionResponse{}, bladZadaniaDevelopera(
			"czynność wymaga wskazania kontenera")
	}
	silnik, err := silnikKontenerow()
	if err != nil {
		return shared.DeveloperContainerActionResponse{}, bladSilnikaKontenerow(
			"czynność na kontenerze")
	}
	defer silnik.Close()

	kontekst, przerwij := context.WithTimeout(ctx, czasSilnikaKontenerow)
	defer przerwij()

	wyjscie := ""
	switch z.Action {
	case shared.ContainerActionKindStart:
		err = silnik.ContainerStart(kontekst, identyfikator, container.StartOptions{})
	case shared.ContainerActionKindStop:
		err = silnik.ContainerStop(kontekst, identyfikator, container.StopOptions{})
	case shared.ContainerActionKindRestart:
		err = silnik.ContainerRestart(kontekst, identyfikator, container.StopOptions{})
	case shared.ContainerActionKindRemove:
		err = silnik.ContainerRemove(kontekst, identyfikator, container.RemoveOptions{Force: true})
	case shared.ContainerActionKindLogs:
		wyjscie, err = logKontenera(kontekst, silnik, identyfikator, z.Tail)
	default:
		return shared.DeveloperContainerActionResponse{}, bladZadaniaDevelopera(
			"nieznana czynność na kontenerze: " + string(z.Action))
	}
	if err != nil {
		if brakSilnikaKontenerow(err) {
			return shared.DeveloperContainerActionResponse{}, bladSilnikaKontenerow(
				"czynność na kontenerze")
		}
		return shared.DeveloperContainerActionResponse{}, bladWykonaniaDevelopera(
			"silnik kontenerów odmówił czynności " + string(z.Action) + ": " + err.Error())
	}

	odpowiedz := shared.DeveloperContainerActionResponse{}
	if wyjscie != "" {
		odpowiedz.Output = wskaznikTekstu(wyjscie)
	}
	// Usunięty kontener nie ma już opisu — odpowiedź niesie sam identyfikator i stan.
	if z.Action == shared.ContainerActionKindRemove {
		odpowiedz.Container = shared.ContainerInfo{
			Id:     identyfikator,
			Name:   nazwaKontenera(nil, identyfikator),
			Status: shared.ContainerStatusExited,
		}
		return odpowiedz, nil
	}
	odpowiedz.Container = opisKontenera(kontekst, silnik, identyfikator)
	return odpowiedz, nil
}

// opisKontenera odczytuje stan kontenera po wykonanej czynności, do złożenia
// odpowiedzi bieżącym opisem, nie samą deklaracją sukcesu.
func opisKontenera(ctx context.Context, silnik *client.Client,
	identyfikator string) shared.ContainerInfo {

	wykaz, err := silnik.ContainerList(ctx, container.ListOptions{
		All:     true,
		Filters: filters.NewArgs(filters.Arg("id", identyfikator)),
	})
	if err == nil && len(wykaz) > 0 {
		return kontenerKontraktu(wykaz[0])
	}
	return shared.ContainerInfo{
		Id:     identyfikator,
		Name:   nazwaKontenera(nil, identyfikator),
		Status: shared.ContainerStatusCreated,
	}
}

// logKontenera czyta ogon logu kontenera o żądanej głębokości wierszy,
// oczyszczony z nagłówków multipleksowania strumieni silnika.
func logKontenera(ctx context.Context, silnik *client.Client, identyfikator string,
	ogon *int) (string, error) {

	ile := najwiecejWierszyLoguKontenera
	if ogon != nil && *ogon > 0 {
		ile = *ogon
	}
	strumien, err := silnik.ContainerLogs(ctx, identyfikator, container.LogsOptions{
		ShowStdout: true,
		ShowStderr: true,
		Tail:       itoa(ile),
	})
	if err != nil {
		return "", err
	}
	defer strumien.Close()

	bajty, err := io.ReadAll(io.LimitReader(strumien, 4<<20))
	if err != nil && !errors.Is(err, io.EOF) {
		return "", err
	}
	return oczyscLogKontenera(bajty), nil
}

// oczyscLogKontenera zdejmuje ośmiobajtowe nagłówki multipleksowania, którymi
// silnik znakuje porcje strumienia przeplatającego wyjście i diagnostykę.
func oczyscLogKontenera(bajty []byte) string {
	zapis := strings.Builder{}
	for i := 0; i+8 <= len(bajty); {
		if bajty[i] > 2 || bajty[i+1] != 0 || bajty[i+2] != 0 || bajty[i+3] != 0 {
			// Strumień bez multipleksowania (kontener z terminalem) — reszta
			// jest zwykłym tekstem.
			zapis.Write(bajty[i:])
			break
		}
		dlugosc := int(bajty[i+4])<<24 | int(bajty[i+5])<<16 | int(bajty[i+6])<<8 | int(bajty[i+7])
		poczatek := i + 8
		koniec := poczatek + dlugosc
		if dlugosc < 0 || koniec > len(bajty) {
			koniec = len(bajty)
		}
		zapis.Write(bajty[poczatek:koniec])
		i = koniec
	}
	return zapis.String()
}

// BudujObraz obsługuje `developer.image.build`: pakuje kontekst budowania
// i prowadzi budowanie obrazu silnikiem, czytając dziennik postępu.
func (a *adapterDevelopera) BudujObraz(ctx context.Context,
	z shared.DeveloperImageBuildRequest) (shared.DeveloperImageBuildResponse, error) {

	okno, sciezka, err := a.plikOkna(z.WindowId, z.Dockerfile)
	if err != nil {
		return shared.DeveloperImageBuildResponse{}, err
	}
	if err := sprawdzZmianeSystemu(okno.TrybUprawnien, "budowanie obrazu kontenera"); err != nil {
		return shared.DeveloperImageBuildResponse{}, err
	}
	znacznik := strings.TrimSpace(z.Tag)
	if znacznik == "" {
		return shared.DeveloperImageBuildResponse{}, bladZadaniaDevelopera(
			"budowanie obrazu wymaga znacznika")
	}
	if z.Push != nil && *z.Push {
		return shared.DeveloperImageBuildResponse{}, bladZadaniaDevelopera(
			"wypchnięcie obrazu do rejestru wymaga poświadczeń rejestru, których " +
				"kontrakt tej komendy nie niesie; obraz zbuduje się, a wypchnięcie " +
				"idzie osobną drogą")
	}

	silnik, err := silnikKontenerow()
	if err != nil {
		return shared.DeveloperImageBuildResponse{}, bladSilnikaKontenerow("budowanie obrazu")
	}
	defer silnik.Close()

	kontekst, przerwij := context.WithTimeout(ctx, czasBudowaniaObrazu)
	defer przerwij()

	katalog := filepath.Dir(sciezka)
	archiwum, err := archiwumKontekstuBudowania(katalog)
	if err != nil {
		return shared.DeveloperImageBuildResponse{}, err
	}
	defer archiwum.Close()

	nastawy := dockerBuildOptions(filepath.Base(sciezka), znacznik, z.BuildArgs)
	odpowiedz, err := silnik.ImageBuild(kontekst, archiwum, nastawy)
	if err != nil {
		if brakSilnikaKontenerow(err) {
			return shared.DeveloperImageBuildResponse{}, bladSilnikaKontenerow("budowanie obrazu")
		}
		return shared.DeveloperImageBuildResponse{}, bladWykonaniaDevelopera(
			"silnik kontenerów odmówił budowania: " + err.Error())
	}
	defer odpowiedz.Body.Close()

	identyfikator, err := odczytajPrzebiegBudowaniaObrazu(odpowiedz.Body)
	if err != nil {
		return shared.DeveloperImageBuildResponse{}, bladWykonaniaDevelopera(
			"budowanie obrazu " + znacznik + " zawiodło: " + err.Error())
	}
	if identyfikator == "" {
		identyfikator = znacznik
	}
	return shared.DeveloperImageBuildResponse{ImageId: identyfikator}, nil
}

// odczytajPrzebiegBudowaniaObrazu czyta strumień postępu i wyławia identyfikator
// gotowego obrazu albo pierwszy błąd budowania.
func odczytajPrzebiegBudowaniaObrazu(zrodlo io.Reader) (string, error) {
	czytnik := json.NewDecoder(zrodlo)
	identyfikator := ""
	for {
		var krok struct {
			Stream string `json:"stream"`
			Error  string `json:"error"`
			Aux    struct {
				ID string `json:"ID"`
			} `json:"aux"`
		}
		if err := czytnik.Decode(&krok); err != nil {
			if errors.Is(err, io.EOF) {
				return identyfikator, nil
			}
			return identyfikator, err
		}
		if krok.Error != "" {
			return "", errors.New(strings.TrimSpace(krok.Error))
		}
		if krok.Aux.ID != "" {
			identyfikator = krok.Aux.ID
		}
		if wskazanie, jest := strings.CutPrefix(strings.TrimSpace(krok.Stream),
			"Successfully built "); jest {
			identyfikator = strings.TrimSpace(wskazanie)
		}
	}
}

// KompozycjaKontenerow obsługuje `developer.compose.up`: czyta plik stosu
// i podnosi jego usługi kontenerami przez silnik.
func (a *adapterDevelopera) KompozycjaKontenerow(ctx context.Context,
	z shared.DeveloperComposeUpRequest) (shared.DeveloperComposeUpResponse, error) {

	okno, sciezka, err := a.plikOkna(z.WindowId, z.File)
	if err != nil {
		return shared.DeveloperComposeUpResponse{}, err
	}
	if err := sprawdzZmianeSystemu(okno.TrybUprawnien, "uruchomienie stosu kontenerów"); err != nil {
		return shared.DeveloperComposeUpResponse{}, err
	}

	opis, err := wczytajOpisKompozycji(sciezka)
	if err != nil {
		return shared.DeveloperComposeUpResponse{}, err
	}
	wybrane := usluigWybrane(opis, z.Services)
	if len(wybrane) == 0 {
		return shared.DeveloperComposeUpResponse{}, bladZadaniaDevelopera(
			"plik " + z.File + " nie opisuje ani jednej usługi do uruchomienia")
	}

	silnik, err := silnikKontenerow()
	if err != nil {
		return shared.DeveloperComposeUpResponse{}, bladSilnikaKontenerow("uruchomienie stosu")
	}
	defer silnik.Close()

	kontekst, przerwij := context.WithTimeout(ctx, czasBudowaniaObrazu)
	defer przerwij()

	stos := nazwaStosu(sciezka)
	if z.Down != nil && *z.Down {
		return a.zatrzymajStos(kontekst, silnik, stos, wybrane)
	}
	return a.podniesStos(kontekst, silnik, stos, wybrane, opis)
}

// opisKompozycji jest odczytanym plikiem `docker-compose.yml` w zakresie pól,
// którym ten moduł się posługuje, nie każdą konstrukcją formatu Compose.
type opisKompozycji struct {
	Services map[string]struct {
		Image       string            `yaml:"image"`
		Command     any               `yaml:"command"`
		Environment any               `yaml:"environment"`
		Ports       []string          `yaml:"ports"`
		DependsOn   any               `yaml:"depends_on"`
		Labels      map[string]string `yaml:"labels"`
		Restart     string            `yaml:"restart"`
	} `yaml:"services"`
}

// wczytajOpisKompozycji czyta plik stosu Compose z katalogu okna i rozbiera
// go do zakresu pól, którym moduł się posługuje.
func wczytajOpisKompozycji(sciezka string) (opisKompozycji, error) {
	bajty, err := os.ReadFile(sciezka)
	if err != nil {
		return opisKompozycji{}, bladZasobuDevelopera(
			"nie można odczytać pliku stosu " + sciezka + ": " + err.Error())
	}
	var opis opisKompozycji
	if err := yaml.Unmarshal(bajty, &opis); err != nil {
		return opisKompozycji{}, bladZadaniaDevelopera(
			"plik " + sciezka + " nie jest czytelnym opisem stosu: " + err.Error())
	}
	return opis, nil
}

// usluigWybrane zwraca nazwy usług do uruchomienia w kolejności ustalonej:
// wszystkie usługi opisu, gdy żądanie nie wskazało zawężenia wprost.
func usluigWybrane(opis opisKompozycji, zadane []string) []string {
	wszystkie := make([]string, 0, len(opis.Services))
	for nazwa := range opis.Services {
		wszystkie = append(wszystkie, nazwa)
	}
	sort.Strings(wszystkie)
	if len(zadane) == 0 {
		return wszystkie
	}
	wybrane := make([]string, 0, len(zadane))
	for _, nazwa := range wszystkie {
		for _, zadana := range zadane {
			if nazwa == strings.TrimSpace(zadana) {
				wybrane = append(wybrane, nazwa)
				break
			}
		}
	}
	return wybrane
}

// nazwaStosu wywodzi nazwę stosu z katalogu pliku — tak samo, jak robi to
// Compose. Dzięki temu kontenery założone przez moduł i przez zewnętrzny klient
// noszą tę samą etykietę i widać je w jednym wykazie.
func nazwaStosu(sciezka string) string {
	nazwa := filepath.Base(filepath.Dir(sciezka))
	nazwa = strings.ToLower(strings.Map(func(znak rune) rune {
		if (znak >= 'a' && znak <= 'z') || (znak >= 'A' && znak <= 'Z') ||
			(znak >= '0' && znak <= '9') || znak == '_' || znak == '-' {
			return znak
		}
		return '-'
	}, nazwa))
	if nazwa == "" || nazwa == "-" {
		return "danaco"
	}
	return nazwa
}

// etykietaStosu jest kluczem etykiety, po którym rdzeń rozpoznaje kontenery
// należące do jednego stosu Compose podniesionego oknem.
const etykietaStosu = "com.docker.compose.project"

// podniesStos zakłada i startuje kontenery usług wybranych z opisu stosu,
// budując najpierw obrazy usług, które opis wskazuje budowaniem.
func (a *adapterDevelopera) podniesStos(ctx context.Context, silnik *client.Client,
	stos string, usluigi []string, opis opisKompozycji) (shared.DeveloperComposeUpResponse, error) {

	zapis := strings.Builder{}
	uruchomione := make([]shared.ContainerInfo, 0, len(usluigi))

	for _, nazwa := range usluigi {
		usluga := opis.Services[nazwa]
		if strings.TrimSpace(usluga.Image) == "" {
			// Usługa budowana z `build:` wymaga wcześniejszego zbudowania obrazu.
			zapis.WriteString(nazwa + ": pominięta — usługa nie wskazuje obrazu (`image`), " +
				"a budowanie z `build` idzie komendą developer.image.build\n")
			continue
		}
		kontener := stos + "-" + nazwa

		// Kontener o tej nazwie z poprzedniego biegu zostaje usunięty przed startem.
		_ = silnik.ContainerRemove(ctx, kontener, container.RemoveOptions{Force: true})

		nastawy := &container.Config{
			Image:  usluga.Image,
			Labels: map[string]string{etykietaStosu: stos, "com.docker.compose.service": nazwa},
			Env:    zmienneUslugi(usluga.Environment),
		}
		if polecenie := polecenieUslugi(usluga.Command); len(polecenie) > 0 {
			nastawy.Cmd = polecenie
		}
		for klucz, wartosc := range usluga.Labels {
			nastawy.Labels[klucz] = wartosc
		}

		gospodarz, err := nastawyGospodarza(usluga.Ports, usluga.Restart)
		if err != nil {
			zapis.WriteString(nazwa + ": " + err.Error() + "\n")
			continue
		}

		zalozony, err := silnik.ContainerCreate(ctx, nastawy, gospodarz,
			&network.NetworkingConfig{}, nil, kontener)
		if err != nil {
			zapis.WriteString(nazwa + ": nie można założyć kontenera — " + err.Error() + "\n")
			continue
		}
		if err := silnik.ContainerStart(ctx, zalozony.ID, container.StartOptions{}); err != nil {
			zapis.WriteString(nazwa + ": nie można uruchomić kontenera — " + err.Error() + "\n")
			continue
		}
		zapis.WriteString(nazwa + ": uruchomiona\n")
		uruchomione = append(uruchomione, opisKontenera(ctx, silnik, zalozony.ID))
	}

	odpowiedz := shared.DeveloperComposeUpResponse{Services: uruchomione}
	if zapis.Len() > 0 {
		odpowiedz.Output = wskaznikTekstu(strings.TrimRight(zapis.String(), "\n"))
	}
	return odpowiedz, nil
}

// zatrzymajStos zatrzymuje i usuwa kontenery stosu rozpoznane etykietą stosu,
// bez względu na to, które usługi opis stosu dziś wymienia.
func (a *adapterDevelopera) zatrzymajStos(ctx context.Context, silnik *client.Client,
	stos string, usluigi []string) (shared.DeveloperComposeUpResponse, error) {

	zapis := strings.Builder{}
	zatrzymane := make([]shared.ContainerInfo, 0, len(usluigi))
	for _, nazwa := range usluigi {
		kontener := stos + "-" + nazwa
		opis := opisKontenera(ctx, silnik, kontener)
		if err := silnik.ContainerRemove(ctx, kontener,
			container.RemoveOptions{Force: true}); err != nil {
			zapis.WriteString(nazwa + ": " + err.Error() + "\n")
			continue
		}
		opis.Status = shared.ContainerStatusExited
		zatrzymane = append(zatrzymane, opis)
		zapis.WriteString(nazwa + ": zatrzymana i usunięta\n")
	}
	odpowiedz := shared.DeveloperComposeUpResponse{Services: zatrzymane}
	if zapis.Len() > 0 {
		odpowiedz.Output = wskaznikTekstu(strings.TrimRight(zapis.String(), "\n"))
	}
	return odpowiedz, nil
}

// zmienneUslugi sprowadza obie postacie zapisu zmiennych (wykaz i mapa) do
// jednego kształtu `KLUCZ=wartość`, którego wymaga silnik.
func zmienneUslugi(zapis any) []string {
	switch wartosc := zapis.(type) {
	case []any:
		zmienne := make([]string, 0, len(wartosc))
		for _, pozycja := range wartosc {
			if tekst, jest := pozycja.(string); jest {
				zmienne = append(zmienne, tekst)
			}
		}
		return zmienne
	case map[string]any:
		zmienne := make([]string, 0, len(wartosc))
		for klucz, pozycja := range wartosc {
			zmienne = append(zmienne, klucz+"="+tekstWartosciYaml(pozycja))
		}
		sort.Strings(zmienne)
		return zmienne
	default:
		return nil
	}
}

// polecenieUslugi sprowadza obie postacie zapisu polecenia usługi Compose —
// tekst i wykaz — do jednego wykazu argumentów procesu.
func polecenieUslugi(zapis any) []string {
	switch wartosc := zapis.(type) {
	case string:
		return strings.Fields(wartosc)
	case []any:
		polecenie := make([]string, 0, len(wartosc))
		for _, pozycja := range wartosc {
			polecenie = append(polecenie, tekstWartosciYaml(pozycja))
		}
		return polecenie
	default:
		return nil
	}
}

// tekstWartosciYaml sprowadza wartość odczytaną z pliku YAML do tekstu, bez
// względu na to, czy plik zapisał ją jako napis, liczbę czy wartość logiczną.
func tekstWartosciYaml(wartosc any) string {
	if wartosc == nil {
		return ""
	}
	if tekst, jest := wartosc.(string); jest {
		return tekst
	}
	bajty, err := json.Marshal(wartosc)
	if err != nil {
		return ""
	}
	return strings.Trim(string(bajty), `"`)
}

// nastawyGospodarza składa przypisania portów i nastawę wznowienia kontenera
// usługi, do przekazania silnikowi przy jego założeniu.
func nastawyGospodarza(porty []string, wznowienie string) (*container.HostConfig, error) {
	nastawy := &container.HostConfig{PortBindings: map[nat.Port][]nat.PortBinding{}}
	nastawy.PublishAllPorts = false
	for _, zapis := range porty {
		wewnetrzny, zewnetrzny, err := rozbierzPrzypisaniePortu(zapis)
		if err != nil {
			return nil, err
		}
		nastawy.PortBindings[wewnetrzny] = []nat.PortBinding{{HostPort: zewnetrzny}}
	}
	if wznowienie != "" && wznowienie != "no" {
		nastawy.RestartPolicy = container.RestartPolicy{
			Name: container.RestartPolicyMode(wznowienie),
		}
	}
	return nastawy, nil
}

// rozbierzPrzypisaniePortu rozbiera zapis `8080:80` albo `8080:80/tcp` na
// port gospodarza, port kontenera i protokół.
func rozbierzPrzypisaniePortu(zapis string) (nat.Port, string, error) {
	tresc := strings.TrimSpace(zapis)
	protokol := "tcp"
	if ukosnik := strings.LastIndexByte(tresc, '/'); ukosnik > 0 {
		protokol = tresc[ukosnik+1:]
		tresc = tresc[:ukosnik]
	}
	czesci := strings.Split(tresc, ":")
	switch len(czesci) {
	case 1:
		return nat.Port(czesci[0] + "/" + protokol), "", nil
	case 2:
		return nat.Port(czesci[1] + "/" + protokol), czesci[0], nil
	case 3:
		// Zapis z adresem gospodarza niesie adres, który tu jest pomijany.
		return nat.Port(czesci[2] + "/" + protokol), czesci[1], nil
	default:
		return "", "", errors.New("nieczytelne przypisanie portu: " + zapis)
	}
}

// dockerBuildOptions składa nastawy budowania obrazu przekazywane silnikowi:
// znacznik obrazu, ścieżkę pliku Dockerfile i argumenty budowania.
func dockerBuildOptions(dockerfile, znacznik string,
	parametry json.RawMessage) build.ImageBuildOptions {

	nastawy := build.ImageBuildOptions{
		Tags:       []string{znacznik},
		Dockerfile: dockerfile,
		Remove:     true,
	}
	if len(parametry) == 0 {
		return nastawy
	}
	wpisy := map[string]string{}
	if err := json.Unmarshal(parametry, &wpisy); err != nil {
		return nastawy
	}
	nastawy.BuildArgs = make(map[string]*string, len(wpisy))
	for klucz, wartosc := range wpisy {
		nastawy.BuildArgs[klucz] = wskaznikTekstu(wartosc)
	}
	return nastawy
}

// archiwumKontekstuBudowania pakuje katalog kontekstu budowania strumieniem
// do archiwum tar, z regułami `.dockerignore` uwzględnionymi.
func archiwumKontekstuBudowania(katalog string) (io.ReadCloser, error) {
	pomijaj, err := regulyPomijaniaKontekstu(katalog)
	if err != nil {
		return nil, err
	}
	odczyt, zapis := io.Pipe()
	go func() {
		archiwum := tar.NewWriter(zapis)
		blad := filepath.Walk(katalog, func(sciezka string, opis os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			wzgledna, err := filepath.Rel(katalog, sciezka)
			if err != nil {
				return err
			}
			if wzgledna == "." {
				return nil
			}
			wzgledna = filepath.ToSlash(wzgledna)
			if pomijaj(wzgledna) {
				if opis.IsDir() {
					return filepath.SkipDir
				}
				return nil
			}
			// Dowiązania i pliki urządzeń zostają poza archiwum kontekstu.
			if !opis.Mode().IsRegular() && !opis.IsDir() {
				return nil
			}
			naglowek, err := tar.FileInfoHeader(opis, "")
			if err != nil {
				return err
			}
			naglowek.Name = wzgledna
			if err := archiwum.WriteHeader(naglowek); err != nil {
				return err
			}
			if opis.IsDir() {
				return nil
			}
			plik, err := os.Open(sciezka)
			if err != nil {
				return err
			}
			defer plik.Close()
			_, err = io.Copy(archiwum, plik)
			return err
		})
		if blad != nil {
			_ = archiwum.Close()
			_ = zapis.CloseWithError(blad)
			return
		}
		if err := archiwum.Close(); err != nil {
			_ = zapis.CloseWithError(err)
			return
		}
		_ = zapis.Close()
	}()
	return odczyt, nil
}

// regulyPomijaniaKontekstu czyta `.dockerignore` i oddaje sito ścieżek.
// Brak pliku znaczy sito domyślne, które i tak odrzuca katalog repozytorium
// Gita — nikt nie buduje obrazu po to, żeby mieć w nim historię zmian.
func regulyPomijaniaKontekstu(katalog string) (func(string) bool, error) {
	wzorce := []string{".git"}
	bajty, err := os.ReadFile(filepath.Join(katalog, ".dockerignore"))
	if err == nil {
		for _, wiersz := range strings.Split(string(bajty), "\n") {
			tresc := strings.TrimSpace(wiersz)
			if tresc == "" || strings.HasPrefix(tresc, "#") {
				continue
			}
			wzorce = append(wzorce, strings.TrimSuffix(tresc, "/"))
		}
	} else if !os.IsNotExist(err) {
		return nil, bladWykonaniaDevelopera(
			"nie można odczytać .dockerignore: " + err.Error())
	}
	return func(wzgledna string) bool {
		for _, wzorzec := range wzorce {
			if wzgledna == wzorzec || strings.HasPrefix(wzgledna, wzorzec+"/") {
				return true
			}
			if dopasowane, err := filepath.Match(wzorzec, wzgledna); err == nil && dopasowane {
				return true
			}
		}
		return false
	}, nil
}
