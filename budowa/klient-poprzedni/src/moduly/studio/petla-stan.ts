import type {
  StudioAgentConflict,
  StudioAgentSettings,
  StudioAgentSlot,
  StudioChain,
  StudioChainStep,
  StudioTaskPlan,
} from '../../../../shared/contract';

/** Stan okna pętli wykonawczej, jeden na okno: wynik wsadu dla jednego dokumentu jest bilansem tej pozycji, nie zbiorczym potwierdzeniem gotowości. */
export interface PozycjaWsadu {
  /** Dokument objęty wsadem. */
  idDokumentu: string;
  /** Stan pozycji w przebiegu. */
  stan: 'oczekuje' | 'w-biegu' | 'udana' | 'odrzucona' | 'przerwana';
  /** Powód odrzucenia albo przerwania; puste, gdy się udało. */
  powod: string;
}

/** Nastawy pętli odczytane z rdzenia jako pełny obiekt, albo null, gdy odczytu jeszcze wcale nie wykonano. */
export type OdczytNastaw = StudioAgentSettings | null;

export interface StanPetli {
  /** Rozkład prowadzony w oknie; `null`, gdy żadnego nie ma. */
  rozklad(): StudioTaskPlan | null;
  ustawRozklad(rozklad: StudioTaskPlan | null): void;
  /** Pozostałe rozkłady dokumentu — do przełączenia się między zleceniami. */
  rozklady(): readonly StudioTaskPlan[];
  ustawRozklady(rozklady: readonly StudioTaskPlan[]): void;

  /** Nastawy pętli odczytane z rdzenia; `null` znaczy „jeszcze nie pytałem". */
  nastawy(): OdczytNastaw;
  ustawNastawy(nastawy: OdczytNastaw): void;
  /** Powód odmowy uruchomienia oddany przez rdzeń; `null`, gdy żadnej nie było. */
  odmowaPrzebiegu(): string | null;
  ustawOdmowePrzebiegu(powod: string | null): void;

  /** Wykonawcy pracujący nad dokumentem wraz z zajętymi fragmentami. */
  obsada(): readonly StudioAgentSlot[];
  ustawObsade(obsada: readonly StudioAgentSlot[]): void;
  /** Spięcia o ten sam fragment wraz z odłożonym brzmieniem. */
  spiecia(): readonly StudioAgentConflict[];
  ustawSpiecia(spiecia: readonly StudioAgentConflict[]): void;

  /** Łańcuchy zapisane w rdzeniu. */
  lancuchy(): readonly StudioChain[];
  ustawLancuchy(lancuchy: readonly StudioChain[]): void;
  /** Łańcuch składany w warsztacie — kroki w kolejności wykonania. */
  skladane(): readonly StudioChainStep[];
  ustawSkladane(kroki: readonly StudioChainStep[]): void;
  /** Nazwa łańcucha składanego. */
  nazwaSkladanego(): string;
  ustawNazweSkladanego(nazwa: string): void;
  /** Łańcuch wskazany do zmiany; puste znaczy nowy. */
  zmieniany(): string;
  ustawZmieniany(idLancucha: string): void;

  /** Pozycje wsadu wraz z wynikiem każdej osobno. */
  wsad(): readonly PozycjaWsadu[];
  ustawWsad(pozycje: readonly PozycjaWsadu[]): void;
  /** Przestawia jedną pozycję wsadu, nie ruszając pozostałych. */
  przestawPozycjeWsadu(idDokumentu: string, stan: PozycjaWsadu['stan'], powod: string): void;
  /** Czy Operator zażądał przerwania wsadu. */
  wsadPrzerwany(): boolean;
  ustawPrzerwanieWsadu(przerwany: boolean): void;

  /** Czy okno pętli jest otwarte. Domyślnie NIE — powierzchnia należy do dokumentu. */
  otwarte(): boolean;
  ustawOtwarte(otwarte: boolean): void;
  /** Czy Operator wybrał stałą kolumnę zamiast nakładki na żądanie. */
  stalaKolumna(): boolean;
  ustawStalaKolumne(stala: boolean): void;

