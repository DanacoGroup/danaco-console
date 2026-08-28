// Arkusz potwierdzenia usunięcia jest wzięty wprost z toru sesji, bo rodzina klas biblioteczna ma jeden arkusz przy potwierdzeniu sesji, a drugi arkusz tej samej czynności rozjechałby się przy poprawce.
import '../powloka/usuniecie-sesji.css';

import type { HistoryDeleteResponse } from '../../../shared/contract';
import { utworzDymekObjasnienia } from '../komponenty/dymek';
import { oznaczFaze } from '../komponenty/faza-okna';
import { opisOdmowyBledu } from '../komponenty/odmowa';
import type { Wynik } from '../protokol/kanal';

/**
 * Potwierdzenie wyczyszczenia całej historii rozmowy okna dostaje ten sam wyjątek co kasowanie sesji i mówi wprost, ile pozycji zniknie, biorąc liczbę z ostatniego odczytu historii, a nie z wykazu widocznego na ekranie.
 */
export interface WskazanieCzyszczenia {
  /** Okno, którego historia ginie — pokazywane, bo to ono jest bytem czynności. */
  okno: string;
  /** Ile pozycji stoi w wykazie na ekranie. */
  wWidoku: number;
  /** Ile pozycji naliczył rdzeń przy ostatnim odczycie. */
  razem: number;
  /** Czy odczyt doszedł do skutku choć raz — bez niego `razem` nic nie znaczy. */
  odczytany: boolean;
}

/** Wysyłka żądania usunięcia historii podana z zewnątrz panelu; potwierdzenie czynności już wcześniej padło. */
export type WysylkaCzyszczenia = () => Promise<Wynik<HistoryDeleteResponse>>;

/** Czynność nazwana w treści odmowy rdzenia — wzorem stałej czynności usunięcia używanej w torze sesji. */
const CZYNNOSC_CZYSZCZENIA = 'Wyczyszczenie historii rozmowy okna';

// Objaśnienie mówi prawdę o braku kosza: historia rozmowy nie ma kosza, bo dane kasowane są wprost w transakcji, więc przepisanie zdania o koszu sesji byłoby tu kłamstwem o skutku czynności.
const OBJASNIENIE =
  'Czyszczenie zabiera CAŁĄ historię rozmowy tego okna, także pozycje starsze niż widoczne ' +
  'w wykazie. Historia rozmowy nie ma kosza — kosz rdzenia (30 dni) dotyczy usuniętych SESJI, ' +
  'nie wypowiedzi okna. Do usunięcia pojedynczych pozycji służy „Usuń wskazane", która ' +
  'potwierdzenia nie wymaga.';

/**
 * Otwiera potwierdzenie i prowadzi czyszczenie do końca, oddając odpowiedź rdzenia po udanej komendzie albo wartość pustą, gdy Operator lub rdzeń odmówił czynności.
 */
export function otworzCzyszczenieHistorii(
  wskazanie: WskazanieCzyszczenia,
  wyslij: WysylkaCzyszczenia,
): Promise<HistoryDeleteResponse | null> {
  const rama = zlozRame(wskazanie);
  document.body.append(rama.modal);
  // Natywne otwarcie modalu nie działa w sprawdzianach, więc plik ma otwarcie zapasowe wzorem sesji.
  if (typeof rama.modal.showModal === 'function') rama.modal.showModal();
  else rama.modal.open = true;
  rama.anuluj.focus();
  oznaczFaze(rama.modal, rama.pas, 'gotowe');

  return new Promise<HistoryDeleteResponse | null>((rozwiaz) => {
    let wToku = false;

    /** Zdjęcie okna z drzewa; okno ukryte nawarstwiałoby się przy każdym otwarciu. */
    function zdejmij(): boolean {
      // Rdzeń już kasuje — zamknięcie okna nie ma czego odwołać.
      if (wToku) return false;
      if (typeof rama.modal.close === 'function') rama.modal.close();
      rama.modal.open = false;
      rama.modal.remove();
      return true;
    }

    rama.modal.addEventListener('cancel', (zdarzenie) => {
      zdarzenie.preventDefault();
      zdejmij();
    });

    rama.anuluj.addEventListener('click', () => {
      const bylaOdpowiedz = rama.wyczysc.hidden;
      if (!zdejmij()) {
        powiedzDlaczegoNieTeraz(rama, ODMOWA_ZAMKNIECIA);
        return;
      }
      if (!bylaOdpowiedz) rozwiaz(null);
    });

    rama.wyczysc.addEventListener('click', () => {
      if (wToku) {
        powiedzDlaczegoNieTeraz(rama, ODMOWA_POWTORZENIA);
        return;
      }
      wToku = true;
      zapowiedzWywolanie(rama);
      void wyslij().then((wynik) => {
        wToku = false;
        rozwiaz(pokazOdpowiedz(rama, wynik));
      });
    });
  });
}

