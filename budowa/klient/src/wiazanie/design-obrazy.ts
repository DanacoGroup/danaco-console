// Obróbka obrazów w oknie Design: odczyt danych zasobu, przekształcenia,
// poprawki, powiększenie, zdjęcie tła, wektoryzacja, złożenie i rozbiór na warstwy.
import {
  Command,
  ImageAdjustKind,
  ImageBlendMode,
  ImageTransformKind,
  ImageVectorizeMode,
} from '../../../shared/contract.ts';
import type { DesignAsset } from '../../../shared/contract.ts';
import type { Kanal } from '../protokol/kanal.ts';
import { wywolaj } from '../protokol/wywolanie.ts';
import { dolozCzynnosciPanelu, zapytajWSzufladzie } from './czynnosci-okna.ts';
import { oglos } from './ogloszenie.ts';

const NAGLOWEK = 'Obrazy';

interface Otoczenie {
  kanal: Kanal;
  korzen: Element;
  idOkna: () => string;
  odswiez: () => void;
}

export function zwiazObrazyDesignu(
  kanal: Kanal,
  korzen: Element,
  idOkna: () => string,
  odswiez: () => void,
  dolacz: (zdejmij: () => void) => void,
  przy: AddEventListenerOptions,
): void {
  const otoczenie: Otoczenie = { kanal, korzen, idOkna, odswiez };
  const zdejmij = dolozCzynnosciPanelu(korzen, 'panel-assets', 'Obróbka obrazu', [
    {
      naglowek: 'Wgląd',
      pozycje: [{ kod: 'dane', nazwa: 'Odczytaj dane obrazu…' }],
    },
    {
      naglowek: 'Postać',
      pozycje: [
        { kod: 'postac', nazwa: 'Zmień postać zapisu…' },
        { kod: 'przeksztalc', nazwa: 'Przekształć obraz…' },
        { kod: 'popraw', nazwa: 'Popraw obraz…' },
        { kod: 'powieksz', nazwa: 'Powiększ obraz…' },
      ],
    },
    {
      naglowek: 'Praca na treści',
      pozycje: [
        { kod: 'tlo', nazwa: 'Zdejmij tło…' },
        { kod: 'wektor', nazwa: 'Zamień na wektor…' },
        { kod: 'zloz', nazwa: 'Złóż dwa obrazy…' },
        { kod: 'warstwy', nazwa: 'Rozbierz na warstwy…' },
      ],
    },
  ], (kod) => {
    void wykonaj(otoczenie, kod, 'panel-assets');
  }, przy);
  if (zdejmij !== null) dolacz(zdejmij);
}

/* Zasób wskazuje wykaz okna: obraz spoza niego nie ma jak trafić do rdzenia,
   bo okno nie sięga do plików maszyny Operatora. */
async function wyborZasobow(
  kanal: Kanal,
  idOkna: string,
): Promise<ReadonlyArray<readonly [string, string]>> {
  const wykaz = await wywolaj(kanal, Command.DesignAssetList, { windowId: idOkna });
  return (wykaz.wynik?.assets ?? []).map((zasob: DesignAsset) =>
    [zasob.id, `${zasob.name ?? zasob.id} · ${zasob.format ?? zasob.kind}`] as const);
}

async function wskazZasob(
  otoczenie: Otoczenie,
  panel: string,
  tytul: string,
  wykonanie: string,
  dodatkowe: Parameters<typeof zapytajWSzufladzie>[2]['pola'] = [],
  opis?: string,
): Promise<Record<string, string> | null> {
  const zasoby = await wyborZasobow(otoczenie.kanal, otoczenie.idOkna());
  if (zasoby.length === 0) {
    oglos(NAGLOWEK, 'To okno nie ma jeszcze żadnego zasobu obrazowego.', 'ostrzezenie');
    return null;
  }
  return zapytajWSzufladzie(otoczenie.korzen, panel, {
    tytul,
    ...(opis === undefined ? {} : { opis }),
    pola: [{ klucz: 'zasob', etykieta: 'Obraz', wybor: zasoby }, ...dodatkowe],
    wykonanie,
  });
}

async function wykonaj(otoczenie: Otoczenie, kod: string, panel: string): Promise<void> {
  if (kod === 'dane') return odczytajDane(otoczenie, panel);
  if (kod === 'postac') return zmienPostac(otoczenie, panel);
  if (kod === 'przeksztalc') return przeksztalc(otoczenie, panel);
  if (kod === 'popraw') return popraw(otoczenie, panel);
  if (kod === 'powieksz') return powieksz(otoczenie, panel);
  if (kod === 'tlo') return zdejmijTlo(otoczenie, panel);
  if (kod === 'wektor') return zamienNaWektor(otoczenie, panel);
  if (kod === 'zloz') return zlozObrazy(otoczenie, panel);
  if (kod === 'warstwy') return rozbierzNaWarstwy(otoczenie, panel);
  /* Czynność spoza obsłużonych odmawia zamiast milczeć: pozycja bez gałęzi
     wyglądałaby jak działająca. */
  oglos(NAGLOWEK, `Czynność „${kod}" nie jest prowadzona przez to okno.`, 'ostrzezenie');
}

