import { elementIkony } from '../ikony/ikony';
import { utworzMagistrale, type Odsubskrybuj } from '../polaczenie/magistrala-zdarzen';

/** Liczba ustawień okna pokazywana przy uchwycie szuflady jako plakietka sygnałowa z liczbą pozycji podsumowania. */
const LICZBA_USTAWIEN = 8;

/** Szuflada sterowania jest oprawą czyniącą kontrolki widocznymi: zwinięta pokazuje podsumowanie wartości, rozwinięta — komplet kontrolek, a stan ogłasza magistralą wspólną dla obu uchwytów. */
export interface Szuflada {
  /** Element montowany w kolumnie sterowania. */
  element: HTMLElement;
  /** Czy komplet kontrolek jest rozwinięty. */
  otwarta(): boolean;
  /** Ustawia stan rozwinięcia. */
  ustaw(otwarta: boolean): void;
  /** Odwraca stan rozwinięcia. */
  przelacz(): void;
  /** Subskrypcja zmian stanu rozwinięcia. */
  naZmiane(sluchacz: (otwarta: boolean) => void): Odsubskrybuj;
}

/** Warstwy szuflady: podsumowanie widoczne po zwinięciu oraz komplet kontrolek widoczny po rozwinięciu, powiązane wspólnym identyfikatorem okna. */
export interface WarstwySzuflady {
  /** Podsumowanie ośmiu wartości — warstwa zwiniętej szuflady. */
  podsumowanie: HTMLElement;
  /** Komplet kontrolek okna — warstwa rozwiniętej szuflady. */
  tresc: HTMLElement;
  /** Identyfikator okna; wchodzi w identyfikatory dostępności. */
  idOkna: string;
  /** Stan rozwinięcia w chwili zamontowania. */
  otwartaNaStarcie: boolean;
}

export function utworzSzuflade(warstwy: WarstwySzuflady): Szuflada {
  const zmiany = utworzMagistrale<boolean>();
  const idWarstw = `dc-widok-ster-warstwy-${warstwy.idOkna}`;
  let otwarta = warstwy.otwartaNaStarcie;

  const element = document.createElement('div');
  element.className = 'dc-widok-ster__szuflada';

  const uchwyt = document.createElement('button');
  uchwyt.type = 'button';
  uchwyt.className = 'dc-widok-ster__uchwyt';
  uchwyt.setAttribute('aria-controls', idWarstw);

  const nazwa = document.createElement('span');
  nazwa.className = 'dc-widok-ster__uchwyt-nazwa';

  // Licznik na plakietce sygnałowej (`--sygnal`).
  const licznik = document.createElement('span');
  licznik.className = 'dn-plakietka dn-plakietka--sygnal dc-widok-ster__licznik';
  licznik.textContent = String(LICZBA_USTAWIEN);

  const strzalka = elementIkony('strzalka-dol', {
    rozmiar: 16,
    klasa: 'dn-ikona dc-widok-ster__strzalka',
  });

  uchwyt.append(
    elementIkony('ustawienia', { rozmiar: 18 }),
    nazwa,
    licznik,
    strzalka,
  );

  const warstwyElement = document.createElement('div');
  warstwyElement.className = 'dc-widok-ster__warstwy';
  warstwyElement.id = idWarstw;
  warstwyElement.append(warstwy.podsumowanie, warstwy.tresc);

  element.append(uchwyt, warstwyElement);

  uchwyt.addEventListener('click', () => ustaw(!otwarta));

  // Klawisz wyjścia zwija szufladę i wraca ogniskiem na uchwyt, bez trwałego zamknięcia.
  element.addEventListener('keydown', (zdarzenie) => {
    if (zdarzenie.key !== 'Escape' || !otwarta) return;
    ustaw(false);
    uchwyt.focus();
  });

  function ustaw(nowa: boolean): void {
    otwarta = nowa;
    odrysuj();
    zmiany.oglos(otwarta);
  }

  function odrysuj(): void {
    element.dataset.otwarta = otwarta ? 'tak' : 'nie';
    uchwyt.setAttribute('aria-expanded', String(otwarta));
    nazwa.textContent = otwarta ? 'Ustawienia okna' : 'Rozwiń ustawienia okna';
    uchwyt.title = otwarta
      ? 'Zwiń komplet kontrolek — wartości zostaną w podsumowaniu'
      : 'Rozwiń komplet ośmiu kontrolek sterowania oknem';
  }

  odrysuj();

  return {
    element,
    otwarta: () => otwarta,
    ustaw,
    przelacz: () => ustaw(!otwarta),
    naZmiane: (sluchacz) => zmiany.subskrybuj(sluchacz),
  };
}
