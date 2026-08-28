import { utworzPanelRodzin } from './panel-rodzin';
import { sekcjeMaterialu } from './sekcje-rodzin';
import { utworzDymekObjasnienia } from '../../komponenty/dymek';
import { utworzNaglowekOkna } from '../../komponenty/naglowek-okna';
import { utworzWierszOdpowiedzi } from '../../modele/kontrolki-formularza';
import { utworzCzynnosciMaterialu, type CzynnosciMaterialu } from './czynnosci-materialu';
import {
  KLASY_DYMKA,
  KODY_OKIEN,
  OBJASNIENIA,
  STANY_PUSTE,
  WYKAZY,
} from './etykiety-browser';
import { nazwaRodzaju, type MonitorZmian, type Przechwycenie } from './material-sesji';
import { KLASA_PRZYCISKU, przyciskCzynnosci } from './przyciski-browser';
import { utworzStanOkna, type StanOkna } from './stan-okna';
import type { StanPrzegladania } from './stan-przegladania';
import { utworzSterWyboru } from './ster-wyboru';

/**
 * Capture & Monitor Panel — okno przechwytywania i monitorowania modułu
 * Browser, warstwa trzecia, otwierane z menu operacji. Jedna odpowiedzialność:
 * złożenie okna oraz wykazy materiału i monitorów.
 */
export interface OknoCaptureMonitor {
  element: HTMLElement;
  odswiez(): void;
}

/** Nastawa widoku wykazu materiału sesji — filtr rodzaju pozycji pokazywanych aktualnie w liście panelu. */
const WIDOKI = [
  { wartosc: 'wszystko', etykieta: 'Cały materiał', opis: 'Zrzuty, archiwa i treści stron.' },
  { wartosc: 'zrzut', etykieta: 'Zrzuty', opis: 'Pozycje z odnośnikiem do zrzutu ekranu.' },
  { wartosc: 'archiwum', etykieta: 'Archiwa', opis: 'Pozycje ze źródłem strony.' },
  { wartosc: 'tresc', etykieta: 'Treści', opis: 'Pozycje z samą treścią renderowaną.' },
];

export function utworzOknoCaptureMonitor(stan: StanPrzegladania): OknoCaptureMonitor {
  const okno = utworzStanOkna(STANY_PUSTE.materialy);
  const odpowiedz = utworzWierszOdpowiedzi();
  const powiedz = (tresc: string, ok: boolean): void => odpowiedz.pokaz(tresc, ok);
  const czynnosci = utworzCzynnosciMaterialu(stan, powiedz);

  const podglad = document.createElement('pre');
  podglad.className = 'mb-material__podglad';
  podglad.hidden = true;

  const lista = document.createElement('ul');
  lista.className = 'mb-material__lista';

  const monitory = document.createElement('ul');
  monitory.className = 'mb-material__monitory';

  const widok = utworzSterWyboru({
    nastawa: 'Widok materiału',
    pozycje: WIDOKI,
    klasa: 'mb-material__ster',
    podpis: false,
    naZmiane: () => odswiez(),
  });

  const uwaga = document.createElement('p');
  uwaga.className = 'dn-pole-opis mb-uwaga';
  uwaga.textContent = WYKAZY.material;

  const pasek = document.createElement('div');
  pasek.className = 'mb-panel__pasek';
  pasek.append(
    przyciskCzynnosci('Przechwyć bieżącą stronę', KLASA_PRZYCISKU.glowny, () =>
      void czynnosci.przechwyc(),
    ),
    utworzDymekObjasnienia(OBJASNIENIA.przechwycenie, KLASY_DYMKA),
    przyciskCzynnosci('Załóż monitor strony', KLASA_PRZYCISKU.zarys, () =>
      czynnosci.zalozMonitor(),
    ),
    utworzDymekObjasnienia(OBJASNIENIA.monitorZmian, KLASY_DYMKA),
    widok.element,
    przyciskCzynnosci('→ Wyślij do Library', KLASA_PRZYCISKU.zarys, () =>
      void czynnosci.przekaz('library'),
    ),
    przyciskCzynnosci('→ Wyślij do Research', KLASA_PRZYCISKU.zarys, () =>
      void czynnosci.przekaz('research'),
    ),
    utworzDymekObjasnienia(OBJASNIENIA.przekazanie, KLASY_DYMKA),
  );

  // Bez tego panelu operator nie miałby z okna drogi do komend, które rdzeń już obsługuje.
  const rodziny = utworzPanelRodzin(sekcjeMaterialu(stan));

  okno.tresc.append(lista, monitory, podglad, rodziny.element, odpowiedz.element);

  const element = document.createElement('section');
  element.className = 'mb-okno mb-okno--pomocnicze';
  element.dataset['okno'] = KODY_OKIEN.materialy;
  element.setAttribute('aria-label', 'Capture & Monitor Panel — materiał sesji przeglądania');
  element.append(
    utworzNaglowekOkna({ tytul: 'Capture & Monitor Panel', klasa: 'mb-okno__naglowek' }),
    pasek,
    uwaga,
    okno.element,
  );

  /** Pokazuje treść wybranej pozycji materiału w podglądzie panelu, bez opuszczania widoku wykazu. */
  function pokaz(pozycja: Przechwycenie): void {
    const migawka = pozycja.migawka;
    podglad.hidden = false;
    podglad.textContent =
      `${migawka.url}\n${(migawka.title ?? '').trim()}\n\n` +
      `${migawka.text ?? '(rdzeń nie oddał treści renderowanej)'}`;
    powiedz(
      `Podgląd pozycji: ${nazwaRodzaju(pozycja.rodzaj)} strony ${migawka.url}.`,
      true,
    );
  }

  function odswiez(): void {
    const wybrany = widok.wartosc();
    const pozycje = stan.material
      .przechwycenia()
      .filter((pozycja) => wybrany === 'wszystko' || pozycja.rodzaj === wybrany);
    lista.replaceChildren(...pozycje.map((pozycja) => wierszPozycji(pozycja, czynnosci, pokaz)));
    monitory.replaceChildren(
      ...stan.material.monitory().map((monitor) => wierszMonitora(monitor, stan, czynnosci)),
    );
    nanieStan(okno, stan, pozycje.length + stan.material.monitory().length);
  }

  okno.ustawPonowienie(odswiez);
  odswiez();

  return { element, odswiez };
}

