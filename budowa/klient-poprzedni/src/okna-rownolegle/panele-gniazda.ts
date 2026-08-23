import {
  paneleNieotwieralne,
  paneleOtwieralne,
} from '../okna-pomocnicze/panele-otwieralne';
import { wytworniaPanelu } from '../okna-pomocnicze/wytwornia-paneli';
import { utworzMagistrale, type Odsubskrybuj } from '../polaczenie/magistrala-zdarzen';
import type { Kanal } from '../protokol/kanal';
import { ikonaPanelu } from './ikony-paneli';
import { utworzKolumnaPaneli } from './kolumna-paneli';
import type { PozycjaMenu } from './menu-paneli';
import { utworzPanelWStosie, type PanelWStosie } from './panel-w-stosie';
import { SZEROKOSC_MIN_PANELU, SZEROKOSC_PANELU_DOMYSLNA } from './rodzaje-obszaru';
import { utworzStanPaneli } from './stan-paneli';
import { utworzSterowaniePanelami } from './sterowanie-panelami';
import { utworzUchwytSzerokosci } from './uchwyt-szerokosci';
import { utworzWidokPelnoekranowy } from './widok-pelnoekranowy';

/**
 * Panele jednego gniazda — od przycisku w nagłówku po kolumnę obok rozmowy.
 *
 * Panel jest własnością rozmowy, a nie ekranu: przy dwóch i czterech oknach
 * każde okno ma własny zestaw paneli. Dlatego powstaje jeden egzemplarz na
 * gniazdo i nie ma tu ani jednej zmiennej na poziomie modułu — inaczej dwa
 * gniazda podzieliłyby się stanem.
 *
 * Panele stoją osobno od `gniazdo-okna.ts`, bo gniazdo odpowiada za jedno
 * miejsce na scenie (opis, nagłówek, rolę, fasadę okna), a panele to siedem
 * współpracujących bytów: stan, wytwórnia, stos, uchwyt, pełny ekran, menu
 * i skróty.
 *
 * Wykaz pozycji przychodzi z rejestru okien pomocniczych po kodzie modułu
 * gniazda, ale do menu wchodzą tylko pozycje mające wytwórnię
 * (`panele-otwieralne.ts`). Pozycja bez wytwórni nie dostaje wiersza
 * nieczynnego — zamiast tego zdanie pod wykazem mówi, ilu pozycji brakuje
 * i dlaczego. Wiersz, w który nie da się kliknąć, byłby bramką.
 *
 * Każdy panel stoi na komendzie rdzenia, więc bez kanału żadna wytwórnia nie
 * ma czym zbudować okna. Wykaz jest wtedy pusty, a menu mówi wprost, czego
 * brakuje (`BRAK_KANALU`) — to brak środka, nie odmowa.
 *
 * Plik nie liczy szerokości (`szerokosci-gniazda.ts`), nie rysuje panelu
 * (`panel-w-stosie.ts`) i nie wie, gdzie na scenie stoi jego gniazdo. Oddaje
 * elementy i mówi, kiedy układ się zmienił — resztę robi gniazdo.
 */
export interface PaneleGniazda {
  /** Kolumna paneli osadzana obok kolumny rozmowy. */
  kolumna: HTMLElement;
  /** Uchwyt ręcznego ustawiania szerokości; chowa się przy pustym stosie. */
  uchwyt: HTMLElement;
  /** Widok pełnoekranowy panelu — kładziony nad kolumnami gniazda. */
  pelnyEkran: HTMLElement;
  /** Sterowanie osadzane w nagłówku okna rozmowy: skróty i menu `⋮`. */
  sterowanie: HTMLElement;
  /** Żądana szerokość kolumny paneli w pikselach; `0` = stos pusty. */
  szerokosc(): number;
  /** Przestawia moduł gniazda — zmienia wykaz pozycji w menu. */
  ustawModul(kod: string): void;
  /** Podaje okno wykonania; pusty napis = rdzeń nie dał gniazdu okna. */
  ustawOkno(kod: string): void;
  /** Zgłasza każdą zmianę układu paneli — gniazdo przelicza kolumny. */
  naZmiane(sluchacz: () => void): Odsubskrybuj;
  /** Zamyka panele wraz z ich subskrypcjami oraz nasłuchy menu. */
  zamknij(): void;
}

