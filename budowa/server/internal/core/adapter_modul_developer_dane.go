// Odpowiedzialność pliku: komendy Data Console (developer.data.*) na sterownikach PostgreSQL, MySQL
// i SQLite wkompilowanych w rdzeń; wiersz połączenia niesie odwołanie do sejfu, nie hasło.
package core

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/jackc/pgx/v5/stdlib"
	_ "modernc.org/sqlite"

	"danacoconsole/server/internal/dane"
	"danacoconsole/shared"
)

const (
	przedrostekPolaczeniaDanych = "conn-"
	czasZapytaniaDanych         = 60 * time.Second
	// Konsola pokazuje próbkę wyniku, nie eksport; pełny wynik nie mieści się w pamięci rdzenia ani w oknie.
	najwiecejWierszyWyniku = 1000
)

func (a *adapterDevelopera) UstawPolaczenieDanych(ctx context.Context,
	z shared.DeveloperDataConnectionSetRequest) (shared.DeveloperDataConnectionSetResponse, error) {

	okno, err := a.oknoDevelopera(z.WindowId)
	if err != nil {
		return shared.DeveloperDataConnectionSetResponse{}, err
	}
	if a.repozytorium == nil {
		return shared.DeveloperDataConnectionSetResponse{}, bladZasobuDevelopera(
			"serwer nie ma miejsca na opisy połączeń bazodanowych")
	}
	nazwa := strings.TrimSpace(z.Name)
	if nazwa == "" {
		return shared.DeveloperDataConnectionSetResponse{}, bladZadaniaDevelopera(
			"połączenie wymaga nazwy")
	}
	baza := strings.TrimSpace(z.Database)
	if baza == "" {
		return shared.DeveloperDataConnectionSetResponse{}, bladZadaniaDevelopera(
			"połączenie wymaga wskazania bazy")
	}
	if !znanySilnikDanych(z.Engine) {
		return shared.DeveloperDataConnectionSetResponse{}, bladZadaniaDevelopera(
			"nieznany silnik bazy danych: " + string(z.Engine))
	}
	// Baza SQLite jest plikiem, więc jej wskazanie przechodzi przez bramę obszaru.
	if z.Engine == shared.DataEngineSqlite {
		sciezka, err := sciezkaWObszarze(a.korzenieOkna(okno), baza, plikIstnieje)
		if err != nil {
			return shared.DeveloperDataConnectionSetResponse{}, err
		}
		baza = sciezka
	}

	kod := strings.TrimSpace(tekstWskazaniaDevelopera(z.ConnectionId))
	if kod == "" {
		kod = nowyIdentyfikator(przedrostekPolaczeniaDanych)
	}
	wiersz := dane.PolaczenieDanych{
		Kod:           kod,
		OknoKod:       okno.Id,
		Nazwa:         nazwa,
		Silnik:        z.Engine,
		Host:          z.Host,
		Baza:          baza,
		Uzytkownik:    z.User,
		Poswiadczenie: z.CredentialRef,
		TylkoOdczyt:   z.ReadOnly == nil || *z.ReadOnly,
	}
	if z.Port != nil {
		port := int64(*z.Port)
		wiersz.Port = &port
	}
	if err := a.repozytorium.ZapiszPolaczenieDanych(ctx, wiersz); err != nil {
		return shared.DeveloperDataConnectionSetResponse{}, bladZapisuDevelopera(err,
			"nie można zapisać połączenia "+nazwa+": "+err.Error())
	}

	zapisane, err := a.repozytorium.PolaczenieDanychPoKodzie(ctx, kod)
	if err != nil {
		return shared.DeveloperDataConnectionSetResponse{}, bladWykonaniaDevelopera(
			"połączenie " + nazwa + " zapisało się, lecz nie daje się odczytać")
	}
	return shared.DeveloperDataConnectionSetResponse{
		Connection: polaczenieKontraktu(zapisane),
	}, nil
}

