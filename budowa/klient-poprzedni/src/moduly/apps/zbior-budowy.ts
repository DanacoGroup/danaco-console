import {
  ChangeKind,
  type AppDeployment,
  type AppStage,
  type AppsBuildChangedEvent,
} from '../../../../shared/contract';

/**
 * Zbiór etapów budowy i wdrożeń produktu — wszystko, co moduł wie o przebiegu
 * pracy.
 *
 * Stoi osobno od stanu modułu, bo to inna odpowiedzialność: stan modułu
 * prowadzi okno rdzenia i architekturę oraz ogłasza zmiany oknom, a ten zbiór
 * wyłącznie gromadzi to, co przynoszą zdarzenia i odpowiedzi. Rozdział pozwala
 * sprawdzić gromadzenie bez dotykania kanału.
 *
 * Wdrożenia mają trzy źródła, etapy jedno. Wdrożenie wchodzi tu zdarzeniem
 * `apps.build.changed`, odpowiedzią na `apps.deployment.run` albo wykazem
 * z `apps.deployment.list`. Etapów budowy żadna komenda odczytu nie zwraca,
 * więc ich pusty zbiór na starcie jest stanem prawdziwym, a nie brakiem
 * odczytu. `czyWdrozeniaCzytane()` odróżnia „rdzeń powiedział, że nie ma nic"
 * od „nikt jeszcze nie pytał".
 *
 * Odpowiedź komendy jest starsza niż zdarzenie. `apps.deployment.run` kończy
 * się, gdy przebieg ruszy, a nie gdy się skończy, więc jej odpowiedź niesie
 * migawkę ze stanem `pending`; przejścia `running` → `succeeded`/`failed`
 * przychodzą wyłącznie zdarzeniem. Ciąg dalszy obietnicy `await` biegnie
 * mikrozadaniem, czyli po obsłużeniu ramek, które padły w tej samej turze pętli
 * zdarzeń — zapis z odpowiedzi trafiłby na koniec i cofnął stan końcowy
 * z powrotem do `pending`.
 *
 * Dlatego źródła są rozróżnione, a nie uporządkowane po stanie: ustawianie
 * stanów w drabinkę „który jest dalszy" byłoby zgadywaniem cyklu życia po
 * stronie klienta. Rozstrzyga pochodzenie — strumień zdarzeń jest jedynym
 * źródłem postępu przebiegu, a odpowiedź komendy wyłącznie zasiewa pozycję,
 * której strumień jeszcze nie zgłosił, żeby wiersz pojawił się natychmiast po
 * naciśnięciu. Wdrożenia znanego ze zdarzenia odpowiedź już nie cofnie.
 *
 * Etap bez identyfikatora nie jest etapem. Zdarzenie niesie pole `stage` jako
 * wymagane, więc przy zmianie dotyczącej samego wdrożenia rdzeń wypełnia je
 * zaślepką `{ id: "", name: "", status: "", windowId: <okno> }`. Wpisana do
 * wykazu dawałaby bezimienny wiersz etapu i zdanie „etapów 1" przy zerze etapów
 * prawdziwych. Zaślepkę mijamy; wdrożenie z tej samej ramki bierzemy normalnie.
 *
 * Zbiór liczy też ramki, które minął. Sam pusty wykaz nie odróżnia „nie
 * przyszło jeszcze nic" od „przyszło, ale etapów w tym nie było". Rachunek
 * ramek daje oknu liczby z ramek, które padły, więc zdanie o pustce składa się
 * z nich albo nie pada wcale, zamiast być wpisanym na stałe w napis
 * o zachowaniu rdzenia (wzorzec: `moduly/katalog-okien.ts`).
 */

/**
 * Ile ramek `apps.build.changed` przyszło i co niosły.
 *
 * Liczby dotyczą wyłącznie ramek wchłoniętych przez ten zbiór, więc mówią
 * o tym, co zobaczył ten egzemplarz modułu. Żadna z nich nie orzeka, co rdzeń
 * rozgłasza, a czego nie.
 */
export interface RachunekRamek {
  /** Wszystkie wchłonięte ramki zdarzenia. */
  wszystkie: number;
  /** Ramki, w których pole `stage` niosło identyfikator. */
  zEtapem: number;
  /** Ramki niosące wdrożenie. */
  zWdrozeniem: number;
  /** Ramki, które zdjęły etap ze zbioru (`change: deleted`). */
  zdjeteEtapy: number;
}

