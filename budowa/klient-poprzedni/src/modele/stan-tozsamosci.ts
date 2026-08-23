import {
  ChangeKind,
  ConfigAxis,
  type IdentityCategory,
  type IdentityDocument,
  type IdentityEffectiveGetResponse,
  type IdentityMode,
} from '../../../shared/contract';
import type { Odsubskrybuj } from '../polaczenie/magistrala-zdarzen';
import type { Kanal, Wynik } from '../protokol/kanal';
import { wskazanieKompletne, type WskazanieOsi } from './wybor-osi';
import { utworzZrodloTozsamosci, type ZrodloTozsamosci } from './zrodlo-tozsamosci';

/**
 * Stan zasad i tożsamości modelu — katalog kategorii, treści zapisane dla
 * jednej osi oraz nakładka obowiązująca.
 *
 * Trzy rzeczy w jednym stanie, ponieważ edytor potrzebuje wszystkich naraz:
 * katalog mówi, jakie kategorie istnieją i który tryb proponują, zapisy mówią,
 * co w nich stoi, a nakładka — co z tego trafi do modelu.
 *
 * Oś wyznacza zakres odczytu: zapisy pobieramy dla osi wskazanej, a nie dla
 * wszystkich naraz, więc kategoria bez zapisu na tej osi ma treść pustą,
 * a obowiązuje dla niej treść z osi szerszej. Oś, która wymaga bytu, a bytu nie
 * ma, jest traktowana jak platforma — adresowanie osi bez bytu nie ma
 * w kontrakcie znaczenia.
 */
export interface StanTozsamosci {
  /** Kategorie w kolejności warstw i porządku wewnątrz warstwy. */
  kategorie(): readonly IdentityCategory[];
  /** Zapis treści kategorii dla osi czynnej; `null`, gdy zapisu nie ma. */
  dokument(kategoria: string): IdentityDocument | null;
  /** Oś i byt, dla których czytamy i zapisujemy treści. */
  wskazanie(): WskazanieOsi;
  /** Zmienia oś i wczytuje jej zapisy oraz nakładkę. */
  ustawWskazanie(wskazanie: WskazanieOsi): Promise<void>;
  /** Kategoria czynna w edytorze; `null`, gdy katalog jest pusty. */
  wybrana(): IdentityCategory | null;
  /** Ustawia kategorię czynną. */
  wybierz(kategoria: string): void;
  /** Nakładka obowiązująca; `null`, gdy jej jeszcze nie odczytano. */
  nakladka(): IdentityEffectiveGetResponse | null;
  /** Wczytuje katalog, zapisy osi oraz nakładkę. */
  odswiez(): Promise<void>;
  /** `identity.document.set` — zapis treści kategorii wraz z trybem. */
  zapisz(kategoria: string, tresc: string, tryb: IdentityMode): Promise<Wynik<unknown>>;
  /** `identity.document.remove` — zdjęcie zapisu kategorii z osi czynnej. */
  usun(dokument: string): Promise<Wynik<unknown>>;
  /** Subskrypcja przeliczenia stanu. */
  naZmiane(sluchacz: () => void): void;
  /** Odłącza subskrypcję zdarzeń kanału. */
  rozlacz(): void;
}

export function utworzStanTozsamosci(kanal: Kanal): StanTozsamosci {
  const zrodlo: ZrodloTozsamosci = utworzZrodloTozsamosci(kanal);

  let kategorie: IdentityCategory[] = [];
  let dokumenty: IdentityDocument[] = [];
  let nakladka: IdentityEffectiveGetResponse | null = null;
  let wybrana: string | null = null;
  let wskazanie: WskazanieOsi = { os: ConfigAxis.Platform, bytOsi: '' };

  const sluchacze: Array<() => void> = [];
  const oglos = (): void => {
    for (const sluchacz of [...sluchacze]) sluchacz();
  };

  /** Oś skuteczna: wskazanie bez bytu znaczy platformę. */
  const skuteczne = (): WskazanieOsi =>
    wskazanieKompletne(wskazanie) ? wskazanie : { os: ConfigAxis.Platform, bytOsi: '' };

  /** Żądanie adresujące oś czynną; byt pusty pomijamy, nie wysyłamy pustki. */
  const adres = (): { axis: ConfigAxis; axisId?: string } => {
    const os = skuteczne();
    return os.bytOsi === '' ? { axis: os.os } : { axis: os.os, axisId: os.bytOsi };
  };

  /** Czy zapis dotyczy osi czynnej — zdarzenie z innej osi wykazu nie rusza. */
  const naOsiCzynnej = (dokument: IdentityDocument): boolean => {
    const os = skuteczne();
    return dokument.axis === os.os && (dokument.axisId ?? '') === os.bytOsi;
  };

  const odsubskrybuj: Odsubskrybuj = zrodlo.naZmiane((tresc) => {
    if (!naOsiCzynnej(tresc.document)) return;
    dokumenty =
      tresc.change === ChangeKind.Deleted
        ? dokumenty.filter((zapis) => zapis.id !== tresc.document.id)
        : [
            ...dokumenty.filter((zapis) => zapis.categoryId !== tresc.document.categoryId),
            tresc.document,
          ];
    oglos();
    void wczytajNakladke();
  });

  async function wczytajNakladke(): Promise<void> {
    const wynik = await zrodlo.nakladka(adres());
    nakladka = wynik.udany ? (wynik.wynik ?? null) : null;
    oglos();
  }

  async function wczytajZapisy(): Promise<void> {
    dokumenty = await zrodlo.dokumenty(adres());
    oglos();
  }

  async function wczytaj(): Promise<void> {
    kategorie = await zrodlo.kategorie({});
    if (wybrana === null || !kategorie.some((kategoria) => kategoria.id === wybrana)) {
      wybrana = kategorie[0]?.id ?? null;
    }
    await Promise.all([wczytajZapisy(), wczytajNakladke()]);
  }

  return {
    kategorie: () => kategorie,

    dokument: (kategoria) =>
      dokumenty.find((zapis) => zapis.categoryId === kategoria) ?? null,

    wskazanie: () => wskazanie,

    async ustawWskazanie(nowe) {
      wskazanie = nowe;
      await Promise.all([wczytajZapisy(), wczytajNakladke()]);
    },

    wybrana: () => kategorie.find((kategoria) => kategoria.id === wybrana) ?? null,

    wybierz(kategoria) {
      wybrana = kategoria;
      oglos();
    },

    nakladka: () => nakladka,

    odswiez: wczytaj,

    async zapisz(kategoria, tresc, tryb) {
      const wynik = await zrodlo.zapisz({
        categoryId: kategoria,
        content: tresc,
        mode: tryb,
        ...adres(),
      });
      const dokument = wynik.wynik?.document;
      if (wynik.udany && dokument !== undefined) {
        dokumenty = [
          ...dokumenty.filter((zapis) => zapis.categoryId !== dokument.categoryId),
          dokument,
        ];
        oglos();
        await wczytajNakladke();
      }
      return wynik;
    },

    async usun(dokument) {
      const wynik = await zrodlo.usun({ documentId: dokument });
      if (wynik.udany) {
        dokumenty = dokumenty.filter((zapis) => zapis.id !== dokument);
        oglos();
        await wczytajNakladke();
      }
      return wynik;
    },

    naZmiane: (sluchacz) => void sluchacze.push(sluchacz),

    rozlacz: odsubskrybuj,
  };
}
