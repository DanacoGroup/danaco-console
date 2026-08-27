import type {
  ExportFormat,
  ResearchCredibility,
  ResearchFindingStatus,
  ResearchReportSection,
  ResearchSourceKind,
} from '../../../../shared/contract';

/** Zlecenie okna Research w kształcie, w jakim składa je Operator: zakres badania definiowany w Research Workspace, zanim trafi do kontraktu. */
export interface ZlecenieZakresu {
  zakres: string;
  /** Etapy badania po jednym w wierszu; puste wiersze są pomijane. */
  etapy: readonly string[];
}

/** Źródło katalogowane w Sources Manager wraz z metadanymi, pochodzeniem i oceną wiarygodności, zanim trafi do kontraktu. */
export interface ZlecenieZrodla {
  idOkna: string;
  tytul: string;
  rodzaj: ResearchSourceKind;
  adres: string;
  pochodzenie: string;
  wiarygodnosc: ResearchCredibility;
  idPlikuRepozytorium: string;
}

/** Ustalenie zapisywane w Findings Panel wraz z powiązanymi źródłami i stanem; puste pole identyfikatora ustalenia zakłada nowy wpis. */
export interface ZlecenieUstalenia {
  idOkna: string;
  tresc: string;
  idUstalenia: string;
  idZrodel: readonly string[];
  stan: ResearchFindingStatus;
}

/** Raport składany w Report Builder z tytułu, zaznaczonych ustaleń i sekcji redakcji, zanim trafi do kontraktu jako żądanie zapisu. */
export interface ZlecenieRaportu {
  idOkna: string;
  idRaportu: string;
  tytul: string;
  idUstalen: readonly string[];
  sekcje: readonly ResearchReportSection[];
}

/** Eksport raportu z Export Panel: format dokumentowy, ścieżka docelowa oraz nastawa zapisu do repozytorium. */
export interface ZlecenieEksportu {
  idRaportu: string;
  format: ExportFormat;
  sciezka: string;
  doRepozytorium: boolean;
}
