package core

// CzyTuraWBiegu mówi, czy w oknie trwa tura odpowiedzi modelu. Wiedzę tę ma wyłącznie domena rozmowy — metoda wystawia ją stanowi okna zamiast zakładać drugi rejestr tur. Odczyt idzie pod tym samym zamkiem, co zapis.
func (a *adapterRozmowy) CzyTuraWBiegu(idOkna string) bool {
	if a == nil {
		return false
	}
	a.mu.Lock()
	defer a.mu.Unlock()
	_, biegnie := a.biegnace[idOkna]
	return biegnie
}