async function odczytajDane(otoczenie: Otoczenie, panel: string): Promise<void> {
  const wartosci = await wskazZasob(otoczenie, panel, 'Dane obrazu', 'Odczytaj dane');
  if (wartosci === null) return;
  const wynik = await wywolaj(otoczenie.kanal, Command.ImageInspect, {
    assetId: wartosci.zasob ?? '',
  });
  if (!wynik.udany || wynik.wynik === undefined) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił odczytu danych obrazu.', 'ostrzezenie');
    return;
  }
  const dane = wynik.wynik;
  oglos(NAGLOWEK, `${dane.format} · ${dane.width}×${dane.height} · ${dane.sizeBytes} bajtów`
    + `${dane.colorSpace === undefined ? '' : ` · ${dane.colorSpace}`}`);
}

/* Wynik obróbki wraca nowym zasobem, więc wykaz odświeża się sam, a stary
   obraz zostaje nietknięty. */
async function zmienPostac(otoczenie: Otoczenie, panel: string): Promise<void> {
  const wartosci = await wskazZasob(otoczenie, panel, 'Postać zapisu obrazu', 'Zmień postać', [
    {
      klucz: 'postac',
      etykieta: 'Postać docelowa',
      wybor: [['webp', 'WebP'], ['png', 'PNG'], ['jpeg', 'JPEG'], ['avif', 'AVIF']],
    },
    { klucz: 'jakosc', etykieta: 'Jakość', wartosc: '85' },
  ], 'Rdzeń odda nowy zasób; obraz źródłowy zostaje nietknięty.');
  if (wartosci === null) return;
  const jakosc = Number.parseInt(wartosci.jakosc ?? '', 10);
  const wynik = await wywolaj(otoczenie.kanal, Command.ImageConvert, {
    windowId: otoczenie.idOkna(),
    assetId: wartosci.zasob ?? '',
    format: wartosci.postac ?? 'webp',
    ...(Number.isFinite(jakosc) ? { quality: jakosc } : {}),
  });
  if (!wynik.udany || wynik.wynik === undefined) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił zmiany postaci.', 'ostrzezenie');
    return;
  }
  oglos(NAGLOWEK, `Obraz zapisany jako ${wartosci.postac ?? ''}; ${wynik.wynik.sizeBytes} bajtów`
    + `${wynik.wynik.savedBytes === undefined
      ? '' : `, zaoszczędzono ${wynik.wynik.savedBytes}`}.`);
  otoczenie.odswiez();
}

async function przeksztalc(otoczenie: Otoczenie, panel: string): Promise<void> {
  const wartosci = await wskazZasob(otoczenie, panel, 'Przekształcenie obrazu', 'Przekształć', [
    {
      klucz: 'rodzaj',
      etykieta: 'Przekształcenie',
      wybor: [
        [ImageTransformKind.Resize, 'Zmiana rozmiaru'],
        [ImageTransformKind.Crop, 'Wycięcie'],
        [ImageTransformKind.Rotate, 'Obrót'],
        [ImageTransformKind.FlipHorizontal, 'Odbicie w poziomie'],
        [ImageTransformKind.FlipVertical, 'Odbicie w pionie'],
        [ImageTransformKind.Thumbnail, 'Miniatura'],
      ],
    },
    { klucz: 'szerokosc', etykieta: 'Szerokość' },
    { klucz: 'wysokosc', etykieta: 'Wysokość' },
    { klucz: 'stopnie', etykieta: 'Stopnie obrotu' },
  ]);
  if (wartosci === null) return;
  const szerokosc = Number.parseInt(wartosci.szerokosc ?? '', 10);
  const wysokosc = Number.parseInt(wartosci.wysokosc ?? '', 10);
  const stopnie = Number.parseInt(wartosci.stopnie ?? '', 10);
  const wynik = await wywolaj(otoczenie.kanal, Command.ImageTransform, {
    windowId: otoczenie.idOkna(),
    assetId: wartosci.zasob ?? '',
    operation: (wartosci.rodzaj ?? ImageTransformKind.Resize) as ImageTransformKind,
    ...(Number.isFinite(szerokosc) ? { width: szerokosc } : {}),
    ...(Number.isFinite(wysokosc) ? { height: wysokosc } : {}),
    ...(Number.isFinite(stopnie) ? { degrees: stopnie } : {}),
    keepAspect: true,
  });
  if (!wynik.udany) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił przekształcenia.', 'ostrzezenie');
    return;
  }
  oglos(NAGLOWEK, 'Obraz przekształcony; nowy zasób stoi w wykazie.');
  otoczenie.odswiez();
}

