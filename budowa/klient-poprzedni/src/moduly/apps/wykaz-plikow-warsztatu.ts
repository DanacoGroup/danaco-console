import type { DeveloperFile } from '../../../../shared/contract';

/**
 * Wykaz plików warsztatu okna, czyli treść oddana przez komendę wykazu plików
 * warsztatu. Pokazuje wszystkie pliki okna, a każdy wiersz jest przyciskiem,
 * którego naciśnięcie wstawia ścieżkę i treść pliku do formularza okna.
 */
export interface WykazPlikowWarsztatu {
  element: HTMLElement;
  /** Nanosi pliki; oddaje ich liczbę. */
  nanies(pliki: readonly DeveloperFile[]): number;
}

export function utworzWykazPlikowWarsztatu(
  tytul: string,
  naWybor: (plik: DeveloperFile) => void,
): WykazPlikowWarsztatu {
  const naglowek = document.createElement('p');
  naglowek.className = 'mp-wykaz__tytul';
  naglowek.textContent = tytul;

  const lista = document.createElement('ul');
  lista.className = 'mp-wykaz';

  const element = document.createElement('div');
  element.className = 'mp-wykaz__powloka';
  element.append(naglowek, lista);

  return {
    element,
    nanies(pliki) {
      lista.replaceChildren(...pliki.map((plik) => wiersz(plik, naWybor)));
      return pliki.length;
    },
  };
}

/**
 * Buduje wiersz jednego pliku warsztatu jako przycisk, który niesie ścieżkę
 * pliku oraz zdanie metadanych. Naciśnięcie przycisku oddaje cały opis pliku
 * procedurze wyboru podanej przez okno wołające.
 */
function wiersz(plik: DeveloperFile, naWybor: (plik: DeveloperFile) => void): HTMLElement {
  const nazwa = document.createElement('span');
  nazwa.className = 'mp-wykaz__nazwa';
  nazwa.textContent = plik.path;

  const opis = document.createElement('span');
  opis.className = 'mp-wykaz__zaleznosci';
  opis.textContent = opiszMetadane(plik);

  const przycisk = document.createElement('button');
  przycisk.type = 'button';
  przycisk.className = 'dn-btn dn-btn--duch dn-btn--sm';
  przycisk.dataset['plik'] = plik.path;
  przycisk.append(nazwa, opis);
  przycisk.addEventListener('click', () => naWybor(plik));

  const element = document.createElement('li');
  element.className = 'mp-wykaz__wiersz';
  element.append(przycisk);
  return element;
}

/**
 * Składa zdanie metadanych pliku z języka, rozmiaru w bajtach oraz oznaczenia
 * wersji, biorąc każdą część wyłącznie z pól oddanych przez rdzeń. Brak
 * rozmiaru zostaje nazwany brakiem, a nie zastąpiony wartością zerową.
 */
function opiszMetadane(plik: DeveloperFile): string {
  const czesci: string[] = [];
  if (plik.language !== undefined) czesci.push(plik.language);
  czesci.push(
    plik.sizeBytes === undefined ? 'rozmiar nieoddany przez rdzeń' : `${plik.sizeBytes} B`,
  );
  if (plik.versionId !== undefined) czesci.push(`wersja ${plik.versionId}`);
  return czesci.join(' · ');
}
