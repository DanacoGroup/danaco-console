import {
  StudioPageOrientation,
  type StudioComment,
  type StudioDiffHunk,
  type StudioTrackedChange,
} from '../../../../shared/contract';
import {
  marginesyKartki,
  naPunkty,
  NASTAWY_MARGINESOW,
  nosnikPoOznaczeniu,
  nosnikWlasny,
  opiszStrone,
  zastosujNastaweMarginesow,
  rozdzielNaStrony,
  szerokoscKartkiMm,
  szerokoscPolaMm,
  szerokoscPolaPunkty,
  wysokoscKartkiMm,
  wysokoscPolaPunkty,
  zPunktow,
  type StronaPracy,
} from './nastawy-strony';
import { utworzLinijkePozioma, type LinijkaPozioma } from './linijka-pozioma';
import { utworzLinijkePionowa, type LinijkaPionowa } from './linijka-pionowa';
import {
  zerowaWciecieAkapitu,
  type TabulatorAkapitu,
  type WciecieAkapitu,
} from './linijka-podzialka';
import {
  utworzPamiecNastawWidoku,
  type MagazynNastawWidoku,
  type NastawyOperatoraWidoku,
  type PamiecNastawWidoku,
} from './widok-nastawy-operatora';
import {
  policzSkaleWidoku,
  przytnijSkale,
  skalaDokumentu,
  zapamietajSkaleDokumentu,
} from './widok-skali';
import {
  kartkaPrzyPrzewinieciu,
  kolumnyUkladu,
  przewiniecieDoKartki,
  rozlozKartkiWRzedy,
  opiszUkladKartek,
  stronaRozkladowki,
} from './widok-ukladu-stron';
import { utworzPasekWidoku, type PasekWidoku } from './widok-pasek-widoku';
import { czytajNastawyDruku, wydrukuj, type NastawyDruku } from './widok-druku';
import {
  krotnoscStopnia,
  type NastawyAkapitow,
  type NastawyWizualne,
} from './nastawy-wizualne';
import {
  czytajBloki,
  dlugoscDoPunktu,
  serializujPowierzchnie,
  wyrysBloku,
  zapiszBlok,
  type BlokTresci,
  type RodzajBloku,
} from './zapis-formatowany';
import { rozdzielNaOdcinki, opiszZmiane } from './zmiany-modelu';

/**
 * Powierzchnia dokumentu — kartka, na której Operator i model piszą to samo
 * pismo.
 *
 * ── Jedna powierzchnia, cztery tryby widoku ─────────────────────────────────
 * Formatowany (domyślny), źródłowy ze znacznikami, podgląd wydania i różnica
 * nałożona na treść. Tryb zmienia to, CO widać, a nie to, na czym się pracuje:
 * treść jest jedna i mieszka w stanie modułu, a nie w kontrolce.
 *
 * ── Dwa różne rachunki położenia i po co dwa ────────────────────────────────
 * Zaznaczenie jedzie do rdzenia jako zakres ZNAKÓW TREŚCI (`selectionStart`,
 * `selectionEnd` komendy `studio.contextual.op`), więc liczy się je zapisem
 * markdown (`dlugoscDoPunktu`). Kursor po przerysowaniu wraca natomiast na
 * miejsce w TEKŚCIE WIDOCZNYM bloku, bo tam Operator go widział. Dwie miary,
 * dwa zastosowania — zlanie ich przesuwałoby kursor o długość znaczników.
 *
 * ── Paginacja ───────────────────────────────────────────────────────────────
 * Bloki mierzy przeglądarka, rozdziela je `rozdzielNaStrony`. Podział jest
 * prawdziwy: blok, który nie mieści się na kartce, schodzi na następną. Rdzeń
 * liczy własną paginację w `studio.preview.render` i oddaje ją obrazami stron —
 * podgląd wydania sięga po nią osobno, bo tylko rdzeń zna swój silnik składu.
 */

/** Tryb widoku powierzchni. */
export type TrybWidoku = 'formatowany' | 'zrodlowy' | 'wydanie' | 'roznica';

/**
 * Tryb pokazywania adiustacji — trzy, wzorem pakietu biurowego.
 *
 * `cala` znaczy zmiany widoczne w miejscu wraz z decyzją; `po-zmianach` —
 * treść taka, jaka będzie po przyjęciu wszystkiego, bez znaczników;
 * `okienko` — znaczniki w treści zwięzłe, a wykaz zmian obok, w panelu okna.
 * To są tryby JEDNEGO okna, nie trzy okna.
 */
export type TrybAdiustacji = 'cala' | 'po-zmianach' | 'okienko';

/** Zakres zaznaczenia w znakach treści. */
export interface ZakresZaznaczenia {
  poczatek: number;
  koniec: number;
}

/**
 * Jedna kartka wyrysowana przez rdzeń — numer strony i obraz w zapisie `data:`.
 *
 * Zapis `data:` jest tu jedyną drogą: pole `uri` zasobu magazynu jest ścieżką
 * w systemie plików RDZENIA, więc przeglądarka nie wczyta spod niego niczego —
 * także wtedy, gdy zasób powstał bez zarzutu. Bajty przychodzą odpowiedzią
 * komendy `design.asset.content.get` i stąd `data:`.
 */
export interface KartkaRdzenia {
  numer: number;
  zrodlo: string;
}

/** Czynności, którymi powierzchnia rozmawia z oknem. */
export interface CzynnosciPowierzchni {
  /** Treść zmieniona pisaniem w powierzchni. */
  naTresc(tresc: string): void;
  /** Zaznaczenie albo położenie kursora zmienione w powierzchni. */
  naZaznaczenie(zakres: ZakresZaznaczenia | null): void;
  /** Decyzja o jednej zmianie śledzonej podjęta znacznikiem w treści. */
  naDecyzjeZmiany(kod: string, przyjmij: boolean): void;
  /**
   * Nastawy strony przestawione WEWNĄTRZ powierzchni — chwytem na linijce albo
   * paskiem widoku.
   *
   * Pole nieobowiązkowe, bo powierzchnia umie te nastawy prowadzić sama: chwyt
   * marginesu ma działać także wtedy, gdy okno o nim nie wie. Okno, które pole
   * podaje, dostaje nastawy do swojego profilu wydania; okno, które go nie podaje,
   * zachowuje się jak dotąd i niczego nie traci.
   */
  naStrone?(strona: StronaPracy): void;
  /** Wcięcia akapitu przestawione chwytem na linijce. */
  naWciecieAkapitu?(numerBloku: number, wciecie: WciecieAkapitu): void;
  /** Tabulatory akapitu założone albo zdjęte na linijce. */
  naTabulatoryAkapitu?(numerBloku: number, tabulatory: readonly TabulatorAkapitu[]): void;
}