func (a *adapterDevelopera) WykazPolaczenDanych(ctx context.Context,
	z shared.DeveloperDataConnectionListRequest) (shared.DeveloperDataConnectionListResponse, error) {

	okno, err := a.oknoDevelopera(z.WindowId)
	if err != nil {
		return shared.DeveloperDataConnectionListResponse{}, err
	}
	if a.repozytorium == nil {
		return shared.DeveloperDataConnectionListResponse{}, bladZasobuDevelopera(
			"serwer nie ma miejsca na opisy połączeń bazodanowych")
	}
	wiersze, err := a.repozytorium.PolaczeniaDanych(ctx, okno.Id)
	if err != nil {
		return shared.DeveloperDataConnectionListResponse{}, bladWykonaniaDevelopera(
			"nie można odczytać połączeń okna " + okno.Id + ": " + err.Error())
	}
	polaczenia := make([]shared.DataConnection, 0, len(wiersze))
	for _, wiersz := range wiersze {
		polaczenia = append(polaczenia, polaczenieKontraktu(wiersz))
	}
	return shared.DeveloperDataConnectionListResponse{Connections: polaczenia}, nil
}

// Odpowiedź niesie odwołanie do sejfu, nigdy hasło.
func polaczenieKontraktu(wiersz dane.PolaczenieDanych) shared.DataConnection {
	opis := shared.DataConnection{
		Id:            wiersz.Kod,
		Name:          wiersz.Nazwa,
		Engine:        wiersz.Silnik,
		Host:          wiersz.Host,
		Database:      wiersz.Baza,
		User:          wiersz.Uzytkownik,
		CredentialRef: wiersz.Poswiadczenie,
		ReadOnly:      wiersz.TylkoOdczyt,
	}
	if wiersz.Port != nil {
		port := int(*wiersz.Port)
		opis.Port = &port
	}
	return opis
}

func znanySilnikDanych(silnik shared.DataEngine) bool {
	switch silnik {
	case shared.DataEnginePostgres, shared.DataEngineMysql, shared.DataEngineSqlite:
		return true
	default:
		return false
	}
}

// Głębokość zawęża odczyt: wykaz kolumn wszystkich tabel dużej bazy liczy tysiące wierszy.
func (a *adapterDevelopera) SchematDanych(ctx context.Context,
	z shared.DeveloperDataSchemaGetRequest) (shared.DeveloperDataSchemaGetResponse, error) {

	polaczenie, baza, err := a.otworzPolaczenieDanych(ctx, z.ConnectionId)
	if err != nil {
		return shared.DeveloperDataSchemaGetResponse{}, err
	}
	defer baza.Close()

	glebokosc := 2
	if z.Depth != nil && *z.Depth > 0 {
		glebokosc = *z.Depth
	}
	wskazanie := strings.TrimSpace(tekstWskazaniaDevelopera(z.Path))

	kontekst, przerwij := context.WithTimeout(ctx, czasZapytaniaDanych)
	defer przerwij()

	wezly, err := wezlySchematu(kontekst, baza, polaczenie, wskazanie, glebokosc)
	if err != nil {
		return shared.DeveloperDataSchemaGetResponse{}, bladWykonaniaDevelopera(
			"nie można odczytać schematu połączenia " + polaczenie.Nazwa + ": " + err.Error())
	}
	return shared.DeveloperDataSchemaGetResponse{Nodes: wezly}, nil
}