export interface OpcjePaneliGniazda {
  /** Kanał do rdzenia; `null` znaczy „nie ma czym otworzyć ani jednego panelu". */
  kanal: Kanal | null;
  /** Kod modułu gniazda — po nim idzie spis pozycji. */
  modul: string;
  /** Okno wykonania gniazda; pusty napis = rdzeń go nie dał. */
  okno: string;
  /** Przedrostek klas modułu przekazywany panelom. */
  przedrostek: string;
  /** Szerokość gniazda w chwili pytania — sufit uchwytu. */
  dostepnaSzerokosc(): number;
  /**
   * Sekcje doklejane w menu `⋮` za kreską — dziś czynności sesji.
   *
   * Przechodzą tędy nietknięte do sterowania panelami. Panele gniazda nie
   * budują ich i nie zaglądają do środka: czynność sesji nie jest panelem,
   * a ten plik odpowiada wyłącznie za panele.
   */
  sekcjeDalsze?: readonly HTMLElement[];
}

export function utworzPaneleGniazda(opcje: OpcjePaneliGniazda): PaneleGniazda {
  let modul = opcje.modul;
  let okno = opcje.okno;

  const zmiany = utworzMagistrale<void>();
  const stan = utworzStanPaneli(SZEROKOSC_PANELU_DOMYSLNA);
  const kolumna = utworzKolumnaPaneli();
  const pelny = utworzWidokPelnoekranowy(() => stan.ustawPelnyEkran(null));

  /** Panele powołane do życia, po kodzie pozycji. */
  const zbudowane = new Map<string, PanelWStosie>();

  const uchwyt = utworzUchwytSzerokosci({
    etykieta: 'Szerokość kolumny paneli tej rozmowy',
    wartosc: SZEROKOSC_PANELU_DOMYSLNA,
    minimum: SZEROKOSC_MIN_PANELU,
    maksimum: () => opcje.dostepnaSzerokosc(),
    naZmiane: (px) => stan.ustawSzerokosc(px),
  });

  const sterowanie = utworzSterowaniePanelami({
    pozycje: pozycjeMenu(),
    nieotwieralne: nieotwieralnych(),
    czyOtwarty: (kod) => stan.czyOtwarty(kod),
    naWybor: (kod) => stan.przelacz(kod),
    ...(opcje.sekcjeDalsze === undefined ? {} : { sekcjeDalsze: opcje.sekcjeDalsze }),
  });

  /**
   * Pozycje wchodzące do menu i na skróty.
   *
   * Bez kanału wykaz jest pusty, a nie wygaszony: pozycja, której nie ma czym
   * zbudować, nie dostaje wiersza. Żadnej pozycji nie podajemy skrótu
   * klawiszowego, bo dla paneli nie jest zarejestrowany ani jeden skrót —
   * napis obiecujący klawisz, który nic nie robi, byłby atrapą.
   */
  function pozycjeMenu(): readonly PozycjaMenu[] {
    if (opcje.kanal === null) return [];
    return paneleOtwieralne(modul).map((pozycja) => ({
      kod: pozycja.kod,
      nazwa: pozycja.nazwa,
      przeznaczenie: pozycja.przeznaczenie,
      ikona: ikonaPanelu(pozycja.kod),
    }));
  }

  /** Ile pozycji spisu modułu nie da się dziś otworzyć. */
  function nieotwieralnych(): number {
    return paneleNieotwieralne(modul).length;
  }

  /**
   * Nazwa pozycji widziana przez Operatora.
   *
   * Kod pozycji jest nazwą dla kodu, nie dla człowieka: „podglad-bash" nie mówi
   * nic komuś, kto patrzy na nagłówek panelu. Kod zostaje nazwą zapasową tylko
   * wtedy, gdy pozycji nie ma już w spisie — a wtedy jest jedyną prawdą, jaką
   * mamy o tym panelu.
   */
  function nazwaPozycji(kod: string): string {
    return paneleOtwieralne(modul).find((pozycja) => pozycja.kod === kod)?.nazwa ?? kod;
  }

  /** Powołuje panel do życia; `null`, gdy nie ma czym go zbudować. */
  function powolaj(kod: string): PanelWStosie | null {
    const wytwornia = wytworniaPanelu(kod);
    if (wytwornia === null || opcje.kanal === null) return null;

    const nazwa = nazwaPozycji(kod);
    const panel = wytwornia({
      kanal: opcje.kanal,
      okno,
      modul,
      przedrostek: opcje.przedrostek,
    });

    return utworzPanelWStosie({
      kod,
      tytul: nazwa,
      panel,
      naZamkniecie: () => stan.zamknij(kod),
      naPelnyEkran: () => stan.ustawPelnyEkran(kod),
    });
  }

  /**
   * Przeliczenie całego zestawu po każdej zmianie stanu.
   *
   * Kolejność ma znaczenie: najpierw zdejmujemy panele zamknięte (wraz z ich
   * subskrypcjami), potem powołujemy nowe, a dopiero na końcu przenosimy treść
   * na pełny ekran. Odwrotna kolejność zamknęłaby panel, którego treść wisi
   * już na scenie, i zostawiła po nim pusty prostokąt.
   */
  function przelicz(): void {
    const otwarte = stan.otwarte();

    for (const [kod, panel] of [...zbudowane]) {
      if (otwarte.includes(kod)) continue;
      panel.zamknij();
      zbudowane.delete(kod);
    }

    for (const kod of otwarte) {
      if (zbudowane.has(kod)) continue;
      const panel = powolaj(kod);
      // Kod otwarty bez wytwórni znaczy rozjazd spisu z kodem. Cisza
      // udawałaby, że pozycji nie było — a była, tylko nic z niej nie wyszło.
      if (panel === null) {
        console.warn('[panele gniazda] pozycja otwarta nie ma czym stanąć', kod);
        continue;
      }
      zbudowane.set(kod, panel);
      panel.odswiez();
    }

    const wStosie = otwarte
      .map((kod) => zbudowane.get(kod))
      .filter((panel): panel is PanelWStosie => panel !== undefined);
    kolumna.ustaw(wStosie);

    const naScenie = stan.pelnyEkran();
    const wyniesiony = naScenie === null ? undefined : zbudowane.get(naScenie);
    if (wyniesiony === undefined) pelny.ukryj();
    else pelny.pokaz(nazwaPozycji(wyniesiony.kod), wyniesiony.tresc);

    uchwyt.ustawWartosc(stan.szerokosc());
    uchwyt.element.hidden = kolumna.pusta();
    sterowanie.odswiez();
    zmiany.oglos();
  }

  stan.naZmiane(przelicz);
  przelicz();

  return {
    kolumna: kolumna.element,
    uchwyt: uchwyt.element,
    pelnyEkran: pelny.element,
    sterowanie: sterowanie.element,

    // Stos pusty nie zajmuje szerokości. Kolumna o zerowej szerokości i tak by
    // się nie pokazała, ale kontrakt szerokości liczy się z tej liczby —
    // niezerowa przy pustym stosie odebrałaby rozmowie miejsce bez powodu.
    szerokosc: () => (kolumna.pusta() ? 0 : stan.szerokosc()),

    ustawModul(kod) {
      if (kod === modul) return;
      modul = kod;
      sterowanie.ustawPozycje(pozycjeMenu(), nieotwieralnych());
      // Panele otwarte zostają: moduł zmienia spis pozycji, a nie to, co
      // Operator już postawił obok rozmowy. Zdjęcie panelu jest jego decyzją.
      przelicz();
    },

    ustawOkno(kod) {
      if (kod === okno) return;
      okno = kod;
      // Panele już stojące wzięły okno przy powołaniu i nie umieją go zmienić
      // w locie — powstają na nowo, żeby pytały rdzeń o okno właściwe.
      for (const [, panel] of zbudowane) panel.zamknij();
      zbudowane.clear();
      przelicz();
    },

    naZmiane: (sluchacz) => zmiany.subskrybuj(sluchacz),

    zamknij() {
      for (const [, panel] of zbudowane) panel.zamknij();
      zbudowane.clear();
      sterowanie.zamknij();
    },
  };
}
