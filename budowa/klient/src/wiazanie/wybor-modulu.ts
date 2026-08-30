/**
 * Wiązanie ekranu wyboru modułu środowiska z rdzeniem. Znacznik niesie
 * biblioteka Właściciela; ten plik wypełnia nagłówek i kafle rejestrem rdzenia,
 * a kafel bez pokrycia zdejmuje.
 */

import { Command, type Environment, type Module } from '../../../shared/contract.ts';
import type { Kanal } from '../protokol/kanal.ts';
import { wywolaj } from '../protokol/wywolanie.ts';

/** Węzły ekranu wyboru; brak płótna znaczy, że w ramie stoi inne okno. */
interface WezlyWyboru {
  plotno: HTMLElement;
  siatka: HTMLElement;
}

/** Wypełnia ekran wyboru modułu wskazanego środowiska; prawda znaczy, że znacznik stał i wiązanie stanęło. */
export async function zwiazWyborModulu(kanal: Kanal, kodSrodowiska: string): Promise<boolean> {
  const wezly = zbierzWezly();
  if (wezly === null) return false;

  const kafleWzorcowe = zdejmijKafleWzorcowe(wezly.siatka);
  zdejmijMiaryBezZrodla(wezly.plotno);

  const srodowiska = await wywolaj(kanal, Command.EnvironmentList, { includeModules: true });
  if (!srodowiska.udany || srodowiska.wynik === undefined) return true;
  const srodowisko = srodowiska.wynik.environments.find((s) => s.code === kodSrodowiska);
  if (srodowisko === undefined) return true;
  opiszSrodowisko(wezly.plotno, srodowisko);

  const moduly = await wywolaj(kanal, Command.ModuleList, {});
  if (!moduly.udany || moduly.wynik === undefined) return true;
  wypelnijKafle(wezly.siatka, kafleWzorcowe, srodowisko, moduly.wynik.modules);
  return true;
}

/** Zbiera węzły ekranu wyboru. */
function zbierzWezly(): WezlyWyboru | null {
  const plotno = document.querySelector<HTMLElement>('.pd-plotno');
  const siatka = plotno?.querySelector<HTMLElement>('.pd-siatka') ?? null;
  if (plotno === null || siatka === null) return null;
  return { plotno, siatka };
}

/** Zdejmuje kafle prototypu w spis po kodzie modułu; każdy niesie własną ikonę, więc powielanie jednego wzoru dałoby dziewięć kafli o tym samym znaku. */
function zdejmijKafleWzorcowe(siatka: HTMLElement): Map<string, HTMLElement> {
  const spis = new Map<string, HTMLElement>();
  for (const kafel of siatka.querySelectorAll<HTMLElement>('.pd-kafel')) {
    const kod = (kafel.dataset.modul ?? '').toLowerCase();
    if (kod !== '') spis.set(kod, kafel.cloneNode(true) as HTMLElement);
  }
  return spis;
}

/** Wpisuje w nagłówek nazwę, motto i opis środowiska; motto puste zdejmuje wiersz, bo katalog wydań stanowi, że pustego motta nie zastępuje się tekstem widoku. */
function opiszSrodowisko(plotno: HTMLElement, srodowisko: Environment): void {
  const tytul = plotno.querySelector('.pd-tytul');
  if (tytul !== null) tytul.textContent = srodowisko.name;
  const motto = plotno.querySelector('.pd-motto');
  if (motto !== null) {
    if (srodowisko.motto === undefined || srodowisko.motto === '') motto.remove();
    else motto.textContent = srodowisko.motto;
  }
  const opis = plotno.querySelector('.pd-opis');
  if (opis !== null && srodowisko.description !== undefined) {
    opis.textContent = srodowisko.description;
  }
}

/**
 * Stawia kafel na każdy moduł środowiska. Kolejność bierze się z wykazu kodów
 * środowiska, bo to ona rozstrzyga układ nawigacji. Moduł, dla którego prototyp
 * nie ma kafla, nie staje: kafel bez własnego znaku byłby kaflem dorobionym.
 */
function wypelnijKafle(
  siatka: HTMLElement,
  wzorcowe: Map<string, HTMLElement>,
  srodowisko: Environment,
  moduly: Module[],
): void {
  const katalog = new Map(moduly.map((modul) => [modul.code, modul]));
  siatka.replaceChildren();
  for (const kod of srodowisko.moduleCodes ?? []) {
    const modul = katalog.get(kod);
    const wzor = wzorcowe.get(kod);
    if (modul === undefined || wzor === undefined) continue;
    siatka.appendChild(zbudujKafel(wzor, modul));
  }
}

/** Zwraca klon kafla prototypu opisany nazwą i opisem modułu; kafel niesie kod modułu, bo po nim rozstrzyga się otwarcie okna. */
function zbudujKafel(wzor: HTMLElement, modul: Module): HTMLElement {
  const kafel = wzor.cloneNode(true) as HTMLElement;
  kafel.dataset.modul = modul.code;
  const nazwa = kafel.querySelector('.pd-kafel-nazwa');
  if (nazwa !== null) nazwa.textContent = modul.name;
  const opis = kafel.querySelector('.pd-kafel-opis');
  if (opis !== null) {
    if (modul.description === undefined) opis.remove();
    else opis.textContent = modul.description;
  }
  kafel.querySelector('.pd-kafel-meta')?.remove();
  return kafel;
}

/* Miary bez pokrycia w kontrakcie: liczba sesji kafla oraz licznik projektów
   i sesji listwy. Kontrakt nie wiąże sesji ze środowiskiem ani z modułem, więc
   liczba stojąca w prototypie jest treścią przykładową, nie miarą. */
function zdejmijMiaryBezZrodla(plotno: HTMLElement): void {
  plotno.querySelector('.pd-listwa-meta')?.remove();
}
