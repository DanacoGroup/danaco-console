/* Znacznik przedsionka niesie biblioteka Właściciela; ten plik wypełnia nagłówek
   i kafle odpowiedzią `environment.enter`, a kafel bez pokrycia zdejmuje. */

import { Command, type Environment, type Module } from '../../../shared/contract.ts';
import type { Kanal } from '../protokol/kanal.ts';
import { tozsamoscKlienta } from '../protokol/tozsamosc-klienta.ts';
import { wywolaj } from '../protokol/wywolanie.ts';

interface WezlyWyboru {
  plotno: HTMLElement;
  siatka: HTMLElement;
}

/**
 * Wypełnia przedsionek wskazanego środowiska; prawda znaczy, że znacznik stał
 * i wiązanie stanęło. Wejście idzie jedną komendą, bo `environment.enter`
 * oddaje naraz środowisko, jego moduły i karty sesji w nim otwarte.
 */
export async function zwiazWyborModulu(kanal: Kanal, srodowisko: Environment): Promise<boolean> {
  const wezly = zbierzWezly();
  if (wezly === null) return false;

  const kafleWzorcowe = zdejmijKafleWzorcowe(wezly.siatka);
  const wynik = await wywolaj(kanal, Command.EnvironmentEnter, {
    environmentId: srodowisko.id,
    clientId: tozsamoscKlienta().id,
  });
  if (!wynik.udany || wynik.wynik === undefined) {
    wezly.siatka.replaceChildren();
    return true;
  }
  opiszSrodowisko(wezly.plotno, wynik.wynik.environment);
  wypelnijKafle(wezly.siatka, kafleWzorcowe, wynik.wynik.environment, wynik.wynik.modules);
  return true;
}

function zbierzWezly(): WezlyWyboru | null {
  const plotno = document.querySelector<HTMLElement>('.pd-plotno');
  const siatka = plotno?.querySelector<HTMLElement>('.pd-siatka') ?? null;
  if (plotno === null || siatka === null) return null;
  return { plotno, siatka };
}

/* Kafle prototypu idą do spisu po kodzie modułu: każdy niesie własny znak, więc
   powielenie jednego wzoru dałoby dziewięć kafli o tym samym znaku. */
function zdejmijKafleWzorcowe(siatka: HTMLElement): Map<string, HTMLElement> {
  const spis = new Map<string, HTMLElement>();
  for (const kafel of siatka.querySelectorAll<HTMLElement>('.pd-kafel')) {
    const kod = (kafel.dataset.modul ?? '').toLowerCase();
    if (kod !== '') spis.set(kod, kafel.cloneNode(true) as HTMLElement);
  }
  return spis;
}

/* Motto puste zdejmuje wiersz: katalog wydań stanowi, że pustego motta nie
   zastępuje się tekstem ułożonym po stronie widoku. */
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

/* Kolejność bierze się z wykazu kodów środowiska: to ona rozstrzyga układ
   nawigacji. Moduł bez kafla w prototypie nie staje — kafel bez własnego znaku
   byłby kaflem dorobionym. */
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
    const wzor = wzorcowe.get(kod.toLowerCase());
    if (modul === undefined || wzor === undefined) continue;
    siatka.appendChild(zbudujKafel(wzor, modul));
  }
}

/* Kafel niesie kod modułu, bo po nim rozstrzyga się otwarcie okna. Miara kafla
   schodzi: kontrakt nie wiąże sesji z modułem. */
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
