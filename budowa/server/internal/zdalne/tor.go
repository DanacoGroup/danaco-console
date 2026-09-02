// Odpowiedzialność pliku: rozstrzygnięcie toru dla okna. Jedno wejście —
// Przeloz — prowadzi od tożsamości okna do gotowego wiersza poleceń transportu
// albo do odmowy trójczęściowej nazywającej brakującą część drogi.
package zdalne

import (
	"context"
	"fmt"
	"os/exec"
	"strings"
)

// Funkcja Przeloz przekłada polecenie procesu okna na gotowe wywołanie SSH do hosta wykonania tego okna.
func Przeloz(ctx context.Context, idOkna string, p Polecenie) (Uruchomienie, error) {
	if strings.TrimSpace(p.Program) == "" {
		return Uruchomienie{}, fmt.Errorf("zdalne: polecenie okna %s nie niesie programu, "+
			"więc nie ma czego uruchomić na hoście zdalnym", idOkna)
	}

	nazwa, err := hostOkna(ctx, idOkna)
	if err != nil {
		return Uruchomienie{}, err
	}
	if nazwa == "" {
		return Uruchomienie{}, fmt.Errorf("zdalne: proces okna %s nie ruszył na hoście "+
			"zdalnym, bo żaden host nie jest wskazany — ustawienie host_wykonania jest "+
			"puste na poziomie okna i globalnym; host wskazuje się w oknie komunikacji "+
			"polem „Host wykonania”", idOkna)
	}

	host, err := hostZRejestru(ctx, nazwa)
	if err != nil {
		return Uruchomienie{}, err
	}

	sciezkaSSH, err := exec.LookPath("ssh")
	if err != nil {
		return Uruchomienie{}, fmt.Errorf("zdalne: proces okna %s nie ruszył na hoście %q, "+
			"bo maszyna serwera nie ma programu ssh w PATH — tor jedzie wyłącznie po SSH "+
			"i wymaga jego instalacji: %w", idOkna, nazwa, err)
	}

	return zbudujUruchomienie(sciezkaSSH, host, p)
}
