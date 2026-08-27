import { utworzOknoWarsztatu, type OknoWarsztatu } from './okno-warsztatu';
import type { StanProduktu } from './stan-produktu';
import { WARSZTAT_BACKEND } from './warsztaty-apps';

/**
 * Okno wiodące warstwy usług, nazwane Backend Workspace, powstaje z wiązania
 * wspólnej fabryki `okno-warsztatu.ts` z opisem warsztatu `WARSZTAT_BACKEND`,
 * dzięki czemu formularz warsztatu stoi w drzewie jeden raz.
 */
export function utworzOknoBackendWorkspace(stan: StanProduktu): OknoWarsztatu {
  return utworzOknoWarsztatu(stan, WARSZTAT_BACKEND);
}
