/**
 * Nazwy kanoniczne operacji kontekstowych modułu Developer: jedno źródło nazw
 * dla panelu akcji, żeby przycisk nie rozjechał się z komendą, którą wywołuje.
 * Panel obejmuje osiem operacji, a dziewiąta należy do paska narzędzi promptu.
 */
import { ContextualOpKind } from '../../../../shared/contract';
import { przyciskAkcji } from '../../modele/kontrolki-formularza-braki';

/**
 * Operacja kontekstowa panelu akcji: klucz stały, nazwa kanoniczna trafiająca
 * na przycisk, zakres opisujący jej działanie oraz rodzaj ze słownika
 * kontraktu, z którym jedzie żądanie do rdzenia.
 */
export interface AkcjaKanoniczna {
  /** Klucz stały — nazwa widoczna bywa zmieniana, klucz nie. */
  kod: string;
  /** Nazwa kanoniczna — ta i tylko ta trafia na przycisk. */
  nazwa: string;
  /** Co ta operacja robi. */
  zakres: string;
  /** Rodzaj operacji ze słownika kontraktu — z nim jedzie żądanie. */
  rodzaj: ContextualOpKind;
}

/**
 * Osiem operacji panelu akcji, a po nich dziewiąta — konwersja języka —
 * o innym pochodzeniu. Pierwsze cztery stoją na pasku wprost, pozostałe pięć
 * należy do rozwinięcia, a kolejność tablicy tę różnicę zachowuje.
 */
export const AKCJE_KANONICZNE: readonly AkcjaKanoniczna[] = [
  {
    kod: 'generuj',
    nazwa: 'Generuj',
    zakres: 'nowy fragment z opisu albo z kontekstu repozytorium',
    rodzaj: ContextualOpKind.Generate,
  },
  {
    kod: 'refaktoryzuj',
    nazwa: 'Refaktoryzuj',
    zakres: 'przebudowa treści bez zmiany zachowania',
    rodzaj: ContextualOpKind.Refactor,
  },
  {
    kod: 'wyjasnij',
    nazwa: 'Wyjaśnij',
    zakres: 'objaśnienie treści pliku (nazwa zlecenia: „Przejrzyj")',
    rodzaj: ContextualOpKind.Explain,
  },
  {
    kod: 'napisz-test',
    nazwa: 'Napisz test',
    zakres: 'sprawdzian dla treści pliku (nazwa zlecenia: „Wygeneruj testy")',
    rodzaj: ContextualOpKind.WriteTest,
  },
  {
    kod: 'udokumentuj',
    nazwa: 'Udokumentuj',
    zakres: 'komentarze zgodne z konwencją języka pliku',
    rodzaj: ContextualOpKind.Document,
  },
  {
    kod: 'napraw',
    nazwa: 'Napraw',
    zakres: 'poprawka wedle zgłoszenia budowania albo lintera',
    rodzaj: ContextualOpKind.Fix,
  },
  {
    kod: 'zoptymalizuj',
    nazwa: 'Zoptymalizuj',
    zakres: 'przebudowa treści pod koszt wykonania',
    rodzaj: ContextualOpKind.Optimize,
  },
  {
    kod: 'zmien-nazwe-symbolu',
    nazwa: 'Zmień nazwę symbolu globalnie',
    zakres:
      'zmiana nazwy w całym repozytorium, z podglądem; rdzeń kieruje ją do serwera języka, ' +
      'a nie do modelu — poprawność tej czynności da się rozstrzygnąć',
    rodzaj: ContextualOpKind.RenameSymbol,
  },
  {
    kod: 'konwertuj-jezyk',
    nazwa: 'Konwertuj język',
    zakres:
      'przepisanie treści na inny język albo inną bibliotekę; operacja należy do paska ' +
      'narzędzi promptu, nie do ośmiu operacji panelu akcji',
    rodzaj: ContextualOpKind.Convert,
  },
];

/**
 * Liczba pozycji stojących na pasku wprost; pozostałe należą do rozwinięcia.
 * Samego menu próg nie buduje, więc wszystkie dziewięć przycisków stoi płasko.
 */
export const PROG_PASKA = 4;

/**
 * Osadza dziewięć operacji jako przyciski wykonujące `developer.contextual.op`.
 *
 * @param gospodarz miejsce paska akcji okna.
 * @param wykonaj droga do rdzenia; okno podaje ją, bo tylko ono ma stan treści,
 *   w którym pokazuje wynik i odmowę.
 */
export function osadzAkcjeKanoniczne(
  gospodarz: HTMLElement,
  wykonaj: (akcja: AkcjaKanoniczna) => void,
): HTMLButtonElement[] {
  const przyciski = AKCJE_KANONICZNE.map((akcja) => {
    const przycisk = przyciskAkcji(akcja.nazwa, 'dn-btn dn-btn--zarys');
    przycisk.title = akcja.zakres;
    przycisk.dataset['akcja'] = akcja.kod;
    przycisk.addEventListener('click', () => wykonaj(akcja));
    return przycisk;
  });
  gospodarz.append(...przyciski);
  return przyciski;
}
