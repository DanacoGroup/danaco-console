// Komendy `terminal.host.save`, `terminal.host.list` i `terminal.host.remove`
// — książka hostów okna Session Manager. Wpis hosta jest nastawą Operatora,
// nie stanem widoku — ma przeżyć zamknięcie okna i restart rdzenia.
package core

import (
	"context"
	"errors"
	"strings"

	"danacoconsole/server/internal/dane"
	"danacoconsole/server/internal/protocol"
	"danacoconsole/shared"
)

// przedrostekHosta znakuje identyfikator wpisu książki hostów okna Session
// Managera terminala rdzenia.
const przedrostekHosta = "thost-"

// ZapiszHosta obsługuje `terminal.host.save`, zapisując albo zmieniając wpis
// w książce hostów okna terminala.
func (a *adapterTerminala) ZapiszHosta(ctx context.Context,
	z shared.TerminalHostSaveRequest) (shared.TerminalHostSaveResponse, error) {

	dziennik, err := a.dziennikWyposazenia()
	if err != nil {
		return shared.TerminalHostSaveResponse{}, err
	}
	nazwa := strings.TrimSpace(z.Host.Name)
	cel := strings.TrimSpace(z.Host.Target)
	if nazwa == "" || cel == "" {
		return shared.TerminalHostSaveResponse{}, bladZadaniaTerminala(
			"wpis książki hostów wymaga nazwy i adresu celu")
	}
	if z.Host.Port != nil && (*z.Host.Port < 1 || *z.Host.Port > 65535) {
		return shared.TerminalHostSaveResponse{}, bladZadaniaTerminala(
			"port połączenia leży poza zakresem portów (1–65535)")
	}

	kod := strings.TrimSpace(z.Host.Id)
	powstal := kod == ""
	if powstal {
		kod = nowyIdentyfikator(przedrostekHosta)
	} else if _, err := dziennik.Host(ctx, kod); err != nil {
		if !errors.Is(err, dane.ErrBrakWiersza) {
			return shared.TerminalHostSaveResponse{}, err
		}
		// Wpis wskazany identyfikatorem, którego nie ma, powstaje pod tym
		// identyfikatorem.
		powstal = true
	}

	wiersz := dane.HostTerminala{
		Kod:             kod,
		Nazwa:           nazwa,
		Cel:             cel,
		Grupa:           strings.TrimSpace(wartoscTekstu(z.Host.Group)),
		KatalogRoboczy:  strings.TrimSpace(wartoscTekstu(z.Host.WorkingDir)),
		KluczKod:        wskaznikTekstu(strings.TrimSpace(wartoscTekstu(z.Host.KeyId))),
		HostPosredniKod: wskaznikTekstu(strings.TrimSpace(wartoscTekstu(z.Host.JumpHostId))),
		Notatka:         wartoscTekstu(z.Host.Note),
	}
	if z.Host.Port != nil {
		port := int64(*z.Host.Port)
		wiersz.Port = &port
	}
	if err := dziennik.ZapiszHosta(ctx, wiersz); err != nil {
		return shared.TerminalHostSaveResponse{}, err
	}
	// Odczyt po zapisie, a nie odbicie żądania: znaczniki czasu nadaje baza.
	zapisany, err := dziennik.Host(ctx, kod)
	if err != nil {
		return shared.TerminalHostSaveResponse{}, err
	}
	return shared.TerminalHostSaveResponse{Host: hostKontraktu(zapisany), Created: powstal}, nil
}

// WykazHostow obsługuje `terminal.host.list`, oddając całą książkę hostów okna
// wskazanego w żądaniu komendy.
func (a *adapterTerminala) WykazHostow(ctx context.Context,
	z shared.TerminalHostListRequest) (shared.TerminalHostListResponse, error) {

	dziennik, err := a.dziennikWyposazenia()
	if err != nil {
		return shared.TerminalHostListResponse{}, err
	}
	wiersze, err := dziennik.Hosty(ctx, dane.FiltrHostowTerminala{
		Grupa: strings.TrimSpace(wartoscTekstu(z.Group)),
		Fraza: strings.TrimSpace(wartoscTekstu(z.Query)),
	})
	if err != nil {
		return shared.TerminalHostListResponse{}, err
	}
	wykaz := make([]shared.TerminalHost, 0, len(wiersze))
	for _, wiersz := range wiersze {
		wykaz = append(wykaz, hostKontraktu(wiersz))
	}
	return shared.TerminalHostListResponse{Hosts: wykaz, Total: len(wykaz)}, nil
}

// UsunHosta obsługuje `terminal.host.remove`.
//
// Karty już otwarte do tego hosta biegną dalej — tak mówi kontrakt — bo adres
// przeszedł do karty przy jej otwarciu i nie jest odczytywany z książki przy
// każdym poleceniu.
func (a *adapterTerminala) UsunHosta(ctx context.Context,
	z shared.TerminalHostRemoveRequest) (shared.TerminalHostRemoveResponse, error) {

	dziennik, err := a.dziennikWyposazenia()
	if err != nil {
		return shared.TerminalHostRemoveResponse{}, err
	}
	kod := strings.TrimSpace(z.HostId)
	if kod == "" {
		return shared.TerminalHostRemoveResponse{}, bladZadaniaTerminala(
			"usunięcie wpisu wymaga jego wskazania")
	}
	usuniety, err := dziennik.UsunHosta(ctx, kod)
	if err != nil {
		return shared.TerminalHostRemoveResponse{}, err
	}
	return shared.TerminalHostRemoveResponse{Removed: usuniety}, nil
}

// hostKontraktu przekłada wiersz książki hostów z bazy danych na byt kontraktu
// widoczny oknu Session Manager.
func hostKontraktu(w dane.HostTerminala) shared.TerminalHost {
	wpis := shared.TerminalHost{
		Id:         w.Kod,
		Name:       w.Nazwa,
		Target:     w.Cel,
		Group:      wskaznikTekstu(w.Grupa),
		WorkingDir: wskaznikTekstu(w.KatalogRoboczy),
		KeyId:      w.KluczKod,
		JumpHostId: w.HostPosredniKod,
		Note:       wskaznikTekstu(w.Notatka),
		CreatedAt:  chwilaBazy(w.Utworzono),
		UpdatedAt:  chwilaBazy(w.Zaktualizowano),
	}
	wpis.Port = wskaznikMalej(w.Port)
	return wpis
}

// dziennikWyposazenia oddaje repozytorium modułu albo odmowę nazywającą powód.
// Wyposażenie Terminala jest z natury trwałe: jego wartość polega na tym, że
// przeżywa posiedzenie.
func (a *adapterTerminala) dziennikWyposazenia() (dane.RepozytoriumTerminala, error) {
	if a.repozytorium == nil {
		return nil, protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeInternalError,
			"moduł Terminal: rdzeń pracuje bez dziennika, więc nie prowadzi książki hostów, "+
				"biblioteki skryptów ani wykazu kluczy"))
	}
	return a.repozytorium, nil
}
