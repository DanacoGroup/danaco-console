package models

// MetadaneKonta opisują konto albo profil poświadczeń, którym kanał wykonał
// wywołanie. Fragment tego rodzaju nadaje kanał także w chwili przełączenia
// konta w trakcie sesji — rotacja ma być widoczna, a nie milcząca.
//
// Struktura niesie wyłącznie odwołania: nazwę profilu, katalog poświadczeń,
// etykietę konta. Treść poświadczenia nie trafia tu nigdy.
type MetadaneKonta struct {
	// Konto — identyfikator konta w rejestrze platformy.
	Konto string `json:"account,omitempty"`
	// Profil — nazwa profilu poświadczeń kanału.
	Profil string `json:"profile,omitempty"`
	// Odwolanie — odwołanie do danych dostępowych; nigdy ich treść.
	Odwolanie string `json:"credentialRef,omitempty"`
	// Powod — dlaczego kanał używa właśnie tego konta; wypełniany przy zmianie
	// konta, na przykład po wyczerpaniu limitu.
	Powod string `json:"reason,omitempty"`
	// Kolejne — czy pozostały konta, na które kanał może się przełączyć.
	Kolejne int `json:"remaining,omitempty"`
}
