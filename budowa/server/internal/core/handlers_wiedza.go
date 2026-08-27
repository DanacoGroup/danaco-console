// Plik wpina trzy komendy rodziny knowledge.* jako port Wiedza. Rodzina nie ma okna w kliencie:
// wszystkie trzy komendy są narzędziami modelu, drogą, którą model sięga po wiedzę Operatora.
package core

import (
	"context"

	"danacoconsole/shared"
)

// Wiedza jest portem rodziny knowledge.* obsługującym wskaźnik znaczenia i wyszukiwanie wiedzy modelu.
type Wiedza interface {
	// Wskaznik buduje wskaźnik znaczenia dla wskazanych zakresów.
	Wskaznik(ctx context.Context, z shared.KnowledgeIndexRequest) (shared.KnowledgeIndexResponse, error)
	// Szukaj oddaje fragmenty najbliższe pytaniu wraz ze źródłem, bez którego model cytowałby ślepo.
	Szukaj(ctx context.Context, z shared.KnowledgeSearchRequest) (shared.KnowledgeSearchResponse, error)
	// SzukajObrazu oddaje obrazy najbliższe zdaniu, licząc je z innej przestrzeni niż wektory tekstu.
	SzukajObrazu(ctx context.Context, z shared.KnowledgeImageSearchRequest) (shared.KnowledgeImageSearchResponse, error)
}

// zarejestrujWiedze wpina trzy komendy rodziny knowledge.* obsługujące wskaźnik i wyszukiwanie wiedzy.
func zarejestrujWiedze(r *Rejestr, m Wiedza) {
	if r == nil || m == nil {
		return
	}
	r.Zarejestruj(shared.CommandKnowledgeIndex, obsluz(m.Wskaznik))
	r.Zarejestruj(shared.CommandKnowledgeSearch, obsluz(m.Szukaj))
	r.Zarejestruj(shared.CommandKnowledgeImageSearch, obsluz(m.SzukajObrazu))
}
