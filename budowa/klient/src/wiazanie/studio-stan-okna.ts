// Stan okna Studia zasilany zdarzeniami warstwy wspólnej; znacznik należy do
// Właściciela, a plakietkę kolejki powiela z plakietki stanu.

import {
  ChangeKind,
  Command,
  EventType,
  PermissionMode,
  ProgressStatus,
  QueueStatus,
} from '../../../shared/contract.ts';
import type { Odsubskrybuj } from '../polaczenie/magistrala-zdarzen.ts';
import { zglosUchwyt } from '../polaczenie/rozdzielacz-zdarzen.ts';
import type { Kanal } from '../protokol/kanal.ts';
import { wywolaj } from '../protokol/wywolanie.ts';

interface WezlyStanu {
  stan: HTMLElement;
  plakietki: HTMLElement[];
  tetno: Element | null;
  tryb: HTMLElement | null;
}

const NAZWY_BIEGU: Readonly<Record<ProgressStatus, string>> = {
  [ProgressStatus.Pending]: 'oczekuje',
  [ProgressStatus.Running]: 'pracuje',
  [ProgressStatus.Paused]: 'wstrzymany',
  [ProgressStatus.Stopped]: 'zatrzymany',
  [ProgressStatus.Done]: 'zakończony',
  [ProgressStatus.Failed]: 'błędny',
};

const NAZWY_KOLEJKI: Readonly<Record<QueueStatus, string>> = {
  [QueueStatus.Idle]: 'bezczynna',
  [QueueStatus.Running]: 'pracuje',
  [QueueStatus.Paused]: 'wstrzymana',
  [QueueStatus.Stopped]: 'zatrzymana',
  [QueueStatus.Done]: 'wyczerpana',
};

const NAZWY_TRYBU: Readonly<Record<PermissionMode, string>> = {
  [PermissionMode.Manual]: 'Ręczny',
  [PermissionMode.AcceptEdits]: 'Przyjmij zmiany',
  [PermissionMode.Plan]: 'Plan',
  [PermissionMode.Auto]: 'Auto',
  [PermissionMode.DontAsk]: 'Bez zapytań',
  [PermissionMode.BypassPermissions]: 'Z pominięciem zgód',
};

// Wzór tętna zdjęty raz: plakietka opróżniona ze stanu przykładowego go nie odda.
let wzorTetna: Element | null = null;

export function zdejmijTrescPrzykladowaStanu(korzen: ParentNode): boolean {
  const wezly = zbierzWezly(korzen);
  if (wezly === null) return false;
  opiszBieg(wezly, null, '');
  opiszTryb(wezly, null);
  return true;
}

export function zwiazStanOkna(kanal: Kanal, idOkna: string, korzen: ParentNode): Odsubskrybuj | null {
  const znalezione = zbierzWezly(korzen);
  if (znalezione === null) return null;
  const wezly: WezlyStanu = znalezione;
  const odlaczenia: Odsubskrybuj[] = [];

  const odswiez = async (): Promise<void> => {
    const wynik = await wywolaj(kanal, Command.WindowStateGet, { windowId: idOkna });
    if (!wynik.udany || wynik.wynik === undefined) return;
    opiszBieg(wezly, wynik.wynik.processStatus, '');
    opiszTryb(wezly, wynik.wynik.window.permissionMode);
  };

  // Nazwę etapu niesie wyłącznie telemetria tury; stan okna niesie sam stan.
  let krokTury = '';

  odlaczenia.push(
    zglosUchwyt(EventType.ProgressChanged, (tresc) => {
      if (tresc.windowId !== idOkna) return;
      krokTury = tresc.status === ProgressStatus.Running ? (tresc.stepLabel ?? '') : '';
      opiszBieg(wezly, tresc.status, krokTury);
    }),
    zglosUchwyt(EventType.WindowStateChanged, (tresc) => {
      if (tresc.windowId !== idOkna) return;
      naniesStanOkna(wezly, tresc.state, krokTury);
    }),
    // Zasięg nastawy rozstrzyga rdzeń, więc klient pyta o stan efektywny okna.
    zglosUchwyt(EventType.ConfigChanged, () => {
      void odswiez();
    }),
    zglosUchwyt(EventType.QueueChanged, (tresc) => {
      if (tresc.queue.windowIds?.includes(idOkna) !== true) return;
      opiszKolejke(wezly, tresc.change === ChangeKind.Deleted ? null : tresc.queue.status);
    }),
  );

  void odswiez();
  return () => {
    for (const odlacz of odlaczenia) odlacz();
  };
}

