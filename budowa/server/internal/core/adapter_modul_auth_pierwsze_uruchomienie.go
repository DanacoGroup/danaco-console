// Odpowiedzialność pliku: znacznik bramki postawionej BEZ poczty — jedyna
// pamięć o tym, że pierwsze konto Właściciela powstało na świeżej instalacji,
// w której konta nadawczego nie było czym wskazać.
//
// ── DLACZEGO TEN ZNACZNIK ISTNIEJE ──────────────────────────────────────────
// Klient uruchamia okno i pierwsze, co widzi, to rejestracja. Świeża instalka
// nie ma konta nadawczego platformy (`mailer.host`, `mailer.address` bez
// wartości domyślnej), więc dopóki założenie bramki wymagało poczty, pierwszego
// konta nie dawało się założyć w ogóle: „Załóż konto" oddawało
// `internal_error`, a odmowa odsyłała do okna Konfiguracji, do którego bez
// konta nie sposób wejść. Poczta jest potrzebna do pisania do INNYCH ludzi —
// do potwierdzania ich adresów i do odzyskiwania hasła listem — a nie do
// postawienia bramki na własnym urządzeniu.
//
// ── CO ZNACZNIK MÓWI, A CZEGO NIE ───────────────────────────────────────────
// Mówi: „konto założono, listu nie było komu nadać, adres pozostaje
// niepotwierdzony". Nie mówi, że adres jest w porządku — dlatego konto zostaje
// w bazie NIEPOTWIERDZONE i nikt tego stanu nie udaje. Znacznik zdejmuje
// wyłącznie jeden warunek: bramki nie zamyka brak potwierdzenia, którego
// platforma nie miała czym wysłać. Bramkę nadal otwiera hasło i tylko hasło.
//
// ── DLACZEGO W SEJFIE, A NIE W BAZIE ────────────────────────────────────────
// Tabela `potwierdzenie_tozsamosci` przyjmuje dwa cele i tylko dwa
// (`CHECK (cel IN ('weryfikacja', 'odzyskanie'))`), a `konto_wlasciciela` ma
// jedną kolumnę stanu — trzeciego stanu nie ma gdzie zapisać bez zmiany
// schematu, którego rdzeń nie jest w tej chwili właścicielem. Sejf poświadczeń
// jest tym samym magazynem trwałym, w którym leży już sekret kotwicy, i chodzi
// do niego ten sam przedrostek `auth:` — więc pamięć bramki zostaje w jednym
// miejscu, a nie w dwóch.
package core

import (
	"context"

	"danacoconsole/shared"
)

// bytZnacznikaBezPoczty jest kluczem wpisu w sejfie. Nazwa mówi wprost, czego
// zabrakło, bo wpis czyta człowiek zaglądający do sejfu równie często jak rdzeń.
const bytZnacznikaBezPoczty = przedrostekBytuSejfu + "bramka-bez-poczty"

// zapiszZnacznikBezPoczty zapamiętuje adres, którego nikt nie potwierdził.
//
// Treścią wpisu jest sam adres — po to, żeby odmowy i przyszłe potwierdzenie
// mogły powiedzieć, O KTÓRY adres chodzi, bez sięgania po wiersz konta.
// Niepowodzenie zapisu jest odmową rejestracji, nie ciszą: bramka bez tego
// wpisu byłaby kontem, którego hasło działa, a bramka i tak nie wpuszcza.
func (a *adapterUwierzytelnienia) zapiszZnacznikBezPoczty(ctx context.Context, email string) error {
	if a.sejf == nil {
		return bladBramki(shared.ErrorCodeInternalError,
			"sejfu poświadczeń nie wpięto; stanu bramki bez poczty nie ma gdzie zapisać")
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
