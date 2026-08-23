package konfig

import "danacoconsole/shared"

// Kontekst wskazuje byt obowiązujący na każdym z dziewięciu poziomów zasięgu
// w chwili rozstrzygania. Pole puste znaczy, że poziom nie dotyczy tego
// wywołania — jest wtedy pomijany, a rozstrzyganie schodzi na poziom szerszy.
// Poziomy globalny i aplikacja obowiązują zawsze i nie mają bytu, więc ich klucz
// zasięgu jest pusty: platformy ani programu nie ma czym zawęzić.
type Kontekst struct {
	Srodowisko  string // srodowisko.id
	Modul       string // modul.id
	ParaModulow string // para_modulow.id
	Projekt     string // projekt.id
	KartaSesji  string // karta_sesji.id
	Rola        string // rola_multitasking.id
	Okno        string // okno_komunikacji.id — poziom najwęższy

	// Model i Konto wskazują byt osi rozstrzygania. Pole puste znaczy, że oś
	// nie dotyczy tego wywołania — jest wtedy pomijana, a rozstrzyganie schodzi
	// na oś platformy. Oś platformy obowiązuje zawsze i nie ma bytu, więc jej
	// klucz jest pusty.
	Model string // kanal_modelu.kod albo identyfikator modelu
	Konto string // konto.id
}

// Adres wskazuje jedno miejsce zapisu ustawienia: poziom zasięgu wraz z bytem
// tego poziomu oraz oś wraz z bytem osi. Odpowiada czwórce kolumn
// ustawienie.poziom_zasiegu_id + klucz_zasiegu + os + klucz_osi.
type Adres struct {
	Poziom       Poziom
	KluczZasiegu string
	Os           Os
	KluczOsi     string
}

// Adres zwraca klucz zasięgu dla poziomu oraz informację, czy poziom obowiązuje
// w tym kontekście.
func (k Kontekst) Adres(poziom Poziom) (string, bool) {
	switch poziom {
	case shared.ConfigScopeApplication, shared.ConfigScopeGlobal:
		return "", true
	case shared.ConfigScopeEnvironment:
		return k.Srodowisko, k.Srodowisko != ""
	case shared.ConfigScopeModule:
		return k.Modul, k.Modul != ""
	case shared.ConfigScopeModulePair:
		return k.ParaModulow, k.ParaModulow != ""
	case shared.ConfigScopeProject:
		return k.Projekt, k.Projekt != ""
	case shared.ConfigScopeSession:
		return k.KartaSesji, k.KartaSesji != ""
	case shared.ConfigScopeRole:
		return k.Rola, k.Rola != ""
	case shared.ConfigScopeWindow:
		return k.Okno, k.Okno != ""
	}
	return "", false
}

// adresOsi zwraca byt osi oraz informację, czy oś obowiązuje w tym kontekście.
// Oś platformy obowiązuje zawsze i nie ma bytu.
func (k Kontekst) adresOsi(os Os) (string, bool) {
	switch OsLubPlatforma(os) {
	case OsPlatformy:
		return "", true
	case OsModelu:
		return k.Model, k.Model != ""
	case OsKonta:
		return k.Konto, k.Konto != ""
	}
	return "", false
}

// Adresy zwraca miejsca zapisu obowiązujące w kontekście, uporządkowane od
// najwęższego do najszerszego. To jest dokładna kolejność rozstrzygania:
// wygrywa pierwszy adres, pod którym ustawienie jest zapisane.
//
// Porządek jest dwupoziomowy: najpierw poziom zasięgu, a w ramach poziomu —
// oś: konto, model, platforma. Kontekst bez modelu i bez konta daje dokładnie
// osiem adresów osi platformy.
func (k Kontekst) Adresy() []Adres {
	adresy := make([]Adres, 0, len(poziomyOdNajwezszego)*len(osieOdNajwezszej))
	for _, poziom := range poziomyOdNajwezszego {
		klucz, obowiazuje := k.Adres(poziom)
		if !obowiazuje {
			continue
		}
		for _, os := range osieOdNajwezszej {
			kluczOsi, dotyczy := k.adresOsi(os)
			if !dotyczy {
				continue
			}
			adresy = append(adresy, Adres{
				Poziom: poziom, KluczZasiegu: klucz, Os: os, KluczOsi: kluczOsi,
			})
		}
	}
	return adresy
}
