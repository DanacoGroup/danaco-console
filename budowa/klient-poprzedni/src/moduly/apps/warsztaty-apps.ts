import { AppWorkspaceLayer } from '../../../../shared/contract';
import { BRAKI_BACKEND, BRAKI_FRONTEND, KODY_OKIEN, NAZWY_OKIEN } from './etykiety-apps';
import type { OpisWarsztatu } from './okno-warsztatu';

/**
 * Czym różnią się oba warsztaty modułu Apps.
 *
 * Opisy stoją osobno od ramy warsztatu, bo rama jest czynnością, a to jest
 * wykazem: dwie warstwy z wyliczenia kontraktu wraz z tym, co przypisano
 * każdej z nich. Dopisanie trzeciej warstwy byłoby zmianą kontraktu, nie
 * zmianą widoku.
 */
export const WARSZTAT_FRONTEND: OpisWarsztatu = {
  kodOkna: KODY_OKIEN.FrontendWorkspace,
  tytul: NAZWY_OKIEN[KODY_OKIEN.FrontendWorkspace] ?? 'Frontend Workspace',
  warstwa: AppWorkspaceLayer.Frontend,
  tytulWykazu: 'Drzewo komponentów warstwy interfejsu',
  objasnienieSciezki:
    'Ścieżka pliku warstwy interfejsu względem katalogu roboczego okna. ' +
    'Kontrakt wymaga jej w każdym zapisie warsztatu.',
  braki: BRAKI_FRONTEND,
};

export const WARSZTAT_BACKEND: OpisWarsztatu = {
  kodOkna: KODY_OKIEN.BackendWorkspace,
  tytul: NAZWY_OKIEN[KODY_OKIEN.BackendWorkspace] ?? 'Backend Workspace',
  warstwa: AppWorkspaceLayer.Backend,
  tytulWykazu: 'Mapa zależności usług',
  objasnienieSciezki:
    'Ścieżka pliku usługi względem katalogu roboczego okna. Konfiguracja usługi ' +
    'idzie tą samą komendą co jej kod — kontrakt nie rozróżnia obu zapisów.',
  braki: BRAKI_BACKEND,
};
