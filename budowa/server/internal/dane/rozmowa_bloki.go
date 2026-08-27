// Plik dokleja zapisane bloki wiadomości do historii okna oddawanej
// kontraktem klienta, bez zmiany samego kontraktu.
package dane

import (
	"context"
	"encoding/json"

	"danacoconsole/shared"
)

// blokKontraktu to jedna pozycja wykazu bloków w metadanych wiadomości,
// niosąca rodzaj, treść, dane i chwilę powstania.
type blokKontraktu struct {
	Rodzaj shared.ChunkKind `json:"kind"`
	Tresc  string           `json:"text,omitempty"`
	Dane   json.RawMessage  `json:"data,omitempty"`
	Chwila int64            `json:"at,omitempty"`
}

// doklejBloki dopisuje bloki okna do wiadomości historii. Błąd odczytu bloków
// nie unieważnia historii — tekst rozmowy ma się wyświetlić także wtedy, gdy
// warstwa bloków zawiodła; wykaz wraca wówczas nietknięty.
func (u *UtrwalaczRozmowy) doklejBloki(ctx context.Context, idOkna string, wykaz []shared.Message) []shared.Message {
	if u.zestaw == nil || u.zestaw.Bloki == nil || len(wykaz) == 0 {
		return wykaz
	}
	bloki, err := u.zestaw.Bloki.ListaOkna(ctx, idOkna)
	if err != nil || len(bloki) == 0 {
		return wykaz
	}
	naWiadomosc := pogrupujBloki(bloki)
	for i := range wykaz {
		zapisane, sa := naWiadomosc[wykaz[i].Id]
		if !sa {
			continue
		}
		wykaz[i].Metadata = dolaczBlokiDoMetadanych(wykaz[i].Metadata, zapisane)
	}
	return wykaz
}

// pogrupujBloki rozkłada wykaz bloków okna na wiadomości. Kolejność w obrębie
// wiadomości pochodzi z zapytania (kolejnosc, id) i zostaje zachowana.
func pogrupujBloki(bloki []BlokWiadomosci) map[string][]blokKontraktu {
	naWiadomosc := map[string][]blokKontraktu{}
	for _, blok := range bloki {
		naWiadomosc[blok.WiadomoscKod] = append(naWiadomosc[blok.WiadomoscKod], blokKontraktu{
			Rodzaj: blok.Rodzaj,
			Tresc:  blok.Tresc,
			Dane:   blok.Ladunek,
			Chwila: blok.Chwila,
		})
	}
	return naWiadomosc
}

// dolaczBlokiDoMetadanych dopisuje klucz `blocks` do obszaru metadanych,
// zachowując pola już złożone (persona, sourceWindowId, tokeny). Obszar
// nieczytelny albo niesklejalny zostaje jak był — bloki znikają z odpowiedzi,
// ale atrybucja i treść nie.
func dolaczBlokiDoMetadanych(meta json.RawMessage, bloki []blokKontraktu) json.RawMessage {
	obszar := map[string]json.RawMessage{}
	if len(meta) > 0 {
		if err := json.Unmarshal(meta, &obszar); err != nil {
			return meta
		}
	}
	surowe, err := json.Marshal(bloki)
	if err != nil {
		return meta
	}
	obszar["blocks"] = surowe
	zlozone, err := json.Marshal(obszar)
	if err != nil {
		return meta
	}
	return zlozone
}
