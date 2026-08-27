import type { ResearchReportSection } from '../../../../shared/contract';
import { przycisk, utworzWierszOdpowiedzi } from '../../modele/kontrolki-formularza';
import { AKCJE_RAPORTU } from './akcje-okien';
import {
  ustawStanRaportu,
  wykonajAkcjeRaportu,
  zlozRaport,
  type KontekstRaportu,
} from './czynnosci-raportu';
import { KODY_OKIEN } from './kody-okien';
import { utworzKreatorRaportu } from './kreator-raportu';
import { utworzRameBadania } from './rama-badania';
import type { StanBadania } from './stan-badania';
import { utworzStanOknaBadania } from './stan-okna-badania';

/**
 * Report Builder jest oknem kreatorem, składającym raport z ustaleń oraz redagującym jego treść, z sekcją podglądu osadzoną trwale i modalem kompozycji.
 */
export interface OknoReportBuilder {
  element: HTMLElement;
  odswiez(): void;
  /** Zamyka kreator i zdejmuje modal z dokumentu. */
  rozlacz(): void;
}

export function utworzOknoReportBuilder(
  stan: StanBadania,
  przejdz: (kodOkna: string) => void,
): OknoReportBuilder {
  const kontekst: KontekstRaportu = {
    stan,
    okno: utworzStanOknaBadania(),
    odpowiedz: utworzWierszOdpowiedzi(),
    kreator: utworzKreatorRaportu(() => void zlozRaport(kontekst)),
    przejdz,
  };

  const otworz = przycisk('Otwórz kreator raportu', 'dn-btn dn-btn--sm dn-btn--atrament');
  otworz.addEventListener('click', () => kontekst.kreator.otworz(stan.zakres()));

  const podglad = document.createElement('ol');
  podglad.className = 'mr-raport';

  // Konspekt stoi obok treści: struktura po lewej, treść po prawej; kliknięcie prowadzi ognisko.
  const konspekt = document.createElement('ol');
  konspekt.className = 'mr-konspekt';
  konspekt.setAttribute('aria-label', 'Konspekt raportu');

  const pokrycie = document.createElement('p');
  pokrycie.className = 'dn-pole-opis mr-raport__pokrycie';

  const podzial = document.createElement('div');
  podzial.className = 'mr-raport__podzial';
  podzial.append(konspekt, podglad);

  kontekst.okno.tresc.append(otworz, kontekst.odpowiedz.element, pokrycie, podzial);

  const rama = utworzRameBadania(
    KODY_OKIEN.raport,
    'Report Builder',
    'kreator',
    AKCJE_RAPORTU,
    (akcja) => void wykonajAkcjeRaportu(kontekst, akcja),
  );
  rama.cialo.append(kontekst.okno.element);

  function odswiez(): void {
    const sekcje = stan.raport()?.sections ?? [];
    podglad.replaceChildren(...sekcje.map(wierszSekcji));
    konspekt.replaceChildren(
      ...sekcje.map((sekcja) => pozycjaKonspektu(sekcja, () => ogniskujSekcje(sekcja.id))),
    );
    pokrycie.textContent = opisPokrycia(stan);
    ustawStanRaportu(kontekst);
  }

  /** Prowadzi wzrok do sekcji w podglądzie treści; sekcja zdjęta nie istnieje. */
  function ogniskujSekcje(identyfikator: string): void {
    const cel = podglad.querySelector<HTMLElement>(`[data-sekcja="${identyfikator}"]`);
    cel?.scrollIntoView({ block: 'nearest' });
  }

  odswiez();

  return {
    element: rama.element,
    odswiez,
    rozlacz() {
      kontekst.kreator.zamknij();
      kontekst.kreator.element.remove();
    },
  };
}

/**
 * Funkcja tworzy pozycję konspektu z nazwą sekcji oraz znacznikiem informującym, czy sekcja ma już wpisaną treść.
 */
function pozycjaKonspektu(
  sekcja: ResearchReportSection,
  naWybor: () => void,
): HTMLElement {
  const element = document.createElement('li');
  element.className = 'mr-konspekt__pozycja';

  const nazwa = sekcja.title === '' ? 'sekcja bez tytułu' : sekcja.title;
  const kontrolka = przycisk(nazwa, 'dn-btn dn-btn--sm dn-btn--duch');
  kontrolka.addEventListener('click', naWybor);

  const stan = document.createElement('span');
  const zTrescia = (sekcja.content ?? '').trim() !== '';
  stan.className = zTrescia
    ? 'dn-plakietka dn-plakietka--sukces'
    : 'dn-plakietka dn-plakietka--informacja';
  stan.textContent = zTrescia ? 'z treścią' : 'pusta';

  element.append(kontrolka, stan);
  return element;
}

/**
 * Pokrycie ustaleń w dokumencie — ile z zebranych weszło do sekcji.
 *
 * Liczone po zbiorze `findingIds` sekcji, bo to samo ustalenie może zasilać
 * więcej niż jedną sekcję i suma po sekcjach byłaby zawyżona.
 */
function opisPokrycia(stan: StanBadania): string {
  const sekcje = stan.raport()?.sections ?? [];
  const ustalenia = stan.ustalenia();
  if (ustalenia.length === 0) return 'Nie ma jeszcze ustaleń, którymi raport miałby być pokryty.';

  const objete = new Set<string>();
  for (const sekcja of sekcje) {
    for (const ustalenie of sekcja.findingIds ?? []) objete.add(ustalenie);
  }
  const wykorzystane = ustalenia.filter((ustalenie) => objete.has(ustalenie.id)).length;
  const procent = Math.round((wykorzystane / ustalenia.length) * 100);
  return (
    `Pokrycie ustaleń: ${String(wykorzystane)} z ${String(ustalenia.length)} (${String(procent)}%). ` +
    (wykorzystane === ustalenia.length
      ? 'Poza dokumentem nie zostało żadne.'
      : `Poza dokumentem zostaje ${String(ustalenia.length - wykorzystane)}.`)
  );
}

/** Funkcja tworzy jedną sekcję raportu widoczną w podglądzie treści okna, złożoną z tytułu oraz treści sekcji. */
function wierszSekcji(sekcja: ResearchReportSection): HTMLElement {
  const element = document.createElement('li');
  element.className = 'mr-raport__sekcja';
  element.dataset['sekcja'] = sekcja.id;

  const tytul = document.createElement('h5');
  tytul.className = 'mr-raport__tytul';
  tytul.textContent = sekcja.title === '' ? 'sekcja bez tytułu' : sekcja.title;

  const tresc = document.createElement('p');
  tresc.className = 'mr-raport__tresc';
  tresc.textContent = sekcja.content ?? '';

  element.append(tytul, tresc);
  return element;
}
