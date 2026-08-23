import { StudioPageOrientation, type StudioPageSetup } from '../../../../shared/contract';

/**
 * Nastawy strony okna pracy z dokumentem — kartka, marginesy, nagłówek, stopka,
 * numeracja i skala.
 *
 * Kształt nastaw NIE jest wymyślony: to `StudioPageSetup` z kontraktu, ten sam,
 * który jedzie w profilu wydania (`studio.export.profile.save`) i w renderze
 * podglądu (`studio.preview.render`). Okno pracy pokazuje więc kartkę o tych
 * samych nastawach, którymi rdzeń wyda dokument — inaczej „widzisz to, co
 * dostaniesz" byłoby obietnicą bez pokrycia.
 *
 * Rozmiar nośnika w milimetrach stoi tutaj, bo kontrakt niesie samo oznaczenie
 * („A4"), a nie jego wymiary. Wykaz jest zamknięty i nazwany — nośnik spoza
 * wykazu bierze wymiary A4 i okno mówi o tym wprost, zamiast rysować kartkę
 * o zgadniętym rozmiarze.
 *
 * Plik nie zna DOM: wejściem są nastawy, wyjściem liczby i zdania.
 */

/** Rodzaj nośnika — arkusz do pisania albo koperta do nadruku. */
export type RodzajNosnika = 'arkusz' | 'koperta' | 'wlasny';

/** Wymiary nośnika w milimetrach. */
export interface WymiaryNosnika {
  oznaczenie: string;
  szerokoscMm: number;
  wysokoscMm: number;
  rodzaj: RodzajNosnika;
}

/**
 * Nośniki wbudowane — wykaz ZAPASOWY, nie źródło prawdy.
 *
 * Źródłem prawdy jest rdzeń: `studio.page.paper.list` oddaje szereg A, szereg B,
 * Letter, Legal, Tabloid oraz koperty DL, C4, C5 i C6 wraz z wymiarami, czytając
 * je z JEDNEGO wspólnego wykazu rdzenia (`nosniki_druku_wspolne.go`). Okno pracy
 * woła tę komendę przy wczytaniu dokumentu i wchłania odpowiedź
 * (`wchlonNosnikiRdzenia`), więc kartka rysuje się wymiarami rdzenia, a nie
 * wymiarami stąd.
 *
 * Ten wykaz stoi tutaj wyłącznie dlatego, że kartkę trzeba narysować, ZANIM
 * odpowiedź rdzenia dojedzie — rysowanie w rozmiarze zgadniętym byłoby kłamstwem
 * o nośniku. Poprawki wymiarów wnosi się do wykazu rdzenia; wpis tutaj ma się do
 * niego równać, a nie odwrotnie. Sprawdzian `nastawy-strony.test.ts` pilnuje, żeby
 * wchłonięcie naprawdę nadpisywało wymiary wbudowane.
 */
