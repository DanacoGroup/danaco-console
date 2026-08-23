// Odpowiedzialność pliku: obiekty osadzone w dokumencie — obrazy, logo,
// kształty, ikony i pola tekstowe. Cztery czynności kontraktu: wstawienie,
// wykaz, postać obiektu i usunięcie.
//
// ── Pochodzenie obiektu jest obowiązkowe ────────────────────────────────────
// Zlecenie mówi wprost: wstawienie obrazu bez zapisanego pochodzenia jest
// brakiem, nie skrótem. Dlatego każdy obiekt niosący bajty ma zapisane, SKĄD
// jest — zasób magazynu rdzenia, węzeł modułu Design, plik Biblioteki, baza
// zdjęciowa albo adres w sieci. Obraz bez pochodzenia jest za tydzień obrazem,
// o którym nikt nie wie, czy wolno go było użyć.
//
// ── Bajty biorą się z magazynu zasobów rdzenia ──────────────────────────────
// Obiekt nie nosi bajtów w swoim wierszu: nosi wskazanie zasobu, a bajty leżą
// w magazynie pod sumą kontrolną. Ta droga jest jedna dla całego modułu
// (`odlozTrescStudia`, `bajtyZasobuStudia`) i tu się jej nie zakłada drugi raz.
//
// ── Kształty i ikony są rachunkiem modułu Design ─────────────────────────────
// Ikona wchodzi z KATALOGU Designu (`ikonaKataloguDesignu`,
// `svgIkonyKataloguDesignu`) — tego samego, którym jedzie
// `design.icon.library.search`. Drugiego katalogu ikon w Studiu nie ma i mieć
// nie będzie. Kształt rysowany na miejscu opisuje się rodzajem, rozmiarem
// i wyglądem, a nie własnym rachunkiem ścieżek; kształt wymagający ścieżek
// edytowalnych wskazuje się węzłem Designu (`designNodeId`), który powstaje
// przez `design.vector.shape.add`. Do plików modułu Design ten odcinek nie
// wchodzi.
package core

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"os"
	"strconv"
	"strings"

	"danacoconsole/server/internal/dane"
	"danacoconsole/shared"
)

// obiektRozmiarDomyslnyMm to rozmiar obiektu, którego Operator nie podał —
// szerokość połowy kolumny tekstu pisma A4.
const (
	obiektSzerokoscDomyslnaMm = 80.0
	obiektWysokoscDomyslnaMm  = 60.0
	// obiektGranicaBajtow chroni bazę i magazyn przed obrazem, którego nikt nie
	// zamierzał wstawiać do pisma.
	obiektGranicaBajtow = 64 << 20
)

// WstawObiekt wstawia obiekt w miejsce kursora (`studio.object.insert`).
func (a *adapterStudia) WstawObiekt(ctx context.Context,
	z shared.StudioObjectInsertRequest) (shared.StudioObjectInsertResponse, error) {

	stan, err := a.postacWczytaj(ctx, z.DocumentId)
	if err != nil {
		return shared.StudioObjectInsertResponse{}, err
	}
	autor := postacAutor(z.Author)
	if err := obiektSprawdzRodzaj(z.Kind); err != nil {
		return shared.StudioObjectInsertResponse{}, err
	}

	dlugosc := postacDlugosc(&stan.forma)
	miejsce := z.Offset
	if miejsce < 0 {
		miejsce = 0
	}
	if miejsce > dlugosc {
		miejsce = dlugosc
	}
	if odcinki, pominiete := postacOdcinkiDozwolone(&stan.forma, miejsce, miejsce, autor); len(odcinki) == 0 {
		return shared.StudioObjectInsertResponse{}, bladWskazaniaStudio(
			"wstawienie obiektu " + postacZapisZakresu(miejsce, miejsce) +
				" zatrzymane przez blokadę fragmentu: " + postacNazwaBlokad(pominiete))
	}

	obiekt := shared.StudioDocumentObject{
		Id: nowyIdentyfikator(przedrostekObiektuPostaci), Kind: z.Kind,
		AltText: z.AltText, InnerText: z.InnerText, Caption: z.Caption,
		WidthMm: z.WidthMm, HeightMm: z.HeightMm, Wrap: z.Wrap, Anchor: z.Anchor,
		ShapeKind: z.ShapeKind, FillColor: z.FillColor, StrokeColor: z.StrokeColor,
		AnchorOffset: postacWskaznikLiczby(miejsce),
	}
	bilans := shared.StudioActionBalance{Skipped: []shared.StudioSkippedItem{}}
	if err := a.obiektUstalPochodzenie(ctx, z, &obiekt, &bilans); err != nil {
		return shared.StudioObjectInsertResponse{}, err
	}
	obiektDomyslneWymiary(&obiekt)
	if obiekt.Anchor == nil {
		zakotwiczenie := shared.StudioAnchorKind(shared.StudioAnchorKindParagraph)
		obiekt.Anchor = &zakotwiczenie
	}
	if obiekt.Wrap == nil {
		oplywanie := shared.StudioTextWrap(shared.StudioTextWrapInline)
		obiekt.Wrap = &oplywanie
	}
	obiekt.ZOrder = postacWskaznikLiczby(obiektNastepnaWarstwa(&stan.forma))

	wiersz, err := obiektDoWiersza(stan.dokument.ID, obiekt)
	if err != nil {
		return shared.StudioObjectInsertResponse{}, err
	}
	skladnica, err := a.postacSkladnica()
	if err != nil {
		return shared.StudioObjectInsertResponse{}, err
	}
	zapisany, err := skladnica.ZapiszObiektDokumentu(ctx, wiersz)
	if err != nil {
		return shared.StudioObjectInsertResponse{}, bladStudio(err)
	}

	// Obiekt wchodzi do drzewa blokiem nietekstowym: nie zajmuje ani jednego
	// znaku treści, więc nie przesuwa zaznaczeń, przypisów ani blokad.
	tabelaWstawBlokWMiejscu(&stan.forma, shared.StudioDocumentBlock{
		Id: nowyIdentyfikator(przedrostekBlokuPostaci), Kind: blokPostaciObiekt,
		ObjectId: postacWskaznikTekstu(obiekt.Id),
	}, miejsce)

	// Obraz wstawiony unieważnia spis ilustracji — spis, który go nie zna,
	// pokazuje stan sprzed wstawienia.
	if err := a.aparatZnaczNieswiezoscRodzaju(ctx, stan,
		shared.StudioApparatusKindFigureIndex); err != nil {

		return shared.StudioObjectInsertResponse{}, err
	}
	if err := a.postacWczytajWiersze(ctx, stan); err != nil {
		return shared.StudioObjectInsertResponse{}, err
	}

	bilans.Applied = 1
	bilans.Note = postacWskaznikTekstu(obiektNazwaRodzaju(z.Kind) + " wstawiony " +
		postacZapisZakresu(miejsce, miejsce) + "; pochodzenie: " +
		obiektZapisPochodzenia(obiekt))
	stan.opisCzynnosci = "wstawienie obiektu: " + obiektNazwaRodzaju(z.Kind)

	forma, bilansGotowy, zmiana, err := a.postacZakoncz(ctx, stan, autor,
		shared.StudioChangeKindWstawienie, shared.StudioActionKindObjectChange,
		miejsce, miejsce, bilans)
	if err != nil {
		return shared.StudioObjectInsertResponse{}, err
	}
	zlozony := postacZlozObiekty([]dane.ObiektDokumentuStudia{zapisany})
	if len(zlozony) == 0 {
		return shared.StudioObjectInsertResponse{}, postacBladZaplecza(
			"obiekt zapisany, ale nie da się go złożyć do odpowiedzi")
	}
	return shared.StudioObjectInsertResponse{
		Form: forma, Balance: bilansGotowy, Change: zmiana,
		ActionId: stan.czynnosc, Object: zlozony[0],
	}, nil
}

