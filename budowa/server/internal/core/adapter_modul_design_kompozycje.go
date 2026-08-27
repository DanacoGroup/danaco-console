// Odpowiedzialność pliku: obszar kompozycji Design Board (`design.board.update`)
// modułu Design — przekład kontraktu na `dane.KompozycjaDesignu` /
// `dane.WarstwaKompozycji` i z powrotem.
package core

import (
	"context"

	"danacoconsole/server/internal/dane"
	"danacoconsole/shared"
)

// Przedrostki identyfikatorów bytów obszaru kompozycji. Zasoby i prompty mają
// własne przedrostki w `adapter_modul_design.go` — tu stoją wyłącznie byty,
// którymi zajmuje się ten plik.
const (
	przedrostekKompozycjiDesign = "plansza-"
	przedrostekWarstwyDesign    = "warstwa-"
)

// ZapiszKompozycje zakłada kompozycję Design Board albo nadpisuje zastaną po
// `BoardId` i podmienia komplet jej warstw. Lista warstw pusta jest stanem
// poprawnym — Operator wyczyścił płótno, a nie zapomniał go wypełnić.
func (a *adapterDesignu) ZapiszKompozycje(ctx context.Context,
	z shared.DesignBoardUpdateRequest) (shared.DesignBoardUpdateResponse, error) {

	if z.WindowId == "" {
		return shared.DesignBoardUpdateResponse{}, bladWskazaniaDesignu(
			"komenda design.board.update bez wskazania okna")
	}

	kod := wartoscTekstu(z.BoardId)
	if kod == "" {
		kod = nowyIdentyfikator(przedrostekKompozycjiDesign)
	}

	kompozycja := dane.KompozycjaDesignu{Kod: kod, Okno: z.WindowId, Nazwa: z.Name}
	warstwy := przelozWarstwyDoZapisu(z.Layers)

	zapisana, err := a.repozytorium.ZapiszKompozycje(ctx, kompozycja, warstwy)
	if err != nil {
		return shared.DesignBoardUpdateResponse{}, bladDesignu(err)
	}

	board, err := a.zlozBoard(ctx, zapisana)
	if err != nil {
		return shared.DesignBoardUpdateResponse{}, bladDesignu(err)
	}
	return shared.DesignBoardUpdateResponse{Board: board}, nil
}

// Kompozycje zwraca kompozycje okna wraz z warstwami — obsługuje
// `design.board.list`. `Total` jest długością wykazu, bo wykaz jest pełny;
// kompozycja bez warstw zostaje w wykazie.
func (a *adapterDesignu) Kompozycje(ctx context.Context,
	z shared.DesignBoardListRequest) (shared.DesignBoardListResponse, error) {

	if z.WindowId == "" {
		return shared.DesignBoardListResponse{}, bladWskazaniaDesignu(
			"komenda design.board.list bez wskazania okna")
	}

	wiersze, err := a.repozytorium.Kompozycje(ctx, z.WindowId)
	if err != nil {
		return shared.DesignBoardListResponse{}, bladDesignu(err)
	}
	kompozycje := make([]shared.DesignBoard, 0, len(wiersze))
	for _, wiersz := range wiersze {
		// Warstwy idą osobnym odczytem na kompozycję, tak samo, jak składa je
		// `ZapiszKompozycje`.
		board, err := a.zlozBoard(ctx, wiersz)
		if err != nil {
			return shared.DesignBoardListResponse{}, bladDesignu(err)
		}
		kompozycje = append(kompozycje, board)
	}
	return shared.DesignBoardListResponse{Boards: kompozycje, Total: len(kompozycje)}, nil
}

// przelozWarstwyDoZapisu przekłada warstwy kontraktu na wiersze warstwy
// danych. Warstwa bez `Id` dostaje identyfikator tu, bo repozytorium wymaga
// kodu przed zapisem.
func przelozWarstwyDoZapisu(warstwyZadania []shared.DesignBoardLayer) []dane.WarstwaKompozycji {
	warstwy := make([]dane.WarstwaKompozycji, 0, len(warstwyZadania))
	for numer, warstwa := range warstwyZadania {
		kod := warstwa.Id
		if kod == "" {
			kod = nowyIdentyfikator(przedrostekWarstwyDesign)
		}
		kolejnosc := numer + 1
		if warstwa.Order != nil {
			kolejnosc = *warstwa.Order
		}
		warstwy = append(warstwy, dane.WarstwaKompozycji{
			Kod:         kod,
			ZasobID:     warstwa.AssetId,
			X:           warstwa.X,
			Y:           warstwa.Y,
			Szerokosc:   warstwa.Width,
			Wysokosc:    warstwa.Height,
			Kolejnosc:   kolejnosc,
			Zablokowana: warstwa.Locked != nil && *warstwa.Locked,
			Adnotacja:   warstwa.Note,
		})
	}
	return warstwy
}

// zlozBoard składa `DesignBoard` kontraktu z kompozycji i jej warstw odczytanych
// osobno — `ZapiszKompozycje` oddaje samą kompozycję, warstwy repozytorium
// niesie oddzielną metodą (`Warstwy`), tak jak dzieli je schemat
// `migracja_048_design.sql`.
func (a *adapterDesignu) zlozBoard(ctx context.Context,
	kompozycja dane.KompozycjaDesignu) (shared.DesignBoard, error) {

	warstwy, err := a.repozytorium.Warstwy(ctx, kompozycja.ID)
	if err != nil {
		return shared.DesignBoard{}, err
	}
	return shared.DesignBoard{
		Id:        kompozycja.Kod,
		WindowId:  kompozycja.Okno,
		Name:      kompozycja.Nazwa,
		Layers:    przelozWarstwyKontraktu(warstwy),
		UpdatedAt: chwilaBazy(kompozycja.Zaktualizowano),
	}, nil
}

// przelozWarstwyKontraktu przekłada wiersze warstwy danych na warstwy kontraktu
// w kolejności, w jakiej repozytorium już je zwróciło (ORDER BY kolejnosc, id).
func przelozWarstwyKontraktu(warstwy []dane.WarstwaKompozycji) []shared.DesignBoardLayer {
	lista := make([]shared.DesignBoardLayer, 0, len(warstwy))
	for _, warstwa := range warstwy {
		kolejnosc := warstwa.Kolejnosc
		zablokowana := warstwa.Zablokowana
		lista = append(lista, shared.DesignBoardLayer{
			Id:      warstwa.Kod,
			AssetId: warstwa.ZasobID,
			X:       warstwa.X,
			Y:       warstwa.Y,
			Width:   warstwa.Szerokosc,
			Height:  warstwa.Wysokosc,
			Order:   &kolejnosc,
			Locked:  &zablokowana,
			Note:    warstwa.Adnotacja,
		})
	}
	return lista
}