const ODMOWA_POWTORZENIA =
  'Rdzeń już czyści historię tego okna. Drugie wywołanie nie miałoby czego usunąć — czekamy na jego odpowiedź.';

const ODMOWA_ZAMKNIECIA =
  'Rdzeń już czyści historię tego okna. Zamknięcie okna niczego nie odwoła, a Operator zostałby bez odpowiedzi — czekamy na nią tutaj.';

/**
 * Zdanie o stracie niesie trzy różne prawdy osobnymi zdaniami, bo rdzeń mógł policzyć całość, policzyć zero albo nie policzyć jeszcze nic.
 */
function zdanieStraty(wskazanie: WskazanieCzyszczenia): string {
  const koniec =
    ' Czynności nie da się cofnąć: historia rozmowy nie ma kosza.';
  if (!wskazanie.odczytany) {
    return (
      `Odczyt historii tego okna jeszcze nie wrócił z rdzenia, więc pełnej liczby NIE ZNAMY. ` +
      `W wykazie stoi ${wskazanie.wWidoku} pozycji i jest to liczba DOLNA — zginie ich co ` +
      `najmniej tyle, a najpewniej więcej.${koniec}`
    );
  }
  if (wskazanie.razem === 0) {
    return (
      'Rdzeń nie ma ani jednej zapisanej pozycji historii tego okna — czyszczenie nie ' +
      'zabierze niczego. Odpowiedź rdzenia i tak padnie tutaj, wraz z liczbą usuniętych.'
    );
  }
  return (
    `Zniknie ${wskazanie.razem} pozycji historii tego okna — wszystkie zapisane wypowiedzi, ` +
    `także starsze niż ${wskazanie.wWidoku} widocznych w wykazie.${koniec}`
  );
}

/** Części modalu, po które sięga przebieg czynności czyszczenia — od dialogu przez pas stanu po parę przycisków. */
interface RamaCzyszczenia {
  modal: HTMLDialogElement;
  pas: HTMLElement;
  skutek: HTMLElement;
  ostrzezenie: HTMLElement;
  wyczysc: HTMLButtonElement;
  anuluj: HTMLButtonElement;
}

/**
 * Budowa modalu — skład co do węzła jak w potwierdzeniu sesji: nagłówek
 * z dymkiem, ciało z bytem czynności, ostrzeżeniem, pasem stanów i akapitem
 * skutku, stopka z odmową i czynnością niebezpieczną.
 */
function zlozRame(wskazanie: WskazanieCzyszczenia): RamaCzyszczenia {
  const modal = document.createElement('dialog');
  modal.className = 'dn-modal dn-usuwanie';

  const naglowek = document.createElement('div');
  naglowek.className = 'dn-modal-naglowek';
  const tytul = document.createElement('h2');
  tytul.className = 'dn-modal-tytul';
  tytul.textContent = 'Wyczyścić całą historię tego okna?';
  naglowek.append(tytul, utworzDymekObjasnienia(OBJASNIENIE));

  // Wykaz jednopozycyjny nazywa okno jako byt czynności, jak tor sesji nazywa tytuły w wykazie.
  const wykaz = document.createElement('ul');
  wykaz.className = 'dn-usuwanie__wykaz';
  const pozycja = document.createElement('li');
  pozycja.className = 'dn-usuwanie__pozycja';
  pozycja.textContent = `Historia rozmowy okna ${wskazanie.okno}`;
  pozycja.dataset['okno'] = wskazanie.okno;
  wykaz.append(pozycja);

  const ostrzezenie = document.createElement('p');
  ostrzezenie.className = 'dn-usuwanie__ostrzezenie';
  ostrzezenie.textContent = zdanieStraty(wskazanie);

  const pas = document.createElement('p');
  pas.className = 'dn-pusty-stan dn-usuwanie__stan';
  pas.hidden = true;

  const skutek = document.createElement('p');
  skutek.className = 'dn-usuwanie__skutek';
  skutek.hidden = true;
  skutek.setAttribute('role', 'status');

  const cialo = document.createElement('div');
  cialo.className = 'dn-modal-cialo';
  cialo.append(wykaz, ostrzezenie, pas, skutek);

  const anuluj = document.createElement('button');
  anuluj.type = 'button';
  anuluj.className = 'dn-btn dn-btn--zarys';
  // Napis mówi, co się stanie po naciśnięciu, a nie że okno zniknie — jak
  // „Zostaw sesje" w torze sesji.
  anuluj.textContent = 'Zostaw historię';

  const wyczysc = document.createElement('button');
  wyczysc.type = 'button';
  wyczysc.className = 'dn-btn dn-btn--niebezpieczny';
  wyczysc.textContent = 'Wyczyść trwale';

  const stopka = document.createElement('div');
  stopka.className = 'dn-modal-stopka';
  stopka.append(anuluj, wyczysc);

  modal.append(naglowek, cialo, stopka);
  return { modal, pas, skutek, ostrzezenie, wyczysc, anuluj };
}

