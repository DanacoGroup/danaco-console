import type { LibraryVersion } from '../../../../shared/contract';
import { utworzDymekObjasnienia } from '../../komponenty/dymek';
import { opisOdmowy } from '../../komponenty/odmowa';
import { utworzRameOkna } from '../../komponenty/rama-okna';
import { pobierzPlik, przycisk, utworzWierszOdpowiedzi } from '../../modele/kontrolki-formularza';
import { utworzDolozenieWersji } from './dolozenie-wersji';
import { TRESC_NIEZNANA, type WerdyktTresci } from './dostepnosc-tresci';
import { raportZmian } from './eksporty-biblioteki';
import { BEZ_KOMENDY_WERSJE } from './etykiety-biblioteki';
import { nazwaZakresu, utworzFiltrWersji } from './filtr-wersji';
import { utworzPanelAkcji } from './panel-akcji';
import type { StanBiblioteki } from './stan-biblioteki';
import { utworzStanOkna } from './stan-okna';
import { utworzWierszWersji } from './wiersz-wersji';
import type { ZrodloOtoczenia } from './zrodlo-otoczenia';

/**
 * Panel Versioning Panel łączy trzy czynności: przegląd wersji, dołożenie
 * kolejnej i przywrócenie wcześniejszej, z odmową rdzenia pokazaną wprost
 * zamiast pustej tabeli, a filtr historii i raport zmian nie wysyłają komend.
 */
export interface OknoWersji {
  element: HTMLElement;
  odswiez(): void;
}

