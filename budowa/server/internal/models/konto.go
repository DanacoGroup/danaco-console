package models

// MetadaneKonta opisuje konto albo profil poświadczeń, którym kanał wykonał wywołanie; niesie wyłącznie
// odwołania — nazwę profilu, katalog poświadczeń, etykietę konta, nigdy treść poświadczenia.
type MetadaneKonta struct {
	// Konto — identyfikator konta w rejestrze platformy.
	Konto string `json:"account,omitempty"`
	// Profil — nazwa profilu poświadczeń kanału.
	Profil string `json:"profile,omitempty"`
	// Odwolanie — odwołanie do danych dostępowych; nigdy ich treść.
	Odwolanie string `json:"credentialRef,omitempty"`
	// Powod — dlaczego kanał używa tego konta, na przykład po wyczerpaniu limitu poprzedniego.
	Powod string `json:"reason,omitempty"`
	// Kolejne — czy pozostały konta, na które kanał może się przełączyć.
	Kolejne int `json:"remaining,omitempty"`
}