// WykazObiektow oddaje obiekty osadzone w dokumencie (`studio.object.list`).
func (a *adapterStudia) WykazObiektow(ctx context.Context,
	z shared.StudioObjectListRequest) (shared.StudioObjectListResponse, error) {

	stan, err := a.postacWczytaj(ctx, z.DocumentId)
	if err != nil {
		return shared.StudioObjectListResponse{}, err
	}
	if z.Kind != nil {
		if err := obiektSprawdzRodzaj(*z.Kind); err != nil {
			return shared.StudioObjectListResponse{}, err
		}
	}
	obiekty := make([]shared.StudioDocumentObject, 0, len(stan.forma.Objects))
	for _, obiekt := range stan.forma.Objects {
		if z.Kind != nil && obiekt.Kind != *z.Kind {
			continue
		}
		// Obiekt stojący pod blokadą fragmentu mówi to o sobie wprost: okno ma
		// pokazać, czego model nie tknie, PRZED próbą, a nie po odmowie.
		miejsce := aparatWartoscLiczby(obiekt.AnchorOffset)
		if len(postacBlokadyZakresu(&stan.forma, miejsce, miejsce)) > 0 {
			obiekt.Locked = postacWskaznikPrawdy(true)
		}
		obiekty = append(obiekty, obiekt)
	}
	return shared.StudioObjectListResponse{Objects: obiekty}, nil
}

