import { BRAKI } from './etykiety-designu';
import { przyciskBrakuDrogi } from './brak-drogi';
import type { OsWyrownania, SzablonUkladu } from './zapis-kompozycji';
import type { StanKompozycji } from './stan-kompozycji';

/**
 * Przybornik Design Board — narzędzia panelu akcji modułu nad kanwą.
 *
 * Jedna odpowiedzialność: kontrolki zmieniające kompozycję i widok kanwy.
 *
 * Grupy dzielą się tym, co narzędzie zmienia: widok kanwy (powiększenie,
 * przesunięcie, siatka) · układ warstw zaznaczonych (wyrównanie,
 * rozmieszczenie) · skład kompozycji (szablony, elementy pomocnicze,
 * zaznaczenie wszystkiego) · czynności bez drogi w kontrakcie (wersjonowanie,
 * eksport, kursor współpracy).
 *
 * Grupa bez drogi w kontrakcie stoi w przyborniku, a nie poza nim, żeby komplet
 * narzędzi był widoczny razem z tym, które z nich nie mają wykonania — nazwane,
 * klikalne i mówiące dlaczego.
 */
export interface NarzedziaPlanszy {
  element: HTMLElement;
  odswiez(): void;
}

/** Wyrównania z panelu akcji, w kolejności czytania układu. */
const WYROWNANIA: readonly (readonly [OsWyrownania, string])[] = [
  ['lewo', 'Do lewej'],
  ['srodek', 'Do środka w poziomie'],
  ['prawo', 'Do prawej'],
  ['gora', 'Do góry'],
  ['srodek-pion', 'Do środka w pionie'],
  ['dol', 'Do dołu'],
];

/** Cztery szablony układu wymienione w panelu akcji modułu Design. */
const SZABLONY: readonly (readonly [SzablonUkladu, string])[] = [
  ['tablica-nastroju', 'Tablica nastroju'],
  ['siatka-porownawcza', 'Siatka porównawcza'],
  ['arkusz-brandingowy', 'Arkusz brandingowy'],
  ['formaty-spolecznosciowe', 'Formaty społecznościowe'],
];

/** Biblioteka elementów pomocniczych — warstwy bez zasobu. */
const ELEMENTY_POMOCNICZE: readonly string[] = ['Ramka', 'Podpis', 'Strzałka', 'Pole notatki'];

/**
 * Presety ramek obszaru roboczego — szerokości z punktów łamania produktu.
 *
 * Liczby nie są wzięte z cudzych urządzeń: to punkty łamania kierunku
 * projektowego platformy, te same, na których stoi cały jej układ. Dzięki temu
 * ramka makiety ma dokładnie tę szerokość, przy której produkt zmienia postać,
 * a nie szerokość telefonu, który akurat był w sprzedaży.
 *
 * Wysokości preset nie ustawia — punkty łamania jej nie rozstrzygają.
 */
const RAMKI_URZADZEN: readonly (readonly [string, number])[] = [
  ['Telefon poziomo', 640],
  ['Tablet', 960],
  ['Biurko', 1280],
  ['Szerokie biurko', 1600],
];

export function utworzNarzedziaPlanszy(stan: StanKompozycji): NarzedziaPlanszy {
  const licznik = document.createElement('span');
  licznik.className = 'dn-plakietka';

  const siatka = przyciskNarzedzia('Siatka i linie pomocnicze', () =>
    stan.ustawWidok({ siatka: !stan.widok().siatka }),
  );

  const element = document.createElement('div');
  element.className = 'dn-przybornik md-przybornik';
  element.setAttribute('aria-label', 'Narzędzia kompozycji');
  element.append(
    grupa('Widok kanwy', [
      przyciskNarzedzia('Powiększ', () => zmienPowiekszenie(stan, 1.25)),
      przyciskNarzedzia('Pomniejsz', () => zmienPowiekszenie(stan, 0.8)),
      przyciskNarzedzia('Widok naturalny', () =>
        stan.ustawWidok({ powiekszenie: 1, przesuniecieX: 0, przesuniecieY: 0 }),
      ),
      siatka,
    ]),
    grupa('Wyrównanie i rozmieszczenie', [
      ...WYROWNANIA.map(([os, nazwa]) => przyciskNarzedzia(nazwa, () => stan.wyrownaj(os))),
      przyciskNarzedzia('Rozmieść w poziomie', () => stan.rozmiesc(false)),
      przyciskNarzedzia('Rozmieść w pionie', () => stan.rozmiesc(true)),
    ]),
    grupa('Skład kompozycji', [
      ...SZABLONY.map(([kod, nazwa]) => przyciskNarzedzia(nazwa, () => stan.ulozSzablon(kod))),
      ...ELEMENTY_POMOCNICZE.map((nazwa) =>
        przyciskNarzedzia(`+ ${nazwa}`, () => stan.dolozElement(nazwa)),
      ),
      przyciskNarzedzia('Zaznacz wszystkie', () => stan.zaznaczWszystkie()),
      licznik,
    ]),
    grupa('Ramki obszaru roboczego', [
      ...RAMKI_URZADZEN.map(([nazwa, szerokosc]) =>
        przyciskNarzedzia(`+ ${nazwa} (${String(szerokosc)} px)`, () =>
          stan.dolozRamke(`${nazwa} — ${String(szerokosc)} px`, szerokosc),
        ),
      ),
    ]),
    grupa('Bez drogi w kontrakcie', [
      przyciskBrakuDrogi(BRAKI.wersjeKompozycji),
      przyciskBrakuDrogi(BRAKI.eksportKompozycji),
      przyciskBrakuDrogi(BRAKI.kursorWspolpracy),
    ]),
  );

  return {
    element,

    odswiez() {
      const widok = stan.widok();
      siatka.dataset['wlaczona'] = String(widok.siatka);
      siatka.setAttribute('aria-pressed', String(widok.siatka));
      licznik.textContent =
        `${stan.zaznaczone().length} z ${stan.warstwy().length} warstw zaznaczonych · ` +
        `powiększenie ${Math.round(widok.powiekszenie * 100)}%`;
    },
  };
}

/** Powiększenie w granicach czytelności kanwy. */
function zmienPowiekszenie(stan: StanKompozycji, mnoznik: number): void {
  const nowe = stan.widok().powiekszenie * mnoznik;
  stan.ustawWidok({ powiekszenie: Math.min(4, Math.max(0.25, nowe)) });
}

function przyciskNarzedzia(nazwa: string, czynnosc: () => void): HTMLButtonElement {
  const element = document.createElement('button');
  element.type = 'button';
  element.className = 'dn-btn dn-btn--duch dn-btn--sm';
  element.textContent = nazwa;
  element.addEventListener('click', czynnosc);
  return element;
}

function grupa(nazwa: string, kontrolki: readonly HTMLElement[]): HTMLElement {
  const tytul = document.createElement('span');
  tytul.className = 'md-przybornik__tytul';
  tytul.textContent = nazwa;

  const element = document.createElement('div');
  element.className = 'md-przybornik__grupa';
  element.append(tytul, ...kontrolki);
  return element;
}
