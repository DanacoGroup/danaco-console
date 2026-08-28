package session

import "danacoconsole/shared"

// Rola okna rozstrzyga miejsce okna w pętli koordynator-wykonawca; wartości należą do kontraktu.

// roleDopuszczalne to zbiór ról znanych kontraktowi, wykorzystywany do sprawdzania poprawności wartości.
var roleDopuszczalne = map[shared.WindowRole]bool{
	shared.WindowRoleStandalone:  true,
	shared.WindowRoleExecutor:    true,
	shared.WindowRoleCoordinator: true,
}

// RolaDomyslna obowiązuje, gdy okno nie wskazało roli. Okno samodzielne stoi
// poza pętlą, więc jest bezpiecznym stanem wyjściowym.
const RolaDomyslna = shared.WindowRoleStandalone

// normalizujRole doprowadza rolę i powiązanie koordynatora do stanu spójnego:
// rola nieznana albo pusta staje się rolą domyślną, a wskazanie koordynatora
// zachowuje wyłącznie okno wykonawcy. Nieznana wartość nie zatrzymuje otwarcia
// okna.
func normalizujRole(u Ustawienia) Ustawienia {
	if !roleDopuszczalne[u.RolaOkna] {
		u.RolaOkna = RolaDomyslna
	}
	if u.RolaOkna != shared.WindowRoleExecutor {
		u.OknoKoordynatora = ""
	}
	return u
}

// Metoda CzyWykonawca mówi, czy zakończenie tury tego okna wybudza koordynatora nadrzędnego wobec niego.
func (o Okno) CzyWykonawca() bool {
	return o.RolaOkna == shared.WindowRoleExecutor
}

// Metoda CzyKoordynator mówi, czy okno przyjmuje wybudzenia zgłaszane przez swoich wykonawców podległych.
func (o Okno) CzyKoordynator() bool {
	return o.RolaOkna == shared.WindowRoleCoordinator
}

// Metoda CzySamodzielne mówi, czy okno stoi poza pętlą koordynator-wykonawca i pracuje całkiem niezależnie.
func (o Okno) CzySamodzielne() bool {
	return o.RolaOkna == shared.WindowRoleStandalone
}

// RolaZnana odpowiada, czy wartość roli pochodzi ze słownika kontraktu.
// Korzysta z niej warstwa wyżej, gdy chce odróżnić wartość przyjętą wprost od
// zastąpionej wartością domyślną.
func RolaZnana(r shared.WindowRole) bool {
	return roleDopuszczalne[r]
}
