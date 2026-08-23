import { ChangeKind } from '../../../../shared/contract';
import type { Kanal } from '../../protokol/kanal';
import { utworzOdczytBadania } from './odczyt-badania';
import { utworzPamiecBadania, type PamiecBadania } from './pamiec-badania';
import { utworzWskazanieLektury, type WskazanieLektury } from './wskazanie-lektury';
import { utworzZaznaczenie, type ZaznaczeniePozycji } from './zaznaczenie-pozycji';
import { utworzZrodloLektury, type ZrodloLektury } from './zrodlo-lektury';
import { utworzZrodloOdkrywania, type ZrodloOdkrywania } from './zrodlo-odkrywania';
import { utworzZrodloOknaBadania, type ZrodloOknaBadania } from './zrodlo-okna-badania';
import { utworzZrodloResearch, type ZrodloResearch } from './zrodlo-research';

/**
 * Jedno badanie na cały moduł Research.
 *
 * Siedem okien — Research Workspace, Discovery Panel, Sources Manager,
 * Reading View, Findings Panel, Report Builder, Export Panel — pracuje na tym
 * samym oknie badania, tym samym kompletem źródeł i tych samych ustaleniach.
 * Gdyby każde okno prowadziło swój zbiór, powiązanie źródło↔ustalenie↔sekcja
 * raportu dotyczyłoby sześciu różnych bytów.
 *
 * Wykazy narastają z odpowiedzi. Komendy odczytu wykazu źródeł i ustaleń są już
 * w kontrakcie, ale rdzeń nie ma dla nich uchwytów, więc stan trzyma to, co rdzeń
 * potwierdził w tym połączeniu, i nie dopowiada reszty. Po dobudowie uchwytów
 * odczyt dołoży się tutaj, obok `odswiez` — kształt stanu tego nie wymaga.
 *
 * Plik jest złożeniem pamięci, odczytu i dwóch źródeł komend; sam nie trzyma
 * ani jednej wartości.
 */
export type { FazaBadania } from './pamiec-badania';

export interface StanBadania extends PamiecBadania {
  /** Komendy obszaru `research.*`. */
  zrodlo: ZrodloResearch;
  /** Komendy obszaru `window.*` — okno badania, jego stan i panel akcji. */
  okna: ZrodloOknaBadania;
  /** Wyszukiwanie zasilające Discovery Panel — obszary wiedzy i repozytorium. */
  odkrywanie: ZrodloOdkrywania;
  /** Treść materiału czytanego w Reading View — obszar repozytorium. */
  lekturaZrodla: ZrodloLektury;
  /** Zaznaczenie źródeł — wspólne Sources Manager i Findings Panel. */
  wybraneZrodla: ZaznaczeniePozycji;
  /** Zaznaczenie ustaleń — wspólne Findings Panel i Report Builder. */
  wybraneUstalenia: ZaznaczeniePozycji;
  /** Wskazanie materiału do lektury — wspólne Sources Manager i Reading View. */
  lektura: WskazanieLektury;
  /** Odczytuje okno badania i jego stan; wolno wołać wielokrotnie. */
  odswiez(idSesji: string): Promise<void>;
  /** Odłącza subskrypcje kanału. */
  rozlacz(): void;
}

export function utworzStanBadania(kanal: Kanal): StanBadania {
  const zrodlo = utworzZrodloResearch(kanal);
  const okna = utworzZrodloOknaBadania(kanal);
  const odkrywanie = utworzZrodloOdkrywania(kanal);
  const lekturaZrodla = utworzZrodloLektury(kanal);
  const pamiec = utworzPamiecBadania();
  const odswiez = utworzOdczytBadania(okna, pamiec);

  // Zaznaczenie jest nastawą wspólną oknom, więc ogłasza się tą samą drogą co
  // treść. Pamięć go nie zna — nie pochodzi z rdzenia i nie jest treścią
  // badania — więc ogłoszenie ma tu własny rejestr słuchaczy, a `obserwuj`
  // niżej wpisuje słuchacza do obu. Wołający ma jedną subskrypcję na cały stan
  // i nie musi wiedzieć, która zmiana skąd pochodzi.
  const sluchaczeNastawy = new Set<() => void>();
  const ogloszNastawe = (): void => {
    for (const sluchacz of [...sluchaczeNastawy]) sluchacz();
  };

  // Zdarzenie rdzenia jest jedynym odświeżeniem poza własnym działaniem:
  // raport zmieniony na innym urządzeniu konta dociera tą samą drogą.
  //
  // Rdzeń rozgłasza `research.report.changed` bez wskazania sesji, więc
  // zdarzenie dochodzi na każde połączenie. Porównanie `report.windowId`
  // z oknem badania odsiewa raporty cudzych okien; bez niego Export Panel
  // wydałby dokument, którego to okno nie składało.
  const odsubskrybuj = zrodlo.naZmianeRaportu((tresc) => {
    if (pamiec.idOkna() === '' || tresc.report.windowId !== pamiec.idOkna()) return;
    pamiec.wchlonRaport(tresc.change === ChangeKind.Deleted ? null : tresc.report);
  });

  return {
    ...pamiec,
    zrodlo,
    okna,
    odkrywanie,
    lekturaZrodla,
    wybraneZrodla: utworzZaznaczenie(ogloszNastawe),
    wybraneUstalenia: utworzZaznaczenie(ogloszNastawe),
    lektura: utworzWskazanieLektury(ogloszNastawe),
    odswiez,

    /** Jedna subskrypcja na cały stan: treść badania i nastawa zaznaczenia. */
    obserwuj(sluchacz) {
      const odepnijPamiec = pamiec.obserwuj(sluchacz);
      sluchaczeNastawy.add(sluchacz);
      return () => {
        odepnijPamiec();
        sluchaczeNastawy.delete(sluchacz);
      };
    },

    rozlacz() {
      odsubskrybuj();
      zrodlo.rozlacz();
      okna.rozlacz();
      pamiec.zapomnij();
      sluchaczeNastawy.clear();
    },
  };
}