// Katalog systemowy jest u każdego silnika inny; wspólny jest kształt odpowiedzi.
func wezlySchematu(ctx context.Context, baza *sql.DB, polaczenie dane.PolaczenieDanych,
	wskazanie string, glebokosc int) ([]shared.DataSchemaNode, error) {

	wezly := make([]shared.DataSchemaNode, 0, 64)
	zapytanie, argumenty := zapytanieOSchemat(polaczenie.Silnik, wskazanie)
	wiersze, err := baza.QueryContext(ctx, zapytanie, argumenty...)
	if err != nil {
		return nil, err
	}
	defer wiersze.Close()

	tabele := map[string]bool{}
	for wiersze.Next() {
		var tabela, kolumna, typ, dopuszczaPuste string
		if err := wiersze.Scan(&tabela, &kolumna, &typ, &dopuszczaPuste); err != nil {
			return nil, err
		}
		if !tabele[tabela] {
			tabele[tabela] = true
			wezly = append(wezly, shared.DataSchemaNode{
				Path: tabela,
				Name: tabela,
				Kind: shared.SchemaNodeKindTable,
			})
		}
		if glebokosc < 2 {
			continue
		}
		pusta := strings.EqualFold(dopuszczaPuste, "YES") || dopuszczaPuste == "1"
		wezly = append(wezly, shared.DataSchemaNode{
			Path:       tabela + "/" + kolumna,
			ParentPath: wskaznikTekstu(tabela),
			Name:       kolumna,
			Kind:       shared.SchemaNodeKindColumn,
			DataType:   wskaznikTekstu(typ),
			Nullable:   wskaznikPrawdy(pusta),
		})
	}
	if err := wiersze.Err(); err != nil {
		return nil, err
	}
	sort.SliceStable(wezly, func(i, j int) bool { return wezly[i].Path < wezly[j].Path })
	return wezly, nil
}

// Wskazanie tabeli wchodzi parametrem: nazwa wklejona do polecenia jest wstrzyknięciem SQL.
func zapytanieOSchemat(silnik shared.DataEngine, wskazanie string) (string, []any) {
	switch silnik {
	case shared.DataEnginePostgres:
		return `SELECT table_name, column_name, data_type, is_nullable
		        FROM information_schema.columns
		        WHERE table_schema NOT IN ('pg_catalog','information_schema')
		          AND ($1 = '' OR table_name = $1)
		        ORDER BY table_name, ordinal_position`, []any{wskazanie}
	case shared.DataEngineMysql:
		return `SELECT table_name, column_name, data_type, is_nullable
		        FROM information_schema.columns
		        WHERE table_schema = DATABASE()
		          AND (? = '' OR table_name = ?)
		        ORDER BY table_name, ordinal_position`, []any{wskazanie, wskazanie}
	default:
		// SQLite nie ma information_schema; ten odczyt daje pragma_table_info.
		return `SELECT m.name, p.name, p.type,
		               CASE WHEN p."notnull" = 0 THEN 'YES' ELSE 'NO' END
		        FROM sqlite_master m
		        JOIN pragma_table_info(m.name) p
		        WHERE m.type = 'table' AND m.name NOT LIKE 'sqlite_%'
		          AND (? = '' OR m.name = ?)
		        ORDER BY m.name, p.cid`, []any{wskazanie, wskazanie}
	}
}

func (a *adapterDevelopera) WykonajZapytanieDanych(ctx context.Context,
	z shared.DeveloperDataQueryRunRequest) (shared.DeveloperDataQueryRunResponse, error) {

	polaczenie, baza, err := a.otworzPolaczenieDanych(ctx, z.ConnectionId)
	if err != nil {
		return shared.DeveloperDataQueryRunResponse{}, err
	}
	defer baza.Close()

	polecenie := strings.TrimSpace(z.Sql)
	if polecenie == "" {
		return shared.DeveloperDataQueryRunResponse{}, bladZadaniaDevelopera(
			"wykonanie wymaga treści polecenia SQL")
	}
	zmieniajace := poleceniZmieniajace(polecenie)
	if polaczenie.TylkoOdczyt && zmieniajace {
		return shared.DeveloperDataQueryRunResponse{}, bladDostepuDevelopera(
			"połączenie " + polaczenie.Nazwa + " jest oznaczone jako tylko do odczytu, " +
				"a to polecenie zmienia dane; naprawa: zdjąć nastawę tylko do odczytu " +
				"z opisu połączenia")
	}
	if z.Explain != nil && *z.Explain {
		polecenie = przedrostekPlanu(polaczenie.Silnik) + " " + polecenie
		zmieniajace = false
	}

	limit := najwiecejWierszyWyniku
	if z.Limit != nil && *z.Limit > 0 && *z.Limit < limit {
		limit = *z.Limit
	}

	kontekst, przerwij := context.WithTimeout(ctx, czasZapytaniaDanych)
	defer przerwij()

	poczatek := time.Now()
	wynik, err := wykonajPolecenieDanych(kontekst, baza, polecenie, zmieniajace,
		z.Transaction != nil && *z.Transaction, limit)
	if err != nil {
		return shared.DeveloperDataQueryRunResponse{}, bladWykonaniaDevelopera(
			"silnik odmówił wykonania polecenia: " + err.Error())
	}
	wynik.DurationMs = time.Since(poczatek).Milliseconds()
	return shared.DeveloperDataQueryRunResponse{Result: wynik}, nil
}

