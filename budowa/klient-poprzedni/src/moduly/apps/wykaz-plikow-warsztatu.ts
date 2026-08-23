import type { DeveloperFile } from '../../../../shared/contract';

/**
 * Wykaz plików warsztatu okna — treść, którą oddaje `apps.workspace.list`.
 *
 * Wykaz pokazuje wszystkie pliki warsztatu okna, więc widać, że wpisana ścieżka
 * należy już do istniejącego pliku, zanim zapis go nadpisze.
 *
 * Wiersz jest przyciskiem: naciśnięcie wstawia ścieżkę i treść pliku do
 * formularza okna. Wiersz nieklikalny kazałby przepisywać ścieżkę ręcznie,
 * a przepisana ścieżka bywa innym plikiem — klucz warsztatu to
 * `(okno, warstwa, ścieżka)` co do znaku.
 *
 * Wykaz niczego nie streszcza: rozmiar i wersja idą tak, jak oddał je rdzeń,
 * a brak pola nazywamy brakiem, zamiast podstawiać zero, które wyglądałoby jak
 * zmierzona wartość.
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

/** Jeden plik warsztatu: przycisk ze ścieżką i zdaniem metadanych. */
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

/** Zdanie metadanych pliku: język, rozmiar, wersja — wyłącznie z pól rdzenia. */
function opiszMetadane(plik: DeveloperFile): string {
  const czesci: string[] = [];
  if (plik.language !== undefined) czesci.push(plik.language);
  czesci.push(
    plik.sizeBytes === undefined ? 'rozmiar nieoddany przez rdzeń' : `${plik.sizeBytes} B`,
  );
  if (plik.versionId !== undefined) czesci.push(`wersja ${plik.versionId}`);
  return czesci.join(' · ');
}
