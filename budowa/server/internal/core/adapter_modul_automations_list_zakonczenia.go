// Odpowiedzialność pliku: list o zakończonym przebiegu automatyki — wyzwalacz
// przy przejściu w stan końcowy i odwzorowanie zmiennych listu 7.
package core

import (
	"context"
	"fmt"
	netmail "net/mail"
	"strconv"
	"strings"
	"time"

	"danacoconsole/server/internal/dane"
	"danacoconsole/server/internal/konfiguracja"
	"danacoconsole/server/internal/mail"
	"danacoconsole/server/internal/nadajnik"
	"danacoconsole/shared"
)

/*
stanyListu przekłada stan przebiegu na słowo, którego używa list.

Rdzeń zna trzy stany końcowe, dostawa listu dopuszcza dwa. Powodzenie jest
„zakończony”; niepowodzenie i przerwanie idą jako „zatrzymany”, bo oba znaczą
dla Operatora to samo: bieg stanął przed metą i czeka na jego rozstrzygnięcie.
*/
var stanyListu = map[string]string{
	shared.AutomationExecutionStatusSucceeded: "zakończony",
	shared.AutomationExecutionStatusFailed:    "zatrzymany",
	shared.AutomationExecutionStatusStopped:   "zatrzymany",
}

// listZakonczeniaPrzebiegu wysyła list 7, gdy przebieg właśnie wszedł w stan
// końcowy. Znacznik `Powiadomiono` zamyka drogę po pierwszym nadaniu.
func (a *adapterAutomatyk) listZakonczeniaPrzebiegu(ctx context.Context,
	przebieg dane.Przebieg) (bool, error) {

	if !czyStanKoncowyPrzebiegu(przebieg.Stan) || przebieg.Powiadomiono {
		return false, nil
	}
	if a.nastawyListow == nil {
		return false, nil
	}
	adresat, nazwa, err := a.adresatPrzebiegu(ctx, przebieg)
	if err != nil || adresat == "" {
		return false, err
	}
	nadawca := a.nastawyListow(ctx)
	if err := nadawca.Brak(); err != nil {
		/* Brak konta nadawczego nie jest błędem przebiegu: bieg się zakończył,
		   a listu nie ma czym wysłać. Znacznik zostaje zdjęty, żeby list poszedł
		   po uzupełnieniu nastaw. */
		return false, nil
	}
	komplet, err := kompletListow()
	if err != nil {
		return false, err
	}
	teraz := time.Now()
	wiadomosc, err := komplet.Build(mail.KindRunFinished, mail.Envelope{
		From: nadawca.Nadawca(),
		To:   netmail.Address{Address: adresat},
		Date: teraz,
	}, mail.Content{Values: wartosciListuPrzebiegu(przebieg, nazwa, adresat, a.adresKonsoli, teraz)},
		znakiMarki())
	if err != nil {
		return false, fmt.Errorf("nie udało się złożyć listu o przebiegu %s: %w", przebieg.Kod, err)
	}
	if _, err := nadajnik.Wyslij(nadawca, wiadomosc); err != nil {
		return false, fmt.Errorf("nie udało się wysłać listu o przebiegu %s: %w", przebieg.Kod, err)
	}
	return true, nil
}

// adresatPrzebiegu wyprowadza odbiorcę listu łańcuchem przebieg → automatyka →
// konto. Pusty adres znaczy automatykę bez wskazania konta i list nie wychodzi.
func (a *adapterAutomatyk) adresatPrzebiegu(ctx context.Context,
	przebieg dane.Przebieg) (string, string, error) {

	wlasciciel, err := a.repozytorium.WlascicielAutomatyki(ctx, przebieg.AutomatykaID)
	if err != nil {
		return "", "", err
	}
	if wlasciciel.KontoId == 0 || a.konta == nil {
		return "", wlasciciel.Nazwa, nil
	}
	konto, err := a.konta.KontoPoId(ctx, wlasciciel.KontoId)
	if err != nil {
		return "", wlasciciel.Nazwa, nil
	}
	return konto.Email, wlasciciel.Nazwa, nil
}

// wartosciListuPrzebiegu składa komplet zmiennych listu 7.
func wartosciListuPrzebiegu(przebieg dane.Przebieg, nazwa, adresat, adresKonsoli string,
	teraz time.Time) map[string]string {

	return map[string]string{
		"recipient_address": adresat,
		"support_address":   adresWsparcia,
		"year":              strconv.Itoa(teraz.Year()),
		"run_name":          nazwaAlboKod(nazwa, przebieg.AutomatykaKod),
		"run_status":        stanyListu[przebieg.Stan],
		"run_step":          fmt.Sprintf("%d z %d", przebieg.EtapBiezacy, przebieg.Etapow),
		"run_duration":      trwaniePrzebiegu(przebieg),
		"started_at":        czasListu(przebieg.Rozpoczeto, teraz),
		"run_url":           konfiguracja.AdresPrzebiegu(adresKonsoli, przebieg.Kod),
		/* Środowiska przebieg nie zna: migracja 039 stanowi wprost, że automatyka
		   nie jest bytem sesji ani okna. Nazwa produktu jest jedynym wskazaniem,
		   które nie byłoby zmyślone. */
		"environment_name": "Danaco Console",
	}
}

// nazwaAlboKod oddaje nazwę automatyki, a przy jej braku identyfikator — list
// musi nazwać bieg, którego dotyczy.
func nazwaAlboKod(nazwa, kod string) string {
	if strings.TrimSpace(nazwa) == "" {
		return kod
	}
	return nazwa
}

// trwaniePrzebiegu liczy czas między początkiem a końcem biegu. Postać jest
// zwięzła — Operator czyta ją z tabeli listu, nie z dziennika.
func trwaniePrzebiegu(przebieg dane.Przebieg) string {
	poczatek, err := time.Parse(formatZnacznikaBazy, przebieg.Rozpoczeto)
	if err != nil || przebieg.Zakonczono == nil {
		return "—"
	}
	koniec, err := time.Parse(formatZnacznikaBazy, *przebieg.Zakonczono)
	if err != nil {
		return "—"
	}
	trwanie := koniec.Sub(poczatek).Round(time.Second)
	if trwanie < time.Minute {
		return fmt.Sprintf("%d s", int(trwanie.Seconds()))
	}
	if trwanie < time.Hour {
		return fmt.Sprintf("%d min %d s", int(trwanie.Minutes()), int(trwanie.Seconds())%60)
	}
	return fmt.Sprintf("%d h %d min", int(trwanie.Hours()), int(trwanie.Minutes())%60)
}

/*
czasListu podaje czas rozpoczęcia w postaci czytelnej dla Operatora.

Rdzeń zapisuje czas w UTC i strefy Operatora nie zna, więc list niesie skrót
strefy zapisu, nie strefy odbiorcy. Znacznik nieczytelny zastępuje chwila
złożenia listu — wiersz pusty czyta się jak błąd składania.
*/
func czasListu(znacznik string, zapas time.Time) string {
	chwila, err := time.Parse(formatZnacznikaBazy, znacznik)
	if err != nil {
		return zapas.Format(postacCzasuListu)
	}
	return chwila.Format(postacCzasuListu)
}
