// Wpięcie trzech komend rodziny `knowledge.*` — portu `Wiedza`, przez który
// rejestr komend rdzenia dociera do adaptera (`adapter_modul_wiedza.go`).
//
// Rodzina nie ma okna w kliencie: wszystkie trzy komendy są narzędziami modelu
// (`shared.NARZEDZIA_MODELU` niesie `danaco_knowledge_index`,
// `danaco_knowledge_search` i `danaco_knowledge_image_search`), czyli drogą,
// którą model podłączony przez CLI sięga po wiedzę Operatora zamiast
// odpowiadać z pamięci. Dlatego port nie rozgłasza żadnego zdarzenia zmiany — nie ma okna, które by je odebrało.
//
// Wskaźnik nie odświeża się sam przy wgraniu pliku. Osadzenie dokumentu to
// sekundy pracy procesora, a pierwsze pobiera rząd gigabajta wag; wpięcie go
// w `library.file.upload` zamieniłoby wgranie pliku w komendę, która czasem
// trwa minutę i czasem odmawia z powodu braku sieci. Budowanie wskaźnika jest
// czynnością osobną i świadomą, zgodnie z opisem `knowledge.index`
// w kontrakcie.
package core

import (
	"context"

	"danacoconsole/shared"
)

// Wiedza jest portem rodziny `knowledge.*`.
type Wiedza interface {
	// Wskaznik buduje wskaźnik znaczenia dla wskazanych zakresów.
	Wskaznik(ctx context.Context, z shared.KnowledgeIndexRequest) (shared.KnowledgeIndexResponse, error)
	// Szukaj oddaje fragmenty najbliższe pytaniu wraz ze źródłem — bez źródła
	// model cytowałby bez możliwości sprawdzenia.
	Szukaj(ctx context.Context, z shared.KnowledgeSearchRequest) (shared.KnowledgeSearchResponse, error)
	// SzukajObrazu oddaje obrazy najbliższe zdaniu, które je opisuje. Osobna
	// komenda, nie zakres tamtej: obraz wraca bez fragmentu do zacytowania,
	// a liczy go model z innej przestrzeni niż wektory tekstu.
	SzukajObrazu(ctx context.Context, z shared.KnowledgeImageSearchRequest) (shared.KnowledgeImageSearchResponse, error)
}

// zarejestrujWiedze wpina trzy komendy rodziny `knowledge.*`.
func zarejestrujWiedze(r *Rejestr, m Wiedza) {
	if r == nil || m == nil {
		return
	}
	r.Zarejestruj(shared.CommandKnowledgeIndex, obsluz(m.Wskaznik))
	r.Zarejestruj(shared.CommandKnowledgeSearch, obsluz(m.Szukaj))
	r.Zarejestruj(shared.CommandKnowledgeImageSearch, obsluz(m.SzukajObrazu))
}
