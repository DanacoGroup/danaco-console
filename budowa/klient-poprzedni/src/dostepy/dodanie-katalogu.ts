import { elementIkony } from '../ikony/ikony';
import { Cel, czyPowlokaNatywna, naWskazanieZZasobnika, wskazKatalog } from './dialog-katalogu';
import { przyciskKarty } from './elementy-karty';
import { utworzKomunikatCzynnosci } from './komunikat-czynnosci';
import type { StanDostepow } from './stan-dostepow';

/**
 * Dodanie punktu dostępu do katalogu na „Mój komputer".
 *
 * Droga pierwsza i właściwa: natywne okno wyboru powłoki. Ścieżka wraca
 * istniejąca i rozwinięta, więc punkt zakładany na niej nie wskazuje miejsca,
 * którego nie ma.
 *
 * Droga druga: pole ścieżki wpisanej z ręki. Nie jest atrapą ani zapasem na
 * gorsze czasy — interfejs bywa otwarty w przeglądarce, bez powłoki, a wtedy
 * natywnego okna po prostu nie ma. Pole zostaje czynne zawsze, także
 * w powłoce: Operator, który zna ścieżkę, nie musi jej odklikiwać.
 *
 * Przycisk natywnego wyboru nie jest wygaszany poza powłoką. Naciśnięcie daje
 * odpowiedź: zdanie o tym, że okno systemowe należy do powłoki, i przeniesienie
 * uwagi do pola ścieżki.
 */
export interface DodanieKatalogu {
  /** Element osadzany pod wykazem punktów. */
  element: HTMLElement;
  /** Odłącza nasłuch wyboru z zasobnika powłoki. */
  rozlacz(): void;
}

export function utworzDodanieKatalogu(stan: StanDostepow): DodanieKatalogu {
  const komunikat = utworzKomunikatCzynnosci();

  const pole = document.createElement('input');
  pole.type = 'text';
  pole.className = 'dd-dodanie__pole';
  pole.placeholder = 'ścieżka katalogu, np. C:\\Projekty\\Lex';
  pole.setAttribute('aria-label', 'Ścieżka katalogu do udostępnienia');
  pole.addEventListener('keydown', (zdarzenie) => {
    if (zdarzenie.key === 'Enter') void zaloz(pole.value.trim());
  });

  const zPowloki = przyciskKarty('Dodaj katalog z Mój komputer', 'dn-btn dn-btn--zarys dn-btn--sm');
  zPowloki.prepend(elementIkony('folder', { rozmiar: 16 }));
  zPowloki.addEventListener('click', () => void zWyboruPowloki());

  const zPola = przyciskKarty('Dodaj ze ścieżki', 'dn-btn dn-btn--atrament dn-btn--sm');
  zPola.addEventListener('click', () => void zaloz(pole.value.trim()));

  const wiersz = document.createElement('div');
  wiersz.className = 'dd-dodanie__wiersz';
  wiersz.append(zPowloki, pole, zPola);

  const objasnienie = document.createElement('p');
  objasnienie.className = 'dd-dodanie__objasnienie';
  objasnienie.textContent =
    'Katalog dodany tutaj staje się punktem dostępu platformy. Model sięgnie do niego dopiero po nadaniu dostępu oknu rozmowy — dodanie punktu samo w sobie niczego nie otwiera.';

  const element = document.createElement('section');
  element.className = 'dd-dodanie';
  element.append(naglowek(), wiersz, objasnienie, komunikat.element);

  /** Wybór natywny; poza powłoką kieruje Operatora do pola ścieżki. */
  async function zWyboruPowloki(): Promise<void> {
    if (!czyPowlokaNatywna()) {
      komunikat.pokaz(
        'Natywne okno wyboru należy do powłoki Danaco Console. W przeglądarce wpisz ścieżkę katalogu w polu obok.',
        false,
      );
      pole.focus();
      return;
    }
    const sciezka = await wskazKatalog(Cel.PunktDostepu);
    if (sciezka === '') {
      komunikat.pokaz('Nie wskazano katalogu — nic się nie zmieniło.', true);
      return;
    }
    await zaloz(sciezka);
  }

  /**
   * Założenie punktu na wskazanej ścieżce.
   *
   * Katalog dodany tutaj leży na maszynie bieżącej — na tej, na której stoi
   * powłoka. Schemat wymaga wskazania urządzenia
   * (`migracja_013_punkty_dostepu.sql`), a jego identyfikator poda komenda
   * `device.list`, gdy trafi do kontraktu. Do tego czasu drugi argument
   * `zalozKatalogLokalny` zostaje pusty i rdzeń odmawia z powodem — bez atrapy
   * identyfikatora.
   */
  async function zaloz(sciezka: string): Promise<void> {
    if (sciezka === '') {
      komunikat.pokaz('Podaj ścieżkę katalogu albo wskaż go oknem powłoki.', false);
      pole.focus();
      return;
    }
    komunikat.pokaz('Zakładam punkt dostępu…', true);
    const wynik = await stan.zalozKatalogLokalny(sciezka);
    komunikat.zWyniku(
      wynik,
      `Katalog ${sciezka} jest punktem dostępu. Nadaj go oknu, żeby model mógł tam sięgnąć.`,
    );
    if (wynik.udany) pole.value = '';
  }

  // Zasobnik powłoki ma własną pozycję wyboru katalogu i rozgłasza wynik
  // zdarzeniem. Bez tego nasłuchu wybór z zasobnika nie miałby konsumenta.
  const odsubskrybuj = naWskazanieZZasobnika((sciezka) => {
    pole.value = sciezka;
    void zaloz(sciezka);
  });

  return { element, rozlacz: odsubskrybuj };
}

function naglowek(): HTMLElement {
  const element = document.createElement('h4');
  element.className = 'dd-dodanie__tytul';
  element.textContent = 'Dodaj katalog lokalny';
  return element;
}

