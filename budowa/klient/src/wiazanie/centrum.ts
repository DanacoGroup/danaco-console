/**
 * Wiązanie Centrum dowodzenia z rdzeniem. Znacznik niesie biblioteka
 * Właściciela — ten plik nic nie buduje: wypełnia wykaz sesji odpowiedzią
 * rdzenia, zakłada sesje na żądanie i wprowadza w okno modułu.
 */

import { Command, type Environment, type Module, type Session } from '../../../shared/contract.ts';
import type { Kanal, Wynik } from '../protokol/kanal.ts';
import { wywolaj } from '../protokol/wywolanie.ts';
import { zwiazWyborModulu } from './wybor-modulu.ts';
import { zwiazOkno } from './okno-modulu.ts';
import { zwiazStudio } from './studio.ts';

/** Kod modułu, którego wnętrze wchodzi do wydania; pozostałe moduły stoją w szynie, lecz okna w tym wydaniu nie mają. */
const KOD_MODULU_WYDANIA = 'studio';

/* Drogi powrotu na stronę główną, które niesie znacznik Właściciela: przycisk
   pasa narzędzi, pozycja menu aplikacji i karta główna okna. */
const POWROT_NA_STRONE_GLOWNA =
  '[aria-label="Centrum dowodzenia"], .dn-karta-widoku--glowna, [data-wyjscie-modulu]';

/** Węzły Centrum, na których wiązanie pracuje. Brak któregokolwiek znaczy, że okno Centrum nie stoi. */
interface WezlyCentrum {
  obszar: HTMLElement;
  wykazSesji: HTMLElement;
}

/** Wiązanie stoi raz na dokument: powłoka może wstawić okno ponownie, a podwójny nasłuch dawałby podwójne sesje. */
let zwiazane = false;

