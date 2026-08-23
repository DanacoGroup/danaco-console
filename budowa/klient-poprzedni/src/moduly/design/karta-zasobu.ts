import { DesignAssetKind, type DesignAsset } from '../../../../shared/contract';

/**
 * Jedna pozycja wykazu zasobów Assets Panel — zasób wraz z jego metadanymi.
 *
 * Karta jest przyciskiem, nie prostokątem: wybór zasobu przestawia naraz panel
 * metadanych, kanwę Design Board i pole obrazu referencyjnego Prompt Buildera,
 * więc musi być osiągalna klawiaturą.
 *
 * Zasób niesie `uri` tylko wtedy, gdy rdzeń go zna; bez niego karta pokazuje
 * pole zastępcze z rodzajem zasobu zamiast miniatury.
 */

/** Polskie nazwy rodzajów zasobu — słownik kontraktu, nie wykaz własny. */
const NAZWY_RODZAJOW: Readonly<Record<DesignAssetKind, string>> = {
  [DesignAssetKind.Image]: 'grafika rastrowa',
  [DesignAssetKind.Vector]: 'grafika wektorowa',
  [DesignAssetKind.Composition]: 'kompozycja tablicy',
  // Wynik pracy modelu wchodzi do Assets Panel tą samą drogą, co zasób wniesiony
  // ręcznie, więc wykaz musi znać także jego rodzaje — inaczej karta pokazałaby
  // pustkę w miejscu nazwy. Słownik pokrywa cały typ kontraktu: brak wartości
  // zatrzymuje kompilację.
  [DesignAssetKind.Document]: 'dokument',
  [DesignAssetKind.Audio]: 'nagranie dźwiękowe',
  [DesignAssetKind.Video]: 'film',
  [DesignAssetKind.Archive]: 'archiwum',
};

export function nazwaRodzaju(rodzaj: DesignAssetKind): string {
  return NAZWY_RODZAJOW[rodzaj];
}

/** Nazwa zasobu widoczna w wykazie; brak nazwy zastępuje identyfikator. */
export function nazwaZasobu(zasob: DesignAsset): string {
  const nazwa = (zasob.name ?? '').trim();
  return nazwa === '' ? zasob.id : nazwa;
}

/**
 * Czy zasób pasuje do frazy zawężającej wykaz — po nazwie albo po etykiecie.
 *
 * Zawężenie jest miejscowe, nie polem żądania: `design.asset.list` zawęża
 * polami `windowId`, `kind`, `tags`, `favoriteOnly` i `limit` (`tor-komendy.ts`),
 * frazy wśród nich nie ma. Predykat stoi tu w jednej kopii, bo czytają go dwa
 * widoki tego samego zbioru — okno Assets Panel i panel `zasoby-designu` stosu
 * paneli pomocniczych; dwie kopie dawałyby na tę samą frazę różne wyniki.
 *
 * @param fraza fraza już przycięta i sprowadzona do małych liter; pustą
 *   („nie zawężaj") rozstrzyga wywołujący, nie ta funkcja.
 */
export function czyPasujeDoFrazy(zasob: DesignAsset, fraza: string): boolean {
  const etykiety = (zasob.tags ?? []).join(' ').toLowerCase();
  return nazwaZasobu(zasob).toLowerCase().includes(fraza) || etykiety.includes(fraza);
}

export function utworzKarteZasobu(
  zasob: DesignAsset,
  czynny: boolean,
  naWybor: (idZasobu: string) => void,
): HTMLElement {
  const podglad = document.createElement('span');
  podglad.className = 'md-zasob__podglad';
  podglad.dataset['rodzaj'] = zasob.kind;
  podglad.textContent = zasob.uri === undefined ? 'bez podglądu' : (zasob.format ?? zasob.kind);

  const tytul = document.createElement('span');
  tytul.className = 'md-zasob__tytul';
  tytul.textContent = nazwaZasobu(zasob);

  const opis = document.createElement('span');
  opis.className = 'md-zasob__opis';
  opis.textContent = opisZasobu(zasob);

  const znaczniki = document.createElement('span');
  znaczniki.className = 'md-zasob__znaczniki';
  znaczniki.append(...plakietkiZasobu(zasob));

  const element = document.createElement('button');
  element.type = 'button';
  element.className = 'dn-karta dn-karta--klikalna md-zasob';
  element.dataset['zasob'] = zasob.id;
  element.dataset['czynny'] = String(czynny);
  element.setAttribute('aria-pressed', String(czynny));
  element.append(podglad, tytul, opis, znaczniki);
  element.addEventListener('click', () => naWybor(zasob.id));
  return element;
}

/** Zdanie metadanych: rodzaj, format, wymiary, chwila powstania. */
export function opisZasobu(zasob: DesignAsset): string {
  const czesci = [nazwaRodzaju(zasob.kind)];
  if (zasob.format !== undefined) czesci.push(zasob.format);
  if (zasob.width !== undefined && zasob.height !== undefined) {
    czesci.push(`${zasob.width}×${zasob.height}`);
  }
  czesci.push(new Date(zasob.createdAt).toLocaleString('pl'));
  return czesci.join(' · ');
}

/** Plakietki zasobu: ulubiony, wariant, etykiety. */
function plakietkiZasobu(zasob: DesignAsset): HTMLElement[] {
  const plakietki: HTMLElement[] = [];
  if (zasob.favorite === true) plakietki.push(plakietka('ulubiony', 'dn-plakietka--sygnal'));
  if (zasob.variantOfAssetId !== undefined) {
    plakietki.push(plakietka(`wariant ${zasob.variantOfAssetId}`, 'dn-plakietka--informacja'));
  }
  // Etykieta zasobu nie niesie stanu, więc idzie plakietką bazową — barwna
  // odmiana jest zarezerwowana dla plakietek znaczących stan.
  for (const etykieta of zasob.tags ?? []) plakietki.push(plakietka(etykieta, ''));
  return plakietki;
}

function plakietka(tresc: string, wariant: string): HTMLElement {
  const element = document.createElement('span');
  element.className = 'dn-plakietka';
  if (wariant !== '') element.classList.add(wariant);
  element.textContent = tresc;
  return element;
}
