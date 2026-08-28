//go:build windows

package session

// wstrzymaj zgłasza brak wsparcia platformy Windows dla wstrzymania całego drzewa procesów i pozostawia proces nietknięty, uczciwie odzwierciedlając możliwości maszyny.
func (d *drzewoProcesow) wstrzymaj() (bool, error) {
	return false, nil
}

// wznow zgłasza brak wsparcia platformy Windows dla wznowienia wstrzymanego drzewa procesów i pozostawia proces nietknięty, zgodnie z ograniczeniem systemu.
func (d *drzewoProcesow) wznow() (bool, error) {
	return false, nil
}