const NOSNIKI_WBUDOWANE: readonly WymiaryNosnika[] = [
  { oznaczenie: 'A0', szerokoscMm: 841, wysokoscMm: 1189, rodzaj: 'arkusz' },
  { oznaczenie: 'A1', szerokoscMm: 594, wysokoscMm: 841, rodzaj: 'arkusz' },
  { oznaczenie: 'A2', szerokoscMm: 420, wysokoscMm: 594, rodzaj: 'arkusz' },
  { oznaczenie: 'A3', szerokoscMm: 297, wysokoscMm: 420, rodzaj: 'arkusz' },
  { oznaczenie: 'A4', szerokoscMm: 210, wysokoscMm: 297, rodzaj: 'arkusz' },
  { oznaczenie: 'A5', szerokoscMm: 148, wysokoscMm: 210, rodzaj: 'arkusz' },
  { oznaczenie: 'A6', szerokoscMm: 105, wysokoscMm: 148, rodzaj: 'arkusz' },
  { oznaczenie: 'B1', szerokoscMm: 707, wysokoscMm: 1000, rodzaj: 'arkusz' },
  { oznaczenie: 'B2', szerokoscMm: 500, wysokoscMm: 707, rodzaj: 'arkusz' },
  { oznaczenie: 'B3', szerokoscMm: 353, wysokoscMm: 500, rodzaj: 'arkusz' },
  { oznaczenie: 'B4', szerokoscMm: 250, wysokoscMm: 353, rodzaj: 'arkusz' },
  { oznaczenie: 'B5', szerokoscMm: 176, wysokoscMm: 250, rodzaj: 'arkusz' },
  { oznaczenie: 'Letter', szerokoscMm: 216, wysokoscMm: 279, rodzaj: 'arkusz' },
  { oznaczenie: 'Legal', szerokoscMm: 216, wysokoscMm: 356, rodzaj: 'arkusz' },
  { oznaczenie: 'Tabloid', szerokoscMm: 279, wysokoscMm: 432, rodzaj: 'arkusz' },
  { oznaczenie: 'DL', szerokoscMm: 220, wysokoscMm: 110, rodzaj: 'koperta' },
  { oznaczenie: 'C4', szerokoscMm: 324, wysokoscMm: 229, rodzaj: 'koperta' },
  { oznaczenie: 'C5', szerokoscMm: 229, wysokoscMm: 162, rodzaj: 'koperta' },
  { oznaczenie: 'C6', szerokoscMm: 162, wysokoscMm: 114, rodzaj: 'koperta' },
];

/**
 * Nośniki widoczne w oknie — wbudowane, uzupełnione tym, co oddał rdzeń.
 *
 * Zmienna, a nie stała, bo wykaz rdzenia dojeżdża po zbudowaniu okna. Wykaz jest
 * zawsze pełny: dopóki rdzeń nie odpowie, obowiązują wpisy wbudowane.
 */
let nosnikiBiezace: readonly WymiaryNosnika[] = NOSNIKI_WBUDOWANE;

/** Wykaz nośników obowiązujący teraz. */
export function nosniki(): readonly WymiaryNosnika[] {
  return nosnikiBiezace;
}

/**
 * Wchłania wykaz nośników rdzenia (`studio.page.paper.list`).
 *
 * Wpis rdzenia o oznaczeniu znanym NADPISUJE wymiary wbudowane — rdzeń jest
 * źródłem prawdy. Wpis o oznaczeniu nieznanym dochodzi na koniec. Wpisy
 * wbudowane, których rdzeń nie zna, zostają: wykaz widoczny ma być zawsze pełny,
 * a pozycja zniknięta z okna wyglądałaby na usterkę.
 *
 * Rodzaj nośnika bierze się z odpowiedzi rdzenia, gdy ta go niesie (pole `kind`
 * kontraktu: `sheet` albo `envelope`), a nie z wpisu wbudowanego. Inaczej koperta
 * dołożona kiedyś do wykazu rdzenia wchodziłaby do okna jako arkusz i nadruk
 * koperty nie miałby się na czym wykonać.
 */
export function wchlonNosnikiRdzenia(
  wykaz: readonly {
    oznaczenie: string;
    szerokoscMm: number;
    wysokoscMm: number;
    rodzaj?: RodzajNosnika;
  }[],
): readonly WymiaryNosnika[] {
  const zlozone = new Map<string, WymiaryNosnika>();
  for (const nosnik of NOSNIKI_WBUDOWANE) zlozone.set(nosnik.oznaczenie.toLowerCase(), nosnik);
  for (const pozycja of wykaz) {
    if (pozycja.szerokoscMm <= 0 || pozycja.wysokoscMm <= 0) continue;
    const klucz = pozycja.oznaczenie.trim().toLowerCase();
    if (klucz === '') continue;
    const znany = zlozone.get(klucz);
    zlozone.set(klucz, {
      oznaczenie: pozycja.oznaczenie.trim(),
      szerokoscMm: pozycja.szerokoscMm,
      wysokoscMm: pozycja.wysokoscMm,
      rodzaj: pozycja.rodzaj ?? znany?.rodzaj ?? 'arkusz',
    });
  }
  nosnikiBiezace = [...zlozone.values()];
  return nosnikiBiezace;
}

