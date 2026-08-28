import { utworzMagistrale, type Odsubskrybuj } from '../polaczenie/magistrala-zdarzen';

/**
 * Zestaw paneli otwartych obok jednej rozmowy stoi w domknięciu fabryki, jedno na gniazdo, trzymając wyłącznie zbiór kodów i dwie liczby układu, bez budowania paneli ani sprawdzania, co da się otworzyć.
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
 * Tworzy stan paneli jednego gniazda; zmiana woła się raz na zmianę, nie raz na pole, żeby subskrybent nie odrysowywał układu w stanie przejściowym.
 */
export function utworzStanPaneli(szerokoscPoczatkowa: number): StanPaneli {
  const zmiany = utworzMagistrale<void>();

  // Kolejność otwierania jest treścią, nie skutkiem ubocznym, dlatego trzyma ją tablica, a nie zbiór.
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
    // Panel zamknięty nie może zostać pełnym ekranem — to byłoby wskazaniem na nieistniejący byt.
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
      // Na pełny ekran idzie wyłącznie panel, który stoi; kod spoza zestawu jest przycinany, to nie błąd.
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
