import { Command } from '../../../shared/contract';
import { elementIkony } from '../ikony/ikony';
import { ETYKIETA_ROLI_OKNA, ETYKIETA_STANU_SESJI } from './etykiety-pulpitu';
import type { KolumnaSrodowiska, SesjaMatrycy } from './model-danych';
import { utworzSekcje } from './naglowek-sekcji';
import type { WejscieDoSesji } from './zdarzenia-pulpitu';

/**
 * Matryca sesji: środowiska rdzenia obok siebie, a pod każdym jego sesje
 * z odczytu `session.list`. Kolumny powstają z `environment.list`, nie
 * z wykazu kodów klienta. Kafel sesji jest kontrolką: naciśnięcie nadaje
 * zamiar `session.open`.
 */
export interface SekcjaMatrycy {
  element: HTMLElement;
  /** Miejsce montażu pasa relacji, pod kolumnami środowisk macierzy sesji. */
  podKolumnami: HTMLElement;
  odswiez(matryca: KolumnaSrodowiska[], pozaSrodowiskami: SesjaMatrycy[]): void;
}

/** Buduje matrycę sesji pulpitu operacyjnego Mission Control wraz z miejscem montażu pasa relacji sesji. */
export function utworzSekcjeMatrycy(
  matryca: KolumnaSrodowiska[],
  pozaSrodowiskami: SesjaMatrycy[],
  nadaj: (wejscie: WejscieDoSesji) => void,
): SekcjaMatrycy {
  const { element, cialo } = utworzSekcje(
    'mc-sekcja--matryca',
    {
      tytul: 'Matryca sesji',
      // Liczby środowisk nie ma w napisie: wykaz oddaje rdzeń, nie wykaz kodów klienta.
      dopisek: 'Środowiska obok siebie, sesje pod każdym, powiązania pod spodem.',
      ikona: 'menu',
    },
    'mc-tytul-matryca',
  );

  const kolumny = document.createElement('div');
  kolumny.className = 'mc-matryca';

  const nieprzypisane = document.createElement('div');
  nieprzypisane.className = 'mc-matryca__poza';

  const podKolumnami = document.createElement('div');
  podKolumnami.className = 'mc-sekcja--matryca__pod';

  cialo.append(kolumny, nieprzypisane, podKolumnami);

  const odswiez = (dane: KolumnaSrodowiska[], poza: SesjaMatrycy[]): void => {
    // Brak kolumn mówi o sobie: kolumny powstają z `environment.list`, nie z wykazu kodów klienta.
    kolumny.replaceChildren(
      ...(dane.length === 0
        ? [wierszBezKolumn()]
        : dane.map((kolumna) => wyrysujKolumne(kolumna, nadaj))),
    );
    nieprzypisane.replaceChildren(...wyrysujPozaSrodowiskami(poza));
  };
  odswiez(matryca, pozaSrodowiskami);

  return { element, podKolumnami, odswiez };
}

/** Matryca przed odczytem środowisk z rdzenia — nazwany brak stanu ładowania, a nie pusty prostokąt karty. */
function wierszBezKolumn(): HTMLElement {
  const pusto = document.createElement('p');
  pusto.className = 'mc-kolumna__pusto mc-matryca__bez-kolumn';
  pusto.textContent =
    'Wykaz środowisk nie przyszedł jeszcze z rdzenia — kolumny powstają z odczytu environment.list.';
  return pusto;
}

/** Sesje bez przypisanej kolumny środowiska — osobny wiersz informacyjny pod matrycą, nie kontrolka wejścia. */
function wyrysujPozaSrodowiskami(poza: SesjaMatrycy[]): HTMLElement[] {
  if (poza.length === 0) return [];

  const etykieta = document.createElement('p');
  etykieta.className = 'mc-matryca__poza-etykieta';
  // Sesja trafia tu, gdy odczyt nie wskazał środowiska albo wykaz środowisk jeszcze nie przyszedł.
  etykieta.textContent =
    'Sesje bez kolumny — odczyt nie wskazał ich środowiska albo wykaz środowisk jeszcze nie przyszedł:';

  const wykaz = document.createElement('ul');
  wykaz.className = 'mc-matryca__poza-wykaz';
  for (const sesja of poza) {
    const wiersz = document.createElement('li');
    wiersz.className = 'mc-matryca__poza-sesja';
    wiersz.textContent = `${sesja.tytul} · ${ETYKIETA_STANU_SESJI[sesja.status]}`;
    wykaz.append(wiersz);
  }
  return [etykieta, wykaz];
}

