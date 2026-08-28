import { LibraryPreviewKind, type LibraryPreview } from '../../../../shared/contract';

/**
 * Wyrenderowanie podglądu pliku w postaci podanej przez rdzeń. Kontrakt zna
 * cztery rodzaje podglądu: tekst w polu `text`, grafikę oraz strony PDF jako
 * odnośnik `imageRef`, a plik binarny bez treści. Skrócenie podglądu wypisuje
 * zdanie pod treścią.
 */
export function utworzTrescPodgladu(podglad: LibraryPreview): HTMLElement {
  const element = document.createElement('div');
  element.className = 'ml-podglad__tresc';
  element.dataset['rodzaj'] = podglad.kind;

  element.append(cialoPodgladu(podglad));
  if (podglad.truncated === true) {
    element.append(
      zdanie('Podgląd skrócony przez rdzeń — to jest fragment pliku, nie jego całość.'),
    );
  }
  return element;
}

/**
 * Zdanie o stanie podglądu widoczne pod jego treścią: akapit klasy
 * `dn-pole-opis` niosący komunikat o skróceniu podglądu albo o braku treści
 * dla danego rodzaju pliku.
 */
function zdanie(tresc: string): HTMLElement {
  const element = document.createElement('p');
  element.className = 'dn-pole-opis';
  element.textContent = tresc;
  return element;
}

/**
 * Ciało podglądu właściwe dla rodzaju z kontraktu: blok tekstu dla podglądu
 * tekstowego, zdanie zastępcze dla pliku binarnego, a dla grafiki i stron PDF
 * sam odnośnik `imageRef` podany przez rdzeń.
 */
function cialoPodgladu(podglad: LibraryPreview): HTMLElement {
  if (podglad.kind === LibraryPreviewKind.Text) {
    const tekst = document.createElement('pre');
    tekst.className = 'ml-podglad__tekst';
    tekst.textContent = podglad.text ?? '';
    return tekst;
  }
  if (podglad.kind === LibraryPreviewKind.Binary) {
    return zdanie('Plik binarny — rdzeń nie zwraca dla niego treści do pokazania.');
  }
  const odnosnik = (podglad.imageRef ?? '').trim();
  if (odnosnik === '') {
    return zdanie('Rdzeń nie podał odnośnika do podglądu graficznego.');
  }
  const wskazanie = document.createElement('p');
  wskazanie.className = 'ml-podglad__odnosnik';
  wskazanie.textContent = `Podgląd graficzny pod odnośnikiem rdzenia: ${odnosnik}`;
  return wskazanie;
}
