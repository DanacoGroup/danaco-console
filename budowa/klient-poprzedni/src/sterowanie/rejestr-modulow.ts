import { Command, type Module } from '../../../shared/contract';
import type { Kanal } from '../protokol/kanal';
import { utworzMagistrale, type Odsubskrybuj } from '../polaczenie/magistrala-zdarzen';

/**
 * Katalog modułów rdzenia — czytająca pamięć podręczna nad `module.list`.
 *
 * Wykaz idzie z rdzenia, nie z `KnownModuleIds` kontraktu: tam stoją cztery
 * pozycje — `talkin`, `workspace`, `codestudio`, `multitaskingai` — a są to
 * środowiska, nie moduły. Katalog rdzenia niesie piętnaście modułów (tabela
 * `modul`, zasilana przez `migracja_007_zaczyn_slownikow.sql`): `studio`,
 * `library`, `agents`, `research` i pozostałe. Podstawienie wykazu środowisk
 * pod pole pytające o moduł pokazywałoby każdy prawdziwy moduł jako „spoza
 * wykazu”, a Operatorowi proponowałoby środowiska tam, gdzie pyta się o moduł.
 *
 * Rejestr niczego nie zapisuje — jest wyłącznie odczytem. Wykaz jest pusty do
 * chwili odpowiedzi rdzenia; `odpowiedzOtrzymana` rozróżnia „jeszcze nie wiem”
 * od „katalog jest pusty”, bo widok musi te dwa stany rozróżnić.
 */
export interface RejestrModulow {
  /** Moduły znane w tej chwili; pusta lista, dopóki rdzeń nie odpowie. */
  moduly(): Module[];
  /** Czy rdzeń odpowiedział na `module.list` choć raz. */
  odpowiedzOtrzymana(): boolean;
  /** Zamawia wykaz z rdzenia. */
  odswiez(): void;
  /** Subskrypcja zmian wykazu. */
  naZmiane(sluchacz: (moduly: Module[]) => void): Odsubskrybuj;
}

export function utworzRejestrModulow(kanal: Kanal): RejestrModulow {
  const zmiany = utworzMagistrale<Module[]>();
  let wykaz: Module[] = [];
  let odpowiedziano = false;

  return {
    moduly: () => wykaz,

    odpowiedzOtrzymana: () => odpowiedziano,

    odswiez() {
      kanal.wyslij(Command.ModuleList, {}, (wynik) => {
        // Odpowiedź odnotowujemy także przy odmowie: pytanie zostało
        // rozstrzygnięte, więc wskaźnik ładowania ma zgasnąć, a widok pokazać
        // stan pusty zamiast wiecznego czekania.
        odpowiedziano = true;
        const odebrane = wynik.wynik?.modules;
        wykaz = wynik.udany && odebrane !== undefined ? odebrane : wykaz;
        zmiany.oglos(wykaz);
      });
    },

    naZmiane: (sluchacz) => zmiany.subskrybuj(sluchacz),
  };
}
