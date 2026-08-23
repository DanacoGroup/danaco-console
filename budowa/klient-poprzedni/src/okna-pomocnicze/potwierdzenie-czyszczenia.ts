// Arkusz potwierdzenia usunięcia bierzemy wprost z toru sesji, zamiast pisać
// drugi. Rodzina `.dn-usuwanie__*` jest biblioteczna (przedrostek `dn-`), a
// jedyny jej arkusz stoi przy potwierdzeniu sesji; dwa arkusze o jednej
// czynności rozjechałyby się przy pierwszej poprawce.
import '../powloka/usuniecie-sesji.css';

import type { HistoryDeleteResponse } from '../../../shared/contract';
import { utworzDymekObjasnienia } from '../komponenty/dymek';
import { oznaczFaze } from '../komponenty/faza-okna';
import { opisOdmowyBledu } from '../komponenty/odmowa';
import type { Wynik } from '../protokol/kanal';

/**
 * Potwierdzenie wyczyszczenia całej historii rozmowy okna.
 *
 * Czyszczenie całej historii okna dostaje ten sam wyjątek, który ma kasowanie
 * sesji — i tylko ono. „Usuń wskazane" potwierdzenia nie ma: Operator wskazał
 * pozycje własną ręką, więc wskazanie jest zgodą.
 *
 * `powloka/potwierdzenie-usuniecia.ts` nie da się tu użyć wprost, bo jest
 * związany z sesją kształtem: przyjmuje `WskazanieUsuniecia[]` (identyfikator
 * i tytuł sesji), buduje wykaz „jedna pozycja na sesję", oddaje
 * `RozliczenieUsuniecia` po identyfikatorach i rozlicza je zdaniami
 * z `rozliczenie-usuniecia.ts`, które odmieniają słowo „sesja". Czyszczenie
 * historii nie ma ani wykazu bytów (okno jest jedno), ani rozliczenia po
 * identyfikatorach — rdzeń oddaje samą liczbę `deleted`.
 *
 * Wygląd jest za to powielony co do znaku: ten sam natywny `<dialog>`, ta sama
 * rama biblioteki (`.dn-modal*`), ta sama rodzina klas `.dn-usuwanie__*` z tego
 * samego arkusza, ten sam znak objaśnienia w nagłówku, ten sam pas stanów na
 * `oznaczFaze`, ta sama para przycisków i ta sama zasada „odmowa rdzenia nie
 * zamyka modalu".
 *
 * Potwierdzenie mówi, co zniknie, wraz z liczbą — pytanie „czy na pewno?" bez
 * liczby jest klikane odruchowo. Liczba pochodzi z `total` ostatniego
 * `history.load`, czyli z rdzenia, a nie z długości wykazu na ekranie. Gdy
 * odczytu jeszcze nie było, modal mówi to wprost i podaje liczbę widoczną jako
 * dolną — tą samą granicą, którą niesie zapowiedź retencji
 * (`nastawa-retencji.ts`).
 */

/** Co panel wie o rozmiarze straty w chwili otwarcia potwierdzenia. */
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

/** Wysyłka `history.delete` podana z zewnątrz; potwierdzenie już padło. */
export type WysylkaCzyszczenia = () => Promise<Wynik<HistoryDeleteResponse>>;

/** Czynność nazwana w odmowie rdzenia — wzorem `CZYNNOSC_USUNIECIA` toru sesji. */
const CZYNNOSC_CZYSZCZENIA = 'Wyczyszczenie historii rozmowy okna';

// Objaśnienie mówi prawdę o koszu, a raczej o jego braku. Sesja ma kosz na
// trzydzieści dni i potwierdzenie sesji o tym mówi. Historia rozmowy kosza nie
// ma: `dane/historia.go` kasuje wiersze wprost w transakcji i sprząta po nich
// bloki, więc przepisanie tamtego zdania byłoby tu kłamstwem o skutku.
const OBJASNIENIE =
  'Czyszczenie zabiera CAŁĄ historię rozmowy tego okna, także pozycje starsze niż widoczne ' +
  'w wykazie. Historia rozmowy nie ma kosza — kosz rdzenia (30 dni) dotyczy usuniętych SESJI, ' +
  'nie wypowiedzi okna. Do usunięcia pojedynczych pozycji służy „Usuń wskazane", która ' +
  'potwierdzenia nie wymaga.';

/**
 * Otwiera potwierdzenie i prowadzi czyszczenie do końca.
 *
 * Oddaje odpowiedź rdzenia, gdy komenda przeszła — także odpowiedź „usunięto
 * 0", bo to też jest odpowiedź i panel ma ją powtórzyć. Oddaje `null`, gdy
 * Operator odmówił albo gdy rdzeń odmówił; treść odmowy została wtedy pokazana
 * w modalu.
 */
export function otworzCzyszczenieHistorii(
  wskazanie: WskazanieCzyszczenia,
  wyslij: WysylkaCzyszczenia,
): Promise<HistoryDeleteResponse | null> {
  const rama = zlozRame(wskazanie);
  document.body.append(rama.modal);
  // `showModal` daje nakładkę, stos okien, pułapkę ogniska i Escape. Środowisko
  // sprawdzianów DOM go nie implementuje — stąd otwarcie zapasowe, wzorem
  // `powloka/potwierdzenie-usuniecia.ts`.
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
 * Zdanie o stracie — treść właściwa całego potwierdzenia.
 *
 * Trzy różne prawdy, trzy różne zdania, bo mieszanie ich zamieniłoby liczbę
 * w ozdobę: rdzeń policzył całość, rdzeń policzył zero, rdzeń nie policzył
 * jeszcze nic. Zdanie ostatnie podaje liczbę widoczną i nazywa ją dolną —
 * tak samo jak zapowiedź retencji nazywa swoją.
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

/** Części modalu, po które sięga przebieg czynności. */
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

  // Wykaz jednopozycyjny — bytem czynności jest okno i ma być nazwane, tak jak
  // w torze sesji nazwane są tytuły sesji. Pusty wykaz zostawiłby stratę bez
  // adresata.
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
 * Stan `ladowanie`: czynność biegnie, a oba przyciski zostają klikalne. Przed
 * powtórzeniem broni strażnik `wToku`, a nie `disabled` — przycisk odpowiada
 * zdaniem, dlaczego nie ma czego zrobić.
 */
function zapowiedzWywolanie(rama: RamaCzyszczenia): void {
  zapowiedzPowod(rama.wyczysc, ODMOWA_POWTORZENIA);
  zapowiedzPowod(rama.anuluj, ODMOWA_ZAMKNIECIA);
  rama.pas.textContent = 'Rdzeń czyści historię rozmowy tego okna…';
  oznaczFaze(rama.modal, rama.pas, 'ladowanie');
}

/** Powód bezskuteczności zapowiedziany pod kursorem i czytnikowi, nie blokadą. */
function zapowiedzPowod(kontrolka: HTMLButtonElement, powod: string): void {
  kontrolka.title = powod;
  kontrolka.setAttribute('aria-description', powod);
}

/** Zdjęcie zapowiedzi po odpowiedzi rdzenia — przyciski znów mają skutek. */
function zdejmijPowod(kontrolka: HTMLButtonElement): void {
  kontrolka.title = '';
  kontrolka.removeAttribute('aria-description');
}

/** Odpowiedź na kliknięcie, które w tej chwili nie ma czego wykonać. */
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

  // Nic nie zginęło — to stan pusty czynności i tak go nazywamy. Zapis, który
  // się ostał, nie ma prawa wyglądać jak zapis skasowany.
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
