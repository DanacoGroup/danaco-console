import './aktualizacja.css';

import { czyPowlokaNatywna } from '../powloka/powloka-natywna';
import { tozsamoscKlienta } from '../protokol/tozsamosc-klienta';
import {
  drogaAktualizacji,
  nasluchujPostepu,
  odmowaPrzemijajaca,
  wykonajAktualizacje,
  type DrogaAktualizacji,
  type PostepAktualizacji,
} from './most-aktualizacji';
import { nowszeWydanie, type Wydanie } from './wykaz-wydan';
import { utworzWykazWydanWidok } from './wykaz-wydan-widok';

/**
 * Baner „Aktualizuje" — pas w oknie aplikacji, który działa jak przycisk: jedno
 * kliknięcie zakłada wydanie i uruchamia aplikację ponownie. Nie ma okna „czy na
 * pewno" — kliknięcie jest zgodą, a jedynym miejscem, gdzie produkt o cokolwiek
 * pyta, zostaje logowanie.
 *
 * Pas pojawia się wyłącznie wtedy, gdy jest co zakładać. Brak sieci, brak wydań
 * i wydanie nie nowsze od zainstalowanego znaczą to samo: baneru nie ma. Sam
 * z siebie nic nie aktualizuje — dopóki nikt nie kliknie, aplikacja pracuje na
 * wersji zainstalowanej.
 *
 * Odmowa zostaje na widoku: gdy powłoka nie zdoła założyć wydania (zła suma,
 * brak uprawnień, starsza powłoka bez tego polecenia), baner nie znika, tylko
 * zamienia się w zdanie mówiące, co poszło nie tak i co z tym zrobić.
 *
 * Zdanie zachęty rozstrzyga powłoka (`drogaAktualizacji()`), a nie samo
 * rozpoznanie środowiska: na kopii z pakietu `.deb` powłoka odmawia podmiany
 * zaraz po kliknięciu. Gdy powłoka milczy, pas obiecuje tylko to, co wiadomo na
 * pewno — że kliknięcie założy wydanie; o restarcie mówi dopiero zdanie
 * powodzenia, które i tak układa powłoka.
 *
 * Tor postępu rysuje się wyłącznie wtedy, gdy powłoka poda liczbę pobranych
 * bajtów oraz całość. Gdy poda same bajty, jest licznik megabajtów bez toru.
 * Gdy nie poda nic, pas mówi to wprost i pokazuje upływ czasu od kliknięcia —
 * jedyną liczbę, którą zna. Pasek udający procenty byłby atrapą.
 */

/** Jak długo czekać z pierwszym pytaniem o wydania — żeby nie konkurować z uruchomieniem. */
const ZWLOKA_PIERWSZEGO_PYTANIA = 8_000;

/** Co ile pytać ponownie. Doba: wydania nie powstają częściej, a produkt bywa
 *  otwarty tygodniami. */
const ODSTEP_PYTANIA = 24 * 60 * 60 * 1000;

/** Co ile odświeżać upływ czasu w pasie, gdy powłoka postępu nie podaje. */
const TETNO_LICZNIKA = 1_000;

export interface BanerAktualizacji {
  /** Wpina baner do dokumentu i zaczyna pytać o wydania. */
  uruchom(): void;
  /** Zatrzymuje pytanie i zdejmuje baner. */
  zatrzymaj(): void;
}

/**
 * Buduje baner nad gospodarzem dokumentu.
 *
 * `gospodarz` to element, na którego górze pas ma stanąć. Domyślnie `document.body`
 * — baner należy do okna aplikacji, nie do trasy ani modułu, więc nie mieszka
 * w routerze i nie znika przy zmianie widoku.
 */
