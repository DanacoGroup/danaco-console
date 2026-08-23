import { utworzOknoWarsztatu, type OknoWarsztatu } from './okno-warsztatu';
import type { StanProduktu } from './stan-produktu';
import { WARSZTAT_FRONTEND } from './warsztaty-apps';

/**
 * Frontend Workspace — okno wiodące warstwy interfejsu.
 *
 * Plik jest wiązaniem, nie drugim widokiem: formularz warsztatu stoi raz, w
 * `okno-warsztatu.ts`. Kod okna nie pada tu literałem — ramę Apps woła
 * `utworzRameApps(opis.kodOkna, …)` ze zmiennej, żeby `KODY_OKIEN` zostało
 * jedynym miejscem, w którym kody okien Apps są wypisane.
 */
export function utworzOknoFrontendWorkspace(stan: StanProduktu): OknoWarsztatu {
  return utworzOknoWarsztatu(stan, WARSZTAT_FRONTEND);
}