// Słowo nieznane liczy się jako zmieniające: brama ma mylić się w stronę odmowy.
func poleceniZmieniajace(polecenie string) bool {
	pola := strings.Fields(strings.ToUpper(polecenie))
	if len(pola) == 0 {
		return true
	}
	switch pola[0] {
	case "SELECT", "WITH", "SHOW", "EXPLAIN", "DESCRIBE", "DESC", "PRAGMA":
		return false
	default:
		return true
	}
}

func przedrostekPlanu(silnik shared.DataEngine) string {
	switch silnik {
	case shared.DataEnginePostgres:
		return "EXPLAIN (ANALYZE false, VERBOSE true)"
	case shared.DataEngineSqlite:
		return "EXPLAIN QUERY PLAN"
	default:
		return "EXPLAIN"
	}
}

// Transakcja jest na życzenie żądania, bo część poleceń DDL nie działa w transakcji wymuszonej.
func wykonajPolecenieDanych(ctx context.Context, baza *sql.DB, polecenie string,
	zmieniajace, wTransakcji bool, limit int) (shared.DataQueryResult, error) {

	if zmieniajace {
		if wTransakcji {
			transakcja, err := baza.BeginTx(ctx, nil)
			if err != nil {
				return shared.DataQueryResult{}, err
			}
			wynik, err := transakcja.ExecContext(ctx, polecenie)
			if err != nil {
				_ = transakcja.Rollback()
				return shared.DataQueryResult{}, err
			}
			if err := transakcja.Commit(); err != nil {
				return shared.DataQueryResult{}, err
			}
			return wynikPolecenia(wynik), nil
		}
		wynik, err := baza.ExecContext(ctx, polecenie)
		if err != nil {
			return shared.DataQueryResult{}, err
		}
		return wynikPolecenia(wynik), nil
	}

	wiersze, err := baza.QueryContext(ctx, polecenie)
	if err != nil {
		return shared.DataQueryResult{}, err
	}
	defer wiersze.Close()
	return wynikOdczytu(wiersze, limit)
}

func wynikPolecenia(wynik sql.Result) shared.DataQueryResult {
	odpowiedz := shared.DataQueryResult{
		Columns: []string{},
		Rows:    json.RawMessage("[]"),
	}
	if liczba, err := wynik.RowsAffected(); err == nil {
		zmienione := int(liczba)
		odpowiedz.AffectedRows = &zmienione
	}
	return odpowiedz
}

// Wartości idą do JSON jako tekst; rzutowanie na liczby w rdzeniu gubiłoby precyzję typów silnika.
func wynikOdczytu(wiersze *sql.Rows, limit int) (shared.DataQueryResult, error) {
	kolumny, err := wiersze.Columns()
	if err != nil {
		return shared.DataQueryResult{}, err
	}
	zebrane := make([][]any, 0, 64)
	przyciete := false

	for wiersze.Next() {
		if len(zebrane) >= limit {
			przyciete = true
			break
		}
		komorki := make([]any, len(kolumny))
		wskazniki := make([]any, len(kolumny))
		for i := range komorki {
			wskazniki[i] = &komorki[i]
		}
		if err := wiersze.Scan(wskazniki...); err != nil {
			return shared.DataQueryResult{}, err
		}
		wiersz := make([]any, len(kolumny))
		for i, komorka := range komorki {
			wiersz[i] = wartoscKomorki(komorka)
		}
		zebrane = append(zebrane, wiersz)
	}
	if err := wiersze.Err(); err != nil {
		return shared.DataQueryResult{}, err
	}

	tresc, err := json.Marshal(zebrane)
	if err != nil {
		return shared.DataQueryResult{}, err
	}
	return shared.DataQueryResult{
		Columns:   kolumny,
		Rows:      tresc,
		RowCount:  len(zebrane),
		Truncated: przyciete,
	}, nil
}