function naniesStanOkna(wezly: WezlyStanu, stan: unknown, krok: string): void {
  if (typeof stan !== 'object' || stan === null) return;
  const pola = stan as { processStatus?: ProgressStatus; window?: { permissionMode?: PermissionMode } };
  if (pola.processStatus !== undefined) opiszBieg(wezly, pola.processStatus, krok);
  if (pola.window?.permissionMode !== undefined) opiszTryb(wezly, pola.window.permissionMode);
}

function opiszBieg(wezly: WezlyStanu, stan: ProgressStatus | null, krok: string): void {
  for (const plakietka of wezly.plakietki) {
    plakietka.hidden = stan === null;
    plakietka.replaceChildren();
    if (stan === null) continue;
    if (wezly.tetno !== null && stan === ProgressStatus.Running) {
      plakietka.appendChild(wezly.tetno.cloneNode(true));
    }
    plakietka.appendChild(
      document.createTextNode(krok === '' ? NAZWY_BIEGU[stan] : `${NAZWY_BIEGU[stan]} · ${krok}`),
    );
  }
}

function opiszKolejke(wezly: WezlyStanu, stan: QueueStatus | null): void {
  const stojaca = wezly.stan.querySelector<HTMLElement>('[data-kolejka]');
  if (stan === null) {
    stojaca?.remove();
    return;
  }
  const plakietka = stojaca ?? zbudujPlakietkeKolejki(wezly);
  if (plakietka === null) return;
  plakietka.textContent = `kolejka: ${NAZWY_KOLEJKI[stan]}`;
}

function zbudujPlakietkeKolejki(wezly: WezlyStanu): HTMLElement | null {
  const wzor = wezly.plakietki[0];
  if (wzor === undefined) return null;
  const plakietka = wzor.cloneNode(false) as HTMLElement;
  plakietka.hidden = false;
  plakietka.dataset.kolejka = '';
  wezly.stan.insertBefore(plakietka, wzor.nextSibling);
  return plakietka;
}

function opiszTryb(wezly: WezlyStanu, tryb: PermissionMode | null): void {
  const znak = wezly.tryb;
  if (znak === null) return;
  const napis = [...znak.childNodes].find((wezel) => wezel.nodeType === Node.TEXT_NODE);
  if (napis === undefined) return;
  napis.textContent = tryb === null ? '' : NAZWY_TRYBU[tryb];
}

// Stan biegu stoi w dwóch miejscach znacznika: na wstążce i w nagłówku okna komunikacji.
function zbierzWezly(korzen: ParentNode): WezlyStanu | null {
  const stan = korzen.querySelector('.st-wstazka-stan');
  if (!(stan instanceof HTMLElement)) return null;
  const wybor = '.st-wstazka-stan .dn-plakietka--sygnal, .sta-kom-stan .dn-plakietka--sygnal';
  const plakietki = [...korzen.querySelectorAll<HTMLElement>(wybor)];
  if (plakietki.length === 0) return null;
  for (const plakietka of plakietki) {
    const tetno = plakietka.querySelector('.pt-tetno');
    if (tetno !== null) wzorTetna ??= tetno.cloneNode(true) as Element;
  }
  const tryb = korzen.querySelector('.sta-chip--tryb');
  return { stan, plakietki, tetno: wzorTetna, tryb: tryb instanceof HTMLElement ? tryb : null };
}
