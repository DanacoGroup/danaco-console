import { describe, expect, it } from 'vitest';

import {
  AutomationAlertTrigger,
  Command,
  QueueAction,
  QueueScope,
} from '../../../../shared/contract';
import type { Kanal, Wynik } from '../../protokol/kanal';
import { utworzZrodloAutomations } from './zrodlo-automations';

/** Droga z okna do rdzenia dla trzech rodzin modułu Automations, mierzona wprost z kontraktu. */

/** Kanał próbny zapamiętuje komendy i oddaje odpowiedź pustą, mierząc warstwę kliencką, nie zachowanie rdzenia. */
function kanalProbny(): { kanal: Kanal; wyslane: string[] } {
  const wyslane: string[] = [];
  const kanal = {
    wyslij(komenda: string, _zadanie: unknown, przyWyniku?: (wynik: Wynik<unknown>) => void) {
      wyslane.push(komenda);
      przyWyniku?.({ udany: true, wynik: {} });
      return `zadanie-${wyslane.length}`;
    },
    naZdarzenie: () => () => undefined,
    naDowolny: () => () => undefined,
    sesja: () => ({}) as ReturnType<Kanal['sesja']>,
    dziennikNieznanych: () => ({}) as ReturnType<Kanal['dziennikNieznanych']>,
  } as unknown as Kanal;
  return { kanal, wyslane };
}

/** Woła każdą czynność źródła po kolei raz, dając pełny przelot trzech rodzin komend modułu Automations. */
async function przelotAutomations(kanal: Kanal): Promise<void> {
  const zrodlo = utworzZrodloAutomations(kanal);
  const automatyka = 'automat-przykladowa';
  const kolejka = '1';
  const zlecenie = 'zlec-1';
  const przebieg = 'przebieg-1';
  const harmonogram = 'harm-1';

  // Rdzeń modułu — komendy, którymi okna pracują bez otwierania panelu.
  await zrodlo.zapiszAutomatyke({ name: 'Raport tygodniowy' });
  await zrodlo.automatyki({});
  await zrodlo.ustawHarmonogram({ workflowId: automatyka, cron: '0 7 * * 1' });
  await zrodlo.ustawZaleznosci({ workflowId: automatyka });
  await zrodlo.przebiegi({ workflowId: automatyka });
  await zrodlo.zalozKolejke({ sessionId: 'sesja-1' });
  await zrodlo.dzialanieKolejki({ queueId: kolejka, action: QueueAction.Start });
  // Akcja kolejki posuwa kolejkę samą, na przykład kolejkę przeglądu bez automatyki.
  await zrodlo.dzialanieNaKolejce({ queueId: kolejka, action: QueueAction.Pause });
  await zrodlo.wykazKolejek({});
  await zrodlo.powiazKolejke({ queueId: kolejka, workflowId: automatyka });
  await zrodlo.odczytajHarmonogramy({});

  // Workflow Builder — wersje, etykiety, zmienne, kanwa, szablony, publikacja.
  await zrodlo.symulujPrzeplyw({ workflowId: automatyka });
  await zrodlo.wersjeAutomatyki({ workflowId: automatyka });
  await zrodlo.przywrocWersje({ workflowId: automatyka, version: 1 });
  await zrodlo.porownajWersje({ workflowId: automatyka, fromVersion: 1, toVersion: 2 });
  await zrodlo.ustawEtykiety({ workflowId: automatyka, tags: [] });
  await zrodlo.ustawZmienne({ workflowId: automatyka, variables: [] });
  await zrodlo.ustawNotatkeKroku({ workflowId: automatyka, stepId: 'krok-1', note: '' });
  await zrodlo.ustawUkladKanwy({ workflowId: automatyka, positions: [] });
  await zrodlo.zapiszSzablon({ workflowId: automatyka, name: 'Wzorzec' });
  await zrodlo.szablony({});
  await zrodlo.zastosujSzablon({ templateId: 'szablon-1', name: 'Z szablonu' });
  await zrodlo.opublikujAutomatyke({ workflowId: automatyka });
  await zrodlo.udostepnijAutomatyke({ workflowId: automatyka, shared: true });

  // Scheduler — okna wykonania, uruchomienie wsteczne, nadzór, webhook.
  await zrodlo.ustawOknaWykonania({ scheduleId: harmonogram, windows: [] });
  await zrodlo.uruchomWstecznie({ scheduleId: harmonogram, fromAt: 0, toAt: 1 });
  await zrodlo.historiaWyzwolen({});
  await zrodlo.ustawNadzorUruchomien({ scheduleId: harmonogram, toleranceSeconds: 0 });
  await zrodlo.adresWebhooka({ workflowId: automatyka });

  // Queue Manager — jedenaście działań na zleceniu, wykaz, polityka, martwe, głębokość.
  await zrodlo.dodajZlecenie({ queueId: kolejka, payload: {} });
  await zrodlo.zdejmijZlecenie({ queueId: kolejka, itemId: zlecenie });
  await zrodlo.odlozZlecenie({ queueId: kolejka, itemId: zlecenie, delaySeconds: 60 });
  await zrodlo.podzielZlecenie({ queueId: kolejka, itemId: zlecenie, payloads: [] });
  await zrodlo.scalZlecenia({ queueId: kolejka, itemIds: [zlecenie] });
  await zrodlo.skierujZlecenie({ queueId: kolejka, itemId: zlecenie });
  await zrodlo.rozgalezZlecenie({ queueId: kolejka, itemId: zlecenie, branches: [] });
  await zrodlo.uwarunkujZlecenie({ queueId: kolejka, itemId: zlecenie, condition: '' });
  await zrodlo.zleceniaKolejki({ queueId: kolejka });
  await zrodlo.ustawPolitykeKolejki({
    queueId: kolejka, policy: { scope: QueueScope.Local },
  });
  await zrodlo.zadaniaMartwe({});
  await zrodlo.glebokoscKolejki({ fromAt: 0, bucketSeconds: 3600 });

  // Execution Monitor — log, kroki, punkty, ładunki, alarmy, budżety, skarbiec, audyt.
  await zrodlo.dziennikPrzebiegu({ executionId: przebieg });
  await zrodlo.krokiPrzebiegu({ executionId: przebieg });
  await zrodlo.punktyWznowienia({ executionId: przebieg });
  await zrodlo.wznowPrzebieg({ executionId: przebieg });
  await zrodlo.podgladLadunku({ executionId: przebieg, stepId: 'krok-1' });
  await zrodlo.odtworzPrzebieg({ executionId: przebieg });
  await zrodlo.ustawReguleAlarmowania({
    workflowId: automatyka, trigger: AutomationAlertTrigger.Failure, channels: [],
  });
  await zrodlo.regulyAlarmowania({});
  await zrodlo.ustawBudzetyPrzebiegu({ workflowId: automatyka });
  await zrodlo.zapiszPoswiadczenie({ name: 'klucz', value: 'wartosc' });
  await zrodlo.poswiadczenia({});
  await zrodlo.usunPoswiadczenie({ secretRef: 'sejf:klucz' });
  await zrodlo.odczytajAudyt({});
}