/**
 * Wchłania odpowiedź komendy `studio.page.paper.list` wprost.
 *
 * Przekład stoi tutaj, a nie u wołacza, bo tu leży wiedza o kształcie wykazu:
 * kontrakt nazywa pola po angielsku (`name`, `widthMm`, `heightMm`, `kind`),
 * a rysowanie kartki liczy w polach polskich. Dwa przekłady tego samego kształtu
 * rozjechałyby się przy pierwszym dołożonym polu.
 */
export function wchlonNosnikiKontraktu(
  papiery: readonly { name: string; widthMm: number; heightMm: number; kind?: string }[],
): readonly WymiaryNosnika[] {
  return wchlonNosnikiRdzenia(
    papiery.map((papier) => ({
      oznaczenie: papier.name,
      szerokoscMm: papier.widthMm,
      wysokoscMm: papier.heightMm,
      ...(papier.kind === 'envelope' ? { rodzaj: 'koperta' as RodzajNosnika } : {}),
    })),
  );
}

/**
 * Nośnik własny podany wymiarami — Operator nie jest uwięziony w wykazie.
 *
 * Wymiar niedodatni albo nieskończony oddaje `null`, a nie kartkę o rozmiarze
 * zerowym: kartka bez wymiaru nie jest nastawą, jest usterką.
 */
export function nosnikWlasny(
  szerokoscMm: number,
  wysokoscMm: number,
): WymiaryNosnika | null {
  if (!Number.isFinite(szerokoscMm) || !Number.isFinite(wysokoscMm)) return null;
  if (szerokoscMm <= 0 || wysokoscMm <= 0) return null;
  return {
    oznaczenie: `własny ${szerokoscMm}×${wysokoscMm} mm`,
    szerokoscMm,
    wysokoscMm,
    rodzaj: 'wlasny',
  };
}

/** Nośnik domyślny — pismo urzędowe w tym kraju idzie na A4. */
const NOSNIK_DOMYSLNY: WymiaryNosnika = {
  oznaczenie: 'A4',
  szerokoscMm: 210,
  wysokoscMm: 297,
  rodzaj: 'arkusz',
};

/**
 * Wykaz nośników w postaci stałej — zgodność z wołaczami sprzed wykazu zmiennego.
 *
 * Wstążka czyta go raz, przy składaniu kontrolki; przestawienie wykazu z rdzenia
 * widać przez `nosniki()`.
 */
export const NOSNIKI: readonly WymiaryNosnika[] = NOSNIKI_WBUDOWANE;

/**
 * Punkty na milimetr przy 96 punktach na cal.
 *
 * 96 dpi jest odniesieniem CSS dla jednostki `px`, więc kartka rysowana tym
 * przelicznikiem ma na ekranie rozmiar, który po wydruku w skali 100 % zgadza
 * się z nośnikiem.
 */
const PUNKTOW_NA_MM = 96 / 25.4;

/** Jednostka podziałki linijek — wybór Operatora, nie rozstrzygnięcie kodu. */
export type JednostkaMiary = 'mm' | 'cal';

/** Nastawy strony wraz z wartościami domyślnymi wypełnionymi. */
export interface StronaPracy {
  nosnik: WymiaryNosnika;
  orientacja: StudioPageOrientation;
  marginesGoraMm: number;
  marginesDolMm: number;
  marginesLewyMm: number;
  marginesPrawyMm: number;
  naglowek: string;
  stopka: string;
  numeracja: boolean;
  /** Skala widoku w procentach — nastawa `scale` kontraktu. */
  skala: number;
  /**
   * Margines na oprawę w milimetrach — pas doliczany do marginesu wewnętrznego.
   *
   * Kontrakt tego pola nie ma (`StudioPageSetup` niesie cztery marginesy i nic
   * poza nimi), więc oprawa żyje przez sesję okna i jest zgłoszona jako brak
   * pozycji kontraktu. Rysowana jest prawdziwie: pas oprawy zabiera pole pisania.
   */
  marginesOprawyMm: number;
  /** Strona, przy której stoi oprawa. */
  stronaOprawy: 'wewnatrz' | 'gora';
  /**
   * Marginesy odbicia dla druku dwustronnego.
   *
   * Włączone znaczy, że margines „lewy" jest marginesem WEWNĘTRZNYM: na stronie
   * nieparzystej stoi po lewej, na parzystej po prawej. Bez tego oprawa wypadałaby
   * raz w rowku, raz na krawędzi.
   */
  marginesyOdbicia: boolean;
  /** Jednostka podziałki linijek i pól wymiarowych. */
  jednostka: JednostkaMiary;
  /** Widoczna granica marginesu w treści — przełącznik Operatora. */
  graniceMarginesow: boolean;
}

