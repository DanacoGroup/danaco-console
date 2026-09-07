// Notatki na przeglądanej stronie: założenie, poprawa treści oraz wątki, do
// których notatka należy. Panel wystawiał dotąd sam wykaz bez żadnej drogi.

import { Command, type BrowserNote } from '../../../shared/contract.ts';
import type { Kanal } from '../protokol/kanal.ts';
import { wywolaj } from '../protokol/wywolanie.ts';
import {
  dolozPrzycisk,
  odmowa,
  pasNadCialem,
  polePasa,
  powiedz,
} from './browser-wspolne.ts';
import { postawGrupeSzyny, szynaModulu } from './browser-wytwory.ts';

export function zwiazNotatkiPrzegladania(
  kanal: Kanal,
  korzen: Element,
  idOkna: string,
  przy: AddEventListenerOptions,
  odswiezNotatki: () => void,
): void {
  const szyna = szynaModulu(korzen);
  const cialo = korzen.querySelector('#panel-notes .sta-okno-tresc');
  const pas = pasNadCialem(cialo, 'Notatki strony');
  const poleTresci = polePasa(pas, 'Treść notatki');
  const poleWatku = polePasa(pas, 'Nazwa wątku');
  dolozPrzycisk(pas, 'Zapisz notatkę', 'notatkaZaloz');
  void postawWatki(kanal, szyna, idOkna);

  korzen.addEventListener('click', (zdarzenie) => {
    const cel = zdarzenie.target;
    if (!(cel instanceof Element)) return;

    if (cel.closest('[data-notatka-zaloz]') !== null) {
      void zalozNotatke(kanal, idOkna, poleTresci, odswiezNotatki);
      return;
    }
    const poprawa = cel.closest<HTMLElement>('[data-notatka-popraw]');
    if (poprawa !== null) {
      otworzPoprawe(kanal, poprawa, odswiezNotatki);
      return;
    }
    const doWatku = cel.closest<HTMLElement>('[data-notatka-do-watku]');
    if (doWatku !== null) {
      void wstawNotatkeDoWatku(kanal, szyna, idOkna, doWatku.dataset.notatkaDoWatku ?? '',
        (poleWatku?.value ?? '').trim(), odswiezNotatki);
    }
  }, przy);
}

/** Dokłada w wierszu notatki drogi poprawy treści i przypisania do wątku. */
export function dolozCzynnosciNotatki(wiersz: HTMLElement, notatka: BrowserNote): void {
  const tresc = document.createElement('span');
  tresc.className = 'brw-cytat';
  tresc.dataset.notatkaTresc = notatka.id;
  tresc.textContent = notatka.content;
  const rzad = document.createElement('div');
  rzad.className = 'sta-chip-rzad';
  dolozPrzycisk(rzad, 'Popraw', 'notatkaPopraw', notatka.id);
  dolozPrzycisk(rzad, 'Do wątku', 'notatkaDoWatku', notatka.id);
  wiersz.replaceChildren(tresc, rzad);
}

async function zalozNotatke(
  kanal: Kanal,
  idOkna: string,
  pole: HTMLInputElement | null,
  odswiez: () => void,
): Promise<void> {
  const content = (pole?.value ?? '').trim();
  if (content === '') {
    odmowa(undefined, 'Notatka bez treści nie powstanie — pole treści jest puste.');
    return;
  }
  const wynik = await wywolaj(kanal, Command.BrowserNoteAdd, { windowId: idOkna, content });
  if (!wynik.udany) {
    odmowa(wynik.blad, 'Rdzeń odmówił zapisu notatki.');
    return;
  }
  if (pole !== null) pole.value = '';
  powiedz('Notatka stoi w wykazie.');
  odswiez();
}

/* Poprawa idzie w miejscu, na wierszu wykazu: osobne okno dialogowe nie ma
   gdzie stanąć, a treść notatki bywa dłuższa niż pole paska. */
function otworzPoprawe(kanal: Kanal, przycisk: HTMLElement, odswiez: () => void): void {
  const noteId = przycisk.dataset.notatkaPopraw ?? '';
  const wiersz = przycisk.closest('.br-notatka');
  const wezel = wiersz?.querySelector<HTMLElement>(`[data-notatka-tresc="${noteId}"]`) ?? null;
  if (noteId === '' || wezel === null || wezel.isContentEditable) return;
  const przed = wezel.textContent ?? '';
  wezel.contentEditable = 'true';
  wezel.focus();
  wezel.addEventListener('keydown', (zdarzenie) => {
    if (zdarzenie.key === 'Escape') {
      wezel.textContent = przed;
      wezel.removeAttribute('contenteditable');
      return;
    }
    if (zdarzenie.key !== 'Enter' || zdarzenie.shiftKey) return;
    zdarzenie.preventDefault();
    wezel.removeAttribute('contenteditable');
    void zapiszPoprawe(kanal, noteId, (wezel.textContent ?? '').trim(), przed, odswiez);
  });
}

async function zapiszPoprawe(
  kanal: Kanal,
  noteId: string,
  content: string,
  przed: string,
  odswiez: () => void,
): Promise<void> {
  if (content === '' || content === przed) {
    odswiez();
    return;
  }
  const wynik = await wywolaj(kanal, Command.BrowserNoteUpdate, { noteId, content });
  if (!wynik.udany) {
    odmowa(wynik.blad, 'Rdzeń odmówił poprawy notatki.');
  } else {
    powiedz('Notatka poprawiona.');
  }
  odswiez();
}

/* Kontrakt nie zna dokładania notatki do wątku: zapis idzie całym wątkiem,
   więc skład bierze się z wykazu wątków i dochodzi do niego notatka wskazana. */
async function wstawNotatkeDoWatku(
  kanal: Kanal,
  szyna: Element | null,
  idOkna: string,
  idNotatki: string,
  nazwa: string,
  odswiez: () => void,
): Promise<void> {
  if (idNotatki === '') return;
  if (nazwa === '') {
    odmowa(undefined, 'Wątek wymaga nazwy — pole nazwy wątku jest puste.');
    return;
  }
  const wykaz = await wywolaj(kanal, Command.BrowserNoteThreadList, { windowId: idOkna });
  const watek = wykaz.wynik?.threads.find((pozycja) => pozycja.name === nazwa);
  const sklad = new Set(watek?.noteIds ?? []);
  sklad.add(idNotatki);
  const wynik = await wywolaj(kanal, Command.BrowserNoteThreadSet, {
    windowId: idOkna,
    threadId: watek?.id,
    name: nazwa,
    noteIds: [...sklad],
  });
  if (!wynik.udany) {
    odmowa(wynik.blad, 'Rdzeń odmówił zapisu wątku notatek.');
    return;
  }
  powiedz(`Notatka stoi w wątku „${wynik.wynik?.thread.name ?? nazwa}”.`);
  void postawWatki(kanal, szyna, idOkna);
  odswiez();
}

async function postawWatki(kanal: Kanal, szyna: Element | null, idOkna: string): Promise<void> {
  if (szyna === null) return;
  const wynik = await wywolaj(kanal, Command.BrowserNoteThreadList, { windowId: idOkna });
  if (!wynik.udany) return;
  const watki = wynik.wynik?.threads ?? [];
  postawGrupeSzyny(szyna, 'watki', 'Wątki notatek', watki.map((watek) => ({
    tytul: `${watek.name} — notatek: ${watek.noteIds?.length ?? 0}`,
    cechy: { watekNotatek: watek.id },
  })));
}
