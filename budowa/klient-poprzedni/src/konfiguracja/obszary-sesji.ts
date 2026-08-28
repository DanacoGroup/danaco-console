import {
  ConfigAxis,
  SessionConfigArea,
  SessionConfigSource,
  WorkingDirectoryDegradation,
  type SessionConfigOrigin,
  type WorkingDirectoryResolution,
} from '../../../shared/contract';
import { opisAdresu, type AdresUstawienia } from './adres-ustawienia';

// Słownik obszarów konfiguracji sesji i nazewnictwo pochodzenia obszaru, z kontraktu.

/** Obszary w kolejności pól kontraktu — tej samej, w której odpowiada rdzeń, więc wykaz i odpowiedź się zgadzają. */
export const OBSZARY_SESJI: readonly SessionConfigArea[] = [
  SessionConfigArea.Model,
  SessionConfigArea.Account,
  SessionConfigArea.Provider,
  SessionConfigArea.SystemPrompt,
  SessionConfigArea.Tools,
  SessionConfigArea.Permissions,
  SessionConfigArea.Mcp,
  SessionConfigArea.Hooks,
  SessionConfigArea.Skills,
  SessionConfigArea.Memory,
  SessionConfigArea.ProjectContext,
  SessionConfigArea.ConversationContext,
  SessionConfigArea.WorkingDirectory,
  SessionConfigArea.AdditionalDirectories,
  SessionConfigArea.Environment,
  SessionConfigArea.InputOutput,
  SessionConfigArea.Runtime,
  SessionConfigArea.SessionLifecycle,
];

/** Nazwa obszaru konfiguracji sesji pokazywana Operatorowi w oknie ustawień sesji, zamiast kodu z kontraktu. */
export const NAZWY_OBSZAROW: Readonly<Record<SessionConfigArea, string>> = {
  [SessionConfigArea.Model]: 'model i nakład rozumowania',
  [SessionConfigArea.Account]: 'konto i sposób jego wyboru',
  [SessionConfigArea.Provider]: 'dostawca, kanał i transport',
  [SessionConfigArea.SystemPrompt]: 'prompt systemowy',
  [SessionConfigArea.Tools]: 'narzędzia dopuszczone modelowi',
  [SessionConfigArea.Permissions]: 'tryb uprawnień i reguły zgody',
  [SessionConfigArea.Mcp]: 'serwery MCP',
  [SessionConfigArea.Hooks]: 'zaczepy cyklu życia kanału',
  [SessionConfigArea.Skills]: 'umiejętności',
  [SessionConfigArea.Memory]: 'pamięć trwała',
  [SessionConfigArea.ProjectContext]: 'kontekst projektu',
  [SessionConfigArea.ConversationContext]: 'kontekst rozmowy',
  [SessionConfigArea.WorkingDirectory]: 'katalog roboczy',
  [SessionConfigArea.AdditionalDirectories]: 'katalogi dodatkowe',
  [SessionConfigArea.Environment]: 'zmienne środowiskowe kanału',
  [SessionConfigArea.InputOutput]: 'postać wejścia i wyjścia',
  [SessionConfigArea.Runtime]: 'zasięg wykonania i powłoka',
  [SessionConfigArea.SessionLifecycle]: 'granice biegu sesji',
};

/** Rejestr, z którego wzięta jest treść obszaru konfiguracji — nazwany Operatorowi zamiast pokazany jako kod. */
export const NAZWY_ZRODEL: Readonly<Record<SessionConfigSource, string>> = {
  [SessionConfigSource.Default]: 'wartość domyślna kontraktu',
  [SessionConfigSource.SettingsCatalog]: 'katalog ustawień',
  [SessionConfigSource.IdentityCatalog]: 'katalog tożsamości',
  [SessionConfigSource.AccountRegistry]: 'rejestr kont',
  [SessionConfigSource.ChannelRegistry]: 'rejestr kanałów modelu',
  [SessionConfigSource.AccessRegistry]: 'rejestr nadań dostępu',
  [SessionConfigSource.Override]: 'zapis własny konfiguracji sesji',
};

