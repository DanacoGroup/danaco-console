/**
 * Baner „Aktualizuje” to pas w górnej części okna aplikacji, który jednym
 * kliknięciem zakłada dostępne wydanie i uruchamia aplikację ponownie.
 */
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

/** Ustala czas oczekiwania przed pierwszym zapytaniem o dostępne wydania, aby nie konkurować z uruchomieniem aplikacji. */
const ZWLOKA_PIERWSZEGO_PYTANIA = 8_000;

/** Ustala odstęp między kolejnymi zapytaniami o wydania: doba, ponieważ wydania nie powstają częściej, a aplikacja bywa otwarta tygodniami. */
const ODSTEP_PYTANIA = 24 * 60 * 60 * 1000;

/** Ustala częstotliwość odświeżania wyświetlanego upływu czasu w pasie, gdy powłoka nie podaje postępu pobierania. */
const TETNO_LICZNIKA = 1_000;

export interface BanerAktualizacji {
  /** Wpina baner do dokumentu i zaczyna pytać o wydania. */
  uruchom(): void;
  /** Zatrzymuje pytanie i zdejmuje baner. */
  zatrzymaj(): void;
}

/**
 * Buduje baner nad wskazanym elementem, domyślnie nad treścią dokumentu,
 * ponieważ pas należy do okna aplikacji, a nie do pojedynczej trasy.
 */
export function utworzBanerAktualizacji(gospodarz: HTMLElement = document.body): BanerAktualizacji {
  // Obudowa trzyma pas i przycisk wykazu jako rodzeństwo, nie zagnieżdżenie.
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

  /** Zdanie zachęty odpowiada wiedzy o drodze aktualizacji: trzy stany dają trzy różne zdania. */
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
    // Pas wraca do bycia przyciskiem po odmowie, inaczej zostaje trwale wyłączony.
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

  /** Tor postępu powstaje wyłącznie przy znanej całości pobrania, inaczej pokazuje sam licznik. */
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
      // Powłoka milczy, więc pas pokazuje wprost upływ czasu zamiast zmyślonego procentu.
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
    // Praca trwa, więc pas przestaje przyjmować kliknięcia do jej zakończenia.
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
    // Odmowa trwała wyłącza pas, odmowa przemijająca pozostawia możliwość ponowienia.
    pas.disabled = !odmowaPrzemijajaca(kod);
    pas.removeAttribute('aria-busy');
    wynikNaWidoku = true;
  }

  /** Droga niedostępna rozpoznana przed pobraniem otrzymuje inny znak niż odmowa czynności. */
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
    // Obieg pytania pomija czynność trwającą albo już zakończoną, a wynik zostaje na widoku.
    if (zdjety || wTrakcie || zalozone) return;
    const wydanie = await nowszeWydanie(tozsamoscKlienta().wersja);
    if (zdjety) return;
    if (wydanie === null) {
      // Pusty wykaz nie kasuje wyniku poprzedniej czynności widocznego na pasie.
      if (wynikNaWidoku) return;
      obudowa.hidden = true;
      return;
    }
    await ustalDroge();
    // Drugie sprawdzenie po czekaniu chroni przed nadpisaniem stanu podjętej czynności.
    if (zdjety || wTrakcie || zalozone) return;
    proponowane = wydanie;
    pokaz(wydanie);
  }

  async function zaloz(wydanie: Wydanie): Promise<void> {
    // Zapora wejścia dopuszcza jedną aktualizację naraz, bo druga pisałaby w ten sam plik.
    if (wTrakcie || zalozone) return;

    // Droga niemożliwa rozstrzyga się przed pobraniem, zanim ruszy zbędna praca.
    if (droga !== undefined && droga !== null && !droga.mozliwa) {
      pokazDrogeNiedostepna(droga.zdanie);
      return;
    }

    wTrakcie = true;
    pokazPrace(wydanie);

    // Nasłuch działa wyłącznie na czas jednego przebiegu pobierania.
    const nasluch = nasluchujPostepu((postep) => {
      // Zdarzenie spóźnione albo dotyczące zdjętego banera nie przerysowuje pasa.
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
    // Zdanie powodzenia pochodzi od powłoki, która zna stan pliku i rdzenia.
    pokazPrzebieg(wynik.przebieg.zdanie);
  }

  pas.onclick = (): void => {
    if (proponowane === null) return;
    void zaloz(proponowane);
  };

  // Chronologia wydań rozwija się dopiero na żądanie i pyta kanał od nowa.
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
      // Poza powłoką natywną baner nie ma czynności do wykonania po kliknięciu.
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