/** Nastawy domyślne: A4 pionowo, marginesy po 20 mm, numeracja stron. */
export function domyslnaStrona(): StronaPracy {
  return {
    nosnik: NOSNIK_DOMYSLNY,
    orientacja: StudioPageOrientation.Pionowa,
    marginesGoraMm: 20,
    marginesDolMm: 20,
    marginesLewyMm: 20,
    marginesPrawyMm: 20,
    naglowek: '',
    stopka: '',
    numeracja: true,
    skala: 100,
    marginesOprawyMm: 0,
    stronaOprawy: 'wewnatrz',
    marginesyOdbicia: false,
    jednostka: 'mm',
    graniceMarginesow: true,
  };
}

/** Nastawa gotowa marginesów — cztery liczby pod jedną nazwą. */
export interface NastawaMarginesow {
  kod: string;
  nazwa: string;
  goraMm: number;
  dolMm: number;
  lewyMm: number;
  prawyMm: number;
  /** Oprawa nastawy; zero znaczy „bez oprawy". */
  oprawaMm: number;
  opis: string;
}

/**
 * Nastawy gotowe marginesów, wzorem pakietu biurowego.
 *
 * „Własne" nastawy w wykazie nie ma i mieć nie może: własne są tym, co Operator
 * wpisze w cztery pola, a nie kolejną pozycją do wybrania.
 */
export const NASTAWY_MARGINESOW: readonly NastawaMarginesow[] = [
  {
    kod: 'waskie',
    nazwa: 'Wąskie',
    goraMm: 13,
    dolMm: 13,
    lewyMm: 13,
    prawyMm: 13,
    oprawaMm: 0,
    opis: 'Po 13 mm z każdej strony — więcej tekstu na kartce, mniej miejsca na uwagi na marginesie.',
  },
  {
    kod: 'normalne',
    nazwa: 'Normalne',
    goraMm: 20,
    dolMm: 20,
    lewyMm: 20,
    prawyMm: 20,
    oprawaMm: 0,
    opis: 'Po 20 mm — postać pisma urzędowego.',
  },
  {
    kod: 'szerokie',
    nazwa: 'Szerokie',
    goraMm: 25,
    dolMm: 25,
    lewyMm: 50,
    prawyMm: 50,
    oprawaMm: 0,
    opis: 'Góra i dół 25 mm, boki po 50 mm — wiersz krótszy, czyta się szybciej.',
  },
  {
    kod: 'oprawa',
    nazwa: 'Do oprawy',
    goraMm: 20,
    dolMm: 20,
    lewyMm: 20,
    prawyMm: 20,
    oprawaMm: 12,
    opis:
      'Marginesy 20 mm i pas oprawy 12 mm po stronie wewnętrznej wraz z marginesami odbicia — ' +
      'do wydruku dwustronnego zszywanego albo klejonego.',
  },
];

/** Nakłada nastawę gotową na nastawy strony. */
export function zastosujNastaweMarginesow(
  strona: StronaPracy,
  nastawa: NastawaMarginesow,
): StronaPracy {
  return {
    ...strona,
    marginesGoraMm: nastawa.goraMm,
    marginesDolMm: nastawa.dolMm,
    marginesLewyMm: nastawa.lewyMm,
    marginesPrawyMm: nastawa.prawyMm,
    marginesOprawyMm: nastawa.oprawaMm,
    marginesyOdbicia: nastawa.oprawaMm > 0 ? true : strona.marginesyOdbicia,
  };
}

