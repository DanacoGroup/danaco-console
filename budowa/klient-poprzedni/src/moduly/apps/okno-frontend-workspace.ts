import { utworzOknoWarsztatu, type OknoWarsztatu } from './okno-warsztatu';
import type { StanProduktu } from './stan-produktu';
import { WARSZTAT_FRONTEND } from './warsztaty-apps';

/**
 * Frontend Workspace — okno wiodące warstwy interfejsu.
 *
 * Plik jest wiązaniem, nie drugim widokiem: formularz warsztatu stoi raz, w
 * `okno-warsztatu.ts`, a stąd dostaje wyłącznie stan produktu i opis warsztatu
 * frontendu.
 */
export function utworzOknoFrontendWorkspace(stan: StanProduktu): OknoWarsztatu {
  return utworzOknoWarsztatu(stan, WARSZTAT_FRONTEND);
}
