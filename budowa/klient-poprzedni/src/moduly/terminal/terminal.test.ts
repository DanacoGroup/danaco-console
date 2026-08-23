import { describe, expect, it } from 'vitest';
import {
  AutomationStepKind,
  Command,
  TerminalKeyType,
  TerminalShell,
  type AutomationStep,
} from '../../../../shared/contract';
import type { Kanal } from '../../protokol/kanal';
import { czytajKonfiguracjeSsh, utworzKsiazkeHostow, zapiszKonfiguracjeSsh } from './ksiazka-hostow';
import { utworzBiblioteke } from './biblioteka-skryptow';
import { utworzZrodloTerminala } from './zrodlo-terminala';

/**
 * Sprawdziany modułu Terminal — czynność, nie kształt pliku.
 *
 * Pilnowane jest to, co przy poprawce łatwo zepsuć po cichu, a co rozstrzyga
 * o tym, czy okno naprawdę dochodzi do rdzenia:
 *
 *   1. każda komenda rodziny wychodzi POD SWOJĄ NAZWĄ z kontraktu i z ładunkiem,
 *      który dostała — pomyłka w nazwie kończy się dziś odmową `not_found`
 *      dopiero na żywym rdzeniu, a tu widać ją od razu;
 *   2. odpowiedź o kształcie innym niż kontraktowy nie przechodzi jako wynik
 *      udany — okno, które przyjęłoby wykaz bez pola `hosts`, pokazałoby pustkę
 *      zamiast odmowy;
 *   3. wykaz książki hostów i biblioteki jest ZASTĘPOWANY wykazem z rdzenia,
 *      a nie scalany: wpis zdjęty w rdzeniu ma zniknąć także z ekranu;
 *   4. czytanie konfiguracji OpenSSH bierze port do wpisu, bo rdzeń podaje go
 *      programowi ssh przełącznikiem — port zgubiony po drodze kierowałby kartę
 *      na port domyślny bez słowa.
 *
 * Rdzeń jest atrapą: sprawdzian pyta o zachowanie modułu, nie serwera.
 */

interface Zapis {
  komenda: string;
  zadanie: Record<string, unknown>;
}

function atrapaKanalu(zapisy: Zapis[], odpowiedzi: Record<string, unknown> = {}): Kanal {
  return {
    wyslij(
      komenda: Command,
      zadanie: Record<string, unknown>,
      przyWyniku?: (wynik: unknown) => void,
    ): string {
      zapisy.push({ komenda: String(komenda), zadanie });
      const wynik = odpowiedzi[String(komenda)];
      przyWyniku?.(
        wynik === undefined
          ? { udany: false, blad: { code: 'not_found', message: 'atrapa nie zna tej komendy' } }
          : { udany: true, wynik },
      );
      return 'x';
    },
    naZdarzenie: () => () => undefined,
    naDowolny: () => () => undefined,
    sesja: () => ({ id: () => 'sesja-1' }) as never,
    dziennikNieznanych: () => ({}) as never,
  } as unknown as Kanal;
}