// UstawPostacObiektu ustawia postać obiektu (`studio.object.format.set`) —
// rozmiar, przycięcie, opływanie tekstem, położenie i zakotwiczenie, warstwę,
// tekst zastępczy, obramowanie, wypełnienie i obrót.
func (a *adapterStudia) UstawPostacObiektu(ctx context.Context,
	z shared.StudioObjectFormatSetRequest) (shared.StudioObjectFormatSetResponse, error) {

	stan, err := a.postacWczytaj(ctx, z.DocumentId)
	if err != nil {
		return shared.StudioObjectFormatSetResponse{}, err
	}
	autor := postacAutor(z.Author)
	obiekt, jest := obiektZnajdz(&stan.forma, z.ObjectId)
	if !jest {
		return shared.StudioObjectFormatSetResponse{}, bladWskazaniaStudio(
			"zmiana postaci obiektu „" + strings.TrimSpace(z.ObjectId) +
				"”, którego dokument nie ma")
	}
	miejsce := aparatWartoscLiczby(obiekt.AnchorOffset)
	if odcinki, pominiete := postacOdcinkiDozwolone(&stan.forma, miejsce, miejsce, autor); len(odcinki) == 0 {
		return shared.StudioObjectFormatSetResponse{}, bladWskazaniaStudio(
			"zmiana postaci obiektu zatrzymana przez blokadę fragmentu: " +
				postacNazwaBlokad(pominiete))
	}

	bilans := shared.StudioActionBalance{Skipped: []shared.StudioSkippedItem{}}
	if len(z.Crop) > 0 {
		var przyciecie shared.StudioObjectCrop
		if err := json.Unmarshal(z.Crop, &przyciecie); err != nil {
			return shared.StudioObjectFormatSetResponse{}, bladWskazaniaStudio(
				"przycięcie w żądaniu jest nieczytelne: " + err.Error())
		}
		if obiekt.Kind != shared.StudioObjectKindImage && obiekt.Kind != shared.StudioObjectKindLogo {
			bilans.Skipped = append(bilans.Skipped, shared.StudioSkippedItem{
				Reason: "przycięcie dotyczy wyłącznie obrazu",
				Detail: postacWskaznikTekstu(obiektNazwaRodzaju(obiekt.Kind) +
					" nie ma czego przycinać — przycięcie obowiązuje dla obrazu i logo"),
			})
		} else {
			obiekt.Crop = &przyciecie
			bilans.Applied++
		}
	}
	if len(z.Border) > 0 {
		var obramowanie shared.StudioBorder
		if err := json.Unmarshal(z.Border, &obramowanie); err != nil {
			return shared.StudioObjectFormatSetResponse{}, bladWskazaniaStudio(
				"obramowanie w żądaniu jest nieczytelne: " + err.Error())
		}
		obiekt.Border = &obramowanie
		bilans.Applied++
	}

	// Zmiana rozmiaru z zachowaniem proporcji przelicza wymiar drugi z wymiarów
	// zastanych: Operator, który podał samą szerokość, nie spodziewa się obrazu
	// rozciągniętego w pionie.
	staraSzerokosc, staraWysokosc := obiektWymiary(obiekt)
	if z.WidthMm != nil {
		obiekt.WidthMm = z.WidthMm
		bilans.Applied++
	}
	if z.HeightMm != nil {
		obiekt.HeightMm = z.HeightMm
		bilans.Applied++
	}
	if z.KeepAspect != nil && *z.KeepAspect && staraSzerokosc > 0 && staraWysokosc > 0 {
		switch {
		case z.WidthMm != nil && z.HeightMm == nil:
			obiekt.HeightMm = postacWskaznikMiary(*z.WidthMm * staraWysokosc / staraSzerokosc)
		case z.HeightMm != nil && z.WidthMm == nil:
			obiekt.WidthMm = postacWskaznikMiary(*z.HeightMm * staraSzerokosc / staraWysokosc)
		}
	}
	for _, cecha := range []struct {
		podane bool
		zapis  func()
	}{
		{z.Wrap != nil, func() { obiekt.Wrap = z.Wrap }},
		{z.Anchor != nil, func() { obiekt.Anchor = z.Anchor }},
		{z.PositionXMm != nil, func() { obiekt.PositionXMm = z.PositionXMm }},
		{z.PositionYMm != nil, func() { obiekt.PositionYMm = z.PositionYMm }},
		{z.ZOrder != nil, func() { obiekt.ZOrder = z.ZOrder }},
		{z.AltText != nil, func() { obiekt.AltText = z.AltText }},
		{z.FillColor != nil, func() { obiekt.FillColor = z.FillColor }},
		{z.StrokeColor != nil, func() { obiekt.StrokeColor = z.StrokeColor }},
		{z.StrokeWidthPt != nil, func() { obiekt.StrokeWidthPt = z.StrokeWidthPt }},
		{z.Shadow != nil, func() { obiekt.Shadow = z.Shadow }},
		{z.RotationDeg != nil, func() { obiekt.RotationDeg = z.RotationDeg }},
		{z.InnerText != nil, func() { obiekt.InnerText = z.InnerText }},
	} {
		if cecha.podane {
			cecha.zapis()
			bilans.Applied++
		}
	}
	if z.AnchorOffset != nil {
		nowe := *z.AnchorOffset
		if nowe < 0 {
			nowe = 0
		}
		if dlugosc := postacDlugosc(&stan.forma); nowe > dlugosc {
			nowe = dlugosc
		}
		obiekt.AnchorOffset = postacWskaznikLiczby(nowe)
		bilans.Applied++
	}
	podpisZmieniony := false
	if z.Caption != nil {
		obiekt.Caption = z.Caption
		podpisZmieniony = true
		bilans.Applied++
	}
	if bilans.Applied == 0 {
		return shared.StudioObjectFormatSetResponse{}, bladWskazaniaStudio(
			"ustawienie postaci obiektu bez ani jednej cechy do ustawienia")
	}

	wiersz, err := obiektDoWiersza(stan.dokument.ID, obiekt)
	if err != nil {
		return shared.StudioObjectFormatSetResponse{}, err
	}
	skladnica, err := a.postacSkladnica()
	if err != nil {
		return shared.StudioObjectFormatSetResponse{}, err
	}
	zapisany, err := skladnica.ZapiszObiektDokumentu(ctx, wiersz)
	if err != nil {
		return shared.StudioObjectFormatSetResponse{}, bladStudio(err)
	}
	if podpisZmieniony {
		if err := a.aparatZnaczNieswiezoscRodzaju(ctx, stan,
			shared.StudioApparatusKindFigureIndex); err != nil {

			return shared.StudioObjectFormatSetResponse{}, err
		}
	}
	if err := a.postacWczytajWiersze(ctx, stan); err != nil {
		return shared.StudioObjectFormatSetResponse{}, err
	}

	bilans.Note = postacWskaznikTekstu("postać obiektu ustawiona; cech zmienionych " +
		strconv.Itoa(bilans.Applied))
	stan.opisCzynnosci = "zmiana postaci obiektu"

	forma, bilansGotowy, zmiana, err := a.postacZakoncz(ctx, stan, autor,
		shared.StudioChangeKindFormatowanie, shared.StudioActionKindObjectChange,
		miejsce, miejsce, bilans)
	if err != nil {
		return shared.StudioObjectFormatSetResponse{}, err
	}
	zlozony := postacZlozObiekty([]dane.ObiektDokumentuStudia{zapisany})
	if len(zlozony) == 0 {
		return shared.StudioObjectFormatSetResponse{}, postacBladZaplecza(
			"obiekt zapisany, ale nie da się go złożyć do odpowiedzi")
	}
	return shared.StudioObjectFormatSetResponse{
		Form: forma, Balance: bilansGotowy, Change: zmiana,
		ActionId: stan.czynnosc, Object: zlozony[0],
	}, nil
}

