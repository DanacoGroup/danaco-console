// Obserwacje stron: założenie monitora zmian pod adresem z paska i zdjęcie
// monitora stojącego. Sprawdzenie monitora stoi już w szynie wytworów.

import { Command } from '../../../shared/contract.ts';
import type { Kanal } from '../protokol/kanal.ts';
import { wywolaj } from '../protokol/wywolanie.ts';
import {
  adresWymagany,
  dolozPrzycisk,
  odmowa,
  polePasa,
  potwierdzone,
  powiedz,
  zalozPas,
} from './browser-wspolne.ts';

export function zwiazObserwacjeStron(
  kanal: Kanal,
  korzen: Element,
  idOkna: string,
  przy: AddEventListenerOptions,
  odswiezObserwacje: () => void,
): void {
  const poleOkresu = postawPas(korzen);

  korzen.addEventListener('click', (zdarzenie) => {
    const cel = zdarzenie.target;
    if (!(cel instanceof Element)) return;

    if (cel.closest('[data-obserwacja-zaloz]') !== null) {
      void zalozObserwacje(kanal, korzen, idOkna, poleOkresu, odswiezObserwacje);
      return;
    }
    const zdjecie = cel.closest<HTMLElement>('[data-obserwacja-zdejmij]');
    if (zdjecie !== null) {
      void zdejmijObserwacje(kanal, zdjecie, odswiezObserwacje);
    }
  }, przy);
}

function postawPas(korzen: Element): HTMLInputElement | null {
  const panel = korzen.querySelector('#panel-browser');
  const plotno = panel?.querySelector('.brw-plotno') ?? null;
  if (panel === null || plotno === null) return null;
  const pas = zalozPas('Obserwacje strony');
  const pole = polePasa(pas, 'Co ile sekund sprawdzać');
  dolozPrzycisk(pas, 'Załóż monitor', 'obserwacjaZaloz');
  panel.insertBefore(pas, plotno);
  return pole;
}

async function zalozObserwacje(
  kanal: Kanal,
  korzen: Element,
  idOkna: string,
  poleOkresu: HTMLInputElement | null,
  odswiez: () => void,
): Promise<void> {
  const url = adresWymagany(korzen);
  if (url === '') return;
  const okres = Number.parseInt((poleOkresu?.value ?? '').trim(), 10);
  const wynik = await wywolaj(kanal, Command.BrowserMonitorAdd, {
    windowId: idOkna,
    url,
    intervalSeconds: Number.isFinite(okres) && okres > 0 ? okres : undefined,
  });
  if (!wynik.udany) {
    odmowa(wynik.blad, 'Rdzeń odmówił założenia monitora strony.');
    return;
  }
  const monitor = wynik.wynik?.monitor;
  powiedz(`Monitor strony stoi; sprawdzenie co ${monitor?.intervalSeconds ?? '?'} s.`);
  odswiez();
}

/* Zdjęcie monitora jest nieodwracalne — historia porównań schodzi razem z nim,
   więc czynność wymaga drugiego naciśnięcia. */
async function zdejmijObserwacje(
  kanal: Kanal,
  przycisk: HTMLElement,
  odswiez: () => void,
): Promise<void> {
  const monitorId = przycisk.dataset.obserwacjaZdejmij ?? '';
  if (monitorId === '') return;
  if (!potwierdzone(przycisk, 'Zdjęcie monitora kasuje jego historię. Naciśnij drugi raz.')) return;
  const wynik = await wywolaj(kanal, Command.BrowserMonitorRemove, { monitorId });
  if (!wynik.udany) {
    odmowa(wynik.blad, 'Rdzeń odmówił zdjęcia monitora.');
    return;
  }
  powiedz(wynik.wynik?.removed === true
    ? 'Monitor zdjęty.'
    : 'Rdzeń przyjął żądanie, ale monitora nie zdjął.');
  odswiez();
}
