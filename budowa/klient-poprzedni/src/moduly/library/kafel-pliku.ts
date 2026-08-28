import type { LibraryFile } from '../../../../shared/contract';
import type { TrafienieZnaczenia } from './magazyn-biblioteki';

/**
 * Kafel pliku obsługuje widok siatki i widok galerii biblioteki. Niesie zaznaczenie do
 * czynności zbiorczych, wskazanie pliku czynnego oraz znak rodziny treści zajmujący
 * miejsce przeznaczone na miniaturę.
 */
export interface OpisKafla {
  zaznaczony: boolean;
  czynny: boolean;
  /** Trafienie wskaźnika znaczenia; `null`, gdy plik nie pochodzi z takiego wyszukiwania. */
  znaczenie: TrafienieZnaczenia | null;
  naZaznaczenie(): void;
  naWskazanie(): void;
}

export function utworzKafelPliku(plik: LibraryFile, opis: OpisKafla): HTMLElement {
  const zaznacz = document.createElement('input');
  zaznacz.type = 'checkbox';
  zaznacz.className = 'dn-check ml-kafel__zaznacz';
  zaznacz.checked = opis.zaznaczony;
  zaznacz.setAttribute('aria-label', `Zaznacz plik ${plik.name}`);
  zaznacz.addEventListener('change', () => opis.naZaznaczenie());

  const znak = document.createElement('span');
  znak.className = 'ml-kafel__znak';
  znak.dataset['rodzina'] = rodzinaTresci(plik.mimeType);
  znak.textContent = skrotRodzaju(plik.mimeType);
  znak.setAttribute('aria-hidden', 'true');

  const nazwa = document.createElement('button');
  nazwa.type = 'button';
  nazwa.className = 'ml-kafel__nazwa';
  nazwa.textContent = plik.name;
  nazwa.addEventListener('click', () => opis.naWskazanie());

  const metryka = document.createElement('span');
  metryka.className = 'ml-kafel__metryka';
  metryka.textContent = metrykaKafla(plik);

  const element = document.createElement('li');
  element.className = 'ml-kafel';
  element.dataset['plik'] = plik.id;
  element.dataset['czynny'] = opis.czynny ? 'tak' : 'nie';
  element.append(zaznacz, znak, nazwa, metryka);

  if (plik.tags !== undefined && plik.tags.length > 0) {
    const etykiety = document.createElement('span');
    etykiety.className = 'ml-kafel__etykiety';
    for (const kod of plik.tags) {
      const plakietka = document.createElement('span');
      plakietka.className = 'dn-plakietka';
      plakietka.textContent = kod;
      etykiety.append(plakietka);
    }
    element.append(etykiety);
  }

  const fragment = zdanieZnaczenia(opis.znaczenie);
  if (fragment !== '') {
    const cytat = document.createElement('span');
    cytat.className = 'ml-kafel__znaczenie';
    cytat.textContent = fragment;
    element.append(cytat);
  }

  return element;
}

/**
 * Składa zdanie o tym, dlaczego plik znalazł się w wyniku wyszukiwania po znaczeniu:
 * podaje fragment będący podstawą trafienia, a trafność dopisuje wyłącznie wtedy, gdy
 * podał ją rdzeń.
 */
export function zdanieZnaczenia(znaczenie: TrafienieZnaczenia | null): string {
  if (znaczenie === null) return '';
  const fragment = znaczenie.fragment.trim();
  const skrot = fragment.length > 180 ? `${fragment.slice(0, 180)}…` : fragment;
  const trafnosc = znaczenie.trafnosc === null ? '' : ` (trafność ${znaczenie.trafnosc}/100)`;
  if (skrot === '') return `dopasowanie po znaczeniu${trafnosc} — rdzeń nie podał fragmentu`;
  return `dopasowanie po znaczeniu${trafnosc}: „${skrot}"`;
}

/**
 * Wyznacza rodzinę treści z typu MIME pliku. Rodzina jest nośnikiem znaku kafla, więc
 * typ pusty oraz typ spoza obrazu, dźwięku, obrazu ruchomego i tekstu dostają wartość
 * zastępczą.
 */
export function rodzinaTresci(mimeType: string | undefined): string {
  const rodzaj = (mimeType ?? '').toLowerCase();
  if (rodzaj === '') return 'nieznana';
  const rodzina = rodzaj.split('/')[0] ?? '';
  if (rodzina === 'image' || rodzina === 'audio' || rodzina === 'video' || rodzina === 'text') {
    return rodzina;
  }
  return 'inna';
}

/**
 * Składa skrót rodzaju treści widoczny w miejscu miniatury. Skrót bierze się z podtypu
 * MIME, zapisany jest wielkimi literami i nigdy nie przekracza czterech znaków.
 */
function skrotRodzaju(mimeType: string | undefined): string {
  const rodzaj = (mimeType ?? '').toLowerCase();
  if (rodzaj === '') return '—';
  const podtyp = rodzaj.split('/')[1] ?? '';
  const czlon = podtyp.split(/[+.;]/)[0] ?? '';
  if (czlon === '') return '—';
  return czlon.slice(0, 4).toUpperCase();
}

/**
 * Składa metrykę kafla z typu treści, rozmiaru w bajtach, modułu źródłowego pliku oraz
 * daty ostatniej zmiany; człony nieznane zostają pominięte.
 */
function metrykaKafla(plik: LibraryFile): string {
  const czesci: string[] = [];
  if (plik.mimeType !== undefined && plik.mimeType !== '') czesci.push(plik.mimeType);
  if (plik.sizeBytes !== undefined) czesci.push(`${plik.sizeBytes} B`);
  if (plik.sourceModuleId !== undefined) czesci.push(`z modułu ${plik.sourceModuleId}`);
  czesci.push(new Date(plik.updatedAt).toLocaleDateString('pl'));
  return czesci.join(' · ');
}