/** Powierzchnia dokumentu wraz z jej sterowaniem. */
export interface PowierzchniaDokumentu {
  /** Element osadzany w oknie — obszar przewijany z kartkami. */
  element: HTMLElement;
  /**
   * Przerysowuje kartki z treści; kursor wraca na swoje miejsce.
   *
   * Treść ta sama co ostatnio wyrysowana NIE jest rysowana drugi raz. Powód jest
   * praktyczny i wprost o pracę Operatora: stan modułu ogłasza zmianę przy każdym
   * naciśnięciu klawisza, a przerysowanie zabiera kursor. `wymus` przełamuje ten
   * warunek tam, gdzie zmieniło się coś INNEGO niż treść — zmiany śledzone,
   * komentarze, nastawy strony, tryb widoku.
   */
  pokaz(tresc: string, wymus?: boolean): void;
  /** Przestawia tryb widoku. */
  ustawTryb(tryb: TrybWidoku): void;
  tryb(): TrybWidoku;
  /**
   * Podaje kartki wyrysowane przez rdzeń — obrazy stron podglądu wydania.
   *
   * Wykaz pusty ZDEJMUJE wyrys rdzenia i wraca do kartek liczonych w oknie. Jest
   * to droga potrzebna, a nie ozdobna: treść zmieniona po renderze czyni obrazy
   * rdzenia nieaktualnymi, a pokazywanie starych kartek jako podglądu bieżącej
   * treści byłoby kłamstwem o dokumencie.
   */
  ustawKartkiRdzenia(kartki: readonly KartkaRdzenia[]): void;
  /** Liczba kartek oddanych przez rdzeń; zero znaczy „rdzeń stron nie oddał". */
  kartekZRdzenia(): number;
  /** Przestawia nastawy strony i przelicza podział. */
  ustawStrone(strona: StronaPracy): void;
  /** Przestawia nastawy wizualne treści. */
  ustawNastawy(nastawy: NastawyWizualne): void;
  /**
   * Liczba kartek w rzędzie — widok jednej albo wielu stron obok siebie.
   *
   * Jeden znaczy „jedna kartka w rzędzie", więcej — układ „obok siebie". Wybór
   * rozkładówki idzie osobno (`ustawUkladKartek`), bo rozkładówka nie jest liczbą
   * kartek: jest książką, w której strona pierwsza stoi sama po prawej.
   */
  ustawKolumny(kolumny: number): void;
  /** Nastawy widoku Operatora obowiązujące teraz. */
  nastawyWidoku(): NastawyOperatoraWidoku;
  /** Przestawia nastawy widoku Operatora i zapamiętuje je. */
  przestawWidok(zmiana: Partial<NastawyOperatoraWidoku>): void;
  /** Skala widoku w procentach; nastawa gotowa liczy się z pola widoku. */
  ustawSkale(procent: number): void;
  /** Nakłada nastawę gotową skali; nastawa nieznana nie zmienia nic i oddaje `false`. */
  ustawNastaweSkali(kod: string): boolean;
  /** Skala pamiętana przy dokumencie — wołane po wczytaniu dokumentu. */
  wczytajSkaleDokumentu(idDokumentu: string): void;
  /** Przewija do wskazanej kartki i oddaje jej numer po przycięciu do zakresu. */
  skoczDoKartki(numer: number): number;
  /** Numer kartki widocznej teraz. */
  kartkaBiezaca(): number;
  /** Linijki — widoczność jest przełącznikiem Operatora. */
  ustawLinijki(widoczne: boolean): void;
  /** Wcięcia akapitu, w którym stoi kursor. */
  wciecieAkapitu(): WciecieAkapitu;
  /** Tabulatory akapitu, w którym stoi kursor. */
  tabulatoryAkapitu(): readonly TabulatorAkapitu[];
  /** Nastawy strony obowiązujące w powierzchni — wraz z nastawami jej własnymi. */
  strona(): StronaPracy;
  /** Zmiany śledzone nakładane na treść. */
  ustawZmiany(zmiany: readonly StudioTrackedChange[]): void;
  /** Przestawia tryb pokazywania adiustacji. */
  ustawAdiustacje(tryb: TrybAdiustacji): void;
  adiustacja(): TrybAdiustacji;
  /** Komentarze, których kotwice stoją w treści. */
  ustawKomentarze(komentarze: readonly StudioComment[]): void;
  /** Położenie kotwicy komentarza w punktach powierzchni; `null`, gdy jej nie ma. */
  polozenieKotwicy(idKomentarza: string): number | null;
  /**
   * Przenosi kursor do zmiany następnej albo poprzedniej i oddaje jej kod.
   *
   * Kolejność jest kolejnością w treści, nie kolejnością zapisu w bazie:
   * „następna zmiana" znaczy następna od miejsca, w którym Operator stoi.
   */
  skoczDoZmiany(wPrzod: boolean): string | null;
  /** Fragmenty różnicy nakładane w trybie różnicy. */
  ustawFragmenty(fragmenty: readonly StudioDiffHunk[]): void;
  /** Zaznaczenie bieżące w znakach treści; `null` znaczy „kursor bez zaznaczenia". */
  zaznaczenie(): ZakresZaznaczenia | null;
  /** Numer bloku, w którym stoi kursor; `-1`, gdy kursor jest poza treścią. */
  blokKursora(): number;
  /** Prostokąt kursora we współrzędnych powierzchni — miejsce wiersza polecenia. */
  polozenieKursora(): { x: number; y: number } | null;
  /** Przestawia styl nazwany bloku, w którym stoi kursor. */
  ustawStylBloku(rodzaj: RodzajBloku): void;
  /** Bloki treści bieżącej — czyta je wstążka i pasek stanu. */
  bloki(): readonly BlokTresci[];
  /** Liczba stron po podziale. */
  liczbaStron(): number;
  /** Zdanie o kartce i podziale — do paska stanu. */
  opis(): string;
}

/** Ile bloków wolno przerysować bez podziału na strony. */
const GRANICA_PODZIALU = 400;

