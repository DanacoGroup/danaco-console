import type { LibraryFile } from '../../../../shared/contract';
import type { StanTresci } from './dostepnosc-tresci';
import { zdanieZnaczenia } from './kafel-pliku';
import type { TrafienieZnaczenia } from './magazyn-biblioteki';

/**
 * Jeden wiersz wykazu plików repozytorium.
 *
 * Wiersz niesie trzy rzeczy naraz: zaznaczenie do czynności zbiorczych,
 * wskazanie pliku czynnego dla pozostałych trzech okien i metrykę pliku.
 *
 * O braku treści wiersz mówi wyłącznie wtedy, gdy rdzeń sam ją orzekł
 * (`dostepnosc-tresci.ts`, werdykt `brak`), i powtarza jego powód. Plik, o który
 * nikt jeszcze nie pytał, nie dostaje żadnego znaku.
 *
 * Odmowa podglądu nie jest orzeczeniem o treści: rdzeń odmawia kodem
 * `not_found`, gdy pliku nie ma, i kodem `internal_error`, gdy nośnik nie oddał
 * treści spod odwołania. Oba są zdaniem o pliku i o nośniku, nie o zawartości
 * repozytorium, więc odmowa dostaje własne zdanie — powód rdzenia bez
 * dopisanego zarzutu.
 *
 * Metryka nie jest świadkiem treści: brak `versionId` czy `checksum` niczego
 * o niej nie orzeka i w drugą stronę tak samo. Świadkiem jest odpowiedź rdzenia.
 */
export interface OpisWiersza {
  zaznaczony: boolean;
  czynny: boolean;
  /** Odpowiedź rdzenia o treści tego pliku; nie domysł wiersza. */
  tresc: StanTresci;
  /**
   * Trafienie wskaźnika znaczenia wraz z fragmentem, na którym się oparło;
   * `null`, gdy plik nie pochodzi z wyszukiwania po znaczeniu.
   */
  znaczenie: TrafienieZnaczenia | null;
  naZaznaczenie(): void;
  naWskazanie(): void;
}

export function utworzWierszPliku(plik: LibraryFile, opis: OpisWiersza): HTMLElement {
  const zaznacz = document.createElement('input');
  zaznacz.type = 'checkbox';
  zaznacz.className = 'dn-check';
  zaznacz.checked = opis.zaznaczony;
  zaznacz.setAttribute('aria-label', `Zaznacz plik ${plik.name}`);
  zaznacz.addEventListener('change', () => opis.naZaznaczenie());

  const nazwa = document.createElement('button');
  nazwa.type = 'button';
  nazwa.className = 'ml-plik__nazwa';
  nazwa.textContent = plik.name;
  nazwa.addEventListener('click', () => opis.naWskazanie());

  const metryka = document.createElement('span');
  metryka.className = 'ml-plik__metryka';
  metryka.textContent = opisMetryki(plik);

  const element = document.createElement('li');
  element.className = 'ml-plik';
  element.dataset['plik'] = plik.id;
  element.dataset['czynny'] = opis.czynny ? 'tak' : 'nie';
  element.dataset['tresc'] = opis.tresc.werdykt;
  element.append(zaznacz, nazwa, metryka);

  const zdanie = zdanieOTresci(opis.tresc);
  if (zdanie !== '') {
    const powod = document.createElement('span');
    powod.className = 'ml-plik__powod';
    powod.textContent = zdanie;
    element.append(powod);
  }

  // Fragment, na którym oparło się dopasowanie po znaczeniu, stoi przy wierszu,
  // bo bez niego trafność jest liczbą bez podstawy do sprawdzenia.
  const fragment = zdanieZnaczenia(opis.znaczenie);
  if (fragment !== '') {
    const cytat = document.createElement('span');
    cytat.className = 'ml-plik__znaczenie';
    cytat.textContent = fragment;
    element.append(cytat);
  }

  return element;
}

/**
 * Zdanie wiersza o treści pliku — jedno na werdykt, żadne bez werdyktu.
 *
 * `brak` jest jedynym werdyktem, po którym wolno powiedzieć „bez treści
 * w repozytorium": tylko wtedy rdzeń o treści orzekł. `odmowa` powtarza sam
 * powód i mówi, że treść pozostaje nieznana. `odwolanie` mówi, że rdzeń wskazał
 * miejsce treści, a jej samej nie podał. `nieznana` i `osiagalna` nie dają
 * zdania: pierwsza nie ma czego powiedzieć, druga niczego nie zarzuca.
 */
function zdanieOTresci(tresc: StanTresci): string {
  if (tresc.werdykt === 'brak') return `bez treści w repozytorium — ${tresc.powod}`;
  if (tresc.werdykt === 'odmowa') return `treść nieznana — ${tresc.powod}`;
  if (tresc.werdykt === 'odwolanie') return `treści nie widziałem — ${tresc.powod}`;
  return '';
}

/** Metryka wiersza: ścieżka, moduł wytwórcy, rozmiar i etykiety. */
function opisMetryki(plik: LibraryFile): string {
  const czesci: string[] = [];
  if (plik.path !== undefined && plik.path !== '') czesci.push(plik.path);
  if (plik.sourceModuleId !== undefined) czesci.push(`z modułu ${plik.sourceModuleId}`);
  if (plik.sizeBytes !== undefined) czesci.push(`${plik.sizeBytes} B`);
  const etykiety = plik.tags ?? [];
  if (etykiety.length > 0) czesci.push(etykiety.join(' · '));
  return czesci.join(' — ');
}