/** Wiąże Centrum z rdzeniem; kanał z obiektu globalnego, bo biblioteka nie jest modułem. Prawda znaczy, że znacznik stał i wiązanie stanęło. */
export function zwiazCentrum(kanal: Kanal | undefined = globalThis.DanacoKanal): boolean {
  if (zwiazane || kanal === undefined) return false;
  const wezly = zbierzWezly();
  if (wezly === null) return false;
  zwiazane = true;

  const wzorWiersza = zdejmijWzorWiersza(wezly.wykazSesji);
  zdejmijTresciPrzykladowe();
  const odswiez = (): void => {
    void odswiezWykaz(kanal, wezly.wykazSesji, wzorWiersza);
  };
  odswiez();
  void wypelnijSrodowiska(kanal, wezly.obszar);
  const katalogModulow = new Map<string, Module>();
  void wczytajModuly(kanal, katalogModulow);
  const katalogSrodowisk = new Map<string, Environment>();
  void wczytajSrodowiska(kanal, katalogSrodowisk);

  wezly.obszar.addEventListener('click', (zdarzenie) => {
    const cel = zdarzenie.target;
    if (!(cel instanceof Element)) return;
    if (cel.closest('[data-okno-nowe]') !== null) {
      void zalozSesje(kanal, wezly.wykazSesji, wzorWiersza);
      return;
    }
  });

  /* Czynności sesji słuchają na dokumencie, nie na wnętrzu okna: biblioteka
     menu przenosi treść menu poza wiersz, więc zdarzenie nie przechodzi przez
     panel, w którym wiersz stoi. */
  document.addEventListener('click', (zdarzenie) => {
    const cel = zdarzenie.target;
    if (!(cel instanceof Element)) return;
    const czynnosc = cel.closest<HTMLElement>('[data-poz-akcja]');
    const wykonaj = CZYNNOSCI_SESJI[czynnosc?.dataset.pozAkcja ?? ''];
    const idSesji = czynnosc?.dataset.idSesji;
    if (wykonaj === undefined || idSesji === undefined) return;
    void wykonaj(kanal, idSesji).then((wynik) => {
      // Odmowa rdzenia wychodzi na wierzch: czynność, która milczy po
      // niepowodzeniu, zostawia Operatora przy wykazie sprzed czynności bez
      // słowa, dlaczego się nie zmienił.
      if (!wynik.udany) oglos('Czynność sesji', wynik.blad?.message ?? 'Rdzeń odmówił wykonania.');
      odswiez();
    });
  });

  /* Wnętrze okna zmienia się w miejscu, więc kolejne wejście podmienia element
     wstawiony poprzednio, nie ten zdjęty przy pierwszym. Szyna stoi poza
     wnętrzem, więc nasłuch obejmuje cały dokument.

     Faza przechwytywania jest konieczna: grot wejścia biblioteki zatrzymuje
     zdarzenie na sobie, więc w fazie bąbelkowania nasłuch nigdy by go nie
     zobaczył. */
  let wnetrze: Element = wezly.obszar;
  // Nazwa środowiska przechodzi z Centrum przez okno wyboru aż do nagłówka
  // modułu; kafel wyboru nie stoi w szynie, więc sam jej nie niesie.
  let nazwaSrodowiska = '';
  document.addEventListener('click', (zdarzenie) => {
    const cel = zdarzenie.target;
    if (!(cel instanceof Element)) return;

    /* Powrót na stronę główną wraca tym samym elementem, który Centrum
       opuściło: jego nasłuchy stoją nietknięte, a wykaz sesji odświeża się
       odpowiedzią rdzenia, bo w module mogły powstać nowe. */
    if (cel.closest(POWROT_NA_STRONE_GLOWNA) !== null) {
      if (wnetrze === wezly.obszar) return;
      wnetrze.replaceWith(wezly.obszar);
      wnetrze = wezly.obszar;
      odswiez();
      return;
    }

    /* Wykaz sesji stoi w panelu bocznym Centrum; osobnego okna rejestru sesji
       to wydanie nie niesie, więc czynność nazywa to wprost. */
    if (cel.closest('[data-otwarz-historie]') !== null) {
      /* Zatrzymanie zdarzenia zdejmuje komunikat biblioteki, który zapowiada
         otwarcie rejestru sesji. Dwa zdania naraz, z których jedno jest
         nieprawdziwe, są gorsze niż milczenie. */
      zdarzenie.stopPropagation();
      oglos('Historia sesji', 'Sesje konta stoją w panelu bocznym Centrum. '
        + 'Osobne okno rejestru sesji nie wchodzi do tego wydania.');
      return;
    }

    const kodNowejSesji = cel.closest<HTMLElement>('[data-nowa-sesja-srodowisko]')?.dataset
      .nowaSesjaSrodowisko;
    if (kodNowejSesji !== undefined) {
      void zalozSesje(kanal, wezly.wykazSesji, wzorWiersza);
      const wstawione = wstawWnetrze(wnetrze, 'dn-tresc-przedsionek');
      if (wstawione === null) return;
      wnetrze = wstawione;
      nazwaSrodowiska = katalogSrodowisk.get(kodNowejSesji)?.name ?? '';
      void zwiazWyborModulu(kanal, kodNowejSesji);
      return;
    }

    /* Środowisko bez modułów w rejestrze nie ma czego rozwinąć. Bez tego zdania
       pozycja szyny odsyłałaby do listy, która nigdy nie stanie. */
    const przelacznik = cel.closest<HTMLElement>('.dn-szyna-poz--srodowisko');
    if (przelacznik !== null) {
      const srodowisko = katalogSrodowisk.get(przelacznik.dataset.srodowisko ?? '');
      if (srodowisko !== undefined && (srodowisko.moduleCodes?.length ?? 0) === 0) {
        zdarzenie.stopPropagation();
        oglos(srodowisko.name, 'Rejestr rdzenia nie wskazuje dla tego środowiska '
          + 'ani jednego modułu, więc lista nie ma czego rozwinąć.');
      }
    }

    const kodSrodowiska = kodSrodowiskaWejscia(cel);
    if (kodSrodowiska !== '') {
      const wstawione = wstawWnetrze(wnetrze, 'dn-tresc-przedsionek');
      if (wstawione === null) return;
      wnetrze = wstawione;
      nazwaSrodowiska = nazwaKartySrodowiska(cel);
      void zwiazWyborModulu(kanal, kodSrodowiska);
      return;
    }

    const kod = kodModulu(cel);
    if (kod === '') return;
    const wstawione = wstawWnetrze(wnetrze, 'dn-tresc-' + kod);
    if (wstawione === null) {
      /* Moduł bez wnętrza nie staje się modułem bieżącym: oznaczenie w szynie
         mówiłoby, że Operator w nim pracuje. */
      zdarzenie.stopPropagation();
      zapowiedzModul(katalogModulow.get(kod));
      return;
    }
    wnetrze = wstawione;
    const srodowisko = nazwaSrodowiskaWejscia(cel) || nazwaSrodowiska;
    if (kod === KOD_MODULU_WYDANIA) zwiazStudio(kanal, srodowisko);
    else zwiazOkno(kanal, kod, srodowisko);
  }, true);

  return true;
}

