/**
 * Wiązanie Centrum dowodzenia z rdzeniem. Znacznik niesie biblioteka
 * Właściciela — ten plik nic nie buduje: wypełnia wykaz sesji odpowiedzią
 * rdzenia, zakłada sesje na żądanie i wprowadza w okno modułu.
 */

import { Command, type Environment, type Session } from '../../../shared/contract.ts';
import type { Kanal } from '../protokol/kanal.ts';
import { wywolaj } from '../protokol/wywolanie.ts';
import { zwiazStudio } from './studio.ts';

/** Kod modułu, którego wnętrze wchodzi do wydania; pozostałe moduły stoją w szynie, lecz okna w tym wydaniu nie mają. */
const KOD_MODULU_WYDANIA = 'studio';

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
  void odswiezWykaz(kanal, wezly.wykazSesji, wzorWiersza);
  void wypelnijSrodowiska(kanal, wezly.obszar);

  wezly.obszar.addEventListener('click', (zdarzenie) => {
    const cel = zdarzenie.target;
    if (!(cel instanceof Element)) return;
    if (cel.closest('[data-okno-nowe]') !== null) {
      void zalozSesje(kanal, wezly.wykazSesji, wzorWiersza);
    }
  });

  // Szyna stoi poza obszarem, więc nasłuch modułu obejmuje cały dokument.
  document.addEventListener('click', (zdarzenie) => {
    const cel = zdarzenie.target;
    if (!(cel instanceof Element)) return;
    if (kodModulu(cel) !== KOD_MODULU_WYDANIA) return;
    wejdzWModul(kanal, wezly.obszar, nazwaSrodowiskaWejscia(cel));
  });

  return true;
}

/** Kod modułu wskazanego pozycją szyny albo kaflem; oba zapisy sprowadza do małych liter, bo rejestr rdzenia trzyma kody małymi. */
function kodModulu(cel: Element): string {
  const pozycja = cel.closest('.dn-szyna-poz--modul');
  const wskazanie =
    pozycja !== null
      ? (pozycja.getAttribute('data-modul') ?? '')
      : (cel.closest('.cd-kafel, .dn-kafel--modul')?.getAttribute('data-komponent') ?? '');
  return wskazanie.toLowerCase();
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

/** Zwraca klon wzoru wiersza opisany nazwą sesji; menu czynności klonu zdejmuje się wraz z odwołaniem do nieistniejącej treści. */
function zbudujWiersz(wzor: HTMLElement, sesja: Session): HTMLElement {
  const wiersz = wzor.cloneNode(true) as HTMLElement;
  const nazwa = wiersz.querySelector('.dn-obszar-pozycja-nazwa');
  // Sesja bez nazwy dostaje nazwany stan pusty: identyfikator jest oznaczeniem magazynu, nie nazwą pracy Operatora.
  if (nazwa !== null) nazwa.textContent = sesja.title ?? 'Sesja bez nazwy';
  wiersz.querySelector('[data-menu-tresc]')?.remove();
  wiersz.querySelector('.dn-obszar-pozycja')?.setAttribute('data-id-sesji', sesja.id);
  return wiersz;
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

/** Wprowadza w okno modułu: wnętrze Centrum ustępuje wnętrzu modułu z szablonu, po czym wiązanie modułu obejmuje stojący znacznik. */
function wejdzWModul(kanal: Kanal, obszar: HTMLElement, nazwaSrodowiska: string): void {
  const szablon = document.getElementById('dn-tresc-studio');
  if (!(szablon instanceof HTMLTemplateElement)) return;
  const wnetrze = szablon.content.firstElementChild;
  if (wnetrze === null) return;
  obszar.replaceWith(wnetrze.cloneNode(true));
  zwiazStudio(kanal, nazwaSrodowiska);
}

/** Zdejmuje treść przykładową bez pokrycia w rdzeniu: karty okien poza główną i komponenty własne. Pusty wykaz odsłania stan pusty ze znacznika. */
function zdejmijTresciPrzykladowe(): void {
  const karty = document.querySelector('.dn-karty-lista');
  if (karty !== null) {
    for (const karta of [...karty.querySelectorAll('.dn-karta-widoku')].slice(1)) karta.remove();
  }
  document.getElementById('cd-wlasne')?.replaceChildren();
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
  karta.querySelector('.cd-karta-meta')?.remove();
}

/** Liczba modułów wraz z odmianą rzeczownika; polszczyzna rozróżnia trzy formy, a karta niesie tę miarę zdaniem, nie samą liczbą. */
function miaraModulow(ile: number): string {
  const reszta = ile % 10;
  const setka = ile % 100;
  if (ile === 1) return '1 moduł';
  if (reszta >= 2 && reszta <= 4 && (setka < 12 || setka > 14)) return `${ile} moduły`;
  return `${ile} modułów`;
}
