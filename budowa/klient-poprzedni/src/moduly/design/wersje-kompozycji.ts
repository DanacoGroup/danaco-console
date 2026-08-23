import { Command, type DesignBoardVersion } from '../../../../shared/contract';
import { opisOdmowyBledu } from '../../komponenty/odmowa';
import {
  poleTekstowe,
  poleWyboru,
  przycisk,
  ustawPozycje,
  utworzWierszOdpowiedzi,
} from '../../modele/kontrolki-formularza';
import type { StanKompozycji } from './stan-kompozycji';
import type { StanDesignu } from './stan-designu';

/**
 * Wersjonowanie i wyrys kompozycji w Design Board —
 * `design.board.version.save`, `design.board.version.list`,
 * `design.board.version.restore`, `design.board.export`.
 *
 * ── Po co wersje ────────────────────────────────────────────────────────────
 * `design.board.update` zapisuje układ BIEŻĄCY i zastępuje poprzedni, więc
 * ciągu postaci tablicy nie było skąd wziąć: praca sprzed godziny znikała przy
 * pierwszym przesunięciu warstwy.
 *
 * ── Przywrócenie nie kasuje stanu porzuconego ───────────────────────────────
 * Rdzeń zakłada przy przywróceniu wersję z układu sprzed przywrócenia i oddaje
 * ją w odpowiedzi. Panel pokazuje jej identyfikator wprost, żeby droga powrotna
 * była widoczna od razu, a nie po ponownym odczycie wykazu.
 *
 * ── Wyrys jest plikiem, nie zapowiedzią ─────────────────────────────────────
 * Odpowiedź melduje nazwę pliku i jego typ treści. Kompozycja bez ani jednej
 * warstwy z bajtami kończy się odmową rdzenia i panel pokazuje to zdanie
 * w całości — pusty prostokąt podany jako plik wyglądałby jak plik uszkodzony.
 */
export interface WersjeKompozycji {
  element: HTMLElement;
  /** Zleca odczyt wersji kompozycji wskazanej na kanwie. */
  wczytaj(): Promise<void>;
}

/** Formaty wyrysu, które rdzeń składa biblioteką wkompilowaną. */
const FORMATY_WYRYSU = [
  { wartosc: 'png', etykieta: 'PNG — warstwy złożone na jedno płótno' },
  { wartosc: 'pdf', etykieta: 'PDF — wyrys osadzony w dokumencie' },
  { wartosc: 'svg', etykieta: 'SVG — warstwy z treścią osadzoną' },
] as const;