// UsunObiekt usuwa obiekt osadzony w dokumencie (`studio.object.remove`).
func (a *adapterStudia) UsunObiekt(ctx context.Context,
	z shared.StudioObjectRemoveRequest) (shared.StudioObjectRemoveResponse, error) {

	stan, err := a.postacWczytaj(ctx, z.DocumentId)
	if err != nil {
		return shared.StudioObjectRemoveResponse{}, err
	}
	autor := postacAutor(z.Author)
	kod := strings.TrimSpace(z.ObjectId)
	obiekt, jest := obiektZnajdz(&stan.forma, kod)
	if !jest {
		return shared.StudioObjectRemoveResponse{}, bladWskazaniaStudio(
			"usunięcie obiektu „" + kod + "”, którego dokument nie ma")
	}
	miejsce := aparatWartoscLiczby(obiekt.AnchorOffset)
	if odcinki, pominiete := postacOdcinkiDozwolone(&stan.forma, miejsce, miejsce, autor); len(odcinki) == 0 {
		return shared.StudioObjectRemoveResponse{}, bladWskazaniaStudio(
			"usunięcie obiektu zatrzymane przez blokadę fragmentu: " +
				postacNazwaBlokad(pominiete))
	}

	bilans := shared.StudioActionBalance{Skipped: []shared.StudioSkippedItem{}}
	for _, zalezny := range aparatOdwolaniaDoCelu(&stan.forma, kod) {
		bilans.Skipped = append(bilans.Skipped, shared.StudioSkippedItem{
			Reason: "odwołanie straciło cel",
			Detail: postacWskaznikTekstu(aparatNazwaRodzaju(zalezny.Kind) + " „" + zalezny.Id +
				"” prowadził do usuwanego obiektu i został oznaczony jako nieświeży"),
		})
		if err := a.aparatZapiszNieswiezosc(ctx, stan, zalezny.Id, true); err != nil {
			return shared.StudioObjectRemoveResponse{}, err
		}
	}

	skladnica, err := a.postacSkladnica()
	if err != nil {
		return shared.StudioObjectRemoveResponse{}, err
	}
	usunieto, err := skladnica.UsunObiektDokumentu(ctx, kod)
	if err != nil {
		return shared.StudioObjectRemoveResponse{}, bladStudio(err)
	}
	if !usunieto {
		return shared.StudioObjectRemoveResponse{}, bladWskazaniaStudio(
			"obiektu „" + kod + "” nie udało się usunąć — wiersza już nie ma")
	}
	obiektUsunBlok(&stan.forma, kod)
	if err := a.aparatZnaczNieswiezoscRodzaju(ctx, stan,
		shared.StudioApparatusKindFigureIndex); err != nil {

		return shared.StudioObjectRemoveResponse{}, err
	}
	if err := a.postacWczytajWiersze(ctx, stan); err != nil {
		return shared.StudioObjectRemoveResponse{}, err
	}

	// Bajty obiektu zostają w magazynie zasobów rdzenia: ten sam zasób bywa
	// wstawiony w kilku dokumentach, więc usunięcie obiektu nie ma prawa kasować
	// treści, na którą powołuje się ktoś inny. Bilans mówi to wprost.
	if obiekt.AssetId != nil && *obiekt.AssetId != "" {
		bilans.Skipped = append(bilans.Skipped, shared.StudioSkippedItem{
			Reason: "zasób magazynu zostaje",
			Detail: postacWskaznikTekstu("bajty obiektu leżą w magazynie zasobów rdzenia " +
				"pod wskazaniem " + *obiekt.AssetId + " i nie zostały usunięte — ten sam " +
				"zasób bywa wstawiony w innych dokumentach"),
		})
	}
	bilans.Applied = 1
	bilans.Note = postacWskaznikTekstu(obiektNazwaRodzaju(obiekt.Kind) + " usunięty z dokumentu")
	stan.opisCzynnosci = "usunięcie obiektu: " + obiektNazwaRodzaju(obiekt.Kind)

	forma, bilansGotowy, _, err := a.postacZakoncz(ctx, stan, autor,
		shared.StudioChangeKindUsuniecie, shared.StudioActionKindObjectChange,
		miejsce, miejsce, bilans)
	if err != nil {
		return shared.StudioObjectRemoveResponse{}, err
	}
	return shared.StudioObjectRemoveResponse{
		Removed: true, Form: forma, Balance: bilansGotowy, ActionId: stan.czynnosc,
	}, nil
}

// ── Pochodzenie i bajty ─────────────────────────────────────────────────────

