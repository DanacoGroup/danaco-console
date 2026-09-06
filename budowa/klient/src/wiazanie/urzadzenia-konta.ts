// Wykaz urządzeń w oknie Ustawień: wiersze z rdzenia i odłączanie tokenu.
import { Command } from '../../../shared/contract.ts';
import type { Device } from '../../../shared/contract.ts';
import type { Kanal } from '../protokol/kanal.ts';
import { wywolaj } from '../protokol/wywolanie.ts';
import { oglos } from './ogloszenie.ts';

const NAGLOWEK = 'Urządzenia';

export function zwiazUrzadzenia(kanal: Kanal): void {
  void wypelnijWykaz(kanal);
  document.addEventListener('click', (zdarzenie) => {
    const cel = zdarzenie.target;
    if (!(cel instanceof Element)) return;
    const przycisk = cel.closest<HTMLElement>('[data-odlacz]');
    if (przycisk === null) return;
    const wiersz = przycisk.closest<HTMLElement>('.us-urzadzenie');
    void odlacz(kanal, wiersz?.dataset.urzadzenieId ?? '');
  }, true);
}

// Wiersz prototypu służy za wzór odbitki; instalacja ma tyle urządzeń, ile
// faktycznie się na niej uwierzytelniło.
async function wypelnijWykaz(kanal: Kanal): Promise<void> {
  const tabela = document.querySelector<HTMLElement>('#sekcja-urzadzenia .us-tabela tbody');
  if (tabela === null) return;
  const wzor = tabela.querySelector<HTMLElement>('.us-urzadzenie');
  if (wzor === null) return;
  const wynik = await wywolaj(kanal, Command.DeviceList, {});
  if (!wynik.udany) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń nie podał wykazu urządzeń.', 'ostrzezenie');
    return;
  }
  const urzadzenia = wynik.wynik?.devices ?? [];
  tabela.replaceChildren();
  for (const urzadzenie of urzadzenia) {
    tabela.appendChild(wiersz(wzor, urzadzenie));
  }
}

function wiersz(wzor: HTMLElement, urzadzenie: Device): HTMLElement {
  const wezel = wzor.cloneNode(true) as HTMLElement;
  wezel.dataset.urzadzenieId = urzadzenie.deviceId;
  const komorki = wezel.querySelectorAll<HTMLElement>('td');
  if (komorki.length > 0) komorki[0].textContent = urzadzenie.name ?? urzadzenie.deviceId;
  // Token niewazny znaczy urzadzenie znane, ale odlaczone — wiersz zostaje.
  if (komorki.length > 1) komorki[1].textContent = urzadzenie.hasToken ? 'połączone' : 'odłączone';
  if (komorki.length > 2) komorki[2].textContent = chwila(urzadzenie.lastSeenAt);
  // Bieżące zostaje bez przycisku: odłączenie siebie to wylogowanie.
  if (urzadzenie.current) {
    for (const przycisk of wezel.querySelectorAll('[data-odlacz]')) przycisk.remove();
  }
  return wezel;
}

async function odlacz(kanal: Kanal, deviceId: string): Promise<void> {
  if (deviceId === '') return;
  const wynik = await wywolaj(kanal, Command.DeviceRevoke, { deviceId });
  if (!wynik.udany) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił odłączenia urządzenia.', 'ostrzezenie');
    return;
  }
  oglos(NAGLOWEK, 'Urządzenie odłączone; jego token dostępu stracił ważność.');
  void wypelnijWykaz(kanal);
}

// Brak znacznika zostawia myślnik: pusty wiersz czyta się jak błąd składania.
function chwila(znacznik: number | undefined): string {
  if (znacznik === undefined) return '—';
  return new Date(znacznik).toLocaleString('pl-PL');
}
