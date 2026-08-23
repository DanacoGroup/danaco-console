package session

// Wstrzymanie i wznowienie drzewa procesu.
//
// ── Po co, skoro jest ubicie ─────────────────────────────────────────────────
// Rdzeń umiał proces wyłącznie ZAKOŃCZYĆ. Zadanie długie — budowanie, wielka
// migracja, transkodowanie — dało się więc tylko ubić, czyli wyrzucić wykonaną
// pracę. Wstrzymanie oddaje maszynę bez utraty postępu: proces przestaje
// dostawać czas procesora, jego pamięć i otwarte pliki zostają, a wznowienie
// podejmuje pracę w miejscu, w którym stanęła.
//
// ── Dlaczego to nie jest czynność wszystkich systemów ────────────────────────
// Na systemach uniksowych są do tego sygnały SIGSTOP i SIGCONT wysyłane całej
// grupie procesów — dokładnie tej samej, którą obejmuje ubicie. Windows nie ma
// dla obcego procesu odpowiednika: wstrzymanie idzie tam per WĄTEK, a Job Object
// takiej czynności nie zna. Kontrakt komendy `terminal.process.suspend` ma
// dlatego pole `supported`, a nie tylko odpowiedź udaną albo nieudaną —
// „system tego nie umie" to co innego niż „nie udało się".
//
// Rozstrzygnięcie stoi w plikach zależnych od platformy
// (`wstrzymanie_unix.go`, `wstrzymanie_windows.go`), tak samo jak ubicie.

// Wstrzymaj zatrzymuje przejęty proces wraz z całym jego potomstwem. Drugi
// zwracany wynik mówi, czy platforma tę czynność w ogóle zna; fałsz znaczy, że
// proces został NIETKNIĘTY, a nie że czynność zawiodła.
func (d *DrzewoProcesu) Wstrzymaj() (bool, error) {
	if d == nil {
		return false, nil
	}
	return d.drzewo.wstrzymaj()
}

// Wznow podejmuje pracę procesu wstrzymanego. Wynik czyta się tak samo jak
// w Wstrzymaj.
func (d *DrzewoProcesu) Wznow() (bool, error) {
	if d == nil {
		return false, nil
	}
	return d.drzewo.wznow()
}
