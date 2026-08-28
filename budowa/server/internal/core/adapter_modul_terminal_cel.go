// Plik obsługuje adres celu karty powłoki zdalnej: pola `remoteTarget`, `remotePort` i `hostId` żądania `terminal.session.open`. Wpis książki hostów jest wskazaniem, nie kopią: karta bierze zeń adres, port, katalog i ścieżkę klucza.
package core

import (
	"context"
	"errors"
	"strings"

	"danacoconsole/server/internal/dane"
	"danacoconsole/shared"
)

// wskazCelZdalny wypełnia adres celu karty z żądania albo z wpisu książki hostów. Wskazania sprzeczne kończą się odmową, a nie cichym pierwszeństwem jednego z nich.
func (a *adapterTerminala) wskazCelZdalny(ctx context.Context, karta *kartaTerminala,
	z shared.TerminalSessionOpenRequest) error {

	cel := strings.TrimSpace(wartoscTekstu(z.RemoteTarget))
	kodHosta := strings.TrimSpace(wartoscTekstu(z.HostId))

	if kodHosta != "" && cel != "" {
		return bladZadaniaTerminala(
			"karta wskazuje naraz wpis książki hostów i własny adres celu — " +
				"zostaw jedno z dwojga, bo wpis niesie już adres")
	}
	if kodHosta != "" {
		return a.celZWpisuKsiazki(ctx, karta, kodHosta)
	}
	karta.celZdalny = cel
	if z.RemotePort != nil && *z.RemotePort > 0 {
		karta.portZdalny = *z.RemotePort
	}
	// Droga zastępcza: karta bez pola `remoteTarget`, ale ze zmienną `SSH_TARGET`, działa jak dotąd.
	if karta.celZdalny == "" {
		karta.celZdalny = strings.TrimSpace(karta.srodowisko[zmiennaCeluSSH])
	}
	return nil
}

// celZWpisuKsiazki bierze adres, port, katalog i ścieżkę klucza z wpisu książki hostów wskazanego identyfikatorem.
func (a *adapterTerminala) celZWpisuKsiazki(ctx context.Context, karta *kartaTerminala,
	kodHosta string) error {

	if a.repozytorium == nil {
		return bladZadaniaTerminala("rdzeń pracuje bez dziennika, więc nie ma książki hostów; " +
			"podaj adres celu wprost polem remoteTarget")
	}
	wpis, err := a.repozytorium.Host(ctx, kodHosta)
	if errors.Is(err, dane.ErrBrakWiersza) {
		return bladBrakuZasobuTerminala("wpis książki hostów " + kodHosta)
	}
	if err != nil {
		return err
	}
	karta.hostKod = wpis.Kod
	karta.celZdalny = wpis.Cel
	if wpis.Port != nil && *wpis.Port > 0 {
		karta.portZdalny = int(*wpis.Port)
	}
	if karta.katalog == "" {
		karta.katalog = strings.TrimSpace(wpis.KatalogRoboczy)
	}
	if wpis.KluczKod == nil || strings.TrimSpace(*wpis.KluczKod) == "" {
		return nil
	}
	klucz, err := a.repozytorium.Klucz(ctx, *wpis.KluczKod)
	if errors.Is(err, dane.ErrBrakWiersza) {
		// Wpis wskazuje klucz zdjęty z wykazu; karta rusza wtedy z kluczem domyślnym konfiguracji maszyny.
		return nil
	}
	if err != nil {
		return err
	}
	karta.kluczSciezka = klucz.Sciezka
	return nil
}