// obiektUstalPochodzenie rozstrzyga, skąd obiekt jest, i odkłada jego bajty
// w magazynie zasobów rdzenia, gdy przyszły plikiem albo wprost z klienta.
//
// Kolejność jest rozstrzygnięciem: wskazanie zasobu bije wszystko inne, bo
// zasób już leży w magazynie i nie ma po co odkładać go drugi raz.
func (a *adapterStudia) obiektUstalPochodzenie(ctx context.Context,
	z shared.StudioObjectInsertRequest, obiekt *shared.StudioDocumentObject,
	bilans *shared.StudioActionBalance) error {

	switch obiekt.Kind {
	case shared.StudioObjectKindShape:
		return a.obiektPochodzenieKsztaltu(ctx, z, obiekt)
	case shared.StudioObjectKindIcon:
		return a.obiektPochodzenieIkony(ctx, z, obiekt)
	case shared.StudioObjectKindTextbox:
		if z.InnerText == nil || strings.TrimSpace(*z.InnerText) == "" {
			return bladWskazaniaStudio("pole tekstowe bez treści — pole tekstowe wymaga " +
				"brzmienia (pole innerText)")
		}
		obiekt.Source = obiektWskaznikZrodla(shared.StudioObjectSourceDrawn)
		return nil
	case shared.StudioObjectKindChart:
		// Odmowa nazwana, nie cicha: rdzeń nie ma rachunku wykresu i nie udaje,
		// że ma. Wykres składa się w module Design i wstawia jako obraz albo
		// węzeł Designu.
		return bladWskazaniaStudio("wstawienie wykresu — rdzeń nie ma rachunku wykresu " +
			"po stronie modułu Studio. Droga, która działa: złożyć wykres w module Design " +
			"i wstawić go jako obiekt rodzaju image ze wskazaniem zasobu (assetId) albo " +
			"węzła Designu (designNodeId). Brak jest po stronie rdzenia, nie po stronie " +
			"Operatora")
	}
	return a.obiektPochodzenieObrazu(ctx, z, obiekt, bilans)
}

// obiektPochodzenieObrazu ustala pochodzenie obrazu i logo.
func (a *adapterStudia) obiektPochodzenieObrazu(ctx context.Context,
	z shared.StudioObjectInsertRequest, obiekt *shared.StudioDocumentObject,
	bilans *shared.StudioActionBalance) error {

	switch {
	case obiektPodane(z.AssetId):
		kod := strings.TrimSpace(*z.AssetId)
		if _, err := a.sciezkaZasobu(ctx, kod); err != nil {
			return err
		}
		obiekt.AssetId = postacWskaznikTekstu(kod)
		obiekt.Source = obiektWskaznikZrodla(shared.StudioObjectSourceCoreAsset)
		if z.Source != nil && *z.Source == shared.StudioObjectSourcePhotoBank {
			// Zdjęcie z bazy zdjęciowej leży w magazynie jak każdy inny zasób;
			// pochodzenie „baza zdjęciowa" zostaje, bo mówi o prawach do obrazu.
			obiekt.Source = obiektWskaznikZrodla(shared.StudioObjectSourcePhotoBank)
		}
		return nil

	case obiektPodane(z.DesignNodeId):
		kod := strings.TrimSpace(*z.DesignNodeId)
		if a.zasoby == nil {
			return postacBladZaplecza("wstawienie obiektu z modułu Design nie ma drogi — " +
				"rdzeń złożony bez repozytorium Designu")
		}
		if _, err := a.zasoby.SciezkaWektorowaDesignuPoKodzie(ctx, kod); err != nil {
			return bladWskazaniaStudio("węzła modułu Design „" + kod + "” nie ma w rdzeniu; " +
				"kształt zakłada się komendą design.vector.shape.add, a wynik wstawia się " +
				"tutaj jego wskazaniem")
		}
		obiekt.DesignNodeId = postacWskaznikTekstu(kod)
		obiekt.Source = obiektWskaznikZrodla(shared.StudioObjectSourceDesignModule)
		return nil

	case obiektPodane(z.LibraryFileId):
		kod := strings.TrimSpace(*z.LibraryFileId)
		if a.biblioteka == nil {
			return postacBladZaplecza("wstawienie obiektu z Biblioteki nie ma drogi — " +
				"rdzeń złożony bez repozytorium Biblioteki")
		}
		plik, err := a.biblioteka.Plik(ctx, kod)
		if err != nil {
			return bladWskazaniaStudio("pliku Biblioteki „" + kod + "” nie ma w rdzeniu")
		}
		obiekt.LibraryFileId = postacWskaznikTekstu(kod)
		obiekt.Source = obiektWskaznikZrodla(shared.StudioObjectSourceLibraryFile)
		if plik.TrescOdwolanie == nil {
			bilans.Skipped = append(bilans.Skipped, shared.StudioSkippedItem{
				Reason: "plik Biblioteki bez odwołania do bajtów",
				Detail: postacWskaznikTekstu("obiekt niesie pochodzenie z pliku " + kod +
					", ale bajtów w magazynie nie ma — podgląd obiektu będzie pusty, " +
					"dopóki plik nie dostanie treści"),
			})
		}
		return nil

	case obiektPodane(z.BytesBase64):
		bajty, err := base64.StdEncoding.DecodeString(strings.TrimSpace(*z.BytesBase64))
		if err != nil {
			return bladWskazaniaStudio("bajty obiektu w zapisie base64 są nieczytelne: " +
				err.Error())
		}
		return a.obiektOdlozBajty(ctx, bajty, obiektNazwaPliku(z), obiekt,
			shared.StudioObjectSourceFile)

	case obiektPodane(z.Path):
		sciezka := strings.TrimSpace(*z.Path)
		wiadomosc, err := os.Stat(sciezka)
		if err != nil {
			return bladWskazaniaStudio("pliku „" + sciezka + "” rdzeń nie widzi: " + err.Error())
		}
		if wiadomosc.Size() > obiektGranicaBajtow {
			return bladWskazaniaStudio("plik „" + sciezka + "” ma " +
				strconv.FormatInt(wiadomosc.Size(), 10) + " bajtów i przekracza granicę " +
				strconv.Itoa(obiektGranicaBajtow>>20) + " MB przyjętą dla obiektu dokumentu")
		}
		bajty, err := os.ReadFile(sciezka)
		if err != nil {
			return bladWskazaniaStudio("pliku „" + sciezka + "” nie da się odczytać: " +
				err.Error())
		}
		return a.obiektOdlozBajty(ctx, bajty, obiektNazwaPliku(z), obiekt,
			shared.StudioObjectSourceFile)

	case obiektPodane(z.SourceUrl):
		// Adres zapisuje się jako pochodzenie, ale bajtów rdzeń stąd nie
		// pobiera: pobranie treści ze sieci ma w rdzeniu własną drogę
		// (`studio.ingest.url`), a druga byłaby drugą prawdą o tym, co i skąd
		// weszło do dokumentu.
		obiekt.SourceUrl = postacWskaznikTekstu(strings.TrimSpace(*z.SourceUrl))
		obiekt.Source = obiektWskaznikZrodla(shared.StudioObjectSourceWeb)
		bilans.Skipped = append(bilans.Skipped, shared.StudioSkippedItem{
			Reason: "bajty ze sieci nie weszły tą drogą",
			Detail: postacWskaznikTekstu("obiekt niesie zapisane pochodzenie (adres), ale " +
				"treści rdzeń tą komendą nie pobiera; obraz ze sieci wciąga się przez " +
				"studio.ingest.url, a jego zasób wstawia się tutaj polem assetId"),
		})
		return nil
	}

	return bladWskazaniaStudio("wstawienie " + obiektNazwaRodzaju(obiekt.Kind) +
		" bez zapisanego pochodzenia. Wskaż jedno z: zasób magazynu rdzenia (assetId), " +
		"węzeł modułu Design (designNodeId), plik Biblioteki (libraryFileId), ścieżkę " +
		"pliku (path), bajty wprost (bytesBase64) albo adres źródła (sourceUrl). " +
		"Obiekt bez pochodzenia jest brakiem, nie skrótem: za tydzień nikt nie odtworzy, " +
		"na czym pismo się opiera")
}

