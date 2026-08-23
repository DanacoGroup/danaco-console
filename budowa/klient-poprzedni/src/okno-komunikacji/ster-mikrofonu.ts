import { elementIkony } from '../ikony/ikony';
import type { PozycjaMenu } from '../komponenty/menu-drzewo';
import type { Dyktowanie } from './dyktowanie/dyktowanie';
import type { UrzadzenieDzwieku } from './dyktowanie/urzadzenia-dzwieku';
import { utworzSterNastawy, type SterPaska } from './ster-nastawy';

/**
 * Mikrofon — ster paska zlecenia. Jest wyłącznie drogą do gotowej fasady
 * `Dyktowanie`: nie nagrywa, nie zna `MediaRecorder` ani komendy
 * `speech.transcribe` i nie składa wyniku.
 *
 * Komponent ma dwa uchwyty w jednym pudełku: przycisk nagrywania oraz uchwyt
 * drzewa. Nagrywanie jest czynnością, a wybór mikrofonu nastawą; gest
 * „przytrzymaj, aby nagrać" (`dyktowanie/przytrzymanie.ts`) zajmuje
 * `pointerdown` przycisku, a uchwyt menu otwiera się na `click` — jeden
 * przycisk pełniący obie role zaczynałby nagranie przy każdym otwarciu wykazu.
 * Menu jest jedno i biblioteczne (`komponenty/menu-drzewo.ts` przez wspólną
 * obudowę `ster-nastawy.ts`).
 *
 * Drzewo ma dwa poziomy: gałąź „Mikrofon" niesie wykaz urządzeń, a przełącznik
 * „Przytrzymaj, aby nagrać" stoi poziom wyżej, bo dotyczy gestu, nie sprzętu.
 *
 * Mikrofonu bez silnika nie stawiamy wcale — `dostepnosc()` jest pytaniem
 * zadanym zanim cokolwiek powstanie, stąd obietnica w wyniku i `null` zamiast
 * steru. Wyszarzona ikona albo ikona odmawiająca po naciśnięciu obiecywałaby
 * zdolność, której nie ma; krótszy pasek niczego nie obiecuje.
 */

/** Nazwa rodzajowa nastawy — idzie do `aria-label`, nie na ekran. */
const NASTAWA = 'Mikrofon';

/** Klucz gałęzi z wykazem urządzeń; przedrostek oddziela je od przełączników. */
const GALAZ_URZADZEN = 'mikrofon:urzadzenia';
const PRZEDROSTEK_URZADZENIA = 'mikrofon:urzadzenie:';
const KLUCZ_PRZYTRZYMANIA = 'mikrofon:przytrzymanie';

/** Pozycja „cokolwiek wybrał system" — pusty `deviceId` przyjmuje `getUserMedia`. */
const KLUCZ_DOMYSLNEGO = `${PRZEDROSTEK_URZADZENIA}`;
export const NAZWA_DOMYSLNEGO = 'Wejście domyślne systemu';

/** Wartości na uchwycie w trakcie pracy — sygnał, że mikrofon pracuje. */
const WARTOSC_NAGRYWA = 'Nagrywa…';
const WARTOSC_KONCZY = 'Rozpoznaje…';

/**
 * Zdanie o nagraniu bez mowy, w tonie spokojnym: cisza jest prawidłowym
 * wynikiem, a nie awarią (powód pełny stoi przy
 * `dyktowanie/wynik-dyktowania.ts`). Zdanie mówi, co się stało, dlaczego i co
 * z tym zrobić.
 */
export const ZDANIE_BEZ_MOWY =
  'W nagraniu nie było mowy. Silnik przetworzył je do końca i nie usłyszał ani jednego ' +
  'słowa — bywa tak przy pustym pokoju albo wyciszonym wejściu. Nagraj jeszcze raz, ' +
  'mówiąc bliżej mikrofonu.';

/** Zależności steru — wąskie i wstrzykiwane. */
export interface ZaleznosciSteruMikrofonu {
  /** Gotowa fasada ogniwa mowa→tekst; ster jej nie zakłada. */
  dyktowanie: Dyktowanie;
  /**
   * Dokłada rozpoznany tekst do pola wypowiedzi. Wymagana, nie opcjonalna:
   * dyktowanie, którego wynik nie ma dokąd trafić, jest przyciskiem
   * meldującym pracę bez skutku. Kto nie ma pola wypowiedzi, ten nie stawia
   * mikrofonu.
   */
  wstawTekst(tekst: string): void;
}

/** Ster paska ze zdejmowalnymi nasłuchami sprzętu i wyniku. */
export type SterMikrofonu = SterPaska & { rozlacz(): void };

