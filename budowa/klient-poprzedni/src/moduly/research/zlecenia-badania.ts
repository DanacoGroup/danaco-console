import type {
  ExportFormat,
  ResearchCredibility,
  ResearchFindingStatus,
  ResearchReportSection,
  ResearchSourceKind,
} from '../../../../shared/contract';

/**
 * Zlecenia okien Research w kształcie, w jakim składa je Operator.
 *
 * Jedna odpowiedzialność: nazwanie tego, co okno zebrało z formularza, zanim
 * warstwa kontraktu przełoży to na treść żądania. Rozdzielenie jest celowe:
 * okno mówi po polsku o zakresie badania i wiarygodności źródła, a plik źródła
 * komend mówi nazwami kontraktu — dzięki temu zmiana nazwy pola w kontrakcie
 * dotyka jednego pliku, nie pięciu okien.
 */

/** Zakres badania definiowany w Research Workspace. */
export interface ZlecenieZakresu {
  zakres: string;
  /** Etapy badania po jednym w wierszu; puste wiersze są pomijane. */
  etapy: readonly string[];
}

/** Źródło katalogowane w Sources Manager wraz z metadanymi i oceną. */
export interface ZlecenieZrodla {
  idOkna: string;
  tytul: string;
  rodzaj: ResearchSourceKind;
  adres: string;
  pochodzenie: string;
  wiarygodnosc: ResearchCredibility;
  idPlikuRepozytorium: string;
}

/** Ustalenie zapisywane w Findings Panel; `idUstalenia` puste zakłada nowe. */
export interface ZlecenieUstalenia {
  idOkna: string;
  tresc: string;
  idUstalenia: string;
  idZrodel: readonly string[];
  stan: ResearchFindingStatus;
}

/** Raport składany w Report Builder. */
export interface ZlecenieRaportu {
  idOkna: string;
  idRaportu: string;
  tytul: string;
  idUstalen: readonly string[];
  sekcje: readonly ResearchReportSection[];
}

/** Eksport raportu z Export Panel. */
export interface ZlecenieEksportu {
  idRaportu: string;
  format: ExportFormat;
  sciezka: string;
  doRepozytorium: boolean;
}
