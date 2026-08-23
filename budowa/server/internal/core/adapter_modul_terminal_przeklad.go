// Odpowiedzialność pliku: przełożenie bytów modułu Terminal na kształt
// kontraktu oraz zawężenie wykazu procesów filtrem żądania.
//
// Filtr działa tak samo na stanie żywym i na dzienniku. Wykaz procesów składa
// się z dwóch źródeł, więc gdyby każde zawężało po swojemu, ten sam filtr dałby
// niespójny wynik. Dziennik dostaje filtr wprost w zapytaniu, a stan żywy — tę
// samą czwórkę warunków sprawdzoną na widoku kontraktu.
package core

import (
	"danacoconsole/server/internal/dane"
	"danacoconsole/shared"
)

// kartaKontraktu przekłada kartę powłoki na byt kontraktu.
//
// Pola `pid` i `exitCode` zostają puste, bo karta nie jest procesem: procesem
// jest każde wykonane w niej polecenie i to ono ma PID oraz kod wyjścia
// (adapter_modul_terminal_powloki.go). Oba pola są w kontrakcie opcjonalne.
func kartaKontraktu(k *kartaTerminala) shared.TerminalSession {
	return shared.TerminalSession{
		Id:              k.kod,
		WindowId:        k.oknoKod,
		Shell:           k.powloka,
		Title:           wskaznikTekstu(k.tytul),
		WorkingDir:      wskaznikTekstu(k.katalog),
		Status:          k.stan,
		CreatedAt:       k.utworzono.UnixMilli(),
		ParentSessionId: wskaznikTekstu(k.idSesji),
		RemoteTarget:    wskaznikTekstu(k.celZdalny),
	}
}

// kartaWierszaKontraktu przekłada wiersz dziennika kart na byt kontraktu. Na
// tym stoi ta część `terminal.session.list`, która sięga po karty zamknięte —
// rejestr pamięci ich nie trzyma po restarcie rdzenia.
func kartaWierszaKontraktu(w dane.KartaTerminala) shared.TerminalSession {
	return shared.TerminalSession{
		Id:           w.Kod,
		WindowId:     w.OknoKod,
		Shell:        w.Powloka,
		Title:        w.Tytul,
		WorkingDir:   w.KatalogRoboczy,
		Status:       w.Stan,
		CreatedAt:    chwilaBazy(w.Utworzono),
		RemoteTarget: wskaznikTekstu(w.CelZdalny),
	}
}

// procesKontraktu przekłada proces rejestru żywego na byt kontraktu.
//
// Zużycia procesora i pamięci nie wypełniamy. Kontrakt ma na nie pola
// opcjonalne, lecz nie ma tu miernika zużycia obcego procesu — wpisanie zera
// znaczyłoby „nic nie zużywa" zamiast „brak danych".
func procesKontraktu(p *procesTerminala) shared.TerminalProcess {
	stan, kodWyjscia, zakonczono := p.Migawka()
	widok := shared.TerminalProcess{
		Id:        p.kod,
		SessionId: wskaznikTekstu(p.kartaKod),
		WindowId:  p.oknoKod,
		Command:   p.polecenie,
		Initiator: p.inicjator,
		Status:    stan,
		StartedAt: p.uruchomiono.UnixMilli(),
	}
	if p.pid > 0 {
		pid := p.pid
		widok.Pid = &pid
	}
	if p.pidNadrzedny > 0 {
		rodzic := p.pidNadrzedny
		widok.ParentPid = &rodzic
	}
	widok.ExitCode = kodWyjscia
	if !zakonczono.IsZero() {
		chwila := zakonczono.UnixMilli()
		widok.FinishedAt = &chwila
	}
	return widok
}

// procesWierszaKontraktu przekłada wiersz dziennika na byt kontraktu.
func procesWierszaKontraktu(w dane.ProcesTerminala) shared.TerminalProcess {
	widok := shared.TerminalProcess{
		Id:        w.Kod,
		SessionId: w.KartaKod,
		WindowId:  w.OknoKod,
		Command:   w.Polecenie,
		Initiator: w.Inicjator,
		Status:    w.Stan,
		Pid:       wskaznikMalej(w.Pid),
		ParentPid: wskaznikMalej(w.PidNadrzedny),
		ExitCode:  wskaznikMalej(w.KodWyjscia),
		StartedAt: chwilaBazy(w.Uruchomiono),
	}
	if w.Zakonczono != nil {
		chwila := chwilaBazy(*w.Zakonczono)
		widok.FinishedAt = &chwila
	}
	return widok
}

// filtrProcesow czyta zawężenia z żądania kontraktu.
func filtrProcesow(z shared.TerminalProcessListRequest) dane.FiltrProcesow {
	filtr := dane.FiltrProcesow{
		OknoKod:  wartoscTekstu(z.WindowId),
		KartaKod: wartoscTekstu(z.SessionId),
	}
	if z.Status != nil {
		filtr.Stan = *z.Status
	}
	if z.Initiator != nil {
		filtr.Inicjator = *z.Initiator
	}
	return filtr
}

// filtrPrzepuszcza sprawdza widok kontraktu tą samą czwórką warunków, którą
// dziennik dostaje w zapytaniu.
func filtrPrzepuszcza(f dane.FiltrProcesow, p shared.TerminalProcess) bool {
	if f.OknoKod != "" && p.WindowId != f.OknoKod {
		return false
	}
	if f.KartaKod != "" && wartoscTekstu(p.SessionId) != f.KartaKod {
		return false
	}
	if f.Stan != "" && p.Status != f.Stan {
		return false
	}
	return f.Inicjator == "" || p.Initiator == f.Inicjator
}

// wskaznikMalej zwęża liczbę bazy do rozmiaru pola kontraktu.
func wskaznikMalej(wartosc *int64) *int {
	if wartosc == nil {
		return nil
	}
	mniejsza := int(*wartosc)
	return &mniejsza
}

// wskaznikDuzej rozszerza liczbę kontraktu do rozmiaru kolumny bazy.
func wskaznikDuzej(wartosc *int) *int64 {
	if wartosc == nil {
		return nil
	}
	wieksza := int64(*wartosc)
	return &wieksza
}