/**
 * Stan ładowania: czynność biegnie, oba przyciski zostają klikalne, a przed powtórzeniem broni strażnik pracy w toku, nie blokada przycisku.
 */
function zapowiedzWywolanie(rama: RamaCzyszczenia): void {
  zapowiedzPowod(rama.wyczysc, ODMOWA_POWTORZENIA);
  zapowiedzPowod(rama.anuluj, ODMOWA_ZAMKNIECIA);
  rama.pas.textContent = 'Rdzeń czyści historię rozmowy tego okna…';
  oznaczFaze(rama.modal, rama.pas, 'ladowanie');
}

/** Powód bezskuteczności przycisku zapowiedziany pod kursorem myszy i czytnikowi ekranu, nigdy samą blokadą przycisku. */
function zapowiedzPowod(kontrolka: HTMLButtonElement, powod: string): void {
  kontrolka.title = powod;
  kontrolka.setAttribute('aria-description', powod);
}

/** Zdjęcie zapowiedzi bezskuteczności po odpowiedzi rdzenia — przyciski odzyskują wtedy realny skutek dla Operatora. */
function zdejmijPowod(kontrolka: HTMLButtonElement): void {
  kontrolka.title = '';
  kontrolka.removeAttribute('aria-description');
}

/** Odpowiedź na kliknięcie przycisku, który w danej chwili nie ma czego wykonać, wyjaśniona wprost Operatorowi. */
function powiedzDlaczegoNieTeraz(rama: RamaCzyszczenia, powod: string): void {
  rama.skutek.textContent = powod;
  rama.skutek.hidden = false;
}

/**
 * Odpowiedź rdzenia na pasie stanu. Modal zostaje otwarty w obu razach, żeby
 * Operator przeczytał, co się stało — także przy odmowie, która nie ma prawa
 * wyglądać jak wykonanie.
 */
function pokazOdpowiedz(
  rama: RamaCzyszczenia,
  wynik: Wynik<HistoryDeleteResponse>,
): HistoryDeleteResponse | null {
  rama.wyczysc.hidden = true;
  zdejmijPowod(rama.anuluj);
  zdejmijPowod(rama.wyczysc);
  rama.skutek.hidden = true;
  rama.skutek.textContent = '';
  rama.anuluj.textContent = 'Zamknij';
  rama.anuluj.focus();

  const odpowiedz = wynik.wynik;
  if (!wynik.udany || odpowiedz === undefined) {
    rama.pas.textContent = opisOdmowyBledu(CZYNNOSC_CZYSZCZENIA, wynik.blad);
    oznaczFaze(rama.modal, rama.pas, 'blad');
    return null;
  }

  rama.ostrzezenie.hidden = true;
  const zdanie = `Usunięto ${odpowiedz.deleted} pozycji historii tego okna.`;

  // Nic nie zginęło — to stan pusty czynności, a ostały się zapis nie wygląda jak skasowany.
  if (odpowiedz.deleted === 0) {
    rama.pas.textContent = `${zdanie} Nie było czego usuwać — historia tego okna była już pusta.`;
    oznaczFaze(rama.modal, rama.pas, 'puste');
    return odpowiedz;
  }

  rama.skutek.textContent = zdanie;
  rama.skutek.hidden = false;
  oznaczFaze(rama.modal, rama.pas, 'gotowe');
  return odpowiedz;
}
