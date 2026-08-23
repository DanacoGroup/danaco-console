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
 * Versioning Panel — okno zarządcy modułu (`library.versioning-panel`).
 *
 * Trzy funkcje Operatora mają tu drogę do rdzenia: przegląd wersji
 * (`library.version.list`), dołożenie kolejnej (`library.version.add`,
 * `dolozenie-wersji.ts`) i przywrócenie wcześniejszej
 * (`library.version.restore`). Odmowa rdzenia zostaje pokazana wprost, zamiast
 * pustej tabeli. Dołożenie wersji jest tu potrzebne, bo wgranie pliku zakłada
 * nowy dokument — bez niego historia nie rośnie ponad jeden wpis, a „Przywróć"
 * nie ma dokąd wracać.
 *
 * Filtr historii i raport zmian nie wysyłają komendy: `library.version.list`
 * nie przyjmuje pola zawężającego, a raport składa się z odpowiedzi, którą okno
 * już ma. Wywóz historii wraz z treścią każdej wersji nie powstaje — treści
 * wersji niebieżącej nie oddaje żadna komenda kontraktu.
 *
 * To nie jest Session Repository ze Studia: tam wersje żyją w toku sesji, tu —
 * w repozytorium biblioteki, i narastają przy zmianie dokumentu w dowolnym
 * module, także przy zmianie wykonanej przez model. Dlatego panel odświeża się
 * także zdarzeniem `library.file.changed`, a nie wyłącznie własnym działaniem.
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

  /**
   * Raport zmian składa się z historii pokazanej w oknie, nie z osobnego
   * odczytu: wywóz tego, co Operator widzi, jest jedynym wywozem, którego
   * treść da się potwierdzić przed zapisaniem pliku. Zawężenie filtrem wchodzi
   * do raportu i jest w nim nazwane.
   */
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
  // Wykaz wersji i werdykt o treści, z którymi jest narysowany. Trzymane, bo
  // odpowiedź rdzenia o treści przychodzi później niż sama historia (osobnym
  // wywołaniem, także z okna File Preview) — bez tego wiersze pokazywałyby
  // nieaktualny stan aż do następnego odczytu historii.
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
    // Historia zawężona filtrem i historia pusta to dwa różne stany: pierwszy
    // orzeka o zawężeniu okna, drugi o odpowiedzi rdzenia.
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

    // Która wersja jest teraz bieżąca, mówi odpowiedź, nie kliknięcie. Rdzeń
    // oddaje plik po przywróceniu, a w nim `versionId` — jedyne miejsce, w
    // którym stoi wersja naprawdę wskazana. Rozbieżność z klikniętym wierszem
    // jest tu ogłaszana odmową.
    if (poPrzywroceniu.versionId !== wersja.id) {
      odpowiedz.pokaz(
        `Przywrócenie wersji ${nazwa} (${wersja.id}) wróciło powodzeniem, ale dokument ` +
          `oddany przez rdzeń wskazuje wersję ${poPrzywroceniu.versionId ?? 'nieokreśloną'}. ` +
          'Nie potwierdzam przywrócenia.',
        false,
      );
      return;
    }
    // Przywrócenie przestawia wersję bieżącą dokumentu, więc poprzednia
    // odpowiedź rdzenia o jego treści przestała o nim mówić. Zdanie „wersja
    // przywrócona" bez tego pytania potwierdzałoby powrót do treści, której
    // repozytorium może nie mieć.
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
    // Zdanie „rdzeń oddaje treść dokumentu" jest potwierdzeniem i należy się
    // wyłącznie werdyktowi `osiagalna`. Po odmowie odczytu okno nie potwierdza
    // powrotu do treści, ale też nie orzeka jej braku.
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
        // Historia się nie zmieniła, ale werdykt o treści mógł dojść później —
        // wtedy wystarczy przerysować wiersze, bez powtórnego odczytu wersji.
        if (stan.tresc(plik.id).werdykt !== werdyktWykazu) rysuj();
        return;
      }
      ostatni = plik.id;
      void wczytaj();
    },
  };
}