export interface ZbiorBudowy {
  etapy(): readonly AppStage[];
  wdrozenia(): readonly AppDeployment[];
  /** Rachunek ramek, które zbiór wchłonął — źródło zdania o powodzie pustki. */
  ramki(): RachunekRamek;
  /** Wchłania zdarzenie zmiany budowy; usunięcie zdejmuje etap ze zbioru. */
  wchlonZdarzenie(tresc: AppsBuildChangedEvent): void;
  /**
   * Zasiewa wdrożenie migawką z odpowiedzi na `apps.deployment.run`.
   * Pozycja znana już ze strumienia zdarzeń zostaje nietknięta (patrz nagłówek).
   */
  wchlonOdpowiedzWdrozenia(wdrozenie: AppDeployment): void;
  /**
   * Wchłania wykaz z odpowiedzi `apps.deployment.list` — historia z bazy.
   *
   * Obowiązuje ta sama reguła pierwszeństwa co przy odpowiedzi komendy:
   * wdrożenie, o którym mówił już strumień zdarzeń, zostaje nietknięte. Odczyt
   * jest migawką z chwili zapytania, a zdarzenie niesie stan z chwili zmiany —
   * wpuszczenie odczytu na wierzch cofałoby `succeeded` do `running`.
   */
  wchlonWykazWdrozen(wdrozenia: readonly AppDeployment[]): void;
  /** Czy wykaz wdrożeń był już czytany z rdzenia — odróżnia „pusto" od „nie pytano". */
  czyWdrozeniaCzytane(): boolean;
}

export function utworzZbiorBudowy(): ZbiorBudowy {
  let etapy: AppStage[] = [];
  let wdrozenia: AppDeployment[] = [];
  /** Wdrożenia, o których mówił już strumień zdarzeń — odpowiedzi ich nie ruszają. */
  const zgloszoneZdarzeniem = new Set<string>();
  /** Czy `apps.deployment.list` wróciło już z odpowiedzią udaną. */
  let czytane = false;
  const rachunek: RachunekRamek = { wszystkie: 0, zEtapem: 0, zWdrozeniem: 0, zdjeteEtapy: 0 };

  function zapiszWdrozenie(wdrozenie: AppDeployment): void {
    // Najnowsze na przedzie: dziennik wydań i tabela wdrożeń czytają go od
    // góry, a rdzeń nie podaje porządku wykazu.
    wdrozenia = [wdrozenie, ...wdrozenia.filter((inne) => inne.id !== wdrozenie.id)];
  }

  return {
    etapy: () => etapy,
    wdrozenia: () => wdrozenia,
    // Kopia, nie odnośnik: rachunek jest odczytem, a nie polem do zapisu
    // z zewnątrz.
    ramki: () => ({ ...rachunek }),

    wchlonZdarzenie(tresc) {
      rachunek.wszystkie += 1;
      if (tresc.stage.id !== '') rachunek.zEtapem += 1;
      if (tresc.deployment !== undefined) rachunek.zWdrozeniem += 1;
      if (tresc.stage.id !== '') {
        if (tresc.change === ChangeKind.Deleted) {
          rachunek.zdjeteEtapy += 1;
          etapy = etapy.filter((etap) => etap.id !== tresc.stage.id);
        } else {
          etapy = [...etapy.filter((inny) => inny.id !== tresc.stage.id), tresc.stage].sort(
            (pierwszy, drugi) => (pierwszy.order ?? 0) - (drugi.order ?? 0),
          );
        }
      }
      if (tresc.deployment !== undefined) {
        zgloszoneZdarzeniem.add(tresc.deployment.id);
        zapiszWdrozenie(tresc.deployment);
      }
    },

    wchlonOdpowiedzWdrozenia(wdrozenie) {
      if (zgloszoneZdarzeniem.has(wdrozenie.id)) return;
      zapiszWdrozenie(wdrozenie);
    },

    wchlonWykazWdrozen(wykaz) {
      czytane = true;
      // Wykaz idzie od najstarszego do przodu, bo `zapiszWdrozenie` kładzie
      // każdą pozycję NA CZELE. Rdzeń oddaje wdrożenia „od najnowszego"
      // (kontrakt, `AppsDeploymentListResponse.deployments`), więc odwrócenie
      // zachowuje ten porządek zamiast wywracać go na drugą stronę.
      for (const wdrozenie of [...wykaz].reverse()) {
        if (zgloszoneZdarzeniem.has(wdrozenie.id)) continue;
        zapiszWdrozenie(wdrozenie);
      }
    },

    czyWdrozeniaCzytane: () => czytane,
  };
}
