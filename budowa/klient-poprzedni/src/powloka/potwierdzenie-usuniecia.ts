import { oznaczFaze } from '../komponenty/faza-okna';
import { opisOdmowyBledu } from '../komponenty/odmowa';
import type { Wynik } from '../protokol/kanal';
import {
  czyRozliczeniePuste,
  zdanieNieznalezionych,
  zdanieUsunietych,
  type TytulSesji,
} from './rozliczenie-usuniecia';
import {
  zlozRameUsuniecia,
  type RamaUsuniecia,
  type WskazanieUsuniecia,
} from './rama-usuniecia';
import { CZYNNOSC_USUNIECIA, type RozliczenieUsuniecia } from './usuniecie-sesji';

/** Wysyłka `session.delete` podana z zewnątrz, wywoływana dopiero po tym, jak Operator potwierdzi usunięcie. */
export type WysylkaUsuniecia = (
  idSesji: readonly string[],
) => Promise<Wynik<RozliczenieUsuniecia>>;

export type { WskazanieUsuniecia };

// Potwierdzenie trwałego usunięcia sesji: pokazuje, co zginie, i pokazuje, co odpowiedział rdzeń.

/** Otwiera potwierdzenie i prowadzi czynność usunięcia sesji do końca, zwracając rozliczenie samego rdzenia. */
export function otworzUsuniecieSesji(
  wskazania: readonly WskazanieUsuniecia[],
  wyslij: WysylkaUsuniecia,
  tytul: TytulSesji,
): Promise<RozliczenieUsuniecia | null> {
  const rama = zlozRameUsuniecia(wskazania);
  const { modal, pas } = rama;
  document.body.append(modal);
  // `showModal` daje nakładkę, stos okien i Escape; brak w sprawdzianach DOM ma otwarcie zapasowe.
  if (typeof modal.showModal === 'function') modal.showModal();
  else modal.open = true;
  rama.anuluj.focus();
  oznaczFaze(modal, pas, 'gotowe');

  return new Promise<RozliczenieUsuniecia | null>((rozwiaz) => {
    let wToku = false;

    /** Zdjęcie okna z drzewa; okno ukryte nawarstwiałoby się przy każdej karcie. */
    function zdejmij(): boolean {
      // Rdzeń już usuwa sesje — zamknięcie okna nie ma czego odwołać.
      if (wToku) return false;
      if (typeof modal.close === 'function') modal.close();
      modal.open = false;
      modal.remove();
      return true;
    }

    modal.addEventListener('cancel', (zdarzenie) => {
      zdarzenie.preventDefault();
      zdejmij();
    });

    rama.anuluj.addEventListener('click', () => {
      const bylOdpowiedzia = rama.usun.hidden;
      if (!zdejmij()) {
        powiedzDlaczegoNieTeraz(rama, ODMOWA_ZAMKNIECIA);
        return;
      }
      // „Zostaw sesje” to odmowa Operatora; „Zamknij” po odpowiedzi rdzenia niczego już nie rozstrzyga.
      if (!bylOdpowiedzia) rozwiaz(null);
    });

    rama.usun.addEventListener('click', () => {
      if (wToku) {
        powiedzDlaczegoNieTeraz(rama, ODMOWA_POWTORZENIA);
        return;
      }
      wToku = true;
      zapowiedzWywolanie(rama);
      void wyslij(wskazania.map((wskazanie) => wskazanie.id)).then((wynik) => {
        wToku = false;
        rozwiaz(pokazOdpowiedz(rama, wynik, tytul));
      });
    });
  });
}

const ODMOWA_POWTORZENIA =
  'Rdzeń już usuwa wskazane sesje. Drugie wywołanie nie miałoby czego usunąć — czekamy na jego odpowiedź.';

const ODMOWA_ZAMKNIECIA =
  'Rdzeń już usuwa wskazane sesje. Zamknięcie okna niczego nie odwoła, a Operator zostałby bez odpowiedzi — czekamy na nią tutaj.';

/** Stan `ladowanie`: czynność biegnie, oba przyciski zostają klikalne i odpowiadają zdaniem zamiast ciszą. */
function zapowiedzWywolanie(rama: RamaUsuniecia): void {
  zapowiedzPowod(rama.usun, ODMOWA_POWTORZENIA);
  zapowiedzPowod(rama.anuluj, ODMOWA_ZAMKNIECIA);
  rama.pas.textContent = 'Rdzeń usuwa wskazane sesje…';
  oznaczFaze(rama.modal, rama.pas, 'ladowanie');
}

/** Powód bezskuteczności przycisku, zapowiedziany pod kursorem i czytnikowi ekranu, a nie blokadą kontrolki. */
function zapowiedzPowod(kontrolka: HTMLButtonElement, powod: string): void {
  kontrolka.title = powod;
  kontrolka.setAttribute('aria-description', powod);
}

/** Zdjęcie zapowiedzi bezskuteczności po odpowiedzi rdzenia, gdy przyciski znów mają jakikolwiek skutek. */
function zdejmijPowod(kontrolka: HTMLButtonElement): void {
  kontrolka.title = '';
  kontrolka.removeAttribute('aria-description');
}

/** Odpowiedź na kliknięcie kontrolki, która w tej chwili nie ma czego wykonać — zdanie idzie do akapitu skutku. */
function powiedzDlaczegoNieTeraz(rama: RamaUsuniecia, powod: string): void {
  rama.skutek.textContent = powod;
  rama.skutek.hidden = false;
}

/**
 * Odpowiedź rdzenia na pasie stanu. Zwraca rozliczenie albo `null` przy
 * odmowie — modal zostaje otwarty w obu razach, żeby Operator przeczytał, co
 * się stało, zanim okno zniknie.
 */
function pokazOdpowiedz(
  rama: RamaUsuniecia,
  wynik: Wynik<RozliczenieUsuniecia>,
  tytul: TytulSesji,
): RozliczenieUsuniecia | null {
  rama.usun.hidden = true;
  // Odpowiedź przyszła: zapowiedź bezskuteczności przestaje być prawdziwa i znika spod odpowiedzi.
  zdejmijPowod(rama.anuluj);
  zdejmijPowod(rama.usun);
  rama.skutek.hidden = true;
  rama.skutek.textContent = '';
  rama.anuluj.textContent = 'Zamknij';
  rama.anuluj.focus();

  const rozliczenie = wynik.wynik;
  if (!wynik.udany || rozliczenie === undefined) {
    rama.pas.textContent = opisOdmowyBledu(CZYNNOSC_USUNIECIA, wynik.blad);
    oznaczFaze(rama.modal, rama.pas, 'blad');
    return null;
  }

  rama.ostrzezenie.hidden = true;
  const pominiete = zdanieNieznalezionych(rozliczenie, tytul);
  const zdanie = zdanieUsunietych(rozliczenie, tytul);
  const tresc = pominiete === null ? zdanie : `${zdanie} ${pominiete}`;

  // Nic nie zginęło — to stan pusty czynności, zapis, który się ostał, nie ma wyglądać jak skasowany.
  if (czyRozliczeniePuste(rozliczenie)) {
    rama.pas.textContent = tresc;
    oznaczFaze(rama.modal, rama.pas, 'puste');
    return rozliczenie;
  }

  // Powodzenie zdejmuje pas stanów, zdanie o stracie zostaje we własnym akapicie do zamknięcia okna.
  rama.skutek.textContent = tresc;
  rama.skutek.hidden = false;
  oznaczFaze(rama.modal, rama.pas, 'gotowe');
  return rozliczenie;
}