export function utworzBanerAktualizacji(gospodarz: HTMLElement = document.body): BanerAktualizacji {
  // Obudowa trzyma pas i przycisk wykazu jako rodzeństwo. Zagnieżdżenie
  // przycisku wykazu w banerze byłoby niepoprawnym znacznikiem i pułapką dla
  // czytnika ekranu: jedno kliknięcie trafiałoby w dwa sterowniki naraz.
  const obudowa = document.createElement('div');
  obudowa.className = 'da-pas';
  obudowa.hidden = true;

  const pas = document.createElement('button');
  pas.type = 'button';
  pas.className = 'da-baner';

  const przyciskWykazu = document.createElement('button');
  przyciskWykazu.type = 'button';
  przyciskWykazu.className = 'da-pas__wykaz dn-btn dn-btn--sm dn-btn--duch';
  przyciskWykazu.textContent = 'Wykaz wydań';
  przyciskWykazu.setAttribute('aria-expanded', 'false');

  const wykaz = utworzWykazWydanWidok({ wersjaBiezaca: tozsamoscKlienta().wersja });
  wykaz.element.hidden = true;

  obudowa.append(pas, przyciskWykazu, wykaz.element);

  let budzik: number | undefined;
  let licznik: number | undefined;
  let zdjety = false;
  /** Czy w tej chwili trwa zakładanie wydania (powłoka jeszcze nie odpowiedziała). */
  let wTrakcie = false;
  /** Czy wydanie zostało już założone i czekamy wyłącznie na ponowny start. */
  let zalozone = false;
  /** Czy pas trzyma wynik czynności (powodzenie albo odmowę), a nie zachętę. */
  let wynikNaWidoku = false;
  /** Wydanie, które pas w tej chwili proponuje. */
  let proponowane: Wydanie | null = null;
  /** Odpowiedź powłoki o drodze; `undefined` — jeszcze nie pytano, `null` — powłoka milczy. */
  let droga: DrogaAktualizacji | null | undefined;
  /** Ostatni postęp od powłoki; `null` — powłoka nie powiedziała nic. */
  let ostatniPostep: PostepAktualizacji | null = null;
  /** Chwila kliknięcia — źródło jedynej liczby, którą pas zna na pewno. */
  let poczatekPracy = 0;

  /** Pyta powłokę o drogę jeden raz. Kopia nie zmienia się w trakcie życia okna. */
  async function ustalDroge(): Promise<void> {
    if (droga !== undefined) return;
    droga = await drogaAktualizacji();
  }

  /**
   * Zdanie zachęty — tyle obietnicy, ile jest pokrycia. Trzy stany wiedzy
   * o drodze dają trzy różne zdania; milczenie powłoki nie jest podstawą ani do
   * obietnicy restartu, ani do zapowiedzi niepowodzenia.
   */
  function zachetaDoKlikniecia(): string {
    if (droga === undefined || droga === null) return 'Kliknij, żeby je założyć.';
    if (droga.mozliwa) return 'Kliknij, żeby je założyć — aplikacja uruchomi się ponownie.';
    return 'Kliknij, żeby zobaczyć, jak je założyć — tej kopii powłoka nie podmieni.';
  }

  /** Buduje zawartość pasa od nowa: znak + opis. Jedyne miejsce, które pisze w pasie. */
  function zloz(znakTresc: string, opisTresc: string, klasaZnaku = ''): void {
    pas.textContent = '';
    const znak = document.createElement('span');
    znak.className = klasaZnaku === '' ? 'da-baner__znak' : `da-baner__znak ${klasaZnaku}`;
    znak.textContent = znakTresc;

    const opis = document.createElement('span');
    opis.className = 'da-baner__opis';
    opis.textContent = opisTresc;

    pas.append(znak, opis);
    obudowa.hidden = false;
  }

  function pokaz(wydanie: Wydanie): void {
    zloz('Aktualizuje', `Wydanie ${wydanie.wersja} z ${wydanie.data} jest gotowe. ${zachetaDoKlikniecia()}`);
    // Pas wraca do bycia przyciskiem. Bez tego wiersza `disabled` ustawione
    // w `pokazOdmowe()` zostawałoby do końca życia okna i baner pokazywałby
    // zachętę „Kliknij" przy kliknięciu, które nic nie robi.
    pas.disabled = false;
    pas.removeAttribute('aria-busy');
    wynikNaWidoku = false;
  }

  /** Formatuje bajty jako megabajty w zapisie polskim (przecinek dziesiętny). */
  function megabajty(bajtow: number): string {
    return (bajtow / 1_048_576).toFixed(1).replace('.', ',');
  }

  /** Upływ od kliknięcia w postaci `m:ss`. Jedyna liczba znana bez powłoki. */
  function uplyw(): string {
    const sekundy = Math.max(0, Math.floor((Date.now() - poczatekPracy) / 1000));
    const minuty = Math.floor(sekundy / 60);
    return `${minuty}:${String(sekundy % 60).padStart(2, '0')}`;
  }

  /**
   * Tor postępu; powstaje wyłącznie przy znanej całości pobrania.
   *
   * Bez `calosc` nie ma procentu, a tor bez procentu to pas, którego długość
   * niczego nie oznacza — wtedy idzie sam licznik megabajtów. Znaczniki są
   * `<span>`-ami, bo zawartość mieszka wewnątrz `<button>`, który dopuszcza
   * wyłącznie treść frazową; klasy biblioteki (`dn-postep`) działają tak samo
   * niezależnie od znacznika.
   */
  function torPostepu(postep: PostepAktualizacji): HTMLSpanElement | null {
    if (postep.calosc === null) return null;
    const udzial = Math.min(1, Math.max(0, postep.pobrano / postep.calosc));

    const obszar = document.createElement('span');
    obszar.className = 'dn-postep da-baner__postep';

    const tor = document.createElement('span');
    tor.className = 'dn-postep-tor';
    const wartosc = document.createElement('span');
    wartosc.className = 'dn-postep-wartosc';
    wartosc.style.width = `${(udzial * 100).toFixed(1)}%`;
    tor.append(wartosc);

    const etykieta = document.createElement('span');
    etykieta.className = 'dn-postep-etykieta';
    etykieta.textContent =
      `${megabajty(postep.pobrano)} z ${megabajty(postep.calosc)} MB · ` +
      `${Math.round(udzial * 100)}%`;

    obszar.append(tor, etykieta);
    return obszar;
  }

  /** Przerysowuje pas w trakcie pracy — z danych powłoki albo z ich braku. */
  function rysujPrace(wydanie: Wydanie): void {
    const rozmiar = wydanie.rozmiar ? ` (${wydanie.rozmiar} wg wykazu)` : '';

    if (ostatniPostep === null) {
      // Powłoka milczy, więc pas mówi o tym wprost, zamiast rysować tor bez
      // danych. Upływ czasu widocznie się zmienia, więc okno nie wygląda na
      // zawieszone, a żaden procent nie jest zmyślony.
      zloz(
        'Aktualizuję…',
        `Pobieram wydanie ${wydanie.wersja}${rozmiar}. ` +
          `Powłoka nie podaje postępu — trwa to już ${uplyw()}.`,
      );
      return;
    }

    if (ostatniPostep.etap === 'sprawdzanie-sumy') {
      zloz('Aktualizuję…', `Sprawdzam sumę kontrolną pobranego pliku — trwa to już ${uplyw()}.`);
      return;
    }
    if (ostatniPostep.etap === 'zakladanie') {
      zloz('Aktualizuję…', `Zakładam wydanie ${wydanie.wersja} — trwa to już ${uplyw()}.`);
      return;
    }

    const tor = torPostepu(ostatniPostep);
    if (tor === null) {
      zloz(
        'Aktualizuję…',
        `Pobieram wydanie ${wydanie.wersja} — ${megabajty(ostatniPostep.pobrano)} MB. ` +
          'Serwer nie podał całości, więc udziału nie znam.',
      );
      return;
    }
    zloz('Aktualizuję…', `Pobieram wydanie ${wydanie.wersja}${rozmiar}.`);
    pas.append(tor);
  }

  function pokazPrace(wydanie: Wydanie): void {
    ostatniPostep = null;
    poczatekPracy = Date.now();
    rysujPrace(wydanie);
    // Praca trwa — pas przestaje przyjmować kliknięcia. Dwa kliknięcia pod rząd
    // uruchomiłyby dwa równoległe pobrania piszące w ten sam plik roboczy,
    // każde liczące sumę z własnego strumienia.
    pas.disabled = true;
    pas.setAttribute('aria-busy', 'true');
    zatrzymajLicznik();
    licznik = window.setInterval(() => {
      rysujPrace(wydanie);
    }, TETNO_LICZNIKA);
  }

  function zatrzymajLicznik(): void {
    if (licznik !== undefined) window.clearInterval(licznik);
    licznik = undefined;
  }

  function pokazOdmowe(kod: string, powod: string): void {
    zloz('Aktualizacja nieudana', powod, 'da-baner__znak--odmowa');
    // Przy odmowie trwałej pas przestaje być przyciskiem: powtórzone kliknięcie
    // dałoby tę samą odmowę co do słowa, a Operator ma zamiast tego przeczytać
    // powód. Przy odmowie przemijającej (brak łączności, pobieranie przerwane)
    // ponowienie jest jedyną sensowną czynnością i baner jej nie odbiera.
    pas.disabled = !odmowaPrzemijajaca(kod);
    pas.removeAttribute('aria-busy');
    wynikNaWidoku = true;
  }

  /**
   * Droga niedostępna rozpoznana przed pobraniem — inny znak niż odmowa.
   *
   * „Aktualizacja nieudana" byłoby tu nieprawdą, bo nic się nie zaczęło.
   * Powłoka z góry mówi, że tej kopii nie podmieni (pakiet `.deb` idzie przez
   * `dpkg`), więc Operator dostaje jej zdanie i wie, co zrobić sam. Żaden bajt
   * nie idzie po nic.
   */
  function pokazDrogeNiedostepna(zdanie: string): void {
    zloz('Założysz to ręcznie', zdanie, 'da-baner__znak--odmowa');
    pas.disabled = true;
    pas.removeAttribute('aria-busy');
    wynikNaWidoku = true;
  }

  function pokazPrzebieg(zdanie: string): void {
    zloz('Zaktualizowano', zdanie);
    pas.disabled = true;
    pas.removeAttribute('aria-busy');
    wynikNaWidoku = true;
  }

  async function zapytaj(): Promise<void> {
    // Obieg pytania nie depcze czynności, która trwa. Bez tej zapory budzik
    // nadpisałby napis „Aktualizuję…" zachętą i pas znów przyjmowałby
    // kliknięcia, dokładając kolejne pobranie do tego samego pliku roboczego.
    // Po powodzeniu wynik też zostaje na widoku, bo restart bywa odłożony albo
    // nieudany, a to jedyna informacja Operatora o tym, co się stało.
    if (zdjety || wTrakcie || zalozone) return;
    const wydanie = await nowszeWydanie(tozsamoscKlienta().wersja);
    if (zdjety) return;
    if (wydanie === null) {
      // Pusty wykaz (albo wykaz bez nowszego wydania) znaczy „nie ma co
      // zakładać", a nie „skasuj zdanie o tym, co się przed chwilą stało".
      // Po udanej podmianie z nieudanym restartem to jedyna informacja
      // Operatora o przebiegu.
      if (wynikNaWidoku) return;
      obudowa.hidden = true;
      return;
    }
    await ustalDroge();
    // Drugie sprawdzenie po czekaniu: między pytaniem o wykaz a pytaniem
    // o drogę mija czas sieciowy, a pas jest przez ten czas klikalny, więc
    // zakładanie mogło już ruszyć. Bez tego wiersza `pokaz()` nadpisałoby
    // „Aktualizuję…" zachętą i przywróciło `disabled = false`, zdejmując zaporę
    // dwukliku.
    if (zdjety || wTrakcie || zalozone) return;
    proponowane = wydanie;
    pokaz(wydanie);
  }

  async function zaloz(wydanie: Wydanie): Promise<void> {
    // Zapora wejścia: jedna aktualizacja naraz, bo druga pisałaby w ten sam
    // plik roboczy obok aplikacji. Powłoka ma tę samą zaporę u siebie, bo
    // kliknięcia to nie jedyna droga do polecenia; tutaj stoi po to, żeby
    // Operator w ogóle nie zobaczył drugiego przebiegu.
    if (wTrakcie || zalozone) return;

    // Droga niemożliwa rozstrzyga się przed pobraniem. Powłoka sprawdza to samo
    // u siebie, ale rozstrzygnięcie tutaj oszczędza Operatorowi migającego
    // „Aktualizuję…", po którym natychmiast przychodzi odmowa.
    if (droga !== undefined && droga !== null && !droga.mozliwa) {
      pokazDrogeNiedostepna(droga.zdanie);
      return;
    }

    wTrakcie = true;
    pokazPrace(wydanie);

    // Nasłuch zakładamy na czas jednego przebiegu. Zdarzenie dotyczy tego
    // pobrania i po jego końcu nie ma czego słuchać.
    const nasluch = nasluchujPostepu((postep) => {
      // Zdarzenie spóźnione o przebieg albo przychodzące do banera już zdjętego
      // nie ma czego przerysować — i nie ma prawa wskrzesić zdjętego pasa.
      if (!wTrakcie || zdjety) return;
      ostatniPostep = postep;
      rysujPrace(wydanie);
    });

    let wynik;
    try {
      wynik = await wykonajAktualizacje(wydanie);
    } finally {
      wTrakcie = false;
      zatrzymajLicznik();
      void nasluch.then((odwolaj) => {
        odwolaj();
      });
    }
    if (!wynik.udana) {
      pokazOdmowe(wynik.kod, wynik.powod);
      return;
    }
    zalozone = true;
    // Powłoka odkłada restart o `restart_za_ms`, więc jest chwila na pokazanie
    // zdania powodzenia. Układa je powłoka — ona jedna wie, co się stało
    // z plikiem i z rdzeniem w tle.
    pokazPrzebieg(wynik.przebieg.zdanie);
  }

  pas.onclick = (): void => {
    if (proponowane === null) return;
    void zaloz(proponowane);
  };

  // Ujawnianie stopniowe: chronologia wydań rozwija się dopiero, gdy Operator
  // o nią poprosi, i dopiero wtedy pyta kanał. Zwinięcie jej nie kasuje —
  // kolejne przywołanie pokazuje to, co już wiadomo, i pyta na nowo.
  przyciskWykazu.onclick = (): void => {
    const rozwiniety = przyciskWykazu.getAttribute('aria-expanded') === 'true';
    if (rozwiniety) {
      przyciskWykazu.setAttribute('aria-expanded', 'false');
      wykaz.element.hidden = true;
      return;
    }
    przyciskWykazu.setAttribute('aria-expanded', 'true');
    wykaz.element.hidden = false;
    void wykaz.odswiez();
  };

  return {
    uruchom(): void {
      // Poza powłoką natywną (interfejs otwarty w przeglądarce, podgląd
      // dewelopera) baner nie miałby czego zrobić po kliknięciu: podmiana pliku
      // i restart są własnością powłoki. Pytanie kanału wydań byłoby wtedy
      // ruchem w sieć po nic, a jego niepowodzenie ląduje w dzienniku konsoli
      // jako błąd zasobu i wygląda jak usterka produktu.
      //
      // Sprawdzamy wprost obecność powłoki, a nie wykonalność podmiany — o tę
      // drugą pyta się powłoki przez `drogaAktualizacji()`.
      if (!czyPowlokaNatywna()) return;
      gospodarz.prepend(obudowa);
      budzik = window.setTimeout(function pytaj() {
        void zapytaj();
        budzik = window.setTimeout(pytaj, ODSTEP_PYTANIA);
      }, ZWLOKA_PIERWSZEGO_PYTANIA);
    },
    zatrzymaj(): void {
      zdjety = true;
      if (budzik !== undefined) window.clearTimeout(budzik);
      zatrzymajLicznik();
      obudowa.remove();
    },
  };
}
