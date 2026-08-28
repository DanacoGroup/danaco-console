import { Command } from '../../../../shared/contract';
import type {
  CzynnoscWarsztatu,
  PoleCzynnosci,
  WartosciCzynnosci,
} from '../studio/czynnosci-warsztatu';

// Katalog czynności warsztatów modułu Design zapisuje pięć warsztatów danymi zamiast formularzy.

/** Warsztat, do którego czynność należy. Okno filtruje wykaz czynności katalogu po tym polu i pokazuje wyłącznie czynności własnej grupy warsztatu. */
export type GrupaWarsztatu =
  | 'fotografia'
  | 'wektor'
  | 'druk'
  | 'bazy'
  | 'publikacja'
  | 'makieta';

/** Czynność warsztatu Design wraz z przypisaniem do warsztatu, po którym okno rozdziela wykaz katalogu na osobne listy dla każdego z pięciu warsztatów. */
export interface CzynnoscWarsztatuDesignu extends CzynnoscWarsztatu {
  grupa: GrupaWarsztatu;
}

/** Odczyt pola tekstowego z wartości czynności. Wartość pusta albo złożona z samych odstępów zwraca undefined, a nie pusty napis, co odróżnia brak od pustki. */
function tekst(wartosci: WartosciCzynnosci, kod: string): string | undefined {
  const wartosc = (wartosci[kod] ?? '').trim();
  return wartosc === '' ? undefined : wartosc;
}

/** Odczyt pola liczbowego z wartości czynności. Wartość, której nie da się zamienić na liczbę skończoną, zwraca undefined zamiast liczby niepoprawnej. */
function liczba(wartosci: WartosciCzynnosci, kod: string): number | undefined {
  const wartosc = tekst(wartosci, kod);
  if (wartosc === undefined) return undefined;
  const odczytana = Number(wartosc);
  return Number.isFinite(odczytana) ? odczytana : undefined;
}

/** Odczyt pola przełącznika z wartości czynności. Pole nieustawione zwraca undefined, a nie fałsz, ponieważ brak ustawienia i wyłączenie niosą inne znaczenie. */
function przelacznik(wartosci: WartosciCzynnosci, kod: string): boolean | undefined {
  const wartosc = wartosci[kod];
  if (wartosc === undefined) return undefined;
  return wartosc === 'tak';
}

/** Rozbiór wartości pola listy tekstu oddzielonej przecinkami na wykaz pozycji, z którego usunięte zostają pozycje puste powstałe z odstępów wokół przecinków. */
function wykaz(wartosci: WartosciCzynnosci, kod: string): string[] {
  const wartosc = tekst(wartosci, kod);
  if (wartosc === undefined) return [];
  return wartosc
    .split(',')
    .map((pozycja) => pozycja.trim())
    .filter((pozycja) => pozycja !== '');
}

/** Rozbiór wartości pola listy liczb oddzielonej przecinkami na wykaz liczb. Pozycja, której nie da się zamienić na liczbę skończoną, wypada z wykazu. */
function wykazLiczb(wartosci: WartosciCzynnosci, kod: string): number[] {
  return wykaz(wartosci, kod)
    .map((pozycja) => Number(pozycja))
    .filter((pozycja) => Number.isFinite(pozycja));
}

/** Dokłada do podstawy żądania wyłącznie te pola dodatkowe, których wartość jest podana — pole nieustawione oraz pustą listę pomija, zamiast wysyłać je puste. */
function zPolami(
  podstawa: Record<string, unknown>,
  dodatki: Record<string, unknown>,
): Record<string, unknown> {
  const zadanie = { ...podstawa };
  for (const [nazwa, wartosc] of Object.entries(dodatki)) {
    if (wartosc === undefined) continue;
    if (Array.isArray(wartosc) && wartosc.length === 0) continue;
    zadanie[nazwa] = wartosc;
  }
  return zadanie;
}

/**
 * Podpowiedź pola wykazu czynności — przykładowy zapis, którego Operator się
 * spodziewa; nazwa komendy pochodzi ze stałych kontraktu, nie z napisu
 * wpisanego wprost.
 */
const PODPOWIEDZ_CZYNNOSCI = `[{"command":"${Command.DesignPhotoEnhance}","settings":{"autoLevels":true}}]`;

/** Pole materiału czynności — wskazanie zasobu z magazynu okna. Zasób wskazany polem pozostaje niezmieniony, wynik czynności jest jego osobnym wariantem. */
const POLE_ZASOBU: PoleCzynnosci = {
  kod: 'assetId',
  etykieta: 'Materiał',
  rodzaj: 'zasob',
  opis:
    'Zasób z magazynu okna. Materiał zostaje NIETKNIĘTY — wynik jest jego wariantem, ' +
    'więc pomyłka nie kosztuje oryginału.',
};

/** Pole maski czynności — zasób, którego jasność rozstrzyga o obszarze objętym czynnością: punkt biały znaczy obszar objęty, czarny — obszar pominięty. */
const POLE_MASKI: PoleCzynnosci = {
  kod: 'maskAssetId',
  etykieta: 'Maska',
  rodzaj: 'zasob',
  opis: 'Zasób maski: punkt biały znaczy obszar objęty, czarny — pominięty.',
};

/** Pole kanału modelu obrazowego, stosowane w czynnościach o wariancie neuronowym, dla których wybór kanału rozstrzyga, którą drogą rdzeń liczy wynik. */
const POLE_KANALU: PoleCzynnosci = {
  kod: 'channelId',
  etykieta: 'Kanał modelu obrazowego',
  rodzaj: 'tekst',
  opis:
    'Wskazanie kanału. Odpowiedź mówi polem „policzone przez", którą drogą rdzeń ' +
    'naprawdę policzył wynik — kanałem czy rachunkiem wkompilowanym.',
};

/** Zakresy obszarów zapisane wierszami w postaci x,y,szerokość,wysokość. Wiersz o mniej niż czterech liczbach nie tworzy obszaru i zostaje pominięty. */
function obszaryZWierszy(wartosci: WartosciCzynnosci, kod: string): Record<string, number>[] {
  const wartosc = tekst(wartosci, kod);
  if (wartosc === undefined) return [];
  const obszary: Record<string, number>[] = [];
  for (const wiersz of wartosc.split('\n')) {
    const liczby = wiersz
      .split(',')
      .map((pozycja) => Number(pozycja.trim()))
      .filter((pozycja) => Number.isFinite(pozycja));
    if (liczby.length < 4) continue;
    obszary.push({
      x: liczby[0] as number,
      y: liczby[1] as number,
      width: liczby[2] as number,
      height: liczby[3] as number,
    });
  }
  return obszary;
}