/** Marginesy poziome kartki po doliczeniu oprawy i odbicia. */
export interface MarginesyKartki {
  goraMm: number;
  dolMm: number;
  lewyMm: number;
  prawyMm: number;
}

/**
 * Marginesy kartki o wskazanym numerze.
 *
 * Numer strony rozstrzyga wyłącznie przy marginesach odbicia: strona nieparzysta
 * jest stroną prawą rozkładówki, więc jej margines wewnętrzny stoi po lewej.
 * Oprawa dochodzi do marginesu wewnętrznego, a przy oprawie u góry — do górnego.
 */
export function marginesyKartki(strona: StronaPracy, numerStrony: number): MarginesyKartki {
  const oprawa = Math.max(0, strona.marginesOprawyMm);
  if (strona.stronaOprawy === 'gora') {
    return {
      goraMm: strona.marginesGoraMm + oprawa,
      dolMm: strona.marginesDolMm,
      lewyMm: strona.marginesLewyMm,
      prawyMm: strona.marginesPrawyMm,
    };
  }
  const parzysta = strona.marginesyOdbicia && numerStrony % 2 === 0;
  const wewnatrz = (parzysta ? strona.marginesPrawyMm : strona.marginesLewyMm) + oprawa;
  const zewnatrz = parzysta ? strona.marginesLewyMm : strona.marginesPrawyMm;
  return {
    goraMm: strona.marginesGoraMm,
    dolMm: strona.marginesDolMm,
    lewyMm: parzysta ? zewnatrz : wewnatrz,
    prawyMm: parzysta ? wewnatrz : zewnatrz,
  };
}

/** Odnajduje nośnik po oznaczeniu; oznaczenie nieznane oddaje `null`. */
export function nosnikPoOznaczeniu(oznaczenie: string): WymiaryNosnika | null {
  const szukane = oznaczenie.trim().toLowerCase();
  return nosnikiBiezace.find((nosnik) => nosnik.oznaczenie.toLowerCase() === szukane) ?? null;
}

/** Milimetry na cale — jednostka linijki jest wyborem Operatora. */
export function naCale(milimetry: number): number {
  return milimetry / 25.4;
}

/** Cale na milimetry. */
export function zCali(cale: number): number {
  return cale * 25.4;
}

/** Punkty ekranu na milimetry — droga powrotna chwytu linijki. */
export function zPunktow(punkty: number): number {
  return punkty / PUNKTOW_NA_MM;
}

/** Długość w jednostce wybranej, zapisana zwięźle. */
export function opiszDlugosc(milimetry: number, jednostka: JednostkaMiary): string {
  return jednostka === 'cal'
    ? `${naCale(milimetry).toFixed(2)}″`
    : `${Math.round(milimetry * 10) / 10} mm`;
}

/**
 * Składa nastawy okna z nastaw profilu wydania.
 *
 * Pole niepodane bierze wartość domyślną, a nie zero: margines zerowy jest
 * nastawą, którą Operator może wybrać świadomie, więc „nie podano" i „podano
 * zero" nie mogą znaczyć tego samego.
 */
export function stronaZProfilu(nastawy: StudioPageSetup | undefined): StronaPracy {
  const strona = domyslnaStrona();
  if (nastawy === undefined) return strona;
  if (nastawy.pageSize !== undefined) {
    strona.nosnik = nosnikPoOznaczeniu(nastawy.pageSize) ?? strona.nosnik;
  }
  if (nastawy.orientation !== undefined) strona.orientacja = nastawy.orientation;
  if (nastawy.marginTop !== undefined) strona.marginesGoraMm = nastawy.marginTop;
  if (nastawy.marginBottom !== undefined) strona.marginesDolMm = nastawy.marginBottom;
  if (nastawy.marginLeft !== undefined) strona.marginesLewyMm = nastawy.marginLeft;
  if (nastawy.marginRight !== undefined) strona.marginesPrawyMm = nastawy.marginRight;
  if (nastawy.header !== undefined) strona.naglowek = nastawy.header;
  if (nastawy.footer !== undefined) strona.stopka = nastawy.footer;
  if (nastawy.pageNumbers !== undefined) strona.numeracja = nastawy.pageNumbers;
  if (nastawy.scale !== undefined) strona.skala = nastawy.scale;
  return strona;
}