  /** Powiadamia widoki o każdej zmianie stanu. */
  obserwuj(sluchacz: () => void): () => void;
}

export function utworzStanPetli(): StanPetli {
  const sluchacze = new Set<() => void>();
  let biezacyRozklad: StudioTaskPlan | null = null;
  let wszystkieRozklady: readonly StudioTaskPlan[] = [];
  let odczytaneNastawy: OdczytNastaw = null;
  let powodOdmowy: string | null = null;
  let obsadaOkna: readonly StudioAgentSlot[] = [];
  let spieciaOkna: readonly StudioAgentConflict[] = [];
  let lancuchyOkna: readonly StudioChain[] = [];
  let krokiSkladane: readonly StudioChainStep[] = [];
  let nazwaLancucha = '';
  let idZmienianego = '';
  let pozycjeWsadu: readonly PozycjaWsadu[] = [];
  let przerwanie = false;
  // Okno pętli wchodzi na żądanie i schodzi, gdy nie jest używane, bo powierzchnia należy do dokumentu.
  let oknoOtwarte = false;
  let trybStalejKolumny = false;

  function oglos(): void {
    for (const sluchacz of [...sluchacze]) sluchacz();
  }

  return {
    rozklad: () => biezacyRozklad,
    rozklady: () => wszystkieRozklady,
    nastawy: () => odczytaneNastawy,
    odmowaPrzebiegu: () => powodOdmowy,
    obsada: () => obsadaOkna,
    spiecia: () => spieciaOkna,
    lancuchy: () => lancuchyOkna,
    skladane: () => krokiSkladane,
    nazwaSkladanego: () => nazwaLancucha,
    zmieniany: () => idZmienianego,
    wsad: () => pozycjeWsadu,
    wsadPrzerwany: () => przerwanie,
    otwarte: () => oknoOtwarte,
    stalaKolumna: () => trybStalejKolumny,

    ustawRozklad(rozklad) {
      biezacyRozklad = rozklad;
      oglos();
    },

    ustawRozklady(rozklady) {
      wszystkieRozklady = rozklady;
      oglos();
    },

    ustawNastawy(nastawy) {
      odczytaneNastawy = nastawy;
      oglos();
    },

    ustawOdmowePrzebiegu(powod) {
      if (powodOdmowy === powod) return;
      powodOdmowy = powod;
      oglos();
    },

    ustawObsade(obsada) {
      obsadaOkna = obsada;
      oglos();
    },

    ustawSpiecia(spiecia) {
      spieciaOkna = spiecia;
      oglos();
    },

    ustawLancuchy(lancuchy) {
      lancuchyOkna = lancuchy;
      oglos();
    },

    ustawSkladane(kroki) {
      krokiSkladane = kroki;
      oglos();
    },

    ustawNazweSkladanego(nazwa) {
      if (nazwaLancucha === nazwa) return;
      nazwaLancucha = nazwa;
      oglos();
    },

    ustawZmieniany(idLancucha) {
      if (idZmienianego === idLancucha) return;
      idZmienianego = idLancucha;
      oglos();
    },

    ustawWsad(pozycje) {
      pozycjeWsadu = pozycje;
      oglos();
    },

    przestawPozycjeWsadu(idDokumentu, stan, powod) {
      pozycjeWsadu = pozycjeWsadu.map((pozycja) =>
        pozycja.idDokumentu === idDokumentu ? { idDokumentu, stan, powod } : pozycja,
      );
      oglos();
    },

    ustawPrzerwanieWsadu(przerwany) {
      if (przerwanie === przerwany) return;
      przerwanie = przerwany;
      oglos();
    },

    ustawOtwarte(otwarte) {
      if (oknoOtwarte === otwarte) return;
      oknoOtwarte = otwarte;
      oglos();
    },

    ustawStalaKolumne(stala) {
      if (trybStalejKolumny === stala) return;
      trybStalejKolumny = stala;
      oglos();
    },

    obserwuj(sluchacz) {
      sluchacze.add(sluchacz);
      return () => sluchacze.delete(sluchacz);
    },
  };
}