func wartoscKomorki(komorka any) any {
	switch wartosc := komorka.(type) {
	case nil:
		return nil
	case []byte:
		return string(wartosc)
	case time.Time:
		return wartosc.UTC().Format(time.RFC3339Nano)
	case bool, int64, float64:
		return wartosc
	default:
		return fmt.Sprint(wartosc)
	}
}

// Rejestr zastosowanych kroków stoi w bazie docelowej, żeby uruchomienie nie powtarzało kroków.
func (a *adapterDevelopera) UruchomMigracjeDanych(ctx context.Context,
	z shared.DeveloperDataMigrationRunRequest) (shared.DeveloperDataMigrationRunResponse, error) {

	okno, err := a.oknoDevelopera(z.WindowId)
	if err != nil {
		return shared.DeveloperDataMigrationRunResponse{}, err
	}
	proba := z.DryRun != nil && *z.DryRun
	if !proba {
		if err := sprawdzZmianeSystemu(okno.TrybUprawnien, "uruchomienie migracji bazy"); err != nil {
			return shared.DeveloperDataMigrationRunResponse{}, err
		}
	}
	kierunek := strings.ToLower(strings.TrimSpace(z.Direction))
	if kierunek != "up" && kierunek != "down" {
		return shared.DeveloperDataMigrationRunResponse{}, bladZadaniaDevelopera(
			"kierunek migracji ma brzmieć `up` albo `down`; otrzymano " + z.Direction)
	}
	if kierunek == "down" {
		return shared.DeveloperDataMigrationRunResponse{}, bladZadaniaDevelopera(
			"wycofywanie migracji wymaga plików odwrotnych, których serwer w repozytorium " +
				"nie zastał; dziś moduł prowadzi wyłącznie kierunek `up`")
	}

	kroki, katalog, err := a.krokiMigracji(okno.Id, z.Target)
	if err != nil {
		return shared.DeveloperDataMigrationRunResponse{}, err
	}

	polaczenie, baza, err := a.otworzPolaczenieDanych(ctx, z.ConnectionId)
	if err != nil {
		return shared.DeveloperDataMigrationRunResponse{}, err
	}
	defer baza.Close()
	if polaczenie.TylkoOdczyt && !proba {
		return shared.DeveloperDataMigrationRunResponse{}, bladDostepuDevelopera(
			"połączenie " + polaczenie.Nazwa + " jest oznaczone jako tylko do odczytu, " +
				"a migracja zmienia schemat")
	}

	kontekst, przerwij := context.WithTimeout(ctx, czasZapytaniaDanych)
	defer przerwij()

	zastosowane, err := zastosowaneMigracjeDocelowe(kontekst, baza, proba)
	if err != nil {
		return shared.DeveloperDataMigrationRunResponse{}, bladWykonaniaDevelopera(
			"nie można odczytać rejestru migracji w bazie docelowej: " + err.Error())
	}

	wykonane := make([]string, 0, len(kroki))
	oczekujace := make([]string, 0, len(kroki))
	opis := strings.Builder{}
	for _, krok := range kroki {
		if zastosowane[krok] {
			continue
		}
		if proba {
			oczekujace = append(oczekujace, krok)
			continue
		}
		tresc, err := os.ReadFile(filepath.Join(katalog, krok))
		if err != nil {
			return shared.DeveloperDataMigrationRunResponse{}, bladZasobuDevelopera(
				"nie można odczytać migracji " + krok + ": " + err.Error())
		}
		if err := zastosujMigracjeDocelowa(kontekst, baza, krok, string(tresc)); err != nil {
			// Kolejne kroki pisano z założeniem, że ten doszedł do skutku.
			opis.WriteString(krok + ": " + err.Error() + "\n")
			oczekujace = append(oczekujace, krok)
			break
		}
		wykonane = append(wykonane, krok)
		opis.WriteString(krok + ": zastosowana\n")
	}

	odpowiedz := shared.DeveloperDataMigrationRunResponse{Applied: wykonane, Pending: oczekujace}
	if opis.Len() > 0 {
		odpowiedz.Output = wskaznikTekstu(strings.TrimRight(opis.String(), "\n"))
	}
	return odpowiedz, nil
}

