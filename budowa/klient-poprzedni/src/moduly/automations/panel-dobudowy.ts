/**
 * Wspólny kształt paneli dopełniających okna modułu Automations.
 *
 * Cztery panele (wersje, nadzór harmonogramu, zlecenia kolejki, dozór
 * przebiegu) mają tę samą budowę: tytuł, zdanie o przeznaczeniu, garść pól
 * i pas przycisków, a pod nimi jeden nośnik stanu treści. Zamiast czterech razy
 * tego samego rusztowania stoi tu jedno.
 *
 * Panel nie jest oknem operacyjnym i nie udaje nim być. Okna modułu wylicza
 * `kody-okien.ts` i jest ich pięć; panel osadza się WEWNĄTRZ okna, do którego
 * należy jego praca — wersje w Workflow Builderze, zlecenia w Queue Managerze.
 * Szósty kafel na siatce byłby szóstym oknem, którego dokument projektowy nie
 * zna.
 *
 * Odpowiedź rdzenia pokazuje się w całości, zapisem strukturalnym. Panel jest
 * powierzchnią roboczą Operatora nad czynnościami, których kontrakt oddaje
 * bardzo różne kształty — wykaz wersji, różnicę pól, ładunek kroku, punkty
 * wznowienia. Rysunek zmyślony osobno dla każdej z nich pokazywałby MNIEJ, niż
 * rdzeń oddał, a to jest gorsze niż surowy zapis: Operator ma widzieć odpowiedź,
 * a nie jej streszczenie napisane przez okno.
 */
import { przyciskAkcji, wiersz } from '../../modele/kontrolki-formularza';
import type { Wynik } from '../../protokol/kanal';
import { utworzStanTresci, type StanTresci } from './stany-okna';

/** Panel dopełniający jedno okno operacyjne. */
export interface PanelDobudowy {
  /** Element osadzany wewnątrz okna. */
  element: HTMLElement;
  /** Nośnik stanu treści panelu — ładowanie, odmowa, potwierdzenie. */
  tresc: StanTresci;
  /** Dokłada wiersz pola do siatki panelu i oddaje kontrolkę. */
  dodajPole<T extends HTMLElement>(etykieta: string, kontrolka: T, objasnienie?: string): T;
  /** Dokłada przycisk czynności do pasa akcji panelu. */
  dodajCzynnosc(etykieta: string, czynnosc: () => void): HTMLButtonElement;
}

/** Składa panel o podanym tytule i przeznaczeniu. */
export function utworzPanelDobudowy(tytul: string, przeznaczenie: string): PanelDobudowy {
  const naglowek = document.createElement('h4');
  naglowek.className = 'da-panel__tytul';
  naglowek.textContent = tytul;

  const opis = document.createElement('p');
  opis.className = 'da-panel__opis';
  opis.textContent = przeznaczenie;

  const pola = document.createElement('div');
  pola.className = 'da-panel__pola';

  const akcje = document.createElement('div');
  akcje.className = 'da-panel__akcje';

  const tresc = utworzStanTresci();

  const element = document.createElement('section');
  element.className = 'da-panel';
  element.setAttribute('aria-label', tytul);
  element.append(naglowek, opis, pola, akcje, tresc.element);

  return {
    element,
    tresc,
    dodajPole(etykieta, kontrolka, objasnienie) {
      const opcje: { klasa: string; objasnienie?: string } = { klasa: 'da-panel__pole' };
      if (objasnienie !== undefined) opcje.objasnienie = objasnienie;
      pola.append(wiersz(etykieta, kontrolka, opcje));
      return kontrolka;
    },
    dodajCzynnosc(etykieta, czynnosc) {
      const przycisk = przyciskAkcji(etykieta, 'dn-btn dn-btn--zarys');
      przycisk.addEventListener('click', czynnosc);
      akcje.append(przycisk);
      return przycisk;
    },
  };
}

/**
 * Wykonanie czynności panelu: ładowanie, a po odpowiedzi albo odmowa nazwana
 * powodem rdzenia, albo potwierdzenie wraz z zapisem odpowiedzi.
 *
 * Jedno miejsce, bo wszystkie czterdzieści trzy czynności kończą się tak samo
 * i różnią się wyłącznie zdaniem. Powtórzenie tego bloku czterdzieści trzy razy
 * dałoby czterdzieści trzy okazje do pomylenia odmowy z powodzeniem.
 */
export function wykonajCzynnoscPanelu<T>(
  panel: PanelDobudowy,
  opisPracy: string,
  obietnica: Promise<Wynik<T>>,
  zdanieOdmowy: string,
  zdaniePowodzenia: string,
): void {
  panel.tresc.ladowanie(opisPracy);
  void obietnica.then((wynik) => {
    if (!wynik.udany || wynik.wynik === undefined) {
      panel.tresc.blad(zdanieOdmowy, wynik.blad);
      return;
    }
    panel.tresc.potwierdzenie(zdaniePowodzenia, true);
    panel.tresc.tresc().replaceChildren(zapisOdpowiedzi(wynik.wynik));
  });
}

/**
 * Zapis odpowiedzi rdzenia w postaci czytelnej dla człowieka.
 *
 * Wartości wrażliwych tu nie ma i być nie może: skarbiec oddaje wyłącznie
 * referencje, a ładunki kroków wracają zredagowane przez rdzeń. Panel niczego
 * nie maskuje po swojej stronie, bo maskowanie po stronie okna dawałoby
 * złudzenie ochrony — wartość, która dotarła do przeglądarki, jest już
 * wyniesiona z serwera.
 */
function zapisOdpowiedzi(wynik: unknown): HTMLElement {
  const element = document.createElement('pre');
  element.className = 'da-panel__odpowiedz';
  element.textContent = JSON.stringify(wynik, null, 2);
  return element;
}

/**
 * Odczyt zapisu strukturalnego z pola tekstowego.
 *
 * Pole puste daje `undefined` — pole nieobecne w żądaniu znaczy co innego niż
 * pole o wartości pustej. Zapis nieczytelny daje `null`, a wołający odmawia
 * wysłania: żądanie z uszkodzonym ładunkiem odbiłoby się od rdzenia komunikatem
 * o kopercie, a Operator ma zobaczyć, że to on pomylił nawias.
 */
export function odczytajZapis(tekst: string): unknown | null | undefined {
  const oczyszczony = tekst.trim();
  if (oczyszczony === '') return undefined;
  try {
    return JSON.parse(oczyszczony) as unknown;
  } catch {
    return null;
  }
}

/** Liczba z pola liczbowego; pole puste albo nieczytelne daje `undefined`. */
export function odczytajLiczbe(tekst: string): number | undefined {
  const oczyszczony = tekst.trim();
  if (oczyszczony === '') return undefined;
  const liczba = Number(oczyszczony);
  return Number.isFinite(liczba) ? liczba : undefined;
}

/** Chwila z pola tekstowego w postaci daty; puste daje `undefined`. */
export function odczytajChwile(tekst: string): number | undefined {
  const oczyszczony = tekst.trim();
  if (oczyszczony === '') return undefined;
  const chwila = Date.parse(oczyszczony);
  return Number.isFinite(chwila) ? chwila : undefined;
}

/** Wykaz z pola tekstowego rozdzielonego przecinkami; puste daje wykaz pusty. */
export function odczytajWykaz(tekst: string): string[] {
  return tekst
    .split(',')
    .map((pozycja) => pozycja.trim())
    .filter((pozycja) => pozycja !== '');
}
