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

/** Wysyłka `session.delete` podana z zewnątrz; potwierdzenie już padło. */
export type WysylkaUsuniecia = (
  idSesji: readonly string[],
) => Promise<Wynik<RozliczenieUsuniecia>>;

export type { WskazanieUsuniecia };

/**
 * Potwierdzenie trwałego usunięcia sesji — jedna powierzchnia całej czynności.
 *
 * Jedna odpowiedzialność: pokazać, co zginie, przyjąć potwierdzenie i pokazać,
 * co rdzeń odpowiedział. Nazwy komendy ten plik nie zna — wysyłkę dostaje
 * z zewnątrz (`usun`), więc daje się poddać próbie bez rdzenia.
 *
 * Operator widzi stratę, zanim kliknie. Kontrakt żąda pola `confirm`, a rdzeń
 * bez niego odmawia wykonania (`adapter_sesje_usuwanie.go`). Potwierdzenie ma
 * więc treść, a nie samo pytanie „czy jesteś pewien": modal wypisuje tytuły
 * wskazanych sesji i mówi wprost, że zapis ginie razem z wiadomościami, oknami
 * i artefaktami.
 *
 * Okno stoi na natywnym `<dialog>` (wzorem `strona-glowna/pytanie-o-nazwe.ts`),
 * więc warstwa tła, pułapka ogniska i Escape należą do przeglądarki. Escape
 * w trakcie wywołania jest wstrzymany: rdzeń już usuwa sesje, więc zamknięcie
 * okna nie odwołałoby niczego, a Operator zostałby bez odpowiedzi.
 *
 * Trzy stany obowiązkowe stoją na jednym pasie `oznaczFaze` z biblioteki:
 * `ladowanie` w czasie wywołania, `puste` gdy rdzeń nie usunął niczego, `blad`
 * przy odmowie — z kodem i treścią rdzenia z `opisOdmowyBledu`. Odmowa nie
 * zamyka modalu: przycisk czynności znika, zostaje samo „Zamknij", więc nie da
 * się wziąć odmowy za skutek.
 */

/**
 * Otwiera potwierdzenie i prowadzi czynność do końca.
 *
 * Oddaje rozliczenie rdzenia, gdy komenda przeszła — także rozliczenie, w
 * którym nic nie zginęło, bo to też jest odpowiedź rdzenia i pas kart ma ją
 * powtórzyć. Oddaje `null`, gdy Operator odmówił potwierdzenia albo gdy rdzeń
 * odmówił wykonania; treść odmowy została wtedy pokazana w modalu.
 *
 * Obietnica rozstrzyga się z odpowiedzią rdzenia, nie z zamknięciem okna. Modal
 * zostaje otwarty do przeczytania skutku, ale pas kart ma odpowiedź
 * natychmiast — wiązanie rozstrzygnięcia z zamknięciem okna kazałoby pasowi
 * milczeć tak długo, jak długo Operator czyta.
 */
export function otworzUsuniecieSesji(
  wskazania: readonly WskazanieUsuniecia[],
  wyslij: WysylkaUsuniecia,
  tytul: TytulSesji,
): Promise<RozliczenieUsuniecia | null> {
  const rama = zlozRameUsuniecia(wskazania);
  const { modal, pas } = rama;
  document.body.append(modal);
  // `showModal` daje nakładkę, stos okien, pułapkę ogniska i Escape. Środowisko
  // sprawdzianów DOM go nie implementuje (wzorem `moduly/design/modal-kreatora.ts`),
  // a treść potwierdzenia ma być mierzalna także tam — stąd otwarcie zapasowe.
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
      // „Zostaw sesje” przed komendą to odmowa Operatora; „Zamknij” po
      // odpowiedzi rdzenia niczego już nie rozstrzyga — rozliczenie poszło
      // do pasa w chwili, gdy przyszło.
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

/**
 * Stan `ladowanie`: czynność biegnie, oba przyciski zostają klikalne.
 *
 * Przed powtórzeniem czynności broni strażnik `wToku` w obsłudze kliknięcia,
 * a przed zamknięciem okna ten sam strażnik w `zdejmij`, więc `disabled` nie
 * dokładałby ochrony — dokładałby ciszę. Zamiast tego przycisk odpowiada
 * zdaniem, dlaczego w tej chwili nie ma czego zrobić.
 *
 * Powód idzie dwiema drogami, wzorem `przyciskBezKomendy`: `title` pod kursorem
 * i `aria-description` dla czytnika ekranu.
 */
function zapowiedzWywolanie(rama: RamaUsuniecia): void {
  zapowiedzPowod(rama.usun, ODMOWA_POWTORZENIA);
  zapowiedzPowod(rama.anuluj, ODMOWA_ZAMKNIECIA);
  rama.pas.textContent = 'Rdzeń usuwa wskazane sesje…';
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

/**
 * Odpowiedź na kliknięcie, które w tej chwili nie ma czego wykonać.
 *
 * Zdanie idzie do akapitu skutku (`role="status"`), więc czytnik ekranu
 * ogłasza je od razu, a pas stanu dalej trzyma fazę `ladowanie`. Milczenie
 * byłoby tu gorsze od blokady: Operator wziąłby brak reakcji za zawieszenie.
 */
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
  // Odpowiedź przyszła: zapowiedź bezskuteczności przestaje być prawdziwa,
  // a zdanie o niej nie ma prawa zostać pod odpowiedzią rdzenia.
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

  // Nic nie zginęło — to stan pusty czynności i tak jest nazwany. Zapis, który
  // się ostał, nie ma wyglądać jak zapis skasowany.
  if (czyRozliczeniePuste(rozliczenie)) {
    rama.pas.textContent = tresc;
    oznaczFaze(rama.modal, rama.pas, 'puste');
    return rozliczenie;
  }

  // Powodzenie zdejmuje pas stanów, a zdanie o stracie zostaje we własnym
  // akapicie do chwili zamknięcia okna — Operator ma je przeczytać.
  rama.skutek.textContent = tresc;
  rama.skutek.hidden = false;
  oznaczFaze(rama.modal, rama.pas, 'gotowe');
  return rozliczenie;
}
