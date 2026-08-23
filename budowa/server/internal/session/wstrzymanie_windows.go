//go:build windows

package session

// Windows nie ma dla obcego procesu odpowiednika sygnałów SIGSTOP i SIGCONT.
// Wstrzymanie idzie tam per WĄTEK (SuspendThread), a Job Object — uchwyt, którym
// rdzeń obejmuje całe drzewo — takiej czynności nie zna. Przejście po wątkach
// wszystkich procesów drzewa nie jest tu odpowiednikiem: wątek utworzony między
// przebiegiem a końcem czynności zostałby biegnący, więc „wstrzymane" znaczyłoby
// coś innego niż na systemach uniksowych.
//
// Dlatego czynność zgłasza się jako niewspierana, zamiast robić coś podobnego.
// Kontrakt komendy `terminal.process.suspend` ma na to pole `supported`, a proces
// zostaje NIETKNIĘTY — Operator dostaje prawdę o możliwościach maszyny, a nie
// odpowiedź udaną, po której proces dalej zajmuje procesor.

// wstrzymaj zgłasza brak wsparcia platformy; proces zostaje nietknięty.
func (d *drzewoProcesow) wstrzymaj() (bool, error) {
	return false, nil
}

// wznow zgłasza brak wsparcia platformy; proces zostaje nietknięty.
func (d *drzewoProcesow) wznow() (bool, error) {
	return false, nil
}