async function popraw(otoczenie: Otoczenie, panel: string): Promise<void> {
  const wartosci = await wskazZasob(otoczenie, panel, 'Poprawka obrazu', 'Popraw obraz', [
    {
      klucz: 'rodzaj',
      etykieta: 'Poprawka',
      wybor: [
        [ImageAdjustKind.Brightness, 'Jasność'],
        [ImageAdjustKind.Contrast, 'Kontrast'],
        [ImageAdjustKind.Saturation, 'Nasycenie'],
        [ImageAdjustKind.Sharpen, 'Wyostrzenie'],
        [ImageAdjustKind.Blur, 'Rozmycie'],
        [ImageAdjustKind.Denoise, 'Odszumienie'],
        [ImageAdjustKind.Grayscale, 'Skala szarości'],
        [ImageAdjustKind.AutoLevels, 'Wyrównanie poziomów'],
      ],
    },
    { klucz: 'natezenie', etykieta: 'Natężenie poprawki', podpowiedz: 'Puste bierze wartość rdzenia' },
  ]);
  if (wartosci === null) return;
  const natezenie = Number.parseFloat(wartosci.natezenie ?? '');
  const wynik = await wywolaj(otoczenie.kanal, Command.ImageAdjust, {
    windowId: otoczenie.idOkna(),
    assetId: wartosci.zasob ?? '',
    operation: (wartosci.rodzaj ?? ImageAdjustKind.AutoLevels) as ImageAdjustKind,
    ...(Number.isFinite(natezenie) ? { amount: natezenie } : {}),
  });
  if (!wynik.udany) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił poprawki.', 'ostrzezenie');
    return;
  }
  oglos(NAGLOWEK, 'Poprawka wykonana; nowy zasób stoi w wykazie.');
  otoczenie.odswiez();
}

async function powieksz(otoczenie: Otoczenie, panel: string): Promise<void> {
  const wartosci = await wskazZasob(otoczenie, panel, 'Powiększenie obrazu', 'Powiększ', [
    { klucz: 'krotnosc', etykieta: 'Ile razy', wartosc: '2' },
    {
      klucz: 'twarze',
      etykieta: 'Poprawa twarzy',
      wybor: [['nie', 'Bez poprawy'], ['tak', 'Z poprawą twarzy']],
    },
  ]);
  if (wartosci === null) return;
  const krotnosc = Number.parseFloat(wartosci.krotnosc ?? '');
  const wynik = await wywolaj(otoczenie.kanal, Command.ImageUpscale, {
    windowId: otoczenie.idOkna(),
    assetId: wartosci.zasob ?? '',
    ...(Number.isFinite(krotnosc) ? { scale: krotnosc } : {}),
    faces: wartosci.twarze === 'tak',
  });
  if (!wynik.udany) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił powiększenia.', 'ostrzezenie');
    return;
  }
  oglos(NAGLOWEK, 'Obraz powiększony; nowy zasób stoi w wykazie.');
  otoczenie.odswiez();
}

async function zdejmijTlo(otoczenie: Otoczenie, panel: string): Promise<void> {
  const wartosci = await wskazZasob(otoczenie, panel, 'Zdjęcie tła', 'Zdejmij tło');
  if (wartosci === null) return;
  const wynik = await wywolaj(otoczenie.kanal, Command.ImageBackgroundRemove, {
    windowId: otoczenie.idOkna(),
    assetId: wartosci.zasob ?? '',
  });
  if (!wynik.udany) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił zdjęcia tła.', 'ostrzezenie');
    return;
  }
  oglos(NAGLOWEK, 'Tło zdjęte; nowy zasób stoi w wykazie.');
  otoczenie.odswiez();
}

