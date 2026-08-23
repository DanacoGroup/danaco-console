package models

import "context"

// Ujscie odbiera fragmenty strumienia w kolejności ich powstania. Adapter kanału
// nie wie, kto jest po drugiej stronie: warstwa sesji, kolejka, test albo
// koordynator czytający strumień wykonawcy.
type Ujscie interface {
	Fragment(ctx context.Context, f Fragment) error
}

// UjscieFunkcji czyni ujściem zwykłą funkcję — wywołujący nie musi definiować
// typu tylko po to, by odebrać strumień.
type UjscieFunkcji func(ctx context.Context, f Fragment) error

// Fragment przekazuje fragment funkcji opakowanej w ujście.
func (u UjscieFunkcji) Fragment(ctx context.Context, f Fragment) error {
	if u == nil {
		return nil
	}
	return u(ctx, f)
}

// NadajProwenancje wysyła fragment prowenancji jako pierwszy fragment
// strumienia. Niepowodzenie kodowania prowenancji nie przerywa wywołania: opis
// wywołania jest przejrzystością, nie bramką, więc strumień idzie dalej bez
// niego.
func NadajProwenancje(ctx context.Context, u Ujscie, z Zapytanie, p Prowenancja) error {
	if u == nil {
		return nil
	}
	f, err := FragmentProwenancji(z, p)
	if err != nil {
		return nil
	}
	return u.Fragment(ctx, f)
}

// NadajKonto wysyła fragment metadanych konta. Tak samo jak prowenancja jest
// opisem wywołania i nie przerywa strumienia własnym niepowodzeniem.
func NadajKonto(ctx context.Context, u Ujscie, z Zapytanie, k MetadaneKonta) error {
	if u == nil {
		return nil
	}
	f, err := FragmentKonta(z, k)
	if err != nil {
		return nil
	}
	return u.Fragment(ctx, f)
}

// NadajTekst wysyła porcję tekstu odpowiedzi modelu.
func NadajTekst(ctx context.Context, u Ujscie, z Zapytanie, tekst string) error {
	if u == nil {
		return nil
	}
	return u.Fragment(ctx, FragmentTekstu(z, tekst))
}

// NadajObraz wysyła treść wizualną odpowiedzi. W odróżnieniu od prowenancji
// i konta niepowodzenie kodowania nie jest tu drobiazgiem: obraz jest całym
// wynikiem wywołania, więc błąd wraca do wołającego, zamiast zostawić strumień
// pustym i milczącym.
func NadajObraz(ctx context.Context, u Ujscie, z Zapytanie, t TrescObrazu) error {
	if u == nil {
		return nil
	}
	f, err := FragmentObrazu(z, t)
	if err != nil {
		return err
	}
	return u.Fragment(ctx, f)
}

// NadajBlad wysyła fragment błędu technicznego kanału i zwraca ten sam błąd
// jako wynik wywołania — warstwa wyżej dostaje go obiema drogami: w strumieniu
// dla użytkownika i w wyniku dla obsługi.
func NadajBlad(ctx context.Context, u Ujscie, z Zapytanie, err error) error {
	if err == nil {
		return nil
	}
	blad := BladKanalu(err)
	if u != nil {
		_ = u.Fragment(ctx, FragmentBledu(z, blad))
	}
	return err
}