describe('droga z okna do rdzenia — moduł Automations', () => {
  it('ma drogę do każdej komendy rodzin automation.*, queue.* i schedule.*', async () => {
    const { kanal, wyslane } = kanalProbny();
    await przelotAutomations(kanal);

    const rodziny = ['automation.', 'queue.', 'schedule.'];
    const zRodzin = Object.values(Command).filter((nazwa) =>
      rodziny.some((rodzina) => nazwa.startsWith(rodzina)),
    );
    const bezDrogi = zRodzin.filter((nazwa) => !wyslane.includes(nazwa));

    expect(
      bezDrogi,
      `komendy modułu Automations bez drogi z okna: ${bezDrogi.join(', ')}`,
    ).toEqual([]);
  });

  it('nie wysyła pól opcjonalnych, których Operator nie wskazał', async () => {
    // Pole nieobecne w żądaniu znaczy co innego niż pole o wartości pustej lub zerowej.
    const wyslaneZadania: Record<string, unknown> = {};
    const kanal = {
      wyslij(komenda: string, zadanie: unknown, przyWyniku?: (wynik: Wynik<unknown>) => void) {
        wyslaneZadania[komenda] = zadanie;
        przyWyniku?.({ udany: true, wynik: {} });
        return 'zadanie-1';
      },
      naZdarzenie: () => () => undefined,
      naDowolny: () => () => undefined,
      sesja: () => ({}) as ReturnType<Kanal['sesja']>,
      dziennikNieznanych: () => ({}) as ReturnType<Kanal['dziennikNieznanych']>,
    } as unknown as Kanal;

    const zrodlo = utworzZrodloAutomations(kanal);
    await zrodlo.dodajZlecenie({ queueId: '1', payload: {} });

    expect(wyslaneZadania[Command.QueueItemEnqueue]).not.toHaveProperty('idempotencyKey');
    expect(wyslaneZadania[Command.QueueItemEnqueue]).not.toHaveProperty('priority');
    expect(wyslaneZadania[Command.QueueItemEnqueue]).not.toHaveProperty('scheduledAt');
  });
});
