// Odpowiedzialność pliku: znacznik bramki postawionej bez poczty — jedyna
// pamięć o tym, że pierwsze konto Właściciela powstało na świeżej instalacji,
// w której konta nadawczego nie było czym wskazać.
package core

import (
	"context"

	"danacoconsole/shared"
)

// bytZnacznikaBezPoczty jest kluczem wpisu w sejfie. Nazwa mówi wprost, czego
// zabrakło, bo wpis czyta człowiek zaglądający do sejfu równie często jak rdzeń.
const bytZnacznikaBezPoczty = przedrostekBytuSejfu + "bramka-bez-poczty"

// zapiszZnacznikBezPoczty zapamiętuje adres, którego nikt nie potwierdził.
// Treścią wpisu jest sam adres, żeby odmowy i przyszłe potwierdzenie mogły
// powiedzieć, o który adres chodzi, bez sięgania po wiersz konta.
func (a *adapterUwierzytelnienia) zapiszZnacznikBezPoczty(ctx context.Context, email string) error {
	if a.sejf == nil {
		return bladBramki(shared.ErrorCodeInternalError,
			"Magazyn haseł jest niedostępny.")
	}
	if _, err := a.sejf.Zapisz(ctx, bytZnacznikaBezPoczty, email); err != nil {
		return err
	}
	return nil
}

// znacznikBezPoczty mówi, czy bramkę postawiono bez poczty, i oddaje adres
// czekający na potwierdzenie. Brak wpisu to odpowiedź „nie", nie awaria.
func (a *adapterUwierzytelnienia) znacznikBezPoczty(ctx context.Context) (string, bool) {
	if a == nil || a.sejf == nil {
		return "", false
	}
	email, jest := a.sejf.Odczytaj(ctx, bytZnacznikaBezPoczty)
	if !jest {
		return "", false
	}
	return email, true
}

// zdejmijZnacznikBezPoczty kasuje wpis, gdy przestał być prawdą: adres został
// potwierdzony albo rejestracja się cofnęła. Niepowodzenie kasowania nie ma
// komu wrócić — czynność wołająca zdała już sprawę ze swojego skutku.
func (a *adapterUwierzytelnienia) zdejmijZnacznikBezPoczty(ctx context.Context) {
	if a == nil || a.sejf == nil {
		return
	}
	_ = a.sejf.Usun(ctx, bytZnacznikaBezPoczty)
}