/** Jedna kolumna matrycy sesji pulpitu operacyjnego Mission Control: nagłówek środowiska z mottem i wykazem. */
function wyrysujKolumne(
  kolumna: KolumnaSrodowiska,
  nadaj: (wejscie: WejscieDoSesji) => void,
): HTMLElement {
  const element = document.createElement('div');
  element.className = 'mc-kolumna';
  element.dataset.srodowisko = kolumna.id;

  const naglowek = document.createElement('div');
  naglowek.className = 'mc-kolumna__naglowek';

  const nazwa = document.createElement('h3');
  nazwa.className = 'mc-kolumna__nazwa';
  nazwa.textContent = kolumna.nazwa;

  const motto = document.createElement('p');
  motto.className = 'mc-kolumna__motto';
  motto.textContent = kolumna.motto;

  const licznik = document.createElement('span');
  licznik.className = 'dn-plakietka dn-plakietka--rola mc-kolumna__licznik';
  licznik.textContent = `${kolumna.sesje.length} sesji`;

  naglowek.append(nazwa, motto, licznik);

  // Kolumna bez sesji mówi to wprost — pusta przestrzeń wygląda jak usterka rysowania.
  const lista = document.createElement('div');
  lista.className = 'mc-kolumna__sesje';
  if (kolumna.sesje.length === 0) {
    const pusto = document.createElement('p');
    pusto.className = 'mc-kolumna__pusto';
    pusto.textContent = 'Żadna sesja nie biegnie w tym środowisku.';
    lista.append(pusto);
  } else {
    lista.append(...kolumna.sesje.map((sesja) => kafelSesji(sesja, kolumna, nadaj)));
  }

  element.append(naglowek, lista);
  return element;
}

/** Kafel jednej sesji w kolumnie środowiska matrycy pulpitu operacyjnego — kontrolka wejścia do tej sesji. */
function kafelSesji(
  sesja: SesjaMatrycy,
  kolumna: KolumnaSrodowiska,
  nadaj: (wejscie: WejscieDoSesji) => void,
): HTMLButtonElement {
  const kafel = document.createElement('button');
  kafel.type = 'button';
  kafel.className = 'dn-karta dn-karta--klikalna mc-sesja';
  kafel.dataset.sesja = sesja.id;
  kafel.title = `Wejdź do sesji: ${sesja.tytul}`;

  const gora = document.createElement('span');
  gora.className = 'mc-sesja__gora';

  // Wskaźnik pracy w tle — kropka pulsująca przy tytule sesji.
  const kropka = document.createElement('span');
  kropka.className = sesja.pracaWTle
    ? 'dn-kropka dn-kropka--tetno mc-sesja__kropka'
    : 'dn-kropka mc-sesja__kropka';
  kropka.setAttribute('aria-hidden', 'true');

  const tytul = document.createElement('span');
  tytul.className = 'mc-sesja__tytul';
  tytul.textContent = sesja.tytul;

  gora.append(kropka, tytul);

  const opis = document.createElement('span');
  opis.className = 'mc-sesja__opis';
  // Stan zawsze słowem, nigdy samą barwą kropki; rola `null` znaczy brak okna wiodącego.
  opis.textContent = [
    ETYKIETA_STANU_SESJI[sesja.status],
    sesja.rolaOkna === null ? 'rola okna nieznana' : ETYKIETA_ROLI_OKNA[sesja.rolaOkna],
    `${sesja.okna} ${odmianaOkna(sesja.okna)}`,
    sesja.pracaWTle ? 'praca w tle' : 'bez pracy w tle',
  ].join(' · ');

  const wejscie = document.createElement('span');
  wejscie.className = 'mc-sesja__wejscie';
  wejscie.append(
    document.createTextNode('Wejdź'),
    elementIkony('strzalka-prawo', { rozmiar: 16 }),
  );

  kafel.append(gora, opis, wejscie);
  kafel.addEventListener('click', () => {
    nadaj({
      komenda: Command.SessionOpen,
      sessionId: sesja.id,
      srodowisko: kolumna.id,
      tytul: sesja.tytul,
    });
  });

  return kafel;
}

/** Odmiana słowa „okno" dobrana do liczby okien komunikacji danej sesji: „1 okno", „3 okna", „5 okien". */
function odmianaOkna(liczba: number): string {
  if (liczba === 1) return 'okno';
  const dziesiatki = liczba % 100;
  const jednosci = liczba % 10;
  const mnoga = jednosci >= 2 && jednosci <= 4 && !(dziesiatki >= 12 && dziesiatki <= 14);
  return mnoga ? 'okna' : 'okien';
}
