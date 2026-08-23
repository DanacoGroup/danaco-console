import { WindowRole } from '../../../shared/contract';
import { elementIkony, type NazwaIkony } from '../ikony/ikony';
import { profilModulu } from '../okno-komunikacji/rejestr-profilow';
import { nazwaModuluGniazda, tytulGniazda } from './etykiety-ukladu';
import { numerGniazda, type IdGniazda } from './identyfikatory';
import type { OdczytLacznosci, PortPonawiania } from './lacznosc-okna';
import { utworzPlakietkeLacznosci, type PlakietkaLacznosci } from './plakietka-lacznosci';
import { utworzPlakietkeRoli } from './plakietka-roli';
import { utworzPlakietkeStanu } from './plakietka-stanu';
import type { StanPary } from './stan-pary';
import type { OpisKierunku } from './wiez-koordynacji';

/** Nagłówek gniazda — tożsamość okna, rola, kierunek zlecenia i stan pętli. */
export interface NaglowekGniazda {
  element: HTMLElement;
  /**
   * Gniazdo na sterowanie panelami, przy prawej krawędzi nagłówka.
   *
   * Wybór okien otwartych obok należy do rozmowy, więc sterowanie nimi stoi
   * w prawym górnym rogu okna rozmowy — czyli tutaj.
   *
   * Miejsce jest własne, a nie doklejane do nagłówka, bo do nagłówka dokleja
   * się już przycisk szuflady sterowania (`aplikacja/wiazanie-gniazda.ts` podaje
   * `akcjePaska: gniazdo.naglowek.element`). Dwaj wołający dopisujący na koniec
   * tego samego elementu dawaliby kolejność zależną od tego, kto zdążył
   * pierwszy; sterowanie panelami stoi więc przed tym, co dokleja się na koniec.
   */
  sterowanie: HTMLElement;
  /** Ustawia rolę okna widoczną w nagłówku. */
  ustawRole(rola: WindowRole): void;
  /** Ustawia opis kierunku zlecenia; `null` usuwa wiersz kierunku. */
  ustawKierunek(opis: OpisKierunku | null): void;
  /** Ustawia stan pętli; `null` chowa plakietkę stanu. */
  ustawStan(stan: StanPary | null): void;
  /**
   * Nazywa moduł, w którym to okno pracuje.
   *
   * Etykieta niesie wartość bieżącą, nie napis rodzajowy: sam numer gniazda
   * mówi, które to miejsce na scenie, ale nic o tym, co w nim pracuje. Moduł
   * jest pierwszą rzeczą, jaką okno naprawdę zna — przychodzi z rdzenia
   * zdarzeniem `window.changed` i jest tą samą wartością, po której scena liczy
   * swoją figurę.
   *
   * Model i katalog roboczy tu nie stoją: gniazdo ich nie zna — nie ma ich
   * w opisie okna (`okno-komunikacji/opis-okna.ts`) ani w żadnym zdarzeniu
   * docierającym do układu. Wpisanie czegokolwiek „na oko" byłoby atrapą.
   */
  ustawModul(kod: string): void;
  /** Przyjmuje odczyt łączności; plakietka sama rozstrzyga, czy się pokazać. */
  ustawLacznosc(odczyt: OdczytLacznosci): void;
  /** Zaznacza chwilę przekazania zlecenia — nagłówek błyska. */
  blysnij(): void;
  /** Zatrzymuje odliczanie plakietki łączności; obowiązkowe przy zamykaniu. */
  zamknij(): void;
}

/** Ustawienia nagłówka; każde ma wartość domyślną. */
export interface OpcjeNaglowka {
  /**
   * Dojście do przebiegu ponowienia dla plakietki łączności.
   *
   * Pominięte znaczy, że transport nie wystawia numeru próby — plakietka mówi
   * to wtedy wprost zamiast liczyć próby drugi raz.
   */
  ponowienie?: PortPonawiania;
}

/** Czas widocznego podświetlenia przekazania, w milisekundach. */
const CZAS_BLYSKU = 900;

/**
 * Nagłówek gniazda układu okien równoległych.
 *
 * Rola jest widoczna na pierwszy rzut oka: numer gniazda krojem technicznym,
 * tytuł krojem szeryfowym (nagłówek — nie treść ciągła), plakietka roli,
 * kierunek zlecenia i stan pętli. Nagłówek nie steruje oknem — sterowanie
 * ośmioma ustawieniami należy do kompletu `sterowanie/`.
 */
