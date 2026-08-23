import { elementIkony, type NazwaIkony } from '../ikony/ikony';
import { SKROTY_ETYKIETA, etykietaSkrotu } from './etykiety-paneli';
import type { PozycjaMenu } from './menu-paneli';
import './menu.css';

/**
 * Rząd ikon-skrótów stojący przy `⋮` w nagłówku okna rozmowy.
 *
 * Zestaw skrótów przychodzi z zewnątrz w `pozycje`; ten moduł go nie wybiera.
 * Skrót powstaje wyłącznie dla pozycji otwieralnej — skrót do panelu, którego
 * nie ma czym zbudować, byłby atrapą. Gdy żadna pozycja nie jest otwieralna,
 * rząd zostaje pusty i nie zajmuje miejsca (`:empty` w `menu.css`); brak jest
 * nazwany w zdaniu pod wykazem menu, a nie udawany martwą ikoną.
 *
 * Dymek niesie nazwę i skrót w jednym wierszu — „Terminal  Ctrl+`" — przez
 * `title` i `aria-label`. Komponent `.dn-tooltip` z `komponenty/drobne.css`
 * wymaga własnego elementu treści wewnątrz przycisku i chmurki pozycjonowanej
 * nad nim, co w rzędzie ikon nagłówka kolidowałoby z listą menu na tej samej
 * warstwie; dymek przeglądarki wystarcza i nie dokłada warstwy.
 *
 * Skrót klawiaturowy jest tu wyłącznie napisem — nic w tym module nie nasłuchuje
 * klawiszy.
 *
 * Znacznik „coś nowego" zapala się tylko na żądanie z zewnątrz (`znacznik`);
 * moduł nie ma własnego źródła sygnału.
 *
 * Moduł nie zna kolumny paneli i nie otwiera niczego sam — woła `naWybor`. Stan
 * „otwarty" czyta z zewnątrz, żeby rząd i menu pokazywały ten sam stan.
 */
export interface OpcjeSkrotow {
  pozycje: readonly PozycjaMenu[];
  ikonaPozycji(kod: string): NazwaIkony;
  czyOtwarty(kod: string): boolean;
  naWybor(kod: string): void;
}

export interface SkrotyPaneli {
  element: HTMLElement;
  odswiez(): void;
  /** Zapala lub gasi znacznik „coś nowego" na skrócie danej pozycji. */
  znacznik(kod: string, widoczny: boolean): void;
}

export function utworzSkrotyPaneli(opcje: OpcjeSkrotow): SkrotyPaneli {
  const element = document.createElement('div');
  element.className = 'dn-skroty-paneli';
  element.setAttribute('role', 'group');
  element.setAttribute('aria-label', SKROTY_ETYKIETA);

  const przyciski = new Map<string, HTMLButtonElement>();

  for (const pozycja of opcje.pozycje) {
    const przycisk = document.createElement('button');
    przycisk.type = 'button';
    przycisk.className = 'dn-btn-ikona dn-skroty-paneli__przycisk';
    przycisk.dataset.panel = pozycja.kod;
    przycisk.append(elementIkony(opcje.ikonaPozycji(pozycja.kod), { rozmiar: 16 }));
    przycisk.addEventListener('click', () => opcje.naWybor(pozycja.kod));
    przyciski.set(pozycja.kod, przycisk);
    element.append(przycisk);
  }

  function odswiez(): void {
    for (const pozycja of opcje.pozycje) {
      const przycisk = przyciski.get(pozycja.kod);
      if (przycisk === undefined) continue;
      const otwarty = opcje.czyOtwarty(pozycja.kod);
      // Przełącznik, nie czynność jednorazowa — stan niesie aria-pressed.
      przycisk.setAttribute('aria-pressed', String(otwarty));
      const opis = etykietaSkrotu(pozycja.nazwa, otwarty);
      // Nazwa i skrót w jednym wierszu dymka.
      const dymek = pozycja.skrot === undefined ? opis : `${opis}  ${pozycja.skrot}`;
      przycisk.setAttribute('aria-label', dymek);
      przycisk.title = dymek;
    }
  }

  odswiez();

  return {
    element,
    odswiez,
    znacznik(kod, widoczny) {
      const przycisk = przyciski.get(kod);
      if (przycisk === undefined) return;
      if (widoczny) przycisk.dataset.znacznik = 'tak';
      else delete przycisk.dataset.znacznik;
    },
  };
}