/** Składa nastawy kontraktu z nastaw okna — droga do profilu wydania. */
export function profilZeStrony(strona: StronaPracy): StudioPageSetup {
  return {
    pageSize: strona.nosnik.oznaczenie,
    orientation: strona.orientacja,
    marginTop: strona.marginesGoraMm,
    marginBottom: strona.marginesDolMm,
    marginLeft: strona.marginesLewyMm,
    marginRight: strona.marginesPrawyMm,
    header: strona.naglowek,
    footer: strona.stopka,
    pageNumbers: strona.numeracja,
    scale: strona.skala,
  };
}

/** Szerokość kartki w milimetrach z uwzględnieniem orientacji. */
export function szerokoscKartkiMm(strona: StronaPracy): number {
  return strona.orientacja === StudioPageOrientation.Pozioma
    ? strona.nosnik.wysokoscMm
    : strona.nosnik.szerokoscMm;
}

/** Wysokość kartki w milimetrach z uwzględnieniem orientacji. */
export function wysokoscKartkiMm(strona: StronaPracy): number {
  return strona.orientacja === StudioPageOrientation.Pozioma
    ? strona.nosnik.szerokoscMm
    : strona.nosnik.wysokoscMm;
}

/** Milimetry na punkty ekranu. */
export function naPunkty(milimetry: number): number {
  return Math.round(milimetry * PUNKTOW_NA_MM);
}

/**
 * Wysokość pola pisania w punktach — kartka bez marginesów.
 *
 * Wartość jest granicą podziału na strony: blok, który się w niej nie mieści,
 * schodzi na stronę następną.
 */
export function wysokoscPolaPunkty(strona: StronaPracy, numerStrony = 1): number {
  const marginesy = marginesyKartki(strona, numerStrony);
  return naPunkty(wysokoscKartkiMm(strona) - marginesy.goraMm - marginesy.dolMm);
}

/** Szerokość pola pisania w punktach. */
export function szerokoscPolaPunkty(strona: StronaPracy, numerStrony = 1): number {
  const marginesy = marginesyKartki(strona, numerStrony);
  return naPunkty(szerokoscKartkiMm(strona) - marginesy.lewyMm - marginesy.prawyMm);
}

/** Szerokość pola pisania w milimetrach — miara linijki poziomej. */
export function szerokoscPolaMm(strona: StronaPracy, numerStrony = 1): number {
  const marginesy = marginesyKartki(strona, numerStrony);
  return szerokoscKartkiMm(strona) - marginesy.lewyMm - marginesy.prawyMm;
}

/**
 * Zdanie o kartce dla paska stanu.
 *
 * Mówi nośnik, orientację i marginesy, bo to one rozstrzygają, gdzie skończy
 * się strona pierwsza — a Operator ma to wiedzieć, zanim zacznie pisać.
 */
export function opiszStrone(strona: StronaPracy): string {
  const orientacja =
    strona.orientacja === StudioPageOrientation.Pozioma ? 'poziomo' : 'pionowo';
  const oprawa =
    strona.marginesOprawyMm > 0
      ? ` · oprawa ${strona.marginesOprawyMm} mm ${strona.stronaOprawy === 'gora' ? 'u góry' : 'wewnątrz'}`
      : '';
  const odbicia = strona.marginesyOdbicia ? ' · marginesy odbicia' : '';
  return (
    `${strona.nosnik.oznaczenie}${strona.nosnik.rodzaj === 'koperta' ? ' (koperta)' : ''} ` +
    `${orientacja} · ${szerokoscKartkiMm(strona)}×` +
    `${wysokoscKartkiMm(strona)} mm · marginesy ${strona.marginesGoraMm}/${strona.marginesDolMm}/` +
    `${strona.marginesLewyMm}/${strona.marginesPrawyMm} mm${oprawa}${odbicia} · skala ${strona.skala} %`
  );
}

