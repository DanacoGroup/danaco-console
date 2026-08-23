import { utworzMagistrale, type Odsubskrybuj } from '../polaczenie/magistrala-zdarzen';

/**
 * Zestaw paneli otwartych obok jednej rozmowy.
 *
 * Przełącznik paneli stoi w nagłówku okna rozmowy, więc każde okno ma własny
 * zestaw paneli: panel należy do rozmowy, nie do ekranu.
 *
 * Moduł nie ma żadnej zmiennej na poziomie modułu — jeden egzemplarz stanu na
 * moduł znaczyłby jeden zestaw paneli na całą aplikację. Cały stan siedzi
 * w domknięciu fabryki, tak jak w `sterowanie/stan-sterowania.ts`; cztery
 * gniazda tworzą cztery niezależne egzemplarze i żaden nie widzi pozostałych.
 *
 * Moduł nie buduje paneli, nie zna ich nazw i nie sprawdza, czy kod pozycji
 * jest otwieralny — trzyma wyłącznie zbiór kodów i dwie liczby układu.
 * Rozstrzyganie, co da się otworzyć, należy do spisu okien pomocniczych.
 *
 * Otwarcie panelu już otwartego nie jest błędem, a wartość spoza zakresu jest
 * przycinana zamiast wstrzymywać wykonanie: powtórne otwarcie nic nie zmienia
 * i nie budzi subskrybentów.
 */
export interface StanPaneli {
  /** Kody paneli otwartych, w kolejności otwierania. */
  otwarte(): readonly string[];
  czyOtwarty(kod: string): boolean;
  otworz(kod: string): void;
  zamknij(kod: string): void;
  przelacz(kod: string): void;
  /** Kod panelu pokazanego na pełnym ekranie; `null` = nikt. */
  pelnyEkran(): string | null;
  ustawPelnyEkran(kod: string | null): void;
  /** Żądana szerokość kolumny paneli w pikselach. */
  szerokosc(): number;
  ustawSzerokosc(px: number): void;
  naZmiane(sluchacz: () => void): Odsubskrybuj;
}

/**
 * Tworzy stan paneli jednego gniazda.
 *
 * `naZmiane` woła się raz na zmianę, nie raz na pole: zamknięcie panelu
 * stojącego na pełnym ekranie zdejmuje też pełny ekran, a subskrybent dostaje
 * jedno powiadomienie, nie dwa. Inaczej odbiorca odrysowywałby układ w stanie
 * przejściowym, w którym pełny ekran wskazuje panel już zamknięty.
 */
export function utworzStanPaneli(szerokoscPoczatkowa: number): StanPaneli {
  const zmiany = utworzMagistrale<void>();

  // Kolejność otwierania jest treścią, nie skutkiem ubocznym — kolumna paneli
  // ustawia je w tej kolejności, więc trzyma to tablica, a nie zbiór.
  const otwarte: string[] = [];
  let pelnyEkran: string | null = null;
  let szerokosc = szerokoscPoczatkowa;

  function oglos(): void {
    zmiany.oglos(undefined);
  }

  function otworz(kod: string): void {
    if (otwarte.includes(kod)) return;
    otwarte.push(kod);
    oglos();
  }

  function zamknij(kod: string): void {
    const miejsce = otwarte.indexOf(kod);
    if (miejsce < 0) return;
    otwarte.splice(miejsce, 1);
    // Panel zamknięty nie może zostać pełnym ekranem — pełny ekran bez panelu
    // byłby wskazaniem na nieistniejący byt.
    if (pelnyEkran === kod) pelnyEkran = null;
    oglos();
  }

  return {
    otwarte: () => [...otwarte],

    czyOtwarty: (kod) => otwarte.includes(kod),

    otworz,

    zamknij,

    przelacz(kod) {
      if (otwarte.includes(kod)) zamknij(kod);
      else otworz(kod);
    },

    pelnyEkran: () => pelnyEkran,

    ustawPelnyEkran(kod) {
      // Na pełny ekran idzie wyłącznie panel, który stoi. Kod spoza zestawu
      // jest przycinany do „nikt" — nie jest to błąd i nie wstrzymuje pracy.
      const nowy = kod !== null && otwarte.includes(kod) ? kod : null;
      if (nowy === pelnyEkran) return;
      pelnyEkran = nowy;
      oglos();
    },

    szerokosc: () => szerokosc,

    ustawSzerokosc(px) {
      if (px === szerokosc) return;
      szerokosc = px;
      oglos();
    },

    naZmiane: (sluchacz) => zmiany.subskrybuj(sluchacz),
  };
}
