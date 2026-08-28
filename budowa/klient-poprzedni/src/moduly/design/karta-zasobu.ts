import { DesignAssetKind, type DesignAsset } from '../../../../shared/contract';

/**
 * Karta zasobu jest jedną pozycją wykazu zasobów Assets Panel i niesie zasób
 * wraz z jego metadanymi. Moduł podaje także polskie nazwy rodzajów zasobu
 * wzięte ze słownika kontraktu, a nie z wykazu własnego.
 */
const NAZWY_RODZAJOW: Readonly<Record<DesignAssetKind, string>> = {
  [DesignAssetKind.Image]: 'grafika rastrowa',
  [DesignAssetKind.Vector]: 'grafika wektorowa',
  [DesignAssetKind.Composition]: 'kompozycja tablicy',
  // Wynik pracy modelu wchodzi do wykazu tą samą drogą, co zasób wniesiony ręcznie.
  [DesignAssetKind.Document]: 'dokument',
  [DesignAssetKind.Audio]: 'nagranie dźwiękowe',
  [DesignAssetKind.Video]: 'film',
  [DesignAssetKind.Archive]: 'archiwum',
};

export function nazwaRodzaju(rodzaj: DesignAssetKind): string {
  return NAZWY_RODZAJOW[rodzaj];
}

/**
 * Zwraca nazwę zasobu widoczną w wykazie. Gdy zasób nie ma nazwy albo nazwa
 * składa się z samych odstępów, w jej miejsce wchodzi identyfikator zasobu.
 */
export function nazwaZasobu(zasob: DesignAsset): string {
  const nazwa = (zasob.name ?? '').trim();
  return nazwa === '' ? zasob.id : nazwa;
}

/**
 * Rozstrzyga, czy zasób pasuje do frazy zawężającej wykaz, porównując frazę
 * z nazwą zasobu oraz z jego etykietami. Fraza przychodzi przycięta i sprowadzona
 * do małych liter, a znaczenie frazy pustej rozstrzyga wywołujący.
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

/**
 * Składa zdanie metadanych zasobu z rodzaju, formatu, wymiarów oraz chwili
 * powstania. Człony rozdziela kropka środkowa, a człony nieznane wypadają
 * ze zdania zamiast zostawiać puste miejsce.
 */
export function opisZasobu(zasob: DesignAsset): string {
  const czesci = [nazwaRodzaju(zasob.kind)];
  if (zasob.format !== undefined) czesci.push(zasob.format);
  if (zasob.width !== undefined && zasob.height !== undefined) {
    czesci.push(`${zasob.width}×${zasob.height}`);
  }
  czesci.push(new Date(zasob.createdAt).toLocaleString('pl'));
  return czesci.join(' · ');
}

/**
 * Składa plakietki zasobu w stałej kolejności: znacznik ulubionego, wskazanie
 * zasobu źródłowego dla wariantu oraz etykiety nadane zasobowi. Kolejność jest
 * stała, więc karty w wykazie czyta się tak samo.
 */
function plakietkiZasobu(zasob: DesignAsset): HTMLElement[] {
  const plakietki: HTMLElement[] = [];
  if (zasob.favorite === true) plakietki.push(plakietka('ulubiony', 'dn-plakietka--sygnal'));
  if (zasob.variantOfAssetId !== undefined) {
    plakietki.push(plakietka(`wariant ${zasob.variantOfAssetId}`, 'dn-plakietka--informacja'));
  }
  // Etykieta zasobu nie niesie stanu, więc idzie plakietką bazową.
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
