import { utworzOknoWarsztatu, type OknoWarsztatu } from './okno-warsztatu';
import type { StanProduktu } from './stan-produktu';
import { WARSZTAT_BACKEND } from './warsztaty-apps';

/**
 * Backend Workspace — okno wiodące warstwy usług.
 *
 * Plik jest wiązaniem, nie drugim widokiem: formularz warsztatu stoi raz,
 * w `okno-warsztatu.ts`. Osobny plik daje jedno miejsce w drzewie, po którym
 * widać, że okno `backend-workspace` jest zbudowane, a nie tylko wymienione
 * w słowniku kodów.
 *
 * Kod okna nie pada tu literałem — ramę woła `utworzRameApps(opis.kodOkna, …)`
 * ze zmiennej, żeby `KODY_OKIEN` pozostał jedynym miejscem, w którym kody Apps
 * są zapisane. Okna Studio rozdziela ta sama zasada: `okno-studio.ts` jest
 * wspólną ramą, a poszczególne okna mają własne pliki.
 */
export function utworzOknoBackendWorkspace(stan: StanProduktu): OknoWarsztatu {
  return utworzOknoWarsztatu(stan, WARSZTAT_BACKEND);
}
