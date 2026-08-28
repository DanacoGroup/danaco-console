// Rachunek najbliższego uruchomienia automatyki z notacji cron: pole
// AutomationSchedule.nextRunAt i pozycja podgląd kolejnych uruchomień panelu
// akcji Scheduler. Rachunek idzie w czasie UTC.
package core

import (
	"strings"
	"time"

	"danacoconsole/server/internal/dane"
	"danacoconsole/shared"
)

// horyzontRachunkuCron ogranicza poszukiwanie terminu. Zapis, który przez rok
// nie trafia w żadną chwilę (na przykład 30 lutego), nie ma najbliższego
// uruchomienia i tak też zostaje pokazany.
const horyzontRachunkuCron = 366

// chwilaNastepnegoUruchomienia zwraca znacznik bazy najbliższego uruchomienia
// harmonogramu. Harmonogram wyłączony nie ma terminu — i nie udaje, że ma.
func chwilaNastepnegoUruchomienia(cron *string, wyzwalacze []dane.WyzwalaczAutomatyki,
	czynny bool, teraz time.Time) *string {

	if !czynny {
		return nil
	}
	zapisy := []string{}
	if cron != nil && strings.TrimSpace(*cron) != "" {
		zapisy = append(zapisy, *cron)
	}
	for _, wyzwalacz := range wyzwalacze {
		if wyzwalacz.Czynny && wyzwalacz.Rodzaj == shared.AutomationTriggerKindCron {
			zapisy = append(zapisy, wyzwalacz.Wyrazenie)
		}
	}
	najblizsze := najblizszaChwila(zapisy, teraz)
	if najblizsze == nil {
		return nil
	}
	znacznik := najblizsze.Format(formatZnacznikaBazy)
	return &znacznik
}

// najblizszaChwila wybiera najwcześniejszy termin spośród zapisów cron.
// Harmonogram ma jedną cykliczność podstawową i dowolną liczbę wyzwalaczy
// czasowych — obowiązuje ten, który wypadnie pierwszy.
func najblizszaChwila(zapisy []string, teraz time.Time) *time.Time {
	var najblizsza *time.Time
	for _, zapis := range zapisy {
		chwila, jest := nastepnaChwilaCron(zapis, teraz)
		if !jest {
			continue
		}
		if najblizsza == nil || chwila.Before(*najblizsza) {
			kopia := chwila
			najblizsza = &kopia
		}
	}
	return najblizsza
}

// nastepnaChwilaCron wylicza pierwszą chwilę po teraz pasującą do zapisu cron harmonogramu automatyki.
func nastepnaChwilaCron(zapis string, teraz time.Time) (time.Time, bool) {
	pola, poprawny := polaCronZTekstu(zapis)
	if !poprawny {
		return time.Time{}, false
	}
	poczatek := teraz.UTC().Truncate(time.Minute).Add(time.Minute)
	for dzien := 0; dzien < horyzontRachunkuCron; dzien++ {
		data := poczatek.AddDate(0, 0, dzien)
		if !pola.dzienPasuje(data) {
			continue
		}
		od := time.Date(data.Year(), data.Month(), data.Day(), 0, 0, 0, 0, time.UTC)
		if dzien == 0 {
			od = poczatek
		}
		if chwila, jest := pola.pierwszaChwilaDnia(od); jest {
			return chwila, true
		}
	}
	return time.Time{}, false
}

// polaCron to zbiory dopuszczalnych wartości pięciu pól zapisu cron harmonogramu automatyki cyklicznej.
type polaCron struct {
	minuty         map[int]bool
	godziny        map[int]bool
	dniMiesiaca    map[int]bool
	miesiace       map[int]bool
	dniTygodnia    map[int]bool
	dzienDowolny   bool
	tydzienDowolny bool
}

// dzienPasuje orzeka, czy data mieści się w polach dnia. Reguła klasyczna:
// gdy oba pola dnia są wskazane, wystarczy trafienie jednego z nich.
func (p polaCron) dzienPasuje(data time.Time) bool {
	if !p.miesiace[int(data.Month())] {
		return false
	}
	dzien := p.dniMiesiaca[data.Day()]
	tydzien := p.dniTygodnia[int(data.Weekday())]
	if p.dzienDowolny && p.tydzienDowolny {
		return true
	}
	if p.dzienDowolny {
		return tydzien
	}
	if p.tydzienDowolny {
		return dzien
	}
	return dzien || tydzien
}

// pierwszaChwilaDnia szuka najwcześniejszej godziny i minuty dnia nie wcześniejszej niż wskazana chwila.
func (p polaCron) pierwszaChwilaDnia(od time.Time) (time.Time, bool) {
	for godzina := od.Hour(); godzina < 24; godzina++ {
		if !p.godziny[godzina] {
			continue
		}
		poczatek := 0
		if godzina == od.Hour() {
			poczatek = od.Minute()
		}
		for minuta := poczatek; minuta < 60; minuta++ {
			if p.minuty[minuta] {
				return time.Date(od.Year(), od.Month(), od.Day(), godzina, minuta, 0, 0, time.UTC), true
			}
		}
	}
	return time.Time{}, false
}
