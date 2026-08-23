import { czyTrasa, TRASA_POCZATKOWA, Trasa } from './trasy';

/** Widok najwyższego rzędu obsługujący jedną trasę. */
export interface WidokTrasy {
  /** Element widoku; router wstawia go do gospodarza raz, przy pierwszym wejściu. */
  element: HTMLElement;
  /** Wywoływane przy każdym wejściu na trasę, także powrotnym. */
  przyWejsciu?(): void;
}

/** Budowniczy widoku — wywoływany raz, przy pierwszym wejściu na trasę. */
export type BudowniczyWidoku = () => WidokTrasy;

/** Odbiorca zmiany trasy. */
export type SluchaczTrasy = (trasa: Trasa) => void;

export interface Router {
  /** Dopisuje budowniczego widoku obsługującego trasę. */
  zarejestruj(trasa: Trasa, budowniczy: BudowniczyWidoku): void;
  /** Przechodzi na wskazaną trasę. */
  pokaz(trasa: Trasa): void;
  /** Buduje widok trasy, nie pokazując go — po to, by zdążył się przygotować. */
  przygotuj(trasa: Trasa): void;
  /** Trasa obecnie pokazywana. */
  biezaca(): Trasa;
  /** Subskrypcja zmiany trasy. */
  naZmiane(sluchacz: SluchaczTrasy): void;
  /** Otwiera trasę zapisaną w adresie dokumentu albo trasę początkową. */
  uruchom(): void;
}

/**
 * Przełączanie widoków najwyższego rzędu — bez biblioteki zewnętrznej.
 *
 * Jedna odpowiedzialność: rozstrzygnięcie, który widok jest na wierzchu.
 * Router nie zna żadnego widoku z osobna; przyjmuje budowniczych i wywołuje
 * ich najwyżej raz.
 *
 * Widok opuszczony nie znika: element zostaje w dokumencie z atrybutem
 * `hidden`, zamiast być usuwanym i budowanym na nowo. Dzięki temu karty
 * opuszczonego środowiska trwają w tle i odtwarzają pełny stan po powrocie,
 * a żadna subskrypcja kanału się nie gubi — powrót nie zakłada drugiego
 * okna komunikacji.
 *
 * Trasa widnieje w adresie dokumentu, więc odświeżenie strony wraca tam,
 * gdzie Operator był, a przycisk „wstecz" przeglądarki działa bez kodu
 * dodatkowego. Adres nieznany nie zatrzymuje uruchomienia.
 */
export function utworzRouter(gospodarz: HTMLElement): Router {
  const budowniczowie = new Map<Trasa, BudowniczyWidoku>();
  const widoki = new Map<Trasa, WidokTrasy>();
  const sluchacze: SluchaczTrasy[] = [];

  let biezaca: Trasa = TRASA_POCZATKOWA;
  /** Chroni przed zapętleniem: zmiana adresu wywołana przez sam router. */
  let wlasnaZmianaAdresu = false;

  /** Widok trasy; budowany przy pierwszym wejściu, potem brany z pamięci. */
  function widok(trasa: Trasa): WidokTrasy | null {
    const zbudowany = widoki.get(trasa);
    if (zbudowany !== undefined) return zbudowany;

    const budowniczy = budowniczowie.get(trasa);
    if (budowniczy === undefined) return null;

    const nowy = budowniczy();
    nowy.element.classList.add('dn-aplikacja__widok');
    nowy.element.hidden = true;
    gospodarz.append(nowy.element);
    widoki.set(trasa, nowy);
    return nowy;
  }

  function pokaz(trasa: Trasa): void {
    const wybrany = widok(trasa);
    if (wybrany === null) return;

    for (const [nazwa, inny] of widoki) inny.element.hidden = nazwa !== trasa;

    biezaca = trasa;
    gospodarz.dataset.trasa = trasa;
    zapiszAdres(trasa);

    wybrany.przyWejsciu?.();
    for (const sluchacz of [...sluchacze]) sluchacz(trasa);
  }

  /** Zapisuje trasę w adresie dokumentu bez wywoływania własnej obsługi. */
  function zapiszAdres(trasa: Trasa): void {
    if (odczytajAdres() === trasa) return;
    wlasnaZmianaAdresu = true;
    window.location.hash = `#${trasa}`;
  }

  window.addEventListener('hashchange', () => {
    if (wlasnaZmianaAdresu) {
      wlasnaZmianaAdresu = false;
      return;
    }
    const trasa = odczytajAdres();
    if (trasa !== null && trasa !== biezaca) pokaz(trasa);
  });

  return {
    zarejestruj: (trasa, budowniczy) => void budowniczowie.set(trasa, budowniczy),
    pokaz,
    przygotuj: (trasa) => void widok(trasa),
    biezaca: () => biezaca,
    naZmiane: (sluchacz) => void sluchacze.push(sluchacz),
    uruchom: () => pokaz(odczytajAdres() ?? TRASA_POCZATKOWA),
  };
}

/** Trasa zapisana w adresie dokumentu; `null`, gdy adres jej nie niesie. */
function odczytajAdres(): Trasa | null {
  const nazwa = window.location.hash.replace(/^#/u, '');
  return czyTrasa(nazwa) ? nazwa : null;
}
