package session

// Wstrzymaj zatrzymuje przejęty proces wraz z całym jego potomstwem i zwraca informację, czy platforma zna tę czynność, przy czym fałsz oznacza proces nietknięty, a nie niepowodzenie.
func (d *DrzewoProcesu) Wstrzymaj() (bool, error) {
	if d == nil {
		return false, nil
	}
	return d.drzewo.wstrzymaj()
}

// Wznow podejmuje pracę procesu wstrzymanego wcześniej funkcją Wstrzymaj i zwraca informację o wsparciu platformy na tych samych zasadach co ona.
func (d *DrzewoProcesu) Wznow() (bool, error) {
	if d == nil {
		return false, nil
	}
	return d.drzewo.wznow()
}
