import { pokazKomunikat } from '../aplikacja/komunikaty';

// Odpowiedź kontrolek paska bez gotowej funkcji: naciśnięcie zawsze musi dać wyjaśnienie.

/** Powód nieczynności jednej kontrolki paska górnego, wypisywany razem z tytułem sterującym treścią komunikatu. */
interface Powod {
  tytul: string;
  tresc: string;
}

/**
 * Wykaz kontrolek bez gotowej funkcji.
 *
 * Rejestr sterowany danymi: kontrolka odzyskująca funkcję znika stąd, a nie
 * dostaje gałęzi w kodzie.
 */
const POWODY = {
  powiadomienia: {
    tytul: 'Powiadomienia',
    tresc:
      'Kontrakt rdzenia nie ma ani jednej komendy powiadomień (280 komend, zero '
      + 'trafień na „notif"), więc dzwonek nie ma czego otworzyć. Licznik stoi na '
      + 'zerze, dopóki komenda nie powstanie. Materiał na wykaz byłby: dziennik '
      + 'komunikatów klienta — dziś dymek żyje sześć sekund i znika bez śladu.',
  },
  polecenie: {
    tytul: 'Pole poleceń',
    tresc:
      'Wpisana treść nie pasuje do niczego, co pole umie znaleźć, a wykonania '
      + 'polecenia nie ma: centrum poleceń nie jest zbudowane, '
      + 'a kontrakt nie niesie komendy wyszukiwania globalnego — są wyłącznie '
      + 'modułowe library.file.search i knowledge.search. Wiadomość do modelu '
      + 'wyślij z okna komunikacji.',
  },
} as const satisfies Record<string, Powod>;

/** Nazwa kontrolki paska pozostającej bez gotowej funkcji, używana jako klucz wykazu powodów jej niegotowości. */
export type NieczynnaKontrolka = keyof typeof POWODY;

/**
 * Mówi Operatorowi, dlaczego naciśnięta kontrolka nic nie zrobiła.
 *
 * Waga „ostrzeżenie", nie „błąd": nic się nie zepsuło — brakuje funkcji.
 */
export function wyjasnijNieczynnosc(kontrolka: NieczynnaKontrolka): void {
  const powod = POWODY[kontrolka];
  pokazKomunikat({ tytul: powod.tytul, tresc: powod.tresc, waga: 'ostrz' });
}

/**
 * Opis kontrolki dla technologii wspomagających i dla dymka systemowego.
 *
 * Niegotowość jest sygnalizowana opisowo, więc czytnik ekranu dowiaduje się
 * o niej razem z nazwą kontrolki — nie dopiero po naciśnięciu.
 */
export function opisNieczynnosci(kontrolka: NieczynnaKontrolka, nazwa: string): string {
  return `${nazwa} — nieczynne: ${POWODY[kontrolka].tresc}`;
}