export function utworzOknoWersji(stan: StanBiblioteki, otoczenie: ZrodloOtoczenia): OknoWersji {
  const rama = utworzRameOkna({
    kod: 'versioning-panel',
    tytul: 'Versioning Panel',
    rola: 'zarządca',
    modul: 'Library',
    dodatkiNaglowka: [
      utworzDymekObjasnienia(
        'Historia wersji dokumentu wybranego w Library Explorer. Wersje narastają przy każdej ' +
          'zmianie w dowolnym module, w tym przy zmianie wykonanej przez model.',
        { powloka: 'ml-dymek', znak: 'ml-dymek__znak' },
      ),
    ],
  });
  const okno = utworzStanOkna();
  const odpowiedz = utworzWierszOdpowiedzi();
  const panel = utworzPanelAkcji(otoczenie, stan, BEZ_KOMENDY_WERSJE);
  const dolozenie = utworzDolozenieWersji(stan, () => void wczytaj());

  const lista = document.createElement('ul');
  lista.className = 'ml-wersje';

  const filtr = utworzFiltrWersji(() => wypelnij(wersjeWykazu));

  const odswiezHistorie = przycisk('Odśwież historię', 'dn-btn dn-btn--sm dn-btn--zarys');
  odswiezHistorie.dataset['czynnosc'] = 'wersje';
  odswiezHistorie.addEventListener('click', () => void wczytaj());

  const raport = przycisk('Raport zmian', 'dn-btn dn-btn--sm dn-btn--duch');
  raport.dataset['czynnosc'] = 'raport-zmian';
  raport.addEventListener('click', () => wywiezRaport());

  // Raport zmian składa się z historii pokazanej w oknie, bo tę treść da się potwierdzić przed zapisem.
  function wywiezRaport(): void {
    const plik = stan.czynny();
    if (plik === null) {
      odpowiedz.pokaz('Wskaż dokument w Library Explorer — raport dotyczy jednej historii.', false);
      return;
    }
    const wybrane = filtr.zastosuj(wersjeWykazu);
    if (wybrane.length === 0) {
      odpowiedz.pokaz('Historia po zawężeniu jest pusta — raport nie miałby ani jednej pozycji.', false);
      return;
    }
    pobierzPlik(`historia-${plik.id}.md`, raportZmian(plik, wybrane), 'text/markdown');
    odpowiedz.pokaz(
      `Raport zmian złożony z pozycji widocznych w oknie: ${wybrane.length} ` +
        `(zawężenie: ${nazwaZakresu(filtr.zakres())}). Treści wersji raport nie zawiera.`,
      true,
    );
  }

  let ostatni = '';
  // Wykaz wersji i werdykt o treści są trzymane osobno, bo odpowiedź o treści przychodzi po historii.
  let wersjeWykazu: readonly LibraryVersion[] = [];
  let werdyktWykazu: WerdyktTresci = 'nieznana';

  async function wczytaj(): Promise<void> {
    const plik = stan.czynny();
    if (plik === null) return;
    okno.ladowanie(`Rdzeń odczytuje wersje pliku „${plik.name}".`);
    const wynik = await stan.zrodlo.wersje(plik.id);
    if (!wynik.udany || wynik.wynik === undefined) {
      okno.blad(opisOdmowy('Odczyt wersji dokumentu', wynik.blad?.code, wynik.blad?.message));
      return;
    }
    wypelnij(wynik.wynik.versions);
  }

  /** Rysuje same wiersze — bez ruszania fazy okna, którą ustawia odczyt. */
  function rysuj(): void {
    const plik = stan.czynny();
    const tresc = plik === null ? TRESC_NIEZNANA : stan.tresc(plik.id);
    werdyktWykazu = tresc.werdykt;
    lista.replaceChildren(
      ...filtr.zastosuj(wersjeWykazu).map((wersja) =>
        utworzWierszWersji(wersja, {
          biezaca: plik?.versionId === wersja.id,
          tresc,
          naPrzywrocenie: () => void przywroc(wersja),
        }),
      ),
    );
  }

  function wypelnij(wersje: readonly LibraryVersion[]): void {
    wersjeWykazu = wersje;
    rysuj();
    if (wersje.length === 0) {
      okno.puste(
        'Dokument bez historii',
        'Rdzeń nie ma ani jednej wersji tego dokumentu. Dołóż wersję formularzem niżej ' +
          'albo zmień dokument w dowolnym module — historia narasta sama.',
      );
      return;
    }
    // Historia zawężona filtrem i pusta to dwa różne stany: zawężenie okna, a nie odpowiedź rdzenia.
    if (filtr.zastosuj(wersje).length === 0) {
      okno.puste(
        'Zawężenie filtru nie zostawiło ani jednej wersji',
        `Rdzeń oddał wersji: ${wersje.length}, a zawężenie „${nazwaZakresu(filtr.zakres())}" ` +
          'nie obejmuje żadnej z nich. Sprawcę zmiany zapisuje moduł, który ją wykonał — ' +
          'wersja bez wypełnionego pola sprawcy nie wchodzi do żadnej z dwóch grup.',
      );
      return;
    }
    okno.gotowe();
  }

  async function przywroc(wersja: LibraryVersion): Promise<void> {
    const plik = stan.czynny();
    if (plik === null) return;
    odpowiedz.pokaz(`Przywracanie wersji ${wersja.label ?? wersja.id}…`, true);
    const wynik = await stan.zrodlo.przywroc(plik.id, wersja.id);
    if (!wynik.udany || wynik.wynik === undefined) {
      odpowiedz.pokaz(
        opisOdmowy('Przywrócenie wersji', wynik.blad?.code, wynik.blad?.message),
        false,
      );
      return;
    }
    const nazwa = wersja.label ?? wersja.id;
    const poPrzywroceniu = wynik.wynik.file;
    stan.wchlon(poPrzywroceniu);
    void wczytaj();

    // Którą wersję rdzeń uznał za bieżącą, mówi odpowiedź po przywróceniu, nie kliknięty wiersz.
    if (poPrzywroceniu.versionId !== wersja.id) {
      odpowiedz.pokaz(
        `Przywrócenie wersji ${nazwa} (${wersja.id}) wróciło powodzeniem, ale dokument ` +
          `oddany przez rdzeń wskazuje wersję ${poPrzywroceniu.versionId ?? 'nieokreśloną'}. ` +
          'Nie potwierdzam przywrócenia.',
        false,
      );
      return;
    }
    // Przywrócenie przestawia wersję bieżącą, więc poprzednia odpowiedź o treści przestała o niej mówić.
    odpowiedz.pokaz(`Wersja ${nazwa} wskazana jako bieżąca — pytam rdzeń o treść…`, true);
    const stanTresci = await stan.zbadajTresc(plik.id);
    if (stanTresci.werdykt === 'brak') {
      odpowiedz.pokaz(
        `Wersja ${nazwa} stoi jako bieżąca, ale dokument nadal nie ma treści ` +
          `w repozytorium — ${stanTresci.powod}.`,
        false,
      );
      return;
    }
    if (stanTresci.werdykt === 'odwolanie') {
      odpowiedz.pokaz(
        `Wersja ${nazwa} stoi jako bieżąca, ale na pytanie o treść rdzeń oddał wskazanie ` +
          `zamiast bajtów — ${stanTresci.powod}. Powrotu do treści nie potwierdzam.`,
        false,
      );
      return;
    }
    // Zdanie o oddaniu treści dokumentu należy się wyłącznie werdyktowi osiągalna, nie każdej odpowiedzi.
    if (stanTresci.werdykt !== 'osiagalna') {
      odpowiedz.pokaz(
        `Wersja ${nazwa} wskazana jako bieżąca, ale rdzeń nie odpowiedział o treści ` +
          `dokumentu — ${stanTresci.powod}. Powrót do treści pozostaje niepotwierdzony.`,
        false,
      );
      return;
    }
    odpowiedz.pokaz(`Wersja ${nazwa} przywrócona; rdzeń oddaje treść dokumentu.`, true);
  }

  okno.tresc.append(lista);
  rama.pasek.append(odswiezHistorie, raport);
  rama.narzedzia.append(filtr.element);
  rama.cialo.append(okno.element, dolozenie.element, odpowiedz.element, panel.element);

  return {
    element: rama.element,

    odswiez() {
      const plik = stan.czynny();
      if (plik === null) {
        wersjeWykazu = [];
        lista.replaceChildren();
        ostatni = '';
        okno.puste('Nic nie wybrano', 'Wskaż dokument w Library Explorer, aby zobaczyć wersje.');
        return;
      }
      if (plik.id === ostatni) {
        // Historia się nie zmieniła, ale werdykt o treści mógł dojść później — wystarczy przerysować wiersze.
        if (stan.tresc(plik.id).werdykt !== werdyktWykazu) rysuj();
        return;
      }
      ostatni = plik.id;
      void wczytaj();
    },
  };
}