/**
 * Co z nastaw strony dojeżdża do rdzenia, a co zostaje w oknie.
 *
 * Zdanie stoi przy nastawach, a nie w sprawozdaniu, bo Operator ma wiedzieć, co
 * przeżyje zamknięcie okna, ZANIM to zamknie.
 *
 * Zdanie zmieniło treść wraz z odcinkiem postaci dokumentu: droga trwałości nie
 * jest już profilem wydania, a `studio.page.setup.set` wraz z rodziną
 * `studio.page.*` i `studio.section.*`. Oprawa, marginesy odbicia, format własny
 * w milimetrach, kolumny, sekcje, znak wodny oraz nagłówek i stopka osobne dla
 * pierwszej strony i stron parzystych mają w kontrakcie własne pola i dojeżdżają.
 * Jednostka linijki i granice marginesów idą osobno — nastawami widoku
 * (`studio.view.set`), bo dotyczą patrzenia, nie dokumentu.
 */
export const POWOD_NIETRWALOSCI_STRONY =
  'Nastawy strony jadą do rdzenia komendą studio.page.setup.set: nośnik z wykazu albo format ' +
  'własny w milimetrach, orientacja, cztery marginesy, margines na oprawę, marginesy odbicia, ' +
  'nastawa gotowa i kolumny — osobno dla każdej sekcji. Nagłówek, stopka i numeracja idą ' +
  'studio.page.headerfooter.set i studio.page.numbering.set, znak wodny studio.page.watermark.set. ' +
  'Jednostka linijki i granice marginesów należą do nastaw WIDOKU (studio.view.set), a nie do ' +
  'dokumentu: opisują, jak Operator patrzy, a nie jak dokument wyjdzie z drukarki. Skala widoku ' +
  'jedzie oboma polami — profil wydania niesie skalę WYDRUKU, nastawy widoku skalę PATRZENIA.';

/**
 * Rozdziela bloki na strony wedle ich wysokości.
 *
 * Wysokość bloku podaje wołający (`zmierz`), bo zmierzyć ją potrafi wyłącznie
 * przeglądarka — a rachunek podziału ma dać się sprawdzić bez niej. Blok wyższy
 * od całej strony zostaje na stronie własnej: dzielenie go w środku wymagałoby
 * łamania wiersza wewnątrz akapitu, czego ta warstwa nie umie i czego nie udaje.
 *
 * Podział jawny (`podzial-strony`) kończy stronę niezależnie od miejsca, jakie
 * na niej zostało — to nastawa Operatora, a nie skutek rachunku.
 */
export function rozdzielNaStrony(
  liczba: number,
  wysokoscPola: number,
  zmierz: (numer: number) => number,
  czyPodzial: (numer: number) => boolean,
): number[][] {
  const strony: number[][] = [];
  let biezaca: number[] = [];
  let zajete = 0;

  for (let numer = 0; numer < liczba; numer += 1) {
    if (czyPodzial(numer)) {
      biezaca.push(numer);
      strony.push(biezaca);
      biezaca = [];
      zajete = 0;
      continue;
    }
    const wysokosc = zmierz(numer);
    if (biezaca.length > 0 && zajete + wysokosc > wysokoscPola) {
      strony.push(biezaca);
      biezaca = [numer];
      zajete = wysokosc;
      continue;
    }
    biezaca.push(numer);
    zajete += wysokosc;
  }

  // Strona pusta na końcu jest stroną prawdziwą wyłącznie po podziale jawnym:
  // dokument kończący się podziałem ma stronę następną, choćby jeszcze pustą.
  if (biezaca.length > 0 || strony.length === 0) strony.push(biezaca);
  return strony;
}
