// Odpowiedzialność pliku: adres celu karty powłoki zdalnej — pola
// `remoteTarget`, `remotePort` i `hostId` żądania `terminal.session.open`.
//
// ── Dlaczego adres przestał jechać zmienną środowiska ────────────────────────
// Karta zdalna brała dotąd adres ze zmiennej `SSH_TARGET`, bo kontrakt nie miał
// na niego pola. Miało to dwa skutki, których żaden nie był zamierzony: portu
// nie dało się podać osobno (zmienna niesie jeden napis, a `ssh` chce `-p`),
// a zmienne środowiska karty z zamysłu NIE MAJĄ kolumny w bazie — więc karta
// zdalna odtworzona po restarcie rdzenia traciła adres i pierwsze polecenie
// kończyło się odmową. Kontrakt ma dziś pola wprost, a karta — kolumny
// (migracja 251).
//
// Zmienna `SSH_TARGET` zostaje jako droga zastępcza, nie jako droga główna:
// karta założona przed tą zmianą i klient, który jeszcze nie przestawił się na
// nowe pola, mają dalej działać.
//
// ── Wpis książki hostów jest wskazaniem, nie kopią ──────────────────────────
// Pole `hostId` bierze z wpisu adres, port, katalog roboczy i ŚCIEŻKĘ klucza.
// Materiału klucza nie tyka nikt: `ssh` dostaje ścieżkę przełącznikiem `-i`,
// a plik czyta sam, na maszynie rdzenia.
package core

import (
	"context"
	"errors"
	"strings"

	"danacoconsole/server/internal/dane"
	"danacoconsole/shared"
)

// wskazCelZdalny wypełnia adres celu karty z żądania albo z wpisu książki hostów.
//
// Wskazania sprzeczne kończą się odmową, a nie cichym pierwszeństwem jednego
// z nich: Operator, który podał i wpis książki, i własny adres, ma dwa różne
// zamiary i nie ma powodu zgadywać, który z nich jest ważniejszy.
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
	// Droga zastępcza: karta bez pola `remoteTarget`, ale ze zmienną `SSH_TARGET`
	// w środowisku, zachowuje się jak dotąd.
	if karta.celZdalny == "" {
		karta.celZdalny = strings.TrimSpace(karta.srodowisko[zmiennaCeluSSH])
	}
	return nil
}

// celZWpisuKsiazki bierze adres, port, katalog i ścieżkę klucza z książki hostów.
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
		// Wpis wskazuje klucz zdjęty z wykazu. Karta rusza z kluczem domyślnym
		// konfiguracji maszyny — tak samo jak wpis bez wskazania klucza —
		// bo odmowa byłaby tu surowsza od skutku.
		return nil
	}
	if err != nil {
		return err
	}
	karta.kluczSciezka = klucz.Sciezka
	return nil
}