func (a *adapterDevelopera) krokiMigracji(oknoKod string, cel *string) ([]string, string, error) {
	wskazanie := strings.TrimSpace(tekstWskazaniaDevelopera(cel))
	if wskazanie == "" {
		wskazanie = "migrations"
	}
	_, katalog, err := a.plikOkna(oknoKod, wskazanie)
	if err != nil {
		return nil, "", err
	}
	opis, err := os.Stat(katalog)
	if err != nil || !opis.IsDir() {
		return nil, "", bladZasobuDevelopera(
			"katalog migracji " + wskazanie + " nie istnieje w repozytorium okna")
	}
	wpisy, err := os.ReadDir(katalog)
	if err != nil {
		return nil, "", bladWykonaniaDevelopera(
			"nie można odczytać katalogu migracji: " + err.Error())
	}
	kroki := make([]string, 0, len(wpisy))
	for _, wpis := range wpisy {
		if !wpis.IsDir() && strings.EqualFold(filepath.Ext(wpis.Name()), ".sql") {
			kroki = append(kroki, wpis.Name())
		}
	}
	if len(kroki) == 0 {
		return nil, "", bladZasobuDevelopera(
			"katalog " + wskazanie + " nie zawiera ani jednego pliku .sql")
	}
	// Kolejność narzucają nazwy z numerem wiodącym, nie porządek katalogu.
	sort.Strings(kroki)
	return kroki, katalog, nil
}

// Rejestr stoi w bazie Operatora, nie w bazie produktu.
const nazwaRejestruMigracjiDocelowej = "danaco_migracja"

func zastosowaneMigracjeDocelowe(ctx context.Context, baza *sql.DB,
	proba bool) (map[string]bool, error) {

	if !proba {
		if _, err := baza.ExecContext(ctx, `CREATE TABLE IF NOT EXISTS `+
			nazwaRejestruMigracjiDocelowej+` (nazwa VARCHAR(255) PRIMARY KEY)`); err != nil {
			return nil, err
		}
	}
	wiersze, err := baza.QueryContext(ctx, `SELECT nazwa FROM `+nazwaRejestruMigracjiDocelowej)
	if err != nil {
		// Próba (`dryRun`) niczego nie zakłada; brak rejestru znaczy bazę bez zastosowanych migracji.
		if proba {
			return map[string]bool{}, nil
		}
		return nil, err
	}
	defer wiersze.Close()

	zastosowane := map[string]bool{}
	for wiersze.Next() {
		var nazwa string
		if err := wiersze.Scan(&nazwa); err != nil {
			return nil, err
		}
		zastosowane[nazwa] = true
	}
	return zastosowane, wiersze.Err()
}

// Krok i wpis do rejestru idą w jednej transakcji; krok bez wpisu wykonałby się drugi raz.
func zastosujMigracjeDocelowa(ctx context.Context, baza *sql.DB, nazwa, tresc string) error {
	transakcja, err := baza.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	if _, err := transakcja.ExecContext(ctx, tresc); err != nil {
		_ = transakcja.Rollback()
		return err
	}
	if _, err := transakcja.ExecContext(ctx,
		`INSERT INTO `+nazwaRejestruMigracjiDocelowej+` (nazwa) VALUES (`+
			znakParametru(baza)+`)`, nazwa); err != nil {
		_ = transakcja.Rollback()
		return err
	}
	return transakcja.Commit()
}

// PostgreSQL numeruje parametry (`$1`); pozostałe silniki używają znaku zapytania.
func znakParametru(baza *sql.DB) string {
	if fmt.Sprintf("%T", baza.Driver()) == "*stdlib.Driver" {
		return "$1"
	}
	return "?"
}

