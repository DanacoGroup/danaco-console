/**
 * Nagłówek strony głównej — nadtytuł produktu, nazwa widoku i zdanie roli.
 *
 * Trójka napisów pochodzi z makiety Centrum dowodzenia
 * (`design/05-okna/przeplyw/03-centrum-dowodzenia.html`) i ze schematu strony
 * w opracowaniu („DANACO CONSOLE — CENTRUM DOWODZENIA" nad strefami wyboru).
 * Znaku marki nagłówek nie powtarza: godło i nazwa produktu stoją już w pasku
 * górnym, a strona ma nad strefami tytuł czytany, nie drugi sygnet.
 *
 * Stopień tytułu jest niższy od stopnia tytułów kart środowisk (20 wobec 24 px)
 * z rozmysłu: środkiem ciężkości strony pozostaje strefa pierwsza, nie napis
 * nad nią.
 */

const PRODUKT = 'Danaco Console';

const NAZWA_WIDOKU = 'Centrum dowodzenia';

const ROLA =
  'Przedpokój przed strefą roboczą. Strona główna nie otwiera przestrzeni ' +
  'roboczej bezpośrednio — prowadzą z niej trzy niezależne ścieżki: wybór ' +
  'środowiska, budowa komponentu własnego oraz ustawienia. Trzy strefy różnią ' +
  'się formą prezentacji celowo: hierarchia skali (karty → kafle → listwa) ' +
  'komunikuje wagę, zanim Operator przeczyta którykolwiek napis.';

export function utworzNaglowekStrony(): HTMLElement {
  const element = document.createElement('header');
  element.className = 'dn-strona__naglowek';

  const nadtytul = document.createElement('p');
  nadtytul.className = 'dn-etykieta-mono dn-strona__nadtytul';
  nadtytul.textContent = PRODUKT;

  const tytul = document.createElement('h1');
  tytul.className = 'dn-strona__tytul';
  tytul.textContent = NAZWA_WIDOKU;

  const rola = document.createElement('p');
  rola.className = 'dn-strona__rola';
  rola.textContent = ROLA;

  element.append(nadtytul, tytul, rola);
  return element;
}