/** Kod modułu wskazanego pozycją szyny, kaflem Centrum albo kaflem przedsionka; zapis sprowadza do małych liter, bo rejestr rdzenia trzyma kody małymi. */
function kodModulu(cel: Element): string {
  const wskazanie =
    cel.closest('.dn-szyna-poz--modul')?.getAttribute('data-modul') ??
    cel.closest('.pd-kafel')?.getAttribute('data-modul') ??
    cel.closest('.cd-kafel, .dn-kafel--modul')?.getAttribute('data-komponent') ??
    '';
  return wskazanie.toLowerCase();
}

/** Wczytuje rejestr środowisk do spisu po kodzie; nazwa i wykaz modułów rozstrzygają, co pozycja szyny może otworzyć. */
async function wczytajSrodowiska(kanal: Kanal, spis: Map<string, Environment>): Promise<void> {
  const wynik = await wywolaj(kanal, Command.EnvironmentList, { includeModules: true });
  if (!wynik.udany || wynik.wynik === undefined) return;
  for (const srodowisko of wynik.wynik.environments) spis.set(srodowisko.code, srodowisko);
}

/** Wczytuje rejestr modułów do spisu po kodzie; opisy modułów są jedynym źródłem zapowiedzi okna, którego wydanie jeszcze nie niesie. */
async function wczytajModuly(kanal: Kanal, spis: Map<string, Module>): Promise<void> {
  const wynik = await wywolaj(kanal, Command.ModuleList, {});
  if (!wynik.udany || wynik.wynik === undefined) return;
  for (const modul of wynik.wynik.modules) spis.set(modul.code, modul);
}

/**
 * Zapowiada moduł, którego okna to wydanie nie niesie. Komunikat nazywa
 * niegotowość wprost i podaje opis modułu z rejestru rdzenia; twierdzenie, że
 * okno się otwiera, byłoby nieprawdą, a milczenie zostawiłoby pozycję martwą.
 */
function zapowiedzModul(modul: Module | undefined): void {
  if (modul === undefined) return;
  const opis = modul.description ?? '';
  const zdanie = opis === '' ? '' : opis + ' ';
  oglos(modul.name, zdanie + 'Okno tego modułu nie wchodzi do tego wydania.');
}

/** Powiadomienie biblioteki; jej brak zostawia czynność bez komunikatu, bo dorabianie własnego byłoby stawianiem elementu. */
function oglos(tytul: string, tresc: string): void {
  const most = globalThis as { dnToast?: (t: string, o: string, r: string, ms: number) => void };
  most.dnToast?.(tytul, tresc, 'informacja', 4200);
}

/** Nazwa środowiska z karty Centrum, wpisana tam wcześniej rejestrem rdzenia. */
function nazwaKartySrodowiska(cel: Element): string {
  const karta = cel.closest('.dn-karta-srodowiska');
  return karta?.querySelector('.dn-karta-srodowiska-tytul')?.textContent?.trim() ?? '';
}

/** Kod środowiska, do którego prowadzi naciśnięty grot karty Centrum; pustka znaczy, że kliknięcie karty nie dotyczyło. */
function kodSrodowiskaWejscia(cel: Element): string {
  const wejscie = cel.closest('[data-wejdz]');
  if (wejscie === null) return '';
  return wejscie.closest<HTMLElement>('.dn-karta-srodowiska')?.dataset.srodowisko ?? '';
}

/** Zbiera węzły Centrum; brak wykazu sesji znaczy, że w ramie stoi inne okno. */
function zbierzWezly(): WezlyCentrum | null {
  const obszar = document.querySelector<HTMLElement>('.dn-rama-prawa .dn-obszar');
  const wykazSesji = document.getElementById('wykaz-sesji');
  if (obszar === null || wykazSesji === null) return null;
  return { obszar, wykazSesji };
}

/** Zdejmuje wzór wiersza z treści przykładowej; kształt wiersza bierze się ze znacznika, nie z kodu. */
function zdejmijWzorWiersza(wykaz: HTMLElement): HTMLElement | null {
  const wiersz = wykaz.querySelector<HTMLElement>('.dn-panel-wiersz');
  return wiersz === null ? null : (wiersz.cloneNode(true) as HTMLElement);
}

/** Wczytuje wykaz sesji rdzenia i wstawia go w miejsce treści przykładowej. */
async function odswiezWykaz(
  kanal: Kanal,
  wykaz: HTMLElement,
  wzor: HTMLElement | null,
): Promise<void> {
  if (wzor === null) return;
  const wynik = await wywolaj(kanal, Command.SessionList, {});
  if (!wynik.udany || wynik.wynik === undefined) return;
  wykaz.replaceChildren();
  for (const sesja of wynik.wynik.sessions) {
    wykaz.appendChild(zbudujWiersz(wzor, sesja));
  }
}