func (a *adapterDevelopera) otworzPolaczenieDanych(ctx context.Context,
	kod string) (dane.PolaczenieDanych, *sql.DB, error) {

	wskazanie := strings.TrimSpace(kod)
	if wskazanie == "" {
		return dane.PolaczenieDanych{}, nil, bladZadaniaDevelopera(
			"czynność wymaga wskazania połączenia bazodanowego")
	}
	if a.repozytorium == nil {
		return dane.PolaczenieDanych{}, nil, bladZasobuDevelopera(
			"serwer nie ma opisów połączeń bazodanowych")
	}
	polaczenie, err := a.repozytorium.PolaczenieDanychPoKodzie(ctx, wskazanie)
	if errors.Is(err, dane.ErrBrakWiersza) {
		return dane.PolaczenieDanych{}, nil, bladZasobuDevelopera(
			"nie ma połączenia bazodanowego " + wskazanie)
	}
	if err != nil {
		return dane.PolaczenieDanych{}, nil, bladWykonaniaDevelopera(
			"nie można odczytać połączenia " + wskazanie + ": " + err.Error())
	}

	sterownik, odwolanie, err := adresPolaczeniaDanych(polaczenie)
	if err != nil {
		return dane.PolaczenieDanych{}, nil, err
	}
	baza, err := sql.Open(sterownik, odwolanie)
	if err != nil {
		return dane.PolaczenieDanych{}, nil, bladWykonaniaDevelopera(
			"nie można otworzyć połączenia " + polaczenie.Nazwa + ": " + err.Error())
	}
	// sql.Open niczego jeszcze nie łączy; łączność sprawdza dopiero Ping.
	kontekst, przerwij := context.WithTimeout(ctx, 10*time.Second)
	defer przerwij()
	if err := baza.PingContext(kontekst); err != nil {
		baza.Close()
		return dane.PolaczenieDanych{}, nil, bladWykonaniaDevelopera(
			"połączenie " + polaczenie.Nazwa + " nie odpowiada: " + err.Error())
	}
	return polaczenie, baza, nil
}

// Połączenie bez odwołania do sejfu jedzie bez hasła (SQLite).
func adresPolaczeniaDanych(polaczenie dane.PolaczenieDanych) (string, string, error) {
	haslo := ""
	if polaczenie.Poswiadczenie != nil && strings.TrimSpace(*polaczenie.Poswiadczenie) != "" {
		haslo = os.Getenv(strings.TrimSpace(*polaczenie.Poswiadczenie))
	}
	uzytkownik := ""
	if polaczenie.Uzytkownik != nil {
		uzytkownik = *polaczenie.Uzytkownik
	}
	host := "localhost"
	if polaczenie.Host != nil && strings.TrimSpace(*polaczenie.Host) != "" {
		host = strings.TrimSpace(*polaczenie.Host)
	}

	switch polaczenie.Silnik {
	case shared.DataEnginePostgres:
		port := int64(5432)
		if polaczenie.Port != nil {
			port = *polaczenie.Port
		}
		odwolanie := "postgres://" + uwierzytelnienieAdresu(uzytkownik, haslo) + host + ":" +
			strconv.FormatInt(port, 10) + "/" + polaczenie.Baza + "?sslmode=prefer"
		return "pgx", odwolanie, nil

	case shared.DataEngineMysql:
		port := int64(3306)
		if polaczenie.Port != nil {
			port = *polaczenie.Port
		}
		poswiadczenie := uzytkownik
		if haslo != "" {
			poswiadczenie += ":" + haslo
		}
		if poswiadczenie != "" {
			poswiadczenie += "@"
		}
		odwolanie := poswiadczenie + "tcp(" + host + ":" + strconv.FormatInt(port, 10) + ")/" +
			polaczenie.Baza + "?parseTime=true"
		return "mysql", odwolanie, nil

	case shared.DataEngineSqlite:
		return "sqlite", filepath.ToSlash(polaczenie.Baza), nil

	default:
		return "", "", bladZadaniaDevelopera(
			"nieznany silnik bazy danych: " + string(polaczenie.Silnik))
	}
}

func uwierzytelnienieAdresu(uzytkownik, haslo string) string {
	if uzytkownik == "" {
		return ""
	}
	if haslo == "" {
		return uzytkownik + "@"
	}
	return uzytkownik + ":" + haslo + "@"
}