/** Powód rozejścia katalogu roboczego nazwany Operatorowi, zamiast pokazany jako kod z kontraktu rdzenia. */
export const NAZWY_ROZEJSCIA: Readonly<Record<WorkingDirectoryDegradation, string>> = {
  [WorkingDirectoryDegradation.None]: 'katalog ustawiony jest katalogiem używanym',
  [WorkingDirectoryDegradation.Missing]: 'katalog ustawiony nie istnieje',
  [WorkingDirectoryDegradation.NotPermitted]: 'proces nie ma prawa do katalogu ustawionego',
  [WorkingDirectoryDegradation.OutsideGrantedRoots]:
    'katalog leży poza korzeniami nadań dostępu okna',
  [WorkingDirectoryDegradation.DeviceUnavailable]: 'urządzenie z katalogiem jest nieosiągalne',
  [WorkingDirectoryDegradation.TemplateFailed]: 'wzorca katalogu sesji nie dało się rozwinąć',
  [WorkingDirectoryDegradation.FallbackUsed]: 'obowiązuje katalog zastępczy',
};

/** Nazwa obszaru konfiguracji sesji; obszar spoza wykazu kontraktu pokazuje wprost swój własny kod źródłowy. */
export function nazwaObszaru(obszar: string): string {
  return NAZWY_OBSZAROW[obszar as SessionConfigArea] ?? obszar;
}

/** Nazwa rejestru źródłowego wartości obszaru sesji; rejestr spoza wykazu kontraktu pokazuje własny kod. */
export function nazwaZrodla(zrodlo: string): string {
  return NAZWY_ZRODEL[zrodlo as SessionConfigSource] ?? zrodlo;
}

/** Nazwa powodu rozejścia katalogu roboczego sesji; powód spoza wykazu kontraktu pokazuje wprost własny kod. */
export function nazwaRozejscia(powod: string): string {
  return NAZWY_ROZEJSCIA[powod as WorkingDirectoryDegradation] ?? powod;
}

/**
 * Adres, pod którym rdzeń znalazł wartość obszaru; `null`, gdy pochodzenie
 * poziomu nie wskazuje — tak jest przy wartości domyślnej.
 */
export function adresPochodzenia(pochodzenie: SessionConfigOrigin): AdresUstawienia | null {
  if (pochodzenie.scope === undefined) return null;
  return {
    zasieg: pochodzenie.scope,
    bytZasiegu: pochodzenie.scopeId ?? '',
    os: pochodzenie.axis ?? ConfigAxis.Platform,
    bytOsi: pochodzenie.axisId ?? '',
  };
}

/**
 * Skąd wzięła się wartość obszaru, jednym zdaniem: rejestr źródłowy, a gdy
 * rdzeń wskazał poziom — także poziom, byt i oś. Brak zapisu nie jest dziurą:
 * zdanie nazywa wartość domyślną wprost.
 */
export function opisPochodzenia(pochodzenie: SessionConfigOrigin): string {
  const adres = adresPochodzenia(pochodzenie);
  const rejestr = nazwaZrodla(pochodzenie.source);
  return adres === null ? rejestr : `${rejestr} · ${opisAdresu(adres)}`;
}

/** Czy obszar konfiguracji jest zapisem własnym sesji, a nie wartością odziedziczoną z rejestru nadrzędnego. */
export function czyZapisWlasny(pochodzenie: SessionConfigOrigin): boolean {
  return pochodzenie.source === SessionConfigSource.Override;
}

/**
 * Zdanie o katalogu roboczym. `checkedAt` równe zeru znaczy, że nic nie było
 * sprawdzane — odczyt składa wtedy rozejście z samej konfiguracji, bez
 * dotykania dysku. Okno mówi to wprost, zamiast podawać brak sprawdzenia za
 * sprawdzenie pomyślne.
 */
export function opisKatalogu(katalog: WorkingDirectoryResolution): string {
  const uzywany = katalog.effectivePath === '' ? '(nie wskazano)' : katalog.effectivePath;
  const czlony = [`używany: ${uzywany}`];
  if (katalog.requestedPath !== undefined && katalog.requestedPath !== '') {
    czlony.push(`ustawiony: ${katalog.requestedPath}`);
  }
  czlony.push(nazwaRozejscia(katalog.reason));
  if (katalog.detail !== undefined && katalog.detail !== '') czlony.push(katalog.detail);
  if (!Number.isFinite(katalog.checkedAt) || katalog.checkedAt === 0) {
    czlony.push('dysku nie sprawdzano — rozejście złożone z samej konfiguracji');
  }
  return czlony.join(' · ');
}
