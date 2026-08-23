import {
  ChangeKind,
  Command,
  ConfigScope,
  EventType,
  ExecutionEnv,
  PermissionMode,
  type ConfigEntry,
  type Envelope,
  type Window,
} from '../../../shared/contract';
import type { Kanal } from '../protokol/kanal';

/**
 * Przyrząd pomiarowy sprawdzianów — nie jest częścią produktu i nie ma ani
 * jednego wołacza poza plikami `*.test.ts`.
 *
 * Odwzorowuje rdzeń, żeby sprawdzian mierzył zachowanie klienta wobec tego, co
 * rdzeń robi naprawdę. Każda gałąź poniżej ma wskazane miejsce w kodzie Go,
 * które odwzorowuje.
 *
 * Odwzorowane zachowania:
 *
 *  · rozgłoszenie idzie także do nadawcy — `transport/rozgloszenie.go`:
 *    `Rozglos` woła `rozglosPoza(konto, "", k)` z pustym identyfikatorem
 *    pomijanym. Na tym stoi zbieżność równoległych egzemplarzy stanu okna.
 *
 *  · `config.set` → odpowiedź `{entry}` i zdarzenie `config.changed`
 *    o rodzaju `updated` (`core/handlers_config.go`).
 *
 *  · `config.reset` → odpowiedź `{entries}` i zdarzenie `config.changed`
 *    o rodzaju `deleted`, niosące wpis ze starą wartością
 *    (`core/handlers_config.go`; `core/adapter_ustawienia.go`, funkcja
 *    `Przywroc` — usuwa wiersz i oddaje to, co usunęła).
 *
 *  · `config.get` z poziomem → surowe wpisy tego poziomu, niezależnie od osi
 *    (`core/adapter_ustawienia.go`, gałąź `ListaPoziomu`).
 *
 *  · `config.get` bez poziomu → polityka efektywna: po jednym zwycięskim wpisie
 *    na klucz, z polem `scope` niosącym poziom, na którym wartość znaleziono
 *    (tamże, gałąź `z.Scope == nil` → `PolitykaEfektywna(kontekstZasiegu(ScopeId))`
 *    → `konfig/odwzorowanie_kontraktu.go`, `WpisKontraktu`). Kontekst tej gałęzi
 *    ma wypełnione wyłącznie okno (`kontekstZasiegu`), więc rozstrzyganie
 *    obejmuje poziom okna i poziom globalny osi platformy. Karty sesji ani osi
 *    modelu ta droga nie widzi.
 *
 *  · `window.update` → odpowiedź `{window}` i zdarzenie `window.changed`.
 */

/** Okno wyjściowe sprawdzianów; każde zbudowanie daje własny egzemplarz. */
export function oknoProbne(nadpisania: Partial<Window> = {}): Window {
  return {
    id: 'okno-1',
    sessionId: 'sesja-1',
    moduleId: 'studio',
    windowRole: 'standalone',
    modelChannelId: 'kan-1',
    permissionMode: PermissionMode.Manual,
    workingDirs: [],
    executionEnv: ExecutionEnv.Local,
    ...nadpisania,
  } as unknown as Window;
}

/** Adres wpisu w bazie atrapy: poziom, byt poziomu, oś, byt osi, klucz. */
function adres(wpis: ConfigEntry): string {
  return [wpis.scope, wpis.scopeId ?? '', wpis.axis ?? '', wpis.axisId ?? '', wpis.key].join(
    '|',
  );
}

/** Rdzeń odwzorowany na potrzeby pomiaru. */
export interface AtrapaRdzenia {
  kanal: Kanal;
  /** Okno tak, jak widzi je „rdzeń". */
  okno(): Window;
  /** Zapisuje wpis wprost do bazy, z pominięciem komendy — ustawienie sceny. */
  wstaw(wpis: ConfigEntry): void;
  /** Ile razy padła komenda o wskazanej nazwie. */
  liczba(komenda: string): number;
  /** Wszystkie wysłane żądania w kolejności, wraz z treścią. */
  wyslane(): readonly { komenda: string; zadanie: unknown }[];
}

