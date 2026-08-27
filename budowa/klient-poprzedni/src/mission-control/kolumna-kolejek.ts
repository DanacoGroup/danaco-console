import { QueueStatus } from '../../../shared/contract';
import { ETYKIETA_BRAKU_ZRODLA, ETYKIETA_STANU_KOLEJKI } from './etykiety-pulpitu';
import type { KolejkaPulpitu } from './model-danych';
import { utworzSekcje } from './naglowek-sekcji';
import { utworzPasekTransportu } from './przyciski-transportu';
import { utworzStanPusty } from './stan-pusty';
import type { ZamiarKolejki } from './zdarzenia-pulpitu';

/**
 * Kolumna kolejek jest drugą z trzech kolumn pulpitu i pokazuje kolejki znane
 * ze zdarzeń zmiany, każdą z własnym paskiem transportu. Miary, których kontrakt
 * nie niesie, kolumna wypisuje wprost jako brak źródła.
 */
export interface KolumnaKolejek {
  element: HTMLElement;
  odswiez(kolejki: KolejkaPulpitu[]): void;
}

/**
 * Buduje kolumnę kolejek wraz ze sterowaniem: wykaz wierszy kolejek, stan pusty
 * na wypadek pustego wykazu oraz odczyt przyjmujący nowy wykaz kolejek. Zamiary
 * sterowania kolumna przekazuje podanej funkcji.
 */
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

/**
 * Buduje wiersz jednej kolejki: nazwę, plakietkę stanu, liczniki okien i obiegów
 * naprawczych oraz pasek przycisków transportu. Identyfikator i stan kolejki
 * wiersz niesie w polach danych elementu.
 */
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

/**
 * Buduje pojedynczy licznik kolejki: widoczny tekst miary wraz z objaśnieniem
 * podawanym w podpowiedzi elementu. Objaśnienie nazywa wielkość mierzoną oraz
 * jej źródło.
 */
function licznik(tekst: string, objasnienie: string): HTMLElement {
  const element = document.createElement('span');
  element.className = 'mc-kolejka__licznik';
  element.textContent = tekst;
  element.title = objasnienie;
  return element;
}

/**
 * Dobiera odmianę plakietki do stanu kolejki: bieg daje odmianę sukcesu,
 * wstrzymanie odmianę ostrzeżenia, zatrzymanie odmianę błędu, a stan pozostały
 * odmianę informacyjną.
 */
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

/**
 * Odmienia rzeczownik nazywający okno przez liczbę, zgodnie z regułą polskiej
 * liczby mnogiej wraz z wyjątkiem obejmującym liczebniki od jedenastu do
 * czternastu.
 */
function odmianaOkna(liczba: number): string {
  if (liczba === 1) {
    return 'okno';
  }
  const dziesiatki = liczba % 100;
  const jednosci = liczba % 10;
  const mnoga = jednosci >= 2 && jednosci <= 4 && !(dziesiatki >= 12 && dziesiatki <= 14);
  return mnoga ? 'okna' : 'okien';
}

/**
 * Odmienia rzeczownik nazywający obieg przez liczbę, zgodnie z regułą polskiej
 * liczby mnogiej wraz z wyjątkiem obejmującym liczebniki od jedenastu do
 * czternastu.
 */
function odmianaObiegu(liczba: number): string {
  if (liczba === 1) {
    return 'obieg';
  }
  const dziesiatki = liczba % 100;
  const jednosci = liczba % 10;
  const mnoga = jednosci >= 2 && jednosci <= 4 && !(dziesiatki >= 12 && dziesiatki <= 14);
  return mnoga ? 'obiegi' : 'obiegów';
}