export function utworzNaglowekGniazda(
  id: IdGniazda,
  rola: WindowRole,
  opcje: OpcjeNaglowka = {},
): NaglowekGniazda {
  const element = document.createElement('header');
  element.className = 'dn-okna__naglowek';
  element.dataset.gniazdo = id;

  const numer = document.createElement('span');
  numer.className = 'dn-okna__numer';
  numer.textContent = String(numerGniazda(id));

  const tytul = document.createElement('h2');
  tytul.className = 'dn-okna__tytul';
  tytul.textContent = tytulGniazda(id);

  // Moduł stoi przy tytule, nie zamiast niego: numer mówi, które to miejsce
  // na scenie, moduł — co w nim pracuje. Obie odpowiedzi są potrzebne naraz.
  const modul = document.createElement('span');
  modul.className = 'dn-okna__modul';
  modul.hidden = true;

  const plakietkaRoli = utworzPlakietkeRoli(rola);

  const kierunek = document.createElement('span');
  kierunek.className = 'dn-okna__kierunek';
  kierunek.hidden = true;

  const rozpychacz = document.createElement('span');
  rozpychacz.className = 'dn-okna__rozpychacz';

  const plakietkaStanu = utworzPlakietkeStanu('brak-pary');
  plakietkaStanu.element.hidden = true;

  // Gniazdo sterowania panelami stoi za rozpychaczem, czyli przy prawej
  // krawędzi, a przed plakietką stanu: stan pętli jest odczytem, sterowanie
  // panelami czynnością, i to czynność ma być bliżej ręki.
  const sterowanie = document.createElement('div');
  sterowanie.className = 'dn-okna__sterowanie';

  // Plakietka łączności stoi przed plakietką stanu pary, bo brak łącza
  // unieważnia sens stanu pętli: „Para gotowa" przy zerwanym łączu jest
  // zdaniem prawdziwym i mylącym naraz.
  const plakietkaLacznosci: PlakietkaLacznosci = utworzPlakietkeLacznosci(
    opcje.ponowienie === undefined ? {} : { ponowienie: opcje.ponowienie },
  );

  element.append(
    numer,
    tytul,
    modul,
    plakietkaRoli.element,
    kierunek,
    rozpychacz,
    sterowanie,
    plakietkaLacznosci.element,
    plakietkaStanu.element,
  );

  let zegar: number | undefined;

  return {
    element,
    sterowanie,

    ustawRole(nowa) {
      element.dataset.rola = nowa;
      plakietkaRoli.ustaw(nowa);
    },

    ustawKierunek(opis) {
      if (opis === null) {
        kierunek.hidden = true;
        kierunek.replaceChildren();
        return;
      }
      kierunek.hidden = false;
      kierunek.replaceChildren(...trescKierunku(opis));
    },

    ustawStan(stan) {
      if (stan === null) {
        plakietkaStanu.element.hidden = true;
        return;
      }
      plakietkaStanu.element.hidden = false;
      plakietkaStanu.ustaw(stan);
    },

    ustawModul(kod) {
      const profil = profilModulu(kod);
      // Kod pusty znaczy „rdzeń jeszcze nie nazwał modułu okna" — wtedy wiersz
      // znika w całości. Napis o module wspólnym byłby podaniem stanu zastanego
      // za rozstrzygnięcie o oknie.
      const nazwany = profil.kod.length > 0;
      modul.hidden = !nazwany;
      modul.textContent = nazwany ? nazwaModuluGniazda(profil.nazwa) : '';
      if (nazwany) element.dataset.modul = profil.kod;
      else delete element.dataset.modul;
      tytul.title = nazwany
        ? `${tytulGniazda(id)} — moduł ${profil.nazwa}`
        : tytulGniazda(id);
    },

    ustawLacznosc: (odczyt) => plakietkaLacznosci.ustaw(odczyt),

    blysnij() {
      element.classList.remove('dn-okna__naglowek--przekazanie');
      // Wymuszenie ponownego przeliczenia układu, żeby powtórne wywołanie
      // uruchomiło animację od nowa zamiast ją pominąć.
      void element.offsetWidth;
      element.classList.add('dn-okna__naglowek--przekazanie');
      window.clearTimeout(zegar);
      zegar = window.setTimeout(
        () => element.classList.remove('dn-okna__naglowek--przekazanie'),
        CZAS_BLYSKU,
      );
    },

    zamknij() {
      window.clearTimeout(zegar);
      plakietkaLacznosci.rozlacz();
    },
  };
}

/** Węzły wiersza kierunku: grot po stronie okna partnera, zwrócony z biegiem zlecenia. */
function trescKierunku(opis: OpisKierunku): Node[] {
  const ikona: NazwaIkony = opis.wPrawo ? 'strzalka-prawo' : 'strzalka-lewo';
  const grot = elementIkony(ikona, { rozmiar: 16 });
  const napis = document.createElement('span');
  napis.textContent = opis.tekst;
  return opis.grotPoPrawej ? [napis, grot] : [grot, napis];
}