// obiektPochodzenieKsztaltu ustala pochodzenie kształtu.
func (a *adapterStudia) obiektPochodzenieKsztaltu(ctx context.Context,
	z shared.StudioObjectInsertRequest, obiekt *shared.StudioDocumentObject) error {

	if obiektPodane(z.DesignNodeId) {
		kod := strings.TrimSpace(*z.DesignNodeId)
		if a.zasoby == nil {
			return postacBladZaplecza("wstawienie kształtu z modułu Design nie ma drogi — " +
				"rdzeń złożony bez repozytorium Designu")
		}
		if _, err := a.zasoby.SciezkaWektorowaDesignuPoKodzie(ctx, kod); err != nil {
			return bladWskazaniaStudio("kształtu modułu Design „" + kod + "” nie ma " +
				"w rdzeniu; kształt o ścieżkach edytowalnych zakłada komenda " +
				"design.vector.shape.add i dopiero jej wynik wstawia się do dokumentu")
		}
		obiekt.DesignNodeId = postacWskaznikTekstu(kod)
		obiekt.Source = obiektWskaznikZrodla(shared.StudioObjectSourceDesignModule)
		return nil
	}
	if z.ShapeKind == nil {
		nazwy := make([]string, 0, 8)
		for _, rodzaj := range shared.WartosciStudioShapeKind() {
			nazwy = append(nazwy, string(rodzaj))
		}
		return bladWskazaniaStudio("wstawienie kształtu bez wskazania jego rodzaju " +
			"(pole shapeKind); wykaz: " + strings.Join(nazwy, ", ") +
			". Kształt o ścieżkach edytowalnych wskazuje się węzłem Designu (designNodeId)")
	}
	obiekt.ShapeKind = z.ShapeKind
	obiekt.Source = obiektWskaznikZrodla(shared.StudioObjectSourceDrawn)
	return nil
}

// obiektPochodzenieIkony bierze ikonę z KATALOGU modułu Design i odkłada jej
// rysunek w magazynie zasobów rdzenia.
//
// Katalog jest jeden — ten sam, którym jedzie `design.icon.library.search`.
// Drugi katalog ikon w Studiu rozjechałby się z Designem przy pierwszym
// uzupełnieniu wykazu.
func (a *adapterStudia) obiektPochodzenieIkony(ctx context.Context,
	z shared.StudioObjectInsertRequest, obiekt *shared.StudioDocumentObject) error {

	if !obiektPodane(z.IconName) {
		return bladWskazaniaStudio("wstawienie ikony bez wskazania jej nazwy (pole " +
			"iconName); nazwy bierze się z katalogu ikon rdzenia " +
			"(design.icon.library.search)")
	}
	nazwa := strings.TrimSpace(*z.IconName)
	wzor, jest := ikonaKataloguDesignu(nazwa)
	if !jest {
		wzor, jest = wzorDlaPojeciaDesignu(nazwa)
	}
	if !jest {
		return bladWskazaniaStudio("ikony „" + nazwa + "” nie ma w katalogu ikon rdzenia; " +
			"nazwę wybiera się z wykazu komendy design.icon.library.search, a ikonę " +
			"nieobecną w katalogu składa design.icon.generate")
	}
	svg := svgIkonyKataloguDesignu(wzor, 0, 0)
	if err := a.obiektOdlozBajty(ctx, []byte(svg), "ikona-"+wzor.Nazwa+".svg", obiekt,
		shared.StudioObjectSourceDesignModule); err != nil {

		return err
	}
	if obiekt.AltText == nil {
		obiekt.AltText = postacWskaznikTekstu("ikona " + wzor.Nazwa)
	}
	return nil
}

