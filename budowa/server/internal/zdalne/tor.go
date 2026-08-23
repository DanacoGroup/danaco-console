// Odpowiedzialność pliku: rozstrzygnięcie toru dla okna. Jedno wejście —
// Przeloz — prowadzi od tożsamości okna do gotowego wiersza poleceń transportu
// albo do odmowy trójczęściowej nazywającej brakującą część drogi.
package zdalne

import (
	"fmt"
	"os/exec"
	"strings"
)

// Przeloz przekłada polecenie procesu okna na wywołanie SSH do hosta wykonania
// tego okna. Droga ma pięć ogniw i każde brakujące jest osobną, nazwaną odmową:
//
//  1. baza rdzenia zasilona (Zasil — wpięcie kompozycji);
//  2. host wskazany ustawieniem `host_wykonania` (okno albo poziom globalny);
//  3. host wpisany do wykazu `host_zdalny`;
//  4. zgoda Operatora na tym wierszu wydana;
//  5. program `ssh` obecny na maszynie rdzenia.
//
// Odmowa nie jest bramką wobec Operatora: każda mówi, co się nie stało,
// dlaczego, i którym ruchem Operator to zmienia. Zgoda per host
// chroni maszyny Operatora — rdzeń nie zainicjuje połączenia z maszyną,
// której mu nie oddano.
func Przeloz(idOkna string, p Polecenie) (Uruchomienie, error) {
	if strings.TrimSpace(p.Program) == "" {
		return Uruchomienie{}, fmt.Errorf("zdalne: polecenie okna %s nie niesie programu, "+
			"więc nie ma czego uruchomić na hoście zdalnym", idOkna)
	}

	nazwa, err := hostOkna(idOkna)
	if err != nil {
		return Uruchomienie{}, err
	}
	if nazwa == "" {
		return Uruchomienie{}, fmt.Errorf("zdalne: proces okna %s nie ruszył na hoście "+
			"zdalnym, bo żaden host nie jest wskazany — ustawienie host_wykonania jest "+
			"puste na poziomie okna i globalnym; host wskazuje się w oknie komunikacji "+
			"polem „Host wykonania”", idOkna)
	}

	host, err := hostZRejestru(nazwa)
	if err != nil {
		return Uruchomienie{}, err
	}

	sciezkaSSH, err := exec.LookPath("ssh")
	if err != nil {
		return Uruchomienie{}, fmt.Errorf("zdalne: proces okna %s nie ruszył na hoście %q, "+
			"bo maszyna rdzenia nie ma programu ssh w PATH — tor jedzie wyłącznie po SSH "+
			"i wymaga jego instalacji: %w", idOkna, nazwa, err)
	}

	return zbudujUruchomienie(sciezkaSSH, host, p), nil
}
