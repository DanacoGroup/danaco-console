import { Command, type LibraryFile } from '../../../../shared/contract';
import type { WykazBrakow } from './braki-kontraktu';
import { czynnosciZbiorcze } from './biblioteka-czynnosci';
import { pozycjaPliku, wykazWersji } from './biblioteka-pozycja';
import { utworzRameOkna } from '../../komponenty/rama-okna';
import {
  pole,
  przyciskAkcji as przycisk,
  wiersz,
  wybor,
  wykaz,
} from '../../modele/kontrolki-formularza';
import type { StanProjektu } from './stan-projektu';
import { utworzStanTresci } from './stany-okna';
import type { ZrodloWorkspace } from './zrodlo-workspace';

/**
 * Project Library to okno zarządcy modułu Workspace ograniczone do zakresu projektu: przegląd
 * z wyszukiwaniem, wgranie i praca zbiorowa na zaznaczeniu.
 */
export interface OknoBiblioteki {
  element: HTMLElement;
  odswiez(): void;
}

/** Moduły docelowe udostępnienia to te, które pracują na materiale projektu, wymienione jako pary kodu i etykiety modułu. */
const MODULY_DOCELOWE: ReadonlyArray<[string, string]> = [
  ['studio', 'Studio — edycja dokumentu'],
  ['research', 'Research — materiał badania'],
  ['translate', 'Translate — tłumaczenie'],
  ['library', 'Library — repozytorium wiedzy'],
  ['developer', 'Developer — praca w kodzie'],
];

export function utworzOknoBiblioteki(
  zrodlo: ZrodloWorkspace,
  stan: StanProjektu,
  braki: WykazBrakow,
): OknoBiblioteki {
  const rama = utworzRameOkna({
    tytul: 'Project Library',
    kod: 'project-library',
    rola: 'zarządca',
    przeznaczenie: 'Zasoby projektu: przegląd, wgranie, wersje i udostępnienie do modułu zewnętrznego.',
    przedrostek: 'dw',
  });
  const tresc = utworzStanTresci();

  const szukaj = pole('Fraza wyszukiwania', 'fragment nazwy albo ścieżki');
  const modul = wybor('Moduł docelowy udostępnienia', MODULY_DOCELOWE.map(([kod, opis]) => [kod, opis]));
  const etykieta = pole('Etykieta zbiorcza', 'np. pisma-procesowe');

  const odswiez = przycisk('Odśwież wykaz');
  const wgraj = przycisk('+ Wgraj plik', 'dn-btn dn-btn--atrament');
  const plik = document.createElement('input');
  plik.type = 'file';
  plik.hidden = true;
  const udostepnij = przycisk('Udostępnij zaznaczone');
  const oznacz = przycisk('Nadaj etykietę zaznaczonym');

  rama.akcje.append(
    wgraj,
    odswiez,
    udostepnij,
    oznacz,
    braki.przyciskBraku(
      'Pobierz zaznaczone',
      'Pobranie treści plików projektu',
      Command.LibraryFilePreview,
    ),
    braki.przyciskBraku(
      'Przenieś zaznaczone',
      'Przeniesienie pliku między katalogami',
      Command.LibraryFileMove,
    ),
    braki.przyciskBraku(
      'Usuń zaznaczone',
      'Usunięcie pliku z repozytorium',
      Command.LibraryFileDelete,
    ),
    plik,
  );
  rama.cialo.append(
    wiersz('Szukaj w projekcie', szukaj, { klasa: 'dw-wiersz', objasnienie: 'Fraza zawęża wykaz po nazwie i ścieżce pliku.' }),
    wiersz('Moduł docelowy', modul, { klasa: 'dw-wiersz', objasnienie: 'Udostępnienie przenosi zaznaczone pliki wraz z projektem.' }),
    wiersz('Etykieta', etykieta, { klasa: 'dw-wiersz', objasnienie: 'Etykieta zbiorcza idzie komendą library.tag.set po jednym pliku.' }),
    tresc.element,
  );

  const zaznaczone = new Set<string>();

  function pozycja(wpis: LibraryFile): HTMLElement {
    return pozycjaPliku(wpis, {
      zaznaczony: zaznaczone.has(wpis.id),
      przelacz(zaznacz) {
        if (zaznacz) zaznaczone.add(wpis.id);
        else zaznaczone.delete(wpis.id);
      },
      pokazWersje(): void {
        tresc.ladowanie(`Odczyt wersji pliku ${wpis.name}…`);
        void zrodlo.wersjePliku(wpis.id).then((wynik) => {
          if (!wynik.udany || wynik.wynik === undefined) {
            tresc.blad(`Rdzeń nie oddał wersji pliku ${wpis.name}.`, wynik.blad);
            return;
          }
          if (wynik.wynik.length === 0) {
            tresc.pusto(`Plik ${wpis.name} nie ma zapisanych wersji.`);
            return;
          }
          tresc.tresc().append(wykazWersji(wpis.name, wynik.wynik));
        });
      },
    });
  }

  /** Odczyt wykazu zasobów oddaje to, co narysował, bo czynności zbiorcze sprawdzają po nim skutek pracy. */
  async function odczytaj(): Promise<readonly LibraryFile[] | null> {
    const projekt = stan.projekt();
    if (projekt === '') {
      rama.ustawZnacznik('');
      tresc.pusto('Wskaż projekt na pulpicie, aby zobaczyć jego zasoby.');
      return null;
    }
    tresc.ladowanie('Odczyt zasobów projektu…');
    const wynik = await zrodlo.biblioteka({ projectId: projekt, query: szukaj.value });
    if (!wynik.udany || wynik.wynik === undefined) {
      // Licznik gaśnie razem z wykazem: liczba sprzed odmowy mówiłaby o zasobach, których nie odczytano.
      rama.ustawZnacznik('');
      tresc.blad('Rdzeń nie oddał zasobów projektu.', wynik.blad);
      return null;
    }
    // Licznik mówi, ile pozycji rdzeń oddał na pytanie: przy frazie to liczba trafień, nie plików.
    rama.ustawZnacznik(`pliki: ${wynik.wynik.length}`);
    if (wynik.wynik.length === 0) {
      tresc.pusto(
        szukaj.value === ''
          ? 'Projekt nie ma jeszcze plików. Pierwszy możesz wgrać przyciskiem „+ Wgraj plik”.'
          : `Żaden plik projektu nie odpowiada frazie „${szukaj.value}”.`,
      );
      return wynik.wynik;
    }
    const lista = wykaz('Zasoby projektu', 'dw-wykaz');
    for (const wpis of wynik.wynik) lista.append(pozycja(wpis));
    tresc.tresc().append(lista);
    return wynik.wynik;
  }

  const zbiorcze = czynnosciZbiorcze({ zrodlo, stan, tresc, zaznaczone, odczytaj });

  wgraj.addEventListener('click', () => plik.click());
  plik.addEventListener('change', () => {
    const wybrany = plik.files?.[0];
    if (wybrany !== undefined) zbiorcze.wgraj(wybrany);
  });
  udostepnij.addEventListener('click', () => zbiorcze.udostepnij(modul.value));
  oznacz.addEventListener('click', () => zbiorcze.oznacz(etykieta.value));

  szukaj.addEventListener('change', () => void odczytaj());
  odswiez.addEventListener('click', () => void odczytaj());
  // Zmiana projektu unieważnia zaznaczenie: identyfikatory plików są ścieżkami w obrębie projektu.
  stan.naZmiane(() => {
    zaznaczone.clear();
    void odczytaj();
  });

  return {
    element: rama.element,
    odswiez() {
      void odczytaj();
    },
  };
}