// obiektOdlozBajty utrwala bajty obiektu w magazynie zasobów rdzenia i zapisuje
// jego pochodzenie.
func (a *adapterStudia) obiektOdlozBajty(ctx context.Context, bajty []byte, nazwa string,
	obiekt *shared.StudioDocumentObject, zrodlo shared.StudioObjectSource) error {

	if len(bajty) == 0 {
		return bladWskazaniaStudio("wstawienie obiektu o zerowej długości — nie ma czego " +
			"osadzić w dokumencie")
	}
	if len(bajty) > obiektGranicaBajtow {
		return bladWskazaniaStudio("obiekt o " + strconv.Itoa(len(bajty)) +
			" bajtach przekracza granicę " + strconv.Itoa(obiektGranicaBajtow>>20) +
			" MB przyjętą dla obiektu dokumentu")
	}
	zasob, err := a.odlozTrescStudia(ctx, bajty, nazwa, obiektFormatPliku(nazwa), "")
	if err != nil {
		return err
	}
	obiekt.AssetId = postacWskaznikTekstu(zasob.Id)
	obiekt.Source = obiektWskaznikZrodla(zrodlo)
	return nil
}

// ── Drobne rachunki ─────────────────────────────────────────────────────────

// obiektPodane mówi, czy pole nieobowiązkowe niesie wartość.
func obiektPodane(pole *string) bool {
	return pole != nil && strings.TrimSpace(*pole) != ""
}

// obiektWskaznikZrodla oddaje wskaźnik na pochodzenie obiektu.
func obiektWskaznikZrodla(zrodlo shared.StudioObjectSource) *shared.StudioObjectSource {
	kopia := zrodlo
	return &kopia
}

// obiektZnajdz odnajduje obiekt w postaci dokumentu.
func obiektZnajdz(forma *shared.StudioDocumentForm, kod string) (shared.StudioDocumentObject, bool) {
	szukany := strings.TrimSpace(kod)
	for _, obiekt := range forma.Objects {
		if obiekt.Id == szukany {
			return obiekt, true
		}
	}
	return shared.StudioDocumentObject{}, false
}

// obiektUsunBlok zdejmuje z drzewa blok wskazujący obiekt.
func obiektUsunBlok(forma *shared.StudioDocumentForm, kod string) {
	nowe := make([]shared.StudioDocumentBlock, 0, len(forma.Blocks))
	for _, blok := range forma.Blocks {
		if blok.Kind == blokPostaciObiekt && blok.ObjectId != nil && *blok.ObjectId == kod {
			continue
		}
		nowe = append(nowe, blok)
	}
	forma.Blocks = nowe
}

// obiektNastepnaWarstwa oddaje warstwę o jeden wyżej od najwyższej zastanej —
// obiekt wstawiony staje NAD tym, co już jest, a nie pod spodem.
func obiektNastepnaWarstwa(forma *shared.StudioDocumentForm) int {
	najwyzsza := 0
	for _, obiekt := range forma.Objects {
		if obiekt.ZOrder != nil && *obiekt.ZOrder >= najwyzsza {
			najwyzsza = *obiekt.ZOrder + 1
		}
	}
	return najwyzsza
}

// obiektWymiary oddaje wymiary obiektu w milimetrach.
func obiektWymiary(obiekt shared.StudioDocumentObject) (float64, float64) {
	szerokosc, wysokosc := 0.0, 0.0
	if obiekt.WidthMm != nil {
		szerokosc = *obiekt.WidthMm
	}
	if obiekt.HeightMm != nil {
		wysokosc = *obiekt.HeightMm
	}
	return szerokosc, wysokosc
}

// obiektDomyslneWymiary nadaje obiektowi rozmiar, gdy Operator go nie podał.
//
// Obiekt o zerowym rozmiarze byłby obiektem niewidzialnym — Operator zobaczyłby
// odpowiedź „wstawiono" i puste miejsce w dokumencie.
func obiektDomyslneWymiary(obiekt *shared.StudioDocumentObject) {
	if obiekt.WidthMm == nil || *obiekt.WidthMm <= 0 {
		obiekt.WidthMm = postacWskaznikMiary(obiektSzerokoscDomyslnaMm)
	}
	if obiekt.HeightMm == nil || *obiekt.HeightMm <= 0 {
		obiekt.HeightMm = postacWskaznikMiary(obiektWysokoscDomyslnaMm)
	}
}

// obiektSprawdzRodzaj odrzuca rodzaj obiektu, którego kontrakt nie zna.
func obiektSprawdzRodzaj(rodzaj shared.StudioObjectKind) error {
	for _, znany := range shared.WartosciStudioObjectKind() {
		if rodzaj == znany {
			return nil
		}
	}
	nazwy := make([]string, 0, 6)
	for _, znany := range shared.WartosciStudioObjectKind() {
		nazwy = append(nazwy, string(znany))
	}
	return bladWskazaniaStudio("obiekt rodzaju „" + string(rodzaj) +
		"”, którego rdzeń nie zna; wykaz: " + strings.Join(nazwy, ", "))
}

// obiektNazwaRodzaju nazywa rodzaj obiektu pełnym słowem.
func obiektNazwaRodzaju(rodzaj shared.StudioObjectKind) string {
	switch rodzaj {
	case shared.StudioObjectKindImage:
		return "obraz"
	case shared.StudioObjectKindShape:
		return "kształt"
	case shared.StudioObjectKindIcon:
		return "ikona"
	case shared.StudioObjectKindTextbox:
		return "pole tekstowe"
	case shared.StudioObjectKindLogo:
		return "logo"
	case shared.StudioObjectKindChart:
		return "wykres"
	default:
		return "obiekt"
	}
}

