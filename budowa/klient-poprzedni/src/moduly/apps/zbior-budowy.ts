import {
  ChangeKind,
  type AppDeployment,
  type AppStage,
  type AppsBuildChangedEvent,
} from '../../../../shared/contract';

/**
 * Zbiór etapów budowy i wdrożeń produktu: gromadzi wszystko, co moduł wie
 * o przebiegu, składając odpowiedzi komend ze strumieniem zdarzeń rdzenia.
 *
 * Rachunek ramek poniżej liczy ramki zmiany budowy wchłonięte przez ten zbiór.
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
   * Zasiewa wdrożenie migawką z odpowiedzi komendy uruchomienia.
   */
  wchlonOdpowiedzWdrozenia(wdrozenie: AppDeployment): void;
  /**
   * Wchłania wykaz wdrożeń z odpowiedzi komendy wykazu.
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
    // Najnowsze na przedzie, bo dziennik i tabela czytają wykaz od góry.
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
      // Wykaz idzie od najstarszego, bo zapis kładzie każdą pozycję na czele.
      for (const wdrozenie of [...wykaz].reverse()) {
        if (zgloszoneZdarzeniem.has(wdrozenie.id)) continue;
        zapiszWdrozenie(wdrozenie);
      }
    },

    czyWdrozeniaCzytane: () => czytane,
  };
}