describe('źródło modułu Terminal — droga z okna do komendy', () => {
  it('wysyła każdą komendę wyposażenia pod jej nazwą z kontraktu', async () => {
    const zapisy: Zapis[] = [];
    const zrodlo = utworzZrodloTerminala(atrapaKanalu(zapisy));

    await zrodlo.zamknijKarte({ sessionId: 'term-1' });
    await zrodlo.karty({ windowId: 'okno-1' });
    await zrodlo.odczytajPlik({ sessionId: 'term-1', path: 'package.json' });
    await zrodlo.wstrzymaj({ processId: 'tproc-1' });
    await zrodlo.zapiszHosta({
      host: { id: '', name: 'web', target: 'operator@web', createdAt: 0, updatedAt: 0 },
    });
    await zrodlo.hosty({});
    await zrodlo.usunHosta({ hostId: 'thost-1' });
    await zrodlo.zapiszSkrypt({
      script: {
        id: '',
        name: 'wydanie',
        kind: 'script',
        shell: TerminalShell.Bash,
        content: 'echo x',
        version: 0,
        createdAt: 0,
        updatedAt: 0,
      },
    });
    await zrodlo.skrypty({});
    await zrodlo.usunSkrypt({ scriptId: 'tscr-1' });
    await zrodlo.sprawdzSkrypt({ content: 'echo x', shell: TerminalShell.Bash });
    await zrodlo.otworzTunel({ windowId: 'okno-1', kind: 'local' });
    await zrodlo.tunele({});
    await zrodlo.zamknijTunel({ tunnelId: 'ttun-1' });
    await zrodlo.wytworzKlucz({ name: 'klucz', keyType: TerminalKeyType.Ed25519 });
    await zrodlo.wciagnijKlucz({ name: 'klucz', path: '/tmp/id_ed25519' });
    await zrodlo.klucze();
    await zrodlo.usunKlucz({ keyId: 'tkey-1' });
    await zrodlo.zalozObserwacje({ sessionId: 'term-1', pattern: '*.go', command: 'go build ./...' });
    await zrodlo.zatrzymajObserwacje({ watchId: 'twch-1' });
    await zrodlo.obserwacje({});

    expect(zapisy.map((zapis) => zapis.komenda)).toEqual([
      Command.TerminalSessionClose,
      Command.TerminalSessionList,
      Command.TerminalFileRead,
      Command.TerminalProcessSuspend,
      Command.TerminalHostSave,
      Command.TerminalHostList,
      Command.TerminalHostRemove,
      Command.TerminalScriptSave,
      Command.TerminalScriptList,
      Command.TerminalScriptRemove,
      Command.TerminalScriptLint,
      Command.TerminalTunnelOpen,
      Command.TerminalTunnelList,
      Command.TerminalTunnelClose,
      Command.TerminalKeyGenerate,
      Command.TerminalKeyImport,
      Command.TerminalKeyList,
      Command.TerminalKeyRemove,
      Command.TerminalWatchStart,
      Command.TerminalWatchStop,
      Command.TerminalWatchList,
    ]);
  });

  it('przenosi ładunek żądania nietknięty', async () => {
    const zapisy: Zapis[] = [];
    const zrodlo = utworzZrodloTerminala(atrapaKanalu(zapisy));

    await zrodlo.zalozObserwacje({
      sessionId: 'term-7',
      pattern: 'src/**/*.ts',
      command: 'npm test',
      debounceMs: 250,
      recursive: true,
    });

    expect(zapisy[0]?.zadanie).toEqual({
      sessionId: 'term-7',
      pattern: 'src/**/*.ts',
      command: 'npm test',
      debounceMs: 250,
      recursive: true,
    });
  });

  it('nie przepuszcza odpowiedzi o kształcie spoza kontraktu jako wyniku udanego', async () => {
    const zrodlo = utworzZrodloTerminala(
      atrapaKanalu([], {
        // Odpowiedź bez pola `hosts` jest odpowiedzią cudzej komendy albo rdzenia
        // starszego niż kontrakt. Przepuszczona dałaby oknu pusty wykaz, czyli
        // zdanie „książka jest pusta" zamiast „nie udało się zapytać".
        [Command.TerminalHostList]: { total: 0 },
        [Command.TerminalKeyList]: { keys: [{ id: 'tkey-1' }], total: 1 },
      }),
    );

    const hosty = await zrodlo.hosty({});
    expect(hosty.udany).toBe(false);

    const klucze = await zrodlo.klucze();
    expect(klucze.udany).toBe(true);
    expect(klucze.wynik).toHaveLength(1);
  });

  it('oddaje prawdę o skutku usunięcia, a nie samo powodzenie wywołania', async () => {
    const zrodlo = utworzZrodloTerminala(
      atrapaKanalu([], { [Command.TerminalHostRemove]: { removed: false } }),
    );

    const wynik = await zrodlo.usunHosta({ hostId: 'thost-nieistniejacy' });

    // Wywołanie się udało, a wpisu nie było — to dwie różne rzeczy i okno musi
    // widzieć obie.
    expect(wynik.udany).toBe(true);
    expect(wynik.wynik).toBe(false);
  });
});

describe('książka hostów jako widok na dziennik rdzenia', () => {
  it('zastępuje wykaz treścią z rdzenia zamiast go scalać', () => {
    const ksiazka = utworzKsiazkeHostow();
    ksiazka.dodaj({ nazwa: 'stary', cel: 'operator@stary', grupa: '', katalog: '', notatka: '' });

    ksiazka.zastap([
      { id: 'thost-1', nazwa: 'nowy', cel: 'operator@nowy', grupa: '', katalog: '', notatka: '' },
    ]);

    expect(ksiazka.wpisy().map((wpis) => wpis.nazwa)).toEqual(['nowy']);
  });

  it('bierze port z konfiguracji OpenSSH i oddaje go z powrotem', () => {
    const wpisy = czytajKonfiguracjeSsh(
      ['Host wydanie', '  HostName wydanie.example', '  User operator', '  Port 2222'].join('\n'),
    );

    expect(wpisy).toHaveLength(1);
    expect(wpisy[0]?.cel).toBe('operator@wydanie.example');
    expect(wpisy[0]?.port).toBe(2222);
    expect(zapiszKonfiguracjeSsh(wpisy)).toContain('Port 2222');
  });

  it('pomija port niepoprawny zamiast wpisywać zmyśloną wartość', () => {
    const wpisy = czytajKonfiguracjeSsh(
      ['Host dziwny', '  HostName dziwny.example', '  Port nie-liczba'].join('\n'),
    );

    expect(wpisy[0]?.port).toBeUndefined();
  });
});

