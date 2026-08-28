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

/** Nagłówek gniazda — tożsamość okna, rola, kierunek zlecenia i stan pętli — widoczne jednym spojrzeniem w każdym gnieździe sceny. */
export interface NaglowekGniazda {
  element: HTMLElement;
  /** Gniazdo na sterowanie panelami przy prawej krawędzi nagłówka, osobne od przycisku szuflady. */
  sterowanie: HTMLElement;
  /** Ustawia rolę okna widoczną w nagłówku. */
  ustawRole(rola: WindowRole): void;
  /** Ustawia opis kierunku zlecenia; `null` usuwa wiersz kierunku. */
  ustawKierunek(opis: OpisKierunku | null): void;
  /** Ustawia stan pętli; `null` chowa plakietkę stanu. */
  ustawStan(stan: StanPary | null): void;
  /** Nazywa moduł, w którym to okno pracuje — pierwsza rzecz, jaką okno naprawdę zna. */
  ustawModul(kod: string): void;
  /** Przyjmuje odczyt łączności; plakietka sama rozstrzyga, czy się pokazać. */
  ustawLacznosc(odczyt: OdczytLacznosci): void;
  /** Zaznacza chwilę przekazania zlecenia — nagłówek błyska. */
  blysnij(): void;
  /** Zatrzymuje odliczanie plakietki łączności; obowiązkowe przy zamykaniu. */
  zamknij(): void;
}

/** Ustawienia nagłówka gniazda; każde pole ma wartość domyślną stosowaną, gdy gospodarz jej nie poda wprost. */
export interface OpcjeNaglowka {
  /** Dojście do przebiegu ponowienia dla plakietki łączności; pominięte znaczy brak numeru próby. */
  ponowienie?: PortPonawiania;
}

/** Czas widocznego podświetlenia przekazania zlecenia w nagłówku gniazda, wyrażony w milisekundach czasu. */
const CZAS_BLYSKU = 900;

/**
 * Nagłówek gniazda układu okien równoległych pokazuje jednym spojrzeniem numer gniazda, tytuł, plakietkę roli, kierunek zlecenia i stan pętli, nie sterując przy tym samym oknem.
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

  // Moduł stoi przy tytule, nie zamiast niego: numer mówi, które to miejsce, moduł — co w nim pracuje.
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

  // Sterowanie panelami stoi bliżej ręki niż plakietka stanu, bo jest czynnością, a stan tylko odczytem.
  const sterowanie = document.createElement('div');
  sterowanie.className = 'dn-okna__sterowanie';

  // Plakietka łączności stoi przed plakietką stanu pary, bo brak łącza unieważnia sens stanu pętli.
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
      // Kod pusty znaczy, że rdzeń jeszcze nie nazwał modułu okna, więc wiersz o module znika w całości.
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
      // Wymuszenie ponownego przeliczenia układu, żeby powtórne wywołanie uruchomiło animację od nowa.
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

/** Węzły wiersza kierunku: grot po stronie okna partnera, zwrócony zgodnie z kierunkiem biegu zlecenia. */
function trescKierunku(opis: OpisKierunku): Node[] {
  const ikona: NazwaIkony = opis.wPrawo ? 'strzalka-prawo' : 'strzalka-lewo';
  const grot = elementIkony(ikona, { rozmiar: 16 });
  const napis = document.createElement('span');
  napis.textContent = opis.tekst;
  return opis.grotPoPrawej ? [napis, grot] : [grot, napis];
}
