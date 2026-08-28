import type {
  Channel,
  DesignAsset,
  DesignAssetListResponse,
  Module,
} from '../../../../shared/contract';
import { opisOdmowyBledu } from '../../komponenty/odmowa';
import type { Wynik } from '../../protokol/kanal';
import { KOD_MODULU } from './etykiety-designu';
import type { ZapytanieZasobow } from './zrodlo-designu';
import type { ZrodloZaplecza } from './zrodlo-zaplecza';

/**
 * Zapis modułu Design wraz z odczytami, które go wypełniają. Odpowiada
 * wyłącznie za prawdę o zasobach, oknie modułu i rejestrach zaplecza.
 */
export type FazaZasobow = 'spoczynek' | 'odczyt' | 'gotowe' | 'blad';

export interface ZapisDesignu {
  zbior: DesignAsset[];
  wybor: string | null;
  kanaly: Channel[];
  cele: Module[];
  faza: FazaZasobow;
  powod: string;
  /** Czy powód mówi o braku rozstrzygnięcia, a nie o odmowie rdzenia. */
  bezRozstrzygniecia: boolean;
  warunki: ZapytanieZasobow;
  oknoModulu: string;
  zdanieOOknie: string;
}

export function pustyZapisDesignu(): ZapisDesignu {
  return {
    zbior: [],
    wybor: null,
    kanaly: [],
    cele: [],
    faza: 'spoczynek',
    powod: '',
    bezRozstrzygniecia: false,
    // Granica wyjściowa chroni odpowiedź przed wykazem, którego nie uniesie.
    warunki: { idOkna: '', rodzaj: '', etykiety: [], tylkoUlubione: false, granica: 60 },
    oknoModulu: '',
    zdanieOOknie: 'Okno modułu nie zostało jeszcze ustalone z rdzenia.',
  };
}

/** Wciąga zasób po własnej zmianie albo po zdarzeniu `design.asset.changed`, aktualizując zbiór zasobów modułu. */
export function wchlonZasob(zapis: ZapisDesignu, zasob: DesignAsset): void {
  const pozycja = zapis.zbior.findIndex((wpis) => wpis.id === zasob.id);
  if (pozycja === -1) zapis.zbior = [...zapis.zbior, zasob];
  else zapis.zbior = zapis.zbior.map((wpis) => (wpis.id === zasob.id ? zasob : wpis));
  if (zapis.wybor === null) zapis.wybor = zasob.id;
  zapis.faza = 'gotowe';
}

/**
 * Zdejmuje zasób usunięty w rdzeniu wraz z jego wyborem, odbierając zdarzenie
 * rodzaju `deleted` zdarzenia `design.asset.changed`.
 */
export function usunZasob(zapis: ZapisDesignu, idZasobu: string): void {
  zapis.zbior = zapis.zbior.filter((wpis) => wpis.id !== idZasobu);
  if (zapis.wybor === idZasobu) zapis.wybor = zapis.zbior[0]?.id ?? null;
}

/**
 * Przyjęcie odpowiedzi na `design.asset.list`. Odmowa zostawia zbiór
 * nietknięty i zapisuje powód, odróżniając pusty wykaz po odmowie od pustego
 * wykazu na świeżej instalacji.
 */
export function przyjmijOdczytZasobow(
  zapis: ZapisDesignu,
  wynik: Wynik<DesignAssetListResponse>,
): void {
  zapis.bezRozstrzygniecia = false;
  if (!wynik.udany || wynik.wynik === undefined) {
    zapis.faza = 'blad';
    zapis.powod = opisOdmowyBledu('Odczyt zasobów', wynik.blad);
    return;
  }
  zapis.zbior = wynik.wynik.assets;
  if (zapis.wybor !== null && !zapis.zbior.some((wpis) => wpis.id === zapis.wybor)) {
    zapis.wybor = null;
  }
  zapis.faza = 'gotowe';
  zapis.powod = '';
}

/**
 * Ustala okno modułu w sesji i odczytuje jego stan.
 *
 * Komendy obszaru `design.*` wymagają `windowId`, a klient go nie wymyśla:
 * bierze okno tej sesji, którego `moduleId` jest kodem tego modułu. Brak
 * takiego okna jest stanem opisanym, nie wyjątkiem.
 */
export async function ustalOknoModulu(
  zapis: ZapisDesignu,
  zaplecze: ZrodloZaplecza,
  idSesji: string,
): Promise<void> {
  if (idSesji === '') {
    zapis.zdanieOOknie = 'Sesja nie jest jeszcze otwarta — okno modułu nie ma identyfikatora.';
    return;
  }
  const wykaz = await zaplecze.okna(idSesji);
  if (!wykaz.udany || wykaz.wynik === undefined) {
    zapis.zdanieOOknie = opisOdmowyBledu('Odczyt okien sesji', wykaz.blad);
    return;
  }
  const wlasne = wykaz.wynik.windows.find((okno) => okno.moduleId === KOD_MODULU);
  if (wlasne === undefined) {
    zapis.oknoModulu = '';
    zapis.zdanieOOknie =
      'Sesja nie ma jeszcze okna modułu Design — komendy obszaru wymagają jego identyfikatora.';
    return;
  }
  zapis.oknoModulu = wlasne.id;
  zapis.warunki = { ...zapis.warunki, idOkna: wlasne.id };
  const stan = await zaplecze.stanOkna(wlasne.id);
  zapis.zdanieOOknie =
    stan.udany && stan.wynik !== undefined
      ? `Okno ${wlasne.id}: proces ${stan.wynik.processStatus}, wiadomości ${stan.wynik.messageCount}.`
      : opisOdmowyBledu('Odczyt stanu okna', stan.blad);
}

/**
 * Odczyt rejestrów zaplecza: kanałów modelu i katalogu modułów.
 *
 * Oba są cudze i niezależne, więc idą równolegle i odmowa jednego nie gasi
 * drugiego.
 */
export async function odczytajZaplecze(
  zapis: ZapisDesignu,
  zaplecze: ZrodloZaplecza,
): Promise<void> {
  const [silniki, moduly] = await Promise.all([
    zaplecze.silniki(),
    zaplecze.modulyDocelowe(),
  ]);
  zapis.kanaly = silniki.udany && silniki.wynik !== undefined ? silniki.wynik.channels : [];
  zapis.cele =
    moduly.udany && moduly.wynik !== undefined
      ? moduly.wynik.modules.filter((modul) => modul.code !== KOD_MODULU)
      : [];
}
