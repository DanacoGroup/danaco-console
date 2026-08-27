/** Nagłówek strony głównej wyświetla nadtytuł produktu, nazwę bieżącego widoku i jedno zdanie opisujące rolę strony, w stopniu niższym niż tytuły kart środowisk. */

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
