// Wiązanie karty modułu Library z rdzeniem: rozdziela panele prototypu między
// wywołania rodziny `library.*` i trzyma wspólne wskazanie pliku oraz kolekcji.
import { Command, EventType } from '../../../shared/contract.ts';
import type { Odsubskrybuj } from '../polaczenie/magistrala-zdarzen.ts';
import { zglosUchwyt } from '../polaczenie/rozdzielacz-zdarzen.ts';
import type { Kanal } from '../protokol/kanal.ts';
import { wywolaj } from '../protokol/wywolanie.ts';
import { zwiazEksplorator } from './library-eksplorator.ts';
import { zwiazEtykiety } from './library-etykiety.ts';
import { zwiazPodglad } from './library-podglad.ts';
import { zwiazWersje } from './library-wersje.ts';
import { zwiazZadania } from './library-zadania.ts';
import { utworzKontekst, wartosc } from './library-wspolne.ts';
import {
  uzgodnijPrzelacznikiPaneli,
  wpiszPole,
  zdejmijSterowanieWspolne,
  zdejmijTrescWspolna,
} from './okno-modulu.ts';

interface WiazanieKarty {
  korzen: Element;
  odlaczenia: Odsubskrybuj[];
}

const WIAZANIA = new Map<string, WiazanieKarty>();

export function zwiazLibrary(
  kanal: Kanal,
  nazwaSrodowiska: string,
  idOknaStojacego: string,
  wskazanieKorzenia: Element | string,
): boolean {
  const korzen = korzenKarty(wskazanieKorzenia);
  if (korzen === null) return false;
  const idKarty = korzen.getAttribute('data-karta') ?? '';
  if (idKarty === '') return false;
  if (WIAZANIA.get(idKarty)?.korzen === korzen) return false;
  const obszar = korzen.querySelector('.sta-obszar');
  if (!(obszar instanceof HTMLElement)) return false;
  zwolnijLibrary(idKarty);

  const odlaczenia: Odsubskrybuj[] = [];
  const sterowanie = new AbortController();
  const przy: AddEventListenerOptions = { signal: sterowanie.signal };
  odlaczenia.push(() => {
    sterowanie.abort();
  });
  WIAZANIA.set(idKarty, { korzen, odlaczenia });

  zdejmijTrescWspolna(obszar);
  zdejmijSterowanieWspolne(obszar);
  zdejmijTrescPrzykladowa(obszar);
  opiszGlowe(obszar, nazwaSrodowiska);

  const kontekst = utworzKontekst();
  zwiazEksplorator(kanal, obszar, kontekst, przy);
  zwiazEtykiety(kanal, obszar, kontekst, przy);
  zwiazPodglad(kanal, obszar, idOknaStojacego, kontekst, przy);
  zwiazWersje(kanal, obszar, kontekst, przy);
  zwiazZadania(kanal, obszar, kontekst, przy);
  uzgodnijPrzelacznikiPaneli(obszar, []);

  odlaczenia.push(
    zglosUchwyt(EventType.LibraryFileChanged, () => {
      kontekst.odswiez();
    }),
  );
  void opiszZakres(kanal, obszar);
  return true;
}

export function zwolnijLibrary(idKarty: string): void {
  const wiazanie = WIAZANIA.get(idKarty);
  if (wiazanie === undefined) return;
  for (const odlacz of wiazanie.odlaczenia) odlacz();
  WIAZANIA.delete(idKarty);
}

function korzenKarty(wskazanie: Element | string): Element | null {
  if (wskazanie instanceof Element) return wskazanie;
  return document.querySelector(`[data-karta="${wskazanie}"]`);
}

// Wpisy rozmowy, znaczniki źródeł i chipy czynności są treścią przykładową
// prototypu — rodzina `library.*` nie ma komendy, która by je odtworzyła.
function zdejmijTrescPrzykladowa(obszar: HTMLElement): void {
  obszar.querySelector('.sta-kom-stan')?.remove();
  obszar.querySelector('.sta-kom-kontekst')?.replaceChildren();
  obszar.querySelector('.sta-kontekst-akcji')?.replaceChildren();
  obszar.querySelector('.sta-kom-monitor')?.remove();
}

function opiszGlowe(obszar: HTMLElement, nazwaSrodowiska: string): void {
  const pola = [...obszar.querySelectorAll('.sta-kom-naglowek .sta-kom-pole')];
  wpiszPole(pola, 'Środowisko', nazwaSrodowiska);
  wpiszPole(pola, 'Model', '');
  wpiszPole(pola, 'Wysiłek', '');
}

async function opiszZakres(kanal: Kanal, obszar: HTMLElement): Promise<void> {
  const znak = obszar.querySelector('.sta-okno-belka .sta-chip');
  if (znak === null) return;
  const odpowiedz = await wywolaj(kanal, Command.LibraryStatsGet, {});
  const miary = wartosc(odpowiedz)?.stats ?? null;
  if (miary === null) {
    znak.remove();
    return;
  }
  znak.textContent = `Materiały: ${String(miary.fileCount)}`;
}
