// Odpowiedzialność pliku: znacznik bramki postawionej bez poczty — pamięć
// o koncie założonym tam, gdzie konta nadawczego nie było czym wskazać.
package core

import (
	"context"
	"strconv"

	"danacoconsole/shared"
)

// bytZnacznikaBezPoczty jest kluczem wpisu zastanego, sprzed rozdzielenia
// znacznika na konta. Powstać mógł wyłącznie na instalacji jednokontowej.
const bytZnacznikaBezPoczty = przedrostekBytuSejfu + "bramka-bez-poczty"

// bytZnacznikaKonta składa klucz wpisu konta; konto zerowe klucza nie ma.
func bytZnacznikaKonta(kontoId int64) string {
	if kontoId == 0 {
		return ""
	}
	return bytZnacznikaBezPoczty + ":" + strconv.FormatInt(kontoId, 10)
}

// zapiszZnacznikBezPoczty zapamiętuje adres, którego nikt nie potwierdził;
// treścią wpisu jest sam adres, żeby odmowa wskazała go bez wiersza konta.
func (a *adapterUwierzytelnienia) zapiszZnacznikBezPoczty(ctx context.Context,
	kontoId int64, email string) error {

	if a.sejf == nil {
		return bladBramki(shared.ErrorCodeInternalError,
			"Magazyn haseł jest niedostępny.")
	}
	byt := bytZnacznikaKonta(kontoId)
	if byt == "" {
		return bladBramki(shared.ErrorCodeInternalError,
			"Znacznik bramki nie ma wskazanego konta.")
	}
	if _, err := a.sejf.Zapisz(ctx, byt, email); err != nil {
		return err
	}
	return nil
}

// znacznikBezPoczty oddaje adres czekający na potwierdzenie w tym koncie.
func (a *adapterUwierzytelnienia) znacznikBezPoczty(ctx context.Context,
	kontoId int64) (string, bool) {

	if a == nil || a.sejf == nil {
		return "", false
	}
	byt := bytZnacznikaKonta(kontoId)
	if byt == "" {
		return "", false
	}
	if email, jest := a.sejf.Odczytaj(ctx, byt); jest {
		return email, true
	}
	return a.przejmijZnacznikZastany(ctx, kontoId, byt)
}

// przejmijZnacznikZastany przenosi wpis zastany pod klucz konta najstarszego
// i kasuje go. Czytany jako zapasowy dla każdego konta otwierałby bramkę kontu,
// do którego list z kodem wyszedł. Konta nieodczytanego nie zastępuje wtedy
// najstarsze: bramka otwarta z niepewności jest otwarta dla każdego.
func (a *adapterUwierzytelnienia) przejmijZnacznikZastany(ctx context.Context,
	kontoId int64, byt string) (string, bool) {

	if a.konto == nil {
		return "", false
	}
	najstarsze, err := a.konto.Konto(ctx)
	if err != nil {
		if dziennik := dziennikZKontekstu(ctx); dziennik != nil {
			dziennik.Printf("bramka: konto najstarsze nieodczytane przy znaczniku bez poczty: %v", err)
		}
		return "", false
	}
	if najstarsze.Id != kontoId {
		return "", false
	}
	email, jest := a.sejf.Odczytaj(ctx, bytZnacznikaBezPoczty)
	if !jest {
		return "", false
	}
	if _, err := a.sejf.Zapisz(ctx, byt, email); err != nil {
		return email, true
	}
	_ = a.sejf.Usun(ctx, bytZnacznikaBezPoczty)
	return email, true
}

// zdejmijZnacznikBezPoczty kasuje wpis konta; niepowodzenie nie ma komu wrócić.
func (a *adapterUwierzytelnienia) zdejmijZnacznikBezPoczty(ctx context.Context, kontoId int64) {
	if a == nil || a.sejf == nil {
		return
	}
	if byt := bytZnacznikaKonta(kontoId); byt != "" {
		_ = a.sejf.Usun(ctx, byt)
	}
}