export const CZYNNOSCI_WARSZTATOW_DESIGNU: readonly CzynnoscWarsztatuDesignu[] = [
  // ── Warsztat fotografii ────────────────────────────────────────────────────
  {
    grupa: 'fotografia',
    komenda: Command.DesignPhotoCrop,
    nazwa: 'Kadruj i prostuj',
    opis:
      'Kadr, prostowanie horyzontu i proporcje. Prostowanie idzie PRZED kadrem, ' +
      'więc w narożach nie zostają puste trójkąty.',
    pola: [
      POLE_ZASOBU,
      { kod: 'x', etykieta: 'Lewy brzeg kadru (px)', rodzaj: 'liczba' },
      { kod: 'y', etykieta: 'Górny brzeg kadru (px)', rodzaj: 'liczba' },
      { kod: 'width', etykieta: 'Szerokość kadru (px)', rodzaj: 'liczba' },
      { kod: 'height', etykieta: 'Wysokość kadru (px)', rodzaj: 'liczba' },
      {
        kod: 'aspectRatio',
        etykieta: 'Proporcje kadru',
        rodzaj: 'wybor',
        pozycje: [
          { wartosc: '', etykieta: 'bez zmiany' },
          { wartosc: '1:1', etykieta: '1:1 — kwadrat' },
          { wartosc: '4:5', etykieta: '4:5 — pion społecznościowy' },
          { wartosc: '16:9', etykieta: '16:9 — panorama' },
          { wartosc: '3:2', etykieta: '3:2 — klatka małoobrazkowa' },
        ],
      },
      { kod: 'straightenDeg', etykieta: 'Prostowanie horyzontu (°)', rodzaj: 'liczba' },
    ],
    zloz(wartosci, otoczenie) {
      return {
        zadanie: zPolami(
          { assetId: wartosci['assetId'] ?? '', windowId: otoczenie.idOkna },
          {
            x: liczba(wartosci, 'x'),
            y: liczba(wartosci, 'y'),
            width: liczba(wartosci, 'width'),
            height: liczba(wartosci, 'height'),
            aspectRatio: tekst(wartosci, 'aspectRatio'),
            straightenDeg: liczba(wartosci, 'straightenDeg'),
          },
        ),
      };
    },
  },
  {
    grupa: 'fotografia',
    komenda: Command.DesignPhotoTransform,
    nazwa: 'Obróć, odbij, wyprostuj perspektywę',
    opis: 'Obrót, odbicia oraz korekcja perspektywy i dystorsji obiektywu.',
    pola: [
      POLE_ZASOBU,
      { kod: 'rotateDeg', etykieta: 'Obrót (°)', rodzaj: 'liczba' },
      { kod: 'flipHorizontal', etykieta: 'Odbij w poziomie', rodzaj: 'logiczne' },
      { kod: 'flipVertical', etykieta: 'Odbij w pionie', rodzaj: 'logiczne' },
      {
        kod: 'lensDistortion',
        etykieta: 'Korekcja obiektywu',
        rodzaj: 'liczba',
        opis: 'Wartość dodatnia prostuje beczkę, ujemna poduszkę.',
      },
      {
        kod: 'perspective',
        etykieta: 'Naroża po korekcji perspektywy',
        rodzaj: 'wielowiersz',
        podpowiedz: '0,0\n1200,40\n1180,800\n20,760',
        opis:
          'Cztery wiersze „x,y" w kolejności: lewy górny, prawy górny, prawy dolny, ' +
          'lewy dolny. Trzy punkty nie wyznaczają czworokąta, więc rdzeń odmówi.',
      },
    ],
    zloz(wartosci, otoczenie) {
      const naroza = (tekst(wartosci, 'perspective') ?? '')
        .split('\n')
        .map((wiersz) =>
          wiersz
            .split(',')
            .map((pozycja) => Number(pozycja.trim()))
            .filter((pozycja) => Number.isFinite(pozycja)),
        )
        .filter((punkt) => punkt.length >= 2)
        .map((punkt) => ({ x: punkt[0] as number, y: punkt[1] as number }));
      return {
        zadanie: zPolami(
          { assetId: wartosci['assetId'] ?? '', windowId: otoczenie.idOkna },
          {
            rotateDeg: liczba(wartosci, 'rotateDeg'),
            flipHorizontal: przelacznik(wartosci, 'flipHorizontal'),
            flipVertical: przelacznik(wartosci, 'flipVertical'),
            lensDistortion: liczba(wartosci, 'lensDistortion'),
            perspective: naroza,
          },
        ),
      };
    },
  },
  {
    grupa: 'fotografia',
    komenda: Command.DesignPhotoResample,
    nazwa: 'Przelicz rozdzielczość',
    opis: 'Lanczos, dwuliniowy albo najbliższy sąsiad, z zachowaniem proporcji.',
    pola: [
      POLE_ZASOBU,
      { kod: 'width', etykieta: 'Szerokość docelowa (px)', rodzaj: 'liczba' },
      { kod: 'height', etykieta: 'Wysokość docelowa (px)', rodzaj: 'liczba' },
      {
        kod: 'filter',
        etykieta: 'Filtr',
        rodzaj: 'wybor',
        pozycje: [
          { wartosc: '', etykieta: 'Lanczos (domyślnie)' },
          { wartosc: 'lanczos', etykieta: 'Lanczos — najostrzejszy' },
          { wartosc: 'bilinear', etykieta: 'Dwuliniowy — łagodny' },
          { wartosc: 'nearest', etykieta: 'Najbliższy sąsiad — twarde piksele' },
        ],
      },
      { kod: 'keepAspectRatio', etykieta: 'Zachowaj proporcje', rodzaj: 'logiczne' },
    ],
    zloz(wartosci, otoczenie) {
      return {
        zadanie: zPolami(
          { assetId: wartosci['assetId'] ?? '', windowId: otoczenie.idOkna },
          {
            width: liczba(wartosci, 'width'),
            height: liczba(wartosci, 'height'),
            filter: tekst(wartosci, 'filter'),
            keepAspectRatio: przelacznik(wartosci, 'keepAspectRatio'),
          },
        ),
      };
    },
  },
  {
    grupa: 'fotografia',
    komenda: Command.DesignPhotoUpscale,
    nazwa: 'Powiększ',
    opis:
      'Powiększenie 2×, 4× albo 8× z wyostrzeniem po powiększeniu. ' +
      'Odpowiedź mówi, którą drogą rdzeń policzył wynik.',
    pola: [
      POLE_ZASOBU,
      {
        kod: 'factor',
        etykieta: 'Krotność',
        rodzaj: 'wybor',
        pozycje: [
          { wartosc: '2', etykieta: '2×' },
          { wartosc: '4', etykieta: '4×' },
          { wartosc: '8', etykieta: '8×' },
        ],
      },
      { kod: 'sharpenAfter', etykieta: 'Wyostrz po powiększeniu', rodzaj: 'logiczne' },
      POLE_KANALU,
    ],
    zloz(wartosci, otoczenie) {
      const krotnosc = liczba(wartosci, 'factor');
      if (krotnosc === undefined) {
        return { odmowa: 'Powiększenie wymaga wskazania krotności: 2, 4 albo 8.' };
      }
      return {
        zadanie: zPolami(
          {
            assetId: wartosci['assetId'] ?? '',
            factor: krotnosc,
            windowId: otoczenie.idOkna,
          },
          {
            sharpenAfter: przelacznik(wartosci, 'sharpenAfter'),
            channelId: tekst(wartosci, 'channelId'),
          },
        ),
      };
    },
  },
  {
    grupa: 'fotografia',
    komenda: Command.DesignPhotoEnhance,
    nazwa: 'Popraw jakość',
    opis:
      'Auto-poziomy, auto-kontrast, odszumienie i wyostrzenie. Odpowiedź niesie ' +
      'wykaz kroków, które NAPRAWDĘ weszły.',
    pola: [
      POLE_ZASOBU,
      { kod: 'autoLevels', etykieta: 'Auto-poziomy', rodzaj: 'logiczne' },
      { kod: 'autoContrast', etykieta: 'Auto-kontrast', rodzaj: 'logiczne' },
      { kod: 'denoise', etykieta: 'Odszumienie (0–1)', rodzaj: 'liczba' },
      { kod: 'sharpen', etykieta: 'Wyostrzenie (0–1)', rodzaj: 'liczba' },
    ],
    zloz(wartosci, otoczenie) {
      return {
        zadanie: zPolami(
          { assetId: wartosci['assetId'] ?? '', windowId: otoczenie.idOkna },
          {
            autoLevels: przelacznik(wartosci, 'autoLevels'),
            autoContrast: przelacznik(wartosci, 'autoContrast'),
            denoise: liczba(wartosci, 'denoise'),
            sharpen: liczba(wartosci, 'sharpen'),
          },
        ),
      };
    },
  },
  {
    grupa: 'fotografia',
    komenda: Command.DesignPhotoColorCorrect,
    nazwa: 'Koryguj barwę',
    opis:
      'Jasność, kontrast, ekspozycja, temperatura, gamma, nasycenie, odcień i krzywe. ' +
      'Odpowiedź niesie ZMIERZONE przesunięcie histogramu.',
    pola: [
      POLE_ZASOBU,
      { kod: 'brightness', etykieta: 'Jasność (−1…1)', rodzaj: 'liczba' },
      { kod: 'contrast', etykieta: 'Kontrast (−1…1)', rodzaj: 'liczba' },
      {
        kod: 'exposure',
        etykieta: 'Ekspozycja (działki)',
        rodzaj: 'liczba',
        opis: 'Jedna działka to podwojenie ilości światła.',
      },
      { kod: 'temperature', etykieta: 'Temperatura barwowa (−1…1)', rodzaj: 'liczba' },
      { kod: 'gamma', etykieta: 'Gamma', rodzaj: 'liczba' },
      { kod: 'saturation', etykieta: 'Nasycenie (−1…1)', rodzaj: 'liczba' },
      { kod: 'hueShiftDeg', etykieta: 'Przesunięcie odcienia (°)', rodzaj: 'liczba' },
      { kod: 'lightness', etykieta: 'Jasność percepcyjna (−1…1)', rodzaj: 'liczba' },
      {
        kod: 'curves',
        etykieta: 'Krzywe tonalne',
        rodzaj: 'wielowiersz',
        podpowiedz: '{"rgb":[{"wejscie":0,"wyjscie":16},{"wejscie":255,"wyjscie":240}]}',
        opis: 'Obiekt JSON kanał → wykaz punktów. Kanały: rgb, r, g, b.',
      },
    ],
    zloz(wartosci, otoczenie) {
      const zapis = tekst(wartosci, 'curves');
      let krzywe: unknown;
      if (zapis !== undefined) {
        try {
          krzywe = JSON.parse(zapis);
        } catch {
          return {
            odmowa:
              'Krzywe tonalne nie są poprawnym zapisem JSON. Rdzeń sprawdza je po swojej ' +
              'stronie, ale zapis nieczytelny nie dojedzie tam w ogóle.',
          };
        }
      }
      return {
        zadanie: zPolami(
          { assetId: wartosci['assetId'] ?? '', windowId: otoczenie.idOkna },
          {
            brightness: liczba(wartosci, 'brightness'),
            contrast: liczba(wartosci, 'contrast'),
            exposure: liczba(wartosci, 'exposure'),
            temperature: liczba(wartosci, 'temperature'),
            gamma: liczba(wartosci, 'gamma'),
            saturation: liczba(wartosci, 'saturation'),
            hueShiftDeg: liczba(wartosci, 'hueShiftDeg'),
            lightness: liczba(wartosci, 'lightness'),
            curves: krzywe,
          },
        ),
      };
    },
  },
  {
    grupa: 'fotografia',
    komenda: Command.DesignPhotoFilterApply,
    nazwa: 'Nałóż filtr',
    opis: 'Rozmycie, wyostrzenie, ziarno, winieta, sepia, monochrom, poświata, cień.',
    pola: [
      POLE_ZASOBU,
      {
        kod: 'filter',
        etykieta: 'Filtr',
        rodzaj: 'wybor',
        pozycje: [
          { wartosc: 'blur', etykieta: 'Rozmycie' },
          { wartosc: 'sharpen', etykieta: 'Wyostrzenie' },
          { wartosc: 'grain', etykieta: 'Ziarno' },
          { wartosc: 'vignette', etykieta: 'Winieta' },
          { wartosc: 'sepia', etykieta: 'Sepia' },
          { wartosc: 'monochrome', etykieta: 'Monochrom' },
          { wartosc: 'glow', etykieta: 'Poświata' },
          { wartosc: 'shadow', etykieta: 'Cień' },
        ],
      },
      { kod: 'amount', etykieta: 'Siła (0–1)', rodzaj: 'liczba' },
    ],
    zloz(wartosci, otoczenie) {
      const filtr = tekst(wartosci, 'filter');
      if (filtr === undefined) {
        return { odmowa: 'Nałożenie filtru wymaga wskazania, który filtr ma wejść.' };
      }
      return {
        zadanie: zPolami(
          { assetId: wartosci['assetId'] ?? '', filter: filtr, windowId: otoczenie.idOkna },
          { amount: liczba(wartosci, 'amount') },
        ),
      };
    },
  },
  {
    grupa: 'fotografia',
    komenda: Command.DesignPhotoRetouch,
    nazwa: 'Retuszuj obszary',
    opis: 'Leczenie z otoczenia albo klonowanie ze wskazanego punktu.',
    pola: [
      POLE_ZASOBU,
      {
        kod: 'regions',
        etykieta: 'Obszary',
        rodzaj: 'wielowiersz',
        podpowiedz: '120,80,40,40\n300,220,25,25',
        opis: 'Jeden obszar na wiersz: „x,y,szerokość,wysokość" w pikselach.',
      },
      {
        kod: 'mode',
        etykieta: 'Sposób',
        rodzaj: 'wybor',
        pozycje: [
          { wartosc: '', etykieta: 'Leczenie z otoczenia (domyślnie)' },
          { wartosc: 'heal', etykieta: 'Leczenie z otoczenia' },
          { wartosc: 'clone', etykieta: 'Klonowanie ze źródła' },
        ],
      },
      { kod: 'sourceX', etykieta: 'Źródło klonowania — X', rodzaj: 'liczba' },
      { kod: 'sourceY', etykieta: 'Źródło klonowania — Y', rodzaj: 'liczba' },
    ],
    zloz(wartosci, otoczenie) {
      const obszary = obszaryZWierszy(wartosci, 'regions');
      if (obszary.length === 0) {
        return {
          odmowa:
            'Retusz wymaga co najmniej jednego obszaru w postaci „x,y,szerokość,wysokość".',
        };
      }
      const zrodloX = liczba(wartosci, 'sourceX');
      const zrodloY = liczba(wartosci, 'sourceY');
      const tryb = tekst(wartosci, 'mode');
      if (tryb === 'clone' && (zrodloX === undefined || zrodloY === undefined)) {
        return { odmowa: 'Klonowanie wymaga wskazania punktu źródłowego — obu współrzędnych.' };
      }
      return {
        zadanie: zPolami(
          {
            assetId: wartosci['assetId'] ?? '',
            regions: obszary,
            windowId: otoczenie.idOkna,
          },
          {
            mode: tryb,
            source:
              zrodloX === undefined || zrodloY === undefined
                ? undefined
                : { x: zrodloX, y: zrodloY },
          },
        ),
      };
    },
  },
  {
    grupa: 'fotografia',
    komenda: Command.DesignPhotoInpaint,
    nazwa: 'Domaluj obszar',
    opis:
      'Wypełnia obszar treścią z otoczenia. Odpowiedź mówi, którą drogą rdzeń policzył wynik.',
    pola: [
      POLE_ZASOBU,
      POLE_MASKI,
      {
        kod: 'regions',
        etykieta: 'Obszary (gdy maski nie ma)',
        rodzaj: 'wielowiersz',
        podpowiedz: '400,120,180,90',
        opis: 'Jeden obszar na wiersz: „x,y,szerokość,wysokość".',
      },
      { kod: 'prompt', etykieta: 'Opis treści', rodzaj: 'tekst' },
      POLE_KANALU,
    ],
    zloz(wartosci, otoczenie) {
      const maska = tekst(wartosci, 'maskAssetId');
      const obszary = obszaryZWierszy(wartosci, 'regions');
      if (maska === undefined && obszary.length === 0) {
        return {
          odmowa:
            'Domalowanie wymaga maski albo obszarów — rdzeń nie zgaduje, który fragment ' +
            'zdjęcia domalować.',
        };
      }
      return {
        zadanie: zPolami(
          { assetId: wartosci['assetId'] ?? '', windowId: otoczenie.idOkna },
          {
            maskAssetId: maska,
            regions: obszary,
            prompt: tekst(wartosci, 'prompt'),
            channelId: tekst(wartosci, 'channelId'),
          },
        ),
      };
    },
  },
  {
    grupa: 'fotografia',
    komenda: Command.DesignPhotoExpand,
    nazwa: 'Rozszerz kadr',
    opis: 'Dokłada obszar poza pierwotną ramkę, wypełniając go z brzegu obrazu.',
    pola: [
      POLE_ZASOBU,
      { kod: 'left', etykieta: 'Z lewej (px)', rodzaj: 'liczba' },
      { kod: 'right', etykieta: 'Z prawej (px)', rodzaj: 'liczba' },
      { kod: 'top', etykieta: 'Od góry (px)', rodzaj: 'liczba' },
      { kod: 'bottom', etykieta: 'Od dołu (px)', rodzaj: 'liczba' },
      { kod: 'prompt', etykieta: 'Opis treści nowego obszaru', rodzaj: 'tekst' },
      POLE_KANALU,
    ],
    zloz(wartosci, otoczenie) {
      return {
        zadanie: zPolami(
          { assetId: wartosci['assetId'] ?? '', windowId: otoczenie.idOkna },
          {
            left: liczba(wartosci, 'left'),
            right: liczba(wartosci, 'right'),
            top: liczba(wartosci, 'top'),
            bottom: liczba(wartosci, 'bottom'),
            prompt: tekst(wartosci, 'prompt'),
            channelId: tekst(wartosci, 'channelId'),
          },
        ),
      };
    },
  },
  {
    grupa: 'fotografia',
    komenda: Command.DesignPhotoBackgroundRemove,
    nazwa: 'Odetnij tło',
    opis:
      'Zostawia kanał krycia. Odpowiedź niesie ZMIERZONY udział punktów przezroczystych — ' +
      'po nim widać, ile obrazu zniknęło.',
    pola: [POLE_ZASOBU, { kod: 'tolerance', etykieta: 'Tolerancja (0–1)', rodzaj: 'liczba' }, POLE_KANALU],
    zloz(wartosci, otoczenie) {
      return {
        zadanie: zPolami(
          { assetId: wartosci['assetId'] ?? '', windowId: otoczenie.idOkna },
          {
            tolerance: liczba(wartosci, 'tolerance'),
            channelId: tekst(wartosci, 'channelId'),
          },
        ),
      };
    },
  },
  {
    grupa: 'fotografia',
    komenda: Command.DesignPhotoSelectObject,
    nazwa: 'Zaznacz obiekt',
    opis: 'Oddaje maskę jako zasób wraz ze zmierzonym udziałem punktów zaznaczonych.',
    pola: [
      POLE_ZASOBU,
      { kod: 'pointX', etykieta: 'Punkt wewnątrz obiektu — X', rodzaj: 'liczba' },
      { kod: 'pointY', etykieta: 'Punkt wewnątrz obiektu — Y', rodzaj: 'liczba' },
      { kod: 'tolerance', etykieta: 'Tolerancja (0–1)', rodzaj: 'liczba' },
    ],
    zloz(wartosci, otoczenie) {
      const x = liczba(wartosci, 'pointX');
      const y = liczba(wartosci, 'pointY');
      if (x === undefined || y === undefined) {
        return { odmowa: 'Zaznaczenie wymaga punktu wewnątrz obiektu — obu współrzędnych.' };
      }
      return {
        zadanie: zPolami(
          { assetId: wartosci['assetId'] ?? '', point: { x, y }, windowId: otoczenie.idOkna },
          { tolerance: liczba(wartosci, 'tolerance') },
        ),
      };
    },
  },
  {
    grupa: 'fotografia',
    komenda: Command.DesignPhotoMaskSet,
    nazwa: 'Nałóż maskę',
    opis: 'Maska nieniszcząca: źródło zostaje, a maska wchodzi do łańcucha edycji.',
    pola: [
      POLE_ZASOBU,
      POLE_MASKI,
      { kod: 'invert', etykieta: 'Odwróć maskę', rodzaj: 'logiczne' },
      { kod: 'featherPx', etykieta: 'Miękkość krawędzi (px)', rodzaj: 'liczba' },
    ],
    zloz(wartosci, otoczenie) {
      const maska = tekst(wartosci, 'maskAssetId');
      if (maska === undefined) {
        return { odmowa: 'Nałożenie maski wymaga wskazania zasobu maski.' };
      }
      return {
        zadanie: zPolami(
          {
            assetId: wartosci['assetId'] ?? '',
            maskAssetId: maska,
            windowId: otoczenie.idOkna,
          },
          {
            invert: przelacznik(wartosci, 'invert'),
            featherPx: liczba(wartosci, 'featherPx'),
          },
        ),
      };
    },
  },
  {
    grupa: 'fotografia',
    komenda: Command.DesignPhotoVectorize,
    nazwa: 'Obrysuj kontury',
    opis: 'Zamienia raster w rysunek wektorowy przez progowanie i obrysowanie obszarów.',
    pola: [
      POLE_ZASOBU,
      { kod: 'colors', etykieta: 'Liczba barw (2–64)', rodzaj: 'liczba' },
      { kod: 'threshold', etykieta: 'Próg (0–1)', rodzaj: 'liczba' },
      { kod: 'smoothing', etykieta: 'Wygładzenie (px)', rodzaj: 'liczba' },
    ],
    zloz(wartosci, otoczenie) {
      return {
        zadanie: zPolami(
          { assetId: wartosci['assetId'] ?? '', windowId: otoczenie.idOkna },
          {
            colors: liczba(wartosci, 'colors'),
            threshold: liczba(wartosci, 'threshold'),
            smoothing: liczba(wartosci, 'smoothing'),
          },
        ),
      };
    },
  },
  {
    grupa: 'fotografia',
    komenda: Command.DesignPhotoBatchApply,
    nazwa: 'Puść wsad',
    opis: 'Ten sam zestaw czynności na wielu zasobach. Zasób nieudany wraca w bilansie.',
    pola: [
      { kod: 'assetIds', etykieta: 'Materiały', rodzaj: 'zasoby' },
      { kod: 'presetId', etykieta: 'Nastawa', rodzaj: 'tekst' },
      {
        kod: 'operations',
        etykieta: 'Czynności (gdy nastawy nie ma)',
        rodzaj: 'wielowiersz',
        podpowiedz: PODPOWIEDZ_CZYNNOSCI,
        opis: 'Wykaz JSON czynności w kolejności wykonania.',
      },
    ],
    zloz(wartosci, otoczenie) {
      const materialy = wykaz(wartosci, 'assetIds');
      if (materialy.length === 0) {
        return { odmowa: 'Wsad wymaga wskazania co najmniej jednego materiału.' };
      }
      const zapis = tekst(wartosci, 'operations');
      let czynnosci: unknown;
      if (zapis !== undefined) {
        try {
          czynnosci = JSON.parse(zapis);
        } catch {
          return { odmowa: 'Wykaz czynności nie jest poprawnym zapisem JSON.' };
        }
      }
      if (czynnosci === undefined && tekst(wartosci, 'presetId') === undefined) {
        return {
          odmowa:
            'Wsad wymaga czynności albo nastawy — rdzeń nie zgaduje, co ma z tymi zasobami zrobić.',
        };
      }
      return {
        zadanie: zPolami(
          { assetIds: materialy, windowId: otoczenie.idOkna },
          { presetId: tekst(wartosci, 'presetId'), operations: czynnosci },
        ),
      };
    },
  },
  {
    grupa: 'fotografia',
    komenda: Command.DesignPhotoLayerComposite,
    nazwa: 'Złóż warstwy',
    opis: 'Warstwy rastrowe z trybami mieszania i kryciem. Warstwa pominięta wraca w bilansie.',
    pola: [
      {
        kod: 'layers',
        etykieta: 'Warstwy',
        rodzaj: 'wielowiersz',
        podpowiedz: '[{"assetId":"zasob-1"},{"assetId":"zasob-2","opacity":0.5,"blendMode":"multiply"}]',
        opis: 'Wykaz JSON warstw w kolejności od spodu.',
      },
      { kod: 'width', etykieta: 'Szerokość płótna (px)', rodzaj: 'liczba' },
      { kod: 'height', etykieta: 'Wysokość płótna (px)', rodzaj: 'liczba' },
    ],
    zloz(wartosci, otoczenie) {
      const zapis = tekst(wartosci, 'layers');
      if (zapis === undefined) {
        return { odmowa: 'Kompozycja warstw wymaga wykazu warstw.' };
      }
      let warstwy: unknown;
      try {
        warstwy = JSON.parse(zapis);
      } catch {
        return { odmowa: 'Wykaz warstw nie jest poprawnym zapisem JSON.' };
      }
      return {
        zadanie: zPolami(
          { layers: warstwy, windowId: otoczenie.idOkna },
          { width: liczba(wartosci, 'width'), height: liczba(wartosci, 'height') },
        ),
      };
    },
  },
  {
    grupa: 'fotografia',
    komenda: Command.DesignPhotoPresetSave,
    nazwa: 'Zapisz nastawę',
    opis: 'Utrwala zestaw czynności pod nazwą — do powtórzenia wsadem.',
    pola: [
      { kod: 'name', etykieta: 'Nazwa nastawy', rodzaj: 'tekst' },
      {
        kod: 'operations',
        etykieta: 'Czynności',
        rodzaj: 'wielowiersz',
        podpowiedz: PODPOWIEDZ_CZYNNOSCI,
      },
      { kod: 'presetId', etykieta: 'Nastawa nadpisywana', rodzaj: 'tekst' },
    ],
    zloz(wartosci, otoczenie) {
      const nazwa = tekst(wartosci, 'name');
      const zapis = tekst(wartosci, 'operations');
      if (nazwa === undefined || zapis === undefined) {
        return { odmowa: 'Nastawa wymaga nazwy i wykazu czynności.' };
      }
      let czynnosci: unknown;
      try {
        czynnosci = JSON.parse(zapis);
      } catch {
        return { odmowa: 'Wykaz czynności nie jest poprawnym zapisem JSON.' };
      }
      return {
        zadanie: zPolami(
          { windowId: otoczenie.idOkna, name: nazwa, operations: czynnosci },
          { presetId: tekst(wartosci, 'presetId') },
        ),
      };
    },
  },
  {
    grupa: 'fotografia',
    komenda: Command.DesignPhotoPresetList,
    nazwa: 'Wykaz nastaw',
    opis: 'Nastawy zapisane w tym oknie.',
    pola: [],
    zloz(_wartosci, otoczenie) {
      return { zadanie: { windowId: otoczenie.idOkna } };
    },
  },
  {
    grupa: 'fotografia',
    komenda: Command.DesignPhotoMetadataGet,
    nazwa: 'Zmierz właściwości',
    opis: 'Wymiary, format, model barwny, kanał krycia, rozdzielczość i dane EXIF z pliku.',
    pola: [POLE_ZASOBU],
    zloz(wartosci) {
      return { zadanie: { assetId: wartosci['assetId'] ?? '' } };
    },
  },
  {
    grupa: 'fotografia',
    komenda: Command.DesignPhotoHistoryGet,
    nazwa: 'Łańcuch edycji',
    opis: 'Czynności, którymi zasób powstał, wraz z nastawami i drogą rachunku.',
    pola: [POLE_ZASOBU, { kod: 'limit', etykieta: 'Granica liczby ogniw', rodzaj: 'liczba' }],
    zloz(wartosci) {
      return {
        zadanie: zPolami(
          { assetId: wartosci['assetId'] ?? '' },
          { limit: liczba(wartosci, 'limit') },
        ),
      };
    },
  },

  // ── Warsztat wektora ───────────────────────────────────────────────────────
  {
    grupa: 'wektor',
    komenda: Command.DesignVectorPathSet,
    nazwa: 'Narysuj ścieżkę',
    opis: 'Zapisuje węzły ścieżki. Wskazanie ścieżki zmienia ją; brak zakłada nową.',
    pola: [
      { kod: 'boardId', etykieta: 'Kompozycja', rodzaj: 'tekst' },
      {
        kod: 'nodes',
        etykieta: 'Węzły',
        rodzaj: 'wielowiersz',
        podpowiedz: '[{"x":0,"y":0,"kind":"corner"},{"x":100,"y":0,"kind":"corner"}]',
        opis: 'Wykaz JSON węzłów. Uchwyty (handleIn/handleOut) są ODSUNIĘCIAMI od węzła.',
      },
      { kod: 'pathId', etykieta: 'Ścieżka zmieniana', rodzaj: 'tekst' },
      { kod: 'layerId', etykieta: 'Warstwa', rodzaj: 'tekst' },
      { kod: 'closed', etykieta: 'Zamknięta', rodzaj: 'logiczne' },
      { kod: 'name', etykieta: 'Nazwa ścieżki', rodzaj: 'tekst' },
    ],
    zloz(wartosci) {
      const kompozycja = tekst(wartosci, 'boardId');
      const zapis = tekst(wartosci, 'nodes');
      if (kompozycja === undefined || zapis === undefined) {
        return { odmowa: 'Ścieżka wymaga wskazania kompozycji i wykazu węzłów.' };
      }
      let wezly: unknown;
      try {
        wezly = JSON.parse(zapis);
      } catch {
        return { odmowa: 'Wykaz węzłów nie jest poprawnym zapisem JSON.' };
      }
      return {
        zadanie: zPolami(
          { boardId: kompozycja, nodes: wezly },
          {
            pathId: tekst(wartosci, 'pathId'),
            layerId: tekst(wartosci, 'layerId'),
            closed: przelacznik(wartosci, 'closed'),
            name: tekst(wartosci, 'name'),
          },
        ),
      };
    },
  },
  {
    grupa: 'wektor',
    komenda: Command.DesignVectorPathList,
    nazwa: 'Wykaz ścieżek',
    opis: 'Ścieżki kompozycji wraz z węzłami.',
    pola: [
      { kod: 'boardId', etykieta: 'Kompozycja', rodzaj: 'tekst' },
      { kod: 'layerId', etykieta: 'Zawężenie do warstwy', rodzaj: 'tekst' },
    ],
    zloz(wartosci) {
      const kompozycja = tekst(wartosci, 'boardId');
      if (kompozycja === undefined) {
        return { odmowa: 'Wykaz ścieżek wymaga wskazania kompozycji.' };
      }
      return {
        zadanie: zPolami({ boardId: kompozycja }, { layerId: tekst(wartosci, 'layerId') }),
      };
    },
  },
  {
    grupa: 'wektor',
    komenda: Command.DesignVectorPathRemove,
    nazwa: 'Usuń ścieżkę',
    opis: 'Usuwa ścieżkę. Ścieżki, której nie ma, rdzeń nie odmawia — odpowiada „nie było".',
    pola: [{ kod: 'pathId', etykieta: 'Ścieżka', rodzaj: 'tekst' }],
    zloz(wartosci) {
      const sciezka = tekst(wartosci, 'pathId');
      if (sciezka === undefined) return { odmowa: 'Usunięcie wymaga wskazania ścieżki.' };
      return { zadanie: { pathId: sciezka } };
    },
  },
  {
    grupa: 'wektor',
    komenda: Command.DesignVectorShapeAdd,
    nazwa: 'Dołóż kształt',
    opis:
      'Prostokąt, elipsa, wielokąt, gwiazda albo odcinek. Kształt powstaje OD RAZU jako ' +
      'węzły ścieżki, więc da się go ciągnąć piórem od pierwszej chwili.',
    pola: [
      { kod: 'boardId', etykieta: 'Kompozycja', rodzaj: 'tekst' },
      {
        kod: 'kind',
        etykieta: 'Kształt',
        rodzaj: 'wybor',
        pozycje: [
          { wartosc: 'rectangle', etykieta: 'Prostokąt' },
          { wartosc: 'ellipse', etykieta: 'Elipsa' },
          { wartosc: 'polygon', etykieta: 'Wielokąt' },
          { wartosc: 'star', etykieta: 'Gwiazda' },
          { wartosc: 'line', etykieta: 'Odcinek' },
        ],
      },
      { kod: 'x', etykieta: 'X', rodzaj: 'liczba' },
      { kod: 'y', etykieta: 'Y', rodzaj: 'liczba' },
      { kod: 'width', etykieta: 'Szerokość', rodzaj: 'liczba' },
      { kod: 'height', etykieta: 'Wysokość', rodzaj: 'liczba' },
      { kod: 'cornerRadius', etykieta: 'Zaokrąglenie naroży', rodzaj: 'liczba' },
      { kod: 'points', etykieta: 'Wierzchołki / ramiona', rodzaj: 'liczba' },
      { kod: 'innerRadius', etykieta: 'Promień wewnętrzny gwiazdy', rodzaj: 'liczba' },
      { kod: 'layerId', etykieta: 'Warstwa', rodzaj: 'tekst' },
    ],
    zloz(wartosci) {
      const kompozycja = tekst(wartosci, 'boardId');
      const rodzaj = tekst(wartosci, 'kind');
      const x = liczba(wartosci, 'x');
      const y = liczba(wartosci, 'y');
      const szerokosc = liczba(wartosci, 'width');
      const wysokosc = liczba(wartosci, 'height');
      if (
        kompozycja === undefined ||
        rodzaj === undefined ||
        x === undefined ||
        y === undefined ||
        szerokosc === undefined ||
        wysokosc === undefined
      ) {
        return {
          odmowa: 'Kształt wymaga kompozycji, rodzaju, położenia i obu wymiarów.',
        };
      }
      return {
        zadanie: zPolami(
          { boardId: kompozycja, kind: rodzaj, x, y, width: szerokosc, height: wysokosc },
          {
            cornerRadius: liczba(wartosci, 'cornerRadius'),
            points: liczba(wartosci, 'points'),
            innerRadius: liczba(wartosci, 'innerRadius'),
            layerId: tekst(wartosci, 'layerId'),
          },
        ),
      };
    },
  },
  {
    grupa: 'wektor',
    komenda: Command.DesignVectorBoolean,
    nazwa: 'Operacja logiczna',
    opis:
      'Suma, różnica, część wspólna, wykluczenie. Ścieżki źródłowe usunięte wracają ' +
      'w wykazie — pole puste tam, gdzie zniknęły, byłoby ciszą.',
    pola: [
      {
        kod: 'pathIds',
        etykieta: 'Ścieżki',
        rodzaj: 'tekst',
        podpowiedz: 'sciezka-1, sciezka-2',
        opis: 'Co najmniej dwie. Kolejność rozstrzyga przy różnicy.',
      },
      {
        kod: 'operation',
        etykieta: 'Operacja',
        rodzaj: 'wybor',
        pozycje: [
          { wartosc: 'union', etykieta: 'Suma' },
          { wartosc: 'subtract', etykieta: 'Różnica' },
          { wartosc: 'intersect', etykieta: 'Część wspólna' },
          { wartosc: 'exclude', etykieta: 'Wykluczenie' },
        ],
      },
      { kod: 'keepSources', etykieta: 'Zostaw ścieżki źródłowe', rodzaj: 'logiczne' },
    ],
    zloz(wartosci) {
      const sciezki = wykaz(wartosci, 'pathIds');
      const operacja = tekst(wartosci, 'operation');
      if (sciezki.length < 2 || operacja === undefined) {
        return { odmowa: 'Operacja logiczna wymaga dwóch ścieżek i wskazania operacji.' };
      }
      return {
        zadanie: zPolami(
          { pathIds: sciezki, operation: operacja },
          { keepSources: przelacznik(wartosci, 'keepSources') },
        ),
      };
    },
  },
  {
    grupa: 'wektor',
    komenda: Command.DesignVectorTextPath,
    nazwa: 'Tekst na ścieżce',
    opis:
      'Układa tekst wzdłuż ścieżki albo zamienia go w kontury. Tekst źródłowy ZOSTAJE ' +
      'w nazwie ścieżki, bo kontury nie wiedzą, że były literami.',
    pola: [
      { kod: 'boardId', etykieta: 'Kompozycja', rodzaj: 'tekst' },
      { kod: 'text', etykieta: 'Treść', rodzaj: 'tekst' },
      { kod: 'fontFamily', etykieta: 'Krój pisma', rodzaj: 'tekst', podpowiedz: 'Go Regular' },
      { kod: 'fontSize', etykieta: 'Rozmiar pisma', rodzaj: 'liczba' },
      { kod: 'pathId', etykieta: 'Ścieżka nośna', rodzaj: 'tekst' },
      { kod: 'outline', etykieta: 'Zamień w kontury', rodzaj: 'logiczne' },
      { kod: 'x', etykieta: 'X (bez ścieżki)', rodzaj: 'liczba' },
      { kod: 'y', etykieta: 'Y (bez ścieżki)', rodzaj: 'liczba' },
    ],
    zloz(wartosci) {
      const kompozycja = tekst(wartosci, 'boardId');
      const tresc = tekst(wartosci, 'text');
      const kroj = tekst(wartosci, 'fontFamily');
      const rozmiar = liczba(wartosci, 'fontSize');
      if (
        kompozycja === undefined ||
        tresc === undefined ||
        kroj === undefined ||
        rozmiar === undefined
      ) {
        return { odmowa: 'Tekst na ścieżce wymaga kompozycji, treści, kroju i rozmiaru pisma.' };
      }
      return {
        zadanie: zPolami(
          { boardId: kompozycja, text: tresc, fontFamily: kroj, fontSize: rozmiar },
          {
            pathId: tekst(wartosci, 'pathId'),
            outline: przelacznik(wartosci, 'outline'),
            x: liczba(wartosci, 'x'),
            y: liczba(wartosci, 'y'),
          },
        ),
      };
    },
  },
  {
    grupa: 'wektor',
    komenda: Command.DesignVectorOptimize,
    nazwa: 'Oczyść zapis',
    opis: 'Skraca zapis współrzędnych. Ubytek bajtów jest ZMIERZONY, nie oszacowany.',
    pola: [
      { kod: 'pathIds', etykieta: 'Ścieżki', rodzaj: 'tekst', podpowiedz: 'sciezka-1, sciezka-2' },
      { kod: 'boardId', etykieta: 'Kompozycja (gdy ścieżek nie wskazano)', rodzaj: 'tekst' },
      { kod: 'precision', etykieta: 'Miejsca po przecinku', rodzaj: 'liczba' },
    ],
    zloz(wartosci) {
      const sciezki = wykaz(wartosci, 'pathIds');
      const kompozycja = tekst(wartosci, 'boardId');
      if (sciezki.length === 0 && kompozycja === undefined) {
        return {
          odmowa:
            'Czyszczenie wymaga wskazania ścieżek albo kompozycji — rdzeń nie przepisuje ' +
            'wszystkiego, co ma w bazie.',
        };
      }
      return {
        zadanie: zPolami(
          {},
          {
            pathIds: sciezki,
            boardId: kompozycja,
            precision: liczba(wartosci, 'precision'),
          },
        ),
      };
    },
  },
  {
    grupa: 'wektor',
    komenda: Command.DesignVectorSymbolSet,
    nazwa: 'Zapisz symbol',
    opis: 'Definicja do wielokrotnego użycia ze ścieżek i warstw.',
    pola: [
      { kod: 'boardId', etykieta: 'Kompozycja', rodzaj: 'tekst' },
      { kod: 'name', etykieta: 'Nazwa symbolu', rodzaj: 'tekst' },
      { kod: 'symbolId', etykieta: 'Symbol zmieniany', rodzaj: 'tekst' },
      { kod: 'pathIds', etykieta: 'Ścieżki', rodzaj: 'tekst' },
      { kod: 'layerIds', etykieta: 'Warstwy', rodzaj: 'tekst' },
    ],
    zloz(wartosci) {
      const kompozycja = tekst(wartosci, 'boardId');
      const nazwa = tekst(wartosci, 'name');
      if (kompozycja === undefined || nazwa === undefined) {
        return { odmowa: 'Symbol wymaga wskazania kompozycji i nazwy.' };
      }
      return {
        zadanie: zPolami(
          { boardId: kompozycja, name: nazwa },
          {
            symbolId: tekst(wartosci, 'symbolId'),
            pathIds: wykaz(wartosci, 'pathIds'),
            layerIds: wykaz(wartosci, 'layerIds'),
          },
        ),
      };
    },
  },
  {
    grupa: 'wektor',
    komenda: Command.DesignVectorSymbolList,
    nazwa: 'Wykaz symboli',
    opis: 'Symbole kompozycji wraz z liczbą wystąpień.',
    pola: [{ kod: 'boardId', etykieta: 'Kompozycja', rodzaj: 'tekst' }],
    zloz(wartosci) {
      const kompozycja = tekst(wartosci, 'boardId');
      if (kompozycja === undefined) return { odmowa: 'Wykaz symboli wymaga kompozycji.' };
      return { zadanie: { boardId: kompozycja } };
    },
  },
  {
    grupa: 'wektor',
    komenda: Command.DesignVectorExport,
    nazwa: 'Wydaj wektor',
    opis: 'SVG, PDF albo EPS. Wydanie wektorowe ZOSTAJE wektorem — nic tu nie rasteryzuje.',
    pola: [
      { kod: 'boardId', etykieta: 'Kompozycja', rodzaj: 'tekst' },
      {
        kod: 'target',
        etykieta: 'Postać',
        rodzaj: 'wybor',
        pozycje: [
          { wartosc: 'svg', etykieta: 'SVG' },
          { wartosc: 'pdf', etykieta: 'PDF' },
          { wartosc: 'eps', etykieta: 'EPS' },
        ],
      },
      { kod: 'pathIds', etykieta: 'Ścieżki (brak bierze całość)', rodzaj: 'tekst' },
      { kod: 'frameId', etykieta: 'Ramka zawężająca', rodzaj: 'tekst' },
    ],
    zloz(wartosci) {
      const kompozycja = tekst(wartosci, 'boardId');
      const postac = tekst(wartosci, 'target');
      if (kompozycja === undefined || postac === undefined) {
        return { odmowa: 'Wydanie wektora wymaga kompozycji i postaci wydania.' };
      }
      return {
        zadanie: zPolami(
          { boardId: kompozycja, target: postac },
          { pathIds: wykaz(wartosci, 'pathIds'), frameId: tekst(wartosci, 'frameId') },
        ),
      };
    },
  },

  // ── Warsztat druku ─────────────────────────────────────────────────────────
  {
    grupa: 'druk',
    komenda: Command.DesignPrintProfileSet,
    nazwa: 'Zapisz profil wydania',
    opis: 'Przestrzeń barw, norma, spad, znaczniki, rozdzielczość, nośnik.',
    pola: [
      {
        kod: 'colorSpace',
        etykieta: 'Przestrzeń barw',
        rodzaj: 'wybor',
        pozycje: [
          { wartosc: 'cmyk', etykieta: 'CMYK' },
          { wartosc: 'rgb', etykieta: 'RGB' },
          { wartosc: 'grayscale', etykieta: 'Skala szarości' },
          { wartosc: 'spot', etykieta: 'Barwy dodatkowe' },
        ],
      },
      { kod: 'name', etykieta: 'Nazwa profilu', rodzaj: 'tekst' },
      {
        kod: 'standard',
        etykieta: 'Norma',
        rodzaj: 'wybor',
        pozycje: [
          { wartosc: '', etykieta: 'bez normy' },
          { wartosc: 'pdfX1a', etykieta: 'PDF/X-1a' },
          { wartosc: 'pdfX3', etykieta: 'PDF/X-3' },
          { wartosc: 'pdfX4', etykieta: 'PDF/X-4' },
          { wartosc: 'pdfA2b', etykieta: 'PDF/A-2b' },
        ],
      },
      { kod: 'bleedMm', etykieta: 'Spad (mm)', rodzaj: 'liczba' },
      { kod: 'cropMarks', etykieta: 'Znaczniki cięcia', rodzaj: 'logiczne' },
      { kod: 'registrationMarks', etykieta: 'Znaczniki pasowania', rodzaj: 'logiczne' },
      { kod: 'colorBar', etykieta: 'Pasek kontrolny barw', rodzaj: 'logiczne' },
      { kod: 'dpi', etykieta: 'Rozdzielczość (dpi)', rodzaj: 'liczba' },
      { kod: 'iccProfile', etykieta: 'Profil ICC', rodzaj: 'tekst' },
      { kod: 'overprintBlack', etykieta: 'Nadruk czerni', rodzaj: 'logiczne' },
      { kod: 'paperSize', etykieta: 'Nośnik', rodzaj: 'tekst', podpowiedz: 'A4' },
      { kod: 'profileId', etykieta: 'Profil nadpisywany', rodzaj: 'tekst' },
    ],
    zloz(wartosci, otoczenie) {
      const przestrzen = tekst(wartosci, 'colorSpace');
      if (przestrzen === undefined) {
        return { odmowa: 'Profil wydania wymaga wskazania przestrzeni barw.' };
      }
      const profil = zPolami(
        { colorSpace: przestrzen },
        {
          name: tekst(wartosci, 'name'),
          standard: tekst(wartosci, 'standard'),
          bleedMm: liczba(wartosci, 'bleedMm'),
          cropMarks: przelacznik(wartosci, 'cropMarks'),
          registrationMarks: przelacznik(wartosci, 'registrationMarks'),
          colorBar: przelacznik(wartosci, 'colorBar'),
          dpi: liczba(wartosci, 'dpi'),
          iccProfile: tekst(wartosci, 'iccProfile'),
          overprintBlack: przelacznik(wartosci, 'overprintBlack'),
          paperSize: tekst(wartosci, 'paperSize'),
        },
      );
      return {
        zadanie: zPolami(
          { windowId: otoczenie.idOkna, profile: profil },
          { profileId: tekst(wartosci, 'profileId') },
        ),
      };
    },
  },
  {
    grupa: 'druk',
    komenda: Command.DesignPrintProfileList,
    nazwa: 'Wykaz profili',
    opis: 'Profile okna, nośniki znane rdzeniowi i profile ICC OBECNE na serwerze.',
    pola: [],
    zloz(_wartosci, otoczenie) {
      return { zadanie: { windowId: otoczenie.idOkna } };
    },
  },
  {
    grupa: 'druk',
    komenda: Command.DesignPrintPaperList,
    nazwa: 'Wykaz nośników',
    opis: 'Rozmiary nośników znane rdzeniowi, z zawężeniem rodziną.',
    pola: [
      {
        kod: 'family',
        etykieta: 'Rodzina',
        rodzaj: 'tekst',
        podpowiedz: 'iso-a',
      },
    ],
    zloz(wartosci) {
      return { zadanie: zPolami({}, { family: tekst(wartosci, 'family') }) };
    },
  },
  {
    grupa: 'druk',
    komenda: Command.DesignPrintPreflight,
    nazwa: 'Kontrola przeddrukowa',
    opis:
      'Mierzy materiał wobec profilu i oddaje zastrzeżenia w kolejności WAGI. ' +
      'Wada o wadze błędu zablokuje wydanie.',
    pola: [
      { kod: 'boardId', etykieta: 'Kompozycja', rodzaj: 'tekst' },
      { kod: 'assetId', etykieta: 'Zasób', rodzaj: 'zasob' },
      { kod: 'profileId', etykieta: 'Profil', rodzaj: 'tekst' },
    ],
    zloz(wartosci) {
      const kompozycja = tekst(wartosci, 'boardId');
      const zasob = tekst(wartosci, 'assetId');
      if (kompozycja === undefined && zasob === undefined) {
        return {
          odmowa: 'Kontrola mierzy MATERIAŁ — wskaż kompozycję albo zasób.',
        };
      }
      return {
        zadanie: zPolami(
          {},
          { boardId: kompozycja, assetId: zasob, profileId: tekst(wartosci, 'profileId') },
        ),
      };
    },
  },
  {
    grupa: 'druk',
    komenda: Command.DesignPrintExport,
    nazwa: 'Wydaj do druku',
    opis:
      'PDF, TIFF albo EPS. Wada o wadze błędu ODMAWIA wydania; pominięcie kontroli ' +
      'jest jawnym wyborem i wraca w odpowiedzi.',
    pola: [
      {
        kod: 'format',
        etykieta: 'Format',
        rodzaj: 'wybor',
        pozycje: [
          { wartosc: 'pdf', etykieta: 'PDF' },
          { wartosc: 'tiff', etykieta: 'TIFF' },
          { wartosc: 'eps', etykieta: 'EPS' },
        ],
      },
      { kod: 'boardId', etykieta: 'Kompozycja', rodzaj: 'tekst' },
      { kod: 'assetId', etykieta: 'Zasób', rodzaj: 'zasob' },
      { kod: 'frameId', etykieta: 'Ramka zawężająca', rodzaj: 'tekst' },
      { kod: 'profileId', etykieta: 'Profil', rodzaj: 'tekst' },
      {
        kod: 'skipPreflight',
        etykieta: 'Pomiń kontrolę przeddrukową',
        rodzaj: 'logiczne',
        opis: 'Pominięcie wraca w odpowiedzi — nikt potem nie powie, że nie wiedział.',
      },
    ],
    zloz(wartosci, otoczenie) {
      const format = tekst(wartosci, 'format');
      if (format === undefined) return { odmowa: 'Wydanie wymaga wskazania formatu.' };
      return {
        zadanie: zPolami(
          { format, windowId: otoczenie.idOkna },
          {
            boardId: tekst(wartosci, 'boardId'),
            assetId: tekst(wartosci, 'assetId'),
            frameId: tekst(wartosci, 'frameId'),
            profileId: tekst(wartosci, 'profileId'),
            skipPreflight: przelacznik(wartosci, 'skipPreflight'),
          },
        ),
      };
    },
  },
  {
    grupa: 'druk',
    komenda: Command.DesignLargeformatTile,
    nazwa: 'Podziel na kafle',
    opis:
      'Kafle wielkiego formatu z zakładką na sklejenie. Odpowiedź niesie ZMIERZONĄ ' +
      'rozdzielczość skuteczną w rozmiarze docelowym.',
    pola: [
      POLE_ZASOBU,
      { kod: 'tileWidthMm', etykieta: 'Szerokość kafla (mm)', rodzaj: 'liczba' },
      { kod: 'tileHeightMm', etykieta: 'Wysokość kafla (mm)', rodzaj: 'liczba' },
      { kod: 'overlapMm', etykieta: 'Zakładka (mm)', rodzaj: 'liczba' },
      { kod: 'targetWidthMm', etykieta: 'Szerokość całości (mm)', rodzaj: 'liczba' },
      { kod: 'targetHeightMm', etykieta: 'Wysokość całości (mm)', rodzaj: 'liczba' },
      { kod: 'markers', etykieta: 'Znaczniki sklejenia', rodzaj: 'logiczne' },
    ],
    zloz(wartosci, otoczenie) {
      const szerokosc = liczba(wartosci, 'tileWidthMm');
      const wysokosc = liczba(wartosci, 'tileHeightMm');
      if (szerokosc === undefined || wysokosc === undefined) {
        return { odmowa: 'Podział wymaga obu wymiarów kafla w milimetrach.' };
      }
      return {
        zadanie: zPolami(
          {
            assetId: wartosci['assetId'] ?? '',
            tileWidthMm: szerokosc,
            tileHeightMm: wysokosc,
            windowId: otoczenie.idOkna,
          },
          {
            overlapMm: liczba(wartosci, 'overlapMm'),
            targetWidthMm: liczba(wartosci, 'targetWidthMm'),
            targetHeightMm: liczba(wartosci, 'targetHeightMm'),
            markers: przelacznik(wartosci, 'markers'),
          },
        ),
      };
    },
  },
  {
    grupa: 'druk',
    komenda: Command.DesignChartRender,
    nazwa: 'Wyrysuj wykres',
    opis:
      'Wykres z SERII DANYCH, nie z obrazu — po zmianie liczb da się go przerysować.',
    pola: [
      {
        kod: 'kind',
        etykieta: 'Rodzaj',
        rodzaj: 'wybor',
        pozycje: [
          { wartosc: 'column', etykieta: 'Kolumnowy' },
          { wartosc: 'bar', etykieta: 'Słupkowy poziomy' },
          { wartosc: 'line', etykieta: 'Liniowy' },
          { wartosc: 'area', etykieta: 'Warstwowy' },
          { wartosc: 'pie', etykieta: 'Kołowy' },
          { wartosc: 'donut', etykieta: 'Pierścieniowy' },
          { wartosc: 'scatter', etykieta: 'Punktowy' },
          { wartosc: 'stackedBar', etykieta: 'Skumulowany' },
        ],
      },
      {
        kod: 'series',
        etykieta: 'Serie danych',
        rodzaj: 'wielowiersz',
        podpowiedz: '[{"name":"Przychód","values":[12,18,24]}]',
      },
      { kod: 'categories', etykieta: 'Podpisy kategorii', rodzaj: 'tekst', podpowiedz: 'I, II, III' },
      { kod: 'title', etykieta: 'Tytuł', rodzaj: 'tekst' },
      { kod: 'width', etykieta: 'Szerokość', rodzaj: 'liczba' },
      { kod: 'height', etykieta: 'Wysokość', rodzaj: 'liczba' },
      {
        kod: 'format',
        etykieta: 'Format',
        rodzaj: 'wybor',
        pozycje: [
          { wartosc: '', etykieta: 'SVG (domyślnie)' },
          { wartosc: 'svg', etykieta: 'SVG' },
          { wartosc: 'png', etykieta: 'PNG' },
        ],
      },
      { kod: 'tokenSetId', etykieta: 'Zestaw żetonów', rodzaj: 'tekst' },
      { kod: 'legend', etykieta: 'Legenda', rodzaj: 'logiczne' },
    ],
    zloz(wartosci, otoczenie) {
      const rodzaj = tekst(wartosci, 'kind');
      const zapis = tekst(wartosci, 'series');
      if (rodzaj === undefined || zapis === undefined) {
        return { odmowa: 'Wykres wymaga rodzaju i serii danych.' };
      }
      let serie: unknown;
      try {
        serie = JSON.parse(zapis);
      } catch {
        return { odmowa: 'Serie danych nie są poprawnym zapisem JSON.' };
      }
      return {
        zadanie: zPolami(
          { kind: rodzaj, series: serie, windowId: otoczenie.idOkna },
          {
            categories: wykaz(wartosci, 'categories'),
            title: tekst(wartosci, 'title'),
            width: liczba(wartosci, 'width'),
            height: liczba(wartosci, 'height'),
            format: tekst(wartosci, 'format'),
            tokenSetId: tekst(wartosci, 'tokenSetId'),
            legend: przelacznik(wartosci, 'legend'),
          },
        ),
      };
    },
  },
  {
    grupa: 'druk',
    komenda: Command.DesignDiagramRender,
    nazwa: 'Wyrysuj schemat',
    opis: 'Węzły i połączenia. Węzeł, którego układ nie umieścił, wraca w bilansie.',
    pola: [
      {
        kod: 'kind',
        etykieta: 'Rodzaj',
        rodzaj: 'wybor',
        pozycje: [
          { wartosc: 'flow', etykieta: 'Przepływ' },
          { wartosc: 'org', etykieta: 'Organizacyjny' },
          { wartosc: 'sequence', etykieta: 'Sekwencja' },
          { wartosc: 'mindmap', etykieta: 'Mapa myśli' },
          { wartosc: 'timeline', etykieta: 'Oś czasu' },
        ],
      },
      {
        kod: 'nodes',
        etykieta: 'Węzły',
        rodzaj: 'wielowiersz',
        podpowiedz: '[{"id":"a","label":"Start"},{"id":"b","label":"Koniec","parentId":"a"}]',
      },
      {
        kod: 'edges',
        etykieta: 'Połączenia',
        rodzaj: 'wielowiersz',
        podpowiedz: '[{"fromId":"a","toId":"b"}]',
      },
      { kod: 'title', etykieta: 'Tytuł', rodzaj: 'tekst' },
      {
        kod: 'direction',
        etykieta: 'Kierunek',
        rodzaj: 'wybor',
        pozycje: [
          { wartosc: '', etykieta: 'pionowo (domyślnie)' },
          { wartosc: 'vertical', etykieta: 'Pionowo' },
          { wartosc: 'horizontal', etykieta: 'Poziomo' },
        ],
      },
      {
        kod: 'format',
        etykieta: 'Format',
        rodzaj: 'wybor',
        pozycje: [
          { wartosc: '', etykieta: 'SVG (domyślnie)' },
          { wartosc: 'svg', etykieta: 'SVG' },
          { wartosc: 'png', etykieta: 'PNG' },
        ],
      },
      { kod: 'tokenSetId', etykieta: 'Zestaw żetonów', rodzaj: 'tekst' },
    ],
    zloz(wartosci, otoczenie) {
      const rodzaj = tekst(wartosci, 'kind');
      const zapis = tekst(wartosci, 'nodes');
      if (rodzaj === undefined || zapis === undefined) {
        return { odmowa: 'Schemat wymaga rodzaju i wykazu węzłów.' };
      }
      let wezly: unknown;
      let polaczenia: unknown;
      try {
        wezly = JSON.parse(zapis);
      } catch {
        return { odmowa: 'Wykaz węzłów nie jest poprawnym zapisem JSON.' };
      }
      const zapisPolaczen = tekst(wartosci, 'edges');
      if (zapisPolaczen !== undefined) {
        try {
          polaczenia = JSON.parse(zapisPolaczen);
        } catch {
          return { odmowa: 'Wykaz połączeń nie jest poprawnym zapisem JSON.' };
        }
      }
      return {
        zadanie: zPolami(
          { kind: rodzaj, nodes: wezly, windowId: otoczenie.idOkna },
          {
            edges: polaczenia,
            title: tekst(wartosci, 'title'),
            direction: tekst(wartosci, 'direction'),
            format: tekst(wartosci, 'format'),
            tokenSetId: tekst(wartosci, 'tokenSetId'),
          },
        ),
      };
    },
  },

  // ── Przeglądarka baz zdjęciowych ───────────────────────────────────────────
  {
    grupa: 'bazy',
    komenda: Command.DesignStockSearch,
    nazwa: 'Szukaj w bazach',
    opis:
      'Pyta dostawców i mówi, KOGO zapytano i KTO nie odpowiedział. Dostawcy bez klucza ' +
      'pracują od pierwszej minuty; dostawcy z kluczem czytają go z sejfu rdzenia.',
    pola: [
      { kod: 'query', etykieta: 'Fraza', rodzaj: 'tekst' },
      {
        kod: 'provider',
        etykieta: 'Dostawca',
        rodzaj: 'wybor',
        pozycje: [
          { wartosc: '', etykieta: 'wszyscy' },
          { wartosc: 'openverse', etykieta: 'Openverse (bez klucza)' },
          { wartosc: 'wikimedia', etykieta: 'Wikimedia Commons (bez klucza)' },
          { wartosc: 'met', etykieta: 'Met Museum (bez klucza)' },
          { wartosc: 'nasa', etykieta: 'NASA (bez klucza)' },
          { wartosc: 'unsplash', etykieta: 'Unsplash (klucz z sejfu)' },
          { wartosc: 'pexels', etykieta: 'Pexels (klucz z sejfu)' },
          { wartosc: 'pixabay', etykieta: 'Pixabay (klucz z sejfu)' },
          { wartosc: 'smithsonian', etykieta: 'Smithsonian (klucz z sejfu)' },
        ],
      },
      { kod: 'limit', etykieta: 'Granica na dostawcę', rodzaj: 'liczba' },
    ],
    zloz(wartosci, otoczenie) {
      const fraza = tekst(wartosci, 'query');
      if (fraza === undefined) return { odmowa: 'Wyszukanie wymaga frazy.' };
      return {
        zadanie: zPolami(
          { query: fraza, windowId: otoczenie.idOkna },
          { provider: tekst(wartosci, 'provider'), limit: liczba(wartosci, 'limit') },
        ),
      };
    },
  },
  {
    grupa: 'bazy',
    komenda: Command.DesignStockImport,
    nazwa: 'Wciągnij zasób',
    opis:
      'Wciąga bajty i ZAPISUJE LICENCJĘ razem z zasobem. Zasób z bazy zewnętrznej bez ' +
      'zapisanej licencji jest usterką, nie zasobem — dlatego brak licencji odmawia.',
    pola: [
      {
        kod: 'provider',
        etykieta: 'Dostawca',
        rodzaj: 'tekst',
        opis: 'Nazwa dostawcy z wyniku wyszukania.',
      },
      { kod: 'externalId', etykieta: 'Identyfikator u dostawcy', rodzaj: 'tekst' },
    ],
    zloz(wartosci, otoczenie) {
      const dostawca = tekst(wartosci, 'provider');
      const identyfikator = tekst(wartosci, 'externalId');
      if (dostawca === undefined || identyfikator === undefined) {
        return { odmowa: 'Wciągnięcie wymaga dostawcy i identyfikatora zasobu u niego.' };
      }
      return {
        zadanie: {
          provider: dostawca,
          externalId: identyfikator,
          windowId: otoczenie.idOkna,
        },
      };
    },
  },

  // ── Warsztat publikacji ────────────────────────────────────────────────────
  {
    grupa: 'publikacja',
    komenda: Command.DesignTemplateSave,
    nazwa: 'Zapisz szablon albo publikację',
    opis:
      'Szablon materiału. Wykaz STRON czyni z niego publikację wielostronicową; brak ' +
      'stron zostawia szablon jednostronicowy.',
    pola: [
      { kod: 'name', etykieta: 'Nazwa', rodzaj: 'tekst' },
      {
        kod: 'kind',
        etykieta: 'Rodzaj materiału',
        rodzaj: 'wybor',
        pozycje: [
          { wartosc: 'print', etykieta: 'Druk' },
          { wartosc: 'social', etykieta: 'Społecznościowy' },
          { wartosc: 'banner', etykieta: 'Baner' },
          { wartosc: 'presentation', etykieta: 'Prezentacja' },
          { wartosc: 'email', etykieta: 'Poczta' },
        ],
      },
      { kod: 'width', etykieta: 'Szerokość', rodzaj: 'liczba' },
      { kod: 'height', etykieta: 'Wysokość', rodzaj: 'liczba' },
      { kod: 'description', etykieta: 'Do czego służy', rodzaj: 'tekst' },
      { kod: 'templateId', etykieta: 'Szablon nadpisywany', rodzaj: 'tekst' },
      {
        kod: 'pages',
        etykieta: 'Strony publikacji',
        rodzaj: 'wielowiersz',
        podpowiedz: '[{"number":1,"name":"okładka","layers":[]},{"number":2,"layers":[]}]',
        opis: 'Wykaz JSON stron. Numery liczy się od jednego i nie mogą się powtarzać.',
      },
    ],
    zloz(wartosci, otoczenie) {
      const nazwa = tekst(wartosci, 'name');
      const rodzaj = tekst(wartosci, 'kind');
      const szerokosc = liczba(wartosci, 'width');
      const wysokosc = liczba(wartosci, 'height');
      if (
        nazwa === undefined ||
        rodzaj === undefined ||
        szerokosc === undefined ||
        wysokosc === undefined
      ) {
        return { odmowa: 'Szablon wymaga nazwy, rodzaju i obu wymiarów materiału.' };
      }
      const zapis = tekst(wartosci, 'pages');
      let strony: unknown;
      if (zapis !== undefined) {
        try {
          strony = JSON.parse(zapis);
        } catch {
          return { odmowa: 'Wykaz stron nie jest poprawnym zapisem JSON.' };
        }
      }
      return {
        zadanie: zPolami(
          {
            windowId: otoczenie.idOkna,
            name: nazwa,
            kind: rodzaj,
            width: szerokosc,
            height: wysokosc,
          },
          {
            description: tekst(wartosci, 'description'),
            templateId: tekst(wartosci, 'templateId'),
            pages: strony,
          },
        ),
      };
    },
  },
  {
    grupa: 'publikacja',
    komenda: Command.DesignTemplateList,
    nazwa: 'Wykaz szablonów',
    opis: 'Szablony okna, z zawężeniem rodzaju.',
    pola: [
      {
        kod: 'kind',
        etykieta: 'Rodzaj',
        rodzaj: 'wybor',
        pozycje: [
          { wartosc: '', etykieta: 'wszystkie' },
          { wartosc: 'print', etykieta: 'Druk' },
          { wartosc: 'social', etykieta: 'Społecznościowy' },
          { wartosc: 'banner', etykieta: 'Baner' },
          { wartosc: 'presentation', etykieta: 'Prezentacja' },
          { wartosc: 'email', etykieta: 'Poczta' },
        ],
      },
    ],
    zloz(wartosci, otoczenie) {
      return {
        zadanie: zPolami({ windowId: otoczenie.idOkna }, { kind: tekst(wartosci, 'kind') }),
      };
    },
  },
  {
    grupa: 'publikacja',
    komenda: Command.DesignTemplateApply,
    nazwa: 'Zastosuj szablon',
    opis:
      'Zakłada kompozycję ze szablonu wraz z podstawieniami treści. Nazwa, której szablon ' +
      'nie ma, wraca w wykazie niedopasowanych — cicha zgoda dałaby połowę kampanii.',
    pola: [
      { kod: 'templateId', etykieta: 'Szablon', rodzaj: 'tekst' },
      { kod: 'boardId', etykieta: 'Kompozycja docelowa', rodzaj: 'tekst' },
      { kod: 'pageNumber', etykieta: 'Strona publikacji', rodzaj: 'liczba' },
      {
        kod: 'replacements',
        etykieta: 'Podstawienia treści',
        rodzaj: 'wielowiersz',
        podpowiedz: '{"logo":"zasob-1","zdjęcie produktu":"zasob-2"}',
        opis: 'Obiekt JSON: nazwa warstwy → identyfikator zasobu.',
      },
    ],
    zloz(wartosci, otoczenie) {
      const szablon = tekst(wartosci, 'templateId');
      if (szablon === undefined) return { odmowa: 'Zastosowanie wymaga wskazania szablonu.' };
      const zapis = tekst(wartosci, 'replacements');
      let podstawienia: unknown;
      if (zapis !== undefined) {
        try {
          podstawienia = JSON.parse(zapis);
        } catch {
          return { odmowa: 'Podstawienia nie są poprawnym zapisem JSON.' };
        }
      }
      return {
        zadanie: zPolami(
          { templateId: szablon, windowId: otoczenie.idOkna },
          {
            boardId: tekst(wartosci, 'boardId'),
            pageNumber: liczba(wartosci, 'pageNumber'),
            replacements: podstawienia,
          },
        ),
      };
    },
  },
  {
    grupa: 'publikacja',
    komenda: Command.DesignPrintExport,
    nazwa: 'Wydaj publikację',
    opis:
      'Strony szablonu wychodzą JEDNYM plikiem PDF, przez tę samą kontrolę przeddrukową. ' +
      'Odpowiedź niesie ZMIERZONĄ liczbę stron w pliku.',
    pola: [
      { kod: 'templateId', etykieta: 'Szablon publikacji', rodzaj: 'tekst' },
      {
        kod: 'format',
        etykieta: 'Format',
        rodzaj: 'wybor',
        pozycje: [
          { wartosc: 'pdf', etykieta: 'PDF — jedyny wielostronicowy' },
          { wartosc: 'tiff', etykieta: 'TIFF (jedna strona)' },
          { wartosc: 'eps', etykieta: 'EPS (jedna strona)' },
        ],
      },
      {
        kod: 'pageOrder',
        etykieta: 'Kolejność stron',
        rodzaj: 'tekst',
        podpowiedz: '1, 4, 2, 3',
        opis: 'Numery stron w kolejności wydania. Brak bierze kolejność numerów.',
      },
      {
        kod: 'binding',
        etykieta: 'Oprawa',
        rodzaj: 'wybor',
        pozycje: [
          { wartosc: '', etykieta: 'bez oprawy' },
          { wartosc: 'brak', etykieta: 'Bez oprawy' },
          { wartosc: 'zeszytowa', etykieta: 'Zeszytowa (stron ÷ 4)' },
          { wartosc: 'klejona', etykieta: 'Klejona' },
          { wartosc: 'spiralna', etykieta: 'Spiralna' },
        ],
      },
      { kod: 'profileId', etykieta: 'Profil wydania', rodzaj: 'tekst' },
      { kod: 'skipPreflight', etykieta: 'Pomiń kontrolę przeddrukową', rodzaj: 'logiczne' },
    ],
    zloz(wartosci, otoczenie) {
      const szablon = tekst(wartosci, 'templateId');
      const format = tekst(wartosci, 'format');
      if (szablon === undefined || format === undefined) {
        return { odmowa: 'Wydanie publikacji wymaga szablonu i formatu.' };
      }
      return {
        zadanie: zPolami(
          { templateId: szablon, format, windowId: otoczenie.idOkna },
          {
            pageOrder: wykazLiczb(wartosci, 'pageOrder'),
            binding: tekst(wartosci, 'binding'),
            profileId: tekst(wartosci, 'profileId'),
            skipPreflight: przelacznik(wartosci, 'skipPreflight'),
          },
        ),
      };
    },
  },
  {
    grupa: 'publikacja',
    komenda: Command.DesignCampaignSetBuild,
    nazwa: 'Zbuduj komplet kampanii',
    opis: 'Kompozycja w każdym zamówionym rozmiarze. Rozmiar nieudany wraca w bilansie.',
    pola: [
      { kod: 'boardId', etykieta: 'Kompozycja źródłowa', rodzaj: 'tekst' },
      {
        kod: 'sizes',
        etykieta: 'Rozmiary',
        rodzaj: 'wielowiersz',
        podpowiedz: '[{"name":"post","width":1080,"height":1080},{"name":"baner","width":1200,"height":628}]',
      },
      {
        kod: 'format',
        etykieta: 'Format',
        rodzaj: 'wybor',
        pozycje: [
          { wartosc: '', etykieta: 'PNG (domyślnie)' },
          { wartosc: 'png', etykieta: 'PNG' },
          { wartosc: 'jpeg', etykieta: 'JPEG' },
          { wartosc: 'pdf', etykieta: 'PDF' },
        ],
      },
    ],
    zloz(wartosci, otoczenie) {
      const kompozycja = tekst(wartosci, 'boardId');
      const zapis = tekst(wartosci, 'sizes');
      if (kompozycja === undefined || zapis === undefined) {
        return { odmowa: 'Komplet kampanii wymaga kompozycji i wykazu rozmiarów.' };
      }
      let rozmiary: unknown;
      try {
        rozmiary = JSON.parse(zapis);
      } catch {
        return { odmowa: 'Wykaz rozmiarów nie jest poprawnym zapisem JSON.' };
      }
      return {
        zadanie: zPolami(
          { boardId: kompozycja, sizes: rozmiary, windowId: otoczenie.idOkna },
          { format: tekst(wartosci, 'format') },
        ),
      };
    },
  },
  // ── Barwa (w warsztacie fotografii) ────────────────────────────────────────
  {
    grupa: 'fotografia',
    komenda: Command.DesignColorPaletteGenerate,
    nazwa: 'Zbuduj paletę z harmonii',
    opis: 'Dopełnienie, triada, tetrada, analogia albo monochrom wokół barwy wiodącej.',
    pola: [
      { kod: 'baseColor', etykieta: 'Barwa wiodąca', rodzaj: 'tekst', podpowiedz: '#1f6feb' },
      {
        kod: 'harmony',
        etykieta: 'Harmonia',
        rodzaj: 'wybor',
        pozycje: [
          { wartosc: 'complementary', etykieta: 'Dopełnienie' },
          { wartosc: 'analogous', etykieta: 'Analogia' },
          { wartosc: 'triad', etykieta: 'Triada' },
          { wartosc: 'tetrad', etykieta: 'Tetrada' },
          { wartosc: 'splitComplementary', etykieta: 'Dopełnienie rozdzielone' },
          { wartosc: 'mono', etykieta: 'Monochrom' },
        ],
      },
      { kod: 'count', etykieta: 'Liczba barw', rodzaj: 'liczba' },
    ],
    zloz(wartosci, otoczenie) {
      const wiodaca = tekst(wartosci, 'baseColor');
      const harmonia = tekst(wartosci, 'harmony');
      if (wiodaca === undefined || harmonia === undefined) {
        return { odmowa: 'Paleta wymaga barwy wiodącej i reguły harmonii.' };
      }
      return {
        zadanie: zPolami(
          { baseColor: wiodaca, harmony: harmonia },
          { count: liczba(wartosci, 'count'), windowId: otoczenie.idOkna },
        ),
      };
    },
  },
  {
    grupa: 'fotografia',
    komenda: Command.DesignColorPaletteExtract,
    nazwa: 'Wyciągnij paletę z obrazu',
    opis: 'Barwy dominujące wraz ze ZMIERZONYM udziałem każdej w obrazie.',
    pola: [POLE_ZASOBU, { kod: 'count', etykieta: 'Liczba barw', rodzaj: 'liczba' }],
    zloz(wartosci) {
      return {
        zadanie: zPolami(
          { assetId: wartosci['assetId'] ?? '' },
          { count: liczba(wartosci, 'count') },
        ),
      };
    },
  },
  {
    grupa: 'fotografia',
    komenda: Command.DesignColorContrastCheck,
    nazwa: 'Zmierz kontrast',
    opis: 'Współczynnik kontrastu wedle WCAG 2.1 wraz z oceną progów AA i AAA.',
    pola: [
      { kod: 'foreground', etykieta: 'Barwa pierwszoplanowa', rodzaj: 'tekst' },
      { kod: 'background', etykieta: 'Barwa tła', rodzaj: 'tekst' },
      { kod: 'fontSize', etykieta: 'Rozmiar pisma', rodzaj: 'liczba' },
      { kod: 'bold', etykieta: 'Pismo pogrubione', rodzaj: 'logiczne' },
    ],
    zloz(wartosci) {
      const pierwszy = tekst(wartosci, 'foreground');
      const tlo = tekst(wartosci, 'background');
      if (pierwszy === undefined || tlo === undefined) {
        return { odmowa: 'Pomiar kontrastu wymaga obu barw.' };
      }
      return {
        zadanie: zPolami(
          { foreground: pierwszy, background: tlo },
          { fontSize: liczba(wartosci, 'fontSize'), bold: przelacznik(wartosci, 'bold') },
        ),
      };
    },
  },
  {
    grupa: 'fotografia',
    komenda: Command.DesignColorVisionSimulate,
    nazwa: 'Symuluj wadę widzenia',
    opis: 'Obraz widziany przy protanopii, deuteranopii, tritanopii albo achromatopsji.',
    pola: [
      POLE_ZASOBU,
      {
        kod: 'vision',
        etykieta: 'Wada widzenia',
        rodzaj: 'wybor',
        pozycje: [
          { wartosc: 'protanopia', etykieta: 'Protanopia' },
          { wartosc: 'deuteranopia', etykieta: 'Deuteranopia' },
          { wartosc: 'tritanopia', etykieta: 'Tritanopia' },
          { wartosc: 'achromatopsia', etykieta: 'Achromatopsja' },
        ],
      },
    ],
    zloz(wartosci, otoczenie) {
      const wada = tekst(wartosci, 'vision');
      if (wada === undefined) return { odmowa: 'Symulacja wymaga wskazania wady widzenia.' };
      return {
        zadanie: { assetId: wartosci['assetId'] ?? '', vision: wada, windowId: otoczenie.idOkna },
      };
    },
  },
  {
    grupa: 'fotografia',
    komenda: Command.DesignColorConvert,
    nazwa: 'Przelicz barwę',
    opis: 'Ten sam kolor we wszystkich przestrzeniach naraz: hex, RGB, HSL, Lab, CMYK.',
    pola: [
      { kod: 'value', etykieta: 'Zapis barwy', rodzaj: 'tekst', podpowiedz: '#1f6feb' },
      {
        kod: 'space',
        etykieta: 'Przestrzeń wejściowa',
        rodzaj: 'wybor',
        pozycje: [
          { wartosc: '', etykieta: 'rozpoznaj sam' },
          { wartosc: 'hex', etykieta: 'Szesnastkowo' },
          { wartosc: 'rgb', etykieta: 'RGB' },
          { wartosc: 'hsl', etykieta: 'HSL' },
          { wartosc: 'lab', etykieta: 'Lab' },
          { wartosc: 'cmyk', etykieta: 'CMYK' },
        ],
      },
    ],
    zloz(wartosci) {
      const zapis = tekst(wartosci, 'value');
      if (zapis === undefined) return { odmowa: 'Przeliczenie wymaga zapisu barwy.' };
      return { zadanie: zPolami({ value: zapis }, { space: tekst(wartosci, 'space') }) };
    },
  },
  {
    grupa: 'fotografia',
    komenda: Command.DesignColorAccessibilityAudit,
    nazwa: 'Zbadaj dostępność zestawu',
    opis: 'Mierzy kontrast wszystkich PAR barwnych żetonów i oddaje pary łamiące próg.',
    pola: [
      { kod: 'tokenSetId', etykieta: 'Zestaw żetonów', rodzaj: 'tekst' },
      { kod: 'minimumRatio', etykieta: 'Próg kontrastu', rodzaj: 'liczba' },
    ],
    zloz(wartosci) {
      const zestaw = tekst(wartosci, 'tokenSetId');
      if (zestaw === undefined) return { odmowa: 'Badanie wymaga wskazania zestawu żetonów.' };
      return {
        zadanie: zPolami(
          { tokenSetId: zestaw },
          { minimumRatio: liczba(wartosci, 'minimumRatio') },
        ),
      };
    },
  },
  {
    grupa: 'wektor',
    komenda: Command.DesignColorGradientSet,
    nazwa: 'Ustaw gradient',
    opis:
      'Gradient rozstrzyga się CELEM: powtórne wywołanie na tej samej ścieżce ZMIENIA ' +
      'gradient, zamiast dokładać drugi obok.',
    pola: [
      { kod: 'boardId', etykieta: 'Kompozycja', rodzaj: 'tekst' },
      {
        kod: 'gradient',
        etykieta: 'Gradient',
        rodzaj: 'wielowiersz',
        podpowiedz: '{"kind":"linear","angle":90,"stops":[{"offset":0,"color":"#000"},{"offset":1,"color":"#fff"}]}',
      },
      { kod: 'pathId', etykieta: 'Ścieżka', rodzaj: 'tekst' },
      { kod: 'layerId', etykieta: 'Warstwa', rodzaj: 'tekst' },
    ],
    zloz(wartosci) {
      const kompozycja = tekst(wartosci, 'boardId');
      const zapis = tekst(wartosci, 'gradient');
      if (kompozycja === undefined || zapis === undefined) {
        return { odmowa: 'Gradient wymaga kompozycji i nastaw gradientu.' };
      }
      let gradient: unknown;
      try {
        gradient = JSON.parse(zapis);
      } catch {
        return { odmowa: 'Nastawy gradientu nie są poprawnym zapisem JSON.' };
      }
      return {
        zadanie: zPolami(
          { boardId: kompozycja, gradient },
          { pathId: tekst(wartosci, 'pathId'), layerId: tekst(wartosci, 'layerId') },
        ),
      };
    },
  },

  // ── Ikony i kroje (w warsztacie wektora) ───────────────────────────────────
  {
    grupa: 'wektor',
    komenda: Command.DesignIconLibrarySearch,
    nazwa: 'Szukaj ikony',
    opis: 'Katalog ikon WKOMPILOWANY w binarium rdzenia — jest zawsze i jest jeden.',
    pola: [
      { kod: 'query', etykieta: 'Fraza', rodzaj: 'tekst' },
      { kod: 'set', etykieta: 'Zestaw', rodzaj: 'tekst' },
      { kod: 'limit', etykieta: 'Granica', rodzaj: 'liczba' },
    ],
    zloz(wartosci) {
      return {
        zadanie: zPolami(
          {},
          {
            query: tekst(wartosci, 'query'),
            set: tekst(wartosci, 'set'),
            limit: liczba(wartosci, 'limit'),
          },
        ),
      };
    },
  },
  {
    grupa: 'wektor',
    komenda: Command.DesignIconSet,
    nazwa: 'Zapisz ikonę własną',
    opis:
      'Ikona własna okna. Odpowiedź niesie miejsca, w których ikona NIE trzyma siatki — ' +
      'bilans zamiast ciszy.',
    pola: [
      { kod: 'name', etykieta: 'Nazwa', rodzaj: 'tekst' },
      { kod: 'svg', etykieta: 'Treść SVG', rodzaj: 'wielowiersz' },
      { kod: 'iconId', etykieta: 'Ikona zmieniana', rodzaj: 'tekst' },
      { kod: 'gridSize', etykieta: 'Siatka', rodzaj: 'liczba' },
      { kod: 'strokeWidth', etykieta: 'Grubość obrysu', rodzaj: 'liczba' },
      { kod: 'tags', etykieta: 'Etykiety', rodzaj: 'tekst' },
    ],
    zloz(wartosci, otoczenie) {
      const nazwa = tekst(wartosci, 'name');
      const svg = tekst(wartosci, 'svg');
      if (nazwa === undefined || svg === undefined) {
        return { odmowa: 'Ikona wymaga nazwy i treści SVG.' };
      }
      return {
        zadanie: zPolami(
          { windowId: otoczenie.idOkna, name: nazwa, svg },
          {
            iconId: tekst(wartosci, 'iconId'),
            gridSize: liczba(wartosci, 'gridSize'),
            strokeWidth: liczba(wartosci, 'strokeWidth'),
            tags: wykaz(wartosci, 'tags'),
          },
        ),
      };
    },
  },
  {
    grupa: 'wektor',
    komenda: Command.DesignIconGenerate,
    nazwa: 'Zbuduj ikony dla pojęć',
    opis: 'Pojęcie, dla którego ikona NIE powstała, wraca w bilansie — nie w ciszy.',
    pola: [
      { kod: 'concepts', etykieta: 'Pojęcia', rodzaj: 'tekst', podpowiedz: 'dom, szukaj, kosz' },
      { kod: 'gridSize', etykieta: 'Siatka', rodzaj: 'liczba' },
      { kod: 'strokeWidth', etykieta: 'Grubość obrysu', rodzaj: 'liczba' },
      { kod: 'styleReferenceIconId', etykieta: 'Ikona wzorcowa', rodzaj: 'tekst' },
      POLE_KANALU,
    ],
    zloz(wartosci, otoczenie) {
      const pojecia = wykaz(wartosci, 'concepts');
      if (pojecia.length === 0) {
        return { odmowa: 'Budowa ikon wymaga wskazania pojęć — rdzeń ich nie wymyśla.' };
      }
      return {
        zadanie: zPolami(
          { windowId: otoczenie.idOkna, concepts: pojecia },
          {
            gridSize: liczba(wartosci, 'gridSize'),
            strokeWidth: liczba(wartosci, 'strokeWidth'),
            styleReferenceIconId: tekst(wartosci, 'styleReferenceIconId'),
            channelId: tekst(wartosci, 'channelId'),
          },
        ),
      };
    },
  },
  {
    grupa: 'wektor',
    komenda: Command.DesignIconSpriteBuild,
    nazwa: 'Zbuduj pakiet ikon',
    opis: 'Pakiet SVG z symbolami albo KRÓJ IKONOWY złożony przez rdzeń z konturów.',
    pola: [
      { kod: 'iconIds', etykieta: 'Ikony', rodzaj: 'tekst', podpowiedz: 'rdzen-24/dom, ikona-1' },
      {
        kod: 'kind',
        etykieta: 'Postać',
        rodzaj: 'wybor',
        pozycje: [
          { wartosc: 'spriteSvg', etykieta: 'Pakiet SVG' },
          { wartosc: 'webfont', etykieta: 'Krój ikonowy' },
        ],
      },
      { kod: 'name', etykieta: 'Nazwa pakietu', rodzaj: 'tekst' },
    ],
    zloz(wartosci, otoczenie) {
      const ikony = wykaz(wartosci, 'iconIds');
      const postac = tekst(wartosci, 'kind');
      if (ikony.length === 0 || postac === undefined) {
        return { odmowa: 'Pakiet wymaga ikon i wskazania postaci.' };
      }
      return {
        zadanie: zPolami(
          { iconIds: ikony, kind: postac },
          { name: tekst(wartosci, 'name'), windowId: otoczenie.idOkna },
        ),
      };
    },
  },
  {
    grupa: 'wektor',
    komenda: Command.DesignFaviconBuild,
    nazwa: 'Zbuduj ikonę witryny',
    opis: 'Komplet rozmiarów wraz z manifestem. Odpowiedź niesie rozmiary, które POWSTAŁY.',
    pola: [
      POLE_ZASOBU,
      { kod: 'sizes', etykieta: 'Rozmiary', rodzaj: 'tekst', podpowiedz: '16, 32, 180' },
      { kod: 'includeManifest', etykieta: 'Dołóż manifest', rodzaj: 'logiczne' },
    ],
    zloz(wartosci, otoczenie) {
      return {
        zadanie: zPolami(
          { assetId: wartosci['assetId'] ?? '', windowId: otoczenie.idOkna },
          {
            sizes: wykazLiczb(wartosci, 'sizes'),
            includeManifest: przelacznik(wartosci, 'includeManifest'),
          },
        ),
      };
    },
  },
  {
    grupa: 'wektor',
    komenda: Command.DesignFontPairSuggest,
    nazwa: 'Zestawienia krojów',
    opis:
      'Propozycje par krojów. Pole „droga" mówi, czy powstały kanałem modelu, czy regułą ' +
      'rdzenia — a reguła zestawia wyłącznie kroje, którymi rdzeń NAPRAWDĘ dysponuje.',
    pola: [
      { kod: 'mood', etykieta: 'Charakter', rodzaj: 'tekst' },
      { kod: 'baseFont', etykieta: 'Krój, do którego szukamy pary', rodzaj: 'tekst' },
      { kod: 'count', etykieta: 'Liczba propozycji', rodzaj: 'liczba' },
      POLE_KANALU,
    ],
    zloz(wartosci, otoczenie) {
      return {
        zadanie: zPolami(
          { windowId: otoczenie.idOkna },
          {
            mood: tekst(wartosci, 'mood'),
            baseFont: tekst(wartosci, 'baseFont'),
            count: liczba(wartosci, 'count'),
            channelId: tekst(wartosci, 'channelId'),
          },
        ),
      };
    },
  },
  {
    grupa: 'wektor',
    komenda: Command.DesignFontPreview,
    nazwa: 'Podgląd kroju',
    opis:
      'Podgląd z PRAWDZIWYCH konturów glifów. Pole „dostępny" mówi prawdę: kroju, ' +
      'którego rdzeń nie ma, nie podstawia innym.',
    pola: [
      { kod: 'fontFamily', etykieta: 'Krój', rodzaj: 'tekst', podpowiedz: 'Go Regular' },
      { kod: 'sampleText', etykieta: 'Tekst próbny', rodzaj: 'tekst' },
      { kod: 'sizes', etykieta: 'Rozmiary', rodzaj: 'tekst', podpowiedz: '12, 16, 24' },
      { kod: 'subsetText', etykieta: 'Podzbiór glifów', rodzaj: 'tekst' },
    ],
    zloz(wartosci) {
      const kroj = tekst(wartosci, 'fontFamily');
      if (kroj === undefined) return { odmowa: 'Podgląd wymaga wskazania kroju.' };
      return {
        zadanie: zPolami(
          { fontFamily: kroj },
          {
            sampleText: tekst(wartosci, 'sampleText'),
            sizes: wykazLiczb(wartosci, 'sizes'),
            subsetText: tekst(wartosci, 'subsetText'),
          },
        ),
      };
    },
  },
  {
    grupa: 'wektor',
    komenda: Command.DesignFontGlyphsGet,
    nazwa: 'Glify kroju',
    opis: 'Znaki kroju wraz z konturem każdego. Znak, którego krój nie ma, nie wchodzi.',
    pola: [
      { kod: 'fontFamily', etykieta: 'Krój', rodzaj: 'tekst' },
      { kod: 'from', etykieta: 'Od punktu kodowego', rodzaj: 'liczba' },
      { kod: 'to', etykieta: 'Do punktu kodowego', rodzaj: 'liczba' },
      { kod: 'limit', etykieta: 'Granica', rodzaj: 'liczba' },
    ],
    zloz(wartosci) {
      const kroj = tekst(wartosci, 'fontFamily');
      if (kroj === undefined) return { odmowa: 'Wykaz glifów wymaga wskazania kroju.' };
      return {
        zadanie: zPolami(
          { fontFamily: kroj },
          {
            from: liczba(wartosci, 'from'),
            to: liczba(wartosci, 'to'),
            limit: liczba(wartosci, 'limit'),
          },
        ),
      };
    },
  },

  // ── Warsztat makiety ───────────────────────────────────────────────────────
  {
    grupa: 'makieta',
    komenda: Command.DesignFrameSet,
    nazwa: 'Ustaw ramkę',
    opis: 'Ramka jest EKRANEM: wyznacza obszar wydania i to jej rozmiar się sprawdza.',
    pola: [
      { kod: 'boardId', etykieta: 'Kompozycja', rodzaj: 'tekst' },
      { kod: 'name', etykieta: 'Nazwa ramki', rodzaj: 'tekst' },
      { kod: 'frameId', etykieta: 'Ramka zmieniana', rodzaj: 'tekst' },
      { kod: 'width', etykieta: 'Szerokość', rodzaj: 'liczba' },
      { kod: 'height', etykieta: 'Wysokość', rodzaj: 'liczba' },
      { kod: 'x', etykieta: 'X na kanwie', rodzaj: 'liczba' },
      { kod: 'y', etykieta: 'Y na kanwie', rodzaj: 'liczba' },
      {
        kod: 'devicePreset',
        etykieta: 'Nastawa urządzenia',
        rodzaj: 'tekst',
        podpowiedz: 'telefon',
        opis: 'Wypełnia wymiary; wskazanie wymiarów wprost ją bije.',
      },
    ],
    zloz(wartosci) {
      const kompozycja = tekst(wartosci, 'boardId');
      const nazwa = tekst(wartosci, 'name');
      if (kompozycja === undefined || nazwa === undefined) {
        return { odmowa: 'Ramka wymaga kompozycji i nazwy.' };
      }
      return {
        zadanie: zPolami(
          { boardId: kompozycja, name: nazwa },
          {
            frameId: tekst(wartosci, 'frameId'),
            width: liczba(wartosci, 'width'),
            height: liczba(wartosci, 'height'),
            x: liczba(wartosci, 'x'),
            y: liczba(wartosci, 'y'),
            devicePreset: tekst(wartosci, 'devicePreset'),
          },
        ),
      };
    },
  },
  {
    grupa: 'makieta',
    komenda: Command.DesignFrameList,
    nazwa: 'Wykaz ramek',
    opis: 'Ramki kompozycji wraz z nastawami urządzeń znanymi rdzeniowi.',
    pola: [{ kod: 'boardId', etykieta: 'Kompozycja', rodzaj: 'tekst' }],
    zloz(wartosci) {
      const kompozycja = tekst(wartosci, 'boardId');
      if (kompozycja === undefined) return { odmowa: 'Wykaz ramek wymaga kompozycji.' };
      return { zadanie: { boardId: kompozycja } };
    },
  },
  {
    grupa: 'makieta',
    komenda: Command.DesignFrameRemove,
    nazwa: 'Usuń ramkę',
    opis:
      'Usuwa ramkę i ZWALNIA jej warstwy — warstwy zostają w kompozycji, a odpowiedź ' +
      'mówi które. Cisza o tym kazałaby oknu zgadywać, czy ich szukać dalej.',
    pola: [{ kod: 'frameId', etykieta: 'Ramka', rodzaj: 'tekst' }],
    zloz(wartosci) {
      const ramka = tekst(wartosci, 'frameId');
      if (ramka === undefined) return { odmowa: 'Usunięcie wymaga wskazania ramki.' };
      return { zadanie: { frameId: ramka } };
    },
  },
  {
    grupa: 'makieta',
    komenda: Command.DesignLayoutAuto,
    nazwa: 'Ułóż automatycznie',
    opis: 'Przelicza położenia warstw ramki i ZAPISUJE je — to zmiana, nie podgląd.',
    pola: [
      { kod: 'frameId', etykieta: 'Ramka', rodzaj: 'tekst' },
      {
        kod: 'direction',
        etykieta: 'Kierunek',
        rodzaj: 'wybor',
        pozycje: [
          { wartosc: 'vertical', etykieta: 'Pionowo' },
          { wartosc: 'horizontal', etykieta: 'Poziomo' },
        ],
      },
      { kod: 'gap', etykieta: 'Odstęp', rodzaj: 'liczba' },
      { kod: 'paddingTop', etykieta: 'Odstęp od góry', rodzaj: 'liczba' },
      { kod: 'paddingRight', etykieta: 'Odstęp od prawej', rodzaj: 'liczba' },
      { kod: 'paddingBottom', etykieta: 'Odstęp od dołu', rodzaj: 'liczba' },
      { kod: 'paddingLeft', etykieta: 'Odstęp od lewej', rodzaj: 'liczba' },
      {
        kod: 'align',
        etykieta: 'Wyrównanie w poprzek',
        rodzaj: 'wybor',
        pozycje: [
          { wartosc: '', etykieta: 'do początku' },
          { wartosc: 'start', etykieta: 'Do początku' },
          { wartosc: 'center', etykieta: 'Do środka' },
          { wartosc: 'end', etykieta: 'Do końca' },
          { wartosc: 'stretch', etykieta: 'Rozciągnij' },
        ],
      },
      { kod: 'layerIds', etykieta: 'Warstwy w kolejności', rodzaj: 'tekst' },
    ],
    zloz(wartosci) {
      const ramka = tekst(wartosci, 'frameId');
      const kierunek = tekst(wartosci, 'direction');
      if (ramka === undefined || kierunek === undefined) {
        return { odmowa: 'Układ automatyczny wymaga ramki i kierunku układania.' };
      }
      const uklad = zPolami(
        { direction: kierunek },
        {
          gap: liczba(wartosci, 'gap'),
          paddingTop: liczba(wartosci, 'paddingTop'),
          paddingRight: liczba(wartosci, 'paddingRight'),
          paddingBottom: liczba(wartosci, 'paddingBottom'),
          paddingLeft: liczba(wartosci, 'paddingLeft'),
          align: tekst(wartosci, 'align'),
        },
      );
      return {
        zadanie: zPolami(
          { frameId: ramka, layout: uklad },
          { layerIds: wykaz(wartosci, 'layerIds') },
        ),
      };
    },
  },
  {
    grupa: 'makieta',
    komenda: Command.DesignConstraintSet,
    nazwa: 'Ustaw więzy',
    opis:
      'Więzy responsywne warstw ramki. Odpowiedź niesie liczbę więzi NAPRAWDĘ zmienionych, ' +
      'nie liczbę nadesłanych.',
    pola: [
      { kod: 'frameId', etykieta: 'Ramka', rodzaj: 'tekst' },
      {
        kod: 'constraints',
        etykieta: 'Więzy',
        rodzaj: 'wielowiersz',
        podpowiedz: '[{"layerId":"warstwa-1","horizontal":"stretch","vertical":"start"}]',
      },
    ],
    zloz(wartosci) {
      const ramka = tekst(wartosci, 'frameId');
      const zapis = tekst(wartosci, 'constraints');
      if (ramka === undefined || zapis === undefined) {
        return { odmowa: 'Więzy wymagają ramki i wykazu więzi.' };
      }
      let wiezy: unknown;
      try {
        wiezy = JSON.parse(zapis);
      } catch {
        return { odmowa: 'Wykaz więzi nie jest poprawnym zapisem JSON.' };
      }
      return { zadanie: { frameId: ramka, constraints: wiezy } };
    },
  },
  {
    grupa: 'makieta',
    komenda: Command.DesignFrameResizeApply,
    nazwa: 'Zmień rozmiar ramki',
    opis: 'Przelicza warstwy z ich więzów i oddaje warstwy PO przeliczeniu.',
    pola: [
      { kod: 'frameId', etykieta: 'Ramka', rodzaj: 'tekst' },
      { kod: 'width', etykieta: 'Nowa szerokość', rodzaj: 'liczba' },
      { kod: 'height', etykieta: 'Nowa wysokość', rodzaj: 'liczba' },
    ],
    zloz(wartosci) {
      const ramka = tekst(wartosci, 'frameId');
      const szerokosc = liczba(wartosci, 'width');
      const wysokosc = liczba(wartosci, 'height');
      if (ramka === undefined || szerokosc === undefined || wysokosc === undefined) {
        return { odmowa: 'Zmiana rozmiaru wymaga ramki i obu wymiarów.' };
      }
      return { zadanie: { frameId: ramka, width: szerokosc, height: wysokosc } };
    },
  },
  {
    grupa: 'makieta',
    komenda: Command.DesignGridSet,
    nazwa: 'Ustaw siatkę',
    opis: 'Siatka ramki albo całej kompozycji. Musi mieć do czego przylgnąć.',
    pola: [
      { kod: 'columns', etykieta: 'Kolumny', rodzaj: 'liczba' },
      { kod: 'gutter', etykieta: 'Rynna', rodzaj: 'liczba' },
      { kod: 'margin', etykieta: 'Margines', rodzaj: 'liczba' },
      { kod: 'baseline', etykieta: 'Siatka bazowa', rodzaj: 'liczba' },
      { kod: 'frameId', etykieta: 'Ramka', rodzaj: 'tekst' },
      { kod: 'boardId', etykieta: 'Kompozycja', rodzaj: 'tekst' },
    ],
    zloz(wartosci) {
      const ramka = tekst(wartosci, 'frameId');
      const kompozycja = tekst(wartosci, 'boardId');
      if (ramka === undefined && kompozycja === undefined) {
        return { odmowa: 'Siatka musi mieć do czego przylgnąć — wskaż ramkę albo kompozycję.' };
      }
      const siatka = zPolami(
        {},
        {
          columns: liczba(wartosci, 'columns'),
          gutter: liczba(wartosci, 'gutter'),
          margin: liczba(wartosci, 'margin'),
          baseline: liczba(wartosci, 'baseline'),
        },
      );
      if (Object.keys(siatka).length === 0) {
        return { odmowa: 'Siatka bez kolumn, rynny, marginesu i linii bazowej nie jest siatką.' };
      }
      return {
        zadanie: zPolami({ grid: siatka }, { frameId: ramka, boardId: kompozycja }),
      };
    },
  },
  {
    grupa: 'makieta',
    komenda: Command.DesignComponentSave,
    nazwa: 'Zapisz komponent',
    opis: 'Komponent okna wraz z wariantami stanu. Zmiana dochodzi do wszystkich instancji.',
    pola: [
      { kod: 'name', etykieta: 'Nazwa', rodzaj: 'tekst' },
      { kod: 'componentId', etykieta: 'Komponent nadpisywany', rodzaj: 'tekst' },
      {
        kod: 'variants',
        etykieta: 'Warianty',
        rodzaj: 'wielowiersz',
        podpowiedz: '[{"name":"spoczynek"},{"name":"hover"}]',
      },
      { kod: 'tokenSetId', etykieta: 'Zestaw żetonów', rodzaj: 'tekst' },
    ],
    zloz(wartosci, otoczenie) {
      const nazwa = tekst(wartosci, 'name');
      if (nazwa === undefined) return { odmowa: 'Komponent wymaga nazwy.' };
      const zapis = tekst(wartosci, 'variants');
      let warianty: unknown;
      if (zapis !== undefined) {
        try {
          warianty = JSON.parse(zapis);
        } catch {
          return { odmowa: 'Wykaz wariantów nie jest poprawnym zapisem JSON.' };
        }
      }
      return {
        zadanie: zPolami(
          { windowId: otoczenie.idOkna, name: nazwa },
          {
            componentId: tekst(wartosci, 'componentId'),
            variants: warianty,
            tokenSetId: tekst(wartosci, 'tokenSetId'),
          },
        ),
      };
    },
  },
  {
    grupa: 'makieta',
    komenda: Command.DesignComponentList,
    nazwa: 'Wykaz komponentów',
    opis: 'Komponenty okna wraz z liczbą instancji każdego.',
    pola: [{ kod: 'componentId', etykieta: 'Zawężenie do komponentu', rodzaj: 'tekst' }],
    zloz(wartosci, otoczenie) {
      return {
        zadanie: zPolami(
          { windowId: otoczenie.idOkna },
          { componentId: tekst(wartosci, 'componentId') },
        ),
      };
    },
  },
  {
    grupa: 'makieta',
    komenda: Command.DesignComponentInstanceAdd,
    nazwa: 'Dołóż instancję komponentu',
    opis: 'Zakłada warstwę instancji i oddaje komponent z NOWĄ liczbą instancji.',
    pola: [
      { kod: 'componentId', etykieta: 'Komponent', rodzaj: 'tekst' },
      { kod: 'boardId', etykieta: 'Kompozycja', rodzaj: 'tekst' },
      { kod: 'variant', etykieta: 'Wariant', rodzaj: 'tekst' },
      { kod: 'frameId', etykieta: 'Ramka', rodzaj: 'tekst' },
      { kod: 'x', etykieta: 'X', rodzaj: 'liczba' },
      { kod: 'y', etykieta: 'Y', rodzaj: 'liczba' },
    ],
    zloz(wartosci) {
      const komponent = tekst(wartosci, 'componentId');
      const kompozycja = tekst(wartosci, 'boardId');
      if (komponent === undefined || kompozycja === undefined) {
        return { odmowa: 'Instancja wymaga komponentu i kompozycji.' };
      }
      return {
        zadanie: zPolami(
          { componentId: komponent, boardId: kompozycja },
          {
            variant: tekst(wartosci, 'variant'),
            frameId: tekst(wartosci, 'frameId'),
            x: liczba(wartosci, 'x'),
            y: liczba(wartosci, 'y'),
          },
        ),
      };
    },
  },
  {
    grupa: 'makieta',
    komenda: Command.DesignPrototypeLinkSet,
    nazwa: 'Ustaw przejście prototypu',
    opis: 'Połączenie między ramkami: wyzwalacz, przejście, czas.',
    pola: [
      { kod: 'boardId', etykieta: 'Kompozycja', rodzaj: 'tekst' },
      { kod: 'fromFrameId', etykieta: 'Ramka źródłowa', rodzaj: 'tekst' },
      { kod: 'toFrameId', etykieta: 'Ramka docelowa', rodzaj: 'tekst' },
      {
        kod: 'trigger',
        etykieta: 'Wyzwalacz',
        rodzaj: 'wybor',
        pozycje: [
          { wartosc: 'click', etykieta: 'Kliknięcie' },
          { wartosc: 'hover', etykieta: 'Najechanie' },
          { wartosc: 'drag', etykieta: 'Przeciągnięcie' },
          { wartosc: 'afterDelay', etykieta: 'Po czasie' },
        ],
      },
      {
        kod: 'transition',
        etykieta: 'Przejście',
        rodzaj: 'wybor',
        pozycje: [
          { wartosc: 'instant', etykieta: 'Natychmiast' },
          { wartosc: 'dissolve', etykieta: 'Przenikanie' },
          { wartosc: 'slideLeft', etykieta: 'Wsuwanie w lewo' },
          { wartosc: 'slideRight', etykieta: 'Wsuwanie w prawo' },
          { wartosc: 'slideUp', etykieta: 'Wsuwanie w górę' },
          { wartosc: 'slideDown', etykieta: 'Wsuwanie w dół' },
          { wartosc: 'push', etykieta: 'Wypychanie' },
        ],
      },
      { kod: 'linkId', etykieta: 'Połączenie zmieniane', rodzaj: 'tekst' },
      { kod: 'durationMs', etykieta: 'Czas przejścia (ms)', rodzaj: 'liczba' },
      { kod: 'layerId', etykieta: 'Warstwa wyzwalająca', rodzaj: 'tekst' },
    ],
    zloz(wartosci) {
      const kompozycja = tekst(wartosci, 'boardId');
      const od = tekst(wartosci, 'fromFrameId');
      const do_ = tekst(wartosci, 'toFrameId');
      const wyzwalacz = tekst(wartosci, 'trigger');
      const przejscie = tekst(wartosci, 'transition');
      if (
        kompozycja === undefined ||
        od === undefined ||
        do_ === undefined ||
        wyzwalacz === undefined ||
        przejscie === undefined
      ) {
        return {
          odmowa: 'Przejście wymaga kompozycji, obu ramek, wyzwalacza i rodzaju przejścia.',
        };
      }
      return {
        zadanie: zPolami(
          {
            boardId: kompozycja,
            fromFrameId: od,
            toFrameId: do_,
            trigger: wyzwalacz,
            transition: przejscie,
          },
          {
            linkId: tekst(wartosci, 'linkId'),
            durationMs: liczba(wartosci, 'durationMs'),
            layerId: tekst(wartosci, 'layerId'),
          },
        ),
      };
    },
  },
  {
    grupa: 'makieta',
    komenda: Command.DesignPrototypeGet,
    nazwa: 'Odczyt prototypu',
    opis:
      'Ramki i przejścia. Odpowiedź niesie ramki NIEOSIĄGALNE — prototyp z ramką ' +
      'osieroconą wygląda bez tego jak prototyp kompletny.',
    pola: [
      { kod: 'boardId', etykieta: 'Kompozycja', rodzaj: 'tekst' },
      { kod: 'startFrameId', etykieta: 'Ramka początkowa', rodzaj: 'tekst' },
    ],
    zloz(wartosci) {
      const kompozycja = tekst(wartosci, 'boardId');
      if (kompozycja === undefined) return { odmowa: 'Odczyt prototypu wymaga kompozycji.' };
      return {
        zadanie: zPolami(
          { boardId: kompozycja },
          { startFrameId: tekst(wartosci, 'startFrameId') },
        ),
      };
    },
  },
  {
    grupa: 'makieta',
    komenda: Command.DesignPrototypeLinkRemove,
    nazwa: 'Usuń przejście',
    opis: 'Usuwa połączenie między ramkami.',
    pola: [{ kod: 'linkId', etykieta: 'Połączenie', rodzaj: 'tekst' }],
    zloz(wartosci) {
      const polaczenie = tekst(wartosci, 'linkId');
      if (polaczenie === undefined) return { odmowa: 'Usunięcie wymaga wskazania połączenia.' };
      return { zadanie: { linkId: polaczenie } };
    },
  },
  {
    grupa: 'makieta',
    komenda: Command.DesignMockupGenerate,
    nazwa: 'Zbuduj makietę z opisu',
    opis:
      'Rdzeń rozkłada opis na SEKCJE i układa je pionowo w ramce. Wynik jest UKŁADEM ' +
      'warstw, nie obrazkiem — da się go dalej przesuwać.',
    pola: [
      { kod: 'boardId', etykieta: 'Kompozycja', rodzaj: 'tekst' },
      {
        kod: 'prompt',
        etykieta: 'Opis ekranu',
        rodzaj: 'wielowiersz',
        podpowiedz: 'nagłówek, nawigacja, bohater, karty, stopka',
      },
      { kod: 'devicePreset', etykieta: 'Nastawa urządzenia', rodzaj: 'tekst' },
      { kod: 'tokenSetId', etykieta: 'Zestaw żetonów', rodzaj: 'tekst' },
      POLE_KANALU,
    ],
    zloz(wartosci, otoczenie) {
      const kompozycja = tekst(wartosci, 'boardId');
      const opis = tekst(wartosci, 'prompt');
      if (kompozycja === undefined || opis === undefined) {
        return { odmowa: 'Makieta wymaga kompozycji i opisu ekranu — rdzeń go nie wymyśla.' };
      }
      return {
        zadanie: zPolami(
          { windowId: otoczenie.idOkna, boardId: kompozycja, prompt: opis },
          {
            devicePreset: tekst(wartosci, 'devicePreset'),
            tokenSetId: tekst(wartosci, 'tokenSetId'),
            channelId: tekst(wartosci, 'channelId'),
          },
        ),
      };
    },
  },
  {
    grupa: 'makieta',
    komenda: Command.DesignMockupImport,
    nazwa: 'Wczytaj makietę ze zrzutu',
    opis:
      'Rdzeń WYKRYWA układ obszarów na zrzucie i zaznacza linie tekstu. Liter nie czyta — ' +
      'to wymagałoby biblioteki, której w drzewie nie ma, i jest zgłoszone jako brak. ' +
      'Obszary nierozłożone wracają w bilansie.',
    pola: [
      { kod: 'boardId', etykieta: 'Kompozycja', rodzaj: 'tekst' },
      POLE_ZASOBU,
      { kod: 'recognizeText', etykieta: 'Znacz linie tekstu', rodzaj: 'logiczne' },
      { kod: 'devicePreset', etykieta: 'Nastawa urządzenia', rodzaj: 'tekst' },
    ],
    zloz(wartosci, otoczenie) {
      const kompozycja = tekst(wartosci, 'boardId');
      if (kompozycja === undefined) return { odmowa: 'Wczytanie makiety wymaga kompozycji.' };
      return {
        zadanie: zPolami(
          {
            windowId: otoczenie.idOkna,
            boardId: kompozycja,
            assetId: wartosci['assetId'] ?? '',
          },
          {
            recognizeText: przelacznik(wartosci, 'recognizeText'),
            devicePreset: tekst(wartosci, 'devicePreset'),
          },
        ),
      };
    },
  },
  {
    grupa: 'publikacja',
    komenda: Command.DesignProductMockupRender,
    nazwa: 'Wyrysuj makietę produktową',
    opis: 'Nakłada projekt na zdjęcie produktu z obrotem i kryciem.',
    pola: [
      { kod: 'designAssetId', etykieta: 'Projekt', rodzaj: 'zasob' },
      { kod: 'productAssetId', etykieta: 'Zdjęcie produktu', rodzaj: 'zasob' },
      { kod: 'x', etykieta: 'X nałożenia', rodzaj: 'liczba' },
      { kod: 'y', etykieta: 'Y nałożenia', rodzaj: 'liczba' },
      { kod: 'width', etykieta: 'Szerokość nałożenia', rodzaj: 'liczba' },
      { kod: 'height', etykieta: 'Wysokość nałożenia', rodzaj: 'liczba' },
      { kod: 'rotation', etykieta: 'Obrót (°)', rodzaj: 'liczba' },
      { kod: 'opacity', etykieta: 'Krycie (0–1)', rodzaj: 'liczba' },
    ],
    zloz(wartosci, otoczenie) {
      const projekt = tekst(wartosci, 'designAssetId');
      const produkt = tekst(wartosci, 'productAssetId');
      const x = liczba(wartosci, 'x');
      const y = liczba(wartosci, 'y');
      const szerokosc = liczba(wartosci, 'width');
      const wysokosc = liczba(wartosci, 'height');
      if (
        projekt === undefined ||
        produkt === undefined ||
        x === undefined ||
        y === undefined ||
        szerokosc === undefined ||
        wysokosc === undefined
      ) {
        return {
          odmowa:
            'Makieta produktowa wymaga obu zasobów oraz położenia i wymiarów nałożenia.',
        };
      }
      return {
        zadanie: zPolami(
          {
            designAssetId: projekt,
            productAssetId: produkt,
            x,
            y,
            width: szerokosc,
            height: wysokosc,
            windowId: otoczenie.idOkna,
          },
          {
            rotation: liczba(wartosci, 'rotation'),
            opacity: liczba(wartosci, 'opacity'),
          },
        ),
      };
    },
  },
];

/** Czynności jednego warsztatu modułu Design. Okno bierze z katalogu wyłącznie czynności własnej grupy warsztatu i pomija pozostałe. */
export function czynnosciWarsztatu(grupa: GrupaWarsztatu): readonly CzynnoscWarsztatuDesignu[] {
  return CZYNNOSCI_WARSZTATOW_DESIGNU.filter((czynnosc) => czynnosc.grupa === grupa);
}
