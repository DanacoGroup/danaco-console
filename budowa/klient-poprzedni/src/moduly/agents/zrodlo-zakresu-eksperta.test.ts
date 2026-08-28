import { describe, expect, it } from 'vitest';

import {
  AgentPermissionGroup,
  Command,
  IsolationTechnicalScope,
} from '../../../../shared/contract';
import type { Kanal, Wynik } from '../../protokol/kanal';
import { utworzZrodloZakresuEksperta } from './zrodlo-zakresu-eksperta';

/**
 * Kanał próbny sprawdzianów zakresu działania eksperta: zapamiętuje wysłane
 * komendy wraz z żądaniami i oddaje odpowiedź pustą albo tę, którą wskazano
 * przy zakładaniu kanału.
 */
function kanalProbny(odpowiedzi: Record<string, unknown> = {}): {
  kanal: Kanal;
  wyslane: string[];
  zadania: Record<string, Record<string, unknown>>;
} {
  const wyslane: string[] = [];
  const zadania: Record<string, Record<string, unknown>> = {};
  const kanal = {
    wyslij(komenda: string, zadanie: unknown, przyWyniku?: (wynik: Wynik<unknown>) => void) {
      wyslane.push(komenda);
      zadania[komenda] = zadanie as Record<string, unknown>;
      przyWyniku?.({ udany: true, wynik: odpowiedzi[komenda] ?? {} });
      return `zadanie-${wyslane.length}`;
    },
    naZdarzenie: () => () => undefined,
    naDowolny: () => () => undefined,
    sesja: () => ({}) as ReturnType<Kanal['sesja']>,
    dziennikNieznanych: () => ({}) as ReturnType<Kanal['dziennikNieznanych']>,
  } as unknown as Kanal;
  return { kanal, wyslane, zadania };
}

/**
 * Wywołuje każdą czynność źródła zakresu jeden raz, składając pełny przelot
 * obszaru; wykaz komend wysłanych powstaje z tego przelotu, a nie z pamięci.
 */
async function przelotZakresu(kanal: Kanal): Promise<void> {
  const zrodlo = utworzZrodloZakresuEksperta(kanal);

  await zrodlo.odlaczUmiejetnosc('ekspert-1', 'korekta');
  await zrodlo.konektory('ekspert-1');
  await zrodlo.odlaczKonektor('ekspert-1', 'konektor-1');
  await zrodlo.skonfigurujKonektor({
    idEksperta: 'ekspert-1',
    idKonektora: 'konektor-1',
    czynny: false,
  });
  await zrodlo.wersja('ekspert-1', 'wersja-1');
  await zrodlo.przypisania();
  await zrodlo.ustawModuly('ekspert-1', ['terminal']);
  await zrodlo.izolacja('ekspert-1');
  await zrodlo.zapiszIzolacje('ekspert-1', [
    { scope: IsolationTechnicalScope.NetworkAccess, isolated: true },
  ]);
  await zrodlo.polityka('ekspert-1');
  await zrodlo.ustawPodagentow('ekspert-1', true, 5);
  await zrodlo.zdejmijUprawnienie({ idEksperta: 'ekspert-1' });
  await zrodlo.zakresyNarzedzi();
  await zrodlo.zapiszZakresNarzedzia({
    idProfilu: 'profil-1',
    nazwaPozycji: 'danaco:session.list',
    dostepna: false,
  });
}

describe('źródło zakresu działania eksperta', () => {
  it('ma drogę z okna do każdej komendy zakresu', async () => {
    const { kanal, wyslane } = kanalProbny();
    await przelotZakresu(kanal);

    const oczekiwane = [
      Command.AgentSkillRemove,
      Command.AgentConnectorList,
      Command.AgentConnectorRemove,
      Command.AgentConnectorConfigure,
      Command.AgentVersionGet,
      Command.AgentAssignmentList,
      Command.AgentModulesSet,
      Command.AgentIsolationGet,
      Command.AgentIsolationSet,
      Command.AgentPolicyGet,
      Command.AgentSubagentSet,
      Command.AgentPermissionRemove,
      Command.ToolsScopeList,
      Command.ToolsScopeSet,
    ];
    expect([...wyslane].sort()).toEqual([...oczekiwane].sort());
  });

  it('wysyła wykaz modułów także pusty — pusty znaczy brak ograniczenia', async () => {
    const { kanal, zadania } = kanalProbny();
    const zrodlo = utworzZrodloZakresuEksperta(kanal);

    await zrodlo.ustawModuly('ekspert-1', []);
    expect(zadania[Command.AgentModulesSet]?.['moduleCodes']).toEqual([]);
  });

  it('zdejmuje wpisy wszystkich grup, gdy grupy nie wskazano', async () => {
    const { kanal, zadania } = kanalProbny();
    const zrodlo = utworzZrodloZakresuEksperta(kanal);

    await zrodlo.zdejmijUprawnienie({ idEksperta: 'ekspert-1' });
    const bezGrupy = zadania[Command.AgentPermissionRemove];
    expect(bezGrupy).toEqual({ agentId: 'ekspert-1' });

    await zrodlo.zdejmijUprawnienie({
      idEksperta: 'ekspert-1',
      grupa: AgentPermissionGroup.Modules,
      zakres: 'studio',
    });
    expect(zadania[Command.AgentPermissionRemove]).toEqual({
      agentId: 'ekspert-1',
      group: AgentPermissionGroup.Modules,
      scope: 'studio',
    });
  });

  it('nie wysyła granicy podagentów przy wyłączeniu — zero zapisuje rdzeń', async () => {
    const { kanal, zadania } = kanalProbny();
    const zrodlo = utworzZrodloZakresuEksperta(kanal);

    await zrodlo.ustawPodagentow('ekspert-1', false);
    expect(zadania[Command.AgentSubagentSet]).toEqual({ agentId: 'ekspert-1', enabled: false });
  });

  it('odmawia wywołania przy niepoprawnym JSON konfiguracji zamiast rzucać wyjątkiem', async () => {
    const { kanal, wyslane } = kanalProbny();
    const zrodlo = utworzZrodloZakresuEksperta(kanal);

    const wynik = await zrodlo.skonfigurujKonektor({
      idEksperta: 'ekspert-1',
      idKonektora: 'konektor-1',
      konfiguracja: '{ to nie jest json',
    });
    expect(wynik.udany).toBe(false);
    expect(wyslane).toEqual([]);
  });
});