/** Jedna pozycja materiału sesji przeglądania wraz z czynnościami dostępnymi wprost z poziomu wykazu panelu. */
function wierszPozycji(
  pozycja: Przechwycenie,
  czynnosci: CzynnosciMaterialu,
  pokaz: (pozycja: Przechwycenie) => void,
): HTMLElement {
  const migawka = pozycja.migawka;
  const tytul = (migawka.title ?? '').trim();

  const opis = document.createElement('span');
  opis.className = 'mb-material__opis';
  opis.textContent =
    `${nazwaRodzaju(pozycja.rodzaj)} · ${tytul === '' ? migawka.url : tytul} · ` +
    `${czasPozycji(migawka.capturedAt)}`;

  const element = document.createElement('li');
  element.className = 'mb-material__wiersz';
  element.dataset['rodzaj'] = pozycja.rodzaj;
  element.append(
    opis,
    przyciskCzynnosci('Podgląd', KLASA_PRZYCISKU.duch, () => pokaz(pozycja)),
    przyciskCzynnosci('Otwórz ponownie', KLASA_PRZYCISKU.duch, () => void czynnosci.otworz(migawka)),
  );
  return element;
}

/** Jeden monitor zmian strony wraz z wynikiem ostatniego sprawdzenia wykonanego przez rdzeń w tej sesji. */
function wierszMonitora(
  monitor: MonitorZmian,
  stan: StanPrzegladania,
  czynnosci: CzynnosciMaterialu,
): HTMLElement {
  const opis = document.createElement('span');
  opis.className = 'mb-material__opis';
  opis.textContent = `monitor · ${monitor.adres} · ${opisWyniku(monitor)}`;

  const element = document.createElement('li');
  element.className = 'mb-material__wiersz';
  element.dataset['wynik'] = monitor.wynik;
  element.append(
    opis,
    przyciskCzynnosci('Sprawdź teraz', KLASA_PRZYCISKU.duch, () =>
      void czynnosci.sprawdzMonitor(monitor),
    ),
    przyciskCzynnosci('Zdejmij monitor', KLASA_PRZYCISKU.duch, () =>
      stan.material.usunMonitor(monitor.adres),
    ),
  );
  return element;
}

/** Stan monitora zmian wyrażony słowem — sam znacznik barwy nie mówi operatorowi, co się faktycznie stało. */
function opisWyniku(monitor: MonitorZmian): string {
  if (monitor.wynik === 'nietkniety') return 'jeszcze nie sprawdzany';
  if (monitor.wynik === 'bez-zmian') return `bez zmian, sprawdzony ${czasPozycji(monitor.sprawdzonyO)}`;
  return `ZMIANA (${monitor.roznicaZnakow > 0 ? '+' : ''}${monitor.roznicaZnakow} znaków)`;
}

/** Godzina zdarzenia w postaci lokalnej dla operatora; czas rdzenia liczony jest w milisekundach epoki. */
function czasPozycji(znacznik: number): string {
  if (znacznik <= 0) return 'bez znacznika czasu';
  return new Date(znacznik).toLocaleTimeString('pl-PL');
}

/**
 * Trzy stany obowiązkowe panelu: czekanie, odmowa, pustka.
 *
 * Panel czyta stan okna przeglądarki, bo bez ustalonego okna nie ma czego
 * przechwycić — a wtedy pustka ma powód, nie tylko brak treści.
 */
function nanieStan(okno: StanOkna, stan: StanPrzegladania, pozycji: number): void {
  if (stan.faza() === 'odczyt') {
    okno.ladowanie('Rdzeń ustala okno przeglądarki tej sesji…');
    return;
  }
  if (stan.faza() === 'blad') {
    okno.blad(stan.powod());
    return;
  }
  if (pozycji === 0) {
    okno.pusteZOpisu();
    return;
  }
  okno.gotowe();
}