describe('plan zadania powłoki — jedna rodzina komend, druga powierzchnia', () => {
  /**
   * Zapora rozstrzygnięcia, nie sprawdzian kształtu.
   *
   * Opracowanie modułu opisuje w oknie Task & Schedule runner zadań z wyrażeniami
   * cron, a kontrakt wiąże cykliczność z automatyką. Rozjazd rozstrzygnięto na
   * korzyść kontraktu: rodzina zostaje jedna, a okno Terminala jest jej drugą
   * powierzchnią. Gdyby ktoś kiedyś wniósł rodzinę `terminal.schedule.*`, byłyby
   * dwie prawdy o jednym harmonogramie — i to jest właśnie ta pomyłka, której
   * sprawdzian ma nie przepuścić.
   */
  it('zakłada plan komendami rodziny automatyk, a nie własną rodziną terminala', async () => {
    const zapisy: Zapis[] = [];
    const zrodlo = utworzZrodloTerminala(atrapaKanalu(zapisy));

    await zrodlo.zapiszAutomatyke({
      name: 'kopia zapasowa',
      steps: [
        {
          id: 'polecenie-powloki',
          kind: AutomationStepKind.Command,
          command: Command.TerminalCommandExec,
          params: { sessionId: 'term-1', command: 'tar -czf kopia.tgz dane/' },
        },
      ],
    });
    await zrodlo.ustawHarmonogram({ workflowId: 'wf-1', cron: '0 3 * * *', enabled: true });

    expect(zapisy.map((zapis) => zapis.komenda)).toEqual([
      Command.AutomationWorkflowSave,
      Command.AutomationScheduleSet,
    ]);
    for (const zapis of zapisy) {
      expect(zapis.komenda.startsWith('terminal.schedule')).toBe(false);
    }
  });

  it('niesie polecenie powłoki w treści żądania kroku, a nie obok komendy', async () => {
    const zapisy: Zapis[] = [];
    const zrodlo = utworzZrodloTerminala(atrapaKanalu(zapisy));

    await zrodlo.zapiszAutomatyke({
      name: 'kopia zapasowa',
      steps: [
        {
          id: 'polecenie-powloki',
          kind: AutomationStepKind.Command,
          command: Command.TerminalCommandExec,
          params: { sessionId: 'term-9', command: 'tar -czf kopia.tgz dane/' },
        },
      ],
    });

    const kroki = (zapisy[0]?.zadanie as { steps?: AutomationStep[] }).steps ?? [];
    // Krok wołający komendę kontraktu przechodzi tą samą bramą uprawnień
    // i egzekutorem izolacji, co polecenie wydane ręcznie z karty. Krok niosący
    // polecenie powłoki wprost byłby drugą drogą do powłoki, bez tych sprawdzeń.
    expect(kroki[0]?.kind).toBe(AutomationStepKind.Command);
    expect(kroki[0]?.command).toBe(Command.TerminalCommandExec);
    expect(kroki[0]?.params).toEqual({ sessionId: 'term-9', command: 'tar -czf kopia.tgz dane/' });
  });
});

describe('biblioteka skryptów jako widok na dziennik rdzenia', () => {
  it('zastępuje wykaz treścią z rdzenia wraz z numerami wersji', () => {
    const biblioteka = utworzBiblioteke();
    biblioteka.zapisz({
      nazwa: 'stara',
      rodzaj: 'skrypt',
      powloka: TerminalShell.Bash,
      tresc: 'echo stara',
      tagi: '',
      wersja: 1,
      ostatnieUruchomienie: 0,
    });

    biblioteka.zastap([
      {
        id: 'tscr-1',
        nazwa: 'wydanie',
        rodzaj: 'skrypt',
        powloka: TerminalShell.Bash,
        tresc: 'echo wydanie',
        tagi: 'wydanie',
        wersja: 3,
        ostatnieUruchomienie: 0,
      },
    ]);

    expect(biblioteka.pozycje().map((pozycja) => pozycja.nazwa)).toEqual(['wydanie']);
    expect(biblioteka.znajdz('wydanie')?.wersja).toBe(3);
  });
});