/**
 * Buduje ster mikrofonu albo oddaje `null`, gdy dyktowania nie da się wykonać.
 *
 * Obietnica, bo `dostepnosc()` pyta rdzeń o silnik mowy. Wołający montuje ster
 * dopiero po jej rozstrzygnięciu — pasek ma na to przygotowane gniazdo.
 */
export async function utworzSterMikrofonu(
  zaleznosci: ZaleznosciSteruMikrofonu,
): Promise<SterMikrofonu | null> {
  const { dyktowanie } = zaleznosci;

  const gotowosc = await dyktowanie.dostepnosc();
  if (!gotowosc.dostepne) return null;

  /** Urządzenia widziane przez przeglądarkę; puste do pierwszego odczytu. */
  let urzadzenia: readonly UrzadzenieDzwieku[] = [];
  /** Zdanie o niepełnym wykazie (nazwy zastępcze, brak sprzętu); puste = wykaz pełny. */
  let powodWykazu = '';
  /** Czy nagrywanie idzie gestem trzymania, czy przełączeniem kliknięciem. */
  let przytrzymanie = false;

  const ster = utworzSterNastawy({
    nastawa: NASTAWA,
    ikona: 'mikrofon',
    wykonaj: (klucz) => {
      wybierz(klucz);
      // Wybór urządzenia i gestu jest lokalny — nie ma komendy, nie ma odmowy
      // rdzenia i nie ma na co czekać, więc obietnica jest spełniona
      // natychmiast.
      return Promise.resolve();
    },
    odswiez: () => odswiez(),
  });

  // Przycisk nagrywania stoi przed uchwytem, w tym samym pudełku steru: jeden
  // komponent paska, dwa sposoby użycia. Najpierw czynność, potem jej nastawy.
  const nagraj = document.createElement('button');
  nagraj.type = 'button';
  nagraj.className = 'dn-btn-ikona dc-ster-mikrofonu__nagraj';
  ster.element.prepend(nagraj);

  /** Odpięcie gestu trzymania; `null`, gdy stoi tryb przełączania kliknięciem. */
  let odepnijGest: (() => void) | null = null;

  function nacisniecie(): void {
    if (dyktowanie.stan() === 'nagrywa') dyktowanie.zakoncz();
    else if (dyktowanie.stan() === 'bezczynne') dyktowanie.rozpocznij();
    // Stan `konczy` jest przelotny i już domyka nagranie — drugie polecenie
    // w tej chwili nie ma czego zacząć ani czego skończyć.
  }

  /**
   * Przestawia sposób obsługi przycisku.
   *
   * Dwa tryby wykluczają się w zdarzeniach, nie w umowie: gest zajmuje
   * `pointerdown`, a przełączanie `click`, który po każdym puszczeniu i tak
   * przychodzi. Trzymanie obu naraz zaczynałoby nagranie i kończyło je tym
   * samym ruchem.
   */
  function podepnijObsluge(): void {
    odepnijGest?.();
    odepnijGest = null;
    nagraj.removeEventListener('click', nacisniecie);
    if (przytrzymanie) odepnijGest = dyktowanie.podepnijPrzycisk(nagraj);
    else nagraj.addEventListener('click', nacisniecie);
  }

  function wybierz(klucz: string): void {
    if (klucz === KLUCZ_PRZYTRZYMANIA) {
      przytrzymanie = !przytrzymanie;
      podepnijObsluge();
      return;
    }
    if (klucz.startsWith(PRZEDROSTEK_URZADZENIA)) {
      dyktowanie.wybierzUrzadzenie(klucz.slice(PRZEDROSTEK_URZADZENIA.length));
    }
  }

  /** Nazwa wskazanego urządzenia — wartość na uchwyt. */
  function nazwaWybranego(): string {
    const wybrane = dyktowanie.urzadzenie();
    if (wybrane === '') return NAZWA_DOMYSLNEGO;
    return urzadzenia.find((u) => u.id === wybrane)?.nazwa ?? wybrane;
  }

  /** Wartość na uchwycie: stan pracy, a poza pracą — wskazany sprzęt. */
  function wartosc(): string {
    const stan = dyktowanie.stan();
    if (stan === 'nagrywa') return WARTOSC_NAGRYWA;
    if (stan === 'konczy') return WARTOSC_KONCZY;
    return nazwaWybranego();
  }

  /**
   * Drzewo mikrofonu — dwa poziomy.
   *
   * Wejście domyślne systemu stoi na wykazie zawsze, także przy pustym odczycie
   * sprzętu: `getUserMedia` bez wskazania działa również wtedy, gdy przeglądarka
   * nie chce pokazać etykiet. Gałąź bez ani jednej pozycji byłaby ślepym
   * zaułkiem, a wykaz bez wejścia domyślnego — kłamstwem o możliwościach.
   */
  function drzewo(): PozycjaMenu[] {
    const wybrane = dyktowanie.urzadzenie();
    const wejscia: PozycjaMenu[] = [
      {
        rodzaj: 'wybor',
        klucz: KLUCZ_DOMYSLNEGO,
        nazwa: NAZWA_DOMYSLNEGO,
        opis: 'Nagrywa tym wejściem, które system ma ustawione jako podstawowe.',
        wybrany: wybrane === '',
      },
      ...urzadzenia.map<PozycjaMenu>((urzadzenie) => ({
        rodzaj: 'wybor',
        klucz: `${PRZEDROSTEK_URZADZENIA}${urzadzenie.id}`,
        nazwa: urzadzenie.nazwa,
        opis: urzadzenie.domyslne
          ? 'Wejście wskazane w systemie jako podstawowe.'
          : 'Wejście dźwięku widziane przez przeglądarkę.',
        wybrany: urzadzenie.id === wybrane,
      })),
    ];

    return [
      {
        rodzaj: 'galaz',
        klucz: GALAZ_URZADZEN,
        nazwa: NASTAWA,
        // Powód niepełnego wykazu wypiera opis rodzajowy: zdanie o zastępczych
        // nazwach albo o braku sprzętu jest tu jedyną rzeczą, której nie widać
        // po samej liście.
        opis: powodWykazu === '' ? `Nagrywa: ${nazwaWybranego()}` : powodWykazu,
        ikona: 'mikrofon',
        dzieci: wejscia,
      },
      {
        rodzaj: 'przelacznik',
        klucz: KLUCZ_PRZYTRZYMANIA,
        nazwa: 'Przytrzymaj, aby nagrać',
        opis: przytrzymanie
          ? 'Nagrywa, dopóki trzymasz przycisk; Escape w trakcie porzuca nagranie.'
          : 'Włączone: nagrywa tylko podczas trzymania. Wyłączone: pierwsze naciśnięcie zaczyna, drugie kończy.',
        wlaczony: przytrzymanie,
      },
    ];
  }

  function odswiez(): void {
    const nagrywa = dyktowanie.stan() === 'nagrywa';
    nagraj.replaceChildren(
      elementIkony(nagrywa ? 'zatrzymaj' : 'mikrofon', { rozmiar: 14 }),
    );
    nagraj.setAttribute('aria-pressed', String(nagrywa));
    // Etykieta przycisku mówi o czynności, a nie o stanie: czytnik ekranu
    // odczytuje ją jako to, co się stanie po naciśnięciu.
    const etykieta = nagrywa
      ? 'Zakończ nagrywanie i rozpoznaj mowę'
      : przytrzymanie
        ? 'Przytrzymaj, aby nagrać'
        : 'Nagraj wypowiedź';
    nagraj.setAttribute('aria-label', etykieta);
    nagraj.title = etykieta;
    ster.ustaw(wartosc(), drzewo());
  }

  /** Odczyt sprzętu — wołany przy złożeniu i po pierwszym udanym rozpoznaniu. */
  function wczytajUrzadzenia(): void {
    void dyktowanie.urzadzenia().then((wykaz) => {
      urzadzenia = wykaz.urzadzenia;
      powodWykazu = wykaz.powod;
      odswiez();
    });
  }

  const odsubskrybujStan = dyktowanie.naStan(() => odswiez());

  const odsubskrybujWynik = dyktowanie.naWynik((wynik) => {
    if (wynik.stan === 'rozpoznano') {
      // Tekst dokładany do pola, nie podmieniający jego treści: dyktowana jest
      // dalsza część zdania zaczętego ręką.
      zaleznosci.wstawTekst(wynik.tekst);
      ster.zdanie('');
      // Nazwy urządzeń odsłaniają się dopiero po pierwszej zgodzie na mikrofon,
      // a ta właśnie padła. Bez tego odczytu zostałoby „Mikrofon 2" na stałe.
      wczytajUrzadzenia();
      return;
    }
    if (wynik.stan === 'bez_mowy') {
      ster.zdanie(ZDANIE_BEZ_MOWY, 'spokojne');
      return;
    }
    ster.zdanie(wynik.powod, 'odmowa');
  });

  podepnijObsluge();
  wczytajUrzadzenia();
  odswiez();

  return {
    element: ster.element,
    odswiez,
    rozlacz() {
      odsubskrybujStan();
      odsubskrybujWynik();
      odepnijGest?.();
      odepnijGest = null;
      nagraj.removeEventListener('click', nacisniecie);
    },
  };
}
