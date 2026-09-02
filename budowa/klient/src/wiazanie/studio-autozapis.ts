// Autozapis dokumentu Studia w pasie stanu kanwy: nastawy rdzenia, zapis
// samoczynny w podanym przez nie odstępie i przełączenie znakiem stanu.

import {
  Command,
  EventType,
  StudioBackupReason,
  type StudioAutosaveSettings,
  type StudioDocument,
} from '../../../shared/contract.ts';
import type { Odsubskrybuj } from '../polaczenie/magistrala-zdarzen.ts';
import { zglosUchwyt } from '../polaczenie/rozdzielacz-zdarzen.ts';
import type { Kanal } from '../protokol/kanal.ts';
import { wywolaj } from '../protokol/wywolanie.ts';
import { oglos } from './ogloszenie.ts';
import { nazwijOdmowe, wskazDokumentOkna } from './studio-dokument.ts';

interface WezlyAutozapisu {
  znak: HTMLElement;
  napis: ChildNode;
  kanwa: HTMLElement;
}

export function zdejmijTrescPrzykladowaAutozapisu(korzen: ParentNode): boolean {
  const wezly = zbierzWezly(korzen);
  if (wezly === null) return false;
  wezly.napis.textContent = '';
  return true;
}

// Odstęp zapisu bierze się z nastaw rdzenia; nastawy bez odstępu zostawiają
// sam opis stanu, bo klient nie ma czym rozstrzygnąć, jak często zapisywać.
export function zwiazAutozapis(
  kanal: Kanal,
  idOkna: string,
  korzen: ParentNode,
): Odsubskrybuj | null {
  const znalezione = zbierzWezly(korzen);
  if (znalezione === null) return null;
  const wezly: WezlyAutozapisu = znalezione;
  const sterowanie = new AbortController();
  let dokument: StudioDocument | null = null;
  let nastawy: StudioAutosaveSettings | null = null;
  let odliczanie = 0;

  const opisz = (): void => {
    wezly.napis.textContent = opisStanu(nastawy);
  };

  const zatrzymajOdliczanie = (): void => {
    if (odliczanie === 0) return;
    globalThis.clearInterval(odliczanie);
    odliczanie = 0;
  };

  const zapisz = async (): Promise<void> => {
    if (dokument === null || nastawy?.enabled !== true) return;
    const wynik = await wywolaj(kanal, Command.StudioAutosaveRun, {
      documentId: dokument.id,
      content: wezly.kanwa.textContent ?? '',
      trigger: StudioBackupReason.Interval,
    });
    if (!wynik.udany || wynik.wynik === undefined) {
      oglos('Studio', nazwijOdmowe('Zapis samoczynny', wynik.blad), 'ostrzezenie');
      return;
    }
    nastawy = {
      ...nastawy,
      lastSaveAt: wynik.wynik.savedAt,
      lastSaveFailed: !wynik.wynik.saved,
      lastFailureReason: wynik.wynik.failureReason,
    };
    opisz();
    if (!wynik.wynik.saved) {
      oglos('Studio', `Zapis samoczynny nieudany — ${wynik.wynik.failureReason ?? 'bez powodu'}.`,
        'ostrzezenie');
    }
  };

  const ustawOdliczanie = (): void => {
    zatrzymajOdliczanie();
    const odstep = nastawy?.enabled === true ? nastawy.intervalSeconds ?? 0 : 0;
    if (odstep <= 0) return;
    odliczanie = globalThis.setInterval(() => {
      void zapisz();
    }, odstep * 1000);
  };

  const przyjmijNastawy = (przyjete: StudioAutosaveSettings): void => {
    nastawy = przyjete;
    opisz();
    ustawOdliczanie();
  };

  const przelacz = async (): Promise<void> => {
    if (nastawy === null) return;
    const wynik = await wywolaj(kanal, Command.StudioAutosaveSet, {
      documentId: dokument?.id,
      windowId: idOkna,
      enabled: !nastawy.enabled,
    });
    if (!wynik.udany || wynik.wynik === undefined) {
      oglos('Studio', nazwijOdmowe('Przestawienie autozapisu', wynik.blad), 'ostrzezenie');
      return;
    }
    przyjmijNastawy(wynik.wynik.settings);
  };

  wezly.znak.setAttribute('role', 'button');
  wezly.znak.tabIndex = 0;
  wezly.znak.addEventListener('click', () => {
    void przelacz();
  }, { signal: sterowanie.signal });
  wezly.znak.addEventListener('keydown', (zdarzenie) => {
    if (zdarzenie.key !== 'Enter' && zdarzenie.key !== ' ') return;
    zdarzenie.preventDefault();
    void przelacz();
  }, { signal: sterowanie.signal });

  const odlacz = zglosUchwyt(EventType.StudioDocumentChanged, (tresc) => {
    if (tresc.document.windowId !== idOkna) return;
    dokument = tresc.document;
  });

  void wskazDokumentOkna(kanal, idOkna).then(async (otwarty) => {
    dokument = otwarty;
    const wynik = await wywolaj(kanal, Command.StudioAutosaveGet, {
      documentId: otwarty?.id,
      windowId: idOkna,
    });
    if (!wynik.udany || wynik.wynik === undefined) return;
    przyjmijNastawy(wynik.wynik.settings);
  });

  return () => {
    sterowanie.abort();
    zatrzymajOdliczanie();
    odlacz();
  };
}

function opisStanu(nastawy: StudioAutosaveSettings | null): string {
  if (nastawy === null) return '';
  if (!nastawy.enabled) return ' autozapis wyłączony';
  if (nastawy.lastSaveFailed === true) return ' zapis samoczynny nieudany';
  if (nastawy.lastSaveAt === undefined) return ' autozapis czynny';
  const godzina = new Date(nastawy.lastSaveAt).toLocaleTimeString('pl-PL', {
    hour: '2-digit',
    minute: '2-digit',
  });
  return ` zapisano samoczynnie ${godzina}`;
}

// Znak autozapisu to pozycja pasa stanu niosąca tętno; napis stoi osobnym węzłem.
function zbierzWezly(korzen: ParentNode): WezlyAutozapisu | null {
  const kanwa = korzen.querySelector('.dn-kanwa');
  const znak = [...(korzen.querySelector('.st-status')?.children ?? [])].find(
    (pozycja) => pozycja.querySelector('.pt-tetno') !== null,
  );
  if (!(znak instanceof HTMLElement) || !(kanwa instanceof HTMLElement)) return null;
  const napis = [...znak.childNodes].find((wezel) => wezel.nodeType === Node.TEXT_NODE);
  return napis === undefined ? null : { znak, napis, kanwa };
}