async function zamienNaWektor(otoczenie: Otoczenie, panel: string): Promise<void> {
  const wartosci = await wskazZasob(otoczenie, panel, 'Zamiana na wektor', 'Zamień na wektor', [
    {
      klucz: 'sposob',
      etykieta: 'Sposób odrysu',
      wybor: [
        [ImageVectorizeMode.Outline, 'Po obrysie'],
        [ImageVectorizeMode.Centerline, 'Po osi kształtu'],
        [ImageVectorizeMode.Posterize, 'Płaskimi plamami'],
      ],
    },
    { klucz: 'barwy', etykieta: 'Liczba barw', wartosc: '8' },
  ]);
  if (wartosci === null) return;
  const barwy = Number.parseInt(wartosci.barwy ?? '', 10);
  const wynik = await wywolaj(otoczenie.kanal, Command.ImageVectorize, {
    windowId: otoczenie.idOkna(),
    assetId: wartosci.zasob ?? '',
    mode: (wartosci.sposob ?? ImageVectorizeMode.Outline) as ImageVectorizeMode,
    ...(Number.isFinite(barwy) ? { colors: barwy } : {}),
  });
  if (!wynik.udany) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił wektoryzacji.', 'ostrzezenie');
    return;
  }
  oglos(NAGLOWEK, 'Obraz odrysowany wektorem; nowy zasób stoi w wykazie.');
  otoczenie.odswiez();
}

/* Złożenie potrzebuje dwóch zasobów, więc szuflada pyta o oba z tego samego
   wykazu — obraz nakładany może być tym samym, co podkład. */
async function zlozObrazy(otoczenie: Otoczenie, panel: string): Promise<void> {
  const zasoby = await wyborZasobow(otoczenie.kanal, otoczenie.idOkna());
  if (zasoby.length < 2) {
    oglos(NAGLOWEK, 'Do złożenia potrzeba dwóch zasobów obrazowych.', 'ostrzezenie');
    return;
  }
  const wartosci = await zapytajWSzufladzie(otoczenie.korzen, panel, {
    tytul: 'Złożenie obrazów',
    pola: [
      { klucz: 'podklad', etykieta: 'Obraz spodni', wybor: zasoby },
      { klucz: 'naklad', etykieta: 'Obraz nakładany', wybor: zasoby },
      { klucz: 'poziomo', etykieta: 'Odsunięcie w poziomie', wartosc: '0' },
      { klucz: 'pionowo', etykieta: 'Odsunięcie w pionie', wartosc: '0' },
      { klucz: 'krycie', etykieta: 'Krycie', wartosc: '1' },
      {
        klucz: 'mieszanie',
        etykieta: 'Sposób mieszania',
        wybor: [
          [ImageBlendMode.Normal, 'Zwykłe'],
          [ImageBlendMode.Multiply, 'Mnożenie'],
          [ImageBlendMode.Screen, 'Rozjaśnianie'],
          [ImageBlendMode.Overlay, 'Nakładka'],
        ],
      },
    ],
    wykonanie: 'Złóż obrazy',
  });
  if (wartosci === null) return;
  const poziomo = Number.parseInt(wartosci.poziomo ?? '', 10);
  const pionowo = Number.parseInt(wartosci.pionowo ?? '', 10);
  const krycie = Number.parseFloat(wartosci.krycie ?? '');
  const wynik = await wywolaj(otoczenie.kanal, Command.ImageCompose, {
    windowId: otoczenie.idOkna(),
    baseAssetId: wartosci.podklad ?? '',
    overlayAssetId: wartosci.naklad ?? '',
    ...(Number.isFinite(poziomo) ? { x: poziomo } : {}),
    ...(Number.isFinite(pionowo) ? { y: pionowo } : {}),
    ...(Number.isFinite(krycie) ? { opacity: krycie } : {}),
    blendMode: (wartosci.mieszanie ?? ImageBlendMode.Normal) as ImageBlendMode,
  });
  if (!wynik.udany) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił złożenia obrazów.', 'ostrzezenie');
    return;
  }
  oglos(NAGLOWEK, 'Obrazy złożone; nowy zasób stoi w wykazie.');
  otoczenie.odswiez();
}

async function rozbierzNaWarstwy(otoczenie: Otoczenie, panel: string): Promise<void> {
  const wartosci = await wskazZasob(otoczenie, panel, 'Rozbiór na warstwy', 'Rozbierz obraz', [
    { klucz: 'warstwy', etykieta: 'Najwyżej warstw', wartosc: '4' },
  ]);
  if (wartosci === null) return;
  const warstwy = Number.parseInt(wartosci.warstwy ?? '', 10);
  const wynik = await wywolaj(otoczenie.kanal, Command.ImageLayersSplit, {
    windowId: otoczenie.idOkna(),
    assetId: wartosci.zasob ?? '',
    ...(Number.isFinite(warstwy) ? { maxLayers: warstwy } : {}),
  });
  if (!wynik.udany || wynik.wynik === undefined) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił rozbioru na warstwy.', 'ostrzezenie');
    return;
  }
  oglos(NAGLOWEK, `Obraz rozebrany na ${wynik.wynik.layers} warstw; `
    + `${wynik.wynik.assets.length} zasobów stoi w wykazie.`);
  otoczenie.odswiez();
}
