/**
 * Nazwy kanoniczne operacji kontekstowych modułu Developer.
 *
 * Jedno źródło nazw dla panelu akcji: ta sama operacja nazwana w dwóch
 * miejscach inaczej rozjeżdża przycisk z komendą, którą wywołuje. Panel
 * obejmuje osiem operacji; dziewiąta, „Konwertuj język", należy do paska
 * narzędzi promptu i ma to odnotowane we własnym zakresie.
 *
 * ── Co się tu zmieniło ─────────────────────────────────────────────────────
 * Wykaz stał wcześniej jako dziewięć przycisków bez drogi do rdzenia: kontrakt
 * niósł nazwę `developer.contextual.op`, lecz rdzeń nie miał dla niej uchwytu.
 * Rdzeń ma ją dziś, więc przyciski wołają — a przycisk tłumaczący swój brak
 * jest właściwy dokładnie do chwili, w której brak zniknie.
 *
 * ── Dlaczego zmiana nazwy symbolu nie idzie do modelu ──────────────────────
 * Zmiana nazwy w całym repozytorium jest czynnością ROZSTRZYGALNĄ i robi ją
 * serwer języka, który zna graf odwołań. Model dałby wynik prawdopodobny zamiast
 * poprawnego — i to w czynności, której poprawność da się sprawdzić. Rdzeń
 * kieruje więc `renameSymbol` do serwera języka, a panel podaje ją tą samą
 * komendą, bo dla Operatora jest to ta sama pozycja paska.
 */
import { ContextualOpKind } from '../../../../shared/contract';
import { przyciskAkcji } from '../../modele/kontrolki-formularza-braki';

/** Operacja kontekstowa: nazwa kanoniczna wraz z rodzajem z kontraktu. */
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
 * o innym pochodzeniu.
 *
 * Pierwsze cztery (`[Generuj] [Refaktoryzuj] [Wyjaśnij] [Napisz test]`) stoją
 * na pasku wprost; pozostałe pięć należy pod `[Więcej ▾]`. Kolejność tablicy tę
 * różnicę zachowuje, a `PROG_PASKA` ją nazywa.
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
 * Ile pozycji stoi na pasku wprost; reszta należy do `[Więcej ▾]`.
 *
 * Samo menu nie jest zbudowane i próg go nie buduje: jedyne menu biblioteczne
 * (`komponenty/menu-drzewo.ts`) jest sterem nastawy — na jego uchwycie stoi
 * wartość bieżąca, a wybór liścia ją zmienia. Operacje jednorazowe wartości
 * bieżącej nie mają, więc uchwyt pokazywałby napis, który niczego nie
 * odzwierciedla. Wszystkie dziewięć stoi więc płasko.
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