export function utworzWersjeKompozycji(
  stan: StanDesignu,
  kompozycja: StanKompozycji,
): WersjeKompozycji {
  let wersje: readonly DesignBoardVersion[] = [];

  const nazwa = poleTekstowe({
    etykieta: 'Nazwa wersji',
    podpowiedz: 'np. układ przed przebudową nagłówka',
    opis: 'Pole name żądania design.board.version.save; puste bierze znacznik czasu.',
  });
  const uzasadnienie = poleTekstowe({
    etykieta: 'Uzasadnienie wersji',
    podpowiedz: 'dlaczego ten stan wart jest zapamiętania',
    opis: 'Pole note; puste zostawia wersję bez uzasadnienia.',
  });
  const wybor = poleWyboru(
    { etykieta: 'Wersja', opis: 'Pole versionId żądania design.board.version.restore.' },
    [],
  );
  const format = poleWyboru(
    { etykieta: 'Format wyrysu', opis: 'Pole format żądania design.board.export.' },
    [...FORMATY_WYRYSU],
  );
  const skala = poleTekstowe({
    etykieta: 'Krotność skali wyrysu',
    podpowiedz: 'np. 2',
    opis: 'Pole scale; puste bierze skalę naturalną.',
  });

  const zapisz = przycisk('Zapisz wersję układu', 'dn-btn dn-btn--sm dn-btn--atrament');
  const odczytaj = przycisk('Odczytaj wersje', 'dn-btn dn-btn--sm dn-btn--zarys');
  const przywroc = przycisk('Przywróć wersję', 'dn-btn dn-btn--sm dn-btn--zarys');
  const wyrysuj = przycisk('Wyrysuj kompozycję do pliku', 'dn-btn dn-btn--sm dn-btn--zarys');
  const odpowiedz = utworzWierszOdpowiedzi();

  const pasek = document.createElement('div');
  pasek.className = 'md-czynnosci__pasek';
  pasek.append(zapisz, odczytaj, przywroc, wyrysuj);

  const element = document.createElement('div');
  element.className = 'md-wersje';
  element.append(
    nazwa.element,
    uzasadnienie.element,
    wybor.element,
    format.element,
    skala.element,
    pasek,
    odpowiedz.element,
  );

  zapisz.addEventListener('click', () => void zapiszWersje());
  odczytaj.addEventListener('click', () => void wczytaj());
  przywroc.addEventListener('click', () => void przywrocWersje());
  wyrysuj.addEventListener('click', () => void wyrysujKompozycje());

  /**
   * Sprawdza, czy kompozycja jest już w rdzeniu.
   *
   * Wszystkie cztery komendy tego panelu wskazują kompozycję identyfikatorem,
   * który nadaje rdzeń przy pierwszym `design.board.update`. Kanwa nigdy
   * niezapisana nie ma czego wersjonować ani wyrysowywać, a odmowa rdzenia
   * mówiłaby wtedy o kompozycji nieznanej zamiast o zapisie, którego zabrakło.
   */
  function bezKompozycji(komenda: string): boolean {
    if (kompozycja.idKompozycji() !== '') return false;
    odpowiedz.pokaz(
      `Komenda ${komenda} wskazuje kompozycję identyfikatorem nadanym przez rdzeń — ` +
        'najpierw zapisz kompozycję przyciskiem „Zapisz kompozycję w rdzeniu".',
      false,
    );
    return true;
  }

  function pokazWersje(): void {
    ustawPozycje(
      wybor.kontrolka,
      wersje.map((wersja) => ({
        wartosc: wersja.id,
        etykieta: `${wersja.name ?? wersja.id} — warstw: ${wersja.layerCount}`,
      })),
    );
  }

  /** Krotność z pola; zero znaczy „nie wskazano" i pole do rdzenia nie idzie. */
  function krotnosc(): number {
    const liczba = Number.parseFloat(skala.kontrolka.value.trim());
    return Number.isFinite(liczba) && liczba > 0 ? liczba : 0;
  }

  async function zapiszWersje(): Promise<void> {
    if (bezKompozycji(Command.DesignBoardVersionSave)) return;
    odpowiedz.pokaz('Zapis wersji układu…', true);
    const wynik = await stan.czuwanie.prowadz(
      'zapis wersji kompozycji',
      stan.zrodlo.zapiszWersje({
        idKompozycji: kompozycja.idKompozycji(),
        nazwa: nazwa.kontrolka.value,
        uzasadnienie: uzasadnienie.kontrolka.value,
      }),
      {
        wToku: (zdanie) => odpowiedz.pokaz(zdanie, true),
        cisza: (zdanie) => odpowiedz.pokaz(zdanie, false),
        powrot: (zdanie) => odpowiedz.pokaz(zdanie, false),
        spozniona: (zdanie) => odpowiedz.pokaz(zdanie, false),
      },
    );
    if (wynik === null) return;
    if (!wynik.udany || wynik.wynik === undefined) {
      odpowiedz.pokaz(opisOdmowyBledu('Zapis wersji kompozycji', wynik.blad), false);
      return;
    }
    const wersja = wynik.wynik.version;
    wersje = [wersja, ...wersje];
    pokazWersje();
    wybor.kontrolka.value = wersja.id;
    // Liczba warstw z odpowiedzi, nie z kanwy: wersja opisuje układ ZAPISANY
    // w rdzeniu, a kanwa mogła zmienić się po ostatnim zapisie kompozycji.
    odpowiedz.pokaz(
      `Rdzeń utrwalił wersję ${wersja.id} — warstw w migawce: ${wersja.layerCount}.`,
      true,
    );
  }

  async function wczytaj(): Promise<void> {
    if (bezKompozycji(Command.DesignBoardVersionList)) return;
    odpowiedz.pokaz('Odczyt wersji kompozycji…', true);
    const wynik = await stan.czuwanie.prowadz(
      'odczyt wersji kompozycji',
      stan.zrodlo.wersje(kompozycja.idKompozycji(), 0),
      {
        wToku: (zdanie) => odpowiedz.pokaz(zdanie, true),
        cisza: (zdanie) => odpowiedz.pokaz(zdanie, false),
        powrot: (zdanie) => odpowiedz.pokaz(zdanie, false),
        spozniona: (zdanie) => odpowiedz.pokaz(zdanie, false),
      },
    );
    if (wynik === null) return;
    if (!wynik.udany || wynik.wynik === undefined) {
      odpowiedz.pokaz(opisOdmowyBledu('Odczyt wersji kompozycji', wynik.blad), false);
      return;
    }
    wersje = wynik.wynik.versions;
    pokazWersje();
    odpowiedz.pokaz(
      wersje.length === 0
        ? 'Ta kompozycja nie ma jeszcze ani jednej zapisanej wersji.'
        : `Wersji kompozycji: ${wynik.wynik.total}.`,
      true,
    );
  }

  async function przywrocWersje(): Promise<void> {
    if (wybor.kontrolka.value === '') {
      odpowiedz.pokaz('Wskaż wersję — najpierw odczytaj wersje albo zapisz pierwszą.', false);
      return;
    }
    odpowiedz.pokaz('Przywracanie układu z wersji…', true);
    const wynik = await stan.czuwanie.prowadz(
      'przywrócenie wersji kompozycji',
      stan.zrodlo.przywrocWersje(wybor.kontrolka.value),
      {
        wToku: (zdanie) => odpowiedz.pokaz(zdanie, true),
        cisza: (zdanie) => odpowiedz.pokaz(zdanie, false),
        powrot: (zdanie) => odpowiedz.pokaz(zdanie, false),
        spozniona: (zdanie) => odpowiedz.pokaz(zdanie, false),
      },
    );
    if (wynik === null) return;
    if (!wynik.udany || wynik.wynik === undefined) {
      odpowiedz.pokaz(opisOdmowyBledu('Przywrócenie wersji kompozycji', wynik.blad), false);
      return;
    }
    const plansza = wynik.wynik.board;
    // Kanwa bierze układ z odpowiedzi rdzenia, nie odtwarza go z migawki
    // trzymanej w oknie: prawdą o kompozycji jest to, co zapisała baza.
    kompozycja.wczytaj(plansza.id, plansza.name ?? '', plansza.layers ?? []);
    const odlozona = wynik.wynik.supersededVersion;
    odpowiedz.pokaz(
      `Kanwa niesie układ z wersji ${wybor.kontrolka.value} — warstw: ` +
        `${(plansza.layers ?? []).length}.` +
        (odlozona === undefined
          ? ''
          : ` Układ sprzed przywrócenia został odłożony jako wersja ${odlozona.id}.`),
      true,
    );
  }

  async function wyrysujKompozycje(): Promise<void> {
    if (bezKompozycji(Command.DesignBoardExport)) return;
    odpowiedz.pokaz('Wyrys kompozycji…', true);
    const wynik = await stan.czuwanie.prowadz(
      'wyrys kompozycji',
      stan.zrodlo.wyrysujKompozycje({
        idKompozycji: kompozycja.idKompozycji(),
        format: format.kontrolka.value,
        skala: krotnosc(),
        // Obszaru okno nie narzuca: kadr wskazuje się zaznaczeniem na kanwie,
        // a wyrys bez wskazania bierze całość — tak stanowi kontrakt.
        obszar: null,
      }),
      {
        wToku: (zdanie) => odpowiedz.pokaz(zdanie, true),
        cisza: (zdanie) => odpowiedz.pokaz(zdanie, false),
        powrot: (zdanie) => odpowiedz.pokaz(zdanie, false),
        spozniona: (zdanie) => odpowiedz.pokaz(zdanie, false),
      },
    );
    if (wynik === null) return;
    if (!wynik.udany || wynik.wynik === undefined) {
      odpowiedz.pokaz(opisOdmowyBledu('Wyrys kompozycji', wynik.blad), false);
      return;
    }
    const wyrys = wynik.wynik;
    odpowiedz.pokaz(
      `Rdzeń wyrysował plik „${wyrys.fileName}" typu ${wyrys.mediaType} — ` +
        `${wyrys.contentBase64.length} znaków treści w zapisie base64.`,
      true,
    );
  }

  return { element, wczytaj };
}