export function utworzAtrapeRdzenia(poczatkowe: Window = oknoProbne()): AtrapaRdzenia {
  let okno = poczatkowe;
  const baza = new Map<string, ConfigEntry>();
  const sluchacze: ((koperta: Envelope) => void)[] = [];
  const wyslane: { komenda: string; zadanie: unknown }[] = [];

  function rozglos(typ: string, tresc: unknown): void {
    const koperta = { id: 'zdarzenie', type: typ, payload: tresc } as unknown as Envelope;
    for (const sluchacz of [...sluchacze]) sluchacz(koperta);
  }

  /**
   * Polityka efektywna: zwycięża poziom najwęższy. Kontekst gałęzi bez poziomu
   * zna wyłącznie okno, więc wykaz adresów to okno i poziom globalny.
   */
  function politykaEfektywna(idOkna: string): ConfigEntry[] {
    const kolejnosc = [
      { scope: ConfigScope.Window, scopeId: idOkna },
      { scope: ConfigScope.Global, scopeId: '' },
    ];
    const zwyciezcy = new Map<string, ConfigEntry>();
    for (const poziom of kolejnosc) {
      for (const wpis of baza.values()) {
        if (wpis.scope !== poziom.scope) continue;
        if ((wpis.scopeId ?? '') !== poziom.scopeId) continue;
        // Oś platformy: wpis osi węższej nie wchodzi do tego rozstrzygania.
        if ((wpis.axisId ?? '') !== '') continue;
        if (zwyciezcy.has(wpis.key)) continue;
        zwyciezcy.set(wpis.key, wpis);
      }
    }
    return [...zwyciezcy.values()];
  }

  const kanal = {
    wyslij(komenda: string, zadanie: Record<string, unknown>, przyWyniku?: (w: unknown) => void) {
      wyslane.push({ komenda, zadanie });

      if (komenda === Command.WindowUpdate) {
        okno = { ...okno, ...zadanie } as Window;
        przyWyniku?.({ udany: true, wynik: { window: okno } });
        rozglos(EventType.WindowChanged, { change: ChangeKind.Updated, window: okno });
        return 'x';
      }

      if (komenda === Command.ConfigSet) {
        const wpis = wpisZadania(zadanie);
        baza.set(adres(wpis), wpis);
        przyWyniku?.({ udany: true, wynik: { entry: wpis } });
        rozglos(EventType.ConfigChanged, { change: ChangeKind.Updated, entry: wpis });
        return 'x';
      }

      if (komenda === Command.ConfigReset) {
        const wzorzec = wpisZadania({ ...zadanie, value: null });
        const klucz = adres(wzorzec);
        const usuwany = baza.get(klucz);
        baza.delete(klucz);
        const wpisy = usuwany === undefined ? [] : [usuwany];
        przyWyniku?.({ udany: true, wynik: { entries: wpisy } });
        for (const stary of wpisy) {
          rozglos(EventType.ConfigChanged, { change: ChangeKind.Deleted, entry: stary });
        }
        return 'x';
      }

      if (komenda === Command.ConfigGet) {
        const wpisy =
          zadanie['scope'] === undefined
            ? politykaEfektywna(String(zadanie['scopeId'] ?? ''))
            : [...baza.values()].filter(
                (wpis) =>
                  wpis.scope === zadanie['scope'] &&
                  (wpis.scopeId ?? '') === String(zadanie['scopeId'] ?? ''),
              );
        przyWyniku?.({ udany: true, wynik: { entries: wpisy } });
        return 'x';
      }

      przyWyniku?.({ udany: true, wynik: {} });
      return 'x';
    },

    naZdarzenie(zdarzenie: string, sluchacz: (tresc: unknown, koperta: Envelope) => void) {
      const opakowany = (koperta: Envelope): void => {
        if (koperta.type !== zdarzenie) return;
        sluchacz((koperta as unknown as { payload: unknown }).payload, koperta);
      };
      sluchacze.push(opakowany);
      return () => {
        const miejsce = sluchacze.indexOf(opakowany);
        if (miejsce >= 0) sluchacze.splice(miejsce, 1);
      };
    },

    naDowolny(sluchacz: (koperta: Envelope) => void) {
      sluchacze.push(sluchacz);
      return () => undefined;
    },

    sesja: () => ({ id: () => 'sesja-1' }),
    dziennikNieznanych: () => ({}),
  } as unknown as Kanal;

  return {
    kanal,
    okno: () => okno,
    wstaw(wpis) {
      baza.set(adres(wpis), wpis);
    },
    liczba: (komenda) => wyslane.filter((pozycja) => pozycja.komenda === komenda).length,
    wyslane: () => wyslane,
  };
}

/** Wpis konfiguracji złożony z treści komendy; człony puste są pomijane. */
function wpisZadania(zadanie: Record<string, unknown>): ConfigEntry {
  const wpis: ConfigEntry = {
    key: String(zadanie['key'] ?? ''),
    value: zadanie['value'],
    scope: zadanie['scope'] as ConfigScope,
  };
  if (zadanie['scopeId'] !== undefined) wpis.scopeId = String(zadanie['scopeId']);
  if (zadanie['axis'] !== undefined) wpis.axis = zadanie['axis'] as ConfigEntry['axis'];
  if (zadanie['axisId'] !== undefined) wpis.axisId = String(zadanie['axisId']);
  return wpis;
}