/**
 * Zwraca klon wzoru wiersza opisany nazwą sesji. Menu czynności zostaje, bo
 * niesie czynności o pokryciu w kontrakcie; jego odwołanie dostaje
 * identyfikator sesji, żeby dwa wiersze nie wskazywały tego samego menu.
 */
function zbudujWiersz(wzor: HTMLElement, sesja: Session): HTMLElement {
  const wiersz = wzor.cloneNode(true) as HTMLElement;
  wiersz.dataset.idSesji = sesja.id;
  const nazwa = wiersz.querySelector('.dn-obszar-pozycja-nazwa');
  // Sesja bez nazwy dostaje nazwany stan pusty: identyfikator jest oznaczeniem magazynu, nie nazwą pracy Operatora.
  if (nazwa !== null) nazwa.textContent = sesja.title ?? 'Sesja bez nazwy';
  wiersz.querySelector('.dn-obszar-pozycja')?.setAttribute('data-id-sesji', sesja.id);
  const menu = wiersz.querySelector('[data-menu-tresc]');
  const wyzwalacz = wiersz.querySelector('[data-menu]');
  if (menu !== null && wyzwalacz !== null) {
    const oznaczenie = 'menu-sesji-' + sesja.id;
    menu.id = oznaczenie;
    wyzwalacz.setAttribute('data-menu', oznaczenie);
  }
  zdejmijCzynnosciBezZrodla(wiersz);
  /* Identyfikator sesji siada na samej pozycji menu, nie tylko na wierszu:
     biblioteka menu przenosi treść menu poza wiersz, więc szukanie sesji
     w przodkach pozycji nic by nie znalazło. */
  for (const pozycja of wiersz.querySelectorAll<HTMLElement>('[data-poz-akcja]')) {
    pozycja.dataset.idSesji = sesja.id;
  }
  return wiersz;
}

/* Czynności menu, dla których kontrakt ma komendę. Pozycje spoza tego spisu
   znikają: pozycja menu, która nic nie robi, jest obietnicą bez pokrycia.

   Spis trzyma wywołania, nie same nazwy komend: każda z tych komend bierze
   wykaz sesji, a nie pojedyncze wskazanie, i tylko wywołanie zapisane przy
   swojej komendzie daje się sprawdzić kontraktem przy budowaniu. */
const CZYNNOSCI_SESJI: Record<string, (kanal: Kanal, idSesji: string) => Promise<Wynik<unknown>>> = {
  archiwizuj: (kanal, idSesji) =>
    wywolaj(kanal, Command.SessionArchive, { sessionIds: [idSesji] }),
  // Potwierdzenie nieodwracalności niesie sama pozycja menu: nazywa usunięcie
  // trwałym, a rejestr sesji drugiego pytania nie stawia.
  usun: (kanal, idSesji) =>
    wywolaj(kanal, Command.SessionDelete, { sessionIds: [idSesji], confirm: true }),
  wyjmij: (kanal, idSesji) =>
    wywolaj(kanal, Command.SessionProjectClear, { sessionIds: [idSesji] }),
};

/** Zdejmuje z menu wiersza pozycje bez komendy w kontrakcie wraz z rozdzielnikami, które po nich zostały. */
function zdejmijCzynnosciBezZrodla(wiersz: HTMLElement): void {
  for (const pozycja of wiersz.querySelectorAll('[data-menu-tresc] .sta-menu-poz')) {
    const czynnosc = pozycja.getAttribute('data-poz-akcja') ?? '';
    if (CZYNNOSCI_SESJI[czynnosc] === undefined) pozycja.remove();
  }
  for (const rozdzielnik of wiersz.querySelectorAll('[data-menu-tresc] .sta-menu-sep')) {
    const przed = rozdzielnik.previousElementSibling;
    if (przed === null || przed.classList.contains('sta-menu-sep')) rozdzielnik.remove();
  }
}

/** Zakłada sesję w rdzeniu i odświeża wykaz jej odpowiedzią. */
async function zalozSesje(
  kanal: Kanal,
  wykaz: HTMLElement,
  wzor: HTMLElement | null,
): Promise<void> {
  const wynik = await wywolaj(kanal, Command.SessionCreate, {});
  if (!wynik.udany) return;
  await odswiezWykaz(kanal, wykaz, wzor);
}

