import { utworzSekcjeBraku } from './sekcja-braku';
import type { SekcjaUstawien } from './sekcje';

/**
 * Sekcja „Konta modeli" — pozycja odsyłająca, nie druga rama.
 *
 * Rodzina `account.*` (add, list, remove, update, default.set) jest obsłużona
 * w komplecie w `budowa/client/src/modele/`: panel kont, wykaz kont, źródło
 * kont. Drugi formularz do tych samych pięciu komend byłby drugim oknem do tych
 * samych danych rdzenia, a te rozjechałyby się przy pierwszej zmianie.
 *
 * Sekcja bez pól nie jest brakiem do uzupełnienia, tylko podziałem
 * odpowiedzialności. Warunkiem zbudowania czegokolwiek w tym miejscu jest
 * wymaganie, którego `modele/` nie potrafi unieść — osobne uprawnienia, inny
 * kontrakt albo inny odbiorca.
 *
 * Sekcja istnieje, bo projekt okna ją wymienia: Operator, który jej tu nie
 * znajdzie, szukałby jej dalej w tym oknie, a jedno zdanie z drogą jest krótsze
 * niż to szukanie i uczciwsze niż milczenie.
 */
export function utworzSekcjeKontaModeli(): SekcjaUstawien {
  return utworzSekcjeBraku({
    wstep:
      'Konta modeli — klucze dostawców, konto domyślne i przypisania — prowadzi ' +
      'panel kont w oknie modeli. Tutaj ich nie ma świadomie: te same dane ' +
      'w dwóch oknach rozjechałyby się przy pierwszej zmianie.',
    pomiar: [
      'Rodzina `account.*` liczy pięć komend (add, list, remove, update, ' +
        'default.set) i wszystkie pięć jest obsłużonych w ' +
        '`budowa/client/src/modele/`.',
      'To samo rozstrzygnięcie zdjęło z tego rejestru pięć sekcji zbudowanych ' +
        'na `account.*` i `identity.*` — uzasadnienie stoi w nagłówku ' +
        '`sekcje.ts` i nie zostaje tu powtórzone.',
    ],
    odeslanie:
      'Droga: okno modeli → panel kont. Zmiana wykonana tam obowiązuje ' +
      'wszędzie — nie ma drugiego miejsca, w którym trzeba ją powtórzyć.',
  });
}
