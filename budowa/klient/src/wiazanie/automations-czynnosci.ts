// Czynności modułu Automations: uruchomienie automatyki, sterowanie kolejką
// i włączanie automatyki. Okno wystawiało dotąd same wykazy rdzenia.
import { Command, QueueAction } from '../../../shared/contract.ts';
import type { AutomationWorkflow, Queue } from '../../../shared/contract.ts';
import type { Kanal } from '../protokol/kanal.ts';
import { wywolaj } from '../protokol/wywolanie.ts';
import { oglos } from './ogloszenie.ts';

const NAGLOWEK = 'Automatyki';

// Nasłuch stoi na korzeniu okna, więc przetrwa odświeżenie wykazu.
export function zwiazCzynnosciAutomatyk(kanal: Kanal, korzen: Element, odswiez: () => void): void {
  korzen.addEventListener('click', (zdarzenie) => {
    const cel = zdarzenie.target;
    if (!(cel instanceof Element)) return;
    const przycisk = cel.closest<HTMLElement>('[data-czynnosc-automatyki]');
    if (przycisk === null) return;
    zdarzenie.stopPropagation();
    const kod = przycisk.closest<HTMLElement>('[data-automatyka]')?.dataset.automatyka ?? '';
    void wykonaj(kanal, przycisk.dataset.czynnoscAutomatyki ?? '', kod, odswiez);
  }, true);
}

async function wykonaj(kanal: Kanal, czynnosc: string, kod: string, odswiez: () => void): Promise<void> {
  if (kod === '') return;
  if (czynnosc === 'uruchom') return uruchom(kanal, kod, odswiez);
  if (czynnosc === 'wlacz' || czynnosc === 'wylacz') {
    return przestaw(kanal, kod, czynnosc === 'wlacz', odswiez);
  }
  if (czynnosc === 'zatrzymaj') return zatrzymaj(kanal, kod, odswiez);
}

// Kolejka jest jedyną drogą wykonania: rdzeń nie zna komendy „wykonaj teraz".
async function uruchom(kanal: Kanal, workflowId: string, odswiez: () => void): Promise<void> {
  const kolejka = await kolejkaAutomatyk(kanal);
  if (kolejka === '') {
    oglos(NAGLOWEK, 'Rdzeń nie prowadzi żadnej kolejki, do której można wstawić bieg.', 'ostrzezenie');
    return;
  }
  const wynik = await wywolaj(kanal, Command.AutomationQueueAction, {
    queueId: kolejka,
    action: QueueAction.Enqueue,
    workflowId,
  });
  if (!wynik.udany) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił uruchomienia automatyki.', 'ostrzezenie');
    return;
  }
  oglos(NAGLOWEK, 'Automatyka stoi w kolejce wykonania.');
  odswiez();
}

async function zatrzymaj(kanal: Kanal, workflowId: string, odswiez: () => void): Promise<void> {
  const kolejka = await kolejkaAutomatyk(kanal);
  if (kolejka === '') return;
  const wynik = await wywolaj(kanal, Command.AutomationQueueAction, {
    queueId: kolejka,
    action: QueueAction.Stop,
    workflowId,
  });
  if (!wynik.udany) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił zatrzymania biegu.', 'ostrzezenie');
    return;
  }
  oglos(NAGLOWEK, 'Bieg zatrzymany.');
  odswiez();
}

// Kontrakt nie zna komendy przełączenia; zapis idzie całą automatyką, a nazwa
// jest w nim wymagana i bierze się z wykazu.
async function przestaw(kanal: Kanal, workflowId: string, enabled: boolean, odswiez: () => void): Promise<void> {
  const nazwa = await nazwaAutomatyki(kanal, workflowId);
  if (nazwa === '') {
    oglos(NAGLOWEK, 'Rdzeń nie podał nazwy automatyki; zapis bez niej nie przejdzie.', 'ostrzezenie');
    return;
  }
  const wynik = await wywolaj(kanal, Command.AutomationWorkflowSave, { workflowId, name: nazwa, enabled });
  if (!wynik.udany) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił zapisu automatyki.', 'ostrzezenie');
    return;
  }
  oglos(NAGLOWEK, enabled ? 'Automatyka włączona.' : 'Automatyka wyłączona.');
  odswiez();
}

async function kolejkaAutomatyk(kanal: Kanal): Promise<string> {
  const wynik = await wywolaj(kanal, Command.QueueList, {});
  if (!wynik.udany) return '';
  const kolejki: Queue[] = wynik.wynik?.queues ?? [];
  return kolejki[0]?.id ?? '';
}

async function nazwaAutomatyki(kanal: Kanal, workflowId: string): Promise<string> {
  const wynik = await wywolaj(kanal, Command.AutomationWorkflowList, {});
  if (!wynik.udany) return '';
  const wykaz: AutomationWorkflow[] = wynik.wynik?.workflows ?? [];
  return wykaz.find((w) => w.id === workflowId)?.name ?? '';
}