/** Nazwa środowiska, przez które Operator wszedł w moduł; pozycja szyny stoi w grupie środowiska, kafel Centrum nie należy do żadnej. */
function nazwaSrodowiskaWejscia(cel: Element): string {
  const grupa = cel.closest('.dn-szyna-poz--modul')?.closest('.dn-szyna-moduly');
  if (grupa === null || grupa === undefined) return '';
  const przelacznik = document.querySelector(
    `.dn-szyna-poz--srodowisko[aria-controls="${grupa.id}"]`,
  );
  return przelacznik?.getAttribute('aria-label') ?? '';
}

/** Podmienia wnętrze okna na blok ze wskazanego szablonu i oddaje element wstawiony; pustka znaczy, że szablonu nie ma albo jest pusty. */
function wstawWnetrze(stojace: Element, gniazdo: string): Element | null {
  const szablon = document.getElementById(gniazdo);
  if (!(szablon instanceof HTMLTemplateElement)) return null;
  const blok = szablon.content.firstElementChild;
  if (blok === null) return null;
  const wstawione = blok.cloneNode(true) as Element;
  stojace.replaceWith(wstawione);
  return wstawione;
}

/** Zdejmuje treść przykładową bez pokrycia w rdzeniu: karty okien poza główną i komponenty własne. Pusty wykaz odsłania stan pusty ze znacznika. */
function zdejmijTresciPrzykladowe(): void {
  const karty = document.querySelector('.dn-karty-lista');
  if (karty !== null) {
    for (const karta of [...karty.querySelectorAll('.dn-karta-widoku')].slice(1)) karta.remove();
  }
  document.getElementById('cd-wlasne')?.replaceChildren();
  // Odsyłacz do pliku prototypu prowadzi poza produkt i w wydaniu nie stoi.
  document.querySelector('.cd-modul-odnosnik')?.remove();
}

/** Wypełnia karty środowisk rejestrem rdzenia; karta bez pokrycia znika, bo znacznik niesie ich cztery, a rejestr rozstrzyga, ile ich jest. */
async function wypelnijSrodowiska(kanal: Kanal, obszar: HTMLElement): Promise<void> {
  const wynik = await wywolaj(kanal, Command.EnvironmentList, { includeModules: true });
  if (!wynik.udany || wynik.wynik === undefined) return;
  const rejestr = new Map(wynik.wynik.environments.map((s) => [s.code, s]));
  for (const karta of obszar.querySelectorAll<HTMLElement>('.dn-karta-srodowiska')) {
    const srodowisko = rejestr.get(karta.dataset.srodowisko ?? '');
    if (srodowisko === undefined) {
      karta.remove();
      continue;
    }
    opiszSrodowisko(karta, srodowisko);
  }
}

/** Wpisuje w kartę nazwę, opis i liczbę modułów środowiska; stopka z liczbą sesji znika, bo kontrakt nie wiąże sesji ze środowiskiem. */
function opiszSrodowisko(karta: HTMLElement, srodowisko: Environment): void {
  const tytul = karta.querySelector('.dn-karta-srodowiska-tytul');
  if (tytul !== null) tytul.textContent = srodowisko.name;
  const opis = karta.querySelector('.dn-karta-srodowiska-opis');
  if (opis !== null && srodowisko.description !== undefined) {
    opis.textContent = srodowisko.description;
  }
  const miara = karta.querySelector('.dn-karta-srodowiska-motto .cd-metryka-czlon');
  if (miara !== null) miara.textContent = miaraModulow(srodowisko.moduleCodes?.length ?? 0);
  oczyscStopke(karta);
}

/* Stopka karty niesie dwie rzeczy naraz: liczbę sesji środowiska i grot wejścia.
   Kontrakt nie wiąże sesji ze środowiskiem, więc liczba znika, a grot zostaje —
   zdjęcie całej stopki zabrałoby Operatorowi drogę do przedsionka. */
function oczyscStopke(karta: HTMLElement): void {
  const stopka = karta.querySelector('.cd-karta-meta');
  if (stopka === null) return;
  for (const wezel of [...stopka.childNodes]) {
    if (wezel instanceof Element && wezel.closest('[data-wejdz]') !== null) continue;
    wezel.remove();
  }
}

/** Liczba modułów wraz z odmianą rzeczownika; polszczyzna rozróżnia trzy formy, a karta niesie tę miarę zdaniem, nie samą liczbą. */
function miaraModulow(ile: number): string {
  const reszta = ile % 10;
  const setka = ile % 100;
  if (ile === 1) return '1 moduł';
  if (reszta >= 2 && reszta <= 4 && (setka < 12 || setka > 14)) return `${ile} moduły`;
  return `${ile} modułów`;
}
