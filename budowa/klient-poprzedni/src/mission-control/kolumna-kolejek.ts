import { QueueStatus } from '../../../shared/contract';
import { ETYKIETA_BRAKU_ZRODLA, ETYKIETA_STANU_KOLEJKI } from './etykiety-pulpitu';
import type { KolejkaPulpitu } from './model-danych';
import { utworzSekcje } from './naglowek-sekcji';
import { utworzPasekTransportu } from './przyciski-transportu';
import { utworzStanPusty } from './stan-pusty';
import type { ZamiarKolejki } from './zdarzenia-pulpitu';

/**
 * Druga z trzech kolumn pulpitu — kolejki ze sterowaniem.
 *
 * Jedna odpowiedzialność: kolejki znane ze zdarzeń `queue.changed`, każda
 * z własnym paskiem transportu (`queue.action`). Kontrakt nie niesie roli
 * kolejki, liczby zadań ani odczytu `queue.list`, więc wykaz buduje się
 * wyłącznie ze zdarzeń, a brakujące miary są wypisane jako brak źródła.
 *
 * Bieg naprawczy nie ma limitu, dlatego licznik obiegów stoi przy każdej
 * kolejce — przejrzystość zastępuje tu bramę. Sterowanie jest ręczne, żaden
 * przycisk nie jest wyszarzony; stan kolejki poznaje się po plakietce ze słowem.
 */
export interface KolumnaKolejek {
  element: HTMLElement;
  odswiez(kolejki: KolejkaPulpitu[]): void;
}

/** Buduje kolumnę kolejek ze sterowaniem. */
export function utworzKolumneKolejek(
  kolejki: KolejkaPulpitu[],
  nadaj: (zamiar: ZamiarKolejki) => void,
): KolumnaKolejek {
  const { element, cialo } = utworzSekcje(
    'mc-sekcja--kolejki',
    {
      tytul: 'Kolejki',
      dopisek: 'Sterowanie ręczne. Bez limitu obiegów, zatrzymanie zawsze czynne.',
      ikona: 'ustawienia',
    },
    'mc-tytul-kolejki',
  );

  const wykaz = document.createElement('ul');
  wykaz.className = 'mc-kolejki';
  cialo.append(wykaz);

  const odswiez = (dane: KolejkaPulpitu[]): void => {
    if (dane.length === 0) {
      wykaz.replaceChildren(
        utworzStanPusty(
          'Żadna kolejka nie jest znana',
          'Kontrakt nie ma odczytu queue.list — wykaz buduje się ze zdarzeń queue.changed, a żadne dotąd nie nadeszło.',
        ),
      );
      return;
    }
    wykaz.replaceChildren(...dane.map((kolejka) => wiersz(kolejka, nadaj)));
  };
  odswiez(kolejki);

  return { element, odswiez };
}

/** Jedna kolejka: nazwa, stan, liczniki, pasek przycisków transportu. */
function wiersz(kolejka: KolejkaPulpitu, nadaj: (zamiar: ZamiarKolejki) => void): HTMLLIElement {
  const element = document.createElement('li');
  element.className = 'mc-kolejka';
  element.dataset.kolejka = kolejka.id;
  element.dataset.stan = kolejka.status;

  const gora = document.createElement('div');
  gora.className = 'mc-kolejka__gora';

  const nazwa = document.createElement('span');
  nazwa.className = 'mc-kolejka__rola';
  nazwa.textContent = kolejka.nazwa;

  const stan = document.createElement('span');
  stan.className = `dn-plakietka ${klasaStanu(kolejka.status)} mc-kolejka__stan`;
  stan.textContent = ETYKIETA_STANU_KOLEJKI[kolejka.status];

  gora.append(nazwa, stan);

  const liczniki = document.createElement('div');
  liczniki.className = 'mc-kolejka__liczniki';
  liczniki.append(
    licznik(
      kolejka.okna === null ? 'okna: nie podano' : `${kolejka.okna} ${odmianaOkna(kolejka.okna)}`,
      'Okna komunikacji obsługiwane przez kolejkę',
    ),
    licznik(
      kolejka.obiegi === null
        ? 'obiegi: nie podano'
        : `${kolejka.obiegi} ${odmianaObiegu(kolejka.obiegi)}`,
      'Licznik obiegów naprawczych — bez limitu',
    ),
    licznik(
      `zadania: ${ETYKIETA_BRAKU_ZRODLA}`,
      'Kontrakt nie niesie liczby zadań czekających w kolejce. Brakujący odczyt zgłoszony.',
    ),
  );

  element.append(gora, liczniki, utworzPasekTransportu(kolejka, nadaj));
  return element;
}

/** Pojedynczy licznik z objaśnieniem w podpowiedzi. */
function licznik(tekst: string, objasnienie: string): HTMLElement {
  const element = document.createElement('span');
  element.className = 'mc-kolejka__licznik';
  element.textContent = tekst;
  element.title = objasnienie;
  return element;
}

/** Klasa plakietki stanu kolejki. */
function klasaStanu(status: QueueStatus): string {
  switch (status) {
    case QueueStatus.Running:
      return 'dn-plakietka--sukces';
    case QueueStatus.Paused:
      return 'dn-plakietka--ostrzezenie';
    case QueueStatus.Stopped:
      return 'dn-plakietka--blad';
    default:
      return 'dn-plakietka--informacja';
  }
}

/** „1 okno", „3 okna", „7 okien". */
function odmianaOkna(liczba: number): string {
  if (liczba === 1) {
    return 'okno';
  }
  const dziesiatki = liczba % 100;
  const jednosci = liczba % 10;
  const mnoga = jednosci >= 2 && jednosci <= 4 && !(dziesiatki >= 12 && dziesiatki <= 14);
  return mnoga ? 'okna' : 'okien';
}

/** „0 obiegów", „1 obieg", „3 obiegi", „7 obiegów". */
function odmianaObiegu(liczba: number): string {
  if (liczba === 1) {
    return 'obieg';
  }
  const dziesiatki = liczba % 100;
  const jednosci = liczba % 10;
  const mnoga = jednosci >= 2 && jednosci <= 4 && !(dziesiatki >= 12 && dziesiatki <= 14);
  return mnoga ? 'obiegi' : 'obiegów';
}