export function utworzPowierzchnieDokumentu(
  strona: StronaPracy,
  nastawy: NastawyWizualne,
  akapity: NastawyAkapitow,
  czynnosci: CzynnosciPowierzchni,
  /**
   * Magazyn nastaw widoku — pominięty znaczy magazyn tego urządzenia.
   *
   * Podanie magazynu opartego na `studio.view.get` i `.set`
   * (`strona-magazyn-widoku.ts`) przenosi nastawy widoku do rdzenia BEZ zmiany
   * ani jednego wołacza wewnątrz powierzchni — po to od początku brały magazyn
   * podany, a nie sięgały po `localStorage` z globalnej przestrzeni. `null`
   * znaczy „bez zapisu": nastawy działają przez sesję i nikt nie obiecuje, że
   * przeżyją zamknięcie.
   */
  magazynWidoku?: MagazynNastawWidoku | null,
): PowierzchniaDokumentu {
  let trybBiezacy: TrybWidoku = 'formatowany';
  let stronaBiezaca = strona;
  let nastawyBiezace = nastawy;
  let zmianySledzone: readonly StudioTrackedChange[] = [];
  let fragmentyRoznicy: readonly StudioDiffHunk[] = [];
  let blokiBiezace: BlokTresci[] = [];
  let trescBiezaca = '';
  let liczbaStronBiezaca = 1;
  let trybAdiustacji: TrybAdiustacji = 'cala';
  let komentarzeBiezace: readonly StudioComment[] = [];
  let kartkaWidoczna = 1;

  /**
   * Nastawy strony przestawione w powierzchni — chwytem linijki albo paskiem.
   *
   * Trzymane osobno od nastaw, które przyszły z okna, bo okno pcha CAŁE nastawy
   * strony (`ustawStrone`) ze swojej kopii. Bez tego wykazu chwyt marginesu
   * przeżyłby dokładnie do następnego naciśnięcia czegokolwiek na wstążce, a to
   * wyglądałoby na usterkę linijki. Nastawa okna wygrywa wyłącznie z tym polem,
   * które okno naprawdę zmieniło — porównanie z ostatnią kopią okna niżej.
   */
  let nadpisaniaStrony: Partial<StronaPracy> = {};
  let stronaOkna = strona;

  // Rozróżnienie „pominięto" od „podano null" jest tu istotne: `undefined` ma dać
  // magazyn domyślny, a `null` ma znaczyć „bez zapisu". Wywołanie z jednym
  // argumentem załatwia oba przypadki, bo `utworzPamiecNastawWidoku` ma własną
  // wartość domyślną, a `null` przechodzi przez nią nietknięty.
  const pamiecWidoku: PamiecNastawWidoku =
    magazynWidoku === undefined
      ? utworzPamiecNastawWidoku()
      : utworzPamiecNastawWidoku(magazynWidoku);
  let widok = pamiecWidoku.biezace();
  /** Dokument, przy którym pamiętana jest skala; puste znaczy „bez dokumentu". */
  let dokumentSkali = '';

  /** Wcięcia akapitów, numerem bloku — tak samo jak nastawy akapitu. */
  const wcieciaBlokow = new Map<number, WciecieAkapitu>();
  /** Tabulatory akapitów, numerem bloku. */
  const tabulatoryBlokow = new Map<number, readonly TabulatorAkapitu[]>();
  /** Szerokości kolumn tabel w milimetrach, numerem bloku tabeli. */
  const kolumnyTabel = new Map<number, readonly number[]>();

  const pole = document.createElement('div');
  pole.className = 'ms-praca__kartki';

  const zrodlo = document.createElement('textarea');
  zrodlo.className = 'dn-pole-kontrolka ms-praca__zrodlo';
  zrodlo.spellcheck = false;
  zrodlo.rows = 24;
  zrodlo.setAttribute('aria-label', 'Treść dokumentu w zapisie źródłowym ze znacznikami');

  /* ── Linijki i pasek widoku ──────────────────────────────────────────────── */

  const linijkaPozioma: LinijkaPozioma = utworzLinijkePozioma(stronaBiezaca, {
    naMargines: (ktory, milimetry) =>
      przestawStroneWewnatrz(ktory === 'lewy' ? { marginesLewyMm: milimetry } : { marginesPrawyMm: milimetry }),
    naWciecie: (wciecie) => ustawWciecieKursora(wciecie),
    naTabulatory: (tabulatory) => ustawTabulatoryKursora(tabulatory),
    naKrawedzKolumny: (numer, milimetry) => ustawKrawedzKolumny(numer, milimetry),
  });

  const linijkaPionowa: LinijkaPionowa = utworzLinijkePionowa(stronaBiezaca, {
    naMargines: (ktory, milimetry) =>
      przestawStroneWewnatrz(ktory === 'gora' ? { marginesGoraMm: milimetry } : { marginesDolMm: milimetry }),
  });

  const pasek: PasekWidoku = utworzPasekWidoku(widok, {
    naSkale: (procent) => ustawSkale(procent),
    naNastaweSkali: (kod) => ustawNastaweSkali(kod),
    naUklad: (uklad, wRzedzie) => przestawWidok({ ukladKartek: uklad, kartekWRzedzie: wRzedzie }),
    naPrzewijanie: (tryb) => przestawWidok({ przewijanie: tryb }),
    naLinijki: (widoczne) => przestawWidok({ linijkiWidoczne: widoczne }),
    naJednostke: (jednostka) => przestawWidok({ jednostka }),
    naGraniceMarginesow: (widoczne) => przestawWidok({ graniceMarginesow: widoczne }),
    naTrybZrodlowy: (wlaczony) => ustawTrybWidoku(wlaczony ? 'zrodlowy' : 'formatowany'),
    naPodgladWydruku: (wlaczony) => ustawTrybWidoku(wlaczony ? 'wydanie' : 'formatowany'),
    naKartke: (numer) => void skoczDoKartki(numer),
    naTrybDokumentow: (tryb) => przestawWidok({ trybDokumentow: tryb }),
    naKierunekPodzialu: (kierunek) => przestawWidok({ kierunekPodzialu: kierunek }),
    naUkladPaneli: (uklad) => ustawUkladPaneli(uklad),

    naNosnik: (oznaczenie) => {
      const wybrany = nosnikPoOznaczeniu(oznaczenie);
      // Nośnik nieznany nie zmienia kartki: rysowanie jej w rozmiarze zgadniętym
      // byłoby kłamstwem o nośniku, na którym pismo wyjdzie z drukarki.
      if (wybrany === null) return;
      przestawStroneWewnatrz({ nosnik: wybrany });
    },

    naNosnikWlasny: (szerokosc, wysokosc) => {
      const wlasny = nosnikWlasny(szerokosc, wysokosc);
      if (wlasny === null) return;
      przestawStroneWewnatrz({ nosnik: wlasny });
    },

    naOrientacje: (pozioma) =>
      przestawStroneWewnatrz({
        orientacja: pozioma ? StudioPageOrientation.Pozioma : StudioPageOrientation.Pionowa,
      }),

    naNastaweMarginesow: (kod) => {
      const nastawa = NASTAWY_MARGINESOW.find((pozycja) => pozycja.kod === kod);
      if (nastawa === undefined) return;
      const zlozona = zastosujNastaweMarginesow(stronaBiezaca, nastawa);
      przestawStroneWewnatrz({
        marginesGoraMm: zlozona.marginesGoraMm,
        marginesDolMm: zlozona.marginesDolMm,
        marginesLewyMm: zlozona.marginesLewyMm,
        marginesPrawyMm: zlozona.marginesPrawyMm,
        marginesOprawyMm: zlozona.marginesOprawyMm,
        marginesyOdbicia: zlozona.marginesyOdbicia,
      });
    },

    naMargines: (ktory, milimetry) => {
      if (ktory === 'gora') przestawStroneWewnatrz({ marginesGoraMm: milimetry });
      if (ktory === 'dol') przestawStroneWewnatrz({ marginesDolMm: milimetry });
      if (ktory === 'lewy') przestawStroneWewnatrz({ marginesLewyMm: milimetry });
      if (ktory === 'prawy') przestawStroneWewnatrz({ marginesPrawyMm: milimetry });
    },

    naOprawe: (milimetry) => przestawStroneWewnatrz({ marginesOprawyMm: milimetry }),
    naStroneOprawy: (strona) => przestawStroneWewnatrz({ stronaOprawy: strona }),
    naOdbicia: (wlaczone) => przestawStroneWewnatrz({ marginesyOdbicia: wlaczone }),
    naDruk: (nastawy) => void drukuj(nastawy),
    naSzybkiDruk: () => void drukuj(czytajNastawyDruku()),
  });

  /** Zdanie o ostatnim wydruku — czyta je pasek stanu przez `opis`. */
  let zdanieDruku = '';

  /**
   * Drukuje kartki powierzchni.
   *
   * Wydruk bierze TE kartki, które Operator widzi w podglądzie — nie drugi wyrys
   * i nie surowy tekst. Przygotowanie polega więc na trzech rzeczach: zejściu do
   * trybu wydania (żeby wydruk zgadzał się z podglądem), oznaczeniu stron poza
   * zakresem i zdjęciu adiustacji, gdy Operator drukuje pismo do wysłania. Wszystkie
   * trzy są odwracane po zamknięciu okna drukarki.
   */
  function drukuj(nastawy: NastawyDruku): void {
    const wynik = wydrukuj(nastawy, {
      liczbaStron: () => liczbaStronBiezaca,
      kartkaBiezaca: () => kartkaWidoczna,
      przygotuj: (nastawyDruku, strony) => {
        const trybPrzedDrukiem = trybBiezacy;
        const adiustacjaPrzedDrukiem = trybAdiustacji;
        if (trybPrzedDrukiem !== 'wydanie') ustawTrybWidoku('wydanie');
        if (nastawyDruku.adiustacja === 'po-zmianach' && adiustacjaPrzedDrukiem !== 'po-zmianach') {
          trybAdiustacji = 'po-zmianach';
          element.dataset['adiustacja'] = 'po-zmianach';
          pokaz(trescBiezaca, true);
        }
        const objete = new Set(strony);
        for (const kartka of Array.from(pole.querySelectorAll<HTMLElement>('[data-strona]'))) {
          const numer = Number(kartka.dataset['strona'] ?? '0');
          kartka.dataset['druk'] = objete.has(numer) ? 'tak' : 'pomin';
        }
        element.dataset['druk'] = 'tak';
        element.style.setProperty('--ms-skala-druku', String(nastawyDruku.skala / 100));

        return () => {
          delete element.dataset['druk'];
          element.style.removeProperty('--ms-skala-druku');
          for (const kartka of Array.from(pole.querySelectorAll<HTMLElement>('[data-strona]'))) {
            delete kartka.dataset['druk'];
          }
          if (trybAdiustacji !== adiustacjaPrzedDrukiem) {
            trybAdiustacji = adiustacjaPrzedDrukiem;
            element.dataset['adiustacja'] = adiustacjaPrzedDrukiem;
          }
          if (trybBiezacy !== trybPrzedDrukiem) ustawTrybWidoku(trybPrzedDrukiem);
          else pokaz(trescBiezaca, true);
        };
      },
    });
    zdanieDruku = wynik.zdanie;
  }

  const naroznik = document.createElement('div');
  naroznik.className = 'ms-praca__naroznik';

  const plansza = document.createElement('div');
  plansza.className = 'ms-praca__plansza';
  plansza.append(naroznik, linijkaPozioma.element, linijkaPionowa.element, pole);

  const element = document.createElement('div');
  element.className = 'ms-praca__powierzchnia';
  element.dataset['tryb'] = trybBiezacy;
  element.dataset['adiustacja'] = trybAdiustacji;
  element.dataset['przewijanie'] = widok.przewijanie;
  element.dataset['granice'] = widok.graniceMarginesow ? 'tak' : 'nie';
  element.dataset['uklad'] = widok.ukladKartek;
  element.append(pasek.element, plansza, zrodlo);

  /* ── Nastawy jako zmienne arkusza ────────────────────────────────────────── */

  function przypnijNastawy(): void {
    const styl = element.style;
    styl.setProperty('--ms-kartka-szerokosc', `${naPunkty(szerokoscKartki())}px`);
    styl.setProperty('--ms-kartka-wysokosc', `${naPunkty(wysokoscKartki())}px`);
    styl.setProperty('--ms-margines-gora', `${naPunkty(stronaBiezaca.marginesGoraMm)}px`);
    styl.setProperty('--ms-margines-dol', `${naPunkty(stronaBiezaca.marginesDolMm)}px`);
    styl.setProperty('--ms-margines-lewy', `${naPunkty(stronaBiezaca.marginesLewyMm)}px`);
    styl.setProperty('--ms-margines-prawy', `${naPunkty(stronaBiezaca.marginesPrawyMm)}px`);
    styl.setProperty('--ms-skala', String(stronaBiezaca.skala / 100));
    styl.setProperty('--ms-kroj', nastawyBiezace.krój);
    styl.setProperty('--ms-stopien', `${nastawyBiezace.stopien}pt`);
    styl.setProperty('--ms-interlinia', String(nastawyBiezace.interlinia));
    styl.setProperty('--ms-wciecie', `${naPunkty(nastawyBiezace.wciecieMm)}px`);
    styl.setProperty('--ms-odstep-akapitu', `${naPunkty(nastawyBiezace.odstepMm)}px`);
    styl.setProperty(
      '--ms-kolumny',
      String(kolumnyUkladu(widok.ukladKartek, widok.kartekWRzedzie)),
    );
    element.dataset['przewijanie'] = widok.przewijanie;
    element.dataset['granice'] = stronaBiezaca.graniceMarginesow ? 'tak' : 'nie';
    element.dataset['uklad'] = widok.ukladKartek;
    linijkaPozioma.ustawStrone(stronaBiezaca);
    linijkaPionowa.ustawStrone(stronaBiezaca);
    linijkaPozioma.ustawWidocznosc(widok.linijkiWidoczne && trybBiezacy !== 'zrodlowy');
    linijkaPionowa.ustawWidocznosc(widok.linijkiWidoczne && trybBiezacy !== 'zrodlowy');
    naroznik.hidden = !widok.linijkiWidoczne || trybBiezacy === 'zrodlowy';
  }

  function szerokoscKartki(): number {
    return szerokoscKartkiMm(stronaBiezaca);
  }

  function wysokoscKartki(): number {
    return wysokoscKartkiMm(stronaBiezaca);
  }

  /* ── Nastawy strony i widoku przestawiane wewnątrz powierzchni ───────────── */

  /**
   * Przestawia nastawy strony chwytem linijki albo paskiem widoku.
   *
   * Zmiana idzie do nadpisań, a nie tylko do nastaw bieżących: okno pcha własną
   * kopię nastaw strony przy każdej swojej zmianie i bez nadpisań zabierałoby
   * Operatorowi to, co właśnie ustawił chwytem.
   */
  function przestawStroneWewnatrz(zmiana: Partial<StronaPracy>): void {
    nadpisaniaStrony = { ...nadpisaniaStrony, ...zmiana };
    stronaBiezaca = { ...stronaBiezaca, ...zmiana };
    if (zmiana.skala !== undefined) {
      widok = pamiecWidoku.przestaw({ skala: zmiana.skala });
      if (dokumentSkali !== '') zapamietajSkaleDokumentu(dokumentSkali, zmiana.skala);
    }
    przypnijNastawy();
    czynnosci.naStrone?.(stronaBiezaca);
    pokaz(trescBiezaca, true);
  }

  /** Przestawia nastawy widoku Operatora, zapamiętuje je i przerysowuje. */
  function przestawWidok(zmiana: Partial<NastawyOperatoraWidoku>): void {
    widok = pamiecWidoku.przestaw(zmiana);
    if (zmiana.jednostka !== undefined) {
      przestawStroneWewnatrz({ jednostka: zmiana.jednostka });
      return;
    }
    if (zmiana.graniceMarginesow !== undefined) {
      przestawStroneWewnatrz({ graniceMarginesow: zmiana.graniceMarginesow });
      return;
    }
    przypnijNastawy();
    pokaz(trescBiezaca, true);
  }

  /** Skala widoku w procentach — nastawa strony i nastawa widoku naraz. */
  function ustawSkale(procent: number): void {
    przestawStroneWewnatrz({ skala: przytnijSkale(procent) });
  }

  /**
   * Nakłada nastawę gotową skali.
   *
   * Pole widoku mierzy przeglądarka, więc rachunek dostaje wymiary zmierzone
   * TERAZ. Nastawa nieznana albo pole jeszcze nieosadzone nie zmienia nic i oddaje
   * `false` — skala zgadnięta byłaby gorsza od nastawy niewykonanej.
   */
  function ustawNastaweSkali(kod: string): boolean {
    const prostokat = element.getBoundingClientRect();
    const szerokoscLinijki = widok.linijkiWidoczne ? linijkaPionowa.element.offsetWidth : 0;
    const policzona = policzSkaleWidoku(kod, {
      szerokoscWidokuPx: Math.max(0, prostokat.width - szerokoscLinijki - 32),
      wysokoscWidokuPx: Math.max(0, prostokat.height - pasek.element.offsetHeight - 32),
      szerokoscKartkiMm: szerokoscKartki(),
      wysokoscKartkiMm: wysokoscKartki(),
      szerokoscTekstuMm: szerokoscPolaMm(stronaBiezaca, kartkaWidoczna),
      kartekWRzedzie: kolumnyUkladu(widok.ukladKartek, widok.kartekWRzedzie),
      odstepPx: 24,
    });
    if (policzona === null) return false;
    ustawSkale(policzona);
    return true;
  }

  /** Przestawia tryb widoku wraz z zapisem przełącznika trybu źródłowego. */
  function ustawTrybWidoku(tryb: TrybWidoku): void {
    trybBiezacy = tryb;
    element.dataset['tryb'] = tryb;
    pole.contentEditable = tryb === 'wydanie' ? 'false' : 'true';
    widok = pamiecWidoku.przestaw({ trybZrodlowy: tryb === 'zrodlowy' });
    przypnijNastawy();
    pokaz(trescBiezaca, true);
  }

  /**
   * Wcięcia akapitu, w którym stoi kursor.
   *
   * Bez kursora w treści nie ma akapitu, do którego wcięcie by należało — więc
   * linijka pokazuje wcięcia zerowe, a nie wcięcia bloku poprzedniego.
   */
  function wciecieAkapituKursora(): WciecieAkapitu {
    const numer = numerBlokuKursora();
    return (numer < 0 ? undefined : wcieciaBlokow.get(numer)) ?? zerowaWciecieAkapitu();
  }

  function numerBlokuKursora(): number {
    const miejsce = zapamietajKursor();
    return miejsce === null ? -1 : miejsce.blok;
  }

  function ustawWciecieKursora(wciecie: WciecieAkapitu): void {
    const numer = numerBlokuKursora();
    if (numer < 0) return;
    wcieciaBlokow.set(numer, wciecie);
    czynnosci.naWciecieAkapitu?.(numer, wciecie);
    linijkaPozioma.ustawWciecie(wciecie);
    pokaz(trescBiezaca, true);
  }

  function ustawTabulatoryKursora(tabulatory: readonly TabulatorAkapitu[]): void {
    const numer = numerBlokuKursora();
    if (numer < 0) return;
    tabulatoryBlokow.set(numer, tabulatory);
    czynnosci.naTabulatoryAkapitu?.(numer, tabulatory);
  }

  /**
   * Przestawia szerokość kolumn tabeli, w której stoi kursor.
   *
   * Krawędzie idą w milimetrach od lewej krawędzi pola pisania, więc szerokość
   * kolumny jest różnicą krawędzi sąsiednich. Nastawa jest nastawą WIDOKU: składnia
   * tabeli w treści szerokości kolumn nie niesie i kontrakt nie ma jej gdzie
   * zapisać — okno mówi o tym w objaśnieniu chwytu, zamiast udawać zapis.
   */
  function ustawKrawedzKolumny(numer: number, milimetry: number): void {
    const blok = numerBlokuKursora();
    if (blok < 0) return;
    const krawedzie = [...(kolumnyTabel.get(blok) ?? [])];
    if (numer < 0 || numer >= krawedzie.length) return;
    krawedzie[numer] = milimetry;
    kolumnyTabel.set(blok, krawedzie);
    pokaz(trescBiezaca, true);
  }

  /**
   * Panele okna jako nakładki albo jako stałe kolumny.
   *
   * Rozstrzygnięcie Właściciela: powierzchnia należy do dokumentu, więc nakładka
   * jest postacią domyślną, a stałe kolumny zostają trybem do wyboru. Znacznik
   * idzie na wspólnego przodka powierzchni i paneli, bo to on rozkłada kolumny —
   * powierzchnia sięga po niego dopiero teraz, gdy Operator o to prosi, i nie
   * zakłada, że okno w ogóle takiego przodka ma.
   */
  function ustawUkladPaneli(uklad: 'nakladka' | 'kolumny'): void {
    const przodek = element.closest<HTMLElement>('.ms-praca__kolumny');
    if (przodek === null) return;
    przodek.dataset['ukladPaneli'] = uklad;
  }

  /** Przewija do wskazanej kartki i oddaje jej numer po przycięciu. */
  function skoczDoKartki(numer: number): number {
    const docelowa = Math.min(Math.max(1, Math.round(numer)), Math.max(1, liczbaStronBiezaca));
    const kartka = pole.querySelector<HTMLElement>(`[data-strona='${docelowa}']`);
    // Przewinięcie do elementu, gdy powłoka je zna; inaczej rachunkiem z wysokości
    // kartki. Powłoka bez układu (środowisko sprawdzianu) nie ma
    // `scrollIntoView` i wtedy skok liczy się sam, zamiast przerywać czynność.
    if (kartka !== null && typeof kartka.scrollIntoView === 'function') {
      kartka.scrollIntoView({ block: 'start' });
    } else {
      element.scrollTop = przewiniecieDoKartki(
        docelowa,
        naPunkty(wysokoscKartki()) * (stronaBiezaca.skala / 100),
        24,
      );
    }
    kartkaWidoczna = docelowa;
    odswiezPasek();
    return docelowa;
  }

  function odswiezPasek(): void {
    pasek.odswiez(widok, trybBiezacy === 'wydanie', liczbaStronBiezaca, kartkaWidoczna);
  }

  /* ── Kursor: zapamiętanie i powrót ───────────────────────────────────────── */

  /** Miejsce kursora liczone w tekście widocznym bloku. */
  interface MiejsceKursora {
    blok: number;
    znak: number;
  }

  function elementyBlokow(): HTMLElement[] {
    return Array.from(pole.querySelectorAll<HTMLElement>('[data-numer-bloku]'));
  }

  function blokWezla(wezel: Node | null): HTMLElement | null {
    let biezacy: Node | null = wezel;
    while (biezacy !== null && biezacy !== pole) {
      if (biezacy instanceof HTMLElement && biezacy.dataset['numerBloku'] !== undefined) {
        return biezacy;
      }
      biezacy = biezacy.parentNode;
    }
    return null;
  }

  function zapamietajKursor(): MiejsceKursora | null {
    const wybor = document.getSelection();
    if (wybor === null || wybor.rangeCount === 0) return null;
    const zakres = wybor.getRangeAt(0);
    const blok = blokWezla(zakres.startContainer);
    if (blok === null) return null;
    const przed = document.createRange();
    przed.selectNodeContents(blok);
    przed.setEnd(zakres.startContainer, zakres.startOffset);
    return {
      blok: Number(blok.dataset['numerBloku'] ?? '0'),
      znak: przed.toString().length,
    };
  }

  function przywrocKursor(miejsce: MiejsceKursora | null): void {
    if (miejsce === null) return;
    const blok = elementyBlokow().find(
      (element) => Number(element.dataset['numerBloku'] ?? '-1') === miejsce.blok,
    );
    if (blok === undefined) return;
    const chodzik = document.createTreeWalker(blok, NodeFilter.SHOW_TEXT);
    let policzone = 0;
    let ostatni: Text | null = null;
    while (chodzik.nextNode() !== null) {
      const wezel = chodzik.currentNode as Text;
      ostatni = wezel;
      const dlugosc = (wezel.textContent ?? '').length;
      if (policzone + dlugosc >= miejsce.znak) {
        ustawKursorNa(wezel, miejsce.znak - policzone);
        return;
      }
      policzone += dlugosc;
    }
    if (ostatni !== null) ustawKursorNa(ostatni, (ostatni.textContent ?? '').length);
  }

  function ustawKursorNa(wezel: Text, offset: number): void {
    const wybor = document.getSelection();
    if (wybor === null) return;
    const zakres = document.createRange();
    zakres.setStart(wezel, Math.max(0, Math.min(offset, (wezel.textContent ?? '').length)));
    zakres.collapse(true);
    wybor.removeAllRanges();
    wybor.addRange(zakres);
  }

  /* ── Wyrys ───────────────────────────────────────────────────────────────── */

  /**
   * Buduje jedną kartkę wraz z nagłówkiem, polem pisania i stopką.
   *
   * Marginesy liczy `marginesyKartki` z numeru strony, bo przy marginesach odbicia
   * i przy oprawie kartka parzysta ma inne marginesy niż nieparzysta. Bez tego pas
   * oprawy wypadałby raz w rowku, a raz na krawędzi.
   */
  function utworzKartke(numer: number, ile: number): { kartka: HTMLElement; pole: HTMLElement } {
    const kartka = document.createElement('div');
    kartka.className = 'ms-strona';
    kartka.dataset['strona'] = String(numer);
    kartka.dataset['polozenie'] = stronaRozkladowki(numer);
    const marginesy = marginesyKartki(stronaBiezaca, numer);
    kartka.style.paddingTop = `${naPunkty(marginesy.goraMm)}px`;
    kartka.style.paddingBottom = `${naPunkty(marginesy.dolMm)}px`;
    kartka.style.paddingLeft = `${naPunkty(marginesy.lewyMm)}px`;
    kartka.style.paddingRight = `${naPunkty(marginesy.prawyMm)}px`;

    const naglowek = document.createElement('div');
    naglowek.className = 'ms-strona__naglowek';
    naglowek.contentEditable = 'false';
    naglowek.textContent = stronaBiezaca.naglowek;

    const cialo = document.createElement('div');
    cialo.className = 'ms-strona__pole';

    const stopka = document.createElement('div');
    stopka.className = 'ms-strona__stopka';
    stopka.contentEditable = 'false';
    const numeracja = stronaBiezaca.numeracja ? `Strona ${numer} z ${ile}` : '';
    stopka.textContent =
      stronaBiezaca.stopka === '' ? numeracja : `${stronaBiezaca.stopka} · ${numeracja}`;

    kartka.append(naglowek, cialo, stopka);
    return { kartka, pole: cialo };
  }

  /** Nakłada nastawy akapitu i stopień stylu nazwanego na element bloku. */
  function przypnijAkapit(element: HTMLElement, numer: number, rodzaj: RodzajBloku): void {
    const nastawa = akapity.dla(numer);
    element.dataset['numerBloku'] = String(numer);
    element.dataset['wyrownanie'] = nastawa.wyrownanie;
    element.style.setProperty('--ms-krotnosc', String(krotnoscStopnia(rodzaj)));
    if (nastawa.barwa !== '') element.style.color = nastawa.barwa;

    // Wcięcia chwycone na linijce są nastawą TEGO akapitu, więc idą na jego
    // element, a nie na kartkę: wcięcie ustawione dla jednego akapitu nie może
    // przesuwać całego pisma.
    const wciecie = wcieciaBlokow.get(numer);
    if (wciecie !== undefined) {
      element.style.marginLeft = `${naPunkty(wciecie.leweMm)}px`;
      element.style.marginRight = `${naPunkty(wciecie.praweMm)}px`;
      element.style.textIndent = `${naPunkty(wciecie.pierwszyWierszMm)}px`;
    }
    const tabulatory = tabulatoryBlokow.get(numer);
    if (tabulatory !== undefined && tabulatory.length > 0) {
      element.style.tabSize = String(Math.max(1, Math.round(tabulatory[0]?.milimetry ?? 8)));
      element.dataset['tabulatory'] = String(tabulatory.length);
    }
    const krawedzie = kolumnyTabel.get(numer);
    if (krawedzie !== undefined && rodzaj === 'tabela') {
      const wiersz = element.querySelector('tr');
      if (wiersz !== null) {
        let poprzednia = 0;
        Array.from(wiersz.children).forEach((komorka, kolumna) => {
          const krawedz = krawedzie[kolumna];
          if (krawedz === undefined || !(komorka instanceof HTMLElement)) return;
          komorka.style.width = `${naPunkty(krawedz - poprzednia)}px`;
          poprzednia = krawedz;
        });
      }
    }
  }

  /**
   * Nakłada zmiany śledzone na treść bloku.
   *
   * Odcinek objęty zmianą dostaje własny element ze znakiem autora i dwoma
   * przyciskami decyzji. Decyzja stoi w miejscu zmiany, a nie w osobnym oknie —
   * to jest sedno rozstrzygnięcia: Operator rozstrzyga tam, gdzie patrzy.
   */
  function nalozZmiany(element: HTMLElement, blok: BlokTresci, poczatek: number): string[] {
    if (trybAdiustacji === 'po-zmianach') return [];
    const zapis = zapiszBlok(blok);
    const odcinki = rozdzielNaOdcinki(zapis, przesunieteZmiany(poczatek, zapis.length));
    if (odcinki.length <= 1 && odcinki[0]?.zmiana === null) return [];

    const objete: string[] = [];
    element.replaceChildren();
    for (const odcinek of odcinki) {
      if (odcinek.zmiana === null) {
        element.append(document.createTextNode(odcinek.tekst));
        continue;
      }
      objete.push(odcinek.zmiana.id);
      const znacznik = document.createElement('span');
      znacznik.className = 'ms-zmiana';
      znacznik.dataset['zmiana'] = odcinek.zmiana.id;
      znacznik.dataset['autor'] = odcinek.zmiana.author;
      znacznik.title = opiszZmiane(odcinek.zmiana);
      znacznik.append(document.createTextNode(odcinek.tekst));
      // Znaczniki decyzji stoją w treści tylko przy adiustacji całej; w okienku
      // recenzowania decyzja jest w wykazie obok i dwa jej miejsca naraz kazałyby
      // zgadywać, które jest właściwe.
      if (trybAdiustacji === 'cala') {
        znacznik.append(
          przyciskDecyzji(odcinek.zmiana.id, true),
          przyciskDecyzji(odcinek.zmiana.id, false),
        );
      }
      element.append(znacznik);
    }
    return objete;
  }

  /**
   * Pasek zatwierdzenia zmiany modelu — „Gotowe" i „Cofnij" pod zmianą.
   *
   * Stoi pod blokiem, w którym zmiana leży, i zamyka pracę modelu jednym
   * naciśnięciem. „Gotowe" to przyjęcie zmian bloku (`studio.tracking.decide`
   * z `accept`), „Cofnij" to ich odrzucenie — treść wraca do postaci sprzed
   * operacji. Pasek nie jest trzecim miejscem decyzji: woła tę samą drogę, co
   * znaczniki w treści i wykaz w okienku recenzowania.
   */
  function pasekZatwierdzenia(kody: readonly string[]): HTMLElement {
    const zdanie = document.createElement('span');
    zdanie.className = 'ms-zatwierdzenie__zdanie';
    zdanie.textContent =
      kody.length === 1
        ? 'Model naniósł tu zmianę — zamknij ją albo wycofaj.'
        : `Model naniósł tu zmian: ${kody.length} — zamknij je albo wycofaj.`;

    const gotowe = document.createElement('button');
    gotowe.type = 'button';
    gotowe.className = 'dn-btn dn-btn--sm dn-btn--sygnal';
    gotowe.textContent = 'Gotowe';
    gotowe.contentEditable = 'false';
    gotowe.title = 'Przyjmuje zmiany tego miejsca i zakłada wersję w repozytorium sesji.';
    gotowe.addEventListener('click', () => {
      for (const kod of kody) czynnosci.naDecyzjeZmiany(kod, true);
    });

    const cofnij = document.createElement('button');
    cofnij.type = 'button';
    cofnij.className = 'dn-btn dn-btn--sm dn-btn--zarys';
    cofnij.textContent = 'Cofnij';
    cofnij.contentEditable = 'false';
    cofnij.title = 'Odrzuca zmiany tego miejsca — wraca treść sprzed pracy modelu.';
    cofnij.addEventListener('click', () => {
      for (const kod of kody) czynnosci.naDecyzjeZmiany(kod, false);
    });

    const pasek = document.createElement('div');
    pasek.className = 'ms-zatwierdzenie';
    pasek.contentEditable = 'false';
    pasek.dataset['zmiany'] = kody.join(' ');
    pasek.append(zdanie, gotowe, cofnij);
    return pasek;
  }

  /**
   * Kotwice komentarzy — miejsce, w które celuje dymek na marginesie.
   *
   * Kotwica jest znakiem w treści, a nie wpisem w wykazie: komentarz przypięty
   * do fragmentu ma pokazywać, do którego. Komentarz bez zakresu kotwicy nie
   * dostaje — dotyczy dokumentu, nie miejsca.
   */
  function nalozKotwice(element: HTMLElement, numer: number, poczatek: number, dlugosc: number): void {
    for (const komentarz of komentarzeBiezace) {
      const od = komentarz.selectionStart;
      if (od === undefined) continue;
      if (od < poczatek || od > poczatek + dlugosc) continue;
      if (komentarz.parentCommentId !== undefined) continue;
      const kotwica = document.createElement('span');
      kotwica.className = 'ms-kotwica';
      kotwica.contentEditable = 'false';
      kotwica.dataset['komentarz'] = komentarz.id;
      kotwica.dataset['blok'] = String(numer);
      kotwica.textContent = '❝';
      kotwica.title = `Komentarz ${komentarz.author}: ${komentarz.body}`;
      element.append(kotwica);
    }
  }

  /** Zmiany przycięte do zakresu bloku i przeliczone na jego początek. */
  function przesunieteZmiany(poczatek: number, dlugosc: number): StudioTrackedChange[] {
    const wynik: StudioTrackedChange[] = [];
    for (const zmiana of zmianySledzone) {
      const od = zmiana.rangeStart - poczatek;
      const doZnaku = zmiana.rangeEnd - poczatek;
      if (doZnaku < 0 || od > dlugosc) continue;
      wynik.push({
        ...zmiana,
        rangeStart: Math.max(0, od),
        rangeEnd: Math.min(dlugosc, doZnaku),
      });
    }
    return wynik;
  }

  function przyciskDecyzji(kod: string, przyjmij: boolean): HTMLElement {
    const przycisk = document.createElement('button');
    przycisk.type = 'button';
    przycisk.className = 'ms-zmiana__decyzja';
    przycisk.contentEditable = 'false';
    przycisk.dataset['decyzja'] = przyjmij ? 'przyjmij' : 'odrzuc';
    przycisk.textContent = przyjmij ? '✓' : '✕';
    przycisk.title = przyjmij
      ? 'Przyjmij tę zmianę — treść zostaje w postaci po zmianie, rdzeń zakłada wersję.'
      : 'Odrzuć tę zmianę — w to miejsce wraca treść sprzed niej.';
    przycisk.addEventListener('mousedown', (zdarzenie) => zdarzenie.preventDefault());
    przycisk.addEventListener('click', () => czynnosci.naDecyzjeZmiany(kod, przyjmij));
    return przycisk;
  }

  /** Numery bloków objętych fragmentami różnicy, liczone wierszami. */
  function blokiFragmentow(): Map<number, StudioDiffHunk> {
    const objete = new Map<number, StudioDiffHunk>();
    if (fragmentyRoznicy.length === 0) return objete;
    let wiersz = 1;
    blokiBiezace.forEach((blok, numer) => {
      const wysokosc = zapiszBlok(blok).split('\n').length;
      const od = wiersz;
      const doWiersza = wiersz + wysokosc - 1;
      for (const fragment of fragmentyRoznicy) {
        const start = fragment.startLine ?? 0;
        const koniec = fragment.endLine ?? start;
        if (start === 0) continue;
        if (koniec >= od && start <= doWiersza) objete.set(numer, fragment);
      }
      wiersz = doWiersza + 2;
    });
    return objete;
  }

  /** Treść, z której powstał wyrys bieżący; puste znaczy „jeszcze nie rysowano". */
  let trescWyrysowana: string | null = null;

  /**
   * Kartki wyrysowane przez RDZEŃ — obrazy stron podglądu wydania.
   *
   * Wykaz pusty znaczy „rdzeń jeszcze nie oddał stron", nie „stron nie ma".
   * Rozróżnienie jest ważne, bo od niego zależy, co pokazuje pasek stanu: liczbę
   * kartek rdzenia albo liczbę kartek liczonych w oknie.
   */
  let kartkiRdzenia: readonly KartkaRdzenia[] = [];

  /**
   * Rysuje kartki oddane przez rdzeń, w tym samym układzie rzędów, co kartki
   * liczone w oknie.
   *
   * Obraz wchodzi jako `<img>` o źródle `data:` — bajty przyszły odpowiedzią
   * komendy, więc nie ma tu ani jednego żądania do sieci. Szerokość kartki jest
   * szerokością pola, żeby skala widoku i nastawy gotowe („do szerokości strony")
   * działały na wyrysie rdzenia tak samo jak na kartce liczonej w oknie.
   */
  function wyrysujKartkiRdzenia(): void {
    liczbaStronBiezaca = kartkiRdzenia.length;
    const rzedy = rozlozKartkiWRzedy(
      kartkiRdzenia.length,
      widok.ukladKartek,
      widok.kartekWRzedzie,
    );
    const wyrys: HTMLElement[] = [];
    for (const rzad of rzedy) {
      const wiersz = document.createElement('div');
      wiersz.className = 'ms-praca__rzad';
      for (const miejsceRzedu of rzad) {
        if (miejsceRzedu === null) {
          const puste = document.createElement('div');
          puste.className = 'ms-strona ms-strona--puste';
          puste.setAttribute('aria-hidden', 'true');
          wiersz.append(puste);
          continue;
        }
        const kartka = kartkiRdzenia.find((pozycja) => pozycja.numer === miejsceRzedu);
        if (kartka === undefined) continue;
        const oprawa = document.createElement('div');
        oprawa.className = 'ms-strona ms-strona--rdzen';
        oprawa.dataset['kartka'] = String(kartka.numer);
        oprawa.dataset['zrodloWyrysu'] = 'rdzen';
        const obraz = document.createElement('img');
        obraz.className = 'ms-strona__obraz';
        obraz.src = kartka.zrodlo;
        obraz.alt = `Strona ${kartka.numer} z ${kartkiRdzenia.length} — wyrys rdzenia`;
        obraz.loading = 'lazy';
        oprawa.append(obraz);
        wiersz.append(oprawa);
      }
      wyrys.push(wiersz);
    }
    pole.replaceChildren(...wyrys);
    pole.style.setProperty('--ms-szerokosc-pola', `${szerokoscPolaPunkty(stronaBiezaca)}px`);
  }

  function pokaz(tresc: string, wymus = false): void {
    if (!wymus && trescWyrysowana === tresc) {
      // Treść w oknie jest już tą treścią — przerysowanie zabrałoby kursor
      // i nic nie zmieniło. To jest przypadek każdego naciśnięcia klawisza:
      // powierzchnia oddała treść stanowi, a stan wraca z nią tutaj.
      trescBiezaca = tresc;
      return;
    }
    trescBiezaca = tresc;
    trescWyrysowana = tresc;
    blokiBiezace = czytajBloki(tresc);
    akapity.przytnij(blokiBiezace.length);
    if (zrodlo.value !== tresc) zrodlo.value = tresc;
    if (trybBiezacy === 'zrodlowy') {
      liczbaStronBiezaca = 1;
      return;
    }

    // Podgląd wydania rysuje kartki RDZENIA, gdy je ma. Wyrys rdzenia jest
    // wynikiem jego silnika składu, więc pokazuje typografię, paginację, nagłówek
    // i stopkę tak, jak wyjdą w wydaniu — czego kartka złożona przez przeglądarkę
    // nie umie obiecać. Brak obrazów nie jest usterką: wyrys rdzenia dojeżdża po
    // wywołaniu `studio.preview.render`, a do tego czasu podgląd pokazuje kartki
    // liczone tutaj i mówi to wprost paskiem stanu.
    if (trybBiezacy === 'wydanie' && kartkiRdzenia.length > 0) {
      wyrysujKartkiRdzenia();
      odswiezPasek();
      return;
    }

    const miejsce = zapamietajKursor();
    const objete = trybBiezacy === 'roznica' ? blokiFragmentow() : new Map<number, StudioDiffHunk>();

    // Krok pierwszy: bloki na jednej kartce, żeby przeglądarka je zmierzyła.
    const pierwsza = utworzKartke(1, 1);
    pole.replaceChildren(pierwsza.kartka);
    // Wpis układu, a nie sam element: pasek zatwierdzenia jest elementem bez
    // bloku treści, więc podział na strony musi wiedzieć, który wpis niesie blok,
    // a który stoi obok niego. Bez tego numer wpisu rozjechałby się z numerem
    // bloku i podział jawny wypadałby w niewłaściwym miejscu.
    const wpisy: { element: HTMLElement; blok: number }[] = [];
    let poczatek = 0;
    blokiBiezace.forEach((blok, numer) => {
      const wyrys = wyrysBloku(blok);
      przypnijAkapit(wyrys, numer, blok.rodzaj);
      const dlugosc = zapiszBlok(blok).length;
      let zmianyBloku: string[] = [];
      if (blok.rodzaj !== 'tabela' && blok.rodzaj !== 'kod') {
        zmianyBloku = nalozZmiany(wyrys, blok, poczatek);
        nalozKotwice(wyrys, numer, poczatek, dlugosc);
      }
      const fragment = objete.get(numer);
      if (fragment !== undefined) {
        wyrys.dataset['fragment'] = String(fragment.index);
        wyrys.dataset['rodzajFragmentu'] = fragment.kind;
      }
      wpisy.push({ element: wyrys, blok: numer });
      pierwsza.pole.append(wyrys);
      // Pasek zatwierdzenia jest osobnym wpisem układu, więc wchodzi do podziału
      // na strony jak każdy inny — inaczej stałby na kartce, z której zmiana
      // właśnie zeszła.
      if (zmianyBloku.length > 0 && trybAdiustacji === 'cala') {
        const pasek = pasekZatwierdzenia(zmianyBloku);
        wpisy.push({ element: pasek, blok: -1 });
        pierwsza.pole.append(pasek);
      }
      poczatek += dlugosc + 2;
    });

    // Krok drugi: podział na kartki wedle zmierzonych wysokości.
    const wysokosc = wysokoscPolaPunkty(stronaBiezaca);
    const strony =
      wpisy.length > GRANICA_PODZIALU
        ? [wpisy.map((_, numer) => numer)]
        : rozdzielNaStrony(
            wpisy.length,
            wysokosc,
            (numer) => wpisy[numer]?.element.offsetHeight ?? 0,
            (numer) => {
              const wpis = wpisy[numer];
              if (wpis === undefined || wpis.blok < 0) return false;
              return blokiBiezace[wpis.blok]?.rodzaj === 'podzial-strony';
            },
          );
    liczbaStronBiezaca = strony.length;

    const kartki = new Map<number, HTMLElement>();
    strony.forEach((numery, indeks) => {
      const kartka = utworzKartke(indeks + 1, strony.length);
      for (const numer of numery) {
        const wpis = wpisy[numer];
        if (wpis !== undefined) kartka.pole.append(wpis.element);
      }
      kartki.set(indeks + 1, kartka.kartka);
    });

    // Kartki idą w rzędach wedle układu wybranego przez Operatora. Rozkładówka
    // stawia w rzędzie pierwszym puste miejsce okładki — element bez kartki, żeby
    // strona pierwsza stanęła po prawej, tak jak w oprawionym pismie.
    const rzedy = rozlozKartkiWRzedy(strony.length, widok.ukladKartek, widok.kartekWRzedzie);
    const wyrys: HTMLElement[] = [];
    for (const rzad of rzedy) {
      const wiersz = document.createElement('div');
      wiersz.className = 'ms-praca__rzad';
      for (const miejsceRzedu of rzad) {
        if (miejsceRzedu === null) {
          const puste = document.createElement('div');
          puste.className = 'ms-strona ms-strona--puste';
          puste.setAttribute('aria-hidden', 'true');
          wiersz.append(puste);
          continue;
        }
        const kartka = kartki.get(miejsceRzedu);
        if (kartka !== undefined) wiersz.append(kartka);
      }
      wyrys.push(wiersz);
    }
    pole.replaceChildren(...wyrys);
    pole.style.setProperty('--ms-szerokosc-pola', `${szerokoscPolaPunkty(stronaBiezaca)}px`);
    przywrocKursor(miejsce);
    odswiezPasek();
    odswiezLinijkiZKursora();
  }

  /**
   * Przestawia znaczniki linijek wedle miejsca kursora i zaznaczenia.
   *
   * Linijka pokazuje wcięcia akapitu, w którym stoi kursor, jego położenie
   * i granice zaznaczenia — inaczej byłaby miarką bez związku z tym, co Operator
   * właśnie pisze. Krawędzie kolumn pojawiają się wyłącznie w tabeli: chwyt
   * kolumny nad akapitem nie miałby czego przestawiać.
   */
  function odswiezLinijkiZKursora(): void {
    if (!widok.linijkiWidoczne || trybBiezacy === 'zrodlowy') return;
    linijkaPozioma.ustawWciecie(wciecieAkapituKursora());
    const numer = numerBlokuKursora();
    linijkaPozioma.ustawTabulatory(numer < 0 ? [] : (tabulatoryBlokow.get(numer) ?? []));
    linijkaPozioma.ustawKrawedzieKolumn(krawedzieTabeliKursora(numer));

    const wybor = document.getSelection();
    if (wybor === null || wybor.rangeCount === 0) {
      linijkaPozioma.ustawKursor(null);
      linijkaPionowa.ustawKursor(null);
      linijkaPozioma.ustawZaznaczenie(null);
      return;
    }
    const zakres = wybor.getRangeAt(0);
    const blok = blokWezla(zakres.startContainer);
    const kartka = blok?.closest<HTMLElement>('.ms-strona__pole') ?? null;
    if (kartka === null) {
      linijkaPozioma.ustawKursor(null);
      linijkaPionowa.ustawKursor(null);
      return;
    }
    const krotnosc = stronaBiezaca.skala / 100;
    const polePisania = kartka.getBoundingClientRect();
    const kursor = zakres.getBoundingClientRect();
    linijkaPozioma.ustawKursor(zPunktow((kursor.left - polePisania.left) / krotnosc));
    linijkaPionowa.ustawKursor(zPunktow((kursor.top - polePisania.top) / krotnosc));
    linijkaPozioma.ustawZaznaczenie(
      zakres.collapsed
        ? null
        : {
            odMm: zPunktow((kursor.left - polePisania.left) / krotnosc),
            doMm: zPunktow((kursor.right - polePisania.left) / krotnosc),
          },
    );
  }

  /**
   * Krawędzie kolumn tabeli, w której stoi kursor.
   *
   * Zmierzone z wyrysu, a nie zgadnięte z liczby kolumn: tabela ustawia szerokości
   * wedle treści, więc krawędź na linijce ma stać tam, gdzie naprawdę stoi krawędź
   * komórki. Blok, który tabelą nie jest, oddaje wykaz pusty.
   */
  function krawedzieTabeliKursora(numer: number): readonly number[] {
    if (numer < 0) return [];
    if (blokiBiezace[numer]?.rodzaj !== 'tabela') return [];
    const zapamietane = kolumnyTabel.get(numer);
    if (zapamietane !== undefined) return zapamietane;
    const tabela = pole.querySelector<HTMLElement>(`[data-numer-bloku='${numer}']`);
    const wiersz = tabela?.querySelector('tr');
    if (wiersz === null || wiersz === undefined) return [];
    const krotnosc = stronaBiezaca.skala / 100;
    const lewa = tabela?.getBoundingClientRect().left ?? 0;
    const krawedzie: number[] = [];
    const komorki = Array.from(wiersz.children);
    // Krawędź OSTATNIEJ komórki jest krawędzią tabeli, nie granicą między
    // kolumnami — chwytanie jej zmieniałoby szerokość całej tabeli, a to jest
    // czynność inna i nie ta, o którą prosi zlecenie.
    for (const komorka of komorki.slice(0, -1)) {
      const prostokat = komorka.getBoundingClientRect();
      krawedzie.push(zPunktow((prostokat.right - lewa) / krotnosc));
    }
    kolumnyTabel.set(numer, krawedzie);
    return krawedzie;
  }

  /* ── Pisanie ─────────────────────────────────────────────────────────────── */

  pole.contentEditable = 'true';
  pole.spellcheck = true;
  pole.setAttribute('role', 'textbox');
  pole.setAttribute('aria-multiline', 'true');
  pole.setAttribute('aria-label', 'Treść dokumentu na kartce');

  pole.addEventListener('input', () => {
    const zebrane = zbierzTresc();
    trescBiezaca = zebrane;
    czynnosci.naTresc(zebrane);
  });

  /** Zbiera treść ze wszystkich kartek — bloki idą w kolejności kartek. */
  function zbierzTresc(): string {
    const czesci: string[] = [];
    for (const kartka of Array.from(pole.querySelectorAll<HTMLElement>('.ms-strona__pole'))) {
      const zapis = serializujPowierzchnie(kartka);
      if (zapis !== '') czesci.push(zapis);
    }
    return czesci.join('\n\n');
  }

  zrodlo.addEventListener('input', () => {
    trescBiezaca = zrodlo.value;
    czynnosci.naTresc(zrodlo.value);
  });

  /** Zaznaczenie w zapisie treści; `null`, gdy kursor stoi bez zaznaczenia. */
  function zaznaczenie(): ZakresZaznaczenia | null {
    if (trybBiezacy === 'zrodlowy') {
      const od = zrodlo.selectionStart;
      const doZnaku = zrodlo.selectionEnd;
      return od === doZnaku ? null : { poczatek: od, koniec: doZnaku };
    }
    const wybor = document.getSelection();
    if (wybor === null || wybor.rangeCount === 0) return null;
    const zakres = wybor.getRangeAt(0);
    const od = pozycjaWTresci(zakres.startContainer, zakres.startOffset);
    const doZnaku = pozycjaWTresci(zakres.endContainer, zakres.endOffset);
    if (od === null || doZnaku === null || od === doZnaku) return null;
    return { poczatek: Math.min(od, doZnaku), koniec: Math.max(od, doZnaku) };
  }

  /**
   * Położenie punktu w zapisie treści.
   *
   * Liczone jest sumą zapisów bloków poprzedzających i zapisem bloku bieżącego
   * do punktu — bo tak samo liczy je `zapiszBloki`, którym treść jedzie do
   * rdzenia. Punkt poza treścią oddaje `null`, a nie zero: zero jest położeniem
   * prawdziwym i pomyłka tutaj wysłałaby operację na początek dokumentu.
   */
  function pozycjaWTresci(wezel: Node, offset: number): number | null {
    const blok = blokWezla(wezel);
    if (blok === null) return null;
    const numer = Number(blok.dataset['numerBloku'] ?? '-1');
    if (numer < 0) return null;

    let poczatek = 0;
    for (let i = 0; i < numer; i += 1) {
      const wczesniejszy = blokiBiezace[i];
      if (wczesniejszy !== undefined) poczatek += zapiszBlok(wczesniejszy).length + 2;
    }
    const blokTresci = blokiBiezace[numer];
    const przedrostek = blokTresci === undefined ? 0 : dlugoscPrzedrostka(blokTresci);
    const wewnatrz = dlugoscDoPunktu(blok, wezel, offset) ?? 0;
    return poczatek + przedrostek + wewnatrz;
  }

  /** Długość znacznika, który zapis bloku stawia przed jego treścią. */
  function dlugoscPrzedrostka(blok: BlokTresci): number {
    const zapis = zapiszBlok(blok);
    const widoczny = blok.wiersze[0] ?? '';
    if (widoczny === '') return 0;
    const pozycja = zapis.indexOf(widoczny);
    return pozycja < 0 ? 0 : pozycja;
  }

  for (const rodzaj of ['keyup', 'mouseup', 'select']) {
    pole.addEventListener(rodzaj, () => {
      czynnosci.naZaznaczenie(zaznaczenie());
      odswiezLinijkiZKursora();
    });
    zrodlo.addEventListener(rodzaj, () => czynnosci.naZaznaczenie(zaznaczenie()));
  }

  // Kartka widoczna liczy się z przewinięcia, bo licznik „kartka 3 z 12" ma mówić
  // o tym, na co Operator patrzy, a nie o tym, gdzie ostatnio kliknął.
  element.addEventListener('scroll', () => {
    const nowa = kartkaPrzyPrzewinieciu(
      element.scrollTop,
      naPunkty(wysokoscKartki()) * (stronaBiezaca.skala / 100),
      24,
      liczbaStronBiezaca,
    );
    if (nowa === kartkaWidoczna) return;
    kartkaWidoczna = nowa;
    odswiezPasek();
  });

  przypnijNastawy();
  odswiezPasek();

  return {
    element,
    pokaz,
    tryb: () => trybBiezacy,
    ustawTryb: ustawTrybWidoku,
    kartekZRdzenia: () => kartkiRdzenia.length,

    /**
     * Kartki rdzenia pchnięte z okna po renderze podglądu.
     *
     * Przerysowanie idzie z wymuszeniem, bo zmieniła się rzecz INNA niż treść —
     * bez wymuszenia warunek „treść ta sama" zatrzymałby wyrys i Operator
     * zobaczyłby kartki okna, choć rdzeń już oddał swoje.
     */
    ustawKartkiRdzenia(kartki) {
      kartkiRdzenia = [...kartki];
      pokaz(trescBiezaca, true);
    },

    /**
     * Nastawy strony pchnięte z okna.
     *
     * Nadpisania Operatora zostają, ale wyłącznie te, których okno tym pchnięciem
     * NIE zmieniło. Pole zmienione w oknie jest nastawą świeższą niż nadpisanie
     * i wygrywa — inaczej wstążka przestałaby działać po pierwszym chwycie na
     * linijce.
     */
    ustawStrone(nowa) {
      for (const klucz of Object.keys(nadpisaniaStrony) as (keyof StronaPracy)[]) {
        if (nowa[klucz] !== stronaOkna[klucz]) delete nadpisaniaStrony[klucz];
      }
      stronaOkna = nowa;
      stronaBiezaca = { ...nowa, ...nadpisaniaStrony };
      przypnijNastawy();
      pokaz(trescBiezaca, true);
    },

    ustawNastawy(nowe) {
      nastawyBiezace = nowe;
      przypnijNastawy();
    },

    ustawKolumny(ile) {
      const wRzedzie = Math.max(1, Math.round(ile));
      przestawWidok({
        ukladKartek: wRzedzie > 1 ? 'obok' : 'jedna',
        kartekWRzedzie: wRzedzie,
      });
    },

    nastawyWidoku: () => widok,
    przestawWidok,
    ustawSkale,
    ustawNastaweSkali,

    /**
     * Skala pamiętana przy dokumencie.
     *
     * Dokument bez zapamiętanej skali zostaje przy skali ostatniej Operatora —
     * powrót do stu procent przy każdym nowym pismie byłby nastawą narzuconą.
     */
    wczytajSkaleDokumentu(idDokumentu) {
      dokumentSkali = idDokumentu;
      const zapamietana = skalaDokumentu(idDokumentu);
      if (zapamietana !== null) ustawSkale(zapamietana);
    },

    skoczDoKartki,
    kartkaBiezaca: () => kartkaWidoczna,

    ustawLinijki(widoczne) {
      przestawWidok({ linijkiWidoczne: widoczne });
    },

    wciecieAkapitu: wciecieAkapituKursora,

    tabulatoryAkapitu() {
      const numer = numerBlokuKursora();
      return numer < 0 ? [] : (tabulatoryBlokow.get(numer) ?? []);
    },

    strona: () => stronaBiezaca,

    ustawZmiany(zmiany) {
      zmianySledzone = zmiany;
      pokaz(trescBiezaca, true);
    },

    ustawAdiustacje(tryb) {
      trybAdiustacji = tryb;
      element.dataset['adiustacja'] = tryb;
      pokaz(trescBiezaca, true);
    },

    adiustacja: () => trybAdiustacji,

    ustawKomentarze(komentarze) {
      komentarzeBiezace = komentarze;
      pokaz(trescBiezaca, true);
    },

    polozenieKotwicy(idKomentarza) {
      const kotwica = pole.querySelector<HTMLElement>(`[data-komentarz='${idKomentarza}']`);
      if (kotwica === null) return null;
      const prostokat = kotwica.getBoundingClientRect();
      const powloka = element.getBoundingClientRect();
      return prostokat.top - powloka.top + element.scrollTop;
    },

    skoczDoZmiany(wPrzod) {
      const znaczniki = Array.from(pole.querySelectorAll<HTMLElement>('[data-zmiana]'));
      if (znaczniki.length === 0) return null;
      const miejsce = zapamietajKursor();
      const biezacy = miejsce === null ? -1 : miejsce.blok;
      const uporzadkowane = wPrzod ? znaczniki : [...znaczniki].reverse();
      const wybrany =
        uporzadkowane.find((znacznik) => {
          const blok = blokWezla(znacznik);
          const numer = blok === null ? -1 : Number(blok.dataset['numerBloku'] ?? '-1');
          return wPrzod ? numer > biezacy : numer < biezacy && numer >= 0;
        }) ?? uporzadkowane[0];
      if (wybrany === undefined) return null;
      wybrany.scrollIntoView({ block: 'center' });
      const tekst = wybrany.firstChild;
      if (tekst !== null && tekst.nodeType === Node.TEXT_NODE) ustawKursorNa(tekst as Text, 0);
      return wybrany.dataset['zmiana'] ?? null;
    },

    ustawFragmenty(fragmenty) {
      fragmentyRoznicy = fragmenty;
      if (trybBiezacy === 'roznica') pokaz(trescBiezaca, true);
    },

    zaznaczenie,

    blokKursora() {
      const miejsce = zapamietajKursor();
      return miejsce === null ? -1 : miejsce.blok;
    },

    polozenieKursora() {
      const wybor = document.getSelection();
      if (wybor === null || wybor.rangeCount === 0) return null;
      const zakres = wybor.getRangeAt(0);
      if (blokWezla(zakres.startContainer) === null) return null;
      const prostokat = zakres.getBoundingClientRect();
      const powloka = element.getBoundingClientRect();
      return {
        x: prostokat.left - powloka.left + element.scrollLeft,
        y: prostokat.bottom - powloka.top + element.scrollTop,
      };
    },

    ustawStylBloku(rodzaj) {
      const miejsce = zapamietajKursor();
      if (miejsce === null) return;
      const blok = blokiBiezace[miejsce.blok];
      if (blok === undefined) return;
      blokiBiezace[miejsce.blok] = { rodzaj, wiersze: blok.wiersze };
      const zapis = blokiBiezace.map((pozycja) => zapiszBlok(pozycja)).join('\n\n');
      czynnosci.naTresc(zapis);
      pokaz(zapis, true);
    },

    bloki: () => blokiBiezace,
    liczbaStron: () => liczbaStronBiezaca,

    opis() {
      return (
        `${opiszStrone(stronaBiezaca)} · ` +
        `${opiszUkladKartek(widok.ukladKartek, widok.kartekWRzedzie, widok.przewijanie, liczbaStronBiezaca)} · ` +
        `bloków ${blokiBiezace.length} · ${linijkaPozioma.opis()}` +
        (zdanieDruku === '' ? '' : ` · ${zdanieDruku}`)
      );
    },
  };
}
