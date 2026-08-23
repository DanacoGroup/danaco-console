import { ChangeKind, type Agent } from '../../../../shared/contract';
import type { Kanal } from '../../protokol/kanal';
import { opisOdmowy } from '../../komponenty/odmowa';
import { utworzZrodloAgentow, type ZrodloAgentow } from './zrodlo-agentow';
import { testowanyAgent } from './testowany-agent';

/**
 * Jeden zbiór ekspertów na cały moduł Agents.
 *
 * Pięć okien modułu — Agent Builder, Model Configuration, Skills Manager,
 * Connectors Manager, Permissions Center — pracuje na tym samym ekspercie
 * czynnym. Gdyby każde okno prowadziło własny wykaz i własny wybór, zmiana
 * modelu bazowego nie odświeżałaby biblioteki, a przypisanie umiejętności
 * dotyczyłoby innego eksperta niż ten pokazany w edytorze.
 *
 * Poza własnym działaniem odświeża wyłącznie zdarzenie: zmiana dokonana na innym
 * urządzeniu konta dociera zdarzeniem `agent.changed` i stan wciąga ją tak samo
 * jak własną. Nie ma tu odpytywania w pętli.
 */
export type FazaStanu = 'ladowanie' | 'gotowe' | 'blad';

export interface StanAgentow {
  /** Źródło komend obszaru `agent.*` — okna wołają je wprost. */
  zrodlo: ZrodloAgentow;
  /** Biblioteka ekspertów w kolejności nadanej przez rdzeń. */
  eksperci(): readonly Agent[];
  /** Ekspert czynny albo `null`, gdy żaden nie jest wybrany. */
  wybrany(): Agent | null;
  /** Wybiera eksperta; `null` zdejmuje wybór. */
  wybierz(idEksperta: string | null): void;
  /** Faza ostatniego odczytu biblioteki. */
  faza(): FazaStanu;
  /** Powód odmowy ostatniego odczytu; pusty, gdy odczyt się udał. */
  powod(): string;
  /** Fraza wyszukiwania obowiązująca w bibliotece. */
  fraza(): string;
  /** Odczytuje bibliotekę z rdzenia; fraza pominięta zostaje bez zmian. */
  odswiez(fraza?: string): Promise<void>;
  /** Wciąga eksperta po własnej zmianie, bez czekania na zdarzenie. */
  wchlon(ekspert: Agent): void;
  /** Powiadamia widok o każdej zmianie stanu. */
  obserwuj(sluchacz: () => void): () => void;
  /** Odłącza subskrypcję zdarzeń rdzenia. */
  rozlacz(): void;
}

export function utworzStanAgentow(kanal: Kanal): StanAgentow {
  const zrodlo = utworzZrodloAgentow(kanal);
  const sluchacze = new Set<() => void>();

  let biblioteka: Agent[] = [];
  let wybor: string | null = null;
  let stanFazy: FazaStanu = 'ladowanie';
  let stanPowodu = '';
  let szukana = '';

  /**
   * Ogłoszenie zmiany stanu — najpierw testowany ekspert, potem widok.
   *
   * Testowany ekspert sięga poza ten moduł. Czat modułu Agents jest czatem
   * testowym i nie ma pamięci sesyjnej: rozmowa znika przy zamknięciu okna oraz
   * przy zmianie testowanego agenta (`profil-modulu.ts` → `pamiecSesyjna: false`
   * dla `agents`). Okno rozmowy stoi obok widoku modułu, na scenie sesji, więc
   * bez tego wiersza nie miałoby jak się dowiedzieć o przełączeniu eksperta.
   *
   * Kolejność jest treścią: rozgłos idzie przed słuchaczami okien modułu, żeby
   * rozmowa zdążyła się wyczyścić i nazwać powód, zanim okna modułu przerysują
   * się na nowego eksperta. Odwrotna kolejność pokazywałaby przez moment czat
   * poprzedniego eksperta w oknach już opisanych nazwiskiem następnego.
   */
  function oglos(): void {
    const wybrany = biblioteka.find((wpis) => wpis.id === wybor) ?? null;
    testowanyAgent.ustaw(wybrany?.id ?? '', wybrany?.name ?? '');
    for (const sluchacz of [...sluchacze]) sluchacz();
  }

  function wchlon(ekspert: Agent): void {
    const pozycja = biblioteka.findIndex((wpis) => wpis.id === ekspert.id);
    if (pozycja === -1) biblioteka = [...biblioteka, ekspert];
    else biblioteka = biblioteka.map((wpis) => (wpis.id === ekspert.id ? ekspert : wpis));
    if (wybor === null) wybor = ekspert.id;
    stanFazy = 'gotowe';
    oglos();
  }

  function usunZeZbioru(idEksperta: string): void {
    biblioteka = biblioteka.filter((wpis) => wpis.id !== idEksperta);
    if (wybor === idEksperta) wybor = biblioteka[0]?.id ?? null;
    oglos();
  }

  const odsubskrybuj = zrodlo.naZmiane((tresc) => {
    if (tresc.change === ChangeKind.Deleted) {
      usunZeZbioru(tresc.agent.id);
      return;
    }
    wchlon(tresc.agent);
  });

  return {
    zrodlo,

    eksperci: () => biblioteka,

    wybrany: () => biblioteka.find((wpis) => wpis.id === wybor) ?? null,

    wybierz(idEksperta) {
      if (wybor === idEksperta) return;
      wybor = idEksperta;
      oglos();
    },

    faza: () => stanFazy,

    powod: () => stanPowodu,

    fraza: () => szukana,

    async odswiez(fraza) {
      if (fraza !== undefined) szukana = fraza;
      stanFazy = 'ladowanie';
      stanPowodu = '';
      oglos();

      const wynik = await zrodlo.wykaz(szukana, false);
      if (!wynik.udany || wynik.wynik === undefined) {
        stanFazy = 'blad';
        stanPowodu = opisOdmowy(
          'Odczyt biblioteki ekspertów',
          wynik.blad?.code,
          wynik.blad?.message,
        );
        oglos();
        return;
      }
      biblioteka = wynik.wynik.agents;
      if (wybor !== null && !biblioteka.some((wpis) => wpis.id === wybor)) wybor = null;
      if (wybor === null) wybor = biblioteka[0]?.id ?? null;
      stanFazy = 'gotowe';
      stanPowodu = '';
      oglos();
    },

    wchlon,

    obserwuj(sluchacz) {
      sluchacze.add(sluchacz);
      return () => sluchacze.delete(sluchacz);
    },

    rozlacz() {
      odsubskrybuj();
      sluchacze.clear();
    },
  };
}
