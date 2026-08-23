package core

// CzyTuraWBiegu mówi, czy w oknie trwa tura odpowiedzi modelu.
//
// Wiedzę tę ma wyłącznie domena rozmowy: tura jest wpisem w wykazie biegów
// adaptera, zakładanym przy `message.send` i usuwanym po domknięciu strumienia.
// Metoda wystawia ją stanowi okna (`window.state.get`, pole `streaming`),
// zamiast zakładać drugi rejestr tur obok istniejącego.
//
// Odczyt idzie pod tym samym zamkiem, co zapis — pytanie o stan może przyjść
// w trakcie tury, z innego połączenia tego samego konta.
func (a *adapterRozmowy) CzyTuraWBiegu(idOkna string) bool {
	if a == nil {
		return false
	}
	a.mu.Lock()
	defer a.mu.Unlock()
	_, biegnie := a.biegnace[idOkna]
	return biegnie
}
