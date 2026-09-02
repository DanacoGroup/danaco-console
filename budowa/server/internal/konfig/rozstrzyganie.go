package konfig

// Pochodzenie mówi, skąd wzięła się wartość rozstrzygnięta: z zapisu,
// z wartości domyślnej albo z braku definicji i zapisu.
type Pochodzenie string

// Wartości typu Pochodzenie, obejmujące pochodzenie zapisu, wartości
// domyślnej oraz braku definicji i zapisu.
const (
	// Wartość zapisana na poziomie zasięgu wskazanym w wyniku rozstrzygnięcia;
	// jest to najczęstsze pochodzenie zwracanej wartości.
	PochodzenieZapis Pochodzenie = "zapis"
	// Wartość domyślna z rejestru definicji, zwracana przy braku zapisu na
	// którymkolwiek z ośmiu poziomów zasięgu.
	pochodzenieDomyslna Pochodzenie = "domyslna"
	// Klucz bez definicji i bez zapisu w żadnym z poziomów zasięgu; odczyt
	// nie jest błędem, zwracana jest wartość pusta.
	PochodzenieNieznane Pochodzenie = "nieznane"
)

// Wynik to rozstrzygnięcie jednego ustawienia wraz ze wskazaniem źródła,
// z którego pochodzi zwrócona wartość.
type Wynik struct {
	Klucz        string
	Wartosc      string
	Rodzaj       Rodzaj
	Poziom       Poziom      // poziom, z którego pochodzi wartość; pusty dla domyślnej
	KluczZasiegu string      // byt tego poziomu
	Os           Os          // oś, z której pochodzi wartość; pusta dla domyślnej
	KluczOsi     string      // byt tej osi
	Pochodzenie  Pochodzenie // zapis · domyślna · nieznane
	Objasnienie  string      // opis z rejestru definicji
}

// Rozstrzygacz rozstrzyga ustawienia po dziewięciu poziomach zasięgu,
// sięgając po zapis, a w jego braku po wartość domyślną.
type Rozstrzygacz struct {
	zrodlo  Zrodlo
	rejestr *Rejestr

	// rozglos jest drugą stroną drogi Rozstrzygnij, zawiadamiającą
	// nasłuchujących w rozgloszenie.go.
	rozglos rozglosnia
}

// Nowy buduje rozstrzygacz. Brak źródła daje puste źródło pamięciowe, brak
// rejestru — rejestr wbudowany. Rozstrzygacz powstaje zawsze; niekompletne
// zależności dają politykę domyślną, nie odmowę startu.
func Nowy(zrodlo Zrodlo, rejestr *Rejestr) *Rozstrzygacz {
	if zrodlo == nil {
		zrodlo = NoweZrodloPamieciowe()
	}
	if rejestr == nil {
		rejestr = RejestrWbudowany()
	}
	return &Rozstrzygacz{zrodlo: zrodlo, rejestr: rejestr}
}

// Rejestr zwraca rejestr definicji użyty przez rozstrzygacz do wyznaczenia
// wartości domyślnych ustawień.
func (r *Rozstrzygacz) Rejestr() *Rejestr {
	if r == nil {
		return nil
	}
	return r.rejestr
}

// Rozstrzygnij zwraca wartość ustawienia obowiązującą w kontekście oraz poziom,
// z którego ta wartość pochodzi. Rozstrzyganie idzie od poziomu najwęższego do
// najszerszego; wygrywa pierwszy poziom, na którym ustawienie jest zapisane.
func (r *Rozstrzygacz) Rozstrzygnij(kontekst Kontekst, klucz string) Wynik {
	zapisy, _ := r.zapisy(kontekst)
	return r.rozstrzygnijZZapisow(kontekst, klucz, zapisy)
}

// zapisy pobiera zapisane ustawienia dla adresów kontekstu. Błąd źródła nie
// zatrzymuje rozstrzygania: wpisy odczytane mimo błędu wchodzą do
// rozstrzygnięcia, a ustawienia bez zapisu schodzą na wartość domyślną.
func (r *Rozstrzygacz) zapisy(kontekst Kontekst) (map[kluczWpisu]Wpis, error) {
	indeks := make(map[kluczWpisu]Wpis)
	if r == nil || r.zrodlo == nil {
		return indeks, nil
	}
	wpisy, err := r.zrodlo.Wpisy(kontekst.KontoOperatora, kontekst.Adresy())
	for _, wpis := range wpisy {
		if wpis.Klucz == "" {
			continue
		}
		indeks[kluczem(wpis.Adres(), wpis.Klucz)] = wpis
	}
	return indeks, err
}

// rozstrzygnijZZapisow wykonuje samą regułę pierwszeństwa najwęższego na już
// pobranym zbiorze zapisów. Dzięki temu podgląd polityki efektywnej rozstrzyga
// komplet ustawień na jednym odczycie źródła.
func (r *Rozstrzygacz) rozstrzygnijZZapisow(kontekst Kontekst, klucz string, zapisy map[kluczWpisu]Wpis) Wynik {
	for _, adres := range kontekst.Adresy() {
		wpis, jest := zapisy[kluczem(adres, klucz)]
		if !jest {
			continue
		}
		return Wynik{
			Klucz:        klucz,
			Wartosc:      wpis.Wartosc,
			Rodzaj:       RodzajLubTekst(wpis.Rodzaj),
			Poziom:       adres.Poziom,
			KluczZasiegu: adres.KluczZasiegu,
			Os:           OsLubPlatforma(adres.Os),
			KluczOsi:     adres.KluczOsi,
			Pochodzenie:  PochodzenieZapis,
			Objasnienie:  r.objasnienie(klucz),
		}
	}
	return r.wartoscDomyslna(klucz)
}

// wartoscDomyslna zwraca wartość z rejestru definicji, a dla klucza spoza
// rejestru — wartość pustą oznaczoną jako nieznana. Żaden z tych przypadków
// nie jest błędem.
func (r *Rozstrzygacz) wartoscDomyslna(klucz string) Wynik {
	if definicja, jest := r.Rejestr().Definicja(klucz); jest {
		return Wynik{
			Klucz:       klucz,
			Wartosc:     definicja.Domyslna,
			Rodzaj:      definicja.Rodzaj,
			Poziom:      PoziomBrak,
			Pochodzenie: pochodzenieDomyslna,
			Objasnienie: definicja.Objasnienie,
		}
	}
	return Wynik{
		Klucz:       klucz,
		Rodzaj:      RodzajTekst,
		Poziom:      PoziomBrak,
		Pochodzenie: PochodzenieNieznane,
	}
}

func (r *Rozstrzygacz) objasnienie(klucz string) string {
	definicja, _ := r.Rejestr().Definicja(klucz)
	return definicja.Objasnienie
}