// obiektZapisPochodzenia opisuje pochodzenie obiektu zdaniem dla Operatora.
func obiektZapisPochodzenia(obiekt shared.StudioDocumentObject) string {
	czesci := make([]string, 0, 3)
	if obiekt.Source != nil {
		czesci = append(czesci, obiektNazwaZrodla(*obiekt.Source))
	}
	switch {
	case obiekt.AssetId != nil && *obiekt.AssetId != "":
		czesci = append(czesci, "zasób "+*obiekt.AssetId)
	case obiekt.DesignNodeId != nil && *obiekt.DesignNodeId != "":
		czesci = append(czesci, "węzeł Designu "+*obiekt.DesignNodeId)
	case obiekt.LibraryFileId != nil && *obiekt.LibraryFileId != "":
		czesci = append(czesci, "plik Biblioteki "+*obiekt.LibraryFileId)
	case obiekt.SourceUrl != nil && *obiekt.SourceUrl != "":
		czesci = append(czesci, *obiekt.SourceUrl)
	}
	if len(czesci) == 0 {
		return "rysunek złożony w dokumencie"
	}
	return strings.Join(czesci, ", ")
}

// obiektNazwaZrodla nazywa pochodzenie obiektu pełnym słowem.
func obiektNazwaZrodla(zrodlo shared.StudioObjectSource) string {
	switch zrodlo {
	case shared.StudioObjectSourceFile:
		return "plik Operatora"
	case shared.StudioObjectSourceCoreAsset:
		return "magazyn zasobów rdzenia"
	case shared.StudioObjectSourceDesignModule:
		return "moduł Design"
	case shared.StudioObjectSourcePhotoBank:
		return "baza zdjęciowa"
	case shared.StudioObjectSourceLibraryFile:
		return "Biblioteka plików"
	case shared.StudioObjectSourceWeb:
		return "adres w sieci"
	case shared.StudioObjectSourceDrawn:
		return "rysunek złożony w dokumencie"
	default:
		return string(zrodlo)
	}
}

// obiektNazwaPliku składa nazwę, pod którą bajty obiektu wejdą do magazynu.
func obiektNazwaPliku(z shared.StudioObjectInsertRequest) string {
	if obiektPodane(z.Path) {
		sciezka := strings.TrimSpace(*z.Path)
		if wskazanie := strings.LastIndexAny(sciezka, "/\\"); wskazanie >= 0 &&
			wskazanie+1 < len(sciezka) {

			return sciezka[wskazanie+1:]
		}
		return sciezka
	}
	return "obiekt-dokumentu-" + string(z.Kind)
}

// obiektFormatPliku odczytuje format z nazwy pliku; brak rozszerzenia znaczy
// zapis dwójkowy nienazwany.
func obiektFormatPliku(nazwa string) string {
	if wskazanie := strings.LastIndex(nazwa, "."); wskazanie >= 0 && wskazanie+1 < len(nazwa) {
		return strings.ToLower(nazwa[wskazanie+1:])
	}
	return "bin"
}

// obiektDoWiersza przekłada obiekt na wiersz warstwy danych.
//
// Kolumny biorą to, po czym się pyta — rodzaj, pochodzenie, zakotwiczenie,
// warstwę, tekst zastępczy, podpis. Reszta postaci idzie polem JSON. Zapis tego,
// co ma kolumnę, także w JSON-ie dałby dwie prawdy o jednym wierszu.
func obiektDoWiersza(dokumentID int64,
	obiekt shared.StudioDocumentObject) (dane.ObiektDokumentuStudia, error) {

	postac := shared.StudioDocumentObject{
		WidthMm: obiekt.WidthMm, HeightMm: obiekt.HeightMm, Crop: obiekt.Crop,
		Wrap: obiekt.Wrap, PositionXMm: obiekt.PositionXMm, PositionYMm: obiekt.PositionYMm,
		Border: obiekt.Border, ShapeKind: obiekt.ShapeKind, FillColor: obiekt.FillColor,
		StrokeColor: obiekt.StrokeColor, StrokeWidthPt: obiekt.StrokeWidthPt,
		Shadow: obiekt.Shadow, RotationDeg: obiekt.RotationDeg,
		InnerText: obiekt.InnerText, Caption: obiekt.Caption,
	}
	zapis, err := json.Marshal(postac)
	if err != nil {
		return dane.ObiektDokumentuStudia{}, postacBladZaplecza(
			"postaci obiektu nie da się zapisać: " + err.Error())
	}
	wiersz := dane.ObiektDokumentuStudia{
		Kod: obiekt.Id, DokumentID: dokumentID, Rodzaj: string(obiekt.Kind),
		ZasobKod: obiekt.AssetId, DesignWezelKod: obiekt.DesignNodeId,
		BibliotekaPlikKod: obiekt.LibraryFileId, AdresZrodla: obiekt.SourceUrl,
		ZakotwiczeniePozycja: int64(aparatWartoscLiczby(obiekt.AnchorOffset)),
		Warstwa:              int64(aparatWartoscLiczby(obiekt.ZOrder)),
		TekstZastepczy:       obiekt.AltText, TekstWewnetrzny: obiekt.InnerText,
		Podpis: obiekt.Caption, PostacJSON: postacWskaznikTekstu(string(zapis)),
	}
	if obiekt.Source != nil {
		zrodlo := string(*obiekt.Source)
		wiersz.Zrodlo = &zrodlo
	}
	if obiekt.Anchor != nil {
		wiersz.Zakotwiczenie = string(*obiekt.Anchor)
	}
	return wiersz, nil
}
