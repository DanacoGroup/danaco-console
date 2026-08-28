import { opisOdmowyBledu } from '../../komponenty/odmowa';
import { poleTekstowe, przycisk, utworzWierszOdpowiedzi } from '../../modele/kontrolki-formularza';
import { nazwaRodzaju, nazwaZasobu } from './karta-zasobu';
import { rozbijEtykiety } from './przyciecie-pol';
import type { StanDesignu } from './stan-designu';
import { wczytajPlik, type WczytanyPlik } from './wczytanie-pliku';

/**
 * Wgranie zasobu do Assets Panel — droga zasobu wnoszonego przez Operatora, osobna od generowania
 * i niewymagająca kanału modelu.
 */
export interface WgranieZasobu {
  element: HTMLElement;
}

export function utworzWgranieZasobu(stan: StanDesignu): WgranieZasobu {
  const pole = document.createElement('input');
  pole.type = 'file';
  pole.className = 'dn-pole-kontrolka md-wgranie__pole';
  // Rodzaje obrazu, bo rodzaj zasobu zna wyłącznie rastr i wektor — kompozycja nie bywa plikiem.
  pole.accept = 'image/*';
  pole.multiple = true;
  pole.id = 'md-wgranie-plik';

  const etykietaPola = document.createElement('label');
  etykietaPola.className = 'dn-pole-etykieta';
  etykietaPola.htmlFor = pole.id;
  etykietaPola.textContent = 'Plik zasobu do wgrania';

  const etykiety = poleTekstowe({
    etykieta: 'Etykiety nadawane przy wgraniu',
    podpowiedz: 'kampania, zima',
    opis:
      'Rozdzielone przecinkiem; pole może zostać puste. Etykiety idą do rdzenia razem ' +
      'z zasobem i od razu zawężają wykaz w filtrze powyżej.',
  });

  // Objaśnienie stoi przy polu, nie w dymku: mówi, co opuszcza przeglądarkę, zanim Operator wskaże plik.
  const objasnienie = document.createElement('p');
  objasnienie.className = 'dn-pole-opis';
  objasnienie.textContent =
    'Treść pliku czyta KLIENT i wysyła ją bajtami do magazynu rdzenia. Ścieżki nie wysyłamy — ' +
    'rdzeń stoi na innej maszynie i otworzyłby pod nią cudzy plik albo żaden. Zasób przestaje ' +
    'zależeć od pliku na dysku: nadpisanie go ani skasowanie nie rusza już tego, co w rdzeniu.';

  const wgraj = przycisk('Wgraj wskazany plik', 'dn-btn dn-btn--sm dn-btn--atrament');
  const odpowiedz = utworzWierszOdpowiedzi();

  // Płyta upuszczania niesie własny opis, nie obramowanie — bez opisu jest dla Operatora niewidoczna.
  const plyta = document.createElement('div');
  plyta.className = 'md-wgranie__plyta';
  plyta.dataset['nad'] = 'nie';
  plyta.textContent = 'Albo upuść tutaj plik obrazu — droga jest ta sama co przycisku.';

  const element = document.createElement('section');
  element.className = 'md-wgranie';
  element.setAttribute('aria-label', 'Wgranie zasobu do Assets Panel');
  element.append(
    etykietaPola,
    pole,
    objasnienie,
    etykiety.element,
    wgraj,
    plyta,
    odpowiedz.element,
  );

  wgraj.addEventListener('click', () => {
    const pliki = [...(pole.files ?? [])];
    if (pliki.length === 0) {
      odpowiedz.pokaz('Wskaż plik w polu powyżej albo upuść go na płytę pod przyciskiem.', false);
      return;
    }
    void wgrajPliki(pliki);
  });

  // Bez blokady zdarzenia przeglądarka otwiera upuszczony plik w karcie, a praca Operatora znika.
  plyta.addEventListener('dragover', (zdarzenie) => {
    zdarzenie.preventDefault();
    plyta.dataset['nad'] = 'tak';
  });
  plyta.addEventListener('dragleave', () => void (plyta.dataset['nad'] = 'nie'));
  plyta.addEventListener('drop', (zdarzenie) => {
    zdarzenie.preventDefault();
    plyta.dataset['nad'] = 'nie';
    const pliki = [...(zdarzenie.dataTransfer?.files ?? [])];
    if (pliki.length === 0) {
      odpowiedz.pokaz('Upuszczenie nie niosło pliku — nie ma czego wgrać.', false);
      return;
    }
    void wgrajPliki(pliki);
  });

  /** Wgrywa wskazane pliki po kolei i składa z wyników jedno sprawozdanie. */
  async function wgrajPliki(pliki: readonly File[]): Promise<void> {
    if (stan.idOkna() === '') {
      odpowiedz.pokaz(`Wgranie zasobu wymaga okna modułu. ${stan.opisOkna()}`, false);
      return;
    }
    const nadawane = rozbijEtykiety(etykiety.kontrolka.value);
    odpowiedz.pokaz(`Wgrywanie ${pliki.length} plików do rdzenia…`, true);
    const weszly: string[] = [];
    const odmowy: string[] = [];
    for (const plik of pliki) {
      const zdanie = await wgrajJeden(plik, nadawane);
      if (zdanie.udane) weszly.push(zdanie.zdanie);
      else odmowy.push(zdanie.zdanie);
    }
    odpowiedz.pokaz(sprawozdanie(weszly, odmowy), odmowy.length === 0);
  }

  /** Jedno wgranie: odczyt pliku, wywołanie, wciągnięcie zasobu do zbioru modułu. */
  async function wgrajJeden(
    plik: File,
    nadawane: readonly string[],
  ): Promise<{ udane: boolean; zdanie: string }> {
    let wczytany: WczytanyPlik;
    try {
      wczytany = await wczytajPlik(plik);
    } catch (blad) {
      // Nieudany odczyt pliku jest usterką przeglądarki, nie odmową rdzenia — zdanie musi to rozróżnić.
      const powod = blad instanceof Error ? blad.message : 'przeglądarka nie podała przyczyny';
      return { udane: false, zdanie: `„${plik.name}": pliku nie udało się odczytać — ${powod}` };
    }
    const wynik = await stan.zrodlo.wgrajZasob({
      idOkna: stan.idOkna(),
      nazwa: wczytany.nazwa,
      rodzaj: wczytany.rodzaj,
      trescBase64: wczytany.trescBase64,
      format: wczytany.format,
      szerokosc: wczytany.szerokosc,
      wysokosc: wczytany.wysokosc,
      etykiety: nadawane,
    });
    if (!wynik.udany || wynik.wynik === undefined) {
      return { udane: false, zdanie: `„${plik.name}": ${opisOdmowyBledu('Wgranie zasobu', wynik.blad)}` };
    }
    // Zasób wchodzi do zbioru modułu tą samą drogą co wynik generowania, bez ponownego odczytu listy.
    const zasob = wynik.wynik.asset;
    stan.wchlon(zasob);
    return {
      udane: true,
      // Zdanie mówi, co oddał rdzeń, nie co okno wysłało; rozmiar pliku dokładamy jako pomiar klienta.
      zdanie:
        `„${nazwaZasobu(zasob)}" (${nazwaRodzaju(zasob.kind)}, ` +
        `wysłano ${wczytany.bajtow} B) — zasób ${zasob.id}`,
    };
  }

  return { element };
}

/**
 * Sprawozdanie z całej partii: co weszło i co zostało odmówione. Odmowy
 * wypisujemy co do jednej — sama liczba nie mówi, który plik przepadł
 * ani z jakiego powodu.
 */
function sprawozdanie(weszly: readonly string[], odmowy: readonly string[]): string {
  const czesci: string[] = [];
  if (weszly.length > 0) czesci.push(`Rdzeń przyjął ${weszly.length}: ${weszly.join(' · ')}.`);
  if (odmowy.length > 0) czesci.push(`Odmówiono ${odmowy.length}: ${odmowy.join(' · ')}.`);
  if (czesci.length === 0) return 'Nie było czego wgrać.';
  return czesci.join(' ');
}
