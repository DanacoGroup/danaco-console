import {
  ChangeKind,
  type AppWorkspaceLayer,
  type AppsWorkspaceChangedEvent,
  type DeveloperFile,
} from '../../../../shared/contract';

/**
 * Pliki warsztatu obu warstw — wszystko, co moduł wie o zawartości warsztatów.
 * Jedno miejsce spotkania odczytu `apps.workspace.list` i zdarzenia
 * `apps.workspace.changed`, żeby okna nie prowadziły dwóch rozjeżdżających się wykazów.
 */
export interface WarsztatPlikow {
  /** Pliki jednej warstwy; kolejność z odpowiedzi rdzenia, dopiski na końcu. */
  pliki(warstwa: AppWorkspaceLayer): readonly DeveloperFile[];
  /** Czy warstwa była już czytana z rdzenia — odróżnia „pusto" od „nie pytano". */
  czyCzytana(warstwa: AppWorkspaceLayer): boolean;
  /** Wchłania wykaz z odpowiedzi `apps.workspace.list` — zastępuje warstwę. */
  wchlonWykaz(warstwa: AppWorkspaceLayer, pliki: readonly DeveloperFile[]): void;
  /** Wchłania zdarzenie `apps.workspace.changed`; usunięcie zdejmuje plik. */
  wchlonZdarzenie(tresc: AppsWorkspaceChangedEvent): void;
  /** Wchłania plik z odpowiedzi `apps.workspace.update` — zapis własnego okna. */
  wchlonZapis(warstwa: AppWorkspaceLayer, plik: DeveloperFile): void;
}

export function utworzWarsztatPlikow(): WarsztatPlikow {
  const wedlugWarstwy = new Map<string, DeveloperFile[]>();
  const czytane = new Set<string>();

  /** Wstawia plik na miejsce pliku o tej samej ścieżce albo dokłada na koniec. */
  function wstaw(warstwa: AppWorkspaceLayer, plik: DeveloperFile): void {
    const biezace = wedlugWarstwy.get(warstwa) ?? [];
    const pozycja = biezace.findIndex((inny) => inny.path === plik.path);
    if (pozycja === -1) wedlugWarstwy.set(warstwa, [...biezace, plik]);
    else {
      wedlugWarstwy.set(
        warstwa,
        biezace.map((inny) => (inny.path === plik.path ? plik : inny)),
      );
    }
  }

  return {
    pliki: (warstwa) => wedlugWarstwy.get(warstwa) ?? [],
    czyCzytana: (warstwa) => czytane.has(warstwa),

    wchlonWykaz(warstwa, pliki) {
      // Zastępuje, nie dokłada: odpowiedź rdzenia jest pełnym stanem warstwy
      // w chwili odczytu.
      wedlugWarstwy.set(warstwa, [...pliki]);
      czytane.add(warstwa);
    },

    wchlonZdarzenie(tresc) {
      // Zdarzenie bez pola pliku nie niesie nic, co dałoby się wstawić.
      if (tresc.file === undefined) return;
      if (tresc.change === ChangeKind.Deleted) {
        const biezace = wedlugWarstwy.get(tresc.layer) ?? [];
        wedlugWarstwy.set(
          tresc.layer,
          biezace.filter((inny) => inny.path !== tresc.file?.path),
        );
        return;
      }
      wstaw(tresc.layer, tresc.file);
    },

    wchlonZapis(warstwa, plik) {
      wstaw(warstwa, plik);
    },
  };
}
